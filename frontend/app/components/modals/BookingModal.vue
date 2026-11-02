<template>
  <div v-if="show" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content" style="max-width: 600px;">
      <AppModalHeader
        :title="isQueueMode ? 'ต่อคิวแผงติดจอง' : 'จองนัดหมายดูแผงเช่า'"
        :subtitle="`แผง ${stallShortNumber(stall?.stall_number) || '-'} · โซน ${stall?.zone_code || stall?.zone?.code || 'ทั่วไป'}`"
        icon="pin"
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
        <h4 style="font-size: var(--font-size-xl); font-weight: 700; color: var(--tx-059669); margin-bottom: 8px;">
          {{ isQueueMode ? 'คุณเข้าสู่ลำดับคิวเรียบร้อยแล้ว!' : 'ส่งคำขอจองดูแผงเรียบร้อยแล้ว!' }}
        </h4>
        <div style="font-size: var(--font-size-sm); color: var(--text-muted); line-height: 1.6; max-width: 440px; margin: 0 auto 24px;">
          เลขที่คำขอจอง: <strong style="color: var(--tx-0f172a);">{{ successBookingNumber }}</strong><br/>
          <template v-if="isQueueMode">
            สถานะของคุณ: <strong style="color: var(--tx-d97706);">รอคิว (Waiting List)</strong><br/>
            เมื่อผู้ถือแผงคนก่อนหน้าสละสิทธิ์หรือหมดเวลา ระบบจะเลื่อนคิวของคุณขึ้นมาถือแผงและส่งการแจ้งเตือนให้เลือกวันนัดชมทันที
          </template>
          <template v-else>
            ขั้นตอนถัดไป: <strong>เจ้าหน้าที่ตลาดจะตรวจสอบและอนุมัติวันเวลานัดหมายเข้าชม</strong> จากนั้นคุณสามารถเข้าดูแผงหน้างาน และกดยืนยันความสนใจในระบบได้ทันที
          </template>
        </div>
        <div style="display: flex; justify-content: center; gap: 12px;">
          <AppButton to="/bookings" variant="primary" size="md" @click="$emit('close')">ดูรายการการจองของฉัน</AppButton>
          <AppButton variant="secondary" size="md" @click="$emit('close')">ปิดหน้าต่าง</AppButton>
        </div>
      </div>

      <!-- Form State -->
      <form v-else novalidate @submit.prevent="handleSubmit">
        <div class="modal-body">
          <!-- Queue Mode Notice -->
          <div v-if="isQueueMode" style="background: var(--bg-fffbeb); border: 1px solid var(--bd-fde68a); border-radius: 10px; padding: 12px 16px; margin-bottom: 16px; display: flex; align-items: center; gap: 12px;">
            <div style="color: var(--tx-d97706); display: flex; align-items: center; justify-content: center; flex-shrink: 0;">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 22h14"></path><path d="M5 2h14"></path><path d="M17 22v-4.172a2 2 0 0 0-.586-1.414L12 12l-4.414 4.414A2 2 0 0 0 7 17.828V22"></path><path d="M7 2v4.172a2 2 0 0 0 .586 1.414L12 12l4.414-4.414A2 2 0 0 0 17 6.172V2"></path></svg>
            </div>
            <div style="font-size: var(--font-size-sm); color: var(--tx-92400e);">
              <strong>แผงนี้มีผู้จองถือสิทธิ์อยู่:</strong> ท่านกำลังทำรายการ <strong>"ต่อคิว"</strong> (จะได้คิวถัดไปอัตโนมัติ โดยไม่ต้องเลือกวันนัดหมายในขั้นตอนนี้)
            </div>
          </div>

          <!-- Stall Summary Info Banner -->
          <div style="background: rgba(16, 185, 129, 0.08); border: 1px solid var(--bd-a7f3d0); padding: 14px; border-radius: var(--radius-md); margin-bottom: 20px; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 10px;">
            <div>
              <div style="font-weight: 700; font-size: var(--font-size-md); color: var(--tx-0f172a);">แผง {{ stallShortNumber(stall?.stall_number) }} - {{ stall?.category || stall?.category_name || 'ทั่วไป' }}</div>
              <div style="font-size: var(--font-size-xs); color: var(--text-muted); margin-top: 2px;">
                ขนาด: {{ stall?.size_width || 2 }} × {{ stall?.size_length || 2 }} ม. ({{ stall?.area_sqm || 4 }} ตร.ม.)
              </div>
            </div>
            <div style="text-align: right;">
              <div style="font-size: var(--font-size-lg); font-weight: 800; color: var(--tx-059669);">฿{{ Number(stall?.monthly_rate || 2500).toLocaleString() }}</div>
              <div style="font-size: var(--font-size-xs); color: var(--text-dim);">บาท / เดือน (รายวัน: ฿{{ stall?.daily_rate || 100 }})</div>
            </div>
          </div>

          <h4 style="font-size: var(--font-size-sm); font-weight: 700; color: var(--text-muted); margin-bottom: 12px; text-transform: uppercase; letter-spacing: 0.05em;">
            1. ข้อมูลผู้ประสงค์เข้าดูแผง
          </h4>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label" for="b_name">ชื่อ-นามสกุล <b class="req-star">*</b></label>
              <AppInput id="b_name" v-model="form.tenant_name" type="text" placeholder="เช่น นายสมเกียรติ มั่นคง" :error="errors.tenant_name" @blur="validateField('tenant_name')" />
            </div>

            <div class="form-group">
              <label class="form-label" for="b_phone">เบอร์โทรศัพท์ติดต่อ <b class="req-star">*</b></label>
              <AppInput id="b_phone" v-model="form.tenant_phone" type="tel" placeholder="08x-xxx-xxxx" :error="errors.tenant_phone" @blur="validateField('tenant_phone')" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label" for="b_email">อีเมล</label>
              <AppInput id="b_email" v-model="form.tenant_email" type="email" placeholder="name@example.com" />
            </div>

            <div class="form-group">
              <label class="form-label" for="b_id_card">เลขประจำตัวประชาชน (ถ้ามี)</label>
              <AppInput id="b_id_card" v-model="form.tenant_id_card" type="text" placeholder="13 หลัก" />
            </div>
          </div>

          <h4 style="font-size: var(--font-size-sm); font-weight: 700; color: var(--text-muted); margin: 16px 0 12px; text-transform: uppercase; letter-spacing: 0.05em;">
            {{ isQueueMode ? '2. สินค้าที่ต้องการขาย' : '2. ข้อมูลนัดหมายและสินค้าที่ต้องการขาย' }}
          </h4>

          <!-- Only show date/time selection if NOT queue mode (Page 2 Rule) -->
          <div v-if="!isQueueMode" class="form-row">
            <div class="form-group">
              <label class="form-label" for="b_date">วันที่ต้องการเข้าชมแผง <b class="req-star">*</b></label>
              <AppInput id="b_date" v-model="form.visit_date" type="date" :min="minDate" :error="errors.visit_date" @blur="validateField('visit_date')" />
            </div>

            <div class="form-group">
              <label class="form-label" for="b_time">ช่วงเวลาที่สะดวก <b class="req-star">*</b></label>
              <AppSelect v-model="form.visit_time_slot" :options="visitTimeOptions" />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label" for="b_business">ประเภทสินค้า / ร้านค้าที่ตั้งใจจะขาย</label>
            <AppInput id="b_business" v-model="form.intended_business" type="text" placeholder="เช่น ร้านก๋วยเตี๋ยวเรือ, เสื้อผ้าแฟชั่นเกาหลี" />
          </div>
        </div>

        <div class="modal-footer">
          <AppButton type="button" variant="secondary" size="md" @click="$emit('close')">ยกเลิก</AppButton>
          <AppButton type="submit" :variant="isQueueMode ? 'primary' : 'success'" size="md" :loading="submitting">
            {{ isQueueMode ? 'ยืนยันการต่อคิวแผงนี้' : 'ยืนยันการจองดูแผง' }}
          </AppButton>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  show: boolean
  stall: any
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

const { post } = useApi()
const { user } = useAuth()
const { showError, showWarning } = useAppAlert()

const submitting = ref(false)
const isSuccess = ref(false)
const successBookingNumber = ref('')

const isQueueMode = computed(() => {
  return props.stall?.status === 'reserved' || props.stall?.status === 'occupied' || Boolean(props.stall?.is_queue)
})

const minDate = new Date().toISOString().slice(0, 10)

const visitTimeOptions = [
  { value: '09:00 - 10:30', label: 'ช่วงเช้า (09:00 - 10:30 น.)' },
  { value: '10:30 - 12:00', label: 'ช่วงสาย (10:30 - 12:00 น.)' },
  { value: '13:30 - 15:00', label: 'ช่วงบ่าย (13:30 - 15:00 น.)' },
  { value: '15:00 - 16:30', label: 'ช่วงบ่ายแก่ (15:00 - 16:30 น.)' }
]

const form = ref({
  tenant_name: 'คุณสมศักดิ์ ค้าขายดี',
  tenant_phone: '081-999-8877',
  tenant_email: 'somsak@example.com',
  tenant_id_card: '1-1004-88991-23-4',
  visit_date: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString().slice(0, 10), // tomorrow
  visit_time_slot: '10:30 - 12:00',
  intended_business: ''
})

const errors = reactive({
  tenant_name: '',
  tenant_phone: '',
  visit_date: ''
})

const validateField = (field: keyof typeof errors) => {
  if (field === 'tenant_name') {
    errors.tenant_name = form.value.tenant_name.trim() ? '' : 'กรุณาระบุชื่อ-นามสกุล'
  } else if (field === 'tenant_phone') {
    errors.tenant_phone = form.value.tenant_phone.trim() ? '' : 'กรุณาระบุเบอร์โทรศัพท์ติดต่อ'
  } else if (field === 'visit_date') {
    errors.visit_date = !isQueueMode.value && !form.value.visit_date ? 'กรุณาเลือกวันที่ต้องการเข้าชมแผง' : ''
  }
}

const validateForm = () => {
  validateField('tenant_name')
  validateField('tenant_phone')
  validateField('visit_date')
  return !errors.tenant_name && !errors.tenant_phone && !errors.visit_date
}

watch(() => props.show, (newVal) => {
  if (newVal) {
    isSuccess.value = false
    errors.tenant_name = ''
    errors.tenant_phone = ''
    errors.visit_date = ''
    const cleanedName = String(user.value?.full_name || 'คุณสมชาย ใจดี').replace(/\([^)]*\)/g, '').trim()
    form.value.tenant_name = cleanedName || 'คุณสมชาย ใจดี'
    form.value.tenant_email = user.value?.email || 'somchai@example.com'
    form.value.intended_business = props.stall?.category || props.stall?.category_name || ''
  }
})

const handleSubmit = async () => {
  if (!props.stall) return
  if (!validateForm()) return
  submitting.value = true

  try {
    const res: any = await post('/bookings', {
      stall_ident: props.stall.stall_number,
      applicant_name: form.value.tenant_name,
      applicant_phone: form.value.tenant_phone,
      applicant_email: form.value.tenant_email,
      business_type: form.value.intended_business,
      visit_time: isQueueMode.value ? 'รอเลือกเวลาเมื่อถึงคิว' : form.value.visit_time_slot,
      stall_id: props.stall.id,
      stall_number: props.stall.stall_number,
      market_name: props.stall.market_name,
      tenant_name: form.value.tenant_name,
      tenant_phone: form.value.tenant_phone,
      tenant_email: form.value.tenant_email,
      tenant_id_card: form.value.tenant_id_card,
      visit_date: isQueueMode.value ? 'รอเลือกวันเมื่อถึงคิว' : form.value.visit_date,
      visit_time_slot: isQueueMode.value ? undefined : form.value.visit_time_slot,
      intended_business: form.value.intended_business,
      is_queue: isQueueMode.value
    })

    const created = res?.data || res
    successBookingNumber.value = created?.booking_number || created?.booking_no || ''
    isSuccess.value = true
    emit('success')
  } catch (err: any) {
    if (err?.status === 409 || err?.statusCode === 409 || err?.data?.code === 409) {
      await showWarning(err?.data?.detail || 'คุณได้ทำการจองแผงนี้ไว้แล้ว ไม่สามารถจองซ้ำได้', 'ไม่สามารถจองแผงซ้ำได้')
    } else {
      await showError(err?.data?.detail || 'เกิดข้อผิดพลาดในการส่งคำขอจอง', 'ส่งคำขอไม่สำเร็จ')
    }
  } finally {
    submitting.value = false
  }
}
</script>
