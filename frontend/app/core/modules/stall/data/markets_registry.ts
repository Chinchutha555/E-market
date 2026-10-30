// ============================================================================
// BMA Markets GIS Registry (12 ตลาด กทม.)
// ============================================================================
// Provides base market metadata (ID, slug, name, district, GPS center, camera
// settings). Full GIS zone polygons and floor layouts are stored and served
// exclusively by the Backend (`backend/data/markets/*.json` & `market_gis_config`
// via `GET /api/v1/markets/:slug/gis-config`).
// ============================================================================

import type { MarketFacility, ProjectZone } from './market_gis_types'

export interface MarketGisConfig {
  id: number
  slug: string
  name: string
  district: string
  center: [number, number] // [lng, lat]
  overviewZoom: number
  bearing?: number
  pitch?: number
  description: string
  projects: ProjectZone[]
  facilities?: MarketFacility[]
}

// ============================================================================
// ALL 12 BMA MARKETS MASTER REGISTRY
// ============================================================================
export const BMA_MARKETS_MASTER: Record<string, MarketGisConfig> = {
  '1': {
    id: 1,
    slug: 'chatuchak',
    name: 'ตลาดนัดจตุจักร',
    district: 'เขตจตุจักร',
    center: [100.55036, 13.80045],
    overviewZoom: 16.6,
    bearing: 29.5,
    pitch: 0,
    description: 'ตลาดนัดระดับโลก แหล่งรวมสินค้าหัตถกรรม แฟชั่น ของแต่งบ้าน และต้นไม้ มี 30 โครงการ',
    projects: []
  },
  '2': {
    id: 2,
    slug: 'pracha-niwet-1',
    name: 'ตลาดประชานิเวศน์ 1',
    district: 'เขตจตุจักร',
    center: [100.5454943, 13.8409024],
    overviewZoom: 18.2,
    bearing: 0,
    pitch: 0,
    description: 'ตลาดชุมชนคุณภาพ ศูนย์อาหารปรุงสำเร็จและของฝากเลื่องชื่อย่านประชาชื่น',
    projects: []
  },
  '3': {
    id: 3,
    slug: 'min-buri',
    name: 'ตลาดมีนบุรี',
    district: 'เขตมีนบุรี',
    center: [100.7254896, 13.8150875],
    overviewZoom: 17.8,
    bearing: 15,
    pitch: 0,
    description: 'ศูนย์กลางการค้าและตลาดนัดใหญ่ที่สุดของกรุงเทพฯ ฝั่งตะวันออก มี 2 ชั้น (ชั้นล่างและชั้นบน)',
    projects: []
  },
  '4': {
    id: 4,
    slug: 'singhanat',
    name: 'ตลาดสิงหา',
    district: 'เขตคลองเตย',
    center: [100.5967214, 13.7113411],
    overviewZoom: 18.8,
    bearing: 10,
    pitch: 0,
    description: 'ตลาดชุมชนใจกลางสุขุมวิท 22 รองรับกลุ่มคนทำงานออฟฟิศ นักท่องเที่ยว และผู้อยู่อาศัยในพื้นที่',
    projects: []
  },
  '5': {
    id: 5,
    slug: 'thewarat',
    name: 'ตลาดเทวราช',
    district: 'เขตดุสิต',
    center: [100.5027484, 13.7712845],
    overviewZoom: 18.4,
    bearing: 22,
    pitch: 0,
    description: 'ตลาดเทวราช (ตลาดเทเวศร์) ริมคลองผดุงกรุงเกษม แหล่งจำหน่ายไม้ดอกไม้ประดับ อาหารสด และอาหารปรุงสำเร็จเก่าแก่ของดุสิต มี 4 โซนหลักรวม 475 แผงค้า',
    projects: []
  },
  '6': {
    id: 6,
    slug: 'rat-burana',
    name: 'ตลาดราษฎร์บูรณะ',
    district: 'เขตราษฎร์บูรณะ',
    center: [100.506194, 13.6811691],
    overviewZoom: 18.2,
    bearing: 0,
    pitch: 0,
    description: 'ตลาดสดและสินค้าชุมชนฝั่งธนบุรี บริเวณศูนย์ราชการราษฎร์บูรณะ',
    projects: []
  },
  '7': {
    id: 7,
    slug: 'bang-kapi',
    name: 'ตลาดบางกะปิ',
    district: 'เขตบางกะปิ',
    center: [100.6462099, 13.7671914],
    overviewZoom: 18.2,
    bearing: 12,
    pitch: 0,
    description: 'ตลาดสดและศูนย์การค้าชุมชนขนาดใหญ่ย่านลาดพร้าว-รามคำแหง ติดแนวรถไฟฟ้าสายสีเหลือง',
    projects: []
  },
  '8': {
    id: 8,
    slug: 'wongwian-lek',
    name: 'ตลาดพระเครื่องวงเวียนเล็ก',
    district: 'เขตคลองสาน',
    center: [100.4965505, 13.7351518],
    overviewZoom: 18.6,
    bearing: 5,
    pitch: 0,
    description: 'ศูนย์รวมพระเครื่อง วัตถุมงคล ของเก่า และของสะสมเก่าแก่ระดับตำนานฝั่งธนบุรี',
    projects: []
  },
  '9': {
    id: 9,
    slug: 'nong-chok',
    name: 'ตลาดหนองจอก',
    district: 'เขตหนองจอก',
    center: [100.866388, 13.8562316],
    overviewZoom: 18.4,
    bearing: 20,
    pitch: 0,
    description: 'ตลาดสดชานเมืองริมคลองแสนแสบ แหล่งรวมอาหารสด พืชผักการเกษตรท้องถิ่น และอาหารฮาลาล',
    projects: []
  },
  '10': {
    id: 10,
    slug: 'ratchadaphisek',
    name: 'ตลาดรัชดาภิเษก',
    district: 'เขตธนบุรี',
    center: [100.4770889, 13.7212038],
    overviewZoom: 18.5,
    bearing: 10,
    pitch: 0,
    description: 'ตลาดรัชดาภิเษก (ตลาดพลู) ศูนย์รวมสตรีทฟู้ด ขนมหวาน และอาหารโบราณชื่อดังของฝั่งธนบุรี',
    projects: []
  },
  '11': {
    id: 11,
    slug: 'thonburi-sanamluang-2',
    name: 'ตลาดธนบุรี (สนามหลวง 2)',
    district: 'เขตทวีวัฒนา',
    center: [100.3528837, 13.7469945],
    overviewZoom: 17.2,
    bearing: 0,
    pitch: 0,
    description: 'ตลาดนัดครบวงจรขนาดใหญ่ 11 โซน ครอบคลุมไม้ดอกไม้ประดับ สัตว์เลี้ยง ปลาสวยงาม พระเครื่อง และศูนย์อาหาร',
    projects: []
  },
  '12': {
    id: 12,
    slug: 'bang-khae-phirom',
    name: 'ตลาดบางแคภิรมย์',
    district: 'เขตบางแค',
    center: [100.3836051, 13.6863954],
    overviewZoom: 18.4,
    bearing: 0,
    pitch: 0,
    description: 'ตลาดแห่งใหม่ในสวนสาธารณะบางแคภิรมย์ อาคารโปร่งสบาย 6 โซนสินค้าครบครัน',
    projects: []
  }
}

export function getMarketGisConfig(marketIdOrSlug?: string | number | null): MarketGisConfig {
  const fallback = BMA_MARKETS_MASTER['1']!
  if (!marketIdOrSlug) return fallback
  const key = String(marketIdOrSlug).trim()

  // Match by direct ID
  if (BMA_MARKETS_MASTER[key]) return BMA_MARKETS_MASTER[key]!

  // Match by slug
  const match = Object.values(BMA_MARKETS_MASTER).find(
    m => m.slug.toLowerCase() === key.toLowerCase() || m.name === key
  )
  return match || fallback
}
