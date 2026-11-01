<template>
  <div class="app-file-upload">
    <span v-if="label" class="app-file-upload-label">
      {{ label }}
      <b v-if="requiredText" class="app-file-upload-required">* {{ requiredText }}</b>
    </span>
    <div
      class="app-file-upload-panel"
      :class="{ 'has-error': resolvedError, 'is-dragging': isDragging }"
      role="button"
      tabindex="0"
      :aria-label="modelValue ? 'เปลี่ยนไฟล์' : selectText"
      @click="openFilePicker"
      @keydown.enter.prevent="openFilePicker"
      @keydown.space.prevent="openFilePicker"
      @dragover.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @drop.prevent="handleDrop"
    >
      <div class="app-file-upload-icon" :class="`is-${fileKind}`" aria-hidden="true">
        <svg v-if="fileKind === 'image'" viewBox="0 0 24 24"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><path d="m21 15-5-5L5 21"/></svg>
        <svg v-else-if="fileKind === 'sheet'" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8M11 9v12"/></svg>
        <svg v-else-if="fileKind === 'archive'" viewBox="0 0 24 24"><path d="M21 8v13H3V8M1 3h22v5H1zM10 12h4"/></svg>
        <svg v-else-if="fileKind === 'file-text'" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg>
        <svg v-else viewBox="0 0 24 24"><path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/></svg>
      </div>
      <div class="app-file-upload-copy">
        <strong>{{ modelValue?.name || selectText }}</strong>
        <span v-if="modelValue">{{ formatFileSize(modelValue.size) }} · คลิกเพื่อเปลี่ยนไฟล์</span>
        <span v-else>{{ description || 'คลิกเพื่อเลือกไฟล์ หรือลากไฟล์มาวางที่นี่' }}</span>
        <small>{{ resolvedHint }}</small>
      </div>
      <AppButton v-if="modelValue" size="sm" variant="danger-outline" @click.stop="removeFile">ลบไฟล์</AppButton>
      <input ref="fileInput" type="file" :accept="accept" class="app-file-upload-input" @change="handleInput" />
    </div>
    <span v-if="resolvedError" class="app-file-upload-error">{{ resolvedError }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import AppButton from './AppButton.vue'

const props = withDefaults(defineProps<{
  modelValue?: File | null
  label?: string
  requiredText?: string
  accept?: string
  maxSizeMb?: number
  hint?: string
  description?: string
  selectText?: string
  error?: string
}>(), {
  modelValue: null,
  label: '',
  requiredText: '',
  accept: '',
  maxSizeMb: 10,
  hint: '',
  description: '',
  selectText: 'เลือกไฟล์',
  error: ''
})

const emit = defineEmits<{
  'update:modelValue': [file: File | null]
  change: [file: File | null]
  error: [message: string]
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const internalError = ref('')
const isDragging = ref(false)
const resolvedError = computed(() => props.error || internalError.value)
const resolvedHint = computed(() => props.hint || `ขนาดไฟล์ไม่เกิน ${props.maxSizeMb}MB`)
const fileKind = computed<'image' | 'file-text' | 'sheet' | 'archive' | 'paperclip'>(() => {
  const type = props.modelValue?.type || ''
  const name = props.modelValue?.name.toLowerCase() || ''
  if (type.startsWith('image/')) return 'image'
  if (type === 'application/pdf' || /\.(pdf|doc|docx|txt)$/.test(name) || /pdf|document|text/.test(props.accept)) return 'file-text'
  if (/spreadsheet|csv|excel/.test(type) || /\.(csv|xls|xlsx)$/.test(name)) return 'sheet'
  if (/zip|compressed|archive/.test(type) || /\.(zip|rar|7z)$/.test(name)) return 'archive'
  return 'paperclip'
})

const openFilePicker = () => fileInput.value?.click()
const formatFileSize = (bytes: number) => bytes >= 1024 * 1024 ? `${(bytes / 1024 / 1024).toFixed(2)} MB` : `${Math.max(1, Math.round(bytes / 1024))} KB`

const matchesAccept = (file: File) => {
  if (!props.accept) return true
  return props.accept.split(',').map(value => value.trim()).some(rule => {
    if (rule.startsWith('.')) return file.name.toLowerCase().endsWith(rule.toLowerCase())
    if (rule.endsWith('/*')) return file.type.startsWith(rule.slice(0, -1))
    return file.type === rule
  })
}

const selectFile = (file: File | null) => {
  if (!file) return
  let message = ''
  if (!matchesAccept(file)) message = 'ชนิดไฟล์ไม่ถูกต้อง กรุณาเลือกไฟล์ตามรูปแบบที่กำหนด'
  else if (file.size > props.maxSizeMb * 1024 * 1024) message = `ขนาดไฟล์ต้องไม่เกิน ${props.maxSizeMb}MB`
  internalError.value = message
  emit('error', message)
  if (message) return
  emit('update:modelValue', file)
  emit('change', file)
}

const handleInput = (event: Event) => {
  const input = event.target as HTMLInputElement
  selectFile(input.files?.[0] || null)
  input.value = ''
}
const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  selectFile(event.dataTransfer?.files?.[0] || null)
}
const removeFile = () => {
  internalError.value = ''
  emit('error', '')
  emit('update:modelValue', null)
  emit('change', null)
  if (fileInput.value) fileInput.value.value = ''
}
</script>

<style scoped>
.app-file-upload{display:flex;flex-direction:column;gap:7px}.app-file-upload-label{color:var(--tx-475569);font-size:.875rem;font-weight:600}.app-file-upload-required{color:var(--tx-ef4444);font-weight:600}.app-file-upload-panel{position:relative;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:8px;min-height:150px;padding:18px;border:1px dashed var(--bd-cbd5e1);border-radius:12px;background:var(--bg-f8fafc);text-align:center;cursor:pointer;transition:border-color .15s ease,background .15s ease,box-shadow .15s ease}.app-file-upload-panel:hover,.app-file-upload-panel.is-dragging{border-color:var(--bd-93c5fd);background:var(--bg-eff6ff)}.app-file-upload-panel:focus-visible{outline:none;border-color:#2563eb;box-shadow:0 0 0 3px rgba(37,99,235,.12)}.app-file-upload-panel.has-error{border-color:#ef4444;background:var(--bg-fffafa)}.app-file-upload-icon{display:grid;place-items:center;width:52px;height:52px;border:1px solid var(--bd-dbeafe);border-radius:14px;background:var(--bg-eff6ff);color:var(--tx-2563eb)}.app-file-upload-icon svg{width:26px;height:26px;fill:none;stroke:currentColor;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}.app-file-upload-icon.is-image{color:var(--tx-2563eb)}.app-file-upload-icon.is-file-text{color:var(--tx-dc2626);background:var(--bg-fef2f2);border-color:var(--bd-fecaca)}.app-file-upload-icon.is-sheet{color:var(--tx-15803d);background:var(--bg-f0fdf4);border-color:var(--bd-bbf7d0)}.app-file-upload-icon.is-archive{color:var(--tx-b45309);background:var(--bg-fffbeb);border-color:var(--bd-fde68a)}.app-file-upload-copy{display:flex;flex-direction:column;align-items:center;gap:3px;min-width:0}.app-file-upload-copy strong{max-width:100%;overflow:hidden;color:var(--tx-1e293b);font-size:.9375rem;font-weight:700;text-overflow:ellipsis;white-space:nowrap}.app-file-upload-copy span,.app-file-upload-copy small{color:var(--tx-64748b);font-size:.8125rem;font-weight:500}.app-file-upload-error{color:var(--tx-dc2626);font-size:.8125rem;font-weight:500}.app-file-upload-input{position:absolute;width:1px;height:1px;opacity:0;pointer-events:none}@media(max-width:560px){.app-file-upload-panel{min-height:140px;padding:14px}}
</style>
