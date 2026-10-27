<template>
  <div>
    <AppHeader />
    <main class="content-body manage-page" :class="{ 'is-map-page': isMapOverviewMode }">
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
      <section v-if="isInspectDashboardMode" class="officer-summary-dashboard" aria-label="สรุปข้อมูลแดชบอร์ดตลาด">
        <!-- 1. Operational KPI Strip (แยกสถิติ "การจองจริง" ออกจาก "รายการรอคิว" ชัดเจน) -->
        <div class="map-summary-grid dash-kpi-ribbon">
          <div class="map-summary-card summary-reserved">
            <span class="map-summary-icon" aria-hidden="true">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"></circle><path d="M12 7v5l3 2"></path></svg>
            </span>
            <span class="map-summary-content">
              <small class="map-summary-title">การจองจริง (รออนุมัติ {{ realActionCounts.pendingBookings }})</small>
              <span class="map-summary-value-row">
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
              <span class="map-summary-value-row">
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
              <span class="map-summary-value-row">
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
              <span class="map-summary-value-row">
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
      </section>

      <section v-if="!isInspectDashboardMode" class="workspace-grid" :class="{ 'is-fullscreen-map-view': isMapOverviewMode }">
        <article class="panel gis-panel" :class="{ 'is-fullscreen-map-panel': isMapOverviewMode }">
          <!-- 1. Heading และ Toolbar ปกติสำหรับโหมดจัดการผัง GIS -->
          <div v-if="!isMapOverviewMode" class="panel-heading gis-panel-heading">
            <div class="gis-panel-header-left">
              <h2>{{ marketName || (isInspectMode ? 'ผังแผงเช่า' : 'ผังแผงค้า') }}</h2>
              <p>ทั้งหมด {{ totalMarketStalls.toLocaleString('th-TH') }} แผง</p>
            </div>
            <div class="gis-heading-right">
              <div id="step1-map-toolbar" class="gis-map-toolbar"></div>
            </div>
          </div>

          <!-- ภาพรวมตลาด: KPI ลอยเหนือผัง และเครื่องมือแบบแนวตั้งด้านขวา -->
          <div v-if="isMapOverviewMode" class="market-map-kpi-strip" aria-label="สรุปสถานะแผงเช่า">
            <button type="button" class="market-map-kpi is-total" :class="{ active: dashboardStatusFilter === 'all' }" @click="setDashboardStatusFilter('all')">
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="4" y="4" width="6" height="6"/><rect x="14" y="4" width="6" height="6"/><rect x="4" y="14" width="6" height="6"/><rect x="14" y="14" width="6" height="6"/></svg></span>
              <span class="market-map-kpi-copy"><small>ทั้งหมด</small><strong>{{ currentDisplayedSummary.total.toLocaleString('th-TH') }} <em>แผง</em></strong></span>
            </button>
            <button type="button" class="market-map-kpi is-available" :class="{ active: dashboardStatusFilter === 'available' }" @click="setDashboardStatusFilter('available')">
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="9"/><path d="m8 12 2.5 2.5L16 9"/></svg></span>
              <span class="market-map-kpi-copy"><small>ว่าง</small><strong>{{ currentDisplayedSummary.available.toLocaleString('th-TH') }} <em>แผง</em></strong></span>
            </button>
            <button type="button" class="market-map-kpi is-reserved" :class="{ active: dashboardStatusFilter === 'reserved' }" @click="setDashboardStatusFilter('reserved')">
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg></span>
              <span class="market-map-kpi-copy"><small>ติดจอง</small><strong>{{ currentDisplayedSummary.reserved.toLocaleString('th-TH') }} <em>แผง</em></strong></span>
            </button>
            <button type="button" class="market-map-kpi is-occupied" :class="{ active: dashboardStatusFilter === 'occupied' }" @click="setDashboardStatusFilter('occupied')">
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="8" r="4"/><path d="M5 21a7 7 0 0 1 14 0"/></svg></span>
              <span class="market-map-kpi-copy"><small>มีผู้เช่า</small><strong>{{ currentDisplayedSummary.occupied.toLocaleString('th-TH') }} <em>แผง</em></strong></span>
            </button>
            <div class="market-map-kpi is-overdue">
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3 2.5 20h19L12 3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg></span>
              <span class="market-map-kpi-copy"><small>ค้างชำระ</small><strong>{{ marketMapRiskSummary.overdue.toLocaleString('th-TH') }} <em>แผง</em></strong></span>
            </div>
            <div class="market-map-kpi is-expiring">
              <span class="market-map-kpi-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="5" width="18" height="16" rx="2"/><path d="M16 3v4M8 3v4M3 10h18"/></svg></span>
              <span class="market-map-kpi-copy"><small>ใกล้หมดสัญญา</small><strong>{{ marketMapRiskSummary.expiring.toLocaleString('th-TH') }} <em>แผง</em></strong></span>
            </div>
          </div>

          <div v-if="isMapOverviewMode" class="market-map-tool-rail" aria-label="เครื่องมือแผนที่">
            <div id="step1-map-toolbar" class="gis-map-toolbar"></div>
            <button
              type="button"
              class="market-map-view-button"
              :class="{ active: dashboardStatusFilter === 'all' }"
              title="แสดงแผงทั้งหมด"
              aria-label="แสดงแผงทั้งหมด"
              @click="setDashboardStatusFilter('all')"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
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
            :stallStatusFilter="dashboardStatusFilter"
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
            <template #overlay>
              <transition name="slide-fade">
                <div
                  v-if="showInspectStallOverlay && stallEdit.ident"
                  id="map-selected-stall-drawer"
                  class="inspect-selected-stall-drawer"
                  :class="{
                    'is-management-mode': !isInspectMode,
                    'is-edit-mode': stallDrawerMode === 'edit',
                    'is-full-detail-mode': stallDrawerMode === 'full-detail'
                  }"
                >
                  <!-- 1. Header -->
                  <div class="issd-header">
                    <div class="issd-header-top">
                      <div class="issd-title-group">
                        <button
                          v-if="(stallDrawerMode === 'edit' && !isCreatingNewStall) || stallDrawerMode === 'full-detail'"
                          type="button"
                          class="issd-back-btn"
                          :title="stallDrawerMode === 'full-detail' ? 'กลับไปแบบสรุปย่อ' : 'กลับไปดูรายละเอียด'"
                          aria-label="ย้อนกลับ"
                          @click="stallDrawerMode === 'full-detail' ? backToMiniDetail() : cancelStallEdit()"
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

                  <!-- 1.1 In-Map Detail Tabs Bar -->
                  <div v-if="stallDrawerMode === 'full-detail'" class="issd-detail-tabs-bar" role="tablist">
                    <button
                      type="button"
                      role="tab"
                      class="issd-tab-btn"
                      :class="{ active: stallDetailTab === 'info' }"
                      @click="stallDetailTab = 'info'"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18M9 21V9"/>
                      </svg>
                      <span>ข้อมูลแผง</span>
                    </button>
                    <button
                      type="button"
                      role="tab"
                      class="issd-tab-btn"
                      :class="{ active: stallDetailTab === 'contracts' }"
                      @click="stallDetailTab = 'contracts'"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/>
                      </svg>
                      <span>สัญญา</span>
                      <span v-if="stallContractsList.length" class="issd-tab-count">{{ stallContractsList.length }}</span>
                    </button>
                    <button
                      type="button"
                      role="tab"
                      class="issd-tab-btn"
                      :class="{ active: stallDetailTab === 'subtenants' }"
                      @click="stallDetailTab = 'subtenants'"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                      </svg>
                      <span>ผู้เช่าช่วง</span>
                      <span v-if="stallSubtenants.length" class="issd-tab-count">{{ stallSubtenants.length }}</span>
                    </button>
                    <button
                      type="button"
                      role="tab"
                      class="issd-tab-btn"
                      :class="{ active: stallDetailTab === 'payments' }"
                      @click="stallDetailTab = 'payments'"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <rect x="2" y="5" width="20" height="14" rx="2"/><line x1="2" y1="10" x2="22" y2="10"/>
                      </svg>
                      <span>ชำระเงิน</span>
                      <span v-if="stallPayments.length" class="issd-tab-count">{{ stallPayments.length }}</span>
                    </button>
                    <button
                      type="button"
                      role="tab"
                      class="issd-tab-btn"
                      :class="{ active: stallDetailTab === 'maintenance' }"
                      @click="stallDetailTab = 'maintenance'"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
                      </svg>
                      <span>ซ่อม/ขออนุญาต</span>
                      <span v-if="stallMaintenance.length" class="issd-tab-count">{{ stallMaintenance.length }}</span>
                    </button>
                    <button
                      type="button"
                      role="tab"
                      class="issd-tab-btn"
                      :class="{ active: stallDetailTab === 'bookings' }"
                      @click="stallDetailTab = 'bookings'"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
                      </svg>
                      <span>การจอง</span>
                      <span v-if="stallBookingsAndApps.length" class="issd-tab-count">{{ stallBookingsAndApps.length }}</span>
                    </button>
                  </div>

                  <!-- 2. Body -->
                  <!-- 2A. Mini Summary Mode -->
                  <div v-if="stallDrawerMode === 'detail'" class="issd-body">
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

                    <!-- Current Occupancy / Booking Summary -->
                    <div class="issd-occupancy-box" :class="`is-${stallEdit.stall_status || 'available'}`">
                      <div class="issd-occupancy-head">
                        <span>{{ inspectStallOccupancy.heading }}</span>
                        <small v-if="inspectStallOccupancy.refCode">{{ inspectStallOccupancy.refCode }}</small>
                      </div>
                      <strong>{{ inspectStallOccupancy.primaryText }}</strong>
                      <p v-if="inspectStallOccupancy.subText">{{ inspectStallOccupancy.subText }}</p>
                    </div>
                  </div>

                  <!-- 2B. Full Detail Mode (In-Map Detail Overlay) -->
                  <div v-else-if="stallDrawerMode === 'full-detail'" class="issd-body issd-body-full">
                    <!-- Tab: ข้อมูลแผงเช่า (info) -->
                    <template v-if="stallDetailTab === 'info'">
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

                      <!-- KPI Strip 4 ช่อง -->
                      <div class="issd-kpi-grid">
                        <div class="issd-kpi-card">
                          <span>ขนาดพื้นที่</span>
                          <strong>{{ inspectStallSizeSqm }} ตร.ม.</strong>
                          <em>กว้าง {{ stallEdit.width || 2 }} × ลึก {{ stallEdit.depth || 2.5 }} ม.</em>
                        </div>
                        <div class="issd-kpi-card is-green">
                          <span>ค่าเช่ารายเดือน</span>
                          <strong>฿{{ Number(stallEdit.monthly_rent || 1800).toLocaleString('th-TH') }}</strong>
                          <em>ต่อเดือน</em>
                        </div>
                        <div class="issd-kpi-card">
                          <span>ค่าเช่ารายวัน</span>
                          <strong>฿{{ Number(stallEdit.daily_rent || 150).toLocaleString('th-TH') }}</strong>
                          <em>ต่อวัน</em>
                        </div>
                        <div class="issd-kpi-card">
                          <span>ค่าธรรมเนียมแรกเข้า</span>
                          <strong>฿{{ Number(inspectStallEntryFee).toLocaleString('th-TH') }}</strong>
                          <em>ชำระเมื่อทำสัญญา</em>
                        </div>
                      </div>

                      <!-- Occupancy Card -->
                      <div class="issd-occupancy-box" :class="`is-${stallEdit.stall_status || 'available'}`">
                        <div class="issd-occupancy-head">
                          <span>{{ inspectStallOccupancy.heading }}</span>
                          <small v-if="inspectStallOccupancy.refCode">{{ inspectStallOccupancy.refCode }}</small>
                        </div>
                        <strong>{{ inspectStallOccupancy.primaryText }}</strong>
                        <p v-if="inspectStallOccupancy.subText">{{ inspectStallOccupancy.subText }}</p>
                      </div>

                      <!-- Detailed Specifications Grid -->
                      <div class="issd-fact-section">
                        <span class="issd-section-title">ข้อมูลจำเพาะแผง</span>
                        <div class="issd-fact-grid">
                          <div class="issd-fact-cell">
                            <dt>รหัสแผง</dt>
                            <dd>{{ stallEdit.ident }}</dd>
                          </div>
                          <div class="issd-fact-cell">
                            <dt>ชื่อแผง</dt>
                            <dd>{{ stallEdit.name || '-' }}</dd>
                          </div>
                          <div class="issd-fact-cell">
                            <dt>โครงการ / โซน</dt>
                            <dd>{{ cleanProjectName(selectedProject?.name || stallEdit.project_name || '-') }}</dd>
                          </div>
                          <div class="issd-fact-cell">
                            <dt>ชั้น</dt>
                            <dd>{{ formatFloorDisplay(stallEdit.floor) }}</dd>
                          </div>
                          <div class="issd-fact-cell full">
                            <dt>หมวดสินค้าที่อนุญาต</dt>
                            <dd>{{ stallEdit.category || selectedProject?.category || 'แผงค้ามาตรฐาน' }}</dd>
                          </div>
                          <div v-if="stallEdit.description" class="issd-fact-cell full">
                            <dt>รายละเอียดเพิ่มเติม</dt>
                            <dd>{{ stallEdit.description }}</dd>
                          </div>
                        </div>
                      </div>
                    </template>

                    <!-- Tab: การจองและคำร้อง (bookings) -->
                    <template v-else-if="stallDetailTab === 'bookings'">
                      <div v-if="stallBookingsAndApps.length > 0" class="issd-booking-list">
                        <div
                          v-for="item in stallBookingsAndApps"
                          :key="item.key"
                          class="issd-booking-card"
                        >
                          <div class="issd-card-top-row">
                            <span class="issd-card-type-tag">{{ item.typeLabel }}</span>
                            <span class="issd-status-pill" :class="item.badgeClass">{{ item.statusText }}</span>
                          </div>
                          <strong class="issd-card-title">{{ item.applicantName }}</strong>
                          <div class="issd-card-meta-row">
                            <span>รหัส: <b>{{ item.refNumber }}</b></span>
                            <span>{{ item.dateText }}</span>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <circle cx="12" cy="12" r="10"></circle>
                          <line x1="12" y1="8" x2="12" y2="12"></line>
                          <line x1="12" y1="16" x2="12.01" y2="16"></line>
                        </svg>
                        <p>ยังไม่มีรายการจองหรือคำร้องสำหรับแผงนี้</p>
                      </div>
                    </template>

                    <!-- Tab: สัญญา (contracts) -->
                    <template v-else-if="stallDetailTab === 'contracts'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">สัญญาเช่า</strong>
                          <span class="issd-tab-subtitle">จัดการสัญญาเช่าของแผง {{ stallEdit.ident }}</span>
                        </div>
                      </div>

                      <div v-if="stallContractsList.length > 0" class="issd-contract-list">
                        <div
                          v-for="contract in stallContractsList"
                          :key="contract.id || contract.contract_number"
                          class="issd-contract-card"
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
                              {{ contract.status === 'active' ? 'ใช้งานอยู่' : (contract.status === 'expired' ? 'หมดอายุ' : ((contract.status === 'cancelled' || contract.status === 'terminated') ? 'ยกเลิกแล้ว' : (contract.status || 'ร่างสัญญา'))) }}
                            </span>
                          </div>

                          <div class="issd-contract-details-grid">
                            <div class="issd-contract-detail-item">
                              <span class="detail-label">ผู้เช่า:</span>
                              <strong class="detail-value">{{ contract.tenant_name || contract.name || '-' }}</strong>
                            </div>
                            <div v-if="contract.tenant_phone || contract.phone" class="issd-contract-detail-item">
                              <span class="detail-label">เบอร์โทร:</span>
                              <span class="detail-value">{{ contract.tenant_phone || contract.phone }}</span>
                            </div>
                            <div v-if="contract.shop_name" class="issd-contract-detail-item">
                              <span class="detail-label">ร้านค้า:</span>
                              <span class="detail-value">{{ contract.shop_name }}</span>
                            </div>
                            <div class="issd-contract-detail-item">
                              <span class="detail-label">ค่าเช่า:</span>
                              <strong class="detail-value is-rent">฿{{ Number(contract.monthly_rental || contract.monthly_rent || 0).toLocaleString('th-TH') }}/ด.</strong>
                            </div>
                            <div v-if="contract.security_deposit || contract.deposit_amount" class="issd-contract-detail-item">
                              <span class="detail-label">เงินประกัน:</span>
                              <span class="detail-value">฿{{ Number(contract.security_deposit || contract.deposit_amount).toLocaleString('th-TH') }}</span>
                            </div>
                            <div class="issd-contract-detail-item full-span">
                              <span class="detail-label">ระยะสัญญา:</span>
                              <span class="detail-value">{{ contract.start_date || '-' }} ถึง {{ contract.end_date || '-' }}</span>
                            </div>
                          </div>

                          <div v-if="contract.note || contract.terms_and_conditions" class="issd-subtenant-note" style="margin-top: 6px;">
                            {{ contract.note || contract.terms_and_conditions }}
                          </div>

                          <div class="issd-card-actions-row">
                            <button
                              type="button"
                              class="issd-mini-action-btn is-primary"
                              @click="openEditContractModal(contract)"
                            >
                              แก้ไขสัญญา
                            </button>
                            <button
                              v-if="contract.status === 'active'"
                              type="button"
                              class="issd-mini-action-btn is-danger"
                              @click="terminateContract(contract)"
                            >
                              ยกเลิกสัญญา
                            </button>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                          <polyline points="14 2 14 8 20 8"></polyline>
                        </svg>
                        <p>ยังไม่มีประวัติสัญญาสำหรับแผงนี้</p>
                        <span class="issd-empty-sub">กดปุ่ม “ออกสัญญาเช่าใหม่สำหรับแผงนี้” ด้านล่างเพื่อเริ่มทำสัญญาและกำหนดผู้เช่า</span>
                      </div>
                    </template>

                    <!-- Tab: ผู้เช่าช่วง (subtenants) -->
                    <template v-else-if="stallDetailTab === 'subtenants'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">ผู้เช่าช่วง</strong>
                          <span class="issd-tab-subtitle">จัดการรายชื่อผู้เช่าช่วงของแผง {{ stallEdit.ident }}</span>
                        </div>
                      </div>

                      <div v-if="stallSubtenants.length > 0" class="issd-subtenant-list">
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
                                <span v-if="sub.idCard" class="issd-subtenant-idcard">เลข ปชช. {{ sub.idCard }}</span>
                              </div>
                            </div>
                            <span class="issd-status-pill" :class="sub.status === 'active' ? 'is-available' : 'is-occupied'">
                              {{ sub.status === 'active' ? 'กำลังเช่าช่วง' : 'สิ้นสุดสัญญา' }}
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
                            <span>สัญญา: {{ sub.startDate || '-' }} ถึง {{ sub.endDate || '-' }}</span>
                            <span v-if="sub.subRentAmount" class="is-rent">฿{{ Number(sub.subRentAmount).toLocaleString('th-TH') }}/ด.</span>
                          </div>
                          <div v-if="sub.note" class="issd-subtenant-note">
                            {{ sub.note }}
                          </div>
                          <div class="issd-card-actions-row">
                            <button type="button" class="issd-mini-action-btn" @click="toggleSubtenantStatus(sub.id)">
                              {{ sub.status === 'active' ? 'สิ้นสุดสัญญา' : 'กลับมาเช่าต่อ' }}
                            </button>
                            <button type="button" class="issd-mini-action-btn is-danger" @click="deleteSubtenant(sub.id)">
                              ลบ
                            </button>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/>
                        </svg>
                        <p>ยังไม่มีข้อมูลผู้เช่าช่วงสำหรับแผงนี้</p>
                        <AppButton size="sm" variant="outline" @click="openAddSubtenantModal">
                          + เพิ่มผู้เช่าช่วง
                        </AppButton>
                      </div>
                    </template>

                    <!-- Tab: ประวัติชำระเงิน (payments) -->
                    <template v-else-if="stallDetailTab === 'payments'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">ประวัติการชำระเงิน</strong>
                          <span class="issd-tab-subtitle">รอบบิลและประวัติการชำระเงินของแผง {{ stallEdit.ident }}</span>
                        </div>
                      </div>

                      <div v-if="stallPayments.length > 0" class="issd-payment-list">
                        <div
                          v-for="pay in stallPayments"
                          :key="pay.id"
                          class="issd-payment-card"
                        >
                          <div class="issd-card-top-row">
                            <div class="issd-pay-period">
                              <strong>{{ pay.periodMonth }}</strong>
                              <span>{{ pay.title }}</span>
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
                            <span>วันที่ชำระ: {{ pay.paidDate || '-' }}</span>
                            <button type="button" class="issd-mini-action-btn is-danger" @click="deletePayment(pay.id)">ลบ</button>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <rect x="2" y="5" width="20" height="14" rx="2"/><line x1="2" y1="10" x2="22" y2="10"/>
                        </svg>
                        <p>ยังไม่มีประวัติการชำระเงินสำหรับแผงนี้</p>
                        <AppButton size="sm" variant="outline" @click="openAddPaymentModal">
                          + บันทึกการชำระเงิน
                        </AppButton>
                      </div>
                    </template>

                    <!-- Tab: ประวัติการซ่อม / ขออนุญาต (maintenance) -->
                    <template v-else-if="stallDetailTab === 'maintenance'">
                      <div class="issd-tab-header-row">
                        <div>
                          <strong class="issd-tab-title">ประวัติการซ่อม / ขออนุญาต</strong>
                          <span class="issd-tab-subtitle">การแจ้งซ่อมแซมและขออนุญาตของแผง {{ stallEdit.ident }}</span>
                        </div>
                      </div>

                      <div v-if="stallMaintenance.length > 0" class="issd-maintenance-list">
                        <div
                          v-for="m in stallMaintenance"
                          :key="m.id"
                          class="issd-maintenance-card"
                        >
                          <div class="issd-card-top-row">
                            <span class="issd-maintenance-type-pill" :class="`type-${m.type}`">
                              {{ formatMaintenanceType(m.type) }}
                            </span>
                            <span class="issd-status-pill" :class="m.status === 'completed' ? 'is-available' : (m.status === 'rejected' ? 'is-occupied' : 'is-reserved')">
                              {{ formatMaintenanceStatus(m.status) }}
                            </span>
                          </div>
                          <strong class="issd-card-title">{{ m.title }}</strong>
                          <p class="issd-maintenance-desc">{{ m.description }}</p>
                          <div class="issd-maintenance-meta-grid">
                            <div><span>ผู้แจ้ง:</span> <strong>{{ m.requesterName }}</strong></div>
                            <div><span>วันที่แจ้ง:</span> <strong>{{ m.requestDate }}</strong></div>
                            <div v-if="m.technicianName"><span>ผู้รับผิดชอบ:</span> <strong>{{ m.technicianName }}</strong></div>
                            <div v-if="m.estimatedCost"><span>ค่าใช้จ่าย:</span> <strong class="is-rent">฿{{ Number(m.estimatedCost).toLocaleString('th-TH') }}</strong></div>
                          </div>
                          <div class="issd-card-actions-row">
                            <button
                              v-if="m.status !== 'completed'"
                              type="button"
                              class="issd-mini-action-btn is-success"
                              @click="updateMaintenanceStatus(m.id, 'completed')"
                            >
                              ทำเครื่องหมายว่าเสร็จสิ้น
                            </button>
                            <button type="button" class="issd-mini-action-btn is-danger" @click="deleteMaintenance(m.id)">
                              ลบ
                            </button>
                          </div>
                        </div>
                      </div>
                      <div v-else class="issd-empty-box">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
                        </svg>
                        <p>ยังไม่มีประวัติการซ่อมหรือขออนุญาตสำหรับแผงนี้</p>

                        <AppButton size="sm" variant="outline" @click="openAddMaintenanceModal">
                          + แจ้งซ่อม / ขออนุญาต
                        </AppButton>
                      </div>
                    </template>
                  </div>

                  <!-- View: Edit in Map -->
                  <div v-else class="issd-body issd-body-edit">
                    <form class="issd-form-grid" novalidate @submit.prevent="saveStall">
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

                  <!-- 3. Footer -->
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
                  <div v-else-if="stallDrawerMode === 'full-detail'" class="issd-footer issd-footer-full">
                    <!-- Tab: ข้อมูลแผง (info) -->
                    <template v-if="stallDetailTab === 'info'">
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
                        ออกสัญญาเช่าใหม่สำหรับแผงนี้
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
                        แจ้งซ่อม / ขออนุญาต
                      </AppButton>
                    </template>

                    <!-- Tab: การจองและคำร้อง (bookings) -->
                    <template v-else-if="stallDetailTab === 'bookings'">
                      <AppButton
                        variant="secondary"
                        fullWidth
                        to="/bookings"
                      >
                        <template #prefix>
                          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                            <line x1="16" y1="2" x2="16" y2="6"></line>
                            <line x1="8" y1="2" x2="8" y2="6"></line>
                            <line x1="3" y1="10" x2="21" y2="10"></line>
                          </svg>
                        </template>
                        ดูรายการจองและคำร้องทั้งหมด
                      </AppButton>
                    </template>
                  </div>
                  <div v-else class="issd-footer issd-footer-edit">
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
              </transition>
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
const isInspectDashboardMode = computed(() => !isMapOverviewMode.value && String(route.query.mode || '') === 'inspect')
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
const dashboardStatusFilter = ref<'all' | 'available' | 'reserved' | 'occupied'>('all')
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
    { key: 'occupied', count: realMarketSummary.value.occupied, color: '#2563eb' }
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
}

const stallDrawerMode = ref<'detail' | 'edit' | 'full-detail'>('detail')
const stallDetailTab = ref<'info' | 'contracts' | 'subtenants' | 'payments' | 'maintenance' | 'bookings'>('info')
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
const stallDetailTabOptions: { value: 'info' | 'contracts' | 'subtenants' | 'payments' | 'maintenance' | 'bookings'; label: string }[] = [
  { value: 'info', label: 'ข้อมูลแผง' },
  { value: 'contracts', label: 'สัญญา' },
  { value: 'subtenants', label: 'ผู้เช่าช่วง' },
  { value: 'payments', label: 'ประวัติชำระเงิน' },
  { value: 'maintenance', label: 'ซ่อม/ขออนุญาต' },
  { value: 'bookings', label: 'การจองและคำร้อง' }
]

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
          customStallContracts.value[upper] = parsed
          return
        }
      } catch {}
    }
  }
  const status = String(stallEdit.value.stall_status || '').toLowerCase()
  if (status === 'occupied') {
    const tenant = stallEdit.value.name || 'คุณสมศักดิ์ มั่งคั่ง'
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
    customStallContracts.value[upper] = [defaultContract]
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
        createdAt: new Date().toISOString()
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
        createdAt: '2026-03-01T10:15:00.000Z'
      }
    ]
  } else {
    stallMaintenance.value = []
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
        createdAt: '2026-03-02T11:20:00.000Z'
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
        createdAt: '2026-02-01T14:45:00.000Z'
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
        createdAt: '2026-01-05T10:00:00.000Z'
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

const stallBookingsAndApps = computed(() => {
  const target = String(stallEdit.value.ident || stallEdit.value.stall_ident || '').trim().toUpperCase()
  if (!target) return []
  const items: Array<{
    key: string
    typeLabel: string
    refNumber: string
    applicantName: string
    dateText: string
    statusText: string
    badgeClass: string
  }> = []

  marketBookings.value.forEach((b: any) => {
    const codes = String(b.stall_number || b.stall_ident || b.stall_code || '')
      .split(',')
      .map(s => s.trim().toUpperCase())
    if (!codes.includes(target)) return
    const st = String(b.status || '').toLowerCase()
    const statusMap: Record<string, { label: string; cls: string }> = {
      pending: { label: 'รอตรวจสอบวันนัด', cls: 'is-reserved' },
      approved: { label: 'นัดหมายสำเร็จ', cls: 'is-available' },
      interested: { label: 'สนใจเช่าแผง', cls: 'is-available' },
      waiting_queue: { label: 'คิวสำรอง', cls: 'is-reserved' },
      cancelled: { label: 'ยกเลิก', cls: 'is-occupied' },
      rejected: { label: 'ปฏิเสธ', cls: 'is-occupied' }
    }
    const mapped = statusMap[st] || { label: b.status || 'รอดำเนินการ', cls: 'is-reserved' }
    items.push({
      key: `booking-${b.id || b.booking_number}`,
      typeLabel: 'นัดเข้าชมพื้นที่',
      refNumber: b.booking_number || `#${b.id}`,
      applicantName: b.tenant_name || b.applicant_name || '-',
      dateText: b.visit_date ? `${b.visit_date}${b.visit_time_slot ? ' · ' + b.visit_time_slot : ''}` : formatStallDate(b.created_at),
      statusText: mapped.label,
      badgeClass: mapped.cls
    })
  })

  marketApplications.value.forEach((a: any) => {
    const codes = String(a.stall_number || a.stall_ident || a.stall_code || '')
      .split(',')
      .map(s => s.trim().toUpperCase())
    if (!codes.includes(target)) return
    const st = String(a.status || '').toLowerCase()
    const isApproved = st === 'approved'
    const isRejected = st === 'rejected'
    items.push({
      key: `app-${a.id || a.application_number}`,
      typeLabel: 'คำร้องขอเช่าแผง',
      refNumber: a.application_number || `#${a.id}`,
      applicantName: a.applicant_name || '-',
      dateText: formatStallDate(a.created_at),
      statusText: isApproved ? 'อนุมัติแล้ว' : (isRejected ? 'ไม่อนุมัติ' : 'รอพิจารณาอนุมัติ'),
      badgeClass: isApproved ? 'is-available' : (isRejected ? 'is-occupied' : 'is-reserved')
    })
  })

  return items
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

const loadStallExtraDetails = async () => {
  const code = String(stallEdit.value.ident || stallEdit.value.stall_ident || '').trim()
  if (!code) return
  loadStallContracts(code)
  loadStallSubtenants(code)
  loadStallMaintenance(code)
  loadStallPayments(code)
  const stored = stalls.value.find(s => String(s.ident || s.stall_ident).toUpperCase() === code.toUpperCase())
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

const openStallFullDetail = () => {
  stallDrawerMode.value = 'full-detail'
  stallDetailTab.value = 'info'
  void loadStallExtraDetails()
}

const backToMiniDetail = () => {
  stallDrawerMode.value = 'detail'
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
const emptyStall = () => ({ id: 0, ident: '', name: '', description: '', category: '', zone_id: 1, zone_name: 'โซน A', project_name: selectedProject.value?.name || '', floor: '1', width: 2, depth: 2, daily_rent: 150, monthly_rent: 1800, stall_status: 'available', coordinates_geojson: '', published: false })
const stallEdit = ref(emptyStall())
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

// ข้อมูลตัวเลขสรุปที่แสดงบนแถบฟิลเตอร์ด้านบน: สอดคล้องกับโซนที่เลือก หรือทั้งตลาด
const currentDisplayedSummary = computed(() => {
  if (selectedProject.value) {
    const p = selectedProject.value
    let total = 0
    let available = 0
    let reserved = 0
    let occupied = 0
    const floorMap = getAllFloorLayoutsForProject(p)
    floorMap.forEach((layout) => {
      layout.stalls.forEach((stall: any) => {
        total += 1
        const code = String(stall.code || stall.ident || stall.stall_ident || '').trim()
        const status = resolveRealStallStatus(code, p, stall)
        if (status === 'reserved') reserved += 1
        else if (status === 'occupied') occupied += 1
        else available += 1
      })
    })
    return { total, available, reserved, occupied }
  }
  return {
    total: totalMarketStalls.value,
    available: stallStatusSummary.value.available,
    reserved: stallStatusSummary.value.reserved,
    occupied: stallStatusSummary.value.occupied
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

const selectProject = (project: ProjectZone) => { selectedProject.value = project }
const clearSelectedProject = () => {
  selectedProject.value = null
  closeStall()
  if (route.query.project || route.query.stall) {
    const nextQuery: Record<string, any> = { ...route.query }
    delete nextQuery.project
    delete nextQuery.stall
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
const closeStall = () => {
  showStallEditor.value = false
  showInspectStallOverlay.value = false
  showInspectPhotoModal.value = false
  resetStallErrors()
  stallDrawerMode.value = 'detail'
  stallDetailTab.value = 'info'
  pendingStallLayout.value = null
  contracts.value = []
  stallHistory.value = []
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
    stallDrawerMode.value = 'detail'
  }
}
const selectStall = async (mapStall: any) => {
  pendingStallLayout.value = null
  resetStallErrors()
  const code = String(mapStall.code || mapStall.ident || '')
  const stored = stalls.value.find(s => String(s.ident || s.stall_ident).toUpperCase() === code.toUpperCase())
  const geometry = mapStall.coordinates?.length ? JSON.stringify({ type: 'Polygon', coordinates: [mapStall.coordinates] }) : ''
  contracts.value = []
  stallHistory.value = []
  stallDetailTab.value = 'info'
  loadStallContracts(code)
  if (!stored?.id) {
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
    stallDrawerMode.value = 'detail'
    showInspectStallOverlay.value = true
    showStallEditor.value = false
    return
  }
  stallEdit.value = {
    ...emptyStall(),
    ...stored,
    coordinates_geojson: stored.coordinates_geojson || geometry,
    stall_status: resolveRealStallStatus(code, selectedProject.value, stored)
  }
  stallDrawerMode.value = 'detail'
  showInspectStallOverlay.value = true
  showStallEditor.value = false
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
}

const handleBeforeUnload = (event: BeforeUnloadEvent) => {
  if (!isInspectMode.value && hasUnsavedChanges.value) {
    event.preventDefault()
    event.returnValue = ''
  }
}

onMounted(() => {
  void loadData()
  loadAllCustomContracts()
  if (import.meta.client) {
    window.addEventListener('beforeunload', handleBeforeUnload)
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
.dash-donut-legend-list{display:grid;gap:8px}.dash-legend-row{width:100%;display:flex;align-items:center;justify-content:space-between;gap:10px;padding:9px 12px;border:1px solid #eef2f6;border-radius:10px;background:#f8fafc}.legend-left{display:inline-flex;align-items:center;gap:8px;color:#334155;font-size:.875rem;font-weight:600}.legend-left i{width:10px;height:10px;border-radius:50%;display:inline-block}.dot-available{background:#10b981}.dot-reserved{background:#f59e0b}.dot-occupied{background:#2563eb}.legend-right{display:inline-flex;align-items:center;gap:8px}.legend-right strong{color:#0f172a;font-size:.9375rem;font-weight:700}
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
.inspect-selected-stall-drawer{position:absolute;top:62px;right:14px;width:338px;max-height:calc(100% - 76px);display:flex;flex-direction:column;pointer-events:auto;z-index:25;background:rgba(255,255,255,.98);backdrop-filter:blur(14px);border:1px solid #e2e8f0;border-radius:18px;box-shadow:0 16px 38px rgba(15,23,42,.16);overflow:hidden;transition:width .24s cubic-bezier(.16,1,.3,1)}
.inspect-selected-stall-drawer.is-edit-mode{width:min(440px,calc(100% - 28px))}
.inspect-selected-stall-drawer.is-full-detail-mode{width:min(560px,calc(100vw - 36px));z-index:60}
@media(min-width:769px){.inspect-selected-stall-drawer.is-management-mode{bottom:100px;max-height:calc(100% - 176px)}}
.issd-header{padding:14px 16px 11px;border-bottom:1px solid #f1f5f9;background:#fff;flex-shrink:0}.issd-header-top{display:flex;align-items:center;justify-content:space-between;gap:10px}.issd-title-group{display:flex;align-items:center;gap:7px;flex-wrap:wrap;min-width:0}.issd-title-group strong{font-size:.9rem;font-weight:700;color:#0f172a}.issd-project-sub{margin-top:3px;font-size:.75rem;font-weight:500;color:#64748b;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
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
}
.issd-tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  height: 31px;
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
.issd-body {
  flex: 1 1 auto;
  overflow-y: auto;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  scrollbar-width: thin;
  scrollbar-color: rgba(148, 163, 184, 0.35) transparent;
}
.issd-body::-webkit-scrollbar {
  width: 5px;
}
.issd-body::-webkit-scrollbar-track {
  background: transparent;
}
.issd-body::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.35);
  border-radius: 999px;
}
.issd-body::-webkit-scrollbar-thumb:hover {
  background: rgba(100, 116, 139, 0.55);
}
.issd-body-full{display:flex;flex-direction:column;gap:14px;padding:14px 16px}
.issd-photo-banner.is-full{height:124px}
.issd-kpi-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:10px;margin:0}

.issd-kpi-card{padding:10px 12px;border-radius:11px;border:1px solid #e2e8f0;background:#f8fafc;display:grid;gap:2px}
.issd-kpi-card span{color:#64748b;font-size:.75rem;font-weight:500}
.issd-kpi-card strong{color:#0f172a;font-size:1.15rem;font-weight:700;line-height:1.2}
.issd-kpi-card.is-green strong{color:#10b981}
.issd-kpi-card em{font-style:normal;color:#64748b;font-size:.75rem;font-weight:500}
.issd-fact-section{display:flex;flex-direction:column;gap:8px}
.issd-section-title{font-size:.8125rem;font-weight:700;color:#1e293b}
.issd-fact-grid{display:grid;grid-template-columns:repeat(2,1fr);gap:8px;margin:0}
.issd-fact-cell{padding:9px 12px;border:1px solid #e2e8f0;border-radius:9px;background:#fff}
.issd-fact-cell.full{grid-column:1/-1}
.issd-fact-cell dt{color:#64748b;font-size:.75rem;font-weight:500}
.issd-fact-cell dd{margin:3px 0 0;color:#1e293b;font-size:.8125rem;font-weight:600;word-break:break-word}
.issd-booking-list,.issd-contract-list{display:flex;flex-direction:column;gap:10px}
.issd-booking-card,.issd-contract-card{padding:12px 14px;border:1px solid #e2e8f0;border-radius:12px;background:#fff;display:flex;flex-direction:column;gap:6px;transition:border-color .15s ease,box-shadow .15s ease}
.issd-booking-card:hover,.issd-contract-card:hover{border-color:#cbd5e1;box-shadow:0 2px 8px rgba(15,23,42,.04)}
.issd-card-top-row{display:flex;align-items:center;justify-content:space-between;gap:8px}
.issd-card-type-tag{font-size:.75rem;font-weight:600;color:#2563eb;background:#eff6ff;padding:2px 8px;border-radius:6px}
.issd-card-title,.issd-contract-num{font-size:.875rem;font-weight:700;color:#0f172a}
.issd-contract-tenant{font-size:.8125rem;color:#475569}
.issd-contract-tenant strong{color:#1e293b;font-weight:600}
.issd-card-meta-row{display:flex;align-items:center;justify-content:space-between;gap:10px;font-size:.75rem;color:#64748b;flex-wrap:wrap}
.issd-card-meta-row .is-rent{color:#10b981;font-weight:700;font-size:.8125rem}
.issd-timeline-list{display:flex;flex-direction:column;padding:4px 4px 10px}
.issd-timeline-item{position:relative;display:flex;gap:12px;padding-bottom:16px}
.issd-timeline-item:not(:last-child)::before{content:'';position:absolute;left:5px;top:13px;bottom:0;width:1px;background:#e2e8f0}
.issd-timeline-dot{position:relative;z-index:1;width:11px;height:11px;margin-top:4px;flex:0 0 11px;border-radius:50%;background:#3b82f6;box-shadow:0 0 0 3px #dbeafe}
.issd-timeline-copy{flex:1;min-width:0}
.issd-timeline-copy strong{display:block;font-size:.8125rem;color:#1e293b;font-weight:600}
.issd-timeline-copy time{display:block;color:#64748b;font-size:.75rem;margin-top:2px}
.issd-timeline-copy p{margin:6px 0 0;padding:6px 10px;border-radius:8px;background:#f8fafc;color:#475569;font-size:.75rem;line-height:1.45}
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
.issd-status-pill{display:inline-flex;align-items:center;padding:2px 9px;border-radius:999px;font-size:.75rem;font-weight:600;border:1px solid transparent}.issd-status-pill.is-available{background:#dcfce7;color:#15803d;border-color:#bbf7d0}.issd-status-pill.is-reserved{background:#fef3c7;color:#d97706;border-color:#fde68a}.issd-status-pill.is-occupied{background:#f1f5f9;color:#475569;border-color:#cbd5e1}
.issd-photo-banner{position:relative;height:108px;flex-shrink:0;border-radius:11px;overflow:hidden;border:1px solid #e2e8f0;background:#0f172a;cursor:pointer}.issd-photo-banner img{width:100%;height:100%;object-fit:cover;display:block;transition:transform .25s ease}.issd-photo-banner:hover img{transform:scale(1.05)}.issd-photo-count-badge{position:absolute;right:8px;bottom:8px;display:inline-flex;align-items:center;gap:4px;padding:2px 8px;border-radius:7px;background:rgba(15,23,42,.82);backdrop-filter:blur(4px);color:#fff;font-size:.75rem;font-weight:600}
.issd-detail-list{margin:0;flex-shrink:0;border:1px solid #e2e8f0;border-radius:12px;background:#fff;overflow:hidden}.issd-detail-row{display:flex;align-items:center;justify-content:space-between;gap:10px;padding:7px 12px;font-size:.8125rem}.issd-detail-row+.issd-detail-row{border-top:1px solid #f1f5f9}.issd-detail-row dt{color:#64748b;font-size:.75rem;font-weight:500}.issd-detail-row dd{margin:0;color:#1e293b;font-size:.8125rem;font-weight:600;text-align:right}.issd-detail-row dd.is-rent{color:#15803d;font-weight:700}.issd-detail-row dd.is-bold{color:#0f172a;font-weight:700}
.issd-occupancy-box{padding:9px 12px;flex-shrink:0;border-radius:12px;border:1px solid #e2e8f0;background:#f8fafc;display:grid;gap:2px}.issd-occupancy-box.is-available{background:#f0fdf4;border-color:#bbf7d0}.issd-occupancy-box.is-reserved{background:#fffbeb;border-color:#fde68a}.issd-occupancy-box.is-occupied{background:#f8fafc;border-color:#cbd5e1}.issd-occupancy-head{display:flex;align-items:center;justify-content:space-between;gap:8px}.issd-occupancy-head span{font-size:.75rem;font-weight:600;color:#64748b}.issd-occupancy-head small{font-size:.75rem;font-weight:700;color:#2563eb}.issd-occupancy-box strong{font-size:.8125rem;font-weight:700;color:#0f172a;line-height:1.35}.issd-occupancy-box p{margin:0;font-size:.75rem;color:#475569;line-height:1.35}
.issd-footer{padding:10px 14px 14px;background:#fff;border-top:1px solid #f1f5f9;box-shadow:0 -4px 12px rgba(15,23,42,.04);flex-shrink:0}.issd-cta-btn{width:100%;box-sizing:border-box;padding:10px 16px;border:0;border-radius:11px;background:linear-gradient(180deg,#3b82f6 0%,#2563eb 100%);color:#fff;font-family:inherit;font-size:.875rem;font-weight:700;text-decoration:none;display:inline-flex;align-items:center;justify-content:center;gap:8px;cursor:pointer;box-shadow:0 4px 14px rgba(37,99,235,.32);transition:transform .15s ease,box-shadow .15s ease}.issd-cta-btn:hover{transform:translateY(-1px);box-shadow:0 6px 18px rgba(37,99,235,.4)}
.issd-footer-full{display:flex;align-items:center;justify-content:space-between;gap:10px}
.issd-footer-full :deep(.app-button),
.issd-footer-full > button,
.issd-footer-full > a {
  flex: 1;
}
.inspect-photo-Modal,.inspect-photo-modal{width:min(820px,94vw);background:#0f172a;border-radius:18px;overflow:hidden;box-shadow:0 30px 60px -20px rgba(0,0,0,.6);display:flex;flex-direction:column}.ipm-head{display:flex;align-items:center;justify-content:space-between;padding:14px 18px;color:#fff}.ipm-head strong{display:block;font-size:1rem;font-weight:700}.ipm-head span{display:block;font-size:.8125rem;color:#94a3b8;margin-top:2px}.ipm-close{width:34px;height:34px;border-radius:10px;border:0;background:rgba(255,255,255,.12);color:#fff;display:inline-flex;align-items:center;justify-content:center;cursor:pointer}.ipm-close svg{width:18px;height:18px}.ipm-stage{position:relative;height:min(480px,58vh);background:#000;display:flex;align-items:center;justify-content:center}.ipm-main-img{width:100%;height:100%;object-fit:contain}.ipm-nav{position:absolute;top:50%;transform:translateY(-50%);width:40px;height:40px;border-radius:50%;border:0;background:rgba(255,255,255,.9);color:#0f172a;font-size:1.35rem;font-weight:700;cursor:pointer;display:grid;place-items:center}.ipm-nav.is-prev{left:14px}.ipm-nav.is-next{right:14px}.ipm-counter{position:absolute;bottom:12px;left:50%;transform:translateX(-50%);background:rgba(15,23,42,.78);color:#fff;font-size:.75rem;font-weight:600;padding:3px 12px;border-radius:999px}.ipm-thumbs{display:flex;justify-content:center;gap:8px;padding:12px 18px 16px}.ipm-thumb{width:76px;height:52px;border-radius:8px;overflow:hidden;border:2px solid transparent;padding:0;background:#1e293b;cursor:pointer;opacity:.6}.ipm-thumb.active{border-color:#60a5fa;opacity:1}.ipm-thumb img{width:100%;height:100%;object-fit:cover;display:block}
@media(max-width:768px){.inspect-selected-stall-drawer{position:fixed;left:0;right:0;bottom:0;top:auto;width:100%!important;max-height:85vh;border-radius:20px 20px 0 0;z-index:99}.inspect-selected-stall-drawer.is-full-detail-mode{width:100%!important;max-height:90vh}}
.map-edit-status{display:flex;gap:10px;align-items:flex-start;padding:13px;background:#eff6ff;border:1px solid #bfdbfe;border-radius:10px;color:#1e40af}.map-edit-status strong,.map-edit-status small{display:block}.map-edit-status strong{font-size:.875rem}.map-edit-status small{margin-top:3px;font-size:.75rem;color:#476587}.map-edit-dot{width:10px;height:10px;flex:0 0 10px;border-radius:50%;background:#3b5bff;margin-top:5px;box-shadow:0 0 0 4px #dbeafe}
.zone-form,.facility-form{display:grid;grid-template-columns:1fr 1fr;gap:13px}.zone-form label,.facility-form label{display:grid;align-content:start;gap:6px;color:#475569;font-size:.8125rem;font-weight:600}.zone-form .full,.facility-form .full{grid-column:1/-1}.req-star{color:#ef4444;font-weight:700}.stall-detail-summary{display:grid;gap:5px;padding:16px;margin:20px 0 14px;border:1px solid #dbeafe;border-radius:12px;background:#f8fbff}.stall-detail-summary strong{color:#1e293b;font-size:1rem}.stall-detail-summary>span:last-child{color:#64748b;font-size:.875rem}.stall-detail-summary .soft-badge{justify-self:start}.status-available{background:#dcfce7;color:#15803d!important}.status-reserved{background:#fef3c7;color:#d97706!important}.status-occupied{background:#f1f5f9;color:#64748b!important}.stall-detail-tabs{display:flex!important;width:100%;margin:0 0 14px}.stall-detail-tabs :deep(.app-segmented-option){flex:1}.stall-detail-grid{display:grid;grid-template-columns:1fr 1fr;gap:10px;margin:0}.stall-detail-grid>div{padding:12px 14px;border:1px solid #e2e8f0;border-radius:10px;background:#fff}.stall-detail-grid .full{grid-column:1/-1}.stall-detail-grid dt{color:#64748b;font-size:.8125rem;font-weight:500}.stall-detail-grid dd{margin:4px 0 0;color:#1e293b;font-size:.9375rem;font-weight:600;word-break:break-word}.detail-section.is-tab-content{margin-top:0}.detail-section-heading{display:flex;align-items:center;justify-content:space-between;gap:12px}.detail-section-heading h3{margin:0}
.issd-tab-header-row {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-bottom: 12px;
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
  transition: border-color .15s ease, box-shadow .15s ease;
}
.issd-subtenant-card:hover,
.issd-payment-card:hover,
.issd-maintenance-card:hover {
  border-color: #cbd5e1;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

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

.is-fullscreen-map-panel :deep(.map-floating-top-controls) {
  display: none !important;
}

@media (min-width: 769px) {
  .is-fullscreen-map-panel :deep(.stall-mode-panel) {
    top: 68px !important;
    left: 18px !important;
    bottom: auto !important;
    height: fit-content !important;
    max-height: calc(100% - 84px) !important;
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

.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer) {
  top: 68px !important;
  right: 18px !important;
  max-height: calc(100% - 84px) !important;
}

.is-fullscreen-map-panel :deep(.inspect-selected-stall-drawer.is-full-detail-mode) {
  width: min(560px, calc(100vw - 36px)) !important;
}

/* Market information overview: six KPI cards floating above the GIS map. */
.market-map-kpi-strip {
  position: absolute;
  top: 16px;
  left: 20px;
  right: 76px;
  z-index: 110;
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 12px;
  pointer-events: none;
}

.market-map-kpi {
  min-width: 0;
  min-height: 82px;
  padding: 13px 14px;
  display: flex;
  align-items: center;
  gap: 12px;
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
  padding: 12px 13px;
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
  text-overflow: ellipsis;
  white-space: nowrap;
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
.market-map-kpi.is-overdue .market-map-kpi-icon,
.market-map-kpi.is-overdue .market-map-kpi-copy em { background: #fef2f2; color: #dc2626; }
.market-map-kpi.is-expiring .market-map-kpi-icon,
.market-map-kpi.is-expiring .market-map-kpi-copy em { background: #f5f3ff; color: #7c3aed; }

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
  top: 114px !important;
}

@media (max-width: 1280px) {
  .market-map-kpi-strip {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
  .is-fullscreen-map-panel :deep(.map-project-guidance),
  .is-fullscreen-map-panel :deep(.map-floating-floor-preview-pill) {
    top: 208px !important;
  }
}

@media (max-width: 720px) {
  .market-map-kpi-strip {
    top: 10px;
    left: 10px;
    right: 10px;
    display: flex;
    gap: 8px;
    overflow-x: auto;
    scrollbar-width: none;
  }
  .market-map-kpi-strip::-webkit-scrollbar { display: none; }
  .market-map-kpi { min-width: 184px; min-height: 72px; padding: 10px 12px; }
  .market-map-kpi.active { padding: 9px 11px; }
  .market-map-tool-rail { right: 10px; border-right: 1px solid #e2e8f0; border-radius: 14px; }
  .is-fullscreen-map-panel :deep(.map-project-guidance),
  .is-fullscreen-map-panel :deep(.map-floating-floor-preview-pill) { top: 96px !important; }
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
</style>
