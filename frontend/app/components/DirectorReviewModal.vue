<template>
  <AppModal :model-value="open" :title="kind === 'market' ? 'ตรวจสอบหนี้ค้าง' : 'ผลการพิจารณาคำร้อง'" :description="kind === 'market' ? 'ระบุผลตรวจของงานการคลังก่อนลงนาม' : 'ระบุผลพิจารณาก่อนลงนาม'" @update:model-value="onOpenChange">
    <div class="review-form">
      <fieldset v-if="kind === 'market'" class="review-group">
        <legend>2.2 งานการคลัง/ฝ่ายการคลัง <span class="req-star" aria-hidden="true">*</span></legend>
        <div class="review-choices">
          <AppRadio v-model="market.debt_status" name="market-debt-status" value="none" label="ไม่มีหนี้ค้างชำระ" />
          <AppRadio v-model="market.debt_status" name="market-debt-status" value="has" label="มีหนี้ค้าง" />
        </div>
        <label v-if="market.debt_status === 'has'" class="review-extra">จำนวนหนี้ค้าง (บาท) <span class="req-star" aria-hidden="true">*</span>
          <AppInput v-model="market.debt_amount" type="number" min="0.01" step="0.01" required />
        </label>
      </fieldset>
      <template v-else>
        <fieldset class="review-group">
          <legend>3. ผู้มีอำนาจอนุมัติ <span class="req-star" aria-hidden="true">*</span></legend>
          <div class="review-choices">
            <AppRadio v-model="bureau.decision" name="bureau-decision" value="approve" label="อนุมัติ" />
            <AppRadio v-model="bureau.decision" name="bureau-decision" value="reject" label="ไม่อนุมัติ" />
          </div>

        </fieldset>
        <div class="review-group">
          <AppCheckbox v-model="bureau.office_checked" label="สำนักงาน เลขระบบ และวันที่" />
          <div v-if="bureau.office_checked" class="review-fields">
            <label>เลขระบบ<AppInput v-model="bureau.system_number" maxlength="255" placeholder="ระบุเลขระบบ" /></label>
            <label>วันที่<AppInput v-model="bureau.document_date" type="date" /></label>
          </div>
        </div>
      </template>
      <p v-if="error" class="review-error" role="alert">{{ error }}</p>
    </div>
    <template #footer>
      <AppButton size="md" variant="secondary" @click="$emit('cancel')">ยกเลิก</AppButton>
      <AppButton size="md" variant="primary" @click="confirm">ลงนามต่อ</AppButton>
    </template>
  </AppModal>
</template>
<script setup lang="ts">
type MarketReview = { debt_status: string; debt_amount: string | number }
type BureauReview = { decision: string; other_details: string; office_checked: boolean; system_number: string; document_date: string }
const props = defineProps<{ open: boolean; kind: 'market' | 'bureau'; application?: any; initialDecision?: 'approve' | 'reject' }>()
const emit = defineEmits<{ cancel: []; confirm: [review: MarketReview | BureauReview] }>()
const market = ref<MarketReview>({ debt_status: '', debt_amount: '' })
const bureau = ref<BureauReview>({ decision: '', other_details: '', office_checked: false, system_number: '', document_date: '' })
const error = ref('')
watch(() => props.open, (open) => {
  if (!open) return
  market.value = { debt_status: '', debt_amount: '', ...(props.application?.market_director_review || {}) }
  bureau.value = { decision: '', other_details: '', office_checked: false, system_number: '', document_date: '', ...(props.application?.bureau_director_review || {}), ...(props.initialDecision ? { decision: props.initialDecision } : {}) }
  error.value = ''
})
const onOpenChange = (value: boolean) => { if (!value) emit('cancel') }
const confirm = () => {
  error.value = ''
  if (props.kind === 'market') {
    if (!['none', 'has'].includes(market.value.debt_status)) { error.value = 'กรุณาเลือกข้อมูลหนี้ค้าง'; return }
    const amount = Number(market.value.debt_amount)
    if (market.value.debt_status === 'has' && (!Number.isFinite(amount) || amount <= 0 || amount > 1000000000000)) { error.value = 'กรุณาระบุจำนวนหนี้ค้าง'; return }
    emit('confirm', { ...market.value, debt_amount: market.value.debt_status === 'has' ? amount : 0 })
    return
  }
  if (!bureau.value.decision) { error.value = 'กรุณาเลือกผลการพิจารณา'; return }
  if ([...bureau.value.other_details.trim()].length > 255 || [...bureau.value.system_number.trim()].length > 255) { error.value = 'รายละเอียดแต่ละช่องต้องไม่เกิน 255 ตัวอักษร'; return }
  emit('confirm', { ...bureau.value })
}
</script>
<style scoped>
.review-form{display:grid;gap:16px;color:var(--tx-334155)}.review-group{min-width:0;margin:0;padding:0;border:0}.review-group legend,.review-extra,.review-fields label{display:block;margin-bottom:8px;color:var(--tx-334155);font-size:var(--font-size-sm);font-weight:600}.review-choices{display:flex;flex-wrap:wrap;gap:12px 24px}.review-extra{margin-top:12px}.review-extra :deep(.app-input),.review-fields :deep(.app-input){margin-top:8px}.review-fields{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px;margin-top:12px}.review-error,.req-star{color:var(--tx-dc2626)}.review-error{margin:0;font-size:var(--font-size-sm)}@media(max-width:480px){.review-fields{grid-template-columns:1fr}}
</style>
