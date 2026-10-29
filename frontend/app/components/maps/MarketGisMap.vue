<template>
  <div class="gis-map-wrapper">
    <!-- Map Header / Overlay Controls -->
    <div class="gis-overlay-header">
      <div class="gis-controls-left">
        <span class="gis-badge">
          <span class="badge-dot" style="background-color: var(--primary);"></span>
          <span>GIS Map: ผังแผงค้าในตลาด</span>
          <span class="gis-engine-pill">MapLibre GL</span>
        </span>

        <!-- Zone Selector -->
        <select v-model="selectedZoneId" class="gis-select" @change="applyFilters">
          <option value="">ทุกโซน (All Zones)</option>
          <option v-for="z in zones" :key="z.id" :value="z.id">
            {{ z.code }} - {{ z.name }}
          </option>
        </select>

        <!-- Only Available Checkbox -->
        <AppCheckbox
          v-model="onlyAvailable"
          label="เฉพาะแผงที่ว่างพร้อมจอง"
          @change="applyFilters"
        />
      </div>

      <div class="gis-controls-right">
        <!-- 2-Mode Segmented Control: แผนที่ vs ดาวเทียม -->
        <div class="gis-mode-segmented">
          <button
            type="button"
            class="gis-mode-tab"
            :class="{ active: currentMode === 'street' }"
            @click="setMapMode('street')"
            title="มุมมองแผนที่ (Street Map)"
          >
            แผนที่
          </button>
          <button
            type="button"
            class="gis-mode-tab"
            :class="{ active: currentMode === 'satellite' }"
            @click="setMapMode('satellite')"
            title="มุมมองภาพถ่ายดาวเทียม (Satellite Imagery)"
          >
            ดาวเทียม
          </button>
        </div>

      </div>
    </div>

    <!-- The MapLibre Canvas Container -->
    <div id="gis-map-canvas" class="map-canvas"></div>

    <!-- Selected Stall Drawer / Floating Popup Card -->
    <transition name="slide-up">
      <div v-if="selectedStall" class="gis-stall-card card">
        <div class="gis-card-header">
          <div>
            <div style="display: flex; align-items: center; gap: 8px;">
              <span class="stall-code" style="font-size: 1.35rem; color: var(--primary);">{{ selectedStall.stall_number }}</span>
              <span :class="getStatusBadgeClass(selectedStall.status)">
                {{ getStatusLabel(selectedStall.status) }}
              </span>
            </div>
            <div style="font-size: 0.8rem; color: var(--text-muted); margin-top: 2px;">
              โซน {{ selectedStall.zone_code }}: {{ selectedStall.category || selectedStall.zone_name }}
            </div>
          </div>
          <button class="btn-logout" @click="selectedStall = null" aria-label="ปิด">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>

        <div class="gis-card-body">
          <div class="gis-detail-grid">
            <div>
              <span class="gis-label">ขนาดพื้นที่:</span>
              <strong>{{ selectedStall.size_width || 3 }} × {{ selectedStall.size_length || 4 }} ม. ({{ selectedStall.area_sqm || selectedStall.size_sqm }} ตร.ม.)</strong>
            </div>
            <div>
              <span class="gis-label">ค่าเช่ารายเดือน:</span>
              <strong style="color: #059669; font-size: 1.05rem;">฿{{ Number(selectedStall.monthly_rate).toLocaleString() }}</strong>
            </div>
            <div>
              <span class="gis-label">ค่าเช่ารายวัน:</span>
              <strong>฿{{ Number(selectedStall.daily_rate).toLocaleString() }}</strong>
            </div>
            <div>
              <span class="gis-label">พิกัด GIS:</span>
              <span style="font-family: monospace; font-size: 0.78rem;">{{ Number(selectedStall.lat || 13.80).toFixed(5) }}, {{ Number(selectedStall.lng || 100.55).toFixed(5) }}</span>
            </div>
          </div>

          <div
            v-if="selectedStall.has_active_booking"
            class="gis-booking-alert"
            style="cursor: pointer; display: flex; align-items: center; justify-content: space-between;"
            @click="$emit('viewBooking', selectedStall)"
          >
            <span style="display: inline-flex; align-items: center; gap: 6px;">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
              <span>มีคำขอจองเข้าดูแผงนี้อยู่แล้ว (เลขที่: {{ selectedStall.active_booking_number }})</span>
            </span>
            <span style="text-decoration: underline; font-weight: 600; font-size: 0.78rem;">คลิกดู &rarr;</span>
          </div>

          <div class="gis-card-footer">
            <button
              v-if="selectedStall.status === 'available'"
              class="btn btn-primary"
              style="width: 100%; padding: 12px; font-size: 0.95rem; display: flex; align-items: center; justify-content: center; gap: 6px;"
              @click="$emit('bookStall', selectedStall)"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
              <span>กดจองดูแผงนี้ (Book Stall Visit)</span>
            </button>

            <!-- Button to view active booking details -->
            <button
              v-if="selectedStall.has_active_booking"
              class="btn btn-primary"
              style="width: 100%; padding: 11px; font-size: 0.92rem; background: linear-gradient(135deg, #8b5cf6 0%, #6d28d9 100%); margin-top: 6px; display: flex; align-items: center; justify-content: center; gap: 8px;"
              @click="$emit('viewBooking', selectedStall)"
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
              <span>ดูรายละเอียดการจอง ({{ selectedStall.active_booking_number }})</span>
            </button>

            <div v-if="selectedStall.status !== 'available' && !selectedStall.has_active_booking" style="text-align: center; font-size: 0.85rem; color: var(--text-dim); padding: 8px;">
              แผงนี้ไม่ว่างสำหรับจอง (สถานะ: {{ getStatusLabel(selectedStall.status) }})
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Map Legend -->
    <div class="gis-legend">
      <div class="legend-title">สัญลักษณ์สถานะแผง (GIS Legend)</div>
      <div class="legend-items">
        <div class="legend-item">
          <span class="legend-color" style="background-color: #10b981; border: 1px solid #059669;"></span>
          <span>ว่างพร้อมจอง (Available)</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #f59e0b; border: 1px solid #d97706;"></span>
          <span>มีผู้เช่าแล้ว (Occupied)</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #8b5cf6; border: 1px solid #7c3aed;"></span>
          <span>ติดจอง/รออนุมัติ (Pending)</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #f43f5e; border: 1px solid #e11d48;"></span>
          <span>ปิดปรับปรุง (Maintenance)</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { addStandardMapControls } from '~/utils/mapControls'
const { showWarning } = useAppAlert()
const props = defineProps<{
  zones: any[]
}>()

const emit = defineEmits<{
  (e: 'bookStall', stall: any): void
  (e: 'viewBooking', stall: any): void
}>()

const { get } = useApi()

let maplibregl: any = null
let map: any = null
let hoverPopup: any = null

const selectedZoneId = ref('')
const onlyAvailable = ref(false)
const currentMode = ref<'street' | 'satellite'>('street')
const selectedStall = ref<any>(null)
const geojsonData = ref<any>(null)

// Market Center: [lng, lat] for MapLibre
const marketCenter: [number, number] = [100.5504, 13.8000]

// Base Unified MapLibre Raster Style containing both Street & Satellite layers
const mapLibreStyle: any = {
  version: 8,
  sources: {
    'osm-street': {
      type: 'raster',
      tiles: [
        'https://tile.openstreetmap.org/{z}/{x}/{y}.png'
      ],
      tileSize: 256,
      attribution: '&copy; OpenStreetMap contributors',
      maxzoom: 19
    },
    'esri-satellite': {
      type: 'raster',
      tiles: [
        'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}'
      ],
      tileSize: 256,
      attribution: 'Tiles &copy; Esri'
    },
    'esri-labels': {
      type: 'raster',
      tiles: [
        'https://server.arcgisonline.com/ArcGIS/rest/services/Reference/World_Boundaries_and_Places/MapServer/tile/{z}/{y}/{x}'
      ],
      tileSize: 256,
      attribution: 'Tiles &copy; Esri'
    }
  },
  layers: [
    {
      id: 'esri-satellite-layer',
      type: 'raster',
      source: 'esri-satellite',
      layout: { visibility: 'none' },
      minzoom: 0,
      maxzoom: 20
    },
    {
      id: 'osm-street-layer',
      type: 'raster',
      source: 'osm-street',
      layout: { visibility: 'visible' },
      minzoom: 0,
      maxzoom: 19
    },
    {
      id: 'esri-labels-layer',
      type: 'raster',
      source: 'esri-labels',
      layout: { visibility: 'none' },
      minzoom: 0,
      maxzoom: 20
    }
  ]
}

onMounted(async () => {
  if (process.client) {
    maplibregl = await import('maplibre-gl')
    initMap()
    await loadStallsGeoJson()
  }
})

onBeforeUnmount(() => {
  if (map) {
    map.remove()
    map = null
  }
})

const initMap = () => {
  if (!maplibregl || !document.getElementById('gis-map-canvas')) return

  map = new maplibregl.Map({
    container: 'gis-map-canvas',
    style: mapLibreStyle,
    center: marketCenter,
    zoom: 17.5,
    minZoom: 14,
    maxZoom: 21,
    pitch: 25
  })

  // Add zoom and pitch navigation controls
  addStandardMapControls(maplibregl, map, {
    onFit: () => recenterMap(),
    fitTitle: 'กลับสู่จุดศูนย์กลางตลาด',
    notify: showWarning
  })

  hoverPopup = new maplibregl.Popup({
    closeButton: false,
    closeOnClick: false,
    offset: 12
  })

  map.on('load', () => {
    // Add market perimeter boundary
    map.addSource('market-bounds-src', {
      type: 'geojson',
      data: {
        type: 'Feature',
        geometry: {
          type: 'Polygon',
          coordinates: [[
            [100.5496, 13.7993],
            [100.5518, 13.7993],
            [100.5518, 13.8010],
            [100.5496, 13.8010],
            [100.5496, 13.7993]
          ]]
        }
      }
    })

    map.addLayer({
      id: 'market-bounds-fill',
      type: 'fill',
      source: 'market-bounds-src',
      paint: {
        'fill-color': '#2563eb',
        'fill-opacity': 0.05
      }
    })

    map.addLayer({
      id: 'market-bounds-line',
      type: 'line',
      source: 'market-bounds-src',
      paint: {
        'line-color': '#2563eb',
        'line-width': 2,
        'line-dasharray': [4, 4]
      }
    })

    // Setup stalls source and layers
    setupStallsLayers()
  })
}

const setupStallsLayers = () => {
  if (!map) return

  const filtered = getFilteredGeoJson()

  if (!map.getSource('stalls-src')) {
    map.addSource('stalls-src', {
      type: 'geojson',
      data: filtered
    })

    // Stalls Fill Layer
    map.addLayer({
      id: 'stalls-fill',
      type: 'fill',
      source: 'stalls-src',
      paint: {
        'fill-color': [
          'case',
          ['boolean', ['get', 'has_active_booking'], false], '#8b5cf6',
          ['==', ['get', 'status'], 'available'], '#10b981',
          ['==', ['get', 'status'], 'occupied'], '#f59e0b',
          ['==', ['get', 'status'], 'reserved'], '#8b5cf6',
          ['==', ['get', 'status'], 'maintenance'], '#f43f5e',
          '#64748b'
        ],
        'fill-opacity': 0.82
      }
    })

    // Stalls Outline Layer
    map.addLayer({
      id: 'stalls-line',
      type: 'line',
      source: 'stalls-src',
      paint: {
        'line-color': '#ffffff',
        'line-width': 1.5
      }
    })

    // Interactive Hover on Stalls
    map.on('mousemove', 'stalls-fill', (e: any) => {
      if (!e.features || !e.features.length) return
      map.getCanvas().style.cursor = 'pointer'
      const f = e.features[0]
      const p = f.properties

      const coordinates = e.lngLat
      const statusLabel = getStatusLabel(p.status)

      hoverPopup.setLngLat(coordinates).setHTML(`
        <div style="font-family: 'Bai Jamjuree', sans-serif; padding: 4px; font-size: 12px; line-height: 1.4;">
          <div style="font-weight: 700; color: #0f172a; font-size: 13px;">
            แผง ${p.stall_number} <span style="font-weight: 500; font-size: 11px; color: #64748b;">(โซน ${p.zone_code})</span>
          </div>
          <div style="color: #475569; font-size: 11px; margin-bottom: 2px;">
            ${p.category || p.zone_name || ''}
          </div>
          <div style="display: flex; justify-content: space-between; gap: 8px; margin-top: 4px;">
            <span style="font-weight: 700; color: #059669;">฿${Number(p.monthly_rate).toLocaleString()}/ด.</span>
            <span style="font-size: 10px; padding: 1px 6px; border-radius: 4px; background: #f1f5f9; color: #334155;">${statusLabel}</span>
          </div>
        </div>
      `).addTo(map)
    })

    map.on('mouseleave', 'stalls-fill', () => {
      map.getCanvas().style.cursor = ''
      hoverPopup.remove()
    })

    // Click on Stall
    map.on('click', 'stalls-fill', (e: any) => {
      if (!e.features || !e.features.length) return
      const p = e.features[0].properties
      selectedStall.value = p
      if (p.lng && p.lat) {
        map.flyTo({
          center: [Number(p.lng), Number(p.lat)],
          zoom: 18.5,
          speed: 1.2
        })
      }
    })
  } else {
    (map.getSource('stalls-src') as any).setData(filtered)
  }
}

const getFilteredGeoJson = () => {
  if (!geojsonData.value) {
    return { type: 'FeatureCollection', features: [] }
  }

  let features = geojsonData.value.features || []
  if (selectedZoneId.value) {
    features = features.filter((f: any) => String(f.properties.zone_id) === String(selectedZoneId.value))
  }
  if (onlyAvailable.value) {
    features = features.filter((f: any) => f.properties.status === 'available')
  }

  return {
    type: 'FeatureCollection',
    features
  }
}

const loadStallsGeoJson = async () => {
  try {
    const data = await get('/stalls/gis-geojson')
    geojsonData.value = data
    if (map && map.isStyleLoaded()) {
      setupStallsLayers()
    }
  } catch (err) {
    console.error('Failed to load GIS GeoJSON:', err)
  }
}

const setMapMode = (mode: 'street' | 'satellite') => {
  if (currentMode.value === mode || !map) return
  currentMode.value = mode

  if (mode === 'street') {
    map.setLayoutProperty('osm-street-layer', 'visibility', 'visible')
    map.setLayoutProperty('esri-satellite-layer', 'visibility', 'none')
    map.setLayoutProperty('esri-labels-layer', 'visibility', 'none')
  } else {
    map.setLayoutProperty('osm-street-layer', 'visibility', 'none')
    map.setLayoutProperty('esri-satellite-layer', 'visibility', 'visible')
    map.setLayoutProperty('esri-labels-layer', 'visibility', 'visible')
  }
}

const applyFilters = () => {
  if (map && map.getSource('stalls-src')) {
    (map.getSource('stalls-src') as any).setData(getFilteredGeoJson())
  }
}

const recenterMap = () => {
  if (map) {
    map.flyTo({
      center: marketCenter,
      zoom: 17.5,
      pitch: 25,
      bearing: 0,
      speed: 1.2
    })
  }
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'available': return 'badge badge-available'
    case 'occupied': return 'badge badge-occupied'
    case 'reserved': return 'badge badge-reserved'
    case 'maintenance': return 'badge badge-maintenance'
    default: return 'badge'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'available': return 'ว่างพร้อมจอง'
    case 'occupied': return 'มีผู้เช่าแล้ว'
    case 'reserved': return 'ติดจอง'
    case 'maintenance': return 'ปิดปรับปรุง'
    default: return status || 'ว่าง'
  }
}

defineExpose({
  loadStallsGeoJson
})
</script>

<style scoped>
.gis-map-wrapper {
  position: relative;
  width: 100%;
  height: 640px;
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--border-subtle);
  box-shadow: var(--shadow-lg);
}

.map-canvas {
  width: 100%;
  height: 100%;
  background: #0f172a;
  position: relative;
}

/* Overlay Header Controls */
.gis-overlay-header {
  position: absolute;
  top: 16px;
  left: 16px;
  right: 16px;
  z-index: 10;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  pointer-events: none;
}

.gis-controls-left, .gis-controls-right {
  display: flex;
  align-items: center;
  gap: 10px;
  pointer-events: auto;
}

.gis-badge {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(12px);
  padding: 8px 16px;
  border-radius: var(--radius-full);
  font-size: 0.85rem;
  font-weight: 600;
  border: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  gap: 8px;
  color: #0f172a;
  box-shadow: var(--shadow-md);
}

/* 2-Mode Segmented Control (Pill Style) */
.gis-mode-segmented {
  display: inline-flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(8px);
  padding: 3px;
  border-radius: 9999px;
  border: 1px solid rgba(226, 232, 240, 0.9);
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.12);
}

.gis-mode-tab {
  padding: 6px 16px;
  border-radius: 9999px;
  border: none;
  background: transparent;
  font-family: inherit;
  font-size: 0.84rem;
  font-weight: 600;
  color: #64748b;
  cursor: pointer;
  transition: all 0.18s ease;
  line-height: 1.2;
}

.gis-mode-tab:hover {
  color: #1e293b;
}

.gis-mode-tab.active {
  background: #ffffff;
  color: #2563eb;
  font-weight: 700;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.2);
}

.gis-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(8px);
  padding: 6px 14px;
  border-radius: 9999px;
  border: 1px solid rgba(226, 232, 240, 0.9);
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.12);
  font-family: inherit;
  font-size: 0.84rem;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  transition: all 0.18s ease;
  line-height: 1.2;
}

.gis-btn svg {
  color: #2563eb;
}

.gis-btn:hover {
  background: #ffffff;
  color: #1e293b;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.16);
  transform: translateY(-1px);
}

/* Floating Stall Details Card */
.gis-stall-card {
  position: absolute;
  bottom: 20px;
  left: 20px;
  width: 380px;
  max-width: calc(100% - 40px);
  z-index: 20;
  background: #ffffff;
  backdrop-filter: blur(16px);
  border: 1px solid var(--border-card);
  box-shadow: var(--shadow-popup);
  padding: 20px;
  border-radius: var(--radius-xl);
  color: #0f172a;
}

.gis-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 14px;
  border-bottom: 1px solid #f1f5f9;
  padding-bottom: 10px;
}

.gis-detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  font-size: 0.85rem;
  margin-bottom: 14px;
}

.gis-label {
  display: block;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.gis-booking-alert {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #2563eb;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  font-size: 0.78rem;
  margin-bottom: 12px;
}

/* Legend */
.gis-legend {
  position: absolute;
  bottom: 20px;
  right: 20px;
  z-index: 10;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(12px);
  border: 1px solid var(--border-subtle);
  padding: 12px 16px;
  border-radius: var(--radius-md);
  font-size: 0.78rem;
  box-shadow: var(--shadow-md);
  color: #0f172a;
}

.legend-title {
  font-weight: 700;
  margin-bottom: 8px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.legend-items {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #334155;
}

.legend-color {
  width: 14px;
  height: 14px;
  border-radius: 3px;
  display: inline-block;
}

/* Transitions */
.slide-up-enter-active, .slide-up-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-up-enter-from, .slide-up-leave-to {
  transform: translateY(20px);
  opacity: 0;
}

@media (max-width: 768px) {
  .gis-map-wrapper {
    height: 520px;
  }
  .gis-overlay-header {
    position: static;
    padding: 10px 12px;
    background: #ffffff;
    border-bottom: 1px solid var(--border-subtle);
  }
  .gis-controls-left, .gis-controls-right {
    width: 100%;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 8px;
  }
  .gis-stall-card {
    position: fixed !important;
    bottom: 0 !important;
    left: 0 !important;
    right: 0 !important;
    width: 100% !important;
    max-width: 100% !important;
    max-height: 80vh !important;
    border-radius: 20px 20px 0 0 !important;
    z-index: 99 !important;
    box-shadow: 0 -10px 30px rgba(15, 23, 42, 0.25) !important;
  }
  .gis-legend {
    display: none;
  }
}
</style>
