// สถานะกลางของรอบจดมิเตอร์ ใช้ข้อมูลแผง สัญญา มิเตอร์ และอัตราค่าน้ำ-ไฟร่วมกัน
// ใช้ร่วมกันระหว่าง: แผงรายการบิล, ลิ้นชักรายแผง, การ์ดกรอง และการไฮไลต์แผงบนแผนที่
import { computed, reactive, ref } from 'vue'

export type BillStatus = 'pending' | 'overdue' | 'paid'
export type UtilityStatus = 'unrecorded' | BillStatus
export type TempBillStatus = 'unissued' | BillStatus
export type CostFilter = 'all' | 'todo' | 'pending' | 'overdue' | 'paid'
export type AmountRange = 'all' | 'lt500' | '500-1000' | 'gt1000'

export interface UtilityBill {
  stall: string
  tenant: string
  contractId?: string
  contractNumber?: string
  shopName?: string
  waterMeterCode?: string
  electricMeterCode?: string
  prevWater: number
  prevElectric: number
  currWater: number | null
  currElectric: number | null
  status: UtilityStatus
}

export interface TempBill {
  id: string
  stall: string
  customer: string
  phone: string
  startDate: string
  endDate: string
  dailyRate: number
  extra: number
  discount: number
  status: BillStatus
  issuedAt: string
}

export const BILL_STATUS_LABEL: Record<UtilityStatus, string> = {
  unrecorded: 'ยังไม่จดมิเตอร์',
  pending: 'รอชำระ',
  overdue: 'เกินกำหนด',
  paid: 'ชำระแล้ว'
}

const hash = (s: string) => {
  let h = 0
  for (let i = 0; i < s.length; i++) h = (h * 31 + s.charCodeAt(i)) >>> 0
  return h
}

const toDay = (iso: string) => new Date(`${iso}T00:00:00`).getTime()
export const daysBetween = (start: string, end: string) => {
  if (!start || !end) return 0
  const d = Math.round((toDay(end) - toDay(start)) / 86400000) + 1
  return d > 0 ? d : 0
}
const isoToday = () => new Date().toISOString().slice(0, 10)
const isoShift = (days: number) => new Date(Date.now() + days * 86400000).toISOString().slice(0, 10)

// ----- state (module-level singleton) -----
const rates = reactive({ water: 18, electric: 8 })

// ----- อัตราตามที่ตั้งไว้ในหน้าตั้งค่าน้ำ-ไฟ (รายแผง -> โครงการ -> มาตรฐาน) -----
export type ChargeKind = 'water' | 'electric'
export interface ChargeRule { mode: 'meter' | 'flat' | 'none'; rate: number; flat: number; surcharge: number }
type RuleResolver = (code: string, kind: ChargeKind) => ChargeRule | null
let ruleResolver: RuleResolver | null = null
const resolverVersion = ref(0)
export const registerRuleResolver = (fn: RuleResolver) => { ruleResolver = fn; resolverVersion.value++ }
// แผง -> โครงการ (ลงทะเบียนจากหน้าข้อมูลตลาด เพื่อให้หาอัตราระดับโครงการได้)
const stallProjects = reactive<Record<string, string>>({})
export const registerStallProjects = (map: Record<string, string>) => { Object.assign(stallProjects, map) }
export const projectOfStall = (code: string) => stallProjects[code.toUpperCase()] || ''

const fallbackRule = (kind: ChargeKind): ChargeRule => ({ mode: 'meter', rate: rates[kind], flat: 0, surcharge: 0 })
export const chargeRuleOf = (code: string, kind: ChargeKind): ChargeRule => {
  void resolverVersion.value
  return ruleResolver?.(code, kind) || fallbackRule(kind)
}
// ค่าใช้จ่ายของมิเตอร์หนึ่งประเภท: ตามมิเตอร์ = หน่วย x อัตรา + ค่าบริการ, เหมาจ่าย = ยอดเหมา, ไม่คิดค่าใช้บริการ = 0 (null = ยังไม่มีเลขมิเตอร์)
export const chargeFor = (rule: ChargeRule, units: number | null): number | null => {
  if (rule.mode === 'none') return 0
  if (rule.mode === 'flat') return rule.flat
  if (units == null) return null
  return Math.max(0, units) * (rule.rate + (rule.surcharge || 0))
}
// ประเภทที่ต้องจดมิเตอร์ (เฉพาะตามมิเตอร์)
export const needsReading = (code: string, kind: ChargeKind) => chargeRuleOf(code, kind).mode === 'meter'
const billingMonth = ref(isoToday().slice(0, 7))
const amountRange = ref<AmountRange>('all')
const utilityBills = reactive<Record<string, UtilityBill>>({})
// แผงที่ข้ามการจดในรอบนี้ (พร้อมเหตุผล) key = เดือน:รหัสแผง
const utilitySkips = reactive<Record<string, string>>({})
const tempBills = ref<TempBill[]>([])
let tempSeeded = false

const monthKey = (code: string) => `${billingMonth.value}:${code.toUpperCase()}`

export const useStallBilling = () => {
  // ---------- สาธารณูปโภค ----------
  const ensureUtilityBills = (occupied: Array<{
    code: string
    tenant: string
    contractId?: string
    contractNumber?: string
    shopName?: string
    prevWater?: number | null
    prevElectric?: number | null
    waterMeterCode?: string
    electricMeterCode?: string
    currWater?: number | null
    currElectric?: number | null
    status?: UtilityStatus
  }>) => {
    occupied.forEach((context) => {
      const { code, tenant } = context
      const key = monthKey(code)
      const existing = utilityBills[key]
      if (existing) {
        existing.tenant = tenant || existing.tenant
        existing.contractId = context.contractId || existing.contractId
        existing.contractNumber = context.contractNumber || existing.contractNumber
        existing.shopName = context.shopName || existing.shopName
        existing.waterMeterCode = context.waterMeterCode || existing.waterMeterCode
        existing.electricMeterCode = context.electricMeterCode || existing.electricMeterCode
        if (context.currWater != null) existing.currWater = context.currWater
        if (context.currElectric != null) existing.currElectric = context.currElectric
        if (context.status) existing.status = context.status
        return
      }
      const prevWater = context.prevWater ?? (100 + (hash(code) % 400))
      const prevElectric = context.prevElectric ?? (1000 + (hash(code + 'e') % 3000))
      utilityBills[key] = {
        stall: code.toUpperCase(),
        tenant,
        contractId: context.contractId,
        contractNumber: context.contractNumber,
        shopName: context.shopName,
        waterMeterCode: context.waterMeterCode,
        electricMeterCode: context.electricMeterCode,
        prevWater,
        prevElectric,
        currWater: context.currWater ?? null,
        currElectric: context.currElectric ?? null,
        status: context.status || 'unrecorded'
      }
    })
  }

  const utilityList = computed(() =>
    Object.entries(utilityBills)
      .filter(([k]) => k.startsWith(`${billingMonth.value}:`))
      .map(([, v]) => v)
      .sort((a, b) => a.stall.localeCompare(b.stall, 'th', { numeric: true }))
  )
  const getUtility = (code: string) => utilityBills[monthKey(code)] || null
  const waterUnits = (b: UtilityBill) => (b.currWater == null ? null : b.currWater - b.prevWater)
  const electricUnits = (b: UtilityBill) => (b.currElectric == null ? null : b.currElectric - b.prevElectric)
  const waterCharge = (b: UtilityBill) => chargeFor(chargeRuleOf(b.stall, 'water'), waterUnits(b))
  const electricCharge = (b: UtilityBill) => chargeFor(chargeRuleOf(b.stall, 'electric'), electricUnits(b))
  const utilityTotal = (b: UtilityBill) => {
    const w = waterCharge(b)
    const e = electricCharge(b)
    return w == null || e == null ? null : w + e
  }
  const getUtilitySkip = (code: string) => utilitySkips[monthKey(code)] || ''
  const skipUtility = (code: string, reason: string) => {
    utilitySkips[monthKey(code)] = reason
  }
  const clearUtilitySkip = (code: string) => {
    delete utilitySkips[monthKey(code)]
  }
  const saveUtilityReadings = (code: string, water: number, electric: number) => {
    const b = getUtility(code)
    if (!b) return
    clearUtilitySkip(code)
    b.currWater = water
    b.currElectric = electric
    if (b.status === 'unrecorded') b.status = 'pending'
  }
  const markUtilityPaid = (code: string) => {
    const b = getUtility(code)
    if (b && b.status !== 'unrecorded') b.status = 'paid'
  }

  // ---------- บิลแผงเช่าชั่วคราว ----------
  const seedTempBills = (vacantCodes: string[]) => {
    if (tempSeeded || vacantCodes.length < 3) return
    tempSeeded = true
    const mk = (i: number, customer: string, phone: string, from: number, to: number, status: BillStatus): TempBill => ({
      id: `TB-${String(i).padStart(4, '0')}`,
      stall: vacantCodes[i].toUpperCase(),
      customer,
      phone,
      startDate: isoShift(from),
      endDate: isoShift(to),
      dailyRate: 150,
      extra: 0,
      discount: 0,
      status,
      issuedAt: isoShift(from - 1)
    })
    tempBills.value = [
      mk(0, 'สมหญิง ขายดี', '081-234-5678', -12, -10, 'paid'),
      mk(1, 'ประยุทธ์ ใจกล้า', '089-111-2233', -4, -2, 'overdue'),
      mk(2, 'ลลิตา ค้าขาย', '086-777-8899', 2, 5, 'pending')
    ]
  }
  const tempDays = (b: Pick<TempBill, 'startDate' | 'endDate'>) => daysBetween(b.startDate, b.endDate)
  const tempTotal = (b: Pick<TempBill, 'startDate' | 'endDate' | 'dailyRate' | 'extra' | 'discount'>) =>
    Math.max(0, tempDays(b) * b.dailyRate + (Number(b.extra) || 0) - (Number(b.discount) || 0))
  const tempConflict = (stall: string, start: string, end: string) =>
    tempBills.value.find(b => b.stall === stall.toUpperCase() && !(end < b.startDate || start > b.endDate)) || null
  const issueTempBill = (data: Omit<TempBill, 'id' | 'status' | 'issuedAt'>) => {
    const bill: TempBill = {
      ...data,
      stall: data.stall.toUpperCase(),
      id: `TB-${String(tempBills.value.length + 1).padStart(4, '0')}`,
      status: 'pending',
      issuedAt: isoToday()
    }
    tempBills.value = [bill, ...tempBills.value]
    return bill
  }
  const markTempPaid = (id: string) => {
    const b = tempBills.value.find(x => x.id === id)
    if (b) b.status = 'paid'
  }
  const tempBillsOf = (code: string) => tempBills.value.filter(b => b.stall === code.toUpperCase())

  // ---------- ตัวกรองการ์ด ----------
  const matchesFilter = (status: UtilityStatus, filter: CostFilter) =>
    filter === 'all' || (filter === 'todo' ? status === 'unrecorded' : status === filter)

  const matchesAmount = (amount: number | null) => {
    const r = amountRange.value
    if (r === 'all') return true
    if (amount == null) return false
    return r === 'lt500' ? amount < 500 : r === '500-1000' ? amount >= 500 && amount <= 1000 : amount > 1000
  }

  return {
    rates,
    amountRange,
    matchesAmount,
    billingMonth,
    tempBills,
    utilityList,
    ensureUtilityBills,
    getUtility,
    getUtilitySkip,
    skipUtility,
    clearUtilitySkip,
    waterUnits,
    electricUnits,
    utilityTotal,
    waterCharge,
    electricCharge,
    saveUtilityReadings,
    markUtilityPaid,
    seedTempBills,
    tempDays,
    tempTotal,
    tempConflict,
    issueTempBill,
    markTempPaid,
    tempBillsOf,
    matchesFilter
  }
}
