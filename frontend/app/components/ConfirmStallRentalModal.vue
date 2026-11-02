<template>
  <Teleport to="body">
    <div v-if="open" class="rent-overlay" @click.self="$emit('cancel')">
      <div class="rent-dialog" role="dialog" aria-modal="true" aria-labelledby="rent-confirm-title">
        <AppModalHeader title="ยืนยันการเช่า" :subtitle="'แผง ' + (stall || '-') + ' · ' + (booking?.tenant_name || booking?.applicant_name || '-')" icon="document" @close="$emit('cancel')" />
        <div class="rent-body"><div class="rent-fields">
          <p v-if="submitted && !stall" class="rent-error">ไม่พบแผงที่ต้องการยืนยันการเช่า</p>
          <label>ประเภทการค้า <span class="req-star" aria-hidden="true">*</span></label>
          <p v-if="typesLoading" class="rent-message">กำลังโหลดประเภทการค้า...</p>
          <div v-else-if="typesError" class="rent-message" role="alert">
            <p class="rent-error">{{ typesError }}</p>
            <AppButton size="sm" variant="outline" @click="loadTypes">ลองอีกครั้ง</AppButton>
          </div>
          <p v-else-if="!tradeTypes.length" class="rent-message">ยังไม่มีประเภทการค้า กรุณาเพิ่มประเภทในการตั้งค่าของตลาดก่อน</p>
          <div v-else class="rent-choices">
            <AppCheckbox v-for="item in tradeTypes" :key="item.TradeType" :model-value="rentType === item.TradeType" :label="(item.Title || legacyLabels[item.TradeType] || item.TradeType) + (item.RequiresContract === false ? ' (ไม่ต้องทำสัญญา)' : ' (ต้องทำสัญญา)')" @update:model-value="rentType = $event ? item.TradeType : ''" />
          </div>
          <p v-if="submitted && !rentType && !typesLoading && !typesError && tradeTypes.length" class="rent-error">กรุณาเลือกประเภทการค้า</p>
        </div></div>
        <div class="rent-actions"><AppButton variant="secondary" size="md" @click="$emit('cancel')">ยกเลิก</AppButton><AppButton variant="primary" size="md" :loading="loading" :disabled="typesLoading || !!typesError || !tradeTypes.length" @click="confirmRental">ยืนยันการเช่า</AppButton></div>
      </div>
    </div>
  </Teleport>
</template>
<script setup lang="ts">
const props = defineProps<{ open: boolean; booking: any; loading?: boolean }>()
const emit = defineEmits<{ cancel: []; confirm: [value: { stallIdent: string; rentType: string }] }>()
const { get } = useApi()
const stall = ref('')
const rentType = ref('')
const submitted = ref(false)
const typesLoading = ref(false)
const typesError = ref('')
const tradeTypes = ref<Array<{ TradeType: string; Title: string; RequiresContract: boolean | null }>>([])
const legacyLabels: Record<string, string> = { trial: 'ทดลองขาย', temporary: 'ทำการค้าแบบชั่วคราว', contract: 'ทำการค้าแบบมีสัญญา' }
let loadSequence = 0
const stallOptions = computed(() => {
  const b = props.booking
  if (!b) return []
  const items = Array.isArray(b.sub_bookings) && b.sub_bookings.length ? b.sub_bookings : [b]
  return items.filter((item: any) => !['cancelled', 'rejected', 'not_interested', 'no_show', 'terminated', 'interested', 'lease_applied', 'completed'].includes(String(item.status || '')))
    .flatMap((item: any) => String(item.stall_ident || item.stall_number || '').split(',').map((code: string) => code.trim()).filter(Boolean))
    .map((code: string) => ({ value: code, label: code }))
})
const loadTypes = async () => {
  const sequence = ++loadSequence
  tradeTypes.value = []
  typesError.value = ''
  rentType.value = ''
  const selected = props.booking?.sub_bookings?.find((item: any) => String(item.stall_ident || item.stall_number || '').trim() === stall.value)
  const marketId = selected?.market_id || props.booking?.market_id
  if (!marketId) { typesError.value = 'ไม่พบตลาดของการจอง กรุณาตรวจสอบข้อมูลการจอง'; typesLoading.value = false; return }
  typesLoading.value = true
  try {
    const response: any = await get('/markets/' + encodeURIComponent(String(marketId)) + '/trade-conditions')
    if (sequence === loadSequence) tradeTypes.value = Array.isArray(response.data) ? response.data : []
  } catch {
    if (sequence === loadSequence) typesError.value = 'ไม่สามารถโหลดประเภทการค้าได้'
  } finally {
    if (sequence === loadSequence) typesLoading.value = false
  }
}
const confirmRental = () => {
  submitted.value = true
  if (!stall.value || !rentType.value || props.loading || typesLoading.value || !tradeTypes.value.some(item => item.TradeType === rentType.value)) return
  emit('confirm', { stallIdent: stall.value, rentType: rentType.value })
}
watch(() => props.open, (open) => {
  if (!open) { loadSequence++; return }
  submitted.value = false
  const code = String(props.booking?.stall_ident || props.booking?.stall_number || '').trim()
  stall.value = stallOptions.value.length === 1 ? stallOptions.value[0].value : (stallOptions.value.some((option: any) => option.value === code) ? code : '')
  loadTypes()
})
</script>
<style scoped>
.rent-overlay { position: fixed; inset: 0; z-index: 2000; display: grid; place-items: center; padding: 16px; background: rgba(15,23,42,.55); }
.rent-dialog { width: min(100%, 520px); max-height:min(90vh,90dvh); display:flex; flex-direction:column; overflow:hidden; border-radius: 20px; background: var(--bg-ffffff); box-shadow: 0 20px 35px rgba(15,23,42,.2); }
.rent-body{flex:1 1 auto;min-height:0;overflow-y:auto;overscroll-behavior:contain;padding:24px 28px}
.rent-icon { width: 54px; height: 54px; display: grid; place-items: center; margin: 0 auto 16px; color: var(--tx-2563eb); background: var(--bg-eff6ff); border-radius: 14px; }
h3 { margin: 0; text-align: center; font-size: var(--font-size-xl); color: var(--tx-0f172a); } p { color: var(--tx-64748b); text-align: center; font-size: var(--font-size-base); line-height: 1.5; }
.rent-message { margin:0; text-align:left; font-size:var(--font-size-sm); }
.rent-error { margin: 0; color: var(--tx-dc2626); text-align: left; }
.rent-fields { display: grid; gap: 10px; } label { font-size: var(--font-size-sm); font-weight: 600; color: var(--tx-334155); }
.rent-choices { display: flex; flex-wrap: wrap; gap: 12px 20px; } .rent-actions { flex:0 0 auto;display: flex; justify-content: flex-end; gap: 10px; padding:16px 28px; border-top:1px solid var(--bd-e2e8f0); background:var(--bg-f8fafc); }
</style>
