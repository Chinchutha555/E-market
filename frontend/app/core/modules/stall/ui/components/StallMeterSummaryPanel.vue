<template>
  <section class="msp" role="dialog" aria-label="สรุปการจดมิเตอร์">
    <header class="msp-head">
      <div>
        <h3>สรุปการจดมิเตอร์</h3>
        <p>ภาพรวมทุกโครงการ</p>
      </div>
      <button type="button" class="msp-close" aria-label="ปิดสรุป" @click="$emit('close')">
        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M6 6l12 12M18 6L6 18" /></svg>
      </button>
    </header>

    <div class="msp-period-row">
      <StallMeterPeriodPicker :month="month" :current-month="currentMonth" :due="dueText" @update:month="$emit('update:month', $event)" />
    </div>

    <div v-if="!detailProject" class="msp-body">
      <template v-if="projects.length">
        <div class="msp-overall">
          <div class="msp-overall-top">
            <div>
              <strong class="msp-percent">{{ overallPercent }}%</strong>
              <span class="msp-overall-text">จดแล้ว {{ totals.recorded.toLocaleString('th-TH') }} จาก {{ totals.total.toLocaleString('th-TH') }} แผง</span>
            </div>
            <span class="msp-chip" :class="isPast ? 'is-muted' : 'is-info'">{{ isPast ? 'งวดย้อนหลัง' : daysLeftText }}</span>
          </div>
          <div class="msp-bar" role="progressbar" :aria-valuenow="overallPercent" aria-valuemin="0" aria-valuemax="100"><i :style="{ width: `${overallPercent}%` }"></i></div>
          <p class="msp-remaining">{{ totals.total - totals.recorded > 0 ? `เหลือ ${(totals.total - totals.recorded).toLocaleString('th-TH')} แผง` : 'จดครบทุกแผงแล้ว' }}</p>
        </div>

        <div class="msp-stats" aria-label="ความพร้อมออกบิลและการใช้">
          <div class="msp-stat is-ready"><span>พร้อมออกบิล</span><strong>{{ readyStalls.toLocaleString('th-TH') }} <small>แผง</small></strong></div>
          <div class="msp-stat"><span>ใช้น้ำรวม</span><strong>{{ totals.water.toLocaleString('th-TH') }} <small>หน่วย</small></strong></div>
          <div class="msp-stat"><span>ใช้ไฟรวม</span><strong>{{ totals.electric.toLocaleString('th-TH') }} <small>หน่วย</small></strong></div>
        </div>

        <section v-if="hasFollowUp" class="msp-follow" aria-label="เรื่องที่ต้องตามต่อ">
          <div class="msp-section-head"><h4>ต้องตามต่อ</h4></div>
          <ul class="msp-follow-list">
            <li v-if="noMeterCount" class="msp-follow-item">
              <div><strong>ยังไม่มีมิเตอร์</strong><small>{{ noMeterCount.toLocaleString('th-TH') }} แผง จดไม่ได้จนกว่าจะเพิ่มมิเตอร์</small></div>
              <AppButton variant="outline" size="sm" @click="$emit('open-meter-settings')">ตั้งมิเตอร์</AppButton>
            </li>
            <li v-if="abnormalRows.length" class="msp-follow-item is-column">
              <div><strong>ผิดปกติ</strong><small>{{ abnormalAll.toLocaleString('th-TH') }} แผง ใช้สูงกว่าเกณฑ์ของโครงการ</small></div>
              <ul class="msp-abnormal-list">
                <li v-for="item in abnormalRows" :key="item.code">
                  <button type="button" @click="$emit('open-stall', item.code)">
                    <span class="msp-abn-name">{{ item.label }}<small>{{ item.project }}</small></span>
                    <span class="msp-abn-use">{{ item.text }}</span>
                  </button>
                </li>
              </ul>
            </li>
            <li v-if="skipEntries.length" class="msp-follow-item is-column">
              <div><strong>ข้ามไป</strong><small>{{ skippedAll.toLocaleString('th-TH') }} แผง แยกตามเหตุผล</small></div>
              <div class="msp-reasons">
                <span v-for="entry in skipEntries" :key="entry.reason" class="msp-pill is-skipped">{{ entry.reason }} {{ entry.count }}</span>
              </div>
            </li>
          </ul>
        </section>

        <div class="msp-section-head">
          <h4>รายโครงการ</h4>
          <span>{{ readyCount }} / {{ projects.length }} พร้อมออกบิล</span>
        </div>
        <ul class="msp-list">
          <li v-for="row in sortedProjects" :key="row.id" class="msp-row is-clickable" role="button" tabindex="0" :aria-label="`ดูตารางที่จดของ${row.name}`" @click="openDetail(row.id)" @keydown.enter.prevent="openDetail(row.id)">
            <div class="msp-row-main">
              <div class="msp-row-title">
                <strong>{{ row.name }}</strong>
                <small>{{ row.method }}</small>
              </div>
              <AppButton variant="outline" size="sm" @click.stop="$emit('open-project', row.id)">เปิดผัง</AppButton>
            </div>
            <div class="msp-row-progress">
              <div class="msp-bar is-sm" role="progressbar" :aria-valuenow="percent(row.recorded, row.total)" aria-valuemin="0" aria-valuemax="100"><i :class="{ 'is-done': row.ready }" :style="{ width: `${percent(row.recorded, row.total)}%` }"></i></div>
              <span class="msp-row-go" aria-hidden="true">ดูตาราง ›</span>
              <span class="msp-row-count">{{ row.recorded.toLocaleString('th-TH') }}/{{ row.total.toLocaleString('th-TH') }} แผง</span>
            </div>
            <div class="msp-row-tags">
              <span class="msp-pill" :class="row.ready ? 'is-ready' : 'is-pending'">{{ row.ready ? 'พร้อมออกบิล' : 'ยังไม่ครบ' }}</span>
              <span v-if="row.abnormal" class="msp-pill is-abnormal">ผิดปกติ {{ row.abnormal }}</span>
              <span v-if="row.skipped" class="msp-pill is-skipped">ข้าม {{ row.skipped }}</span>
            </div>
          </li>
        </ul>
      </template>
      <AppEmptyState v-else class="meter-empty" variant="compact" icon="document" title="ยังไม่มีแผงที่ต้องจดมิเตอร์" description="แผงที่มีผู้เช่าจะแสดงที่นี่" />
    </div>

    <div v-else class="msp-body msp-detail">
      <div class="msp-detail-head">
        <button type="button" class="msp-back" @click="detailId = ''">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="15 18 9 12 15 6"></polyline></svg>
          ทุกโครงการ
        </button>
        <AppButton variant="outline" size="sm" @click="$emit('open-project', detailProject.id)">เปิดผัง</AppButton>
      </div>
      <div class="msp-detail-title">
        <h4>{{ detailProject.name }}</h4>
        <span>จดแล้ว {{ detailProject.recorded.toLocaleString('th-TH') }} จาก {{ detailProject.total.toLocaleString('th-TH') }} แผง</span>
      </div>
      <AppSecondaryTabs v-model="detailFilter" :options="detailTabs" aria-label="กรองสถานะการจด" />
      <div v-if="pagedRows.length" class="msp-table-wrap">
        <table class="msp-table">
          <thead>
            <tr>
              <th>แผง</th>
              <th class="num">น้ำ</th>
              <th class="num">ไฟ</th>
              <th class="num">ยอด</th>
              <th>สถานะ</th>
              <th class="act"><span class="sr-only">ดำเนินการ</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in pagedRows" :key="row.code">
              <td><strong>{{ row.label }}</strong><small>{{ row.tenant || '-' }}</small></td>
              <td class="num"><template v-if="row.recorded"><strong>{{ fmt(row.waterUsed) }}</strong><small>{{ fmt(row.waterPrev) }} → {{ fmt(row.waterCurr) }}</small></template><span v-else class="muted">-</span></td>
              <td class="num"><template v-if="row.recorded"><strong :class="{ 'is-danger': row.abnormal }">{{ fmt(row.electricUsed) }}</strong><small>{{ fmt(row.electricPrev) }} → {{ fmt(row.electricCurr) }}</small></template><span v-else class="muted">-</span></td>
              <td class="num"><strong v-if="row.recorded">฿{{ fmt(row.amount) }}</strong><span v-else class="muted">-</span></td>
              <td><span class="msp-pill" :class="`is-${row.state}`">{{ stateLabel[row.state] }}</span></td>
              <td class="act"><AppButton variant="outline" size="sm" @click="$emit('open-stall', row.code)">{{ row.recorded ? 'ดู' : 'จด' }}</AppButton></td>
            </tr>
          </tbody>
        </table>
      </div>
      <AppEmptyState v-else class="meter-empty" variant="compact" icon="search" title="ไม่พบแผงตามตัวกรอง" description="เลือกสถานะอื่น" />
      <AppPagination v-if="filteredRows.length" v-model="detailPage" v-model:items-per-page="detailPerPage" :total-items="filteredRows.length" unit-label="แผง" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import AppButton from '~/components/AppButton.vue'
import AppEmptyState from '~/components/AppEmptyState.vue'
import AppPagination from '~/components/AppPagination.vue'
import AppSecondaryTabs from '~/components/AppSecondaryTabs.vue'
import StallMeterPeriodPicker from './StallMeterPeriodPicker.vue'
import type { MeterProjectRow, MeterStallRow } from './StallMeterSummaryModal.vue'

const props = defineProps<{
  projects: MeterProjectRow[]
  /** งวดที่ดูอยู่ YYYY-MM */
  month: string
  /** งวดปัจจุบัน YYYY-MM */
  currentMonth: string
  /** รายการแผงทุกโครงการ ใช้หาแผงผิดปกติ ไม่มีมิเตอร์ และแผงพร้อมออกบิล */
  stalls: MeterStallRow[]
  /** ค่าเฉลี่ยหน่วยที่ใช้ของแต่ละโครงการ เทียบค่าผิดปกติ */
  averages: Record<string, { electric: number; water: number }>
  /** เหตุผลที่ข้าม -> จำนวนแผง */
  skipReasons: Record<string, number>
}>()
defineEmits<{ close: []; 'open-project': [id: string]; 'open-stall': [code: string]; 'open-meter-settings': []; 'update:month': [value: string] }>()

// ตารางที่จดของโครงการที่เลือก (แตะแถวโครงการเพื่อดู)
const detailId = ref('')
const detailFilter = ref('all')
const detailPage = ref(1)
const detailPerPage = ref(10)
const detailProject = computed(() => props.projects.find(row => row.id === detailId.value) || null)
const detailStalls = computed(() => props.stalls
  .filter(row => row.projectId === detailId.value)
  .sort((a, b) => a.label.localeCompare(b.label, 'th', { numeric: true })))
const stateLabel: Record<MeterStallRow['state'], string> = { recorded: 'จดแล้ว', abnormal: 'ผิดปกติ', pending: 'ยังไม่จด', skipped: 'ข้าม', 'no-meter': 'ยังไม่มีมิเตอร์', 'no-data': 'ไม่มีข้อมูล' }
const countState = (state: MeterStallRow['state']) => detailStalls.value.filter(row => row.state === state).length
const detailTabs = computed(() => [
  { value: 'all', label: 'ทั้งหมด', count: detailStalls.value.length },
  { value: 'recorded', label: 'จดแล้ว', count: countState('recorded') },
  { value: 'pending', label: 'ยังไม่จด', count: countState('pending') },
  { value: 'abnormal', label: 'ผิดปกติ', count: countState('abnormal') },
  { value: 'skipped', label: 'ข้าม', count: countState('skipped') },
  ...(['no-meter', 'no-data'] as const).filter(state => countState(state) > 0).map(state => ({ value: state, label: stateLabel[state], count: countState(state) }))
])
const filteredRows = computed(() => (detailFilter.value === 'all' ? detailStalls.value : detailStalls.value.filter(row => row.state === detailFilter.value)))
const pagedRows = computed(() => filteredRows.value.slice((detailPage.value - 1) * detailPerPage.value, detailPage.value * detailPerPage.value))
const fmt = (value: number | null | undefined) => (value == null ? '-' : Number(value).toLocaleString('th-TH', { maximumFractionDigits: 2 }))
const openDetail = (id: string) => { detailId.value = id; detailFilter.value = 'all'; detailPage.value = 1 }
watch([detailFilter, detailPerPage], () => { detailPage.value = 1 })
const percent = (recorded: number, total: number) => (total ? Math.round((recorded / total) * 100) : 0)

const totals = computed(() => props.projects.reduce((sum, row) => ({
  total: sum.total + row.total,
  recorded: sum.recorded + row.recorded,
  abnormal: sum.abnormal + row.abnormal,
  skipped: sum.skipped + row.skipped,
  water: sum.water + row.waterUnits,
  electric: sum.electric + row.electricUnits
}), { total: 0, recorded: 0, abnormal: 0, skipped: 0, water: 0, electric: 0 }))
const readyStalls = computed(() => props.stalls.filter(row => row.state === 'recorded').length)
const noMeterCount = computed(() => props.stalls.filter(row => row.state === 'no-meter').length)
const abnormalAll = computed(() => props.stalls.filter(row => row.abnormal).length)
const skippedAll = computed(() => props.stalls.filter(row => row.state === 'skipped').length)
const skipEntries = computed(() => Object.entries(props.skipReasons).map(([reason, count]) => ({ reason, count })).sort((a, b) => b.count - a.count))
const projectName = (id: string) => props.projects.find(row => row.id === id)?.name || ''
// ผิดปกติ 5 แผงแรก เรียงจากใช้เกินค่าเฉลี่ยมากสุด พร้อมค่าที่ใช้เทียบค่าเฉลี่ยของโครงการ
const abnormalRows = computed(() => props.stalls
  .filter(row => row.abnormal)
  .map((row) => {
    const avg = props.averages[row.projectId] || { electric: 0, water: 0 }
    const electricRatio = avg.electric > 0 ? (row.electricUsed ?? 0) / avg.electric : 0
    const waterRatio = avg.water > 0 ? (row.waterUsed ?? 0) / avg.water : 0
    const isElectric = electricRatio >= waterRatio
    const used = (isElectric ? row.electricUsed : row.waterUsed) ?? 0
    const mean = isElectric ? avg.electric : avg.water
    return { code: row.code, label: row.label, project: projectName(row.projectId), ratio: Math.max(electricRatio, waterRatio), text: `${isElectric ? 'ไฟ' : 'น้ำ'} ${used.toLocaleString('th-TH')} (เฉลี่ย ${Math.round(mean).toLocaleString('th-TH')})` }
  })
  .sort((a, b) => b.ratio - a.ratio)
  .slice(0, 5))
const hasFollowUp = computed(() => noMeterCount.value > 0 || abnormalRows.value.length > 0 || skipEntries.value.length > 0)
const overallPercent = computed(() => percent(totals.value.recorded, totals.value.total))
const readyCount = computed(() => props.projects.filter(row => row.ready).length)
// โครงการที่จดน้อยสุดขึ้นก่อน เพื่อให้เห็นที่ต้องเร่งก่อน
const sortedProjects = computed(() => [...props.projects].sort((a, b) => percent(a.recorded, a.total) - percent(b.recorded, b.total) || b.total - a.total))

const monthDate = computed(() => { const [y, m] = (props.month || '').split('-').map(Number); return new Date(y || new Date().getFullYear(), (m || 1) - 1, 1) })
const dueDate = computed(() => new Date(monthDate.value.getFullYear(), monthDate.value.getMonth() + 1, 0))
const dueText = computed(() => new Intl.DateTimeFormat('th-TH', { day: 'numeric', month: 'short' }).format(dueDate.value))
const isPast = computed(() => props.month < props.currentMonth)
const daysLeft = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return Math.max(0, Math.round((dueDate.value.getTime() - today.getTime()) / 86400000))
})
const daysLeftText = computed(() => (daysLeft.value === 0 ? 'ครบกำหนดวันนี้' : `เหลือ ${daysLeft.value} วัน`))
</script>

<style scoped>
.msp { display: flex; flex-direction: column; min-height: 0; height: 100%; color: var(--tx-334155); }
.msp-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; padding: 16px 20px; border-bottom: 1px solid var(--bd-e2e8f0); }
.msp-head h3 { margin: 0; color: var(--tx-1e293b); font-size: var(--font-size-lg); font-weight: 700; }
.msp-head p { margin: 3px 0 0; color: var(--tx-64748b); font-size: var(--font-size-sm); font-weight: 500; }
.msp-close { display: grid; place-items: center; flex: 0 0 auto; width: 32px; height: 32px; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-64748b); cursor: pointer; }
.msp-close:hover { background: var(--bg-f1f5f9); color: var(--tx-0f172a); }
.msp-period-row { position: relative; z-index: 5; flex: 0 0 auto; padding: 14px 20px 0; }
.msp-body { min-height: 0; flex: 1 1 auto; overflow: auto; display: flex; flex-direction: column; gap: 16px; padding: 16px 20px 20px; }
.msp-body > * { flex-shrink: 0; }

.msp-overall { padding: 16px; border: 1px solid var(--bd-dbeafe); border-radius: 14px; background: var(--bg-eff6ff); }
.msp-overall-top { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 12px; }
.msp-percent { margin-right: 10px; color: var(--tx-2563eb); font-size: 2rem; font-weight: 800; line-height: 1; font-variant-numeric: tabular-nums; }
.msp-overall-text { color: var(--tx-334155); font-size: var(--font-size-base); font-weight: 600; }
.msp-chip { flex: 0 0 auto; padding: 4px 12px; border-radius: 999px; font-size: var(--font-size-xs); font-weight: 700; white-space: nowrap; }
.msp-chip.is-info { background: var(--bg-ffffff); color: var(--tx-2563eb); border: 1px solid var(--bd-dbeafe); }
.msp-chip.is-muted { background: var(--bg-e2e8f0); color: var(--tx-475569); }
.msp-bar { height: 10px; overflow: hidden; border-radius: 999px; background: var(--bg-ffffff); }
.msp-bar i { display: block; height: 100%; border-radius: inherit; background: #2563eb; transition: width .3s ease; }
.msp-bar.is-sm { flex: 1 1 auto; height: 8px; background: var(--bg-e2e8f0); }
.msp-bar i.is-done { background: #10b981; }
.msp-remaining { margin: 10px 0 0; color: var(--tx-64748b); font-size: var(--font-size-sm); font-weight: 500; }

.msp-stats { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 8px; }
.msp-stat { display: flex; flex-direction: column; gap: 2px; padding: 10px 12px; border: 1px solid var(--bd-eef2f6); border-radius: 10px; background: var(--bg-f8fafc); }
.msp-stat > span { color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 600; }
.msp-stat strong { color: var(--tx-0f172a); font-size: var(--font-size-lg); font-weight: 800; line-height: 1.3; font-variant-numeric: tabular-nums; overflow-wrap: anywhere; }
.msp-stat strong small { color: var(--tx-94a3b8); font-size: var(--font-size-xs); font-weight: 600; }
.msp-stat.is-ready strong { color: var(--tx-15803d); }

.msp-section-head { display: flex; align-items: baseline; justify-content: space-between; gap: 10px; }
.msp-section-head h4 { margin: 0; color: var(--tx-0f172a); font-size: var(--font-size-base); font-weight: 700; }
.msp-section-head span { color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 600; }

.msp-list { list-style: none; margin: -4px 0 0; padding: 0; display: flex; flex-direction: column; gap: 10px; }
.msp-row { display: flex; flex-direction: column; gap: 10px; padding: 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-ffffff); }
.msp-row-main { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.msp-row-title { min-width: 0; display: flex; flex-direction: column; }
.msp-row-title strong { color: var(--tx-0f172a); font-size: var(--font-size-base); font-weight: 700; }
.msp-row-title small { color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.msp-row-progress { display: flex; align-items: center; gap: 12px; }
.msp-row-count { flex: 0 0 auto; color: var(--tx-475569); font-size: var(--font-size-sm); font-weight: 600; font-variant-numeric: tabular-nums; }
.msp-row-tags { display: flex; flex-wrap: wrap; gap: 6px; }

.msp-pill { display: inline-flex; padding: 3px 10px; border-radius: 7px; font-size: var(--font-size-xs); font-weight: 600; white-space: nowrap; }
.msp-pill.is-ready { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.msp-pill.is-pending { background: var(--bg-fef3c7); color: var(--tx-d97706); }
.msp-pill.is-abnormal { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.msp-pill.is-skipped { background: var(--bg-f1f5f9); color: var(--tx-64748b); }
.meter-empty.app-empty-state { min-height: 0; margin-top: 4px; padding: 22px 16px; border: 1px dashed var(--bd-d5e0ed); border-radius: 16px; background: var(--bg-ffffff); }

.msp-follow { display: flex; flex-direction: column; gap: 10px; }
.msp-follow-list { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 8px; }
.msp-follow-item { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 12px 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; background: var(--bg-ffffff); }
.msp-follow-item.is-column { flex-direction: column; align-items: stretch; }
.msp-follow-item strong { display: block; color: var(--tx-0f172a); font-size: var(--font-size-base); font-weight: 700; }
.msp-follow-item small { display: block; margin-top: 2px; color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 500; }
.msp-abnormal-list { list-style: none; margin: 0; padding: 0; border: 1px solid var(--bd-f1f5f9); border-radius: 10px; overflow: hidden; }
.msp-abnormal-list > li + li { border-top: 1px solid var(--bd-f1f5f9); }
.msp-abnormal-list button { display: flex; align-items: center; justify-content: space-between; gap: 10px; width: 100%; padding: 9px 12px; border: 0; background: var(--bg-ffffff); font: inherit; text-align: left; cursor: pointer; }
.msp-abnormal-list button:hover { background: var(--bg-fef2f2); }
.msp-abn-name { min-width: 0; color: var(--tx-0f172a); font-size: var(--font-size-sm); font-weight: 700; }
.msp-abn-name small { display: block; margin-top: 0; color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 500; }
.msp-abn-use { flex: 0 0 auto; color: var(--tx-dc2626); font-size: var(--font-size-sm); font-weight: 700; font-variant-numeric: tabular-nums; }
.msp-reasons { display: flex; flex-wrap: wrap; gap: 6px; }

.msp-row.is-clickable { cursor: pointer; transition: border-color .15s ease, background-color .15s ease; }
.msp-row.is-clickable:hover { border-color: #93c5fd; background: var(--bg-f8fafc); }
.msp-row.is-clickable:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 2px; }
.msp-row-go { flex: 0 0 auto; color: var(--tx-2563eb); font-size: var(--font-size-xs); font-weight: 600; }
.msp-detail { gap: 12px; }
.msp-detail-head { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.msp-back { display: inline-flex; align-items: center; gap: 6px; padding: 6px 10px 6px 8px; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-475569); font: inherit; font-size: var(--font-size-sm); font-weight: 600; cursor: pointer; }
.msp-back:hover { background: var(--bg-f8fafc); color: var(--tx-0f172a); }
.msp-detail-title h4 { margin: 0; color: var(--tx-0f172a); font-size: var(--font-size-lg); font-weight: 700; }
.msp-detail-title span { color: var(--tx-64748b); font-size: var(--font-size-sm); font-weight: 500; }
.msp-table-wrap { overflow-x: auto; -webkit-overflow-scrolling: touch; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; }
.msp-table { width: 100%; min-width: 520px; border-collapse: collapse; }
.msp-table th { padding: 12px 10px; border-bottom: 1px solid var(--bd-f1f5f9); color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 600; text-align: left; white-space: nowrap; }
.msp-table td { padding: 12px 10px; border-bottom: 1px solid var(--bd-f1f5f9); vertical-align: middle; font-size: var(--font-size-sm); }
.msp-table tbody tr:last-child td { border-bottom: 0; }
.msp-table tbody tr:hover { background: var(--bg-f8fafc); }
.msp-table th.num, .msp-table td.num { text-align: right; font-variant-numeric: tabular-nums; }
.msp-table th.act, .msp-table td.act { width: 1%; text-align: right; }
.msp-table strong { display: block; color: var(--tx-1e293b); font-weight: 700; }
.msp-table strong.is-danger { color: var(--tx-dc2626); }
.msp-table small { display: block; margin-top: 2px; color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 500; }
.msp-table .muted { color: var(--tx-94a3b8); }
.msp-pill.is-recorded { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.msp-pill.is-no-meter { background: var(--bg-fef3c7); color: var(--tx-b45309); }
.msp-pill.is-no-data { background: var(--bg-f1f5f9); color: var(--tx-64748b); }
.sr-only { position: absolute; width: 1px; height: 1px; overflow: hidden; clip: rect(0 0 0 0); white-space: nowrap; }

@media (max-width: 640px) {
  .msp-head, .msp-body, .msp-period-row { padding-left: 14px; padding-right: 14px; }
  .msp-stats { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .msp-percent { font-size: 1.75rem; }
}
</style>
