<template>
  <div class="app-detail-page">
    <AppHeader />

    <main class="content-body detail-container">
      <AppBreadcrumb
        :items="crumbItems"
        :current="application ? `รายละเอียด ${applicationNumber}` : 'รายละเอียดคำร้อง'"
        :back-to="backTo"
      />

      <div v-if="loading" class="detail-skeleton" aria-label="กำลังโหลดรายละเอียดคำร้อง">
        <div class="skeleton-line wide"></div>
        <div class="skeleton-line short"></div>
        <div class="skeleton-panel"></div>
        <div class="skeleton-grid"><div></div><div></div><div></div></div>
      </div>

      <div v-else-if="!application" class="detail-empty">
        <svg width="38" height="38" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
        <h1>ไม่พบคำร้อง</h1>
        <p>คำร้องนี้อาจถูกลบหรือหมายเลขคำร้องไม่ถูกต้อง</p>
      </div>

      <template v-else>
        <section class="detail-hero">
          <div class="detail-title-wrap">
            <div class="detail-title-icon">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line></svg>
            </div>
            <div>
              <div class="eyebrow">{{ isStall ? 'รายละเอียดคำร้องขอเช่าแผงค้า' : 'รายละเอียดคำร้องทั่วไป' }}</div>
              <h1>#{{ applicationNumber }}</h1>
              <div v-if="createdAtLabel" class="created-at">วันที่ยื่นคำร้อง {{ createdAtLabel }}</div>
            </div>
          </div>

          <div class="hero-actions">
            <AppButton v-if="isStall" variant="outline" size="sm" @click="showApplicationDocument = true"><template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg></template>ตัวอย่างคำร้อง</AppButton>
            <span class="status-badge" :class="statusClass">{{ statusLabel }}</span>
          </div>
        </section>

        <section class="next-action-card" :class="`next-action-${nextAction.tone}`" aria-labelledby="next-action-title">
          <div class="next-action-icon">
            <svg v-if="nextAction.tone === 'waiting'" width="21" height="21" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"></circle><polyline points="12 7 12 12 15 14"></polyline></svg>
            <svg v-else-if="nextAction.tone === 'complete'" width="21" height="21" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
            <svg v-else-if="nextAction.tone === 'danger'" width="21" height="21" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
            <svg v-else width="21" height="21" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"></path><path d="m13 6 6 6-6 6"></path></svg>
          </div>
          <div class="next-action-copy">
            <div class="next-action-meta">
              <span>ขั้นตอนถัดไป</span>
              <span class="next-action-owner" :class="{ 'is-mine': isMyTurn }">{{ isMyTurn ? 'คุณต้องทำ' : nextAction.owner }}</span>
            </div>
            <h2 id="next-action-title">{{ nextAction.title }}</h2>
            <p>{{ nextAction.description }}</p>
          </div>
          <div class="next-action-buttons">
            <!-- 3 Officer Review Actions (เมื่อมีสิทธิ์พิจารณาคำร้องนี้) -->
            <template v-if="canReview">
              <AppButton
                variant="primary"
                size="md"
                :loading="actionSubmitting"
                :disabled="actionSubmitting"
                @click="processApproval()"
              >
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"></polyline></svg>
                </template>
                <span>{{ approveButtonLabel }}</span>
              </AppButton>

              <AppButton
                variant="secondary"
                size="md"
                :disabled="actionSubmitting"
                @click="openRejectModal('return_revision')"
              >
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"></polyline><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path></svg>
                </template>
                <span>ส่งกลับแก้ไข</span>
              </AppButton>

              <AppButton
                variant="danger-outline"
                size="md"
                :disabled="actionSubmitting"
                @click="openRejectModal('reject')"
              >
                <template #prefix>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
                </template>
                <span>ไม่อนุมัติ</span>
              </AppButton>
            </template>

            <!-- Back navigation button -->
            <AppButton :to="backTo" variant="secondary" size="md">
              <template #prefix>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>
              </template>
              <span>{{ returnUrl ? 'กลับหน้าที่เปิดมา' : 'กลับสู่รายการคำร้อง' }}</span>
            </AppButton>
          </div>
        </section>

        <section class="progress-card" aria-label="ขั้นตอนการพิจารณาคำร้อง">
          <div class="booking-stepper-container">
            <div class="stepper-track">
              <template v-for="(step, index) in steps" :key="step.number">
                <div class="stepper-step" :class="{ completed: isStepCompleted(step.number), active: isStepActive(step.number) }">
                  <div class="step-circle">
                    <svg v-if="isStepCompleted(step.number)" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"></polyline></svg>
                    <span v-else>{{ step.number }}</span>
                  </div>
                  <div class="step-label-text">{{ step.label }}</div>
                </div>
                <div v-if="index < steps.length - 1" class="stepper-connector" :class="{ filled: isConnectorFilled(step.number) }"></div>
              </template>
            </div>
          </div>
        </section>

        <section class="application-information" aria-label="ข้อมูลคำร้อง">
          <div class="application-information-grid">
            <div class="application-information-main">
              <article class="summary-section">
                <div class="section-heading">
                  <span class="section-icon" aria-hidden="true">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
                  </span>
                  <h3>{{ isStall ? 'ข้อมูลแผงที่ขอเช่า' : 'ข้อมูลคำร้อง' }}</h3>
                  <div v-if="isStall" class="stall-heading"><span v-for="stall in requestedStalls" :key="stall.code">แผง {{ stall.code }}</span></div>
                </div>
                <template v-if="isStall">
                  <div v-for="stall in requestedStalls" :key="stall.code" class="requested-stall">
                    <dl class="stall-grid">
                      <div><dt>โซน</dt><dd>{{ stall.zone || '-' }}</dd></div>
                      <div><dt>หมายเลขแผง</dt><dd>{{ stall.code }} <span v-if="stall.floor" class="stall-floor">{{ stall.floor }}</span></dd></div>
                      <div><dt>ประเภทพื้นที่</dt><dd>{{ stall.category || '-' }}</dd></div>
                      <div><dt>ขนาดพื้นที่ (ตร.ม.)</dt><dd>{{ stall.size || '-' }}</dd></div>
                      <div><dt>ค่าธรรมเนียมสิทธิแรกเข้า (บาท)</dt><dd>{{ formatAmount(stall.entryFee) }}</dd></div>
                      <div><dt>ค่าเช่าต่อเดือน (บาท)</dt><dd>{{ formatAmount(stall.monthlyRent) }}</dd></div>
                      <div><dt>ค่าส่วนกลาง (บาท)</dt><dd>{{ formatAmount(stall.commonFee) }}</dd></div>
                      <div><dt>ค่าประกันล่วงหน้า (บาท)</dt><dd>{{ formatAmount(stall.depositFee) }}</dd></div>
                    </dl>
                    <div class="stall-photos"><div>รูปภาพ</div><div class="stall-photo-list"><img v-for="photo in stall.photos" :key="photo" :src="photo" :alt="'รูปแผง ' + stall.code"></div></div>
                  </div>
                  <dl v-if="application.details" class="detail-list stall-extra"><div class="full"><dt>รายละเอียด</dt><dd>{{ application.details }}</dd></div></dl>
                </template>
                <dl v-else class="detail-list">
                  <div v-if="application.market_name"><dt>ตลาด</dt><dd>{{ application.market_name }}</dd></div>
                  <div v-if="isStall"><dt>แผงค้า</dt><dd>{{ application.stall_number || application.stall_ident || '-' }}</dd></div>
                  <div v-if="isStall"><dt>ชื่อร้านค้า</dt><dd>{{ application.shop_name || '-' }}</dd></div>
                  <div v-if="isStall"><dt>ประเภทสินค้า</dt><dd>{{ application.product_type || application.product_category || 'ทั่วไป' }}</dd></div>
                  <div v-if="!isStall"><dt>หัวข้อคำร้อง</dt><dd>{{ application.subject || application.shop_name || 'คำร้องทั่วไป' }}</dd></div>
                  <div v-if="application.details" class="full"><dt>รายละเอียด</dt><dd>{{ application.details }}</dd></div>
                </dl>
              </article>

              <article class="summary-section">
                <div class="section-heading">
                  <span class="section-icon" aria-hidden="true">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="9" y1="15" x2="15" y2="15"></line><line x1="9" y1="11" x2="12" y2="11"></line></svg>
                  </span>
                  <h3>ข้อมูลการยื่นคำร้อง</h3>
                </div>
                <dl class="detail-list application-detail-list">
                  <div><dt>ชื่อ–นามสกุล</dt><dd>{{ application.applicant_name || '-' }}</dd></div>
                  <div><dt>เบอร์โทรศัพท์</dt><dd>{{ application.applicant_phone || '-' }}</dd></div>
                  <div v-if="application.applicant_id_card"><dt>เลขประจำตัวประชาชน</dt><dd>{{ application.applicant_id_card }}</dd></div>
                  <div v-if="application.applicant_email"><dt>อีเมล</dt><dd>{{ application.applicant_email }}</dd></div>
                  <div v-if="application.applicant_address && application.applicant_address !== '-'" class="full"><dt>ที่อยู่</dt><dd>{{ application.applicant_address }}</dd></div>
                  <div v-if="isStall"><dt>ชื่อร้านค้า</dt><dd>{{ application.shop_name || '-' }}</dd></div>
                  <div v-if="isStall"><dt>ประเภทสินค้า</dt><dd>{{ application.product_type || application.product_category || '-' }}</dd></div>
                  <div><dt>วันที่ยื่นคำร้อง</dt><dd>{{ createdAtLabel || '-' }}</dd></div>
                  <div><dt>เลขที่คำร้อง</dt><dd>{{ applicationNumber || '-' }}</dd></div>
                  <template v-if="isStall">
                    <div><dt>ระยะเวลาเช่า</dt><dd>{{ application.requested_period_months || application.lease_term_months ? `${application.requested_period_months || application.lease_term_months} เดือน` : '-' }}</dd></div>
                    <div><dt>วันที่ขอเริ่มสัญญา</dt><dd>{{ application.requested_start_date || '-' }}</dd></div>
                    <div v-if="application.stall_monthly_rate || application.monthly_rate"><dt>ค่าเช่า/เดือน</dt><dd>฿{{ Number(application.stall_monthly_rate || application.monthly_rate).toLocaleString() }}</dd></div>
                    <div v-if="application.proposed_rent_type"><dt>ประเภทการเช่า</dt><dd>{{ ({ trial: 'ทดลองขาย', temporary: 'ทำการค้าแบบชั่วคราว', contract: 'ทำการค้าแบบมีสัญญา', monthly: 'รายเดือน', daily: 'รายวัน', yearly: 'รายปี', annual: 'รายปี' } as Record<string, string>)[String(application.proposed_rent_type).toLowerCase()] || application.proposed_rent_type }}</dd></div>
                  </template>
                  <div><dt>สำเนาบัตรประชาชน</dt><dd>{{ application.has_id_card_copy ? 'แนบแล้ว' : 'ไม่ได้แนบ' }}</dd></div>


                </dl>

                <div v-if="canRevise" class="submitted-notice is-revision">
                  <span>เจ้าหน้าที่ส่งคำร้องกลับให้แก้ไข ตอนนี้คุณสามารถแก้ไขข้อมูลที่ยื่นได้</span>
                  <AppButton :to="`/applications?revise=${application.id}`" variant="primary" size="sm">แก้ไขคำร้อง</AppButton>
                </div>

              </article>

            </div>

            <article class="summary-section comments-section">
              <div class="section-heading timeline-section-heading">
                <div class="timeline-heading-text">
                  <span class="timeline-eyebrow">รายการกิจกรรมทั้งหมด</span>
                  <h3>ติดตามการดำเนินการยื่นคำร้อง</h3>
                </div>
              </div>

              <div v-if="timeline.length" class="timeline-list">
                <div v-for="(ev, idx) in timeline" :key="idx" class="timeline-item">
                  <span class="timeline-dot" :class="ev.tone" aria-hidden="true"></span>
                  <div class="timeline-copy">
                    <strong>{{ ev.title }}</strong>
                    <time v-if="ev.time">{{ ev.time }}</time>
                    <div class="timeline-actor">{{ ev.actor }}</div>
                    <p v-if="ev.note">{{ ev.note }}</p>
                  </div>
                </div>
              </div>
              <p v-else class="empty-text">ยังไม่มีประวัติการดำเนินการสำหรับคำร้องนี้</p>
            </article>
          </div>
        </section>
      </template>

      <!-- Reject / Return Revision Modal -->
      <Teleport to="body">
        <div v-if="showRejectModal" class="detail-modal-overlay" @click.self="closeRejectModal">
          <div class="detail-modal-content">
            <div class="detail-modal-header" :class="{ 'is-danger': rejectAction === 'reject' }">
              <div class="modal-icon-badge">
                <svg v-if="rejectAction === 'return_revision'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="1 4 1 10 7 10"></polyline><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path></svg>
                <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
              </div>
              <div class="modal-header-titles">
                <h3 class="modal-title">{{ rejectAction === 'return_revision' ? 'ส่งกลับให้ผู้เช่าแก้ไขคำร้อง' : 'ไม่อนุมัติคำร้อง' }}</h3>
                <p class="modal-sub">คำร้องหมายเลข {{ applicationNumber }}</p>
              </div>
              <button type="button" class="btn-modal-close" @click="closeRejectModal" aria-label="ปิด">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              </button>
            </div>

            <div class="detail-modal-body">
              <div class="quick-reasons-block">
                <label class="modal-input-label">ข้อความสำเร็จรูป:</label>
                <div class="quick-chips">
                  <button
                    v-for="(r, idx) in quickRejectReasons"
                    :key="idx"
                    type="button"
                    class="quick-chip"
                    @click="setQuickReason(r)"
                  >
                    {{ r }}
                  </button>
                </div>
              </div>

              <div class="modal-form-group">
                <label class="modal-input-label required">{{ rejectAction === 'return_revision' ? 'สิ่งที่ต้องแก้ไข / คำแนะนำเพิ่มเติม' : 'เหตุผลการไม่อนุมัติ' }}</label>
                <textarea
                  v-model="rejectComment"
                  class="reason-textarea"
                  rows="4"
                  :placeholder="rejectAction === 'return_revision' ? 'ระบุสิ่งที่ต้องแก้ไข เพื่อแจ้งให้ผู้เช่าทราบ...' : 'ระบุเหตุผลการไม่อนุมัติ เพื่อบันทึกในระบบ...'"
                ></textarea>
                <p v-if="rejectErrorMsg" class="modal-error-msg">{{ rejectErrorMsg }}</p>
              </div>
            </div>

            <div class="detail-modal-footer">
              <AppButton variant="secondary" size="md" :disabled="rejectSubmitting" @click="closeRejectModal">
                ยกเลิก
              </AppButton>
              <AppButton
                :variant="rejectAction === 'reject' ? 'danger' : 'primary'"
                size="md"
                :loading="rejectSubmitting"
                :disabled="rejectSubmitting"
                @click="submitReject"
              >
                <span>{{ rejectAction === 'return_revision' ? 'ยืนยันส่งกลับแก้ไข' : 'ยืนยันไม่อนุมัติ' }}</span>
              </AppButton>
            </div>
          </div>
        </div>
      </Teleport>
    </main>


  </div>
  <ApplicationDocumentPreview :application="showApplicationDocument ? application : null" @close="showApplicationDocument = false" @signed="Object.assign(application, $event)" />
  <OfficerSignatureModal :open="showOfficerSignature" @cancel="showOfficerSignature = false" @confirm="submitOfficerApproval" />
</template>

<script setup lang="ts">
import OfficerSignatureModal from '~/components/OfficerSignatureModal.vue'
import ApplicationDocumentPreview from '~/components/ApplicationDocumentPreview.vue'
const showApplicationDocument = ref(false)
import { loadMarketGisJson, mergeMarketProjects } from '../../stall/data/market_loader'
import { getAllFloorLayoutsForProject } from '../../stall/data/market_stall_layout'

const route = useRoute()
const { get, post } = useApi()
const { showConfirm, showSuccess, showError } = useAppAlert()
const { returnUrl, backTo, crumbItems } = useReturnTo({ label: 'รายการคำร้อง', to: '/applications' })

const loading = ref(true)
const application = ref<any | null>(null)
const requestedStalls = ref<any[]>([])
const formatAmount = (value: number | null | undefined) => value == null ? '-' : Number(value).toLocaleString('th-TH')
const defaultStallPhotos = ['/images/stall_3d_perspective.jpg', '/images/market_stall_sample.jpg', '/images/market_hall_front.jpg']
const loadRequestedStalls = async () => {
  const app = application.value
  if (!app || !isStall.value) { requestedStalls.value = []; return }
  const codes = String(app.stall_ident || app.stall_number || '').split(',').map((v: string) => v.trim()).filter(Boolean)
  const market = String(app.market_id || '1')
  const [stallsResult, base, configResult] = await Promise.all([
    get('/stalls').catch(() => []),
    loadMarketGisJson(market).catch(() => ({ projects: [] })),
    get(`/markets/${encodeURIComponent(market)}/gis-config/draft`).catch(() => null)
  ])
  const records = Array.isArray(stallsResult) ? stallsResult : (stallsResult?.data || [])
  const config = configResult?.data || configResult || {}
  const projects = mergeMarketProjects(base.projects || [], Array.isArray(config.projects) ? config.projects : [])
  requestedStalls.value = codes.map((code: string) => {
    const record = Array.isArray(records) ? records.find((v: any) => String(v.ident || '').toUpperCase() === code.toUpperCase()) : null
    let item: any = null
    let project: any = null
    for (const candidate of projects) {
      const direct = (candidate.stalls || []).find((v: any) => String(v.code || v.ident || '').toUpperCase() === code.toUpperCase())
      const generated = Array.from(getAllFloorLayoutsForProject(candidate).values()).flatMap((layout: any) => layout.stalls).find((v: any) => String(v.code || v.ident || '').toUpperCase() === code.toUpperCase())
      if (direct || generated) { item = direct || generated; project = candidate; break }
    }
    const rent = item?.monthlyRent || record?.monthly_rent || app.stall_monthly_rate || 1800
    const photos = Array.isArray(item?.photos) && item.photos.length ? item.photos : defaultStallPhotos
    const floor = item?.floorName || item?.floor || record?.floor
    return { code, zone: item?.zone || project?.name || record?.project_name || record?.zone_name, floor: floor ? (String(floor).startsWith('ชั้น') ? String(floor) : 'ชั้น ' + floor) : '', category: item?.areaType || project?.category || record?.category, size: item?.size || (Number(item?.widthM || record?.width) * Number(item?.depthM || record?.depth)) || record?.area_sqm, entryFee: item?.entryFee || 75000, monthlyRent: rent, commonFee: item?.commonFee || 500, depositFee: item?.depositFee || rent, photos }
  })
}


const applicationNumber = computed(() => application.value?.application_number || application.value?.application_no || '')

const isStall = computed(() => {
  const app = application.value
  if (!app) return false
  if (app.application_type === 'stall_lease') return true
  if (app.application_type === 'general') return false
  if (String(app.application_number || app.application_no || '').startsWith('GEN-')) return false
  return Boolean(app.booking_id || (app.stall_number && app.stall_number !== '-') || (app.stall_ident && app.stall_ident !== '-'))
})

const status = computed(() => String(application.value?.status || application.value?.application_state || ''))

const statusLabel = computed(() => {
  switch (status.value) {
    case 'pending_officer': return 'รอเจ้าหน้าที่ตรวจสอบ'
    case 'pending_debt_check': return 'รอตรวจสอบหนี้'
    case 'debt_pending': return 'มีหนี้ค้างชำระ'
    case 'pending_market_director': return 'รอ ผอ.ตลาดพิจารณา'
    case 'pending_bureau_director': return 'รอ ผอ.สำนักตลาดพิจารณา'
    case 'approved': return 'อนุมัติแล้ว'
    case 'contract_pending': return 'รอตรวจสัญญา'
    case 'completed': return 'ทำสัญญาแล้ว'
    case 'revision_required': return 'ส่งกลับแก้ไข'
    case 'rejected': return 'ไม่อนุมัติ'
    default: return status.value || 'รอดำเนินการ'
  }
})

const statusClass = computed(() => {
  switch (status.value) {
    case 'approved':
    case 'completed': return 'tone-success'
    case 'rejected': return 'tone-danger'
    case 'pending_market_director':
    case 'contract_pending': return 'tone-blue'
    case 'pending_bureau_director': return 'tone-purple'
    default: return 'tone-warning'
  }
})

const createdAtLabel = computed(() => {
  const raw = application.value?.created_at
  if (!raw) return ''
  const d = new Date(raw)
  if (Number.isNaN(d.getTime())) return String(raw).slice(0, 10)
  return d.toLocaleString('th-TH', { day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' }) + ' น.'
})

const { isTenant, role } = useAuth()

// แก้ไขได้เฉพาะผู้เช่าเมื่อเจ้าหน้าที่ส่งคำร้องกลับ (revision_required)
const canRevise = computed(() => isTenant.value && isStall.value && status.value === 'revision_required')

// ขั้นตอนถัดไปเป็นหน้าที่ของผู้ใช้ที่กำลังดูอยู่หรือไม่ (ใช้แสดงป้าย "คุณต้องทำ" เท่านั้น)
const isMyTurn = computed(() => {
  const r = String(role.value || '')
  switch (status.value) {
    case 'pending_officer':
    case 'pending_debt_check':
    case 'contract_pending': return ['staff', 'admin'].includes(r)
    case 'approved': return isStall.value && ['staff', 'admin'].includes(r)
    case 'pending_market_director': return r === 'market_director'
    case 'pending_bureau_director': return r === 'bureau_director'
    case 'revision_required': return r === 'tenant'
    default: return false
  }
})

// 4 ขั้นตอน เหมือนหน้ารายการคำร้อง (ApplicationMain)
const steps = computed(() => [
  {
    number: 1,
    label: status.value === 'revision_required'
      ? (isTenant.value ? 'แก้ไขคำร้องขอเช่า' : 'รอผู้เช่าแก้ไขคำร้อง')
      : 'ยื่นคำร้องขอเช่า'
  },
  { number: 2, label: 'จนท. ตรวจสอบ/หนี้' },
  { number: 3, label: 'ผอ.ตลาด พิจารณา' },
  { number: 4, label: 'ผอ.สำนัก อนุมัติ' }
])

const isStepCompleted = (step: number): boolean => {
  const s = status.value
  if (s === 'revision_required') return false
  if (step === 1) return true
  if (step === 2) return ['pending_market_director', 'pending_bureau_director', 'approved', 'contract_pending', 'completed'].includes(s)
  if (step === 3) return ['pending_bureau_director', 'approved', 'contract_pending', 'completed'].includes(s)
  if (step === 4) return ['approved', 'contract_pending', 'completed'].includes(s)
  return false
}

const isStepActive = (step: number): boolean => {
  const s = status.value
  if (s === 'revision_required') return step === 1
  if (step === 2) return ['pending_officer', 'pending_debt_check', 'debt_pending'].includes(s)
  if (step === 3) return s === 'pending_market_director'
  if (step === 4) return s === 'pending_bureau_director'
  return false
}

const isConnectorFilled = (afterStep: number): boolean => {
  if (afterStep === 1) return isStepCompleted(2) || (isStepCompleted(1) && status.value !== 'revision_required')
  return isStepCompleted(afterStep + 1)
}

const nextAction = computed(() => {
  switch (status.value) {
    case 'pending_officer':
      return { tone: 'waiting', owner: 'เจ้าหน้าที่ตลาด', title: 'รอเจ้าหน้าที่ตรวจสอบเอกสาร', description: 'เจ้าหน้าที่กำลังตรวจสอบความครบถ้วนของเอกสารและข้อมูลผู้ยื่น จากนั้นจะส่งต่อให้ ผอ.ตลาดพิจารณา' }
    case 'pending_debt_check':
      return { tone: 'waiting', owner: 'เจ้าหน้าที่การคลัง / ตลาด', title: 'รอตรวจสอบประวัติหนี้ค้าง', description: 'ระบบกำลังรอการตรวจสอบประวัติหนี้ค้างชำระของผู้ยื่นคำร้อง' }
    case 'debt_pending':
      return { tone: 'danger', owner: 'ผอ.ตลาด / การคลัง', title: 'มีหนี้ค้างชำระ', description: 'ต้องเคลียร์หนี้ค้างชำระให้เรียบร้อยก่อนจึงจะพิจารณาคำร้องต่อได้' }
    case 'pending_market_director':
      return { tone: 'waiting', owner: 'ผอ.ตลาด', title: 'รอ ผอ.ตลาดให้ความเห็นชอบ', description: 'คำร้องผ่านการตรวจสอบเอกสารแล้ว กำลังรอ ผอ.ตลาดพิจารณาและให้ความเห็นชอบ' }
    case 'pending_bureau_director':
      return { tone: 'waiting', owner: 'ผอ.สำนักตลาด', title: 'รอ ผอ.สำนักตลาดอนุมัติขั้นสุดท้าย', description: 'คำร้องผ่านความเห็นชอบจาก ผอ.ตลาดแล้ว กำลังรอ ผอ.สำนักตลาดอนุมัติขั้นสุดท้าย' }
    case 'approved':
      return isStall.value
        ? { tone: 'action', owner: 'เจ้าหน้าที่ตลาด', title: 'รอจัดทำสัญญาเช่า', description: 'คำร้องได้รับอนุมัติแล้ว เจ้าหน้าที่จะจัดทำสัญญาเช่าแผงค้าในขั้นตอนถัดไป' }
        : { tone: 'complete', owner: 'ผอ.สำนักตลาด', title: 'คำร้องได้รับอนุมัติแล้ว', description: 'คำร้องทั่วไปนี้ผ่านการพิจารณาครบทุกขั้นตอนแล้ว' }
    case 'contract_pending':
      return { tone: 'waiting', owner: 'เจ้าหน้าที่ตลาด', title: 'รอตรวจสอบสัญญาเช่า', description: 'สัญญาเช่าถูกจัดทำแล้ว กำลังรอเจ้าหน้าที่ตรวจสอบและอนุมัติสัญญา' }
    case 'completed':
      return { tone: 'complete', owner: '-', title: 'ดำเนินการครบถ้วนแล้ว', description: 'คำร้องนี้ได้ทำสัญญาเช่าเรียบร้อยแล้ว' }
    case 'revision_required':
      return { tone: 'action', owner: 'ผู้ยื่นคำร้อง', title: 'ผู้ยื่นต้องแก้ไขเอกสาร', description: 'คำร้องถูกส่งกลับให้แก้ไข ผู้ยื่นต้องแก้ไขข้อมูลหรือเอกสารตามหมายเหตุแล้วส่งตรวจสอบอีกครั้ง' }
    case 'rejected':
      return { tone: 'danger', owner: '-', title: 'คำร้องไม่ได้รับอนุมัติ', description: 'คำร้องนี้ไม่ผ่านการพิจารณา ดูเหตุผลได้ในบันทึกความเห็นด้านล่าง' }
    default:
      return { tone: 'waiting', owner: '-', title: 'รอดำเนินการ', description: 'คำร้องอยู่ระหว่างการดำเนินการ' }
  }
})

const comments = computed(() => {
  const app = application.value
  const list: { role: string; text: string; step: string }[] = []
  if (!app) return list
  if (app.officer_comment) list.push({ role: 'เจ้าหน้าที่ตลาด', text: app.officer_comment, step: 'ขั้นตรวจสอบเอกสาร' })
  if (app.market_director_comment) list.push({ role: 'ผอ.ตลาด', text: app.market_director_comment, step: 'ขั้นพิจารณา ผอ.ตลาด' })
  if (app.bureau_director_comment) list.push({ role: 'ผอ.สำนักตลาด', text: app.bureau_director_comment, step: 'ขั้นพิจารณา ผอ.สำนัก' })
  if (app.review_note && !list.some(l => l.text === app.review_note)) list.push({ role: 'บันทึกคำร้อง', text: app.review_note, step: 'หมายเหตุคำร้อง' })
  return list
})

const formatThaiDateTime = (value: string | undefined | null) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value)
  return new Intl.DateTimeFormat('th-TH', { dateStyle: 'long', timeStyle: 'short', timeZone: 'Asia/Bangkok' }).format(date).replace(' น.', '') + ' น.'
}

// ประวัติการดำเนินการ: รูปแบบไทม์ไลน์เดียวกับหน้า bookings
const timeline = computed(() => {
  const app = application.value
  const list: { title: string; time: string; actor: string; note: string; tone: string }[] = []
  if (!app) return list

  const stallSuffix = (app.stall_number || app.stall_ident) ? ` (แผง ${app.stall_number || app.stall_ident})` : ''
  const timeCreated = formatThaiDateTime(app.submitted_at || app.created_at)

  // 1. ลงทะเบียนยื่นคำร้อง
  list.push({
    title: isStall.value
      ? `ลงทะเบียนยื่นคำร้องขอเช่าแผงค้า${stallSuffix}`
      : `ยื่นคำร้องทั่วไป`,
    time: timeCreated,
    actor: `ผู้เช่า · ${app.applicant_name || '-'}`,
    note: app.details || (isStall.value ? `ยื่นคำร้องขอเช่าแผงค้าเรียบร้อยแล้ว รอเจ้าหน้าที่ตลาดตรวจสอบเอกสาร` : 'ยื่นคำร้องเรียบร้อยแล้ว'),
    tone: 'tenant'
  })

  const reviewedTime = formatThaiDateTime(app.reviewed_at || app.updated_at || app.created_at)

  // 2. เจ้าหน้าที่ตรวจสอบเอกสาร / หนี้
  const hasOfficerDone = app.officer_comment || ['pending_market_director', 'pending_bureau_director', 'approved', 'contract_pending', 'completed'].includes(status.value)
  if (hasOfficerDone) {
    list.push({
      title: `เจ้าหน้าที่ตรวจสอบเอกสาร${stallSuffix}`,
      time: reviewedTime || timeCreated,
      actor: 'เจ้าหน้าที่ · เจ้าหน้าที่ตลาด',
      note: app.officer_comment || 'ตรวจสอบเอกสารครบถ้วน ส่งต่อ ผอ.ตลาด',
      tone: 'officer'
    })
  } else if (status.value === 'revision_required') {
    list.push({
      title: `เจ้าหน้าที่ส่งกลับให้แก้ไขคำร้อง${stallSuffix}`,
      time: reviewedTime || timeCreated,
      actor: 'เจ้าหน้าที่ · เจ้าหน้าที่ตลาด',
      note: app.officer_comment || app.review_note || 'เอกสารไม่ครบถ้วน ส่งกลับให้ผู้เช่าแก้ไข',
      tone: 'warning'
    })
  }

  // 3. ผอ.ตลาด ให้ความเห็นชอบ
  const hasMarketDone = app.market_director_comment || ['pending_bureau_director', 'approved', 'contract_pending', 'completed'].includes(status.value)
  if (hasMarketDone) {
    list.push({
      title: `ผอ.ตลาด พิจารณาให้ความเห็นชอบ${stallSuffix}`,
      time: reviewedTime || timeCreated,
      actor: 'ผอ.ตลาด · ผู้อำนวยการตลาด',
      note: app.market_director_comment || 'เห็นชอบคำร้อง ส่งต่อ ผอ.สำนักตลาด พิจารณาอนุมัติขั้นสุดท้าย',
      tone: 'officer'
    })
  }

  // 4. ผอ.สำนักตลาด พิจารณาอนุมัติ / ไม่อนุมัติ
  const isApproved = ['approved', 'contract_pending', 'completed'].includes(status.value)
  const isRejectedStatus = status.value === 'rejected'
  if (app.bureau_director_comment || isApproved || isRejectedStatus) {
    list.push({
      title: isApproved
        ? `ผอ.สำนักตลาด พิจารณาอนุมัติคำร้อง${stallSuffix}`
        : `ผอ.สำนักตลาด ไม่อนุมัติคำร้อง${stallSuffix}`,
      time: reviewedTime || timeCreated,
      actor: 'ผอ.สำนักตลาด · ผู้อำนวยการสำนักตลาด',
      note: app.bureau_director_comment || (isApproved ? 'อนุมัติคำร้องเรียบร้อย' : 'คำร้องไม่ผ่านการอนุมัติ'),
      tone: isApproved ? 'officer' : 'danger'
    })
  }

  // 5. ทำสัญญาแล้ว
  if (status.value === 'completed') {
    list.push({
      title: `จัดทำสัญญาเช่าเรียบร้อยแล้ว${stallSuffix}`,
      time: formatThaiDateTime(app.updated_at || app.created_at),
      actor: 'เจ้าหน้าที่ · เจ้าหน้าที่ตลาด',
      note: 'จัดทำสัญญาเช่าและลงนามในระบบเรียบร้อยแล้ว',
      tone: 'officer'
    })
  }

  // 6. หมายเหตุอื่นๆ (ถ้ามี)
  if (app.review_note && !list.some(s => s.note === app.review_note)) {
    list.push({
      title: 'หมายเหตุการดำเนินการ',
      time: reviewedTime || timeCreated,
      actor: 'ระบบ',
      note: app.review_note,
      tone: 'system'
    })
  }

  return list
})
const actionSubmitting = ref(false)

const canReview = computed(() => {
  const app = application.value
  if (!app) return false
  const r = String(role.value || '')
  return (app.status === 'pending_officer' && ['staff', 'admin'].includes(r)) ||
    (app.status === 'pending_market_director' && r === 'market_director') ||
    (app.status === 'pending_bureau_director' && r === 'bureau_director')
})

const approveButtonLabel = computed(() => {
  const s = status.value
  if (s === 'pending_officer') return 'ตรวจสอบรับรอง'
  if (s === 'pending_market_director') return 'ผอ. อนุมัติ'
  if (s === 'pending_bureau_director') return 'ผอ.สำนัก อนุมัติ'
  return 'อนุมัติคำร้อง'
})

const defaultApproveComment = computed(() => {
  const s = status.value
  if (s === 'pending_officer') return 'เจ้าหน้าที่ตรวจสอบเอกสารครบถ้วน ส่งต่อ ผอ.ตลาด'
  if (s === 'pending_market_director') return 'ผอ.ตลาด พิจารณาอนุมัติคำร้อง'
  if (s === 'pending_bureau_director') return 'ผอ.สำนักตลาด พิจารณาอนุมัติคำร้องเรียบร้อย'
  return 'อนุมัติเรียบร้อย'
})

const showOfficerSignature = ref(false)
const submitOfficerApproval = async (signature: string) => { showOfficerSignature.value = false; await processApproval(signature) }
const processApproval = async (signature = '') => {
  if (isStall.value && status.value === 'pending_officer' && !signature) { showOfficerSignature.value = true; return }
  if (!application.value || actionSubmitting.value) return
  const app = application.value
  const appNo = applicationNumber.value
  const confirmed = await showConfirm(
    `คุณต้องการพิจารณาอนุมัติคำร้อง ${appNo} หรือไม่?`,
    {
      title: 'ยืนยันการอนุมัติคำร้อง',
      type: 'info',
      confirmText: 'ยืนยันอนุมัติ',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return
  actionSubmitting.value = true
  try {
    await post(`/applications/${app.id}/action`, {
      action: 'approve',
      comment: defaultApproveComment.value,
      signature_base64: signature
    })
    await showSuccess(`อนุมัติคำร้อง ${appNo} เรียบร้อยแล้ว`)
    await fetchDetail()
  } catch (err: any) {
    await showError(err?.data?.detail || err?.data?.message || err?.message || 'เกิดข้อผิดพลาดในการอนุมัติคำร้อง')
  } finally {
    actionSubmitting.value = false
  }
}

// Reject / Return Revision Modal State
const showRejectModal = ref(false)
const rejectAction = ref<'return_revision' | 'reject'>('return_revision')
const rejectComment = ref('')
const rejectErrorMsg = ref('')
const rejectSubmitting = ref(false)

const quickRejectReasons = [
  'เอกสารสำเนาบัตรประชาชนไม่ชัดเจน',
  'ประเภทสินค้าไม่สอดคล้องกับผังโซนที่จัดสรร',
  'ข้อมูลร้านค้าหรือเบอร์โทรศัพท์ติดต่อไม่ถูกต้อง',
  'ขอให้แนบเอกสารเพิ่มเติมเพื่อประกอบการพิจารณา'
]

const openRejectModal = (action: 'return_revision' | 'reject') => {
  rejectAction.value = action
  rejectComment.value = ''
  rejectErrorMsg.value = ''
  showRejectModal.value = true
}

const closeRejectModal = () => {
  showRejectModal.value = false
  rejectComment.value = ''
  rejectErrorMsg.value = ''
}

const setQuickReason = (reason: string) => {
  rejectComment.value = reason
  rejectErrorMsg.value = ''
}

const submitReject = async () => {
  if (!rejectComment.value.trim()) {
    rejectErrorMsg.value = rejectAction.value === 'return_revision'
      ? 'กรุณาระบุรายละเอียดหรือสิ่งที่ต้องแก้ไข'
      : 'กรุณาระบุเหตุผลการไม่อนุมัติคำร้อง'
    return
  }
  if (!application.value) return
  const app = application.value
  const appNo = applicationNumber.value
  rejectSubmitting.value = true
  try {
    await post(`/applications/${app.id}/action`, {
      action: rejectAction.value,
      comment: rejectComment.value.trim(),
      reason: rejectComment.value.trim()
    })
    closeRejectModal()
    await showSuccess(
      rejectAction.value === 'return_revision'
        ? `ส่งกลับคำร้อง ${appNo} ให้ผู้เช่าแก้ไขเรียบร้อยแล้ว`
        : `บันทึกผลการไม่อนุมัติคำร้อง ${appNo} เรียบร้อยแล้ว`
    )
    await fetchDetail()
  } catch (err: any) {
    rejectErrorMsg.value = err?.data?.detail || err?.data?.message || err?.message || 'เกิดข้อผิดพลาดในการดำเนินการ'
  } finally {
    rejectSubmitting.value = false
  }
}

const fetchDetail = async () => {
  try {
    const res: any = await get(`/applications/${route.params.id}`)
    application.value = res?.id ? res : (res?.data || null)
    await loadRequestedStalls()
  } catch {
    application.value = null
  } finally {
    loading.value = false
  }
}

onMounted(fetchDetail)
</script>

<style scoped>
.app-detail-page { min-height: 100vh; background: #f8fafc; color: #0f172a; }
.detail-container { display: flex; flex-direction: column; gap: 16px; }

.detail-hero, .progress-card, .detail-empty, .detail-skeleton, .summary-section, .next-action-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
}

/* Hero */
.detail-hero { display: flex; align-items: center; justify-content: space-between; gap: 20px; padding: 22px 24px; }
.detail-title-wrap, .hero-actions { display: flex; align-items: center; gap: 14px; }
.detail-title-wrap { min-width: 0; }
.detail-title-icon { width: 46px; height: 46px; display: grid; place-items: center; border-radius: 12px; background: #eff6ff; border: 1px solid #dbeafe; color: #2563eb; flex-shrink: 0; }
.eyebrow { display: block; color: #64748b; font-size: 0.78rem; font-weight: 600; margin-bottom: 3px; }
.detail-hero h1 { margin: 0; font-size: 1.65rem; line-height: 1.1; font-weight: 800; overflow-wrap: anywhere; }
.created-at { margin-top: 5px; color: #94a3b8; font-size: 0.8125rem; }

.status-badge { padding: 6px 13px; border-radius: 999px; font-size: 0.8125rem; font-weight: 700; white-space: nowrap; }
.tone-warning { background: #fef3c7; color: #b45309; }
.tone-success { background: #dcfce7; color: #15803d; }
.tone-blue { background: #e0f2fe; color: #0369a1; }
.tone-purple { background: #ede9fe; color: #6d28d9; }
.tone-danger { background: #fee2e2; color: #dc2626; }

/* Next action (same design as the bookings detail page) */
.next-action-card { display: flex; align-items: center; flex-wrap: wrap; gap: 14px; min-height: 76px; padding: 15px 18px; border: 1px solid #bfdbfe; border-left: 4px solid #2563eb; border-radius: 14px; background: #fff; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03); }
.next-action-icon { width: 40px; height: 40px; flex: 0 0 40px; display: grid; place-items: center; border-radius: 11px; background: #eff6ff; color: #2563eb; }
.next-action-copy { flex: 1; min-width: 0; }
.next-action-meta { display: flex; align-items: center; gap: 10px; margin-bottom: 6px; color: #64748b; font-size: 0.75rem; font-weight: 650; }
.next-action-owner { padding: 2px 7px; border-radius: 999px; background: #f1f5f9; color: #64748b; font-size: 0.75rem; line-height: 1.35; font-weight: 600; white-space: nowrap; }
.next-action-owner.is-mine { background: #dbeafe; color: #1d4ed8; }
.next-action-copy h2 { margin: 0; color: #0f172a; font-size: 1rem; line-height: 1.35; font-weight: 700; }
.next-action-copy p { margin: 6px 0 0; color: #64748b; font-size: 0.8rem; line-height: 1.45; }
.next-action-buttons { flex: 0 0 auto; display: flex; align-items: center; justify-content: flex-end; gap: 9px; flex-wrap: wrap; margin-left: auto; }
.next-action-waiting { border-color: #bae6fd; border-left-color: #0ea5e9; }
.next-action-waiting .next-action-icon { background: #f0f9ff; color: #0284c7; }
.next-action-complete { border-color: #e2e8f0; border-left-color: #94a3b8; }
.next-action-complete .next-action-icon { background: #f1f5f9; color: #64748b; }
.next-action-danger { border-color: #fecaca; border-left-color: #dc2626; }
.next-action-danger .next-action-icon { background: #fee2e2; color: #dc2626; }

/* Section order: next action first, then hero (as on the bookings detail page) */
.detail-container > .next-action-card { order: 1; }
.detail-container > .detail-hero { order: 2; }
.detail-container > .progress-card { order: 3; }
.detail-container > .application-information { order: 4; }
/* Progress (4-step stepper, same as the applications list) */
.progress-card { padding: 8px; }
.booking-stepper-container { margin: 0; padding: 14px 18px; background: #f8fafc; border-radius: 10px; border: 1px solid #f1f5f9; }
.stepper-track { display: flex; align-items: center; justify-content: space-between; position: relative; max-width: 680px; margin: 0 auto; }
.stepper-step { display: flex; flex-direction: column; align-items: center; gap: 6px; z-index: 2; position: relative; }
.step-circle { width: 26px; height: 26px; border-radius: 50%; background: #e2e8f0; color: #64748b; font-size: 0.75rem; font-weight: 700; display: flex; align-items: center; justify-content: center; transition: all 0.25s ease; }
.step-label-text { font-size: 0.75rem; font-weight: 600; color: #64748b; white-space: nowrap; }
.stepper-step.completed .step-circle { background: #10b981; color: #ffffff; }
.stepper-step.completed .step-label-text { color: #0f172a; }
.stepper-step.active .step-circle { background: #2563eb; color: #ffffff; box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.15); }
.stepper-step.active .step-label-text { color: #2563eb; font-weight: 700; }
.stepper-connector { flex: 1; height: 3px; background: #e2e8f0; margin: -18px 8px 0 8px; z-index: 1; transition: background 0.25s ease; }
.stepper-connector.filled { background: #10b981; }
/* Information */
.requested-stall { margin-bottom: 16px; }
.stall-heading { display: flex; flex-wrap: wrap; justify-content: flex-end; gap: 6px; margin-left: auto; }
.stall-heading span { display: inline-block; padding: 6px 14px; border-radius: 6px; background: #2563eb; color: white; font-size: .875rem; font-weight: 700; }
.stall-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 18px 16px; margin: 18px 0; }
.stall-grid dt, .stall-photos > div:first-child { color: #64748b; font-size: .8125rem; }
.stall-grid dd { margin: 5px 0 0; color: #1e293b; font-size: .9375rem; font-weight: 700; overflow-wrap: anywhere; }
.stall-floor { display: inline-block; padding: 2px 6px; border: 1px solid #bfdbfe; border-radius: 4px; background: #eff6ff; color: #2563eb; font-size: .75rem; }
.stall-photos { padding-top: 16px; border-top: 1px solid #f1f5f9; }
.stall-photo-list { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 10px; }
.stall-photo-list img { width: 110px; height: 78px; border-radius: 8px; object-fit: cover; }
.stall-extra { margin-top: 16px; }
@media (max-width: 1200px) { .stall-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }

.application-information-grid { display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1fr); gap: 16px; align-items: stretch; }
.application-information-main { display: grid; gap: 16px; min-width: 0; align-content: start; }
.summary-section { padding: 24px; }
.section-heading { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; padding-bottom: 12px; border-bottom: 1px solid #f1f5f9; }
.section-heading h3 { margin: 0; font-size: 1.125rem; font-weight: 700; color: #0f172a; }
.section-icon { width: 40px; height: 40px; flex: 0 0 40px; display: grid; place-items: center; border-radius: 10px; background: #eff6ff; color: #2563eb; }
.detail-list { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px 20px; margin: 0; }
.detail-list .full { grid-column: 1 / -1; }
.application-detail-list { grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 16px 24px; }
.application-detail-list dd { overflow-wrap: anywhere; }
.detail-list dt { font-size: 0.8125rem; line-height: 1.4; color: #64748b; font-weight: 500; }
.detail-list dd { margin: 3px 0 0; font-size: 0.875rem; font-weight: 700; line-height: 1.45; color: #0f172a; word-break: break-word; }

.timeline-section-heading { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 16px; padding-bottom: 14px; border-bottom: 1px solid #e2e8f0; }
.timeline-heading-text { display: flex; flex-direction: column; gap: 2px; }
.timeline-eyebrow { font-size: 0.75rem; color: #64748b; font-weight: 600; }
.timeline-section-heading h3 { margin: 0; font-size: 1.15rem; font-weight: 800; color: #0f172a; }

.timeline-list { display: flex; flex-direction: column; list-style: none; margin: 0; padding: 4px 8px 0 0; max-height: 600px; overflow-y: auto; overscroll-behavior: contain; }
.timeline-item { position: relative; display: flex; gap: 14px; padding-bottom: 22px; }
.timeline-item:last-child { padding-bottom: 0; }
.timeline-item:not(:last-child)::before { content: ''; position: absolute; left: 5px; top: 14px; bottom: 0; width: 1px; background: #e2e8f0; }
.timeline-dot { position: relative; z-index: 1; width: 11px; height: 11px; margin-top: 5px; flex: 0 0 11px; border-radius: 50%; background: #3b82f6; box-shadow: 0 0 0 4px #dbeafe; }
.timeline-dot.tenant { background: #3b82f6; box-shadow: 0 0 0 4px #dbeafe; }
.timeline-dot.officer { background: #10b981; box-shadow: 0 0 0 4px #d1fae5; }
.timeline-dot.director { background: #10b981; box-shadow: 0 0 0 4px #d1fae5; }
.timeline-dot.warning { background: #f59e0b; box-shadow: 0 0 0 4px #fef3c7; }
.timeline-dot.danger { background: #ef4444; box-shadow: 0 0 0 4px #fee2e2; }
.timeline-dot.system { background: #8b5cf6; box-shadow: 0 0 0 4px #ede9fe; }

.timeline-copy { flex: 1; min-width: 0; }
.timeline-copy strong { display: block; font-size: 0.925rem; font-weight: 700; color: #0f172a; line-height: 1.35; }
.timeline-copy time { display: block; color: #64748b; font-size: 0.78rem; margin-top: 4px; }
.timeline-copy .timeline-actor { display: block; color: #64748b; font-size: 0.78rem; margin-top: 2px; }
.timeline-copy p { margin: 8px 0 0; padding: 10px 14px; border-radius: 8px; background: #f8fafc; color: #334155; font-size: 0.84rem; line-height: 1.5; border: 1px solid #f1f5f9; overflow-wrap: anywhere; }
.empty-text { margin: 0; font-size: 0.875rem; color: #94a3b8; }


.submitted-notice { display: flex; align-items: center; flex-wrap: wrap; gap: 8px 12px; margin-top: 14px; padding: 10px 14px; border-radius: 10px; background: #f8fafc; border: 1px solid #e2e8f0; color: #475569; font-size: 0.8125rem; line-height: 1.45; }
.submitted-notice.is-revision { justify-content: space-between; background: #fffbeb; border-color: #fde68a; color: #92400e; }
/* Empty & skeleton */
.detail-empty { display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 60px 16px; color: #64748b; }
.detail-empty h1 { margin: 0; font-size: 1.25rem; color: #1e293b; }
.detail-empty p { margin: 0; font-size: 0.875rem; }
.detail-skeleton { padding: 24px; display: grid; gap: 14px; }
.skeleton-line, .skeleton-panel, .skeleton-grid > div { border-radius: 8px; background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%); background-size: 200% 100%; animation: app-detail-shimmer 1.5s infinite; }
.skeleton-line { height: 22px; }
.skeleton-line.wide { width: 46%; }
.skeleton-line.short { width: 28%; height: 16px; }
.skeleton-panel { height: 84px; }
.skeleton-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 14px; }
.skeleton-grid > div { height: 120px; }
@keyframes app-detail-shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

@media (max-width: 900px) {
  .application-information-grid, .detail-list { grid-template-columns: 1fr; }

  .detail-hero { flex-direction: column; align-items: flex-start; }
  .next-action-buttons { margin-left: 0; width: 100%; }
}

@media (max-width: 1200px) and (min-width: 901px), (max-width: 700px) { .application-detail-list { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (max-width: 540px) { .application-detail-list { grid-template-columns: 1fr; } }

/* Action Modal Styles */
.detail-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
}

.detail-modal-content {
  width: 100%;
  max-width: 520px;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  animation: modalPopIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes modalPopIn {
  from { opacity: 0; transform: scale(0.96) translateY(8px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.detail-modal-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 18px 20px;
  border-bottom: 1px solid #f1f5f9;
  background: #fffbeb;
}

.detail-modal-header.is-danger {
  background: #fef2f2;
}

.modal-icon-badge {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  background: #fef3c7;
  color: #d97706;
  flex-shrink: 0;
}

.detail-modal-header.is-danger .modal-icon-badge {
  background: #fee2e2;
  color: #dc2626;
}

.modal-header-titles {
  flex: 1;
  min-width: 0;
}

.modal-title {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 700;
  color: #0f172a;
}

.modal-sub {
  margin: 2px 0 0;
  font-size: 0.8125rem;
  color: #64748b;
}

.btn-modal-close {
  width: 32px;
  height: 32px;
  display: grid;
  place-items: center;
  border-radius: 8px;
  border: 1px solid transparent;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-modal-close:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.detail-modal-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-input-label {
  display: block;
  font-size: 0.84rem;
  font-weight: 600;
  color: #334155;
  margin-bottom: 6px;
}

.modal-input-label.required::after {
  content: ' *';
  color: #dc2626;
}

.quick-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.quick-chip {
  padding: 4px 10px;
  border-radius: 6px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #475569;
  font-size: 0.78rem;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}

.quick-chip:hover {
  background: #eff6ff;
  border-color: #bfdbfe;
  color: #2563eb;
}

.reason-textarea {
  width: 100%;
  box-sizing: border-box;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  font-family: inherit;
  font-size: 0.875rem;
  color: #0f172a;
  outline: none;
  resize: vertical;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.reason-textarea:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
}

.modal-error-msg {
  margin: 5px 0 0;
  font-size: 0.8rem;
  color: #dc2626;
}

.detail-modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 20px;
  border-top: 1px solid #f1f5f9;
  background: #f8fafc;
}

</style>
