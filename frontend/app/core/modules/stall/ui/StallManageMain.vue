<template>
  <div>
    <AppHeader />
    <main
      class="content-body manage-page"
      :class="{ 'is-map-page': isMapOverviewMode, 'has-left-panel': isMapOverviewMode && !!leftPanelWidth }"
      :style="isMapOverviewMode && leftPanelWidth ? { '--lp-w': leftPanelWidth } : undefined"
    >
      <!-- สวิตช์มุมมองแผนที่ (ข้อมูลแผง / การชำระเงิน / การจดมิเตอร์) แบบ Capsule Pill ลอยมุมซ้ายบน -->
      <div v-if="isMapOverviewMode" ref="viewDropdownWrapRef" class="map-mode-dropdown-wrap">
        <button
          type="button"
          class="map-mode-trigger-btn"
          :class="{ 'is-open': isViewDropdownOpen }"
          aria-haspopup="listbox"
          :aria-expanded="isViewDropdownOpen"
          @click="isViewDropdownOpen = !isViewDropdownOpen"
        >
          <span class="trigger-icon-box" :class="`tone-${currentModeOption.tone}`" aria-hidden="true">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" v-html="currentModeOption.svgIcon"></svg>
          </span>
          <span class="trigger-text-group">
            <strong class="trigger-title">{{ currentModeOption.label }}</strong>
          </span>
          <svg class="trigger-chevron" :class="{ 'is-open': isViewDropdownOpen }" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>

        <!-- Dropdown Menu List -->
        <transition name="dropdown-pop-fade">
          <div v-if="isViewDropdownOpen" class="map-mode-dropdown-menu" role="listbox" aria-label="เลือกมุมมองแผนที่">
            <button
              v-for="opt in viewModeOptions"
              :key="opt.value"
              type="button"
              class="map-mode-menu-item"
              :class="{ 'is-selected': pageMode === opt.value }"
              role="option"
              :aria-selected="pageMode === opt.value"
              @click="selectPageMode(opt.value)"
            >
              <span class="item-icon-box" :class="`tone-${opt.tone}`" aria-hidden="true">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" v-html="opt.svgIcon"></svg>
              </span>
              <span class="item-text-group">
                <strong class="item-title">{{ opt.label }}</strong>
                <small class="item-desc">{{ opt.desc }}</small>
              </span>
              <span v-if="pageMode === opt.value" class="item-check" aria-hidden="true">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </span>
            </button>
          </div>
        </transition>
      </div>

      <AppBreadcrumb
        v-if="!isStaff && !isMapOverviewMode"
        :market="marketName"
        :marketId="marketId"
        :items="breadcrumbItems"
        :current="breadcrumbCurrent"
        :backTo="`/stalls?market=${marketId}`"
        :onBack="handleBreadcrumbBack"
      />

      <!-- Page Heading มาตรฐานเดียวกับระบบ สำหรับโหมดจัดการผัง -->
      <header v-if="!isInspectMode && !isMapOverviewMode" class="manage-page-heading">
        <div class="heading-left">
          <h1 class="page-title">จัดการผัง GIS</h1>
          <p class="page-desc">จัดการโซน ตำแหน่งแผงค้า และข้อมูลผังแผงดิจิทัล (GIS)</p>
        </div>
        <div class="manage-page-actions">
          <template v-if="isManageDataLoading">
            <span class="header-status-skeleton" aria-label="กำลังโหลดสถานะผัง"></span>
            <span class="header-action-skeleton" aria-hidden="true"></span>
          </template>
          <template v-else>
            <span class="draft-status" :class="{ published: !hasUnsavedChanges && gisPublished, 'is-dirty': hasUnsavedChanges }">
              <i></i>
              {{ hasUnsavedChanges ? 'มีการแก้ไขที่ยังไม่บันทึก' : (gisPublished ? 'บันทึกและเผยแพร่แล้ว' : 'ไม่มีการแก้ไข') }}
            </span>
            <AppButton
              size="md"
              :variant="hasUnsavedChanges || !gisPublished ? 'primary' : 'outline'"
              :disabled="gisPublished && !hasUnsavedChanges"
              :loading="savingGIS"
              @click="confirmAndSaveGIS"
            >
              <template #prefix>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
                  <polyline points="17 21 17 13 7 13 7 21"></polyline>
                  <polyline points="7 3 7 8 15 8"></polyline>
                </svg>
              </template>
              {{ hasUnsavedChanges ? 'บันทึกการแก้ไข' : (gisPublished ? 'บันทึกการแก้ไข' : 'เผยแพร่ผัง') }}
            </AppButton>
          </template>
        </div>
      </header>

      <section v-if="isInspectDashboardMode" class="staff-dashboard-intro" aria-labelledby="staff-dashboard-title">
        <div class="dashboard-intro-left">
          <h1 id="staff-dashboard-title">แดชบอร์ด</h1>
          <p class="dashboard-updated-sub">ข้อมูลอัปเดตล่าสุด ณ วันที่ {{ dashboardUpdatedDateText }}</p>
        </div>
      </section>

      <!-- ================================================================ -->
      <!-- SUMMARY DASHBOARD VIEW (สรุปข้อมูลของจริงสำหรับเจ้าหน้าที่ตลาด)     -->
      <!-- ================================================================ -->


      <!-- ฟังก์ชันการจดมิเตอร์น้ำ-ไฟ และการออกบิลแผงค้า (จดมิเตอร์ / ออกบิล / กรอกตัวเลข) -->
      <template v-if="isMapOverviewMode && (pageMode === 'meters' || pageMode === 'payments')">
        <button
          v-if="!costTool && pageMode === 'meters'"
          type="button"
          class="map-meter-float-btn"
          title="เปิดตารางจดมิเตอร์น้ำ-ไฟ"
          @click="costTool = 'utility'"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m12 14 4-4"></path>
            <path d="M3.34 19a10 10 0 1 1 17.32 0"></path>
          </svg>
          <span>ตารางจดมิเตอร์น้ำ-ไฟ</span>
        </button>

        <button
          v-if="!costTool && pageMode === 'payments'"
          type="button"
          class="map-meter-float-btn"
          title="เปิดรายการบิลแผงค้า"
          @click="costTool = 'bill'"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 2v20l2-1 2 1 2-1 2 1 2-1 2 1 2-1 2 1V2l-2 1-2-1-2 1-2-1-2 1-2-1-2 1z"></path>
            <path d="M8 8h8M8 12h8M8 16h5"></path>
          </svg>
          <span>รายการบิลแผงค้า</span>
        </button>

        <div v-if="costTool" class="map-cost-float" role="dialog" :aria-label="costTool === 'utility' ? 'รายการค่าน้ำ ค่าไฟ' : 'รายการบิลแผงเช่าชั่วคราว'">
          <StallUtilityBillingPanel
            v-if="costTool === 'utility'"
            :filter="costFilter"
            :active-code="costDrawerOpen ? costStallCode : ''"
            @close="closeCostTool"
            @select="openCostDrawer($event, 'utility')"
          />
          <StallTempBillPanel
            v-else
            :filter="costFilter"
            :active-code="costDrawerOpen ? costStallCode : ''"
            @close="closeCostTool"
            @select="openCostDrawer($event, 'temp')"
            @new="openCostDrawer('', 'temp')"
          />
        </div>

        <StallCostDrawer
          v-if="costDrawerOpen"
          :key="`${costDrawerKind}-${costStallCode}`"
          :kind="costDrawerKind"
          :code="costStallCode"
          :lookup="lookupCostStall"
          @close="closeCostDrawer"
        />
      </template>

      <section
        v-if="isInspectDashboardMode"
        id="map-dashboard-panel"
        class="officer-summary-dashboard"
        aria-label="สรุปข้อมูลแดชบอร์ดตลาด"
      >
        <!-- มุมมอง 1: ข้อมูลแผงและสัญญา (Stalls View) -->
        <template v-if="pageMode === 'stalls'">
        <div class="map-summary-grid dash-kpi-ribbon">
          <div class="map-summary-card summary-reserved">
            <span class="map-summary-icon" aria-hidden="true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"></circle><path d="M12 7v5l3 2"></path></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">การจองจริง (รออนุมัติ {{ realActionCounts.pendingBookings }})</small>
              <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                <span class="kpi-skeleton-num"></span>
                <span class="kpi-skeleton-pill"></span>
              </div>
              <span v-else class="map-summary-value-row">
                <strong class="map-summary-count">{{ realActionCounts.realActiveBookings.toLocaleString('th-TH') }}</strong>
                <em class="map-summary-unit-pill">ใบจอง</em>
              </span>
            </span>
          </div>

          <div class="map-summary-card summary-total">
            <span class="map-summary-icon" aria-hidden="true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">รายการรอคิว (Waiting List)</small>
              <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                <span class="kpi-skeleton-num"></span>
                <span class="kpi-skeleton-pill"></span>
              </div>
              <span v-else class="map-summary-value-row">
                <strong class="map-summary-count">{{ realActionCounts.waitlistCount.toLocaleString('th-TH') }}</strong>
                <em class="map-summary-unit-pill">คิว</em>
              </span>
            </span>
          </div>

          <div class="map-summary-card summary-available">
            <span class="map-summary-icon" aria-hidden="true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">คำร้องรอตรวจ / สัญญาที่มีผล ({{ marketContracts.length }})</small>
              <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                <span class="kpi-skeleton-num"></span>
                <span class="kpi-skeleton-pill"></span>
              </div>
              <span v-else class="map-summary-value-row">
                <strong class="map-summary-count">{{ realActionCounts.pendingOfficerApps.toLocaleString('th-TH') }}</strong>
                <em class="map-summary-unit-pill">รายการ</em>
              </span>
            </span>
          </div>

          <div class="map-summary-card summary-occupied" style="cursor: default;">
            <span class="map-summary-icon" aria-hidden="true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="5" width="20" height="14" rx="2"></rect><line x1="2" y1="10" x2="22" y2="10"></line></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">มูลค่าค่าเช่ารวม/เดือน</small>
              <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                <span class="kpi-skeleton-num" style="width: 80px;"></span>
              </div>
              <span v-else class="map-summary-value-row">
                <strong class="map-summary-count">฿{{ realMarketSummary.totalMonthlyRent.toLocaleString('th-TH') }}</strong>
              </span>
            </span>
          </div>
        </div>

        <!-- 2. Two Visual Chart Cards: Donut Chart (สัดส่วนสถานะพื้นที่) & Bar Chart (สัดส่วนตามหมวดสินค้า) -->
        <div class="dash-two-col-grid">
          <!-- Left Card: Donut Chart สัดส่วนสถานะแผงและมูลค่าค่าเช่า -->
          <article class="dash-system-card">
            <div class="dash-card-header">
              <div>
                <h2 class="dash-card-title">สัดส่วนสถานะพื้นที่</h2>
                <p class="dash-card-subtitle">{{ realMarketSummary.zonesCount }} โซน · เฉลี่ย ฿{{ realMarketSummary.avgRent.toLocaleString('th-TH') }}/แผง</p>
              </div>
            </div>

            <div class="dash-donut-layout">
              <div class="dash-donut-wrap">
                <svg viewBox="0 0 140 140" class="dash-donut-svg" aria-hidden="true">
                  <circle cx="70" cy="70" r="54" fill="none" stroke="#f1f5f9" stroke-width="16" />
                  <circle
                    v-for="seg in donutSegments"
                    :key="seg.key"
                    cx="70"
                    cy="70"
                    r="54"
                    fill="none"
                    :stroke="seg.color"
                    stroke-width="16"
                    :stroke-dasharray="`${seg.dash} ${donutCircumference}`"
                    :stroke-dashoffset="-seg.offset"
                    transform="rotate(-90 70 70)"
                  />
                </svg>
                <div class="dash-donut-center">
                  <strong>{{ realMarketSummary.total.toLocaleString('th-TH') }}</strong>
                  <span>แผงทั้งหมด</span>
                </div>
              </div>

              <div class="dash-donut-legend-list">
                <div class="dash-legend-row">
                  <span class="legend-left"><i class="dot-available"></i>ว่าง</span>
                  <div class="legend-right">
                    <strong>{{ realMarketSummary.available.toLocaleString('th-TH') }}</strong>
                    <span class="dash-status-badge badge-green">{{ realMarketSummary.availablePct }}%</span>
                  </div>
                </div>
                <div class="dash-legend-row">
                  <span class="legend-left"><i class="dot-reserved"></i>ติดจอง</span>
                  <div class="legend-right">
                    <strong>{{ realMarketSummary.reserved.toLocaleString('th-TH') }}</strong>
                    <span class="dash-status-badge badge-amber">{{ realMarketSummary.reservedPct }}%</span>
                  </div>
                </div>
                <div class="dash-legend-row">
                  <span class="legend-left"><i class="dot-occupied"></i>มีผู้เช่า</span>
                  <div class="legend-right">
                    <strong>{{ realMarketSummary.occupied.toLocaleString('th-TH') }}</strong>
                    <span class="dash-status-badge badge-gray">{{ realMarketSummary.occupancyPct }}%</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="dash-mini-rev-row">
              <div class="dash-mini-rev-item">
                <span>ค่าเช่าปัจจุบัน</span>
                <strong class="is-blue">฿{{ realMarketSummary.occupiedMonthlyRent.toLocaleString('th-TH') }}</strong>
              </div>
              <div class="dash-mini-rev-item">
                <span>มูลค่าแผงว่าง</span>
                <strong class="is-green">฿{{ realMarketSummary.availableMonthlyRent.toLocaleString('th-TH') }}</strong>
              </div>
            </div>
          </article>

          <!-- Right Card: Horizontal Bar Chart สัดส่วนแผงแยกตามหมวดสินค้า -->
          <article class="dash-system-card">
            <div class="dash-card-header">
              <div>
                <h2 class="dash-card-title">หมวดสินค้า</h2>
                <p class="dash-card-subtitle">{{ categoryChartRows.length }} หมวดในตลาด</p>
              </div>
            </div>

            <div class="dash-bar-chart-list">
              <div v-for="cat in categoryChartRows" :key="cat.name" class="dash-bar-row">
                <div class="dash-bar-top">
                  <span class="dash-bar-name">
                    <i :style="{ backgroundColor: cat.color }"></i>
                    <span :title="cat.name">{{ compactCategoryLabel(cat.name) }}</span>
                  </span>
                  <span class="dash-bar-meta">
                    <strong>{{ cat.total.toLocaleString('th-TH') }}</strong> · {{ cat.pct }}%
                  </span>
                </div>
                <div class="dash-bar-track">
                  <span class="dash-bar-fill" :style="{ width: `${cat.barPct}%`, backgroundColor: cat.color }"></span>
                </div>
              </div>
            </div>
          </article>
        </div>

        <!-- 3. Zone / Project Breakdown Table (มาตรฐาน Table Design Standards ของระบบ) -->
        <article class="dash-system-card dash-table-card">
          <div class="dash-table-header">
            <div class="dash-table-header-left">
              <h2 class="dash-card-title">ตารางสรุปรายโซน</h2>
              <p class="dash-card-subtitle">{{ filteredZoneRows.length }} โซน</p>
            </div>
            <div class="dash-table-header-actions">
              <div class="dash-category-select">
                <AppSelect
                  v-model="dashZoneCategoryFilter"
                  :options="dashCategoryFilterOptions"
                  placeholder="ทุกหมวดสินค้า"
                />
              </div>
              <div class="dash-search-box">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
                <AppInput
                  v-model="dashZoneSearchQuery"
                  placeholder="ค้นหาโซน..."
                />
              </div>
            </div>
          </div>

          <div class="dash-table-wrap">
            <table class="dash-std-table">
              <thead>
                <tr>
                  <th>โซน</th>
                  <th>หมวดสินค้า</th>
                  <th class="text-right">ทั้งหมด</th>
                  <th class="text-center">ว่าง</th>
                  <th class="text-center">ติดจอง</th>
                  <th class="text-center">มีผู้เช่า</th>
                  <th style="width: 150px;">อัตราเช่า</th>
                  <th class="text-right">ค่าเช่าเฉลี่ย</th>
                  <th class="text-right">จัดการ</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="zone in paginatedZoneRows" :key="zone.id">
                  <td>
                    <div class="zone-name-cell">
                      <span class="zone-color-dot" :style="{ backgroundColor: zone.color }"></span>
                      <div>
                        <strong>{{ zone.name }}</strong>
                        <small v-if="zone.floorsCount > 1">{{ zone.floorsCount }} ชั้น</small>
                      </div>
                    </div>
                  </td>
                  <td>{{ zone.category }}</td>
                  <td class="text-right"><strong>{{ zone.total.toLocaleString('th-TH') }}</strong></td>
                  <td class="text-center">
                    <span class="dash-status-badge badge-green">{{ zone.available.toLocaleString('th-TH') }}</span>
                  </td>
                  <td class="text-center">
                    <span class="dash-status-badge" :class="zone.reserved > 0 ? 'badge-amber' : 'badge-gray'">{{ zone.reserved.toLocaleString('th-TH') }}</span>
                  </td>
                  <td class="text-center">
                    <span class="dash-status-badge badge-gray">{{ zone.occupied.toLocaleString('th-TH') }}</span>
                  </td>
                  <td>
                    <div class="zone-occ-cell">
                      <div class="zone-occ-bar">
                        <span :style="{ width: `${zone.occupancyPct}%` }"></span>
                      </div>
                      <strong>{{ zone.occupancyPct }}%</strong>
                    </div>
                  </td>
                  <td class="text-right"><strong>฿{{ zone.avgMonthlyRent.toLocaleString('th-TH') }}</strong></td>
                  <td class="text-right">
                    <button type="button" class="dash-table-view-btn" @click="openZoneOnMap(zone.id)">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                      <span>ดูผัง</span>
                    </button>
                  </td>
                </tr>
                <tr v-if="filteredZoneRows.length === 0">
                  <td colspan="9" class="dash-empty-row">ไม่พบข้อมูลโซน</td>
                </tr>
              </tbody>
            </table>
          </div>

          <AppPagination
            v-if="filteredZoneRows.length > 0"
            v-model="dashZonePage"
            :total-items="filteredZoneRows.length"
            :items-per-page="dashZonePageSize"
            :page-size-options="[5, 10, 20, 50]"
            unit-label="โซน"
            @update:items-per-page="val => { dashZonePageSize = val; dashZonePage = 1 }"
          />
        </article>
      </template>

      <!-- มุมมอง 2: การเงินและการจัดเก็บ (Payments View) -->
      <template v-else-if="pageMode === 'payments'">
        <!-- 1. Four Financial Metric Cards -->
        <div class="dash-kpi-ribbon" style="display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px;">
          <div class="map-summary-card summary-total">
            <span class="map-summary-icon" style="background: #eff6ff; color: #2563eb;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="5" width="20" height="14" rx="2"></rect><line x1="2" y1="10" x2="22" y2="10"></line></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ยอดเรียกเก็บประจำเดือน</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem;">฿284,500</strong>
                <em class="map-summary-unit-pill" style="background: #eff6ff; color: #2563eb;">42 บิล</em>
              </span>
              <span class="map-summary-meta">อัตราเฉลี่ย ฿6,773 / แผง</span>
            </span>
          </div>

          <div class="map-summary-card summary-available">
            <span class="map-summary-icon" style="background: #ecfdf5; color: #059669;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">จัดเก็บสำเร็จแล้ว</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: #059669;">฿192,000</strong>
                <em class="map-summary-unit-pill" style="background: #ecfdf5; color: #059669;">67.5%</em>
              </span>
              <span class="map-summary-meta">ชำระแล้ว 28 จาก 42 รายการ</span>
            </span>
          </div>

          <div class="map-summary-card summary-reserved">
            <span class="map-summary-icon" style="background: #fffbeb; color: #d97706;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">รอชำระในกำหนด</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: #d97706;">฿58,500</strong>
                <em class="map-summary-unit-pill" style="background: #fffbeb; color: #d97706;">9 บิล</em>
              </span>
              <span class="map-summary-meta">ครบกำหนดสิ้นเดือนนี้</span>
            </span>
          </div>

          <div class="map-summary-card summary-action-required">
            <span class="map-summary-icon" style="background: #fef2f2; color: #dc2626;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ค้างชำระเกินกำหนด</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: #dc2626;">฿34,000</strong>
                <em class="map-summary-unit-pill" style="background: #fee2e2; color: #dc2626;">5 บิล</em>
              </span>
              <span class="map-summary-meta">ต้องออกหนังสือติดตามหนี้</span>
            </span>
          </div>
        </div>

        <!-- 2. Breakdown Section -->
        <div class="dash-card">
          <h3 class="dash-card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><path d="M12 2a10 10 0 0 1 10 10H12V2z"></path></svg>
            โครงสร้างรายรับแยกตามประเภทค่าบริการ
          </h3>
          <div class="dash-breakdown-list">
            <div class="dash-breakdown-row">
              <div class="dash-breakdown-label">
                <span class="dot-indicator" style="background: #2563eb;"></span>
                <span>ค่าเช่าแผงรายเดือน</span>
              </div>
              <div class="dash-breakdown-val">
                <strong>฿226,800</strong>
                <small>79.7%</small>
              </div>
            </div>
            <div class="dash-progress-bar"><div class="dash-progress-fill" style="width: 79.7%; background: #2563eb;"></div></div>

            <div class="dash-breakdown-row">
              <div class="dash-breakdown-label">
                <span class="dot-indicator" style="background: #0284c7;"></span>
                <span>ค่ากระแสไฟฟ้า</span>
              </div>
              <div class="dash-breakdown-val">
                <strong>฿32,450</strong>
                <small>11.4%</small>
              </div>
            </div>
            <div class="dash-progress-bar"><div class="dash-progress-fill" style="width: 11.4%; background: #0284c7;"></div></div>

            <div class="dash-breakdown-row">
              <div class="dash-breakdown-label">
                <span class="dot-indicator" style="background: #0d9488;"></span>
                <span>ค่าน้ำประปา</span>
              </div>
              <div class="dash-breakdown-val">
                <strong>฿15,250</strong>
                <small>5.4%</small>
              </div>
            </div>
            <div class="dash-progress-bar"><div class="dash-progress-fill" style="width: 5.4%; background: #0d9488;"></div></div>

            <div class="dash-breakdown-row">
              <div class="dash-breakdown-label">
                <span class="dot-indicator" style="background: #8b5cf6;"></span>
                <span>ค่าบริการส่วนกลางและเก็บขยะ</span>
              </div>
              <div class="dash-breakdown-val">
                <strong>฿10,000</strong>
                <small>3.5%</small>
              </div>
            </div>
            <div class="dash-progress-bar"><div class="dash-progress-fill" style="width: 3.5%; background: #8b5cf6;"></div></div>
          </div>
        </div>

        <!-- 3. Top Overdue Stalls Table -->
        <div class="dash-card">
          <div class="dash-card-head-row">
            <h3 class="dash-card-title">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#dc2626" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
              รายการแผงค้างชำระที่ต้องติดตามด่วน (5 แผง)
            </h3>
          </div>
          <div class="dash-table-wrap">
            <table class="dash-mini-table">
              <thead>
                <tr>
                  <th>แผง</th>
                  <th>ผู้เช่า</th>
                  <th style="text-align: right;">ยอดค้าง</th>
                  <th style="text-align: center;">สถานะ</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="st in overdueStallTableItems" :key="st.code">
                  <td>
                    <span class="table-stall-pill">{{ st.code }}</span>
                  </td>
                  <td>
                    <div class="table-tenant-name">{{ st.tenant }}</div>
                    <small class="table-tenant-sub">{{ st.zone }}</small>
                  </td>
                  <td style="text-align: right; font-weight: 700; color: #dc2626;">
                    {{ st.amountText }}
                  </td>
                  <td style="text-align: center;">
                    <span class="table-status-pill pill-overdue">{{ st.daysText }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </template>

      <!-- มุมมอง 3: การจดมิเตอร์น้ำและไฟฟ้า (Meters View) -->
      <template v-else-if="pageMode === 'meters'">
        <!-- 1. Four Meter Progress Cards -->
        <div class="dash-kpi-ribbon" style="display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px;">
          <div class="map-summary-card summary-total">
            <span class="map-summary-icon" style="background: #eff6ff; color: #2563eb;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m12 14 4-4"></path><path d="M3.34 19a10 10 0 1 1 17.32 0"></path></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ความคืบหน้าการจด</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: #2563eb;">83.3%</strong>
                <em class="map-summary-unit-pill" style="background: #eff6ff; color: #2563eb;">35/42 จุด</em>
              </span>
              <span class="map-summary-meta">เหลือ 7 จุดที่ต้องจด</span>
            </span>
          </div>

          <div class="map-summary-card summary-available">
            <span class="map-summary-icon" style="background: #ecfdf5; color: #059669;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2.69l5.66 5.66a8 8 0 1 1-11.31 0z"></path></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ปริมาณการใช้น้ำรวม</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: #059669;">1,428</strong>
                <em class="map-summary-unit-pill" style="background: #ecfdf5; color: #059669;">ลบ.ม.</em>
              </span>
              <span class="map-summary-meta">เฉลี่ย 34 ลบ.ม. / แผง</span>
            </span>
          </div>

          <div class="map-summary-card summary-reserved">
            <span class="map-summary-icon" style="background: #fffbeb; color: #d97706;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ปริมาณการใช้ไฟฟ้ารวม</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: #d97706;">12,850</strong>
                <em class="map-summary-unit-pill" style="background: #fffbeb; color: #d97706;">หน่วย (kWh)</em>
              </span>
              <span class="map-summary-meta">เฉลี่ย 306 หน่วย / แผง</span>
            </span>
          </div>

          <div class="map-summary-card summary-action-required">
            <span class="map-summary-icon" style="background: #f5f3ff; color: #7c3aed;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ตรวจพบค่าผิดปกติ</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: #7c3aed;">2</strong>
                <em class="map-summary-unit-pill" style="background: #f5f3ff; color: #7c3aed;">จุด</em>
              </span>
              <span class="map-summary-meta">น้ำพุ่งสูง 1 • ไฟไม่เดิน 1</span>
            </span>
          </div>
        </div>

        <!-- 2. Anomalies Alert Table -->
        <div class="dash-card">
          <h3 class="dash-card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#7c3aed" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
            จุดที่ตรวจพบความผิดปกติ / ต้องตรวจสอบหน้างาน (2 จุด)
          </h3>
          <div class="dash-table-wrap">
            <table class="dash-mini-table">
              <thead>
                <tr>
                  <th>แผง</th>
                  <th>ประเภท</th>
                  <th>ปริมาณการใช้</th>
                  <th>สาเหตุที่น่าสงสัย</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td><span class="table-stall-pill">A-04</span></td>
                  <td><span class="table-type-badge badge-water">มิเตอร์น้ำ</span></td>
                  <td><strong>85 ลบ.ม.</strong> <small style="color: #dc2626;">(+240%)</small></td>
                  <td><span class="table-status-pill pill-anomaly">สงสัยท่อแตก/รั่ว</span></td>
                </tr>
                <tr>
                  <td><span class="table-stall-pill">B-12</span></td>
                  <td><span class="table-type-badge badge-electric">มิเตอร์ไฟ</span></td>
                  <td><strong>0 kWh</strong> <small style="color: #64748b;">(คงที่)</small></td>
                  <td><span class="table-status-pill pill-anomaly">สงสัยมิเตอร์หยุดเดิน</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 3. Progress by Zone -->
        <div class="dash-card">
          <h3 class="dash-card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
            ความคืบหน้าการจดมิเตอร์แยกตามโซน
          </h3>
          <div class="dash-breakdown-list">
            <div class="dash-breakdown-row">
              <span>โซน A (แผงอาหารสด)</span>
              <strong>100% <small>(15/15 จุด - ครบแล้ว)</small></strong>
            </div>
            <div class="dash-progress-bar"><div class="dash-progress-fill" style="width: 100%; background: #059669;"></div></div>

            <div class="dash-breakdown-row">
              <span>โซน B (สินค้าแฟชั่น/ทั่วไป)</span>
              <strong>92.8% <small>(13/14 จุด - เหลือ 1 จุด)</small></strong>
            </div>
            <div class="dash-progress-bar"><div class="dash-progress-fill" style="width: 92.8%; background: #2563eb;"></div></div>

            <div class="dash-breakdown-row">
              <span>โซน C (อาหารปรุงสำเร็จ)</span>
              <strong>53.8% <small>(7/13 จุด - เหลือ 6 จุด)</small></strong>
            </div>
            <div class="dash-progress-bar"><div class="dash-progress-fill" style="width: 53.8%; background: #d97706;"></div></div>
          </div>
        </div>
      </template>
      </section>

      <section v-if="!isInspectDashboardMode" class="workspace-grid" :class="{ 'is-fullscreen-map-view': isMapOverviewMode }">
        <article class="panel gis-panel" :class="{ 'is-fullscreen-map-panel': isMapOverviewMode }">
          <!-- 1. Heading และ Toolbar ปกติสำหรับโหมดจัดการผัง GIS -->
          <div v-if="!isMapOverviewMode" class="panel-heading gis-panel-heading">
            <div class="gis-panel-header-left">
              <h2>{{ marketName || (isInspectMode ? 'ผังแผงเช่า' : 'ผังแผงค้า') }}</h2>
              <p v-if="isManageDataLoading" style="display: flex; align-items: center; gap: 6px; height: 18px;">
                <span class="kpi-skeleton-num" style="width: 80px; height: 14px;"></span>
              </p>
              <p v-else>ทั้งหมด {{ totalMarketStalls.toLocaleString('th-TH') }} แผง</p>
            </div>
            <div class="gis-heading-right">
              <div id="step1-map-toolbar" class="gis-map-toolbar"></div>
            </div>
          </div>

          <!-- ภาพรวมตลาด: KPI ลอยเหนือผัง ตามมุมมองที่เลือก (การชำระเงิน / การจดมิเตอร์ / ข้อมูลแผง) -->
          <div
            v-if="isMapOverviewMode && pageMode === 'payments'"
            class="market-map-kpi-strip"
            :class="{ 'has-side-panel': !!rightSidePanelWidth }"
            :style="rightSidePanelWidth ? { '--side-w': `${rightSidePanelWidth}px` } : undefined"
            aria-label="สรุปสถานะการชำระเงิน"
          >
            <button
              type="button"
              class="market-map-kpi tone-blue"
              :class="{ active: selectedPaymentFilter === 'due' }"
              @click="togglePaymentFilter('due')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                  <line x1="16" y1="2" x2="16" y2="6"></line>
                  <line x1="8" y1="2" x2="8" y2="6"></line>
                  <line x1="3" y1="10" x2="21" y2="10"></line>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>ต้องชำระ</small>
                <strong>{{ paymentSummary.dueCount.toLocaleString('th-TH') }} <em>รายการ</em></strong>
                <small class="kpi-sub">฿{{ paymentSummary.dueAmount.toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</small>
              </span>
            </button>

            <button
              type="button"
              class="market-map-kpi tone-green"
              :class="{ active: selectedPaymentFilter === 'paid' }"
              @click="togglePaymentFilter('paid')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                  <polyline points="22 4 12 14.01 9 11.01"></polyline>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>ชำระครบแล้ว</small>
                <strong>{{ paymentSummary.paidCount.toLocaleString('th-TH') }} <em>รายการ</em></strong>
                <small class="kpi-sub">฿{{ paymentSummary.paidAmount.toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</small>
              </span>
            </button>

            <button
              type="button"
              class="market-map-kpi tone-amber"
              :class="{ active: selectedPaymentFilter === 'pending' }"
              @click="togglePaymentFilter('pending')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>รอชำระ</small>
                <strong>{{ paymentSummary.pendingCount.toLocaleString('th-TH') }} <em>รายการ</em></strong>
                <small class="kpi-sub">฿{{ paymentSummary.pendingAmount.toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</small>
              </span>
            </button>

            <button
              type="button"
              class="market-map-kpi tone-red"
              :class="{ active: selectedPaymentFilter === 'overdue' }"
              @click="togglePaymentFilter('overdue')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="8" x2="12" y2="12"></line>
                  <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>ค้างชำระ</small>
                <strong>{{ paymentSummary.overdueCount.toLocaleString('th-TH') }} <em>รายการ</em></strong>
                <small class="kpi-sub">฿{{ paymentSummary.overdueAmount.toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</small>
              </span>
            </button>

            <button
              type="button"
              class="market-map-kpi tone-rose"
              :class="{ active: selectedPaymentFilter === 'total_overdue' }"
              @click="togglePaymentFilter('total_overdue')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="12" y1="1" x2="12" y2="23"></line>
                  <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>ยอดค้างรวม</small>
                <strong>{{ paymentSummary.overdueAmount.toLocaleString('th-TH') }} <em>บาท</em></strong>
                <small class="kpi-sub">จาก {{ paymentSummary.overdueCount }} รายการ</small>
              </span>
            </button>
          </div>

          <div
            v-else-if="isMapOverviewMode && pageMode === 'meters'"
            class="market-map-kpi-strip"
            :class="{ 'has-side-panel': !!rightSidePanelWidth }"
            :style="rightSidePanelWidth ? { '--side-w': `${rightSidePanelWidth}px` } : undefined"
            aria-label="สรุปสถานะการจดมิเตอร์"
          >
            <button
              type="button"
              class="market-map-kpi tone-blue"
              :class="{ active: selectedMeterFilter === 'total' }"
              @click="toggleMeterFilter('total')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m12 14 4-4"></path>
                  <path d="M3.34 19a10 10 0 1 1 17.32 0"></path>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>ต้องจด</small>
                <strong>{{ meterSummary.totalCount.toLocaleString('th-TH') }} <em>จุด</em></strong>
                <small class="kpi-sub">น้ำ {{ meterSummary.waterCount }} / ไฟ {{ meterSummary.electricCount }}</small>
              </span>
            </button>

            <button
              type="button"
              class="market-map-kpi tone-green"
              :class="{ active: selectedMeterFilter === 'recorded' }"
              @click="toggleMeterFilter('recorded')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                  <polyline points="22 4 12 14.01 9 11.01"></polyline>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>จดแล้ว</small>
                <strong>{{ meterSummary.recordedCount.toLocaleString('th-TH') }} <em>จุด</em></strong>
                <small class="kpi-sub">{{ meterSummary.completionRate }}% เสร็จสิ้น</small>
              </span>
            </button>

            <button
              type="button"
              class="market-map-kpi tone-amber"
              :class="{ active: selectedMeterFilter === 'pending' }"
              @click="toggleMeterFilter('pending')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>รอจด</small>
                <strong>{{ meterSummary.pendingCount.toLocaleString('th-TH') }} <em>จุด</em></strong>
                <small class="kpi-sub">ครบกำหนด {{ meterSummary.dueDateText }}</small>
              </span>
            </button>

            <button
              type="button"
              class="market-map-kpi tone-red"
              :class="{ active: selectedMeterFilter === 'overdue' }"
              @click="toggleMeterFilter('overdue')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path>
                  <line x1="12" y1="9" x2="12" y2="13"></line>
                  <line x1="12" y1="17" x2="12.01" y2="17"></line>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>เกินกำหนด</small>
                <strong>{{ meterSummary.overdueCount.toLocaleString('th-TH') }} <em>จุด</em></strong>
                <small class="kpi-sub">{{ meterSummary.overdueCount === 0 ? 'ไม่มีรายการตกค้าง' : `${meterSummary.overdueCount} จุด` }}</small>
              </span>
            </button>

            <button
              type="button"
              class="market-map-kpi tone-violet"
              :class="{ active: selectedMeterFilter === 'anomaly' }"
              @click="toggleMeterFilter('anomaly')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>ค่าผิดปกติ</small>
                <strong>{{ meterSummary.anomalyCount.toLocaleString('th-TH') }} <em>จุด</em></strong>
                <small class="kpi-sub">ใช้น้ำ/ไฟ สูงเกินเกณฑ์</small>
              </span>
            </button>
          </div>

          <div
            v-else-if="isMapOverviewMode"
            class="market-map-kpi-strip"
            :class="{ 'has-side-panel': !!rightSidePanelWidth }"
            :style="rightSidePanelWidth ? { '--side-w': `${rightSidePanelWidth}px` } : undefined"
            aria-label="สรุปสถานะแผงเช่า"
          >
            <button
              type="button"
              class="market-map-kpi is-total"
              :class="{ active: dashboardStatusFilter === 'all', 'is-loading': isManageDataLoading }"
              :disabled="isManageDataLoading"
              @click="setDashboardStatusFilter('all')"
            >
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="6" height="6"/><rect x="14" y="4" width="6" height="6"/><rect x="4" y="14" width="6" height="6"/><rect x="14" y="14" width="6" height="6"/></svg></span>
              <span class="market-map-kpi-copy">
                <small>ทั้งหมด</small>
                <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                  <span class="kpi-skeleton-num"></span>
                  <span class="kpi-skeleton-pill"></span>
                </div>
                <strong v-else>{{ currentDisplayedSummary.total.toLocaleString('th-TH') }} <em>แผง</em></strong>
              </span>
            </button>
            <button
              type="button"
              class="market-map-kpi is-available"
              :class="{ active: dashboardStatusFilter === 'available', 'is-loading': isManageDataLoading }"
              :disabled="isManageDataLoading"
              @click="setDashboardStatusFilter('available')"
            >
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="9"/><path d="m8 12 2.5 2.5L16 9"/></svg></span>
              <span class="market-map-kpi-copy">
                <small>ว่าง</small>
                <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                  <span class="kpi-skeleton-num"></span>
                  <span class="kpi-skeleton-pill"></span>
                </div>
                <strong v-else>{{ currentDisplayedSummary.available.toLocaleString('th-TH') }} <em>แผง</em></strong>
              </span>
            </button>
            <button
              type="button"
              class="market-map-kpi is-reserved"
              :class="{ active: dashboardStatusFilter === 'reserved', 'is-loading': isManageDataLoading }"
              :disabled="isManageDataLoading"
              @click="setDashboardStatusFilter('reserved')"
            >
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg></span>
              <span class="market-map-kpi-copy">
                <small>ติดจอง</small>
                <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                  <span class="kpi-skeleton-num"></span>
                  <span class="kpi-skeleton-pill"></span>
                </div>
                <strong v-else>{{ currentDisplayedSummary.reserved.toLocaleString('th-TH') }} <em>แผง</em></strong>
              </span>
            </button>
            <button
              type="button"
              class="market-map-kpi is-occupied"
              :class="{ active: dashboardStatusFilter === 'occupied', 'is-loading': isManageDataLoading }"
              :disabled="isManageDataLoading"
              @click="setDashboardStatusFilter('occupied')"
            >
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="8" r="4"/><path d="M5 21a7 7 0 0 1 14 0"/></svg></span>
              <span class="market-map-kpi-copy">
                <small>มีผู้เช่า</small>
                <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                  <span class="kpi-skeleton-num"></span>
                  <span class="kpi-skeleton-pill"></span>
                </div>
                <strong v-else>{{ currentDisplayedSummary.occupied.toLocaleString('th-TH') }} <em>แผง</em></strong>
              </span>
            </button>
            <button
              type="button"
              class="market-map-kpi is-legal-case"
              :class="{ active: dashboardStatusFilter === 'legal_case', 'is-loading': isManageDataLoading }"
              :disabled="isManageDataLoading"
              @click="setDashboardStatusFilter('legal_case')"
            >
              <span class="market-map-kpi-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m14 5 5 5"/><path d="m12 7 5 5"/><path d="M3 21h8"/><path d="m6 18 8-8"/><path d="m15 3 6 6"/>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>มีคดี</small>
                <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                  <span class="kpi-skeleton-num"></span>
                  <span class="kpi-skeleton-pill"></span>
                </div>
                <strong v-else>{{ currentDisplayedSummary.legalCases.toLocaleString('th-TH') }} <em>แผง</em></strong>
              </span>
            </button>
            <button
              type="button"
              class="market-map-kpi is-action-required"
              :class="{ active: dashboardStatusFilter === 'action_required', 'is-loading': isManageDataLoading }"
              :disabled="isManageDataLoading"
              @click="setDashboardStatusFilter('action_required')"
            >
              <span class="market-map-kpi-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
                  <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
                  <path d="m9 14 2 2 4-4"/>
                </svg>
              </span>
              <span class="market-map-kpi-copy">
                <small>ต้องดำเนินการ</small>
                <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                  <span class="kpi-skeleton-num"></span>
                  <span class="kpi-skeleton-pill"></span>
                </div>
                <strong v-else>{{ currentDisplayedSummary.actionRequired.toLocaleString('th-TH') }} <em>แผง</em></strong>
              </span>
            </button>
          </div>

          <MarketProjectsOverviewMap
            ref="mapRef"
            :marketId="marketId"
            :managedProjects="projects"
            :managedFacilities="facilities"
            :selectedProjectId="selectedProject?.id"
            :selectedStallCodes="activeHighlightedStallCodes"
            :stallStatusOverrides="stallStatusOverrides"
            :stallStatusFilter="dashboardStatusFilter === 'legal_case' ? 'all' : dashboardStatusFilter"
            :actionRequiredStallCodes="actionRequiredStallCodes"
            :operationalFilterMeta="stallOperationalFilterMeta"
            :activeFilterCodes="mapActiveFilterCodes"
            :isDashboardOpen="dashboardPanelOpen"
            :dashboardTitle="dashboardViewTitle"
            :dashboardSubtitle="dashboardViewSubtitle"
            @update:isDashboardOpen="dashboardPanelOpen = $event"
            @toggleDashboard="dashboardPanelOpen = !dashboardPanelOpen"
            @update:activePanel="activeMapSidePanel = $event"
            :onlyAvailable="false"
            :categoryFilter="null"
            :demandViewMode="demandViewMode"
            :demandLevelFilter="demandLevelFilter"
            :stallDemandMap="stallDemandMap"
            @update:demandLevelFilter="demandLevelFilter = $event"
            mapStyle="standard"
            managementMode
            :inspectOnly="isInspectMode"
            :isOverviewPage="isMapOverviewMode"
            :filterContext="pageMode === 'costs' ? 'costs' : 'data'"
            :filterBadge="costFilterBadge"
            :hideRail="costDrawerOpen"
            :hideProjectPanel="isMapOverviewMode && showInspectStallOverlay"
            :expandProjectPanel="isMapOverviewMode && isProjectPanelMorphing"
            :returnProjectPanel="isMapOverviewMode && isProjectPanelReturning"
            @clearStallSelection="handleClearStallSelection"
            @reset-cost-filters="resetCostFilters"
            @selectProject="selectProject"
            @requestExitProject="handleRequestExitProject"
            @exitProject="clearSelectedProject"
            @selectStall="selectStall"
            @updateProjectPolygon="updateProjectPolygon"
            @updateStallPolygon="updateStallPolygon"
            @requestAddProject="openZoneEditor"
            @requestAddStallLayout="openDrawnStall"
            @requestAddFacility="openFacilityEditor"
            @requestDeleteFacility="deleteFacility"
            @requestDeleteStall="deleteStall"
            @requestDeleteProject="deleteProject"
          >
            <template #dashboard>
              <div class="ov-dash-container">
                <!-- MODE 1: แดชบอร์ดข้อมูลแผงและสัญญาเช่า (stalls) -->
                <div v-if="pageMode === 'stalls'" class="ov-dash-mode-view">
                  <!-- การ์ดสถานะแผงสำหรับโหมด Responsive (ย้ายมาจากแถบลอยบนแผนที่) -->
                  <section class="ov-dash-section-card ov-dash-responsive-kpis" aria-labelledby="responsive-stall-status-heading">
                    <div class="ov-dash-section-head">
                      <div>
                        <h4 id="responsive-stall-status-heading" class="ov-dash-section-title">สถานะแผงค้าในผังตลาด</h4>
                        <p class="ov-dash-section-note">กดเลือกสถานะเพื่อกรองแผงบนแผนที่</p>
                      </div>
                    </div>
                    <div class="ov-dash-status-pill-grid">
                      <!-- 1. ทั้งหมด -->
                      <button
                        type="button"
                        class="ov-dash-status-pill is-total"
                        :class="{ active: dashboardStatusFilter === 'all', 'is-loading': isManageDataLoading }"
                        :disabled="isManageDataLoading"
                        @click="setDashboardStatusFilter('all')"
                      >
                        <span class="ov-dash-status-pill-icon" aria-hidden="true">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect x="4" y="4" width="6" height="6"/><rect x="14" y="4" width="6" height="6"/>
                            <rect x="4" y="14" width="6" height="6"/><rect x="14" y="14" width="6" height="6"/>
                          </svg>
                        </span>
                        <div class="ov-dash-status-pill-info">
                          <span class="ov-dash-status-pill-label">ทั้งหมด</span>
                          <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                            <span class="kpi-skeleton-num"></span>
                          </div>
                          <strong v-else class="ov-dash-status-pill-val">
                            {{ currentDisplayedSummary.total.toLocaleString('th-TH') }} <small>แผง</small>
                          </strong>
                        </div>
                      </button>

                      <!-- 2. ว่าง -->
                      <button
                        type="button"
                        class="ov-dash-status-pill is-available"
                        :class="{ active: dashboardStatusFilter === 'available', 'is-loading': isManageDataLoading }"
                        :disabled="isManageDataLoading"
                        @click="setDashboardStatusFilter('available')"
                      >
                        <span class="ov-dash-status-pill-icon" aria-hidden="true">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="9"/><path d="m8 12 2.5 2.5L16 9"/>
                          </svg>
                        </span>
                        <div class="ov-dash-status-pill-info">
                          <span class="ov-dash-status-pill-label">ว่าง</span>
                          <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                            <span class="kpi-skeleton-num"></span>
                          </div>
                          <strong v-else class="ov-dash-status-pill-val">
                            {{ currentDisplayedSummary.available.toLocaleString('th-TH') }} <small>แผง</small>
                          </strong>
                        </div>
                      </button>

                      <!-- 3. ติดจอง -->
                      <button
                        type="button"
                        class="ov-dash-status-pill is-reserved"
                        :class="{ active: dashboardStatusFilter === 'reserved', 'is-loading': isManageDataLoading }"
                        :disabled="isManageDataLoading"
                        @click="setDashboardStatusFilter('reserved')"
                      >
                        <span class="ov-dash-status-pill-icon" aria-hidden="true">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/>
                          </svg>
                        </span>
                        <div class="ov-dash-status-pill-info">
                          <span class="ov-dash-status-pill-label">ติดจอง</span>
                          <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                            <span class="kpi-skeleton-num"></span>
                          </div>
                          <strong v-else class="ov-dash-status-pill-val">
                            {{ currentDisplayedSummary.reserved.toLocaleString('th-TH') }} <small>แผง</small>
                          </strong>
                        </div>
                      </button>

                      <!-- 4. มีผู้เช่า -->
                      <button
                        type="button"
                        class="ov-dash-status-pill is-occupied"
                        :class="{ active: dashboardStatusFilter === 'occupied', 'is-loading': isManageDataLoading }"
                        :disabled="isManageDataLoading"
                        @click="setDashboardStatusFilter('occupied')"
                      >
                        <span class="ov-dash-status-pill-icon" aria-hidden="true">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="8" r="4"/><path d="M5 21a7 7 0 0 1 14 0"/>
                          </svg>
                        </span>
                        <div class="ov-dash-status-pill-info">
                          <span class="ov-dash-status-pill-label">มีผู้เช่า</span>
                          <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                            <span class="kpi-skeleton-num"></span>
                          </div>
                          <strong v-else class="ov-dash-status-pill-val">
                            {{ currentDisplayedSummary.occupied.toLocaleString('th-TH') }} <small>แผง</small>
                          </strong>
                        </div>
                      </button>

                      <!-- 5. มีคดี -->
                      <button
                        type="button"
                        class="ov-dash-status-pill is-legal-case"
                        :class="{ active: dashboardStatusFilter === 'legal_case', 'is-loading': isManageDataLoading }"
                        :disabled="isManageDataLoading"
                        @click="setDashboardStatusFilter('legal_case')"
                      >
                        <span class="ov-dash-status-pill-icon" aria-hidden="true">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="m14 5 5 5"/><path d="m12 7 5 5"/><path d="M3 21h8"/><path d="m6 18 8-8"/><path d="m15 3 6 6"/>
                          </svg>
                        </span>
                        <div class="ov-dash-status-pill-info">
                          <span class="ov-dash-status-pill-label">มีคดี</span>
                          <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                            <span class="kpi-skeleton-num"></span>
                          </div>
                          <strong v-else class="ov-dash-status-pill-val">
                            {{ currentDisplayedSummary.legalCases.toLocaleString('th-TH') }} <small>แผง</small>
                          </strong>
                        </div>
                      </button>

                      <!-- 6. ต้องดำเนินการ -->
                      <button
                        type="button"
                        class="ov-dash-status-pill is-action-required"
                        :class="{ active: dashboardStatusFilter === 'action_required', 'is-loading': isManageDataLoading }"
                        :disabled="isManageDataLoading"
                        @click="setDashboardStatusFilter('action_required')"
                      >
                        <span class="ov-dash-status-pill-icon" aria-hidden="true">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
                            <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
                            <path d="m9 14 2 2 4-4"/>
                          </svg>
                        </span>
                        <div class="ov-dash-status-pill-info">
                          <span class="ov-dash-status-pill-label">ต้องดำเนินการ</span>
                          <div v-if="isManageDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                            <span class="kpi-skeleton-num"></span>
                          </div>
                          <strong v-else class="ov-dash-status-pill-val">
                            {{ currentDisplayedSummary.actionRequired.toLocaleString('th-TH') }} <small>แผง</small>
                          </strong>
                        </div>
                      </button>
                    </div>
                  </section>

                  <section class="ov-dash-section-card ov-project-status-card" aria-labelledby="project-status-heading">
                    <div class="ov-dash-section-head">
                      <div>
                        <h4 id="project-status-heading" class="ov-dash-section-title">สัดส่วนสถานะแผงรายโครงการ</h4>
                        <p class="ov-dash-section-note">เลือกโครงการเพื่อดูสัดส่วนสถานะแผง</p>
                      </div>
                    </div>
                    <AppSelect
                      :model-value="selectedDashboardProjectId"
                      :options="dashboardProjectOptions"
                      placeholder="เลือกโครงการ"
                      @update:model-value="handleDashboardProjectChange"
                    />
                    <div v-if="selectedDashboardProjectStatus" class="ov-project-donut-focus">
                      <div class="ov-project-donut-wrap">
                        <div
                          class="ov-project-donut-chart"
                          :style="selectedDashboardProjectDonutStyle"
                          role="img"
                          :aria-label="`${selectedDashboardProjectStatus.name}: ว่าง ${selectedDashboardProjectStatus.available}, มีผู้เช่า ${selectedDashboardProjectStatus.occupied}, ติดจอง ${selectedDashboardProjectStatus.reserved}`"
                        >
                          <span>
                            <strong>{{ selectedDashboardProjectStatus.total.toLocaleString('th-TH') }}</strong>
                            <small>แผงทั้งหมด</small>
                          </span>
                        </div>
                        <strong class="ov-project-donut-name">{{ selectedDashboardProjectStatus.name }}</strong>
                      </div>
                      <div class="ov-project-donut-stats">
                        <div class="is-available">
                          <span><i></i>ว่าง</span>
                          <strong>{{ selectedDashboardProjectStatus.available.toLocaleString('th-TH') }}</strong>
                          <small>{{ selectedDashboardProjectStatus.availablePct.toFixed(1) }}%</small>
                        </div>
                        <div class="is-occupied">
                          <span><i></i>มีผู้เช่า</span>
                          <strong>{{ selectedDashboardProjectStatus.occupied.toLocaleString('th-TH') }}</strong>
                          <small>{{ selectedDashboardProjectStatus.occupiedPct.toFixed(1) }}%</small>
                        </div>
                        <div class="is-reserved">
                          <span><i></i>ติดจอง</span>
                          <strong>{{ selectedDashboardProjectStatus.reserved.toLocaleString('th-TH') }}</strong>
                          <small>{{ selectedDashboardProjectStatus.reservedPct.toFixed(1) }}%</small>
                        </div>
                      </div>
                    </div>
                    <div v-else class="ov-project-donut-empty">
                      ยังไม่มีข้อมูลโครงการ
                    </div>
                  </section>

                  <section class="ov-dash-section-card" aria-labelledby="all-project-donuts-heading">
                    <div class="ov-dash-section-head">
                      <div>
                        <h4 id="all-project-donuts-heading" class="ov-dash-section-title">ภาพรวมทุกโครงการ</h4>
                        <p class="ov-dash-section-note">เปรียบเทียบสัดส่วนสถานะแผง</p>
                      </div>
                    </div>
                    <div class="ov-project-mini-donut-grid">
                      <article
                        v-for="project in stallProjectStatusBreakdown"
                        :key="`donut-${project.id}`"
                        class="ov-project-mini-donut-card"
                        :class="{ active: selectedDashboardProjectId === project.id }"
                        tabindex="0"
                        :aria-label="`${project.name}: ว่าง ${project.available}, มีผู้เช่า ${project.occupied}, ติดจอง ${project.reserved}`"
                        @click="selectDashboardProject(project.id)"
                        @keydown.enter.prevent="selectDashboardProject(project.id)"
                        @keydown.space.prevent="selectDashboardProject(project.id)"
                      >
                        <span
                          v-if="project.badgeText"
                          class="ov-project-mini-badge"
                          :class="`is-${project.badgeTone}`"
                        >{{ project.badgeText }}</span>
                        <div class="ov-project-mini-donut" :style="projectDonutStyle(project)">
                          <span>
                            <strong>{{ project.total.toLocaleString('th-TH') }}</strong>
                            <small>แผง</small>
                          </span>
                        </div>
                        <strong class="ov-project-mini-name">{{ project.name }}</strong>
                        <div class="ov-project-mini-tooltip" role="tooltip">
                          <span class="is-available">ว่าง <strong>{{ project.available }}</strong><small>{{ project.availablePct.toFixed(1) }}%</small></span>
                          <span class="is-occupied">มีผู้เช่า <strong>{{ project.occupied }}</strong><small>{{ project.occupiedPct.toFixed(1) }}%</small></span>
                          <span class="is-reserved">ติดจอง <strong>{{ project.reserved }}</strong><small>{{ project.reservedPct.toFixed(1) }}%</small></span>
                        </div>
                      </article>
                    </div>
                  </section>

                  <section class="ov-dash-section-card" aria-labelledby="contract-health-heading">
                    <div class="ov-dash-section-head">
                      <h4 id="contract-health-heading" class="ov-dash-section-title">สถานะสัญญาเช่า</h4>
                    </div>
                    <div class="ov-contract-health-layout">
                      <div
                        class="ov-contract-donut"
                        :style="{ background: `conic-gradient(#16a34a 0 ${(contractHealthSummary.normalCount / Math.max(1, contractHealthSummary.activeCount)) * 100}%, #f59e0b 0 ${((contractHealthSummary.normalCount + contractHealthSummary.expiringCount) / Math.max(1, contractHealthSummary.activeCount)) * 100}%, #ef4444 0 100%)` }"
                      >
                        <span><strong>{{ contractHealthSummary.activeCount }}</strong><small>สัญญา</small></span>
                      </div>
                      <div class="ov-contract-health-legend">
                        <span><i class="is-active"></i>อยู่ในสัญญา <strong>{{ contractHealthSummary.normalCount }}</strong></span>
                        <span><i class="is-expiring"></i>ใกล้หมดอายุ <strong>{{ contractHealthSummary.expiringCount }}</strong></span>
                        <span><i class="is-expired"></i>หมดอายุ <strong>{{ contractHealthSummary.expiredCount }}</strong></span>
                      </div>
                    </div>
                    <div class="ov-contract-renewal-row">
                      <span>ต่อสัญญาแล้ว <strong>{{ stallContractRenewalSummary.submitted }}</strong></span>
                      <span>รอตอบรับ <strong>{{ stallContractRenewalSummary.pending }}</strong></span>
                      <span>ไม่ต่อ <strong>{{ stallContractRenewalSummary.declined }}</strong></span>
                    </div>
                  </section>

                  <section class="ov-dash-section-card" aria-labelledby="return-status-heading">
                    <div class="ov-dash-section-head">
                      <h4 id="return-status-heading" class="ov-dash-section-title">การส่งคืนแผงค้า</h4>
                    </div>
                    <div class="ov-status-summary-grid is-five">
                      <div><span>แจ้งไม่ต่อ</span><strong>{{ stallReturnStatusSummary.notified }}</strong></div>
                      <div><span>รอส่งคืน</span><strong>{{ stallReturnStatusSummary.waiting }}</strong></div>
                      <div><span>รอตรวจรับ</span><strong>{{ stallReturnStatusSummary.inspection }}</strong></div>
                      <div class="is-success"><span>รับคืนแล้ว</span><strong>{{ stallReturnStatusSummary.returned }}</strong></div>
                      <div class="is-danger"><span>เกินกำหนด</span><strong>{{ stallReturnStatusSummary.overdue }}</strong></div>
                    </div>
                  </section>

                  <section class="ov-dash-section-card" aria-labelledby="legal-status-heading">
                    <div class="ov-dash-section-head">
                      <h4 id="legal-status-heading" class="ov-dash-section-title">การดำเนินคดี</h4>
                    </div>
                    <div class="ov-status-summary-grid">
                      <div><span>รอฝ่ายกฎหมาย</span><strong>{{ stallLegalStatusSummary.awaiting }}</strong></div>
                      <div class="is-pending"><span>กำลังดำเนินคดี</span><strong>{{ stallLegalStatusSummary.proceeding }}</strong></div>
                      <div class="is-danger"><span>บังคับคดี</span><strong>{{ stallLegalStatusSummary.enforcement }}</strong></div>
                      <div class="is-success"><span>สิ้นสุดคดี</span><strong>{{ stallLegalStatusSummary.closed }}</strong></div>
                    </div>
                  </section>

                  <section class="ov-dash-section-card" aria-labelledby="maintenance-status-heading">
                    <div class="ov-dash-section-head">
                      <h4 id="maintenance-status-heading" class="ov-dash-section-title">งานซ่อมบำรุง</h4>
                    </div>
                    <div class="ov-maintenance-stacked-bar" aria-hidden="true">
                      <span class="is-pending" :style="{ flex: stallMaintenanceDashboardSummary.pending }"></span>
                      <span class="is-repairing" :style="{ flex: stallMaintenanceDashboardSummary.repairing }"></span>
                      <span class="is-completed" :style="{ flex: stallMaintenanceDashboardSummary.completed }"></span>
                      <span class="is-overdue" :style="{ flex: stallMaintenanceDashboardSummary.overdue }"></span>
                    </div>
                    <div class="ov-maintenance-legend">
                      <span><i class="is-pending"></i>รอดำเนินการ <strong>{{ stallMaintenanceDashboardSummary.pending }}</strong></span>
                      <span><i class="is-repairing"></i>กำลังซ่อม <strong>{{ stallMaintenanceDashboardSummary.repairing }}</strong></span>
                      <span><i class="is-completed"></i>ซ่อมเสร็จ <strong>{{ stallMaintenanceDashboardSummary.completed }}</strong></span>
                      <span><i class="is-overdue"></i>เกินกำหนด <strong>{{ stallMaintenanceDashboardSummary.overdue }}</strong></span>
                    </div>
                  </section>

                </div>

                <!-- MODE 2: แดชบอร์ดการเงินและการจัดเก็บ (payments) -->
                <div v-else-if="pageMode === 'payments'" class="ov-dash-mode-view">
                  <section class="ov-dash-summary-card" aria-labelledby="payment-summary-heading">
                    <div class="ov-stall-section-head">
                      <div>
                        <h4 id="payment-summary-heading">ภาพรวมการชำระเงิน</h4>
                        <p class="ov-stall-section-note">รอบเรียกเก็บปัจจุบัน</p>
                      </div>
                    </div>
                    <div class="ov-dash-summary-main">
                      <span>ยอดเรียกเก็บ</span>
                      <strong>฿284,500 <small>42 บิล</small></strong>
                    </div>
                    <div class="ov-dash-summary-grid">
                      <div class="is-success">
                        <span>จัดเก็บแล้ว</span>
                        <strong>฿192,000</strong>
                      </div>
                      <div class="is-pending">
                        <span>รอชำระ</span>
                        <strong>฿58,500</strong>
                      </div>
                      <div class="is-danger">
                        <span>ค้างชำระ</span>
                        <strong>฿34,000</strong>
                      </div>
                    </div>
                  </section>

                  <div class="ov-dash-section-card">
                    <div class="ov-dash-section-head">
                      <h4 class="ov-dash-section-title">รายได้และโครงสร้างรายรับ</h4>
                    </div>
                    <div class="ov-dash-income-total">
                      <span>รายได้ค่าเช่าต่อเดือน</span>
                      <strong>฿{{ contractHealthSummary.totalMonthlyRent.toLocaleString('th-TH') }}</strong>
                    </div>
                    <div class="ov-dash-category-list">
                      <div class="ov-dash-category-row">
                        <div class="ov-dash-cat-left">
                          <span class="ov-dash-cat-name">ค่าเช่าพื้นที่แผง</span>
                          <span class="ov-dash-cat-pct">฿211,000 (74.2%)</span>
                        </div>
                        <div class="ov-dash-cat-bar"><span style="width: 74.2%"></span></div>
                      </div>
                      <div class="ov-dash-category-row">
                        <div class="ov-dash-cat-left">
                          <span class="ov-dash-cat-name">ค่าไฟฟ้า</span>
                          <span class="ov-dash-cat-pct">฿46,500 (16.3%)</span>
                        </div>
                        <div class="ov-dash-cat-bar"><span style="width: 16.3%"></span></div>
                      </div>
                      <div class="ov-dash-category-row">
                        <div class="ov-dash-cat-left">
                          <span class="ov-dash-cat-name">ค่าน้ำประปา</span>
                          <span class="ov-dash-cat-pct">฿19,000 (6.7%)</span>
                        </div>
                        <div class="ov-dash-cat-bar"><span style="width: 6.7%"></span></div>
                      </div>
                      <div class="ov-dash-category-row">
                        <div class="ov-dash-cat-left">
                          <span class="ov-dash-cat-name">ค่าส่วนกลาง</span>
                          <span class="ov-dash-cat-pct">฿8,000 (2.8%)</span>
                        </div>
                        <div class="ov-dash-cat-bar"><span style="width: 2.8%"></span></div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- MODE 3: แดชบอร์ดการจดมิเตอร์น้ำและไฟฟ้า (meters) -->
                <div v-else-if="pageMode === 'meters'" class="ov-dash-mode-view">
                  <section class="ov-dash-summary-card" aria-labelledby="meter-summary-heading">
                    <div class="ov-stall-section-head">
                      <div>
                        <h4 id="meter-summary-heading">ภาพรวมการจดมิเตอร์</h4>
                        <p class="ov-stall-section-note">รอบจดมิเตอร์ปัจจุบัน</p>
                      </div>
                    </div>
                    <div class="ov-dash-meter-progress">
                      <div class="ov-dash-meter-progress-head">
                        <span>จดแล้ว</span>
                        <strong>38 จาก 42 <small>มิเตอร์</small></strong>
                      </div>
                      <div class="ov-dash-health-track">
                        <span class="track-good" style="width: 90.5%"></span>
                        <span class="track-warn" style="width: 9.5%"></span>
                      </div>
                    </div>
                    <div class="ov-dash-summary-grid">
                      <div class="is-pending">
                        <span>รอดำเนินการ</span>
                        <strong>4 <small>มิเตอร์</small></strong>
                      </div>
                      <div>
                        <span>น้ำประปา</span>
                        <strong>1,240 <small>ลบ.ม.</small></strong>
                      </div>
                      <div>
                        <span>ไฟฟ้า</span>
                        <strong>9,850 <small>kWh</small></strong>
                      </div>
                    </div>
                  </section>

                  <div class="ov-dash-section-card">
                    <div class="ov-dash-section-head">
                      <div class="ov-dash-section-title-row">
                        <h4 class="ov-dash-section-title">มิเตอร์ใช้ผิดปกติ</h4>
                        <span class="ov-dash-count-badge badge-red">2</span>
                      </div>
                    </div>
                    <div class="ov-dash-watchlist">
                      <div class="ov-dash-watchlist-row urgency-urgent" @click="focusExpiringContractStall('18-042')" role="button" tabindex="0">
                        <div class="ov-dash-w-left">
                          <span class="ov-dash-code-tag">18-042</span>
                          <span class="ov-dash-w-tenant">มิเตอร์น้ำ #W-18042</span>
                        </div>
                        <div class="ov-dash-w-right">
                          <span class="ov-dash-w-pill pill-urgent">+320% (85 ลบ.ม.)</span>
                          <span class="ov-dash-w-icon" aria-hidden="true">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                          </span>
                        </div>
                      </div>
                      <div class="ov-dash-watchlist-row urgency-warning" @click="focusExpiringContractStall('12-108')" role="button" tabindex="0">
                        <div class="ov-dash-w-left">
                          <span class="ov-dash-code-tag">12-108</span>
                          <span class="ov-dash-w-tenant">มิเตอร์ไฟ #E-12108</span>
                        </div>
                        <div class="ov-dash-w-right">
                          <span class="ov-dash-w-pill pill-warning">+210% (1,420 kWh)</span>
                          <span class="ov-dash-w-icon" aria-hidden="true">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </template>
            <template #filter-costs>
              <div class="cf-wrap">
                <p class="cf-hint">{{ costTool === 'bill' ? 'กรองบิลแผงเช่าชั่วคราว' : 'กรองบิลค่าน้ำ ค่าไฟ' }} (เลือกปุ่มมุมซ้ายล่างเพื่อสลับชนิดบิล)</p>
                <div v-if="costTool !== 'bill'" class="cf-group">
                  <span class="cf-label">รอบบิล</span>
                  <AppSelect v-model="costMonth" :options="costMonthOptions" />
                </div>
                <div class="cf-group">
                  <span class="cf-label">สถานะบิล</span>
                  <div class="cf-chips">
                    <button v-for="c in costCards" :key="c.filter" type="button" class="cf-chip" :class="{ active: costFilter === c.filter }" :aria-pressed="costFilter === c.filter" @click="costFilter = c.filter">
                      {{ c.label }} <span>{{ c.count.toLocaleString('th-TH') }}</span>
                    </button>
                  </div>
                </div>
                <div class="cf-group">
                  <span class="cf-label">ช่วงยอดเงินต่อบิล</span>
                  <div class="cf-chips">
                    <button v-for="o in costAmountOptions" :key="o.value" type="button" class="cf-chip" :class="{ active: costAmount === o.value }" :aria-pressed="costAmount === o.value" @click="costAmount = o.value">{{ o.label }}</button>
                  </div>
                </div>
                <p class="cf-result">แสดง <strong>{{ costFilterCodes.length ? costFilterCodes.length.toLocaleString('th-TH') : 'ทั้งหมด' }}</strong>{{ costFilterCodes.length ? ' แผงที่ไฮไลต์บนแผนที่' : '' }}</p>
              </div>
            </template>
            <template #overlay>
              <div
                v-if="showInspectStallOverlay && stallEdit.ident"
                id="map-selected-stall-drawer"
                class="inspect-selected-stall-drawer"
                :class="{
                  'is-management-mode': !isInspectMode,
                  'is-edit-mode': stallDrawerMode === 'edit',
                  'is-full-detail-mode': stallDrawerMode === 'full-detail',
                  'is-drawer-expanded': stallDrawerExpanded,
                  'is-drawer-closing': isStallDrawerClosing,
                  'is-project-morph-swap': isProjectMorphSwap,
                  'is-left-dock': isMapOverviewMode
                }"
              >
                  <!-- 1. Header -->
                  <div class="issd-header">
                    <div v-if="isMapOverviewMode" class="issd-return-row">
                      <button
                        type="button"
                        class="issd-return-btn"
                        :aria-label="stallDrawerMode === 'edit' ? 'รายละเอียดแผง' : (selectedProject?.name ? `กลับไป ${selectedProject.name}` : 'กลับ')"
                        @click="stallDrawerMode === 'edit' && !isCreatingNewStall ? cancelStallEdit() : closeStall()"
                      >
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="15 18 9 12 15 6"></polyline></svg>
                        <span>{{ stallDrawerMode === 'edit' ? 'รายละเอียดแผง' : (selectedProject?.name || 'กลับ') }}</span>
                      </button>
                    </div>
                    <div class="issd-header-top">
                      <div class="issd-title-group">
                        <button
                          v-if="!isMapOverviewMode && (stallDrawerMode === 'edit' && !isCreatingNewStall)"
                          type="button"
                          class="issd-back-btn"
                          title="กลับไปดูรายละเอียด"
                          aria-label="ย้อนกลับ"
                          @click="cancelStallEdit()"
                        >
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="15 18 9 12 15 6"></polyline>
                          </svg>
                        </button>
                        <strong>{{ stallDrawerMode === 'edit' ? (isCreatingNewStall ? 'เพิ่มแผงเช่าใหม่' : `แก้ไขแผง ${stallEdit.ident}`) : `แผง ${stallEdit.ident}` }}</strong>
                        <span
                          v-if="stallDrawerMode !== 'edit'"
                          class="issd-status-pill"
                          :class="`is-${stallEdit.stall_status || 'available'}`"
                        >
                          {{ stallEdit.stall_status === 'reserved' ? 'ติดจอง' : (stallEdit.stall_status === 'occupied' ? 'มีผู้เช่า' : 'ว่าง') }}
                        </span>
                      </div>
                      <button
                        v-if="!isMapOverviewMode"
                        type="button"
                        class="issd-close-btn"
                        title="ปิดหน้าต่าง"
                        aria-label="ปิด"
                        @click="closeStall"
                      >
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                      </button>
                    </div>
                    <div class="issd-project-sub">
                      <template v-if="stallDrawerMode === 'edit'">
                        {{ isCreatingNewStall ? 'ตั้งค่ารหัส ขนาด และอัตราค่าเช่าสำหรับแผงใหม่' : `รหัสแผง ${stallEdit.ident} · แก้ไขข้อมูลและอัตราค่าเช่า` }}
                      </template>
                      <template v-else>
                        {{ cleanProjectName(selectedProject?.name || stallEdit.project_name || 'โซนแผงค้า') }}<template v-if="(selectedProject?.floors?.length || 0) > 1"> · {{ formatFloorDisplay(stallEdit.floor) }}</template>
                      </template>
                    </div>
                  </div>

                  <!-- 2. Body -->
                  <!-- 2A. Mini Summary Mode -->
                  <!-- 2A-1. โหมดจัดการ GIS: Drawer แบบเดิม -->
                  <div v-if="stallDrawerMode === 'detail' && !isInspectMode" class="issd-body">
                    <!-- Stall Photo Banner -->
                    <div
                      class="issd-photo-banner"
                      title="คลิกเพื่อดูภาพแผงเช่าทั้งหมด"
                      @click="openInspectPhotoPreview(0)"
                    >
                      <img :src="inspectStallPhotos[0]" :alt="`ภาพแผง ${stallEdit.ident}`" />
                      <div class="issd-photo-count-badge">
                        <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                        <span>{{ inspectStallPhotos.length }} รูป</span>
                      </div>
                    </div>

                    <!-- Stall Specifications & Rates -->
                    <dl class="issd-detail-list">
                      <div class="issd-detail-row">
                        <dt>หมวดสินค้า</dt>
                        <dd>{{ stallEdit.category || selectedProject?.category || 'ทั่วไป' }}</dd>
                      </div>
                      <div class="issd-detail-row">
                        <dt>ขนาดพื้นที่</dt>
                        <dd>{{ stallEdit.width }} × {{ stallEdit.depth }} ม. ({{ inspectStallSizeSqm }} ตร.ม.)</dd>
                      </div>
                      <div class="issd-detail-row">
                        <dt>ค่าเช่ารายเดือน</dt>
                        <dd class="is-rent">{{ Number(stallEdit.monthly_rent || 1800).toLocaleString('th-TH') }} บาท</dd>
                      </div>
                      <div class="issd-detail-row">
                        <dt>ค่าเช่ารายวัน</dt>
                        <dd>{{ Number(stallEdit.daily_rent || 150).toLocaleString('th-TH') }} บาท</dd>
                      </div>
                      <div class="issd-detail-row">
                        <dt>ค่าธรรมเนียมแรกเข้า</dt>
                        <dd class="is-bold">{{ Number(inspectStallEntryFee).toLocaleString('th-TH') }} บาท</dd>
                      </div>
                    </dl>
                  </div>

                  <!-- 2A-2. โหมดดูข้อมูล (inspect): Drawer ใหม่ -->
                  <div v-else-if="stallDrawerMode === 'detail'" class="issd-body">
                    <!-- Stall Photo Banner -->
                    <div
                      class="issd-photo-banner"
                      title="คลิกเพื่อดูภาพแผงเช่าทั้งหมด"
                      @click="openInspectPhotoPreview(0)"
                    >
                      <img :src="inspectStallPhotos[0]" :alt="`ภาพแผง ${stallEdit.ident}`" />
                      <div class="issd-photo-count-badge">
                        <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                        <span>{{ inspectStallPhotos.length }} รูป</span>
                      </div>
                    </div>

                    <!-- 1. Urgent Action Alert (Expandable in place without navigating away) -->
                    <div v-if="stallActionRichItems.length > 0" class="issd-mini-alert-card">
                      <button
                        type="button"
                        class="issd-mini-alert-banner"
                        :aria-expanded="isMiniActionAlertExpanded"
                        @click="isMiniActionAlertExpanded = !isMiniActionAlertExpanded"
                      >
                        <div class="issd-mini-alert-head">
                          <span class="issd-aac-icon-wrap">
                            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                              <rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect>
                              <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path>
                              <path d="m9 14 2 2 4-4"></path>
                            </svg>
                          </span>
                          <strong>สิ่งที่ต้องดำเนินการ</strong>
                          <span class="issd-aac-count-badge">{{ stallActionRichItems.length }} รายการ</span>
                        </div>
                        <span class="issd-aac-toggle-btn" aria-hidden="true">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" :class="{ 'is-flipped': isMiniActionAlertExpanded }">
                            <polyline points="6 9 12 15 18 9"></polyline>
                          </svg>
                        </span>
                      </button>

                      <!-- Expanded Action Item List right here inside mini mode -->
                      <transition name="fade">
                        <div v-show="isMiniActionAlertExpanded" class="issd-mini-alert-list">
                          <div
                            v-for="act in stallActionRichItems"
                            :key="act.key"
                            class="issd-mini-alert-row"
                            :class="`is-${act.severity}`"
                            role="button"
                            tabindex="0"
                            @click="openStallFullDetail(act.targetTab)"
                            @keydown.enter="openStallFullDetail(act.targetTab)"
                          >
                            <span class="issd-aac-bar" :class="`is-${act.severity}`" aria-hidden="true"></span>
                            <div class="issd-aac-row-text">
                              <strong>{{ act.category }}<b v-if="act.count" class="issd-aac-count" :class="`is-${act.severity}`">{{ act.count }}</b></strong>
                              <small>{{ act.title }}</small>
                            </div>
                            <svg class="issd-mini-alert-chevron" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 18 15 12 9 6"></polyline></svg>
                          </div>
                        </div>
                      </transition>
                    </div>

                    <!-- 2. Primary Occupant & Contract / Booking Section (แยกส่วนผู้เช่า/ผู้จอง กับ สัญญา/การจอง ชัดเจน) -->
                    <!-- 2.1 กรณีแผงว่าง (available) -->
                    <div
                      v-if="stallEdit.stall_status === 'available'"
                      class="issd-card-section issd-occupant-card is-available"
                    >
                      <div class="issd-available-line">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
                        <strong>พร้อมเปิดรับจอง</strong>
                      </div>
                    </div>

                    <!-- 2.2 กรณีติดจอง (reserved หรือมีรายการจอง/คำร้อง) -->
                    <div
                      v-else-if="stallEdit.stall_status === 'reserved' || (stallEdit.stall_status !== 'occupied' && (currentStallBooking || currentStallApplication))"
                      class="issd-card-section issd-occupant-card is-reserved"
                    >
                      <!-- ส่วนผู้จอง / ผู้ยื่นคำร้อง -->
                      <div class="issd-person-row">
                        <div class="issd-person-avatar is-reserved">
                          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                        </div>
                        <div class="issd-person-info">
                          <span class="issd-person-role">{{ currentStallApplication ? 'ผู้ยื่นคำร้องขอเช่า' : 'ผู้จองเข้าชมพื้นที่' }}</span>
                          <strong class="issd-person-name">{{ currentStallBooking?.tenant_name || currentStallBooking?.applicant_name || currentStallApplication?.applicant_name || 'ผู้ขอจองเข้าชมพื้นที่' }}</strong>
                        </div>
                        <a
                          v-if="currentStallBooking?.tenant_phone || currentStallBooking?.applicant_phone || currentStallApplication?.applicant_phone"
                          :href="`tel:${currentStallBooking?.tenant_phone || currentStallBooking?.applicant_phone || currentStallApplication?.applicant_phone}`"
                          class="issd-phone-link"
                          title="คลิกเพื่อโทรออก"
                        >
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                          <span>{{ currentStallBooking?.tenant_phone || currentStallBooking?.applicant_phone || currentStallApplication?.applicant_phone }}</span>
                        </a>
                      </div>

                      <!-- ส่วนข้อมูลการจอง / คำร้อง -->
                      <div class="issd-subdoc-box is-reserved">
                        <div class="issd-subdoc-row">
                          <span class="subdoc-k">{{ currentStallApplication ? 'เลขที่คำร้อง' : 'เลขที่การจอง' }}</span>
                          <strong class="subdoc-v is-code">{{ currentStallApplication?.application_number || currentStallBooking?.booking_number || '-' }}</strong>
                        </div>
                        <div class="issd-subdoc-row">
                          <span class="subdoc-k">{{ currentStallBooking?.visit_date ? 'วันนัดชมแผง' : 'สถานะรายการ' }}</span>
                          <span class="subdoc-v">
                            <template v-if="currentStallBooking?.visit_date">
                              {{ formatThaiDateOnly(currentStallBooking.visit_date) }} {{ currentStallBooking.visit_time || currentStallBooking.visit_time_slot || '' }}
                            </template>
                            <template v-else>
                              {{ currentStallBooking ? getBookingStatusBadgeLabel(currentStallBooking.booking_status || currentStallBooking.status) : 'รอพิจารณาอนุมัติ' }}
                            </template>
                          </span>
                        </div>
                      </div>
                      <div v-if="currentStallApplication" class="issd-occupant-action-strip">
                        <NuxtLink
                          :to="withStallReturn(currentStallApplication.id ? `/applications/${currentStallApplication.id}` : `/applications?search=${encodeURIComponent(currentStallApplication.application_number || stallEdit.ident)}`)"
                          class="issd-occupant-action-btn is-app"
                        >
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                          <span>ดูคำร้องขอเช่าแผง</span>
                          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
                        </NuxtLink>
                      </div>
                      <div v-else-if="currentStallBooking" class="issd-occupant-action-strip">
                        <NuxtLink
                          :to="withStallReturn(currentStallBooking.id ? `/bookings/${currentStallBooking.id}` : `/bookings?search=${encodeURIComponent(currentStallBooking.booking_number || stallEdit.ident)}`)"
                          class="issd-occupant-action-btn is-booking"
                        >
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                          <span>ดูการจองนัดหมาย</span>
                          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
                        </NuxtLink>
                      </div>
                    </div>

                    <!-- 2.3 กรณีมีผู้เช่าตามสัญญา (occupied หรือมี active contract) -->
                    <div
                      v-else-if="stallEdit.stall_status === 'occupied' || currentStallContract"
                      class="issd-card-section issd-occupant-card is-occupied"
                    >
                      <!-- ส่วนผู้เช่า -->
                      <div class="issd-person-row">
                        <div class="issd-person-avatar is-occupied">
                          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                        </div>
                        <div class="issd-person-info">
                          <span class="issd-person-role">ผู้เช่าปัจจุบัน</span>
                          <strong class="issd-person-name">{{ currentStallContract?.tenant_name || currentStallContract?.name || (!String(stallEdit.name || '').startsWith('แผง') && stallEdit.name) || 'คุณสมศักดิ์ มั่งคั่ง' }}</strong>
                        </div>
                        <a
                          v-if="currentStallContract?.tenant_phone"
                          :href="`tel:${currentStallContract.tenant_phone}`"
                          class="issd-phone-link"
                          title="คลิกเพื่อโทรออก"
                        >
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                          <span>{{ currentStallContract.tenant_phone }}</span>
                        </a>
                      </div>

                      <!-- ส่วนข้อมูลสัญญาเช่า -->
                      <div class="issd-subdoc-box">
                        <div class="issd-subdoc-row">
                          <span class="subdoc-k">เลขที่สัญญา</span>
                          <strong class="subdoc-v is-code">{{ currentStallContract?.contract_number || '-' }}</strong>
                        </div>
                        <div class="issd-subdoc-row">
                          <span class="subdoc-k">ระยะสัญญา</span>
                          <span class="subdoc-v">
                            {{ formatThaiDateOnly(currentStallContract?.start_date) }} – {{ formatThaiDateOnly(currentStallContract?.end_date) }}
                            <em class="subdoc-duration">({{ getShortContractDurationText(currentStallContract?.start_date, currentStallContract?.end_date) }})</em>
                          </span>
                        </div>
                      </div>
                    </div>

                    <!-- 3. Operational Snapshot (4 มิติสำคัญ — แสดงข้อมูลสถิติเท่านั้น ไม่สามารถกดได้) -->
                    <div class="issd-mini-ops-grid">
                      <div
                        v-if="stallEdit.stall_status !== 'available'"
                        class="issd-mini-ops-card"
                        :class="overduePaymentCount > 0 ? 'is-alert' : 'is-ok'"
                      >
                        <div class="ops-icon">
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="5" width="20" height="14" rx="2"></rect><line x1="2" y1="10" x2="22" y2="10"></line></svg>
                        </div>
                        <div class="ops-info">
                          <span class="ops-label">การเงิน</span>
                          <strong :class="overduePaymentCount > 0 ? 'text-danger' : 'text-success'">
                            {{ overduePaymentCount > 0 ? `ค้าง ${overduePaymentCount} งวด` : 'ปกติ' }}
                          </strong>
                        </div>
                      </div>

                      <div
                        v-if="stallEdit.stall_status !== 'available'"
                        class="issd-mini-ops-card"
                        :class="isContractExpiringSoon ? 'is-warning' : (currentStallContract ? 'is-ok' : '')"
                      >
                        <div class="ops-icon">
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line></svg>
                        </div>
                        <div class="ops-info">
                          <span class="ops-label">สัญญาเช่า</span>
                          <strong :class="isContractExpiringSoon ? 'text-warning' : (currentStallContract ? 'text-success' : '')">
                            {{ currentStallContract ? (isContractExpiringSoon ? `เหลือ ${contractDaysRemaining} วัน` : 'ปกติ') : 'ไม่มี' }}
                          </strong>
                        </div>
                      </div>

                      <div
                        class="issd-mini-ops-card"
                        :class="openMaintenanceCount > 0 ? 'is-warning' : ''"
                      >
                        <div class="ops-icon">
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"></path></svg>
                        </div>
                        <div class="ops-info">
                          <span class="ops-label">งานซ่อม</span>
                          <strong :class="openMaintenanceCount > 0 ? 'text-warning' : ''">
                            {{ openMaintenanceCount > 0 ? `กำลังซ่อม ${openMaintenanceCount}` : 'ไม่มี' }}
                          </strong>
                        </div>
                      </div>

                      <div
                        class="issd-mini-ops-card"
                        :class="(currentStallWaitlistCount > 0 || currentStallBooking) ? 'is-info-tint' : ''"
                      >
                        <div class="ops-icon">
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                        </div>
                        <div class="ops-info">
                          <span class="ops-label">จอง / คิวรอ</span>
                          <strong>
                            {{ currentStallWaitlistCount > 0 ? `${currentStallWaitlistCount} คิว` : (currentStallBooking ? '1 รายการ' : 'ไม่มี') }}
                          </strong>
                        </div>
                      </div>
                    </div>

                    <!-- 4. รายละเอียดแผงเช่า (Stall Specifications & Rates) -->
                    <dl class="issd-detail-list">
                      <div class="issd-detail-row">
                        <dt>หมวดสินค้า</dt>
                        <dd>{{ stallEdit.category || selectedProject?.category || 'ทั่วไป' }}</dd>
                      </div>
                      <div class="issd-detail-row">
                        <dt>ขนาดพื้นที่</dt>
                        <dd>{{ Number(stallEdit.width) || 2 }} × {{ Number(stallEdit.depth) || 2.5 }} ม. ({{ inspectStallSizeSqm }} ตร.ม.)</dd>
                      </div>
                      <div class="issd-detail-row">
                        <dt>ค่าเช่า / แรกเข้า</dt>
                        <dd>
                          <span class="is-rent">฿{{ Number(stallEdit.monthly_rent || 1800).toLocaleString('th-TH') }}/ด.</span>
                          <span class="issd-sep-dot">·</span>
                          <span class="is-bold">แรกเข้า ฿{{ Number(inspectStallEntryFee).toLocaleString('th-TH') }}</span>
                        </dd>
                      </div>
                    </dl>
                  </div>

                  <!-- 2B. Full Detail Mode (In-Map Detail Overlay - 2-Column Sidebar Layout) -->
                  <div
                    v-else-if="stallDrawerMode === 'full-detail'"
                    class="issd-full-layout"
                  >
                    <!-- Left Sidebar Navigation -->
                    <aside class="issd-full-sidebar" role="navigation" aria-label="แถบเมนูรายละเอียด">
                      <nav class="issd-sidebar-nav" role="tablist">
                        <button
                          type="button"
                          role="tab"
                          class="issd-sidebar-tab"
                          :class="{ active: stallDetailTab === 'info' }"
                          @click="stallDetailTab = 'info'"
                        >
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18M9 21V9"/>
                          </svg>
                          <span class="issd-sidebar-tab-label">ข้อมูลแผง</span>
                        </button>
                        <button
                          type="button"
                          role="tab"
                          class="issd-sidebar-tab"
                          :class="{ active: stallDetailTab === 'contracts' }"
                          @click="stallDetailTab = 'contracts'"
                        >
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/>
                          </svg>
                          <span class="issd-sidebar-tab-label">สัญญา</span>
                          <span v-if="tabActionCounts.contracts" class="issd-sidebar-tab-count is-action" :class="{ 'is-danger': 'contracts' === 'payments' }">{{ tabActionCounts.contracts }}</span>
                        </button>
                        <button
                          type="button"
                          role="tab"
                          class="issd-sidebar-tab"
                          :class="{ active: stallDetailTab === 'subtenants' }"
                          @click="stallDetailTab = 'subtenants'"
                        >
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                          </svg>
                          <span class="issd-sidebar-tab-label">ผู้เช่าช่วง</span>
                          <span v-if="tabActionCounts.subtenants" class="issd-sidebar-tab-count is-action" :class="{ 'is-danger': 'subtenants' === 'payments' }">{{ tabActionCounts.subtenants }}</span>
                        </button>
                        <button
                          type="button"
                          role="tab"
                          class="issd-sidebar-tab"
                          :class="{ active: stallDetailTab === 'payments' }"
                          @click="stallDetailTab = 'payments'"
                        >
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect x="2" y="5" width="20" height="14" rx="2"/><line x1="2" y1="10" x2="22" y2="10"/>
                          </svg>
                          <span class="issd-sidebar-tab-label">ชำระเงิน</span>
                          <span v-if="tabActionCounts.payments" class="issd-sidebar-tab-count is-action" :class="{ 'is-danger': 'payments' === 'payments' }">{{ tabActionCounts.payments }}</span>
                        </button>
                        <button
                          type="button"
                          role="tab"
                          class="issd-sidebar-tab"
                          :class="{ active: stallDetailTab === 'maintenance' }"
                          @click="stallDetailTab = 'maintenance'"
                        >
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
                          </svg>
                          <span class="issd-sidebar-tab-label">งานซ่อม</span>
                          <span v-if="tabActionCounts.maintenance" class="issd-sidebar-tab-count is-action" :class="{ 'is-danger': 'maintenance' === 'payments' }">{{ tabActionCounts.maintenance }}</span>
                        </button>
                        <button
                          type="button"
                          role="tab"
                          class="issd-sidebar-tab"
                          :class="{ active: stallDetailTab === 'bookings' }"
                          @click="stallDetailTab = 'bookings'"
                        >
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
                          </svg>
                          <span class="issd-sidebar-tab-label">การจอง</span>
                          <span v-if="tabActionCounts.bookings" class="issd-sidebar-tab-count is-action" :class="{ 'is-danger': 'bookings' === 'payments' }">{{ tabActionCounts.bookings }}</span>
                        </button>
                        <button
                          type="button"
                          role="tab"
                          class="issd-sidebar-tab"
                          :class="{ active: stallDetailTab === 'history' }"
                          @click="stallDetailTab = 'history'"
                        >
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/><path d="M12 7v5l4 2"/>
                          </svg>
                          <span class="issd-sidebar-tab-label">ประวัติ</span>
                          <span v-if="tabActionCounts.history" class="issd-sidebar-tab-count is-action" :class="{ 'is-danger': 'history' === 'payments' }">{{ tabActionCounts.history }}</span>
                        </button>
                      </nav>
                    </aside>

                    <!-- Right Column: Main Content & Tab Footer -->
                    <div class="issd-full-main">
                      <!-- 1.1 Action Alert Card (at top of content area) -->
                      <div v-if="fullDetailSettling" class="issd-full-skeleton" aria-hidden="true">
                        <span></span><span></span><span></span>
                      </div>
                      <div
                        v-if="stallActionRichItems.length > 0 && !fullDetailSettling"
                        class="issd-action-alert-card"
                      >
                        <div
                          class="issd-aac-header"
                          @click="isActionAlertCollapsed = !isActionAlertCollapsed"
                          role="button"
                          tabindex="0"
                          @keydown.enter="isActionAlertCollapsed = !isActionAlertCollapsed"
                        >
                          <div class="issd-aac-title-group">
                            <div class="issd-aac-icon-wrap">
                              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                                <rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect>
                                <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path>
                                <path d="m9 14 2 2 4-4"></path>
                              </svg>
                            </div>
                            <div class="issd-aac-title-text">
                              <div class="issd-aac-title-row">
                                <strong>สิ่งที่ต้องดำเนินการ</strong>
                                <span class="issd-aac-count-badge">{{ stallActionRichItems.length }} รายการ</span>
                              </div>
                              
                            </div>
                          </div>
                          <button
                            type="button"
                            class="issd-aac-toggle-btn"
                            :aria-label="isActionAlertCollapsed ? 'ขยายรายการงาน' : 'ย่อรายการงาน'"
                            :title="isActionAlertCollapsed ? 'คลิกเพื่อขยาย' : 'คลิกเพื่อย่อ'"
                          >
                            <svg
                              width="14"
                              height="14"
                              viewBox="0 0 24 24"
                              fill="none"
                              stroke="currentColor"
                              stroke-width="2.2"
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              :class="{ 'is-flipped': isActionAlertCollapsed }"
                            >
                              <polyline points="6 9 12 15 18 9"></polyline>
                            </svg>
                          </button>
                        </div>

                        <div v-show="!isActionAlertCollapsed" class="issd-aac-list">
                          <div
                            v-for="act in stallActionRichItems"
                            :key="act.key"
                            class="issd-aac-row"
                            :class="`is-${act.severity}`"
                            @click="jumpToActionTab(act.targetTab)"
                            role="button"
                            tabindex="0"
                            :title="`คลิกเพื่อไปยังแท็บ: ${act.actionLabel}`"
                            @keydown.enter="jumpToActionTab(act.targetTab)"
                          >
                            <div class="issd-aac-row-left">
                              <span class="issd-aac-bar" :class="`is-${act.severity}`" aria-hidden="true"></span>
                              <div class="issd-aac-row-text">
                                <strong>{{ act.category }}<b v-if="act.count" class="issd-aac-count" :class="`is-${act.severity}`">{{ act.count }}</b></strong>
                                <small>{{ act.title }}</small>
                              </div>
                            </div>
                            <div class="issd-aac-row-right">
                              <span class="issd-aac-action-text">{{ act.actionLabel }}</span>
                              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round">
                                <polyline points="9 18 15 12 9 6"></polyline>
                              </svg>
                            </div>
                          </div>
                        </div>
                      </div>

                      <!-- Scrollable Tab Content Body -->
                      <div
                        v-if="!fullDetailSettling"
                        class="issd-body issd-body-full"
                        :class="{ 'is-fixed-tab-layout': stallDetailTab !== 'info' }"
                      >
                    <!-- Tab: ข้อมูลแผงเช่า (info) -->
                    <template v-if="stallDetailTab === 'info'">
                      <!-- 1. กรณีแผงว่าง (available) — กล่องเปิดรับจอง อยู่บนสุด ใต้สิ่งที่ต้องดำเนินการ -->
                      <div
                        v-if="stallEdit.stall_status === 'available'"
                        class="issd-card-section issd-occupant-card is-available"
                      >
                        <div class="issd-available-line">
                          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
                          <strong>พร้อมเปิดรับจอง</strong>
                        </div>
                      </div>

                      <!-- Photo Banner -->
                      <div
                        class="issd-photo-banner is-full"
                        title="คลิกเพื่อดูภาพแผงเช่าทั้งหมด"
                        @click="openInspectPhotoPreview(0)"
                      >
                        <img :src="inspectStallPhotos[0]" :alt="`ภาพแผง ${stallEdit.ident}`" />
                        <div class="issd-photo-count-badge">
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                          <span>{{ inspectStallPhotos.length }} รูป (คลิกดูภาพขยาย)</span>
                        </div>
                      </div>

                      <!-- KPI Strip 3 ช่อง (ข้อมูลหลักเท่านั้น) -->
                      <div class="issd-kpi-grid is-three">
                        <div class="issd-kpi-card">
                          <span>ขนาดพื้นที่</span>
                          <strong>{{ inspectStallSizeSqm }} ตร.ม.</strong>
                          <em>{{ stallEdit.width || 2 }} × {{ stallEdit.depth || 2.5 }} ม.</em>
                        </div>
                        <div class="issd-kpi-card is-green">
                          <span>ค่าเช่า</span>
                          <strong>฿{{ Number(stallEdit.monthly_rent || 1800).toLocaleString('th-TH') }}</strong>
                          <em>ต่อเดือน</em>
                        </div>
                        <div class="issd-kpi-card">
                          <span>ค่าแรกเข้า</span>
                          <strong>฿{{ Number(inspectStallEntryFee).toLocaleString('th-TH') }}</strong>
                          <em>ครั้งเดียว</em>
                        </div>
                      </div>

                      <!-- Detailed Specifications (กล่องรายละเอียดสไตล์ KPI Strip — ไม่มีหัวข้อ) -->
                      <div class="issd-kpi-grid is-three">
                        <div v-if="stallEdit.name && !String(stallEdit.name).startsWith('แผง')" class="issd-kpi-card">
                          <span>ชื่อแผง</span>
                          <strong>{{ stallEdit.name }}</strong>
                          <em>ป้ายร้านค้า</em>
                        </div>
                        <div class="issd-kpi-card">
                          <span>ชั้น</span>
                          <strong>{{ formatFloorDisplay(stallEdit.floor) }}</strong>
                          <em>{{ stallEdit.building ? `อาคาร ${stallEdit.building}` : 'ระดับพื้นที่' }}</em>
                        </div>
                        <div class="issd-kpi-card">
                          <span>โซน</span>
                          <strong>{{ stallEdit.zone_name || selectedProject?.name || 'โซนกลาง' }}</strong>
                          <em>ผังโครงการ</em>
                        </div>
                        <div class="issd-kpi-card">
                          <span>หมวดสินค้า</span>
                          <strong>{{ stallEdit.category || selectedProject?.category || 'แผงค้ามาตรฐาน' }}</strong>
                          <em>ประเภทสินค้า</em>
                        </div>
                        <div v-if="stallEdit.description" class="issd-kpi-card is-full-span">
                          <span>รายละเอียดเพิ่มเติม</span>
                          <strong class="is-desc">{{ stallEdit.description }}</strong>
                        </div>
                      </div>

                      <!-- กล่องข้อมูลผู้เช่า / ผู้จอง / สถานะอื่นๆ (ย้ายมาอยู่ข้างล่างสุด) -->
                      <!-- 2. กรณีติดจอง (reserved หรือมีรายการจอง/คำร้อง) -->
                      <div
                        v-if="stallEdit.stall_status === 'reserved' || (stallEdit.stall_status !== 'occupied' && stallEdit.stall_status !== 'available' && (currentStallBooking || currentStallApplication))"
                        class="issd-card-section issd-occupant-card is-reserved"
                      >
                        <!-- ส่วนผู้จอง / ผู้ยื่นคำร้อง -->
                        <div class="issd-person-row">
                          <div class="issd-person-avatar is-reserved">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                          </div>
                          <div class="issd-person-info">
                            <span class="issd-person-role">{{ currentStallApplication ? 'ผู้ยื่นคำร้องขอเช่า' : 'ผู้จองเข้าชมพื้นที่' }}</span>
                            <strong class="issd-person-name">{{ currentStallBooking?.tenant_name || currentStallBooking?.applicant_name || currentStallApplication?.applicant_name || 'ผู้ขอจองเข้าชมพื้นที่' }}</strong>
                          </div>
                          <a
                            v-if="currentStallBooking?.tenant_phone || currentStallBooking?.applicant_phone || currentStallApplication?.applicant_phone"
                            :href="`tel:${currentStallBooking?.tenant_phone || currentStallBooking?.applicant_phone || currentStallApplication?.applicant_phone}`"
                            class="issd-phone-link"
                            title="คลิกเพื่อโทรออก"
                          >
                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                            <span>{{ currentStallBooking?.tenant_phone || currentStallBooking?.applicant_phone || currentStallApplication?.applicant_phone }}</span>
                          </a>
                        </div>

                        <!-- ส่วนข้อมูลการจอง / คำร้อง -->
                        <div class="issd-subdoc-box is-reserved">
                          <div class="issd-subdoc-row">
                            <span class="subdoc-k">{{ currentStallApplication ? 'เลขที่คำร้อง' : 'เลขที่การจอง' }}</span>
                            <strong class="subdoc-v is-code">{{ currentStallApplication?.application_number || currentStallBooking?.booking_number || '-' }}</strong>
                          </div>
                          <div class="issd-subdoc-row">
                            <span class="subdoc-k">{{ currentStallBooking?.visit_date ? 'วันนัดชมแผง' : 'สถานะรายการ' }}</span>
                            <span class="subdoc-v">
                              <template v-if="currentStallBooking?.visit_date">
                                {{ formatThaiDateOnly(currentStallBooking.visit_date) }} {{ currentStallBooking.visit_time || currentStallBooking.visit_time_slot || '' }}
                              </template>
                              <template v-else>
                                {{ currentStallBooking ? getBookingStatusBadgeLabel(currentStallBooking.booking_status || currentStallBooking.status) : 'รอพิจารณาอนุมัติ' }}
                              </template>
                            </span>
                          </div>
                        </div>
                        <div v-if="currentStallApplication" class="issd-occupant-action-strip">
                          <NuxtLink
                            :to="withStallReturn(currentStallApplication.id ? `/applications/${currentStallApplication.id}` : `/applications?search=${encodeURIComponent(currentStallApplication.application_number || stallEdit.ident)}`)"
                            class="issd-occupant-action-btn is-app"
                          >
                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                            <span>ดูคำร้องขอเช่าแผง</span>
                            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
                          </NuxtLink>
                        </div>
                        <div v-else-if="currentStallBooking" class="issd-occupant-action-strip">
                          <NuxtLink
                            :to="withStallReturn(currentStallBooking.id ? `/bookings/${currentStallBooking.id}` : `/bookings?search=${encodeURIComponent(currentStallBooking.booking_number || stallEdit.ident)}`)"
                            class="issd-occupant-action-btn is-booking"
                          >
                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                            <span>ดูการจองนัดหมาย</span>
                            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
                          </NuxtLink>
                        </div>
                      </div>

                      <!-- 3. กรณีมีผู้เช่าตามสัญญา (occupied หรือมี active contract) -->
                      <div
                        v-else-if="stallEdit.stall_status === 'occupied' || currentStallContract"
                        class="issd-card-section issd-occupant-card is-occupied"
                      >
                        <div class="issd-person-row">
                          <div class="issd-person-avatar is-occupied">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                          </div>
                          <div class="issd-person-info">
                            <span class="issd-person-role">ผู้เช่าปัจจุบัน</span>
                            <strong class="issd-person-name">{{ currentStallContract?.tenant_name || currentStallContract?.name || (!String(stallEdit.name || '').startsWith('แผง') && stallEdit.name) || '-' }}</strong>
                          </div>
                          <a
                            v-if="currentStallContract?.tenant_phone"
                            :href="`tel:${currentStallContract.tenant_phone}`"
                            class="issd-phone-link"
                            title="คลิกเพื่อโทรออก"
                          >
                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                            <span>{{ currentStallContract.tenant_phone }}</span>
                          </a>
                        </div>
                        <div class="issd-subdoc-box">
                          <div class="issd-subdoc-row">
                            <span class="subdoc-k">เลขที่สัญญา</span>
                            <strong class="subdoc-v is-code">{{ currentStallContract?.contract_number || '-' }}</strong>
                          </div>
                          <div class="issd-subdoc-row">
                            <span class="subdoc-k">ระยะสัญญา</span>
                            <span class="subdoc-v">
                              {{ formatThaiDateOnly(currentStallContract?.start_date) }} – {{ formatThaiDateOnly(currentStallContract?.end_date) }}
                              <em class="subdoc-duration">({{ getShortContractDurationText(currentStallContract?.start_date, currentStallContract?.end_date) }})</em>
                            </span>
                          </div>
                        </div>
                      </div>

                      <div v-else-if="stallEdit.stall_status && stallEdit.stall_status !== 'available'" class="issd-occupancy-box" :class="`is-${stallEdit.stall_status}`">
                        <div class="issd-occupancy-head">
                          <span>{{ inspectStallOccupancy.heading }}</span>
                          <small v-if="inspectStallOccupancy.refCode">{{ inspectStallOccupancy.refCode }}</small>
                        </div>
                        <strong>{{ inspectStallOccupancy.primaryText }}</strong>
                        <p v-if="inspectStallOccupancy.subText">{{ inspectStallOccupancy.subText }}</p>
                      </div>
                    </template>

                    <!-- Tab: ประวัติการดำเนินการ (history) -->
                    <template v-else-if="stallDetailTab === 'history'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">ประวัติ</strong>
                        </div>
                      </div>

                      <div v-if="stallActivityTimeline.length > 0" class="issd-scroll-list">
                        <div class="issd-timeline-list issd-timeline-panel">
                          <div v-for="event in stallActivityTimeline" :key="event.key" class="issd-timeline-item">
                            <span class="issd-timeline-dot" :class="event.dotClass || ''"></span>
                            <div class="issd-timeline-copy">
                              <div class="issd-timeline-head-row">
                                <strong>{{ event.title }}</strong>
                                <span v-if="event.category" class="issd-timeline-cat-pill">{{ event.category }}</span>
                              </div>
                              <time>{{ event.date }}</time>
                              <p>{{ event.detail }}</p>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/><path d="M12 7v5l4 2"/>
                        </svg>
                        <p>ไม่มีประวัติ</p>
                      </div>
                    </template>

                    <!-- Tab: การจองและคำร้อง (bookings) -->
                    <template v-else-if="stallDetailTab === 'bookings'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">การจอง</strong>
                        </div>
                      </div>

                      <AppSegmentedControl
                        v-model="bookingTypeFilter"
                        :options="[
                          { value: 'all', label: `ทั้งหมด (${stallBookingsAndApps.length})` },
                          { value: 'booking', label: `การจอง (${stallBookingsAndApps.filter(i => i.typeKey === 'booking').length})` },
                          { value: 'application', label: `คำร้อง (${stallBookingsAndApps.filter(i => i.typeKey === 'application').length})` }
                        ]"
                        aria-label="กรองตามประเภทรายการ"
                      />

                      <div v-if="filteredBookingsAndApps.length > 0" class="issd-booking-list issd-scroll-list">
                        <div
                          v-for="item in filteredBookingsAndApps"
                          :key="item.key"
                          class="issd-booking-card"
                        >
                          <div class="issd-booking-card-header">
                            <div class="issd-booking-type-group">
                              <span class="issd-card-type-tag" :class="`is-${item.typeKey || 'booking'}`">
                                <svg v-if="item.typeKey === 'application'" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                                <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                                <span>{{ item.typeLabel }}</span>
                              </span>
                              <span class="issd-booking-ref-code">{{ item.refNumber }}</span>
                            </div>
                            <span class="issd-status-pill" :class="item.badgeClass">{{ item.statusText }}</span>
                          </div>

                          <div class="issd-booking-applicant-row">
                            <div class="issd-booking-applicant-left">
                              <div class="issd-booking-avatar">
                                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                              </div>
                              <div class="issd-booking-applicant-info">
                                <span class="issd-booking-role-label">{{ item.typeKey === 'application' ? 'ผู้ยื่นคำร้อง' : 'ผู้จอง' }}</span>
                                <strong class="issd-card-title">{{ item.applicantName }}</strong>
                              </div>
                            </div>
                            <a
                              v-if="item.phone"
                              :href="`tel:${item.phone}`"
                              class="issd-phone-link"
                              title="คลิกเพื่อโทรออก"
                            >
                              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
                              <span>{{ item.phone }}</span>
                            </a>
                          </div>

                          <!-- Footer Bar: วันนัดชม/วันที่ทำรายการ (ซ้าย) + ปุ่มดูรายละเอียด (ขวา) -->
                          <div class="issd-booking-footer-bar">
                            <div class="issd-booking-info-item">
                              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                              <span><strong>{{ item.dateText }}</strong></span>
                            </div>
                            <NuxtLink
                              :to="withStallReturn(item.targetUrl)"
                              class="issd-booking-direct-btn"
                              :class="`is-${item.typeKey}`"
                            >
                              <svg v-if="item.typeKey === 'application'" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                              <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                              <span>ดูรายละเอียด</span>
                              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"/></svg>
                            </NuxtLink>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <circle cx="12" cy="12" r="10"></circle>
                          <line x1="12" y1="8" x2="12" y2="12"></line>
                          <line x1="12" y1="16" x2="12.01" y2="16"></line>
                        </svg>
                        <p>{{ bookingTypeFilter === 'application' ? 'ไม่มีคำร้อง' : (bookingTypeFilter === 'booking' ? 'ไม่มีรายการจอง' : 'ไม่มีรายการจองหรือคำร้อง') }}</p>
                      </div>
                    </template>

                    <!-- Tab: สัญญา (contracts) -->
                    <template v-else-if="stallDetailTab === 'contracts'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">สัญญา</strong>
                        </div>
                      </div>

                      <AppSegmentedControl
                        v-model="contractSubTab"
                        :options="[{ value: 'active', label: `ปัจจุบัน (${activeStallContracts.length})` }, { value: 'history', label: `ประวัติ (${historyStallContracts.length})` }]"
                        aria-label="เลือกประเภทสัญญา"
                      />

                      <div v-if="displayedStallContracts.length > 0" class="issd-contract-list issd-scroll-list">
                        <div
                          v-for="contract in displayedStallContracts"
                          :key="contract.id || contract.contract_number"
                          class="issd-contract-card"
                          :class="{ 'is-current': contract.status === 'active' }"
                        >
                          <div class="issd-card-top-row">
                            <div class="issd-contract-top-left">
                              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="issd-contract-icon">
                                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                                <polyline points="14 2 14 8 20 8"></polyline>
                              </svg>
                              <strong class="issd-contract-num">{{ contract.contract_number || `#${contract.id}` }}</strong>
                            </div>
                            <span
                              class="issd-status-pill"
                              :class="contract.status === 'active' ? 'is-available' : ((contract.status === 'cancelled' || contract.status === 'terminated') ? 'is-cancelled' : 'is-reserved')"
                            >
                              {{ contract.status === 'active' ? 'ใช้งานอยู่' : (contract.status === 'expired' ? 'สิ้นสุดสัญญา' : ((contract.status === 'cancelled' || contract.status === 'terminated') ? 'ยกเลิกแล้ว' : (contract.status || 'ร่างสัญญา'))) }}
                            </span>
                          </div>

                          <div class="issd-contract-details-grid">
                            <div class="issd-contract-detail-item">
                              <span class="detail-label">ผู้เช่า</span>
                              <strong class="detail-value">{{ contract.tenant_name || contract.name || '-' }}</strong>
                            </div>
                            <div class="issd-contract-detail-item">
                              <span class="detail-label">ค่าเช่า</span>
                              <strong class="detail-value is-rent">฿{{ Number(contract.monthly_rental || contract.monthly_rent || 0).toLocaleString('th-TH') }}/ด.</strong>
                            </div>
                            <div v-if="contract.security_deposit || contract.deposit_amount" class="issd-contract-detail-item">
                              <span class="detail-label">เงินประกัน</span>
                              <strong class="detail-value">฿{{ Number(contract.security_deposit || contract.deposit_amount).toLocaleString('th-TH') }}</strong>
                            </div>
                            <div class="issd-contract-detail-item is-wide">
                              <span class="detail-label">ระยะสัญญา</span>
                              <strong class="detail-value">{{ formatThaiDateOnly(contract.start_date) }} – {{ formatThaiDateOnly(contract.end_date) }} <em class="detail-sub">({{ getShortContractDurationText(contract.start_date, contract.end_date) }})</em></strong>
                            </div>
                          </div>

                          <div class="issd-card-actions-row">
                            <AppButton size="sm" variant="outline" @click="openEditContractModal(contract)">แก้ไขสัญญา</AppButton>
                            <AppButton v-if="contract.status === 'active'"
                              size="sm" variant="danger-outline" @click="terminateContract(contract)">ยกเลิกสัญญา</AppButton>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                          <polyline points="14 2 14 8 20 8"></polyline>
                        </svg>
                        <p>{{ contractSubTab === 'active' ? 'ไม่มีสัญญาปัจจุบัน' : 'ไม่มีประวัติสัญญา' }}</p>
                      </div>
                    </template>

                    <!-- Tab: ผู้เช่าช่วง (subtenants) -->
                    <template v-else-if="stallDetailTab === 'subtenants'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">ผู้เช่าช่วง</strong>
                        </div>
                      </div>

                      <div v-if="stallSubtenants.length > 0" class="issd-subtenant-list issd-scroll-list">
                        <div
                          v-for="sub in stallSubtenants"
                          :key="sub.id"
                          class="issd-subtenant-card"
                        >
                          <div class="issd-card-top-row">
                            <div class="issd-subtenant-header-info">
                              <div v-if="sub.photoUrl" class="issd-subtenant-avatar">
                                <img :src="sub.photoUrl" :alt="sub.name" />
                              </div>
                              <div v-else class="issd-subtenant-avatar is-placeholder">
                                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                  <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/>
                                </svg>
                              </div>
                              <div class="issd-subtenant-text-block">
                                <strong class="issd-card-title">{{ sub.name }}</strong>
                              </div>
                            </div>
                            <span class="issd-status-pill" :class="sub.status === 'active' ? 'is-available' : 'is-occupied'">
                              {{ sub.status === 'active' ? 'เช่าอยู่' : 'สิ้นสุด' }}
                            </span>
                          </div>
                          <div class="issd-subtenant-details">
                            <div class="issd-detail-item">
                              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
                              <span>{{ sub.phone }}</span>
                            </div>
                            <div v-if="sub.businessType" class="issd-detail-item">
                              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 2 3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4Z"/><path d="M3 6h18"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
                              <span>{{ sub.businessType }}</span>
                            </div>
                          </div>
                          <div class="issd-card-meta-row">
                            <span>{{ sub.startDate || '-' }} – {{ sub.endDate || '-' }}</span>
                            <span v-if="sub.subRentAmount" class="is-rent">฿{{ Number(sub.subRentAmount).toLocaleString('th-TH') }}/ด.</span>
                          </div>
                          <div class="issd-card-actions-row">
                            <AppButton size="sm" variant="secondary" @click="toggleSubtenantStatus(sub.id)">{{ sub.status === 'active' ? 'สิ้นสุดสัญญา' : 'กลับมาเช่าต่อ' }}</AppButton>
                            <AppButton size="sm" variant="danger-outline" @click="deleteSubtenant(sub.id)">ลบ</AppButton>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/>
                        </svg>
                        <p>ไม่มีผู้เช่าช่วง</p>
                      </div>
                    </template>

                    <!-- Tab: ประวัติชำระเงิน (payments) -->
                    <template v-else-if="stallDetailTab === 'payments'">
                      <div class="issd-tab-header-row is-inline-controls">
                        <div class="issd-tab-header-copy">
                          <strong class="issd-tab-title">ชำระเงิน</strong>
                        </div>
                        <div v-if="stallContractsList.length" class="issd-contract-inline-filter">
                          <span class="issd-inline-filter-label">สัญญา</span>
                          <div class="issd-inline-filter-select">
                            <AppSelect v-model="selectedPaymentContract" :options="paymentContractOptions" />
                          </div>
                        </div>
                      </div>

                      <div v-if="stallPayments.length" class="issd-finance-summary">
                        <div><span>เรียกเก็บ</span><strong>฿{{ paymentTotal.toLocaleString('th-TH') }}</strong></div>
                        <div><span>ชำระแล้ว</span><strong class="is-success">฿{{ paidPaymentTotal.toLocaleString('th-TH') }}</strong></div>
                        <div><span>คงค้าง</span><strong :class="{ 'is-alert': outstandingPaymentTotal > 0 }">฿{{ outstandingPaymentTotal.toLocaleString('th-TH') }}</strong></div>
                      </div>

                      <div v-if="filteredStallPayments.length > 0" class="issd-payment-list issd-scroll-list">
                        <div
                          v-for="pay in filteredStallPayments"
                          :key="pay.id"
                          class="issd-payment-card"
                        >
                          <div class="issd-card-top-row">
                            <div class="issd-pay-period">
                              <strong>{{ pay.periodMonth }}</strong>
                            </div>
                            <span class="issd-status-pill" :class="pay.status === 'paid' ? 'is-available' : (pay.status === 'overdue' ? 'is-occupied' : 'is-reserved')">
                              {{ pay.status === 'paid' ? 'ชำระแล้ว' : (pay.status === 'overdue' ? 'เกินกำหนด' : 'รอชำระ') }}
                            </span>
                          </div>
                          <div class="issd-pay-amount-row">
                            <span class="issd-pay-method">{{ formatPaymentMethod(pay.paymentMethod) }}<template v-if="pay.receiptNumber"> · {{ pay.receiptNumber }}</template></span>
                            <strong class="is-rent">฿{{ Number(pay.amount).toLocaleString('th-TH') }}</strong>
                          </div>
                          <div class="issd-card-meta-row">
                            <span>{{ pay.paidDate || 'ยังไม่ชำระ' }}</span>
                            <AppButton size="sm" variant="danger-outline" @click="deletePayment(pay.id)">ลบ</AppButton>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <rect x="2" y="5" width="20" height="14" rx="2"/><line x1="2" y1="10" x2="22" y2="10"/>
                        </svg>
                        <p>ไม่มีประวัติการชำระเงิน</p>
                      </div>
                    </template>

                    <!-- Tab: ประวัติการซ่อม / ขออนุญาต (maintenance) -->
                    <template v-else-if="stallDetailTab === 'maintenance'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">งานซ่อม</strong>
                        </div>
                      </div>

                      <AppSegmentedControl
                        v-model="maintenanceSourceFilter"
                        :options="[{ value: 'all', label: 'ทั้งหมด' }, { value: 'tenant', label: 'ผู้เช่า' }, { value: 'staff', label: 'เจ้าหน้าที่' }]"
                        aria-label="กรองรายการตามผู้แจ้ง"
                      />

                      <div v-if="filteredStallMaintenance.length > 0" class="issd-maintenance-list issd-scroll-list">
                        <div
                          v-for="m in filteredStallMaintenance"
                          :key="m.id"
                          class="issd-maintenance-card"
                        >
                          <div class="issd-card-top-row">
                            <div class="issd-maintenance-tags">
                              <span class="issd-maintenance-type-pill" :class="`type-${m.type}`">{{ formatMaintenanceType(m.type) }}</span>
                            </div>
                            <span class="issd-status-pill" :class="m.status === 'completed' ? 'is-available' : (m.status === 'rejected' ? 'is-occupied' : 'is-reserved')">
                              {{ formatMaintenanceStatus(m.status) }}
                            </span>
                          </div>
                          <strong class="issd-card-title">{{ m.title }}</strong>
                          <div class="issd-maintenance-meta-grid">
                            <div><span>ผู้แจ้ง</span> <strong>{{ m.requesterName }}</strong></div>
                            <div><span>วันที่</span> <strong>{{ m.requestDate }}</strong></div>
                            <div v-if="m.estimatedCost"><span>ค่าใช้จ่าย</span> <strong class="is-rent">฿{{ Number(m.estimatedCost).toLocaleString('th-TH') }}</strong></div>
                          </div>
                          <div class="issd-card-actions-row">
                            <AppButton v-if="m.status !== 'completed'"
                              size="sm" variant="success" @click="updateMaintenanceStatus(m.id, 'completed')">เสร็จสิ้น</AppButton>
                            <AppButton size="sm" variant="danger-outline" @click="deleteMaintenance(m.id)">ลบ</AppButton>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
                        </svg>
                        <p>ไม่มีงานซ่อม</p>
                      </div>
                    </template>
                  </div>

                  <!-- 3B. Tab Action Footer (Inside Full Main) -->
                  <div v-if="!fullDetailSettling && stallDetailTab !== 'info' && stallDetailTab !== 'history'" class="issd-footer issd-footer-full">
                    <!-- Tab: ข้อมูลแผง (info) / ประวัติการดำเนินการ (history) -->
                    <template v-if="stallDetailTab === 'info' || stallDetailTab === 'history'">
                      <AppButton
                        variant="secondary"
                        fullWidth
                        @click="openStallEditorFromOverlay"
                      >
                        <template #prefix>
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M12 20h9"></path>
                            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path>
                          </svg>
                        </template>
                        แก้ไขข้อมูลแผงบนผัง
                      </AppButton>
                      <AppButton
                        v-if="stallEdit.stall_status !== 'available'"
                        variant="primary"
                        fullWidth
                        @click="stallDetailTab = 'contracts'"
                      >
                        <template #prefix>
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                            <polyline points="14 2 14 8 20 8"></polyline>
                            <line x1="16" y1="13" x2="8" y2="13"></line>
                            <line x1="16" y1="17" x2="8" y2="17"></line>
                          </svg>
                        </template>
                        จัดการสัญญา
                      </AppButton>
                    </template>

                    <!-- Tab: สัญญา (contracts) -->
                    <template v-else-if="stallDetailTab === 'contracts'">
                      <AppButton
                        variant="primary"
                        fullWidth
                        @click="openCreateContractModal"
                      >
                        <template #prefix>
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <line x1="12" y1="5" x2="12" y2="19"></line>
                            <line x1="5" y1="12" x2="19" y2="12"></line>
                          </svg>
                        </template>
                        เพิ่มสัญญา
                      </AppButton>
                    </template>

                    <!-- Tab: ผู้เช่าช่วง (subtenants) -->
                    <template v-else-if="stallDetailTab === 'subtenants'">
                      <AppButton
                        variant="primary"
                        fullWidth
                        @click="openAddSubtenantModal"
                      >
                        <template #prefix>
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <line x1="12" y1="5" x2="12" y2="19"></line>
                            <line x1="5" y1="12" x2="19" y2="12"></line>
                          </svg>
                        </template>
                        เพิ่มผู้เช่าช่วง
                      </AppButton>
                    </template>

                    <!-- Tab: ประวัติชำระเงิน (payments) -->
                    <template v-else-if="stallDetailTab === 'payments'">
                      <AppButton
                        variant="primary"
                        fullWidth
                        @click="openAddPaymentModal"
                      >
                        <template #prefix>
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <line x1="12" y1="5" x2="12" y2="19"></line>
                            <line x1="5" y1="12" x2="19" y2="12"></line>
                          </svg>
                        </template>
                        บันทึกการชำระเงิน
                      </AppButton>
                    </template>

                    <!-- Tab: ซ่อม/ขออนุญาต (maintenance) -->
                    <template v-else-if="stallDetailTab === 'maintenance'">
                      <AppButton
                        variant="primary"
                        fullWidth
                        @click="openAddMaintenanceModal"
                      >
                        <template #prefix>
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <line x1="12" y1="5" x2="12" y2="19"></line>
                            <line x1="5" y1="12" x2="19" y2="12"></line>
                          </svg>
                        </template>
                        เพิ่มการแจ้งซ่อม / ขออนุญาต
                      </AppButton>
                    </template>

                    <!-- Tab: การจองและคำร้อง (bookings) -->
                    <template v-else-if="stallDetailTab === 'bookings'">
                      <div class="issd-dual-btn-group">
                        <AppButton
                          variant="secondary"
                          :to="withStallReturn('/bookings')"
                          size="md"
                        >
                          <template #prefix>
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                              <line x1="16" y1="2" x2="16" y2="6"></line>
                              <line x1="8" y1="2" x2="8" y2="6"></line>
                              <line x1="3" y1="10" x2="21" y2="10"></line>
                            </svg>
                          </template>
                          ดูรายการจองทั้งหมด
                        </AppButton>
                        <AppButton
                          variant="primary"
                          :to="withStallReturn('/applications')"
                          size="md"
                        >
                          <template #prefix>
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                              <polyline points="14 2 14 8 20 8"></polyline>
                            </svg>
                          </template>
                          ดูรายการคำร้องทั้งหมด
                        </AppButton>
                      </div>
                    </template>
                  </div>
                </div>
              </div>

                  <!-- View: Edit in Map -->
                  <div v-else-if="stallDrawerMode === 'edit'" class="issd-body issd-body-edit">
                    <form class="issd-form-grid" novalidate @submit.prevent="saveStall">
                      <div class="full issd-photo-editor">
                        <div class="issd-photo-editor-head">
                          <span>รูปภาพแผง</span>
                          <small>{{ editablePhotos.length }}/{{ MAX_STALL_PHOTOS }} รูป</small>
                        </div>
                        <div class="issd-photo-editor-grid">
                          <div v-for="(src, idx) in editablePhotos" :key="`${idx}-${src.slice(-16)}`" class="issd-photo-thumb">
                            <img :src="src" :alt="`ภาพแผง ${idx + 1}`" />
                            <button type="button" class="issd-photo-remove" aria-label="ลบรูปนี้" title="ลบรูปนี้" @click="removeStallPhoto(idx)">
                              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                            </button>
                          </div>
                          <label v-if="editablePhotos.length < MAX_STALL_PHOTOS" class="issd-photo-add">
                            <input type="file" accept="image/*" multiple hidden @change="handleStallPhotoUpload" />
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                            <span>เพิ่มรูป</span>
                          </label>
                        </div>
                      </div>
                      <label>
                        <span>รหัสแผง <b v-if="isCreatingNewStall" class="req-star">*</b></span>
                        <AppInput v-model="stallEdit.ident" placeholder="เช่น A-001" :disabled="!isCreatingNewStall" :error="stallErrors.ident" @blur="validateStallField('ident')" />
                      </label>
                      <label>
                        <span>ชื่อแผง</span>
                        <AppInput v-model="stallEdit.name" placeholder="ระบุชื่อแผง (ถ้ามี)" />
                      </label>
                      <label>
                        <span>โซน</span>
                        <AppSelect v-model="stallEdit.zone_id" :options="zoneOptions" :disabled="!isCreatingNewStall" />
                      </label>
                      <label>
                        <span>โครงการ</span>
                        <AppInput v-model="stallEdit.project_name" placeholder="ชื่อโครงการ" disabled />
                      </label>
                      <label>
                        <span>ชั้น</span>
                        <AppSelect v-model="stallEdit.floor" :options="floorOptions" :disabled="(selectedProject?.floors?.length || 1) <= 1" />
                      </label>
                      <label>
                        <span>สถานะ <b class="req-star">*</b></span>
                        <AppSelect v-model="stallEdit.stall_status" :options="stallStatusOptions" :disabled="stallEdit.stall_status === 'occupied'" />
                      </label>
                      <label>
                        <span>กว้าง (เมตร) <b class="req-star">*</b></span>
                        <AppInput v-model="stallEdit.width" type="number" placeholder="เช่น 2.0" :error="stallErrors.width" @blur="validateStallField('width')" />
                      </label>
                      <label>
                        <span>ลึก (เมตร) <b class="req-star">*</b></span>
                        <AppInput v-model="stallEdit.depth" type="number" placeholder="เช่น 2.0" :error="stallErrors.depth" @blur="validateStallField('depth')" />
                      </label>
                      <label>
                        <span>ค่าเช่ารายวัน (บาท)</span>
                        <AppInput v-model="stallEdit.daily_rent" type="number" placeholder="เช่น 150" />
                      </label>
                      <label>
                        <span>ค่าเช่ารายเดือน (บาท) <b class="req-star">*</b></span>
                        <AppInput v-model="stallEdit.monthly_rent" type="number" placeholder="เช่น 1800" :error="stallErrors.monthly_rent" @blur="validateStallField('monthly_rent')" />
                      </label>
                      <label class="full">
                        <span>หมวดสินค้า</span>
                        <AppSelect
                          v-model="stallEdit.category"
                          :options="stallCategoryOptions"
                          placeholder="เลือกหมวดสินค้า"
                        />
                      </label>
                      <label class="full">
                        <span>รายละเอียด</span>
                        <textarea v-model.trim="stallEdit.description" rows="2" placeholder="รายละเอียดเพิ่มเติมของแผง"></textarea>
                      </label>
                    </form>
                  </div>

                  <!-- 3. Footer for Detail (Mini) and Edit Modes -->
                  <div v-if="stallDrawerMode === 'detail'" class="issd-footer">
                    <button
                      type="button"
                      class="issd-cta-btn"
                      @click="isInspectMode ? openStallFullDetail() : openStallEditorFromOverlay()"
                    >
                      <span>{{ isInspectMode ? 'ดูรายละเอียด' : 'แก้ไขข้อมูลแผง' }}</span>
                      <svg v-if="isInspectMode" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                        <line x1="5" y1="12" x2="19" y2="12"></line>
                        <polyline points="12 5 19 12 12 19"></polyline>
                      </svg>
                      <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M12 20h9"></path>
                        <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path>
                      </svg>
                    </button>
                  </div>
                  <div v-else-if="stallDrawerMode === 'edit'" class="issd-footer issd-footer-edit">
                    <AppButton
                      v-if="!isCreatingNewStall"
                      variant="danger-outline"
                      size="sm"
                      @click="deleteStall()"
                    >
                      <template #prefix>
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M10 11v5M14 11v5"/>
                        </svg>
                      </template>
                      ลบแผง
                    </AppButton>
                    <div class="issd-footer-right">
                      <AppButton variant="secondary" size="sm" @click="cancelStallEdit">
                        ยกเลิก
                      </AppButton>
                      <AppButton :loading="savingStall" size="sm" @click="saveStall">
                        {{ isCreatingNewStall ? 'เพิ่มแผงเช่า' : 'บันทึกการแก้ไข' }}
                      </AppButton>
                    </div>
                  </div>
              </div>
            </template>
          </MarketProjectsOverviewMap>
        </article>
      </section>
    </main>

    <Teleport to="body">
      <AppModal
        v-model="showZoneEditor"
        title="เพิ่มโซนบนผัง"
        :description="`พื้นที่ ${drawnZoneCoordinates.length} จุดที่วาดบนแผนที่จะถูกใช้เป็นขอบเขตโซน`"
      >
        <form class="zone-form" novalidate @submit.prevent="createZone">
          <label>
            <span>ชื่อโซน <b class="req-star">*</b></span>
            <AppInput v-model="zoneForm.name" placeholder="เช่น โซนร้านอาหารใหม่" :error="zoneErrors.name" @blur="validateZoneField('name')" />
          </label>
          <label>
            <span>รหัสโซน <b class="req-star">*</b></span>
            <AppInput v-model="zoneForm.sectionNumber" placeholder="เช่น E" :error="zoneErrors.sectionNumber" @blur="validateZoneField('sectionNumber')" />
          </label>
          <label>
            <span>จำนวนชั้น <b class="req-star">*</b></span>
            <AppInput v-model="zoneForm.floorCount" type="number" placeholder="1 - 20" :error="zoneErrors.floorCount" @blur="validateZoneField('floorCount')" />
          </label>
          <label>
            <span>หมวดสินค้า <b class="req-star">*</b></span>
            <AppSelect
              v-model="zoneForm.category"
              :options="marketCategoryOptions"
              placeholder="เลือกหมวดสินค้า"
              :error="zoneErrors.category"
              @change="validateZoneField('category')"
            />
          </label>
        </form>
        <template #footer>
          <AppButton variant="secondary" @click="dismissZoneEditor">ยกเลิก</AppButton>
          <AppButton @click="createZone">บันทึกโซน</AppButton>
        </template>
      </AppModal>

      <AppModal v-model="showFacilityEditor" title="เพิ่มจุดอำนวยความสะดวก" description="ตำแหน่งถูกกำหนดจากจุดที่เลือกบนผังแล้ว">
        <form class="facility-form" novalidate @submit.prevent="createFacility">
          <label>
            <span>ชื่อจุดบริการ <b class="req-star">*</b></span>
            <AppInput v-model="facilityForm.name" placeholder="เช่น ห้องน้ำอาคาร A" :error="facilityNameError" @blur="validateFacilityName" />
          </label>
          <label>
            <span>ประเภท <b class="req-star">*</b></span>
            <AppSelect v-model="facilityForm.type" :options="facilityTypeOptions" />
          </label>
        </form>
        <template #footer>
          <AppButton variant="secondary" @click="closeFacilityEditor">ยกเลิก</AppButton>
          <AppButton @click="createFacility">เพิ่มจุดบริการ</AppButton>
        </template>
      </AppModal>

      <!-- Modal: เพิ่มผู้เช่าช่วง -->
      <AppModal
        v-model="showAddSubtenantModal"
        title="เพิ่มผู้เช่าช่วง"
        :description="`บันทึกข้อมูลผู้เช่าช่วงสำหรับแผงค้า ${stallEdit.ident || '-'}`"
      >
        <form class="issd-modal-form" novalidate @submit.prevent="saveSubtenant">
          <div class="issd-photo-upload-wrap full">
            <span class="issd-upload-label">รูปถ่ายผู้เช่าช่วง / หน้าร้าน</span>
            <div class="issd-photo-upload-box">
              <div v-if="subtenantForm.photoUrl" class="issd-preview-thumb">
                <img :src="subtenantForm.photoUrl" alt="ภาพตัวอย่าง" />
                <button
                  type="button"
                  class="issd-remove-photo-btn"
                  title="ลบรูป"
                  aria-label="ลบรูป"
                  @click="removeSubtenantPhoto"
                >
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              </div>
              <div v-else class="issd-empty-thumb">
                <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"></path>
                  <circle cx="12" cy="13" r="4"></circle>
                </svg>
              </div>
              <div class="issd-upload-actions">
                <div class="issd-upload-btn-row">
                  <label class="issd-file-choose-btn">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                      <polyline points="17 8 12 3 7 8"></polyline>
                      <line x1="12" y1="3" x2="12" y2="15"></line>
                    </svg>
                    <span>{{ subtenantForm.photoUrl ? 'เปลี่ยนรูปภาพ' : 'เลือกรูปถ่าย' }}</span>
                    <input
                      type="file"
                      accept="image/*"
                      class="issd-file-hidden-input"
                      @change="handleSubtenantImageUpload"
                    />
                  </label>
                  <button
                    v-if="subtenantForm.photoUrl"
                    type="button"
                    class="issd-remove-text-btn"
                    @click="removeSubtenantPhoto"
                  >
                    ลบรูป
                  </button>
                </div>
                <span class="issd-upload-hint">รองรับไฟล์ JPG, PNG, WebP ขนาดไม่เกิน 2MB</span>
              </div>
            </div>
          </div>
          <label>
            <span>ชื่อ-นามสกุล ผู้เช่าช่วง <b class="req-star">*</b></span>
            <AppInput v-model="subtenantForm.name" placeholder="เช่น คุณวราภรณ์ สดใส" :error="subtenantErrors.name" />
          </label>
          <label>
            <span>เบอร์โทรศัพท์ติดต่อ <b class="req-star">*</b></span>
            <AppInput v-model="subtenantForm.phone" placeholder="เช่น 089-123-4567" :error="subtenantErrors.phone" />
          </label>
          <label>
            <span>เลขบัตรประชาชน</span>
            <AppInput v-model="subtenantForm.idCard" placeholder="x-xxxx-xxxxx-xx-x" />
          </label>
          <label>
            <span>ประเภทสินค้า / กิจการ <b class="req-star">*</b></span>
            <AppInput v-model="subtenantForm.businessType" placeholder="เช่น ร้านเครื่องดื่ม, เสื้อผ้าแฟชั่น" :error="subtenantErrors.businessType" />
          </label>
          <label>
            <span>วันที่เริ่มสัญญาเช่าช่วง <b class="req-star">*</b></span>
            <AppInput v-model="subtenantForm.startDate" type="date" :error="subtenantErrors.startDate" />
          </label>
          <label>
            <span>วันที่สิ้นสุดสัญญา</span>
            <AppInput v-model="subtenantForm.endDate" type="date" />
          </label>
          <label>
            <span>ค่าเช่าช่วงต่อเดือน (บาท)</span>
            <AppInput v-model="subtenantForm.subRentAmount" type="number" placeholder="เช่น 2200" />
          </label>
          <label class="full">
            <span>หมายเหตุ / เงื่อนไขเพิ่มเติม</span>
            <textarea
              v-model="subtenantForm.note"
              rows="2"
              class="issd-textarea"
              placeholder="เช่น ยินยอมให้เช่าช่วงเฉพาะวันเสาร์-อาทิตย์ ได้รับความยินยอมจากผู้เช่าหลักแล้ว"
            ></textarea>
          </label>
        </form>
        <template #footer>
          <AppButton variant="secondary" @click="showAddSubtenantModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" @click="saveSubtenant">บันทึกผู้เช่าช่วง</AppButton>
        </template>
      </AppModal>

      <!-- Modal: แจ้งซ่อม / ขออนุญาต -->
      <AppModal
        v-model="showAddMaintenanceModal"
        title="แจ้งซ่อม / ขออนุญาตดำเนินการ"
        :description="`บันทึกรายการสำหรับแผงค้า ${stallEdit.ident || '-'}`"
      >
        <form class="issd-modal-form" novalidate @submit.prevent="saveMaintenance">
          <label>
            <span>ประเภทรายการ <b class="req-star">*</b></span>
            <AppSelect v-model="maintenanceForm.type" :options="maintenanceTypeOptions" />
          </label>
          <label>
            <span>สถานะการดำเนินการ</span>
            <AppSelect v-model="maintenanceForm.status" :options="maintenanceStatusOptions" />
          </label>
          <label class="full">
            <span>หัวข้อเรื่อง / อาการที่ต้องการแก้ไข <b class="req-star">*</b></span>
            <AppInput v-model="maintenanceForm.title" placeholder="เช่น ขอซ่อมแซมกระเบื้องปูพื้น, ติดตั้งมิเตอร์ไฟฟ้า" :error="maintenanceErrors.title" />
          </label>
          <label class="full">
            <span>รายละเอียดคำขอ / ขอบเขตงาน</span>
            <textarea
              v-model="maintenanceForm.description"
              rows="2"
              class="issd-textarea"
              placeholder="ระบุรายละเอียด เช่น จุดที่ชำรุด, อุปกรณ์ที่ขอติดตั้งเพิ่ม"
            ></textarea>
          </label>
          <label>
            <span>ผู้แจ้งเรื่อง <b class="req-star">*</b></span>
            <AppInput v-model="maintenanceForm.requesterName" placeholder="เช่น ผู้เช่าแผง, เจ้าหน้าที่ตรวจผัง" :error="maintenanceErrors.requesterName" />
          </label>
          <label>
            <span>วันที่แจ้งเรื่อง <b class="req-star">*</b></span>
            <AppInput v-model="maintenanceForm.requestDate" type="date" :error="maintenanceErrors.requestDate" />
          </label>
          <label>
            <span>ช่าง / ผู้รับผิดชอบ</span>
            <AppInput v-model="maintenanceForm.technicianName" placeholder="เช่น ช่างสมพงษ์ (งานอาคาร)" />
          </label>
          <label>
            <span>วันที่นัดหมายเข้าดำเนินงาน</span>
            <AppInput v-model="maintenanceForm.scheduledDate" type="date" />
          </label>
          <label>
            <span>ประมาณการค่าใช้จ่าย (บาท)</span>
            <AppInput v-model="maintenanceForm.estimatedCost" type="number" placeholder="เช่น 500" />
          </label>
        </form>
        <template #footer>
          <AppButton variant="secondary" @click="showAddMaintenanceModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" @click="saveMaintenance">บันทึกรายการ</AppButton>
        </template>
      </AppModal>

      <!-- Modal: บันทึกการชำระเงิน -->
      <AppModal
        v-model="showAddPaymentModal"
        title="บันทึกการชำระเงิน"
        :description="`บันทึกประวัติรับชำระสำหรับแผงค้า ${stallEdit.ident || '-'}`"
      >
        <form class="issd-modal-form" novalidate @submit.prevent="savePayment">
          <label>
            <span>งวดเดือน / รอบบิล <b class="req-star">*</b></span>
            <AppInput v-model="paymentForm.periodMonth" placeholder="เช่น งวดเมษายน 2569" :error="paymentErrors.periodMonth" />
          </label>
          <label>
            <span>รายการ / วัตถุประสงค์</span>
            <AppInput v-model="paymentForm.title" placeholder="เช่น ค่าเช่าแผงรายเดือน, ค่าไฟส่วนกลาง" />
          </label>
          <label>
            <span>จำนวนเงิน (บาท) <b class="req-star">*</b></span>
            <AppInput v-model="paymentForm.amount" type="number" placeholder="เช่น 1800" :error="paymentErrors.amount" />
          </label>
          <label>
            <span>วันที่ชำระเงิน <b class="req-star">*</b></span>
            <AppInput v-model="paymentForm.paidDate" type="date" :error="paymentErrors.paidDate" />
          </label>
          <label>
            <span>ช่องทางการชำระ</span>
            <AppSelect v-model="paymentForm.paymentMethod" :options="paymentMethodOptions" />
          </label>
          <label>
            <span>เลขที่ใบเสร็จ / อ้างอิงสลิป</span>
            <AppInput v-model="paymentForm.receiptNumber" placeholder="เช่น REC-6904-01" />
          </label>
        </form>
        <template #footer>
          <AppButton variant="secondary" @click="showAddPaymentModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" @click="savePayment">บันทึกการชำระเงิน</AppButton>
        </template>
      </AppModal>

      <!-- Modal: จัดการสัญญาเช่า (สร้างใหม่ / แก้ไข) -->
      <AppModal
        v-model="showContractModal"
        :title="isEditingContract ? 'แก้ไขสัญญาเช่า' : 'ออกสัญญาเช่าใหม่'"
        :description="`จัดการสัญญาเช่าสำหรับแผงค้า ${stallEdit.ident || '-'}`"
      >
        <form class="issd-modal-form" novalidate @submit.prevent="saveContract">
          <label>
            <span>เลขที่สัญญา <b class="req-star">*</b></span>
            <AppInput v-model="contractForm.contractNumber" placeholder="เช่น CNT-2026-A-001" :error="contractErrors.contractNumber" />
          </label>
          <label>
            <span>สถานะสัญญา</span>
            <AppSelect v-model="contractForm.status" :options="contractStatusOptions" />
          </label>
          <label>
            <span>ชื่อ-นามสกุล ผู้เช่า <b class="req-star">*</b></span>
            <AppInput v-model="contractForm.tenantName" placeholder="เช่น คุณสมศักดิ์ มั่งคั่ง" :error="contractErrors.tenantName" />
          </label>
          <label>
            <span>เบอร์โทรศัพท์ติดต่อ <b class="req-star">*</b></span>
            <AppInput v-model="contractForm.tenantPhone" placeholder="เช่น 081-234-5678" :error="contractErrors.tenantPhone" />
          </label>
          <label>
            <span>เลขบัตรประชาชน</span>
            <AppInput v-model="contractForm.tenantIdCard" placeholder="x-xxxx-xxxxx-xx-x" />
          </label>
          <label>
            <span>ชื่อร้านค้า / ป้ายร้าน</span>
            <AppInput v-model="contractForm.shopName" placeholder="เช่น ส้มตำรสเด็ด, ร้านเครื่องดื่มใจดี" />
          </label>
          <label>
            <span>วันที่เริ่มสัญญา <b class="req-star">*</b></span>
            <AppInput v-model="contractForm.startDate" type="date" :error="contractErrors.startDate" />
          </label>
          <label>
            <span>วันที่สิ้นสุดสัญญา <b class="req-star">*</b></span>
            <AppInput v-model="contractForm.endDate" type="date" :error="contractErrors.endDate" />
          </label>
          <label>
            <span>ค่าเช่าต่อเดือน (บาท) <b class="req-star">*</b></span>
            <AppInput v-model="contractForm.monthlyRent" type="number" placeholder="เช่น 2000" :error="contractErrors.monthlyRent" />
          </label>
          <label>
            <span>เงินประกันความเสียหาย (บาท)</span>
            <AppInput v-model="contractForm.depositAmount" type="number" placeholder="เช่น 4000" />
          </label>
          <label class="full">
            <span>เงื่อนไขและข้อตกลงเพิ่มเติม</span>
            <textarea
              v-model="contractForm.note"
              rows="2"
              class="issd-textarea"
              placeholder="เช่น ชำระค่าเช่าทุกวันที่ 5 ของเดือน, ห้ามประกอบอาหารก่อควันรุนแรง"
            ></textarea>
          </label>
        </form>
        <template #footer>
          <AppButton variant="secondary" @click="showContractModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" @click="saveContract">{{ isEditingContract ? 'บันทึกการแก้ไข' : 'ยืนยันออกสัญญาเช่า' }}</AppButton>
        </template>
      </AppModal>


      <transition name="modal-bounce">
        <div
          v-if="showUnsavedLeaveModal"
          class="unsaved-leave-backdrop"
          role="dialog"
          aria-modal="true"
          @click.self="resolveUnsavedLeave('cancel')"
        >
          <div class="unsaved-leave-card">
            <button
              type="button"
              class="unsaved-leave-close"
              aria-label="อยู่หน้านี้ต่อ"
              title="อยู่หน้านี้ต่อ"
              @click="resolveUnsavedLeave('cancel')"
            >
              <CloseIcon />
            </button>
            <div class="unsaved-leave-icon">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.3" stroke-linecap="round" stroke-linejoin="round">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
            </div>
            <div class="unsaved-leave-content">
              <h3>มีการแก้ไขที่ยังไม่ได้บันทึก</h3>
              <p>คุณมีการแก้ไขข้อมูลผังตลาดที่ยังไม่ได้บันทึก ต้องการบันทึกการแก้ไขก่อนออกจากหน้านี้หรือไม่?</p>
            </div>
            <div class="unsaved-leave-actions">
              <AppButton variant="outline" :disabled="savingGIS" @click="resolveUnsavedLeave('discard')">
                ไม่บันทึก
              </AppButton>
              <AppButton variant="primary" :loading="savingGIS" @click="resolveUnsavedLeave('save')">
                บันทึกการแก้ไข
              </AppButton>
            </div>
          </div>
        </div>
      </transition>

      <transition name="modal-bounce">
        <div
          v-if="showInspectPhotoModal"
          class="unsaved-leave-backdrop inspect-photo-backdrop"
          role="dialog"
          aria-modal="true"
          @click.self="showInspectPhotoModal = false"
        >
          <div class="inspect-photo-modal">
            <div class="ipm-head">
              <div>
                <strong>ภาพตัวอย่างแผง {{ stallEdit.ident }}</strong>
                <span>{{ selectedProject?.name || stallEdit.project_name }} · {{ formatFloorDisplay(stallEdit.floor) }}</span>
              </div>
              <button type="button" class="ipm-close" aria-label="ปิดตัวอย่างภาพ" @click="showInspectPhotoModal = false">
                <CloseIcon />
              </button>
            </div>
            <div class="ipm-stage">
              <img :src="inspectStallPhotos[inspectPhotoIndex]" :alt="`ภาพแผง ${stallEdit.ident}`" class="ipm-main-img" />
              <button type="button" class="ipm-nav is-prev" aria-label="ภาพก่อนหน้า" @click="stepInspectPhoto(-1)">‹</button>
              <button type="button" class="ipm-nav is-next" aria-label="ภาพถัดไป" @click="stepInspectPhoto(1)">›</button>
              <span class="ipm-counter">{{ inspectPhotoIndex + 1 }} / {{ inspectStallPhotos.length }}</span>
            </div>
            <div class="ipm-thumbs">
              <button
                v-for="(src, idx) in inspectStallPhotos"
                :key="idx"
                type="button"
                class="ipm-thumb"
                :class="{ active: inspectPhotoIndex === idx }"
                @click="inspectPhotoIndex = idx"
              >
                <img :src="src" :alt="`รูปที่ ${idx + 1}`" />
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import AppButton from '~/components/AppButton.vue'
import AppInput from '~/components/AppInput.vue'
import AppModal from '~/components/AppModal.vue'
import AppPagination from '~/components/AppPagination.vue'
import AppSegmentedControl from '~/components/AppSegmentedControl.vue'
import AppSelect from '~/components/AppSelect.vue'
import StallUtilityBillingPanel from './components/StallUtilityBillingPanel.vue'
import StallTempBillPanel from './components/StallTempBillPanel.vue'
import StallCostDrawer from './components/StallCostDrawer.vue'
import { useStallBilling, type CostFilter } from '../composables/useStallBilling'
import { getMarketGisConfig } from '../data/markets_registry'
import { loadMarketGisJson, mergeMarketProjects, clearMarketGisCache, cloneProjects } from '../data/market_loader'
import { getAllFloorLayoutsForProject, isOverviewBoundaryProject } from '../data/market_stall_layout'
import type { ProjectZone, MarketFacility } from '../data/market_gis_types'

const route = useRoute()
const router = useRouter()
const { get, post, put, del } = useApi()
const { showSuccess, showError, showConfirm } = useAppAlert()
const { role } = useAuth()
const { applyMarketOverride } = useMarketOverrides()
const isStaff = computed(() => role.value === 'staff')
const mapRef = ref<any>(null)
const marketId = computed(() => String(route.query.market || '1'))
const isMapOverviewMode = computed(() => route.path === '/stalls/map' || route.query.mode === 'map')
const isInspectMode = computed(() => isMapOverviewMode.value || String(route.query.mode || '') === 'inspect')
export type MapPageViewMode = 'stalls' | 'payments' | 'meters'
const pageMode = ref<MapPageViewMode>('stalls')
const isViewDropdownOpen = ref(false)
const viewDropdownWrapRef = ref<HTMLElement | null>(null)
const dashboardPanelOpen = ref(false)
const activeMapSidePanel = ref<'dashboard' | 'search' | 'filter' | 'layers' | null>(null)
const selectedDashboardProjectId = ref('')

const viewModeOptions: {
  value: MapPageViewMode
  label: string
  desc: string
  tone: 'blue' | 'emerald' | 'amber'
  svgIcon: string
}[] = [
  {
    value: 'stalls',
    label: 'ข้อมูลแผง',
    desc: 'สถานะแผง สัญญา ผู้เช่า',
    tone: 'blue',
    svgIcon: '<polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"></polygon><line x1="8" y1="2" x2="8" y2="18"></line><line x1="16" y1="6" x2="16" y2="22"></line>'
  },
  {
    value: 'payments',
    label: 'การชำระเงิน',
    desc: 'ยอดเรียกเก็บและรายการค้างชำระ',
    tone: 'emerald',
    svgIcon: '<rect x="2" y="5" width="20" height="14" rx="2"></rect><line x1="2" y1="10" x2="22" y2="10"></line>'
  },
  {
    value: 'meters',
    label: 'การจดมิเตอร์',
    desc: 'บันทึกและตรวจสอบมิเตอร์น้ำ-ไฟ',
    tone: 'amber',
    svgIcon: '<path d="m12 14 4-4"></path><path d="M3.34 19a10 10 0 1 1 17.32 0"></path>'
  }
]

const currentModeOption = computed(() => {
  return viewModeOptions.find(o => o.value === pageMode.value) || viewModeOptions[0]
})

const selectPageMode = (mode: MapPageViewMode) => {
  pageMode.value = mode
  isViewDropdownOpen.value = false
  selectedPaymentFilter.value = null
  selectedMeterFilter.value = null
  dashboardStatusFilter.value = 'all'
  if (import.meta.client) nextTick(() => setTimeout(() => window.dispatchEvent(new Event('resize')), 100))
}

const handleClickOutsideDropdown = (e: MouseEvent) => {
  if (isViewDropdownOpen.value && viewDropdownWrapRef.value && !viewDropdownWrapRef.value.contains(e.target as Node)) {
    isViewDropdownOpen.value = false
  }
}

watch(pageMode, () => {
  if (import.meta.client) nextTick(() => setTimeout(() => window.dispatchEvent(new Event('resize')), 100))
})
const isInspectDashboardMode = computed(() => !isMapOverviewMode.value && String(route.query.mode || '') === 'inspect')
const showDashboardSection = computed(() =>
  isInspectDashboardMode.value || (isMapOverviewMode.value && dashboardPanelOpen.value)
)
const returnProjectId = computed(() => String(route.query.project || '').trim())
const returnStallCode = computed(() => String(route.query.stall || '').trim())
const resolvedMarketID = ref(0)
const apiMarketName = ref('')
const apiMarketRecord = ref<any>(null)
const marketBookings = ref<any[]>([])
const marketWaitlistItems = ref<any[]>([])
const marketApplications = ref<any[]>([])
const marketContracts = ref<any[]>([])
const dashboardZoneFilter = ref<string>('all')
const dashboardUpdatedDateText = computed(() => {
  try {
    return new Intl.DateTimeFormat('th-TH', { day: 'numeric', month: 'long', year: 'numeric' }).format(new Date())
  } catch {
    return '25 ตุลาคม 2569'
  }
})
const dashboardViewTitle = computed(() => {
  return 'แดชบอร์ด'
})
const dashboardViewSubtitle = computed(() => {
  if (pageMode.value === 'payments') return 'สรุปรอบปัจจุบัน'
  if (pageMode.value === 'meters') return 'สรุปรอบปัจจุบัน'
  return 'สรุปสถานะแผงและสัญญา'
})
const compactCategoryLabel = (name: string) => {
  const normalized = String(name || '').trim()
  return normalized.length > 28 ? `${normalized.slice(0, 28)}…` : normalized
}
const marketName = computed(() => applyMarketOverride({ id: resolvedMarketID.value, slug: marketId.value, name: apiMarketName.value || getMarketGisConfig(marketId.value).name })?.name || getMarketGisConfig(marketId.value).name)
const stalls = ref<any[]>([])
const projects = ref<ProjectZone[]>([])
const facilities = ref<MarketFacility[]>([])
const savedProjectsSnapshot = ref<ProjectZone[]>([])
const savedFacilitiesSnapshot = ref<MarketFacility[]>([])
const savedStallsSnapshot = ref<any[]>([])
const hasUnsavedChanges = ref(false)
const pendingStallUpserts = ref<any[]>([])
const pendingStallDeleteIds = ref<number[]>([])
const dashboardStatusFilter = ref<'all' | 'available' | 'reserved' | 'occupied' | 'action_required' | 'legal_case'>('all')
const setDashboardStatusFilter = (status: typeof dashboardStatusFilter.value) => {
  dashboardStatusFilter.value = dashboardStatusFilter.value === status ? 'all' : status
}
const visibleProjects = computed(() => projects.value.filter(p => !isOverviewBoundaryProject(p, projects.value)))

const resolveRealStallStatus = (codeRaw: string, project?: ProjectZone | null, rawStall?: any): 'available' | 'reserved' | 'occupied' => {
  const code = String(codeRaw || '').trim()
  const upperCode = code.toUpperCase()
  const override = stallStatusOverrides.value[code] || stallStatusOverrides.value[upperCode]
  if (override) return override
  if (project?.layoutMode === 'manual') {
    const manualStatus = String(rawStall?.status || rawStall?.stall_status || '').toLowerCase()
    if (manualStatus === 'reserved' || manualStatus === 'occupied') return manualStatus
  }
  return 'available'
}

const totalMarketStalls = computed(() => visibleProjects.value.reduce((total, project) => {
  let count = 0
  getAllFloorLayoutsForProject(project).forEach((layout) => {
    count += layout.stalls.length
  })
  return total + (count || Number(project.totalStalls) || 0)
}, 0))

const stallStatusSummary = computed(() => {
  const summary = { available: 0, reserved: 0, occupied: 0 }

  visibleProjects.value.forEach((project) => {
    let layoutStallsCount = 0
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        layoutStallsCount += 1
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim()
        const status = resolveRealStallStatus(code, project, stall)
        if (status === 'reserved') summary.reserved += 1
        else if (status === 'occupied') summary.occupied += 1
        else summary.available += 1
      })
    })

    const projectTotal = Number(project.totalStalls) || 0
    if (projectTotal > layoutStallsCount) {
      summary.available += projectTotal - layoutStallsCount
    }
  })

  return summary
})

const openZoneOnMap = (zoneId: string) => {
  const target = projects.value.find(p => p.id === zoneId) || null
  selectedProject.value = target
  setDashboardViewMode('map')
}

const realMarketSummary = computed(() => {
  let total = 0
  let available = 0
  let reserved = 0
  let occupied = 0
  let occupiedMonthlyRent = 0
  let availableMonthlyRent = 0
  let reservedMonthlyRent = 0
  let totalMonthlyRent = 0

  const dbRentByCode: Record<string, number> = {}
  stalls.value.forEach((s: any) => {
    const code = String(s.ident || s.stall_ident || s.stall_number || '').trim().toUpperCase()
    const rent = Number(s.monthly_rent || s.monthlyRent || 0)
    if (code && rent > 0) dbRentByCode[code] = rent
  })
  marketContracts.value.forEach((c: any) => {
    const st = String(c.contract_state || c.status || '').toLowerCase()
    if (st === 'terminated' || st === 'cancelled' || st === 'expired') return
    const code = String(c.stall_ident || c.stall_number || '').trim().toUpperCase()
    const rent = Number(c.monthly_rental || c.monthly_rent || 0)
    if (code && rent > 0) dbRentByCode[code] = rent
  })

  visibleProjects.value.forEach((project) => {
    let layoutStallsCount = 0
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        layoutStallsCount += 1
        const rawCode = String(stall.code || stall.ident || stall.stall_ident || '').trim()
        const upperCode = rawCode.toUpperCase()
        const status = resolveRealStallStatus(rawCode, project, stall)
        const rent = Number(dbRentByCode[upperCode] || stall.monthlyRent || stall.monthly_rent || 1800)

        total += 1
        totalMonthlyRent += rent
        if (status === 'occupied') {
          occupied += 1
          occupiedMonthlyRent += rent
        } else if (status === 'reserved') {
          reserved += 1
          reservedMonthlyRent += rent
        } else {
          available += 1
          availableMonthlyRent += rent
        }
      })
    })

    const projectTotal = Number(project.totalStalls) || 0
    if (projectTotal > layoutStallsCount) {
      const diff = projectTotal - layoutStallsCount
      total += diff
      available += diff
      availableMonthlyRent += diff * 1800
      totalMonthlyRent += diff * 1800
    }
  })

  const safeDenom = Math.max(1, total)
  const occupancyPct = total > 0 ? Math.round((occupied / safeDenom) * 100) : 0
  const reservedPct = total > 0 ? Math.round((reserved / safeDenom) * 100) : 0
  const availablePct = total > 0 ? Math.max(0, 100 - occupancyPct - reservedPct) : 0
  const avgRent = total > 0 ? Math.round(totalMonthlyRent / safeDenom) : 0

  return {
    zonesCount: visibleProjects.value.length,
    total,
    available,
    reserved,
    occupied,
    occupancyPct,
    reservedPct,
    availablePct,
    occupiedMonthlyRent,
    availableMonthlyRent,
    reservedMonthlyRent,
    totalMonthlyRent,
    avgRent
  }
})

const tenantNameOf = (code: string) => {
  const upper = String(code || '').trim().toUpperCase()
  const c = marketContracts.value.find((item: any) =>
    String(item?.stall_number || item?.stall_ident || item?.stall_code || '').split(',').map(x => x.trim().toUpperCase()).includes(upper))
  return c?.tenant_name || 'ผู้เช่าตามสัญญา'
}

const selectedPaymentFilter = ref<'due' | 'paid' | 'pending' | 'overdue' | 'total_overdue' | null>(null)
const togglePaymentFilter = (filter: 'due' | 'paid' | 'pending' | 'overdue' | 'total_overdue') => {
  if (selectedPaymentFilter.value === filter) {
    selectedPaymentFilter.value = null
  } else {
    selectedPaymentFilter.value = filter
  }
}

const selectedMeterFilter = ref<'total' | 'recorded' | 'pending' | 'overdue' | 'anomaly' | null>(null)
const toggleMeterFilter = (filter: 'total' | 'recorded' | 'pending' | 'overdue' | 'anomaly') => {
  if (selectedMeterFilter.value === filter) {
    selectedMeterFilter.value = null
  } else {
    selectedMeterFilter.value = filter
  }
}

const allMarketStallsList = computed(() => {
  const result: {
    code: string
    upperCode: string
    projectId: string
    projectName: string
    status: 'available' | 'reserved' | 'occupied'
    monthlyRent: number
    tenantName: string
  }[] = []

  const seen = new Set<string>()
  visibleProjects.value.forEach((p) => {
    getAllFloorLayoutsForProject(p).forEach((layout) => {
      layout.stalls.forEach((s: any) => {
        const rawCode = String(s.code || s.ident || s.stall_ident || '').trim()
        if (!rawCode) return
        const upper = rawCode.toUpperCase()
        if (seen.has(upper)) return
        seen.add(upper)

        const st = resolveRealStallStatus(rawCode, p, s)
        const rent = Number(s.monthlyRent || s.monthly_rent || 1800)
        const tName = tenantNameOf(rawCode)

        result.push({
          code: rawCode,
          upperCode: upper,
          projectId: p.id,
          projectName: p.name,
          status: st,
          monthlyRent: rent,
          tenantName: tName
        })
      })
    })
  })
  return result
})

const paymentDatasets = computed(() => {
  const stallsList = allMarketStallsList.value
  const pool = [
    ...stallsList.filter(s => s.status === 'occupied'),
    ...stallsList.filter(s => s.status !== 'occupied')
  ]
  const target = pool.slice(0, 42)

  const paidStalls = target.slice(0, 28)
  const pendingStalls = target.slice(28, 37)
  const overdueStalls = target.slice(37, 42)

  return {
    dueStalls: target,
    paidStalls,
    pendingStalls,
    overdueStalls
  }
})

const paymentStallCodeMap = computed(() => ({
  due: paymentDatasets.value.dueStalls.map(s => s.upperCode),
  paid: paymentDatasets.value.paidStalls.map(s => s.upperCode),
  pending: paymentDatasets.value.pendingStalls.map(s => s.upperCode),
  overdue: paymentDatasets.value.overdueStalls.map(s => s.upperCode)
}))

const paymentSummary = computed(() => ({
  dueCount: paymentDatasets.value.dueStalls.length || 42,
  dueAmount: 284500,
  paidCount: paymentDatasets.value.paidStalls.length || 28,
  paidAmount: 192000,
  pendingCount: paymentDatasets.value.pendingStalls.length || 9,
  pendingAmount: 58500,
  overdueCount: paymentDatasets.value.overdueStalls.length || 5,
  overdueAmount: 34000
}))

const meterDatasets = computed(() => {
  const target = paymentDatasets.value.dueStalls
  const recordedStalls = target.slice(0, 35)
  const pendingStalls = target.slice(35, 42)
  const overdueStalls: typeof target = []
  const anomalyStalls = [recordedStalls[2], recordedStalls[11]].filter(Boolean)

  return {
    totalStalls: target,
    recordedStalls,
    pendingStalls,
    overdueStalls,
    anomalyStalls
  }
})

const meterStallCodeMap = computed(() => ({
  total: meterDatasets.value.totalStalls.map(s => s.upperCode),
  recorded: meterDatasets.value.recordedStalls.map(s => s.upperCode),
  pending: meterDatasets.value.pendingStalls.map(s => s.upperCode),
  overdue: meterDatasets.value.overdueStalls.map(s => s.upperCode),
  anomaly: meterDatasets.value.anomalyStalls.map(s => s.upperCode)
}))

const meterSummary = computed(() => ({
  totalCount: meterDatasets.value.totalStalls.length || 42,
  waterCount: meterDatasets.value.totalStalls.length || 42,
  electricCount: meterDatasets.value.totalStalls.length || 42,
  recordedCount: meterDatasets.value.recordedStalls.length || 35,
  completionRate: 83.3,
  pendingCount: meterDatasets.value.pendingStalls.length || 7,
  dueDateText: '30 ต.ค.',
  overdueCount: meterDatasets.value.overdueStalls.length || 0,
  anomalyCount: meterDatasets.value.anomalyStalls.length || 2
}))

const mapActiveFilterCodes = computed<string[] | null>(() => {
  if (pageMode.value === 'stalls') {
    if (dashboardStatusFilter.value === 'legal_case') {
      return legalCaseStallCodes.value
    }
    if (dashboardStatusFilter.value === 'action_required') {
      return actionRequiredStallCodes.value
    }
    return null
  }

  if (pageMode.value === 'payments') {
    if (!selectedPaymentFilter.value) return null
    if (selectedPaymentFilter.value === 'due') return paymentStallCodeMap.value.due
    if (selectedPaymentFilter.value === 'paid') return paymentStallCodeMap.value.paid
    if (selectedPaymentFilter.value === 'pending') return paymentStallCodeMap.value.pending
    if (selectedPaymentFilter.value === 'overdue' || selectedPaymentFilter.value === 'total_overdue') {
      return paymentStallCodeMap.value.overdue
    }
    return null
  }

  if (pageMode.value === 'meters') {
    if (!selectedMeterFilter.value) return null
    if (selectedMeterFilter.value === 'total') return meterStallCodeMap.value.total
    if (selectedMeterFilter.value === 'recorded') return meterStallCodeMap.value.recorded
    if (selectedMeterFilter.value === 'pending') return meterStallCodeMap.value.pending
    if (selectedMeterFilter.value === 'overdue') return meterStallCodeMap.value.overdue
    if (selectedMeterFilter.value === 'anomaly') return meterStallCodeMap.value.anomaly
    return null
  }

  return null
})

const overdueStallTableItems = computed(() => {
  const overdues = paymentDatasets.value.overdueStalls
  const defaultItems = [
    { code: 'A-02', tenant: 'คุณสมชาย พาณิชย์', zone: 'โซน A', amountText: '฿10,200', daysText: 'เกิน 45 วัน' },
    { code: 'A-09', tenant: 'คุณวิภา รุ่งเรือง', zone: 'โซน A', amountText: '฿7,800', daysText: 'เกิน 32 วัน' },
    { code: 'B-03', tenant: 'คุณมานพ ค้าขาย', zone: 'โซน B', amountText: '฿6,500', daysText: 'เกิน 18 วัน' },
    { code: 'B-14', tenant: 'คุณอรทัย สดใส', zone: 'โซน B', amountText: '฿5,400', daysText: 'เกิน 12 วัน' },
    { code: 'C-08', tenant: 'คุณธวัชชัย มีชัย', zone: 'โซน C', amountText: '฿4,100', daysText: 'เกิน 8 วัน' }
  ]
  if (!overdues.length) return defaultItems
  return overdues.slice(0, 5).map((s, idx) => ({
    code: s.code,
    tenant: s.tenantName || defaultItems[idx]?.tenant || 'ผู้เช่า',
    zone: s.projectName || defaultItems[idx]?.zone || 'แผงค้า',
    amountText: defaultItems[idx]?.amountText || '฿6,800',
    daysText: defaultItems[idx]?.daysText || 'เกินกำหนด'
  }))
})

const realActionCounts = computed(() => {
  const isWaitingQueueRow = (b: any) => {
    const st = String(b?.booking_status || b?.status || '').toLowerCase()
    return st === 'waiting_queue' || (!b?.holding_stall && Number(b?.queue_position || b?.queue_number || 0) > 0)
  }
  const realBookingsOnly = marketBookings.value.filter((b: any) => !isWaitingQueueRow(b))
  const pendingBookings = realBookingsOnly.filter((b: any) => {
    const st = String(b.booking_status || b.status || '').toLowerCase()
    return st === 'pending' || st === 'pending_approval' || st === 'pending_slot'
  }).length
  const activeVisits = realBookingsOnly.filter((b: any) => {
    const st = String(b.booking_status || b.status || '').toLowerCase()
    return st === 'approved' || st === 'interested' || st === 'visited'
  }).length
  const realActiveBookings = pendingBookings + activeVisits
  const waitlistCount = marketWaitlistItems.value.filter((w: any) => w.status === 'waiting' || w.status === 'offered').length
  const pendingOfficerApps = marketApplications.value.filter((a: any) => {
    const st = String(a.application_state || a.status || '').toLowerCase()
    return st === 'pending' || st === 'pending_officer' || st === 'submitted'
  }).length
  const pendingDirectorApps = marketApplications.value.filter((a: any) => {
    const st = String(a.application_state || a.status || '').toLowerCase()
    return st === 'pending_director' || st === 'officer_approved'
  }).length
  const approvedAppsWaitingContract = marketApplications.value.filter((a: any) => {
    const st = String(a.application_state || a.status || '').toLowerCase()
    return st === 'approved' || st === 'director_approved'
  }).length

  return {
    pendingBookings,
    activeVisits,
    realActiveBookings,
    waitlistCount,
    pendingOfficerApps,
    pendingDirectorApps,
    approvedAppsWaitingContract,
    totalUrgent: pendingBookings + pendingOfficerApps + approvedAppsWaitingContract
  }
})

interface ExpiringContractItem {
  code: string
  tenant: string
  zone: string
  zoneId?: string
  contractNo: string
  endDateText: string
  daysLeft: number
  daysLeftText: string
  urgency: 'expired' | 'urgent' | 'warning'
  renewalStatus: 'submitted' | 'pending_submit' | 'declined'
  renewalStatusText: string
  monthlyRent: number
}

const dashboardProjectForData = computed(() => visibleProjects.value.find(
  project => project.id === selectedDashboardProjectId.value
) || null)

const dashboardProjectStallCodes = computed(() => {
  const project = dashboardProjectForData.value
  if (!project) return null
  const codes = new Set<string>()
  getAllFloorLayoutsForProject(project).forEach((layout) => {
    layout.stalls.forEach((stall: any) => {
      const code = String(stall.code || stall.ident || stall.stall_ident || '').trim().toUpperCase()
      if (code) codes.add(code)
    })
  })
  return codes
})

const getDashboardRecordCodes = (record: any) => String(
  record?.stall_ident || record?.stall_number || record?.stall_code || record?.ident || record?.code || ''
).split(',').map(code => code.trim().toUpperCase()).filter(Boolean)

const recordMatchesDashboardProject = (record: any) => {
  const codes = dashboardProjectStallCodes.value
  if (!codes) return true
  return getDashboardRecordCodes(record).some(code => codes.has(code))
}

const dashboardContracts = computed(() => marketContracts.value.filter(recordMatchesDashboardProject))

const contractHealthSummary = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const cutoff = new Date(today)
  cutoff.setDate(cutoff.getDate() + 90)
  let normalCount = 0
  let expiringCount = 0
  let expiredCount = 0
  let totalMonthlyRent = 0

  dashboardContracts.value.forEach((contract: any) => {
    const status = String(contract.contract_state || contract.status || '').toLowerCase()
    const endValue = contract.end_date || contract.contract_end_date
    const endDate = endValue ? new Date(endValue) : null
    const hasValidEnd = endDate instanceof Date && !Number.isNaN(endDate.getTime())
    if (['expired', 'terminated', 'cancelled'].includes(status) || (hasValidEnd && endDate < today)) expiredCount += 1
    else if (hasValidEnd && endDate <= cutoff) expiringCount += 1
    else normalCount += 1
    totalMonthlyRent += Number(contract.monthly_rent || contract.rent_amount || contract.rent || 0)
  })

  if (dashboardContracts.value.length === 0) {
    normalCount = selectedDashboardProjectStatus.value?.occupied || 0
  }
  const activeCount = normalCount + expiringCount + expiredCount
  const renewedCount = dashboardContracts.value.filter((contract: any) => ['approved', 'renewed', 'completed'].includes(
    String(contract.renewal_status || contract.renew_status || '').toLowerCase()
  )).length
  const renewalRate = dashboardContracts.value.length ? (renewedCount / dashboardContracts.value.length) * 100 : 0
  return {
    activeCount,
    normalCount,
    expiringCount,
    expiredCount,
    renewalRate,
    totalMonthlyRent
  }
})

const expiringContractList = computed<ExpiringContractItem[]>(() => {
  const occupiedStallList: { code: string; zone: string; zoneId: string }[] = []
  visibleProjects.value.forEach((p) => {
    getAllFloorLayoutsForProject(p).forEach((layout) => {
      layout.stalls.forEach((st: any) => {
        const rawCode = String(st.code || st.ident || st.stall_ident || '').trim()
        if (resolveRealStallStatus(rawCode, p, st) === 'occupied' && rawCode) {
          occupiedStallList.push({ code: rawCode, zone: p.name, zoneId: p.id })
        }
      })
    })
  })

  const sampleData: {
    tenant: string
    codeFallback: string
    zoneFallback: string
    endDateText: string
    daysLeft: number
    daysLeftText: string
    urgency: 'expired' | 'urgent' | 'warning'
    renewalStatus: 'submitted' | 'pending_submit' | 'declined'
    renewalStatusText: string
    monthlyRent: number
  }[] = [
    {
      tenant: 'ร้านสมชาย ผลไม้สด',
      codeFallback: '18-042',
      zoneFallback: 'โครงการที่ 18 (ผักผลไม้)',
      endDateText: '26 ต.ค. 2569',
      daysLeft: -3,
      daysLeftText: 'หมดอายุแล้ว 3 วัน',
      urgency: 'expired',
      renewalStatus: 'pending_submit',
      renewalStatusText: 'ยังไม่ยื่นต่อสัญญา (เสี่ยงคืนแผง)',
      monthlyRent: 4200
    },
    {
      tenant: 'ร้านชาญชัย เครื่องหนัง',
      codeFallback: '26-004',
      zoneFallback: 'โครงการที่ 26 (หัตถกรรม)',
      endDateText: '28 ต.ค. 2569',
      daysLeft: -1,
      daysLeftText: 'หมดอายุแล้ว 1 วัน',
      urgency: 'expired',
      renewalStatus: 'pending_submit',
      renewalStatusText: 'ยังไม่ยื่นต่อสัญญา',
      monthlyRent: 3800
    },
    {
      tenant: 'ร้านรุ่งเรือง ผ้าไหมไทย',
      codeFallback: '24-001',
      zoneFallback: 'โครงการที่ 24 (แฟชั่นเสื้อผ้า)',
      endDateText: '2 พ.ย. 2569',
      daysLeft: 4,
      daysLeftText: 'เหลืออีก 4 วัน',
      urgency: 'urgent',
      renewalStatus: 'submitted',
      renewalStatusText: 'ยื่นคำร้องต่อสัญญาแล้ว (รอตรวจ)',
      monthlyRent: 5500
    },
    {
      tenant: 'ร้านมาลี อาหารตามสั่ง',
      codeFallback: '12-108',
      zoneFallback: 'โครงการที่ 12 (อาหารปรุงสำเร็จ)',
      endDateText: '12 พ.ย. 2569',
      daysLeft: 14,
      daysLeftText: 'เหลืออีก 14 วัน',
      urgency: 'warning',
      renewalStatus: 'declined',
      renewalStatusText: 'แจ้งไม่ต่อสัญญา (รอคืนแผง)',
      monthlyRent: 4800
    },
    {
      tenant: 'ร้านพรพรรณ เครื่องเงินโบราณ',
      codeFallback: '07-015',
      zoneFallback: 'โครงการที่ 7 (ของเก่า/ของสะสม)',
      endDateText: '25 พ.ย. 2569',
      daysLeft: 27,
      daysLeftText: 'เหลืออีก 27 วัน',
      urgency: 'warning',
      renewalStatus: 'submitted',
      renewalStatusText: 'ยื่นคำร้องต่อสัญญาแล้ว',
      monthlyRent: 3500
    }
  ]

  return sampleData.map((item, idx) => {
    const matched = occupiedStallList[idx]
    return {
      code: matched?.code || item.codeFallback,
      tenant: item.tenant,
      zone: matched?.zone || item.zoneFallback,
      zoneId: matched?.zoneId,
      contractNo: `CNT-2568-${String(idx + 104).padStart(4, '0')}`,
      endDateText: item.endDateText,
      daysLeft: item.daysLeft,
      daysLeftText: item.daysLeftText,
      urgency: item.urgency,
      renewalStatus: item.renewalStatus,
      renewalStatusText: item.renewalStatusText,
      monthlyRent: item.monthlyRent
    }
  })
})

const focusExpiringContractStall = async (code: string, zoneId?: string) => {
  const upper = code.trim().toUpperCase()
  activeHighlightedStallCodes.value = [upper]

  let targetProject: ProjectZone | null = null
  if (zoneId) {
    targetProject = projects.value.find(p => p.id === zoneId || p.name === zoneId || p.name?.includes(zoneId)) || null
  }
  if (!targetProject) {
    for (const p of visibleProjects.value) {
      for (const fl of (p.floors || [])) {
        const layout = getAllFloorLayoutsForProject(p).find(l => l.floor === fl.id)
        if (layout?.stalls?.some((s: any) => String(s.code || s.ident || '').toUpperCase() === upper)) {
          targetProject = p
          break
        }
      }
      if (targetProject) break
    }
  }

  if (targetProject && selectedProject.value?.id !== targetProject.id) {
    selectProject(targetProject)
  }

  await nextTick()
  selectStall({ code: upper, ident: upper })
}

const dashZoneCategoryFilter = ref('all')
const dashZoneSearchQuery = ref('')
const dashZonePage = ref(1)
const dashZonePageSize = ref(10)

const zoneRows = computed(() => {
  const dbRentByCode: Record<string, number> = {}
  stalls.value.forEach((s: any) => {
    const code = String(s.ident || s.stall_ident || s.stall_number || '').trim().toUpperCase()
    const rent = Number(s.monthly_rent || s.monthlyRent || 0)
    if (code && rent > 0) dbRentByCode[code] = rent
  })
  marketContracts.value.forEach((c: any) => {
    const st = String(c.contract_state || c.status || '').toLowerCase()
    if (st === 'terminated' || st === 'cancelled' || st === 'expired') return
    const code = String(c.stall_ident || c.stall_number || '').trim().toUpperCase()
    const rent = Number(c.monthly_rental || c.monthly_rent || 0)
    if (code && rent > 0) dbRentByCode[code] = rent
  })

  return visibleProjects.value.map((project, idx) => {
    let layoutStallsCount = 0
    let avail = 0
    let res = 0
    let occ = 0
    let rentSum = 0

    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        layoutStallsCount += 1
        const rawCode = String(stall.code || stall.ident || stall.stall_ident || '').trim()
        const upperCode = rawCode.toUpperCase()
        const status = resolveRealStallStatus(rawCode, project, stall)
        const rent = Number(dbRentByCode[upperCode] || stall.monthlyRent || stall.monthly_rent || 1800)
        rentSum += rent

        if (status === 'occupied') occ += 1
        else if (status === 'reserved') res += 1
        else avail += 1
      })
    })

    const projectTotal = Math.max(layoutStallsCount, Number(project.totalStalls) || 0)
    if (projectTotal > layoutStallsCount) {
      const diff = projectTotal - layoutStallsCount
      avail += diff
      rentSum += diff * 1800
    }

    const occupancyPct = projectTotal > 0 ? Math.round((occ / projectTotal) * 100) : 0
    const avgMonthlyRent = projectTotal > 0 ? Math.round(rentSum / projectTotal) : 1800
    const floorsCount = Array.isArray(project.floors) && project.floors.length > 0 ? project.floors.length : 1

    return {
      id: project.id,
      name: cleanProjectName(project.name) || `โซนที่ ${idx + 1}`,
      sectionNumber: project.sectionNumber || String(idx + 1).padStart(2, '0'),
      category: String(project.category || '').trim() || 'ทั่วไป',
      color: project.color || '#2563eb',
      floorsCount,
      total: projectTotal,
      available: avail,
      reserved: res,
      occupied: occ,
      occupancyPct,
      avgMonthlyRent
    }
  })
})

const donutCircumference = 2 * Math.PI * 54
const donutSegments = computed(() => {
  const total = Math.max(1, realMarketSummary.value.total)
  const items = [
    { key: 'available', count: realMarketSummary.value.available, color: '#10b981' },
    { key: 'reserved', count: realMarketSummary.value.reserved, color: '#f59e0b' },
    { key: 'occupied', count: realMarketSummary.value.occupied, color: '#475569' }
  ]
  let acc = 0
  return items
    .filter(item => item.count > 0)
    .map((item) => {
      const dash = (item.count / total) * donutCircumference
      const seg = { ...item, dash, offset: acc }
      acc += dash
      return seg
    })
})

const categoryChartRows = computed(() => {
  const map = new Map<string, { name: string; total: number; color: string }>()
  zoneRows.value.forEach((z) => {
    const key = z.category || 'ทั่วไป'
    const existing = map.get(key)
    if (existing) {
      existing.total += z.total
    } else {
      map.set(key, { name: key, total: z.total, color: z.color || '#2563eb' })
    }
  })
  const sorted = Array.from(map.values()).sort((a, b) => b.total - a.total).slice(0, 5)
  const maxTotal = Math.max(1, sorted[0]?.total || 1)
  const marketTotal = Math.max(1, realMarketSummary.value.total)
  return sorted.map(item => ({
    ...item,
    pct: Math.round((item.total / marketTotal) * 100),
    barPct: Math.max(6, Math.round((item.total / maxTotal) * 100))
  }))
})

const dashCategoryFilterOptions = computed(() => {
  const uniqueCats = Array.from(new Set(zoneRows.value.map(z => z.category).filter(Boolean)))
  return [
    { value: 'all', label: 'ทุกหมวดสินค้า' },
    ...uniqueCats.map(cat => ({ value: cat, label: cat }))
  ]
})

const filteredZoneRows = computed(() => {
  const q = dashZoneSearchQuery.value.trim().toLowerCase()
  return zoneRows.value.filter((row) => {
    if (dashZoneCategoryFilter.value !== 'all' && row.category !== dashZoneCategoryFilter.value) {
      return false
    }
    if (!q) return true
    return (
      row.name.toLowerCase().includes(q) ||
      row.sectionNumber.toLowerCase().includes(q) ||
      row.category.toLowerCase().includes(q)
    )
  })
})

watch([dashZoneCategoryFilter, dashZoneSearchQuery], () => {
  dashZonePage.value = 1
})

const paginatedZoneRows = computed(() => {
  const start = (dashZonePage.value - 1) * dashZonePageSize.value
  return filteredZoneRows.value.slice(start, start + dashZonePageSize.value)
})

const dashActivityTab = ref<'all' | 'booking' | 'application'>('all')

const mapBookingStatusMeta = (rawStatus: string) => {
  const s = String(rawStatus || '').toLowerCase()
  if (s === 'pending' || s === 'pending_approval') return { label: 'รออนุมัติวันนัด', badgeClass: 'badge-amber' }
  if (s === 'approved') return { label: 'อนุมัติวันนัดแล้ว', badgeClass: 'badge-green' }
  if (s === 'interested' || s === 'visited') return { label: 'เข้าชมแล้ว / สนใจเช่า', badgeClass: 'badge-green' }
  if (s === 'waiting_queue') return { label: 'รอคิวสำรอง', badgeClass: 'badge-amber' }
  if (s === 'rejected' || s === 'expired' || s === 'no_show') return { label: 'ปฏิเสธ / หมดเวลา', badgeClass: 'badge-red' }
  if (s === 'cancelled' || s === 'canceled') return { label: 'ยกเลิกแล้ว', badgeClass: 'badge-gray' }
  return { label: rawStatus || 'รอดำเนินการ', badgeClass: 'badge-amber' }
}

const mapAppStatusMeta = (rawStatus: string) => {
  const s = String(rawStatus || '').toLowerCase()
  if (s === 'pending' || s === 'pending_officer' || s === 'submitted') return { label: 'รอ จนท. ตรวจสอบ', badgeClass: 'badge-amber' }
  if (s === 'pending_director' || s === 'officer_approved') return { label: 'รอผู้บริหารอนุมัติ', badgeClass: 'badge-amber' }
  if (s === 'approved' || s === 'director_approved' || s === 'completed') return { label: 'อนุมัติแล้ว', badgeClass: 'badge-green' }
  if (s === 'rejected' || s === 'revision_requested') return { label: 'แก้ไข / ปฏิเสธ', badgeClass: 'badge-red' }
  if (s === 'cancelled' || s === 'canceled') return { label: 'ยกเลิกแล้ว', badgeClass: 'badge-gray' }
  return { label: rawStatus || 'รอตรวจสอบ', badgeClass: 'badge-amber' }
}

const realRecentActivities = computed(() => {
  const bookingItems = marketBookings.value.map((b: any) => {
    const meta = mapBookingStatusMeta(b.status)
    const visitText = b.visit_date
      ? `${b.visit_date}${b.visit_time_slot ? ' (' + b.visit_time_slot + ')' : ''}`
      : (b.created_at ? String(b.created_at).slice(0, 10) : '-')
    return {
      key: `booking-${b.id || b.booking_number}`,
      kind: 'booking' as const,
      refNumber: b.booking_number || `#BK-${b.id}`,
      typeLabel: 'การจองนัดดูพื้นที่',
      stallCode: b.stall_number || b.stall_ident || b.stall_code || '-',
      personName: b.tenant_name || b.applicant_name || 'ผู้ขอจองพื้นที่',
      phone: b.tenant_phone || b.phone || '',
      dateText: visitText,
      statusLabel: meta.label,
      badgeClass: meta.badgeClass,
      actionUrl: '/bookings',
      sortTime: new Date(b.created_at || b.visit_date || 0).getTime() || Number(b.id || 0)
    }
  })

  const appItems = marketApplications.value.map((a: any) => {
    const meta = mapAppStatusMeta(a.status)
    const dateStr = a.submitted_at || a.created_at
      ? String(a.submitted_at || a.created_at).slice(0, 10)
      : '-'
    return {
      key: `app-${a.id || a.application_number}`,
      kind: 'application' as const,
      refNumber: a.application_number || `#APP-${a.id}`,
      typeLabel: 'คำร้องขอเช่าแผง',
      stallCode: a.stall_number || a.stall_ident || a.stall_code || '-',
      personName: a.applicant_name || a.tenant_name || 'ผู้ยื่นคำร้องขอเช่า',
      phone: a.phone || a.applicant_phone || '',
      dateText: dateStr,
      statusLabel: meta.label,
      badgeClass: meta.badgeClass,
      actionUrl: '/applications',
      sortTime: new Date(a.created_at || a.submitted_at || 0).getTime() || Number(a.id || 0)
    }
  })

  return [...bookingItems, ...appItems].sort((x, y) => y.sortTime - x.sortTime)
})

const filteredRecentActivities = computed(() => {
  if (dashActivityTab.value === 'all') return realRecentActivities.value
  return realRecentActivities.value.filter(item => item.kind === dashActivityTab.value)
})
const showUnsavedLeaveModal = ref(false)
let unsavedLeaveResolver: ((choice: 'save' | 'discard' | 'cancel') => void) | null = null
const selectedProject = ref<ProjectZone | null>(null)
const contracts = ref<any[]>([])
const stallHistory = ref<any[]>([])
const gisPublished = ref(false)
const isManageDataLoading = ref(true)
const savingGIS = ref(false)
const savingStall = ref(false)
const showStallEditor = ref(false)
const showInspectStallOverlay = ref(false)
const showInspectPhotoModal = ref(false)
const emptyStall = () => ({ id: 0, ident: '', name: '', description: '', category: '', zone_id: 1, zone_name: 'โซน A', project_name: selectedProject.value?.name || '', floor: '1', width: 2, depth: 2, daily_rent: 150, monthly_rent: 1800, stall_status: 'available', coordinates_geojson: '', published: false })
const stallEdit = ref(emptyStall())
const inspectPhotoIndex = ref(0)
const DEFAULT_INSPECT_STALL_PHOTOS = [
  '/images/stall_3d_perspective.jpg',
  '/images/market_stall_sample.jpg',
  '/images/market_hall_front.jpg'
]
const inspectStallPhotos = computed(() => {
  const customPhotos = (stallEdit.value as any)?.photos
  if (Array.isArray(customPhotos) && customPhotos.length > 0) return customPhotos
  return DEFAULT_INSPECT_STALL_PHOTOS
})
const inspectStallSizeSqm = computed(() => {
  const w = Number(stallEdit.value.width) || 2
  const d = Number(stallEdit.value.depth) || 2.5
  const area = Number((w * d).toFixed(1))
  return Number.isInteger(area) ? String(area) : area.toFixed(1)
})
const inspectStallEntryFee = computed(() => {
  const explicit = Number((stallEdit.value as any)?.entry_fee || (stallEdit.value as any)?.entryFee || 0)
  if (explicit > 0) return explicit
  const rent = Number(stallEdit.value.monthly_rent) || 1800
  return rent * 20
})
const formatFloorDisplay = (floorVal?: string | number) => {
  const raw = String(floorVal || '1').trim()
  if (!raw) return 'ชั้น 1'
  return raw.startsWith('ชั้น') ? raw : `ชั้น ${raw}`
}
const cleanProjectName = (rawVal?: string) => {
  return String(rawVal || '').replace(/\s*\([^)]*\)/g, '').trim() || 'โซนแผงค้า'
}
// ===== รูปภาพแผง (เพิ่ม/ลบในโหมดแก้ไข GIS) =====
const MAX_STALL_PHOTOS = 6
const getStallPhotosKey = (code: string) => `stall_photos_${marketId.value}_${String(code).toUpperCase()}`
const editablePhotos = computed<string[]>(() => {
  const list = (stallEdit.value as any)?.photos
  return Array.isArray(list) ? list : []
})
const persistStallPhotos = () => {
  const code = String(stallEdit.value.ident || '').trim()
  if (!import.meta.client || !code) return
  try {
    localStorage.setItem(getStallPhotosKey(code), JSON.stringify(editablePhotos.value))
  } catch {
    void showError('พื้นที่จัดเก็บรูปภาพเต็ม กรุณาลบรูปเดิมหรือใช้รูปที่เล็กลง', 'บันทึกรูปไม่สำเร็จ')
  }
}
const loadStallPhotos = () => {
  const code = String(stallEdit.value.ident || '').trim()
  if (!import.meta.client || !code) return
  try {
    const raw = localStorage.getItem(getStallPhotosKey(code))
    if (!raw) return
    const parsed = JSON.parse(raw)
    if (Array.isArray(parsed) && parsed.length) (stallEdit.value as any).photos = parsed
  } catch {}
}
// โหลดรูปที่ผู้ใช้อัปโหลดไว้ทุกครั้งที่เลือกแผงใหม่ (ต้องอยู่หลังประกาศ stallEdit และ loadStallPhotos)
watch(stallEdit, () => loadStallPhotos())
const readFileAsDataUrl = (file: File) => new Promise<string>((resolve, reject) => {
  const reader = new FileReader()
  reader.onload = () => resolve(String(reader.result || ''))
  reader.onerror = () => reject(reader.error)
  reader.readAsDataURL(file)
})
const handleStallPhotoUpload = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files || [])
  input.value = ''
  if (!files.length) return
  const next = [...editablePhotos.value]
  for (const file of files) {
    if (next.length >= MAX_STALL_PHOTOS) {
      void showError(`เพิ่มรูปได้สูงสุด ${MAX_STALL_PHOTOS} รูปต่อแผง`, 'เพิ่มรูปครบแล้ว')
      break
    }
    if (!file.type.startsWith('image/')) {
      void showError('กรุณาเลือกไฟล์รูปภาพเท่านั้น (JPG, PNG, WebP)', 'ไฟล์ไม่ถูกต้อง')
      continue
    }
    if (file.size > 2 * 1024 * 1024) {
      void showError(`รูป "${file.name}" ใหญ่เกิน 2 MB`, 'ไฟล์มีขนาดใหญ่เกินไป')
      continue
    }
    try { next.push(await readFileAsDataUrl(file)) } catch {}
  }
  ;(stallEdit.value as any).photos = next
  persistStallPhotos()
}
const removeStallPhoto = (idx: number) => {
  const next = editablePhotos.value.filter((_, i) => i !== idx)
  ;(stallEdit.value as any).photos = next
  persistStallPhotos()
}
const openInspectPhotoPreview = (idx = 0) => {
  inspectPhotoIndex.value = idx
  showInspectPhotoModal.value = true
}
const stepInspectPhoto = (delta: number) => {
  const len = inspectStallPhotos.value.length
  if (len < 2) return
  inspectPhotoIndex.value = (inspectPhotoIndex.value + delta + len) % len
}
const inspectStallOccupancy = computed(() => {
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  const status = String(stallEdit.value.stall_status || 'available').toLowerCase()

  const matchesStall = (item: any) => {
    const raw = String(item?.stall_number || item?.stall_ident || item?.stall_code || '')
    return raw.split(',').map(s => s.trim().toUpperCase()).includes(code)
  }

  const activeContract = contracts.value[0] || marketContracts.value.find(matchesStall)
  const matchingApp = marketApplications.value.find(matchesStall)
  const matchingBooking = marketBookings.value.find(matchesStall)

  if (status === 'occupied' || activeContract) {
    const tenantName = activeContract?.tenant_name || matchingApp?.applicant_name || 'มีผู้เช่าตามสัญญาปัจจุบัน'
    const refCode = activeContract?.contract_number || activeContract?.contract_no || ''
    const dateRange = (activeContract?.start_date || activeContract?.end_date)
      ? `ระยะสัญญา ${activeContract?.start_date || '-'} ถึง ${activeContract?.end_date || '-'}`
      : 'อยู่ระหว่างอายุสัญญาเช่า'
    return {
      heading: 'ผู้เช่าปัจจุบัน',
      refCode,
      primaryText: tenantName,
      subText: dateRange
    }
  }

  if (status === 'reserved' || matchingApp || matchingBooking) {
    if (matchingApp) {
      return {
        heading: 'คำร้องขอเช่าแผง',
        refCode: matchingApp.application_number || `#${matchingApp.id}`,
        primaryText: matchingApp.applicant_name || 'ผู้ยื่นคำร้องขอเช่า',
        subText: 'อยู่ระหว่างพิจารณาอนุมัติคำร้องขอเช่า'
      }
    }
    if (matchingBooking) {
      return {
        heading: 'การจองนัดเข้าชมพื้นที่',
        refCode: matchingBooking.booking_number || `#${matchingBooking.id}`,
        primaryText: matchingBooking.tenant_name || matchingBooking.applicant_name || 'ผู้จองเข้าชมพื้นที่',
        subText: matchingBooking.visit_date
          ? `วันนัดเข้าชม ${matchingBooking.visit_date}${matchingBooking.visit_time_slot ? ' · ' + matchingBooking.visit_time_slot : ''}`
          : 'อยู่ระหว่างดำเนินการนัดเข้าชมพื้นที่'
      }
    }
    return {
      heading: 'สถานะการจอง',
      refCode: '',
      primaryText: 'ติดจอง / รอดำเนินการทำสัญญา',
      subText: 'กดดูรายละเอียดเพื่อตรวจสอบรายการจองและคำร้อง'
    }
  }

  return {
    heading: 'สถานะการใช้งานแผง',
    refCode: '',
    primaryText: 'แผงว่าง พร้อมเปิดรับการจองและทำสัญญา',
    subText: stallEdit.value.description ? String(stallEdit.value.description) : ''
  }
})
const activeHighlightedStallCodes = computed(() => {
  if (isMapOverviewMode.value && pageMode.value === 'costs') {
    if (costDrawerOpen.value && costStallCode.value) return [costStallCode.value]
    return costFilterCodes.value
  }
  if (showInspectStallOverlay.value && stallEdit.value.ident) {
    return [String(stallEdit.value.ident)]
  }
  return returnStallCode.value ? [returnStallCode.value] : []
})
export interface StallSubtenant {
  id: string
  stallIdent: string
  name: string
  phone: string
  idCard?: string
  businessType: string
  startDate: string
  endDate: string
  subRentAmount?: number
  note?: string
  photoUrl?: string
  status: 'active' | 'ended'
  createdAt: string
  contractNumber?: string
  sequence?: number
}

export interface StallMaintenanceItem {
  id: string
  stallIdent: string
  type: 'repair' | 'renovation' | 'utility' | 'other'
  title: string
  description: string
  requesterName: string
  technicianName?: string
  requestDate: string
  scheduledDate?: string
  estimatedCost?: number
  status: 'pending' | 'in_progress' | 'completed' | 'rejected'
  createdAt: string
  source?: 'tenant' | 'staff'
}

export interface StallPaymentItem {
  id: string
  stallIdent: string
  invoiceNumber: string
  periodMonth: string
  title: string
  amount: number
  paidDate?: string
  paymentMethod: 'qr' | 'transfer' | 'cash' | 'credit'
  receiptNumber?: string
  status: 'paid' | 'pending' | 'overdue'
  createdAt: string
  contractNumber?: string
}

const stallDrawerMode = ref<'detail' | 'edit' | 'full-detail'>('detail')
const stallDetailTab = ref<'info' | 'contracts' | 'subtenants' | 'payments' | 'maintenance' | 'bookings' | 'history'>('info')
const isActionAlertCollapsed = ref(false)
const isMiniActionAlertExpanded = ref(false)
const contractSubTab = ref<'active' | 'history'>('active')
const bookingTypeFilter = ref<'all' | 'booking' | 'application'>('all')
const showZoneEditor = ref(false)
const showFacilityEditor = ref(false)
const facilityNameError = ref('')
const mapVersion = ref(0)
const zoneOptions = [{ value: 1, label: 'โซน A' }, { value: 2, label: 'โซน B' }, { value: 3, label: 'โซน C' }, { value: 4, label: 'โซน D' }]
const floorOptions = computed(() => {
  const fls = selectedProject.value?.floors
  if (fls && fls.length > 0) return fls.map(f => ({ value: String(f.id), label: f.name || `ชั้น ${f.id}` }))
  return [{ value: '1', label: 'ชั้น 1' }, { value: '2', label: 'ชั้น 2' }, { value: '3', label: 'ชั้น 3' }]
})
const stallStatusOptions = [{ value: 'available', label: 'ว่าง' }, { value: 'reserved', label: 'ติดจอง' }, { value: 'occupied', label: 'มีผู้เช่า' }]
const stallStatusLabel = computed(() => stallStatusOptions.find(option => option.value === stallEdit.value.stall_status)?.label || 'ไม่ระบุสถานะ')
const currentStallContract = computed(() => {
  if (stallEdit.value.stall_status === 'available') return null
  return stallContractsList.value.find((contract: any) => {
    const st = String(contract.status || contract.contract_state || '').toLowerCase()
    return st === 'active' || st === 'ใช้งานอยู่'
  }) || null
})
const overduePaymentCount = computed(() => {
  if (stallEdit.value.stall_status === 'available') return 0
  return stallPayments.value.filter(payment => payment.status === 'overdue').length
})
const overduePaymentSum = computed(() => {
  if (stallEdit.value.stall_status === 'available') return 0
  return stallPayments.value
    .filter(item => item.status === 'overdue')
    .reduce((sum, item) => sum + Number(item.amount || 0), 0)
})
const openMaintenanceCount = computed(() => stallMaintenance.value.filter(item => !['completed', 'rejected'].includes(item.status)).length)
const paymentTotal = computed(() => stallPayments.value.reduce((sum, item) => sum + Number(item.amount || 0), 0))
const paidPaymentTotal = computed(() => stallPayments.value.filter(item => item.status === 'paid').reduce((sum, item) => sum + Number(item.amount || 0), 0))
const outstandingPaymentTotal = computed(() => Math.max(0, paymentTotal.value - paidPaymentTotal.value))

const currentStallBooking = computed(() => {
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (!code) return null
  const fromMarket = marketBookings.value.find((b: any) => {
    const raw = String(b.stall_number || b.stall_ident || b.stall_code || '')
    const matches = raw.split(',').map(s => s.trim().toUpperCase()).includes(code)
    if (!matches) return false
    const st = String(b.booking_status || b.status || '').toLowerCase()
    return !['rejected', 'cancelled', 'not_interested', 'no_show', 'terminated'].includes(st)
  })
  if (fromMarket) return fromMarket
  if (stallEdit.value.stall_status === 'available') return null
  const bk = stallBookingsAndApps.value.find(item => item.typeKey === 'booking')
  if (bk) {
    return {
      id: bk.id || 118,
      booking_number: bk.refNumber,
      stall_ident: code,
      stall_number: code,
      tenant_name: bk.applicantName,
      applicant_name: bk.applicantName,
      tenant_phone: bk.phone,
      applicant_phone: bk.phone,
      visit_date: '2026-11-01',
      visit_time_slot: '13.00 - 14.00 น.',
      booking_status: bk.statusText.includes('ยื่นคำร้อง') ? 'lease_applied' : 'approved',
      status: 'A'
    }
  }
  return null
})

const currentStallApplication = computed(() => {
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (!code) return null
  const fromMarket = marketApplications.value.find((a: any) => {
    const raw = String(a.stall_number || a.stall_ident || a.stall_code || '')
    const matches = raw.split(',').map(s => s.trim().toUpperCase()).includes(code)
    if (!matches) return false
    const st = String(a.application_state || a.status || '').toLowerCase()
    return !['rejected', 'cancelled'].includes(st)
  })
  if (fromMarket) return fromMarket
  if (stallEdit.value.stall_status === 'available') return null
  const app = stallBookingsAndApps.value.find(item => item.typeKey === 'application')
  if (app) {
    return {
      id: app.id || 33,
      application_number: app.refNumber,
      stall_ident: code,
      stall_number: code,
      applicant_name: app.applicantName,
      applicant_phone: app.phone,
      tenant_name: app.applicantName,
      tenant_phone: app.phone,
      shop_name: stallEdit.value.category || 'ร้านค้าทั่วไป',
      application_state: app.statusText.includes('อนุมัติแล้ว') ? 'approved' : 'pending_officer',
      status: 'A'
    }
  }
  return null
})

const currentStallWaitlistCount = computed(() => {
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (!code) return 0
  return marketWaitlistItems.value.filter((w: any) => {
    const raw = String(w.stall_number || w.stall_ident || '')
    return raw.split(',').map(s => s.trim().toUpperCase()).includes(code) && (w.status === 'waiting' || w.status === 'offered')
  }).length
})

const contractDaysRemaining = computed(() => {
  if (!currentStallContract.value?.end_date) return null
  const end = new Date(currentStallContract.value.end_date).getTime()
  if (isNaN(end)) return null
  return Math.ceil((end - Date.now()) / (1000 * 60 * 60 * 24))
})

const isContractExpiringSoon = computed(() => {
  if (stallEdit.value.stall_status === 'available') return false
  const days = contractDaysRemaining.value
  return days !== null && days >= 0 && days <= 60
})

const formatThaiDateOnly = (dateStr?: string | null) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return String(dateStr)
  return new Intl.DateTimeFormat('th-TH', { day: 'numeric', month: 'short', year: 'numeric', timeZone: 'Asia/Bangkok' }).format(d)
}

const getContractDurationText = (startDateStr?: string | null, endDateStr?: string | null) => {
  if (!startDateStr || !endDateStr) return '-'
  const start = new Date(startDateStr)
  const end = new Date(endDateStr)
  if (isNaN(start.getTime()) || isNaN(end.getTime())) return '-'

  const diffMs = end.getTime() - start.getTime()
  if (diffMs <= 0) return '-'

  const totalDays = Math.round(diffMs / (1000 * 60 * 60 * 24)) + 1
  const totalMonths = Math.round(totalDays / 30.4375)

  if (totalMonths >= 12) {
    const years = Math.floor(totalMonths / 12)
    const remMonths = totalMonths % 12
    if (remMonths === 0) {
      return `${years} ปี (${totalMonths} เดือน)`
    }
    return `${years} ปี ${remMonths} เดือน`
  } else if (totalMonths > 0) {
    return `${totalMonths} เดือน (${totalDays} วัน)`
  }
  return `${totalDays} วัน`
}

const getShortContractDurationText = (startDateStr?: string | null, endDateStr?: string | null) => {
  if (!startDateStr || !endDateStr) return '-'
  const start = new Date(startDateStr)
  const end = new Date(endDateStr)
  if (isNaN(start.getTime()) || isNaN(end.getTime())) return '-'

  const diffMs = end.getTime() - start.getTime()
  if (diffMs <= 0) return '-'

  const totalDays = Math.round(diffMs / (1000 * 60 * 60 * 24)) + 1
  const totalMonths = Math.round(totalDays / 30.4375)

  if (totalMonths >= 12) {
    const years = Math.floor(totalMonths / 12)
    const remMonths = totalMonths % 12
    if (remMonths === 0) {
      return `${years} ปี`
    }
    return `${years} ปี ${remMonths} เดือน`
  } else if (totalMonths > 0) {
    return `${totalMonths} เดือน`
  }
  return `${totalDays} วัน`
}

const getBookingStatusBadgeLabel = (status: string) => {
  switch (String(status || '').toLowerCase()) {
    case 'pending':
    case 'pending_approval': return 'รอเจ้าหน้าที่อนุมัติ'
    case 'pending_slot': return 'รอเลือกวันนัดชม'
    case 'approved': return 'นัดหมายสำเร็จ'
    case 'interested': return 'รอผู้เช่ายื่นคำร้อง'
    case 'lease_applied':
    case 'completed': return 'ยื่นคำร้องสำเร็จ'
    case 'rejected':
    case 'cancelled': return 'ยกเลิก'
    default: return status || 'รอดำเนินการ'
  }
}

const stallActionItems = computed(() => {
  const items: string[] = []
  if (overduePaymentCount.value > 0) {
    const amt = overduePaymentSum.value > 0 ? ` · ฿${overduePaymentSum.value.toLocaleString('th-TH')}` : ''
    items.push(`ติดตามยอดค้างชำระ ${overduePaymentCount.value} งวด${amt}`)
  }
  if (openMaintenanceCount.value > 0) {
    items.push(`ติดตามงานซ่อม/คำร้องที่ยังไม่ปิด ${openMaintenanceCount.value} รายการ`)
  }
  if (currentStallBooking.value) {
    const st = String(currentStallBooking.value.booking_status || currentStallBooking.value.status || '').toLowerCase()
    if (st === 'pending' || st === 'pending_approval') {
      items.push('มีคำขอจองนัดหมาย รอเจ้าหน้าที่ตรวจสอบและอนุมัติ')
    } else if (st === 'approved') {
      items.push(`มีนัดเข้าชมแผง ${currentStallBooking.value.visit_date || ''}`)
    } else if (st === 'interested') {
      items.push('รอผู้เช่ายื่นคำร้องขอเช่าแผง')
    }
  }
  if (currentStallApplication.value) {
    items.push(`มีคำร้องขอเช่า (${currentStallApplication.value.application_number || 'รอตรวจ'}) อยู่ระหว่างพิจารณา`)
  }
  if (isContractExpiringSoon.value && contractDaysRemaining.value !== null) {
    items.push(`สัญญาใกล้หมดอายุ (เหลือ ${contractDaysRemaining.value} วัน)`)
  }
  return items.slice(0, 3)
})

const stallActionRichItems = computed(() => {
  const list: Array<{
    key: string
    category: string
    severity: 'danger' | 'warning' | 'info'
    title: string
    count?: number
    description: string
    actionLabel: string
    targetTab: 'info' | 'contracts' | 'subtenants' | 'payments' | 'maintenance' | 'bookings' | 'history'
  }> = []

  if (overduePaymentCount.value > 0) {
    const amt = overduePaymentSum.value > 0 ? ` (฿${overduePaymentSum.value.toLocaleString('th-TH')})` : ''
    list.push({
      key: 'act-overdue',
      category: 'ค้างชำระ',
      severity: 'danger',
      title: 'เกินกำหนด',
      count: overduePaymentCount.value,
      description: `มียอดค้างชำระเกินกำหนด กรุณาตรวจสอบรอบบิลและติดตามการชำระเงิน`,
      actionLabel: 'ไปที่ชำระเงิน',
      targetTab: 'payments'
    })
  }

  if (openMaintenanceCount.value > 0) {
    const openTitles = stallMaintenance.value
      .filter(item => !['completed', 'rejected'].includes(item.status))
      .map(item => item.title)
      .slice(0, 2)
      .join(', ')
    list.push({
      key: 'act-maintenance',
      category: 'งานซ่อม',
      severity: 'warning',
      title: 'รอดำเนินการ',
      count: openMaintenanceCount.value,
      description: openTitles ? `รายการล่าสุด: ${openTitles}` : 'มีรายการแจ้งซ่อมหรือขออนุญาตที่ยังไม่ได้ปิดงาน',
      actionLabel: 'ไปที่งานซ่อม',
      targetTab: 'maintenance'
    })
  }

  if (currentStallBooking.value) {
    const b = currentStallBooking.value
    const st = String(b.booking_status || b.status || '').toLowerCase()
    const applicant = b.tenant_name || b.applicant_name || 'ผู้ขอจอง'
    const refNo = b.booking_number || `#${b.id}`
    if (st === 'pending' || st === 'pending_approval') {
      list.push({
        key: 'act-booking-pending',
        category: 'การจอง',
        severity: 'warning',
        title: 'รอตรวจสอบนัดชม',
        count: 1,
        description: `ผู้จอง: ${applicant} รอเจ้าหน้าที่ตรวจสอบและอนุมัติวันนัดหมาย`,
        actionLabel: 'ไปที่การจอง',
        targetTab: 'bookings'
      })
    } else if (st === 'approved') {
      list.push({
        key: 'act-booking-approved',
        category: 'การจอง',
        severity: 'info',
        title: 'นัดเข้าชมแผง',
        count: 1,
        description: `ผู้จอง: ${applicant} นัดเข้าชมวันที่ ${b.visit_date || '-'}${b.visit_time_slot ? ' เวลา ' + b.visit_time_slot : ''}`,
        actionLabel: 'ไปที่การจอง',
        targetTab: 'bookings'
      })
    } else if (st === 'interested') {
      list.push({
        key: 'act-booking-interested',
        category: 'การจอง',
        severity: 'info',
        title: 'รอยื่นคำร้อง',
        count: 1,
        description: `${applicant} เข้าชมแผงแล้วและสนใจเช่า อยู่ระหว่างรอยื่นคำร้องขอเช่า`,
        actionLabel: 'ไปที่การจอง',
        targetTab: 'bookings'
      })
    }
  }

  if (currentStallApplication.value) {
    const a = currentStallApplication.value
    const st = String(a.application_state || a.status || '').toLowerCase()
    if (st !== 'approved' && st !== 'completed') {
      list.push({
        key: 'act-application',
        category: 'คำร้อง',
        severity: 'warning',
        title: 'รอพิจารณา',
        count: 1,
        description: `ผู้ยื่นคำร้อง: ${a.applicant_name || '-'} อยู่ระหว่างรอเจ้าหน้าที่พิจารณาอนุมัติ`,
        actionLabel: 'พิจารณาคำร้อง',
        targetTab: 'bookings'
      })
    }
  }

  if (isContractExpiringSoon.value && contractDaysRemaining.value !== null) {
    list.push({
      key: 'act-contract-expiring',
      category: 'สัญญา',
      severity: 'warning',
      title: 'ใกล้ครบกำหนด',
      description: `สัญญาเลขที่ ${currentStallContract.value?.contract_number || '-'} จะสิ้นสุดในวันที่ ${currentStallContract.value?.end_date || '-'}`,
      actionLabel: 'ไปที่สัญญา',
      targetTab: 'contracts'
    })
  }

  return list
})

// Badge ตัวเลขบนแท็บ: แสดงเฉพาะแท็บที่มีงานต้องดำเนินการ
const tabActionCounts = computed(() => ({
  info: 0,
  contracts: isContractExpiringSoon.value ? 1 : 0,
  subtenants: 0,
  payments: overduePaymentCount.value,
  maintenance: openMaintenanceCount.value,
  bookings: stallActionRichItems.value.filter(item => item.targetTab === 'bookings').length,
  history: 0
}))

const stallDetailTabOptions: { value: 'info' | 'contracts' | 'subtenants' | 'payments' | 'maintenance' | 'bookings' | 'history'; label: string }[] = [
  { value: 'info', label: 'ข้อมูลแผง' },
  { value: 'contracts', label: 'สัญญา' },
  { value: 'subtenants', label: 'ผู้เช่าช่วง' },
  { value: 'payments', label: 'ประวัติชำระเงิน' },
  { value: 'maintenance', label: 'ซ่อม/ขออนุญาต' },
  { value: 'bookings', label: 'การจองและคำร้อง' },
  { value: 'history', label: 'ประวัติการดำเนินการ' }
]

const maintenanceSourceFilter = ref<'all' | 'tenant' | 'staff'>('all')
const getMaintenanceSource = (item: StallMaintenanceItem) => {
  if (item.source) return item.source
  const requester = String(item.requesterName || '').toLowerCase()
  return requester.includes('ผู้เช่า') ? 'tenant' : 'staff'
}

const selectedPaymentContract = ref('all')
const paymentContractOptions = computed(() => [
  { value: 'all', label: 'ทุกสัญญา' },
  ...stallContractsList.value.map((contract: any) => ({
    value: contract.contract_number,
    label: `${contract.contract_number}${contract.status === 'active' ? ' · ปัจจุบัน' : ' · ประวัติ'}`
  }))
])
const filteredStallPayments = computed(() => selectedPaymentContract.value === 'all'
  ? stallPayments.value
  : stallPayments.value.filter(item => item.contractNumber === selectedPaymentContract.value))
const filteredStallMaintenance = computed(() => {
  if (maintenanceSourceFilter.value === 'all') return stallMaintenance.value
  return stallMaintenance.value.filter(item => getMaintenanceSource(item) === maintenanceSourceFilter.value)
})

const formatStallDate = (value?: string) => {
  if (!value) return '-'
  try {
    return new Intl.DateTimeFormat('th-TH', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
  } catch {
    return String(value)
  }
}

// Contract Management State & Handlers
const customStallContracts = ref<Record<string, any[]>>({})
const showContractModal = ref(false)
const isEditingContract = ref(false)
const contractForm = ref({
  id: '',
  contractNumber: '',
  tenantName: '',
  tenantPhone: '',
  tenantIdCard: '',
  shopName: '',
  startDate: '',
  endDate: '',
  monthlyRent: '' as number | string,
  depositAmount: '' as number | string,
  status: 'active' as string,
  note: ''
})
const contractErrors = ref({
  contractNumber: '',
  tenantName: '',
  tenantPhone: '',
  startDate: '',
  endDate: '',
  monthlyRent: ''
})

const contractStatusOptions = [
  { value: 'active', label: 'ใช้งานอยู่ (ปกติ)' },
  { value: 'draft', label: 'ร่างสัญญา' },
  { value: 'expired', label: 'หมดสัญญา' },
  { value: 'cancelled', label: 'ยกเลิกแล้ว' }
]

const getContractsStorageKey = (code: string) => `stall_contracts_custom_${marketId.value}_${code.toUpperCase()}`

const loadStallContracts = (code: string) => {
  if (!code) return
  const upper = code.toUpperCase()
  const key = getContractsStorageKey(upper)
  if (import.meta.client) {
    const raw = localStorage.getItem(key)
    if (raw) {
      try {
        const parsed = JSON.parse(raw)
        if (Array.isArray(parsed) && parsed.length > 0) {
          customStallContracts.value[upper] = parsed.map((c: any) =>
            String(c.tenant_name || '').startsWith('แผง')
              ? { ...c, tenant_name: 'คุณสมศักดิ์ มั่งคั่ง' }
              : c
          )
          return
        }
      } catch {}
    }
  }
  const status = String(stallEdit.value.stall_status || '').toLowerCase()
  if (status === 'occupied') {
    const rawName = String(stallEdit.value.name || '').trim()
    const tenant = (rawName && !rawName.startsWith('แผง')) ? rawName : 'คุณสมศักดิ์ มั่งคั่ง'
    const rent = Number(stallEdit.value.monthly_rent || 2000)
    const defaultContract = {
      id: `cnt-${upper}-01`,
      contract_number: `CNT-2026-${upper}-01`,
      stall_ident: upper,
      stall_number: upper,
      tenant_name: tenant,
      tenant_phone: '081-889-9221',
      tenant_id_card: '3-1006-00421-44-1',
      shop_name: stallEdit.value.category || 'ร้านค้าทั่วไป',
      start_date: '2026-01-01',
      end_date: '2026-12-31',
      monthly_rental: rent,
      monthly_rent: rent,
      security_deposit: rent * 2,
      status: 'active',
      note: 'สัญญาเช่าแผงค้ามาตรฐาน ชำระทุกวันที่ 5 ของเดือน',
      created_at: '2026-01-01T08:00:00.000Z'
    }
    customStallContracts.value[upper] = [
      defaultContract,
      {
        id: `cnt-${upper}-history-02`, contract_number: `CNT-2025-${upper}-02`, stall_ident: upper, stall_number: upper,
        tenant_name: 'คุณกมลชนก แสงทอง', tenant_phone: '089-245-7810', shop_name: 'ร้านเครื่องดื่มสมุนไพร',
        start_date: '2025-01-01', end_date: '2025-12-31', monthly_rental: Math.max(1200, rent - 200),
        security_deposit: Math.max(2400, (rent - 200) * 2), status: 'expired',
        note: 'สิ้นสุดสัญญาตามกำหนดและส่งคืนพื้นที่เรียบร้อย', created_at: '2025-01-01T08:00:00.000Z'
      },
      {
        id: `cnt-${upper}-history-01`, contract_number: `CNT-2024-${upper}-01`, stall_ident: upper, stall_number: upper,
        tenant_name: 'คุณประวิทย์ วัฒนกิจ', tenant_phone: '086-441-9082', shop_name: 'ร้านของฝากและสินค้าแห้ง',
        start_date: '2024-03-01', end_date: '2024-11-30', monthly_rental: Math.max(1000, rent - 400),
        security_deposit: Math.max(2000, (rent - 400) * 2), status: 'terminated',
        note: 'ยกเลิกก่อนกำหนดตามคำขอของผู้เช่า', created_at: '2024-03-01T08:00:00.000Z'
      }
    ]
  } else {
    const rent = Number(stallEdit.value.monthly_rent || 2000)
    customStallContracts.value[upper] = [
      {
        id: `cnt-${upper}-history-02`, contract_number: `CNT-2025-${upper}-02`, stall_ident: upper, stall_number: upper,
        tenant_name: 'คุณกมลชนก แสงทอง', tenant_phone: '089-245-7810', shop_name: 'ร้านเครื่องดื่มสมุนไพร',
        start_date: '2025-01-01', end_date: '2025-12-31', monthly_rental: Math.max(1200, rent - 200),
        security_deposit: Math.max(2400, (rent - 200) * 2), status: 'expired',
        note: 'สิ้นสุดสัญญาตามกำหนด', created_at: '2025-01-01T08:00:00.000Z'
      }
    ]
  }
}

const loadAllCustomContracts = () => {
  if (!import.meta.client) return
  try {
    const prefix = `stall_contracts_custom_${marketId.value}_`
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.startsWith(prefix)) {
        const stallCode = key.replace(prefix, '').toUpperCase()
        const raw = localStorage.getItem(key)
        if (raw) {
          const parsed = JSON.parse(raw)
          if (Array.isArray(parsed)) {
            customStallContracts.value[stallCode] = parsed
          }
        }
      }
    }
  } catch {}
}

const syncStallDataFromContracts = (code: string) => {
  const upperCode = code.toUpperCase()
  const currentContracts = stallContractsList.value
  const activeContract = currentContracts.find((c: any) => {
    const st = String(c.status || c.contract_state || '').toLowerCase()
    return st === 'active' || st === 'ใช้งานอยู่'
  })

  if (activeContract) {
    stallEdit.value.stall_status = 'occupied'
    stallEdit.value.name = activeContract.tenant_name || activeContract.name || stallEdit.value.name
    if (activeContract.monthly_rental || activeContract.monthly_rent) {
      stallEdit.value.monthly_rent = Number(activeContract.monthly_rental || activeContract.monthly_rent)
    }

    const foundIdx = stalls.value.findIndex(s => String(s.ident || s.stall_ident).toUpperCase() === upperCode)
    if (foundIdx >= 0) {
      stalls.value[foundIdx] = {
        ...stalls.value[foundIdx],
        stall_status: 'occupied',
        status: 'occupied',
        name: stallEdit.value.name,
        monthly_rent: stallEdit.value.monthly_rent
      }
    }

    projects.value.forEach(p => {
      getAllFloorLayoutsForProject(p).forEach(layout => {
        const layoutStall = layout.stalls.find((s: any) => String(s.code || s.ident).toUpperCase() === upperCode)
        if (layoutStall) {
          layoutStall.status = 'occupied'
          layoutStall.stall_status = 'occupied'
          layoutStall.name = stallEdit.value.name
          layoutStall.monthlyRent = stallEdit.value.monthly_rent
        }
      })
      if (Array.isArray(p.stalls)) {
        const pStall = p.stalls.find((s: any) => String(s.code || s.ident).toUpperCase() === upperCode)
        if (pStall) {
          pStall.status = 'occupied'
          pStall.stall_status = 'occupied'
          pStall.name = stallEdit.value.name
          pStall.monthlyRent = stallEdit.value.monthly_rent
        }
      }
    })
    loadStallPayments(upperCode)
    loadStallMaintenance(upperCode)
  } else {
    const hasBookingOrApp = (marketBookings.value || []).some((b: any) => {
      const st = String(b.booking_status || b.status || '').toLowerCase()
      if (['rejected', 'cancelled', 'completed', 'terminated'].includes(st)) return false
      return String(b.stall_ident || b.stall_number || '').toUpperCase().split(',').map(s => s.trim()).includes(upperCode)
    })

    const targetStatus = hasBookingOrApp ? 'reserved' : 'available'
    stallEdit.value.stall_status = targetStatus
    if (targetStatus === 'available') {
      stallEdit.value.name = ''
    }

    const foundIdx = stalls.value.findIndex(s => String(s.ident || s.stall_ident).toUpperCase() === upperCode)
    if (foundIdx >= 0) {
      stalls.value[foundIdx] = {
        ...stalls.value[foundIdx],
        stall_status: targetStatus,
        status: targetStatus,
        name: targetStatus === 'available' ? '' : stalls.value[foundIdx].name
      }
    }

    projects.value.forEach(p => {
      getAllFloorLayoutsForProject(p).forEach(layout => {
        const layoutStall = layout.stalls.find((s: any) => String(s.code || s.ident).toUpperCase() === upperCode)
        if (layoutStall) {
          layoutStall.status = targetStatus
          layoutStall.stall_status = targetStatus
          if (targetStatus === 'available') layoutStall.name = ''
        }
      })
      if (Array.isArray(p.stalls)) {
        const pStall = p.stalls.find((s: any) => String(s.code || s.ident).toUpperCase() === upperCode)
        if (pStall) {
          pStall.status = targetStatus
          pStall.stall_status = targetStatus
          if (targetStatus === 'available') pStall.name = ''
        }
      }
    })
    loadStallPayments(upperCode)
    loadStallMaintenance(upperCode)
    loadStallSubtenants(upperCode)
  }
}

const openCreateContractModal = () => {
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  isEditingContract.value = false
  const seq = (stallContractsList.value.length + 1).toString().padStart(2, '0')
  const defaultNum = `CNT-${new Date().getFullYear()}-${code || 'STALL'}-${seq}`
  const today = new Date().toISOString().slice(0, 10)
  const nextYear = new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString().slice(0, 10)
  const rent = Number(stallEdit.value.monthly_rent) || 1800

  contractForm.value = {
    id: '',
    contractNumber: defaultNum,
    tenantName: stallEdit.value.name || '',
    tenantPhone: '',
    tenantIdCard: '',
    shopName: stallEdit.value.category || '',
    startDate: today,
    endDate: nextYear,
    monthlyRent: rent,
    depositAmount: rent * 2,
    status: 'active',
    note: ''
  }
  contractErrors.value = {
    contractNumber: '',
    tenantName: '',
    tenantPhone: '',
    startDate: '',
    endDate: '',
    monthlyRent: ''
  }
  showContractModal.value = true
}

const openEditContractModal = (contract: any) => {
  isEditingContract.value = true
  contractForm.value = {
    id: String(contract.id || ''),
    contractNumber: contract.contract_number || contract.contract_no || '',
    tenantName: contract.tenant_name || contract.name || '',
    tenantPhone: contract.tenant_phone || contract.phone || '',
    tenantIdCard: contract.tenant_id_card || contract.id_card || '',
    shopName: contract.shop_name || contract.business_name || '',
    startDate: contract.start_date || '',
    endDate: contract.end_date || '',
    monthlyRent: contract.monthly_rental || contract.monthly_rent || stallEdit.value.monthly_rent || 1800,
    depositAmount: contract.security_deposit || contract.deposit_amount || 0,
    status: contract.status || 'active',
    note: contract.note || contract.terms_and_conditions || ''
  }
  contractErrors.value = {
    contractNumber: '',
    tenantName: '',
    tenantPhone: '',
    startDate: '',
    endDate: '',
    monthlyRent: ''
  }
  showContractModal.value = true
}

const saveContract = async () => {
  contractErrors.value = {
    contractNumber: contractForm.value.contractNumber.trim() ? '' : 'กรุณาระบุเลขที่สัญญา',
    tenantName: contractForm.value.tenantName.trim() ? '' : 'กรุณาระบุชื่อ-นามสกุลผู้เช่า',
    tenantPhone: contractForm.value.tenantPhone.trim() ? '' : 'กรุณาระบุเบอร์โทรศัพท์ติดต่อ',
    startDate: contractForm.value.startDate ? '' : 'กรุณาระบุวันเริ่มสัญญา',
    endDate: contractForm.value.endDate ? '' : 'กรุณาระบุวันสิ้นสุดสัญญา',
    monthlyRent: Number(contractForm.value.monthlyRent) > 0 ? '' : 'กรุณาระบุค่าเช่าต่อเดือนที่ถูกต้อง'
  }
  if (Object.values(contractErrors.value).some(Boolean)) return

  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (!code) return

  const contractId = isEditingContract.value && contractForm.value.id
    ? contractForm.value.id
    : `CNT-${Date.now()}`

  const contractData = {
    id: contractId,
    contract_number: contractForm.value.contractNumber.trim(),
    stall_ident: code,
    stall_number: code,
    tenant_name: contractForm.value.tenantName.trim(),
    tenant_phone: contractForm.value.tenantPhone.trim(),
    tenant_id_card: contractForm.value.tenantIdCard.trim(),
    shop_name: contractForm.value.shopName.trim() || contractForm.value.tenantName.trim(),
    start_date: contractForm.value.startDate,
    end_date: contractForm.value.endDate,
    monthly_rental: Number(contractForm.value.monthlyRent) || 0,
    monthly_rent: Number(contractForm.value.monthlyRent) || 0,
    security_deposit: Number(contractForm.value.depositAmount) || 0,
    status: contractForm.value.status,
    note: contractForm.value.note?.trim() || '',
    created_at: new Date().toISOString()
  }

  const currentList = [...(customStallContracts.value[code] || [])]
  const existingIdx = currentList.findIndex(c => String(c.id) === String(contractId) || c.contract_number === contractData.contract_number)
  if (existingIdx >= 0) {
    currentList[existingIdx] = contractData
  } else {
    currentList.unshift(contractData)
  }
  customStallContracts.value[code] = currentList

  if (import.meta.client) {
    localStorage.setItem(getContractsStorageKey(code), JSON.stringify(currentList))
  }

  const cIdx = contracts.value.findIndex(c => String(c.id) === String(contractId) || c.contract_number === contractData.contract_number)
  if (cIdx >= 0) {
    contracts.value[cIdx] = { ...contracts.value[cIdx], ...contractData }
  } else {
    contracts.value.unshift(contractData)
  }

  const mcIdx = marketContracts.value.findIndex(c => String(c.id) === String(contractId) || c.contract_number === contractData.contract_number)
  if (mcIdx >= 0) {
    marketContracts.value[mcIdx] = { ...marketContracts.value[mcIdx], ...contractData }
  } else {
    marketContracts.value.unshift(contractData)
  }

  syncStallDataFromContracts(code)

  try {
    if (isEditingContract.value) {
      await put(`/contracts/${contractId}`, contractData)
    } else {
      await post('/contracts', contractData)
    }
  } catch {}

  showContractModal.value = false
  await showSuccess(
    isEditingContract.value ? 'บันทึกการแก้ไขสัญญาเช่าและอัปเดตข้อมูลแผงเรียบร้อย' : 'ออกสัญญาเช่าใหม่และอัปเดตสถานะแผงเป็น "มีผู้เช่า" เรียบร้อยแล้ว',
    'สำเร็จ'
  )
}

const terminateContract = async (contract: any) => {
  const confirmed = await showConfirm(
    `คุณต้องการยกเลิกหรือสิ้นสุดสัญญา "${contract.contract_number || `#${contract.id}`}" ก่อนกำหนดหรือไม่? ระบบจะปรับปรุงสถานะแผงเป็น "ว่าง" โดยอัตโนมัติหากไม่มีสัญญาอื่นที่ใช้งานอยู่`,
    {
      title: 'ยืนยันยกเลิกสัญญาเช่า',
      type: 'error',
      icon: 'trash',
      confirmText: 'ยืนยันยกเลิกสัญญา',
      cancelText: 'ไม่ยกเลิก'
    }
  )
  if (!confirmed) return

  const code = String(stallEdit.value.ident || contract.stall_ident || '').trim().toUpperCase()
  const contractId = String(contract.id || contract.contract_number)

  const currentList = [...(customStallContracts.value[code] || stallContractsList.value)]
  const idx = currentList.findIndex(c => String(c.id || c.contract_number) === contractId)
  if (idx >= 0) {
    currentList[idx] = {
      ...currentList[idx],
      status: 'cancelled',
      contract_state: 'cancelled',
      terminated_at: new Date().toISOString()
    }
    customStallContracts.value[code] = currentList
    if (import.meta.client && code) {
      localStorage.setItem(getContractsStorageKey(code), JSON.stringify(currentList))
    }
  } else {
    const updated = {
      ...contract,
      status: 'cancelled',
      contract_state: 'cancelled',
      terminated_at: new Date().toISOString()
    }
    customStallContracts.value[code] = [updated, ...currentList]
    if (import.meta.client && code) {
      localStorage.setItem(getContractsStorageKey(code), JSON.stringify(customStallContracts.value[code]))
    }
  }

  const cIdx = contracts.value.findIndex(c => String(c.id || c.contract_number) === contractId)
  if (cIdx >= 0) {
    contracts.value[cIdx] = { ...contracts.value[cIdx], status: 'cancelled', contract_state: 'cancelled' }
  }
  const mcIdx = marketContracts.value.findIndex(c => String(c.id || c.contract_number) === contractId)
  if (mcIdx >= 0) {
    marketContracts.value[mcIdx] = { ...marketContracts.value[mcIdx], status: 'cancelled', contract_state: 'cancelled' }
  }

  syncStallDataFromContracts(code)

  try {
    await post(`/contracts/${contract.id}/terminate`)
  } catch {}

  await showSuccess('สัญญาเช่าถูกยกเลิกแล้ว และสถานะแผงถูกปรับปรุงเรียบร้อย', 'สำเร็จ')
}

// 1. Subtenants Management
const stallSubtenants = ref<StallSubtenant[]>([])
const showAddSubtenantModal = ref(false)
const subtenantForm = ref({
  name: '',
  phone: '',
  idCard: '',
  businessType: '',
  startDate: '',
  endDate: '',
  subRentAmount: '' as number | string,
  note: '',
  photoUrl: ''
})
const subtenantErrors = ref({
  name: '',
  phone: '',
  businessType: '',
  startDate: ''
})

const handleSubtenantImageUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) {
    void showError('กรุณาเลือกไฟล์รูปภาพเท่านั้น (JPG, PNG, WebP)', 'ไฟล์ไม่ถูกต้อง')
    return
  }
  if (file.size > 2 * 1024 * 1024) {
    void showError('ขนาดไฟล์รูปภาพต้องไม่เกิน 2 MB', 'ไฟล์มีขนาดใหญ่เกินไป')
    return
  }
  const reader = new FileReader()
  reader.onload = () => {
    subtenantForm.value.photoUrl = String(reader.result || '')
  }
  reader.readAsDataURL(file)
}

const removeSubtenantPhoto = () => {
  subtenantForm.value.photoUrl = ''
}

const getSubtenantsStorageKey = (code: string) => `stall_subtenants_${marketId.value}_${code.toUpperCase()}`

const loadStallSubtenants = (code: string) => {
  if (!code) {
    stallSubtenants.value = []
    return
  }
  const key = getSubtenantsStorageKey(code)
  if (import.meta.client) {
    const raw = localStorage.getItem(key)
    if (raw) {
      try {
        stallSubtenants.value = JSON.parse(raw)
        return
      } catch {}
    }
  }
  const status = String(stallEdit.value.stall_status || '').toLowerCase()
  if (status === 'occupied') {
    stallSubtenants.value = [
      {
        id: `sub-${code}-01`,
        stallIdent: code,
        name: 'คุณวราภรณ์ สดใส (ผู้เช่าช่วง)',
        phone: '089-456-7890',
        idCard: '1-1004-99882-12-1',
        businessType: stallEdit.value.category || 'ร้านเครื่องดื่มและชาชงสด',
        startDate: '2026-01-01',
        endDate: '2026-12-31',
        subRentAmount: Number(stallEdit.value.monthly_rent || 1800) + 400,
        note: 'ได้รับความยินยอมเป็นลายลักษณ์อักษรจากผู้เช่าสัญญาหลักแล้ว',
        photoUrl: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=crop&w=200&q=80',
        status: 'active',
        createdAt: new Date().toISOString(),
        contractNumber: `CNT-2026-${code.toUpperCase()}-01`,
        sequence: 2
      },
      {
        id: `sub-${code}-history-01`, stallIdent: code, name: 'คุณชุติมา พรเจริญ', phone: '082-590-4412',
        idCard: '1-1004-88421-13-9', businessType: 'อาหารและเครื่องดื่ม', startDate: '2025-02-01', endDate: '2025-08-31',
        subRentAmount: Number(stallEdit.value.monthly_rent || 1800) + 200,
        note: 'สิ้นสุดการเช่าช่วงตามกำหนด', status: 'ended', createdAt: '2025-02-01T08:00:00.000Z',
        contractNumber: `CNT-2025-${code.toUpperCase()}-02`, sequence: 1
      }
    ]
  } else {
    stallSubtenants.value = []
  }
}

const openAddSubtenantModal = () => {
  subtenantForm.value = {
    name: '',
    phone: '',
    idCard: '',
    businessType: stallEdit.value.category || '',
    startDate: new Date().toISOString().slice(0, 10),
    endDate: '',
    subRentAmount: stallEdit.value.monthly_rent ? Number(stallEdit.value.monthly_rent) : '',
    note: '',
    photoUrl: ''
  }
  subtenantErrors.value = { name: '', phone: '', businessType: '', startDate: '' }
  showAddSubtenantModal.value = true
}

const saveSubtenant = () => {
  subtenantErrors.value = {
    name: subtenantForm.value.name ? '' : 'กรุณาระบุชื่อ-นามสกุล',
    phone: subtenantForm.value.phone ? '' : 'กรุณาระบุเบอร์โทรศัพท์',
    businessType: subtenantForm.value.businessType ? '' : 'กรุณาระบุประเภทสินค้าหรือกิจการ',
    startDate: subtenantForm.value.startDate ? '' : 'กรุณาระบุวันที่เริ่มสัญญา'
  }
  if (Object.values(subtenantErrors.value).some(Boolean)) return

  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  const newSub: StallSubtenant = {
    id: `sub-${Date.now()}`,
    stallIdent: code,
    name: subtenantForm.value.name.trim(),
    phone: subtenantForm.value.phone.trim(),
    idCard: subtenantForm.value.idCard?.trim() || undefined,
    businessType: subtenantForm.value.businessType.trim(),
    startDate: subtenantForm.value.startDate,
    endDate: subtenantForm.value.endDate || '-',
    subRentAmount: Number(subtenantForm.value.subRentAmount) || undefined,
    note: subtenantForm.value.note?.trim() || undefined,
    photoUrl: subtenantForm.value.photoUrl?.trim() || undefined,
    status: 'active',
    createdAt: new Date().toISOString()
  }

  stallSubtenants.value.unshift(newSub)
  if (import.meta.client && code) {
    localStorage.setItem(getSubtenantsStorageKey(code), JSON.stringify(stallSubtenants.value))
  }
  showAddSubtenantModal.value = false
  void showSuccess('บันทึกข้อมูลผู้เช่าช่วงเรียบร้อยแล้ว', 'เพิ่มผู้เช่าช่วงสำเร็จ')
}

const toggleSubtenantStatus = (id: string) => {
  const item = stallSubtenants.value.find(s => s.id === id)
  if (!item) return
  item.status = item.status === 'active' ? 'ended' : 'active'
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (import.meta.client && code) {
    localStorage.setItem(getSubtenantsStorageKey(code), JSON.stringify(stallSubtenants.value))
  }
  void showSuccess(`เปลี่ยนสถานะผู้เช่าช่วงเป็น "${item.status === 'active' ? 'กำลังเช่าช่วง' : 'สิ้นสุดสัญญา'}" แล้ว`, 'อัปเดตสถานะสำเร็จ')
}

const deleteSubtenant = async (id: string) => {
  const confirmed = await showConfirm(
    'คุณต้องการลบข้อมูลผู้เช่าช่วงรายนี้ใช่หรือไม่? การกระทำนี้ไม่สามารถย้อนกลับได้',
    {
      title: 'ยืนยันการลบผู้เช่าช่วง',
      type: 'danger',
      confirmText: 'ลบข้อมูล',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return
  stallSubtenants.value = stallSubtenants.value.filter(s => s.id !== id)
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (import.meta.client && code) {
    localStorage.setItem(getSubtenantsStorageKey(code), JSON.stringify(stallSubtenants.value))
  }
  void showSuccess('ลบข้อมูลผู้เช่าช่วงเรียบร้อยแล้ว', 'ลบสำเร็จ')
}

// 2. Maintenance & Permission Requests Management
const stallMaintenance = ref<StallMaintenanceItem[]>([])
const showAddMaintenanceModal = ref(false)
const maintenanceForm = ref({
  type: 'repair' as 'repair' | 'renovation' | 'utility' | 'other',
  title: '',
  description: '',
  requesterName: '',
  technicianName: '',
  requestDate: '',
  scheduledDate: '',
  estimatedCost: '' as number | string,
  status: 'pending' as 'pending' | 'in_progress' | 'completed' | 'rejected'
})
const maintenanceErrors = ref({
  title: '',
  requesterName: '',
  requestDate: ''
})

const getMaintenanceStorageKey = (code: string) => `stall_maintenance_${marketId.value}_${code.toUpperCase()}`

const loadStallMaintenance = (code: string) => {
  if (!code) {
    stallMaintenance.value = []
    return
  }
  const key = getMaintenanceStorageKey(code)
  if (import.meta.client) {
    const raw = localStorage.getItem(key)
    if (raw) {
      try {
        stallMaintenance.value = JSON.parse(raw)
        return
      } catch {}
    }
  }
  const status = String(stallEdit.value.stall_status || '').toLowerCase()
  const hasActiveContract = (customStallContracts.value[code.toUpperCase()] || []).some((c: any) => c.status === 'active') ||
    stallContractsList.value.some((c: any) => c.status === 'active')
  if (status === 'occupied' || hasActiveContract) {
    stallMaintenance.value = [
      {
        id: `maint-${code}-01`,
        stallIdent: code,
        type: 'repair',
        title: 'ขอซ่อมแซมกระเบื้องปูพื้นชำรุด',
        description: 'กระเบื้องหน้าแผงหลุดร่อนจำนวน 2 แผ่น ช่างส่วนกลางเข้าตรวจสอบแล้ว',
        requesterName: stallEdit.value.name ? `${stallEdit.value.name} (ผู้เช่า)` : 'ผู้เช่าแผง',
        technicianName: 'ช่างสมพงษ์ (งานอาคาร)',
        requestDate: '2026-02-14',
        scheduledDate: '2026-02-18',
        estimatedCost: 350,
        status: 'completed',
        createdAt: '2026-02-14T09:30:00.000Z'
      },
      {
        id: `maint-${code}-02`,
        stallIdent: code,
        type: 'utility',
        title: 'ขออนุญาตติดตั้งมิเตอร์ไฟฟ้าและเบรกเกอร์ย่อย',
        description: 'ขอเพิ่มปลั๊กพ่วงกำลังไฟ 15A สำหรับตู้แช่เครื่องดื่ม',
        requesterName: stallEdit.value.name ? `${stallEdit.value.name} (ผู้เช่า)` : 'ผู้เช่าแผง',
        technicianName: 'ฝ่ายช่างไฟฟ้าตลาด',
        requestDate: '2026-03-01',
        scheduledDate: '2026-03-05',
        estimatedCost: 1200,
        status: 'completed',
        createdAt: '2026-03-01T10:15:00.000Z',
        source: 'tenant'
      },
      {
        id: `maint-${code}-03`, stallIdent: code, type: 'repair', title: 'ตรวจเช็กระบบระบายน้ำตามแผน',
        description: 'เจ้าหน้าที่อาคารตรวจสภาพท่อระบายน้ำหลังแผงเพื่อป้องกันการอุดตัน',
        requesterName: 'ฝ่ายอาคารและสถานที่', technicianName: 'ทีมซ่อมบำรุงตลาด', requestDate: '2026-10-25',
        scheduledDate: '2026-10-30', estimatedCost: 0, status: 'in_progress', createdAt: '2026-10-25T09:00:00.000Z', source: 'staff'
      }
    ]
  } else {
    stallMaintenance.value = [
      {
        id: `maint-${code}-03`,
        stallIdent: code,
        type: 'repair',
        title: 'ตรวจระบบระบายน้ำประจำเดือน',
        description: 'ตรวจท่อระบายน้ำหลังแผง',
        requesterName: 'ฝ่ายอาคารและสถานที่',
        technicianName: 'ทีมซ่อมบำรุงตลาด',
        requestDate: '2026-10-25',
        scheduledDate: '2026-10-30',
        estimatedCost: 0,
        status: 'in_progress',
        createdAt: '2026-10-25T09:00:00.000Z',
        source: 'staff'
      }
    ]
  }
}

const openAddMaintenanceModal = () => {
  maintenanceForm.value = {
    type: 'repair',
    title: '',
    description: '',
    requesterName: '',
    technicianName: '',
    requestDate: new Date().toISOString().slice(0, 10),
    scheduledDate: '',
    estimatedCost: '',
    status: 'pending'
  }
  maintenanceErrors.value = { title: '', requesterName: '', requestDate: '' }
  showAddMaintenanceModal.value = true
}

const saveMaintenance = () => {
  maintenanceErrors.value = {
    title: maintenanceForm.value.title ? '' : 'กรุณาระบุหัวข้อเรื่องการซ่อมหรือคำขอ',
    requesterName: maintenanceForm.value.requesterName ? '' : 'กรุณาระบุผู้แจ้งเรื่อง',
    requestDate: maintenanceForm.value.requestDate ? '' : 'กรุณาระบุวันที่แจ้ง'
  }
  if (Object.values(maintenanceErrors.value).some(Boolean)) return

  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  const newItem: StallMaintenanceItem = {
    id: `maint-${Date.now()}`,
    stallIdent: code,
    type: maintenanceForm.value.type,
    title: maintenanceForm.value.title.trim(),
    description: maintenanceForm.value.description.trim(),
    requesterName: maintenanceForm.value.requesterName.trim(),
    technicianName: maintenanceForm.value.technicianName?.trim() || undefined,
    requestDate: maintenanceForm.value.requestDate,
    scheduledDate: maintenanceForm.value.scheduledDate || undefined,
    estimatedCost: Number(maintenanceForm.value.estimatedCost) || undefined,
    status: maintenanceForm.value.status,
    createdAt: new Date().toISOString()
  }

  stallMaintenance.value.unshift(newItem)
  if (import.meta.client && code) {
    localStorage.setItem(getMaintenanceStorageKey(code), JSON.stringify(stallMaintenance.value))
  }
  showAddMaintenanceModal.value = false
  void showSuccess('บันทึกรายการซ่อม / ขออนุญาตเรียบร้อยแล้ว', 'บันทึกสำเร็จ')
}

const updateMaintenanceStatus = (id: string, newStatus: 'completed' | 'in_progress' | 'pending' | 'rejected') => {
  const item = stallMaintenance.value.find(m => m.id === id)
  if (!item) return
  item.status = newStatus
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (import.meta.client && code) {
    localStorage.setItem(getMaintenanceStorageKey(code), JSON.stringify(stallMaintenance.value))
  }
  void showSuccess('อัปเดตสถานะการดำเนินการเรียบร้อยแล้ว', 'อัปเดตสำเร็จ')
}

const deleteMaintenance = async (id: string) => {
  const confirmed = await showConfirm(
    'คุณต้องการลบรายการนี้ใช่หรือไม่?',
    {
      title: 'ยืนยันการลบรายการ',
      type: 'danger',
      confirmText: 'ลบรายการ',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return
  stallMaintenance.value = stallMaintenance.value.filter(m => m.id !== id)
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (import.meta.client && code) {
    localStorage.setItem(getMaintenanceStorageKey(code), JSON.stringify(stallMaintenance.value))
  }
  void showSuccess('ลบรายการเรียบร้อยแล้ว', 'ลบสำเร็จ')
}

const formatMaintenanceType = (t: string) => {
  const map: Record<string, string> = {
    repair: 'แจ้งซ่อมแซม',
    renovation: 'ขออนุญาตต่อเติม',
    utility: 'ระบบไฟฟ้า/ประปา',
    other: 'คำขออื่นๆ'
  }
  return map[t] || 'แจ้งซ่อม'
}

const formatMaintenanceStatus = (s: string) => {
  const map: Record<string, string> = {
    pending: 'รอดำเนินการ',
    in_progress: 'กำลังดำเนินการ',
    completed: 'เสร็จสิ้น',
    rejected: 'ไม่อนุมัติ'
  }
  return map[s] || s
}

// 3. Payment History Management
const stallPayments = ref<StallPaymentItem[]>([])
const showAddPaymentModal = ref(false)
const paymentForm = ref({
  periodMonth: '',
  title: 'ค่าเช่าแผงรายเดือน',
  amount: '' as number | string,
  paidDate: '',
  paymentMethod: 'qr' as 'qr' | 'transfer' | 'cash' | 'credit',
  receiptNumber: ''
})
const paymentErrors = ref({
  periodMonth: '',
  amount: '',
  paidDate: ''
})

const getPaymentsStorageKey = (code: string) => `stall_payments_${marketId.value}_${code.toUpperCase()}`

const loadStallPayments = (code: string) => {
  if (!code) {
    stallPayments.value = []
    return
  }
  const key = getPaymentsStorageKey(code)
  if (import.meta.client) {
    const raw = localStorage.getItem(key)
    if (raw) {
      try {
        stallPayments.value = JSON.parse(raw)
        return
      } catch {}
    }
  }
  const status = String(stallEdit.value.stall_status || '').toLowerCase()
  const hasActiveContract = (customStallContracts.value[code.toUpperCase()] || []).some((c: any) => c.status === 'active') ||
    stallContractsList.value.some((c: any) => c.status === 'active')
  if (status === 'occupied' || hasActiveContract) {
    const rent = Number(stallEdit.value.monthly_rent || 1800)
    stallPayments.value = [
      {
        id: `pay-${code}-01`,
        stallIdent: code,
        invoiceNumber: `INV-2026-03-${code}`,
        periodMonth: 'งวดมีนาคม 2569',
        title: 'ค่าเช่าแผงประจำงวดเดือนมีนาคม',
        amount: rent,
        paidDate: '2026-03-02',
        paymentMethod: 'qr',
        receiptNumber: `REC-6903-${code}-01`,
        status: 'paid',
        createdAt: '2026-03-02T11:20:00.000Z',
        contractNumber: `CNT-2026-${code.toUpperCase()}-01`
      },
      {
        id: `pay-${code}-02`,
        stallIdent: code,
        invoiceNumber: `INV-2026-02-${code}`,
        periodMonth: 'งวดกุมภาพันธ์ 2569',
        title: 'ค่าเช่าแผงประจำงวดเดือนกุมภาพันธ์',
        amount: rent,
        paidDate: '2026-02-01',
        paymentMethod: 'transfer',
        receiptNumber: `REC-6902-${code}-01`,
        status: 'paid',
        createdAt: '2026-02-01T14:45:00.000Z',
        contractNumber: `CNT-2026-${code.toUpperCase()}-01`
      },
      {
        id: `pay-${code}-03`,
        stallIdent: code,
        invoiceNumber: `INV-2026-01-${code}`,
        periodMonth: 'งวดมกราคม 2569',
        title: 'ค่าธรรมเนียมแรกเข้าและเงินประกัน',
        amount: Number(inspectStallEntryFee.value || 36000),
        paidDate: '2026-01-05',
        paymentMethod: 'transfer',
        receiptNumber: `REC-6901-${code}-00`,
        status: 'paid',
        createdAt: '2026-01-05T10:00:00.000Z',
        contractNumber: `CNT-2026-${code.toUpperCase()}-01`
      },
      {
        id: `pay-${code}-04`, stallIdent: code, invoiceNumber: `INV-2026-10-${code}`,
        periodMonth: 'งวดตุลาคม 2569', title: 'ค่าเช่าแผงประจำงวดเดือนตุลาคม', amount: rent,
        paymentMethod: 'transfer', status: 'overdue', createdAt: '2026-10-05T08:00:00.000Z',
        contractNumber: `CNT-2026-${code.toUpperCase()}-01`
      },
      {
        id: `pay-${code}-history-01`, stallIdent: code, invoiceNumber: `INV-2025-12-${code}`,
        periodMonth: 'งวดธันวาคม 2568', title: 'ค่าเช่างวดสุดท้ายของสัญญาเดิม', amount: Math.max(1200, rent - 200),
        paidDate: '2025-12-03', paymentMethod: 'cash', receiptNumber: `REC-6812-${code}-01`, status: 'paid',
        createdAt: '2025-12-03T10:00:00.000Z', contractNumber: `CNT-2025-${code.toUpperCase()}-02`
      }
    ]
  } else {
    stallPayments.value = []
  }
}

const openAddPaymentModal = () => {
  paymentForm.value = {
    periodMonth: 'งวดปัจจุบัน',
    title: 'ค่าเช่าแผงรายเดือน',
    amount: stallEdit.value.monthly_rent ? Number(stallEdit.value.monthly_rent) : 1800,
    paidDate: new Date().toISOString().slice(0, 10),
    paymentMethod: 'qr',
    receiptNumber: `REC-${Date.now().toString().slice(-6)}`
  }
  paymentErrors.value = { periodMonth: '', amount: '', paidDate: '' }
  showAddPaymentModal.value = true
}

const savePayment = () => {
  paymentErrors.value = {
    periodMonth: paymentForm.value.periodMonth ? '' : 'กรุณาระบุงวดเดือนหรือรอบบิล',
    amount: paymentForm.value.amount ? '' : 'กรุณาระบุจำนวนเงิน',
    paidDate: paymentForm.value.paidDate ? '' : 'กรุณาระบุวันที่ชำระ'
  }
  if (Object.values(paymentErrors.value).some(Boolean)) return

  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  const newPay: StallPaymentItem = {
    id: `pay-${Date.now()}`,
    stallIdent: code,
    invoiceNumber: `INV-${Date.now().toString().slice(-6)}`,
    periodMonth: paymentForm.value.periodMonth.trim(),
    title: paymentForm.value.title,
    amount: Number(paymentForm.value.amount),
    paidDate: paymentForm.value.paidDate,
    paymentMethod: paymentForm.value.paymentMethod,
    receiptNumber: paymentForm.value.receiptNumber?.trim() || undefined,
    status: 'paid',
    createdAt: new Date().toISOString()
  }

  stallPayments.value.unshift(newPay)
  if (import.meta.client && code) {
    localStorage.setItem(getPaymentsStorageKey(code), JSON.stringify(stallPayments.value))
  }
  showAddPaymentModal.value = false
  void showSuccess('บันทึกการชำระเงินเรียบร้อยแล้ว', 'บันทึกสำเร็จ')
}

const deletePayment = async (id: string) => {
  const confirmed = await showConfirm(
    'คุณต้องการลบรายการชำระเงินนี้ใช่หรือไม่?',
    {
      title: 'ยืนยันการลบรายการชำระเงิน',
      type: 'danger',
      confirmText: 'ลบรายการ',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return
  stallPayments.value = stallPayments.value.filter(p => p.id !== id)
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (import.meta.client && code) {
    localStorage.setItem(getPaymentsStorageKey(code), JSON.stringify(stallPayments.value))
  }
  void showSuccess('ลบรายการชำระเงินเรียบร้อยแล้ว', 'ลบสำเร็จ')
}

const formatPaymentMethod = (m: string) => {
  const map: Record<string, string> = {
    qr: 'สแกน QR PromptPay',
    transfer: 'โอนผ่านบัญชี',
    cash: 'เงินสด',
    credit: 'บัตรเครดิต/เดบิต'
  }
  return map[m] || m
}

const maintenanceTypeOptions = [
  { value: 'repair', label: 'แจ้งซ่อมแซม' },
  { value: 'renovation', label: 'ขออนุญาตต่อเติม / ดัดแปลง' },
  { value: 'utility', label: 'ระบบไฟฟ้า / ประปา' },
  { value: 'other', label: 'คำขออื่นๆ' }
]

const maintenanceStatusOptions = [
  { value: 'pending', label: 'รอดำเนินการ' },
  { value: 'in_progress', label: 'กำลังดำเนินการ' },
  { value: 'completed', label: 'เสร็จสิ้น' },
  { value: 'rejected', label: 'ไม่อนุมัติ' }
]

const paymentMethodOptions = [
  { value: 'qr', label: 'สแกน QR PromptPay' },
  { value: 'transfer', label: 'โอนผ่านบัญชีธนาคาร' },
  { value: 'cash', label: 'เงินสด' },
  { value: 'credit', label: 'บัตรเครดิต / เดบิต' }
]

const filteredBookingsAndApps = computed(() =>
  bookingTypeFilter.value === 'all'
    ? stallBookingsAndApps.value
    : stallBookingsAndApps.value.filter(item => item.typeKey === bookingTypeFilter.value)
)
const stallBookingsAndApps = computed(() => {
  const target = String(stallEdit.value.ident || stallEdit.value.stall_ident || '').trim().toUpperCase()
  if (!target) return []
  const items: Array<{
    key: string
    id?: string | number
    typeKey: 'booking' | 'application' | 'waitlist'
    typeLabel: string
    refNumber: string
    applicantName: string
    phone?: string
    dateLabel: string
    dateText: string
    statusText: string
    badgeClass: string
    targetUrl: string
  }> = []

  marketBookings.value.forEach((b: any) => {
    const codes = String(b.stall_number || b.stall_ident || b.stall_code || '')
      .split(',')
      .map(s => s.trim().toUpperCase())
    if (!codes.includes(target)) return
    const st = String(b.booking_status || b.status || '').toLowerCase()
    const statusMap: Record<string, { label: string; cls: string }> = {
      pending: { label: 'รอตรวจสอบวันนัด', cls: 'is-reserved' },
      pending_approval: { label: 'รออนุมัตินัดหมาย', cls: 'is-reserved' },
      pending_slot: { label: 'รอเลือกวันนัดชม', cls: 'is-reserved' },
      approved: { label: 'นัดหมายสำเร็จ', cls: 'is-available' },
      interested: { label: 'สนใจเช่าแผง', cls: 'is-available' },
      lease_applied: { label: 'ยื่นคำร้องขอเช่าแล้ว', cls: 'is-info' },
      completed: { label: 'ทำสัญญาเช่าสำเร็จ', cls: 'is-available' },
      waiting_queue: { label: 'รอคิวสำรอง', cls: 'is-reserved' },
      waiting: { label: 'รอคิวสำรอง', cls: 'is-reserved' },
      offered: { label: 'ได้รับสิทธิ์แทน', cls: 'is-info' },
      no_show: { label: 'ไม่มาตามนัด', cls: 'is-cancelled' },
      not_interested: { label: 'ไม่ประสงค์เช่าต่อ', cls: 'is-occupied' },
      withdrawn: { label: 'ถอนคิวแล้ว', cls: 'is-occupied' },
      expired: { label: 'หมดเวลาใช้สิทธิ์', cls: 'is-occupied' },
      cancelled: { label: 'ยกเลิกแล้ว', cls: 'is-occupied' },
      rejected: { label: 'ไม่อนุมัติ', cls: 'is-cancelled' }
    }
    const mapped = statusMap[st] || { label: 'รอดำเนินการ', cls: 'is-reserved' }
    const bookingRef = b.booking_number || `#${b.id}`
    items.push({
      key: `booking-${b.id || b.booking_number}`,
      id: b.id,
      typeKey: 'booking',
      typeLabel: 'นัดเข้าชมพื้นที่',
      refNumber: bookingRef,
      applicantName: b.tenant_name || b.applicant_name || '-',
      phone: b.tenant_phone || b.applicant_phone || '',
      dateLabel: b.visit_date ? 'วันนัดชมแผง' : 'วันที่ทำรายการ',
      dateText: b.visit_date ? `${formatThaiDateOnly(b.visit_date)}${b.visit_time_slot ? ' · ' + b.visit_time_slot : ''}` : formatStallDate(b.created_at),
      statusText: mapped.label,
      badgeClass: mapped.cls,
      targetUrl: b.id ? `/bookings/${b.id}` : `/bookings?search=${encodeURIComponent(b.booking_number || target)}`
    })
  })

  marketApplications.value.forEach((a: any) => {
    const codes = String(a.stall_number || a.stall_ident || a.stall_code || '')
      .split(',')
      .map(s => s.trim().toUpperCase())
    if (!codes.includes(target)) return
    const st = String(a.application_state || a.status || '').toLowerCase()
    const isApproved = st === 'approved' || st === 'completed'
    const isRejected = st === 'rejected' || st === 'cancelled'
    const appRef = a.application_number || `#${a.id}`
    items.push({
      key: `app-${a.id || a.application_number}`,
      id: a.id,
      typeKey: 'application',
      typeLabel: 'คำร้องขอเช่าแผง',
      refNumber: appRef,
      applicantName: a.applicant_name || a.tenant_name || '-',
      phone: a.applicant_phone || a.tenant_phone || '',
      dateLabel: 'วันที่ยื่นคำร้อง',
      dateText: formatStallDate(a.created_at),
      statusText: isApproved ? 'อนุมัติแล้ว' : (isRejected ? 'ไม่อนุมัติ' : 'รอพิจารณาอนุมัติ'),
      badgeClass: isApproved ? 'is-available' : (isRejected ? 'is-cancelled' : 'is-reserved'),
      targetUrl: a.id ? `/applications/${a.id}` : `/applications?search=${encodeURIComponent(a.application_number || target)}`
    })
  })

  if (items.length === 0 && String(stallEdit.value.stall_status || '').toLowerCase() !== 'available') {
    const rawName = String(stallEdit.value.name || '').trim()
    const tenant = (rawName && !rawName.startsWith('แผง')) ? rawName : (currentStallContract.value?.tenant_name || 'สมชาย ใจดี')
    const phone = currentStallContract.value?.tenant_phone || '081-234-5678'
    const cleanCode = target.replace(/[^A-Z0-9]/g, '')
    const refCode = cleanCode.slice(-6) || '015124'
    items.push(
      {
        key: `mock-booking-${target}-01`,
        id: 118,
        typeKey: 'booking',
        typeLabel: 'นัดเข้าชมพื้นที่',
        refNumber: `BK-20261027-${refCode}`,
        applicantName: tenant,
        phone: phone,
        dateLabel: 'วันนัดชมแผง',
        dateText: '2026-11-01 · 13.00 - 14.00 น.',
        statusText: 'ยื่นคำร้องขอเช่าแล้ว',
        badgeClass: 'is-info',
        targetUrl: '/bookings'
      },
      {
        key: `mock-application-${target}-01`,
        id: 33,
        typeKey: 'application',
        typeLabel: 'คำร้องขอเช่าแผง',
        refNumber: `APP-20261027-${refCode}`,
        applicantName: tenant,
        phone: phone,
        dateLabel: 'วันที่ยื่นคำร้อง',
        dateText: '27 ต.ค. 2569 22:06',
        statusText: 'รอพิจารณาอนุมัติ',
        badgeClass: 'is-reserved',
        targetUrl: `/applications?search=${encodeURIComponent(target)}`
      }
    )
  }

  return items
})

const stallActivityTimeline = computed(() => {
  const events: Array<{
    key: string
    category: string
    dotClass?: string
    title: string
    date: string
    detail: string
  }> = [
    ...stallContractsList.value.map((contract: any) => ({
      key: `contract-${contract.id || contract.contract_number}`,
      category: 'สัญญาเช่า',
      dotClass: contract.status === 'active' ? 'is-green' : 'is-gray',
      title: contract.status === 'active' ? 'เริ่มสัญญาเช่าปัจจุบัน' : (contract.status === 'terminated' || contract.status === 'cancelled' ? 'ยกเลิกสัญญาเช่า' : 'สิ้นสุดสัญญาเช่าเดิม'),
      date: contract.start_date ? `${formatThaiDateOnly(contract.start_date)} ถึง ${formatThaiDateOnly(contract.end_date)} · ระยะสัญญา ${getContractDurationText(contract.start_date, contract.end_date)}` : (contract.created_at || '-'),
      detail: `เลขที่สัญญา ${contract.contract_number || '-'} · ผู้เช่า: ${contract.tenant_name || '-'} · อัตราค่าเช่า ฿${Number(contract.monthly_rental || contract.monthly_rent || 0).toLocaleString('th-TH')}/เดือน`
    })),
    ...stallPayments.value.map(pay => ({
      key: `pay-${pay.id}`,
      category: 'การชำระเงิน',
      dotClass: pay.status === 'paid' ? 'is-green' : 'is-amber',
      title: pay.status === 'paid' ? `รับชำระเงิน (${pay.periodMonth})` : `แจ้งหนี้ค้างชำระ (${pay.periodMonth})`,
      date: pay.paidDate || formatStallDate(pay.createdAt),
      detail: `${pay.title} · จำนวน ฿${Number(pay.amount || 0).toLocaleString('th-TH')}${pay.receiptNumber ? ' · ใบเสร็จ ' + pay.receiptNumber : ''}`
    })),
    ...stallMaintenance.value.map(item => ({
      key: `maint-${item.id}`,
      category: 'งานซ่อม/ขออนุญาต',
      dotClass: item.status === 'completed' ? 'is-green' : 'is-amber',
      title: `${formatMaintenanceType(item.type)}: ${item.title}`,
      date: item.requestDate,
      detail: `ผู้แจ้ง: ${item.requesterName || '-'} · สถานะ: ${formatMaintenanceStatus(item.status)}${item.technicianName ? ' · ผู้รับผิดชอบ: ' + item.technicianName : ''}`
    })),
    ...stallBookingsAndApps.value.map(item => ({
      key: `bk-${item.key}`,
      category: item.typeLabel,
      dotClass: 'is-blue',
      title: `${item.typeLabel} (${item.refNumber})`,
      date: item.dateText,
      detail: `ผู้ทำรายการ: ${item.applicantName} · สถานะ: ${item.statusText}`
    }))
  ]
  return events
})

const stallContractsList = computed(() => {
  const target = String(stallEdit.value.ident || stallEdit.value.stall_ident || '').trim().toUpperCase()
  if (!target) return []
  const custom = customStallContracts.value[target] || []
  const contractMap = new Map<string, any>()
  ;[...marketContracts.value, ...contracts.value].forEach((c: any) => {
    const code = String(c.stall_ident || c.stall_number || c.stall_code || '').toUpperCase()
    if (!code.split(',').map(s => s.trim()).includes(target) && code !== target) return
    const key = String(c.id || c.contract_number || c.contract_no || Math.random())
    contractMap.set(key, c)
  })
  custom.forEach((c: any) => {
    const key = String(c.id || c.contract_number || c.contract_no || Math.random())
    contractMap.set(key, c)
  })
  return Array.from(contractMap.values())
})

const activeStallContracts = computed(() =>
  stallContractsList.value.filter((c: any) => {
    const st = String(c.status || c.contract_state || '').toLowerCase()
    return st === 'active' || st === 'draft' || st === 'ใช้งานอยู่'
  })
)

const historyStallContracts = computed(() =>
  stallContractsList.value.filter((c: any) => {
    const st = String(c.status || c.contract_state || '').toLowerCase()
    return st !== 'active' && st !== 'draft' && st !== 'ใช้งานอยู่'
  })
)

const displayedStallContracts = computed(() =>
  contractSubTab.value === 'active' ? activeStallContracts.value : historyStallContracts.value
)

const lastLoadedExtraStallCode = ref('')

const loadStallExtraDetails = async (forceReloadLocal = false) => {
  const code = String(stallEdit.value.ident || stallEdit.value.stall_ident || '').trim()
  if (!code) return
  const upperCode = code.toUpperCase()
  if (forceReloadLocal || lastLoadedExtraStallCode.value !== upperCode) {
    loadStallContracts(code)
    loadStallSubtenants(code)
    loadStallMaintenance(code)
    loadStallPayments(code)
    lastLoadedExtraStallCode.value = upperCode
  }
  const stored = stalls.value.find(s => String(s.ident || s.stall_ident).toUpperCase() === upperCode)
  const stallId = stored?.id || (stallEdit.value as any)?.id
  try {
    const [cRes, hRes] = await Promise.allSettled([
      get(`/contracts?stall_ident=${encodeURIComponent(code)}`),
      stallId ? get(`/stalls/${stallId}/history`) : Promise.resolve([])
    ])
    if (cRes.status === 'fulfilled' && Array.isArray(unwrap(cRes.value))) {
      contracts.value = unwrap(cRes.value)
    }
    if (hRes.status === 'fulfilled' && Array.isArray(unwrap(hRes.value))) {
      stallHistory.value = unwrap(hRes.value)
    }
  } catch {}
}

const fullDetailSettling = ref(false)
const stallDrawerExpanded = ref(false)
const isStallDrawerClosing = ref(false)
const isProjectPanelMorphing = ref(false)
const isProjectMorphSwap = ref(false)
const isProjectPanelReturning = ref(false)
const STALL_DRAWER_MORPH_DURATION = 420
const STALL_DRAWER_RETURN_DURATION = 620
let stallDrawerResizeTimer: ReturnType<typeof setTimeout> | null = null
let stallDrawerSwapTimer: ReturnType<typeof setTimeout> | null = null
const clearStallDrawerResizeTimer = () => {
  if (!stallDrawerResizeTimer) return
  clearTimeout(stallDrawerResizeTimer)
  stallDrawerResizeTimer = null
}
const clearStallDrawerSwapTimer = () => {
  if (!stallDrawerSwapTimer) return
  clearTimeout(stallDrawerSwapTimer)
  stallDrawerSwapTimer = null
}
const expandStallDrawerOnNextFrame = () => {
  if (!import.meta.client) {
    stallDrawerExpanded.value = true
    return
  }
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      stallDrawerExpanded.value = true
    })
  })
}
const showSelectedStallPanel = (replaceExisting = false) => {
  // ระหว่างที่กรอบโครงการกำลังขยาย ให้เปลี่ยนเฉพาะข้อมูลเป้าหมาย
  // และปล่อย animation เดิมทำงานต่อจนจบ ไม่เริ่มขยายซ้ำอีกรอบ
  if (isProjectPanelMorphing.value && !showInspectStallOverlay.value) return
  clearStallDrawerResizeTimer()
  clearStallDrawerSwapTimer()
  isProjectPanelReturning.value = false
  showStallEditor.value = false

  if (replaceExisting) {
    // เปลี่ยนแผงขณะ drawer เปิดอยู่: คงกรอบเดิมและแทนที่เฉพาะข้อมูล ไม่เล่น morph ซ้ำ
    isProjectPanelMorphing.value = false
    isProjectMorphSwap.value = false
    isStallDrawerClosing.value = false
    stallDrawerExpanded.value = true
    showInspectStallOverlay.value = true
    return
  }

  if (import.meta.client && isMapOverviewMode.value && isInspectMode.value && selectedProject.value) {
    // ขยาย panel โครงการจริงก่อน แล้วสลับเฉพาะเนื้อหาภายในเมื่อกรอบถึงขนาดปลายทาง
    showInspectStallOverlay.value = false
    stallDrawerExpanded.value = false
    isProjectPanelMorphing.value = true
    stallDrawerResizeTimer = window.setTimeout(() => {
      isProjectMorphSwap.value = true
      stallDrawerExpanded.value = true
      showInspectStallOverlay.value = true
      isProjectPanelMorphing.value = false
      stallDrawerResizeTimer = null
      stallDrawerSwapTimer = window.setTimeout(() => {
        isProjectMorphSwap.value = false
        stallDrawerSwapTimer = null
      }, 240)
    }, STALL_DRAWER_MORPH_DURATION)
    return
  }

  showInspectStallOverlay.value = true
  if (isInspectMode.value) expandStallDrawerOnNextFrame()
}
const openStallFullDetail = (tab?: 'info' | 'contracts' | 'subtenants' | 'payments' | 'maintenance' | 'bookings' | 'history') => {
  clearStallDrawerResizeTimer()
  isStallDrawerClosing.value = false
  stallDetailTab.value = tab || 'info'
  if (import.meta.client) {
    // คงเนื้อหาสรุปเดิมไว้ระหว่างที่กรอบขยาย แล้วจึงสลับเป็นรายละเอียดเต็ม
    fullDetailSettling.value = true
    stallDrawerExpanded.value = true
    stallDrawerResizeTimer = window.setTimeout(() => {
      stallDrawerMode.value = 'full-detail'
      fullDetailSettling.value = false
      stallDrawerResizeTimer = null
      void loadStallExtraDetails(false)
    }, STALL_DRAWER_MORPH_DURATION)
  } else {
    stallDrawerExpanded.value = true
    stallDrawerMode.value = 'full-detail'
    void loadStallExtraDetails(false)
  }
}

// URL สำหรับกลับมายังแผนที่ในสภาพเดิม (โซน + แผง + โหมด drawer + แท็บ ที่เปิดค้างไว้)
const stallReturnUrl = computed(() => router.resolve({
  path: route.path,
  query: {
    ...route.query,
    ...(selectedProject.value?.id ? { project: String(selectedProject.value.id) } : {}),
    stall: String(stallEdit.value.ident || ''),
    drawer: stallDrawerMode.value === 'full-detail' ? 'full' : 'detail',
    tab: stallDetailTab.value
  }
}).fullPath)

const withStallReturn = (url: string) => `${url}${url.includes('?') ? '&' : '?'}return=${encodeURIComponent(stallReturnUrl.value)}`

const restoreReturnedStall = async () => {
  const code = returnStallCode.value
  if (!code || !selectedProject.value || showInspectStallOverlay.value) return
  const found = (selectedProject.value.stalls || []).find((s: any) =>
    String(s.code || s.ident || s.stall_ident || '').toUpperCase() === code.toUpperCase()
  )
  await selectStall(found ? { ...found, code: found.code || found.ident || code } : { code })
  if (String(route.query.drawer || '') === 'full') {
    const allowedTabs = ['info', 'contracts', 'subtenants', 'payments', 'maintenance', 'bookings', 'history']
    const tab = String(route.query.tab || 'info')
    openStallFullDetail((allowedTabs.includes(tab) ? tab : 'info') as any)
  }
}

const jumpToActionTab = (targetTab?: string) => {
  if (targetTab && ['info', 'contracts', 'subtenants', 'payments', 'maintenance', 'bookings', 'history'].includes(targetTab)) {
    stallDetailTab.value = targetTab as any
  }
}

const backToMiniDetail = () => {
  clearStallDrawerResizeTimer()
  fullDetailSettling.value = false
  stallDrawerExpanded.value = false
  if (import.meta.client) {
    // คงหน้ารายละเอียดไว้ระหว่างที่กรอบยุบ แล้วค่อยคืนเป็นการ์ดสรุป
    stallDrawerResizeTimer = window.setTimeout(() => {
      stallDrawerMode.value = 'detail'
      stallDrawerResizeTimer = null
    }, STALL_DRAWER_MORPH_DURATION)
  } else {
    stallDrawerMode.value = 'detail'
  }
}
const stallDetailUrl = computed(() => {
  const params = new URLSearchParams({ market: marketId.value, stall: String(stallEdit.value.ident || stallEdit.value.stall_ident || '') })
  if (selectedProject.value?.id) params.set('project', selectedProject.value.id)
  if (isInspectMode.value) {
    params.set('mode', 'inspect')
    params.set('view', 'map')
  }
  return `/stalls/detail?${params.toString()}`
})
const DEFAULT_MARKET_CATEGORIES = [
  'อาหารและเครื่องดื่ม',
  'อาหารสด ผักผลไม้ และเนื้อสัตว์',
  'เสื้อผ้าและเครื่องประดับ',
  'ของใช้เบ็ดเตล็ดและของใช้ในบ้าน',
  'ต้นไม้และอุปกรณ์จัดสวน',
  'สัตว์เลี้ยงและอุปกรณ์',
  'หัตถกรรมและของแต่งบ้าน',
  'ศิลปะและของสะสม',
  'สินค้าโอทอปและของฝาก'
]
const marketCategoryOptions = computed(() => {
  const cfg = getMarketGisConfig(marketId.value)
  const overridden = applyMarketOverride({ id: resolvedMarketID.value || cfg.id, slug: marketId.value || cfg.slug })
  const configuredCats: string[] = Array.isArray(overridden?.product_categories)
    ? overridden.product_categories
    : []
  const gisCats = [
    ...(cfg?.projects || []).map(p => String(p.category || '').trim()),
    ...projects.value.map(p => String(p.category || '').trim())
  ].filter(c => c && c !== 'ยังไม่ระบุหมวดสินค้า')
  const merged = Array.from(new Set([...configuredCats, ...gisCats])).filter(Boolean)
  const list = merged.length > 0 ? merged : DEFAULT_MARKET_CATEGORIES
  return list.map(cat => ({ value: cat, label: cat }))
})
const stallCategoryOptions = computed(() => {
  const base = marketCategoryOptions.value
  const current = String(stallEdit.value.category || '').trim()
  if (current && !base.some(opt => opt.value === current)) {
    return [{ value: current, label: current }, ...base]
  }
  return base
})
const facilityTypeOptions = [{ value: 'toilet', label: 'ห้องน้ำ' }, { value: 'parking', label: 'ที่จอดรถ' }, { value: 'exit', label: 'ทางเข้า–ออก' }, { value: 'service', label: 'จุดบริการ' }, { value: 'waste', label: 'จุดทิ้งขยะ' }]
const zoneForm = ref({ name: '', sectionNumber: '', category: '', floorCount: 1 })
const drawnZoneCoordinates = ref<[number, number][]>([])
const pendingStallLayout = ref<{ projectId: string; coordinates: [number, number][]; floorId: string } | null>(null)
const isCreatingNewStall = computed(() => Boolean(pendingStallLayout.value))
const facilityForm = ref<{ name: string; type: string; coordinates: [number, number]; projectId?: string }>({ name: '', type: 'toilet', coordinates: [100.5, 13.75] })
const unwrap = (value: any) => value?.data ?? value
const stallStatusOverrides = computed(() => {
  const overrides: Record<string, 'available' | 'reserved' | 'occupied'> = {}
  const setStatus = (rawCode: string, status: 'available' | 'reserved' | 'occupied') => {
    const clean = String(rawCode || '').trim()
    if (!clean) return
    overrides[clean] = status
    overrides[clean.toUpperCase()] = status
  }

  stalls.value.forEach((s: any) => {
    const code = String(s.ident || s.stall_ident || s.stall_number || '').trim()
    const rawStatus = String(s.stall_status || s.status || '').toLowerCase()
    if (!code) return
    if (rawStatus === 'reserved' || rawStatus === 'occupied') {
      setStatus(code, rawStatus)
    }
  })

  const inactiveBookingStatuses = new Set([
    'waiting_queue',
    'rejected',
    'not_interested',
    'cancelled',
    'no_show',
    'lease_applied',
    'completed',
    'terminated',
    'closed'
  ])
  marketBookings.value.forEach((b: any) => {
    const st = String(b.booking_status || b.status || '').toLowerCase()
    if (!st || inactiveBookingStatuses.has(st)) return
    const rawCodes = String(b.stall_ident || b.stall_number || b.stall_code || '')
    rawCodes.split(',').forEach((c) => {
      const clean = c.trim()
      if (clean && overrides[clean.toUpperCase()] !== 'occupied') {
        setStatus(clean, 'reserved')
      }
    })
  })

  const inactiveAppStatuses = new Set(['rejected', 'cancelled', 'terminated', 'closed', 'completed'])
  marketApplications.value.forEach((a: any) => {
    const st = String(a.application_state || a.status || '').toLowerCase()
    if (!st || inactiveAppStatuses.has(st)) return
    const rawCodes = String(a.stall_ident || a.stall_number || a.stall_code || '')
    rawCodes.split(',').forEach((c) => {
      const clean = c.trim()
      if (clean && overrides[clean.toUpperCase()] !== 'occupied') {
        setStatus(clean, 'reserved')
      }
    })
  })

  const allCustomContracts = Object.values(customStallContracts.value).flat()
  const contractMap = new Map<string, any>()
  ;[...marketContracts.value, ...contracts.value].forEach((c: any) => {
    const key = String(c.id || c.contract_number || c.contract_no || '')
    if (key) contractMap.set(key, c)
  })
  allCustomContracts.forEach((c: any) => {
    const key = String(c.id || c.contract_number || c.contract_no || '')
    if (key) contractMap.set(key, c)
  })

  contractMap.forEach((c: any) => {
    const st = String(c.contract_state || c.status || '').toLowerCase()
    if (st === 'terminated' || st === 'cancelled' || st === 'expired') return
    const rawCodes = String(c.stall_ident || c.stall_number || c.stall_code || '')
    rawCodes.split(',').forEach((code) => {
      const clean = code.trim()
      if (clean) setStatus(clean, 'occupied')
    })
  })

  return overrides
})

const demandViewMode = ref(false)
const demandLevelFilter = ref<'all' | 'high' | 'medium' | 'low' | 'none'>('all')

const stallDemandMap = computed<Record<string, {
  demandLevel: 'high' | 'medium' | 'low' | 'none'
  demandScore: number
  waitlistCount: number
  bookingCount: number
}>>(() => {
  const result: Record<string, {
    demandLevel: 'high' | 'medium' | 'low' | 'none'
    demandScore: number
    waitlistCount: number
    bookingCount: number
  }> = {}

  // 1. Map stalls from API
  const stallApiMap = new Map<string, any>()
  stalls.value.forEach((s: any) => {
    const code = String(s.code || s.ident || s.stall_ident || s.stall_number || '').trim().toUpperCase()
    if (code) stallApiMap.set(code, s)
  })

  // 2. Count active bookings per stall
  const inactiveBookingStatuses = new Set([
    'rejected',
    'not_interested',
    'cancelled',
    'no_show',
    'completed',
    'terminated',
    'closed'
  ])
  const activeBookingsByStall = new Map<string, number>()
  marketBookings.value.forEach((b: any) => {
    const st = String(b.booking_status || b.status || '').toLowerCase()
    if (!st || inactiveBookingStatuses.has(st)) return
    const rawCodes = String(b.stall_ident || b.stall_number || b.stall_code || '')
    rawCodes.split(',').forEach((c) => {
      const clean = c.trim().toUpperCase()
      if (clean) activeBookingsByStall.set(clean, (activeBookingsByStall.get(clean) || 0) + 1)
    })
  })

  // 3. Count active waitlist per stall
  const activeWaitlistByStall = new Map<string, number>()
  marketWaitlistItems.value.forEach((w: any) => {
    const st = String(w.status || '').toLowerCase()
    if (st === 'waiting' || st === 'offered' || st === 'active') {
      const rawCodes = String(w.stall_ident || w.stall_number || w.stall_code || '')
      rawCodes.split(',').forEach((c) => {
        const clean = c.trim().toUpperCase()
        if (clean) activeWaitlistByStall.set(clean, (activeWaitlistByStall.get(clean) || 0) + 1)
      })
    }
  })

  // 4. Iterate all project stalls (using getAllFloorLayoutsForProject correctly with Map)
  visibleProjects.value.forEach((p) => {
    const floorMap = getAllFloorLayoutsForProject(p)
    floorMap.forEach((layout) => {
      ;(layout.stalls || []).forEach((st: any) => {
        const code = String(st.code || st.ident || st.stall_ident || st.stall_number || '').trim().toUpperCase()
        if (!code || result[code]) return

        const apiStall = stallApiMap.get(code)
        const apiWaitlistCount = Number(apiStall?.waitlist_count || apiStall?.WaitlistCount || 0)
        const listWaitlistCount = activeWaitlistByStall.get(code) || 0
        const waitlistCount = Math.max(apiWaitlistCount, listWaitlistCount)

        const bookingFromList = activeBookingsByStall.get(code) || 0
        const isReservedOrOccupied = stallStatusOverrides.value[code] === 'reserved' ||
          st.status === 'reserved' ||
          apiStall?.status === 'reserved' ||
          apiStall?.stall_status === 'reserved'
        const bookingCount = Math.max(bookingFromList, isReservedOrOccupied ? 1 : 0)

        const demandScore = bookingCount + (waitlistCount * 2)

        let demandLevel: 'high' | 'medium' | 'low' | 'none' = 'none'
        if (waitlistCount >= 3 || demandScore >= 7) {
          demandLevel = 'high'
        } else if (waitlistCount >= 1 || demandScore >= 3) {
          demandLevel = 'medium'
        } else if (bookingCount >= 1) {
          demandLevel = 'low'
        } else {
          demandLevel = 'none'
        }

        result[code] = {
          demandLevel,
          demandScore,
          waitlistCount,
          bookingCount
        }
      })
    })
  })

  return result
})

const demandSummary = computed(() => {
  let high = 0
  let medium = 0
  let low = 0
  let none = 0
  let total = 0

  Object.values(stallDemandMap.value).forEach((d) => {
    total++
    if (d.demandLevel === 'high') high++
    else if (d.demandLevel === 'medium') medium++
    else if (d.demandLevel === 'low') low++
    else none++
  })

  return { total, high, medium, low, none }
})

// แผงที่อยู่ระหว่างการฟ้องร้อง/ดำเนินคดี รองรับชื่อฟิลด์จากทั้งข้อมูลสัญญาและข้อมูลแผง
const legalCaseStallCodes = computed<string[]>(() => {
  const codes = new Set<string>()
  const activeLegalStatuses = new Set(['filed', 'pending', 'active', 'litigation', 'in_court', 'court', 'judgment', 'enforcement'])
  const closedLegalStatuses = new Set(['closed', 'dismissed', 'resolved', 'withdrawn', 'cancelled'])
  const collect = (record: any, fallbackCode = '') => {
    const legalStatus = String(record?.legal_status || record?.lawsuit_status || record?.case_status || '').trim().toLowerCase()
    const hasCaseNumber = Boolean(String(record?.court_case_number || record?.case_number || '').trim())
    const explicitlyLegal = record?.is_litigated === true || record?.has_lawsuit === true || record?.legal_case === true
    const isActiveCase = !closedLegalStatuses.has(legalStatus)
      && (activeLegalStatuses.has(legalStatus) || hasCaseNumber || explicitlyLegal)
    if (!isActiveCase) return
    const rawCodes = String(record?.stall_ident || record?.stall_number || record?.stall_code || record?.ident || record?.code || fallbackCode)
    rawCodes.split(',').forEach((rawCode) => {
      const code = rawCode.trim().toUpperCase()
      if (code) codes.add(code)
    })
  }

  marketContracts.value.forEach((record: any) => collect(record))
  marketApplications.value.forEach((record: any) => collect(record))
  stalls.value.forEach((record: any) => collect(record))
  visibleProjects.value.forEach((project) => {
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((record: any) => collect(record))
    })
  })
  return Array.from(codes)
})

type StallOperationalMeta = {
  contractStatus?: 'active' | 'expiring' | 'expired'
  contractExpiryDays?: number
  renewalStatus?: 'pending' | 'reviewing' | 'approved' | 'declined'
  returnStatus?: 'pending' | 'inspection' | 'returned' | 'overdue'
  legalStatus?: 'legal_review' | 'proceeding' | 'enforcement' | 'closed'
  maintenanceStatus?: 'pending' | 'repairing' | 'completed' | 'overdue'
  requestTypes?: Array<'contract_name_change' | 'transfer' | 'contract_cancel' | 'sublease'>
}

const stallOperationalFilterMeta = computed<Record<string, StallOperationalMeta>>(() => {
  const result: Record<string, StallOperationalMeta> = {}
  const ensure = (code: string) => {
    const cleanCode = code.trim().toUpperCase()
    if (!cleanCode) return null
    if (!result[cleanCode]) result[cleanCode] = { requestTypes: [] }
    return result[cleanCode]
  }
  const recordCodes = (record: any) => String(record?.stall_ident || record?.stall_number || record?.stall_code || record?.ident || record?.code || '')
    .split(',').map(code => code.trim().toUpperCase()).filter(Boolean)
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const allContracts = [...marketContracts.value, ...Object.values(customStallContracts.value).flat()]
  allContracts.forEach((contract: any) => {
    const rawStatus = String(contract.contract_state || contract.status || '').toLowerCase()
    const endValue = contract.end_date || contract.contract_end_date
    const endDate = endValue ? new Date(endValue) : null
    const expiryDays = endDate && !Number.isNaN(endDate.getTime())
      ? Math.ceil((endDate.getTime() - today.getTime()) / 86_400_000)
      : undefined
    const contractStatus: StallOperationalMeta['contractStatus'] = ['expired', 'terminated', 'cancelled'].includes(rawStatus) || (expiryDays !== undefined && expiryDays < 0)
      ? 'expired'
      : (expiryDays !== undefined && expiryDays <= 90 ? 'expiring' : 'active')

    const renewalRaw = String(contract.renewal_status || contract.renew_status || '').toLowerCase()
    const renewalMap: Record<string, StallOperationalMeta['renewalStatus']> = {
      pending: 'pending', pending_renewal: 'pending', waiting: 'pending',
      reviewing: 'reviewing', under_review: 'reviewing', submitted: 'reviewing',
      approved: 'approved', renewed: 'approved',
      declined: 'declined', non_renewal: 'declined', not_renewing: 'declined'
    }
    const returnRaw = String(contract.return_status || contract.handover_status || contract.vacate_status || '').toLowerCase()
    const returnMap: Record<string, StallOperationalMeta['returnStatus']> = {
      waiting: 'pending', pending_return: 'pending', vacating: 'pending', notified: 'pending',
      inspection: 'inspection', pending_inspection: 'inspection',
      returned: 'returned', completed: 'returned', accepted: 'returned',
      overdue: 'overdue', late: 'overdue'
    }
    const legalRaw = String(contract.legal_status || contract.lawsuit_status || contract.case_status || '').toLowerCase()
    const legalMap: Record<string, StallOperationalMeta['legalStatus']> = {
      pending_legal: 'legal_review', awaiting_legal: 'legal_review', legal_review: 'legal_review',
      in_progress: 'proceeding', litigation: 'proceeding', court: 'proceeding', filed: 'proceeding', active: 'proceeding',
      enforcement: 'enforcement', execution: 'enforcement', judgment_enforcement: 'enforcement',
      closed: 'closed', resolved: 'closed', dismissed: 'closed', withdrawn: 'closed'
    }

    recordCodes(contract).forEach((code) => {
      const meta = ensure(code)
      if (!meta) return
      meta.contractStatus = contractStatus
      meta.contractExpiryDays = expiryDays
      if (renewalMap[renewalRaw]) meta.renewalStatus = renewalMap[renewalRaw]
      if (returnMap[returnRaw]) meta.returnStatus = returnMap[returnRaw]
      if (legalMap[legalRaw]) meta.legalStatus = legalMap[legalRaw]
    })
  })

  marketApplications.value.forEach((application: any) => {
    const rawType = String(application.request_type || application.application_type || application.type || '').toLowerCase()
    let requestType: NonNullable<StallOperationalMeta['requestTypes']>[number] | undefined
    if (['contract_name_change', 'name_change', 'change_contract_name'].includes(rawType)) requestType = 'contract_name_change'
    else if (['transfer', 'rights_transfer', 'transfer_right'].includes(rawType)) requestType = 'transfer'
    else if (['contract_cancel', 'cancel_contract', 'termination'].includes(rawType)) requestType = 'contract_cancel'
    else if (['sublease', 'subtenant', 'sub_tenant'].includes(rawType)) requestType = 'sublease'
    if (!requestType) return
    recordCodes(application).forEach((code) => {
      const meta = ensure(code)
      if (meta && !meta.requestTypes?.includes(requestType!)) meta.requestTypes?.push(requestType!)
    })
  })

  visibleProjects.value.forEach((project) => {
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        const code = String(stall.code || '').trim().toUpperCase()
        const meta = ensure(code)
        if (!meta || !import.meta.client) return
        const stored = localStorage.getItem(getMaintenanceStorageKey(code))
        if (!stored) return
        try {
          const records = JSON.parse(stored)
          if (!Array.isArray(records) || !records.length) return
          const statuses = records.map((record: any) => {
            const status = String(record.status || '').toLowerCase()
            const dueValue = record.scheduledDate || record.dueDate
            const dueDate = dueValue ? new Date(dueValue) : null
            const overdue = dueDate && !Number.isNaN(dueDate.getTime()) && dueDate < today && !['completed', 'done', 'cancelled', 'rejected'].includes(status)
            if (overdue) return 'overdue'
            if (['in_progress', 'repairing', 'assigned'].includes(status)) return 'repairing'
            if (['completed', 'done'].includes(status)) return 'completed'
            return 'pending'
          }) as NonNullable<StallOperationalMeta['maintenanceStatus']>[]
          meta.maintenanceStatus = statuses.includes('overdue') ? 'overdue'
            : statuses.includes('repairing') ? 'repairing'
              : statuses.includes('pending') ? 'pending' : 'completed'
        } catch {}
      })
    })
  })
  return result
})

// รายชื่อรหัสแผงที่มีสิ่งที่เจ้าหน้าที่ต้องดำเนินการ (ค้างชำระ / ใกล้หมดสัญญา / งานซ่อมค้าง / รออนุมัติจองหรือคำร้อง)
const actionRequiredStallCodes = computed<string[]>(() => {
  const actionSet = new Set<string>()
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const expiringCutoff = new Date(today)
  expiringCutoff.setDate(expiringCutoff.getDate() + 60)

  // 1. การจองที่รอดำเนินการหรืออยู่ระหว่างติดตามผล
  marketBookings.value.forEach((b: any) => {
    const st = String(b.booking_status || b.status || '').toLowerCase()
    if (['pending', 'pending_approval', 'pending_slot', 'approved', 'interested'].includes(st)) {
      const rawCodes = String(b.stall_ident || b.stall_number || b.stall_code || '')
      rawCodes.split(',').forEach((c) => {
        const clean = c.trim().toUpperCase()
        if (clean) actionSet.add(clean)
      })
    }
  })

  // 2. คำร้องขอเช่าที่รอพิจารณาหรือรอทำสัญญา
  marketApplications.value.forEach((a: any) => {
    const st = String(a.application_state || a.status || '').toLowerCase()
    if (st && !['completed', 'rejected', 'cancelled', 'terminated', 'closed'].includes(st)) {
      const rawCodes = String(a.stall_ident || a.stall_number || a.stall_code || '')
      rawCodes.split(',').forEach((c) => {
        const clean = c.trim().toUpperCase()
        if (clean) actionSet.add(clean)
      })
    }
  })

  // 3. สัญญาเช่าที่ค้างชำระหรือใกล้หมดสัญญา
  const allCustomContracts = Object.values(customStallContracts.value).flat()
  ;[...marketContracts.value, ...allCustomContracts].forEach((contract: any) => {
    const contractStatus = String(contract.contract_state || contract.status || '').toLowerCase()
    if (['terminated', 'cancelled', 'expired'].includes(contractStatus)) return
    const rawCodes = String(contract.stall_ident || contract.stall_number || contract.stall_code || '')
    const paymentStatus = String(contract.payment_status || contract.rent_status || contract.billing_status || '').toLowerCase()
    const dueDateValue = contract.payment_due_date || contract.due_date || contract.next_due_date
    const dueDate = dueDateValue ? new Date(dueDateValue) : null
    const isPastDue = dueDate instanceof Date && !Number.isNaN(dueDate.getTime()) && dueDate < today
    const isOverdue = paymentStatus === 'overdue' || paymentStatus === 'past_due' || (isPastDue && paymentStatus !== 'paid')

    const endDateValue = contract.end_date || contract.contract_end_date
    const endDate = endDateValue ? new Date(endDateValue) : null
    const isExpiring = endDate instanceof Date && !Number.isNaN(endDate.getTime()) && endDate >= today && endDate <= expiringCutoff

    if (isOverdue || isExpiring) {
      rawCodes.split(',').forEach((c) => {
        const clean = c.trim().toUpperCase()
        if (clean) actionSet.add(clean)
      })
    }
  })

  // 4. แผงที่มีผู้เช่า (occupied) ซึ่งมีรายการค้างชำระหรืองานซ่อมที่ยังไม่ปิด
  const checkStallHasLocalOrDefaultAction = (upperCode: string) => {
    if (showInspectStallOverlay.value && String(stallEdit.value.ident || '').trim().toUpperCase() === upperCode) {
      return stallActionRichItems.value.length > 0
    }
    if (import.meta.client) {
      let hasOverduePay = true
      const rawPay = localStorage.getItem(getPaymentsStorageKey(upperCode))
      if (rawPay) {
        try {
          const parsedPay = JSON.parse(rawPay)
          hasOverduePay = Array.isArray(parsedPay) && parsedPay.some((p: any) => p.status === 'overdue')
        } catch {}
      }
      let hasOpenMaint = true
      const rawMaint = localStorage.getItem(getMaintenanceStorageKey(upperCode))
      if (rawMaint) {
        try {
          const parsedMaint = JSON.parse(rawMaint)
          hasOpenMaint = Array.isArray(parsedMaint) && parsedMaint.some((m: any) => !['completed', 'rejected'].includes(m.status))
        } catch {}
      }
      return hasOverduePay || hasOpenMaint
    }
    return true
  }

  visibleProjects.value.forEach((project) => {
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim()
        if (!code) return
        const upper = code.toUpperCase()
        const status = resolveRealStallStatus(code, project, stall)
        if (status === 'occupied' && checkStallHasLocalOrDefaultAction(upper)) {
          actionSet.add(upper)
        } else if (status === 'reserved') {
          actionSet.add(upper)
        }
      })
    })
  })

  return Array.from(actionSet)
})

// ข้อมูลตัวเลขสรุปที่แสดงบนแถบฟิลเตอร์ด้านบน: สอดคล้องกับโซนที่เลือก หรือทั้งตลาด
const currentDisplayedSummary = computed(() => {
  const actionSet = new Set(actionRequiredStallCodes.value)
  const legalCaseSet = new Set(legalCaseStallCodes.value)
  if (selectedProject.value) {
    const p = selectedProject.value
    let total = 0
    let available = 0
    let reserved = 0
    let occupied = 0
    let actionRequired = 0
    let legalCases = 0
    const floorMap = getAllFloorLayoutsForProject(p)
    floorMap.forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        total += 1
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim()
        const status = resolveRealStallStatus(code, p, stall)
        if (status === 'reserved') reserved += 1
        else if (status === 'occupied') occupied += 1
        else available += 1
        if (code && actionSet.has(code.toUpperCase())) {
          actionRequired += 1
        }
        if (code && legalCaseSet.has(code.toUpperCase())) legalCases += 1
      })
    })
    return { total, available, reserved, occupied, actionRequired, legalCases }
  }
  let marketActionCount = 0
  let marketLegalCaseCount = 0
  visibleProjects.value.forEach((project) => {
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim().toUpperCase()
        if (code && actionSet.has(code)) {
          marketActionCount += 1
        }
        if (code && legalCaseSet.has(code)) marketLegalCaseCount += 1
      })
    })
  })
  return {
    total: totalMarketStalls.value,
    available: stallStatusSummary.value.available,
    reserved: stallStatusSummary.value.reserved,
    occupied: stallStatusSummary.value.occupied,
    actionRequired: marketActionCount,
    legalCases: marketLegalCaseCount
  }
})

const stallProjectStatusBreakdown = computed(() => {
  const actionCodes = new Set(actionRequiredStallCodes.value.map(code => String(code).trim().toUpperCase()))
  const rows = visibleProjects.value.map((project) => {
    let available = 0
    let reserved = 0
    let occupied = 0
    let actionRequired = 0
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim()
        const status = resolveRealStallStatus(code, project, stall)
        if (status === 'reserved') reserved += 1
        else if (status === 'occupied') occupied += 1
        else available += 1
        if (code && actionCodes.has(code.toUpperCase())) actionRequired += 1
      })
    })
    const total = available + reserved + occupied
    const availablePct = total ? (available / total) * 100 : 0
    const reservedPct = total ? (reserved / total) * 100 : 0
    const occupiedPct = total ? (occupied / total) * 100 : 0
    return {
      id: project.id,
      name: cleanProjectName(project.name),
      order: Number(project.sectionNumber) || Number(String(project.name || '').match(/\d+/)?.[0]) || Number.MAX_SAFE_INTEGER,
      total,
      available,
      reserved,
      occupied,
      actionRequired,
      availablePct,
      reservedPct,
      occupiedPct,
      badgeText: actionRequired > 0 ? `ติดตาม ${actionRequired}` : (availablePct >= 70 ? 'ว่างสูง' : ''),
      badgeTone: actionRequired > 0 ? 'danger' : 'available'
    }
  })
    .sort((a, b) => a.order - b.order || a.name.localeCompare(b.name, 'th', { numeric: true }))

  return rows
})

const dashboardProjectOptions = computed(() => stallProjectStatusBreakdown.value.map(project => ({
  value: project.id,
  label: `${project.name} · ${project.total.toLocaleString('th-TH')} แผง`
})))
const selectDashboardProject = async (projectId: string) => {
  const id = String(projectId || '')
  const project = visibleProjects.value.find(item => item.id === id)
  if (!project) return
  selectedDashboardProjectId.value = id
  if (selectedProject.value?.id === id) return
  await nextTick()
  if (typeof mapRef.value?.enterProject === 'function') {
    mapRef.value.enterProject(project)
  } else {
    selectProject(project)
  }
}
const handleDashboardProjectChange = (value: string | number) => {
  void selectDashboardProject(String(value))
}
const selectedDashboardProjectStatus = computed(() => stallProjectStatusBreakdown.value.find(
  project => project.id === selectedDashboardProjectId.value
) || stallProjectStatusBreakdown.value[0] || null)
const projectDonutStyle = (project: { total: number; availablePct: number; occupiedPct: number }) => {
  if (!project || project.total === 0) return { background: '#e2e8f0' }
  const availableEnd = project.availablePct
  const occupiedEnd = availableEnd + project.occupiedPct
  return {
    background: `conic-gradient(#10b981 0 ${availableEnd}%, #64748b ${availableEnd}% ${occupiedEnd}%, #f59e0b ${occupiedEnd}% 100%)`
  }
}
const selectedDashboardProjectDonutStyle = computed(() => {
  const project = selectedDashboardProjectStatus.value
  return project ? projectDonutStyle(project) : { background: '#e2e8f0' }
})

watch(
  [() => selectedProject.value?.id, stallProjectStatusBreakdown],
  ([mapProjectId, projects]) => {
    if (mapProjectId && projects.some(project => project.id === mapProjectId)) {
      selectedDashboardProjectId.value = mapProjectId
      return
    }
    if (!projects.some(project => project.id === selectedDashboardProjectId.value)) {
      selectedDashboardProjectId.value = projects[0]?.id || ''
    }
  },
  { immediate: true }
)

const stallContractRenewalSummary = computed(() => {
  const statuses = dashboardContracts.value.map((contract: any) => String(
    contract.renewal_status || contract.renew_status || ''
  ).toLowerCase())
  const submitted = statuses.filter(status => ['approved', 'renewed', 'completed', 'submitted'].includes(status)).length
  const declined = statuses.filter(status => ['declined', 'non_renewal', 'not_renewing'].includes(status)).length
  const pending = statuses.filter(status => ['pending', 'pending_renewal', 'reviewing', 'under_review', 'waiting'].includes(status)).length
  return { submitted, pending, declined }
})

const stallReturnStatusSummary = computed(() => {
  const summary = { notified: 0, waiting: 0, inspection: 0, returned: 0, overdue: 0 }
  dashboardContracts.value.forEach((contract: any) => {
    const status = String(contract.return_status || contract.handover_status || contract.vacate_status || '').toLowerCase()
    if (['notified', 'notice', 'non_renewal'].includes(status)) summary.notified += 1
    else if (['waiting', 'pending_return', 'vacating'].includes(status)) summary.waiting += 1
    else if (['inspection', 'pending_inspection'].includes(status)) summary.inspection += 1
    else if (['returned', 'completed', 'accepted'].includes(status)) summary.returned += 1
    else if (['overdue', 'late'].includes(status)) summary.overdue += 1
  })
  return summary
})

const stallLegalStatusSummary = computed(() => {
  const summary = { awaiting: 0, proceeding: 0, enforcement: 0, closed: 0 }
  const collect = (record: any) => {
    const status = String(record?.legal_status || record?.lawsuit_status || record?.case_status || '').toLowerCase()
    if (['pending_legal', 'awaiting_legal', 'legal_review'].includes(status)) summary.awaiting += 1
    else if (['in_progress', 'litigation', 'court', 'filed'].includes(status)) summary.proceeding += 1
    else if (['enforcement', 'execution', 'judgment_enforcement'].includes(status)) summary.enforcement += 1
    else if (['closed', 'resolved', 'dismissed', 'withdrawn'].includes(status)) summary.closed += 1
  }
  dashboardContracts.value.forEach(collect)
  marketApplications.value.filter(recordMatchesDashboardProject).forEach(collect)
  stalls.value.filter(recordMatchesDashboardProject).forEach(collect)
  return summary
})

const stallMaintenanceDashboardSummary = computed(() => {
  const summary = { pending: 0, repairing: 0, completed: 0, overdue: 0 }
  const visited = new Set<string>()
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  visibleProjects.value
    .filter(project => !selectedDashboardProjectId.value || project.id === selectedDashboardProjectId.value)
    .forEach((project) => {
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim().toUpperCase()
        if (!code || visited.has(code) || !import.meta.client) return
        visited.add(code)
        const stored = localStorage.getItem(`stall_maintenance_${marketId.value}_${code}`)
        if (!stored) return
        try {
          const records = JSON.parse(stored)
          if (!Array.isArray(records)) return
          records.forEach((record: any) => {
            const status = String(record.status || '').toLowerCase()
            const dueValue = record.scheduledDate || record.dueDate
            const dueDate = dueValue ? new Date(dueValue) : null
            const isOverdue = dueDate instanceof Date && !Number.isNaN(dueDate.getTime()) && dueDate < today
              && !['completed', 'rejected', 'cancelled'].includes(status)
            if (isOverdue) summary.overdue += 1
            else if (['in_progress', 'repairing', 'assigned'].includes(status)) summary.repairing += 1
            else if (['completed', 'done'].includes(status)) summary.completed += 1
            else if (!['rejected', 'cancelled'].includes(status)) summary.pending += 1
          })
        } catch {}
      })
    })
    })
  return summary
})

const marketMapRiskSummary = computed(() => {
  const overdueStalls = new Set<string>()
  const expiringStalls = new Set<string>()
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const expiringCutoff = new Date(today)
  expiringCutoff.setDate(expiringCutoff.getDate() + 30)

  marketContracts.value.forEach((contract: any) => {
    const contractStatus = String(contract.contract_state || contract.status || '').toLowerCase()
    if (['terminated', 'cancelled', 'expired'].includes(contractStatus)) return

    const stallCode = String(contract.stall_ident || contract.stall_number || contract.stall_code || contract.id || '').trim().toUpperCase()
    if (!stallCode) return

    const paymentStatus = String(contract.payment_status || contract.rent_status || contract.billing_status || '').toLowerCase()
    const dueDateValue = contract.payment_due_date || contract.due_date || contract.next_due_date
    const dueDate = dueDateValue ? new Date(dueDateValue) : null
    const isPastDue = dueDate instanceof Date && !Number.isNaN(dueDate.getTime()) && dueDate < today
    if (paymentStatus === 'overdue' || paymentStatus === 'past_due' || (isPastDue && paymentStatus !== 'paid')) {
      overdueStalls.add(stallCode)
    }

    const endDateValue = contract.end_date || contract.contract_end_date
    const endDate = endDateValue ? new Date(endDateValue) : null
    if (endDate instanceof Date && !Number.isNaN(endDate.getTime()) && endDate >= today && endDate <= expiringCutoff) {
      expiringStalls.add(stallCode)
    }
  })

  return { overdue: overdueStalls.size, expiring: expiringStalls.size }
})

const currentDisplayedDemandSummary = computed(() => {
  if (selectedProject.value) {
    const p = selectedProject.value
    let total = 0
    let high = 0
    let medium = 0
    let low = 0
    let none = 0
    const floorMap = getAllFloorLayoutsForProject(p)
    floorMap.forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        total += 1
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim().toUpperCase()
        const demand = stallDemandMap.value[code]
        const level = demand?.demandLevel || 'none'
        if (level === 'high') high += 1
        else if (level === 'medium') medium += 1
        else if (level === 'low') low += 1
        else none += 1
      })
    })
    return { total, high, medium, low, none }
  }
  return demandSummary.value
})

const pointInPolygonRing = ([lng, lat]: [number, number], ring: [number, number][] = []): boolean => {
  let inside = false
  for (let i = 0, j = ring.length - 1; i < ring.length; j = i++) {
    const [xi, yi] = ring[i] || [0, 0]
    const [xj, yj] = ring[j] || [0, 0]
    if ((yi > lat) !== (yj > lat) && lng < ((xj - xi) * (lat - yi)) / (yj - yi) + xi) inside = !inside
  }
  return inside
}

const syncSnapshotsFromCurrent = () => {
  savedProjectsSnapshot.value = cloneProjects(projects.value)
  savedFacilitiesSnapshot.value = cloneProjects(facilities.value)
  savedStallsSnapshot.value = cloneProjects(stalls.value)
  pendingStallUpserts.value = []
  pendingStallDeleteIds.value = []
  hasUnsavedChanges.value = false
}

const revertUnsavedChanges = () => {
  mapRef.value?.finishActiveTools?.()
  projects.value = cloneProjects(savedProjectsSnapshot.value)
  facilities.value = cloneProjects(savedFacilitiesSnapshot.value)
  stalls.value = cloneProjects(savedStallsSnapshot.value)
  if (selectedProject.value) {
    selectedProject.value = projects.value.find(p => p.id === selectedProject.value?.id) || null
  }
  pendingStallUpserts.value = []
  pendingStallDeleteIds.value = []
  hasUnsavedChanges.value = false
}

const queueStallUpsert = (payload: any) => {
  const code = String(payload.ident || payload.stall_ident || '').trim().toUpperCase()
  if (!code) return
  const idx = pendingStallUpserts.value.findIndex(item => String(item.ident || item.stall_ident || '').trim().toUpperCase() === code)
  if (idx >= 0) pendingStallUpserts.value[idx] = { ...pendingStallUpserts.value[idx], ...payload }
  else pendingStallUpserts.value.push({ ...payload })
}

const persistGISChanges = async (publishState = true) => {
  mapRef.value?.finishActiveTools?.()
  const deleteIds = [...new Set(pendingStallDeleteIds.value.filter(id => Number(id) > 0))]
  const upserts = [...pendingStallUpserts.value]

  for (const id of deleteIds) {
    await del(`/stalls/${id}`).catch(() => {})
  }
  for (const item of upserts) {
    try {
      const response = item.id
        ? await put(`/stalls/${item.id}`, stallPayload(item))
        : await post('/stalls', stallPayload(item))
      const saved = unwrap(response)
      if (saved?.id) {
        const sIdx = stalls.value.findIndex(s => String(s.ident || s.stall_ident).toUpperCase() === String(saved.ident || item.ident).toUpperCase())
        if (sIdx >= 0) stalls.value[sIdx] = { ...stalls.value[sIdx], ...saved }
      }
    } catch {}
  }

  await put(`/markets/${encodeURIComponent(marketId.value)}/gis-config`, {
    projects: projects.value,
    facilities: facilities.value,
    published: publishState
  })
  clearMarketGisCache(marketId.value)
  gisPublished.value = publishState
  syncSnapshotsFromCurrent()
}

const askUnsavedLeaveDecision = (): Promise<'save' | 'discard' | 'cancel'> => {
  showUnsavedLeaveModal.value = true
  return new Promise((resolve) => {
    unsavedLeaveResolver = resolve
  })
}

const resolveUnsavedLeave = async (choice: 'save' | 'discard' | 'cancel') => {
  if (choice === 'save') {
    savingGIS.value = true
    try {
      await persistGISChanges(true)
      showUnsavedLeaveModal.value = false
      void showSuccess('บันทึกการแก้ไขและอัปเดตผังตลาดเรียบร้อยแล้ว', 'บันทึกสำเร็จ')
      if (unsavedLeaveResolver) {
        unsavedLeaveResolver('save')
        unsavedLeaveResolver = null
      }
    } catch {
      void showError('ไม่สามารถบันทึกผังได้ กรุณาตรวจสอบสิทธิ์และข้อมูลพิกัด', 'บันทึกไม่สำเร็จ')
    } finally {
      savingGIS.value = false
    }
    return
  }
  showUnsavedLeaveModal.value = false
  if (choice === 'discard') {
    revertUnsavedChanges()
  }
  if (unsavedLeaveResolver) {
    unsavedLeaveResolver(choice)
    unsavedLeaveResolver = null
  }
}

const ensureCanLeaveOrBack = async (): Promise<boolean> => {
  if (isInspectMode.value || !hasUnsavedChanges.value) return true
  const decision = await askUnsavedLeaveDecision()
  return decision !== 'cancel'
}

const selectProject = (project: ProjectZone) => {
  const switchingProject = Boolean(
    selectedProject.value?.id
    && selectedProject.value.id !== project.id
    && showInspectStallOverlay.value
  )
  selectedProject.value = project
  // ถ้าเลือกโครงการใหม่ขณะเปิดรายละเอียดแผง ให้ยุบ drawer เดิมกลับเป็น
  // panel โครงการก่อน เพื่อไม่ให้ข้อมูลแผงของโครงการเก่าค้างซ้อนกับโครงการใหม่
  if (switchingProject) closeStall()
}
const clearSelectedProject = () => {
  selectedProject.value = null
  closeStall()
  if (route.query.project || route.query.stall) {
    const nextQuery: Record<string, any> = { ...route.query }
    delete nextQuery.project
    delete nextQuery.stall
    delete nextQuery.drawer
    delete nextQuery.tab
    void router.replace({ query: nextQuery })
  }
}
const handleRequestExitProject = () => {
  clearSelectedProject()
}
const handleBreadcrumbBack = async () => {
  if (selectedProject.value) {
    clearSelectedProject()
    return
  }
  await navigateTo(`/stalls?market=${encodeURIComponent(marketId.value)}`)
}
const modeLabel = computed(() => isInspectMode.value ? 'ดูข้อมูลแผงเช่า' : 'จัดการแผงเช่า')
const breadcrumbItems = computed(() => selectedProject.value
  ? [{ label: modeLabel.value, onClick: handleRequestExitProject }]
  : [])
const breadcrumbCurrent = computed(() => selectedProject.value
  ? selectedProject.value.name
  : modeLabel.value)
const recomputeDirtyState = () => {
  const projectsChanged = JSON.stringify(projects.value) !== JSON.stringify(savedProjectsSnapshot.value)
  const facilitiesChanged = JSON.stringify(facilities.value) !== JSON.stringify(savedFacilitiesSnapshot.value)
  const stallsChanged = JSON.stringify(stalls.value) !== JSON.stringify(savedStallsSnapshot.value)
  hasUnsavedChanges.value = Boolean(projectsChanged || facilitiesChanged || stallsChanged || pendingStallDeleteIds.value.length > 0)
}

const updateProjectPolygon = async (project: ProjectZone) => {
  const index = projects.value.findIndex(item => item.id === project.id)
  if (index >= 0) projects.value[index] = { ...projects.value[index], ...project }
  if (selectedProject.value?.id === project.id) {
    selectedProject.value = projects.value[index] || project
  }
  recomputeDirtyState()
}
const updateStallPolygon = async (project: ProjectZone, updatedStall: any) => {
  const index = projects.value.findIndex(item => item.id === project.id)
  if (index >= 0) {
    projects.value[index] = { ...projects.value[index], ...project }
    if (selectedProject.value?.id === project.id) {
      selectedProject.value = projects.value[index] || project
    }
  }
  const code = String(updatedStall?.code || updatedStall?.ident || '')
  const stored = stalls.value.find(s => String(s.ident || s.stall_ident).toUpperCase() === code.toUpperCase())
  if (stored && Array.isArray(updatedStall?.coordinates)) {
    const updatedRecord = {
      ...stored,
      width: updatedStall.widthM || stored.width || 2,
      depth: updatedStall.depthM || stored.depth || 2,
      coordinates_geojson: JSON.stringify({ type: 'Polygon', coordinates: [updatedStall.coordinates] })
    }
    const sIdx = stalls.value.findIndex(s => String(s.ident || s.stall_ident).toUpperCase() === code.toUpperCase())
    if (sIdx >= 0) stalls.value[sIdx] = updatedRecord
    queueStallUpsert(updatedRecord)
  }
  recomputeDirtyState()
}
const openZoneEditor = (coordinates: [number, number][]) => {
  drawnZoneCoordinates.value = coordinates
  zoneForm.value = { name: '', sectionNumber: '', category: marketCategoryOptions.value[0]?.value || '', floorCount: 1 }
  zoneErrors.name = ''
  zoneErrors.sectionNumber = ''
  zoneErrors.floorCount = ''
  zoneErrors.category = ''
  showZoneEditor.value = true
}
const dismissZoneEditor = () => {
  showZoneEditor.value = false
  drawnZoneCoordinates.value = []
  zoneErrors.name = ''
  zoneErrors.sectionNumber = ''
  zoneErrors.floorCount = ''
  zoneErrors.category = ''
}
const zoneErrors = reactive({ name: '', sectionNumber: '', floorCount: '', category: '' })
const validateZoneField = (field: 'name' | 'sectionNumber' | 'floorCount' | 'category') => {
  if (field === 'name') {
    zoneErrors.name = String(zoneForm.value.name || '').trim() ? '' : 'กรุณาระบุชื่อโซน'
  } else if (field === 'sectionNumber') {
    zoneErrors.sectionNumber = String(zoneForm.value.sectionNumber || '').trim() ? '' : 'กรุณาระบุรหัสโซน'
  } else if (field === 'floorCount') {
    const count = Number(zoneForm.value.floorCount)
    zoneErrors.floorCount = (!Number.isFinite(count) || count < 1 || count > 20) ? 'กรุณาระบุจำนวนชั้น 1 - 20' : ''
  } else if (field === 'category') {
    zoneErrors.category = String(zoneForm.value.category || '').trim() ? '' : 'กรุณาเลือกหมวดสินค้า'
  }
}
const validateZoneForm = () => {
  validateZoneField('name')
  validateZoneField('sectionNumber')
  validateZoneField('floorCount')
  validateZoneField('category')
  return !zoneErrors.name && !zoneErrors.sectionNumber && !zoneErrors.floorCount && !zoneErrors.category
}
const calcPolygonAreaSqm = (ring: [number, number][]) => {
  if (ring.length < 3) return 0
  const latRad = ((ring[0]?.[1] || 13.8) * Math.PI) / 180
  const mPerDegLng = 111320 * Math.cos(latRad)
  const mPerDegLat = 110540
  let twiceArea = 0
  for (let i = 0; i < ring.length - 1; i++) {
    const x1 = (ring[i]?.[0] || 0) * mPerDegLng
    const y1 = (ring[i]?.[1] || 0) * mPerDegLat
    const x2 = (ring[i + 1]?.[0] || 0) * mPerDegLng
    const y2 = (ring[i + 1]?.[1] || 0) * mPerDegLat
    twiceArea += x1 * y2 - x2 * y1
  }
  return Math.max(1, Math.round(Math.abs(twiceArea) / 2))
}
const ZONE_COLOR_PALETTE = ['#3b82f6', '#e06456', '#4eb5b7', '#7e8ee8', '#10b981', '#f59e0b', '#8b5cf6', '#ec4899']
const resolveCategoryColor = (categoryName: string) => {
  const existing = projects.value.find(p => p.category === categoryName && p.color)
    || getMarketGisConfig(marketId.value).projects.find(p => p.category === categoryName && p.color)
  if (existing?.color) return existing.color
  return ZONE_COLOR_PALETTE[projects.value.length % ZONE_COLOR_PALETTE.length] || '#3b82f6'
}
const createZone = async () => {
  if (!validateZoneForm()) return
  const ring = drawnZoneCoordinates.value.map(point => [...point] as [number, number])
  if (ring.length < 3) return
  ring.push([...ring[0]!])
  const centroid: [number, number] = [ring.slice(0, -1).reduce((sum, point) => sum + point[0], 0) / (ring.length - 1), ring.slice(0, -1).reduce((sum, point) => sum + point[1], 0) / (ring.length - 1)]
  const id = `zone-${Date.now()}`
  const floorCount = Math.min(20, Math.max(1, Number(zoneForm.value.floorCount) || 1))
  const floors = Array.from({ length: floorCount }, (_, index) => ({ id: String(index + 1), name: `ชั้น ${index + 1}`, shortLabel: `${index + 1}F`, level: index + 1 }))
  const chosenCategory = String(zoneForm.value.category || '').trim() || marketCategoryOptions.value[0]?.value || 'ทั่วไป'
  const project: ProjectZone = { id, name: String(zoneForm.value.name).trim(), sectionNumber: String(zoneForm.value.sectionNumber).trim().toUpperCase(), category: chosenCategory, color: resolveCategoryColor(chosenCategory), stallsAvailable: 0, totalStalls: 0, realAreaSqm: calcPolygonAreaSqm(ring), centroid, coordinates: ring, floors, stalls: [], layoutMode: 'manual' }
  projects.value = [...projects.value, project]
  selectedProject.value = project
  hasUnsavedChanges.value = true
  dismissZoneEditor()
  void showSuccess('โซนถูกเพิ่มลงในผังแล้ว กดปุ่ม “บันทึกการแก้ไข” ด้านบนเพื่อบันทึกผัง', 'เพิ่มโซนในผังแล้ว')
}
const openFacilityEditor = (coordinates: [number, number], projectId?: string) => {
  const hostProject = projectId
    ? projects.value.find(p => p.id === projectId)
    : (selectedProject.value || projects.value.find(p => pointInPolygonRing(coordinates, p.coordinates)))
  facilityForm.value = { name: '', type: 'toilet', coordinates, projectId: hostProject?.id }
  facilityNameError.value = ''
  showFacilityEditor.value = true
}
const closeFacilityEditor = () => {
  showFacilityEditor.value = false
  facilityNameError.value = ''
}
const validateFacilityName = () => {
  facilityNameError.value = facilityForm.value.name.trim() ? '' : 'กรุณาระบุชื่อจุดบริการ'
  return !facilityNameError.value
}
const createFacility = async () => {
  if (!validateFacilityName()) return
  facilities.value = [
    ...facilities.value,
    {
      id: `facility-${Date.now()}`,
      name: facilityForm.value.name,
      type: facilityForm.value.type,
      coordinates: facilityForm.value.coordinates,
      projectId: facilityForm.value.projectId
    }
  ]
  hasUnsavedChanges.value = true
  closeFacilityEditor()
  void showSuccess('เพิ่มจุดบริการลงในผังแล้ว กดปุ่ม “บันทึกการแก้ไข” ด้านบนเพื่อบันทึกผัง', 'เพิ่มจุดบริการในผังแล้ว')
}
const deleteFacility = async (facility: MarketFacility) => {
  const confirmed = await showConfirm(`ต้องการลบจุดบริการ “${facility.name}” ออกจากผังตลาดใช่หรือไม่`, { title: 'ยืนยันการลบจุดบริการ', type: 'error', icon: 'trash', confirmText: 'ลบจุดบริการ' })
  if (!confirmed) return
  facilities.value = facilities.value.filter(item => item.id !== facility.id)
  hasUnsavedChanges.value = true
  void showSuccess(`ลบจุดบริการ “${facility.name}” ออกจากผังแล้ว กดปุ่ม “บันทึกการแก้ไข” ด้านบนเพื่อบันทึกผัง`, 'ลบจุดบริการในผังแล้ว')
}
const stallErrors = reactive({ ident: '', width: '', depth: '', monthly_rent: '' })
const resetStallErrors = () => { stallErrors.ident = ''; stallErrors.width = ''; stallErrors.depth = ''; stallErrors.monthly_rent = '' }
const validateStallField = (field: keyof typeof stallErrors) => {
  if (field === 'ident') {
    stallErrors.ident = String(stallEdit.value.ident || '').trim() ? '' : 'กรุณาระบุรหัสแผง'
  } else if (field === 'width') {
    const val = Number(stallEdit.value.width)
    stallErrors.width = !Number.isFinite(val) || val <= 0 ? 'กรุณาระบุความกว้างมากกว่า 0 เมตร' : ''
  } else if (field === 'depth') {
    const val = Number(stallEdit.value.depth)
    stallErrors.depth = !Number.isFinite(val) || val <= 0 ? 'กรุณาระบุความลึกมากกว่า 0 เมตร' : ''
  } else if (field === 'monthly_rent') {
    const val = Number(stallEdit.value.monthly_rent)
    stallErrors.monthly_rent = !Number.isFinite(val) || val < 0 || String(stallEdit.value.monthly_rent ?? '').trim() === '' ? 'กรุณาระบุค่าเช่ารายเดือน' : ''
  }
}
const validateStallForm = () => {
  validateStallField('ident')
  validateStallField('width')
  validateStallField('depth')
  validateStallField('monthly_rent')
  return !stallErrors.ident && !stallErrors.width && !stallErrors.depth && !stallErrors.monthly_rent
}
const openDrawnStall = (project: ProjectZone, coordinates: [number, number][], floorId: string, dimensions?: { widthM: number; depthM: number }) => { const closed = [...coordinates.map(point => [...point] as [number, number]), [...coordinates[0]!]]; const floor = project.floors?.find(item => item.id === floorId) || project.floors?.[0] || { id: '1', name: 'ชั้น 1', level: 1 }; const count = (project.stalls || []).filter((stall: any) => String(stall.floor || '1') === String(floor.id)).length + 1; pendingStallLayout.value = { projectId: project.id, coordinates: closed, floorId: String(floor.id) }; stallEdit.value = { ...emptyStall(), ident: `${project.sectionNumber}-${floor.id}-${String(count).padStart(3, '0')}`, project_name: project.name, category: project.category || marketCategoryOptions.value[0]?.value || '', floor: String(floor.id), coordinates_geojson: JSON.stringify({ type: 'Polygon', coordinates: [closed] }), width: dimensions?.widthM || 2, depth: dimensions?.depthM || 2 }; resetStallErrors(); stallDrawerMode.value = 'edit'; showInspectStallOverlay.value = true; showStallEditor.value = false }
const deleteProject = async (project: ProjectZone, stats?: { total: number; available: number; reserved: number; occupied: number }) => {
  const gisStalls = Array.isArray(project.stalls) ? project.stalls : []
  const dbStallsInZone = stalls.value.filter(s => String(s.project_name || '').trim() === String(project.name || '').trim())
  const nonAvailableFromGis = gisStalls.filter((s: any) => {
    const code = String(s.code || s.ident || '')
    const effStatus = stallStatusOverrides.value[code] || s.status || 'available'
    return effStatus !== 'available'
  }).length
  const nonAvailableFromDb = dbStallsInZone.filter((s: any) => {
    const effStatus = s.stall_status || s.status || 'available'
    return effStatus !== 'available'
  }).length
  const nonAvailableFromStats = (stats?.reserved || 0) + (stats?.occupied || 0)
  const blockedStallCount = Math.max(nonAvailableFromGis, nonAvailableFromDb, nonAvailableFromStats)

  if (blockedStallCount > 0) {
    await showError(
      `ไม่สามารถลบโซน “${project.name}” ได้ เนื่องจากมีแผงที่ติดจองหรือมีผู้เช่าอยู่ ${blockedStallCount} แผง (ต้องไม่มีแผงในโซน หรือแผงทั้งหมดต้องเป็นสถานะว่างเท่านั้นจึงจะลบโซนได้)`,
      'ไม่สามารถลบโซนได้'
    )
    return
  }

  const totalStallsCount = Math.max(gisStalls.length, dbStallsInZone.length, stats?.total || 0)
  const confirmMessage = totalStallsCount > 0
    ? `โซน “${project.name}” มีแผงสถานะว่าง ${totalStallsCount} แผง (ไม่มีผู้เช่าหรือติดจอง) ต้องการลบโซนนี้พร้อมแผงและจุดบริการภายในโซนออกจากผังตลาดใช่หรือไม่`
    : `ต้องการลบโซน “${project.name}” และจุดบริการภายในโซนออกจากผังตลาดใช่หรือไม่`

  const confirmed = await showConfirm(confirmMessage, { title: 'ยืนยันการลบโซน', type: 'error', icon: 'trash', confirmText: 'ลบโซน' })
  if (!confirmed) return

  const dbStallsToDelete = [...dbStallsInZone]
  dbStallsToDelete.forEach(s => {
    if (s.id) pendingStallDeleteIds.value.push(Number(s.id))
  })
  projects.value = projects.value.filter(item => item.id !== project.id)
  facilities.value = facilities.value.filter(f => (
    f.projectId !== project.id && !pointInPolygonRing(f.coordinates, project.coordinates)
  ))
  stalls.value = stalls.value.filter(s => String(s.project_name || '').trim() !== String(project.name || '').trim())
  selectedProject.value = null
  hasUnsavedChanges.value = true
  void showSuccess(`ลบโซน “${project.name}” ออกจากผังแล้ว กดปุ่ม “บันทึกการแก้ไข” ด้านบนเพื่อบันทึกผัง`, 'ลบโซนในผังแล้ว')
}
const stallPayload = (form: any) => ({ ...form, market_id: resolvedMarketID.value, market_slug: marketId.value, market_name: marketName.value, zone_id: Number(form.zone_id), width: Number(form.width), depth: Number(form.depth), daily_rent: Number(form.daily_rent), monthly_rent: Number(form.monthly_rent) })
const finishCloseStall = (preserveProjectMorph = false) => {
  clearStallDrawerResizeTimer()
  clearStallDrawerSwapTimer()
  fullDetailSettling.value = false
  stallDrawerExpanded.value = false
  isStallDrawerClosing.value = false
  if (!preserveProjectMorph) {
    isProjectPanelMorphing.value = false
    isProjectPanelReturning.value = false
  }
  isProjectMorphSwap.value = false
  showStallEditor.value = false
  showInspectStallOverlay.value = false
  showInspectPhotoModal.value = false
  resetStallErrors()
  stallDrawerMode.value = 'detail'
  stallDetailTab.value = 'info'
  pendingStallLayout.value = null
}
const closeStall = () => {
  clearStallDrawerResizeTimer()
  const shouldCollapseBeforeClose = import.meta.client
    && isMapOverviewMode.value
    && showInspectStallOverlay.value
    && stallDrawerMode.value === 'full-detail'

  if (shouldCollapseBeforeClose) {
    // ซ่อนเฉพาะเนื้อหา แล้วส่งกรอบต่อให้ panel โครงการยุบกลับจาก 720px → 380px
    isStallDrawerClosing.value = true
    stallDrawerResizeTimer = window.setTimeout(() => {
      stallDrawerResizeTimer = null
      isProjectPanelMorphing.value = true
      isProjectPanelReturning.value = true
      finishCloseStall(true)
      requestAnimationFrame(() => {
        requestAnimationFrame(() => {
          isProjectPanelMorphing.value = false
          stallDrawerResizeTimer = window.setTimeout(() => {
            isProjectPanelReturning.value = false
            stallDrawerResizeTimer = null
          }, STALL_DRAWER_RETURN_DURATION)
        })
      })
    }, 120)
    return
  }
  finishCloseStall()
}
const openStallEditorFromOverlay = () => {
  resetStallErrors()
  stallDrawerMode.value = 'edit'
}
const cancelStallEdit = () => {
  if (isCreatingNewStall.value) {
    pendingStallLayout.value = null
    closeStall()
  } else {
    const code = String(stallEdit.value.ident || '')
    const stored = stalls.value.find(s => String(s.ident || s.stall_ident).toUpperCase() === code.toUpperCase())
    if (stored) {
      stallEdit.value = {
        ...emptyStall(),
        ...stored,
        stall_status: resolveRealStallStatus(code, selectedProject.value, stored)
      }
    }
    resetStallErrors()
    stallDrawerMode.value = isInspectMode.value ? 'full-detail' : 'detail'
    stallDrawerExpanded.value = stallDrawerMode.value === 'full-detail'
  }
}
const selectStall = async (mapStall: any) => {
  if (isMapOverviewMode.value && pageMode.value === 'costs') {
    handleCostStallClick(mapStall)
    return
  }
  const replaceExistingStall = showInspectStallOverlay.value && stallDrawerMode.value === 'full-detail'
  pendingStallLayout.value = null
  clearStallDrawerSwapTimer()
  isStallDrawerClosing.value = false
  isProjectPanelReturning.value = false
  resetStallErrors()
  isMiniActionAlertExpanded.value = false
  const code = String(mapStall.code || mapStall.ident || '')
  const stored = stalls.value.find(s => String(s.ident || s.stall_ident).toUpperCase() === code.toUpperCase())
  const geometry = mapStall.coordinates?.length ? JSON.stringify({ type: 'Polygon', coordinates: [mapStall.coordinates] }) : ''
  contracts.value = []
  stallHistory.value = []
  stallDetailTab.value = 'info'
  if (!stored?.id) {
    const rawMapStatus = String(mapStall.status || mapStall.stall_status || 'available').toLowerCase()
    if (rawMapStatus === 'available') {
      const upperCode = code.toUpperCase()
      const generatedMockId = `cnt-${upperCode}-01`
      const existingCustom = customStallContracts.value[upperCode] || []
      const cleanedCustom = existingCustom.filter((contract: any) => String(contract.id) !== generatedMockId)
      if (cleanedCustom.length !== existingCustom.length) {
        if (cleanedCustom.length) customStallContracts.value[upperCode] = cleanedCustom
        else delete customStallContracts.value[upperCode]
      }
    }
    stallEdit.value = {
      ...emptyStall(),
      ident: code,
      name: mapStall.name || '',
      category: mapStall.areaType || selectedProject.value?.category || '',
      project_name: selectedProject.value?.name || mapStall.zone || '',
      floor: String(mapStall.floor || '1'),
      width: mapStall.widthM || 2,
      depth: mapStall.depthM || 2,
      monthly_rent: mapStall.monthlyRent || 1800,
      stall_status: resolveRealStallStatus(code, selectedProject.value, mapStall),
      coordinates_geojson: geometry
    }
    loadStallContracts(code)
    loadStallSubtenants(code)
    loadStallMaintenance(code)
    loadStallPayments(code)
    lastLoadedExtraStallCode.value = code.toUpperCase()
    stallDrawerMode.value = isInspectMode.value ? 'full-detail' : 'detail'
    stallDrawerExpanded.value = false
    showSelectedStallPanel(replaceExistingStall)
    if (isInspectMode.value) {
      void loadStallExtraDetails(false)
    }
    return
  }
  stallEdit.value = {
    ...emptyStall(),
    ...stored,
    coordinates_geojson: stored.coordinates_geojson || geometry,
    stall_status: resolveRealStallStatus(code, selectedProject.value, stored)
  }
  loadStallContracts(code)
  loadStallSubtenants(code)
  loadStallMaintenance(code)
  loadStallPayments(code)
  lastLoadedExtraStallCode.value = code.toUpperCase()
  stallDrawerMode.value = isInspectMode.value ? 'full-detail' : 'detail'
  stallDrawerExpanded.value = false
  showSelectedStallPanel(replaceExistingStall)
  if (isInspectMode.value) {
    void loadStallExtraDetails(false)
  }
  get(`/stalls/${stored.id}/contracts`)
    .then((res) => {
      contracts.value = Array.isArray(unwrap(res)) ? unwrap(res) : []
    })
    .catch(() => {})
}
const deleteStall = async (projectArg?: ProjectZone | null, mapStallArg?: any) => {
  const targetProject = projectArg || selectedProject.value || projects.value.find(p => String(p.name || '').trim() === String(stallEdit.value.project_name || '').trim())
  const code = String(mapStallArg?.code || mapStallArg?.ident || stallEdit.value.ident || '').trim()
  if (!code) return

  const stored = stalls.value.find(s => String(s.ident || s.stall_ident).toUpperCase() === code.toUpperCase())
  const gisStall = targetProject?.stalls?.find((s: any) => String(s.code || s.ident).toUpperCase() === code.toUpperCase())
  const effectiveStatus = stallStatusOverrides.value[code] || stallStatusOverrides.value[code.toUpperCase()] || stored?.stall_status || stored?.status || gisStall?.status || mapStallArg?.status || 'available'

  if (effectiveStatus !== 'available') {
    await showError(
      `ไม่สามารถลบแผง “${code}” ได้ เนื่องจากแผงนี้ติดจองหรือมีผู้เช่าอยู่ (ต้องเป็นสถานะว่างเท่านั้นจึงจะลบแผงได้)`,
      'ไม่สามารถลบแผงได้'
    )
    return
  }

  const confirmed = await showConfirm(
    `ต้องการลบแผง “${code}” ออกจากโซน “${targetProject?.name || stallEdit.value.project_name || '-'}” ใช่หรือไม่`,
    { title: 'ยืนยันการลบแผง', type: 'error', icon: 'trash', confirmText: 'ลบแผง' }
  )
  if (!confirmed) return

  if (targetProject) {
    const pIdx = projects.value.findIndex(p => p.id === targetProject.id)
    if (pIdx >= 0) {
      const sourceStalls = Array.isArray(projectArg?.stalls) && projectArg.stalls.length
        ? projectArg.stalls
        : (projects.value[pIdx]!.stalls || [])
      const nextStalls = sourceStalls.filter((s: any) => String(s.code || s.ident).toUpperCase() !== code.toUpperCase())
      projects.value[pIdx] = {
        ...projects.value[pIdx]!,
        layoutMode: 'manual',
        stalls: nextStalls,
        totalStalls: nextStalls.length,
        stallsAvailable: nextStalls.filter((s: any) => (s.status || 'available') === 'available').length
      }
      if (selectedProject.value?.id === targetProject.id) {
        selectedProject.value = projects.value[pIdx]!
      }
    }
  }

  const dbId = stored?.id || stallEdit.value?.id
  if (dbId) pendingStallDeleteIds.value.push(Number(dbId))
  pendingStallUpserts.value = pendingStallUpserts.value.filter(item => String(item.ident || item.stall_ident).toUpperCase() !== code.toUpperCase())
  stalls.value = stalls.value.filter(s => String(s.ident || s.stall_ident).toUpperCase() !== code.toUpperCase() && (!dbId || s.id !== dbId))
  hasUnsavedChanges.value = true
  if (showStallEditor.value) closeStall()

  void showSuccess(`ลบแผง “${code}” ออกจากผังแล้ว กดปุ่ม “บันทึกการแก้ไข” ด้านบนเพื่อบันทึกผัง`, 'ลบแผงในผังแล้ว')
}
const saveStall = async () => {
  if (!validateStallForm()) return
  savingStall.value = true
  try {
    const draftSaved = { ...stallEdit.value }
    const existingStallIdx = stalls.value.findIndex(s => (draftSaved.id && s.id === draftSaved.id) || String(s.ident || s.stall_ident).toUpperCase() === String(draftSaved.ident).toUpperCase())
    if (existingStallIdx >= 0) {
      const nextStalls = [...stalls.value]
      nextStalls[existingStallIdx] = { ...nextStalls[existingStallIdx], ...draftSaved }
      stalls.value = nextStalls
    } else {
      stalls.value = [...stalls.value, draftSaved]
    }
    queueStallUpsert(draftSaved)

    if (pendingStallLayout.value) {
      const pending = pendingStallLayout.value
      const projectIndex = projects.value.findIndex(project => project.id === pending.projectId)
      if (projectIndex >= 0) {
        const project = projects.value[projectIndex]!
        const chosenFloorId = String(stallEdit.value.floor || pending.floorId || '1')
        const floor = project.floors?.find(item => String(item.id) === chosenFloorId) || project.floors?.[0] || { id: chosenFloorId, name: `ชั้น ${chosenFloorId}`, level: Number(chosenFloorId) || 1 }
        const rent = Number(stallEdit.value.monthly_rent) || 1800
        const manualStall = {
          code: draftSaved.ident || stallEdit.value.ident,
          projectId: project.id,
          zone: project.name,
          status: stallEdit.value.stall_status || draftSaved.stall_status || 'available',
          areaType: stallEdit.value.category || 'ทั่วไป',
          size: Number((Number(stallEdit.value.width) * Number(stallEdit.value.depth)).toFixed(1)),
          widthM: Number(stallEdit.value.width),
          depthM: Number(stallEdit.value.depth),
          entryFee: 0,
          monthlyRent: rent,
          commonFee: 0,
          depositFee: rent,
          photo: '',
          coordinates: pending.coordinates,
          floor: String(floor.id),
          floorName: floor.name,
          floorLevel: floor.level
        }
        const allStalls = [...(project.stalls || []), manualStall]
        projects.value[projectIndex] = {
          ...project,
          layoutMode: 'manual',
          stalls: allStalls,
          totalStalls: allStalls.length,
          stallsAvailable: allStalls.filter((stall: any) => stall.status === 'available').length
        }
        selectedProject.value = projects.value[projectIndex]!
        pendingStallLayout.value = null
      }
    } else if (selectedProject.value) {
      const pIdx = projects.value.findIndex(p => p.id === selectedProject.value?.id)
      if (pIdx >= 0 && Array.isArray(projects.value[pIdx]?.stalls)) {
        const targetProj = projects.value[pIdx]!
        const targetCode = String(draftSaved.ident || stallEdit.value.ident)
        const updatedStalls = (targetProj.stalls || []).map((st: any) => {
          if (String(st.code).toUpperCase() !== targetCode.toUpperCase()) return st
          return {
            ...st,
            status: stallEdit.value.stall_status || draftSaved.stall_status || st.status,
            floor: String(stallEdit.value.floor || st.floor || '1'),
            widthM: Number(stallEdit.value.width) || st.widthM,
            depthM: Number(stallEdit.value.depth) || st.depthM,
            size: Number((Number(stallEdit.value.width || st.widthM) * Number(stallEdit.value.depth || st.depthM)).toFixed(1)),
            monthlyRent: Number(stallEdit.value.monthly_rent) || st.monthlyRent
          }
        })
        projects.value[pIdx] = {
          ...targetProj,
          stalls: updatedStalls,
          stallsAvailable: updatedStalls.filter((s: any) => s.status === 'available').length
        }
        selectedProject.value = projects.value[pIdx]!
      }
    }
    hasUnsavedChanges.value = true
    closeStall()
    void showSuccess('อัปเดตข้อมูลแผงบนผังแล้ว กดปุ่ม “บันทึกการแก้ไข” ด้านบนเพื่อบันทึกผัง', 'อัปเดตแผงในผังแล้ว')
  } finally {
    savingStall.value = false
  }
}
const saveGIS = async (publish: boolean) => {
  savingGIS.value = true
  try {
    await persistGISChanges(publish)
    void showSuccess(publish ? 'บันทึกการแก้ไขและอัปเดตผังตลาดเรียบร้อยแล้ว' : 'เก็บการแก้ไขเป็นฉบับร่างแล้ว', publish ? 'บันทึกสำเร็จ' : 'บันทึกฉบับร่างสำเร็จ')
  } catch {
    void showError('ไม่สามารถบันทึกผังได้ กรุณาตรวจสอบสิทธิ์และข้อมูลพิกัด', 'บันทึกไม่สำเร็จ')
  } finally {
    savingGIS.value = false
  }
}
const confirmAndSaveGIS = async () => {
  const confirmed = await showConfirm(
    hasUnsavedChanges.value
      ? 'คุณได้ทำการแก้ไขข้อมูลผังตลาด ต้องการบันทึกการแก้ไขและเผยแพร่ผังล่าสุดใช่หรือไม่?'
      : 'ต้องการบันทึกและเผยแพร่ผังตลาดใช่หรือไม่?',
    {
      title: 'ยืนยันการบันทึกการแก้ไข',
      type: 'info',
      confirmText: 'บันทึกการแก้ไข',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return
  await saveGIS(true)
}
const formatDate = (value: string) => new Intl.DateTimeFormat('th-TH', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))

const loadData = async () => {
  isManageDataLoading.value = true
  const cfg = await loadMarketGisJson(marketId.value)
  apiMarketName.value = cfg.name
  try {
    const [marketRes, stallRes, gisRes, bookingsRes, waitlistRes, appsRes, contractsRes] = await Promise.all([
      get('/markets').catch(() => null),
      get(`/stalls?market=${encodeURIComponent(marketId.value)}`).catch(() => null),
      get(`/markets/${encodeURIComponent(marketId.value)}/gis-config/draft`).catch(() => null),
      get('/bookings').catch(() => null),
      get('/bookings/waitlist/mine').catch(() => null),
      get('/applications').catch(() => null),
      get('/contracts').catch(() => null)
    ])
    const markets = unwrap(marketRes)
    const market = (Array.isArray(markets) ? markets : []).find((item: any) => String(item.id) === marketId.value || item.slug === marketId.value || String(item.id) === String(cfg.id))
    apiMarketRecord.value = market || null
    resolvedMarketID.value = Number(market?.id || cfg.id)
    apiMarketName.value = market?.name || cfg.name
    const stallData = unwrap(stallRes)
    stalls.value = Array.isArray(stallData)
      ? stallData.filter((s: any) => !s.market_id || Number(s.market_id) === resolvedMarketID.value)
      : []
    const gis = unwrap(gisRes) || {}
    projects.value = mergeMarketProjects(cfg.projects, gis.projects)
    facilities.value = Array.isArray(gis.facilities) ? gis.facilities : []
    gisPublished.value = Boolean(gis.published)

    const allBookings = unwrap(bookingsRes)
    const allWaitlist = unwrap(waitlistRes)
    const allApps = unwrap(appsRes)
    const allContracts = unwrap(contractsRes)
    const targetMktName = String(apiMarketName.value || cfg.name || '').trim()
    const knownCodes = new Set<string>()
    projects.value.forEach((p) => {
      getAllFloorLayoutsForProject(p).forEach((layout) => {
        layout.stalls.forEach((st: any) => {
          const c = String(st.code || st.ident || '').trim().toUpperCase()
          if (c) knownCodes.add(c)
        })
      })
      ;(p.stalls || []).forEach((st: any) => {
        const c = String(st.code || st.ident || '').trim().toUpperCase()
        if (c) knownCodes.add(c)
      })
    })
    const matchesCurrentMarket = (item: any) => {
      if (!item) return false
      if (Number(item.market_id) > 0) {
        return Number(item.market_id) === resolvedMarketID.value
      }
      const itemMktName = String(item.market_name || '').trim()
      if (itemMktName && targetMktName) {
        return itemMktName.includes(targetMktName) || targetMktName.includes(itemMktName)
      }
      const rawStall = String(item.stall_number || item.stall_ident || item.stall_code || '').trim()
      if (rawStall && knownCodes.size > 0) {
        return rawStall.split(',').some(code => knownCodes.has(code.trim().toUpperCase()))
      }
      return false
    }
    marketBookings.value = Array.isArray(allBookings) ? allBookings.filter(matchesCurrentMarket) : []
    marketWaitlistItems.value = Array.isArray(allWaitlist) ? allWaitlist.filter(matchesCurrentMarket) : []
    marketApplications.value = Array.isArray(allApps) ? allApps.filter(matchesCurrentMarket) : []
    marketContracts.value = Array.isArray(allContracts) ? allContracts.filter(matchesCurrentMarket) : []
    loadAllCustomContracts()
  } catch {
    projects.value = cloneProjects(cfg.projects)
    facilities.value = []
  } finally {
    syncSnapshotsFromCurrent()
    isManageDataLoading.value = false
  }
  const targetProjectId = returnProjectId.value || selectedProject.value?.id
  if (targetProjectId) selectedProject.value = projects.value.find(p => p.id === targetProjectId) || null
  if (returnStallCode.value && selectedProject.value) {
    await nextTick()
    await restoreReturnedStall()
  }
}

const handleBeforeUnload = (event: BeforeUnloadEvent) => {
  if (!isInspectMode.value && hasUnsavedChanges.value) {
    event.preventDefault()
    event.returnValue = ''
  }
}


// ================================================================
// โหมดจัดการค่าใช้จ่าย (ข้อมูลจำลอง): แผงรายการ + ลิ้นชักรายแผง + การ์ดกรอง
// ================================================================
const billing = useStallBilling()
const costTool = ref<'utility' | 'bill' | null>(null)
const costFilter = ref<CostFilter>('all')
const costDrawerOpen = ref(false)
const costDrawerKind = ref<'utility' | 'temp'>('utility')
const costStallCode = ref('')
const clickedCostStalls = new Map<string, { occupied: boolean; dailyRate: number }>()

const rightSidePanelWidth = computed(() => {
  if (costDrawerOpen.value) return 380
  // แผงแดชบอร์ด
  if (activeMapSidePanel.value === 'dashboard') return 380
  if (activeMapSidePanel.value) return 380
  return 0
})


const occupiedCostStalls = computed(() => {
  const seen = new Set<string>()
  Object.entries(stallStatusOverrides.value).forEach(([code, st]) => {
    if (st === 'occupied') seen.add(code.toUpperCase())
  })
  return [...seen].map(code => ({ code, tenant: tenantNameOf(code) }))
})
watch([occupiedCostStalls, () => billing.billingMonth.value], () => billing.ensureUtilityBills(occupiedCostStalls.value), { immediate: true })
watch(stalls, (list) => {
  const vacant = list
    .map((x: any) => String(x.ident || x.stall_ident || '').trim())
    .filter((code: string) => code && resolveRealStallStatus(code, null, null) !== 'occupied')
  billing.seedTempBills(vacant)
}, { immediate: true })

const lookupCostStall = (code: string) => {
  const upper = code.toUpperCase()
  const stored = stalls.value.find((x: any) => String(x.ident || x.stall_ident || '').toUpperCase() === upper)
  if (stored) return { occupied: resolveRealStallStatus(upper, null, stored) === 'occupied', dailyRate: Number(stored.daily_rent) || 150 }
  return clickedCostStalls.get(upper) || null
}
const openCostDrawer = (code: string, kind: 'utility' | 'temp') => {
  costStallCode.value = code.toUpperCase()
  costDrawerKind.value = kind
  costDrawerOpen.value = true
}
const closeCostDrawer = () => {
  costDrawerOpen.value = false
  costStallCode.value = ''
}
const toggleCostTool = (tool: 'utility' | 'bill') => {
  costTool.value = costTool.value === tool ? null : tool
  costFilter.value = 'all'
  closeCostDrawer()
}
const closeCostTool = () => {
  costTool.value = null
  costFilter.value = 'all'
  closeCostDrawer()
}
const handleCostStallClick = (mapStall: any) => {
  const code = String(mapStall?.code || mapStall?.ident || '').trim().toUpperCase()
  if (!code) return
  const status = resolveRealStallStatus(code, selectedProject.value, mapStall)
  const occupied = status === 'occupied'
  const stored = stalls.value.find((x: any) => String(x.ident || x.stall_ident || '').toUpperCase() === code)
  clickedCostStalls.set(code, { occupied, dailyRate: Number(stored?.daily_rent) || 150 })
  if (occupied) billing.ensureUtilityBills([{ code, tenant: tenantNameOf(code) }])
  // ไม่ได้เปิดแผงไหนอยู่: เลือกชนิดตามสถานะแผง (มีผู้เช่า = ค่าน้ำไฟ, ไม่มี = บิลชั่วคราว)
  openCostDrawer(code, costTool.value ? (costTool.value === 'utility' ? 'utility' : 'temp') : (occupied ? 'utility' : 'temp'))
}

const costCards = computed(() => {
  const isBill = costTool.value === 'bill'
  const icons = {
    all: '<rect x="4" y="4" width="6" height="6"/><rect x="14" y="4" width="6" height="6"/><rect x="4" y="14" width="6" height="6"/><rect x="14" y="14" width="6" height="6"/>',
    todo: '<path d="M12 20h9"/><path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4z"/>',
    pending: '<circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/>',
    overdue: '<path d="M10.3 3.9 1.8 18a2 2 0 0 0 1.7 3h17a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0z"/><path d="M12 9v4M12 17h.01"/>',
    paid: '<circle cx="12" cy="12" r="9"/><path d="m8 12 2.5 2.5L16 9"/>'
  }
  const baht = (v: number) => `฿${Math.round(v).toLocaleString('th-TH')}`
  let rows: { status: string; amount: number }[]
  if (isBill) {
    rows = billing.tempBills.value.map(b => ({ status: b.status, amount: billing.tempTotal(b) }))
  } else {
    rows = billing.utilityList.value.map(b => ({ status: b.status, amount: billing.utilityTotal(b) || 0 }))
  }
  const by = (st: string) => rows.filter(r => r.status === st)
  const sum = (list: { amount: number }[]) => list.reduce((a, r) => a + r.amount, 0)
  const unissued = Math.max(0, (currentDisplayedSummary.value.available || 0) - billing.tempBills.value.length)
  return [
    { filter: 'all' as CostFilter, tone: 'is-total', icon: icons.all, label: isBill ? 'บิลทั้งหมด' : 'ทั้งหมด', count: rows.length, unit: isBill ? 'บิล' : 'แผง', sub: `ยอดรวม ${baht(sum(rows))}` },
    isBill
      ? { filter: 'todo' as CostFilter, tone: 'is-occupied', icon: icons.todo, label: 'ยังไม่ออกบิล', count: unissued, unit: 'แผงว่าง', sub: 'คลิกแผงเพื่อออกบิล' }
      : { filter: 'todo' as CostFilter, tone: 'is-occupied', icon: icons.todo, label: 'ยังไม่จดมิเตอร์', count: by('unrecorded').length, unit: 'แผง', sub: 'ต้องจดก่อนออกบิล' },
    { filter: 'pending' as CostFilter, tone: 'is-reserved', icon: icons.pending, label: 'รอชำระ', count: by('pending').length, unit: isBill ? 'บิล' : 'แผง', sub: `รอเก็บ ${baht(sum(by('pending')))}` },
    { filter: 'overdue' as CostFilter, tone: 'is-overdue', icon: icons.overdue, label: 'เกินกำหนด', count: by('overdue').length, unit: isBill ? 'บิล' : 'แผง', sub: `ค้าง ${baht(sum(by('overdue')))}` },
    { filter: 'paid' as CostFilter, tone: 'is-available', icon: icons.paid, label: 'ชำระแล้ว', count: by('paid').length, unit: isBill ? 'บิล' : 'แผง', sub: `เก็บได้ ${baht(sum(by('paid')))}` }
  ]
})

const costMonth = billing.billingMonth
const costAmount = billing.amountRange
const costAmountOptions = [
  { value: 'all', label: 'ทั้งหมด' },
  { value: 'lt500', label: 'ต่ำกว่า ฿500' },
  { value: '500-1000', label: '฿500 – 1,000' },
  { value: 'gt1000', label: 'มากกว่า ฿1,000' }
] as { value: any; label: string }[]
const costMonthOptions = computed(() => {
  const base = new Date()
  return Array.from({ length: 3 }, (_, i) => {
    const d = new Date(base.getFullYear(), base.getMonth() - i, 1)
    return {
      value: `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`,
      label: new Intl.DateTimeFormat('th-TH', { month: 'long', year: 'numeric' }).format(d)
    }
  })
})
// ตัวเลขบนปุ่มตัวกรองนับเฉพาะตัวกรองที่ตั้งในแผงตัวกรอง (ช่วงยอดเงิน) ไม่นับสถานะที่เลือกจากการ์ดด้านบน
const costFilterBadge = computed(() => (costAmount.value !== 'all' ? 1 : 0))
const resetCostFilters = () => {
  costFilter.value = 'all'
  costAmount.value = 'all'
}
// ความกว้างของแผงที่เปิดอยู่ทางซ้าย: ใช้เลื่อนปุ่มลอยไปอยู่ข้างแผง และให้แผงชิดล่าง
const leftPanelWidth = computed(() => {
  if (pageMode.value === 'costs') return costTool.value ? '640px' : ''
  if (dashboardPanelOpen.value) return '380px'
  // ให้คอนโทรลลอยเคลื่อนพร้อมกรอบที่กำลัง morph ไม่กระโดดตอนสลับเนื้อหา
  if (isProjectPanelMorphing.value) return '720px'
  if (showInspectStallOverlay.value) return stallDrawerExpanded.value ? '720px' : (stallDrawerMode.value === 'edit' ? '440px' : '380px')
  return selectedProject.value ? '380px' : ''
})
const handleClearStallSelection = () => {
  // ขณะดูรายละเอียดแผง การคลิกพื้นที่อื่นบนแผนที่ต้องไม่ยุบ/ปิดกล่อง
  // ให้กลับไปกล่องโครงการได้จากปุ่มย้อนกลับใน panel เท่านั้น
  if (showInspectStallOverlay.value) return
}
const onStallEsc = (e: KeyboardEvent) => {
  if (e.key !== 'Escape' || showInspectPhotoModal.value) return
  if (showInspectStallOverlay.value) closeStall()
  else handleClearStallSelection()
}

const costFilterCodes = computed<string[]>(() => {
  if (costFilter.value === 'all' && costAmount.value === 'all') return []
  if (costTool.value === 'bill') {
    return billing.tempBills.value
      .filter(b => billing.matchesFilter(b.status, costFilter.value) && billing.matchesAmount(billing.tempTotal(b)))
      .map(b => b.stall)
  }
  return billing.utilityList.value
    .filter(b => billing.matchesFilter(b.status, costFilter.value) && billing.matchesAmount(billing.utilityTotal(b)))
    .map(b => b.stall)
})

onMounted(() => {
  void loadData()
  loadAllCustomContracts()
  if (import.meta.client) {
    window.addEventListener('beforeunload', handleBeforeUnload)
    window.addEventListener('keydown', onStallEsc)
    window.addEventListener('click', handleClickOutsideDropdown)
    if (isMapOverviewMode.value) {
      nextTick(() => {
        setTimeout(() => {
          window.dispatchEvent(new Event('resize'))
        }, 150)
      })
    }
  }
})

onBeforeUnmount(() => {
  if (import.meta.client) {
    window.removeEventListener('beforeunload', handleBeforeUnload)
    window.removeEventListener('keydown', onStallEsc)
    window.removeEventListener('click', handleClickOutsideDropdown)
  }
})

onBeforeRouteLeave(async () => {
  const canProceed = await ensureCanLeaveOrBack()
  return canProceed
})

onBeforeRouteUpdate(async (to, from) => {
  if (to.query.market !== from.query.market) {
    const canProceed = await ensureCanLeaveOrBack()
    return canProceed
  }
  return true
})

watch(marketId, () => { selectedProject.value = null; closeStall(); void loadData() })

watch(() => route.fullPath, () => {
  if (isMapOverviewMode.value && import.meta.client) {
    nextTick(() => {
      setTimeout(() => {
        window.dispatchEvent(new Event('resize'))
      }, 100)
    })
  }
})

const CloseIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [h('path', { d: 'm6 6 12 12M18 6 6 18' })])
</script>

<style scoped>
.gis-panel{padding:26px 28px!important}.gis-panel-heading{align-items:center;margin-bottom:22px!important}.gis-panel-heading h2{font-size:1.25rem!important;font-weight:800;color:#0f172a;margin:0}.gis-panel-heading p{margin:6px 0 0!important}.gis-heading-right{display:flex;align-items:center;gap:12px;margin-left:auto;flex-wrap:wrap}.gis-map-toolbar{display:flex;flex:0 1 auto;min-width:0}.gis-map-toolbar :deep(.master-map-top-bar){width:auto;margin:0}.gis-map-toolbar :deep(.master-map-tools-right){margin-left:0;flex-wrap:nowrap}@media(max-width:700px){.gis-panel-heading{align-items:flex-start}.gis-heading-right,.gis-map-toolbar{width:100%;margin-left:0}.gis-map-toolbar :deep(.master-map-top-bar),.gis-map-toolbar :deep(.master-map-tools-right){width:100%;flex-wrap:wrap}.gis-map-toolbar :deep(.search-input-group){flex:1;width:auto}}
.staff-dashboard-intro{display:flex;align-items:center;justify-content:space-between;gap:16px;flex-wrap:wrap;margin:8px 0 22px}.dashboard-intro-left h1{margin:0;color:#0f172a;font-size:1.5rem;font-weight:700;letter-spacing:-.01em;line-height:1.3}.dashboard-updated-sub{margin:4px 0 0;color:#64748b;font-size:.875rem;font-weight:500}.dashboard-intro-actions{display:flex;align-items:center;gap:12px;flex-wrap:wrap}
.officer-summary-dashboard{display:grid;gap:20px;margin-bottom:12px}.dash-kpi-ribbon{margin:0!important}.dash-kpi-ribbon .map-summary-card{cursor:default}.dash-kpi-ribbon .map-summary-card:hover{transform:none;box-shadow:0 1px 3px rgba(15,23,42,.03)}.dash-two-col-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:20px;align-items:stretch}
.dash-system-card{background:#fff;border:1px solid #eef2f6;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.03);padding:22px 24px;display:flex;flex-direction:column;gap:16px}
.dash-card-header{display:flex;align-items:center;justify-content:space-between;gap:14px;flex-wrap:wrap}.dash-card-title{margin:0;color:#1e293b;font-size:1.05rem;font-weight:700;line-height:1.3}.dash-card-subtitle{margin:3px 0 0;color:#64748b;font-size:.8125rem;font-weight:500}
.dash-status-badge{display:inline-flex;align-items:center;justify-content:center;padding:3px 10px;border-radius:7px;font-size:.75rem;font-weight:600;white-space:nowrap}.dash-status-badge.badge-amber{background:#fef3c7;color:#d97706}.dash-status-badge.badge-green{background:#dcfce7;color:#15803d}.dash-status-badge.badge-red{background:#fee2e2;color:#dc2626}.dash-status-badge.badge-gray{background:#f1f5f9;color:#64748b}
.dash-donut-layout{display:grid;grid-template-columns:148px minmax(0,1fr);gap:22px;align-items:center;margin:auto 0}
.dash-donut-wrap{position:relative;width:148px;height:148px;display:flex;align-items:center;justify-content:center;margin:0 auto}.dash-donut-svg{width:148px;height:148px;display:block}.dash-donut-center{position:absolute;inset:0;display:flex;flex-direction:column;align-items:center;justify-content:center;text-align:center;pointer-events:none}.dash-donut-center strong{color:#0f172a;font-size:1.25rem;font-weight:800;line-height:1.1}.dash-donut-center span{color:#64748b;font-size:.75rem;font-weight:600;margin-top:3px}
.dash-donut-legend-list{display:grid;gap:8px}.dash-legend-row{width:100%;display:flex;align-items:center;justify-content:space-between;gap:10px;padding:9px 12px;border:1px solid #eef2f6;border-radius:10px;background:#f8fafc}.legend-left{display:inline-flex;align-items:center;gap:8px;color:#334155;font-size:.875rem;font-weight:600}.legend-left i{width:10px;height:10px;border-radius:50%;display:inline-block}.dot-available{background:#10b981}.dot-reserved{background:#f59e0b}.dot-occupied{background:#475569}.legend-right{display:inline-flex;align-items:center;gap:8px}.legend-right strong{color:#0f172a;font-size:.9375rem;font-weight:700}
.dash-mini-rev-row{display:grid;grid-template-columns:1fr 1fr;gap:10px;padding-top:12px;border-top:1px solid #f1f5f9}.dash-mini-rev-item{display:flex;align-items:center;justify-content:space-between;gap:8px;padding:10px 14px;border-radius:10px;background:#f8fafc;border:1px solid #eef2f6}.dash-mini-rev-item span{color:#64748b;font-size:.8125rem;font-weight:500}.dash-mini-rev-item strong{font-size:1rem;font-weight:700;color:#0f172a}.dash-mini-rev-item strong.is-blue{color:#2563eb}.dash-mini-rev-item strong.is-green{color:#059669}
.dash-bar-chart-list{display:grid;gap:14px;margin:auto 0}.dash-bar-row{display:grid;gap:6px}.dash-bar-top{display:flex;align-items:center;justify-content:space-between;gap:10px}.dash-bar-name{display:inline-flex;align-items:center;gap:8px;color:#1e293b;font-size:.875rem;font-weight:600;min-width:0;max-width:62%;overflow:hidden;white-space:nowrap}.dash-bar-name>span{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.dash-bar-name i{width:10px;height:10px;border-radius:3px;flex-shrink:0}.dash-bar-meta{color:#64748b;font-size:.8125rem;font-weight:500;white-space:nowrap}.dash-bar-meta strong{color:#0f172a;font-weight:700}.dash-bar-track{height:10px;border-radius:999px;background:#f1f5f9;overflow:hidden}.dash-bar-fill{display:block;height:100%;border-radius:999px;transition:width .3s ease}
@media(max-width:640px){.dash-donut-layout{grid-template-columns:1fr}.dash-mini-rev-row{grid-template-columns:1fr}}
.dash-table-card{padding:24px 28px;gap:16px}.dash-table-header{display:flex;align-items:center;justify-content:space-between;gap:16px;flex-wrap:wrap}.dash-table-header-actions{display:flex;align-items:center;gap:10px;flex-wrap:wrap}.dash-category-select{min-width:210px}
.dash-search-box{position:relative;display:inline-flex;align-items:center;min-width:240px}.dash-search-box svg{position:absolute;left:12px;color:#94a3b8;pointer-events:none;z-index:1}.dash-search-box :deep(.app-input){min-height:40px;padding-left:36px;border-color:#e2e8f0;border-radius:8px}
.dash-filter-tabs{display:flex;align-items:center;gap:24px;border-bottom:1px solid #f1f5f9;overflow-x:auto}.dash-filter-tab{position:relative;padding:10px 4px 12px;border:0;background:transparent;color:#64748b;font:inherit;font-size:.9375rem;font-weight:600;cursor:pointer;white-space:nowrap;transition:color .15s}.dash-filter-tab:hover{color:#1e293b}.dash-filter-tab.active{color:#2563eb}.dash-filter-tab.active::after{content:'';position:absolute;left:0;right:0;bottom:-1px;height:3px;border-radius:2px;background:#2563eb}
.dash-table-wrap{width:100%;overflow-x:auto}.dash-std-table{width:100%;border-collapse:collapse;text-align:left}.dash-std-table thead th{padding:14px 16px;color:#64748b;font-size:.8125rem;font-weight:500;border-bottom:1px solid #f1f5f9;white-space:nowrap;background:#fff}.dash-std-table tbody td{padding:16px;color:#334155;font-size:.875rem;font-weight:500;border-bottom:1px solid #f1f5f9;vertical-align:middle}.dash-std-table tbody tr:hover{background:#f8fafc}.dash-std-table .text-right{text-align:right}.dash-std-table .text-center{text-align:center}.dash-std-table strong{color:#1e293b;font-weight:700}.dash-std-table .code-primary{color:#2563eb}
.zone-name-cell{display:flex;align-items:center;gap:12px}.zone-color-dot{width:12px;height:12px;border-radius:4px;flex-shrink:0}.zone-name-cell strong{display:block;color:#1e293b;font-size:.875rem;font-weight:600}.zone-name-cell small,.applicant-cell small{display:block;margin-top:2px;color:#64748b;font-size:.75rem;font-weight:500}.applicant-cell strong{display:block;color:#1e293b;font-size:.875rem;font-weight:600}
.zone-occ-cell{display:flex;align-items:center;gap:10px}.zone-occ-bar{flex:1;height:8px;border-radius:999px;background:#e2e8f0;overflow:hidden;min-width:80px}.zone-occ-bar span{display:block;height:100%;border-radius:999px;background:#2563eb}.zone-occ-cell strong{font-size:.8125rem;min-width:36px;text-align:right}
.dash-table-view-btn{display:inline-flex;align-items:center;gap:6px;padding:6px 14px;border:1px solid #e2e8f0;border-radius:8px;background:#fff;color:#334155;font:inherit;font-size:.825rem;font-weight:600;text-decoration:none;cursor:pointer;white-space:nowrap;transition:background .15s,border-color .15s,color .15s}.dash-table-view-btn:hover{background:#eff6ff;border-color:#93c5fd;color:#2563eb}
.dash-empty-row{padding:32px 16px!important;text-align:center;color:#64748b!important;font-size:.875rem}
@media(max-width:1080px){.dash-two-col-grid{grid-template-columns:1fr}}@media(max-width:768px){.dash-revenue-grid{grid-template-columns:1fr}.dash-queue-item{flex-direction:column;align-items:flex-start}.dash-queue-right{width:100%;justify-content:space-between}}
.workspace-grid{display:grid;grid-template-columns:minmax(0,1fr)!important;margin-top:0}.workspace-grid[style*="display: none"]{display:none!important}.gis-panel-heading h2{font-size:1.15rem}.map-summary-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:18px;margin:4px 0 24px}.map-summary-card{background:#fff;border:1px solid #eef2f6;border-radius:14px;padding:16px 20px;display:flex;align-items:center;gap:16px;cursor:pointer;font-family:inherit;text-align:left;position:relative;transition:all .2s cubic-bezier(.4,0,.2,1);box-shadow:0 1px 3px rgba(15,23,42,.03)}.map-summary-card:hover{transform:translateY(-2px);box-shadow:0 6px 16px rgba(15,23,42,.07)}.map-summary-card:focus-visible{outline:2px solid #2563eb;outline-offset:2px}.map-summary-icon{width:46px;height:46px;border-radius:12px;display:flex;align-items:center;justify-content:center;flex-shrink:0;transition:all .2s ease}.map-summary-content{display:flex;flex-direction:column;min-width:0}.map-summary-title{font-size:.8125rem;color:#64748b;font-weight:500;line-height:1.25;margin-bottom:6px;white-space:nowrap}.map-summary-value-row{display:flex;align-items:baseline;gap:8px}.map-summary-count{font-size:1.45rem;font-weight:800;color:#0f172a;line-height:1.1}.map-summary-unit-pill{font-size:.78rem;font-weight:600;padding:2px 8px;border-radius:6px;line-height:1.2;font-style:normal}.summary-total .map-summary-icon{background:#eff6ff;color:#2563eb;border:1px solid #dbeafe}.summary-total .map-summary-unit-pill{background:#eff6ff;color:#1d4ed8}.summary-total.active{border-color:#2563eb;background:#f8faff;box-shadow:0 0 0 2px rgba(37,99,235,.18)}.summary-available .map-summary-icon{background:#ecfdf5;color:#059669;border:1px solid #d1fae5}.summary-available .map-summary-unit-pill{background:#dcfce7;color:#15803d}.summary-available.active{border-color:#10b981;background:#f0fdf4;box-shadow:0 0 0 2px rgba(16,185,129,.2)}.summary-reserved .map-summary-icon{background:#fffbeb;color:#d97706;border:1px solid #fef3c7}.summary-reserved .map-summary-unit-pill{background:#fef3c7;color:#b45309}.summary-reserved.active{border-color:#f59e0b;background:#fffdf7;box-shadow:0 0 0 2px rgba(245,158,11,.2)}.summary-occupied .map-summary-icon{background:#f1f5f9;color:#64748b;border:1px solid #e2e8f0}.summary-occupied .map-summary-unit-pill{background:#f1f5f9;color:#64748b}.summary-occupied.active{border-color:#94a3b8;background:#f8fafc;box-shadow:0 0 0 2px rgba(100,116,139,.18)}@media(max-width:800px){.map-summary-grid{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:440px){.map-summary-grid{grid-template-columns:1fr}}.header-actions{display:inline-flex;align-items:center;gap:10px}.draft-status{display:inline-flex;align-items:center;gap:7px;height:34px;box-sizing:border-box;padding:0 12px;border-radius:999px;background:#f1f5f9;color:#475569;font-size:.8125rem;font-weight:600;transition:all .2s ease}.draft-status i{width:8px;height:8px;border-radius:50%;background:currentColor}.draft-status.published{background:#ecfdf5;color:#15803d;border:1px solid #bbf7d0}.draft-status.is-dirty{background:#fff7ed;color:#c2410c;border:1px solid #fed7aa}.draft-status.is-dirty i{background:#ea580c;box-shadow:0 0 0 3px rgba(234,88,12,.18)}.header-status-skeleton,.header-action-skeleton{display:block;border-radius:999px;background:linear-gradient(90deg,#eef2f7 25%,#f8fafc 50%,#eef2f7 75%);background-size:200% 100%;animation:headerSkeletonShimmer 1.35s ease-in-out infinite}.header-status-skeleton{width:102px;height:34px}.header-action-skeleton{width:128px;height:34px;border-radius:7px}@keyframes headerSkeletonShimmer{0%{background-position:100% 0}100%{background-position:-100% 0}}@media(prefers-reduced-motion:reduce){.header-status-skeleton,.header-action-skeleton{animation:none}}
.unsaved-leave-backdrop{position:fixed;inset:0;background:rgba(15,23,42,.65);backdrop-filter:blur(8px);-webkit-backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;padding:20px;z-index:99999}.unsaved-leave-card{background:#fff;border-radius:18px;width:100%;max-width:430px;padding:28px 24px 22px;text-align:center;box-shadow:0 24px 48px -12px rgba(15,23,42,.28),0 0 0 1px rgba(226,232,240,.8);position:relative}.unsaved-leave-close{position:absolute;top:14px;right:14px;width:32px;height:32px;border-radius:8px;border:0;background:transparent;color:#64748b;display:inline-flex;align-items:center;justify-content:center;cursor:pointer;transition:background .15s,color .15s}.unsaved-leave-close:hover{background:#f1f5f9;color:#0f172a}.unsaved-leave-close svg{width:18px;height:18px}.unsaved-leave-icon{width:56px;height:56px;border-radius:50%;margin:0 auto 16px;display:flex;align-items:center;justify-content:center;background:#fef3c7;color:#d97706;box-shadow:0 0 0 6px #fffbeb}.unsaved-leave-content h3{font-size:1.15rem;font-weight:700;color:#0f172a;margin:0 0 8px}.unsaved-leave-content p{font-size:.9375rem;color:#475569;line-height:1.55;margin:0 0 22px}.unsaved-leave-actions{display:grid;grid-template-columns:1fr 1fr;gap:10px}
.inspect-selected-stall-drawer{position:absolute;top:62px;right:14px;width:min(380px,calc(100vw - 28px));max-height:calc(100% - 76px);display:flex;flex-direction:column;pointer-events:auto;z-index:125;background:#ffffff;border:1px solid #e2e8f0;border-radius:18px;box-shadow:0 16px 38px rgba(15,23,42,.16);overflow:hidden;contain:layout paint;will-change:width,transform,opacity;transition:width .46s cubic-bezier(.16,1,.3,1),box-shadow .36s ease,background-color .28s ease}
.inspect-selected-stall-drawer.is-edit-mode{width:min(440px,calc(100% - 28px))}
.inspect-selected-stall-drawer.is-full-detail-mode{width:min(760px,calc(100vw - 36px));z-index:130}
@media(min-width:769px){.inspect-selected-stall-drawer.is-management-mode{bottom:100px;max-height:calc(100% - 176px)}}
.issd-header{padding:16px 18px 14px;border-bottom:1px solid #f1f5f9;background:#fff;flex-shrink:0}.issd-header-top{display:flex;align-items:center;justify-content:space-between;gap:10px}.issd-title-group{display:flex;align-items:center;gap:8px;flex-wrap:wrap;min-width:0}.issd-title-group strong{font-size:1rem;font-weight:700;color:#0f172a}.issd-project-sub{margin-top:4px;font-size:.8125rem;font-weight:500;color:#64748b;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.issd-back-btn{width:28px;height:28px;border-radius:8px;border:1px solid #e2e8f0;background:#fff;color:#475569;display:inline-flex;align-items:center;justify-content:center;cursor:pointer;flex-shrink:0;transition:background .15s,border-color .15s,color .15s}.issd-back-btn:hover{background:#f1f5f9;color:#0f172a;border-color:#cbd5e1}
.issd-close-btn{width:28px;height:28px;border-radius:8px;border:1px solid #e2e8f0;background:#fff;color:#475569;display:inline-flex;align-items:center;justify-content:center;cursor:pointer;flex-shrink:0;transition:background .15s,border-color .15s,color .15s}.issd-close-btn:hover{background:#fef2f2;border-color:#fca5a5;color:#ef4444}
.issd-detail-tabs-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  flex-shrink: 0;
  overflow: visible;
}
.issd-tab-btn {
  flex: 1 1 auto;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  height: 32px;
  padding: 0 10px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  color: #475569;
  font-family: inherit;
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.03);
}
.issd-tab-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
  border-color: #cbd5e1;
}
.issd-tab-btn.active {
  background: #2563eb;
  color: #ffffff;
  border-color: #1d4ed8;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.22);
  font-weight: 700;
}
.issd-tab-btn.active svg {
  color: #ffffff;
}
.issd-tab-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 17px;
  height: 17px;
  padding: 0 5px;
  border-radius: 999px;
  background: #f1f5f9;
  color: #475569;
  font-size: 0.75rem;
  font-weight: 700;
  line-height: 1;
}
.issd-tab-btn.active .issd-tab-count {
  background: rgba(255, 255, 255, 0.28);
  color: #ffffff;
}
.issd-tab-count.is-urgent {
  background: #fef3c7;
  color: #d97706;
}
.issd-tab-btn.has-urgent-badge {
  border-color: #fcd34d;
  background: #fffbeb;
  color: #b45309;
}
.issd-body {
  flex: 1 1 auto;
  min-height: 0;
  overflow-y: auto;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  scrollbar-width: thin;
  scrollbar-color: rgba(148, 163, 184, 0.35) transparent;
}
.issd-body::-webkit-scrollbar,
.issd-scroll-list::-webkit-scrollbar {
  width: 5px;
}
.issd-body::-webkit-scrollbar-track,
.issd-scroll-list::-webkit-scrollbar-track {
  background: transparent;
}
.issd-body::-webkit-scrollbar-thumb,
.issd-scroll-list::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.35);
  border-radius: 999px;
}
.issd-body::-webkit-scrollbar-thumb:hover,
.issd-scroll-list::-webkit-scrollbar-thumb:hover {
  background: rgba(100, 116, 139, 0.55);
}
.issd-body-full {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 14px 16px;
}
.issd-body-full.is-fixed-tab-layout {
  overflow: hidden;
}
.issd-scroll-list {
  flex: 1 1 auto;
  min-height: 0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-right: 2px;
  scrollbar-width: thin;
  scrollbar-color: rgba(148, 163, 184, 0.35) transparent;
}
.issd-photo-banner.is-full{height:124px}
.issd-kpi-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:10px;margin:0;flex-shrink:0}

.issd-kpi-card{padding:10px 12px;border-radius:11px;border:1px solid #e2e8f0;background:#f8fafc;display:grid;gap:2px}
.issd-kpi-card span{color:#64748b;font-size:.75rem;font-weight:500}
.issd-kpi-card strong{color:#0f172a;font-size:1.15rem;font-weight:700;line-height:1.2}
.issd-kpi-card.is-green strong{color:#10b981}
.issd-kpi-card em{font-style:normal;color:#64748b;font-size:.75rem;font-weight:500}
.issd-fact-section{display:flex;flex-direction:column;gap:8px;flex-shrink:0}
.issd-section-title{font-size:.8125rem;font-weight:700;color:#1e293b}
.issd-fact-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:8px;margin:0}
.issd-fact-cell{padding:9px 12px;border:1px solid #e2e8f0;border-radius:9px;background:#fff}
.issd-fact-cell.full{grid-column:1/-1}
.issd-fact-cell dt{color:#64748b;font-size:.75rem;font-weight:500}
.issd-fact-cell dd{margin:3px 0 0;color:#1e293b;font-size:.8125rem;font-weight:600;word-break:break-word}
.issd-ops-snapshot{display:grid;gap:10px;padding:13px;border:1px solid #dbeafe;border-radius:12px;background:#f8fbff;flex-shrink:0}
.issd-ops-heading{display:flex;align-items:flex-start;justify-content:space-between;gap:12px}.issd-ops-heading>div{display:grid;gap:2px}.issd-ops-heading strong{font-size:.9375rem;color:#0f172a}.issd-ops-heading span:not(.issd-status-pill){font-size:.75rem;color:#64748b}
.issd-ops-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:1px;background:#e2e8f0;border:1px solid #e2e8f0;border-radius:9px;overflow:hidden}.issd-ops-grid>div{display:grid;gap:2px;padding:9px 10px;background:#fff}.issd-ops-grid span{font-size:.75rem;color:#64748b}.issd-ops-grid strong{font-size:.8125rem;color:#1e293b;line-height:1.35}.issd-ops-grid strong.is-alert{color:#dc2626}
.issd-action-summary{display:grid;gap:5px;padding:9px 10px;border-radius:9px;background:#fffbeb;border:1px solid #fde68a}.issd-action-summary>strong{font-size:.8125rem;color:#92400e}.issd-action-summary ul{display:grid;gap:3px;margin:0;padding-left:18px;color:#92400e;font-size:.75rem;line-height:1.45}

/* Action Task Cards in "สิ่งที่ต้องดำเนินการ" Tab */
.issd-action-cards-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.issd-action-task-card {
  padding: 13px 15px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  gap: 6px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}
.issd-action-task-card.is-danger {
  border-color: #fecaca;
  background: #fffafa;
}
.issd-action-task-card.is-warning {
  border-color: #fde68a;
  background: #fffdf7;
}
.issd-action-task-card.is-info {
  border-color: #bfdbfe;
  background: #f8fbff;
}
.issd-action-task-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.issd-action-cat-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 9px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}
.issd-action-cat-badge.is-danger {
  background: #fee2e2;
  color: #dc2626;
}
.issd-action-cat-badge.is-warning {
  background: #fef3c7;
  color: #d97706;
}
.issd-action-cat-badge.is-info {
  background: #eff6ff;
  color: #2563eb;
}
.issd-action-priority-text {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
}
.issd-action-task-title {
  font-size: 0.875rem;
  font-weight: 700;
  color: #0f172a;
}
.issd-action-task-desc {
  margin: 0;
  font-size: 0.8125rem;
  color: #475569;
  line-height: 1.45;
}
.issd-action-task-footer {
  display: flex;
  justify-content: flex-end;
  padding-top: 6px;
  border-top: 1px dashed #e2e8f0;
}
.issd-empty-box.is-clean-state {
  background: #f0fdf4;
  border-color: #bbf7d0;
}
.issd-empty-box.is-clean-state svg {
  color: #16a34a;
}
.issd-empty-box.is-clean-state p {
  color: #15803d;
  font-weight: 700;
}
.issd-empty-sub {
  font-size: 0.75rem;
  color: #64748b;
}

.issd-booking-list,.issd-contract-list{display:flex;flex-direction:column;gap:10px}
.issd-booking-card,.issd-contract-card{padding:14px 16px;border:1px solid #e2e8f0;border-radius:12px;background:#fff;display:flex;flex-direction:column;gap:10px;flex-shrink:0;transition:border-color .15s ease,box-shadow .15s ease}
.issd-contract-card.is-current{border-color:#93c5fd;box-shadow:inset 3px 0 0 #2563eb}.issd-current-label{display:inline-flex;padding:2px 7px;border-radius:6px;background:#eff6ff;color:#2563eb;font-size:.75rem;font-weight:600}
.issd-booking-card:hover,.issd-contract-card:hover{border-color:#cbd5e1;box-shadow:0 3px 10px rgba(15,23,42,.05)}
.issd-card-top-row{display:flex;align-items:center;justify-content:space-between;gap:8px}
.issd-card-type-tag{display:inline-flex;align-items:center;gap:5px;font-size:.75rem;font-weight:600;color:#2563eb;background:#eff6ff;padding:3px 9px;border-radius:6px}
.issd-card-type-tag.is-application{color:#7c3aed;background:#f5f3ff}
.issd-booking-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.issd-booking-type-group {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex-wrap: wrap;
}
.issd-booking-ref-code {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.8125rem;
  font-weight: 700;
  color: #334155;
}
.issd-booking-applicant-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}
.issd-booking-applicant-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}
.issd-booking-avatar {
  width: 34px;
  height: 34px;
  border-radius: 9px;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  color: #475569;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.issd-booking-applicant-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}
.issd-booking-role-label {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 500;
  line-height: 1.2;
}
.issd-booking-phone {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 500;
}
.issd-booking-footer-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 9px;
  background: #f8fafc;
  border: 1px solid #f1f5f9;
  flex-wrap: wrap;
}
.issd-booking-info-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  color: #64748b;
}
.issd-booking-info-item svg {
  color: #64748b;
  flex-shrink: 0;
}
.issd-booking-info-item strong {
  color: #1e293b;
  font-weight: 600;
}

/* Booking & Application Direct Action Button */
.issd-booking-direct-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 11px;
  border-radius: 7px;
  font-size: 0.75rem;
  font-weight: 600;
  text-decoration: none;
  border: 1px solid #bfdbfe;
  background: #ffffff;
  color: #1d4ed8;
  flex-shrink: 0;
  transition: all 0.15s ease;
}
.issd-booking-direct-btn:hover {
  background: #2563eb;
  color: #ffffff;
  border-color: #2563eb;
}
.issd-booking-direct-btn.is-application {
  border-color: #ddd6fe;
  background: #f5f3ff;
  color: #6d28d9;
}
.issd-booking-direct-btn.is-application:hover {
  background: #7c3aed;
  color: #ffffff;
  border-color: #7c3aed;
}

.issd-dual-btn-group {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}
.issd-dual-btn-group > * {
  flex: 1;
}

/* Mini Drawer Occupant Action Links */
.issd-occupant-action-strip {
  padding-top: 8px;
  border-top: 1px dashed #e2e8f0;
  display: flex;
}
.issd-occupant-action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 6px 10px;
  border-radius: 7px;
  font-size: 0.75rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.15s ease;
}
.issd-occupant-action-btn.is-app {
  background: #f5f3ff;
  border: 1px solid #ddd6fe;
  color: #6d28d9;
}
.issd-occupant-action-btn.is-app:hover {
  background: #ede9fe;
  border-color: #c4b5fd;
}
.issd-occupant-action-btn.is-booking {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #1d4ed8;
}
.issd-occupant-action-btn.is-booking:hover {
  background: #dbeafe;
  border-color: #93c5fd;
}

/* ==========================================================================
   Full Detail Action Alert Card (Modern High-End Task Card)
   ========================================================================== */
.issd-action-alert-card {
  margin: 10px 14px 4px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.08);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  overflow: hidden;
  transition: all 0.2s ease;
}

.issd-aac-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 9px 12px;
  background: #fffbeb;
  cursor: pointer;
  user-select: none;
  transition: background 0.15s ease;
}

.issd-aac-header:hover {
  background: #fef3c7;
}

.issd-aac-title-group {
  display: flex;
  align-items: center;
  gap: 9px;
  min-width: 0;
}

.issd-aac-icon-wrap {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: #fef3c7;
  color: #d97706;
  border: 1px solid #fde68a;
  display: grid;
  place-items: center;
  flex-shrink: 0;
}

.issd-aac-title-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.issd-aac-title-row {
  display: flex;
  align-items: center;
  gap: 7px;
  flex-wrap: wrap;
}

.issd-aac-title-row strong {
  font-size: 0.84375rem;
  font-weight: 700;
  color: #92400e;
}

.issd-aac-count-badge {
  display: inline-flex;
  align-items: center;
  padding: 1px 7px;
  border-radius: 999px;
  background: #fde68a;
  color: #78350f;
  font-size: 0.75rem;
  font-weight: 700;
}

.issd-aac-sub {
  font-size: 0.75rem;
  color: #b45309;
}

.issd-aac-toggle-btn {
  width: 26px;
  height: 26px;
  border-radius: 6px;
  border: 1px solid rgba(217, 119, 6, 0.25);
  background: #ffffff;
  color: #92400e;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.issd-aac-toggle-btn:hover {
  background: #fef3c7;
  border-color: #d97706;
}

.issd-aac-toggle-btn svg {
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.issd-aac-toggle-btn svg.is-flipped {
  transform: rotate(-90deg);
}

.issd-aac-list {
  display: flex;
  flex-direction: column;
  background: #ffffff;
}

.issd-aac-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  border-top: 1px solid #f1f5f9;
  cursor: pointer;
  transition: background 0.15s ease;
}

.issd-aac-row:hover {
  background: #f8fafc;
}

.issd-aac-row:hover .issd-aac-action-text {
  color: #1d4ed8;
}

.issd-aac-row:hover .issd-aac-row-right svg {
  transform: translateX(2px);
}

.issd-aac-row.is-danger:hover {
  background: #fef2f2;
}

.issd-aac-row.is-warning:hover {
  background: #fffbeb;
}

.issd-aac-row.is-info:hover {
  background: #eff6ff;
}

.issd-aac-row-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex: 1;
}

.issd-aac-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 7px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 700;
  white-space: nowrap;
  flex-shrink: 0;
}

.issd-aac-tag.is-danger {
  background: #fee2e2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.issd-aac-tag.is-warning {
  background: #fef3c7;
  color: #d97706;
  border: 1px solid #fde68a;
}

.issd-aac-tag.is-info {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
}

.issd-aac-row-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.issd-aac-row-right {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  flex-shrink: 0;
  color: #2563eb;
  font-size: 0.78125rem;
  font-weight: 600;
}

.issd-aac-row-right svg {
  transition: transform 0.15s ease;
}

.issd-aac-action-text {
  white-space: nowrap;
}
.issd-card-title,.issd-contract-num{font-size:.875rem;font-weight:700;color:#0f172a}
.issd-contract-tenant{font-size:.8125rem;color:#475569}
.issd-contract-tenant strong{color:#1e293b;font-weight:600}
.issd-card-meta-row{display:flex;align-items:center;justify-content:space-between;gap:10px;font-size:.75rem;color:#64748b;flex-wrap:wrap}
.issd-card-meta-row .is-rent{color:#10b981;font-weight:700;font-size:.8125rem}
.issd-timeline-list{display:flex;flex-direction:column;padding:4px 4px 10px}
.issd-timeline-panel{padding:12px 14px 4px;border:1px solid #e2e8f0;border-radius:12px;background:#fff}
.issd-timeline-item{position:relative;display:flex;gap:12px;padding-bottom:16px}
.issd-timeline-item:not(:last-child)::before{content:'';position:absolute;left:5px;top:13px;bottom:0;width:1px;background:#e2e8f0}
.issd-timeline-dot{position:relative;z-index:1;width:11px;height:11px;margin-top:4px;flex:0 0 11px;border-radius:50%;background:#3b82f6;box-shadow:0 0 0 3px #dbeafe}
.issd-timeline-dot.is-green{background:#10b981;box-shadow:0 0 0 3px #dcfce7}
.issd-timeline-dot.is-amber{background:#f59e0b;box-shadow:0 0 0 3px #fef3c7}
.issd-timeline-dot.is-gray{background:#64748b;box-shadow:0 0 0 3px #f1f5f9}
.issd-timeline-copy{flex:1;min-width:0}
.issd-timeline-head-row{display:flex;align-items:center;justify-content:space-between;gap:8px;flex-wrap:wrap}
.issd-timeline-cat-pill{display:inline-flex;padding:2px 8px;border-radius:6px;background:#f1f5f9;color:#475569;font-size:.75rem;font-weight:600}
.issd-timeline-copy strong{display:block;font-size:.8125rem;color:#1e293b;font-weight:600}
.issd-timeline-copy time{display:block;color:#64748b;font-size:.75rem;margin-top:2px}
.issd-timeline-copy p{margin:6px 0 0;padding:7px 10px;border-radius:8px;background:#f8fafc;border:1px solid #f1f5f9;color:#475569;font-size:.75rem;line-height:1.45}
.issd-empty-box{display:flex;flex-direction:column;align-items:center;justify-content:center;gap:8px;padding:32px 16px;text-align:center;color:#94a3b8;background:#f8fafc;border-radius:12px;border:1px dashed #e2e8f0}
.issd-empty-box svg{color:#cbd5e1}
.issd-empty-box p{margin:0;font-size:.8125rem;color:#64748b}
.issd-body-edit{padding:12px 14px}
.issd-form-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:10px;margin:0}
.issd-form-grid .full{grid-column:1/-1}
.issd-form-grid label{display:grid;gap:5px;color:#475569;font-size:.8125rem;font-weight:600}
.issd-form-grid textarea{width:100%;box-sizing:border-box;border:1px solid #cbd5e1;border-radius:8px;padding:8px 10px;color:#1e293b;background:#fff;font:inherit;resize:vertical;min-height:56px}
.issd-form-grid textarea:focus{outline:none;border-color:#2563eb;box-shadow:0 0 0 3px rgba(37,99,235,.12)}
.issd-footer-edit{display:flex;align-items:center;justify-content:space-between;gap:8px}
.issd-footer-right{display:flex;align-items:center;gap:8px;margin-left:auto}
@media(max-width:480px){.issd-form-grid{grid-template-columns:1fr}}
.issd-status-pill{display:inline-flex;align-items:center;padding:3px 10px;border-radius:999px;font-size:.75rem;font-weight:600;border:1px solid transparent;white-space:nowrap}.issd-status-pill.is-available{background:#dcfce7;color:#15803d;border-color:#bbf7d0}.issd-status-pill.is-reserved{background:#fef3c7;color:#d97706;border-color:#fde68a}.issd-status-pill.is-info{background:#eff6ff;color:#2563eb;border-color:#bfdbfe}.issd-status-pill.is-occupied{background:#f1f5f9;color:#475569;border-color:#cbd5e1}
.issd-photo-banner{position:relative;height:124px;flex-shrink:0;border-radius:12px;overflow:hidden;border:1px solid #e2e8f0;background:#0f172a;cursor:pointer}.issd-photo-banner img{width:100%;height:100%;object-fit:cover;display:block;transition:transform .25s ease}.issd-photo-banner:hover img{transform:scale(1.05)}.issd-photo-count-badge{position:absolute;right:10px;bottom:10px;display:inline-flex;align-items:center;gap:5px;padding:3px 9px;border-radius:8px;background:rgba(15,23,42,.82);backdrop-filter:blur(4px);color:#fff;font-size:.75rem;font-weight:600}
.issd-detail-list{margin:0;flex-shrink:0;border:1px solid #e2e8f0;border-radius:12px;background:#fff;overflow:hidden}.issd-detail-row{display:flex;align-items:center;justify-content:space-between;gap:14px;padding:10px 14px;font-size:.8125rem}.issd-detail-row+.issd-detail-row{border-top:1px solid #f1f5f9}.issd-detail-row dt{color:#64748b;font-size:.8125rem;font-weight:500;flex-shrink:0;white-space:nowrap}.issd-detail-row dd{margin:0;color:#1e293b;font-size:.8125rem;font-weight:600;text-align:right;max-width:65%;line-height:1.4}.issd-detail-row dd.is-rent{color:#15803d;font-weight:700}.issd-detail-row dd.is-bold{color:#0f172a;font-weight:700}
.issd-occupancy-box{padding:11px 14px;flex-shrink:0;border-radius:12px;border:1px solid #e2e8f0;background:#f8fafc;display:grid;gap:3px}.issd-occupancy-box.is-available{background:#f0fdf4;border-color:#bbf7d0}.issd-occupancy-box.is-reserved{background:#fffbeb;border-color:#fde68a}.issd-occupancy-box.is-occupied{background:#f8fafc;border-color:#cbd5e1}.issd-occupancy-head{display:flex;align-items:center;justify-content:space-between;gap:8px}.issd-occupancy-head span{font-size:.75rem;font-weight:600;color:#64748b}.issd-occupancy-head small{font-size:.75rem;font-weight:700;color:#2563eb}.issd-occupancy-box strong{font-size:.8125rem;font-weight:700;color:#0f172a;line-height:1.35}.issd-occupancy-box p{margin:0;font-size:.75rem;color:#475569;line-height:1.35}
.issd-footer{padding:12px 18px 16px;background:#fff;border-top:1px solid #f1f5f9;box-shadow:0 -4px 12px rgba(15,23,42,.04);flex-shrink:0}.issd-cta-btn{width:100%;box-sizing:border-box;padding:11px 16px;border:0;border-radius:11px;background:linear-gradient(180deg,#3b82f6 0%,#2563eb 100%);color:#fff;font-family:inherit;font-size:.875rem;font-weight:700;text-decoration:none;display:inline-flex;align-items:center;justify-content:center;gap:8px;cursor:pointer;box-shadow:0 4px 14px rgba(37,99,235,.32);transition:transform .15s ease,box-shadow .15s ease}.issd-cta-btn:hover{transform:translateY(-1px);box-shadow:0 6px 18px rgba(37,99,235,.4)}
.issd-footer-full{display:flex;align-items:center;justify-content:space-between;gap:10px}
.issd-footer-full :deep(.app-button),
.issd-footer-full > button,
.issd-footer-full > a {
  flex: 1;
}
.inspect-photo-Modal,.inspect-photo-modal{width:min(820px,94vw);background:#0f172a;border-radius:18px;overflow:hidden;box-shadow:0 30px 60px -20px rgba(0,0,0,.6);display:flex;flex-direction:column}.ipm-head{display:flex;align-items:center;justify-content:space-between;padding:14px 18px;color:#fff}.ipm-head strong{display:block;font-size:1rem;font-weight:700}.ipm-head span{display:block;font-size:.8125rem;color:#94a3b8;margin-top:2px}.ipm-close{width:34px;height:34px;border-radius:10px;border:0;background:rgba(255,255,255,.12);color:#fff;display:inline-flex;align-items:center;justify-content:center;cursor:pointer}.ipm-close svg{width:18px;height:18px}.ipm-stage{position:relative;height:min(480px,58vh);background:#000;display:flex;align-items:center;justify-content:center}.ipm-main-img{width:100%;height:100%;object-fit:contain}.ipm-nav{position:absolute;top:50%;transform:translateY(-50%);width:40px;height:40px;border-radius:50%;border:0;background:rgba(255,255,255,.9);color:#0f172a;font-size:1.35rem;font-weight:700;cursor:pointer;display:grid;place-items:center}.ipm-nav.is-prev{left:14px}.ipm-nav.is-next{right:14px}.ipm-counter{position:absolute;bottom:12px;left:50%;transform:translateX(-50%);background:rgba(15,23,42,.78);color:#fff;font-size:.75rem;font-weight:600;padding:3px 12px;border-radius:999px}.ipm-thumbs{display:flex;justify-content:center;gap:8px;padding:12px 18px 16px}.ipm-thumb{width:76px;height:52px;border-radius:8px;overflow:hidden;border:2px solid transparent;padding:0;background:#1e293b;cursor:pointer;opacity:.6}.ipm-thumb.active{border-color:#60a5fa;opacity:1}.ipm-thumb img{width:100%;height:100%;object-fit:cover;display:block}
@media(max-width:768px){.inspect-selected-stall-drawer{position:fixed;left:0;right:0;bottom:0;top:58px;width:100%!important;height:calc(100vh - 58px);max-height:calc(100vh - 58px);border-radius:0;z-index:1000}.inspect-selected-stall-drawer.is-full-detail-mode{width:100%!important;top:58px;height:calc(100vh - 58px);max-height:calc(100vh - 58px);border-radius:0}}
.map-edit-status{display:flex;gap:10px;align-items:flex-start;padding:13px;background:#eff6ff;border:1px solid #bfdbfe;border-radius:10px;color:#1e40af}.map-edit-status strong,.map-edit-status small{display:block}.map-edit-status strong{font-size:.875rem}.map-edit-status small{margin-top:3px;font-size:.75rem;color:#476587}.map-edit-dot{width:10px;height:10px;flex:0 0 10px;border-radius:50%;background:#3b5bff;margin-top:5px;box-shadow:0 0 0 4px #dbeafe}
.zone-form,.facility-form{display:grid;grid-template-columns:1fr 1fr;gap:13px}.zone-form label,.facility-form label{display:grid;align-content:start;gap:6px;color:#475569;font-size:.8125rem;font-weight:600}.zone-form .full,.facility-form .full{grid-column:1/-1}.req-star{color:#ef4444;font-weight:700}.stall-detail-summary{display:grid;gap:5px;padding:16px;margin:20px 0 14px;border:1px solid #dbeafe;border-radius:12px;background:#f8fbff}.stall-detail-summary strong{color:#1e293b;font-size:1rem}.stall-detail-summary>span:last-child{color:#64748b;font-size:.875rem}.stall-detail-summary .soft-badge{justify-self:start}.status-available{background:#dcfce7;color:#15803d!important}.status-reserved{background:#fef3c7;color:#d97706!important}.status-occupied{background:#f1f5f9;color:#64748b!important}.stall-detail-tabs{display:flex!important;width:100%;margin:0 0 14px}.stall-detail-tabs :deep(.app-segmented-option){flex:1}.stall-detail-grid{display:grid;grid-template-columns:1fr 1fr;gap:10px;margin:0}.stall-detail-grid>div{padding:12px 14px;border:1px solid #e2e8f0;border-radius:10px;background:#fff}.stall-detail-grid .full{grid-column:1/-1}.stall-detail-grid dt{color:#64748b;font-size:.8125rem;font-weight:500}.stall-detail-grid dd{margin:4px 0 0;color:#1e293b;font-size:.9375rem;font-weight:600;word-break:break-word}.detail-section.is-tab-content{margin-top:0}.detail-section-heading{display:flex;align-items:center;justify-content:space-between;gap:12px}.detail-section-heading h3{margin:0}
.issd-tab-header-row {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-bottom: 2px;
  flex-shrink: 0;
}
.issd-tab-header-row.is-inline-controls {
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.issd-tab-header-copy {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.issd-contract-inline-filter {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
.issd-inline-filter-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  white-space: nowrap;
}
.issd-inline-filter-select {
  width: 195px;
}
.issd-tab-header-row.is-inline-controls {
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.issd-tab-header-copy {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.issd-contract-inline-filter {
  width: 210px;
  flex-shrink: 0;
}
.issd-tab-title {
  display: block;
  font-size: 0.9375rem;
  font-weight: 700;
  color: #0f172a;
}
.issd-tab-subtitle {
  display: block;
  font-size: 0.75rem;
  color: #64748b;
  margin-top: 1px;
}

.issd-subtenant-list,
.issd-payment-list,
.issd-maintenance-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.issd-subtenant-card,
.issd-payment-card,
.issd-maintenance-card {
  padding: 12px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
  transition: border-color .15s ease, box-shadow .15s ease;
}
.issd-subtenant-card:hover,
.issd-payment-card:hover,
.issd-maintenance-card:hover {
  border-color: #cbd5e1;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}
.issd-contract-context{display:grid;grid-template-columns:auto 1fr;gap:2px 10px;align-items:center;padding:10px 12px;border:1px solid #dbeafe;border-radius:10px;background:#f8fbff;flex-shrink:0}.issd-contract-context span{font-size:.75rem;color:#64748b}.issd-contract-context strong{font-size:.8125rem;color:#2563eb}.issd-contract-context small{grid-column:2;font-size:.75rem;color:#475569}
.issd-finance-summary{display:grid;grid-template-columns:repeat(3,1fr);gap:8px;flex-shrink:0}.issd-finance-summary>div{display:grid;gap:2px;padding:9px 10px;border:1px solid #e2e8f0;border-radius:9px;background:#fff}.issd-finance-summary span{font-size:.75rem;color:#64748b}.issd-finance-summary strong{font-size:.8125rem;color:#1e293b}.issd-finance-summary .is-success{color:#15803d}.issd-finance-summary .is-alert{color:#dc2626}
.issd-source-filter{display:flex;align-items:center;gap:6px;padding-bottom:2px;overflow-x:auto;flex-shrink:0}.issd-source-filter button{min-height:34px;padding:6px 12px;border:1px solid #e2e8f0;border-radius:8px;background:#fff;color:#64748b;font:inherit;font-size:.8125rem;font-weight:600;white-space:nowrap;cursor:pointer;flex-shrink:0;transition:all .15s ease}.issd-source-filter button:hover{background:#f8fafc;color:#1e293b}.issd-source-filter button.active{border-color:#93c5fd;background:#eff6ff;color:#2563eb}.issd-maintenance-tags{display:flex;align-items:center;gap:6px;flex-wrap:wrap}.issd-source-pill{display:inline-flex;padding:2px 8px;border-radius:6px;background:#f1f5f9;color:#475569;font-size:.75rem;font-weight:600}

.issd-subtenant-header-info {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}
.issd-subtenant-avatar {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  overflow: hidden;
  flex-shrink: 0;
  border: 1px solid #e2e8f0;
  background: #f1f5f9;
}
.issd-subtenant-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.issd-subtenant-avatar.is-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  background: #f8fafc;
}
.issd-subtenant-text-block {
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.issd-subtenant-idcard {
  font-size: 0.75rem;
  color: #64748b;
  margin-top: 1px;
}

/* Photo Upload in Modal */
.issd-photo-upload-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.issd-upload-label {
  font-size: 0.8125rem;
  font-weight: 600;
  color: #475569;
}
.issd-photo-upload-box {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px;
  border-radius: 10px;
  border: 1px dashed #cbd5e1;
  background: #f8fafc;
}
.issd-preview-thumb {
  position: relative;
  width: 56px;
  height: 56px;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  flex-shrink: 0;
}
.issd-preview-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.issd-remove-photo-btn {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: rgba(15, 23, 42, 0.75);
  color: #ffffff;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  transition: background 0.15s ease;
}
.issd-remove-photo-btn:hover {
  background: #dc2626;
}
.issd-empty-thumb {
  width: 56px;
  height: 56px;
  border-radius: 10px;
  background: #f1f5f9;
  border: 1px dashed #cbd5e1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  flex-shrink: 0;
}
.issd-upload-actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.issd-upload-btn-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.issd-file-choose-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 0.8125rem;
  font-weight: 600;
  color: #1e293b;
  cursor: pointer;
  transition: all 0.15s ease;
}
.issd-file-choose-btn:hover {
  background: #f1f5f9;
  border-color: #94a3b8;
}
.issd-remove-text-btn {
  border: none;
  background: transparent;
  color: #ef4444;
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  padding: 4px 6px;
}
.issd-remove-text-btn:hover {
  text-decoration: underline;
}
.issd-file-hidden-input {
  display: none;
}
.issd-upload-hint {
  font-size: 0.75rem;
  color: #64748b;
}

.issd-subtenant-details {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 0.8125rem;
  color: #475569;
}
.issd-detail-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: #475569;
}
.issd-detail-item svg {
  color: #64748b;
  flex-shrink: 0;
}
.issd-subtenant-note {
  font-size: 0.75rem;
  color: #64748b;
  background: #f8fafc;
  padding: 6px 10px;
  border-radius: 8px;
  border: 1px solid #f1f5f9;
  line-height: 1.4;
}

.issd-card-actions-row {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding-top: 6px;
  border-top: 1px dashed #f1f5f9;
}
.issd-mini-action-btn {
  border: 1px solid #e2e8f0;
  background: #fff;
  color: #475569;
  font-family: inherit;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all .15s ease;
}
.issd-mini-action-btn:hover {
  background: #f8fafc;
  color: #0f172a;
  border-color: #cbd5e1;
}
.issd-mini-action-btn.is-success {
  color: #15803d;
  border-color: #bbf7d0;
  background: #f0fdf4;
}
.issd-mini-action-btn.is-success:hover {
  background: #dcfce7;
  border-color: #86efac;
}
.issd-mini-action-btn.is-danger {
  color: #ef4444;
  border-color: #fecaca;
  background: #fff;
}
.issd-mini-action-btn.is-danger:hover {
  background: #fef2f2;
  border-color: #f87171;
}

.issd-contract-top-left {
  display: flex;
  align-items: center;
  gap: 8px;
}
.issd-contract-icon {
  color: #2563eb;
  flex-shrink: 0;
}
.issd-contract-details-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 6px 12px;
  background: #f8fafc;
  padding: 8px 10px;
  border-radius: 8px;
  margin-top: 4px;
}
.issd-contract-detail-item {
  display: flex;
  align-items: baseline;
  gap: 6px;
  font-size: 0.75rem;
}
.issd-contract-detail-item.full-span {
  grid-column: 1 / -1;
}
.issd-contract-detail-item .detail-label {
  color: #64748b;
  font-weight: 500;
  flex-shrink: 0;
}
.issd-contract-detail-item .detail-value {
  color: #1e293b;
  font-weight: 600;
}
.issd-contract-detail-item .detail-value.is-rent {
  color: #15803d;
  font-weight: 700;
}
.issd-mini-action-btn.is-primary {
  color: #2563eb;
  border-color: #bfdbfe;
  background: #eff6ff;
}
.issd-mini-action-btn.is-primary:hover {
  background: #dbeafe;
  border-color: #93c5fd;
}
.issd-status-pill.is-cancelled {
  background: #fee2e2;
  color: #dc2626;
  border-color: #fecaca;
}

.issd-pay-period strong {
  display: block;
  font-size: 0.875rem;
  font-weight: 700;
  color: #0f172a;
}
.issd-pay-period span {
  display: block;
  font-size: 0.75rem;
  color: #64748b;
}
.issd-pay-amount-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin: 2px 0;
}
.issd-pay-method {
  font-size: 0.8125rem;
  color: #475569;
}

.issd-maintenance-type-pill {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}
.issd-maintenance-type-pill.type-repair {
  background: #fef2f2;
  color: #dc2626;
}
.issd-maintenance-type-pill.type-renovation {
  background: #eff6ff;
  color: #2563eb;
}
.issd-maintenance-type-pill.type-utility {
  background: #fef3c7;
  color: #d97706;
}
.issd-maintenance-type-pill.type-other {
  background: #f1f5f9;
  color: #475569;
}
.issd-maintenance-desc {
  margin: 2px 0 4px;
  font-size: 0.8125rem;
  color: #475569;
  line-height: 1.45;
}
.issd-maintenance-meta-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 4px 10px;
  font-size: 0.75rem;
  color: #64748b;
  background: #f8fafc;
  padding: 8px 10px;
  border-radius: 8px;
  margin: 2px 0;
}
.issd-maintenance-meta-grid strong {
  color: #1e293b;
  font-weight: 600;
}
@media(max-width:480px){.issd-ops-grid,.issd-finance-summary{grid-template-columns:1fr}.issd-contract-context{grid-template-columns:1fr}.issd-contract-context small{grid-column:auto}}

.issd-modal-form {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}
.issd-modal-form label {
  display: grid;
  align-content: start;
  gap: 6px;
  color: #475569;
  font-size: 0.8125rem;
  font-weight: 600;
}
.issd-modal-form .full {
  grid-column: 1 / -1;
}
.issd-textarea {
  width: 100%;
  box-sizing: border-box;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 12px;
  font-family: inherit;
  font-size: 0.875rem;
  color: #1e293b;
  background: #fff;
  resize: vertical;
  outline: none;
  min-height: 64px;
  transition: border-color .15s ease, box-shadow .15s ease;
}
.issd-textarea:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
}
@media (max-width: 580px) {
  .issd-modal-form {
    grid-template-columns: 1fr;
  }
}

.manage-page-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin: 8px 0 16px;
  flex-wrap: wrap;
}

.manage-page-heading .page-title {
  margin: 0;
  color: #1e293b;
  font-size: 1.5rem;
  font-weight: 800;
  line-height: 1.3;
}

.manage-page-heading .page-desc {
  margin: 4px 0 0;
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
  line-height: 1.5;
}

.manage-page-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.panel,
.gis-panel {
  background: #ffffff;
  border: 1px solid #eef2f6;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
  padding: 24px 28px !important;
}

.panel-heading,
.gis-panel-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px !important;
  flex-wrap: wrap;
}

.gis-panel-header-left {
  flex: 0 1 auto;
  min-width: 0;
}

.gis-panel-header-left h2 {
  font-size: 1.125rem !important;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  line-height: 1.3;
}

.gis-panel-header-left p {
  font-size: 0.8125rem !important;
  font-weight: 500;
  color: #64748b;
  margin: 4px 0 0 !important;
  line-height: 1.4;
}

.gis-heading-right {
  display: flex !important;
  align-items: center !important;
  gap: 10px !important;
  margin-left: auto !important;
  flex-shrink: 0 !important;
}

.gis-map-toolbar {
  display: flex !important;
  align-items: center !important;
  flex-shrink: 0 !important;
}

.gis-map-toolbar :deep(.master-map-top-bar) {
  width: auto !important;
  margin: 0 !important;
  display: flex !important;
  align-items: center !important;
}

.gis-map-toolbar :deep(.master-map-tools-right) {
  display: flex !important;
  align-items: center !important;
  gap: 8px !important;
  margin-left: 0 !important;
  flex-wrap: nowrap !important;
}

.gis-map-toolbar :deep(.search-input-group) {
  width: 250px !important;
  min-width: 140px !important;
  max-width: 280px !important;
  flex-shrink: 1 !important;
}

.gis-map-toolbar :deep(.map-advanced-filter-wrap) {
  flex-shrink: 0 !important;
  display: inline-flex !important;
  position: relative !important;
}

.gis-map-toolbar :deep(.map-advanced-filter-btn) {
  display: inline-flex !important;
  align-items: center !important;
  gap: 6px !important;
  height: 38px !important;
  padding: 6px 14px !important;
  white-space: nowrap !important;
  flex-shrink: 0 !important;
  visibility: visible !important;
  opacity: 1 !important;
}

.manage-page{padding-bottom:48px}
.tool-panel{align-self:start;position:sticky;top:86px}
.vertex-list,.form-stack{display:grid;gap:10px}
.vertex-row{display:grid;grid-template-columns:24px 1fr 1fr 30px;gap:6px;align-items:center}
.vertex-row span{font-size:.75rem;color:#64748b;text-align:center}
.vertex-row button,.facility-list button,.icon-button{border:0;background:transparent;color:#64748b;cursor:pointer;padding:5px}
.vertex-row svg,.facility-list svg,.icon-button svg{width:18px;height:18px}
.form-stack label,.form-grid label{display:grid;gap:6px;color:#475569;font-size:.8125rem;font-weight:600}
.coordinate-grid{display:grid;grid-template-columns:1fr 1fr;gap:8px}
.facility-list{display:grid;gap:7px;margin-top:15px}
.facility-list>div{display:flex;justify-content:space-between;padding:9px 10px;background:#f8fafc;border-radius:9px}
.facility-list span{display:grid}
.facility-list small{font-size:.75rem;color:#64748b}
.check-list{display:grid;gap:8px}
.check-list>div{display:flex;gap:10px;padding:12px;background:#f8fafc;border-radius:10px}
.check-list svg{width:20px;color:#16a34a}
.check-list span{display:grid}
.check-list strong{font-size:.875rem}
.check-list small{font-size:.75rem;color:#64748b}
.stacked{display:grid;margin-top:16px}
.empty-state{padding:20px;text-align:center;color:#64748b;background:#f8fafc;border-radius:10px;font-size:.875rem}
.compact{padding:12px}

.form-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:13px;margin:20px 0}
.form-grid .full{grid-column:1/-1}
input,textarea{width:100%;box-sizing:border-box;border:1px solid #cbd5e1;border-radius:8px;padding:9px 10px;color:#1e293b;background:#fff;font:inherit}
.detail-section{border-top:1px solid #e2e8f0;padding-top:17px;margin-top:17px}
.detail-section h3{font-size:.9375rem;color:#1e293b}
.contract-list{display:grid;gap:8px}
.contract-list article{display:flex;justify-content:space-between;gap:10px;padding:11px;background:#f8fafc;border-radius:9px}
.contract-list article div{display:grid}
.contract-list span{font-size:.8125rem;color:#64748b}
.soft-badge{background:#eff6ff;color:#2563eb!important;padding:4px 8px;border-radius:7px;height:max-content}
.history-list{display:grid}
.history-list>div{display:flex;gap:10px;padding:6px 0}
.history-list i{width:9px;height:9px;border-radius:50%;background:#3b82f6;margin-top:5px}
.history-list span{display:grid}
.history-list small{font-size:.75rem;color:#64748b}

@media (max-width: 768px) {
  .manage-page-heading {
    align-items: flex-start;
    flex-direction: column;
    gap: 12px;
  }
  .manage-page-actions {
    width: 100%;
    justify-content: space-between;
  }
}

/* แถบสลับโหมดหน้าข้อมูลตลาด (ดูข้อมูล / จัดการค่าใช้จ่าย) */
.manage-page.is-map-page { --mp-top: 114px; --mp-top-2row: 190px; }
/* สวิตช์มุมมองแผนที่แบบ Capsule Pill ลอยมุมซ้ายบน */
.map-mode-dropdown-wrap {
  position: absolute;
  top: 18px;
  left: 18px;
  z-index: 125;
}
.map-mode-trigger-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 48px;
  padding: 0 14px 0 8px;
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(226, 232, 240, 0.95);
  border-radius: 9999px;
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  text-align: left;
  font-family: inherit;
}
.map-mode-trigger-btn:hover {
  background: #ffffff;
  border-color: #cbd5e1;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.12);
  transform: translateY(-1px);
}
.map-mode-trigger-btn.is-open {
  background: #ffffff;
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.16), 0 10px 28px rgba(15, 23, 42, 0.12);
}
.trigger-icon-box {
  width: 32px;
  height: 32px;
  border-radius: 9999px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s ease;
}
.trigger-icon-box.tone-blue { background: #eff6ff; color: #2563eb; }
.trigger-icon-box.tone-emerald { background: #ecfdf5; color: #059669; }
.trigger-icon-box.tone-amber { background: #fffbeb; color: #d97706; }

.trigger-text-group {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}
.trigger-eyebrow {
  font-size: 0.8125rem;
  color: #64748b;
  font-weight: 500;
  line-height: 1;
}
.trigger-title {
  font-size: 0.875rem;
  color: #0f172a;
  font-weight: 700;
  line-height: 1;
}
.trigger-chevron {
  color: #64748b;
  transition: transform 0.2s ease;
  margin-left: 2px;
}
.trigger-chevron.is-open {
  transform: rotate(180deg);
  color: #2563eb;
}

/* Dropdown Menu */
.map-mode-dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  width: 320px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  box-shadow: 0 12px 32px rgba(15, 23, 42, 0.16);
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  z-index: 140;
}
.map-mode-menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 12px;
  border: 0;
  background: transparent;
  cursor: pointer;
  text-align: left;
  transition: background 0.15s ease;
  width: 100%;
  font-family: inherit;
}
.map-mode-menu-item:hover {
  background: #f8fafc;
}
.map-mode-menu-item.is-selected {
  background: #eff6ff;
}
.item-icon-box {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.item-icon-box.tone-blue { background: #eff6ff; color: #2563eb; }
.item-icon-box.tone-emerald { background: #ecfdf5; color: #059669; }
.item-icon-box.tone-amber { background: #fffbeb; color: #d97706; }
.map-mode-menu-item.is-selected .item-icon-box.tone-blue { background: #2563eb; color: #ffffff; }
.map-mode-menu-item.is-selected .item-icon-box.tone-emerald { background: #059669; color: #ffffff; }
.map-mode-menu-item.is-selected .item-icon-box.tone-amber { background: #d97706; color: #ffffff; }

.item-text-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}
.item-title {
  font-size: 0.9375rem;
  color: #0f172a;
  font-weight: 700;
}
.map-mode-menu-item.is-selected .item-title {
  color: #1d4ed8;
}
.item-desc {
  font-size: 0.8125rem;
  color: #64748b;
  font-weight: 500;
}
.item-check {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-left: auto;
}

.dropdown-pop-fade-enter-active,
.dropdown-pop-fade-leave-active {
  transition: opacity 0.18s ease, transform 0.18s cubic-bezier(0.16, 1, 0.3, 1);
}
.dropdown-pop-fade-enter-from,
.dropdown-pop-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.98);
}

@media (max-width: 1000px) {
  .manage-page .officer-summary-dashboard.is-map-floating { top: var(--mp-top-2row); }
}
@media (max-width: 720px) {
  .map-mode-dropdown-wrap { top: 10px; left: 10px; }
  .manage-page.is-map-page { --mp-top: 148px; --mp-top-2row: 148px; }
}

.market-map-kpi.tone-blue .market-map-kpi-icon { background: #eff6ff; color: #2563eb; }
.market-map-kpi.tone-green .market-map-kpi-icon { background: #ecfdf5; color: #059669; }
.market-map-kpi.tone-amber .market-map-kpi-icon { background: #fffbeb; color: #d97706; }
.market-map-kpi.tone-red .market-map-kpi-icon { background: #fef2f2; color: #dc2626; }
.market-map-kpi.tone-rose .market-map-kpi-icon { background: #fff1f2; color: #e11d48; }
.market-map-kpi.tone-violet .market-map-kpi-icon { background: #f5f3ff; color: #7c3aed; }

.dash-breakdown-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 10px;
}
.dash-breakdown-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}
.dash-breakdown-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #334155;
  font-weight: 500;
}
.dot-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dash-breakdown-val {
  display: flex;
  align-items: baseline;
  gap: 6px;
}
.dash-breakdown-val strong {
  color: #0f172a;
  font-weight: 700;
}
.dash-breakdown-val small {
  color: #64748b;
  font-size: 0.8125rem;
}
.dash-progress-bar {
  width: 100%;
  height: 6px;
  background: #f1f5f9;
  border-radius: 9999px;
  overflow: hidden;
}
.dash-progress-fill {
  height: 100%;
  border-radius: 9999px;
  transition: width 0.4s ease;
}

.dash-mini-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
  margin-top: 8px;
}
.dash-mini-table th {
  padding: 8px 10px;
  text-align: left;
  font-size: 0.8125rem;
  font-weight: 600;
  color: #64748b;
  border-bottom: 1px solid #f1f5f9;
  background: #f8fafc;
}
.dash-mini-table td {
  padding: 10px 10px;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}
.table-stall-pill {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 6px;
  background: #eff6ff;
  color: #2563eb;
  font-weight: 700;
  font-size: 0.8125rem;
}
.table-tenant-name {
  color: #1e293b;
  font-weight: 600;
  font-size: 0.875rem;
}
.table-tenant-sub {
  color: #64748b;
  font-size: 0.75rem;
}
.table-status-pill {
  display: inline-block;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}
.table-status-pill.pill-overdue {
  background: #fee2e2;
  color: #dc2626;
}
.table-status-pill.pill-anomaly {
  background: #fef2f2;
  color: #b91c1c;
}
.table-type-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}
.table-type-badge.badge-water {
  background: #eff6ff;
  color: #1d4ed8;
}
.table-type-badge.badge-electric {
  background: #fef3c7;
  color: #b45309;
}

/* ปุ่มและแผงแดชบอร์ดลอยบนแผนที่ (อยู่เหนือปุ่มคำอธิบายสีโครงการ) */
.map-dashboard-fab {
  position: absolute;
  left: 18px;
  bottom: 18px;
  z-index: 46;
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  padding: 0;
  border: 2px solid #ffffff;
  border-radius: 12px;
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 4px 18px rgba(37, 99, 235, 0.45), 0 2px 6px rgba(15, 23, 42, 0.2);
  cursor: pointer;
  transition: transform 0.18s ease, background 0.18s ease;
}
.map-dashboard-fab:hover { transform: translateY(-1px); background: #1d4ed8; }
.map-dashboard-fab.active { background: #1e3a8a; box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.3), 0 4px 14px rgba(15, 23, 42, 0.25); }
.map-dashboard-fab:focus-visible { outline: 3px solid #93c5fd; outline-offset: 2px; }

.officer-summary-dashboard.is-map-floating {
  position: absolute;
  top: var(--mp-top); /* ใต้การ์ด KPI */
  left: 18px;
  bottom: 18px; /* ชิดล่าง ปุ่มแดชบอร์ดเลื่อนไปอยู่ข้างแผง */
  z-index: 60;
  grid-template-columns: minmax(0, 1fr);
  box-sizing: border-box;
  width: min(520px, calc(100% - 36px));
  max-width: 100%;
  margin: 0;
  padding: 0 16px 16px;
  gap: 14px;
  align-content: start;
  overflow-x: hidden;
  overflow-y: auto;
  background: rgba(255, 255, 255, 0.98);
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  box-shadow: 0 18px 44px rgba(15, 23, 42, 0.22);
  -webkit-overflow-scrolling: touch;
}
.is-map-floating .map-dash-panel-head {
  position: sticky;
  top: 0;
  z-index: 2;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin: 0 -16px;
  padding: 16px 16px 12px;
  background: #ffffff;
  border-bottom: 1px solid #f1f5f9;
}
.is-map-floating .map-dash-panel-head h2 { margin: 0; color: #0f172a; font-size: 1.125rem; font-weight: 700; }
.is-map-floating .map-dash-panel-head p { margin: 2px 0 0; color: #64748b; font-size: .8125rem; font-weight: 500; }
.map-dash-panel-close {
  flex: 0 0 auto;
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #ffffff;
  color: #475569;
  cursor: pointer;
}
.map-dash-panel-close:hover { background: #f1f5f9; }
.is-map-floating .dash-kpi-ribbon { display: grid !important; grid-template-columns: repeat(2, minmax(0, 1fr)) !important; gap: 10px !important; }
.is-map-floating .dash-two-col-grid { grid-template-columns: minmax(0, 1fr); gap: 14px; }
.is-map-floating .dash-table-wrap { overflow-x: auto; -webkit-overflow-scrolling: touch; }
@media (max-width: 720px) {
  .map-dashboard-fab { left: 10px; bottom: 12px; }
  .officer-summary-dashboard.is-map-floating { top: var(--mp-top); left: 10px; bottom: 10px; width: calc(100% - 20px); }
}

/* โหมดจัดการค่าใช้จ่าย: ปุ่มลอย 2 ปุ่ม แผงรายการ และการ์ดกรอง */
.map-dashboard-fab.is-cost-fab.fab-utility { bottom: 18px; }
.map-dashboard-fab.is-cost-fab.fab-bill { bottom: 72px; }
.map-cost-float {
  position: absolute;
  top: var(--mp-top);
  left: 18px;
  bottom: 18px;
  z-index: 60;
  box-sizing: border-box;
  width: min(640px, calc(100% - 36px));
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.98);
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  box-shadow: 0 18px 44px rgba(15, 23, 42, 0.22);
}
.market-map-kpi.is-overdue .market-map-kpi-icon { background: #fee2e2; color: #dc2626; }
.market-map-kpi-copy small.kpi-sub { color: #64748b; font-size: 0.8125rem; font-weight: 500; }
@media (max-width: 1000px) { .map-cost-float { top: var(--mp-top-2row); } }
@media (max-width: 720px) {
  .map-dashboard-fab.is-cost-fab.fab-utility { bottom: 12px; }
  .map-dashboard-fab.is-cost-fab.fab-bill { bottom: 66px; }
  .map-cost-float { top: var(--mp-top); left: 10px; bottom: 10px; width: calc(100% - 20px); }
}

/* ตัวกรองค่าใช้จ่าย (อยู่ในแผงสไลด์ของแผนที่) */
.cf-wrap { display: flex; flex-direction: column; gap: 16px; }
.cf-hint { margin: 0; color: #64748b; font-size: 0.875rem; font-weight: 500; line-height: 1.5; }
.cf-group { display: flex; flex-direction: column; gap: 10px; }
.cf-label { color: #475569; font-size: 0.875rem; font-weight: 600; }
.cf-chips { display: flex; flex-wrap: wrap; gap: 8px; }
.cf-chip { display: inline-flex; align-items: center; gap: 6px; min-height: 38px; padding: 6px 14px; border: 1px solid #e2e8f0; border-radius: 9999px; background: #fff; color: #475569; font-family: inherit; font-size: 0.875rem; font-weight: 600; cursor: pointer; }
.cf-chip span { color: #94a3b8; font-size: 0.8125rem; }
.cf-chip:hover { background: #f8fafc; }
.cf-chip.active { border-color: #2563eb; background: #2563eb; color: #fff; box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3); }
.cf-chip.active span { color: #dbeafe; }
.cf-result { margin: 0; padding: 10px 12px; border-radius: 10px; background: #eff6ff; color: #475569; font-size: 0.875rem; font-weight: 500; }
.cf-result strong { color: #2563eb; font-weight: 700; }

/* กันตัวอักษรในการ์ด KPI ล้นกรอบ */
.market-map-kpi { overflow: hidden; }
.market-map-kpi-icon { width: 34px; height: 34px; flex: 0 0 34px; }
.market-map-kpi-copy { flex: 1 1 auto; min-width: 0; overflow: visible; }
.market-map-kpi-copy small { max-width: 100%; }
.market-map-kpi-copy strong { min-width: 0; max-width: 100%; overflow: visible; gap: 5px; font-size: clamp(1rem, 1.15vw, 1.3rem); }
.market-map-kpi-copy strong em { flex: 0 0 auto; padding: 2px 6px; white-space: nowrap; }

/* ลิ้นชักแผงจอดชิดซ้ายแทนแผงโครงการ (หน้าข้อมูลตลาดบนจอ Desktop) */
@media (min-width: 769px) {
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock) {
    left: 18px !important;
    right: auto !important;
    width: min(380px, calc(100% - 36px)) !important;
  }
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-edit-mode) { width: min(440px, calc(100% - 36px)) !important; }
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode) {
    width: min(720px, calc(100% - 100px)) !important;
  }
  .is-fullscreen-map-panel :deep(.issd-mobile-only) {
    display: none !important;
  }
}

/* แผงโครงการ ↔ แผงรายละเอียดแผง: หน้าตาเดียวกัน ตำแหน่งเดียวกัน สลับทันทีโดยไม่มีอนิเมชัน */
@media (min-width: 769px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel) {
    width: min(380px, calc(100% - 36px)) !important;
    box-sizing: border-box !important;
    padding: 18px 20px !important;
    border: 1px solid #e2e8f0 !important;
    border-radius: 18px !important;
    background: #ffffff !important;
    box-shadow: 0 16px 38px rgba(15, 23, 42, 0.16) !important;
    transition: width .42s cubic-bezier(.2, .8, .2, 1), box-shadow .3s ease, background-color .18s ease !important;
  }
}
.is-fullscreen-map-panel :deep(.stall-mode-panel.is-expanding-to-stall) {
  width: min(720px, calc(100% - 100px)) !important;
  background: #f8fafc !important;
  overflow: hidden !important;
  transition: width .42s cubic-bezier(.2, .8, .2, 1), box-shadow .3s ease, background-color .18s ease !important;
  will-change: width;
  pointer-events: none;
}
.is-fullscreen-map-panel :deep(.stall-mode-panel.is-exiting-project) {
  opacity: 0 !important;
  transition: opacity .18s ease !important;
  pointer-events: none;
}
.is-fullscreen-map-panel :deep(.stall-mode-panel.is-expanding-to-stall > *) {
  animation: projectPanelContentOut .18s ease .22s both;
}
@keyframes projectPanelContentOut {
  from { opacity: 1; transform: translateX(0); }
  to { opacity: 0; transform: translateX(-8px); }
}
.is-fullscreen-map-panel :deep(.stall-mode-panel.is-returning-from-stall > *) {
  animation: projectPanelContentIn .18s ease both;
}
.is-fullscreen-map-panel :deep(.stall-mode-panel.is-returning-from-stall) {
  transition: width .62s cubic-bezier(.22, 1, .36, 1), box-shadow .4s ease, background-color .24s ease !important;
}
@keyframes projectPanelContentIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
/* แผงโครงการเข้า/ออก: ไม่ใช้ transition ของ Vue (ที่เลื่อน 16px แยกต่างหาก) ใช้แอนิเมชันเดียวกับแผงแผง ตอนเข้าเลื่อนเข้าจากซ้ายแบบเดียวกัน (ไม่เฟด) ตอนออกหายทันทีเหมือนกัน */
.is-fullscreen-map-panel :deep(.stall-mode-panel.is-swapped-out) { visibility: hidden; pointer-events: none; }
@media (min-width: 769px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel.smp-slide-enter-from),
  .is-fullscreen-map-panel :deep(.stall-mode-panel.smp-slide-leave-to) { opacity: 1; transform: none; }
  .is-fullscreen-map-panel :deep(.stall-mode-panel.smp-slide-enter-active),
  .is-fullscreen-map-panel :deep(.stall-mode-panel.smp-slide-leave-active) { transition: none; }
}

/* แผงแผงชิดซ้ายบนหน้าข้อมูลตลาด: สไลด์เข้า/ออกจากขอบซ้ายอย่างนุ่มนวล (GPU-accelerated slide) */
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.slide-fade-enter-active) {
  transition: transform 0.28s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.24s ease-out !important;
  will-change: transform, opacity;
}
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.slide-fade-leave-active) {
  transition: transform 0.24s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.2s ease-in !important;
  will-change: transform, opacity;
  pointer-events: none;
}
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.slide-fade-enter-from),
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.slide-fade-leave-to) {
  transform: translateX(calc(-100% - 24px)) !important;
  opacity: 0 !important;
}
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.slide-fade-enter-to),
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.slide-fade-leave-from) {
  transform: translateX(0) !important;
  opacity: 1 !important;
}
/* ตอนเปิดรายละเอียดแผง ให้เริ่มจากตำแหน่งและขนาด panel เดิม ไม่เลื่อนเข้าจากนอกจอ */
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.slide-fade-enter-active),
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.slide-fade-leave-active) {
  transition: transform .42s cubic-bezier(.2, .8, .2, 1), opacity .16s ease !important;
}
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.slide-fade-enter-from) {
  transform: scaleX(.528) !important;
  opacity: 1 !important;
}
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.slide-fade-enter-to) {
  transform: scaleX(1) !important;
  opacity: 1 !important;
}
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.slide-fade-leave-from),
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.slide-fade-leave-to) {
  transform: scaleX(.528) !important;
  opacity: 0 !important;
}
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap.slide-fade-enter-active) {
  transition: none !important;
}
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap.slide-fade-enter-from),
.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap.slide-fade-enter-to) {
  transform: scaleX(1) !important;
  opacity: 1 !important;
}
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap .issd-header,
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap .issd-full-sidebar,
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap .issd-full-main {
  animation: stallPanelContentFadeIn .16s ease-out both;
}
@keyframes stallPanelContentFadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-drawer-closing > * {
  opacity: 0;
  transition: opacity .12s ease !important;
}
@media (max-width: 768px) {
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.slide-fade-enter-from),
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.slide-fade-leave-to) {
    transform: translateY(100%) !important;
    opacity: 0 !important;
  }
}

/* ปุ่มย้อนกลับ: ตำแหน่งและขนาดตรงกันทุกพิกเซลระหว่างแผงโครงการกับแผงแผง */
.inspect-selected-stall-drawer.is-left-dock .issd-header { padding: 18px 20px 14px; }
.is-fullscreen-map-panel :deep(.smp-back-btn),
.issd-return-btn { box-sizing: border-box; height: 34px; min-height: 34px; line-height: 1; }

/* เส้นคั่นใต้ปุ่มย้อนกลับของแผงแผง (เหมือนแผงโครงการ) */
.issd-return-row { width: 100%; padding-bottom: 14px; margin-bottom: 14px; border-bottom: 1px solid #e8eef6; }
.issd-return-row .issd-return-btn { margin: 0; }

/* รายละเอียดแผงแบบเต็มบนจอ Desktop: หัวแผงกว้างเท่าแถบเมนูข้าง แล้วการ์ด "สิ่งที่ต้องดำเนินการ" อยู่บนสุดของคอลัมน์ขวา */
@media (min-width: 769px) {
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode {
    display: grid !important;
    grid-template-columns: 188px minmax(0, 1fr);
    grid-template-rows: auto minmax(0, 1fr) auto;
    background: #f8fafc;
  }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-full-layout { display: contents; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-full-main { grid-column: 2; grid-row: 1 / 4; display: flex; flex-direction: column; min-width: 0; min-height: 0; overflow-y: auto; overflow-x: hidden; -webkit-overflow-scrolling: touch; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-header {
    grid-column: 1;
    grid-row: 1;
    padding: 18px 14px 14px;
    border-right: 1px solid #eef2f6;
    background: #ffffff;
  }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-full-sidebar { grid-column: 1; grid-row: 2 / 4; width: auto; }
  /* การ์ดต้องดำเนินการอยู่บนสุดของคอลัมน์ขวา ชิดกับเนื้อหาด้านล่าง ไม่เว้นช่องว่าง */
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-action-alert-card { margin: 16px 18px 12px !important; flex: 0 0 auto; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-body-full { flex: 1 0 auto; min-height: 0; overflow: visible; }
  /* เลื่อนทั้งคอลัมน์ขวาพร้อมกัน (การ์ด + เนื้อหา) แถบปุ่มท้ายหน้าตรึงอยู่ด้านล่าง */
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-footer-full { position: sticky; bottom: 0; z-index: 2; margin-top: auto; flex: 0 0 auto; }
  /* หัวแผงในคอลัมน์แคบ: ชื่อแผงบรรทัดบน แล้วป้ายสถานะ + ชื่อโครงการอยู่บรรทัดเดียวกัน */
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-header {
    display: grid;
    grid-template-columns: auto minmax(0, 1fr);
    column-gap: 8px;
    row-gap: 8px;
    align-content: start;
    align-items: center;
  }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-return-row { grid-column: 1 / -1; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-header-top,
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-title-group { display: contents; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-title-group strong { grid-column: 1 / -1; font-size: 1.1875rem; line-height: 1.3; word-break: break-word; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-status-pill { grid-column: 1; justify-self: start; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-project-sub { grid-column: 2; margin: 0; font-size: 0.875rem; white-space: normal; overflow: visible; text-overflow: clip; line-height: 1.3; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-return-row { padding-bottom: 12px; margin-bottom: 12px; }
}

@media (prefers-reduced-motion: reduce) {
  .inspect-selected-stall-drawer { transition: none !important; }
}

/* มีแผงเปิดอยู่ทางซ้าย: ปุ่มลอยเลื่อนมาอยู่ข้างแผง (ชิดล่างเสมอ) ส่วนแผงชิดล่างของแผนที่ */
.has-left-panel .map-dashboard-fab { left: calc(18px + min(var(--lp-w), calc(100% - 36px)) + 12px); transition: left 0.2s ease; }
@media (max-width: 720px) { .has-left-panel .map-dashboard-fab { display: none; } }

/* ตัวเลือกชั้นอาคาร: แคปซูลมีตัวอักษร ที่มุมขวาล่าง ทางซ้ายของชุดปุ่มซูม ขอบล่างเสมอกับกลุ่มซูม */
.is-fullscreen-map-panel :deep(.map-floating-top-controls) {
  display: flex !important;
  top: auto !important;
  left: auto !important;
  right: 62px !important;
  bottom: var(--ov-ctrl-bottom, 39px) !important; /* วัดจากขอบล่างของชุดปุ่มซูมจริง */
  z-index: 40 !important;
  transition: transform 0.26s cubic-bezier(0.16, 1, 0.3, 1);
}
.is-fullscreen-map-panel :deep(.market-floor-switcher) {
  flex-direction: row;
  align-items: center;
  gap: 4px;
  padding: 4px;
  border: 1px solid #bfdbfe;
  border-radius: 999px;
  overflow: visible;
  background: rgba(255, 255, 255, 0.97);
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.14);
}
.is-fullscreen-map-panel :deep(.market-floor-tab) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: auto;
  min-width: 0;
  height: 38px;
  padding: 0 18px;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: #64748b;
  font-size: 0.9375rem;
  font-weight: 700;
  line-height: 1;
  white-space: nowrap;
  transition: background 0.15s ease, color 0.15s ease;
}
.is-fullscreen-map-panel :deep(.market-floor-tab:hover:not(.active)) { background: #f1f5f9; color: #0f172a; }
.is-fullscreen-map-panel :deep(.market-floor-tab.active) { background: #2563eb; color: #ffffff; box-shadow: 0 3px 10px rgba(37, 99, 235, 0.3); }


/* แผงโครงการและแผงรายละเอียดแผงบน Desktop: ความสูงเท่ากัน (ยืดจากใต้การ์ด KPI ถึงเหนือปุ่มแดชบอร์ด) */
@media (min-width: 769px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel.stall-mode-panel),
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-left-dock) {
    top: var(--mp-top) !important;
    bottom: 18px !important;
    height: auto !important;
    max-height: none !important;
  }
}
@media (min-width: 769px) and (max-width: 1000px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel.stall-mode-panel),
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-left-dock) { top: var(--mp-top-2row) !important; }
}

/* ปุ่มย้อนกลับของแผงโครงการและแผงแผงใช้สไตล์เดียวกัน (เท่า AppButton sm) */
.is-fullscreen-map-panel :deep(.smp-back-btn),
.issd-return-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  align-self: flex-start;
  max-width: 100%;
  height: 34px;
  margin: 0 0 12px;
  padding: 0 14px 0 10px;
  border: 1px solid #2563eb;
  border-radius: 7px;
  background: #2563eb;
  color: #ffffff;
  font-family: inherit;
  font-size: 0.875rem;
  font-weight: 600;
  line-height: 1;
  cursor: pointer;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.22);
  transition: background 0.15s ease, box-shadow 0.15s ease;
}
.issd-return-btn span { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.is-fullscreen-map-panel :deep(.smp-back-btn:hover),
.issd-return-btn:hover { background: #1d4ed8; box-shadow: 0 6px 14px rgba(37, 99, 235, 0.3); transform: none; }
.issd-return-btn:focus-visible { outline: 3px solid #93c5fd; outline-offset: 2px; }
.is-fullscreen-map-panel :deep(.smp-back-section) { padding-bottom: 14px; }
.is-fullscreen-map-panel :deep(.smp-back-btn) { margin: 0; }

/* ตัวอักษรของแผงโครงการให้อ่านง่ายตามมาตรฐานของระบบ */
.is-fullscreen-map-panel :deep(.smp-title) { font-size: 1.25rem; }
.is-fullscreen-map-panel :deep(.smp-dot) { width: 14px; height: 14px; }
.is-fullscreen-map-panel :deep(.smp-category) { margin-top: -4px; font-size: 0.9375rem; line-height: 1.5; }
.is-fullscreen-map-panel :deep(.smp-area-label) { font-size: 0.9375rem; }
.is-fullscreen-map-panel :deep(.smp-area-value) { font-size: 2rem; }
.is-fullscreen-map-panel :deep(.smp-area-value span) { font-size: 1rem; }
.is-fullscreen-map-panel :deep(.smp-stat-label) { font-size: 0.9375rem; }
.is-fullscreen-map-panel :deep(.smp-stat strong) { font-size: 1.5rem; }
.is-fullscreen-map-panel :deep(.smp-legend) { font-size: 0.9375rem; }
.is-fullscreen-map-panel :deep(.smp-hint) { font-size: 0.875rem; line-height: 1.55; }
.inspect-selected-stall-drawer.is-left-dock .issd-project-sub { font-size: 0.9375rem; }
.inspect-selected-stall-drawer.is-left-dock .issd-title-group strong { font-size: 1.1875rem; }

/* ================================================================ */
/* FULL-BLEED GIS MAP OVERVIEW VIEW (/stalls/map)                     */
/* ================================================================ */
.content-body.manage-page.is-map-page {
  padding: 0 !important;
  max-width: 100% !important;
  margin: 0 !important;
  height: calc(100vh - 64px);
  min-height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
  background: #f8fafc;
}

.is-fullscreen-map-view {
  height: 100% !important;
  flex: 1 !important;
  display: flex !important;
  flex-direction: column !important;
  margin: 0 !important;
  padding: 0 !important;
}

.is-fullscreen-map-panel {
  height: 100% !important;
  flex: 1 !important;
  padding: 0 !important;
  margin: 0 !important;
  border: none !important;
  border-radius: 0 !important;
  box-shadow: none !important;
  background: transparent !important;
  display: flex !important;
  flex-direction: column !important;
  position: relative !important;
}

.is-fullscreen-map-panel :deep(.projects-master-container) {
  height: 100% !important;
  flex: 1 !important;
  gap: 0 !important;
  display: flex !important;
  flex-direction: column !important;
  position: relative !important;
}

.is-fullscreen-map-panel :deep(.maplibre-overview-wrapper) {
  height: 100% !important;
  flex: 1 !important;
  border-radius: 0 !important;
  border: none !important;
  box-shadow: none !important;
}

/* ซ่อนปุ่ม แผนที่/ดาวเทียม, ปรับกึ่งกลาง และ คำอธิบายสีโครงการ (เพราะมีในแถบเครื่องมือชั้นข้อมูลด้านขวาแล้ว) */
.is-fullscreen-map-panel :deep(.map-floating-top-controls),
.is-fullscreen-map-panel :deep(.map-floating-style-switcher),
.is-fullscreen-map-panel :deep(.map-floating-recenter-btn),
.is-fullscreen-map-panel :deep(.map-floating-legend-container) {
  display: none !important;
}

@media (min-width: 769px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel) {
    top: var(--mp-top) !important;
    left: 18px !important;
    bottom: 18px !important;
    height: auto !important;
    max-height: none !important;
    z-index: 30 !important;
  }
}

.is-fullscreen-map-panel :deep(.map-project-guidance) {
  top: 16px !important;
  left: 18px !important;
  min-height: 48px !important;
  box-sizing: border-box !important;
  border-radius: 12px !important;
  padding: 8px 14px !important;
  background: rgba(255, 255, 255, 0.96) !important;
  backdrop-filter: blur(12px) !important;
  -webkit-backdrop-filter: blur(12px) !important;
  border: 1px solid rgba(226, 232, 240, 0.95) !important;
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08) !important;
  z-index: 100 !important;
}

.is-fullscreen-map-panel :deep(.map-floating-floor-preview-pill) {
  top: 16px !important;
  left: 18px !important;
}

@media (min-width: 769px) {
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer) {
    top: var(--mp-top) !important;
    right: 18px !important;
    max-height: calc(100% - (var(--mp-top) + 16px)) !important;
    z-index: 125 !important;
  }

  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-full-detail-mode) {
    width: min(840px, calc(100vw - 36px)) !important;
    z-index: 130 !important;
  }
}

/* ========================================================= */
/* OVERVIEW MAP DASHBOARD (แผงสไลด์แดชบอร์ดด้านขวาของแผนที่)  */
/* ========================================================= */
.ov-dash-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
}
.ov-dash-mode-view {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* 1. Hero Revenue / Health Box */
.ov-dash-hero-box {
  background: linear-gradient(135deg, #f8fafc 0%, #eff6ff 100%);
  border: 1px solid #dbeafe;
  border-radius: 14px;
  padding: 16px 18px;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.04);
}
.ov-dash-hero-box.finance {
  background: linear-gradient(135deg, #f0fdf4 0%, #eff6ff 100%);
  border-color: #bbf7d0;
}
.ov-dash-hero-box.meter {
  background: linear-gradient(135deg, #f0fdfa 0%, #eff6ff 100%);
  border-color: #99f6e4;
}
.ov-dash-hero-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}
.ov-dash-hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #1e40af;
  font-size: 0.8125rem;
  font-weight: 600;
  margin-bottom: 3px;
}
.ov-dash-hero-value {
  font-size: 1.75rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.02em;
  line-height: 1.15;
}
.ov-dash-hero-value small {
  font-size: 0.875rem;
  font-weight: 600;
  color: #64748b;
  letter-spacing: normal;
}
.ov-dash-hero-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 0.75rem;
  font-weight: 700;
  border-radius: 9999px;
  white-space: nowrap;
}
.ov-dash-hero-pill.green {
  background: #dcfce7;
  color: #15803d;
}
.ov-dash-health-bar-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.ov-dash-health-labels {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}
.health-label-good { color: #16a34a; }
.health-label-warn { color: #d97706; }
.health-label-danger { color: #dc2626; }
.ov-dash-health-track {
  display: flex;
  width: 100%;
  height: 7px;
  background: #e2e8f0;
  border-radius: 9999px;
  overflow: hidden;
}
.ov-dash-health-track span {
  height: 100%;
  transition: width 0.3s ease;
}
.ov-dash-health-track .track-good { background: #22c55e; }
.ov-dash-health-track .track-warn { background: #f59e0b; }
.ov-dash-health-track .track-danger { background: #ef4444; }

/* 2. 4 Contract Intelligence Cards (2x2 Grid) */
.ov-dash-kpi-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}
.ov-dash-kpi-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  transition: border-color 0.2s, box-shadow 0.2s;
}
.ov-dash-kpi-card:hover {
  border-color: #cbd5e1;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}
.ov-dash-kpi-card.card-blue { border-left: 3px solid #2563eb; }
.ov-dash-kpi-card.card-amber { border-left: 3px solid #d97706; }
.ov-dash-kpi-card.card-red { border-left: 3px solid #dc2626; }
.ov-dash-kpi-card.card-green { border-left: 3px solid #16a34a; }

.ov-dash-kpi-top {
  display: flex;
  align-items: center;
  gap: 8px;
}
.ov-dash-kpi-icon {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.ov-dash-kpi-icon.blue { background: #eff6ff; color: #2563eb; }
.ov-dash-kpi-icon.amber { background: #fef3c7; color: #d97706; }
.ov-dash-kpi-icon.red { background: #fee2e2; color: #dc2626; }
.ov-dash-kpi-icon.green { background: #dcfce7; color: #16a34a; }

.ov-dash-kpi-title {
  font-size: 0.8125rem;
  color: #64748b;
  font-weight: 600;
  white-space: nowrap;
}
.ov-dash-kpi-number {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.15;
}
.ov-dash-kpi-number small {
  font-size: 0.8125rem;
  font-weight: 600;
  color: #64748b;
}
.ov-dash-kpi-number.text-amber { color: #d97706; }
.ov-dash-kpi-number.text-red { color: #dc2626; }
.ov-dash-kpi-number.text-green { color: #16a34a; }

.ov-dash-legal-card {
  width: 100%; min-height: 64px; display: flex; align-items: center; gap: 10px;
  padding: 10px 12px; border: 1px solid #fecaca; border-radius: 12px;
  background: #fff7f7; color: #991b1b; font-family: inherit; text-align: left; cursor: pointer;
  transition: border-color .15s ease, background .15s ease, box-shadow .15s ease;
}
.ov-dash-legal-card:hover,.ov-dash-legal-card.active { border-color: #ef4444; background: #fef2f2; box-shadow: 0 2px 8px rgba(220,38,38,.1); }
.ov-dash-legal-icon { width: 34px; height: 34px; flex: 0 0 34px; display: grid; place-items: center; border-radius: 9px; background: #fee2e2; color: #dc2626; }
.ov-dash-legal-icon svg { width: 18px; height: 18px; }
.ov-dash-legal-copy { flex: 1 1 auto; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.ov-dash-legal-copy small { font-size: .8125rem; color: #7f1d1d; font-weight: 600; }
.ov-dash-legal-copy strong { font-size: 1.125rem; color: #b91c1c; font-weight: 800; line-height: 1.2; }
.ov-dash-legal-copy em { font-size: .8125rem; font-style: normal; font-weight: 600; }
.ov-dash-legal-action { flex: 0 0 auto; color: #dc2626; font-size: .8125rem; font-weight: 600; }

/* การ์ดสรุปร่วมสำหรับแดชบอร์ดข้อมูลแผง การชำระเงิน และการจดมิเตอร์ */
.ov-dash-summary-card,
.ov-stall-overview-card,
.ov-stall-priority-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 14px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}

.ov-stall-card-head,
.ov-stall-section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.ov-stall-card-head h4,
.ov-stall-section-head h4 {
  margin: 0;
  color: #0f172a;
  font-size: 0.9375rem;
  font-weight: 700;
  line-height: 1.35;
}

.ov-stall-card-head p {
  margin: 2px 0 0;
  color: #64748b;
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-stall-section-note {
  margin: 2px 0 0;
  color: #64748b;
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-stall-occupancy {
  display: inline-flex;
  align-items: center;
  min-height: 26px;
  padding: 3px 9px;
  border-radius: 999px;
  background: #eff6ff;
  color: #2563eb;
  font-size: 0.75rem;
  font-weight: 700;
  white-space: nowrap;
}

.ov-stall-status-track {
  display: flex;
  width: 100%;
  height: 8px;
  margin: 12px 0;
  overflow: hidden;
  border-radius: 999px;
  background: #e2e8f0;
}

.ov-stall-status-track span {
  min-width: 0;
  height: 100%;
  transition: width 0.25s ease;
}

.ov-stall-status-track .is-available,
.ov-stall-status-actions i.is-available { background: #10b981; }
.ov-stall-status-track .is-reserved,
.ov-stall-status-actions i.is-reserved { background: #f59e0b; }
.ov-stall-status-track .is-occupied,
.ov-stall-status-actions i.is-occupied { background: #64748b; }

.ov-stall-status-actions {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 7px;
}

.ov-stall-status-actions button {
  min-width: 0;
  min-height: 58px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: space-between;
  gap: 4px;
  border: 1px solid #eef2f6;
  border-radius: 10px;
  background: #f8fafc;
  color: #334155;
  font-family: inherit;
  text-align: left;
  cursor: pointer;
  transition: border-color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease;
}

.ov-stall-status-actions button:hover,
.ov-stall-status-actions button.active {
  border-color: #93c5fd;
  background: #eff6ff;
  box-shadow: 0 1px 4px rgba(37, 99, 235, 0.08);
}

.ov-stall-status-actions button span {
  display: flex;
  align-items: center;
  gap: 5px;
  max-width: 100%;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 600;
  white-space: nowrap;
}

.ov-stall-status-actions button i {
  width: 7px;
  height: 7px;
  flex: 0 0 7px;
  border-radius: 50%;
}

.ov-stall-status-actions button strong {
  color: #0f172a;
  font-size: 1.125rem;
  font-weight: 800;
  line-height: 1.1;
}

.ov-stall-priority-card,
.ov-dash-summary-card {
  display: flex;
  flex-direction: column;
  gap: 9px;
}

.ov-stall-priority-row {
  width: 100%;
  min-height: 44px;
  padding: 7px 9px;
  display: grid;
  grid-template-columns: 30px minmax(0, 1fr) auto 16px;
  align-items: center;
  gap: 8px;
  border: 1px solid #eef2f6;
  border-radius: 10px;
  background: #f8fafc;
  color: #334155;
  font-family: inherit;
  text-align: left;
}

button.ov-stall-priority-row {
  cursor: pointer;
  transition: border-color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease;
}

button.ov-stall-priority-row:hover,
button.ov-stall-priority-row.active {
  border-color: #93c5fd;
  background: #eff6ff;
  box-shadow: 0 1px 4px rgba(37, 99, 235, 0.08);
}

.ov-stall-priority-icon {
  width: 30px;
  height: 30px;
  display: grid;
  place-items: center;
  border-radius: 8px;
}

.ov-stall-priority-icon svg {
  width: 17px;
  height: 17px;
}

.ov-stall-priority-row.is-action .ov-stall-priority-icon { background: #eff6ff; color: #2563eb; }
.ov-stall-priority-row.is-legal .ov-stall-priority-icon { background: #fee2e2; color: #dc2626; }
.ov-stall-priority-row.is-expiring .ov-stall-priority-icon { background: #fef3c7; color: #d97706; }

.ov-stall-priority-label {
  min-width: 0;
  overflow: hidden;
  color: #334155;
  font-size: 0.8125rem;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ov-stall-priority-row strong {
  min-width: 24px;
  color: #0f172a;
  font-size: 1rem;
  font-weight: 800;
  text-align: right;
}

.ov-stall-priority-row.is-legal strong { color: #dc2626; }
.ov-stall-priority-row.is-expiring strong { color: #b45309; }

.ov-stall-priority-chevron {
  width: 15px;
  height: 15px;
  color: #94a3b8;
}

.ov-dash-summary-main {
  min-height: 50px;
  padding: 9px 11px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border-radius: 10px;
  background: #eff6ff;
}

.ov-dash-summary-main span,
.ov-dash-summary-grid span {
  color: #64748b;
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-dash-summary-main strong {
  color: #1d4ed8;
  font-size: 1.125rem;
  font-weight: 800;
  white-space: nowrap;
}

.ov-dash-summary-main strong small {
  margin-left: 4px;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-dash-utilization {
  min-height: 58px;
  padding: 10px 11px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
  border-radius: 10px;
  background: #eff6ff;
}

.ov-dash-utilization-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 10px;
}

.ov-dash-utilization-head span {
  color: #475569;
  font-size: 0.8125rem;
  font-weight: 600;
}

.ov-dash-utilization-head strong {
  color: #1d4ed8;
  font-size: 1.125rem;
  font-weight: 800;
}

.ov-dash-utilization-track {
  width: 100%;
  height: 6px;
  overflow: hidden;
  border-radius: 999px;
  background: #dbeafe;
}

.ov-dash-utilization-track span {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: #2563eb;
  transition: width 0.25s ease;
}

.ov-dash-summary-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.ov-dash-summary-grid > div {
  min-width: 0;
  padding: 9px 10px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  border: 1px solid #eef2f6;
  border-radius: 10px;
  background: #f8fafc;
}

.ov-dash-summary-grid strong {
  color: #0f172a;
  font-size: 1rem;
  font-weight: 800;
  white-space: nowrap;
}

.ov-dash-summary-grid .is-success strong { color: #15803d; }
.ov-dash-summary-grid .is-pending strong { color: #b45309; }
.ov-dash-summary-grid .is-danger strong,
.ov-dash-summary-grid .is-expired strong { color: #dc2626; }

.ov-dash-summary-grid small {
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-dash-meter-progress {
  min-height: 58px;
  padding: 10px 11px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
  border-radius: 10px;
  background: #eff6ff;
}

.ov-dash-meter-progress-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 10px;
}

.ov-dash-meter-progress-head span {
  color: #64748b;
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-dash-meter-progress-head strong {
  color: #1d4ed8;
  font-size: 1rem;
  font-weight: 800;
  white-space: nowrap;
}

.ov-dash-meter-progress-head small {
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-dash-income-total {
  min-height: 42px;
  margin-bottom: 10px;
  padding: 8px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  border-radius: 9px;
  background: #eff6ff;
}

.ov-dash-income-total span {
  color: #475569;
  font-size: 0.8125rem;
  font-weight: 600;
}

.ov-dash-income-total strong {
  color: #1d4ed8;
  font-size: 1rem;
  font-weight: 800;
  white-space: nowrap;
}

.ov-project-status-legend {
  margin: -2px 0 9px;
  display: flex;
  align-items: center;
  gap: 12px;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-project-status-legend span,
.ov-contract-health-legend span,
.ov-maintenance-legend span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.ov-project-status-legend i,
.ov-contract-health-legend i,
.ov-maintenance-legend i {
  width: 8px;
  height: 8px;
  flex: 0 0 8px;
  border-radius: 3px;
}

.ov-project-status-legend i.is-available,
.ov-project-stacked-bar .is-available { background: #10b981; }
.ov-project-status-legend i.is-occupied,
.ov-project-stacked-bar .is-occupied { background: #64748b; }
.ov-project-status-legend i.is-reserved,
.ov-project-stacked-bar .is-reserved { background: #f59e0b; }

.ov-project-status-list {
  display: flex;
  flex-direction: column;
  gap: 9px;
  max-height: 360px;
  padding-right: 5px;
  overflow-y: auto;
  overscroll-behavior: contain;
  scrollbar-gutter: stable;
}

.ov-project-status-list::-webkit-scrollbar { width: 6px; }
.ov-project-status-list::-webkit-scrollbar-track { background: transparent; }
.ov-project-status-list::-webkit-scrollbar-thumb { border-radius: 999px; background: #cbd5e1; }
.ov-project-status-list::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

.ov-project-donut-focus {
  margin-top: 14px;
  padding: 16px 14px;
  display: grid;
  grid-template-columns: 132px minmax(0, 1fr);
  align-items: center;
  gap: 16px;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  background: linear-gradient(145deg, #ffffff 0%, #f8fafc 100%);
}

.ov-project-donut-wrap {
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.ov-project-donut-chart {
  position: relative;
  width: 126px;
  height: 126px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  box-shadow: 0 7px 20px rgba(15, 23, 42, 0.1);
  transition: background 0.25s ease;
}

.ov-project-donut-chart::after {
  position: absolute;
  inset: 17px;
  border-radius: 50%;
  background: #ffffff;
  box-shadow: inset 0 0 0 1px #f1f5f9;
  content: '';
}

.ov-project-donut-chart > span {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1px;
}

.ov-project-donut-chart strong {
  color: #0f172a;
  font-size: 1.5rem;
  font-weight: 800;
  line-height: 1.15;
}

.ov-project-donut-chart small {
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-project-donut-name {
  max-width: 100%;
  overflow: hidden;
  color: #334155;
  font-size: 0.875rem;
  font-weight: 700;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ov-project-donut-stats {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 7px;
}

.ov-project-donut-stats > div {
  min-height: 42px;
  padding: 6px 8px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: center;
  gap: 1px 7px;
  border: 1px solid #eef2f6;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.88);
}

.ov-project-donut-stats span {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 6px;
  color: #64748b;
  font-size: 0.8125rem;
  font-weight: 600;
  white-space: nowrap;
}

.ov-project-donut-stats i {
  width: 8px;
  height: 8px;
  flex: 0 0 8px;
  border-radius: 3px;
}

.ov-project-donut-stats strong {
  color: #1e293b;
  font-size: 1rem;
  font-weight: 800;
}

.ov-project-donut-stats small {
  grid-column: 1 / -1;
  padding-left: 14px;
  color: #94a3b8;
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-project-donut-stats .is-available i { background: #10b981; }
.ov-project-donut-stats .is-occupied i { background: #64748b; }
.ov-project-donut-stats .is-reserved i { background: #f59e0b; }
.ov-project-donut-stats .is-available strong { color: #047857; }
.ov-project-donut-stats .is-reserved strong { color: #b45309; }

.ov-project-donut-empty {
  margin-top: 14px;
  padding: 28px 16px;
  border: 1px dashed #cbd5e1;
  border-radius: 12px;
  background: #f8fafc;
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
  text-align: center;
}

.ov-project-mini-donut-grid {
  max-height: 382px;
  padding: 2px 4px 4px 2px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  overflow-y: auto;
  overscroll-behavior: contain;
  scrollbar-gutter: stable;
}

.ov-project-mini-donut-grid::-webkit-scrollbar { width: 6px; }
.ov-project-mini-donut-grid::-webkit-scrollbar-track { background: transparent; }
.ov-project-mini-donut-grid::-webkit-scrollbar-thumb { border-radius: 999px; background: #cbd5e1; }

.ov-project-mini-donut-card {
  position: relative;
  min-width: 0;
  min-height: 158px;
  padding: 14px 8px 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  border-radius: 13px;
  background: linear-gradient(155deg, #ffffff 0%, #f8fafc 100%);
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
  outline: none;
  transition: border-color 0.18s ease, box-shadow 0.18s ease, transform 0.18s ease;
}

.ov-project-mini-donut-card:hover,
.ov-project-mini-donut-card:focus-visible,
.ov-project-mini-donut-card.active {
  z-index: 2;
  border-color: #93c5fd;
  box-shadow: 0 7px 18px rgba(37, 99, 235, 0.13);
  transform: translateY(-2px);
}

.ov-project-mini-donut-card.active {
  background: linear-gradient(155deg, #ffffff 0%, #eff6ff 100%);
  box-shadow: 0 0 0 1px #bfdbfe, 0 7px 18px rgba(37, 99, 235, 0.13);
}

.ov-project-mini-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
  padding: 3px 7px;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
  line-height: 1.2;
}

.ov-project-mini-badge.is-danger { background: #fee2e2; color: #dc2626; }
.ov-project-mini-badge.is-available { background: #dcfce7; color: #15803d; }

.ov-project-mini-donut {
  position: relative;
  width: 86px;
  height: 86px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.09);
  transition: transform 0.2s ease;
}

.ov-project-mini-donut-card:hover .ov-project-mini-donut { transform: scale(1.04); }

.ov-project-mini-donut::after {
  position: absolute;
  inset: 12px;
  border-radius: 50%;
  background: #ffffff;
  box-shadow: inset 0 0 0 1px #f1f5f9;
  content: '';
}

.ov-project-mini-donut > span {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  line-height: 1.05;
}

.ov-project-mini-donut strong { color: #0f172a; font-size: 1.125rem; font-weight: 800; }
.ov-project-mini-donut small { color: #64748b; font-size: 0.75rem; font-weight: 600; }

.ov-project-mini-name {
  width: 100%;
  overflow: hidden;
  color: #334155;
  font-size: 0.8125rem;
  font-weight: 700;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ov-project-mini-tooltip {
  position: absolute;
  inset: auto 7px 7px;
  z-index: 4;
  padding: 8px 9px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  border-radius: 9px;
  background: rgba(15, 23, 42, 0.94);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.22);
  opacity: 0;
  pointer-events: none;
  transform: translateY(5px);
  transition: opacity 0.16s ease, transform 0.16s ease;
}

.ov-project-mini-donut-card:hover .ov-project-mini-tooltip,
.ov-project-mini-donut-card:focus-visible .ov-project-mini-tooltip {
  opacity: 1;
  transform: translateY(0);
}

.ov-project-mini-tooltip span {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto;
  align-items: center;
  gap: 6px;
  color: #e2e8f0;
  font-size: 0.75rem;
  font-weight: 500;
}

.ov-project-mini-tooltip strong { color: #ffffff; font-size: 0.8125rem; font-weight: 700; }
.ov-project-mini-tooltip small { color: #cbd5e1; font-size: 0.75rem; font-weight: 500; }
.ov-project-mini-tooltip .is-available::before,
.ov-project-mini-tooltip .is-occupied::before,
.ov-project-mini-tooltip .is-reserved::before { display: none; }

@media (max-width: 420px) {
  .ov-project-donut-focus { grid-template-columns: 1fr; }
  .ov-project-donut-stats { width: 100%; }
  .ov-project-mini-donut-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

.ov-project-status-row {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.ov-project-status-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  color: #475569;
  font-size: 0.8125rem;
  font-weight: 600;
}

.ov-project-status-meta span {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ov-project-status-meta strong {
  color: #1e293b;
  font-size: 0.8125rem;
  font-weight: 700;
}

.ov-project-stacked-bar,
.ov-maintenance-stacked-bar {
  display: flex;
  width: 100%;
  height: 8px;
  overflow: hidden;
  border-radius: 999px;
  background: #e2e8f0;
}

.ov-project-stacked-bar span {
  min-width: 0;
  height: 100%;
  cursor: help;
  transition: width 0.25s ease;
}

.ov-contract-health-layout {
  display: grid;
  grid-template-columns: 92px minmax(0, 1fr);
  align-items: center;
  gap: 16px;
}

.ov-contract-donut {
  position: relative;
  width: 88px;
  height: 88px;
  display: grid;
  place-items: center;
  border-radius: 50%;
}

.ov-contract-donut::after {
  position: absolute;
  inset: 12px;
  border-radius: 50%;
  background: #ffffff;
  content: '';
}

.ov-contract-donut > span {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  line-height: 1.1;
}

.ov-contract-donut strong {
  color: #0f172a;
  font-size: 1.125rem;
  font-weight: 800;
}

.ov-contract-donut small {
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-contract-health-legend {
  display: flex;
  flex-direction: column;
  gap: 7px;
  color: #475569;
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-contract-health-legend span strong {
  margin-left: auto;
  color: #1e293b;
  font-size: 0.8125rem;
  font-weight: 700;
}

.ov-contract-health-legend i.is-active { background: #16a34a; }
.ov-contract-health-legend i.is-expiring { background: #f59e0b; }
.ov-contract-health-legend i.is-expired { background: #ef4444; }

.ov-contract-renewal-row {
  margin-top: 12px;
  padding-top: 10px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 6px;
  border-top: 1px solid #f1f5f9;
}

.ov-contract-renewal-row span {
  min-width: 0;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 500;
  text-align: center;
}

.ov-contract-renewal-row strong {
  display: block;
  margin-top: 2px;
  color: #1e293b;
  font-size: 0.9375rem;
  font-weight: 800;
}

.ov-status-summary-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 7px;
}

.ov-status-summary-grid.is-five {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.ov-status-summary-grid > div {
  min-width: 0;
  min-height: 48px;
  padding: 7px 9px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2px;
  border: 1px solid #eef2f6;
  border-radius: 9px;
  background: #f8fafc;
}

.ov-status-summary-grid span {
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 500;
  line-height: 1.25;
}

.ov-status-summary-grid strong {
  color: #1e293b;
  font-size: 1rem;
  font-weight: 800;
}

.ov-status-summary-grid .is-success strong { color: #15803d; }
.ov-status-summary-grid .is-pending strong { color: #b45309; }
.ov-status-summary-grid .is-danger strong { color: #dc2626; }

.ov-maintenance-stacked-bar {
  height: 10px;
  margin-bottom: 11px;
}

.ov-maintenance-stacked-bar span {
  min-width: 0;
  height: 100%;
  transition: flex 0.25s ease;
}

.ov-maintenance-stacked-bar .is-pending,
.ov-maintenance-legend i.is-pending { background: #f59e0b; }
.ov-maintenance-stacked-bar .is-repairing,
.ov-maintenance-legend i.is-repairing { background: #2563eb; }
.ov-maintenance-stacked-bar .is-completed,
.ov-maintenance-legend i.is-completed { background: #16a34a; }
.ov-maintenance-stacked-bar .is-overdue,
.ov-maintenance-legend i.is-overdue { background: #ef4444; }

.ov-maintenance-legend {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 7px 12px;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 500;
}

.ov-maintenance-legend strong {
  margin-left: auto;
  color: #1e293b;
  font-size: 0.8125rem;
  font-weight: 700;
}

/* เพิ่มความชัดเจนของตัวอักษรในแดชบอร์ด โดยคงสัดส่วนการ์ดเดิม */
.ov-dash-container .ov-dash-section-title,
.ov-dash-container .ov-stall-card-head h4,
.ov-dash-container .ov-stall-section-head h4 {
  font-size: 1rem;
}

.ov-dash-container .ov-dash-section-note,
.ov-dash-container .ov-stall-card-head p,
.ov-dash-container .ov-stall-section-note,
.ov-dash-container .ov-project-status-legend,
.ov-dash-container .ov-contract-renewal-row span,
.ov-dash-container .ov-status-summary-grid span,
.ov-dash-container .ov-maintenance-legend,
.ov-dash-container .ov-contract-donut small,
.ov-dash-container .ov-dash-summary-grid small,
.ov-dash-container .ov-dash-meter-progress-head small {
  font-size: 0.8125rem;
  line-height: 1.4;
}

.ov-dash-container .ov-project-status-meta,
.ov-dash-container .ov-project-status-meta strong,
.ov-dash-container .ov-contract-health-legend,
.ov-dash-container .ov-contract-health-legend span strong,
.ov-dash-container .ov-maintenance-legend strong,
.ov-dash-container .ov-dash-summary-main span,
.ov-dash-container .ov-dash-summary-grid span,
.ov-dash-container .ov-dash-meter-progress-head span,
.ov-dash-container .ov-dash-income-total span {
  font-size: 0.875rem;
}

.ov-dash-container .ov-status-summary-grid strong,
.ov-dash-container .ov-dash-summary-grid strong,
.ov-dash-container .ov-dash-meter-progress-head strong,
.ov-dash-container .ov-dash-income-total strong {
  font-size: 1.0625rem;
}

/* 3. Section Cards */
.ov-dash-section-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 14px 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
}
.ov-dash-section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}
.ov-dash-section-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.ov-dash-section-title {
  margin: 0;
  font-size: 0.95rem;
  font-weight: 700;
  color: #0f172a;
}
.ov-dash-count-badge {
  font-size: 0.75rem;
  font-weight: 700;
  padding: 1px 7px;
  border-radius: 9999px;
}
.ov-dash-count-badge.badge-amber {
  background: #fef3c7;
  color: #b45309;
}
.ov-dash-count-badge.badge-red {
  background: #fee2e2;
  color: #dc2626;
}

/* 4. Streamlined Watchlist Rows */
.ov-dash-watchlist {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.ov-dash-watchlist-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 9px 12px;
  background: #f8fafc;
  border: 1px solid #eef2f6;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s ease;
}
.ov-dash-watchlist-row:hover {
  background: #ffffff;
  border-color: #93c5fd;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.06);
  transform: translateY(-1px);
}
.ov-dash-watchlist-row.urgency-expired {
  border-left: 3px solid #ef4444;
  background: #fef2f2;
}
.ov-dash-watchlist-row.urgency-expired:hover {
  background: #ffffff;
  border-color: #fca5a5;
}
.ov-dash-watchlist-row.urgency-urgent {
  border-left: 3px solid #f97316;
  background: #fff7ed;
}
.ov-dash-watchlist-row.urgency-urgent:hover {
  background: #ffffff;
  border-color: #fdba74;
}
.ov-dash-watchlist-row.urgency-warning {
  border-left: 3px solid #f59e0b;
  background: #fffbeb;
}
.ov-dash-watchlist-row.urgency-warning:hover {
  background: #ffffff;
  border-color: #fde68a;
}

.ov-dash-w-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex: 1 1 auto;
}
.ov-dash-code-tag {
  font-size: 0.8125rem;
  font-weight: 800;
  color: #1e293b;
  background: #e2e8f0;
  padding: 2px 6px;
  border-radius: 5px;
  letter-spacing: 0.02em;
  flex-shrink: 0;
}
.ov-dash-w-tenant {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ov-dash-w-right {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}
.ov-dash-w-pill {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 9999px;
  white-space: nowrap;
}
.ov-dash-w-pill.pill-expired {
  background: #fee2e2;
  color: #dc2626;
  font-weight: 700;
}
.ov-dash-w-pill.pill-urgent {
  background: #ffedd5;
  color: #c2410c;
  font-weight: 700;
}
.ov-dash-w-pill.pill-warning {
  background: #fef3c7;
  color: #b45309;
}
.ov-dash-w-pill.pill-renew-submitted {
  background: #dcfce7;
  color: #15803d;
}
.ov-dash-w-pill.pill-renew-pending_submit {
  background: #f1f5f9;
  color: #64748b;
}
.ov-dash-w-pill.pill-renew-declined {
  background: #fee2e2;
  color: #991b1b;
}

.ov-dash-w-icon {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  background: #f1f5f9;
  transition: all 0.15s ease;
}
.ov-dash-watchlist-row:hover .ov-dash-w-icon {
  color: #2563eb;
  background: #eff6ff;
}

/* 5. Minimalist 4-Stat Strip */
.ov-dash-pipeline-strip {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8px;
  background: #f8fafc;
  border: 1px solid #eef2f6;
  border-radius: 10px;
  padding: 10px 8px;
}
.ov-dash-pipe-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 2px;
}
.ov-dash-pipe-stat small {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 500;
  white-space: nowrap;
}
.ov-dash-pipe-stat strong {
  font-size: 1.25rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.1;
}
.ov-dash-pipe-stat strong.text-blue { color: #2563eb; }
.ov-dash-pipe-stat strong.text-amber { color: #d97706; }
.ov-dash-pipe-stat strong.text-purple { color: #7c3aed; }
.ov-dash-pipe-stat strong.text-green { color: #16a34a; }

/* 6. Category Breakdown List */
.ov-dash-category-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.ov-dash-category-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-bottom: 6px;
  border-bottom: 1px solid #f1f5f9;
}
.ov-dash-category-row:last-child {
  border-bottom: none;
  padding-bottom: 0;
}
.ov-dash-cat-left {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.ov-dash-cat-name {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8125rem;
  font-weight: 600;
  color: #334155;
}
.ov-dash-cat-name .cat-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}
.ov-dash-cat-pct {
  font-size: 0.75rem;
  font-weight: 700;
  color: #64748b;
}
.ov-dash-cat-bar {
  width: 100%;
  height: 5px;
  background: #f1f5f9;
  border-radius: 9999px;
  overflow: hidden;
}
.ov-dash-cat-bar span {
  display: block;
  height: 100%;
  border-radius: 9999px;
  background: #2563eb;
  transition: width 0.3s ease;
}

/* ========================================================= */
/* Responsive Stall Status KPI Cards inside Dashboard        */
/* ========================================================= */
.ov-dash-responsive-kpis {
  display: none;
}

@media (max-width: 1200px) {
  .ov-dash-responsive-kpis {
    display: block !important;
  }
}

.ov-dash-status-pill-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  margin-top: 6px;
  width: 100%;
  box-sizing: border-box;
}

.ov-dash-status-pill {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #ffffff;
  cursor: pointer;
  text-align: left;
  font-family: inherit;
  transition: all 0.16s ease;
  position: relative;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.03);
  min-width: 0;
  box-sizing: border-box;
  width: 100%;
}

.ov-dash-status-pill:hover {
  transform: translateY(-1px);
  border-color: #cbd5e1;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.06);
}

.ov-dash-status-pill:focus-visible {
  outline: 2px solid #2563eb;
  outline-offset: 1px;
}

.ov-dash-status-pill-icon {
  flex: 0 0 32px;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: grid;
  place-items: center;
  transition: all 0.16s ease;
}

.ov-dash-status-pill-icon svg {
  width: 17px;
  height: 17px;
}

.ov-dash-status-pill-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1 1 auto;
}

.ov-dash-status-pill-label {
  font-size: 0.8125rem;
  color: #64748b;
  font-weight: 500;
  line-height: 1.25;
  white-space: normal;
  word-break: break-word;
}

.ov-dash-status-pill-val {
  font-size: 1.0625rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.2;
  margin-top: 2px;
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.ov-dash-status-pill-val small {
  font-size: 0.75rem;
  font-weight: 500;
  color: #64748b;
}

/* Color schemes for each status */
.ov-dash-status-pill.is-total .ov-dash-status-pill-icon { background: #eff6ff; color: #2563eb; }
.ov-dash-status-pill.is-total.active {
  border-color: #2563eb;
  background: #f8faff;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.18);
}

.ov-dash-status-pill.is-available .ov-dash-status-pill-icon { background: #ecfdf5; color: #059669; }
.ov-dash-status-pill.is-available.active {
  border-color: #10b981;
  background: #f0fdf4;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.ov-dash-status-pill.is-reserved .ov-dash-status-pill-icon { background: #fffbeb; color: #d97706; }
.ov-dash-status-pill.is-reserved.active {
  border-color: #f59e0b;
  background: #fffdf7;
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.2);
}

.ov-dash-status-pill.is-occupied .ov-dash-status-pill-icon { background: #f1f5f9; color: #64748b; }
.ov-dash-status-pill.is-occupied.active {
  border-color: #94a3b8;
  background: #f8fafc;
  box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.18);
}

.ov-dash-status-pill.is-legal-case .ov-dash-status-pill-icon { background: #fef2f2; color: #ef4444; }
.ov-dash-status-pill.is-legal-case.active {
  border-color: #ef4444;
  background: #fef2f2;
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
}

.ov-dash-status-pill.is-action-required .ov-dash-status-pill-icon { background: #eef2ff; color: #6366f1; }
.ov-dash-status-pill.is-action-required.active {
  border-color: #6366f1;
  background: #f5f3ff;
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
}

/* Market information overview: KPI cards floating above the GIS map */
.market-map-kpi-strip {
  position: absolute;
  top: 18px;
  left: 50%;
  right: auto;
  transform: translateX(-50%);
  width: min(1240px, calc(100% - 320px));
  max-width: 1240px;
  margin: 0;
  z-index: 110;
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 10px;
  pointer-events: none;
  transition: left 0.3s cubic-bezier(0.16, 1, 0.3, 1),
              width 0.3s cubic-bezier(0.16, 1, 0.3, 1),
              transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

/* เมื่อเปิด panel ให้ขยับเฉพาะพื้นที่วาง โดยคงดีไซน์การ์ด ไอคอน และ typography เดิม */
.market-map-kpi-strip.has-side-panel {
  left: calc((100% - var(--side-w, 380px) + 180px) / 2);
  width: min(1050px, calc(100% - var(--side-w, 380px) - 170px));
}

.market-map-kpi {
  min-width: 0;
  min-height: 72px;
  padding: 12px 14px;
  display: flex;
  align-items: center;
  gap: 10px;
  border: 1px solid rgba(226, 232, 240, 0.94);
  border-radius: 15px;
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.11);
  color: #475569;
  font-family: inherit;
  text-align: left;
  pointer-events: auto;
}

button.market-map-kpi {
  cursor: pointer;
  transition: border-color 0.16s ease, box-shadow 0.16s ease, transform 0.16s ease;
}

button.market-map-kpi:hover {
  transform: translateY(-1px);
  border-color: #bfdbfe;
  box-shadow: 0 10px 26px rgba(15, 23, 42, 0.14);
}

.market-map-kpi.active {
  border: 2px solid #3b82f6;
  padding: 11px 13px;
  box-shadow: 0 8px 24px rgba(37, 99, 235, 0.2);
}

.market-map-kpi-icon {
  width: 40px;
  height: 40px;
  flex: 0 0 40px;
  display: grid;
  place-items: center;
  border-radius: 11px;
  background: #eff6ff;
  color: #2563eb;
}

.market-map-kpi-icon svg {
  width: 20px;
  height: 20px;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.market-map-kpi-copy {
  min-width: 0;
  display: grid;
  gap: 2px;
}

.market-map-kpi-copy small {
  overflow: hidden;
  color: #475569;
  font-size: 0.8125rem;
  font-weight: 600;
  line-height: 1.2;
  white-space: normal;
}
.market-map-kpi-copy > small:first-child {
  display: -webkit-box;
  min-height: 31px;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.market-map-kpi-copy strong {
  display: flex;
  align-items: baseline;
  gap: 6px;
  color: #0f172a;
  font-size: clamp(1.15rem, 1.45vw, 1.55rem);
  font-weight: 800;
  line-height: 1.1;
  white-space: nowrap;
}

.market-map-kpi-copy em {
  padding: 3px 7px;
  border-radius: 7px;
  background: #eff6ff;
  color: #2563eb;
  font-size: 0.75rem;
  font-style: normal;
  font-weight: 700;
}

.market-map-kpi.is-available .market-map-kpi-icon,
.market-map-kpi.is-available .market-map-kpi-copy em { background: #ecfdf5; color: #059669; }
.market-map-kpi.is-reserved .market-map-kpi-icon,
.market-map-kpi.is-reserved .market-map-kpi-copy em { background: #fff7ed; color: #ea580c; }
.market-map-kpi.is-occupied .market-map-kpi-icon,
.market-map-kpi.is-occupied .market-map-kpi-copy em { background: #f1f5f9; color: #475569; }
.market-map-kpi.is-action-required {
  background: #eef2ff;
  border: 1px solid #c7d2fe;
  box-shadow: 0 8px 24px rgba(79, 70, 229, 0.14);
}
.market-map-kpi.is-action-required:hover {
  background: #e0e7ff;
  border-color: #818cf8;
  box-shadow: 0 10px 28px rgba(79, 70, 229, 0.2);
}
.market-map-kpi.is-action-required .market-map-kpi-icon {
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
  color: #ffffff;
  box-shadow: 0 4px 10px rgba(79, 70, 229, 0.28);
}
.market-map-kpi.is-action-required .market-map-kpi-copy small {
  color: #4338ca;
  font-weight: 700;
}
.market-map-kpi.is-action-required .market-map-kpi-copy strong {
  color: #312e81;
}
.market-map-kpi.is-action-required .market-map-kpi-copy em {
  background: #ffffff;
  color: #4338ca;
  border: 1px solid #c7d2fe;
}
.market-map-kpi.is-action-required.active {
  border: 2px solid #4f46e5;
  background: #e0e7ff;
  box-shadow: 0 10px 26px rgba(79, 70, 229, 0.26);
}
.market-map-kpi.is-legal-case .market-map-kpi-icon,
.market-map-kpi.is-legal-case .market-map-kpi-copy em { background: #fee2e2; color: #dc2626; }
.market-map-kpi.is-overdue .market-map-kpi-icon,
.market-map-kpi.is-overdue .market-map-kpi-copy em { background: #fef2f2; color: #dc2626; }
.market-map-kpi.is-expiring .market-map-kpi-icon,
.market-map-kpi.is-expiring .market-map-kpi-copy em { background: #f5f3ff; color: #7c3aed; }

/* Market Map KPI Skeleton Loading */
.market-map-kpi.is-loading {
  cursor: wait;
}

.market-map-kpi-skeleton {
  display: flex;
  align-items: baseline;
  gap: 6px;
  height: 24px;
  margin-top: 3px;
}

.kpi-skeleton-num {
  display: inline-block;
  width: 52px;
  height: 20px;
  border-radius: 6px;
  background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
  background-size: 200% 100%;
  animation: kpiSkeletonShimmer 1.35s ease-in-out infinite;
}

.kpi-skeleton-pill {
  display: inline-block;
  width: 32px;
  height: 18px;
  border-radius: 6px;
  background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
  background-size: 200% 100%;
  animation: kpiSkeletonShimmer 1.35s ease-in-out infinite;
}

@keyframes kpiSkeletonShimmer {
  0% { background-position: 100% 0; }
  100% { background-position: -100% 0; }
}

@media (prefers-reduced-motion: reduce) {
  .kpi-skeleton-num,
  .kpi-skeleton-pill {
    animation: none;
  }
}

/* Compact vertical map tools on the right: search, filters, and demand view. */
.market-map-tool-rail {
  position: absolute;
  top: 50%;
  right: 0;
  z-index: 115;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  padding: 4px 0;
  border: 1px solid #e2e8f0;
  border-right: 0;
  border-radius: 14px 0 0 14px;
  background: rgba(255, 255, 255, 0.97);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.14);
  transform: translateY(-50%);
  transition: opacity 0.22s cubic-bezier(0.16, 1, 0.3, 1), transform 0.26s cubic-bezier(0.16, 1, 0.3, 1);
}

.market-map-tool-rail.has-drawer-open {
  opacity: 0;
  pointer-events: none;
  transform: translateY(-50%) translateX(24px);
}

.market-map-tool-rail .gis-map-toolbar :deep(.master-map-top-bar),
.market-map-tool-rail .gis-map-toolbar :deep(.master-map-tools-right) {
  width: auto;
  margin: 0;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 0;
}

.market-map-tool-rail .gis-map-toolbar :deep(.search-input-group),
.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-filter-wrap) {
  width: 52px;
  height: 52px;
}

.market-map-tool-rail .gis-map-toolbar :deep(.search-input-group) {
  transition: width 0.18s ease;
}

.market-map-tool-rail .gis-map-toolbar :deep(.search-input-group:focus-within) {
  position: absolute;
  top: 4px;
  right: 60px;
  width: min(320px, calc(100vw - 100px));
  filter: drop-shadow(0 8px 18px rgba(15, 23, 42, 0.14));
}

.market-map-tool-rail .gis-map-toolbar :deep(.map-search-input) {
  width: 100%;
  height: 52px;
  padding: 8px 10px 8px 50px;
  border: 0;
  border-radius: 10px 0 0 0;
  background: transparent;
  color: transparent;
  cursor: pointer;
}

.market-map-tool-rail .gis-map-toolbar :deep(.search-input-group:focus-within .map-search-input) {
  border: 1px solid #bfdbfe;
  border-radius: 10px;
  background: #ffffff;
  color: #1e293b;
  cursor: text;
}

.market-map-tool-rail .gis-map-toolbar :deep(.map-search-input::placeholder) { color: transparent; }
.market-map-tool-rail .gis-map-toolbar :deep(.search-input-group:focus-within .map-search-input::placeholder) { color: #94a3b8; }
.market-map-tool-rail .gis-map-toolbar :deep(.search-icon) { left: 17px; color: #334155; }

.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-filter-btn),
.market-map-view-button {
  width: 52px;
  height: 52px;
  padding: 0;
  display: grid;
  place-items: center;
  overflow: hidden;
  border: 0;
  border-top: 1px solid #eef2f6;
  border-radius: 0;
  background: transparent;
  color: transparent;
  cursor: pointer;
}

.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-filter-btn svg),
.market-map-view-button svg {
  width: 20px;
  height: 20px;
  min-width: 20px;
  color: #334155;
}

.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-filter-btn:hover),
.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-filter-btn.active),
.market-map-view-button:hover,
.market-map-view-button.active {
  background: #eff6ff;
}

.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-filter-panel) {
  top: 0;
  right: 62px;
}

.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-count) {
  position: absolute;
  top: 5px;
  right: 5px;
  color: #ffffff;
}

.market-map-view-button.active svg { color: #2563eb; }

.is-fullscreen-map-panel :deep(.map-project-guidance),
.is-fullscreen-map-panel :deep(.map-floating-floor-preview-pill) {
  top: var(--mp-top) !important;
}

@media (max-width: 1200px) {
  .market-map-kpi-strip {
    display: none !important;
  }
  .ov-dash-responsive-kpis {
    display: block !important;
  }
  .manage-page.is-map-page {
    --mp-top: 76px !important;
    --mp-top-2row: 76px !important;
  }
  .is-fullscreen-map-panel :deep(.map-project-guidance),
  .is-fullscreen-map-panel :deep(.map-floating-floor-preview-pill) {
    top: 76px !important;
  }
}

@media (min-width: 769px) and (max-width: 1200px) {
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer) {
    top: 76px !important;
    max-height: calc(100% - 92px) !important;
  }
  .is-fullscreen-map-panel :deep(.stall-mode-panel) {
    top: 76px !important;
    max-height: calc(100% - 92px) !important;
  }
}

@media (max-width: 768px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel.stall-mode-panel),
  .is-fullscreen-map-panel :deep(.stall-mode-panel) {
    position: fixed !important;
    left: 0 !important;
    right: 0 !important;
    bottom: 0 !important;
    top: auto !important;
    width: 100% !important;
    max-width: 100% !important;
    height: 68vh !important;
    max-height: calc(100dvh - 58px) !important;
    border-radius: 20px 20px 0 0 !important;
    border: 1px solid #e2e8f0 !important;
    border-bottom: 0 !important;
    box-shadow: 0 -8px 32px rgba(15, 23, 42, 0.18) !important;
    z-index: 200 !important;
    display: flex !important;
    flex-direction: column !important;
    background: #ffffff !important;
    padding: 0 !important;
    overflow: hidden !important;
    box-sizing: border-box !important;
    transition: height 0.28s cubic-bezier(0.16, 1, 0.3, 1), max-height 0.28s cubic-bezier(0.16, 1, 0.3, 1) !important;
  }

  .is-fullscreen-map-panel :deep(.stall-mode-panel.is-sheet-expanded) {
    height: calc(100dvh - 58px) !important;
    max-height: calc(100dvh - 58px) !important;
  }

  /* บนโหมด responsive ซ่อนแถบเครื่องมือ/ตัวกรอง ov-rail ไม่ให้ลอยบังแผง slide up */
  .is-fullscreen-map-panel :deep(.ov-rail.is-project-active) {
    display: none !important;
    opacity: 0 !important;
    pointer-events: none !important;
  }

  /* บนโหมด responsive ซ่อนปุ่มย้อนกลับและปุ่มลูกศร ใช้การรูดแผงลงเพื่อปิดกลับผังรวมแทน */
  .is-fullscreen-map-panel :deep(.smp-desktop-only),
  .is-fullscreen-map-panel :deep(.smp-back-section),
  .is-fullscreen-map-panel :deep(.smp-expand-toggle-btn) {
    display: none !important;
  }

  /* ขีดปรับขนาดตรึงอยู่ด้านบนสุดตลอดเวลา เป็นหัวกล่องขาวทึบ 100% ไม่โปร่งใส ไม่เห็นพื้นหลังทะลุ */
  .is-fullscreen-map-panel :deep(.smp-drag-handle-bar) {
    flex: 0 0 auto !important;
    display: flex !important;
    justify-content: center !important;
    align-items: center !important;
    width: 100% !important;
    padding: 12px 0 10px !important;
    margin: 0 !important;
    background: #ffffff !important;
    border-bottom: 1px solid #f1f5f9 !important;
    border-radius: 20px 20px 0 0 !important;
    cursor: grab !important;
    touch-action: none !important;
    user-select: none !important;
    -webkit-user-select: none !important;
  }

  .is-fullscreen-map-panel :deep(.smp-scroll-content) {
    flex: 1 1 auto !important;
    min-height: 0 !important;
    overflow-y: auto !important;
    -webkit-overflow-scrolling: touch !important;
    padding: 12px 18px 28px !important;
    display: flex !important;
    flex-direction: column !important;
    gap: 12px !important;
    background: #ffffff !important;
  }

  /* หน้าต่างรายละเอียดแผงบนโหมด Responsive: เปิดเต็มหน้าจอใต้ Navbar พร้อมแถบแท็บแนวนอน */
  .inspect-selected-stall-drawer,
  .inspect-selected-stall-drawer.is-full-detail-mode,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer),
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock),
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-full-detail-mode) {
    position: fixed !important;
    top: 58px !important;
    left: 0 !important;
    right: 0 !important;
    bottom: 0 !important;
    width: 100vw !important;
    max-width: 100vw !important;
    height: calc(100vh - 58px) !important;
    height: calc(100dvh - 58px) !important;
    max-height: calc(100vh - 58px) !important;
    max-height: calc(100dvh - 58px) !important;
    border-radius: 0 !important;
    border: none !important;
    box-shadow: none !important;
    z-index: 1000 !important;
    display: flex !important;
    flex-direction: column !important;
    background: #f8fafc !important;
    overflow: hidden !important;
  }

  /* ส่วนหัวคงที่ด้านบน (Sticky Header) */
  .inspect-selected-stall-drawer .issd-header,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-header) {
    display: flex !important;
    flex-direction: column !important;
    gap: 8px !important;
    padding: 14px 16px 10px !important;
    background: #ffffff !important;
    border-bottom: 1px solid #eef2f6 !important;
    border-right: none !important;
    flex: 0 0 auto !important;
  }

  .inspect-selected-stall-drawer .issd-return-row,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-return-row) {
    display: flex !important;
    align-items: center !important;
    justify-content: flex-start !important;
    padding: 0 0 12px !important;
    margin: 0 0 4px !important;
    border-bottom: 1px solid #e8eef6 !important;
    width: 100% !important;
  }

  .inspect-selected-stall-drawer .issd-return-btn,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-return-btn) {
    display: inline-flex !important;
    align-items: center !important;
    gap: 8px !important;
    min-height: 42px !important;
    height: 42px !important;
    padding: 0 16px 0 12px !important;
    border-radius: 10px !important;
    border: 1px solid #bfdbfe !important;
    background: #eff6ff !important;
    color: #2563eb !important;
    font-size: 0.875rem !important;
    font-weight: 600 !important;
    cursor: pointer !important;
    box-shadow: none !important;
    margin: 0 !important;
    -webkit-tap-highlight-color: transparent !important;
  }

  .inspect-selected-stall-drawer .issd-return-btn:hover,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-return-btn:hover) {
    background: #dbeafe !important;
  }

  .issd-mobile-only {
    display: grid !important;
    place-items: center !important;
    width: 38px !important;
    height: 38px !important;
    border-radius: 10px !important;
    border: 1px solid #e2e8f0 !important;
    background: #ffffff !important;
    color: #475569 !important;
    cursor: pointer !important;
  }
  .issd-mobile-only:hover {
    background: #f1f5f9 !important;
    color: #0f172a !important;
  }

  .inspect-selected-stall-drawer .issd-header-top,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-header-top) {
    display: flex !important;
    align-items: center !important;
    justify-content: space-between !important;
    width: 100% !important;
  }

  .inspect-selected-stall-drawer .issd-title-group,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-title-group) {
    display: flex !important;
    align-items: center !important;
    gap: 8px !important;
    flex-wrap: wrap !important;
  }

  .inspect-selected-stall-drawer .issd-title-group strong,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-title-group strong) {
    font-size: 1.15rem !important;
    font-weight: 700 !important;
    color: #0f172a !important;
  }

  .inspect-selected-stall-drawer .issd-project-sub,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-project-sub) {
    font-size: 0.8125rem !important;
    color: #64748b !important;
    margin: 0 !important;
  }

  /* โครงสร้างเนื้อหาเต็มจอ */
  .inspect-selected-stall-drawer .issd-full-layout,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-full-layout) {
    display: flex !important;
    flex-direction: column !important;
    flex: 1 1 auto !important;
    min-height: 0 !important;
    overflow: hidden !important;
  }

  /* แถบแท็บแนวนอน (Horizontal Scrollable Tabs) */
  .inspect-selected-stall-drawer .issd-full-sidebar,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-full-sidebar) {
    width: 100% !important;
    flex: 0 0 auto !important;
    background: #ffffff !important;
    border-bottom: 1px solid #e2e8f0 !important;
    border-right: none !important;
    padding: 0 12px !important;
    margin: 0 !important;
    overflow-x: auto !important;
    -webkit-overflow-scrolling: touch !important;
    scrollbar-width: none !important;
  }

  .inspect-selected-stall-drawer .issd-full-sidebar::-webkit-scrollbar,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-full-sidebar::-webkit-scrollbar) {
    display: none !important;
  }

  .inspect-selected-stall-drawer .issd-sidebar-nav,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-sidebar-nav) {
    display: flex !important;
    flex-direction: row !important;
    align-items: center !important;
    gap: 8px !important;
    padding: 8px 4px !important;
    white-space: nowrap !important;
    width: max-content !important;
    min-width: 100% !important;
  }

  .inspect-selected-stall-drawer .issd-sidebar-tab,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-sidebar-tab) {
    flex: 0 0 auto !important;
    width: auto !important;
    min-height: 42px !important;
    height: 42px !important;
    padding: 0 16px !important;
    border-radius: 999px !important;
    border: 1px solid transparent !important;
    background: #f1f5f9 !important;
    color: #475569 !important;
    font-size: 0.875rem !important;
    font-weight: 600 !important;
    display: inline-flex !important;
    align-items: center !important;
    gap: 8px !important;
    cursor: pointer !important;
    transition: all 0.15s ease !important;
    -webkit-tap-highlight-color: transparent !important;
  }

  .inspect-selected-stall-drawer .issd-sidebar-tab:hover,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-sidebar-tab:hover) {
    background: #e2e8f0 !important;
    color: #1e293b !important;
  }

  .inspect-selected-stall-drawer .issd-sidebar-tab.active,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-sidebar-tab.active) {
    background: #2563eb !important;
    color: #ffffff !important;
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.25) !important;
  }

  .inspect-selected-stall-drawer .issd-sidebar-tab.active svg,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-sidebar-tab.active svg) {
    stroke: #ffffff !important;
  }

  /* คอลัมน์เนื้อหาหลัก เลื่อนแนวตั้งได้เต็มพื้นที่ */
  .inspect-selected-stall-drawer .issd-full-main,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-full-main) {
    flex: 1 1 auto !important;
    min-height: 0 !important;
    overflow-y: auto !important;
    -webkit-overflow-scrolling: touch !important;
    display: flex !important;
    flex-direction: column !important;
    background: #f8fafc !important;
  }

  .inspect-selected-stall-drawer .issd-action-alert-card,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-action-alert-card) {
    margin: 12px 14px 4px !important;
  }

  .inspect-selected-stall-drawer .issd-body-full,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-body-full) {
    padding: 12px 14px 24px !important;
    gap: 12px !important;
  }

  /* การ์ดสเปกและรายละเอียดต่างๆ ให้เต็มความกว้างและโปร่งสบาย */
  .inspect-selected-stall-drawer .issd-card-section,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-card-section) {
    padding: 14px 16px !important;
    border-radius: 14px !important;
  }

  /* ขนาดปุ่มในบัตรรายการ / แท็บต่างๆ บนมือถือ ให้อ่านง่าย แตะสะดวก ไม่กดพลาด */
  .inspect-selected-stall-drawer .issd-card-actions-row button,
  .inspect-selected-stall-drawer .issd-card-actions-row a,
  .inspect-selected-stall-drawer .issd-booking-direct-btn,
  .inspect-selected-stall-drawer .issd-occupant-action-btn,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-card-actions-row button),
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-card-actions-row a),
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-booking-direct-btn),
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-occupant-action-btn) {
    min-height: 42px !important;
    padding: 8px 16px !important;
    font-size: 0.875rem !important;
    font-weight: 600 !important;
    border-radius: 10px !important;
    display: inline-flex !important;
    align-items: center !important;
    justify-content: center !important;
    gap: 6px !important;
  }

  .inspect-selected-stall-drawer .issd-full-main .app-segmented-option,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-full-main .app-segmented-option) {
    min-height: 38px !important;
    padding: 6px 14px !important;
    font-size: 0.875rem !important;
    border-radius: 8px !important;
  }

  /* Footer ปุ่มคำสั่ง ตรึงด้านล่างสุดของหน้าจอ รองรับ safe-area บนมือถือ */
  .inspect-selected-stall-drawer .issd-footer-full,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-footer-full) {
    position: sticky !important;
    bottom: 0 !important;
    background: #ffffff !important;
    border-top: 1px solid #e2e8f0 !important;
    padding: 12px 16px calc(12px + env(safe-area-inset-bottom, 0px)) !important;
    box-shadow: 0 -4px 16px rgba(15, 23, 42, 0.06) !important;
    z-index: 10 !important;
    display: flex !important;
    flex-direction: column !important;
    gap: 10px !important;
    width: 100% !important;
    box-sizing: border-box !important;
  }

  .inspect-selected-stall-drawer .issd-dual-btn-group,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-dual-btn-group) {
    display: flex !important;
    flex-direction: column !important;
    gap: 10px !important;
    width: 100% !important;
    box-sizing: border-box !important;
  }

  .inspect-selected-stall-drawer .issd-dual-btn-group > *,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-dual-btn-group > *) {
    width: 100% !important;
    flex: 0 0 auto !important;
  }

  .inspect-selected-stall-drawer .issd-footer-full button,
  .inspect-selected-stall-drawer .issd-footer-full .app-button,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-footer-full button),
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-footer-full .app-button) {
    min-height: 48px !important;
    height: 48px !important;
    font-size: 0.9375rem !important;
    font-weight: 600 !important;
    border-radius: 12px !important;
    width: 100% !important;
    display: inline-flex !important;
    align-items: center !important;
    justify-content: center !important;
    gap: 8px !important;
    box-sizing: border-box !important;
  }
}

@media (max-width: 720px) {
  .manage-page.is-map-page { --mp-top: 76px !important; --mp-top-2row: 76px !important; }
  .market-map-kpi-strip {
    display: none !important;
  }
  .market-map-tool-rail { right: 10px; border-right: 1px solid #e2e8f0; border-radius: 14px; }
  .is-fullscreen-map-panel :deep(.map-project-guidance),
  .is-fullscreen-map-panel :deep(.map-floating-floor-preview-pill) { top: 76px !important; }
}

/* Mini Drawer: Urgent Action Alert Banner & Expandable Card */
.issd-mini-alert-card {
  width: 100%;
  box-sizing: border-box;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 11px;
  overflow: hidden;
  transition: all 0.2s ease;
}
.issd-mini-alert-card .issd-mini-alert-banner {
  width: 100%;
  box-sizing: border-box;
  padding: 10px 13px;
  background: transparent;
  border: none;
  border-radius: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  cursor: pointer;
  font-family: inherit;
  text-align: left;
  transition: background 0.15s ease;
}
.issd-mini-alert-card .issd-mini-alert-banner:hover {
  background: #fef3c7;
}
.issd-mini-alert-head {
  display: flex;
  align-items: center;
  gap: 7px;
  color: #b45309;
  font-size: 0.8125rem;
  font-weight: 700;
  min-width: 0;
}
.issd-mini-alert-head svg {
  color: #d97706;
  flex-shrink: 0;
}
.issd-mini-alert-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  color: #d97706;
  flex-shrink: 0;
}
.issd-mini-alert-link svg {
  transition: transform 0.2s ease;
}
.issd-mini-alert-link svg.is-flipped {
  transform: rotate(180deg);
}
.issd-mini-alert-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 6px 10px 10px;
  border-top: 1px dashed #fde68a;
  background: #fffdf5;
}
.issd-mini-alert-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 8px 10px;
  background: #ffffff;
  border: 1px solid #fef3c7;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
}
.issd-mini-alert-row:hover {
  background: #fefce8;
  border-color: #fde047;
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(217, 119, 6, 0.08);
}
.issd-mini-alert-row-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex: 1;
}
.issd-mini-alert-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}
.issd-mini-alert-text strong {
  font-size: 0.8125rem;
  color: #1e293b;
  font-weight: 600;
  line-height: 1.25;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.issd-mini-alert-desc {
  font-size: 0.75rem;
  color: #64748b;
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.issd-mini-alert-action-pill {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 0.75rem;
  font-weight: 600;
  color: #d97706;
  white-space: nowrap;
  flex-shrink: 0;
  padding: 3px 8px;
  border-radius: 6px;
  background: #fef3c7;
  border: 1px solid #fde68a;
  transition: all 0.15s ease;
}
.issd-mini-alert-action-pill.is-danger {
  color: #dc2626;
  background: #fee2e2;
  border-color: #fca5a5;
}
.issd-mini-alert-action-pill.is-warning {
  color: #d97706;
  background: #fef3c7;
  border-color: #fde68a;
}
.issd-mini-alert-action-pill.is-info {
  color: #2563eb;
  background: #eff6ff;
  border-color: #bfdbfe;
}

/* Mini Drawer: Primary Occupant Card */
.issd-card-section {
  padding: 14px 16px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  gap: 10px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}
.issd-card-section.is-occupied {
  border-color: #cbd5e1;
  background: #f8fafc;
}
.issd-card-section.is-reserved {
  border-color: #fde68a;
  background: #fffdf5;
}
.issd-card-section.is-available {
  border-color: #bbf7d0;
  background: #f0fdf4;
}
.issd-section-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.issd-section-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8125rem;
  font-weight: 600;
  color: #64748b;
}
.issd-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 9px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  line-height: 1.2;
}
.issd-badge.badge-active,
.issd-badge.badge-success {
  background: #dcfce7;
  color: #15803d;
}
.issd-badge.badge-warning {
  background: #fef3c7;
  color: #b45309;
}
.issd-occupant-name {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  min-width: 0;
}
.issd-occupant-name strong {
  font-size: 0.9375rem;
  font-weight: 700;
  color: #0f172a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-phone-link {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 3px 9px;
  border-radius: 7px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  color: #2563eb;
  font-size: 0.75rem;
  font-weight: 600;
  text-decoration: none;
  flex-shrink: 0;
  transition: background 0.15s ease, border-color 0.15s ease;
}
.issd-phone-link:hover {
  background: #dbeafe;
  border-color: #93c5fd;
}
.issd-occupant-meta {
  display: flex;
  flex-direction: column;
  gap: 7px;
  padding-top: 10px;
  border-top: 1px solid #e2e8f0;
}
.issd-meta-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  font-size: 0.8125rem;
}
.issd-meta-item .meta-k {
  color: #64748b;
  font-weight: 500;
  flex-shrink: 0;
}
.issd-meta-item .meta-v {
  color: #1e293b;
  font-weight: 600;
  text-align: right;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-meta-item .code-font {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  color: #2563eb;
}

.issd-contract-code-pill {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.75rem;
  font-weight: 700;
  color: #2563eb;
  background: #eff6ff;
  border: 1px solid #dbeafe;
  padding: 2px 8px;
  border-radius: 6px;
}

/* Separated Person Row + Contract/Booking Sub-box inside Occupant Card */
.issd-person-row {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}
.issd-person-avatar {
  width: 34px;
  height: 34px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #dbeafe;
}
.issd-person-avatar.is-reserved {
  background: #fef3c7;
  color: #d97706;
  border-color: #fde68a;
}
.issd-person-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
  flex: 1;
}
.issd-person-role {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  line-height: 1.2;
}
.issd-person-name {
  font-size: 0.9375rem;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-subdoc-box {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 9px 11px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 9px;
}
.issd-subdoc-box.is-reserved {
  border-color: #fde68a;
  background: #ffffff;
}
.issd-subdoc-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  font-size: 0.75rem;
  line-height: 1.35;
}
.issd-subdoc-row .subdoc-k {
  color: #64748b;
  font-weight: 500;
  flex-shrink: 0;
}
.issd-subdoc-row .subdoc-v {
  color: #1e293b;
  font-weight: 600;
  text-align: right;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-subdoc-row .subdoc-v.is-code {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-weight: 700;
  color: #2563eb;
}
.issd-subdoc-row .subdoc-duration {
  font-style: normal;
  font-weight: 600;
  color: #2563eb;
  margin-left: 2px;
}
.issd-occupant-compact-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding-top: 8px;
  border-top: 1px solid #e2e8f0;
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 500;
}
.issd-duration-tag {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  background: #eff6ff;
  color: #2563eb;
  font-size: 0.75rem;
  font-weight: 700;
  white-space: nowrap;
}
.issd-sep-dot {
  margin: 0 5px;
  color: #cbd5e1;
}

/* Mini Drawer: 4-KPI Operational Grid */
.issd-mini-ops-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}
.issd-mini-ops-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 11px;
  min-width: 0;
  cursor: default;
  user-select: none;
  font-family: inherit;
  text-align: left;
}
.issd-mini-ops-card .ops-icon {
  width: 32px;
  height: 32px;
  border-radius: 9px;
  background: #f1f5f9;
  color: #64748b;
  display: grid;
  place-items: center;
  flex-shrink: 0;
}
.issd-mini-ops-card .ops-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.issd-mini-ops-card .ops-label {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 500;
  line-height: 1.2;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-mini-ops-card .ops-info strong {
  font-size: 0.8125rem;
  color: #1e293b;
  font-weight: 700;
  line-height: 1.25;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-mini-ops-card .ops-info strong.text-success {
  color: #15803d;
}
.issd-mini-ops-card .ops-info strong.text-warning {
  color: #d97706;
}
.issd-mini-ops-card .ops-info strong.text-danger {
  color: #dc2626;
}
.issd-mini-ops-card.is-ok .ops-icon {
  background: #ecfdf5;
  color: #10b981;
}
.issd-mini-ops-card.is-info-tint .ops-icon {
  background: #eff6ff;
  color: #2563eb;
}
.issd-mini-ops-card.is-alert {
  background: #fef2f2;
  border-color: #fecaca;
}
.issd-mini-ops-card.is-alert .ops-icon {
  background: #fee2e2;
  color: #dc2626;
}
.issd-mini-ops-card.is-warning {
  background: #fffbeb;
  border-color: #fde68a;
}
.issd-mini-ops-card.is-warning .ops-icon {
  background: #fef3c7;
  color: #d97706;
}
.issd-mini-ops-grid.is-full-view {
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  margin-top: 14px;
}
@media (max-width: 768px) {
  .issd-mini-ops-grid.is-full-view {
    grid-template-columns: repeat(2, 1fr);
  }
}
.issd-mini-ops-card.is-clickable {
  cursor: pointer;
  transition: transform 0.16s cubic-bezier(0.2, 0.8, 0.2, 1), box-shadow 0.16s cubic-bezier(0.2, 0.8, 0.2, 1), border-color 0.16s ease;
}
.issd-mini-ops-card.is-clickable:hover {
  transform: translateY(-2px);
  border-color: #cbd5e1;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.07);
}
.issd-mini-ops-card.is-clickable:active {
  transform: translateY(0);
}

/* Floating Filter Bar & Quick Stats in Sleek Pill/Chip Style (Non-overlapping) */
.map-floating-overlay-bar {
  position: absolute;
  top: 16px;
  left: 18px;
  right: 18px;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  pointer-events: none;
}

.map-floating-overlay-bar > * {
  pointer-events: auto;
}

.map-overlay-left-cluster {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  max-width: calc(100% - 310px);
}

/* Mode Switcher Segmented Pill */
.map-view-toggle-pill {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(226, 232, 240, 0.95);
  border-radius: 12px;
  padding: 4px;
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
  height: 40px;
  box-sizing: border-box;
  flex-shrink: 0;
}

.view-toggle-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #475569;
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.18s ease;
  white-space: nowrap;
}

.view-toggle-btn:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.view-toggle-btn.active {
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.25);
}

.view-toggle-btn.is-demand.active {
  background: #ea580c;
  color: #ffffff;
  box-shadow: 0 2px 6px rgba(234, 88, 12, 0.25);
}

.demand-live-pulse-dot {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: #fbbf24;
  box-shadow: 0 0 0 2px rgba(251, 191, 36, 0.4);
  display: inline-block;
}

/* Compact Filter Chips Row */
.map-floating-chips-group {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(226, 232, 240, 0.95);
  border-radius: 12px;
  padding: 4px 6px;
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
  height: 40px;
  box-sizing: border-box;
  overflow-x: auto;
  max-width: 100%;
}

.map-filter-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;

  border-radius: 8px;
  border: 1px solid transparent;
  background: transparent;
  color: #475569;
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.map-filter-chip:hover {
  background: #f8fafc;
  color: #0f172a;
}

.chip-label {
  font-size: 0.8125rem;
}

.chip-badge {
  font-size: 0.8125rem;
  font-weight: 700;
  color: #0f172a;
}

.chip-color-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  flex-shrink: 0;
}

.chip-color-dot.is-available { background: #10b981; }
.chip-color-dot.is-reserved { background: #f59e0b; }
.chip-color-dot.is-occupied { background: #64748b; }
.chip-color-dot.is-high { background: #dc2626; }
.chip-color-dot.is-medium { background: #f97316; }
.chip-color-dot.is-low { background: #eab308; }
.chip-color-dot.is-none { background: #10b981; }

.map-filter-chip.active {
  background: #eff6ff;
  border-color: #bfdbfe;
  color: #1d4ed8;
}

.map-filter-chip.chip-demand-high.active {
  background: #fef2f2;
  border-color: #fca5a5;
  color: #b91c1c;
}

.map-filter-chip.chip-demand-medium.active {
  background: #fff7ed;
  border-color: #fed7aa;
  color: #c2410c;
}

.map-filter-chip.chip-demand-low.active {
  background: #fefce8;
  border-color: #fef08a;
  color: #a16207;
}

.map-filter-chip.chip-demand-none.active {
  background: #f0fdf4;
  border-color: #bbf7d0;
  color: #15803d;
}

.map-floating-overlay-bar .gis-map-toolbar {
  margin-left: auto;
  flex-shrink: 0;
}

.map-floating-overlay-bar .gis-map-toolbar :deep(.map-search-input) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(226, 232, 240, 0.9);
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
}

.map-floating-overlay-bar .gis-map-toolbar :deep(.map-advanced-filter-btn) {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(226, 232, 240, 0.9);
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
}

@media (max-width: 960px) {
  .map-floating-overlay-bar {
    position: static;
    padding: 10px 12px;
    background: #ffffff;
    border-bottom: 1px solid #e2e8f0;
    pointer-events: auto;
    flex-direction: column;
    align-items: stretch;
  }
  .map-overlay-left-cluster {
    max-width: 100%;
  }
  .map-floating-overlay-bar .gis-map-toolbar {
    width: 100%;
    margin-left: 0;
  }
}

/* ==========================================================================
   Stall Drawer Slide-Fade Animation
   ========================================================================== */
.slide-fade-enter-active {
  transition: transform 0.28s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.24s ease-out;
}

.slide-fade-leave-active {
  transition: transform 0.22s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.2s ease-in;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(40px) scale(0.97);
  opacity: 0;
}

@media (max-width: 768px) {
  .slide-fade-enter-active {
    transition: transform 0.32s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.26s ease-out;
  }

  .slide-fade-leave-active {
    transition: transform 0.24s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.22s ease-in;
  }

  .slide-fade-enter-from,
  .slide-fade-leave-to {
    transform: translateY(100%);
    opacity: 0;
  }
}

/* การ์ดผู้เช่า/ผู้จอง: แถวข้อมูลเป็นหัวข้อบน ค่าล่าง ชิดซ้าย ไม่ให้ค่ายาวๆ เบียดกัน */
.inspect-selected-stall-drawer .issd-subdoc-box { gap: 12px; padding: 12px 14px; }
.inspect-selected-stall-drawer .issd-subdoc-row { flex-direction: column; flex-wrap: nowrap; align-items: flex-start; gap: 2px; }
.inspect-selected-stall-drawer .issd-subdoc-row .subdoc-k { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-subdoc-row .subdoc-v { margin-left: 0; text-align: left; white-space: normal; overflow: visible; overflow-wrap: anywhere; }
.inspect-selected-stall-drawer .issd-subdoc-row .subdoc-v.is-code { font-size: 0.9375rem; letter-spacing: 0.01em; }
.inspect-selected-stall-drawer .issd-subdoc-row .subdoc-duration { margin-left: 6px; }

/* การ์ดสรุปการเงิน/สัญญา/งานซ่อม/จอง: ขนาดเท่ากัน ข้อความบรรทัดเดียว ไม่ตกบรรทัด */
.inspect-selected-stall-drawer .issd-mini-ops-grid { grid-auto-rows: 1fr; gap: 10px; }
.inspect-selected-stall-drawer .issd-mini-ops-card { min-height: 66px; padding: 10px 12px; gap: 10px; align-items: center; }
.inspect-selected-stall-drawer .issd-mini-ops-card .ops-icon { width: 36px; height: 36px; flex: 0 0 36px; }
.inspect-selected-stall-drawer .issd-mini-ops-card .ops-info { flex: 1 1 auto; gap: 2px; }
.inspect-selected-stall-drawer .issd-mini-ops-card .ops-label { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-mini-ops-card .ops-info strong { font-size: 0.9375rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }


/* ตัวเลขจำนวนรายการในป้ายหมวด */
.issd-aac-tag { display: inline-flex; align-items: center; gap: 6px; }
.issd-aac-tag-count { display: inline-grid; place-items: center; min-width: 18px; height: 18px; padding: 0 5px; border-radius: 9999px; font-size: 0.75rem; font-weight: 700; line-height: 1; color: #ffffff; }
.issd-aac-tag.is-danger .issd-aac-tag-count { background: #dc2626; }
.issd-aac-tag.is-warning .issd-aac-tag-count { background: #d97706; }
.issd-aac-tag.is-info .issd-aac-tag-count { background: #2563eb; }

/* การ์ดสิ่งที่ต้องดำเนินการ (หน้ารายละเอียดเต็ม): โปร่ง อ่านง่าย แถบสีบอกระดับ ปุ่มลิงก์เป็นกรอบ */
.issd-full-main .issd-action-alert-card { background: #ffffff; border: 1px solid #fde68a; border-radius: 14px; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04); }
.issd-full-main .issd-aac-header { background: #fffbeb; padding: 12px 16px; }
.issd-full-main .issd-aac-row { padding: 12px 16px; gap: 12px; border-top: 1px solid #f1f5f9; }
.issd-full-main .issd-aac-row-left { display: flex; align-items: center; gap: 12px; min-width: 0; }
.issd-aac-bar { flex: 0 0 4px; align-self: stretch; min-height: 34px; border-radius: 9999px; background: #94a3b8; }
.issd-aac-bar.is-danger { background: #dc2626; }
.issd-aac-bar.is-warning { background: #f59e0b; }
.issd-aac-bar.is-info { background: #2563eb; }
.issd-aac-row-text { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.issd-aac-row-text strong { display: inline-flex; align-items: center; gap: 8px; color: #1e293b; font-size: 1rem; font-weight: 700; line-height: 1.3; }
.issd-aac-row-text small { color: #64748b; font-size: 0.875rem; font-weight: 500; line-height: 1.3; }
.issd-aac-count { display: inline-grid; place-items: center; min-width: 22px; height: 22px; padding: 0 6px; border-radius: 9999px; color: #ffffff; font-size: 0.8125rem; font-weight: 700; line-height: 1; background: #94a3b8; }
.issd-aac-count.is-danger { background: #dc2626; }
.issd-aac-count.is-warning { background: #d97706; }
.issd-aac-count.is-info { background: #2563eb; }
.issd-full-main .issd-aac-row-right { flex: 0 0 auto; display: inline-flex; align-items: center; gap: 4px; height: 32px; padding: 0 12px; border: 1px solid #93c5fd; border-radius: 7px; background: #ffffff; color: #2563eb; font-size: 0.875rem; font-weight: 600; white-space: nowrap; transition: background 0.15s ease; }
.issd-full-main .issd-aac-row:hover .issd-aac-row-right { background: #eff6ff; }

/* การ์ดสิ่งที่ต้องดำเนินการ (หน้าสรุปแผง): แถวโปร่ง แถบสีบอกระดับ เหมือนหน้ารายละเอียดเต็ม */
.issd-mini-alert-card { background: #ffffff; border: 1px solid #fde68a; border-radius: 14px; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04); }
.issd-mini-alert-card .issd-mini-alert-banner { background: #fffbeb; padding: 12px 16px; }
.issd-mini-alert-card .issd-mini-alert-list { padding: 0; background: #ffffff; border-top: 1px solid #fde68a; }
.issd-mini-alert-list .issd-mini-alert-row { display: flex; align-items: center; justify-content: flex-start; gap: 12px; padding: 12px 16px; border: 0; border-top: 1px solid #f1f5f9; border-radius: 0; background: #ffffff; }
.issd-mini-alert-list .issd-mini-alert-row:first-child { border-top: 0; }
.issd-mini-alert-list .issd-mini-alert-row:hover { background: #f8fafc; }
.issd-mini-alert-list .issd-aac-bar { flex: 0 0 4px; align-self: stretch; min-height: 34px; }
.issd-mini-alert-list .issd-aac-row-text { flex: 1 1 auto; }
.issd-mini-alert-chevron { flex: 0 0 auto; color: #2563eb; }

/* ตัวกรองสัญญาในแท็บชำระเงิน: ย่อตามความกว้างแผง ไม่ล้นขอบ */
.inspect-selected-stall-drawer .issd-full-main .issd-tab-header-row.is-inline-controls { flex-wrap: wrap; }
.inspect-selected-stall-drawer .issd-full-main .issd-contract-inline-filter { flex: 1 1 220px; min-width: 0; max-width: 100%; justify-content: flex-end; }
.inspect-selected-stall-drawer .issd-full-main .issd-inline-filter-select { flex: 1 1 140px; width: auto; min-width: 0; max-width: 240px; }
.inspect-selected-stall-drawer .issd-full-main .issd-inline-filter-select :deep(.app-select-trigger) { min-width: 0; }

/* สวิตช์โหมด ข้อมูล / ค่าใช้จ่าย: ปรับให้เรียบ สมส่วน ตัวอักษรชัด */
.map-mode-float { padding: 3px; border: 1px solid rgba(203, 213, 225, 0.8); border-radius: 12px; background: rgba(255, 255, 255, 0.97); box-shadow: 0 4px 14px rgba(15, 23, 42, 0.1); }
.map-mode-float :deep(.app-segmented-control) { gap: 2px; padding: 0; border: 0; background: transparent; }
.map-mode-float :deep(.app-segmented-option) { min-width: 116px; min-height: 38px; padding: 0 16px; gap: 8px; justify-content: center; border-radius: 9px; font-size: 0.9375rem; font-weight: 600; color: #475569; }
.map-mode-float :deep(.app-segmented-option svg) { width: 18px; height: 18px; flex: 0 0 18px; }
.map-mode-float :deep(.app-segmented-option:hover:not(.is-active)) { background: #f1f5f9; color: #1e293b; }
.map-mode-float :deep(.app-segmented-option.is-active) { background: #2563eb; color: #ffffff; font-weight: 700; box-shadow: none; }
@media (max-width: 720px) { .map-mode-float :deep(.app-segmented-option) { min-width: 96px; padding: 0 12px; } }



/* สวิตช์โหมดแบบการ์ดเตี้ยเท่าการ์ด KPI: อยู่แถวเดียวกัน ไอคอนบน ชื่อล่าง ตัวที่เลือกเป็นพื้นน้ำเงิน */
.map-mode-float {
  top: 18px;
  left: 20px;
  height: 72px;
  box-sizing: border-box;
  padding: 6px;
  border: 1px solid rgba(226, 232, 240, 0.94);
  border-radius: 15px;
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(12px);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.11);
}
.map-mode-float :deep(.app-segmented-control) { display: flex; gap: 4px; height: 100%; padding: 0; border: 0; border-radius: 0; background: transparent; box-shadow: none; }
.map-mode-float :deep(.app-segmented-control)::before { content: none; }
.map-mode-float :deep(.app-segmented-option) {
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  min-width: 84px;
  min-height: 0;
  height: 100%;
  padding: 0 10px;
  border-radius: 10px;
  background: transparent;
  box-shadow: none;
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 600;
  line-height: 1.1;
}
.map-mode-float :deep(.app-segmented-option svg) { width: 22px; height: 22px; flex: 0 0 22px; }
.map-mode-float :deep(.app-segmented-option:hover:not(.is-active)) { background: #f1f5f9; color: #1e293b; }
.map-mode-float :deep(.app-segmented-option.is-active) { background: #2563eb; color: #ffffff; font-weight: 700; box-shadow: 0 3px 10px rgba(37, 99, 235, 0.3); }
@media (max-width: 720px) {
  .map-mode-float { top: 10px; left: 10px; height: 56px; }
  .map-mode-float :deep(.app-segmented-option) { min-width: 72px; }
  .map-mode-float :deep(.app-segmented-option svg) { width: 18px; height: 18px; flex-basis: 18px; }
  .manage-page.is-map-page { --mp-top: 76px !important; --mp-top-2row: 76px !important; }
}

/* โครงว่างระหว่างรอกล่องขยายจบ */
.issd-full-skeleton { display: flex; flex-direction: column; gap: 14px; padding: 18px 22px; }
.issd-full-skeleton span { display: block; border-radius: 14px; background: linear-gradient(90deg, #eef2f6 25%, #f8fafc 50%, #eef2f6 75%); background-size: 200% 100%; animation: issdSkeleton 1.1s ease-in-out infinite; }
.issd-full-skeleton span:nth-child(1) { height: 56px; }
.issd-full-skeleton span:nth-child(2) { height: 150px; }
.issd-full-skeleton span:nth-child(3) { height: 90px; }
@keyframes issdSkeleton { from { background-position: 200% 0; } to { background-position: -200% 0; } }

/* จอไม่เกิน 1280px: การ์ด KPI กระชับขึ้น (ซ่อนกล่องไอคอน ลดระยะ) ไม่ให้ตัวเลขและป้าย "แผง" ถูกตัด */
@media (max-width: 1280px) {
  .market-map-kpi-strip { gap: 10px; }
  .market-map-kpi { padding: 10px 12px; gap: 8px; }
  .market-map-kpi.active { padding: 9px 11px; }
  .market-map-kpi-icon { display: none; }
  .market-map-kpi-copy strong { font-size: clamp(1.05rem, 2vw, 1.3rem); }
}
</style>

<style>
/* ==========================================================================
   Stall Drawer - Action Alert Card (Modern High-End Task Card)
   Direct unscoped rules scoped strictly to .inspect-selected-stall-drawer to ensure
   100% reliable styling inside child component slots across all browser environments
   ========================================================================== */
.inspect-selected-stall-drawer .issd-action-alert-card,
.issd-action-alert-card {
  margin: 10px 14px 4px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.08);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  overflow: hidden;
  transition: all 0.2s ease;
}

.inspect-selected-stall-drawer .issd-aac-header,
.issd-aac-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 9px 12px;
  background: #fffbeb;
  cursor: pointer;
  user-select: none;
  transition: background 0.15s ease;
}

.inspect-selected-stall-drawer .issd-aac-header:hover,
.issd-aac-header:hover {
  background: #fef3c7;
}

.inspect-selected-stall-drawer .issd-aac-title-group,
.issd-aac-title-group {
  display: flex;
  align-items: center;
  gap: 9px;
  min-width: 0;
}

.inspect-selected-stall-drawer .issd-aac-icon-wrap,
.issd-aac-icon-wrap {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: #fef3c7;
  color: #d97706;
  border: 1px solid #fde68a;
  display: grid;
  place-items: center;
  flex-shrink: 0;
}

.inspect-selected-stall-drawer .issd-aac-title-text,
.issd-aac-title-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.inspect-selected-stall-drawer .issd-aac-title-row,
.issd-aac-title-row {
  display: flex;
  align-items: center;
  gap: 7px;
  flex-wrap: wrap;
}

.inspect-selected-stall-drawer .issd-aac-title-row strong,
.issd-aac-title-row strong {
  font-size: 0.84375rem;
  font-weight: 700;
  color: #92400e;
}

.inspect-selected-stall-drawer .issd-aac-count-badge,
.issd-aac-count-badge {
  display: inline-flex;
  align-items: center;
  padding: 1px 7px;
  border-radius: 999px;
  background: #fde68a;
  color: #78350f;
  font-size: 0.75rem;
  font-weight: 700;
}

.inspect-selected-stall-drawer .issd-aac-sub,
.issd-aac-sub {
  font-size: 0.75rem;
  color: #b45309;
}

.inspect-selected-stall-drawer .issd-aac-toggle-btn,
.issd-aac-toggle-btn {
  width: 26px;
  height: 26px;
  border-radius: 6px;
  border: 1px solid rgba(217, 119, 6, 0.25);
  background: #ffffff;
  color: #92400e;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.inspect-selected-stall-drawer .issd-aac-toggle-btn:hover,
.issd-aac-toggle-btn:hover {
  background: #fef3c7;
  border-color: #d97706;
}

.inspect-selected-stall-drawer .issd-aac-toggle-btn svg,
.issd-aac-toggle-btn svg {
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.inspect-selected-stall-drawer .issd-aac-toggle-btn svg.is-flipped,
.issd-aac-toggle-btn svg.is-flipped {
  transform: rotate(-90deg);
}

.inspect-selected-stall-drawer .issd-aac-list,
.issd-aac-list {
  display: flex;
  flex-direction: column;
  background: #ffffff;
}

.inspect-selected-stall-drawer .issd-aac-row,
.issd-aac-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  border-top: 1px solid #f1f5f9;
  cursor: pointer;
  transition: background 0.15s ease;
}

.inspect-selected-stall-drawer .issd-aac-row:hover,
.issd-aac-row:hover {
  background: #f8fafc;
}

.inspect-selected-stall-drawer .issd-aac-row:hover .issd-aac-action-text,
.issd-aac-row:hover .issd-aac-action-text {
  color: #1d4ed8;
}

.inspect-selected-stall-drawer .issd-aac-row:hover .issd-aac-row-right svg,
.issd-aac-row:hover .issd-aac-row-right svg {
  transform: translateX(2px);
}

.inspect-selected-stall-drawer .issd-aac-row.is-danger:hover,
.issd-aac-row.is-danger:hover {
  background: #fef2f2;
}

.inspect-selected-stall-drawer .issd-aac-row.is-warning:hover,
.issd-aac-row.is-warning:hover {
  background: #fffbeb;
}

.inspect-selected-stall-drawer .issd-aac-row.is-info:hover,
.issd-aac-row.is-info:hover {
  background: #eff6ff;
}

.inspect-selected-stall-drawer .issd-aac-row-left,
.issd-aac-row-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex: 1;
}

.inspect-selected-stall-drawer .issd-aac-tag,
.issd-aac-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 7px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 700;
  white-space: nowrap;
  flex-shrink: 0;
}

.inspect-selected-stall-drawer .issd-aac-tag.is-danger,
.issd-aac-tag.is-danger {
  background: #fee2e2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.inspect-selected-stall-drawer .issd-aac-tag.is-warning,
.issd-aac-tag.is-warning {
  background: #fef3c7;
  color: #d97706;
  border: 1px solid #fde68a;
}

.inspect-selected-stall-drawer .issd-aac-tag.is-info,
.issd-aac-tag.is-info {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
}

.inspect-selected-stall-drawer .issd-aac-row-title,
.issd-aac-row-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.inspect-selected-stall-drawer .issd-aac-row-right,
.issd-aac-row-right {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  flex-shrink: 0;
  color: #2563eb;
  font-size: 0.78125rem;
  font-weight: 600;
}

.inspect-selected-stall-drawer .issd-aac-row-right svg,
.issd-aac-row-right svg {
  transition: transform 0.15s ease;
}

.inspect-selected-stall-drawer .issd-aac-action-text,
.issd-aac-action-text {
  white-space: nowrap;
}

/* ==========================================================================
   Full Detail Mode 2-Column Sidebar & Content Layout
   ========================================================================== */
.inspect-selected-stall-drawer.is-full-detail-mode {
  width: min(760px, calc(100vw - 36px)) !important;
  max-height: calc(100% - 76px) !important;
  z-index: 130 !important;
}

@keyframes issdModeFadeIn {
  from {
    opacity: 0;
    transform: translate3d(8px, 0, 0);
  }
  to {
    opacity: 1;
    transform: translate3d(0, 0, 0);
  }
}

.inspect-selected-stall-drawer .issd-full-layout,
.issd-full-layout {
  display: flex;
  flex: 1 1 auto;
  min-height: 0;
  overflow: hidden;
  background: #ffffff;
  animation: issdModeFadeIn 0.16s cubic-bezier(0.16, 1, 0.3, 1);
}

/* Left Sidebar Navigation */
.inspect-selected-stall-drawer .issd-full-sidebar,
.issd-full-sidebar {
  width: 152px;
  flex-shrink: 0;
  background: #f8fafc;
  border-right: 1px solid #eef2f6;
  display: flex;
  flex-direction: column;
  padding: 12px 8px;
  overflow-y: auto;
}

.inspect-selected-stall-drawer .issd-sidebar-nav,
.issd-sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.inspect-selected-stall-drawer .issd-sidebar-tab,
.issd-sidebar-tab {
  display: flex;
  align-items: center;
  gap: 9px;
  width: 100%;
  padding: 9px 12px;
  border-radius: 9px;
  border: 1px solid transparent;
  background: transparent;
  color: #475569;
  font-family: inherit;
  font-size: 0.84375rem;
  font-weight: 500;
  text-align: left;
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
  box-sizing: border-box;
}

.inspect-selected-stall-drawer .issd-sidebar-tab:hover,
.issd-sidebar-tab:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.inspect-selected-stall-drawer .issd-sidebar-tab.active,
.issd-sidebar-tab.active {
  background: #eff6ff;
  color: #2563eb;
  font-weight: 600;
  border-color: #dbeafe;
  box-shadow: 0 1px 3px rgba(37, 99, 235, 0.08);
}

.inspect-selected-stall-drawer .issd-sidebar-tab svg,
.issd-sidebar-tab svg {
  width: 15px;
  height: 15px;
  flex-shrink: 0;
  color: currentColor;
}

.inspect-selected-stall-drawer .issd-sidebar-tab-label,
.issd-sidebar-tab-label {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.inspect-selected-stall-drawer .issd-sidebar-tab-count,
.issd-sidebar-tab-count {
  margin-left: auto;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 1px 7px;
  border-radius: 999px;
  background: #e2e8f0;
  color: #475569;
  line-height: 1.25;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.inspect-selected-stall-drawer .issd-sidebar-tab.active .issd-sidebar-tab-count,
.issd-sidebar-tab.active .issd-sidebar-tab-count {
  background: #2563eb;
  color: #ffffff;
}

/* Right Main Area */
.inspect-selected-stall-drawer .issd-full-main,
.issd-full-main {
  flex: 1 1 auto;
  min-width: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #ffffff;
}

.inspect-selected-stall-drawer .issd-full-main .issd-action-alert-card,
.issd-full-main .issd-action-alert-card {
  margin: 12px 18px 4px;
}

.inspect-selected-stall-drawer .issd-full-main .issd-body-full,
.issd-full-main .issd-body-full {
  flex: 1 1 auto;
  min-height: 0;
  overflow-y: auto;
  padding: 12px 18px 16px;
}

.inspect-selected-stall-drawer .issd-tab-header-row.is-inline-controls,
.issd-tab-header-row.is-inline-controls {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.inspect-selected-stall-drawer .issd-contract-inline-filter,
.issd-contract-inline-filter {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.inspect-selected-stall-drawer .issd-inline-filter-label,
.issd-inline-filter-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
  white-space: nowrap;
}

.inspect-selected-stall-drawer .issd-inline-filter-select,
.issd-inline-filter-select {
  width: 195px;
}

.inspect-selected-stall-drawer .issd-full-main .issd-footer-full,
.issd-full-main .issd-footer-full {
  padding: 12px 18px 16px;
  background: #ffffff;
  border-top: 1px solid #f1f5f9;
  box-shadow: 0 -4px 12px rgba(15, 23, 42, 0.03);
}

/* Mobile Responsive Adjustments */
@media (max-width: 768px) {
  .inspect-selected-stall-drawer,
  .inspect-selected-stall-drawer.is-full-detail-mode {
    position: fixed !important;
    left: 0 !important;
    right: 0 !important;
    bottom: 0 !important;
    width: 100vw !important;
    max-width: 100vw !important;
    top: 58px !important;
    height: calc(100vh - 58px) !important;
    height: calc(100dvh - 58px) !important;
    max-height: calc(100vh - 58px) !important;
    max-height: calc(100dvh - 58px) !important;
    border-radius: 0 !important;
    border: none !important;
    box-shadow: none !important;
  }

  .inspect-selected-stall-drawer .issd-full-layout,
  .issd-full-layout {
    flex-direction: column;
  }

  .inspect-selected-stall-drawer .issd-full-sidebar,
  .issd-full-sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #eef2f6;
    padding: 8px 12px;
    overflow-x: auto;
    overflow-y: hidden;
  }

  .inspect-selected-stall-drawer .issd-sidebar-nav,
  .issd-sidebar-nav {
    flex-direction: row;
    white-space: nowrap;
    gap: 6px;
  }

  .inspect-selected-stall-drawer .issd-sidebar-tab,
  .issd-sidebar-tab {
    width: auto;
    flex-shrink: 0;
    min-height: 42px !important;
    height: 42px !important;
    padding: 0 16px !important;
    font-size: 0.875rem !important;
  }
}

/* ================================================================ */
/* รายละเอียดแผงแบบเต็ม: ธีมเดียวกันทุกแท็บ (การ์ด ตัวอักษร ปุ่ม สถานะ) */
/* ================================================================ */
.inspect-selected-stall-drawer .issd-full-main { background: #f8fafc; }
.inspect-selected-stall-drawer .issd-full-main .issd-body-full { padding: 20px 22px 24px; gap: 16px; display: flex; flex-direction: column; }
.inspect-selected-stall-drawer .issd-full-main .issd-footer-full { padding: 14px 22px 16px; background: #ffffff; }
.inspect-selected-stall-drawer .issd-full-sidebar { width: 168px; padding: 14px 10px; background: #ffffff; }
.inspect-selected-stall-drawer .issd-sidebar-tab { padding: 10px 12px; font-size: 0.9375rem; font-weight: 600; color: #475569; border-radius: 10px; }
.inspect-selected-stall-drawer .issd-sidebar-tab.active { background: #eff6ff; border-color: transparent; color: #2563eb; box-shadow: none; }

/* หัวข้อแท็บ + ตัวกรอง */
.inspect-selected-stall-drawer .issd-full-main .issd-tab-header-row { display: flex; flex-direction: row; align-items: center; justify-content: space-between; gap: 12px; margin: 0; text-align: left; }
.inspect-selected-stall-drawer .issd-full-main .issd-tab-title { margin: 0; color: #1e293b; font-size: 1.125rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main .app-segmented-control { align-self: flex-start; background: #ffffff; }
.inspect-selected-stall-drawer .issd-full-main .app-segmented-option { font-size: 0.875rem; min-height: 34px; padding: 5px 14px; }
.inspect-selected-stall-drawer .issd-full-main .app-segmented-option.is-active { background: #eff6ff; color: #2563eb; box-shadow: none; }
.inspect-selected-stall-drawer .issd-full-main .issd-scroll-list { display: flex; flex-direction: column; gap: 12px; }

/* การ์ดรายการทุกชนิดหน้าตาเดียวกัน */
.inspect-selected-stall-drawer .issd-full-main :is(.issd-booking-card, .issd-contract-card, .issd-subtenant-card, .issd-payment-card, .issd-maintenance-card) {
  display: flex; flex-direction: column; gap: 14px; padding: 16px 18px;
  background: #ffffff; border: 1px solid #eef2f6; border-radius: 14px; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
}
.inspect-selected-stall-drawer .issd-full-main .issd-contract-card.is-current { background: #ffffff; border-color: #eef2f6; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04); }
.inspect-selected-stall-drawer .issd-full-main :is(.issd-card-top-row, .issd-booking-card-header) { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.inspect-selected-stall-drawer .issd-full-main .issd-card-title { color: #1e293b; font-size: 1rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main :is(.issd-contract-num, .issd-booking-ref-code, .issd-pay-period strong) { color: #1e293b; font-size: 1rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main :is(.detail-label, .issd-booking-role-label, .issd-maintenance-meta-grid span, .issd-pay-method) { color: #64748b; font-size: 0.8125rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-full-main :is(.detail-value, .issd-maintenance-meta-grid strong) { color: #1e293b; font-size: 0.9375rem; font-weight: 600; }
.inspect-selected-stall-drawer .issd-full-main .is-rent { color: #10b981; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main .issd-contract-details-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px 16px; }
.inspect-selected-stall-drawer .issd-full-main .issd-contract-detail-item { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.inspect-selected-stall-drawer .issd-full-main .issd-contract-detail-item.is-wide { grid-column: 1 / -1; }
.inspect-selected-stall-drawer .issd-full-main .issd-maintenance-meta-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 10px 16px; }
.inspect-selected-stall-drawer .issd-full-main .issd-maintenance-meta-grid > div { display: flex; flex-direction: column; gap: 2px; }
.inspect-selected-stall-drawer .issd-full-main .issd-card-actions-row { display: flex; justify-content: flex-end; gap: 8px; padding-top: 12px; border-top: 1px solid #f1f5f9; }
.inspect-selected-stall-drawer .issd-full-main .issd-card-meta-row,
.inspect-selected-stall-drawer .issd-full-main .issd-pay-amount-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; color: #475569; font-size: 0.9375rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-full-main .issd-pay-amount-row strong { font-size: 1.125rem; font-weight: 800; }

/* ป้ายสถานะ (soft pill ตามมาตรฐานระบบ) */
.inspect-selected-stall-drawer .issd-full-main .issd-status-pill { padding: 4px 12px; border-radius: 8px; font-size: 0.8125rem; font-weight: 600; }

/* ผู้เช่า / ผู้จอง */
.inspect-selected-stall-drawer .issd-full-main .issd-booking-applicant-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap; }
.inspect-selected-stall-drawer .issd-full-main .issd-phone-link { font-size: 0.875rem; font-weight: 600; }
.inspect-selected-stall-drawer .issd-full-main .issd-booking-footer-bar { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding-top: 12px; border-top: 1px solid #f1f5f9; color: #475569; font-size: 0.9375rem; }
.inspect-selected-stall-drawer .issd-full-main .issd-booking-direct-btn { display: inline-flex; align-items: center; gap: 6px; height: 34px; padding: 0 12px; border: 1px solid #93c5fd; border-radius: 7px; background: #fff; color: #2563eb; font-size: 0.875rem; font-weight: 600; text-decoration: none; }
.inspect-selected-stall-drawer .issd-full-main .issd-booking-direct-btn:hover { background: #eff6ff; }

/* สรุปการเงิน */
.inspect-selected-stall-drawer .issd-full-main .issd-finance-summary { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; }
.inspect-selected-stall-drawer .issd-full-main .issd-finance-summary > div { display: flex; flex-direction: column; gap: 2px; padding: 12px 14px; background: #fff; border: 1px solid #eef2f6; border-radius: 12px; }
.inspect-selected-stall-drawer .issd-full-main .issd-finance-summary span { color: #64748b; font-size: 0.8125rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-full-main .issd-finance-summary strong { color: #1e293b; font-size: 1.125rem; font-weight: 800; }

/* ข้อมูลแผง: การ์ดตัวเลข + รายการรายละเอียดเป็นแถว */
.inspect-selected-stall-drawer .issd-full-main .issd-kpi-card { background: #fff; border: 1px solid #eef2f6; border-radius: 14px; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-section { display: flex; flex-direction: column; gap: 10px; }
.inspect-selected-stall-drawer .issd-full-main .issd-section-title { color: #1e293b; font-size: 1.125rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-grid { display: flex; flex-direction: column; background: #fff; border: 1px solid #eef2f6; border-radius: 14px; overflow: hidden; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell,
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell.full { display: flex; align-items: baseline; justify-content: space-between; gap: 16px; width: auto; padding: 14px 18px; border: 0; border-bottom: 1px solid #f1f5f9; border-radius: 0; background: transparent; box-shadow: none; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell:last-child { border-bottom: 0; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell dt { color: #64748b; font-size: 0.9375rem; font-weight: 500; white-space: nowrap; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell dd { margin: 0; color: #1e293b; font-size: 0.9375rem; font-weight: 600; text-align: right; }
.inspect-selected-stall-drawer .issd-full-main .issd-photo-banner { border-radius: 14px; overflow: hidden; }

/* ประวัติ (timeline) */
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-panel { background: #fff; border: 1px solid #eef2f6; border-radius: 14px; padding: 8px 18px; }
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-item { padding: 12px 0; }
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-copy strong { color: #1e293b; font-size: 0.9375rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-copy time { color: #94a3b8; font-size: 0.8125rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-copy p { margin: 4px 0 0; color: #475569; font-size: 0.9375rem; line-height: 1.55; }

/* สถานะว่าง */
.inspect-selected-stall-drawer .issd-full-main .issd-empty-box { display: flex; flex-direction: column; align-items: center; gap: 10px; padding: 36px 20px; background: #fff; border: 1px dashed #dbe3ee; border-radius: 14px; color: #94a3b8; text-align: center; }
.inspect-selected-stall-drawer .issd-full-main .issd-empty-box p { margin: 0; color: #64748b; font-size: 0.9375rem; font-weight: 500; }

/* ตัวเลข/รหัสในรายการต้องดำเนินการ แสดงเป็นป้ายแยก ไม่ปนอยู่ในประโยค */
.issd-aac-meta { flex: 0 0 auto; padding: 2px 10px; border: 1px solid rgba(15, 23, 42, 0.1); border-radius: 8px; background: rgba(255, 255, 255, 0.75); color: #475569; font-size: 0.8125rem; font-weight: 600; white-space: nowrap; }

/* ===== Drawer readability pass: larger type, roomier spacing, no redundant blocks ===== */
.inspect-selected-stall-drawer .issd-title-group strong { font-size: 1.0625rem; }
.inspect-selected-stall-drawer .issd-project-sub { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-header { padding: 16px 20px 14px; }
.inspect-selected-stall-drawer .issd-body { gap: 16px; padding: 18px 20px; }
.inspect-selected-stall-drawer .issd-body-full { gap: 16px; padding: 18px 22px 22px; }

.inspect-selected-stall-drawer .issd-available-line { display: flex; align-items: center; gap: 10px; color: #15803d; }
.inspect-selected-stall-drawer .issd-available-line strong { font-size: 1.0625rem; font-weight: 700; }

.inspect-selected-stall-drawer .issd-kpi-grid { gap: 14px; }
.inspect-selected-stall-drawer .issd-kpi-grid.is-three { grid-template-columns: repeat(3, 1fr); }
.inspect-selected-stall-drawer .issd-kpi-card { padding: 16px 18px; gap: 4px; border-radius: 14px; }
.inspect-selected-stall-drawer .issd-kpi-card span,
.inspect-selected-stall-drawer .issd-kpi-card em { font-size: 0.875rem; }
.inspect-selected-stall-drawer .issd-kpi-card strong { font-size: 1.25rem; }
.inspect-selected-stall-drawer .issd-section-title { font-size: 1.0625rem; }
.inspect-selected-stall-drawer .issd-fact-grid { gap: 12px; }
.inspect-selected-stall-drawer .issd-fact-cell { padding: 13px 16px; border-radius: 12px; }
.inspect-selected-stall-drawer .issd-fact-cell dt { font-size: 0.875rem; }
.inspect-selected-stall-drawer .issd-fact-cell dd { font-size: 0.9375rem; margin-top: 3px; }
.inspect-selected-stall-drawer .issd-occupancy-box { padding: 14px 18px; }
.inspect-selected-stall-drawer .issd-occupancy-head span,
.inspect-selected-stall-drawer .issd-occupancy-head small { font-size: 0.875rem; }
.inspect-selected-stall-drawer .issd-occupancy-box strong { font-size: 1rem; }
.inspect-selected-stall-drawer .issd-occupancy-box p { font-size: 0.875rem; }

.inspect-selected-stall-drawer .issd-detail-row { padding: 14px 18px; font-size: 0.9375rem; }
.inspect-selected-stall-drawer .issd-detail-row dt,
.inspect-selected-stall-drawer .issd-detail-row dd { font-size: 0.9375rem; }

.inspect-selected-stall-drawer .issd-card-section { padding: 18px 20px; gap: 14px; }
.inspect-selected-stall-drawer .issd-person-role { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-person-name { font-size: 1rem; }
.inspect-selected-stall-drawer .issd-subdoc-row .subdoc-k,
.inspect-selected-stall-drawer .issd-subdoc-row .subdoc-v { font-size: 0.9375rem; }
.inspect-selected-stall-drawer .issd-phone-link { font-size: 0.875rem; }

.inspect-selected-stall-drawer .issd-mini-ops-grid { gap: 12px; }
.inspect-selected-stall-drawer .issd-mini-ops-card { padding: 14px 16px; gap: 12px; border-radius: 13px; }
.inspect-selected-stall-drawer .issd-mini-ops-card .ops-icon { width: 38px; height: 38px; }
.inspect-selected-stall-drawer .issd-mini-ops-card .ops-label { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-mini-ops-card .ops-info strong { font-size: 0.9375rem; }

.inspect-selected-stall-drawer .issd-mini-alert-card .issd-mini-alert-banner { padding: 14px 16px; }
.inspect-selected-stall-drawer .issd-mini-alert-head { font-size: 0.9375rem; }
.inspect-selected-stall-drawer .issd-mini-alert-link { font-size: 0.875rem; }
.inspect-selected-stall-drawer .issd-mini-alert-list { padding: 8px 12px 12px; gap: 8px; }
.inspect-selected-stall-drawer .issd-mini-alert-row { padding: 11px 12px; }
.inspect-selected-stall-drawer .issd-mini-alert-text strong { font-size: 0.9375rem; white-space: normal; }
.inspect-selected-stall-drawer .issd-mini-alert-desc { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-mini-alert-action-pill { font-size: 0.8125rem; padding: 4px 10px; }

.inspect-selected-stall-drawer .issd-aac-title-row strong { font-size: 1.0625rem; }
.inspect-selected-stall-drawer .issd-aac-count-badge { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-aac-row { padding: 12px 14px; }
.inspect-selected-stall-drawer .issd-aac-tag { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-aac-row-title { font-size: 0.9375rem; }
.inspect-selected-stall-drawer .issd-aac-action-text { font-size: 0.875rem; }

.inspect-selected-stall-drawer .issd-full-sidebar { width: 172px; padding: 12px 8px; }
.inspect-selected-stall-drawer .issd-sidebar-nav { gap: 4px; }
.inspect-selected-stall-drawer .issd-sidebar-tab { font-size: 0.875rem; padding: 10px 12px; gap: 9px; border-radius: 10px; }
.inspect-selected-stall-drawer .issd-sidebar-tab svg { width: 16px; height: 16px; }
.inspect-selected-stall-drawer .issd-sidebar-tab-count { font-size: 0.75rem; padding: 1px 8px; }
.inspect-selected-stall-drawer .issd-tab-title { font-size: 1.0625rem; }
.inspect-selected-stall-drawer .issd-tab-subtitle { font-size: 0.8125rem; }
@media (max-width: 768px) {
  .inspect-selected-stall-drawer .issd-full-sidebar { width: 100%; }
  .inspect-selected-stall-drawer .issd-kpi-grid.is-three { grid-template-columns: repeat(2, 1fr); }
}

/* ===== Fit pass: no clipped text, wider mini drawer, readable tab internals ===== */
.inspect-selected-stall-drawer:not(.is-full-detail-mode) { width: min(440px, calc(100vw - 28px)); }
.inspect-selected-stall-drawer .issd-person-name,
.inspect-selected-stall-drawer .issd-subdoc-row .subdoc-v,
.inspect-selected-stall-drawer .issd-subdoc-row .subdoc-v.is-code,
.inspect-selected-stall-drawer .issd-mini-ops-card .ops-info strong,
.inspect-selected-stall-drawer .issd-mini-ops-card .ops-label,
.inspect-selected-stall-drawer .issd-detail-row dd {
  white-space: normal; overflow: visible; text-overflow: clip; word-break: break-word; overflow-wrap: anywhere;
}
.inspect-selected-stall-drawer .issd-subdoc-row { flex-wrap: wrap; gap: 4px 12px; align-items: baseline; }
.inspect-selected-stall-drawer .issd-subdoc-row .subdoc-v { text-align: right; margin-left: auto; }
.inspect-selected-stall-drawer .issd-person-row { flex-wrap: wrap; row-gap: 10px; }
.inspect-selected-stall-drawer .issd-person-info { flex: 1 1 160px; min-width: 0; }
.inspect-selected-stall-drawer .issd-detail-row dd { max-width: 70%; }

/* Expandable alert list: stack vertically so nothing is cut off */
.inspect-selected-stall-drawer .issd-mini-alert-banner { flex-wrap: wrap; row-gap: 4px; }
.inspect-selected-stall-drawer .issd-mini-alert-head strong { white-space: normal; line-height: 1.35; }
.inspect-selected-stall-drawer .issd-mini-alert-row { flex-direction: column; align-items: stretch; gap: 8px; }
.inspect-selected-stall-drawer .issd-mini-alert-row-left { flex-direction: column; align-items: flex-start; gap: 6px; }
.inspect-selected-stall-drawer .issd-mini-alert-text strong,
.inspect-selected-stall-drawer .issd-mini-alert-desc { white-space: normal; overflow: visible; text-overflow: clip; }
.inspect-selected-stall-drawer .issd-mini-alert-action-pill { align-self: flex-end; }
.inspect-selected-stall-drawer .issd-aac-row-title,
.inspect-selected-stall-drawer .issd-aac-title-row strong { white-space: normal; overflow: visible; text-overflow: clip; }

/* Tab internals (contracts / subtenants / payments / maintenance / bookings / history) */
.inspect-selected-stall-drawer .issd-body-full { font-size: 0.9375rem; }
.inspect-selected-stall-drawer .issd-card-title { font-size: 1rem; }
.inspect-selected-stall-drawer .issd-contract-num,
.inspect-selected-stall-drawer .issd-booking-ref-code { font-size: 0.9375rem; }
.inspect-selected-stall-drawer .issd-contract-detail-item,
.inspect-selected-stall-drawer .issd-detail-item,
.inspect-selected-stall-drawer .issd-subtenant-note,
.inspect-selected-stall-drawer .issd-maintenance-desc,
.inspect-selected-stall-drawer .issd-booking-info-item,
.inspect-selected-stall-drawer .issd-pay-period,
.inspect-selected-stall-drawer .issd-pay-method,
.inspect-selected-stall-drawer .issd-timeline-copy p,
.inspect-selected-stall-drawer .issd-timeline-copy time { font-size: 0.875rem; line-height: 1.5; }
.inspect-selected-stall-drawer .issd-timeline-copy strong { font-size: 0.9375rem; }
.inspect-selected-stall-drawer .issd-contract-detail-item *,
.inspect-selected-stall-drawer .issd-detail-item *,
.inspect-selected-stall-drawer .issd-maintenance-meta-grid *,
.inspect-selected-stall-drawer .issd-booking-info-item * { font-size: inherit; }
.inspect-selected-stall-drawer .issd-pay-amount-row strong { font-size: 1.0625rem; }
.inspect-selected-stall-drawer .issd-status-pill,
.inspect-selected-stall-drawer .issd-card-type-tag,
.inspect-selected-stall-drawer .issd-source-pill,
.inspect-selected-stall-drawer .issd-maintenance-type-pill,
.inspect-selected-stall-drawer .issd-timeline-cat-pill,
.inspect-selected-stall-drawer .issd-booking-role-label,
.inspect-selected-stall-drawer .issd-current-label { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-mini-action-btn,
.inspect-selected-stall-drawer .issd-booking-direct-btn,
.inspect-selected-stall-drawer .issd-inline-filter-label { font-size: 0.875rem; }
.inspect-selected-stall-drawer .issd-contract-card,
.inspect-selected-stall-drawer .issd-subtenant-card,
.inspect-selected-stall-drawer .issd-payment-card,
.inspect-selected-stall-drawer .issd-maintenance-card,
.inspect-selected-stall-drawer .issd-booking-card { padding: 16px 18px; gap: 12px; }
.inspect-selected-stall-drawer .issd-contract-card *,
.inspect-selected-stall-drawer .issd-maintenance-card *,
.inspect-selected-stall-drawer .issd-booking-card *,
.inspect-selected-stall-drawer .issd-subtenant-card * { overflow-wrap: anywhere; }

/* ===== Unified tab theme: one card language for contracts / subtenants / payments / maintenance / bookings ===== */
.inspect-selected-stall-drawer .issd-body-full .issd-tab-header-row { margin: 0 0 4px; }
.inspect-selected-stall-drawer .issd-body-full .issd-tab-title { font-size: 1.0625rem; font-weight: 700; color: #0f172a; }

.inspect-selected-stall-drawer .issd-contract-list,
.inspect-selected-stall-drawer .issd-subtenant-list,
.inspect-selected-stall-drawer .issd-payment-list,
.inspect-selected-stall-drawer .issd-maintenance-list,
.inspect-selected-stall-drawer .issd-booking-list { display: flex; flex-direction: column; gap: 12px; }

.inspect-selected-stall-drawer .issd-contract-card,
.inspect-selected-stall-drawer .issd-subtenant-card,
.inspect-selected-stall-drawer .issd-payment-card,
.inspect-selected-stall-drawer .issd-maintenance-card,
.inspect-selected-stall-drawer .issd-booking-card {
  background: #ffffff; border: 1px solid #e2e8f0; border-radius: 14px;
  padding: 16px 18px; display: flex; flex-direction: column; gap: 12px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}
.inspect-selected-stall-drawer .issd-contract-card.is-current { border-color: #bbf7d0; background: #f7fef9; }

.inspect-selected-stall-drawer .issd-card-top-row,
.inspect-selected-stall-drawer .issd-booking-card-header {
  display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap;
}
.inspect-selected-stall-drawer .issd-card-title { font-size: 1rem; font-weight: 700; color: #0f172a; }
.inspect-selected-stall-drawer .issd-contract-num { font-size: 0.9375rem; font-weight: 700; color: #1e293b; }

/* label above value grid */
.inspect-selected-stall-drawer .issd-contract-details-grid {
  display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px 16px;
  padding-top: 12px; border-top: 1px solid #f1f5f9;
}
.inspect-selected-stall-drawer .issd-contract-detail-item { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.inspect-selected-stall-drawer .issd-contract-detail-item.is-wide { grid-column: 1 / -1; }
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-label { font-size: 0.8125rem; color: #64748b; font-weight: 500; }
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-value { font-size: 0.9375rem; color: #1e293b; font-weight: 600; }
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-value.is-rent { color: #15803d; }
.inspect-selected-stall-drawer .detail-sub { font-style: normal; font-weight: 500; color: #2563eb; }

.inspect-selected-stall-drawer .issd-subtenant-details,
.inspect-selected-stall-drawer .issd-card-meta-row {
  display: flex; flex-wrap: wrap; align-items: center; gap: 6px 16px;
  font-size: 0.875rem; color: #475569;
}
.inspect-selected-stall-drawer .issd-card-meta-row { justify-content: space-between; padding-top: 10px; border-top: 1px solid #f1f5f9; }
.inspect-selected-stall-drawer .issd-card-meta-row .is-rent,
.inspect-selected-stall-drawer .issd-pay-amount-row .is-rent { color: #15803d; font-weight: 700; }
.inspect-selected-stall-drawer .issd-pay-period { display: flex; flex-direction: column; gap: 2px; }
.inspect-selected-stall-drawer .issd-pay-period strong { font-size: 1rem; color: #0f172a; }
.inspect-selected-stall-drawer .issd-pay-amount-row { display: flex; align-items: baseline; justify-content: space-between; gap: 10px; }
.inspect-selected-stall-drawer .issd-pay-amount-row strong { font-size: 1.125rem; }
.inspect-selected-stall-drawer .issd-maintenance-meta-grid { display: flex; flex-wrap: wrap; gap: 6px 18px; font-size: 0.875rem; color: #64748b; }
.inspect-selected-stall-drawer .issd-maintenance-meta-grid strong { color: #1e293b; font-weight: 600; }
.inspect-selected-stall-drawer .issd-finance-summary { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
.inspect-selected-stall-drawer .issd-finance-summary > div { padding: 12px 14px; border: 1px solid #e2e8f0; border-radius: 12px; background: #fff; display: flex; flex-direction: column; gap: 2px; }
.inspect-selected-stall-drawer .issd-finance-summary span { font-size: 0.8125rem; color: #64748b; }
.inspect-selected-stall-drawer .issd-finance-summary strong { font-size: 1.0625rem; color: #0f172a; }
.inspect-selected-stall-drawer .issd-finance-summary strong.is-success { color: #15803d; }
.inspect-selected-stall-drawer .issd-finance-summary strong.is-alert { color: #dc2626; }

/* one shared action row + button style */
.inspect-selected-stall-drawer .issd-card-actions-row { display: flex; gap: 8px; justify-content: flex-end; padding-top: 4px; }
.inspect-selected-stall-drawer .issd-mini-action-btn {
  padding: 7px 14px; border-radius: 8px; border: 1px solid #e2e8f0; background: #fff;
  color: #334155; font-family: inherit; font-weight: 600; cursor: pointer; transition: all .15s ease;
}
.inspect-selected-stall-drawer .issd-mini-action-btn:hover { background: #f8fafc; }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-primary { border-color: #93c5fd; color: #2563eb; }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-primary:hover { background: #eff6ff; }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-success { border-color: #86efac; color: #15803d; }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-success:hover { background: #f0fdf4; }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-danger { border-color: #fca5a5; color: #ef4444; }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-danger:hover { background: #fef2f2; }

/* sub-filters share the same segmented look */
.inspect-selected-stall-drawer .issd-source-filter { display: inline-flex; gap: 4px; padding: 4px; background: #f1f5f9; border-radius: 10px; align-self: flex-start; }
.inspect-selected-stall-drawer .issd-source-filter button {
  padding: 7px 14px; border: none; border-radius: 8px; background: transparent; color: #64748b;
  font-family: inherit; font-size: 0.875rem; font-weight: 600; cursor: pointer;
}
.inspect-selected-stall-drawer .issd-source-filter button.active { background: #fff; color: #2563eb; box-shadow: 0 1px 3px rgba(15,23,42,.08); }

.inspect-selected-stall-drawer .issd-empty-box { padding: 36px 16px; border: 1px dashed #cbd5e1; border-radius: 14px; background: #f8fafc; display: flex; flex-direction: column; align-items: center; gap: 10px; color: #64748b; }
.inspect-selected-stall-drawer .issd-empty-box p { margin: 0; font-size: 0.9375rem; }

/* bookings */
.inspect-selected-stall-drawer .issd-booking-applicant-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap; }
.inspect-selected-stall-drawer .issd-booking-footer-bar { display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap; padding-top: 10px; border-top: 1px solid #f1f5f9; font-size: 0.875rem; color: #475569; }

/* Prevent flex children of the drawer body from collapsing into thin lines */
.inspect-selected-stall-drawer .issd-body > *,
.inspect-selected-stall-drawer .issd-body-full > * { flex-shrink: 0; }
.inspect-selected-stall-drawer .issd-mini-alert-card { flex-shrink: 0; min-height: min-content; }

/* Compact contract card: 3 columns + period row, no nested grey box */
.inspect-selected-stall-drawer .issd-contract-card { padding: 14px 16px; gap: 10px; }
.inspect-selected-stall-drawer .issd-contract-details-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px 14px;
  background: transparent; padding: 10px 0 0; border: none; border-top: 1px solid #f1f5f9; border-radius: 0;
}
.inspect-selected-stall-drawer .issd-contract-card .issd-card-actions-row { padding: 0; margin-top: -2px; }
.inspect-selected-stall-drawer .issd-subtenant-card,
.inspect-selected-stall-drawer .issd-payment-card,
.inspect-selected-stall-drawer .issd-maintenance-card,
.inspect-selected-stall-drawer .issd-booking-card { padding: 14px 16px; gap: 10px; }
.inspect-selected-stall-drawer .issd-contract-list,
.inspect-selected-stall-drawer .issd-subtenant-list,
.inspect-selected-stall-drawer .issd-payment-list,
.inspect-selected-stall-drawer .issd-maintenance-list,
.inspect-selected-stall-drawer .issd-booking-list { gap: 10px; }
.inspect-selected-stall-drawer .issd-body-full { gap: 12px; padding: 14px 20px 16px; }
@media (max-width: 768px) {
  .inspect-selected-stall-drawer .issd-contract-details-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

.inspect-selected-stall-drawer .issd-sidebar-tab-count.is-action { background: #fef3c7; color: #b45309; }
.inspect-selected-stall-drawer .issd-sidebar-tab-count.is-action.is-danger { background: #fee2e2; color: #dc2626; }
.inspect-selected-stall-drawer .issd-sidebar-tab.active .issd-sidebar-tab-count.is-action { background: #d97706; color: #fff; }
.inspect-selected-stall-drawer .issd-sidebar-tab.active .issd-sidebar-tab-count.is-action.is-danger { background: #dc2626; color: #fff; }

/* Mini alert rows: chip + title on one line, divided rows (same look as full-drawer list) */
.inspect-selected-stall-drawer .issd-mini-alert-list { padding: 0; gap: 0; border-top: 1px solid #fde68a; background: #fff; }
.inspect-selected-stall-drawer .issd-mini-alert-row {
  flex-direction: row; align-items: center; justify-content: flex-start; gap: 10px;
  padding: 11px 14px; background: #fff; border: none; border-radius: 0; box-shadow: none; transform: none;
  border-top: 1px solid #f1f5f9; cursor: pointer;
}
.inspect-selected-stall-drawer .issd-mini-alert-row:first-child { border-top: none; }
.inspect-selected-stall-drawer .issd-mini-alert-row.is-danger { background: #fff1f1; }
.inspect-selected-stall-drawer .issd-mini-alert-row:hover { background: #fffbeb; transform: none; box-shadow: none; }
.inspect-selected-stall-drawer .issd-mini-alert-row.is-danger:hover { background: #ffe4e4; }
.inspect-selected-stall-drawer .issd-mini-alert-title { font-size: 0.9375rem; font-weight: 600; color: #1e293b; line-height: 1.35; min-width: 0; overflow-wrap: anywhere; }
.inspect-selected-stall-drawer .issd-mini-alert-row .issd-aac-tag { flex-shrink: 0; }

/* Mini banner header = same look as full drawer (icon box + title + count pill + chevron) */
.inspect-selected-stall-drawer .issd-mini-alert-head { gap: 10px; flex-wrap: wrap; }
.inspect-selected-stall-drawer .issd-mini-alert-head strong { font-size: 1rem; color: #92400e; }
.inspect-selected-stall-drawer .issd-mini-alert-card .issd-aac-toggle-btn { pointer-events: none; }

/* Photo editor in GIS edit drawer */
.inspect-selected-stall-drawer .issd-photo-editor { display: flex; flex-direction: column; gap: 8px; }
.inspect-selected-stall-drawer .issd-photo-editor-head { display: flex; align-items: center; justify-content: space-between; font-size: 0.875rem; font-weight: 600; color: #1e293b; }
.inspect-selected-stall-drawer .issd-photo-editor-head small { font-size: 0.8125rem; color: #64748b; font-weight: 500; }
.inspect-selected-stall-drawer .issd-photo-editor-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.inspect-selected-stall-drawer .issd-photo-thumb { position: relative; aspect-ratio: 4 / 3; border-radius: 10px; overflow: hidden; border: 1px solid #e2e8f0; background: #f1f5f9; }
.inspect-selected-stall-drawer .issd-photo-thumb img { width: 100%; height: 100%; object-fit: cover; display: block; }
.inspect-selected-stall-drawer .issd-photo-remove { position: absolute; top: 5px; right: 5px; width: 24px; height: 24px; border-radius: 999px; border: none; background: rgba(15,23,42,.72); color: #fff; display: grid; place-items: center; cursor: pointer; }
.inspect-selected-stall-drawer .issd-photo-remove:hover { background: #dc2626; }
.inspect-selected-stall-drawer .issd-photo-add { aspect-ratio: 4 / 3; border: 1.5px dashed #93c5fd; border-radius: 10px; background: #f8fbff; color: #2563eb; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 4px; font-size: 0.8125rem; font-weight: 600; cursor: pointer; transition: background .15s ease; }
.inspect-selected-stall-drawer .issd-photo-add:hover { background: #eff6ff; }
</style>
