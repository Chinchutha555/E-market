<template>
  <div v-if="show" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content" style="max-width: 660px;">
      <AppModalHeader
        title="ยื่นคำร้องขอเช่าแผงค้า"
        :subtitle="`แผง ${targetStallNumber} · อ้างอิงการจอง ${booking?.booking_number || '-'}`"
        icon="document"
        @close="$emit('close')"
      />

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
            style="background: var(--bg-eff6ff); border: 1px solid var(--bd-bfdbfe); border-radius: 10px; padding: 12px 14px; margin-bottom: 16px;"
          >
            <div style="display: flex; flex-wrap: wrap; justify-content: space-between; gap: 10px; font-size: 0.85rem; color: var(--tx-1e3a8a);">
              <div>
                <span>แผงค้าที่ยื่นคำร้อง: </span>
                <strong>{{ targetStallNumber }} ({{ selectedStallRecords.length }} แผง)</strong>
              </div>
              <div style="display: flex; gap: 14px; flex-wrap: wrap;">
                <span>ค่าเช่ารวม: <strong>฿{{ calculatedMonthlyRent.toLocaleString() }}/เดือน</strong></span>
                <span>เงินมัดจำรวม: <strong>฿{{ calculatedDeposit.toLocaleString() }}</strong></span>
                <span>รวมแรกเข้า: <strong style="color: var(--tx-1d4ed8);">฿{{ (calculatedMonthlyRent + calculatedDeposit).toLocaleString() }}</strong></span>
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

          <div class="form-group">
            <label class="form-label">คำนำหน้าชื่อ <b class="req-star">*</b></label>
            <AppSelect v-model="applicantTitle" :options="applicantTitleOptions" placeholder="เลือกคำนำหน้า" :error="titleError" />
          </div>

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

          <AppFileUpload
            v-model="idCardCopy"
            class="form-group"
            label="สำเนาบัตรประชาชน (ID Card Copy)"
            accept=".pdf,.jpg,.jpeg,.png,application/pdf,image/jpeg,image/png"
            :max-size-mb="3"
            hint="PDF, JPG หรือ PNG ขนาดไม่เกิน 3MB"
            description="แนบสำเนาบัตรประจำตัวประชาชน พร้อมรับรองสำเนาถูกต้อง"
            :error="idCardCopyError"
            @error="idCardCopyError = $event"
          />

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
import AppFileUpload from '~/components/AppFileUpload.vue'
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

const removeIDCardFile = () => {
  idCardCopy.value = null
  idCardCopyError.value = ''
}

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

const applicantTitle = ref('')
const titleError = ref('')
const applicantTitleOptions = [
  { value: 'นาย', label: 'นาย' },
  { value: 'นาง', label: 'นาง' },
  { value: 'นางสาว', label: 'นางสาว' }
]
const splitApplicantTitle = (value: string) => {
  const name = String(value || '').trim()
  const title = name.match(/^(นางสาว|นาย|นาง)(?=\s|[^\s])/u)?.[0] || ''
  return { title, name: title ? name.slice(title.length).trim() : name }
}
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
  return !titleError.value && !Object.values(errors).some(Boolean)
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
      const parsedName = splitApplicantTitle(base.tenant_name || base.applicant_name || '')
      applicantTitle.value = parsedName.title
      form.value.applicant_name = parsedName.name
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
        applicant_name: applicantTitle.value + ' ' + form.value.applicant_name.trim(),
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
        proposed_rent_type: rec.rental_type || form.value.proposed_rent_type || 'monthly'
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
  color: var(--tx-1e293b);
  margin: 0;
}

.upload-badge-spec {
  font-size: 0.75rem;
  color: var(--tx-64748b);
  background: var(--bg-f1f5f9);
  padding: 3px 8px;
  border-radius: 6px;
  font-weight: 500;
}

.id-card-dropzone {
  display: flex;
  align-items: center;
  gap: 14px;
  background: var(--bg-f8fafc);
  border: 1.5px dashed var(--bd-cbd5e1);
  border-radius: 12px;
  padding: 14px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.id-card-dropzone:hover,
.id-card-dropzone.is-drag-over {
  background: var(--bg-eff6ff);
  border-color: #2563eb;
  transform: translateY(-1px);
}

.id-card-dropzone.has-error {
  background: var(--bg-fef2f2);
  border-color: var(--bd-f87171);
}

.dropzone-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
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
  color: var(--tx-1e293b);
  margin-bottom: 2px;
}

.prompt-link {
  color: var(--tx-2563eb);
  text-decoration: underline;
  text-underline-offset: 2px;
}

.prompt-desc {
  font-size: 0.8125rem;
  color: var(--tx-64748b);
  line-height: 1.4;
}

.dropzone-cta-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-cbd5e1);
  color: var(--tx-334155);
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
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-bfdbfe);
  border-radius: 12px;
  padding: 10px 14px;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.05);
}

.preview-thumb-wrap {
  width: 46px;
  height: 46px;
  border-radius: 8px;
  overflow: hidden;
  background: var(--bg-f1f5f9);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid var(--bd-e2e8f0);
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
  color: var(--tx-dc2626);
  background: var(--bg-fee2e2);
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
  color: var(--tx-1e293b);
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
  color: var(--tx-64748b);
  font-weight: 500;
}

.meta-check {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: var(--tx-16a34a);
  font-weight: 600;
  background: var(--bg-ecfdf5);
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
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-bfdbfe);
  color: var(--tx-2563eb);
  font-size: 0.8125rem;
  font-weight: 600;
  padding: 6px 10px;
  border-radius: 7px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-preview-change:hover {
  background: var(--bg-dbeafe);
}

.btn-preview-delete {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-fecaca);
  color: var(--tx-ef4444);
  border-radius: 7px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-preview-delete:hover {
  background: var(--bg-fee2e2);
}

.upload-error-msg {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  font-size: 0.8125rem;
  color: var(--tx-dc2626);
  background: var(--bg-fef2f2);
  border: 1px solid var(--bd-fecaca);
  padding: 6px 10px;
  border-radius: 7px;
}
</style>
