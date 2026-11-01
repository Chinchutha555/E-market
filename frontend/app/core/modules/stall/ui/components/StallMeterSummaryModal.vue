<template>
  <AppModal
    :model-value="modelValue"
    size="xl"
    icon="document"
    :title="`ตารางสรุปการจด · รอบ ${monthLabel}`"
    description="สรุปความคืบหน้าการจดมิเตอร์น้ำ-ไฟ รายโครงการและรายแผง"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <div class="msm-toolbar">
      <AppSecondaryTabs v-model="tab" :options="tabOptions" aria-label="มุมมองตารางสรุปการจด" />
      <div class="msm-toolbar-actions">
        <div v-if="tab === 'stall'" class="msm-project-select">
          <AppSelect v-model="selectedProjectId" :options="projectOptions" />
        </div>
        <AppButton variant="outline" size="md" @click="exportCsv">ส่งออก Excel</AppButton>
      </div>
    </div>

    <!-- รายโครงการ (เรียงจากโครงการที่จดน้อยสุดก่อน) -->
    <div v-if="tab === 'project'" class="msm-table-wrap">
      <table v-if="sortedProjects.length" class="msm-table">
        <thead>
          <tr>
            <th>โครงการ</th>
            <th>วิธีคิด ไฟ / น้ำ</th>
            <th>จดแล้ว</th>
            <th class="num">ผิดปกติ</th>
            <th class="num">ข้าม</th>
            <th class="num">ไฟ (หน่วย)</th>
            <th class="num">น้ำ (ลบ.ม.)</th>
            <th>สถานะ</th>
            <th class="act"><span class="sr-only">การดำเนินการ</span></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in sortedProjects" :key="row.id">
            <td class="strong">{{ row.name }}</td>
            <td class="muted">{{ row.method }}</td>
            <td>
              <div class="msm-progress">
                <span class="msm-progress-track" aria-hidden="true"><i :style="{ width: `${percent(row.recorded, row.total)}%` }"></i></span>
                <span class="msm-progress-text">{{ row.recorded.toLocaleString('th-TH') }}/{{ row.total.toLocaleString('th-TH') }}</span>
              </div>
            </td>
            <td class="num" :class="row.abnormal ? 'is-danger' : 'is-zero'">{{ row.abnormal }}</td>
            <td class="num" :class="{ 'is-zero': !row.skipped }">{{ row.skipped }}</td>
            <td class="num">{{ fmt(row.electricUnits) }}</td>
            <td class="num">{{ fmt(row.waterUnits) }}</td>
            <td><span class="msm-pill" :class="row.ready ? 'is-ready' : 'is-pending'">{{ row.ready ? 'ครบ · พร้อมออกบิล' : 'ยังไม่ครบ' }}</span></td>
            <td class="act"><AppButton variant="outline" size="sm" @click="emit('open-project', row.id)">เปิดผัง</AppButton></td>
          </tr>
        </tbody>
      </table>
      <AppEmptyState v-else variant="compact" icon="document" title="ยังไม่มีแผงที่ต้องจดมิเตอร์" description="แผงที่มีผู้เช่าจะแสดงในตารางนี้เพื่อจดมิเตอร์น้ำ-ไฟ" />
    </div>

    <!-- รายแผง (ของโครงการที่เลือก) -->
    <div v-else class="msm-stall-pane">
      <AppSecondaryTabs v-model="stallFilter" :options="filterOptions" aria-label="กรองสถานะการจด" />
      <div class="msm-table-wrap">
        <table v-if="visibleStalls.length" class="msm-table">
          <thead>
            <tr>
              <th>แผง (ผู้เช่า)</th>
              <th class="num">ไฟ: เดิม → ใหม่</th>
              <th class="num">ใช้ไป</th>
              <th class="num">น้ำ: เดิม → ใหม่</th>
              <th class="num">ใช้ไป</th>
              <th class="num">ยอดประมาณ</th>
              <th>จดโดย</th>
              <th>สถานะ</th>
              <th class="act"><span class="sr-only">การดำเนินการ</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in visibleStalls" :key="row.code">
              <td>
                <strong class="msm-stall-name">{{ row.label }}</strong>
                <small class="msm-stall-tenant">{{ row.tenant || '-' }}</small>
              </td>
              <td class="num">{{ row.recorded ? `${fmt(row.electricPrev)} → ${fmt(row.electricCurr)}` : `${fmt(row.electricPrev)} → -` }}</td>
              <td class="num" :class="{ 'is-danger': row.abnormal }">{{ row.recorded ? fmt(row.electricUsed) : '-' }}</td>
              <td class="num">{{ row.recorded ? `${fmt(row.waterPrev)} → ${fmt(row.waterCurr)}` : `${fmt(row.waterPrev)} → -` }}</td>
              <td class="num" :class="{ 'is-danger': row.abnormal }">{{ row.recorded ? fmt(row.waterUsed) : '-' }}</td>
              <td class="num strong">{{ row.recorded ? `฿${fmt(row.amount)}` : '-' }}</td>
              <td class="muted">{{ row.recordedBy || '-' }}</td>
              <td><span class="msm-pill" :class="`is-${row.state}`">{{ stateLabel[row.state] }}</span></td>
              <td class="act"><AppButton variant="outline" size="sm" @click="emit('open-stall', row.code)">{{ row.recorded ? 'แก้ไข' : 'จด' }}</AppButton></td>
            </tr>
          </tbody>
        </table>
        <AppEmptyState v-else variant="compact" icon="search" title="ไม่พบแผงตามตัวกรอง" description="ลองเลือกสถานะอื่นหรือโครงการอื่น" />
      </div>
    </div>

    <template #footer>
      <p class="msm-footer-note">
        จดครบ {{ readyCount }} โครงการ · {{ readyCount ? 'พร้อมส่งออกใบแจ้งหนี้' : 'ยังไม่มีโครงการที่พร้อมส่งออกใบแจ้งหนี้' }}
        · โครงการที่ยังไม่ครบจะส่งได้เฉพาะแผงที่จดแล้วและไม่ผิดปกติ
      </p>
      <AppButton variant="primary" size="md" :disabled="!projects.length" @click="emit('send-invoice')">ส่งไปออกใบแจ้งหนี้</AppButton>
    </template>
  </AppModal>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import AppButton from '~/components/AppButton.vue'
import AppEmptyState from '~/components/AppEmptyState.vue'
import AppModal from '~/components/AppModal.vue'
import AppSecondaryTabs from '~/components/AppSecondaryTabs.vue'
import AppSelect from '~/components/AppSelect.vue'

// ตารางสรุปการจดมิเตอร์น้ำ-ไฟ (ข้อมูลจำลองจาก useStallBilling) แสดงรายโครงการและรายแผง
export interface MeterProjectRow {
  id: string
  name: string
  method: string
  recorded: number
  total: number
  abnormal: number
  skipped: number
  electricUnits: number
  waterUnits: number
  ready: boolean
}
export interface MeterStallRow {
  code: string
  projectId: string
  label: string
  tenant: string
  recorded: boolean
  abnormal: boolean
  skipped: boolean
  electricPrev: number
  electricCurr: number | null
  electricUsed: number | null
  waterPrev: number
  waterCurr: number | null
  waterUsed: number | null
  amount: number | null
  recordedBy: string
  state: 'recorded' | 'abnormal' | 'pending' | 'skipped'
}

const props = defineProps<{
  modelValue: boolean
  monthLabel: string
  projects: MeterProjectRow[]
  stalls: MeterStallRow[]
  initialProjectId?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'open-project': [id: string]
  'open-stall': [code: string]
  'send-invoice': []
}>()

const tab = ref<'project' | 'stall'>('project')
const stallFilter = ref<'all' | 'recorded' | 'abnormal' | 'pending' | 'skipped'>('all')
const selectedProjectId = ref('')

const stateLabel: Record<MeterStallRow['state'], string> = {
  recorded: 'จดแล้ว',
  abnormal: 'ผิดปกติ',
  pending: 'ยังไม่จด',
  skipped: 'ข้าม'
}

// เรียงจากโครงการที่จดน้อยสุดก่อน
const sortedProjects = computed(() => [...props.projects].sort((a, b) => {
  const ra = a.total ? a.recorded / a.total : 0
  const rb = b.total ? b.recorded / b.total : 0
  return ra - rb || b.total - a.total
}))
const readyCount = computed(() => props.projects.filter(project => project.ready).length)

const projectOptions = computed(() => props.projects.map(project => ({ value: project.id, label: project.name })))
const selectedProject = computed(() => props.projects.find(project => project.id === selectedProjectId.value))

const projectStalls = computed(() => props.stalls.filter(row => row.projectId === selectedProjectId.value))
const countOf = (state: MeterStallRow['state']) => projectStalls.value.filter(row => row.state === state).length
const filterOptions = computed(() => [
  { value: 'all', label: 'ทั้งหมด', count: projectStalls.value.length },
  { value: 'recorded', label: 'จดแล้ว', count: countOf('recorded') },
  { value: 'abnormal', label: 'ผิดปกติ', count: countOf('abnormal') },
  { value: 'pending', label: 'ยังไม่จด', count: countOf('pending') },
  { value: 'skipped', label: 'ข้าม', count: countOf('skipped') }
])
const visibleStalls = computed(() => stallFilter.value === 'all' ? projectStalls.value : projectStalls.value.filter(row => row.state === stallFilter.value))

const tabOptions = computed(() => [
  { value: 'project', label: 'รายโครงการ' },
  { value: 'stall', label: selectedProject.value ? `รายแผง · ${selectedProject.value.name}` : 'รายแผง' }
])

// เปิดหน้าต่างใหม่ทุกครั้ง: เริ่มที่แท็บรายโครงการ และเตรียมโครงการสำหรับแท็บรายแผง (โครงการที่เปิดอยู่บนแผนที่ ถ้าไม่มีใช้โครงการแรก)
watch(() => props.modelValue, (open) => {
  if (!open) return
  tab.value = 'project'
  stallFilter.value = 'all'
  const preferred = props.projects.find(project => project.id === props.initialProjectId)
  selectedProjectId.value = preferred?.id || sortedProjects.value[0]?.id || ''
})
watch(selectedProjectId, () => { stallFilter.value = 'all' })

const percent = (done: number, total: number) => total > 0 ? Math.min(100, Math.round((done / total) * 100)) : 0
const fmt = (value: number | null | undefined) => (value == null ? '-' : Number(value).toLocaleString('th-TH', { maximumFractionDigits: 2 }))

// ส่งออกตารางที่เปิดอยู่เป็นไฟล์ .csv (เปิดใน Excel ได้ รองรับภาษาไทยด้วย BOM)
const exportCsv = () => {
  const quote = (value: unknown) => `"${String(value ?? '').replace(/"/g, '""')}"`
  let header: string[] = []
  let lines: unknown[][] = []
  if (tab.value === 'project') {
    header = ['โครงการ', 'วิธีคิด ไฟ/น้ำ', 'จดแล้ว', 'ทั้งหมด', 'ผิดปกติ', 'ข้าม', 'ไฟ (หน่วย)', 'น้ำ (ลบ.ม.)', 'สถานะ']
    lines = sortedProjects.value.map(row => [row.name, row.method, row.recorded, row.total, row.abnormal, row.skipped, row.electricUnits, row.waterUnits, row.ready ? 'ครบ · พร้อมออกบิล' : 'ยังไม่ครบ'])
  } else {
    header = ['แผง', 'ผู้เช่า', 'ไฟ เดิม', 'ไฟ ใหม่', 'ไฟ ใช้ไป', 'น้ำ เดิม', 'น้ำ ใหม่', 'น้ำ ใช้ไป', 'ยอดประมาณ', 'จดโดย', 'สถานะ']
    lines = visibleStalls.value.map(row => [row.label, row.tenant, row.electricPrev, row.electricCurr ?? '', row.electricUsed ?? '', row.waterPrev, row.waterCurr ?? '', row.waterUsed ?? '', row.amount ?? '', row.recordedBy, stateLabel[row.state]])
  }
  const csv = '﻿' + [header, ...lines].map(line => line.map(quote).join(',')).join('\r\n')
  const url = URL.createObjectURL(new Blob([csv], { type: 'text/csv;charset=utf-8' }))
  const link = document.createElement('a')
  link.href = url
  link.download = `meter-summary-${tab.value}.csv`
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}
</script>

<style scoped>
.msm-toolbar { display: flex; align-items: center; justify-content: space-between; gap: 12px; flex-wrap: wrap; margin: -4px 0 12px; }
.msm-toolbar-actions { display: flex; align-items: center; gap: 10px; }
.msm-project-select { width: 200px; }
.msm-stall-pane { display: flex; flex-direction: column; gap: 10px; }

.msm-table-wrap { overflow-x: auto; -webkit-overflow-scrolling: touch; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; }
.msm-table { width: 100%; min-width: 860px; border-collapse: collapse; }
.msm-table th { position: sticky; top: 0; padding: 12px 14px; background: var(--bg-f8fafc); color: var(--tx-475569); font-size: 0.8125rem; font-weight: 600; text-align: left; white-space: nowrap; border-bottom: 1px solid var(--bd-e2e8f0); }
.msm-table td { padding: 12px 14px; color: var(--tx-334155); font-size: 0.875rem; font-weight: 500; border-bottom: 1px solid var(--bd-f1f5f9); vertical-align: middle; }
.msm-table tbody tr:last-child td { border-bottom: 0; }
.msm-table tbody tr:hover td { background: var(--bg-f8fafc); }
.msm-table .num { text-align: right; white-space: nowrap; font-variant-numeric: tabular-nums; }
.msm-table .act { width: 1%; text-align: right; white-space: nowrap; }
.msm-table .strong { color: var(--tx-0f172a); font-weight: 700; }
.msm-table .muted { color: var(--tx-64748b); }
.msm-table .is-danger { color: var(--tx-dc2626); font-weight: 700; }
.msm-table .is-zero { color: var(--tx-94a3b8); }

.msm-progress { display: flex; align-items: center; gap: 10px; min-width: 190px; }
.msm-progress-track { flex: 1 1 auto; height: 8px; overflow: hidden; border-radius: 999px; background: var(--bg-e2e8f0); }
.msm-progress-track i { display: block; height: 100%; border-radius: 999px; background: #2563eb; transition: width .3s ease; }
.msm-progress-text { flex: 0 0 auto; min-width: 62px; text-align: right; color: var(--tx-475569); font-size: 0.8125rem; font-weight: 600; font-variant-numeric: tabular-nums; }

.msm-pill { display: inline-block; padding: 4px 12px; border-radius: 8px; font-size: 0.8125rem; font-weight: 600; white-space: nowrap; }
.msm-pill.is-pending { background: var(--bg-fef3c7); color: var(--tx-b45309); }
.msm-pill.is-ready, .msm-pill.is-recorded { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.msm-pill.is-abnormal { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.msm-pill.is-skipped { background: var(--bg-f1f5f9); color: var(--tx-64748b); }

.msm-stall-name { display: block; color: var(--tx-0f172a); font-size: 0.9375rem; font-weight: 700; }
.msm-stall-tenant { display: block; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }

.msm-footer-note { flex: 1 1 auto; margin: 0; color: var(--tx-475569); font-size: 0.8125rem; font-weight: 500; line-height: 1.5; }
.sr-only { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0; }

@media (max-width: 640px) {
  .msm-toolbar { flex-direction: column; align-items: stretch; }
  .msm-toolbar-actions { flex-wrap: wrap; }
  .msm-project-select { width: 100%; }
}
</style>
