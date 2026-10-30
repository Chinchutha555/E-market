<template>
  <div>
    <!-- Top Header -->
    <AppHeader />

    <div class="content-body">
      <!-- Breadcrumb Navigation at top of section body -->
      <AppBreadcrumb
        :market="currentMarketName"
        :marketId="route.query.market || '1'"
        :items="bookingBreadcrumbItems"
        :current="currentBreadcrumbText"
        :backTo="`/stalls?market=${route.query.market || '1'}`"
        :onBack="handleBreadcrumbBack"
      />

      <!-- Market Info Header Banner with Quick Market Switcher -->
      <div class="market-select-banner card glass-card" style="display: flex; justify-content: space-between; align-items: center; border-radius: 16px; padding: 14px 22px; margin-bottom: 32px; flex-wrap: wrap; gap: 14px; border: 1px solid #bfdbfe; background: linear-gradient(135deg, #f0f7ff 0%, #ffffff 100%); box-shadow: 0 4px 16px rgba(37, 99, 235, 0.06);">
        <div style="display: flex; align-items: center; gap: 14px; flex-wrap: wrap;">
          <div style="width: 42px; height: 42px; border-radius: 12px; background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%); color: #ffffff; display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 12px rgba(37, 99, 235, 0.25);">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
              <polyline points="9 22 9 12 15 12 15 22"></polyline>
            </svg>
          </div>
          <div>
            <div style="font-size: 0.74rem; color: #2563eb; font-weight: 800; text-transform: uppercase; letter-spacing: 0.04em;">ตลาดที่ท่านกำลังทำการจองแผงเช่า</div>
            <div style="display: flex; align-items: center; gap: 10px; margin-top: 2px;">
              <strong style="font-size: 1.25rem; color: #0f172a; font-weight: 800;">{{ currentMarketName }}</strong>
              <span style="font-size: 0.8rem; color: #475569; background: #ffffff; padding: 3px 10px; border-radius: 20px; border: 1px solid #dbeafe; display: inline-flex; align-items: center; gap: 5px; font-weight: 600; box-shadow: 0 1px 3px rgba(0,0,0,0.04);">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                  <circle cx="12" cy="10" r="3"></circle>
                </svg>
                <span>{{ currentMarketDistrict }}</span>
              </span>
            </div>
          </div>
        </div>

        <div style="display: flex; align-items: center; gap: 12px; flex-wrap: wrap;">
          <NuxtLink to="/" class="btn btn-sm btn-outline" style="border-color: #cbd5e1; color: #1e40af; background: #ffffff; font-weight: 700; display: inline-flex; align-items: center; gap: 7px; padding: 8px 16px; border-radius: 10px; box-shadow: 0 1px 3px rgba(0,0,0,0.04);">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"></polygon>
              <line x1="8" y1="2" x2="8" y2="18"></line>
              <line x1="16" y1="6" x2="16" y2="22"></line>
            </svg>
            <span>เลือกจองตลาดอื่น</span>
          </NuxtLink>
        </div>
      </div>

      <!-- Stepper Component (3 Steps) -->
      <div class="booking-stepper-wrap">
        <BookingStepper :currentStep="currentStep" @changeStep="setStep" />
      </div>

      <!-- ============================================================ -->
      <!-- STEP 1: เลือกแผงเช่าที่สนใจ (Projects Overview & Sub-Stalls)  -->
      <!-- ============================================================ -->
      <div v-show="currentStep === 1">
        <!-- แจ้งเตือนเมื่อตลาดเต็ม (ใช้ข้อมูลจริง) -->
        <div v-if="isMarketFull" class="market-full-banner">
          <div class="market-full-banner-icon">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
          </div>
          <div class="market-full-banner-content">
            <h3 class="market-full-title">ตลาดนี้มีผู้เช่าเต็มแล้ว (เปิดให้เข้าชมผังแผงเท่านั้น)</h3>
            <p class="market-full-desc">ขณะนี้แผงค้าใน{{ currentMarketName }}มีผู้เช่าครบทุกแผงแล้ว จึงยังไม่สามารถเปิดรับการจองแผงได้ในขณะนี้ ท่านยังสามารถเข้าชมผังแผงและข้อมูลโซนได้ แต่แผงค้าจะไม่สามารถคลิกเลือกเพื่อทำการจองได้</p>
          </div>
        </div>

        <!-- Title & Filter Header -->
        <div class="step1-title-bar">
          <div>
            <h2 style="font-size: 1.25rem; font-weight: 800; color: #0f172a; margin-bottom: 4px;">
              เลือกแผงเช่าที่สนใจ
            </h2>
            <div style="font-size: 0.82rem; color: #64748b;">
              คลิกเลือกโครงการหรือแผงค้าบนแผนที่เพื่อดูรายละเอียดและทำรายการจอง
            </div>
            <div style="height: 3px; width: 48px; background: var(--primary); border-radius: var(--radius-full); margin-top: 6px;"></div>
          </div>
          <div id="step1-map-toolbar" class="step1-map-toolbar"></div>
        </div>

        <!-- ========================================================== -->
        <!-- GIS PROJECTS & STALLS MAP                                  -->
        <!-- ========================================================== -->
        <div style="padding: 0; position: relative; overflow: hidden; margin-bottom: 24px;">

          <!-- The SVG Master Map of Projects Component -->
          <MarketProjectsOverviewMap
            ref="overviewMapRef"
            :marketId="currentMarketId"
            :marketFull="isMarketFull"
            :selectedProjectId="selectedProject?.id"
            :selectedStallCodes="selectedStalls.map(s => s.code)"
            :selectedStallKeys="selectedStalls.map(stallSelectionKey)"
            :myBookedStallCodes="Object.keys(myActiveBookingsByStall)"
            :stallStatusOverrides="stallStatusOverrides"
            v-model:onlyAvailable="onlyAvailable"
            :categoryFilter="selectedCategoryFilter"
            v-model:mapStyle="mapStyle"
            @selectProject="handleSelectProject"
            @exitProject="handleExitProject"
            @selectStall="handleMapStallSelect"
            @update:categoryFilter="val => selectedCategoryFilter = val"
          >
            <!-- Selected Stall Drawer (UI เดิม) ลอยอยู่บนแผนที่ด้านขวา -->
            <template #overlay>
              <transition name="slide-fade">
              <div
                v-if="selectedStalls.length > 0"
                id="map-selected-stall-drawer"
                class="selected-stall-drawer card"
                style="
                  position: absolute;
                  top: 62px;
                  right: 14px;
                  width: 360px;
                  max-height: calc(100% - 76px);
                  display: flex;
                  flex-direction: column;
                  pointer-events: auto;
                  z-index: 25;
                  padding: 0;
                  box-shadow: 0 12px 36px rgba(15, 23, 42, 0.16);
                  backdrop-filter: blur(14px);
                  background: rgba(255, 255, 255, 0.98);
                  border-radius: 16px;
                  border: 1px solid #e2e8f0;
                  overflow: hidden;
                "
              >
                <!-- 1. Drawer Header (Pinned at Top) -->
                <div style="padding: 14px 16px 10px; border-bottom: 1px solid #f1f5f9; background: #ffffff; flex-shrink: 0;">
                  <div style="display: flex; justify-content: space-between; align-items: center;">
                    <strong style="font-size: 1.125rem; color: #0f172a; font-weight: 700;">แผงค้าที่เลือก</strong>
                    <button class="btn-icon-square" style="width: 26px; height: 26px;" @click="clearAllSelectedStalls" title="ยกเลิกการเลือก" aria-label="ล้าง">
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                    </button>
                  </div>
                  <div style="font-size: 0.9375rem; color: #64748b; font-weight: 600; margin-top: 2px;">{{ selectedProjectsSummary }}</div>
                </div>

                <div v-if="activeWaitlistOffer" class="waitlist-claim-notice is-map-drawer" role="status">
                  <div class="waitlist-claim-icon" aria-hidden="true">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <circle cx="12" cy="12" r="9"></circle>
                      <path d="m8.5 12 2.25 2.25L15.5 9.5"></path>
                    </svg>
                  </div>
                  <div class="waitlist-claim-content">
                    <div class="waitlist-claim-title-row">
                      <strong>ถึงคิวของคุณแล้ว</strong>
                      <span class="waitlist-claim-badge">ได้รับสิทธิ์</span>
                    </div>
                    <span>แผง {{ activeWaitlistOffer.stall_number || activeWaitlistOffer.stall_ident }} ถูกกันไว้ให้คุณ{{ waitlistOfferDeadlineText }}</span>
                  </div>
                </div>

                <!-- 2. Scrollable Middle Body (Stalls List & Queue Info) -->
                <div
                  style="
                    flex: 1 1 auto;
                    overflow-y: auto;
                    overflow-x: hidden;
                    overscroll-behavior: contain;
                    -webkit-overflow-scrolling: touch;
                    padding: 12px 14px;
                    display: flex;
                    flex-direction: column;
                    gap: 10px;
                  "
                >
                  <!-- Per-stall cards: แยกกล่องทีละแผง โดยแต่ละแผงใช้ UI รูปภาพบน + กล่อง 3 คอลัมน์ด้านล่าง -->
                  <div class="sel-stall-list" style="display: flex; flex-direction: column; gap: 10px;">
                    <div v-for="s in selectedStalls" :key="stallSelectionKey(s)" class="sel-stall-card" @click="focusClaimedStallOnMap()" style="cursor: pointer;" title="คลิกเพื่อดูตำแหน่งแผงนี้บนแผนที่">
                      <div class="ssc-head">
                        <div class="ssc-title">
                          <div class="ssc-title-main">
                            <strong>แผง {{ stallDisplayNumber(s) }}</strong>
                            <span :class="s.status === 'available' ? 'badge badge-available' : 'badge badge-reserved'" class="ssc-status-badge">
                              {{ s.status === 'available' ? 'แผงว่าง' : 'ติดจอง' }}
                            </span>
                          </div>
                          <span v-if="s.zone || s.floorName" class="ssc-floor-name" :title="s.zone || s.floorName">{{ s.zone }}<template v-if="s.zone && s.floorName && s.floorName !== s.zone"> • {{ s.floorName }}</template><template v-else-if="!s.zone">{{ s.floorName }}</template></span>
                        </div>
                        <button type="button" class="ssc-remove" @click.stop="removeStall(s)" title="นำแผงนี้ออก" aria-label="ลบแผง">
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                        </button>
                      </div>

                      <!-- Stall Photo Banner (Full width with [ 🖼 3 รูป ] badge) -->
                      <div
                        class="stall-photo-banner-interactive"
                        style="position: relative; border-radius: 10px; overflow: hidden; height: 112px; margin-bottom: 8px; border: 1px solid #e2e8f0; cursor: pointer; background: #0f172a;"
                        @click="openPhotoPreview(s)"
                        title="คลิกเพื่อดูภาพตัวอย่างทั้งหมด"
                      >
                        <img
                          :src="getStallPhotos(s)[0]"
                          alt="ภาพแผง"
                          style="width: 100%; height: 100%; object-fit: cover; display: block; transition: transform 0.25s ease;"
                        />
                        <div style="position: absolute; bottom: 6px; right: 6px; background: rgba(15, 23, 42, 0.82); backdrop-filter: blur(4px); color: #ffffff; padding: 2px 8px; border-radius: 6px; font-size: 0.8125rem; font-weight: 600; display: inline-flex; align-items: center; gap: 4px;">
                          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                          <span>{{ getStallPhotos(s).length > 1 ? `${getStallPhotos(s).length} รูป` : '3 รูป' }}</span>
                        </div>
                      </div>

                      <!-- 3-Box Metrics Grid (ขนาด | ค่าเช่า | แรกเข้า) -->
                      <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 10px 6px; text-align: center;">
                        <div>
                          <div style="font-size: 0.8125rem; color: #64748b; font-weight: 600; margin-bottom: 2px;">ขนาด</div>
                          <div style="font-size: 1.0625rem; font-weight: 800; color: #0f172a;">{{ s.size }} <span style="font-size: 0.8125rem; font-weight: normal; color: #64748b;">ตร.ม.</span></div>
                        </div>
                        <div style="border-left: 1px solid #e2e8f0; border-right: 1px solid #e2e8f0;">
                          <div style="font-size: 0.8125rem; color: #64748b; font-weight: 600; margin-bottom: 2px;">ค่าเช่า</div>
                          <div style="font-size: 1.0625rem; font-weight: 800; color: #10b981;">฿{{ s.monthlyRent.toLocaleString() }}<div style="font-size: 0.75rem; font-weight: 500; color: #64748b;">/เดือน</div></div>
                        </div>
                        <div>
                          <div style="font-size: 0.8125rem; color: #64748b; font-weight: 600; margin-bottom: 2px;">แรกเข้า</div>
                          <div style="font-size: 1.0625rem; font-weight: 800; color: #0f172a;">฿{{ s.entryFee.toLocaleString() }}</div>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Waiting List Queue Status (Clean & Minimal) -->
                  <div v-if="currentSelectionStatus === 'reserved' && currentWaitingInfo">
                    <div style="background: #fffbeb; border: 1px solid #fde68a; border-radius: 10px; padding: 10px 12px;">
                      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px;">
                        <span style="font-size: 0.875rem; font-weight: 800; color: #92400e;">
                          {{ selectedWaitingEntry ? 'สถานะคิวของคุณ' : 'เมื่อลงทะเบียน คุณจะได้' }}
                        </span>
                        <span class="badge" style="background: #d97706; color: #ffffff; font-weight: 800; font-size: 0.8125rem; padding: 2px 8px;">
                          คิวสำรองที่ {{ selectedWaitingEntry?.position || currentWaitingInfo.next_queue_number }}
                        </span>
                      </div>
                      <div style="font-size: 0.8125rem; color: #78350f; display: flex; justify-content: space-between; border-top: 1px dashed #fcd34d; padding-top: 6px; margin-top: 4px;">
                        <span>ผู้ถือสิทธิ์ปัจจุบัน:</span>
                        <strong>{{ currentWaitingInfo.current_booker?.masked_name || 'ผู้จองหลัก' }}</strong>
                      </div>
                    </div>
                  </div>
                  <div v-else-if="currentSelectionStatus === 'reserved'" style="background: #fffbeb; border: 1px solid #fde68a; border-radius: 10px; padding: 10px 12px; font-size: 0.8125rem; color: #92400e;">
                    กำลังตรวจสอบลำดับคิวจริงของแผงนี้
                  </div>
                </div>

                <!-- 3. Drawer Sticky / Pinned Footer CTA Bar (Always Visible at Bottom without Scrolling) -->
                <div
                  style="
                    flex-shrink: 0;
                    padding: 12px 14px 14px;
                    background: #ffffff;
                    border-top: 1px solid #eef2f6;
                    box-shadow: 0 -4px 14px rgba(15, 23, 42, 0.06);
                  "
                >
                  <!-- Total Summary Header if > 1 stall -->
                  <div
                    v-if="selectedStalls.length > 1"
                    style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; font-size: 0.875rem;"
                  >
                    <span style="color: #64748b; font-weight: 600;">รวม {{ selectedStalls.length }} แผง ({{ totalSize }} ตร.ม.)</span>
                    <span style="color: #0f172a; font-weight: 700;">
                      <span style="color: #10b981; font-weight: 800; font-size: 1.0625rem;">฿{{ totalMonthlyRent.toLocaleString() }}</span>
                      <span style="font-size: 0.75rem; color: #64748b;">/ด.</span>
                    </span>
                  </div>

                  <!-- Primary CTA Button (Always instantly clickable without scrolling) -->
                  <button
                    type="button"
                    class="btn btn-primary"
                    :disabled="currentSelectionStatus === 'reserved' && myWaitlistByStall[selectedStalls[0]?.code || '']?.status === 'waiting'"
                    style="
                      width: 100%;
                      padding: 11px 16px;
                      font-size: 0.9375rem;
                      border-radius: 10px;
                      font-weight: 700;
                      box-shadow: 0 4px 12px rgba(37, 99, 235, 0.35);
                      display: flex;
                      align-items: center;
                      justify-content: center;
                      gap: 8px;
                      cursor: pointer;
                    "
                    :style="currentSelectionStatus === 'reserved' && myWaitlistByStall[selectedStalls[0]?.code || '']?.status === 'waiting' ? 'background: #94a3b8; border-color: #94a3b8; box-shadow: none; cursor: default;' : ''"
                    @click="handleNextFromStep1"
                  >
                    <span>
                      {{
                        currentSelectionStatus === 'reserved'
                          ? (myWaitlistByStall[selectedStalls[0]?.code || '']?.status === 'waiting'
                              ? `อยู่ในคิวสำรองลำดับที่ ${myWaitlistByStall[selectedStalls[0]?.code || '']?.position || 1} แล้ว`
                              : 'รอคิวแผงนี้')
                          : 'จองแผงที่เลือก'
                      }}
                    </span>
                    <svg v-if="!(currentSelectionStatus === 'reserved' && myWaitlistByStall[selectedStalls[0]?.code || '']?.status === 'waiting')" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <line x1="5" y1="12" x2="19" y2="12"></line>
                      <polyline points="12 5 19 12 12 19"></polyline>
                    </svg>
                  </button>
                </div>
              </div>
              </transition>

              <!-- Mobile Stall Hint Pill (Only on mobile/tablet when inside project and no stalls chosen yet) -->
              <div v-if="selectedProject && selectedStalls.length === 0 && !isSelectedProjectFull" class="mobile-stall-hint-pill">
                <span class="mshp-dot"></span>
                <span>แตะแผงสีเขียวบนผังเพื่อเลือกจอง</span>
              </div>

              <!-- Mobile Project Full Banner (Only on mobile/tablet if project has 0 available) -->
              <div v-if="selectedProject && isSelectedProjectFull" class="mobile-project-full-banner">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="9"></circle>
                  <path d="m8.5 8.5 7 7m0-7-7 7"></path>
                </svg>
                <span>โครงการนี้มีผู้เช่าครบทุกแผงแล้ว</span>
              </div>

              <!-- Empty state occupies the same right-side drawer until a stall is chosen. -->
              <transition name="slide-fade">
                <aside
                  v-if="selectedProject && selectedStalls.length === 0"
                  class="map-selection-empty-drawer"
                  aria-live="polite"
                  aria-label="กรุณาเลือกแผง"
                >
                  <header class="msed-header">
                    <strong>แผงที่เลือก (0)</strong>
                    <span>{{ selectedProject.name }}</span>
                  </header>
                  <div class="msed-content">
                    <template v-if="isSelectedProjectFull">
                      <span class="msed-icon is-full" aria-hidden="true">
                        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"></circle><path d="m8.5 8.5 7 7m0-7-7 7"></path></svg>
                      </span>
                      <h3>โครงการนี้เต็มแล้ว</h3>
                      <p>ขณะนี้มีผู้เช่าครบทุกแผง จึงไม่สามารถเลือกแผงหรือทำรายการจองได้</p>
                    </template>
                    <template v-else>
                    <span class="msed-icon" aria-hidden="true">
                      <svg width="27" height="27" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"></rect><path d="M3 9h18M9 21V9"></path></svg>
                    </span>
                    <h3>กรุณาเลือกแผงที่สนใจ</h3>
                    <p>คลิกแผงบนผังเพื่อดูรายละเอียด และเพิ่มเข้ารายการจอง</p>
                    <div class="msed-status-list">
                      <div><i class="is-available"></i><span><strong>แผงสีเขียว</strong> จองได้ทันที</span></div>
                      <div><i class="is-reserved"></i><span><strong>แผงสีเหลือง</strong> ลงทะเบียนคิวสำรอง</span></div>
                    </div>
                    </template>
                  </div>
                </aside>
              </transition>
            </template>
          </MarketProjectsOverviewMap>
        </div>


        <!-- (ย้ายออกมานอกผังแผงเดิม เพื่อให้ใช้งานได้ทั้งโหมดแผนที่โครงการและโหมดผังแผง) -->
        <!-- Conflict Alert Modal (เมื่อพยายามเลือกแผงว่างปนกับแผงติดจอง) -->
        <transition name="modal-fade">
          <div v-if="showConflictModal" class="modal-overlay" style="z-index: 1000;" @click.self="cancelConflictModal">
            <div class="modal-card card" style="max-width: 520px; width: 92%; padding: 28px; border-radius: 16px; box-shadow: var(--shadow-popup);">
              <div style="display: flex; align-items: flex-start; gap: 16px; margin-bottom: 18px;">
                <div style="width: 48px; height: 48px; border-radius: 12px; background: #fef3c7; color: #d97706; display: flex; align-items: center; justify-content: center; flex-shrink: 0;">
                  <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
                    <line x1="12" y1="9" x2="12" y2="13"></line>
                    <line x1="12" y1="17" x2="12.01" y2="17"></line>
                  </svg>
                </div>
                <div>
                  <h3 style="font-size: 1.15rem; font-weight: 700; color: #0f172a; margin-bottom: 6px;">
                    ไม่สามารถเลือกแผงว่างร่วมกับแผงติดจองได้
                  </h3>
                  <p style="font-size: 0.88rem; color: #475569; line-height: 1.55; margin-bottom: 0;">
                    <span v-if="pendingConflictingStall?.status === 'reserved'">
                      คุณกำลังคลิกเลือก <strong>แผง {{ stallDisplayNumber(pendingConflictingStall) }} (สถานะ: ติดจอง<span v-if="stallWaitingQueueCache[pendingConflictingStall?.code]"> — คิวสำรองลำดับที่ {{ stallWaitingQueueCache[pendingConflictingStall?.code]?.next_queue_number }}</span>)</strong><span v-if="stallWaitingQueueCache[pendingConflictingStall?.code]?.current_booker"> ต่อจาก <strong>{{ stallWaitingQueueCache[pendingConflictingStall?.code]?.current_booker?.masked_name }}</strong></span> แต่รายการที่คุณเลือกไว้ก่อนหน้าเป็น <strong>กลุ่มแผงว่าง</strong>
                    </span>
                    <span v-else>
                      คุณกำลังคลิกเลือก <strong>แผง {{ stallDisplayNumber(pendingConflictingStall) }} (สถานะ: ว่าง)</strong> แต่รายการที่คุณเลือกไว้ก่อนหน้าเป็น <strong>กลุ่มแผงติดจอง (Waiting List)</strong>
                    </span>
                  </p>
                </div>
              </div>

              <!-- Business Rule Explanation Box -->
              <div style="background: #fffbeb; border: 1px solid #fde68a; border-radius: 10px; padding: 14px 16px; font-size: 0.84rem; color: #92400e; margin-bottom: 22px; line-height: 1.5;">
                <div style="font-weight: 700; margin-bottom: 4px; display: flex; align-items: center; gap: 6px;">
                  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
                  <span>เหตุผลทางขั้นตอนการเช่า (Business Rule):</span>
                </div>
                แผงที่ <strong>ว่าง</strong> สามารถนัดหมายเข้าดูและทำสัญญาเช่าได้ทันที ส่วนแผงที่ <strong>ติดจอง</strong> จะเป็นการลงทะเบียนคิวสำรอง (Waiting List) หากรวมในรายการจองเดียวกันจะทำให้สถานะการจองค้าง และไม่สามารถแยกยื่นคำร้องขอเช่าแผงที่พร้อมได้
              </div>

              <!-- Action Buttons -->
              <div style="display: flex; gap: 10px; justify-content: flex-end; flex-wrap: wrap;">
                <button
                  type="button"
                  class="btn btn-outline"
                  style="border-color: #cbd5e1; color: #475569; font-weight: 600;"
                  @click="cancelConflictModal"
                >
                  ยกเลิก (คงรายการเดิมไว้)
                </button>
                <button
                  type="button"
                  class="btn btn-primary"
                  style="background: #d97706; border-color: #d97706; color: #ffffff; font-weight: 700;"
                  @click="resolveConflictAndSelect"
                >
                  ล้างรายการเดิม แล้วเลือกแผง {{ stallDisplayNumber(pendingConflictingStall) }}
                </button>
              </div>
            </div>
          </div>
        </transition>

        <!-- Step 1 Bottom Action Row (Matching Next Button) -->
        <div style="display: flex; justify-content: flex-end; align-items: center; margin-top: 16px;">
          <button
            type="button"
            class="btn"
            :class="isMarketFull ? 'btn-outline' : 'btn-primary'"
            :disabled="isMarketFull"
            :style="isMarketFull ? 'padding: 12px 28px; font-size: 0.92rem; font-weight: 600; border-radius: 10px; opacity: 0.65; cursor: not-allowed; border-color: #cbd5e1; background: #f1f5f9; color: #64748b;' : 'padding: 12px 32px; font-size: 1rem; font-weight: 700; border-radius: 10px; box-shadow: 0 4px 14px rgba(67, 97, 238, 0.35); display: inline-flex; align-items: center; gap: 8px;'"
            @click="handleNextFromStep1"
          >
            <span>{{ isMarketFull ? 'ตลาดเต็มแล้ว (ไม่สามารถจองได้)' : 'ถัดไป' }}</span>
            <span v-if="!isMarketFull">&rarr;</span>
          </button>
        </div>
      </div>

      <!-- ============================================================ -->
      <!-- STEP 2: กรอกข้อมูลเพื่อเยี่ยมชมแผง (Enhanced with AI Smart ID)-->
      <!-- ============================================================ -->
      <div v-if="currentStep === 2">
        <div style="margin-bottom: 24px;">
          <h2 style="font-size: 1.25rem; font-weight: 800; color: #0f172a; margin-bottom: 4px;">
            กรอกข้อมูลเพื่อเยี่ยมชมแผง
          </h2>
          <div style="height: 3px; width: 48px; background: var(--primary); border-radius: var(--radius-full);"></div>
        </div>



        <!-- Waiting List Context Alert Banner for Step 2 -->
        <div v-if="currentSelectionStatus === 'reserved'" style="margin-bottom: 22px; background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%); border: 1.5px solid #fcd34d; border-radius: 12px; padding: 14px 18px; color: #92400e; font-size: 0.85rem; display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 10px;">
          <div style="display: flex; align-items: center; gap: 10px;">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" style="color: #d97706; flex-shrink: 0;"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            <div>
              <strong>คำขอนี้เป็นการลงทะเบียนรับสิทธิ์คิวสำรอง<span v-if="currentWaitingInfo"> (คิวสำรองที่ {{ currentWaitingInfo.next_queue_number }})</span></strong>
              <div style="font-size: 0.78rem; color: #b45309; margin-top: 2px;">
                <template v-if="currentWaitingInfo?.current_booker">ต่อจากผู้จองสิทธิ์ก่อนหน้า: <strong>{{ currentWaitingInfo.current_booker.masked_name }}</strong> (ปิดข้อมูลบางส่วนตามหลัก PDPA)</template>
                <template v-else>กำลังตรวจสอบข้อมูลคิวจากระบบ</template>
              </div>
            </div>
          </div>
          <span class="badge" style="background: #d97706; color: #ffffff; padding: 5px 12px; font-weight: 700; border-radius: 20px;">
            {{ currentWaitingInfo ? `คิวสำรองที่ ${currentWaitingInfo.next_queue_number}` : 'กำลังตรวจสอบคิว' }}
          </span>
        </div>


        <form @submit.prevent="goToStep(3)">
          <!-- Section 1: ข้อมูลผู้เช่า -->
          <div class="card" style="margin-bottom: 24px; padding: 28px;">
            <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #eef2f6; padding-bottom: 16px; margin-bottom: 22px;">
              <div style="display: flex; align-items: center; gap: 12px;">
                <div style="width: 40px; height: 40px; border-radius: 12px; background: #eff6ff; color: #2563eb; display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0;">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                    <circle cx="12" cy="7" r="4"></circle>
                  </svg>
                </div>
                <h3 style="font-size: 1.05rem; font-weight: 700; color: #0f172a; margin: 0;">
                  ข้อมูลผู้เช่า
                </h3>
              </div>
              <span v-if="idVerified" style="font-size: 0.78rem; color: #059669; font-weight: 600; display: inline-flex; align-items: center; gap: 4px;">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
                  <polyline points="9 12 11 14 15 10"></polyline>
                </svg>
                <span>ยืนยันข้อมูลผ่านบัตรประจำตัวประชาชนแล้ว</span>
              </span>
            </div>

            <!-- Row 1: ชื่อ นามสกุล และ เลขบัตรประชาชน 13 หลัก (ดึงจากโปรไฟล์ที่เข้าสู่ระบบ - Disabled) -->
            <div class="form-row-3" style="margin-bottom: 16px;">
              <div class="form-group" style="margin-bottom: 0;">
                <label class="form-label" for="form_first_name">ชื่อ</label>
                <input
                  id="form_first_name"
                  v-model="bookingForm.firstName"
                  type="text"
                  disabled
                  placeholder="เช่น สมชาย"
                  class="form-input"
                />
              </div>

              <div class="form-group" style="margin-bottom: 0;">
                <label class="form-label" for="form_last_name">นามสกุล</label>
                <input
                  id="form_last_name"
                  v-model="bookingForm.lastName"
                  type="text"
                  disabled
                  placeholder="เช่น ใจดี"
                  class="form-input"
                />
              </div>

              <div class="form-group" style="margin-bottom: 0;">
                <label class="form-label" for="form_id_card">เลขประจำตัวประชาชน 13 หลัก</label>
                <input
                  id="form_id_card"
                  v-model="bookingForm.idCard"
                  type="text"
                  disabled
                  maxlength="17"
                  placeholder="1-XXXX-XXXXX-XX-X"
                  class="form-input"
                />
              </div>
            </div>

            <!-- Row 2: เบอร์ติดต่อ อีเมล และ สินค้าที่จะขาย -->
            <div class="form-row-3">
              <div class="form-group" style="margin-bottom: 0;">
                <label class="form-label" for="form_phone">เบอร์ติดต่อ</label>
                <input
                  id="form_phone"
                  v-model="bookingForm.phone"
                  type="tel"
                  disabled
                  placeholder="087-567-4497"
                  class="form-input"
                />
              </div>

              <div class="form-group" style="margin-bottom: 0;">
                <label class="form-label" for="form_email">อีเมล</label>
                <input
                  id="form_email"
                  v-model="bookingForm.email"
                  type="email"
                  disabled
                  placeholder="tenant@emarket.bma.go.th"
                  class="form-input"
                />
              </div>

              <!-- สินค้าที่ต้องการขาย * (เปิดให้เลือกผ่าน Custom UI Dropdown) -->
              <div class="form-group" :class="{ 'has-form-error': !!formErrors.business }" style="margin-bottom: 0;">
                <label class="form-label" for="form_business">
                  สินค้าที่ต้องการขาย <span style="color: #ef4444; font-weight: 700;">*</span>
                </label>
                <AppSelect
                  id="form_business"
                  v-model="bookingForm.business"
                  :options="businessOptions"
                  :error="!!formErrors.business"
                  placeholder="เลือกสินค้าที่ต้องการขาย"
                />
                <div v-if="formErrors.business" class="form-field-error-text">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="12"></line>
                    <line x1="12" y1="16" x2="12.01" y2="16"></line>
                  </svg>
                  <span>{{ formErrors.business }}</span>
                </div>
              </div>
            </div>

            <!-- หมายเหตุหรือข้อสอบถามเพิ่มเติม (เปิดให้กรอก) -->
            <div class="form-group" style="margin-top: 16px; margin-bottom: 0;">
              <label class="form-label" for="form_notes">หมายเหตุหรือข้อสอบถามเพิ่มเติม</label>
              <input
                id="form_notes"
                v-model="bookingForm.notes"
                type="text"
                placeholder="ระบุข้อสอบถาม เช่น สนใจทำเลใกล้ทางเข้า หรือต้องการไฟฟ้ากำลังสูง"
                class="form-input"
              />
            </div>
          </div>


          <!-- Section 2: ข้อมูลการเยี่ยมชมแผง (Date Cards Strip & Time Slots) -->
          <div class="card" style="margin-bottom: 28px; padding: 28px;">
            <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #eef2f6; padding-bottom: 16px; margin-bottom: 22px; flex-wrap: wrap; gap: 12px;">
              <div style="display: flex; align-items: center; gap: 12px;">
                <div style="width: 40px; height: 40px; border-radius: 12px; background: #eff6ff; color: #2563eb; display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0;">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                    <line x1="16" y1="2" x2="16" y2="6"></line>
                    <line x1="8" y1="2" x2="8" y2="6"></line>
                    <line x1="3" y1="10" x2="21" y2="10"></line>
                  </svg>
                </div>
                <div>
                  <h3 style="font-size: 1.05rem; font-weight: 700; color: #0f172a; margin: 0 0 2px 0;">
                    ข้อมูลการเยี่ยมชมแผง
                  </h3>
                  <span style="font-size: 0.78rem; color: #64748b;">
                    * เจ้าหน้าที่สำนักงานตลาดจะพาเดินชมสถานที่จริงตามรอบเวลา
                  </span>
                </div>
              </div>
              <span style="font-size: 0.74rem; background: #eff6ff; color: #1e40af; border: 1px solid #bfdbfe; border-radius: 6px; padding: 4px 10px; font-weight: 600; display: inline-flex; align-items: center; gap: 5px;">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
                  <polyline points="9 12 11 14 15 10"></polyline>
                </svg>
                <span>กำหนดรอบโดยสำนักงานตลาด</span>
              </span>
            </div>

            <!-- Step 1: แถบการ์ดเลือกวัน (Date Cards Strip) -->
            <div :class="{ 'has-form-error': !!formErrors.visitDate }" style="margin-bottom: 22px;">
              <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; flex-wrap: wrap; gap: 8px;">
                <label class="form-label" style="margin-bottom: 0; font-weight: 700; color: #1e293b; display: inline-flex; align-items: center; gap: 8px;">
                  <span style="background: #2563eb; color: #ffffff; width: 22px; height: 22px; border-radius: 50%; display: inline-flex; align-items: center; justify-content: center; font-size: 0.75rem; font-weight: 800;">1</span>
                  <span style="font-size: 0.95rem;">เลือกวันที่เปิดให้เข้าชมแผง <span style="color: #ef4444; font-weight: 700;">*</span></span>
                </label>
                <NuxtLink
                  :to="`/stalls/visit-schedule?market=${route.query.market || '1'}`"
                  target="_blank"
                  style="font-size: 0.76rem; color: #2563eb; background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 6px; padding: 4px 10px; display: inline-flex; align-items: center; gap: 5px; text-decoration: none; font-weight: 600;"
                  title="เปิดหน้าจัดการรอบเข้าชมแผงสำหรับเจ้าหน้าที่"
                >
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="3"></circle>
                    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                  </svg>
                  <span>จัดการรอบเวลา (เจ้าหน้าที่)</span>
                </NuxtLink>
              </div>

              <!-- Strip of Date Cards -->
              <div class="date-cards-strip" :class="{ 'date-strip-error': !!formErrors.visitDate }">
                <div
                  v-for="d in availableDates"
                  :key="d.date"
                  class="date-card"
                  :class="{
                    selected: bookingForm.visitDate === d.date,
                    'is-full-date': d.isAllFull,
                    'is-closed-date': d.isClosed
                  }"
                  :title="bookingForm.visitDate === d.date ? 'คลิกอีกครั้งเพื่อยกเลิกวันที่เลือก' : `คลิกเพื่อเลือกวันที่ ${d.dayName} ${d.dayNum} ${d.monthShort}`"
                  @click="selectVisitDate(d)"
                >
                  <div class="date-card-header">
                    <span class="date-day-name">{{ d.dayName }}</span>
                    <span v-if="d.badge" class="date-badge" :class="d.badgeType">
                      {{ d.badge }}
                    </span>
                  </div>
                  <div class="date-number-wrap">
                    <span class="date-day-num">{{ d.dayNum }}</span>
                    <span class="date-month">{{ d.monthShort }}</span>
                  </div>
                  <div class="date-card-footer">
                    <span class="date-status" :class="d.statusClass">
                      <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 3px;">
                        <circle cx="12" cy="12" r="10"></circle>
                        <polyline points="12 6 12 12 16 14"></polyline>
                      </svg>
                      {{ d.slotsText }}
                    </span>
                    <span v-if="bookingForm.visitDate === d.date" style="color: #2563eb; display: inline-flex;">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                    </span>
                  </div>
                </div>
              </div>

              <!-- Inline error message for date -->
              <div v-if="formErrors.visitDate" class="form-field-error-text" style="margin-top: 8px;">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="8" x2="12" y2="12"></line>
                  <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <span>{{ formErrors.visitDate }}</span>
              </div>

              <!-- Market Advice banner under date cards (Dynamic from active date description or market default) -->
              <div style="font-size: 0.76rem; color: #1e40af; background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 8px; padding: 8px 12px; margin-top: 10px; display: flex; align-items: center; gap: 8px;">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink: 0; color: #2563eb;">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="16" x2="12" y2="12"></line>
                  <line x1="12" y1="8" x2="12.01" y2="8"></line>
                </svg>
                <span>{{ currentSelectedDateAdvice || 'คำแนะนำ: ตลาดเปิดบริการเต็มรูปแบบในวันเสาร์และอาทิตย์ (09:00 - 18:00 น.) ชมบรรยากาศและทราฟฟิกลูกค้าจริงได้ดีที่สุด' }}</span>
              </div>
            </div>

            <!-- Step 2: แถบการ์ดรอบเวลา (Time Slots Strip เชื่อมโยงกับวันที่เลือก) -->
            <div :class="{ 'has-form-error': !!formErrors.visitSlot }">
              <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; flex-wrap: wrap; gap: 8px;">
                <label class="form-label" style="margin-bottom: 0; font-weight: 700; color: #1e293b; display: inline-flex; align-items: center; gap: 6px;">
                  <span style="background: #2563eb; color: #ffffff; width: 20px; height: 20px; border-radius: 50%; display: inline-flex; align-items: center; justify-content: center; font-size: 0.72rem; font-weight: 800;">2</span>
                  <span>เลือกรอบเวลาที่สะดวกเข้าชม (พาชมโดยเจ้าหน้าที่) <span style="color: #ef4444; font-weight: 700;">*</span></span>
                </label>
                <span style="font-size: 0.76rem; color: #64748b;">
                  รอบละ 60 นาที &bull; เฉพาะรอบที่ยังว่าง
                </span>
              </div>

              <!-- Notice when selected date has no available slots -->
              <div v-if="isCurrentDateAllFull" style="background: #fef2f2; border: 1.5px solid #fecaca; border-radius: 10px; padding: 12px 16px; margin-bottom: 12px; display: flex; align-items: center; gap: 10px; color: #991b1b; font-size: 0.85rem;">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink: 0; color: #dc2626;">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="15" y1="9" x2="9" y2="15"></line>
                  <line x1="9" y1="9" x2="15" y2="15"></line>
                </svg>
                <div>
                  <strong>รอบเวลาของวันที่เลือกเต็มทุกรอบแล้ว</strong>
                  <div style="font-size: 0.78rem; color: #b91c1c; margin-top: 2px;">
                    กรุณาคลิกเลือกวันที่อื่นด้านบนที่มีรอบว่าง เพื่อทำการนัดหมายเข้าชม
                  </div>
                </div>
              </div>

              <!-- Notice when user has not yet selected a date -->
              <div v-else-if="!bookingForm.visitDate" style="background: #f8fafc; border: 1.5px dashed #cbd5e1; border-radius: 10px; padding: 12px 16px; margin-bottom: 12px; display: flex; align-items: center; gap: 10px; color: #64748b; font-size: 0.84rem;">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: #94a3b8; flex-shrink: 0;">
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                </svg>
                <span>กรุณาคลิกเลือกวันที่ในข้อ 1 ด้านบนก่อน เพื่อดูและเลือกรอบเวลาที่เปิดให้บริการ</span>
              </div>

              <div class="time-slot-cards-grid" :class="{ 'time-slots-error': !!formErrors.visitSlot }">
                <div
                  v-for="s in currentTimeSlots"
                  :key="s.time"
                  class="time-slot-card"
                  :class="{
                    selected: bookingForm.visitSlot === s.time && !s.isFull && !s.hasTenantConflict,
                    'is-full-slot': s.isFull || s.hasTenantConflict
                  }"
                  :title="bookingForm.visitSlot === s.time ? `คลิกอีกครั้งเพื่อยกเลิกรอบ ${s.time}` : (s.enabled === false ? `${s.time} งดให้บริการรอบเวลานี้` : (s.hasTenantConflict ? `${s.time} คุณมีนัดหมายอื่นในเวลานี้แล้ว` : (s.isFull ? `${s.time} เต็มแล้ว (${s.bookedCount || s.capacity}/${s.capacity} คิว) ห้ามเลือก` : `คลิกเพื่อเลือกรอบ ${s.time}`)))"
                  @click="selectTimeSlot(s)"
                >
                  <div class="slot-card-top">
                    <span class="slot-card-icon" :style="{ color: (s.isFull || s.hasTenantConflict) ? '#94a3b8' : 'var(--primary)' }">
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="10"></circle>
                        <polyline points="12 6 12 12 16 14"></polyline>
                      </svg>
                    </span>
                    <strong class="slot-card-time" :style="{ color: (s.isFull || s.hasTenantConflict) ? '#64748b' : '#0f172a' }">{{ s.time }}</strong>
                    <span v-if="s.enabled === false" style="margin-left: auto; font-size: 0.75rem; color: #64748b; font-weight: 700; background: #f1f5f9; padding: 2px 7px; border-radius: 4px; display: inline-flex; align-items: center; gap: 3px;">
                      งดรอบนี้
                    </span>
                    <span v-else-if="s.hasTenantConflict" style="margin-left: auto; font-size: 0.75rem; color: #d97706; font-weight: 700; background: #fef3c7; border: 1px solid #fde68a; padding: 2px 7px; border-radius: 4px; display: inline-flex; align-items: center; gap: 3px;">
                      มีนัดหมายอื่นแล้ว
                    </span>
                    <span v-else-if="s.isFull" style="margin-left: auto; font-size: 0.75rem; color: #dc2626; font-weight: 700; background: #fee2e2; padding: 2px 7px; border-radius: 4px; display: inline-flex; align-items: center; gap: 3px;">
                      <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                      </svg>
                      เต็มแล้ว
                    </span>
                    <span v-else-if="(s.bookedCount || 0) > 0" class="slot-booked-summary" style="margin-left: auto; font-size: 0.75rem; color: #2563eb; font-weight: 600; background: #eff6ff; border: 1px solid #bfdbfe; padding: 2px 7px; border-radius: 4px; display: inline-flex; align-items: center; gap: 3px;">
                      จองแล้ว {{ s.bookedCount }}/{{ s.capacity }} คิว
                    </span>
                  </div>
                  <div class="slot-card-capacity" :class="s.capacityClass">
                    {{ s.capacityText }}
                  </div>
                  <div class="slot-card-note">{{ s.note }}</div>
                </div>
              </div>

              <!-- Inline error message for slot -->
              <div v-if="formErrors.visitSlot" class="form-field-error-text" style="margin-top: 8px;">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="8" x2="12" y2="12"></line>
                  <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <span>{{ formErrors.visitSlot }}</span>
              </div>
            </div>
          </div>

          <!-- Step Navigation Buttons -->
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <button type="button" class="btn btn-outline" @click="goToStep(1)">
              <span>ย้อนกลับ</span>
            </button>
            <button type="submit" class="btn btn-primary" style="padding: 10px 24px; font-weight: 700;">
              <span>ถัดไป</span>
              <span>&rarr;</span>
            </button>
          </div>
        </form>
      </div>

      <!-- ============================================================ -->
      <!-- STEP 3: ตรวจสอบและยืนยัน (Enhanced with BA/QA/UX)            -->
      <!-- ============================================================ -->
      <div v-else-if="currentStep === 3" class="step3-review">
        <div style="margin-bottom: 24px;">
          <h2 style="font-size: 1.25rem; font-weight: 800; color: #0f172a; margin-bottom: 4px;">
            ตรวจสอบและยืนยัน
          </h2>
          <div style="height: 3px; width: 48px; background: var(--primary); border-radius: var(--radius-full);"></div>
        </div>

        <!-- Confirmation GIS Map: แสดงแผนที่จริงพร้อมล็อคแผงที่จอง -->
        <div class="card step3-map-card" style="padding: 0; position: relative; overflow: hidden; background: #ffffff; border: 1px solid #e2e8f0; margin-bottom: 24px; border-radius: 16px;">
          <div class="step3-map-heading">
            <div class="step3-map-heading-copy">
              <span class="step3-map-heading-icon" aria-hidden="true"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 5-8 12-8 12S4 15 4 10a8 8 0 1 1 16 0z"></path><circle cx="12" cy="10" r="2.5"></circle></svg></span>
              <div><h3>ตำแหน่งแผงค้า</h3><p>{{ currentMarket?.name || currentMarketName || 'ตลาดนัดจตุจักร' }} · {{ selectedProject?.name || 'โครงการ' }}</p></div>
            </div>
            <div class="step3-map-heading-legend"><span><i class="selected"></i>แผงที่คุณเลือก ({{ selectedStalls.length }})</span><span><i class="available"></i>แผงว่าง</span><span><i class="reserved"></i>ติดจอง</span></div>
          </div>
          <MarketProjectsOverviewMap
            containerId="step3-confirmation-maplibre"
            :marketId="currentMarketId"
            :marketFull="isMarketFull"
            :selectedProjectId="selectedProject?.id"
            :selectedStallCodes="selectedStalls.map(s => s.code)"
            :selectedStallKeys="selectedStalls.map(stallSelectionKey)"
            :stallStatusOverrides="stallStatusOverrides"
            v-model:mapStyle="mapStyle"
            :previewMode="true"
          />
        </div>

        <!-- Card 1: แผงเช่าที่คุณเลือก (List of selected stalls - แยกการ์ดแต่ละแผงพร้อมรูปภาพในตัว) -->
        <div class="card step3-stalls-card" :class="{ 'has-multiple-stalls': selectedStalls.length > 1 }" style="margin-bottom: 24px; padding: 24px;">
          <!-- Card Header -->
          <div style="display: flex; align-items: center; gap: 12px; border-bottom: 1px solid #eef2f6; padding-bottom: 16px; margin-bottom: 20px;">
            <div style="width: 38px; height: 38px; border-radius: 10px; background: #eff6ff; color: #2563eb; display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
                <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
                <line x1="12" y1="22.08" x2="12" y2="12"></line>
              </svg>
            </div>
            <h3 style="font-size: 1.05rem; font-weight: 700; color: #0f172a; margin: 0;">
              แผงเช่าที่คุณเลือก
            </h3>
          </div>

          <!-- Stalls Cards Loop (แยกการ์ดแต่ละแผงพร้อมรูปภาพในตัวตามรูป) -->
          <div class="step3-stalls-scroll">
          <div
            v-for="(stallItem, idx) in selectedStalls"
            :key="stallSelectionKey(stallItem)"
            style="background: #ffffff; border: 1px solid #e2e8f0; border-radius: 14px; padding: 22px 24px; margin-bottom: 20px;"
          >
            <!-- Sub-card Header -->
            <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #f1f5f9; padding-bottom: 14px; margin-bottom: 18px;">
              <span style="background: #2563eb; color: #ffffff; font-weight: 700; font-size: 0.85rem; padding: 4px 14px; border-radius: 6px;">
                แผง {{ stallDisplayNumber(stallItem) }}
              </span>
              <span
                v-if="stallItem.status === 'available'"
                style="background: #ecfdf5; color: #10b981; font-weight: 600; font-size: 0.82rem; padding: 4px 14px; border-radius: 20px;"
              >
                ว่าง
              </span>
              <span
                v-else
                style="background: #fffbeb; color: #d97706; font-weight: 600; font-size: 0.82rem; padding: 4px 14px; border-radius: 20px;"
              >
                ติดจอง
              </span>
            </div>

            <!-- 4-Column Data Grid -->
            <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px 20px; font-size: 0.88rem; margin-bottom: 18px;">
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 4px;">โซน</div>
                <div style="color: #1e293b; font-size: 0.95rem; font-weight: 600;">{{ stallItem.zone || selectedProject?.name }}</div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 4px;">หมายเลขแผง</div>
                <div style="color: #1e293b; font-size: 0.95rem; font-weight: 700; display: inline-flex; align-items: center; gap: 6px;">
                  <span>{{ stallDisplayNumber(stallItem) }}</span>
                  <span v-if="stallItem.floorName" style="font-size: 0.72rem; padding: 1px 7px; border-radius: 4px; background: #eff6ff; color: #2563eb; font-weight: 600; border: 1px solid #bfdbfe;">
                    {{ stallItem.floorName }}
                  </span>
                </div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 4px;">ประเภทพื้นที่</div>
                <div style="color: #1e293b; font-size: 0.95rem; font-weight: 600;">{{ stallItem.areaType || 'ทั่วไป' }}</div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 4px;">ขนาดพื้นที่(ตร.ม.)</div>
                <div style="color: #1e293b; font-size: 0.95rem; font-weight: 600;">{{ stallItem.size }}</div>
              </div>

              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 4px;">ค่าธรรมเนียมสิทธิแรกเข้า (บาท)</div>
                <div style="color: #1e293b; font-size: 0.95rem; font-weight: 600;">{{ (stallItem.entryFee || 75000).toLocaleString() }}</div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 4px;">ค่าเช่าต่อเดือน (บาท)</div>
                <div style="color: #1e293b; font-size: 0.95rem; font-weight: 600;">{{ (stallItem.monthlyRent || 1800).toLocaleString() }}</div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 4px;">ค่าส่วนกลาง (บาท)</div>
                <div style="color: #1e293b; font-size: 0.95rem; font-weight: 600;">{{ (stallItem.commonFee || 500).toLocaleString() }}</div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 4px;">ค่าประกันล่วงหน้า (บาท)</div>
                <div style="color: #1e293b; font-size: 0.95rem; font-weight: 600;">{{ (stallItem.depositFee || stallItem.monthlyRent || 1800).toLocaleString() }}</div>
              </div>
            </div>

            <!-- Photos of THIS specific stall -->
            <div class="step3-photos">
              <div style="color: #64748b; font-size: 0.82rem; margin-bottom: 10px;">รูปภาพ</div>
              <div style="display: flex; gap: 14px; flex-wrap: wrap;">
                <div
                  v-for="(photoUrl, pIdx) in getStallPhotos(stallItem)"
                  :key="pIdx"
                  class="summary-stall-photo-card"
                  @click="openPhotoPreview(stallItem, pIdx)"
                  :title="`คลิกเพื่อดูรูปขนาดใหญ่ (รูปที่ ${pIdx + 1})`"
                >
                  <img
                    :src="photoUrl"
                    :alt="`รูปแผง ${stallDisplayNumber(stallItem)} ภาพที่ ${pIdx + 1}`"
                    style="width: 100%; height: 100%; object-fit: cover;"
                  />
                  <div class="photo-card-hover-overlay">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                      <circle cx="11" cy="11" r="8"></circle>
                      <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                      <line x1="11" y1="8" x2="11" y2="14"></line>
                      <line x1="8" y1="11" x2="14" y2="11"></line>
                    </svg>
                    <span>ดูรูปใหญ่</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Waiting List Details with PDPA Masking if reserved -->
            <div v-if="stallItem.status === 'reserved'" style="margin-top: 14px; background: #fffbeb; border: 1.5px solid #fcd34d; border-radius: 10px; padding: 12px 16px; font-size: 0.82rem; color: #92400e; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 8px;">
              <div style="display: flex; align-items: center; gap: 8px;">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" style="color: #d97706; flex-shrink: 0;"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                <span>
                  <strong>สถานะคิวสำรอง:</strong>
                  <template v-if="stallWaitingQueueCache[stallItem.code]"> คุณจะได้รับสิทธิ์ <strong>คิวสำรองที่ {{ stallWaitingQueueCache[stallItem.code].next_queue_number }}</strong><span v-if="stallWaitingQueueCache[stallItem.code].current_booker"> (ต่อจาก: <strong>{{ stallWaitingQueueCache[stallItem.code].current_booker?.masked_name }}</strong>)</span></template>
                  <template v-else> กำลังตรวจสอบลำดับคิวจริง</template>
                </span>
              </div>
              <span style="font-size: 0.72rem; background: #ecfdf5; border: 1px solid #a7f3d0; color: #047857; font-weight: 700; padding: 3px 10px; border-radius: 12px; display: inline-flex; align-items: center; gap: 4px;">
                <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
                <span>PDPA Masked</span>
              </span>
            </div>
          </div>
          </div>
        </div>

        <!-- 2-Column Review Section (ข้อมูลผู้จอง & การเยี่ยมชมแผง) -->
        <div class="step3-review-grid" style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 24px;">
          <!-- Card 2: ข้อมูลผู้จอง -->
          <div class="card" style="padding: 24px;">
            <div style="display: flex; align-items: center; gap: 12px; border-bottom: 1px solid #eef2f6; padding-bottom: 14px; margin-bottom: 18px;">
              <div style="width: 38px; height: 38px; border-radius: 10px; background: #eff6ff; color: #2563eb; display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
              </div>
              <span style="font-size: 1.05rem; font-weight: 700; color: #0f172a;">ข้อมูลผู้จอง</span>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px; font-size: 0.88rem;">
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 2px;">ชื่อ - สกุล</div>
                <div style="font-weight: 700; color: #0f172a;">{{ bookingForm.firstName }} {{ bookingForm.lastName }}</div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 2px;">เลขบัตรประชาชน</div>
                <div style="font-weight: 700; color: #0f172a;">{{ bookingForm.idCard || '1-1004-98721-34-5' }}</div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 2px;">เบอร์ติดต่อ</div>
                <div style="font-weight: 700; color: #0f172a;">{{ bookingForm.phone }}</div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 2px;">สินค้าที่จะจำหน่าย</div>
                <div style="font-weight: 700; color: #0f172a;">{{ bookingForm.business }}</div>
              </div>
            </div>
          </div>

          <!-- Card 3: การเยี่ยมชมแผง -->
          <div class="card" style="padding: 24px;">
            <div style="display: flex; align-items: center; gap: 12px; border-bottom: 1px solid #eef2f6; padding-bottom: 14px; margin-bottom: 18px;">
              <div style="width: 38px; height: 38px; border-radius: 10px; background: #eff6ff; color: #2563eb; display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0;">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                  <line x1="16" y1="2" x2="16" y2="6"></line>
                  <line x1="8" y1="2" x2="8" y2="6"></line>
                  <line x1="3" y1="10" x2="21" y2="10"></line>
                </svg>
              </div>
              <span style="font-size: 1.05rem; font-weight: 700; color: #0f172a;">การเยี่ยมชมแผง</span>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px; font-size: 0.88rem;">
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 2px;">วันนัดเยี่ยมชมแผง</div>
                <div style="font-weight: 700; color: #0f172a;">{{ formatThaiDate(bookingForm.visitDate) }}</div>
              </div>
              <div>
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 2px;">ช่วงเวลาที่นัดเยี่ยมชมแผง</div>
                <div style="font-weight: 700; color: #0f172a;">{{ bookingForm.visitSlot }}</div>
              </div>
              <div style="grid-column: 1 / -1;">
                <div style="color: #64748b; font-size: 0.8rem; margin-bottom: 2px;">หมายเหตุ</div>
                <div style="color: #475569;">{{ bookingForm.notes || '-' }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- WOW FEATURE 3: Transparent Financial Fee Schedule -->
        <div class="card step3-finance-card" style="margin-bottom: 24px; padding: 24px; border: 1.5px solid #dbeafe; background: #ffffff;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
            <div style="display: flex; align-items: center; gap: 10px;">
              <span class="step3-finance-icon" aria-hidden="true"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="1" y="4" width="22" height="16" rx="2" ry="2"></rect>
                <line x1="10" y1="2" x2="10" y2="22"></line>
              </svg></span>
              <span style="font-size: 1rem; font-weight: 700; color: #0f172a;">ตารางประมาณการค่าใช้จ่ายแบบโปร่งใส (ประมาณการวันทำสัญญา)</span>
            </div>
            <span class="badge" style="background: #ecfdf5; color: #059669; font-weight: 600;">ไม่มีค่าใช้จ่ายแอบแฝง</span>
          </div>

          <div style="display: flex; flex-direction: column; gap: 10px; font-size: 0.88rem;">
            <div style="display: flex; justify-content: space-between; padding-bottom: 8px; border-bottom: 1px dashed #e2e8f0;">
              <span style="color: #475569;">1. ค่าธรรมเนียมสิทธิแรกเข้า (ชำระครั้งแรกครั้งเดียว)</span>
              <strong style="color: #0f172a;">฿{{ totalEntryFee.toLocaleString() }}</strong>
            </div>
            <div style="display: flex; justify-content: space-between; padding-bottom: 8px; border-bottom: 1px dashed #e2e8f0;">
              <span style="color: #475569;">
                2. ค่าเช่าแผงล่วงหน้า (เดือนแรก)
              </span>
              <strong style="color: #0f172a;">
                ฿{{ totalMonthlyRent.toLocaleString() }}
              </strong>
            </div>
            <div style="display: flex; justify-content: space-between; padding-bottom: 8px; border-bottom: 1px dashed #e2e8f0;">
              <div style="display: flex; align-items: center; gap: 8px;">
                <span style="color: #475569;">3. เงินประกันความเสียหายการเช่า</span>
                <span class="badge" style="background: #eff6ff; color: #2563eb; font-size: 0.72rem; padding: 1px 6px;">* คืนเงิน 100% เมื่อสิ้นสุดสัญญา</span>
              </div>
              <strong style="color: #0f172a;">฿{{ totalDepositFee.toLocaleString() }}</strong>
            </div>
            <div style="display: flex; justify-content: space-between; padding-bottom: 8px; border-bottom: 1px dashed #e2e8f0;">
              <span style="color: #475569;">
                4. ค่าบริการรักษาความสะอาดและส่วนกลาง (เดือนแรก)
              </span>
              <strong style="color: #0f172a;">
                ฿{{ totalCommonFee.toLocaleString() }}
              </strong>
            </div>
            <div style="display: flex; justify-content: space-between; padding-top: 6px; font-size: 1.05rem;">
              <strong style="color: #1e3a8a;">ยอดรวมที่ต้องเตรียมชำระในวันทำสัญญา (โดยประมาณ):</strong>
              <strong style="color: #2563eb; font-size: 1.25rem;">฿{{ contractTotal.toLocaleString() }}</strong>
            </div>
          </div>
        </div>

        <!-- 100% Refundable Security Deposit Guarantee Card (BA) -->
        <div class="step3-guarantee-card" style="background: #f0fdf4; border: 1.5px solid #86efac; border-radius: 12px; padding: 14px 20px; margin-bottom: 24px; display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px;">
          <div style="display: flex; align-items: center; gap: 12px;">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: #15803d; flex-shrink: 0;">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <polyline points="10 9 9 9 8 9"></polyline>
            </svg>
            <div>
              <strong style="color: #166534; font-size: 0.95rem;">การรับประกันเงินประกันการเช่า 100% (Refundable Deposit Guarantee)</strong>
              <div style="font-size: 0.8rem; color: #15803d;">
                เงินประกันสัญญา ฿{{ totalDepositFee.toLocaleString() }} ได้รับคืนเต็มจำนวน 100% ภายใน 15 วันทำการหลังส่งมอบพื้นที่คืนตามระเบียบ กทม.
              </div>
            </div>
          </div>
          <span class="badge" style="background: #15803d; color: #ffffff; padding: 4px 10px; font-weight: 700; display: inline-flex; align-items: center; gap: 4px;">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
            <span>คุ้มครองตามกฎหมาย กทม.</span>
          </span>
        </div>


        <!-- Step Navigation Buttons -->
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <button type="button" class="btn btn-outline" @click="goToStep(2)">
            <span>ย้อนกลับ</span>
          </button>
          <button type="button" class="btn btn-primary" :disabled="submitting" @click="submitBooking" style="padding: 10px 24px; font-weight: 700;">
            <span v-if="submitting">กำลังบันทึกการจอง...</span>
            <span v-else style="display: inline-flex; align-items: center; gap: 8px;">
              <span>ยืนยันการจอง</span>
              <span>&rarr;</span>
            </span>
          </button>
        </div>
      </div>

      <!-- ============================================================ -->
      <!-- SUCCESS MODAL: VIP DIGITAL BOOKING PASS (WOW FEATURE 5)      -->
      <!-- ============================================================ -->
      <div v-if="showSuccessModal" class="modal-overlay" @click.self="closeSuccessModal">
        <canvas ref="confettiCanvas" class="confetti-canvas-overlay"></canvas>
        <div class="modal-content vip-booking-pass-card" style="width: min(440px, calc(100vw - 32px)); max-width: 440px; padding: 0; text-align: center; position: relative; overflow: hidden; border-radius: 20px; border: 1px solid #bfdbfe; box-shadow: 0 25px 60px rgba(15, 23, 42, 0.25);">
          <!-- Top Metallic Gradient Banner -->
          <div class="vip-pass-header" style="background: linear-gradient(135deg, #1e3a8a 0%, #2563eb 50%, #3b82f6 100%); color: #ffffff; padding: 24px 28px; position: relative; text-align: left;">
            <button
              class="btn-icon-square"
              style="position: absolute; top: 18px; right: 18px; width: 32px; height: 32px; background: rgba(255,255,255,0.2); border: none; color: #ffffff; backdrop-filter: blur(4px);"
              @click="closeSuccessModal"
              aria-label="ปิด"
            >
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
            </button>

            <div style="display: flex; align-items: center; gap: 12px; margin-bottom: 8px;">
              <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: #bfdbfe; flex-shrink: 0;">
                <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
                <polyline points="9 22 9 12 15 12 15 22"></polyline>
              </svg>
              <div>
                <div style="font-size: 0.72rem; letter-spacing: 0.1em; text-transform: uppercase; color: #bfdbfe; font-weight: 700;">สำนักงานตลาดกรุงเทพมหานคร (BMA)</div>
                <h3 style="font-size: 1.25rem; font-weight: 800; color: #ffffff; margin: 0;">
                  บัตรนัดหมายเข้าชมแผงเช่า
                </h3>
              </div>
            </div>

            <!-- Booking Reference ID Pill with Copy button (QA/UX) -->
            <div style="display: inline-flex; align-items: center; gap: 8px; background: rgba(255,255,255,0.18); border: 1px solid rgba(255,255,255,0.3); padding: 4px 12px; border-radius: 20px; font-size: 0.8rem; margin-top: 6px;">
              <span>รหัสอ้างอิง:</span>
              <strong style="color: #fef08a; letter-spacing: 0.05em;">{{ bookingRefCode }}</strong>
              <button
                type="button"
                style="background: none; border: none; color: #ffffff; cursor: pointer; font-size: 0.75rem; margin-left: 4px; display: inline-flex; align-items: center; gap: 4px; text-decoration: underline;"
                @click="copyRefCode"
              >
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                </svg>
                <span>{{ copiedRef ? 'คัดลอกแล้ว' : 'คัดลอก' }}</span>
              </button>
            </div>
          </div>

          <!-- Body Content -->
          <div style="padding: 24px 28px; background: #ffffff;">
            <!-- Countdown Banner -->
            <div style="background: #fefce8; border: 1px solid #fef08a; border-radius: 10px; padding: 12px 14px; margin-bottom: 20px; display: flex; align-items: center; gap: 12px;">
              <span style="flex: 1; min-width: 0; text-align: left; font-size: 0.82rem; line-height: 1.35; color: #854d0e; font-weight: 600; display: inline-flex; align-items: center; gap: 6px;">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="flex: 0 0 auto; color: #ca8a04;"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                <span>ยื่นคำขอนัดดูแผงสำเร็จ</span>
              </span>
              <span class="badge" style="flex: 0 0 auto; white-space: nowrap; background: #fef3c7; color: #b45309; border: 1px solid #fde68a; font-size: 0.72rem; font-weight: 700; display: inline-flex; align-items: center; gap: 5px; padding: 7px 10px;">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                <span>รอเจ้าหน้าที่ตรวจสอบ</span>
              </span>
            </div>

            <!-- Ticket Details Grid -->
            <div style="display: grid; grid-template-columns: 1fr; gap: 18px; text-align: left; margin-bottom: 20px;">
              <!-- Left: Specs -->
              <div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px; font-size: 0.82rem; display: flex; flex-direction: column; gap: 8px;">
                <div style="display: flex; justify-content: space-between;">
                  <span style="color: #64748b;">ตลาด:</span>
                  <strong style="color: #0f172a;">{{ currentMarketName }}</strong>
                </div>
                <div style="display: flex; justify-content: space-between;">
                  <span style="color: #64748b;">หมายเลขแผง:</span>
                  <strong style="color: #2563eb; font-size: 1rem;">
                    แผง {{ passDisplayStalls.length > 0 ? passDisplayStalls.map(stallDisplayNumber).join(', ') : stallDisplayNumber(selectedStall) }}
                    <span v-if="passDisplayStalls.length > 1" style="font-size: 0.8rem; font-weight: 600; color: #64748b;">({{ passDisplayStalls.length }} แผง)</span>
                  </strong>
                </div>
                <div v-if="passDisplayStalls[0]?.status === 'reserved' || currentSelectionStatus === 'reserved'" style="display: flex; justify-content: space-between; background: #fffbeb; padding: 6px 10px; border-radius: 6px; border: 1px solid #fde68a;">
                  <span style="color: #92400e; font-weight: 600;">สถานะคิวสำรอง:</span>
                  <strong v-if="currentWaitingInfo" style="color: #d97706;">คิวสำรองที่ {{ currentWaitingInfo.next_queue_number }}<span v-if="currentWaitingInfo.current_booker"> (ต่อจาก {{ currentWaitingInfo.current_booker.masked_name }})</span></strong>
                  <strong v-else style="color: #d97706;">กำลังตรวจสอบลำดับคิว</strong>
                </div>
                <div style="display: flex; justify-content: space-between;">
                  <span style="color: #64748b;">วันนัดหมาย:</span>
                  <strong style="color: #0f172a;">{{ formatThaiDate(bookingForm.visitDate) }}</strong>
                </div>
                <div style="display: flex; justify-content: space-between;">
                  <span style="color: #64748b;">ช่วงเวลา:</span>
                  <strong style="color: #0f172a;">{{ bookingForm.visitSlot }}</strong>
                </div>
                <div style="display: flex; justify-content: space-between;">
                  <span style="color: #64748b;">ชื่อผู้จอง:</span>
                  <strong style="color: #0f172a;">{{ bookingForm.firstName }} {{ bookingForm.lastName }}</strong>
                </div>
              </div>

            </div>

            <!-- Bottom Navigation -->
            <div style="display: flex; gap: 10px;">
              <NuxtLink to="/bookings" class="btn btn-primary" style="flex: 1; padding: 12px; font-weight: 700;">
                รายการจองของฉัน &rarr;
              </NuxtLink>
              <NuxtLink to="/" class="btn btn-outline" style="padding: 12px; font-weight: 600;">
                กลับหน้าหลัก
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================ -->
      <!-- MODAL 1: STALL PHOTO PREVIEW (ดูตัวอย่างภาพ / เลื่อนดูหลายรูป) -->
      <!-- ============================================================ -->
      <div v-if="showVirtualModal" class="modal-overlay photo-preview-overlay" @click.self="closePhotoPreview">
        <div class="photo-preview" role="dialog" aria-label="ตัวอย่างภาพแผง">
          <div class="pp-header">
            <div>
              <div class="pp-title">แผง {{ stallDisplayNumber(previewStall) }}</div>
              <div class="pp-sub">{{ previewStall?.zone }}</div>
            </div>
            <button type="button" class="pp-close" @click="closePhotoPreview" aria-label="ปิด">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
            </button>
          </div>

          <div class="pp-stage">
            <button v-if="previewPhotos.length > 1" type="button" class="pp-nav pp-prev" @click="stepPreview(-1)" aria-label="รูปก่อนหน้า">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>
            </button>
            <img :key="previewPhotos[previewIndex]" :src="previewPhotos[previewIndex]" :alt="`ภาพแผง ${stallDisplayNumber(previewStall)} รูปที่ ${previewIndex + 1}`" class="pp-image" />
            <button v-if="previewPhotos.length > 1" type="button" class="pp-nav pp-next" @click="stepPreview(1)" aria-label="รูปถัดไป">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
            </button>
            <span v-if="previewPhotos.length > 1" class="pp-counter">{{ previewIndex + 1 }} / {{ previewPhotos.length }}</span>
          </div>

          <div v-if="previewPhotos.length > 1" class="pp-thumbs">
            <button
              v-for="(src, i) in previewPhotos"
              :key="src"
              type="button"
              class="pp-thumb"
              :class="{ active: i === previewIndex }"
              @click="previewIndex = i"
              :aria-label="`ดูรูปที่ ${i + 1}`"
            >
              <img :src="src" alt="" />
            </button>
          </div>
        </div>
      </div>

      <!-- ============================================================ -->
      <!-- MODAL 2: STALL COMPARISON (WOW FEATURE 7)                    -->
      <!-- ============================================================ -->
      <div v-if="showCompareModal" class="modal-overlay" @click.self="showCompareModal = false">
        <div class="modal-content" style="max-width: 640px; padding: 24px; border-radius: 18px; border: 1px solid #bfdbfe;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
            <div style="display: flex; align-items: center; gap: 8px;">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: #2563eb;">
                <line x1="18" y1="20" x2="18" y2="10"></line>
                <line x1="12" y1="20" x2="12" y2="4"></line>
                <line x1="6" y1="20" x2="6" y2="14"></line>
              </svg>
              <h3 style="font-size: 1.15rem; font-weight: 800; color: #0f172a; margin: 0;">เปรียบเทียบแผงค้า (Stall Comparison)</h3>
            </div>
            <button class="btn-icon-square" style="width: 32px; height: 32px;" @click="showCompareModal = false" aria-label="ปิด">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
            </button>
          </div>

          <!-- Comparison Table -->
          <table style="width: 100%; border-collapse: collapse; font-size: 0.85rem; margin-bottom: 20px;">
            <thead>
              <tr style="background: #f1f5f9; border-bottom: 2px solid #cbd5e1;">
                <th style="padding: 10px; text-align: left; color: #475569;">คุณลักษณะ</th>
                <th style="padding: 10px; text-align: center; color: #2563eb; background: #eff6ff;">แผง 20A (ที่เลือกไว้)</th>
                <th style="padding: 10px; text-align: center; color: #475569;">แผง 20B (แผงข้างเคียง)</th>
              </tr>
            </thead>
            <tbody>
              <tr style="border-bottom: 1px solid #e2e8f0;">
                <td style="padding: 10px; color: #64748b;">ค่าเช่ารายเดือน</td>
                <td style="padding: 10px; text-align: center; font-weight: 800; color: #059669; background: #f0fdf4;">฿1,800</td>
                <td style="padding: 10px; text-align: center; font-weight: 700; color: #0f172a;">฿2,000</td>
              </tr>
              <tr style="border-bottom: 1px solid #e2e8f0;">
                <td style="padding: 10px; color: #64748b;">ขนาดพื้นที่</td>
                <td style="padding: 10px; text-align: center; font-weight: 700; color: #0f172a; background: #eff6ff;">5.0 ตร.ม.</td>
                <td style="padding: 10px; text-align: center; font-weight: 700; color: #0f172a;">5.5 ตร.ม.</td>
              </tr>
              <tr style="border-bottom: 1px solid #e2e8f0;">
                <td style="padding: 10px; color: #64748b;">ตำแหน่งทางสัญจร</td>
                <td style="padding: 10px; text-align: center; font-weight: 700; color: #2563eb; background: #eff6ff;">ติดเวทีหลัก &amp; ทางเดินใหญ่ (ทำเลทอง)</td>
                <td style="padding: 10px; text-align: center; color: #475569;">แผงถัดไป</td>
              </tr>
              <tr style="border-bottom: 1px solid #e2e8f0;">
                <td style="padding: 10px; color: #64748b;">ปริมาณคนเดินผ่าน</td>
                <td style="padding: 10px; text-align: center; font-weight: 700; color: #059669; background: #eff6ff;">สูงมาก (High Traffic)</td>
                <td style="padding: 10px; text-align: center; color: #475569;">ปานกลาง (Medium)</td>
              </tr>
              <tr>
                <td style="padding: 10px; color: #64748b;">ทิศแดด/ลม</td>
                <td style="padding: 10px; text-align: center; font-weight: 700; color: #059669; background: #eff6ff;">ร่มตลอดวัน มีลมพัด</td>
                <td style="padding: 10px; text-align: center; color: #475569;">แดดบ่ายเล็กน้อย</td>
              </tr>
            </tbody>
          </table>

          <div style="text-align: right;">
            <button class="btn btn-primary" style="padding: 10px 24px; font-weight: 700;" @click="showCompareModal = false">
              คงการเลือกแผง 20A &rarr;
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppSelect, { type SelectOption } from '~/components/AppSelect.vue'
import {
  getBookingVisitDates,
  setCachedScheduleBookings,
  type AvailableDateSlot as ScheduleAvailableDateSlot
} from '../composables/useMarketVisitSchedule'
import { BMA_MARKETS_MASTER, getMarketGisConfig } from '../data/markets_registry'
import { useNotifications } from '~/composables/useNotifications'

const route = useRoute()
const router = useRouter()
const { get, post } = useApi()
const { user } = useAuth()
const { showAlert, showWarning, showError, showSuccess, showConfirm } = useAppAlert()
const { fetchNotifications } = useNotifications()

const overviewMapRef = ref<any>(null)

const focusClaimedStallOnMap = (): void => {
  if (!selectedStalls.value.length) return
  nextTick(() => {
    overviewMapRef.value?.resize?.()
    const focused = overviewMapRef.value?.focusSelectedStalls?.()
    if (!focused) {
      window.setTimeout(() => {
        overviewMapRef.value?.resize?.()
        overviewMapRef.value?.focusSelectedStalls?.()
      }, 250)
    }
  })
}

const businessOptions = [
  { value: 'เสื้อผ้าแฟชั่น', label: 'เสื้อผ้าแฟชั่นและเครื่องประดับ' },
  { value: 'อาหารปรุงสุก', label: 'อาหารปรุงสุกและเครื่องดื่ม' },
  { value: 'งานศิลปะและของสะสม', label: 'งานศิลปะและของสะสม' },
  { value: 'ของตกแต่งบ้านและงานคราฟต์', label: 'ของตกแต่งบ้านและงานคราฟต์' },
  { value: 'เซรามิกและเบญจรงค์', label: 'เซรามิกและเบญจรงค์' },
  { value: 'ของใช้เบ็ดเตล็ด', label: 'ของใช้เบ็ดเตล็ด' }
]

const { applyMarketOverride } = useMarketOverrides()
const initialMarketConfig = getMarketGisConfig(String(route.query.market || '1'))
const initialOverriddenMarket = applyMarketOverride({
  id: initialMarketConfig.id,
  slug: initialMarketConfig.slug,
  name: initialMarketConfig.name,
  district: initialMarketConfig.district
})

const currentMarketName = ref(initialOverriddenMarket.name || initialMarketConfig.name)
const currentMarketDistrict = ref(initialOverriddenMarket.district || initialMarketConfig.district)
const isMarketFull = ref(false)
const currentMarket = ref<any>(null)

const allMarketOptions = computed<SelectOption[]>(() => {
  return Object.values(BMA_MARKETS_MASTER).map(m => ({
    value: String(m.id),
    label: `${m.name} (${m.district})`
  }))
})

const handleMarketSelect = (val: string | number) => {
  if (!val) return
  router.push({
    query: {
      ...route.query,
      market: String(val)
    }
  })
  selectedProject.value = null
  bookedProject.value = null
  selectedStalls.value = []
}

const currentStep = ref<number>(1)
const selectedCategoryFilter = ref<string | null>(null)
// เริ่มต้นเป็น null เพื่อให้แผนที่แสดงภาพรวมตลาดทั้งหมด (ไม่เจาะจงซูมเข้าโครงการใดโครงการหนึ่งตั้งแต่ต้น)
const selectedProject = ref<{
  id: string
  name: string
  category: string
  color: string
  stallsAvailable: number
  totalStalls: number
} | null>(null)
const bookedProject = ref<{
  id: string
  name: string
  category?: string
  color?: string
  stallsAvailable?: number
  totalStalls?: number
} | null>(null)
const isSelectedProjectFull = computed(() => selectedProject.value?.stallsAvailable === 0)

// Step 1 starts with all stalls displayed by default
const onlyAvailable = ref(false)
const showTrafficHeatmap = ref(false)
const zoneFilter = ref<'all' | 'hotspot' | 'stage'>('all')
const detailsOpen = ref(true)
const zoomLevel = ref(1)
const mapStyle = ref<'standard' | 'dark' | 'satellite'>('standard')
const submitting = ref(false)
const showSuccessModal = ref(false)
const showVirtualModal = ref(false)

// ---- Stall photo preview (ดูตัวอย่างภาพ) ----
// แผงที่ยังไม่มีรูปของตัวเอง ใช้ชุดภาพตัวอย่างนี้ (เปลี่ยนเป็นรูปจริงจาก backend ได้ผ่าน field photos)
const DEFAULT_STALL_GALLERY = [
  '/images/stall_3d_perspective.jpg',
  '/images/market_stall_sample.jpg',
  '/images/market_hall_front.jpg'
]
const previewStallCode = ref('')
const previewIndex = ref(0)

const getStallPhotos = (s?: StallModel | null): string[] => {
  if (!s) return DEFAULT_STALL_GALLERY
  if (Array.isArray((s as any).photos) && (s as any).photos.length > 0) return (s as any).photos
  if (s.photo) {
    return [
      s.photo,
      '/images/market_stall_sample.jpg',
      '/images/stall_3d_perspective.jpg'
    ]
  }
  return DEFAULT_STALL_GALLERY
}
const currentPreviewStall = ref<StallModel | null>(null)
const previewStall = computed(() => currentPreviewStall.value || selectedStalls.value.find(s => s.code === previewStallCode.value) || selectedStall.value)
const previewPhotos = computed(() => getStallPhotos(previewStall.value))

const openPhotoPreview = (s: StallModel, index = 0) => {
  currentPreviewStall.value = s
  previewStallCode.value = s.code
  previewIndex.value = index
  showVirtualModal.value = true
}
const closePhotoPreview = () => { showVirtualModal.value = false }
const stepPreview = (delta: number) => {
  const n = previewPhotos.value.length
  if (n < 2) return
  previewIndex.value = (previewIndex.value + delta + n) % n
}
const showCompareModal = ref(false)

// BA commercial choices
const paymentSchedule = ref<'monthly' | 'quarterly' | 'annual'>('monthly')
const bmaRulesOpen = ref(false)
const copiedRef = ref(false)


const idVerified = ref(false)

// Time Slots & Officer Available Dates relational model
interface TimeSlotItem {
  time: string
  capacityText: string
  capacityClass: string
  note: string
  isFull: boolean
  enabled?: boolean
  capacity?: number
  bookedCount?: number
  bookedStalls?: string[]
  remainingQueues: number
  hasTenantConflict?: boolean
  tenantConflictInfo?: {
    bookingNumber: string
    marketName: string
    projectName: string
    stallNumbers: string[]
  } | null
}

interface AvailableDateSlot {
  date: string
  dayName: string
  dayNum: string
  monthShort: string
  badge?: string
  badgeType?: string
  slotsText: string
  statusClass: string
  isPopular?: boolean
  description?: string
  isAllFull?: boolean
  isClosed?: boolean
  totalCapacity?: number
  totalBooked?: number
  totalRemainingQueues?: number
  slots: TimeSlotItem[]
}

const currentMarketId = computed(() => String(route.query.market || '1'))
const allMarketBookings = ref<any[]>([])

const normalizePersonName = (val: any): string =>
  String(val || '')
    .replace(/^(คุณ|นาย|นางสาว|นาง|น\.ส\.)\s*/g, '')
    .replace(/\([^)]*\)/g, '')
    .replace(/\s+/g, ' ')
    .trim()
    .toLowerCase()

const normalizeDigitsOnly = (val: any): string => String(val || '').replace(/\D/g, '')

const formErrors = ref({
  business: '',
  visitDate: '',
  visitSlot: ''
})

// เริ่มต้น: ข้อมูลโปรไฟล์ผู้เช่าดึงมาแสดงอัตโนมัติ (ห้ามว่าง) ส่วนฟิลด์ที่ต้องกรอก/เลือกเองให้เป็นค่าว่างทั้งหมด
const bookingForm = ref({
  firstName: 'สมชาย',
  lastName: 'ใจดี',
  idCard: '1-1004-98721-34-5',
  phone: '087-567-4497',
  email: 'somchai.jai@thaigov.or.th',
  business: '',
  notes: '',
  visitDate: '',
  visitSlot: ''
})

const isBookingOwnedByCurrentTenant = (b: any): boolean => {
  if (!b || ['rejected', 'not_interested', 'cancelled', 'no_show'].includes(b.status)) return false
  const currentUser = user.value as any
  if (currentUser?.id !== undefined && currentUser?.id !== null && b.tenant_user_id !== undefined && b.tenant_user_id !== null) {
    return String(b.tenant_user_id) === String(currentUser.id)
  }

  const demoOwnerNames: Record<string, string> = {
    tenant: 'สมชาย ใจดี',
    tenant2: 'นภาพร รัตนศิลป์',
    tenant3: 'วิชัย เจริญทรัพย์'
  }
  const demoName = normalizePersonName(demoOwnerNames[String(currentUser?.username || '')])
  const bName = normalizePersonName(b.tenant_name || b.applicant_name)
  const bPhone = normalizeDigitsOnly(b.tenant_phone || b.applicant_phone)
  const bIdCard = normalizeDigitsOnly(b.tenant_id_card || b.citizen_id)

  const formFullName = normalizePersonName(`${bookingForm.value?.firstName || ''} ${bookingForm.value?.lastName || ''}`)
  const userFullName = normalizePersonName(currentUser?.full_name)
  const candidateNames = new Set([userFullName, demoName, formFullName].filter(Boolean))

  const formPhone = normalizeDigitsOnly(bookingForm.value?.phone)
  const userPhone = normalizeDigitsOnly(currentUser?.phone || currentUser?.phone_number)
  const candidatePhones = new Set([userPhone, formPhone].filter(p => p.length >= 9))

  const formIdCard = normalizeDigitsOnly(bookingForm.value?.idCard)
  const userIdCard = normalizeDigitsOnly(currentUser?.citizen_id || currentUser?.id_card)
  const candidateIdCards = new Set([userIdCard, formIdCard].filter(id => id.length === 13))

  if (bName && candidateNames.has(bName)) return true
  if (bPhone.length >= 9 && candidatePhones.has(bPhone)) return true
  if (bIdCard.length === 13 && candidateIdCards.has(bIdCard)) return true
  return false
}

const loadAvailableDatesFromSchedule = (): AvailableDateSlot[] => {
  return getBookingVisitDates(
    currentMarketId.value,
    14,
    allMarketBookings.value,
    currentMarketName.value,
    isBookingOwnedByCurrentTenant
  )
}

const availableDates = ref<AvailableDateSlot[]>(loadAvailableDatesFromSchedule())

// Auto-sync whenever marketId or marketName changes
watch([currentMarketId, currentMarketName], () => {
  availableDates.value = loadAvailableDatesFromSchedule()
  if (bookingForm.value.visitDate) {
    const cur = availableDates.value.find(d => d.date === bookingForm.value.visitDate)
    if (!cur || cur.isClosed) {
      const firstOpen = availableDates.value.find(d => !d.isClosed && !d.isAllFull)
      bookingForm.value.visitDate = firstOpen?.date || ''
    }
  }
})

const currentSelectedDateAdvice = computed(() => {
  const d = availableDates.value.find(item => item.date === bookingForm.value.visitDate)
  return d?.description || ''
})

const minDate = new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString().slice(0, 10)
const defaultVisitDate = (availableDates.value.find(d => d.isPopular && !d.isClosed) || availableDates.value.find(d => !d.isClosed) || availableDates.value[0])?.date || minDate

// Dynamic Time Slots computed from selected date (Relational Link)
const currentTimeSlots = computed<TimeSlotItem[]>(() => {
  const selectedDateObj = availableDates.value.find(d => d.date === bookingForm.value.visitDate)
  if (selectedDateObj && selectedDateObj.slots && selectedDateObj.slots.length > 0) {
    return selectedDateObj.slots
  }
  return [
    { time: '09.30 - 10.30 น.', capacityText: 'ว่าง 2 คิว', capacityClass: 'cap-med', note: 'ช่วงเช้า อากาศสบาย', isFull: false, remainingQueues: 2 },
    { time: '11.00 - 12.00 น.', capacityText: 'ว่าง 4 คิว', capacityClass: 'cap-high', note: 'ช่วงสาย บรรยากาศเริ่มคึกคัก', isFull: false, remainingQueues: 4 },
    { time: '13.00 - 14.00 น.', capacityText: 'ว่าง 2 คิว', capacityClass: 'cap-med', note: 'ช่วงบ่าย มีเจ้าหน้าที่พาชม', isFull: false, remainingQueues: 2 },
    { time: '15.30 - 16.30 น.', capacityText: 'ว่าง 3 คิว', capacityClass: 'cap-high', note: 'ช่วงเย็น ชมทราฟฟิกจริง', isFull: false, remainingQueues: 3 }
  ]
})

const isCurrentDateAllFull = computed(() => {
  const selectedDateObj = availableDates.value.find(d => d.date === bookingForm.value.visitDate)
  if (!selectedDateObj) return false
  return selectedDateObj.isAllFull || (selectedDateObj.slots && selectedDateObj.slots.every(s => s.isFull))
})

const selectVisitDate = (d: AvailableDateSlot) => {
  // Treat the selected date as a toggle. Clearing the date must also clear its
  // dependent time slot so an old appointment time cannot be submitted.
  if (bookingForm.value.visitDate === d.date) {
    bookingForm.value.visitDate = ''
    bookingForm.value.visitSlot = ''
    formErrors.value.visitDate = ''
    formErrors.value.visitSlot = ''
    return
  }
  if (d.isClosed) {
    formErrors.value.visitDate = `วันที่ ${d.dayName} ที่ ${d.dayNum} ${d.monthShort} ปิดรับนัดหมายเข้าชม กรุณาเลือกวันอื่นที่เปิดให้บริการ`
    return
  }
  bookingForm.value.visitDate = d.date
  formErrors.value.visitDate = ''
}

const selectTimeSlot = (slot: TimeSlotItem) => {
  // A second click on the active slot deselects it, matching the date-card
  // behaviour and letting the tenant change their mind before proceeding.
  if (bookingForm.value.visitSlot === slot.time) {
    bookingForm.value.visitSlot = ''
    formErrors.value.visitSlot = ''
    return
  }
  if (!bookingForm.value.visitDate) {
    formErrors.value.visitDate = 'กรุณาคลิกเลือกวันที่เปิดให้เข้าชมแผงในข้อ 1 ด้านบนก่อนเลือกรอบเวลา'
    return
  }
  if (slot.enabled === false) {
    formErrors.value.visitSlot = `รอบเวลา ${slot.time} งดให้บริการในวันที่เลือก กรุณาเลือกรอบเวลาอื่น`
    return
  }
  if (slot.hasTenantConflict) {
    const cInfo = slot.tenantConflictInfo
    const placeText = cInfo
      ? `${cInfo.projectName} (${cInfo.marketName}${cInfo.stallNumbers.length ? ' • แผง ' + cInfo.stallNumbers.join(', ') : ''})`
      : 'โครงการอื่น'
    formErrors.value.visitSlot = `คุณมีนัดหมายเข้าชมแผงที่ "${placeText}" ในรอบเวลา ${slot.time} อยู่แล้ว ไม่สามารถนัดเวลาชนกันได้ กรุณาเลือกรอบเวลาอื่น`
    return
  }
  if (slot.isFull) {
    formErrors.value.visitSlot = `รอบเวลา ${slot.time} ของวันที่เลือก มีผู้ลงทะเบียนเต็มโควตาแล้ว (${slot.bookedCount || slot.capacity}/${slot.capacity} คิว) กรุณาเลือกรอบเวลาอื่น`
    return
  }
  bookingForm.value.visitSlot = slot.time
  formErrors.value.visitSlot = ''
}

watch(() => bookingForm.value.business, (val) => {
  if (val) formErrors.value.business = ''
})

// Auto-switch slot if current slot is full or conflicts on the newly selected date
watch(() => bookingForm.value.visitDate, (newDate) => {
  if (!newDate) {
    bookingForm.value.visitSlot = ''
    return
  }
  const selectedDateObj = availableDates.value.find(d => d.date === newDate)
  if (selectedDateObj) {
    const slots = selectedDateObj.slots || []
    const currentSlotObj = slots.find(s => s.time === bookingForm.value.visitSlot)
    // หากเคยเลือกรอบไว้แล้วและรอบนั้นเต็มหรือชนนัดตัวเองในวันใหม่ ให้สลับไปรอบแรกที่ยังว่าง
    if (bookingForm.value.visitSlot && (!currentSlotObj || currentSlotObj.isFull || currentSlotObj.hasTenantConflict)) {
      const firstAvailable = slots.find(s => !s.isFull && !s.hasTenantConflict)
      bookingForm.value.visitSlot = firstAvailable ? firstAvailable.time : ''
    }
  }
})

const bookingRefCode = ref('BMA-JJ-2026-08942')

// Auto-fill from authenticated user profile
watch(user, (newU: any) => {
  if (newU) {
    if (newU.full_name) {
      const cleanedFullName = String(newU.full_name)
        .replace(/\([^)]*\)/g, '')
        .replace(/^(คุณ|นาย|นางสาว|นาง|น\.ส\.)\s*/g, '')
        .trim()
      const parts = cleanedFullName.split(/\s+/).filter(Boolean)
      if (parts.length > 0) bookingForm.value.firstName = parts[0]
      if (parts.length > 1) bookingForm.value.lastName = parts.slice(1).join(' ')
    } else if (newU.first_name) {
      bookingForm.value.firstName = newU.first_name
      if (newU.last_name) bookingForm.value.lastName = newU.last_name
    }
    if (newU.citizen_id || newU.id_card) {
      bookingForm.value.idCard = formatThaiId(newU.citizen_id || newU.id_card)
    }
    if (newU.phone_number || newU.phone) {
      bookingForm.value.phone = newU.phone_number || newU.phone
    }
    if (newU.email) {
      bookingForm.value.email = newU.email
    }
  }
}, { immediate: true })

// Signature Pad State & Handlers
const signatureCanvas = ref<HTMLCanvasElement | null>(null)
const signatureDrawn = ref(false)
let isDrawing = false
let sigCtx: CanvasRenderingContext2D | null = null

// Representative stalls array matching screenshot
interface StallModel {
  selectionId?: string
  code: string
  projectId?: string
  zone: string
  floorName?: string
  status: 'available' | 'occupied' | 'reserved'
  areaType: string
  size: number
  entryFee: number
  monthlyRent: number
  commonFee: number
  depositFee: number
  photo: string
  photos?: string[]
  backendId?: number
}

interface BackendStallRecord {
  id: number
  ident?: string
  stall_ident?: string
  stall_number?: string
  project_name?: string
  monthly_rent?: number
  status?: string
  stall_status?: string
  waitlist_count?: number
  offered_to_tenant_id?: number
  offered_expires_at?: string
}

const backendStallsByIdent = ref<Record<string, BackendStallRecord>>({})
const myWaitlistByStall = ref<Record<string, any>>({})

const activeWaitlistOffer = computed<any | null>(() => {
  const targetCode = String(selectedStalls.value[0]?.code || route.query.stall || '').trim()
  if (!targetCode) return null
  const item = myWaitlistByStall.value[targetCode]
  return item?.status === 'offered' ? item : null
})

const waitlistOfferDeadlineText = computed(() => {
  const raw = activeWaitlistOffer.value?.expires_at || activeWaitlistOffer.value?.offered_expires_at
  if (!raw) return ''
  const expiresAt = new Date(raw)
  if (Number.isNaN(expiresAt.getTime())) return ''
  return `ถึงวันที่ ${expiresAt.toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: 'numeric' })} เวลา ${expiresAt.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' })} น.`
})

const getBackendStallIdent = (stall: BackendStallRecord): string => (
  String(stall.ident || stall.stall_ident || stall.stall_number || '').trim()
)

const isStallOfferedToMe = (ident?: string, stall?: BackendStallRecord): boolean => {
  const currentUserId = Number((user.value as any)?.id || 0)
  if (currentUserId > 0 && stall?.offered_to_tenant_id && Number(stall.offered_to_tenant_id) === currentUserId) {
    return true
  }
  if (ident && myWaitlistByStall.value[ident]?.status === 'offered') {
    return true
  }
  return false
}

const normalizeBackendStallStatus = (status?: string, stall?: BackendStallRecord): StallModel['status'] => {
  const ident = stall ? getBackendStallIdent(stall) : ''
  if (status === 'offered') {
    return isStallOfferedToMe(ident, stall) ? 'available' : 'reserved'
  }
  if (status === 'available' || status === 'reserved' || status === 'occupied') return status
  return 'available'
}

const getBackendRecordStatus = (stall?: BackendStallRecord): string | undefined => (
  // `status` is the lifecycle flag for the database row (normally "A").
  // Availability is carried by `stall_status`, which changes on booking.
  stall?.stall_status || stall?.status
)

const stallStatusOverrides = computed<Record<string, StallModel['status']>>(() => {
  const statuses: Record<string, StallModel['status']> = {}
  Object.entries(backendStallsByIdent.value).forEach(([ident, stall]) => {
    const rawStatus = getBackendRecordStatus(stall)
    if (!rawStatus) return
    const status = normalizeBackendStallStatus(rawStatus, stall)
    if (status !== 'available') statuses[ident] = status
  })
  return statuses
})

const refreshBackendStalls = async (): Promise<void> => {
  try {
    const marketParam = currentMarket.value?.slug || currentMarketId.value
    const url = marketParam ? `/stalls?market=${encodeURIComponent(marketParam)}` : '/stalls'
    const [response, waitlistRes] = await Promise.all([
      get<any>(url),
      get<any>('/bookings/waitlist/mine').catch(() => null)
    ])
    const waitListItems: any[] = Array.isArray(waitlistRes)
      ? waitlistRes
      : (Array.isArray(waitlistRes?.data) ? waitlistRes.data : [])
    const wMap: Record<string, any> = {}
    waitListItems.forEach((w: any) => {
      const code = String(w.stall_number || w.stall_ident || '').trim()
      if (code && ['waiting', 'offered'].includes(w.status)) wMap[code] = w
    })
    myWaitlistByStall.value = wMap

    const list: BackendStallRecord[] = Array.isArray(response)
      ? response
      : (Array.isArray(response?.data) ? response.data : [])
    const records: Record<string, BackendStallRecord> = {}
    list.forEach(stall => {
      const ident = getBackendStallIdent(stall)
      const rawStatus = getBackendRecordStatus(stall)
      if (!ident || !rawStatus) return
      records[ident] = stall
    })
    backendStallsByIdent.value = records
  } catch (error) {
    console.warn('Could not refresh stall statuses:', error)
  }
}


const handleBackendStallSync = () => {
  void refreshBackendStalls()
  void refreshMyActiveBookings()
}

// Track active bookings belonging to the current tenant so the same stall cannot be booked/queued twice
const myActiveBookingsByStall = ref<Record<string, any>>({})
const sessionBookedStallMap = ref<Record<string, any>>({})
const confirmedPassStalls = ref<StallModel[]>([])
const passDisplayStalls = computed<StallModel[]>(() => (
  selectedStalls.value.length > 0 ? selectedStalls.value : confirmedPassStalls.value
))

const isQueueBookingRecord = (b?: any): boolean => {
  if (!b) return false
  if (['rejected', 'cancelled', 'not_interested', 'no_show', 'lease_applied', 'completed', 'terminated'].includes(b.status)) {
    return false
  }
  return b.status === 'waiting_queue' || (!b.holding_stall && Number(b.queue_position || b.queue_number) > 0)
}

const getDuplicateBookingStatusLabel = (status?: string, existingBooking?: any): string => {
  if (isQueueBookingRecord(existingBooking)) {
    const qPos = Number(existingBooking?.queue_position || existingBooking?.queue_number) || 1
    return status === 'approved'
      ? `คิวสำรองที่ ${qPos} • นัดหมายสำเร็จ`
      : `คิวสำรองที่ ${qPos} • รอเจ้าหน้าที่อนุมัติวันนัด`
  }
  switch (status) {
    case 'waiting_queue': return 'อยู่ในคิวสำรอง'
    case 'pending_approval':
    case 'pending': return 'รอเจ้าหน้าที่อนุมัติวันนัด'
    case 'pending_slot': return 'รอเลือกวันนัดใหม่'
    case 'approved': return 'นัดหมายสำเร็จ (รอเข้าชมแผง)'
    case 'pending_decision': return 'เข้าชมแผงแล้ว (รอตัดสินใจ)'
    case 'interested': return 'ยืนยันสนใจเช่าแผงแล้ว'
    case 'lease_applied':
    case 'completed': return 'ส่งคำขอเช่าแผงแล้ว'
    default: return 'มีรายการจองที่กำลังดำเนินการอยู่'
  }
}

const showDuplicateBookingWarning = (stallCode: string, existingBooking?: any) => {
  const bNo = existingBooking?.booking_number || existingBooking?.booking_no || ''
  const statusLabel = getDuplicateBookingStatusLabel(existingBooking?.status, existingBooking)
  const refSuffix = bNo ? ` (รหัสการจอง ${bNo} • สถานะ: ${statusLabel})` : ` (${statusLabel})`
  showWarning(
    `คุณได้ทำการจองหรืออยู่ในคิวสำรองของแผง ${stallCode} ไว้แล้ว${refSuffix} ไม่สามารถทำรายการจองแผงเดิมซ้ำได้ กรุณาตรวจสอบที่หน้า "รายการจองของฉัน"`,
    'ไม่สามารถจองแผงซ้ำได้'
  )
}

const refreshMyActiveBookings = async (): Promise<void> => {
  try {
    const response = await get<any>('/bookings')
    const list: any[] = Array.isArray(response)
      ? response
      : (Array.isArray(response?.data) ? response.data : [])

    // The regular endpoint is deliberately scoped to the signed-in tenant.
    // Fetch only anonymous per-slot totals separately, so slot availability is
    // shared by every tenant without exposing another tenant's booking details.
    let scheduleBookings = list
    try {
      const targetMarket = currentMarket.value?.id ? String(currentMarket.value.id) : (currentMarketId.value || '1')
      const capacityResponse = await get<any>(
        `/bookings/visit-capacity?market_id=${encodeURIComponent(targetMarket)}&market=${encodeURIComponent(currentMarketId.value)}`
      )
      const capacityItems: any[] = Array.isArray(capacityResponse)
        ? capacityResponse
        : (Array.isArray(capacityResponse?.data) ? capacityResponse.data : [])

      scheduleBookings = [
        ...capacityItems.flatMap((item) => {
          const bookedCount = Math.max(0, Number(item?.booked_count) || 0)
          return Array.from({ length: bookedCount }, () => ({
            market_id: item.market_id || targetMarket,
            market_name: item.market_name || currentMarketName.value,
            visit_date: item.visit_date,
            visit_time_slot: item.visit_time_slot,
            status: 'pending_approval',
            __visit_capacity_record: true
          }))
        }),
        // Keep the tenant's own records for the "you already have an
        // appointment at this time" validation; capacity is counted above.
        ...list.map(booking => ({ ...booking, __tenant_private_record: true }))
      ]
    } catch {
      // Older backends do not yet provide the aggregate endpoint. Retain the
      // tenant's own count rather than presenting an unusable booking screen.
    }

    allMarketBookings.value = scheduleBookings
    setCachedScheduleBookings(scheduleBookings)

    const map: Record<string, any> = { ...sessionBookedStallMap.value }
    list.forEach((b) => {
      const codes = String(b.stall_number || b.stall_ident || '')
        .split(',')
        .map(c => c.trim())
        .filter(Boolean)
      codes.forEach((code) => {
        if (['rejected', 'not_interested', 'cancelled', 'no_show'].includes(b.status)) {
          if (map[code]?.id === b.id || map[code]?.booking_number === b.booking_number) {
            delete map[code]
            delete sessionBookedStallMap.value[code]
          }
          return
        }
        if (isBookingOwnedByCurrentTenant(b) || sessionBookedStallMap.value[code]) {
          map[code] = b
        }
      })
    })
    myActiveBookingsByStall.value = map
    availableDates.value = loadAvailableDatesFromSchedule()

    if (bookingForm.value.visitDate && bookingForm.value.visitSlot) {
      const selectedDateObj = availableDates.value.find(d => d.date === bookingForm.value.visitDate)
      const currentSlotObj = selectedDateObj?.slots?.find(s => s.time === bookingForm.value.visitSlot)
      if (currentSlotObj?.isFull || currentSlotObj?.hasTenantConflict) {
        const firstAvailable = selectedDateObj?.slots?.find(s => !s.isFull && !s.hasTenantConflict)
        bookingForm.value.visitSlot = firstAvailable ? firstAvailable.time : ''
      }
    }
  } catch (err) {
    myActiveBookingsByStall.value = { ...sessionBookedStallMap.value }
    availableDates.value = loadAvailableDatesFromSchedule()
  }
}

// Upper hall stalls grid
const upperHallStalls = ref<StallModel[]>([])
// Lower hall stalls grid
const lowerHallStalls = ref<StallModel[]>([])

// Multi-Stall Selection & Status Handling (BA & QA Requirement)
const selectedStalls = ref<StallModel[]>([])
const selectedStall = computed(() => selectedStalls.value[0] || null)
const selectedProjectsSummary = computed(() => {
  const uniqueZones = Array.from(new Set(selectedStalls.value.map(s => s.zone?.trim()).filter(Boolean)))
  if (uniqueZones.length > 0) return uniqueZones.join(', ')
  return selectedProject.value?.name || ''
})
const stallSelectionKey = (stall: StallModel): string => stall.selectionId || stall.code
const stallDisplayNumber = (stall?: any): string => {
  if (!stall) return '-'
  const sourceNumber = stall.sourceProperties?.number ?? stall.sourceProperties?.Number
  return String(sourceNumber ?? stall.name ?? stall.code ?? '-').trim() || '-'
}

// Current Selection Status ('available' or 'reserved')
const currentSelectionStatus = computed<'available' | 'reserved' | null>(() => {
  if (selectedStalls.value.length === 0) return null
  return selectedStalls.value[0].status === 'reserved' ? 'reserved' : 'available'
})

const selectedWaitingEntry = computed<any | null>(() => {
  const code = selectedStalls.value[0]?.code
  if (!code) return null
  const item = myWaitlistByStall.value[code]
  return item?.status === 'waiting' ? item : null
})

// Status Conflict Warning Modal State
const showConflictModal = ref(false)
const pendingConflictingStall = ref<StallModel | null>(null)

// Waiting List Queue & PDPA Masked Booker Interface
interface WaitingQueueBooker {
  queue_number: number
  booking_number: string
  masked_name: string
  masked_phone: string
  visit_date?: string
  visit_time_slot?: string
  status: string
  status_label: string
  created_at?: string
  days_left_estimate: number
}

interface WaitingQueueData {
  stall_id?: number | null
  stall_number: string
  queue_count?: number
  current_position?: number
  total_waiting_count: number
  next_queue_number: number
  has_waiting_list: boolean
  current_booker?: WaitingQueueBooker | null
  queues: WaitingQueueBooker[]
}

const stallWaitingQueueCache = ref<Record<string, WaitingQueueData>>({})
const loadingWaitingQueue = ref(false)
const activeWaitingStallCode = ref<string>('')

const currentWaitingInfo = computed<WaitingQueueData | null>(() => {
  const code = activeWaitingStallCode.value || selectedStalls.value[0]?.code
  if (!code) return null
  return stallWaitingQueueCache.value[code] || null
})

const fetchWaitingQueueForStall = async (stallCode: string) => {
  if (!stallCode) return
  if (stallWaitingQueueCache.value[stallCode]) return
  loadingWaitingQueue.value = true
  try {
    const rawRes = await get<any>(`/bookings/waiting-queue/${encodeURIComponent(stallCode)}`)
    const data = rawRes?.data || rawRes
    if (data) {
      const maskedName = data.current_booker?.masked_name || data.masked_current_booker
      const currentBooker = (maskedName && maskedName !== '-') ? {
        masked_name: maskedName,
        masked_phone: data.current_booker?.masked_phone || data.masked_booker_phone || '',
        status: data.current_booker?.status || data.status || 'reserved',
        status_label: 'ติดจอง',
        days_left_estimate: 0
      } : (data.current_booker || null)

      stallWaitingQueueCache.value[stallCode] = {
        stall_number: data.stall_number || data.stall_ident || stallCode,
        total_waiting_count: Number(data.total_waiting_count ?? data.queue_count ?? 0),
        next_queue_number: Number(data.next_queue_number ?? data.current_position ?? 1),
        has_waiting_list: Boolean(data.has_waiting_list ?? ((data.queue_count ?? 0) > 0)),
        current_booker: currentBooker,
        queues: data.queues || []
      }
    }
  } catch (err) {
    // Do not invent a queue owner or position when the source cannot be read.
    delete stallWaitingQueueCache.value[stallCode]
  } finally {
    loadingWaitingQueue.value = false
  }
}

// Check if a stall is selected
const isStallSelected = (code: string) => {
  return selectedStalls.value.some(s => s.code === code)
}

// Multi-stall aggregate totals
const totalMonthlyRent = computed(() => {
  if (selectedStalls.value.length === 0) return 1800
  return selectedStalls.value.reduce((sum, s) => sum + (s.monthlyRent || 1800), 0)
})

const totalEntryFee = computed(() => {
  if (selectedStalls.value.length === 0) return 75000
  return selectedStalls.value.reduce((sum, s) => sum + (s.entryFee || 75000), 0)
})

const totalDepositFee = computed(() => {
  if (selectedStalls.value.length === 0) return 1800
  return selectedStalls.value.reduce((sum, s) => sum + (s.depositFee || 1800), 0)
})

const totalCommonFee = computed(() => {
  if (selectedStalls.value.length === 0) return 500
  return selectedStalls.value.reduce((sum, s) => sum + (s.commonFee || 500), 0)
})

const totalSize = computed(() => {
  if (selectedStalls.value.length === 0) return 5
  return selectedStalls.value.reduce((sum, s) => sum + (s.size || 5), 0)
})

// Commercial financial calculations (BA & Dev perspective)
const quarterlyRent = computed(() => {
  return Math.round(totalMonthlyRent.value * 3 * 0.95) // 5% discount
})

const quarterlySavings = computed(() => {
  return (totalMonthlyRent.value * 3) - quarterlyRent.value
})

const annualRent = computed(() => {
  return Math.round(totalMonthlyRent.value * 12 * 0.9) // 10% discount
})

const annualSavings = computed(() => {
  return ((totalMonthlyRent.value * 12) - annualRent.value) + totalCommonFee.value // 10% discount + 1 month common fee waived
})

const contractTotal = computed(() => {
  return totalEntryFee.value + totalMonthlyRent.value + totalDepositFee.value + totalCommonFee.value
})

// Dynamic Breadcrumb based on View Mode
const bookingBreadcrumbItems = computed(() => {
  if (currentStep.value === 1 && selectedProject.value) {
    return [{ label: 'จองแผงเช่า', onClick: () => { selectedProject.value = null } }]
  }
  if (currentStep.value === 2) {
    return [{ label: 'จองแผงเช่า', onClick: () => { currentStep.value = 1 } }]
  }
  if (currentStep.value === 3) {
    return [
      { label: 'จองแผงเช่า', onClick: () => { currentStep.value = 1 } },
      { label: 'กรอกข้อมูลเพื่อเยี่ยมชมแผง', onClick: () => { currentStep.value = 2 } }
    ]
  }
  return []
})

const currentBreadcrumbText = computed(() => {
  if (currentStep.value === 1) {
    return selectedProject.value ? selectedProject.value.name : 'จองแผงเช่า'
  }
  if (currentStep.value === 2) return 'กรอกข้อมูลเพื่อเยี่ยมชมแผง'
  return 'ตรวจสอบและยืนยัน'
})

const handleBreadcrumbBack = computed(() => {
  if (currentStep.value === 3) return () => { currentStep.value = 2 }
  if (currentStep.value === 2) return () => { currentStep.value = 1 }
  if (currentStep.value === 1 && selectedProject.value) return () => { selectedProject.value = null }
  return undefined
})

// Input Formatters (QA & UX)
function formatThaiId(val: string) {
  const raw = val.replace(/[^0-9]/g, '').slice(0, 13)
  if (raw.length <= 1) return raw
  if (raw.length <= 5) return `${raw.slice(0, 1)}-${raw.slice(1)}`
  if (raw.length <= 10) return `${raw.slice(0, 1)}-${raw.slice(1, 5)}-${raw.slice(5)}`
  if (raw.length <= 12) return `${raw.slice(0, 1)}-${raw.slice(1, 5)}-${raw.slice(5, 10)}-${raw.slice(10)}`
  return `${raw.slice(0, 1)}-${raw.slice(1, 5)}-${raw.slice(5, 10)}-${raw.slice(10, 12)}-${raw.slice(12, 13)}`
}

const handleIdCardInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  bookingForm.value.idCard = formatThaiId(target.value)
}

const handlePhoneInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  const raw = target.value.replace(/[^0-9]/g, '').slice(0, 10)
  if (raw.length <= 3) {
    bookingForm.value.phone = raw
  } else if (raw.length <= 6) {
    bookingForm.value.phone = `${raw.slice(0, 3)}-${raw.slice(3)}`
  } else {
    bookingForm.value.phone = `${raw.slice(0, 3)}-${raw.slice(3, 6)}-${raw.slice(6)}`
  }
}

const generateStallsForProject = (project: any, realStall20a?: any) => {
  const upper: StallModel[] = []
  const lower: StallModel[] = []

  const isProj23 = project.id === 'proj-23' || project.name.includes('23')
  const numOnly = project.name.replace(/[^0-9]/g, '') || '1'
  const prefix = isProj23 ? '' : `P${numOnly}-`

  for (let r = 1; r <= 5; r++) {
    for (let c = 1; c <= 14; c++) {
      let code = ''
      if (isProj23) {
        code = r === 2 && c === 12 ? '20A' : `U-${r}${String.fromCharCode(64 + c)}`
      } else {
        code = r === 2 && c === 12 ? `${prefix}20A` : `${prefix}U${r}${String.fromCharCode(64 + c)}`
      }

      const backendStatus = getBackendRecordStatus(backendStallsByIdent.value[code])
      const status: StallModel['status'] = isMarketFull.value
        ? 'occupied'
        : (backendStatus ? normalizeBackendStallStatus(backendStatus) : 'available')

      upper.push({
        code,
        zone: project.name,
        projectId: project.id,
        status,
        areaType: code.includes('20A') || code.includes('01') ? 'ติดเวที' : 'ทั่วไป',
        size: 5,
        entryFee: 75000,
        monthlyRent: 1800,
        commonFee: 500,
        depositFee: 1800,
        photo: '/images/stall_3d_perspective.jpg',
        backendId: isProj23 && code === '20A' && realStall20a ? realStall20a.id : undefined
      })
    }
  }

  for (let r = 1; r <= 4; r++) {
    for (let c = 1; c <= 14; c++) {
      let code = ''
      if (isProj23) {
        code = `L-${r}${String.fromCharCode(64 + c)}`
      } else {
        code = `${prefix}L${r}${String.fromCharCode(64 + c)}`
      }

      const backendStatus = getBackendRecordStatus(backendStallsByIdent.value[code])
      const status: StallModel['status'] = isMarketFull.value
        ? 'occupied'
        : (backendStatus ? normalizeBackendStallStatus(backendStatus) : 'available')

      lower.push({
        code,
        zone: project.name,
        projectId: project.id,
        status,
        areaType: 'ทั่วไป',
        size: 5,
        entryFee: 75000,
        monthlyRent: 1800,
        commonFee: 500,
        depositFee: 1800,
        photo: '/images/stall_3d_perspective.jpg'
      })
    }
  }

  upperHallStalls.value = upper
  lowerHallStalls.value = lower

  // Do NOT auto-select stall! Let the tenant choose themselves.
}

const handleSelectProject = async (project: any) => {
  if (selectedProject.value?.id === project.id) return
  selectedProject.value = project
  if (selectedStalls.value.length === 0) {
    bookedProject.value = project
  }
  generateStallsForProject(project)
}

const handleExitProject = () => {
  selectedProject.value = null
}

// Stall clicked inside a project polygon on the GIS map
const handleMapStallSelect = (stall: StallModel) => {
  const backendStall = backendStallsByIdent.value[stall.code]
  const effectiveStatus = backendStall ? normalizeBackendStallStatus(getBackendRecordStatus(backendStall), backendStall) : stall.status
  const resolvedProjectId = stall.projectId || selectedProject.value?.id || bookedProject.value?.id
  const resolvedZone = stall.zone || backendStall?.project_name || selectedProject.value?.name || bookedProject.value?.name || ''
  selectStall({ ...stall, zone: resolvedZone, projectId: resolvedProjectId, status: effectiveStatus, backendId: backendStall?.id })
}

const handleNextFromStep1 = () => {
  if (isMarketFull.value) {
    showWarning(`ตลาด${currentMarketName.value} มีผู้เช่าเต็มแล้ว ขณะนี้ยังไม่สามารถทำรายการจองได้`, 'ตลาดเต็มแล้ว')
    return
  }
  if (selectedStalls.value.length === 0) {
    showWarning('กรุณาคลิกเลือกโครงการบนแผนที่ แล้วเลือกแผงย่อยที่ต้องการก่อนไปขั้นตอนถัดไป', 'แจ้งเตือนการเลือกแผง')
    return
  }

  if (currentSelectionStatus.value === 'reserved') {
    void joinSelectedWaitlists()
    return
  }
  const duplicateStall = selectedStalls.value.find(s => Boolean(myActiveBookingsByStall.value[s.code]))
  if (duplicateStall) {
    const existingBooking = myActiveBookingsByStall.value[duplicateStall.code]
    selectedStalls.value = selectedStalls.value.filter(s => !myActiveBookingsByStall.value[s.code])
    showDuplicateBookingWarning(duplicateStall.code, existingBooking)
    return
  }
  goToStep(2)
}

// Drawing logic
const getCanvasPos = (clientX: number, clientY: number) => {
  if (!signatureCanvas.value) return { x: 0, y: 0 }
  const rect = signatureCanvas.value.getBoundingClientRect()
  const scaleX = signatureCanvas.value.width / rect.width
  const scaleY = signatureCanvas.value.height / rect.height
  return {
    x: (clientX - rect.left) * scaleX,
    y: (clientY - rect.top) * scaleY
  }
}

const startDrawing = (e: MouseEvent) => {
  if (!signatureCanvas.value) return
  sigCtx = signatureCanvas.value.getContext('2d')
  if (!sigCtx) return
  isDrawing = true
  const pos = getCanvasPos(e.clientX, e.clientY)
  sigCtx.beginPath()
  sigCtx.moveTo(pos.x, pos.y)
  sigCtx.strokeStyle = '#1e3a8a'
  sigCtx.lineWidth = 2.5
  sigCtx.lineCap = 'round'
  sigCtx.lineJoin = 'round'
}

const draw = (e: MouseEvent) => {
  if (!isDrawing || !sigCtx) return
  const pos = getCanvasPos(e.clientX, e.clientY)
  sigCtx.lineTo(pos.x, pos.y)
  sigCtx.stroke()
  signatureDrawn.value = true
}

const stopDrawing = () => {
  isDrawing = false
}

const handleTouchStart = (e: TouchEvent) => {
  if (e.touches.length > 0) {
    const t = e.touches[0]
    if (!signatureCanvas.value) return
    sigCtx = signatureCanvas.value.getContext('2d')
    if (!sigCtx) return
    isDrawing = true
    const pos = getCanvasPos(t.clientX, t.clientY)
    sigCtx.beginPath()
    sigCtx.moveTo(pos.x, pos.y)
    sigCtx.strokeStyle = '#1e3a8a'
    sigCtx.lineWidth = 2.5
    sigCtx.lineCap = 'round'
    sigCtx.lineJoin = 'round'
  }
}

const handleTouchMove = (e: TouchEvent) => {
  if (e.touches.length > 0 && isDrawing && sigCtx) {
    const t = e.touches[0]
    const pos = getCanvasPos(t.clientX, t.clientY)
    sigCtx.lineTo(pos.x, pos.y)
    sigCtx.stroke()
    signatureDrawn.value = true
  }
}

const clearSignature = () => {
  if (!signatureCanvas.value) return
  const ctx = signatureCanvas.value.getContext('2d')
  if (ctx) {
    ctx.clearRect(0, 0, signatureCanvas.value.width, signatureCanvas.value.height)
  }
  signatureDrawn.value = false
}

const autoSign = () => {
  if (!signatureCanvas.value) return
  const ctx = signatureCanvas.value.getContext('2d')
  if (!ctx) return
  clearSignature()
  ctx.strokeStyle = '#1e3a8a'
  ctx.lineWidth = 2.5
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'

  // Draw cursive signature
  ctx.beginPath()
  ctx.moveTo(80, 70)
  ctx.bezierCurveTo(95, 30, 110, 30, 125, 65)
  ctx.bezierCurveTo(140, 90, 155, 35, 175, 55)
  ctx.bezierCurveTo(195, 75, 210, 45, 230, 60)
  ctx.bezierCurveTo(250, 75, 270, 35, 290, 50)
  ctx.bezierCurveTo(310, 65, 340, 40, 370, 60)
  ctx.moveTo(90, 85)
  ctx.bezierCurveTo(180, 80, 270, 95, 380, 78)
  ctx.stroke()

  signatureDrawn.value = true
}

// Confetti Engine
const confettiCanvas = ref<HTMLCanvasElement | null>(null)
const launchConfetti = () => {
  nextTick(() => {
    if (!confettiCanvas.value) return
    const canvas = confettiCanvas.value
    canvas.width = canvas.parentElement?.clientWidth || 560
    canvas.height = canvas.parentElement?.clientHeight || 480
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const particles: Array<{
      x: number
      y: number
      vx: number
      vy: number
      size: number
      color: string
      rotation: number
      vRot: number
    }> = []

    const colors = ['#3b82f6', '#10b981', '#f59e0b', '#ec4899', '#8b5cf6', '#06b6d4']
    for (let i = 0; i < 70; i++) {
      particles.push({
        x: canvas.width / 2 + (Math.random() - 0.5) * 60,
        y: canvas.height / 3 + (Math.random() - 0.5) * 40,
        vx: (Math.random() - 0.5) * 12,
        vy: Math.random() * -10 - 2,
        size: Math.random() * 8 + 4,
        color: colors[Math.floor(Math.random() * colors.length)],
        rotation: Math.random() * 360,
        vRot: (Math.random() - 0.5) * 10
      })
    }

    let frames = 0
    const render = () => {
      if (!ctx) return
      ctx.clearRect(0, 0, canvas.width, canvas.height)
      particles.forEach(p => {
        p.x += p.vx
        p.y += p.vy
        p.vy += 0.35
        p.vx *= 0.98
        p.rotation += p.vRot

        ctx.save()
        ctx.translate(p.x, p.y)
        ctx.rotate((p.rotation * Math.PI) / 180)
        ctx.fillStyle = p.color
        ctx.fillRect(-p.size / 2, -p.size / 2, p.size, p.size * 0.6)
        ctx.restore()
      })
      frames++
      if (frames < 120) {
        requestAnimationFrame(render)
      } else {
        ctx.clearRect(0, 0, canvas.width, canvas.height)
      }
    }
    render()
  })
}


// Pass Actions (QA & UX)
const copyRefCode = async () => {
  try {
    if (navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(bookingRefCode.value)
    }
    copiedRef.value = true
    setTimeout(() => { copiedRef.value = false }, 2500)
  } catch {}
}

const printPass = () => {
  window.print()
}

const addToCalendar = () => {
  const stallCodes = passDisplayStalls.value.map(s => s.code).join(', ') || '20A'
  const title = `นัดหมายเยี่ยมชมแผง ${stallCodes} - ${currentMarketName.value}`
  const details = `นัดหมายเข้าชมพื้นที่แผงเช่ากับเจ้าหน้าที่ รหัสอ้างอิง: ${bookingRefCode.value}`
  const dateStr = (bookingForm.value.visitDate || '').replace(/-/g, '')
  const icsData = `BEGIN:VCALENDAR\nVERSION:2.0\nBEGIN:VEVENT\nSUMMARY:${title}\nDESCRIPTION:${details}\nDTSTART:${dateStr}T130000\nDTEND:${dateStr}T140000\nEND:VEVENT\nEND:VCALENDAR`
  const blob = new Blob([icsData], { type: 'text/calendar;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.setAttribute('download', `visit-stall-${bookingRefCode.value}.ics`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const closeSuccessModal = () => {
  showSuccessModal.value = false
  if (selectedStalls.value.length === 0) {
    currentStep.value = 1
  }
}

const selectStall = (s: StallModel) => {
  if (isMarketFull.value) {
    showWarning(`ตลาด${currentMarketName.value} มีผู้เช่าเต็มแล้ว ขณะนี้ยังไม่สามารถเลือกหรือจองแผงได้`, 'ตลาดนี้เต็มแล้ว')
    return
  }
  if (s.status === 'occupied') {
    return
  }
  // Toggle off if clicking an already selected stall
  const selectedKey = stallSelectionKey(s)
  const existingIdx = selectedStalls.value.findIndex(item => stallSelectionKey(item) === selectedKey)
  if (existingIdx !== -1) {
    selectedStalls.value.splice(existingIdx, 1)
    if (selectedStalls.value.length === 0) {
      bookedProject.value = null
    }
    if (activeWaitingStallCode.value === s.code) {
      activeWaitingStallCode.value = selectedStalls.value[0]?.code || ''
    }
    return
  }

  // Duplicate Booking Prevention: block selecting a stall the current tenant already booked or queued for
  const existingActiveBooking = myActiveBookingsByStall.value[s.code]
  if (existingActiveBooking) {
    showDuplicateBookingWarning(s.code, existingActiveBooking)
    return
  }

  const existingWaitlist = myWaitlistByStall.value[s.code]
  if (existingWaitlist && existingWaitlist.status === 'waiting') {
    // Selecting an already-waitlisted stall is a read-only preview action.
    // Keep duplicate queue submission disabled in the drawer CTA instead of
    // interrupting the user with a modal.
    if (!s.projectId && selectedProject.value?.id) {
      s.projectId = selectedProject.value.id
    }
    if (selectedProject.value) {
      bookedProject.value = selectedProject.value
    }
    selectedStalls.value = [{ ...s, status: 'reserved' }]
    activeWaitingStallCode.value = s.code
    void fetchWaitingQueueForStall(s.code)
    return
  }

  // Rule 1 & 2: แผงติดจองไม่สามารถเพิ่มเข้าใบจองได้ และ Waiting list ผูกกับแผงรายแผง ไม่ใช่ตะกร้าหลายแผง
  if (currentSelectionStatus.value && s.status !== currentSelectionStatus.value) {
    if (s.status === 'reserved') {
      showWarning('แผงติดจองไม่สามารถเพิ่มเข้าใบจองได้ กรุณาดำเนินการจองแผงว่างที่เลือกไว้ให้เสร็จสิ้นหรือล้างรายการเดิมก่อน', 'ไม่สามารถเพิ่มแผงติดจองเข้าใบจองได้')
    } else {
      showWarning('ไม่สามารถเลือกแผงว่างปนกับแผงที่กำลังจะเข้าคิวรอได้ กรุณากดรอคิวแผงนี้หรือนำแผงติดจองออกก่อน', 'สถานะแผงไม่ตรงกัน')
    }
    return
  }
  if (!s.projectId && selectedProject.value?.id) {
    s.projectId = selectedProject.value.id
  }
  if (selectedProject.value) {
    bookedProject.value = selectedProject.value
  }
  if (s.status === 'reserved') {
    // Waiting list ผูกกับแผงรายแผง ไม่ใช่ตะกร้าหลายแผง
    selectedStalls.value = [s]
    activeWaitingStallCode.value = s.code
    fetchWaitingQueueForStall(s.code)
  } else {
    selectedStalls.value.push(s)
    activeWaitingStallCode.value = ''
  }
}

const joinSelectedWaitlists = async () => {
  if (selectedStalls.value.length === 0 || currentSelectionStatus.value !== 'reserved') return
  const stall = selectedStalls.value[0]!
  const fullName = `${bookingForm.value.firstName || ''} ${bookingForm.value.lastName || ''}`.trim() || 'ผู้เช่า'
  const phone = String(bookingForm.value.phone || '').trim()
  try {
    const rawRes: any = await post('/bookings/waitlist', {
      stall_id: stall.backendId,
      stall_ident: stall.code,
      market_id: Number(currentMarket.value?.id) || undefined,
      market_name: currentMarketName.value,
      project_name: stall.zone || selectedProject.value?.name || '',
      tenant_name: fullName,
      tenant_phone: phone
    })
    const created = rawRes?.data || rawRes
    const pos = Number(created?.position) || 1
    selectedStalls.value = []
    bookedProject.value = null
    activeWaitingStallCode.value = ''
    delete stallWaitingQueueCache.value[stall.code]
    await Promise.all([refreshBackendStalls(), refreshMyActiveBookings(), fetchNotifications()])
    await showSuccess(
      `ลงชื่อรอคิวแผง ${stall.code} สำเร็จ (ลำดับคิวที่ ${pos}) ระบบจะแจ้งเตือนทันทีเมื่อแผงว่างลง คุณสามารถดูสถานะและออกจากคิวได้ที่แท็บ "แผงที่สนใจ"`,
      'รอคิวแผงนี้สำเร็จ'
    )
  } catch (error: any) {
    const message = error?.data?.detail || error?.data?.error?.message || error?.data?.message || 'ไม่สามารถเข้าคิวรอแผงนี้ได้'
    await showWarning(message, 'ไม่สามารถเข้าคิวได้')
  }
}

const removeStall = (stall: StallModel) => {
  const selectedKey = stallSelectionKey(stall)
  selectedStalls.value = selectedStalls.value.filter(item => stallSelectionKey(item) !== selectedKey)
  if (selectedStalls.value.length === 0) {
    bookedProject.value = null
  }
  if (activeWaitingStallCode.value === stall.code) {
    activeWaitingStallCode.value = selectedStalls.value[0]?.code || ''
  }
}

const clearAllSelectedStalls = () => {
  selectedStalls.value = []
  bookedProject.value = null
  activeWaitingStallCode.value = ''
}

const resolveConflictAndSelect = () => {
  if (pendingConflictingStall.value) {
    const stall = pendingConflictingStall.value
    const existingActiveBooking = myActiveBookingsByStall.value[stall.code]
    if (existingActiveBooking) {
      pendingConflictingStall.value = null
      showConflictModal.value = false
      showDuplicateBookingWarning(stall.code, existingActiveBooking)
      return
    }
    if (!stall.projectId && selectedProject.value?.id) {
      stall.projectId = selectedProject.value.id
    }
    if (selectedProject.value) {
      bookedProject.value = selectedProject.value
    }
    selectedStalls.value = [stall]
    if (stall.status === 'reserved') {
      activeWaitingStallCode.value = stall.code
      fetchWaitingQueueForStall(stall.code)
    }
    pendingConflictingStall.value = null
  }
  showConflictModal.value = false
}

const cancelConflictModal = () => {
  pendingConflictingStall.value = null
  showConflictModal.value = false
}

const getStallClass = (s: StallModel) => {
  const isSelected = isStallSelected(s.code)
  const isConflict = currentSelectionStatus.value && s.status !== currentSelectionStatus.value && s.status !== 'occupied'
  let isDimmed = onlyAvailable.value && s.status !== 'available'
  if (zoneFilter.value === 'hotspot') {
    const isHotspot = s.code === '20A' || s.code === '20B' || s.code.endsWith('A') || s.code.endsWith('N')
    if (!isHotspot) isDimmed = true
  } else if (zoneFilter.value === 'stage') {
    const isStage = s.code === '20A' || s.code === '20B' || s.code.startsWith('U-2')
    if (!isStage) isDimmed = true
  }
  return {
    'is-available': !isMarketFull.value && s.status === 'available',
    'is-reserved': !isMarketFull.value && s.status === 'reserved',
    'is-occupied': isMarketFull.value || s.status === 'occupied',
    'is-market-full': isMarketFull.value,
    'is-selected': isSelected,
    'is-dimmed': isDimmed,
    'is-conflict-mode': isConflict,
    'traffic-hotspot': showTrafficHeatmap.value && (s.code === '20A' || s.code === '20B' || s.areaType === 'ติดเวที')
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'available': return 'ว่าง'
    case 'reserved': return 'ติดจอง'
    default: return 'กำลังใช้งาน'
  }
}

// Pan & Zoom Engine for Stalls Floorplan
const stallsPanX = ref(0)
const stallsPanY = ref(0)
const isStallsDragging = ref(false)
const hasStallsDragged = ref(false)
let stallsDragStart = { x: 0, y: 0 }

const zoomIn = () => {
  if (zoomLevel.value < 2.4) {
    zoomLevel.value = Math.min(2.4, +(zoomLevel.value + 0.2).toFixed(2))
  }
}

const zoomOut = () => {
  if (zoomLevel.value > 0.6) {
    zoomLevel.value = Math.max(0.6, +(zoomLevel.value - 0.2).toFixed(2))
  }
}

const recenterMap = () => {
  zoomLevel.value = 1
  stallsPanX.value = 0
  stallsPanY.value = 0
}

const handleStallsWheel = (e: WheelEvent) => {
  const factor = e.deltaY < 0 ? 1.15 : 0.88
  zoomLevel.value = Math.min(2.4, Math.max(0.6, +(zoomLevel.value * factor).toFixed(2)))
}

const onStallsMouseDown = (e: MouseEvent) => {
  if (e.button !== 0) return
  isStallsDragging.value = true
  hasStallsDragged.value = false
  stallsDragStart = { x: e.clientX - stallsPanX.value, y: e.clientY - stallsPanY.value }
}

const onStallsMouseMove = (e: MouseEvent) => {
  if (!isStallsDragging.value) return
  const dx = e.clientX - stallsDragStart.x
  const dy = e.clientY - stallsDragStart.y
  if (Math.abs(dx - stallsPanX.value) > 4 || Math.abs(dy - stallsPanY.value) > 4) {
    hasStallsDragged.value = true
  }
  stallsPanX.value = dx
  stallsPanY.value = dy
}

const onStallsMouseUp = () => {
  isStallsDragging.value = false
}

const validateStep2 = (): boolean => {
  let isValid = true
  formErrors.value = {
    business: '',
    visitDate: '',
    visitSlot: ''
  }

  if (!bookingForm.value.business) {
    formErrors.value.business = 'กรุณาเลือกสินค้าที่ต้องการขาย'
    isValid = false
  }

  if (!bookingForm.value.visitDate) {
    formErrors.value.visitDate = 'กรุณาเลือกวันที่เปิดให้เข้าชมแผง'
    isValid = false
  }

  if (!bookingForm.value.visitSlot) {
    formErrors.value.visitSlot = 'กรุณาเลือกรอบเวลาเข้าชมแผง'
    isValid = false
  } else {
    const currentSlotObj = currentTimeSlots.value.find(s => s.time === bookingForm.value.visitSlot)
    if (currentSlotObj?.hasTenantConflict) {
      const conflictMarket = currentSlotObj.tenantConflictInfo?.marketName || 'ตลาดอื่น'
      const conflictProject = currentSlotObj.tenantConflictInfo?.projectName || ''
      formErrors.value.visitSlot = `คุณมีนัดหมายเข้าชมแผงที่ ${conflictMarket}${conflictProject ? ` (${conflictProject})` : ''} ในวันและเวลานี้แล้ว กรุณาเลือกรอบเวลาอื่น`
      isValid = false
    } else if (currentSlotObj?.isFull) {
      formErrors.value.visitSlot = `รอบเวลา ${bookingForm.value.visitSlot} เต็มแล้ว กรุณาเลือกรอบเวลาอื่น`
      isValid = false
    }
  }

  if (!isValid) {
    nextTick(() => {
      const firstErrorEl = document.querySelector('.has-form-error')
      if (firstErrorEl) {
        firstErrorEl.scrollIntoView({ behavior: 'smooth', block: 'center' })
      }
    })
  }

  return isValid
}

const goToStep = (step: number) => {
  if (step === currentStep.value) return
  if (step > 1 && isMarketFull.value) {
    showWarning(`ตลาด${currentMarketName.value} มีผู้เช่าเต็มแล้ว ขณะนี้ยังไม่สามารถทำรายการจองได้`, 'ตลาดเต็มแล้ว')
    return
  }
  if (step > 1 && selectedStalls.value.length === 0) {
    showWarning('กรุณาเลือกแผงบนแผนที่ก่อนถึงจะไปขั้นตอนถัดไปได้', 'แจ้งเตือนการเลือกแผง')
    return
  }
  if (step > 1) {
    const duplicateStall = selectedStalls.value.find(s => Boolean(myActiveBookingsByStall.value[s.code]))
    if (duplicateStall) {
      const existingBooking = myActiveBookingsByStall.value[duplicateStall.code]
      selectedStalls.value = selectedStalls.value.filter(s => !myActiveBookingsByStall.value[s.code])
      showDuplicateBookingWarning(duplicateStall.code, existingBooking)
      return
    }
  }
  if (step === 3) {
    if (!validateStep2()) {
      return
    }
  }
  currentStep.value = step
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

watch(currentStep, (newStep) => {
  if (newStep === 2) {
    void refreshMyActiveBookings()
  }
  if (newStep === 1 || newStep === 3) {
    nextTick(() => {
      window.dispatchEvent(new Event('resize'))
      if (newStep === 1 && selectedStalls.value.length > 0) {
        focusClaimedStallOnMap()
      }
    })
  }
})

const setStep = (step: number) => {
  goToStep(step)
}

const formatThaiDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear() + 543}`
}


const submitBooking = async () => {
  if (selectedStalls.value.length === 0) {
    showWarning('ไม่พบรายการแผงที่เลือก กรุณาเลือกแผงเช่าก่อนดำเนินการ', 'แจ้งเตือนการเลือกแผง')
    return
  }
  submitting.value = true
  try {
    await Promise.all([refreshBackendStalls(), refreshMyActiveBookings(), fetchNotifications()])

    // Prevent submitting a duplicate booking for a stall the current tenant already holds or queues for
    const duplicateStall = selectedStalls.value.find(stall => Boolean(myActiveBookingsByStall.value[stall.code]))
    if (duplicateStall) {
      const existingBooking = myActiveBookingsByStall.value[duplicateStall.code]
      selectedStalls.value = selectedStalls.value.filter(stall => !myActiveBookingsByStall.value[stall.code])
      if (selectedStalls.value.length === 0) {
        bookedProject.value = null
      }
      showDuplicateBookingWarning(duplicateStall.code, existingBooking)
      return
    }

    const refreshedSelection = selectedStalls.value.map(stall => {
      const backendStall = backendStallsByIdent.value[stall.code]
      return backendStall
        ? { ...stall, status: normalizeBackendStallStatus(getBackendRecordStatus(backendStall), backendStall), backendId: backendStall.id }
        : stall
    })
    const statusChanged = refreshedSelection.some((stall, index) => stall.status !== selectedStalls.value[index]?.status)
    const firstSelectableStatus = refreshedSelection.find(stall => stall.status !== 'occupied')?.status
    selectedStalls.value = refreshedSelection.filter(stall => (
      stall.status !== 'occupied' && stall.status === firstSelectableStatus
    ))
    if (statusChanged) {
      showWarning(
        selectedStalls.value.length > 0
          ? 'สถานะของแผงที่เลือกมีการเปลี่ยนแปลงโดยเจ้าหน้าที่ กรุณาตรวจสอบรายการอีกครั้งก่อนยืนยัน'
          : 'แผงที่เลือกถูกล็อกหรือมีผู้เช่าแล้ว รายการดังกล่าวถูกนำออกจากการจอง',
        'สถานะแผงมีการเปลี่ยนแปลง'
      )
      return
    }

    const appointmentGroupId = `BKG-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
    const resolvedMarketNumericId = Number(
      currentMarket.value?.id ||
      getMarketGisConfig(currentMarketId.value)?.id ||
      currentMarketId.value
    )
    const numericMarketId = Number.isInteger(resolvedMarketNumericId) && resolvedMarketNumericId > 0
      ? resolvedMarketNumericId
      : 1
    const fullName = `${bookingForm.value.firstName || 'สมชาย'} ${bookingForm.value.lastName || 'ใจดี'}`.trim()
    const contactPhone = (bookingForm.value.phone || '081-234-5678').trim()
    const contactEmail = (bookingForm.value.email || '').trim()
    const idCard = (bookingForm.value.idCard || '1-1004-00234-56-7').trim()
    const visitDate = bookingForm.value.visitDate || new Date().toISOString().split('T')[0]
    const visitTime = bookingForm.value.visitSlot || '13:00 - 14:00'
    const firstStall = selectedStalls.value[0]!

    const stallsPayload = selectedStalls.value.map((stallItem) => {
      const backendStall = backendStallsByIdent.value[stallItem.code]
      const stallId = stallItem.backendId || backendStall?.id
      const numericStallId = typeof stallId === 'number' && Number.isInteger(stallId)
        ? stallId
        : (typeof stallId === 'string' && /^\d+$/.test(stallId) ? Number(stallId) : undefined)
      return {
        stall_id: numericStallId,
        stall_ident: stallItem.code,
        stall_number: stallItem.code,
        project_name: stallItem.zone || backendStall?.project_name || selectedProject.value?.name || 'โครงการ 01',
        monthly_rent: Number(stallItem.monthlyRent || 1800),
        deposit_fee: Number(stallItem.depositFee || stallItem.monthlyRent || 1800)
      }
    })
    const uniqueProjectNames = Array.from(new Set(stallsPayload.map(s => s.project_name).filter(Boolean)))

    const res: any = await post('/bookings', {
      stalls: stallsPayload,
      stall_ident: firstStall.code,
      stall_number: firstStall.code,
      applicant_name: fullName,
      applicant_phone: contactPhone,
      applicant_email: contactEmail,
      business_type: bookingForm.value.business || 'อาหารและเครื่องดื่ม',
      visit_date: visitDate,
      visit_time: visitTime,
      note: bookingForm.value.notes || '',
      market_id: numericMarketId,
      project_name: uniqueProjectNames.join(', ') || firstStall.zone || selectedProject.value?.name || 'โครงการ 01',
      market_name: currentMarketName.value || 'ตลาดนัดจตุจักร',
      tenant_name: fullName,
      tenant_user_id: (user.value as any)?.id || undefined,
      tenant_phone: contactPhone,
      tenant_email: contactEmail || null,
      tenant_id_card: idCard,
      intended_business: bookingForm.value.business || 'อาหารและเครื่องดื่ม',
      visit_time_slot: visitTime,
      is_queue: false,
      appointment_group_id: appointmentGroupId
    })

    const createdRecord = res?.data || res
    const bNo = createdRecord?.booking_no || createdRecord?.booking_number || res?.booking_no || res?.booking_number
    const succeededStalls = [...selectedStalls.value]
    succeededStalls.forEach((stallItem) => {
      sessionBookedStallMap.value[stallItem.code] = {
        ...(createdRecord || {}),
        stall_number: stallItem.code,
        booking_number: bNo,
        holding_stall: true,
        queue_position: 0,
        status: createdRecord?.status || 'pending_approval'
      }
      delete stallWaitingQueueCache.value[stallItem.code]
    })

    await Promise.all([refreshBackendStalls(), refreshMyActiveBookings(), fetchNotifications()])

    confirmedPassStalls.value = [...succeededStalls]
    selectedStalls.value = []
    bookedProject.value = null
    activeWaitingStallCode.value = ''

    const randomDigits = Math.floor(10000 + Math.random() * 90000)
    bookingRefCode.value = bNo || `BMA-JJ-2026-${randomDigits}`
    showSuccessModal.value = true
    launchConfetti()
  } catch (err: any) {
    showError(err?.data?.detail || err?.data?.error?.message || err?.data?.message || 'เกิดข้อผิดพลาดในการบันทึกคำขอจอง กรุณาลองใหม่อีกครั้ง')
  } finally {
    submitting.value = false
  }
}

// Keyboard shortcuts (QA & Dev)
const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    if (showSuccessModal.value) closeSuccessModal()
    showVirtualModal.value = false
    showCompareModal.value = false
  }
  if (showVirtualModal.value) {
    if (e.key === 'ArrowLeft') stepPreview(-1)
    if (e.key === 'ArrowRight') stepPreview(1)
  }
}

onMounted(async () => {
  window.addEventListener('keydown', handleKeyDown)
  window.addEventListener('focus', handleBackendStallSync)

  // Listen for schedule updates from VisitScheduleMain in same window or across tabs
  const handleScheduleSync = () => {
    availableDates.value = loadAvailableDatesFromSchedule()
  }
  window.addEventListener('market-visit-schedule-updated', handleScheduleSync)
  const handleStorageSync = (e: StorageEvent) => {
    if (e.key?.startsWith('visit-schedule:')) {
      handleScheduleSync()
    }
  }
  window.addEventListener('storage', handleStorageSync)

  const syncMarketData = async (marketQuery: string | number) => {
    const queryStr = String(marketQuery || '1')
    const localMaster = getMarketGisConfig(queryStr)
    if (localMaster) {
      const overriddenLocal = applyMarketOverride({
        id: localMaster.id,
        slug: localMaster.slug,
        name: localMaster.name,
        district: localMaster.district
      })
      currentMarketName.value = overriddenLocal.name || localMaster.name
      currentMarketDistrict.value = overriddenLocal.district || localMaster.district
    }

    try {
      const marketsRes: any = await get('/markets')
      const marketsList = Array.isArray(marketsRes) ? marketsRes : (Array.isArray(marketsRes?.data) ? marketsRes.data : [])
      const found = marketsList.find((m: any) => String(m.id) === queryStr || m.slug === queryStr)
      if (found) {
        const resolved = applyMarketOverride(found)
        currentMarket.value = resolved
        currentMarketName.value = resolved.name
        currentMarketDistrict.value = resolved.district
        // Treat zero as authoritative only when the API actually supplies the
        // field. `null`/missing values must not accidentally close the market.
        const hasAvailableCount = found.available_stalls !== undefined && found.available_stalls !== null
        if (found.status === 'full' || (hasAvailableCount && Number(found.available_stalls) === 0)) {
          isMarketFull.value = true
          if (selectedProject.value) selectedProject.value.stallsAvailable = 0
        } else {
          isMarketFull.value = false
        }
      }
    } catch (e) {
      console.warn('Could not fetch market:', e)
    }
  }

  await syncMarketData(route.query.market || '1')
  watch(currentMarketId, async (newId) => {
    await syncMarketData(newId)
    await Promise.all([refreshBackendStalls(), refreshMyActiveBookings(), fetchNotifications()])
  })

  if (user.value) {
    const demoTenantProfiles: Record<string, { phone: string; idCard: string }> = {
      tenant: { phone: '081-234-5678', idCard: '1-1004-99887-12-3' },
      tenant2: { phone: '089-111-2233', idCard: '1-1005-44221-55-6' },
      tenant3: { phone: '086-778-9900', idCard: '1-1006-88772-33-4' }
    }
    const demoProfile = demoTenantProfiles[String((user.value as any).username || '')]
    if (user.value.full_name) {
      const cleanedFullName = String(user.value.full_name)
        .replace(/\([^)]*\)/g, '')
        .replace(/^(คุณ|นาย|นางสาว|นาง|น\.ส\.)\s*/g, '')
        .trim()
      const parts = cleanedFullName.split(/\s+/).filter(Boolean)
      bookingForm.value.firstName = parts[0] || 'สมชาย'
      bookingForm.value.lastName = parts.slice(1).join(' ') || 'ใจดี'
    }
    if (user.value.email) {
      bookingForm.value.email = user.value.email
    }
    if ((user.value as any).phone || (user.value as any).phone_number || demoProfile) {
      bookingForm.value.phone = (user.value as any).phone || (user.value as any).phone_number || demoProfile?.phone || bookingForm.value.phone
    }
    if ((user.value as any).citizen_id || (user.value as any).id_card || demoProfile) {
      bookingForm.value.idCard = formatThaiId((user.value as any).citizen_id || (user.value as any).id_card || demoProfile?.idCard || bookingForm.value.idCard)
    }
  }

  await Promise.all([refreshBackendStalls(), refreshMyActiveBookings()])
  const realStall20a = backendStallsByIdent.value['20A'] || null

  // Initialize stalls for selected project if any was preselected
  if (selectedProject.value) {
    generateStallsForProject(selectedProject.value, realStall20a)
  }
  selectedStalls.value = []

  // If a specific stall was requested via URL query (?stall=20A), select it (including when offered from waitlist!)
  if (route.query.stall && !isMarketFull.value) {
    const targetCode = String(route.query.stall).trim()
    const backendRec = backendStallsByIdent.value[targetCode]
    const found = upperHallStalls.value.find(s => s.code === targetCode) || lowerHallStalls.value.find(s => s.code === targetCode)
    const effStatus = backendRec
      ? normalizeBackendStallStatus(getBackendRecordStatus(backendRec), backendRec)
      : (found?.status || 'available')
    if (effStatus !== 'occupied' && !myActiveBookingsByStall.value[targetCode]) {
      const targetStallModel: StallModel = found
        ? { ...found, status: effStatus, backendId: backendRec?.id || found.backendId }
        : {
            code: targetCode,
            zone: String(route.query.project || backendRec?.project_name || 'โซนแผงค้า'),
            status: effStatus,
            areaType: 'ทั่วไป',
            size: 5,
            entryFee: 36000,
            monthlyRent: Number(backendRec?.monthly_rent || 1800),
            commonFee: 500,
            depositFee: Number(backendRec?.monthly_rent || 1800),
            photo: '/images/stall_3d_perspective.jpg',
            backendId: backendRec?.id
          }
      selectedStalls.value = [targetStallModel]
      if (backendRec?.market_id && String(backendRec.market_id) !== String(route.query.market || '')) {
        router.replace({
          query: {
            ...route.query,
            market: String(backendRec.market_id)
          }
        })
      }
      if (targetStallModel.status === 'reserved') {
        activeWaitingStallCode.value = targetCode
        fetchWaitingQueueForStall(targetCode)
      }
      focusClaimedStallOnMap()
      if (String(route.query.step || '') === '2' && effStatus === 'available') {
        currentStep.value = 2
      }
    }
  }

  // Prefetch demo waiting list queues
  fetchWaitingQueueForStall('A-03')
  fetchWaitingQueueForStall('B-04')
  fetchWaitingQueueForStall('02A')
  fetchWaitingQueueForStall('U-2J')
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
  window.removeEventListener('focus', handleBackendStallSync)
})
</script>

<style scoped>
.booking-stepper-wrap { margin-top: 40px; }

.step1-title-bar { display: flex; align-items: flex-end; justify-content: space-between; gap: 18px; margin-bottom: 16px; }
.step1-map-toolbar { flex: 0 1 700px; width: min(100%, 700px); margin-left: auto; display: flex; align-items: flex-end; justify-content: flex-end; }
.waitlist-claim-notice { display: flex; align-items: center; gap: 10px; padding: 11px 14px; background: #eff6ff; color: #1e3a8a; }
.waitlist-claim-notice.is-map-drawer { flex-shrink: 0; border-bottom: 1px solid #bfdbfe; }
.waitlist-claim-icon { width: 34px; height: 34px; flex: 0 0 34px; display: grid; place-items: center; border-radius: 9px; background: #dbeafe; color: #2563eb; }
.waitlist-claim-content { min-width: 0; display: flex; flex: 1; flex-direction: column; gap: 2px; }
.waitlist-claim-content strong { color: #1e3a8a; font-size: 0.9375rem; font-weight: 700; }
.waitlist-claim-content span { color: #475569; font-size: 0.8125rem; font-weight: 500; }
.waitlist-claim-title-row { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.waitlist-claim-badge { flex: 0 0 auto; padding: 3px 8px; border-radius: 999px; background: #dbeafe; color: #1d4ed8 !important; font-size: 0.75rem !important; font-weight: 700 !important; }
@media (max-width: 760px) { .step1-title-bar { align-items: flex-start; flex-direction: column; gap: 12px; } .step1-map-toolbar { width: 100%; max-width: none; } }

.map-selection-empty-drawer {
  position: absolute;
  top: 62px;
  right: 14px;
  z-index: 25;
  display: flex;
  width: 360px;
  min-height: 355px;
  max-height: calc(100% - 76px);
  flex-direction: column;
  overflow: hidden;
  /* Capture clicks so they do not pass through to a random map feature behind the drawer. */
  pointer-events: auto;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  background: rgba(255, 255, 255, .98);
  box-shadow: 0 12px 36px rgba(15, 23, 42, .16);
  backdrop-filter: blur(14px);
}
.msed-header { padding: 15px 16px 13px; border-bottom: 1px solid #eef2f6; }
.msed-header strong { display: block; color: #0f172a; font-size: 1.125rem; font-weight: 800; }
.msed-header span { display: block; margin-top: 3px; overflow: hidden; color: #64748b; font-size: .9375rem; font-weight: 600; text-overflow: ellipsis; white-space: nowrap; }
.msed-content { display: flex; flex: 1; flex-direction: column; align-items: center; padding: 30px 23px 22px; text-align: center; }
.msed-icon { display: grid; width: 58px; height: 58px; place-items: center; border: 1px solid #bfdbfe; border-radius: 17px; background: #eff6ff; color: #2563eb; box-shadow: 0 7px 16px rgba(37, 99, 235, .12); }
.msed-content h3 { margin: 15px 0 5px; color: #0f172a; font-size: 1.25rem; font-weight: 800; }
.msed-content > p { max-width: 290px; margin: 0; color: #64748b; font-size: .9375rem; line-height: 1.55; }
.msed-icon.is-full { border-color: #fecaca; background: #fef2f2; color: #dc2626; box-shadow: 0 7px 16px rgba(220, 38, 38, .1); }
.msed-status-list { width: 100%; margin-top: 21px; padding: 10px 12px; border: 1px solid #e2e8f0; border-radius: 10px; background: #f8fafc; text-align: left; }
.msed-status-list > div { display: flex; align-items: center; gap: 8px; color: #475569; font-size: .875rem; line-height: 1.5; }
.msed-status-list > div + div { margin-top: 8px; }
.msed-status-list i { width: 9px; height: 9px; flex: 0 0 auto; border-radius: 50%; }
.msed-status-list i.is-available { background: #10b981; }
.msed-status-list i.is-reserved { background: #f59e0b; }
.msed-status-list strong { color: #334155; font-weight: 800; }

/* Form input disabled style */
.form-input:disabled,
.form-select:disabled {
  background-color: #f8fafc !important;
  color: #475569 !important;
  border-color: #e2e8f0 !important;
  cursor: not-allowed !important;
  opacity: 0.92;
}

/* Form Field Inline Error Styles */
.form-field-error-text {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #dc2626;
  font-size: 0.8rem;
  font-weight: 500;
  margin-top: 6px;
  animation: formErrorFadeIn 0.2s ease-out;
}

/* Step 3 Summary Stall Photo Card */
.summary-stall-photo-card {
  position: relative;
  width: 140px;
  height: 105px;
  border-radius: 12px;
  overflow: hidden;
  border: 1.5px solid #e2e8f0;
  background: #f8fafc;
  cursor: pointer;
  transition: all 0.22s cubic-bezier(0.4, 0, 0.2, 1);
}

.summary-stall-photo-card:hover {
  transform: translateY(-3px) scale(1.02);
  border-color: #3b82f6;
  box-shadow: 0 8px 18px rgba(37, 99, 235, 0.18);
}

.photo-card-hover-overlay {
  position: absolute;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(2px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 600;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.summary-stall-photo-card:hover .photo-card-hover-overlay {
  opacity: 1;
}

.date-strip-error {
  outline: 2px dashed #f87171 !important;
  outline-offset: 4px;
  border-radius: 12px;
  background: #fff5f5;
  padding: 4px;
}

.time-slots-error {
  outline: 2px dashed #f87171 !important;
  outline-offset: 4px;
  border-radius: 12px;
  background: #fff5f5;
  padding: 4px;
}

@keyframes formErrorFadeIn {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Selected stalls bar under the GIS map */

/* Stall photo preview (lightbox) */
.photo-preview-overlay { background: rgba(15, 23, 42, 0.72); backdrop-filter: blur(4px); z-index: 1200; }
.photo-preview { width: min(860px, 94vw); background: #0f172a; border-radius: 18px; overflow: hidden; box-shadow: 0 30px 60px -20px rgba(0, 0, 0, 0.6); display: flex; flex-direction: column; }
.pp-header { display: flex; align-items: center; justify-content: space-between; padding: 14px 18px; color: #ffffff; }
.pp-title { font-size: 1rem; font-weight: 700; }
.pp-sub { font-size: 0.76rem; color: #94a3b8; margin-top: 1px; }
.pp-close { width: 34px; height: 34px; border-radius: 10px; border: none; background: rgba(255, 255, 255, 0.12); color: #ffffff; display: inline-flex; align-items: center; justify-content: center; cursor: pointer; transition: background 0.15s ease; }
.pp-close:hover { background: rgba(255, 255, 255, 0.22); }
.pp-stage { position: relative; height: min(520px, 60vh); background: #000000; display: flex; align-items: center; justify-content: center; }
.pp-image { width: 100%; height: 100%; object-fit: contain; animation: pp-fade 0.2s ease; }
@keyframes pp-fade { from { opacity: 0; } to { opacity: 1; } }
.pp-nav { position: absolute; top: 50%; transform: translateY(-50%); width: 42px; height: 42px; border-radius: 50%; border: none; background: rgba(255, 255, 255, 0.9); color: #0f172a; display: inline-flex; align-items: center; justify-content: center; cursor: pointer; box-shadow: 0 4px 14px rgba(0, 0, 0, 0.3); transition: transform 0.15s ease, background 0.15s ease; }
.pp-nav:hover { background: #ffffff; transform: translateY(-50%) scale(1.08); }
.pp-prev { left: 14px; }
.pp-next { right: 14px; }
.pp-counter { position: absolute; bottom: 12px; left: 50%; transform: translateX(-50%); background: rgba(15, 23, 42, 0.75); color: #ffffff; font-size: 0.74rem; font-weight: 600; padding: 3px 12px; border-radius: 999px; }
.pp-thumbs { display: flex; gap: 8px; padding: 12px 18px 16px; overflow-x: auto; justify-content: center; }
.pp-thumb { flex: 0 0 auto; width: 76px; height: 54px; border-radius: 8px; overflow: hidden; border: 2px solid transparent; padding: 0; background: #1e293b; cursor: pointer; opacity: 0.6; transition: all 0.15s ease; }
.pp-thumb img { width: 100%; height: 100%; object-fit: cover; display: block; }
.pp-thumb:hover { opacity: 0.9; }
.pp-thumb.active { border-color: #60a5fa; opacity: 1; }

/* Interactive Photo Banner */
.stall-photo-banner-interactive {
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}
.stall-photo-banner-interactive:hover {
  box-shadow: 0 6px 16px rgba(15, 23, 42, 0.16);
}
.stall-photo-banner-interactive:hover img {
  transform: scale(1.05);
}

/* Per-stall cards inside selected-stall drawer */
.sel-stall-list { display: flex; flex-direction: column; gap: 10px; margin-bottom: 12px; }
.sel-stall-card { border: 1px solid #e2e8f0; border-radius: 12px; padding: 10px; background: #ffffff; transition: border-color 0.15s ease, box-shadow 0.15s ease; }
.sel-stall-card:hover { border-color: #93c5fd; box-shadow: 0 2px 8px rgba(37, 99, 235, 0.08); }
.ssc-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 8px; margin-bottom: 10px; }
.ssc-title { display: flex; flex: 1; flex-direction: column; align-items: flex-start; gap: 3px; min-width: 0; }
.ssc-title-main { display: flex; align-items: center; gap: 7px; min-width: 0; }
.ssc-title strong { font-size: 1rem; color: #0f172a; white-space: nowrap; }
.ssc-status-badge { flex-shrink: 0; padding: 3px 8px !important; font-size: 0.8125rem !important; }
.ssc-floor-name { display: block; max-width: 280px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #64748b; font-size: 0.8125rem; font-weight: 600; line-height: 1.4; }
.ssc-remove { width: 24px; height: 24px; border-radius: 8px; border: 1px solid #e2e8f0; background: #ffffff; color: #64748b; display: inline-flex; align-items: center; justify-content: center; cursor: pointer; flex-shrink: 0; transition: all 0.15s ease; }
.ssc-remove:hover { background: #fef2f2; border-color: #fca5a5; color: #ef4444; }
.ssc-body { display: grid; grid-template-columns: 84px 1fr; gap: 10px; align-items: stretch; }
.ssc-photo { position: relative; border-radius: 8px; overflow: hidden; border: 1px solid #e2e8f0; cursor: pointer; min-height: 72px; padding: 0; background: #f1f5f9; }
.ssc-photo:hover img { transform: scale(1.05); }
.ssc-photo img { transition: transform 0.2s ease; }
.ssc-photo img { width: 100%; height: 100%; object-fit: cover; display: block; }
.ssc-photo-tag { position: absolute; bottom: 4px; right: 4px; display: inline-flex; align-items: center; gap: 3px; background: rgba(15, 23, 42, 0.78); color: #ffffff; font-size: 0.75rem; font-weight: 600; padding: 2px 6px; border-radius: 5px; }
.ssc-metrics { display: flex; flex-direction: column; justify-content: center; gap: 4px; min-width: 0; }
.ssc-metric { display: flex; align-items: baseline; justify-content: space-between; gap: 8px; font-size: 0.8125rem; color: #64748b; }
.ssc-metric b { font-size: 0.9375rem; font-weight: 700; color: #0f172a; white-space: nowrap; }
.ssc-metric b small { font-size: 0.75rem; font-weight: 500; color: #64748b; }
.ssc-metric b.is-money { color: #059669; }
.sel-stall-total { display: flex; flex-direction: column; gap: 2px; margin: -2px 0 12px; padding: 8px 10px; border-radius: 10px; background: #f8fafc; font-size: 0.8125rem; color: #64748b; }
.sel-stall-total span:first-child { font-weight: 700; color: #334155; }

.map-selection-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  margin: -8px 0 24px;
  padding: 16px 20px;
  background: #ffffff;
  border: 1px solid #bfdbfe;
  border-radius: 14px;
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.08);
}

.msb-left {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
  flex: 1 1 320px;
}

.msb-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.92rem;
  font-weight: 600;
  color: #1e293b;
}

.msb-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  max-height: 72px;
  overflow-y: auto;
}

.msb-right {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
}

.msb-metric {
  display: flex;
  flex-direction: column;
  font-size: 0.72rem;
  color: #64748b;
}

.msb-metric strong {
  font-size: 0.98rem;
  font-weight: 700;
  color: #1e293b;
}

.msb-metric strong.is-money {
  color: #10b981;
}

.msb-clear {
  border-color: #cbd5e1;
  color: #475569;
  background: #f1f5f9;
  font-weight: 600;
  border-radius: 8px;
}

.msb-cta {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 10px;
  font-weight: 700;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

/* Map Canvas Controls */
.map-ctrl-btn {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #334155;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.06);
  transition: all 0.15s ease;
}

.map-ctrl-btn:hover {
  background: #f8fafc;
  color: var(--primary);
  border-color: #cbd5e1;
}

/* View Mode Segmented Tab Switcher */
.view-mode-tabs-container {
  display: flex;
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  padding: 4px;
  border-radius: 20px;
  gap: 4px;
}

.view-mode-tab-btn {
  border: none;
  background: transparent;
  padding: 6px 14px;
  border-radius: 16px;
  font-size: 0.82rem;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.view-mode-tab-btn:hover {
  color: var(--primary);
}

.view-mode-tab-btn.active {
  background: #ffffff;
  color: var(--primary);
  font-weight: 700;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
}

/* Zone Filter Pills Bar */
.zone-filter-bar {
  display: flex;
  gap: 6px;
  background: #f1f5f9;
  padding: 4px;
  border-radius: 20px;
  border: 1px solid #cbd5e1;
}

.zone-pill-btn {
  border: none;
  background: transparent;
  padding: 6px 14px;
  border-radius: 16px;
  font-size: 0.8rem;
  font-weight: 600;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s ease;
}

.zone-pill-btn.active {
  background: #ffffff;
  color: var(--primary);
  font-weight: 700;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

/* Floorplan Landmark Banners */
.floorplan-landmark-tag {
  background: rgba(15, 23, 42, 0.88);
  backdrop-filter: blur(4px);
  color: #ffffff;
  padding: 5px 16px;
  border-radius: 20px;
  font-size: 0.78rem;
  font-weight: 700;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.floorplan-landmark-tag.entrance {
  background: #1e3a8a;
  border: 1px solid #60a5fa;
}

.floorplan-landmark-tag.stage {
  background: #0f172a;
  border: 1px solid #f59e0b;
  color: #fef3c7;
}

/* Building Layout Perimeter Shape */
.building-polygon-border {
  border: 2px solid #cbd5e1;
  background: #ffffff;
  border-radius: 12px;
  padding: 16px 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
  position: relative;
  transition: all 0.3s ease;
}

.building-polygon-border.with-heatmap-glow {
  box-shadow: 0 0 24px rgba(245, 158, 11, 0.25);
  border-color: #f59e0b;
}

/* Stalls Grid Sections */
.stalls-grid-section {
  display: grid;
  grid-template-columns: repeat(14, 28px);
  gap: 8px;
  justify-content: center;
}

.market-aisle-divider {
  margin: 14px 0;
  padding: 6px;
  background: #f1f5f9;
  border-radius: 6px;
  text-align: center;
  font-size: 0.72rem;
  font-weight: 600;
  color: #94a3b8;
  letter-spacing: 0.05em;
  transition: all 0.3s ease;
}

.market-aisle-divider.traffic-hot {
  background: linear-gradient(90deg, #fee2e2 0%, #fef3c7 50%, #fee2e2 100%);
  color: #b45309;
  font-weight: 800;
  border: 1px dashed #f59e0b;
  box-shadow: 0 0 10px rgba(245, 158, 11, 0.2);
}

/* Individual Stall Box */
.stall-box-item {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 0.65rem;
  font-weight: 700;
  position: relative;
  transition: all 0.15s ease;
  user-select: none;
}

.stall-mini-code {
  display: none;
}

.stall-box-item.is-occupied {
  background-color: #cbd5e1;
  border: 1px solid #94a3b8;
  cursor: not-allowed;
}

.stall-box-item.is-reserved {
  background-color: #f59e0b;
  border: 1px solid #d97706;
}

.stall-box-item.is-reserved:hover {
  transform: scale(1.15);
  box-shadow: 0 0 10px rgba(245, 158, 11, 0.5);
}

.stall-box-item.is-available {
  background-color: #10b981;
  border: 1px solid #059669;
}

.stall-box-item.is-available:hover {
  transform: scale(1.15);
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.5);
}

.stall-box-item.is-selected {
  outline: 3px solid var(--primary);
  outline-offset: 2px;
  transform: scale(1.2);
  z-index: 5;
}

.stall-box-item.traffic-hotspot {
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.6);
  border-color: #ef4444;
}

.stall-box-item.is-dimmed {
  opacity: 0.2;
}

.selected-check {
  color: #ffffff;
  font-size: 0.8rem;
  font-weight: 800;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}


.input-verified-glow {
  border-color: #10b981 !important;
  background-color: #f0fdf4 !important;
}

/* Date Cards Strip */
.date-cards-strip {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}

@media (max-width: 768px) {
  .date-cards-strip {
    display: flex;
    overflow-x: auto;
    padding-bottom: 8px;
    scroll-snap-type: x mandatory;
    -webkit-overflow-scrolling: touch;
  }
  .date-card {
    min-width: 135px;
    flex-shrink: 0;
    scroll-snap-align: start;
  }
}

.date-card {
  border: 1.5px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px 14px;
  background: #ffffff;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  text-align: left;
}

.date-card:hover {
  border-color: #93c5fd;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.08);
}

.date-card.selected {
  border-color: var(--primary, #2563eb);
  background: #eff6ff;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
}

.date-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
  gap: 4px;
}

.date-day-name {
  font-size: 0.8rem;
  font-weight: 700;
  color: #475569;
}

.date-card.selected .date-day-name {
  color: #1e40af;
}

.date-badge {
  font-size: 0.68rem;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 6px;
  white-space: nowrap;
}

.date-badge.badge-primary {
  background: #dbeafe;
  color: #1d4ed8;
}

.date-badge.badge-success {
  background: #dcfce7;
  color: #15803d;
}

.date-badge.badge-purple {
  background: #f3e8ff;
  color: #7e22ce;
}

.date-badge.badge-neutral {
  background: #f1f5f9;
  color: #475569;
}

.date-number-wrap {
  display: flex;
  align-items: baseline;
  gap: 6px;
  margin-bottom: 8px;
}

.date-day-num {
  font-size: 1.5rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1;
}

.date-card.selected .date-day-num {
  color: #1d4ed8;
}

.date-month {
  font-size: 0.78rem;
  color: #64748b;
  font-weight: 600;
}

.date-card-footer {
  font-size: 0.72rem;
  border-top: 1px dashed #e2e8f0;
  padding-top: 6px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.date-card.selected .date-card-footer {
  border-top-color: #bfdbfe;
}

.date-status.status-open {
  color: #059669;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
}

.date-status.status-limited {
  color: #d97706;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
}

.date-status.status-closed {
  color: #94a3b8;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
}

.date-card.is-closed-date {
  opacity: 0.6;
  background: #f8fafc;
  border-color: #e2e8f0;
  cursor: not-allowed;
}

.date-card.is-closed-date:hover {
  transform: none;
  box-shadow: none;
  border-color: #cbd5e1;
}

/* Interactive Time Slot Cards */
.time-slot-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.time-slot-card {
  border: 1.5px solid #e2e8f0;
  border-radius: 12px;
  padding: 14px;
  background: #ffffff;
  cursor: pointer;
  transition: all 0.2s ease;
}

.time-slot-card:hover {
  border-color: #93c5fd;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.08);
}

.time-slot-card.selected {
  border-color: var(--primary);
  background: #eff6ff;
  box-shadow: 0 0 0 2px rgba(67, 97, 238, 0.2);
}

.time-slot-card.is-full-slot {
  background: #f8fafc !important;
  border-color: #e2e8f0 !important;
  opacity: 0.62;
  cursor: not-allowed !important;
  box-shadow: none !important;
  transform: none !important;
  position: relative;
}

.time-slot-card.is-full-slot:hover {
  border-color: #fca5a5 !important;
  background: #fff5f5 !important;
}

.cap-full {
  background: #fee2e2;
  color: #dc2626;
  font-weight: 700;
}

.slot-card-top {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.slot-card-icon {
  font-size: 1.2rem;
}

.slot-card-time {
  font-size: 0.92rem;
  color: #0f172a;
}

.slot-card-capacity {
  font-size: 0.75rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 6px;
  display: inline-block;
  margin-bottom: 6px;
}

.cap-low {
  background: #fee2e2;
  color: #dc2626;
}

.cap-med {
  background: #fef3c7;
  color: #d97706;
}

.cap-high {
  background: #dcfce7;
  color: #15803d;
}

.slot-card-note {
  font-size: 0.74rem;
  color: #64748b;
  line-height: 1.3;
}
.slot-booked-summary { white-space: nowrap; }

/* Payment Schedule Grid (BA Perspective) */
.payment-schedule-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 14px;
}

.schedule-card {
  border: 1.5px solid #cbd5e1;
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  position: relative;
  background: #ffffff;
  transition: all 0.2s ease;
}

.schedule-card:hover {
  border-color: #93c5fd;
  transform: translateY(-2px);
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.08);
}

.schedule-card.active {
  border-color: var(--primary);
  background: #f8faff;
  box-shadow: 0 0 0 2px rgba(67, 97, 238, 0.25);
}

.schedule-badge {
  position: absolute;
  top: -10px;
  right: 12px;
  background: #10b981;
  color: #ffffff;
  font-size: 0.72rem;
  font-weight: 800;
  padding: 2px 8px;
  border-radius: 10px;
  box-shadow: 0 2px 6px rgba(16, 185, 129, 0.3);
}

.schedule-badge.gold {
  background: linear-gradient(135deg, #d97706 0%, #b45309 100%);
}

.schedule-top {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 0.92rem;
  color: #0f172a;
}

.schedule-price {
  font-size: 1.25rem;
  font-weight: 800;
  color: #1e40af;
  margin-bottom: 4px;
}

.schedule-price .unit {
  font-size: 0.78rem;
  font-weight: 500;
  color: #64748b;
}

.schedule-desc {
  font-size: 0.76rem;
  color: #64748b;
  line-height: 1.35;
}

/* 3D Virtual Tour Hotspot Pins */
.virtual-pin {
  position: absolute;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  z-index: 10;
}

.pin-dot {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #ef4444;
  border: 3px solid #ffffff;
  box-shadow: 0 0 10px rgba(239, 68, 68, 0.8);
  animation: pulsePin 1.5s infinite;
}

@keyframes pulsePin {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.7);
  }
  50% {
    transform: scale(1.15);
    box-shadow: 0 0 0 8px rgba(239, 68, 68, 0);
  }
}

.pin-label {
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(4px);
  color: #ffffff;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 700;
  white-space: nowrap;
  pointer-events: none;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

/* Confetti Overlay */
.confetti-canvas-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 20;
}

/* Bottom Map Style Switcher */
.map-style-switcher {
  position: absolute;
  bottom: 16px;
  left: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  background: #ffffff;
  padding: 6px 10px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  z-index: 10;
}

.style-thumb {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  cursor: pointer;
  padding: 4px 6px;
  border-radius: 8px;
  transition: all 0.15s ease;
}

.style-thumb.active {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
}

.thumb-preview {
  width: 44px;
  height: 28px;
  border-radius: 4px;
  border: 1px solid #cbd5e1;
}

.thumb-preview.standard {
  background: linear-gradient(135deg, #e0f2fe 0%, #bae6fd 100%);
}

.thumb-preview.dark {
  background: #1e293b;
}

.thumb-preview.satellite {
  background: linear-gradient(135deg, #334155 0%, #15803d 100%);
}

.style-thumb span {
  font-size: 0.68rem;
  font-weight: 600;
  color: #475569;
}

/* Bottom Legend */
.map-bottom-legend {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 16px;
  background: #ffffff;
  padding: 8px 18px;
  border-radius: 20px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  z-index: 10;
}

.legend-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8rem;
  color: #475569;
  font-weight: 500;
}

.legend-box {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}

.legend-box.occupied {
  background: #475569;
}

.legend-box.reserved {
  background: #f59e0b;
}

.legend-box.available {
  background: #10b981;
}

/* Transitions & Animations */
.slide-fade-enter-active, .slide-fade-leave-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-fade-enter-from, .slide-fade-leave-to {
  transform: translateX(20px);
  opacity: 0;
}

@keyframes bounceLeft {
  0%, 100% {
    transform: translateX(0);
  }
  50% {
    transform: translateX(-7px);
  }
}

.stalls-interactive-canvas.is-panning {
  cursor: grabbing !important;
}

/* Floating Zoom Controls */
.map-floating-zoom-ctrl {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(8px);
  padding: 6px;
  border-radius: 12px;
  border: 1px solid #bfdbfe;
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.12);
}

.zoom-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  color: #1e3a8a;
  font-size: 1.15rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
  line-height: 1;
}

.zoom-btn:hover {
  background: #eff6ff;
  border-color: #3b82f6;
  color: #2563eb;
  transform: scale(1.06);
}

.zoom-btn.reset-btn {
  font-size: 0.85rem;
}

.zoom-level-pill {
  font-size: 0.68rem;
  font-weight: 800;
  color: #475569;
  text-align: center;
  padding: 3px 2px;
  cursor: pointer;
  background: #f1f5f9;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  user-select: none;
  transition: all 0.15s ease;
}

.zoom-level-pill:hover {
  background: #e2e8f0;
  color: #1e3a8a;
}

.hand-pointing-icon {
  animation: bounceLeft 1.5s infinite ease-in-out;
}

/* Multi-Stall Selection Status Bar */
.selection-status-indicator {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 18px;
  border-radius: 12px;
  margin-bottom: 12px;
  font-size: 0.85rem;
  line-height: 1.4;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.06);
  position: relative;
  z-index: 10;
}

.status-available-mode {
  background: linear-gradient(135deg, #ecfdf5 0%, #f0fdf4 100%);
  border: 1.5px solid #6ee7b7;
  color: #065f46;
}

.status-reserved-mode {
  background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%);
  border: 1.5px solid #fcd34d;
  color: #92400e;
}

.selection-status-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.selection-status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

.status-available-mode .selection-status-dot {
  background: #10b981;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.25);
  animation: pulse 1.8s infinite;
}

.status-reserved-mode .selection-status-dot {
  background: #f59e0b;
  box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.25);
  animation: pulse 1.8s infinite;
}

.selection-status-right {
  flex-shrink: 0;
}

.btn-clear-selection {
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 8px;
  padding: 5px 12px;
  font-size: 0.78rem;
  font-weight: 700;
  color: inherit;
  cursor: pointer;
  transition: all 0.15s ease;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.btn-clear-selection:hover {
  background: #ffffff;
  border-color: rgba(0, 0, 0, 0.25);
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
}

/* Selected Stall Chips */
.selected-stalls-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.stall-chip {
  border-radius: 9999px;
  font-size: 0.78rem;
  font-weight: 700;
  padding: 4px 10px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  transition: all 0.18s ease;
}

.stall-chip.chip-available {
  background: #d1fae5;
  color: #065f46;
  border: 1px solid #6ee7b7;
}

.stall-chip.chip-reserved {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fcd34d;
}

.stall-chip-remove {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.1);
  border: none;
  font-size: 9px;
  line-height: 1;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: inherit;
  padding: 0;
  transition: all 0.15s ease;
}

.stall-chip-remove:hover {
  background: rgba(0, 0, 0, 0.25);
  color: #ffffff;
  transform: scale(1.1);
}

/* Conflict Mode Styling for Floorplan Stalls */
.stall-box-item.is-conflict-mode {
  opacity: 0.38 !important;
  filter: grayscale(0.65) !important;
  position: relative;
  cursor: not-allowed;
}

.stall-box-item.is-conflict-mode:hover {
  opacity: 0.7 !important;
  filter: grayscale(0) !important;
  outline: 2px dashed #f59e0b;
  outline-offset: 2px;
  cursor: pointer;
}

/* Floorplan Stall Queue Pill */
.stall-queue-pill {
  position: absolute;
  top: 1px;
  right: 1px;
  background: #d97706;
  color: #ffffff;
  font-size: 0.52rem;
  font-weight: 800;
  line-height: 1;
  padding: 1.5px 3px;
  border-radius: 3px;
  letter-spacing: -0.02em;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  pointer-events: none;
  z-index: 2;
}

/* Custom Smooth Scrollbar for Selected Stall Drawer */
.selected-stall-drawer {
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 #f8fafc;
}

.selected-stall-drawer::-webkit-scrollbar {
  width: 6px;
}

.selected-stall-drawer::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 8px;
  margin: 6px 0;
}

.selected-stall-drawer::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 8px;
}

.selected-stall-drawer::-webkit-scrollbar-thumb:hover {
  background-color: #94a3b8;
}

/* Responsive adjustments for Booking Canvas & Drawers */
@media (max-width: 992px) {
  .stall-view-top-bar {
    position: static !important;
    right: auto !important;
    left: auto !important;
    padding: 12px 14px !important;
    background: #ffffff !important;
    border-bottom: 1px solid #e2e8f0;
    flex-wrap: wrap;
    gap: 8px;
  }

  .zone-filter-bar {
    overflow-x: auto;
    width: 100%;
    padding-bottom: 4px;
    -webkit-overflow-scrolling: touch;
  }
}

.mobile-stall-hint-pill {
  display: none;
}

.mobile-project-full-banner {
  display: none;
}

@media (max-width: 992px) {
  .map-selection-empty-drawer {
    display: none !important;
  }

  .mobile-stall-hint-pill {
    position: absolute;
    top: 60px;
    left: 12px;
    z-index: 20;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    background: rgba(15, 23, 42, 0.84);
    backdrop-filter: blur(8px);
    color: #ffffff;
    border-radius: 9999px;
    font-size: 0.74rem;
    font-weight: 600;
    box-shadow: 0 4px 12px rgba(15, 23, 42, 0.2);
    pointer-events: none;
  }

  .mobile-stall-hint-pill .mshp-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #10b981;
    box-shadow: 0 0 6px #10b981;
  }

  .mobile-project-full-banner {
    position: absolute;
    top: 56px;
    left: 12px;
    right: 12px;
    z-index: 25;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 9px 14px;
    background: rgba(254, 242, 242, 0.98);
    backdrop-filter: blur(8px);
    border: 1px solid #fca5a5;
    border-radius: 12px;
    color: #b91c1c;
    font-size: 0.78rem;
    font-weight: 600;
    box-shadow: 0 4px 14px rgba(220, 38, 38, 0.15);
  }

  .selected-stall-drawer {
    position: fixed !important;
    bottom: 0 !important;
    top: auto !important;
    left: 0 !important;
    right: 0 !important;
    width: 100% !important;
    max-height: 80vh !important;
    border-radius: 20px 20px 0 0 !important;
    z-index: 99 !important;
    box-shadow: 0 -10px 30px rgba(15, 23, 42, 0.25) !important;
    border-bottom: none !important;
    border-left: none !important;
    border-right: none !important;
  }
}

@media (max-width: 768px) {
  .map-style-switcher {
    bottom: 10px !important;
    left: 10px !important;
    padding: 4px 6px !important;
  }

  .style-thumb {
    padding: 2px 6px !important;
    font-size: 0.7rem !important;
  }

  .thumb-preview {
    width: 16px !important;
    height: 16px !important;
  }

  .view-mode-tabs-container {
    width: 100%;
    display: flex;
  }

  .view-mode-tab-btn {
    flex: 1;
    justify-content: center;
    padding: 6px 8px !important;
    font-size: 0.76rem !important;
  }
}

/* แจ้งเตือนตลาดเต็ม */
.market-full-banner {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  background: #fef2f2;
  border: 1.5px solid #fecaca;
  border-radius: 14px;
  padding: 16px 20px;
  margin-bottom: 20px;
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.05);
}

.market-full-banner-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: #fee2e2;
  color: #dc2626;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.market-full-banner-content {
  flex: 1;
}

.market-full-title {
  font-size: 0.98rem;
  font-weight: 800;
  color: #991b1b;
  margin-bottom: 4px;
}

.market-full-desc {
  font-size: 0.84rem;
  color: #b91c1c;
  line-height: 1.55;
  margin: 0;
}

.stall-box-item.is-market-full {
  cursor: not-allowed !important;
  opacity: 0.75;
}

.stall-box-item.is-market-full:hover {
  transform: none !important;
  box-shadow: none !important;
}

.is-full-guide {
  border-color: #fca5a5 !important;
  background: rgba(254, 242, 242, 0.95) !important;
}

/* Compact Step 3 review layout: keep the GIS component prominent and reduce repeated detail cards. */
.step3-review { color: #0f172a; }
.step3-review { display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1fr); column-gap: 24px; row-gap: 24px; align-items: start; }
.step3-review > .step3-map-card, .step3-review > .step3-stalls-card, .step3-review > .step3-review-grid, .step3-review > .step3-finance-card, .step3-review > .step3-guarantee-card { margin-bottom: 0 !important; }
.step3-review > div:first-child,
.step3-review > .step3-review-grid,
.step3-review > .step3-finance-card,
.step3-review > .step3-guarantee-card,
.step3-review > div:last-child { grid-column: 1 / -1; }
.step3-review > .step3-map-card { grid-column: 2; grid-row: 2; width: auto; height: auto; align-self: stretch; margin-bottom: 0 !important; }
.step3-review > .step3-stalls-card { grid-column: 1; grid-row: 2; width: auto; height: 760px !important; align-self: start; margin-bottom: 0 !important; }
.step3-stalls-card.has-multiple-stalls { height: 760px !important; max-height: 760px; overflow: hidden; display: flex; flex-direction: column; overscroll-behavior: contain; }
.step3-stalls-scroll { min-height: 0; }
.step3-stalls-card.has-multiple-stalls .step3-stalls-scroll { flex: 1 1 auto; min-height: 0; overflow-y: auto; overscroll-behavior: contain; padding-right: 4px; }
.step3-review > .step3-review-grid { grid-column: 1 / -1; grid-row: 3; grid-template-columns: 1fr 1fr; gap: 24px !important; align-self: start; }
.step3-review > .step3-finance-card { grid-column: 1 / -1; grid-row: 4; align-self: start; }
.step3-review > .step3-guarantee-card { grid-column: 1 / -1; grid-row: 5; }
.step3-review > div:last-child { grid-row: 6; }
.step3-review > div:first-child { margin-bottom: 16px !important; }
.step3-map-card { margin-bottom: 0 !important; }
.step3-map-heading { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin: 0 24px; padding: 20px 0 12px; border-bottom: 1px solid #eef2f6; }
.step3-map-heading-copy { display: flex; align-items: center; gap: 14px; min-width: 0; }
.step3-map-heading-icon { width: 40px; height: 40px; display: inline-flex; align-items: center; justify-content: center; flex: 0 0 40px; border: 0; border-radius: 10px; background: #eff6ff; color: #2563eb; }
.step3-map-heading-icon svg { width: 20px; height: 20px; }
.step3-map-heading h3 { margin: 0; color: #0f172a; font-size: var(--font-size-lg); font-weight: 700; }
.step3-map-heading p { margin: 4px 0 0; color: #64748b; font-size: var(--font-size-sm); }
.step3-map-heading-legend { display: inline-flex; align-items: center; gap: 14px; color: #475569; font-size: var(--font-size-sm); white-space: nowrap; }
.step3-map-heading-legend span { display: inline-flex; align-items: center; gap: 6px; }
.step3-map-heading-legend i { width: 14px; height: 14px; display: inline-block; border-radius: 4px; }
.step3-map-heading-legend i.selected { background: #2563eb; box-shadow: 0 0 0 3px #dbeafe; }
.step3-map-heading-legend i.available { background: #10b981; }
.step3-map-heading-legend i.reserved { background: #f59e0b; }
.step3-map-card :deep(.preview-top-bar) { display: none !important; }
.step3-map-card :deep(.projects-master-container.is-preview-mode) { border-radius: 16px; overflow: hidden; }
.step3-map-card :deep(.projects-master-container.is-preview-mode) { padding: 0 !important; gap: 0 !important; }
.step3-map-card :deep(.projects-master-container.is-preview-mode .maplibre-overview-wrapper) { width: calc(100% - 48px); height: 620px !important; min-height: 620px !important; margin: 24px 24px 24px; }
.step3-map-card :deep(.preview-top-bar) { min-height: 126px; padding: 24px 28px 16px; }
.step3-map-card :deep(.preview-head-left) { gap: 14px; }
.step3-map-card :deep(.preview-head-icon-box) { width: 48px; height: 48px; border-radius: 13px; }
.step3-map-card :deep(.preview-head-title) { font-size: var(--font-size-xl); }
.step3-map-card :deep(.preview-head-subtitle) { font-size: var(--font-size-sm); }
.step3-map-card :deep(.preview-booked-pill-btn) { border-radius: 14px; }
.step3-map-card :deep(.preview-legend-row) { gap: 14px; font-size: var(--font-size-sm); }
.step3-stalls-card { margin-bottom: 16px !important; padding: 24px !important; border-radius: 16px !important; }
.step3-stalls-card > div:first-child > div:first-child { width: 40px !important; height: 40px !important; border-radius: 10px !important; }
.step3-stalls-card > div:first-child h3 { font-size: var(--font-size-lg) !important; }
.step3-stalls-card > div[style*="border: 1px solid #e2e8f0"] { padding: 16px !important; margin-bottom: 12px !important; }
.step3-stalls-card > div[style*="border: 1px solid #e2e8f0"] > div[style*="display: grid"] { gap: 10px 14px !important; }
.step3-stalls-card > div[style*="display: grid"] { gap: 12px 18px !important; margin-bottom: 14px !important; }
.step3-stalls-card > div[style*="display: grid"] > div > div:first-child { font-size: var(--font-size-xs) !important; }
.step3-stalls-card > div[style*="display: grid"] > div > div:last-child { font-size: var(--font-size-sm) !important; }
.step3-photos { display: block !important; margin-top: 16px; padding-top: 14px; border-top: 1px solid #f1f5f9; }
.step3-photos > div:last-child { display: flex; gap: 10px !important; flex-wrap: nowrap !important; overflow-x: auto; padding-bottom: 2px; }
.step3-photos .summary-stall-photo-card { width: 112px !important; height: 78px !important; flex: 0 0 112px; border-radius: 9px; }
.step3-review-grid { gap: 16px !important; margin-bottom: 16px !important; }
.step3-review-grid > .card { padding: 24px !important; }
.step3-review-grid > .card > div[style*="display: grid"] { row-gap: 18px !important; }
.step3-review-grid { gap: 16px !important; }
.step3-review-grid > .card > div:first-child > div:first-child { width: 40px !important; height: 40px !important; border-radius: 10px !important; }
.step3-review-grid > .card > div:first-child > div:first-child svg { width: 20px; height: 20px; }
.step3-review-grid > .card > div:first-child > span { font-size: var(--font-size-lg) !important; font-weight: 700 !important; }
.step3-review-grid [style*="font-size: 0.8rem"] { font-size: var(--font-size-xs) !important; }
.step3-review-grid [style*="font-size: 0.88rem"] { font-size: var(--font-size-sm) !important; }
.step3-finance-card { margin-bottom: 16px !important; padding: 20px !important; }
.step3-review > .step3-finance-card { margin-top: 0 !important; }
.step3-finance-icon { width: 40px; height: 40px; display: inline-flex; align-items: center; justify-content: center; border-radius: 10px; background: #eff6ff; color: #2563eb; flex: 0 0 40px; }
.step3-finance-icon svg { width: 20px; height: 20px; }
.step3-finance-card > div:last-child { font-size: var(--font-size-sm) !important; }
.step3-guarantee-card { margin-bottom: 16px !important; padding: 12px 16px !important; }
.step3-review > .step3-guarantee-card { margin-top: 0 !important; }
.step3-review > div:last-child { margin-top: 0 !important; }
@media (max-width: 900px) {
  .step3-review { display: block; }
  .step3-review > .step3-stalls-card, .step3-stalls-card.has-multiple-stalls { height: auto !important; max-height: none; overflow: visible; }
  .step3-stalls-card.has-multiple-stalls .step3-stalls-scroll { overflow: visible; }
  .step3-review > .step3-map-card, .step3-review > .step3-stalls-card { margin-bottom: 16px !important; }
  .step3-map-card :deep(.preview-top-bar) { min-height: 112px; padding: 18px 18px 12px; }
  .step3-map-heading { align-items: flex-start; flex-direction: column; margin: 0 16px; padding: 18px 0 12px; }
  .step3-map-heading-legend { flex-wrap: wrap; white-space: normal; }
  .step3-map-card :deep(.projects-master-container.is-preview-mode .maplibre-overview-wrapper) { width: calc(100% - 32px); height: 500px !important; min-height: 500px; margin: 18px 16px 16px; }
  .step3-review-grid { grid-template-columns: 1fr !important; }
}
</style>
