<template>
  <div v-if="show" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content" style="max-width: 660px;">
      <div class="modal-header">
        <div style="display: flex; align-items: center; gap: 10px;">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--primary);">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="16" y1="13" x2="8" y2="13"></line>
            <line x1="16" y1="17" x2="8" y2="17"></line>
          </svg>
          <div>
            <h3 style="font-size: 1.15rem; font-weight: 700;">ยื่นคำร้องขอเช่าแผงค้า (Lease Application)</h3>
            <div style="font-size: 0.8rem; color: var(--primary);">
              แผงเลขที่: {{ targetStallNumber }} (อ้างอิงการจอง: {{ booking?.booking_number || '-' }})
            </div>
          </div>
        </div>
        <button class="btn-logout" @click="$emit('close')" aria-label="ปิด">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
        </button>
      </div>

      <!-- Success State -->
      <div v-if="isSuccess" class="modal-body" style="text-align: center; padding: 40px 20px;">
        <div style="margin-bottom: 16px; color: var(--status-avail-text); display: flex; justify-content: center;">
          <div style="width: 56px; height: 56px; border-radius: 50%; background: rgba(16, 185, 129, 0.15); display: flex; align-items: center; justify-content: center;">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
          </div>
        </div>
        <h4 style="font-size: 1.25rem; font-weight: 700; color: var(--status-avail-text); margin-bottom: 8px;">
          ยื่นคำร้องขอเช่าแผงเรียบร้อยแล้ว!
        </h4>
        <div style="font-size: 0.9rem; color: var(--text-muted); line-height: 1.6; max-width: 460px; margin: 0 auto 24px;">
          เลขที่คำร้อง: <strong style="color: var(--primary); font-size: 1.05rem;">{{ successAppNumber }}</strong><br/>
          ขั้นตอนถัดไป: คำร้องจะถูกส่งให้ <strong>เจ้าหน้าที่ตรวจสอบ</strong> &rarr; <strong>ผอ.ตลาด อนุมัติ</strong> &rarr; <strong>ผอ.สำนักตลาด อนุมัติ</strong> เพื่อจัดทำสัญญาเช่าฉบับสมบูรณ์
        </div>
        <div style="display: flex; justify-content: center; gap: 12px;">
          <NuxtLink to="/applications" class="btn btn-primary" @click="$emit('close')">
            ติดตามสถานะคำร้องของฉัน
          </NuxtLink>
          <button class="btn btn-secondary" @click="$emit('close')">ปิดหน้าต่าง</button>
        </div>
      </div>

      <!-- Form State -->
      <form v-else novalidate @submit.prevent="handleSubmit">
        <div class="modal-body">
          <!-- Financial Summary calculated from remaining active stalls -->
          <div
            style="background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 10px; padding: 12px 14px; margin-bottom: 16px;"
          >
            <div style="display: flex; flex-wrap: wrap; justify-content: space-between; gap: 10px; font-size: 0.85rem; color: #1e3a8a;">
              <div>
                <span>แผงค้าที่ยื่นคำร้อง: </span>
                <strong>{{ targetStallNumber }} ({{ selectedStallRecords.length }} แผง)</strong>
              </div>
              <div style="display: flex; gap: 14px; flex-wrap: wrap;">
                <span>ค่าเช่ารวม: <strong>฿{{ calculatedMonthlyRent.toLocaleString() }}/เดือน</strong></span>
                <span>เงินมัดจำรวม: <strong>฿{{ calculatedDeposit.toLocaleString() }}</strong></span>
                <span>รวมแรกเข้า: <strong style="color: #1d4ed8;">฿{{ (calculatedMonthlyRent + calculatedDeposit).toLocaleString() }}</strong></span>
              </div>
            </div>
          </div>

          <h4 style="font-size: 0.9375rem; font-weight: 700; color: var(--primary); margin-bottom: 14px; display: flex; align-items: center; gap: 8px;">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
            <span>1. ข้อมูลผู้ยื่นคำร้อง (Applicant Details)</span>
          </h4>

          <!-- ลำดับที่ 1: ชื่อ-นามสกุล -> ลำดับที่ 2: อายุ -->
          <div class="form-row">
            <div class="form-group">
              <label class="form-label" for="app_name">ชื่อ-นามสกุล <b class="req-star">*</b></label>
              <AppInput
                id="app_name"
                v-model="form.applicant_name"
                type="text"
                placeholder="ระบุชื่อและนามสกุลจริง"
                :error="errors.applicant_name"
                @blur="validateField('applicant_name')"
              />
            </div>

            <div class="form-group">
              <label class="form-label" for="app_age">อายุ (ปี)</label>
              <AppInput
                id="app_age"
                v-model="form.applicant_age"
                type="number"
                min="1"
                max="120"
                placeholder="ระบุอายุ เช่น 35"
              />
            </div>
          </div>

          <!-- ลำดับที่ 3: เลขบัตรประชาชน -> ลำดับที่ 4: เบอร์โทร -->
          <div class="form-row">
            <div class="form-group">
              <label class="form-label" for="app_id">เลขประจำตัวประชาชน (13 หลัก) <b class="req-star">*</b></label>
              <AppInput
                id="app_id"
                v-model="form.applicant_id_card"
                type="text"
                placeholder="x-xxxx-xxxxx-xx-x"
                maxlength="17"
                :error="errors.applicant_id_card"
                @blur="validateField('applicant_id_card')"
              />
            </div>

            <div class="form-group">
              <label class="form-label" for="app_phone">เบอร์โทรศัพท์ติดต่อ <b class="req-star">*</b></label>
              <AppInput
                id="app_phone"
                v-model="form.applicant_phone"
                type="tel"
                placeholder="08x-xxx-xxxx"
                :error="errors.applicant_phone"
                @blur="validateField('applicant_phone')"
              />
            </div>
          </div>

          <!-- ลำดับที่ 5: อีเมล -> ที่อยู่ตามบัตรประชาชน -->
          <div class="form-row">
            <div class="form-group">
              <label class="form-label" for="app_email">อีเมล</label>
              <AppInput
                id="app_email"
                v-model="form.applicant_email"
                type="email"
                placeholder="example@domain.com"
              />
            </div>

            <div class="form-group">
              <label class="form-label" for="app_address">ที่อยู่ตามบัตรประชาชน / ที่อยู่ติดต่อ</label>
              <AppInput
                id="app_address"
                v-model="form.applicant_address"
                type="text"
                placeholder="บ้านเลขที่ ถนน แขวง/ตำบล เขต/อำเภอ จังหวัด"
              />
            </div>
          </div>

          <!-- ลำดับที่ 6: สำเนาบัตรประชาชน (Modern Document Attachment UI) -->
          <div class="form-group upload-section-container">
            <div class="upload-header-bar">
              <label class="upload-label" for="app_id_copy">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="3" y="4" width="18" height="16" rx="2"></rect>
                  <circle cx="8" cy="10" r="2"></circle>
                  <line x1="14" y1="9" x2="18" y2="9"></line>
                  <line x1="14" y1="13" x2="18" y2="13"></line>
                  <line x1="6" y1="16" x2="18" y2="16"></line>
                </svg>
                <span>สำเนาบัตรประชาชน (ID Card Copy)</span>
              </label>
              <span class="upload-badge-spec">PDF, JPG หรือ PNG (ไม่เกิน 3 MB)</span>
            </div>

            <input
              ref="idCardInputRef"
              id="app_id_copy"
              type="file"
              class="sr-only"
              accept=".pdf,.jpg,.jpeg,.png,application/pdf,image/jpeg,image/png"
              @change="onIDCardCopyChange"
            />

            <!-- การ์ดพรีวิวไฟล์เมื่อเลือกเอกสารแล้ว -->
            <div v-if="idCardCopy" class="id-card-preview-box">
              <div class="preview-thumb-wrap">
                <img v-if="idCardPreviewUrl" :src="idCardPreviewUrl" alt="ID Card Preview" class="preview-thumb-img" />
                <div v-else class="preview-pdf-badge">
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                    <polyline points="14 2 14 8 20 8"></polyline>
                    <line x1="16" y1="13" x2="8" y2="13"></line>
                    <line x1="16" y1="17" x2="8" y2="17"></line>
                  </svg>
                  <span>PDF</span>
                </div>
              </div>

              <div class="preview-details">
                <div class="preview-filename" :title="idCardCopy.name">{{ idCardCopy.name }}</div>
                <div class="preview-meta">
                  <span class="meta-size">{{ formatFileSize(idCardCopy.size) }}</span>
                  <span class="meta-check">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                    <span>พร้อมแนบไฟล์</span>
                  </span>
                </div>
              </div>

              <div class="preview-actions">
                <button type="button" class="btn-preview-change" @click="triggerFileInput" title="เปลี่ยนไฟล์">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 2v6h-6"></path>
                    <path d="M3 12a9 9 0 0 1 15-6.7L21 8"></path>
                    <path d="M3 22v-6h6"></path>
                    <path d="M21 12a9 9 0 0 1-15 6.7L3 16"></path>
                  </svg>
                  <span>เปลี่ยน</span>
                </button>
                <button type="button" class="btn-preview-delete" @click="removeIDCardFile" title="ลบไฟล์">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                </button>
              </div>
            </div>

            <!-- กล่อง Dropzone แบบคลีนเมื่อยังไม่ได้เลือกไฟล์ -->
            <div
              v-else
              class="id-card-dropzone"
              :class="{ 'is-drag-over': isDragging, 'has-error': !!idCardCopyError }"
              @click="triggerFileInput"
              @dragover.prevent="isDragging = true"
              @dragleave.prevent="isDragging = false"
              @drop.prevent="handleFileDrop"
            >
              <div class="dropzone-icon-box">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                  <polyline points="17 8 12 3 7 8"></polyline>
                  <line x1="12" y1="3" x2="12" y2="15"></line>
                </svg>
              </div>
              <div class="dropzone-prompt">
                <div class="prompt-title">
                  <span class="prompt-link">คลิกเพื่อเลือกไฟล์</span> หรือลากไฟล์มาวางที่นี่
                </div>
                <div class="prompt-desc">แนบสำเนาบัตรประจำตัวประชาชน พร้อมรับรองสำเนาถูกต้อง</div>
              </div>
              <div class="dropzone-cta-btn">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
                </svg>
                <span>เลือกไฟล์</span>
              </div>
            </div>

            <!-- ข้อความแจ้งเตือน Error -->
            <div v-if="idCardCopyError" class="upload-error-msg">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              <span>{{ idCardCopyError }}</span>
            </div>
          </div>

          <h4 style="font-size: 0.9rem; font-weight: 700; color: var(--primary); margin: 16px 0 12px;">
            2. ข้อมูลร้านค้าและสินค้าที่ขอจำหน่าย
          </h4>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label" for="app_shop">ชื่อร้านค้า (Shop Name) <b class="req-star">*</b></label>
              <AppInput id="app_shop" v-model="form.shop_name" type="text" placeholder="เช่น ชานมไข่มุก เฮงเฮง" :error="errors.shop_name" @blur="validateField('shop_name')" />
            </div>

            <div class="form-group">
              <label class="form-label" for="app_prod">ประเภทสินค้าที่จะจำหน่าย <b class="req-star">*</b></label>
              <AppInput id="app_prod" v-model="form.product_type" type="text" placeholder="เช่น เครื่องดื่ม, ของทอด, เสื้อผ้า" :error="errors.product_type" @blur="validateField('product_type')" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label" for="app_start">วันที่ประสงค์เริ่มเช่า <b class="req-star">*</b></label>
              <AppInput id="app_start" v-model="form.requested_start_date" type="date" :error="errors.requested_start_date" @blur="validateField('requested_start_date')" />
            </div>

            <div class="form-group">
              <label class="form-label">ระยะเวลาขอเช่า (เดือน) <b class="req-star">*</b></label>
              <AppSelect
                v-model="form.requested_period_months"
                :options="periodOptions"
                placeholder="เลือกระยะเวลาขอเช่า"
              />
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <AppButton type="button" variant="secondary" size="md" @click="$emit('close')">ยกเลิก</AppButton>
          <AppButton type="submit" variant="primary" size="md" :loading="submitting" :disabled="submitting">
            <span v-if="submitting">กำลังส่งคำร้อง...</span>
            <span v-else>ยื่นคำร้องขอเช่าแผง {{ targetStallNumber }}</span>
          </AppButton>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  show: boolean
  booking: any
  stall: any
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

const { post } = useApi()
const { showError } = useAppAlert()

const submitting = ref(false)
const idCardCopy = ref<File | null>(null)
const idCardCopyError = ref('')
const idCardInputRef = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
const idCardPreviewUrl = ref('')

const triggerFileInput = () => {
  idCardInputRef.value?.click()
}

const handleFileSelect = (file: File | null) => {
  idCardCopyError.value = ''
  if (idCardPreviewUrl.value) {
    URL.revokeObjectURL(idCardPreviewUrl.value)
    idCardPreviewUrl.value = ''
  }
  if (!file) {
    idCardCopy.value = null
    return
  }
  if (!['application/pdf', 'image/jpeg', 'image/png'].includes(file.type) || file.size > 3 * 1024 * 1024) {
    idCardCopy.value = null
    if (idCardInputRef.value) idCardInputRef.value.value = ''
    idCardCopyError.value = 'กรุณาเลือกไฟล์ PDF, JPG หรือ PNG ขนาดไม่เกิน 3 MB'
    return
  }
  idCardCopy.value = file
  if (file.type.startsWith('image/')) {
    idCardPreviewUrl.value = URL.createObjectURL(file)
  }
}

const onIDCardCopyChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0] || null
  handleFileSelect(file)
}

const handleFileDrop = (event: DragEvent) => {
  isDragging.value = false
  const file = event.dataTransfer?.files?.[0] || null
  handleFileSelect(file)
}

const removeIDCardFile = () => {
  if (idCardPreviewUrl.value) {
    URL.revokeObjectURL(idCardPreviewUrl.value)
    idCardPreviewUrl.value = ''
  }
  idCardCopy.value = null
  idCardCopyError.value = ''
  if (idCardInputRef.value) idCardInputRef.value.value = ''
}

const formatFileSize = (bytes: number) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

onBeforeUnmount(() => {
  if (idCardPreviewUrl.value) {
    URL.revokeObjectURL(idCardPreviewUrl.value)
  }
})

const readIDCardCopy = (file: File) => new Promise<string>((resolve, reject) => {
  const reader = new FileReader()
  reader.onerror = () => reject(new Error('ไม่สามารถอ่านไฟล์สำเนาบัตรประชาชน'))
  reader.onload = () => resolve(String(reader.result).split(',')[1] || '')
  reader.readAsDataURL(file)
})
const isSuccess = ref(false)
const successAppNumber = ref('')
const selectedSubBookingId = ref<number | string>('')

const periodOptions = [
  { value: 3, label: '3 เดือน (ทดลองขาย)' },
  { value: 6, label: '6 เดือน' },
  { value: 12, label: '1 ปี (มาตรฐาน)' },
  { value: 24, label: '2 ปี' }
]

const selectableSubBookings = computed(() => {
  const subs = props.booking?.sub_bookings
  if (Array.isArray(subs) && subs.length > 0) {
    const interestedSubs = subs.filter((s: any) => s.status === 'interested')
    if (interestedSubs.length > 0) return interestedSubs
    return subs.filter((s: any) => !['cancelled', 'rejected', 'not_interested', 'no_show', 'terminated', 'expired'].includes(s.status))
  }
  return props.booking ? [props.booking] : []
})

const subBookingSelectOptions = computed(() => {
  const totalRate = selectableSubBookings.value.reduce((sum: number, s: any) => sum + Number(s.stall_monthly_rate || 3600), 0)
  const allCodes = selectableSubBookings.value.map((s: any) => s.stall_number).join(', ')
  const items = selectableSubBookings.value.map((s: any) => ({
    value: s.id,
    label: `เฉพาะแผง ${s.stall_number}${s.stall_zone_name ? ` (${s.stall_zone_name})` : ''} — ฿${Number(s.stall_monthly_rate || 3600).toLocaleString()}/เดือน`
  }))
  if (selectableSubBookings.value.length > 1) {
    return [
      {
        value: 'all',
        label: `ทุกแผงที่เหลือในใบจองนี้ (${selectableSubBookings.value.length} แผง: ${allCodes}) — รวม ฿${totalRate.toLocaleString()}/เดือน`
      },
      ...items
    ]
  }
  return items
})

const selectedStallRecords = computed(() => {
  if (selectedSubBookingId.value === 'all' && selectableSubBookings.value.length > 0) {
    return selectableSubBookings.value
  }
  if (selectableSubBookings.value.length > 1 && selectedSubBookingId.value) {
    const found = selectableSubBookings.value.find((s: any) => String(s.id) === String(selectedSubBookingId.value))
    if (found) return [found]
  }
  if (selectableSubBookings.value.length === 1) {
    return [selectableSubBookings.value[0]]
  }
  return props.booking ? [props.booking] : []
})

const activeBookingRecord = computed(() => {
  return selectedStallRecords.value[0] || props.booking
})

const calculatedMonthlyRent = computed(() => {
  return selectedStallRecords.value.reduce((sum: number, s: any) => {
    const rate = Number(s?.stall_monthly_rate || s?.monthly_rent || 3600)
    return sum + rate
  }, 0)
})

const calculatedDeposit = computed(() => {
  return selectedStallRecords.value.reduce((sum: number, s: any) => {
    const rate = Number(s?.stall_monthly_rate || s?.monthly_rent || 3600)
    const dep = Number(s?.deposit_fee || rate)
    return sum + dep
  }, 0)
})

const targetStallNumber = computed(() => {
  if (selectedStallRecords.value.length > 1) {
    return selectedStallRecords.value.map((s: any) => s.stall_number).join(', ')
  }
  return activeBookingRecord.value?.stall_number || props.stall?.stall_number || ''
})

const targetStallId = computed(() => {
  return activeBookingRecord.value?.stall_id || activeBookingRecord.value?.stall_number || props.stall?.id || 1
})

const form = ref({
  applicant_name: '',
  applicant_age: '',
  applicant_id_card: '',
  applicant_phone: '',
  applicant_email: '',
  applicant_address: '',
  shop_name: '',
  product_type: '',
  requested_start_date: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().slice(0, 10),
  requested_period_months: 12,
  proposed_rent_type: 'monthly'
})

const errors = reactive({
  applicant_name: '',
  applicant_id_card: '',
  applicant_phone: '',
  shop_name: '',
  product_type: '',
  requested_start_date: ''
})

const resetErrors = () => {
  errors.applicant_name = ''
  errors.applicant_id_card = ''
  errors.applicant_phone = ''
  errors.shop_name = ''
  errors.product_type = ''
  errors.requested_start_date = ''
}

const validateField = (field: keyof typeof errors) => {
  if (field === 'applicant_name') {
    errors.applicant_name = form.value.applicant_name.trim() ? '' : 'กรุณาระบุชื่อ-นามสกุล'
  } else if (field === 'applicant_id_card') {
    errors.applicant_id_card = form.value.applicant_id_card.trim() ? '' : 'กรุณาระบุเลขประจำตัวประชาชน'
  } else if (field === 'applicant_phone') {
    errors.applicant_phone = form.value.applicant_phone.trim() ? '' : 'กรุณาระบุเบอร์โทรศัพท์ติดต่อ'
  } else if (field === 'shop_name') {
    errors.shop_name = form.value.shop_name.trim() ? '' : 'กรุณาระบุชื่อร้านค้า'
  } else if (field === 'product_type') {
    errors.product_type = form.value.product_type.trim() ? '' : 'กรุณาระบุประเภทสินค้าที่จะจำหน่าย'
  } else if (field === 'requested_start_date') {
    errors.requested_start_date = form.value.requested_start_date ? '' : 'กรุณาเลือกวันที่ประสงค์เริ่มเช่า'
  }
}

const validateForm = () => {
  validateField('applicant_name')
  validateField('applicant_id_card')
  validateField('applicant_phone')
  validateField('shop_name')
  validateField('product_type')
  validateField('requested_start_date')
  return !Object.values(errors).some(Boolean)
}

watch(() => props.show, (newVal) => {
  if (newVal) {
    isSuccess.value = false
    removeIDCardFile()
    form.value.applicant_age = ''
    resetErrors()
    if (selectableSubBookings.value.length > 1) {
      selectedSubBookingId.value = 'all'
    } else if (selectableSubBookings.value.length === 1) {
      selectedSubBookingId.value = selectableSubBookings.value[0].id
    } else {
      selectedSubBookingId.value = props.booking?.id || ''
    }
    const base = activeBookingRecord.value || props.booking
    if (base) {
      form.value.applicant_name = base.tenant_name || ''
      form.value.applicant_phone = base.tenant_phone || ''
      form.value.applicant_email = base.tenant_email || ''
      form.value.applicant_id_card = base.tenant_id_card || ''
      form.value.product_type = base.intended_business || ''
      form.value.shop_name = base.intended_business || ''
    }
  } else {
    removeIDCardFile()
  }
})

const handleSubmit = async () => {
  if (!validateForm()) return
  if (idCardCopyError.value) return
  const age = form.value.applicant_age === '' ? null : Number(form.value.applicant_age)
  if (age !== null && (!Number.isInteger(age) || age < 1 || age > 120)) { await showError('กรุณาระบุอายุระหว่าง 1 ถึง 120 ปี', 'ข้อมูลไม่ถูกต้อง'); return }
  submitting.value = true
  try {
    const idCardCopyBase64 = idCardCopy.value ? await readIDCardCopy(idCardCopy.value) : ''
    const periodMonths = Number(form.value.requested_period_months) || 12
    const base = activeBookingRecord.value || props.booking
    const mId = base?.market_id || props.stall?.market_id
    const mName = base?.market_name || props.stall?.market_name || 'ตลาด กทม.'

    const targetRecords = selectedStallRecords.value.length > 0 ? selectedStallRecords.value : (base ? [base] : [])
    let lastAppNo = ''

    for (const rec of targetRecords) {
      const sIdent = rec.stall_number || rec.stall_ident || props.stall?.stall_number || '-'
      const sId = rec.stall_id || props.stall?.id
      const bId = rec.id || base?.id
      const monthlyRate = Number(rec.stall_monthly_rate || rec.monthly_rent || 3600)

      const payload = {
        application_type: 'stall_lease',
        booking_id: bId,
        stall_id: sId,
        stall_ident: sIdent,
        stall_number: sIdent,
        market_id: mId,
        market_name: mName,
        applicant_name: form.value.applicant_name.trim(),
        applicant_age: age,
        id_card_copy_base64: idCardCopyBase64,
        applicant_phone: form.value.applicant_phone.trim(),
        applicant_id_card: form.value.applicant_id_card.trim(),
        applicant_email: form.value.applicant_email.trim(),
        applicant_address: form.value.applicant_address.trim() || '-',
        shop_name: form.value.shop_name.trim(),
        product_category: form.value.product_type.trim() || 'สินค้าทั่วไป',
        product_type: form.value.product_type.trim() || 'สินค้าทั่วไป',
        requested_start_date: form.value.requested_start_date,
        lease_term_months: periodMonths,
        requested_period_months: periodMonths,
        monthly_rate: monthlyRate,
        proposed_rent_type: form.value.proposed_rent_type || 'monthly'
      }

      const res: any = await post('/applications', payload)
      const data = res?.data || res
      if (data?.application_no || data?.application_number) {
        lastAppNo = data.application_no || data.application_number
      }
    }

    successAppNumber.value = lastAppNo || 'APP-COMPLETED'
    isSuccess.value = true
    emit('success')
  } catch (err: any) {
    await showError(err?.data?.detail || err?.data?.error?.message || err?.data?.message || err?.message || 'เกิดข้อผิดพลาดในการยื่นคำร้อง', 'ยื่นคำร้องไม่สำเร็จ')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

.upload-section-container {
  margin-top: 6px;
  margin-bottom: 8px;
}

.upload-header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  flex-wrap: wrap;
  gap: 6px;
}

.upload-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.upload-badge-spec {
  font-size: 0.75rem;
  color: #64748b;
  background: #f1f5f9;
  padding: 3px 8px;
  border-radius: 6px;
  font-weight: 500;
}

.id-card-dropzone {
  display: flex;
  align-items: center;
  gap: 14px;
  background: #f8fafc;
  border: 1.5px dashed #cbd5e1;
  border-radius: 12px;
  padding: 14px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.id-card-dropzone:hover,
.id-card-dropzone.is-drag-over {
  background: #eff6ff;
  border-color: #2563eb;
  transform: translateY(-1px);
}

.id-card-dropzone.has-error {
  background: #fef2f2;
  border-color: #f87171;
}

.dropzone-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.id-card-dropzone:hover .dropzone-icon-box {
  background: #2563eb;
  color: #ffffff;
}

.dropzone-prompt {
  flex: 1;
  min-width: 0;
}

.prompt-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 2px;
}

.prompt-link {
  color: #2563eb;
  text-decoration: underline;
  text-underline-offset: 2px;
}

.prompt-desc {
  font-size: 0.8125rem;
  color: #64748b;
  line-height: 1.4;
}

.dropzone-cta-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  color: #334155;
  font-size: 0.8125rem;
  font-weight: 600;
  padding: 7px 13px;
  border-radius: 8px;
  flex-shrink: 0;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.id-card-preview-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #ffffff;
  border: 1px solid #bfdbfe;
  border-radius: 12px;
  padding: 10px 14px;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.05);
}

.preview-thumb-wrap {
  width: 46px;
  height: 46px;
  border-radius: 8px;
  overflow: hidden;
  background: #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid #e2e8f0;
}

.preview-thumb-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-pdf-badge {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #dc2626;
  background: #fee2e2;
  width: 100%;
  height: 100%;
  gap: 1px;
}

.preview-pdf-badge span {
  font-size: 0.75rem;
  font-weight: 700;
  line-height: 1;
}

.preview-details {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.preview-filename {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.preview-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.75rem;
}

.meta-size {
  color: #64748b;
  font-weight: 500;
}

.meta-check {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #16a34a;
  font-weight: 600;
  background: #ecfdf5;
  padding: 1px 7px;
  border-radius: 6px;
}

.preview-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.btn-preview-change {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #2563eb;
  font-size: 0.8125rem;
  font-weight: 600;
  padding: 6px 10px;
  border-radius: 7px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-preview-change:hover {
  background: #dbeafe;
}

.btn-preview-delete {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: #ffffff;
  border: 1px solid #fecaca;
  color: #ef4444;
  border-radius: 7px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-preview-delete:hover {
  background: #fee2e2;
}

.upload-error-msg {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  font-size: 0.8125rem;
  color: #dc2626;
  background: #fef2f2;
  border: 1px solid #fecaca;
  padding: 6px 10px;
  border-radius: 7px;
}
</style>
