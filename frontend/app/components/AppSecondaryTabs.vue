<template>
  <div class="secondary-tabs-scroll">
    <div class="secondary-tabs" role="tablist" :aria-label="ariaLabel" @keydown="handleKeydown">
      <button
        v-for="option in options"
        :key="option.value"
        :ref="(element) => setTabRef(option.value, element)"
        type="button"
        class="secondary-tab"
        :class="{ active: option.value === modelValue }"
        role="tab"
        :aria-selected="option.value === modelValue"
        :tabindex="option.value === modelValue ? 0 : -1"
        :disabled="option.disabled"
        @click="selectTab(option)"
      >
        <span>{{ option.label }}</span>
        <span v-if="option.count !== undefined" class="secondary-tab-count">{{ option.count }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ComponentPublicInstance } from 'vue'

export interface AppSecondaryTabOption {
  value: string
  label: string
  count?: string | number
  disabled?: boolean
}

const props = withDefaults(defineProps<{
  modelValue: string
  options: AppSecondaryTabOption[]
  ariaLabel?: string
}>(), { ariaLabel: 'แท็บรอง' })

const emit = defineEmits<{
  'update:modelValue': [value: string]
  change: [value: string]
}>()
const tabRefs = new Map<string, HTMLButtonElement>()

const setTabRef = (value: string, element: Element | ComponentPublicInstance | null) => {
  if (element instanceof HTMLButtonElement) tabRefs.set(value, element)
}
const selectTab = (option: AppSecondaryTabOption) => {
  if (option.disabled || option.value === props.modelValue) return
  emit('update:modelValue', option.value)
  emit('change', option.value)
}
const handleKeydown = (event: KeyboardEvent) => {
  if (!['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key)) return
  const enabled = props.options.filter(option => !option.disabled)
  if (!enabled.length) return
  const currentIndex = Math.max(0, enabled.findIndex(option => option.value === props.modelValue))
  let nextIndex = currentIndex
  if (event.key === 'ArrowRight') nextIndex = (currentIndex + 1) % enabled.length
  if (event.key === 'ArrowLeft') nextIndex = (currentIndex - 1 + enabled.length) % enabled.length
  if (event.key === 'Home') nextIndex = 0
  if (event.key === 'End') nextIndex = enabled.length - 1
  event.preventDefault()
  const next = enabled[nextIndex]
  selectTab(next)
  nextTick(() => tabRefs.get(next.value)?.focus())
}
</script>

<style scoped>
.secondary-tabs-scroll{max-width:100%;overflow-x:auto;overflow-y:hidden;scrollbar-width:none;-webkit-overflow-scrolling:touch}.secondary-tabs-scroll::-webkit-scrollbar{display:none}.secondary-tabs{display:inline-flex;align-items:flex-end;gap:24px;min-width:max-content;border-bottom:1px solid var(--bd-dbe5f1)}.secondary-tab{position:relative;display:inline-flex;align-items:center;justify-content:center;gap:8px;min-height:44px;padding:0 4px 12px;border:0;border-radius:0;background:transparent;color:var(--tx-334155);font-family:inherit;font-size:var(--font-size-base);font-weight:600;line-height:1;white-space:nowrap;cursor:pointer;transition:color .18s ease}.secondary-tab:hover:not(:disabled):not(.active){color:var(--tx-3b82f6)}.secondary-tab.active{color:var(--tx-2563eb)}.secondary-tab.active::after{content:'';position:absolute;left:0;right:0;bottom:-1px;height:3px;border-radius:0;background:#2563eb}.secondary-tab:focus-visible{outline:3px solid rgba(37,99,235,.18);outline-offset:2px;border-radius:2px}.secondary-tab:disabled{cursor:not-allowed;color:var(--tx-cbd5e1)}.secondary-tab-count{display:inline-flex;align-items:center;justify-content:center;min-width:23px;height:23px;padding:0 7px;border-radius:999px;background:var(--bg-e2e8f0);color:var(--tx-475569);font-size:var(--font-size-xs);font-weight:700}.secondary-tab.active .secondary-tab-count{background:#2563eb;color:#fff}.secondary-tab:disabled .secondary-tab-count{background:var(--bg-f1f5f9);color:var(--tx-cbd5e1)}@media(max-width:640px){.secondary-tabs{gap:18px}.secondary-tab{min-height:42px;padding-bottom:10px;font-size:var(--font-size-sm)}.secondary-tab-count{min-width:22px;height:22px}}
</style>
