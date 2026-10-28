<template>
  <div class="app-segmented-control" role="group" :aria-label="ariaLabel">
    <button
      v-for="option in options"
      :key="option.value"
      type="button"
      class="app-segmented-option"
      :class="{ 'is-active': option.value === modelValue }"
      :aria-pressed="option.value === modelValue"
      :title="option.title || option.label"
      @click="selectOption(option.value)"
    >
      <svg v-if="option.icon === 'card'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <rect x="3" y="4" width="18" height="16" rx="2.5"></rect>
        <line x1="3" y1="9.5" x2="21" y2="9.5"></line>
        <line x1="7" y1="14.5" x2="13" y2="14.5"></line>
      </svg>
      <svg v-else-if="option.icon === 'table'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <rect x="3" y="3" width="18" height="18" rx="2"></rect>
        <path d="M3 9h18"></path>
        <path d="M3 15h18"></path>
        <path d="M12 3v18"></path>
      </svg>
      <svg v-else-if="option.icon === 'grid'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <rect x="3" y="3" width="7" height="7" rx="1.5"></rect>
        <rect x="14" y="3" width="7" height="7" rx="1.5"></rect>
        <rect x="14" y="14" width="7" height="7" rx="1.5"></rect>
        <rect x="3" y="14" width="7" height="7" rx="1.5"></rect>
      </svg>
      <svg v-else-if="option.icon === 'list'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <line x1="8" y1="6" x2="21" y2="6"></line>
        <line x1="8" y1="12" x2="21" y2="12"></line>
        <line x1="8" y1="18" x2="21" y2="18"></line>
        <line x1="3" y1="6" x2="3.01" y2="6"></line>
        <line x1="3" y1="12" x2="3.01" y2="12"></line>
        <line x1="3" y1="18" x2="3.01" y2="18"></line>
      </svg>
      <svg v-else-if="option.icon === 'market'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
        <polyline points="9 22 9 12 15 12 15 22"></polyline>
      </svg>
      <svg v-else-if="option.icon === 'map'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"></polygon>
        <line x1="8" y1="2" x2="8" y2="18"></line>
        <line x1="16" y1="6" x2="16" y2="22"></line>
      </svg>
      <svg v-else-if="option.icon === 'receipt'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <path d="M4 2v20l2-1 2 1 2-1 2 1 2-1 2 1 2-1 2 1V2l-2 1-2-1-2 1-2-1-2 1-2-1-2 1z"></path>
        <path d="M8 8h8M8 12h8M8 16h5"></path>
      </svg>
      <span>{{ option.label }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
type SegmentedIcon = 'market' | 'grid' | 'list' | 'card' | 'table' | 'map' | 'receipt'

interface SegmentedOption {
  value: string
  label: string
  title?: string
  icon?: SegmentedIcon
}

defineProps<{
  modelValue: string
  options: SegmentedOption[]
  ariaLabel?: string
}>()

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void
  (event: 'change', value: string): void
}>()

const selectOption = (value: string) => {
  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style scoped>
.app-segmented-control {
  display: inline-flex;
  align-items: center;
  padding: 3px;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
}

.app-segmented-option {
  min-height: 32px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 5px 12px;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: #64748b;
  font-family: inherit;
  font-size: 0.8rem;
  font-weight: 600;
  white-space: nowrap;
  cursor: pointer;
  transition: color 0.18s ease, background 0.18s ease, box-shadow 0.18s ease;
}

.app-segmented-option:hover:not(.is-active) {
  color: #334155;
}

.app-segmented-option.is-active {
  background: #ffffff;
  color: #0f172a;
  font-weight: 700;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.1);
}

.app-segmented-option:focus-visible {
  outline: 2px solid #93c5fd;
  outline-offset: 1px;
}

.app-segmented-option svg {
  width: 15px;
  height: 15px;
  flex-shrink: 0;
}
</style>
