<template>
  <section class="msp" role="dialog" aria-label="ตารางสรุปการจดมิเตอร์">
    <header class="msp-head">
      <div>
        <h3>ตารางสรุปการจด</h3>
        <p>ความคืบหน้าการจดมิเตอร์น้ำ-ไฟรายโครงการ</p>
      </div>
      <button type="button" class="msp-close" aria-label="ปิดตาราง" @click="$emit('close')">
        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M6 6l12 12M18 6L6 18" /></svg>
      </button>
    </header>
    <div class="msp-body">
      <table v-if="projects.length" class="msp-table">
        <thead><tr><th>โครงการ</th><th>จดแล้ว</th><th>สถานะ</th><th></th></tr></thead>
        <tbody>
          <tr v-for="row in projects" :key="row.id">
            <td><strong>{{ row.name }}</strong><small>{{ row.method }}</small></td>
            <td>
              <div class="msp-progress"><i :style="{ width: `${percent(row.recorded, row.total)}%` }"></i></div>
              <small>{{ row.recorded }}/{{ row.total }} แผง</small>
            </td>
            <td><span class="msp-pill" :class="row.ready ? 'is-ready' : 'is-pending'">{{ row.ready ? 'พร้อมออกบิล' : 'ยังไม่ครบ' }}</span></td>
            <td><AppButton variant="outline" size="sm" @click="$emit('open-project', row.id)">เปิดผัง</AppButton></td>
          </tr>
        </tbody>
      </table>
      <AppEmptyState v-else variant="compact" icon="document" title="ยังไม่มีแผงที่ต้องจดมิเตอร์" description="แผงที่มีผู้เช่าจะแสดงในตารางนี้" />
    </div>
  </section>
</template>

<script setup lang="ts">
import AppButton from '~/components/AppButton.vue'
import AppEmptyState from '~/components/AppEmptyState.vue'
import type { MeterProjectRow } from './StallMeterSummaryModal.vue'

defineProps<{ projects: MeterProjectRow[] }>()
defineEmits<{ close: []; 'open-project': [id: string] }>()
const percent = (recorded: number, total: number) => total ? Math.round((recorded / total) * 100) : 0
</script>

<style scoped>
.msp { display: flex; flex-direction: column; min-height: 0; height: 100%; color: var(--tx-334155); }
.msp-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; padding: 16px; border-bottom: 1px solid var(--bd-e2e8f0); }
.msp-head h3 { margin: 0; color: var(--tx-1e293b); font-size: 1.0625rem; font-weight: 700; }
.msp-head p { margin: 3px 0 0; color: var(--tx-64748b); font-size: .8125rem; font-weight: 500; }
.msp-close { display: grid; place-items: center; width: 32px; height: 32px; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-64748b); cursor: pointer; }
.msp-body { min-height: 0; flex: 1 1 auto; overflow: auto; padding: 0 16px 16px; }
.msp-table { width: 100%; border-collapse: collapse; }
.msp-table th { padding: 12px 8px; border-bottom: 1px solid var(--bd-f1f5f9); color: var(--tx-64748b); font-size: .75rem; font-weight: 600; text-align: left; white-space: nowrap; }
.msp-table td { padding: 12px 8px; border-bottom: 1px solid var(--bd-f1f5f9); vertical-align: middle; font-size: .8125rem; }
.msp-table strong, .msp-table small { display: block; }
.msp-table strong { color: var(--tx-1e293b); font-weight: 700; }
.msp-table small { margin-top: 2px; color: var(--tx-64748b); font-size: .75rem; }
.msp-progress { width: 72px; height: 7px; overflow: hidden; border-radius: 999px; background: var(--bg-e2e8f0); }
.msp-progress i { display: block; height: 100%; border-radius: inherit; background: #2563eb; }
.msp-pill { display: inline-flex; padding: 4px 8px; border-radius: 7px; font-size: .75rem; font-weight: 600; white-space: nowrap; }
.msp-pill.is-ready { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.msp-pill.is-pending { background: var(--bg-fef3c7); color: var(--tx-d97706); }
</style>
