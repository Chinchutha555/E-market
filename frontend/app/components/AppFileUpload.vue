<template>
  <div class="app-file-upload">
    <span v-if="label" class="app-file-upload-label">
      {{ label }}
      <b v-if="requiredText" class="app-file-upload-required">* {{ requiredText }}</b>
    </span>

    <!-- Mode 1: Premium Document Card View (When file is attached) -->
    <div
      v-if="modelValue && (imagePreviewUrl || isPdf) && showPreview"
      class="app-file-upload-card"
      :class="{ 'has-error': resolvedError }"
    >
      <!-- Card Header -->
      <div class="doc-card-header">
        <div class="doc-card-meta">
          <div class="doc-icon-badge" :class="isPdf ? 'is-pdf' : 'is-img'">
            <svg v-if="!isPdf" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="18" height="18" rx="2"/>
              <circle cx="8.5" cy="8.5" r="1.5"/>
              <path d="m21 15-5-5L5 21"/>
            </svg>
            <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <path d="M14 2v6h6M8 13h8M8 17h8"/>
            </svg>
          </div>
          <div class="doc-title-group">
            <span class="doc-name" :title="modelValue.name">{{ modelValue.name }}</span>
            <div class="doc-sub-info">
              <span class="doc-size">{{ formatFileSize(modelValue.size) }}</span>
            </div>
          </div>
        </div>

        <div class="doc-header-actions">
          <AppButton size="sm" variant="outline" @click="openFilePicker">
            <template #prefix>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
            </template>
            เปลี่ยนไฟล์
          </AppButton>
          <AppButton size="sm" variant="danger-outline" @click="removeFile">
            <template #prefix>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
            </template>
            ลบไฟล์
          </AppButton>
        </div>
      </div>

      <!-- Card Body / Image Preview Showcase -->
      <div
        class="doc-preview-viewport"
        role="button"
        tabindex="0"
        title="คลิกเพื่อดูเอกสารขนาดเต็ม"
        @click="showFullModal = true"
        @keydown.enter.prevent="showFullModal = true"
        @keydown.space.prevent="showFullModal = true"
      >
        <div v-if="imagePreviewUrl" class="doc-img-frame">
          <img
            :src="imagePreviewUrl"
            :alt="modelValue.name"
            class="doc-img-display"
          />
          <div class="doc-img-overlay">
            <span class="doc-zoom-btn">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="11" cy="11" r="8"></circle>
                <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                <line x1="11" y1="8" x2="11" y2="14"></line>
                <line x1="8" y1="11" x2="14" y2="11"></line>
              </svg>
              คลิกเพื่อดูภาพขนาดเต็ม
            </span>
          </div>
        </div>
        <div v-else-if="isPdf" class="doc-pdf-frame">
          <div class="doc-pdf-icon-lg">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <path d="M14 2v6h6M8 13h8M8 17h8"/>
            </svg>
            <span>เอกสาร PDF</span>
          </div>
          <span class="doc-zoom-btn">คลิกเพื่อเปิดดูตัวอย่าง PDF</span>
        </div>
      </div>

      <!-- Card Footer Helper -->
      <div class="doc-card-footer">
        <span class="doc-footer-hint">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="16" x2="12" y2="12"></line>
            <line x1="12" y1="8" x2="12.01" y2="8"></line>
          </svg>
          กรุณาตรวจสอบความชัดเจนของเอกสารก่อนยื่นคำร้อง
        </span>
      </div>

      <input ref="fileInput" type="file" :accept="accept" class="app-file-upload-input" @change="handleInput" />
    </div>

    <!-- Mode 2: Standard Panel Mode (Empty state or generic files) -->
    <div
      v-else
      class="app-file-upload-panel"
      :class="{ 'has-error': resolvedError, 'is-dragging': isDragging, 'has-file': Boolean(modelValue) }"
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
        <svg v-if="fileKind === 'image'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><path d="m21 15-5-5L5 21"/></svg>
        <svg v-else-if="fileKind === 'sheet'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8M11 9v12"/></svg>
        <svg v-else-if="fileKind === 'archive'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 8v13H3V8M1 3h22v5H1zM10 12h4"/></svg>
        <svg v-else-if="fileKind === 'file-text'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/></svg>
      </div>
      <div class="app-file-upload-copy">
        <strong>{{ modelValue?.name || selectText }}</strong>
        <span v-if="modelValue">{{ formatFileSize(modelValue.size) }} · คลิกเพื่อเปลี่ยนไฟล์</span>
        <span v-else>{{ description || 'คลิกเพื่อเลือกไฟล์ หรือลากไฟล์มาวางที่นี่' }}</span>
        <small>{{ resolvedHint }}</small>
      </div>

      <div v-if="modelValue" class="app-file-upload-actions" @click.stop>
        <AppButton size="sm" variant="outline" @click="openFilePicker">
          <template #prefix>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="17 8 12 3 7 8"/>
              <line x1="12" y1="3" x2="12" y2="15"/>
            </svg>
          </template>
          เปลี่ยนไฟล์
        </AppButton>
        <AppButton size="sm" variant="danger-outline" @click="removeFile">
          <template #prefix>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="3 6 5 6 21 6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
          </template>
          ลบไฟล์
        </AppButton>
      </div>

      <input ref="fileInput" type="file" :accept="accept" class="app-file-upload-input" @change="handleInput" />
    </div>

    <span v-if="resolvedError" class="app-file-upload-error">{{ resolvedError }}</span>

    <!-- Modal Lightbox for full preview -->
    <Teleport to="body">
      <div v-if="showFullModal" class="app-file-upload-modal-overlay" @click.self="showFullModal = false">
        <section class="app-file-upload-modal" role="dialog" aria-modal="true" aria-labelledby="preview-file-title">
          <AppModalHeader
            :title="label || 'ตัวอย่างไฟล์ที่แนบ'"
            :subtitle="modelValue?.name || ''"
            icon="document"
            @close="showFullModal = false"
          />
          <div class="app-file-upload-modal-body">
            <img
              v-if="imagePreviewUrl"
              :src="imagePreviewUrl"
              :alt="modelValue?.name || 'ตัวอย่างไฟล์'"
              class="app-file-upload-modal-img"
            />
            <iframe
              v-else-if="pdfPreviewUrl"
              :src="pdfPreviewUrl"
              class="app-file-upload-modal-iframe"
              title="ตัวอย่างเอกสาร PDF"
            ></iframe>
            <p v-else class="app-file-upload-no-preview">ไม่สามารถแสดงตัวอย่างไฟล์นี้ได้</p>
          </div>
          <footer class="app-file-upload-modal-footer">
            <AppButton variant="secondary" size="md" @click="showFullModal = false">ปิด</AppButton>
          </footer>
        </section>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
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
  showPreview?: boolean
}>(), {
  modelValue: null,
  label: '',
  requiredText: '',
  accept: '',
  maxSizeMb: 10,
  hint: '',
  description: '',
  selectText: 'เลือกไฟล์',
  error: '',
  showPreview: true
})

const emit = defineEmits<{
  'update:modelValue': [file: File | null]
  change: [file: File | null]
  error: [message: string]
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const internalError = ref('')
const isDragging = ref(false)
const showFullModal = ref(false)
const imagePreviewUrl = ref('')
const pdfPreviewUrl = ref('')

const resolvedError = computed(() => props.error || internalError.value)
const resolvedHint = computed(() => props.hint || `ขนาดไฟล์ไม่เกิน ${props.maxSizeMb}MB`)

const fileKind = computed<'image' | 'file-text' | 'sheet' | 'archive' | 'paperclip'>(() => {
  const type = props.modelValue?.type || ''
  const name = props.modelValue?.name.toLowerCase() || ''
  if (type.startsWith('image/') || /\.(jpg|jpeg|png|webp|gif|svg)$/i.test(name)) return 'image'
  if (type === 'application/pdf' || /\.(pdf|doc|docx|txt)$/.test(name) || /pdf|document|text/.test(props.accept)) return 'file-text'
  if (/spreadsheet|csv|excel/.test(type) || /\.(csv|xls|xlsx)$/.test(name)) return 'sheet'
  if (/zip|compressed|archive/.test(type) || /\.(zip|rar|7z)$/.test(name)) return 'archive'
  return 'paperclip'
})

const isPdf = computed(() => {
  const type = props.modelValue?.type || ''
  const name = props.modelValue?.name.toLowerCase() || ''
  return type === 'application/pdf' || name.endsWith('.pdf')
})

const revokeUrls = () => {
  if (imagePreviewUrl.value && imagePreviewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(imagePreviewUrl.value)
    imagePreviewUrl.value = ''
  }
  if (pdfPreviewUrl.value && pdfPreviewUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(pdfPreviewUrl.value)
    pdfPreviewUrl.value = ''
  }
}

watch(() => props.modelValue, (file) => {
  revokeUrls()
  if (!file) return
  if (typeof window !== 'undefined') {
    if (file.type.startsWith('image/') || /\.(jpg|jpeg|png|webp|gif|svg)$/i.test(file.name)) {
      try {
        imagePreviewUrl.value = URL.createObjectURL(file)
      } catch {
        imagePreviewUrl.value = ''
      }
    } else if (file.type === 'application/pdf' || file.name.toLowerCase().endsWith('.pdf')) {
      try {
        pdfPreviewUrl.value = URL.createObjectURL(file)
      } catch {
        pdfPreviewUrl.value = ''
      }
    }
  }
}, { immediate: true })

onBeforeUnmount(() => {
  revokeUrls()
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
  revokeUrls()
  internalError.value = ''
  emit('error', '')
  emit('update:modelValue', null)
  emit('change', null)
  if (fileInput.value) fileInput.value.value = ''
}
</script>

<style scoped>
.app-file-upload {
  display: flex;
  flex-direction: column;
  gap: 7px;
}
.app-file-upload-label {
  color: var(--tx-475569);
  font-size: .875rem;
  font-weight: 600;
}
.app-file-upload-required {
  color: var(--tx-ef4444);
  font-weight: 600;
}

/* Card Mode */
.app-file-upload-card {
  display: flex;
  flex-direction: column;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
  transition: border-color .18s ease, box-shadow .18s ease;
}
.app-file-upload-card:hover {
  border-color: var(--bd-93c5fd);
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.08);
}
.app-file-upload-card.has-error {
  border-color: #ef4444;
  background: var(--bg-fffafa);
}
.doc-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--bg-f8fafc);
  border-bottom: 1px solid var(--bd-eef2f6);
  flex-wrap: wrap;
  gap: 10px;
}
.doc-card-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
  flex: 1;
}
.doc-icon-badge {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  flex-shrink: 0;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  border: 1px solid var(--bd-dbeafe);
}
.doc-icon-badge.is-pdf {
  background: var(--bg-fef2f2);
  color: var(--tx-dc2626);
  border-color: var(--bd-fecaca);
}
.doc-title-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.doc-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--tx-1e293b);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 260px;
}
.doc-sub-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.75rem;
}
.doc-size {
  color: var(--tx-64748b);
  font-weight: 500;
}
.doc-status-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 1px 7px;
  border-radius: 6px;
  background: var(--bg-dcfce7);
  color: var(--tx-15803d);
  font-weight: 600;
  font-size: 0.75rem;
}
.doc-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.doc-preview-viewport {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 180px;
  max-height: 260px;
  padding: 14px;
  background: #f1f5f9;
  cursor: pointer;
  overflow: hidden;
}
.doc-img-frame {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  max-width: 100%;
  max-height: 232px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.1);
  background: #ffffff;
}
.doc-img-display {
  display: block;
  max-width: 100%;
  max-height: 232px;
  object-fit: contain;
}
.doc-img-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.35);
  opacity: 0;
  transition: opacity .18s ease;
}
.doc-preview-viewport:hover .doc-img-overlay {
  opacity: 1;
}
.doc-zoom-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  background: rgba(15, 23, 42, 0.85);
  color: #ffffff;
  font-size: 0.8125rem;
  font-weight: 600;
  backdrop-filter: blur(4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}
.doc-pdf-frame {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 20px;
}
.doc-pdf-icon-lg {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  color: var(--tx-dc2626);
  font-weight: 600;
  font-size: 0.875rem;
}

.doc-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background: var(--bg-ffffff);
  border-top: 1px solid var(--bd-eef2f6);
  flex-wrap: wrap;
  gap: 8px;
}
.doc-footer-hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--tx-64748b);
  font-size: 0.75rem;
  font-weight: 500;
}

/* Panel Mode (Empty State) */
.app-file-upload-panel {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  min-height: 150px;
  padding: 18px;
  border: 1px dashed var(--bd-cbd5e1);
  border-radius: 12px;
  background: var(--bg-f8fafc);
  text-align: center;
  cursor: pointer;
  transition: border-color .15s ease, background .15s ease, box-shadow .15s ease;
}
.app-file-upload-panel:hover,
.app-file-upload-panel.is-dragging {
  border-color: var(--bd-93c5fd);
  background: var(--bg-eff6ff);
}
.app-file-upload-panel:focus-visible {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, .12);
}
.app-file-upload-panel.has-error {
  border-color: #ef4444;
  background: var(--bg-fffafa);
}
.app-file-upload-icon {
  display: grid;
  place-items: center;
  width: 52px;
  height: 52px;
  border: 1px solid var(--bd-dbeafe);
  border-radius: 14px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}
.app-file-upload-icon svg {
  width: 26px;
  height: 26px;
}
.app-file-upload-icon.is-image {
  color: var(--tx-2563eb);
}
.app-file-upload-icon.is-file-text {
  color: var(--tx-dc2626);
  background: var(--bg-fef2f2);
  border-color: var(--bd-fecaca);
}
.app-file-upload-icon.is-sheet {
  color: var(--tx-15803d);
  background: var(--bg-f0fdf4);
  border-color: var(--bd-bbf7d0);
}
.app-file-upload-icon.is-archive {
  color: var(--tx-b45309);
  background: var(--bg-fffbeb);
  border-color: var(--bd-fde68a);
}

.app-file-upload-copy {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  min-width: 0;
}
.app-file-upload-copy strong {
  max-width: 100%;
  overflow: hidden;
  color: var(--tx-1e293b);
  font-size: .9375rem;
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.app-file-upload-copy span,
.app-file-upload-copy small {
  color: var(--tx-64748b);
  font-size: .8125rem;
  font-weight: 500;
}
.app-file-upload-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 4px;
}
.app-file-upload-error {
  color: var(--tx-dc2626);
  font-size: .8125rem;
  font-weight: 500;
}
.app-file-upload-input {
  position: absolute;
  width: 1px;
  height: 1px;
  opacity: 0;
  pointer-events: none;
}

/* Lightbox Modal */
.app-file-upload-modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 3000;
  display: grid;
  place-items: center;
  padding: 16px;
  background: rgba(15, 23, 42, .65);
  backdrop-filter: blur(2px);
}
.app-file-upload-modal {
  width: min(800px, 100%);
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-ffffff);
  border-radius: 16px;
  box-shadow: 0 20px 35px rgba(15, 23, 42, .25);
}
.app-file-upload-modal-body {
  flex: 1 1 auto;
  min-height: 200px;
  max-height: 70vh;
  overflow: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: var(--bg-f8fafc);
}
.app-file-upload-modal-img {
  max-width: 100%;
  max-height: 65vh;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 4px 14px rgba(0, 0, 0, .08);
}
.app-file-upload-modal-iframe {
  width: 100%;
  height: 65vh;
  border: 0;
  border-radius: 8px;
}
.app-file-upload-no-preview {
  color: var(--tx-64748b);
  font-size: .875rem;
}
.app-file-upload-modal-footer {
  flex: 0 0 auto;
  display: flex;
  justify-content: flex-end;
  padding: 14px 20px;
  border-top: 1px solid var(--bd-eef2f6);
  background: var(--bg-ffffff);
}

@media (max-width: 560px) {
  .doc-card-header {
    flex-direction: column;
    align-items: flex-start;
  }
  .doc-header-actions {
    width: 100%;
    justify-content: flex-end;
  }
  .doc-card-footer {
    flex-direction: column;
    align-items: flex-start;
  }
  .app-file-upload-panel {
    min-height: 140px;
    padding: 14px;
  }
}
</style>
