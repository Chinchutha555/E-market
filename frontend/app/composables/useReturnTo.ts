/**
 * อ่านพารามิเตอร์ `return` (path ภายในระบบเท่านั้น) เพื่อให้ Breadcrumb / ปุ่มย้อนกลับ
 * พากลับไปยังหน้าต้นทางในสภาพเดิม เช่น แผนที่ที่เปิด drawer แผงค้าค้างไว้
 */
export const useReturnTo = (fallback: { label: string; to: string }) => {
  const route = useRoute()

  const returnUrl = computed(() => {
    const raw = typeof route.query.return === 'string' ? route.query.return : ''
    return raw.startsWith('/') && !raw.startsWith('//') ? raw : ''
  })

  const returnLabel = computed(() => {
    const path = returnUrl.value.split('?')[0]
    if (path === '/stalls/map') return 'แผนผังแผงค้า'
    if (path === '/stalls/manage') return 'จัดการแผงเช่า'
    return fallback.label
  })

  const backTo = computed(() => returnUrl.value || fallback.to)
  const crumbItems = computed(() => [{ label: returnLabel.value, to: backTo.value }])

  return { returnUrl, backTo, crumbItems }
}
