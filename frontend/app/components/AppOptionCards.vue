<template>
  <div class="app-option-cards" role="radiogroup" :aria-label="ariaLabel">
    <button
      v-for="option in options"
      :key="option.value"
      type="button"
      class="app-option-card"
      :class="{ active: option.value === modelValue }"
      role="radio"
      :aria-checked="option.value === modelValue"
      @click="select(option.value)"
    >
      <span class="app-option-dot" aria-hidden="true"></span>
      <span class="app-option-text">
        <strong>{{ option.label }}</strong>
        <small v-if="option.description">{{ option.description }}</small>
      </span>
    </button>
  </div>
</template>

<script setup lang="ts">
export interface OptionCard { value: string; label: string; description?: string }

const props = withDefaults(defineProps<{ modelValue: string; options: OptionCard[]; ariaLabel?: string }>(), { ariaLabel: 'ตัวเลือก' })
const emit = defineEmits<{ 'update:modelValue': [value: string]; change: [value: string] }>()
const select = (value: string) => {
  if (value === props.modelValue) return
  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style scoped>
.app-option-cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 12px; }
.app-option-card { display: flex; align-items: flex-start; gap: 12px; padding: 14px 16px; border: 1.5px solid var(--bd-e2e8f0); border-radius: 12px; background: var(--bg-ffffff); font: inherit; text-align: left; cursor: pointer; transition: border-color .15s ease, background-color .15s ease, box-shadow .15s ease; }
.app-option-card:hover { border-color: var(--bd-cbd5e1); background: var(--bg-f8fafc); }
.app-option-card:focus-visible { outline: 2px solid #93c5fd; outline-offset: 2px; }
.app-option-card.active { border-color: #3b82f6; background: var(--bg-eff6ff); box-shadow: 0 2px 8px rgba(37, 99, 235, 0.1); }
.app-option-dot { flex: 0 0 auto; width: 20px; height: 20px; margin-top: 1px; box-sizing: border-box; border: 1.5px solid var(--bd-cbd5e1); border-radius: 50%; background: var(--bg-ffffff); }
.app-option-card.active .app-option-dot { border-color: #2563eb; background: radial-gradient(circle, #ffffff 0 3.5px, #2563eb 4px); }
.app-option-text { display: flex; flex-direction: column; gap: 3px; min-width: 0; }
.app-option-text strong { color: var(--tx-1e293b); font-size: var(--font-size-base); font-weight: 700; }
.app-option-card.active .app-option-text strong { color: var(--tx-1d4ed8); }
.app-option-text small { color: var(--tx-64748b); font-size: var(--font-size-xs); font-weight: 500; line-height: 1.45; }
</style>
