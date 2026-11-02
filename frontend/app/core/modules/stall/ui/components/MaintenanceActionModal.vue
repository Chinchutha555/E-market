<template>
  <AppModal
    :model-value="modelValue"
    icon="tool"
    :title="modalTitle"
    :description="itemTitle ? `รายการ: ${itemTitle}` : modalDescription"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <form class="maintenance-action-form" novalidate @submit.prevent="submit">
      <template v-if="mode === 'accept'">
        <label>
          <span>ผู้รับผิดชอบงาน <b class="req-star" aria-hidden="true">*</b></span>
          <AppSelect v-model="acceptForm.assigneeType" :options="assigneeOptions" :error="fieldErrors.assigneeType" />
        </label>
        <label>
          <span>กำหนดเสร็จ <b class="req-star" aria-hidden="true">*</b></span>
          <AppInput v-model="acceptForm.scheduledDate" type="date" required :error="fieldErrors.scheduledDate" />
        </label>
        <label>
          <span>ค่าใช้จ่ายรับผิดชอบโดย</span>
          <AppSelect v-model="acceptForm.expensePayer" :options="payerOptions" />
        </label>
        <label>
          <span>ความเร่งด่วน</span>
          <AppSelect v-model="acceptForm.priority" :options="priorityOptions" />
        </label>
        <label class="is-wide">
          <span>ข้อความถึงผู้เช่า (ถ้ามี)</span>
          <textarea v-model="acceptForm.tenantMessage" rows="3" placeholder="เช่น ช่างจะเข้าตรวจวันที่ ..."></textarea>
        </label>
      </template>

      <template v-else-if="mode === 'complete'">
        <label>
          <span>ผู้ดำเนินการแก้ไข <b class="req-star" aria-hidden="true">*</b></span>
          <AppSelect v-model="completeForm.executorType" :options="assigneeOptions" placeholder="เลือกประเภทผู้ดำเนินการ" :error="fieldErrors.executorType" />
        </label>
        <label>
          <span>ชื่อผู้แก้ไข / บริษัท <b class="req-star" aria-hidden="true">*</b></span>
          <AppInput v-model="completeForm.executorName" placeholder="ชื่อช่างหรือบริษัท" required :error="fieldErrors.executorName" />
        </label>
        <label>
          <span>วันที่แก้เสร็จ</span>
          <AppInput v-model="completeForm.completionDate" type="date" />
        </label>
        <label>
          <span>ค่าใช้จ่ายจริง (ถ้ามี)</span>
          <AppInput v-model="completeForm.actualCost" type="number" placeholder="0" />
        </label>
        <label class="is-wide">
          <span>สิ่งที่แก้ไข</span>
          <textarea v-model="completeForm.repairSummary" rows="3" placeholder="ระบุสิ่งที่ดำเนินการแก้ไข"></textarea>
        </label>

        <AppImageUpload
          v-model="completeForm.afterPhotos"
          class="is-wide"
          label="รูปหลังแก้"
          required-text="อย่างน้อย 1 รูป"
          multiple
          select-text="เพิ่มรูป"
          :max-files="5"
          :error="fieldErrors.afterPhotos"
          preview-alt="รูปหลังแก้"
        />
      </template>

      <template v-else>
        <label class="is-wide">
          <span>เหตุผล <b class="req-star" aria-hidden="true">*</b></span>
          <AppSelect v-model="rejectForm.reason" :options="rejectReasonOptions" placeholder="เลือกเหตุผล" :error="fieldErrors.reason" />
        </label>
        <label class="is-wide">
          <span>รายละเอียดเพิ่มเติม <b class="req-star" aria-hidden="true">*</b></span>
          <textarea
            v-model="rejectForm.detail"
            rows="4"
            required
            placeholder="ข้อความนี้จะแสดงให้ผู้เช่าเห็น"
            :class="{ 'has-error': fieldErrors.detail }"
            :aria-invalid="Boolean(fieldErrors.detail)"
          ></textarea>
          <span v-if="fieldErrors.detail" class="field-error">{{ fieldErrors.detail }}</span>
        </label>
      </template>
    </form>

    <template #footer>
      <AppButton size="md" variant="secondary" @click="emit('update:modelValue', false)">ยกเลิก</AppButton>
      <AppButton size="md" :variant="confirmVariant" @click="submit">
        {{ confirmText }}
      </AppButton>
    </template>
  </AppModal>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import AppButton from '~/components/AppButton.vue'
import AppImageUpload, { type AppImageFile } from '~/components/AppImageUpload.vue'
import AppInput from '~/components/AppInput.vue'
import AppModal from '~/components/AppModal.vue'
import AppSelect from '~/components/AppSelect.vue'

export type MaintenanceActionMode = 'accept' | 'complete' | 'reject'
export type MaintenancePhoto = AppImageFile
export interface MaintenanceActionItem {
  title: string
  assigneeType?: 'market_technician' | 'contractor' | 'tenant'
  assigneeName?: string
  scheduledDate?: string
  expensePayer?: 'market' | 'tenant' | 'shared'
  priority?: 'normal' | 'urgent' | 'critical'
  tenantMessage?: string
  executorType?: 'market_technician' | 'contractor' | 'tenant'
  executorName?: string
  completionDate?: string
  actualCost?: number
  repairSummary?: string
  afterPhotos?: MaintenancePhoto[]
  rejectionReason?: string
  rejectionDetail?: string
}
export type MaintenanceActionPayload =
  | { mode: 'accept'; assigneeType: 'market_technician' | 'contractor' | 'tenant'; scheduledDate: string; expensePayer: 'market' | 'tenant' | 'shared'; priority: 'normal' | 'urgent' | 'critical'; tenantMessage: string }
  | { mode: 'complete'; executorType: 'market_technician' | 'contractor' | 'tenant'; executorName: string; completionDate: string; actualCost: number | string; repairSummary: string; afterPhotos: MaintenancePhoto[] }
  | { mode: 'reject'; reason: string; detail: string }

const props = defineProps<{
  modelValue: boolean
  mode: MaintenanceActionMode
  item?: MaintenanceActionItem | null
  itemTitle?: string
}>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  submit: [payload: MaintenanceActionPayload]
}>()

const assigneeOptions = [
  { value: 'market_technician', label: 'ช่างสำนักงานตลาด' },
  { value: 'contractor', label: 'ผู้รับเหมา' },
  { value: 'tenant', label: 'ให้ผู้เช่าดำเนินการเอง' }
]
const payerOptions = [
  { value: 'market', label: 'สำนักงานตลาด' },
  { value: 'tenant', label: 'ผู้เช่า' },
  { value: 'shared', label: 'ร่วมรับผิดชอบค่าใช้จ่าย' }
]
const priorityOptions = [
  { value: 'normal', label: 'ปกติ' },
  { value: 'urgent', label: 'เร่งด่วน' },
  { value: 'critical', label: 'ฉุกเฉิน' }
]
const rejectReasonOptions = [
  { value: 'tenant_responsibility', label: 'เป็นความรับผิดชอบของผู้เช่าตามสัญญา' },
  { value: 'insufficient_information', label: 'ข้อมูลหรือหลักฐานไม่เพียงพอ' },
  { value: 'duplicate', label: 'เป็นรายการซ้ำ' },
  { value: 'outside_scope', label: 'อยู่นอกขอบเขตบริการของตลาด' },
  { value: 'other', label: 'เหตุผลอื่น' }
]

const acceptForm = reactive({
  assigneeType: 'market_technician' as 'market_technician' | 'contractor' | 'tenant',
  scheduledDate: '',
  expensePayer: 'market' as 'market' | 'tenant' | 'shared',
  priority: 'normal' as 'normal' | 'urgent' | 'critical',
  tenantMessage: ''
})
const completeForm = reactive({
  executorType: '' as '' | 'market_technician' | 'contractor' | 'tenant',
  executorName: '',
  completionDate: '',
  actualCost: '' as number | string,
  repairSummary: '',
  afterPhotos: [] as MaintenancePhoto[]
})
const rejectForm = reactive({ reason: '', detail: '' })
const validationAttempted = ref(false)

const resetForm = () => {
  validationAttempted.value = false
  const item = props.item
  Object.assign(acceptForm, {
    assigneeType: item?.assigneeType || 'market_technician',
    scheduledDate: item?.scheduledDate || '',
    expensePayer: item?.expensePayer || 'market',
    priority: item?.priority || 'normal',
    tenantMessage: item?.tenantMessage || ''
  })
  Object.assign(completeForm, {
    executorType: item?.executorType || item?.assigneeType || '',
    executorName: item?.executorName || item?.assigneeName || '',
    completionDate: item?.completionDate || new Date().toISOString().slice(0, 10),
    actualCost: item?.actualCost ?? '',
    repairSummary: item?.repairSummary || '',
    afterPhotos: [...(item?.afterPhotos || [])]
  })
  Object.assign(rejectForm, { reason: item?.rejectionReason || '', detail: item?.rejectionDetail || '' })
}
watch(() => [props.modelValue, props.mode, props.item], () => { if (props.modelValue) resetForm() }, { immediate: true })

const modalTitle = computed(() => ({ accept: 'รับเรื่องแจ้งซ่อม', complete: 'บันทึกผลการซ่อม', reject: 'ไม่รับเรื่องแจ้งซ่อม' })[props.mode])
const modalDescription = computed(() => ({ accept: 'มอบหมายผู้รับผิดชอบและกำหนดวันแล้วเสร็จ', complete: 'บันทึกผลและหลักฐานก่อนปิดงานซ่อม', reject: 'ระบุเหตุผลเพื่อแจ้งกลับผู้เช่า' })[props.mode])
const confirmText = computed(() => ({ accept: 'ยืนยันรับเรื่อง', complete: 'ยืนยันเสร็จสิ้น', reject: 'ยืนยันไม่รับเรื่อง' })[props.mode])
const confirmVariant = computed<'primary' | 'success' | 'danger'>(() => ({ accept: 'primary', complete: 'success', reject: 'danger' })[props.mode] as 'primary' | 'success' | 'danger')
const fieldErrors = computed<Record<string, string>>(() => {
  if (!validationAttempted.value) return {}
  if (props.mode === 'accept') {
    return {
      assigneeType: acceptForm.assigneeType ? '' : 'กรุณาเลือกผู้รับผิดชอบงาน',
      scheduledDate: acceptForm.scheduledDate ? '' : 'กรุณาระบุกำหนดเสร็จ'
    }
  }
  if (props.mode === 'complete') {
    return {
      executorType: completeForm.executorType ? '' : 'กรุณาเลือกประเภทผู้ดำเนินการ',
      executorName: completeForm.executorName.trim() ? '' : 'กรุณาระบุชื่อผู้แก้ไขหรือบริษัท',
      afterPhotos: completeForm.afterPhotos.length ? '' : 'กรุณาแนบรูปหลังแก้อย่างน้อย 1 รูป'
    }
  }
  return {
    reason: rejectForm.reason ? '' : 'กรุณาเลือกเหตุผล',
    detail: rejectForm.detail.trim() ? '' : 'กรุณาระบุรายละเอียดเพิ่มเติม'
  }
})
const hasValidationErrors = computed(() => Object.values(fieldErrors.value).some(Boolean))

const submit = () => {
  validationAttempted.value = true
  if (hasValidationErrors.value) return
  if (props.mode === 'accept') {
    emit('submit', { mode: 'accept', ...acceptForm })
  } else if (props.mode === 'complete' && completeForm.executorType) {
    emit('submit', { mode: 'complete', ...completeForm, executorType: completeForm.executorType, afterPhotos: [...completeForm.afterPhotos] })
  } else if (props.mode === 'reject') {
    emit('submit', { mode: 'reject', reason: rejectForm.reason, detail: rejectForm.detail })
  }
}
</script>

<style scoped>
.maintenance-action-form{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:16px}.maintenance-action-form label{display:flex;flex-direction:column;gap:7px;min-width:0;color:var(--tx-334155);font-size:var(--font-size-sm);font-weight:600}.maintenance-action-form .is-wide{grid-column:1/-1}.maintenance-action-form textarea{box-sizing:border-box;width:100%;min-height:96px;padding:11px 14px;resize:vertical;border:1px solid var(--bd-cbd5e1);border-radius:10px;background:var(--bg-ffffff);color:var(--tx-1e293b);font:500 .9375rem/1.5 inherit;outline:none}.maintenance-action-form textarea:focus{border-color:#2563eb;box-shadow:0 0 0 3px rgba(37,99,235,.12)}.maintenance-action-form textarea.has-error{border-color:#ef4444;box-shadow:0 0 0 1px #ef4444;background:var(--bg-fffafa)}.req-star{color:var(--tx-ef4444)}.field-error{display:block;color:var(--tx-dc2626);font-size:var(--font-size-xs);font-weight:500;line-height:1.3}@media(max-width:640px){.maintenance-action-form{grid-template-columns:1fr}.maintenance-action-form .is-wide{grid-column:auto}}
</style>
