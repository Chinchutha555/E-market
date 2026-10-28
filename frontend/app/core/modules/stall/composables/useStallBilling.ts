// สถานะกลางของงานจัดการค่าใช้จ่าย (ข้อมูลจำลอง ยังไม่เชื่อม API)
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
const billingMonth = ref(isoToday().slice(0, 7))
const amountRange = ref<AmountRange>('all')
const utilityBills = reactive<Record<string, UtilityBill>>({})
const tempBills = ref<TempBill[]>([])
let tempSeeded = false

const monthKey = (code: string) => `${billingMonth.value}:${code.toUpperCase()}`

export const useStallBilling = () => {
  // ---------- สาธารณูปโภค ----------
  const ensureUtilityBills = (occupied: { code: string; tenant: string }[]) => {
    occupied.forEach(({ code, tenant }) => {
      const key = monthKey(code)
      if (utilityBills[key]) return
      const h = hash(key)
      const prevWater = 100 + (hash(code) % 400)
      const prevElectric = 1000 + (hash(code + 'e') % 3000)
      const k = h % 10
      const status: UtilityStatus = k < 3 ? 'unrecorded' : k < 6 ? 'pending' : k < 7 ? 'overdue' : 'paid'
      utilityBills[key] = {
        stall: code.toUpperCase(),
        tenant,
        prevWater,
        prevElectric,
        currWater: status === 'unrecorded' ? null : prevWater + 6 + (h % 19),
        currElectric: status === 'unrecorded' ? null : prevElectric + 40 + (h % 90),
        status
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
  const utilityTotal = (b: UtilityBill) => {
    const w = waterUnits(b)
    const e = electricUnits(b)
    return w == null || e == null ? null : w * rates.water + e * rates.electric
  }
  const saveUtilityReadings = (code: string, water: number, electric: number) => {
    const b = getUtility(code)
    if (!b) return
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
    waterUnits,
    electricUnits,
    utilityTotal,
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
