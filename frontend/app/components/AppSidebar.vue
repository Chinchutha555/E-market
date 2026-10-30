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
      <!-- เจ้าหน้าที่ตลาด: เข้าถึงงานของตลาดตนเองจาก Sidebar โดยตรง -->
      <template v-if="user?.role === 'staff'">
        <!-- 1. ข้อมูลตลาด (แผนที่ผังตลาดแบบเต็มจอ) -->
        <NuxtLink
          :to="staffMarketMapLink"
          class="nav-icon-link"
          :class="{ active: route.path === '/stalls/map' }"
          :title="isExpanded ? '' : 'ข้อมูลตลาด'"
        >
          <span class="nav-icon-wrap">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"></polygon>
              <line x1="8" y1="2" x2="8" y2="18"></line>
              <line x1="16" y1="6" x2="16" y2="22"></line>
            </svg>
          </span>
          <span v-if="isExpanded" class="nav-label">ข้อมูลตลาด</span>
        </NuxtLink>

      </template>

      <!-- หน้าหลักสำหรับผู้ดูแลระบบและผู้เช่า -->
      <NuxtLink
        v-else
        to="/"
        class="nav-icon-link"
        :class="{ active: isHomeActive }"
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
        @click="acknowledgeMenu('bookings')"
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
        <span v-if="displayBookingActionCount > 0" class="sidebar-action-badge" :aria-label="`มีรายการจองที่ต้องจัดการ ${displayBookingActionCount} รายการ`">{{ formatBadgeCount(displayBookingActionCount) }}</span>
      </NuxtLink>

      <!-- 3. Applications -->
      <NuxtLink
        to="/applications"
        class="nav-icon-link"
        :class="{ active: route.path.startsWith('/applications') }"
        :title="isExpanded ? '' : 'คำร้อง'"
        @click="acknowledgeMenu('applications')"
      >
        <span class="nav-icon-wrap">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 11l18-6v14L3 13z"></path>
            <path d="M3 13l2 7h4l-2-6"></path>
            <path d="M21 9a3 3 0 0 1 0 6"></path>
          </svg>
        </span>
        <span v-if="isExpanded" class="nav-label">คำร้อง</span>
        <span v-if="applicationActionCount > 0" class="sidebar-action-badge" :aria-label="`มีคำร้องที่ต้องจัดการ ${applicationActionCount} รายการ`">{{ formatBadgeCount(applicationActionCount) }}</span>
      </NuxtLink>

      <NuxtLink to="/debts" class="nav-icon-link" :class="{ active: route.path.startsWith('/debts') }" :title="isExpanded ? '' : 'หนี้ค้างชำระ'">
        <span class="nav-icon-wrap"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="5" width="18" height="14" rx="2"/><path d="M3 10h18M7 15h3"/></svg></span>
        <span v-if="isExpanded" class="nav-label">หนี้ค้างชำระ</span>
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

      <!-- เจ้าหน้าที่ตลาด: จัดการตลาด (กลุ่มเมนูย่อย) อยู่ล่างสุดของเมนูหลัก -->
      <template v-if="user?.role === 'staff'">
        <!-- 2. จัดการตลาด (กลุ่มเมนูย่อย) -->
        <button
          type="button"
          class="nav-icon-link staff-nav-group-toggle"
          :class="{ active: !isExpanded && isMarketManagementRoute }"
          :aria-expanded="isExpanded ? isMarketManagementExpanded : undefined"
          :title="isExpanded ? '' : 'จัดการตลาด'"
          @click="toggleMarketManagement"
        >
          <span class="nav-icon-wrap">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M4 10h16v10H4z"></path>
              <path d="M3 10l2-6h14l2 6"></path>
              <path d="M3 10a3 3 0 0 0 6 0 3 3 0 0 0 6 0 3 3 0 0 0 6 0"></path>
              <path d="M9 20v-5h6v5"></path>
            </svg>
          </span>
          <span v-if="isExpanded" class="nav-label">จัดการตลาด</span>
          <svg v-if="isExpanded" class="staff-nav-group-chevron" :class="{ 'is-open': isMarketManagementExpanded }" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>

        <div v-if="isExpanded && isMarketManagementExpanded" class="staff-nav-submenu">
          <NuxtLink
            :to="staffMarketInfoLink"
            class="nav-icon-link staff-nav-submenu-link"
            :class="{ active: route.path === '/stalls' }"
            :title="isExpanded ? '' : 'รายละเอียดตลาด'"
          >
            <span class="nav-icon-wrap">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M3 21h18"></path>
                <path d="M5 21V5l7-3 7 3v16"></path>
                <path d="M9 21v-5h6v5"></path>
                <path d="M9 8h.01"></path>
                <path d="M15 8h.01"></path>
                <path d="M9 12h.01"></path>
                <path d="M15 12h.01"></path>
              </svg>
            </span>
            <span v-if="isExpanded" class="nav-label">รายละเอียดตลาด</span>
          </NuxtLink>

          <NuxtLink
            :to="staffLayoutLink"
            class="nav-icon-link staff-nav-submenu-link"
            :class="{ active: route.path === '/stalls/manage' && route.query.mode !== 'inspect' }"
            :title="isExpanded ? '' : 'จัดการผัง GIS'"
          >
            <span class="nav-icon-wrap">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18-6 3V5l6-3 6 3 6-3v12"></path>
                <path d="M9 2v16"></path>
                <path d="M15 5v6"></path>
                <path d="m14 18 6.5-6.5a1.7 1.7 0 0 0-2.4-2.4L11.6 15.6 11 19z"></path>
              </svg>
            </span>
            <span v-if="isExpanded" class="nav-label">จัดการผัง GIS</span>
          </NuxtLink>

          <NuxtLink
            :to="staffVisitManagementLink"
            class="nav-icon-link staff-nav-submenu-link"
            :class="{ active: route.path.startsWith('/stalls/visit-schedule') && route.query.section !== 'holidays' }"
            :title="isExpanded ? '' : 'จัดการรอบเข้าชม'"
          >
            <span class="nav-icon-wrap">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="9"></circle>
                <path d="M12 7v5l3 2"></path>
              </svg>
            </span>
            <span v-if="isExpanded" class="nav-label">จัดการรอบเข้าชม</span>
          </NuxtLink>

          <NuxtLink
            :to="`${staffVisitManagementLink}&section=holidays`"
            class="nav-icon-link staff-nav-submenu-link"
            :class="{ active: route.query.section === 'holidays' }"
            :title="isExpanded ? '' : 'จัดการวันหยุด'"
          >
            <span class="nav-icon-wrap">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="4" width="18" height="18" rx="2"></rect>
                <line x1="16" y1="2" x2="16" y2="6"></line>
                <line x1="8" y1="2" x2="8" y2="6"></line>
                <line x1="3" y1="10" x2="21" y2="10"></line>
              </svg>
            </span>
            <span v-if="isExpanded" class="nav-label">จัดการวันหยุด</span>
          </NuxtLink>
        </div>
      </template>

      <!-- 5. Settings (Admin Only) -->
      <NuxtLink
        v-if="user?.role === 'admin'"
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

      <!-- 6. Application Types (Admin Only) -->
      <NuxtLink
        v-if="user?.role === 'admin'"
        to="/applications/types"
        class="nav-icon-link"
        :class="{ active: route.path.startsWith('/applications/types') }"
        :title="isExpanded ? '' : 'ประเภทคำร้อง'"
      >
        <span class="nav-icon-wrap">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="4" y1="21" x2="4" y2="14"></line>
            <line x1="4" y1="10" x2="4" y2="3"></line>
            <line x1="12" y1="21" x2="12" y2="12"></line>
            <line x1="12" y1="8" x2="12" y2="3"></line>
            <line x1="20" y1="21" x2="20" y2="16"></line>
            <line x1="20" y1="12" x2="20" y2="3"></line>
            <line x1="1" y1="14" x2="7" y2="14"></line>
            <line x1="9" y1="8" x2="15" y2="8"></line>
            <line x1="17" y1="16" x2="23" y2="16"></line>
          </svg>
        </span>
        <span v-if="isExpanded" class="nav-label">ประเภทคำร้อง</span>
      </NuxtLink>
    </nav>

  </aside>
</template>

<script setup lang="ts">
const route = useRoute()
const { user } = useAuth()
const { isExpanded, expand, collapse } = useSidebar()
const { notifications, activeRole } = useNotifications()

const isHomeActive = computed(() => {
  if (route.path === '/') return true
  // สำหรับแอดมินและผู้เช่า: หน้ารายละเอียดตลาดและหน้าจัดการตลาดจะเชื่อมโยงมาจากหน้าหลัก (ตาม breadcrumb)
  if (user.value?.role !== 'staff') {
    return route.path.startsWith('/stalls') || route.path.startsWith('/markets')
  }
  return false
})
const staffMarketId = computed(() => (user.value as any)?.assigned_markets?.[0]?.id)
const staffMarketMapLink = computed(() =>
  staffMarketId.value ? `/stalls/map?market=${staffMarketId.value}` : '/staff'
)
const staffMarketInfoLink = computed(() =>
  staffMarketId.value ? `/stalls?market=${staffMarketId.value}` : '/staff'
)
const staffLayoutLink = computed(() =>
  staffMarketId.value ? `/stalls/manage?market=${staffMarketId.value}` : '/staff'
)
const staffVisitManagementLink = computed(() =>
  staffMarketId.value ? `/stalls/visit-schedule?market=${staffMarketId.value}` : '/staff'
)
const isMarketManagementRoute = computed(() =>
  route.path === '/stalls' ||
  (route.path === '/stalls/manage' && route.query.mode !== 'inspect') ||
  route.path.startsWith('/stalls/visit-schedule')
)
const isMarketManagementExpanded = ref(isMarketManagementRoute.value)
const toggleMarketManagement = () => {
  if (!isExpanded.value) {
    isMarketManagementExpanded.value = true
    expand()
    return
  }
  isMarketManagementExpanded.value = !isMarketManagementExpanded.value
}
const sidebarSeenByMenu = ref<Record<'bookings' | 'applications', number>>({
  bookings: 0,
  applications: 0
})

const sidebarAttentionStorageKey = computed(() => {
  const userId = (user.value as any)?.id || 'anonymous'
  return `e-market:sidebar-attention:${userId}:${activeRole.value}`
})

const loadSidebarAttentionState = () => {
  if (!process.client) return
  try {
    const saved = JSON.parse(localStorage.getItem(sidebarAttentionStorageKey.value) || '{}')
    sidebarSeenByMenu.value = {
      bookings: Number(saved.bookings) || 0,
      applications: Number(saved.applications) || 0
    }
  } catch {
    sidebarSeenByMenu.value = { bookings: 0, applications: 0 }
  }
}

const notificationMenu = (notification: any) => {
  const status = String(notification?.status_code || notification?.type || '').toLowerCase()
  if (status === 'lease_applied' || status === 'completed') return 'applications'
  const target = String(notification?.target_url || '')
  if (target.startsWith('/applications')) return 'applications'
  return target.startsWith('/bookings') ? 'bookings' : ''
}

const menuNotifications = (menu: 'bookings' | 'applications') => notifications.value.filter((notification: any) =>
  notification.recipient_role === activeRole.value &&
  notificationMenu(notification) === menu
)

const menuActionCount = (menu: 'bookings' | 'applications') => computed(() => menuNotifications(menu).filter((notification: any) =>
  Number(notification.id) > sidebarSeenByMenu.value[menu]
).length)

const acknowledgeMenu = (menu: 'bookings' | 'applications') => {
  const latestId = menuNotifications(menu).reduce((maxId: number, notification: any) =>
    Math.max(maxId, Number(notification.id) || 0), 0)
  sidebarSeenByMenu.value = {
    ...sidebarSeenByMenu.value,
    [menu]: latestId
  }
  if (process.client) {
    localStorage.setItem(sidebarAttentionStorageKey.value, JSON.stringify(sidebarSeenByMenu.value))
  }
}

const bookingActionCount = menuActionCount('bookings')
const applicationActionCount = menuActionCount('applications')
const sidebarBadgePreviewCount = computed(() => {
  if (!import.meta.dev) return 0
  const count = Number(route.query.sidebarBadge)
  return Number.isFinite(count) && count > 0 ? Math.floor(count) : 0
})
const displayBookingActionCount = computed(() => sidebarBadgePreviewCount.value || bookingActionCount.value)
const formatBadgeCount = (count: number) => count >= 10 ? '9+' : String(count)

onMounted(loadSidebarAttentionState)

watch(sidebarAttentionStorageKey, loadSidebarAttentionState)

watch(() => route.fullPath, () => {
  if (isMarketManagementRoute.value) {
    isMarketManagementExpanded.value = true
  }
  if (process.client && window.innerWidth <= 768) {
    collapse()
  }
})
</script>

<style scoped>
.sidebar {
  border-right: 1px solid rgba(226, 232, 240, 0.85);
  box-shadow: 2px 0 16px rgba(15, 23, 42, 0.03);
}

.brand-icon-box {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.28);
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.sidebar-action-badge {
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-left: auto;
  border-radius: 999px;
  background: #ef4444;
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 700;
  line-height: 1;
  box-shadow: 0 2px 6px rgba(239, 68, 68, 0.25);
}

.sidebar:not(.expanded) .sidebar-action-badge {
  position: absolute;
  top: 5px;
  right: 5px;
  min-width: 17px;
  height: 17px;
  padding: 0 4px;
  font-size: 0.6875rem;
  border: 2px solid #ffffff;
}

.brand-icon-box:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 6px 18px rgba(37, 99, 235, 0.38);
}

.nav-icon-link {
  border-radius: 12px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.nav-icon-link:hover {
  background: #f8fafc;
  color: #2563eb;
  transform: translateX(2px);
}

.nav-icon-link.active {
  background: #eff6ff;
  color: #2563eb;
  box-shadow: inset 0 0 0 1px rgba(191, 219, 254, 0.5);
}

.staff-nav-group-toggle {
  width: 100%;
  margin: 4px 0 0;
  padding: 0 14px;
  border: 0;
  background: transparent;
  font-family: inherit;
  color: #475569;
  font-weight: 600;
  cursor: pointer;
  text-align: inherit;
}

.staff-nav-group-toggle:hover {
  background: #f8fafc;
  color: #2563eb;
  transform: none;
}

.staff-nav-group-toggle .nav-label {
  flex: 1;
}

.staff-nav-group-chevron {
  flex-shrink: 0;
  transition: transform 0.18s ease;
}

.staff-nav-group-chevron.is-open {
  transform: rotate(180deg);
}

.staff-nav-submenu {
  display: grid;
  gap: 2px;
  margin: -1px 0 9px 23px;
  padding-left: 8px;
  border-left: 1px solid #e2e8f0;
}

.sidebar.expanded .staff-nav-submenu .staff-nav-submenu-link {
  height: 38px;
  padding-left: 10px;
  border-radius: 10px;
}

.sidebar.expanded .staff-nav-submenu .staff-nav-submenu-link .nav-icon-wrap {
  display: none;
}

.sidebar.expanded .staff-nav-submenu .staff-nav-submenu-link:hover {
  background: transparent;
  transform: none;
}

.sidebar.expanded .staff-nav-submenu .staff-nav-submenu-link.active {
  background: transparent;
  box-shadow: none;
  color: #2563eb;
}

.sidebar.expanded .staff-nav-submenu .staff-nav-submenu-link.active::before {
  display: none;
}

.sidebar.expanded .staff-nav-submenu .staff-nav-submenu-link.active .nav-label {
  color: #2563eb;
  font-weight: 600;
}

@media (max-width: 768px) {
  .brand-subtitle {
    display: none !important;
  }

  .brand-text-block {
    display: flex !important;
    flex-direction: column !important;
    justify-content: center !important;
  }
}
</style>
