<template>
  <aside class="sidebar" :class="{ 'expanded': isExpanded }">
    <!-- Brand Header -->
    <div class="sidebar-brand-wrapper">
      <NuxtLink to="/" class="sidebar-brand" title="e-Market สำนักงานตลาดกรุงเทพมหานคร">
        <div class="brand-icon-box">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path>
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <path d="M16 10a4 4 0 0 1-8 0"></path>
          </svg>
        </div>
        <div v-if="isExpanded" class="brand-text-block">
          <span class="brand-title">e-Market</span>
          <span class="brand-subtitle">สำนักงานตลาด กทม.</span>
        </div>
      </NuxtLink>
    </div>

    <!-- Nav Items -->
    <nav class="sidebar-nav">
      <!-- 1. Home: 12 Markets Directory / Stalls Map -->
      <NuxtLink
        to="/"
        class="nav-icon-link"
        :class="{ active: route.path === '/' || route.path.startsWith('/stalls') }"
        :title="isExpanded ? '' : 'หน้าแรก'"
      >
        <span class="nav-icon-wrap">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
            <polyline points="9 22 9 12 15 12 15 22"></polyline>
          </svg>
        </span>
        <span v-if="isExpanded" class="nav-label">หน้าแรก</span>
      </NuxtLink>

      <!-- 2. Bookings -->
      <NuxtLink
        to="/bookings"
        class="nav-icon-link"
        :class="{ active: route.path.startsWith('/bookings') }"
        :title="isExpanded ? '' : 'การจองแผงเช่า'"
      >
        <span class="nav-icon-wrap">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="4" width="18" height="18" rx="2"></rect>
            <line x1="16" y1="2" x2="16" y2="6"></line>
            <line x1="8" y1="2" x2="8" y2="6"></line>
            <line x1="3" y1="10" x2="21" y2="10"></line>
          </svg>
        </span>
        <span v-if="isExpanded" class="nav-label">การจองแผงเช่า</span>
      </NuxtLink>

      <!-- 3. Applications -->
      <NuxtLink
        to="/applications"
        class="nav-icon-link"
        :class="{ active: route.path.startsWith('/applications') }"
        :title="isExpanded ? '' : 'คำร้อง'"
      >
        <span class="nav-icon-wrap">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 11l18-6v14L3 13z"></path>
            <path d="M3 13l2 7h4l-2-6"></path>
            <path d="M21 9a3 3 0 0 1 0 6"></path>
          </svg>
        </span>
        <span v-if="isExpanded" class="nav-label">คำร้อง</span>
      </NuxtLink>

      <!-- 4. Contracts -->
      <NuxtLink
        to="/contracts"
        class="nav-icon-link"
        :class="{ active: route.path.startsWith('/contracts') }"
        :title="isExpanded ? '' : (user?.role === 'tenant' ? 'แผงเช่าและสัญญาของฉัน' : 'สัญญาการเช่า')"
      >
        <span class="nav-icon-wrap">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="16" y1="13" x2="8" y2="13"></line>
            <line x1="16" y1="17" x2="8" y2="17"></line>
          </svg>
        </span>
        <span v-if="isExpanded" class="nav-label">
          {{ user?.role === 'tenant' ? 'แผงเช่า & สัญญาของฉัน' : 'สัญญาการเช่า' }}
        </span>
      </NuxtLink>

      <!-- 5. Settings (Admin Only) -->
      <NuxtLink
        v-if="user?.role !== 'tenant'"
        to="/settings"
        class="nav-icon-link"
        :class="{ active: route.path === '/settings' }"
        :title="isExpanded ? '' : 'ตั้งค่าโซนตลาด'"
      >
        <span class="nav-icon-wrap">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
          </svg>
        </span>
        <span v-if="isExpanded" class="nav-label">ตั้งค่าโซนตลาด</span>
      </NuxtLink>
    </nav>

    <!-- Footer Logout Button -->
    <div class="sidebar-footer">
      <button class="nav-icon-link logout-btn" @click="handleLogout" :title="isExpanded ? '' : 'ออกจากระบบ (Logout)'">
        <span class="nav-icon-wrap">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
            <polyline points="16 17 21 12 16 7"></polyline>
            <line x1="21" y1="12" x2="9" y2="12"></line>
          </svg>
        </span>
        <span v-if="isExpanded" class="nav-label" style="color: #ef4444;">ออกจากระบบ</span>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
const route = useRoute()
const { user, logout } = useAuth()
const { isExpanded } = useSidebar()
const { showConfirm } = useAppAlert()

const handleLogout = async () => {
  const confirmed = await showConfirm('คุณต้องการออกจากระบบหรือไม่?', 'ยืนยันออกจากระบบ')
  if (confirmed) {
    await logout()
  }
}
</script>
