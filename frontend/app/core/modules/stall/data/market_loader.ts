// ============================================================================
// Market GIS Data Loader (แยก 1 ตลาด = 1 ไฟล์ JSON)
// ============================================================================
// Loads individual market GIS layout files on-demand from /data/markets/<slug>.json
// with instant caching and fallback to static registry.
// ============================================================================

import { getMarketGisConfig, type MarketGisConfig } from './markets_registry'
import type { ProjectZone } from './market_gis_types'

export function cloneProjects<T>(items: T): T {
  if (!items) return items
  if (typeof structuredClone === 'function') return structuredClone(items)
  return JSON.parse(JSON.stringify(items))
}

/**
 * Merges base market projects (from /data/markets/<slug>.json) with saved GIS projects from API.
 * Prevents custom-drawn zones (zone-*) from wiping out the original market projects when
 * the initial draft was created before base JSON projects were loaded.
 */
function normalizeProjectNames(list: ProjectZone[] = []): ProjectZone[] {
  return list.map(p => ({
    ...p,
    name: String(p.name || '').replace(/\s*\(โซนแผงค้าตลาดสิงหา\)/g, '').trim()
  }))
}

export function mergeMarketProjects(baseProjects: ProjectZone[] = [], savedProjects: ProjectZone[] = [], clonePayload = true): ProjectZone[] {
  const copy = (projects: ProjectZone[]) => clonePayload ? cloneProjects(projects) : projects.slice()
  if (!Array.isArray(savedProjects) || savedProjects.length === 0) {
    return normalizeProjectNames(copy(baseProjects))
  }
  if (!Array.isArray(baseProjects) || baseProjects.length === 0) {
    return normalizeProjectNames(copy(savedProjects))
  }

  const coreBaseProjects = baseProjects.filter(p => !String(p.id).startsWith('zone-'))
  const baseIds = new Set(coreBaseProjects.map(p => p.id))
  const savedHasAnyBaseProject = savedProjects.some(p => baseIds.has(p.id))

  // If savedProjects only contains newly drawn custom zones (none of the base projects),
  // combine all coreBaseProjects with the custom zones so original market zones never disappear.
  if (!savedHasAnyBaseProject && coreBaseProjects.length > 0) {
    return normalizeProjectNames([...copy(coreBaseProjects), ...copy(savedProjects)])
  }

  return normalizeProjectNames(copy(savedProjects))
}

export const MARKET_SLUG_MAP: Record<string, string> = {
  '1': 'chatuchak',
  'chatuchak': 'chatuchak',
  '2': 'pracha-niwet-1',
  'pracha-niwet-1': 'pracha-niwet-1',
  '3': 'min-buri',
  'min-buri': 'min-buri',
  '4': 'singhanat',
  'singhanat': 'singhanat',
  '5': 'thewarat',
  'thewarat': 'thewarat',
  '6': 'rat-burana',
  'rat-burana': 'rat-burana',
  '7': 'bang-kapi',
  'bang-kapi': 'bang-kapi',
  '8': 'wongwian-lek',
  'wongwian-lek': 'wongwian-lek',
  '9': 'nong-chok',
  'nong-chok': 'nong-chok',
  '10': 'ratchadaphisek',
  'ratchadaphisek': 'ratchadaphisek',
  '11': 'thonburi-sanamluang-2',
  'thonburi-sanamluang-2': 'thonburi-sanamluang-2',
  '12': 'bang-khae-phirom',
  'bang-khae-phirom': 'bang-khae-phirom'
}

const marketCache = new Map<string, MarketGisConfig>()
// คำขอที่กำลังโหลดอยู่ ใช้ร่วมกันเมื่อหลายคอมโพเนนต์ขอตลาดเดียวกันพร้อมกัน
const marketInflight = new Map<string, Promise<MarketGisConfig>>()
// facilities ที่มากับ gis-config เดียวกัน (กันต้องโหลดก้อน ~7MB ซ้ำอีกรอบ)
const marketFacilities = new Map<string, any[]>()
export const getLoadedMarketFacilities = (marketIdOrSlug: string | number): any[] | undefined =>
  marketFacilities.get(getMarketSlug(marketIdOrSlug))

// ใช้ผลของคำขอเดียวกันร่วมกันภายในช่วงสั้นๆ (หน้าและแผนที่ขอ gis-config/draft ก้อนใหญ่พร้อมกัน)
const sharedRequests = new Map<string, { at: number; promise: Promise<any> }>()
export function sharedRequest<T>(key: string, factory: () => Promise<T>, ttlMs = 4000): Promise<T> {
  const hit = sharedRequests.get(key)
  if (hit && Date.now() - hit.at < ttlMs) return hit.promise
  const promise = factory()
  sharedRequests.set(key, { at: Date.now(), promise })
  promise.catch(() => { if (sharedRequests.get(key)?.promise === promise) sharedRequests.delete(key) })
  return promise
}

export function clearMarketGisCache(marketIdOrSlug?: string | number | null): void {
  if (!marketIdOrSlug) {
    marketCache.clear()
    marketInflight.clear()
    marketFacilities.clear()
    sharedRequests.clear()
    return
  }
  const slug = getMarketSlug(marketIdOrSlug)
  marketCache.delete(slug)
  marketInflight.delete(slug)
  marketFacilities.delete(slug)
}

/**
 * Resolve the slug filename for a market ID or slug
 */
export function getMarketSlug(marketIdOrSlug: string | number): string {
  const key = String(marketIdOrSlug).trim().toLowerCase()
  return MARKET_SLUG_MAP[key] || 'chatuchak'
}

function resolveBackendApiBase(): string {
  try {
    const config = useRuntimeConfig()
    const base = String(config?.public?.apiBase || config?.public?.apiUrl || '').trim()
    if (base) return base.replace(/\/+$/, '')
  } catch {}
  const envBase = String(
    (import.meta as any)?.env?.NUXT_PUBLIC_API_BASE ||
    (import.meta as any)?.env?.NUXT_PUBLIC_API_URL ||
    'http://localhost:8080/api/v1'
  ).trim()
  return envBase.replace(/\/+$/, '')
}

/**
 * Fetch and load a single market GIS layout on-demand from the Backend API
 * (`GET /api/v1/markets/<slug>/gis-config`) with immediate cache and fallback to static registry.
 */
export function loadMarketGisJson(marketIdOrSlug?: string | number | null, forceReload = false, signal?: AbortSignal): Promise<MarketGisConfig> {
  const slug = getMarketSlug(marketIdOrSlug || '1')

  if (!forceReload && marketCache.has(slug)) {
    return Promise.resolve(marketCache.get(slug)!)
  }
  // A page-scoped request must be independently abortable. Do not attach it to
  // the shared in-flight promise because aborting one hidden page would reject
  // every other consumer waiting for the same market.
  if (signal) return fetchMarketGisJson(slug, signal)
  if (!forceReload && marketInflight.has(slug)) {
    return marketInflight.get(slug)!
  }

  const pending = fetchMarketGisJson(slug).finally(() => {
    if (marketInflight.get(slug) === pending) marketInflight.delete(slug)
  })
  marketInflight.set(slug, pending)
  return pending
}

async function fetchMarketGisJson(slug: string, signal?: AbortSignal): Promise<MarketGisConfig> {
  const baseConfig = getMarketGisConfig(slug)

  // Fetch authoritative GIS projects & facilities from Backend API
  if (typeof window !== 'undefined' && typeof fetch === 'function') {
    try {
      const apiBase = resolveBackendApiBase()
      const res = await fetch(`${apiBase}/markets/${encodeURIComponent(slug)}/gis-config`, { cache: 'default', signal })
      if (res.ok) {
        const payload = await res.json()
        const gisData = payload?.data ?? payload
        marketFacilities.set(slug, Array.isArray(gisData?.facilities) ? gisData.facilities : [])
        const apiProjects = Array.isArray(gisData?.projects) ? gisData.projects : []
        if (apiProjects.length > 0) {
          const mergedConfig: MarketGisConfig = {
            ...baseConfig,
            // API payload is cached and treated as read-only by map/booking views.
            // Keep nested geometry by reference here; cloning a 10k-stall market
            // blocks the main thread for seconds. The GIS editor creates its own
            // mutable copy when it loads this cached config.
            projects: mergeMarketProjects(baseConfig.projects, apiProjects, false)
          }
          marketCache.set(slug, mergedConfig)
          return mergedConfig
        }
      }
    } catch (e) {
      if (signal?.aborted) throw e
      console.warn(`[market_loader] Failed to fetch GIS config from backend for ${slug}, using static fallback`, e)
    }
  }

  // Fallback to static registry
  marketCache.set(slug, baseConfig)
  return baseConfig
}
