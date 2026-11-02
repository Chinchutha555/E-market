<template>
  <div class="contracts-page">
    <AppHeader />

    <div class="content-body">
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

            <AppButton variant="outline" size="md" title="ออกรายงาน" @click="exportReport">
              <template #prefix><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"/><polyline points="16 6 12 2 8 6"/><line x1="12" y1="2" x2="12" y2="15"/></svg></template>
              ออกรายงาน
            </AppButton>
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
            @click="selectedStatus = ''; filterContracts(); currentPage = 1" @keydown.enter.prevent="selectedStatus = ''; filterContracts(); currentPage = 1" @keydown.space.prevent="selectedStatus = ''; filterContracts(); currentPage = 1"
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

          <!-- 3. รอดำเนินการ (Pending) -->
          <div
            class="kpi-ribbon-card kpi-pending"
            :class="{ active: selectedStatus === 'pending' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'pending'; pendingContractFilter = 'all'; filterContracts(); currentPage = 1" @keydown.enter.prevent="selectedStatus = 'pending'; pendingContractFilter = 'all'; filterContracts(); currentPage = 1" @keydown.space.prevent="selectedStatus = 'pending'; pendingContractFilter = 'all'; filterContracts(); currentPage = 1"
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

          <!-- 2. ใช้งานอยู่ (Active) -->
          <div
            class="kpi-ribbon-card kpi-applied"
            :class="{ active: selectedStatus === 'active' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'active'; filterContracts(); currentPage = 1" @keydown.enter.prevent="selectedStatus = 'active'; filterContracts(); currentPage = 1" @keydown.space.prevent="selectedStatus = 'active'; filterContracts(); currentPage = 1"
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

          <!-- 4. หมดสัญญา (Expired) -->
          <div
            class="kpi-ribbon-card kpi-expired"
            :class="{ active: selectedStatus === 'expired' }"
            role="button"
            tabindex="0"
            @click="selectedStatus = 'expired'; filterContracts(); currentPage = 1" @keydown.enter.prevent="selectedStatus = 'expired'; filterContracts(); currentPage = 1" @keydown.space.prevent="selectedStatus = 'expired'; filterContracts(); currentPage = 1"
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
            @click="selectedStatus = 'terminated'; filterContracts(); currentPage = 1" @keydown.enter.prevent="selectedStatus = 'terminated'; filterContracts(); currentPage = 1" @keydown.space.prevent="selectedStatus = 'terminated'; filterContracts(); currentPage = 1"
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

        <div v-if="selectedStatus === 'pending'" class="contract-subfilters" aria-label="ขั้นตอนสัญญาที่รอดำเนินการ">
          <button v-for="option in pendingContractOptions" :key="option.value" type="button" class="contract-subfilter" :class="{ active: pendingContractFilter === option.value }" :aria-pressed="pendingContractFilter === option.value" @click="pendingContractFilter = option.value; filterContracts(); currentPage = 1">
            {{ option.label }} <span>{{ option.count }}</span>
          </button>
        </div>

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
        <div v-else-if="viewMode === 'table'" class="table-container" role="region" aria-label="ตารางสัญญา" tabindex="0">
          <table class="custom-table">
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

        <!-- Card View -->
        <div v-else class="contract-groups">
          <section v-for="group in groupedContracts" :key="group.key" class="contract-group">
            <button type="button" class="contract-group-heading" :aria-expanded="isContractGroupOpen(group.key)" @click="toggleContractGroup(group.key)">
              <svg class="group-chevron" :class="{ closed: !isContractGroupOpen(group.key) }" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>
              <span class="group-icon"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11H3z"/><path d="M9 21v-8h6v8"/></svg></span>
              <strong>{{ group.name }}</strong><span class="group-count">{{ group.items.length }} สัญญา</span>
            </button>
            <div v-show="isContractGroupOpen(group.key)" class="contract-group-body">
              <article v-for="c in group.items" :key="c.id" class="flow-contract-card" :class="contractBorderClass(c.status)">
                <div class="flow-card-header">
                  <div class="flow-card-identity">
                    <span class="flow-icon"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg></span>
                    <div><div class="flow-label">เลขที่สัญญา</div><strong class="flow-number">{{ c.contract_number }}</strong><div class="flow-date">วันที่ยื่น {{ formatContractDate(c.created_at) }}</div></div>
                  </div>
                  <div class="flow-card-header-actions">
                    <AppButton variant="outline" size="sm" @click="viewContractDoc(c)">ตัวอย่างเอกสาร</AppButton>
                    <span class="status-pill-badge" :class="getStatusBadgeClass(c.status)">{{ getStatusLabel(c.status) }}</span>
                    <button type="button" class="flow-collapse" :aria-expanded="isContractCardOpen(c.id)" :aria-label="isContractCardOpen(c.id) ? 'ย่อรายละเอียดสัญญา' : 'ขยายรายละเอียดสัญญา'" @click="toggleContractCard(c.id)">
                      <svg :class="{ closed: !isContractCardOpen(c.id) }" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>
                    </button>
                  </div>
                </div>
                <div class="flow-meta"><span>แผง <strong>{{ c.stall_ident || '-' }}</strong></span><span>{{ contractProject(c) }}</span><span>{{ contractMarket(c) }}</span><span>{{ c.tenant_name || '-' }} <a v-if="c.tenant_phone" :href="'tel:' + c.tenant_phone">{{ c.tenant_phone }}</a></span></div>
                <div v-show="isContractCardOpen(c.id)" class="flow-card-body">
                  <div class="flow-stepper" aria-label="ขั้นตอนสัญญา">
                    <div v-for="(step, index) in contractStepLabels" :key="step" class="flow-step" :class="contractStepClass(c, index)">
                      <span class="flow-step-circle">
                        <svg v-if="contractStepClass(c, index) === 'done'" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                        <span v-else>{{ index + 1 }}</span>
                      </span>
                      <span>{{ step }}</span>
                    </div>
                  </div>
                  <div class="flow-details">
                    <div><span>ระยะสัญญา</span><strong>{{ formatContractDate(c.start_date) }} – {{ formatContractDate(c.end_date) }}</strong></div>
                    <div><span>ร้านค้า</span><strong>{{ c.shop_name || '-' }}</strong><small>{{ c.business_type || '-' }}</small></div>
                    <div><span>ผู้เช่า</span><strong>{{ c.tenant_name || '-' }}</strong><small>{{ c.tenant_phone || '-' }}</small></div>
                    <div><span>ค่าเช่า</span><strong>฿{{ Number(c.rental_fee || 0).toLocaleString() }} {{ c.billing_cycle === 'daily' ? '/ วัน' : '/ เดือน' }}</strong><small>เงินประกัน ฿{{ Number(c.security_deposit || 0).toLocaleString() }}</small></div>
                  </div>
                </div>
                <div class="flow-card-footer">
                  <div class="flow-progress"><span>ความคืบหน้า</span><strong>{{ contractProgressText(c) }}</strong></div>
                  <div class="flow-actions">
                    <AppButton v-if="canReviewContract(c) && c.status !== 'pending_bureau_director'" variant="danger-outline" size="md" @click="rejectContract(c)">ส่งกลับแก้ไข</AppButton>
                    <AppButton v-if="isTenant && c.status === 'revision_required'" variant="primary" size="md" @click="openContractRevision(c)">แก้ไขและยื่นใหม่</AppButton>
                    <AppButton variant="secondary" size="md" @click="viewContractDoc(c)">
                      <template #prefix><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8z"/></svg></template>
                      รายละเอียด
                    </AppButton>
                    <AppButton v-if="canReviewContract(c)" variant="primary" size="md" @click="approveContract(c)">{{ c.status === 'pending_bureau_director' ? 'ลงนามและอนุมัติ' : 'ลงนามและส่งต่อ' }}</AppButton>
                    <AppButton v-if="canTerminateContract(c)" variant="danger-outline" size="md" @click="terminateContract(c)">สิ้นสุดสัญญา</AppButton>
                  </div>
                </div>
              </article>
            </div>
          </section>
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
            <p v-if="editingContract?.review_reason" class="contract-return-note">เหตุผลที่ส่งกลับ: {{ editingContract.review_reason }}</p>
            <AppFileUpload v-model="contractDocument" label="เอกสารประกอบสัญญา" required-text="จำเป็น" accept=".pdf,.jpg,.jpeg,.png,application/pdf,image/jpeg,image/png" :max-size-mb="3" hint="PDF, JPG หรือ PNG ขนาดไม่เกิน 3MB" :error="contractDocumentError" @error="contractDocumentError = $event" />
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
            <AppButton type="submit" variant="primary" size="md" :disabled="saving">
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
              <h2 style="font-size: 1.375rem; font-weight: 700;">หนังสือสัญญาเช่าพื้นที่แผงค้าในตลาด</h2>
              <div style="font-size: var(--font-size-sm); color: var(--tx-475569); margin-top: 4px;">ระบบจัดการตลาด e-Market Management System</div>
              <div style="font-size: var(--font-size-sm); font-weight: 600; color: var(--tx-059669); margin-top: 6px;">
                เลขที่สัญญา: {{ currentDoc.contract_number }}
              </div>
            </div>

            <div style="line-height: 1.8; font-size: var(--font-size-base); text-align: justify;">
              <p>
                สัญญาฉบับนี้ทำขึ้น ณ <strong>ตลาดกลางชุมชน e-Market</strong> เมื่อวันที่ <strong>{{ currentDoc.created_at ? currentDoc.created_at.slice(0, 10) : currentDoc.start_date }}</strong>
                ระหว่าง <strong>บริษัท บริหารจัดการตลาด จำกัด</strong> (ซึ่งต่อไปนี้เรียกว่า "ผู้ให้เช่า") ฝ่ายหนึ่ง กับ
              </p>

              <div class="contract-grid" style="margin-top: 14px; background: var(--bg-f8fafc); padding: 14px; border-radius: 8px; border: 1px solid var(--bd-e2e8f0);">
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
              <div class="contract-grid" style="margin-top: 8px; background: var(--bg-f8fafc); padding: 14px; border-radius: 8px; border: 1px solid var(--bd-e2e8f0);">
                <div><strong>อัตราค่าเช่า:</strong> ฿{{ Number(currentDoc.rental_fee).toLocaleString() }} ({{ currentDoc.billing_cycle === 'monthly' ? 'รายเดือน' : 'รายวัน' }})</div>
                <div><strong>เงินประกันสัญญา:</strong> ฿{{ Number(currentDoc.security_deposit || 0).toLocaleString() }} บาท</div>
                <div><strong>เงินจ่ายล่วงหน้า:</strong> ฿{{ Number(currentDoc.advance_payment || 0).toLocaleString() }} บาท</div>
                <div><strong>สถานะสัญญา:</strong> {{ getStatusLabel(currentDoc.status) }}</div>
              </div>

              <div class="section-title">ข้อ ๔. ข้อกำหนดและเงื่อนไขเพิ่มเติม</div>
              <p style="color: var(--tx-475569);">
                {{ currentDoc.terms_and_conditions || 'ผู้เช่าต้องปฏิบัติตามกฎระเบียบความสะอาด สุขอนามัย และความปลอดภัยของตลาดอย่างเคร่งครัด' }}
              </p>

              <div class="contract-sign-box">
                <div v-for="signer in contractSignatures" :key="signer.kind" class="sign-line">
                  <img v-if="signer.url" :src="signer.url" :alt="'ลายเซ็น ' + signer.label" class="contract-signature-image" />
                  <span v-else>...........................................................</span><br>
                  ( {{ signer.label }} )
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
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
const { get, post, apiFetch } = useApi()
const { showSuccess, showError, showConfirm } = useAppAlert()

const isTenant = computed(() => authStore.isTenant || authStore.currentUser?.role === 'tenant')
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
const viewMode = ref<'table' | 'cards'>('table')
const viewModeOptions = [
  { value: 'table', label: 'ตาราง', icon: 'table' as const },
  { value: 'cards', label: 'การ์ด', icon: 'card' as const }
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
        : (a.status === 'completed' || a.status === 'contract_pending') && (String(a.stall_number || a.stall_ident) === stallKey || String(a.stall_id) === stallKey)
      )
      if (!matchedApp || (matchedApp.status !== 'completed' && matchedApp.status !== 'contract_pending')) {
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
const submitContract = async () => {
  if (!validateContractForm()) return
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
.page-header-row { display: flex; justify-content: space-between; align-items: flex-start; margin: 0 0 22px; }
.main-page-title { margin: 0; color: #0f172a; font-size: 1.5rem; font-weight: 800; }
.page-subtitle { margin: 4px 0 0; color: #64748b; font-size: 0.9375rem; }
.table-card-title { font-size: 1.125rem; font-weight: 600; }
.table-card-subtitle { font-size: 0.8125rem; }
.search-input { height: 40px; }
.table-container { overflow-x: auto; border: 0; border-radius: 0; }
.custom-table { min-width: 1050px; }
.custom-table th { padding: 14px 16px; color: #64748b; background: #fff; border-bottom: 1px solid #f1f5f9; font-size: 0.8125rem; font-weight: 500; }
.custom-table td { padding: 16px; font-size: 0.875rem; }
.custom-table tbody tr:hover { background: #f8fafc; }
.action-heading, .action-cell { text-align: right; }
.table-actions-cell { display: inline-flex; flex-wrap: wrap; justify-content: flex-end; }
.cell-sub-cycle, .cell-sub-category { font-size: 0.8125rem; }
.cell-phone-link, .cell-duration-text { font-size: 0.875rem; }
.contract-groups { display: flex; flex-direction: column; gap: 16px; margin-bottom: 20px; }
.contract-group { overflow: hidden; border: 1px solid #e2e8f0; border-radius: 14px; background: #fff; }
.contract-group-heading { width: 100%; min-height: 58px; display: flex; align-items: center; gap: 10px; padding: 12px 16px; color: #1e293b; background: #f8fafc; border: 0; text-align: left; font: inherit; cursor: pointer; }
.contract-group-heading strong { font-size: 1rem; font-weight: 700; }
.group-icon { display: inline-flex; color: #2563eb; }
.group-chevron { color: #64748b; transform: rotate(180deg); transition: transform .15s; }
.group-chevron.closed { transform: none; }
.group-count { padding: 2px 10px; border: 1px solid #e2e8f0; border-radius: 999px; background: #fff; color: #64748b; font-size: 0.8125rem; font-weight: 600; }
.contract-group-body { display: flex; flex-direction: column; gap: 12px; padding: 12px; }
.flow-contract-card { padding: 20px 22px; border: 1px solid #e2e8f0; border-left: 4px solid #2563eb; border-radius: 14px; background: #fff; box-shadow: 0 1px 3px rgba(0,0,0,.03); }
.flow-contract-card.border-active { border-left-color: #10b981; }
.flow-contract-card.border-pending { border-left-color: #f59e0b; }
.flow-contract-card.border-danger { border-left-color: #ef4444; }
.flow-contract-card.border-muted { border-left-color: #94a3b8; }
.flow-card-header, .flow-card-identity, .flow-card-header-actions, .flow-meta, .flow-card-footer, .flow-actions { display: flex; align-items: center; }
.flow-card-header { justify-content: space-between; gap: 16px; }
.flow-card-identity { gap: 14px; min-width: 0; }
.flow-icon { display: grid; place-items: center; width: 48px; height: 48px; flex: 0 0 48px; color: #2563eb; background: #eff6ff; border: 1px solid #dbeafe; border-radius: 12px; }
.flow-label, .flow-date { color: #64748b; font-size: 0.8125rem; }
.flow-number { display: block; color: #0f172a; font-size: 1.125rem; font-weight: 800; overflow-wrap: anywhere; }
.flow-card-header-actions { justify-content: flex-end; flex-wrap: wrap; gap: 10px; }
.flow-collapse { display: grid; place-items: center; width: 32px; height: 32px; color: #64748b; border: 0; border-radius: 8px; background: transparent; cursor: pointer; }
.flow-collapse:hover { background: #f1f5f9; }
.flow-collapse svg { transform: rotate(180deg); transition: transform .15s; }
.flow-collapse svg.closed { transform: none; }
.flow-meta { flex-wrap: wrap; gap: 8px 0; margin-top: 16px; color: #64748b; font-size: 0.875rem; }
.flow-meta > span { padding: 0 12px; border-right: 1px solid #e2e8f0; }
.flow-meta > span:first-child { padding-left: 0; }
.flow-meta > span:last-child { border: 0; }
.flow-meta strong, .flow-meta a { color: #2563eb; }
.flow-meta a { margin-left: 5px; text-decoration: none; }
.flow-card-body { margin-top: 18px; }
.flow-stepper { display: flex; overflow-x: auto; gap: 8px; padding: 16px; border: 1px solid #f1f5f9; border-radius: 10px; background: #f8fafc; }
.flow-step { position: relative; min-width: 125px; flex: 1; display: flex; flex-direction: column; align-items: center; gap: 6px; color: #64748b; text-align: center; font-size: 0.8125rem; font-weight: 600; }
.flow-step::after { position: absolute; top: 15px; left: calc(50% + 23px); width: calc(100% - 38px); height: 2px; background: #cbd5e1; content: ''; }
.flow-step:last-child::after { display: none; }
.flow-step.done::after { background: #10b981; }
.flow-step-circle { position: relative; z-index: 1; display: grid; place-items: center; width: 32px; height: 32px; border: 1px solid #cbd5e1; border-radius: 50%; background: #fff; font-size: 0.75rem; }
.flow-step.done .flow-step-circle { border-color: #10b981; background: #10b981; color: #fff; }
.flow-step.current .flow-step-circle { border-color: #2563eb; background: #2563eb; color: #fff; }
.flow-step.current { color: #2563eb; }
.flow-details { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 18px; padding: 20px 0; }
.flow-details > div { min-width: 0; display: flex; flex-direction: column; gap: 4px; }
.flow-details span, .flow-details small { color: #64748b; font-size: 0.8125rem; }
.flow-details strong { color: #1e293b; font-size: 0.9375rem; font-weight: 700; overflow-wrap: anywhere; }
.flow-card-footer { justify-content: space-between; flex-wrap: wrap; gap: 12px; padding-top: 14px; border-top: 1px solid #f1f5f9; }
.flow-progress { display: flex; flex-wrap: wrap; gap: 8px; color: #64748b; font-size: 0.8125rem; }
.flow-progress strong { color: #334155; font-weight: 600; }
.flow-actions { justify-content: flex-end; flex-wrap: wrap; gap: 8px; }
@media (max-width: 800px) { .flow-details { grid-template-columns: repeat(2, minmax(0, 1fr)); } .flow-card-header { align-items: flex-start; } }
@media (max-width: 550px) { .flow-contract-card { padding: 16px; } .flow-card-header { flex-wrap: wrap; } .flow-card-header-actions { width: 100%; justify-content: flex-start; } .flow-card-footer { align-items: flex-start; } .flow-actions { justify-content: flex-start; } }
.contract-return-note { color: #b45309; background: #fffbeb; padding: 12px; border-radius: 8px; }
.contract-signature-image { max-width: 150px; max-height: 55px; object-fit: contain; }
.contracts-page {
  min-height: 100vh;
  background: var(--bg-f8fafc);
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

/* Standard Table Card Container */
.table-card-container {
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-eef2f6);
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
  font-size: var(--font-size-xl);
  font-weight: 700;
  color: var(--tx-1e293b);
  margin: 0;
  letter-spacing: -0.01em;
}

.table-card-subtitle {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
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
  color: var(--tx-94a3b8);
  pointer-events: none;
}

.search-input {
  width: 250px;
  height: 38px;
  padding: 0 34px 0 36px;
  font-size: var(--font-size-sm);
  font-family: inherit;
  color: var(--tx-1e293b);
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 8px;
  transition: all 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.search-input::placeholder {
  color: var(--tx-94a3b8);
}

.clear-search-btn {
  position: absolute;
  right: 10px;
  background: none;
  border: none;
  color: var(--tx-94a3b8);
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: color 0.15s ease;
}

.clear-search-btn:hover {
  color: var(--tx-475569);
}

/* Header Buttons */
.btn-export-report {
  height: 38px;
  padding: 0 14px;
  font-size: var(--font-size-sm);
  font-family: inherit;
  font-weight: 500;
  color: var(--tx-475569);
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-export-report:hover {
  background: var(--bg-f8fafc);
  border-color: var(--bd-cbd5e1);
  color: var(--tx-1e293b);
}

.btn-table-primary-action {
  height: 38px;
  padding: 0 16px;
  font-size: var(--font-size-sm);
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
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-eef2f6);
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
  font-size: 1.375rem;
  font-weight: 800;
  color: var(--tx-0f172a);
  line-height: 1;
}

.kpi-unit-pill {
  font-size: var(--font-size-xs);
  font-weight: 600;
  padding: 2px 7px;
  border-radius: 5px;
  line-height: 1.2;
}

/* KPI Color Variations */
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

.kpi-applied .kpi-icon-box {
  background: var(--bg-f0fdf4);
  color: var(--tx-16a34a);
  border: 1px solid var(--bd-dcfce7);
}
.kpi-applied .kpi-unit-pill {
  background: var(--bg-dcfce7);
  color: var(--tx-15803d);
}
.kpi-applied.active {
  border-color: #16a34a;
  background: var(--bg-f7fee7);
  box-shadow: 0 0 0 2px rgba(22, 163, 74, 0.18);
}

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

.kpi-terminated .kpi-icon-box {
  background: var(--bg-fef2f2);
  color: var(--tx-dc2626);
  border: 1px solid var(--bd-fee2e2);
}
.kpi-terminated .kpi-unit-pill {
  background: var(--bg-fee2e2);
  color: var(--tx-b91c1c);
}
.kpi-terminated.active {
  border-color: #ef4444;
  background: var(--bg-fffafa);
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
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  color: var(--tx-94a3b8);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.empty-title {
  font-size: var(--font-size-lg);
  font-weight: 700;
  color: var(--tx-1e293b);
  margin-bottom: 6px;
}

.empty-desc {
  font-size: var(--font-size-sm);
  color: var(--tx-64748b);
  max-width: 440px;
  margin: 0 auto;
}

.btn-empty-reset {
  margin-top: 16px;
  padding: 7px 16px;
  font-size: var(--font-size-xs);
  font-family: inherit;
  font-weight: 600;
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-bfdbfe);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-empty-reset:hover {
  background: var(--bg-dbeafe);
}

/* Custom Table Styles */
.table-container {
  overflow-x: auto;
  border: 1px solid var(--bd-f1f5f9);
  border-radius: 12px;
  margin-bottom: 18px;
}

.custom-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: var(--font-size-sm);
}

.custom-table th {
  padding: 13px 16px;
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--tx-64748b);
  background: var(--bg-f8fafc);
  border-bottom: 1px solid var(--bd-eef2f6);
  white-space: nowrap;
}

.custom-table td {
  padding: 15px 16px;
  border-bottom: 1px solid var(--bd-f1f5f9);
  color: var(--tx-334155);
  vertical-align: middle;
}

.custom-table tbody tr:hover {
  background: var(--bg-f8fafc);
}

.table-code-blue {
  color: var(--tx-2563eb);
  font-weight: 700;
  font-size: var(--font-size-sm);
}

.clickable-code {
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 2px;
}

.clickable-code:hover {
  color: var(--tx-1d4ed8);
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

.cell-expiring-tag {
  font-size: var(--font-size-xs);
  color: var(--tx-dc2626);
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 2px;
}

.cell-bold-amount {
  font-weight: 700;
  color: var(--tx-1e293b);
  font-size: var(--font-size-base);
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
  background: var(--bg-ffffff);
  border: 1px solid var(--bd-eef2f6);
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
  border-color: var(--bd-cbd5e1);
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
  font-size: var(--font-size-xs);
  color: var(--tx-94a3b8);
  font-weight: 500;
}

.contract-card-number {
  font-size: var(--font-size-base);
  font-weight: 700;
  cursor: pointer;
  line-height: 1.3;
}

.contract-card-number:hover {
  text-decoration: underline;
}

.contract-card-cycle {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
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
  font-size: var(--font-size-xs);
  font-weight: 700;
  color: var(--tx-dc2626);
  background: var(--bg-fee2e2);
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
  font-size: var(--font-size-xs);
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 7px;
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

.pill-gray {
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
}

.pill-red {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
}

/* Circle Status Icons */
.green-circle {
  background: var(--bg-dcfce7);
  color: var(--tx-15803d);
}

.amber-circle {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
}

.blue-circle {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}

.gray-circle {
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
}

.red-circle {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
}

/* Shop Banner in Card */
.contract-shop-banner {
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-f1f5f9);
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
  font-size: var(--font-size-base);
  font-weight: 700;
  color: var(--tx-1e293b);
  line-height: 1.3;
}

.shop-category-subtitle {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  margin-top: 2px;
}

.stall-badge-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: var(--bg-eff6ff);
  color: var(--tx-1d4ed8);
  border: 1px solid var(--bd-bfdbfe);
  padding: 3px 9px;
  border-radius: 6px;
  font-size: var(--font-size-xs);
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
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  font-weight: 500;
  margin-bottom: 2px;
}

.meta-card-value {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--tx-1e293b);
  line-height: 1.3;
}

.meta-card-phone {
  font-size: var(--font-size-xs);
  color: var(--tx-2563eb);
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
  color: var(--tx-475569);
  font-size: var(--font-size-xs);
}

.meta-card-amount {
  font-size: var(--font-size-lg);
  font-weight: 800;
  color: var(--tx-1e293b);
  line-height: 1.2;
}

.meta-card-amount .unit-text {
  font-size: var(--font-size-xs);
  font-weight: 500;
  color: var(--tx-64748b);
  margin-left: 2px;
}

/* Card Footer Actions */
.contract-card-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
  border-top: 1px solid var(--bd-f1f5f9);
  padding-top: 12px;
  margin-top: auto;
}

/* Badge helpers */
.badge-occupied {
  background: var(--bg-f1f5f9);
  color: var(--tx-334155);
  border: 1px solid var(--bd-e2e8f0);
  padding: 3px 8px;
  border-radius: 6px;
  font-size: var(--font-size-xs);
  font-weight: 600;
}

/* Field Errors */
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
/* Match the application list toolbar and table. */
.table-card-title { font-size: 1.125rem; font-weight: 600; }
.search-input { height: 40px; }
.table-container { border: 0; border-radius: 0; }
.custom-table th { padding: 14px 16px; color: #64748b; background: #fff; border-bottom: 1px solid #f1f5f9; font-size: 0.8125rem; font-weight: 500; }
.custom-table td { padding: 16px; font-size: 0.875rem; }
.cell-sub-cycle, .cell-sub-category { font-size: 0.8125rem; }
.cell-phone-link, .cell-duration-text { font-size: 0.875rem; }
.contract-subfilters { display: flex; flex-wrap: wrap; gap: 8px; margin: -8px 0 20px; }
.contract-subfilter { display: inline-flex; align-items: center; gap: 6px; padding: 6px 12px; border: 1px solid #e2e8f0; border-radius: 999px; background: #fff; color: #64748b; font: inherit; font-size: 0.875rem; font-weight: 600; cursor: pointer; }
.contract-subfilter span { display: inline-grid; place-items: center; min-width: 22px; height: 22px; padding: 0 5px; border-radius: 999px; background: #f1f5f9; font-size: 0.75rem; }
.contract-subfilter.active { border-color: #93c5fd; background: #eff6ff; color: #2563eb; }
.contract-subfilter.active span { background: #dbeafe; }
.status-pill-badge { font-size: 0.8125rem; padding: 4px 12px; }
.contract-group-heading:focus-visible, .flow-collapse:focus-visible, .contract-subfilter:focus-visible { outline: 2px solid #2563eb; outline-offset: -2px; }
</style>
