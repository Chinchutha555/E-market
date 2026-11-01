// ตั้งค่าอัตราค่าน้ำ-ไฟ มิเตอร์ และประวัติการจดรายรอบ (ซิงก์กับ backend)
// ลำดับการใช้อัตรา: รายแผง -> รายโครงการ -> ค่ากลางของตลาด
import { reactive, ref, watch } from 'vue'
import { projectOfStall, registerRuleResolver, useStallBilling } from './useStallBilling'
import { useApi } from '~/composables/useApi'

export type RuleMode = 'inherit' | 'meter' | 'flat' | 'none'
export type UtilityKind = 'electric' | 'water'
export interface UtilityRule { mode: RuleMode; rate: number; minMonthly: number; flat: number; flatDaily: number; surcharge?: number }
export interface UtilityRates { electric: UtilityRule; water: UtilityRule }
export interface MeterDevice {
  id: string
  name: string
  kind: UtilityKind
  code: string
  currentReading: number | null
  /** รองรับข้อมูลเดิมก่อนเปลี่ยนชื่อ field เป็น currentReading */
  initialReading?: number | null
  startDate?: string
  status: 'active' | 'inactive'
}
export interface MeterBinding { electric: string; water: string; devices?: MeterDevice[] }
export interface MeterReadingRecord {
  billingMonth: string
  stall: string
  contractId: string
  contractNumber: string
  tenant: string
  shopName: string
  waterMeterCode: string
  electricMeterCode: string
  prevWater: number
  prevElectric: number
  currWater: number | null
  currElectric: number | null
  status: 'unrecorded' | 'pending' | 'overdue' | 'paid'
  recordedAt: string
}
export type RuleSource = 'stall' | 'project' | 'default'

export const RULE_MODE_LABEL: Record<RuleMode, string> = { inherit: 'ตามระดับบน', meter: 'ตามมิเตอร์', flat: 'เหมาจ่าย', none: 'ไม่คิดค่าใช้บริการ' }
export const UTILITY_UNIT: Record<UtilityKind, string> = { electric: 'kWh', water: 'ลบ.ม.' }

export type RateScope = 'market' | 'project' | 'stall'
interface MarketSettings {
  scope?: RateScope
  defaults: UtilityRates
  projects: Record<string, UtilityRates>
  stalls: Record<string, UtilityRates & { projectId: string }>
  meters: Record<string, MeterBinding>
  readings?: Record<string, MeterReadingRecord>
}
const store = reactive<Record<string, MarketSettings>>({})
// สถานะซิงก์กับ backend ต่อ market: โหลดครั้งเดียว แล้วบันทึกอัตโนมัติ (หน่วงเล็กน้อย) เมื่อค่าที่บันทึกแล้วเปลี่ยน
const loadPromises: Record<string, Promise<void>> = {}
const hydrated: Record<string, boolean> = {}
const syncTimers: Record<string, ReturnType<typeof setTimeout>> = {}
const resolverRegistered: Record<string, boolean> = {}
export const utilitySyncError = ref('')
const key = (marketId: string) => `emarket-utility-settings-v2-${marketId}`
const clone = <T>(value: T): T => JSON.parse(JSON.stringify(value))
const blankRule = (mode: RuleMode = 'inherit', rate = 0): UtilityRule => ({ mode, rate, minMonthly: 0, flat: 0, flatDaily: 0 })
const inheritAll = (): UtilityRates => ({ electric: blankRule(), water: blankRule() })
const isInheritAll = (value: UtilityRates) => value.electric.mode === 'inherit' && value.water.mode === 'inherit'

export const describeRule = (rule: UtilityRule, kind: UtilityKind) => {
  if (rule.mode === 'none') return `ไม่คิด${kind === 'electric' ? 'ค่าไฟฟ้า' : 'ค่าน้ำประปา'}`
  if (rule.mode === 'meter') {
    const min = rule.minMonthly > 0 ? ` · ขั้นต่ำ ฿${rule.minMonthly.toLocaleString('th-TH')}/เดือน` : ''
    const extra = (rule.surcharge || 0) > 0 ? ` + ฿${rule.surcharge!.toLocaleString('th-TH')}/${UTILITY_UNIT[kind]}` : ''
    return `ตามมิเตอร์ ฿${rule.rate}/${UTILITY_UNIT[kind]}${extra}${min}`
  }
  if (rule.mode === 'flat') {
    const daily = rule.flatDaily > 0 ? ` · ฿${rule.flatDaily.toLocaleString('th-TH')}/วัน` : ''
    return `เหมาจ่าย ฿${rule.flat.toLocaleString('th-TH')}/เดือน${daily}`
  }
  return RULE_MODE_LABEL.inherit
}

export const useUtilitySettings = (marketId: string) => {
  const { rates } = useStallBilling()

  if (!store[marketId]) {
    let saved: MarketSettings | null = null
    try { saved = JSON.parse(localStorage.getItem(key(marketId)) || 'null') } catch { saved = null }
    store[marketId] = saved?.defaults
      ? saved
      : {
          defaults: { electric: blankRule('meter', rates.electric), water: blankRule('meter', rates.water) },
          projects: {}, stalls: {}, meters: {}, readings: {}
        }
  }
  const state = store[marketId]
  const { get, put } = useApi()

  const pushRemote = async () => {
    try {
      await put(`/markets/${marketId}/utility-settings`, { settings: clone(state) })
      utilitySyncError.value = ''
      return true
    } catch (err) {
      const error = err as any
      const apiMessage = error?.data?.message || error?.response?._data?.message
      const status = error?.statusCode || error?.response?.status
      utilitySyncError.value = apiMessage || (!status
        ? 'เชื่อมต่อ backend ไม่ได้ กรุณาเปิดบริการ API แล้วลองบันทึกอีกครั้ง'
        : 'บันทึกการตั้งค่าน้ำไฟและมิเตอร์ขึ้นระบบไม่สำเร็จ')
      console.warn('utility settings sync failed', err)
      return false
    }
  }
  const syncNow = async () => {
    clearTimeout(syncTimers[marketId])
    return pushRemote()
  }
  if (!loadPromises[marketId]) {
    // ค่าที่แคชไว้ในเบราว์เซอร์ใช้แสดงก่อน แล้วค่าจาก backend ทับเมื่อโหลดเสร็จ
    watch(() => store[marketId], (value) => {
      try { localStorage.setItem(key(marketId), JSON.stringify(value)) } catch { /* ignore */ }
      if (!hydrated[marketId]) return
      clearTimeout(syncTimers[marketId])
      syncTimers[marketId] = setTimeout(pushRemote, 600)
    }, { deep: true })
    loadPromises[marketId] = (async () => {
      try {
        const res: any = await get(`/markets/${marketId}/utility-settings`)
        const remote = (res?.data ?? res)?.settings
        if (remote?.defaults) {
          const target = state
          target.scope = remote.scope
          target.defaults = remote.defaults
          target.projects = remote.projects || {}
          target.stalls = remote.stalls || {}
          target.meters = remote.meters || {}
          target.readings = remote.readings || {}
        } else if (Object.keys(state.meters).length || Object.keys(state.projects).length || Object.keys(state.stalls).length) {
          // backend ยังไม่มีข้อมูล แต่เครื่องนี้เคยตั้งไว้: ส่งขึ้นระบบครั้งแรก
          hydrated[marketId] = true
          await pushRemote()
        }
      } catch (err) {
        console.warn('utility settings load failed', err)
      } finally {
        hydrated[marketId] = true
      }
    })()
  }
  const ready = loadPromises[marketId]
  // ค่าเริ่มต้นคือ ทั้งตลาดแบบเดียวกัน (เปลี่ยนเมื่อผู้ใช้เลือกรูปแบบอื่นและบันทึก)
  if (!state.scope) state.scope = 'market'
  const getScope = (): RateScope => state.scope || 'market'
  const setScope = (value: RateScope) => { state.scope = value }

  const getDefaults = () => clone(state.defaults)
  const setDefaults = (value: UtilityRates) => {
    state.defaults = clone(value)
    if (value.electric.mode === 'meter') rates.electric = value.electric.rate
    if (value.water.mode === 'meter') rates.water = value.water.rate
  }

  const getProjectRates = (projectId: string): UtilityRates => clone(state.projects[projectId] || inheritAll())
  const setProjectRates = (projectId: string, value: UtilityRates) => {
    if (isInheritAll(value)) delete state.projects[projectId]
    else state.projects[projectId] = clone(value)
  }
  const getStallRates = (code: string): UtilityRates => {
    const { projectId: _p, ...rest } = state.stalls[code.toUpperCase()] || ({} as any)
    return clone(rest.electric ? rest : inheritAll())
  }
  const setStallRates = (code: string, projectId: string, value: UtilityRates) => {
    const id = code.toUpperCase()
    if (isInheritAll(value)) delete state.stalls[id]
    else state.stalls[id] = { ...clone(value), projectId }
  }
  const customStallCount = (projectId: string) => Object.values(state.stalls).filter(s => s.projectId === projectId).length

  // อัตราที่ใช้จริงของแผง พร้อมบอกว่ามาจากระดับไหน
  const resolve = (projectId: string, code: string, kind: UtilityKind): { rule: UtilityRule; source: RuleSource } => {
    if (state.scope === 'market') return { rule: state.defaults[kind], source: 'default' }
    const stall = state.scope !== 'market' ? state.stalls[code.toUpperCase()]?.[kind] : undefined
    if (stall && stall.mode !== 'inherit') return { rule: stall, source: 'stall' }
    const project = state.projects[projectId]?.[kind]
    if (project && project.mode !== 'inherit') return { rule: project, source: 'project' }
    return { rule: state.defaults[kind], source: 'default' }
  }
  const resolveProject = (projectId: string, kind: UtilityKind) => {
    if (state.scope === 'market') return { rule: state.defaults[kind], source: 'default' as RuleSource }
    const project = state.projects[projectId]?.[kind]
    return project && project.mode !== 'inherit' ? { rule: project, source: 'project' as RuleSource } : { rule: state.defaults[kind], source: 'default' as RuleSource }
  }

  const normalizeDevices = (binding?: MeterBinding): MeterDevice[] => {
    if (binding?.devices?.length) return binding.devices.map(({ initialReading, startDate: _startDate, ...device }) => ({
      ...device,
      currentReading: device.currentReading ?? initialReading ?? null
    }))
    const devices: MeterDevice[] = []
    if (binding?.electric) devices.push({ id: `legacy-electric`, name: 'ไฟฟ้า', kind: 'electric', code: binding.electric, currentReading: null, status: 'active' })
    if (binding?.water) devices.push({ id: `legacy-water`, name: 'น้ำประปา', kind: 'water', code: binding.water, currentReading: null, status: 'active' })
    return devices
  }
  const getMeter = (code: string): MeterBinding => {
    const binding = state.meters[code.toUpperCase()] || { electric: '', water: '' }
    return { ...clone(binding), devices: normalizeDevices(binding) }
  }
  const setMeter = (code: string, value: MeterBinding) => {
    const id = code.toUpperCase()
    const devices = normalizeDevices(value).map(device => ({ ...device, code: device.code.trim().toUpperCase() }))
    const electric = devices.find(device => device.status === 'active' && device.kind === 'electric')?.code || value.electric.trim().toUpperCase()
    const water = devices.find(device => device.status === 'active' && device.kind === 'water')?.code || value.water.trim().toUpperCase()
    if (!electric && !water && !devices.length) delete state.meters[id]
    else state.meters[id] = { electric, water, devices }
  }
  // อัปเดตมิเตอร์หลายแผงในครั้งเดียว ลด reactive recalculation ระหว่างโหลดข้อมูลโครงการใหญ่
  const setMeters = (values: Record<string, MeterBinding>) => {
    const next = { ...state.meters }
    Object.entries(values).forEach(([code, value]) => {
      const id = code.toUpperCase()
      const devices = normalizeDevices(value).map(device => ({ ...device, code: device.code.trim().toUpperCase() }))
      const electric = devices.find(device => device.status === 'active' && device.kind === 'electric')?.code || value.electric.trim().toUpperCase()
      const water = devices.find(device => device.status === 'active' && device.kind === 'water')?.code || value.water.trim().toUpperCase()
      if (!electric && !water && !devices.length) delete next[id]
      else next[id] = { electric, water, devices }
    })
    state.meters = next
  }
  const duplicateOf = (code: string, kind: UtilityKind, meterCode: string) => {
    if (!meterCode) return ''
    return Object.entries(state.meters).find(([c, m]) => c !== code.toUpperCase() && m[kind] === meterCode)?.[0] || ''
  }
  const readingKey = (billingMonth: string, code: string) => `${billingMonth}:${code.trim().toUpperCase()}`
  const getReading = (billingMonth: string, code: string): MeterReadingRecord | null => {
    const value = state.readings?.[readingKey(billingMonth, code)]
    return value ? clone(value) : null
  }
  const setReading = (value: MeterReadingRecord) => {
    if (!state.readings) state.readings = {}
    state.readings[readingKey(value.billingMonth, value.stall)] = clone({ ...value, stall: value.stall.toUpperCase() })
  }
  // ให้การคิดเงินในหน้าอื่นใช้อัตราที่ตั้งไว้ (รายแผง -> โครงการ -> มาตรฐาน) ลงทะเบียนครั้งเดียวต่อตลาด
  if (!resolverRegistered[marketId]) {
  resolverRegistered[marketId] = true
  registerRuleResolver((code, kind) => {
    const { rule } = resolve(projectOfStall(code), code, kind)
    if (rule.mode === 'inherit') return null
    return { mode: rule.mode, rate: rule.rate, flat: rule.flat, surcharge: rule.surcharge || 0 }
  })
  }

  return { ready, getScope, setScope, getDefaults, setDefaults, getProjectRates, setProjectRates, getStallRates, setStallRates, customStallCount, resolve, resolveProject, getMeter, setMeter, setMeters, duplicateOf, getReading, setReading, syncNow }
}
