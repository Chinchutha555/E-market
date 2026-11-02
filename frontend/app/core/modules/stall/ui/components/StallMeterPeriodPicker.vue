<template>
  <div ref="periodRef" class="mpp-period" :class="{ 'is-past': !isCurrentPeriod }" role="group" aria-label="งวดการจดมิเตอร์">
    <button type="button" class="mpp-period-nav" :disabled="!olderMonth" aria-label="งวดก่อนหน้า" @click="pickMonth(olderMonth)">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="15 18 9 12 15 6"></polyline></svg>
    </button>
    <div class="mpp-period-main">
      <small>งวดการจด<span v-if="isCurrentPeriod" class="mpp-period-badge">งวดปัจจุบัน</span><span v-else class="mpp-period-badge is-past">งวดย้อนหลัง</span></small>
      <button type="button" class="mpp-period-pick" :aria-expanded="pickerOpen" aria-haspopup="dialog" aria-label="เลือกปีและเดือน" @click="togglePicker">
        <strong>{{ periodLabel }}</strong>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" :class="{ 'is-open': pickerOpen }"><polyline points="6 9 12 15 18 9"></polyline></svg>
      </button>
      <span v-if="due" class="mpp-period-due">ครบกำหนด {{ due }}</span>
    </div>
    <button type="button" class="mpp-period-nav" :disabled="!newerMonth" aria-label="งวดถัดไป" @click="pickMonth(newerMonth)">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 18 15 12 9 6"></polyline></svg>
    </button>

    <div v-if="pickerOpen" class="mpp-picker" role="dialog" aria-label="เลือกปีและเดือนของงวด">
      <div class="mpp-picker-year">
        <button type="button" class="mpp-picker-nav" :disabled="pickerYear <= minYear" aria-label="ปีก่อนหน้า" @click="pickerYear -= 1">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="15 18 9 12 15 6"></polyline></svg>
        </button>
        <strong>ปี {{ pickerYear + 543 }}</strong>
        <button type="button" class="mpp-picker-nav" :disabled="pickerYear >= currentYear" aria-label="ปีถัดไป" @click="pickerYear += 1">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 18 15 12 9 6"></polyline></svg>
        </button>
      </div>
      <div class="mpp-picker-grid">
        <button
          v-for="item in pickerMonths"
          :key="item.value"
          type="button"
          class="mpp-picker-month"
          :class="{ 'is-selected': item.value === month, 'is-current': item.value === currentMonth }"
          :disabled="item.disabled"
          :aria-pressed="item.value === month"
          @click="selectFromPicker(item.value)"
        >{{ item.label }}</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

// ตัวเลือกงวดการจดมิเตอร์ (เดือนบิล): ลูกศรเลื่อนทีละเดือน แตะชื่อเดือนเพื่อเลือกปีและเดือน ใช้ร่วมกันในการ์ดโครงการและหน้าสรุปการจด
const props = defineProps<{
  /** งวดที่เลือกอยู่ YYYY-MM */
  month: string
  /** งวดปัจจุบัน YYYY-MM (เลือกงวดหลังจากนี้ไม่ได้) */
  currentMonth: string
  /** วันครบกำหนดของงวด เช่น "30 พ.ย." */
  due?: string
}>()
const emit = defineEmits<{ 'update:month': [value: string] }>()

// งวด (เดือนบิล) รูปแบบ YYYY-MM ปุ่มลูกศรเลื่อนทีละเดือน แตะชื่อเดือนเพื่อเลือกปีและเดือน (เลือกงวดอนาคตไม่ได้)
const MIN_YEARS_BACK = 5
const parseMonth = (value?: string) => { const [y, m] = (value || '').split('-').map(Number); return { y: y || new Date().getFullYear(), m: m || 1 } }
const toMonth = (y: number, m: number) => `${y}-${String(m).padStart(2, '0')}`
const shiftMonth = (value: string, delta: number) => { const { y, m } = parseMonth(value); const d = new Date(y, m - 1 + delta, 1); return toMonth(d.getFullYear(), d.getMonth() + 1) }
const currentYear = computed(() => parseMonth(props.currentMonth).y)
const minYear = computed(() => currentYear.value - MIN_YEARS_BACK)
const minMonth = computed(() => toMonth(minYear.value, 1))
const olderMonth = computed(() => { const value = shiftMonth(props.month, -1); return value >= minMonth.value ? value : null })
const newerMonth = computed(() => { const value = shiftMonth(props.month, 1); return value <= props.currentMonth ? value : null })
const isCurrentPeriod = computed(() => props.month === props.currentMonth)
const monthName = (y: number, m: number, style: 'long' | 'short') => new Intl.DateTimeFormat('th-TH', { month: style }).format(new Date(y, m - 1, 1))
const periodLabel = computed(() => { const { y, m } = parseMonth(props.month); return `${monthName(y, m, 'long')} ${y + 543}` })
const pickMonth = (value: string | null) => { if (value) emit('update:month', value) }

const pickerOpen = ref(false)
const pickerYear = ref(parseMonth(props.month).y)
const periodRef = ref<HTMLElement | null>(null)
const pickerMonths = computed(() => Array.from({ length: 12 }, (_, index) => {
  const value = toMonth(pickerYear.value, index + 1)
  return { value, label: monthName(pickerYear.value, index + 1, 'short'), disabled: value > props.currentMonth || value < minMonth.value }
}))
const togglePicker = () => {
  if (!pickerOpen.value) pickerYear.value = parseMonth(props.month).y
  pickerOpen.value = !pickerOpen.value
}
const selectFromPicker = (value: string) => { pickMonth(value); pickerOpen.value = false }
const onDocumentPointer = (event: MouseEvent) => { if (pickerOpen.value && periodRef.value && !periodRef.value.contains(event.target as Node)) pickerOpen.value = false }
const onDocumentKey = (event: KeyboardEvent) => { if (event.key === 'Escape') pickerOpen.value = false }
onMounted(() => { document.addEventListener('mousedown', onDocumentPointer); document.addEventListener('keydown', onDocumentKey) })
onBeforeUnmount(() => { document.removeEventListener('mousedown', onDocumentPointer); document.removeEventListener('keydown', onDocumentKey) })
</script>

<style scoped>
.mpp-period { position: relative; display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 16px 12px; border: 1px solid var(--bd-dbeafe); border-radius: 14px; background: var(--bg-eff6ff); }
.mpp-period.is-past { border-color: var(--bd-e2e8f0); background: var(--bg-f1f5f9); }
.mpp-period-badge.is-past { background: var(--bg-e2e8f0); color: var(--tx-475569); }
.mpp-back-current { display: flex; align-items: center; justify-content: center; gap: 6px; width: 100%; margin: -4px 0 12px; padding: 8px 10px; border: 1px dashed var(--bd-cbd5e1); border-radius: 10px; background: var(--bg-ffffff); color: var(--tx-2563eb); font: inherit; font-size: var(--font-size-sm); font-weight: 600; cursor: pointer; }
.mpp-back-current:hover { background: var(--bg-eff6ff); }
.mpp-period-nav { display: grid; place-items: center; flex: 0 0 34px; width: 34px; height: 34px; border: 1px solid var(--bd-dbeafe); border-radius: 9px; background: var(--bg-ffffff); color: var(--tx-2563eb); cursor: pointer; transition: background-color .15s ease, border-color .15s ease; }
.mpp-period-nav:hover:not(:disabled) { border-color: #93c5fd; background: var(--bg-eff6ff); }
.mpp-period-nav:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 2px; }
.mpp-period-nav:disabled { color: var(--tx-cbd5e1); cursor: not-allowed; }
.mpp-period-main { min-width: 0; flex: 1 1 auto; display: flex; flex-direction: column; align-items: center; gap: 6px; text-align: center; }
.mpp-period-main small { display: inline-flex; align-items: center; gap: 6px; color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 600; }
.mpp-period-badge { padding: 2px 10px; border-radius: 999px; background: var(--bg-dbeafe); color: var(--tx-2563eb); font-size: 0.75rem; font-weight: 700; }
.mpp-period-main strong { color: var(--tx-0f172a); font-size: var(--font-size-lg); font-weight: 800; line-height: 1.35; }
.mpp-period-pick { display: inline-flex; align-items: center; gap: 6px; padding: 2px 6px; border: 0; border-radius: 8px; background: transparent; color: inherit; font: inherit; cursor: pointer; transition: background-color .15s ease; }
.mpp-period-pick:hover { background: rgba(37, 99, 235, 0.08); }
.mpp-period-pick:focus-visible { outline: 2px solid var(--bd-93c5fd); outline-offset: 2px; }
.mpp-period-pick svg { color: var(--tx-2563eb); transition: transform .18s ease; }
.mpp-period-pick svg.is-open { transform: rotate(180deg); }
.mpp-picker { position: absolute; z-index: 20; top: calc(100% + 8px); left: 0; right: 0; padding: 12px; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-ffffff); box-shadow: 0 12px 30px rgba(15, 23, 42, 0.16); }
.mpp-picker-year { display: flex; align-items: center; justify-content: space-between; gap: 8px; margin-bottom: 10px; }
.mpp-picker-year strong { color: var(--tx-0f172a); font-size: var(--font-size-base); font-weight: 700; }
.mpp-picker-nav { display: grid; place-items: center; width: 30px; height: 30px; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-475569); cursor: pointer; }
.mpp-picker-nav:hover:not(:disabled) { border-color: #93c5fd; color: var(--tx-2563eb); background: var(--bg-eff6ff); }
.mpp-picker-nav:disabled { color: var(--tx-cbd5e1); cursor: not-allowed; }
.mpp-picker-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 6px; }
.mpp-picker-month { padding: 9px 4px; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-334155); font: inherit; font-size: var(--font-size-sm); font-weight: 600; cursor: pointer; transition: background-color .15s ease, border-color .15s ease; }
.mpp-picker-month:hover:not(:disabled):not(.is-selected) { border-color: #93c5fd; background: var(--bg-eff6ff); }
.mpp-picker-month.is-current:not(.is-selected) { border-color: #93c5fd; color: var(--tx-2563eb); }
.mpp-picker-month.is-selected { border-color: #2563eb; background: #2563eb; color: #ffffff; }
.mpp-picker-month:disabled { color: var(--tx-cbd5e1); background: var(--bg-f8fafc); cursor: not-allowed; }
.mpp-period-due { color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 500; }
</style>
