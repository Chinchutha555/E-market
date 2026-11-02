<template>
  <AppModal :model-value="!!application && !showPreview" title="จัดการเอกสาร" :description="'คำร้อง ' + (application?.application_number || application?.application_no || '-')" @update:model-value="close">
    <div class="document-step">
      <div class="step-number">1</div>
      <div class="step-content">
        <h3>ดาวน์โหลดคำร้อง</h3>
        <p>ดาวน์โหลดคำร้องเป็น PDF แล้วนำให้ผู้เช่าลงนาม</p>
        <AppButton variant="outline" size="md" @click="showPreview = true">
          <template #prefix><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg></template>
          เปิดเอกสาร
        </AppButton>
      </div>
    </div>
    <div class="document-step">
      <div class="step-number">2</div>
      <div class="step-content">
        <h3>นำเข้าเอกสารลงนาม</h3>
        <p>แนบคำร้องที่ผู้เช่าลงนามแล้ว เพื่อปิดคำร้อง</p>
        <AppFileUpload v-model="signedFile" label="เอกสารลงนาม" required-text="จำเป็น" accept=".pdf,.jpg,.jpeg,.png,application/pdf,image/jpeg,image/png" :max-size-mb="10" hint="PDF, JPG หรือ PNG ขนาดไม่เกิน 10MB" :error="fileError" @error="fileError = $event" />
      </div>
    </div>
    <p v-if="submitError" class="submit-error" role="alert">{{ submitError }}</p>
    <template #footer>
      <AppButton variant="secondary" size="md" :disabled="submitting" @click="close">ปิด</AppButton>
      <AppButton variant="primary" size="md" :loading="submitting" :disabled="submitting" @click="submit">บันทึกเอกสาร</AppButton>
    </template>
  </AppModal>
  <ApplicationDocumentPreview :application="showPreview ? application : null" @close="showPreview = false" />
</template>

<script setup lang="ts">
import AppFileUpload from '~/components/AppFileUpload.vue'
import ApplicationDocumentPreview from '~/components/ApplicationDocumentPreview.vue'

const props = defineProps<{ application: any | null }>()
const emit = defineEmits<{ close: []; completed: [] }>()
const { post } = useApi()
const { showSuccess } = useAppAlert()
const showPreview = ref(false)
const signedFile = ref<File | null>(null)
const fileError = ref('')
const submitError = ref('')
const submitting = ref(false)

watch(() => props.application, () => {
  showPreview.value = false
  signedFile.value = null
  fileError.value = ''
  submitError.value = ''
})
const close = () => { if (!submitting.value) emit('close') }
const submit = async () => {
  if (!signedFile.value) { fileError.value = 'กรุณาแนบเอกสารลงนาม'; return }
  if (fileError.value || submitting.value || !props.application?.id) return
  submitting.value = true
  submitError.value = ''
  try {
    const file = signedFile.value
    const encoded = await new Promise<string>((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(String(reader.result).split(',')[1] || '')
      reader.onerror = () => reject(new Error('อ่านไฟล์ไม่สำเร็จ'))
      reader.readAsDataURL(file)
    })
    await post('/applications/' + props.application.id + '/signed-document', { file: encoded })
    emit('completed')
    await showSuccess('บันทึกเอกสารและปิดคำร้องแล้ว', 'บันทึกสำเร็จ')
  } catch (error: any) {
    submitError.value = error?.data?.detail || error?.data?.message || error?.message || 'บันทึกเอกสารไม่สำเร็จ'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.document-step { display: flex; align-items: flex-start; gap: 14px; padding: 4px 0 20px; }
.document-step + .document-step { padding-top: 20px; border-top: 1px solid #e2e8f0; }
.step-number { display: grid; place-items: center; width: 32px; height: 32px; flex: 0 0 32px; border-radius: 50%; background: #eff6ff; color: #2563eb; font-size: var(--font-size-sm); font-weight: 700; }
.step-content { min-width: 0; flex: 1; }
.step-content h3 { margin: 2px 0 4px; color: #1e293b; font-size: var(--font-size-md); font-weight: 700; }
.step-content p { margin: 0 0 14px; color: #64748b; font-size: var(--font-size-base); line-height: 1.5; }
.submit-error { margin: 0; color: #dc2626; font-size: var(--font-size-sm); }
</style>