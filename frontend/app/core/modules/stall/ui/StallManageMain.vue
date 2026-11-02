<template>
  <div>
    <AppHeader />
    <main
      class="content-body manage-page"
      :class="{ 'is-map-page': isMapOverviewMode, 'has-left-panel': isMapOverviewMode && !!leftPanelWidth, 'has-right-panel': isMapOverviewMode && !!rightSidePanelWidth }"
      :style="isMapOverviewMode ? { ...(leftPanelWidth ? { '--lp-w': leftPanelWidth } : {}), ...(rightSidePanelWidth ? { '--side-w': `${rightSidePanelWidth}px` } : {}) } : undefined"
    >
      <!-- สวิตช์มุมมองแผนที่ (ข้อมูลแผง / การชำระเงิน / การจดมิเตอร์) แบบ Capsule Pill ลอยมุมซ้ายบน -->
      <div v-if="isMapOverviewMode" ref="viewDropdownWrapRef" class="map-mode-dropdown-wrap" :class="{ 'is-compact': !!rightSidePanelWidth }">
        <button
          type="button"
          class="map-mode-trigger-btn"
          :class="{ 'is-open': isViewDropdownOpen }"
          :title="currentModeOption.label"
          :aria-label="`มุมมองแผนที่: ${currentModeOption.label}`"
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
          title="ดูข้อมูลมิเตอร์ในแผนที่"
          aria-label="ดูข้อมูลมิเตอร์ในแผนที่"
          @click="meterSummaryOpen = true"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <rect x="3" y="3" width="18" height="18" rx="2"></rect>
            <path d="M3 9h18"></path>
            <path d="M3 15h18"></path>
            <path d="M12 3v18"></path>
          </svg>
        </button>

        <div v-if="meterSummaryOpen" class="map-cost-float meter-summary-float">
          <StallMeterSummaryPanel
            :projects="meterSummaryProjects"
            v-model:month="costMonth"
            :current-month="meterCurrentMonth"
            :stalls="meterSummaryData.stalls"
            :averages="meterSummaryData.averages"
            :skip-reasons="meterSkipReasonCounts"
            @open-stall="openMeterStall"
            @open-meter-settings="goMeterSettings"
            @close="meterSummaryOpen = false"
            @open-project="openMeterProject"
          />
        </div>


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

        <StallMeterRecordCard
          v-if="pageMode === 'meters' && meterCardCode"
          :code="meterCardCode"
          :label="meterCardLabel"
          :project-name="selectedProject ? cleanProjectName(selectedProject.name) : ''"
          :read-only="isPastPeriod"
          :pending-left="meterPendingLeft"
          :has-next="meterHasNext"
          :average-electric="meterAverages.electric"
          :average-water="meterAverages.water"
          @back="closeMeterCard"
          :can-go-back="meterCanGoBack"
          @prev="openPrevMeterStall"
          @close="closeMeterCardToOverview"
          @saved="handleMeterRecorded($event)"
          @skipped="handleMeterSkipped"
          :notice="meterNotice"
        />

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
                  <td colspan="9" class="dash-empty-row"><AppEmptyState variant="compact" icon="search" title="ไม่พบข้อมูลโซน" description="ลองเปลี่ยนตัวกรองเพื่อดูข้อมูลอื่น" /></td>
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
            <span class="map-summary-icon" style="background: var(--bg-eff6ff); color: var(--tx-2563eb);">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="5" width="20" height="14" rx="2"></rect><line x1="2" y1="10" x2="22" y2="10"></line></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ยอดเรียกเก็บประจำเดือน</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem;">฿284,500</strong>
                <em class="map-summary-unit-pill" style="background: var(--bg-eff6ff); color: var(--tx-2563eb);">42 บิล</em>
              </span>
              <span class="map-summary-meta">อัตราเฉลี่ย ฿6,773 / แผง</span>
            </span>
          </div>

          <div class="map-summary-card summary-available">
            <span class="map-summary-icon" style="background: var(--bg-ecfdf5); color: var(--tx-059669);">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">จัดเก็บสำเร็จแล้ว</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: var(--tx-059669);">฿192,000</strong>
                <em class="map-summary-unit-pill" style="background: var(--bg-ecfdf5); color: var(--tx-059669);">67.5%</em>
              </span>
              <span class="map-summary-meta">ชำระแล้ว 28 จาก 42 รายการ</span>
            </span>
          </div>

          <div class="map-summary-card summary-reserved">
            <span class="map-summary-icon" style="background: var(--bg-fffbeb); color: var(--tx-d97706);">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">รอชำระในกำหนด</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: var(--tx-d97706);">฿58,500</strong>
                <em class="map-summary-unit-pill" style="background: var(--bg-fffbeb); color: var(--tx-d97706);">9 บิล</em>
              </span>
              <span class="map-summary-meta">ครบกำหนดสิ้นเดือนนี้</span>
            </span>
          </div>

          <div class="map-summary-card summary-action-required">
            <span class="map-summary-icon" style="background: var(--bg-fef2f2); color: var(--tx-dc2626);">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ค้างชำระเกินกำหนด</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: var(--tx-dc2626);">฿34,000</strong>
                <em class="map-summary-unit-pill" style="background: var(--bg-fee2e2); color: var(--tx-dc2626);">5 บิล</em>
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
                  <td style="text-align: right; font-weight: 700; color: var(--tx-dc2626);">
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
            <span class="map-summary-icon" style="background: var(--bg-eff6ff); color: var(--tx-2563eb);">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m12 14 4-4"></path><path d="M3.34 19a10 10 0 1 1 17.32 0"></path></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ความคืบหน้าการจด</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: var(--tx-2563eb);">{{ meterSummary.completionRate }}%</strong>
                <em class="map-summary-unit-pill" style="background: var(--bg-eff6ff); color: var(--tx-2563eb);">{{ meterSummary.recordedCount }}/{{ meterSummary.totalCount }} แผง</em>
              </span>
              <span class="map-summary-meta">เหลือ {{ meterSummary.pendingCount }} แผงที่ต้องจด</span>
            </span>
          </div>

          <div class="map-summary-card summary-available">
            <span class="map-summary-icon" style="background: var(--bg-ecfdf5); color: var(--tx-059669);">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2.69l5.66 5.66a8 8 0 1 1-11.31 0z"></path></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ปริมาณการใช้น้ำรวม</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: var(--tx-059669);">{{ meterUsageSummary.waterTotal.toLocaleString('th-TH') }}</strong>
                <em class="map-summary-unit-pill" style="background: var(--bg-ecfdf5); color: var(--tx-059669);">ลบ.ม.</em>
              </span>
              <span class="map-summary-meta">เฉลี่ย {{ meterUsageSummary.waterAverage.toLocaleString('th-TH') }} ลบ.ม. / แผงที่จดแล้ว</span>
            </span>
          </div>

          <div class="map-summary-card summary-reserved">
            <span class="map-summary-icon" style="background: var(--bg-fffbeb); color: var(--tx-d97706);">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ปริมาณการใช้ไฟฟ้ารวม</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: var(--tx-d97706);">{{ meterUsageSummary.electricTotal.toLocaleString('th-TH') }}</strong>
                <em class="map-summary-unit-pill" style="background: var(--bg-fffbeb); color: var(--tx-d97706);">หน่วย (kWh)</em>
              </span>
              <span class="map-summary-meta">เฉลี่ย {{ meterUsageSummary.electricAverage.toLocaleString('th-TH') }} หน่วย / แผงที่จดแล้ว</span>
            </span>
          </div>

          <div class="map-summary-card summary-action-required">
            <span class="map-summary-icon" style="background: var(--bg-f5f3ff); color: var(--tx-7c3aed);">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">ตรวจพบค่าผิดปกติ</small>
              <span class="map-summary-value-row">
                <strong class="map-summary-count" style="font-size: 1.25rem; color: var(--tx-7c3aed);">{{ meterSummary.anomalyCount }}</strong>
                <em class="map-summary-unit-pill" style="background: var(--bg-f5f3ff); color: var(--tx-7c3aed);">จุด</em>
              </span>
              <span class="map-summary-meta">{{ meterSummary.anomalyCount ? 'สูงเกินค่าเฉลี่ยของโครงการ' : 'ไม่พบค่าที่ต้องตรวจสอบ' }}</span>
            </span>
          </div>
        </div>

        <!-- 2. Anomalies Alert Table -->
        <div class="dash-card">
          <h3 class="dash-card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#7c3aed" stroke-width="2"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
            แผงที่ตรวจพบความผิดปกติ / ต้องตรวจสอบหน้างาน ({{ meterSummary.anomalyCount }} แผง)
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
                <tr v-for="row in meterAnomalyRows" :key="row.code">
                  <td><span class="table-stall-pill">{{ row.label }}</span></td>
                  <td><span class="table-type-badge badge-electric">น้ำ / ไฟ</span></td>
                  <td><strong>{{ (row.electricUsed || 0).toLocaleString('th-TH') }} kWh</strong> · {{ (row.waterUsed || 0).toLocaleString('th-TH') }} ลบ.ม.</td>
                  <td><span class="table-status-pill pill-anomaly">สูงเกินค่าเฉลี่ยโครงการ</span></td>
                </tr>
                <tr v-if="meterAnomalyRows.length === 0">
                  <td colspan="4" class="dash-empty-row">ไม่พบค่ามิเตอร์ผิดปกติในรอบนี้</td>
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
            <template v-for="project in meterLiveProjects" :key="project.id">
              <div class="dash-breakdown-row">
                <span>{{ project.name }}</span>
                <strong>{{ project.total ? Math.round((project.recorded / project.total) * 100) : 0 }}% <small>({{ project.recorded }}/{{ project.total }} จุด{{ project.recorded === project.total ? ' · ครบแล้ว' : ` · เหลือ ${project.total - project.recorded} จุด` }})</small></strong>
              </div>
              <div class="dash-progress-bar"><div class="dash-progress-fill" :style="{ width: `${project.total ? (project.recorded / project.total) * 100 : 0}%`, background: project.recorded === project.total ? '#059669' : '#2563eb' }"></div></div>
            </template>
            <div v-if="meterLiveProjects.length === 0" class="dash-empty-row">ยังไม่มีแผงที่ต้องจดมิเตอร์ในรอบนี้</div>
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
            v-if="isMapOverviewMode && pageMode === 'payments' && !meterCardCode"
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
            class="market-map-kpi-strip is-meters"
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
              <span class="market-map-kpi-icon" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="6" height="6"/><rect x="14" y="4" width="6" height="6"/><rect x="4" y="14" width="6" height="6"/><rect x="14" y="14" width="6" height="6"/></svg></span>
              <span class="market-map-kpi-copy">
                <small>ทั้งหมด</small>
                <strong>{{ meterSummary.totalCount.toLocaleString('th-TH') }} <em>แผง</em></strong>
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
                <strong>{{ meterSummary.recordedCount.toLocaleString('th-TH') }} <em>แผง</em></strong>
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
                <strong>{{ meterSummary.pendingCount.toLocaleString('th-TH') }} <em>แผง</em></strong>
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
                <strong>{{ meterSummary.overdueCount.toLocaleString('th-TH') }} <em>แผง</em></strong>
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
                <strong>{{ meterSummary.anomalyCount.toLocaleString('th-TH') }} <em>แผง</em></strong>
              </span>
            </button>
          </div>

          <div
            v-else-if="isMapOverviewMode && !meterCardCode"
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
              :class="{ active: dashboardStatusFilter === 'action_required', 'is-loading': isActionDataLoading }"
              :disabled="isActionDataLoading"
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
                <div v-if="isActionDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
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
            :managedDataVersion="mapVersion"
            :selectedProjectId="selectedProject?.id"
            :selectedStallCodes="activeHighlightedStallCodes"
            :stallStatusOverrides="stallStatusOverrides"
            :stallStatusFilter="dashboardStatusFilter === 'legal_case' ? 'all' : dashboardStatusFilter"
            :statusFilterContext="dashboardStatusFilter"
            :actionRequiredStallCodes="actionRequiredStallCodes"
            :operationalFilterMeta="stallOperationalFilterMeta"
            :activeFilterCodes="mapActiveFilterCodes"
            :meterStates="meterMapStates"
            @market-data-loaded="handleMapDataLoaded"
            :isDashboardOpen="dashboardPanelOpen"
            :dashboardTitle="dashboardViewTitle"
            :dashboardSubtitle="dashboardViewSubtitle"
            @update:isDashboardOpen="dashboardPanelOpen = $event"
            @toggleDashboard="dashboardPanelOpen = !dashboardPanelOpen"
            @update:activePanel="activeMapSidePanel = $event"
            @clear-status-filter="dashboardStatusFilter = 'all'"
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
            :filterContext="isCostMode ? 'costs' : 'data'"
            :filterBadge="costFilterBadge"
            :hideRail="costDrawerOpen"
            :hideProjectPanel="isMapOverviewMode && (showInspectStallOverlay || !!meterCardCode) && !isStallDrawerClosing"
            :returnProjectPanel="isMapOverviewMode && isProjectPanelReturning"
            @clearStallSelection="handleClearStallSelection"
            @operational-filter-change="onOperationalFilterChange"
            @filtered-summary="onFilteredSummary"
            @reset-cost-filters="resetCostFilters"
            @selectProject="selectProject"
            @requestExitProject="handleRequestExitProject"
            @exitProject="clearSelectedProject"
            @selectStall="selectStall"
            @updateProjectPolygon="updateProjectPolygon"
            @updateStallPolygon="updateStallPolygon"
            @requestAddProject="openZoneEditor"
            @requestEditProject="openProjectEditor"
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
                        :class="{ active: dashboardStatusFilter === 'action_required', 'is-loading': isActionDataLoading }"
                        :disabled="isActionDataLoading"
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
                          <div v-if="isActionDataLoading" class="market-map-kpi-skeleton" aria-hidden="true">
                            <span class="kpi-skeleton-num"></span>
                          </div>
                          <strong v-else class="ov-dash-status-pill-val">
                            {{ currentDisplayedSummary.actionRequired.toLocaleString('th-TH') }} <small>แผง</small>
                          </strong>
                        </div>
                      </button>
                    </div>
                  </section>

                  <template v-if="isManageDataLoading">
                    <section v-for="n in 4" :key="`dash-skeleton-${n}`" class="ov-dash-section-card ov-dash-skeleton-card" aria-busy="true" aria-label="กำลังโหลดข้อมูลแดชบอร์ด">
                      <span class="kpi-skeleton-num dash-skel-title"></span>
                      <span v-if="n === 1" class="dash-skel-donut" aria-hidden="true"></span>
                      <div v-for="r in 3" :key="r" class="dash-skel-row" aria-hidden="true">
                        <span class="dash-skel-dot"></span>
                        <span class="kpi-skeleton-num dash-skel-line"></span>
                        <span class="kpi-skeleton-pill"></span>
                      </div>
                    </section>
                  </template>
                  <template v-else>
                  <section class="ov-dash-section-card" aria-labelledby="all-project-donuts-heading">
                    <div class="ov-dash-section-head">
                      <div>
                        <h4 id="all-project-donuts-heading" class="ov-dash-section-title">จำนวนแผงแต่ละโครงการ</h4>
                      </div>
                    </div>
                    <div class="ov-hbar-scroll">
                    <div class="ov-hbar-key" aria-hidden="true">
                      <span class="is-available"><i></i>ว่าง</span>
                      <span class="is-reserved"><i></i>ติดจอง</span>
                      <span class="is-occupied"><i></i>มีผู้เช่า</span>
                    </div>
                    <ol class="ov-hbar-list" aria-label="จำนวนแผงแยกตามสถานะของแต่ละโครงการ">
                      <li v-for="project in sortedProjectBars" :key="`hbar-${project.id}`">
                        <button
                          type="button"
                          class="ov-hbar-row"
                          :class="{ active: selectedDashboardProjectId === project.id }"
                          :aria-pressed="selectedDashboardProjectId === project.id"
                          :aria-label="`${project.name}: รวม ${project.total} แผง, ว่าง ${project.available}, ติดจอง ${project.reserved}, มีผู้เช่า ${project.occupied}`"
                          :title="`${project.name} · ว่าง ${project.available.toLocaleString('th-TH')} · ติดจอง ${project.reserved.toLocaleString('th-TH')} · มีผู้เช่า ${project.occupied.toLocaleString('th-TH')}`"
                          @click="selectDashboardProject(project.id)"
                        >
                          <span class="ov-hbar-name">{{ project.name }}</span>
                          <span class="ov-hbar-total">{{ project.total.toLocaleString('th-TH') }} แผง</span>
                          <span class="ov-hbar-track" aria-hidden="true">
                            <i v-if="project.available" class="is-available" :style="{ width: `${(project.available / Math.max(1, project.total)) * 100}%` }"></i>
                            <i v-if="project.reserved" class="is-reserved" :style="{ width: `${(project.reserved / Math.max(1, project.total)) * 100}%` }"></i>
                            <i v-if="project.occupied" class="is-occupied" :style="{ width: `${(project.occupied / Math.max(1, project.total)) * 100}%` }"></i>
                          </span>
                          <span class="ov-hbar-chips">
                            <span class="ov-hbar-chip is-available" :class="{ 'is-zero': !project.available }"><strong>{{ project.available.toLocaleString('th-TH') }}</strong></span>
                            <span class="ov-hbar-chip is-reserved" :class="{ 'is-zero': !project.reserved }"><strong>{{ project.reserved.toLocaleString('th-TH') }}</strong></span>
                            <span class="ov-hbar-chip is-occupied" :class="{ 'is-zero': !project.occupied }"><strong>{{ project.occupied.toLocaleString('th-TH') }}</strong></span>
                          </span>
                        </button>
                      </li>
                    </ol>
                    </div>
                  </section>

                  <section class="ov-dash-section-card ov-popular-card" aria-labelledby="popular-stalls-heading">
                    <div class="ov-dash-section-head">
                      <h4 id="popular-stalls-heading" class="ov-dash-section-title">{{ popularMode === 'stall' ? 'แผงยอดนิยม' : 'โครงการยอดนิยม' }}</h4>
                      <AppSegmentedControl v-model="popularMode" :options="popularModeOptions" aria-label="มุมมองอันดับยอดนิยม" />
                    </div>
                    <div class="ov-popular-toolbar" role="group" aria-label="ช่วงเวลา">
                      <button v-for="days in popularDayOptions" :key="days" type="button" class="ov-popular-range" :class="{ 'is-active': popularDays === days }" :aria-pressed="popularDays === days" @click="popularDays = days">{{ days }} วันล่าสุด</button>
                    </div>
                    <template v-if="popularMode === 'stall'">
                      <ol v-if="popularStallRanking.length" class="ov-popular-list">
                        <li v-for="(item, index) in popularStallRanking" :key="item.code">
                          <button type="button" class="ov-popular-row" :class="[item.hasMedal ? `is-rank-${item.rank}` : '', { 'is-selected': isPopularStallActive(item.code) }]" :aria-pressed="isPopularStallActive(item.code)" :title="`${item.code} · ${item.projectName || ''}`" :aria-label="`อันดับ ${item.rank} แผง ${item.number} ${item.projectName} จองแล้ว ${item.count} ครั้ง`" @click="filterStallFromDashboard(item.code)">
                            <span class="ov-popular-rank">{{ item.rank }}</span>
                            <span class="ov-popular-main"><strong>แผง {{ item.number }}</strong><small>{{ item.projectName || '-' }}</small></span>
                            <span class="ov-popular-count"><b>{{ item.count }}</b> ครั้ง</span>
                          </button>
                        </li>
                      </ol>
                      <p v-if="popularStallAllTied" class="ov-popular-note">ยังไม่มีแผงที่ถูกจองซ้ำ จัดเรียงตามการจองล่าสุด</p>
                      <AppEmptyState v-else-if="!popularStallRanking.length" variant="compact" icon="heart" :title="`ยังไม่มีการจองใน ${popularDays} วันล่าสุด`" description="เมื่อมีผู้จองแผง อันดับจะแสดงที่นี่" />
                    </template>
                    <template v-else>
                      <ol v-if="popularProjectRanking.length" class="ov-popular-list">
                        <li v-for="(item, index) in popularProjectRanking" :key="item.id">
                          <button type="button" class="ov-popular-row" :class="[item.hasMedal ? `is-rank-${item.rank}` : '', { 'is-selected': isPopularProjectActive(item.id) }]" :aria-pressed="isPopularProjectActive(item.id)" :aria-label="`อันดับ ${item.rank} ${item.name} จองแล้ว ${item.count} ครั้ง จาก ${item.stallCount} แผง คิดเป็น ${item.percent.toFixed(1)} เปอร์เซ็นต์ของแผงทั้งโครงการ`" @click="openPopularProject(item.id)">
                            <span class="ov-popular-rank">{{ item.rank }}</span>
                            <span class="ov-popular-main"><strong>{{ item.name }}</strong><small>{{ item.stallCount }} แผง · {{ item.percent >= 10 ? item.percent.toFixed(0) : item.percent.toFixed(1) }}% ของโครงการ</small></span>
                            <span class="ov-popular-count"><b>{{ item.count }}</b> ครั้ง</span>
                          </button>
                        </li>
                      </ol>
                      <p v-if="popularProjectAllTied" class="ov-popular-note">ทุกโครงการถูกจองเท่ากัน ยังไม่มีโครงการที่นำอย่างชัดเจน</p>
                      <AppEmptyState v-else-if="!popularProjectRanking.length" variant="compact" icon="heart" :title="`ยังไม่มีการจองใน ${popularDays} วันล่าสุด`" description="เมื่อมีผู้จองแผง อันดับจะแสดงที่นี่" />
                    </template>
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
                        <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('contractStatus', 'active')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('contractStatus', 'active') }" @click="toggleDashboardFilter('contractStatus', 'active')" @keydown.enter.prevent="toggleDashboardFilter('contractStatus', 'active')"><i class="is-active"></i>อยู่ในสัญญา <strong>{{ contractHealthSummary.normalCount }}</strong></span>
                        <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('contractStatus', 'expiring')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('contractStatus', 'expiring') }" @click="toggleDashboardFilter('contractStatus', 'expiring')" @keydown.enter.prevent="toggleDashboardFilter('contractStatus', 'expiring')"><i class="is-expiring"></i>ใกล้หมดอายุ <strong>{{ contractHealthSummary.expiringCount }}</strong></span>
                        <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('contractStatus', 'expired')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('contractStatus', 'expired') }" @click="toggleDashboardFilter('contractStatus', 'expired')" @keydown.enter.prevent="toggleDashboardFilter('contractStatus', 'expired')"><i class="is-expired"></i>หมดอายุ <strong>{{ contractHealthSummary.expiredCount }}</strong></span>
                      </div>
                    </div>
                    <div v-if="renewalBarTotal" class="ov-renewal-bar" aria-hidden="true">
                      <i class="is-renewed" :style="{ width: `${(stallContractRenewalSummary.submitted / renewalBarTotal) * 100}%` }"></i>
                      <i class="is-awaiting" :style="{ width: `${(stallContractRenewalSummary.pending / renewalBarTotal) * 100}%` }"></i>
                      <i class="is-declined" :style="{ width: `${(stallContractRenewalSummary.declined / renewalBarTotal) * 100}%` }"></i>
                    </div>
                    <div class="ov-contract-renewal-row">
                      <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('renewalStatus', 'approved')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('renewalStatus', 'approved') }" @click="toggleDashboardFilter('renewalStatus', 'approved')" @keydown.enter.prevent="toggleDashboardFilter('renewalStatus', 'approved')">ต่อสัญญาแล้ว <strong>{{ stallContractRenewalSummary.submitted }}</strong></span>
                      <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('renewalStatus', 'awaiting')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('renewalStatus', 'awaiting') }" @click="toggleDashboardFilter('renewalStatus', 'awaiting')" @keydown.enter.prevent="toggleDashboardFilter('renewalStatus', 'awaiting')">รอตอบรับ <strong>{{ stallContractRenewalSummary.pending }}</strong></span>
                      <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('renewalStatus', 'declined')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('renewalStatus', 'declined') }" @click="toggleDashboardFilter('renewalStatus', 'declined')" @keydown.enter.prevent="toggleDashboardFilter('renewalStatus', 'declined')">ไม่ต่อ <strong>{{ stallContractRenewalSummary.declined }}</strong></span>
                    </div>
                  </section>

                  <section class="ov-dash-section-card ov-operational-status-card" aria-labelledby="return-status-heading">
                    <div class="ov-dash-section-head">
                      <h4 id="return-status-heading" class="ov-dash-section-title">การส่งคืนแผงค้า</h4>
                      <small class="ov-dash-section-meta">{{ stallReturnTotal }} รายการ</small>
                    </div>
                    <div class="ov-stage-flow is-four">
                      <div role="button" tabindex="0" :aria-pressed="isDashFilterActive('returnStatus', 'notified')" :class="{ 'is-filterable': true, 'is-zero': stallReturnStatusSummary.notified === 0, 'is-selected': isDashFilterActive('returnStatus', 'notified') }" @click="toggleDashboardFilter('returnStatus', 'notified')" @keydown.enter.prevent="toggleDashboardFilter('returnStatus', 'notified')"><span>แจ้งไม่ต่อ</span><strong>{{ stallReturnStatusSummary.notified }}</strong></div>
                      <div role="button" tabindex="0" :aria-pressed="isDashFilterActive('returnStatus', 'pending')" :class="{ 'is-filterable': true, 'is-zero': stallReturnStatusSummary.waiting === 0, 'is-selected': isDashFilterActive('returnStatus', 'pending') }" @click="toggleDashboardFilter('returnStatus', 'pending')" @keydown.enter.prevent="toggleDashboardFilter('returnStatus', 'pending')"><span>รอส่งคืน</span><strong>{{ stallReturnStatusSummary.waiting }}</strong></div>
                      <div role="button" tabindex="0" :aria-pressed="isDashFilterActive('returnStatus', 'inspection')" :class="{ 'is-filterable': true, 'is-zero': stallReturnStatusSummary.inspection === 0, 'is-selected': isDashFilterActive('returnStatus', 'inspection') }" @click="toggleDashboardFilter('returnStatus', 'inspection')" @keydown.enter.prevent="toggleDashboardFilter('returnStatus', 'inspection')"><span>รอตรวจรับ</span><strong>{{ stallReturnStatusSummary.inspection }}</strong></div>
                      <div class="is-success" role="button" tabindex="0" :aria-pressed="isDashFilterActive('returnStatus', 'returned')" :class="{ 'is-filterable': true, 'is-zero': stallReturnStatusSummary.returned === 0, 'is-selected': isDashFilterActive('returnStatus', 'returned') }" @click="toggleDashboardFilter('returnStatus', 'returned')" @keydown.enter.prevent="toggleDashboardFilter('returnStatus', 'returned')"><span>รับคืนแล้ว</span><strong>{{ stallReturnStatusSummary.returned }}</strong></div>
                    </div>
                    <div class="ov-stage-alert" role="button" tabindex="0" :aria-pressed="isDashFilterActive('returnStatus', 'overdue')" :class="{ 'is-filterable': true, 'is-zero': stallReturnStatusSummary.overdue === 0, 'is-selected': isDashFilterActive('returnStatus', 'overdue') }" @click="toggleDashboardFilter('returnStatus', 'overdue')" @keydown.enter.prevent="toggleDashboardFilter('returnStatus', 'overdue')">
                      <span>เกินกำหนดส่งคืน</span><strong>{{ stallReturnStatusSummary.overdue }}</strong>
                    </div>
                  </section>

                  <section class="ov-dash-section-card ov-operational-status-card" aria-labelledby="legal-status-heading">
                    <div class="ov-dash-section-head">
                      <h4 id="legal-status-heading" class="ov-dash-section-title">การดำเนินคดี</h4>
                      <small class="ov-dash-section-meta">{{ stallLegalTotal }} คดี</small>
                    </div>
                    <div class="ov-stage-flow is-four">
                      <div role="button" tabindex="0" :aria-pressed="isDashFilterActive('legalStatus', 'legal_review')" :class="{ 'is-filterable': true, 'is-zero': stallLegalStatusSummary.awaiting === 0, 'is-selected': isDashFilterActive('legalStatus', 'legal_review') }" @click="toggleDashboardFilter('legalStatus', 'legal_review')" @keydown.enter.prevent="toggleDashboardFilter('legalStatus', 'legal_review')"><span>รอฝ่ายกฎหมาย</span><strong>{{ stallLegalStatusSummary.awaiting }}</strong></div>
                      <div class="is-pending" role="button" tabindex="0" :aria-pressed="isDashFilterActive('legalStatus', 'proceeding')" :class="{ 'is-filterable': true, 'is-zero': stallLegalStatusSummary.proceeding === 0, 'is-selected': isDashFilterActive('legalStatus', 'proceeding') }" @click="toggleDashboardFilter('legalStatus', 'proceeding')" @keydown.enter.prevent="toggleDashboardFilter('legalStatus', 'proceeding')"><span>กำลังดำเนินคดี</span><strong>{{ stallLegalStatusSummary.proceeding }}</strong></div>
                      <div class="is-danger" role="button" tabindex="0" :aria-pressed="isDashFilterActive('legalStatus', 'enforcement')" :class="{ 'is-filterable': true, 'is-zero': stallLegalStatusSummary.enforcement === 0, 'is-selected': isDashFilterActive('legalStatus', 'enforcement') }" @click="toggleDashboardFilter('legalStatus', 'enforcement')" @keydown.enter.prevent="toggleDashboardFilter('legalStatus', 'enforcement')"><span>บังคับคดี</span><strong>{{ stallLegalStatusSummary.enforcement }}</strong></div>
                      <div class="is-success" role="button" tabindex="0" :aria-pressed="isDashFilterActive('legalStatus', 'closed')" :class="{ 'is-filterable': true, 'is-zero': stallLegalStatusSummary.closed === 0, 'is-selected': isDashFilterActive('legalStatus', 'closed') }" @click="toggleDashboardFilter('legalStatus', 'closed')" @keydown.enter.prevent="toggleDashboardFilter('legalStatus', 'closed')"><span>สิ้นสุดคดี</span><strong>{{ stallLegalStatusSummary.closed }}</strong></div>
                    </div>
                  </section>

                  <section class="ov-dash-section-card ov-operational-status-card" aria-labelledby="maintenance-status-heading">
                    <div class="ov-dash-section-head">
                      <h4 id="maintenance-status-heading" class="ov-dash-section-title">งานซ่อมบำรุง</h4>
                    </div>
                    <div class="ov-maintenance-donut-layout">
                      <div
                        class="ov-maintenance-donut"
                        :style="stallMaintenanceDonutStyle"
                        role="img"
                        :aria-label="`งานซ่อมทั้งหมด ${stallMaintenanceTotal} รายการ: รอดำเนินการ ${stallMaintenanceDashboardSummary.pending}, กำลังซ่อม ${stallMaintenanceDashboardSummary.repairing}, ซ่อมเสร็จ ${stallMaintenanceDashboardSummary.completed}, เกินกำหนด ${stallMaintenanceDashboardSummary.overdue}`"
                      >
                        <span><strong>{{ stallMaintenanceTotal.toLocaleString('th-TH') }}</strong><small>งานทั้งหมด</small></span>
                      </div>
                      <div class="ov-maintenance-legend">
                        <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('maintenanceStatus', 'pending')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('maintenanceStatus', 'pending') }" @click="toggleDashboardFilter('maintenanceStatus', 'pending')" @keydown.enter.prevent="toggleDashboardFilter('maintenanceStatus', 'pending')"><i class="is-pending"></i>รอดำเนินการ <strong>{{ stallMaintenanceDashboardSummary.pending }}</strong></span>
                        <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('maintenanceStatus', 'repairing')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('maintenanceStatus', 'repairing') }" @click="toggleDashboardFilter('maintenanceStatus', 'repairing')" @keydown.enter.prevent="toggleDashboardFilter('maintenanceStatus', 'repairing')"><i class="is-repairing"></i>กำลังซ่อม <strong>{{ stallMaintenanceDashboardSummary.repairing }}</strong></span>
                        <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('maintenanceStatus', 'completed')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('maintenanceStatus', 'completed') }" @click="toggleDashboardFilter('maintenanceStatus', 'completed')" @keydown.enter.prevent="toggleDashboardFilter('maintenanceStatus', 'completed')"><i class="is-completed"></i>ซ่อมเสร็จ <strong>{{ stallMaintenanceDashboardSummary.completed }}</strong></span>
                        <span role="button" tabindex="0" :aria-pressed="isDashFilterActive('maintenanceStatus', 'overdue')" :class="{ 'is-filterable': true, 'is-selected': isDashFilterActive('maintenanceStatus', 'overdue') }" @click="toggleDashboardFilter('maintenanceStatus', 'overdue')" @keydown.enter.prevent="toggleDashboardFilter('maintenanceStatus', 'overdue')"><i class="is-overdue"></i>เกินกำหนด <strong>{{ stallMaintenanceDashboardSummary.overdue }}</strong></span>
                      </div>
                    </div>
                  </section>
                  </template>

                </div>

                <!-- MODE 2: แดชบอร์ดการเงินและการจัดเก็บ (payments) -->
                <div v-else-if="pageMode === 'payments'" class="ov-dash-mode-view">
                  <template v-if="isManageDataLoading">
                    <section v-for="n in 2" :key="`dash-skeleton-${n}`" class="ov-dash-section-card ov-dash-skeleton-card" aria-busy="true" aria-label="กำลังโหลดข้อมูลแดชบอร์ด">
                      <span class="kpi-skeleton-num dash-skel-title"></span>
                      <span v-if="n === 1" class="dash-skel-donut" aria-hidden="true"></span>
                      <div v-for="r in 3" :key="r" class="dash-skel-row" aria-hidden="true">
                        <span class="dash-skel-dot"></span>
                        <span class="kpi-skeleton-num dash-skel-line"></span>
                        <span class="kpi-skeleton-pill"></span>
                      </div>
                    </section>
                  </template>
                  <template v-else>
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
                  </template>
                </div>

                <!-- MODE 3: แดชบอร์ดการจดมิเตอร์น้ำและไฟฟ้า (meters) -->
                <div v-else-if="pageMode === 'meters'" class="ov-dash-mode-view">
                  <template v-if="isManageDataLoading">
                    <section v-for="n in 2" :key="`dash-skeleton-${n}`" class="ov-dash-section-card ov-dash-skeleton-card" aria-busy="true" aria-label="กำลังโหลดข้อมูลแดชบอร์ด">
                      <span class="kpi-skeleton-num dash-skel-title"></span>
                      <span v-if="n === 1" class="dash-skel-donut" aria-hidden="true"></span>
                      <div v-for="r in 3" :key="r" class="dash-skel-row" aria-hidden="true">
                        <span class="dash-skel-dot"></span>
                        <span class="kpi-skeleton-num dash-skel-line"></span>
                        <span class="kpi-skeleton-pill"></span>
                      </div>
                    </section>
                  </template>
                  <template v-else>
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
                  </template>
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
            <template v-if="pageMode === 'meters'" #project-body>
              <StallMeterProjectCard :rows="meterProjectStalls" :stats="meterProjectStats" :cycle="meterCycle" :read-only="isPastPeriod" v-model:month="costMonth" :current-month="meterCurrentMonth" @next="openNextMeterStall" />
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
                        :aria-label="selectedProject?.name ? `กลับไป ${selectedProject.name}` : 'กลับ'"
                        @click="closeStall()"
                      >
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="15 18 9 12 15 6"></polyline></svg>
                        <span>{{ selectedProject?.name || 'กลับ' }}</span>
                      </button>
                    </div>
                    <div class="issd-header-top">
                      <div class="issd-title-group">
                        <strong>{{ `แผง ${stallDisplayNumber}` }}</strong>
                        <span
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
                        {{ cleanProjectName(selectedProject?.name || stallEdit.project_name || 'โซนแผงค้า') }}<template v-if="(selectedProject?.floors?.length || 0) > 1"> · {{ formatFloorDisplay(stallEdit.floor) }}</template>
                    </div>
                  </div>

                  <!-- 2. Body -->
                  <!-- 2A. Mini Summary Mode -->
                  <!-- 2A-1. โหมดจัดการ GIS: Drawer แบบเดิม -->
                  <div v-if="(stallDrawerMode === 'detail' || stallDrawerMode === 'edit') && !isInspectMode" class="issd-body">
                    <!-- Stall Photo Banner -->
                    <div
                      class="issd-photo-banner"
                      title="คลิกเพื่อดูภาพแผงเช่าทั้งหมด"
                      @click="openInspectPhotoPreview(0)"
                    >
                      <img :src="inspectStallPhotos[0]" :alt="`ภาพแผง ${stallDisplayNumber}`" />
                      <div class="issd-photo-count-badge">
                        <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                        <span>{{ inspectStallPhotos.length }} รูป</span>
                      </div>
                    </div>

                    <!-- 3 ช่องสรุป: ขนาด | ค่าเช่า | แรกเข้า (เหมือนการ์ดฝั่งผู้เช่า) -->
                    <div class="issd-metrics3">
                      <div>
                        <span class="issd-metric-label">ขนาด</span>
                        <strong class="issd-metric-value">{{ inspectStallSizeSqm }} <small>ตร.ม.</small></strong>
                      </div>
                      <div>
                        <span class="issd-metric-label">ค่าเช่า</span>
                        <strong class="issd-metric-value is-rent">฿{{ Number(stallEdit.monthly_rent || 1800).toLocaleString('th-TH') }}<small>/เดือน</small></strong>
                      </div>
                      <div>
                        <span class="issd-metric-label">แรกเข้า</span>
                        <strong class="issd-metric-value">฿{{ Number(inspectStallEntryFee).toLocaleString('th-TH') }}</strong>
                      </div>
                    </div>

                    <!-- ข้อมูลที่เหลือทั้งหมด -->
                    <dl class="issd-detail-list">
                      <div class="issd-detail-row">
                        <dt>การรับจอง</dt>
                        <dd><span class="issd-bookable-pill" :class="{ 'is-paused': stallEdit.bookable === false }">{{ stallEdit.bookable === false ? 'ปิดรับจอง' : 'เปิดรับจอง' }}</span></dd>
                      </div>
                      <div class="issd-detail-row">
                        <dt>หมวดสินค้า</dt>
                        <dd>{{ stallEdit.category || selectedProject?.category || 'ทั่วไป' }}</dd>
                      </div>
                      <div class="issd-detail-row">
                        <dt>ขนาดแผง</dt>
                        <dd>{{ stallEdit.width }} × {{ stallEdit.depth }} ม.</dd>
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
                      <img :src="inspectStallPhotos[0]" :alt="`ภาพแผง ${stallDisplayNumber}`" />
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
                          <span class="issd-aac-count-badge">{{ stallActionTotalCount }} รายการ</span>
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
                      <button type="button" class="issd-close-btn issd-full-close" title="ปิดและกลับไปโครงการ" aria-label="ปิดรายละเอียดแผง กลับไปโครงการ" @click="closeStall">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                      </button>
                      <!-- 1.1 Action Alert Card (at top of content area) -->
                      <div v-if="fullDetailSettling" class="issd-full-skeleton" aria-hidden="true">
                        <span></span><span></span><span></span>
                      </div>
                      <div
                        v-if="stallActionRichItems.length > 0 && !fullDetailSettling"
                        class="issd-action-alert-card"
                      >
                        <div class="issd-action-card-rail" role="group" aria-label="รายการสิ่งที่ต้องดำเนินการ">
                          <button
                            v-for="act in stallActionRichItems"
                            :key="act.key"
                            type="button"
                            class="issd-action-summary-card"
                            :class="`is-${act.severity}`"
                            @click="jumpToActionTab(act.targetTab)"
                            :title="`คลิกเพื่อไปยังแท็บ: ${act.actionLabel}`"
                          >
                            <div class="issd-action-summary-copy">
                              <span class="issd-action-summary-title">
                                <strong>{{ act.category }}</strong>
                                <b v-if="act.count" class="issd-aac-count" :class="`is-${act.severity}`">{{ act.count }}</b>
                              </span>
                              <small>{{ act.title }}</small>
                            </div>
                            <svg class="issd-action-summary-chevron" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 18 15 12 9 6"></polyline></svg>
                          </button>
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
                        <img :src="inspectStallPhotos[0]" :alt="`ภาพแผง ${stallDisplayNumber}`" />
                        <div class="issd-photo-count-badge">
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                          <span>{{ inspectStallPhotos.length }} รูป</span>
                        </div>
                      </div>

                      <!-- KPI Strip 3 ช่อง (ข้อมูลหลักเท่านั้น) -->
                      <div class="issd-kpi-grid is-three is-primary-stats">
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

                      <!-- รายละเอียดพื้นที่: ใช้การ์ดเดียวและรายการข้อมูลกะทัดรัด เพื่อลดการ์ดซ้อนและรองรับข้อความยาว -->
                      <section class="issd-stall-detail-card" aria-labelledby="stall-space-detail-heading">
                        <div class="issd-stall-detail-head">
                          <span class="issd-stall-detail-icon" aria-hidden="true">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18M3 9h6M3 15h6M15 8h2M15 12h2M15 16h2"/></svg>
                          </span>
                          <strong id="stall-space-detail-heading">รายละเอียดพื้นที่</strong>
                        </div>
                        <dl class="issd-stall-facts-grid">
                          <div class="issd-stall-fact">
                            <dt>ชื่อแผง</dt>
                            <dd>{{ stallDisplayNumber || stallEdit.name || '-' }}</dd>
                          </div>
                          <div class="issd-stall-fact">
                            <dt>ชั้น</dt>
                            <dd>{{ formatFloorDisplay(stallEdit.floor) }}<small>{{ stallEdit.building ? `อาคาร ${stallEdit.building}` : 'ระดับพื้นที่' }}</small></dd>
                          </div>
                          <div class="issd-stall-fact is-wide">
                            <dt>โครงการ</dt>
                            <dd>{{ selectedProject?.name || stallEdit.project_name || stallEdit.zone_name || '-' }}<small>ผังโครงการ</small></dd>
                          </div>
                          <div class="issd-stall-fact is-wide">
                            <dt>หมวดสินค้า</dt>
                            <dd>{{ stallEdit.category || selectedProject?.category || 'แผงค้ามาตรฐาน' }}<small>ประเภทสินค้า</small></dd>
                          </div>
                          <div v-if="stallEdit.description" class="issd-stall-fact is-wide">
                            <dt>รายละเอียดเพิ่มเติม</dt>
                            <dd>{{ stallEdit.description }}</dd>
                          </div>
                        </dl>
                      </section>

                      <!-- กล่องข้อมูลผู้เช่า / คำร้องขอเช่า / สถานะอื่นๆ (ไม่แสดงผู้ขอเข้าชมในแท็บข้อมูลแผง) -->
                      <!-- 2. กรณีมีคำร้องขอเช่า -->
                      <div
                        v-if="stallEdit.stall_status !== 'occupied' && !currentStallContract && currentStallApplication"
                        class="issd-card-section issd-occupant-card is-reserved"
                      >
                        <!-- ส่วนผู้ยื่นคำร้อง -->
                        <div class="issd-person-row">
                          <div class="issd-person-avatar is-reserved">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                          </div>
                          <div class="issd-person-info">
                            <span class="issd-person-role">ผู้ยื่นคำร้องขอเช่า</span>
                            <strong class="issd-person-name">{{ currentStallApplication.applicant_name || 'ผู้ยื่นคำร้องขอเช่า' }}</strong>
                          </div>
                          <a
                            v-if="currentStallApplication.applicant_phone"
                            :href="`tel:${currentStallApplication.applicant_phone}`"
                            class="issd-phone-link"
                            title="คลิกเพื่อโทรออก"
                          >
                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                            <span>{{ currentStallApplication.applicant_phone }}</span>
                          </a>
                        </div>

                        <!-- ส่วนข้อมูลคำร้อง -->
                        <div class="issd-subdoc-box is-reserved">
                          <div class="issd-subdoc-row">
                            <span class="subdoc-k">เลขที่คำร้อง</span>
                            <strong class="subdoc-v is-code">{{ currentStallApplication.application_number || '-' }}</strong>
                          </div>
                          <div class="issd-subdoc-row">
                            <span class="subdoc-k">สถานะรายการ</span>
                            <span class="subdoc-v">รอพิจารณาอนุมัติ</span>
                          </div>
                        </div>
                        <div class="issd-occupant-action-strip">
                          <NuxtLink
                            :to="withStallReturn(currentStallApplication.id ? `/applications/${currentStallApplication.id}` : `/applications?search=${encodeURIComponent(currentStallApplication.application_number || stallEdit.ident)}`)"
                            class="issd-occupant-action-btn is-app"
                          >
                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                            <span>ดูคำร้องขอเช่าแผง</span>
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

                      <div v-else-if="stallEdit.stall_status && stallEdit.stall_status !== 'available' && !currentStallBooking" class="issd-occupancy-box" :class="`is-${stallEdit.stall_status}`">
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
                      <AppEmptyState v-else class="issd-tab-empty-state" variant="compact" icon="history" title="ยังไม่มีประวัติ" description="กิจกรรมของแผงจะถูกบันทึกไว้ที่นี่" />
                    </template>

                    <!-- Tab: การจองและคำร้อง (bookings) -->
                    <template v-else-if="stallDetailTab === 'bookings'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">การจอง</strong>
                        </div>
                      </div>

                      <AppSecondaryTabs
                        v-model="bookingTypeFilter"
                        :options="[
                          { value: 'all', label: 'ทั้งหมด', count: stallBookingsAndApps.length },
                          { value: 'booking', label: 'การจอง', count: stallBookingsAndApps.filter(i => i.typeKey === 'booking').length },
                          { value: 'application', label: 'คำร้อง', count: stallBookingsAndApps.filter(i => i.typeKey === 'application').length }
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
                            <div class="issd-booking-kind" :class="`is-${item.typeKey}`">
                              <span class="issd-booking-kind-icon" aria-hidden="true">
                                <svg v-if="item.typeKey === 'application'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="8" y1="13" x2="16" y2="13"/><line x1="8" y1="17" x2="14" y2="17"/></svg>
                                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                              </span>
                              <span class="issd-booking-kind-copy">
                                <strong>{{ item.typeLabel }}</strong>
                                <small>{{ item.refNumber }}</small>
                              </span>
                            </div>
                            <span class="issd-status-pill issd-booking-status" :class="item.badgeClass">
                              <span class="issd-booking-status-dot" aria-hidden="true"></span>
                              {{ item.statusText }}
                            </span>
                          </div>

                          <div class="issd-booking-schedule">
                            <div class="issd-booking-date-line">
                              <strong>{{ item.typeKey === 'application' ? `ยื่นเมื่อ ${formatThaiDateOnly(item.rawDate)}` : item.dateText }}</strong>
                              <span v-if="item.typeKey === 'booking' && getBookingCountdownText(item.rawDate)" class="issd-booking-countdown">{{ getBookingCountdownText(item.rawDate) }}</span>
                              <span v-else-if="item.typeKey === 'application' && getApplicationWaitingText(item.rawDate)" class="issd-booking-countdown is-application">{{ getApplicationWaitingText(item.rawDate) }}</span>
                            </div>
                          </div>

                          <div class="issd-booking-person-row">
                            <div class="issd-booking-applicant-left">
                              <div class="issd-booking-avatar" aria-hidden="true">
                                {{ item.applicantName.trim().charAt(0) || 'ผ' }}
                              </div>
                              <div class="issd-booking-applicant-info">
                                <strong class="issd-card-title">{{ item.applicantName }}</strong>
                                <a v-if="item.phone" :href="`tel:${item.phone}`" class="issd-booking-phone">{{ item.phone }}</a>
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
                            <div class="issd-booking-info-item">
                              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                              <span><strong>{{ item.dateText }}</strong></span>
                            </div>
                            <AppButton
                              :to="withStallReturn(item.targetUrl)"
                              size="sm"
                              variant="primary"
                              class="issd-booking-detail-btn"
                            >
                              <template #prefix>
                                <svg v-if="item.typeKey === 'application'" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                                <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                              </template>
                              {{ item.typeKey === 'application' ? 'ตรวจสอบคำร้อง' : 'ตรวจสอบการจอง' }}
                              <template #suffix>
                                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
                              </template>
                            </AppButton>
                          </div>
                          <div v-if="item.typeKey === 'application' && item.sourceRef" class="issd-application-source">
                            มาจากนัดเข้าชม {{ item.sourceRef }}
                          </div>
                        </div>
                      </div>
                      <AppEmptyState
                        v-else
                        class="issd-tab-empty-state"
                        variant="compact"
                        :icon="bookingTypeFilter === 'booking' ? 'calendar' : 'document'"
                        :title="bookingTypeFilter === 'application' ? 'ยังไม่มีคำร้อง' : (bookingTypeFilter === 'booking' ? 'ยังไม่มีรายการจอง' : 'ยังไม่มีรายการจองหรือคำร้อง')"
                        description="รายการที่เกี่ยวข้องกับแผงนี้จะแสดงที่นี่"
                      />
                    </template>

                    <!-- Tab: สัญญา (contracts) -->
                    <template v-else-if="stallDetailTab === 'contracts'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">สัญญา</strong>
                        </div>
                      </div>

                      <AppSecondaryTabs
                        v-model="contractSubTab"
                        :options="[{ value: 'active', label: 'ปัจจุบัน', count: activeStallContracts.length }, { value: 'history', label: 'ประวัติ', count: historyStallContracts.length }]"
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
                            <span
                              class="issd-status-pill issd-contract-status"
                              :class="contract.status === 'active' ? 'is-available' : ((contract.status === 'cancelled' || contract.status === 'terminated') ? 'is-cancelled' : 'is-reserved')"
                            >
                              <span class="issd-contract-status-dot" aria-hidden="true"></span>
                              {{ contract.status === 'active' ? 'ใช้งานอยู่' : (contract.status === 'expired' ? 'สิ้นสุดสัญญา' : ((contract.status === 'cancelled' || contract.status === 'terminated') ? 'ยกเลิกแล้ว' : (contract.status || 'ร่างสัญญา'))) }}
                            </span>
                            <strong class="issd-contract-num">{{ contract.contract_number || `#${contract.id}` }}</strong>
                          </div>

                          <div class="issd-compact-card-body">
                            <div class="issd-contract-overview">
                              <div class="issd-contract-identity">
                                <strong>{{ contract.tenant_name || contract.name || '-' }}</strong>
                              </div>
                              <div class="issd-contract-period-row">
                                <span>{{ formatThaiDateOnly(contract.start_date) }} – {{ formatThaiDateOnly(contract.end_date) }} · {{ getShortContractDurationText(contract.start_date, contract.end_date) }}</span>
                                <strong v-if="getContractRemainingDays(contract.end_date) !== null">เหลือ {{ getContractRemainingDays(contract.end_date) }} วัน</strong>
                              </div>
                              <div class="issd-contract-progress" aria-hidden="true">
                                <span :style="{ width: `${getContractProgressPercent(contract.start_date, contract.end_date)}%` }"></span>
                              </div>
                              <div class="issd-contract-money-summary">
                                <div>
                                  <span>ค่าเช่า</span>
                                  <strong>฿{{ Number(contract.monthly_rental || contract.monthly_rent || 0).toLocaleString('th-TH') }} <small>/ เดือน</small></strong>
                                </div>
                                <div>
                                  <span>เงินประกัน</span>
                                  <strong>฿{{ Number(contract.security_deposit || contract.deposit_amount || 0).toLocaleString('th-TH') }} <small class="is-paid">· รับแล้ว</small></strong>
                                </div>
                              </div>
                            </div>

                            <div v-if="isInspectMode && ['active', 'expired'].includes(contract.status)" class="issd-workflow">
                              <strong class="issd-workflow-title">ติดตามสัญญา</strong>
                              <div class="issd-workflow-grid">
                                <label v-for="field in contractWorkflowFields" :key="field.key" class="issd-workflow-field">
                                  <span>{{ field.label }}</span>
                                  <AppSelect :model-value="workflowValue(contract, field.key)" :options="field.options" @update:model-value="setWorkflowValue(contract, field.key, $event)" />
                                </label>
                              </div>
                              <AppButton v-if="isWorkflowDirty(contract)" size="sm" variant="primary" :loading="contractWorkflowSaving === String(contract.id)" @click="saveContractWorkflow(contract)">บันทึกสถานะ</AppButton>
                            </div>

                            <div class="issd-card-actions-row">
                              <AppButton size="sm" variant="outline" @click="openEditContractModal(contract)">ดูสัญญา</AppButton>
                              <AppButton v-if="contract.status === 'active'"
                                size="sm" variant="danger-outline" @click="terminateContract(contract)">ยกเลิกสัญญา</AppButton>
                            </div>
                          </div>
                        </div>
                      </div>
                      <AppEmptyState
                        v-else
                        class="issd-tab-empty-state"
                        variant="compact"
                        icon="document"
                        :title="contractSubTab === 'active' ? 'ยังไม่มีสัญญาปัจจุบัน' : 'ยังไม่มีประวัติสัญญา'"
                        description="ข้อมูลสัญญาของแผงนี้จะแสดงที่นี่"
                      />
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
                      <AppEmptyState v-else class="issd-tab-empty-state" variant="compact" icon="users" title="ยังไม่มีผู้เช่าช่วง" description="เพิ่มผู้เช่าช่วงเพื่อเริ่มจัดการข้อมูล" />
                    </template>

                    <!-- Tab: ประวัติชำระเงิน (payments) -->
                    <template v-else-if="stallDetailTab === 'payments'">
                      <div class="issd-tab-header-row is-inline-controls">
                        <div class="issd-tab-header-copy">
                          <strong class="issd-tab-title">ชำระเงิน</strong>
                        </div>
                        <div v-if="stallContractsList.length" class="issd-contract-inline-filter">
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
                      <AppEmptyState v-else class="issd-tab-empty-state" variant="compact" icon="document" title="ยังไม่มีประวัติการชำระเงิน" description="รายการชำระเงินของแผงนี้จะแสดงที่นี่" />
                    </template>

                    <!-- Tab: ประวัติการซ่อม / ขออนุญาต (maintenance) -->
                    <template v-else-if="stallDetailTab === 'maintenance'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">งานซ่อม</strong>
                        </div>
                      </div>

                      <AppSecondaryTabs
                        v-model="maintenanceSourceFilter"
                        :options="[
                          { value: 'all', label: 'ทั้งหมด', count: stallMaintenance.length },
                          { value: 'tenant', label: 'ผู้เช่า', count: stallMaintenance.filter(item => getMaintenanceSource(item) === 'tenant').length },
                          { value: 'staff', label: 'เจ้าหน้าที่', count: stallMaintenance.filter(item => getMaintenanceSource(item) === 'staff').length }
                        ]"
                        aria-label="กรองรายการตามผู้แจ้ง"
                      />

                      <div v-if="filteredStallMaintenance.length > 0" class="issd-maintenance-list issd-scroll-list">
                        <div
                          v-for="m in filteredStallMaintenance"
                          :key="m.id"
                          class="issd-maintenance-card"
                          :class="`is-${m.status}`"
                        >
                          <div class="issd-card-top-row">
                            <span class="issd-status-pill issd-maintenance-status" :class="m.status === 'completed' ? 'is-available' : (m.status === 'rejected' ? 'is-occupied' : 'is-reserved')">
                              <span class="issd-maintenance-status-dot" aria-hidden="true"></span>
                              {{ formatMaintenanceStatus(m.status) }}
                            </span>
                            <div class="issd-maintenance-type-label">
                              <span aria-hidden="true">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                  <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.8-3.8a6 6 0 0 1-7.9 7.9l-6.9 6.9a2.1 2.1 0 0 1-3-3l6.9-6.9a6 6 0 0 1 7.9-7.9z" />
                                </svg>
                              </span>
                              {{ formatMaintenanceType(m.type) }}
                            </div>
                          </div>
                          <div class="issd-maintenance-identity">
                            <strong>{{ m.title }}</strong>
                          </div>
                          <div class="issd-compact-card-body">
                            <div v-if="m.status === 'pending'" class="issd-maintenance-meta-grid">
                              <div class="issd-maintenance-meta-item"><span>ผู้แจ้ง</span><strong>{{ m.requesterName }}</strong></div>
                              <div class="issd-maintenance-meta-item"><span>วันที่แจ้ง</span><strong>{{ formatMaintenanceDate(m.requestDate) }}</strong></div>
                              <div class="issd-maintenance-meta-item"><span>รอรับเรื่อง</span><strong class="is-waiting">{{ maintenanceWaitingDays(m) }} วัน</strong></div>
                            </div>

                            <div v-else-if="m.status === 'in_progress'" class="issd-maintenance-meta-grid">
                              <div class="issd-maintenance-meta-item"><span>ผู้รับผิดชอบ</span><strong>{{ m.assigneeName || m.technicianName || formatMaintenanceAssignee(m.assigneeType) }}</strong></div>
                              <div class="issd-maintenance-meta-item"><span>รับเรื่องเมื่อ</span><strong>{{ formatMaintenanceDate(m.acceptedAt || m.requestDate) }}</strong></div>
                              <div class="issd-maintenance-meta-item"><span>กำหนดเสร็จ</span><strong>{{ formatMaintenanceDate(m.scheduledDate) }}</strong></div>
                            </div>

                            <div v-else-if="m.status === 'completed'" class="issd-maintenance-meta-grid">
                              <div class="issd-maintenance-meta-item"><span>แก้ไขโดย</span><strong>{{ m.executorName || m.technicianName || 'ทีมซ่อมบำรุงตลาด' }}</strong></div>
                              <div class="issd-maintenance-meta-item"><span>เสร็จเมื่อ</span><strong class="is-success">{{ formatMaintenanceDate(m.completionDate || m.completedAt || m.scheduledDate) }}</strong></div>
                              <div class="issd-maintenance-meta-item"><span>ใช้เวลา</span><strong>{{ maintenanceDurationDays(m) }} วัน</strong></div>
                            </div>

                            <div v-else class="issd-maintenance-meta-grid">
                              <div class="issd-maintenance-meta-item"><span>ผู้แจ้ง</span><strong>{{ m.requesterName }}</strong></div>
                              <div class="issd-maintenance-meta-item"><span>วันที่แจ้ง</span><strong>{{ formatMaintenanceDate(m.requestDate) }}</strong></div>
                              <div class="issd-maintenance-meta-item"><span>ปิดเรื่องโดย</span><strong>{{ m.rejectedBy || 'เจ้าหน้าที่ตลาด' }}</strong></div>
                            </div>

                            <div v-if="m.status === 'in_progress' && (m.expensePayer || m.priority || m.tenantMessage)" class="issd-maintenance-assignment-note">
                              <span v-if="m.expensePayer">ผู้รับผิดชอบค่าใช้จ่าย: <strong>{{ formatMaintenancePayer(m.expensePayer) }}</strong></span>
                              <span v-if="m.priority">ความเร่งด่วน: <strong>{{ formatMaintenancePriority(m.priority) }}</strong></span>
                              <p v-if="m.tenantMessage">ข้อความถึงผู้เช่า: {{ m.tenantMessage }}</p>
                            </div>

                            <div v-if="m.status === 'completed'" class="issd-maintenance-result">
                              <div class="issd-maintenance-photo-compare">
                                <figure>
                                  <div class="issd-maintenance-photo-frame is-before">
                                    <img v-if="m.beforePhotos?.[0]?.dataUrl" :src="m.beforePhotos[0].dataUrl" alt="รูปก่อนแก้" />
                                    <span v-else>ไม่มีรูปก่อนแก้</span>
                                  </div>
                                  <figcaption>ก่อนแก้ไข</figcaption>
                                </figure>
                                <figure>
                                  <div class="issd-maintenance-photo-frame is-after">
                                    <img v-if="m.afterPhotos?.[0]?.dataUrl" :src="m.afterPhotos[0].dataUrl" alt="รูปหลังแก้" />
                                    <span v-else>{{ m.afterPhotos?.length ? `รูปหลังแก้ ${m.afterPhotos.length} รูป` : 'ไม่มีรูปหลังแก้' }}</span>
                                  </div>
                                  <figcaption>หลังแก้ไข</figcaption>
                                </figure>
                              </div>
                              <p v-if="m.repairSummary"><strong>สิ่งที่แก้:</strong> {{ m.repairSummary }}</p>
                              <p v-if="m.actualCost !== undefined"><strong>ค่าใช้จ่ายจริง:</strong> ฿{{ Number(m.actualCost).toLocaleString('th-TH') }}</p>
                            </div>

                            <div v-if="m.status === 'rejected'" class="issd-maintenance-rejection-summary">
                              <strong>เหตุผลที่ไม่รับเรื่อง:</strong>
                              <span>{{ formatMaintenanceRejectReason(m.rejectionReason) }}<template v-if="m.rejectionDetail"> · {{ m.rejectionDetail }}</template></span>
                            </div>

                            <div class="issd-maintenance-footer-row">
                              <div v-if="m.status === 'pending'" class="issd-card-actions-row">
                                <AppButton size="sm" variant="secondary" @click="openMaintenanceAction(m, 'reject')">ไม่รับเรื่อง</AppButton>
                                <AppButton size="sm" variant="primary" @click="openMaintenanceAction(m, 'accept')">รับเรื่อง</AppButton>
                              </div>
                              <div v-else-if="m.status === 'in_progress'" class="issd-card-actions-row">
                                <AppButton size="sm" variant="primary" @click="openMaintenanceAction(m, 'complete')">บันทึกว่าเสร็จสิ้น</AppButton>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                      <AppEmptyState v-else class="issd-tab-empty-state" variant="compact" icon="inbox" title="ยังไม่มีงานซ่อม" description="รายการแจ้งซ่อมและขออนุญาตจะแสดงที่นี่" />
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

                  <!-- 3. Footer for Detail (Mini) and Edit Modes -->
                  <div v-if="stallDrawerMode === 'detail' || stallDrawerMode === 'edit'" class="issd-footer">
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
              </div>
            </template>
          </MarketProjectsOverviewMap>
        </article>
      </section>
    </main>

    <MaintenanceActionModal
      :model-value="Boolean(activeMaintenanceAction)"
      :mode="activeMaintenanceAction?.mode || 'accept'"
      :item="activeMaintenanceItem"
      :item-title="activeMaintenanceItem?.title"
      @update:model-value="handleMaintenanceModalVisibility"
      @submit="submitMaintenanceAction"
    />

    <Teleport to="body">
      <AppModal
        v-model="showZoneEditor"
        icon="map"
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
          <AppButton variant="secondary" size="md" @click="dismissZoneEditor">ยกเลิก</AppButton>
          <AppButton size="md" @click="createZone">บันทึกโซน</AppButton>
        </template>
      </AppModal>

      <AppModal
        v-model="showStallEditModal"
        size="lg"
        icon="pin"
        :title="isCreatingNewStall ? 'เพิ่มแผงเช่าใหม่' : `แก้ไขแผง ${stallDisplayNumber}`"
        :description="isCreatingNewStall ? 'ตั้งค่ารหัส ขนาด และอัตราค่าเช่าสำหรับแผงใหม่' : 'ข้อมูล ขนาด และราคาของแผง'"
      >
        <form class="issd-form-grid stall-edit-form" novalidate @submit.prevent="saveStall">
                              <h4 class="issd-form-section full">ข้อมูลแผง</h4>
                              <AppImageUpload
                                v-model="editablePhotoFiles"
                                class="full"
                                label="รูปภาพแผง"
                                :hint="`${editablePhotos.length}/${MAX_STALL_PHOTOS} รูป · JPG, PNG, WebP ไม่เกิน 2MB ต่อรูป`"
                                multiple
                                :max-files="MAX_STALL_PHOTOS"
                                preview-alt="ภาพแผง"
                              />
                              <label>
                                <span>รหัสแผง <b v-if="isCreatingNewStall" class="req-star">*</b></span>
                                <AppInput v-model="stallEdit.ident" placeholder="เช่น A-001" :disabled="!isCreatingNewStall" :error="stallErrors.ident" @blur="validateStallField('ident')" />
                              </label>
                              <label>
                                <span>ชื่อแผง</span>
                                <AppInput v-model="stallEdit.name" placeholder="ระบุชื่อแผง (ถ้ามี)" />
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
        
                              <h4 class="issd-form-section full">ขนาดและตำแหน่ง</h4>
                              <label>
                                <span>กว้าง (เมตร) <b class="req-star">*</b></span>
                                <AppInput v-model="stallEdit.width" type="number" placeholder="เช่น 2.0" :error="stallErrors.width" @blur="validateStallField('width')" />
                              </label>
                              <label>
                                <span>ลึก (เมตร) <b class="req-star">*</b></span>
                                <AppInput v-model="stallEdit.depth" type="number" placeholder="เช่น 2.0" :error="stallErrors.depth" @blur="validateStallField('depth')" />
                              </label>
                              <label>
                                <span>โครงการ</span>
                                <AppInput v-model="stallEdit.project_name" placeholder="ชื่อโครงการ" disabled />
                              </label>
                              <label>
                                <span>ชั้น</span>
                                <AppSelect v-model="stallEdit.floor" :options="floorOptions" :disabled="(selectedProject?.floors?.length || 1) <= 1" />
                              </label>
        
                              <h4 class="issd-form-section full">ราคาและการรับจอง</h4>
                              <label>
                                <span>ค่าเช่ารายเดือน (บาท) <b class="req-star">*</b></span>
                                <AppInput v-model="stallEdit.monthly_rent" type="number" placeholder="เช่น 1800" :error="stallErrors.monthly_rent" @blur="validateStallField('monthly_rent')" />
                              </label>
                              <label>
                                <span>ค่าเช่ารายวัน (บาท)</span>
                                <AppInput v-model="stallEdit.daily_rent" type="number" placeholder="เช่น 150" />
                              </label>
                              <label>
                                <span>การรับจอง</span>
                                <AppSelect v-model="stallBookableModel" :options="stallBookableOptions" />
                              </label>
                              <label>
                                <span>สถานะ</span>
                                <AppSelect v-model="stallEdit.stall_status" :options="stallStatusOptions" :disabled="stallEdit.stall_status === 'occupied'" />
                              </label>
                            </form>
        <template #footer>
          <AppButton v-if="!isCreatingNewStall" class="stall-edit-delete" variant="danger-outline" size="md" @click="deleteStall()">
            <template #prefix>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6M10 11v5M14 11v5"/></svg>
            </template>
            ลบแผง
          </AppButton>
          <AppButton variant="secondary" size="md" @click="cancelStallEdit">ยกเลิก</AppButton>
          <AppButton :loading="savingStall" size="md" @click="saveStall">{{ isCreatingNewStall ? 'เพิ่มแผงเช่า' : 'บันทึกการแก้ไข' }}</AppButton>
        </template>
      </AppModal>

      <AppModal
        v-model="showProjectEditor"
        icon="map"
        title="แก้ไขข้อมูลโครงการ"
        description="ชื่อ พื้นที่ และขนาดแผงมาตรฐานของโครงการนี้"
      >
        <form class="zone-form" novalidate @submit.prevent="saveProjectEdit">
          <label>
            <span>ชื่อโครงการ <b class="req-star">*</b></span>
            <AppInput v-model="projectForm.name" placeholder="เช่น โครงการที่ 17" :error="projectErrors.name" @blur="validateProjectField('name')" />
          </label>
          <label>
            <span>พื้นที่ในโครงการ (ตร.ม.)</span>
            <AppInput v-model="projectForm.area" type="number" placeholder="เช่น 1800" :error="projectErrors.area" @blur="validateProjectField('area')" />
          </label>
          <label>
            <span>ขนาดแผงมาตรฐาน กว้าง (เมตร)</span>
            <AppInput v-model="projectForm.stallWidth" type="number" placeholder="เช่น 2.5" :error="projectErrors.stallWidth" @blur="validateProjectField('stallWidth')" />
          </label>
          <label>
            <span>ขนาดแผงมาตรฐาน ลึก (เมตร)</span>
            <AppInput v-model="projectForm.stallDepth" type="number" placeholder="เช่น 2.0" :error="projectErrors.stallDepth" @blur="validateProjectField('stallDepth')" />
          </label>
        </form>
        <template #footer>
          <AppButton variant="secondary" size="md" @click="showProjectEditor = false">ยกเลิก</AppButton>
          <AppButton size="md" @click="saveProjectEdit">บันทึก</AppButton>
        </template>
      </AppModal>

      <AppModal v-model="showFacilityEditor" icon="pin" title="เพิ่มจุดอำนวยความสะดวก" description="ตำแหน่งถูกกำหนดจากจุดที่เลือกบนผังแล้ว">
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
          <AppButton variant="secondary" size="md" @click="closeFacilityEditor">ยกเลิก</AppButton>
          <AppButton size="md" @click="createFacility">เพิ่มจุดบริการ</AppButton>
        </template>
      </AppModal>

      <!-- Modal: เพิ่มผู้เช่าช่วง -->
      <AppModal
        v-model="showAddSubtenantModal"
        icon="users"
        title="เพิ่มผู้เช่าช่วง"
        :description="`บันทึกข้อมูลผู้เช่าช่วงสำหรับแผงค้า ${stallDisplayNumber || '-'}`"
      >
        <form class="issd-modal-form" novalidate @submit.prevent="saveSubtenant">
          <AppImageUpload
            v-model="subtenantForm.photoUrl"
            class="full"
            label="รูปถ่ายผู้เช่าช่วง / หน้าร้าน"
            preview-alt="ภาพตัวอย่างผู้เช่าช่วงหรือหน้าร้าน"
          />
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
          <AppButton variant="secondary" size="md" @click="showAddSubtenantModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" size="md" @click="saveSubtenant">บันทึกผู้เช่าช่วง</AppButton>
        </template>
      </AppModal>

      <!-- Modal: แจ้งซ่อม / ขออนุญาต -->
      <AppModal
        v-model="showAddMaintenanceModal"
        icon="tool"
        title="แจ้งซ่อม / ขออนุญาตดำเนินการ"
        :description="`บันทึกรายการสำหรับแผงค้า ${stallDisplayNumber || '-'}`"
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
          <AppButton variant="secondary" size="md" @click="showAddMaintenanceModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" size="md" @click="saveMaintenance">บันทึกรายการ</AppButton>
        </template>
      </AppModal>

      <!-- Modal: บันทึกการชำระเงิน -->
      <AppModal
        v-model="showAddPaymentModal"
        icon="payment"
        title="บันทึกการชำระเงิน"
        :description="`บันทึกประวัติรับชำระสำหรับแผงค้า ${stallDisplayNumber || '-'}`"
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
          <AppButton variant="secondary" size="md" @click="showAddPaymentModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" size="md" @click="savePayment">บันทึกการชำระเงิน</AppButton>
        </template>
      </AppModal>

      <!-- Modal: จัดการสัญญาเช่า (สร้างใหม่ / แก้ไข) -->
      <AppModal
        v-model="showContractModal"
        icon="document"
        :title="isEditingContract ? 'แก้ไขสัญญาเช่า' : 'ออกสัญญาเช่าใหม่'"
        :description="`จัดการสัญญาเช่าสำหรับแผงค้า ${stallDisplayNumber || '-'}`"
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
          <AppButton variant="secondary" size="md" @click="showContractModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" size="md" @click="saveContract">{{ isEditingContract ? 'บันทึกการแก้ไข' : 'ยืนยันออกสัญญาเช่า' }}</AppButton>
        </template>
      </AppModal>



      <transition name="modal-bounce">
        <div
          v-if="showInspectPhotoModal"
          class="unsaved-leave-backdrop inspect-photo-backdrop"
          role="dialog"
          aria-modal="true"
          @click.self="showInspectPhotoModal = false"
        >
          <div class="inspect-photo-modal">
            <AppModalHeader :title="`ภาพตัวอย่างแผง ${stallDisplayNumber}`" :subtitle="`${selectedProject?.name || stallEdit.project_name} · ${formatFloorDisplay(stallEdit.floor)}`" icon="document" @close="showInspectPhotoModal = false" />
            <div class="ipm-stage">
              <img :src="inspectStallPhotos[inspectPhotoIndex]" :alt="`ภาพแผง ${stallDisplayNumber}`" class="ipm-main-img" />
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
import AppImageUpload, { type AppImageFile } from '~/components/AppImageUpload.vue'
import AppInput from '~/components/AppInput.vue'
import AppModal from '~/components/AppModal.vue'
import AppPagination from '~/components/AppPagination.vue'
import AppSegmentedControl from '~/components/AppSegmentedControl.vue'
import AppSelect from '~/components/AppSelect.vue'
import { type MeterProjectRow, type MeterStallRow } from './components/StallMeterSummaryModal.vue'
import MaintenanceActionModal, { type MaintenanceActionMode, type MaintenanceActionPayload } from './components/MaintenanceActionModal.vue'
import { chargeRuleOf, needsReading, registerStallProjects, useStallBilling, type CostFilter } from '../composables/useStallBilling'
import { useUtilitySettings } from '../composables/useUtilitySettings'
import { getMarketGisConfig } from '../data/markets_registry'
import { mergeMarketProjects, clearMarketGisCache, cloneProjects } from '../data/market_loader'
import { getAllFloorLayoutsForProject, isOverviewBoundaryProject } from '../data/market_stall_layout'
import type { ProjectZone, MarketFacility } from '../data/market_gis_types'

const StallUtilityBillingPanel = defineAsyncComponent(() => import('./components/StallUtilityBillingPanel.vue'))
const StallTempBillPanel = defineAsyncComponent(() => import('./components/StallTempBillPanel.vue'))
const StallCostDrawer = defineAsyncComponent(() => import('./components/StallCostDrawer.vue'))
const StallMeterSummaryPanel = defineAsyncComponent(() => import('./components/StallMeterSummaryPanel.vue'))
const StallMeterProjectCard = defineAsyncComponent(() => import('./components/StallMeterProjectCard.vue'))
const StallMeterRecordCard = defineAsyncComponent(() => import('./components/StallMeterRecordCard.vue'))

const route = useRoute()
const router = useRouter()
const { get, post, put, patch, del, apiFetch } = useApi()
const { showSuccess, showError, showConfirm, showWarning } = useAppAlert()
const { role } = useAuth()
const { applyMarketOverride } = useMarketOverrides()
const isStaff = computed(() => role.value === 'staff')
const mapRef = ref<any>(null)
const marketId = computed(() => String(route.query.market || '1'))
const isMapOverviewMode = computed(() => route.path === '/stalls/map' || route.query.mode === 'map')
const isInspectMode = computed(() => isMapOverviewMode.value || String(route.query.mode || '') === 'inspect')
export type MapPageViewMode = 'stalls' | 'payments' | 'meters'
const pageMode = ref<MapPageViewMode>('stalls')
// มุมมองที่ทำงานกับค่าใช้จ่าย (จดมิเตอร์ / การชำระเงิน): กดแผงบนแผนที่แล้วเปิดฟอร์มค่าน้ำ-ไฟ/บิล ไม่เปิดรายละเอียดแผงปกติ
const isCostMode = computed(() => pageMode.value === 'meters' || pageMode.value === 'payments')
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
    svgIcon: '<path d="m2 7 4.41-4.41A2 2 0 0 1 7.83 2h8.34a2 2 0 0 1 1.42.59L22 7"></path><path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"></path><path d="M15 22v-4a2 2 0 0 0-2-2h-2a2 2 0 0 0-2 2v4"></path><path d="M2 7h20"></path><path d="M22 7v3a2 2 0 0 1-2 2 2.7 2.7 0 0 1-1.59-.63.7.7 0 0 0-.82 0A2.7 2.7 0 0 1 16 12a2.7 2.7 0 0 1-1.59-.63.7.7 0 0 0-.82 0A2.7 2.7 0 0 1 12 12a2.7 2.7 0 0 1-1.59-.63.7.7 0 0 0-.82 0A2.7 2.7 0 0 1 8 12a2.7 2.7 0 0 1-1.59-.63.7.7 0 0 0-.82 0A2.7 2.7 0 0 1 4 12a2 2 0 0 1-2-2V7"></path>'
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
  const modeChanged = pageMode.value !== mode
  pageMode.value = mode
  isViewDropdownOpen.value = false
  selectedPaymentFilter.value = null
  selectedMeterFilter.value = null
  dashboardStatusFilter.value = 'all'
  // สลับมุมมองแล้วเริ่มจากภาพรวมตลาด ไม่พาโครงการที่เปิดค้างไว้ไปด้วย ให้เลือกโครงการใหม่เอง
  if (modeChanged && selectedProject.value) clearSelectedProject()
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
const marketDebts = ref<any[]>([])
const extractStallCodesFromRecord = (record: any) => {
  const codes = new Set<string>()
  String(record?.stall_ident || record?.stall_number || record?.stall_code || record?.ident || record?.code || '')
    .split(',')
    .map(code => code.trim().toUpperCase())
    .filter(Boolean)
    .forEach(code => codes.add(code))
  const descriptionCodes = String(record?.description || '').toUpperCase()
    .match(/[A-Z]{1,10}-(?:\d+[A-Z]?)(?:-\d+[A-Z]?)+/g) || []
  descriptionCodes.forEach(code => codes.add(code))
  return Array.from(codes)
}
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

// Keep the KPI ribbon and the project drawer on the same authoritative GIS
// dataset. Some registry entries only contain market metadata while the map
// receives the complete projects from the backend.
const handleMapDataLoaded = (payload: { projects?: ProjectZone[]; facilities?: MarketFacility[] }) => {
  if (Array.isArray(payload.projects) && payload.projects.length) {
    projects.value = payload.projects
  }
  if (Array.isArray(payload.facilities)) {
    facilities.value = payload.facilities
  }
}

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

// ตัวเลขสรุปการจดมิเตอร์: เปิดโครงการอยู่นับเฉพาะโครงการนั้น ไม่ได้เปิดนับทั้งตลาด
const meterScopeRows = computed<MeterStallRow[]>(() => (pageMode.value === 'meters' && selectedProject.value ? meterProjectStalls.value : meterLiveRows.value))
const meterDatasets = computed(() => {
  const target = meterScopeRows.value
  const recordedStalls = target.filter(row => row.recorded)
  const pendingStalls = target.filter(row => row.state === 'pending')
  const overdueStalls = pendingStalls.filter(row => billing.getUtility(row.code)?.status === 'overdue')
  const anomalyStalls = target.filter(row => row.abnormal)

  return {
    totalStalls: target,
    recordedStalls,
    pendingStalls,
    overdueStalls,
    anomalyStalls
  }
})

const meterStallCodeMap = computed(() => ({
  total: meterDatasets.value.totalStalls.map(s => s.code.toUpperCase()),
  recorded: meterDatasets.value.recordedStalls.map(s => s.code.toUpperCase()),
  pending: meterDatasets.value.pendingStalls.map(s => s.code.toUpperCase()),
  overdue: meterDatasets.value.overdueStalls.map(s => s.code.toUpperCase()),
  anomaly: meterDatasets.value.anomalyStalls.map(s => s.code.toUpperCase())
}))

const meterSummary = computed(() => {
  const rows = meterDatasets.value.totalStalls
  const dueDate = new Date(`${billing.billingMonth.value}-01T00:00:00`)
  dueDate.setMonth(dueDate.getMonth() + 1, 0)
  const dueDateText = new Intl.DateTimeFormat('th-TH', { day: 'numeric', month: 'short' }).format(dueDate)
  const totalCount = rows.length
  const recordedCount = meterDatasets.value.recordedStalls.length
  return {
    totalCount,
    waterCount: rows.filter(row => needsReading(row.code, 'water')).length,
    electricCount: rows.filter(row => needsReading(row.code, 'electric')).length,
    recordedCount,
    completionRate: totalCount ? Number(((recordedCount / totalCount) * 100).toFixed(1)) : 0,
    pendingCount: meterDatasets.value.pendingStalls.length,
    dueDateText,
    overdueCount: meterDatasets.value.overdueStalls.length,
    anomalyCount: meterDatasets.value.anomalyStalls.length
  }
})

// โหมดจดมิเตอร์: ระบายสีแผนที่ตามสถานะการจด (แผงที่ไม่มีบิลไม่อยู่ในรายการ จึงจางลง)
const meterMapStates = computed<Record<string, string> | null>(() => {
  if (pageMode.value !== 'meters') return null
  const states: Record<string, string> = {}
  meterLiveRows.value.forEach((row) => { states[row.code.toUpperCase()] = row.state })
  return states
})

// สรุปแผงทั้งโครงการที่เปิดอยู่ในการ์ดจดมิเตอร์ (บิลมีเฉพาะแผงที่มีผู้เช่า จึงใช้แยกแผงว่างกับแผงเหมาจ่ายได้)
const meterProjectStats = computed(() => {
  if (pageMode.value !== 'meters' || !selectedProject.value) return null
  const codes = new Set<string>()
  getAllFloorLayoutsForProject(selectedProject.value).forEach((layout) => {
    layout.stalls.forEach((stall: any) => {
      const code = stallCodeOf(stall)
      if (code) codes.add(code)
    })
  })
  const withTenant = billing.utilityList.value.filter(bill => codes.has(bill.stall.toUpperCase())).length
  return { total: codes.size, vacant: codes.size - withTenant, flat: Math.max(0, withTenant - meterProjectStalls.value.length) }
})
// งวดที่เลือกเป็นงวดย้อนหลัง (เก่ากว่าเดือนปัจจุบัน) ดูข้อมูลที่จดไว้ได้อย่างเดียว
const currentMonthKey = () => { const d = new Date(); return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}` }
const meterCurrentMonth = computed(() => currentMonthKey())
const isPastPeriod = computed(() => billing.billingMonth.value < currentMonthKey())
// รอบจดปัจจุบันและวันครบกำหนด (สิ้นเดือนของรอบบิล)
const meterCycle = computed(() => {
  const start = new Date(`${billing.billingMonth.value}-01T00:00:00`)
  return {
    label: new Intl.DateTimeFormat('th-TH', { month: 'long', year: 'numeric' }).format(start),
    due: meterSummary.value.dueDateText
  }
})

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

// แผงยอดนิยม: จัดอันดับการจองในช่วง 30 วันล่าสุด (นับการจองที่ยังไม่ถูกยกเลิก/ปฏิเสธ) สลับดูได้ 2 มุมมอง: ตามแผง / ตามโครงการ
// แต่ละมุมมองมีช่วงเวลาของตัวเอง (แผงกับโครงการไม่ใช้ช่วงเดียวกัน)
type PopularDays = 7 | 30 | 90
const stallPopularDays = ref<PopularDays>(30)
const projectPopularDays = ref<PopularDays>(30)
// ช่วงเวลาของมุมมองที่เปิดอยู่ (ปุ่มเลือกช่วงอ่านและเขียนค่าของมุมมองนั้นเท่านั้น)
const popularDays = computed<PopularDays>({
  get: () => (popularMode.value === 'stall' ? stallPopularDays.value : projectPopularDays.value),
  set: (value) => {
    if (popularMode.value === 'stall') stallPopularDays.value = value
    else projectPopularDays.value = value
  }
})
const popularDayOptions: Array<7 | 30 | 90> = [7, 30, 90]
const popularMode = ref<'stall' | 'project'>('stall')
const popularModeOptions = [
  { value: 'stall', label: 'แผง' },
  { value: 'project', label: 'โครงการ' }
]
const recentActiveBookings = (days: PopularDays) => {
  const since = Date.now() - days * 86_400_000
  const excluded = new Set(['cancelled', 'rejected', 'not_interested', 'no_show', 'terminated'])
  return marketBookings.value.filter((booking: any) => {
    if (excluded.has(String(booking.booking_status || booking.status || '').toLowerCase())) return false
    const time = new Date(booking.created_at || booking.visit_date || 0).getTime()
    return !(Number.isFinite(time) && time > 0 && time < since)
  })
}
const bookingTime = (booking: any) => {
  const time = new Date(booking.created_at || booking.visit_date || 0).getTime()
  return Number.isFinite(time) ? time : 0
}
// เลขอันดับและเหรียญ 1-3 ตามลำดับแถวเสมอ (แถวที่เท่ากันเรียงตามการจองล่าสุด)
const withTies = <T extends { count: number }>(rows: T[]) => rows.map((row, index) => ({ ...row, rank: index + 1, hasMedal: index < 3 }))
const popularStallRanking = computed(() => {
  const counts = new Map<string, { code: string; count: number; projectName: string; last: number }>()
  recentActiveBookings(stallPopularDays.value).filter(recordMatchesDashboardProject).forEach((booking: any) => {
    getDashboardRecordCodes(booking).forEach((code) => {
      const item = counts.get(code) || { code, count: 0, projectName: String(booking.project_name || ''), last: 0 }
      item.count += 1
      item.last = Math.max(item.last, bookingTime(booking))
      counts.set(code, item)
    })
  })
  const top = [...counts.values()].sort((a, b) => b.count - a.count || b.last - a.last).slice(0, 5)
  if (!top.length) return []
  // ชื่อแผงเรียกตามเลขจากผังสำรวจ (เช่น "แผง 129") และชื่อโครงการจากผัง เหมือนหัวแผงรายละเอียด ไม่ใช้รหัสภายในเต็ม
  const wanted = new Set(top.map(item => item.code))
  const info = new Map<string, { number: string; projectName: string }>()
  visibleProjects.value.forEach((project) => {
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim().toUpperCase()
        if (!wanted.has(code) || info.has(code)) return
        const number = String(stall.sourceProperties?.number ?? stall.sourceProperties?.Number ?? '').trim()
        info.set(code, { number, projectName: cleanProjectName(project.name) })
      })
    })
  })
  return withTies(top.map(item => ({
    ...item,
    number: info.get(item.code)?.number || item.code,
    projectName: info.get(item.code)?.projectName || cleanProjectName(item.projectName)
  })))
})
const popularStallAllTied = computed(() => popularStallRanking.value.length > 1 && popularStallRanking.value.every(item => item.count === popularStallRanking.value[0]!.count))

// มุมมองตามโครงการ: โครงการที่ถูกจองมากที่สุด (นับทุกโครงการ ไม่ผูกกับโครงการที่เลือกในแดชบอร์ด)
const popularProjectRanking = computed(() => {
  const bookings = recentActiveBookings(projectPopularDays.value)
  if (!bookings.length) return []
  const projectOfCode = new Map<string, { id: string; name: string }>()
  const stallTotals = new Map<string, number>()
  visibleProjects.value.forEach((project) => {
    let total = 0
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim().toUpperCase()
        total += 1
        if (code && !projectOfCode.has(code)) projectOfCode.set(code, { id: String(project.id), name: cleanProjectName(project.name) })
      })
    })
    stallTotals.set(String(project.id), total)
  })
  const groups = new Map<string, { id: string; name: string; count: number; stalls: Set<string>; last: number }>()
  bookings.forEach((booking: any) => {
    getDashboardRecordCodes(booking).forEach((code) => {
      const project = projectOfCode.get(code)
      if (!project) return
      const group = groups.get(project.id) || { id: project.id, name: project.name, count: 0, stalls: new Set<string>(), last: 0 }
      group.count += 1
      group.stalls.add(code)
      group.last = Math.max(group.last, bookingTime(booking))
      groups.set(project.id, group)
    })
  })
  // เรียงตามสัดส่วนแผงที่ถูกจองเทียบจำนวนแผงทั้งโครงการ (โครงการใหญ่จึงไม่ได้เปรียบเพราะมีแผงเยอะ) แล้วตามจำนวนครั้ง
  const rows = [...groups.values()].map((group) => {
    const total = stallTotals.get(group.id) || 0
    const percent = total > 0 ? (group.stalls.size / total) * 100 : 0
    return { id: group.id, name: group.name, count: group.count, stallCount: group.stalls.size, total, percent, last: group.last }
  })
  rows.sort((a, b) => b.percent - a.percent || b.count - a.count || b.last - a.last)
  const top = rows.slice(0, 5)
  return top.map((row, index) => ({ ...row, rank: index + 1, hasMedal: index < 3 }))
})
const popularProjectAllTied = computed(() => popularProjectRanking.value.length > 1 && popularProjectRanking.value.every(item => item.percent === popularProjectRanking.value[0]!.percent && item.count === popularProjectRanking.value[0]!.count))

// การ์ดแผงยอดนิยม: แต่ละมุมมอง (แผง / โครงการ) เก็บการเลือกของตัวเองแยกจากแผนที่
// ไฮไลต์เกิดจากการกดแถวในการ์ดนี้เท่านั้น ไม่ดูจากช่องค้นหาหรือโครงการที่เปิดบนแผนที่
const popularStallPick = ref('')
const popularProjectPick = ref('')
watch(popularMode, () => {
  popularStallPick.value = ''
  popularProjectPick.value = ''
})
const isPopularStallActive = (code: string) => popularMode.value === 'stall' && !!popularStallPick.value && popularStallPick.value === String(code || '').trim().toUpperCase()
const isPopularProjectActive = (id: string) => popularMode.value === 'project' && !!popularProjectPick.value && popularProjectPick.value === String(id || '')
const openPopularProject = (id: string) => {
  const project = visibleProjects.value.find(item => item.id === id)
  if (!project) return
  if (isPopularProjectActive(id)) {
    popularProjectPick.value = ''
    mapRef.value?.exitProject?.()
    return
  }
  popularProjectPick.value = String(id)
  if (typeof mapRef.value?.enterProject === 'function') mapRef.value.enterProject(project)
  else selectProject(project)
}
const filterStallFromDashboard = async (code: string) => {
  const clean = String(code || '').trim().toUpperCase()
  if (isPopularStallActive(clean)) {
    popularStallPick.value = ''
    mapRef.value?.clearAppliedSearch?.()
    return
  }
  popularStallPick.value = clean
  await mapRef.value?.searchStallFromDashboard?.(clean)
}

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
    totalMonthlyRent += Number(contract.monthly_rental || contract.monthly_rent || contract.rent_amount || contract.rent || 0)
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
const selectedProject = ref<ProjectZone | null>(null)
const contracts = ref<any[]>([])
const stallHistory = ref<any[]>([])
const gisPublished = ref(false)
const isManageDataLoading = ref(true)
const isActionDataLoading = ref(true)
let loadDataController: AbortController | null = null
let loadDataRun = 0
let reloadAfterActivation = false
let isManagePageActive = true
let loadedMarketId = ''
let hasActivatedManagePage = false
const savingGIS = ref(false)
const savingStall = ref(false)
const showStallEditor = ref(false)
const showInspectStallOverlay = ref(false)
watch(showInspectStallOverlay, (open) => {
  if (open) isViewDropdownOpen.value = false
})
const showInspectPhotoModal = ref(false)
const emptyStall = () => ({ id: 0, ident: '', name: '', description: '', category: '', zone_id: 1, zone_name: selectedProject.value?.name || '', project_name: selectedProject.value?.name || '', floor: '1', width: 2, depth: 2, daily_rent: 150, monthly_rent: 1800, stall_status: 'available', bookable: true, coordinates_geojson: '', published: false })
const stallEdit = ref(emptyStall())
const selectedStallSurveyNumber = ref('')
const stallDisplayNumber = computed(() => selectedStallSurveyNumber.value || String(stallEdit.value.ident || ''))
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
const editablePhotoFiles = computed<AppImageFile[]>({
  get: () => editablePhotos.value.map((dataUrl, index) => ({ name: `stall-photo-${index + 1}`, dataUrl })),
  set: (photos) => {
    ;(stallEdit.value as any).photos = photos.map(photo => photo.dataUrl).filter(Boolean)
    persistStallPhotos()
  }
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
  if (isMapOverviewMode.value && isCostMode.value) {
    if (meterCardCode.value) return [meterCardCode.value]
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
  acceptedAt?: string
  assigneeType?: 'market_technician' | 'contractor' | 'tenant'
  assigneeName?: string
  expensePayer?: 'market' | 'tenant' | 'shared'
  priority?: 'normal' | 'urgent' | 'critical'
  tenantMessage?: string
  completedAt?: string
  completionDate?: string
  executorType?: 'market_technician' | 'contractor' | 'tenant'
  executorName?: string
  actualCost?: number
  repairSummary?: string
  beforePhotos?: Array<{ name: string; dataUrl?: string }>
  afterPhotos?: Array<{ name: string; dataUrl?: string }>
  rejectionReason?: string
  rejectionDetail?: string
  rejectedAt?: string
  rejectedBy?: string
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
const isMiniActionAlertExpanded = ref(false)
const contractSubTab = ref<'active' | 'history'>('active')
const bookingTypeFilter = ref<'all' | 'booking' | 'application'>('all')
const showZoneEditor = ref(false)
const showFacilityEditor = ref(false)
const facilityNameError = ref('')
const mapVersion = ref(0)
const touchManagedMapData = () => { mapVersion.value += 1 }
const zoneOptions = [{ value: 1, label: 'โซน A' }, { value: 2, label: 'โซน B' }, { value: 3, label: 'โซน C' }, { value: 4, label: 'โซน D' }]
const floorOptions = computed(() => {
  const fls = selectedProject.value?.floors
  if (fls && fls.length > 0) return fls.map(f => ({ value: String(f.id), label: f.name || `ชั้น ${f.id}` }))
  return [{ value: '1', label: 'ชั้น 1' }, { value: '2', label: 'ชั้น 2' }, { value: '3', label: 'ชั้น 3' }]
})
const stallBookableOptions = [{ value: 'open', label: 'เปิดรับจอง' }, { value: 'paused', label: 'ปิดรับจอง' }]
const stallBookableModel = computed({
  get: () => (stallEdit.value.bookable === false ? 'paused' : 'open'),
  set: (value: string | number) => { stallEdit.value.bookable = value !== 'paused' }
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

const getBookingCountdownText = (dateStr?: string | null) => {
  if (!dateStr) return ''
  const target = new Date(dateStr)
  if (isNaN(target.getTime())) return ''
  const now = new Date()
  const targetDay = new Date(target.getFullYear(), target.getMonth(), target.getDate()).getTime()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  const days = Math.round((targetDay - today) / 86400000)
  if (days < 0) return ''
  if (days === 0) return 'วันนี้'
  return `อีก ${days} วัน`
}

const getApplicationWaitingText = (dateStr?: string | null) => {
  if (!dateStr) return ''
  const submitted = new Date(dateStr)
  if (isNaN(submitted.getTime())) return ''
  const days = Math.max(0, Math.floor((Date.now() - submitted.getTime()) / 86400000))
  return days === 0 ? 'ยื่นวันนี้' : `รอ ${days} วัน`
}

const getContractRemainingDays = (endDateStr?: string | null) => {
  if (!endDateStr) return null
  const end = new Date(endDateStr)
  if (isNaN(end.getTime())) return null
  const now = new Date()
  const endDay = new Date(end.getFullYear(), end.getMonth(), end.getDate()).getTime()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  return Math.max(0, Math.ceil((endDay - today) / 86400000))
}

const getContractProgressPercent = (startDateStr?: string | null, endDateStr?: string | null) => {
  if (!startDateStr || !endDateStr) return 0
  const start = new Date(startDateStr).getTime()
  const end = new Date(endDateStr).getTime()
  if (!Number.isFinite(start) || !Number.isFinite(end) || end <= start) return 0
  return Math.min(100, Math.max(0, ((Date.now() - start) / (end - start)) * 100))
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
    } else if (st === 'pending_slot') {
      list.push({
        key: 'act-booking-slot',
        category: 'การจอง',
        severity: 'warning',
        title: 'รอเลือกวันนัดชม',
        count: 1,
        description: `ผู้จอง: ${applicant} ยังไม่ได้เลือกวันและเวลานัดเข้าชมแผง`,
        actionLabel: 'ไปที่การจอง',
        targetTab: 'bookings'
      })
    } else if (st === 'approved') {
      list.push({
        key: 'act-booking-approved',
        category: 'การจอง',
        severity: 'warning',
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
        severity: 'warning',
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
    if (['approved', 'director_approved'].includes(st) && !currentStallContract.value) {
      list.push({
        key: 'act-application-contract',
        category: 'คำร้อง',
        severity: 'warning',
        title: 'รอออกสัญญา',
        count: 1,
        description: `คำร้อง ${a.application_number || a.application_no || '-'} อนุมัติแล้วและยังไม่มีสัญญาเช่า`,
        actionLabel: 'ไปที่คำร้อง',
        targetTab: 'bookings'
      })
    } else if (st !== 'approved' && st !== 'director_approved' && st !== 'completed') {
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

const stallActionTotalCount = computed(() => stallActionRichItems.value.reduce(
  (total, item) => total + (item.count || 1),
  0
))

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
  // สัญญาในหน้าข้อมูลแผงต้องมาจาก backend เท่านั้น
  // ข้อมูลจำลอง/ฉบับร่างใน localStorage เคยทำให้สัญญาปัจจุบันและประวัติถูกนับซ้ำ
  delete customStallContracts.value[upper]
}

const loadAllCustomContracts = () => {
  // คงฟังก์ชันไว้สำหรับจุดเรียกเดิม แต่ไม่โหลดสัญญาจำลองกลับมาปนกับข้อมูลจริง
  customStallContracts.value = {}
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
    if (foundIdx >= 0 && (
      String(stalls.value[foundIdx]?.stall_status || stalls.value[foundIdx]?.status) !== 'occupied'
      || stalls.value[foundIdx]?.name !== stallEdit.value.name
      || Number(stalls.value[foundIdx]?.monthly_rent) !== Number(stallEdit.value.monthly_rent)
    )) {
      stalls.value[foundIdx] = {
        ...stalls.value[foundIdx],
        stall_status: 'occupied',
        status: 'occupied',
        name: stallEdit.value.name,
        monthly_rent: stallEdit.value.monthly_rent
      }
    }

    // แผงมีรหัสผูกกับโครงการที่กำลังเปิดอยู่ จึงไม่ต้องวนทุกโครงการ/ทุกแผงในตลาด
    const projectsToSync = selectedProject.value ? [selectedProject.value] : []
    projectsToSync.forEach(p => {
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
    if (foundIdx >= 0 && String(stalls.value[foundIdx]?.stall_status || stalls.value[foundIdx]?.status) !== targetStatus) {
      stalls.value[foundIdx] = {
        ...stalls.value[foundIdx],
        stall_status: targetStatus,
        status: targetStatus,
        name: targetStatus === 'available' ? '' : stalls.value[foundIdx].name
      }
    }

    const projectsToSync = selectedProject.value ? [selectedProject.value] : []
    projectsToSync.forEach(p => {
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
const maintenanceByStall = shallowRef<Record<string, StallMaintenanceItem[]>>({})
const maintenanceIndexVersion = ref(0)
const showAddMaintenanceModal = ref(false)
const activeMaintenanceAction = ref<{ id: string; mode: MaintenanceActionMode } | null>(null)
const activeMaintenanceItem = computed(() => {
  if (!activeMaintenanceAction.value) return null
  return stallMaintenance.value.find(item => item.id === activeMaintenanceAction.value?.id) || null
})
const maintenanceRejectReasonOptions = [
  { value: 'tenant_responsibility', label: 'เป็นความรับผิดชอบของผู้เช่าตามสัญญา' },
  { value: 'insufficient_information', label: 'ข้อมูลหรือหลักฐานไม่เพียงพอ' },
  { value: 'duplicate', label: 'เป็นรายการซ้ำ' },
  { value: 'outside_scope', label: 'อยู่นอกขอบเขตบริการของตลาด' },
  { value: 'other', label: 'เหตุผลอื่น' }
]
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

const replaceMaintenanceIndex = (records: StallMaintenanceItem[]) => {
  const next: Record<string, StallMaintenanceItem[]> = {}
  records.forEach((record) => {
    const code = String(record.stallIdent || (record as any).stall_ident || '').trim().toUpperCase()
    if (!code) return
    ;(next[code] ||= []).push({ ...record, stallIdent: code, id: String(record.id) })
  })
  maintenanceByStall.value = next
  maintenanceIndexVersion.value += 1
}

const syncCurrentMaintenanceIndex = () => {
  const code = String(stallEdit.value.ident || '').trim().toUpperCase()
  if (!code) return
  maintenanceByStall.value = { ...maintenanceByStall.value, [code]: [...stallMaintenance.value] }
  maintenanceIndexVersion.value += 1
}

const maintenancePayload = (item: StallMaintenanceItem) => ({
  ...item,
  market_id: resolvedMarketID.value,
  stall_ident: item.stallIdent
})

const persistMaintenanceItem = async (item: StallMaintenanceItem) => {
  const numericID = Number(item.id)
  const response = Number.isFinite(numericID) && numericID > 0
    ? await put(`/maintenance-records/${numericID}`, maintenancePayload(item))
    : await post('/maintenance-records', maintenancePayload(item))
  return unwrap(response) as StallMaintenanceItem
}

const formatMaintenanceDate = (value?: string) => {
  if (!value) return '–'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('th-TH', { day: 'numeric', month: 'short', year: 'numeric' }).format(date)
}

const maintenanceWaitingDays = (item: StallMaintenanceItem) => {
  const started = new Date(item.requestDate || item.createdAt)
  if (Number.isNaN(started.getTime())) return 0
  return Math.max(0, Math.ceil((Date.now() - started.getTime()) / 86400000))
}

const maintenanceDurationDays = (item: StallMaintenanceItem) => {
  const started = new Date(item.acceptedAt || item.requestDate || item.createdAt)
  const ended = new Date(item.completedAt || item.completionDate || '')
  if (Number.isNaN(started.getTime()) || Number.isNaN(ended.getTime())) return 0
  return Math.max(1, Math.ceil((ended.getTime() - started.getTime()) / 86400000))
}

const formatMaintenanceAssignee = (value?: StallMaintenanceItem['assigneeType']) => ({
  market_technician: 'ช่างสำนักงานตลาด',
  contractor: 'ผู้รับเหมา',
  tenant: 'ผู้เช่าดำเนินการเอง'
} as Record<string, string>)[value || ''] || '–'

const formatMaintenancePayer = (value?: StallMaintenanceItem['expensePayer']) => ({
  market: 'สำนักงานตลาด',
  tenant: 'ผู้เช่า',
  shared: 'ร่วมรับผิดชอบ'
} as Record<string, string>)[value || ''] || '–'

const formatMaintenancePriority = (value?: StallMaintenanceItem['priority']) => ({
  normal: 'ปกติ',
  urgent: 'เร่งด่วน',
  critical: 'ฉุกเฉิน'
} as Record<string, string>)[value || ''] || 'ปกติ'

const formatMaintenanceRejectReason = (value?: string) => maintenanceRejectReasonOptions.find(option => option.value === value)?.label || value || '–'

const openMaintenanceAction = (item: StallMaintenanceItem, mode: MaintenanceActionMode) => {
  activeMaintenanceAction.value = { id: item.id, mode }
}

const closeMaintenanceAction = () => { activeMaintenanceAction.value = null }
const handleMaintenanceModalVisibility = (visible: boolean) => { if (!visible) closeMaintenanceAction() }

const submitMaintenanceAction = async (payload: MaintenanceActionPayload) => {
  const item = activeMaintenanceItem.value
  if (!item) return
  const draft: StallMaintenanceItem = { ...item }

  if (payload.mode === 'accept') {
    draft.status = 'in_progress'
    draft.assigneeType = payload.assigneeType
    draft.assigneeName = formatMaintenanceAssignee(payload.assigneeType)
    draft.scheduledDate = payload.scheduledDate
    draft.expensePayer = payload.expensePayer
    draft.priority = payload.priority
    draft.tenantMessage = payload.tenantMessage.trim() || undefined
    draft.acceptedAt = new Date().toISOString()
  } else if (payload.mode === 'complete') {
    draft.status = 'completed'
    draft.executorType = payload.executorType
    draft.executorName = payload.executorName.trim()
    draft.completionDate = payload.completionDate || new Date().toISOString().slice(0, 10)
    draft.completedAt = new Date().toISOString()
    draft.actualCost = payload.actualCost === '' ? undefined : Number(payload.actualCost)
    draft.repairSummary = payload.repairSummary.trim() || undefined
    draft.afterPhotos = [...payload.afterPhotos]
  } else {
    draft.status = 'rejected'
    draft.rejectionReason = payload.reason
    draft.rejectionDetail = payload.detail.trim()
    draft.rejectedAt = new Date().toISOString()
    draft.rejectedBy = 'เจ้าหน้าที่ตลาด'
  }

  try {
    const saved = await persistMaintenanceItem(draft)
    Object.assign(item, saved, { id: String(saved.id || item.id) })
    syncCurrentMaintenanceIndex()
    closeMaintenanceAction()
  } catch {
    void showError('ไม่สามารถบันทึกการดำเนินงานซ่อมลงระบบได้', 'บันทึกไม่สำเร็จ')
    return
  }
  const messages = {
    accept: ['รับเรื่องและมอบหมายงานซ่อมเรียบร้อยแล้ว', 'รับเรื่องสำเร็จ'],
    complete: ['บันทึกผลการซ่อมและปิดงานเรียบร้อยแล้ว', 'เสร็จสิ้น'],
    reject: ['บันทึกการไม่รับเรื่องและข้อความถึงผู้เช่าแล้ว', 'บันทึกสำเร็จ']
  } as const
  void showSuccess(messages[payload.mode][0], messages[payload.mode][1])
}

const loadStallMaintenance = (code: string) => {
  if (!code) {
    stallMaintenance.value = []
    return
  }
  stallMaintenance.value = [...(maintenanceByStall.value[code.toUpperCase()] || [])]
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

const saveMaintenance = async () => {
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

  try {
    const saved = await persistMaintenanceItem(newItem)
    stallMaintenance.value.unshift({ ...saved, id: String(saved.id), stallIdent: code })
    syncCurrentMaintenanceIndex()
  } catch {
    void showError('ไม่สามารถบันทึกรายการซ่อมลงระบบได้', 'บันทึกไม่สำเร็จ')
    return
  }
  showAddMaintenanceModal.value = false
  void showSuccess('บันทึกรายการซ่อม / ขออนุญาตเรียบร้อยแล้ว', 'บันทึกสำเร็จ')
}

const updateMaintenanceStatus = async (id: string, newStatus: 'completed' | 'in_progress' | 'pending' | 'rejected') => {
  const item = stallMaintenance.value.find(m => m.id === id)
  if (!item) return
  const draft = { ...item, status: newStatus }
  try {
    const saved = await persistMaintenanceItem(draft)
    Object.assign(item, saved, { id: String(saved.id || item.id) })
    syncCurrentMaintenanceIndex()
  } catch {
    void showError('ไม่สามารถอัปเดตสถานะงานซ่อมได้', 'อัปเดตไม่สำเร็จ')
    return
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
  try {
    await del(`/maintenance-records/${encodeURIComponent(id)}`)
    stallMaintenance.value = stallMaintenance.value.filter(m => m.id !== id)
    syncCurrentMaintenanceIndex()
  } catch {
    void showError('ไม่สามารถลบรายการซ่อมได้', 'ลบไม่สำเร็จ')
    return
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
    rejected: 'ไม่รับเรื่อง'
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
  const upperCode = code.trim().toUpperCase()
  stallPayments.value = marketDebts.value
    .filter(debt => extractStallCodesFromRecord(debt).includes(upperCode))
    .map((debt: any) => {
      const rawStatus = String(debt.status || '').trim().toLowerCase()
      const isOutstanding = ['outstanding', 'overdue', 'past_due', 'ค้างชำระ'].includes(rawStatus)
      return {
        id: `debt-${debt.id}`,
        stallIdent: upperCode,
        invoiceNumber: `DEBT-${debt.id}`,
        periodMonth: 'รายการหนี้',
        title: String(debt.description || 'รายการค้างชำระ'),
        amount: Number(debt.amount_satang || 0) / 100,
        paymentMethod: 'transfer' as const,
        status: isOutstanding ? 'overdue' as const : 'paid' as const,
        createdAt: String(debt.created_at || ''),
        paidDate: isOutstanding ? undefined : String(debt.updated_at || debt.created_at || '')
      }
    })
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
    rawDate?: string
    rawStatus?: string
    sourceRef?: string
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
      rawDate: b.visit_date || b.created_at || '',
      rawStatus: st,
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
      rawDate: a.created_at || '',
      rawStatus: st,
      sourceRef: a.booking_number || a.source_booking_number || a.booking?.booking_number || '',
      dateLabel: 'วันที่ยื่นคำร้อง',
      dateText: formatStallDate(a.created_at),
      statusText: isApproved ? 'อนุมัติแล้ว' : (isRejected ? 'ไม่อนุมัติ' : 'รอพิจารณาอนุมัติ'),
      badgeClass: isApproved ? 'is-available' : (isRejected ? 'is-cancelled' : 'is-reserved'),
      targetUrl: a.id ? `/applications/${a.id}` : `/applications?search=${encodeURIComponent(a.application_number || target)}`
    })
  })

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
  const contractMap = new Map<string, any>()
  ;[...marketContracts.value, ...contracts.value].forEach((c: any) => {
    const code = String(c.stall_ident || c.stall_number || c.stall_code || '').toUpperCase()
    if (!code.split(',').map(s => s.trim()).includes(target) && code !== target) return
    const key = String(c.contract_number || c.contract_no || c.id || '')
    if (!key) return
    contractMap.set(key, c)
  })
  return Array.from(contractMap.values())
})

const activeStallContracts = computed(() =>
  stallContractsList.value.filter((c: any) => {
    const st = String(c.status || c.contract_state || '').toLowerCase()
    return st === 'active' || st === 'ใช้งานอยู่'
  })
)

const historyStallContracts = computed(() =>
  stallContractsList.value.filter((c: any) => {
    const st = String(c.status || c.contract_state || '').toLowerCase()
    return st !== 'active' && st !== 'ใช้งานอยู่'
  })
)

// ติดตามสัญญา (ต่อสัญญา / ส่งคืนแผง / คดี): บันทึกลงฐานข้อมูลผ่าน PATCH /contracts/:id/workflow แล้วแดชบอร์ดและตัวกรองอ่านค่าจริงจากสัญญา
type ContractWorkflowKey = 'renewal_status' | 'return_status' | 'legal_status'
const NO_WORKFLOW = 'none'
const contractWorkflowFields: Array<{ key: ContractWorkflowKey; label: string; options: Array<{ value: string; label: string }> }> = [
  { key: 'renewal_status', label: 'การต่อสัญญา', options: [
    { value: NO_WORKFLOW, label: 'ไม่มี' }, { value: 'pending', label: 'รอต่ออายุ' }, { value: 'reviewing', label: 'อยู่ระหว่างพิจารณา' },
    { value: 'approved', label: 'ต่อสัญญาแล้ว' }, { value: 'declined', label: 'ไม่ประสงค์ต่อ' }
  ] },
  { key: 'return_status', label: 'การส่งคืนแผง', options: [
    { value: NO_WORKFLOW, label: 'ไม่มี' }, { value: 'notified', label: 'แจ้งไม่ต่อ' }, { value: 'pending_return', label: 'รอส่งคืน' },
    { value: 'inspection', label: 'รอตรวจรับ' }, { value: 'returned', label: 'รับคืนแล้ว' }, { value: 'overdue', label: 'เกินกำหนดส่งคืน' }
  ] },
  { key: 'legal_status', label: 'การดำเนินคดี', options: [
    { value: NO_WORKFLOW, label: 'ไม่มี' }, { value: 'legal_review', label: 'รอฝ่ายกฎหมาย' }, { value: 'in_progress', label: 'กำลังดำเนินคดี' },
    { value: 'enforcement', label: 'บังคับคดี' }, { value: 'closed', label: 'สิ้นสุดคดี' }
  ] }
]
const contractWorkflowDraft = ref<Record<string, Partial<Record<ContractWorkflowKey, string>>>>({})
const contractWorkflowSaving = ref<string>('')
const workflowSaved = (contract: any, key: ContractWorkflowKey) => String(contract?.[key] || '') || NO_WORKFLOW
const workflowValue = (contract: any, key: ContractWorkflowKey) => contractWorkflowDraft.value[String(contract?.id)]?.[key] ?? workflowSaved(contract, key)
const setWorkflowValue = (contract: any, key: ContractWorkflowKey, value: unknown) => {
  const id = String(contract?.id)
  contractWorkflowDraft.value = { ...contractWorkflowDraft.value, [id]: { ...(contractWorkflowDraft.value[id] || {}), [key]: String(value) } }
}
const isWorkflowDirty = (contract: any) => contractWorkflowFields.some(field => workflowValue(contract, field.key) !== workflowSaved(contract, field.key))
const saveContractWorkflow = async (contract: any) => {
  if (!contract?.id) return
  contractWorkflowSaving.value = String(contract.id)
  try {
    const body = Object.fromEntries(contractWorkflowFields.map((field) => {
      const value = workflowValue(contract, field.key)
      return [field.key, value === NO_WORKFLOW ? '' : value]
    }))
    const saved: any = unwrap(await patch(`/contracts/${contract.id}/workflow`, body)) || {}
    const next = Object.fromEntries(contractWorkflowFields.map(field => [field.key, String(saved[field.key] ?? '')]))
    Object.assign(contract, next)
    marketContracts.value = marketContracts.value.map((item: any) => String(item.id) === String(contract.id) ? { ...item, ...next } : item)
    const drafts = { ...contractWorkflowDraft.value }
    delete drafts[String(contract.id)]
    contractWorkflowDraft.value = drafts
    void showSuccess('บันทึกสถานะติดตามสัญญาแล้ว แดชบอร์ดอัปเดตตามข้อมูลนี้', 'บันทึกสำเร็จ')
  } catch {
    void showError('ไม่สามารถบันทึกสถานะติดตามสัญญาได้ กรุณาตรวจสอบสิทธิ์แล้วลองใหม่', 'บันทึกไม่สำเร็จ')
  } finally {
    contractWorkflowSaving.value = ''
  }
}

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
      // สถานะมีผู้เช่าต้องสัมพันธ์กับสัญญา active ที่ backend ส่งกลับมาเสมอ
      syncStallDataFromContracts(code)
    }
    if (hRes.status === 'fulfilled' && Array.isArray(unwrap(hRes.value))) {
      stallHistory.value = unwrap(hRes.value)
    }
  } catch {}
}

let scheduledStallDetailLoadToken = 0
const scheduleStallExtraDetailsAfterOpen = (code: string) => {
  const token = ++scheduledStallDetailLoadToken
  if (!import.meta.client) {
    void loadStallExtraDetails(false)
    return
  }
  // ปล่อยให้ issdExpandFromCard (0.75s) และ content fade วาดจบก่อน
  // เพื่อไม่ให้ response จาก API กระตุก animation ที่ได้รับอนุมัติไว้
  window.setTimeout(() => {
    if (token !== scheduledStallDetailLoadToken) return
    if (!showInspectStallOverlay.value) return
    const currentCode = String(stallEdit.value.ident || stallEdit.value.stall_ident || '').trim().toUpperCase()
    if (currentCode !== code.trim().toUpperCase()) return
    void loadStallExtraDetails(false)
  }, 880)
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
    // กล่องเดิมค้างอยู่: เนื้อหาโครงการหายทันที แล้วกล่องรายละเอียด
    // ขยายขอบขวาจาก 380px เป็น 720px ด้วย issdExpandFromCard
    isProjectMorphSwap.value = true
    stallDrawerExpanded.value = true
    isProjectPanelMorphing.value = true
    showInspectStallOverlay.value = true
    stallDrawerResizeTimer = window.setTimeout(() => {
      isProjectPanelMorphing.value = false
      isProjectMorphSwap.value = false
      stallDrawerResizeTimer = null
    }, 850)
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
    // การกลับจากหน้ารายละเอียดภายนอกต้องคืน state ทันที ไม่เล่น morph ซ้ำ
    // เพราะ openStallFullDetail จะล้าง timer ของ showSelectedStallPanel และทำให้
    // panel โครงการค้างอยู่แทนรายละเอียดแผง
    clearStallDrawerResizeTimer()
    clearStallDrawerSwapTimer()
    isProjectPanelMorphing.value = false
    isProjectPanelReturning.value = false
    isProjectMorphSwap.value = false
    isStallDrawerClosing.value = false
    fullDetailSettling.value = false
    showInspectStallOverlay.value = true
    stallDrawerExpanded.value = true
    stallDrawerMode.value = 'full-detail'
    stallDetailTab.value = (allowedTabs.includes(tab) ? tab : 'info') as any
    await loadStallExtraDetails(false)
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
const showStallEditModal = computed({
  get: () => stallDrawerMode.value === 'edit',
  set: (open: boolean) => { if (!open && stallDrawerMode.value === 'edit') cancelStallEdit() }
})
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
    // occupied ไม่ใช้ค่าค้างในข้อมูลแผงเป็นแหล่งจริง เพราะอาจทำให้แสดงผู้เช่า
    // ทั้งที่ไม่มีสัญญา active; สัญญาด้านล่างเป็นผู้กำหนด occupied เท่านั้น
    if (rawStatus === 'reserved') setStatus(code, 'reserved')
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

  const contractMap = new Map<string, any>()
  ;[...marketContracts.value, ...contracts.value].forEach((c: any) => {
    const key = String(c.id || c.contract_number || c.contract_no || '')
    if (key) contractMap.set(key, c)
  })
  contractMap.forEach((c: any) => {
    const st = String(c.contract_state || c.status || '').toLowerCase()
    const target = ['active', 'approved', 'effective', 'signed', 'ใช้งานอยู่'].includes(st)
      ? 'occupied'
      : ['pending', 'draft', 'contract_pending', 'รออนุมัติ'].includes(st) ? 'reserved' : ''
    if (!target) return
    const rawCodes = String(c.stall_ident || c.stall_number || c.stall_code || '')
    rawCodes.split(',').forEach((code) => {
      const clean = code.trim()
      if (clean) setStatus(clean, target)
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
  lastBookingAt?: number
  renewalStatus?: 'pending' | 'reviewing' | 'approved' | 'declined'
  returnStatus?: 'notified' | 'pending' | 'inspection' | 'returned' | 'overdue'
  legalStatus?: 'legal_review' | 'proceeding' | 'enforcement' | 'closed'
  maintenanceStatus?: 'pending' | 'repairing' | 'completed' | 'overdue'
  requestTypes?: Array<'contract_name_change' | 'transfer' | 'contract_cancel' | 'sublease'>
  searchTerms?: string[]
  searchReferences?: Array<{ type: 'booking' | 'application'; number: string }>
}

// ตัวกรองสถานะที่ใช้อยู่ในแผนที่ (ส่งกลับจากแผนที่) ใช้ไฮไลต์ตัวเลขบนแดชบอร์ดที่กำลังกรองอยู่
type DashboardFilterKey = 'contractStatus' | 'renewalStatus' | 'returnStatus' | 'legalStatus' | 'maintenanceStatus'
const operationalFilterState = ref<Record<string, string>>({})
const onOperationalFilterChange = (filters: Record<string, string>) => { operationalFilterState.value = filters }
const isDashFilterActive = (key: DashboardFilterKey, value: string) => operationalFilterState.value[key] === value
// กดซ้ำเพื่อล้าง และกรองทีละมิติ (สัญญา/ต่ออายุ/ส่งคืน/คดี/ซ่อม) เพื่อให้ผลลัพธ์อ่านง่าย
const toggleDashboardFilter = (key: DashboardFilterKey, value: string) => {
  const next = isDashFilterActive(key, value) ? 'all' : value
  mapRef.value?.applyDashboardFilter?.({
    zone: next === 'all' ? 'all' : (selectedDashboardProjectId.value || 'all'),
    contractStatus: 'all', renewalStatus: 'all', returnStatus: 'all', legalStatus: 'all', maintenanceStatus: 'all',
    [key]: next
  })
}

const stallOperationalFilterMeta = computed<Record<string, StallOperationalMeta>>(() => {
  const result: Record<string, StallOperationalMeta> = {}
  const ensure = (code: string) => {
    const cleanCode = code.trim().toUpperCase()
    if (!cleanCode) return null
    if (!result[cleanCode]) result[cleanCode] = { requestTypes: [], searchTerms: [], searchReferences: [] }
    return result[cleanCode]
  }
  const recordCodes = (record: any) => String(record?.stall_ident || record?.stall_number || record?.stall_code || record?.ident || record?.code || '')
    .split(',').map(code => code.trim().toUpperCase()).filter(Boolean)
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const addSearchTerms = (record: any, type: 'booking' | 'application', terms: unknown[]) => {
    recordCodes(record).forEach((code) => {
      const meta = ensure(code)
      if (!meta) return
      terms.map(term => String(term || '').trim()).filter(Boolean).forEach((term) => {
        if (!meta.searchTerms?.includes(term)) meta.searchTerms?.push(term)
        if (!meta.searchReferences?.some(reference => reference.type === type && reference.number === term)) {
          meta.searchReferences?.push({ type, number: term })
        }
      })
    })
  }

  marketBookings.value.forEach((booking: any) => addSearchTerms(booking, 'booking', [
    booking.booking_number,
    booking.booking_no,
    booking.reference_number
  ]))
  // วันที่จองล่าสุดของแต่ละแผง (ไม่นับที่ยกเลิก/ปฏิเสธ) ใช้กรอง "มีการจองในช่วง"
  const inactiveBookingStatuses = new Set(['cancelled', 'rejected', 'not_interested', 'no_show', 'terminated'])
  marketBookings.value.forEach((booking: any) => {
    if (inactiveBookingStatuses.has(String(booking.booking_status || booking.status || '').toLowerCase())) return
    const time = new Date(booking.created_at || booking.visit_date || 0).getTime()
    if (!Number.isFinite(time) || time <= 0) return
    recordCodes(booking).forEach((code) => {
      const meta = ensure(code)
      if (meta) meta.lastBookingAt = Math.max(meta.lastBookingAt || 0, time)
    })
  })

  marketApplications.value.forEach((application: any) => addSearchTerms(application, 'application', [
    application.application_number,
    application.document_number,
    application.reference_number,
    application.booking_number,
    application.source_booking_number
  ]))

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
      reviewing: 'reviewing', under_review: 'reviewing',
      approved: 'approved', renewed: 'approved', completed: 'approved', submitted: 'approved',
      declined: 'declined', non_renewal: 'declined', not_renewing: 'declined'
    }
    const returnRaw = String(contract.return_status || contract.handover_status || contract.vacate_status || '').toLowerCase()
    const returnMap: Record<string, StallOperationalMeta['returnStatus']> = {
      notified: 'notified', notice: 'notified', non_renewal: 'notified',
      waiting: 'pending', pending_return: 'pending', vacating: 'pending',
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

  // งานซ่อมถูกโหลดจาก Backend ครั้งเดียวเป็น index ตามรหัสแผง
  // ไม่ต้องเดินครบทุก polygon และอ่าน localStorage ซ้ำรายแผงอีกต่อไป
  void maintenanceIndexVersion.value
  Object.entries(maintenanceByStall.value).forEach(([code, records]) => {
    const meta = ensure(code)
    if (!meta || !records.length) return
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
  })
  return result
})

// รายชื่อรหัสแผงที่มีสิ่งที่เจ้าหน้าที่ต้องดำเนินการ (ค้างชำระ / ใกล้หมดสัญญา / งานซ่อมค้าง / รออนุมัติจองหรือคำร้อง)
const actionRequiredStallCodes = computed<string[]>(() => {
  const actionSet = new Set<string>()
  const addCode = (rawCode: unknown) => {
    const code = String(rawCode || '').trim().toUpperCase()
    if (code) actionSet.add(code)
  }
  const addRecordCodes = (record: any) => {
    extractStallCodesFromRecord(record).forEach(addCode)
  }
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const expiringCutoff = new Date(today)
  expiringCutoff.setDate(expiringCutoff.getDate() + 60)
  // วันที่จาก Backend อยู่ในรูป YYYY-MM-DD ซึ่ง JavaScript แปลเป็น UTC
  // เปิดขอบเขตถึงสิ้นวันเพื่อให้สัญญาที่ครบ 60 วันพอดีไม่หลุดในเขตเวลาไทย
  expiringCutoff.setHours(23, 59, 59, 999)
  const bookingActionStates = new Set(['pending', 'pending_approval', 'pending_slot', 'approved', 'interested'])
  const applicationTerminalStates = new Set(['completed', 'rejected', 'cancelled', 'terminated', 'closed'])
  const allCustomContracts = Object.values(customStallContracts.value).flat()
  const allContractsForAction = [...marketContracts.value, ...allCustomContracts]
  const activeContractCodes = new Set<string>()
  allContractsForAction.forEach((contract: any) => {
    const status = String(contract.contract_state || contract.status || '').toLowerCase()
    if (status === 'active' || status === 'ใช้งานอยู่') {
      extractStallCodesFromRecord(contract).forEach(code => activeContractCodes.add(code))
    }
  })

  // 1. การจองทุกสถานะที่กระบวนการยังไม่จบ รวมถึงนัดที่อนุมัติแล้ว
  marketBookings.value.forEach((b: any) => {
    const st = String(b.booking_status || b.status || '').toLowerCase()
    if (bookingActionStates.has(st)) addRecordCodes(b)
  })

  // 2. คำร้องทุกสถานะที่ยังต้องตรวจ อนุมัติ หรือออกสัญญาต่อ
  marketApplications.value.forEach((a: any) => {
    const st = String(a.application_state || a.status || '').toLowerCase()
    if (!st || applicationTerminalStates.has(st)) return
    const codes = extractStallCodesFromRecord(a)
    const isApproved = ['approved', 'director_approved'].includes(st)
    codes.forEach((code) => {
      if (!isApproved || !activeContractCodes.has(code)) addCode(code)
    })
  })

  // 3. สัญญาเช่าที่ค้างชำระหรือใกล้หมดสัญญา
  allContractsForAction.forEach((contract: any) => {
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

    if (isOverdue || isExpiring) rawCodes.split(',').forEach(addCode)
  })

  // 4. หนี้ค้างจาก Backend รองรับทั้งฟิลด์รหัสแผงและรหัสที่อยู่ในรายละเอียดหนี้
  marketDebts.value.forEach((debt: any) => {
    const status = String(debt.status || '').trim().toLowerCase()
    if (!['outstanding', 'overdue', 'past_due', 'ค้างชำระ'].includes(status)) return
    addRecordCodes(debt)
  })

  // 5. งานซ่อมของทุกแผงที่ยังไม่ปิดงาน
  void maintenanceIndexVersion.value
  Object.entries(maintenanceByStall.value).forEach(([code, records]) => {
    const hasOpenWork = records.some((item: any) => {
      const status = String(item.status || '').trim().toLowerCase()
      return !['completed', 'done', 'rejected', 'cancelled', 'closed'].includes(status)
    })
    if (hasOpenWork) addCode(code)
  })

  return Array.from(actionSet)
})

// ข้อมูลตัวเลขสรุปที่แสดงบนแถบฟิลเตอร์ด้านบน: สอดคล้องกับโซนที่เลือก หรือทั้งตลาด
const filteredSummary = ref<{ active: boolean; total: number; available: number; reserved: number; occupied: number; paused: number; codes: Set<string> }>({ active: false, total: 0, available: 0, reserved: 0, occupied: 0, paused: 0, codes: new Set() })
const onFilteredSummary = (summary: typeof filteredSummary.value) => { filteredSummary.value = summary }
// การ์ดตัวเลขตามตัวกรองบนผัง: มีตัวกรองใช้อยู่ให้นับเฉพาะแผงที่ตรงตัวกรอง
const currentDisplayedSummary = computed(() => {
  const base = baseDisplayedSummary.value
  const f = filteredSummary.value
  if (!f.active || selectedProject.value) return base
  const countIn = (codes: string[]) => codes.reduce((n, code) => n + (f.codes.has(String(code).trim().toUpperCase()) ? 1 : 0), 0)
  return {
    ...base,
    total: f.total,
    available: f.available,
    reserved: f.reserved,
    occupied: f.occupied,
    actionRequired: countIn(actionRequiredStallCodes.value),
    legalCases: countIn(legalCaseStallCodes.value)
  }
})
const baseDisplayedSummary = computed(() => {
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
      chartLabel: Number.isFinite(Number(project.sectionNumber)) && Number(project.sectionNumber) > 0
        ? String(Number(project.sectionNumber)).padStart(2, '0')
        : (String(project.name || '').match(/\d+/)?.[0]
            ? String(Number(String(project.name || '').match(/\d+/)?.[0])).padStart(2, '0')
            : cleanProjectName(project.name)),
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

const allProjectsChartSummary = computed(() => stallProjectStatusBreakdown.value.reduce(
  (summary, project) => {
    summary.available += project.available
    summary.reserved += project.reserved
    summary.occupied += project.occupied
    return summary
  },
  { available: 0, reserved: 0, occupied: 0 }
))

const niceChartMaximum = (value: number) => {
  if (value <= 0) return 100
  const roughStep = value / 3
  const magnitude = 10 ** Math.floor(Math.log10(roughStep))
  const normalized = roughStep / magnitude
  const niceStep = normalized <= 1 ? 1 : normalized <= 2 ? 2 : normalized <= 5 ? 5 : 10
  return Math.ceil(value / (niceStep * magnitude)) * niceStep * magnitude
}

// กราฟแท่งแนวนอนของภาพรวมโครงการ: ชื่อโครงการซ้าย แท่งแบ่งสีตามสถานะ จำนวนแผงรวมที่ปลายแท่ง
const sortedProjectBars = computed(() => [...stallProjectStatusBreakdown.value].sort((a, b) => a.order - b.order))

const allProjectsChartMax = computed(() => niceChartMaximum(
  Math.max(0, ...stallProjectStatusBreakdown.value.map(project => project.total))
))

const allProjectsChartTicks = computed(() => {
  const maximum = allProjectsChartMax.value
  return [0, maximum / 3, (maximum / 3) * 2, maximum]
    .map(value => Math.round(value))
    .filter((value, index, values) => index === 0 || value !== values[index - 1])
})

const selectDashboardProject = async (projectId: string) => {
  const id = String(projectId || '')
  const project = visibleProjects.value.find(item => item.id === id)
  if (!project) return
  if (selectedDashboardProjectId.value === id) {
    if (selectedProject.value?.id === id && typeof mapRef.value?.exitProject === 'function') {
      mapRef.value.exitProject()
    }
    selectedDashboardProjectId.value = ''
    return
  }
  selectedDashboardProjectId.value = id
  if (selectedProject.value?.id === id) return
  await nextTick()
  if (typeof mapRef.value?.enterProject === 'function') {
    mapRef.value.enterProject(project)
  } else {
    selectProject(project)
  }
}
const selectedDashboardProjectStatus = computed(() => stallProjectStatusBreakdown.value.find(
  project => project.id === selectedDashboardProjectId.value
) || null)

watch(
  [() => selectedProject.value?.id, stallProjectStatusBreakdown],
  ([mapProjectId, projects]) => {
    if (mapProjectId && projects.some(project => project.id === mapProjectId)) {
      selectedDashboardProjectId.value = mapProjectId
      return
    }
    if (selectedDashboardProjectId.value && !projects.some(project => project.id === selectedDashboardProjectId.value)) {
      selectedDashboardProjectId.value = ''
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

const renewalBarTotal = computed(() => stallContractRenewalSummary.value.submitted + stallContractRenewalSummary.value.pending + stallContractRenewalSummary.value.declined)
const stallReturnTotal = computed(() => Object.values(stallReturnStatusSummary.value).reduce((sum, value) => sum + value, 0))
const stallLegalTotal = computed(() => Object.values(stallLegalStatusSummary.value).reduce((sum, value) => sum + value, 0))

const stallMaintenanceDashboardSummary = computed(() => {
  const summary = { pending: 0, repairing: 0, completed: 0, overdue: 0 }
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  void maintenanceIndexVersion.value
  let projectCodes: Set<string> | null = null
  if (selectedDashboardProjectId.value) {
    projectCodes = new Set<string>()
    const selected = visibleProjects.value.find(project => project.id === selectedDashboardProjectId.value)
    if (selected) {
      getAllFloorLayoutsForProject(selected).forEach((layout) => {
        layout.stalls.forEach((stall: any) => projectCodes?.add(String(stall.code || '').trim().toUpperCase()))
      })
    }
  }
  Object.entries(maintenanceByStall.value).forEach(([code, records]) => {
    if (projectCodes && !projectCodes.has(code)) return
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
  })
  return summary
})

const stallMaintenanceTotal = computed(() => {
  const summary = stallMaintenanceDashboardSummary.value
  return summary.pending + summary.repairing + summary.completed + summary.overdue
})

const stallMaintenanceDonutStyle = computed(() => {
  const summary = stallMaintenanceDashboardSummary.value
  const total = stallMaintenanceTotal.value
  if (!total) return { background: '#e2e8f0' }
  const pendingEnd = (summary.pending / total) * 100
  const repairingEnd = pendingEnd + ((summary.repairing / total) * 100)
  const completedEnd = repairingEnd + ((summary.completed / total) * 100)
  return {
    background: `conic-gradient(#f59e0b 0 ${pendingEnd}%, #2563eb ${pendingEnd}% ${repairingEnd}%, #16a34a ${repairingEnd}% ${completedEnd}%, #ef4444 ${completedEnd}% 100%)`
  }
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
  touchManagedMapData()
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

const ensureCanLeaveOrBack = async (): Promise<boolean> => {
  if (isInspectMode.value || !hasUnsavedChanges.value) return true
  const leave = await showConfirm('ข้อมูลผังตลาดมีการแก้ไขที่ยังไม่ได้บันทึก ออกจากหน้านี้และทิ้งการแก้ไขหรือไม่', {
    title: 'มีการแก้ไขที่ยังไม่ได้บันทึก',
    confirmText: 'ทิ้งและออกจากหน้านี้'
  })
  if (leave) revertUnsavedChanges()
  return leave
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
// These handlers only run after an actual GIS edit. Avoid serializing the entire
// 13MB project tree merely to determine whether the page is dirty.
const markGISChanged = () => {
  hasUnsavedChanges.value = true
  touchManagedMapData()
}

const updateProjectPolygon = async (project: ProjectZone) => {
  const index = projects.value.findIndex(item => item.id === project.id)
  if (index >= 0) projects.value[index] = { ...projects.value[index], ...project }
  if (selectedProject.value?.id === project.id) {
    selectedProject.value = projects.value[index] || project
  }
  markGISChanged()
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
      coordinates_geojson: JSON.stringify({ type: 'Polygon', coordinates: [updatedStall.coordinates, ...(updatedStall.polygonHoles || [])] })
    }
    const sIdx = stalls.value.findIndex(s => String(s.ident || s.stall_ident).toUpperCase() === code.toUpperCase())
    if (sIdx >= 0) stalls.value[sIdx] = updatedRecord
    queueStallUpsert(updatedRecord)
  }
  markGISChanged()
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
const showProjectEditor = ref(false)
const projectEditId = ref('')
const projectForm = ref({ name: '', area: '' as string | number, stallWidth: '' as string | number, stallDepth: '' as string | number })
const projectErrors = reactive({ name: '', area: '', stallWidth: '', stallDepth: '' })
const openProjectEditor = (projectId: string) => {
  const project = projects.value.find(item => item.id === projectId)
  if (!project) return
  projectEditId.value = projectId
  projectForm.value = {
    name: String(project.name || ''),
    area: project.realAreaSqm || '',
    stallWidth: project.standardStallWidthM || 2.5,
    stallDepth: project.standardStallDepthM || 2
  }
  projectErrors.name = ''
  projectErrors.area = ''
  projectErrors.stallWidth = ''
  projectErrors.stallDepth = ''
  showProjectEditor.value = true
}
const validateProjectField = (field: 'name' | 'area' | 'stallWidth' | 'stallDepth') => {
  if (field === 'name') {
    projectErrors.name = String(projectForm.value.name || '').trim() ? '' : 'กรุณาระบุชื่อโครงการ'
    return
  }
  const raw = projectForm.value[field]
  if (raw === '' || raw === null || raw === undefined) { projectErrors[field] = ''; return }
  const value = Number(raw)
  projectErrors[field] = Number.isFinite(value) && value > 0 ? '' : 'ต้องมากกว่า 0'
}
const saveProjectEdit = () => {
  ;(['name', 'area', 'stallWidth', 'stallDepth'] as const).forEach(validateProjectField)
  if (projectErrors.name || projectErrors.area || projectErrors.stallWidth || projectErrors.stallDepth) return
  const idx = projects.value.findIndex(item => item.id === projectEditId.value)
  if (idx < 0) return
  const current = projects.value[idx]!
  const area = Number(projectForm.value.area)
  projects.value[idx] = {
    ...current,
    name: String(projectForm.value.name).trim(),
    realAreaSqm: area > 0 ? area : current.realAreaSqm,
    standardStallWidthM: Number(projectForm.value.stallWidth) || undefined,
    standardStallDepthM: Number(projectForm.value.stallDepth) || undefined
  }
  if (selectedProject.value?.id === current.id) selectedProject.value = projects.value[idx]!
  markGISChanged()
  showProjectEditor.value = false
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
  markGISChanged()
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
  markGISChanged()
  closeFacilityEditor()
  void showSuccess('เพิ่มจุดบริการลงในผังแล้ว กดปุ่ม “บันทึกการแก้ไข” ด้านบนเพื่อบันทึกผัง', 'เพิ่มจุดบริการในผังแล้ว')
}
const deleteFacility = async (facility: MarketFacility) => {
  const confirmed = await showConfirm(`ต้องการลบจุดบริการ “${facility.name}” ออกจากผังตลาดใช่หรือไม่`, { title: 'ยืนยันการลบจุดบริการ', type: 'error', icon: 'trash', confirmText: 'ลบจุดบริการ' })
  if (!confirmed) return
  facilities.value = facilities.value.filter(item => item.id !== facility.id)
  markGISChanged()
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
const openDrawnStall = (project: ProjectZone, coordinates: [number, number][], floorId: string, dimensions?: { widthM: number; depthM: number }) => { const closed = [...coordinates.map(point => [...point] as [number, number]), [...coordinates[0]!]]; const floor = project.floors?.find(item => item.id === floorId) || project.floors?.[0] || { id: '1', name: 'ชั้น 1', level: 1 }; const count = (project.stalls || []).filter((stall: any) => String(stall.floor || '1') === String(floor.id)).length + 1; pendingStallLayout.value = { projectId: project.id, coordinates: closed, floorId: String(floor.id) }; stallEdit.value = { ...emptyStall(), ident: `${project.sectionNumber}-${floor.id}-${String(count).padStart(3, '0')}`, zone_name: project.name, project_name: project.name, category: project.category || marketCategoryOptions.value[0]?.value || '', floor: String(floor.id), coordinates_geojson: JSON.stringify({ type: 'Polygon', coordinates: [closed] }), width: dimensions?.widthM || 2, depth: dimensions?.depthM || 2 }; resetStallErrors(); stallDrawerMode.value = 'edit'; showInspectStallOverlay.value = true; showStallEditor.value = false }
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
  markGISChanged()
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
  selectedStallSurveyNumber.value = ''
}
const closeStall = () => {
  clearStallDrawerResizeTimer()
  const shouldCollapseBeforeClose = import.meta.client
    && isMapOverviewMode.value
    && showInspectStallOverlay.value
    && stallDrawerMode.value === 'full-detail'

  if (shouldCollapseBeforeClose) {
    // ส่วนขวาหุบกลับด้วย clip-path เหลือ 380px ขณะที่เนื้อหาโครงการเฟดกลับเข้ามาใต้กล่อง
    isStallDrawerClosing.value = true
    isProjectPanelReturning.value = true
    stallDrawerResizeTimer = window.setTimeout(() => {
      stallDrawerResizeTimer = null
      finishCloseStall(true)
      stallDrawerResizeTimer = window.setTimeout(() => {
        isProjectPanelReturning.value = false
        stallDrawerResizeTimer = null
      }, 400)
    }, 640)
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
        zone_name: selectedProject.value?.name || stored.project_name || stored.zone_name || '',
        project_name: selectedProject.value?.name || stored.project_name || '',
        stall_status: resolveRealStallStatus(code, selectedProject.value, stored)
      }
    }
    resetStallErrors()
    stallDrawerMode.value = isInspectMode.value ? 'full-detail' : 'detail'
    stallDrawerExpanded.value = stallDrawerMode.value === 'full-detail'
  }
}
const selectStall = async (mapStall: any) => {
  if (isMapOverviewMode.value && isCostMode.value) {
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
  const sourceNumber = mapStall.sourceProperties?.number ?? mapStall.sourceProperties?.Number
  selectedStallSurveyNumber.value = String(sourceNumber ?? '').trim()
  const stored = stalls.value.find(s => String(s.ident || s.stall_ident).toUpperCase() === code.toUpperCase())
  const geometry = mapStall.coordinates?.length ? JSON.stringify({ type: 'Polygon', coordinates: [mapStall.coordinates, ...(mapStall.polygonHoles || [])] }) : ''
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
      name: selectedStallSurveyNumber.value || mapStall.name || '',
      category: mapStall.areaType || selectedProject.value?.category || '',
      zone_name: selectedProject.value?.name || mapStall.zone || '',
      project_name: selectedProject.value?.name || mapStall.zone || '',
      floor: String(mapStall.floor || '1'),
      width: mapStall.widthM || 2,
      depth: mapStall.depthM || 2,
      monthly_rent: mapStall.monthlyRent || 1800,
      stall_status: resolveRealStallStatus(code, selectedProject.value, mapStall),
      bookable: mapStall.bookable !== false,
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
      scheduleStallExtraDetailsAfterOpen(code)
    }
    return
  }
  stallEdit.value = {
    ...emptyStall(),
    ...stored,
    zone_name: selectedProject.value?.name || stored.project_name || stored.zone_name || mapStall.zone || '',
    project_name: selectedProject.value?.name || stored.project_name || mapStall.zone || '',
    name: selectedStallSurveyNumber.value || stored.name || mapStall.name || '',
    // ข้อมูลแผงใช้ค่าจากไฟล์ GIS (mapStall) เป็นหลัก ฐานข้อมูลแผงเป็นแค่ค่าสำรอง
    category: mapStall.areaType || stored.category || '',
    floor: String(mapStall.floor || stored.floor || '1'),
    width: mapStall.widthM || stored.width || 2,
    depth: mapStall.depthM || stored.depth || 2,
    monthly_rent: mapStall.monthlyRent || stored.monthly_rent || 1800,
    coordinates_geojson: stored.coordinates_geojson || geometry,
    bookable: mapStall.bookable !== false,
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
    scheduleStallExtraDetailsAfterOpen(code)
  }
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
  markGISChanged()
  if (showStallEditor.value) closeStall()

  void showSuccess(`ลบแผง “${code}” ออกจากผังแล้ว กดปุ่ม “บันทึกการแก้ไข” ด้านบนเพื่อบันทึกผัง`, 'ลบแผงในผังแล้ว')
}
const quickMissingFields = computed(() => {
  const missing: string[] = []
  if (!Number(stallEdit.value.monthly_rent)) missing.push('ค่าเช่า')
  if (!Number(stallEdit.value.width) || !Number(stallEdit.value.depth)) missing.push('ขนาด')
  return missing
})
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
          bookable: stallEdit.value.bookable !== false,
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
            monthlyRent: Number(stallEdit.value.monthly_rent) || st.monthlyRent,
            bookable: stallEdit.value.bookable !== false
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
    markGISChanged()
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
  const requestedMarketId = marketId.value
  loadDataController?.abort()
  const controller = new AbortController()
  loadDataController = controller
  const run = ++loadDataRun
  const isStale = () => controller.signal.aborted || run !== loadDataRun
  isManageDataLoading.value = true
  isActionDataLoading.value = true
  // The map component owns the authoritative /gis-config request. The page only
  // needs the synchronous registry metadata plus the editable draft below.
  const cfg = getMarketGisConfig(requestedMarketId)
  apiMarketName.value = cfg.name
  projects.value = cfg.projects.slice()
  facilities.value = []
  try {
    const scopedMarketID = Number(cfg.id)
    const scopedGet = <T = any>(endpoint: string, query?: Record<string, any>) =>
      apiFetch<T>(endpoint, { method: 'GET', query, signal: controller.signal })

    // Start all requests together, but only block the KPI skeleton on the four
    // datasets it actually needs. Draft GIS and operational history can hydrate
    // after the summary is already usable.
    const coreRequest = Promise.all([
      scopedGet('/stalls', { market_id: scopedMarketID }).catch(() => null),
      scopedGet('/bookings', { market_id: scopedMarketID }).catch(() => null),
      scopedGet('/contracts', { market_id: scopedMarketID }).catch(() => null)
    ])
    const detailRequest = Promise.all([
      // Published GIS is already owned by the map in inspect mode. Only the GIS
      // editor needs the much larger draft payload.
      isInspectMode.value
        ? Promise.resolve(null)
        : scopedGet(`/markets/${encodeURIComponent(requestedMarketId)}/gis-config/draft`).catch(() => null),
      scopedGet('/bookings/waitlist/mine', { market_id: scopedMarketID }).catch(() => null),
      scopedGet('/applications', { market_id: scopedMarketID }).catch(() => null),
      scopedGet('/maintenance-records', { market_id: scopedMarketID }).catch(() => null),
      scopedGet('/debts', { market_id: scopedMarketID }).catch(() => null)
    ])

    const [stallRes, bookingsRes, contractsRes] = await coreRequest
    if (isStale()) return
    // The route registry already resolves the selected market. Loading the
    // complete market list here delayed every cold map refresh for no benefit.
    apiMarketRecord.value = { id: cfg.id, name: cfg.name, slug: requestedMarketId }
    resolvedMarketID.value = Number(cfg.id)
    apiMarketName.value = cfg.name
    const stallData = unwrap(stallRes)
    stalls.value = Array.isArray(stallData)
      ? stallData.filter((s: any) => !s.market_id || Number(s.market_id) === resolvedMarketID.value)
      : []

    const allBookings = unwrap(bookingsRes)
    const allContracts = unwrap(contractsRes)
    const targetMktName = String(apiMarketName.value || cfg.name || '').trim()
    const matchesCurrentMarket = (item: any) => {
      if (!item) return false
      if (Number(item.market_id) > 0) {
        return Number(item.market_id) === resolvedMarketID.value
      }
      const itemMktName = String(item.market_name || '').trim()
      if (itemMktName && targetMktName) {
        return itemMktName.includes(targetMktName) || targetMktName.includes(itemMktName)
      }
      return false
    }
    marketBookings.value = Array.isArray(allBookings) ? allBookings.filter(matchesCurrentMarket) : []
    marketContracts.value = Array.isArray(allContracts) ? allContracts.filter(matchesCurrentMarket) : []
    isManageDataLoading.value = false

    const targetProjectId = returnProjectId.value || selectedProject.value?.id
    if (targetProjectId) selectedProject.value = projects.value.find(p => p.id === targetProjectId) || null

    const [gisRes, waitlistRes, appsRes, maintenanceRes, debtsRes] = await detailRequest
    if (isStale()) return
    const gis = unwrap(gisRes) || {}
    if (Array.isArray(gis.projects) && gis.projects.length) {
      projects.value = mergeMarketProjects(cfg.projects, gis.projects, false)
      facilities.value = Array.isArray(gis.facilities) ? gis.facilities : []
      gisPublished.value = Boolean(gis.published)
      touchManagedMapData()
    }
    const allWaitlist = unwrap(waitlistRes)
    const allApps = unwrap(appsRes)
    marketWaitlistItems.value = Array.isArray(allWaitlist) ? allWaitlist.filter(matchesCurrentMarket) : []
    marketApplications.value = Array.isArray(allApps) ? allApps.filter(matchesCurrentMarket) : []
    const maintenanceRecords = unwrap(maintenanceRes)
    if (Array.isArray(maintenanceRecords)) replaceMaintenanceIndex(maintenanceRecords)
    const debts = unwrap(debtsRes)
    marketDebts.value = Array.isArray(debts)
      ? debts.filter((debt: any) => Number(debt.market_id) === resolvedMarketID.value)
      : []
    loadAllCustomContracts()
    isActionDataLoading.value = false
    // Snapshots are only required by the editable GIS page. Cloning the complete
    // market in inspect mode delayed the KPI ribbon even after every API resolved.
    if (!isInspectMode.value) syncSnapshotsFromCurrent()
  } catch (error: any) {
    if (isStale() || error?.name === 'AbortError') return
    projects.value = cfg.projects.slice()
    facilities.value = []
    isManageDataLoading.value = false
    isActionDataLoading.value = false
  } finally {
    if (isStale()) {
      reloadAfterActivation = true
      return
    }
    isManageDataLoading.value = false
    isActionDataLoading.value = false
    reloadAfterActivation = false
    loadedMarketId = requestedMarketId
    if (loadDataController === controller) loadDataController = null
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
// โหมดจัดการค่าใช้จ่าย: เชื่อมแผง สัญญา มิเตอร์ อัตราค่าน้ำ-ไฟ และรอบจด
// ================================================================
const billing = useStallBilling()
// โหลดอัตราที่ตั้งไว้ของตลาด (ลงทะเบียนตัวคิดอัตราให้การคิดเงินด้านล่างใช้)
const utilitySettings = useUtilitySettings(marketId.value)
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


const contractStallCodes = (contract: any) => String(contract?.stall_ident || contract?.stall_number || contract?.stall_code || '')
  .split(',')
  .map(code => code.trim().toUpperCase())
  .filter(Boolean)
const isEffectiveContract = (contract: any) => {
  const status = String(contract?.contract_state || contract?.status || '').trim().toLowerCase()
  return ['active', 'approved', 'effective', 'signed', 'ใช้งานอยู่'].includes(status)
}
const activeMeterContracts = computed(() => {
  const byStall = new Map<string, any>()
  // รายการจดมิเตอร์ใช้สัญญาจาก backend เท่านั้น ข้อมูลร่างใน localStorage ห้ามสร้างผู้เช่าหรือมิเตอร์จริง
  marketContracts.value.forEach((contract: any) => {
    if (!isEffectiveContract(contract)) return
    contractStallCodes(contract).forEach(code => {
      if (!byStall.has(code)) byStall.set(code, contract)
    })
  })
  return byStall
})
const meterProjectMap = computed<Record<string, string>>(() => {
  const map: Record<string, string> = {}
  visibleProjects.value.forEach(project => getAllFloorLayoutsForProject(project).forEach(layout => layout.stalls.forEach((stall: any) => {
    const code = String(stall?.code || stall?.ident || stall?.stall_ident || '').trim().toUpperCase()
    if (code) map[code] = String(project.id)
  })))
  return map
})
const testReadingOf = (code: string, kind: 'water' | 'electric') => {
  const seed = [...`${code}:${kind}`].reduce((sum, char) => ((sum * 31) + char.charCodeAt(0)) >>> 0, 0)
  return kind === 'water' ? 100 + (seed % 400) : 1000 + (seed % 3000)
}
const meterRelations = computed(() => {
  const projectMap = meterProjectMap.value
  return [...activeMeterContracts.value.entries()].map(([code, contract]) => {
    const projectId = projectMap[code] || ''
    const binding = utilitySettings.getMeter(code)
    const devices = binding.devices || []
    const activeDevice = (kind: 'water' | 'electric') => devices.find(device => device.kind === kind && device.status === 'active')
    const water = activeDevice('water')
    const electric = activeDevice('electric')
    const contractId = String(contract.id || '')
    const persisted = utilitySettings.getReading(billing.billingMonth.value, code)
    // ห้ามนำเลขที่จดของผู้เช่ารายเดิมมาใช้กับสัญญาใหม่ในแผงเดียวกัน
    const saved = persisted && (!persisted.contractId || persisted.contractId === contractId) ? persisted : null
    return {
      code,
      projectId,
      tenant: contract.tenant_name || 'ผู้เช่าตามสัญญา',
      contractId,
      contractNumber: contract.contract_number || contract.contract_no || '',
      shopName: contract.shop_name || contract.business_name || '',
      prevWater: saved?.prevWater ?? water?.currentReading ?? testReadingOf(code, 'water'),
      prevElectric: saved?.prevElectric ?? electric?.currentReading ?? testReadingOf(code, 'electric'),
      currWater: saved?.currWater ?? null,
      currElectric: saved?.currElectric ?? null,
      status: saved?.status || 'unrecorded',
      waterMeterCode: water?.code || binding.water,
      electricMeterCode: electric?.code || binding.electric
    }
  }).filter(item => item.projectId && (
    utilitySettings.resolve(item.projectId, item.code, 'water').rule.mode === 'meter'
    || utilitySettings.resolve(item.projectId, item.code, 'electric').rule.mode === 'meter'
  ))
})
const ensureMeterRelations = async () => {
  await utilitySettings.ready
  const projectMap = meterProjectMap.value
  registerStallProjects(projectMap)
  let createdMeter = false
  const meterUpdates: Record<string, ReturnType<typeof utilitySettings.getMeter>> = {}
  activeMeterContracts.value.forEach((_contract, code) => {
    const projectId = projectMap[code]
    if (!projectId) return
    const binding = utilitySettings.getMeter(code)
    const devices = [...(binding.devices || [])]
    ;(['electric', 'water'] as const).forEach(kind => {
      if (utilitySettings.resolve(projectId, code, kind).rule.mode !== 'meter') return
      if (devices.some(device => device.kind === kind && device.status === 'active')) return
      const prefix = kind === 'electric' ? 'EM' : 'WM'
      devices.push({
        id: `trial-${kind}-${code}`,
        name: kind === 'electric' ? 'มิเตอร์ไฟฟ้า' : 'มิเตอร์น้ำประปา',
        kind,
        code: `${prefix}-${code}`,
        currentReading: testReadingOf(code, kind),
        status: 'active'
      })
      createdMeter = true
    })
    if (devices.length !== (binding.devices || []).length) meterUpdates[code] = { ...binding, devices }
  })
  if (Object.keys(meterUpdates).length) {
    utilitySettings.setMeters(meterUpdates)
    createdMeter = true
  }
  billing.ensureUtilityBills(meterRelations.value)
  if (createdMeter) void utilitySettings.syncNow()
}
watch([activeMeterContracts, () => billing.billingMonth.value, pageMode], () => {
  if (pageMode.value === 'meters' || pageMode.value === 'payments') void ensureMeterRelations()
}, { immediate: true })
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
  if (occupied) {
    const relation = meterRelations.value.find(item => item.code === code)
    if (relation) billing.ensureUtilityBills([relation])
  }
  // จดมิเตอร์: กดแผงไหนก็เปิดการ์ดของแผงนั้นฝั่งซ้าย (แทนการ์ดโครงการ) แผงว่างขึ้นการ์ดพร้อมข้อความว่าไม่ต้องจด (รอกำหนดว่าจะแสดงอะไรเพิ่ม)
  if (pageMode.value === 'meters') {
    openMeterCard(code)
    return
  }
  // ไม่ได้เปิดแผงไหนอยู่: เลือกชนิดตามสถานะแผง (มีผู้เช่า = ค่าน้ำไฟ, ไม่มี = บิลชั่วคราว)
  openCostDrawer(code, costTool.value ? (costTool.value === 'utility' ? 'utility' : 'temp') : (occupied ? 'utility' : 'temp'))
}

// การจดมิเตอร์น้ำ-ไฟจากสัญญาที่มีผลและมิเตอร์ที่ผูกกับแผง
// กดโครงการ -> การ์ดสรุปการจดของโครงการ (ซ้าย) / กดแผง -> การ์ดจดรายแผงแทนการ์ดโครงการ / บันทึกแล้วไปแผงถัดไปที่ยังไม่จด
// ตารางสรุปการจด (ไอคอนมุมล่างซ้าย) ใช้ดูภาพรวมทุกโครงการ
const meterSummaryOpen = ref(false)
const meterCardCode = ref('')
const meterMonthLabel = computed(() => new Intl.DateTimeFormat('th-TH', { month: 'short', year: 'numeric' }).format(new Date(`${billing.billingMonth.value}-01T00:00:00`)))
const METER_ABNORMAL_FACTOR = 2.5
const stallCodeOf = (stall: any) => String(stall?.code || stall?.ident || stall?.stall_ident || '').trim().toUpperCase()
// แผง -> โครงการ ให้การคิดเงินหาอัตราระดับโครงการได้ (ทำเมื่อเข้าโหมดจดมิเตอร์/ค่าใช้จ่าย)
let registeredProjectsKey = ''
watch([isCostMode, visibleProjects], ([on]) => {
  if (!on) return
  // สร้างแผนที่แผง->โครงการใหม่เฉพาะเมื่อรายการโครงการเปลี่ยนจริง (ไล่ทุกแผงใช้เวลา)
  const key = visibleProjects.value.map(p => `${p.id}:${p.floors?.length ?? 0}`).join('|')
  if (key === registeredProjectsKey) return
  registeredProjectsKey = key
  const map: Record<string, string> = {}
  visibleProjects.value.forEach(project => getAllFloorLayoutsForProject(project).forEach(layout => layout.stalls.forEach((stall: any) => {
    const code = stallCodeOf(stall)
    if (code) map[code] = String(project.id)
  })))
  registerStallProjects(map)
}, { immediate: true })
type MeterData = { projects: MeterProjectRow[]; stalls: MeterStallRow[]; averages: Record<string, { electric: number; water: number }> }
const EMPTY_METER_DATA: MeterData = { projects: [], stalls: [], averages: {} }

// onlyProjectId: คำนวณเฉพาะโครงการเดียว (การ์ดโครงการ) ไม่ระบุ = ทุกโครงการ (ตารางสรุป)
const buildMeterData = (onlyProjectId?: string): MeterData => {
  const info = new Map<string, { projectId: string; projectName: string; number: string }>()
  visibleProjects.value.forEach((project) => {
    if (onlyProjectId && String(project.id) !== onlyProjectId) return
    getAllFloorLayoutsForProject(project).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        const code = stallCodeOf(stall)
        if (!code || info.has(code)) return
        const number = String(stall.sourceProperties?.number ?? stall.sourceProperties?.Number ?? '').trim()
        info.set(code, { projectId: String(project.id), projectName: cleanProjectName(project.name), number })
      })
    })
  })
  // เฉพาะแผงที่มีมิเตอร์ที่ต้องจดอย่างน้อย 1 ประเภท (เหมาจ่าย/ไม่คิดค่าใช้บริการทั้งคู่ไม่ต้องจด)
  const bills = billing.utilityList.value.filter(bill => info.has(bill.stall.toUpperCase()) && (needsReading(bill.stall, 'water') || needsReading(bill.stall, 'electric')))
  // ค่าเฉลี่ยหน่วยของแต่ละโครงการ ใช้ตรวจค่าที่ใช้สูงผิดปกติ (เทียบเมื่อจดแล้วอย่างน้อย 3 แผง)
  const samples = new Map<string, { electric: number[]; water: number[] }>()
  bills.forEach((bill) => {
    if (bill.status === 'unrecorded') return
    const id = info.get(bill.stall.toUpperCase())!.projectId
    const bucket = samples.get(id) || { electric: [], water: [] }
    bucket.electric.push(billing.electricUnits(bill) ?? 0)
    bucket.water.push(billing.waterUnits(bill) ?? 0)
    samples.set(id, bucket)
  })
  const mean = (list: number[]) => (list.length >= 3 ? list.reduce((sum, value) => sum + value, 0) / list.length : 0)
  const averages: MeterData['averages'] = {}
  samples.forEach((bucket, id) => { averages[id] = { electric: mean(bucket.electric), water: mean(bucket.water) } })

  const stalls: MeterStallRow[] = bills.map((bill) => {
    const meta = info.get(bill.stall.toUpperCase())!
    const recorded = bill.status !== 'unrecorded'
    const eUsed = billing.electricUnits(bill)
    const wUsed = billing.waterUnits(bill)
    const projectAverage = averages[meta.projectId] || { electric: 0, water: 0 }
    const abnormal = recorded && (
      (projectAverage.electric > 0 && (eUsed ?? 0) > projectAverage.electric * METER_ABNORMAL_FACTOR)
      || (projectAverage.water > 0 && (wUsed ?? 0) > projectAverage.water * METER_ABNORMAL_FACTOR)
    )
    const skipReason = recorded ? '' : billing.getUtilitySkip(bill.stall)
    // แผงที่ต้องจดมิเตอร์แต่ยังไม่ได้ตั้งเลขมิเตอร์ จดไม่ได้จนกว่าจะเพิ่มมิเตอร์
    const binding = utilitySettings.getMeter(bill.stall)
    const missingMeter = !recorded && (['water', 'electric'] as const).some(kind => needsReading(bill.stall, kind) && !binding[kind])
    return {
      code: bill.stall,
      projectId: meta.projectId,
      label: `แผง ${meta.number || bill.stall}`,
      tenant: bill.tenant,
      recorded,
      abnormal,
      skipped: !!skipReason,
      electricPrev: bill.prevElectric,
      electricCurr: bill.currElectric,
      electricUsed: eUsed,
      waterPrev: bill.prevWater,
      waterCurr: bill.currWater,
      waterUsed: wUsed,
      amount: billing.utilityTotal(bill),
      recordedBy: '',
      // งวดย้อนหลังดูได้อย่างเดียว: แผงที่ไม่มีเลขที่จดไว้ในงวดนั้นเป็น "ไม่มีข้อมูลงวดนี้"
      state: recorded ? (abnormal ? 'abnormal' : 'recorded') : isPastPeriod.value ? 'no-data' : missingMeter ? 'no-meter' : (skipReason ? 'skipped' : 'pending')
    }
  })
  const groups = new Map<string, MeterProjectRow>()
  stalls.forEach((row) => {
    const project = [...info.values()].find(item => item.projectId === row.projectId)
    const sampleCode = row.code
    const modeLabel = (kind: 'water' | 'electric') => ({ meter: 'ตามมิเตอร์', flat: 'เหมาจ่าย', none: 'ไม่คิดค่าใช้บริการ' }[chargeRuleOf(sampleCode, kind).mode])
    const currentMethod = `${modeLabel('electric')} / ${modeLabel('water')}`
    const group = groups.get(row.projectId) || { id: row.projectId, name: project?.projectName || row.projectId, method: currentMethod, recorded: 0, total: 0, abnormal: 0, skipped: 0, electricUnits: 0, waterUnits: 0, ready: false }
    if (group.total > 0 && group.method !== currentMethod) group.method = 'หลายรูปแบบตามแผง'
    group.total += 1
    if (row.recorded) { group.recorded += 1; group.electricUnits += row.electricUsed || 0; group.waterUnits += row.waterUsed || 0 }
    if (row.abnormal) group.abnormal += 1
    if (row.skipped) group.skipped += 1
    groups.set(row.projectId, group)
  })
  const projects = [...groups.values()].map(group => ({ ...group, ready: group.total > 0 && group.recorded === group.total && group.abnormal === 0 }))
  return { projects, stalls, averages }
}
const meterSummaryData = computed<MeterData>(() => (meterSummaryOpen.value ? buildMeterData() : EMPTY_METER_DATA))
const meterProjectData = computed<MeterData>(() => (pageMode.value === 'meters' && selectedProject.value ? buildMeterData(String(selectedProject.value.id)) : EMPTY_METER_DATA))
// ระหว่างสลับตลาด/โครงการ Vue อาจประเมิน computed ก่อนข้อมูลรอบใหม่พร้อม จึงต้องมี fallback ทุกจุดที่อ่านผลลัพธ์
// เหตุผลที่ข้าม -> จำนวนแผง (ใช้ในสรุปการจด)
const meterSkipReasonCounts = computed(() => {
  const counts: Record<string, number> = {}
  meterSummaryData.value.stalls.filter(row => row.state === 'skipped').forEach((row) => {
    const reason = billing.getUtilitySkip(row.code) || 'ไม่ระบุเหตุผล'
    counts[reason] = (counts[reason] || 0) + 1
  })
  return counts
})
const goMeterSettings = () => { void navigateTo(`/stalls/utility-settings?market=${encodeURIComponent(String(marketId.value))}&section=meters`) }
const meterSummaryProjects = computed<MeterProjectRow[]>(() => meterSummaryData.value?.projects ?? [])
const meterProjectStalls = computed<MeterStallRow[]>(() => meterProjectData.value?.stalls ?? [])
const meterProjectAverages = computed<MeterData['averages']>(() => meterProjectData.value?.averages ?? {})
const meterLiveData = computed<MeterData>(() => buildMeterData())
const meterLiveRows = computed<MeterStallRow[]>(() => meterLiveData.value?.stalls ?? [])
const meterLiveProjects = computed<MeterProjectRow[]>(() => meterLiveData.value?.projects ?? [])
const meterAnomalyRows = computed<MeterStallRow[]>(() => meterScopeRows.value.filter(row => row.abnormal))
const meterUsageSummary = computed(() => {
  const recorded = meterScopeRows.value.filter(row => row.recorded)
  const waterValues = recorded.map(row => row.waterUsed).filter((value): value is number => value != null)
  const electricValues = recorded.map(row => row.electricUsed).filter((value): value is number => value != null)
  const waterTotal = waterValues.reduce((sum, value) => sum + value, 0)
  const electricTotal = electricValues.reduce((sum, value) => sum + value, 0)
  return {
    waterTotal,
    electricTotal,
    waterAverage: waterValues.length ? Math.round(waterTotal / waterValues.length) : 0,
    electricAverage: electricValues.length ? Math.round(electricTotal / electricValues.length) : 0
  }
})
const meterOrderedStalls = computed(() => [...meterProjectStalls.value].sort((a, b) => a.label.localeCompare(b.label, 'th', { numeric: true })))
const meterCardRow = computed(() => meterProjectStalls.value.find(row => row.code.toUpperCase() === meterCardCode.value) || null)
// ชื่อแผงบนการ์ด: แผงที่มีผู้เช่าใช้ชื่อจากรายการจด แผงว่างหาเลขแผงจากผังของโครงการที่เปิดอยู่
const meterCardLabel = computed(() => {
  if (meterCardRow.value) return meterCardRow.value.label
  const code = meterCardCode.value
  if (!code) return ''
  let number = ''
  if (selectedProject.value) {
    getAllFloorLayoutsForProject(selectedProject.value).forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        if (!number && stallCodeOf(stall) === code) number = String(stall.sourceProperties?.number ?? stall.sourceProperties?.Number ?? '').trim()
      })
    })
  }
  return `แผง ${number || code}`
})
const meterPendingLeft = computed(() => meterProjectStalls.value.filter(row => row.state === 'pending').length)
const meterAverages = computed(() => meterProjectAverages.value[String(selectedProject.value?.id)] || { electric: 0, water: 0 })
// แผงถัดไปที่ยังไม่จด (ไล่ตามลำดับเลขแผง วนกลับไปต้นรายการ ไม่นับแผงที่ข้ามแล้ว)
const findNextMeterStall = (afterCode: string) => {
  const list = meterOrderedStalls.value
  const start = list.findIndex(row => row.code.toUpperCase() === afterCode.toUpperCase())
  for (let i = 1; i <= list.length; i++) {
    const row = list[(start + i + list.length) % list.length]
    if (row && row.state === 'pending' && row.code.toUpperCase() !== afterCode.toUpperCase()) return row
  }
  return null
}
const meterHasNext = computed(() => !!findNextMeterStall(meterCardCode.value))
// ประวัติแผงที่เพิ่งจดหรือดูผ่านการ์ด ใช้ย้อนกลับไปแผงก่อนหน้า (ล้างเมื่อปิดการ์ด)
const meterCardHistory = ref<string[]>([])
const meterCanGoBack = computed(() => meterCardHistory.value.length > 0)
const openMeterCard = (code: string) => {
  closeCostDrawer()
  const next = String(code).toUpperCase()
  if (meterCardCode.value && meterCardCode.value !== next) meterCardHistory.value.push(meterCardCode.value)
  meterCardCode.value = next
}
const openPrevMeterStall = () => {
  const prev = meterCardHistory.value.pop()
  if (!prev) return
  closeCostDrawer()
  meterCardCode.value = prev
}
watch(meterCardCode, (value) => { if (!value) meterCardHistory.value = [] })
const closeMeterCard = () => { meterCardCode.value = '' }
// ปุ่มปิดบนการ์ดจดรายแผง: กลับหน้ารวมโครงการเลย (ไม่ค้างที่การ์ดสรุปโครงการ)
const closeMeterCardToOverview = () => {
  meterCardCode.value = ''
  clearSelectedProject()
}
const openNextMeterStall = () => {
  const next = findNextMeterStall(meterCardCode.value)
  if (next) openMeterCard(next.code)
  else void showSuccess('จดครบทุกแผงของโครงการนี้แล้ว', 'จดครบแล้ว')
}
const persistRecordedMeterReadings = (code: string) => {
  const bill = billing.getUtility(code)
  if (!bill) return
  const binding = utilitySettings.getMeter(code)
  const devices = (binding.devices || []).map(device => {
    if (device.status !== 'active') return device
    if (device.kind === 'water' && needsReading(code, 'water') && bill.currWater != null) return { ...device, currentReading: bill.currWater }
    if (device.kind === 'electric' && needsReading(code, 'electric') && bill.currElectric != null) return { ...device, currentReading: bill.currElectric }
    return device
  })
  utilitySettings.setMeter(code, { ...binding, devices })
  utilitySettings.setReading({
    billingMonth: billing.billingMonth.value,
    stall: bill.stall,
    contractId: bill.contractId || '',
    contractNumber: bill.contractNumber || '',
    tenant: bill.tenant,
    shopName: bill.shopName || '',
    waterMeterCode: bill.waterMeterCode || '',
    electricMeterCode: bill.electricMeterCode || '',
    prevWater: bill.prevWater,
    prevElectric: bill.prevElectric,
    currWater: bill.currWater,
    currElectric: bill.currElectric,
    status: bill.status,
    recordedAt: new Date().toISOString()
  })
  void utilitySettings.syncNow()
}
// แจ้งผลบนการ์ดของแผงถัดไป ให้รู้ว่าแผงก่อนหน้าถูกบันทึกหรือข้ามแล้ว และตอนนี้อยู่แผงใหม่ (หายเองใน 6 วินาที)
const meterNotice = ref('')
let meterNoticeTimer: ReturnType<typeof setTimeout> | undefined
const showMeterNotice = (text: string) => {
  meterNotice.value = text
  clearTimeout(meterNoticeTimer)
  meterNoticeTimer = setTimeout(() => { meterNotice.value = '' }, 6000)
}
watch(meterCardCode, (value) => { if (!value) { clearTimeout(meterNoticeTimer); meterNotice.value = '' } })
const meterLabelOf = (code: string) => meterProjectStalls.value.find(row => row.code.toUpperCase() === code.toUpperCase())?.label || `แผง ${code}`
const handleMeterSkipped = (code: string) => {
  const label = meterLabelOf(code)
  handleMeterRecorded(code, `ข้าม${label} แล้ว`)
}
const handleMeterRecorded = (code: string, noticeText?: string) => {
  const doneLabel = meterLabelOf(code)
  persistRecordedMeterReadings(code)
  const next = findNextMeterStall(code)
  if (next) { openMeterCard(next.code); showMeterNotice(noticeText || `บันทึก${doneLabel} แล้ว`); return }
  closeMeterCard()
  void showSuccess('จดครบทุกแผงของโครงการนี้แล้ว', 'จดครบแล้ว')
}
// จากตารางสรุปการจด: ปิดตาราง เข้าโครงการของแผงนั้น แล้วเปิดการ์ดจดแผง
const openMeterStall = async (code: string) => {
  meterSummaryOpen.value = false
  const upper = String(code).toUpperCase()
  const project = visibleProjects.value.find(item => getAllFloorLayoutsForProject(item).some(layout => layout.stalls.some((stall: any) => stallCodeOf(stall) === upper)))
  if (project && selectedProject.value?.id !== project.id) {
    if (typeof mapRef.value?.enterProject === 'function') mapRef.value.enterProject(project)
    else selectProject(project)
    await nextTick()
  }
  openMeterCard(upper)
}
const openMeterProject = (id: string) => {
  meterSummaryOpen.value = false
  closeMeterCard()
  const project = visibleProjects.value.find(item => String(item.id) === String(id))
  if (!project) return
  if (typeof mapRef.value?.enterProject === 'function') mapRef.value.enterProject(project)
  else selectProject(project)
}
// ยังไม่เชื่อมระบบออกใบแจ้งหนี้ (ข้อมูลจำลอง): แจ้งให้ชัดเจน ไม่ทำเหมือนส่งสำเร็จ
const sendMeterInvoice = () => {
  void showWarning('ฟังก์ชันส่งไปออกใบแจ้งหนี้ยังไม่เชื่อมต่อกับระบบออกใบแจ้งหนี้ ขณะนี้ใช้ข้อมูลจำลอง', 'ยังไม่พร้อมใช้งาน')
}
watch(pageMode, () => { meterCardCode.value = ''; meterSummaryOpen.value = false })
watch(() => selectedProject.value?.id, () => { meterCardCode.value = '' })

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
  if (meterCardCode.value) return '380px'
  if (isCostMode.value && costTool.value) return '640px'
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
  if (meterCardCode.value) { closeMeterCard(); return }
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

const bindManagePageListeners = () => {
  if (!import.meta.client) return
  window.removeEventListener('beforeunload', handleBeforeUnload)
  window.removeEventListener('keydown', onStallEsc)
  window.removeEventListener('click', handleClickOutsideDropdown)
  window.addEventListener('beforeunload', handleBeforeUnload)
  window.addEventListener('keydown', onStallEsc)
  window.addEventListener('click', handleClickOutsideDropdown)
}

const unbindManagePageListeners = () => {
  if (!import.meta.client) return
  window.removeEventListener('beforeunload', handleBeforeUnload)
  window.removeEventListener('keydown', onStallEsc)
  window.removeEventListener('click', handleClickOutsideDropdown)
}

onMounted(() => {
  isManagePageActive = true
  void loadData()
  loadAllCustomContracts()
  if (import.meta.client) {
    bindManagePageListeners()
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
  isManagePageActive = false
  loadDataController?.abort()
  loadDataController = null
  loadDataRun += 1
  unbindManagePageListeners()
})

onDeactivated(() => {
  isManagePageActive = false
  unbindManagePageListeners()
  if (loadDataController) {
    reloadAfterActivation = true
    loadDataController.abort()
    loadDataController = null
  }
})

onActivated(() => {
  isManagePageActive = true
  bindManagePageListeners()
  // KeepAlive also calls onActivated on the initial mount. loadData() already
  // starts in onMounted, so do not abort and restart that first request.
  if (!hasActivatedManagePage) {
    hasActivatedManagePage = true
    return
  }
  if (reloadAfterActivation || loadedMarketId !== marketId.value) {
    void loadData()
    return
  }
  if (isMapOverviewMode.value && import.meta.client) {
    nextTick(() => window.requestAnimationFrame(() => mapRef.value?.resize?.()))
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

watch(marketId, () => {
  if (!isManagePageActive || route.path !== '/stalls/map') {
    reloadAfterActivation = true
    return
  }
  selectedProject.value = null
  closeStall()
  void loadData()
})

watch(() => route.fullPath, () => {
  if (isManagePageActive && route.path === '/stalls/map' && isMapOverviewMode.value && import.meta.client) {
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
.gis-panel{padding:26px 28px!important}.gis-panel-heading{align-items:center;margin-bottom:22px!important}.gis-panel-heading h2{font-size:1.25rem!important;font-weight:800;color:var(--tx-0f172a);margin:0}.gis-panel-heading p{margin:6px 0 0!important}.gis-heading-right{display:flex;align-items:center;gap:12px;margin-left:auto;flex-wrap:wrap}.gis-map-toolbar{display:flex;flex:0 1 auto;min-width:0}.gis-map-toolbar :deep(.master-map-top-bar){width:auto;margin:0}.gis-map-toolbar :deep(.master-map-tools-right){margin-left:0;flex-wrap:nowrap}@media(max-width:700px){.gis-panel-heading{align-items:flex-start}.gis-heading-right,.gis-map-toolbar{width:100%;margin-left:0}.gis-map-toolbar :deep(.master-map-top-bar),.gis-map-toolbar :deep(.master-map-tools-right){width:100%;flex-wrap:wrap}.gis-map-toolbar :deep(.search-input-group){flex:1;width:auto}}
.staff-dashboard-intro{display:flex;align-items:center;justify-content:space-between;gap:16px;flex-wrap:wrap;margin:8px 0 22px}.dashboard-intro-left h1{margin:0;color:var(--tx-0f172a);font-size:1.5rem;font-weight:700;letter-spacing:-.01em;line-height:1.3}.dashboard-updated-sub{margin:4px 0 0;color:var(--tx-64748b);font-size:.875rem;font-weight:500}.dashboard-intro-actions{display:flex;align-items:center;gap:12px;flex-wrap:wrap}
.officer-summary-dashboard{display:grid;gap:20px;margin-bottom:12px}.dash-kpi-ribbon{margin:0!important}.dash-kpi-ribbon .map-summary-card{cursor:default}.dash-kpi-ribbon .map-summary-card:hover{transform:none;box-shadow:0 1px 3px rgba(15,23,42,.03)}.dash-two-col-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:20px;align-items:stretch}
.dash-system-card{background:var(--bg-ffffff);border:1px solid var(--bd-eef2f6);border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.03);padding:22px 24px;display:flex;flex-direction:column;gap:16px}
.dash-card-header{display:flex;align-items:center;justify-content:space-between;gap:14px;flex-wrap:wrap}.dash-card-title{margin:0;color:var(--tx-1e293b);font-size:1.05rem;font-weight:700;line-height:1.3}.dash-card-subtitle{margin:3px 0 0;color:var(--tx-64748b);font-size:.8125rem;font-weight:500}
.dash-status-badge{display:inline-flex;align-items:center;justify-content:center;padding:3px 10px;border-radius:7px;font-size:.75rem;font-weight:600;white-space:nowrap}.dash-status-badge.badge-amber{background:var(--bg-fef3c7);color:var(--tx-d97706)}.dash-status-badge.badge-green{background:var(--bg-dcfce7);color:var(--tx-15803d)}.dash-status-badge.badge-red{background:var(--bg-fee2e2);color:var(--tx-dc2626)}.dash-status-badge.badge-gray{background:var(--bg-f1f5f9);color:var(--tx-64748b)}
.dash-donut-layout{display:grid;grid-template-columns:148px minmax(0,1fr);gap:22px;align-items:center;margin:auto 0}
.dash-donut-wrap{position:relative;width:148px;height:148px;display:flex;align-items:center;justify-content:center;margin:0 auto}.dash-donut-svg{width:148px;height:148px;display:block}.dash-donut-center{position:absolute;inset:0;display:flex;flex-direction:column;align-items:center;justify-content:center;text-align:center;pointer-events:none}.dash-donut-center strong{color:var(--tx-0f172a);font-size:1.25rem;font-weight:800;line-height:1.1}.dash-donut-center span{color:var(--tx-64748b);font-size:.75rem;font-weight:600;margin-top:3px}
.dash-donut-legend-list{display:grid;gap:8px}.dash-legend-row{width:100%;display:flex;align-items:center;justify-content:space-between;gap:10px;padding:9px 12px;border:1px solid var(--bd-eef2f6);border-radius:10px;background:var(--bg-f8fafc)}.legend-left{display:inline-flex;align-items:center;gap:8px;color:var(--tx-334155);font-size:.875rem;font-weight:600}.legend-left i{width:10px;height:10px;border-radius:50%;display:inline-block}.dot-available{background:#10b981}.dot-reserved{background:#f59e0b}.dot-occupied{background:#475569}.legend-right{display:inline-flex;align-items:center;gap:8px}.legend-right strong{color:var(--tx-0f172a);font-size:.9375rem;font-weight:700}
.dash-mini-rev-row{display:grid;grid-template-columns:1fr 1fr;gap:10px;padding-top:12px;border-top:1px solid var(--bd-f1f5f9)}.dash-mini-rev-item{display:flex;align-items:center;justify-content:space-between;gap:8px;padding:10px 14px;border-radius:10px;background:var(--bg-f8fafc);border:1px solid var(--bd-eef2f6)}.dash-mini-rev-item span{color:var(--tx-64748b);font-size:.8125rem;font-weight:500}.dash-mini-rev-item strong{font-size:1rem;font-weight:700;color:var(--tx-0f172a)}.dash-mini-rev-item strong.is-blue{color:var(--tx-2563eb)}.dash-mini-rev-item strong.is-green{color:var(--tx-059669)}
.dash-bar-chart-list{display:grid;gap:14px;margin:auto 0}.dash-bar-row{display:grid;gap:6px}.dash-bar-top{display:flex;align-items:center;justify-content:space-between;gap:10px}.dash-bar-name{display:inline-flex;align-items:center;gap:8px;color:var(--tx-1e293b);font-size:.875rem;font-weight:600;min-width:0;max-width:62%;overflow:hidden;white-space:nowrap}.dash-bar-name>span{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.dash-bar-name i{width:10px;height:10px;border-radius:3px;flex-shrink:0}.dash-bar-meta{color:var(--tx-64748b);font-size:.8125rem;font-weight:500;white-space:nowrap}.dash-bar-meta strong{color:var(--tx-0f172a);font-weight:700}.dash-bar-track{height:10px;border-radius:999px;background:var(--bg-f1f5f9);overflow:hidden}.dash-bar-fill{display:block;height:100%;border-radius:999px;transition:width .3s ease}
@media(max-width:640px){.dash-donut-layout{grid-template-columns:1fr}.dash-mini-rev-row{grid-template-columns:1fr}}
.dash-table-card{padding:24px 28px;gap:16px}.dash-table-header{display:flex;align-items:center;justify-content:space-between;gap:16px;flex-wrap:wrap}.dash-table-header-actions{display:flex;align-items:center;gap:10px;flex-wrap:wrap}.dash-category-select{min-width:210px}
.dash-search-box{position:relative;display:inline-flex;align-items:center;min-width:240px}.dash-search-box svg{position:absolute;left:12px;color:var(--tx-94a3b8);pointer-events:none;z-index:1}.dash-search-box :deep(.app-input){min-height:40px;padding-left:36px;border-color:var(--bd-e2e8f0);border-radius:8px}
.dash-filter-tabs{display:flex;align-items:center;gap:24px;border-bottom:1px solid var(--bd-f1f5f9);overflow-x:auto}.dash-filter-tab{position:relative;padding:10px 4px 12px;border:0;background:transparent;color:var(--tx-64748b);font:inherit;font-size:.9375rem;font-weight:600;cursor:pointer;white-space:nowrap;transition:color .15s}.dash-filter-tab:hover{color:var(--tx-1e293b)}.dash-filter-tab.active{color:var(--tx-2563eb)}.dash-filter-tab.active::after{content:'';position:absolute;left:0;right:0;bottom:-1px;height:3px;border-radius:2px;background:#2563eb}
.dash-table-wrap{width:100%;overflow-x:auto}.dash-std-table{width:100%;border-collapse:collapse;text-align:left}.dash-std-table thead th{padding:14px 16px;color:var(--tx-64748b);font-size:.8125rem;font-weight:500;border-bottom:1px solid var(--bd-f1f5f9);white-space:nowrap;background:var(--bg-ffffff)}.dash-std-table tbody td{padding:16px;color:var(--tx-334155);font-size:.875rem;font-weight:500;border-bottom:1px solid var(--bd-f1f5f9);vertical-align:middle}.dash-std-table tbody tr:hover{background:var(--bg-f8fafc)}.dash-std-table .text-right{text-align:right}.dash-std-table .text-center{text-align:center}.dash-std-table strong{color:var(--tx-1e293b);font-weight:700}.dash-std-table .code-primary{color:var(--tx-2563eb)}
.zone-name-cell{display:flex;align-items:center;gap:12px}.zone-color-dot{width:12px;height:12px;border-radius:4px;flex-shrink:0}.zone-name-cell strong{display:block;color:var(--tx-1e293b);font-size:.875rem;font-weight:600}.zone-name-cell small,.applicant-cell small{display:block;margin-top:2px;color:var(--tx-64748b);font-size:.75rem;font-weight:500}.applicant-cell strong{display:block;color:var(--tx-1e293b);font-size:.875rem;font-weight:600}
.zone-occ-cell{display:flex;align-items:center;gap:10px}.zone-occ-bar{flex:1;height:8px;border-radius:999px;background:var(--bg-e2e8f0);overflow:hidden;min-width:80px}.zone-occ-bar span{display:block;height:100%;border-radius:999px;background:#2563eb}.zone-occ-cell strong{font-size:.8125rem;min-width:36px;text-align:right}
.dash-table-view-btn{display:inline-flex;align-items:center;gap:6px;padding:6px 14px;border:1px solid var(--bd-e2e8f0);border-radius:8px;background:var(--bg-ffffff);color:var(--tx-334155);font:inherit;font-size:.825rem;font-weight:600;text-decoration:none;cursor:pointer;white-space:nowrap;transition:background .15s,border-color .15s,color .15s}.dash-table-view-btn:hover{background:var(--bg-eff6ff);border-color:var(--bd-93c5fd);color:var(--tx-2563eb)}
.dash-empty-row{padding:32px 16px!important;text-align:center;color:var(--tx-64748b)!important;font-size:.875rem}
@media(max-width:1080px){.dash-two-col-grid{grid-template-columns:1fr}}@media(max-width:768px){.dash-revenue-grid{grid-template-columns:1fr}.dash-queue-item{flex-direction:column;align-items:flex-start}.dash-queue-right{width:100%;justify-content:space-between}}
.workspace-grid{display:grid;grid-template-columns:minmax(0,1fr)!important;margin-top:0}.workspace-grid[style*="display: none"]{display:none!important}.gis-panel-heading h2{font-size:1.15rem}.map-summary-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:18px;margin:4px 0 24px}.map-summary-card{background:var(--bg-ffffff);border:1px solid var(--bd-eef2f6);border-radius:14px;padding:16px 20px;display:flex;align-items:center;gap:16px;cursor:pointer;font-family:inherit;text-align:left;position:relative;transition:all .2s cubic-bezier(.4,0,.2,1);box-shadow:0 1px 3px rgba(15,23,42,.03)}.map-summary-card:hover{transform:translateY(-2px);box-shadow:0 6px 16px rgba(15,23,42,.07)}.map-summary-card:focus-visible{outline:2px solid #2563eb;outline-offset:2px}.map-summary-icon{width:46px;height:46px;border-radius:12px;display:flex;align-items:center;justify-content:center;flex-shrink:0;transition:all .2s ease}.map-summary-content{display:flex;flex-direction:column;min-width:0}.map-summary-title{font-size:.8125rem;color:var(--tx-64748b);font-weight:500;line-height:1.25;margin-bottom:6px;white-space:nowrap}.map-summary-value-row{display:flex;align-items:baseline;gap:8px}.map-summary-count{font-size:1.45rem;font-weight:800;color:var(--tx-0f172a);line-height:1.1}.map-summary-unit-pill{font-size:.78rem;font-weight:600;padding:2px 8px;border-radius:6px;line-height:1.2;font-style:normal}.summary-total .map-summary-icon{background:var(--bg-eff6ff);color:var(--tx-2563eb);border:1px solid var(--bd-dbeafe)}.summary-total .map-summary-unit-pill{background:var(--bg-eff6ff);color:var(--tx-1d4ed8)}.summary-total.active{border-color:#2563eb;background:var(--bg-f8faff);box-shadow:0 0 0 2px rgba(37,99,235,.18)}.summary-available .map-summary-icon{background:var(--bg-ecfdf5);color:var(--tx-059669);border:1px solid var(--bd-d1fae5)}.summary-available .map-summary-unit-pill{background:var(--bg-dcfce7);color:var(--tx-15803d)}.summary-available.active{border-color:#10b981;background:var(--bg-f0fdf4);box-shadow:0 0 0 2px rgba(16,185,129,.2)}.summary-reserved .map-summary-icon{background:var(--bg-fffbeb);color:var(--tx-d97706);border:1px solid var(--bd-fef3c7)}.summary-reserved .map-summary-unit-pill{background:var(--bg-fef3c7);color:var(--tx-b45309)}.summary-reserved.active{border-color:#f59e0b;background:var(--bg-fffdf7);box-shadow:0 0 0 2px rgba(245,158,11,.2)}.summary-occupied .map-summary-icon{background:var(--bg-f1f5f9);color:var(--tx-64748b);border:1px solid var(--bd-e2e8f0)}.summary-occupied .map-summary-unit-pill{background:var(--bg-f1f5f9);color:var(--tx-64748b)}.summary-occupied.active{border-color:var(--bd-94a3b8);background:var(--bg-f8fafc);box-shadow:0 0 0 2px rgba(100,116,139,.18)}@media(max-width:800px){.map-summary-grid{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:440px){.map-summary-grid{grid-template-columns:1fr}}.header-actions{display:inline-flex;align-items:center;gap:10px}.draft-status{display:inline-flex;align-items:center;gap:7px;height:34px;box-sizing:border-box;padding:0 12px;border-radius:999px;background:var(--bg-f1f5f9);color:var(--tx-475569);font-size:.8125rem;font-weight:600;transition:all .2s ease}.draft-status i{width:8px;height:8px;border-radius:50%;background:currentColor}.draft-status.published{background:var(--bg-ecfdf5);color:var(--tx-15803d);border:1px solid var(--bd-bbf7d0)}.draft-status.is-dirty{background:var(--bg-fff7ed);color:var(--tx-c2410c);border:1px solid var(--bd-fed7aa)}.draft-status.is-dirty i{background:#ea580c;box-shadow:0 0 0 3px rgba(234,88,12,.18)}.header-status-skeleton,.header-action-skeleton{display:block;border-radius:999px;background:linear-gradient(90deg,var(--bg-eef2f7) 25%,var(--bg-f8fafc) 50%,var(--bg-eef2f7) 75%);background-size:200% 100%;animation:headerSkeletonShimmer 1.35s ease-in-out infinite}.header-status-skeleton{width:102px;height:34px}.header-action-skeleton{width:128px;height:34px;border-radius:7px}@keyframes headerSkeletonShimmer{0%{background-position:100% 0}100%{background-position:-100% 0}}@media(prefers-reduced-motion:reduce){.header-status-skeleton,.header-action-skeleton{animation:none}}
.unsaved-leave-backdrop{position:fixed;inset:0;background:rgba(15,23,42,.65);backdrop-filter:blur(8px);-webkit-backdrop-filter:blur(8px);display:flex;align-items:center;justify-content:center;padding:20px;z-index:99999}.unsaved-leave-card{background:var(--bg-ffffff);border-radius:20px;width:100%;max-width:480px;overflow:hidden;box-shadow:0 24px 48px -12px rgba(15,23,42,.28),0 0 0 1px rgba(226,232,240,.8);position:relative}.unsaved-leave-content{padding:24px 28px}.unsaved-leave-content p{font-size:.9375rem;color:var(--tx-475569);line-height:1.55;margin:0}.unsaved-leave-actions{display:grid;grid-template-columns:1fr 1fr;gap:10px;padding:16px 28px;border-top:1px solid var(--bd-e2e8f0);background:var(--bg-f8fafc)}
.inspect-selected-stall-drawer{position:absolute;top:62px;right:14px;width:min(380px,calc(100vw - 28px));max-height:calc(100% - 76px);display:flex;flex-direction:column;pointer-events:auto;z-index:125;background:var(--bg-ffffff);border:1px solid var(--bd-e2e8f0);border-radius:18px;box-shadow:0 16px 38px rgba(15,23,42,.16);overflow:hidden;contain:layout paint;will-change:width,transform,opacity;transition:width .46s cubic-bezier(.16,1,.3,1),box-shadow .36s ease,background-color .28s ease}
.inspect-selected-stall-drawer.is-edit-mode{width:min(440px,calc(100% - 28px))}
.inspect-selected-stall-drawer.is-full-detail-mode{width:min(760px,calc(100vw - 36px));z-index:130}
@media(min-width:769px){.inspect-selected-stall-drawer.is-management-mode{bottom:100px;max-height:calc(100% - 176px)}}
.issd-header{padding:16px 18px 14px;border-bottom:1px solid var(--bd-f1f5f9);background:var(--bg-ffffff);flex-shrink:0}.issd-header-top{display:flex;align-items:center;justify-content:space-between;gap:10px}.issd-title-group{display:flex;align-items:center;gap:8px;flex-wrap:wrap;min-width:0}.issd-title-group strong{font-size:1rem;font-weight:700;color:var(--tx-0f172a)}.issd-project-sub{margin-top:4px;font-size:.8125rem;font-weight:500;color:var(--tx-64748b);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.issd-back-btn{width:28px;height:28px;border-radius:8px;border:1px solid var(--bd-e2e8f0);background:var(--bg-ffffff);color:var(--tx-475569);display:inline-flex;align-items:center;justify-content:center;cursor:pointer;flex-shrink:0;transition:background .15s,border-color .15s,color .15s}.issd-back-btn:hover{background:var(--bg-f1f5f9);color:var(--tx-0f172a);border-color:var(--bd-cbd5e1)}
.issd-close-btn{width:28px;height:28px;border-radius:8px;border:1px solid var(--bd-e2e8f0);background:var(--bg-ffffff);color:var(--tx-475569);display:inline-flex;align-items:center;justify-content:center;cursor:pointer;flex-shrink:0;transition:background .15s,border-color .15s,color .15s}.issd-close-btn:hover{background:var(--bg-fef2f2);border-color:var(--bd-fca5a5);color:var(--tx-ef4444)}
.issd-detail-tabs-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  background: var(--bg-f8fafc);
  border-bottom: 1px solid var(--bd-e2e8f0);
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
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  color: var(--tx-475569);
  font-family: inherit;
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.03);
}
.issd-tab-btn:hover {
  background: var(--bg-f1f5f9);
  color: var(--tx-0f172a);
  border-color: var(--bd-cbd5e1);
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
  background: var(--bg-f1f5f9);
  color: var(--tx-475569);
  font-size: 0.75rem;
  font-weight: 700;
  line-height: 1;
}
.issd-tab-btn.active .issd-tab-count {
  background: var(--bg-rgba-255-255-255-0p28);
  color: #ffffff;
}
.issd-tab-count.is-urgent {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
}
.issd-tab-btn.has-urgent-badge {
  border-color: #fcd34d;
  background: var(--bg-fffbeb);
  color: var(--tx-b45309);
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

.issd-kpi-card{padding:10px 12px;border-radius:11px;border:1px solid var(--bd-e2e8f0);background:var(--bg-f8fafc);display:grid;gap:2px}
.issd-kpi-card span{color:var(--tx-64748b);font-size:.75rem;font-weight:500}
.issd-kpi-card strong{color:var(--tx-0f172a);font-size:1.15rem;font-weight:700;line-height:1.2}
.issd-kpi-card.is-green strong{color:var(--tx-10b981)}
.issd-kpi-card em{font-style:normal;color:var(--tx-64748b);font-size:.75rem;font-weight:500}
.issd-fact-section{display:flex;flex-direction:column;gap:8px;flex-shrink:0}
.issd-section-title{font-size:.8125rem;font-weight:700;color:var(--tx-1e293b)}
.issd-fact-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:8px;margin:0}
.issd-fact-cell{padding:9px 12px;border:1px solid var(--bd-e2e8f0);border-radius:9px;background:var(--bg-ffffff)}
.issd-fact-cell.full{grid-column:1/-1}
.issd-fact-cell dt{color:var(--tx-64748b);font-size:.75rem;font-weight:500}
.issd-fact-cell dd{margin:3px 0 0;color:var(--tx-1e293b);font-size:.8125rem;font-weight:600;word-break:break-word}
.issd-ops-snapshot{display:grid;gap:10px;padding:13px;border:1px solid var(--bd-dbeafe);border-radius:12px;background:var(--bg-f8fbff);flex-shrink:0}
.issd-ops-heading{display:flex;align-items:flex-start;justify-content:space-between;gap:12px}.issd-ops-heading>div{display:grid;gap:2px}.issd-ops-heading strong{font-size:.9375rem;color:var(--tx-0f172a)}.issd-ops-heading span:not(.issd-status-pill){font-size:.75rem;color:var(--tx-64748b)}
.issd-ops-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:1px;background:var(--bg-e2e8f0);border:1px solid var(--bd-e2e8f0);border-radius:9px;overflow:hidden}.issd-ops-grid>div{display:grid;gap:2px;padding:9px 10px;background:var(--bg-ffffff)}.issd-ops-grid span{font-size:.75rem;color:var(--tx-64748b)}.issd-ops-grid strong{font-size:.8125rem;color:var(--tx-1e293b);line-height:1.35}.issd-ops-grid strong.is-alert{color:var(--tx-dc2626)}
.issd-action-summary{display:grid;gap:5px;padding:9px 10px;border-radius:9px;background:var(--bg-fffbeb);border:1px solid var(--bd-fde68a)}.issd-action-summary>strong{font-size:.8125rem;color:var(--tx-92400e)}.issd-action-summary ul{display:grid;gap:3px;margin:0;padding-left:18px;color:var(--tx-92400e);font-size:.75rem;line-height:1.45}

/* Action Task Cards in "สิ่งที่ต้องดำเนินการ" Tab */
.issd-action-cards-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.issd-action-task-card {
  padding: 13px 15px;
  border-radius: 12px;
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  display: flex;
  flex-direction: column;
  gap: 6px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}
.issd-action-task-card.is-danger {
  border-color: var(--bd-fecaca);
  background: var(--bg-fffafa);
}
.issd-action-task-card.is-warning {
  border-color: var(--bd-fde68a);
  background: var(--bg-fffdf7);
}
.issd-action-task-card.is-info {
  border-color: var(--bd-bfdbfe);
  background: var(--bg-f8fbff);
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
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
}
.issd-action-cat-badge.is-warning {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
}
.issd-action-cat-badge.is-info {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}
.issd-action-priority-text {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--tx-64748b);
}
.issd-action-task-title {
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--tx-0f172a);
}
.issd-action-task-desc {
  margin: 0;
  font-size: 0.8125rem;
  color: var(--tx-475569);
  line-height: 1.45;
}
.issd-action-task-footer {
  display: flex;
  justify-content: flex-end;
  padding-top: 6px;
  border-top: 1px dashed var(--bd-e2e8f0);
}
.issd-empty-box.is-clean-state {
  background: var(--bg-f0fdf4);
  border-color: var(--bd-bbf7d0);
}
.issd-empty-box.is-clean-state svg {
  color: var(--tx-16a34a);
}
.issd-empty-box.is-clean-state p {
  color: var(--tx-15803d);
  font-weight: 700;
}
.issd-empty-sub {
  font-size: 0.75rem;
  color: var(--tx-64748b);
}

.issd-booking-list,.issd-contract-list{display:flex;flex-direction:column;gap:10px}
.issd-booking-card,.issd-contract-card{padding:14px 16px;border:1px solid var(--bd-e2e8f0);border-radius:12px;background:var(--bg-ffffff);display:flex;flex-direction:column;gap:10px;flex-shrink:0;transition:border-color .15s ease,box-shadow .15s ease}
.issd-contract-card.is-current{border-color:var(--bd-93c5fd);box-shadow:inset 3px 0 0 #2563eb}.issd-current-label{display:inline-flex;padding:2px 7px;border-radius:6px;background:var(--bg-eff6ff);color:var(--tx-2563eb);font-size:.75rem;font-weight:600}
.issd-booking-card:hover,.issd-contract-card:hover{border-color:var(--bd-cbd5e1);box-shadow:0 3px 10px rgba(15,23,42,.05)}
.issd-card-top-row{display:flex;align-items:center;justify-content:space-between;gap:8px}
.issd-card-type-tag{display:inline-flex;align-items:center;gap:5px;font-size:.75rem;font-weight:600;color:var(--tx-2563eb);background:var(--bg-eff6ff);padding:3px 9px;border-radius:6px}
.issd-card-type-tag.is-application{color:var(--tx-7c3aed);background:var(--bg-f5f3ff)}
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
  font-size: 0.8125rem;
  font-weight: 700;
  color: var(--tx-334155);
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
  background: var(--bg-f1f5f9);
  border: 1px solid var(--bd-e2e8f0);
  color: var(--tx-475569);
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
  color: var(--tx-64748b);
  font-weight: 500;
  line-height: 1.2;
}
.issd-booking-phone {
  font-size: 0.75rem;
  color: var(--tx-64748b);
  font-weight: 500;
}
.issd-booking-footer-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 9px;
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-f1f5f9);
  flex-wrap: wrap;
}
.issd-booking-info-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  color: var(--tx-64748b);
}
.issd-booking-info-item svg {
  color: var(--tx-64748b);
  flex-shrink: 0;
}
.issd-booking-info-item strong {
  color: var(--tx-1e293b);
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
  border: 1px solid var(--bd-bfdbfe);
  background: var(--bg-ffffff);
  color: var(--tx-1d4ed8);
  flex-shrink: 0;
  transition: all 0.15s ease;
}
.issd-booking-direct-btn:hover {
  background: #2563eb;
  color: #ffffff;
  border-color: #2563eb;
}
.issd-booking-direct-btn.is-application {
  border-color: var(--bd-ddd6fe);
  background: var(--bg-f5f3ff);
  color: var(--tx-6d28d9);
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
  border-top: 1px dashed var(--bd-e2e8f0);
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
  background: var(--bg-f5f3ff);
  border: 1px solid var(--bd-ddd6fe);
  color: var(--tx-6d28d9);
}
.issd-occupant-action-btn.is-app:hover {
  background: var(--bg-ede9fe);
  border-color: var(--bd-c4b5fd);
}
.issd-occupant-action-btn.is-booking {
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-bfdbfe);
  color: var(--tx-1d4ed8);
}
.issd-occupant-action-btn.is-booking:hover {
  background: var(--bg-dbeafe);
  border-color: var(--bd-93c5fd);
}

/* ==========================================================================
   Full Detail Action Alert Card (Modern High-End Task Card)
   ========================================================================== */
.issd-action-alert-card {
  margin: 10px 14px 4px;
  background: var(--bg-fffbeb);
  border: 1px solid var(--bd-fde68a);
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
  background: var(--bg-fffbeb);
  cursor: pointer;
  user-select: none;
  transition: background 0.15s ease;
}

.issd-aac-header:hover {
  background: var(--bg-fef3c7);
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
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
  border: 1px solid var(--bd-fde68a);
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
  color: var(--tx-92400e);
}

.issd-aac-count-badge {
  display: inline-flex;
  align-items: center;
  padding: 1px 7px;
  border-radius: 999px;
  background: #fde68a;
  color: var(--tx-78350f);
  font-size: 0.75rem;
  font-weight: 700;
}

.issd-aac-sub {
  font-size: 0.75rem;
  color: var(--tx-b45309);
}

.issd-aac-toggle-btn {
  width: 26px;
  height: 26px;
  border-radius: 6px;
  border: 1px solid rgba(217, 119, 6, 0.25);
  background: var(--bg-ffffff);
  color: var(--tx-92400e);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.issd-aac-toggle-btn:hover {
  background: var(--bg-fef3c7);
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
  background: var(--bg-ffffff);
}

.issd-aac-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  border-top: 1px solid var(--bd-f1f5f9);
  cursor: pointer;
  transition: background 0.15s ease;
}

.issd-aac-row:hover {
  background: var(--bg-f8fafc);
}

.issd-aac-row:hover .issd-aac-action-text {
  color: var(--tx-1d4ed8);
}

.issd-aac-row:hover .issd-aac-row-right svg {
  transform: translateX(2px);
}

.issd-aac-row.is-danger:hover {
  background: var(--bg-fef2f2);
}

.issd-aac-row.is-warning:hover {
  background: var(--bg-fffbeb);
}

.issd-aac-row.is-info:hover {
  background: var(--bg-eff6ff);
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
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
  border: 1px solid var(--bd-fecaca);
}

.issd-aac-tag.is-warning {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
  border: 1px solid var(--bd-fde68a);
}

.issd-aac-tag.is-info {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  border: 1px solid var(--bd-bfdbfe);
}

.issd-aac-row-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--tx-1e293b);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.issd-aac-row-right {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  flex-shrink: 0;
  color: var(--tx-2563eb);
  font-size: 0.78125rem;
  font-weight: 600;
}

.issd-aac-row-right svg {
  transition: transform 0.15s ease;
}

.issd-aac-action-text {
  white-space: nowrap;
}
.issd-card-title,.issd-contract-num{font-size:.875rem;font-weight:700;color:var(--tx-0f172a)}
.issd-contract-tenant{font-size:.8125rem;color:var(--tx-475569)}
.issd-contract-tenant strong{color:var(--tx-1e293b);font-weight:600}
.issd-card-meta-row{display:flex;align-items:center;justify-content:space-between;gap:10px;font-size:.75rem;color:var(--tx-64748b);flex-wrap:wrap}
.issd-card-meta-row .is-rent{color:var(--tx-10b981);font-weight:700;font-size:.8125rem}
.issd-timeline-list{display:flex;flex-direction:column;padding:4px 4px 10px}
.issd-timeline-panel{padding:12px 14px 4px;border:1px solid var(--bd-e2e8f0);border-radius:12px;background:var(--bg-ffffff)}
.issd-timeline-item{position:relative;display:flex;gap:12px;padding-bottom:16px}
.issd-timeline-item:not(:last-child)::before{content:'';position:absolute;left:5px;top:13px;bottom:0;width:1px;background:var(--bg-e2e8f0)}
.issd-timeline-dot{position:relative;z-index:1;width:11px;height:11px;margin-top:4px;flex:0 0 11px;border-radius:50%;background:#3b82f6;box-shadow:0 0 0 3px #dbeafe}
.issd-timeline-dot.is-green{background:#10b981;box-shadow:0 0 0 3px #dcfce7}
.issd-timeline-dot.is-amber{background:#f59e0b;box-shadow:0 0 0 3px #fef3c7}
.issd-timeline-dot.is-gray{background:#64748b;box-shadow:0 0 0 3px #f1f5f9}
.issd-timeline-copy{flex:1;min-width:0}
.issd-timeline-head-row{display:flex;align-items:center;justify-content:space-between;gap:8px;flex-wrap:wrap}
.issd-timeline-cat-pill{display:inline-flex;padding:2px 8px;border-radius:6px;background:var(--bg-f1f5f9);color:var(--tx-475569);font-size:.75rem;font-weight:600}
.issd-timeline-copy strong{display:block;font-size:.8125rem;color:var(--tx-1e293b);font-weight:600}
.issd-timeline-copy time{display:block;color:var(--tx-64748b);font-size:.75rem;margin-top:2px}
.issd-timeline-copy p{margin:6px 0 0;padding:7px 10px;border-radius:8px;background:var(--bg-f8fafc);border:1px solid var(--bd-f1f5f9);color:var(--tx-475569);font-size:.75rem;line-height:1.45}
.issd-empty-box{display:flex;flex-direction:column;align-items:center;justify-content:center;gap:8px;padding:32px 16px;text-align:center;color:var(--tx-94a3b8);background:var(--bg-f8fafc);border-radius:12px;border:1px dashed var(--bd-e2e8f0)}
.issd-empty-box svg{color:var(--tx-cbd5e1)}
.issd-empty-box p{margin:0;font-size:.8125rem;color:var(--tx-64748b)}
.issd-body-edit{padding:12px 14px}
.issd-form-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:10px;margin:0}
.issd-form-grid .full{grid-column:1/-1}
.stall-edit-form{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px 14px}.stall-edit-form label{display:grid;gap:5px}.stall-edit-delete{margin-right:auto}@media(max-width:560px){.stall-edit-form{grid-template-columns:1fr}}
.issd-form-section{margin:8px 0 0;padding-top:12px;border-top:1px solid var(--bd-f1f5f9);color:var(--tx-0f172a);font-size:var(--font-size-sm);font-weight:700}
.issd-form-grid>.issd-form-section:first-child{margin-top:0;padding-top:0;border-top:0}
.issd-form-grid label{display:grid;gap:5px;color:var(--tx-475569);font-size:.8125rem;font-weight:600}
.issd-form-grid textarea{width:100%;box-sizing:border-box;border:1px solid var(--bd-cbd5e1);border-radius:10px;padding:9px 14px;color:var(--tx-1e293b);background:var(--bg-ffffff);font:500 .88rem/1.4 inherit;resize:vertical;min-height:64px;transition:border-color .15s,box-shadow .15s}
.issd-form-grid textarea::placeholder{color:var(--tx-94a3b8)}
.issd-form-grid textarea:focus{outline:none;border-color:#2563eb;box-shadow:0 0 0 3px rgba(37,99,235,.12)}
.issd-footer-edit{display:flex;align-items:center;justify-content:space-between;gap:8px}
.issd-footer-right{display:flex;align-items:center;gap:8px;margin-left:auto}
@media(max-width:480px){.issd-form-grid{grid-template-columns:1fr}}
.issd-status-pill{display:inline-flex;align-items:center;padding:3px 10px;border-radius:999px;font-size:.75rem;font-weight:600;border:1px solid transparent;white-space:nowrap}.issd-status-pill.is-available{background:var(--bg-dcfce7);color:var(--tx-15803d);border-color:var(--bd-bbf7d0)}.issd-status-pill.is-reserved{background:var(--bg-fef3c7);color:var(--tx-d97706);border-color:var(--bd-fde68a)}.issd-status-pill.is-info{background:var(--bg-eff6ff);color:var(--tx-2563eb);border-color:var(--bd-bfdbfe)}.issd-status-pill.is-occupied{background:var(--bg-f1f5f9);color:var(--tx-475569);border-color:var(--bd-cbd5e1)}
.issd-photo-banner{position:relative;height:124px;flex-shrink:0;border-radius:12px;overflow:hidden;border:1px solid var(--bd-e2e8f0);background:#0f172a;cursor:pointer}.issd-photo-banner img{width:100%;height:100%;object-fit:cover;display:block;transition:transform .25s ease}.issd-photo-banner:hover img{transform:scale(1.05)}.issd-photo-count-badge{position:absolute;right:10px;bottom:10px;display:inline-flex;align-items:center;gap:5px;padding:3px 9px;border-radius:8px;background:rgba(15,23,42,.82);backdrop-filter:blur(4px);color:#fff;font-size:.75rem;font-weight:600}
.issd-detail-list{margin:0;flex-shrink:0;border:1px solid var(--bd-e2e8f0);border-radius:12px;background:var(--bg-ffffff);overflow:hidden}.issd-detail-row{display:flex;align-items:center;justify-content:space-between;gap:14px;padding:10px 14px;font-size:.8125rem}.issd-detail-row+.issd-detail-row{border-top:1px solid var(--bd-f1f5f9)}.issd-detail-row dt{color:var(--tx-64748b);font-size:.8125rem;font-weight:500;flex-shrink:0;white-space:nowrap}.issd-detail-row dd{margin:0;color:var(--tx-1e293b);font-size:.8125rem;font-weight:600;text-align:right;max-width:65%;line-height:1.4}.issd-detail-row dd.is-rent{color:var(--tx-15803d);font-weight:700}.issd-detail-row dd.is-bold{color:var(--tx-0f172a);font-weight:700}
.issd-occupancy-box{padding:11px 14px;flex-shrink:0;border-radius:12px;border:1px solid var(--bd-e2e8f0);background:var(--bg-f8fafc);display:grid;gap:3px}.issd-occupancy-box.is-available{background:var(--bg-f0fdf4);border-color:var(--bd-bbf7d0)}.issd-occupancy-box.is-reserved{background:var(--bg-fffbeb);border-color:var(--bd-fde68a)}.issd-occupancy-box.is-occupied{background:var(--bg-f8fafc);border-color:var(--bd-cbd5e1)}.issd-occupancy-head{display:flex;align-items:center;justify-content:space-between;gap:8px}.issd-occupancy-head span{font-size:.75rem;font-weight:600;color:var(--tx-64748b)}.issd-occupancy-head small{font-size:.75rem;font-weight:700;color:var(--tx-2563eb)}.issd-occupancy-box strong{font-size:.8125rem;font-weight:700;color:var(--tx-0f172a);line-height:1.35}.issd-occupancy-box p{margin:0;font-size:.75rem;color:var(--tx-475569);line-height:1.35}
.issd-footer{padding:12px 18px 16px;background:var(--bg-ffffff);border-top:1px solid var(--bd-f1f5f9);box-shadow:0 -4px 12px rgba(15,23,42,.04);flex-shrink:0}.issd-cta-btn{width:100%;box-sizing:border-box;padding:11px 16px;border:0;border-radius:11px;background:linear-gradient(180deg,#3b82f6 0%,#2563eb 100%);color:#fff;font-family:inherit;font-size:.875rem;font-weight:700;text-decoration:none;display:inline-flex;align-items:center;justify-content:center;gap:8px;cursor:pointer;box-shadow:0 4px 14px rgba(37,99,235,.32);transition:transform .15s ease,box-shadow .15s ease}.issd-cta-btn:hover{transform:translateY(-1px);box-shadow:0 6px 18px rgba(37,99,235,.4)}
.issd-metrics3{display:grid;grid-template-columns:1fr 1fr 1fr;margin-top:10px;padding:10px 6px;text-align:center;background:var(--bg-f8fafc);border:1px solid var(--bd-e2e8f0);border-radius:10px}.issd-metrics3>div{display:flex;flex-direction:column;gap:2px;min-width:0}.issd-metrics3>div+div{border-left:1px solid var(--bd-e2e8f0)}.issd-metric-label{font-size:var(--font-size-xs);font-weight:600;color:var(--tx-64748b)}.issd-metric-value{font-size:var(--font-size-lg);font-weight:800;color:var(--tx-0f172a);line-height:1.25}.issd-metric-value.is-rent{color:var(--tx-10b981)}.issd-metric-value small{display:block;font-size:var(--font-size-xs);font-weight:500;color:var(--tx-64748b)}.issd-bookable-pill{display:inline-block;padding:3px 10px;border-radius:6px;font-size:var(--font-size-xs);font-weight:600;background:var(--bg-dcfce7);color:var(--tx-15803d)}.issd-bookable-pill.is-paused{background:var(--bg-f1f5f9);color:var(--tx-64748b)}.issd-quick-note{margin-top:10px;padding:8px 12px;border-radius:8px;font-size:var(--font-size-xs);font-weight:600}.issd-quick-note.is-warn{background:var(--bg-fef3c7);color:var(--tx-b45309)}
.issd-footer-full{display:flex;align-items:center;justify-content:space-between;gap:10px}
.issd-footer-full :deep(.app-button),
.issd-footer-full > button,
.issd-footer-full > a {
  flex: 1;
}
.inspect-photo-Modal,.inspect-photo-modal{width:min(820px,94vw);background:#0f172a;border-radius:18px;overflow:hidden;box-shadow:0 30px 60px -20px rgba(0,0,0,.6);display:flex;flex-direction:column}.ipm-head{display:flex;align-items:center;justify-content:space-between;padding:14px 18px;color:#fff}.ipm-head strong{display:block;font-size:1rem;font-weight:700}.ipm-head span{display:block;font-size:.8125rem;color:var(--tx-94a3b8);margin-top:2px}.ipm-close{width:34px;height:34px;border-radius:10px;border:0;background:var(--bg-rgba-255-255-255-0p12);color:#fff;display:inline-flex;align-items:center;justify-content:center;cursor:pointer}.ipm-close svg{width:18px;height:18px}.ipm-stage{position:relative;height:min(480px,58vh);background:#000;display:flex;align-items:center;justify-content:center}.ipm-main-img{width:100%;height:100%;object-fit:contain}.ipm-nav{position:absolute;top:50%;transform:translateY(-50%);width:40px;height:40px;border-radius:50%;border:0;background:var(--bg-rgba-255-255-255-0p9);color:var(--tx-0f172a);font-size:1.35rem;font-weight:700;cursor:pointer;display:grid;place-items:center}.ipm-nav.is-prev{left:14px}.ipm-nav.is-next{right:14px}.ipm-counter{position:absolute;bottom:12px;left:50%;transform:translateX(-50%);background:rgba(15,23,42,.78);color:#fff;font-size:.75rem;font-weight:600;padding:3px 12px;border-radius:999px}.ipm-thumbs{display:flex;justify-content:center;gap:8px;padding:12px 18px 16px}.ipm-thumb{width:76px;height:52px;border-radius:8px;overflow:hidden;border:2px solid transparent;padding:0;background:#1e293b;cursor:pointer;opacity:.6}.ipm-thumb.active{border-color:#60a5fa;opacity:1}.ipm-thumb img{width:100%;height:100%;object-fit:cover;display:block}
@media(max-width:768px){.inspect-selected-stall-drawer{position:fixed;left:0;right:0;bottom:0;top:58px;width:100%!important;height:calc(100vh - 58px);max-height:calc(100vh - 58px);border-radius:0;z-index:1000}.inspect-selected-stall-drawer.is-full-detail-mode{width:100%!important;top:58px;height:calc(100vh - 58px);max-height:calc(100vh - 58px);border-radius:0}}
.map-edit-status{display:flex;gap:10px;align-items:flex-start;padding:13px;background:var(--bg-eff6ff);border:1px solid var(--bd-bfdbfe);border-radius:10px;color:var(--tx-1e40af)}.map-edit-status strong,.map-edit-status small{display:block}.map-edit-status strong{font-size:.875rem}.map-edit-status small{margin-top:3px;font-size:.75rem;color:var(--tx-476587)}.map-edit-dot{width:10px;height:10px;flex:0 0 10px;border-radius:50%;background:#3b5bff;margin-top:5px;box-shadow:0 0 0 4px #dbeafe}
.zone-form,.facility-form{display:grid;grid-template-columns:1fr 1fr;gap:13px}.zone-form label,.facility-form label{display:grid;align-content:start;gap:6px;color:var(--tx-475569);font-size:.8125rem;font-weight:600}.zone-form .full,.facility-form .full{grid-column:1/-1}.req-star{color:var(--tx-ef4444);font-weight:700}.stall-detail-summary{display:grid;gap:5px;padding:16px;margin:20px 0 14px;border:1px solid var(--bd-dbeafe);border-radius:12px;background:var(--bg-f8fbff)}.stall-detail-summary strong{color:var(--tx-1e293b);font-size:1rem}.stall-detail-summary>span:last-child{color:var(--tx-64748b);font-size:.875rem}.stall-detail-summary .soft-badge{justify-self:start}.status-available{background:var(--bg-dcfce7);color:var(--tx-15803d)!important}.status-reserved{background:var(--bg-fef3c7);color:var(--tx-d97706)!important}.status-occupied{background:var(--bg-f1f5f9);color:var(--tx-64748b)!important}.stall-detail-tabs{display:flex!important;width:100%;margin:0 0 14px}.stall-detail-tabs :deep(.app-segmented-option){flex:1}.stall-detail-grid{display:grid;grid-template-columns:1fr 1fr;gap:10px;margin:0}.stall-detail-grid>div{padding:12px 14px;border:1px solid var(--bd-e2e8f0);border-radius:10px;background:var(--bg-ffffff)}.stall-detail-grid .full{grid-column:1/-1}.stall-detail-grid dt{color:var(--tx-64748b);font-size:.8125rem;font-weight:500}.stall-detail-grid dd{margin:4px 0 0;color:var(--tx-1e293b);font-size:.9375rem;font-weight:600;word-break:break-word}.detail-section.is-tab-content{margin-top:0}.detail-section-heading{display:flex;align-items:center;justify-content:space-between;gap:12px}.detail-section-heading h3{margin:0}
.issd-tab-header-row {
  display: flex;
  flex-direction: column;
  gap: 0;
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
  color: var(--tx-64748b);
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
  color: var(--tx-0f172a);
}
.issd-tab-subtitle {
  display: block;
  font-size: 0.75rem;
  color: var(--tx-64748b);
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
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  background: var(--bg-ffffff);
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
  transition: border-color .15s ease, box-shadow .15s ease;
}
.issd-subtenant-card:hover,
.issd-payment-card:hover,
.issd-maintenance-card:hover {
  border-color: var(--bd-cbd5e1);
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}
.issd-contract-context{display:grid;grid-template-columns:auto 1fr;gap:2px 10px;align-items:center;padding:10px 12px;border:1px solid var(--bd-dbeafe);border-radius:10px;background:var(--bg-f8fbff);flex-shrink:0}.issd-contract-context span{font-size:.75rem;color:var(--tx-64748b)}.issd-contract-context strong{font-size:.8125rem;color:var(--tx-2563eb)}.issd-contract-context small{grid-column:2;font-size:.75rem;color:var(--tx-475569)}
.issd-finance-summary{display:grid;grid-template-columns:repeat(3,1fr);gap:8px;flex-shrink:0}.issd-finance-summary>div{display:grid;gap:2px;padding:9px 10px;border:1px solid var(--bd-e2e8f0);border-radius:9px;background:var(--bg-ffffff)}.issd-finance-summary span{font-size:.75rem;color:var(--tx-64748b)}.issd-finance-summary strong{font-size:.8125rem;color:var(--tx-1e293b)}.issd-finance-summary .is-success{color:var(--tx-15803d)}.issd-finance-summary .is-alert{color:var(--tx-dc2626)}
.issd-source-filter{display:flex;align-items:center;gap:6px;padding-bottom:2px;overflow-x:auto;flex-shrink:0}.issd-source-filter button{min-height:34px;padding:6px 12px;border:1px solid var(--bd-e2e8f0);border-radius:8px;background:var(--bg-ffffff);color:var(--tx-64748b);font:inherit;font-size:.8125rem;font-weight:600;white-space:nowrap;cursor:pointer;flex-shrink:0;transition:all .15s ease}.issd-source-filter button:hover{background:var(--bg-f8fafc);color:var(--tx-1e293b)}.issd-source-filter button.active{border-color:var(--bd-93c5fd);background:var(--bg-eff6ff);color:var(--tx-2563eb)}.issd-maintenance-tags{display:flex;align-items:center;gap:6px;flex-wrap:wrap}.issd-source-pill{display:inline-flex;padding:2px 8px;border-radius:6px;background:var(--bg-f1f5f9);color:var(--tx-475569);font-size:.75rem;font-weight:600}

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
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-f1f5f9);
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
  color: var(--tx-94a3b8);
  background: var(--bg-f8fafc);
}
.issd-subtenant-text-block {
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.issd-subtenant-idcard {
  font-size: 0.75rem;
  color: var(--tx-64748b);
  margin-top: 1px;
}

.issd-subtenant-details {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 0.8125rem;
  color: var(--tx-475569);
}
.issd-detail-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--tx-475569);
}
.issd-detail-item svg {
  color: var(--tx-64748b);
  flex-shrink: 0;
}
.issd-subtenant-note {
  font-size: 0.75rem;
  color: var(--tx-64748b);
  background: var(--bg-f8fafc);
  padding: 6px 10px;
  border-radius: 8px;
  border: 1px solid var(--bd-f1f5f9);
  line-height: 1.4;
}

.issd-card-actions-row {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding-top: 6px;
  border-top: 1px dashed var(--bd-f1f5f9);
}
.issd-mini-action-btn {
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  color: var(--tx-475569);
  font-family: inherit;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all .15s ease;
}
.issd-mini-action-btn:hover {
  background: var(--bg-f8fafc);
  color: var(--tx-0f172a);
  border-color: var(--bd-cbd5e1);
}
.issd-mini-action-btn.is-success {
  color: var(--tx-15803d);
  border-color: var(--bd-bbf7d0);
  background: var(--bg-f0fdf4);
}
.issd-mini-action-btn.is-success:hover {
  background: var(--bg-dcfce7);
  border-color: var(--bd-86efac);
}
.issd-mini-action-btn.is-danger {
  color: var(--tx-ef4444);
  border-color: var(--bd-fecaca);
  background: var(--bg-ffffff);
}
.issd-mini-action-btn.is-danger:hover {
  background: var(--bg-fef2f2);
  border-color: var(--bd-f87171);
}

.issd-contract-top-left {
  display: flex;
  align-items: center;
  gap: 8px;
}
.issd-contract-icon {
  color: var(--tx-2563eb);
  flex-shrink: 0;
}
.issd-contract-details-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 6px 12px;
  background: var(--bg-f8fafc);
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
  color: var(--tx-64748b);
  font-weight: 500;
  flex-shrink: 0;
}
.issd-contract-detail-item .detail-value {
  color: var(--tx-1e293b);
  font-weight: 600;
}
.issd-contract-detail-item .detail-value.is-rent {
  color: var(--tx-15803d);
  font-weight: 700;
}
.issd-mini-action-btn.is-primary {
  color: var(--tx-2563eb);
  border-color: var(--bd-bfdbfe);
  background: var(--bg-eff6ff);
}
.issd-mini-action-btn.is-primary:hover {
  background: var(--bg-dbeafe);
  border-color: var(--bd-93c5fd);
}
.issd-status-pill.is-cancelled {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
  border-color: var(--bd-fecaca);
}

.issd-pay-period strong {
  display: block;
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--tx-0f172a);
}
.issd-pay-period span {
  display: block;
  font-size: 0.75rem;
  color: var(--tx-64748b);
}
.issd-pay-amount-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin: 1px 0;
}
.issd-pay-method {
  font-size: 0.8125rem;
  color: var(--tx-475569);
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
  background: var(--bg-fef2f2);
  color: var(--tx-dc2626);
}
.issd-maintenance-type-pill.type-renovation {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}
.issd-maintenance-type-pill.type-utility {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
}
.issd-maintenance-type-pill.type-other {
  background: var(--bg-f1f5f9);
  color: var(--tx-475569);
}
.issd-maintenance-desc {
  margin: 2px 0 4px;
  font-size: 0.8125rem;
  color: var(--tx-475569);
  line-height: 1.45;
}
.issd-maintenance-meta-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 4px 10px;
  font-size: 0.75rem;
  color: var(--tx-64748b);
  background: var(--bg-f8fafc);
  padding: 8px 10px;
  border-radius: 8px;
  margin: 2px 0;
}
.issd-maintenance-meta-grid strong {
  color: var(--tx-1e293b);
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
  color: var(--tx-475569);
  font-size: 0.8125rem;
  font-weight: 600;
}
.issd-modal-form .full {
  grid-column: 1 / -1;
}
.issd-textarea {
  width: 100%;
  box-sizing: border-box;
  border: 1px solid var(--bd-cbd5e1);
  border-radius: 10px;
  padding: 9px 14px;
  font-family: inherit;
  font-size: 0.875rem;
  font-weight: 500;
  line-height: 1.5;
  color: var(--tx-1e293b);
  background: var(--bg-ffffff);
  resize: vertical;
  outline: none;
  min-height: 64px;
  transition: border-color .15s ease, box-shadow .15s ease;
}
.issd-textarea::placeholder {
  color: var(--tx-94a3b8);
  font-weight: 500;
  opacity: 1;
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
  color: var(--tx-1e293b);
  font-size: 1.5rem;
  font-weight: 800;
  line-height: 1.3;
}

.manage-page-heading .page-desc {
  margin: 4px 0 0;
  color: var(--tx-64748b);
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
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-eef2f6);
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
  color: var(--tx-1e293b);
  margin: 0;
  line-height: 1.3;
}

.gis-panel-header-left p {
  font-size: 0.8125rem !important;
  font-weight: 500;
  color: var(--tx-64748b);
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
.vertex-row span{font-size:.75rem;color:var(--tx-64748b);text-align:center}
.vertex-row button,.facility-list button,.icon-button{border:0;background:transparent;color:var(--tx-64748b);cursor:pointer;padding:5px}
.vertex-row svg,.facility-list svg,.icon-button svg{width:18px;height:18px}
.form-stack label,.form-grid label{display:grid;gap:6px;color:var(--tx-475569);font-size:.8125rem;font-weight:600}
.coordinate-grid{display:grid;grid-template-columns:1fr 1fr;gap:8px}
.facility-list{display:grid;gap:7px;margin-top:15px}
.facility-list>div{display:flex;justify-content:space-between;padding:9px 10px;background:var(--bg-f8fafc);border-radius:9px}
.facility-list span{display:grid}
.facility-list small{font-size:.75rem;color:var(--tx-64748b)}
.check-list{display:grid;gap:8px}
.check-list>div{display:flex;gap:10px;padding:12px;background:var(--bg-f8fafc);border-radius:10px}
.check-list svg{width:20px;color:var(--tx-16a34a)}
.check-list span{display:grid}
.check-list strong{font-size:.875rem}
.check-list small{font-size:.75rem;color:var(--tx-64748b)}
.stacked{display:grid;margin-top:16px}
.empty-state{padding:20px;text-align:center;color:var(--tx-64748b);background:var(--bg-f8fafc);border-radius:10px;font-size:.875rem}
.compact{padding:12px}

.form-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:13px;margin:20px 0}
.form-grid .full{grid-column:1/-1}
input,textarea{width:100%;box-sizing:border-box;border:1px solid var(--bd-cbd5e1);border-radius:8px;padding:9px 10px;color:var(--tx-1e293b);background:var(--bg-ffffff);font:inherit}
.detail-section{border-top:1px solid var(--bd-e2e8f0);padding-top:17px;margin-top:17px}
.detail-section h3{font-size:.9375rem;color:var(--tx-1e293b)}
.contract-list{display:grid;gap:8px}
.contract-list article{display:flex;justify-content:space-between;gap:10px;padding:11px;background:var(--bg-f8fafc);border-radius:9px}
.contract-list article div{display:grid}
.contract-list span{font-size:.8125rem;color:var(--tx-64748b)}
.soft-badge{background:var(--bg-eff6ff);color:var(--tx-2563eb)!important;padding:4px 8px;border-radius:7px;height:max-content}
.history-list{display:grid}
.history-list>div{display:flex;gap:10px;padding:6px 0}
.history-list i{width:9px;height:9px;border-radius:50%;background:#3b82f6;margin-top:5px}
.history-list span{display:grid}
.history-list small{font-size:.75rem;color:var(--tx-64748b)}

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
  z-index: 180;
}
.map-mode-trigger-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 48px;
  padding: 0 14px 0 8px;
  background: var(--bg-rgba-255-255-255-0p96);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--bd-rgba-226-232-240-0p95);
  border-radius: 9999px;
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  text-align: left;
  font-family: inherit;
}
.map-mode-trigger-btn:hover {
  background: var(--bg-ffffff);
  border-color: var(--bd-cbd5e1);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.12);
  transform: translateY(-1px);
}
.map-mode-trigger-btn.is-open {
  background: var(--bg-ffffff);
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
.trigger-icon-box.tone-blue { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.trigger-icon-box.tone-emerald { background: var(--bg-ecfdf5); color: var(--tx-059669); }
.trigger-icon-box.tone-amber { background: var(--bg-fffbeb); color: var(--tx-d97706); }

.trigger-text-group {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}
.trigger-eyebrow {
  font-size: 0.8125rem;
  color: var(--tx-64748b);
  font-weight: 500;
  line-height: 1;
}
.trigger-title {
  font-size: 0.875rem;
  color: var(--tx-0f172a);
  font-weight: 700;
  line-height: 1;
}
.trigger-chevron {
  color: var(--tx-64748b);
  transition: transform 0.2s ease;
  margin-left: 2px;
}
.trigger-chevron.is-open {
  transform: rotate(180deg);
  color: var(--tx-2563eb);
}

/* Dropdown Menu */
.map-mode-dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  width: 320px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
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
  background: var(--bg-f8fafc);
}
.map-mode-menu-item.is-selected {
  background: var(--bg-eff6ff);
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
.item-icon-box.tone-blue { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.item-icon-box.tone-emerald { background: var(--bg-ecfdf5); color: var(--tx-059669); }
.item-icon-box.tone-amber { background: var(--bg-fffbeb); color: var(--tx-d97706); }
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
  color: var(--tx-0f172a);
  font-weight: 700;
}
.map-mode-menu-item.is-selected .item-title {
  color: var(--tx-1d4ed8);
}
.item-desc {
  font-size: 0.8125rem;
  color: var(--tx-64748b);
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
@media (max-width: 768px) {
  /* มือถือ: แผงข้าง (แดชบอร์ด/ค้นหา/ตัวกรอง/ชั้นข้อมูล) เปิดทับเต็มจอ ซ่อนตัวเลือกมุมมองไม่ให้บังหัวแผง */
  .map-mode-dropdown-wrap.is-compact { display: none; }
}
@media (max-width: 720px) {
  .map-mode-dropdown-wrap { top: 10px; left: 10px; }
  .manage-page.is-map-page { --mp-top: 148px; --mp-top-2row: 148px; }
}

.market-map-kpi.tone-blue .market-map-kpi-icon, .market-map-kpi.tone-blue .market-map-kpi-copy em { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.market-map-kpi.tone-green .market-map-kpi-icon, .market-map-kpi.tone-green .market-map-kpi-copy em { background: var(--bg-ecfdf5); color: var(--tx-059669); }
.market-map-kpi.tone-amber .market-map-kpi-icon, .market-map-kpi.tone-amber .market-map-kpi-copy em { background: var(--bg-fffbeb); color: var(--tx-d97706); }
.market-map-kpi.tone-red .market-map-kpi-icon, .market-map-kpi.tone-red .market-map-kpi-copy em { background: var(--bg-fef2f2); color: var(--tx-dc2626); }
.market-map-kpi.tone-rose .market-map-kpi-icon, .market-map-kpi.tone-rose .market-map-kpi-copy em { background: var(--bg-fff1f2); color: var(--tx-e11d48); }
.market-map-kpi.tone-violet .market-map-kpi-icon, .market-map-kpi.tone-violet .market-map-kpi-copy em { background: var(--bg-f5f3ff); color: var(--tx-7c3aed); }

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
  color: var(--tx-334155);
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
  color: var(--tx-0f172a);
  font-weight: 700;
}
.dash-breakdown-val small {
  color: var(--tx-64748b);
  font-size: 0.8125rem;
}
.dash-progress-bar {
  width: 100%;
  height: 6px;
  background: var(--bg-f1f5f9);
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
  color: var(--tx-64748b);
  border-bottom: 1px solid var(--bd-f1f5f9);
  background: var(--bg-f8fafc);
}
.dash-mini-table td {
  padding: 10px 10px;
  border-bottom: 1px solid var(--bd-f1f5f9);
  vertical-align: middle;
}
.table-stall-pill {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 6px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  font-weight: 700;
  font-size: 0.8125rem;
}
.table-tenant-name {
  color: var(--tx-1e293b);
  font-weight: 600;
  font-size: 0.875rem;
}
.table-tenant-sub {
  color: var(--tx-64748b);
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
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
}
.table-status-pill.pill-anomaly {
  background: var(--bg-fef2f2);
  color: var(--tx-b91c1c);
}
.table-type-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}
.table-type-badge.badge-water {
  background: var(--bg-eff6ff);
  color: var(--tx-1d4ed8);
}
.table-type-badge.badge-electric {
  background: var(--bg-fef3c7);
  color: var(--tx-b45309);
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
  border: 2px solid var(--bd-ffffff);
  border-radius: 12px;
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 4px 18px rgba(37, 99, 235, 0.45), 0 2px 6px rgba(15, 23, 42, 0.2);
  cursor: pointer;
  transition: transform 0.18s ease, background 0.18s ease;
}
.map-dashboard-fab:hover { transform: translateY(-1px); background: #1d4ed8; }
.map-dashboard-fab.active { background: #1e3a8a; box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.3), 0 4px 14px rgba(15, 23, 42, 0.25); }
.map-dashboard-fab:focus-visible { outline: 3px solid var(--bd-93c5fd); outline-offset: 2px; }

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
  background: var(--bg-rgba-255-255-255-0p98);
  border: 1px solid var(--bd-e2e8f0);
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
  background: var(--bg-ffffff);
  border-bottom: 1px solid var(--bd-f1f5f9);
}
.is-map-floating .map-dash-panel-head h2 { margin: 0; color: var(--tx-0f172a); font-size: 1.125rem; font-weight: 700; }
.is-map-floating .map-dash-panel-head p { margin: 2px 0 0; color: var(--tx-64748b); font-size: .8125rem; font-weight: 500; }
.map-dash-panel-close {
  flex: 0 0 auto;
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 8px;
  background: var(--bg-ffffff);
  color: var(--tx-475569);
  cursor: pointer;
}
.map-dash-panel-close:hover { background: var(--bg-f1f5f9); }
.is-map-floating .dash-kpi-ribbon { display: grid !important; grid-template-columns: repeat(2, minmax(0, 1fr)) !important; gap: 10px !important; }
.is-map-floating .dash-two-col-grid { grid-template-columns: minmax(0, 1fr); gap: 14px; }
.is-map-floating .dash-table-wrap { overflow-x: auto; -webkit-overflow-scrolling: touch; }
@media (max-width: 720px) {
  .map-dashboard-fab { left: 10px; bottom: 12px; }
  .officer-summary-dashboard.is-map-floating { top: var(--mp-top); left: 10px; bottom: 10px; width: calc(100% - 20px); }
}

/* ปุ่มเปิดตารางจดมิเตอร์น้ำ-ไฟ: ไอคอนอย่างเดียว พื้นสีเข้ม ลอยมุมล่างซ้ายของแผนที่ ขนาดเท่าปุ่ม md (40px) */
.is-fullscreen-map-panel .map-meter-float-btn,
.map-meter-float-btn {
  position: absolute;
  left: 18px;
  bottom: 18px;
  z-index: 46;
  width: 40px;
  height: 40px;
  padding: 0;
  display: grid;
  place-items: center;
  border: 1px solid #1d4ed8;
  border-radius: 8px;
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.35);
  cursor: pointer;
  transition: transform 0.18s ease, background-color 0.18s ease, border-color 0.18s ease;
}
.map-meter-float-btn svg { width: 18px; height: 18px; flex: 0 0 auto; }
.map-meter-float-btn:hover { transform: translateY(-1px); background: #1d4ed8; border-color: #1e40af; }
/* มีการ์ดฝั่งซ้ายเปิดอยู่: ขยับปุ่มไปอยู่ทางขวาของการ์ด ไม่ซ้อนทับ (เหมือนปุ่มแดชบอร์ด) */
.has-left-panel .map-meter-float-btn { left: calc(18px + min(var(--lp-w), calc(100% - 36px)) + 12px); transition: left 0.2s ease, transform 0.18s ease, background-color 0.18s ease, border-color 0.18s ease; }
.map-meter-float-btn:focus-visible { outline: 3px solid rgba(37, 99, 235, 0.25); outline-offset: 2px; }
@media (max-width: 720px) {
  .map-meter-float-btn { left: 10px; bottom: 12px; }
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
  background: var(--bg-rgba-255-255-255-0p98);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 16px;
  box-shadow: 0 18px 44px rgba(15, 23, 42, 0.22);
}
.market-map-kpi.is-overdue .market-map-kpi-icon { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.market-map-kpi-copy small.kpi-sub { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
@media (max-width: 1000px) { .map-cost-float { top: var(--mp-top-2row); } }
@media (max-width: 720px) {
  .map-dashboard-fab.is-cost-fab.fab-utility { bottom: 12px; }
  .map-dashboard-fab.is-cost-fab.fab-bill { bottom: 66px; }
  .map-cost-float { top: var(--mp-top); left: 10px; bottom: 10px; width: calc(100% - 20px); }
}

/* ตัวกรองค่าใช้จ่าย (อยู่ในแผงสไลด์ของแผนที่) */
.cf-wrap { display: flex; flex-direction: column; gap: 16px; }
.cf-hint { margin: 0; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; line-height: 1.5; }
.cf-group { display: flex; flex-direction: column; gap: 10px; }
.cf-label { color: var(--tx-475569); font-size: 0.875rem; font-weight: 600; }
.cf-chips { display: flex; flex-wrap: wrap; gap: 8px; }
.cf-chip { display: inline-flex; align-items: center; gap: 6px; min-height: 38px; padding: 6px 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 9999px; background: var(--bg-ffffff); color: var(--tx-475569); font-family: inherit; font-size: 0.875rem; font-weight: 600; cursor: pointer; }
.cf-chip span { color: var(--tx-94a3b8); font-size: 0.8125rem; }
.cf-chip:hover { background: var(--bg-f8fafc); }
.cf-chip.active { border-color: #2563eb; background: #2563eb; color: #fff; box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3); }
.cf-chip.active span { color: #dbeafe; }
.cf-result { margin: 0; padding: 10px 12px; border-radius: 10px; background: var(--bg-eff6ff); color: var(--tx-475569); font-size: 0.875rem; font-weight: 500; }
.cf-result strong { color: var(--tx-2563eb); font-weight: 700; }

/* กันตัวอักษรในการ์ด KPI ล้นกรอบ */
.market-map-kpi { overflow: hidden; }
.market-map-kpi-copy { flex: 1 1 auto; min-width: 0; overflow: visible; }
.market-map-kpi-copy small { max-width: 100%; }
.market-map-kpi-copy strong { min-width: 0; max-width: 100%; overflow: visible; gap: 5px; font-size: clamp(1rem, 1.15vw, 1.3rem); }
.market-map-kpi-copy strong em { flex: 0 0 auto; padding: 2px 6px; white-space: nowrap; }

/* ลิ้นชักแผงจอดชิดซ้ายแทนแผงโครงการ (หน้าข้อมูลตลาดบนจอ Desktop) */
@media (min-width: 769px) {
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock) {
    top: var(--mp-top) !important;
    bottom: 18px !important;
    left: 18px !important;
    right: auto !important;
    height: auto !important;
    max-height: none !important;
    z-index: 150 !important;
    width: min(380px, calc(100% - 36px)) !important;
  }
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-edit-mode) { width: min(440px, calc(100% - 36px)) !important; }
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode) {
    width: min(720px, calc(100% - 100px)) !important;
  }
  /* เปิดตัวกรอง/แดชบอร์ดทางขวาพร้อมกัน: แผงรายละเอียดหดให้พอดีพื้นที่ที่เหลือ ไม่ทับแผงขวา */
  .manage-page.has-right-panel .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode) {
    width: min(720px, calc(100% - var(--side-w, 380px) - 84px)) !important;
    min-width: 420px;
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
    border: 1px solid var(--bd-e2e8f0) !important;
    border-radius: 18px !important;
    background: var(--bg-ffffff) !important;
    box-shadow: 0 16px 38px rgba(15, 23, 42, 0.16) !important;
    transition: width .42s cubic-bezier(.2, .8, .2, 1), box-shadow .3s ease, background-color .18s ease !important;
  }
}
.is-fullscreen-map-panel :deep(.stall-mode-panel.is-expanding-to-stall) {
  z-index: 140 !important;
  background: transparent !important;
  border-color: transparent !important;
  box-shadow: none !important;
  pointer-events: none;
}
/* แผงโครงการถูกดันออกไปทางขวาพร้อมเฟด จังหวะเดียวกับที่แผงแผงเลื่อนเข้ามา (ตอนปิดก็กลับเข้ามาย้อนทิศ) */
.is-fullscreen-map-panel :deep(.stall-mode-panel.is-expanding-to-stall > *) {
  animation: projectPanelPushOut .7s cubic-bezier(.22, 1, .36, 1) both;
  will-change: transform, opacity;
}
@keyframes projectPanelPushOut {
  from { transform: translateX(0); opacity: 1; }
  to { transform: translateX(340px); opacity: 0; }
}
.is-fullscreen-map-panel :deep(.stall-mode-panel.is-returning-from-stall) {
  animation: projectPanelFadeIn .3s ease .3s both;
  will-change: opacity;
}
@keyframes projectPanelFadeIn {
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
@media (min-width: 769px) {
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode {
    transform: translateX(0);
    clip-path: inset(0 0 0 0 round 18px);
    will-change: clip-path;
  }
  /* เปิด: เริ่มเป็นการ์ดขาวขนาดเดียวกับการ์ดโครงการ (ตรึงขอบซ้ายไว้ที่เดิม กว้าง 380px) แล้วขยายขอบขวาออกไปจนสุด ด้วย clip-path ไม่ยืด width */
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap {
    animation: issdExpandFromCard .75s cubic-bezier(.3, .7, .25, 1) both;
  }
  /* ปิด: ขอบขวาหดกลับเป็นขนาดการ์ดโครงการ */
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-drawer-closing {
    clip-path: inset(0 calc(100% - 380px) 0 0 round 18px);
    background-color: var(--bg-ffffff);
    opacity: 0;
    transition: clip-path .5s cubic-bezier(.4, 0, .2, 1), background-color .5s ease, opacity .12s ease .48s !important;
    pointer-events: none;
  }
}
@keyframes issdExpandFromCard {
  from { clip-path: inset(0 calc(100% - 380px) 0 0 round 18px); background-color: var(--bg-ffffff); }
  to { clip-path: inset(0 0 0 0 round 18px); background-color: var(--bg-f8fafc); }
}
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap .issd-header,
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap .issd-full-sidebar,
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-project-morph-swap .issd-full-main {
  animation: stallPanelContentFadeIn .5s ease-out both;
}
@keyframes stallPanelContentFadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
/* ปิด: ยกแต่ละส่วนของเนื้อหาขึ้นเป็นเลเยอร์ของตัวเอง ให้ clip-path หุบโดยไม่ต้องวาดเนื้อหาทั้งหมดใหม่ทุกเฟรม (กันกระตุก) */
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-drawer-closing > *,
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-drawer-closing .issd-header,
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-drawer-closing .issd-full-sidebar {
  will-change: opacity;
  transform: translateZ(0);
}
/* ปิด: เนื้อหาคงอยู่จนขอบขวาหุบถึงซ้ายแล้วค่อยเฟดออก เห็นเป็นการหุบจากขวาไปซ้ายชัดเจน */
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-drawer-closing > *,
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-drawer-closing .issd-header,
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-drawer-closing .issd-full-sidebar {
  opacity: 0;
  transition: opacity .15s ease .4s !important;
}
/* ปิด: ฝั่งเนื้อหาด้านขวาเลื่อนเข้าซ้ายไปพร้อมขอบที่หุบ (720 -> 380px) และจางลงระหว่างทาง ไอคอนกับข้อความจึงเคลื่อนตามขอบ ไม่ใช่แค่ถูกตัดทิ้ง */
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode.is-drawer-closing .issd-full-main {
  transform: translateX(-340px);
  opacity: 0;
  transition: transform .5s cubic-bezier(.4, 0, .2, 1), opacity .32s ease-in !important;
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
.issd-return-row { width: 100%; padding-bottom: 14px; margin-bottom: 14px; border-bottom: 1px solid var(--bd-e8eef6); }
.issd-return-row .issd-return-btn { margin: 0; }

/* รายละเอียดแผงแบบเต็มบนจอ Desktop: หัวแผงกว้างเท่าแถบเมนูข้าง แล้วการ์ด "สิ่งที่ต้องดำเนินการ" อยู่บนสุดของคอลัมน์ขวา */
@media (min-width: 769px) {
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode {
    display: grid !important;
    grid-template-columns: 188px minmax(0, 1fr);
    grid-template-rows: auto minmax(0, 1fr) auto;
    background: var(--bg-f8fafc);
  }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-full-layout { display: contents; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-full-main { grid-column: 2; grid-row: 1 / 4; display: flex; flex-direction: column; min-width: 0; min-height: 0; overflow-y: auto; overflow-x: hidden; -webkit-overflow-scrolling: touch; }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-header {
    grid-column: 1;
    grid-row: 1;
    padding: 18px 20px 14px; /* ขอบในซ้าย-ขวาเท่าแผงโครงการ (20px) */
    border-right: 1px solid var(--bd-eef2f6);
    background: var(--bg-ffffff);
  }
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-full-sidebar { grid-column: 1; grid-row: 2 / 4; width: auto; }
  /* การ์ดต้องดำเนินการอยู่บนสุดของคอลัมน์ขวา ชิดกับเนื้อหาด้านล่าง ไม่เว้นช่องว่าง */
  .inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode .issd-action-alert-card { margin: 20px 18px 0 !important; flex: 0 0 auto; }
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
  border: 1px solid var(--bd-bfdbfe);
  border-radius: 999px;
  overflow: visible;
  background: var(--bg-rgba-255-255-255-0p97);
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
  color: var(--tx-64748b);
  font-size: 0.9375rem;
  font-weight: 700;
  line-height: 1;
  white-space: nowrap;
  transition: background 0.15s ease, color 0.15s ease;
}
.is-fullscreen-map-panel :deep(.market-floor-tab:hover:not(.active)) { background: var(--bg-f1f5f9); color: var(--tx-0f172a); }
.is-fullscreen-map-panel :deep(.market-floor-tab.active) { background: #2563eb; color: #ffffff; box-shadow: 0 3px 10px rgba(37, 99, 235, 0.3); }


/* แผงโครงการบน Desktop: ความสูงยืดจากใต้การ์ด KPI ถึงเหนือปุ่มแดชบอร์ด */
@media (min-width: 769px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel.stall-mode-panel) {
    top: var(--mp-top) !important;
    bottom: 18px !important;
    height: auto !important;
    max-height: none !important;
  }
}
@media (min-width: 769px) and (max-width: 1000px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel.stall-mode-panel) { top: var(--mp-top-2row) !important; }
}

/* แผงโครงการที่ใช้เนื้อหาเฉพาะโหมด (จดมิเตอร์): สูงเท่าเนื้อหา ไม่ยืดเต็มความสูง */
.is-fullscreen-map-panel { --meter-card-h: min(660px, calc(100% - var(--mp-top, 114px) - 18px)); }
@media (min-width: 769px) {
  /* การ์ดโครงการและการ์ดจดรายแผงสูงเท่ากัน ชิดล่าง ไม่กระโดดตอนสลับ */
  .is-fullscreen-map-panel :deep(.stall-mode-panel.has-custom-body.stall-mode-panel) {
    top: auto !important;
    bottom: 18px !important;
    height: var(--meter-card-h) !important;
    max-height: none !important;
  }
}
@media (max-width: 768px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel.has-custom-body.stall-mode-panel) {
    height: 62dvh !important;
    max-height: calc(100dvh - 58px) !important;
  }
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
.issd-return-btn:focus-visible { outline: 3px solid var(--bd-93c5fd); outline-offset: 2px; }
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
  background: var(--bg-f8fafc);
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

/* ซ่อนปุ่ม แผนที่/ดาวเทียม, ปรับกึ่งกลาง และ คำอธิบายสีโครงการ (เพราะมีในแถบเครื่องมือชั้นข้อมูลด้านขวาแล้ว) — ไม่ซ่อนกล่อง map-floating-top-controls เพราะมีแถบเลือกชั้นอยู่ข้างใน */
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
  background: var(--bg-rgba-255-255-255-0p96) !important;
  backdrop-filter: blur(12px) !important;
  -webkit-backdrop-filter: blur(12px) !important;
  border: 1px solid var(--bd-rgba-226-232-240-0p95) !important;
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
  background: linear-gradient(135deg, var(--bg-f8fafc) 0%, var(--bg-eff6ff) 100%);
  border: 1px solid var(--bd-dbeafe);
  border-radius: 14px;
  padding: 16px 18px;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.04);
}
.ov-dash-hero-box.finance {
  background: linear-gradient(135deg, var(--bg-f0fdf4) 0%, var(--bg-eff6ff) 100%);
  border-color: var(--bd-bbf7d0);
}
.ov-dash-hero-box.meter {
  background: linear-gradient(135deg, var(--bg-f0fdfa) 0%, var(--bg-eff6ff) 100%);
  border-color: var(--bd-99f6e4);
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
  color: var(--tx-1e40af);
  font-size: 0.8125rem;
  font-weight: 600;
  margin-bottom: 3px;
}
.ov-dash-hero-value {
  font-size: 1.75rem;
  font-weight: 800;
  color: var(--tx-0f172a);
  letter-spacing: -0.02em;
  line-height: 1.15;
}
.ov-dash-hero-value small {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--tx-64748b);
  letter-spacing: normal;
}
.ov-dash-hero-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  background: var(--bg-dbeafe);
  color: var(--tx-1d4ed8);
  font-size: 0.75rem;
  font-weight: 700;
  border-radius: 9999px;
  white-space: nowrap;
}
.ov-dash-hero-pill.green {
  background: var(--bg-dcfce7);
  color: var(--tx-15803d);
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
.health-label-good { color: var(--tx-16a34a); }
.health-label-warn { color: var(--tx-d97706); }
.health-label-danger { color: var(--tx-dc2626); }
.ov-dash-health-track {
  display: flex;
  width: 100%;
  height: 7px;
  background: var(--bg-e2e8f0);
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
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  transition: border-color 0.2s, box-shadow 0.2s;
}
.ov-dash-kpi-card:hover {
  border-color: var(--bd-cbd5e1);
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
.ov-dash-kpi-icon.blue { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.ov-dash-kpi-icon.amber { background: var(--bg-fef3c7); color: var(--tx-d97706); }
.ov-dash-kpi-icon.red { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.ov-dash-kpi-icon.green { background: var(--bg-dcfce7); color: var(--tx-16a34a); }

.ov-dash-kpi-title {
  font-size: 0.8125rem;
  color: var(--tx-64748b);
  font-weight: 600;
  white-space: nowrap;
}
.ov-dash-kpi-number {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--tx-0f172a);
  line-height: 1.15;
}
.ov-dash-kpi-number small {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--tx-64748b);
}
.ov-dash-kpi-number.text-amber { color: var(--tx-d97706); }
.ov-dash-kpi-number.text-red { color: var(--tx-dc2626); }
.ov-dash-kpi-number.text-green { color: var(--tx-16a34a); }

.ov-dash-legal-card {
  width: 100%; min-height: 64px; display: flex; align-items: center; gap: 10px;
  padding: 10px 12px; border: 1px solid var(--bd-fecaca); border-radius: 12px;
  background: var(--bg-fff7f7); color: var(--tx-991b1b); font-family: inherit; text-align: left; cursor: pointer;
  transition: border-color .15s ease, background .15s ease, box-shadow .15s ease;
}
.ov-dash-legal-card:hover,.ov-dash-legal-card.active { border-color: #ef4444; background: var(--bg-fef2f2); box-shadow: 0 2px 8px rgba(220,38,38,.1); }
.ov-dash-legal-icon { width: 34px; height: 34px; flex: 0 0 34px; display: grid; place-items: center; border-radius: 9px; background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.ov-dash-legal-icon svg { width: 18px; height: 18px; }
.ov-dash-legal-copy { flex: 1 1 auto; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.ov-dash-legal-copy small { font-size: .8125rem; color: var(--tx-7f1d1d); font-weight: 600; }
.ov-dash-legal-copy strong { font-size: 1.125rem; color: var(--tx-b91c1c); font-weight: 800; line-height: 1.2; }
.ov-dash-legal-copy em { font-size: .8125rem; font-style: normal; font-weight: 600; }
.ov-dash-legal-action { flex: 0 0 auto; color: var(--tx-dc2626); font-size: .8125rem; font-weight: 600; }

/* การ์ดสรุปร่วมสำหรับแดชบอร์ดข้อมูลแผง การชำระเงิน และการจดมิเตอร์ */
.ov-dash-summary-card,
.ov-stall-overview-card,
.ov-stall-priority-card {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
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
  color: var(--tx-0f172a);
  font-size: 0.9375rem;
  font-weight: 700;
  line-height: 1.35;
}

.ov-stall-card-head p {
  margin: 2px 0 0;
  color: var(--tx-64748b);
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-stall-section-note {
  margin: 2px 0 0;
  color: var(--tx-64748b);
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-stall-occupancy {
  display: inline-flex;
  align-items: center;
  min-height: 26px;
  padding: 3px 9px;
  border-radius: 999px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
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
  background: var(--bg-e2e8f0);
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
  border: 1px solid var(--bd-eef2f6);
  border-radius: 10px;
  background: var(--bg-f8fafc);
  color: var(--tx-334155);
  font-family: inherit;
  text-align: left;
  cursor: pointer;
  transition: border-color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease;
}

.ov-stall-status-actions button:hover,
.ov-stall-status-actions button.active {
  border-color: var(--bd-93c5fd);
  background: var(--bg-eff6ff);
  box-shadow: 0 1px 4px rgba(37, 99, 235, 0.08);
}

.ov-stall-status-actions button span {
  display: flex;
  align-items: center;
  gap: 5px;
  max-width: 100%;
  color: var(--tx-64748b);
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
  color: var(--tx-0f172a);
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
  border: 1px solid var(--bd-eef2f6);
  border-radius: 10px;
  background: var(--bg-f8fafc);
  color: var(--tx-334155);
  font-family: inherit;
  text-align: left;
}

button.ov-stall-priority-row {
  cursor: pointer;
  transition: border-color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease;
}

button.ov-stall-priority-row:hover,
button.ov-stall-priority-row.active {
  border-color: var(--bd-93c5fd);
  background: var(--bg-eff6ff);
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

.ov-stall-priority-row.is-action .ov-stall-priority-icon { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.ov-stall-priority-row.is-legal .ov-stall-priority-icon { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.ov-stall-priority-row.is-expiring .ov-stall-priority-icon { background: var(--bg-fef3c7); color: var(--tx-d97706); }

.ov-stall-priority-label {
  min-width: 0;
  overflow: hidden;
  color: var(--tx-334155);
  font-size: 0.8125rem;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ov-stall-priority-row strong {
  min-width: 24px;
  color: var(--tx-0f172a);
  font-size: 1rem;
  font-weight: 800;
  text-align: right;
}

.ov-stall-priority-row.is-legal strong { color: var(--tx-dc2626); }
.ov-stall-priority-row.is-expiring strong { color: var(--tx-b45309); }

.ov-stall-priority-chevron {
  width: 15px;
  height: 15px;
  color: var(--tx-94a3b8);
}

.ov-dash-summary-main {
  min-height: 50px;
  padding: 9px 11px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border-radius: 10px;
  background: var(--bg-eff6ff);
}

.ov-dash-summary-main span,
.ov-dash-summary-grid span {
  color: var(--tx-64748b);
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-dash-summary-main strong {
  color: var(--tx-1d4ed8);
  font-size: 1.125rem;
  font-weight: 800;
  white-space: nowrap;
}

.ov-dash-summary-main strong small {
  margin-left: 4px;
  color: var(--tx-64748b);
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
  background: var(--bg-eff6ff);
}

.ov-dash-utilization-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 10px;
}

.ov-dash-utilization-head span {
  color: var(--tx-475569);
  font-size: 0.8125rem;
  font-weight: 600;
}

.ov-dash-utilization-head strong {
  color: var(--tx-1d4ed8);
  font-size: 1.125rem;
  font-weight: 800;
}

.ov-dash-utilization-track {
  width: 100%;
  height: 6px;
  overflow: hidden;
  border-radius: 999px;
  background: var(--bg-dbeafe);
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
  border: 1px solid var(--bd-eef2f6);
  border-radius: 10px;
  background: var(--bg-f8fafc);
}

.ov-dash-summary-grid strong {
  color: var(--tx-0f172a);
  font-size: 1rem;
  font-weight: 800;
  white-space: nowrap;
}

.ov-dash-summary-grid .is-success strong { color: var(--tx-15803d); }
.ov-dash-summary-grid .is-pending strong { color: var(--tx-b45309); }
.ov-dash-summary-grid .is-danger strong,
.ov-dash-summary-grid .is-expired strong { color: var(--tx-dc2626); }

.ov-dash-summary-grid small {
  color: var(--tx-64748b);
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
  background: var(--bg-eff6ff);
}

.ov-dash-meter-progress-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 10px;
}

.ov-dash-meter-progress-head span {
  color: var(--tx-64748b);
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-dash-meter-progress-head strong {
  color: var(--tx-1d4ed8);
  font-size: 1rem;
  font-weight: 800;
  white-space: nowrap;
}

.ov-dash-meter-progress-head small {
  color: var(--tx-64748b);
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
  background: var(--bg-eff6ff);
}

.ov-dash-income-total span {
  color: var(--tx-475569);
  font-size: 0.8125rem;
  font-weight: 600;
}

.ov-dash-income-total strong {
  color: var(--tx-1d4ed8);
  font-size: 1rem;
  font-weight: 800;
  white-space: nowrap;
}

.ov-project-status-legend {
  margin: -2px 0 9px;
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--tx-64748b);
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
.ov-project-status-list::-webkit-scrollbar-thumb { border-radius: 999px; background: var(--bg-cbd5e1); }
.ov-project-status-list::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

.ov-project-donut-focus {
  margin-top: 14px;
  padding: 16px 14px;
  display: grid;
  grid-template-columns: 132px minmax(0, 1fr);
  align-items: center;
  gap: 16px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  background: linear-gradient(145deg, var(--bg-ffffff) 0%, var(--bg-f8fafc) 100%);
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
  background: var(--bg-ffffff);
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
  color: var(--tx-0f172a);
  font-size: 1.5rem;
  font-weight: 800;
  line-height: 1.15;
}

.ov-project-donut-chart small {
  color: var(--tx-64748b);
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-project-donut-name {
  max-width: 100%;
  overflow: hidden;
  color: var(--tx-334155);
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
  border: 1px solid var(--bd-eef2f6);
  border-radius: 10px;
  background: var(--bg-rgba-255-255-255-0p88);
}

.ov-project-donut-stats span {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--tx-64748b);
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
  color: var(--tx-1e293b);
  font-size: 1rem;
  font-weight: 800;
}

.ov-project-donut-stats small {
  grid-column: 1 / -1;
  padding-left: 14px;
  color: var(--tx-94a3b8);
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-project-donut-stats .is-available i { background: #10b981; }
.ov-project-donut-stats .is-occupied i { background: #64748b; }
.ov-project-donut-stats .is-reserved i { background: #f59e0b; }
.ov-project-donut-stats .is-available strong { color: var(--tx-047857); }
.ov-project-donut-stats .is-reserved strong { color: var(--tx-b45309); }

.ov-project-donut-empty {
  margin-top: 14px;
  padding: 28px 16px;
  border: 1px dashed var(--bd-cbd5e1);
  border-radius: 12px;
  background: var(--bg-f8fafc);
  color: var(--tx-64748b);
  font-size: 0.875rem;
  font-weight: 500;
  text-align: center;
}

.ov-project-column-legend {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 8px 14px;
  margin: 14px 0 10px;
  text-align: center;
}

.ov-project-column-legend span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--tx-475569);
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-project-column-legend i {
  width: 11px;
  height: 11px;
  flex: 0 0 11px;
  border-radius: 3px;
}

.ov-project-column-legend strong {
  margin-left: 2px;
  color: var(--tx-0f172a);
  font-size: 0.875rem;
  font-weight: 800;
}

.ov-project-column-legend .is-available i { background: #10b981; }
.ov-project-column-legend .is-reserved i { background: #f59e0b; }
.ov-project-column-legend .is-occupied i { background: #475569; }

.ov-project-column-chart {
  position: relative;
  display: grid;
  grid-template-columns: 38px minmax(0, 1fr);
  min-height: 278px;
}

.ov-project-column-y-axis {
  position: relative;
  height: 238px;
  margin-top: 8px;
  color: var(--tx-64748b);
  font-size: 0.75rem;
  font-weight: 500;
  pointer-events: none;
}

.ov-project-column-y-axis span {
  position: absolute;
  right: 8px;
  line-height: 1;
  transform: translateY(50%);
}

.ov-project-column-scroll {
  min-width: 0;
  padding-bottom: 9px;
  overflow-x: auto;
  overflow-y: hidden;
  overscroll-behavior-x: contain;
  scrollbar-color: #cbd5e1 #f1f5f9;
  scrollbar-width: thin;
}

.ov-project-column-scroll:focus-visible {
  border-radius: 8px;
  outline: 2px solid var(--bd-93c5fd);
  outline-offset: 2px;
}

.ov-project-column-scroll::-webkit-scrollbar { height: 8px; }
.ov-project-column-scroll::-webkit-scrollbar-track { border-radius: 999px; background: var(--bg-f1f5f9); }
.ov-project-column-scroll::-webkit-scrollbar-thumb { border-radius: 999px; background: var(--bg-cbd5e1); }
.ov-project-column-scroll::-webkit-scrollbar-thumb:hover { background: #94a3b8; }

.ov-project-column-grid {
  position: relative;
  min-width: max(100%, calc(var(--project-count) * 96px));
  height: 268px;
  padding: 8px 4px 0;
  display: grid;
  grid-template-columns: repeat(var(--project-count), minmax(84px, 1fr));
  align-items: end;
  gap: 10px;
  background-image: repeating-linear-gradient(
    to top,
    var(--bg-e2e8f0) 0,
    var(--bg-e2e8f0) 1px,
    transparent 1px,
    transparent var(--grid-step)
  );
  background-position: 0 8px;
  background-size: 100% 238px;
  background-repeat: no-repeat;
}

.ov-project-column-item {
  position: relative;
  height: 260px;
  min-width: 0;
  padding: 0 3px 31px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  border: 0;
  border-radius: 9px 9px 0 0;
  background: transparent;
  color: inherit;
  font-family: inherit;
  cursor: pointer;
}

.ov-project-column-item::after {
  position: absolute;
  right: 0;
  bottom: 27px;
  left: 0;
  height: 3px;
  border-radius: 3px;
  background: transparent;
  content: '';
}

.ov-project-column-item:hover,
.ov-project-column-item:focus-visible {
  background: var(--bg-rgba-239-246-255-0p72);
  outline: none;
}

.ov-project-column-item.active {
  box-shadow: inset 0 0 0 2px #2563eb;
  background: var(--bg-rgba-239-246-255-0p76);
}

.ov-project-column-item.active::after { background: #2563eb; }

.ov-project-column-total {
  position: absolute;
  bottom: calc(34px + (var(--bar-ratio, 0) * 205px));
  left: 50%;
  z-index: 2;
  color: var(--tx-334155);
  font-size: 0.75rem;
  font-weight: 800;
  line-height: 1;
  transform: translateX(-50%);
  white-space: nowrap;
}

.ov-project-column-stack {
  width: min(30px, 68%);
  height: calc(var(--bar-ratio, 0) * 205px);
  min-height: 2px;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  overflow: hidden;
  border-radius: 6px 6px 2px 2px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.08);
  transition: filter 0.16s ease, transform 0.16s ease;
}

.ov-project-column-item:hover .ov-project-column-stack,
.ov-project-column-item:focus-visible .ov-project-column-stack { filter: saturate(1.08); transform: translateY(-2px); }

.ov-project-column-stack i {
  width: 100%;
  min-height: 0;
  display: block;
}

.ov-project-column-stack .is-available { background: #10b981; }
.ov-project-column-stack .is-reserved { background: #f59e0b; }
.ov-project-column-stack .is-occupied { background: #475569; }

.ov-project-column-label {
  position: absolute;
  right: 1px;
  bottom: 5px;
  left: 1px;
  overflow: hidden;
  color: var(--tx-475569);
  font-size: 0.75rem;
  font-weight: 600;
  line-height: 1.2;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ov-project-column-item.active .ov-project-column-label { color: var(--tx-1d4ed8); font-weight: 800; }

.ov-project-column-tooltip {
  position: absolute;
  right: 2px;
  bottom: 40px;
  left: 2px;
  z-index: 5;
  width: 158px;
  padding: 8px 10px;
  display: none;
  flex-direction: column;
  gap: 2px;
  border-radius: 9px;
  background: rgba(15, 23, 42, 0.96);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.24);
  color: #ffffff;
  text-align: left;
  transform: translateX(-30%);
}

.ov-project-column-tooltip strong { font-size: 0.8125rem; font-weight: 700; }
.ov-project-column-tooltip-stats {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.ov-project-column-tooltip small {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  color: var(--tx-cbd5e1);
  font-size: 0.75rem;
  font-weight: 500;
  white-space: nowrap;
}
.ov-project-column-tooltip small::before {
  width: 6px;
  height: 6px;
  flex: 0 0 6px;
  border-radius: 2px;
  content: '';
}
.ov-project-column-tooltip small.is-available::before { background: #10b981; }
.ov-project-column-tooltip small.is-reserved::before { background: #f59e0b; }
.ov-project-column-tooltip small.is-occupied::before { background: #64748b; }
.ov-project-column-tooltip small b {
  margin-left: auto;
  color: #ffffff;
  font-weight: 700;
}
.ov-project-column-item:first-child .ov-project-column-tooltip {
  right: auto;
  left: 2px;
  transform: none;
}
.ov-project-column-item:last-child .ov-project-column-tooltip {
  right: 2px;
  left: auto;
  transform: none;
}
.ov-project-column-item:hover .ov-project-column-tooltip,
.ov-project-column-item:focus-visible .ov-project-column-tooltip { display: flex; }

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
.ov-project-mini-donut-grid::-webkit-scrollbar-thumb { border-radius: 999px; background: var(--bg-cbd5e1); }

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
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 13px;
  background: linear-gradient(155deg, var(--bg-ffffff) 0%, var(--bg-f8fafc) 100%);
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
  outline: none;
  transition: border-color 0.18s ease, box-shadow 0.18s ease, transform 0.18s ease;
}

.ov-project-mini-donut-card:hover,
.ov-project-mini-donut-card:focus-visible,
.ov-project-mini-donut-card.active {
  z-index: 2;
  border-color: var(--bd-93c5fd);
  box-shadow: 0 7px 18px rgba(37, 99, 235, 0.13);
  transform: translateY(-2px);
}

.ov-project-mini-donut-card.active {
  background: linear-gradient(155deg, var(--bg-ffffff) 0%, var(--bg-eff6ff) 100%);
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

.ov-project-mini-badge.is-danger { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.ov-project-mini-badge.is-available { background: var(--bg-dcfce7); color: var(--tx-15803d); }

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
  background: var(--bg-ffffff);
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

.ov-project-mini-donut strong { color: var(--tx-0f172a); font-size: 1.125rem; font-weight: 800; }
.ov-project-mini-donut small { color: var(--tx-64748b); font-size: 0.75rem; font-weight: 600; }

.ov-project-mini-name {
  width: 100%;
  overflow: hidden;
  color: var(--tx-334155);
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
.ov-project-mini-tooltip small { color: var(--tx-cbd5e1); font-size: 0.75rem; font-weight: 500; }
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
  color: var(--tx-475569);
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
  color: var(--tx-1e293b);
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
  background: var(--bg-e2e8f0);
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
  background: var(--bg-ffffff);
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
  color: var(--tx-0f172a);
  font-size: 1.125rem;
  font-weight: 800;
}

.ov-contract-donut small {
  color: var(--tx-64748b);
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-contract-health-legend {
  display: flex;
  flex-direction: column;
  gap: 7px;
  color: var(--tx-475569);
  font-size: 0.8125rem;
  font-weight: 500;
}

.ov-contract-health-legend span strong {
  margin-left: auto;
  color: var(--tx-1e293b);
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
  border-top: 1px solid var(--bd-f1f5f9);
}

.ov-contract-renewal-row span {
  min-width: 0;
  color: var(--tx-64748b);
  font-size: 0.75rem;
  font-weight: 500;
  text-align: center;
}

.ov-contract-renewal-row strong {
  display: block;
  margin-top: 2px;
  color: var(--tx-1e293b);
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
  border: 1px solid var(--bd-eef2f6);
  border-radius: 9px;
  background: var(--bg-f8fafc);
}

.ov-status-summary-grid span {
  color: var(--tx-64748b);
  font-size: 0.75rem;
  font-weight: 500;
  line-height: 1.25;
}

.ov-status-summary-grid strong {
  color: var(--tx-1e293b);
  font-size: 1rem;
  font-weight: 800;
}

.ov-status-summary-grid .is-success strong { color: var(--tx-15803d); }
.ov-status-summary-grid .is-pending strong { color: var(--tx-b45309); }
.ov-status-summary-grid .is-danger strong { color: var(--tx-dc2626); }

.ov-maintenance-legend i.is-pending { background: #f59e0b; }
.ov-maintenance-legend i.is-repairing { background: #2563eb; }
.ov-maintenance-legend i.is-completed { background: #16a34a; }
.ov-maintenance-legend i.is-overdue { background: #ef4444; }

.ov-maintenance-donut-layout {
  display: grid;
  grid-template-columns: 104px minmax(0, 1fr);
  align-items: center;
  gap: 18px;
}

.ov-maintenance-donut {
  position: relative;
  width: 100px;
  height: 100px;
  display: grid;
  place-items: center;
  border-radius: 50%;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.08);
}

.ov-maintenance-donut::after {
  position: absolute;
  inset: 14px;
  border-radius: 50%;
  background: var(--bg-ffffff);
  box-shadow: inset 0 0 0 1px #f1f5f9;
  content: '';
}

.ov-maintenance-donut > span {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1px;
  line-height: 1.1;
}

.ov-maintenance-donut strong {
  color: var(--tx-0f172a);
  font-size: 1.25rem;
  font-weight: 800;
}

.ov-maintenance-donut small {
  color: var(--tx-64748b);
  font-size: 0.75rem;
  font-weight: 600;
}

.ov-maintenance-legend {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 7px;
  color: var(--tx-64748b);
  font-size: 0.75rem;
  font-weight: 500;
}

.ov-maintenance-legend strong {
  margin-left: auto;
  color: var(--tx-1e293b);
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

.ov-dash-container .ov-operational-status-card .ov-status-summary-grid > div {
  min-height: 58px;
  padding: 9px 11px;
}

.ov-dash-container .ov-operational-status-card .ov-status-summary-grid span,
.ov-dash-container .ov-operational-status-card .ov-maintenance-legend {
  font-size: 0.875rem;
  line-height: 1.45;
}

.ov-dash-container .ov-operational-status-card .ov-status-summary-grid strong {
  font-size: 1.125rem;
}

.ov-dash-container .ov-operational-status-card .ov-maintenance-legend strong {
  font-size: 0.9375rem;
  font-weight: 800;
}

/* 3. Section Cards */
.ov-dash-section-card {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
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
  color: var(--tx-0f172a);
}
.ov-dash-count-badge {
  font-size: 0.75rem;
  font-weight: 700;
  padding: 1px 7px;
  border-radius: 9999px;
}
.ov-dash-count-badge.badge-amber {
  background: var(--bg-fef3c7);
  color: var(--tx-b45309);
}
.ov-dash-count-badge.badge-red {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
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
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-eef2f6);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s ease;
}
.ov-dash-watchlist-row:hover {
  background: var(--bg-ffffff);
  border-color: var(--bd-93c5fd);
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.06);
  transform: translateY(-1px);
}
.ov-dash-watchlist-row.urgency-expired {
  border-left: 3px solid #ef4444;
  background: var(--bg-fef2f2);
}
.ov-dash-watchlist-row.urgency-expired:hover {
  background: var(--bg-ffffff);
  border-color: var(--bd-fca5a5);
}
.ov-dash-watchlist-row.urgency-urgent {
  border-left: 3px solid #f97316;
  background: var(--bg-fff7ed);
}
.ov-dash-watchlist-row.urgency-urgent:hover {
  background: var(--bg-ffffff);
  border-color: var(--bd-fdba74);
}
.ov-dash-watchlist-row.urgency-warning {
  border-left: 3px solid #f59e0b;
  background: var(--bg-fffbeb);
}
.ov-dash-watchlist-row.urgency-warning:hover {
  background: var(--bg-ffffff);
  border-color: var(--bd-fde68a);
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
  color: var(--tx-1e293b);
  background: var(--bg-e2e8f0);
  padding: 2px 6px;
  border-radius: 5px;
  letter-spacing: 0.02em;
  flex-shrink: 0;
}
.ov-dash-w-tenant {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--tx-1e293b);
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
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
  font-weight: 700;
}
.ov-dash-w-pill.pill-urgent {
  background: var(--bg-ffedd5);
  color: var(--tx-c2410c);
  font-weight: 700;
}
.ov-dash-w-pill.pill-warning {
  background: var(--bg-fef3c7);
  color: var(--tx-b45309);
}
.ov-dash-w-pill.pill-renew-submitted {
  background: var(--bg-dcfce7);
  color: var(--tx-15803d);
}
.ov-dash-w-pill.pill-renew-pending_submit {
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
}
.ov-dash-w-pill.pill-renew-declined {
  background: var(--bg-fee2e2);
  color: var(--tx-991b1b);
}

.ov-dash-w-icon {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--tx-94a3b8);
  background: var(--bg-f1f5f9);
  transition: all 0.15s ease;
}
.ov-dash-watchlist-row:hover .ov-dash-w-icon {
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
}

/* 5. Minimalist 4-Stat Strip */
.ov-dash-pipeline-strip {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8px;
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-eef2f6);
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
  color: var(--tx-64748b);
  font-weight: 500;
  white-space: nowrap;
}
.ov-dash-pipe-stat strong {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--tx-0f172a);
  line-height: 1.1;
}
.ov-dash-pipe-stat strong.text-blue { color: var(--tx-2563eb); }
.ov-dash-pipe-stat strong.text-amber { color: var(--tx-d97706); }
.ov-dash-pipe-stat strong.text-purple { color: var(--tx-7c3aed); }
.ov-dash-pipe-stat strong.text-green { color: var(--tx-16a34a); }

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
  border-bottom: 1px solid var(--bd-f1f5f9);
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
  color: var(--tx-334155);
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
  color: var(--tx-64748b);
}
.ov-dash-cat-bar {
  width: 100%;
  height: 5px;
  background: var(--bg-f1f5f9);
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

@media (max-width: 1024px) {
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
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  background: var(--bg-ffffff);
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
  border-color: var(--bd-cbd5e1);
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
  color: var(--tx-64748b);
  font-weight: 500;
  line-height: 1.25;
  white-space: normal;
  word-break: break-word;
}

.ov-dash-status-pill-val {
  font-size: 1.0625rem;
  font-weight: 800;
  color: var(--tx-0f172a);
  line-height: 1.2;
  margin-top: 2px;
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.ov-dash-status-pill-val small {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--tx-64748b);
}

/* Color schemes for each status */
.ov-dash-status-pill.is-total .ov-dash-status-pill-icon { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.ov-dash-status-pill.is-total.active {
  border-color: #2563eb;
  background: var(--bg-f8faff);
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.18);
}

.ov-dash-status-pill.is-available .ov-dash-status-pill-icon { background: var(--bg-ecfdf5); color: var(--tx-059669); }
.ov-dash-status-pill.is-available.active {
  border-color: #10b981;
  background: var(--bg-f0fdf4);
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.ov-dash-status-pill.is-reserved .ov-dash-status-pill-icon { background: var(--bg-fffbeb); color: var(--tx-d97706); }
.ov-dash-status-pill.is-reserved.active {
  border-color: #f59e0b;
  background: var(--bg-fffdf7);
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.2);
}

.ov-dash-status-pill.is-occupied .ov-dash-status-pill-icon { background: var(--bg-f1f5f9); color: var(--tx-64748b); }
.ov-dash-status-pill.is-occupied.active {
  border-color: var(--bd-94a3b8);
  background: var(--bg-f8fafc);
  box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.18);
}

.ov-dash-status-pill.is-legal-case .ov-dash-status-pill-icon { background: var(--bg-fef2f2); color: var(--tx-ef4444); }
.ov-dash-status-pill.is-legal-case.active {
  border-color: #ef4444;
  background: var(--bg-fef2f2);
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
}

.ov-dash-status-pill.is-action-required .ov-dash-status-pill-icon { background: var(--bg-eef2ff); color: var(--tx-6366f1); }
.ov-dash-status-pill.is-action-required.active {
  border-color: #6366f1;
  background: var(--bg-f5f3ff);
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
}

/* Market information overview: KPI cards floating above the GIS map */
.market-map-kpi-strip {
  position: absolute;
  top: 18px;
  left: calc(50% + 34px);
  right: auto;
  transform: translateX(-50%);
  width: min(1560px, calc(100% - 320px));
  max-width: 1560px;
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

/* โหมดจดมิเตอร์มี 5 ใบ: ใช้ขนาดการ์ดเท่าแถบ 6 ใบ แต่ลดความกว้างแถบเหลือ 5/6 เพื่อให้การ์ดอยู่กึ่งกลางแผนที่ */
.market-map-kpi-strip.is-meters { grid-template-columns: repeat(5, minmax(0, 1fr)); }
.market-map-kpi-strip.is-meters:not(.has-side-panel) { width: min(1300px, calc((100% - 320px) * 5 / 6)); max-width: 1300px; }
.market-map-kpi-strip.is-meters.has-side-panel { width: min(875px, calc((100% - var(--side-w, 380px) - 170px) * 5 / 6)); }

.market-map-kpi {
  min-width: 0;
  min-height: 80px;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  border: 1px solid var(--bd-rgba-226-232-240-0p94);
  border-radius: 15px;
  background: var(--bg-rgba-255-255-255-0p96);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.11);
  color: var(--tx-475569);
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
  border-color: var(--bd-bfdbfe);
  box-shadow: 0 10px 26px rgba(15, 23, 42, 0.14);
}

.market-map-kpi.active {
  border: 2px solid #3b82f6;
  padding: 13px 15px;
  box-shadow: 0 8px 24px rgba(37, 99, 235, 0.2);
}

.market-map-kpi-icon {
  width: 40px;
  height: 40px;
  flex: 0 0 40px;
  display: grid;
  place-items: center;
  border-radius: 11px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
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
  color: var(--tx-475569);
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
  gap: 8px;
  color: var(--tx-0f172a);
  font-size: 1.25rem;
  font-weight: 800;
  line-height: 1.1;
  white-space: nowrap;
}

.market-map-kpi-copy em {
  padding: 4px 10px;
  border-radius: 7px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  font-size: 0.75rem;
  font-style: normal;
  font-weight: 700;
}

.market-map-kpi.is-available .market-map-kpi-icon,
.market-map-kpi.is-available .market-map-kpi-copy em { background: var(--bg-ecfdf5); color: var(--tx-059669); }
.market-map-kpi.is-reserved .market-map-kpi-icon,
.market-map-kpi.is-reserved .market-map-kpi-copy em { background: var(--bg-fff7ed); color: var(--tx-ea580c); }
.market-map-kpi.is-occupied .market-map-kpi-icon,
.market-map-kpi.is-occupied .market-map-kpi-copy em { background: var(--bg-f1f5f9); color: var(--tx-475569); }
.market-map-kpi.is-action-required {
  background: var(--bg-eef2ff);
  border: 1px solid var(--bd-c7d2fe);
  box-shadow: 0 8px 24px rgba(79, 70, 229, 0.14);
}
.market-map-kpi.is-action-required:hover {
  background: var(--bg-e0e7ff);
  border-color: var(--bd-818cf8);
  box-shadow: 0 10px 28px rgba(79, 70, 229, 0.2);
}
.market-map-kpi.is-action-required .market-map-kpi-icon {
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
  color: #ffffff;
  box-shadow: 0 4px 10px rgba(79, 70, 229, 0.28);
}
.market-map-kpi.is-action-required .market-map-kpi-copy small {
  color: var(--tx-4338ca);
  font-weight: 700;
}
.market-map-kpi.is-action-required .market-map-kpi-copy strong {
  color: var(--tx-312e81);
}
.market-map-kpi.is-action-required .market-map-kpi-copy em {
  background: var(--bg-ffffff);
  color: var(--tx-4338ca);
  border: 1px solid var(--bd-c7d2fe);
}
.market-map-kpi.is-action-required.active {
  border: 2px solid #4f46e5;
  background: var(--bg-e0e7ff);
  box-shadow: 0 10px 26px rgba(79, 70, 229, 0.26);
}
.market-map-kpi.is-legal-case .market-map-kpi-icon,
.market-map-kpi.is-legal-case .market-map-kpi-copy em { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.market-map-kpi.is-overdue .market-map-kpi-icon,
.market-map-kpi.is-overdue .market-map-kpi-copy em { background: var(--bg-fef2f2); color: var(--tx-dc2626); }
.market-map-kpi.is-expiring .market-map-kpi-icon,
.market-map-kpi.is-expiring .market-map-kpi-copy em { background: var(--bg-f5f3ff); color: var(--tx-7c3aed); }

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
  background: linear-gradient(90deg, var(--bg-f1f5f9) 25%, var(--bg-e2e8f0) 50%, var(--bg-f1f5f9) 75%);
  background-size: 200% 100%;
  animation: kpiSkeletonShimmer 1.35s ease-in-out infinite;
}

.ov-dash-skeleton-card {
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-height: 150px;
}
.dash-skel-title { width: 140px; height: 18px; }
.dash-skel-donut {
  align-self: center;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 14px solid var(--bg-f1f5f9);
  box-sizing: border-box;
}
.dash-skel-row { display: flex; align-items: center; gap: 10px; }
.dash-skel-dot { width: 12px; height: 12px; border-radius: 50%; background: var(--bg-e2e8f0); flex: 0 0 auto; }
.dash-skel-line { flex: 1 1 auto; width: auto; height: 14px; }

.kpi-skeleton-pill {
  display: inline-block;
  width: 32px;
  height: 18px;
  border-radius: 6px;
  background: linear-gradient(90deg, var(--bg-f1f5f9) 25%, var(--bg-e2e8f0) 50%, var(--bg-f1f5f9) 75%);
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
  border: 1px solid var(--bd-e2e8f0);
  border-right: 0;
  border-radius: 14px 0 0 14px;
  background: var(--bg-rgba-255-255-255-0p97);
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
  border: 1px solid var(--bd-bfdbfe);
  border-radius: 10px;
  background: var(--bg-ffffff);
  color: var(--tx-1e293b);
  cursor: text;
}

.market-map-tool-rail .gis-map-toolbar :deep(.map-search-input::placeholder) { color: transparent; }
.market-map-tool-rail .gis-map-toolbar :deep(.search-input-group:focus-within .map-search-input::placeholder) { color: var(--tx-94a3b8); }
.market-map-tool-rail .gis-map-toolbar :deep(.search-icon) { left: 17px; color: var(--tx-334155); }

.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-filter-btn),
.market-map-view-button {
  width: 52px;
  height: 52px;
  padding: 0;
  display: grid;
  place-items: center;
  overflow: hidden;
  border: 0;
  border-top: 1px solid var(--bd-eef2f6);
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
  color: var(--tx-334155);
}

.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-filter-btn:hover),
.market-map-tool-rail .gis-map-toolbar :deep(.map-advanced-filter-btn.active),
.market-map-view-button:hover,
.market-map-view-button.active {
  background: var(--bg-eff6ff);
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

.market-map-view-button.active svg { color: var(--tx-2563eb); }

.is-fullscreen-map-panel :deep(.map-project-guidance),
.is-fullscreen-map-panel :deep(.map-floating-floor-preview-pill) {
  top: var(--mp-top) !important;
}

/* ซ่อนแถบการ์ดเฉพาะ iPad/มือถือ (ไม่เกิน 1024px) — Desktop ที่ย่อหน้าต่างยังแสดงการ์ดเสมอ */
@media (max-width: 1024px) {
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

@media (min-width: 769px) and (max-width: 1024px) {
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
    border: 1px solid var(--bd-e2e8f0) !important;
    border-bottom: 0 !important;
    box-shadow: 0 -8px 32px rgba(15, 23, 42, 0.18) !important;
    z-index: 200 !important;
    display: flex !important;
    flex-direction: column !important;
    background: var(--bg-ffffff) !important;
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
    background: var(--bg-ffffff) !important;
    border-bottom: 1px solid var(--bd-f1f5f9) !important;
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
    background: var(--bg-ffffff) !important;
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
    background: var(--bg-f8fafc) !important;
    overflow: hidden !important;
  }

  /* ส่วนหัวคงที่ด้านบน (Sticky Header) */
  .inspect-selected-stall-drawer .issd-header,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-header) {
    display: flex !important;
    flex-direction: column !important;
    gap: 8px !important;
    padding: 14px 16px 10px !important;
    background: var(--bg-ffffff) !important;
    border-bottom: 1px solid var(--bd-eef2f6) !important;
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
    border-bottom: 1px solid var(--bd-e8eef6) !important;
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
    border: 1px solid var(--bd-bfdbfe) !important;
    background: var(--bg-eff6ff) !important;
    color: var(--tx-2563eb) !important;
    font-size: 0.875rem !important;
    font-weight: 600 !important;
    cursor: pointer !important;
    box-shadow: none !important;
    margin: 0 !important;
    -webkit-tap-highlight-color: transparent !important;
  }

  .inspect-selected-stall-drawer .issd-return-btn:hover,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-return-btn:hover) {
    background: var(--bg-dbeafe) !important;
  }

  .issd-mobile-only {
    display: grid !important;
    place-items: center !important;
    width: 38px !important;
    height: 38px !important;
    border-radius: 10px !important;
    border: 1px solid var(--bd-e2e8f0) !important;
    background: var(--bg-ffffff) !important;
    color: var(--tx-475569) !important;
    cursor: pointer !important;
  }
  .issd-mobile-only:hover {
    background: var(--bg-f1f5f9) !important;
    color: var(--tx-0f172a) !important;
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
    color: var(--tx-0f172a) !important;
  }

  .inspect-selected-stall-drawer .issd-project-sub,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-project-sub) {
    font-size: 0.8125rem !important;
    color: var(--tx-64748b) !important;
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
    background: var(--bg-ffffff) !important;
    border-bottom: 1px solid var(--bd-e2e8f0) !important;
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
    background: var(--bg-f1f5f9) !important;
    color: var(--tx-475569) !important;
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
    background: var(--bg-e2e8f0) !important;
    color: var(--tx-1e293b) !important;
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
    background: var(--bg-f8fafc) !important;
  }

  .inspect-selected-stall-drawer .issd-action-alert-card,
  .is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer .issd-action-alert-card) {
    margin: 12px 14px 0 !important;
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
    background: var(--bg-ffffff) !important;
    border-top: 1px solid var(--bd-e2e8f0) !important;
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
  .market-map-tool-rail { right: 10px; border-right: 1px solid var(--bd-e2e8f0); border-radius: 14px; }
  .is-fullscreen-map-panel :deep(.map-project-guidance),
  .is-fullscreen-map-panel :deep(.map-floating-floor-preview-pill) { top: 76px !important; }
}

/* Mini Drawer: Urgent Action Alert Banner & Expandable Card */
.issd-mini-alert-card {
  width: 100%;
  box-sizing: border-box;
  background: var(--bg-fffbeb);
  border: 1px solid var(--bd-fde68a);
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
  background: var(--bg-fef3c7);
}
.issd-mini-alert-head {
  display: flex;
  align-items: center;
  gap: 7px;
  color: var(--tx-b45309);
  font-size: 0.8125rem;
  font-weight: 700;
  min-width: 0;
}
.issd-mini-alert-head svg {
  color: var(--tx-d97706);
  flex-shrink: 0;
}
.issd-mini-alert-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--tx-d97706);
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
  border-top: 1px dashed var(--bd-fde68a);
  background: var(--bg-fffdf5);
}
.issd-mini-alert-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 8px 10px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-fef3c7);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
}
.issd-mini-alert-row:hover {
  background: var(--bg-fefce8);
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
  color: var(--tx-1e293b);
  font-weight: 600;
  line-height: 1.25;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.issd-mini-alert-desc {
  font-size: 0.75rem;
  color: var(--tx-64748b);
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
  color: var(--tx-d97706);
  white-space: nowrap;
  flex-shrink: 0;
  padding: 3px 8px;
  border-radius: 6px;
  background: var(--bg-fef3c7);
  border: 1px solid var(--bd-fde68a);
  transition: all 0.15s ease;
}
.issd-mini-alert-action-pill.is-danger {
  color: var(--tx-dc2626);
  background: var(--bg-fee2e2);
  border-color: var(--bd-fca5a5);
}
.issd-mini-alert-action-pill.is-warning {
  color: var(--tx-d97706);
  background: var(--bg-fef3c7);
  border-color: var(--bd-fde68a);
}
.issd-mini-alert-action-pill.is-info {
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
  border-color: var(--bd-bfdbfe);
}

/* Mini Drawer: Primary Occupant Card */
.issd-card-section {
  padding: 14px 16px;
  border-radius: 12px;
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  display: flex;
  flex-direction: column;
  gap: 10px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}
.issd-card-section.is-occupied {
  border-color: var(--bd-cbd5e1);
  background: var(--bg-f8fafc);
}
.issd-card-section.is-reserved {
  border-color: var(--bd-fde68a);
  background: var(--bg-fffdf5);
}
.issd-card-section.is-available {
  border-color: var(--bd-bbf7d0);
  background: var(--bg-f0fdf4);
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
  color: var(--tx-64748b);
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
  background: var(--bg-dcfce7);
  color: var(--tx-15803d);
}
.issd-badge.badge-warning {
  background: var(--bg-fef3c7);
  color: var(--tx-b45309);
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
  color: var(--tx-0f172a);
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
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-bfdbfe);
  color: var(--tx-2563eb);
  font-size: 0.75rem;
  font-weight: 600;
  text-decoration: none;
  flex-shrink: 0;
  transition: background 0.15s ease, border-color 0.15s ease;
}
.issd-phone-link:hover {
  background: var(--bg-dbeafe);
  border-color: var(--bd-93c5fd);
}
.issd-occupant-meta {
  display: flex;
  flex-direction: column;
  gap: 7px;
  padding-top: 10px;
  border-top: 1px solid var(--bd-e2e8f0);
}
.issd-meta-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  font-size: 0.8125rem;
}
.issd-meta-item .meta-k {
  color: var(--tx-64748b);
  font-weight: 500;
  flex-shrink: 0;
}
.issd-meta-item .meta-v {
  color: var(--tx-1e293b);
  font-weight: 600;
  text-align: right;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-meta-item .code-font {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  color: var(--tx-2563eb);
}

.issd-contract-code-pill {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-dbeafe);
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
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  border: 1px solid var(--bd-dbeafe);
}
.issd-person-avatar.is-reserved {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
  border-color: var(--bd-fde68a);
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
  color: var(--tx-64748b);
  line-height: 1.2;
}
.issd-person-name {
  font-size: 0.9375rem;
  font-weight: 700;
  color: var(--tx-0f172a);
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
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 9px;
}
.issd-subdoc-box.is-reserved {
  border-color: var(--bd-fde68a);
  background: var(--bg-ffffff);
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
  color: var(--tx-64748b);
  font-weight: 500;
  flex-shrink: 0;
}
.issd-subdoc-row .subdoc-v {
  color: var(--tx-1e293b);
  font-weight: 600;
  text-align: right;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-subdoc-row .subdoc-v.is-code {
  font-weight: 700;
  color: var(--tx-2563eb);
}
.issd-subdoc-row .subdoc-duration {
  font-style: normal;
  font-weight: 600;
  color: var(--tx-2563eb);
  margin-left: 2px;
}
.issd-occupant-compact-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding-top: 8px;
  border-top: 1px solid var(--bd-e2e8f0);
  font-size: 0.75rem;
  color: var(--tx-64748b);
  font-weight: 500;
}
.issd-duration-tag {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  font-size: 0.75rem;
  font-weight: 700;
  white-space: nowrap;
}
.issd-sep-dot {
  margin: 0 5px;
  color: var(--tx-cbd5e1);
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
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
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
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
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
  color: var(--tx-64748b);
  font-weight: 500;
  line-height: 1.2;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-mini-ops-card .ops-info strong {
  font-size: 0.8125rem;
  color: var(--tx-1e293b);
  font-weight: 700;
  line-height: 1.25;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.issd-mini-ops-card .ops-info strong.text-success {
  color: var(--tx-15803d);
}
.issd-mini-ops-card .ops-info strong.text-warning {
  color: var(--tx-d97706);
}
.issd-mini-ops-card .ops-info strong.text-danger {
  color: var(--tx-dc2626);
}
.issd-mini-ops-card.is-ok .ops-icon {
  background: var(--bg-ecfdf5);
  color: var(--tx-10b981);
}
.issd-mini-ops-card.is-info-tint .ops-icon {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}
.issd-mini-ops-card.is-alert {
  background: var(--bg-fef2f2);
  border-color: var(--bd-fecaca);
}
.issd-mini-ops-card.is-alert .ops-icon {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
}
.issd-mini-ops-card.is-warning {
  background: var(--bg-fffbeb);
  border-color: var(--bd-fde68a);
}
.issd-mini-ops-card.is-warning .ops-icon {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
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
  border-color: var(--bd-cbd5e1);
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
  background: var(--bg-rgba-255-255-255-0p96);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--bd-rgba-226-232-240-0p95);
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
  color: var(--tx-475569);
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.18s ease;
  white-space: nowrap;
}

.view-toggle-btn:hover {
  background: var(--bg-f1f5f9);
  color: var(--tx-1e293b);
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
  background: var(--bg-rgba-255-255-255-0p96);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--bd-rgba-226-232-240-0p95);
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
  color: var(--tx-475569);
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.map-filter-chip:hover {
  background: var(--bg-f8fafc);
  color: var(--tx-0f172a);
}

.chip-label {
  font-size: 0.8125rem;
}

.chip-badge {
  font-size: 0.8125rem;
  font-weight: 700;
  color: var(--tx-0f172a);
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
  background: var(--bg-eff6ff);
  border-color: var(--bd-bfdbfe);
  color: var(--tx-1d4ed8);
}

.map-filter-chip.chip-demand-high.active {
  background: var(--bg-fef2f2);
  border-color: var(--bd-fca5a5);
  color: var(--tx-b91c1c);
}

.map-filter-chip.chip-demand-medium.active {
  background: var(--bg-fff7ed);
  border-color: var(--bd-fed7aa);
  color: var(--tx-c2410c);
}

.map-filter-chip.chip-demand-low.active {
  background: var(--bg-fefce8);
  border-color: var(--bd-fef08a);
  color: var(--tx-a16207);
}

.map-filter-chip.chip-demand-none.active {
  background: var(--bg-f0fdf4);
  border-color: var(--bd-bbf7d0);
  color: var(--tx-15803d);
}

.map-floating-overlay-bar .gis-map-toolbar {
  margin-left: auto;
  flex-shrink: 0;
}

.map-floating-overlay-bar .gis-map-toolbar :deep(.map-search-input) {
  background: var(--bg-rgba-255-255-255-0p95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid var(--bd-rgba-226-232-240-0p9);
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
}

.map-floating-overlay-bar .gis-map-toolbar :deep(.map-advanced-filter-btn) {
  background: var(--bg-rgba-255-255-255-0p95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid var(--bd-rgba-226-232-240-0p9);
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.08);
}

@media (max-width: 960px) {
  .map-floating-overlay-bar {
    position: static;
    padding: 10px 12px;
    background: var(--bg-ffffff);
    border-bottom: 1px solid var(--bd-e2e8f0);
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

/* การ์ดสิ่งที่ต้องดำเนินการ (หน้ารายละเอียดเต็ม): สรุปเป็นการ์ดเล็กหนึ่งแถวและเลื่อนแนวนอน */
.issd-full-main .issd-action-alert-card {
  display: block;
  min-height: 0;
  overflow: visible;
  background: transparent;
  border: 0;
  border-radius: 0;
  box-shadow: none;
}
.issd-action-card-rail {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  min-width: 0;
  padding: 0 2px;
}
.issd-action-summary-card {
  display: flex;
  align-items: stretch;
  gap: 9px;
  width: 100%;
  min-width: 0;
  min-height: 64px;
  padding: 9px 14px;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 10px;
  background: var(--bg-ffffff);
  color: var(--tx-334155);
  font-family: inherit;
  text-align: left;
  cursor: pointer;
  transition: border-color .15s ease, background .15s ease;
}
.issd-action-summary-card:last-child:nth-child(odd) { grid-column: 1 / -1; }
.issd-action-summary-card:hover { border-color: var(--bd-93c5fd); background: var(--bg-f8fbff); }
.issd-action-summary-card.is-danger { border-color: var(--bd-fecaca); background: var(--bg-fffafa); }
.issd-action-summary-card.is-warning { border-color: var(--bd-fde68a); background: var(--bg-fffdf5); }
.issd-action-summary-card.is-info { border-color: var(--bd-bfdbfe); background: var(--bg-f8fbff); }
.issd-action-summary-card.is-danger:hover { border-color: var(--bd-fca5a5); background: var(--bg-fef2f2); }
.issd-action-summary-card.is-warning:hover { border-color: #fcd34d; background: var(--bg-fffbeb); }
.issd-action-summary-card.is-info:hover { border-color: var(--bd-93c5fd); background: var(--bg-eff6ff); }
.issd-action-summary-card:focus-visible { outline: 3px solid rgba(37, 99, 235, .18); outline-offset: 2px; }
.issd-action-summary-copy { display: flex; flex: 1 1 auto; flex-direction: column; justify-content: center; gap: 3px; min-width: 0; }
.issd-action-summary-title { display: flex; align-items: center; gap: 7px; min-width: 0; }
.issd-action-summary-title strong { overflow: hidden; color: var(--tx-1e293b); font-size: .9375rem; font-weight: 700; text-overflow: ellipsis; white-space: nowrap; }
.issd-action-summary-copy small { overflow: hidden; color: var(--tx-64748b); font-size: .8125rem; font-weight: 500; text-overflow: ellipsis; white-space: nowrap; }
.issd-action-summary-chevron { align-self: center; flex: 0 0 auto; color: var(--tx-2563eb); }
.issd-aac-bar { flex: 0 0 4px; align-self: stretch; min-height: 34px; border-radius: 9999px; background: #94a3b8; }
.issd-aac-bar.is-danger { background: #dc2626; }
.issd-aac-bar.is-warning { background: #f59e0b; }
.issd-aac-bar.is-info { background: #2563eb; }
.issd-aac-row-text { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.issd-aac-row-text strong { display: inline-flex; align-items: center; gap: 8px; color: var(--tx-1e293b); font-size: 1rem; font-weight: 700; line-height: 1.3; }
.issd-aac-row-text small { color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; line-height: 1.3; }
.issd-aac-count { display: inline-grid; place-items: center; min-width: 22px; height: 22px; padding: 0 6px; border-radius: 9999px; color: #ffffff; font-size: 0.8125rem; font-weight: 700; line-height: 1; background: #94a3b8; }
.issd-aac-count.is-danger { background: #dc2626; }
.issd-aac-count.is-warning { background: #d97706; }
.issd-aac-count.is-info { background: #2563eb; }
@media (max-width: 768px) {
  .issd-action-card-rail { grid-template-columns: minmax(0, 1fr); padding: 0 2px; }
  .issd-action-summary-card:last-child:nth-child(odd) { grid-column: auto; }
}

/* การ์ดสิ่งที่ต้องดำเนินการ (หน้าสรุปแผง): แถวโปร่ง แถบสีบอกระดับ เหมือนหน้ารายละเอียดเต็ม */
.issd-mini-alert-card { background: var(--bg-ffffff); border: 1px solid var(--bd-fde68a); border-radius: 14px; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04); }
.issd-mini-alert-card .issd-mini-alert-banner { background: var(--bg-fffbeb); padding: 12px 16px; }
.issd-mini-alert-card .issd-mini-alert-list { padding: 0; background: var(--bg-ffffff); border-top: 1px solid var(--bd-fde68a); }
.issd-mini-alert-list .issd-mini-alert-row { display: flex; align-items: center; justify-content: flex-start; gap: 12px; padding: 12px 16px; border: 0; border-top: 1px solid var(--bd-f1f5f9); border-radius: 0; background: var(--bg-ffffff); }
.issd-mini-alert-list .issd-mini-alert-row:first-child { border-top: 0; }
.issd-mini-alert-list .issd-mini-alert-row:hover { background: var(--bg-f8fafc); }
.issd-mini-alert-list .issd-aac-bar { flex: 0 0 4px; align-self: stretch; min-height: 34px; }
.issd-mini-alert-list .issd-aac-row-text { flex: 1 1 auto; }
.issd-mini-alert-chevron { flex: 0 0 auto; color: var(--tx-2563eb); }

/* ตัวกรองสัญญาในแท็บชำระเงิน: ย่อตามความกว้างแผง ไม่ล้นขอบ */
.inspect-selected-stall-drawer .issd-full-main .issd-tab-header-row.is-inline-controls { flex-wrap: wrap; }
.inspect-selected-stall-drawer .issd-full-main .issd-contract-inline-filter { flex: 1 1 220px; min-width: 0; max-width: 100%; justify-content: flex-end; }
.inspect-selected-stall-drawer .issd-full-main .issd-inline-filter-select { flex: 1 1 140px; width: auto; min-width: 0; max-width: 240px; }
.inspect-selected-stall-drawer .issd-full-main .issd-inline-filter-select :deep(.app-select-trigger) { min-width: 0; }

/* สวิตช์โหมด ข้อมูล / ค่าใช้จ่าย: ปรับให้เรียบ สมส่วน ตัวอักษรชัด */
.map-mode-float { padding: 3px; border: 1px solid var(--bd-rgba-203-213-225-0p8); border-radius: 12px; background: var(--bg-rgba-255-255-255-0p97); box-shadow: 0 4px 14px rgba(15, 23, 42, 0.1); }
.map-mode-float :deep(.app-segmented-control) { gap: 2px; padding: 0; border: 0; background: transparent; }
.map-mode-float :deep(.app-segmented-option) { min-width: 116px; min-height: 38px; padding: 0 16px; gap: 8px; justify-content: center; border-radius: 9px; font-size: 0.9375rem; font-weight: 600; color: var(--tx-475569); }
.map-mode-float :deep(.app-segmented-option svg) { width: 18px; height: 18px; flex: 0 0 18px; }
.map-mode-float :deep(.app-segmented-option:hover:not(.is-active)) { background: var(--bg-f1f5f9); color: var(--tx-1e293b); }
.map-mode-float :deep(.app-segmented-option.is-active) { background: #2563eb; color: #ffffff; font-weight: 700; box-shadow: none; }
@media (max-width: 720px) { .map-mode-float :deep(.app-segmented-option) { min-width: 96px; padding: 0 12px; } }



/* สวิตช์โหมดแบบการ์ดเตี้ยเท่าการ์ด KPI: อยู่แถวเดียวกัน ไอคอนบน ชื่อล่าง ตัวที่เลือกเป็นพื้นน้ำเงิน */
.map-mode-float {
  top: 18px;
  left: 20px;
  height: 72px;
  box-sizing: border-box;
  padding: 6px;
  border: 1px solid var(--bd-rgba-226-232-240-0p94);
  border-radius: 15px;
  background: var(--bg-rgba-255-255-255-0p96);
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
  color: var(--tx-64748b);
  font-size: 0.875rem;
  font-weight: 600;
  line-height: 1.1;
}
.map-mode-float :deep(.app-segmented-option svg) { width: 22px; height: 22px; flex: 0 0 22px; }
.map-mode-float :deep(.app-segmented-option:hover:not(.is-active)) { background: var(--bg-f1f5f9); color: var(--tx-1e293b); }
.map-mode-float :deep(.app-segmented-option.is-active) { background: #2563eb; color: #ffffff; font-weight: 700; box-shadow: 0 3px 10px rgba(37, 99, 235, 0.3); }
@media (max-width: 720px) {
  .map-mode-float { top: 10px; left: 10px; height: 56px; }
  .map-mode-float :deep(.app-segmented-option) { min-width: 72px; }
  .map-mode-float :deep(.app-segmented-option svg) { width: 18px; height: 18px; flex-basis: 18px; }
  .manage-page.is-map-page { --mp-top: 76px !important; --mp-top-2row: 76px !important; }
}

/* โครงว่างระหว่างรอกล่องขยายจบ */
.issd-full-skeleton { display: flex; flex-direction: column; gap: 14px; padding: 18px 22px; }
.issd-full-skeleton span { display: block; border-radius: 14px; background: linear-gradient(90deg, var(--bg-eef2f6) 25%, var(--bg-f8fafc) 50%, var(--bg-eef2f6) 75%); background-size: 200% 100%; animation: issdSkeleton 1.1s ease-in-out infinite; }
.issd-full-skeleton span:nth-child(1) { height: 56px; }
.issd-full-skeleton span:nth-child(2) { height: 150px; }
.issd-full-skeleton span:nth-child(3) { height: 90px; }
@keyframes issdSkeleton { from { background-position: 200% 0; } to { background-position: -200% 0; } }

/* การ์ด KPI ปรับตามความกว้างของแถบการ์ดเอง (ไม่ใช่ความกว้างจอ) จึงใช้ได้ทั้งตอนย่อหน้าต่างและตอนเปิดแผงข้าง */
.market-map-kpi-strip { container-type: inline-size; container-name: kpi; }
.market-map-kpi-copy small { max-width: 100%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
/* ป้ายชื่อการ์ดบรรทัดเดียวเสมอ (เดิมสำรองที่ 2 บรรทัดทำให้ห่างจากตัวเลขเกินไป) */
.market-map-kpi-copy > small:first-child {
  display: block;
  min-height: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  -webkit-line-clamp: unset;
}
.market-map-kpi-copy { gap: 5px; } /* ระยะระหว่างชื่อการ์ดกับตัวเลข */
@container kpi (max-width: 1130px) {
  /* ยังไม่เปิดแผงข้าง: เก็บไอคอนไว้และลดระยะในการ์ด โดยคงป้ายหน่วย "แผง" ไว้ */
  .market-map-kpi { padding: 10px 6px 10px 10px; gap: 12px; }
  .market-map-kpi.active { padding: 9px 5px 9px 9px; }
  .market-map-kpi-copy small { font-size: 0.75rem; } /* 12px = ขนาดต่ำสุดของระบบ ให้ "ต้องดำเนินการ" แสดงครบในพื้นที่ที่เหลือ */
  /* คงรูปแบบป้ายหน่วยให้เห็นพื้นหลังชัดเจน แม้ยังไม่ได้เปิดแผงด้านข้าง */
  .market-map-kpi-strip:not(.has-side-panel) .market-map-kpi .market-map-kpi-copy strong { font-size: 1.25rem; gap: 6px; }
  .market-map-kpi-strip:not(.has-side-panel) .market-map-kpi .market-map-kpi-copy strong em {
    display: inline-block;
    padding: 4px 10px;
    font-size: 0.75rem;
    font-weight: 700;
    align-self: flex-end;
  }
  /* เปิดแผงข้างแล้ว แถบแคบลง: ซ่อนไอคอนแทน แล้วคืนป้ายหน่วย "แผง" */
  .market-map-kpi-strip.has-side-panel .market-map-kpi { padding: 10px 12px; gap: 0; }
  .market-map-kpi-strip.has-side-panel .market-map-kpi.active { padding: 9px 11px; }
  .market-map-kpi-strip.has-side-panel .market-map-kpi-icon { display: none; }
  .market-map-kpi-strip.has-side-panel .market-map-kpi-copy small { font-size: 0.8125rem; }
  .market-map-kpi-strip.has-side-panel .market-map-kpi-copy strong { font-size: 1.25rem; line-height: 1.15; }
  .market-map-kpi-strip.has-side-panel .market-map-kpi-copy strong em { display: block; }
}
@container kpi (max-width: 1060px) {
  /* แคบจริงๆ แม้ยังไม่เปิดแผงข้าง (เช่นตอนเมนูข้างขยาย): ซ่อนไอคอนด้วย ให้ตัวเลขกับป้ายหน่วยไม่ถูกตัด */
  .market-map-kpi-icon { display: none; }
  .market-map-kpi { padding: 10px 12px; gap: 0; }
  /* ไอคอนหายแล้วมีที่เหลือ: คืนขนาดตัวเลขและชื่อให้อ่านง่ายขึ้น */
  .market-map-kpi-strip:not(.has-side-panel) .market-map-kpi .market-map-kpi-copy strong { font-size: 1.25rem; line-height: 1.15; }
  .market-map-kpi-copy small { font-size: 0.8125rem; }
}
@container kpi (max-width: 800px) {
  /* แคบมาก: ซ่อนป้ายหน่วย "แผง" เหลือชื่อและตัวเลข */
  .market-map-kpi { padding: 8px 10px; }
  .market-map-kpi-copy small { font-size: 0.8125rem; }
  .market-map-kpi-copy strong { font-size: 1.25rem; }
  .market-map-kpi-copy strong em { display: none !important; }
}

/* แถบการ์ดอยู่กึ่งกลางเสมอและเว้นทางซ้ายพอสำหรับตัวเลือกมุมมองเต็ม ไม่ให้ชนเมื่อย่อหน้าต่าง */
/* แถบการ์ดอยู่กึ่งกลางของพื้นที่ที่เหลือ โดยเว้นซ้ายให้ตัวเลือกมุมมอง (ไม่ซ้อนกันทุกความกว้างจอ) */
/* การ์ดด้านบนอยู่กึ่งกลางพื้นที่แผนที่ที่มองเห็น (เริ่มหลังแถบเมนู 68px ถึงขอบขวาของแผนที่) เว้นสองข้างเท่ากัน 190px */
.market-map-kpi-strip:not(.has-side-panel) { left: calc(50% + 34px); right: auto; width: min(1560px, calc(100% - 320px)); max-width: 1560px; margin: 0; transform: translateX(-50%); }

/* Desktop (กว้างกว่า 1024px) ตอนเปิดแผงข้างทางขวา (ตัวกรอง/ค้นหา/ชั้นข้อมูล/แดชบอร์ด): ตัวเลือกมุมมองย่อเหลือเฉพาะไอคอน */
/* iPad และมือถือ (ไม่เกิน 1024px) ไม่ย่อ แสดงชื่อเหมือนเดิม */
@media (min-width: 1025px) {
  .map-mode-dropdown-wrap.is-compact .map-mode-trigger-btn { justify-content: center; width: 48px; padding: 0; gap: 0; }
  .map-mode-dropdown-wrap.is-compact .map-mode-trigger-btn .trigger-text-group,
  .map-mode-dropdown-wrap.is-compact .map-mode-trigger-btn .trigger-chevron { display: none; }
  /* แถบการ์ดอยู่กึ่งกลางของพื้นที่ที่เหลือทางซ้ายของแผงข้างพอดี (ตัวเลือกมุมมองย่อแล้วจึงไม่ต้องเว้นที่ชดเชยไปทางขวา) */
  .market-map-kpi-strip.has-side-panel {
    left: calc((100% - var(--side-w, 380px)) / 2);
    width: min(1050px, calc(100% - var(--side-w, 380px) - 170px));
  }
}

/* Desktop จอไม่กว้าง (1025–1280px): แถบการ์ดไม่เว้นสองข้างเท่ากัน แต่ใช้พื้นที่ทางขวาของปุ่มมุมมองให้เต็ม เพื่อให้การ์ดทั้ง 6 ใบยังอ่านได้ */
@media (min-width: 1025px) and (max-width: 1280px) {
  .market-map-kpi-strip:not(.has-side-panel) { left: calc(50% + 34px); transform: translateX(-50%); width: calc(100% - 320px); max-width: none; }
  /* เปิดแผงข้างบนจอเล็ก: การ์ดคงขนาดอ่านได้ (ไม่บีบ) แล้วเลื่อนแนวนอนดูใบที่เหลือได้ */
  .market-map-kpi-strip.has-side-panel {
    left: 85px;
    transform: none;
    width: calc(100% - var(--side-w, 380px) - 101px);
    max-width: none;
    grid-template-columns: repeat(6, minmax(118px, 1fr));
    overflow-x: auto;
    overscroll-behavior-x: contain;
    pointer-events: auto;
    scrollbar-width: none;
  }
  .market-map-kpi-strip.has-side-panel::-webkit-scrollbar { display: none; }
}

/* รายละเอียดแผงแบบเต็มที่กว้างน้อย (เช่นตอนมีแผงขวาเปิดอยู่): แถบเมนูกลายเป็นแถบแท็บบนสุด เนื้อหาใช้พื้นที่เต็ม */
.inspect-selected-stall-drawer.is-left-dock.is-full-detail-mode { container-type: inline-size; container-name: detail; }
@container detail (max-width: 600px) {
  .inspect-selected-stall-drawer .issd-full-layout { flex-direction: column; }
  .inspect-selected-stall-drawer .issd-full-sidebar { width: 100%; border-right: 0; border-bottom: 1px solid var(--bd-eef2f6); padding: 0 8px; overflow-x: auto; overflow-y: hidden; scrollbar-width: none; }
  .inspect-selected-stall-drawer .issd-full-sidebar::-webkit-scrollbar { display: none; }
  .inspect-selected-stall-drawer .issd-sidebar-nav { flex-direction: row; gap: 2px; white-space: nowrap; }
  .inspect-selected-stall-drawer .issd-sidebar-tab { width: auto; flex: 0 0 auto; padding: 11px 10px; border: 0; border-radius: 0; border-bottom: 3px solid transparent; background: transparent; box-shadow: none; }
  .inspect-selected-stall-drawer .issd-sidebar-tab.active { background: transparent; color: var(--tx-2563eb); border-bottom-color: #2563eb; }
  .inspect-selected-stall-drawer .issd-sidebar-tab svg { display: none; }
  .inspect-selected-stall-drawer .issd-kpi-grid.is-three { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .inspect-selected-stall-drawer .issd-full-main .issd-body-full { padding: 14px 14px 18px; }
  .inspect-selected-stall-drawer .issd-full-main .issd-footer-full { padding: 12px 14px 14px; }
  .inspect-selected-stall-drawer .issd-full-main .issd-tab-header-row.is-inline-controls { flex-wrap: wrap; }
}

/* ตัวเลขสถานะบนแดชบอร์ดที่กดกรองแผนที่ได้: หน้าตาเป็นปุ่ม (มีกรอบ ชี้แล้วเปลี่ยนสี) */
.ov-dash-hint { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.ov-dash-mode-view .is-filterable { cursor: pointer; border: 1px solid var(--bd-e2e8f0); border-radius: 10px; transition: background-color .15s ease, border-color .15s ease, box-shadow .15s ease; }
.ov-contract-health-legend span.is-filterable,
.ov-maintenance-legend span.is-filterable { padding: 6px 10px; background: var(--bg-ffffff); }
.ov-contract-renewal-row span.is-filterable { padding: 8px 4px; background: var(--bg-ffffff); }
.ov-stage-flow > div.is-filterable { border-color: var(--bd-e2e8f0); }
.ov-dash-mode-view .is-filterable:hover { background-color: var(--bg-eff6ff); border-color: #93c5fd; box-shadow: 0 1px 4px rgba(37, 99, 235, 0.15); }
.ov-dash-mode-view .is-filterable:focus-visible { outline: 2px solid #93c5fd; outline-offset: 1px; }
.ov-dash-mode-view .is-filterable.is-selected { background-color: var(--bg-eff6ff); border-color: #2563eb; box-shadow: inset 0 0 0 1px #2563eb; }

/* แผงยอดนิยมบนแดชบอร์ด: ตารางอันดับ แถวละ 2 บรรทัด (ชื่อแผง/โครงการ) เหรียญ 1-3 จำนวนครั้งอยู่ขวา กดแล้วเปิดแผงนั้นบนแผนที่ */
.ov-popular-toolbar { display: flex; align-items: center; gap: 6px; flex-wrap: nowrap; margin: -2px 0 6px; overflow-x: auto; scrollbar-width: none; }
.ov-popular-toolbar::-webkit-scrollbar { display: none; }
.ov-popular-range { flex: 0 0 auto; padding: 3px 10px; border: 1px solid var(--bd-e2e8f0); border-radius: 999px; background: var(--bg-ffffff); color: var(--tx-64748b); font: inherit; font-size: 0.8125rem; font-weight: 600; white-space: nowrap; cursor: pointer; transition: background-color .15s ease, color .15s ease, border-color .15s ease; }
.ov-popular-range:hover { background: var(--bg-f8fafc); color: var(--tx-0f172a); }
.ov-popular-range.is-active { background: var(--bg-eff6ff); border-color: #93c5fd; color: var(--tx-2563eb); }
/* แผงยอดนิยม: แถวอันดับ 2 บรรทัด (ชื่อแผง/โครงการ) เหรียญ 1-3 จำนวนครั้งอยู่ขวา */
.ov-popular-list { list-style: none; margin: 4px 0 0; padding: 0; display: flex; flex-direction: column; }
.ov-popular-list > li + li { border-top: 1px solid var(--bd-f1f5f9); }
.ov-popular-row { width: 100%; display: grid; grid-template-columns: 30px minmax(0, 1fr) auto; align-items: center; gap: 14px; padding: 10px 8px; border: 0; border-radius: 10px; background: transparent; color: var(--tx-0f172a); font: inherit; text-align: left; cursor: pointer; transition: background-color .15s ease; }
.ov-popular-row:hover { background: var(--bg-f8fafc); }
.ov-popular-row.is-selected { background: var(--bg-eff6ff); box-shadow: inset 0 0 0 1.5px #2563eb; }
.ov-popular-row:focus-visible { outline: 2px solid #93c5fd; outline-offset: 1px; }
.ov-popular-rank { width: 30px; height: 30px; display: grid; place-items: center; border-radius: 50%; background: var(--bg-f1f5f9); color: var(--tx-475569); font-size: 0.875rem; font-weight: 800; }
.ov-popular-row.is-rank-1 .ov-popular-rank { background: #f59e0b; color: #ffffff; box-shadow: 0 2px 6px rgba(245, 158, 11, 0.3); }
.ov-popular-row.is-rank-2 .ov-popular-rank { background: #94a3b8; color: #ffffff; }
.ov-popular-row.is-rank-3 .ov-popular-rank { background: #c2713a; color: #ffffff; }
.ov-popular-main { min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.ov-popular-main strong { font-size: 1rem; font-weight: 700; color: var(--tx-0f172a); line-height: 1.35; }
.ov-popular-main small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; line-height: 1.35; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.ov-popular-count { display: flex; align-items: baseline; gap: 4px; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; white-space: nowrap; }
.ov-popular-count b { color: var(--tx-0f172a); font-size: 1.25rem; font-weight: 800; }

/* แดชบอร์ด: หัวข้อที่เป็นขั้นตอน (ส่งคืนแผง / คดี) แสดงเป็นช่องเรียงตามลำดับ (ไม่มีลูกศร) ตัวเลข 0 จางลง */
.ov-dash-section-meta { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.ov-stage-flow { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 8px; }
.ov-stage-flow > div { position: relative; min-width: 0; padding: 10px 4px 9px; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 2px; border: 1px solid var(--bd-eef2f6); border-top: 3px solid var(--bd-cbd5e1); border-radius: 10px; background: var(--bg-f8fafc); text-align: center; }
.ov-stage-flow > div > span { color: var(--tx-64748b); font-size: 0.75rem; font-weight: 500; line-height: 1.25; }
.ov-stage-flow > div > strong { color: var(--tx-1e293b); font-size: 1.25rem; font-weight: 800; line-height: 1.2; }
.ov-stage-flow > div.is-pending { border-top-color: #f59e0b; }
.ov-stage-flow > div.is-success { border-top-color: #16a34a; }
.ov-stage-flow > div.is-danger { border-top-color: #ef4444; }
.ov-stage-flow > div.is-zero > strong { opacity: 0.4; font-weight: 700; }
.ov-stage-alert { margin-top: 10px; display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 9px 12px; border: 1px solid var(--bd-eef2f6); border-radius: 10px; background: var(--bg-f8fafc); color: var(--tx-475569); font-size: 0.8125rem; font-weight: 600; }
.ov-stage-alert strong { color: var(--tx-dc2626); font-size: 1.0625rem; font-weight: 800; }
.ov-stage-alert.is-zero { opacity: 0.6; }
.ov-stage-alert.is-zero strong { color: var(--tx-64748b); }
.ov-stage-alert.is-selected, .ov-stage-flow > div.is-selected { border-color: #2563eb; background-color: var(--bg-eff6ff); }
/* แถบสัดส่วนการต่อสัญญา */
.ov-renewal-bar { margin-top: 12px; display: flex; height: 8px; overflow: hidden; border-radius: 999px; background: var(--bg-e2e8f0); }
.ov-renewal-bar i { display: block; height: 100%; transition: width .3s ease; }
.ov-renewal-bar .is-renewed { background: #16a34a; }
.ov-renewal-bar .is-awaiting { background: #f59e0b; }
.ov-renewal-bar .is-declined { background: #94a3b8; }
.ov-renewal-bar + .ov-contract-renewal-row { margin-top: 10px; }
.ov-popular-note { margin: 8px 0 0; padding: 8px 10px; border-radius: 8px; background: var(--bg-f8fafc); color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; line-height: 1.4; }

/* ชื่อแกนตั้งของกราฟภาพรวมโครงการ = จำนวนแผง (แกนนอนเป็นชื่อโครงการใต้แท่งอยู่แล้ว) */
.ov-chart-axis-title { color: var(--tx-64748b); font-size: 0.75rem; font-weight: 600; }
.ov-chart-axis-title.is-y { margin: 14px 0 10px 2px; }

/* ภาพรวมโครงการ: กราฟแท่งแนวนอน */
.ov-hbar-total { color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; white-space: nowrap; font-variant-numeric: tabular-nums; }
.ov-hbar-chips { grid-column: 1 / -1; display: flex; gap: 6px; }
.ov-hbar-chip { display: inline-flex; align-items: baseline; justify-content: center; min-width: 52px; padding: 2px 10px; border-radius: 7px; font-size: 0.8125rem; font-weight: 600; white-space: nowrap; font-variant-numeric: tabular-nums; }
.ov-hbar-chip strong { font-size: 0.9375rem; font-weight: 700; }
.ov-hbar-chip.is-available { background: rgba(16, 185, 129, 0.14); color: var(--tx-059669); }
.ov-hbar-chip.is-reserved { background: rgba(245, 158, 11, 0.16); color: #d97706; }
.ov-hbar-chip.is-occupied { background: rgba(71, 85, 105, 0.14); color: var(--tx-475569); }
/* ค่า 0: คงสีตามประเภทไว้ แต่จางลง */
.ov-hbar-chip.is-zero { opacity: 0.55; }
.ov-hbar-chip.is-zero strong { font-weight: 600; }
.ov-hbar-key > span { display: inline-flex; align-items: center; gap: 6px; }
.ov-hbar-key i { width: 10px; height: 10px; border-radius: 3px; }
.ov-hbar-key .is-available i { background: #10b981; }
.ov-hbar-key .is-reserved i { background: #f59e0b; }
.ov-hbar-key .is-occupied i { background: #475569; }
.ov-hbar-scroll { max-height: 440px; margin-top: 8px; overflow-y: auto; overscroll-behavior: contain; scrollbar-width: thin; }
.ov-hbar-key { position: sticky; top: 0; z-index: 2; display: flex; flex-wrap: wrap; gap: 4px 16px; padding: 8px 4px 10px; background: var(--bg-ffffff); border-bottom: 1px solid var(--bd-f1f5f9); color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.ov-hbar-list { list-style: none; margin: 0; padding: 12px 4px 2px 0; }
.ov-hbar-row { width: 100%; display: grid; grid-template-columns: minmax(0, 1fr) auto; align-items: center; gap: 7px 12px; padding: 12px 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; background: var(--bg-f8fafc); color: inherit; font: inherit; text-align: left; cursor: pointer; transition: background-color .15s ease, border-color .15s ease; }
.ov-hbar-row:hover { background: var(--bg-ffffff); border-color: #93c5fd; }
.ov-hbar-list li + li { margin-top: 10px; }
.ov-hbar-row:focus-visible { outline: 2px solid #93c5fd; outline-offset: 1px; }
.ov-hbar-row.active { background: var(--bg-eff6ff); border-color: #2563eb; }
.ov-hbar-name { min-width: 0; overflow: hidden; color: var(--tx-334155); font-size: 0.9375rem; font-weight: 600; line-height: 1.35; text-overflow: ellipsis; white-space: nowrap; }
.ov-hbar-row.active .ov-hbar-name { color: var(--tx-2563eb); font-weight: 700; }
.ov-hbar-track { grid-column: 1 / -1; display: flex; align-items: stretch; gap: 2px; height: 10px; overflow: hidden; border-radius: 999px; background: var(--bg-f1f5f9); }
.ov-hbar-track i { display: block; height: 100%; min-width: 2px; transition: width .3s ease; }
.ov-hbar-track .is-occupied { background: #475569; }
.ov-hbar-track .is-reserved { background: #f59e0b; }
.ov-hbar-track .is-available { background: #10b981; }
/* คำอธิบายสีของกราฟภาพรวมโครงการ: อยู่ใต้รายการ ตรึงไว้ไม่เลื่อนไปกับแท่ง */
.ov-project-column-legend.is-footer { margin: 10px 0 0; padding-top: 10px; border-top: 1px solid var(--bd-f1f5f9); justify-content: center; }

/* เส้นคั่นใต้หัวข้อของทุกการ์ดบนแดชบอร์ด */
.ov-dash-mode-view .ov-dash-section-head,
.ov-dash-mode-view .ov-stall-section-head { margin-bottom: 12px; padding-bottom: 10px; border-bottom: 1px solid var(--bd-f1f5f9); }

/* ติดตามสัญญาในการ์ดสัญญา (เจ้าหน้าที่): ต่อสัญญา / ส่งคืนแผง / คดี */
.issd-workflow { display: flex; flex-direction: column; gap: 10px; padding-top: 12px; border-top: 1px dashed var(--bd-e2e8f0); }
.issd-workflow-title { color: var(--tx-0f172a); font-size: 0.875rem; font-weight: 700; }
.issd-workflow-grid { display: grid; grid-template-columns: 1fr; gap: 10px; }
.issd-workflow-field { display: flex; flex-direction: column; gap: 4px; min-width: 0; }
.issd-workflow-field > span { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
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
  background: var(--bg-fffbeb);
  border: 1px solid var(--bd-fde68a);
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
  background: var(--bg-fffbeb);
  cursor: pointer;
  user-select: none;
  transition: background 0.15s ease;
}

.inspect-selected-stall-drawer .issd-aac-header:hover,
.issd-aac-header:hover {
  background: var(--bg-fef3c7);
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
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
  border: 1px solid var(--bd-fde68a);
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
  color: var(--tx-92400e);
}

.inspect-selected-stall-drawer .issd-aac-count-badge,
.issd-aac-count-badge {
  display: inline-flex;
  align-items: center;
  padding: 1px 7px;
  border-radius: 999px;
  background: #fde68a;
  color: var(--tx-78350f);
  font-size: 0.75rem;
  font-weight: 700;
}

.inspect-selected-stall-drawer .issd-aac-sub,
.issd-aac-sub {
  font-size: 0.75rem;
  color: var(--tx-b45309);
}

.inspect-selected-stall-drawer .issd-aac-toggle-btn,
.issd-aac-toggle-btn {
  width: 26px;
  height: 26px;
  border-radius: 6px;
  border: 1px solid rgba(217, 119, 6, 0.25);
  background: var(--bg-ffffff);
  color: var(--tx-92400e);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.inspect-selected-stall-drawer .issd-aac-toggle-btn:hover,
.issd-aac-toggle-btn:hover {
  background: var(--bg-fef3c7);
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
  background: var(--bg-ffffff);
}

.inspect-selected-stall-drawer .issd-aac-row,
.issd-aac-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  border-top: 1px solid var(--bd-f1f5f9);
  cursor: pointer;
  transition: background 0.15s ease;
}

.inspect-selected-stall-drawer .issd-aac-row:hover,
.issd-aac-row:hover {
  background: var(--bg-f8fafc);
}

.inspect-selected-stall-drawer .issd-aac-row:hover .issd-aac-action-text,
.issd-aac-row:hover .issd-aac-action-text {
  color: var(--tx-1d4ed8);
}

.inspect-selected-stall-drawer .issd-aac-row:hover .issd-aac-row-right svg,
.issd-aac-row:hover .issd-aac-row-right svg {
  transform: translateX(2px);
}

.inspect-selected-stall-drawer .issd-aac-row.is-danger:hover,
.issd-aac-row.is-danger:hover {
  background: var(--bg-fef2f2);
}

.inspect-selected-stall-drawer .issd-aac-row.is-warning:hover,
.issd-aac-row.is-warning:hover {
  background: var(--bg-fffbeb);
}

.inspect-selected-stall-drawer .issd-aac-row.is-info:hover,
.issd-aac-row.is-info:hover {
  background: var(--bg-eff6ff);
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
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
  border: 1px solid var(--bd-fecaca);
}

.inspect-selected-stall-drawer .issd-aac-tag.is-warning,
.issd-aac-tag.is-warning {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
  border: 1px solid var(--bd-fde68a);
}

.inspect-selected-stall-drawer .issd-aac-tag.is-info,
.issd-aac-tag.is-info {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  border: 1px solid var(--bd-bfdbfe);
}

.inspect-selected-stall-drawer .issd-aac-row-title,
.issd-aac-row-title {
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--tx-1e293b);
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
  color: var(--tx-2563eb);
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
  background: var(--bg-ffffff);
  animation: issdModeFadeIn 0.16s cubic-bezier(0.16, 1, 0.3, 1);
}

/* Left Sidebar Navigation */
.inspect-selected-stall-drawer .issd-full-sidebar,
.issd-full-sidebar {
  width: 152px;
  flex-shrink: 0;
  background: var(--bg-f8fafc);
  border-right: 1px solid var(--bd-eef2f6);
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
  color: var(--tx-475569);
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
  background: var(--bg-f1f5f9);
  color: var(--tx-0f172a);
}

.inspect-selected-stall-drawer .issd-sidebar-tab.active,
.issd-sidebar-tab.active {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  font-weight: 600;
  border-color: var(--bd-dbeafe);
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
  background: var(--bg-e2e8f0);
  color: var(--tx-475569);
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
  background: var(--bg-ffffff);
  position: relative;
  padding-top: 24px;
}
/* ปุ่มปิดมุมขวาบนของแผง: เว้นแถบบน 24px ไม่ให้ทับการ์ดแรก */
.issd-full-main > .issd-full-close { position: absolute; top: 6px; right: 10px; z-index: 4; }

.inspect-selected-stall-drawer .issd-full-main .issd-action-alert-card,
.issd-full-main .issd-action-alert-card {
  margin: 12px 18px 0;
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
  color: var(--tx-64748b);
  white-space: nowrap;
}

.inspect-selected-stall-drawer .issd-inline-filter-select,
.issd-inline-filter-select {
  width: 195px;
}

.inspect-selected-stall-drawer .issd-full-main .issd-footer-full,
.issd-full-main .issd-footer-full {
  padding: 12px 18px 16px;
  background: var(--bg-ffffff);
  border-top: 1px solid var(--bd-f1f5f9);
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
    border-bottom: 1px solid var(--bd-eef2f6);
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
.inspect-selected-stall-drawer .issd-full-main { background: var(--bg-f8fafc); }
.inspect-selected-stall-drawer .issd-full-main .issd-body-full { padding: 20px 22px 24px; gap: 16px; display: flex; flex-direction: column; }
.inspect-selected-stall-drawer .issd-full-main .issd-footer-full { padding: 14px 22px 16px; background: var(--bg-ffffff); }
.inspect-selected-stall-drawer .issd-full-sidebar { width: 168px; padding: 14px 10px; background: var(--bg-ffffff); }
.inspect-selected-stall-drawer .issd-sidebar-tab { padding: 10px 12px; font-size: 0.9375rem; font-weight: 600; color: var(--tx-475569); border-radius: 10px; }
.inspect-selected-stall-drawer .issd-sidebar-tab.active { background: var(--bg-eff6ff); border-color: transparent; color: var(--tx-2563eb); box-shadow: none; }

/* หัวข้อแท็บ + ตัวกรอง */
.inspect-selected-stall-drawer .issd-full-main .issd-tab-header-row { display: flex; flex-direction: row; align-items: center; justify-content: space-between; gap: 12px; margin: 0; text-align: left; }
.inspect-selected-stall-drawer .issd-full-main .issd-tab-title { margin: 0; color: var(--tx-1e293b); font-size: 1.125rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main .app-segmented-control { align-self: flex-start; background: var(--bg-ffffff); }
.inspect-selected-stall-drawer .issd-full-main .app-segmented-option { font-size: 0.875rem; min-height: 34px; padding: 5px 14px; }
.inspect-selected-stall-drawer .issd-full-main .app-segmented-option.is-active { background: var(--bg-eff6ff); color: var(--tx-2563eb); box-shadow: none; }
.inspect-selected-stall-drawer .issd-full-main .issd-scroll-list { display: flex; flex-direction: column; gap: 12px; }

/* การ์ดรายการทุกชนิดหน้าตาเดียวกัน */
.inspect-selected-stall-drawer .issd-full-main :is(.issd-booking-card, .issd-contract-card, .issd-subtenant-card, .issd-payment-card, .issd-maintenance-card) {
  display: flex; flex-direction: column; gap: 14px; padding: 16px 18px;
  background: var(--bg-ffffff); border: 1px solid var(--bd-eef2f6); border-radius: 14px; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
}
.inspect-selected-stall-drawer .issd-full-main .issd-contract-card.is-current { background: var(--bg-ffffff); border-color: var(--bd-eef2f6); box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04); }
.inspect-selected-stall-drawer .issd-full-main :is(.issd-card-top-row, .issd-booking-card-header) { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.inspect-selected-stall-drawer .issd-full-main .issd-card-title { color: var(--tx-1e293b); font-size: 1rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main :is(.issd-contract-num, .issd-booking-ref-code, .issd-pay-period strong) { color: var(--tx-1e293b); font-size: 1rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main :is(.detail-label, .issd-booking-role-label, .issd-maintenance-meta-grid span, .issd-pay-method) { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-full-main :is(.detail-value, .issd-maintenance-meta-grid strong) { color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 600; }
.inspect-selected-stall-drawer .issd-full-main .is-rent { color: var(--tx-10b981); font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main .issd-contract-details-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px 16px; }
.inspect-selected-stall-drawer .issd-full-main .issd-contract-detail-item { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.inspect-selected-stall-drawer .issd-full-main .issd-contract-detail-item.is-wide { grid-column: 1 / -1; }
.inspect-selected-stall-drawer .issd-full-main .issd-maintenance-meta-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 10px 16px; }
.inspect-selected-stall-drawer .issd-full-main .issd-maintenance-meta-grid > div { display: flex; flex-direction: column; gap: 2px; }
.inspect-selected-stall-drawer .issd-full-main .issd-card-actions-row { display: flex; justify-content: flex-end; gap: 8px; padding-top: 12px; border-top: 1px solid var(--bd-f1f5f9); }
.inspect-selected-stall-drawer .issd-full-main .issd-card-meta-row,
.inspect-selected-stall-drawer .issd-full-main .issd-pay-amount-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; color: var(--tx-475569); font-size: 0.9375rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-full-main .issd-pay-amount-row strong { font-size: 1.125rem; font-weight: 800; }

/* ป้ายสถานะ (soft pill ตามมาตรฐานระบบ) */
.inspect-selected-stall-drawer .issd-full-main .issd-status-pill { padding: 4px 12px; border-radius: 8px; font-size: 0.8125rem; font-weight: 600; }

/* ผู้เช่า / ผู้จอง */
.inspect-selected-stall-drawer .issd-full-main .issd-booking-applicant-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap; }
.inspect-selected-stall-drawer .issd-full-main .issd-phone-link { font-size: 0.875rem; font-weight: 600; }
.inspect-selected-stall-drawer .issd-full-main .issd-booking-footer-bar { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding-top: 12px; border-top: 1px solid var(--bd-f1f5f9); color: var(--tx-475569); font-size: 0.9375rem; }
.inspect-selected-stall-drawer .issd-full-main .issd-booking-direct-btn { display: inline-flex; align-items: center; gap: 6px; height: 34px; padding: 0 12px; border: 1px solid var(--bd-93c5fd); border-radius: 7px; background: var(--bg-ffffff); color: var(--tx-2563eb); font-size: 0.875rem; font-weight: 600; text-decoration: none; }
.inspect-selected-stall-drawer .issd-full-main .issd-booking-direct-btn:hover { background: var(--bg-eff6ff); }

/* สรุปการเงิน */
.inspect-selected-stall-drawer .issd-full-main .issd-finance-summary { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; }
.inspect-selected-stall-drawer .issd-full-main .issd-finance-summary > div { display: flex; flex-direction: column; gap: 2px; padding: 12px 14px; background: var(--bg-ffffff); border: 1px solid var(--bd-eef2f6); border-radius: 12px; }
.inspect-selected-stall-drawer .issd-full-main .issd-finance-summary span { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-full-main .issd-finance-summary strong { color: var(--tx-1e293b); font-size: 1.125rem; font-weight: 800; }

/* ข้อมูลแผง: การ์ดตัวเลข + รายการรายละเอียดเป็นแถว */
.inspect-selected-stall-drawer .issd-full-main .issd-kpi-card { background: var(--bg-ffffff); border: 1px solid var(--bd-eef2f6); border-radius: 14px; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-section { display: flex; flex-direction: column; gap: 10px; }
.inspect-selected-stall-drawer .issd-full-main .issd-section-title { color: var(--tx-1e293b); font-size: 1.125rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-grid { display: flex; flex-direction: column; background: var(--bg-ffffff); border: 1px solid var(--bd-eef2f6); border-radius: 14px; overflow: hidden; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell,
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell.full { display: flex; align-items: baseline; justify-content: space-between; gap: 16px; width: auto; padding: 14px 18px; border: 0; border-bottom: 1px solid var(--bd-f1f5f9); border-radius: 0; background: transparent; box-shadow: none; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell:last-child { border-bottom: 0; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell dt { color: var(--tx-64748b); font-size: 0.9375rem; font-weight: 500; white-space: nowrap; }
.inspect-selected-stall-drawer .issd-full-main .issd-fact-cell dd { margin: 0; color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 600; text-align: right; }
.inspect-selected-stall-drawer .issd-full-main .issd-photo-banner { border-radius: 14px; overflow: hidden; }

/* ประวัติ (timeline) */
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-panel { background: var(--bg-ffffff); border: 1px solid var(--bd-eef2f6); border-radius: 14px; padding: 8px 18px; }
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-item { padding: 12px 0; }
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-copy strong { color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-copy time { color: var(--tx-94a3b8); font-size: 0.8125rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-full-main .issd-timeline-copy p { margin: 4px 0 0; color: var(--tx-475569); font-size: 0.9375rem; line-height: 1.55; }

/* สถานะว่าง */
.inspect-selected-stall-drawer .issd-full-main .issd-empty-box { display: flex; flex-direction: column; align-items: center; gap: 10px; padding: 36px 20px; background: var(--bg-ffffff); border: 1px dashed var(--bd-dbe3ee); border-radius: 14px; color: var(--tx-94a3b8); text-align: center; }
.inspect-selected-stall-drawer .issd-full-main .issd-empty-box p { margin: 0; color: var(--tx-64748b); font-size: 0.9375rem; font-weight: 500; }

/* ตัวเลข/รหัสในรายการต้องดำเนินการ แสดงเป็นป้ายแยก ไม่ปนอยู่ในประโยค */
.issd-aac-meta { flex: 0 0 auto; padding: 2px 10px; border: 1px solid var(--bd-rgba-15-23-42-0p1); border-radius: 8px; background: var(--bg-rgba-255-255-255-0p75); color: var(--tx-475569); font-size: 0.8125rem; font-weight: 600; white-space: nowrap; }

/* ===== Drawer readability pass: larger type, roomier spacing, no redundant blocks ===== */
.inspect-selected-stall-drawer .issd-title-group strong { font-size: 1.0625rem; }
.inspect-selected-stall-drawer .issd-project-sub { font-size: 0.8125rem; }
.inspect-selected-stall-drawer .issd-header { padding: 16px 20px 14px; }
.inspect-selected-stall-drawer .issd-body { gap: 16px; padding: 18px 20px; }
.inspect-selected-stall-drawer .issd-body-full { gap: 16px; padding: 18px 22px 22px; }

.inspect-selected-stall-drawer .issd-available-line { display: flex; align-items: center; gap: 10px; color: var(--tx-15803d); }
.inspect-selected-stall-drawer .issd-available-line strong { font-size: 1.0625rem; font-weight: 700; }

.inspect-selected-stall-drawer .issd-kpi-grid { gap: 14px; }
.inspect-selected-stall-drawer .issd-kpi-grid.is-three { grid-template-columns: repeat(3, 1fr); }
.inspect-selected-stall-drawer .issd-kpi-card { padding: 16px 18px; gap: 4px; border-radius: 14px; }
.inspect-selected-stall-drawer .issd-kpi-card span,
.inspect-selected-stall-drawer .issd-kpi-card em { font-size: 0.875rem; }
.inspect-selected-stall-drawer .issd-kpi-card strong { font-size: 1.25rem; }
.inspect-selected-stall-drawer .issd-full-main .issd-kpi-grid.is-primary-stats { gap: 12px; }
.inspect-selected-stall-drawer .issd-full-main .issd-kpi-grid.is-primary-stats .issd-kpi-card {
  min-height: 104px;
  align-content: center;
  padding: 14px 18px;
  border-color: var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  box-shadow: 0 1px 2px rgba(15, 23, 42, .025);
}
.inspect-selected-stall-drawer .issd-full-main .issd-kpi-grid.is-primary-stats .issd-kpi-card span { color: var(--tx-64748b); font-size: .8125rem; font-weight: 600; }
.inspect-selected-stall-drawer .issd-full-main .issd-kpi-grid.is-primary-stats .issd-kpi-card strong { margin-top: 2px; font-size: 1.375rem; line-height: 1.25; }
.inspect-selected-stall-drawer .issd-full-main .issd-kpi-grid.is-primary-stats .issd-kpi-card em { color: var(--tx-64748b); font-size: .8125rem; }

.inspect-selected-stall-drawer .issd-stall-detail-card {
  overflow: hidden;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  background: var(--bg-ffffff);
  box-shadow: 0 1px 3px rgba(15, 23, 42, .03);
}
.inspect-selected-stall-drawer .issd-stall-detail-head {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 52px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--bd-eef2f6);
}
.inspect-selected-stall-drawer .issd-stall-detail-icon {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  flex: 0 0 34px;
  border: 1px solid var(--bd-dbeafe);
  border-radius: 9px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}
.inspect-selected-stall-drawer .issd-stall-detail-icon svg { width: 18px; height: 18px; }
.inspect-selected-stall-drawer .issd-stall-detail-head strong { color: var(--tx-1e293b); font-size: 1rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-stall-facts-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1px;
  margin: 0;
  background: var(--bg-eef2f6);
}
.inspect-selected-stall-drawer .issd-stall-fact {
  display: grid;
  align-content: start;
  gap: 5px;
  min-width: 0;
  min-height: 78px;
  padding: 13px 16px;
  background: var(--bg-ffffff);
}
.inspect-selected-stall-drawer .issd-stall-fact.is-wide,
.inspect-selected-stall-drawer .issd-stall-fact:only-child,
.inspect-selected-stall-drawer .issd-stall-fact:last-child:nth-child(odd) {
  grid-column: 1 / -1;
  min-height: 0;
}
.inspect-selected-stall-drawer .issd-stall-fact dt { color: var(--tx-64748b); font-size: .8125rem; font-weight: 600; }
.inspect-selected-stall-drawer .issd-stall-fact dd { margin: 0; color: var(--tx-1e293b); font-size: .9375rem; font-weight: 700; line-height: 1.45; overflow-wrap: anywhere; }
.inspect-selected-stall-drawer .issd-stall-fact dd small { display: block; margin-top: 2px; color: var(--tx-94a3b8); font-size: .8125rem; font-weight: 500; }
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
  .inspect-selected-stall-drawer .issd-kpi-grid.is-primary-stats .issd-kpi-card:last-child { grid-column: 1 / -1; }
  .inspect-selected-stall-drawer .issd-stall-facts-grid { grid-template-columns: minmax(0, 1fr); }
  .inspect-selected-stall-drawer .issd-stall-fact.is-wide { grid-column: auto; }
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
.inspect-selected-stall-drawer .issd-body-full .issd-tab-header-row {
  margin: 0 0 8px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--bd-e2e8f0);
}
.inspect-selected-stall-drawer .issd-body-full .issd-tab-title {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--tx-0f172a);
  font-size: 1.125rem;
  font-weight: 700;
  line-height: 1.35;
}
.inspect-selected-stall-drawer .issd-body-full .issd-tab-title::before {
  width: 4px;
  height: 22px;
  border-radius: 2px;
  background: #2563eb;
  content: '';
  flex: 0 0 auto;
}

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
  background: var(--bg-ffffff); border: 1px solid var(--bd-e2e8f0); border-radius: 14px;
  padding: 16px 18px; display: flex; flex-direction: column; gap: 12px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}
.inspect-selected-stall-drawer .issd-contract-card.is-current { border-color: var(--bd-bbf7d0); background: var(--bg-f7fef9); }

.inspect-selected-stall-drawer .issd-card-top-row,
.inspect-selected-stall-drawer .issd-booking-card-header {
  display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap;
}
.inspect-selected-stall-drawer .issd-card-title { font-size: 1rem; font-weight: 700; color: var(--tx-0f172a); }
.inspect-selected-stall-drawer .issd-contract-num { font-size: 0.9375rem; font-weight: 700; color: var(--tx-1e293b); }

/* label above value grid */
.inspect-selected-stall-drawer .issd-contract-details-grid {
  display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px 16px;
  padding-top: 12px; border-top: 1px solid var(--bd-f1f5f9);
}
.inspect-selected-stall-drawer .issd-contract-detail-item { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.inspect-selected-stall-drawer .issd-contract-detail-item.is-wide { grid-column: 1 / -1; }
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-label { font-size: 0.8125rem; color: var(--tx-64748b); font-weight: 500; }
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-value { font-size: 0.9375rem; color: var(--tx-1e293b); font-weight: 600; }
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-value.is-rent { color: var(--tx-15803d); }
.inspect-selected-stall-drawer .detail-sub { font-style: normal; font-weight: 500; color: var(--tx-2563eb); }

.inspect-selected-stall-drawer .issd-subtenant-details,
.inspect-selected-stall-drawer .issd-card-meta-row {
  display: flex; flex-wrap: wrap; align-items: center; gap: 6px 16px;
  font-size: 0.875rem; color: var(--tx-475569);
}
.inspect-selected-stall-drawer .issd-card-meta-row { justify-content: space-between; padding-top: 10px; border-top: 1px solid var(--bd-f1f5f9); }
.inspect-selected-stall-drawer .issd-card-meta-row .is-rent,
.inspect-selected-stall-drawer .issd-pay-amount-row .is-rent { color: var(--tx-15803d); font-weight: 700; }
.inspect-selected-stall-drawer .issd-pay-period { display: flex; flex-direction: column; gap: 2px; }
.inspect-selected-stall-drawer .issd-pay-period strong { font-size: 1rem; color: var(--tx-0f172a); }
.inspect-selected-stall-drawer .issd-pay-amount-row { display: flex; align-items: baseline; justify-content: space-between; gap: 10px; }
.inspect-selected-stall-drawer .issd-pay-amount-row strong { font-size: 1.125rem; }
.inspect-selected-stall-drawer .issd-maintenance-meta-grid { display: flex; flex-wrap: wrap; gap: 6px 18px; font-size: 0.875rem; color: var(--tx-64748b); }
.inspect-selected-stall-drawer .issd-maintenance-meta-grid strong { color: var(--tx-1e293b); font-weight: 600; }
.inspect-selected-stall-drawer .issd-finance-summary { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
.inspect-selected-stall-drawer .issd-finance-summary > div { padding: 12px 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 12px; background: var(--bg-ffffff); display: flex; flex-direction: column; gap: 2px; }
.inspect-selected-stall-drawer .issd-finance-summary span { font-size: 0.8125rem; color: var(--tx-64748b); }
.inspect-selected-stall-drawer .issd-finance-summary strong { font-size: 1.0625rem; color: var(--tx-0f172a); }
.inspect-selected-stall-drawer .issd-finance-summary strong.is-success { color: var(--tx-15803d); }
.inspect-selected-stall-drawer .issd-finance-summary strong.is-alert { color: var(--tx-dc2626); }

/* one shared action row + button style */
.inspect-selected-stall-drawer .issd-card-actions-row { display: flex; gap: 8px; justify-content: flex-end; padding-top: 4px; }
.inspect-selected-stall-drawer .issd-mini-action-btn {
  padding: 7px 14px; border-radius: 8px; border: 1px solid var(--bd-e2e8f0); background: var(--bg-ffffff);
  color: var(--tx-334155); font-family: inherit; font-weight: 600; cursor: pointer; transition: all .15s ease;
}
.inspect-selected-stall-drawer .issd-mini-action-btn:hover { background: var(--bg-f8fafc); }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-primary { border-color: var(--bd-93c5fd); color: var(--tx-2563eb); }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-primary:hover { background: var(--bg-eff6ff); }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-success { border-color: var(--bd-86efac); color: var(--tx-15803d); }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-success:hover { background: var(--bg-f0fdf4); }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-danger { border-color: var(--bd-fca5a5); color: var(--tx-ef4444); }
.inspect-selected-stall-drawer .issd-mini-action-btn.is-danger:hover { background: var(--bg-fef2f2); }

/* sub-filters share the same segmented look */
.inspect-selected-stall-drawer .issd-source-filter { display: inline-flex; gap: 4px; padding: 4px; background: var(--bg-f1f5f9); border-radius: 10px; align-self: flex-start; }
.inspect-selected-stall-drawer .issd-source-filter button {
  padding: 7px 14px; border: none; border-radius: 8px; background: transparent; color: var(--tx-64748b);
  font-family: inherit; font-size: 0.875rem; font-weight: 600; cursor: pointer;
}
.inspect-selected-stall-drawer .issd-source-filter button.active { background: var(--bg-ffffff); color: var(--tx-2563eb); box-shadow: 0 1px 3px rgba(15,23,42,.08); }

.inspect-selected-stall-drawer .issd-empty-box { padding: 36px 16px; border: 1px dashed var(--bd-cbd5e1); border-radius: 14px; background: var(--bg-f8fafc); display: flex; flex-direction: column; align-items: center; gap: 10px; color: var(--tx-64748b); }
.inspect-selected-stall-drawer .issd-empty-box p { margin: 0; font-size: 0.9375rem; }

/* Framed empty state is intentionally limited to the stall-detail sidebar tabs. */
.inspect-selected-stall-drawer .issd-tab-empty-state.app-empty-state {
  min-height: 0;
  margin-top: 4px;
  padding: 22px 16px;
  border: 1px dashed var(--bd-d5e0ed);
  border-radius: 16px;
  background: var(--bg-ffffff);
}

/* bookings */
.inspect-selected-stall-drawer .issd-booking-applicant-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap; }
.inspect-selected-stall-drawer .issd-booking-footer-bar { display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap; padding-top: 10px; border-top: 1px solid var(--bd-f1f5f9); font-size: 0.875rem; color: var(--tx-475569); }

/* Prevent flex children of the drawer body from collapsing into thin lines */
.inspect-selected-stall-drawer .issd-body > *,
.inspect-selected-stall-drawer .issd-body-full > * { flex-shrink: 0; }
.inspect-selected-stall-drawer .issd-mini-alert-card { flex-shrink: 0; min-height: min-content; }

/* Contract card: shared system card language with distinct header, content and action footer. */
.inspect-selected-stall-drawer .issd-contract-card {
  gap: 0;
  padding: 0;
  overflow: hidden;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  background: var(--bg-ffffff);
  box-shadow: 0 1px 3px rgba(15, 23, 42, .03);
}
.inspect-selected-stall-drawer .issd-contract-card.is-current { border-color: var(--bd-bbf7d0); background: var(--bg-ffffff); box-shadow: inset 3px 0 0 #22c55e, 0 1px 3px rgba(15, 23, 42, .03); }
.inspect-selected-stall-drawer .issd-contract-card .issd-card-top-row {
  padding: 13px 16px;
  border-bottom: 1px solid var(--bd-eef2f6);
  background: var(--bg-ffffff);
}
.inspect-selected-stall-drawer .issd-contract-top-left { gap: 10px; }
.inspect-selected-stall-drawer .issd-contract-icon {
  box-sizing: content-box;
  width: 18px;
  height: 18px;
  padding: 8px;
  flex: 0 0 18px;
  border: 1px solid var(--bd-dbeafe);
  border-radius: 9px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}
.inspect-selected-stall-drawer .issd-contract-num { font-size: 1rem; font-weight: 700; color: var(--tx-1e293b); overflow-wrap: anywhere; }
.inspect-selected-stall-drawer .issd-contract-details-grid {
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px 20px;
  padding: 14px 16px;
  border: 0;
  background: var(--bg-ffffff);
}
.inspect-selected-stall-drawer .issd-contract-detail-item { gap: 4px; }
.inspect-selected-stall-drawer .issd-contract-detail-item.is-wide { grid-column: 1 / -1; padding-top: 12px; border-top: 1px solid var(--bd-f1f5f9); }
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-label { color: var(--tx-64748b); font-size: .8125rem; font-weight: 600; }
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-value { color: var(--tx-1e293b); font-size: .9375rem; font-weight: 700; line-height: 1.45; }
.inspect-selected-stall-drawer .issd-contract-card .issd-card-actions-row {
  gap: 8px;
  justify-content: flex-end;
  margin: 0;
  padding: 11px 16px;
  border-top: 1px solid var(--bd-eef2f6);
  background: var(--bg-f8fafc);
}
.inspect-selected-stall-drawer .issd-subtenant-card,
.inspect-selected-stall-drawer .issd-payment-card,
.inspect-selected-stall-drawer .issd-maintenance-card,
.inspect-selected-stall-drawer .issd-booking-card { padding: 14px 16px; gap: 10px; }
.inspect-selected-stall-drawer .issd-contract-list,
.inspect-selected-stall-drawer .issd-subtenant-list,
.inspect-selected-stall-drawer .issd-payment-list,
.inspect-selected-stall-drawer .issd-maintenance-list,
.inspect-selected-stall-drawer .issd-booking-list { gap: 10px; }
.inspect-selected-stall-drawer .issd-body-full { gap: 12px; padding: 12px 20px 16px; }
@media (max-width: 768px) {
  .inspect-selected-stall-drawer .issd-contract-details-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}

.inspect-selected-stall-drawer .issd-sidebar-tab-count.is-action { background: var(--bg-fef3c7); color: var(--tx-b45309); }
.inspect-selected-stall-drawer .issd-sidebar-tab-count.is-action.is-danger { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.inspect-selected-stall-drawer .issd-sidebar-tab.active .issd-sidebar-tab-count.is-action { background: #d97706; color: #fff; }
.inspect-selected-stall-drawer .issd-sidebar-tab.active .issd-sidebar-tab-count.is-action.is-danger { background: #dc2626; color: #fff; }

/* Mini alert rows: chip + title on one line, divided rows (same look as full-drawer list) */
.inspect-selected-stall-drawer .issd-mini-alert-list { padding: 0; gap: 0; border-top: 1px solid var(--bd-fde68a); background: var(--bg-ffffff); }
.inspect-selected-stall-drawer .issd-mini-alert-row {
  flex-direction: row; align-items: center; justify-content: flex-start; gap: 10px;
  padding: 11px 14px; background: var(--bg-ffffff); border: none; border-radius: 0; box-shadow: none; transform: none;
  border-top: 1px solid var(--bd-f1f5f9); cursor: pointer;
}
.inspect-selected-stall-drawer .issd-mini-alert-row:first-child { border-top: none; }
.inspect-selected-stall-drawer .issd-mini-alert-row.is-danger { background: var(--bg-fff1f1); }
.inspect-selected-stall-drawer .issd-mini-alert-row:hover { background: var(--bg-fffbeb); transform: none; box-shadow: none; }
.inspect-selected-stall-drawer .issd-mini-alert-row.is-danger:hover { background: var(--bg-ffe4e4); }
.inspect-selected-stall-drawer .issd-mini-alert-title { font-size: 0.9375rem; font-weight: 600; color: var(--tx-1e293b); line-height: 1.35; min-width: 0; overflow-wrap: anywhere; }
.inspect-selected-stall-drawer .issd-mini-alert-row .issd-aac-tag { flex-shrink: 0; }

/* Mini banner header = same look as full drawer (icon box + title + count pill + chevron) */
.inspect-selected-stall-drawer .issd-mini-alert-head { gap: 10px; flex-wrap: wrap; }
.inspect-selected-stall-drawer .issd-mini-alert-head strong { font-size: 1rem; color: var(--tx-92400e); }
.inspect-selected-stall-drawer .issd-mini-alert-card .issd-aac-toggle-btn { pointer-events: none; }

/* Maintenance card: use the same header / content / footer anatomy as system cards. */
.inspect-selected-stall-drawer .issd-maintenance-card {
  gap: 0;
  padding: 0;
  overflow: hidden;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  background: var(--bg-ffffff);
  box-shadow: 0 1px 3px rgba(15, 23, 42, .03);
}
.inspect-selected-stall-drawer .issd-maintenance-card .issd-card-top-row {
  align-items: flex-start;
  padding: 14px 16px;
  border-bottom: 1px solid var(--bd-eef2f6);
  background: var(--bg-ffffff);
}
.inspect-selected-stall-drawer .issd-maintenance-heading {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  min-width: 0;
}
.inspect-selected-stall-drawer .issd-maintenance-icon {
  display: grid;
  place-items: center;
  box-sizing: border-box;
  width: 38px;
  height: 38px;
  flex: 0 0 38px;
  border: 1px solid var(--bd-dbeafe);
  border-radius: 10px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}
.inspect-selected-stall-drawer .issd-maintenance-icon svg { width: 20px; height: 20px; }
.inspect-selected-stall-drawer .issd-maintenance-heading-copy {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 7px;
  min-width: 0;
}
.inspect-selected-stall-drawer .issd-maintenance-card .issd-card-title {
  color: var(--tx-1e293b);
  font-size: 1rem;
  font-weight: 700;
  line-height: 1.4;
}
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-type-pill,
.inspect-selected-stall-drawer .issd-maintenance-card .issd-status-pill {
  padding: 4px 10px;
  border-radius: 8px;
  font-size: .8125rem;
  font-weight: 600;
  line-height: 1.35;
}
.inspect-selected-stall-drawer .issd-maintenance-card .issd-status-pill { flex: 0 0 auto; }
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 14px 20px;
  margin: 0;
  padding: 14px 16px;
  border: 0;
  border-radius: 0;
  background: var(--bg-ffffff);
}
.inspect-selected-stall-drawer .issd-maintenance-meta-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}
.inspect-selected-stall-drawer .issd-maintenance-meta-item span {
  color: var(--tx-64748b);
  font-size: .8125rem;
  font-weight: 600;
}
.inspect-selected-stall-drawer .issd-maintenance-meta-item strong {
  color: var(--tx-1e293b);
  font-size: .9375rem;
  font-weight: 700;
  line-height: 1.45;
}
.inspect-selected-stall-drawer .issd-maintenance-meta-item strong.is-rent { color: var(--tx-15803d); }
.inspect-selected-stall-drawer .issd-maintenance-card .issd-card-actions-row {
  gap: 8px;
  justify-content: flex-end;
  margin: 0;
  padding: 11px 16px;
  border-top: 1px solid var(--bd-eef2f6);
  background: var(--bg-f8fafc);
}
@media (max-width: 560px) {
  .inspect-selected-stall-drawer .issd-maintenance-card .issd-card-top-row { align-items: stretch; }
  .inspect-selected-stall-drawer .issd-maintenance-card .issd-status-pill { align-self: flex-start; margin-left: 48px; }
}

/* Booking card: compact two-row layout without nested containers. */
.inspect-selected-stall-drawer .issd-booking-card {
  gap: 0;
  padding: 0;
  overflow: hidden;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  background: var(--bg-ffffff);
  box-shadow: 0 1px 3px rgba(15, 23, 42, .03);
}
.inspect-selected-stall-drawer .issd-booking-card-header {
  padding: 12px 14px;
  border-bottom: 1px solid var(--bd-eef2f6);
  background: var(--bg-ffffff);
}
.inspect-selected-stall-drawer .issd-booking-type-group { gap: 10px; }
.inspect-selected-stall-drawer .issd-booking-ref-code {
  color: var(--tx-1e293b);
  font-size: .9375rem;
  font-weight: 700;
  overflow-wrap: anywhere;
}
.inspect-selected-stall-drawer .issd-booking-card .issd-card-type-tag,
.inspect-selected-stall-drawer .issd-booking-card .issd-status-pill {
  padding: 4px 10px;
  border-radius: 8px;
  font-size: .8125rem;
  font-weight: 600;
  line-height: 1.35;
}
.inspect-selected-stall-drawer .issd-booking-compact-body {
  display: grid;
  grid-template-columns: minmax(180px, 1fr) auto minmax(180px, auto) auto;
  align-items: center;
  gap: 10px 16px;
  padding: 12px 14px;
}
.inspect-selected-stall-drawer .issd-booking-avatar {
  width: 36px;
  height: 36px;
  border-color: var(--bd-dbeafe);
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}
.inspect-selected-stall-drawer .issd-booking-applicant-info { gap: 2px; }
.inspect-selected-stall-drawer .issd-booking-role-label {
  color: var(--tx-64748b);
  font-size: .8125rem;
  font-weight: 500;
}
.inspect-selected-stall-drawer .issd-booking-applicant-info .issd-card-title {
  color: var(--tx-1e293b);
  font-size: .9375rem;
  line-height: 1.35;
}
.inspect-selected-stall-drawer .issd-booking-compact-body .issd-phone-link,
.inspect-selected-stall-drawer .issd-booking-compact-body .issd-booking-info-item {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  min-height: 34px;
  color: var(--tx-475569);
  font-size: .875rem;
  font-weight: 600;
  white-space: nowrap;
}
.inspect-selected-stall-drawer .issd-booking-compact-body .issd-phone-link {
  padding: 0 10px;
  border: 1px solid var(--bd-dbeafe);
  border-radius: 8px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  text-decoration: none;
}
.inspect-selected-stall-drawer .issd-booking-compact-body .issd-phone-link:hover { background: var(--bg-dbeafe); }
.inspect-selected-stall-drawer .issd-booking-detail-btn { justify-self: end; }
@media (max-width: 980px) {
  .inspect-selected-stall-drawer .issd-booking-compact-body {
    grid-template-columns: minmax(180px, 1fr) auto;
  }
  .inspect-selected-stall-drawer .issd-booking-detail-btn { justify-self: end; }
}
@media (max-width: 560px) {
  .inspect-selected-stall-drawer .issd-booking-card-header { align-items: flex-start; }
  .inspect-selected-stall-drawer .issd-booking-card-header > .issd-status-pill { margin-left: 0; }
  .inspect-selected-stall-drawer .issd-booking-compact-body { grid-template-columns: 1fr; }
  .inspect-selected-stall-drawer .issd-booking-compact-body .issd-phone-link,
  .inspect-selected-stall-drawer .issd-booking-detail-btn { justify-self: stretch; }
  .inspect-selected-stall-drawer .issd-booking-detail-btn { width: 100%; }
}

/* Compact operational cards: simple header and one flexible content row. */
.inspect-selected-stall-drawer .issd-contract-card,
.inspect-selected-stall-drawer .issd-maintenance-card {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 0;
}
.inspect-selected-stall-drawer .issd-contract-card > .issd-card-top-row,
.inspect-selected-stall-drawer .issd-maintenance-card > .issd-card-top-row {
  align-items: center;
  min-height: 54px;
  padding: 10px 14px;
  border-bottom: 1px solid var(--bd-eef2f6);
}
.inspect-selected-stall-drawer .issd-contract-card .issd-contract-top-left,
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-heading { align-items: center; gap: 10px; }
.inspect-selected-stall-drawer .issd-contract-card .issd-contract-icon,
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-icon {
  box-sizing: border-box;
  width: 34px;
  height: 34px;
  flex: 0 0 34px;
  border-radius: 9px;
}
.inspect-selected-stall-drawer .issd-contract-card .issd-contract-icon { padding: 7px; }
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-icon { padding: 0; }
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-icon svg { width: 18px; height: 18px; }
.inspect-selected-stall-drawer .issd-maintenance-heading-copy {
  flex-direction: row;
  align-items: center;
  gap: 9px;
  flex-wrap: wrap;
}
.inspect-selected-stall-drawer .issd-maintenance-card .issd-card-title { order: -1; }
.inspect-selected-stall-drawer .issd-compact-card-body {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 12px;
  min-width: 0;
  padding: 10px 14px;
}
.inspect-selected-stall-drawer .issd-contract-card .issd-contract-details-grid,
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  align-items: center;
  gap: 0;
  min-width: 0;
  margin: 0;
  padding: 0;
  border: 0;
  background: transparent;
}
.inspect-selected-stall-drawer .issd-contract-detail-item,
.inspect-selected-stall-drawer .issd-maintenance-meta-item {
  flex: 1 1 130px;
  min-width: 110px;
  padding: 0 14px;
  border-left: 1px solid var(--bd-eef2f6);
}
.inspect-selected-stall-drawer .issd-contract-detail-item:first-child,
.inspect-selected-stall-drawer .issd-maintenance-meta-item:first-child {
  padding-left: 0;
  border-left: 0;
}
.inspect-selected-stall-drawer .issd-contract-detail-item:nth-child(odd),
.inspect-selected-stall-drawer .issd-maintenance-meta-item:nth-child(odd) {
  padding-left: 0;
  border-left: 0;
}
.inspect-selected-stall-drawer .issd-contract-detail-item.is-wide {
  grid-column: auto;
  flex-basis: 250px;
  padding-top: 0;
  border-top: 0;
}
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-label,
.inspect-selected-stall-drawer .issd-maintenance-meta-item span {
  font-size: .8125rem;
  line-height: 1.3;
}
.inspect-selected-stall-drawer .issd-contract-detail-item .detail-value,
.inspect-selected-stall-drawer .issd-maintenance-meta-item strong {
  font-size: .875rem;
  line-height: 1.4;
  overflow-wrap: normal;
  word-break: normal;
}
.inspect-selected-stall-drawer .issd-compact-card-body > .issd-card-actions-row {
  display: flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  margin: 0;
  padding: 0;
  border: 0;
  background: transparent;
}
@media (max-width: 1200px) {
  .inspect-selected-stall-drawer .issd-compact-card-body { align-items: stretch; flex-direction: column; }
  .inspect-selected-stall-drawer .issd-contract-card .issd-contract-details-grid,
  .inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 0;
  }
  .inspect-selected-stall-drawer .issd-contract-detail-item,
  .inspect-selected-stall-drawer .issd-maintenance-meta-item {
    min-width: 0;
    padding: 7px 14px;
    border-left: 1px solid var(--bd-eef2f6);
  }
  .inspect-selected-stall-drawer .issd-contract-detail-item:nth-child(odd),
  .inspect-selected-stall-drawer .issd-maintenance-meta-item:nth-child(odd) {
    padding-left: 0;
    border-left: 0;
  }
  .inspect-selected-stall-drawer .issd-contract-detail-item .detail-value,
  .inspect-selected-stall-drawer .issd-maintenance-meta-item strong {
    overflow-wrap: normal;
    word-break: normal;
  }
  .inspect-selected-stall-drawer .issd-compact-card-body > .issd-card-actions-row { justify-content: flex-end; padding-left: 0; }
}
@media (max-width: 560px) {
  .inspect-selected-stall-drawer .issd-contract-card > .issd-card-top-row,
  .inspect-selected-stall-drawer .issd-maintenance-card > .issd-card-top-row { align-items: flex-start; }
  .inspect-selected-stall-drawer .issd-maintenance-card .issd-status-pill { margin-left: 0; }
  .inspect-selected-stall-drawer .issd-contract-card .issd-contract-details-grid,
  .inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-grid {
    display: grid;
    grid-template-columns: 1fr;
  }
  .inspect-selected-stall-drawer .issd-contract-detail-item,
  .inspect-selected-stall-drawer .issd-maintenance-meta-item {
    flex-basis: auto;
    min-width: 0;
    padding: 8px 0 0;
    border-top: 1px solid var(--bd-eef2f6);
    border-left: 0;
  }
  .inspect-selected-stall-drawer .issd-contract-detail-item:first-child,
  .inspect-selected-stall-drawer .issd-maintenance-meta-item:first-child { padding-top: 0; border-top: 0; }
  .inspect-selected-stall-drawer .issd-contract-detail-item:nth-child(even),
  .inspect-selected-stall-drawer .issd-maintenance-meta-item:nth-child(even) { padding-left: 0; border-left: 0; }
}

/* Booking card pattern: status/ref, prominent schedule, applicant/action. */
.inspect-selected-stall-drawer .issd-booking-card {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 12px 14px;
  overflow: visible;
}
.inspect-selected-stall-drawer .issd-booking-card-header {
  align-items: center;
  padding: 0;
  border: 0;
}
.inspect-selected-stall-drawer .issd-booking-kind {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  color: var(--tx-0369a1);
}
.inspect-selected-stall-drawer .issd-booking-kind.is-application { color: var(--tx-6d28d9); }
.inspect-selected-stall-drawer .issd-booking-kind-icon {
  display: grid;
  place-items: center;
  width: 38px;
  height: 38px;
  flex: 0 0 38px;
  border-radius: 10px;
  background: var(--bg-e0f2fe);
}
.inspect-selected-stall-drawer .issd-booking-kind.is-application .issd-booking-kind-icon { background: var(--bg-f3e8ff); }
.inspect-selected-stall-drawer .issd-booking-kind-icon svg { width: 21px; height: 21px; }
.inspect-selected-stall-drawer .issd-booking-kind-copy { display: flex; flex-direction: column; gap: 1px; min-width: 0; }
.inspect-selected-stall-drawer .issd-booking-kind-copy strong {
  color: currentColor;
  font-size: 1rem;
  font-weight: 800;
  line-height: 1.4;
}
.inspect-selected-stall-drawer .issd-booking-kind-copy small {
  color: var(--tx-64748b);
  font-family: inherit;
  font-size: .8125rem;
  font-weight: 700;
  letter-spacing: .015em;
}
.inspect-selected-stall-drawer .issd-booking-status {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 999px;
}
.inspect-selected-stall-drawer .issd-booking-status-dot {
  width: 8px;
  height: 8px;
  flex: 0 0 8px;
  border-radius: 999px;
  background: currentColor;
}
.inspect-selected-stall-drawer .issd-booking-ref-code {
  color: var(--tx-64748b);
  font-family: inherit;
  font-size: .875rem;
  font-weight: 700;
  letter-spacing: .015em;
}
.inspect-selected-stall-drawer .issd-booking-schedule {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 3px 0 5px;
}
.inspect-selected-stall-drawer .issd-booking-date-line {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.inspect-selected-stall-drawer .issd-booking-date-line > strong {
  color: var(--tx-0f172a);
  font-size: 1rem;
  font-weight: 800;
  line-height: 1.3;
}
.inspect-selected-stall-drawer .issd-booking-countdown {
  display: inline-flex;
  padding: 4px 10px;
  border-radius: 8px;
  background: var(--bg-f1f5f9);
  color: var(--tx-475569);
  font-size: .8125rem;
  font-weight: 600;
}
.inspect-selected-stall-drawer .issd-booking-countdown.is-application { background: var(--bg-fff7ed); color: var(--tx-9a3412); }
.inspect-selected-stall-drawer .issd-booking-schedule p {
  margin: 0;
  color: var(--tx-475569);
  font-size: .875rem;
  font-weight: 500;
  line-height: 1.5;
}
.inspect-selected-stall-drawer .issd-booking-person-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-top: 7px;
  border-top: 1px solid var(--bd-eef2f6);
}
.inspect-selected-stall-drawer .issd-booking-person-row > .issd-phone-link,
.inspect-selected-stall-drawer .issd-booking-person-row > .issd-booking-info-item { display: none; }
.inspect-selected-stall-drawer .issd-booking-person-row .issd-booking-avatar {
  width: 38px;
  height: 38px;
  border: 0;
  border-radius: 999px;
  background: var(--bg-f1f5f9);
  color: var(--tx-334155);
  font-size: 1rem;
  font-weight: 700;
}
.inspect-selected-stall-drawer .issd-booking-person-row .issd-booking-applicant-info { gap: 1px; }
.inspect-selected-stall-drawer .issd-booking-person-row .issd-card-title {
  font-size: 1rem;
  font-weight: 800;
  line-height: 1.4;
}
.inspect-selected-stall-drawer .issd-booking-phone {
  color: var(--tx-2563eb);
  font-size: .875rem;
  font-weight: 500;
  text-decoration: none;
}
.inspect-selected-stall-drawer .issd-booking-phone:hover { text-decoration: underline; }
.inspect-selected-stall-drawer .issd-booking-person-row .issd-booking-detail-btn { margin-left: auto; }
.inspect-selected-stall-drawer .issd-application-source {
  margin-top: 10px;
  padding: 8px 11px;
  border-radius: 10px;
  background: var(--bg-f8fafc);
  color: var(--tx-475569);
  font-size: .8125rem;
  font-weight: 500;
}
@media (max-width: 560px) {
  .inspect-selected-stall-drawer .issd-booking-card-header { align-items: flex-start; flex-direction: column; }
  .inspect-selected-stall-drawer .issd-booking-person-row { align-items: flex-start; flex-wrap: wrap; }
  .inspect-selected-stall-drawer .issd-booking-person-row .issd-booking-detail-btn { width: 100%; margin-left: 0; }
}

/* Contract card pattern: status, tenant, timeline, finance summary and actions. */
.inspect-selected-stall-drawer .issd-contract-card,
.inspect-selected-stall-drawer .issd-contract-card.is-current {
  gap: 0;
  padding: 14px 16px;
  border-color: var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  box-shadow: 0 1px 3px rgba(15, 23, 42, .03);
}
.inspect-selected-stall-drawer .issd-contract-card > .issd-card-top-row {
  min-height: 0;
  padding: 0;
  border: 0;
}
.inspect-selected-stall-drawer .issd-contract-status,
.inspect-selected-stall-drawer .issd-maintenance-status {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 5px 11px;
  border-radius: 999px;
}
.inspect-selected-stall-drawer .issd-contract-status-dot,
.inspect-selected-stall-drawer .issd-maintenance-status-dot {
  width: 7px;
  height: 7px;
  flex: 0 0 7px;
  border-radius: 999px;
  background: currentColor;
}
.inspect-selected-stall-drawer .issd-contract-card .issd-contract-num {
  color: var(--tx-64748b);
  font-family: inherit;
  font-size: .8125rem;
  font-weight: 700;
  letter-spacing: .015em;
}
.inspect-selected-stall-drawer .issd-contract-card > .issd-compact-card-body {
  gap: 12px;
  padding: 12px 0 0;
}
.inspect-selected-stall-drawer .issd-contract-overview { display: flex; flex-direction: column; gap: 10px; }
.inspect-selected-stall-drawer .issd-contract-identity { display: flex; flex-direction: column; gap: 3px; }
.inspect-selected-stall-drawer .issd-contract-identity strong { color: var(--tx-0f172a); font-size: 1rem; font-weight: 800; }
.inspect-selected-stall-drawer .issd-contract-identity span,
.inspect-selected-stall-drawer .issd-contract-period-row { color: var(--tx-475569); font-size: .875rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-contract-period-row { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.inspect-selected-stall-drawer .issd-contract-period-row strong { color: var(--tx-334155); font-size: .8125rem; white-space: nowrap; }
.inspect-selected-stall-drawer .issd-contract-progress { height: 6px; overflow: hidden; border-radius: 999px; background: var(--bg-e2e8f0); }
.inspect-selected-stall-drawer .issd-contract-progress span { display: block; height: 100%; border-radius: inherit; background: #16a34a; }
.inspect-selected-stall-drawer .issd-contract-money-summary {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  padding: 10px 12px;
  border-radius: 12px;
  background: var(--bg-f8fafc);
}
.inspect-selected-stall-drawer .issd-contract-money-summary > div { display: flex; flex-direction: column; gap: 3px; }
.inspect-selected-stall-drawer .issd-contract-money-summary span { color: var(--tx-64748b); font-size: .8125rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-contract-money-summary strong { color: var(--tx-0f172a); font-size: .9375rem; font-weight: 800; }
.inspect-selected-stall-drawer .issd-contract-money-summary small { color: var(--tx-64748b); font-size: .8125rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-contract-money-summary small.is-paid { color: var(--tx-15803d); }
.inspect-selected-stall-drawer .issd-contract-card .issd-card-actions-row { justify-content: flex-end; padding: 0; border: 0; background: transparent; }

/* Maintenance card pattern: status/type, issue context, summary and actions. */
.inspect-selected-stall-drawer .issd-maintenance-card {
  gap: 0;
  padding: 8px 12px;
  border-color: var(--bd-e2e8f0);
  background: var(--bg-ffffff);
}
.inspect-selected-stall-drawer .issd-maintenance-card > .issd-card-top-row {
  min-height: 0;
  margin: 0;
  padding: 0;
  border: 0;
  flex: 0 0 auto;
}
.inspect-selected-stall-drawer .issd-maintenance-type-label {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  color: var(--tx-475569);
  font-size: .8125rem;
  font-weight: 600;
}
.inspect-selected-stall-drawer .issd-maintenance-type-label > span { display: grid; place-items: center; }
.inspect-selected-stall-drawer .issd-maintenance-type-label svg { width: 16px; height: 16px; }
.inspect-selected-stall-drawer .issd-maintenance-identity {
  display: flex;
  flex: 0 0 auto;
  flex-direction: column;
  min-height: 0;
  gap: 0;
  margin: 0;
  padding: 0;
}
.inspect-selected-stall-drawer .issd-maintenance-identity strong { color: var(--tx-0f172a); font-size: 1rem; font-weight: 800; line-height: 1.35; }
.inspect-selected-stall-drawer .issd-maintenance-identity span { color: var(--tx-475569); font-size: .875rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-maintenance-card > .issd-compact-card-body {
  flex: 0 0 auto;
  min-height: 0;
  gap: 2px;
  margin: 0;
  padding: 0;
}
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin: 2px 0;
  padding: 10px 14px;
  border-radius: 12px;
  background: var(--bg-f8fafc);
}
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-item {
  min-width: 0;
  padding: 0;
  border: 0;
}
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-item span { color: var(--tx-64748b); font-size: .8125rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-item strong { color: var(--tx-1e293b); font-size: .875rem; font-weight: 700; }
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-item strong.is-waiting { color: var(--tx-c2410c); }
.inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-item strong.is-success { color: var(--tx-15803d); }
.inspect-selected-stall-drawer .issd-maintenance-footer-row { display: flex; align-items: center; justify-content: flex-end; gap: 8px; margin-top: 12px; }
.inspect-selected-stall-drawer .issd-maintenance-footer-row .issd-card-actions-row { padding: 0; border: 0; background: transparent; }
.inspect-selected-stall-drawer .issd-maintenance-assignment-note {
  display: flex;
  align-items: center;
  gap: 6px 14px;
  flex-wrap: wrap;
  padding: 9px 11px;
  border-radius: 9px;
  background: var(--bg-f8fafc);
  color: var(--tx-64748b);
  font-size: .8125rem;
  font-weight: 500;
}
.inspect-selected-stall-drawer .issd-maintenance-assignment-note strong { color: var(--tx-334155); font-weight: 700; }
.inspect-selected-stall-drawer .issd-maintenance-assignment-note p { width: 100%; margin: 0; color: var(--tx-475569); line-height: 1.5; }
.inspect-selected-stall-drawer .issd-maintenance-result { display: flex; flex-direction: column; gap: 8px; }
.inspect-selected-stall-drawer .issd-maintenance-photo-compare { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 10px; }
.inspect-selected-stall-drawer .issd-maintenance-photo-compare figure { min-width: 0; margin: 0; }
.inspect-selected-stall-drawer .issd-maintenance-photo-frame {
  display: grid;
  place-items: center;
  height: 112px;
  overflow: hidden;
  border: 1px dashed var(--bd-cbd5e1);
  border-radius: 10px;
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
  font-size: .8125rem;
}
.inspect-selected-stall-drawer .issd-maintenance-photo-frame.is-after { border-color: var(--bd-86efac); background: var(--bg-f0fdf4); color: var(--tx-15803d); }
.inspect-selected-stall-drawer .issd-maintenance-photo-frame img { width: 100%; height: 100%; object-fit: cover; }
.inspect-selected-stall-drawer .issd-maintenance-photo-compare figcaption { margin-top: 5px; color: var(--tx-64748b); font-size: .75rem; font-weight: 500; }
.inspect-selected-stall-drawer .issd-maintenance-result p,
.inspect-selected-stall-drawer .issd-maintenance-rejection-summary {
  margin: 0;
  padding: 9px 11px;
  border-radius: 9px;
  background: var(--bg-f8fafc);
  color: var(--tx-475569);
  font-size: .8125rem;
  font-weight: 500;
  line-height: 1.5;
}
.inspect-selected-stall-drawer .issd-maintenance-rejection-summary { display: flex; gap: 5px; flex-wrap: wrap; }
.inspect-selected-stall-drawer .issd-maintenance-rejection-summary strong { color: var(--tx-334155); }
@media (max-width: 560px) {
  .inspect-selected-stall-drawer .issd-contract-period-row,
  .inspect-selected-stall-drawer .issd-maintenance-footer-row { align-items: flex-start; flex-direction: column; }
  .inspect-selected-stall-drawer .issd-contract-money-summary,
  .inspect-selected-stall-drawer .issd-maintenance-card .issd-maintenance-meta-grid,
  .inspect-selected-stall-drawer .issd-maintenance-photo-compare { grid-template-columns: 1fr; }
  .inspect-selected-stall-drawer .issd-maintenance-footer-row .issd-card-actions-row { width: 100%; justify-content: flex-end; }
}

/* Photo editor in GIS edit drawer */
.inspect-selected-stall-drawer .issd-photo-editor { display: flex; flex-direction: column; gap: 8px; }
.inspect-selected-stall-drawer .issd-photo-editor-head { display: flex; align-items: center; justify-content: space-between; font-size: 0.875rem; font-weight: 600; color: var(--tx-1e293b); }
.inspect-selected-stall-drawer .issd-photo-editor-head small { font-size: 0.8125rem; color: var(--tx-64748b); font-weight: 500; }
.inspect-selected-stall-drawer .issd-photo-editor-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.inspect-selected-stall-drawer .issd-photo-thumb { position: relative; aspect-ratio: 4 / 3; border-radius: 10px; overflow: hidden; border: 1px solid var(--bd-e2e8f0); background: var(--bg-f1f5f9); }
.inspect-selected-stall-drawer .issd-photo-thumb img { width: 100%; height: 100%; object-fit: cover; display: block; }
.inspect-selected-stall-drawer .issd-photo-remove { position: absolute; top: 5px; right: 5px; width: 24px; height: 24px; border-radius: 999px; border: none; background: rgba(15,23,42,.72); color: #fff; display: grid; place-items: center; cursor: pointer; }
.inspect-selected-stall-drawer .issd-photo-remove:hover { background: #dc2626; }
.inspect-selected-stall-drawer .issd-photo-add { aspect-ratio: 4 / 3; border: 1.5px dashed var(--bd-93c5fd); border-radius: 10px; background: var(--bg-f8fbff); color: var(--tx-2563eb); display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 4px; font-size: 0.8125rem; font-weight: 600; cursor: pointer; transition: background .15s ease; }
.inspect-selected-stall-drawer .issd-photo-add:hover { background: var(--bg-eff6ff); }
</style>
