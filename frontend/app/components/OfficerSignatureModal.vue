<template>
  <Teleport to="body">
    <div v-if="open" class="signature-overlay" @click.self="$emit('cancel')">
      <section class="signature-modal" role="dialog" aria-modal="true" aria-labelledby="signature-title">
        <AppModalHeader title="ลงนามเพื่ออนุมัติคำร้อง" subtitle="เพิ่มลายเซ็นก่อนดำเนินการอนุมัติ" icon="edit" @close="$emit('cancel')" />

        <div class="signature-body">
          <div class="signature-heading">
            <div>
              <strong>ลายเซ็นเจ้าหน้าที่</strong>
              <p>เขียนด้วยเมาส์หรือปากกาบนพื้นที่ด้านล่าง</p>
            </div>
            <div class="signature-heading-actions">
              <AppButton
                v-if="hasSavedSignature && !hasSignature"
                size="sm"
                variant="ghost"
                @click="restoreSavedSignature"
              >
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                    <polyline points="7 10 12 15 17 10"/>
                    <line x1="12" y1="15" x2="12" y2="3"/>
                  </svg>
                </template>
                ใช้ลายเซ็นที่บันทึกไว้
              </AppButton>
              <AppButton
                size="sm"
                variant="ghost"
                :disabled="!hasSignature"
                @click="clearCurrentSignature"
              >
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M3 12a9 9 0 1 0 3-6.7L3 8"/>
                    <path d="M3 3v5h5"/>
                  </svg>
                </template>
                ล้างลายเซ็น
              </AppButton>
              <span class="signature-status" :class="{ ready: hasSignature }">
                {{ signatureStatusText }}
              </span>
            </div>
          </div>

          <div class="canvas-wrap">
            <canvas
              ref="canvas"
              class="signature-canvas"
              width="900"
              height="300"
              aria-label="พื้นที่เขียนลายเซ็น"
              @pointerdown="startDrawing"
              @pointermove="draw"
              @pointerup="stopDrawing"
              @pointercancel="stopDrawing"
              @lostpointercapture="stopDrawing"
            ></canvas>
            <span v-if="!hasSignature" class="canvas-hint">ลงลายเซ็นที่นี่</span>
          </div>

          <AppImageUpload
            v-model="imageUrl"
            class="signature-upload"
            variant="button"
            :show-remove="false"
            label="แนบรูปลายเซ็น"
            select-text="เลือกรูปลายเซ็น"
            change-text="เปลี่ยนรูปลายเซ็น"
            accept="image/png,image/jpeg"
            :max-size-mb="3"
            hint="PNG หรือ JPEG ขนาดไม่เกิน 3MB"
            preview-alt="ตัวอย่างลายเซ็นที่เลือก"
          />

          <div class="signature-save-bar">
            <label class="signature-save-toggle">
              <input
                type="checkbox"
                v-model="rememberSignature"
                class="signature-checkbox"
                @change="handleRememberToggle"
              />
              <span class="signature-save-label">บันทึกลายเซ็นไว้ใช้ในการอนุมัติครั้งถัดไป</span>
            </label>
            <AppButton
              v-if="hasSavedSignature"
              size="sm"
              variant="ghost"
              class="btn-delete-saved"
              @click="deleteSavedSignature"
            >
              <template #prefix>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
              </template>
              ลบลายเซ็นที่บันทึกไว้
            </AppButton>
          </div>

          <p v-if="error" class="signature-error" role="alert">{{ error }}</p>
        </div>

        <footer class="signature-footer">
          <AppButton variant="secondary" size="md" @click="$emit('cancel')">ยกเลิก</AppButton>
          <AppButton variant="primary" size="md" @click="confirmSignature">ดำเนินการต่อ</AppButton>
        </footer>
      </section>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import AppImageUpload from './AppImageUpload.vue'

const props = defineProps<{ open: boolean }>()
useBodyScrollLock(toRef(props, 'open'))
const emit = defineEmits<{ cancel: []; confirm: [value: string] }>()

const { token } = useAuth()

const canvas = ref<HTMLCanvasElement | null>(null)
const imageUrl = ref('')
const imageBase64 = ref('')
const hasSignature = ref(false)
const drawing = ref(false)
const error = ref('')
const isUsingSaved = ref(false)
const hasSavedSignature = ref(false)
const rememberSignature = ref(true)
const isRestoring = ref(false)

const getStorageKey = (): string | null => {
  if (typeof window === 'undefined' || !token.value) return null
  try {
    const payload = token.value.split('.')[1]
    if (!payload) return null
    const claims = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
    const id = Number(claims.user_id)
    return Number.isSafeInteger(id) && id > 0 ? 'officer_saved_signature_v2_' + id : null
  } catch {
    return null
  }
}

const getSavedSignature = (): string | null => {
  const key = getStorageKey()
  if (!key) return null
  try { return localStorage.getItem(key) } catch { return null }
}

const saveSignatureToStorage = (dataUrl: string) => {
  const key = getStorageKey()
  if (!key || !dataUrl) return
  try { localStorage.setItem(key, dataUrl) } catch {}
}

const removeSavedSignatureFromStorage = () => {
  const key = getStorageKey()
  if (!key) return
  try { localStorage.removeItem(key) } catch {}
}
const normalizeDataUrl = (val: string) => {
  if (!val) return ''
  if (val.startsWith('data:image/')) return val
  return `data:image/png;base64,${val}`
}

const signatureStatusText = computed(() => {
  if (isUsingSaved.value) return 'ใช้ลายเซ็นที่บันทึกไว้'
  if (hasSignature.value) return 'เพิ่มลายเซ็นแล้ว'
  return 'ยังไม่มีลายเซ็น'
})

const drawImageToCanvas = (dataUrl: string) => {
  if (!canvas.value) return
  const el = canvas.value
  const ctx = el.getContext('2d')
  if (!ctx) return
  const img = new Image()
  img.onload = () => {
    if (!canvas.value) return
    ctx.clearRect(0, 0, el.width, el.height)
    const cw = el.width
    const ch = el.height
    const maxW = cw * 0.8
    const maxH = ch * 0.8
    const scale = Math.min(maxW / img.width, maxH / img.height, 1)
    const dw = img.width * scale
    const dh = img.height * scale
    const dx = (cw - dw) / 2
    const dy = (ch - dh) / 2
    ctx.drawImage(img, dx, dy, dw, dh)
  }
  img.src = dataUrl
}

const point = (event: PointerEvent) => {
  const el = canvas.value!
  const rect = el.getBoundingClientRect()
  return {
    x: (event.clientX - rect.left) * el.width / rect.width,
    y: (event.clientY - rect.top) * el.height / rect.height
  }
}

const startDrawing = (event: PointerEvent) => {
  if (!canvas.value) return
  if (imageUrl.value) {
    canvas.value.getContext('2d')?.clearRect(0, 0, canvas.value.width, canvas.value.height)
    hasSignature.value = false
  }
  imageUrl.value = ''
  imageBase64.value = ''
  isUsingSaved.value = false
  error.value = ''
  drawing.value = true
  canvas.value.setPointerCapture(event.pointerId)
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return
  const p = point(event)
  ctx.beginPath()
  ctx.moveTo(p.x, p.y)
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.lineWidth = 4
  ctx.strokeStyle = '#0f172a'
}

const draw = (event: PointerEvent) => {
  if (!drawing.value || !canvas.value) return
  const ctx = canvas.value.getContext('2d')
  if (!ctx) return
  const p = point(event)
  ctx.lineTo(p.x, p.y)
  ctx.stroke()
  hasSignature.value = true
}

const stopDrawing = () => {
  drawing.value = false
}

const clearCurrentSignature = () => {
  const el = canvas.value
  el?.getContext('2d')?.clearRect(0, 0, el.width, el.height)
  imageUrl.value = ''
  imageBase64.value = ''
  hasSignature.value = false
  isUsingSaved.value = false
  error.value = ''
}

const restoreSavedSignature = () => {
  const saved = getSavedSignature()
  if (!saved) return
  const dataUrl = normalizeDataUrl(saved)
  isRestoring.value = true
  imageUrl.value = dataUrl
  imageBase64.value = saved.startsWith('data:') ? (saved.split(',')[1] || '') : saved
  hasSignature.value = true
  isUsingSaved.value = true
  error.value = ''
  drawImageToCanvas(dataUrl)
  nextTick(() => { isRestoring.value = false })
}

const deleteSavedSignature = () => {
  removeSavedSignatureFromStorage()
  hasSavedSignature.value = false
  if (isUsingSaved.value) {
    clearCurrentSignature()
  }
}

const handleRememberToggle = () => {
  if (typeof window !== 'undefined') {
    const key = getStorageKey()
    if (key) localStorage.setItem(key + '_remember', String(rememberSignature.value))
  }
  if (!rememberSignature.value) {
    deleteSavedSignature()
  } else if (hasSignature.value) {
    const fullDataUrl = imageUrl.value || (canvas.value?.toDataURL('image/png') || '')
    if (fullDataUrl) {
      saveSignatureToStorage(fullDataUrl)
      hasSavedSignature.value = true
    }
  }
}

const confirmSignature = () => {
  if (!hasSignature.value) {
    error.value = 'กรุณาเขียนลายเซ็นหรือแนบรูปลายเซ็น'
    return
  }
  const data = imageBase64.value || canvas.value?.toDataURL('image/png').split(',')[1] || ''
  if (!data) {
    error.value = 'ไม่พบข้อมูลลายเซ็น'
    return
  }
  if (rememberSignature.value) {
    const fullDataUrl = imageUrl.value || (canvas.value?.toDataURL('image/png') || '')
    if (fullDataUrl) {
      saveSignatureToStorage(fullDataUrl)
      hasSavedSignature.value = true
    }
  }
  emit('confirm', data)
}

watch(imageUrl, (value) => {
  if (isRestoring.value) return
  if (value) {
    imageBase64.value = value.split(',')[1] || ''
    hasSignature.value = Boolean(imageBase64.value)
    isUsingSaved.value = false
    error.value = ''
    drawImageToCanvas(value)
    if (rememberSignature.value) {
      saveSignatureToStorage(value)
      hasSavedSignature.value = true
    }
  } else if (!drawing.value) {
    imageBase64.value = ''
    hasSignature.value = false
    isUsingSaved.value = false
    const el = canvas.value
    el?.getContext('2d')?.clearRect(0, 0, el.width, el.height)
  }
})

const initSignature = () => {
  error.value = ''
  if (typeof window !== 'undefined') {
    const key = getStorageKey()
    const pref = key ? localStorage.getItem(key + '_remember') : null
    if (pref !== null) {
      rememberSignature.value = pref === 'true'
    }
  }
  const saved = getSavedSignature()
  hasSavedSignature.value = Boolean(saved)
  if (saved && rememberSignature.value) {
    restoreSavedSignature()
  } else {
    clearCurrentSignature()
  }
}

watch(() => props.open, async (open) => {
  if (open) {
    await nextTick()
    initSignature()
  }
})
</script>

<style scoped>
.signature-overlay {
  position: fixed;
  inset: 0;
  z-index: 2500;
  display: grid;
  place-items: center;
  padding: 16px;
  background: rgba(15, 23, 42, .55);
}

.signature-modal {
  width: min(620px, 100%);
  max-height: 92vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-ffffff);
  border-radius: 16px;
  box-shadow: 0 20px 35px rgba(15, 23, 42, .2);
}

.signature-body {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  overscroll-behavior: contain;
  padding: 20px 24px 24px;
}

.signature-heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.signature-heading strong {
  color: var(--tx-1e293b);
  font-size: .9375rem;
}

.signature-heading p {
  margin: 4px 0 0;
  color: var(--tx-64748b);
  font-size: .875rem;
  line-height: 1.45;
}

.signature-status {
  flex-shrink: 0;
  padding: 4px 10px;
  border-radius: 6px;
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
  font-size: .8125rem;
  font-weight: 600;
}

.signature-status.ready {
  background: var(--bg-dcfce7);
  color: var(--tx-15803d);
}

.signature-heading-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 6px;
}

.canvas-wrap {
  position: relative;
}

.signature-canvas {
  display: block;
  width: 100%;
  height: 180px;
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 10px;
  background: var(--bg-ffffff);
  touch-action: none;
  cursor: crosshair;
}

.canvas-hint {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  color: var(--tx-94a3b8);
  font-size: .875rem;
  pointer-events: none;
}

.signature-upload {
  margin-top: 18px;
}

.signature-save-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 14px;
  padding: 10px 14px;
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-eef2f6);
  border-radius: 10px;
}

.signature-save-toggle {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.signature-checkbox {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  accent-color: var(--tx-2563eb);
  cursor: pointer;
}

.signature-save-label {
  font-size: .875rem;
  color: var(--tx-334155);
  font-weight: 500;
}

.btn-delete-saved {
  color: var(--tx-dc2626) !important;
  font-size: .8125rem !important;
}

.btn-delete-saved:hover {
  background: var(--bg-fee2e2) !important;
}

.signature-error {
  margin: 10px 0 0;
  color: var(--tx-dc2626);
  font-size: .875rem;
}

.signature-footer {
  flex: 0 0 auto;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 16px 24px;
  border-top: 1px solid var(--bd-eef2f6);
  background: var(--bg-f8fafc);
}

@media (max-width: 540px) {
  .signature-header,
  .signature-body {
    padding-left: 16px;
    padding-right: 16px;
  }
  .signature-heading {
    flex-wrap: wrap;
  }
  .signature-heading-actions {
    width: 100%;
    justify-content: flex-start;
  }
  .signature-save-bar {
    flex-direction: column;
    align-items: flex-start;
  }
  .signature-footer {
    padding: 14px 16px;
  }
}
</style>
