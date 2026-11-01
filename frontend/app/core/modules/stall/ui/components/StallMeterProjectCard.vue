<template>
  <section class="mpc" aria-label="สรุปการจดมิเตอร์ของโครงการ">
    <template v-if="rows.length">
      <div class="mpc-progress-card">
        <div class="mpc-progress-head">
          <div>
            <strong class="mpc-percent">{{ percent }}%</strong>
            <span class="mpc-percent-label">จดแล้ว</span>
          </div>
          <span class="mpc-count">{{ recorded.toLocaleString('th-TH') }} จาก {{ rows.length.toLocaleString('th-TH') }} แผง</span>
        </div>
        <div class="mpc-track" role="progressbar" :aria-valuenow="percent" aria-valuemin="0" aria-valuemax="100"><i :style="{ width: `${percent}%` }"></i></div>
      </div>

      <div class="mpc-tiles">
        <button v-for="tile in tiles" :key="tile.key" type="button" class="mpc-tile" :class="[`is-${tile.key}`, { 'is-active': filter === tile.key }]" :aria-pressed="filter === tile.key" @click="toggleFilter(tile.key)">
          <span>{{ tile.label }}</span>
          <strong>{{ tile.count }}</strong>
        </button>
      </div>

      <section class="mpc-billing" aria-label="สรุปความพร้อมออกบิล">
        <div class="mpc-billing-item">
          <span>พร้อมออกบิล</span>
          <strong>{{ billableRows.length.toLocaleString('th-TH') }} <small>แผง</small></strong>
        </div>
        <div class="mpc-billing-item is-amount">
          <span>ยอดประมาณ</span>
          <strong>฿{{ fmt(billableAmount) }}</strong>
        </div>
        <p>คำนวณเฉพาะแผงที่จดแล้วและไม่พบค่าผิดปกติ</p>
      </section>

      <div class="mpc-list-head">
        <h4>งานที่ต้องดำเนินการ</h4>
        <div class="mpc-list-tools">
          <span>{{ filteredRows.length.toLocaleString('th-TH') }} แผง</span>
          <button v-if="filter" type="button" class="mpc-clear" @click="filter = ''">ล้างตัวกรอง</button>
        </div>
      </div>
      <ul v-if="visibleRows.length" class="mpc-list">
        <li v-for="row in visibleRows" :key="row.code">
          <button type="button" class="mpc-row" @click="emit('open-stall', row.code)">
            <span class="mpc-row-main">
              <strong>{{ row.label }}</strong>
              <small>{{ row.tenant || '-' }}</small>
            </span>
            <span class="mpc-pill" :class="`is-${row.state}`">{{ stateLabel[row.state] }}</span>
          </button>
        </li>
      </ul>
      <AppEmptyState v-else-if="filter && !visibleRows.length" variant="compact" icon="search" title="ไม่พบแผงตามตัวกรอง" description="กดช่องสถานะเดิมอีกครั้งเพื่อล้างตัวกรอง" />
      <AppEmptyState v-else-if="!visibleRows.length" variant="compact" icon="document" title="จดครบทุกแผงแล้ว" description="แผงที่จดเสร็จแล้วจะไม่แสดงซ้ำในรายการงาน" />
      <div v-if="visibleRows.length" class="mpc-pagination" aria-label="แบ่งหน้ารายการแผงที่ต้องดำเนินการ">
        <span>แสดง {{ pageStart }}–{{ pageEnd }} จาก {{ filteredRows.length.toLocaleString('th-TH') }} แผง</span>
        <div class="mpc-pagination-actions">
          <button type="button" :disabled="page === 1" @click="page -= 1">ก่อนหน้า</button>
          <span>หน้า {{ page }} / {{ pageCount }}</span>
          <button type="button" :disabled="page === pageCount" @click="page += 1">ถัดไป</button>
        </div>
      </div>
    </template>
    <AppEmptyState v-else variant="compact" icon="document" title="โครงการนี้ไม่มีแผงที่ต้องจดมิเตอร์" description="แผงที่มีผู้เช่าตามสัญญาเท่านั้นที่ต้องจดมิเตอร์น้ำ-ไฟ" />

    <div v-if="rows.length" class="mpc-footer">
      <AppButton variant="primary" size="md" full-width :disabled="!pending" @click="emit('next')">
        {{ pending ? `จดแผงถัดไป (เหลือ ${pending} แผง)` : 'จดครบทุกแผงแล้ว' }}
      </AppButton>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import AppButton from '~/components/AppButton.vue'
import AppEmptyState from '~/components/AppEmptyState.vue'
import type { MeterStallRow } from './StallMeterSummaryModal.vue'

// การ์ดสรุปการจดมิเตอร์ของโครงการ (แสดงแทนข้อมูลแผงในการ์ดโครงการ ฝั่งซ้าย เมื่ออยู่มุมมองจดมิเตอร์)
const props = defineProps<{ rows: MeterStallRow[] }>()
const emit = defineEmits<{ 'open-stall': [code: string]; next: [] }>()

type Filter = '' | 'pending' | 'abnormal' | 'skipped'
const filter = ref<Filter>('')
const stateLabel: Record<MeterStallRow['state'], string> = { recorded: 'จดแล้ว', abnormal: 'ผิดปกติ', pending: 'ยังไม่จด', skipped: 'ข้าม' }

const count = (state: MeterStallRow['state']) => props.rows.filter(row => row.state === state).length
const recorded = computed(() => props.rows.filter(row => row.recorded).length)
const pending = computed(() => count('pending'))
const percent = computed(() => (props.rows.length ? Math.round((recorded.value / props.rows.length) * 100) : 0))
const billableRows = computed(() => props.rows.filter(row => row.recorded && !row.abnormal))
const billableAmount = computed(() => billableRows.value.reduce((sum, row) => sum + (row.amount || 0), 0))

const tiles = computed(() => [
  { key: 'pending' as const, label: 'ยังไม่จด', count: pending.value },
  { key: 'abnormal' as const, label: 'ผิดปกติ', count: count('abnormal') },
  { key: 'skipped' as const, label: 'ข้าม', count: count('skipped') }
])
const toggleFilter = (key: Filter) => { filter.value = filter.value === key ? '' : key }

// การ์ดซ้ายแสดงเฉพาะงานที่ยังต้องดำเนินการและจำกัดจำนวน เพื่อไม่ให้ยาวเมื่อโครงการมีหลายร้อยแผง
const ORDER: Record<MeterStallRow['state'], number> = { pending: 0, abnormal: 1, skipped: 2, recorded: 3 }
const LIST_LIMIT = 6
const page = ref(1)
const actionRows = computed(() => props.rows.filter(row => row.state !== 'recorded'))
const filteredRows = computed(() => {
  const list = filter.value ? actionRows.value.filter(row => row.state === filter.value) : actionRows.value
  return [...list].sort((a, b) => ORDER[a.state] - ORDER[b.state] || a.label.localeCompare(b.label, 'th', { numeric: true }))
})
const pageCount = computed(() => Math.max(1, Math.ceil(filteredRows.value.length / LIST_LIMIT)))
const visibleRows = computed(() => filteredRows.value.slice((page.value - 1) * LIST_LIMIT, page.value * LIST_LIMIT))
const pageStart = computed(() => (filteredRows.value.length ? (page.value - 1) * LIST_LIMIT + 1 : 0))
const pageEnd = computed(() => Math.min(page.value * LIST_LIMIT, filteredRows.value.length))
watch(filter, () => { page.value = 1 })
watch(pageCount, (total) => { if (page.value > total) page.value = total })
const fmt = (value: number) => Number(value).toLocaleString('th-TH', { maximumFractionDigits: 2 })
</script>

<style scoped>
.mpc { display: flex; flex-direction: column; gap: 14px; min-width: 0; }

.mpc-progress-card { padding: 14px 16px; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; background: var(--bg-f8fafc); }
.mpc-progress-head { display: flex; align-items: flex-end; justify-content: space-between; gap: 10px; margin-bottom: 10px; }
.mpc-percent { color: var(--tx-0f172a); font-size: 2rem; font-weight: 800; line-height: 1; }
.mpc-percent-label { margin-left: 6px; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 600; }
.mpc-count { color: var(--tx-475569); font-size: 0.875rem; font-weight: 600; }
.mpc-track { height: 8px; overflow: hidden; border-radius: 999px; background: var(--bg-e2e8f0); }
.mpc-track i { display: block; height: 100%; border-radius: 999px; background: #2563eb; transition: width .3s ease; }

.mpc-tiles { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 8px; }
.mpc-tile { display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 10px 4px; border: 1px solid var(--bd-e2e8f0); border-top: 3px solid var(--bd-cbd5e1); border-radius: 10px; background: var(--bg-ffffff); color: var(--tx-64748b); font: inherit; font-size: 0.8125rem; font-weight: 600; cursor: pointer; transition: background-color .15s ease, border-color .15s ease; }
.mpc-tile strong { color: var(--tx-0f172a); font-size: 1.375rem; font-weight: 800; line-height: 1.2; }
.mpc-tile.is-pending { border-top-color: #f59e0b; }
.mpc-tile.is-abnormal { border-top-color: #ef4444; }
.mpc-tile.is-abnormal strong { color: var(--tx-dc2626); }
.mpc-tile:hover { background: var(--bg-f8fafc); }
.mpc-tile.is-active { background: var(--bg-eff6ff); border-color: #2563eb; }

.mpc-billing { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 8px; padding: 12px; border: 1px solid var(--bd-dbeafe); border-radius: 12px; background: var(--bg-eff6ff); }
.mpc-billing-item { min-width: 0; display: flex; flex-direction: column; gap: 3px; padding: 9px 10px; border-radius: 9px; background: var(--bg-ffffff); }
.mpc-billing-item > span { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.mpc-billing-item strong { color: var(--tx-0f172a); font-size: 1.25rem; font-weight: 800; line-height: 1.2; }
.mpc-billing-item strong small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.mpc-billing-item.is-amount strong { color: var(--primary); }
.mpc-billing p { grid-column: 1 / -1; margin: 1px 2px 0; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; line-height: 1.45; }

.mpc-list-head { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.mpc-list-head h4 { margin: 0; color: var(--tx-0f172a); font-size: 0.9375rem; font-weight: 700; }
.mpc-list-tools { display: flex; align-items: center; gap: 5px; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.mpc-clear { padding: 2px 8px; border: 0; border-radius: 6px; background: transparent; color: var(--tx-2563eb); font: inherit; font-size: 0.8125rem; font-weight: 600; cursor: pointer; }
.mpc-clear:hover { background: var(--bg-eff6ff); }
.mpc-list { list-style: none; margin: 0; padding: 0; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; overflow: hidden; }
.mpc-list > li + li { border-top: 1px solid var(--bd-f1f5f9); }
.mpc-row { width: 100%; display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 10px 14px; border: 0; background: var(--bg-ffffff); font: inherit; text-align: left; cursor: pointer; transition: background-color .15s ease; }
.mpc-row:hover { background: var(--bg-f8fafc); }
.mpc-row:focus-visible { outline: 2px solid #93c5fd; outline-offset: -2px; }
.mpc-row-main { min-width: 0; display: flex; flex-direction: column; }
.mpc-row-main strong { color: var(--tx-0f172a); font-size: 0.9375rem; font-weight: 700; }
.mpc-row-main small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.mpc-pill { flex: 0 0 auto; padding: 3px 10px; border-radius: 8px; font-size: 0.8125rem; font-weight: 600; white-space: nowrap; }
.mpc-pill.is-pending { background: var(--bg-fef3c7); color: var(--tx-b45309); }
.mpc-pill.is-recorded { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.mpc-pill.is-abnormal { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.mpc-pill.is-skipped { background: var(--bg-f1f5f9); color: var(--tx-64748b); }
.mpc-pagination { display: flex; align-items: center; justify-content: space-between; gap: 8px; color: var(--tx-64748b); font-size: 0.75rem; font-weight: 500; }
.mpc-pagination-actions { display: flex; align-items: center; gap: 4px; }
.mpc-pagination-actions button { padding: 4px 7px; border: 1px solid var(--bd-e2e8f0); border-radius: 6px; background: var(--bg-ffffff); color: var(--tx-475569); font: inherit; font-size: 0.75rem; cursor: pointer; }
.mpc-pagination-actions button:hover:not(:disabled) { border-color: var(--bd-93c5fd); color: var(--tx-2563eb); background: var(--bg-eff6ff); }
.mpc-pagination-actions button:disabled { cursor: not-allowed; opacity: .45; }
.mpc-pagination-actions > span { min-width: 58px; text-align: center; }

.mpc-footer { position: sticky; bottom: -1px; margin: 0 -2px; padding: 12px 2px 2px; background: linear-gradient(to top, var(--bg-ffffff) 78%, transparent); }
</style>
