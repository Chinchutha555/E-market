<template>
  <component
    :is="componentTag"
    :to="to"
    :href="href"
    :type="isButton ? type : undefined"
    :disabled="isButton ? (disabled || loading) : undefined"
    :aria-disabled="!isButton && (disabled || loading) ? 'true' : undefined"
    class="app-btn"
    :class="[
      `app-btn-${variant}`,
      `app-btn-${size}`,
      {
        'is-loading': loading,
        'is-disabled': disabled,
        'is-full-width': fullWidth
      }
    ]"
    @click="handleClick"
  >
    <!-- Loading Spinner -->
    <svg
      v-if="loading"
      class="app-btn-spinner"
      width="15"
      height="15"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2.5"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
    </svg>

    <!-- Prefix Icon / Slot -->
    <span v-if="$slots.prefix" class="app-btn-icon app-btn-prefix">
      <slot name="prefix" />
    </span>

    <!-- Main Content -->
    <span class="app-btn-text">
      <slot>{{ label }}</slot>
    </span>

    <!-- Suffix Icon / Slot -->
    <span v-if="$slots.suffix" class="app-btn-icon app-btn-suffix">
      <slot name="suffix" />
    </span>
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export type ButtonVariant =
  | 'primary'
  | 'secondary'
  | 'success'
  | 'danger'
  | 'danger-outline'
  | 'purple'
  | 'outline'
  | 'ghost'

export type ButtonSize = 'sm' | 'base' | 'md' | 'lg'

const props = withDefaults(
  defineProps<{
    label?: string
    variant?: ButtonVariant
    size?: ButtonSize
    type?: 'button' | 'submit' | 'reset'
    to?: string
    href?: string
    disabled?: boolean
    loading?: boolean
    fullWidth?: boolean
  }>(),
  {
    variant: 'primary',
    size: 'md',
    type: 'button',
    disabled: false,
    loading: false,
    fullWidth: false
  }
)

const emit = defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()

const isNuxtLink = computed(() => !!props.to)
const isAnchor = computed(() => !props.to && !!props.href)
const isButton = computed(() => !props.to && !props.href)

const componentTag = computed(() => {
  if (isNuxtLink.value) return resolveComponent('NuxtLink')
  if (isAnchor.value) return 'a'
  return 'button'
})

const handleClick = (e: MouseEvent) => {
  if (props.disabled || props.loading) {
    e.preventDefault()
    e.stopPropagation()
    return
  }
  emit('click', e)
}
</script>

<style scoped>
.app-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  box-sizing: border-box;
  font-family: inherit;
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  white-space: nowrap;
  outline: none;
  transition: all 0.15s ease;
  user-select: none;
  vertical-align: middle;
  line-height: 1;
}

.app-btn:active:not(.is-disabled):not(.is-loading) {
  transform: translateY(1px);
}

.app-btn.is-full-width {
  width: 100%;
}

.app-btn.is-disabled {
  opacity: 0.55;
  cursor: not-allowed;
  pointer-events: none;
}

.app-btn.is-loading {
  cursor: wait;
}

/* ============================================================
   SIZES (Unified height & padding across the entire system)
   ============================================================ */
.app-btn-sm {
  height: 32px;
  min-height: 32px;
  max-height: 32px;
  padding: 0 10px;
  font-size: 0.78rem;
  border-radius: 6px;
  gap: 5px;
}

.app-btn-base {
  height: 36px;
  min-height: 36px;
  max-height: 36px;
  padding: 0 14px;
  font-size: 0.82rem;
  border-radius: 8px;
  gap: 6px;
}

.app-btn-md {
  height: 40px;
  min-height: 40px;
  max-height: 40px;
  padding: 0 16px;
  font-size: 0.875rem;
  border-radius: 8px;
  gap: 7px;
}

.app-btn-lg {
  height: 46px;
  min-height: 46px;
  max-height: 46px;
  padding: 0 22px;
  font-size: 0.95rem;
  border-radius: 10px;
  gap: 8px;
}

/* ============================================================
   VARIANTS
   ============================================================ */
/* 1. Primary (Solid Blue) */
.app-btn-primary {
  background: #2563eb;
  border: 1px solid #2563eb;
  color: #ffffff;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.25);
}
.app-btn-primary:hover:not(.is-disabled) {
  background: #1d4ed8;
  border-color: #1d4ed8;
  color: #ffffff;
}

/* 2. Secondary (Neutral Outlined/Ghost for Close & Cancel) */
.app-btn-secondary {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-cbd5e1);
  color: var(--tx-334155);
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
}
.app-btn-secondary:hover:not(.is-disabled) {
  background: var(--bg-f8fafc);
  border-color: var(--bd-94a3b8);
  color: var(--tx-0f172a);
}

/* 3. Success (Solid Green) */
.app-btn-success {
  background: #10b981;
  border: 1px solid #10b981;
  color: #ffffff;
  box-shadow: 0 2px 6px rgba(16, 185, 129, 0.25);
}
.app-btn-success:hover:not(.is-disabled) {
  background: #059669;
  border-color: #059669;
  color: #ffffff;
}

/* 4. Danger (Solid Red) */
.app-btn-danger {
  background: #dc2626;
  border: 1px solid #dc2626;
  color: #ffffff;
  box-shadow: 0 2px 6px rgba(220, 38, 38, 0.25);
}
.app-btn-danger:hover:not(.is-disabled) {
  background: #b91c1c;
  border-color: #b91c1c;
  color: #ffffff;
}

/* 5. Danger Outline */
.app-btn-danger-outline {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-fca5a5);
  color: var(--tx-ef4444);
}
.app-btn-danger-outline:hover:not(.is-disabled) {
  background: var(--bg-fef2f2);
  border-color: var(--bd-f87171);
  color: var(--tx-dc2626);
}

/* 6. Purple (Solid Purple) */
.app-btn-purple {
  background: #8b5cf6;
  border: 1px solid #8b5cf6;
  color: #ffffff;
  box-shadow: 0 2px 6px rgba(139, 92, 246, 0.25);
}
.app-btn-purple:hover:not(.is-disabled) {
  background: #7c3aed;
  border-color: #7c3aed;
  color: #ffffff;
}

/* 7. Outline */
.app-btn-outline {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  color: var(--tx-334155);
}
.app-btn-outline:hover:not(.is-disabled) {
  background: var(--bg-f8fafc);
  border-color: var(--bd-cbd5e1);
  color: var(--tx-0f172a);
}

/* 8. Ghost */
.app-btn-ghost {
  background: transparent;
  border: 1px solid transparent;
  color: var(--tx-475569);
}
.app-btn-ghost:hover:not(.is-disabled) {
  background: var(--bg-f1f5f9);
  color: var(--tx-0f172a);
}

/* ============================================================
   ELEMENTS (Icons, text, spinner)
   ============================================================ */
.app-btn-text {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.app-btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.app-btn-spinner {
  animation: appBtnSpin 0.75s linear infinite;
  flex-shrink: 0;
}

@keyframes appBtnSpin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
