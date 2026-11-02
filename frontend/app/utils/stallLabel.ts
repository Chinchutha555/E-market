// ชื่อแผงที่แสดงให้คน: หมายเลขแผงสั้น (เช่น 116) แทนรหัสเต็มของผัง (เช่น JJ-17-116)
// รหัสเต็มใช้เป็นตัวอ้างอิงในระบบเท่านั้น ห้ามแสดงตรงๆ
const stallNumberRegistry = new Map<string, string>()

// หน้าไหนโหลดข้อมูลผัง GIS แล้ว ให้ลงทะเบียนหมายเลขจริงของแต่ละแผงไว้ (exact)
export const registerStallNumbers = (projects: any[] | null | undefined) => {
  if (!Array.isArray(projects)) return
  const add = (stall: any) => {
    const code = String(stall?.code ?? '').trim()
    const props = stall?.sourceProperties
    const number = String(props?.number ?? props?.Number ?? '').trim()
    if (code && number) stallNumberRegistry.set(code, number)
  }
  for (const project of projects) {
    if (Array.isArray(project?.stalls)) project.stalls.forEach(add)
    if (Array.isArray(project?.floors)) project.floors.forEach((floor: any) => Array.isArray(floor?.stalls) && floor.stalls.forEach(add))
  }
}

// รหัสผัง: JJ-{โครงการ}[-{โซน}]-{หมายเลข}[-{ลำดับ}] ใช้เมื่อยังไม่ได้โหลดข้อมูลผัง
const numberFromCode = (code: string): string => {
  const parts = code.split('-')
  if (parts[0] !== 'JJ' || parts.length < 3) return code
  if (parts.length === 3) return parts[2] || code
  if (parts.length === 4) {
    const [, , third, fourth] = parts
    return (fourth || '').length >= 4 && (fourth || '').length > (third || '').length ? third || code : fourth || code
  }
  if (parts.length === 5) return parts[3] || code
  return code
}

export const stallShortNumber = (value?: string | number | null): string => {
  const raw = String(value ?? '').trim()
  if (!raw || raw === '-') return raw
  // หลายแผงในช่องเดียว เช่น "JJ-17-116, JJ-17-117"
  if (raw.includes(',')) return raw.split(',').map(part => stallShortNumber(part)).filter(Boolean).join(', ')
  return stallNumberRegistry.get(raw) || numberFromCode(raw)
}
