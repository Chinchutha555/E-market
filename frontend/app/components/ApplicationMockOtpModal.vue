<template>
  <AppModal :model-value="open" title="ยืนยันการอนุมัติ" description="กรอกรหัส OTP เพื่อยืนยันการลงนาม" @update:model-value="close">
    <p class="mock-notice">ระบบทดสอบ: ใช้รหัส <strong>123456</strong> ยังไม่มีการส่งอีเมล</p>
    <div class="form-label" id="otp-label">รหัส OTP <span class="req-star" aria-hidden="true">*</span></div>
    <div class="otp-inputs" role="group" aria-labelledby="otp-label" @paste.prevent="pasteCode">
      <input
        v-for="(_, index) in digits"
        :id="'application-mock-otp-' + index"
        :key="index"
        :value="digits[index]"
        :aria-label="'รหัส OTP หลักที่ ' + (index + 1)"
        :aria-invalid="!!error"
        type="text"
        inputmode="numeric"
        pattern="[0-9]*"
        maxlength="1"
        :autocomplete="index === 0 ? 'one-time-code' : 'off'"
        @input="inputDigit(index, $event)"
        @keydown="keyDigit(index, $event)"
      />
    </div>
    <p v-if="error" class="otp-error" role="alert">{{ error }}</p>
    <template #footer>
      <AppButton variant="secondary" size="md" @click="close">ยกเลิก</AppButton>
      <AppButton variant="primary" size="md" @click="confirm">ยืนยันการอนุมัติ</AppButton>
    </template>
  </AppModal>
</template>

<script setup lang="ts">
const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ cancel: []; confirm: [code: string] }>()
const digits = ref<string[]>(Array(6).fill(''))
const error = ref('')
const focusDigit = (index: number) => {
  if (index >= 0 && index < 6) document.getElementById('application-mock-otp-' + index)?.focus()
}
watch(() => props.open, async (open) => {
  digits.value = Array(6).fill('')
  error.value = ''
  if (open) { await nextTick(); focusDigit(0) }
})
const close = () => emit('cancel')
const inputDigit = (index: number, event: Event) => {
  const input = event.target as HTMLInputElement
  const value = input.value.replace(/[^0-9]/g, '').slice(-1)
  digits.value[index] = value
  input.value = value
  error.value = ''
  if (value) focusDigit(index + 1)
}
const keyDigit = (index: number, event: KeyboardEvent) => {
  if (event.key === 'Backspace' && !digits.value[index] && index > 0) {
    event.preventDefault()
    digits.value[index - 1] = ''
    focusDigit(index - 1)
  } else if (event.key === 'ArrowLeft') {
    event.preventDefault()
    focusDigit(index - 1)
  } else if (event.key === 'ArrowRight') {
    event.preventDefault()
    focusDigit(index + 1)
  } else if (event.key === 'Enter') {
    event.preventDefault()
    confirm()
  }
}
const pasteCode = (event: ClipboardEvent) => {
  const value = (event.clipboardData?.getData('text') || '').replace(/[^0-9]/g, '').slice(0, 6)
  if (!value) return
  digits.value = Array.from({ length: 6 }, (_, index) => value[index] || '')
  error.value = ''
  nextTick(() => focusDigit(Math.min(value.length, 5)))
}
const confirm = () => {
  const code = digits.value.join('')
  if (code !== '123456') { error.value = 'รหัส OTP ทดสอบไม่ถูกต้อง'; return }
  emit('confirm', code)
}
</script>

<style scoped>
.mock-notice { margin: 0 0 20px; padding: 12px 14px; color: #1e3a8a; background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 8px; }
.otp-inputs { display: flex; gap: 8px; }
.otp-inputs input { width: 48px; min-width: 0; height: 52px; border: 1px solid #cbd5e1; border-radius: 8px; text-align: center; font: inherit; font-size: var(--font-size-xl); font-weight: 600; color: #1e293b; }
.otp-inputs input:focus { outline: 2px solid #93c5fd; border-color: #2563eb; }
.otp-inputs input[aria-invalid="true"] { border-color: #ef4444; }
.otp-error { margin: 8px 0 0; color: #dc2626; font-size: var(--font-size-xs); }
@media (max-width: 420px) { .otp-inputs { gap: 5px; } .otp-inputs input { flex: 1; width: 0; } }
</style>