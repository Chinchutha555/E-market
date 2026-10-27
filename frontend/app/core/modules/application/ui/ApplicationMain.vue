<template>
  <div>
    <AppHeader />

    <div class="content-body">


      <!-- Loading State -->
      <div v-if="loading" style="text-align: center; padding: 60px;">
        <div style="font-size: 1rem; color: #64748b; display: inline-flex; align-items: center; gap: 8px;">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          <span>กำลังโหลดรายการคำร้องขอเช่า...</span>
        </div>
      </div>

      <!-- Standard Table Card Container (Mockup Standard) -->
      <div v-else class="table-card-container">
        <!-- Card Header: Title + Subtitle on Left, Export + Actions + Search on Right -->
        <div class="table-card-header">
          <div class="table-header-left">
            <h2 class="table-card-title">รายการคำร้องทั้งหมด</h2>
            <div class="table-card-subtitle">{{ filteredApplications.length }} รายการ</div>
          </div>

          <div class="table-header-actions">
            <!-- Reset Test Data Button -->
            <AppButton
              variant="danger-outline"
              size="sm"
              @click="handleClearAllData"
              title="ล้างข้อมูลการจองและคำร้องทั้งหมด เพื่อเริ่มทดสอบ Flow ใหม่จาก 0"
            >
              <template #prefix>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
              </template>
              <span>รีเซ็ตข้อมูลทดสอบ</span>
            </AppButton>

            <!-- Outlined Export Report Button -->
            <button class="btn-export-report" @click="exportReport" title="ส่งออกรายงาน">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"></path><polyline points="16 6 12 2 8 6"></polyline><line x1="12" y1="2" x2="12" y2="15"></line></svg>
              <span>ออกรายงาน</span>
            </button>

            <!-- Primary Action Button: ยื่นคำร้องทั่วไป -->
            <button type="button" class="btn-table-primary-action" @click="openCreateGeneralModal" title="ยื่นคำร้องทั่วไป">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
              <span>ยื่นคำร้องทั่วไป</span>
            </button>

            <!-- Secondary Action: ขอเช่าแผงค้า (จองดูแผง) -->
            <NuxtLink to="/bookings" class="btn-table-stall-action" title="ไปหน้าจองดูแผงค้าและยื่นขอเช่า">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
              <span>ขอเช่าแผงค้า</span>
            </NuxtLink>

            <!-- Manage Application Types (Admin & Staff) -->
            <NuxtLink
              v-if="isAdminOrStaff"
              to="/applications/types"
              class="btn-table-manage-types"
              title="จัดการประเภทคำร้อง (ผู้ดูแลระบบและเจ้าหน้าที่)"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
              <span>จัดการประเภทคำร้อง</span>
            </NuxtLink>

            <!-- Search Box with Search Icon Left & Funnel Filter Icon Right -->
            <div class="table-search-funnel-box">
              <span class="search-icon-left">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
              </span>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="ค้นหาเลขที่คำร้อง / ชื่อ / เรื่อง..."
                @input="currentPage = 1"
              />
              <span class="filter-icon-right" title="ตัวกรอง">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"></polygon></svg>
              </span>
            </div>

            <!-- View Mode Switcher: Table vs Cards -->
            <AppSegmentedControl
              v-model="viewMode"
              :options="viewModeOptions"
              aria-label="เลือกมุมมองการแสดงรายการคำร้อง"
            />
          </div>
        </div>

        <!-- Application Type Filter Bar (แยกประเภทคำร้อง) -->
        <div class="application-type-bar">
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

        <!-- Filter Tabs (ทั้งหมด, รอ จนท., รอ ผอ.ตลาด, รอ ผอ.สำนัก, อนุมัติแล้ว, ส่งกลับแก้ไข) -->
        <div class="table-nav-tabs">
          <button
            class="table-nav-tab"
            :class="{ active: activeTab === 'all' }"
            @click="activeTab = 'all'; currentPage = 1"
          >
            ทั้งหมด ({{ applications.length }})
          </button>
          <button
            class="table-nav-tab"
            :class="{ active: activeTab === 'pending_officer' }"
            @click="activeTab = 'pending_officer'; currentPage = 1"
          >
            รอ จนท. ตรวจสอบ ({{ countByStatus('pending_officer') }})
          </button>
          <button
            class="table-nav-tab"
            :class="{ active: activeTab === 'pending_market_director' }"
            @click="activeTab = 'pending_market_director'; currentPage = 1"
          >
            รอ ผอ.ตลาด ({{ countByStatus('pending_market_director') }})
          </button>
          <button
            class="table-nav-tab"
            :class="{ active: activeTab === 'pending_bureau_director' }"
            @click="activeTab = 'pending_bureau_director'; currentPage = 1"
          >
            รอ ผอ.สำนักตลาด ({{ countByStatus('pending_bureau_director') }})
          </button>
          <button
            class="table-nav-tab"
            :class="{ active: activeTab === 'approved' }"
            @click="activeTab = 'approved'; currentPage = 1"
          >
            อนุมัติแล้ว ({{ countByStatus('approved') }})
          </button>
          <button
            class="table-nav-tab"
            :class="{ active: activeTab === 'rejected' }"
            @click="activeTab = 'rejected'; currentPage = 1"
          >
            ส่งกลับ/ไม่อนุมัติ ({{ countByStatus('rejected') + countByStatus('revision_required') }})
          </button>
        </div>

        <!-- View 1: Table View -->
        <div v-if="viewMode === 'table'" class="table-container">
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
                  <span class="status-pill-badge" :class="getStatusBadgeClass(app.status)">
                    {{ getStatusLabel(app.status) }}
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
                    <div class="review-comment-pill" :class="`comment-${getReviewNoteInfo(app).type}`">
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
                      v-if="app.status === 'pending_officer'"
                      class="btn-table-action btn-table-primary"
                      @click="processApproval(app, 'เจ้าหน้าที่ตรวจสอบเอกสารครบถ้วน ส่งต่อ ผอ.ตลาด')"
                      title="เจ้าหน้าที่ตรวจสอบ"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      <span>ตรวจแล้ว</span>
                    </button>

                    <!-- Step 2: Market Director Review -->
                    <button
                      v-else-if="['pending_market_director', 'pending_debt_check', 'debt_pending'].includes(app.status)"
                      class="btn-table-action btn-table-primary"
                      @click="processApproval(app, 'ผอ.ตลาด พิจารณาอนุมัติคำร้อง')"
                      title="ผอ.ตลาด อนุมัติ"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      <span>ผอ. อนุมัติ</span>
                    </button>

                    <!-- Step 3: Bureau Director Review -->
                    <button
                      v-else-if="app.status === 'pending_bureau_director'"
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
                      :to="`/contracts?new_for_stall=${encodeURIComponent(app.stall_number || app.stall_id || '')}&app_id=${app.id}`"
                      class="btn-table-action btn-table-primary"
                      title="จัดทำสัญญาเช่า"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
                      <span>ทำสัญญา</span>
                    </NuxtLink>
                    <span
                      v-else-if="app.status === 'approved' && !isStallLease(app)"
                      class="general-approved-pill"
                    >
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                      <span>อนุมัติแล้ว</span>
                    </span>

                    <!-- Reject / Request Revision -->
                    <button
                      v-if="['pending_officer', 'pending_market_director', 'pending_bureau_director', 'pending_debt_check', 'debt_pending'].includes(app.status)"
                      class="btn-table-action btn-table-cancel"
                      @click="openRejectModal(app)"
                      title="ระบุหมายเหตุไม่อนุมัติคำร้อง"
                    >
                      <span>ไม่อนุมัติ</span>
                    </button>

                    <!-- Clean Edit/Detail Icon per mockup -->
                    <button
                      class="btn-table-icon-edit"
                      title="ดูรายละเอียด/จัดการคำร้อง"
                      @click="openDetails(app)"
                    >
                      <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
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

        <!-- View 2: Cards Grid View -->
        <div v-else-if="viewMode === 'card'" class="application-cards-container">
          <div v-if="paginatedApplications.length > 0" class="application-cards-grid">
            <div
              v-for="app in paginatedApplications"
              :key="app.id"
              class="application-card"
              :class="`status-${app.status}`"
            >
              <!-- Card Top Row: Application Number & Created Date + Type & Status Badge -->
              <div class="app-card-header">
                <div class="app-card-header-left">
                  <div class="app-card-icon-box" :class="{ 'icon-box-general': !isStallLease(app) }">
                    <svg v-if="isStallLease(app)" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <rect x="3" y="3" width="7" height="7"></rect>
                      <rect x="14" y="3" width="7" height="7"></rect>
                      <rect x="14" y="14" width="7" height="7"></rect>
                      <rect x="3" y="14" width="7" height="7"></rect>
                    </svg>
                    <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                      <polyline points="14 2 14 8 20 8"></polyline>
                      <line x1="16" y1="13" x2="8" y2="13"></line>
                      <line x1="16" y1="17" x2="8" y2="17"></line>
                    </svg>
                  </div>
                  <div>
                    <div style="display: flex; align-items: center; gap: 6px; flex-wrap: wrap;">
                      <strong class="table-code-blue app-card-number">{{ app.application_number }}</strong>
                      <span v-if="isStallLease(app)" class="badge-type badge-type-stall">ขอเช่าแผง</span>
                      <span v-else class="badge-type badge-type-general">คำร้องทั่วไป</span>
                    </div>
                    <div class="app-card-date">
                      ยื่นเมื่อ {{ app.created_at?.slice(0, 10) }}
                    </div>
                  </div>
                </div>

                <div class="app-card-header-right">
                  <span class="status-pill-badge" :class="getStatusBadgeClass(app.status)">
                    {{ getStatusLabel(app.status) }}
                  </span>
                </div>
              </div>

              <!-- Action Actor Tag Subheader -->
              <div class="app-card-actor-row">
                <div class="action-actor-tag" :class="getActionActorInfo(app).tagClass">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                    <circle cx="12" cy="7" r="4"></circle>
                  </svg>
                  <span>{{ getActionActorInfo(app).actionText }}</span>
                </div>
              </div>

              <!-- Stall & Shop / Subject Banner -->
              <div class="app-card-shop-banner" :class="{ 'banner-general': !isStallLease(app) }">
                <template v-if="isStallLease(app)">
                  <div class="app-card-shop-info">
                    <div class="app-card-shop-name">{{ app.shop_name }}</div>
                    <div class="app-card-product-type">{{ app.product_type || '-' }}</div>
                  </div>
                  <div class="app-card-stall-badge">
                    <span class="badge badge-occupied">แผง {{ app.stall_number }}</span>
                    <span class="app-card-stall-rate">
                      ฿{{ Number(app.stall_monthly_rate).toLocaleString() }}/ด. ({{ app.requested_period_months }} ด.)
                    </span>
                  </div>
                </template>
                <template v-else>
                  <div class="app-card-shop-info">
                    <div class="app-card-shop-name">{{ app.subject || app.shop_name || 'คำร้องทั่วไป' }}</div>
                    <div class="app-card-product-type">{{ app.market_name || 'ตลาด กทม.' }}</div>
                  </div>
                  <div class="app-card-stall-badge">
                    <span v-if="app.stall_number && app.stall_number !== '-'" class="badge badge-occupied">แผง {{ app.stall_number }}</span>
                    <span v-else class="badge badge-general-pill">คำร้องทั่วไป</span>
                  </div>
                </template>
              </div>

              <!-- Applicant & Meta Info Grid -->
              <div class="app-card-meta-grid">
                <div class="app-card-meta-item">
                  <span class="meta-label">ผู้ยื่นคำร้อง</span>
                  <div class="meta-val-primary">{{ app.applicant_name }}</div>
                  <div class="meta-val-sub">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
                    </svg>
                    <span>{{ app.applicant_phone || '-' }}</span>
                  </div>
                </div>

                <div class="app-card-meta-item">
                  <span class="meta-label">วันที่ยื่น / เริ่ม</span>
                  <div class="meta-val-primary meta-val-date">
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                      <line x1="16" y1="2" x2="16" y2="6"></line>
                      <line x1="8" y1="2" x2="8" y2="6"></line>
                      <line x1="3" y1="10" x2="21" y2="10"></line>
                    </svg>
                    <span>{{ app.requested_start_date || app.created_at?.slice(0, 10) }}</span>
                  </div>
                </div>
              </div>

              <!-- Review Note (if any) -->
              <div v-if="getReviewNoteInfo(app).text" class="app-card-note-section">
                <div class="review-comment-pill" :class="`comment-${getReviewNoteInfo(app).type}`">
                  <div class="comment-header-row">
                    <span class="comment-role-label">{{ getReviewNoteInfo(app).role }}:</span>
                  </div>
                  <div class="comment-text-body">
                    {{ getReviewNoteInfo(app).text }}
                  </div>
                </div>
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

              <!-- Card Action Footer -->
              <div class="app-card-actions">
                <!-- Step 1: Officer Review -->
                <button
                  v-if="app.status === 'pending_officer'"
                  class="btn-table-action btn-table-primary"
                  @click="processApproval(app, 'เจ้าหน้าที่ตรวจสอบเอกสารครบถ้วน ส่งต่อ ผอ.ตลาด')"
                  title="เจ้าหน้าที่ตรวจสอบ"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                  <span>ตรวจแล้ว</span>
                </button>

                <!-- Step 2: Market Director Review -->
                <button
                  v-else-if="['pending_market_director', 'pending_debt_check', 'debt_pending'].includes(app.status)"
                  class="btn-table-action btn-table-primary"
                  @click="processApproval(app, 'ผอ.ตลาด พิจารณาอนุมัติคำร้อง')"
                  title="ผอ.ตลาด อนุมัติ"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                  <span>ผอ. อนุมัติ</span>
                </button>

                <!-- Step 3: Bureau Director Review -->
                <button
                  v-else-if="app.status === 'pending_bureau_director'"
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
                  :to="`/contracts?new_for_stall=${encodeURIComponent(app.stall_number || app.stall_id || '')}&app_id=${app.id}`"
                  class="btn-table-action btn-table-primary"
                  title="จัดทำสัญญาเช่า"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
                  <span>ทำสัญญา</span>
                </NuxtLink>
                <span
                  v-else-if="app.status === 'approved' && !isStallLease(app)"
                  class="general-approved-pill"
                >
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                  <span>อนุมัติแล้ว</span>
                </span>

                <!-- Reject / Request Revision -->
                <button
                  v-if="['pending_officer', 'pending_market_director', 'pending_bureau_director', 'pending_debt_check', 'debt_pending'].includes(app.status)"
                  class="btn-table-action btn-table-cancel"
                  @click="openRejectModal(app)"
                  title="ระบุหมายเหตุไม่อนุมัติคำร้อง"
                >
                  <span>ไม่อนุมัติ</span>
                </button>

                <!-- Detail Button -->
                <button
                  class="btn-card-detail"
                  title="ดูรายละเอียด/จัดการคำร้อง"
                  @click="openDetails(app)"
                >
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="3"></circle>
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                  </svg>
                  <span>รายละเอียด</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Empty State in Card View -->
          <div v-else class="empty-cards-state">
            <svg width="44" height="44" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <line x1="10" y1="9" x2="8" y2="9"></line>
            </svg>
            <div class="empty-cards-title">ไม่พบรายการคำร้องขอเช่าตามเงื่อนไขที่ค้นหา</div>
          </div>
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
              <h3 class="modal-title-danger">ระบุหมายเหตุไม่อนุมัติคำร้อง</h3>
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
            <span>เมื่อยืนยันไม่อนุมัติ ระบบจะบันทึกสถานะเป็น "ไม่อนุมัติ" และจัดเก็บหมายเหตุนี้เพื่อนำไปแสดงในคอลัมน์หมายเหตุการพิจารณา</span>
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
              หมายเหตุ / เหตุผลการไม่อนุมัติ <span style="color: #dc2626;">*</span>
            </label>
            <textarea
              id="reject-comment-input"
              v-model="rejectComment"
              class="reject-textarea"
              :class="{ 'has-error': rejectErrorMsg }"
              rows="3"
              placeholder="กรอกเหตุผลหรือคำอธิบายที่ไม่อนุมัติคำร้องขอเช่าแผง..."
              @input="rejectErrorMsg = ''"
            ></textarea>
            <div v-if="rejectErrorMsg" class="reject-error-text">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              <span>{{ rejectErrorMsg }}</span>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="modal-footer">
          <AppButton
            variant="secondary"
            size="md"
            :disabled="rejectSubmitting"
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
            ยืนยันไม่อนุมัติคำร้อง
          </AppButton>
        </div>
      </div>
    </div>

    <!-- Modal: ยื่นคำร้องทั่วไป (Create General Application) -->
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
          <div v-if="generalFormError" class="general-alert-error">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
            <span>{{ generalFormError }}</span>
          </div>

          <!-- Subject Preset Selector -->
          <div class="form-group">
            <div class="preset-label-header">
              <label class="form-label" style="margin-bottom: 0;">ประเภทคำร้อง / เรื่องที่ขอ <span class="text-danger">*</span></label>
              <NuxtLink
                v-if="isAdminOrStaff"
                to="/applications/types"
                class="btn-manage-types-shortcut"
                title="จัดการรายการประเภทคำร้องในระบบ (เฉพาะผู้ดูแลระบบและเจ้าหน้าที่)"
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
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
                <span>จัดการประเภท</span>
              </NuxtLink>
            </div>
            <div class="general-preset-grid">
              <button
                v-for="sub in generalPresetSubjects"
                :key="sub"
                type="button"
                class="general-preset-btn"
                :class="{ active: generalForm.presetSubject === sub }"
                @click="generalForm.presetSubject = sub"
              >
                {{ sub }}
              </button>
            </div>
          </div>

          <!-- Custom Subject Input if "เรื่องอื่นๆ" -->
          <div v-if="generalForm.presetSubject === 'เรื่องอื่นๆ'" class="form-group">
            <label class="form-label" for="custom-subject">ระบุหัวข้อคำร้อง <span class="text-danger">*</span></label>
            <input
              id="custom-subject"
              v-model="generalForm.customSubject"
              type="text"
              class="form-input"
              placeholder="เช่น ขอติดตั้งมิเตอร์ไฟฟ้าชั่วคราว..."
            />
          </div>

          <!-- Market Selection & Stall Number (optional) -->
          <div class="form-row-2col">
            <div class="form-group">
              <label class="form-label">ตลาดที่เกี่ยวข้อง <span class="text-danger">*</span></label>
              <AppSelect
                v-model="generalForm.market_id"
                :options="marketOptions"
                placeholder="เลือกตลาด..."
              />
            </div>
            <div class="form-group">
              <label class="form-label">หมายเลขแผงค้า (ถ้ามี)</label>
              <input
                v-model="generalForm.stall_number"
                type="text"
                class="form-input"
                placeholder="เช่น A-01, JJ-19-122 หรือเว้นว่าง"
              />
            </div>
          </div>

          <!-- Applicant Name & Phone -->
          <div class="form-row-2col">
            <div class="form-group">
              <label class="form-label">ชื่อ-นามสกุล ผู้ยื่นคำร้อง <span class="text-danger">*</span></label>
              <input
                v-model="generalForm.applicant_name"
                type="text"
                class="form-input"
                :disabled="isTenant"
                :class="{ 'input-disabled': isTenant }"
                placeholder="ชื่อ-นามสกุล..."
              />
            </div>
            <div class="form-group">
              <label class="form-label">เบอร์โทรศัพท์ติดต่อ <span class="text-danger">*</span></label>
              <input
                v-model="generalForm.applicant_phone"
                type="tel"
                class="form-input"
                :disabled="isTenant"
                :class="{ 'input-disabled': isTenant }"
                placeholder="08X-XXX-XXXX"
              />
            </div>
          </div>

          <!-- ID Card & Shop Name -->
          <div class="form-row-2col">
            <div class="form-group">
              <label class="form-label">เลขประจำตัวประชาชน (13 หลัก)</label>
              <input
                v-model="generalForm.applicant_id_card"
                type="text"
                class="form-input"
                maxlength="13"
                placeholder="1XXXXXXXXXXXX"
              />
            </div>
            <div class="form-group">
              <label class="form-label">ชื่อร้านค้า / กิจการ</label>
              <input
                v-model="generalForm.shop_name"
                type="text"
                class="form-input"
                placeholder="ชื่อร้านหรือชื่อกิจการ..."
              />
            </div>
          </div>

          <!-- Details Textarea -->
          <div class="form-group">
            <label class="form-label">รายละเอียดคำร้อง / วัตถุประสงค์ <span class="text-danger">*</span></label>
            <textarea
              v-model="generalForm.details"
              class="form-textarea"
              rows="3"
              placeholder="ระบุเหตุผล ความจำเป็น หรือข้อความคำร้องที่ต้องการให้เจ้าหน้าที่พิจารณา..."
            ></textarea>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="modal-footer">
          <AppButton
            variant="secondary"
            size="md"
            :disabled="generalSubmitting"
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
            ส่งคำร้องทั่วไป
          </AppButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'
import { useNotifications } from '~/composables/useNotifications'

const { get, post } = useApi()
const authStore = useAuthStore()
const { user, currentUser, isTenant, isAdmin, isStaff, role } = useAuth()
const isAdminOrStaff = computed(() => isAdmin.value || isStaff.value || ['admin', 'staff'].includes(role.value || authStore.role || ''))
const { showAlert, showConfirm, showSuccess, showError } = useAppAlert()
const { fetchNotifications } = useNotifications()
const { resetTestData } = useResetTestData()

const handleClearAllData = async () => {
  await resetTestData(async () => {
    await loadApplications()
  })
}

const loading = ref(true)
const applications = ref<any[]>([])
const marketsList = ref<any[]>([])
const searchQuery = ref('')
const activeTab = ref('all')
const typeFilter = ref<'all' | 'stall_lease' | 'general'>('all')
const currentPage = ref(1)
const pageSize = ref(10)

const viewMode = ref<'table' | 'card'>('table')
const viewModeOptions = [
  { value: 'table', label: 'ตาราง', icon: 'table' as const },
  { value: 'card', label: 'การ์ด', icon: 'card' as const }
]

// แยกแยะประเภทคำร้อง: คำร้องขอเช่าแผง vs คำร้องทั่วไป
const isStallLease = (app: any): boolean => {
  if (app.application_type === 'stall_lease') return true
  if (app.application_type === 'general') return false
  if (app.application_number?.startsWith('GEN-') || app.application_no?.startsWith('GEN-')) return false
  return Boolean(app.booking_id || (app.stall_number && app.stall_number !== '-') || (app.stall_ident && app.stall_ident !== '-'))
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

const countByStatus = (status: string) => {
  let list = Array.isArray(applications.value) ? applications.value : []
  if (typeFilter.value === 'stall_lease') {
    list = list.filter(a => isStallLease(a))
  } else if (typeFilter.value === 'general') {
    list = list.filter(a => isGeneralApp(a))
  }
  if (status === 'pending_market_director') {
    return list.filter(a => a.status === 'pending_market_director' || a.status === 'pending_debt_check' || a.status === 'debt_pending').length
  }
  return list.filter(a => a.status === status).length
}

const filteredApplications = computed(() => {
  let list = Array.isArray(applications.value) ? applications.value : []

  // Filter by Type
  if (typeFilter.value === 'stall_lease') {
    list = list.filter(a => isStallLease(a))
  } else if (typeFilter.value === 'general') {
    list = list.filter(a => isGeneralApp(a))
  }

  // Filter by Status Tab
  if (activeTab.value !== 'all') {
    if (activeTab.value === 'rejected') {
      list = list.filter(a => a.status === 'rejected' || a.status === 'revision_required')
    } else if (activeTab.value === 'pending_market_director') {
      list = list.filter(a => a.status === 'pending_market_director' || a.status === 'pending_debt_check' || a.status === 'debt_pending')
    } else {
      list = list.filter(a => a.status === activeTab.value)
    }
  }

  // Search
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(a =>
      (a.application_number && a.application_number.toLowerCase().includes(q)) ||
      (a.applicant_name && a.applicant_name.toLowerCase().includes(q)) ||
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

// Helper: ข้อมูลผู้ดำเนินการ / ใครเป็นคน Action สอดคล้องกับ Status
// ด้านล่างบอกว่าใครคือคน Action / กรณีที่เป็นการรอ ให้บอกว่า "รอใคร"
const getActionActorInfo = (app: any) => {
  const status = app.status || app.application_state
  switch (status) {
    case 'pending_officer':
      return {
        role: 'เจ้าหน้าที่ตลาด',
        actionText: 'จากเจ้าหน้าที่ตลาด',
        tagClass: 'actor-tag-officer-pending'
      }
    case 'pending_debt_check':
      return {
        role: 'เจ้าหน้าที่การคลัง/ตลาด',
        actionText: 'รอ ตรวจสอบหนี้ค้าง',
        tagClass: 'actor-tag-officer'
      }
    case 'debt_pending':
      return {
        role: 'ผอ.ตลาด / การคลัง',
        actionText: 'รอ เคลียร์หนี้ค้างชำระ',
        tagClass: 'actor-tag-warning'
      }
    case 'pending_market_director':
      return {
        role: 'ผอ.ตลาด',
        actionText: 'รอ ผอ.ตลาด',
        tagClass: 'actor-tag-director'
      }
    case 'pending_bureau_director':
      return {
        role: 'ผอ.สำนักตลาด',
        actionText: 'รอ ผอ.สำนักตลาด',
        tagClass: 'actor-tag-director'
      }
    case 'approved':
      return {
        role: 'ผอ.สำนักตลาด',
        actionText: 'ผอ.สำนักตลาด',
        tagClass: 'actor-tag-success'
      }
    case 'revision_required':
      return {
        role: 'เจ้าหน้าที่ตลาด',
        actionText: 'เจ้าหน้าที่ตลาด',
        tagClass: 'actor-tag-warning'
      }
    case 'rejected': {
      let actor = 'ผู้มีอำนาจพิจารณา'
      if (app.bureau_director_comment) {
        actor = 'ผอ.สำนักตลาด'
      } else if (app.market_director_comment) {
        actor = 'ผอ.ตลาด'
      } else if (app.officer_comment) {
        actor = 'เจ้าหน้าที่ตลาด'
      }
      return {
        role: actor,
        actionText: actor,
        tagClass: 'actor-tag-danger'
      }
    }
    default:
      return {
        role: '-',
        actionText: '-',
        tagClass: 'actor-tag-default'
      }
  }
}

// Helper: ดึงข้อความหมายเหตุการพิจารณา (officer_comment, market_director_comment, review_note) สอดคล้องกับว่าใครเป็นคน Action
const getReviewNoteInfo = (app: any) => {
  const status = app.status || app.application_state

  if (status === 'pending_officer') {
    if (app.review_note) {
      return { role: 'คำร้อง', text: app.review_note, type: 'info' }
    }
    return { role: '', text: '', type: 'empty' }
  }

  if (status === 'pending_debt_check') {
    return {
      role: 'จนท.',
      text: app.officer_comment || app.review_note || '',
      type: 'officer'
    }
  }

  if (status === 'debt_pending') {
    return {
      role: 'การคลัง / ผอ.ตลาด',
      text: app.market_director_comment || app.review_note || '',
      type: 'warning'
    }
  }

  if (status === 'pending_market_director') {
    return {
      role: 'จนท.',
      text: app.officer_comment || app.review_note || '',
      type: 'officer'
    }
  }

  if (status === 'pending_bureau_director') {
    return {
      role: 'ผอ.ตลาด',
      text: app.market_director_comment || app.officer_comment || app.review_note || '',
      type: 'director'
    }
  }

  if (status === 'approved') {
    return {
      role: 'ผอ.สำนัก',
      text: app.bureau_director_comment || app.market_director_comment || app.officer_comment || app.review_note || 'อนุมัติเรียบร้อย',
      type: 'success'
    }
  }

  if (status === 'revision_required') {
    return {
      role: 'สิ่งที่ต้องแก้ไข',
      text: app.officer_comment || app.review_note || '',
      type: 'warning'
    }
  }

  if (status === 'rejected') {
    return {
      role: 'เหตุผล',
      text: app.review_note || app.bureau_director_comment || app.market_director_comment || app.officer_comment || 'ไม่อนุมัติคำร้อง',
      type: 'danger'
    }
  }

  const fallback = app.bureau_director_comment || app.market_director_comment || app.officer_comment || app.review_note
  if (fallback) {
    return { role: 'หมายเหตุ', text: fallback, type: 'info' }
  }
  return { role: '', text: '', type: 'empty' }
}

// Helper: รวบรวมประวัติความเห็นทั้งหมดทุกขั้นตอน
const getAllCommentsList = (app: any) => {
  const list: { role: string; text: string; step: string }[] = []
  if (app.officer_comment) {
    list.push({ role: 'เจ้าหน้าที่ตลาด', text: app.officer_comment, step: 'ขั้นตรวจสอบเอกสาร' })
  }
  if (app.market_director_comment) {
    list.push({ role: 'ผอ.ตลาด', text: app.market_director_comment, step: 'ขั้นพิจารณา ผอ.ตลาด' })
  }
  if (app.bureau_director_comment) {
    list.push({ role: 'ผอ.สำนักตลาด', text: app.bureau_director_comment, step: 'ขั้นพิจารณา ผอ.สำนัก' })
  }
  if (app.review_note && !list.some(l => l.text === app.review_note)) {
    list.push({ role: 'บันทึกคำร้อง', text: app.review_note, step: 'หมายเหตุคำร้อง' })
  }
  return list
}

const processApproval = async (app: any, defaultComment: string) => {
  const confirmed = await showConfirm(
    `ยืนยันการอนุมัติคำร้อง ${app.application_number} (ผู้ยื่น: ${app.applicant_name}) ใช่หรือไม่?`,
    {
      title: 'ยืนยันพิจารณาอนุมัติ',
      type: 'info',
      confirmText: 'อนุมัติคำร้อง',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return

  try {
    await post(`/applications/${app.id}/action`, {
      action: 'approve',
      comment: defaultComment
    })
    await loadApplications()
    await showSuccess('บันทึกผลการอนุมัติเรียบร้อยแล้ว', 'อนุมัติสำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาดในการอนุมัติ', 'ข้อผิดพลาด')
  }
}

// State สำหรับโมดอลระบุหมายเหตุไม่อนุมัติคำร้อง
const showRejectModal = ref(false)
const targetAppForReject = ref<any>(null)
const rejectComment = ref('')
const rejectSubmitting = ref(false)
const rejectErrorMsg = ref('')

const quickRejectReasons = [
  'เอกสารหลักฐานไม่ครบถ้วนหรือไม่ชัดเจน',
  'คุณสมบัติไม่ผ่านเกณฑ์การพิจารณาของตลาด',
  'ไม่สามารถติดต่อผู้ยื่นคำร้องได้ตามกำหนด',
  'มีประวัติค้างชำระหนี้สินกับตลาด',
  'ประเภทสินค้าไม่สอดคล้องกับโซนการค้าของแผง'
]

const openRejectModal = (app: any) => {
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
    rejectComment.value = `${rejectComment.value} / ${reason}`
  }
  rejectErrorMsg.value = ''
}

const submitReject = async () => {
  const comment = rejectComment.value.trim()
  if (!comment) {
    rejectErrorMsg.value = 'กรุณาระบุหมายเหตุหรือเหตุผลในการไม่อนุมัติคำร้อง'
    return
  }

  if (!targetAppForReject.value) return

  rejectSubmitting.value = true
  rejectErrorMsg.value = ''

  try {
    await post(`/applications/${targetAppForReject.value.id}/action`, {
      action: 'reject',
      comment: comment,
      reason: comment
    })
    await loadApplications()
    closeRejectModal()
    await showSuccess('บันทึกผลการไม่อนุมัติคำร้องเรียบร้อยแล้ว', 'ไม่อนุมัติสำเร็จ')
  } catch (err: any) {
    rejectErrorMsg.value = err?.data?.detail || err?.data?.message || 'เกิดข้อผิดพลาดในการบันทึกไม่อนุมัติ'
  } finally {
    rejectSubmitting.value = false
  }
}

// State สำหรับโมดอลยื่นคำร้องทั่วไป
const showGeneralModal = ref(false)
const generalSubmitting = ref(false)
const generalFormError = ref('')

// รายการประเภทคำร้อง (ดึงมาจาก API /application-types โดยตรง)
const defaultFallbackSubjects = [
  'ขออนุญาตปรับปรุง/ต่อเติมแผงค้า',
  'ขอเปลี่ยน/เพิ่มประเภทสินค้าที่จำหน่าย',
  'ขอออกหนังสือรับรองการเช่าแผงค้า',
  'ขอผ่อนผันเวลาเปิด-ปิดแผงค้า',
  'แจ้งปัญหา / ร้องเรียนการใช้งานพื้นที่ตลาด',
  'เรื่องอื่นๆ'
]

const dynamicApplicationTypes = ref<any[]>([])

const loadApplicationTypes = async () => {
  try {
    const res = await get<any>('/application-types?active_only=true')
    if (res.status && res.data && Array.isArray(res.data) && res.data.length > 0) {
      dynamicApplicationTypes.value = res.data
    }
  } catch {
    // ใช้ค่า fallback หากไม่สามารถเชื่อมต่อได้
  }
}

const generalPresetSubjects = computed(() => {
  if (dynamicApplicationTypes.value.length > 0) {
    return dynamicApplicationTypes.value.map((t: any) => t.title)
  }
  return defaultFallbackSubjects
})

const generalForm = ref({
  presetSubject: 'ขออนุญาตปรับปรุง/ต่อเติมแผงค้า',
  customSubject: '',
  market_id: 1,
  market_name: 'ตลาดนัดจตุจักร',
  stall_number: '',
  applicant_name: '',
  applicant_phone: '',
  applicant_id_card: '',
  applicant_email: '',
  shop_name: '',
  details: ''
})

// Helper: ดึงข้อมูลโปรไฟล์ผู้ใช้งานที่ล็อกอินอยู่ในปัจจุบัน (เฉพาะกรณี role: tenant)
const getCurrentUserProfile = () => {
  if (typeof window !== 'undefined') {
    authStore.hydrate()
  }

  // ตรวจสอบ role: ถ้าไม่ใช่ tenant (ผู้เช่า) ให้ส่งค่าว่างเพื่อให้เจ้าหน้าที่/ผู้ใช้กรอกเอง
  const activeRole = authStore.role || authStore.currentUser?.role || (typeof window !== 'undefined' ? localStorage.getItem('user_role') : '')
  const isTenantUser = activeRole === 'tenant' || authStore.isTenant || isTenant.value
  if (!isTenantUser) {
    return {
      name: '',
      phone: '',
      idCard: '',
      email: ''
    }
  }

  // ถ้าเป็น role tenant ให้ดึงข้อมูลของผู้เช่ามาแสดงอัตโนมัติ
  let profile: any = authStore.currentUser || currentUser.value || user.value
  if (!profile && typeof window !== 'undefined') {
    try {
      const raw = localStorage.getItem('user_info') || localStorage.getItem('emarket_user')
      if (raw) profile = JSON.parse(raw)
    } catch {}
  }

  const rawName = (profile?.full_name || profile?.name || profile?.fullName || '').trim()
  const cleanName = rawName.replace(/\s*\((ผู้เช่า|tenant)\)\s*/gi, '').trim()
  const rawPhone = (profile?.phone || profile?.phone_number || profile?.tel || profile?.mobile || '').trim()

  const finalName = cleanName || rawName || ''
  const finalPhone = rawPhone || ''
  const idCard = profile?.citizen_id || profile?.id_card || ''
  const email = profile?.email || ''

  return {
    name: finalName,
    phone: finalPhone,
    idCard,
    email
  }
}

const openCreateGeneralModal = () => {
  loadApplicationTypes()
  const profile = getCurrentUserProfile()
  const initialSubject = generalPresetSubjects.value[0] || 'ขออนุญาตปรับปรุง/ต่อเติมแผงค้า'
  generalForm.value = {
    presetSubject: initialSubject,
    customSubject: '',
    market_id: marketsList.value[0]?.id || 1,
    market_name: marketsList.value[0]?.name || 'ตลาดนัดจตุจักร',
    stall_number: '',
    applicant_name: profile.name,
    applicant_phone: profile.phone,
    applicant_id_card: profile.idCard,
    applicant_email: profile.email,
    shop_name: '',
    details: ''
  }
  generalFormError.value = ''
  showGeneralModal.value = true
}

// ตรวจสอบและดึงข้อมูลโปรไฟล์ทันทีเมื่อโมดอลเปิดขึ้นมา (เฉพาะกรณี role: tenant)
watch(showGeneralModal, (isOpen) => {
  if (isOpen) {
    const activeRole = authStore.role || authStore.currentUser?.role || (typeof window !== 'undefined' ? localStorage.getItem('user_role') : '')
    const isTenantUser = activeRole === 'tenant' || authStore.isTenant || isTenant.value
    if (isTenantUser) {
      const profile = getCurrentUserProfile()
      if (!generalForm.value.applicant_name) {
        generalForm.value.applicant_name = profile.name
      }
      if (!generalForm.value.applicant_phone) {
        generalForm.value.applicant_phone = profile.phone
      }
      if (!generalForm.value.applicant_id_card && profile.idCard) {
        generalForm.value.applicant_id_card = profile.idCard
      }
      if (!generalForm.value.applicant_email && profile.email) {
        generalForm.value.applicant_email = profile.email
      }
    }
  }
})

const closeCreateGeneralModal = () => {
  showGeneralModal.value = false
  generalFormError.value = ''
}

const submitGeneralApplication = async () => {
  generalFormError.value = ''
  const form = generalForm.value
  const subject = form.presetSubject === 'เรื่องอื่นๆ' ? form.customSubject.trim() : form.presetSubject
  if (!subject) {
    generalFormError.value = 'กรุณาระบุหัวข้อคำร้อง'
    return
  }
  if (!form.applicant_name.trim() || !form.applicant_phone.trim()) {
    generalFormError.value = 'กรุณาระบุชื่อ-นามสกุล และเบอร์โทรศัพท์ผู้ยื่นคำร้อง'
    return
  }
  if (!form.details.trim()) {
    generalFormError.value = 'กรุณาระบุรายละเอียดคำร้อง / วัตถุประสงค์'
    return
  }

  generalSubmitting.value = true
  try {
    const selectedMarket = marketsList.value.find(m => m.id === form.market_id)
    const marketName = selectedMarket?.name || form.market_name || 'ตลาด กทม.'

    await post('/applications', {
      application_type: 'general',
      subject: subject,
      shop_name: form.shop_name.trim() || subject,
      product_category: 'คำร้องทั่วไป',
      product_type: 'คำร้องทั่วไป',
      stall_ident: form.stall_number.trim() || '-',
      stall_number: form.stall_number.trim() || '-',
      market_id: form.market_id,
      market_name: marketName,
      applicant_name: form.applicant_name.trim(),
      applicant_phone: form.applicant_phone.trim(),
      applicant_id_card: form.applicant_id_card.trim(),
      applicant_email: form.applicant_email.trim(),
      applicant_address: '-',
      details: form.details.trim(),
      review_note: form.details.trim(),
      lease_term_months: 1,
      requested_period_months: 1,
      monthly_rate: 0
    })

    await loadApplications()
    closeCreateGeneralModal()
    await showSuccess(`ยื่นคำร้องทั่วไปเรื่อง "${subject}" เรียบร้อยแล้ว`, 'บันทึกสำเร็จ')
  } catch (err: any) {
    generalFormError.value = err?.data?.detail || err?.data?.error?.message || err?.data?.message || err?.message || 'เกิดข้อผิดพลาดในการยื่นคำร้อง'
  } finally {
    generalSubmitting.value = false
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
  const idCardText = app.applicant_id_card ? `\nเลขประจำตัว ปชช.: ${app.applicant_id_card}` : ''
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
    `ประเภท: ${typeText}\nเลขที่คำร้อง: ${app.application_number || app.application_no}${marketText}${stallText}${subjectText}\nผู้ยื่นคำร้อง: ${app.applicant_name} (โทร. ${app.applicant_phone})${idCardText}${emailText}${addressText}\nสถานะ: ${getStatusLabel(app.status)}\nขั้นตอน: ${getActionActorInfo(app).actionText}${detailText}${commentDetails}`,
    { title: 'รายละเอียดคำร้อง', type: 'info' }
  )
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
      `"${getStatusLabel(app.status)}"`,
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

const getStatusBadgeClass = (status: string) => {
  switch (status) {
    case 'pending_officer':
    case 'pending_debt_check':
    case 'pending_market_director':
    case 'pending_bureau_director':
      return 'badge-soft-blue'
    case 'debt_pending':
      return 'badge-soft-warning'
    case 'approved':
      return 'badge-soft-success'
    case 'revision_required':
      return 'badge-soft-warning'
    case 'rejected':
      return 'badge-soft-danger'
    default:
      return 'badge-soft-gray'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'pending_officer':
    case 'pending_debt_check':
    case 'debt_pending':
    case 'pending_market_director':
    case 'pending_bureau_director':
      return 'รอตรวจสอบ'
    case 'approved':
      return 'อนุมัติ'
    case 'revision_required':
      return 'ส่งกลับแก้ไข'
    case 'rejected':
      return 'ไม่อนุมัติ'
    default:
      return status || '-'
  }
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    authStore.hydrate()
    const saved = localStorage.getItem('applications_view_mode')
    if (saved === 'table' || saved === 'card') {
      viewMode.value = saved
    }
  }
  loadApplications()
  loadMarkets()
  loadApplicationTypes()
})

watch(viewMode, (newVal) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('applications_view_mode', newVal)
  }
})
</script>

<style scoped>
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
</style>
