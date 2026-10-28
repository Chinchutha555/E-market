<template>
  <div class="application-main-page">
    <AppHeader />
    <div class="page-container">
      <!-- Loading Skeleton State -->
      <div v-if="loading" class="skeleton-loading-container">
        <div class="skeleton-line skeleton-header"></div>
        <div class="skeleton-cards-grid">
          <div v-for="n in 4" :key="n" class="skeleton-card"></div>
        </div>
      </div>

      <!-- Main Container (Table Card Container) -->
      <div v-else class="table-card-container">
        <!-- Header Row: Title & Subtitle on Left, Search + Actions on Right (โครงเดียวกับหน้าจอง) -->
        <div class="card-header-row">
          <div class="header-left">
            <h2 class="section-title">รายการคำร้องทั้งหมด</h2>
            <div class="item-count-subtitle">ทั้งหมด {{ filteredApplications.length }} รายการ</div>
          </div>

          <div class="header-right">
            <!-- Search Input -->
            <div class="search-input-wrapper">
              <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="ค้นหาหมายเลขคำร้อง หรือผู้ยื่น"
                class="search-input"
                @input="currentPage = 1"
              />
              <button
                v-if="searchQuery"
                class="clear-search-btn"
                type="button"
                aria-label="ล้างการค้นหา"
                @click="searchQuery = ''; currentPage = 1"
              >
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              </button>
            </div>

            <!-- View Mode Switcher: Cards vs Table -->
            <AppSegmentedControl
              v-model="viewMode"
              :options="viewModeOptions"
              aria-label="เลือกรูปแบบการแสดงรายการคำร้อง"
            />

            <!-- Primary Action: ยื่นคำร้องทั่วไป -->
            <button type="button" class="btn-table-primary-action" @click="openCreateGeneralModal" title="ยื่นคำร้องทั่วไป">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
              <span>ยื่นคำร้องทั่วไป</span>
            </button>

            <!-- Manage Application Types (Admin & Staff) -->
            <NuxtLink
              v-if="isAdminOrStaff"
              to="/applications/types"
              class="btn-table-manage-types"
              title="จัดการประเภทคำร้อง"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="4" y1="21" x2="4" y2="14"></line><line x1="4" y1="10" x2="4" y2="3"></line><line x1="12" y1="21" x2="12" y2="12"></line><line x1="12" y1="8" x2="12" y2="3"></line><line x1="20" y1="21" x2="20" y2="16"></line><line x1="20" y1="12" x2="20" y2="3"></line><line x1="1" y1="14" x2="7" y2="14"></line><line x1="9" y1="8" x2="15" y2="8"></line><line x1="17" y1="16" x2="23" y2="16"></line></svg>
              <span>จัดการประเภทคำร้อง</span>
            </NuxtLink>

            <!-- Reset Test Data Button -->
            <AppButton
              variant="danger-outline"
              size="sm"
              @click="handleClearAllData"
              title="รีเซ็ตข้อมูลทดสอบ"
            >
              <template #prefix>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
              </template>
              <span>รีเซ็ต</span>
            </AppButton>
          </div>
        </div>

        <!-- 4-Card Status Overview KPI Ribbon (สไตล์เดียวกับหน้าจอง) -->
        <div class="officer-kpi-ribbon officer-status-kpi-ribbon">
          <!-- 1. รายการทั้งหมด (Blue Theme) -->
          <div
            class="kpi-ribbon-card kpi-all"
            :class="{ active: activeTab === 'all' }"
            @click="activeTab = 'all'; currentPage = 1"
            role="button"
            tabindex="0"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path><rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">ทั้งหมด</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countAll }}</span>
                <span class="kpi-unit-pill">รายการ</span>
              </div>
            </div>
          </div>

          <!-- 2. รอดำเนินการ (Amber Theme) พร้อม Popover Breakdown -->
          <div
            class="kpi-ribbon-card kpi-pending"
            :class="{ active: activeTab === 'pending' }"
            @click="activeTab = 'pending'; currentPage = 1"
            role="button"
            tabindex="0"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">รอดำเนินการ</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countPendingCombined }}</span>
                <span class="kpi-unit-pill">รายการ</span>
              </div>
            </div>

            <!-- Info Tooltip Button (ดูรายละเอียดขั้นตอนที่ค้าง) -->
            <div
              class="kpi-info-wrapper"
              @mouseenter="showPendingTooltip = true"
              @mouseleave="showPendingTooltip = false"
              @click.stop="showPendingTooltip = !showPendingTooltip"
            >
              <button
                type="button"
                class="kpi-info-btn"
                aria-label="รายละเอียดขั้นตอนรอดำเนินการ"
                title="รายละเอียดขั้นตอนรอดำเนินการ"
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="16" x2="12" y2="12"></line>
                  <line x1="12" y1="8" x2="12.01" y2="8"></line>
                </svg>
              </button>

              <!-- Rich Tooltip Popover -->
              <transition name="fade-slide">
                <div v-if="showPendingTooltip" class="kpi-info-popover" @click.stop>
                  <div class="kpi-popover-header">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
                    <span>ขั้นตอนรอดำเนินการ ({{ countPendingCombined }})</span>
                  </div>
                  <div class="kpi-popover-body">
                    <div class="kpi-popover-item" :class="{ 'is-selected': pendingSubFilter === 'pending_officer' }" @click="pendingSubFilter = 'pending_officer'; activeTab = 'pending'; currentPage = 1">
                      <span class="kpi-status-dot amber"></span>
                      <span class="kpi-status-name">รอ จนท. ตรวจสอบ</span>
                      <span class="kpi-status-val">{{ countPendingOfficer }}</span>
                    </div>
                    <div class="kpi-popover-item" :class="{ 'is-selected': pendingSubFilter === 'pending_debt' }" @click="pendingSubFilter = 'pending_debt'; activeTab = 'pending'; currentPage = 1">
                      <span class="kpi-status-dot blue"></span>
                      <span class="kpi-status-name">รอ ตรวจสอบหนี้</span>
                      <span class="kpi-status-val">{{ countPendingDebt }}</span>
                    </div>
                    <div class="kpi-popover-item" :class="{ 'is-selected': pendingSubFilter === 'pending_market_director' }" @click="pendingSubFilter = 'pending_market_director'; activeTab = 'pending'; currentPage = 1">
                      <span class="kpi-status-dot purple"></span>
                      <span class="kpi-status-name">รอ ผอ.ตลาด</span>
                      <span class="kpi-status-val">{{ countPendingMarketDirector }}</span>
                    </div>
                    <div class="kpi-popover-item" :class="{ 'is-selected': pendingSubFilter === 'pending_bureau_director' }" @click="pendingSubFilter = 'pending_bureau_director'; activeTab = 'pending'; currentPage = 1">
                      <span class="kpi-status-dot green"></span>
                      <span class="kpi-status-name">รอ ผอ.สำนักตลาด</span>
                      <span class="kpi-status-val">{{ countPendingBureauDirector }}</span>
                    </div>
                  </div>
                  <div v-if="pendingSubFilter !== 'all'" class="kpi-popover-reset" @click="pendingSubFilter = 'all'; currentPage = 1">
                    <span>แสดงทั้งหมดในรอดำเนินการ</span>
                  </div>
                </div>
              </transition>
            </div>
          </div>

          <!-- 3. อนุมัติแล้ว (Purple Theme per screenshot) -->
          <div
            class="kpi-ribbon-card kpi-applied"
            :class="{ active: activeTab === 'approved' }"
            @click="activeTab = 'approved'; currentPage = 1"
            role="button"
            tabindex="0"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
                <polyline points="10 9 9 9 8 9"></polyline>
              </svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">อนุมัติแล้ว</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countApproved }}</span>
                <span class="kpi-unit-pill">รายการ</span>
              </div>
            </div>
          </div>

          <!-- 4. ส่งกลับ/ไม่อนุมัติ (Soft Gray Theme) -->
          <div
            class="kpi-ribbon-card kpi-cancelled"
            :class="{ active: activeTab === 'rejected' }"
            @click="activeTab = 'rejected'; currentPage = 1"
            role="button"
            tabindex="0"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">ยกเลิก/ส่งกลับ</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countRejected }}</span>
                <span class="kpi-unit-pill">รายการ</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Officer Sub-filters when "รอดำเนินการ" is Active -->
        <transition name="fade-slide">
          <div v-if="activeTab === 'pending'" class="tenant-subfilter-bar" style="margin-top: -12px; margin-bottom: 22px;">
            <button
              type="button"
              class="subfilter-chip chip-all"
              :class="{ active: pendingSubFilter === 'all' }"
              @click="pendingSubFilter = 'all'; currentPage = 1"
            >
              <span>ทั้งหมด</span>
              <span class="subfilter-badge">{{ countPendingCombined }}</span>
            </button>
            <button
              type="button"
              class="subfilter-chip chip-amber"
              :class="{ active: pendingSubFilter === 'pending_officer' }"
              @click="pendingSubFilter = 'pending_officer'; currentPage = 1"
            >
              <span class="subfilter-dot"></span>
              <span>รอ จนท. ตรวจสอบ</span>
              <span class="subfilter-badge">{{ countPendingOfficer }}</span>
            </button>
            <button
              type="button"
              class="subfilter-chip chip-blue"
              :class="{ active: pendingSubFilter === 'pending_debt' }"
              @click="pendingSubFilter = 'pending_debt'; currentPage = 1"
            >
              <span class="subfilter-dot"></span>
              <span>รอ ตรวจสอบหนี้</span>
              <span class="subfilter-badge">{{ countPendingDebt }}</span>
            </button>
            <button
              type="button"
              class="subfilter-chip chip-purple"
              :class="{ active: pendingSubFilter === 'pending_market_director' }"
              @click="pendingSubFilter = 'pending_market_director'; currentPage = 1"
            >
              <span class="subfilter-dot"></span>
              <span>รอ ผอ.ตลาด</span>
              <span class="subfilter-badge">{{ countPendingMarketDirector }}</span>
            </button>
            <button
              type="button"
              class="subfilter-chip chip-green"
              :class="{ active: pendingSubFilter === 'pending_bureau_director' }"
              @click="pendingSubFilter = 'pending_bureau_director'; currentPage = 1"
            >
              <span class="subfilter-dot"></span>
              <span>รอ ผอ.สำนักตลาด</span>
              <span class="subfilter-badge">{{ countPendingBureauDirector }}</span>
            </button>
          </div>
        </transition>

        <!-- Application Type Filter Bar -->
        <div class="application-type-bar" style="margin-bottom: 22px;">
          <div class="type-filter-pills">
            <button
              type="button"
              class="type-pill"
              :class="{ active: typeFilter === 'all' }"
              @click="typeFilter = 'all'; currentPage = 1"
            >
              <span>ประเภททั้งหมด</span>
              <span class="type-pill-count">{{ applications.length }}</span>
            </button>
            <button
              type="button"
              class="type-pill pill-stall"
              :class="{ active: typeFilter === 'stall_lease' }"
              @click="typeFilter = 'stall_lease'; currentPage = 1"
            >
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="3" width="7" height="7"></rect>
                <rect x="14" y="3" width="7" height="7"></rect>
                <rect x="14" y="14" width="7" height="7"></rect>
                <rect x="3" y="14" width="7" height="7"></rect>
              </svg>
              <span>คำร้องขอเช่าแผงค้า</span>
              <span class="type-pill-count">{{ countByType('stall_lease') }}</span>
            </button>
            <button
              type="button"
              class="type-pill pill-general"
              :class="{ active: typeFilter === 'general' }"
              @click="typeFilter = 'general'; currentPage = 1"
            >
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
              </svg>
              <span>คำร้องทั่วไป</span>
              <span class="type-pill-count">{{ countByType('general') }}</span>
            </button>
          </div>
        </div>

        <!-- View 1: Cards Accordion Grouped View (Mockup Standard per screenshot) -->
        <div v-if="viewMode === 'card'" class="application-cards-container">
          <div v-if="paginatedProjectGroups.length > 0" class="officer-groups-wrapper">
            <div
              v-for="project in paginatedProjectGroups"
              :key="project.key"
              class="project-stall-entry"
            >
              <!-- Accordion Header -->
              <button
                type="button"
                class="officer-project-header"
                :class="{ open: isProjectGroupOpen(project.key) }"
                @click="toggleProjectGroup(project.key)"
              >
                <div class="header-left-content">
                  <svg class="chevron-indicator" :class="{ rotated: isProjectGroupOpen(project.key) }" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="6 9 12 15 18 9"></polyline></svg>
                  <span class="accordion-header-icon" aria-hidden="true">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
                  </span>
                  <span class="market-title">{{ project.projectName }}</span>
                  <span class="market-count-pill">{{ project.applications.length }} รายการคำร้อง</span>
                </div>
              </button>

              <!-- Accordion Body -->
              <div
                v-show="isProjectGroupOpen(project.key)"
                class="market-accordion-body officer-cards-list in-project"
              >
                <div
                  v-for="app in project.applications"
                  :key="app.id"
                  class="booking-item-card"
                  :class="[getCardBorderClass(app.status), { 'collapsed-card': !isCardExpanded(app.id) }]"
                >
                  <!-- Card Top Row -->
                  <div class="card-top-row">
                    <!-- Top Line: Application ID & Status Badge + Chevron -->
                    <div class="card-top-main">
                      <div class="top-row-left">
                        <div class="circle-status-icon blue-circle">
                          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                            <line x1="16" y1="2" x2="16" y2="6"></line>
                            <line x1="8" y1="2" x2="8" y2="6"></line>
                            <line x1="3" y1="10" x2="21" y2="10"></line>
                          </svg>
                        </div>

                        <div class="booking-number-block">
                          <div class="booking-label-row">
                            <span class="booking-label">หมายเลขคำร้อง</span>
                          </div>
                          <div class="booking-id-text">#{{ app.application_number }}</div>
                          <div class="officer-booked-at">
                            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                            <span>วันที่ยื่น {{ formatSubmissionTimestamp(app.created_at) }}</span>
                          </div>
                        </div>
                      </div>

                      <div class="top-row-right">
                        <!-- Aggregate Status Badge Pill -->
                        <span class="status-pill-badge" :class="getStatusPillClass(app.status)">
                          {{ getStatusBadgeLabel(app.status) }}
                        </span>

                        <!-- Card Collapse Toggle Chevron -->
                        <button class="card-collapse-btn" type="button" @click.stop="toggleCardExpanded(app.id)" title="ย่อ/ขยายการ์ด">
                          <svg
                            class="card-chevron"
                            :class="{ rotated: isCardExpanded(app.id) }"
                            width="16"
                            height="16"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                          >
                            <polyline points="6 9 12 15 18 9"></polyline>
                          </svg>
                        </button>
                      </div>
                    </div>

                    <!-- Compact context row: stall, project, market, and applicant -->
                    <div class="card-meta-tags-row compact-officer-meta">
                      <div class="compact-meta-item compact-meta-stall">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                        <span>แผงค้า</span>
                        <strong>{{ app.stall_number && app.stall_number !== '-' ? app.stall_number : 'คำร้องทั่วไป' }}</strong>
                      </div>

                      <div class="compact-meta-item">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V9"></path><path d="M9 22V12h6v10"></path></svg>
                        <span>{{ getApplicationProjectName(app) }}</span>
                      </div>

                      <div class="compact-meta-item">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 21h18M3 7v14M21 7v14M6 7V3h12v4M9 21v-6a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v6"></path></svg>
                        <span>{{ app.market_name || 'ตลาดนัดจตุจักร' }}</span>
                      </div>

                      <div class="compact-meta-item compact-meta-person">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                        <strong>{{ app.applicant_name }}</strong>
                        <a v-if="app.applicant_phone" :href="'tel:' + app.applicant_phone" @click.stop>{{ app.applicant_phone }}</a>
                      </div>
                    </div>
                  </div>

                  <!-- Collapsed Next Step Bar (เมื่อย่อการ์ด) -->
                  <div
                    v-if="!isCardExpanded(app.id)"
                    class="collapsed-next-step"
                    :class="getProgressTone(app.status)"
                  >
                    <div class="collapsed-next-step-content">
                      <span class="collapsed-next-step-label">ความคืบหน้า</span>
                      <strong>{{ getApplicationProgressText(app) }}</strong>
                    </div>

                    <div class="collapsed-actions-group">
                      <AppButton v-if="app.has_id_card_copy" variant="outline" size="sm" @click.stop="downloadIDCardCopy(app)">สำเนาบัตรประชาชน</AppButton>
                      <AppButton
                        v-if="canTenantEdit(app)"
                        variant="primary"
                        size="sm"
                        @click.stop="openRevision(app)"
                        title="แก้ไขคำร้องและส่งเอกสารตรวจสอบอีกครั้ง"
                      >
                        <template #prefix>
                          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                          </svg>
                        </template>
                        <span>แก้ไขคำร้อง</span>
                      </AppButton>
                      <AppButton
                        variant="secondary"
                        size="sm"
                        @click.stop="openDetails(app)"
                        title="ดูรายละเอียดคำร้อง"
                      >
                        <template #prefix>
                          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"></circle><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8z"></path></svg>
                        </template>
                        <span>รายละเอียด</span>
                      </AppButton>

                      <template v-if="canReview(app) && app.status === 'pending_officer'">
                        <AppButton
                          variant="primary"
                          size="sm"
                          @click.stop="processApproval(app, 'เจ้าหน้าที่ตรวจสอบเอกสารครบถ้วน ส่งต่อ ผอ.ตลาด')"
                          title="อนุมัติส่งต่อ"
                        >
                          <template #prefix>
                            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                          </template>
                          <span>อนุมัติส่งต่อ</span>
                        </AppButton>
                      </template>
                      <template v-else-if="canReview(app) && app.status === 'pending_market_director'">
                        <AppButton
                          variant="primary"
                          size="sm"
                          @click.stop="processApproval(app, 'ผอ.ตลาด พิจารณาอนุมัติคำร้อง')"
                          title="ผอ. อนุมัติ"
                        >
                          <template #prefix>
                            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                          </template>
                          <span>ผอ. อนุมัติ</span>
                        </AppButton>
                      </template>
                      <template v-else-if="canReview(app) && app.status === 'pending_bureau_director'">
                        <AppButton
                          variant="primary"
                          size="sm"
                          @click.stop="processApproval(app, 'ผอ.สำนักตลาด พิจารณาอนุมัติคำร้องเรียบร้อย')"
                          title="ผอ.สำนัก อนุมัติ"
                        >
                          <template #prefix>
                            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                          </template>
                          <span>ผอ.สำนัก อนุมัติ</span>
                        </AppButton>
                      </template>
                      <template v-else-if="app.status === 'approved' && isStallLease(app)">
                        <NuxtLink
                          :to="'/contracts?new_for_stall=' + encodeURIComponent(app.stall_number || app.stall_id || '') + '&app_id=' + app.id"
                          class="btn-collapsed-primary"
                          title="จัดทำสัญญาเช่า"
                        >
                          <span>ทำสัญญา</span>
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg>
                        </NuxtLink>
                      </template>
                      <template v-else-if="app.status === 'contract_pending'">
                        <NuxtLink
                          to="/contracts?status=pending"
                          class="btn-collapsed-primary"
                          style="background: #2563eb;"
                          title="ดูสัญญาเช่าที่รอเจ้าหน้าที่ตรวจสอบ"
                        >
                          <span>ดูสัญญา (รอตรวจ)</span>
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg>
                        </NuxtLink>
                      </template>
                      <template v-else-if="app.status === 'completed'">
                        <NuxtLink
                          :to="'/contracts?stall=' + encodeURIComponent(app.stall_number || app.stall_id || '')"
                          class="btn-collapsed-primary"
                          style="background: #10b981;"
                          title="ดูสัญญาเช่า"
                        >
                          <span>ดูสัญญาเช่า</span>
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg>
                        </NuxtLink>
                      </template>
                    </div>
                  </div>

                  <!-- Expanded Card Body -->
                  <div v-show="isCardExpanded(app.id)">
                    <!-- 4-Phase Stepper Container -->
                    <div class="booking-stepper-container">
                      <div class="stepper-track">
                        <!-- Step 1: ยื่นคำร้องขอเช่า -->
                        <div class="stepper-step" :class="{ completed: isAppStepCompleted(app.status, 1), active: isAppStepActive(app.status, 1) }">
                          <div class="step-circle">
                            <svg v-if="isAppStepCompleted(app.status, 1)" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"></polyline></svg>
                            <span v-else>1</span>
                          </div>
                          <div class="step-label-text">
                            {{ app.status === 'revision_required'
                                ? (isTenant ? 'แก้ไขคำร้องขอเช่า' : 'รอผู้เช่าแก้ไขคำร้อง')
                                : 'ยื่นคำร้องขอเช่า'
                            }}
                          </div>
                        </div>

                        <div class="stepper-connector" :class="{ filled: isAppStepCompleted(app.status, 2) || (isAppStepCompleted(app.status, 1) && app.status !== 'revision_required') }"></div>

                        <!-- Step 2: จนท. ตรวจสอบ/หนี้ -->
                        <div class="stepper-step" :class="{ completed: isAppStepCompleted(app.status, 2), active: isAppStepActive(app.status, 2) }">
                          <div class="step-circle">
                            <svg v-if="isAppStepCompleted(app.status, 2)" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"></polyline></svg>
                            <span v-else>2</span>
                          </div>
                          <div class="step-label-text">จนท. ตรวจสอบ/หนี้</div>
                        </div>

                        <div class="stepper-connector" :class="{ filled: isAppStepCompleted(app.status, 3) }"></div>

                        <!-- Step 3: ผอ.ตลาด พิจารณา -->
                        <div class="stepper-step" :class="{ completed: isAppStepCompleted(app.status, 3), active: isAppStepActive(app.status, 3) }">
                          <div class="step-circle">
                            <svg v-if="isAppStepCompleted(app.status, 3)" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"></polyline></svg>
                            <span v-else>3</span>
                          </div>
                          <div class="step-label-text">ผอ.ตลาด พิจารณา</div>
                        </div>

                        <div class="stepper-connector" :class="{ filled: isAppStepCompleted(app.status, 4) }"></div>

                        <!-- Step 4: ผอ.สำนัก อนุมัติ -->
                        <div class="stepper-step" :class="{ completed: isAppStepCompleted(app.status, 4), active: isAppStepActive(app.status, 4) }">
                          <div class="step-circle">
                            <svg v-if="isAppStepCompleted(app.status, 4)" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"></polyline></svg>
                            <span v-else>4</span>
                          </div>
                          <div class="step-label-text">ผอ.สำนัก อนุมัติ</div>
                        </div>
                      </div>
                    </div>

                    <!-- 4-Column Info Grid (Officer Columns Grid) -->
                    <div class="info-columns-grid officer-columns-grid">
                      <!-- Col 1: วันที่ขอเริ่มสัญญา -->
                      <div class="info-column">
                        <div class="info-label">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                          <span>{{ isStallLease(app) ? 'วันที่ขอเริ่มสัญญา' : 'วันที่ยื่นคำร้อง' }}</span>
                        </div>
                        <div class="info-value-row">
                          <span class="info-value-text">{{ formatThaiStartDate(app.requested_start_date || app.created_at) }}</span>
                          <span v-if="getDaysRemainingBadge(app.requested_start_date)" class="days-remaining-pill">
                            {{ getDaysRemainingBadge(app.requested_start_date) }}
                          </span>
                        </div>
                      </div>

                      <!-- Col 2: ประเภท / เรื่อง -->
                      <div class="info-column">
                        <div class="info-label">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                          <span>{{ isStallLease(app) ? 'ประเภทคำร้อง' : 'เรื่องที่ขอ' }}</span>
                        </div>
                        <div class="info-value-row">
                          <span class="info-value-text">{{ isStallLease(app) ? 'ขอเช่าแผงค้า' : (app.subject || app.shop_name || 'คำร้องทั่วไป') }}</span>
                        </div>
                        <div v-if="isStallLease(app) && app.stall_monthly_rate" class="cell-sub-meeting-point">
                          <span>฿{{ Number(app.stall_monthly_rate).toLocaleString() }}/ด. ({{ app.requested_period_months || 12 }} เดือน)</span>
                        </div>
                        <div v-else-if="app.review_note" class="cell-sub-meeting-point">
                          <span>{{ app.review_note.slice(0, 35) }}{{ app.review_note.length > 35 ? '...' : '' }}</span>
                        </div>
                      </div>

                      <!-- Col 3: แผงค้า -->
                      <div class="info-column">
                        <div class="info-label">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                          <span>แผงค้า</span>
                        </div>
                        <div class="info-value-row">
                          <span class="stall-count-pill">
                            <span class="stall-code-strong">{{ app.stall_number && app.stall_number !== '-' ? app.stall_number : 'คำร้องทั่วไป' }}</span>
                          </span>
                        </div>
                        <div v-if="app.shop_name" class="cell-sub-meeting-point">
                          <span>{{ app.shop_name }}</span>
                        </div>
                      </div>

                      <!-- Col 4: ผู้ขอเช่า / ผู้ยื่นคำร้อง -->
                      <div class="info-column">
                        <div class="info-label">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                          <span>ผู้ยื่นคำร้อง</span>
                        </div>
                        <div class="info-value-row">
                          <span class="info-value-text">{{ app.applicant_name }}</span>
                          <a v-if="app.applicant_phone" :href="'tel:' + app.applicant_phone" class="phone-call-btn" title="คลิกเพื่อโทรออกทันที">
                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                            <span>{{ app.applicant_phone }}</span>
                          </a>
                        </div>
                      </div>
                    </div>

                    <!-- Card Footer Row -->
                    <div class="card-footer-row">
                      <div class="footer-left-status" :class="getProgressTone(app.status)">
                        <div class="collapsed-next-step-content">
                          <span class="collapsed-next-step-label">ความคืบหน้า</span>
                          <strong>{{ getApplicationProgressText(app) }}</strong>
                        </div>
                      </div>

                      <div class="footer-right-actions">
                        <AppButton v-if="app.has_id_card_copy" variant="outline" size="md" @click.stop="downloadIDCardCopy(app)">สำเนาบัตรประชาชน</AppButton>
                        <AppButton
                          v-if="canTenantEdit(app)"
                          variant="primary"
                          size="md"
                          @click.stop="openRevision(app)"
                          title="แก้ไขคำร้องและส่งเอกสารตรวจสอบอีกครั้ง"
                        >
                          <template #prefix>
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                            </svg>
                          </template>
                          <span>แก้ไขคำร้อง</span>
                        </AppButton>
                        <!-- Detail Button -->
                        <AppButton
                          variant="secondary"
                          size="md"
                          @click.stop="openDetails(app)"
                          title="ดูรายละเอียดคำร้อง"
                        >
                          <template #prefix>
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"></circle><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8z"></path></svg>
                          </template>
                          <span>รายละเอียด</span>
                        </AppButton>

                        <!-- Step 1: Officer Review -->
                        <template v-if="canReview(app) && app.status === 'pending_officer'">
                          <AppButton
                            variant="primary"
                            size="md"
                            @click.stop="processApproval(app, 'เจ้าหน้าที่ตรวจสอบเอกสารครบถ้วน ส่งต่อ ผอ.ตลาด')"
                            title="เจ้าหน้าที่ตรวจสอบและอนุมัติส่งต่อ"
                          >
                            <template #prefix>
                              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                            </template>
                            <span>อนุมัติส่งต่อ</span>
                          </AppButton>
                          <AppButton
                            variant="outline"
                            size="md"
                            @click.stop="openRejectModal(app, 'return_revision')"
                          >
                            <span>ส่งกลับแก้ไข</span>
                          </AppButton>
                          <AppButton
                            variant="danger-outline"
                            size="md"
                            @click.stop="openRejectModal(app, 'reject')"
                          >
                            <span>ไม่อนุมัติ</span>
                          </AppButton>
                        </template>

                        <!-- Step 2: Market Director Review -->
                        <template v-else-if="canReview(app) && app.status === 'pending_market_director'">
                          <AppButton
                            variant="primary"
                            size="md"
                            @click.stop="processApproval(app, 'ผอ.ตลาด พิจารณาอนุมัติคำร้อง')"
                            title="ผอ.ตลาด อนุมัติ"
                          >
                            <template #prefix>
                              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                            </template>
                            <span>ผอ. อนุมัติ</span>
                          </AppButton>
                          <AppButton
                            variant="outline"
                            size="md"
                            @click.stop="openRejectModal(app, 'return_revision')"
                          >
                            <span>ส่งกลับแก้ไข</span>
                          </AppButton>
                          <AppButton
                            variant="danger-outline"
                            size="md"
                            @click.stop="openRejectModal(app, 'reject')"
                          >
                            <span>ไม่อนุมัติ</span>
                          </AppButton>
                        </template>

                        <!-- Step 3: Bureau Director Review -->
                        <template v-else-if="canReview(app) && app.status === 'pending_bureau_director'">
                          <AppButton
                            variant="primary"
                            size="md"
                            @click.stop="processApproval(app, 'ผอ.สำนักตลาด พิจารณาอนุมัติคำร้องเรียบร้อย')"
                            title="ผอ.สำนัก อนุมัติ"
                          >
                            <template #prefix>
                              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                            </template>
                            <span>ผอ.สำนัก อนุมัติ</span>
                          </AppButton>
                          <AppButton
                            variant="outline"
                            size="md"
                            @click.stop="openRejectModal(app, 'return_revision')"
                          >
                            <span>ส่งกลับแก้ไข</span>
                          </AppButton>
                          <AppButton
                            variant="danger-outline"
                            size="md"
                            @click.stop="openRejectModal(app, 'reject')"
                          >
                            <span>ไม่อนุมัติ</span>
                          </AppButton>
                        </template>

                        <!-- Final Step: Approved -> Contract button styled with blue theme -->
                        <template v-else-if="app.status === 'approved' && isStallLease(app)">
                          <NuxtLink
                            :to="'/contracts?new_for_stall=' + encodeURIComponent(app.stall_number || app.stall_id || '') + '&app_id=' + app.id"
                            class="btn-action-primary-solid"
                            title="จัดทำสัญญาเช่าแผงค้า"
                          >
                            <span>ดูคำร้อง / ทำสัญญา</span>
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg>
                          </NuxtLink>
                        </template>
                        <template v-else-if="app.status === 'contract_pending'">
                          <NuxtLink
                            to="/contracts?status=pending"
                            class="btn-action-primary-solid"
                            style="background: #2563eb;"
                            title="ดูสัญญาเช่าที่รอตรวจสอบ"
                          >
                            <span>ดูสัญญา (รอตรวจ)</span>
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg>
                          </NuxtLink>
                        </template>
                        <template v-else-if="app.status === 'completed'">
                          <NuxtLink
                            :to="'/contracts?stall=' + encodeURIComponent(app.stall_number || app.stall_id || '')"
                            class="btn-action-primary-solid"
                            style="background: #10b981;"
                            title="ดูสัญญาเช่า"
                          >
                            <span>ดูสัญญาเช่า</span>
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg>
                          </NuxtLink>
                        </template>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty State in Card View -->
          <div v-else class="empty-state-box officer-empty-state">
            <div class="empty-state-panel">
              <div class="empty-icon" aria-hidden="true">
                <svg width="44" height="44" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                  <polyline points="14 2 14 8 20 8"></polyline>
                  <line x1="16" y1="13" x2="8" y2="13"></line>
                  <line x1="16" y1="17" x2="8" y2="17"></line>
                </svg>
              </div>
              <div class="empty-title">ไม่พบรายการคำร้องที่ค้นหา</div>
              <div class="empty-desc">ลองตรวจสอบคำค้นหา หรือเปลี่ยนตัวกรองเพื่อดูรายการอื่น</div>
              <button
                v-if="searchQuery || activeTab !== 'all' || typeFilter !== 'all'"
                type="button"
                class="btn btn-outline-primary btn-sm"
                style="margin-top: 14px;"
                @click="searchQuery = ''; activeTab = 'all'; typeFilter = 'all'; pendingSubFilter = 'all'; currentPage = 1"
              >
                ล้างตัวกรองทั้งหมด
              </button>
            </div>
          </div>
        </div>

        <!-- View 2: Table View -->
        <div v-else-if="viewMode === 'table'" class="table-container">
          <table class="custom-table">
            <thead>
              <tr>
                <th style="width: 170px;">เลขที่คำร้อง</th>
                <th style="width: 155px;">แผง / เรื่องที่ขอ</th>
                <th>ผู้ยื่นคำร้อง</th>
                <th>ชื่อร้าน / รายละเอียด</th>
                <th style="width: 120px;">วันที่ยื่น / เริ่ม</th>
                <th style="width: 165px;">สถานะพิจารณา</th>
                <th style="min-width: 220px; max-width: 290px;">หมายเหตุการพิจารณา</th>
                <th style="text-align: right; width: 175px;">การดำเนินการ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="app in paginatedApplications" :key="app.id">
                <td>
                  <div style="display: flex; align-items: center; gap: 6px; flex-wrap: wrap;">
                    <strong class="table-code-blue">{{ app.application_number }}</strong>
                    <span v-if="isStallLease(app)" class="badge-type badge-type-stall">ขอเช่าแผง</span>
                    <span v-else class="badge-type badge-type-general">คำร้องทั่วไป</span>
                  </div>
                  <div style="font-size: 0.75rem; color: #94a3b8; margin-top: 2px;">
                    {{ app.created_at?.slice(0, 10) }}
                  </div>
                </td>
                <td>
                  <template v-if="isStallLease(app)">
                    <span class="badge badge-occupied">
                      แผง {{ app.stall_number }}
                    </span>
                    <div style="font-size: 0.75rem; color: #64748b; margin-top: 2px;">
                      ฿{{ Number(app.stall_monthly_rate).toLocaleString() }}/ด. ({{ app.requested_period_months }} เดือน)
                    </div>
                  </template>
                  <template v-else>
                    <div class="general-subject-title">{{ app.subject || app.shop_name || 'คำร้องทั่วไป' }}</div>
                    <div v-if="app.stall_number && app.stall_number !== '-'" class="general-stall-hint">
                      แผง {{ app.stall_number }}
                    </div>
                    <div v-else class="general-no-stall-hint">
                      คำร้องทั่วไป
                    </div>
                  </template>
                </td>
                <td>
                  <div style="font-weight: 600; color: #1e293b;">{{ app.applicant_name }}</div>
                  <div style="font-size: 0.8125rem; color: #64748b;">{{ app.applicant_phone }}</div>
                </td>
                <td>
                  <template v-if="isStallLease(app)">
                    <div style="font-weight: 600; color: #1e293b;">{{ app.shop_name }}</div>
                    <div style="font-size: 0.8125rem; color: #94a3b8;">{{ app.product_type || '-' }}</div>
                  </template>
                  <template v-else>
                    <div style="font-weight: 600; color: #1e293b;">{{ app.shop_name || app.applicant_name }}</div>
                    <div style="font-size: 0.8125rem; color: #64748b; max-width: 220px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;" :title="app.review_note || app.subject">
                      {{ app.review_note || app.subject || '-' }}
                    </div>
                  </template>
                </td>
                <td style="color: #475569;">{{ app.requested_start_date || app.created_at?.slice(0, 10) }}</td>
                <td>
                  <span class="status-pill-badge" :class="getStatusPillClass(app.status)">
                    {{ getStatusBadgeLabel(app.status) }}
                  </span>
                  <div class="action-actor-tag" :class="getActionActorInfo(app).tagClass">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                      <circle cx="12" cy="7" r="4"></circle>
                    </svg>
                    <span>{{ getActionActorInfo(app).actionText }}</span>
                  </div>
                </td>
                <td>
                  <div v-if="getReviewNoteInfo(app).text" class="review-comment-container">
                    <div class="review-comment-pill" :class="'comment-' + getReviewNoteInfo(app).type">
                      <div class="comment-header-row">
                        <span class="comment-role-label">{{ getReviewNoteInfo(app).role }}:</span>
                      </div>
                      <div class="comment-text-body">
                        {{ getReviewNoteInfo(app).text }}
                      </div>
                    </div>
                    <!-- Badge hint if earlier review comments exist -->
                    <div
                      v-if="getAllCommentsList(app).length > 1"
                      class="more-comments-badge"
                      @click="openDetails(app)"
                      title="คลิกเพื่อดูบันทึกความเห็นครบทุกขั้นตอน"
                    >
                      <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="10"></circle>
                        <line x1="12" y1="16" x2="12" y2="12"></line>
                        <line x1="12" y1="8" x2="12.01" y2="8"></line>
                      </svg>
                      <span>บันทึกครบ {{ getAllCommentsList(app).length }} ขั้นตอน</span>
                    </div>
                  </div>
                  <div v-else class="text-empty-dash">-</div>
                </td>
                <td style="text-align: right;">
                  <div style="display: inline-flex; align-items: center; justify-content: flex-end; gap: 6px;">
                    <!-- Step 1: Officer Review -->
                    <button
                      v-if="canReview(app) && app.status === 'pending_officer'"
                      class="btn-table-action btn-table-primary"
                      @click="processApproval(app, 'เจ้าหน้าที่ตรวจสอบเอกสารครบถ้วน ส่งต่อ ผอ.ตลาด')"
                      title="เจ้าหน้าที่ตรวจสอบ"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      <span>ตรวจแล้ว</span>
                    </button>

                    <!-- Step 2: Market Director Review -->
                    <button
                      v-else-if="canReview(app) && app.status === 'pending_market_director'"
                      class="btn-table-action btn-table-primary"
                      @click="processApproval(app, 'ผอ.ตลาด พิจารณาอนุมัติคำร้อง')"
                      title="ผอ.ตลาด อนุมัติ"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      <span>ผอ. อนุมัติ</span>
                    </button>

                    <!-- Step 3: Bureau Director Review -->
                    <button
                      v-else-if="canReview(app) && app.status === 'pending_bureau_director'"
                      class="btn-table-action btn-table-primary"
                      @click="processApproval(app, 'ผอ.สำนักตลาด พิจารณาอนุมัติคำร้องเรียบร้อย')"
                      title="ผอ.สำนัก อนุมัติ"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      <span>ผอ.สำนัก</span>
                    </button>

                    <!-- Final: Approved -> Generate Contract (Only for stall lease) -->
                    <NuxtLink
                      v-else-if="app.status === 'approved' && isStallLease(app)"
                      :to="'/contracts?new_for_stall=' + encodeURIComponent(app.stall_number || app.stall_id || '') + '&app_id=' + app.id"
                      class="btn-table-action btn-table-primary"
                      title="จัดทำสัญญาเช่า"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
                      <span>ทำสัญญา</span>
                    </NuxtLink>
                    <NuxtLink
                      v-else-if="app.status === 'contract_pending'"
                      to="/contracts?status=pending"
                      class="btn-table-action btn-table-primary"
                      style="background: #2563eb; color: #fff; border-color: #2563eb;"
                      title="ดูสัญญาเช่าที่รอเจ้าหน้าที่ตรวจสอบ"
                    >
                      <span>ดูสัญญา (รอตรวจ)</span>
                    </NuxtLink>
                    <NuxtLink
                      v-else-if="app.status === 'completed'"
                      :to="'/contracts?stall=' + encodeURIComponent(app.stall_number || app.stall_id || '')"
                      class="btn-table-action btn-table-view"
                      title="ดูสัญญาเช่า"
                    >
                      <span>ดูสัญญา</span>
                    </NuxtLink>
                    <span
                      v-else-if="app.status === 'approved' && !isStallLease(app)"
                      class="general-approved-pill"
                    >
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      <span>อนุมัติแล้ว</span>
                    </span>

                    <!-- Officer & Tenant actions -->
                    <AppButton
                      v-if="canTenantEdit(app)"
                      variant="primary"
                      size="sm"
                      @click="openRevision(app)"
                      title="แก้ไขคำร้องและส่งเอกสารตรวจสอบอีกครั้ง"
                    >
                      <template #prefix>
                        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                        </svg>
                      </template>
                      <span>แก้ไขคำร้อง</span>
                    </AppButton>
                    <AppButton v-if="app.has_id_card_copy" variant="outline" size="sm" @click="downloadIDCardCopy(app)">สำเนาบัตรประชาชน</AppButton>
                    <template v-if="canReview(app) && app.status === 'pending_officer'">
                      <button class="btn-table-action btn-table-primary" @click="processApproval(app, 'เจ้าหน้าที่ตรวจสอบเอกสารครบถ้วน ส่งต่อ ผอ.ตลาด')">อนุมัติส่งต่อ</button>
                      <button class="btn-table-action btn-table-outline" @click="openRejectModal(app, 'return_revision')">ส่งกลับแก้ไข</button>
                      <button class="btn-table-action btn-table-cancel" @click="openRejectModal(app, 'reject')">ไม่อนุมัติ</button>
                    </template>
                    <template v-else-if="canReview(app) && (app.status === 'pending_market_director' || app.status === 'pending_bureau_director')">
                      <button class="btn-table-action btn-table-primary" @click="processApproval(app, 'พิจารณาอนุมัติคำร้อง')">อนุมัติ</button>
                      <button class="btn-table-action btn-table-outline" @click="openRejectModal(app, 'return_revision')">ส่งกลับแก้ไข</button>
                      <button class="btn-table-action btn-table-cancel" @click="openRejectModal(app, 'reject')">ไม่อนุมัติ</button>
                    </template>

                    <!-- Detail Button -->
                    <button
                      class="btn-table-icon-edit"
                      title="ดูรายละเอียด/จัดการคำร้อง"
                      @click="openDetails(app)"
                    >
                      <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8z"></path></svg>
                    </button>
                  </div>
                </td>
              </tr>

              <tr v-if="filteredApplications.length === 0">
                <td colspan="8" style="text-align: center; color: #94a3b8; padding: 48px;">
                  ไม่พบรายการคำร้องขอเช่าตามเงื่อนไขที่ค้นหา
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Standard AppPagination Component -->
        <AppPagination
          v-if="filteredApplications.length > 0"
          v-model="currentPage"
          v-model:items-per-page="pageSize"
          :total-items="filteredApplications.length"
          unit-label="คำร้อง"
        />
      </div>
    </div>

    <!-- Dedicated Reject Modal with Remark / Reason Input -->
    <div
      v-if="showRejectModal && targetAppForReject"
      class="modal-overlay"
      @click.self="closeRejectModal"
    >
      <div class="modal-content reject-modal-box">
        <!-- Modal Header -->
        <div class="modal-header">
          <div class="modal-header-title-wrap">
            <div class="modal-icon-badge modal-icon-danger">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="15" y1="9" x2="9" y2="15"></line>
                <line x1="9" y1="9" x2="15" y2="15"></line>
              </svg>
            </div>
            <div>
              <h3 class="modal-title-danger">{{ rejectAction === 'return_revision' ? 'ระบุเหตุผลส่งกลับแก้ไข' : 'ระบุเหตุผลไม่อนุมัติคำร้อง' }}</h3>
              <div class="modal-subtitle">
                คำร้องเลขที่ <strong style="color: #2563eb;">{{ targetAppForReject.application_number }}</strong>
              </div>
            </div>
          </div>
          <button class="btn-modal-close" @click="closeRejectModal" aria-label="ปิด">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Modal Body -->
        <div class="modal-body">
          <!-- Summary Banner -->
          <div class="reject-summary-banner">
            <div class="summary-item">
              <span class="summary-label">ผู้ยื่นคำร้อง:</span>
              <strong class="summary-value">{{ targetAppForReject.applicant_name }}</strong>
            </div>
            <div class="summary-item">
              <span class="summary-label">ชื่อร้าน / แผง:</span>
              <span class="summary-value">{{ targetAppForReject.shop_name }} (แผง {{ targetAppForReject.stall_number || targetAppForReject.stall_ident }})</span>
            </div>
          </div>

          <!-- Alert Note -->
          <div class="reject-alert-box">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink: 0; margin-top: 1px;">
              <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
              <line x1="12" y1="9" x2="12" y2="13"></line>
              <line x1="12" y1="17" x2="12.01" y2="17"></line>
            </svg>
            <span>{{ rejectAction === 'return_revision' ? 'คำร้องจะถูกส่งกลับให้ผู้เช่าแก้ไข พร้อมแสดงหมายเหตุ โดยยังคงการจองไว้' : 'คำร้องจะถูกยกเลิกและบันทึกเหตุผลการไม่อนุมัติ' }}</span>
          </div>

          <!-- Quick Reason Chips -->
          <div style="margin-bottom: 12px;">
            <label class="form-label-sub">เลือกข้อความสำเร็จรูป (คลิกเพื่อเพิ่มลงในช่องหมายเหตุ):</label>
            <div class="quick-preset-chips">
              <button
                v-for="(r, idx) in quickRejectReasons"
                :key="idx"
                type="button"
                class="quick-preset-chip"
                @click="setQuickReason(r)"
              >
                + {{ r }}
              </button>
            </div>
          </div>

          <!-- Textarea Field -->
          <div>
            <label class="form-label" for="reject-comment-input">
              หมายเหตุ / เหตุผล <span style="color: #dc2626;">*</span>
            </label>
            <textarea
              id="reject-comment-input"
              v-model="rejectComment"
              class="form-textarea"
              rows="4"
              placeholder="ระบุเหตุผลหรือรายละเอียดเพิ่มเติมที่ผู้ยื่นคำร้องควรทราบ..."
              :class="{ 'input-error': rejectErrorMsg }"
            ></textarea>
            <div v-if="rejectErrorMsg" class="form-error-msg">
              {{ rejectErrorMsg }}
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="modal-footer">
          <AppButton
            variant="secondary"
            size="md"
            @click="closeRejectModal"
          >
            ยกเลิก
          </AppButton>
          <AppButton
            variant="danger"
            size="md"
            :loading="rejectSubmitting"
            @click="submitReject"
          >
            {{ rejectAction === 'return_revision' ? 'ยืนยันส่งกลับแก้ไข' : 'ยืนยันไม่อนุมัติคำร้อง' }}
          </AppButton>
        </div>
      </div>
    </div>

    <!-- Modal แก้ไขคำร้องขอเช่าแผงค้า (Revision Modal) -->
    <div v-if="revisionTarget" class="modal-overlay" @click.self="closeRevision">
      <div class="modal-content revision-modal-content" role="dialog" aria-modal="true" aria-labelledby="revision-title">
        <!-- Modal Header -->
        <div class="modal-header">
          <div class="revision-header-left">
            <div class="revision-icon-badge">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
              </svg>
            </div>
            <div>
              <h3 id="revision-title" class="revision-modal-title">แก้ไขคำร้องขอเช่าแผงค้า (Edit Lease Application)</h3>
              <div class="revision-modal-subtitle">
                แผงเลขที่: <strong>{{ revisionTargetStallNumber }}</strong>
                <span class="revision-sub-dot">•</span>
                <span>เลขที่คำร้อง: <strong>{{ revisionTarget.application_number || revisionTarget.application_no || '-' }}</strong></span>
              </div>
            </div>
          </div>
          <button type="button" class="btn-modal-close" aria-label="ปิด" @click="closeRevision">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Form Body -->
        <form novalidate @submit.prevent="submitRevision">
          <div class="modal-body revision-modal-body">
            <!-- Review Note Alert (ถ้าเจ้าหน้าที่ส่งกลับแก้ไข) -->
            <div
              v-if="revisionTarget.review_note || revisionTarget.officer_comment"
              class="revision-officer-note-card"
            >
              <div class="revision-note-header">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="8" x2="12" y2="12"></line>
                  <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <span>เหตุผลที่เจ้าหน้าที่ส่งกลับแก้ไข / ข้อเสนอแนะ</span>
              </div>
              <div class="revision-note-body">
                {{ revisionTarget.review_note || revisionTarget.officer_comment }}
              </div>
            </div>

            <!-- Financial Summary Banner (เหมือนตอนยื่นคำร้องครั้งแรกในหน้าจอง) -->
            <div class="revision-stall-summary-banner">
              <div class="summary-banner-inner">
                <div class="summary-stall-col">
                  <span>แผงค้าที่ขอเช่า: </span>
                  <strong>{{ revisionTargetStallNumber }}</strong>
                  <span v-if="revisionTargetMarketName" class="summary-market-tag">({{ revisionTargetMarketName }})</span>
                </div>
                <div class="summary-rates-wrap">
                  <span>ค่าเช่า: <strong>฿{{ revisionMonthlyRent.toLocaleString() }}/เดือน</strong></span>
                  <span>เงินมัดจำ: <strong>฿{{ revisionDeposit.toLocaleString() }}</strong></span>
                  <span>รวมแรกเข้า: <strong class="summary-total-highlight">฿{{ (revisionMonthlyRent + revisionDeposit).toLocaleString() }}</strong></span>
                </div>
              </div>
            </div>

            <!-- ส่วนที่ 1: ข้อมูลผู้ยื่นคำร้อง (Applicant Details) -->
            <h4 class="form-section-heading">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
              <span>1. ข้อมูลผู้ยื่นคำร้อง (Applicant Details)</span>
            </h4>

            <!-- ลำดับที่ 1: ชื่อ-นามสกุล -> ลำดับที่ 2: อายุ -->
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="rev_applicant_name">ชื่อ-นามสกุล <b class="req-star">*</b></label>
                <AppInput
                  id="rev_applicant_name"
                  v-model="revisionForm.applicant_name"
                  type="text"
                  placeholder="ชื่อ-นามสกุล ผู้ยื่นคำร้อง"
                  :error="revisionErrors.applicant_name"
                  @blur="validateRevisionField('applicant_name')"
                />
              </div>

              <div class="form-group">
                <label class="form-label" for="rev_applicant_age">อายุ (ปี)</label>
                <AppInput
                  id="rev_applicant_age"
                  v-model="revisionForm.applicant_age"
                  type="number"
                  min="1"
                  max="120"
                  placeholder="ระบุอายุ เช่น 35"
                />
              </div>
            </div>

            <!-- ลำดับที่ 3: เลขบัตรประชาชน -> ลำดับที่ 4: เบอร์โทร -->
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="rev_applicant_id">เลขประจำตัวประชาชน (13 หลัก) <b class="req-star">*</b></label>
                <AppInput
                  id="rev_applicant_id"
                  v-model="revisionForm.applicant_id_card"
                  type="text"
                  placeholder="x-xxxx-xxxxx-xx-x"
                  maxlength="17"
                  :error="revisionErrors.applicant_id_card"
                  @blur="validateRevisionField('applicant_id_card')"
                />
              </div>

              <div class="form-group">
                <label class="form-label" for="rev_applicant_phone">เบอร์โทรศัพท์ติดต่อ <b class="req-star">*</b></label>
                <AppInput
                  id="rev_applicant_phone"
                  v-model="revisionForm.applicant_phone"
                  type="tel"
                  placeholder="08x-xxx-xxxx"
                  :error="revisionErrors.applicant_phone"
                  @blur="validateRevisionField('applicant_phone')"
                />
              </div>
            </div>

            <!-- ลำดับที่ 5: อีเมล -> ที่อยู่ตามบัตรประชาชน -->
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="rev_applicant_email">อีเมล</label>
                <AppInput
                  id="rev_applicant_email"
                  v-model="revisionForm.applicant_email"
                  type="email"
                  placeholder="example@domain.com"
                  :error="revisionErrors.applicant_email"
                  @blur="validateRevisionField('applicant_email')"
                />
              </div>

              <div class="form-group">
                <label class="form-label" for="rev_applicant_address">ที่อยู่ตามบัตรประชาชน / ที่อยู่ติดต่อ</label>
                <AppInput
                  id="rev_applicant_address"
                  v-model="revisionForm.applicant_address"
                  type="text"
                  placeholder="บ้านเลขที่ ถนน แขวง/ตำบล เขต/อำเภอ จังหวัด"
                />
              </div>
            </div>

            <!-- ลำดับที่ 6: สำเนาบัตรประชาชน (Modern Document Attachment UI) -->
            <div class="form-group upload-section-container">
              <div class="upload-header-bar">
                <label class="upload-label" for="rev_id_copy">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <rect x="3" y="4" width="18" height="16" rx="2"></rect>
                    <circle cx="8" cy="10" r="2"></circle>
                    <line x1="14" y1="9" x2="18" y2="9"></line>
                    <line x1="14" y1="13" x2="18" y2="13"></line>
                    <line x1="6" y1="16" x2="18" y2="16"></line>
                  </svg>
                  <span>สำเนาบัตรประชาชน (ID Card Copy)</span>
                </label>
                <span class="upload-badge-spec">PDF, JPG หรือ PNG (ไม่เกิน 3 MB)</span>
              </div>

              <input
                ref="revIdCardInputRef"
                id="rev_id_copy"
                type="file"
                class="sr-only"
                accept=".pdf,.jpg,.jpeg,.png,application/pdf,image/jpeg,image/png"
                @change="onRevIDCardCopyChange"
              />

              <!-- การ์ดพรีวิวไฟล์เมื่อเลือกเอกสารแล้ว -->
              <div v-if="revIdCardCopy" class="id-card-preview-box">
                <div class="preview-thumb-wrap">
                  <img v-if="revIdCardPreviewUrl" :src="revIdCardPreviewUrl" alt="ID Card Preview" class="preview-thumb-img" />
                  <div v-else class="preview-pdf-badge">
                    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                      <polyline points="14 2 14 8 20 8"></polyline>
                      <line x1="16" y1="13" x2="8" y2="13"></line>
                      <line x1="16" y1="17" x2="8" y2="17"></line>
                    </svg>
                    <span>PDF</span>
                  </div>
                </div>

                <div class="preview-details">
                  <div class="preview-filename" :title="revIdCardCopy.name">{{ revIdCardCopy.name }}</div>
                  <div class="preview-meta">
                    <span class="meta-size">{{ formatFileSize(revIdCardCopy.size) }}</span>
                    <span class="meta-check">
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                      <span>พร้อมแนบไฟล์</span>
                    </span>
                  </div>
                </div>

                <div class="preview-actions">
                  <button type="button" class="btn-preview-change" @click="triggerRevFileInput" title="เปลี่ยนไฟล์">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M21 2v6h-6"></path>
                      <path d="M3 12a9 9 0 0 1 15-6.7L21 8"></path>
                      <path d="M3 22v-6h6"></path>
                      <path d="M21 12a9 9 0 0 1-15 6.7L3 16"></path>
                    </svg>
                    <span>เปลี่ยน</span>
                  </button>
                  <button type="button" class="btn-preview-delete" @click="removeRevIDCardFile" title="ลบไฟล์">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="3 6 5 6 21 6"></polyline>
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                    </svg>
                  </button>
                </div>
              </div>

              <!-- กล่อง Dropzone แบบคลีนเมื่อยังไม่ได้เลือกไฟล์ใหม่ -->
              <div
                v-else
                class="id-card-dropzone"
                :class="{ 'is-drag-over': isRevDragging, 'has-error': !!revIdCardCopyError }"
                @click="triggerRevFileInput"
                @dragover.prevent="isRevDragging = true"
                @dragleave.prevent="isRevDragging = false"
                @drop.prevent="handleRevFileDrop"
              >
                <div class="dropzone-icon-box">
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                    <polyline points="17 8 12 3 7 8"></polyline>
                    <line x1="12" y1="3" x2="12" y2="15"></line>
                  </svg>
                </div>
                <div class="dropzone-prompt">
                  <div class="prompt-title">
                    <span class="prompt-link">คลิกเพื่อเลือกไฟล์</span> หรือลากไฟล์มาวางที่นี่
                  </div>
                  <div class="prompt-desc">แนบสำเนาบัตรประจำตัวประชาชนฉบับแก้ไข พร้อมรับรองสำเนาถูกต้อง</div>
                </div>
                <div class="dropzone-cta-btn">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
                  </svg>
                  <span>เลือกไฟล์</span>
                </div>
              </div>

              <div v-if="revIdCardCopyError" class="upload-error-msg">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="8" x2="12" y2="12"></line>
                  <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <span>{{ revIdCardCopyError }}</span>
              </div>
            </div>

            <!-- ส่วนที่ 2: ข้อมูลร้านค้าและสินค้าที่ขอจำหน่าย -->
            <h4 class="form-section-heading">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path>
                <line x1="3" y1="6" x2="21" y2="6"></line>
                <path d="M16 10a4 4 0 0 1-8 0"></path>
              </svg>
              <span>2. ข้อมูลร้านค้าและสินค้าที่ขอจำหน่าย</span>
            </h4>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="rev_shop_name">ชื่อร้านค้า (Shop Name) <b class="req-star">*</b></label>
                <AppInput
                  id="rev_shop_name"
                  v-model="revisionForm.shop_name"
                  type="text"
                  placeholder="เช่น ชานมไข่มุก เฮงเฮง"
                  :error="revisionErrors.shop_name"
                  @blur="validateRevisionField('shop_name')"
                />
              </div>

              <div class="form-group">
                <label class="form-label" for="rev_product_category">ประเภทสินค้าที่จะจำหน่าย <b class="req-star">*</b></label>
                <AppInput
                  id="rev_product_category"
                  v-model="revisionForm.product_category"
                  type="text"
                  placeholder="เช่น เครื่องดื่ม, ของทอด, เสื้อผ้า"
                  :error="revisionErrors.product_category"
                  @blur="validateRevisionField('product_category')"
                />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="rev_requested_start_date">วันที่ประสงค์เริ่มเช่า <b class="req-star">*</b></label>
                <AppInput
                  id="rev_requested_start_date"
                  v-model="revisionForm.requested_start_date"
                  type="date"
                  :error="revisionErrors.requested_start_date"
                  @blur="validateRevisionField('requested_start_date')"
                />
              </div>

              <div class="form-group">
                <label class="form-label">ระยะเวลาขอเช่า (เดือน) <b class="req-star">*</b></label>
                <AppSelect
                  v-model="revisionForm.lease_term_months"
                  :options="revisionPeriodOptions"
                  placeholder="เลือกระยะเวลาขอเช่า"
                />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="rev_details">รายละเอียดเพิ่มเติม / หมายเหตุการแก้ไข</label>
              <textarea
                id="rev_details"
                v-model.trim="revisionForm.details"
                class="form-textarea-revision"
                rows="2"
                placeholder="ระบุรายละเอียดเพิ่มเติม หรือสิ่งที่ได้ปรับปรุงแก้ไขตามที่เจ้าหน้าที่แจ้ง"
              ></textarea>
            </div>

            <div v-if="revisionError" class="revision-form-error-alert">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              <span>{{ revisionError }}</span>
            </div>
          </div>

          <!-- Modal Footer: ยกเลิก และ ส่งกลับให้เจ้าหน้าที่ตรวจสอบอีกครั้ง -->
          <div class="modal-footer revision-modal-footer">
            <AppButton
              type="button"
              variant="secondary"
              size="md"
              :disabled="revisionSaving"
              @click="closeRevision"
            >
              ยกเลิก
            </AppButton>
            <AppButton
              type="submit"
              variant="primary"
              size="md"
              :loading="revisionSaving"
              :disabled="revisionSaving"
            >
              <template #prefix>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="22" y1="2" x2="11" y2="13"></line>
                  <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
                </svg>
              </template>
              <span>ส่งกลับให้เจ้าหน้าที่ตรวจสอบอีกครั้ง</span>
            </AppButton>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal ยื่นคำร้องทั่วไป -->
    <div
      v-if="showGeneralModal"
      class="modal-overlay"
      @click.self="closeCreateGeneralModal"
    >
      <div class="modal-content general-modal-box">
        <!-- Modal Header -->
        <div class="modal-header">
          <div class="modal-header-title-wrap">
            <div class="modal-icon-badge modal-icon-purple">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="16" y1="13" x2="8" y2="13"></line>
                <line x1="16" y1="17" x2="8" y2="17"></line>
              </svg>
            </div>
            <div>
              <h3 class="modal-title-general">ยื่นคำร้องทั่วไป</h3>
              <div class="modal-subtitle">บันทึกคำร้องทั่วไปสำหรับผู้ค้าหรือประชาชนทั่วไป</div>
            </div>
          </div>
          <button class="btn-modal-close" @click="closeCreateGeneralModal" aria-label="ปิด">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Modal Body -->
        <div class="modal-body general-modal-body">
          <form @submit.prevent="submitGeneralApplication" class="general-app-form">
            <div class="form-grid-2">
              <div class="form-group">
                <label class="form-label" for="gen-applicant-name">ชื่อ-นามสกุล ผู้ยื่นคำร้อง <span class="text-danger">*</span></label>
                <input
                  id="gen-applicant-name"
                  v-model="generalForm.applicant_name"
                  type="text"
                  class="form-input"
                  placeholder="เช่น นายสมชาย ใจดี"
                  required
                />
              </div>

              <div class="form-group">
                <label class="form-label" for="gen-applicant-phone">เบอร์โทรศัพท์ติดต่อ <span class="text-danger">*</span></label>
                <input
                  id="gen-applicant-phone"
                  v-model="generalForm.applicant_phone"
                  type="tel"
                  class="form-input"
                  placeholder="เช่น 0812345678"
                  required
                />
              </div>
            </div>

            <div class="form-grid-2">
              <div class="form-group">
                <label class="form-label" for="gen-applicant-id-card">เลขประจำตัวประชาชน</label>
                <input
                  id="gen-applicant-id-card"
                  v-model="generalForm.applicant_id_card"
                  type="text"
                  class="form-input"
                  placeholder="13 หลัก (ถ้ามี)"
                  maxlength="13"
                />
              </div>

              <div class="form-group">
                <label class="form-label" for="gen-applicant-email">อีเมล</label>
                <input
                  id="gen-applicant-email"
                  v-model="generalForm.applicant_email"
                  type="email"
                  class="form-input"
                  placeholder="name@example.com"
                />
              </div>
            </div>

            <div class="form-grid-2">
              <div class="form-group">
                <label class="form-label">ตลาดที่เกี่ยวข้อง</label>
                <AppSelect
                  v-model="generalForm.market_id"
                  :options="marketOptions"
                  placeholder="เลือกตลาดที่เกี่ยวข้อง"
                />
              </div>

              <div class="form-group">
                <label class="form-label" for="gen-stall-number">รหัสแผงค้า (ถ้ามี)</label>
                <input
                  id="gen-stall-number"
                  v-model="generalForm.stall_number"
                  type="text"
                  class="form-input"
                  placeholder="เช่น A-01 (เว้นว่างได้หากไม่มี)"
                />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="gen-subject">เรื่อง / หัวข้อคำร้อง <span class="text-danger">*</span></label>
              <input
                id="gen-subject"
                v-model="generalForm.subject"
                type="text"
                class="form-input"
                placeholder="เช่น ขอเปลี่ยนชื่อผู้ประกอบการ, ขอขยายเวลาชำระค่าเช่า, ขอซ่อมแซมแผง"
                required
              />
            </div>

            <div class="form-group">
              <label class="form-label" for="gen-details">รายละเอียดคำร้อง <span class="text-danger">*</span></label>
              <textarea
                id="gen-details"
                v-model="generalForm.details"
                class="form-textarea"
                rows="4"
                placeholder="ระบุรายละเอียด ความประสงค์ หรือเหตุผลความจำเป็นประกอบคำร้อง..."
                required
              ></textarea>
            </div>

            <div class="form-group">
              <label class="form-label" for="gen-applicant-address">ที่อยู่ผู้ยื่นคำร้อง</label>
              <input
                id="gen-applicant-address"
                v-model="generalForm.applicant_address"
                type="text"
                class="form-input"
                placeholder="ที่อยู่สำหรับการติดต่อหรือส่งเอกสาร"
              />
            </div>

            <div v-if="generalFormError" class="form-error-alert">
              {{ generalFormError }}
            </div>
          </form>
        </div>

        <!-- Modal Footer -->
        <div class="modal-footer">
          <AppButton
            variant="secondary"
            size="md"
            @click="closeCreateGeneralModal"
          >
            ยกเลิก
          </AppButton>
          <AppButton
            variant="primary"
            size="md"
            :loading="generalSubmitting"
            @click="submitGeneralApplication"
          >
            บันทึกและยื่นคำร้อง
          </AppButton>
        </div>
      </div>
    </div>

    <!-- Slide-over Right Side Panel for History -->
    <AppHistoryPanel
      :is-open="showHistoryPanel"
      title="ประวัติคำร้อง"
      history-type="application"
      :history-data="applicationHistory"
      @close="showHistoryPanel = false"
    />
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'
import { useNotifications } from '~/composables/useNotifications'

const { get, post, apiFetch } = useApi()
const route = useRoute()
const authStore = useAuthStore()
const { user, currentUser, isTenant, role } = useAuth()
const isAdminOrStaff = computed(() => ['admin', 'staff'].includes(role.value || authStore.role || ''))
const canReview = (app: any) => {
  const currentRole = role.value || authStore.role || ''
  return (app.status === 'pending_officer' && ['staff', 'admin'].includes(currentRole)) ||
    (app.status === 'pending_market_director' && currentRole === 'market_director') ||
    (app.status === 'pending_bureau_director' && currentRole === 'bureau_director')
}
const { showAlert, showConfirm, showSuccess, showError } = useAppAlert()
const { fetchNotifications } = useNotifications()
const { resetTestData } = useResetTestData()

const revisionTarget = ref<any>(null)
const revisionSaving = ref(false)
const revisionError = ref('')
const revIdCardCopy = ref<File | null>(null)
const revIdCardCopyError = ref('')
const revIdCardInputRef = ref<HTMLInputElement | null>(null)
const isRevDragging = ref(false)
const revIdCardPreviewUrl = ref('')

const triggerRevFileInput = () => {
  revIdCardInputRef.value?.click()
}

const handleRevFileSelect = (file: File | null) => {
  revIdCardCopyError.value = ''
  if (revIdCardPreviewUrl.value) {
    URL.revokeObjectURL(revIdCardPreviewUrl.value)
    revIdCardPreviewUrl.value = ''
  }
  if (!file) {
    revIdCardCopy.value = null
    return
  }
  if (!['application/pdf', 'image/jpeg', 'image/png'].includes(file.type) || file.size > 3 * 1024 * 1024) {
    revIdCardCopy.value = null
    if (revIdCardInputRef.value) revIdCardInputRef.value.value = ''
    revIdCardCopyError.value = 'กรุณาเลือกไฟล์ PDF, JPG หรือ PNG ขนาดไม่เกิน 3 MB'
    return
  }
  revIdCardCopy.value = file
  if (file.type.startsWith('image/')) {
    revIdCardPreviewUrl.value = URL.createObjectURL(file)
  }
}

const onRevIDCardCopyChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0] || null
  handleRevFileSelect(file)
}

const handleRevFileDrop = (event: DragEvent) => {
  isRevDragging.value = false
  const file = event.dataTransfer?.files?.[0] || null
  handleRevFileSelect(file)
}

const removeRevIDCardFile = () => {
  if (revIdCardPreviewUrl.value) {
    URL.revokeObjectURL(revIdCardPreviewUrl.value)
    revIdCardPreviewUrl.value = ''
  }
  revIdCardCopy.value = null
  revIdCardCopyError.value = ''
  if (revIdCardInputRef.value) revIdCardInputRef.value.value = ''
}

const readRevIDCardCopy = (file: File) => new Promise<string>((resolve, reject) => {
  const reader = new FileReader()
  reader.onerror = () => reject(new Error('ไม่สามารถอ่านไฟล์สำเนาบัตรประชาชน'))
  reader.onload = () => resolve(String(reader.result).split(',')[1] || '')
  reader.readAsDataURL(file)
})

const revisionForm = reactive({
  applicant_name: '',
  applicant_phone: '',
  applicant_id_card: '',
  applicant_email: '',
  applicant_address: '',
  applicant_age: '',
  shop_name: '',
  product_category: '',
  details: '',
  requested_start_date: '',
  lease_term_months: 12
})
const revisionPeriodOptions = [
  { value: 3, label: '3 เดือน (ทดลองขาย)' },
  { value: 6, label: '6 เดือน' },
  { value: 12, label: '1 ปี (มาตรฐาน)' },
  { value: 24, label: '2 ปี' }
]

const revisionErrors = reactive({
  applicant_name: '',
  applicant_id_card: '',
  applicant_phone: '',
  applicant_email: '',
  shop_name: '',
  product_category: '',
  requested_start_date: ''
})

const resetRevisionErrors = () => {
  revisionErrors.applicant_name = ''
  revisionErrors.applicant_id_card = ''
  revisionErrors.applicant_phone = ''
  revisionErrors.applicant_email = ''
  revisionErrors.shop_name = ''
  revisionErrors.product_category = ''
  revisionErrors.requested_start_date = ''
}

const validateRevisionField = (field: keyof typeof revisionErrors) => {
  if (field === 'applicant_name') {
    revisionErrors.applicant_name = revisionForm.applicant_name.trim() ? '' : 'กรุณาระบุชื่อ-นามสกุล'
  } else if (field === 'applicant_id_card') {
    const id = revisionForm.applicant_id_card.trim().replace(/[-\s]/g, '')
    if (!id) {
      revisionErrors.applicant_id_card = 'กรุณาระบุเลขประจำตัวประชาชน'
    } else if (id.length !== 13) {
      revisionErrors.applicant_id_card = 'เลขประจำตัวประชาชนต้องครบ 13 หลัก'
    } else {
      revisionErrors.applicant_id_card = ''
    }
  } else if (field === 'applicant_phone') {
    revisionErrors.applicant_phone = revisionForm.applicant_phone.trim() ? '' : 'กรุณาระบุเบอร์โทรศัพท์ติดต่อ'
  } else if (field === 'applicant_email') {
    if (revisionForm.applicant_email.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(revisionForm.applicant_email.trim())) {
      revisionErrors.applicant_email = 'รูปแบบอีเมลไม่ถูกต้อง'
    } else {
      revisionErrors.applicant_email = ''
    }
  } else if (field === 'shop_name') {
    revisionErrors.shop_name = revisionForm.shop_name.trim() ? '' : 'กรุณาระบุชื่อร้านค้า'
  } else if (field === 'product_category') {
    revisionErrors.product_category = revisionForm.product_category.trim() ? '' : 'กรุณาระบุประเภทสินค้าที่จะจำหน่าย'
  } else if (field === 'requested_start_date') {
    revisionErrors.requested_start_date = revisionForm.requested_start_date ? '' : 'กรุณาเลือกวันที่ประสงค์เริ่มเช่า'
  }
}

const validateRevisionForm = (): boolean => {
  validateRevisionField('applicant_name')
  validateRevisionField('applicant_id_card')
  validateRevisionField('applicant_phone')
  validateRevisionField('applicant_email')
  validateRevisionField('shop_name')
  validateRevisionField('product_category')
  validateRevisionField('requested_start_date')
  return !Object.values(revisionErrors).some(Boolean)
}

const revisionTargetStallNumber = computed(() => {
  if (!revisionTarget.value) return '-'
  return revisionTarget.value.stall_number || revisionTarget.value.stall_ident || revisionTarget.value.stall_id || '-'
})

const revisionTargetMarketName = computed(() => {
  if (!revisionTarget.value) return ''
  return getApplicationProjectName(revisionTarget.value) || revisionTarget.value.market_name || ''
})

const revisionMonthlyRent = computed(() => {
  if (!revisionTarget.value) return 0
  return Number(
    revisionTarget.value.stall_monthly_rate ||
    revisionTarget.value.monthly_rent ||
    revisionTarget.value.monthly_rate ||
    revisionTarget.value.rent_amount ||
    revisionTarget.value.stall?.monthly_rate ||
    revisionTarget.value.stall?.price ||
    0
  )
})

const revisionDeposit = computed(() => {
  if (!revisionTarget.value) return 0
  const dep = Number(
    revisionTarget.value.deposit_amount ||
    revisionTarget.value.deposit ||
    0
  )
  if (dep > 0) return dep
  return revisionMonthlyRent.value > 0 ? revisionMonthlyRent.value * 2 : 0
})

const openRevision = (app: any) => {
  revisionTarget.value = app
  revisionError.value = ''
  resetRevisionErrors()
  removeRevIDCardFile()
  revisionForm.applicant_name = app.applicant_name || ''
  revisionForm.applicant_age = app.applicant_age ? String(app.applicant_age) : ''
  revisionForm.applicant_phone = app.applicant_phone || ''
  revisionForm.applicant_id_card = app.applicant_id_card || ''
  revisionForm.applicant_email = app.applicant_email || ''
  revisionForm.applicant_address = app.applicant_address || ''
  revisionForm.shop_name = app.shop_name || ''
  revisionForm.product_category = app.product_category || app.product_type || ''
  revisionForm.details = app.details || ''
  revisionForm.requested_start_date = app.requested_start_date ? String(app.requested_start_date).slice(0, 10) : ''
  revisionForm.lease_term_months = Number(app.lease_term_months || app.requested_period_months || 12)
}

const closeRevision = () => {
  if (!revisionSaving.value) {
    revisionTarget.value = null
    removeRevIDCardFile()
    resetRevisionErrors()
  }
}

const submitRevision = async () => {
  if (!revisionTarget.value || revisionSaving.value) return
  if (!validateRevisionForm()) return
  if (revIdCardCopyError.value) return

  const age = revisionForm.applicant_age === '' ? null : Number(revisionForm.applicant_age)
  if (age !== null && (!Number.isInteger(age) || age < 1 || age > 120)) {
    await showError('กรุณาระบุอายุระหว่าง 1 ถึง 120 ปี', 'ข้อมูลไม่ถูกต้อง')
    return
  }

  revisionError.value = ''
  revisionSaving.value = true
  try {
    const idCardCopyBase64 = revIdCardCopy.value ? await readRevIDCardCopy(revIdCardCopy.value) : ''
    const payload = {
      ...revisionForm,
      applicant_age: age,
      id_card_copy_base64: idCardCopyBase64,
      product_type: revisionForm.product_category,
      lease_term_months: Number(revisionForm.lease_term_months) || 12,
      requested_period_months: Number(revisionForm.lease_term_months) || 12
    }
    await post(`/applications/${revisionTarget.value.id}/resubmit`, payload)
    const appNo = revisionTarget.value.application_number || revisionTarget.value.application_no || ''
    revisionTarget.value = null
    removeRevIDCardFile()
    await loadApplications()
    await showSuccess(`ส่งกลับเอกสารคำร้อง ${appNo} ให้เจ้าหน้าที่ตรวจสอบเรียบร้อยแล้ว`, 'ส่งคำร้องสำเร็จ')
  } catch (err: any) {
    revisionError.value = err?.data?.error?.message || err?.data?.message || 'ไม่สามารถส่งคำร้องได้ กรุณาลองใหม่อีกครั้ง'
  } finally {
    revisionSaving.value = false
  }
}

const handleClearAllData = async () => {
  await resetTestData(async () => {
    await loadApplications()
  })
}

const loading = ref(true)
const applications = ref<any[]>([])
const marketsList = ref<any[]>([])
const searchQuery = ref(typeof route.query.search === 'string' ? route.query.search : '')

watch(() => route.query.search, (val) => {
  if (typeof val === 'string') {
    searchQuery.value = val
  }
})
const activeTab = ref<'all' | 'pending' | 'approved' | 'rejected'>('all')
const pendingSubFilter = ref<'all' | 'pending_officer' | 'pending_debt' | 'pending_market_director' | 'pending_bureau_director'>('all')
const showPendingTooltip = ref(false)
const typeFilter = ref<'all' | 'stall_lease' | 'general'>('all')
const currentPage = ref(1)
const pageSize = ref(10)

const viewMode = ref<'card' | 'table'>('card')
const viewModeOptions = [
  { value: 'card', label: 'การ์ด', icon: 'card' as const },
  { value: 'table', label: 'ตาราง', icon: 'table' as const }
]

// Accordion & Card Expand State
const openProjectGroups = ref<Record<string, boolean>>({})
const openAppCards = ref<Record<string, boolean>>({})

const isProjectGroupOpen = (key: string): boolean => {
  return openProjectGroups.value[key] !== false
}

const toggleProjectGroup = (key: string) => {
  const current = isProjectGroupOpen(key)
  openProjectGroups.value[key] = !current
}

const isCardExpanded = (id: number | string): boolean => {
  return openAppCards.value[id] !== false
}

const toggleCardExpanded = (id: number | string) => {
  const current = isCardExpanded(id)
  openAppCards.value[id] = !current
}

// แยกแยะประเภทคำร้อง: คำร้องขอเช่าแผง vs คำร้องทั่วไป
const isStallLease = (app: any): boolean => {
  if (app.application_type === 'stall_lease') return true
  if (app.application_type === 'general') return false
  if (app.application_number?.startsWith('GEN-') || app.application_no?.startsWith('GEN-')) return false
  return Boolean(app.booking_id || (app.stall_number && app.stall_number !== '-') || (app.stall_ident && app.stall_ident !== '-'))
}

const canTenantEdit = (app: any): boolean => {
  if (!isTenant.value || !isStallLease(app)) return false
  const s = String(app.status || app.application_state || '').toLowerCase()
  return s === 'revision_required' || s === 'rejected'
}

const isGeneralApp = (app: any): boolean => !isStallLease(app)

const countByType = (type: 'stall_lease' | 'general'): number => {
  const list = Array.isArray(applications.value) ? applications.value : []
  if (type === 'stall_lease') {
    return list.filter(a => isStallLease(a)).length
  }
  return list.filter(a => isGeneralApp(a)).length
}

const loadMarkets = async () => {
  try {
    const res: any = await get('/markets')
    marketsList.value = Array.isArray(res) ? res : (Array.isArray(res?.data) ? res.data : [])
  } catch (err) {
    console.error('Failed to load markets:', err)
  }
}

const marketOptions = computed(() => {
  if (marketsList.value.length === 0) {
    return [{ value: 1, label: 'ตลาดนัดจตุจักร' }]
  }
  return marketsList.value.map(m => ({
    value: m.id,
    label: m.name
  }))
})

const loadApplications = async () => {
  loading.value = true
  try {
    const res: any = await get('/applications')
    applications.value = Array.isArray(res) ? res : (Array.isArray(res?.data) ? res.data : [])
    await fetchNotifications()
  } catch (err) {
    console.error('Failed to load applications:', err)
  } finally {
    loading.value = false
  }
}

// Counts for KPI Ribbon Cards
const countAll = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []
  if (typeFilter.value === 'stall_lease') list = list.filter(a => isStallLease(a))
  else if (typeFilter.value === 'general') list = list.filter(a => isGeneralApp(a))
  return list.length
})

const countPendingCombined = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []
  if (typeFilter.value === 'stall_lease') list = list.filter(a => isStallLease(a))
  else if (typeFilter.value === 'general') list = list.filter(a => isGeneralApp(a))
  return list.filter(a => [
    'pending_officer',
    'pending_debt_check',
    'debt_pending',
    'pending_market_director',
    'pending_bureau_director'
  ].includes(a.status)).length
})

const countPendingOfficer = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []
  if (typeFilter.value === 'stall_lease') list = list.filter(a => isStallLease(a))
  else if (typeFilter.value === 'general') list = list.filter(a => isGeneralApp(a))
  return list.filter(a => a.status === 'pending_officer').length
})

const countPendingDebt = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []
  if (typeFilter.value === 'stall_lease') list = list.filter(a => isStallLease(a))
  else if (typeFilter.value === 'general') list = list.filter(a => isGeneralApp(a))
  return list.filter(a => a.status === 'pending_debt_check' || a.status === 'debt_pending').length
})

const countPendingMarketDirector = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []
  if (typeFilter.value === 'stall_lease') list = list.filter(a => isStallLease(a))
  else if (typeFilter.value === 'general') list = list.filter(a => isGeneralApp(a))
  return list.filter(a => a.status === 'pending_market_director').length
})

const countPendingBureauDirector = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []
  if (typeFilter.value === 'stall_lease') list = list.filter(a => isStallLease(a))
  else if (typeFilter.value === 'general') list = list.filter(a => isGeneralApp(a))
  return list.filter(a => a.status === 'pending_bureau_director').length
})

const countApproved = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []
  if (typeFilter.value === 'stall_lease') list = list.filter(a => isStallLease(a))
  else if (typeFilter.value === 'general') list = list.filter(a => isGeneralApp(a))
  return list.filter(a => a.status === 'approved').length
})

const countRejected = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []
  if (typeFilter.value === 'stall_lease') list = list.filter(a => isStallLease(a))
  else if (typeFilter.value === 'general') list = list.filter(a => isGeneralApp(a))
  return list.filter(a => a.status === 'rejected' || a.status === 'revision_required').length
})

// Filtered Applications
const filteredApplications = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []

  // Filter by Type
  if (typeFilter.value === 'stall_lease') {
    list = list.filter(a => isStallLease(a))
  } else if (typeFilter.value === 'general') {
    list = list.filter(a => isGeneralApp(a))
  }

  // Filter by Active Tab
  if (activeTab.value === 'pending') {
    if (pendingSubFilter.value === 'all') {
      list = list.filter(a => [
        'pending_officer',
        'pending_debt_check',
        'debt_pending',
        'pending_market_director',
        'pending_bureau_director'
      ].includes(a.status))
    } else if (pendingSubFilter.value === 'pending_officer') {
      list = list.filter(a => a.status === 'pending_officer')
    } else if (pendingSubFilter.value === 'pending_debt') {
      list = list.filter(a => a.status === 'pending_debt_check' || a.status === 'debt_pending')
    } else if (pendingSubFilter.value === 'pending_market_director') {
      list = list.filter(a => a.status === 'pending_market_director')
    } else if (pendingSubFilter.value === 'pending_bureau_director') {
      list = list.filter(a => a.status === 'pending_bureau_director')
    }
  } else if (activeTab.value === 'approved') {
    list = list.filter(a => a.status === 'approved')
  } else if (activeTab.value === 'rejected') {
    list = list.filter(a => a.status === 'rejected' || a.status === 'revision_required')
  }

  // Filter by Search Query
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(a =>
      (a.application_number && a.application_number.toLowerCase().includes(q)) ||
      (a.applicant_name && a.applicant_name.toLowerCase().includes(q)) ||
      (a.applicant_phone && a.applicant_phone.toLowerCase().includes(q)) ||
      (a.shop_name && a.shop_name.toLowerCase().includes(q)) ||
      (a.stall_number && a.stall_number.toLowerCase().includes(q)) ||
      (a.subject && a.subject.toLowerCase().includes(q))
    )
  }

  return list
})

const paginatedApplications = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredApplications.value.slice(start, start + pageSize.value)
})

// Project Name extraction for grouping
const getApplicationProjectName = (app: any): string => {
  if (app.project_name) return app.project_name
  if (app.zone_name) return app.zone_name
  if (app.stall_zone) return app.stall_zone
  const stall = app.stall_number || app.stall_ident || ''
  if (stall) {
    const match = stall.match(/(?:JJ-?|PROJ-?)(\d+)/i)
    if (match && match[1]) {
      return `โครงการที่ ${match[1]}`
    }
    const parts = stall.split('-')
    if (parts.length >= 2 && !isNaN(Number(parts[1]))) {
      return `โครงการที่ ${parts[1]}`
    }
  }
  return app.market_name ? app.market_name : 'โครงการทั่วไป'
}

// Group paginated applications by project
interface ApplicationProjectGroup {
  key: string
  projectName: string
  applications: any[]
}

const paginatedProjectGroups = computed<ApplicationProjectGroup[]>(() => {
  const map = new Map<string, ApplicationProjectGroup>()

  paginatedApplications.value.forEach(app => {
    const projName = getApplicationProjectName(app)
    const key = `${app.market_name || 'JJ'}::${projName}`
    if (!map.has(key)) {
      map.set(key, {
        key,
        projectName: projName,
        applications: []
      })
    }
    map.get(key)!.applications.push(app)
  })

  return Array.from(map.values())
})

// Thai Date and Time Formats
const THAI_MONTHS = [
  'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
  'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
]
const THAI_DAYS = ['อาทิตย์', 'จันทร์', 'อังคาร', 'พุธ', 'พฤหัสบดี', 'ศุกร์', 'เสาร์']

const formatSubmissionTimestamp = (dt?: string): string => {
  if (!dt) return '-'
  try {
    const d = new Date(dt)
    if (isNaN(d.getTime())) return dt.slice(0, 10)
    const day = d.getDate()
    const month = THAI_MONTHS[d.getMonth()]
    const year = d.getFullYear() + 543
    const hours = String(d.getHours()).padStart(2, '0')
    const mins = String(d.getMinutes()).padStart(2, '0')
    return `${day} ${month} ${year} เวลา ${hours}.${mins} น.`
  } catch {
    return dt.slice(0, 10)
  }
}

const formatThaiStartDate = (dt?: string): string => {
  if (!dt) return '-'
  try {
    const d = new Date(dt)
    if (isNaN(d.getTime())) return dt.slice(0, 10)
    const dayOfWeek = THAI_DAYS[d.getDay()]
    const day = d.getDate()
    const month = THAI_MONTHS[d.getMonth()]
    const year = d.getFullYear() + 543
    return `${dayOfWeek} ${day} ${month} ${year}`
  } catch {
    return dt.slice(0, 10)
  }
}

const getDaysRemainingBadge = (dt?: string): string => {
  if (!dt) return ''
  try {
    const target = new Date(dt)
    const now = new Date()
    target.setHours(0, 0, 0, 0)
    now.setHours(0, 0, 0, 0)
    const diffTime = target.getTime() - now.getTime()
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    if (diffDays > 0) return `อีก ${diffDays} วัน`
    if (diffDays === 0) return 'วันนี้'
    return ''
  } catch {
    return ''
  }
}

// Stepper Logic (4 Steps)
const isAppStepCompleted = (status: string, step: number): boolean => {
  // เมื่อเอกสารคำร้องถูกตีกลับมาแก้ไข (revision_required) สเต็ปที่ 1 ยังไม่เสร็จสิ้น ต้องรอให้ผู้เช่าส่งกลับ
  if (status === 'revision_required') {
    return false
  }
  if (step === 1) return true
  if (step === 2) {
    return ['pending_market_director', 'pending_bureau_director', 'approved', 'contract_pending', 'completed'].includes(status)
  }
  if (step === 3) {
    return ['pending_bureau_director', 'approved', 'contract_pending', 'completed'].includes(status)
  }
  if (step === 4) {
    return ['approved', 'contract_pending', 'completed'].includes(status)
  }
  return false
}

const isAppStepActive = (status: string, step: number): boolean => {
  // เมื่อถูกส่งกลับแก้ไข สเต็ปด้านบนต้องเป็น Step ที่ 1 (ไม่ใช่ 2)
  if (status === 'revision_required') {
    return step === 1
  }
  if (step === 1) return false
  if (step === 2) {
    return ['pending_officer', 'pending_debt_check', 'debt_pending'].includes(status)
  }
  if (step === 3) {
    return status === 'pending_market_director'
  }
  if (step === 4) {
    return status === 'pending_bureau_director'
  }
  return false
}

// Status Badges, Colors, and Tones
const getStatusBadgeLabel = (status: string): string => {
  switch (status) {
    case 'pending_officer':
      return isTenant.value ? 'รอเจ้าหน้าที่ตรวจเอกสาร' : 'รอเจ้าหน้าที่ตรวจสอบ'
    case 'pending_debt_check':
      return 'รอ ตรวจสอบหนี้'
    case 'debt_pending':
      return 'มีหนี้ค้างชำระ'
    case 'pending_market_director':
      return 'รอ ผอ.ตลาดพิจารณา'
    case 'pending_bureau_director':
      return 'รอ ผอ.สำนักตลาดพิจารณา'
    case 'approved':
      return 'อนุมัติแล้ว'
    case 'contract_pending':
      return 'รอตรวจสัญญา'
    case 'completed':
      return 'ทำสัญญาแล้ว'
    case 'revision_required':
      return isTenant.value ? 'ต้องแก้ไขเอกสาร' : 'ส่งกลับให้ผู้เช่าแก้ไข'
    case 'rejected':
      return 'ไม่อนุมัติ'
    default:
      return status || 'รอดำเนินการ'
  }
}

const getStatusPillClass = (status: string): string => {
  switch (status) {
    case 'pending_officer':
    case 'pending_debt_check':
    case 'debt_pending':
      return 'pill-amber'
    case 'pending_market_director':
      return 'pill-blue'
    case 'pending_bureau_director':
      return 'pill-purple'
    case 'approved':
    case 'completed':
      return 'pill-green'
    case 'contract_pending':
      return 'pill-blue'
    case 'revision_required':
      return 'pill-amber'
    case 'rejected':
      return 'pill-red'
    default:
      return 'pill-gray'
  }
}

const getCardBorderClass = (status: string): string => {
  switch (status) {
    case 'pending_officer':
    case 'pending_debt_check':
    case 'debt_pending':
    case 'revision_required':
      return 'border-status-pending'
    case 'pending_market_director':
    case 'contract_pending':
      return 'border-status-ready'
    case 'pending_bureau_director':
      return 'border-status-applied'
    case 'approved':
    case 'completed':
      return 'border-status-approved'
    default:
      return 'border-status-muted'
  }
}

const getProgressTone = (status: string): string => {
  switch (status) {
    case 'approved':
    case 'completed':
      return 'tone-success'
    case 'contract_pending':
      return 'tone-action'
    case 'pending_officer':
    case 'pending_debt_check':
    case 'debt_pending':
    case 'revision_required':
      return 'tone-waiting'
    case 'pending_market_director':
    case 'pending_bureau_director':
      return 'tone-action'
    default:
      return 'tone-muted'
  }
}

const getApplicationProgressText = (app: any): string => {
  const status = app.status || ''
  switch (status) {
    case 'revision_required': {
      const note = app.review_note || app.officer_comment || 'กรุณาตรวจสอบเอกสารและแก้ไขข้อมูลคำร้อง'
      return isTenant.value
        ? `เจ้าหน้าที่ส่งกลับมาแก้ไข: ${note}`
        : `ส่งกลับให้ผู้เช่าแก้ไข: ${note} (รอผู้เช่าแก้ไขและส่งกลับ)`
    }
    case 'pending_officer':
      return isTenant.value
        ? 'อยู่ระหว่างรอเจ้าหน้าที่ตลาดตรวจสอบเอกสารความถูกต้อง'
        : 'รอเจ้าหน้าที่ตรวจสอบเอกสารและความถูกต้อง (สามารถอนุมัติส่งต่อ, ส่งกลับแก้ไข หรือไม่อนุมัติ)'
    case 'pending_debt_check':
    case 'debt_pending':
      return 'อยู่ระหว่างรอตรวจสอบประวัติหนี้ค้างชำระค่าเช่า'
    case 'pending_market_director':
      return 'เจ้าหน้าที่ตรวจสอบแล้ว รอ ผอ.ตลาด พิจารณาให้ความเห็นชอบ'
    case 'pending_bureau_director':
      return 'ผอ.ตลาด ให้ความเห็นชอบแล้ว รอ ผอ.สำนักตลาด ลงนามอนุมัติ'
    case 'approved':
      return isStallLease(app)
        ? 'อนุมัติคำร้องขอเช่าแล้ว พร้อมดำเนินการจัดทำสัญญาเช่าแผงค้า'
        : 'อนุมัติคำร้องเรียบร้อยแล้ว'
    case 'contract_pending':
      return 'ผู้เช่าจัดทำสัญญาเช่าแล้ว อยู่ระหว่างรอเจ้าหน้าที่ตรวจสอบสัญญา'
    case 'completed':
      return 'สัญญาเช่าได้รับการอนุมัติและมีผลบังคับใช้เรียบร้อยแล้ว'
    case 'rejected':
      return isTenant.value
        ? `คำร้องไม่ผ่านการอนุมัติ: ${app.review_note || app.bureau_director_comment || app.market_director_comment || app.officer_comment || 'ไม่ผ่านเกณฑ์การพิจารณา'} (สามารถกดแก้ไขคำร้องเพื่อยื่นใหม่ได้)`
        : `ไม่อนุมัติคำร้อง: ${app.review_note || app.bureau_director_comment || app.market_director_comment || app.officer_comment || 'ไม่ผ่านเกณฑ์การพิจารณา'}`
    default:
      return 'อยู่ระหว่างดำเนินการตามขั้นตอน'
  }
}

// Actor and Review Helpers for Table View
const getActionActorInfo = (app: any) => {
  const status = app.status || app.application_state
  switch (status) {
    case 'pending_officer':
      return { role: 'เจ้าหน้าที่ตลาด', actionText: 'รอ จนท. ตรวจสอบ', tagClass: 'actor-tag-officer-pending' }
    case 'pending_debt_check':
      return { role: 'เจ้าหน้าที่การคลัง/ตลาด', actionText: 'รอ ตรวจสอบหนี้ค้าง', tagClass: 'actor-tag-officer' }
    case 'debt_pending':
      return { role: 'ผอ.ตลาด / การคลัง', actionText: 'รอ เคลียร์หนี้ค้างชำระ', tagClass: 'actor-tag-warning' }
    case 'pending_market_director':
      return { role: 'ผอ.ตลาด', actionText: 'รอ ผอ.ตลาด', tagClass: 'actor-tag-director' }
    case 'pending_bureau_director':
      return { role: 'ผอ.สำนักตลาด', actionText: 'รอ ผอ.สำนักตลาด', tagClass: 'actor-tag-director' }
    case 'approved':
      return { role: 'ผอ.สำนักตลาด', actionText: 'ผอ.สำนักตลาด', tagClass: 'actor-tag-success' }
    case 'revision_required':
      return { role: 'เจ้าหน้าที่ตลาด', actionText: 'เจ้าหน้าที่ตลาด', tagClass: 'actor-tag-warning' }
    case 'rejected': {
      let actor = 'ผู้มีอำนาจพิจารณา'
      if (app.bureau_director_comment) actor = 'ผอ.สำนักตลาด'
      else if (app.market_director_comment) actor = 'ผอ.ตลาด'
      else if (app.officer_comment) actor = 'เจ้าหน้าที่ตลาด'
      return { role: actor, actionText: actor, tagClass: 'actor-tag-danger' }
    }
    default:
      return { role: '-', actionText: '-', tagClass: 'actor-tag-default' }
  }
}

const getReviewNoteInfo = (app: any) => {
  if (app.bureau_director_comment) {
    return { role: 'ผอ.สำนักตลาด', text: app.bureau_director_comment, type: 'director' }
  }
  if (app.market_director_comment) {
    return { role: 'ผอ.ตลาด', text: app.market_director_comment, type: 'director' }
  }
  if (app.officer_comment) {
    return { role: 'เจ้าหน้าที่', text: app.officer_comment, type: 'officer' }
  }
  if (app.review_note) {
    return { role: 'หมายเหตุ', text: app.review_note, type: 'note' }
  }
  return { role: '', text: '', type: 'none' }
}

const getAllCommentsList = (app: any) => {
  const comments: Array<{ role: string; text: string; step: string }> = []
  if (app.officer_comment) {
    comments.push({ role: 'เจ้าหน้าที่ตลาด', text: app.officer_comment, step: 'ขั้นที่ 1 ตรวจสอบเอกสาร' })
  }
  if (app.market_director_comment) {
    comments.push({ role: 'ผอ.ตลาด', text: app.market_director_comment, step: 'ขั้นที่ 2 ให้ความเห็นชอบ' })
  }
  if (app.bureau_director_comment) {
    comments.push({ role: 'ผอ.สำนักตลาด', text: app.bureau_director_comment, step: 'ขั้นที่ 3 อนุมัติขั้นสุดท้าย' })
  }
  if (comments.length === 0 && app.review_note) {
    comments.push({ role: 'หมายเหตุ', text: app.review_note, step: 'บันทึก' })
  }
  return comments
}

// Approval and Actions
const processApproval = async (app: any, defaultComment: string) => {
  const confirmed = await showConfirm(
    `ยืนยันการอนุมัติคำร้องเลขที่ ${app.application_number} หรือไม่?`,
    'ยืนยันการอนุมัติ',
    { confirmButtonText: 'อนุมัติคำร้อง', cancelButtonText: 'ยกเลิก' }
  )
  if (!confirmed) return

  try {
    await post(`/applications/${app.id}/action`, {
      action: 'approve',
      comment: defaultComment
    })
    await loadApplications()
    await showSuccess('บันทึกการอนุมัติเรียบร้อยแล้ว', 'อนุมัติสำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || err?.data?.message || 'เกิดข้อผิดพลาดในการอนุมัติ', 'เกิดข้อผิดพลาด')
  }
}

// Detail View
const downloadIDCardCopy = async (app: any) => {
  try {
    const file = await apiFetch<Blob>(`/applications/${app.id}/id-card-copy`, { responseType: 'blob' })
    const url = URL.createObjectURL(file)
    const link = document.createElement('a')
    link.href = url
    link.download = 'id-card-copy'
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(url)
  } catch {
    await showError('ไม่สามารถเปิดสำเนาบัตรประชาชนได้ กรุณาตรวจสอบสิทธิ์', 'เกิดข้อผิดพลาด')
  }
}
const openDetails = async (app: any) => {
  const isStall = isStallLease(app)
  const typeText = isStall ? 'คำร้องขอเช่าแผงค้า' : 'คำร้องทั่วไป'
  const stallText = isStall
    ? `\nแผงค้า: ${app.stall_number || app.stall_ident}`
    : (app.stall_number && app.stall_number !== '-' ? `\nแผงที่เกี่ยวข้อง: ${app.stall_number}` : '\nแผงค้า: คำร้องทั่วไป (ไม่มีแผง)')
  const marketText = app.market_name ? `\nตลาด: ${app.market_name}` : ''
  const subjectText = isStall
    ? `\nชื่อร้านค้า: ${app.shop_name} (${app.product_type || app.product_category || 'ทั่วไป'})`
    : `\nหัวข้อคำร้อง: ${app.subject || app.shop_name || 'คำร้องทั่วไป'}`
  const idCardText = `${app.applicant_id_card ? `\nเลขประจำตัว ปชช.: ${app.applicant_id_card}` : ''}${app.applicant_age ? `\nอายุ: ${app.applicant_age} ปี` : ''}`
  const emailText = app.applicant_email ? `\nอีเมล: ${app.applicant_email}` : ''
  const addressText = app.applicant_address && app.applicant_address !== '-' ? `\nที่อยู่: ${app.applicant_address}` : ''
  const detailText = (app.details || app.review_note)
    ? `\nรายละเอียดคำร้อง: ${app.details || app.review_note}`
    : ''

  const comments = getAllCommentsList(app)
  const commentDetails = comments.length > 0
    ? '\n\nบันทึกความเห็นการพิจารณา:\n' + comments.map(c => `• ${c.role} (${c.step}): ${c.text}`).join('\n')
    : ''

  await showAlert(
    `ประเภท: ${typeText}\nเลขที่คำร้อง: ${app.application_number || app.application_no}${marketText}${stallText}${subjectText}\nผู้ยื่นคำร้อง: ${app.applicant_name} (โทร. ${app.applicant_phone})${idCardText}${emailText}${addressText}\nสถานะ: ${getStatusBadgeLabel(app.status)}\nขั้นตอน: ${getActionActorInfo(app).actionText}${detailText}${commentDetails}`,
    { title: 'รายละเอียดคำร้อง', type: 'info' }
  )
}

// Reject Modal State
const rejectAction = ref<'reject' | 'return_revision'>('reject')
const showRejectModal = ref(false)
const targetAppForReject = ref<any>(null)
const rejectComment = ref('')
const rejectErrorMsg = ref('')
const rejectSubmitting = ref(false)

const quickRejectReasons = [
  'เอกสารหลักฐานประกอบคำร้องไม่ครบถ้วน',
  'มีประวัติค้างชำระค่าเช่าหรือค่าบริการ',
  'คุณสมบัติไม่ผ่านเกณฑ์การพิจารณาเช่าแผงค้า',
  'ข้อมูลสำเนาบัตรประชาชนหรือทะเบียนบ้านไม่ชัดเจน',
  'ประเภทสินค้าไม่สอดคล้องกับโซนการค้าของแผง'
]

const openRejectModal = (app: any, action: 'reject' | 'return_revision' = 'reject') => {
  rejectAction.value = action
  targetAppForReject.value = app
  rejectComment.value = ''
  rejectErrorMsg.value = ''
  showRejectModal.value = true
}

const closeRejectModal = () => {
  showRejectModal.value = false
  targetAppForReject.value = null
  rejectComment.value = ''
  rejectErrorMsg.value = ''
}

const setQuickReason = (reason: string) => {
  if (!rejectComment.value) {
    rejectComment.value = reason
  } else if (!rejectComment.value.includes(reason)) {
    rejectComment.value += ' ' + reason
  }
}

const submitReject = async () => {
  if (!targetAppForReject.value) return
  const comment = rejectComment.value.trim()
  if (!comment) {
    rejectErrorMsg.value = 'กรุณาระบุหมายเหตุหรือเหตุผลประกอบ'
    return
  }

  rejectSubmitting.value = true
  rejectErrorMsg.value = ''

  try {
    await post(`/applications/${targetAppForReject.value.id}/action`, {
      action: rejectAction.value,
      comment: comment,
      reason: comment
    })
    const wasRevision = rejectAction.value === 'return_revision'
    await loadApplications()
    closeRejectModal()
    await showSuccess(wasRevision ? 'ส่งคำร้องกลับเพื่อแก้ไขแล้ว' : 'บันทึกผลการไม่อนุมัติคำร้องเรียบร้อยแล้ว', 'บันทึกสำเร็จ')
  } catch (err: any) {
    rejectErrorMsg.value = err?.data?.detail || err?.data?.message || 'เกิดข้อผิดพลาดในการบันทึก'
  } finally {
    rejectSubmitting.value = false
  }
}

// Modal ยื่นคำร้องทั่วไป
const showGeneralModal = ref(false)
const generalSubmitting = ref(false)
const generalFormError = ref('')
const generalForm = ref({
  applicant_name: '',
  applicant_phone: '',
  applicant_id_card: '',
  applicant_email: '',
  applicant_address: '',
  market_id: 1,
  stall_number: '',
  subject: '',
  details: ''
})

const openCreateGeneralModal = () => {
  generalForm.value = {
    applicant_name: (currentUser.value?.first_name ? `${currentUser.value.first_name} ${currentUser.value.last_name || ''}` : '') || user.value?.name || '',
    applicant_phone: currentUser.value?.phone || user.value?.phone || '',
    applicant_id_card: '',
    applicant_email: currentUser.value?.email || user.value?.email || '',
    applicant_address: '',
    market_id: 1,
    stall_number: '',
    subject: '',
    details: ''
  }
  generalFormError.value = ''
  showGeneralModal.value = true
}

const closeCreateGeneralModal = () => {
  showGeneralModal.value = false
  generalFormError.value = ''
}

const submitGeneralApplication = async () => {
  if (!generalForm.value.applicant_name || !generalForm.value.applicant_phone || !generalForm.value.subject || !generalForm.value.details) {
    generalFormError.value = 'กรุณากรอกข้อมูลในช่องที่มีเครื่องหมาย * ให้ครบถ้วน'
    return
  }

  generalSubmitting.value = true
  generalFormError.value = ''

  try {
    await post('/applications/general', {
      applicant_name: generalForm.value.applicant_name,
      applicant_phone: generalForm.value.applicant_phone,
      applicant_id_card: generalForm.value.applicant_id_card,
      applicant_email: generalForm.value.applicant_email,
      applicant_address: generalForm.value.applicant_address,
      market_id: generalForm.value.market_id,
      stall_number: generalForm.value.stall_number || '-',
      subject: generalForm.value.subject,
      details: generalForm.value.details
    })

    await loadApplications()
    closeCreateGeneralModal()
    await showSuccess('ยื่นคำร้องทั่วไปเรียบร้อยแล้ว เจ้าหน้าที่จะดำเนินการตรวจสอบต่อไป', 'ยื่นคำร้องสำเร็จ')
  } catch (err: any) {
    generalFormError.value = err?.data?.detail || err?.data?.message || 'เกิดข้อผิดพลาดในการยื่นคำร้อง'
  } finally {
    generalSubmitting.value = false
  }
const exportReport = () => {
  const list = filteredApplications.value
  if (!list || list.length === 0) {
    showError('ไม่มีรายการคำร้องที่จะส่งออก', 'ไม่พบข้อมูล')
    return
  }

  const headers = [
    'เลขที่คำร้อง',
    'ประเภทคำร้อง',
    'เรื่อง/แผงค้า',
    'ตลาด',
    'ผู้ยื่นคำร้อง',
    'เบอร์โทรศัพท์',
    'เลขประจำตัวประชาชน',
    'ชื่อร้านค้า',
    'ประเภทสินค้า',
    'รายละเอียดคำร้อง',
    'วันที่ยื่น',
    'สถานะ',
    'ขั้นตอนปัจจุบัน',
    'หมายเหตุการพิจารณา'
  ]

  const rows = list.map(app => {
    const isStall = isStallLease(app)
    const typeLabel = isStall ? 'คำร้องขอเช่าแผงค้า' : 'คำร้องทั่วไป'
    const subjectOrStall = isStall ? `แผง ${app.stall_number || app.stall_ident}` : (app.subject || app.shop_name || '-')
    const comments = getAllCommentsList(app).map(c => `${c.role}: ${c.text}`).join(' | ')
    const note = comments || app.details || app.review_note || '-'

    return [
      `"${(app.application_number || app.application_no || '').replace(/"/g, '""')}"`,
      `"${typeLabel}"`,
      `"${subjectOrStall.replace(/"/g, '""')}"`,
      `"${(app.market_name || '-').replace(/"/g, '""')}"`,
      `"${(app.applicant_name || '').replace(/"/g, '""')}"`,
      `"${(app.applicant_phone || '').replace(/"/g, '""')}"`,
      `"${(app.applicant_id_card || '-').replace(/"/g, '""')}"`,
      `"${(app.shop_name || '-').replace(/"/g, '""')}"`,
      `"${(app.product_type || app.product_category || '-').replace(/"/g, '""')}"`,
      `"${(app.details || app.review_note || '-').replace(/"/g, '""')}"`,
      `"${(app.created_at?.slice(0, 10) || '-').replace(/"/g, '""')}"`,
      `"${getStatusBadgeLabel(app.status)}"`,
      `"${getActionActorInfo(app).actionText.replace(/"/g, '""')}"`,
      `"${note.replace(/"/g, '""')}"`
    ].join(',')
  })

  const csvContent = '\uFEFF' + [headers.join(','), ...rows].join('\r\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute('download', `applications_report_${new Date().toISOString().slice(0, 10)}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)

  showSuccess(`ส่งออกรายงานคำร้อง ${list.length} รายการ เป็นไฟล์ CSV เรียบร้อยแล้ว`, 'ส่งออกรายงานสำเร็จ')
}

// History Panel
const showHistoryPanel = ref(false)
const applicationHistory = computed(() => {
  return applications.value.map(a => ({
    id: a.id,
    action: a.status === 'approved' ? 'อนุมัติคำร้อง' : (a.status === 'rejected' ? 'ไม่อนุมัติ' : 'ยื่นคำร้อง'),
    timestamp: a.created_at || a.updated_at,
    note: a.review_note || a.subject || a.application_number
  }))
})

// Watchers
watch(activeTab, () => {
  pendingSubFilter.value = 'all'
  currentPage.value = 1
})

watch(pendingSubFilter, () => {
  currentPage.value = 1
})

watch(typeFilter, () => {
  currentPage.value = 1
})

watch(viewMode, (newVal) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('applications_view_mode', newVal)
  }
  currentPage.value = 1
})

onMounted(async () => {
  if (typeof window !== 'undefined') {
    authStore.hydrate()
    const saved = localStorage.getItem('applications_view_mode')
    if (saved === 'table' || saved === 'card') {
      viewMode.value = saved
    }
  }
  if (route.query.search) {
    searchQuery.value = String(route.query.search)
  }
  await Promise.all([
    loadMarkets(),
    loadApplications()
  ])
})
</script>


<style scoped>

/* -------------------------------------------------------------
 * Booking-Style Cards & KPI Ribbon Styles (Copied from BookingMain)
 * ------------------------------------------------------------- */
.card-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 22px;
  gap: 16px;
  flex-wrap: wrap;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.section-title {
  font-size: 1.35rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.02em;
  margin: 0;
}

.item-count-subtitle {
  font-size: 0.825rem;
  color: #64748b;
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  min-width: 260px;
}

.search-input {
  width: 100%;
  padding: 8px 32px 8px 34px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 9px;
  font-size: 0.83rem;
  color: #1e293b;
  font-family: inherit;
  transition: all 0.2s ease;
  height: 38px;
}

.search-input:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.search-icon {
  position: absolute;
  left: 10px;
  color: #94a3b8;
  pointer-events: none;
}

.clear-search-btn {
  position: absolute;
  right: 8px;
  background: transparent;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 3px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.clear-search-btn:hover {
  color: #0f172a;
}

/* Officer KPI Ribbon */
.officer-kpi-ribbon {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 24px;
}

@media (max-width: 1100px) {
  .officer-kpi-ribbon {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .officer-kpi-ribbon {
    grid-template-columns: 1fr;
  }
}

.kpi-ribbon-card {
  background: #ffffff;
  border: 1px solid #eef2f6;
  border-radius: 14px;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 14px;
  cursor: pointer;
  font-family: inherit;
  text-align: left;
  position: relative;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}

.kpi-ribbon-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(15, 23, 42, 0.07);
}

.kpi-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.kpi-text-box {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.kpi-title {
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 500;
  line-height: 1.2;
  margin-bottom: 4px;
  white-space: nowrap;
}

.kpi-value-row {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.kpi-count {
  font-size: 1.45rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1;
}

.kpi-unit-pill {
  font-size: 0.78rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
  line-height: 1.2;
}

/* 1. ทั้งหมด */
.kpi-all .kpi-icon-box {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #dbeafe;
}
.kpi-all .kpi-unit-pill {
  background: #eff6ff;
  color: #1d4ed8;
}
.kpi-all.active {
  border-color: #2563eb;
  background: #f8faff;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.18);
}

/* 2. รอดำเนินการ */
.kpi-pending .kpi-icon-box {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #fef3c7;
}
.kpi-pending .kpi-unit-pill {
  background: #fef3c7;
  color: #b45309;
}
.kpi-pending.active {
  border-color: #f59e0b;
  background: #fffdf7;
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.2);
}

/* 3. อนุมัติแล้ว */
.kpi-applied .kpi-icon-box {
  background: #f5f3ff;
  color: #7c3aed;
  border: 1px solid #ede9fe;
}
.kpi-applied .kpi-unit-pill {
  background: #ede9fe;
  color: #6d28d9;
}
.kpi-applied.active {
  border-color: #8b5cf6;
  background: #faf5ff;
  box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.2);
}

/* 4. ยกเลิก/ส่งกลับ */
.kpi-cancelled .kpi-icon-box {
  background: #f1f5f9;
  color: #64748b;
  border: 1px solid #e2e8f0;
}
.kpi-cancelled .kpi-unit-pill {
  background: #f1f5f9;
  color: #64748b;
}
.kpi-cancelled.active {
  border-color: #94a3b8;
  background: #f8fafc;
  box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.18);
}

/* Popover Breakdown */
.kpi-info-wrapper {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
}

.kpi-info-btn {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 0;
}

.kpi-info-btn:hover {
  background: #eff6ff;
  border-color: #93c5fd;
  color: #2563eb;
  transform: scale(1.08);
}

.kpi-info-popover {
  position: absolute;
  top: calc(100% + 6px);
  right: -4px;
  width: 250px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 10px 25px -5px rgba(15, 23, 42, 0.15), 0 8px 10px -6px rgba(15, 23, 42, 0.1);
  padding: 12px;
  z-index: 50;
  cursor: default;
  text-align: left;
}

.kpi-popover-header {
  font-size: 0.78rem;
  font-weight: 700;
  color: #1e293b;
  display: flex;
  align-items: center;
  gap: 6px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f1f5f9;
  margin-bottom: 8px;
}

.kpi-popover-body {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.kpi-popover-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 8px;
  border-radius: 8px;
  font-size: 0.8rem;
  color: #334155;
  cursor: pointer;
  transition: background 0.15s ease;
}

.kpi-popover-item:hover {
  background: #f8fafc;
}

.kpi-popover-item.is-selected {
  background: #eff6ff;
  font-weight: 600;
  color: #2563eb;
}

.kpi-status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  margin-right: 6px;
  display: inline-block;
  flex-shrink: 0;
}
.kpi-status-dot.amber { background: #f59e0b; }
.kpi-status-dot.blue { background: #2563eb; }
.kpi-status-dot.purple { background: #8b5cf6; }
.kpi-status-dot.green { background: #10b981; }

.kpi-status-name {
  flex: 1;
  font-size: 0.78rem;
}

.kpi-status-val {
  font-weight: 700;
  font-size: 0.8rem;
}

.kpi-popover-reset {
  margin-top: 8px;
  padding-top: 6px;
  border-top: 1px solid #f1f5f9;
  font-size: 0.75rem;
  color: #2563eb;
  font-weight: 600;
  text-align: center;
  cursor: pointer;
}

/* Sub-filter bar */
.tenant-subfilter-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
}

.subfilter-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 9999px;
  font-size: 0.8rem;
  font-weight: 600;
  color: #475569;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.subfilter-chip:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.subfilter-chip.active {
  background: #eff6ff;
  color: #1d4ed8;
  border-color: #93c5fd;
  box-shadow: 0 1px 3px rgba(37, 99, 235, 0.1);
}

.subfilter-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}
.chip-amber .subfilter-dot { background: #f59e0b; }
.chip-blue .subfilter-dot { background: #2563eb; }
.chip-purple .subfilter-dot { background: #8b5cf6; }
.chip-green .subfilter-dot { background: #10b981; }

.subfilter-badge {
  font-size: 0.75rem;
  padding: 1px 6px;
  border-radius: 10px;
  background: #e2e8f0;
  color: #334155;
  margin-left: 2px;
}

.subfilter-chip.active .subfilter-badge {
  background: #dbeafe;
  color: #1e40af;
}

/* Accordion Project Headers */
.officer-groups-wrapper {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.project-stall-entry {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.officer-project-header {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 18px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s ease;
}

.officer-project-header:hover {
  background: #f1f5f9;
}

.officer-project-header .header-left-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.chevron-indicator {
  color: #64748b;
  transition: transform 0.25s ease;
}

.chevron-indicator.rotated {
  transform: rotate(180deg);
}

.accordion-header-icon {
  color: #2563eb;
  display: flex;
  align-items: center;
}

.market-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: #1e293b;
}

.market-count-pill {
  font-size: 0.78rem;
  font-weight: 600;
  padding: 3px 10px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  color: #475569;
}

.market-accordion-body {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* Booking Item Card (Application) */
.booking-item-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-left: 4px solid #cbd5e1;
  border-radius: 14px;
  padding: 20px 22px;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
}

.booking-item-card:hover {
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.05);
}

.booking-item-card.border-status-pending {
  border-left: 4px solid #f59e0b;
}

.booking-item-card.border-status-ready {
  border-left: 4px solid #2563eb;
}

.booking-item-card.border-status-applied {
  border-left: 4px solid #8b5cf6;
}

.booking-item-card.border-status-approved {
  border-left: 4px solid #10b981;
}

.booking-item-card.border-status-muted {
  border-left: 4px solid #94a3b8;
}

/* Card Top Row */
.card-top-row {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.card-top-main {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  width: 100%;
}

.top-row-left {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.circle-status-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.circle-status-icon.blue-circle {
  background: linear-gradient(135deg, #eff6ff 0%, #e0e7ff 100%);
  color: #2563eb;
  border: 1px solid #bfdbfe;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.08);
}

.booking-number-block {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.booking-label-row {
  display: flex;
  align-items: center;
}

.booking-label {
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 600;
  line-height: 1.2;
}

.booking-id-text {
  font-size: 1.25rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.02em;
  line-height: 1.2;
}

.officer-booked-at {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  margin-top: 3px;
  color: #94a3b8;
  font-size: 0.78rem;
  font-weight: 500;
}

.top-row-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.status-pill-badge {
  font-size: 0.8rem;
  font-weight: 600;
  padding: 4px 14px;
  border-radius: 20px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
}

.pill-green {
  background: #dcfce7;
  color: #15803d;
}

.pill-amber {
  background: #fef3c7;
  color: #d97706;
}

.pill-blue {
  background: #eff6ff;
  color: #2563eb;
}

.pill-purple {
  background: #f5f3ff;
  color: #7c3aed;
}

.pill-gray {
  background: #f1f5f9;
  color: #64748b;
}

.pill-red {
  background: #fee2e2;
  color: #dc2626;
}

.card-collapse-btn {
  background: transparent;
  border: none;
  color: #64748b;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-chevron {
  transition: transform 0.2s ease;
}

.card-chevron.rotated {
  transform: rotate(180deg);
}

/* Context Meta Tags Row */
.card-meta-tags-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  color: #64748b;
  font-size: 0.8rem;
}

.compact-meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.compact-meta-item + .compact-meta-item {
  padding-left: 12px;
  border-left: 1px solid #e2e8f0;
}

.compact-meta-item svg {
  color: #94a3b8;
  flex-shrink: 0;
}

.compact-meta-item strong {
  color: #334155;
  font-weight: 650;
}

.compact-meta-stall strong {
  color: #1d4ed8;
}

.compact-meta-item a {
  color: #64748b;
  text-decoration: none;
}

.compact-meta-item a:hover {
  color: #2563eb;
  text-decoration: underline;
}

/* Stepper Track */
.booking-stepper-container {
  margin: 18px 0 16px 0;
  padding: 12px 16px;
  background: #f8fafc;
  border-radius: 10px;
  border: 1px solid #f1f5f9;
}

.stepper-track {
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  max-width: 680px;
  margin: 0 auto;
}

.stepper-step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  z-index: 2;
  position: relative;
}

.step-circle {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: #e2e8f0;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.25s ease;
}

.step-label-text {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  white-space: nowrap;
}

.stepper-step.completed .step-circle {
  background: #10b981;
  color: #ffffff;
}

.stepper-step.completed .step-label-text {
  color: #0f172a;
}

.stepper-step.active .step-circle {
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.15);
}

.stepper-step.active .step-label-text {
  color: #2563eb;
  font-weight: 700;
}

.stepper-connector {
  flex: 1;
  height: 3px;
  background: #e2e8f0;
  margin: -18px 8px 0 8px;
  z-index: 1;
  transition: background 0.25s ease;
}

.stepper-connector.filled {
  background: #10b981;
}

/* 4-Column Info Grid */
.info-columns-grid {
  display: grid;
  grid-template-columns: 1.25fr 1.15fr 1fr 1.35fr;
  gap: 16px;
  padding: 14px 0 8px 0;
  border-top: 1px solid #f1f5f9;
}

@media (max-width: 900px) {
  .info-columns-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 580px) {
  .info-columns-grid {
    grid-template-columns: 1fr;
  }
}

.info-column {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.76rem;
  color: #64748b;
  font-weight: 500;
}

.info-value-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.info-value-text {
  font-size: 0.92rem;
  font-weight: 700;
  color: #0f172a;
}

.days-remaining-pill {
  font-size: 0.78rem;
  font-weight: 600;
  color: #2563eb;
  background: #eff6ff;
  padding: 2px 8px;
  border-radius: 12px;
  border: 1px solid #dbeafe;
}

.stall-count-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.82rem;
  font-weight: 700;
  color: #1e293b;
  background: #f8fafc;
  padding: 3px 10px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.stall-code-strong {
  font-weight: 700;
  color: #1d4ed8;
}

.cell-sub-meeting-point {
  font-size: 0.76rem;
  color: #64748b;
}

.phone-call-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: #eff6ff;
  color: #2563eb;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 0.78rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s ease;
}

.phone-call-btn:hover {
  background: #2563eb;
  color: #ffffff;
}

/* Card Footer Row */
.card-footer-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 14px;
  margin-top: 10px;
  border-top: 1px solid #f1f5f9;
  flex-wrap: wrap;
  gap: 12px;
}

.footer-left-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.8125rem;
  line-height: 1.45;
}

.footer-left-status strong {
  font-weight: 600;
  color: inherit;
}

.footer-right-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.btn-action-primary-solid {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #2563eb;
  border: 1px solid #2563eb;
  border-radius: 8px;
  font-size: 0.83rem;
  font-weight: 600;
  color: #ffffff;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.15s ease;
  white-space: nowrap;
  box-shadow: 0 1px 2px rgba(37, 99, 235, 0.2);
}

.btn-action-primary-solid:hover {
  background: #1d4ed8;
  border-color: #1d4ed8;
}

/* Collapsed Next Step */
.collapsed-next-step {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-top: 12px;
  padding-top: 11px;
  border-top: 1px solid #f1f5f9;
  font-size: 0.8125rem;
  line-height: 1.45;
}

.collapsed-next-step-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.collapsed-next-step-content::before {
  content: '';
  width: 7px;
  height: 7px;
  flex-shrink: 0;
  border-radius: 50%;
  background: currentColor;
}

.collapsed-next-step-label {
  color: #94a3b8;
  font-weight: 500;
  white-space: nowrap;
}

.collapsed-next-step strong {
  color: inherit;
  font-weight: 600;
}

.collapsed-actions-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.btn-collapsed-primary {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: #2563eb;
  border: 1px solid #2563eb;
  border-radius: 8px;
  font-size: 0.82rem;
  font-weight: 600;
  color: #ffffff;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.btn-collapsed-primary:hover {
  background: #1d4ed8;
  border-color: #1d4ed8;
}

.tone-waiting { color: #d97706; }
.tone-success { color: #15803d; }
.tone-action { color: #2563eb; }
.tone-muted { color: #64748b; }

/* Application Type Bar & Filter Pills */
.application-type-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 14px;
  border-bottom: 1px solid #f1f5f9;
}

.type-filter-pills {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #f8fafc;
  padding: 4px;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
}

.type-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #64748b;
  font-size: 0.8125rem; /* 13px */
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}

.type-pill:hover:not(.active) {
  color: #1e293b;
  background: rgba(255, 255, 255, 0.6);
}

.type-pill.active {
  background: #ffffff;
  color: #0f172a;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.08);
}

.type-pill-count {
  font-size: 0.75rem; /* 12px min floor */
  padding: 1px 6px;
  border-radius: 10px;
  background: #f1f5f9;
  color: #475569;
  font-weight: 700;
}

.type-pill.active .type-pill-count {
  background: #eff6ff;
  color: #2563eb;
}

.type-pill.pill-general.active .type-pill-count {
  background: #f5f3ff;
  color: #6d28d9;
}

/* Secondary Button for Stall Leasing Action */
.btn-table-stall-action {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #ffffff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
  font-size: 0.85rem;
  font-weight: 600;
  padding: 7px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  line-height: 1.3;
  text-decoration: none;
}

.btn-table-stall-action:hover {
  background: #eff6ff;
  border-color: #93c5fd;
  color: #1d4ed8;
}

/* Manage Types Button (Admin Only) */
.btn-table-manage-types {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #ffffff;
  color: #475569;
  border: 1px solid #cbd5e1;
  font-size: 0.85rem;
  font-weight: 600;
  padding: 7px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  line-height: 1.3;
  text-decoration: none;
}

.btn-table-manage-types:hover {
  background: #f8fafc;
  border-color: #94a3b8;
  color: #0f172a;
}

/* Type Badges */
.badge-type {
  font-size: 0.75rem; /* 12px min floor */
  font-weight: 600;
  padding: 2px 7px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  line-height: 1.3;
}

.badge-type-stall {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.badge-type-general {
  background: #f5f3ff;
  color: #6d28d9;
  border: 1px solid #ddd6fe;
}

.general-subject-title {
  font-weight: 600;
  color: #1e293b;
  font-size: 0.875rem;
}

.general-stall-hint {
  font-size: 0.75rem;
  color: #64748b;
  margin-top: 2px;
}

.general-no-stall-hint {
  font-size: 0.75rem;
  color: #94a3b8;
  margin-top: 2px;
}

.general-approved-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 10px;
  background: #ecfdf5;
  color: #047857;
  border: 1px solid #a7f3d0;
  border-radius: 6px;
  font-size: 0.8125rem;
  font-weight: 600;
}

.icon-box-general {
  background: #f5f3ff !important;
  color: #7c3aed !important;
}

.banner-general {
  background: #faf5ff !important;
  border-color: #f3e8ff !important;
}

.badge-general-pill {
  background: #f3e8ff;
  color: #6b21a8;
  border: 1px solid #e9d5ff;
}

/* General Application Modal Styles */
.general-modal-box {
  background: #ffffff;
  border-radius: 16px;
  width: 100%;
  max-width: 620px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.12), 0 8px 10px -6px rgba(0, 0, 0, 0.08);
  border: 1px solid #eef2f6;
  overflow: hidden;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.modal-icon-purple {
  background: #f5f3ff;
  color: #7c3aed;
}

.modal-title-general {
  font-size: 1.0625rem; /* 17px */
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  line-height: 1.3;
}

.general-modal-body {
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 22px 24px;
}

.general-alert-error {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #b91c1c;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 0.8125rem;
}

.general-preset-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  margin-top: 6px;
}

.general-preset-btn {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 0.8125rem; /* 13px */
  font-weight: 500;
  color: #475569;
  text-align: left;
  cursor: pointer;
  transition: all 0.15s ease;
  line-height: 1.35;
  font-family: inherit;
}

.general-preset-btn:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
  color: #1e293b;
}

.general-preset-btn.active {
  background: #eff6ff;
  border-color: #93c5fd;
  color: #1d4ed8;
  font-weight: 600;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-row-2col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.form-input {
  width: 100%;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 9px 12px;
  font-size: 0.875rem; /* 14px */
  color: #1e293b;
  font-family: inherit;
  outline: none;
  box-sizing: border-box;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.form-input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.form-input:disabled,
.form-input.input-disabled {
  background-color: #f1f5f9;
  color: #64748b;
  cursor: not-allowed;
  border-color: #e2e8f0;
  opacity: 0.95;
  user-select: none;
}

.form-textarea {
  width: 100%;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 0.875rem; /* 14px */
  color: #1e293b;
  font-family: inherit;
  outline: none;
  resize: vertical;
  min-height: 75px;
  box-sizing: border-box;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.form-textarea:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.text-danger {
  color: #dc2626;
}

/* Applications Cards Grid Styles */
.application-cards-container {
  margin-top: 18px;
  margin-bottom: 24px;
}

.application-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 18px;
}

.application-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
  position: relative;
}

.application-card:hover {
  border-color: #cbd5e1;
  box-shadow: 0 8px 20px -4px rgba(15, 23, 42, 0.08);
  transform: translateY(-2px);
}

.app-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.app-card-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.app-card-icon-box {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.app-card-number {
  font-size: 0.9375rem; /* 15px */
  font-weight: 700;
  display: block;
  line-height: 1.25;
}

.app-card-date {
  font-size: 0.75rem; /* 12px min floor */
  color: #94a3b8;
  margin-top: 2px;
}

.app-card-actor-row {
  margin-top: -6px;
}

.app-card-shop-banner {
  background: #f8fafc;
  border: 1px solid #f1f5f9;
  border-radius: 10px;
  padding: 12px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.app-card-shop-info {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-card-shop-name {
  font-weight: 700;
  font-size: 0.9375rem; /* 15px */
  color: #1e293b;
  line-height: 1.3;
}

.app-card-product-type {
  font-size: 0.8125rem; /* 13px */
  color: #64748b;
  margin-top: 2px;
}

.app-card-stall-badge {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  flex-shrink: 0;
}

.app-card-stall-rate {
  font-size: 0.75rem; /* 12px min floor */
  color: #64748b;
  font-weight: 500;
  white-space: nowrap;
}

.app-card-meta-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  padding: 0 2px;
}

.meta-label {
  font-size: 0.75rem; /* 12px min floor */
  color: #64748b;
  margin-bottom: 2px;
  display: block;
}

.meta-val-primary {
  font-size: 0.875rem; /* 14px */
  font-weight: 600;
  color: #1e293b;
  line-height: 1.3;
}

.meta-val-sub {
  font-size: 0.8125rem; /* 13px */
  color: #64748b;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 2px;
}

.meta-val-date {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: #334155;
}

.app-card-note-section {
  display: flex;
  flex-direction: column;
  gap: 6px;
  border-top: 1px dashed #f1f5f9;
  padding-top: 10px;
}

.app-card-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
  border-top: 1px solid #f1f5f9;
  padding-top: 14px;
  margin-top: auto;
}

.btn-card-detail {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  font-size: 0.8125rem; /* 13px */
  font-weight: 500;
  color: #334155;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}

.btn-card-detail:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
  color: #0f172a;
}

.empty-cards-state {
  text-align: center;
  padding: 56px 20px;
  color: #94a3b8;
  background: #f8fafc;
  border: 1px dashed #cbd5e1;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.empty-cards-title {
  font-size: 0.9375rem; /* 15px */
  color: #64748b;
  font-weight: 500;
}

/* Actor Tag in สถานะพิจารณา */
.action-actor-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  margin-top: 6px;
  font-size: 0.75rem; /* 12px min floor */
  font-weight: 500;
  padding: 3px 8px;
  border-radius: 6px;
  line-height: 1.3;
}

.actor-tag-officer-pending {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.actor-tag-officer {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.actor-tag-director {
  background: #eef2ff;
  color: #4338ca;
  border: 1px solid #c7d2fe;
}

.actor-tag-success {
  background: #ecfdf5;
  color: #047857;
  border: 1px solid #a7f3d0;
}

.actor-tag-warning {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
}

.actor-tag-danger {
  background: #fef2f2;
  color: #b91c1c;
  border: 1px solid #fecaca;
}

.actor-tag-default {
  background: #f8fafc;
  color: #64748b;
  border: 1px solid #e2e8f0;
}

/* Comment Box in หมายเหตุการพิจารณา */
.review-comment-container {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.review-comment-pill {
  border-radius: 8px;
  padding: 6px 10px;
  font-size: 0.8125rem; /* 13px */
  line-height: 1.4;
  word-break: break-word;
  transition: all 0.15s ease;
}

.comment-officer {
  background: #f0f7ff;
  border: 1px solid #dbeafe;
  color: #1e3a8a;
}

.comment-director {
  background: #f5f3ff;
  border: 1px solid #ede9fe;
  color: #3730a3;
}

.comment-success {
  background: #f0fdf4;
  border: 1px solid #dcfce7;
  color: #166534;
}

.comment-warning {
  background: #fffdf5;
  border: 1px solid #fef3c7;
  color: #92400e;
}

.comment-danger {
  background: #fff5f5;
  border: 1px solid #fee2e2;
  color: #991b1b;
}

.comment-info {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #334155;
}

.comment-header-row {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 2px;
}

.comment-role-label {
  font-weight: 600;
  font-size: 0.75rem; /* 12px */
  letter-spacing: 0.01em;
}

.comment-text-body {
  font-size: 0.8125rem; /* 13px */
  color: inherit;
  opacity: 0.95;
}

.more-comments-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem; /* 12px */
  color: #2563eb;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  transition: background 0.15s ease;
  width: fit-content;
}

.more-comments-badge:hover {
  background: #eff6ff;
  text-decoration: underline;
}

.text-empty-dash {
  color: #94a3b8;
  font-size: 0.875rem;
  padding-left: 4px;
}

/* Reject Modal Styles */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.48);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 16px;
}

.reject-modal-box {
  background: #ffffff;
  border-radius: 16px;
  width: 100%;
  max-width: 520px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.12), 0 8px 10px -6px rgba(0, 0, 0, 0.08);
  border: 1px solid #eef2f6;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  border-bottom: 1px solid #f1f5f9;
}

.modal-header-title-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.modal-icon-badge {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.modal-icon-danger {
  background: #fee2e2;
  color: #dc2626;
}

.modal-title-danger {
  font-size: 1.0625rem; /* 17px */
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  line-height: 1.3;
}

.modal-subtitle {
  font-size: 0.8125rem; /* 13px */
  color: #64748b;
  margin-top: 2px;
}

.btn-modal-close {
  background: transparent;
  border: none;
  cursor: pointer;
  color: #94a3b8;
  padding: 6px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.btn-modal-close:hover {
  background: #f1f5f9;
  color: #334155;
}

.modal-body {
  padding: 20px 24px;
}

.reject-summary-banner {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 10px 14px;
  margin-bottom: 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 0.8125rem; /* 13px */
}

.summary-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.summary-label {
  color: #64748b;
  font-weight: 500;
}

.summary-value {
  color: #1e293b;
}

.reject-alert-box {
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 0.78125rem; /* 12.5px */
  color: #991b1b;
  margin-bottom: 14px;
  display: flex;
  align-items: flex-start;
  gap: 8px;
  line-height: 1.4;
}

.form-label-sub {
  font-size: 0.75rem; /* 12px */
  font-weight: 600;
  color: #64748b;
  margin-bottom: 6px;
  display: block;
}

.quick-preset-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 14px;
}

.quick-preset-chip {
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  color: #475569;
  font-size: 0.75rem; /* 12px */
  padding: 4px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
  line-height: 1.3;
}

.quick-preset-chip:hover {
  background: #e2e8f0;
  color: #0f172a;
  border-color: #cbd5e1;
}

.form-label {
  display: block;
  font-size: 0.8125rem; /* 13px */
  font-weight: 600;
  color: #334155;
  margin-bottom: 6px;
}

.reject-textarea {
  width: 100%;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 0.875rem; /* 14px */
  font-family: inherit;
  color: #1e293b;
  resize: vertical;
  min-height: 80px;
  box-sizing: border-box;
  transition: border-color 0.15s, box-shadow 0.15s;
  outline: none;
}

.reject-textarea:focus {
  border-color: #dc2626;
  box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.1);
}

.reject-textarea.has-error {
  border-color: #ef4444;
}

.reject-error-text {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #dc2626;
  font-size: 0.75rem; /* 12px */
  margin-top: 5px;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 24px;
  background: #f8fafc;
  border-top: 1px solid #f1f5f9;
}
.preset-label-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.btn-manage-types-shortcut {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 0.75rem; /* 12px min floor */
  font-weight: 600;
  color: #2563eb;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  padding: 3px 8px;
  border-radius: 6px;
  text-decoration: none;
  transition: all 0.15s ease;
}

.btn-manage-types-shortcut:hover {
  background: #dbeafe;
  color: #1d4ed8;
  border-color: #93c5fd;
}

/* ==========================================================================
   Tenant Revision Popup Modal Styles (Replicating LeaseApplicationModal.vue)
   ========================================================================== */
.revision-modal-content {
  max-width: 680px;
  width: 95%;
  border-radius: 16px;
  background: #ffffff;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.revision-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.revision-icon-badge {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.revision-modal-title {
  font-size: 1.125rem; /* 18px */
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  line-height: 1.3;
}

.revision-modal-subtitle {
  font-size: 0.8125rem; /* 13px */
  color: #64748b;
  margin-top: 3px;
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.revision-modal-subtitle strong {
  color: #2563eb;
  font-weight: 600;
}

.revision-sub-dot {
  color: #cbd5e1;
}

.revision-modal-body {
  max-height: calc(85vh - 140px);
  overflow-y: auto;
  padding: 20px 24px;
}

/* Officer Note Alert Banner */
.revision-officer-note-card {
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 10px;
  padding: 12px 14px;
  margin-bottom: 16px;
}

.revision-note-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 700;
  font-size: 0.8125rem; /* 13px */
  color: #b45309;
  margin-bottom: 4px;
}

.revision-note-body {
  font-size: 0.875rem; /* 14px */
  color: #92400e;
  line-height: 1.5;
  word-break: break-word;
}

/* Stall Financial Summary Banner */
.revision-stall-summary-banner {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 10px;
  padding: 12px 14px;
  margin-bottom: 16px;
}

.summary-banner-inner {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  font-size: 0.85rem;
  color: #1e3a8a;
}

.summary-stall-col {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.summary-market-tag {
  color: #3b82f6;
  font-size: 0.8125rem;
}

.summary-rates-wrap {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
  align-items: center;
}

.summary-total-highlight {
  color: #1d4ed8;
  font-size: 0.9375rem;
  font-weight: 700;
}

/* Form Section Headings & Inputs */
.form-section-heading {
  font-size: 0.9375rem; /* 15px */
  font-weight: 700;
  color: #2563eb;
  margin: 18px 0 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  border-bottom: 1px solid #f1f5f9;
  padding-bottom: 6px;
}

.form-textarea-revision {
  box-sizing: border-box;
  width: 100%;
  min-height: 60px;
  padding: 9px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  background: #fff;
  color: #1e293b;
  font-family: inherit;
  font-size: 0.875rem; /* 14px */
  line-height: 1.5;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
  resize: vertical;
}

.form-textarea-revision:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
}

.revision-form-error-alert {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  padding: 10px 14px;
  color: #dc2626;
  font-size: 0.8125rem; /* 13px */
  font-weight: 500;
  margin-top: 14px;
}

.revision-modal-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  padding: 14px 24px;
  background: #f8fafc;
  border-top: 1px solid #f1f5f9;
}

/* Upload styles for Revision Modal */
.upload-section-container {
  margin-top: 6px;
  margin-bottom: 8px;
}

.upload-header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  flex-wrap: wrap;
  gap: 6px;
}

.upload-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.upload-badge-spec {
  font-size: 0.75rem;
  color: #64748b;
  background: #f1f5f9;
  padding: 3px 8px;
  border-radius: 6px;
  font-weight: 500;
}

.id-card-dropzone {
  display: flex;
  align-items: center;
  gap: 14px;
  background: #f8fafc;
  border: 1.5px dashed #cbd5e1;
  border-radius: 12px;
  padding: 14px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.id-card-dropzone:hover,
.id-card-dropzone.is-drag-over {
  background: #eff6ff;
  border-color: #2563eb;
  transform: translateY(-1px);
}

.id-card-dropzone.has-error {
  background: #fef2f2;
  border-color: #f87171;
}

.dropzone-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.id-card-dropzone:hover .dropzone-icon-box {
  background: #2563eb;
  color: #ffffff;
}

.dropzone-prompt {
  flex: 1;
  min-width: 0;
}

.prompt-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 2px;
}

.prompt-link {
  color: #2563eb;
  text-decoration: underline;
  text-underline-offset: 2px;
}

.prompt-desc {
  font-size: 0.8125rem;
  color: #64748b;
  line-height: 1.4;
}

.dropzone-cta-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  color: #334155;
  font-size: 0.8125rem;
  font-weight: 600;
  padding: 7px 13px;
  border-radius: 8px;
  flex-shrink: 0;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.id-card-preview-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #ffffff;
  border: 1px solid #bfdbfe;
  border-radius: 12px;
  padding: 10px 14px;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.05);
}

.preview-thumb-wrap {
  width: 46px;
  height: 46px;
  border-radius: 8px;
  overflow: hidden;
  background: #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid #e2e8f0;
}

.preview-thumb-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-pdf-badge {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #dc2626;
  background: #fee2e2;
  width: 100%;
  height: 100%;
  gap: 1px;
}

.preview-pdf-badge span {
  font-size: 0.75rem;
  font-weight: 700;
  line-height: 1;
}

.preview-details {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.preview-filename {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.preview-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.75rem;
}

.meta-size {
  color: #64748b;
  font-weight: 500;
}

.meta-check {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #16a34a;
  font-weight: 600;
  background: #ecfdf5;
  padding: 1px 7px;
  border-radius: 6px;
}

.preview-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.btn-preview-change {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #2563eb;
  font-size: 0.8125rem;
  font-weight: 600;
  padding: 6px 10px;
  border-radius: 7px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-preview-change:hover {
  background: #dbeafe;
}

.btn-preview-delete {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: #ffffff;
  border: 1px solid #fecaca;
  color: #ef4444;
  border-radius: 7px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-preview-delete:hover {
  background: #fee2e2;
}

.upload-error-msg {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  font-size: 0.8125rem;
  color: #dc2626;
  background: #fef2f2;
  border: 1px solid #fecaca;
  padding: 6px 10px;
  border-radius: 7px;
}
</style>

