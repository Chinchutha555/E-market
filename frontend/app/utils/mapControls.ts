// ชุดปุ่มควบคุมแผนที่มาตรฐานของระบบ (ใช้ร่วมกันทุกแผนที่ MapLibre)
// ลำดับจากบนลงล่าง: ตำแหน่งของฉัน → ซูมเข้า (+) → ซูมออก (−) → เข็มทิศ → ปรับมุมมองให้เห็นทั้งพื้นที่
// อยู่ในกลุ่มเดียวกัน ขนาด 40px (สไตล์อยู่ใน assets/css/main.css)

export interface StandardMapControlOptions {
  /** ฟังก์ชันปรับมุมมองกลับไปให้เห็นทั้งพื้นที่ของแผนที่นั้น */
  onFit: () => void
  /** ชื่อของปุ่มปรับมุมมอง (title/aria-label) */
  fitTitle?: string
  /** ใช้แสดงข้อความเมื่อระบุตำแหน่งไม่ได้ (เช่น showWarning ของระบบ) */
  notify?: (message: string, title?: string) => void
  position?: 'bottom-right' | 'bottom-left' | 'top-right' | 'top-left'
}

const FIT_ICON = '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 8V5a2 2 0 0 1 2-2h3M16 3h3a2 2 0 0 1 2 2v3M21 16v3a2 2 0 0 1-2 2h-3M8 21H5a2 2 0 0 1-2-2v-3"></path><rect x="8" y="8" width="8" height="8" rx="1"></rect></svg>'

const geolocateMessage = (code: number) => {
  if (typeof window !== 'undefined' && !window.isSecureContext) return 'ระบุตำแหน่งไม่ได้: ต้องเปิดระบบผ่านการเชื่อมต่อแบบปลอดภัย (HTTPS)'
  if (code === 1) return 'ไม่ได้รับอนุญาตให้ใช้ตำแหน่ง กรุณาอนุญาตการเข้าถึงตำแหน่งในการตั้งค่าเบราว์เซอร์'
  if (code === 3) return 'ระบุตำแหน่งไม่ทันเวลา กรุณาลองอีกครั้ง'
  return 'ไม่สามารถระบุตำแหน่งปัจจุบันได้ในขณะนี้'
}

/** เพิ่มชุดปุ่มควบคุมมาตรฐานลงแผนที่ คืนค่า element ของกลุ่มปุ่ม (ใช้วัดตำแหน่งได้) */
export const addStandardMapControls = (maplibregl: any, map: any, options: StandardMapControlOptions): HTMLElement | null => {
  const position = options.position || 'bottom-right'
  const navControl = new maplibregl.NavigationControl({ visualizePitch: true })
  map.addControl(navControl, position)
  const navGroup: HTMLElement | undefined = (navControl as any)._container
  if (!navGroup) return null

  // ปุ่มปรับมุมมองให้เห็นทั้งพื้นที่ (ล่างสุด)
  const fitBtn = document.createElement('button')
  fitBtn.type = 'button'
  fitBtn.className = 'map-ctrl-recenter-btn'
  fitBtn.title = options.fitTitle || 'ปรับมุมมองให้เห็นทั้งพื้นที่'
  fitBtn.setAttribute('aria-label', fitBtn.title)
  fitBtn.innerHTML = FIT_ICON
  fitBtn.addEventListener('click', () => options.onFit())
  navGroup.appendChild(fitBtn)

  // ปุ่มตำแหน่งของฉัน (บนสุด) ใช้ตำแหน่งจากเบราว์เซอร์ ไม่ส่งไปที่ระบบ
  try {
    const geolocate = new maplibregl.GeolocateControl({
      positionOptions: { enableHighAccuracy: true, timeout: 12000 },
      trackUserLocation: true,
      showUserLocation: true,
      showAccuracyCircle: true,
      fitBoundsOptions: { maxZoom: 19 }
    })
    map.addControl(geolocate, position)
    const geoContainer: HTMLElement | undefined = (geolocate as any)._container
    // ปุ่มของ GeolocateControl ถูกสร้างแบบ async (หลังตรวจสิทธิ์ตำแหน่ง) จึงรอจนกว่าจะมีปุ่มแล้วค่อยย้ายเข้ากลุ่ม
    const attach = (tries = 0) => {
      const geoBtn = geoContainer?.querySelector('button') as HTMLButtonElement | null
      if (geoContainer && geoBtn) {
        geoBtn.title = 'ตำแหน่งของฉัน'
        geoBtn.setAttribute('aria-label', 'แสดงตำแหน่งของฉันบนแผนที่')
        navGroup.insertBefore(geoBtn, navGroup.firstChild)
        geoContainer.style.display = 'none'
      } else if (tries < 50) {
        window.setTimeout(() => attach(tries + 1), 100)
      }
    }
    attach()
    geolocate.on('error', (e: any) => {
      const message = geolocateMessage(Number(e?.code))
      if (options.notify) options.notify(message, 'ตำแหน่งของฉัน')
      else console.warn('[map]', message)
    })
  } catch (error) {
    console.error('[map] geolocate control unavailable:', error)
  }
  return navGroup
}
