<template>
  <div class="app-pagination-container">
    <!-- Left: Items Per Page Selector with AppSelect -->
    <div class="pagination-left">
      <span class="pagination-label">แสดง</span>
      <div class="pagination-select-wrap">
        <AppSelect
          :model-value="effectiveItemsPerPage"
          :options="pageSizeSelectOptions"
          :disabled="disabled"
          @update:model-value="onPageSizeSelect"
        />
      </div>
      <span class="pagination-label">{{ unitLabel || 'รายการ' }}ต่อหน้า</span>
      <span v-if="totalItems !== undefined" class="pagination-total-items">
        (ทั้งหมด {{ totalItems }} {{ unitLabel || 'รายการ' }})
      </span>
    </div>

    <!-- Right: Page Navigation Buttons -->
    <div class="pagination-right">
      <button
        type="button"
        class="page-nav-btn"
        :disabled="disabled || modelValue <= 1"
        @click="goToPage(1)"
        title="หน้าแรก"
        aria-label="หน้าแรก"
      >
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="11 17 6 12 11 7"></polyline>
          <polyline points="18 17 13 12 18 7"></polyline>
        </svg>
      </button>

      <button
        type="button"
        class="page-nav-btn"
        :disabled="disabled || modelValue <= 1"
        @click="goToPage(modelValue - 1)"
        title="หน้าก่อนหน้า"
        aria-label="หน้าก่อนหน้า"
      >
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="15 18 9 12 15 6"></polyline>
        </svg>
      </button>

      <button
        v-for="(p, idx) in visiblePages"
        :key="idx"
        type="button"
        class="page-num-btn"
        :class="{ active: p === modelValue, ellipsis: p === '...' }"
        :disabled="disabled || p === '...'"
        @click="p !== '...' && goToPage(Number(p))"
      >
        {{ p }}
      </button>

      <button
        type="button"
        class="page-nav-btn"
        :disabled="disabled || modelValue >= effectiveTotalPages"
        @click="goToPage(modelValue + 1)"
        title="หน้าถัดไป"
        aria-label="หน้าถัดไป"
      >
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="9 18 15 12 9 6"></polyline>
        </svg>
      </button>

      <button
        type="button"
        class="page-nav-btn"
        :disabled="disabled || modelValue >= effectiveTotalPages"
        @click="goToPage(effectiveTotalPages)"
        title="หน้าสุดท้าย"
        aria-label="หน้าสุดท้าย"
      >
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="13 17 18 12 13 7"></polyline>
          <polyline points="6 17 11 12 6 7"></polyline>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: number
    totalItems?: number
    totalPages?: number
    itemsPerPage?: number
    pageSizeOptions?: number[]
    disabled?: boolean
    unitLabel?: string
  }>(),
  {
    modelValue: 1,
    itemsPerPage: 10,
    pageSizeOptions: () => [5, 10, 20, 50],
    disabled: false,
    unitLabel: 'รายการ'
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', page: number): void
  (e: 'update:itemsPerPage', size: number): void
  (e: 'change', page: number): void
  (e: 'change-page-size', size: number): void
}>()

const effectiveItemsPerPage = computed(() => Math.max(1, props.itemsPerPage || 10))

const pageSizeSelectOptions = computed(() => {
  const list = props.pageSizeOptions && props.pageSizeOptions.length > 0
    ? props.pageSizeOptions
    : [5, 10, 20, 50]
  return list.map(size => ({
    value: size,
    label: String(size)
  }))
})

const effectiveTotalPages = computed(() => {
  if (typeof props.totalItems === 'number') {
    return Math.max(1, Math.ceil(props.totalItems / effectiveItemsPerPage.value))
  }
  return Math.max(1, props.totalPages || 1)
})

const visiblePages = computed(() => {
  const total = effectiveTotalPages.value
  const current = Math.min(Math.max(1, props.modelValue), total)

  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  const pages: (number | string)[] = []
  if (current <= 3) {
    pages.push(1, 2, 3, 4, '...', total)
  } else if (current >= total - 2) {
    pages.push(1, '...', total - 3, total - 2, total - 1, total)
  } else {
    pages.push(1, '...', current - 1, current, current + 1, '...', total)
  }
  return pages
})

const goToPage = (page: number) => {
  if (props.disabled) return
  const target = Math.min(Math.max(1, page), effectiveTotalPages.value)
  if (target !== props.modelValue) {
    emit('update:modelValue', target)
    emit('change', target)
  }
}

const onPageSizeSelect = (val: string | number) => {
  const newSize = Number(val)
  if (!isNaN(newSize) && newSize > 0) {
    emit('update:itemsPerPage', newSize)
    emit('change-page-size', newSize)
    // Always reset to page 1 when changing page size to prevent out-of-range index
    if (props.modelValue !== 1) {
      emit('update:modelValue', 1)
      emit('change', 1)
    }
  }
}
</script>

<style scoped>
.app-pagination-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 18px;
  margin-top: 20px;
  border-top: 1px solid var(--bd-f1f5f9);
  flex-wrap: wrap;
  gap: 14px;
  font-family: inherit;
}

.pagination-left {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--tx-64748b);
  font-size: 0.875rem;
  font-weight: 500;
  flex-wrap: wrap;
}

.pagination-label {
  white-space: nowrap;
  user-select: none;
}

.pagination-total-items {
  color: var(--tx-94a3b8);
  font-size: 0.8125rem;
  font-weight: 500;
  margin-left: 2px;
  white-space: nowrap;
}

.pagination-select-wrap {
  width: 72px;
}

.pagination-select-wrap :deep(.app-select-trigger) {
  min-height: 34px;
  height: 34px;
  padding: 0 8px 0 10px;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--tx-1e293b);
  border-color: var(--bd-cbd5e1);
  background: var(--bg-ffffff);
}

.pagination-select-wrap :deep(.app-select-trigger:hover) {
  border-color: var(--bd-94a3b8);
}

.pagination-select-wrap :deep(.app-select-trigger.is-focused) {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
}

.pagination-select-wrap :deep(.app-select-arrow) {
  color: var(--tx-64748b);
}

.pagination-select-wrap :deep(.app-select-menu) {
  min-width: 72px;
  max-height: 180px;
  border-radius: 8px;
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.12);
}

.pagination-select-wrap :deep(.app-select-option) {
  padding: 7px 10px;
  font-size: 0.85rem;
  font-weight: 600;
  justify-content: center;
}

.pagination-select-wrap :deep(.app-select-check) {
  display: none;
}

.pagination-right {
  display: flex;
  align-items: center;
  gap: 5px;
}

.page-nav-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 8px;
  background: var(--bg-ffffff);
  color: var(--tx-64748b);
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}

.page-nav-btn:hover:not(:disabled) {
  border-color: var(--bd-93c5fd);
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
}

.page-nav-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  border-color: var(--bd-eef2f6);
  background: var(--bg-f8fafc);
}

.page-num-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 34px;
  height: 34px;
  padding: 0 6px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 8px;
  background: var(--bg-ffffff);
  color: var(--tx-334155);
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}

.page-num-btn:hover:not(:disabled):not(.active) {
  border-color: var(--bd-cbd5e1);
  background: var(--bg-f8fafc);
  color: var(--tx-0f172a);
}

.page-num-btn.active {
  background: #2563eb;
  border-color: #2563eb;
  color: #ffffff;
  font-weight: 700;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.22);
  cursor: default;
}

.page-num-btn.ellipsis {
  border-color: transparent;
  background: transparent;
  cursor: default;
  color: var(--tx-94a3b8);
  padding: 0 2px;
  min-width: 24px;
}

@media (max-width: 640px) {
  .app-pagination-container {
    justify-content: center;
    gap: 12px;
  }
}
</style>
