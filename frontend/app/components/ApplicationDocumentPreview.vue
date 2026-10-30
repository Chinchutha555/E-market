<template>
  <Teleport to="body">
    <div v-if="application" class="document-overlay" @click.self="$emit('close')">
      <section class="document-modal" role="dialog" aria-modal="true" aria-labelledby="document-title">
        <header class="document-header"><h2 id="document-title">ตัวอย่างคำร้องขออนุญาตเช่าแผงค้า</h2><button type="button" class="document-close" aria-label="ปิดตัวอย่างเอกสาร" @click="$emit('close')"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18M6 6l12 12"/></svg></button></header>
        <div class="document-tools"><div class="zoom-controls"><AppButton size="sm" variant="outline" :disabled="zoom <= 70" @click="zoom = Math.max(70, zoom - 10)">−</AppButton><span>{{ zoom }}%</span><AppButton size="sm" variant="outline" :disabled="zoom >= 180" @click="zoom = Math.min(180, zoom + 10)">+</AppButton><AppButton size="sm" variant="ghost" @click="zoom = 100">ขนาดปกติ</AppButton></div></div>
        <div class="document-scroll">
          <article class="document-paper" :style="{ zoom: zoom / 100 }">
            <div class="document-number">คำร้องขออนุญาตที่ {{ application.document_number || '.......... / ..........' }}</div>
            <h3>คำร้องขออนุญาตเช่าแผงค้า สำนักงานตลาดกรุงเทพมหานคร</h3>
            <div class="document-date">
              <div>เขียนที่ <span class="fill-line place">{{ application.market_name || '' }}</span></div>
              <div>วันที่ <span class="fill-line date-day">{{ submittedDate.day }}</span> เดือน <span class="fill-line date-month">{{ submittedDate.month }}</span> พ.ศ. <span class="fill-line date-year">{{ submittedDate.year }}</span></div>
              <div>เวลาที่ยื่นคำร้อง <span class="fill-line date-time">{{ submittedDate.time }}</span></div>
            </div>
            <p><strong>เรียน</strong> ผู้อำนวยการสำนักงานตลาดกรุงเทพมหานคร</p>
            <section class="document-section">
              <strong>1. ข้อมูลผู้ยื่นคำร้องขออนุญาต</strong>
              <div class="document-row">ข้าพเจ้า <span class="document-value">{{ applicantTitle || '(ไม่ระบุคำนำหน้า)' }}</span> <span class="document-value">{{ applicantName || '........................................' }}</span> อายุ <span class="document-value">{{ application.applicant_age || '........' }}</span> ปี</div>
              <div class="document-row">เลขประจำตัวประชาชน <span class="document-value">{{ application.applicant_id_card || '........................................' }}</span> โทร. <span class="document-value">{{ application.applicant_phone || '................................' }}</span></div>
              <strong>เอกสารที่ยื่นประกอบคำร้องขออนุญาต</strong>
              <div class="document-indent"><span class="form-check" :class="{ marked: !!application.booking_number }"></span> ใบจองจากระบบ เลขที่ {{ application.booking_number || '........................................' }}</div>
              <div class="document-indent">ประเภท <span class="form-check" :class="{ marked: checked('trial') }"></span> ทดลองขาย <span class="form-check" :class="{ marked: checked('temporary') }"></span> ทำการค้าแบบชั่วคราว <span class="form-check" :class="{ marked: checked('contract') }"></span> ทำการค้าแบบมีสัญญา</div>
              <div class="document-indent"><span class="form-check" :class="{ marked: application.has_id_card_copy }"></span> สำเนาบัตรประจำตัวประชาชน</div>
              <p>ข้าพเจ้าได้อ่านและเข้าใจหลักเกณฑ์ เงื่อนไข และวิธีการให้เช่า ทั้งหมดแล้วทุกประการ</p>
              <div class="document-sign">ลงชื่อ <span class="fill-line signature">{{ savedSignature }}</span> ผู้ขอสิทธิ์<br><span class="tenant-name">( {{ application.applicant_name || '..................................................' }} )</span></div>
            </section>
            <section class="document-section"><strong>2. ส่วนของเจ้าหน้าที่</strong><br><strong>2.1 งานผลประโยชน์และเร่งรัดหนี้สิน/ผู้จัดการตลาด</strong><div class="document-indent">□ ไม่มีแผงค้า<br>□ มีแผงค้า จำนวน ............ แผง เลขที่ ........................ โครงการ ........................ ประเภท ........................<br>□ สัญญาอื่น ๆ (เช่า) สำนักงานตลาดกรุงเทพมหานคร<br>□ อื่น ๆ ..............................................................................................<br>□ คุณสมบัติและเอกสารถูกต้องครบถ้วน ได้แก่ ลงทะเบียนรับเรื่อง/ใบจองจากระบบ ........................</div><div class="document-sign">ลงชื่อ <span class="officer-signature"><img v-if="officerSignatureUrl" :src="officerSignatureUrl" alt="ลายเซ็นเจ้าหน้าที่" /><span v-else>..................................................</span></span> เจ้าหน้าที่<br>(..................................................) .......... / .......... / ..........</div><div class="document-divider"></div><strong>2.2 งานการคลัง/ฝ่ายการคลัง</strong><div class="document-indent">□ ไม่มีหนี้ค้างชำระ<br>□ มีหนี้ค้าง จำนวน ........................ บาท</div><div class="document-sign">ลงชื่อ .................................................. เจ้าหน้าที่<br>(..................................................) .......... / .......... / ..........</div></section>
            <section class="document-section"><strong>3. ผู้มีอำนาจอนุมัติ</strong><div class="document-indent">□ อนุมัติ &nbsp;&nbsp; □ ไม่อนุมัติ &nbsp;&nbsp; □ อื่น ๆ ................................................<br>□ สำนักงาน เลขระบบ วันที่ ............................</div><div class="document-sign">ลงชื่อ .................................................. ผู้มีอำนาจอนุมัติ<br>(..................................................) .......... / .......... / ..........</div></section>
            <section class="document-section"><strong>4. งานการคลัง/ผู้จัดการตลาด</strong><div class="document-indent">□ จัดเก็บเงิน จำนวน ........................ บาท<br>ตามใบเสร็จรับเงินเล่มที่ ............... เลขที่ ............... วันที่ ............... ที่แบบ ...............<br>□ ส่งเจ้าของเรื่อง วันที่ ............................</div><div class="document-sign">ลงชื่อ .................................................. เจ้าหน้าที่<br>(..................................................) .......... / .......... / ..........</div></section>
          </article>
        </div>
        <footer class="document-footer"><AppButton v-if="isTenant && !savedSignature" variant="primary" @click="openSigning">ลงนามคำร้อง</AppButton><span v-else-if="savedSignature" class="signed-note">ลงนามแล้ว</span><AppButton variant="secondary" @click="$emit('close')">ปิด</AppButton></footer>
      </section>
    </div>
    <div v-if="signing" class="sign-overlay" @click.self="signing = false"><section class="sign-dialog" role="dialog" aria-modal="true" aria-labelledby="sign-title"><h2 id="sign-title">ลงนามคำร้อง</h2><p>กรอกชื่อผู้ยื่นเพื่อบันทึกการลงนามในคำร้อง</p><label for="signature-name">ชื่อผู้ลงนาม</label><AppInput id="signature-name" v-model="signatureName" placeholder="ชื่อและนามสกุล" /><div class="sign-actions"><AppButton variant="secondary" :disabled="signingBusy" @click="signing = false">ยกเลิก</AppButton><AppButton variant="primary" :loading="signingBusy" :disabled="!signatureName.trim()" @click="submitSignature">ยืนยันลงนาม</AppButton></div></section></div>
  </Teleport>
</template>
<script setup lang="ts">
const props = defineProps<{ application: any | null }>()
const emit = defineEmits<{ close: []; signed: [value: { tenant_signature: string; tenant_signed_at: string }] }>()
const { isTenant } = useAuth()
const { post, apiFetch } = useApi()
const { showConfirm, showError } = useAppAlert()
const zoom = ref(100)
const officerSignatureUrl = ref('')
let signatureRequest = 0
const loadOfficerSignature = async (app: any) => {
  const requestId = ++signatureRequest
  if (officerSignatureUrl.value) URL.revokeObjectURL(officerSignatureUrl.value)
  officerSignatureUrl.value = ''
  if (!app?.id || !app.has_officer_signature) return
  try {
    const file = await apiFetch<Blob>(`/applications/${app.id}/officer-signature`, { responseType: 'blob' })
    if (requestId === signatureRequest) officerSignatureUrl.value = URL.createObjectURL(file)
  } catch { /* The signature is only shown when the authorized request succeeds. */ }
}
onBeforeUnmount(() => { signatureRequest++; if (officerSignatureUrl.value) URL.revokeObjectURL(officerSignatureUrl.value) })
const signing = ref(false)
const signingBusy = ref(false)
const signatureName = ref('')
const savedSignature = ref('')
watch(() => props.application, (app) => { savedSignature.value = app?.tenant_signature || ''; signing.value = false; zoom.value = 100; if (import.meta.client) loadOfficerSignature(app) }, { immediate: true })
const submittedDate = computed(() => {
  const raw = props.application?.created_at
  const date = raw ? new Date(raw) : null
  if (!date || Number.isNaN(date.getTime())) return { day: '', month: '', year: '', time: '' }
  const parts = new Intl.DateTimeFormat('th-TH-u-ca-buddhist-nu-latn', { timeZone: 'Asia/Bangkok', day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit', hourCycle: 'h23' }).formatToParts(date)
  const part = (type: string) => parts.find(item => item.type === type)?.value || ''
  return { day: part('day'), month: part('month'), year: part('year'), time: part('hour') + ':' + part('minute') + ' น.' }
})
const openSigning = () => { signatureName.value = props.application?.applicant_name || ''; signing.value = true }
const submitSignature = async () => {
  const name = signatureName.value.trim()
  if (!name || !props.application?.id || signingBusy.value) return
  const confirmed = await showConfirm('ยืนยันการลงนามคำร้องในชื่อ ' + name + ' หรือไม่?', { title: 'ยืนยันลงนาม', type: 'info', confirmText: 'ยืนยันลงนาม', cancelText: 'ยกเลิก' })
  if (!confirmed) return
  signingBusy.value = true
  try {
    const result: any = await post(`/applications/${props.application.id}/sign`, { name })
    const data = result?.data || result
    savedSignature.value = data.tenant_signature || name
    emit('signed', { tenant_signature: savedSignature.value, tenant_signed_at: data.tenant_signed_at })
    signing.value = false
  } catch (err: any) { await showError(err?.data?.message || err?.data?.detail || 'ไม่สามารถบันทึกการลงนามได้') }
  finally { signingBusy.value = false }
}

const applicantTitle = computed(() => String(props.application?.applicant_name || '').trim().match(/^(นางสาว|นาย|นาง)(?=\s|[^\s])/u)?.[0] || '')
const applicantName = computed(() => {
  const name = String(props.application?.applicant_name || '').trim()
  return applicantTitle.value ? name.slice(applicantTitle.value.length).trim() : name
})
const checked = (type: string) => props.application?.proposed_rent_type === type

</script>
<style scoped>
.document-overlay { position: fixed; inset: 0; z-index: 2000; display: grid; place-items: center; padding: 16px; background: rgba(15,23,42,.55); }
.document-modal { width: min(900px, 100%); max-height: 94vh; display: flex; flex-direction: column; border-radius: 16px; overflow: hidden; background: #fff; box-shadow: 0 20px 35px rgba(15,23,42,.2); }
.document-header, .document-footer { display: flex; align-items: center; padding: 14px 20px; gap: 12px; }
.document-header { justify-content: space-between; border-bottom: 1px solid #e2e8f0; } .document-header h2 { margin: 0; font-size: 1.125rem; color: #1e293b; }
.document-close { display: grid; place-items: center; width: 34px; height: 34px; border: 0; border-radius: 8px; background: #f1f5f9; color: #475569; cursor: pointer; }
.document-tools { display: flex; justify-content: center; padding: 8px 16px; border-bottom: 1px solid #e2e8f0; }
.zoom-controls { display: flex; align-items: center; gap: 8px; color: #475569; font-size: .875rem; }
.zoom-controls span { min-width: 48px; text-align: center; }
.document-scroll { flex: 1; overflow: auto; padding: 24px; background: #f1f5f9; }
.document-paper { max-width: 740px; margin: auto; padding: 32px; background: white; color: #0f172a; box-shadow: 0 1px 4px rgba(15,23,42,.15); font-size: .8125rem; line-height: 1.65; }
.document-paper h3 { text-align: center; margin: 4px 0 12px; font-size: 1rem; } .document-paper p { margin: 10px 0; }
.document-number { text-align: right; } .document-date { width: min(420px, 100%); margin-left: auto; } .document-date > div { display: flex; align-items: baseline; gap: 5px; } .document-date .place, .document-date .date-time, .document-date .date-month { flex: 1; } .document-date .date-day { min-width: 35px; } .document-date .date-year { min-width: 55px; }
.document-section { border: 1px solid #94a3b8; padding: 12px; margin: 10px 0; } .document-row { overflow-wrap: anywhere; } .document-value { font-weight: 600; } .document-indent { padding-left: 18px; }
.fill-line { display: inline-block; min-height: 1.2em; min-width: 80px; padding: 0 3px; border-bottom: 1px dotted #475569; text-align: center; font-weight: 600; } .fill-line.signature { min-width: 170px; }
.form-check { display: inline-block; position: relative; width: 13px; height: 13px; border: 1px solid #334155; vertical-align: middle; }
.form-check.marked::after { content: ''; position: absolute; width: 5px; height: 8px; left: 3px; top: 0; border: solid #111; border-width: 0 2px 2px 0; transform: rotate(45deg); }
.document-sign { text-align: right; margin-top: 8px; } .tenant-name { display: inline-block; width: 250px; text-align: center; } .officer-signature { display: inline-flex; width: 170px; min-height: 24px; justify-content: center; align-items: flex-end; vertical-align: bottom; } .officer-signature img { max-width: 165px; max-height: 55px; object-fit: contain; } .document-divider { border-top: 1px solid #cbd5e1; margin: 12px -12px; }
.signed-note { margin-right: auto; color: #15803d; font-size: .875rem; font-weight: 600; }
.sign-overlay { position: fixed; inset: 0; z-index: 2100; display: grid; place-items: center; padding: 16px; background: rgba(15,23,42,.45); } .sign-dialog { width: min(450px,100%); background: white; padding: 24px; border-radius: 14px; } .sign-dialog h2 { margin: 0; font-size: 1.125rem; } .sign-dialog p { color: #64748b; font-size: .875rem; } .sign-dialog label { display: block; margin: 12px 0 6px; font-size: .875rem; font-weight: 600; } .sign-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 20px; }
.document-footer { justify-content: flex-end; border-top: 1px solid #e2e8f0; }
@media (max-width: 600px) { .document-scroll { padding: 12px; } .document-paper { padding: 18px; } }
</style>
