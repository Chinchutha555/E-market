<template>
  <header class="top-navbar">
    <div class="navbar-left">
      <!-- Hamburger Button (Expand / Collapse Sidebar) -->
      <button
        type="button"
        class="btn-icon-square hamburger-toggle-btn"
        :class="{ active: isExpanded }"
        @click.stop="toggleSidebar"
        :title="isExpanded ? 'ย่อแถบเมนู (Collapse Sidebar)' : 'ขยายแถบเมนู (Expand Sidebar)'"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="3" y1="12" x2="21" y2="12"></line>
          <line x1="3" y1="6" x2="21" y2="6"></line>
          <line x1="3" y1="18" x2="21" y2="18"></line>
        </svg>
      </button>

      <div v-if="subtitle" class="navbar-subtitle-text">
        {{ subtitle }}
      </div>
    </div>

    <div class="navbar-right">
      <!-- Notification Bell -->
      <button class="notification-bell-btn" title="การแจ้งเตือน">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
          <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
        </svg>
        <span class="notification-badge-dot"></span>
      </button>


      <!-- User Profile Badge -->
      <div class="user-profile-badge" @click="toggleMenu" style="position: relative; cursor: pointer;">
        <div class="user-avatar-circle" :style="user?.role === 'tenant' ? 'border: 2px solid #10b981;' : 'border: 2px solid #3b82f6;'">
          <img :src="user?.role === 'tenant' ? 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=120&q=80' : 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?auto=format&fit=crop&w=120&q=80'" alt="Avatar" />
        </div>
        <div class="user-info-text">
          <span class="user-name-text">{{ user?.full_name || (user?.role === 'tenant' ? 'สมชาย ใจดี' : 'ผู้ดูแลระบบ') }}</span>
          <span class="user-role-text" :style="user?.role === 'tenant' ? 'color: var(--tx-059669); font-weight: 700;' : 'color: var(--tx-2563eb); font-weight: 700;'">
            {{ roleLabel }}
          </span>
        </div>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--tx-94a3b8); margin-left: 2px;">
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>

        <!-- Dropdown Menu -->
        <div v-if="showDropdown" class="card" style="position: absolute; top: 110%; right: 0; min-width: 190px; padding: 8px; z-index: 50; box-shadow: var(--shadow-lg);">
          <div style="padding: 6px 10px; font-size: var(--font-size-xs); color: var(--tx-64748b); border-bottom: 1px solid var(--bd-f1f5f9); margin-bottom: 4px;">
            บทบาทปัจจุบัน: <strong>{{ roleLabel }}</strong>
          </div>
          <button class="btn btn-sm btn-outline" style="width: 100%; border: none; justify-content: flex-start; color: var(--tx-ef4444); display: flex; align-items: center; gap: 8px;" @click="handleLogout">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
            <span>ออกจากระบบ</span>
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const props = defineProps<{
  title?: string
  subtitle?: string
}>()

const { user, logout } = useAuth()
const { isExpanded, toggle: toggleSidebar } = useSidebar()
const showDropdown = ref(false)
const roleLabel = computed(() => {
  if (user.value?.role === 'tenant') return 'ผู้เช่าแผงค้า'
  if (user.value?.role === 'staff') {
    const marketName = (user.value as any)?.assigned_markets?.[0]?.name
    return marketName ? `เจ้าหน้าที่${marketName}` : 'เจ้าหน้าที่ตลาด'
  }
  return 'ผู้ดูแลระบบ (Admin)'
})

const toggleMenu = () => {
  showDropdown.value = !showDropdown.value
}

const { showConfirm } = useAppAlert()

const handleLogout = async () => {
  const confirmed = await showConfirm('คุณต้องการออกจากระบบหรือไม่?', 'ยืนยันออกจากระบบ')
  if (confirmed) {
    await logout()
  }
}
</script>

<style scoped>
.navbar-subtitle-text {
  font-size: var(--font-size-sm);
  color: var(--tx-64748b);
  font-weight: 500;
  margin-left: 6px;
}


@media (max-width: 768px) {
  .navbar-subtitle-text {
    display: none;
  }
}
</style>
