<template>
  <aside
    ref="cardRef"
    class="mrc"
    :class="{ 'is-dragging': sheetDragging }"
    role="dialog"
    :aria-label="`จดมิเตอร์ ${label}`"
    :style="sheetHeight !== null ? { '--mrc-sheet-h': `${sheetHeight}px` } : undefined"
  >
    <!-- มือถือ: ขีดสำหรับลากปรับความสูง แตะเพื่อสลับครึ่งจอ/เต็มจอ ลากลงต่ำสุดเพื่อปิด -->
    <div
      class="mrc-handle"
      role="button"
      tabindex="0"
      aria-label="ปรับความสูงหรือรูดลงเพื่อปิดการ์ดจดมิเตอร์"
      @touchstart="onSheetTouchStart"
      @touchmove="onSheetTouchMove"
      @touchend="onSheetTouchEnd"
    >
      <span class="mrc-handle-pill"></span>
    </div>

    <div class="mrc-top">
      <AppButton variant="primary" size="sm" @click="emit('back')">
        <template #prefix>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="15 18 9 12 15 6"></polyline></svg>
        </template>
        {{ projectName }}
      </AppButton>
      <div class="mrc-top-actions">
        <button v-if="bill" type="button" class="mrc-close" :class="{ 'is-active': tenantOpen }" aria-label="ดูข้อมูลผู้เช่า" title="ผู้เช่า" :aria-expanded="tenantOpen" @click="tenantOpen = !tenantOpen">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
        </button>
        <button type="button" class="mrc-close" aria-label="ปิดและกลับไปหน้ารวมโครงการ" title="ปิด" @click="emit('close')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
        </button>
      </div>

    </div>

    <header class="mrc-head">
      <div class="mrc-head-row">
        <h3>{{ label }}</h3>
        <span v-if="bill" class="mrc-pill" :class="`is-${state}`">{{ stateLabel[state] }}</span>
      </div>
      <p v-if="bill" class="mrc-head-sub"><template v-if="bill.shopName">{{ bill.shopName }} · </template>งวด{{ cycleLabel }}<template v-if="readOnly"> · ดูข้อมูลย้อนหลัง</template></p>
    </header>

    <div class="mrc-body">
      <p v-if="notice" class="mrc-notice" role="status">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
        {{ notice }}
      </p>
      <AppEmptyState v-if="!bill" class="mrc-empty" variant="compact" icon="users" title="ไม่มีผู้เช่า" description="แผงนี้ไม่ต้องจดมิเตอร์" />
      <AppEmptyState v-else-if="state === 'no-data'" class="mrc-empty" variant="compact" icon="history" title="ไม่มีข้อมูลงวดนี้" description="งวดย้อนหลังดูเลขมิเตอร์ที่จดไว้ได้อย่างเดียว" />
      <AppEmptyState v-else-if="state === 'no-meter'" class="mrc-empty" variant="compact" icon="alert" title="ยังไม่มีมิเตอร์" :description="`ยังไม่ได้ตั้ง${missingText} เพิ่มได้ที่หน้าจัดการมิเตอร์`" />
      <AppEmptyState v-else-if="!meters.length" class="mrc-empty" variant="compact" icon="document" title="ไม่ต้องจดมิเตอร์" description="แผงนี้คิดน้ำ-ไฟแบบเหมาจ่ายหรือไม่คิดค่าใช้บริการ" />
      <template v-else>
        <p v-if="readOnly" class="mrc-past-note">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          กำลังดูข้อมูลย้อนหลังของงวด{{ cycleLabel }} แก้ไขไม่ได้
        </p>
        <section v-for="meter in meters" :key="meter.key" class="mrc-meter" :class="{ 'has-error': !!meter.error }">
          <header class="mrc-meter-head">
            <span class="mrc-meter-icon" aria-hidden="true">
              <svg v-if="meter.key === 'water'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2.69l5.66 5.66a8 8 0 1 1-11.31 0z"></path></svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
            </span>
            <div class="mrc-meter-title">
              <h4>{{ meter.title }}</h4>
              <span v-if="meter.meterCode">{{ meter.meterCode }}</span>
            </div>
          </header>
          <div class="mrc-meter-body">
            <div class="mrc-meter-fields">
              <div class="mrc-prev"><span>ครั้งก่อน</span><strong>{{ meter.prev.toLocaleString('th-TH') }}</strong><em class="mrc-rate">฿{{ meter.rate.toLocaleString('th-TH') }}/{{ meter.unit }}</em></div>
              <label class="mrc-field" @keydown.enter.prevent="onEnter">
                <span>เลขครั้งนี้</span>
                <AppInput v-model="form[meter.key]" type="number" inputmode="decimal" placeholder="กรอกเลข" :disabled="readOnly" :error="meter.error" />
              </label>
              <!-- มือถือ: บอกหน่วยที่ใช้ไปใต้ช่องกรอก (เดสก์ท็อปดูจากแถบคำนวณด้านล่าง) -->
              <p class="mrc-used" :class="{ 'is-empty': meter.value == null || meter.value < 0 }">
                <span><strong>{{ meter.value != null && meter.value >= 0 ? meter.usedText : '-' }}</strong> {{ meter.unit }}</span>
                <b class="mrc-used-cost">฿{{ (meter.value != null && meter.value >= 0 && meter.charge != null ? meter.charge : 0).toLocaleString('th-TH') }}</b>
              </p>
            </div>
            <p v-if="meter.warning" class="mrc-note is-warn">{{ meter.warning }}</p>
          </div>
          <footer class="mrc-meter-foot">
            <div class="mrc-foot-calc">
              <small v-if="meter.value != null && meter.value >= 0" class="mrc-formula"><span>{{ meter.usedText }} {{ meter.unit }}</span><i aria-hidden="true">×</i><span>฿{{ meter.rate.toLocaleString('th-TH') }}</span></small>
              <small v-else>กรอกเลขครั้งนี้เพื่อคำนวณ</small>
            </div>
            <strong class="mrc-foot-amount" :class="{ 'is-empty': meter.value == null || meter.value < 0 || meter.charge == null }">฿{{ (meter.value != null && meter.value >= 0 && meter.charge != null ? meter.charge : 0).toLocaleString('th-TH') }}</strong>
          </footer>
        </section>

        <div v-for="auto in autoCharges" :key="auto.key" class="mrc-row mrc-auto"><span>{{ auto.title }}</span><strong>{{ auto.text }}</strong></div>

        <div class="mrc-total"><span>ยอดรวม</span><strong>฿{{ preview.toLocaleString('th-TH') }}</strong></div>
      </template>
    </div>

    <footer v-if="bill && !readOnly && state !== 'no-meter' && meters.length" class="mrc-foot">
      <AppButton variant="outline" size="lg" class="mrc-prev-btn" :disabled="!canGoBack" aria-label="ย้อนกลับไปแผงก่อนหน้า" title="แผงก่อนหน้า" @click="emit('prev')">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="15 18 9 12 15 6"></polyline></svg>
      </AppButton>
      <AppButton variant="outline" size="lg" @click="skipOpen = true">ข้าม</AppButton>
      <AppButton variant="primary" size="lg" class="mrc-save" @click="save">{{ saveLabel }}</AppButton>
    </footer>
    <AppModal v-model="skipOpen" title="ข้ามการจดมิเตอร์" :description="label" icon="document">
      <p class="mrc-reasons-hint">เลือกเหตุผล แล้วกดตกลง</p>
      <div class="mrc-reason-list" :class="{ 'has-error': !!skipError }" role="group" aria-label="เหตุผลที่ข้าม">
        <button v-for="option in skipOptions" :key="option.value" type="button" class="mrc-reason" :class="{ 'is-selected': skipReason === option.value }" :aria-pressed="skipReason === option.value" @click="selectSkipReason(option.value)">
          <span class="mrc-reason-icon" aria-hidden="true">
            <svg v-if="option.icon === 'tool'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94z"></path></svg>
            <svg v-else-if="option.icon === 'lock'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>
          </span>
          <span class="mrc-reason-label">{{ option.label }}</span>
          <span class="mrc-reason-check" aria-hidden="true">
            <svg v-if="skipReason === option.value" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
          </span>
        </button>
      </div>
      <p v-if="skipError" class="mrc-reason-error">{{ skipError }}</p>
      <template #footer>
        <AppButton variant="secondary" size="md" @click="skipOpen = false">ยกเลิก</AppButton>
        <AppButton variant="primary" size="md" @click="confirmSkip">ตกลง</AppButton>
      </template>
    </AppModal>

    <AppModal v-model="tenantOpen" title="ข้อมูลผู้เช่า" icon="users">
      <dl v-if="bill" class="mrc-tenant-detail">
        <div><dt>ผู้เช่าตามสัญญา</dt><dd>{{ bill.tenant }}</dd></div>
        <div v-if="bill.shopName"><dt>ร้านค้า</dt><dd>{{ bill.shopName }}</dd></div>
        <div v-if="bill.contractNumber"><dt>เลขที่สัญญา</dt><dd>{{ bill.contractNumber }}</dd></div>
      </dl>
      <template #footer>
        <AppButton variant="secondary" size="md" @click="tenantOpen = false">ปิด</AppButton>
      </template>
    </AppModal>
  </aside>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import AppButton from '~/components/AppButton.vue'
import AppModal from '~/components/AppModal.vue'
import AppEmptyState from '~/components/AppEmptyState.vue'
import AppInput from '~/components/AppInput.vue'
import { chargeFor, chargeRuleOf, useStallBilling, type ChargeKind } from '../../composables/useStallBilling'
import { useUtilitySettings } from '../../composables/useUtilitySettings'

// การ์ดจดมิเตอร์รายแผง เชื่อมรอบจดกับสัญญา ผู้เช่า มิเตอร์ และอัตราของแผงเดียวกัน
const props = defineProps<{
  code: string
  label: string
  projectName: string
  pendingLeft: number
  /** งวดย้อนหลัง: ดูข้อมูลที่จดไว้ได้อย่างเดียว */
  readOnly?: boolean
  /** มีแผงก่อนหน้าให้ย้อนกลับไปได้ */
  canGoBack?: boolean
  /** ข้อความแจ้งผลของแผงก่อนหน้า เช่น "บันทึกแผง 12 แล้ว" แสดงบนการ์ดของแผงถัดไป */
  notice?: string
  hasNext: boolean
  /** ค่าเฉลี่ยหน่วยที่ใช้ของแผงที่จดแล้วในโครงการ (0 = ยังไม่พอเทียบ) ใช้เตือนค่าที่สูงผิดปกติ */
  averageElectric: number
  averageWater: number
}>()
const emit = defineEmits<{ back: []; close: []; prev: []; saved: [code: string]; skipped: [code: string] }>()

const cardRef = ref<HTMLElement | null>(null)
const tenantOpen = ref(false)
const { getUtility, getUtilitySkip, skipUtility, saveUtilityReadings, billingMonth } = useStallBilling()
const cycleLabel = computed(() => new Intl.DateTimeFormat('th-TH', { month: 'long', year: 'numeric' }).format(new Date(`${billingMonth.value}-01T00:00:00`)))
const { getMeter } = useUtilitySettings(String(useRoute().query.market || '1'))
const ABNORMAL_FACTOR = 2.5

const bill = computed(() => getUtility(props.code))
const form = reactive({ water: '', electric: '' })
const submitted = ref(false)
const skipOpen = ref(false)
const skipError = ref('')
const skipReason = ref('')
const skipOptions = [
  { value: 'มิเตอร์ชำรุด', label: 'มิเตอร์ชำรุด', icon: 'tool' },
  { value: 'ร้านปิด เข้าไม่ได้', label: 'ร้านปิด เข้าไม่ได้', icon: 'lock' },
  { value: 'อ่านค่าไม่ได้', label: 'อ่านค่าไม่ได้', icon: 'eye-off' }
]

const reset = () => {
  submitted.value = false
  skipOpen.value = false
  skipError.value = ''
  skipReason.value = ''
  form.water = bill.value?.currWater == null ? '' : String(bill.value.currWater)
  form.electric = bill.value?.currElectric == null ? '' : String(bill.value.currElectric)
}
watch(() => props.code, () => { tenantOpen.value = false })
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
  return { key, meterCode: getMeter(props.code)[key], title: key === 'water' ? 'มิเตอร์น้ำ' : 'มิเตอร์ไฟ', unit: 'หน่วย', rate: ruleOf(key).rate + (ruleOf(key).surcharge || 0), ...st, usedText: st.value == null ? '-' : Math.max(0, st.value).toLocaleString('th-TH'), charge: chargeFor(ruleOf(key), st.value) }
}))
// ประเภทที่ต้องจดมิเตอร์แต่ยังไม่ได้ผูกเลขมิเตอร์ไว้ (แผงนี้จดไม่ได้จนกว่าจะเพิ่มมิเตอร์)
const missingKinds = computed(() => KINDS.filter(k => ruleOf(k).mode === 'meter' && !getMeter(props.code)[k]))
const missingText = computed(() => missingKinds.value.map(k => (k === 'water' ? 'มิเตอร์น้ำ' : 'มิเตอร์ไฟ')).join(' และ '))
const autoCharges = computed(() => KINDS.filter(k => ruleOf(k).mode !== 'meter').map(key => ({
  key,
  title: key === 'water' ? 'ค่าน้ำ' : 'ค่าไฟ',
  text: ruleOf(key).mode === 'flat' ? `เหมาจ่าย ฿${ruleOf(key).flat.toLocaleString('th-TH')}/เดือน (ไม่ต้องจดมิเตอร์)` : `${key === 'water' ? 'ไม่คิดค่าน้ำประปา' : 'ไม่คิดค่าไฟฟ้า'} (ไม่ต้องจดมิเตอร์)`
})))
const valid = computed(() => meters.value.every(m => form[m.key] !== '' && !(m.value! < 0)))
const preview = computed(() => KINDS.reduce((sum, k) => sum + (chargeFor(ruleOf(k), stateOf(k).value) ?? 0), 0))
const dirty = computed(() => !!bill.value && (form.water !== String(bill.value.currWater ?? '') || form.electric !== String(bill.value.currElectric ?? '')))
const recorded = computed(() => !!bill.value && bill.value.status !== 'unrecorded')

const state = computed<'recorded' | 'pending' | 'skipped' | 'no-meter' | 'no-data'>(() => {
  if (recorded.value) return 'recorded'
  if (props.readOnly) return 'no-data'
  if (missingKinds.value.length) return 'no-meter'
  return getUtilitySkip(props.code) ? 'skipped' : 'pending'
})
const stateLabel = { recorded: 'จดแล้ว', pending: 'ยังไม่จด', skipped: `ข้าม`, 'no-meter': 'ยังไม่มีมิเตอร์', 'no-data': 'ไม่มีข้อมูล' } as const

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
// เลือกเหตุผลก่อน แล้วกดตกลงเพื่อข้าม
const selectSkipReason = (value: string) => { skipReason.value = value; skipError.value = '' }
const confirmSkip = () => {
  if (!skipReason.value) { skipError.value = 'กรุณาเลือกเหตุผล'; return }
  skipUtility(props.code, skipReason.value)
  skipOpen.value = false
  emit('skipped', props.code)
}

// กด Enter ในช่องกรอก: ไปช่องถัดไป ถ้าเป็นช่องสุดท้ายให้บันทึกและไปแผงถัดไปเลย
const onEnter = (event: KeyboardEvent) => {
  const inputs = Array.from(cardRef.value?.querySelectorAll<HTMLInputElement>('input.app-input') || [])
  const index = inputs.indexOf(event.target as HTMLInputElement)
  const next = inputs[index + 1]
  if (next) next.focus()
  else if (bill.value && !props.readOnly && !skipOpen.value && state.value !== 'no-meter') save()
}

// ความสูงของการ์ดบนมือถือ (จำไว้ข้ามแผง เจ้าหน้าที่ปรับครั้งเดียวแล้วใช้ต่อได้) null = ครึ่งจอตามค่าเริ่มต้น
const sheetHeight = useState<number | null>('overview-sheet-height', () => null)
const sheetDragging = ref(false)
let sheetStartY = 0
let sheetStartHeight = 0
let sheetMoved = false
const SHEET_BREAKPOINT = 768
const sheetMax = () => window.innerHeight - 58
const onSheetTouchStart = (event: TouchEvent) => {
  if (window.innerWidth > SHEET_BREAKPOINT || !event.touches.length) return
  sheetStartY = event.touches[0].clientY
  sheetStartHeight = cardRef.value?.offsetHeight ?? window.innerHeight * 0.5
  sheetMoved = false
  sheetDragging.value = true
}
const onSheetTouchMove = (event: TouchEvent) => {
  if (!sheetDragging.value || !event.touches.length) return
  const delta = sheetStartY - event.touches[0].clientY
  if (Math.abs(delta) > 4) sheetMoved = true
  sheetHeight.value = Math.min(sheetMax(), Math.max(window.innerHeight * 0.18, sheetStartHeight + delta))
}
const onSheetTouchEnd = () => {
  if (!sheetDragging.value) return
  sheetDragging.value = false
  if (!sheetMoved) {
    const current = sheetHeight.value ?? cardRef.value?.offsetHeight ?? window.innerHeight * 0.5
    sheetHeight.value = current < window.innerHeight * 0.75 ? Math.round(sheetMax()) : Math.round(window.innerHeight * 0.5)
    return
  }
  // ลากลงต่ำกว่า 22% ของจอ = ปิดการ์ดกลับหน้ารวมโครงการ
  if (sheetHeight.value !== null && sheetHeight.value <= window.innerHeight * 0.22) {
    sheetHeight.value = null
    emit('close')
  }
}

// เปิดการ์ดแล้วโฟกัสช่องแรกทันที เจ้าหน้าที่พิมพ์ได้เลยโดยไม่ต้องแตะ
const focusFirstInput = () => {
  if (recorded.value) return
  nextTick(() => cardRef.value?.querySelector<HTMLInputElement>('input.app-input')?.focus())
}
const onDocumentKey = (event: KeyboardEvent) => { if (event.key === 'Escape') tenantOpen.value = false }
onMounted(() => {
  document.addEventListener('keydown', onDocumentKey)
  focusFirstInput()
})
onBeforeUnmount(() => {
  document.removeEventListener('keydown', onDocumentKey)
})
// เปลี่ยนแผงขณะการ์ดเปิดอยู่: เปลี่ยนแค่ข้อมูล ไม่เปิดการ์ดใหม่ (ไม่เล่นอนิเมชันสไลด์ซ้ำ) แล้วโฟกัสช่องแรกให้พิมพ์ต่อได้เลย
watch(() => props.code, focusFirstInput)
</script>

<style scoped>
.mrc {
  position: absolute;
  z-index: 60;
  top: auto;
  left: 18px;
  bottom: 18px;
  height: var(--meter-card-h, min(660px, calc(100% - 132px)));
  box-sizing: border-box;
  width: min(380px, calc(100% - 36px));
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 18px;
  background: var(--bg-ffffff);
  box-shadow: 0 16px 38px rgba(15, 23, 42, 0.16);
  animation: mrcIn .4s ease both;
}
@keyframes mrcIn { from { opacity: 0; } to { opacity: 1; } }
.mrc-top { flex: 0 0 auto; padding: 14px 20px 10px; border-bottom: 1px solid var(--bd-e8eef6);  display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.mrc-close { display: grid; place-items: center; flex: 0 0 auto; width: 32px; height: 32px; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-64748b); cursor: pointer; transition: background-color .15s ease, color .15s ease; }
.mrc-top { position: relative; }
.mrc-top-actions { display: flex; align-items: center; gap: 8px; }
.mrc-close.is-active { border-color: var(--bd-dbeafe); background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.mrc-tenant-detail { display: flex; flex-direction: column; margin: 0; }
.mrc-tenant-detail > div { display: flex; justify-content: space-between; gap: 16px; padding: 12px 0; border-bottom: 1px solid var(--bd-f1f5f9); }
.mrc-tenant-detail > div:last-child { border-bottom: 0; }
.mrc-tenant-detail dt { color: var(--tx-64748b); font-size: var(--font-size-sm); font-weight: 600; }
.mrc-tenant-detail dd { margin: 0; color: var(--tx-0f172a); font-size: var(--font-size-base); font-weight: 700; text-align: right; overflow-wrap: anywhere; }
.mrc-close:hover { background: var(--bg-f1f5f9); color: var(--tx-0f172a); }
.mrc-close:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 2px; }
.mrc-head { flex: 0 0 auto; padding: 14px 20px 0; }
.mrc-head h3 { margin: 0; color: var(--tx-0f172a); font-size: 1.375rem; font-weight: 800; }
.mrc-head-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.mrc-head-sub { margin: 4px 0 0; color: var(--tx-64748b); font-size: var(--font-size-sm); font-weight: 500; }
.mrc-pill { padding: 3px 10px; border-radius: 8px; font-size: var(--font-size-xs); font-weight: 600; }
.mrc-pill.is-pending, .mrc-pill.is-no-meter { background: var(--bg-fef3c7); color: var(--tx-b45309); }
.mrc-pill.is-recorded { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.mrc-pill.is-abnormal { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.mrc-pill.is-no-data { background: var(--bg-f1f5f9); color: var(--tx-64748b); }
.mrc-notice { display: flex; align-items: center; gap: 8px; margin: 0; padding: 8px 12px; border-radius: 10px; background: var(--bg-ecfdf5); color: var(--tx-15803d); font-size: var(--font-size-sm); font-weight: 600; }
.mrc-notice svg { flex: 0 0 auto; }
.mrc-past-note { display: flex; align-items: center; gap: 8px; margin: 0; padding: 8px 12px; border-radius: 10px; background: var(--bg-f1f5f9); color: var(--tx-475569); font-size: var(--font-size-xs); font-weight: 600; line-height: 1.45; }
.mrc-past-note svg { flex: 0 0 auto; color: var(--tx-64748b); }
.mrc-pill.is-skipped { background: var(--bg-f1f5f9); color: var(--tx-64748b); }

.mrc-body { flex: 1 1 auto; min-height: 0; overflow-y: auto; padding: 12px 20px; display: flex; flex-direction: column; gap: 10px; -webkit-overflow-scrolling: touch; }

.mrc-body > * { flex-shrink: 0; }
.mrc-meter { display: flex; flex-direction: column; overflow: hidden; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-ffffff); }
.mrc-meter-head { display: flex; align-items: center; gap: 12px; padding: 10px 14px; border-bottom: 1px solid var(--bd-e2e8f0); }
.mrc-meter-icon { display: grid; place-items: center; flex: 0 0 32px; width: 32px; height: 32px; border: 1px solid var(--bd-dbeafe); border-radius: 10px; background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.mrc-meter-icon svg { width: 18px; height: 18px; }
.mrc-meter-title { min-width: 0; display: flex; align-items: baseline; flex-wrap: wrap; gap: 2px 10px; }
.mrc-meter-title h4 { margin: 0; color: var(--tx-0f172a); font-size: var(--font-size-base); font-weight: 700; }
.mrc-meter-title span { color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.mrc-meter-body { display: flex; flex-direction: column; gap: 8px; padding: 10px 14px; }
.mrc-meter-fields { display: grid; grid-template-columns: 112px minmax(0, 1fr); gap: 10px; align-items: start; }
.mrc-meter-foot { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 10px 14px; border-top: 1px solid var(--bd-e2e8f0); background: var(--bg-f8fafc); }
.mrc-rate { display: none; color: var(--tx-64748b); font-size: var(--font-size-sm); font-style: normal; font-weight: 500; }
.mrc-used { display: none; margin: 0; color: var(--tx-64748b); font-size: var(--font-size-sm); font-weight: 500; white-space: nowrap; }
.mrc-used-cost { color: var(--tx-0f172a); font-size: var(--font-size-base); font-weight: 800; font-variant-numeric: tabular-nums; }
.mrc-used.is-empty .mrc-used-cost { color: var(--tx-94a3b8); }
.mrc-used strong { color: var(--tx-0f172a); font-weight: 700; font-variant-numeric: tabular-nums; }
.mrc-used.is-empty strong { color: var(--tx-94a3b8); }
.mrc-foot-calc { min-width: 0; display: flex; flex-direction: column; gap: 2px; }
.mrc-foot-calc span { color: var(--tx-475569); font-size: var(--font-size-sm); font-weight: 600; }
.mrc-foot-calc small { color: var(--tx-94a3b8); font-size: var(--font-size-xs); font-weight: 500; }
.mrc-formula { display: inline-flex; align-items: baseline; flex-wrap: wrap; gap: 2px 8px; }
.mrc-formula i { color: var(--tx-94a3b8); font-style: normal; }
.mrc-foot-amount { flex: 0 0 auto; color: var(--tx-334155); font-size: var(--font-size-md); font-weight: 700; font-variant-numeric: tabular-nums; line-height: 1.2; }
.mrc-foot-amount.is-empty { color: var(--tx-94a3b8); }
.mrc-meter.has-error { border-color: #fca5a5; }
.mrc-row { display: flex; align-items: baseline; justify-content: space-between; gap: 10px; color: var(--tx-64748b); font-size: var(--font-size-sm); font-weight: 500; }
.mrc-row strong { color: var(--tx-0f172a); font-weight: 700; }
.mrc-row strong small { color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 500; }
.mrc-prev { display: flex; flex-direction: column; gap: 6px; }
.mrc-prev > span { color: var(--tx-475569); font-size: var(--font-size-sm); font-weight: 600; }
.mrc-prev strong { display: flex; align-items: center; justify-content: center; min-height: 48px; border-radius: 8px; background: var(--bg-f1f5f9); color: var(--tx-0f172a); font-size: var(--font-size-xl); font-weight: 800; font-variant-numeric: tabular-nums; }
.mrc-field { display: flex; flex-direction: column; gap: 6px; }
.mrc-field > span { color: var(--tx-475569); font-size: var(--font-size-sm); font-weight: 600; }
.mrc-field :deep(.app-input) { min-height: 48px; font-size: var(--font-size-xl); font-weight: 700; }
/* กรอบประเหมือนสถานะว่างในแท็บรายละเอียดแผง */
.mrc-empty.app-empty-state { box-sizing: border-box; min-height: 0; margin-top: 4px; padding: 22px 16px; border: 1px dashed var(--bd-d5e0ed); border-radius: 16px; background: var(--bg-ffffff); }
.mrc-note { margin: 0; padding: 8px 10px; border-radius: 8px; font-size: var(--font-size-xs); font-weight: 500; line-height: 1.45; }
.mrc-note.is-warn { background: var(--bg-fef3c7); color: var(--tx-b45309); }

.mrc-total { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 10px 16px; border: 1px solid var(--bd-dbeafe); border-radius: 14px; background: var(--bg-eff6ff); color: var(--tx-334155); font-size: var(--font-size-base); font-weight: 700; }
.mrc-total strong { color: var(--tx-2563eb); font-size: 1.75rem; font-weight: 800; font-variant-numeric: tabular-nums; line-height: 1.15; }


/* ส่วนท้ายเหมือนการ์ดโครงการ: พื้นขาว ปุ่มสูง 46px (lg) ระยะล่าง 18px */
.mrc-foot { flex: 0 0 auto; display: flex; gap: 10px; padding: 12px 20px 18px; background: var(--bg-ffffff); }
.mrc-save { flex: 1 1 auto; }
.mrc-prev-btn { flex: 0 0 auto; padding-left: 12px; padding-right: 12px; }

.mrc-handle { display: none; flex: 0 0 auto; align-items: center; justify-content: center; width: 100%; padding: 12px 0 10px; border-bottom: 1px solid var(--bd-f1f5f9); background: var(--bg-ffffff); cursor: grab; touch-action: none; }
.mrc-handle-pill { display: block; width: 44px; height: 5px; border-radius: 999px; background: var(--bg-cbd5e1); }
@keyframes mrcUp { from { transform: translateY(100%); } to { transform: translateY(0); } }
/* มือถือ: แผ่นล่างกว้างเต็มด้านข้าง สูงครึ่งจอ ลากขีดด้านบนเพื่อขยายหรือรูดปิดได้ */
@media (max-width: 768px) {
  /* ขนาด ตำแหน่ง และเงาเดียวกับแผ่นล่างของการ์ดโครงการ เปิดต่อกันแล้วไม่กระโดด */
  .mrc { position: fixed; z-index: 200; top: auto; left: 0; right: 0; bottom: 0; width: auto; height: var(--mrc-sheet-h, 62dvh); max-height: calc(100dvh - 58px); border: 1px solid var(--bd-e2e8f0); border-bottom: 0; border-radius: 20px 20px 0 0; box-shadow: 0 -8px 32px rgba(15, 23, 42, 0.18); animation: mrcUp .5s cubic-bezier(0.22, 1, 0.36, 1) both; transition: height 0.22s cubic-bezier(0.16, 1, 0.3, 1); }
  .mrc.is-dragging { transition: none; }
  .mrc-handle { display: flex; }
  .mrc-handle { padding: 8px 0 6px; }
  .mrc-top { padding: 8px 16px; }

  /* มือถือ: น้ำและไฟอยู่ในแถวเดียวต่อมิเตอร์ เห็นและกรอกคู่กันได้โดยไม่ต้องเลื่อน */
  .mrc-head { padding: 8px 16px 0; }
  .mrc-head-sub { margin-top: 2px; font-size: var(--font-size-xs); }
  .mrc-head h3 { font-size: 1.125rem; }
  .mrc-body { padding: 12px 16px; gap: 14px; }
  /* การ์ดมิเตอร์: แถวบน (ชื่อ + ช่องกรอก) กับแถวล่าง (เลขก่อน + ใช้ไป) คั่นด้วยเส้นเต็มความกว้าง */
  .mrc-meter { display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1fr); grid-template-areas: 'head input' 'prev used'; align-items: center; }
  .mrc-meter-body, .mrc-meter-fields { display: contents; }
  .mrc-meter-head { grid-area: head; padding: 12px 8px 12px 14px; gap: 10px; border-bottom: 0; }
  .mrc-meter-icon { flex: 0 0 28px; width: 28px; height: 28px; border-radius: 9px; }
  .mrc-meter-icon svg { width: 16px; height: 16px; }
  .mrc-meter-title h4 { font-size: var(--font-size-sm); white-space: nowrap; }
  .mrc-meter-title span { display: none; }
  .mrc-field { grid-area: input; justify-self: end; width: min(calc(100% - 14px), 168px); margin: 12px 14px 12px 0; }
  .mrc-prev { grid-area: prev; align-self: stretch; flex-direction: column; align-items: flex-start; justify-content: center; gap: 2px; white-space: nowrap; padding: 8px 0 8px 14px; border-top: 1px solid var(--bd-e2e8f0); background: var(--bg-f8fafc); }
  .mrc-rate { display: block; }
  .mrc-used { grid-area: used; align-self: stretch; display: flex; flex-direction: column; align-items: flex-end; justify-content: center; gap: 2px; padding: 8px 14px 8px 0; border-top: 1px solid var(--bd-e2e8f0); background: var(--bg-f8fafc); text-align: right; }
  .mrc-meter .mrc-note { grid-column: 1 / -1; margin: 0 14px 10px; }
  .mrc-prev, .mrc-field { position: relative; }
  .mrc-prev > span, .mrc-field > span { position: absolute; width: 1px; height: 1px; overflow: hidden; clip: rect(0 0 0 0); white-space: nowrap; }
  /* เลขครั้งก่อนเป็นข้อความเล็กใต้ช่องกรอก ไม่ใช้กล่องสีเทา */
  .mrc-prev strong { display: inline; min-height: 0; border-radius: 0; background: transparent; color: var(--tx-334155); font-size: var(--font-size-sm); font-weight: 700; }
  .mrc-prev strong::before { content: 'ก่อน '; color: var(--tx-64748b); font-weight: 500; }
  .mrc-field :deep(.app-input) { min-height: 40px; font-size: 1.0625rem; }
  /* จนท.หน้างานต้องการแค่ดูเลขเก่าแล้วกรอกเลขใหม่: ซ่อนแถบคำนวณรายมิเตอร์ ดูแค่ยอดรวมด้านล่าง */
  .mrc-meter-foot { display: none; }
  .mrc-total { padding: 12px 16px; }
  .mrc-total strong { font-size: 1.375rem; }
}
.mrc-auto { padding: 10px 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; background: var(--bg-f8fafc); }
.mrc-reasons-hint { margin: 0 0 10px; color: var(--tx-64748b); font-size: var(--font-size-sm); font-weight: 500; }
.mrc-reason-list { display: flex; flex-direction: column; gap: 10px; }
.mrc-reason { display: flex; align-items: center; gap: 14px; width: 100%; min-height: 60px; padding: 10px 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-ffffff); color: var(--tx-0f172a); font: inherit; text-align: left; cursor: pointer; transition: border-color .15s ease, background-color .15s ease, transform .1s ease; }
.mrc-reason:hover { border-color: #93c5fd; background: var(--bg-eff6ff); }
.mrc-reason:active { transform: translateY(1px); }
.mrc-reason:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 2px; }
.mrc-reason-icon { display: grid; place-items: center; flex: 0 0 40px; width: 40px; height: 40px; border: 1px solid var(--bd-dbeafe); border-radius: 12px; background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.mrc-reason-icon svg { width: 20px; height: 20px; }
.mrc-reason-label { flex: 1 1 auto; min-width: 0; font-size: var(--font-size-base); font-weight: 700; }
.mrc-reason.is-selected { border-color: #2563eb; background: var(--bg-eff6ff); }
.mrc-reason-check { display: grid; place-items: center; flex: 0 0 22px; width: 22px; height: 22px; border: 2px solid var(--bd-cbd5e1); border-radius: 50%; background: var(--bg-ffffff); color: #ffffff; }
.mrc-reason.is-selected .mrc-reason-check { border-color: #2563eb; background: #2563eb; }
.mrc-reason-list.has-error .mrc-reason { border-color: #fca5a5; }
.mrc-reason-error { margin: 8px 0 0; color: var(--tx-dc2626); font-size: var(--font-size-xs); font-weight: 500; }
.mrc-code { color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 600; }
</style>
