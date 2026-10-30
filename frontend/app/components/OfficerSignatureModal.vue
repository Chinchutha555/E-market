<template>
  <Teleport to="body">
    <div v-if="open" class="signature-overlay" @click.self="$emit('cancel')">
      <section class="signature-modal" role="dialog" aria-modal="true" aria-labelledby="signature-title">
        <header class="signature-header">
          <div class="header-content">
            <span class="header-icon" aria-hidden="true"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 19h8M5 19l3.5-.7L19 7.8a2.1 2.1 0 0 0-3-3L5.5 15.3 5 19z"/><path d="m14 7 3 3"/></svg></span>
            <div><h2 id="signature-title">ลงนามเพื่ออนุมัติคำร้อง</h2><p>เพิ่มลายเซ็นก่อนดำเนินการอนุมัติ</p></div>
          </div>
          <button type="button" class="close" aria-label="ปิด" @click="$emit('cancel')"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18M6 6l12 12"/></svg></button>
        </header>

        <div class="signature-body">
          <div class="signature-heading"><div><strong>ลายเซ็นเจ้าหน้าที่</strong><p>เขียนด้วยเมาส์หรือปากกาบนพื้นที่ด้านล่าง</p></div><span class="signature-status" :class="{ ready: hasSignature }">{{ hasSignature ? 'เพิ่มลายเซ็นแล้ว' : 'ยังไม่มีลายเซ็น' }}</span></div>
          <div class="canvas-wrap">
            <canvas ref="canvas" class="signature-canvas" width="900" height="300" aria-label="พื้นที่เขียนลายเซ็น" @pointerdown="startDrawing" @pointermove="draw" @pointerup="stopDrawing" @pointercancel="stopDrawing" @lostpointercapture="stopDrawing"></canvas>
            <span v-if="!hasSignature" class="canvas-hint">ลงลายเซ็นที่นี่</span>
          </div>
          <div v-if="imageUrl" class="uploaded"><span>รูปภาพลายเซ็นที่เลือก</span><img :src="imageUrl" alt="ตัวอย่างลายเซ็นที่เลือก" /></div>
          <div class="signature-tools">
            <div><AppButton size="sm" variant="outline" @click="fileInput?.click()"><template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 16V4m-4 4 4-4 4 4"/><path d="M4 16v3a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-3"/></svg></template>เลือกรูปลายเซ็น</AppButton><span class="upload-hint">PNG หรือ JPEG ไม่เกิน 3 MB</span></div>
            <AppButton size="sm" variant="ghost" :disabled="!hasSignature" @click="clearSignature"><template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 3-6.7L3 8"/><path d="M3 3v5h5"/></svg></template>ล้างลายเซ็น</AppButton>
            <input ref="fileInput" type="file" accept="image/png,image/jpeg" hidden @change="selectFile" />
          </div>
          <p v-if="error" class="signature-error" role="alert">{{ error }}</p>
        </div>

        <footer class="signature-footer"><AppButton variant="secondary" @click="$emit('cancel')">ยกเลิก</AppButton><AppButton variant="primary" :disabled="!hasSignature" @click="confirmSignature">ดำเนินการต่อ</AppButton></footer>
      </section>
    </div>
  </Teleport>
</template>
<script setup lang="ts">
const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ cancel: []; confirm: [value: string] }>()
const canvas = ref<HTMLCanvasElement | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const imageUrl = ref('')
const imageBase64 = ref('')
const hasSignature = ref(false)
const drawing = ref(false)
const error = ref('')
const point = (event: PointerEvent) => { const el = canvas.value!; const rect = el.getBoundingClientRect(); return { x: (event.clientX - rect.left) * el.width / rect.width, y: (event.clientY - rect.top) * el.height / rect.height } }
const startDrawing = (event: PointerEvent) => { if (!canvas.value) return; if (imageUrl.value) { canvas.value.getContext('2d')?.clearRect(0, 0, canvas.value.width, canvas.value.height); hasSignature.value = false } imageUrl.value = ''; imageBase64.value = ''; error.value = ''; drawing.value = true; canvas.value.setPointerCapture(event.pointerId); const ctx = canvas.value.getContext('2d'); if (!ctx) return; const p = point(event); ctx.beginPath(); ctx.moveTo(p.x,p.y); ctx.lineCap = 'round'; ctx.lineJoin = 'round'; ctx.lineWidth = 4; ctx.strokeStyle = '#0f172a' }
const draw = (event: PointerEvent) => { if (!drawing.value || !canvas.value) return; const ctx = canvas.value.getContext('2d'); if (!ctx) return; const p = point(event); ctx.lineTo(p.x,p.y); ctx.stroke(); hasSignature.value = true }
const stopDrawing = () => { drawing.value = false }
const clearSignature = () => { const el = canvas.value; el?.getContext('2d')?.clearRect(0,0,el.width,el.height); imageUrl.value = ''; imageBase64.value = ''; hasSignature.value = false; error.value = ''; if (fileInput.value) fileInput.value.value = '' }
const selectFile = (event: Event) => { const file = (event.target as HTMLInputElement).files?.[0]; if (!file) return; if (!['image/png','image/jpeg'].includes(file.type) || file.size > 3 * 1024 * 1024) { clearSignature(); error.value = 'กรุณาเลือกไฟล์ PNG หรือ JPEG ขนาดไม่เกิน 3 MB'; return } const reader = new FileReader(); reader.onload = () => { const data = String(reader.result || ''); imageUrl.value = data; imageBase64.value = data.split(',')[1] || ''; hasSignature.value = !!imageBase64.value }; reader.onerror = () => { error.value = 'ไม่สามารถอ่านไฟล์ได้' }; reader.readAsDataURL(file) }
const confirmSignature = () => { if (!hasSignature.value) return; const data = imageBase64.value || canvas.value?.toDataURL('image/png').split(',')[1] || ''; if (data) emit('confirm', data) }
watch(() => props.open, async (open) => { if (open) { await nextTick(); clearSignature() } })
</script>
<style scoped>
.signature-overlay { position: fixed; inset: 0; z-index: 2500; display: grid; place-items: center; padding: 16px; background: rgba(15,23,42,.55); }
.signature-modal { width: min(620px,100%); max-height: 92vh; display: flex; flex-direction: column; overflow: hidden; background: #fff; border-radius: 16px; box-shadow: 0 20px 35px rgba(15,23,42,.2); }
.signature-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; padding: 20px 24px; border-bottom: 1px solid #eef2f6; }
.header-content { display: flex; align-items: center; gap: 12px; min-width: 0; } .header-icon { width: 42px; height: 42px; flex: 0 0 42px; display: grid; place-items: center; border-radius: 11px; color: #2563eb; background: #eff6ff; }
h2 { margin: 0; color: #1e293b; font-size: 1.125rem; font-weight: 700; } .header-content p, .signature-heading p { margin: 4px 0 0; color: #64748b; font-size: .875rem; line-height: 1.45; }
.close { width: 34px; height: 34px; flex: 0 0 34px; display: grid; place-items: center; border: 0; border-radius: 8px; background: #f1f5f9; color: #475569; cursor: pointer; } .close:hover { background: #e2e8f0; }
.signature-body { overflow: auto; padding: 20px 24px 24px; } .signature-heading { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; margin-bottom: 12px; } .signature-heading strong { color: #1e293b; font-size: .9375rem; }
.signature-status { flex-shrink: 0; padding: 4px 10px; border-radius: 6px; background: #f1f5f9; color: #64748b; font-size: .8125rem; font-weight: 600; } .signature-status.ready { background: #dcfce7; color: #15803d; }
.canvas-wrap { position: relative; } .signature-canvas { display: block; width: 100%; height: 180px; border: 1px solid #cbd5e1; border-radius: 10px; background: #fff; touch-action: none; cursor: crosshair; } .canvas-hint { position: absolute; left: 50%; top: 50%; transform: translate(-50%,-50%); color: #94a3b8; font-size: .875rem; pointer-events: none; }
.uploaded { margin-top: 12px; padding: 12px; border: 1px solid #e2e8f0; border-radius: 10px; background: #f8fafc; } .uploaded span { display: block; margin-bottom: 6px; color: #64748b; font-size: .8125rem; } .uploaded img { display: block; max-width: 100%; max-height: 100px; object-fit: contain; }
.signature-tools { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 10px; margin-top: 14px; } .signature-tools > div { display: flex; align-items: center; flex-wrap: wrap; gap: 10px; } .upload-hint { color: #64748b; font-size: .8125rem; }
.signature-error { margin: 10px 0 0; color: #dc2626; font-size: .875rem; } .signature-footer { display: flex; justify-content: flex-end; gap: 10px; padding: 16px 24px; border-top: 1px solid #eef2f6; background: #f8fafc; }
@media (max-width: 540px) { .signature-header, .signature-body { padding-left: 16px; padding-right: 16px; } .signature-heading { flex-wrap: wrap; } .signature-tools { align-items: flex-start; } .signature-footer { padding: 14px 16px; } }
</style>
