<template>
  <AppModal :model-value="open" title="ตรวจสอบข้อมูลคำร้อง" description="ระบุผลตรวจสอบก่อนลงนาม" @update:model-value="onOpenChange">
    <div class="review-form">
      <fieldset class="review-group">
        <legend>ข้อมูลแผงค้า <span class="req-star" aria-hidden="true">*</span></legend>
        <div class="review-choices">
          <AppRadio v-model="form.stall_status" name="officer-stall-status" value="none" label="ไม่มีแผงค้า" />
          <AppRadio v-model="form.stall_status" name="officer-stall-status" value="has" label="มีแผงค้า" />
          <AppRadio v-model="form.stall_status" name="officer-stall-status" value="other" label="อื่น ๆ" />
        </div>
        <div v-if="form.stall_status === 'has'" class="review-fields">
          <div><label for="officer-stall-count">จำนวนแผง <span class="req-star" aria-hidden="true">*</span></label><AppInput id="officer-stall-count" v-model="form.stall_count" type="number" min="1" max="100000" step="1" required /></div>
          <div><label for="officer-stall-numbers">เลขที่แผง</label><AppInput id="officer-stall-numbers" v-model="form.stall_numbers" maxlength="255" placeholder="ระบุเลขที่แผง" /></div>
          <div><label for="officer-stall-project">โครงการ</label><AppInput id="officer-stall-project" v-model="form.project" maxlength="255" placeholder="ระบุโครงการ" /></div>
          <div><label for="officer-stall-type">ประเภท</label><AppInput id="officer-stall-type" v-model="form.trade_type" maxlength="255" placeholder="ระบุประเภท" /></div>
        </div>
      </fieldset>
      <AppCheckbox v-model="form.other_contract" label="มีสัญญาอื่น ๆ (เช่า) กับสำนักงานตลาดกรุงเทพมหานคร" />
      <div v-if="form.stall_status === 'other'" class="review-group">
        <AppCheckbox v-model="form.other_checked" label="อื่น ๆ" />
        <div v-if="form.other_checked" class="review-extra"><label for="officer-other">รายละเอียด <span class="req-star" aria-hidden="true">*</span></label><AppInput id="officer-other" v-model="form.other_details" maxlength="255" placeholder="ระบุรายละเอียด" required /></div>
      </div>
      <div class="review-group">
        <AppCheckbox v-model="form.documents_complete" label="คุณสมบัติและเอกสารถูกต้องครบถ้วน" />
        <div v-if="form.documents_complete" class="review-extra"><label for="officer-reference">เลขที่รับเรื่องหรือใบจองจากระบบ</label><AppInput id="officer-reference" v-model="form.booking_reference" maxlength="255" placeholder="ระบุเลขที่" /></div>
      </div>
      <p v-if="error" class="review-error" role="alert">{{ error }}</p>
    </div>
    <template #footer>
      <AppButton size="md" variant="secondary" @click="$emit('cancel')">ยกเลิก</AppButton>
      <AppButton size="md" variant="primary" @click="confirm">ลงนามต่อ</AppButton>
    </template>
  </AppModal>
</template>
<script setup lang="ts">
type OfficerReview = { stall_status: string; stall_count: string | number; stall_numbers: string; project: string; trade_type: string; other_contract: boolean; other_checked: boolean; other_details: string; documents_complete: boolean; booking_reference: string }
const props = defineProps<{ open: boolean; application?: any }>()
const emit = defineEmits<{ cancel: []; confirm: [review: OfficerReview] }>()
const emptyForm = (): OfficerReview => ({ stall_status: '', stall_count: '', stall_numbers: '', project: '', trade_type: '', other_contract: false, other_checked: false, other_details: '', documents_complete: false, booking_reference: '' })
const form = ref<OfficerReview>(emptyForm())
const error = ref('')
watch(() => props.open, (open) => {
  if (!open) return
  form.value = { ...emptyForm(), ...(props.application?.officer_review || {}) }
  error.value = ''
})
const onOpenChange = (value: boolean) => { if (!value) emit('cancel') }
const confirm = () => {
  if (form.value.stall_status === 'other') {
    form.value.other_checked = true
    form.value.stall_status = 'none'
  }
  error.value = ''
  if (!['none', 'has'].includes(form.value.stall_status)) { error.value = 'กรุณาเลือกข้อมูลแผงค้า'; return }
  const count = Number(form.value.stall_count)
  if (form.value.stall_status === 'has' && (!Number.isInteger(count) || count < 1 || count > 100000)) { error.value = 'กรุณาระบุจำนวนแผงตั้งแต่ 1 ถึง 100,000'; return }
  if (form.value.other_checked && !form.value.other_details.trim()) { error.value = 'กรุณาระบุรายละเอียดอื่น ๆ'; return }
  const text = [form.value.stall_numbers, form.value.project, form.value.trade_type, form.value.other_details, form.value.booking_reference]
  if (text.some(value => [...value.trim()].length > 255)) { error.value = 'รายละเอียดแต่ละช่องต้องไม่เกิน 255 ตัวอักษร'; return }
  emit('confirm', { ...form.value, stall_count: form.value.stall_status === 'has' ? count : 0 })
}
</script>
<style scoped>
.review-form{display:grid;gap:16px;color:var(--tx-334155)}.review-group{min-width:0;margin:0;padding:0;border:0}.review-group legend,.review-fields label,.review-extra label{display:block;margin-bottom:8px;color:var(--tx-334155);font-size:var(--font-size-sm);font-weight:600}.review-choices{display:flex;flex-wrap:wrap;gap:12px 24px}.review-fields{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px;margin-top:12px}.review-extra{margin-top:10px}.review-error,.req-star{color:var(--tx-dc2626)}.review-error{margin:0;font-size:var(--font-size-sm)}@media(max-width:480px){.review-fields{grid-template-columns:1fr}}
</style>
