<template>
  <Teleport to="body">
    <div v-if="application" class="document-overlay" @click.self="$emit('close')">
      <section class="document-modal" role="dialog" aria-modal="true" aria-labelledby="document-title">
        <AppModalHeader :title="application.has_signed_document ? 'เอกสารลงนาม' : 'ตัวอย่างคำร้อง'" :subtitle="`คำร้องเลขที่ ${application.application_number || application.document_number || '-'}`" icon="document" @close="$emit('close')" />
        <template v-if="application.has_signed_document">
          <div class="signed-file-body">
            <p v-if="signedFileLoading" class="signed-file-message">กำลังโหลดเอกสาร...</p>
            <div v-else-if="signedFileError" class="signed-file-message" role="alert">
              <p>{{ signedFileError }}</p>
              <AppButton variant="outline" size="sm" @click="loadSignedFile(application)">ลองอีกครั้ง</AppButton>
            </div>
            <iframe v-else-if="signedFileType === 'application/pdf' && signedFileUrl" :src="signedFileUrl" title="เอกสารลงนาม" />
            <div v-else-if="signedFileUrl" class="signed-image-wrap"><img :src="signedFileUrl" alt="เอกสารลงนาม" /></div>
          </div>
          <footer class="document-footer">
            <AppButton variant="primary" size="md" :disabled="!signedFileUrl || signedFileLoading" @click="downloadSignedFile">ดาวน์โหลดเอกสาร</AppButton>
            <AppButton variant="secondary" size="md" @click="$emit('close')">ปิด</AppButton>
          </footer>
        </template>
        <template v-else>
        <div class="document-tools"><div class="zoom-controls"><AppButton size="sm" variant="outline" :disabled="zoom <= 70" @click="zoom = Math.max(70, zoom - 10)">−</AppButton><span>{{ zoom }}%</span><AppButton size="sm" variant="outline" :disabled="zoom >= 180" @click="zoom = Math.min(180, zoom + 10)">+</AppButton><AppButton size="sm" variant="ghost" @click="zoom = 100">ขนาดปกติ</AppButton></div></div>

        <div class="document-scroll">
          <article ref="documentPaper" class="document-paper" :style="{ zoom: zoom / 100 }">
<template v-if="application.application_type !== 'general'">
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
            <section class="document-section"><strong>2. ส่วนของเจ้าหน้าที่</strong><br><strong>2.1 งานผลประโยชน์และเร่งรัดหนี้สิน/ผู้จัดการตลาด</strong>
              <div class="document-indent">
                <div><span class="form-check" :class="{ marked: officerReview?.stall_status === 'none' }"></span> ไม่มีแผงค้า</div>
                <div><span class="form-check" :class="{ marked: officerReview?.stall_status === 'has' }"></span> มีแผงค้า จำนวน {{ officerReview?.stall_status === 'has' ? officerReview.stall_count : '............' }} แผง เลขที่ {{ officerReview?.stall_numbers || '........................' }} โครงการ {{ officerReview?.project || '........................' }} ประเภท {{ officerReview?.trade_type || '........................' }}</div>
                <div><span class="form-check" :class="{ marked: officerReview?.other_contract }"></span> สัญญาอื่น ๆ กับสำนักงานตลาดกรุงเทพมหานคร</div>
                <div><span class="form-check" :class="{ marked: officerReview?.other_checked }"></span> อื่น ๆ เช่น มีเรื่องร้องเรียน/ถูกเพิกถอนสิทธิ/ถูกคัดสิทธิ/ถูกบอกเลิกสัญญา {{ officerReview?.other_details || '........................................' }}</div>
                <div><span class="form-check" :class="{ marked: officerReview?.documents_complete }"></span> คุณสมบัติและเอกสารถูกต้องครบถ้วน ได้แก่ ลงทะเบียนเรียบร้อย/ใบจองจากระบบ {{ officerReview?.booking_reference || '........................' }}</div>
              </div><div class="document-sign">ลงชื่อ <span class="officer-signature"><img v-if="reviewerSignatures[0].url" :src="reviewerSignatures[0].url" alt="ลายเซ็นเจ้าหน้าที่" /><span v-else>..................................................</span></span> เจ้าหน้าที่<br>(..................................................) .......... / .......... / ..........</div><div class="document-divider"></div><strong>2.2 งานการคลัง/ฝ่ายการคลัง</strong>
              <div class="document-indent">
                <div><span class="form-check" :class="{ marked: marketReview?.debt_status === 'none' }"></span> ไม่มีหนี้ค้างชำระ</div>
                <div><span class="form-check" :class="{ marked: marketReview?.debt_status === 'has' }"></span> มีหนี้ค้าง จำนวน {{ marketReview?.debt_status === 'has' ? Number(marketReview.debt_amount).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) : '........................' }} บาท</div>
              </div><div class="document-sign">ลงชื่อ <span class="officer-signature"><img v-if="reviewerSignatures[1].url" :src="reviewerSignatures[1].url" alt="ลายเซ็นฝ่ายการคลัง" /><span v-else>..................................................</span></span> เจ้าหน้าที่ฝ่ายการคลัง<br>(..................................................) .......... / .......... / ..........</div></section>
            <section class="document-section"><strong>3. ผู้มีอำนาจอนุมัติ</strong>
              <div class="document-indent">
                <div><span class="form-check" :class="{ marked: bureauReview?.decision === 'approve' }"></span> อนุมัติ &nbsp;&nbsp; <span class="form-check" :class="{ marked: bureauReview?.decision === 'reject' }"></span> ไม่อนุมัติ</div>
                <div><span class="form-check" :class="{ marked: bureauReview?.office_checked }"></span> สำนักงาน เลขระบบ {{ bureauReview?.system_number || '........................' }} วันที่ {{ formatReviewDate(bureauReview?.document_date) }}</div>
              </div><div class="document-sign">ลงชื่อ <span class="officer-signature"><img v-if="reviewerSignatures[3].url" :src="reviewerSignatures[3].url" alt="ลายเซ็น ผอ.สำนักตลาด" /><span v-else>..................................................</span></span> ผอ.สำนักตลาด<br>(..................................................) .......... / .......... / ..........</div></section>
            <section class="document-section"><strong>4. งานการคลัง/ผู้จัดการตลาด</strong><div class="document-indent">□ จัดเก็บเงิน จำนวน ........................ บาท<br>ตามใบเสร็จรับเงินเล่มที่ ............... เลขที่ ............... วันที่ ............... ที่แบบ ...............<br>□ ส่งเจ้าของเรื่อง วันที่ ............................</div><div class="document-sign">ลงชื่อ .................................................. เจ้าหน้าที่<br>(..................................................) .......... / .......... / ..........</div></section>
            <div class="document-sign document-bottom-sign">ลงชื่อ <span class="officer-signature"><img v-if="reviewerSignatures[2].url" :src="reviewerSignatures[2].url" alt="ลายเซ็น ผอ.ตลาด" /><span v-else>..................................................</span></span> ผอ.ตลาด<br>(..................................................) .......... / .......... / ..........</div>
            </template>
            <template v-else>
              <div class="document-number">เลขที่คำร้อง {{ application.application_number }}</div>
              <h3>คำร้องทั่วไป</h3>
              <p><strong>เรื่อง</strong> {{ application.subject }}</p>
              <p><strong>ผู้ยื่นคำร้อง</strong> {{ application.applicant_name }} อายุ {{ application.applicant_age || '-' }} ปี</p>
              <p>เลขบัตรประชาชน {{ application.applicant_id_card || '-' }} โทร {{ application.applicant_phone || '-' }}</p>
              <p>ตลาด {{ application.market_name || '-' }} แผง {{ application.stall_number || '-' }}</p>
              <p><strong>รายละเอียด</strong> {{ application.details || '-' }}</p>
              <div class="document-sign">ลงชื่อ .................................................. ผู้ยื่นคำร้อง</div>
            </template>
            <div v-if="application.application_type === 'general'" class="document-section">
              <div>ลายเซ็นผู้ตรวจสอบและผู้อนุมัติ</div>
              <div v-for="entry in reviewerSignatures" :key="entry.role" class="document-sign">
                <span class="officer-signature"><img v-if="entry.url" :src="entry.url" :alt="'ลายเซ็น' + entry.label" /><span v-else>..................................................</span></span> {{ entry.label }}
              </div>
            </div>
          </article>
        </div>
        <footer class="document-footer"><AppButton v-if="!isTenant" variant="primary" size="md" :disabled="signaturesLoading || downloadingPdf" :loading="downloadingPdf" @click="downloadPdf">ดาวน์โหลด PDF</AppButton><AppButton v-if="isTenant && !savedSignature" variant="primary" size="md" @click="openSigning">ลงนามคำร้อง</AppButton><span v-else-if="savedSignature" class="signed-note">ลงนามแล้ว</span><AppButton variant="secondary" size="md" @click="$emit('close')">ปิด</AppButton></footer>
        </template>
      </section>
    </div>
    <div v-if="signing" class="sign-overlay" @click.self="signing = false">
      <section class="sign-dialog" role="dialog" aria-modal="true" aria-labelledby="sign-title">
        <AppModalHeader title="ลงนามคำร้อง" subtitle="กรอกชื่อผู้ยื่นเพื่อบันทึกการลงนามในคำร้อง" icon="edit" @close="signing = false" />
        <div class="sign-body"><label for="signature-name">ชื่อผู้ลงนาม <span class="req-star" aria-hidden="true">*</span></label><AppInput id="signature-name" v-model="signatureName" placeholder="ชื่อและนามสกุล" required :error="signatureError" /></div>
        <div class="sign-actions"><AppButton variant="secondary" size="md" :disabled="signingBusy" @click="signing = false">ยกเลิก</AppButton><AppButton variant="primary" size="md" :loading="signingBusy" @click="submitSignature">ยืนยันลงนาม</AppButton></div>
      </section>
    </div>
  </Teleport>
</template>
<script setup lang="ts">
const props = defineProps<{ application: any | null }>()
const emit = defineEmits<{ close: []; signed: [value: { tenant_signature: string; tenant_signed_at: string }] }>()
const { isTenant } = useAuth()
const { post, apiFetch } = useApi()
const { showConfirm, showError } = useAppAlert()
const zoom = ref(100)
const documentPaper = ref<HTMLElement | null>(null)
const downloadingPdf = ref(false)
const signedFileUrl = ref('')
const signedFileType = ref('')
const signedFileLoading = ref(false)
const signedFileError = ref('')
let signedFileRequest = 0
const clearSignedFile = () => {
  if (signedFileUrl.value) URL.revokeObjectURL(signedFileUrl.value)
  signedFileUrl.value = ''
  signedFileType.value = ''
}
const loadSignedFile = async (app: any) => {
  const request = ++signedFileRequest
  clearSignedFile()
  signedFileError.value = ''
  signedFileLoading.value = !!app?.id && !!app?.has_signed_document
  if (!signedFileLoading.value) return
  try {
    const file = await apiFetch<Blob>('/applications/' + app.id + '/signed-document', { responseType: 'blob' })
    if (request !== signedFileRequest) return
    const type = file.type.toLowerCase().split(';')[0]
    if (!['application/pdf', 'image/jpeg', 'image/png'].includes(type)) throw new Error('unsupported file')
    signedFileType.value = type
    signedFileUrl.value = URL.createObjectURL(file)
  } catch {
    if (request === signedFileRequest) signedFileError.value = 'โหลดเอกสารลงนามไม่สำเร็จ'
  } finally {
    if (request === signedFileRequest) signedFileLoading.value = false
  }
}
const downloadSignedFile = () => {
  if (!signedFileUrl.value) return
  const extension = signedFileType.value === 'application/pdf' ? 'pdf' : signedFileType.value === 'image/png' ? 'png' : 'jpg'
  const link = document.createElement('a')
  link.href = signedFileUrl.value
  link.download = 'signed-application-' + String(props.application?.id || 'document') + '.' + extension
  link.click()
}
watch(() => props.application, app => { if (import.meta.client) void loadSignedFile(app) }, { immediate: true })
onBeforeUnmount(() => { signedFileRequest++; clearSignedFile() })
const officerReview = computed(() => props.application?.officer_review || null)
const marketReview = computed(() => props.application?.market_director_review || null)
const bureauReview = computed(() => props.application?.bureau_director_review || null)
const formatReviewDate = (value?: string) => {
  if (!value) return '........................'
  const date = new Date(value + 'T00:00:00Z')
  return Number.isFinite(date.getTime()) ? new Intl.DateTimeFormat('th-TH-u-ca-buddhist-nu-latn', { day: 'numeric', month: 'short', year: 'numeric', timeZone: 'UTC' }).format(date) : value
}
const reviewerSignatures = ref([{ role: 'officer', label: 'เจ้าหน้าที่', url: '' }, { role: 'finance', label: 'ฝ่ายการคลัง', url: '' }, { role: 'market', label: 'ผอ.ตลาด', url: '' }, { role: 'bureau', label: 'ผอ.สำนักตลาด', url: '' }])
const signaturesLoading = ref(false)
const signatureLoadError = ref('')
let reviewRequest = 0
const loadReviewSignatures = async (app: any) => {
  const request = ++reviewRequest
  reviewerSignatures.value.forEach(item => { if (item.url) URL.revokeObjectURL(item.url); item.url = '' })
  signatureLoadError.value = ''
  signaturesLoading.value = !!app?.id
  if (!app?.id || app.has_signed_document) { signaturesLoading.value = false; return }
  const state = String(app.status || app.application_state || '')
  const flags = [
    app.has_officer_signature === true,
    app.has_finance_signature === true,
    app.has_market_director_signature === true,
    app.has_bureau_director_signature === true
  ]
  const results = await Promise.allSettled(reviewerSignatures.value.map(async (item, index) => {
    if (!flags[index]) return
    try {
      const blob = await apiFetch<Blob>('/applications/' + app.id + '/' + item.role + '-signature', { responseType: 'blob' })
      if (request !== reviewRequest) return
      item.url = URL.createObjectURL(blob)
      await new Promise<void>((resolve, reject) => {
        const image = new Image()
        image.onload = () => resolve()
        image.onerror = () => reject(new Error('โหลดรูปไม่สำเร็จ'))
        image.src = item.url
      })
    } catch (err: any) {
      const statusCode = err?.statusCode || err?.status || err?.response?.status
      if (statusCode === 404 && !flags[index]) return
      if (statusCode === 404) return
      throw new Error('โหลดลายเซ็น' + item.label + 'ไม่สำเร็จ กรุณาลองอีกครั้ง')
    }
  }))
  if (request !== reviewRequest) return
  signaturesLoading.value = false
  const failures = results.filter((result): result is PromiseRejectedResult => result.status === 'rejected')
  if (failures.length) signatureLoadError.value = failures.map(result => result.reason?.message || 'โหลดลายเซ็นไม่สำเร็จ').join(' · ')
}
const downloadPdf = async () => {
  if (!documentPaper.value || signaturesLoading.value || downloadingPdf.value) return
  downloadingPdf.value = true
  try {
    await document.fonts.ready
    const html2pdf = (await import('html2pdf.js')).default
    const number = String(props.application?.application_number || props.application?.application_no || props.application?.id || 'document').replace(/[^a-zA-Z0-9_-]/g, '-')
    await html2pdf().set({
      margin: [8, 8, 8, 8],
      filename: 'application-' + number + '.pdf',
      image: { type: 'jpeg', quality: 0.98 },
      html2canvas: {
        scale: 2,
        backgroundColor: '#ffffff',
        onclone: (cloned: Document) => {
          const paper = cloned.querySelector<HTMLElement>('.html2pdf__container .document-paper')
          if (paper) { paper.style.zoom = '1'; paper.style.boxShadow = 'none' }
        }
      },
      jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' },
      pagebreak: { mode: ['css', 'legacy'], avoid: ['.document-sign'] }
    }).from(documentPaper.value).save()
  } catch {
    await showError('ดาวน์โหลด PDF ไม่สำเร็จ กรุณาลองอีกครั้ง')
  } finally {
    downloadingPdf.value = false
  }
}
watch(() => props.application, app => { if (import.meta.client) void loadReviewSignatures(app) }, { immediate: true })
onBeforeUnmount(() => { reviewRequest++; reviewerSignatures.value.forEach(item => { if (item.url) URL.revokeObjectURL(item.url) }) })
const signing = ref(false)
useBodyScrollLock(computed(() => !!props.application || signing.value))
const signingBusy = ref(false)
const signatureName = ref('')
const signatureError = ref('')
const savedSignature = ref('')
watch(() => props.application, (app) => { savedSignature.value = app?.tenant_signature || ''; signing.value = false; zoom.value = 100 }, { immediate: true })
const submittedDate = computed(() => {
  const raw = props.application?.created_at
  const date = raw ? new Date(raw) : null
  if (!date || Number.isNaN(date.getTime())) return { day: '', month: '', year: '', time: '' }
  const parts = new Intl.DateTimeFormat('th-TH-u-ca-buddhist-nu-latn', { timeZone: 'Asia/Bangkok', day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit', hourCycle: 'h23' }).formatToParts(date)
  const part = (type: string) => parts.find(item => item.type === type)?.value || ''
  return { day: part('day'), month: part('month'), year: part('year'), time: part('hour') + ':' + part('minute') + ' น.' }
})
const openSigning = () => { signatureName.value = props.application?.applicant_name || ''; signatureError.value = ''; signing.value = true }
watch(signatureName, value => { if (value.trim()) signatureError.value = '' })
const submitSignature = async () => {
  const name = signatureName.value.trim()
  signatureError.value = name ? '' : 'กรุณากรอกชื่อผู้ลงนาม'
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
.document-modal { width: min(900px, 100%); max-height: 94vh; display: flex; flex-direction: column; border-radius: 16px; overflow: hidden; background: var(--bg-ffffff); box-shadow: 0 20px 35px rgba(15,23,42,.2); }
.document-header, .document-footer { display: flex; align-items: center; padding: 14px 20px; gap: 12px; }
.document-header { justify-content: space-between; border-bottom: 1px solid var(--bd-e2e8f0); } .document-header h2 { margin: 0; font-size: var(--font-size-lg); color: var(--tx-1e293b); }
.document-close { display: grid; place-items: center; width: 34px; height: 34px; border: 0; border-radius: 8px; background: var(--bg-f1f5f9); color: var(--tx-475569); cursor: pointer; }
.signature-load-error { margin: 12px 20px; color: #dc2626; font-size: var(--font-size-sm); }
.signature-load-error button { color: #2563eb; text-decoration: underline; }
.document-tools { display: flex; justify-content: center; padding: 8px 16px; border-bottom: 1px solid var(--bd-e2e8f0); }
.zoom-controls { display: flex; align-items: center; gap: 8px; color: var(--tx-475569); font-size: var(--font-size-sm); }
.zoom-controls span { min-width: 48px; text-align: center; }
.signed-file-body { flex: 1 1 auto; min-height: 0; overflow: auto; background: #f1f5f9; }
.signed-file-body iframe { display: block; width: 100%; height: min(70vh, 800px); border: 0; }
.signed-file-message { padding: 24px; color: #64748b; font-size: var(--font-size-base); text-align: center; }
.signed-image-wrap { display: flex; justify-content: center; padding: 20px; }
.signed-image-wrap img { max-width: 100%; height: auto; object-fit: contain; }
.document-scroll { flex: 1; overflow: auto; padding: 24px; background: var(--bg-f1f5f9); }
.document-paper { max-width: 740px; margin: auto; padding: 32px; background: var(--bg-ffffff); color: var(--tx-0f172a); box-shadow: 0 1px 4px rgba(15,23,42,.15); font-size: var(--font-size-xs); line-height: 1.65; }
.document-paper h3 { text-align: center; margin: 4px 0 12px; font-size: var(--font-size-md); } .document-paper p { margin: 10px 0; }
.document-number { text-align: right; } .document-date { width: min(420px, 100%); margin-left: auto; } .document-date > div { display: flex; align-items: baseline; gap: 5px; } .document-date .place, .document-date .date-time, .document-date .date-month { flex: 1; } .document-date .date-day { min-width: 35px; } .document-date .date-year { min-width: 55px; }
.document-section { border: 1px solid var(--bd-94a3b8); padding: 12px; margin: 10px 0; } .document-row { overflow-wrap: anywhere; } .document-value { font-weight: 600; } .document-indent { padding-left: 18px; }
.fill-line { display: inline-block; min-height: 1.2em; min-width: 80px; padding: 0 3px; border-bottom: 1px dotted #475569; text-align: center; font-weight: 600; } .fill-line.signature { min-width: 170px; }
.form-check { display: inline-block; position: relative; width: 13px; height: 13px; border: 1px solid #334155; vertical-align: middle; }
.form-check.marked::after { content: ''; position: absolute; width: 5px; height: 8px; left: 3px; top: 0; border: solid #111; border-width: 0 2px 2px 0; transform: rotate(45deg); }
.document-sign { text-align: right; margin-top: 8px; } .document-bottom-sign { margin-top: 12px; } .tenant-name { display: inline-block; width: 250px; text-align: center; } .officer-signature { display: inline-flex; width: 170px; min-height: 24px; justify-content: center; align-items: flex-end; vertical-align: bottom; } .officer-signature img { max-width: 165px; max-height: 55px; object-fit: contain; } .document-divider { border-top: 1px solid var(--bd-cbd5e1); margin: 12px -12px; }
.signed-note { margin-right: auto; color: var(--tx-15803d); font-size: var(--font-size-sm); font-weight: 600; }
.sign-overlay { position: fixed; inset: 0; z-index: 2100; display: grid; place-items: center; padding: 16px; background: rgba(15,23,42,.45); } .sign-dialog { width: min(450px,100%); max-height:min(90vh,90dvh); display:flex; flex-direction:column; overflow:hidden; background: var(--bg-ffffff); border-radius: 20px; } .sign-body{flex:1 1 auto;min-height:0;overflow-y:auto;overscroll-behavior:contain;padding:24px 28px}.sign-dialog label { display: block; margin: 0 0 6px; font-size: var(--font-size-sm); font-weight: 600; } .sign-actions { flex:0 0 auto;display: flex; justify-content: flex-end; gap: 10px; padding:16px 28px; border-top:1px solid var(--bd-e2e8f0); background:var(--bg-f8fafc); }
.document-footer { justify-content: flex-end; border-top: 1px solid var(--bd-e2e8f0); }
@media (max-width: 600px) { .signed-file-body { flex: 1 1 auto; min-height: 0; overflow: auto; background: #f1f5f9; }
.signed-file-body iframe { display: block; width: 100%; height: min(70vh, 800px); border: 0; }
.signed-file-message { padding: 24px; color: #64748b; font-size: var(--font-size-base); text-align: center; }
.signed-image-wrap { display: flex; justify-content: center; padding: 20px; }
.signed-image-wrap img { max-width: 100%; height: auto; object-fit: contain; }
.document-scroll { padding: 12px; } .document-paper { padding: 18px; } }
@media print {
  :global(body > *:not(.document-overlay)) { display: none !important; }
  .document-overlay { position: static; display: block; padding: 0; background: white; }
  .document-modal { display: block; width: 100%; max-height: none; overflow: visible; box-shadow: none; }
  .document-tools, .document-footer, :deep(.app-modal-header), .sign-overlay { display: none !important; }
  .signed-file-body { flex: 1 1 auto; min-height: 0; overflow: auto; background: #f1f5f9; }
.signed-file-body iframe { display: block; width: 100%; height: min(70vh, 800px); border: 0; }
.signed-file-message { padding: 24px; color: #64748b; font-size: var(--font-size-base); text-align: center; }
.signed-image-wrap { display: flex; justify-content: center; padding: 20px; }
.signed-image-wrap img { max-width: 100%; height: auto; object-fit: contain; }
.document-scroll { overflow: visible; padding: 0; background: white; }
  .document-paper { zoom: 1 !important; max-width: none; box-shadow: none; }
}
</style>
