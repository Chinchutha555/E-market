<template>
  <section v-if="isStaff && market" class="staff-market-bar" aria-label="พื้นที่ทำงานเจ้าหน้าที่ตลาด">
    <div class="staff-market-context">
      <span>ตลาดที่รับผิดชอบ</span>
      <strong>{{ market.name }}</strong>
      <small>{{ market.code }} · {{ market.district }}</small>
    </div>
    <nav class="staff-market-tabs" aria-label="เมนูตลาด">
      <NuxtLink :to="dashboardLink" :class="{ active: active === 'dashboard' }">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
        แดชบอร์ด
      </NuxtLink>
      <NuxtLink :to="infoLink" :class="{ active: active === 'info' }">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 21h18"/><path d="M6 21V5l6-3 6 3v16"/><path d="M9 9h.01"/><path d="M15 9h.01"/><path d="M9 13h.01"/><path d="M15 13h.01"/></svg>
        รายละเอียดตลาด
      </NuxtLink>
      <NuxtLink :to="layoutLink" class="manage-layout-link">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4Z"/></svg>
        จัดการผัง GIS
      </NuxtLink>
    </nav>
  </section>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'

defineProps<{ active: 'dashboard' | 'info' }>()
const auth = useAuthStore()
const isStaff = computed(() => auth.role === 'staff')
const market = computed(() => auth.currentUser?.assigned_markets?.[0])
const dashboardLink = computed(() => `/stalls/manage?market=${market.value?.id || ''}&mode=inspect`)
const infoLink = computed(() => `/stalls?market=${market.value?.id || ''}`)
const layoutLink = computed(() => `/stalls/manage?market=${market.value?.id || ''}`)
</script>

<style scoped>
.staff-market-bar { margin: 0 0 18px; padding: 14px 16px; border: 1px solid var(--bd-dbe5f1); border-radius: 14px; background: var(--bg-ffffff); display: flex; align-items: center; justify-content: space-between; gap: 20px; box-shadow: 0 1px 3px rgba(15,23,42,.03); }
.staff-market-context { min-width: 0; display: grid; grid-template-columns: auto auto; align-items: baseline; gap: 2px 10px; }
.staff-market-context span { color: var(--tx-64748b); font-size: .8125rem; font-weight: 600; }
.staff-market-context strong { color: var(--tx-1e293b); font-size: 1rem; }
.staff-market-context small { grid-column: 1 / -1; color: var(--tx-64748b); font-size: .8125rem; }
.staff-market-tabs { display: inline-flex; align-items: center; padding: 4px; gap: 4px; background: var(--bg-f1f5f9); border-radius: 10px; }
.staff-market-tabs a { height: 38px; padding: 0 14px; border-radius: 8px; display: inline-flex; align-items: center; gap: 7px; color: var(--tx-64748b); font-size: .875rem; font-weight: 600; text-decoration: none; white-space: nowrap; }
.staff-market-tabs a.active { background: var(--bg-ffffff); color: var(--tx-2563eb); box-shadow: 0 1px 4px rgba(15,23,42,.1); }
.staff-market-tabs a.manage-layout-link { color: var(--tx-1d4ed8); border-left: 1px solid var(--bd-dbe5f1); }
.staff-market-tabs svg { width: 17px; height: 17px; }
@media (max-width: 820px) { .staff-market-bar { align-items: stretch; flex-direction: column; } .staff-market-tabs { width: 100%; overflow-x: auto; } .staff-market-tabs a { flex: 1; justify-content: center; } }
</style>
