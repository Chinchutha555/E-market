<template>
  <div>
    <AppHeader />
    <main class="content-body stall-detail-page">
      <AppBreadcrumb
        :market="marketName"
        :market-id="market"
        :items="breadcrumbItems"
        :current="stallCode ? `รายละเอียดแผง ${stallCode}` : 'รายละเอียดแผงเช่า'"
        :back-to="backToManageUrl"
      />

      <template v-if="loading">
        <section class="detail-header skeleton-card"><span></span><span></span><span></span></section>
        <section class="detail-card skeleton-card tall"></section>
      </template>

      <template v-else-if="stall">
        <section class="detail-hero-card">
          <div class="detail-gallery-col">
            <div class="detail-main-photo">
              <img :src="stallPhotos[activePhotoIndex]" :alt="`ภาพแผง ${stall.ident || stall.stall_ident}`" />
              <span class="detail-photo-badge">{{ activePhotoIndex + 1 }} / {{ stallPhotos.length }} รูป</span>
            </div>
            <div class="detail-thumb-row">
              <button
                v-for="(src, idx) in stallPhotos"
                :key="idx"
                type="button"
                class="detail-thumb-btn"
                :class="{ active: activePhotoIndex === idx }"
                @click="activePhotoIndex = idx"
              >
                <img :src="src" :alt="`รูปตัวอย่างที่ ${idx + 1}`" />
              </button>
            </div>
          </div>

          <div class="detail-summary-col">
            <div class="detail-summary-top">
              <div>
                <div class="detail-eyebrow-row">
                  <span class="eyebrow">{{ marketName }}</span>
                  <span class="status-badge" :class="`status-${stall.stall_status || stall.status || 'available'}`">{{ statusLabel }}</span>
                </div>
                <h1>แผง {{ stallShortNumber(stall.ident || stall.stall_ident) }}{{ stall.name ? ` (${stall.name})` : '' }}</h1>
                <p>{{ stall.project_name || 'ไม่ระบุโครงการ' }} · โซน {{ stall.zone_name || stall.zone_id || '-' }} · {{ formattedFloor }}</p>
              </div>
              <div class="header-actions">
                <AppButton variant="outline" :to="backToManageUrl">
                  {{ route.query.mode === 'inspect' ? 'กลับไปผังแผงเช่า' : 'กลับไปจัดการผัง' }}
                </AppButton>
                <AppButton variant="primary" :to="`/contracts?stall=${encodeURIComponent(stall.ident || stall.stall_ident)}`">
                  จัดการสัญญา
                </AppButton>
              </div>
            </div>

            <div class="detail-kpi-strip">
              <div class="detail-kpi-box">
                <span>ขนาดพื้นที่</span>
                <strong>{{ stallAreaSqm }} <small>ตร.ม.</small></strong>
                <em>กว้าง {{ stall.width || 2 }} × ลึก {{ stall.depth || 2.5 }} ม.</em>
              </div>
              <div class="detail-kpi-box is-green">
                <span>ค่าเช่ารายเดือน</span>
                <strong>฿{{ Number(stall.monthly_rent || 1800).toLocaleString('th-TH') }}</strong>
                <em>ต่อเดือน</em>
              </div>
              <div class="detail-kpi-box">
                <span>ค่าเช่ารายวัน</span>
                <strong>฿{{ Number(stall.daily_rent || 150).toLocaleString('th-TH') }}</strong>
                <em>ต่อวัน</em>
              </div>
              <div class="detail-kpi-box">
                <span>ค่าธรรมเนียมแรกเข้า</span>
                <strong>฿{{ Number(stallEntryFee).toLocaleString('th-TH') }}</strong>
                <em>ชำระเมื่อทำสัญญา</em>
              </div>
            </div>
          </div>
        </section>

        <!-- Right-Edge History Side Panel Trigger (History Panel Standard) -->
        <button
          type="button"
          class="timeline-edge-trigger"
          :class="{ active: showHistoryPanel }"
          :aria-expanded="showHistoryPanel"
          aria-controls="stall-history-side-panel"
          @click="showHistoryPanel = !showHistoryPanel"
        >
          <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 3-6.7L3 8"/><path d="M3 3v5h5"/><path d="M12 7v5l3 2"/></svg>
          <span>ประวัติการแก้ไข</span>
        </button>

        <Transition name="timeline-slide">
          <div v-if="showHistoryPanel" class="timeline-backdrop" @click.self="showHistoryPanel = false">
            <aside id="stall-history-side-panel" class="timeline-panel" role="dialog" aria-modal="true" aria-label="ประวัติการแก้ไขข้อมูลแผง">
              <header class="timeline-header">
                <div>
                  <span class="eyebrow">รายการกิจกรรมทั้งหมด</span>
                  <h2>ประวัติการแก้ไข</h2>
                </div>
                <button type="button" class="close-button" aria-label="ปิดประวัติการแก้ไข" @click="showHistoryPanel = false">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                </button>
              </header>
              <div class="timeline-body">
                <div v-if="history.length" class="timeline-list">
                  <div v-for="item in history" :key="item.id" class="timeline-item">
                    <span class="timeline-dot"></span>
                    <div class="timeline-copy">
                      <strong>{{ item.actor_name || 'ระบบ' }}</strong>
                      <time>{{ formatDate(item.created_at) }}</time>
                      <p>{{ item.changes }}</p>
                    </div>
                  </div>
                </div>
                <AppEmptyState v-else variant="compact" icon="history" title="ยังไม่มีประวัติการแก้ไข" description="การเปลี่ยนแปลงข้อมูลแผงจะแสดงที่นี่" />
              </div>
            </aside>
          </div>
        </Transition>

        <!-- Main Card with Consolidated Tabs -->
        <section class="detail-card">
          <AppPrimaryTabs
            v-model="activeTab"
            :options="detailPrimaryTabs"
            aria-label="ข้อมูลแผง การจอง สัญญา และประวัติ"
          />

          <!-- Tab 1: ข้อมูลแผง -->
          <div v-if="activeTab === 'info'" class="tab-panel">
            <div class="panel-heading">
              <div>
                <h2>ข้อมูลโครงสร้างและรายละเอียดแผงเช่า</h2>
                <p>รายละเอียดโครงสร้าง ขนาดพื้นที่ หมวดสินค้า และอัตราค่าธรรมเนียมของแผง {{ stallShortNumber(stall.ident || stall.stall_ident) }}</p>
              </div>
            </div>
            <dl class="fact-grid">
              <div><dt>รหัสแผงค้า</dt><dd>{{ stallShortNumber(stall.ident || stall.stall_ident) }}</dd></div>
              <div><dt>โครงการ / อาคาร</dt><dd>{{ stall.project_name || 'ไม่ระบุ' }}</dd></div>
              <div><dt>โซน / ระดับชั้น</dt><dd>โซน {{ stall.zone_name || stall.zone_id || '-' }} / {{ formattedFloor }}</dd></div>
              <div><dt>ขนาดกว้าง × ลึก</dt><dd>{{ stall.width || 2 }} × {{ stall.depth || 2.5 }} ม. ({{ stallAreaSqm }} ตร.ม.)</dd></div>
              <div><dt>หมวดสินค้า</dt><dd>{{ stall.category || 'ยังไม่ระบุ' }}</dd></div>
              <div><dt>สถานะปัจจุบัน</dt><dd><span class="status-badge" :class="`status-${stall.stall_status || stall.status || 'available'}`">{{ statusLabel }}</span></dd></div>
              <div><dt>ค่าเช่ารายวัน</dt><dd>฿{{ Number(stall.daily_rent || 150).toLocaleString('th-TH') }} / วัน</dd></div>
              <div><dt>ค่าเช่ารายเดือน</dt><dd>฿{{ Number(stall.monthly_rent || 1800).toLocaleString('th-TH') }} / เดือน</dd></div>
              <div><dt>ค่าธรรมเนียมแรกเข้า</dt><dd>฿{{ Number(stallEntryFee).toLocaleString('th-TH') }}</dd></div>
              <div class="full"><dt>รายละเอียดเพิ่มเติม</dt><dd>{{ stall.description || 'แผงค้ามาตรฐานพร้อมระบบไฟฟ้าแสงสว่างและทางเดินส่วนกลาง' }}</dd></div>
            </dl>
          </div>

          <!-- Tab 2: การจองและคำร้องขอเช่า -->
          <div v-else-if="activeTab === 'bookings'" class="tab-panel">
            <div class="card-heading">
              <div>
                <h2>รายการจองนัดเข้าชมและคำร้องขอเช่าของแผงนี้</h2>
                <p>{{ stallBookingsAndApps.length }} รายการ</p>
              </div>
              <div class="header-actions">
                <AppButton size="sm" variant="outline" to="/bookings">ไปหน้าการจอง</AppButton>
                <AppButton size="sm" variant="outline" to="/applications">ไปหน้าคำร้อง</AppButton>
              </div>
            </div>

            <div v-if="stallBookingsAndApps.length" class="table-wrap">
              <table class="contracts-table">
                <thead>
                  <tr>
                    <th>ประเภท</th>
                    <th>เลขที่อ้างอิง</th>
                    <th>ผู้ยื่นรายการ</th>
                    <th>วันที่นัดหมาย / วันที่ยื่น</th>
                    <th>สถานะ</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in stallBookingsAndApps" :key="item.key">
                    <td><strong>{{ item.typeLabel }}</strong></td>
                    <td>{{ item.refNumber }}</td>
                    <td>{{ item.applicantName }}</td>
                    <td>{{ item.dateText }}</td>
                    <td><span class="status-badge" :class="item.badgeClass">{{ item.statusText }}</span></td>
                  </tr>
                </tbody>
              </table>
            </div>
            <AppEmptyState
              v-else
              variant="compact"
              icon="calendar"
              title="ยังไม่มีรายการจองหรือคำร้อง"
              description="รายการนัดดูพื้นที่และคำร้องขอเช่าของแผงนี้จะแสดงที่นี่"
            />
          </div>

          <!-- Tab 3: สัญญาของแผงนี้ -->
          <div v-else-if="activeTab === 'contracts'" class="tab-panel">
            <div class="card-heading">
              <div>
                <h2>สัญญาของแผงนี้</h2>
                <p>{{ contracts.length }} รายการ</p>
              </div>
              <AppButton size="sm" variant="outline" :to="`/contracts?stall=${encodeURIComponent(stall.ident || stall.stall_ident)}`">
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                </template>
                จัดการสัญญา
              </AppButton>
            </div>

            <div v-if="contracts.length" class="table-wrap">
              <table class="contracts-table">
                <thead>
                  <tr>
                    <th>เลขที่สัญญา</th>
                    <th>ผู้เช่า (Owner)</th>
                    <th>ชื่อร้านค้า</th>
                    <th>ระยะเวลาสัญญา</th>
                    <th>สถานะ</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="contract in contracts" :key="contract.id">
                    <td><strong>{{ contract.contract_number }}</strong></td>
                    <td>{{ contract.tenant_name || '-' }}</td>
                    <td>{{ contract.shop_name || 'ไม่ระบุชื่อร้าน' }}</td>
                    <td>{{ contract.start_date }} – {{ contract.end_date }}</td>
                    <td><span class="status-badge status-available">{{ contract.status || 'มีผลบังคับใช้' }}</span></td>
                  </tr>
                </tbody>
              </table>
            </div>
            <AppEmptyState
              v-else
              variant="compact"
              icon="document"
              title="ยังไม่มีสัญญาของแผงนี้"
              description="สัญญาที่เชื่อมกับแผงค้าจะแสดงที่นี่"
            />
          </div>

          <!-- Tab 4: ประวัติการแก้ไข -->
          <div v-else class="tab-panel">
            <div class="card-heading">
              <div>
                <h2>ประวัติการแก้ไขข้อมูลแผง</h2>
                <p>{{ history.length }} รายการ</p>
              </div>
            </div>
            <div v-if="history.length" class="timeline-list">
              <div v-for="item in history" :key="item.id" class="timeline-item">
                <span class="timeline-dot"></span>
                <div class="timeline-copy">
                  <strong>{{ item.actor_name || 'ระบบ' }}</strong>
                  <time>{{ formatDate(item.created_at) }}</time>
                  <p>{{ item.changes }}</p>
                </div>
              </div>
            </div>
            <AppEmptyState
              v-else
              variant="compact"
              icon="history"
              title="ยังไม่มีประวัติการแก้ไข"
              description="การเปลี่ยนแปลงข้อมูลแผงจะถูกบันทึกไว้ที่นี่"
            />
          </div>
        </section>
      </template>

      <section v-else class="not-found">
        <AppEmptyState variant="page" icon="alert" title="ไม่พบข้อมูลแผง" description="กลับไปเลือกแผงจากหน้าจัดการผังอีกครั้ง">
          <template #actions>
            <AppButton :to="backToManageUrl">{{ route.query.mode === 'inspect' ? 'กลับไปดูข้อมูลแผง (GIS)' : 'กลับไปจัดการผัง' }}</AppButton>
          </template>
        </AppEmptyState>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import AppButton from '~/components/AppButton.vue'
import { getMarketGisConfig } from '../data/markets_registry'
import { loadMarketGisJson, mergeMarketProjects } from '../data/market_loader'
import { getAllFloorLayoutsForProject } from '../data/market_stall_layout'

const route = useRoute()
const { get } = useApi()
const market = computed(() => String(route.query.market || '1'))
const marketName = computed(() => getMarketGisConfig(market.value).name)
const stallCode = computed(() => String(route.query.stall || '').trim())
const backToManageUrl = computed(() => {
  const params = new URLSearchParams({ market: market.value })
  if (route.query.project) params.set('project', String(route.query.project))
  if (stallCode.value) params.set('stall', stallCode.value)
  if (route.query.mode === 'inspect') {
    params.set('mode', 'inspect')
    params.set('view', 'map')
  }
  return `/stalls/manage?${params.toString()}`
})
const breadcrumbItems = computed(() => {
  const modeBaseParams = new URLSearchParams({ market: market.value })
  if (route.query.mode === 'inspect') {
    modeBaseParams.set('mode', 'inspect')
    modeBaseParams.set('view', 'map')
  }
  const modeOverviewUrl = `/stalls/manage?${modeBaseParams.toString()}`

  const items: Array<{ label: string; to?: string }> = [
    {
      label: route.query.mode === 'inspect' ? 'ผังแผงเช่า' : 'จัดการแผงเช่า',
      to: modeOverviewUrl
    }
  ]
  if (stall.value?.project_name) {
    items.push({
      label: String(stall.value.project_name),
      to: backToManageUrl.value
    })
  }
  return items
})
const loading = ref(true)
const activeTab = ref<'info' | 'bookings' | 'contracts' | 'history'>('info')
const showHistoryPanel = ref(false)
const activePhotoIndex = ref(0)
const stall = ref<any>(null)
const contracts = ref<any[]>([])
const history = ref<any[]>([])
const bookingsList = ref<any[]>([])
const applicationsList = ref<any[]>([])

const DEFAULT_DETAIL_PHOTOS = [
  '/images/stall_3d_perspective.jpg',
  '/images/market_stall_sample.jpg',
  '/images/market_hall_front.jpg'
]
const stallPhotos = computed(() => {
  const custom = stall.value?.photos
  if (Array.isArray(custom) && custom.length > 0) return custom
  return DEFAULT_DETAIL_PHOTOS
})
const stallAreaSqm = computed(() => {
  const w = Number(stall.value?.width) || 2
  const d = Number(stall.value?.depth) || 2.5
  const area = Number((w * d).toFixed(1))
  return Number.isInteger(area) ? String(area) : area.toFixed(1)
})
const stallEntryFee = computed(() => {
  const explicit = Number(stall.value?.entry_fee || stall.value?.entryFee || 0)
  if (explicit > 0) return explicit
  const rent = Number(stall.value?.monthly_rent) || 1800
  return rent * 20
})
const stallBookingsAndApps = computed(() => {
  const target = stallCode.value.toUpperCase()
  if (!target) return []
  const items: Array<{
    key: string
    typeLabel: string
    refNumber: string
    applicantName: string
    dateText: string
    statusText: string
    badgeClass: string
  }> = []

  bookingsList.value.forEach((b: any) => {
    const codes = String(b.stall_number || b.stall_ident || '')
      .split(',')
      .map(s => s.trim().toUpperCase())
    if (!codes.includes(target)) return
    const st = String(b.status || '').toLowerCase()
    const statusMap: Record<string, { label: string; cls: string }> = {
      pending: { label: 'รอตรวจสอบวันนัด', cls: 'status-reserved' },
      approved: { label: 'นัดหมายสำเร็จ', cls: 'status-available' },
      interested: { label: 'สนใจเช่าแผง', cls: 'status-available' },
      waiting_queue: { label: 'คิวสำรอง', cls: 'status-reserved' },
      cancelled: { label: 'ยกเลิก', cls: 'status-occupied' },
      rejected: { label: 'ปฏิเสธ', cls: 'status-occupied' }
    }
    const mapped = statusMap[st] || { label: b.status || 'รอดำเนินการ', cls: 'status-reserved' }
    items.push({
      key: `booking-${b.id || b.booking_number}`,
      typeLabel: 'นัดเข้าชมพื้นที่',
      refNumber: b.booking_number || `#${b.id}`,
      applicantName: b.tenant_name || b.applicant_name || '-',
      dateText: b.visit_date ? `${b.visit_date}${b.visit_time_slot ? ' · ' + b.visit_time_slot : ''}` : formatDate(b.created_at),
      statusText: mapped.label,
      badgeClass: mapped.cls
    })
  })

  applicationsList.value.forEach((a: any) => {
    const codes = String(a.stall_number || a.stall_ident || '')
      .split(',')
      .map(s => s.trim().toUpperCase())
    if (!codes.includes(target)) return
    const st = String(a.status || '').toLowerCase()
    const isApproved = st === 'approved'
    const isRejected = st === 'rejected'
    items.push({
      key: `app-${a.id || a.application_number}`,
      typeLabel: 'คำร้องขอเช่าแผง',
      refNumber: a.application_number || `#${a.id}`,
      applicantName: a.applicant_name || '-',
      dateText: formatDate(a.created_at),
      statusText: isApproved ? 'อนุมัติแล้ว' : (isRejected ? 'ไม่อนุมัติ' : 'รอพิจารณาอนุมัติ'),
      badgeClass: isApproved ? 'status-available' : (isRejected ? 'status-occupied' : 'status-reserved')
    })
  })

  return items
})
const detailPrimaryTabs = computed(() => [
  { value: 'info', label: 'ข้อมูลแผงเช่า', icon: 'info' as const },
  { value: 'bookings', label: 'การจองและคำร้องขอเช่า', icon: 'clock' as const, count: stallBookingsAndApps.value.length },
  { value: 'contracts', label: 'สัญญาของแผงนี้', icon: 'contract' as const, count: contracts.value.length },
  { value: 'history', label: 'ประวัติการแก้ไข', icon: 'history' as const, count: history.value.length }
])

const formattedFloor = computed(() => {
  const raw = String(stall.value?.floor || '').trim()
  if (!raw) return 'ชั้น 1'
  return raw.startsWith('ชั้น') ? raw : `ชั้น ${raw}`
})
const statusLabel = computed(() => ({ available: 'แผงว่าง', reserved: 'ติดจอง', occupied: 'มีผู้เช่า' }[stall.value?.stall_status || stall.value?.status || 'available'] || 'ไม่ระบุสถานะ'))
const unwrap = (value: any) => value?.data ?? value
const formatDate = (value?: string) => value ? new Intl.DateTimeFormat('th-TH', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) : '-'

const loadDetail = async () => {
  loading.value = true
  try {
    const [stallsResult, baseMarket, savedConfigResult, bookingsRes, appsRes] = await Promise.all([
      get('/stalls').catch(() => []),
      loadMarketGisJson(market.value),
      get(`/markets/${encodeURIComponent(market.value)}/gis-config/draft`).catch(() => null),
      get('/bookings').catch(() => []),
      get('/applications').catch(() => [])
    ])
    bookingsList.value = Array.isArray(unwrap(bookingsRes)) ? unwrap(bookingsRes) : []
    applicationsList.value = Array.isArray(unwrap(appsRes)) ? unwrap(appsRes) : []

    const result = unwrap(stallsResult)
    const list = Array.isArray(result) ? result : []
    stall.value = list.find((item: any) => String(item.ident || item.stall_ident).toUpperCase() === stallCode.value.toUpperCase()) || null
    if (!stall.value) {
      const savedConfig = unwrap(savedConfigResult) || {}
      const projects = mergeMarketProjects(baseMarket.projects || [], Array.isArray(savedConfig.projects) ? savedConfig.projects : [])
      for (const project of projects) {
        const directStall = (Array.isArray(project.stalls) ? project.stalls : []).find((item: any) => String(item.code || item.ident).toUpperCase() === stallCode.value.toUpperCase())
        const generatedStall = Array.from(getAllFloorLayoutsForProject(project).values()).flatMap(layout => layout.stalls)
          .find(item => String(item.code).toUpperCase() === stallCode.value.toUpperCase())
        const mapStall = directStall || generatedStall
        if (!mapStall) continue
        stall.value = {
          ident: mapStall.code,
          name: '',
          project_name: project.name,
          zone_name: project.sectionNumber,
          floor: mapStall.floorName || mapStall.floor || '1',
          width: mapStall.widthM || 2,
          depth: mapStall.depthM || 2.5,
          monthly_rent: mapStall.monthlyRent || 1800,
          daily_rent: 150,
          entry_fee: mapStall.entryFee || (Number(mapStall.monthlyRent || 1800) * 20),
          category: mapStall.areaType || project.category || 'แผงค้ามาตรฐาน',
          description: '',
          stall_status: mapStall.status || 'available'
        }
        break
      }
    }
    if (!stall.value) return
    const targetIdent = String(stall.value.ident || stall.value.stall_ident || stallCode.value)
    const [contractsResult, historyResult] = await Promise.allSettled([
      get(`/contracts?stall_ident=${encodeURIComponent(targetIdent)}`),
      stall.value.id ? get(`/stalls/${stall.value.id}/history`) : Promise.resolve([])
    ])
    contracts.value = contractsResult.status === 'fulfilled' && Array.isArray(unwrap(contractsResult.value)) ? unwrap(contractsResult.value) : []
    history.value = historyResult.status === 'fulfilled' && Array.isArray(unwrap(historyResult.value)) ? unwrap(historyResult.value) : []
  } finally {
    loading.value = false
  }
}

watch(stallCode, loadDetail, { immediate: true })
</script>

<style scoped>
.stall-detail-page{padding-bottom:48px}.detail-header,.detail-hero-card,.detail-card,.not-found{background:var(--bg-ffffff);border:1px solid var(--bd-eef2f6);border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.03)}.detail-hero-card{display:grid;grid-template-columns:280px minmax(0,1fr);gap:24px;margin:18px 0;padding:22px 24px;align-items:stretch}.detail-gallery-col{display:flex;flex-direction:column;gap:8px}.detail-main-photo{position:relative;height:164px;border-radius:12px;overflow:hidden;border:1px solid var(--bd-e2e8f0);background:#0f172a}.detail-main-photo img{width:100%;height:100%;object-fit:cover;display:block}.detail-photo-badge{position:absolute;right:8px;bottom:8px;padding:3px 9px;border-radius:7px;background:rgba(15,23,42,.8);color:#fff;font-size:var(--font-size-xs);font-weight:600}.detail-thumb-row{display:grid;grid-template-columns:repeat(3,1fr);gap:8px}.detail-thumb-btn{height:48px;border-radius:8px;overflow:hidden;border:2px solid transparent;padding:0;background:var(--bg-f1f5f9);cursor:pointer;opacity:.65;transition:opacity .15s,border-color .15s}.detail-thumb-btn:hover{opacity:.9}.detail-thumb-btn.active{border-color:#2563eb;opacity:1}.detail-thumb-btn img{width:100%;height:100%;object-fit:cover;display:block}
.detail-summary-col{display:flex;flex-direction:column;justify-content:space-between;gap:18px;min-width:0}.detail-summary-top{display:flex;align-items:flex-start;justify-content:space-between;gap:16px;flex-wrap:wrap}.detail-eyebrow-row{display:flex;align-items:center;gap:10px;flex-wrap:wrap;margin-bottom:6px}.detail-summary-top h1{margin:0 0 4px;color:var(--tx-0f172a);font-size:var(--font-size-2xl);font-weight:700}.detail-summary-top p,.card-heading p,.panel-heading p{margin:0;color:var(--tx-64748b);font-size:var(--font-size-sm)}.detail-kpi-strip{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:12px}.detail-kpi-box{padding:12px 14px;border-radius:12px;border:1px solid var(--bd-e2e8f0);background:var(--bg-f8fafc);display:grid;gap:3px}.detail-kpi-box span{color:var(--tx-64748b);font-size:var(--font-size-xs);font-weight:500}.detail-kpi-box strong{color:var(--tx-0f172a);font-size:var(--font-size-xl);font-weight:700;line-height:1.2}.detail-kpi-box strong small{font-size:var(--font-size-xs);font-weight:500;color:var(--tx-64748b)}.detail-kpi-box.is-green strong{color:var(--tx-10b981)}.detail-kpi-box em{font-style:normal;color:var(--tx-64748b);font-size:var(--font-size-xs);font-weight:500}
.detail-header{display:flex;align-items:center;justify-content:space-between;gap:24px;margin:18px 0;padding:24px 26px}.eyebrow{color:var(--tx-2563eb);font-size:var(--font-size-xs);font-weight:700}.detail-header h1{margin:4px 0;color:var(--tx-0f172a);font-size:var(--font-size-2xl)}.detail-header p{margin:0;color:var(--tx-64748b);font-size:var(--font-size-sm)}.header-actions{display:flex;align-items:center;gap:10px;flex-wrap:wrap}.status-badge{padding:5px 12px;border-radius:7px;font-size:var(--font-size-xs);font-weight:600}.status-available{background:var(--bg-dcfce7);color:var(--tx-15803d)}.status-reserved{background:var(--bg-fef3c7);color:var(--tx-d97706)}.status-occupied{background:var(--bg-f1f5f9);color:var(--tx-64748b)}.detail-card{padding:24px 28px}.detail-tabs-bar{display:flex;align-items:center;gap:24px;border-bottom:1px solid var(--bd-f1f5f9);margin-bottom:22px;overflow-x:auto}.detail-tab-btn{position:relative;display:inline-flex;align-items:center;gap:8px;padding:12px 4px 14px;border:0;background:transparent;color:var(--tx-64748b);font-family:inherit;font-size:var(--font-size-base);font-weight:600;cursor:pointer;white-space:nowrap;transition:color .18s ease}.detail-tab-btn:hover{color:var(--tx-1e293b)}.detail-tab-btn.active{color:var(--tx-2563eb)}.detail-tab-btn.active::after{content:'';position:absolute;left:0;right:0;bottom:-1px;height:3px;border-radius:2px;background:#2563eb}.tab-count-pill{display:inline-flex;align-items:center;justify-content:center;min-width:22px;height:20px;padding:0 7px;border-radius:999px;background:var(--bg-f1f5f9);color:var(--tx-475569);font-size:var(--font-size-xs);font-weight:600}.detail-tab-btn.active .tab-count-pill{background:var(--bg-eff6ff);color:var(--tx-2563eb)}.panel-heading,.card-heading{display:flex;align-items:flex-start;justify-content:space-between;gap:12px;margin-bottom:18px}.panel-heading h2,.card-heading h2{margin:0 0 3px;color:var(--tx-1e293b);font-size:var(--font-size-lg);font-weight:600}.fact-grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:12px;margin:0}.fact-grid>div{padding:14px 16px;border:1px solid var(--bd-e2e8f0);border-radius:10px;background:var(--bg-ffffff)}.fact-grid .full{grid-column:1/-1}.fact-grid dt{color:var(--tx-64748b);font-size:var(--font-size-xs);font-weight:500}.fact-grid dd{margin:5px 0 0;color:var(--tx-1e293b);font-size:var(--font-size-base);font-weight:600;word-break:break-word}.table-wrap{overflow-x:auto}.contracts-table{width:100%;border-collapse:collapse}.contracts-table thead th{padding:14px 16px;border-bottom:1px solid var(--bd-f1f5f9);color:var(--tx-64748b);font-size:var(--font-size-xs);font-weight:500;text-align:left;background:var(--bg-ffffff)}.contracts-table tbody td{padding:16px;border-bottom:1px solid var(--bd-f1f5f9);color:var(--tx-334155);font-size:var(--font-size-sm);font-weight:500}.contracts-table tbody tr:hover{background:var(--bg-f8fafc)}.contracts-table strong{color:var(--tx-1e293b);font-weight:700}.empty-state{margin:0;padding:28px 20px;text-align:center;border-radius:10px;background:var(--bg-f8fafc);color:var(--tx-64748b);font-size:var(--font-size-sm)}.not-found{display:grid;justify-items:center;gap:10px;margin-top:18px;padding:40px;text-align:center}.not-found h1{margin:0;font-size:var(--font-size-xl);color:var(--tx-1e293b)}.not-found p{margin:0;color:var(--tx-64748b)}.skeleton-card{min-height:120px;background:linear-gradient(100deg,var(--bg-f1f5f9) 30%,var(--bg-f8fafc) 45%,var(--bg-f1f5f9) 60%);background-size:220% 100%;animation:shimmer 1.4s infinite}.skeleton-card.tall{min-height:360px}@keyframes shimmer{to{background-position:-220% 0}}
.timeline-edge-trigger{position:fixed;top:52%;right:0;z-index:900;width:44px;height:150px;transform:translateY(-50%);display:block;padding:0;overflow:hidden;border:1px solid var(--bd-dbeafe);border-right:0;border-radius:12px 0 0 12px;background:var(--bg-rgba-255-255-255-0p98);color:var(--tx-2563eb);box-shadow:-6px 7px 20px rgba(15,23,42,.11);backdrop-filter:blur(10px);font-family:inherit;font-size:var(--font-size-xs);font-weight:700;white-space:nowrap;cursor:pointer;transition:background .18s ease,width .18s ease,box-shadow .18s ease}.timeline-edge-trigger>svg{position:absolute;top:15px;left:50%;width:17px;height:17px;transform:translateX(-50%)}.timeline-edge-trigger>span{position:absolute;top:56%;left:50%;transform:translate(-50%,-50%) rotate(-90deg);transform-origin:center}.timeline-edge-trigger:hover,.timeline-edge-trigger.active{width:48px;background:var(--bg-eff6ff);color:var(--tx-1d4ed8);box-shadow:-8px 9px 24px rgba(37,99,235,.16)}.timeline-backdrop{position:fixed;inset:0;z-index:1000;display:flex;justify-content:flex-end;background:rgba(15,23,42,.45);backdrop-filter:blur(2px)}.timeline-panel{width:min(480px,92vw);height:100%;background:var(--bg-ffffff);box-shadow:-16px 0 40px rgba(15,23,42,.16);overflow-y:auto}.timeline-header{position:sticky;top:0;z-index:2;display:flex;align-items:center;justify-content:space-between;padding:22px 24px;border-bottom:1px solid var(--bd-e2e8f0);background:var(--bg-ffffff)}.timeline-header h2{margin:0;font-size:var(--font-size-xl);color:var(--tx-1e293b)}.close-button{width:36px;height:36px;display:grid;place-items:center;border:1px solid var(--bd-e2e8f0);border-radius:9px;background:var(--bg-ffffff);color:var(--tx-475569);cursor:pointer}.timeline-body{padding:24px}.timeline-list{display:flex;flex-direction:column}.timeline-item{position:relative;display:flex;gap:14px;padding-bottom:24px}.timeline-item:not(:last-child)::before{content:'';position:absolute;left:5px;top:13px;bottom:0;width:1px;background:var(--bg-e2e8f0)}.timeline-dot{position:relative;z-index:1;width:11px;height:11px;margin-top:5px;flex:0 0 11px;border-radius:50%;background:#3b82f6;box-shadow:0 0 0 4px #dbeafe}.timeline-copy{flex:1;min-width:0}.timeline-copy strong{display:block;font-size:var(--font-size-sm);color:var(--tx-1e293b)}.timeline-copy time{display:block;color:var(--tx-64748b);font-size:var(--font-size-xs);margin-top:3px}.timeline-copy p{margin:8px 0 0;padding:8px 10px;border-radius:8px;background:var(--bg-f8fafc);color:var(--tx-475569);font-size:var(--font-size-xs);line-height:1.45}.timeline-empty{padding:30px 10px;text-align:center;color:var(--tx-94a3b8);font-size:var(--font-size-sm)}.timeline-slide-enter-active,.timeline-slide-leave-active{transition:background-color .28s ease,backdrop-filter .28s ease}.timeline-slide-enter-active .timeline-panel,.timeline-slide-leave-active .timeline-panel{transition:transform .32s cubic-bezier(.22,1,.36,1),box-shadow .32s ease}.timeline-slide-enter-from,.timeline-slide-leave-to{background:var(--bg-rgba-15-23-42-0);backdrop-filter:blur(0)}.timeline-slide-enter-from .timeline-panel,.timeline-slide-leave-to .timeline-panel{transform:translateX(100%);box-shadow:none}
@media(max-width:960px){.detail-hero-card{grid-template-columns:1fr}.detail-kpi-strip{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:800px){.detail-header{align-items:flex-start;flex-direction:column}.fact-grid{grid-template-columns:1fr 1fr}}@media(max-width:520px){.fact-grid,.detail-kpi-strip{grid-template-columns:1fr}.fact-grid .full{grid-column:auto}.header-actions{width:100%;justify-content:space-between}}
</style>
