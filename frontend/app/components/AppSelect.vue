<template>
  <div ref="rootRef" class="app-select-container" :class="{ 'is-open': isOpen, 'is-disabled': disabled }">
    <!-- Trigger Button (ซ่อนในโหมด inline: รายการแสดงอยู่ในฟอร์มเลย) -->
    <button
      v-if="!inline"
      type="button"
      class="app-select-trigger"
      :class="{ 'has-value': !!selectedLabel, 'is-focused': isOpen, 'is-error': Boolean(error) }"
      :disabled="disabled"
      @click="toggleDropdown"
      :aria-expanded="isOpen"
      :aria-invalid="Boolean(error)"
    >
      <span class="app-select-label">
        {{ selectedLabel || placeholder }}
      </span>
      <span class="app-select-arrow" :class="{ 'is-rotated': isOpen }">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>
      </span>
    </button>
    <!-- Dropdown Menu -->
    <transition name="app-select-fade">
      <div v-if="inline || isOpen" class="app-select-menu" :class="{ 'is-inline': inline }">
        <div v-if="searchable" class="app-select-search">
          <input v-model.trim="query" type="search" placeholder="ค้นหา..." aria-label="ค้นหาในรายการ" @click.stop />
        </div>
        <div class="app-select-list" role="listbox" :aria-multiselectable="multiple || undefined">
          <div
            v-for="opt in visibleOptions"
            :key="opt.value"
            class="app-select-option"
            :class="{ 'is-selected': isOptionSelected(opt) }"
            role="option"
            :aria-selected="isOptionSelected(opt)"
            @click.stop="selectOption(opt)"
          >
            <span class="app-select-option-text">{{ opt.label }}</span>
            <span v-if="isOptionSelected(opt)" class="app-select-check">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
            </span>
          </div>
        </div>
      </div>
    </transition>
    <span v-if="typeof error === 'string' && error" class="app-select-error">{{ error }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'

export interface SelectOption {
  value: string | number
  label: string
}

const props = withDefaults(defineProps<{
  modelValue?: string | number | (string | number)[]
  options: (SelectOption | string)[]
  placeholder?: string
  disabled?: boolean
  error?: boolean | string
  // multiple = เลือกได้หลายค่า (modelValue เป็นอาร์เรย์) เมนูไม่ปิดหลังเลือก
  multiple?: boolean
  // inline = แสดงรายการในฟอร์มเลย ไม่ใช่ป๊อปอัป (เหมาะกับ multiple ที่ต้องเห็นตัวเลือกพร้อมกัน)
  inline?: boolean
  // searchable = มีช่องค้นหาในเมนู (เหมาะกับรายการยาว)
  searchable?: boolean
}>(), {
  modelValue: '',
  options: () => [],
  placeholder: 'กรุณาเลือก...',
  disabled: false,
  error: false,
  multiple: false,
  inline: false,
  searchable: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', val: string | number | (string | number)[]): void
  (e: 'change', val: string | number | (string | number)[]): void
}>()

const rootRef = ref<HTMLElement | null>(null)
const isOpen = ref(false)

const normalizedOptions = computed<SelectOption[]>(() => {
  return props.options.map(opt => {
    if (typeof opt === 'string') {
      return { value: opt, label: opt }
    }
    return opt
  })
})

const query = ref('')
const visibleOptions = computed(() => {
  const q = query.value.toLowerCase()
  if (!props.searchable || !q) return normalizedOptions.value
  return normalizedOptions.value.filter(o => o.label.toLowerCase().includes(q))
})

// ค่าที่เลือกในรูปอาร์เรย์เสมอ (single = ค่าเดียว, multiple = หลายค่า)
const selectedValues = computed<string[]>(() => {
  const v = props.modelValue
  if (props.multiple) return Array.isArray(v) ? v.map(String) : []
  return v === '' || v == null ? [] : [String(v)]
})

const isOptionSelected = (opt: SelectOption) => selectedValues.value.includes(String(opt.value))

const selectedOption = computed(() => {
  if (props.multiple) return null
  return normalizedOptions.value.find(o => String(o.value) === String(props.modelValue ?? '')) || null
})

// multiple: แสดงชื่อ 2 รายการแรก ที่เหลือเป็นจำนวน เช่น "แผง 1, แผง 2 และอีก 3"
const selectedLabel = computed(() => {
  if (!props.multiple) return selectedOption.value?.label || ''
  const labels = normalizedOptions.value.filter(o => selectedValues.value.includes(String(o.value))).map(o => o.label)
  if (labels.length <= 2) return labels.join(', ')
  return `${labels[0]}, ${labels[1]} และอีก ${labels.length - 2}`
})

const toggleDropdown = () => {
  if (props.disabled) return
  isOpen.value = !isOpen.value
}

const selectOption = (opt: SelectOption) => {
  if (props.multiple) {
    const key = String(opt.value)
    const next = selectedValues.value.includes(key)
      ? selectedValues.value.filter(v => v !== key)
      : [...selectedValues.value, key]
    const original = normalizedOptions.value.filter(o => next.includes(String(o.value))).map(o => o.value)
    emit('update:modelValue', original)
    emit('change', original)
    return
  }
  emit('update:modelValue', opt.value)
  emit('change', opt.value)
  isOpen.value = false
  // ป้องกัน event จากคอนเทนเนอร์ภายนอกเปิดเมนูซ้ำในรอบ click เดียวกัน
  requestAnimationFrame(() => { isOpen.value = false })
}

const handleClickOutside = (e: MouseEvent) => {
  if (rootRef.value && !rootRef.value.contains(e.target as Node)) {
    isOpen.value = false
  }
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isOpen.value) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.app-select-container {
  position: relative;
  width: 100%;
  font-family: inherit;
}

.app-select-container.is-open {
  z-index: 50;
}

.app-select-trigger {
  width: 100%;
  min-height: 42px;
  padding: 9px 14px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 10px;
  font-family: inherit;
  font-size: 0.88rem;
  color: var(--tx-94a3b8);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  cursor: pointer;
  transition: all 0.18s ease;
  outline: none;
  text-align: left;
}

.app-select-trigger:hover:not(:disabled) {
  border-color: var(--bd-94a3b8);
  background: var(--bg-f8fafc);
}

.app-select-trigger.is-focused {
  border-color: #2563eb;
  background: var(--bg-ffffff);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
}

.app-select-container.is-disabled .app-select-trigger,
.app-select-trigger:disabled {
  background: var(--bg-f8fafc) !important;
  color: var(--tx-64748b) !important;
  border-color: var(--bd-e2e8f0) !important;
  cursor: not-allowed !important;
}

.app-select-trigger.has-value {
  color: var(--tx-1e293b);
  font-weight: 500;
}

.app-select-trigger.is-error {
  border-color: #ef4444 !important;
  box-shadow: 0 0 0 1px #ef4444, 0 1px 2px rgba(239, 68, 68, 0.1) !important;
  background: var(--bg-fffafa);
}

.app-select-error {
  display: block;
  margin-top: 5px;
  color: var(--tx-dc2626);
  font-size: 0.8125rem;
  font-weight: 500;
  line-height: 1.3;
}

.app-select-trigger:disabled {
  background: var(--bg-f8fafc);
  color: var(--tx-94a3b8);
  cursor: not-allowed;
  opacity: 0.75;
}

.app-select-label {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.app-select-arrow {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--tx-64748b);
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  flex-shrink: 0;
}

.app-select-arrow.is-rotated {
  transform: rotate(180deg);
  color: var(--tx-2563eb);
}

.app-select-menu {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  z-index: 100;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.12), 0 4px 10px rgba(15, 23, 42, 0.04);
  padding: 6px;
  max-height: 260px;
  overflow-y: auto;
  overscroll-behavior: contain;
}

.app-select-menu.is-inline { position: static; top: auto; box-shadow: none; max-height: none; overflow: visible; margin-top: 0; }
.app-select-menu.is-inline .app-select-list { display: grid; grid-template-columns: repeat(auto-fill, minmax(130px, 1fr)); gap: 4px; max-height: 260px; overflow-y: auto; overscroll-behavior: contain; padding: 2px; }
.app-select-search { padding: 0 0 6px; background: var(--bg-ffffff); }
.app-select-menu:not(.is-inline) .app-select-search { position: sticky; top: -6px; z-index: 1; }
.app-select-search input { width: 100%; box-sizing: border-box; height: 36px; padding: 0 12px; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-0f172a); font: inherit; font-size: 0.875rem; outline: none; }
.app-select-search input:focus { border-color: #2563eb; box-shadow: 0 0 0 3px rgba(37, 99, 235, .15); }

.app-select-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.app-select-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 9px 12px;
  border-radius: 8px;
  font-size: 0.86rem;
  color: var(--tx-334155);
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
}

.app-select-option:hover {
  background: var(--bg-f1f5f9);
  color: var(--tx-0f172a);
}

.app-select-option.is-selected {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  font-weight: 600;
}

.app-select-check {
  display: inline-flex;
  align-items: center;
  color: var(--tx-2563eb);
  flex-shrink: 0;
  margin-left: 8px;
}

.app-select-fade-enter-active,
.app-select-fade-leave-active {
  transition: opacity 0.18s ease, transform 0.18s cubic-bezier(0.16, 1, 0.3, 1);
}

.app-select-fade-enter-from,
.app-select-fade-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
