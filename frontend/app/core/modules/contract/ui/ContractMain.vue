<template>
  <div class="contracts-page">
    <AppHeader />

    <div class="content-body">
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
        <div class="table-card-header">
          <div class="table-header-left">
            <h2 class="table-card-title">{{ isTenant ? 'แผงเช่าและสัญญาของฉัน' : 'รายการสัญญาเช่าทั้งหมด' }}</h2>
            <div class="table-card-subtitle">{{ filteredContracts.length }} สัญญา</div>
          </div>

          <div class="table-header-actions">
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

            <!-- Outlined Export Report Button -->
            <button type="button" class="btn-export-report" @click="exportReport" title="ส่งออกรายงาน">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"></path><polyline points="16 6 12 2 8 6"></polyline><line x1="12" y1="2" x2="12" y2="15"></line></svg>
              <span>ออกรายงาน</span>
            </button>

            <!-- Primary Action Button (Officers only) -->
            <button v-if="!isTenant" type="button" class="btn-table-primary-action" @click="openCreateModal">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
              <span>สร้างสัญญาใหม่</span>
            </button>
          </div>
        </div>

        <!-- KPI Ribbon Cards (Overview & Status Filter Shortcuts) -->
        <div class="officer-kpi-ribbon">
          <!-- 1. ทั้งหมด -->
          <div
            class="kpi-ribbon-card kpi-all"
            :class="{ active: selectedStatus === '' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = ''; filterContracts(); currentPage = 1"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">สัญญาเช่าทั้งหมด</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ contracts.length }}</span>
                <span class="kpi-unit-pill">ฉบับ</span>
              </div>
            </div>
          </div>

          <!-- 2. ใช้งานอยู่ (Active) -->
          <div
            class="kpi-ribbon-card kpi-applied"
            :class="{ active: selectedStatus === 'active' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'active'; filterContracts(); currentPage = 1"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">ใช้งานอยู่ (Active)</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countContractsByStatus('active') }}</span>
                <span class="kpi-unit-pill">ฉบับ</span>
              </div>
            </div>
          </div>

          <!-- 3. รอดำเนินการ (Pending) -->
          <div
            class="kpi-ribbon-card kpi-pending"
            :class="{ active: selectedStatus === 'pending' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'pending'; filterContracts(); currentPage = 1"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">รอตรวจสัญญา</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countContractsByStatus('pending') }}</span>
                <span class="kpi-unit-pill">ฉบับ</span>
              </div>
            </div>
          </div>

          <!-- 4. หมดสัญญา (Expired) -->
          <div
            class="kpi-ribbon-card kpi-expired"
            :class="{ active: selectedStatus === 'expired' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'expired'; filterContracts(); currentPage = 1"
          >
            <div class="kpi-icon-box">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            </div>
            <div class="kpi-text-box">
              <div class="kpi-title">หมดสัญญาแล้ว</div>
              <div class="kpi-value-row">
                <span class="kpi-count">{{ countContractsByStatus('expired') }}</span>
                <span class="kpi-unit-pill">ฉบับ</span>
              </div>
            </div>
          </div>

          <!-- 5. สิ้นสุด/ยกเลิก (Terminated) -->
          <div
            class="kpi-ribbon-card kpi-terminated"
            :class="{ active: selectedStatus === 'terminated' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'terminated'; filterContracts(); currentPage = 1"
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

        <!-- Empty State -->
        <div v-if="filteredContracts.length === 0" class="empty-state-box">
          <div class="empty-icon" aria-hidden="true">
            <svg width="44" height="44" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line></svg>
          </div>
          <div class="empty-title">
            {{ searchQuery ? `ไม่พบสัญญาเช่าที่ตรงกับ "${searchQuery}"` : (selectedStatus ? `ไม่มีรายการในสถานะ "${getStatusLabel(selectedStatus)}"` : 'ยังไม่มีรายการสัญญาเช่าในระบบ') }}
          </div>
          <div class="empty-desc">
            {{ searchQuery ? 'ลองตรวจสอบคำค้นหา หรือล้างคำค้นหาเพื่อดูรายการทั้งหมด' : (isTenant ? 'เมื่อคำร้องขอเช่าแผงได้รับการอนุมัติ คุณสามารถจัดทำสัญญาเช่าได้ที่นี่' : 'คุณสามารถสร้างสัญญาใหม่ได้โดยกดปุ่ม "สร้างสัญญาใหม่" ด้านบน') }}
          </div>
          <button
            v-if="searchQuery || selectedStatus"
            type="button"
            class="btn-empty-reset"
            @click="searchQuery = ''; selectedStatus = ''; filterContracts(); currentPage = 1"
          >
            ล้างตัวกรองทั้งหมด
          </button>
        </div>

        <!-- Content Mode 1: Table View -->
        <div v-else-if="viewMode === 'table'" class="table-container">
          <table class="custom-table">
            <thead>
              <tr>
                <th style="width: 140px;">เลขที่สัญญา</th>
                <th style="width: 130px;">แผงเช่า</th>
                <th>ชื่อร้านค้า / สินค้า</th>
                <th>ผู้เช่า (Tenant)</th>
                <th style="width: 155px;">ระยะเวลาสัญญา</th>
                <th style="width: 120px;">ค่าเช่า/เดือน</th>
                <th style="width: 110px;">เงินประกัน</th>
                <th style="width: 120px;">สถานะ</th>
                <th style="text-align: right; width: 170px;">การจัดการ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="c in paginatedContracts" :key="c.id">
                <td>
                  <strong class="table-code-blue clickable-code" @click="viewContractDoc(c)" title="คลิกดูเอกสารสัญญา">{{ c.contract_number }}</strong>
                  <div class="cell-sub-cycle">รอบ{{ c.billing_cycle === 'monthly' ? 'รายเดือน' : 'รายวัน' }}</div>
                </td>
                <td>
                  <span class="badge badge-occupied" v-if="c.stall">
                    {{ c.stall.stall_number }} ({{ c.stall.zone_code || 'A' }})
                  </span>
                  <span class="badge badge-occupied" v-else-if="c.stall_ident">
                    {{ c.stall_ident }}
                  </span>
                </td>
                <td>
                  <div class="cell-shop-name">{{ c.business_name || c.shop_name || '-' }}</div>
                  <div class="cell-sub-category">{{ c.business_type || 'ทั่วไป' }}</div>
                </td>
                <td>
                  <div class="cell-tenant-name">{{ c.tenant_name }}</div>
                  <a v-if="c.tenant_phone" :href="`tel:${c.tenant_phone}`" class="cell-phone-link" @click.stop>
                    <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                    <span>{{ c.tenant_phone }}</span>
                  </a>
                </td>
                <td>
                  <div class="cell-duration-text">{{ c.start_date }} ถึง {{ c.end_date }}</div>
                  <div v-if="isExpiringSoon(c.end_date) && c.status === 'active'" class="cell-expiring-tag">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
                    <span>ใกล้หมดอายุ</span>
                  </div>
                </td>
                <td><strong class="cell-bold-amount">฿{{ Number(c.rental_fee).toLocaleString() }}</strong></td>
                <td>฿{{ Number(c.security_deposit || 0).toLocaleString() }}</td>
                <td>
                  <span class="status-pill-badge" :class="getStatusBadgeClass(c.status)">
                    {{ getStatusLabel(c.status) }}
                  </span>
                </td>
                <td style="text-align: right;">
                  <div class="table-actions-cell">
                    <!-- Officer Approve Button -->
                    <AppButton
                      v-if="c.status === 'pending' && isOfficer"
                      variant="primary"
                      size="sm"
                      title="ตรวจสอบและอนุมัติสัญญา"
                      @click="approveContract(c)"
                    >
                      <template #prefix>
                        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      </template>
                      <span>อนุมัติ</span>
                    </AppButton>

                    <!-- View Document Button -->
                    <AppButton
                      variant="secondary"
                      size="sm"
                      title="ดูเอกสารสัญญาฉบับเต็ม"
                      @click="viewContractDoc(c)"
                    >
                      <template #prefix>
                        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line></svg>
                      </template>
                      <span>{{ c.status === 'pending' && isOfficer ? 'ตรวจสัญญา' : 'สัญญา' }}</span>
                    </AppButton>

                    <!-- Officer Terminate Button -->
                    <AppButton
                      v-if="c.status === 'active' && isOfficer"
                      variant="danger-outline"
                      size="sm"
                      title="สิ้นสุด/ยกเลิกสัญญา"
                      @click="terminateContract(c)"
                    >
                      <span>สิ้นสุด</span>
                    </AppButton>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Content Mode 2: Card Grid View (สไตล์เดียวกับ BookingMain / ApplicationMain) -->
        <div v-else-if="viewMode === 'cards'" class="contract-cards-grid">
          <div
            v-for="c in paginatedContracts"
            :key="c.id"
            class="contract-card"
            :class="`status-${c.status}`"
          >
            <!-- Card Header: Circular Status Icon, Contract ID, Billing cycle, Badges -->
            <div class="contract-card-header">
              <div class="contract-header-left">
                <div class="contract-icon-box" :class="getCircleStatusClass(c.status)">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                    <polyline points="14 2 14 8 20 8"></polyline>
                    <line x1="16" y1="13" x2="8" y2="13"></line>
                    <line x1="16" y1="17" x2="8" y2="17"></line>
                  </svg>
                </div>
                <div class="contract-header-meta">
                  <div class="contract-label-text">เลขที่สัญญา</div>
                  <strong class="contract-card-number table-code-blue" @click="viewContractDoc(c)" title="คลิกดูเอกสารสัญญา">
                    {{ c.contract_number }}
                  </strong>
                  <div class="contract-card-cycle">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                    <span>รอบ{{ c.billing_cycle === 'monthly' ? 'รายเดือน' : 'รายวัน' }}</span>
                  </div>
                </div>
              </div>

              <div class="contract-header-right">
                <span v-if="isExpiringSoon(c.end_date) && c.status === 'active'" class="expiring-badge" title="สัญญาใกล้หมดอายุภายใน 30 วัน">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
                  <span>ใกล้หมดอายุ</span>
                </span>
                <span class="status-pill-badge" :class="getStatusBadgeClass(c.status)">
                  {{ getStatusLabel(c.status) }}
                </span>
              </div>
            </div>

            <!-- Shop & Stall Banner -->
            <div class="contract-shop-banner">
              <div class="shop-banner-info">
                <div class="shop-name-title">{{ c.business_name || c.shop_name || 'ร้านค้าทั่วไป' }}</div>
                <div class="shop-category-subtitle">{{ c.business_type || 'สินค้าทั่วไป' }}</div>
              </div>
              <div class="shop-banner-stall">
                <span class="stall-badge-pill" v-if="c.stall">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                  <span>แผง {{ c.stall.stall_number }} ({{ c.stall.zone_code || 'A' }})</span>
                </span>
                <span class="stall-badge-pill" v-else-if="c.stall_ident">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                  <span>แผง {{ c.stall_ident }}</span>
                </span>
              </div>
            </div>

            <!-- Metadata Grid -->
            <div class="contract-card-meta-grid">
              <div class="meta-card-col">
                <span class="meta-card-label">ผู้เช่า (Tenant)</span>
                <div class="meta-card-value">{{ c.tenant_name }}</div>
                <a v-if="c.tenant_phone" :href="`tel:${c.tenant_phone}`" class="meta-card-phone" @click.stop>
                  <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                  <span>{{ c.tenant_phone }}</span>
                </a>
              </div>

              <div class="meta-card-col">
                <span class="meta-card-label">ระยะเวลาสัญญา</span>
                <div class="meta-card-value meta-card-date">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                  <span>{{ c.start_date }} ถึง {{ c.end_date }}</span>
                </div>
              </div>

              <div class="meta-card-col">
                <span class="meta-card-label">ค่าเช่าต่องวด</span>
                <div class="meta-card-amount">฿{{ Number(c.rental_fee).toLocaleString() }}<span class="unit-text">/ด.</span></div>
              </div>

              <div class="meta-card-col">
                <span class="meta-card-label">เงินประกันสัญญา</span>
                <div class="meta-card-value">฿{{ Number(c.security_deposit || 0).toLocaleString() }}</div>
              </div>
            </div>

            <!-- Card Footer Actions -->
            <div class="contract-card-footer">
              <AppButton
                v-if="c.status === 'pending' && isOfficer"
                variant="primary"
                size="sm"
                @click="approveContract(c)"
                title="อนุมัติสัญญาเช่า"
              >
                <template #prefix>
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                </template>
                <span>อนุมัติสัญญา</span>
              </AppButton>

              <AppButton
                variant="secondary"
                size="sm"
                @click="viewContractDoc(c)"
                title="ดูเอกสารสัญญาฉบับเต็ม"
              >
                <template #prefix>
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line></svg>
                </template>
                <span>{{ c.status === 'pending' && isOfficer ? 'ตรวจสัญญา' : 'ดูสัญญา' }}</span>
              </AppButton>

              <AppButton
                v-if="c.status === 'active' && isOfficer"
                variant="danger-outline"
                size="sm"
                @click="terminateContract(c)"
                title="สิ้นสุดสัญญา"
              >
                <span>สิ้นสุดสัญญา</span>
              </AppButton>
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
        <div class="modal-header">
          <h3 style="font-size: 1.2rem; display: flex; align-items: center; gap: 8px;">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--primary);"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
            <span>{{ isTenant ? 'จัดทำสัญญาเช่าพื้นที่แผงค้า (ทำสัญญา)' : 'สร้างสัญญาเช่าพื้นที่ใหม่ (New Rental Contract)' }}</span>
          </h3>
          <button class="btn-logout" @click="showCreateModal = false" aria-label="ปิด">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>

        <form novalidate @submit.prevent="submitContract">
          <div class="modal-body">
            <!-- Notice banner for tenant -->
            <div v-if="isTenant" class="tenant-contract-banner">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line></svg>
              <div>
                <strong style="color: #1e40af; font-size: 0.95rem;">คำร้องขอเช่าแผงค้าได้รับการอนุมัติเรียบร้อยแล้ว</strong>
                <p style="margin: 2px 0 0; color: #3b82f6; font-size: 0.825rem;">กรุณาตรวจสอบข้อมูลสัญญาเช่าด้านล่าง เมื่อกดยืนยันแล้ว สัญญาจะถูกส่งไปยังเจ้าหน้าที่เพื่อตรวจสอบและอนุมัติสัญญา</p>
              </div>
            </div>

            <h4 style="font-size: 0.95rem; margin-bottom: 12px; color: var(--primary);">1. ข้อมูลแผงเช่าและร้านค้า</h4>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_stall_id">เลือกแผงเช่าที่ต้องการเช่า *</label>
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
                <label class="form-label" for="c_business_name">ชื่อร้านค้า (Shop/Business Name) *</label>
                <AppInput id="c_business_name" v-model="formContract.business_name" type="text" placeholder="เช่น ส้มตำรสเด็ด" :error="contractErrors.business_name" @blur="validateContractField('business_name')" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="c_business_type">ประเภทสินค้าที่จำหน่าย</label>
              <AppInput id="c_business_type" v-model="formContract.business_type" type="text" placeholder="เช่น อาหารปรุงสุก, เครื่องดื่ม, เสื้อผ้า" />
            </div>

            <h4 style="font-size: 0.95rem; margin: 20px 0 12px; color: var(--primary);">2. ข้อมูลผู้เช่า (Tenant Information)</h4>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_tenant_name">ชื่อ-นามสกุล ผู้เช่า *</label>
                <AppInput id="c_tenant_name" v-model="formContract.tenant_name" type="text" placeholder="นาย/นาง/น.ส. ..." :error="contractErrors.tenant_name" @blur="validateContractField('tenant_name')" />
              </div>

              <div class="form-group">
                <label class="form-label" for="c_tenant_id_card">เลขประจำตัวประชาชน / Tax ID *</label>
                <AppInput id="c_tenant_id_card" v-model="formContract.tenant_id_card" type="text" placeholder="13 หลัก" :error="contractErrors.tenant_id_card" @blur="validateContractField('tenant_id_card')" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_tenant_phone">เบอร์โทรศัพท์ติดต่อ *</label>
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

            <h4 style="font-size: 0.95rem; margin: 20px 0 12px; color: var(--primary);">3. เงื่อนไขและอัตราค่าเช่า (Lease Terms)</h4>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_start_date">วันที่เริ่มต้นสัญญา *</label>
                <AppInput id="c_start_date" v-model="formContract.start_date" type="date" :error="contractErrors.start_date" @blur="validateContractField('start_date')" />
              </div>

              <div class="form-group">
                <label class="form-label" for="c_end_date">วันที่สิ้นสุดสัญญา *</label>
                <AppInput id="c_end_date" v-model="formContract.end_date" type="date" :error="contractErrors.end_date" @blur="validateContractField('end_date')" />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="c_billing_cycle">รอบการชำระเงิน</label>
                <AppSelect v-model="formContract.billing_cycle" :options="billingCycleOptions" />
              </div>

              <div class="form-group">
                <label class="form-label" for="c_rental_fee">ค่าเช่าต่องวด (บาท) *</label>
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
            <AppButton type="button" variant="secondary" @click="showCreateModal = false">ยกเลิก</AppButton>
            <AppButton type="submit" variant="primary" :disabled="saving">
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
      <div class="modal-content" style="max-width: 800px; background: #fff;">
        <div class="modal-header" style="background: #f8fafc; border-bottom: 1px solid #e2e8f0; color: #0f172a;">
          <h3 style="font-size: 1.15rem;">เอกสารสัญญาเช่าพื้นที่แผงค้า (เลขที่ {{ currentDoc.contract_number }})</h3>
          <div style="display: flex; gap: 8px;">
            <button type="button" class="btn btn-sm btn-primary" style="display: inline-flex; align-items: center; gap: 6px;" @click="printDoc">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 6 2 18 2 18 9"></polyline><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"></path><rect x="6" y="14" width="12" height="8"></rect></svg>
              <span>พิมพ์เอกสาร</span>
            </button>
            <button type="button" class="btn btn-sm btn-secondary" @click="viewDocModal = false">ปิด</button>
          </div>
        </div>

        <div class="modal-body" id="print-area">
          <!-- Officer Review Banner inside Contract Document -->
          <div v-if="currentDoc.status === 'pending'" class="contract-review-card">
            <div class="review-left">
              <div class="review-badge">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                <span>สถานะ: รอดำเนินการ (รอเจ้าหน้าที่ตรวจสอบสัญญา)</span>
              </div>
              <div class="review-desc">ผู้เช่าได้จัดทำสัญญาเช่าเรียบร้อยแล้ว กรุณาตรวจสอบรายละเอียดความถูกต้องก่อนอนุมัติสัญญา</div>
            </div>
            <div v-if="isOfficer" class="review-right">
              <button type="button" class="btn btn-sm btn-outline-danger" @click="rejectContract(currentDoc)">
                ส่งกลับแก้ไข / ปฏิเสธ
              </button>
              <button type="button" class="btn btn-sm btn-primary" style="display: inline-flex; align-items: center; gap: 6px; background: #16a34a; border-color: #16a34a;" @click="approveContract(currentDoc)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span>อนุมัติสัญญาเช่า</span>
              </button>
            </div>
          </div>

          <div class="contract-document">
            <div class="contract-header">
              <h2 style="font-size: 1.4rem; font-weight: 700;">หนังสือสัญญาเช่าพื้นที่แผงค้าในตลาด</h2>
              <div style="font-size: 0.9rem; color: #475569; margin-top: 4px;">ระบบจัดการตลาด e-Market Management System</div>
              <div style="font-size: 0.85rem; font-weight: 600; color: #059669; margin-top: 6px;">
                เลขที่สัญญา: {{ currentDoc.contract_number }}
              </div>
            </div>

            <div style="line-height: 1.8; font-size: 0.92rem; text-align: justify;">
              <p>
                สัญญาฉบับนี้ทำขึ้น ณ <strong>ตลาดกลางชุมชน e-Market</strong> เมื่อวันที่ <strong>{{ currentDoc.created_at ? currentDoc.created_at.slice(0, 10) : currentDoc.start_date }}</strong>
                ระหว่าง <strong>บริษัท บริหารจัดการตลาด จำกัด</strong> (ซึ่งต่อไปนี้เรียกว่า "ผู้ให้เช่า") ฝ่ายหนึ่ง กับ
              </p>

              <div class="contract-grid" style="margin-top: 14px; background: #f8fafc; padding: 14px; border-radius: 8px; border: 1px solid #e2e8f0;">
                <div><strong>ชื่อผู้เช่า:</strong> {{ currentDoc.tenant_name }}</div>
                <div><strong>เลขประจำตัวประชาชน:</strong> {{ currentDoc.tenant_id_card }}</div>
                <div><strong>เบอร์โทรศัพท์:</strong> {{ currentDoc.tenant_phone }}</div>
                <div><strong>อีเมล:</strong> {{ currentDoc.tenant_email || '-' }}</div>
                <div style="grid-column: 1 / -1;"><strong>ที่อยู่:</strong> {{ currentDoc.tenant_address || '-' }}</div>
              </div>

              <div class="section-title">ข้อ ๑. วัตถุประสงค์และสถานที่เช่า</div>
              <p>
                ผู้ให้เช่าตกลงให้เช่า และผู้เช่าตกลงเช่าพื้นที่แผงค้าเลขที่ <strong>{{ currentDoc.stall ? currentDoc.stall.stall_number : (currentDoc.stall_ident || currentDoc.stall_id) }}</strong> 
                เพื่อประกอบกิจการร้านค้า <strong>"{{ currentDoc.business_name || currentDoc.shop_name || 'ร้านค้า' }}"</strong> จำหน่ายสินค้าประเภท <strong>{{ currentDoc.business_type || 'ทั่วไป' }}</strong>
              </p>

              <div class="section-title">ข้อ ๒. ระยะเวลาการเช่า</div>
              <p>
                สัญญาการเช่ามีกำหนดระยะเวลา ตั้งแต่วันที่ <strong>{{ currentDoc.start_date }}</strong> ถึงวันที่ <strong>{{ currentDoc.end_date }}</strong>
              </p>

              <div class="section-title">ข้อ ๓. ค่าเช่าและเงินประกัน</div>
              <div class="contract-grid" style="margin-top: 8px; background: #f8fafc; padding: 14px; border-radius: 8px; border: 1px solid #e2e8f0;">
                <div><strong>อัตราค่าเช่า:</strong> ฿{{ Number(currentDoc.rental_fee).toLocaleString() }} ({{ currentDoc.billing_cycle === 'monthly' ? 'รายเดือน' : 'รายวัน' }})</div>
                <div><strong>เงินประกันสัญญา:</strong> ฿{{ Number(currentDoc.security_deposit || 0).toLocaleString() }} บาท</div>
                <div><strong>เงินจ่ายล่วงหน้า:</strong> ฿{{ Number(currentDoc.advance_payment || 0).toLocaleString() }} บาท</div>
                <div><strong>สถานะสัญญา:</strong> {{ getStatusLabel(currentDoc.status) }}</div>
              </div>

              <div class="section-title">ข้อ ๔. ข้อกำหนดและเงื่อนไขเพิ่มเติม</div>
              <p style="color: #475569;">
                {{ currentDoc.terms_and_conditions || 'ผู้เช่าต้องปฏิบัติตามกฎระเบียบความสะอาด สุขอนามัย และความปลอดภัยของตลาดอย่างเคร่งครัด' }}
              </p>

              <!-- Signature Area -->
              <div class="contract-sign-box">
                <div>
                  <div class="sign-line">
                    ลงชื่อ ...........................................................<br>
                    ( ผู้ให้เช่า / เจ้าหน้าที่ดูแลตลาด )
                  </div>
                </div>
                <div>
                  <div class="sign-line">
                    ลงชื่อ ...........................................................<br>
                    ( {{ currentDoc.tenant_name }} )<br>
                    ผู้เช่า
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'
import AppSegmentedControl from '~/components/AppSegmentedControl.vue'
import AppPagination from '~/components/AppPagination.vue'
import AppButton from '~/components/AppButton.vue'
import AppSelect from '~/components/AppSelect.vue'
import AppInput from '~/components/AppInput.vue'

const route = useRoute()
const authStore = useAuthStore()
const { get, post } = useApi()
const { showSuccess, showError, showConfirm } = useAppAlert()

const isTenant = computed(() => authStore.isTenant || authStore.currentUser?.role === 'tenant')
const isOfficer = computed(() => ['admin', 'staff', 'market_director', 'bureau_director'].includes(authStore.currentUser?.role || ''))

const loading = ref(true)
const saving = ref(false)
const contracts = ref<any[]>([])
const filteredContracts = ref<any[]>([])
const availableStalls = ref<any[]>([])

const selectedStatus = ref('')
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

// View Mode: Table vs Cards
const viewMode = ref<'table' | 'cards'>('table')
const viewModeOptions = [
  { value: 'table', label: 'ตาราง', icon: 'table' as const },
  { value: 'cards', label: 'การ์ด', icon: 'card' as const }
]

const countContractsByStatus = (st: string) => {
  return contracts.value.filter(c => c.status === st).length
}

const paginatedContracts = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredContracts.value.slice(start, start + pageSize.value)
})

const exportReport = async () => {
  await showSuccess(`ส่งออกรายงานสัญญาเช่าพื้นที่ ${filteredContracts.value.length} ฉบับ เรียบร้อยแล้ว`, 'ส่งออกรายงาน')
}

// Create Modal
const showCreateModal = ref(false)
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
    contracts.value = Array.isArray(contractsRes) ? contractsRes : (Array.isArray(contractsRes?.data) ? contractsRes.data : [])
    const rawStalls = Array.isArray(stallsRes) ? stallsRes : (Array.isArray(stallsRes?.data) ? stallsRes.data : [])
    const rawApps = Array.isArray(appsRes) ? appsRes : (Array.isArray(appsRes?.data) ? appsRes.data : [])
    availableStalls.value = rawStalls.filter((s: any) => s.status === 'available' || s.stall_status === 'available')
    if (route.query.stall) searchQuery.value = String(route.query.stall)
    filterContracts()

    // If query contains new_for_stall (from Approved Lease Application)
    if (route.query.new_for_stall) {
      const stallKey = String(route.query.new_for_stall).trim()
      const appId = route.query.app_id ? Number(route.query.app_id) : null
      const matchedApp = rawApps.find((a: any) => appId
        ? Number(a.id) === appId
        : (a.status === 'approved' || a.status === 'contract_pending') && (String(a.stall_number || a.stall_ident) === stallKey || String(a.stall_id) === stallKey)
      )
      if (!matchedApp || (matchedApp.status !== 'approved' && matchedApp.status !== 'contract_pending')) {
        await showError('ไม่พบคำร้องที่อนุมัติขั้นสุดท้ายแล้ว กรุณากลับไปตรวจสอบหน้าคำร้อง', 'ไม่สามารถสร้างสัญญา')
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
    list = list.filter(c => c.status === selectedStatus.value)
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
    contractErrors.rental_fee = String(formContract.value.rental_fee ?? '').trim() === '' || !Number.isFinite(val) || val < 0 ? 'กรุณาระบุค่าเช่าต่องวด' : ''
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

const submitContract = async () => {
  if (!validateContractForm()) return
  saving.value = true
  try {
    const selectedStall = availableStalls.value.find(s => s.id === formContract.value.stall_id)
    const stallIdent = selectedStall?.ident || selectedStall?.stall_number || String(formContract.value.stall_id || '')
    const payload = {
      ...formContract.value,
      application_id: linkedApplicationId.value || undefined,
      stall_ident: stallIdent,
      stall_number: stallIdent,
      shop_name: formContract.value.business_name || formContract.value.tenant_name || 'ร้านค้า',
      monthly_rental: Number(formContract.value.rental_fee) || 0,
      status: isTenant.value ? 'pending' : (formContract.value.status || 'active')
    }
    await post('/contracts', payload)
    showCreateModal.value = false
    await loadData()
    if (isTenant.value) {
      selectedStatus.value = 'pending'
      filterContracts()
      await showSuccess('จัดทำสัญญาเช่าเรียบร้อยแล้ว! สัญญาถูกส่งไปยังเมนูสัญญาให้เจ้าหน้าที่ตรวจสอบสัญญาแล้ว', 'ส่งสัญญาสำเร็จ')
    } else {
      await showSuccess('บันทึกสัญญาเช่าและเปิดใช้งานแผงเช่าเรียบร้อยแล้ว', 'สร้างสัญญาสำเร็จ')
    }
  } catch (err: any) {
    await showError(err?.data?.detail || err?.data?.error?.message || 'เกิดข้อผิดพลาดในการสร้างสัญญา', 'ข้อผิดพลาด')
  } finally {
    saving.value = false
  }
}

const approveContract = async (contract: any) => {
  const confirmed = await showConfirm(
    `คุณต้องการตรวจสอบและอนุมัติสัญญาเช่าเลขที่ ${contract.contract_number} สำหรับแผง ${contract.stall_ident || contract.stall?.stall_number || contract.stall_id} หรือไม่? ระบบจะเปิดใช้งานสัญญาและเปลี่ยนสถานะแผงเป็น "มีผู้เช่า"`,
    {
      title: 'ยืนยันการอนุมัติสัญญาเช่า',
      type: 'info',
      icon: 'shield',
      confirmText: 'อนุมัติสัญญา',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return
  try {
    await post(`/contracts/${contract.id}/approve`)
    viewDocModal.value = false
    await loadData()
    await showSuccess(`อนุมัติสัญญาเช่าเลขที่ ${contract.contract_number} เรียบร้อยแล้ว สัญญามีผลบังคับใช้ทันที`, 'อนุมัติสัญญาสำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || err?.data?.error?.message || 'เกิดข้อผิดพลาดในการอนุมัติสัญญา', 'ข้อผิดพลาด')
  }
}

const rejectContract = async (contract: any) => {
  const confirmed = await showConfirm(
    `คุณต้องการส่งกลับแก้ไขหรือปฏิเสธสัญญาเช่าเลขที่ ${contract.contract_number} หรือไม่?`,
    {
      title: 'ยืนยันปฏิเสธสัญญาเช่า',
      type: 'warning',
      icon: 'triangle-alert',
      confirmText: 'ยืนยันปฏิเสธ',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return
  try {
    await post(`/contracts/${contract.id}/reject`, { reason: 'ข้อมูลสัญญาไม่ถูกต้องหรือไม่ครบถ้วน' })
    viewDocModal.value = false
    await loadData()
    await showSuccess(`ปฏิเสธสัญญาเช่าเลขที่ ${contract.contract_number} เรียบร้อยแล้ว`, 'สำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || err?.data?.error?.message || 'เกิดข้อผิดพลาด', 'ข้อผิดพลาด')
  }
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

const viewContractDoc = (c: any) => {
  currentDoc.value = c
  viewDocModal.value = true
}

const printDoc = () => {
  window.print()
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
    case 'pending': return 'amber-circle'
    case 'expired': return 'gray-circle'
    case 'terminated': return 'red-circle'
    default: return 'blue-circle'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'active': return 'ใช้งานอยู่'
    case 'pending': return 'รอตรวจสัญญา'
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
.contracts-page {
  min-height: 100vh;
  background: #f8fafc;
}

.content-body {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

/* Loading State */
.contracts-loading-state {
  text-align: center;
  padding: 72px 24px;
}

.loading-spinner-row {
  font-size: 0.9375rem;
  color: #64748b;
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

/* Standard Table Card Container */
.table-card-container {
  background: #ffffff;
  border: 1px solid #eef2f6;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
  padding: 24px 28px;
}

/* Table Card Header */
.table-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 20px;
}

.table-header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.table-card-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  letter-spacing: -0.01em;
}

.table-card-subtitle {
  font-size: 0.8125rem;
  color: #64748b;
}

.table-header-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

/* Search Input Box */
.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: #94a3b8;
  pointer-events: none;
}

.search-input {
  width: 250px;
  height: 38px;
  padding: 0 34px 0 36px;
  font-size: 0.875rem;
  font-family: inherit;
  color: #1e293b;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.search-input::placeholder {
  color: #94a3b8;
}

.clear-search-btn {
  position: absolute;
  right: 10px;
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: color 0.15s ease;
}

.clear-search-btn:hover {
  color: #475569;
}

/* Header Buttons */
.btn-export-report {
  height: 38px;
  padding: 0 14px;
  font-size: 0.875rem;
  font-family: inherit;
  font-weight: 500;
  color: #475569;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-export-report:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
  color: #1e293b;
}

.btn-table-primary-action {
  height: 38px;
  padding: 0 16px;
  font-size: 0.875rem;
  font-family: inherit;
  font-weight: 600;
  color: #ffffff;
  background: #2563eb;
  border: none;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 0 1px 2px rgba(37, 99, 235, 0.2);
}

.btn-table-primary-action:hover {
  background: #1d4ed8;
}

/* KPI Ribbon Cards */
.officer-kpi-ribbon {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 12px;
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
  background: #ffffff;
  border: 1px solid #eef2f6;
  border-radius: 14px;
  padding: 13px 16px;
  display: flex;
  align-items: center;
  gap: 13px;
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
  width: 42px;
  height: 42px;
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
  font-size: 0.8125rem;
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
  font-size: 1.35rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1;
}

.kpi-unit-pill {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 2px 7px;
  border-radius: 5px;
  line-height: 1.2;
}

/* KPI Color Variations */
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

.kpi-applied .kpi-icon-box {
  background: #f0fdf4;
  color: #16a34a;
  border: 1px solid #dcfce7;
}
.kpi-applied .kpi-unit-pill {
  background: #dcfce7;
  color: #15803d;
}
.kpi-applied.active {
  border-color: #16a34a;
  background: #f7fee7;
  box-shadow: 0 0 0 2px rgba(22, 163, 74, 0.18);
}

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

.kpi-expired .kpi-icon-box {
  background: #f8fafc;
  color: #64748b;
  border: 1px solid #e2e8f0;
}
.kpi-expired .kpi-unit-pill {
  background: #f1f5f9;
  color: #64748b;
}
.kpi-expired.active {
  border-color: #94a3b8;
  background: #f8fafc;
  box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.18);
}

.kpi-terminated .kpi-icon-box {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fee2e2;
}
.kpi-terminated .kpi-unit-pill {
  background: #fee2e2;
  color: #b91c1c;
}
.kpi-terminated.active {
  border-color: #ef4444;
  background: #fffafa;
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.18);
}

/* Empty State */
.empty-state-box {
  text-align: center;
  padding: 56px 24px;
}

.empty-icon {
  width: 68px;
  height: 68px;
  border-radius: 50%;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #94a3b8;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 6px;
}

.empty-desc {
  font-size: 0.875rem;
  color: #64748b;
  max-width: 440px;
  margin: 0 auto;
}

.btn-empty-reset {
  margin-top: 16px;
  padding: 7px 16px;
  font-size: 0.8125rem;
  font-family: inherit;
  font-weight: 600;
  color: #2563eb;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-empty-reset:hover {
  background: #dbeafe;
}

/* Custom Table Styles */
.table-container {
  overflow-x: auto;
  border: 1px solid #f1f5f9;
  border-radius: 12px;
  margin-bottom: 18px;
}

.custom-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.875rem;
}

.custom-table th {
  padding: 13px 16px;
  font-size: 0.8125rem;
  font-weight: 600;
  color: #64748b;
  background: #f8fafc;
  border-bottom: 1px solid #eef2f6;
  white-space: nowrap;
}

.custom-table td {
  padding: 15px 16px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  vertical-align: middle;
}

.custom-table tbody tr:hover {
  background: #f8fafc;
}

.table-code-blue {
  color: #2563eb;
  font-weight: 700;
  font-size: 0.875rem;
}

.clickable-code {
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 2px;
}

.clickable-code:hover {
  color: #1d4ed8;
}

.cell-sub-cycle {
  font-size: 0.75rem;
  color: #94a3b8;
  margin-top: 2px;
}

.cell-shop-name {
  font-weight: 600;
  color: #1e293b;
}

.cell-sub-category {
  font-size: 0.78rem;
  color: #94a3b8;
  margin-top: 1px;
}

.cell-tenant-name {
  font-weight: 600;
  color: #1e293b;
}

.cell-phone-link {
  font-size: 0.78rem;
  color: #64748b;
  display: inline-flex;
  align-items: center;
  gap: 3px;
  text-decoration: none;
  margin-top: 1px;
}

.cell-phone-link:hover {
  color: #2563eb;
}

.cell-duration-text {
  color: #475569;
  font-size: 0.8125rem;
}

.cell-expiring-tag {
  font-size: 0.75rem;
  color: #dc2626;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 2px;
}

.cell-bold-amount {
  font-weight: 700;
  color: #1e293b;
  font-size: 0.9375rem;
}

.table-actions-cell {
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
}

/* Card Grid View (Cards Mode) */
.contract-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 18px;
  margin-bottom: 20px;
}

.contract-card {
  background: #ffffff;
  border: 1px solid #eef2f6;
  border-radius: 14px;
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
  position: relative;
}

.contract-card:hover {
  border-color: #cbd5e1;
  box-shadow: 0 8px 20px -4px rgba(15, 23, 42, 0.08);
  transform: translateY(-2px);
}

.contract-card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.contract-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.contract-icon-box {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.contract-header-meta {
  display: flex;
  flex-direction: column;
}

.contract-label-text {
  font-size: 0.75rem;
  color: #94a3b8;
  font-weight: 500;
}

.contract-card-number {
  font-size: 0.9375rem;
  font-weight: 700;
  cursor: pointer;
  line-height: 1.3;
}

.contract-card-number:hover {
  text-decoration: underline;
}

.contract-card-cycle {
  font-size: 0.75rem;
  color: #64748b;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 2px;
}

.contract-header-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
}

.expiring-badge {
  font-size: 0.75rem;
  font-weight: 700;
  color: #dc2626;
  background: #fee2e2;
  padding: 2px 8px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

/* Status Pill Badges */
.status-pill-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.78rem;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 7px;
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

.pill-gray {
  background: #f1f5f9;
  color: #64748b;
}

.pill-red {
  background: #fee2e2;
  color: #dc2626;
}

/* Circle Status Icons */
.green-circle {
  background: #dcfce7;
  color: #15803d;
}

.amber-circle {
  background: #fef3c7;
  color: #d97706;
}

.blue-circle {
  background: #eff6ff;
  color: #2563eb;
}

.gray-circle {
  background: #f1f5f9;
  color: #64748b;
}

.red-circle {
  background: #fee2e2;
  color: #dc2626;
}

/* Shop Banner in Card */
.contract-shop-banner {
  background: #f8fafc;
  border: 1px solid #f1f5f9;
  border-radius: 10px;
  padding: 11px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.shop-banner-info {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.shop-name-title {
  font-size: 0.9375rem;
  font-weight: 700;
  color: #1e293b;
  line-height: 1.3;
}

.shop-category-subtitle {
  font-size: 0.8125rem;
  color: #64748b;
  margin-top: 2px;
}

.stall-badge-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
  padding: 3px 9px;
  border-radius: 6px;
  font-size: 0.78rem;
  font-weight: 700;
  white-space: nowrap;
}

/* Card Metadata Grid */
.contract-card-meta-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px 16px;
  padding: 2px 0;
}

.meta-card-col {
  display: flex;
  flex-direction: column;
}

.meta-card-label {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 500;
  margin-bottom: 2px;
}

.meta-card-value {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
  line-height: 1.3;
}

.meta-card-phone {
  font-size: 0.8125rem;
  color: #2563eb;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  text-decoration: none;
  margin-top: 2px;
}

.meta-card-phone:hover {
  text-decoration: underline;
}

.meta-card-date {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #475569;
  font-size: 0.8125rem;
}

.meta-card-amount {
  font-size: 1.05rem;
  font-weight: 800;
  color: #1e293b;
  line-height: 1.2;
}

.meta-card-amount .unit-text {
  font-size: 0.78rem;
  font-weight: 500;
  color: #64748b;
  margin-left: 2px;
}

/* Card Footer Actions */
.contract-card-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
  border-top: 1px solid #f1f5f9;
  padding-top: 12px;
  margin-top: auto;
}

/* Badge helpers */
.badge-occupied {
  background: #f1f5f9;
  color: #334155;
  border: 1px solid #e2e8f0;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 0.78rem;
  font-weight: 600;
}

/* Field Errors */
.field-error-text {
  margin: 4px 0 0;
  font-size: 0.75rem;
  color: #dc2626;
  font-weight: 500;
}

.tenant-contract-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 10px;
  padding: 14px 16px;
  margin-bottom: 20px;
}

.contract-review-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  background: #fefce8;
  border: 1px solid #fef08a;
  border-radius: 10px;
  padding: 14px 18px;
  margin-bottom: 18px;
}

.review-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8125rem;
  font-weight: 700;
  color: #854d0e;
  background: #fef08a;
  padding: 3px 10px;
  border-radius: 6px;
}

.review-desc {
  font-size: 0.8125rem;
  color: #713f12;
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
</style>
