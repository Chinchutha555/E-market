<template>
  <div class="app-image-upload">
    <span v-if="label" class="app-image-upload-label">
      {{ label }}
      <b v-if="requiredText" class="app-image-upload-required">* {{ requiredText }}</b>
    </span>

    <template v-if="multiple">
      <div v-if="imageFiles.length" class="app-image-upload-grid" :class="{ 'has-error': resolvedError }">
        <div v-for="(photo, index) in imageFiles" :key="`${photo.name}-${index}`" class="app-image-upload-tile has-image">
          <img v-if="photo.dataUrl" :src="photo.dataUrl" :alt="`${previewAlt} ${index + 1}`" />
          <span v-else>รูปที่ {{ index + 1 }}</span>
          <AppButton class="app-image-upload-remove" size="sm" variant="danger" :aria-label="`ลบรูปที่ ${index + 1}`" @click="removeImage(index)">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 6 12 12M18 6 6 18" /></svg>
          </AppButton>
        </div>
        <AppButton v-if="imageFiles.length < maxFiles" size="md" variant="outline" class="app-image-upload-tile app-image-upload-add" @click="openFilePicker">
          <template #prefix>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><path d="m21 15-5-5L5 21"/></svg>
          </template>
          เพิ่มรูป
        </AppButton>
      </div>
      <div
        v-else
        class="app-image-upload-panel app-image-upload-multiple-empty is-clickable"
        :class="{ 'has-error': resolvedError }"
        role="button"
        tabindex="0"
        :aria-label="selectText"
        @click="openFilePicker"
        @keydown.enter.prevent="openFilePicker"
        @keydown.space.prevent="openFilePicker"
      >
        <div class="app-image-upload-preview">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><path d="m21 15-5-5L5 21"/></svg>
        </div>
        <div class="app-image-upload-content">
          <span class="app-image-upload-prompt">{{ selectText }}</span>
          <span class="app-image-upload-hint">{{ resolvedMultipleHint }}</span>
        </div>
      </div>
      <span v-if="imageFiles.length" class="app-image-upload-hint">{{ resolvedMultipleHint }}</span>
      <span v-if="resolvedError" class="app-image-upload-error">{{ resolvedError }}</span>
    </template>

    <div v-else-if="variant === 'button'" class="app-image-upload-button-row" :class="{ 'has-error': resolvedError }">
      <div v-if="singleImage" class="app-image-upload-button-preview">
        <img :src="singleImage" :alt="previewAlt" />
      </div>
      <div class="app-image-upload-button-content">
        <AppButton size="md" variant="outline" @click="openFilePicker">
          <template #prefix>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
          </template>
          {{ singleImage ? changeText : selectText }}
        </AppButton>
        <AppButton v-if="singleImage && showRemove" size="md" variant="danger-outline" @click="removeImage()">ลบรูป</AppButton>
        <span class="app-image-upload-hint">{{ resolvedHint }}</span>
        <span v-if="resolvedError" class="app-image-upload-error">{{ resolvedError }}</span>
      </div>
    </div>

    <div
      v-else
      class="app-image-upload-panel is-clickable"
      :class="{ 'has-error': resolvedError }"
      role="button"
      tabindex="0"
      :aria-label="singleImage ? changeText : selectText"
      @click="openFilePicker"
      @keydown.enter.prevent="openFilePicker"
      @keydown.space.prevent="openFilePicker"
    >
      <div class="app-image-upload-preview" :class="{ 'has-image': singleImage }">
        <img v-if="singleImage" :src="singleImage" :alt="previewAlt" />
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z" />
          <circle cx="12" cy="13" r="4" />
        </svg>
      </div>

      <div class="app-image-upload-content">
        <span class="app-image-upload-prompt">{{ singleImage ? changeText : selectText }}</span>
        <span class="app-image-upload-hint">{{ resolvedHint }}</span>
        <span v-if="resolvedError" class="app-image-upload-error">{{ resolvedError }}</span>
        <AppButton v-if="singleImage && showRemove" size="sm" variant="danger-outline" @click.stop="removeImage()">ลบรูป</AppButton>
      </div>
    </div>
    <input ref="fileInput" type="file" :accept="accept" :multiple="multiple" class="app-image-upload-input" @change="handleFileChange" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import AppButton from './AppButton.vue'

export interface AppImageFile { name: string; dataUrl?: string }

const props = withDefaults(defineProps<{
  modelValue?: string | AppImageFile[]
  label?: string
  hint?: string
  accept?: string
  maxSizeMb?: number
  multiple?: boolean
  maxFiles?: number
  requiredText?: string
  error?: string
  selectText?: string
  changeText?: string
  previewAlt?: string
  variant?: 'panel' | 'button'
  showRemove?: boolean
}>(), {
  modelValue: '',
  label: '',
  hint: '',
  accept: 'image/jpeg,image/png,image/webp',
  maxSizeMb: 2,
  multiple: false,
  maxFiles: 5,
  requiredText: '',
  error: '',
  selectText: 'เลือกรูปถ่าย',
  changeText: 'เปลี่ยนรูปภาพ',
  previewAlt: 'ภาพตัวอย่าง',
  variant: 'panel',
  showRemove: true
})

const emit = defineEmits<{
  'update:modelValue': [value: string | AppImageFile[]]
  change: [file: File | File[] | null]
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const errorMessage = ref('')
const singleImage = computed(() => typeof props.modelValue === 'string' ? props.modelValue : '')
const imageFiles = computed(() => Array.isArray(props.modelValue) ? props.modelValue : [])
const resolvedError = computed(() => props.error || errorMessage.value)
const resolvedHint = computed(() => props.hint || `รองรับไฟล์ JPG, PNG, WebP ขนาดไม่เกิน ${props.maxSizeMb}MB`)
const resolvedMultipleHint = computed(() => props.hint || `เพิ่มได้สูงสุด ${props.maxFiles} รูป · JPG, PNG, WebP ไม่เกิน ${props.maxSizeMb}MB ต่อรูป`)

const openFilePicker = () => fileInput.value?.click()

const handleFileChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  const selectedFiles = Array.from(input.files || [])
  if (!selectedFiles.length) return

  const allowedTypes = props.accept.split(',').map(type => type.trim()).filter(Boolean)
  const invalidType = selectedFiles.find(file => !file.type.startsWith('image/') || (allowedTypes.length && !allowedTypes.includes(file.type)))
  if (invalidType) {
    errorMessage.value = 'กรุณาเลือกไฟล์รูปภาพ JPG, PNG หรือ WebP'
    input.value = ''
    return
  }
  const oversizedFile = selectedFiles.find(file => file.size > props.maxSizeMb * 1024 * 1024)
  if (oversizedFile) {
    errorMessage.value = `ขนาดไฟล์ต้องไม่เกิน ${props.maxSizeMb}MB`
    input.value = ''
    return
  }

  const availableSlots = props.multiple ? Math.max(0, props.maxFiles - imageFiles.value.length) : 1
  const files = selectedFiles.slice(0, availableSlots)
  Promise.all(files.map(file => new Promise<AppImageFile>((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve({ name: file.name, dataUrl: String(reader.result || '') })
    reader.onerror = reject
    reader.readAsDataURL(file)
  }))).then(images => {
    errorMessage.value = ''
    if (props.multiple) {
      emit('update:modelValue', [...imageFiles.value, ...images])
      emit('change', files)
    } else {
      emit('update:modelValue', images[0]?.dataUrl || '')
      emit('change', files[0] || null)
    }
  }).catch(() => { errorMessage.value = 'ไม่สามารถอ่านไฟล์รูปภาพได้ กรุณาลองใหม่' })
  input.value = ''
}

const removeImage = (index?: number) => {
  errorMessage.value = ''
  if (props.multiple) {
    emit('update:modelValue', imageFiles.value.filter((_, photoIndex) => photoIndex !== index))
  } else {
    emit('update:modelValue', '')
  }
  emit('change', null)
  if (fileInput.value) fileInput.value.value = ''
}
</script>

<style scoped>
.app-image-upload-multiple-empty{box-sizing:border-box;width:100%;min-height:180px}
.app-image-upload{display:flex;flex-direction:column;gap:7px}.app-image-upload-label{color:var(--tx-475569);font-size:var(--font-size-sm);font-weight:600}.app-image-upload-required{color:var(--tx-ef4444);font-weight:600}.app-image-upload-panel{position:relative;display:flex;align-items:center;gap:12px;padding:18px;border:1px dashed var(--bd-cbd5e1);border-radius:12px;background:var(--bg-f8fafc)}.app-image-upload-panel.is-clickable{flex-direction:column;justify-content:center;cursor:pointer;text-align:center;transition:border-color .15s ease,background .15s ease,box-shadow .15s ease}.app-image-upload-panel.is-clickable:hover{border-color:var(--bd-93c5fd);background:var(--bg-eff6ff)}.app-image-upload-panel.is-clickable:focus-visible{outline:none;border-color:#2563eb;box-shadow:0 0 0 3px rgba(37,99,235,.12)}.app-image-upload-panel.has-error,.app-image-upload-grid.has-error{border-color:#ef4444;background:var(--bg-fffafa)}.app-image-upload-preview{display:grid;place-items:center;width:72px;height:72px;flex:0 0 72px;overflow:hidden;border:1px dashed var(--bd-cbd5e1);border-radius:12px;background:var(--bg-f1f5f9);color:var(--tx-94a3b8)}.app-image-upload-preview.has-image{border-style:solid}.app-image-upload-preview svg{width:32px;height:32px}.app-image-upload-preview img{width:100%;height:100%;object-fit:cover}.app-image-upload-content{display:flex;min-width:0;align-items:center;flex-direction:column;gap:6px}.app-image-upload-prompt{color:var(--tx-334155);font-size:var(--font-size-base);font-weight:600}.app-image-upload-grid{display:flex;gap:10px;flex-wrap:wrap;padding:10px;border:1px solid transparent;border-radius:12px}.app-image-upload-tile{position:relative;box-sizing:border-box;width:112px;height:112px;overflow:hidden;border:1px solid var(--bd-dbe3ee);border-radius:12px;background:var(--bg-f8fafc);color:var(--tx-334155)}.app-image-upload-tile.has-image{display:grid;place-items:center}.app-image-upload-tile img{width:100%;height:100%;object-fit:cover}.app-image-upload-add{width:112px!important;height:112px!important;min-height:112px!important;max-height:112px!important;flex-direction:column!important;gap:7px!important;font-size:var(--font-size-sm)!important}.app-image-upload-add:hover{border-color:var(--bd-93c5fd);background:var(--bg-eff6ff);color:var(--tx-2563eb)}.app-image-upload-add :deep(.app-btn-icon){margin:0}.app-image-upload-add svg{width:28px;height:28px}.app-image-upload-remove{position:absolute;top:5px;right:5px;width:32px!important;padding:0!important}.app-image-upload-remove svg{width:15px;height:15px}.app-image-upload-hint{color:var(--tx-64748b);font-size:var(--font-size-xs);font-weight:500}.app-image-upload-error{color:var(--tx-dc2626);font-size:var(--font-size-xs);font-weight:500}.app-image-upload-button-row{display:flex;align-items:center;gap:12px;min-height:56px;padding:10px 12px;border:1px solid var(--bd-e2e8f0);border-radius:10px;background:var(--bg-f8fafc)}.app-image-upload-button-row.has-error{border-color:#ef4444;background:var(--bg-fffafa)}.app-image-upload-button-preview{width:52px;height:52px;flex:0 0 52px;overflow:hidden;border:1px solid var(--bd-dbe3ee);border-radius:8px;background:var(--bg-ffffff)}.app-image-upload-button-preview img{width:100%;height:100%;object-fit:contain}.app-image-upload-button-content{display:flex;align-items:center;flex-wrap:wrap;gap:8px;min-width:0}.app-image-upload-button-content :deep(.app-btn svg){width:17px;height:17px}.app-image-upload-button-content .app-image-upload-hint,.app-image-upload-button-content .app-image-upload-error{flex-basis:100%}.app-image-upload-input{position:absolute;width:1px;height:1px;opacity:0;pointer-events:none}@media(max-width:560px){.app-image-upload-panel{padding:14px}.app-image-upload-preview{width:64px;height:64px;flex-basis:64px}.app-image-upload-tile{width:96px;height:96px}.app-image-upload-add{width:96px!important;height:96px!important;min-height:96px!important;max-height:96px!important}.app-image-upload-button-row{align-items:flex-start}.app-image-upload-button-content{flex:1}.app-image-upload-button-content :deep(.app-btn){flex:1}}
.app-image-upload-button-row,.app-image-upload-button-row.has-error{min-height:40px;padding:0;border:0;background:transparent}
</style>
