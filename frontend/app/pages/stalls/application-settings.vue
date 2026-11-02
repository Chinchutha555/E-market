<template>
  <div>
    <AppHeader />
    <main class="content-body">
      <header class="page-heading">
        <div><h1 class="page-title">จัดการคำร้อง</h1><p class="page-desc">กำหนดระยะเวลาของประเภทการค้า</p></div>
        <div v-if="!loading && !loadError" class="heading-actions">
          <span class="draft-status" :class="{ published: !isDirty && isConfigured, 'is-dirty': isDirty }"><i></i>{{ isDirty ? 'มีการเปลี่ยนแปลงที่ยังไม่ได้บันทึก' : isConfigured ? 'บันทึกแล้ว' : 'ยังไม่ได้ตั้งค่า' }}</span>
          <AppButton size="md" :variant="isDirty ? 'primary' : 'outline'" :loading="saving" @click="save">
            <template #prefix><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path><polyline points="17 21 17 13 7 13 7 21"></polyline><polyline points="7 3 7 8 15 8"></polyline></svg></template>
            บันทึกการตั้งค่า
          </AppButton>
        </div>
      </header>
      <div v-if="loading" class="settings-card notice" role="status">กำลังโหลดข้อมูล...</div>
      <div v-else-if="loadError" class="settings-card error" role="alert">{{ loadError }}</div>
      <section v-else class="settings-card" aria-labelledby="conditions-title">
        <div class="section-heading">
          <div>
            <h2 id="conditions-title" class="section-title">ระยะเวลาตามประเภทการค้า</h2>
            <p class="section-subtitle">{{ conditions.length }} ประเภทการค้า</p>
          </div>
          <AppButton size="md" variant="outline" :disabled="conditions.length >= 100" @click="addCondition">
            <template #prefix><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 5v14M5 12h14"></path></svg></template>
            เพิ่มประเภทการค้า
          </AppButton>
        </div>
        <p v-if="!conditions.length" class="empty-state">ยังไม่มีประเภทการค้า กดเพิ่มประเภทการค้าเพื่อเริ่มตั้งค่า</p>
        <div class="conditions-grid">
          <section v-for="item in conditions" :key="item.TradeType" class="condition">
            <header class="condition-heading">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="5" width="18" height="16" rx="2"></rect><path d="M16 3v4M8 3v4M3 10h18"></path></svg>
              <h3>{{ item.Title.trim() || 'ประเภทการค้าใหม่' }}</h3>
              <AppButton size="sm" variant="danger-outline" @click="removeCondition(item)">ลบ</AppButton>
            </header>
            <div class="fields">
              <div class="field field-full"><label :for="item.TradeType + '-title'">ชื่อประเภทการค้า <span class="req-star" aria-hidden="true">*</span></label><AppInput :id="item.TradeType + '-title'" v-model="item.Title" maxlength="255" placeholder="ระบุชื่อประเภทการค้า" required /></div>
              <div class="field field-full contract-requirement">
                <AppCheckbox v-model="item.RequiresContract" label="ต้องทำสัญญาหลังอนุมัติคำร้อง" />
                <span>{{ item.RequiresContract ? 'ต้องทำสัญญา' : 'ไม่ต้องทำสัญญา' }}</span>
              </div>
              <div class="field-full options-heading">
                <strong>รูปแบบระยะเวลา</strong>
                <AppButton size="sm" variant="outline" :disabled="item.Options.length >= 100" @click="addOption(item)">เพิ่มระยะเวลา</AppButton>
              </div>
              <section v-for="(option, index) in item.Options" :key="option.ID" class="option-card field-full">
                <div class="option-heading">
                  <strong>ระยะเวลาที่ {{ index + 1 }}</strong>
                  <AppButton v-if="item.Options.length > 1" size="sm" variant="danger-outline" @click="removeOption(item, option)">ลบ</AppButton>
                </div>
                <div class="fields">
                  <div class="field"><label :for="item.TradeType + option.ID + '-duration'">จำนวน <span class="req-star" aria-hidden="true">*</span></label><AppInput :id="item.TradeType + option.ID + '-duration'" v-model="option.Duration" type="number" min="1" max="10000" step="1" required @update:model-value="updateEndDate(option)" /></div>
                  <div class="field"><label :id="item.TradeType + option.ID + '-unit'">หน่วย <span class="req-star" aria-hidden="true">*</span></label><AppSelect v-model="option.Unit" :options="unitOptions" :aria-labelledby="item.TradeType + option.ID + '-unit'" @update:model-value="updateEndDate(option)" /></div>
                  <div class="field"><label :for="item.TradeType + option.ID + '-start'">วันเริ่มรอบสัญญา</label><AppInput :id="item.TradeType + option.ID + '-start'" v-model="option.StartDate" type="date" @update:model-value="updateEndDate(option)" /></div>
                  <div class="field"><label :for="item.TradeType + option.ID + '-end'">วันสิ้นสุดรอบสัญญา</label><AppInput :id="item.TradeType + option.ID + '-end'" :model-value="option.EndDate" type="date" disabled /></div>
                  <div v-if="option.StartDate && option.EndDate" class="term-result field-full" :class="{ 'is-invalid': !durationBetween(option.StartDate, option.EndDate) }">
                    {{ durationBetween(option.StartDate, option.EndDate) || 'วันสิ้นสุดต้องไม่อยู่ก่อนวันเริ่มต้น' }}
                    <span class="remaining-days">{{ remainingText(option.EndDate) }}</span>
                  </div>
                </div>
              </section>
            </div>
          </section>
        </div>
        <p v-if="error" class="error" role="alert">{{ error }}</p>
      </section>
    </main>
  </div>
</template>
<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'
import { bangkokToday, calculateTradeEndDate, tradeDaysRemaining } from '~/utils/tradeDates'
definePageMeta({
  middleware: (to) => {
    if (!import.meta.client) return
    const auth = useAuthStore()
    auth.hydrate()
    if (!auth.isLoggedIn) return navigateTo('/auth/login')
    if (auth.role === 'staff') {
      const assigned = auth.currentUser?.assigned_markets?.some((m: any) => String(m.id) === String(to.query.market || ''))
      if (!assigned) return navigateTo('/staff', { replace: true })
    }
  }
})
type DurationOption = { ID: string; Duration: string | number; Unit: string; StartDate: string; EndDate: string }
type Condition = { TradeType: string; Title: string; RequiresContract: boolean; Options: DurationOption[] }
const route = useRoute()
const { get, put } = useApi()
const { showSuccess, showConfirm } = useAppAlert()
const legacyLabels: Record<string, string> = { trial: 'ทดลองขาย', temporary: 'ทำการค้าแบบชั่วคราว', contract: 'ทำการค้าแบบมีสัญญา' }
const unitOptions = [{ value: 'days', label: 'วัน' }, { value: 'months', label: 'เดือน' }, { value: 'years', label: 'ปี' }]
const conditions = ref<Condition[]>([])
const savedConditions = ref('')
const isConfigured = ref(false)
const isDirty = computed(() => savedConditions.value !== JSON.stringify(conditions.value))
const loading = ref(true)
const saving = ref(false)
const error = ref('')
const loadError = ref('')
const market = computed(() => String(route.query.market || ''))
const load = async () => {
  loading.value = true
  loadError.value = ''
  if (!market.value) { loadError.value = 'ไม่พบตลาดที่เลือก'; loading.value = false; return }
  try {
    const response: any = await get('/markets/' + encodeURIComponent(market.value) + '/trade-conditions')
    const rows = Array.isArray(response.data) ? response.data : []
    conditions.value = rows.map((row: any) => ({
      TradeType: row.TradeType,
      Title: row.Title || legacyLabels[row.TradeType] || row.TradeType,
      RequiresContract: row.RequiresContract ?? true,
      Options: (Array.isArray(row.Options) && row.Options.length ? row.Options : [{ ID: 'legacy', Duration: row.MaxDuration, Unit: row.Unit, StartDate: row.StartDate, EndDate: row.EndDate }]).map((option: any) => ({
        ID: option.ID, Duration: option.Duration, Unit: option.Unit, StartDate: option.StartDate || '', EndDate: option.EndDate || ''
      }))
    }))
    savedConditions.value = JSON.stringify(conditions.value)
    for (const item of conditions.value) for (const option of item.Options) {
      const corrected = calculateTradeEndDate(option.StartDate, option.Duration, option.Unit)
      if (!corrected || !option.EndDate) continue
      const previous = new Date(corrected + 'T00:00:00Z')
      previous.setUTCDate(previous.getUTCDate() + 1)
      if (previous.toISOString().slice(0, 10) === option.EndDate) option.EndDate = corrected
    }
    isConfigured.value = rows.length > 0
  } catch { loadError.value = 'ไม่สามารถโหลดการตั้งค่าได้ กรุณาลองใหม่' }
  finally { loading.value = false }
}
const newOption = (): DurationOption => ({ ID: 'o' + Date.now().toString(36) + Math.random().toString(36).slice(2, 10), Duration: '', Unit: 'days', StartDate: '', EndDate: '' })
const addCondition = () => {
  if (conditions.value.length >= 100) return
  const code = 't' + Date.now().toString(36) + Math.random().toString(36).slice(2, 10)
  conditions.value.push({ TradeType: code, Title: '', RequiresContract: true, Options: [newOption()] })
}
const addOption = (item: Condition) => { if (item.Options.length < 100) item.Options.push(newOption()) }
const removeOption = async (item: Condition, option: DurationOption) => {
  if (item.Options.length <= 1) return
  const confirmed = await showConfirm('ต้องการลบระยะเวลานี้หรือไม่ การเปลี่ยนแปลงจะมีผลเมื่อบันทึก', { title: 'ลบระยะเวลา', type: 'danger', icon: 'trash', confirmText: 'ลบ' })
  if (confirmed) item.Options = item.Options.filter(entry => entry.ID !== option.ID)
}
const removeCondition = async (item: Condition) => {
  const confirmed = await showConfirm('ต้องการลบประเภทการค้านี้หรือไม่ การเปลี่ยนแปลงจะมีผลเมื่อบันทึก', { title: 'ลบประเภทการค้า', type: 'danger', icon: 'trash', confirmText: 'ลบ' })
  if (confirmed) conditions.value = conditions.value.filter(row => row.TradeType !== item.TradeType)
}
const updateEndDate = (option: DurationOption) => {
  option.EndDate = calculateTradeEndDate(option.StartDate, option.Duration, option.Unit)
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

const save = async () => {
  error.value = ''
  const rows = conditions.value.map(item => ({ ...item, Title: item.Title.trim(), Options: item.Options.map(option => ({ ...option, Duration: Number(option.Duration) })) }))
  const titles = rows.map(item => item.Title.toLowerCase())
  if (rows.some(item => !item.Title || [...item.Title].length > 255) || new Set(titles).size !== titles.length) {
    error.value = 'กรุณาระบุชื่อประเภทการค้าไม่ซ้ำกัน และไม่เกิน 255 ตัวอักษร'
    return
  }
  if (rows.some(item => !item.Options.length || item.Options.some(option => !Number.isInteger(option.Duration) || option.Duration < 1 || option.Duration > 10000 || !['days', 'months', 'years'].includes(option.Unit)))) {
    error.value = 'กรุณาระบุจำนวนระยะเวลาเป็นเลขจำนวนเต็มตั้งแต่ 1 ถึง 10,000'
    return
  }
  if (rows.some(item => item.Options.some(option => (option.StartDate || option.EndDate) && !durationBetween(option.StartDate, option.EndDate)))) {
    error.value = 'กรุณาเลือกวันเริ่มและวันสิ้นสุดให้ถูกต้อง'
    return
  }
  const saveSnapshot = JSON.stringify(conditions.value)
  saving.value = true
  try {
    await put('/markets/' + encodeURIComponent(market.value) + '/trade-conditions', { Conditions: rows })
    savedConditions.value = saveSnapshot
    isConfigured.value = rows.length > 0
    await showSuccess('บันทึกเงื่อนไขแล้ว', 'สำเร็จ')
  } catch { error.value = 'ไม่สามารถบันทึกการตั้งค่าได้ กรุณาลองใหม่' }
  finally { saving.value = false }
}
onMounted(() => { load(); watch(market, load) })
</script>
<style scoped>
.page-heading{display:flex;flex-wrap:wrap;align-items:center;justify-content:space-between;gap:12px 20px;margin:8px 0 18px}
.heading-actions{display:inline-flex;flex-wrap:wrap;align-items:center;gap:10px}
.page-title{margin:0;font-size:var(--font-size-2xl);font-weight:700;color:var(--tx-0f172a)}
.page-desc{margin:4px 0 0;font-size:var(--font-size-sm);font-weight:500;color:var(--tx-64748b)}
.draft-status{display:inline-flex;align-items:center;gap:7px;min-height:34px;box-sizing:border-box;padding:0 12px;border-radius:999px;background:var(--bg-f1f5f9);color:var(--tx-475569);font-size:var(--font-size-xs);font-weight:600}
.draft-status i{width:8px;height:8px;border-radius:50%;background:currentColor}
.draft-status.published{background:var(--bg-ecfdf5);color:var(--tx-15803d);border:1px solid var(--bd-bbf7d0)}
.draft-status.is-dirty{background:var(--bg-fff7ed);color:var(--tx-c2410c);border:1px solid var(--bd-fed7aa)}
.draft-status.is-dirty i{background:#ea580c;box-shadow:0 0 0 3px rgba(234,88,12,.18)}
.settings-card{padding:24px;border:1px solid var(--bd-e2e8f0);border-radius:16px;background:var(--bg-ffffff);box-shadow:0 1px 3px rgba(0,0,0,.03)}
.section-heading{display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:12px;margin-bottom:18px;padding-bottom:14px;border-bottom:1px solid var(--bd-f1f5f9)}.empty-state{color:var(--tx-64748b);font-size:var(--font-size-sm)}
.section-title{margin:0;color:var(--tx-1e293b);font-size:var(--font-size-lg);font-weight:700}
.section-subtitle{margin:4px 0 0;color:var(--tx-64748b);font-size:var(--font-size-sm);font-weight:500}
.conditions-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(min(100%,340px),1fr));gap:14px;align-items:start}
.condition{min-width:0;box-sizing:border-box;padding:0 16px 16px;border:1px solid var(--bd-e2e8f0);border-radius:14px;background:var(--bg-f8fafc)}
.condition-heading{display:flex;align-items:center;gap:12px;padding:14px 0 12px;border-bottom:1px solid var(--bd-e2e8f0);color:var(--tx-334155)}
.condition-heading h3{flex:1;min-width:0;overflow-wrap:anywhere;margin:0;color:var(--tx-1e293b);font-size:var(--font-size-lg);font-weight:700}
.fields{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:14px;margin-top:14px}
.options-heading,.option-heading{display:flex;align-items:center;justify-content:space-between;gap:10px}.options-heading strong,.option-heading strong{font-size:var(--font-size-sm);color:var(--tx-334155)}.option-card{padding:14px;border:1px solid var(--bd-e2e8f0);border-radius:10px;background:var(--bg-ffffff)}.option-card .fields{margin-top:12px}
.contract-requirement{display:flex;flex-wrap:wrap;align-items:center;justify-content:space-between;gap:8px}.contract-requirement > span{font-size:var(--font-size-xs);font-weight:600;color:var(--tx-64748b)}
.field{min-width:0}.field-full{grid-column:1/-1}.field label{display:block;margin-bottom:8px;color:var(--tx-334155);font-size:var(--font-size-sm);font-weight:600}
.term-result{padding:10px 12px;border:1px solid var(--bd-e2e8f0);border-radius:8px;background:var(--bg-ffffff);color:var(--tx-334155);font-size:var(--font-size-sm);font-weight:600}.remaining-days{display:block;margin-top:4px;font-weight:500;color:var(--tx-64748b)}.term-result.is-invalid{color:var(--tx-dc2626)}
.req-star,.error{color:var(--tx-dc2626)}.error{margin:14px 0 0;font-size:var(--font-size-sm)}.notice{color:var(--tx-64748b);font-size:var(--font-size-sm)}
@media(max-width:600px){.page-title{font-size:var(--font-size-xl)}.heading-actions{width:100%}.settings-card{padding:16px}}
</style>