<template>
  <header class="top-navbar">
    <div class="navbar-left">
      <!-- Hamburger Button (Expand / Collapse Sidebar) -->
      <button
        class="btn-icon-square hamburger-toggle-btn"
        :class="{ active: isExpanded }"
        @click="toggleSidebar"
        :title="isExpanded ? 'ย่อแถบเมนู (Collapse Sidebar)' : 'ขยายแถบเมนู (Expand Sidebar)'"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="3" y1="12" x2="21" y2="12"></line>
          <line x1="3" y1="6" x2="21" y2="6"></line>
          <line x1="3" y1="18" x2="21" y2="18"></line>
        </svg>
      </button>

      <!-- System Brand Title -->
      <NuxtLink :to="homeLink" class="navbar-brand-tag" title="ระบบบริหารจัดการแผงเช่า สำนักงานตลาดกรุงเทพมหานคร">
        <span class="brand-title-text">e-Market</span>
        <span class="bma-pill-badge">
          <span class="bma-dot"></span>
          <span class="bma-text-full">ระบบบริหารจัดการแผงเช่า สำนักงานตลาดกรุงเทพมหานคร</span>
          <span class="bma-text-short">ระบบบริหารจัดการแผงเช่า กทม.</span>
        </span>
      </NuxtLink>
    </div>

    <div class="navbar-right" ref="headerRightRef">
      <!-- สลับธีมสว่าง/มืด (อยู่ข้างปุ่มแจ้งเตือน) -->
      <button
        type="button"
        class="theme-toggle-btn"
        :title="isDark ? 'เปลี่ยนเป็นโหมดสว่าง' : 'เปลี่ยนเป็นโหมดมืด'"
        :aria-label="isDark ? 'เปลี่ยนเป็นโหมดสว่าง' : 'เปลี่ยนเป็นโหมดมืด'"
        :aria-pressed="isDark"
        @click="toggleTheme"
      >
        <svg v-if="isDark" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="4"></circle><path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"></path></svg>
        <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
      </button>

      <!-- Notification Bell & Dropdown Central Component -->
      <AppNotificationDropdown
        ref="notifDropdownRef"
        @open="showDropdown = false"
      />

      <!-- User Profile Badge -->
      <div class="user-profile-badge" @click.stop="toggleMenu" style="position: relative; cursor: pointer;">
        <div class="user-avatar-circle" :style="user?.role === 'tenant' ? 'border: 2px solid #10b981; box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);' : 'border: 2px solid #3b82f6; box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);'">
          <img :src="user?.role === 'tenant' ? 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=120&q=80' : 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?auto=format&fit=crop&w=120&q=80'" alt="Avatar" />
        </div>
        <div class="user-info-text">
          <span class="user-name-text">{{ user?.full_name || (user?.role === 'tenant' ? 'สมชาย ใจดี' : 'ผู้ดูแลระบบ') }}</span>
          <span class="user-role-text" :style="user?.role === 'tenant' ? 'color: #059669; font-weight: 700;' : 'color: #2563eb; font-weight: 700;'">
            {{ roleLabel }}
          </span>
        </div>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="user-dropdown-arrow" :class="{ open: showDropdown }">
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>

        <!-- Dropdown Menu -->
        <div v-if="showDropdown" class="profile-dropdown-card" @click.stop>
          <div class="dropdown-header">
            <div style="font-size: 0.75rem; color: #64748b; font-weight: 600;">บัญชีผู้ใช้งาน</div>
            <div style="font-weight: 700; font-size: 0.88rem; color: #0f172a; margin-top: 2px;">
              {{ user?.full_name || (user?.role === 'tenant' ? 'สมชาย ใจดี' : 'ผู้ดูแลระบบ') }}
            </div>
            <div style="font-size: 0.75rem; color: #64748b;">
              {{ user?.email || (user?.role === 'tenant' ? 'somchai@tenant.com' : 'admin@emarket.com') }}
            </div>
          </div>
          <div style="padding: 4px;">
            <NuxtLink v-if="user?.role === 'admin'" to="/settings" class="dropdown-item">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="3"></circle>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
              </svg>
              <span>ตั้งค่าโซนตลาด</span>
            </NuxtLink>
            <button class="dropdown-item logout-action" @click="handleLogout">
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
    </div>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
  subtitle?: string
}>()

const { user, logout } = useAuth()
const { isExpanded, toggle: toggleSidebar } = useSidebar()
const { isDark, toggle: toggleTheme } = useTheme()

const showDropdown = ref(false)
const headerRightRef = ref<HTMLElement | null>(null)
const notifDropdownRef = ref<{ closeNotifications: () => void } | null>(null)
const assignedMarket = computed(() => (user.value as any)?.assigned_markets?.[0])
const homeLink = computed(() => user.value?.role === 'staff' && assignedMarket.value
  ? `/stalls/manage?market=${assignedMarket.value.id}&mode=inspect`
  : '/')
const roleLabel = computed(() => {
  if (user.value?.role === 'tenant') return 'ผู้เช่าแผงค้า'
  if (user.value?.role === 'staff') return assignedMarket.value ? `เจ้าหน้าที่${assignedMarket.value.name}` : 'เจ้าหน้าที่ตลาด'
  return 'ผู้ดูแลระบบตลาด กทม. (Admin)'
})

const toggleMenu = () => {
  notifDropdownRef.value?.closeNotifications()
  showDropdown.value = !showDropdown.value
}

const handleClickOutside = (e: MouseEvent) => {
  if (headerRightRef.value && !headerRightRef.value.contains(e.target as Node)) {
    showDropdown.value = false
  }
}

onMounted(() => {
  window.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  window.removeEventListener('click', handleClickOutside)
})

const { showConfirm } = useAppAlert()

const handleLogout = async () => {
  const confirmed = await showConfirm('คุณต้องการออกจากระบบหรือไม่?', 'ยืนยันออกจากระบบ')
  if (confirmed) {
    await logout()
  }
}
</script>

<style scoped>
.navbar-brand-tag {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  transition: opacity 0.15s ease;
  user-select: none;
}

.navbar-brand-tag:hover {
  opacity: 0.88;
}

.brand-title-text {
  font-size: 1.15rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.02em;
  line-height: 1;
}

.bma-pill-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #eff6ff;
  border: 1px solid #dbeafe;
  color: #1e40af;
  padding: 4px 11px;
  border-radius: 9999px;
  font-size: 0.78rem;
  font-weight: 700;
  line-height: 1.2;
}

.bma-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #2563eb;
  flex-shrink: 0;
}

.bma-text-full {
  display: inline;
}

.bma-text-short {
  display: none;
}

@media (max-width: 860px) {
  .bma-text-full { display: none; }
  .bma-text-short { display: inline; }
}

@media (max-width: 768px) {
  .bma-pill-badge { display: none !important; }
}

.user-dropdown-arrow {
  color: #94a3b8;
  margin-left: 2px;
  transition: transform 0.2s ease;
}

.user-dropdown-arrow.open {
  transform: rotate(180deg);
}

.profile-dropdown-card {
  position: absolute;
  top: 115%;
  right: 0;
  min-width: 230px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  box-shadow: 0 12px 32px rgba(15, 23, 42, 0.12);
  z-index: 50;
  overflow: hidden;
  animation: scaleUp 0.18s cubic-bezier(0.16, 1, 0.3, 1);
}

.dropdown-header {
  padding: 12px 14px;
  background: #f8fafc;
  border-bottom: 1px solid #f1f5f9;
}

.dropdown-item {
  width: 100%;
  padding: 8px 12px;
  border: none;
  background: transparent;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.84rem;
  font-weight: 600;
  color: #334155;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  text-decoration: none;
}

.dropdown-item:hover {
  background: #f1f5f9;
  color: #2563eb;
}

.dropdown-item.logout-action {
  color: #ef4444;
}

.dropdown-item.logout-action:hover {
  background: #fff1f2;
  color: #dc2626;
}

.theme-toggle-btn {
  width: 40px;
  height: 40px;
  flex: 0 0 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e2e8f0;
  border-radius: 50%;
  background: #ffffff;
  color: #475569;
  cursor: pointer;
  transition: background 0.2s, color 0.2s, border-color 0.2s;
}
.theme-toggle-btn:hover { background: #f8fafc; color: #0f172a; }
.theme-toggle-btn:focus-visible { outline: 3px solid #93c5fd; outline-offset: 2px; }

@media (max-width: 768px) {
  .top-navbar {
    position: sticky !important;
    top: 0 !important;
    z-index: 35 !important;
    background: #ffffff !important;
  }
}
</style>
