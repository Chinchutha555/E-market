<template>
  <section class="mpc" aria-label="สรุปการจดมิเตอร์ของโครงการ">
    <template v-if="rows.length">
      <div class="mpc-progress-card">
        <StallMeterPeriodPicker v-if="cycle" class="mpc-period-slot" :month="month" :current-month="currentMonth" :due="cycle.due" @update:month="emit('update:month', $event)" />
        <button v-if="readOnly" type="button" class="mpc-back-current" @click="pickMonth(currentMonth)">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="1 4 1 10 7 10"></polyline><path d="M3.51 15a9 9 0 1 0 .49-8.5L1 10"></path></svg>
          กลับงวดปัจจุบัน
        </button>
        <div class="mpc-progress-head">
          <strong class="mpc-count">จดแล้ว {{ recorded.toLocaleString('th-TH') }} จาก {{ rows.length.toLocaleString('th-TH') }} แผง</strong>
          <span class="mpc-percent">{{ percent }}%</span>
        </div>
        <div class="mpc-track" role="progressbar" :aria-valuenow="percent" aria-valuemin="0" aria-valuemax="100"><i :style="{ width: `${percent}%` }"></i></div>
      </div>

      <div v-if="stats || notes.length" class="mpc-stats" aria-label="สรุปแผงในโครงการ">
        <template v-if="stats">
          <div class="mpc-stat">
            <span>แผงทั้งหมด</span>
            <strong>{{ stats.total.toLocaleString('th-TH') }}</strong>
          </div>
          <div class="mpc-stat is-toread">
            <span>ต้องจดมิเตอร์</span>
            <strong>{{ rows.length.toLocaleString('th-TH') }}</strong>
          </div>
          <div class="mpc-stat">
            <span>ไม่มีผู้เช่า</span>
            <strong>{{ stats.vacant.toLocaleString('th-TH') }}</strong>
          </div>
          <div class="mpc-stat">
            <span>เหมาจ่าย</span>
            <strong>{{ stats.flat.toLocaleString('th-TH') }}</strong>
          </div>
        </template>
        <div v-for="note in notes" :key="note.state" class="mpc-stat" :class="`is-${note.state}`">
          <span>{{ note.label }}</span>
          <strong>{{ note.count.toLocaleString('th-TH') }}</strong>
        </div>
      </div>

      <div class="mpc-hint">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
        <span>{{ readOnly ? 'แตะแผงบนผังเพื่อดูเลขมิเตอร์ที่จดไว้ในงวดนี้' : pending ? 'แตะแผงสีส้มบนผังเพื่อเลือกจดมิเตอร์เอง' : 'แตะแผงบนผังเพื่อดูหรือแก้เลขมิเตอร์ที่จดแล้ว' }}</span>
      </div>

      <div v-if="!readOnly" class="mpc-footer">
        <AppButton v-if="!readOnly" variant="primary" size="lg" full-width :disabled="!pending" @click="emit('next')">
          {{ pending ? `จดแผงถัดไป (เหลือ ${pending} แผง)` : 'จดครบทุกแผงแล้ว' }}
        </AppButton>
      </div>

    </template>
    <AppEmptyState v-else class="meter-empty" variant="compact" icon="document" title="ไม่มีแผงที่ต้องจด" description="เฉพาะแผงที่มีผู้เช่าเท่านั้น" />
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppButton from '~/components/AppButton.vue'
import AppEmptyState from '~/components/AppEmptyState.vue'
import StallMeterPeriodPicker from './StallMeterPeriodPicker.vue'
import type { MeterStallRow } from './StallMeterSummaryModal.vue'

// การ์ดสรุปการจดมิเตอร์ของโครงการ (แสดงแทนข้อมูลแผงในการ์ดโครงการ ฝั่งซ้าย เมื่ออยู่มุมมองจดมิเตอร์)
const props = defineProps<{
  rows: MeterStallRow[]
  stats?: { total: number; vacant: number; flat: number } | null
  /** รอบจดปัจจุบัน เช่น { label: 'พฤศจิกายน 2569', due: '30 พ.ย.' } */
  cycle?: { label: string; due: string } | null
  /** งวด (เดือนบิล) ที่เลือกอยู่ รูปแบบ YYYY-MM */
  /** งวดย้อนหลัง: ดูได้อย่างเดียว */
  readOnly?: boolean
  month: string
  /** งวดปัจจุบัน YYYY-MM (เลือกงวดหลังจากนี้ไม่ได้) */
  currentMonth: string
}>()
const emit = defineEmits<{ next: []; 'update:month': [value: string] }>()

const stateLabel: Record<MeterStallRow['state'], string> = { recorded: 'จดแล้ว', abnormal: 'ผิดปกติ', pending: 'ยังไม่จด', skipped: 'ข้าม', 'no-meter': 'ยังไม่มีมิเตอร์', 'no-data': 'ไม่มีข้อมูลงวดนี้' }

const pickMonth = (value: string | null) => { if (value) emit('update:month', value) }

const count = (state: MeterStallRow['state']) => props.rows.filter(row => row.state === state).length
const recorded = computed(() => props.rows.filter(row => row.recorded).length)
const pending = computed(() => count('pending'))
const percent = computed(() => (props.rows.length ? Math.round((recorded.value / props.rows.length) * 100) : 0))

// แผงที่จดไม่ได้หรือต้องตรวจเพิ่ม: บอกจำนวนเฉพาะสถานะที่มีอยู่ (แผนที่ระบายสีให้หาแผงได้เอง)
const notes = computed(() => (['no-data', 'no-meter', 'abnormal', 'skipped'] as const)
  .map(state => ({ state, label: stateLabel[state], count: count(state) }))
  .filter(note => note.count > 0))
</script>

<style scoped>
.mpc { flex: 1 0 auto; display: flex; flex-direction: column; gap: 14px; min-width: 0; }

.mpc-progress-card { padding: 14px 16px; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; background: var(--bg-f8fafc); }
.mpc-progress-head { display: flex; align-items: baseline; justify-content: space-between; gap: 10px; margin-bottom: 10px; }
.mpc-count { color: var(--tx-0f172a); font-size: var(--font-size-base); font-weight: 700; }
.mpc-percent { color: var(--tx-2563eb); font-size: var(--font-size-base); font-weight: 800; }
.mpc-track { height: 8px; overflow: hidden; border-radius: 999px; background: var(--bg-e2e8f0); }
.mpc-track i { display: block; height: 100%; border-radius: 999px; background: #2563eb; transition: width .3s ease; }

.mpc-stats { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
.mpc-stat { display: flex; flex-direction: column; gap: 2px; padding: 8px 10px; border: 1px solid var(--bd-eef2f6); border-radius: 8px; background: var(--bg-f8fafc); }
.mpc-stat > span { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.mpc-stat strong { color: var(--tx-0f172a); font-size: var(--font-size-xl); font-weight: 800; line-height: 1.25; font-variant-numeric: tabular-nums; }
.mpc-stat.is-toread strong { color: var(--tx-2563eb); }
.mpc-stat.is-no-data strong { color: var(--tx-64748b); }
.mpc-stat.is-no-meter strong { color: var(--tx-d97706); }
.mpc-stat.is-abnormal strong { color: var(--tx-dc2626); }
.mpc-stat.is-skipped strong { color: var(--tx-64748b); }
.mpc-period-slot { margin: 0 0 16px; }
/* ปุ่มติดล่าง: ขยายพื้นหลังขาวลงไปถึงขอบล่างของแผงด้วยค่าลบเท่าระยะ padding ล่างของแผง (เดสก์ท็อป 18px มือถือ 28px) ไม่ให้เนื้อหาโผล่ใต้ปุ่ม */
.mpc-footer { --mpc-pad-b: 18px; margin: auto 0 calc(-1 * var(--mpc-pad-b)); position: sticky; bottom: calc(-1 * var(--mpc-pad-b)); padding: 12px 0 var(--mpc-pad-b); background: var(--bg-ffffff); }
@media (max-width: 768px) { .mpc-footer { --mpc-pad-b: 28px; } }
.mpc-hint { display: flex; align-items: flex-start; gap: 8px; padding: 4px 2px 0; color: var(--tx-64748b); font-size: 0.875rem; line-height: 1.45; }
.mpc-hint svg { flex-shrink: 0; margin-top: 2px; color: var(--tx-2563eb); }

.meter-empty.app-empty-state { min-height: 0; margin-top: 4px; padding: 22px 16px; border: 1px dashed var(--bd-d5e0ed); border-radius: 16px; background: var(--bg-ffffff); }
</style>
