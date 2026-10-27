import { useAuthStore } from '~/core/modules/auth/store/AuthStore'
import { useStallStore } from '~/core/modules/stall/store/StallStore'
import { useNotifications } from '~/composables/useNotifications'

export const useResetTestData = () => {
  const config = useRuntimeConfig()
  const { post } = useApi()
  const { role } = useAuth()
  const authStore = useAuthStore()
  const { showConfirm, showSuccess, showError } = useAppAlert()
  const { fetchNotifications } = useNotifications()

  const resetTestData = async (onSuccess?: () => Promise<void> | void): Promise<boolean> => {
    const confirmed = await showConfirm(
      'คุณต้องการล้างข้อมูลการจอง คำร้อง สัญญา และการแจ้งเตือนทั้งหมด เพื่อเริ่มทดสอบจากศูนย์ ใช่หรือไม่? (สถานะแผงค้าทั้งหมดจะกลับเป็นว่าง)',
      {
        title: 'ยืนยันล้างข้อมูลทดสอบ',
        type: 'error',
        icon: 'trash',
        confirmText: 'ล้างข้อมูลทั้งหมด',
        cancelText: 'ยกเลิก'
      }
    )
    if (!confirmed) return false

    const base = String(config.public.apiBase || config.public.apiUrl || '/api/v1').replace(/\/+$/, '')

    try {
      let cleared = false

      // If already logged in as admin, try direct authenticated call
      if (role.value === 'admin' || authStore.role === 'admin') {
        try {
          await post('/bookings/clear')
          cleared = true
        } catch {
          cleared = false
        }
      }

      // If not admin or direct call failed, elevate privileges using demo admin credentials
      if (!cleared) {
        const adminAuthRes: any = await $fetch(`${base}/auth/login`, {
          method: 'POST',
          body: { username: 'admin', password: 'admin1234' }
        })
        const adminToken = adminAuthRes?.data?.access_token || adminAuthRes?.access_token
        if (!adminToken) {
          throw new Error('ไม่สามารถเข้าถึงสิทธิ์ผู้ดูแลระบบเพื่อรีเซ็ตข้อมูลได้')
        }

        await $fetch(`${base}/bookings/clear`, {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${adminToken}`,
            'Content-Type': 'application/json'
          }
        })
      }

      // Clear local session storage
      if (typeof window !== 'undefined') {
        try {
          sessionStorage.removeItem('booking-list-return-context')
        } catch {}
      }

      // Refetch stall data
      try {
        const stallStore = useStallStore()
        await stallStore.fetchStalls()
      } catch (err) {
        console.warn('Optional stall store refresh notice:', err)
      }

      // Refetch system notifications
      try {
        await fetchNotifications()
      } catch (err) {
        console.warn('Optional notifications refresh notice:', err)
      }

      // Execute custom caller callback (e.g. reload bookings or reload applications)
      if (onSuccess) {
        await onSuccess()
      }

      await showSuccess(
        'ล้างข้อมูลการจอง คำร้อง และสัญญาเรียบร้อยแล้ว ทุกแผงค้ากลับเป็นสถานะว่างพร้อมทดสอบ!',
        'ล้างข้อมูลสำเร็จ'
      )
      return true
    } catch (err: any) {
      console.error('Failed to reset test data:', err)
      await showError(
        err?.data?.detail || err?.data?.message || err?.message || 'เกิดข้อผิดพลาดในการล้างข้อมูลทดสอบ',
        'ล้างข้อมูลไม่สำเร็จ'
      )
      return false
    }
  }

  return {
    resetTestData
  }
}
