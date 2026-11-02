<template>
  <div class="contracts-page">
    <AppHeader />

    <div class="page-container">
      <div class="page-header-row">
        <div>
          <h1 class="main-page-title">สัญญา</h1>
          <p class="page-subtitle">ติดตามและจัดการสัญญาเช่าแผง</p>
        </div>
      </div>
      <!-- Loading State -->
      <div v-if="loading" class="contracts-loading-state">
        <div class="loading-spinner-row">
          <svg class="animate-spin" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="2" x2="12" y2="6"></line><line x1="12" y1="18" x2="12" y2="22"></line><line x1="4.93" y1="4.93" x2="7.76" y2="7.76"></line><line x1="16.24" y1="16.24" x2="19.07" y2="19.07"></line><line x1="2" y1="12" x2="6" y2="12"></line><line x1="18" y1="12" x2="22" y2="12"></line><line x1="4.93" y1="19.07" x2="7.76" y2="16.24"></line><line x1="16.24" y1="7.76" x2="19.07" y2="4.93"></line></svg>
          <span>กำลังโหลดข้อมูลสัญญา...</span>
        </div>
      </div>

      <!-- Standard Table/Card Container -->
      <div v-else class="table-card-container">
        <!-- Card Header: Title + Subtitle on Left, Search + View Mode + Actions on Right -->
        <div class="card-header-row">
          <div class="header-left">
            <h2 class="section-title">{{ isTenant ? 'แผงเช่าและสัญญาของฉัน' : 'รายการสัญญาเช่าทั้งหมด' }}</h2>
            <div class="item-count-subtitle">ทั้งหมด {{ filteredContracts.length }} สัญญา</div>
          </div>

          <div class="header-right">
            <!-- Search Input with Clear Button -->
            <div class="search-input-wrapper">
              <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="ค้นหาเลขที่สัญญา, แผงค้า, ผู้เช่า..."
                class="search-input"
                @input="filterContracts(); currentPage = 1"
              />
              <button
                v-if="searchQuery"
                type="button"
                class="clear-search-btn"
                aria-label="ล้างการค้นหา"
                @click="searchQuery = ''; filterContracts(); currentPage = 1"
              >
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              </button>
            </div>

            <!-- View Mode Switcher: Cards vs Table -->
            <AppSegmentedControl
              v-model="viewMode"
              :options="viewModeOptions"
              aria-label="เลือกรูปแบบการแสดงรายการสัญญา"
            />

            <AppButton variant="outline" size="md" title="ออกรายงาน" @click="exportReport">
              <template #prefix><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"/><polyline points="16 6 12 2 8 6"/><line x1="12" y1="2" x2="12" y2="15"/></svg></template>
              ออกรายงาน
            </AppButton>
          </div>
        </div>

        <!-- KPI Ribbon Cards -->
        <div class="officer-kpi-ribbon officer-status-kpi-ribbon">
          <!-- 1. ทั้งหมด -->
          <div
            class="kpi-ribbon-card kpi-all"
            :class="{ active: selectedStatus === '' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = ''; filterContracts(); currentPage = 1"
            @keydown.enter.prevent="selectedStatus = ''; filterContracts(); currentPage = 1"
            @keydown.space.prevent="selectedStatus = ''; filterContracts(); currentPage = 1"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path><rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">ทั้งหมด</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ contracts.length }}</span>
                <span class="kpi-unit-pill">ฉบับ</span>
              </div>
            </div>
          </div>

          <!-- 2. รอดำเนินการ (Pending) -->
          <div
            class="kpi-ribbon-card kpi-pending"
            :class="{ active: selectedStatus === 'pending' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'pending'; pendingContractFilter = 'all'; filterContracts(); currentPage = 1"
            @keydown.enter.prevent="selectedStatus = 'pending'; pendingContractFilter = 'all'; filterContracts(); currentPage = 1"
            @keydown.space.prevent="selectedStatus = 'pending'; pendingContractFilter = 'all'; filterContracts(); currentPage = 1"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">รอดำเนินการ</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countContractsByStatus('pending') }}</span>
                <span class="kpi-unit-pill">ฉบับ</span>
              </div>
            </div>
          </div>

          <!-- 3. ใช้งานอยู่ (Active) -->
          <div
            class="kpi-ribbon-card kpi-applied"
            :class="{ active: selectedStatus === 'active' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'active'; filterContracts(); currentPage = 1"
            @keydown.enter.prevent="selectedStatus = 'active'; filterContracts(); currentPage = 1"
            @keydown.space.prevent="selectedStatus = 'active'; filterContracts(); currentPage = 1"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">ใช้งานอยู่</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countContractsByStatus('active') }}</span>
                <span class="kpi-unit-pill">ฉบับ</span>
              </div>
            </div>
          </div>

          <!-- 4. หมดสัญญาแล้ว (Expired) -->
          <div
            class="kpi-ribbon-card kpi-expired"
            :class="{ active: selectedStatus === 'expired' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'expired'; filterContracts(); currentPage = 1"
            @keydown.enter.prevent="selectedStatus = 'expired'; filterContracts(); currentPage = 1"
            @keydown.space.prevent="selectedStatus = 'expired'; filterContracts(); currentPage = 1"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">หมดสัญญา</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countContractsByStatus('expired') }}</span>
                <span class="kpi-unit-pill">ฉบับ</span>
              </div>
            </div>
          </div>

          <!-- 5. สิ้นสุด/ยกเลิก (Terminated) -->
          <div
            class="kpi-ribbon-card kpi-cancelled"
            :class="{ active: selectedStatus === 'terminated' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'terminated'; filterContracts(); currentPage = 1"
            @keydown.enter.prevent="selectedStatus = 'terminated'; filterContracts(); currentPage = 1"
            @keydown.space.prevent="selectedStatus = 'terminated'; filterContracts(); currentPage = 1"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">ยกเลิก/สิ้นสุด</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countContractsByStatus('terminated') }}</span>
                <span class="kpi-unit-pill">ฉบับ</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Pending Subfilters -->
        <transition name="fade-slide">
          <div v-if="selectedStatus === 'pending'" class="tenant-subfilter-bar" style="margin-top: -12px; margin-bottom: 22px;">
            <button
              v-for="option in pendingContractOptions"
              :key="option.value"
              type="button"
              class="subfilter-chip"
              :class="[getPendingChipClass(option.value), { active: pendingContractFilter === option.value }]"
              :aria-pressed="pendingContractFilter === option.value"
              @click="pendingContractFilter = option.value; filterContracts(); currentPage = 1"
            >
              <span v-if="option.value !== 'all'" class="subfilter-dot"></span>
              <span>{{ option.label }}</span>
              <span class="subfilter-badge">{{ option.count }}</span>
            </button>
          </div>
        </transition>

        <!-- Empty State -->
        <AppEmptyState
          v-if="filteredContracts.length === 0"
          variant="table"
          :icon="searchQuery || selectedStatus ? 'search' : 'document'"
          :title="searchQuery ? `ไม่พบสัญญาเช่าที่ตรงกับ “${searchQuery}”` : (selectedStatus ? `ไม่มีรายการในสถานะ “${selectedContractStatusLabel}”` : 'ยังไม่มีรายการสัญญาเช่า')"
          :description="searchQuery ? 'ลองตรวจสอบคำค้นหาหรือล้างตัวกรอง' : (isTenant ? 'สัญญาจะแสดงเมื่อคำร้องขอเช่าได้รับอนุมัติ' : 'สัญญาที่ผู้เช่าส่งจะแสดงที่นี่')"
        >
          <template v-if="searchQuery || selectedStatus" #actions>
            <AppButton variant="outline" size="sm" @click="searchQuery = ''; selectedStatus = ''; filterContracts(); currentPage = 1">ล้างตัวกรอง</AppButton>
          </template>
        </AppEmptyState>

        <!-- Table View -->
        <div v-else-if="viewMode === 'table'" class="officer-table-container" role="region" aria-label="ตารางสัญญา" tabindex="0">
          <table class="officer-custom-table">
            <thead>
              <tr>
                <th scope="col">เลขที่สัญญา</th>
                <th scope="col">แผง / ตลาด</th>
                <th scope="col">ร้านค้า</th>
                <th scope="col">ผู้เช่า</th>
                <th scope="col">ระยะสัญญา</th>
                <th scope="col">ค่าเช่า</th>
                <th scope="col">เงินประกัน</th>
                <th scope="col">สถานะ</th>
                <th scope="col" class="action-heading">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="c in paginatedContracts" :key="c.id">
                <td><strong class="table-code-blue">{{ c.contract_number }}</strong><div class="cell-sub-cycle">วันที่ {{ formatContractDate(c.created_at) }}</div></td>
                <td><span class="badge badge-occupied">{{ c.stall_ident || '-' }}</span><div class="cell-sub-category">{{ contractMarket(c) }}</div></td>
                <td><div class="cell-shop-name">{{ c.shop_name || '-' }}</div><div class="cell-sub-category">{{ c.business_type || '-' }}</div></td>
                <td><div class="cell-tenant-name">{{ c.tenant_name || '-' }}</div><a v-if="c.tenant_phone" :href="'tel:' + c.tenant_phone" class="cell-phone-link">{{ c.tenant_phone }}</a></td>
                <td><div class="cell-duration-text">{{ formatContractDate(c.start_date) }} – {{ formatContractDate(c.end_date) }}</div></td>
                <td><strong class="cell-bold-amount">฿{{ Number(c.rental_fee || 0).toLocaleString() }}</strong><div class="cell-sub-cycle">{{ c.billing_cycle === 'daily' ? 'ต่อวัน' : 'ต่อเดือน' }}</div></td>
                <td>฿{{ Number(c.security_deposit || 0).toLocaleString() }}</td>
                <td><span class="status-pill-badge" :class="getStatusBadgeClass(c.status)">{{ getStatusLabel(c.status) }}</span></td>
                <td class="action-cell">
                  <div class="table-actions-cell">
                    <AppButton v-if="canReviewContract(c)" variant="primary" size="sm" @click="approveContract(c)">ลงนาม</AppButton>
                    <AppButton v-if="isStaff && c.status === 'draft'" variant="primary" size="sm" @click="openContractDraft(c)">กำหนดรายละเอียด</AppButton>
                    <AppButton v-if="isTenant && c.status === 'revision_required'" variant="primary" size="sm" @click="openContractRevision(c)">แก้ไข</AppButton>
                    <AppButton variant="secondary" size="sm" @click="viewContractDoc(c)">
                      <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8z"/></svg></template>
                      ดู
                    </AppButton>
                    <AppButton v-if="canTerminateContract(c)" variant="danger-outline" size="sm" @click="terminateContract(c)">สิ้นสุดสัญญา</AppButton>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Card View (Accordion grouped view matching Application) -->
        <div v-else class="application-cards-container">
          <div v-if="groupedContracts.length > 0" class="officer-groups-wrapper">
            <div
              v-for="group in groupedContracts"
              :key="group.key"
              class="project-stall-entry"
            >
              <!-- Accordion Header -->
              <button
                type="button"
                class="officer-project-header"
                :class="{ open: isContractGroupOpen(group.key) }"
                @click="toggleContractGroup(group.key)"
              >
                <div class="header-left-content">
                  <svg class="chevron-indicator" :class="{ rotated: isContractGroupOpen(group.key) }" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="6 9 12 15 18 9"></polyline></svg>
                  <span class="accordion-header-icon" aria-hidden="true">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
                  </span>
                  <span class="market-title">{{ group.name }}</span>
                  <span class="market-count-pill">{{ group.items.length }} สัญญา</span>
                </div>
              </button>

              <!-- Accordion Body -->
              <div
                v-show="isContractGroupOpen(group.key)"
                class="market-accordion-body officer-cards-list in-project"
              >
                <div
                  v-for="c in group.items"
                  :key="c.id"
                  class="booking-item-card"
                  :class="[getContractBorderClass(c.status), { 'collapsed-card': !isContractCardOpen(c.id) }]"
                >
                  <!-- Card Top Row -->
                  <div class="card-top-row">
                    <!-- Top Line: Contract Number & Status Badge + Chevron -->
                    <div class="card-top-main">
                      <div class="top-row-left">
                        <div class="circle-status-icon blue-circle">
                          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                            <polyline points="14 2 14 8 20 8"></polyline>
                            <line x1="16" y1="13" x2="8" y2="13"></line>
                            <line x1="16" y1="17" x2="8" y2="17"></line>
                            <polyline points="10 9 9 9 8 9"></polyline>
                          </svg>
                        </div>

                        <div class="booking-number-block">
                          <div class="booking-label-row">
                            <span class="booking-label">เลขที่สัญญา</span>
                          </div>
                          <div class="booking-id-text">#{{ c.contract_number }}</div>
                          <div class="officer-booked-at">
                            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                            <span>วันที่ทำสัญญา {{ formatContractDate(c.created_at) }}</span>
                          </div>
                        </div>
                      </div>

                      <div class="top-row-right">
                        <AppButton variant="outline" size="sm" title="ตัวอย่างเอกสาร" @click.stop="viewContractDoc(c)">
                          <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg></template>
                          ตัวอย่างเอกสาร
                        </AppButton>
                        <span class="status-pill-badge" :class="getStatusBadgeClass(c.status)">
                          {{ getStatusLabel(c.status) }}
                        </span>
                        <button class="card-collapse-btn" type="button" @click.stop="toggleContractCard(c.id)" title="ย่อ/ขยายการ์ด">
                          <svg
                            class="card-chevron"
                            :class="{ rotated: isContractCardOpen(c.id) }"
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

                    <!-- Compact context row: stall, project, market, and applicant/tenant -->
                    <div class="card-meta-tags-row compact-officer-meta">
                      <div class="compact-meta-item compact-meta-stall">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                        <span>แผงค้า</span>
                        <strong>{{ c.stall_ident || '-' }}</strong>
                      </div>

                      <div class="compact-meta-item">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V9"></path><path d="M9 22V12h6v10"></path></svg>
                        <span>{{ contractProject(c) }}</span>
                      </div>

                      <div class="compact-meta-item">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M3 21h18M3 7v14M21 7v14M6 7V3h12v4M9 21v-6a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v6"></path></svg>
                        <span>{{ contractMarket(c) }}</span>
                      </div>

                      <div class="compact-meta-item compact-meta-person">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                        <strong>{{ c.tenant_name || '-' }}</strong>
                        <a v-if="c.tenant_phone" :href="'tel:' + c.tenant_phone" @click.stop>{{ c.tenant_phone }}</a>
                      </div>
                    </div>
                  </div>

                  <!-- Collapsed Next Step Bar -->
                  <div
                    v-if="!isContractCardOpen(c.id)"
                    class="collapsed-next-step"
                    :class="getContractProgressTone(c.status)"
                  >
                    <div class="collapsed-next-step-content">
                      <span class="collapsed-next-step-label">ความคืบหน้า</span>
                      <strong>{{ contractProgressText(c) }}</strong>
                    </div>

                    <div class="collapsed-actions-group">
                      <AppButton v-if="canReviewContract(c) && c.status !== 'pending_bureau_director'" variant="outline" size="sm" @click.stop="rejectContract(c)"><template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 4v6h6"/><path d="M3.5 15a9 9 0 1 0 2-9.4L1 10"/></svg></template>ส่งกลับแก้ไข</AppButton>
                      <AppButton v-if="isStaff && c.status === 'draft'" variant="primary" size="sm" @click.stop="openContractDraft(c)">กำหนดรายละเอียด</AppButton>
                      <AppButton v-if="isTenant && c.status === 'revision_required'" variant="primary" size="sm" @click.stop="openContractRevision(c)">แก้ไขและยื่นใหม่</AppButton>
                      <AppButton variant="secondary" size="sm" @click.stop="viewContractDoc(c)">
                        <template #prefix><svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"></circle><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8z"></path></svg></template>
                        รายละเอียด
                      </AppButton>
                      <AppButton v-if="canReviewContract(c)" variant="primary" size="sm" @click.stop="approveContract(c)">{{ c.status === 'pending_bureau_director' ? 'ลงนามและอนุมัติ' : 'ลงนามและส่งต่อ' }}</AppButton>
                      <AppButton v-if="canTerminateContract(c)" variant="danger-outline" size="sm" @click.stop="terminateContract(c)">สิ้นสุดสัญญา</AppButton>
                    </div>
                  </div>

                  <!-- Expanded Card Body -->
                  <div v-show="isContractCardOpen(c.id)">
                    <!-- 4-Phase Stepper Container -->
                    <div class="booking-stepper-container">
                      <div class="stepper-track">
                        <template v-for="(step, idx) in contractStepLabels" :key="step">
                          <div class="stepper-step" :class="{ completed: isContractStepCompleted(c, idx), active: isContractStepActive(c, idx) }">
                            <div class="step-circle">
                              <svg v-if="isContractStepCompleted(c, idx)" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"></polyline></svg>
                              <span v-else>{{ idx + 1 }}</span>
                            </div>
                            <div class="step-label-text">
                              {{ idx === 0 && c.status === 'revision_required' ? (isTenant ? 'แก้ไขสัญญา' : 'รอผู้เช่าแก้ไข') : step }}
                            </div>
                          </div>
                          <div v-if="idx < contractStepLabels.length - 1" class="stepper-connector" :class="{ filled: isContractStepCompleted(c, idx + 1) || isContractStepCompleted(c, idx) }"></div>
                        </template>
                      </div>
                    </div>

                    <!-- 4-Column Info Grid -->
                    <div class="info-columns-grid officer-columns-grid">
                      <!-- Col 1: ระยะเวลาสัญญา -->
                      <div class="info-column">
                        <div class="info-label">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                          <span>ระยะเวลาสัญญา</span>
                        </div>
                        <div class="info-value-row">
                          <span class="info-value-text">{{ formatContractDate(c.start_date) }} – {{ formatContractDate(c.end_date) }}</span>
                        </div>
                      </div>

                      <!-- Col 2: อัตราค่าเช่า / การชำระ -->
                      <div class="info-column">
                        <div class="info-label">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                          <span>อัตราค่าเช่า</span>
                        </div>
                        <div class="info-value-row">
                          <span class="info-value-text">฿{{ Number(c.rental_fee || 0).toLocaleString('th-TH') }}</span>
                          <span class="days-remaining-pill">{{ c.billing_cycle === 'daily' ? 'ต่อวัน' : 'ต่อเดือน' }}</span>
                        </div>
                        <div class="cell-sub-meeting-point">
                          <span>เงินประกัน ฿{{ Number(c.security_deposit || 0).toLocaleString('th-TH') }}</span>
                        </div>
                      </div>

                      <!-- Col 3: แผงค้า & ร้านค้า -->
                      <div class="info-column">
                        <div class="info-label">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                          <span>แผงค้า / ร้านค้า</span>
                        </div>
                        <div class="info-value-row">
                          <span class="stall-count-pill">
                            <span class="stall-code-strong">{{ c.stall_ident || '-' }}</span>
                          </span>
                        </div>
                        <div v-if="c.shop_name" class="stall-shop-row">
                          <span class="shop-prefix">ชื่อร้าน:</span>
                          <strong class="shop-name-highlight">{{ c.shop_name }}</strong>
                        </div>
                        <div v-if="c.business_type" class="shop-product">ประเภทสินค้า: {{ c.business_type }}</div>
                      </div>

                      <!-- Col 4: ผู้เช่า -->
                      <div class="info-column">
                        <div class="info-label">
                          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                          <span>ผู้เช่า</span>
                        </div>
                        <div class="info-value-row">
                          <span class="info-value-text">{{ c.tenant_name || '-' }}</span>
                        </div>
                        <div v-if="c.tenant_id_card" class="applicant-id-doc-row">
                          <span class="applicant-id-mono">
                            <span class="sub-label">เลขบัตร:</span> {{ c.tenant_id_card }}
                          </span>
                        </div>
                        <div v-if="c.tenant_phone" class="applicant-phone-row">
                          <a :href="'tel:' + c.tenant_phone" class="phone-call-btn" title="คลิกเพื่อโทรออกทันที">
                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                            <span>{{ c.tenant_phone }}</span>
                          </a>
                        </div>
                      </div>
                    </div>

                    <!-- ข้อมูลเพิ่มเติมและรายละเอียดสัญญา -->
                    <div v-if="c.terms_and_conditions || c.review_reason" class="application-inline-meta-bar">
                      <div v-if="c.review_reason" class="inline-details-row text-danger">
                        <span class="details-lead-label">เหตุผลที่ส่งกลับ:</span>
                        <span class="details-lead-text">{{ c.review_reason }}</span>
                      </div>
                      <div v-else-if="c.terms_and_conditions" class="inline-details-row">
                        <span class="details-lead-label">หมายเหตุ/เงื่อนไข:</span>
                        <span class="details-lead-text">{{ c.terms_and_conditions }}</span>
                      </div>
                    </div>

                    <!-- Card Footer Row -->
                    <div class="card-footer-row">
                      <div class="footer-left-status" :class="getContractProgressTone(c.status)">
                        <div class="collapsed-next-step-content">
                          <span class="collapsed-next-step-label">ความคืบหน้า</span>
                          <strong>{{ contractProgressText(c) }}</strong>
                        </div>
                      </div>

                      <div class="footer-right-actions">
                        <AppButton v-if="canReviewContract(c) && c.status !== 'pending_bureau_director'" variant="outline" size="md" @click.stop="rejectContract(c)">
                          <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 4v6h6"/><path d="M3.5 15a9 9 0 1 0 2-9.4L1 10"/></svg></template>
                          ส่งกลับแก้ไข
                        </AppButton>
                        <AppButton v-if="isStaff && c.status === 'draft'" variant="primary" size="md" @click.stop="openContractDraft(c)">กำหนดรายละเอียด</AppButton>
                        <AppButton v-if="isTenant && c.status === 'revision_required'" variant="primary" size="md" @click.stop="openContractRevision(c)">
                          <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg></template>
                          แก้ไขและยื่นใหม่
                        </AppButton>
                        <AppButton variant="secondary" size="md" @click.stop="viewContractDoc(c)" title="ดูรายละเอียดสัญญา">
                          <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"></circle><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8z"></path></svg></template>
                          รายละเอียด
                        </AppButton>
                        <AppButton v-if="canReviewContract(c)" variant="primary" size="md" @click.stop="approveContract(c)">
                          <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg></template>
                          {{ c.status === 'pending_bureau_director' ? 'ลงนามและอนุมัติ' : 'ลงนามและส่งต่อ' }}
                        </AppButton>
                        <AppButton v-if="canTerminateContract(c)" variant="danger-outline" size="md" @click.stop="terminateContract(c)">
                          <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg></template>
                          สิ้นสุดสัญญา
                        </AppButton>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Universal Pagination Footer -->
        <AppPagination
          v-if="filteredContracts.length > 0"
          v-model="currentPage"
          v-model:items-per-page="pageSize"
          :total-items="filteredContracts.length"
          unit-label="สัญญา"
        />
      </div>
    </div>

    <!-- Modal: Create Contract -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal-content" style="max-width: 720px;">
        <AppModalHeader
          :title="isTenant ? 'จัดทำสัญญาเช่าพื้นที่แผงค้า' : 'สร้างสัญญาเช่าพื้นที่ใหม่'"
          :subtitle="isTenant ? 'ตรวจสอบข้อมูลและส่งสัญญาให้เจ้าหน้าที่พิจารณา' : 'บันทึกข้อมูลผู้เช่า แผงค้า และเงื่อนไขสัญญา'"
          icon="edit"
          @close="showCreateModal = false"
        />

        <form novalidate @submit.prevent="submitContract">
          <div class="modal-body">
            <div v-if="editingContract?.status === 'draft'" class="contract-template-source">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg>
              <div>
                <strong>{{ editingContract.contract_template_title || 'รูปแบบสัญญามาตรฐาน' }}</strong>
                <span v-if="editingContract.contract_template_code">รหัส {{ editingContract.contract_template_code }}</span>
                <p>ระบบเติมข้อมูลจากคำร้องให้แล้ว กรุณากรอกเฉพาะช่องที่ยังว่าง</p>
              </div>
            </div>
            <p v-if="editingContract?.review_reason" class="contract-return-note">เหตุผลที่ส่งกลับ: {{ editingContract.review_reason }}</p>
            <AppFileUpload v-if="editingContract?.status !== 'draft'" v-model="contractDocument" label="เอกสารประกอบสัญญา" required-text="จำเป็น" accept=".pdf,.jpg,.jpeg,.png,application/pdf,image/jpeg,image/png" :max-size-mb="3" hint="PDF, JPG หรือ PNG ขนาดไม่เกิน 3MB" :error="contractDocumentError" @error="contractDocumentError = $event" />
            <!-- Notice banner for tenant -->
            <div v-if="isTenant" class="tenant-contract-banner">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line></svg>
              <div>
                <strong style="color: var(--tx-1e40af); font-size: var(--font-size-base);">คำร้องขอเช่าแผงค้าได้รับการอนุมัติเรียบร้อยแล้ว</strong>
                <p style="margin: 2px 0 0; color: var(--tx-3b82f6); font-size: var(--font-size-xs);">กรุณาตรวจสอบข้อมูลสัญญาเช่าด้านล่าง เมื่อกดยืนยันแล้ว สัญญาจะถูกส่งไปยังเจ้าหน้าที่เพื่อตรวจสอบและอนุมัติสัญญา</p>
              </div>
            </div>

            <h4 style="font-size: var(--font-size-base); margin-bottom: 12px; color: var(--primary);">1. ข้อมูลแผงเช่าและร้านค้า</h4>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_stall_id">เลือกแผงเช่าที่ต้องการเช่า <span class="req-star" aria-hidden="true">*</span></label>
                <AppInput
                  v-if="isTenant"
                  id="c_stall_id"
                  :model-value="getTenantStallLabel()"
                  disabled
                />
                <AppSelect
                  v-else
                  v-model="formContract.stall_id"
                  :options="availableStallOptions"
                  placeholder="-- เลือกแผงเช่าที่ว่าง --"
                  @change="onStallSelected"
                />
                <p v-if="contractErrors.stall_id" class="field-error-text">{{ contractErrors.stall_id }}</p>
              </div>

              <div class="form-group">
                <label class="form-label" for="c_business_name">ชื่อร้านค้า (Shop/Business Name) <span class="req-star" aria-hidden="true">*</span></label>
                <AppInput id="c_business_name" v-model="formContract.business_name" type="text" placeholder="เช่น ส้มตำรสเด็ด" :error="contractErrors.business_name" @blur="validateContractField('business_name')" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="c_business_type">ประเภทสินค้าที่จำหน่าย</label>
              <AppInput id="c_business_type" v-model="formContract.business_type" type="text" placeholder="เช่น อาหารปรุงสุก, เครื่องดื่ม, เสื้อผ้า" />
            </div>

            <h4 style="font-size: var(--font-size-base); margin: 20px 0 12px; color: var(--primary);">2. ข้อมูลผู้เช่า (Tenant Information)</h4>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_tenant_name">ชื่อ-นามสกุล ผู้เช่า <span class="req-star" aria-hidden="true">*</span></label>
                <AppInput id="c_tenant_name" v-model="formContract.tenant_name" type="text" placeholder="นาย/นาง/น.ส. ..." :error="contractErrors.tenant_name" @blur="validateContractField('tenant_name')" />
              </div>

              <div class="form-group">
                <label class="form-label" for="c_tenant_id_card">เลขประจำตัวประชาชน / Tax ID <span class="req-star" aria-hidden="true">*</span></label>
                <AppInput id="c_tenant_id_card" v-model="formContract.tenant_id_card" type="text" placeholder="13 หลัก" :error="contractErrors.tenant_id_card" @blur="validateContractField('tenant_id_card')" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_tenant_phone">เบอร์โทรศัพท์ติดต่อ <span class="req-star" aria-hidden="true">*</span></label>
                <AppInput id="c_tenant_phone" v-model="formContract.tenant_phone" type="text" placeholder="08x-xxx-xxxx" :error="contractErrors.tenant_phone" @blur="validateContractField('tenant_phone')" />
              </div>

              <div class="form-group">
                <label class="form-label" for="c_tenant_email">อีเมล (ถ้ามี)</label>
                <AppInput id="c_tenant_email" v-model="formContract.tenant_email" type="email" placeholder="tenant@example.com" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="c_tenant_address">ที่อยู่ตามบัตรประชาชน / ที่อยู่ติดต่อ</label>
              <AppInput id="c_tenant_address" v-model="formContract.tenant_address" type="textarea" :rows="2" placeholder="บ้านเลขที่ หมู่ ซอย ถนน ตำบล อำเภอ จังหวัด" />
            </div>

            <h4 style="font-size: var(--font-size-base); margin: 20px 0 12px; color: var(--primary);">3. เงื่อนไขและอัตราค่าเช่า (Lease Terms)</h4>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_start_date">วันที่เริ่มต้นสัญญา <span class="req-star" aria-hidden="true">*</span></label>
                <AppInput id="c_start_date" v-model="formContract.start_date" type="date" :error="contractErrors.start_date" @blur="validateContractField('start_date')" />
              </div>

              <div class="form-group">
                <label class="form-label" for="c_end_date">วันที่สิ้นสุดสัญญา <span class="req-star" aria-hidden="true">*</span></label>
                <AppInput id="c_end_date" v-model="formContract.end_date" type="date" :error="contractErrors.end_date" @blur="validateContractField('end_date')" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_billing_cycle">รอบการชำระเงิน</label>
                <AppSelect v-model="formContract.billing_cycle" :options="billingCycleOptions" />
              </div>

              <div class="form-group">
                <label class="form-label" for="c_rental_fee">ค่าเช่าต่องวด (บาท) <span class="req-star" aria-hidden="true">*</span></label>
                <AppInput id="c_rental_fee" v-model="formContract.rental_fee" type="number" :error="contractErrors.rental_fee" @blur="validateContractField('rental_fee')" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_security_deposit">เงินประกันความเสียหาย (บาท)</label>
                <AppInput id="c_security_deposit" v-model="formContract.security_deposit" type="number" />
              </div>

              <div class="form-group">
                <label class="form-label" for="c_advance_payment">เงินจ่ายล่วงหน้า (บาท)</label>
                <AppInput id="c_advance_payment" v-model="formContract.advance_payment" type="number" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="c_terms">ข้อกำหนดพิเศษ / เงื่อนไขสัญญา</label>
              <AppInput id="c_terms" v-model="formContract.terms_and_conditions" type="textarea" :rows="2" placeholder="เช่น ห้ามโอนสิทธิ์ให้บุคคลอื่น, ห้ามดัดแปลงโครงสร้าง" />
            </div>
          </div>

          <div class="modal-footer">
            <AppButton type="button" variant="secondary" size="md" @click="showCreateModal = false">ยกเลิก</AppButton>
            <AppButton type="submit" variant="primary" size="md" :disabled="saving" :class="{ 'draft-submit-label': editingContract?.status === 'draft' }">
              <span v-if="!saving && editingContract?.status === 'draft'">บันทึกแบบร่างและส่งให้ผู้เช่าอ่าน</span>
              <span v-if="saving">กำลังส่งสัญญา...</span>
              <span v-else-if="isTenant">ยืนยันและส่งสัญญาให้เจ้าหน้าที่ตรวจสอบ</span>
              <span v-else>ยืนยันและเปิดใช้งานสัญญา</span>
            </AppButton>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal: View Contract Document (Print-ready) -->
    <div v-if="viewDocModal && currentDoc" class="modal-overlay" @click.self="viewDocModal = false">
      <div class="modal-content" style="max-width: 800px; background: var(--bg-ffffff);">
        <AppModalHeader title="เอกสารสัญญาเช่าพื้นที่แผงค้า" :subtitle="`เลขที่สัญญา ${currentDoc.contract_number}`" icon="document" @close="viewDocModal = false">
          <template #actions>
            <AppButton type="button" variant="primary" size="md" :disabled="contractSignaturesLoading || contractSignaturesError || downloadingContract" :loading="downloadingContract" @click="downloadContractPdf">
              <template #prefix><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 9V2h12v7M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect x="6" y="14" width="12" height="8"/></svg></template>
              ดาวน์โหลด PDF
            </AppButton>
          </template>
        </AppModalHeader>

        <div class="modal-body" id="print-area">
          <!-- Officer Review Banner inside Contract Document -->
          <div v-if="canReviewContract(currentDoc) || currentDoc.status === 'revision_required'" class="contract-review-card">
            <div class="review-left">
              <div class="review-badge">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                <span>{{ getStatusLabel(currentDoc.status) }}</span>
              </div>
              <div class="review-desc">{{ currentDoc.status === 'revision_required' ? 'เหตุผล: ' + (currentDoc.review_reason || '-') : 'ตรวจสอบเอกสารและลงนามก่อนส่งต่อ' }}</div>
            </div>
            <div v-if="canReviewContract(currentDoc)" class="review-right">
              <AppButton v-if="currentDoc.status !== 'pending_bureau_director'" type="button" variant="danger-outline" size="md" @click="rejectContract(currentDoc)">
                ส่งกลับแก้ไข
              </AppButton>
              <AppButton type="button" variant="primary" size="md" @click="approveContract(currentDoc)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span>{{ currentDoc.status === 'pending_bureau_director' ? 'ลงนามและอนุมัติ' : 'ลงนามและส่งต่อ' }}</span>
              </AppButton>
            </div>
          </div>

          <AppButton v-if="currentDoc.has_document" variant="outline" size="sm" @click="openContractAttachment(currentDoc)">ดูเอกสารประกอบ</AppButton>
          <div ref="contractPaper" class="contract-document">
            <div class="contract-header">
              <h2 class="contract-doc-title">{{ currentDoc.contract_template_title || 'หนังสือสัญญาเช่าพื้นที่แผงค้าในตลาด' }}</h2>
              <div class="contract-doc-org">ระบบจัดการตลาด e-Market Management System</div>
              <div class="contract-doc-number-pill">
                เลขที่สัญญา: {{ currentDoc.contract_number }}
              </div>
              <div v-if="currentDoc.contract_template_code" class="contract-template-code">รูปแบบ {{ currentDoc.contract_template_code }}</div>
            </div>

            <div class="contract-body">
              <p class="contract-intro-paragraph">
                สัญญาฉบับนี้ทำขึ้น ณ <strong>ตลาดกลางชุมชน e-Market</strong> เมื่อวันที่ <strong>{{ currentDoc.created_at ? currentDoc.created_at.slice(0, 10) : currentDoc.start_date }}</strong>
                ระหว่าง <strong>บริษัท บริหารจัดการตลาด จำกัด</strong> (ซึ่งต่อไปนี้เรียกว่า "ผู้ให้เช่า") ฝ่ายหนึ่ง กับ
              </p>

              <div class="contract-info-card">
                <div class="contract-info-row">
                  <span class="info-label">ชื่อผู้เช่า:</span>
                  <span class="info-value">{{ currentDoc.tenant_name }}</span>
                </div>
                <div class="contract-info-row">
                  <span class="info-label">เลขประจำตัวประชาชน:</span>
                  <span class="info-value">{{ currentDoc.tenant_id_card }}</span>
                </div>
                <div class="contract-info-row">
                  <span class="info-label">เบอร์โทรศัพท์:</span>
                  <span class="info-value">{{ currentDoc.tenant_phone }}</span>
                </div>
                <div class="contract-info-row">
                  <span class="info-label">อีเมล:</span>
                  <span class="info-value">{{ currentDoc.tenant_email || '-' }}</span>
                </div>
                <div class="contract-info-row contract-info-full">
                  <span class="info-label">ที่อยู่:</span>
                  <span class="info-value">{{ currentDoc.tenant_address || '-' }}</span>
                </div>
              </div>

              <div class="contract-clause-heading">ข้อ ๑. วัตถุประสงค์และสถานที่เช่า</div>
              <p>
                ผู้ให้เช่าตกลงให้เช่า และผู้เช่าตกลงเช่าพื้นที่แผงค้าเลขที่ <strong>{{ currentDoc.stall ? currentDoc.stall.stall_number : (currentDoc.stall_ident || currentDoc.stall_id) }}</strong> 
                เพื่อประกอบกิจการร้านค้า <strong>"{{ currentDoc.business_name || currentDoc.shop_name || 'ร้านค้า' }}"</strong> จำหน่ายสินค้าประเภท <strong>{{ currentDoc.business_type || 'ทั่วไป' }}</strong>
              </p>

              <div class="contract-clause-heading">ข้อ ๒. ระยะเวลาการเช่า</div>
              <p>
                สัญญาการเช่ามีกำหนดระยะเวลา ตั้งแต่วันที่ <strong>{{ currentDoc.start_date }}</strong> ถึงวันที่ <strong>{{ currentDoc.end_date }}</strong>
              </p>

              <div class="contract-clause-heading">ข้อ ๓. ค่าเช่าและเงินประกัน</div>
              <div class="contract-info-card">
                <div class="contract-info-row">
                  <span class="info-label">อัตราค่าเช่า:</span>
                  <span class="info-value">฿{{ Number(currentDoc.rental_fee).toLocaleString() }} ({{ currentDoc.billing_cycle === 'monthly' ? 'รายเดือน' : 'รายวัน' }})</span>
                </div>
                <div class="contract-info-row">
                  <span class="info-label">เงินประกันสัญญา:</span>
                  <span class="info-value">฿{{ Number(currentDoc.security_deposit || 0).toLocaleString() }} บาท</span>
                </div>
                <div class="contract-info-row">
                  <span class="info-label">เงินจ่ายล่วงหน้า:</span>
                  <span class="info-value">฿{{ Number(currentDoc.advance_payment || 0).toLocaleString() }} บาท</span>
                </div>
                <div class="contract-info-row">
                  <span class="info-label">สถานะสัญญา:</span>
                  <span class="info-value">{{ getStatusLabel(currentDoc.status) }}</span>
                </div>
              </div>

              <div class="contract-clause-heading">ข้อ ๔. ข้อกำหนดและเงื่อนไขเพิ่มเติม</div>
              <p class="contract-terms-text">
                {{ currentDoc.terms_and_conditions || 'ผู้เช่าต้องปฏิบัติตามกฎระเบียบความสะอาด สุขอนามัย และความปลอดภัยของตลาดอย่างเคร่งครัด' }}
              </p>

              <div class="contract-sign-box">
                <div v-for="signer in contractSignatures" :key="signer.kind" class="sign-line">
                  <img v-if="signer.url" :src="signer.url" :alt="'ลายเซ็น ' + signer.label" class="contract-signature-image" />
                  <span v-else class="dotted-signature-line">...........................................................</span>
                  <div class="signer-name">( {{ signer.label }} )</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <ContractTemplatePreviewModal :template="draftPreviewTemplate" :contract="draftPreviewContract" @close="draftPreviewTemplate = null; draftPreviewContract = null" />
  <OfficerSignatureModal :open="showContractSignature" @cancel="showContractSignature = false" @confirm="onContractSignature" />
  <ApplicationMockOtpModal :open="showContractOtp" @cancel="showContractOtp = false; pendingContractSignature = ''" @confirm="confirmContractOtp" />
  <AppModal v-model="showContractReturn" title="ส่งสัญญากลับแก้ไข" description="ระบุสิ่งที่ผู้เช่าต้องแก้ไข">
    <label class="form-label" for="contract-return-reason">เหตุผล <span class="req-star" aria-hidden="true">*</span></label>
    <AppInput id="contract-return-reason" v-model="contractReturnReason" type="textarea" :error="contractReturnError" />
    <template #footer><AppButton variant="secondary" size="md" @click="showContractReturn = false">ยกเลิก</AppButton><AppButton variant="primary" size="md" @click="submitContractReturn">ส่งกลับแก้ไข</AppButton></template>
  </AppModal>
  <AppModal v-model="showContractAttachment" title="เอกสารประกอบสัญญา" description="เอกสารที่ผู้เช่าแนบ">
    <iframe v-if="attachmentType === 'application/pdf'" :src="attachmentUrl" title="เอกสารประกอบสัญญา" style="width:100%;height:60vh;border:0"></iframe>
    <img v-else-if="attachmentUrl" :src="attachmentUrl" alt="เอกสารประกอบสัญญา" style="max-width:100%" />
    <template #footer><AppButton variant="secondary" size="md" @click="showContractAttachment = false">ปิด</AppButton></template>
  </AppModal>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'
import AppSegmentedControl from '~/components/AppSegmentedControl.vue'
import AppPagination from '~/components/AppPagination.vue'
import AppButton from '~/components/AppButton.vue'
import AppSelect from '~/components/AppSelect.vue'
import AppInput from '~/components/AppInput.vue'
import AppFileUpload from '~/components/AppFileUpload.vue'
import OfficerSignatureModal from '~/components/OfficerSignatureModal.vue'
import ApplicationMockOtpModal from '~/components/ApplicationMockOtpModal.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { get, post, put, apiFetch } = useApi()
const { showSuccess, showError, showConfirm } = useAppAlert()

const isTenant = computed(() => authStore.isTenant || authStore.currentUser?.role === 'tenant')
const isStaff = computed(() => authStore.role === 'staff')
const canReviewContract = (c: any) => (c?.status === 'pending' && authStore.role === 'staff') || (c?.status === 'pending_market_director' && authStore.role === 'market_director') || (c?.status === 'pending_bureau_director' && authStore.role === 'bureau_director')

const loading = ref(true)
const saving = ref(false)
const contracts = ref<any[]>([])
const filteredContracts = ref<any[]>([])
const availableStalls = ref<any[]>([])
const allContractStalls = ref<any[]>([])

const selectedStatus = ref('')
const pendingContractFilter = ref('all')
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

// View Mode: Table vs Cards
const viewMode = ref<'card' | 'table'>('card')
const viewModeOptions = [
  { value: 'card', label: 'ตลาด', title: 'มุมมองแยกกลุ่มตามตลาด', icon: 'market' as const },
  { value: 'table', label: 'ตาราง', title: 'มุมมองตารางรายการสัญญา', icon: 'table' as const }
]

const pendingContractOptions = computed(() => {
  const options = [
    { value: 'all', label: 'ทั้งหมด' },
    { value: 'pending', label: 'รอเจ้าหน้าที่' },
    { value: 'pending_market_director', label: 'รอ ผอ.ตลาด' },
    { value: 'pending_bureau_director', label: 'รอ ผอ.สำนัก' },
    { value: 'revision_required', label: 'ส่งกลับแก้ไข' }
  ]
  const pendingStates = ['pending', 'pending_market_director', 'pending_bureau_director', 'revision_required']
  return options.map(option => ({
    ...option,
    count: contracts.value.filter(contract => option.value === 'all' ? pendingStates.includes(contract.status) : contract.status === option.value).length
  }))
})
const selectedContractStatusLabel = computed(() => {
  if (selectedStatus.value !== 'pending') return getStatusLabel(selectedStatus.value)
  return pendingContractFilter.value === 'all'
    ? 'รอดำเนินการ'
    : pendingContractOptions.value.find(option => option.value === pendingContractFilter.value)?.label || 'รอดำเนินการ'
})
const countContractsByStatus = (st: string) => {
  return contracts.value.filter(c => st === 'pending' ? ['pending', 'pending_market_director', 'pending_bureau_director', 'revision_required'].includes(c.status) : c.status === st).length
}

const paginatedContracts = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredContracts.value.slice(start, start + pageSize.value)
})


const getContractBorderClass = (status: string) => {
  if (status === 'active') return 'border-status-approved'
  if (['pending', 'pending_market_director', 'pending_bureau_director'].includes(status)) return 'border-status-pending'
  if (status === 'revision_required') return 'border-status-pending'
  if (['expired', 'terminated', 'rejected'].includes(status)) return 'border-status-muted'
  return 'border-status-ready'
}

const getContractProgressTone = (status: string) => {
  if (status === 'active') return 'tone-success'
  if (['pending', 'pending_market_director', 'pending_bureau_director', 'revision_required'].includes(status)) return 'tone-waiting'
  if (['expired', 'terminated', 'rejected'].includes(status)) return 'tone-muted'
  return 'tone-action'
}

const isContractStepCompleted = (contract: any, index: number) => {
  const completed = [
    !!contract.has_tenant_signature || contract.status !== 'revision_required',
    !!contract.has_officer_signature || ['pending_market_director', 'pending_bureau_director', 'active', 'expired', 'terminated'].includes(contract.status),
    !!contract.has_market_signature || ['pending_bureau_director', 'active', 'expired', 'terminated'].includes(contract.status),
    !!contract.has_bureau_signature || ['active', 'expired', 'terminated'].includes(contract.status)
  ]
  return completed[index]
}

const isContractStepActive = (contract: any, index: number) => {
  if (isContractStepCompleted(contract, index)) return false
  const activeStep: Record<string, number> = {
    revision_required: 0,
    pending: 1,
    pending_market_director: 2,
    pending_bureau_director: 3
  }
  return activeStep[contract.status] === index
}

const getPendingChipClass = (val: string) => {
  switch (val) {
    case 'all': return 'chip-all'
    case 'pending': return 'chip-amber'
    case 'pending_market_director': return 'chip-purple'
    case 'pending_bureau_director': return 'chip-green'
    case 'revision_required': return 'chip-orange'
    default: return 'chip-all'
  }
}

const contractStepLabels = ['ผู้เช่าส่งสัญญา', 'เจ้าหน้าที่ตรวจ', 'ผอ.ตลาดพิจารณา', 'ผอ.สำนักอนุมัติ']
const openContractGroups = ref<Record<string, boolean>>({})
const openContractCards = ref<Record<string, boolean>>({})
const isContractGroupOpen = (key: string) => openContractGroups.value[key] !== false
const toggleContractGroup = (key: string) => { openContractGroups.value[key] = !isContractGroupOpen(key) }
const isContractCardOpen = (id: string | number) => openContractCards.value[String(id)] !== false
const toggleContractCard = (id: string | number) => { openContractCards.value[String(id)] = !isContractCardOpen(id) }
const canTerminateContract = (c: any) => c?.status === 'active' && ['staff', 'market_director'].includes(authStore.role)
const stallForContract = (c: any) => allContractStalls.value.find((stall: any) =>
  (c.stall_id && String(stall.id) === String(c.stall_id)) ||
  (c.stall_ident && String(stall.ident || stall.stall_number) === String(c.stall_ident))
)
const contractMarket = (c: any) => c.market_name || stallForContract(c)?.market_name || 'ไม่ระบุตลาด'
const contractProject = (c: any) => c.project_name || stallForContract(c)?.project_name || stallForContract(c)?.zone_name || 'ไม่ระบุโครงการ'
const groupedContracts = computed(() => {
  const groups = new Map<string, { key: string; name: string; items: any[] }>()
  for (const contract of paginatedContracts.value) {
    const market = contractMarket(contract)
    const project = contractProject(contract)
    const key = market + '|' + project
    if (!groups.has(key)) groups.set(key, { key, name: project === 'ไม่ระบุโครงการ' ? market : market + ' · ' + project, items: [] })
    groups.get(key)!.items.push(contract)
  }
  return Array.from(groups.values())
})
const formatContractDate = (value: unknown) => {
  if (!value) return '-'
  const date = new Date(String(value))
  if (Number.isNaN(date.getTime())) return String(value)
  return new Intl.DateTimeFormat('th-TH', { day: 'numeric', month: 'short', year: 'numeric', timeZone: 'Asia/Bangkok' }).format(date)
}
const contractBorderClass = (status: string) => {
  if (status === 'active') return 'border-active'
  if (status === 'revision_required' || status === 'rejected' || status === 'terminated') return 'border-danger'
  if (status === 'expired') return 'border-muted'
  return 'border-pending'
}
const contractStepClass = (contract: any, index: number) => {
  const state = String(contract.status || '')
  const completed = [
    !!contract.has_tenant_signature,
    !!contract.has_officer_signature,
    !!contract.has_market_signature,
    !!contract.has_bureau_signature
  ]
  if (state === 'revision_required') return index === 0 ? 'current' : 'waiting'
  if (completed[index]) return 'done'
  const activeStep: Record<string, number> = { revision_required: 0, pending: 1, pending_market_director: 2, pending_bureau_director: 3 }
  return activeStep[state] === index ? 'current' : 'waiting'
}
const contractProgressText = (contract: any) => {
  switch (contract.status) {
	case 'draft': return 'รอเจ้าหน้าที่กำหนดรายละเอียด'
	case 'tenant_review': return 'รอผู้เช่าตรวจสอบแบบร่าง'
    case 'pending': return 'รอเจ้าหน้าที่ตรวจเอกสาร'
    case 'pending_market_director': return 'รอ ผอ.ตลาดพิจารณา'
    case 'pending_bureau_director': return 'รอ ผอ.สำนักลงนาม'
    case 'revision_required': return 'ส่งกลับให้ผู้เช่าแก้ไข'
    case 'active': return 'สัญญาอนุมัติแล้ว'
    case 'expired': return 'สัญญาหมดอายุ'
    case 'terminated': return 'สัญญาสิ้นสุดแล้ว'
    default: return getStatusLabel(contract.status)
  }
}
const exportReport = async () => {
  await showSuccess(`ส่งออกรายงานสัญญาเช่าพื้นที่ ${filteredContracts.value.length} ฉบับ เรียบร้อยแล้ว`, 'ส่งออกรายงาน')
}

// Create Modal
const showCreateModal = ref(false)
const contractDocument = ref<File | null>(null)
const contractDocumentError = ref('')
const editingContract = ref<any>(null)
const showContractSignature = ref(false)
const signaturePurpose = ref<'submit' | 'review'>('submit')
const reviewTarget = ref<any>(null)
const showContractOtp = ref(false)
const pendingContractSignature = ref('')
const showContractReturn = ref(false)
const contractReturnReason = ref('')
const contractReturnError = ref('')
const showContractAttachment = ref(false)
const attachmentUrl = ref('')
const attachmentType = ref('')
const contractSignatures = ref([{ kind: 'tenant', label: 'ผู้เช่า', url: '' }, { kind: 'officer', label: 'เจ้าหน้าที่', url: '' }, { kind: 'market', label: 'ผอ.ตลาด', url: '' }, { kind: 'bureau', label: 'ผอ.สำนักตลาด', url: '' }])
const formContract = ref({
  stall_id: '' as any,
  tenant_name: '',
  tenant_id_card: '',
  tenant_phone: '',
  tenant_email: '',
  tenant_address: '',
  business_name: '',
  business_type: '',
  start_date: new Date().toISOString().slice(0, 10),
  end_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString().slice(0, 10),
  billing_cycle: 'monthly',
  rental_fee: 4000,
  security_deposit: 8000,
  advance_payment: 4000,
  terms_and_conditions: 'ผู้เช่าต้องรักษาความสะอาดและชำระค่าเช่าภายในวันที่ 5 ของทุกเดือน',
  status: 'active'
})

const getTenantStallLabel = () => {
  const s = availableStalls.value.find(s => String(s.id) === String(formContract.value.stall_id))
  if (!s) return formContract.value.stall_id ? `แผง ${formContract.value.stall_id}` : 'แผงค้าที่เลือก'
  return `แผง ${s.stall_number || s.ident} (${s.project_name || s.zone_name || s.zone?.name || `โซน ${s.zone_id || 'A'}`}) - ฿${Number(s.monthly_rate || s.monthly_rent || formContract.value.rental_fee || 2000).toLocaleString()}/ด.`
}

// Document Preview Modal
const viewDocModal = ref(false)
const currentDoc = ref<any>(null)
const draftPreviewTemplate = ref<any>(null)
const draftPreviewContract = ref<any>(null)
const contractPaper = ref<HTMLElement | null>(null)
const downloadingContract = ref(false)
const contractSignaturesLoading = ref(false)
const contractSignaturesError = ref(false)
useBodyScrollLock(computed(() => showCreateModal.value || viewDocModal.value))

const loadData = async () => {
  loading.value = true
  try {
    if (route.query.status) {
      selectedStatus.value = String(route.query.status)
    }
    const [contractsRes, stallsRes, appsRes] = await Promise.all([
      get('/contracts'),
      get('/stalls'),
      get('/applications').catch(() => [])
    ])
    const rawContracts: any[] = Array.isArray(contractsRes) ? contractsRes : (Array.isArray(contractsRes?.data) ? contractsRes.data : [])
    // API ส่งค่าเช่าเป็น monthly_rental (ไม่มี rental_fee / billing_cycle) จึงแปลงให้หน้านี้ใช้ได้ ไม่ให้แสดง ฿NaN
    contracts.value = rawContracts.map((c: any) => ({
      ...c,
      rental_fee: Number(c.rental_fee ?? c.monthly_rental ?? 0),
      billing_cycle: c.billing_cycle || 'monthly'
    }))
    const rawStalls = Array.isArray(stallsRes) ? stallsRes : (Array.isArray(stallsRes?.data) ? stallsRes.data : [])
    const rawApps = Array.isArray(appsRes) ? appsRes : (Array.isArray(appsRes?.data) ? appsRes.data : [])
    allContractStalls.value = rawStalls
    availableStalls.value = rawStalls.filter((s: any) => s.status === 'available' || s.stall_status === 'available')
    if (route.query.stall) searchQuery.value = String(route.query.stall)
    filterContracts()

    // If query contains new_for_stall (from Approved Lease Application)
    if (route.query.new_for_stall) {
      const stallKey = String(route.query.new_for_stall).trim()
      const appId = route.query.app_id ? Number(route.query.app_id) : null
      const matchedApp = rawApps.find((a: any) => appId
        ? Number(a.id) === appId
        : (a.status === 'approved' || a.status === 'completed' || a.status === 'contract_pending') && (String(a.stall_number || a.stall_ident) === stallKey || String(a.stall_id) === stallKey)
      )
      if (!matchedApp || !['approved', 'completed', 'contract_pending'].includes(matchedApp.status)) {
        await showError('ไม่พบคำร้องที่อนุมัติขั้นสุดท้ายแล้ว กรุณากลับไปตรวจสอบหน้าคำร้อง', 'ไม่สามารถสร้างสัญญา')
        return
      }
      if (rawContracts.some((contract: any) => Number(contract.application_id) === Number(matchedApp.id))) {
        await showError('คำร้องนี้มีสัญญาแล้ว กรุณาตรวจสอบรายการสัญญา', 'มีสัญญาแล้ว')
        return
      }
      let target = rawStalls.find((s: any) =>
        String(s.id) === stallKey ||
        String(s.stall_number || s.ident) === stallKey
      )
      if (!target && matchedApp) {
        target = {
          id: matchedApp.stall_id || matchedApp.stall_number,
          stall_number: matchedApp.stall_number || matchedApp.stall_ident,
          ident: matchedApp.stall_number || matchedApp.stall_ident,
          zone_id: 'A',
          monthly_rate: matchedApp.stall_monthly_rate || matchedApp.monthly_rate || 3600
        }
      }
      if (target) {
        if (!availableStalls.value.some((s: any) => String(s.id) === String(target.id))) {
          availableStalls.value = [target, ...availableStalls.value]
        }
        openCreateModal()
        linkedApplicationId.value = Number(matchedApp.id)
        formContract.value.stall_id = target.id
        const rate = matchedApp?.stall_monthly_rate || matchedApp?.monthly_rate || target.monthly_rate || target.monthly_rent || 3600
        formContract.value.rental_fee = rate
        formContract.value.security_deposit = rate * 2
        formContract.value.advance_payment = rate
        if (matchedApp) {
          formContract.value.tenant_name = matchedApp.applicant_name || matchedApp.tenant_name || ''
          formContract.value.tenant_id_card = matchedApp.applicant_id_card || ''
          formContract.value.tenant_phone = matchedApp.applicant_phone || matchedApp.tenant_phone || ''
          formContract.value.tenant_email = matchedApp.applicant_email || ''
          formContract.value.tenant_address = matchedApp.applicant_address || ''
          formContract.value.business_name = matchedApp.shop_name || ''
          formContract.value.business_type = matchedApp.product_type || matchedApp.product_category || ''
          if (matchedApp.requested_start_date) {
            formContract.value.start_date = matchedApp.requested_start_date
          }
        }
      }
    }
  } catch (err) {
    console.error('Failed to load contracts:', err)
  } finally {
    loading.value = false
  }
}

const filterContracts = () => {
  let list = Array.isArray(contracts.value) ? [...contracts.value] : []
  if (selectedStatus.value) {
    list = list.filter(c => selectedStatus.value === 'pending' ? (pendingContractFilter.value === 'all' ? ['pending', 'pending_market_director', 'pending_bureau_director', 'revision_required'].includes(c.status) : c.status === pendingContractFilter.value) : c.status === selectedStatus.value)
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter(c =>
      (c.contract_number && c.contract_number.toLowerCase().includes(q)) ||
      (c.tenant_name && c.tenant_name.toLowerCase().includes(q)) ||
      (c.business_name && c.business_name.toLowerCase().includes(q)) ||
      (c.shop_name && c.shop_name.toLowerCase().includes(q)) ||
      (c.tenant_phone && c.tenant_phone.includes(q)) ||
      (c.stall?.stall_number && String(c.stall.stall_number).toLowerCase().includes(q)) ||
      (c.stall_ident && String(c.stall_ident).toLowerCase().includes(q))
    )
  }
  filteredContracts.value = list
}

const availableStallOptions = computed(() =>
  availableStalls.value.map((s: any) => ({
    value: s.id,
    label: `แผง ${s.stall_number || s.ident} (${s.project_name || s.zone_name || s.zone?.name || `โซน ${s.zone_id}`}) - ฿${Number(s.monthly_rate || s.monthly_rent || 1800).toLocaleString()}/ด.`
  }))
)

const billingCycleOptions = [
  { value: 'monthly', label: 'รายเดือน (Monthly)' },
  { value: 'daily', label: 'รายวัน (Daily)' }
]

const contractErrors = reactive({
  stall_id: '',
  business_name: '',
  tenant_name: '',
  tenant_id_card: '',
  tenant_phone: '',
  start_date: '',
  end_date: '',
  rental_fee: ''
})

const resetContractErrors = () => {
  contractErrors.stall_id = ''
  contractErrors.business_name = ''
  contractErrors.tenant_name = ''
  contractErrors.tenant_id_card = ''
  contractErrors.tenant_phone = ''
  contractErrors.start_date = ''
  contractErrors.end_date = ''
  contractErrors.rental_fee = ''
}

const validateContractField = (field: keyof typeof contractErrors) => {
  if (field === 'stall_id') {
    contractErrors.stall_id = formContract.value.stall_id ? '' : 'กรุณาเลือกแผงเช่าที่ต้องการเช่า'
  } else if (field === 'business_name') {
    contractErrors.business_name = formContract.value.business_name.trim() ? '' : 'กรุณาระบุชื่อร้านค้า'
  } else if (field === 'tenant_name') {
    contractErrors.tenant_name = formContract.value.tenant_name.trim() ? '' : 'กรุณาระบุชื่อ-นามสกุลผู้เช่า'
  } else if (field === 'tenant_id_card') {
    contractErrors.tenant_id_card = formContract.value.tenant_id_card.trim() ? '' : 'กรุณาระบุเลขประจำตัวประชาชน / Tax ID'
  } else if (field === 'tenant_phone') {
    contractErrors.tenant_phone = formContract.value.tenant_phone.trim() ? '' : 'กรุณาระบุเบอร์โทรศัพท์ติดต่อ'
  } else if (field === 'start_date') {
    contractErrors.start_date = formContract.value.start_date ? '' : 'กรุณาเลือกวันที่เริ่มต้นสัญญา'
  } else if (field === 'end_date') {
    contractErrors.end_date = formContract.value.end_date ? '' : 'กรุณาเลือกวันที่สิ้นสุดสัญญา'
  } else if (field === 'rental_fee') {
    const val = Number(formContract.value.rental_fee)
    contractErrors.rental_fee = String(formContract.value.rental_fee ?? '').trim() === '' || !Number.isFinite(val) || val <= 0 ? 'กรุณาระบุค่าเช่าต่องวด' : ''
  }
}

const validateContractForm = () => {
  validateContractField('stall_id')
  validateContractField('business_name')
  validateContractField('tenant_name')
  validateContractField('tenant_id_card')
  validateContractField('tenant_phone')
  validateContractField('start_date')
  validateContractField('end_date')
  validateContractField('rental_fee')
  return !Object.values(contractErrors).some(Boolean)
}

const onStallSelected = () => {
  validateContractField('stall_id')
  const stall = availableStalls.value.find(s => s.id === formContract.value.stall_id)
  if (stall) {
    const rate = stall.monthly_rate || stall.monthly_rent || 2000
    formContract.value.rental_fee = rate
    formContract.value.security_deposit = rate * 2
    formContract.value.advance_payment = rate
    formContract.value.business_type = stall.category || stall.zone_name || ''
  }
}

const openCreateModal = () => {
  editingContract.value = null
  contractDocument.value = null
  contractDocumentError.value = ''
  linkedApplicationId.value = null
  resetContractErrors()
  formContract.value = {
    stall_id: availableStalls.value[0]?.id || '',
    tenant_name: '',
    tenant_id_card: '',
    tenant_phone: '',
    tenant_email: '',
    tenant_address: '',
    business_name: '',
    business_type: '',
    start_date: new Date().toISOString().slice(0, 10),
    end_date: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString().slice(0, 10),
    billing_cycle: 'monthly',
    rental_fee: 4000,
    security_deposit: 8000,
    advance_payment: 4000,
    terms_and_conditions: 'ผู้เช่าต้องรักษาความสะอาดและชำระค่าเช่าภายในวันที่ 5 ของทุกเดือน',
    status: 'active'
  }
  if (formContract.value.stall_id) {
    onStallSelected()
  }
  showCreateModal.value = true
}

const linkedApplicationId = ref<number | null>(null)

const readContractFile = (file: File) => new Promise<string>((resolve, reject) => {
  const reader = new FileReader()
  reader.onload = () => resolve(String(reader.result).split(',')[1] || '')
  reader.onerror = () => reject(new Error('อ่านไฟล์ไม่สำเร็จ'))
  reader.readAsDataURL(file)
})
const openContractRevision = (contract: any) => {
  openCreateModal()
  editingContract.value = contract
  linkedApplicationId.value = contract.application_id || null
  const stall = availableStalls.value.find((item: any) => String(item.ident || item.stall_number) === String(contract.stall_ident))
  if (!stall) availableStalls.value.push({ id: contract.stall_id || contract.stall_ident, ident: contract.stall_ident, stall_number: contract.stall_ident })
  formContract.value = {
    ...formContract.value,
    stall_id: stall?.id || contract.stall_id || contract.stall_ident,
    tenant_name: contract.tenant_name || '',
    tenant_id_card: contract.tenant_id_card || '',
    tenant_phone: contract.tenant_phone || '',
    tenant_email: contract.tenant_email || '',
    tenant_address: contract.tenant_address || '',
    business_type: contract.business_type || '',
    billing_cycle: contract.billing_cycle || 'monthly',
    terms_and_conditions: contract.terms_and_conditions || '',
    advance_payment: contract.advance_payment || 0,
    business_name: contract.shop_name || '',
    start_date: contract.start_date || '',
    end_date: contract.end_date || '',
    rental_fee: contract.monthly_rental || contract.rental_fee || 0,
    security_deposit: contract.security_deposit || 0
  }
}
const openContractDraft = (contract: any) => {
  openContractRevision(contract)
}

const saveContractDraft = async () => {
  if (!editingContract.value || saving.value) return
  saving.value = true
  try {
    await put(`/contracts/${editingContract.value.id}/draft`, {
      tenant_name: formContract.value.tenant_name,
      tenant_id_card: formContract.value.tenant_id_card,
      tenant_phone: formContract.value.tenant_phone,
      tenant_email: formContract.value.tenant_email,
      tenant_address: formContract.value.tenant_address,
      business_type: formContract.value.business_type,
      billing_cycle: formContract.value.billing_cycle,
      terms_and_conditions: formContract.value.terms_and_conditions,
      shop_name: formContract.value.business_name,
      start_date: formContract.value.start_date,
      end_date: formContract.value.end_date,
      monthly_rental: Number(formContract.value.rental_fee),
      security_deposit: Number(formContract.value.security_deposit),
      advance_payment: Number(formContract.value.advance_payment)
    })
    showCreateModal.value = false
    selectedStatus.value = ''
    await loadData()
    await showSuccess('บันทึกแบบร่างแล้ว ผู้เช่าสามารถตรวจสอบรายละเอียดได้', 'บันทึกสำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || err?.data?.message || err?.message || 'บันทึกแบบร่างไม่สำเร็จ')
  } finally {
    saving.value = false
  }
}
const submitContract = async () => {
  if (!validateContractForm()) return
	if (editingContract.value?.status === 'draft') {
		await saveContractDraft()
		return
	}
  if (!contractDocument.value) { contractDocumentError.value = 'กรุณาแนบเอกสารประกอบสัญญา'; return }
  if (contractDocumentError.value || saving.value) return
  signaturePurpose.value = 'submit'
  showContractSignature.value = true
}
const onContractSignature = async (signature: string) => {
  showContractSignature.value = false
  if (signaturePurpose.value === 'review') {
    if (reviewTarget.value?.status === 'pending_bureau_director') {
      pendingContractSignature.value = signature
      showContractOtp.value = true
      return
    }
    await sendContractApproval(signature, '')
    return
  }
  if (!contractDocument.value || saving.value) return
  saving.value = true
  try {
    const document = await readContractFile(contractDocument.value)
    if (editingContract.value) {
      await post('/contracts/' + editingContract.value.id + '/resubmit', {
        TenantName: formContract.value.tenant_name,
        TenantPhone: formContract.value.tenant_phone,
        TenantIDCard: formContract.value.tenant_id_card,
        TenantEmail: formContract.value.tenant_email,
        TenantAddress: formContract.value.tenant_address,
        BusinessType: formContract.value.business_type,
        BillingCycle: formContract.value.billing_cycle,
        TermsAndConditions: formContract.value.terms_and_conditions,
        AdvancePayment: Number(formContract.value.advance_payment),
        ShopName: formContract.value.business_name,
        StartDate: formContract.value.start_date,
        EndDate: formContract.value.end_date,
        MonthlyRental: Number(formContract.value.rental_fee),
        SecurityDeposit: Number(formContract.value.security_deposit),
        Signature: signature,
        Document: document
      })
    } else {
      const selectedStall = availableStalls.value.find(s => String(s.id) === String(formContract.value.stall_id))
      const stallIdent = selectedStall?.ident || selectedStall?.stall_number || String(formContract.value.stall_id || '')
      await post('/contracts', {
        ...formContract.value,
        application_id: linkedApplicationId.value || undefined,
        stall_ident: stallIdent,
        stall_number: stallIdent,
        stall_id: Number(formContract.value.stall_id) || 0,
        shop_name: formContract.value.business_name,
        monthly_rental: Number(formContract.value.rental_fee),
        security_deposit: Number(formContract.value.security_deposit),
        advance_payment: Number(formContract.value.advance_payment),
        signature_base64: signature,
        document_base64: document
      })
    }
    showCreateModal.value = false
    if (route.query.new_for_stall) { const query = { ...route.query }; delete query.new_for_stall; delete query.app_id; await router.replace({ query }) }
    selectedStatus.value = ''
    await loadData()
    await showSuccess('ส่งสัญญาให้เจ้าหน้าที่ตรวจสอบแล้ว', 'ส่งสัญญาสำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || err?.data?.message || err?.message || 'ส่งสัญญาไม่สำเร็จ')
  } finally { saving.value = false }
}
const approveContract = (contract: any) => {
  if (!canReviewContract(contract)) return
  reviewTarget.value = contract
  signaturePurpose.value = 'review'
  showContractSignature.value = true
}
const confirmContractOtp = async (otp: string) => {
  showContractOtp.value = false
  await sendContractApproval(pendingContractSignature.value, otp)
  pendingContractSignature.value = ''
}
const sendContractApproval = async (signature: string, otp: string) => {
  const contract = reviewTarget.value
  if (!contract) return
  try {
    await post('/contracts/' + contract.id + '/action', { action: 'approve', signature, otp })
    viewDocModal.value = false
    await loadData()
    await showSuccess('บันทึกลายเซ็นและส่งสัญญาแล้ว', 'ดำเนินการสำเร็จ')
  } catch (err: any) { await showError(err?.data?.message || err?.data?.detail || 'ดำเนินการไม่สำเร็จ') }
}
const rejectContract = (contract: any) => {
  if (!canReviewContract(contract) || contract.status === 'pending_bureau_director') return
  reviewTarget.value = contract
  contractReturnReason.value = ''
  contractReturnError.value = ''
  showContractReturn.value = true
}
const submitContractReturn = async () => {
  if (!contractReturnReason.value.trim()) { contractReturnError.value = 'กรุณาระบุเหตุผล'; return }
  try {
    await post('/contracts/' + reviewTarget.value.id + '/action', { action: 'return_revision', reason: contractReturnReason.value.trim() })
    showContractReturn.value = false
    viewDocModal.value = false
    await loadData()
    await showSuccess('ส่งสัญญากลับให้ผู้เช่าแก้ไขแล้ว')
  } catch (err: any) { await showError(err?.data?.message || err?.data?.detail || 'ส่งกลับไม่สำเร็จ') }
}
const terminateContract = async (contract: any) => {
  const confirmed = await showConfirm(
    `คุณต้องการสิ้นสุดสัญญา ${contract.contract_number} ก่อนกำหนดหรือไม่? ระบบจะเปลี่ยนสถานะแผงเป็น "ว่าง" โดยอัตโนมัติ`,
    {
      title: 'ยืนยันสิ้นสุดสัญญาก่อนกำหนด',
      type: 'error',
      icon: 'trash',
      confirmText: 'สิ้นสุดสัญญา',
      cancelText: 'ยกเลิก'
    }
  )
  if (confirmed) {
    try {
      await post(`/contracts/${contract.id}/terminate`)
      await loadData()
      await showSuccess('สัญญาเช่าถูกยกเลิกแล้ว และสถานะแผงถูกปรับเป็นว่างเรียบร้อย', 'สำเร็จ')
    } catch (err: any) {
      await showError(err?.data?.detail || 'เกิดข้อผิดพลาด', 'ข้อผิดพลาด')
    }
  }
}

let signatureLoadSequence = 0
const clearContractSignatures = () => {
  contractSignatures.value.forEach(item => { if (item.url) URL.revokeObjectURL(item.url); item.url = '' })
}
const viewContractDoc = async (c: any) => {
  if (c?.status === 'draft') {
    draftPreviewTemplate.value = {
      Code: c.contract_template_code || 'standard_contract',
      Title: c.contract_template_title || 'สัญญาเช่าแผงค้า',
      ContractKind: c.contract_template_kind || 'standard',
      Description: c.contract_template_description || c.terms_and_conditions || '',
      IsActive: true,
      SortOrder: 0
    }
    draftPreviewContract.value = c
    return
  }
  draftPreviewTemplate.value = null
  draftPreviewContract.value = null
  const sequence = ++signatureLoadSequence
  clearContractSignatures()
  currentDoc.value = c
  viewDocModal.value = true
  const flags = [c.has_tenant_signature, c.has_officer_signature, c.has_market_signature, c.has_bureau_signature]
  await Promise.all(contractSignatures.value.map(async (item, index) => {
    if (!flags[index]) return
    try {
      const data = await apiFetch<Blob>('/contracts/' + c.id + '/assets/' + item.kind, { responseType: 'blob' })
      if (sequence !== signatureLoadSequence) return
      item.url = URL.createObjectURL(data)
      await new Promise<void>((resolve, reject) => { const img = new Image(); img.onload = () => resolve(); img.onerror = () => reject(new Error('image')); img.src = item.url })
    } catch { if (sequence === signatureLoadSequence) contractSignaturesError.value = true }
  }))
}
const openContractAttachment = async (contract: any) => {
  try {
    const file = await apiFetch<Blob>('/contracts/' + contract.id + '/assets/document', { responseType: 'blob' })
    if (attachmentUrl.value) URL.revokeObjectURL(attachmentUrl.value)
    attachmentType.value = file.type
    attachmentUrl.value = URL.createObjectURL(file)
    showContractAttachment.value = true
  } catch { await showError('โหลดเอกสารประกอบไม่สำเร็จ') }
}
watch(viewDocModal, open => { if (!open) { signatureLoadSequence++; clearContractSignatures() } })
onBeforeUnmount(() => { signatureLoadSequence++; clearContractSignatures(); if (attachmentUrl.value) URL.revokeObjectURL(attachmentUrl.value) })

const downloadContractPdf = async () => {
  if (!contractPaper.value || contractSignaturesLoading.value || contractSignaturesError.value || downloadingContract.value) return
  downloadingContract.value = true
  try {
    await document.fonts.ready
    const html2pdf = (await import('html2pdf.js')).default
    const number = String(currentDoc.value?.contract_number || 'document').replace(/[^a-zA-Z0-9_-]/g, '-')
    await html2pdf().set({
      margin: [8, 8, 8, 8],
      filename: 'contract-' + number + '.pdf',
      image: { type: 'jpeg', quality: 0.98 },
      html2canvas: { scale: 2, backgroundColor: '#ffffff' },
      jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' },
      pagebreak: { mode: ['css', 'legacy'], avoid: ['.sign-line'] }
    }).from(contractPaper.value).save()
  } catch { await showError('ดาวน์โหลดสัญญาไม่สำเร็จ') }
  finally { downloadingContract.value = false }
}

const isExpiringSoon = (endDateStr: string) => {
  const end = new Date(endDateStr).getTime()
  const now = Date.now()
  const diffDays = (end - now) / (1000 * 3600 * 24)
  return diffDays > 0 && diffDays <= 30
}

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'active': return 'pill-green'
    case 'pending_market_director':
    case 'pending_bureau_director':
    case 'revision_required':
    case 'pending': return 'pill-amber'
    case 'expired': return 'pill-gray'
    case 'terminated': return 'pill-red'
    case 'rejected': return 'pill-red'
    default: return 'pill-gray'
  }
}

const getCircleStatusClass = (status: string) => {
  switch (status) {
    case 'active': return 'green-circle'
    case 'pending_market_director':
    case 'pending_bureau_director':
    case 'revision_required':
    case 'pending': return 'amber-circle'
    case 'expired': return 'gray-circle'
    case 'terminated': return 'red-circle'
    default: return 'blue-circle'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
	case 'draft': return 'แบบร่าง'
	case 'tenant_review': return 'รอผู้เช่าตรวจสอบ'
    case 'active': return 'ใช้งานอยู่'
    case 'pending': return 'รอเจ้าหน้าที่ตรวจ'
    case 'pending_market_director': return 'รอ ผอ.ตลาดลงนาม'
    case 'pending_bureau_director': return 'รอ ผอ.สำนักอนุมัติ'
    case 'revision_required': return 'ส่งกลับแก้ไข'
    case 'expired': return 'หมดสัญญา'
    case 'terminated': return 'ยกเลิกแล้ว'
    case 'rejected': return 'ไม่ผ่านอนุมัติ'
    default: return status
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.draft-submit-label :deep(span:not(:first-child)) {
  display: none;
}
.contract-template-source{display:flex;gap:12px;margin-bottom:16px;padding:14px 16px;border:1px solid #bfdbfe;border-radius:12px;background:#eff6ff;color:#1e40af}.contract-template-source svg{flex:0 0 auto;width:20px;height:20px;margin-top:2px}.contract-template-source strong,.contract-template-source span{display:block}.contract-template-source strong{font-size:var(--font-size-base);font-weight:700}.contract-template-source span{margin-top:2px;font-size:var(--font-size-xs);font-weight:600;color:#2563eb}.contract-template-source p{margin:5px 0 0;font-size:var(--font-size-sm);font-weight:500;color:#475569}.contract-template-code{margin-top:6px;color:#64748b;font-size:var(--font-size-xs);font-weight:600}
/* Page layout, header & card shell: same as the application page */
.contracts-page {
  background-color: var(--bg-f8fafc);
  min-height: 100vh;
  padding-bottom: 60px;
}

.contracts-page .page-container {
  padding: 24px 32px 48px;
  flex: 1;
}

.page-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 8px 0 20px 0;
  gap: 16px;
  flex-wrap: wrap;
}

.main-page-title {
  font-size: 1.55rem;
  font-weight: 700;
  color: var(--tx-0f172a);
  letter-spacing: -0.01em;
  margin: 0;
}

.page-subtitle {
  margin: 4px 0 0;
  color: var(--tx-64748b);
  font-size: var(--font-size-sm);
  font-weight: 500;
}

/* Loading State */
.contracts-loading-state {
  text-align: center;
  padding: 72px 24px;
}

.loading-spinner-row {
  font-size: var(--font-size-base);
  color: var(--tx-64748b);
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Card Shell */
.table-card-container {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 16px;
  padding: 24px 28px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
}

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
  font-size: 1.375rem;
  font-weight: 800;
  color: var(--tx-0f172a);
  letter-spacing: -0.02em;
  margin: 0;
}

.item-count-subtitle {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

/* Search input */
.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  min-width: 260px;
}

.search-input {
  width: 100%;
  padding: 8px 32px 8px 34px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 9px;
  font-size: var(--font-size-xs);
  color: var(--tx-1e293b);
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
  color: var(--tx-94a3b8);
  pointer-events: none;
}

.clear-search-btn {
  position: absolute;
  right: 8px;
  background: transparent;
  border: none;
  color: var(--tx-94a3b8);
  cursor: pointer;
  padding: 3px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.clear-search-btn:hover {
  color: var(--tx-0f172a);
}

/* Officer KPI Ribbon */
.officer-kpi-ribbon {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 24px;
}

@media (max-width: 1200px) {
  .officer-kpi-ribbon {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .officer-kpi-ribbon {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 500px) {
  .officer-kpi-ribbon {
    grid-template-columns: 1fr;
  }
}

.kpi-ribbon-card {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-eef2f6);
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
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
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
  font-size: var(--font-size-2xl);
  font-weight: 800;
  color: var(--tx-0f172a);
  line-height: 1;
}

.kpi-unit-pill {
  font-size: var(--font-size-xs);
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
  line-height: 1.2;
}

/* 1. ทั้งหมด */
.kpi-all .kpi-icon-box {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  border: 1px solid var(--bd-dbeafe);
}
.kpi-all .kpi-unit-pill {
  background: var(--bg-eff6ff);
  color: var(--tx-1d4ed8);
}
.kpi-all.active {
  border-color: #2563eb;
  background: var(--bg-f8faff);
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.18);
}

/* 2. รอดำเนินการ */
.kpi-pending .kpi-icon-box {
  background: var(--bg-fffbeb);
  color: var(--tx-d97706);
  border: 1px solid var(--bd-fef3c7);
}
.kpi-pending .kpi-unit-pill {
  background: var(--bg-fef3c7);
  color: var(--tx-b45309);
}
.kpi-pending.active {
  border-color: #f59e0b;
  background: var(--bg-fffdf7);
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.2);
}

/* 3. ใช้งานอยู่ */
.kpi-applied .kpi-icon-box {
  background: var(--bg-f5f3ff);
  color: var(--tx-7c3aed);
  border: 1px solid var(--bd-ede9fe);
}
.kpi-applied .kpi-unit-pill {
  background: var(--bg-ede9fe);
  color: var(--tx-6d28d9);
}
.kpi-applied.active {
  border-color: #8b5cf6;
  background: var(--bg-faf5ff);
  box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.2);
}

/* 4. หมดสัญญา */
.kpi-expired .kpi-icon-box {
  background: var(--bg-f8fafc);
  color: var(--tx-64748b);
  border: 1px solid var(--bd-e2e8f0);
}
.kpi-expired .kpi-unit-pill {
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
}
.kpi-expired.active {
  border-color: var(--bd-94a3b8);
  background: var(--bg-f8fafc);
  box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.18);
}

/* 5. ยกเลิก/สิ้นสุด */
.kpi-cancelled .kpi-icon-box {
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
  border: 1px solid var(--bd-e2e8f0);
}
.kpi-cancelled .kpi-unit-pill {
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
}
.kpi-cancelled.active {
  border-color: var(--bd-94a3b8);
  background: var(--bg-f8fafc);
  box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.18);
}

/* Sub-filter bar (Pending Stages) */
.tenant-subfilter-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  padding: 6px 10px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-eef2f6);
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
  margin-top: -12px;
  margin-bottom: 22px;
}

.subfilter-chip {
  appearance: none;
  border-radius: 8px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 5px 12px;
  font-family: inherit;
  font-size: var(--font-size-xs);
  font-weight: 500;
  transition: all 0.18s cubic-bezier(0.4, 0, 0.2, 1);
  user-select: none;
  white-space: nowrap;
}

.subfilter-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 19px;
  height: 19px;
  padding: 0 6px;
  border-radius: 9999px;
  font-size: var(--font-size-xs);
  font-weight: 600;
  line-height: 1;
  transition: all 0.18s ease;
}

.subfilter-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

.subfilter-chip.chip-all {
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  color: var(--tx-475569);
}
.subfilter-chip.chip-all .subfilter-badge {
  background: var(--bg-e2e8f0);
  color: var(--tx-475569);
}
.subfilter-chip.chip-all.active {
  background: var(--bg-eff6ff);
  border: 1.5px solid #3b82f6;
  color: var(--tx-1d4ed8);
  font-weight: 600;
  box-shadow: 0 1px 4px rgba(37, 99, 235, 0.12);
}
.subfilter-chip.chip-all.active .subfilter-badge {
  background: var(--bg-dbeafe);
  color: var(--tx-1d4ed8);
  font-weight: 700;
}

.subfilter-chip.chip-amber {
  background: var(--bg-fffbeb);
  border: 1px solid var(--bd-fef3c7);
  color: var(--tx-92400e);
}
.subfilter-chip.chip-amber .subfilter-dot { background: #f59e0b; }
.subfilter-chip.chip-amber .subfilter-badge { background: var(--bg-fef3c7); color: var(--tx-b45309); }
.subfilter-chip.chip-amber.active {
  background: var(--bg-fef3c7);
  border: 1.5px solid #f59e0b;
  color: var(--tx-78350f);
  font-weight: 600;
  box-shadow: 0 1px 4px rgba(245, 158, 11, 0.15);
}
.subfilter-chip.chip-amber.active .subfilter-badge { background: #fde68a; color: var(--tx-78350f); font-weight: 700; }

.subfilter-chip.chip-purple {
  background: var(--bg-faf5ff);
  border: 1px solid var(--bd-f3e8ff);
  color: var(--tx-6b21a8);
}
.subfilter-chip.chip-purple .subfilter-dot { background: #9333ea; }
.subfilter-chip.chip-purple .subfilter-badge { background: var(--bg-f3e8ff); color: var(--tx-7e22ce); }
.subfilter-chip.chip-purple.active {
  background: var(--bg-f3e8ff);
  border: 1.5px solid #9333ea;
  color: var(--tx-581c87);
  font-weight: 600;
  box-shadow: 0 1px 4px rgba(147, 51, 234, 0.15);
}
.subfilter-chip.chip-purple.active .subfilter-badge { background: var(--bg-e9d5ff); color: var(--tx-581c87); font-weight: 700; }

.subfilter-chip.chip-green {
  background: var(--bg-f0fdf4);
  border: 1px solid var(--bd-dcfce7);
  color: var(--tx-166534);
}
.subfilter-chip.chip-green .subfilter-dot { background: #10b981; }
.subfilter-chip.chip-green .subfilter-badge { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.subfilter-chip.chip-green.active {
  background: var(--bg-dcfce7);
  border: 1.5px solid #10b981;
  color: var(--tx-14532d);
  font-weight: 600;
  box-shadow: 0 1px 4px rgba(16, 185, 129, 0.15);
}
.subfilter-chip.chip-green.active .subfilter-badge { background: var(--bg-bbf7d0); color: var(--tx-14532d); font-weight: 700; }

.subfilter-chip.chip-orange {
  background: var(--bg-fff7ed);
  border: 1px solid var(--bd-ffedd5);
  color: var(--tx-9a3412);
}
.subfilter-chip.chip-orange .subfilter-dot { background: #ea580c; }
.subfilter-chip.chip-orange .subfilter-badge { background: var(--bg-ffedd5); color: var(--tx-c2410c); }
.subfilter-chip.chip-orange.active {
  background: var(--bg-ffedd5);
  border: 1.5px solid #ea580c;
  color: var(--tx-7c2d12);
  font-weight: 600;
  box-shadow: 0 1px 4px rgba(234, 88, 12, 0.15);
}
.subfilter-chip.chip-orange.active .subfilter-badge { background: var(--bg-fed7aa); color: var(--tx-7c2d12); font-weight: 700; }

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
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s ease;
}

.officer-project-header:hover {
  background: var(--bg-f1f5f9);
}

.officer-project-header .header-left-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.chevron-indicator {
  color: var(--tx-64748b);
  transition: transform 0.25s ease;
}

.chevron-indicator.rotated {
  transform: rotate(180deg);
}

.accordion-header-icon {
  color: var(--tx-2563eb);
  display: flex;
  align-items: center;
}

.market-title {
  font-size: var(--font-size-base);
  font-weight: 700;
  color: var(--tx-1e293b);
}

.market-count-pill {
  font-size: var(--font-size-xs);
  font-weight: 600;
  padding: 3px 10px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  color: var(--tx-475569);
}

.market-accordion-body {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* Booking Item Card (Application style) */
.booking-item-card {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-left: 4px solid var(--bd-cbd5e1);
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
  border-left: 4px solid var(--bd-94a3b8);
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
  background: linear-gradient(135deg, var(--bg-eff6ff) 0%, var(--bg-e0e7ff) 100%);
  color: var(--tx-2563eb);
  border: 1px solid var(--bd-bfdbfe);
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
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  font-weight: 600;
  line-height: 1.2;
}

.booking-id-text {
  font-size: var(--font-size-xl);
  font-weight: 800;
  color: var(--tx-0f172a);
  letter-spacing: -0.02em;
  line-height: 1.2;
}

.officer-booked-at {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  margin-top: 3px;
  color: var(--tx-94a3b8);
  font-size: var(--font-size-xs);
  font-weight: 500;
}

.top-row-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
  flex-wrap: wrap;
}

.status-pill-badge {
  font-size: var(--font-size-xs);
  font-weight: 600;
  padding: 4px 14px;
  border-radius: 20px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
}

.pill-green {
  background: var(--bg-dcfce7);
  color: var(--tx-15803d);
}

.pill-amber {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
}

.pill-blue {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}

.pill-purple {
  background: var(--bg-f5f3ff);
  color: var(--tx-7c3aed);
}

.pill-gray {
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
}

.pill-red {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
}

.card-collapse-btn {
  background: transparent;
  border: none;
  color: var(--tx-64748b);
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
  color: var(--tx-64748b);
  font-size: var(--font-size-xs);
}

.compact-meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.compact-meta-item + .compact-meta-item {
  padding-left: 12px;
  border-left: 1px solid var(--bd-e2e8f0);
}

.compact-meta-item svg {
  color: var(--tx-94a3b8);
  flex-shrink: 0;
}

.compact-meta-item strong {
  color: var(--tx-334155);
  font-weight: 650;
}

.compact-meta-stall strong {
  color: var(--tx-1d4ed8);
}

.compact-meta-item a {
  color: var(--tx-64748b);
  text-decoration: none;
}

.compact-meta-item a:hover {
  color: var(--tx-2563eb);
  text-decoration: underline;
}

/* Stepper Track */
.booking-stepper-container {
  margin: 18px 0 16px 0;
  padding: 12px 16px;
  background: var(--bg-f8fafc);
  border-radius: 10px;
  border: 1px solid var(--bd-f1f5f9);
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
  background: var(--bg-e2e8f0);
  color: var(--tx-64748b);
  font-size: var(--font-size-xs);
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.25s ease;
}

.step-label-text {
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--tx-64748b);
  white-space: nowrap;
}

.stepper-step.completed .step-circle {
  background: #10b981;
  color: #ffffff;
}

.stepper-step.completed .step-label-text {
  color: var(--tx-0f172a);
}

.stepper-step.active .step-circle {
  background: #2563eb;
  color: #ffffff;
  box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.15);
}

.stepper-step.active .step-label-text {
  color: var(--tx-2563eb);
  font-weight: 700;
}

.stepper-connector {
  flex: 1;
  height: 3px;
  background: var(--bg-e2e8f0);
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
  border-top: 1px solid var(--bd-f1f5f9);
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
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  font-weight: 500;
}

.info-value-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.info-value-text {
  font-size: var(--font-size-base);
  font-weight: 700;
  color: var(--tx-0f172a);
}

.days-remaining-pill {
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
  padding: 2px 8px;
  border-radius: 12px;
  border: 1px solid var(--bd-dbeafe);
}

.stall-count-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: var(--font-size-xs);
  font-weight: 700;
  color: var(--tx-1e293b);
  background: var(--bg-f8fafc);
  padding: 3px 10px;
  border-radius: 8px;
  border: 1px solid var(--bd-e2e8f0);
}

.stall-code-strong {
  font-weight: 700;
  color: var(--tx-1d4ed8);
}

.stall-shop-row {
  display: flex;
  align-items: baseline;
  gap: 5px;
  font-size: var(--font-size-sm);
  color: var(--tx-334155);
}

.shop-prefix {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
}

.shop-name-highlight {
  font-weight: 700;
  color: var(--tx-0f172a);
}

.shop-product {
  color: var(--tx-64748b);
  font-size: var(--font-size-sm);
  overflow-wrap: anywhere;
}

.cell-sub-meeting-point {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
}

.applicant-id-mono {
  font-size: var(--font-size-xs);
  color: var(--tx-475569);
  font-family: inherit;
}

.applicant-id-mono .sub-label {
  color: var(--tx-94a3b8);
}

.applicant-phone-row {
  margin-top: 6px;
}

.phone-call-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  padding: 3px 8px;
  border-radius: 6px;
  font-size: var(--font-size-xs);
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s ease;
}

.phone-call-btn:hover {
  background: #2563eb;
  color: #ffffff;
}

/* Inline Meta Bar */
.application-inline-meta-bar {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px dashed var(--bd-e2e8f0);
}

.inline-details-row {
  display: flex;
  align-items: baseline;
  gap: 6px;
  font-size: var(--font-size-xs);
  line-height: 1.5;
}

.details-lead-label {
  font-weight: 600;
  color: var(--tx-64748b);
  white-space: nowrap;
}

.details-lead-text {
  color: var(--tx-334155);
}

.inline-details-row.text-danger .details-lead-label,
.inline-details-row.text-danger .details-lead-text {
  color: var(--tx-dc2626);
}

/* Card Footer Row */
.card-footer-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 14px;
  margin-top: 10px;
  border-top: 1px solid var(--bd-f1f5f9);
  flex-wrap: wrap;
  gap: 12px;
}

.footer-left-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: var(--font-size-xs);
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

/* Collapsed Next Step */
.collapsed-next-step {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-top: 12px;
  padding-top: 11px;
  border-top: 1px solid var(--bd-f1f5f9);
  font-size: var(--font-size-xs);
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
  color: var(--tx-94a3b8);
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

.tone-waiting { color: var(--tx-d97706); }
.tone-success { color: var(--tx-15803d); }
.tone-action { color: var(--tx-2563eb); }
.tone-muted { color: var(--tx-64748b); }

/* Table View Styles (Matched to Application table) */
.officer-table-container {
  overflow-x: auto;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 12px;
  margin-bottom: 20px;
  background: var(--bg-ffffff);
}

.officer-custom-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: var(--font-size-sm);
}

.officer-custom-table th {
  padding: 13px 16px;
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--tx-64748b);
  background: var(--bg-f8fafc);
  border-bottom: 1px solid var(--bd-eef2f6);
  white-space: nowrap;
}

.officer-custom-table td {
  padding: 15px 16px;
  border-bottom: 1px solid var(--bd-f1f5f9);
  color: var(--tx-334155);
  vertical-align: middle;
}

.officer-custom-table tbody tr:hover {
  background: var(--bg-f8fafc);
}

.table-code-blue {
  color: var(--tx-2563eb);
  font-weight: 700;
  font-size: var(--font-size-sm);
}

.cell-sub-cycle {
  font-size: var(--font-size-xs);
  color: var(--tx-94a3b8);
  margin-top: 2px;
}

.cell-shop-name {
  font-weight: 600;
  color: var(--tx-1e293b);
}

.cell-sub-category {
  font-size: var(--font-size-xs);
  color: var(--tx-94a3b8);
  margin-top: 1px;
}

.cell-tenant-name {
  font-weight: 600;
  color: var(--tx-1e293b);
}

.cell-phone-link {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  display: inline-flex;
  align-items: center;
  gap: 3px;
  text-decoration: none;
  margin-top: 1px;
}

.cell-phone-link:hover {
  color: var(--tx-2563eb);
}

.cell-duration-text {
  color: var(--tx-475569);
  font-size: var(--font-size-xs);
}

.cell-bold-amount {
  font-weight: 700;
  color: var(--tx-0f172a);
  font-size: var(--font-size-base);
}

.action-heading, .action-cell {
  text-align: right;
}

.table-actions-cell {
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
  flex-wrap: wrap;
}

.badge-occupied {
  background: var(--bg-f1f5f9);
  color: var(--tx-334155);
  border: 1px solid var(--bd-e2e8f0);
  padding: 3px 8px;
  border-radius: 6px;
  font-size: var(--font-size-xs);
  font-weight: 600;
}

/* Modal and Print styles (preserved) */
.contract-return-note { color: #b45309; background: #fffbeb; padding: 12px; border-radius: 8px; }
.contract-signature-image { max-width: 150px; max-height: 55px; object-fit: contain; }

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1050;
  padding: 20px;
}

.modal-content {
  background: var(--bg-ffffff);
  border-radius: 16px;
  width: 100%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  padding: 16px 24px;
  background: var(--bg-f8fafc);
  border-top: 1px solid var(--bd-e2e8f0);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  margin-bottom: 16px;
}

.form-label {
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--tx-1e293b);
  margin-bottom: 6px;
}

.req-star {
  color: var(--tx-dc2626);
}

.field-error-text {
  margin: 4px 0 0;
  font-size: var(--font-size-xs);
  color: var(--tx-dc2626);
  font-weight: 500;
}

.tenant-contract-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-bfdbfe);
  border-radius: 10px;
  padding: 14px 16px;
  margin-bottom: 20px;
}

.contract-review-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  background: var(--bg-fefce8);
  border: 1px solid var(--bd-fef08a);
  border-radius: 10px;
  padding: 14px 18px;
  margin-bottom: 18px;
}

.review-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: var(--font-size-xs);
  font-weight: 700;
  color: var(--tx-854d0e);
  background: #fef08a;
  padding: 3px 10px;
  border-radius: 6px;
}

.review-desc {
  font-size: var(--font-size-xs);
  color: var(--tx-713f12);
  margin-top: 4px;
}

.review-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

/* Print Styles */
@media print {
  body * {
    visibility: hidden;
  }
  #print-area, #print-area * {
    visibility: visible;
  }
  #print-area {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
  }
}

/* ============================================================ */
/* CONTRACT DOCUMENT VIEW & PRINT STYLES                        */
/* ============================================================ */
.contract-document {
  max-width: 760px;
  margin: 0 auto;
  padding: 32px 36px;
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  box-shadow: 0 1px 4px rgba(15, 23, 42, 0.04);
  color: var(--tx-1e293b);
  box-sizing: border-box;
}

.contract-header {
  text-align: center;
  margin-bottom: 22px;
  padding-bottom: 18px;
  border-bottom: 1px solid var(--bd-e2e8f0);
}

.contract-doc-title {
  margin: 0;
  font-size: var(--font-size-xl);
  font-weight: 700;
  color: var(--tx-0f172a);
  line-height: 1.35;
  letter-spacing: -0.01em;
}

.contract-doc-org {
  margin-top: 5px;
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  font-weight: 500;
}

.contract-doc-number-pill {
  display: inline-block;
  margin-top: 10px;
  padding: 4px 14px;
  background: var(--bg-ecfdf5);
  border: 1px solid var(--bd-a7f3d0);
  border-radius: 999px;
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--tx-059669);
  letter-spacing: 0.01em;
}

.contract-body {
  font-size: var(--font-size-sm);
  line-height: 1.75;
  color: var(--tx-334155);
  text-align: left;
}

.contract-body p {
  margin: 8px 0;
  line-height: 1.75;
}

.contract-body strong {
  color: var(--tx-0f172a);
  font-weight: 600;
}

.contract-clause-heading {
  font-size: var(--font-size-base);
  font-weight: 700;
  color: var(--tx-0f172a);
  margin-top: 20px;
  margin-bottom: 6px;
  line-height: 1.45;
  letter-spacing: -0.01em;
}

.contract-info-card {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px 20px;
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 10px;
  padding: 14px 18px;
  margin: 10px 0 14px;
  font-size: var(--font-size-sm);
  line-height: 1.6;
}

.contract-info-row {
  display: flex;
  gap: 8px;
  min-width: 0;
  align-items: baseline;
}

.contract-info-full {
  grid-column: 1 / -1;
}

.info-label {
  color: var(--tx-64748b);
  font-weight: 500;
  flex-shrink: 0;
}

.info-value {
  color: var(--tx-1e293b);
  font-weight: 600;
  word-break: break-word;
}

.contract-terms-text {
  color: var(--tx-475569);
  line-height: 1.7;
}

.contract-sign-box {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 32px;
  margin-top: 36px;
  padding-top: 16px;
  text-align: center;
}

.sign-line {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  line-height: 1.8;
}

.dotted-signature-line {
  color: var(--tx-94a3b8);
  letter-spacing: 1px;
}

.signer-name {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--tx-1e293b);
  display: inline-block;
  margin-top: 4px;
}

.contract-signature-image {
  max-height: 60px;
  max-width: 180px;
  object-fit: contain;
  margin: 0 auto 6px;
  display: block;
}

</style>
