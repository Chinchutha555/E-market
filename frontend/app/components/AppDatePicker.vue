<template>
  <div
    ref="rootRef"
    class="app-date-picker"
    :class="[`size-${size}`, { 'is-open': isOpen, 'is-disabled': disabled, 'has-value': !!modelValue }]"
  >
    <div
      class="date-picker-trigger"
      :class="{ 'is-focused': isOpen, 'is-disabled': disabled }"
      role="button"
      tabindex="0"
      :aria-expanded="isOpen"
      :aria-label="ariaLabel || placeholder"
      @click="toggle"
      @keydown.enter.prevent="toggle"
      @keydown.space.prevent="toggle"
      @keydown.esc.prevent="close"
    >
      <span class="date-picker-icon" aria-hidden="true">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="3" y="4" width="18" height="18" rx="2"></rect>
          <line x1="16" y1="2" x2="16" y2="6"></line>
          <line x1="8" y1="2" x2="8" y2="6"></line>
          <line x1="3" y1="10" x2="21" y2="10"></line>
        </svg>
      </span>
      <span class="date-picker-display" :class="{ 'placeholder-text': !modelValue }">{{ displayText || placeholder }}</span>
      <button
        v-if="clearable && modelValue && !disabled"
        type="button"
        class="date-picker-clear"
        title="ล้างวันที่"
        aria-label="ล้างวันที่"
        @click.stop="clearValue"
      >
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <transition name="date-picker-fade">
      <div v-if="isOpen" class="date-picker-popover" :class="`align-${align}`" role="dialog" :aria-label="ariaLabel || placeholder" @keydown.esc.prevent="close">
        <div class="dp-header">
          <button type="button" class="dp-nav" aria-label="เดือนก่อนหน้า" @click="shiftMonth(-1)">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>
          </button>
          <strong class="dp-title">{{ monthTitle }}</strong>
          <button type="button" class="dp-nav" aria-label="เดือนถัดไป" @click="shiftMonth(1)">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
          </button>
        </div>

        <div class="dp-weekdays" aria-hidden="true">
          <span v-for="day in WEEKDAYS" :key="day">{{ day }}</span>
        </div>

        <div class="dp-grid" role="grid">
          <button
            v-for="cell in cells"
            :key="cell.iso"
            type="button"
            class="dp-day"
            :class="{ 'is-outside': !cell.inMonth, 'is-today': cell.iso === todayIso, 'is-selected': cell.iso === modelValue }"
            :disabled="cell.disabled"
            :aria-label="longLabel(cell.iso)"
            :aria-pressed="cell.iso === modelValue"
            @click="pick(cell.iso)"
          >{{ cell.day }}</button>
        </div>

        <div class="dp-footer">
          <button type="button" class="dp-link" @click="clearValue">ล้าง</button>
          <button type="button" class="dp-link" :disabled="todayDisabled" @click="pick(todayIso)">วันนี้</button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

// ตัวเลือกวันที่กลางของระบบ (ค่าเป็น 'YYYY-MM-DD' แสดงผลเป็นปี พ.ศ.) ใช้แทนช่องวันที่ของเบราว์เซอร์
const props = withDefaults(defineProps<{
  modelValue?: string
  placeholder?: string
  disabled?: boolean
  clearable?: boolean
  size?: 'sm' | 'md'
  min?: string
  max?: string
  align?: 'left' | 'right'
  ariaLabel?: string
}>(), {
  modelValue: '',
  placeholder: 'เลือกวันที่',
  disabled: false,
  clearable: true,
  size: 'md',
  min: '',
  max: '',
  align: 'left',
  ariaLabel: undefined
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}>()

const WEEKDAYS = ['อา', 'จ', 'อ', 'พ', 'พฤ', 'ศ', 'ส']
const rootRef = ref<HTMLElement | null>(null)
const isOpen = ref(false)

const pad = (n: number) => String(n).padStart(2, '0')
const toIso = (d: Date) => `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
const fromIso = (iso: string) => {
  const [y, m, d] = iso.split('-').map(Number)
  return new Date(y!, (m || 1) - 1, d || 1)
}
const isValidIso = (iso: string) => /^\d{4}-\d{2}-\d{2}$/.test(iso)

const todayIso = toIso(new Date())
const viewYear = ref(new Date().getFullYear())
const viewMonth = ref(new Date().getMonth())

const syncView = () => {
  const base = isValidIso(props.modelValue) ? fromIso(props.modelValue) : new Date()
  viewYear.value = base.getFullYear()
  viewMonth.value = base.getMonth()
}
watch(() => props.modelValue, syncView, { immediate: true })

const dateFmt = new Intl.DateTimeFormat('th-TH', { day: 'numeric', month: 'short', year: 'numeric' })
const monthFmt = new Intl.DateTimeFormat('th-TH', { month: 'long', year: 'numeric' })
const longFmt = new Intl.DateTimeFormat('th-TH', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })

const displayText = computed(() => (isValidIso(props.modelValue) ? dateFmt.format(fromIso(props.modelValue)) : ''))
const monthTitle = computed(() => monthFmt.format(new Date(viewYear.value, viewMonth.value, 1)))
const longLabel = (iso: string) => longFmt.format(fromIso(iso))

const isOutOfRange = (iso: string) => Boolean((props.min && iso < props.min) || (props.max && iso > props.max))
const todayDisabled = computed(() => isOutOfRange(todayIso))

// ตาราง 6 สัปดาห์ เริ่มวันอาทิตย์ (รวมวันของเดือนก่อน/ถัดไปแบบจาง)
const cells = computed(() => {
  const first = new Date(viewYear.value, viewMonth.value, 1)
  const start = new Date(viewYear.value, viewMonth.value, 1 - first.getDay())
  return Array.from({ length: 42 }, (_, i) => {
    const d = new Date(start.getFullYear(), start.getMonth(), start.getDate() + i)
    const iso = toIso(d)
    return { iso, day: d.getDate(), inMonth: d.getMonth() === viewMonth.value, disabled: isOutOfRange(iso) }
  })
})

const shiftMonth = (delta: number) => {
  const d = new Date(viewYear.value, viewMonth.value + delta, 1)
  viewYear.value = d.getFullYear()
  viewMonth.value = d.getMonth()
}

const open = () => { if (!props.disabled) { syncView(); isOpen.value = true } }
const close = () => { isOpen.value = false }
const toggle = () => (isOpen.value ? close() : open())

const setValue = (value: string) => {
  emit('update:modelValue', value)
  emit('change', value)
}
const pick = (iso: string) => {
  if (isOutOfRange(iso)) return
  setValue(iso)
  close()
}
const clearValue = () => {
  setValue('')
  close()
}

const onDocumentClick = (event: MouseEvent) => {
  if (!isOpen.value) return
  if (rootRef.value && !rootRef.value.contains(event.target as Node)) close()
}
onMounted(() => document.addEventListener('mousedown', onDocumentClick))
onBeforeUnmount(() => document.removeEventListener('mousedown', onDocumentClick))
</script>

<style scoped>
.app-date-picker { position: relative; display: block; width: 100%; min-width: 0; }
.date-picker-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  box-sizing: border-box;
  padding: 0 10px;
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 10px;
  background: var(--bg-ffffff);
  color: var(--tx-1e293b);
  cursor: pointer;
  transition: border-color 0.15s ease, box-shadow 0.15s ease, background-color 0.15s ease;
}
.size-sm .date-picker-trigger { height: 34px; border-radius: 8px; }
.size-md .date-picker-trigger { height: 42px; }
.date-picker-trigger:hover:not(.is-disabled) { border-color: var(--bd-93c5fd); }
.date-picker-trigger.is-focused { border-color: #2563eb; box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.14); }
.date-picker-trigger:focus-visible { outline: none; border-color: #2563eb; box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.14); }
.date-picker-trigger.is-disabled { background: var(--bg-f8fafc); color: var(--tx-64748b); cursor: not-allowed; }
.date-picker-icon { display: inline-flex; flex: 0 0 auto; color: var(--tx-2563eb); }
.date-picker-display { flex: 1 1 auto; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: var(--font-size-sm); font-weight: 500; }
.date-picker-display.placeholder-text { color: var(--tx-94a3b8); }
.date-picker-clear { display: inline-flex; align-items: center; justify-content: center; flex: 0 0 auto; width: 22px; height: 22px; padding: 0; border: 0; border-radius: 6px; background: transparent; color: var(--tx-64748b); cursor: pointer; }
.date-picker-clear:hover { background: var(--bg-f1f5f9); color: var(--tx-0f172a); }

.date-picker-popover {
  position: absolute;
  top: calc(100% + 6px);
  z-index: 100;
  width: 288px;
  max-width: calc(100vw - 24px);
  box-sizing: border-box;
  padding: 12px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  background: var(--bg-ffffff);
  box-shadow: 0 12px 28px -4px rgba(15, 23, 42, 0.16), 0 4px 10px -2px rgba(15, 23, 42, 0.08);
}
.date-picker-popover.align-left { left: 0; }
.date-picker-popover.align-right { right: 0; }

.dp-header { display: flex; align-items: center; justify-content: space-between; gap: 8px; padding-bottom: 8px; }
.dp-title { color: var(--tx-0f172a); font-size: var(--font-size-base); font-weight: 700; }
.dp-nav { display: inline-flex; align-items: center; justify-content: center; width: 32px; height: 32px; padding: 0; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-475569); cursor: pointer; transition: background-color 0.15s ease; }
.dp-nav:hover { background: var(--bg-f1f5f9); }

.dp-weekdays { display: grid; grid-template-columns: repeat(7, 1fr); margin-bottom: 2px; }
.dp-weekdays span { padding: 6px 0; text-align: center; color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 600; }

.dp-grid { display: grid; grid-template-columns: repeat(7, 1fr); gap: 2px; }
.dp-day { aspect-ratio: 1; display: grid; place-items: center; padding: 0; border: 1px solid transparent; border-radius: 8px; background: transparent; color: var(--tx-1e293b); font: inherit; font-size: var(--font-size-sm); font-weight: 500; cursor: pointer; transition: background-color 0.12s ease; }
.dp-day:hover:not(:disabled):not(.is-selected) { background: var(--bg-eff6ff); }
.dp-day.is-outside { color: var(--tx-94a3b8); }
.dp-day.is-today:not(.is-selected) { border-color: #93c5fd; color: var(--tx-2563eb); font-weight: 700; }
.dp-day.is-selected { background: #2563eb; color: #ffffff; font-weight: 700; }
.dp-day:disabled { color: var(--tx-cbd5e1); cursor: not-allowed; }
.dp-day:focus-visible { outline: 2px solid #93c5fd; outline-offset: 1px; }

.dp-footer { display: flex; align-items: center; justify-content: space-between; margin-top: 8px; padding-top: 8px; border-top: 1px solid var(--bd-f1f5f9); }
.dp-link { padding: 4px 8px; border: 0; border-radius: 6px; background: transparent; color: var(--tx-2563eb); font: inherit; font-size: var(--font-size-sm); font-weight: 600; cursor: pointer; }
.dp-link:hover:not(:disabled) { background: var(--bg-eff6ff); }
.dp-link:disabled { color: var(--tx-94a3b8); cursor: not-allowed; }

.date-picker-fade-enter-active, .date-picker-fade-leave-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.date-picker-fade-enter-from, .date-picker-fade-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
