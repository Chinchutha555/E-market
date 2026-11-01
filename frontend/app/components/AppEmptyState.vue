<template>
  <div class="app-empty-state" :class="`app-empty-state--${variant}`" role="status">
    <div class="app-empty-state__icon" aria-hidden="true">
      <svg v-if="icon === 'search'" viewBox="0 0 24 24"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
      <svg v-else-if="icon === 'document'" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
      <svg v-else-if="icon === 'calendar'" viewBox="0 0 24 24"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
      <svg v-else-if="icon === 'history'" viewBox="0 0 24 24"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 1 0 .49-8.5L1 10"/><polyline points="12 7 12 12 15 14"/></svg>
      <svg v-else-if="icon === 'bell'" viewBox="0 0 24 24"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>
      <svg v-else-if="icon === 'heart'" viewBox="0 0 24 24"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78L12 21.23l8.84-8.84a5.5 5.5 0 0 0 0-7.78z"/></svg>
      <svg v-else-if="icon === 'users'" viewBox="0 0 24 24"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
      <svg v-else-if="icon === 'alert'" viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
      <svg v-else viewBox="0 0 24 24"><path d="M21 8v13H3V8"/><path d="M1 3h22v5H1z"/><path d="M10 12h4"/></svg>
    </div>

    <h3 class="app-empty-state__title">{{ title }}</h3>
    <p v-if="description" class="app-empty-state__description">{{ description }}</p>
    <div v-if="$slots.details" class="app-empty-state__details"><slot name="details" /></div>

    <div v-if="$slots.actions || actionLabel" class="app-empty-state__actions">
      <slot name="actions">
        <AppButton :variant="actionVariant" :size="variant === 'compact' || variant === 'table' ? 'sm' : 'md'" @click="$emit('action')">
          {{ actionLabel }}
        </AppButton>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppButton from '~/components/AppButton.vue'

withDefaults(defineProps<{
  title: string
  description?: string
  icon?: 'inbox' | 'search' | 'document' | 'calendar' | 'history' | 'bell' | 'heart' | 'users' | 'alert'
  variant?: 'page' | 'card' | 'table' | 'compact' | 'map'
  actionLabel?: string
  actionVariant?: 'primary' | 'secondary' | 'success' | 'danger' | 'danger-outline' | 'purple' | 'outline' | 'ghost'
}>(), {
  description: '',
  icon: 'inbox',
  variant: 'card',
  actionLabel: '',
  actionVariant: 'outline'
})

defineEmits<{ action: [] }>()
</script>

<style scoped>
.app-empty-state {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 36px 24px;
  color: var(--text-muted);
  text-align: center;
}

.app-empty-state--page { min-height: 320px; padding: 56px 24px; }
.app-empty-state--card { min-height: 240px; }
.app-empty-state--table { min-height: 180px; padding: 28px 18px; }
/* ขนาดเล็ก: ใช้ในแผงย่อยและการ์ดแคบ (แท็บรายละเอียดแผง แดชบอร์ด) ไอคอน/ข้อความเล็กลงตามสัดส่วน */
.app-empty-state--compact { min-height: 0; padding: 20px 16px; }
.app-empty-state--compact .app-empty-state__icon { width: 44px; height: 44px; margin-bottom: 12px; border-radius: 13px; }
.app-empty-state--compact .app-empty-state__icon svg { width: 22px; height: 22px; }
.app-empty-state--compact .app-empty-state__title { font-size: var(--font-size-md); }
.app-empty-state--compact .app-empty-state__description { max-width: 320px; margin-top: 6px; font-size: var(--font-size-sm); line-height: 1.5; }
.app-empty-state--compact .app-empty-state__actions { margin-top: 12px; }
.app-empty-state--map {
  min-height: 180px;
  padding: 32px 18px;
}

.app-empty-state__icon {
  width: 64px;
  height: 64px;
  display: grid;
  place-items: center;
  margin-bottom: 20px;
  border: 1px solid var(--bd-dbeafe);
  border-radius: 18px;
  background: var(--primary-light);
  color: var(--primary);
}

.app-empty-state__icon svg {
  width: 32px;
  height: 32px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.app-empty-state__title {
  margin: 0;
  color: var(--text-main);
  font-size: var(--font-size-xl);
  font-weight: 700;
  line-height: 1.4;
}

.app-empty-state__description {
  max-width: 480px;
  margin: 10px 0 0;
  color: var(--text-muted);
  font-size: var(--font-size-base);
  font-weight: 500;
  line-height: 1.6;
}

.app-empty-state__actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 18px;
}

.app-empty-state__details {
  width: 100%;
  max-width: 480px;
  margin-top: 16px;
}

@media (max-width: 640px) {
  .app-empty-state,
  .app-empty-state--page,
  .app-empty-state--card {
    min-height: 210px;
    padding: 30px 16px;
  }

  .app-empty-state__actions {
    width: 100%;
    flex-wrap: wrap;
  }
}
</style>
