<template>
  <div
    ref="rootRef"
    class="app-time-picker"
    :class="[
      `size-${size}`,
      {
        'is-open': isOpen,
        'is-disabled': disabled,
        'has-value': !!modelValue
      }
    ]"
  >
    <!-- Time Input Trigger -->
    <div
      class="time-picker-trigger"
      :class="{ 'is-focused': isOpen, 'is-disabled': disabled }"
      @click="toggleDropdown"
      role="button"
      tabindex="0"
      :aria-expanded="isOpen"
      :aria-label="ariaLabel || placeholder || 'เลือกเวลา'"
      @keydown.enter.prevent="toggleDropdown"
      @keydown.space.prevent="toggleDropdown"
      @keydown.esc.prevent="closeDropdown"
    >
      <span class="time-picker-icon">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <polyline points="12 6 12 12 16 14"></polyline>
        </svg>
      </span>

      <span class="time-picker-display" :class="{ 'placeholder-text': !modelValue }">
        {{ formattedDisplay || placeholder }}
      </span>

      <!-- Clear button if value exists -->
      <button
        v-if="clearable && modelValue && !disabled"
        type="button"
        class="time-picker-clear"
        @click.stop="clearValue"
        title="ล้างเวลา"
        aria-label="ล้างเวลา"
      >
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <!-- Dropdown Popover -->
    <transition name="time-picker-fade">
      <div v-if="isOpen" class="time-picker-popover" role="dialog" aria-modal="true">
        <!-- Quick Header -->
        <div class="popover-header">
          <div class="header-time-preview">
            <span class="preview-number">{{ tempHour }}</span>
            <span class="preview-colon">:</span>
            <span class="preview-number">{{ tempMinute }}</span>
            <span class="preview-unit">น.</span>
          </div>
          <button type="button" class="close-popover-btn" @click="closeDropdown" title="ปิด">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Quick Preset Chips -->
        <div v-if="showPresets" class="popover-presets">
          <button
            v-for="preset in presets"
            :key="preset"
            type="button"
            class="preset-chip"
            :class="{ active: modelValue === preset }"
            @click="selectPreset(preset)"
          >
            {{ preset }}
          </button>
        </div>

        <!-- Two Columns: Hours & Minutes -->
        <div class="popover-columns">
          <!-- Hours Column -->
          <div class="time-column" ref="hourColRef">
            <div class="column-title">ชั่วโมง</div>
            <div class="column-scroll">
              <button
                v-for="h in availableHours"
                :key="h"
                type="button"
                class="time-item"
                :class="{ active: tempHour === h }"
                @click="selectHour(h)"
              >
                {{ h }}
              </button>
            </div>
          </div>

          <!-- Minutes Column -->
          <div class="time-column" ref="minuteColRef">
            <div class="column-title">นาที</div>
            <div class="column-scroll">
              <button
                v-for="m in availableMinutes"
                :key="m"
                type="button"
                class="time-item"
                :class="{ active: tempMinute === m }"
                @click="selectMinute(m)"
              >
                {{ m }}
              </button>
            </div>
          </div>
        </div>

        <!-- Popover Actions Footer -->
        <div class="popover-footer">
          <button type="button" class="btn-time-now" @click="selectNow">
            เวลาปัจจุบัน
          </button>
          <button type="button" class="btn-time-confirm" @click="confirmSelection">
            ตกลง
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string // formatted as 'HH:mm' e.g. '09:30'
    placeholder?: string
    disabled?: boolean
    clearable?: boolean
    size?: 'sm' | 'md'
    minuteStep?: number // e.g. 5, 10, 15
    showPresets?: boolean
    presets?: string[]
    ariaLabel?: string
  }>(),
  {
    modelValue: '',
    placeholder: '--:--',
    disabled: false,
    clearable: true,
    size: 'sm',
    minuteStep: 5,
    showPresets: true,
    presets: () => [
      '09:00', '09:30', '10:00', '10:30',
      '11:00', '11:30', '13:00', '13:30',
      '14:00', '15:30', '16:30', '17:00'
    ],
    ariaLabel: undefined
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}>()

const rootRef = ref<HTMLElement | null>(null)
const hourColRef = ref<HTMLElement | null>(null)
const minuteColRef = ref<HTMLElement | null>(null)

const isOpen = ref(false)
const tempHour = ref('09')
const tempMinute = ref('00')

// Parse incoming modelValue e.g. '09:30' or '09.30'
const parseTime = (val: string) => {
  if (!val) return { hour: '09', minute: '00', valid: false }
  const clean = val.replace('.', ':').trim()
  const parts = clean.split(':')
  if (parts.length >= 2) {
    const h = parts[0].padStart(2, '0')
    const m = parts[1].slice(0, 2).padStart(2, '0')
    return { hour: h, minute: m, valid: true }
  }
  return { hour: '09', minute: '00', valid: false }
}

// Watch modelValue changes to sync internal temp values
watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      const parsed = parseTime(val)
      tempHour.value = parsed.hour
      tempMinute.value = parsed.minute
    }
  },
  { immediate: true }
)

const formattedDisplay = computed(() => {
  if (!props.modelValue) return ''
  const parsed = parseTime(props.modelValue)
  return `${parsed.hour}:${parsed.minute}`
})

// Available hours (06 to 22 by default, can expand 00-23)
const availableHours = computed(() => {
  const hours: string[] = []
  for (let i = 6; i <= 22; i++) {
    hours.push(i.toString().padStart(2, '0'))
  }
  return hours
})

// Available minutes based on step
const availableMinutes = computed(() => {
  const minutes: string[] = []
  const step = props.minuteStep || 5
  for (let i = 0; i < 60; i += step) {
    minutes.push(i.toString().padStart(2, '0'))
  }
  return minutes
})

const toggleDropdown = () => {
  if (props.disabled) return
  if (isOpen.value) {
    closeDropdown()
  } else {
    openDropdown()
  }
}

const openDropdown = () => {
  isOpen.value = true
  if (props.modelValue) {
    const parsed = parseTime(props.modelValue)
    tempHour.value = parsed.hour
    tempMinute.value = parsed.minute
  } else {
    tempHour.value = '09'
    tempMinute.value = '30'
  }
  nextTick(() => {
    scrollToActive()
  })
}

const closeDropdown = () => {
  isOpen.value = false
}

const scrollToActive = () => {
  if (hourColRef.value) {
    const activeHourBtn = hourColRef.value.querySelector('.time-item.active') as HTMLElement
    if (activeHourBtn) {
      activeHourBtn.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
    }
  }
  if (minuteColRef.value) {
    const activeMinBtn = minuteColRef.value.querySelector('.time-item.active') as HTMLElement
    if (activeMinBtn) {
      activeMinBtn.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
    }
  }
}

const selectHour = (h: string) => {
  tempHour.value = h
  emitValue()
}

const selectMinute = (m: string) => {
  tempMinute.value = m
  emitValue()
}

const selectPreset = (preset: string) => {
  const parsed = parseTime(preset)
  tempHour.value = parsed.hour
  tempMinute.value = parsed.minute
  emit('update:modelValue', `${parsed.hour}:${parsed.minute}`)
  emit('change', `${parsed.hour}:${parsed.minute}`)
  closeDropdown()
}

const selectNow = () => {
  const now = new Date()
  const h = now.getHours().toString().padStart(2, '0')
  const step = props.minuteStep || 5
  const rawMin = now.getMinutes()
  const roundedMin = Math.round(rawMin / step) * step
  const m = (roundedMin >= 60 ? 55 : roundedMin).toString().padStart(2, '0')

  tempHour.value = h
  tempMinute.value = m
  emitValue()
  closeDropdown()
}

const confirmSelection = () => {
  emitValue()
  closeDropdown()
}

const emitValue = () => {
  const val = `${tempHour.value}:${tempMinute.value}`
  emit('update:modelValue', val)
  emit('change', val)
}

const clearValue = () => {
  emit('update:modelValue', '')
  emit('change', '')
}

// Click outside handling
const handleOutsideClick = (e: MouseEvent) => {
  if (rootRef.value && !rootRef.value.contains(e.target as Node)) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleOutsideClick)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleOutsideClick)
})
</script>

<style scoped>
.app-time-picker {
  position: relative;
  display: inline-block;
  user-select: none;
  font-family: inherit;
}

/* Trigger Box */
.time-picker-trigger {
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 8px;
  color: var(--tx-1e293b);
  cursor: pointer;
  transition: all 0.18s ease;
  box-sizing: border-box;
}

.time-picker-trigger:hover:not(.is-disabled) {
  border-color: var(--bd-93c5fd);
  background: var(--bg-f8fafc);
}

.time-picker-trigger.is-focused {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.16);
  background: var(--bg-ffffff);
}

.time-picker-trigger.is-disabled {
  background: var(--bg-f1f5f9);
  color: var(--tx-94a3b8);
  cursor: not-allowed;
  border-color: var(--bd-e2e8f0);
}

/* Sizing */
.size-sm .time-picker-trigger {
  height: 34px;
  padding: 0 8px;
  min-width: 96px;
}

.size-sm .time-picker-display {
  font-size: 0.82rem;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.size-md .time-picker-trigger {
  height: 40px;
  padding: 0 11px;
  min-width: 110px;
}

.size-md .time-picker-display {
  font-size: 0.88rem;
  font-weight: 600;
  letter-spacing: 0.4px;
}

.time-picker-icon {
  display: inline-flex;
  align-items: center;
  color: var(--tx-64748b);
  flex-shrink: 0;
}

.time-picker-trigger.is-focused .time-picker-icon {
  color: var(--tx-2563eb);
}

.time-picker-display {
  flex: 1;
  text-align: center;
  color: var(--tx-1e293b);
  font-variant-numeric: tabular-nums;
}

.time-picker-display.placeholder-text {
  color: var(--tx-94a3b8);
  font-weight: 400;
}

.time-picker-clear {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: var(--tx-94a3b8);
  padding: 2px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.time-picker-clear:hover {
  color: var(--tx-ef4444);
  background: var(--bg-fee2e2);
}

/* Popover Panel */
.time-picker-popover {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  z-index: 100;
  width: 250px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  box-shadow: 0 12px 28px -4px rgba(15, 23, 42, 0.16), 0 4px 10px -2px rgba(15, 23, 42, 0.08);
  padding: 12px;
  box-sizing: border-box;
}

/* Header */
.popover-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--bd-f1f5f9);
}

.header-time-preview {
  display: flex;
  align-items: baseline;
  gap: 3px;
  font-weight: 700;
  color: var(--tx-1e293b);
}

.preview-number {
  font-size: 1.25rem;
  color: var(--tx-2563eb);
  font-variant-numeric: tabular-nums;
}

.preview-colon {
  font-size: 1.1rem;
  color: var(--tx-64748b);
}

.preview-unit {
  font-size: 0.8rem;
  color: var(--tx-64748b);
  margin-left: 2px;
  font-weight: 500;
}

.close-popover-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 6px;
  border: none;
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
  cursor: pointer;
  transition: all 0.15s ease;
}

.close-popover-btn:hover {
  background: var(--bg-e2e8f0);
  color: var(--tx-0f172a);
}

/* Presets */
.popover-presets {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 5px;
  margin-top: 10px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--bd-f1f5f9);
}

.preset-chip {
  padding: 4px 2px;
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--tx-475569);
  cursor: pointer;
  text-align: center;
  transition: all 0.15s ease;
}

.preset-chip:hover {
  background: var(--bg-eff6ff);
  border-color: var(--bd-93c5fd);
  color: var(--tx-2563eb);
}

.preset-chip.active {
  background: #2563eb;
  border-color: #2563eb;
  color: #ffffff;
}

/* Two Columns */
.popover-columns {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-top: 10px;
}

.time-column {
  display: flex;
  flex-direction: column;
}

.column-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--tx-94a3b8);
  text-align: center;
  padding-bottom: 4px;
}

.column-scroll {
  max-height: 140px;
  overflow-y: auto;
  border: 1px solid var(--bd-f1f5f9);
  border-radius: 8px;
  padding: 4px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  scrollbar-width: thin;
}

.time-item {
  width: 100%;
  padding: 6px 0;
  text-align: center;
  border: none;
  background: transparent;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--tx-334155);
  cursor: pointer;
  transition: all 0.12s ease;
  font-variant-numeric: tabular-nums;
}

.time-item:hover {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}

.time-item.active {
  background: #2563eb;
  color: #ffffff;
  font-weight: 700;
}

/* Footer */
.popover-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 10px;
  padding-top: 8px;
  border-top: 1px solid var(--bd-f1f5f9);
}

.btn-time-now {
  background: transparent;
  border: none;
  color: var(--tx-2563eb);
  font-size: 0.78rem;
  font-weight: 600;
  cursor: pointer;
  padding: 4px 6px;
  border-radius: 6px;
  transition: all 0.15s ease;
}

.btn-time-now:hover {
  background: var(--bg-eff6ff);
}

.btn-time-confirm {
  background: #2563eb;
  color: #ffffff;
  border: none;
  border-radius: 6px;
  font-size: 0.78rem;
  font-weight: 600;
  padding: 5px 14px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-time-confirm:hover {
  background: #1d4ed8;
}

/* Transition */
.time-picker-fade-enter-active,
.time-picker-fade-leave-active {
  transition: all 0.16s ease;
}

.time-picker-fade-enter-from,
.time-picker-fade-leave-to {
  opacity: 0;
  transform: translateY(-6px) scale(0.97);
}
</style>
