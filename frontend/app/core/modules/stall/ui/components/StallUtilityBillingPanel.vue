<template>
  <section class="cp-panel" aria-label="รายการบิลค่าน้ำค่าไฟ">
    <header class="cp-head">
      <div>
        <h2>ค่าน้ำ ค่าไฟ</h2>
        <p>{{ rows.length.toLocaleString('th-TH') }} แผง · คลิกแถวหรือคลิกแผงบนแผนที่เพื่อจดมิเตอร์</p>
      </div>
      <button type="button" class="cp-close" aria-label="ปิด" @click="$emit('close')">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 6 12 12M18 6 6 18"></path></svg>
      </button>
    </header>

    <div class="cp-toolbar">
      <div class="cp-month"><AppSelect v-model="billingMonth" :options="monthOptions" /></div>
      <AppButton variant="outline" size="md" @click="openRates">อัตราค่าบริการ</AppButton>
      <label class="cp-search">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
        <input v-model.trim="search" type="search" placeholder="ค้นหาแผง / ผู้เช่า" />
      </label>
    </div>

    <div class="cp-table-wrap">
      <table class="cp-table">
        <thead>
          <tr><th>แผง</th><th class="num">ค่าน้ำ</th><th class="num">ค่าไฟ</th><th class="num">รวม</th><th>สถานะ</th></tr>
        </thead>
        <tbody>
          <tr v-for="r in pagedRows" :key="r.stall" :class="{ active: r.stall === activeCode }" tabindex="0" @click="$emit('select', r.stall)" @keydown.enter="$emit('select', r.stall)">
            <td><strong>{{ r.stall }}</strong><small>{{ r.tenant }}</small></td>
            <td class="num">
              <template v-if="waterUnits(r) != null">{{ money(waterUnits(r)! * rates.water) }}<small>{{ waterUnits(r) }} หน่วย</small></template>
              <span v-else class="muted">-</span>
            </td>
            <td class="num">
              <template v-if="electricUnits(r) != null">{{ money(electricUnits(r)! * rates.electric) }}<small>{{ electricUnits(r) }} หน่วย</small></template>
              <span v-else class="muted">-</span>
            </td>
            <td class="num"><strong>{{ utilityTotal(r) == null ? '-' : money(utilityTotal(r)!) }}</strong></td>
            <td><span class="cp-badge" :class="`is-${r.status}`">{{ BILL_STATUS_LABEL[r.status] }}</span></td>
          </tr>
          <tr v-if="!pagedRows.length"><td colspan="5" class="empty">ไม่พบรายการที่ตรงกับตัวกรอง</td></tr>
        </tbody>
      </table>
    </div>

    <AppPagination
      v-model="page"
      :total-items="rows.length"
      :items-per-page="pageSize"
      :page-size-options="[5, 10, 20, 50]"
      unit-label="แผง"
      @update:items-per-page="val => { pageSize = val; page = 1 }"
    />

    <AppModal v-model="showRates" title="อัตราค่าบริการ" description="ใช้คำนวณบิลของทุกแผงในตลาดนี้">
      <div class="cp-form">
        <label>ค่าน้ำ (บาท/หน่วย)<AppInput v-model="rateForm.water" type="number" /></label>
        <label>ค่าไฟ (บาท/หน่วย)<AppInput v-model="rateForm.electric" type="number" /></label>
      </div>
      <template #footer>
        <AppButton variant="secondary" @click="showRates = false">ยกเลิก</AppButton>
        <AppButton variant="primary" @click="saveRates">บันทึก</AppButton>
      </template>
    </AppModal>
  </section>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import AppButton from '~/components/AppButton.vue'
import AppInput from '~/components/AppInput.vue'
import AppModal from '~/components/AppModal.vue'
import AppPagination from '~/components/AppPagination.vue'
import AppSelect from '~/components/AppSelect.vue'
import { BILL_STATUS_LABEL, useStallBilling, type CostFilter } from '../../composables/useStallBilling'

// ข้อมูลเป็นข้อมูลจำลอง (mock) ยังไม่เชื่อม API มิเตอร์/บิล
const props = defineProps<{ filter: CostFilter; activeCode?: string }>()
defineEmits<{ close: []; select: [code: string] }>()

const { rates, billingMonth, utilityList, waterUnits, electricUnits, utilityTotal, matchesFilter, matchesAmount, amountRange } = useStallBilling()

const monthOptions = computed(() => {
  const base = new Date()
  return Array.from({ length: 3 }, (_, i) => {
    const d = new Date(base.getFullYear(), base.getMonth() - i, 1)
    return {
      value: `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`,
      label: new Intl.DateTimeFormat('th-TH', { month: 'long', year: 'numeric' }).format(d)
    }
  })
})

const money = (v: number) => v.toLocaleString('th-TH', { maximumFractionDigits: 2 })
const search = ref('')
const page = ref(1)
const pageSize = ref(10)
const rows = computed(() => {
  const q = search.value.toLowerCase()
  return utilityList.value.filter(r =>
    matchesFilter(r.status, props.filter) && matchesAmount(utilityTotal(r)) &&
    (!q || r.stall.toLowerCase().includes(q) || r.tenant.toLowerCase().includes(q))
  )
})
const pagedRows = computed(() => rows.value.slice((page.value - 1) * pageSize.value, page.value * pageSize.value))
watch([() => props.filter, search, billingMonth, amountRange], () => { page.value = 1 })

const showRates = ref(false)
const rateForm = reactive({ water: '', electric: '' })
const openRates = () => {
  rateForm.water = String(rates.water)
  rateForm.electric = String(rates.electric)
  showRates.value = true
}
const saveRates = () => {
  if (Number(rateForm.water) >= 0) rates.water = Number(rateForm.water)
  if (Number(rateForm.electric) >= 0) rates.electric = Number(rateForm.electric)
  showRates.value = false
}
</script>

<style scoped>
.cp-panel { display: flex; flex-direction: column; min-height: 0; height: 100%; color: #334155; }
.cp-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; padding: 16px 16px 12px; border-bottom: 1px solid #f1f5f9; }
.cp-head h2 { margin: 0; color: #0f172a; font-size: 1.125rem; font-weight: 700; }
.cp-head p { margin: 2px 0 0; color: #64748b; font-size: .8125rem; font-weight: 500; }
.cp-close { flex: 0 0 auto; width: 34px; height: 34px; display: grid; place-items: center; border: 1px solid #e2e8f0; border-radius: 8px; background: #fff; color: #475569; cursor: pointer; }
.cp-close:hover { background: #f1f5f9; }
.cp-toolbar { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; padding: 12px 16px; }
.cp-month { width: 160px; }
.cp-search { position: relative; display: flex; align-items: center; flex: 1 1 140px; min-width: 0; }
.cp-search svg { position: absolute; left: 12px; color: #94a3b8; }
.cp-search input { box-sizing: border-box; width: 100%; height: 40px; padding: 0 12px 0 36px; border: 1px solid #e2e8f0; border-radius: 8px; font-family: inherit; font-size: .9375rem; font-weight: 500; color: #1e293b; outline: none; }
.cp-search input:focus { border-color: #2563eb; box-shadow: 0 0 0 3px rgba(37,99,235,.12); }
.cp-table-wrap { flex: 1 1 auto; min-height: 0; overflow: auto; padding: 0 16px; -webkit-overflow-scrolling: touch; }
.cp-table { width: 100%; border-collapse: collapse; }
.cp-table th { position: sticky; top: 0; background: #fff; padding: 12px 10px; text-align: left; color: #64748b; font-size: .8125rem; font-weight: 500; border-bottom: 1px solid #f1f5f9; white-space: nowrap; }
.cp-table td { padding: 12px 10px; font-size: .875rem; font-weight: 500; border-bottom: 1px solid #f1f5f9; vertical-align: middle; }
.cp-table tbody tr { cursor: pointer; }
.cp-table tbody tr:hover, .cp-table tbody tr.active { background: #f8fafc; }
.cp-table tbody tr.active { box-shadow: inset 3px 0 0 #2563eb; }
.cp-table tbody tr:focus-visible { outline: 2px solid #93c5fd; outline-offset: -2px; }
.cp-table .num { text-align: right; white-space: nowrap; }
.cp-table td small { display: block; color: #94a3b8; font-size: .8125rem; font-weight: 500; }
.cp-table td strong { color: #1e293b; font-weight: 700; }
.muted { color: #94a3b8; }
.empty { text-align: center; color: #94a3b8; padding: 28px 10px !important; cursor: default; }
.cp-badge { display: inline-block; padding: 4px 10px; border-radius: 8px; font-size: .8125rem; font-weight: 600; white-space: nowrap; }
.cp-badge.is-pending { background: #fef3c7; color: #d97706; }
.cp-badge.is-paid { background: #dcfce7; color: #15803d; }
.cp-badge.is-overdue { background: #fee2e2; color: #dc2626; }
.cp-badge.is-unrecorded { background: #f1f5f9; color: #64748b; }
.cp-form { display: flex; flex-direction: column; gap: 14px; }
.cp-form label { display: flex; flex-direction: column; gap: 6px; color: #475569; font-size: .875rem; font-weight: 600; }
.cp-panel :deep(.app-pagination), .cp-panel > :deep(div:last-of-type) { flex: 0 0 auto; }
</style>
