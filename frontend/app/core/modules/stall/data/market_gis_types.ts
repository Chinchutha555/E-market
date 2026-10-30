export interface ProjectFloor {
  id: string
  name: string
  shortLabel: string
  level: number
  description?: string
  isFull?: boolean
}

export interface ProjectZone {
  id: string
  name: string
  sectionNumber: string
  category: string
  color: string
  stallsAvailable: number
  totalStalls: number
  realAreaSqm: number
  centroid: [number, number]
  coordinates: [number, number][]
  floors?: ProjectFloor[]
  stalls?: any[]
  /** Manual projects start empty; staff draw each stall on the selected floor. */
  layoutMode?: 'manual' | 'generated'
}

export interface MarketFacility {
  id: string
  name: string
  type: string
  coordinates: [number, number]
  geometry?: {
    type: 'Polygon'
    coordinates: [number, number][][]
  }
  description?: string
  projectId?: string
  floorId?: string
  number?: string
  sourceZone?: string
  sourceProperties?: Record<string, unknown>
}
