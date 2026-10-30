// ================================================================
// Shared Sub-Stall Layout Engine for Market Project Polygons
// - แปลง polygon (lng/lat) เป็นระบบพิกัดเมตรภายใน (local metric plane)
// - หมุนให้แนวขอบที่ยาวที่สุดเป็นแนวนอน แล้ววางแผงย่อยแบบ "แผงหันหลังชนกัน 2 แถว + ทางเดิน"
// - เก็บเฉพาะแผงที่อยู่ภายในกรอบโครงการทั้งหมด (มีระยะร่นจากขอบ)
// ================================================================

import type { ProjectZone, ProjectFloor } from './market_gis_types'

export type StallStatus = 'available' | 'reserved' | 'occupied'

export interface MapStall {
  selectionId?: string
  /** Human-readable stall name from the survey (for example `number`). */
  name?: string
  sourceProperties?: Record<string, unknown>
  code: string
  projectId: string
  zone: string
  status: StallStatus
  areaType: string
  size: number
  widthM: number
  depthM: number
  entryFee: number
  monthlyRent: number
  commonFee: number
  depositFee: number
  photo: string
  coordinates: [number, number][]
  polygonHoles?: [number, number][][]
  center: [number, number]
  floor: string
  floorName: string
  floorLevel: number
}

export interface ProjectStallLayout {
  stalls: MapStall[]
  labelPoint: [number, number]
  total: number
  available: number
  reserved: number
  occupied: number
  floorId: string
  floorName: string
}

export interface MarketStallSummary {
  total: number
  available: number
  reserved: number
  occupied: number
  minRent?: number
  maxRent?: number
  minAreaSqm?: number
  maxAreaSqm?: number
}

// Stall module dimensions (meters)
const STALL_W = 2.5 // ความกว้างหน้าแผง
const STALL_D = 2.0 // ความลึกแผง
const AISLE = 2.0 // ทางเดินระหว่างแถว
const CROSS_AISLE = 2.0 // ทางเดินตัดขวาง
const STALLS_PER_BLOCK = 10 // แผงต่อบล็อกก่อนมีทางเดินตัดขวาง
const EDGE_SETBACK = 0.8 // ระยะร่นจากแนวกรอบโครงการ

type Pt = [number, number]

const M_PER_DEG_LAT = 110540

const makeProjector = (origin: Pt) => {
  const mPerDegLng = 111320 * Math.cos((origin[1] * Math.PI) / 180)
  return {
    toLocal: ([lng, lat]: Pt): Pt => [(lng - origin[0]) * mPerDegLng, (lat - origin[1]) * M_PER_DEG_LAT],
    toLngLat: ([x, y]: Pt): Pt => [origin[0] + x / mPerDegLng, origin[1] + y / M_PER_DEG_LAT]
  }
}

const rotate = ([x, y]: Pt, a: number): Pt => {
  const c = Math.cos(a)
  const s = Math.sin(a)
  return [x * c - y * s, x * s + y * c]
}

const pointInPolygon = ([x, y]: Pt, poly: Pt[]): boolean => {
  let inside = false
  for (let i = 0, j = poly.length - 1; i < poly.length; j = i++) {
    const pi = poly[i]
    const pj = poly[j]
    if (!pi || !pj) continue
    const [xi, yi] = pi
    const [xj, yj] = pj
    if ((yi > y) !== (yj > y) && x < ((xj - xi) * (y - yi)) / (yj - yi) + xi) inside = !inside
  }
  return inside
}

const projectFootprintArea = (project: ProjectZone): number => {
  const ring = project.coordinates
  let twiceArea = 0
  for (let i = 0; i < ring.length - 1; i++) {
    const current = ring[i]
    const next = ring[i + 1]
    if (!current || !next) continue
    twiceArea += current[0] * next[1] - next[0] * current[1]
  }
  return Math.abs(twiceArea) / 2
}

export const isOverviewBoundaryProject = (candidate: ProjectZone, source: ProjectZone[]): boolean => {
  if (candidate.id.startsWith('zone-') || candidate.layoutMode === 'manual') return false
  // A large surveyed stall group can span several buildings, so its convex
  // hull may contain smaller projects even though it is a real selectable
  // project. Min Buri floor 1 is one such group (hundreds of actual stalls).
  if ((candidate.stalls?.length || 0) >= 500) return false
  const candidateArea = projectFootprintArea(candidate)
  const containedProjects = source.filter(project => (
    project.id !== candidate.id &&
    candidateArea > projectFootprintArea(project) * 3 &&
    pointInPolygon(project.centroid, candidate.coordinates)
  ))
  return containedProjects.length >= 2
}

const closePolygonRing = (coordinates: Pt[]): Pt[] => {
  if (coordinates.length === 0) return []
  const ring = coordinates.map(([lng, lat]) => [lng, lat] as Pt)
  const first = ring[0]!
  const last = ring[ring.length - 1]!
  if (first[0] !== last[0] || first[1] !== last[1]) ring.push([first[0], first[1]])
  return ring
}

const getPolygonCenter = (closedRing: Pt[]): Pt => {
  const ring = closedRing.length > 1 ? closedRing.slice(0, -1) : closedRing
  if (!ring.length) return [0, 0]

  const mean: Pt = [
    ring.reduce((sum, point) => sum + point[0], 0) / ring.length,
    ring.reduce((sum, point) => sum + point[1], 0) / ring.length
  ]
  if (pointInPolygon(mean, closedRing)) return mean

  const xs = ring.map(point => point[0])
  const ys = ring.map(point => point[1])
  const minX = Math.min(...xs)
  const maxX = Math.max(...xs)
  const minY = Math.min(...ys)
  const maxY = Math.max(...ys)
  const bboxCenter: Pt = [(minX + maxX) / 2, (minY + maxY) / 2]
  if (pointInPolygon(bboxCenter, closedRing)) return bboxCenter

  // Concave survey polygons can have both centroids outside the footprint.
  // A small deterministic grid finds a safe popup anchor inside the stall.
  for (let row = 1; row < 10; row++) {
    for (let col = 1; col < 10; col++) {
      const candidate: Pt = [minX + ((maxX - minX) * col) / 10, minY + ((maxY - minY) * row) / 10]
      if (pointInPolygon(candidate, closedRing)) return candidate
    }
  }

  return ring[0]!
}

const normalizeStallStatus = (status: unknown): StallStatus => (
  status === 'available' || status === 'reserved' || status === 'occupied' ? status : 'available'
)

const distToSegment = ([px, py]: Pt, [ax, ay]: Pt, [bx, by]: Pt): number => {
  const dx = bx - ax
  const dy = by - ay
  const len2 = dx * dx + dy * dy
  let t = len2 ? ((px - ax) * dx + (py - ay) * dy) / len2 : 0
  t = Math.max(0, Math.min(1, t))
  return Math.hypot(px - (ax + t * dx), py - (ay + t * dy))
}

const distToBoundary = (p: Pt, poly: Pt[]): number => {
  let d = Infinity
  for (let i = 0, j = poly.length - 1; i < poly.length; j = i++) {
    const pj = poly[j]
    const pi = poly[i]
    if (pj && pi) {
      d = Math.min(d, distToSegment(p, pj, pi))
    }
  }
  return d
}

const rowLabel = (i: number): string => {
  let s = ''
  let n = i
  do {
    s = String.fromCharCode(65 + (n % 26)) + s
    n = Math.floor(n / 26) - 1
  } while (n >= 0)
  return s
}

// Interior label point (approximate pole of inaccessibility) so badges sit INSIDE the frame
const findLabelPoint = (poly: Pt[]): Pt => {
  const xs = poly.map(p => p[0])
  const ys = poly.map(p => p[1])
  const minX = Math.min(...xs)
  const maxX = Math.max(...xs)
  const minY = Math.min(...ys)
  const maxY = Math.max(...ys)
  let best: Pt = [(minX + maxX) / 2, (minY + maxY) / 2]
  let bestD = -1
  const steps = 40
  for (let i = 0; i <= steps; i++) {
    for (let j = 0; j <= steps; j++) {
      const p: Pt = [minX + ((maxX - minX) * i) / steps, minY + ((maxY - minY) * j) / steps]
      if (!pointInPolygon(p, poly)) continue
      const d = distToBoundary(p, poly)
      if (d > bestD) {
        bestD = d
        best = p
      }
    }
  }
  return best
}

export const buildProjectStallLayout = (project: ProjectZone, floorId?: string): ProjectStallLayout => {
  // Determine floor info
  const floorList: ProjectFloor[] = project.floors && project.floors.length > 0
    ? project.floors
    : [{ id: '1', name: 'ชั้น 1', shortLabel: 'ชั้น 1', level: 1 }]
  
  const currentFloor: ProjectFloor = floorList.find(f => f.id === floorId) || floorList[0]!
  const effectiveFloorId = currentFloor.id
  const hasMultipleFloors = floorList.length > 1

  // Remove duplicated closing vertex
  const ring = project.coordinates.slice(0, -1) as Pt[]
  const origin: Pt = ring[0] || [0, 0]
  const proj = makeProjector(origin)
  const local = ring.map(proj.toLocal)

  // Orientation: align rows with the longest polygon edge (matches real aisle direction)
  let longest = 0
  let angle = 0
  for (let i = 0; i < local.length; i++) {
    const a = local[i]
    const b = local[(i + 1) % local.length]
    if (!a || !b) continue
    const len = Math.hypot(b[0] - a[0], b[1] - a[1])
    if (len > longest) {
      longest = len
      angle = Math.atan2(b[1] - a[1], b[0] - a[0])
    }
  }

  const rotated = local.map(p => rotate(p, -angle))
  const xs = rotated.map(p => p[0])
  const ys = rotated.map(p => p[1])
  const minX = Math.min(...xs)
  const maxX = Math.max(...xs)
  const minY = Math.min(...ys)
  const maxY = Math.max(...ys)

  const cellFits = (x0: number, y0: number, x1: number, y1: number) => {
    const corners: Pt[] = [[x0, y0], [x1, y0], [x1, y1], [x0, y1], [(x0 + x1) / 2, (y0 + y1) / 2]]
    return corners.every(c => pointInPolygon(c, rotated) && distToBoundary(c, rotated) >= EDGE_SETBACK * 0.5)
  }

  const stalls: MapStall[] = []
  const bandH = STALL_D * 2 + AISLE
  let rowIdx = 0

  // Floor rent multiplier: upper floors have slightly lower rent
  const levelMultiplier = currentFloor.level === 2 ? 0.9 : (currentFloor.level >= 3 ? 0.82 : 1.0)

  for (let bandY = minY + EDGE_SETBACK; bandY + STALL_D <= maxY - EDGE_SETBACK; bandY += bandH) {
    for (let r = 0; r < 2; r++) {
      const y0 = bandY + r * STALL_D
      const y1 = y0 + STALL_D
      if (y1 > maxY - EDGE_SETBACK) break

      const rowCells: { x0: number; x1: number }[] = []
      let col = 0
      let x = minX + EDGE_SETBACK
      while (x + STALL_W <= maxX - EDGE_SETBACK) {
        if (col > 0 && col % STALLS_PER_BLOCK === 0) {
          x += CROSS_AISLE
          col = 0
          if (x + STALL_W > maxX - EDGE_SETBACK) break
        }
        const x0 = x
        const x1 = x + STALL_W
        if (cellFits(x0, y0, x1, y1)) rowCells.push({ x0, x1 })
        x += STALL_W
        col++
      }

      if (rowCells.length === 0) continue
      const letter = rowLabel(rowIdx)
      rowCells.forEach((cell, idx) => {
        const floorTag = hasMultipleFloors ? `${currentFloor.shortLabel}-` : ''
        const code = `${project.sectionNumber}-${floorTag}${letter}${String(idx + 1).padStart(2, '0')}`
        const isCorner = idx === 0 || idx === rowCells.length - 1
        // Generated fallback stalls start empty. Authoritative lock/occupancy
        // state is applied later from the staff-managed data source.
        const status: StallStatus = 'available'

        // แผงเรียงชิดติดกันเป็นบล็อกสถาปัตยกรรม (Seamless Architectural Island)
        const ringLocal: Pt[] = [
          [cell.x0, y0],
          [cell.x1, y0],
          [cell.x1, y1],
          [cell.x0, y1],
          [cell.x0, y0]
        ]
        const coords = ringLocal.map(p => proj.toLngLat(rotate(p, angle)))
        const center = proj.toLngLat(rotate([(cell.x0 + cell.x1) / 2, (y0 + y1) / 2], angle))

        const baseMonthly = isCorner ? 2200 : 1800
        const monthlyRent = Math.round((baseMonthly * levelMultiplier) / 50) * 50
        const baseEntry = isCorner ? 90000 : 75000
        const entryFee = Math.round((baseEntry * levelMultiplier) / 1000) * 1000

        stalls.push({
          selectionId: `${project.id}:${effectiveFloorId}:${code}`,
          code,
          projectId: project.id,
          zone: `${project.name} (${currentFloor.name})`,
          status,
          areaType: isCorner ? 'หัวมุม' : 'ทั่วไป',
          size: STALL_W * STALL_D,
          widthM: STALL_W,
          depthM: STALL_D,
          entryFee,
          monthlyRent,
          commonFee: 500,
          depositFee: monthlyRent,
          photo: '/images/stall_3d_perspective.jpg',
          coordinates: coords,
          center,
          floor: effectiveFloorId,
          floorName: currentFloor.name,
          floorLevel: currentFloor.level
        })
      })
      rowIdx++
    }
  }

  const labelPoint = proj.toLngLat(findLabelPoint(local))
  const available = stalls.filter(s => s.status === 'available').length
  const reserved = stalls.filter(s => s.status === 'reserved').length

  return {
    stalls,
    labelPoint,
    total: stalls.length,
    available,
    reserved,
    occupied: stalls.length - available - reserved,
    floorId: effectiveFloorId,
    floorName: currentFloor.name
  }
}

// Pre-compute layouts for all projects & floors (deterministic)
const layoutCache = new Map<string, ProjectStallLayout>()

export const clearLayoutCache = (): void => {
  layoutCache.clear()
}

export const getProjectStallLayout = (project: ProjectZone, floorId?: string): ProjectStallLayout => {
  const effectiveFloorId = floorId || project.floors?.[0]?.id || '1'
  const currentFloor = project.floors?.find(floor => floor.id === effectiveFloorId)

  // Newly managed zones deliberately begin as an empty canvas. Generating
  // placeholder stalls here would make an unfinished floor look bookable.
  if (project.layoutMode === 'manual' && (!project.stalls || project.stalls.length === 0)) {
    return {
      stalls: [],
      labelPoint: project.centroid,
      total: 0,
      available: 0,
      reserved: 0,
      occupied: 0,
      floorId: effectiveFloorId,
      floorName: currentFloor?.name || 'ชั้น 1'
    }
  }
  const hasSurveyStalls = Boolean(project.stalls && project.stalls.length > 0)
  const cacheKey = `${project.id}:${effectiveFloorId}:${project.centroid.join(',')}:${project.totalStalls}:${hasSurveyStalls ? project.stalls!.length : 'gen'}`

  // If the project carries its own pre-surveyed real stalls (e.g. from official GIS shapefile)
  if (project.stalls && project.stalls.length > 0) {
    const cachedSurveyLayout = layoutCache.get(cacheKey)
    if (cachedSurveyLayout) return cachedSurveyLayout

    const matchingStalls = (project.stalls as MapStall[])
      .filter(s => !s.floor || String(s.floor) === String(effectiveFloorId))
      .map((stall, index) => {
        const coordinates = closePolygonRing(stall.coordinates)
        return {
          ...stall,
          coordinates,
          center: getPolygonCenter(coordinates),
          // Statuses embedded in the GIS file are illustrative. Runtime
          // booking data is applied by MarketProjectsOverviewMap.
          status: (project.layoutMode === 'manual' && stall.status ? stall.status : 'available') as StallStatus,
          selectionId: stall.selectionId || `${project.id}:${effectiveFloorId}:${index}`
        }
      })
    const available = matchingStalls.filter(s => s.status === 'available').length
    const reserved = matchingStalls.filter(s => s.status === 'reserved').length
    const surveyLayout: ProjectStallLayout = {
      stalls: matchingStalls,
      labelPoint: project.centroid,
      total: matchingStalls.length,
      available,
      reserved,
      occupied: matchingStalls.length - available - reserved,
      floorId: effectiveFloorId,
      floorName: currentFloor?.name || 'ชั้น 1'
    }
    layoutCache.set(cacheKey, surveyLayout)
    return surveyLayout
  }

  let layout = layoutCache.get(cacheKey)
  if (!layout) {
    layout = buildProjectStallLayout(project, effectiveFloorId)
    layoutCache.set(cacheKey, layout)
  }
  return layout
}

export const getAllFloorLayoutsForProject = (project: ProjectZone): Map<string, ProjectStallLayout> => {
  const map = new Map<string, ProjectStallLayout>()
  const floors = project.floors && project.floors.length > 0
    ? project.floors
    : [{ id: '1', name: 'ชั้น 1', shortLabel: 'ชั้น 1', level: 1 }]
  floors.forEach(f => {
    map.set(f.id, getProjectStallLayout(project, f.id))
  })
  return map
}

export const getMarketStallSummary = (
  projects: ProjectZone[],
  statusOverrides: Record<string, StallStatus> = {},
  marketFull = false
): MarketStallSummary => {
  const summary: MarketStallSummary = { total: 0, available: 0, reserved: 0, occupied: 0 }
  const visibleProjects = projects.filter(project => !isOverviewBoundaryProject(project, projects))
  let minRent = Infinity
  let maxRent = 0
  let minArea = Infinity
  let maxArea = 0

  visibleProjects.forEach(project => {
    getAllFloorLayoutsForProject(project).forEach(layout => {
      layout.stalls.forEach(stall => {
        // A status embedded in the GIS layout is illustrative only. Runtime
        // availability comes from the booking API's explicit overrides.
        const status: StallStatus = marketFull
          ? 'occupied'
          : (statusOverrides[stall.code] || 'available')
        summary.total += 1
        summary[status] += 1
        const rent = Number(stall.monthlyRent) || 0
        if (rent > 0) {
          if (rent < minRent) minRent = rent
          if (rent > maxRent) maxRent = rent
        }
        const area = Number(stall.size) || 0
        if (area >= 1) {
          if (area < minArea) minArea = area
          if (area > maxArea) maxArea = area
        }
      })
    })
  })

  if (Number.isFinite(minRent) && maxRent > 0) {
    summary.minRent = minRent
    summary.maxRent = maxRent
  }
  if (Number.isFinite(minArea) && maxArea > 0) {
    summary.minAreaSqm = Number(minArea.toFixed(1))
    summary.maxAreaSqm = Number(maxArea.toFixed(1))
  }

  return summary
}
