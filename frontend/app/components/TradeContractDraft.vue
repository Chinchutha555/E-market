<template>
  <div>
    <AppHeader />
    <main class="content-body">
      <header class="heading">
        <div><h1>ตั้งค่าประเภทการค้า</h1><p>เพิ่มประเภทการค้าและกำหนดสัญญา</p></div>
        <div v-if="!loading" class="heading-actions">
          <AppButton size="md" variant="outline" @click="addTrade">
            <template #prefix>
              <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 5v14M5 12h14"/></svg>
            </template>
            เพิ่มประเภทการค้า
          </AppButton>
          <span class="save-state" :class="{ dirty: isDirty, saved: !isDirty && isSaved }">
            <i></i>{{ isDirty ? 'มีการแก้ไขที่ยังไม่ได้บันทึก' : isSaved ? 'บันทึกแล้ว' : 'ยังไม่ได้ตั้งค่า' }}
          </span>
          <AppButton size="md" :variant="isDirty ? 'primary' : 'outline'" :loading="saving" :disabled="!isDirty" @click="save">
            <template #prefix>
              <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><path d="M17 21v-8H7v8M7 3v5h8"/></svg>
            </template>
            บันทึกการตั้งค่า
          </AppButton>
        </div>
      </header>

      <AppPrimaryTabs
        v-if="trades.length > 1"
        v-model="activeTradeId"
        :options="tradeTabOptions"
        aria-label="เลือกประเภทการค้า"
      />
      <section class="panel" :class="{ 'primary-tab-content-card': trades.length > 1 }">
        <p v-if="loading" class="empty">กำลังโหลดข้อมูล...</p>
        <p v-else-if="!trades.length" class="empty">ยังไม่มีประเภทการค้า กดเพิ่มประเภทการค้าเพื่อเริ่มต้น</p>
        <div v-else-if="activeTrade" :key="activeTrade.id" class="trade">
          <div class="row"><h3>{{ activeTrade.title.trim() || 'ประเภทการค้า ' + (activeTradeIndex + 1) }}</h3><AppButton size="sm" variant="danger-outline" @click="removeTrade(activeTrade)">ลบประเภทการค้า</AppButton></div>
          <label :for="'trade-' + activeTrade.id">ชื่อประเภทการค้า <span class="required" aria-hidden="true">*</span></label>
          <AppInput :id="'trade-' + activeTrade.id" v-model="activeTrade.title" placeholder="ระบุชื่อประเภทการค้า" maxlength="255" />
          <div class="checkbox"><AppCheckbox v-model="activeTrade.requiresContract" label="ประเภทการค้านี้มีสัญญา" /></div>
          <template v-if="activeTrade.requiresContract">
            <div class="row section-heading"><div><h4>สัญญา</h4><p>เพิ่มสัญญาและระยะเวลาสำหรับประเภทการค้านี้</p></div><AppButton size="sm" variant="outline" @click="addContract(activeTrade)">เพิ่มสัญญา</AppButton></div>
            <p v-if="!activeTrade.contracts.length" class="empty">ยังไม่มีสัญญา</p>
            <section v-for="(contract, contractIndex) in activeTrade.contracts" :key="contract.id" class="contract">
              <div class="row"><h4>{{ contract.title || 'สัญญา ' + (contractIndex + 1) }}</h4><AppButton size="sm" variant="danger-outline" @click="removeContract(activeTrade, contract)">ลบสัญญา</AppButton></div>
              <label :for="'contract-' + contract.id">ชื่อสัญญา <span class="required" aria-hidden="true">*</span></label>
              <AppInput :id="'contract-' + contract.id" v-model="contract.title" placeholder="ระบุชื่อสัญญา" maxlength="255" />
              <div class="file"><AppFileUpload v-model="contract.file" label="ไฟล์สัญญา PDF" required-text="จำเป็น" accept=".pdf,application/pdf" :max-size-mb="10" hint="ไฟล์ PDF ขนาดไม่เกิน 10 MB" /></div>
              <div class="row section-heading"><div><h4>รูปแบบระยะเวลา</h4><p>เพิ่มได้หลายระยะเวลา</p></div><AppButton size="sm" variant="outline" @click="addDuration(contract.durations)">เพิ่มระยะเวลา</AppButton></div>
              <section v-for="(duration, index) in contract.durations" :key="duration.id" class="option-card">
                <div class="option-heading">
                  <strong>ระยะเวลาที่ {{ index + 1 }}</strong>
                  <AppButton v-if="contract.durations.length > 1" size="sm" variant="danger-outline" @click="removeDuration(contract.durations, index)">ลบ</AppButton>
                </div>
                <div class="fields">
                  <div class="field">
                    <label :for="'duration-' + duration.id">จำนวน <span class="required" aria-hidden="true">*</span></label>
                    <AppInput :id="'duration-' + duration.id" v-model="duration.value" type="number" min="1" max="10000" step="1" placeholder="ระบุจำนวน" @update:model-value="updateEndDate(duration)" />
                  </div>
                  <div class="field">
                    <label :id="'unit-' + duration.id">หน่วย <span class="required" aria-hidden="true">*</span></label>
                    <AppSelect v-model="duration.unit" :options="units" :aria-labelledby="'unit-' + duration.id" @update:model-value="updateEndDate(duration)" />
                  </div>
                  <div class="field">
                    <label :for="'start-' + duration.id">วันเริ่มรอบสัญญา</label>
                    <AppInput :id="'start-' + duration.id" v-model="duration.startDate" type="date" @update:model-value="updateEndDate(duration)" />
                  </div>
                  <div class="field">
                    <label :for="'end-' + duration.id">วันสิ้นสุดรอบสัญญา</label>
                    <AppInput :id="'end-' + duration.id" :model-value="duration.endDate" type="date" disabled />
                  </div>
                  <div v-if="duration.startDate && duration.endDate" class="term-result field-full" :class="{ 'is-invalid': !durationBetween(duration.startDate, duration.endDate) }">
                    {{ durationBetween(duration.startDate, duration.endDate) || 'วันสิ้นสุดต้องไม่อยู่ก่อนวันเริ่มต้น' }}
                    <span class="remaining-days">{{ remainingText(duration.endDate) }}</span>
                  </div>
                </div>
              </section>
            </section>
          </template>
          <template v-else>
            <div class="row section-heading"><div><h4>รูปแบบระยะเวลา</h4><p>เพิ่มได้หลายระยะเวลา</p></div><AppButton size="sm" variant="outline" @click="addDuration(activeTrade.durations)">เพิ่มระยะเวลา</AppButton></div>
            <section v-for="(duration, index) in activeTrade.durations" :key="duration.id" class="option-card">
              <div class="option-heading">
                <strong>ระยะเวลาที่ {{ index + 1 }}</strong>
                <AppButton v-if="activeTrade.durations.length > 1" size="sm" variant="danger-outline" @click="removeDuration(activeTrade.durations, index)">ลบ</AppButton>
              </div>
              <div class="fields">
                <div class="field">
                  <label :for="'duration-' + duration.id">จำนวน <span class="required" aria-hidden="true">*</span></label>
                  <AppInput :id="'duration-' + duration.id" v-model="duration.value" type="number" min="1" max="10000" step="1" placeholder="ระบุจำนวน" @update:model-value="updateEndDate(duration)" />
                </div>
                <div class="field">
                  <label :id="'unit-' + duration.id">หน่วย <span class="required" aria-hidden="true">*</span></label>
                  <AppSelect v-model="duration.unit" :options="units" :aria-labelledby="'unit-' + duration.id" @update:model-value="updateEndDate(duration)" />
                </div>
                <div class="field">
                  <label :for="'start-' + duration.id">วันเริ่มรอบสัญญา</label>
                  <AppInput :id="'start-' + duration.id" v-model="duration.startDate" type="date" @update:model-value="updateEndDate(duration)" />
                </div>
                <div class="field">
                  <label :for="'end-' + duration.id">วันสิ้นสุดรอบสัญญา</label>
                  <AppInput :id="'end-' + duration.id" :model-value="duration.endDate" type="date" disabled />
                </div>
                <div v-if="duration.startDate && duration.endDate" class="term-result field-full" :class="{ 'is-invalid': !durationBetween(duration.startDate, duration.endDate) }">
                  {{ durationBetween(duration.startDate, duration.endDate) || 'วันสิ้นสุดต้องไม่อยู่ก่อนวันเริ่มต้น' }}
                  <span class="remaining-days">{{ remainingText(duration.endDate) }}</span>
                </div>
              </div>
            </section>
          </template>
        </div>
        <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import type { AppPrimaryTabOption } from '~/components/AppPrimaryTabs.vue'
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'
import { bangkokToday, calculateTradeEndDate, tradeDaysRemaining } from '~/utils/tradeDates'

type Duration = { id: string; value: string | number; unit: string; startDate: string; endDate: string }
type Contract = { id: string; title: string; file: File | null; fileName?: string; durations: Duration[] }
type Trade = { id: string; title: string; requiresContract: boolean; contracts: Contract[]; durations: Duration[] }

const route = useRoute()
const authStore = useAuthStore()
const { get, put } = useApi()
const { showConfirm, showSuccess } = useAppAlert()

const market = computed(() => String(route.query.market || authStore.currentUser?.assigned_markets?.[0]?.id || '1'))
const legacyLabels: Record<string, string> = { trial: 'ทดลองขาย', temporary: 'ทำการค้าแบบชั่วคราว', contract: 'ทำการค้าแบบมีสัญญา' }
let nextId = 0
const genCode = (prefix: string) => prefix + (++nextId).toString(36) + Date.now().toString(36).slice(-6)
const units = [{ value: 'days', label: 'วัน' }, { value: 'months', label: 'เดือน' }, { value: 'years', label: 'ปี' }]

const trades = ref<Trade[]>([])
const activeTradeId = ref('')
const savedSnapshot = ref('[]')
const isSaved = ref(false)
const loading = ref(true)
const saving = ref(false)
const error = ref('')

const serializeTrades = (list: Trade[]) =>
  JSON.stringify(
    list.map(t => ({
      id: t.id,
      title: t.title,
      requiresContract: t.requiresContract,
      contracts: t.contracts.map(c => ({
        id: c.id,
        title: c.title,
        fileName: c.file?.name || c.fileName || '',
        durations: c.durations
      })),
      durations: t.durations
    }))
  )

const isDirty = computed(() => savedSnapshot.value !== serializeTrades(trades.value))

const tradeTabOptions = computed<AppPrimaryTabOption[]>(() =>
  trades.value.map((trade, index) => ({
    value: String(trade.id),
    label: trade.title.trim() || 'ประเภทการค้า ' + (index + 1),
    icon: 'contract'
  }))
)
const activeTradeIndex = computed(() => {
  const index = trades.value.findIndex(item => String(item.id) === activeTradeId.value)
  return index >= 0 ? index : 0
})
const activeTrade = computed(() => trades.value[activeTradeIndex.value] || null)

const newDuration = (): Duration => ({ id: genCode('o'), value: '', unit: 'months', startDate: '', endDate: '' })

const updateEndDate = (duration: Duration) => {
  duration.endDate = calculateTradeEndDate(duration.startDate, duration.value, duration.unit)
}

const remainingText = (end: string): string => {
  const days = tradeDaysRemaining(end, bangkokToday())
  if (days === null) return ''
  if (days === 0) return 'สิ้นสุดวันนี้'
  return days > 0 ? 'เหลือ ' + days + ' วัน' : 'สิ้นสุดแล้ว ' + Math.abs(days) + ' วัน'
}

const durationBetween = (startText: string, endText: string): string => {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(startText) || !/^\d{4}-\d{2}-\d{2}$/.test(endText)) return ''
  const start = new Date(startText + 'T00:00:00Z')
  const end = new Date(endText + 'T00:00:00Z')
  if (!Number.isFinite(start.getTime()) || !Number.isFinite(end.getTime()) || end < start || start.toISOString().slice(0, 10) !== startText || end.toISOString().slice(0, 10) !== endText) return ''
  end.setUTCDate(end.getUTCDate() + 1)
  const addMonths = (count: number) => {
    const target = new Date(Date.UTC(start.getUTCFullYear(), start.getUTCMonth() + count, 1))
    const lastDay = new Date(Date.UTC(target.getUTCFullYear(), target.getUTCMonth() + 1, 0)).getUTCDate()
    target.setUTCDate(Math.min(start.getUTCDate(), lastDay))
    return target
  }
  let months = (end.getUTCFullYear() - start.getUTCFullYear()) * 12 + end.getUTCMonth() - start.getUTCMonth()
  if (addMonths(months) > end) months--
  const days = Math.round((end.getTime() - addMonths(months).getTime()) / 86400000)
  return 'ระยะเวลารวม ' + Math.floor(months / 12) + ' ปี ' + (months % 12) + ' เดือน ' + days + ' วัน'
}

const addTrade = () => {
  const trade: Trade = { id: genCode('t'), title: '', requiresContract: false, contracts: [], durations: [newDuration()] }
  trades.value.push(trade)
  activeTradeId.value = String(trade.id)
}
const addContract = (trade: Trade) => trade.contracts.push({ id: genCode('c'), title: '', file: null, durations: [newDuration()] })
const addDuration = (durations: Duration[]) => durations.push(newDuration())
const removeDuration = async (durations: Duration[], index: number) => {
  if (durations.length <= 1) return
  if (await showConfirm('ต้องการลบระยะเวลานี้หรือไม่ การเปลี่ยนแปลงจะมีผลเมื่อบันทึก', { title: 'ลบระยะเวลา', type: 'danger', icon: 'trash', confirmText: 'ลบ' })) {
    durations.splice(index, 1)
  }
}
const removeTrade = async (trade: Trade) => {
  if (await showConfirm('ต้องการลบประเภทการค้านี้หรือไม่', { title: 'ลบประเภทการค้า', type: 'danger', icon: 'trash', confirmText: 'ลบ' })) {
    const removedIndex = trades.value.findIndex(item => item.id === trade.id)
    trades.value = trades.value.filter(item => item.id !== trade.id)
    if (activeTradeId.value === String(trade.id)) {
      const nextTrade = trades.value[Math.min(removedIndex, trades.value.length - 1)]
      activeTradeId.value = nextTrade ? String(nextTrade.id) : ''
    }
  }
}
const removeContract = async (trade: Trade, contract: Contract) => {
  if (await showConfirm('ต้องการลบสัญญานี้หรือไม่', { title: 'ลบสัญญา', type: 'danger', icon: 'trash', confirmText: 'ลบ' })) {
    trade.contracts = trade.contracts.filter(item => item.id !== contract.id)
  }
}

const storageKey = computed(() => 'emarket_trade_contracts_' + market.value)

const load = async () => {
  loading.value = true
  error.value = ''
  try {
    let savedContractMap: Record<string, Array<{ id: string; title: string; fileName?: string; durations: Duration[] }>> = {}
    if (import.meta.client) {
      try {
        const raw = localStorage.getItem(storageKey.value)
        if (raw) savedContractMap = JSON.parse(raw) || {}
      } catch {}
    }
    const response: any = await get('/markets/' + encodeURIComponent(market.value) + '/trade-conditions')
    const rows = Array.isArray(response?.data) ? response.data : []
    const mapped: Trade[] = rows.map((row: any) => {
      const tradeId = String(row.TradeType || genCode('t'))
      const rawOptions = Array.isArray(row.Options) && row.Options.length
        ? row.Options
        : [{ ID: genCode('o'), Duration: row.MaxDuration || 1, Unit: row.Unit || 'months', StartDate: row.StartDate || '', EndDate: row.EndDate || '' }]
      const durations: Duration[] = rawOptions.map((opt: any) => {
        const d: Duration = {
          id: String(opt.ID || genCode('o')),
          value: opt.Duration ?? '',
          unit: opt.Unit || 'months',
          startDate: opt.StartDate || '',
          endDate: opt.EndDate || ''
        }
        const corrected = calculateTradeEndDate(d.startDate, d.value, d.unit)
        if (corrected && d.endDate) {
          const prev = new Date(corrected + 'T00:00:00Z')
          prev.setUTCDate(prev.getUTCDate() + 1)
          if (prev.toISOString().slice(0, 10) === d.endDate) d.endDate = corrected
        }
        return d
      })
      const requiresContract = row.RequiresContract ?? false
      const localContracts = savedContractMap[tradeId]
      const contracts: Contract[] = requiresContract
        ? (Array.isArray(localContracts) && localContracts.length
            ? localContracts.map(c => ({
                id: String(c.id || genCode('c')),
                title: c.title || '',
                file: null,
                fileName: c.fileName || '',
                durations: Array.isArray(c.durations) && c.durations.length ? c.durations.map(d => ({ ...d })) : durations.map(d => ({ ...d }))
              }))
            : [{ id: genCode('c'), title: 'สัญญาหลัก', file: null, durations: durations.map(d => ({ ...d })) }])
        : []
      return {
        id: tradeId,
        title: row.Title || legacyLabels[row.TradeType] || row.TradeType || '',
        requiresContract,
        contracts,
        durations: durations.map(d => ({ ...d }))
      }
    })
    trades.value = mapped
    activeTradeId.value = trades.value[0] ? String(trades.value[0].id) : ''
    savedSnapshot.value = serializeTrades(trades.value)
    isSaved.value = trades.value.length > 0
  } catch {
    trades.value = []
    savedSnapshot.value = serializeTrades(trades.value)
    isSaved.value = false
  } finally {
    loading.value = false
  }
}

const save = async () => {
  if (!isDirty.value || saving.value) return
  error.value = ''
  const titles = trades.value.map(item => item.title.trim().toLowerCase())
  if (trades.value.some(item => !item.title.trim() || [...item.title.trim()].length > 255) || new Set(titles).size !== titles.length) {
    error.value = 'กรุณาระบุชื่อประเภทการค้าไม่ซ้ำกัน และไม่เกิน 255 ตัวอักษร'
    return
  }
  for (const item of trades.value) {
    const activeDurations = item.requiresContract
      ? item.contracts.flatMap(c => c.durations)
      : item.durations
    if (item.requiresContract && (!item.contracts.length || item.contracts.some(c => !c.title.trim()))) {
      error.value = 'กรุณาเพิ่มสัญญาและระบุชื่อสัญญาให้ครบถ้วน'
      return
    }
    if (!activeDurations.length || activeDurations.some(d => !Number.isInteger(Number(d.value)) || Number(d.value) < 1 || Number(d.value) > 10000 || !['days', 'months', 'years'].includes(d.unit))) {
      error.value = 'กรุณาระบุจำนวนระยะเวลาเป็นเลขจำนวนเต็มตั้งแต่ 1 ถึง 10,000'
      return
    }
    if (activeDurations.some(d => (d.startDate || d.endDate) && !durationBetween(d.startDate, d.endDate))) {
      error.value = 'กรุณาเลือกวันเริ่มและวันสิ้นสุดให้ถูกต้อง'
      return
    }
  }

  const payloadConditions = trades.value.map(item => {
    const sourceDurations = item.requiresContract
      ? item.contracts.flatMap(c => c.durations)
      : item.durations
    return {
      TradeType: /^[a-z0-9_-]{1,20}$/.test(item.id) ? item.id : genCode('t').slice(0, 20),
      Title: item.title.trim(),
      RequiresContract: item.requiresContract,
      Options: sourceDurations.map(d => ({
        ID: /^[a-z0-9_-]{1,20}$/.test(d.id) ? d.id : genCode('o').slice(0, 20),
        Duration: Number(d.value),
        Unit: d.unit,
        StartDate: d.startDate || '',
        EndDate: d.endDate || ''
      }))
    }
  })

  saving.value = true
  try {
    await put('/markets/' + encodeURIComponent(market.value) + '/trade-conditions', { Conditions: payloadConditions })
    if (import.meta.client) {
      const contractMap: Record<string, any> = {}
      trades.value.forEach((item, idx) => {
        const key = payloadConditions[idx]?.TradeType || item.id
        item.id = key
        contractMap[key] = item.contracts.map(c => ({
          id: c.id,
          title: c.title,
          fileName: c.file?.name || c.fileName || '',
          durations: c.durations
        }))
      })
      try { localStorage.setItem(storageKey.value, JSON.stringify(contractMap)) } catch {}
    }
    savedSnapshot.value = serializeTrades(trades.value)
    isSaved.value = trades.value.length > 0
    await showSuccess('บันทึกการตั้งค่าแล้ว', 'สำเร็จ')
  } catch {
    error.value = 'ไม่สามารถบันทึกการตั้งค่าได้ กรุณาลองใหม่'
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  load()
  watch(market, load)
})
</script>

<style scoped>
.heading,.row,.heading-actions,.option-heading{display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap}.heading{margin:8px 0 18px}.heading h1{margin:0;color:#0f172a;font-size:var(--font-size-2xl);font-weight:700}.heading p,.section-heading p{margin:4px 0 0;color:#64748b;font-size:var(--font-size-sm)}.heading-actions{gap:10px}.heading-actions svg{width:16px;height:16px;fill:none;stroke:currentColor;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}.save-state{display:inline-flex;align-items:center;gap:7px;padding:8px 12px;border:1px solid var(--bd-e2e8f0);border-radius:999px;color:var(--tx-64748b);font-size:var(--font-size-xs);font-weight:600}.save-state i{width:8px;height:8px;border-radius:50%;background:currentColor}.save-state.dirty{border-color:var(--bd-fed7aa);background:var(--bg-fff7ed);color:var(--tx-c2410c)}.save-state.saved{border-color:var(--bd-bbf7d0);background:var(--bg-ecfdf5);color:var(--tx-15803d)}.panel{padding:24px;border:1px solid var(--bd-dbe5f1,#dbe5f1);border-radius:16px;background:#fff;box-shadow:0 1px 3px rgba(0,0,0,.03)}.panel.primary-tab-content-card{border-top-left-radius:0}.row h3,.row h4{margin:0;color:#1e293b}.row h3{font-size:var(--font-size-lg)}.row h4{font-size:var(--font-size-md)}.empty{margin:0;color:#64748b;font-size:var(--font-size-sm)}.trade > .row:first-child,.contract > .row:first-child{margin-bottom:16px}.trade label,.contract label,.field label{display:block;margin:0 0 8px;color:#334155;font-size:var(--font-size-sm);font-weight:600}.required,.form-error{color:#ef4444}.form-error{margin:14px 0 0;font-size:var(--font-size-sm)}.checkbox{margin-top:16px}.section-heading{margin:20px 0 12px;padding-top:16px;border-top:1px solid #e2e8f0}.contract{margin-top:12px;padding:16px;border:1px solid #e2e8f0;border-radius:10px;background:#f8fafc}.file{margin-top:16px}.option-card{margin-top:12px;padding:14px;border:1px solid #e2e8f0;border-radius:10px;background:#fff}.option-heading strong{font-size:var(--font-size-sm);color:#334155}.fields{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:14px;margin-top:12px}.field{min-width:0}.field-full{grid-column:1/-1}.term-result{padding:10px 12px;border:1px solid #e2e8f0;border-radius:8px;background:#f8fafc;color:#334155;font-size:var(--font-size-sm);font-weight:600}.remaining-days{display:block;margin-top:4px;font-weight:500;color:#64748b}.term-result.is-invalid{color:#dc2626}@media(max-width:600px){.heading-actions{width:100%}.panel{padding:16px}.fields{grid-template-columns:1fr}}
</style>
