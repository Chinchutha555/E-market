<template>
  <div ref="notifWrapperRef" class="notification-wrapper">
    <!-- Notification Bell Button -->
    <button
      type="button"
      class="notification-bell-btn"
      :class="{ active: showNotifications, 'has-unread': currentBellUnreadCount > 0 }"
      @click.stop="toggleNotifications"
      title="การแจ้งเตือนและการนัดหมาย"
    >
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
        <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
      </svg>
      <span v-if="currentBellUnreadCount > 0" class="notification-count-badge">
        {{ currentBellUnreadCount > 99 ? '99+' : currentBellUnreadCount }}
      </span>
    </button>

    <!-- Notification Panel Popover (บน Desktop เป็น Popover Dropdown / บน Mobile เป็น Teleport หน้าจอเต็ม) -->
    <Teleport to="body" :disabled="!isMobile">
      <div
        v-if="showNotifications"
        ref="notifCardRef"
        class="notification-dropdown-card"
        :class="{ 'is-mobile-fullscreen': isMobile }"
        @click.stop
      >
        <!-- Panel Header -->
        <div class="notif-panel-header">
          <div class="notif-header-top">
            <!-- Mobile Back Button: แสดงเฉพาะบน Responsive เพื่อปิดกลับหน้าเดิม -->
            <button
              v-if="isMobile"
              type="button"
              class="notif-mobile-back-btn"
              @click="showNotifications = false"
              aria-label="ย้อนกลับ"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="15 18 9 12 15 6"></polyline>
              </svg>
              <span>ย้อนกลับ</span>
            </button>

            <div class="notif-header-title-group">
              <span class="notif-panel-title">การแจ้งเตือนสถานะ</span>
              <span v-if="displayedUnreadCount > 0" class="notif-unread-pill">
                ใหม่ {{ displayedUnreadCount }}
              </span>
            </div>
            <AppButton
              v-if="displayedUnreadCount > 0"
              variant="outline"
              size="sm"
              @click="handleMarkAllRead"
            >
              <template #prefix>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </template>
              <span>อ่านทั้งหมด</span>
            </AppButton>
          </div>

          <!-- Sub-filter: All vs Unread using AppSegmentedControl -->
          <div class="notif-filter-bar">
            <AppSegmentedControl
              v-model="notifFilter"
              :options="notifFilterOptions"
              aria-label="กรองรายการแจ้งเตือน"
            />
          </div>
        </div>

        <!-- Notification Items List -->
        <div class="notif-list-body">
          <div v-if="displayedNotifications.length === 0" class="notif-empty-state">
            <div class="notif-empty-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
                <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
              </svg>
            </div>
            <div class="notif-empty-title">ไม่มีรายการแจ้งเตือน</div>
            <div class="notif-empty-sub">
              {{ notifFilter === 'unread' ? 'คุณได้อ่านการแจ้งเตือนทั้งหมดเรียบร้อยแล้ว' : 'ยังไม่มีความเคลื่อนไหวของสถานะในขณะนี้' }}
            </div>
          </div>

          <button
            v-for="item in displayedNotifications"
            :key="item.id"
            type="button"
            class="notif-item-row"
            :class="{ unread: !item.is_read }"
            @click="handleNotificationClick(item)"
          >
            <!-- Status Icon Circle -->
            <div class="notif-status-icon" :class="getStatusToneClass(item.status_code)">
              <!-- Approved / Completed -->
              <svg v-if="item.status_code === 'approved' || item.status_code === 'completed'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                <polyline points="22 4 12 14.01 9 11.01"></polyline>
              </svg>
              <!-- Interested / Lease Applied -->
              <svg v-else-if="item.status_code === 'interested' || item.status_code === 'lease_applied'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
              </svg>
              <!-- Waiting Queue / Pending Slot -->
              <svg v-else-if="item.status_code === 'waiting_queue' || item.status_code === 'pending_slot'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                <circle cx="9" cy="7" r="4"></circle>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
              </svg>
              <!-- Rejected / Cancelled / No Show -->
              <svg v-else-if="['rejected', 'cancelled', 'no_show', 'not_interested', 'waitlist_cancelled'].includes(item.status_code)" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="15" y1="9" x2="9" y2="15"></line>
                <line x1="9" y1="9" x2="15" y2="15"></line>
              </svg>
              <!-- Default: Pending Approval / Calendar Clock -->
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
              </svg>
            </div>

            <!-- Content -->
            <div class="notif-item-content">
              <div class="notif-item-top-meta">
                <span class="notif-status-badge" :class="getStatusToneClass(item.status_code)">
                  {{ item.status_label }}
                </span>
                <span class="notif-time-text">{{ formatRelativeTime(item.created_at) }}</span>
              </div>

              <div class="notif-item-title">{{ formatNotifText(item.title) }}</div>
              <div class="notif-item-desc">{{ formatNotifText(item.message) }}</div>

              <div v-if="item.market_name || item.stall_number" class="notif-item-tags">
                <span v-if="item.stall_number" class="notif-meta-chip stall">แผง {{ item.stall_number }}</span>
                <span v-if="item.market_name" class="notif-meta-chip market">{{ item.market_name }}</span>
              </div>
            </div>

            <!-- Unread Dot -->
            <span v-if="!item.is_read" class="notif-unread-dot" title="ยังไม่อ่าน"></span>
          </button>
        </div>
      </div>
    </Teleport>

    <!-- Real-time Status Notification Toast Stack (มุมบนขวา เด้งแจ้งเตือนทันทีที่สถานะเปลี่ยน) -->
    <Teleport to="body">
      <div v-if="toastQueue.length > 0" class="notif-toast-stack" aria-live="polite">
        <TransitionGroup name="notif-toast-anim">
          <div
            v-for="toast in toastQueue"
            :key="toast.id"
            class="notif-toast-card"
            :class="getStatusToneClass(toast.status_code)"
            @click="handleNotificationClick(toast)"
          >
            <div class="notif-toast-icon" :class="getStatusToneClass(toast.status_code)">
              <svg v-if="toast.status_code === 'approved' || toast.status_code === 'completed'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                <polyline points="22 4 12 14.01 9 11.01"></polyline>
              </svg>
              <svg v-else-if="toast.status_code === 'interested' || toast.status_code === 'lease_applied'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
              </svg>
              <svg v-else-if="['rejected', 'cancelled', 'no_show', 'not_interested', 'waitlist_cancelled'].includes(toast.status_code)" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="15" y1="9" x2="9" y2="15"></line>
                <line x1="9" y1="9" x2="15" y2="15"></line>
              </svg>
              <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
                <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
              </svg>
            </div>

            <div class="notif-toast-body">
              <div class="notif-toast-meta">
                <span class="notif-toast-role-pill" :class="toast.recipient_role">
                  {{ toast.recipient_role === 'tenant' ? 'แจ้งเตือนผู้เช่า' : 'แจ้งเตือนเจ้าหน้าที่' }}
                </span>
                <span class="notif-status-badge" :class="getStatusToneClass(toast.status_code)">
                  {{ toast.status_label }}
                </span>
              </div>
              <div class="notif-toast-title">{{ formatNotifText(toast.title) }}</div>
              <div class="notif-toast-message">{{ formatNotifText(toast.message) }}</div>
            </div>

            <button
              type="button"
              class="notif-toast-close"
              @click.stop="dismissToast(toast.id)"
              aria-label="ปิดการแจ้งเตือน"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
        </TransitionGroup>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useNotifications, type INotificationItem } from '~/composables/useNotifications'

const emit = defineEmits<{
  (e: 'open'): void
}>()

const router = useRouter()
const route = useRoute()

watch(() => route.fullPath, () => {
  showNotifications.value = false
})
const {
  tenantNotifications,
  officerNotifications,
  tenantUnreadCount,
  officerUnreadCount,
  activeRole,
  toastQueue,
  fetchNotifications,
  markAsRead,
  markAllAsRead,
  dismissToast
} = useNotifications()

const showNotifications = ref(false)
const notifFilter = ref<'all' | 'unread'>('all')
const notifWrapperRef = ref<HTMLElement | null>(null)
const notifCardRef = ref<HTMLElement | null>(null)
const isMobile = ref(false)

const checkMobile = () => {
  if (import.meta.client) {
    isMobile.value = window.innerWidth <= 768
  }
}

const currentBellUnreadCount = computed(() => {
  return activeRole.value === 'tenant' ? tenantUnreadCount.value : officerUnreadCount.value
})

const roleFilteredList = computed(() => {
  return activeRole.value === 'tenant' ? tenantNotifications.value : officerNotifications.value
})

const displayedUnreadCount = computed(() => {
  return activeRole.value === 'tenant' ? tenantUnreadCount.value : officerUnreadCount.value
})

const displayedNotifications = computed(() => {
  if (notifFilter.value === 'unread') {
    return roleFilteredList.value.filter(n => !n.is_read)
  }
  return roleFilteredList.value
})

const notifFilterOptions = computed(() => [
  { value: 'all', label: `ทั้งหมด (${roleFilteredList.value.length})` },
  { value: 'unread', label: `ยังไม่อ่าน (${displayedUnreadCount.value})` }
])

const formatNotifText = (text: string) => {
  if (!text) return ''
  const thaiDays = ['อาทิตย์', 'จันทร์', 'อังคาร', 'พุธ', 'พฤหัสบดี', 'ศุกร์', 'เสาร์']
  const thaiMonths = ['ม.ค.', 'ก.พ.', 'มี.ค.', 'เม.ย.', 'พ.ค.', 'มิ.ย.', 'ก.ค.', 'ส.ค.', 'ก.ย.', 'ต.ค.', 'พ.ย.', 'ธ.ค.']
  return text.replace(/(?:วันที่\s*|วัน)?(\d{4})-(\d{2})-(\d{2})/g, (_match, yyyy, mm, dd) => {
    const d = new Date(Number(yyyy), Number(mm) - 1, Number(dd))
    if (Number.isNaN(d.getTime())) return _match
    return `วัน${thaiDays[d.getDay()]}ที่ ${d.getDate()} ${thaiMonths[d.getMonth()]} ${d.getFullYear() + 543}`
  })
}

const toggleNotifications = async () => {
  showNotifications.value = !showNotifications.value
  if (showNotifications.value) {
    emit('open')
    await fetchNotifications(false)
  }
}

const closeNotifications = () => {
  showNotifications.value = false
}

defineExpose({
  closeNotifications
})

const handleMarkAllRead = async () => {
  await markAllAsRead(activeRole.value)
}

const handleNotificationClick = async (item: INotificationItem) => {
  await markAsRead(item.id)
  showNotifications.value = false

  if ((item.status_code === 'waitlist_offered' || item.status_code === 'waitlist_expiring_soon') && item.stall_number) {
    router.push({
      path: '/stalls/book',
      query: {
        market: item.market_name || 'ตลาดนัดจตุจักร',
        stall: item.stall_number,
        step: '1',
        claimWaitlist: '1'
      }
    })
  } else if (item.booking_id) {
    router.push({ path: `/bookings/${item.booking_id}` })
  } else if (item.target_url) {
    router.push(item.target_url)
  } else {
    router.push({ path: '/bookings' })
  }
}

const getStatusToneClass = (statusCode: string) => {
  switch (statusCode) {
    case 'approved':
    case 'completed':
    case 'waitlist_offered':
      return 'tone-green'
    case 'interested':
    case 'lease_applied':
      return 'tone-blue'
    case 'pending_slot':
    case 'waiting_queue':
      return 'tone-purple'
    case 'rejected':
    case 'cancelled':
    case 'no_show':
    case 'not_interested':
    case 'waitlist_expired':
    case 'waitlist_cancelled':
      return 'tone-red'
    case 'waitlist_expiring_soon':
    case 'pending_approval':
    default:
      return 'tone-amber'
  }
}

const formatRelativeTime = (iso: string) => {
  if (!iso) return ''
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  const diffSec = Math.floor((Date.now() - date.getTime()) / 1000)
  if (diffSec < 60) return 'เมื่อสักครู่'
  const diffMin = Math.floor(diffSec / 60)
  if (diffMin < 60) return `${diffMin} นาทีที่แล้ว`
  const diffHr = Math.floor(diffMin / 60)
  if (diffHr < 24) return `${diffHr} ชั่วโมงที่แล้ว`
  const thaiMonths = ['ม.ค.', 'ก.พ.', 'มี.ค.', 'เม.ย.', 'พ.ค.', 'มิ.ย.', 'ก.ค.', 'ส.ค.', 'ก.ย.', 'ต.ค.', 'พ.ย.', 'ธ.ค.']
  return `${date.getDate()} ${thaiMonths[date.getMonth()]} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')} น.`
}

const handleClickOutside = (e: MouseEvent) => {
  if (isMobile.value) return // บน mobile หน้าต่างเป็นแบบเต็มหน้าจอ ไม่ปิดด้วย click outside แต่มีปุ่มย้อนกลับ
  if (
    notifWrapperRef.value &&
    !notifWrapperRef.value.contains(e.target as Node) &&
    (!notifCardRef.value || !notifCardRef.value.contains(e.target as Node))
  ) {
    showNotifications.value = false
  }
}

onMounted(() => {
  fetchNotifications(false)
  checkMobile()
  if (import.meta.client) {
    window.addEventListener('resize', checkMobile, { passive: true })
    window.addEventListener('click', handleClickOutside)
  }
})

onBeforeUnmount(() => {
  if (import.meta.client) {
    window.removeEventListener('resize', checkMobile)
    window.removeEventListener('click', handleClickOutside)
  }
})
</script>

<style scoped>
.notification-wrapper {
  position: relative;
  display: inline-flex;
  align-items: center;
}

.notification-bell-btn {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  color: #475569;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  transition: all 0.2s;
}

.notification-bell-btn:hover {
  background: #f8fafc;
  color: #0f172a;
}

.notification-bell-btn.has-unread {
  border-color: #bfdbfe;
  background: #f8fbff;
  color: #1e40af;
}

.notification-bell-btn.active {
  background: #eff6ff;
  color: #2563eb;
  border-color: #93c5fd;
}

.notification-count-badge {
  position: absolute;
  top: -7px;
  right: -8px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 999px;
  background: #ef4444;
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 2px solid #ffffff;
  line-height: 1;
  box-shadow: 0 2px 6px rgba(239, 68, 68, 0.35);
}

.notification-dropdown-card {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  width: 410px;
  max-width: calc(100vw - 24px);
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  box-shadow: 0 16px 40px rgba(15, 23, 42, 0.16);
  z-index: 100;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.notif-panel-header {
  padding: 14px 16px 12px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.notif-header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.notif-header-title-group {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.notif-panel-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: #0f172a;
}

.notif-unread-pill {
  background: #fee2e2;
  color: #dc2626;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 999px;
}

.notif-filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.notif-list-body {
  max-height: 380px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.notif-item-row {
  width: 100%;
  text-align: left;
  border: none;
  border-bottom: 1px solid #f1f5f9;
  background: #ffffff;
  padding: 12px 16px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
  cursor: pointer;
  transition: background 0.15s ease;
  position: relative;
  font-family: inherit;
}

.notif-item-row:hover {
  background: #f8fafc;
}

.notif-item-row.unread {
  background: #f0f7ff;
}

.notif-item-row.unread:hover {
  background: #e6f2ff;
}

.notif-status-icon {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 2px;
}

.notif-status-icon.tone-green {
  background: #dcfce7;
  color: #15803d;
}

.notif-status-icon.tone-blue {
  background: #dbeafe;
  color: #2563eb;
}

.notif-status-icon.tone-amber {
  background: #fef3c7;
  color: #d97706;
}

.notif-status-icon.tone-purple {
  background: #f3e8ff;
  color: #7c3aed;
}

.notif-status-icon.tone-red {
  background: #fee2e2;
  color: #dc2626;
}

.notif-item-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.notif-item-top-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.notif-status-badge {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
  line-height: 1.25;
}

.notif-status-badge.tone-green {
  background: #dcfce7;
  color: #15803d;
}

.notif-status-badge.tone-blue {
  background: #dbeafe;
  color: #1d4ed8;
}

.notif-status-badge.tone-amber {
  background: #fef3c7;
  color: #b45309;
}

.notif-status-badge.tone-purple {
  background: #f3e8ff;
  color: #6d28d9;
}

.notif-status-badge.tone-red {
  background: #fee2e2;
  color: #dc2626;
}

.notif-time-text {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 500;
  white-space: nowrap;
}

.notif-item-title {
  font-size: 0.875rem;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.35;
}

.notif-item-desc {
  font-size: 0.8125rem;
  color: #475569;
  font-weight: 500;
  line-height: 1.45;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notif-item-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 2px;
}

.notif-meta-chip {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
}

.notif-meta-chip.stall {
  background: #f1f5f9;
  color: #334155;
}

.notif-meta-chip.market {
  background: #f8fafc;
  color: #64748b;
  border: 1px solid #e2e8f0;
}

.notif-unread-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #2563eb;
  flex-shrink: 0;
  margin-top: 8px;
}

.notif-empty-state {
  padding: 36px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.notif-empty-icon {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: #f1f5f9;
  color: #94a3b8;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 4px;
}

.notif-empty-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: #334155;
}

.notif-empty-sub {
  font-size: 0.8125rem;
  color: #64748b;
  font-weight: 500;
}

/* Real-time Toast Banner Stack */
.notif-toast-stack {
  position: fixed;
  top: 76px;
  right: 24px;
  z-index: 1200;
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 390px;
  max-width: calc(100vw - 32px);
  pointer-events: none;
}

.notif-toast-card {
  pointer-events: auto;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-left: 4px solid #2563eb;
  border-radius: 14px;
  padding: 14px 16px;
  box-shadow: 0 14px 34px rgba(15, 23, 42, 0.16);
  display: flex;
  align-items: flex-start;
  gap: 12px;
  cursor: pointer;
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.notif-toast-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 18px 38px rgba(15, 23, 42, 0.2);
}

.notif-toast-card.tone-green { border-left-color: #16a34a; }
.notif-toast-card.tone-blue { border-left-color: #2563eb; }
.notif-toast-card.tone-amber { border-left-color: #d97706; }
.notif-toast-card.tone-purple { border-left-color: #7c3aed; }
.notif-toast-card.tone-red { border-left-color: #dc2626; }

.notif-toast-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notif-toast-icon.tone-green { background: #dcfce7; color: #15803d; }
.notif-toast-icon.tone-blue { background: #dbeafe; color: #2563eb; }
.notif-toast-icon.tone-amber { background: #fef3c7; color: #d97706; }
.notif-toast-icon.tone-purple { background: #f3e8ff; color: #7c3aed; }
.notif-toast-icon.tone-red { background: #fee2e2; color: #dc2626; }

.notif-toast-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.notif-toast-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.notif-toast-role-pill {
  font-size: 0.75rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 999px;
}

.notif-toast-role-pill.tenant {
  background: #ecfdf5;
  color: #059669;
  border: 1px solid #a7f3d0;
}

.notif-toast-role-pill.officer {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
}

.notif-toast-title {
  font-size: 0.875rem;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.35;
}

.notif-toast-message {
  font-size: 0.8125rem;
  color: #475569;
  font-weight: 500;
  line-height: 1.4;
}

.notif-toast-close {
  border: none;
  background: transparent;
  color: #94a3b8;
  width: 26px;
  height: 26px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
}

.notif-toast-close:hover {
  background: #f1f5f9;
  color: #475569;
}

.notif-toast-anim-enter-active,
.notif-toast-anim-leave-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.notif-toast-anim-enter-from,
.notif-toast-anim-leave-to {
  opacity: 0;
  transform: translateX(24px) scale(0.96);
}

.notif-mobile-back-btn {
  display: none;
}

/* ================================================================ */
/* โหมดเต็มหน้าจอบน Responsive (Mobile Screen Full Page View)        */
/* ================================================================ */
@media (max-width: 768px) {
  .notif-mobile-back-btn {
    display: inline-flex !important;
    align-items: center !important;
    gap: 6px !important;
    min-height: 40px !important;
    height: 40px !important;
    padding: 0 14px 0 10px !important;
    border-radius: 10px !important;
    border: 1px solid #bfdbfe !important;
    background: #eff6ff !important;
    color: #2563eb !important;
    font-size: 0.875rem !important;
    font-weight: 600 !important;
    cursor: pointer !important;
    flex-shrink: 0 !important;
    -webkit-tap-highlight-color: transparent !important;
  }

  .notif-mobile-back-btn:hover {
    background: #dbeafe !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen {
    position: fixed !important;
    top: 58px !important;
    left: 0 !important;
    right: 0 !important;
    bottom: 0 !important;
    width: 100vw !important;
    max-width: 100vw !important;
    height: calc(100vh - 58px) !important;
    height: calc(100dvh - 58px) !important;
    max-height: calc(100dvh - 58px) !important;
    border-radius: 0 !important;
    border: none !important;
    box-shadow: none !important;
    z-index: 32 !important;
    display: flex !important;
    flex-direction: column !important;
    background: #f8fafc !important;
    overflow: hidden !important;
    overscroll-behavior: contain !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-panel-header {
    flex: 0 0 auto !important;
    padding: 14px 16px 12px !important;
    background: #ffffff !important;
    border-bottom: 1px solid #e2e8f0 !important;
    box-shadow: 0 1px 4px rgba(15, 23, 42, 0.04) !important;
    display: flex !important;
    flex-direction: column !important;
    gap: 12px !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-header-top {
    display: flex !important;
    align-items: center !important;
    justify-content: space-between !important;
    gap: 10px !important;
    width: 100% !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-header-title-group {
    display: inline-flex !important;
    align-items: center !important;
    gap: 8px !important;
    flex: 1 1 auto !important;
    min-width: 0 !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-panel-title {
    font-size: 1.0625rem !important;
    font-weight: 700 !important;
    color: #0f172a !important;
    white-space: nowrap !important;
    overflow: hidden !important;
    text-overflow: ellipsis !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-unread-pill {
    font-size: 0.75rem !important;
    padding: 2px 8px !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-filter-bar {
    width: 100% !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-filter-bar :deep(.app-segmented-control) {
    width: 100% !important;
    background: #f1f5f9 !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-filter-bar :deep(.app-segmented-option) {
    flex: 1 !important;
    min-height: 38px !important;
    font-size: 0.875rem !important;
    font-weight: 600 !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-list-body {
    flex: 1 1 auto !important;
    max-height: none !important;
    height: auto !important;
    overflow-y: auto !important;
    -webkit-overflow-scrolling: touch !important;
    padding: 12px 14px calc(24px + env(safe-area-inset-bottom, 0px)) !important;
    display: flex !important;
    flex-direction: column !important;
    gap: 10px !important;
    background: #f8fafc !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-item-row {
    padding: 14px 16px !important;
    gap: 12px !important;
    border-radius: 14px !important;
    border: 1px solid #eef2f6 !important;
    background: #ffffff !important;
    box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03) !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-item-row.unread {
    background: #eff6ff !important;
    border-color: #bfdbfe !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-status-icon {
    width: 38px !important;
    height: 38px !important;
    border-radius: 12px !important;
    flex-shrink: 0 !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-item-title {
    font-size: 0.9375rem !important;
    line-height: 1.4 !important;
    font-weight: 700 !important;
    color: #0f172a !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-item-desc {
    font-size: 0.875rem !important;
    line-height: 1.45 !important;
    -webkit-line-clamp: 3 !important;
    color: #475569 !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-time-text {
    font-size: 0.75rem !important;
  }

  .notification-dropdown-card.is-mobile-fullscreen .notif-empty-state {
    padding: 60px 20px !important;
  }
}
</style>
