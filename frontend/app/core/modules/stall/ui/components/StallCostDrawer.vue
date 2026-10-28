<template>
  <aside class="cd-drawer" role="dialog" :aria-label="kind === 'utility' ? 'จดมิเตอร์รายแผง' : 'ออกบิลแผงเช่าชั่วคราว'">
    <header class="cd-head">
      <div>
        <span class="cd-eyebrow">{{ kind === 'utility' ? 'ค่าน้ำ ค่าไฟ' : 'บิลแผงเช่าชั่วคราว' }}</span>
        <h3>{{ stallCode ? `แผง ${stallCode}` : 'ออกบิลใหม่' }}</h3>
      </div>
      <button type="button" class="cd-close" aria-label="ปิด" @click="$emit('close')">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 6 12 12M18 6 6 18"></path></svg>
      </button>
    </header>

    <div class="cd-body">
      <!-- ================= ค่าน้ำ ค่าไฟ ================= -->
      <template v-if="kind === 'utility'">
        <p v-if="!bill" class="cd-note is-warn">แผงนี้ไม่มีผู้เช่า จึงจดมิเตอร์ไม่ได้ (ไม่มีผู้ใช้บริการ)</p>
        <template v-else>
          <div class="cd-tenant">
            <div><small>ผู้เช่า</small><strong>{{ bill.tenant }}</strong></div>
            <span class="cd-badge" :class="`is-${bill.status}`">{{ BILL_STATUS_LABEL[bill.status] }}</span>
          </div>

          <fieldset class="cd-meter">
            <legend>มิเตอร์น้ำ</legend>
            <div class="cd-row"><span>เลขก่อนหน้า</span><strong>{{ bill.prevWater.toLocaleString('th-TH') }}</strong></div>
            <label>เลขปัจจุบัน<AppInput v-model="form.water" type="number" placeholder="กรอกเลขมิเตอร์" :error="waterError" /></label>
            <div class="cd-row"><span>ใช้ไป</span><strong>{{ units(form.water, bill.prevWater) }} หน่วย × ฿{{ rates.water }}</strong></div>
          </fieldset>

          <fieldset class="cd-meter">
            <legend>มิเตอร์ไฟ</legend>
            <div class="cd-row"><span>เลขก่อนหน้า</span><strong>{{ bill.prevElectric.toLocaleString('th-TH') }}</strong></div>
            <label>เลขปัจจุบัน<AppInput v-model="form.electric" type="number" placeholder="กรอกเลขมิเตอร์" :error="electricError" /></label>
            <div class="cd-row"><span>ใช้ไป</span><strong>{{ units(form.electric, bill.prevElectric) }} หน่วย × ฿{{ rates.electric }}</strong></div>
          </fieldset>

          <div class="cd-total">
            <span>รวมสุทธิ</span>
            <strong>฿{{ utilityPreview.toLocaleString('th-TH') }}</strong>
          </div>
        </template>
      </template>

      <!-- ================= บิลแผงเช่าชั่วคราว ================= -->
      <template v-else>
        <label v-if="!lockedStall" class="cd-field">รหัสแผง
          <AppInput v-model="codeInput" placeholder="เช่น 12-A05" :error="codeError" />
        </label>

        <p v-if="stallInfo && stallInfo.occupied" class="cd-note is-warn">แผงนี้มีผู้เช่าตามสัญญาอยู่ ออกบิลเช่าชั่วคราวไม่ได้</p>

        <template v-else-if="stallInfo">
          <section v-if="existingBills.length" class="cd-existing">
            <h4>บิลของแผงนี้</h4>
            <div v-for="b in existingBills" :key="b.id" class="cd-existing-row">
              <div>
                <strong>{{ b.id }} · {{ b.customer }}</strong>
                <small>{{ b.startDate }} ถึง {{ b.endDate }} · ฿{{ tempTotal(b).toLocaleString('th-TH') }}</small>
              </div>
              <span class="cd-badge" :class="`is-${b.status}`">{{ BILL_STATUS_LABEL[b.status] }}</span>
              <AppButton v-if="b.status !== 'paid'" size="sm" variant="outline" @click="markTempPaid(b.id)">บันทึกชำระ</AppButton>
            </div>
          </section>

          <h4 class="cd-sub">ออกบิลใหม่</h4>
          <label class="cd-field">ชื่อผู้เช่าชั่วคราว<AppInput v-model="temp.customer" placeholder="ชื่อ-นามสกุล" :error="tempErrors.customer" /></label>
          <label class="cd-field">เบอร์โทร<AppInput v-model="temp.phone" placeholder="08x-xxx-xxxx" :error="tempErrors.phone" /></label>
          <div class="cd-two">
            <label class="cd-field">วันที่เริ่ม<AppInput v-model="temp.startDate" type="date" /></label>
            <label class="cd-field">วันที่สิ้นสุด<AppInput v-model="temp.endDate" type="date" :error="tempErrors.endDate" /></label>
          </div>
          <label class="cd-field">อัตรารายวัน (บาท)<AppInput v-model="temp.dailyRate" type="number" /></label>
          <div class="cd-two">
            <label class="cd-field">ค่าธรรมเนียมเพิ่ม<AppInput v-model="temp.extra" type="number" /></label>
            <label class="cd-field">ส่วนลด<AppInput v-model="temp.discount" type="number" /></label>
          </div>

          <p v-if="conflict" class="cd-note is-warn">ช่วงวันที่ซ้อนกับบิล {{ conflict.id }} ({{ conflict.startDate }} ถึง {{ conflict.endDate }}) กรุณาเลือกช่วงอื่น</p>

          <div class="cd-total">
            <span>{{ days }} วัน × ฿{{ Number(temp.dailyRate) || 0 }}</span>
            <strong>฿{{ tempPreview.toLocaleString('th-TH') }}</strong>
          </div>
        </template>
      </template>
    </div>

    <footer class="cd-foot">
      <template v-if="kind === 'utility' && bill">
        <AppButton variant="secondary" @click="$emit('close')">ปิด</AppButton>
        <AppButton v-if="(bill.status === 'pending' || bill.status === 'overdue') && !dirty" variant="success" @click="markUtilityPaid(bill.stall)">บันทึกชำระ</AppButton>
        <AppButton v-else variant="primary" :disabled="!utilityValid" @click="saveUtility">{{ bill.status === 'unrecorded' ? 'บันทึกและออกบิล' : 'บันทึกการแก้ไข' }}</AppButton>
      </template>
      <template v-else-if="kind === 'temp' && stallInfo && !stallInfo.occupied">
        <AppButton variant="secondary" @click="$emit('close')">ปิด</AppButton>
        <AppButton variant="primary" :disabled="!!conflict" @click="issueTemp">ออกบิล</AppButton>
      </template>
      <AppButton v-else variant="secondary" @click="$emit('close')">ปิด</AppButton>
    </footer>
  </aside>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import AppButton from '~/components/AppButton.vue'
import AppInput from '~/components/AppInput.vue'
import { BILL_STATUS_LABEL, useStallBilling } from '../../composables/useStallBilling'

// ข้อมูลเป็นข้อมูลจำลอง (mock) ยังไม่เชื่อม API
const props = defineProps<{
  kind: 'utility' | 'temp'
  code: string
  /** ตรวจรหัสแผง: คืน null ถ้าไม่พบแผงในตลาดนี้ */
  lookup: (code: string) => { occupied: boolean; dailyRate: number } | null
}>()
const emit = defineEmits<{ close: []; issued: [code: string] }>()

const {
  rates, getUtility, saveUtilityReadings, markUtilityPaid,
  tempBillsOf, tempDays, tempTotal, tempConflict, issueTempBill, markTempPaid
} = useStallBilling()

// ---------- utility ----------
const bill = computed(() => (props.kind === 'utility' ? getUtility(props.code) : null))
const form = reactive({ water: '', electric: '' })
const resetForm = () => {
  form.water = bill.value?.currWater == null ? '' : String(bill.value.currWater)
  form.electric = bill.value?.currElectric == null ? '' : String(bill.value.currElectric)
}
watch(() => [props.code, props.kind], resetForm, { immediate: true })

const units = (curr: string, prev: number) => (curr === '' ? '-' : Math.max(0, Number(curr) - prev))
const waterError = computed(() => (form.water !== '' && bill.value && Number(form.water) < bill.value.prevWater ? 'เลขปัจจุบันต้องไม่น้อยกว่าเลขก่อนหน้า' : ''))
const electricError = computed(() => (form.electric !== '' && bill.value && Number(form.electric) < bill.value.prevElectric ? 'เลขปัจจุบันต้องไม่น้อยกว่าเลขก่อนหน้า' : ''))
const utilityValid = computed(() => form.water !== '' && form.electric !== '' && !waterError.value && !electricError.value)
const utilityPreview = computed(() => {
  if (!bill.value || !utilityValid.value) return 0
  return (Number(form.water) - bill.value.prevWater) * rates.water + (Number(form.electric) - bill.value.prevElectric) * rates.electric
})
const dirty = computed(() =>
  !!bill.value && (form.water !== String(bill.value.currWater ?? '') || form.electric !== String(bill.value.currElectric ?? ''))
)
const saveUtility = () => {
  if (!bill.value || !utilityValid.value) return
  saveUtilityReadings(bill.value.stall, Number(form.water), Number(form.electric))
  resetForm()
}

// ---------- temp bill ----------
const lockedStall = computed(() => !!props.code)
const codeInput = ref(props.code)
watch(() => props.code, (c) => { codeInput.value = c })
const stallCode = computed(() => (props.kind === 'temp' ? codeInput.value.trim().toUpperCase() : props.code))
const stallInfo = computed(() => (stallCode.value ? props.lookup(stallCode.value) : null))
const codeError = computed(() => (!lockedStall.value && stallCode.value && !stallInfo.value ? 'ไม่พบรหัสแผงนี้ในตลาด' : ''))
const existingBills = computed(() => tempBillsOf(stallCode.value))

const today = new Date().toISOString().slice(0, 10)
const temp = reactive({ customer: '', phone: '', startDate: today, endDate: today, dailyRate: '150', extra: '0', discount: '0' })
watch(stallInfo, (i) => { if (i) temp.dailyRate = String(i.dailyRate || 150) }, { immediate: true })
const tempErrors = reactive({ customer: '', phone: '', endDate: '' })
const days = computed(() => tempDays(temp))
const tempPreview = computed(() => tempTotal({ ...temp, dailyRate: Number(temp.dailyRate) || 0, extra: Number(temp.extra) || 0, discount: Number(temp.discount) || 0 }))
const conflict = computed(() => (stallInfo.value && days.value > 0 ? tempConflict(stallCode.value, temp.startDate, temp.endDate) : null))

const issueTemp = () => {
  tempErrors.customer = temp.customer.trim() ? '' : 'กรุณากรอกชื่อผู้เช่า'
  tempErrors.phone = temp.phone.trim() ? '' : 'กรุณากรอกเบอร์โทร'
  tempErrors.endDate = days.value > 0 ? '' : 'วันที่สิ้นสุดต้องไม่ก่อนวันที่เริ่ม'
  if (tempErrors.customer || tempErrors.phone || tempErrors.endDate || conflict.value) return
  issueTempBill({
    stall: stallCode.value,
    customer: temp.customer.trim(),
    phone: temp.phone.trim(),
    startDate: temp.startDate,
    endDate: temp.endDate,
    dailyRate: Number(temp.dailyRate) || 0,
    extra: Number(temp.extra) || 0,
    discount: Number(temp.discount) || 0
  })
  temp.customer = ''
  temp.phone = ''
  emit('issued', stallCode.value)
}
</script>

<style scoped>
.cd-drawer { position: absolute; z-index: 65; top: calc(var(--mp-top, 114px) + 58px); right: 68px; bottom: 18px; box-sizing: border-box; width: min(380px, calc(100% - 36px)); display: flex; flex-direction: column; background: #fff; border: 1px solid #e2e8f0; border-radius: 16px; box-shadow: 0 18px 44px rgba(15, 23, 42, .22); overflow: hidden; }
.cd-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; padding: 16px; border-bottom: 1px solid #f1f5f9; }
.cd-eyebrow { color: #2563eb; font-size: .8125rem; font-weight: 600; }
.cd-head h3 { margin: 2px 0 0; color: #0f172a; font-size: 1.125rem; font-weight: 700; }
.cd-close { flex: 0 0 auto; width: 34px; height: 34px; display: grid; place-items: center; border: 1px solid #e2e8f0; border-radius: 8px; background: #fff; color: #475569; cursor: pointer; }
.cd-close:hover { background: #f1f5f9; }
.cd-body { flex: 1 1 auto; min-height: 0; overflow-y: auto; padding: 16px; display: flex; flex-direction: column; gap: 14px; -webkit-overflow-scrolling: touch; }
.cd-foot { display: flex; justify-content: flex-end; gap: 8px; padding: 12px 16px; border-top: 1px solid #f1f5f9; }
.cd-tenant { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 12px 14px; border: 1px solid #eef2f6; border-radius: 12px; background: #f8fafc; }
.cd-tenant small { display: block; color: #64748b; font-size: .8125rem; font-weight: 500; }
.cd-tenant strong { color: #1e293b; font-size: .9375rem; font-weight: 700; }
.cd-meter { margin: 0; padding: 12px 14px 14px; border: 1px solid #e2e8f0; border-radius: 12px; display: flex; flex-direction: column; gap: 10px; }
.cd-meter legend { padding: 0 6px; color: #334155; font-size: .875rem; font-weight: 700; }
.cd-meter label, .cd-field { display: flex; flex-direction: column; gap: 6px; color: #475569; font-size: .875rem; font-weight: 600; }
.cd-row { display: flex; justify-content: space-between; gap: 8px; color: #64748b; font-size: .875rem; font-weight: 500; }
.cd-row strong { color: #1e293b; font-weight: 700; }
.cd-two { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.cd-total { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 12px 14px; border-radius: 12px; background: #eff6ff; color: #475569; font-size: .9375rem; font-weight: 600; }
.cd-total strong { color: #2563eb; font-size: 1.25rem; font-weight: 800; }
.cd-note { margin: 0; padding: 10px 12px; border-radius: 10px; font-size: .875rem; font-weight: 500; line-height: 1.5; }
.cd-note.is-warn { background: #fef3c7; color: #b45309; }
.cd-sub { margin: 4px 0 -4px; color: #1e293b; font-size: .9375rem; font-weight: 700; }
.cd-existing h4 { margin: 0 0 8px; color: #1e293b; font-size: .9375rem; font-weight: 700; }
.cd-existing-row { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; padding: 10px 0; border-bottom: 1px solid #f1f5f9; }
.cd-existing-row > div { flex: 1 1 140px; min-width: 0; }
.cd-existing-row strong { display: block; color: #1e293b; font-size: .875rem; font-weight: 700; }
.cd-existing-row small { color: #94a3b8; font-size: .8125rem; font-weight: 500; }
.cd-badge { display: inline-block; padding: 4px 10px; border-radius: 8px; font-size: .8125rem; font-weight: 600; white-space: nowrap; }
.cd-badge.is-pending { background: #fef3c7; color: #d97706; }
.cd-badge.is-paid { background: #dcfce7; color: #15803d; }
.cd-badge.is-overdue { background: #fee2e2; color: #dc2626; }
.cd-badge.is-unrecorded { background: #f1f5f9; color: #64748b; }
@media (max-width: 1000px) { .cd-drawer { top: calc(var(--mp-top-2row, 208px) + 54px); } }
@media (max-width: 720px) { .cd-drawer { top: auto; right: 10px; left: 10px; bottom: 10px; width: auto; max-height: 70%; } }
</style>
