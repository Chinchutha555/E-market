<template>
  <label class="app-radio" :class="{ 'is-checked': checked, 'is-disabled': disabled }">
    <input
      type="radio"
      class="app-radio-native"
      :name="name"
      :value="value"
      :checked="checked"
      :disabled="disabled"
      @change="emit('update:modelValue', value)"
    />
    <span class="app-radio-dot" aria-hidden="true"></span>
    <span class="app-radio-label"><slot>{{ label }}</slot></span>
  </label>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue?: string | number | null
  value: string | number
  label?: string
  name?: string
  disabled?: boolean
}>()
const emit = defineEmits<{ 'update:modelValue': [value: string | number] }>()
const checked = computed(() => props.modelValue === props.value)
</script>

<style scoped>
.app-radio { position: relative; display: inline-flex; align-items: center; gap: 10px; min-height: 28px; cursor: pointer; color: var(--tx-334155); font-size: var(--font-size-base); font-weight: 500; user-select: none; }
.app-radio-native { position: absolute; opacity: 0; width: 0; height: 0; }
.app-radio-dot { flex: 0 0 auto; width: 20px; height: 20px; box-sizing: border-box; border: 1.5px solid var(--bd-cbd5e1); border-radius: 50%; background: var(--bg-ffffff); transition: border-color .15s ease, background-color .15s ease, box-shadow .15s ease; }
.app-radio:hover .app-radio-dot { border-color: #2563eb; }
.app-radio.is-checked .app-radio-dot { border-color: #2563eb; background: radial-gradient(circle, #ffffff 0 3.5px, #2563eb 4px); }
.app-radio-native:focus-visible + .app-radio-dot { box-shadow: 0 0 0 3px rgba(37, 99, 235, .3); }
.app-radio.is-checked { color: var(--tx-0f172a); font-weight: 600; }
.app-radio.is-disabled { opacity: .5; cursor: not-allowed; }
</style>
