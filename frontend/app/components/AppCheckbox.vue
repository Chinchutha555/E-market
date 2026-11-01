<template>
  <label
    class="app-checkbox-wrapper"
    :class="[
      `size-${size}`,
      {
        'is-checked': isChecked,
        'is-indeterminate': indeterminate,
        'is-disabled': disabled,
        'has-label': !!label || $slots.default
      }
    ]"
    :for="id"
  >
    <span class="app-checkbox-control">
      <input
        :id="id"
        ref="inputRef"
        type="checkbox"
        class="app-checkbox-native"
        :checked="isChecked"
        :disabled="disabled"
        :name="name"
        :value="value"
        :aria-checked="indeterminate ? 'mixed' : isChecked"
        :aria-label="ariaLabel || label"
        @change="handleChange"
        @focus="isFocused = true"
        @blur="isFocused = false"
      />
      <span
        class="app-checkbox-box"
        :class="{ 'is-focused': isFocused }"
        aria-hidden="true"
      >
        <!-- Indeterminate Icon (Dash) -->
        <svg
          v-if="indeterminate"
          class="app-checkbox-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="3"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>

        <!-- Checked Icon (Checkmark) -->
        <svg
          v-else-if="isChecked"
          class="app-checkbox-icon check-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2.8"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <polyline points="20 6 9 17 4 12"></polyline>
        </svg>
      </span>
    </span>

    <span v-if="label || $slots.default" class="app-checkbox-label">
      <slot>{{ label }}</slot>
    </span>
  </label>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: boolean | any[] | null
    value?: any
    label?: string
    id?: string
    name?: string
    ariaLabel?: string
    disabled?: boolean
    indeterminate?: boolean
    size?: 'sm' | 'md' | 'lg'
  }>(),
  {
    modelValue: false,
    value: true,
    label: '',
    id: undefined,
    name: undefined,
    ariaLabel: undefined,
    disabled: false,
    indeterminate: false,
    size: 'md'
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean | any[]): void
  (e: 'change', value: boolean | any[], event: Event): void
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const isFocused = ref(false)

const isChecked = computed<boolean>(() => {
  if (Array.isArray(props.modelValue)) {
    return props.modelValue.includes(props.value)
  }
  return !!props.modelValue
})

const handleChange = (e: Event) => {
  if (props.disabled) return

  const target = e.target as HTMLInputElement
  let newValue: boolean | any[]

  if (Array.isArray(props.modelValue)) {
    const list = [...props.modelValue]
    const idx = list.indexOf(props.value)
    if (target.checked && idx === -1) {
      list.push(props.value)
    } else if (!target.checked && idx !== -1) {
      list.splice(idx, 1)
    }
    newValue = list
  } else {
    newValue = target.checked
  }

  emit('update:modelValue', newValue)
  emit('change', newValue, e)
}
</script>

<style scoped>
.app-checkbox-wrapper {
  display: inline-flex;
  align-items: center;
  position: relative;
  cursor: pointer;
  user-select: none;
  vertical-align: middle;
  line-height: 1.4;
  gap: 8px;
}

.app-checkbox-wrapper.is-disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.app-checkbox-control {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* Hidden native input for full keyboard/a11y support */
.app-checkbox-native {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  margin: 0;
  padding: 0;
  z-index: 1;
  cursor: inherit;
}

/* Custom Checkbox Box */
.app-checkbox-box {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  background: var(--bg-ffffff);
  border: 1.5px solid var(--bd-cbd5e1);
  color: #ffffff;
  transition: all 0.18s cubic-bezier(0.4, 0, 0.2, 1);
  box-sizing: border-box;
}

/* Hover state */
.app-checkbox-wrapper:hover:not(.is-disabled) .app-checkbox-box {
  border-color: var(--bd-93c5fd);
  background: var(--bg-f8fafc);
}

/* Checked or Indeterminate state */
.app-checkbox-wrapper.is-checked .app-checkbox-box,
.app-checkbox-wrapper.is-indeterminate .app-checkbox-box {
  background: #2563eb;
  border-color: #2563eb;
  color: #ffffff;
  box-shadow: 0 1px 3px rgba(37, 99, 235, 0.25);
}

.app-checkbox-wrapper.is-checked:hover:not(.is-disabled) .app-checkbox-box,
.app-checkbox-wrapper.is-indeterminate:hover:not(.is-disabled) .app-checkbox-box {
  background: #1d4ed8;
  border-color: #1d4ed8;
}

/* Focus ring */
.app-checkbox-native:focus-visible ~ .app-checkbox-box,
.app-checkbox-box.is-focused {
  outline: none;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.2);
}

/* Sizes */
.size-sm .app-checkbox-box {
  width: 16px;
  height: 16px;
  border-radius: 5px;
}
.size-sm .app-checkbox-icon {
  width: 11px;
  height: 11px;
}

.size-md .app-checkbox-box {
  width: 19px;
  height: 19px;
  border-radius: 6px;
}
.size-md .app-checkbox-icon {
  width: 13px;
  height: 13px;
}

.size-lg .app-checkbox-box {
  width: 22px;
  height: 22px;
  border-radius: 7px;
}
.size-lg .app-checkbox-icon {
  width: 15px;
  height: 15px;
}

/* Checkmark Icon Animation */
.app-checkbox-icon {
  display: block;
  flex-shrink: 0;
  stroke: #ffffff;
  stroke-dasharray: 24;
  stroke-dashoffset: 0;
  animation: checkmarkIn 0.18s ease-in-out;
}

@keyframes checkmarkIn {
  0% {
    opacity: 0;
    transform: scale(0.65);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

/* Label text */
.app-checkbox-label {
  font-size: 0.88rem;
  color: var(--tx-334155);
  font-weight: 500;
  cursor: inherit;
}
</style>
