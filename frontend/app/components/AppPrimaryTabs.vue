<template>
  <div class="primary-tabs-scroll">
    <div class="primary-tabs" role="tablist" :aria-label="ariaLabel" @keydown="handleKeydown">
      <button
        v-for="option in options"
        :key="option.value"
        :ref="(element) => setTabRef(option.value, element)"
        type="button"
        class="primary-tab"
        :class="{ active: option.value === modelValue }"
        role="tab"
        :aria-selected="option.value === modelValue"
        :tabindex="option.value === modelValue ? 0 : -1"
        :disabled="option.disabled"
        @click="selectTab(option)"
      >
        <svg v-if="option.icon === 'calendar'" viewBox="0 0 24 24" aria-hidden="true"><rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>
        <svg v-else-if="option.icon === 'heart'" viewBox="0 0 24 24" aria-hidden="true"><path d="M20.8 4.6a5.5 5.5 0 0 0-7.8 0L12 5.7l-1.1-1.1a5.5 5.5 0 0 0-7.8 7.8l1.1 1.1L12 21.2l7.8-7.8 1-1a5.5 5.5 0 0 0 0-7.8Z"/></svg>
        <svg v-else-if="option.icon === 'clock'" viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg>
        <svg v-else-if="option.icon === 'contract'" viewBox="0 0 24 24" aria-hidden="true"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8Z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg>
        <svg v-else-if="option.icon === 'history'" viewBox="0 0 24 24" aria-hidden="true"><path d="M3 12a9 9 0 1 0 3-6.7L3 8"/><path d="M3 3v5h5M12 7v5l3 2"/></svg>
        <svg v-else-if="option.icon === 'info'" viewBox="0 0 24 24" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18M9 21V9"/></svg>
        <svg v-else viewBox="0 0 24 24" aria-hidden="true"><path d="M8 6h13M8 12h13M8 18h13"/><path d="M3 6h.01M3 12h.01M3 18h.01"/></svg>
        <span>{{ option.label }}</span>
        <span v-if="option.count !== undefined" class="primary-tab-count">{{ option.count }}</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ComponentPublicInstance } from 'vue'

export interface AppPrimaryTabOption {
  value: string
  label: string
  icon?: 'list' | 'calendar' | 'heart' | 'clock' | 'contract' | 'history' | 'info'
  count?: string | number
  disabled?: boolean
}

const props = withDefaults(defineProps<{ modelValue: string; options: AppPrimaryTabOption[]; ariaLabel?: string }>(), {
  ariaLabel: 'แท็บหลัก'
})
const emit = defineEmits<{ 'update:modelValue': [value: string]; change: [value: string] }>()
const tabRefs = new Map<string, HTMLButtonElement>()

const setTabRef = (value: string, element: Element | ComponentPublicInstance | null) => {
  if (element instanceof HTMLButtonElement) tabRefs.set(value, element)
}
const selectTab = (option: AppPrimaryTabOption) => {
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
.primary-tabs-scroll{position:relative;z-index:2;max-width:100%;margin:0 0 -1px;overflow-x:auto;overflow-y:hidden;scrollbar-width:thin;-webkit-overflow-scrolling:touch}.primary-tabs{display:flex;align-items:flex-end;gap:10px;min-width:max-content;padding:7px 18px 0 0}.primary-tab{--tab-bg:var(--bg-f1f5f9);position:relative;z-index:1;display:inline-flex;align-items:center;justify-content:center;gap:8px;width:auto;height:50px;box-sizing:border-box;padding:0 22px;border:1px solid var(--bd-dbe5f1);border-bottom:0;border-radius:14px 14px 0 0;background:var(--tab-bg);color:var(--tx-64748b);font-family:inherit;font-size:var(--font-size-base);font-weight:600;line-height:1;white-space:nowrap;cursor:pointer;transition:color .18s ease,background-color .18s ease}.primary-tab:hover:not(:disabled):not(.active){--tab-bg:var(--bg-f8fafc);color:var(--tx-334155)}.primary-tab.active{--tab-bg:var(--bg-ffffff);z-index:3;color:var(--tx-2563eb);font-weight:700}.primary-tab:focus-visible{z-index:5;outline:3px solid rgba(37,99,235,.2);outline-offset:-3px}.primary-tab:disabled{cursor:not-allowed;opacity:.5}.primary-tab svg{width:18px;height:18px;flex:0 0 18px;fill:none;stroke:currentColor;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}.primary-tab-count{display:inline-flex;align-items:center;justify-content:center;min-width:22px;height:22px;padding:0 7px;border-radius:999px;background:var(--bg-e2e8f0);color:var(--tx-475569);font-size:var(--font-size-xs);font-weight:700}.primary-tab.active .primary-tab-count{background:#2563eb;color:#fff}@media(max-width:640px){.primary-tabs{padding-top:5px}.primary-tab{height:46px;padding-left:16px;padding-right:16px;font-size:var(--font-size-sm)}}
</style>
