<template>
  <div class="booking-detail-page">
    <AppHeader />

    <main class="content-body detail-container">
      <AppBreadcrumb
        :items="returnCrumbItems"
        :current="booking ? `รายละเอียด ${formatBookingNumber(booking.booking_number)}` : 'รายละเอียดการจอง'"
        :back-to="returnBackTo"
      />

      <div v-if="loading" class="detail-skeleton" aria-label="กำลังโหลดรายละเอียดการจอง">
        <div class="skeleton-line wide"></div>
        <div class="skeleton-line short"></div>
        <div class="skeleton-panel"></div>
        <div class="skeleton-grid"><div></div><div></div><div></div></div>
      </div>

      <div v-else-if="!booking" class="detail-empty">
        <svg width="38" height="38" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
        <h1>ไม่พบรายการจอง</h1>
        <p>รายการนี้อาจถูกลบหรือหมายเลขการจองไม่ถูกต้อง</p>
      </div>

      <template v-else>
        <section class="detail-hero">
          <div class="detail-title-wrap">
            <div class="detail-title-icon">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            </div>
            <div>
              <div class="eyebrow">รายละเอียดการจองดูแผงค้า</div>
              <h1>{{ formatBookingNumber(booking.booking_number) }}</h1>
              <div class="created-at">วันที่จอง {{ formatThaiDateTime(booking.created_at) }}</div>
            </div>
          </div>

          <div class="hero-actions">
            <span
              v-if="isQueueBooking(booking)"
              class="status-badge tone-warning"
              title="ลำดับคิวสำรอง"
            >
              คิวสำรองที่ {{ booking.queue_position || booking.queue_number || 1 }}
            </span>
            <span class="status-badge" :class="statusToneClass">{{ statusLabel }}</span>
          </div>
        </section>

        <section class="next-action-card" :class="`next-action-${nextAction.tone}`" aria-labelledby="next-action-title">
          <div class="next-action-icon">
            <svg v-if="nextAction.tone === 'waiting'" width="21" height="21" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"></circle><polyline points="12 7 12 12 15 14"></polyline></svg>
            <svg v-else-if="nextAction.tone === 'complete'" width="21" height="21" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
            <svg v-else width="21" height="21" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"></path><path d="m13 6 6 6-6 6"></path></svg>
          </div>
          <div class="next-action-copy">
            <div class="next-action-meta">
              <span>ขั้นตอนถัดไป</span>
              <span class="next-action-owner">{{ nextAction.owner }}</span>
            </div>
            <h2 id="next-action-title">{{ nextAction.title }}</h2>
            <p>{{ nextAction.description }}</p>
          </div>
          <div v-if="hasActionButtons" class="next-action-buttons">
            <!-- ============================================================ -->
            <!-- OFFICER ROLE ACTION BUTTONS (ตามหน้ารายการจอง 100%)          -->
            <!-- ============================================================ -->
            <template v-if="currentRole === 'officer'">
              <!-- 1. รอเจ้าหน้าที่อนุมัติ: อนุมัตินัด + ยกเลิก/ยกเลิกคิว -->
              <template v-if="booking.status === 'pending_approval' || booking.status === 'pending'">
                <AppButton
                  variant="primary"
                  size="md"
                  :disabled="actionLoading"
                  @click="openApproveModal"
                  title="อนุมัติวันเวลานัดหมาย"
                >
                  <template #prefix>
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                  </template>
                  <span>อนุมัตินัด</span>
                </AppButton>
                <AppButton
                  variant="danger-outline"
                  size="md"
                  :disabled="actionLoading"
                  @click="isQueueBooking(booking) ? handleOfficerCancel() : openRejectModal()"
                  :title="isQueueBooking(booking) ? 'ยกเลิกการต่อคิว' : 'ปฏิเสธ/ยกเลิกคำขอนัดหมาย'"
                >
                  <template #prefix>
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                  </template>
                  <span>{{ isQueueBooking(booking) ? 'ยกเลิกคิว' : 'ยกเลิก' }}</span>
                </AppButton>
              </template>

              <!-- 2. นัดหมายสำเร็จ: ยืนยันการเช่า (เฉพาะผู้จองหลัก) + ยกเลิก/ยกเลิกคิว -->
              <template v-else-if="booking.status === 'approved'">
                <AppButton
                  v-if="!isQueueBooking(booking)"
                  variant="success"
                  size="md"
                  :disabled="actionLoading"
                  @click="handleMarkVisited"
                  title="บันทึกผลว่าผู้เช่ามาชมแผงแล้ว และสนใจเช่า"
                >
                  <template #prefix>
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                  </template>
                  <span>ยืนยันการเช่า</span>
                </AppButton>
                <AppButton
                  variant="danger-outline"
                  size="md"
                  :disabled="actionLoading"
                  @click="handleOfficerCancel"
                  :title="isQueueBooking(booking) ? 'ยกเลิกการต่อคิว' : 'ยกเลิกนัดหมายนี้'"
                >
                  <template #prefix>
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                  </template>
                  <span>{{ isQueueBooking(booking) ? 'ยกเลิกคิว' : 'ยกเลิก' }}</span>
                </AppButton>
              </template>

              <!-- 3. รอผู้เช่ายื่นคำร้อง / รอเลือกวันนัดชม: ยกเลิก -->
              <template v-else-if="['interested', 'pending_decision', 'pending_slot'].includes(booking.status)">
                <AppButton
                  variant="danger-outline"
                  size="md"
                  :disabled="actionLoading"
                  @click="handleOfficerCancel"
                  title="ยกเลิกการจอง / บันทึกการสละสิทธิ์"
                >
                  <template #prefix>
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                  </template>
                  <span>ยกเลิก</span>
                </AppButton>
              </template>

              <!-- 4. รอคิว: ยกเลิกคิว -->
              <template v-else-if="booking.status === 'waiting_queue'">
                <AppButton
                  variant="danger-outline"
                  size="md"
                  :disabled="actionLoading"
                  @click="handleOfficerCancel"
                  title="ยกเลิกการต่อคิว"
                >
                  <template #prefix>
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                  </template>
                  <span>ยกเลิกคิว</span>
                </AppButton>
              </template>

              <!-- 5. ยื่นคำร้องสำเร็จ: ดูคำร้อง -->
              <template v-else-if="booking.status === 'lease_applied' || booking.status === 'completed'">
                <AppButton
                  to="/applications"
                  variant="primary"
                  size="md"
                >
                  <span>ดูคำร้อง</span>
                  <template #suffix>
                    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
                  </template>
                </AppButton>
              </template>
            </template>

            <!-- ============================================================ -->
            <!-- TENANT ROLE ACTION BUTTONS (ตามหน้ารายการจอง 100%)           -->
            <!-- ============================================================ -->
            <template v-else>
              <!-- 1. ปุ่มยกเลิก / สละสิทธิ์ -->
              <AppButton
                v-if="isQueueBooking(booking) && ['pending_approval', 'pending', 'approved', 'waiting_queue'].includes(booking.status)"
                variant="danger-outline"
                size="md"
                :disabled="actionLoading"
                @click="handleTenantCancel('ผู้เช่ายกเลิกการต่อคิว')"
                title="สละสิทธิ์การต่อคิว"
              >
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                </template>
                <span>สละสิทธิ์</span>
              </AppButton>

              <AppButton
                v-else-if="['pending_approval', 'pending', 'pending_slot', 'approved', 'interested', 'pending_decision'].includes(booking.status)"
                variant="danger-outline"
                size="md"
                :disabled="actionLoading"
                @click="handleTenantCancel('ผู้เช่าขอยกเลิกการจอง')"
                title="ยกเลิกคำขอจอง"
              >
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                </template>
                <span>ยกเลิก</span>
              </AppButton>

              <!-- 2. ปุ่ม Primary Action ตามสถานะ -->
              <AppButton
                v-if="booking.status === 'pending_slot'"
                variant="primary"
                size="md"
                :disabled="actionLoading"
                @click="openSelectSlotModal"
              >
                <span>เลือกวันนัดชม</span>
                <template #suffix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"></polyline></svg>
                </template>
              </AppButton>

              <AppButton
                v-else-if="['interested', 'pending_decision'].includes(booking.status)"
                variant="primary"
                size="md"
                :disabled="actionLoading"
                @click="openLeaseForSubBooking(booking)"
              >
                <span>ยื่นคำร้องขอเช่าแผงนี้</span>
                <template #suffix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"></polyline></svg>
                </template>
              </AppButton>

              <AppButton
                v-else-if="booking.status === 'lease_applied' || booking.status === 'completed'"
                to="/applications"
                variant="purple"
                size="md"
              >
                <span>ติดตามคำร้อง</span>
                <template #suffix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="9 18 15 12 9 6"></polyline></svg>
                </template>
              </AppButton>
            </template>
          </div>
        </section>

        <section class="progress-card" aria-label="ขั้นตอนการจอง">
          <div v-for="(step, index) in steps" :key="step.number" class="progress-fragment">
            <div class="progress-step" :class="{ active: currentStep === step.number, complete: currentStep > step.number }">
              <div class="progress-circle">
                <svg v-if="currentStep > step.number" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span v-else>{{ step.number }}</span>
              </div>
              <div class="progress-copy">
                <strong>{{ step.label }}</strong>
                <span>{{ step.description }}</span>
              </div>
            </div>
            <div v-if="index < steps.length - 1" class="progress-line" :class="{ complete: currentStep > step.number }"></div>
          </div>
        </section>

        <!-- Waiting Queue Context Banner (แสดงข้อมูลผู้จองหลักและเงื่อนไขคิวสำรองทั้งฝั่งผู้เช่าและเจ้าหน้าที่) -->
        <section v-if="isQueueBooking(booking)" class="queue-context-banner" aria-label="ข้อมูลลำดับคิวสำรอง">
          <div class="qcb-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
          </div>
          <div class="qcb-body">
            <div class="qcb-header">
              <strong>{{ currentRole === 'officer' ? `รายการนี้อยู่ในคิวสำรองลำดับที่ ${booking.queue_position || booking.queue_number || 1} ของแผง ${booking.stall_number}` : `คุณอยู่ในคิวสำรองลำดับที่ ${booking.queue_position || booking.queue_number || 1} ของแผง ${booking.stall_number}` }}</strong>
              <span class="qcb-pill">คิวสำรองที่ {{ booking.queue_position || booking.queue_number || 1 }}</span>
            </div>
            <div class="qcb-details">
              <div v-if="currentRole === 'officer' && activeHolder" class="qcb-detail-row">
                <span class="qcb-dot is-holder"></span>
                <span>ผู้ถือสิทธิ์หลักปัจจุบัน (ตัวจริง): <strong>{{ activeHolder.tenant_name }} ({{ formatBookingNumber(activeHolder.booking_number) }})</strong> — สถานะ: <strong>{{ getStatusLabelByCode(activeHolder.status) }}</strong></span>
              </div>
              <div v-else-if="activeHolderMaskedName" class="qcb-detail-row">
                <span class="qcb-dot is-holder"></span>
                <span>ผู้ถือสิทธิ์หลักปัจจุบัน: <strong>{{ activeHolderMaskedName }}</strong> <span class="qcb-pdpa">(ปิดบังข้อมูลบางส่วนตาม PDPA)</span> — สถานะ: <strong>{{ getStatusLabelByCode(activeHolderStatus) }}</strong></span>
              </div>
              <div v-if="(booking.queue_position || booking.queue_number || 1) > 1" class="qcb-detail-row">
                <span class="qcb-dot is-queue"></span>
                <span>มีคิวสำรองก่อนหน้าอีกลำดับละ <strong>{{ (booking.queue_position || booking.queue_number || 1) - 1 }} ราย</strong></span>
              </div>
            </div>
            <p class="qcb-note">
              <template v-if="currentRole === 'officer'">
                คิวสำรองสามารถกดอนุมัตินัดหมายเข้าชมแผงได้ตามปกติ • ปุ่มยืนยันการเช่าจะเปิดให้กดอัตโนมัติทันทีเมื่อผู้จองหลักสละสิทธิ์
              </template>
              <template v-else>
                คุณสามารถเข้าชมแผงตามวันนัดหมายที่ได้รับอนุมัติได้ตามปกติ • ทันทีที่ผู้จองคิวหลักสละสิทธิ์ ระบบจะเลื่อนคิวของคุณขึ้นเป็นผู้จองหลักและเจ้าหน้าที่จะติดต่อคุณทันที
              </template>
            </p>
          </div>
        </section>

        <section class="booking-information" aria-label="ข้อมูลการจอง">
          <div class="booking-information-grid">
            <div class="booking-information-main">
              <article class="summary-section booked-stalls-section">
                <div class="section-heading stalls-heading">
                  <span class="section-icon" aria-hidden="true"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg></span>
                  <div><h3>แผงที่จอง</h3><p>{{ booking.market_name || '-' }} · {{ siblingBookings.length }} แผง</p></div>
                </div>
                <div class="booked-stalls-scroll" role="region" aria-label="รายการแผงที่จอง" tabindex="0">
                  <table class="booked-stalls-table">
                    <thead><tr><th scope="col">แผงค้า / โครงการ</th><th scope="col">ค่าเช่า/เดือน</th><th scope="col">มัดจำ</th><th scope="col" aria-label="การดำเนินการ"></th></tr></thead>
                    <tbody>
                      <tr v-for="sub in paginatedDetailStalls" :key="sub.id || sub.stall_number" :class="{ 'is-ended': isTerminalStallStatus(sub.status) }">
                        <td>
                          <strong>{{ sub.stall_number }}</strong>
                          <span class="stall-secondary">{{ getSubStallProjectName(sub) || '-' }}</span>
                          <span v-if="sub.status !== booking.status" class="status-badge" :class="getSubStatusToneClass(sub.status)">{{ getStatusLabelByCode(sub.status) }}</span>
                          <span v-if="Number(sub.waitlist_count || 0) > 0" class="stall-secondary">คิวรอ {{ sub.waitlist_count }}</span>
                        </td>
                        <td class="stall-money">฿{{ Number(sub.monthly_rent || sub.stall_monthly_rate || 0).toLocaleString() }}</td>
                        <td class="stall-money">฿{{ Number(sub.deposit ?? sub.deposit_fee ?? sub.monthly_rent ?? sub.stall_monthly_rate ?? 0).toLocaleString() }}</td>
                        <td>
                          <AppButton v-if="!isTerminalStallStatus(sub.status) && !['lease_applied', 'completed'].includes(sub.status)" variant="danger-outline" size="sm" :disabled="actionLoading" @click="handleCancelSubBooking(sub)">ยกเลิกแผงนี้</AppButton>
                          <span v-else class="stall-secondary">—</span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <AppPagination v-if="siblingBookings.length > 5" v-model="detailStallPage" v-model:items-per-page="detailStallsPerPage" :total-items="siblingBookings.length" unit-label="แผง" />
                <dl class="booking-cost-summary">
                  <div><dt>รวมค่าเช่า/เดือน (แผงที่เหลือ)</dt><dd>฿{{ Number(appointmentTotalMonthlyRate || 0).toLocaleString() }}</dd></div>
                  <div><dt>รวมเงินมัดจำ (แผงที่เหลือ)</dt><dd>฿{{ Number(appointmentTotalDeposit || 0).toLocaleString() }}</dd></div>
                </dl>
              </article>
            <article class="summary-section">
              <div class="section-heading">
                <span class="section-icon">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
                </span>
                <h3>ข้อมูลนัดหมาย</h3>
              </div>
              <dl class="detail-list">
                <div class="appointment-highlight"><dt>วันนัดหมาย</dt><dd>{{ formatVisitDate(booking) }}</dd></div>
                <div class="appointment-highlight"><dt>ช่วงเวลา</dt><dd>{{ getVisitSlot(booking) }}</dd></div>
                <div class="full"><dt>จุดนัดพบ</dt><dd>{{ booking.meeting_point || 'เจ้าหน้าที่จะแจ้งหลังอนุมัติ' }}</dd></div>
                <div class="full"><dt>เจ้าหน้าที่ผู้พาชม</dt><dd>{{ booking.escort_officer || 'รอเจ้าหน้าที่แจ้ง' }}</dd></div>
                <div class="full"><dt>สินค้าที่ต้องการจำหน่าย</dt><dd>{{ booking.intended_business || '-' }}</dd></div>
              </dl>
            </article>

            <article class="summary-section contact-section">
              <div class="section-heading">
                <span class="section-icon">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                </span>
                <h3>ข้อมูลผู้จอง</h3>
              </div>
              <dl class="detail-list tenant-list">
                <div><dt>ชื่อ–นามสกุล</dt><dd>{{ booking.tenant_name || '-' }}</dd></div>
                <div><dt>เบอร์โทรศัพท์</dt><dd><a :href="`tel:${booking.tenant_phone}`">{{ booking.tenant_phone || '-' }}</a></dd></div>
                <div v-if="booking.tenant_email"><dt>อีเมล</dt><dd>{{ booking.tenant_email || '-' }}</dd></div>
                <div v-if="booking.tenant_id_card"><dt>เลขบัตรประชาชน</dt><dd>{{ booking.tenant_id_card || '-' }}</dd></div>
              </dl>
            </article>

            </div>
            <article class="summary-section booking-map-section">
              <div class="section-heading booking-map-heading">
                <span class="section-icon" aria-hidden="true"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 5-8 12-8 12S4 15 4 10a8 8 0 1 1 16 0z"></path><circle cx="12" cy="10" r="2.5"></circle></svg></span>
                <div><h3>ตำแหน่งแผงค้า</h3><p>{{ booking.market_name || '-' }} · {{ bookingProjectsSummary }}</p></div>
                <div class="booking-map-legend" aria-label="คำอธิบายสถานะแผง">
                  <span><i class="legend-dot selected"></i>แผงที่คุณเลือก ({{ bookingMapStallCodes.length }})</span>
                  <span><i class="legend-dot available"></i>แผงว่าง</span>
                  <span><i class="legend-dot reserved"></i>ติดจอง</span>
                </div>
              </div>
              <ClientOnly>
                <MarketProjectsOverviewMap
                  v-if="bookingMapMarketId"
                  :key="`${booking.id}-${bookingMapMarketId}`"
                  container-id="booking-detail-maplibre"
                  :market-id="bookingMapMarketId"
                  :selected-stall-codes="bookingMapStallCodes"
                  :selected-stall-keys="bookingMapStallKeys"
                  v-model:map-style="bookingMapStyle"
                  :preview-mode="true"
                />
                <div v-else class="booking-map-unavailable">ยังไม่มีข้อมูลแผนที่ของตลาดนี้</div>
                <template #fallback><div class="booking-map-skeleton" aria-label="กำลังโหลดแผนที่"></div></template>
              </ClientOnly>
            </article>
          </div>

          <div v-if="cancellationReason" class="notice danger"><strong>เหตุผลที่ยกเลิก</strong><span>{{ cancellationReason }}</span></div>
        </section>
      </template>
    </main>

    <!-- 1. Officer Approve Modal (มาตรฐานเดียวกับหน้ารายการจอง) -->
    <div v-if="showApproveModal && booking" class="modal-overlay" @click.self="showApproveModal = false">
      <div class="modal-content" style="max-width: 500px;">
        <div class="modal-header">
          <h3 style="font-size: 1.1rem; font-weight: 700;">อนุมัตินัดหมายเข้าดูแผง ({{ formatBookingNumber(booking.booking_number) }})</h3>
          <button class="btn-modal-close" @click="showApproveModal = false" aria-label="ปิด">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>
        <div class="modal-body">
          <p style="font-size: 0.9rem; color: #475569; margin-bottom: 14px;">
            อนุมัตินัดหมายให้คุณ <strong>{{ booking.tenant_name }}</strong> เข้าดูแผง
            <strong>{{ appointmentStallSummary }}</strong>
            ({{ booking.market_name }}) วันที่ {{ formatVisitDate(booking) }} ({{ getVisitSlot(booking) }})
          </p>
          <div class="form-group" style="margin-bottom: 12px;">
            <label class="form-label">จุดนัดพบหน้างาน *</label>
            <input type="text" v-model="approveMeetingPoint" class="form-input" />
          </div>
          <div class="form-group" style="margin-bottom: 12px;">
            <label class="form-label">เจ้าหน้าที่ผู้พาชม *</label>
            <input type="text" v-model="approveEscortOfficer" class="form-input" />
          </div>
          <div class="form-group">
            <label class="form-label">คำแนะนำเพิ่มเติมสำหรับผู้เช่า</label>
            <textarea v-model="officerNotes" rows="2" class="form-textarea" placeholder="เช่น ให้ติดต่อเจ้าหน้าที่ดูแลตลาด ณ อาคารอำนวยการ ก่อนเข้าชม 15 นาที"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <AppButton variant="secondary" size="md" @click="showApproveModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" size="md" :loading="actionLoading" @click="confirmApprove">ยืนยันการอนุมัตินัดหมาย</AppButton>
        </div>
      </div>
    </div>

    <!-- 2. Officer Reject Modal (มาตรฐานเดียวกับหน้ารายการจอง) -->
    <div v-if="showRejectModal && booking" class="modal-overlay" @click.self="showRejectModal = false">
      <div class="modal-content" style="max-width: 500px;">
        <div class="modal-header">
          <h3 style="font-size: 1.1rem; color: #e11d48; font-weight: 700;">ยกเลิกนัดหมายดูแผง ({{ formatBookingNumber(booking.booking_number) }})</h3>
          <button class="btn-modal-close" @click="showRejectModal = false" aria-label="ปิด">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>
        <div class="modal-body">
          <div style="background: #fff1f2; border: 1px solid #fecdd3; border-radius: 8px; padding: 10px 14px; font-size: 0.85rem; color: #9f1239; margin-bottom: 14px; display: flex; align-items: flex-start; gap: 8px;">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink: 0; margin-top: 2px;"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
            <span>ตามกฎของระบบ เมื่อเจ้าหน้าที่ยกเลิกการจอง สิทธิ์การถือแผงจะถูกปล่อย และเลื่อนคิวที่ 1 ขึ้นมาทันที</span>
          </div>
          <div class="form-group" style="margin-bottom: 12px;">
            <label class="form-label">ระบุเหตุผลในการยกเลิก (บังคับเลือก) *</label>
            <AppSelect
              v-model="rejectReasonDropdown"
              :options="rejectReasonOptions"
              placeholder="เลือกเหตุผลในการยกเลิก"
            />
          </div>
          <div class="form-group">
            <label class="form-label">ข้อความคำชี้แจงเพิ่มเติมถึงผู้เช่า</label>
            <textarea v-model="rejectionReasonText" rows="2" class="form-textarea" placeholder="คำอธิบายเพิ่มเติมเพื่อแจ้งให้ผู้เช่าทราบ..."></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <AppButton variant="secondary" size="md" @click="showRejectModal = false">ปิด</AppButton>
          <AppButton variant="danger" size="md" :loading="actionLoading" @click="confirmReject">ยืนยันยกเลิกและปล่อยคิว</AppButton>
        </div>
      </div>
    </div>

    <!-- 3. Tenant Select Slot Modal (มาตรฐานเดียวกับหน้ารายการจอง) -->
    <div v-if="showSelectSlotModal && booking" class="modal-overlay" @click.self="showSelectSlotModal = false">
      <div class="modal-content" style="max-width: 480px;">
        <div class="modal-header">
          <h3 style="font-size: 1.1rem; font-weight: 700; display: flex; align-items: center; gap: 8px;">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            <span>เลือกวันและเวลานัดชมแผง {{ booking.stall_number }}</span>
          </h3>
          <button class="btn-modal-close" @click="showSelectSlotModal = false" aria-label="ปิด">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group" style="margin-bottom: 14px;">
            <label class="form-label">เลือกวันที่เปิดให้เข้าชมแผง *</label>
            <AppSelect
              v-model="slotDate"
              :options="slotDateOptions"
              placeholder="เลือกวันที่เปิดให้เข้าชมแผง"
            />
          </div>
          <div class="form-group">
            <label class="form-label">เลือกช่วงเวลานัดชม (เฉพาะรอบที่ยังว่าง) *</label>
            <AppSelect
              v-model="slotTime"
              :options="slotTimeOptions"
              placeholder="เลือกช่วงเวลานัดชม"
            />
          </div>
        </div>
        <div class="modal-footer">
          <AppButton variant="secondary" size="md" @click="showSelectSlotModal = false">ยกเลิก</AppButton>
          <AppButton variant="primary" size="md" :loading="actionLoading" @click="confirmSelectSlot">บันทึกวันนัดหมาย</AppButton>
        </div>
      </div>
    </div>

    <LeaseApplicationModal
      v-if="booking"
      :show="showLeaseModal"
      :booking="leaseModalTargetBooking || booking"
      :stall="null"
      @close="showLeaseModal = false; leaseModalTargetBooking = null"
      @success="handleApplicationSubmitted"
    />

    <button
      v-if="booking"
      type="button"
      class="timeline-edge-trigger"
      :class="{ active: showTimeline }"
      :aria-expanded="showTimeline"
      aria-controls="booking-timeline-panel"
      @click="showTimeline = !showTimeline"
    >
      <svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 3-6.7L3 8"></path><path d="M3 3v5h5"></path><path d="M12 7v5l3 2"></path></svg>
      <span>ประวัติการดำเนินการ</span>
    </button>

    <Transition name="timeline-slide">
      <div v-if="showTimeline && booking" class="timeline-backdrop" @click.self="showTimeline = false">
        <aside id="booking-timeline-panel" class="timeline-panel" role="dialog" aria-modal="true" aria-label="ไทม์ไลน์สถานะการจอง">
          <header class="timeline-header">
            <div><span class="eyebrow">รายการกิจกรรมทั้งหมด</span><h2>ประวัติการดำเนินการ</h2></div>
            <button type="button" class="close-button" @click="showTimeline = false" aria-label="ปิดไทม์ไลน์">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
            </button>
          </header>
          <div class="timeline-body">
            <div v-if="timelineItems.length" class="timeline-list">
              <div v-for="(item, index) in timelineItems" :key="index" class="timeline-item">
                <span class="timeline-dot" :class="roleClass(item.actor_role)"></span>
                <div class="timeline-copy">
                  <strong>{{ item.action }}</strong>
                  <time>{{ formatThaiDateTime(item.timestamp) }}</time>
                  <div class="timeline-actor">{{ actorRoleLabel(item.actor_role) }} · {{ item.actor_name || '-' }}</div>
                  <p v-if="item.note">{{ item.note }}</p>
                </div>
              </div>
            </div>
            <div v-else class="timeline-empty">ยังไม่มีประวัติการดำเนินการเพิ่มเติม</div>
          </div>
        </aside>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import {
  getBookingVisitDates,
  setCachedScheduleBookings,
  normalizeVisitDateToIso,
  getAppointmentRoundKey
} from '../../stall/composables/useMarketVisitSchedule'
import MarketProjectsOverviewMap from '../../stall/ui/components/MarketProjectsOverviewMap.vue'
import { BMA_MARKETS_MASTER, getMarketGisConfig } from '../../stall/data/markets_registry'
import { useNotifications } from '~/composables/useNotifications'

const route = useRoute()
const router = useRouter()
const { backTo: returnBackTo, crumbItems: returnCrumbItems } = useReturnTo({ label: 'รายการจอง', to: '/bookings' })
const { get, post } = useApi()
const { user, role: authRole } = useAuth()
const { showConfirm, showSuccess, showError } = useAppAlert()
const { fetchNotifications } = useNotifications()

const loading = ref(true)
const booking = ref<any | null>(null)
const allBookings = ref<any[]>([])
const siblingBookings = ref<any[]>([])
const stallWaitingQueueInfo = ref<any | null>(null)
const leaseModalTargetBooking = ref<any | null>(null)
const detailStallPage = ref(1)
const detailStallsPerPage = ref(5)
const paginatedDetailStalls = computed(() => {
  const start = (detailStallPage.value - 1) * detailStallsPerPage.value
  return siblingBookings.value.slice(start, start + detailStallsPerPage.value)
})
watch([detailStallsPerPage, siblingBookings], () => { detailStallPage.value = 1 })
const bookingMapStyle = ref<'standard' | 'dark' | 'satellite'>('standard')
const bookingMapMarketId = computed(() => {
  const item = booking.value
  if (!item) return null
  const market = Object.values(BMA_MARKETS_MASTER).find(m =>
    String(m.id) === String(item.market_id || '') || m.slug === item.market_slug || m.name === item.market_name
  )
  // Booking records from older data may not include market_id/market_slug.
  // Resolve by name first, then use the registry's standard market fallback
  // so the detail page still renders the same GIS map as Step 3.
  return market ? String(market.id) : String(getMarketGisConfig(item.market_name || '1').id)
})
const bookingMapStallCodes = computed<string[]>(() => [...new Set(siblingBookings.value.flatMap(sub =>
  String(sub.stall_number || sub.stall_ident || '').split(',').map(code => code.trim()).filter(Boolean)
))])
const bookingMapStallKeys = computed<string[]>(() => siblingBookings.value.map(sub =>
  sub.selection_id || sub.stall_selection_id || sub.stall_number || sub.stall_ident
).filter(Boolean))
const showTimeline = ref(false)
const showLeaseModal = ref(false)
const showApproveModal = ref(false)
const showRejectModal = ref(false)
const showSelectSlotModal = ref(false)
const actionLoading = ref(false)

// Approve Modal State
const approveMeetingPoint = ref('หน้ากองอำนวยการตลาด เสาธง 1')
const approveEscortOfficer = ref('นายธนากร สมบูรณ์ (เจ้าหน้าที่เวรนำชมแผง)')
const officerNotes = ref('อนุมัตินัดหมาย ให้เตรียมบัตรประชาชนมาแสดงก่อนเข้าชม')

// Reject Modal State
const rejectReasonDropdown = ref('แผงปิดปรับปรุงระบบสาธารณูปโภค')
const rejectionReasonText = ref('')
const rejectReasonOptions = [
  { value: 'แผงปิดปรับปรุงระบบสาธารณูปโภค', label: 'แผงปิดปรับปรุงระบบสาธารณูปโภค' },
  { value: 'แผงมีผู้เช่าแล้ว / มีสัญญาเดิมผูกพัน', label: 'แผงมีผู้เช่าแล้ว / มีสัญญาเดิมผูกพัน' },
  { value: 'ประเภทสินค้าไม่ตรงกับผังการจัดโซนของตลาด', label: 'ประเภทสินค้าไม่ตรงกับผังการจัดโซนของตลาด' },
  { value: 'ผู้ขอจองมีประวัติค้างชำระหรือผิดเงื่อนไข', label: 'ผู้ขอจองมีประวัติค้างชำระหรือผิดเงื่อนไข' },
  { value: 'อื่นๆ (ระบุในหมายเหตุ)', label: 'อื่นๆ (ระบุในหมายเหตุ)' }
]

// Select Slot Modal State
const slotDate = ref(new Date(Date.now() + 86400000).toISOString().slice(0, 10))
const slotTime = ref('09.30 - 10.30 น.')

const modalAvailableDates = computed(() => {
  const mName = booking.value?.market_name || 'ตลาดนัดจตุจักร'
  const otherBookings = allBookings.value.filter(b => b.id !== booking.value?.id)
  return getBookingVisitDates(mName, 14, otherBookings, mName)
})

const slotDateOptions = computed(() => {
  return modalAvailableDates.value
    .filter(d => !d.isClosed && !d.isAllFull)
    .map(d => ({
      value: d.date,
      label: `${d.dayName} ${d.dayNum} ${d.monthShort} (${d.slotsText})`
    }))
})

const slotTimeOptions = computed(() => {
  const selectedDay = modalAvailableDates.value.find(d => d.date === slotDate.value)
  if (selectedDay && selectedDay.slots?.length) {
    return selectedDay.slots
      .filter(s => !s.isFull && s.enabled !== false)
      .map(s => ({
        value: s.time,
        label: `${s.time} (${s.capacityText})`
      }))
  }
  return [
    { value: '09.30 - 10.30 น.', label: '09.30 - 10.30 น. (ช่วงเช้า)' },
    { value: '11.00 - 12.00 น.', label: '11.00 - 12.00 น. (ช่วงสาย)' },
    { value: '13.00 - 14.00 น.', label: '13.00 - 14.00 น. (ช่วงบ่าย)' },
    { value: '15.30 - 16.30 น.', label: '15.30 - 16.30 น. (ช่วงเย็น)' }
  ]
})

watch(slotDate, () => {
  const opts = slotTimeOptions.value
  if (opts.length > 0 && !opts.some(o => o.value === slotTime.value)) {
    slotTime.value = opts[0].value
  }
})

const isQueueBooking = (b: any): boolean => {
  if (!b) return false
  if (['rejected', 'cancelled', 'not_interested', 'no_show', 'lease_applied', 'completed', 'terminated'].includes(b.status)) {
    return false
  }
  return b.status === 'waiting_queue' || (!b.holding_stall && Number(b.queue_position || b.queue_number) > 0)
}

const currentRole = computed<'tenant' | 'officer'>(() => {
  const r = String(user.value?.role || (authRole as any)?.value || authRole || '').toLowerCase()
  const isOfficer = r === 'admin' || r === 'staff' || r === 'officer'
  if (route.query.role === 'officer' && isOfficer) return 'officer'
  if (route.query.role === 'tenant' && !isOfficer) return 'tenant'
  return isOfficer ? 'officer' : 'tenant'
})

const isTerminalStallStatus = (status?: string) => {
  return ['cancelled', 'rejected', 'not_interested', 'no_show', 'terminated'].includes(String(status || ''))
}

const getSubStatusToneClass = (status?: string) => {
  if (isTerminalStallStatus(status)) return 'tone-muted'
  if (['approved', 'interested'].includes(String(status || ''))) return 'tone-success'
  if (['lease_applied', 'completed'].includes(String(status || ''))) return 'tone-purple'
  return 'tone-warning'
}

const activeSiblingBookings = computed(() => {
  return siblingBookings.value.filter(s => !isTerminalStallStatus(s.status))
})

const STALL_PREFIX_PROJECT_MAP: Record<string, string> = {
  'BKA-01': 'โครงการ 01 (แผงร้านค้า)',
  'BKA-02': 'โครงการ 02 (อาคารพาณิชย์)',
  'BKA-03': 'โครงการ 03 (พื้นที่ด้านซอย 127)',
  'BKA-04': 'โครงการ 04 (พื้นที่ด้านซอย 125)',
  'BKH-01': 'โครงการ 01 (5 เบ็ดเตล็ด)',
  'BKH-02': 'โครงการ 02 (4 เบ็ดเตล็ด)',
  'BKH-03': 'โครงการ 03 (3 อาหารสำเร็จรูป)',
  'BKH-04': 'โครงการ 04 (6 เนื้อสัตว์ อาหารทะเล)',
  'BKH-05': 'โครงการ 05 (1 ผัก-ผลไม้)',
  'BKH-06': 'โครงการ 06 (2 อาหารปรุงและอาหารสำเร็จรูป)',
  'JJ-1': 'โครงการที่ 1',
  'JJ-3': 'โครงการที่ 3',
  'JJ-4': 'โครงการที่ 4',
  'JJ-5': 'โครงการที่ 5',
  'JJ-6': 'โครงการที่ 6',
  'JJ-7': 'โครงการที่ 7',
  'JJ-8': 'โครงการที่ 8',
  'JJ-9': 'โครงการที่ 9',
  'JJ-10': 'โครงการที่ 10',
  'JJ-11': 'โครงการที่ 11',
  'JJ-12': 'โครงการที่ 12',
  'JJ-13': 'โครงการที่ 13',
  'JJ-14': 'โครงการที่ 14',
  'JJ-15': 'โครงการที่ 15',
  'JJ-16': 'โครงการที่ 16',
  'JJ-17': 'โครงการที่ 17',
  'JJ-18': 'โครงการที่ 18',
  'JJ-19': 'โครงการที่ 19',
  'JJ-20': 'โครงการที่ 20/2',
  'JJ-21': 'โครงการที่ 21',
  'JJ-22': 'โครงการที่ 22',
  'JJ-23': 'โครงการที่ 23',
  'JJ-24': 'โครงการที่ 24',
  'JJ-25': 'โครงการที่ 25',
  'JJ-26': 'โครงการที่ 26',
  'JJ-27': 'โครงการที่ 27',
  'JJ-28': 'โครงการที่ 28',
  'JJ-29': 'โครงการที่ 29/2',
  'MIN-01': 'โครงการ 01 (โซนหลัก (ชั้นล่าง))',
  'MIN-02': 'โครงการ 02 (แผงค้าเบ็ดเตล็ด (ชั้นบน))',
  'MIN-03': 'โครงการ 03 (โครงการที่ 3 (ชั้นล่าง))',
  'MIN-04': 'โครงการ 04 (โครงการที่ 2 (ชั้นล่าง))',
  'MIN-05': 'โครงการ 05 (โครงการที่ 4 (ชั้นล่าง))',
  'MIN-06': 'โครงการ 06 (โครงการที่ 1 (ชั้นล่าง))',
  'MIN-07': 'โครงการ 07 (โครงการที่ 7 (ชั้นล่าง))',
  'MIN-08': 'โครงการ 08 (โครงการที่ 6 (ชั้นล่าง))',
  'MIN-09': 'โครงการ 09 (โครงการที่ 8 (ชั้นล่าง))',
  'MIN-10': 'โครงการ 10 (โครงการที่ 5 (ชั้นล่าง))',
  'MIN-11': 'โครงการ 11 (แผงค้าต้นไม้ (ชั้นบน))',
  'MIN-12': 'โครงการ 12 (โซนหลัก (ชั้นบน))',
  'NON-01': 'โครงการ 01 (ด้านใน)',
  'NON-02': 'โครงการ 02 (โซนหลัก)',
  'NON-03': 'โครงการ 03 (ด้านนอก)',
  'NON-04': 'โครงการ 04 (นอกตลาด)',
  'PRA-01': 'โครงการ 01 (D)',
  'PRA-02': 'โครงการ 02 (F)',
  'PRA-03': 'โครงการ 03 (C)',
  'PRA-04': 'โครงการ 04 (A)',
  'PRA-05': 'โครงการ 05 (8)',
  'PRA-06': 'โครงการ 06 (6)',
  'PRA-07': 'โครงการ 07 (4)',
  'PRA-08': 'โครงการ 08 (2)',
  'PRA-09': 'โครงการ 09 (E)',
  'PRA-10': 'โครงการ 10 (9)',
  'PRA-11': 'โครงการ 11 (7)',
  'PRA-12': 'โครงการ 12 (5)',
  'PRA-13': 'โครงการ 13 (3)',
  'PRA-14': 'โครงการ 14 (1)',
  'PRA-15': 'โครงการ 15 (โซนหลัก)',
  'PRA-16': 'โครงการ 16 (G)',
  'RBU-01': 'โครงการ 01 (อาคารตลาดราษฎร์บูรณะ)',
  'RDP-01': 'โครงการ 01 (M)',
  'RDP-02': 'โครงการ 02 (B)',
  'RDP-03': 'โครงการ 03 (D)',
  'RDP-04': 'โครงการ 04 (F)',
  'RDP-05': 'โครงการ 05 (E)',
  'RDP-06': 'โครงการ 06 (C)',
  'RDP-07': 'โครงการ 07 (A)',
  'RDP-08': 'โครงการ 08 (G)',
  'RDP-09': 'โครงการ 09 (H)',
  'RDP-10': 'โครงการ 10 (โซนหลัก)',
  'SIN-01': 'โครงการ 01',
  'THO-01': 'โครงการ 01 (4)',
  'THO-02': 'โครงการ 02 (3)',
  'THO-03': 'โครงการ 03 (1)',
  'THO-04': 'โครงการ 04 (11)',
  'THO-05': 'โครงการ 05 (2)',
  'THO-06': 'โครงการ 06 (6)',
  'THO-07': 'โครงการ 07 (7)',
  'THO-08': 'โครงการ 08 (8)',
  'THO-09': 'โครงการ 09 (โซนหลัก)',
  'THO-10': 'โครงการ 10 (5)',
  'THO-11': 'โครงการ 11 (9)',
  'THO-12': 'โครงการ 12 (1-2)',
  'THO-13': 'โครงการ 13 (ด้านนอก)',
  'TWR-01': 'โซน 1 (อาคารตลาดสดด้านใน)',
  'TWR-02': 'โซน 2 (อาคารตลาดสดรอบนอก)',
  'TWR-03': 'โซน 3 (ศูนย์อาหารเทเวศร์)',
  'TWR-04': 'โซน 4 (ลานต้นโพธิ์)',
  'WON-01': 'โครงการ 01 (ฉ)',
  'WON-02': 'โครงการ 02 (ง)',
  'WON-03': 'โครงการ 03 (จ)',
  'WON-04': 'โครงการ 04 (ก)',
  'WON-05': 'โครงการ 05 (ข)',
  'WON-06': 'โครงการ 06 (ค)',
  'WON-07': 'โครงการ 07 (โซนหลัก)'
}

const getSubStallProjectName = (sub: any): string => {
  const rawCode = String(sub?.stall_number || sub?.stall_ident || '').split(',')[0].trim().toUpperCase()
  const parts = rawCode.split('-')
  if (parts.length >= 2) {
    const prefix = `${parts[0]}-${parts[1]}`
    if (STALL_PREFIX_PROJECT_MAP[prefix]) {
      return STALL_PREFIX_PROJECT_MAP[prefix]
    }
  }
  return sub?.project_name || sub?.stall_project_name || sub?.stall_zone_name || booking.value?.project_name || ''
}

const bookingProjectsSummary = computed(() => {
  const source = activeSiblingBookings.value.length > 0 ? activeSiblingBookings.value : siblingBookings.value
  if (source.length > 0) {
    const uniqueProjs = Array.from(new Set(source.map(s => getSubStallProjectName(s)?.trim()).filter(Boolean)))
    if (uniqueProjs.length > 0) return uniqueProjs.join(', ')
  }
  return booking.value?.project_name || '-'
})

const appointmentStallSummary = computed(() => {
  const source = activeSiblingBookings.value.length > 0 ? activeSiblingBookings.value : siblingBookings.value
  if (source.length > 1) {
    const codes = source.map(s => s.stall_number).filter(Boolean).join(', ')
    return `${codes} (${source.length} แผง)`
  }
  if (source.length === 1) {
    return source[0].stall_number || booking.value?.stall_number || '-'
  }
  return booking.value?.stall_number || '-'
})

const appointmentTotalMonthlyRate = computed(() => {
  const source = activeSiblingBookings.value
  if (siblingBookings.value.length > 0) {
    return source.reduce((sum, s) => sum + (Number(s.monthly_rent || s.stall_monthly_rate) || 0), 0)
  }
  return Number(booking.value?.total_monthly_rent || booking.value?.monthly_rent || booking.value?.stall_monthly_rate || 0)
})

const appointmentTotalDeposit = computed(() => {
  const source = activeSiblingBookings.value
  if (siblingBookings.value.length > 0) {
    return source.reduce((sum, s) => sum + (Number(s.deposit ?? s.deposit_fee ?? s.monthly_rent ?? s.stall_monthly_rate) || 0), 0)
  }
  return Number(booking.value?.total_deposit || booking.value?.deposit || booking.value?.deposit_fee || appointmentTotalMonthlyRate.value)
})

const steps = [
  { number: 1, label: 'จองดูแผง', description: 'ส่งข้อมูลการจอง' },
  { number: 2, label: 'นัดชมแผง', description: 'อนุมัติและเข้าชมแผง' },
  { number: 3, label: 'ยื่นคำร้องขอเช่า', description: 'ส่งคำร้องเพื่อดำเนินการเช่า' }
]

const currentStep = computed(() => {
  const status = booking.value?.status
  if (['lease_applied', 'completed', 'interested'].includes(status)) return 3
  if (['approved', 'waiting_queue', 'pending_decision'].includes(status)) return 2
  return 1
})

const getStatusLabelByCode = (status?: string) => {
  const s = status || 'pending_approval'
  const labels: Record<string, string> = {
    pending: 'รอเจ้าหน้าที่อนุมัติ',
    pending_approval: 'รอเจ้าหน้าที่อนุมัติ',
    pending_slot: 'รอเลือกวันนัดชม',
    waiting_queue: 'รอคิว',
    approved: 'นัดหมายสำเร็จ',
    interested: 'รอผู้เช่ายื่นคำร้อง',
    pending_decision: 'รอผู้เช่ายื่นคำร้อง',
    lease_applied: 'ยื่นคำร้องสำเร็จ',
    completed: 'ยื่นคำร้องสำเร็จ',
    cancelled: 'ยกเลิก',
    rejected: 'ยกเลิก',
    not_interested: 'สละสิทธิ์',
    no_show: 'ไม่มาตามนัด',
    terminated: 'สิ้นสุดสัญญาเช่า'
  }
  return labels[s] || s
}

const activeHolder = computed(() => {
  const stallNo = booking.value?.stall_number
  if (!stallNo) return null
  return allBookings.value.find(
    b => b.stall_number === stallNo && b.holding_stall && !isQueueBooking(b) && !['rejected', 'cancelled', 'not_interested', 'no_show', 'terminated'].includes(b.status)
  ) || null
})

const activeHolderMaskedName = computed(() => {
  if (activeHolder.value?.tenant_name) {
    const clean = String(activeHolder.value.tenant_name).trim().replace(/^(นาย|นางสาว|นาง|น\.ส\.|คุณ)\s*/, '')
    const parts = clean.split(/\s+/).filter(Boolean)
    if (parts.length > 1) return `คุณ${parts[0]} ${parts[1].charAt(0)}.***`
    if (parts.length === 1) return `คุณ${parts[0].slice(0, 3)}***`
  }
  const masked = stallWaitingQueueInfo.value?.masked_current_booker
  if (masked && masked !== '-') return masked
  return 'ผู้จองคิวหลัก'
})

const activeHolderStatus = computed(() => {
  return activeHolder.value?.status || stallWaitingQueueInfo.value?.active_holder_status || 'pending_approval'
})

const statusLabel = computed(() => {
  const status = booking.value?.status
  return getStatusLabelByCode(status)
})

const statusToneClass = computed(() => {
  const status = booking.value?.status
  if (['cancelled', 'rejected', 'not_interested', 'no_show', 'terminated'].includes(status)) return 'tone-muted'
  if (['approved', 'interested'].includes(status)) return 'tone-success'
  if (['lease_applied', 'completed'].includes(status)) return 'tone-purple'
  if (status === 'waiting_queue') return 'tone-blue'
  return 'tone-warning'
})

const timelineItems = computed(() => Array.isArray(booking.value?.timeline) ? booking.value.timeline : [])
const cancellationReason = computed(() => booking.value?.rejection_reason || booking.value?.cancel_reason || booking.value?.not_interested_reason || '')
const normalizeBookingOwnerName = (value: unknown) => String(value || '')
  .replace(/^(คุณ|นาย|นางสาว|นาง|น\.ส\.)\s*/g, '')
  .replace(/\([^)]*\)/g, '')
  .replace(/\s+/g, ' ')
  .trim()
  .toLowerCase()
const isBookingOwnedByCurrentTenant = (item: any) => {
  if (currentRole.value !== 'tenant') return true
  const currentUser = user.value as any
  if (currentUser?.id !== undefined && item?.tenant_user_id !== undefined && item?.tenant_user_id !== null) {
    return String(item.tenant_user_id) === String(currentUser.id)
  }
  const bookingName = normalizeBookingOwnerName(item?.tenant_name || item?.applicant_name)
  const currentName = normalizeBookingOwnerName(currentUser?.full_name)
  if (bookingName && currentName && bookingName === currentName) return true
  const demoOwnerNames: Record<string, string> = { tenant: 'สมชาย ใจดี', tenant2: 'นภาพร รัตนศิลป์', tenant3: 'วิชัย เจริญทรัพย์' }
  return bookingName === normalizeBookingOwnerName(demoOwnerNames[String(currentUser?.username || '')])
}
const actionableStatuses = ['pending', 'pending_approval', 'pending_slot', 'waiting_queue', 'approved', 'interested', 'pending_decision', 'lease_applied', 'completed']
const hasActionButtons = computed(() => actionableStatuses.includes(booking.value?.status))

const copyTenantPhone = async () => {
  const phone = String(booking.value?.tenant_phone || '').trim()
  if (!phone) return
  try {
    await navigator.clipboard.writeText(phone)
    await showSuccess('คัดลอกเบอร์โทรศัพท์ผู้จองแล้ว', 'คัดลอกสำเร็จ')
  } catch {
    await showError('ไม่สามารถคัดลอกเบอร์โทรศัพท์ได้', 'ข้อผิดพลาด')
  }
}

const nextAction = computed(() => {
  const item = booking.value
  const status = item?.status
  const isOfficer = currentRole.value === 'officer'

  if (['cancelled', 'rejected', 'not_interested', 'no_show', 'terminated'].includes(status)) {
    return { tone: 'complete', owner: 'รายการสิ้นสุดแล้ว', title: 'ไม่มีขั้นตอนที่ต้องดำเนินการต่อ', description: cancellationReason.value || 'รายการนี้ถูกปิดและไม่มีผลต่อสิทธิ์การจองแล้ว' }
  }
  if (['lease_applied', 'completed'].includes(status)) {
    return { tone: 'complete', owner: 'ดำเนินการแล้ว', title: 'ยื่นคำร้องขอเช่าเรียบร้อย', description: isOfficer ? 'รอตรวจสอบคำร้องในขั้นตอนจัดการคำร้องเช่า' : 'ติดตามผลการพิจารณาได้ที่หน้าคำร้องของฉัน' }
  }
  if (isQueueBooking(item)) {
    const position = item?.queue_position || item?.queue_number || 1
    if (status === 'approved') {
      return isOfficer
        ? {
            tone: 'waiting',
            owner: `คิวสำรองที่ ${position}`,
            title: `คิวสำรองที่ ${position} • รอผลจากผู้จองหลัก (ปุ่มยืนยันการเช่าจะเปิดให้กดเมื่อผู้จองหลักสละสิทธิ์)`,
            description: `${formatVisitDate(item)} · ${getVisitSlot(item)}`
          }
        : {
            tone: 'waiting',
            owner: `คิวสำรองที่ ${position}`,
            title: 'เข้าชมแผงตามนัด • หากเข้าชมแล้ว กรุณารอผลจากผู้จองคิวหลัก',
            description: 'หากคิวหลักสละสิทธิ์ เจ้าหน้าที่จะติดต่อคุณและเปิดสิทธิ์ให้ยื่นคำร้องทันที'
          }
    }
    return isOfficer
      ? {
          tone: 'action',
          owner: 'เจ้าหน้าที่ต้องทำ',
          title: `คิวสำรองที่ ${position} • รอเจ้าหน้าที่กดอนุมัติวันเวลานัดชมแผง`,
          description: `${formatVisitDate(item)} · ${getVisitSlot(item)}`
        }
      : {
          tone: 'waiting',
          owner: `คิวสำรองที่ ${position}`,
          title: `คิวสำรองที่ ${position} • รอเจ้าหน้าที่อนุมัติวันและเวลานัดชมแผง`,
          description: `${formatVisitDate(item)} · ${getVisitSlot(item)}`
        }
  }
  if (['pending', 'pending_approval'].includes(status)) {
    return isOfficer
      ? { tone: 'action', owner: 'เจ้าหน้าที่ต้องทำ', title: 'ตรวจสอบวันเวลา แล้วอนุมัตินัดชมแผง', description: `ผู้เช่าขอนัดชมวันที่ ${formatVisitDate(item)} เวลา ${getVisitSlot(item)}` }
      : { tone: 'waiting', owner: 'เจ้าหน้าที่ดำเนินการ', title: 'คุณยังไม่ต้องดำเนินการ', description: 'เจ้าหน้าที่กำลังตรวจสอบวันและเวลานัดหมาย ระบบจะแจ้งเมื่ออนุมัติแล้ว' }
  }
  if (status === 'pending_slot') {
    return isOfficer
      ? { tone: 'waiting', owner: 'ผู้จองต้องทำ', title: 'รอผู้เช่าเลือกวันและเวลานัดชมแผง', description: 'ถึงคิวของผู้จองแล้ว เมื่อผู้เช่าเลือกเวลาแล้ว รายการจะกลับมาให้เจ้าหน้าที่อนุมัติ' }
      : { tone: 'action', owner: 'คุณต้องทำ', title: 'ถึงคิวของคุณแล้ว กรุณาเลือกวันและเวลานัดชมแผง', description: 'กดปุ่ม "เลือกวันนัดชม" เพื่อระบุวันและช่วงเวลาที่สะดวกเข้าชมแผงค้า' }
  }
  if (status === 'approved') {
    return isOfficer
      ? { tone: 'action', owner: 'เจ้าหน้าที่ต้องทำ', title: 'ยืนยันการเช่าหลังพาชมแผง', description: 'หลังทราบผลจากผู้เช่าแล้ว ให้ยืนยันว่าตกลงเช่าหรือไม่เช่า' }
      : { tone: 'action', owner: 'คุณต้องทำ', title: 'ไปตามวัน เวลา และจุดนัดหมาย', description: `${formatVisitDate(item)} · ${getVisitSlot(item)} · ${item?.meeting_point || 'จุดนัดพบตามที่เจ้าหน้าที่แจ้ง'}` }
  }
  if (['interested', 'pending_decision'].includes(status)) {
    const days = Number(item?.application_deadline_days || 3)
    let deadlineNote = `ภายใน ${days} วัน`
    if (item?.expires_at) {
      const exp = new Date(item.expires_at)
      if (!isNaN(exp.getTime())) {
        const diffHours = Math.max(1, Math.ceil((exp.getTime() - Date.now()) / 3600000))
        const diffDays = Math.ceil(diffHours / 24)
        const thaiMonthsShort = ['ม.ค.', 'ก.พ.', 'มี.ค.', 'เม.ย.', 'พ.ค.', 'มิ.ย.', 'ก.ค.', 'ส.ค.', 'ก.ย.', 'ต.ค.', 'พ.ย.', 'ธ.ค.']
        const dStr = `${exp.getDate()} ${thaiMonthsShort[exp.getMonth()]} ${exp.getFullYear() + 543}`
        deadlineNote = `ภายใน ${days} วัน • ครบกำหนด ${dStr} (${diffHours < 24 ? `เหลืออีก ${diffHours} ชม.` : `เหลืออีก ${diffDays} วัน`})`
      }
    }
    return isOfficer
      ? { tone: 'waiting', owner: 'ผู้จองต้องทำ', title: 'รอผู้เช่ายื่นคำร้องขอเช่าแผง', description: `ผู้เช่าต้องส่งคำร้อง ${deadlineNote}` }
      : { tone: 'action', owner: 'คุณต้องทำ', title: 'ยื่นคำร้องขอเช่าแผง', description: `กรอกและส่งคำร้อง ${deadlineNote}` }
  }
  return { tone: 'waiting', owner: 'ระบบกำลังตรวจสอบ', title: 'ติดตามสถานะจากหน้านี้', description: 'ยังไม่มีขั้นตอนที่คุณต้องดำเนินการในขณะนี้' }
})

const loadBooking = async () => {
  loading.value = true
  try {
    const response: any = await get('/bookings', { _t: Date.now() })
    const items = Array.isArray(response) ? response : (Array.isArray(response?.data) ? response.data : [])
    allBookings.value = items
    setCachedScheduleBookings(items)
    const key = String(route.params.id || '')
    const found = items.find(item => {
      const matchesParent = String(item.id) === key || String(item.booking_number).replace(/^#/, '') === key.replace(/^#/, '')
      const matchesSub = Array.isArray(item.sub_bookings) && item.sub_bookings.some((s: any) => String(s.id) === key)
      return (matchesParent || matchesSub) && isBookingOwnedByCurrentTenant(item)
    }) || null
    if (found) {
      siblingBookings.value = Array.isArray(found.sub_bookings) && found.sub_bookings.length > 0 ? found.sub_bookings : [found]
      booking.value = found
      if (isQueueBooking(found) && found.stall_number) {
        try {
          const qRes: any = await get(`/bookings/waiting-queue/${encodeURIComponent(found.stall_number)}`)
          stallWaitingQueueInfo.value = qRes?.data || qRes || null
        } catch {
          stallWaitingQueueInfo.value = null
        }
      } else {
        stallWaitingQueueInfo.value = null
      }
    } else {
      siblingBookings.value = []
      booking.value = null
      stallWaitingQueueInfo.value = null
    }
    await fetchNotifications()
  } catch (error) {
    console.error('Failed to load booking details:', error)
    siblingBookings.value = []
    booking.value = null
  } finally {
    loading.value = false
  }
}

const openLeaseForSubBooking = (sub: any) => {
  leaseModalTargetBooking.value = sub
  showLeaseModal.value = true
}

const handleCancelSubBooking = async (sub: any) => {
  const confirmed = await showConfirm(
    `คุณต้องการยกเลิก/สละสิทธิ์เฉพาะแผง ${sub.stall_number} ใช่หรือไม่?\n(ระบบจะปล่อยสิทธิ์เฉพาะแผง ${sub.stall_number} ให้ผู้รอคิวลำดับแรกหรือปล่อยว่าง ส่วนแผงอื่นในใบจองนี้ยังคงอยู่ตามปกติ)`,
    { title: `ยืนยันยกเลิกเฉพาะแผง ${sub.stall_number}`, type: 'warning', confirmText: 'ยืนยันยกเลิกแผงนี้', cancelText: 'ย้อนกลับ' }
  )
  if (!confirmed) return
  actionLoading.value = true
  try {
    await post(`/bookings/${sub.id}/cancel`, {
      reason: currentRole.value === 'officer' ? 'เจ้าหน้าที่ยกเลิกเฉพาะแผงนี้' : 'ผู้เช่าขอยกเลิกเฉพาะแผงนี้',
      stall_ident: sub.stall_number,
      cancel_all: false,
      by_officer: currentRole.value === 'officer'
    })
    await loadBooking()
    await showSuccess(`ยกเลิกเฉพาะแผง ${sub.stall_number} เรียบร้อยแล้ว ยอดเงินและมัดจำคำนวณจากแผงที่เหลือโดยอัตโนมัติ`, 'สำเร็จ')
  } catch (error: any) {
    await showError(error?.data?.detail || 'ไม่สามารถดำเนินการได้', 'ข้อผิดพลาด')
  } finally {
    actionLoading.value = false
  }
}

// 1. Officer Approve Modal Handlers
const openApproveModal = () => {
  if (!booking.value) return
  approveMeetingPoint.value = booking.value.meeting_point || 'หน้ากองอำนวยการตลาด เสาธง 1'
  approveEscortOfficer.value = booking.value.escort_officer || 'นายธนากร สมบูรณ์ (เจ้าหน้าที่เวรนำชมแผง)'
  officerNotes.value = 'อนุมัตินัดหมาย ให้เตรียมบัตรประชาชนมาแสดงก่อนเข้าชม'
  showApproveModal.value = true
}

const confirmApprove = async () => {
  if (!booking.value) return
  actionLoading.value = true
  try {
    await post(`/bookings/${booking.value.id}/approve`, {
      meeting_point: approveMeetingPoint.value,
      escort_officer: approveEscortOfficer.value,
      officer_notes: officerNotes.value
    })
    showApproveModal.value = false
    await loadBooking()
    await showSuccess('อนุมัตินัดหมายการเข้าชมเรียบร้อยแล้ว', 'อนุมัติสำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาดในการอนุมัติ', 'ข้อผิดพลาด')
  } finally {
    actionLoading.value = false
  }
}

// 2. Officer Reject Modal Handlers
const openRejectModal = () => {
  if (!booking.value) return
  rejectionReasonText.value = ''
  showRejectModal.value = true
}

const confirmReject = async () => {
  if (!booking.value) return
  actionLoading.value = true
  try {
    await post(`/bookings/${booking.value.id}/reject`, {
      rejection_reason: rejectReasonDropdown.value,
      officer_notes: rejectionReasonText.value
    })
    showRejectModal.value = false
    await loadBooking()
    await showSuccess('ปฏิเสธการนัดหมายและเลื่อนลำดับคิวเรียบร้อย', 'สำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาด', 'ข้อผิดพลาด')
  } finally {
    actionLoading.value = false
  }
}

// 3. Officer Mark Visited (ยืนยันการเช่า)
const handleMarkVisited = async () => {
  if (!booking.value) return
  const b = booking.value
  const confirmed = await showConfirm(
    `ยืนยันว่าคุณ ${b.tenant_name} เข้าชมแผง ${appointmentStallSummary.value} เรียบร้อยแล้วใช่หรือไม่?\n(ระบบจะเปลี่ยนสถานะเป็น "รอผู้เช่ายื่นคำร้อง" เพื่อให้ผู้เช่ายื่นคำร้องขอเช่าแผงที่คงเหลือในใบจองนี้ได้ทันที)`,
    {
      title: 'ยืนยันการเช่า',
      type: 'info',
      confirmText: 'ยืนยันการเช่า',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return
  actionLoading.value = true
  try {
    await post(`/bookings/${b.id}/mark-visited`, {
      officer_notes: 'เจ้าหน้าที่ยืนยันการเช่าเรียบร้อย รอผู้เช่ายื่นคำร้อง'
    })
    await loadBooking()
    await showSuccess('ยืนยันการเช่าเรียบร้อยแล้ว สถานะเปลี่ยนเป็น "รอผู้เช่ายื่นคำร้อง"', 'สำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาด', 'ข้อผิดพลาด')
  } finally {
    actionLoading.value = false
  }
}

// 4. Officer Cancel (approved, interested, pending_slot, waiting_queue)
const handleOfficerCancel = async () => {
  if (!booking.value) return
  const b = booking.value
  const isQueue = isQueueBooking(b)
  const isInterested = b.status === 'interested' || b.status === 'pending_decision'
  const activeCount = activeSiblingBookings.value.length
  const title = isQueue ? 'ยืนยันยกเลิกการต่อคิว' : isInterested ? 'ยืนยันยกเลิกใบจองทั้งใบ' : 'ยืนยันยกเลิกนัดหมายทั้งใบ'
  const message = isQueue
    ? `คุณต้องการ "ยกเลิกการต่อคิว" ของคุณ ${b.tenant_name} (แผง ${appointmentStallSummary.value} คิวสำรองที่ ${b.queue_position || b.queue_number || 1}) ใช่หรือไม่?\n(ระบบจะนำชื่อออกจากคิวและปรับลำดับคิวข้างหลังขึ้นมาแทน)`
    : isInterested
      ? `คุณต้องการยกเลิกใบจอง / บันทึกการสละสิทธิ์ของคุณ ${b.tenant_name} (แผง ${appointmentStallSummary.value}) ใช่หรือไม่?\n(หากต้องการยกเลิกเฉพาะบางแผง สามารถกดปุ่ม "ยกเลิกแผงนี้" แยกรายแผงด้านล่างได้)`
      : `คุณต้องการ "ยกเลิกนัดหมายทั้งใบจอง" ของคุณ ${b.tenant_name} (${activeCount > 1 ? `รวม ${activeCount} แผง: ` : 'แผง '}${appointmentStallSummary.value}) ใช่หรือไม่?\n(หากต้องการยกเลิกเฉพาะบางแผง สามารถกดปุ่ม "ยกเลิกแผงนี้" แยกรายแผงด้านล่างได้)`

  const confirmed = await showConfirm(message, {
    title,
    type: 'warning',
    confirmText: 'ยืนยันยกเลิก',
    cancelText: 'ย้อนกลับ'
  })
  if (!confirmed) return

  actionLoading.value = true
  try {
    await post(`/bookings/${b.id}/cancel`, {
      reason: isQueue
        ? 'เจ้าหน้าที่ยกเลิกการต่อคิวตามคำขอของผู้เช่า'
        : isInterested
          ? 'เจ้าหน้าที่บันทึกยกเลิกสิทธิ์ / ผู้เช่าแจ้งสละสิทธิ์ไม่ยื่นคำร้อง'
          : 'เจ้าหน้าที่ยกเลิกการนัดหมาย / ผู้เช่าไม่มาตามนัด',
      cancel_all: true,
      by_officer: true
    })
    await loadBooking()
    await showSuccess(
      isQueue
        ? 'ยกเลิกการต่อคิวและปรับลำดับคิวเรียบร้อยแล้ว'
        : isInterested
          ? 'ยกเลิกสิทธิ์และปล่อยแผงค้าคืนสู่ระบบเรียบร้อยแล้ว'
          : 'ยกเลิกนัดหมายและปรับลำดับคิวเรียบร้อยแล้ว',
      'สำเร็จ'
    )
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาดในการยกเลิก', 'ข้อผิดพลาด')
  } finally {
    actionLoading.value = false
  }
}

// 5. Tenant Cancel / สละสิทธิ์
const handleTenantCancel = async (reason: string) => {
  if (!booking.value) return
  const b = booking.value
  const isQueue = isQueueBooking(b)
  const activeCount = activeSiblingBookings.value.length
  const isMultiStall = activeCount > 1
  const confirmed = await showConfirm(
    isQueue
      ? `คุณต้องการสละสิทธิ์การต่อคิวแผง ${b.stall_number} (คิวสำรองที่ ${b.queue_position || b.queue_number || 1}) หรือไม่?\n(หากสละสิทธิ์ คิวจะถูกปล่อยและไม่สามารถกู้คืนได้)`
      : isMultiStall
        ? `คุณต้องการยกเลิกใบจองรหัส ${formatBookingNumber(b.booking_number)} ทั้งหมด ${activeCount} แผงที่เหลือ (${appointmentStallSummary.value}) หรือไม่?\n(หากต้องการยกเลิกเฉพาะบางแผง สามารถกดปุ่ม "ยกเลิกแผงนี้" แยกรายแผงด้านล่างได้)`
        : `คุณต้องการยกเลิกรายการนี้หรือไม่? ${b.holding_stall ? '(หากยกเลิก แผงจะถูกส่งต่อให้คิวถัดไปทันที)' : ''}`,
    {
      title: isQueue ? 'ยืนยันสละสิทธิ์คิวสำรอง' : isMultiStall ? `ยืนยันยกเลิกใบจองทั้งใบ (${activeCount} แผง)` : 'ยืนยันยกเลิกการจอง',
      type: 'warning',
      confirmText: isQueue ? 'ยืนยันสละสิทธิ์' : 'ยืนยันยกเลิก',
      cancelText: 'ย้อนกลับ'
    }
  )
  if (!confirmed) return

  actionLoading.value = true
  try {
    await post(`/bookings/${b.id}/cancel`, {
      reason,
      cancel_all: true,
      by_officer: false
    })
    await loadBooking()
    await showSuccess('ยกเลิกรายการเรียบร้อยแล้ว', 'สำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาด', 'ข้อผิดพลาด')
  } finally {
    actionLoading.value = false
  }
}

// 6. Tenant Select Slot Modal Handlers
const openSelectSlotModal = () => {
  if (!booking.value) return
  const b = booking.value
  const existingIso = normalizeVisitDateToIso(b?.visit_date)
  const openDates = modalAvailableDates.value.filter(d => !d.isClosed && !d.isAllFull)
  if (existingIso && openDates.some(d => d.date === existingIso)) {
    slotDate.value = existingIso
  } else if (openDates.length > 0) {
    slotDate.value = openDates[0].date
  }
  const availableSlotsForDate = slotTimeOptions.value
  if (b?.visit_time_slot && availableSlotsForDate.some(o => o.value === b.visit_time_slot)) {
    slotTime.value = b.visit_time_slot
  } else if (availableSlotsForDate.length > 0) {
    slotTime.value = availableSlotsForDate[0].value
  }
  showSelectSlotModal.value = true
}

const confirmSelectSlot = async () => {
  if (!booking.value) return
  actionLoading.value = true
  try {
    await post(`/bookings/${booking.value.id}/select-slot`, {
      visit_date: slotDate.value,
      visit_time_slot: slotTime.value
    })
    showSelectSlotModal.value = false
    await loadBooking()
    await showSuccess('บันทึกวันนัดหมายใหม่เรียบร้อยแล้ว', 'สำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาด', 'ข้อผิดพลาด')
  } finally {
    actionLoading.value = false
  }
}

const handleApplicationSubmitted = async () => {
  showLeaseModal.value = false
  await loadBooking()
}

const formatBookingNumber = (value: string) => value?.startsWith('#') ? value : `#${value || '-'}`

const formatThaiDateTime = (value: string) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('th-TH', { dateStyle: 'long', timeStyle: 'short', timeZone: 'Asia/Bangkok' }).format(date).replace(' น.', '') + ' น.'
}

const formatVisitDate = (item: any) => {
  if (!item?.visit_date) return 'รอยืนยันวันนัดหมาย'
  if (/จันทร์|อังคาร|พุธ|พฤหัสบดี|ศุกร์|เสาร์|อาทิตย์/.test(item.visit_date)) return item.visit_date
  const date = new Date(item.visit_date)
  if (Number.isNaN(date.getTime())) return item.visit_date
  return new Intl.DateTimeFormat('th-TH', { dateStyle: 'full', timeZone: 'Asia/Bangkok' }).format(date)
}

const getVisitSlot = (item: any) => item?.visit_time_slot || 'รอยืนยันเวลา'
const actorRoleLabel = (role: string) => role === 'officer' ? 'เจ้าหน้าที่' : role === 'system' ? 'ระบบ' : 'ผู้เช่า'
const roleClass = (role: string) => role === 'officer' ? 'officer' : role === 'system' ? 'system' : 'tenant'

onMounted(loadBooking)
watch(() => route.params.id, loadBooking)
</script>

<style scoped>
.booking-detail-page { min-height: 100vh; background: #f8fafc; color: #0f172a; }
.detail-container {
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
}
.detail-hero, .progress-card, .summary-shell, .detail-empty, .detail-skeleton { background: #fff; border: 1px solid #e2e8f0; border-radius: 16px; box-shadow: 0 1px 3px rgba(15,23,42,.03); }
.detail-hero { display: flex; align-items: center; justify-content: space-between; gap: 20px; padding: 22px 24px; margin-bottom: 16px; }
.detail-title-wrap, .hero-actions { display: flex; align-items: center; gap: 14px; }
.detail-title-icon { width: 46px; height: 46px; display: grid; place-items: center; border-radius: 12px; background: #eff6ff; border: 1px solid #dbeafe; color: #2563eb; }
.eyebrow { display: block; color: #64748b; font-size: .78rem; font-weight: 600; margin-bottom: 3px; }
.detail-hero h1 { margin: 0; font-size: 1.65rem; line-height: 1.1; }
.created-at { margin-top: 5px; color: #94a3b8; font-size: .8rem; }
.status-badge { padding: 6px 13px; border-radius: 999px; font-size: .82rem; font-weight: 700; white-space: nowrap; }
.tone-warning { background: #fef3c7; color: #b45309; }.tone-success { background: #dcfce7; color: #15803d; }.tone-blue { background: #e0f2fe; color: #0369a1; }.tone-purple { background: #ede9fe; color: #6d28d9; }.tone-muted { background: #f1f5f9; color: #64748b; }
.booking-at-a-glance { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 0; margin-bottom: 16px; overflow: hidden; background: #fff; border: 1px solid #e2e8f0; border-radius: 14px; box-shadow: 0 1px 3px rgba(15,23,42,.03); }
.glance-item { display: flex; flex-direction: column; justify-content: center; min-width: 0; min-height: 68px; padding: 13px 18px; border-right: 1px solid #eef2f6; }
.glance-item:last-child { border-right: 0; }.glance-item > span { margin-bottom: 4px; color: #64748b; font-size: .73rem; font-weight: 600; }.glance-item strong, .glance-item > a { overflow: hidden; color: #1e293b; font-size: .88rem; font-weight: 700; line-height: 1.35; text-decoration: none; text-overflow: ellipsis; white-space: nowrap; }.glance-contact { position: relative; padding-right: 74px; }.glance-contact > a { color: #2563eb; }.glance-contact-actions { position: absolute; right: 14px; bottom: 14px; display: inline-flex; gap: 5px; }.glance-icon-button { display: grid; width: 28px; height: 28px; place-items: center; padding: 0; border: 1px solid #bfdbfe; border-radius: 7px; background: #eff6ff; color: #2563eb; cursor: pointer; }.glance-icon-button:hover { background: #dbeafe; color: #1d4ed8; }
.progress-card { display: flex; align-items: center; padding: 22px 30px; margin-bottom: 16px; }
.progress-fragment { display: contents; }
.progress-step { display: flex; align-items: center; gap: 10px; min-width: 170px; }
.progress-circle { width: 34px; height: 34px; flex: 0 0 34px; display: grid; place-items: center; border-radius: 50%; border: 2px solid #cbd5e1; background: #f8fafc; color: #64748b; font-weight: 700; }
.progress-step.active .progress-circle { border-color: #2563eb; background: #eff6ff; color: #2563eb; box-shadow: 0 0 0 4px #dbeafe; }
.progress-step.complete .progress-circle { border-color: #2563eb; background: #2563eb; color: #fff; }
.progress-copy { display: flex; flex-direction: column; }.progress-copy strong { font-size: .88rem; }.progress-copy span { color: #94a3b8; font-size: .73rem; margin-top: 2px; }
.progress-step.active .progress-copy strong, .progress-step.complete .progress-copy strong { color: #1d4ed8; }
.progress-line { height: 2px; flex: 1; min-width: 50px; margin: 0 18px; background: #e2e8f0; }.progress-line.complete { background: #2563eb; }
.next-action-card { display: flex; align-items: center; gap: 14px; min-height: 76px; padding: 15px 18px; margin-bottom: 16px; border: 1px solid #bfdbfe; border-left: 4px solid #2563eb; border-radius: 14px; background: #fff; box-shadow: 0 1px 3px rgba(15,23,42,.03); }
.next-action-icon { width: 40px; height: 40px; flex: 0 0 40px; display: grid; place-items: center; border-radius: 11px; background: #eff6ff; color: #2563eb; }
.next-action-copy { flex: 1; min-width: 0; }.next-action-meta { display: flex; align-items: center; gap: 10px; margin-bottom: 6px; color: #64748b; font-size: .75rem; font-weight: 650; }.next-action-owner { padding: 2px 7px; border-radius: 999px; background: #dbeafe; color: #1d4ed8; font-size: var(--font-size-2xs); line-height: 1.35; font-weight: 600; white-space: nowrap; }.next-action-copy h2 { margin: 0; color: #0f172a; font-size: 1rem; line-height: 1.35; }.next-action-copy p { margin: 6px 0 0; color: #64748b; font-size: .8rem; line-height: 1.45; }.next-action-buttons { flex: 0 0 auto; display: flex; align-items: center; justify-content: flex-end; gap: 9px; flex-wrap: wrap; }
.next-action-card { flex-wrap: wrap; }
.next-action-note { flex: 0 0 100%; display: flex; align-items: center; gap: 12px; margin-top: 4px; padding: 10px 14px; border-radius: 8px; background: #eff6ff; color: #475569; font-size: .8rem; line-height: 1.35; }
.next-action-note strong { color: #1d4ed8; font-weight: 700; white-space: nowrap; }
.next-action-waiting { border-color: #bae6fd; border-left-color: #0ea5e9; }.next-action-waiting .next-action-icon { background: #f0f9ff; color: #0284c7; }.next-action-waiting .next-action-owner { background: #e0f2fe; color: #0369a1; }
.next-action-complete { border-color: #e2e8f0; border-left-color: #94a3b8; }.next-action-complete .next-action-icon { background: #f1f5f9; color: #64748b; }.next-action-complete .next-action-owner { background: #f1f5f9; color: #64748b; }
.queue-context-banner { display: flex; align-items: flex-start; gap: 14px; padding: 16px 20px; margin-bottom: 16px; border: 1px solid #fde68a; border-left: 4px solid #d97706; border-radius: 14px; background: #fffbeb; color: #78350f; }
.qcb-icon { width: 38px; height: 38px; flex: 0 0 38px; display: grid; place-items: center; border-radius: 10px; background: #fef3c7; color: #d97706; }
.qcb-body { flex: 1; min-width: 0; }
.qcb-header { display: flex; align-items: center; justify-content: space-between; gap: 10px; flex-wrap: wrap; margin-bottom: 6px; }
.qcb-header strong { font-size: 0.9375rem; font-weight: 700; color: #92400e; }
.qcb-pill { padding: 3px 10px; border-radius: 999px; background: #fef3c7; border: 1px solid #fcd34d; color: #b45309; font-size: 0.75rem; font-weight: 700; }
.qcb-details { display: flex; flex-direction: column; gap: 4px; margin-bottom: 6px; font-size: 0.84rem; color: #78350f; }
.qcb-detail-row { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.qcb-dot { width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0; }
.qcb-dot.is-holder { background: #10b981; }
.qcb-dot.is-queue { background: #d97706; }
.qcb-pdpa { font-size: 0.78rem; color: #92400e; }
.qcb-note { margin: 0; padding-top: 6px; border-top: 1px dashed #fcd34d; font-size: 0.8125rem; color: #92400e; line-height: 1.45; }
.summary-shell { padding: 24px; }
.summary-heading { display: flex; justify-content: space-between; align-items: flex-end; gap: 18px; padding: 0 2px 18px; margin-bottom: 20px; border-bottom: 1px solid #eef2f6; }.summary-heading h2 { margin: 0; font-size: 1.18rem; }.summary-heading p { margin: 6px 0 0; color: #64748b; font-size: .8rem; }.summary-market { color: #475569; font-size: .84rem; font-weight: 600; }
.summary-grid { display: grid; grid-template-columns: 1fr 1fr; align-items: stretch; gap: 16px; }
.summary-section { min-width: 0; height: auto; padding: 20px; border: 1px solid #e2e8f0; border-radius: 14px; background: #fff; box-shadow: 0 1px 2px rgba(15,23,42,.025); }.summary-section.full-width { grid-column: 1 / -1; }
.section-heading { display: flex; align-items: center; gap: 10px; padding-bottom: 13px; margin-bottom: 15px; border-bottom: 1px solid #f1f5f9; }.section-heading h3 { margin: 0; color: #1e293b; font-size: .96rem; }.section-icon { width: 34px; height: 34px; display: grid; place-items: center; border-radius: 9px; background: #eff6ff; color: #2563eb; }
.detail-list { display: grid; grid-template-columns: 1fr 1fr; gap: 16px 20px; margin: 0; }.detail-list > div { min-width: 0; }.detail-list .full { grid-column: 1 / -1; }.detail-list dt { color: #64748b; font-size: .76rem; margin-bottom: 4px; }.detail-list dd { margin: 0; color: #1e293b; font-size: .9rem; font-weight: 650; line-height: 1.45; overflow-wrap: anywhere; }.detail-list a { color: #2563eb !important; text-decoration: none; }.primary-value { color: #2563eb !important; }.money-value { color: #059669 !important; }.appointment-highlight { padding: 11px 12px; border: 1px solid #dbeafe; border-radius: 10px; background: #f8fbff; }.appointment-highlight dd { color: #1d4ed8; }
.detail-multi-stalls { margin-top: 16px; padding-top: 14px; border-top: 1px dashed #cbd5e1; display: flex; flex-direction: column; gap: 8px; }
.dms-title { font-size: 0.78rem; font-weight: 700; color: #475569; margin-bottom: 2px; }
.dms-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 10px 12px; background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; flex-wrap: wrap; }
.dms-left { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; font-size: 0.85rem; }
.dms-left strong { color: #1e293b; font-weight: 700; }
.dms-left span { color: #059669; font-weight: 600; font-size: 0.8125rem; }
.dms-actions { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.dms-btn { display: inline-flex; align-items: center; justify-content: center; padding: 5px 11px; border-radius: 7px; font-size: 0.78rem; font-weight: 600; font-family: inherit; cursor: pointer; transition: all 0.15s ease; }
.dms-btn-primary { background: #eff6ff; color: #2563eb; border: 1px solid #bfdbfe; }
.dms-btn-primary:hover { background: #2563eb; color: #ffffff; }
.dms-btn-danger { background: #ffffff; color: #dc2626; border: 1px solid #fca5a5; }
.dms-btn-danger:hover { background: #fef2f2; }
.notice { display: flex; gap: 12px; margin-top: 16px; padding: 13px 15px; border-radius: 10px; font-size: .84rem; }.notice strong { white-space: nowrap; }.notice.success { background: #ecfdf5; border: 1px solid #a7f3d0; color: #065f46; }.notice.danger { background: #fff1f2; border: 1px solid #fecdd3; color: #9f1239; }
.timeline-edge-trigger { position: fixed; top: 27%; right: 0; z-index: 900; width: 44px; height: 188px; transform: translateY(-50%); display: block; padding: 0; overflow: hidden; border: 1px solid #dbeafe; border-right: 0; border-radius: 12px 0 0 12px; background: rgba(255,255,255,.98); color: #2563eb; box-shadow: -6px 7px 20px rgba(15,23,42,.11); backdrop-filter: blur(10px); font-family: inherit; font-size: .72rem; font-weight: 700; white-space: nowrap; cursor: pointer; transition: background .18s ease, width .18s ease, box-shadow .18s ease, transform .18s ease; }
.timeline-edge-trigger > svg { position: absolute; top: 15px; left: 50%; width: 17px; height: 17px; transform: translateX(-50%); }.timeline-edge-trigger > span { position: absolute; top: 56%; left: 50%; transform: translate(-50%,-50%) rotate(-90deg); transform-origin: center; }.timeline-edge-trigger:hover, .timeline-edge-trigger.active { width: 48px; background: #eff6ff; color: #1d4ed8; box-shadow: -8px 9px 24px rgba(37,99,235,.16); }
.timeline-backdrop { position: fixed; inset: 0; z-index: 1000; display: flex; justify-content: flex-end; background: rgba(15,23,42,.45); backdrop-filter: blur(2px); }.timeline-panel { width: min(480px, 92vw); height: 100%; background: #fff; box-shadow: -16px 0 40px rgba(15,23,42,.16); overflow-y: auto; }.timeline-header { position: sticky; top: 0; z-index: 2; display: flex; align-items: center; justify-content: space-between; padding: 22px 24px; border-bottom: 1px solid #e2e8f0; background: #fff; }.timeline-header h2 { margin: 0; font-size: 1.2rem; }.close-button { width: 36px; height: 36px; display: grid; place-items: center; border: 1px solid #e2e8f0; border-radius: 9px; background: #fff; color: #475569; cursor: pointer; }.timeline-body { padding: 24px; }.timeline-list { display: flex; flex-direction: column; }.timeline-item { position: relative; display: flex; gap: 14px; padding-bottom: 24px; }.timeline-item:not(:last-child)::before { content: ''; position: absolute; left: 5px; top: 13px; bottom: 0; width: 1px; background: #e2e8f0; }.timeline-dot { position: relative; z-index: 1; width: 11px; height: 11px; margin-top: 5px; flex: 0 0 11px; border-radius: 50%; background: #3b82f6; box-shadow: 0 0 0 4px #dbeafe; }.timeline-dot.officer { background: #10b981; box-shadow: 0 0 0 4px #d1fae5; }.timeline-dot.system { background: #8b5cf6; box-shadow: 0 0 0 4px #ede9fe; }.timeline-copy { flex: 1; min-width: 0; }.timeline-copy strong { display: block; font-size: .9rem; }.timeline-copy time, .timeline-actor { display: block; color: #64748b; font-size: .75rem; margin-top: 3px; }.timeline-copy p { margin: 8px 0 0; padding: 8px 10px; border-radius: 8px; background: #f8fafc; color: #475569; font-size: .8rem; line-height: 1.45; }.timeline-empty { padding: 30px 10px; text-align: center; color: #94a3b8; font-size: .85rem; }
.timeline-slide-enter-active, .timeline-slide-leave-active { transition: background-color .28s ease, backdrop-filter .28s ease; }.timeline-slide-enter-active .timeline-panel, .timeline-slide-leave-active .timeline-panel { transition: transform .32s cubic-bezier(.22,1,.36,1), box-shadow .32s ease; }.timeline-slide-enter-from, .timeline-slide-leave-to { background: rgba(15,23,42,0); backdrop-filter: blur(0); }.timeline-slide-enter-from .timeline-panel, .timeline-slide-leave-to .timeline-panel { transform: translateX(100%); box-shadow: none; }
.detail-empty, .detail-skeleton { padding: 40px; text-align: center; }.detail-empty svg { color: #94a3b8; }.detail-empty h1 { margin: 12px 0 4px; font-size: 1.2rem; }.detail-empty p { margin: 0; color: #64748b; }.skeleton-line, .skeleton-panel, .skeleton-grid > div { position: relative; overflow: hidden; background: #eef2f7; }.skeleton-line::after, .skeleton-panel::after, .skeleton-grid > div::after { content: ''; position: absolute; inset: 0; transform: translateX(-100%); background: linear-gradient(90deg, transparent, rgba(255,255,255,.75), transparent); animation: shimmer 1.4s infinite; }.skeleton-line { height: 20px; border-radius: 7px; margin-bottom: 10px; }.skeleton-line.wide { width: 38%; }.skeleton-line.short { width: 20%; }.skeleton-panel { height: 120px; border-radius: 12px; margin: 24px 0; }.skeleton-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }.skeleton-grid > div { height: 180px; border-radius: 12px; } @keyframes shimmer { to { transform: translateX(100%); } }
.modal-overlay { position: fixed; inset: 0; z-index: 1100; background: rgba(15, 23, 42, 0.5); backdrop-filter: blur(2px); display: flex; align-items: center; justify-content: center; padding: 16px; }
.modal-content { width: 100%; background: #ffffff; border-radius: 16px; border: 1px solid #e2e8f0; box-shadow: 0 20px 50px rgba(15, 23, 42, 0.2); overflow: hidden; }
.modal-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 18px 22px; border-bottom: 1px solid #eef2f6; background: #ffffff; }
.modal-header h3 { margin: 0; color: #0f172a; }
.btn-modal-close { width: 34px; height: 34px; display: inline-flex; align-items: center; justify-content: center; border-radius: 8px; border: 1px solid #e2e8f0; background: #ffffff; color: #64748b; cursor: pointer; transition: all 0.15s ease; }
.btn-modal-close:hover { background: #f8fafc; color: #0f172a; }
.modal-body { padding: 20px 22px; }
.modal-footer { display: flex; align-items: center; justify-content: flex-end; gap: 10px; padding: 14px 22px; border-top: 1px solid #eef2f6; background: #f8fafc; }
.form-group { display: flex; flex-direction: column; gap: 6px; }
.form-label { font-size: 0.8125rem; font-weight: 600; color: #334155; }
.form-input, .form-textarea { width: 100%; box-sizing: border-box; padding: 10px 12px; border: 1px solid #cbd5e1; border-radius: 8px; font-family: inherit; font-size: 0.875rem; color: #0f172a; background: #ffffff; outline: none; transition: border-color 0.15s ease, box-shadow 0.15s ease; }
.form-input:focus, .form-textarea:focus { border-color: #2563eb; box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12); }
@media (max-width: 900px) { .detail-hero { align-items: flex-start; }.hero-actions { flex-direction: column; align-items: flex-end; }.booking-at-a-glance { grid-template-columns: repeat(2, 1fr); }.glance-item:nth-child(2) { border-right: 0; }.glance-item:nth-child(-n+2) { border-bottom: 1px solid #eef2f6; }.progress-card { align-items: flex-start; overflow-x: auto; }.progress-step { min-width: 145px; }.progress-copy span { display: none; }.summary-grid { grid-template-columns: 1fr; }.summary-section.full-width { grid-column: auto; } }
@media (max-width: 640px) { .detail-hero { flex-direction: column; }.hero-actions { width: 100%; flex-direction: row; justify-content: space-between; align-items: center; }.detail-title-icon { width: 40px; height: 40px; }.detail-hero h1 { font-size: 1.35rem; }.booking-at-a-glance { grid-template-columns: 1fr; }.glance-item, .glance-item:nth-child(2) { min-height: 60px; border-right: 0; border-bottom: 1px solid #eef2f6; }.glance-item:last-child { border-bottom: 0; }.progress-card { padding: 18px 16px; }.progress-step { min-width: 104px; flex-direction: column; text-align: center; }.progress-copy strong { font-size: .78rem; }.progress-line { min-width: 26px; margin: 16px 8px 0; }.next-action-card { align-items: flex-start; flex-wrap: wrap; }.next-action-copy { min-width: calc(100% - 56px); }.next-action-buttons { width: 100%; padding-left: 54px; }.summary-shell { padding: 16px; }.summary-heading { align-items: flex-start; flex-direction: column; }.detail-list, .tenant-list { grid-template-columns: 1fr; }.detail-list .full { grid-column: auto; }.notice { flex-direction: column; gap: 4px; }.timeline-edge-trigger { top: 32%; bottom: auto; width: 40px; height: 180px; transform: translateY(-50%); border-radius: 11px 0 0 11px; font-size: .68rem; }.timeline-edge-trigger:hover { width: 44px; } }

/* Compact detail layout; the original hero and next-action card themes are preserved. */
.booking-information-grid { display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1fr); gap: 16px; align-items: stretch; }
.booking-information-main { display: grid; gap: 16px; min-width: 0; }
.detail-container { display: flex; flex-direction: column; }
.detail-container > .next-action-card { order: 1; }
.detail-container > .detail-hero { order: 2; }
.detail-container > .progress-card { order: 3; }
.detail-container > .queue-context-banner { order: 4; }
.detail-container > .booking-information { order: 5; }
.detail-container > .notice { order: 6; }
.booking-information .summary-section { padding: 24px; border-radius: 16px; }
.booking-information .booked-stalls-section, .booking-information .booking-map-section { padding: 20px; }
.booking-information .section-heading { margin-bottom: 12px; padding-bottom: 12px; }
.booking-information .section-heading h3 { font-size: var(--font-size-lg); font-weight: 700; color: #0f172a; }
.booking-information .section-heading p { margin: 4px 0 0; font-size: var(--font-size-sm); color: #64748b; }
.booking-information .detail-list { gap: 12px 20px; }
.booking-information .detail-list dt { font-size: var(--font-size-xs); line-height: 1.4; color: #64748b; }
.booking-information .detail-list dd { font-size: var(--font-size-sm); font-weight: 700; line-height: 1.45; color: #0f172a; }
.booking-information .tenant-list { grid-template-columns: repeat(2, minmax(0, 1fr)); }
.booking-information-main > .summary-section:not(.booked-stalls-section) .detail-list { grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px 24px; }
.booking-information-main > .summary-section:not(.booked-stalls-section) .detail-list .full { grid-column: auto; }
.booking-information-main > .summary-section:not(.booked-stalls-section) .detail-list > div { padding: 10px 0; border-radius: 0; background: transparent; }
.booking-information-main > .summary-section:not(.booked-stalls-section) .detail-list .appointment-highlight { background: #f8fbff; border: 1px solid #bfdbfe; border-radius: 12px; padding: 11px 12px; }
.booking-information-main > .summary-section:not(.booked-stalls-section) .detail-list .appointment-highlight dt { color: #64748b; }
.booking-information-main > .summary-section:not(.booked-stalls-section) .detail-list .appointment-highlight dd { color: #2563eb; }
.booking-information .section-icon { width: 40px; height: 40px; border-radius: 10px; }
.booking-information .section-icon svg { width: 20px; height: 20px; }
.booked-stalls-scroll { max-width: 100%; overflow-x: auto; -webkit-overflow-scrolling: touch; }
.booked-stalls-table { width: 100%; table-layout: fixed; border-collapse: collapse; text-align: left; font-size: var(--font-size-sm); }
.booked-stalls-table th:nth-child(1), .booked-stalls-table td:nth-child(1) { width: 38%; }
.booked-stalls-table th:nth-child(2), .booked-stalls-table td:nth-child(2) { width: 20%; }
.booked-stalls-table th:nth-child(3), .booked-stalls-table td:nth-child(3) { width: 20%; }
.booked-stalls-table th:nth-child(4), .booked-stalls-table td:nth-child(4) { width: 22%; }
.booked-stalls-table th { padding: 14px 12px; font-size: var(--font-size-xs); font-weight: 500; color: #64748b; border-bottom: 1px solid #f1f5f9; white-space: nowrap; }
.booked-stalls-table td { padding: 16px 12px; border-bottom: 1px solid #f1f5f9; color: #334155; }
.booked-stalls-table tbody tr:hover { background: #f8fafc; }
.booked-stalls-table td:last-child { text-align: right; }
.booked-stalls-table .is-ended { background: #f8fafc; }
.booked-stalls-table .stall-money { font-weight: 700; color: #1e293b; white-space: nowrap; }
.stall-secondary { display: block; margin-top: 4px; font-size: var(--font-size-xs); color: #64748b; }
.booked-stalls-table .status-badge { margin-top: 6px; font-size: var(--font-size-xs); padding: 4px 8px; }
.booking-cost-summary { margin: 12px 0 0; padding: 12px; border-radius: 8px; background: #f8fafc; display: grid; gap: 8px; font-size: var(--font-size-sm); }
.booking-cost-summary > div { display: flex; justify-content: space-between; gap: 12px; }
.booking-cost-summary dt { color: #64748b; }
.booking-cost-summary dd { margin: 0; font-weight: 700; color: #1e293b; white-space: nowrap; }
.booking-map-section { position: relative; overflow: hidden; height: 100%; min-height: 560px; display: flex; flex-direction: column; }
.booking-map-heading { align-items: flex-start; }
.booking-map-legend { display: inline-flex; align-items: center; justify-content: flex-end; flex-wrap: wrap; gap: 12px; margin-left: auto; color: #475569; font-size: var(--font-size-xs); white-space: nowrap; }
.booking-map-legend span { display: inline-flex; align-items: center; gap: 6px; }
.legend-dot { display: inline-block; width: 12px; height: 12px; border-radius: 4px; }
.legend-dot.selected { background: #2563eb; box-shadow: 0 0 0 3px #dbeafe; }
.legend-dot.available { background: #10b981; }
.legend-dot.reserved { background: #f59e0b; }
.booking-map-section :deep(.preview-head-left), .booking-map-section :deep(.preview-legend-row) { display: none; }
.booking-map-section :deep(.preview-booked-pill-btn) { display: none; }
.booking-map-section :deep(.preview-top-bar) { position: absolute; top: 76px; left: 30px; right: 30px; z-index: 30; padding: 0; border: 0 !important; box-shadow: none !important; pointer-events: none; }
.booking-map-section :deep(.preview-head-right) { width: 100%; align-items: flex-start; }
.booking-map-section :deep(.preview-booked-pill-btn) { max-width: min(100%, 420px); white-space: normal; font-size: var(--font-size-xs); pointer-events: auto; }
.booking-map-section :deep(.preview-badge-stalls) { white-space: normal; overflow-wrap: anywhere; }
.booking-map-section :deep(.projects-master-container), .booking-map-section :deep(.maplibre-overview-wrapper) { height: 100% !important; min-height: 0; flex: 1; }
.booking-map-skeleton, .booking-map-unavailable { height: 490px; border-radius: 12px; background: #f1f5f9; }
.booking-map-unavailable { display: grid; place-items: center; color: #64748b; font-size: var(--font-size-sm); }
.booking-map-skeleton { background: linear-gradient(90deg, #f1f5f9, #e2e8f0, #f1f5f9); background-size: 200% 100%; animation: detail-map-shimmer 1.5s infinite; }
@keyframes detail-map-shimmer { to { background-position: -200% 0; } }
.progress-card { padding: 16px 24px; }
.progress-copy strong { font-size: var(--font-size-sm); }
.progress-copy span { font-size: var(--font-size-xs); }
.detail-hero h1 { font-size: var(--font-size-2xl); overflow-wrap: anywhere; }
.detail-title-wrap { min-width: 0; }
.next-action-meta, .created-at, .eyebrow { font-size: var(--font-size-xs); }
.next-action-copy p { font-size: var(--font-size-sm); }
.timeline-edge-trigger { font-size: var(--font-size-2xs); }
@media (max-width: 1100px) { .booking-information-grid { grid-template-columns: minmax(0, 1fr); }.booking-map-section { min-height: 380px; }.booking-map-legend { justify-content: flex-start; margin-left: 0; }.booking-map-section :deep(.preview-top-bar) { top: 72px; left: 22px; right: 22px; }.booking-map-section :deep(.projects-master-container), .booking-map-section :deep(.maplibre-overview-wrapper) { height: 380px !important; } }
@media (max-width: 640px) { .booking-information .summary-section { padding: 16px; }.booking-information .tenant-list { grid-template-columns: minmax(0, 1fr); }.booked-stalls-table { min-width: 520px; }.progress-card { padding: 16px 12px; }.progress-step { min-width: 0; flex: 1; }.progress-copy strong { font-size: var(--font-size-xs); }.progress-line { min-width: 12px; margin-inline: 8px; }.detail-hero h1 { font-size: var(--font-size-xl); }.timeline-edge-trigger { font-size: var(--font-size-2xs); } }
@media (prefers-reduced-motion: reduce) { .booking-map-skeleton { animation: none; } }

</style>
