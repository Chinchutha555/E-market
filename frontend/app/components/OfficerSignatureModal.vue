<template>
  <Teleport to="body">
    <div v-if="open" class="signature-overlay" @click.self="$emit('cancel')">
      <section class="signature-modal" role="dialog" aria-modal="true" aria-labelledby="signature-title">
        <AppModalHeader title="ลงนามเพื่ออนุมัติคำร้อง" subtitle="เพิ่มลายเซ็นก่อนดำเนินการอนุมัติ" icon="edit" @close="$emit('cancel')" />

        <div class="signature-body">
          <div class="signature-heading"><div><strong>ลายเซ็นเจ้าหน้าที่</strong><p>เขียนด้วยเมาส์หรือปากกาบนพื้นที่ด้านล่าง</p></div><div class="signature-heading-actions"><AppButton size="sm" variant="ghost" :disabled="!hasSignature" @click="clearSignature"><template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 3-6.7L3 8"/><path d="M3 3v5h5"/></svg></template>ล้างลายเซ็น</AppButton><span class="signature-status" :class="{ ready: hasSignature }">{{ hasSignature ? 'เพิ่มลายเซ็นแล้ว' : 'ยังไม่มีลายเซ็น' }}</span></div></div>
          <div class="canvas-wrap">
            <canvas ref="canvas" class="signature-canvas" width="900" height="300" aria-label="พื้นที่เขียนลายเซ็น" @pointerdown="startDrawing" @pointermove="draw" @pointerup="stopDrawing" @pointercancel="stopDrawing" @lostpointercapture="stopDrawing"></canvas>
            <span v-if="!hasSignature" class="canvas-hint">ลงลายเซ็นที่นี่</span>
          </div>
          <AppImageUpload v-model="imageUrl" class="signature-upload" variant="button" label="แนบรูปลายเซ็น" select-text="เลือกรูปลายเซ็น" change-text="เปลี่ยนรูปลายเซ็น" accept="image/png,image/jpeg" :max-size-mb="3" hint="PNG หรือ JPEG ขนาดไม่เกิน 3MB" preview-alt="ตัวอย่างลายเซ็นที่เลือก" />
          <p v-if="error" class="signature-error" role="alert">{{ error }}</p>
        </div>

        <footer class="signature-footer"><AppButton variant="secondary" size="md" @click="$emit('cancel')">ยกเลิก</AppButton><AppButton variant="primary" size="md" @click="confirmSignature">ดำเนินการต่อ</AppButton></footer>
      </section>
    </div>
  </Teleport>
</template>
<script setup lang="ts">
import AppImageUpload from './AppImageUpload.vue'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ cancel: []; confirm: [value: string] }>()
const canvas = ref<HTMLCanvasElement | null>(null)
const imageUrl = ref('')
const imageBase64 = ref('')
const hasSignature = ref(false)
const drawing = ref(false)
const error = ref('')
const point = (event: PointerEvent) => { const el = canvas.value!; const rect = el.getBoundingClientRect(); return { x: (event.clientX - rect.left) * el.width / rect.width, y: (event.clientY - rect.top) * el.height / rect.height } }
const startDrawing = (event: PointerEvent) => { if (!canvas.value) return; if (imageUrl.value) { canvas.value.getContext('2d')?.clearRect(0, 0, canvas.value.width, canvas.value.height); hasSignature.value = false } imageUrl.value = ''; imageBase64.value = ''; error.value = ''; drawing.value = true; canvas.value.setPointerCapture(event.pointerId); const ctx = canvas.value.getContext('2d'); if (!ctx) return; const p = point(event); ctx.beginPath(); ctx.moveTo(p.x,p.y); ctx.lineCap = 'round'; ctx.lineJoin = 'round'; ctx.lineWidth = 4; ctx.strokeStyle = '#0f172a' }
const draw = (event: PointerEvent) => { if (!drawing.value || !canvas.value) return; const ctx = canvas.value.getContext('2d'); if (!ctx) return; const p = point(event); ctx.lineTo(p.x,p.y); ctx.stroke(); hasSignature.value = true }
const stopDrawing = () => { drawing.value = false }
const clearSignature = () => { const el = canvas.value; el?.getContext('2d')?.clearRect(0,0,el.width,el.height); imageUrl.value = ''; imageBase64.value = ''; hasSignature.value = false; error.value = '' }
const confirmSignature = () => { if (!hasSignature.value) { error.value = 'กรุณาเขียนลายเซ็นหรือแนบรูปลายเซ็น'; return } const data = imageBase64.value || canvas.value?.toDataURL('image/png').split(',')[1] || ''; if (data) emit('confirm', data) }
watch(imageUrl, (value) => { if (value) { imageBase64.value = value.split(',')[1] || ''; hasSignature.value = Boolean(imageBase64.value); error.value = '' } })
watch(() => props.open, async (open) => { if (open) { await nextTick(); clearSignature() } })
</script>
<style scoped>
.signature-overlay { position: fixed; inset: 0; z-index: 2500; display: grid; place-items: center; padding: 16px; background: rgba(15,23,42,.55); }
.signature-modal { width: min(620px,100%); max-height: 92vh; display: flex; flex-direction: column; overflow: hidden; background: var(--bg-ffffff); border-radius: 16px; box-shadow: 0 20px 35px rgba(15,23,42,.2); }
.signature-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; padding: 20px 24px; border-bottom: 1px solid var(--bd-eef2f6); }
.header-content { display: flex; align-items: center; gap: 12px; min-width: 0; } .header-icon { width: 42px; height: 42px; flex: 0 0 42px; display: grid; place-items: center; border-radius: 11px; color: var(--tx-2563eb); background: var(--bg-eff6ff); }
h2 { margin: 0; color: var(--tx-1e293b); font-size: var(--font-size-lg); font-weight: 700; } .header-content p, .signature-heading p { margin: 4px 0 0; color: var(--tx-64748b); font-size: var(--font-size-sm); line-height: 1.45; }
.close { width: 34px; height: 34px; flex: 0 0 34px; display: grid; place-items: center; border: 0; border-radius: 8px; background: var(--bg-f1f5f9); color: var(--tx-475569); cursor: pointer; } .close:hover { background: var(--bg-e2e8f0); }
.signature-body { flex:1 1 auto;min-height:0;overflow: auto;overscroll-behavior:contain;padding: 20px 24px 24px; } .signature-heading { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; margin-bottom: 12px; } .signature-heading strong { color: var(--tx-1e293b); font-size: var(--font-size-base); }
.signature-status { flex-shrink: 0; padding: 4px 10px; border-radius: 6px; background: var(--bg-f1f5f9); color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 600; } .signature-status.ready { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.signature-heading-actions { display: flex; align-items: center; justify-content: flex-end; flex-wrap: wrap; gap: 6px; }
.canvas-wrap { position: relative; } .signature-canvas { display: block; width: 100%; height: 180px; border: 1px solid var(--bd-cbd5e1); border-radius: 10px; background: var(--bg-ffffff); touch-action: none; cursor: crosshair; } .canvas-hint { position: absolute; left: 50%; top: 50%; transform: translate(-50%,-50%); color: var(--tx-94a3b8); font-size: var(--font-size-sm); pointer-events: none; }
.signature-upload { margin-top: 18px; }
.signature-error { margin: 10px 0 0; color: var(--tx-dc2626); font-size: var(--font-size-sm); } .signature-footer { flex:0 0 auto;display: flex; justify-content: flex-end; gap: 10px; padding: 16px 24px; border-top: 1px solid var(--bd-eef2f6); background: var(--bg-f8fafc); }
@media (max-width: 540px) { .signature-header, .signature-body { padding-left: 16px; padding-right: 16px; } .signature-heading { flex-wrap: wrap; } .signature-heading-actions { width: 100%; justify-content: flex-start; } .signature-footer { padding: 14px 16px; } }
</style>
