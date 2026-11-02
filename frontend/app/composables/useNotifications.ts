import { computed, watch } from 'vue'

export interface INotificationItem {
  id: number
  recipient_role: 'tenant' | 'officer'
  recipient_user_id?: number
  recipient_name: string
  type: string
  status_code: string
  status_label: string
  title: string
  message: string
  booking_id?: number
  booking_number?: string
  stall_number?: string
  market_name?: string
  project_name?: string
  target_url?: string
  is_read: boolean
  created_at: string
}

type NotificationCursor = { lastTopId: number; initialized: boolean }
const notificationCursors = new Map<string, NotificationCursor>()

export const useNotifications = () => {
  const { get, post } = useApi()
  const { user, token, role: authRole } = useAuth()

  const notifications = useState<INotificationItem[]>('global_notifications_list', () => [])
  const tenantUnreadCount = useState<number>('global_notifications_tenant_unread', () => 0)
  const officerUnreadCount = useState<number>('global_notifications_officer_unread', () => 0)
  const loading = useState<boolean>('global_notifications_loading', () => false)
  const toastQueue = useState<INotificationItem[]>('global_notifications_toast_queue', () => [])
  const notificationOwnerKey = useState<string>('global_notifications_owner_key', () => '')

  const activeRole = computed<'tenant' | 'officer'>(() => {
    const currentAuthRole = String(authRole.value || user.value?.role || '').toLowerCase()
    return currentAuthRole === 'tenant' ? 'tenant' : 'officer'
  })
  const sessionUserId = computed<number | null>(() => {
    try {
      const payload = token.value?.split('.')[1]
      if (!payload) return null
      const claims = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
      const id = Number(claims.user_id)
      return Number.isSafeInteger(id) && id > 0 ? id : null
    } catch { return null }
  })

  const sessionKey = computed(() =>
    sessionUserId.value ? String(sessionUserId.value) + ':' + activeRole.value : 'anonymous'
  )
  const clearClientNotificationState = () => {
    notifications.value = []
    tenantUnreadCount.value = 0
    officerUnreadCount.value = 0
    toastQueue.value = []
  }

  // Nuxt useState survives client-side navigation. Scope that state to the
  // authenticated account so logout/login and role changes cannot briefly show
  // another user's notifications or reuse their "latest notification" cursor.
  watch(sessionKey, (nextKey, previousKey) => {
    if (notificationOwnerKey.value === nextKey) return
    clearClientNotificationState()
    if (nextKey === 'anonymous' && previousKey && previousKey !== 'anonymous') {
      notificationCursors.delete(previousKey)
    }
    notificationOwnerKey.value = nextKey
  }, { immediate: true })

  const isNotificationVisibleToActiveUser = (item: INotificationItem) => {
    const currentId = sessionUserId.value
    return currentId != null && item.recipient_user_id != null &&
      String(item.recipient_user_id) === String(currentId)
  }
  const dismissToast = (id: number) => {
    toastQueue.value = toastQueue.value.filter(t => t.id !== id)
  }

  const pushToast = (items: INotificationItem[]) => {
    for (const item of items) {
      if (!toastQueue.value.some(t => t.id === item.id)) {
        toastQueue.value.unshift(item)
        if (toastQueue.value.length > 3) {
          toastQueue.value = toastQueue.value.slice(0, 3)
        }
        if (typeof window !== 'undefined') {
          window.setTimeout(() => {
            dismissToast(item.id)
          }, 7000)
        }
      }
    }
  }

  const fetchNotifications = async (showNewAsToast = true) => {
    const requestedSessionKey = sessionKey.value
    if (requestedSessionKey === 'anonymous') {
      clearClientNotificationState()
      return
    }
    loading.value = true
    try {
      const res: any = await get('/notifications', { role: 'all', _t: Date.now() })
      const list: INotificationItem[] = Array.isArray(res?.data)
        ? res.data
        : Array.isArray(res)
          ? res
          : []
      const visibleList = list.filter(isNotificationVisibleToActiveUser).map(item => ({ ...item, recipient_role: (activeRole.value === 'tenant' ? 'tenant' : 'officer') as 'tenant' | 'officer' }))
      // Ignore a late response belonging to the account that just logged out.
      if (requestedSessionKey !== sessionKey.value) return

      const cursor = notificationCursors.get(requestedSessionKey) || { lastTopId: 0, initialized: false }

      if (cursor.initialized && showNewAsToast && visibleList.length > 0 && cursor.lastTopId > 0) {
        const newlyAdded: INotificationItem[] = []
        for (const n of visibleList) {
          if (n.id === cursor.lastTopId) break
          if (!n.is_read) {
            newlyAdded.push(n)
          }
        }
        if (newlyAdded.length > 0) {
          pushToast(newlyAdded.reverse())
        }
      }

      if (visibleList.length > 0) {
        cursor.lastTopId = visibleList[0].id
      }
      cursor.initialized = true
      notificationCursors.set(requestedSessionKey, cursor)

      notifications.value = visibleList
      tenantUnreadCount.value = visibleList.filter((n: INotificationItem) => n.recipient_role === 'tenant' && !n.is_read).length
      officerUnreadCount.value = visibleList.filter((n: INotificationItem) => n.recipient_role === 'officer' && !n.is_read).length
    } catch (e) {
      console.error('Failed to fetch notifications:', e)
    } finally {
      loading.value = false
    }
  }

  const markAsRead = async (id: number) => {
    const target = notifications.value.find(n => n.id === id)
    const wasUnread = Boolean(target && !target.is_read)
    if (target && !target.is_read) {
      target.is_read = true
      if (target.recipient_role === 'tenant') {
        tenantUnreadCount.value = Math.max(0, tenantUnreadCount.value - 1)
      } else {
        officerUnreadCount.value = Math.max(0, officerUnreadCount.value - 1)
      }
    }
    dismissToast(id)
    try {
      await post(`/notifications/${id}/read`, { user_id: (user.value as any)?.id, role: activeRole.value })
    } catch (error) {
      // Do not leave a false read state when persistence failed.
      if (target && wasUnread) target.is_read = false
      await fetchNotifications(false)
      console.error('Failed to persist notification read state:', error)
    }
  }

  const markAllAsRead = async (roleFilter?: 'tenant' | 'officer' | 'all') => {
    const targetRole = roleFilter || activeRole.value
    const previousReadState = new Map(notifications.value.map(n => [n.id, n.is_read]))
    notifications.value.forEach((n) => {
      if (targetRole === 'all' || n.recipient_role === targetRole) {
        n.is_read = true
      }
    })
    if (targetRole === 'all' || targetRole === 'tenant') tenantUnreadCount.value = 0
    if (targetRole === 'all' || targetRole === 'officer') officerUnreadCount.value = 0

    try {
      await post('/notifications/read-all', { role: targetRole, user_id: (user.value as any)?.id })
    } catch (error) {
      notifications.value.forEach((n) => {
        n.is_read = previousReadState.get(n.id) ?? n.is_read
      })
      await fetchNotifications(false)
      console.error('Failed to persist notification read-all state:', error)
    }
  }

  const tenantNotifications = computed(() =>
    notifications.value.filter(n => n.recipient_role === 'tenant')
  )

  const officerNotifications = computed(() =>
    notifications.value.filter(n => n.recipient_role === 'officer')
  )

  const activeRoleUnreadCount = computed(() =>
    activeRole.value === 'tenant' ? tenantUnreadCount.value : officerUnreadCount.value
  )

  const totalUnreadCount = computed(() =>
    tenantUnreadCount.value + officerUnreadCount.value
  )

  return {
    notifications,
    tenantNotifications,
    officerNotifications,
    tenantUnreadCount,
    officerUnreadCount,
    activeRole,
    activeRoleUnreadCount,
    totalUnreadCount,
    toastQueue,
    loading,
    fetchNotifications,
    markAsRead,
    markAllAsRead,
    dismissToast
  }
}
