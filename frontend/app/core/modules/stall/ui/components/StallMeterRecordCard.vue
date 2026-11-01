<template>
  <aside class="mrc" role="dialog" :aria-label="`จดมิเตอร์ ${label}`">
    <div class="mrc-top">
      <AppButton variant="primary" size="sm" @click="emit('back')">
        <template #prefix>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="15 18 9 12 15 6"></polyline></svg>
        </template>
        {{ projectName }}
      </AppButton>
    </div>

    <header class="mrc-head">
      <h3>{{ label }}</h3>
      <div v-if="bill" class="mrc-sub">
        <span class="mrc-pill" :class="`is-${state}`">{{ stateLabel[state] }}</span>
        <span class="mrc-remaining">{{ pendingLeft > 0 ? `เหลือยังไม่จด ${pendingLeft} แผง` : 'จดครบทุกแผงแล้ว' }}</span>
      </div>
    </header>

    <div class="mrc-body">
      <p v-if="!bill" class="mrc-note is-warn">แผงนี้ไม่มีสัญญาที่มีผล หรือไม่ได้ตั้งค่าให้คิดตามมิเตอร์ จึงยังไม่มีรายการจด</p>
      <template v-else>
        <div class="mrc-tenant">
          <small>ผู้เช่าตามสัญญา</small>
          <strong>{{ bill.tenant }}</strong>
          <span v-if="bill.shopName">{{ bill.shopName }}</span>
          <span v-if="bill.contractNumber">สัญญา {{ bill.contractNumber }}</span>
        </div>

        <section v-for="meter in meters" :key="meter.key" class="mrc-meter" :class="{ 'has-error': !!meter.error }">
          <h4>{{ meter.title }}<small v-if="meter.meterCode" class="mrc-code"> · {{ meter.meterCode }}</small></h4>
          <div class="mrc-row"><span>เลขครั้งก่อน</span><strong>{{ meter.prev.toLocaleString('th-TH') }}</strong></div>
          <label class="mrc-field">
            <span>เลขครั้งนี้</span>
            <AppInput v-model="form[meter.key]" type="number" placeholder="กรอกเลขมิเตอร์" :error="meter.error" />
          </label>
          <div class="mrc-row">
            <span>ใช้ไป</span>
            <strong>{{ meter.usedText }} <small>{{ meter.unit }} × ฿{{ meter.rate }}</small></strong>
          </div>
          <p v-if="meter.warning" class="mrc-note is-warn">{{ meter.warning }}</p>
        </section>

        <div v-for="auto in autoCharges" :key="auto.key" class="mrc-row mrc-auto"><span>{{ auto.title }}</span><strong>{{ auto.text }}</strong></div>

        <div class="mrc-total"><span>ยอดค่าน้ำไฟประมาณ</span><strong>฿{{ preview.toLocaleString('th-TH') }}</strong></div>

        <div v-if="skipOpen" class="mrc-skip">
          <label class="mrc-field">
            <span>เหตุผลที่ข้าม</span>
            <AppSelect v-model="skipReason" :options="skipOptions" placeholder="เลือกเหตุผล" />
          </label>
          <div class="mrc-skip-actions">
            <AppButton variant="secondary" size="md" @click="skipOpen = false">ยกเลิก</AppButton>
            <AppButton variant="primary" size="md" :disabled="!skipReason" @click="confirmSkip">ยืนยันข้าม</AppButton>
          </div>
        </div>
      </template>
    </div>

    <footer v-if="bill && !skipOpen" class="mrc-foot">
      <AppButton variant="outline" size="md" @click="skipOpen = true">ข้าม</AppButton>
      <AppButton variant="primary" size="md" class="mrc-save" @click="save">{{ saveLabel }}</AppButton>
    </footer>
  </aside>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import AppButton from '~/components/AppButton.vue'
import AppInput from '~/components/AppInput.vue'
import AppSelect from '~/components/AppSelect.vue'
import { chargeFor, chargeRuleOf, useStallBilling, type ChargeKind } from '../../composables/useStallBilling'
import { useUtilitySettings } from '../../composables/useUtilitySettings'

// การ์ดจดมิเตอร์รายแผง เชื่อมรอบจดกับสัญญา ผู้เช่า มิเตอร์ และอัตราของแผงเดียวกัน
const props = defineProps<{
  code: string
  label: string
  projectName: string
  pendingLeft: number
  hasNext: boolean
  /** ค่าเฉลี่ยหน่วยที่ใช้ของแผงที่จดแล้วในโครงการ (0 = ยังไม่พอเทียบ) ใช้เตือนค่าที่สูงผิดปกติ */
  averageElectric: number
  averageWater: number
}>()
const emit = defineEmits<{ back: []; saved: [code: string]; skipped: [code: string] }>()

const { getUtility, getUtilitySkip, skipUtility, saveUtilityReadings } = useStallBilling()
const { getMeter } = useUtilitySettings(String(useRoute().query.market || '1'))
const ABNORMAL_FACTOR = 2.5

const bill = computed(() => getUtility(props.code))
const form = reactive({ water: '', electric: '' })
const submitted = ref(false)
const skipOpen = ref(false)
const skipReason = ref('')
const skipOptions = [
  { value: 'มิเตอร์ชำรุด', label: 'มิเตอร์ชำรุด' },
  { value: 'ร้านปิด เข้าไม่ได้', label: 'ร้านปิด เข้าไม่ได้' },
  { value: 'อ่านค่าไม่ได้', label: 'อ่านค่าไม่ได้' }
]

const reset = () => {
  submitted.value = false
  skipOpen.value = false
  skipReason.value = ''
  form.water = bill.value?.currWater == null ? '' : String(bill.value.currWater)
  form.electric = bill.value?.currElectric == null ? '' : String(bill.value.currElectric)
}
watch(() => props.code, reset, { immediate: true })

const used = (curr: string, prev: number) => (curr === '' ? null : Number(curr) - prev)
const meterState = (key: 'water' | 'electric') => {
  const prev = key === 'water' ? bill.value?.prevWater ?? 0 : bill.value?.prevElectric ?? 0
  const average = key === 'water' ? props.averageWater : props.averageElectric
  const value = used(form[key], prev)
  let error = ''
  let warning = ''
  if (submitted.value && form[key] === '') error = `กรุณากรอกเลขมิเตอร์${key === 'water' ? 'น้ำ' : 'ไฟ'}`
  else if (value != null && value < 0) error = 'เลขครั้งนี้ต้องไม่น้อยกว่าครั้งก่อน (กรอกผิด หรือเปลี่ยนมิเตอร์?)'
  else if (value != null && average > 0 && value > average * ABNORMAL_FACTOR) warning = `ใช้สูงกว่า ${ABNORMAL_FACTOR} เท่าของค่าเฉลี่ยโครงการ (เฉลี่ย ${Math.round(average).toLocaleString('th-TH')} หน่วย) กรุณาตรวจซ้ำ`
  return { prev, value, error, warning }
}
const waterState = computed(() => meterState('water'))
const electricState = computed(() => meterState('electric'))
const KINDS = ['water', 'electric'] as const
const ruleOf = (kind: ChargeKind) => chargeRuleOf(props.code, kind)
const stateOf = (kind: ChargeKind) => (kind === 'water' ? waterState.value : electricState.value)
// จดเฉพาะประเภทที่คิดตามมิเตอร์ ประเภทเหมาจ่าย/ไม่คิดค่าใช้บริการไม่ต้องกรอกเลข
const meters = computed(() => KINDS.filter(k => ruleOf(k).mode === 'meter').map((key) => {
  const st = stateOf(key)
  return { key, meterCode: getMeter(props.code)[key], title: key === 'water' ? 'มิเตอร์น้ำ' : 'มิเตอร์ไฟ', unit: 'หน่วย', rate: ruleOf(key).rate, ...st, usedText: st.value == null ? '-' : Math.max(0, st.value).toLocaleString('th-TH') }
}))
const autoCharges = computed(() => KINDS.filter(k => ruleOf(k).mode !== 'meter').map(key => ({
  key,
  title: key === 'water' ? 'ค่าน้ำ' : 'ค่าไฟ',
  text: ruleOf(key).mode === 'flat' ? `เหมาจ่าย ฿${ruleOf(key).flat.toLocaleString('th-TH')}/เดือน (ไม่ต้องจดมิเตอร์)` : `${key === 'water' ? 'ไม่คิดค่าน้ำประปา' : 'ไม่คิดค่าไฟฟ้า'} (ไม่ต้องจดมิเตอร์)`
})))
const valid = computed(() => meters.value.every(m => form[m.key] !== '' && !(m.value! < 0)))
const preview = computed(() => KINDS.reduce((sum, k) => sum + (chargeFor(ruleOf(k), stateOf(k).value) ?? 0), 0))
const dirty = computed(() => !!bill.value && (form.water !== String(bill.value.currWater ?? '') || form.electric !== String(bill.value.currElectric ?? '')))
const recorded = computed(() => !!bill.value && bill.value.status !== 'unrecorded')

const state = computed<'recorded' | 'abnormal' | 'pending' | 'skipped'>(() => {
  if (recorded.value) return 'recorded'
  return getUtilitySkip(props.code) ? 'skipped' : 'pending'
})
const stateLabel = { recorded: 'จดแล้ว', abnormal: 'ผิดปกติ', pending: 'ยังไม่จด', skipped: `ข้าม` } as const

const saveLabel = computed(() => {
  if (recorded.value && !dirty.value) return props.hasNext ? 'แผงถัดไป →' : 'เสร็จสิ้น'
  return props.hasNext ? 'บันทึก แล้วไปแผงถัดไป →' : 'บันทึก'
})

const save = () => {
  if (!bill.value) return
  if (recorded.value && !dirty.value) { emit('saved', props.code); return }
  submitted.value = true
  if (!valid.value) return
  const water = ruleOf('water').mode === 'meter' ? Number(form.water) : bill.value.prevWater
  const electric = ruleOf('electric').mode === 'meter' ? Number(form.electric) : bill.value.prevElectric
  saveUtilityReadings(bill.value.stall, water, electric)
  emit('saved', props.code)
}
const confirmSkip = () => {
  if (!skipReason.value) return
  skipUtility(props.code, skipReason.value)
  emit('skipped', props.code)
}
</script>

<style scoped>
.mrc {
  position: absolute;
  z-index: 60;
  top: var(--mp-top, 114px);
  left: 18px;
  bottom: 18px;
  box-sizing: border-box;
  width: min(380px, calc(100% - 36px));
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 18px;
  background: var(--bg-ffffff);
  box-shadow: 0 16px 38px rgba(15, 23, 42, 0.16);
  animation: mrcIn .25s ease both;
}
@keyframes mrcIn { from { opacity: 0; } to { opacity: 1; } }
.mrc-top { flex: 0 0 auto; padding: 18px 20px 12px; border-bottom: 1px solid var(--bd-e8eef6); }
.mrc-head { flex: 0 0 auto; padding: 14px 20px 0; }
.mrc-head h3 { margin: 0; color: var(--tx-0f172a); font-size: 1.375rem; font-weight: 800; }
.mrc-sub { display: flex; align-items: center; gap: 10px; margin-top: 8px; flex-wrap: wrap; }
.mrc-remaining { color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.mrc-pill { padding: 3px 10px; border-radius: 8px; font-size: 0.8125rem; font-weight: 600; }
.mrc-pill.is-pending { background: var(--bg-fef3c7); color: var(--tx-b45309); }
.mrc-pill.is-recorded { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.mrc-pill.is-abnormal { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.mrc-pill.is-skipped { background: var(--bg-f1f5f9); color: var(--tx-64748b); }

.mrc-body { flex: 1 1 auto; min-height: 0; overflow-y: auto; padding: 14px 20px; display: flex; flex-direction: column; gap: 14px; -webkit-overflow-scrolling: touch; }
.mrc-tenant { display: flex; flex-direction: column; padding: 10px 14px; border: 1px solid var(--bd-eef2f6); border-radius: 12px; background: var(--bg-f8fafc); }
.mrc-tenant small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.mrc-tenant strong { color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 700; }
.mrc-tenant span { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }

.mrc-meter { display: flex; flex-direction: column; gap: 10px; padding: 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; }
.mrc-meter.has-error { border-color: #fca5a5; }
.mrc-meter h4 { margin: 0; color: var(--tx-334155); font-size: 0.9375rem; font-weight: 700; }
.mrc-row { display: flex; align-items: baseline; justify-content: space-between; gap: 10px; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.mrc-row strong { color: var(--tx-0f172a); font-weight: 700; }
.mrc-row strong small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.mrc-field { display: flex; flex-direction: column; gap: 6px; }
.mrc-field > span { color: var(--tx-475569); font-size: 0.875rem; font-weight: 600; }
.mrc-field :deep(.app-input) { min-height: 48px; font-size: 1.25rem; font-weight: 700; }
.mrc-note { margin: 0; padding: 8px 10px; border-radius: 8px; font-size: 0.8125rem; font-weight: 500; line-height: 1.45; }
.mrc-note.is-warn { background: var(--bg-fef3c7); color: var(--tx-b45309); }

.mrc-total { display: flex; align-items: center; justify-content: space-between; padding: 12px 14px; border-radius: 12px; background: var(--bg-eff6ff); color: var(--tx-475569); font-size: 0.875rem; font-weight: 600; }
.mrc-total strong { color: var(--tx-2563eb); font-size: 1.25rem; font-weight: 800; }

.mrc-skip { display: flex; flex-direction: column; gap: 10px; padding: 14px; border: 1px dashed var(--bd-cbd5e1); border-radius: 12px; }
.mrc-skip-actions { display: flex; justify-content: flex-end; gap: 8px; }

.mrc-foot { flex: 0 0 auto; display: flex; gap: 10px; padding: 12px 20px; border-top: 1px solid var(--bd-e2e8f0); background: var(--bg-f8fafc); }
.mrc-save { flex: 1 1 auto; }

@media (max-width: 720px) {
  .mrc { top: auto; left: 10px; right: 10px; bottom: 10px; width: auto; max-height: 78%; }
}
.mrc-auto { padding: 10px 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; background: var(--bg-f8fafc); }
.mrc-code { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
</style>
