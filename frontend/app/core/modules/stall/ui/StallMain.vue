<template>
  <div>
    <!-- Top Header -->
    <AppHeader />

    <div class="content-body market-detail-page">
      <header class="md-page-heading">
        <div>
          <h2>รายละเอียดตลาด</h2>
          <p>ตรวจสอบข้อมูลสำคัญและอัปเดตรายละเอียดตลาดเมื่อมีการเปลี่ยนแปลง</p>
        </div>
        <div v-if="isOfficer" class="md-page-actions">
          <AppButton v-if="isAdmin" :to="gisManageLink" variant="outline" size="md" title="จัดการผังและตำแหน่งแผงค้า GIS">
            <template #prefix>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18-6 3V5l6-3 6 3 6-3v12"></path>
                <path d="M9 2v16"></path>
                <path d="M15 5v6"></path>
                <path d="m14 18 6.5-6.5a1.7 1.7 0 0 0-2.4-2.4L11.6 15.6 11 19z"></path>
              </svg>
            </template>
            จัดการผัง GIS
          </AppButton>
          <AppButton v-if="isAdmin" :to="visitScheduleLink" variant="outline" size="md" title="จัดการรอบเวลาเข้าชมแผง">
            <template #prefix>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="9"></circle>
                <path d="M12 7v5l3 2"></path>
              </svg>
            </template>
            จัดการรอบเข้าชม
          </AppButton>
          <AppButton :to="editMarketLink" variant="outline" size="md" title="แก้ไขข้อมูลตลาด">
            <template #prefix>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 20h9"></path>
                <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"></path>
              </svg>
            </template>
            แก้ไขข้อมูล
          </AppButton>
        </div>
      </header>

      <!-- 1. Header Card: Real Market Photo Backdrop + Modern Glassmorphic Hero -->
      <section class="md-card md-header-card">
        <!-- Real Market Photo Background Layer with Cinematic Frosting -->
        <div class="md-header-bg-frame" aria-hidden="true">
          <img :src="marketHeroImage" :alt="currentMarket.name" class="md-header-bg-img" />
          <div class="md-header-overlay-gradient"></div>
          <div class="md-header-overlay-mesh"></div>
          <div class="md-header-ambient-glow"></div>
          <div class="md-header-light-beam"></div>
          <div class="md-header-grid-pattern"></div>
        </div>

        <div class="md-header-left">
          <!-- Official Subtitle / Eyebrow Pill -->
          <div class="md-header-eyebrow">
            <span class="md-eyebrow-pill">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
              <span>สำนักงานตลาดกรุงเทพมหานคร</span>
            </span>
          </div>

          <!-- Market Name -->
          <div class="md-title-wrap">
            <h1 class="md-title">{{ currentMarket.name }}</h1>
          </div>

          <p class="md-hero-summary">
            {{ currentMarket.description || 'พื้นที่ค้าขายคุณภาพภายใต้การดูแลของกรุงเทพมหานคร พร้อมผังแผงดิจิทัลที่ช่วยให้เลือกทำเลได้อย่างมั่นใจ' }}
          </p>

          <!-- Meta Information Row: วัน, เวลา, พิกัด -->
          <div class="md-meta-row">
            <div class="md-meta-item">
              <span class="md-meta-icon">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
              </span>
              <div class="md-meta-content">
                <span class="md-meta-label">วันเปิดทำการ</span>
                <strong class="md-meta-val">{{ operatingDays }}</strong>
              </div>
            </div>

            <span class="md-meta-divider" aria-hidden="true"></span>

            <div class="md-meta-item">
              <span class="md-meta-icon">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
              </span>
              <div class="md-meta-content">
                <span class="md-meta-label">เวลาเปิด-ปิด</span>
                <strong class="md-meta-val">{{ operatingTime }}</strong>
              </div>
            </div>

            <span v-if="currentMarket.district" class="md-meta-divider" aria-hidden="true"></span>

            <div v-if="currentMarket.district" class="md-meta-item">
              <span class="md-meta-icon">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
              </span>
              <div class="md-meta-content">
                <span class="md-meta-label">ทำเล / เขต</span>
                <strong class="md-meta-val">{{ currentMarket.district }}</strong>
              </div>
            </div>
          </div>

        </div>

        <!-- Right Side: Lightweight booking action -->
        <div class="md-header-right">
          <div class="md-booking-quick-card">
            <div class="md-bqc-simple-summary">
              <div class="md-bqc-status-row">
                <span class="md-bqc-badge" :class="availNow > 0 ? 'is-available' : 'is-full'">
                  <span class="md-bqc-dot" :class="{ 'pulsing': availNow > 0 }"></span>
                  <span>{{ availNow > 0 ? 'เปิดรับจองออนไลน์' : 'ยังไม่มีแผงว่าง' }}</span>
                </span>
              </div>
              <span
                class="md-bqc-summary-count"
                :class="[
                  availNow > 0 ? 'color-available' : 'color-full',
                  { 'is-three-digits': availNow >= 100 && availNow < 1000, 'is-four-plus-digits': availNow >= 1000 }
                ]"
              >
                <strong>{{ availNow }}</strong>
                <small>แผง</small>
              </span>
              <span class="md-bqc-summary-copy">
                <span class="md-bqc-copy-title">
                  <span>{{ availNow > 0 ? 'แผงว่าง' : 'แผงเต็มแล้ว' }}</span>
                  <strong :class="availNow > 0 ? 'is-ready' : 'is-unavailable'">
                    {{ availNow > 0 ? 'พร้อมจองทันที' : 'ปิดรับจองชั่วคราว' }}
                  </strong>
                </span>
                <small>{{ availNow > 0 ? 'จากทั้งหมด' : 'รวมทั้งหมด' }} <strong>{{ totalStalls.toLocaleString() }}</strong> แผง</small>
              </span>
            </div>

            <!-- Full-Width High-Affordance CTA Button -->
            <NuxtLink
              id="md-bqc-primary-cta"
              :to="bookLink"
              class="md-bqc-action-btn"
              :class="availNow > 0 ? 'btn-active' : 'btn-full'"
            >
              <span class="md-bqc-btn-label">
                {{ isOfficer ? 'จัดการแผงเช่า GIS' : (availNow > 0 ? 'จองแผงเช่าออนไลน์' : 'ดูผังแผงค้า') }}
              </span>
              <svg class="md-bqc-btn-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <line x1="5" y1="12" x2="19" y2="12"></line>
                <polyline points="12 5 19 12 12 19"></polyline>
              </svg>
            </NuxtLink>

          </div>
        </div>
      </section>

      <!-- 2. Floorplan Section -->
      <section id="md-floorplan" class="md-card md-floorplan-section">
        <h2 class="md-section-title">ภาพแผนผังตลาด</h2>
        <span class="md-section-accent"></span>
        <p class="md-section-sub">ดูภาพรวมการจัดโซนแผงค้า แล้วเข้าสู่ผังแผงดิจิทัลเพื่อเลือกแผงที่ต้องการ</p>

        <!-- แจ้งเตือนเมื่อตลาดเต็ม -->
        <div v-if="availNow <= 0" class="md-full-market-alert">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
          <div class="md-full-alert-content">
            <strong>ตลาดนี้มีผู้เช่าเต็มแล้ว (ยังไม่สามารถทำการจองได้)</strong>
            <span>ขณะนี้แผงค้าทุกแผงมีผู้เช่าครบแล้ว ท่านยังสามารถเข้าชมผังแผงดิจิทัลเพื่อดูตำแหน่งและโซนแผงค้าได้</span>
          </div>
        </div>

        <div
          class="md-floorplan-img-link"
          role="button"
          tabindex="0"
          @click="isFloorplanLightboxOpen = true"
          @keydown.enter="isFloorplanLightboxOpen = true"
          title="คลิกเพื่อดูรูปผังตลาดขนาดใหญ่"
          style="cursor: zoom-in;"
        >
          <img :src="currentMarket.floorplan_image_url || '/images/chatuchak_isometric_plan.jpg'" :alt="`ภาพแผนผัง${currentMarket.name}`" class="md-floorplan-img" />
          <span class="md-floorplan-hint">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line><line x1="11" y1="8" x2="11" y2="14"></line><line x1="8" y1="11" x2="14" y2="11"></line></svg>
            <span>คลิกเพื่อดูรูปผังตลาดขนาดเต็ม</span>
          </span>
        </div>
      </section>

      <!-- 3. Main Grid: Stats + Detail (ซ้าย) / Photos + Video (ขวา) -->
      <div class="md-main-grid">
        <div class="md-main-left">
          <!-- 2x2 Stat Cards -->
          <div class="md-stats-grid">
            <div class="md-card md-stat-card tone-green">
              <div class="md-stat-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="6" width="20" height="12" rx="2"></rect><circle cx="12" cy="12" r="2.5"></circle><path d="M6 12h.01M18 12h.01"></path></svg>
              </div>
              <div class="md-stat-body">
                <div class="md-stat-label">ราคาเช่า</div>
                <div class="md-stat-value">{{ priceRange }} <span class="md-stat-unit">บาท/เดือน</span></div>
              </div>
            </div>

            <div class="md-card md-stat-card tone-blue">
              <div class="md-stat-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 3H3v18h18V3z"></path><path d="M3 9h18M9 21V9"></path></svg>
              </div>
              <div class="md-stat-body">
                <div class="md-stat-label">ขนาดพื้นที่</div>
                <div class="md-stat-value">{{ areaRange }} <span class="md-stat-unit">ตร.ม.</span></div>
              </div>
            </div>

            <div :class="['md-card', 'md-stat-card', `tone-${availLevel}`]">
              <div class="md-stat-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
              </div>
              <div class="md-stat-body">
                <div class="md-stat-label">แผงว่างพร้อมให้เช่า</div>
                <div class="md-stat-value">
                  <template v-if="availNow > 0">
                    <span class="md-avail-num">{{ availNow.toLocaleString() }}</span>
                    <span class="md-stat-unit">/ {{ totalStalls.toLocaleString() }} แผง</span>
                  </template>
                  <span v-else class="md-avail-num">เต็ม</span>
                </div>
                <div class="md-avail-track">
                  <span class="md-avail-bar" :style="{ width: availBarWidth }"></span>
                </div>
              </div>
            </div>

            <div class="md-card md-stat-card tone-slate">
              <div class="md-stat-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 17h14M5 17a2 2 0 1 0 4 0M15 17a2 2 0 1 0 4 0"></path><path d="M3 17V11l2-5h14l2 5v6"></path><line x1="3" y1="11" x2="21" y2="11"></line></svg>
              </div>
              <div class="md-stat-body">
                <div class="md-stat-label">ที่จอดรถ</div>
                <div class="md-stat-value">{{ parkingCount }} <span class="md-stat-unit">คัน</span></div>
              </div>
            </div>
          </div>

          <!-- Market Detail Card -->
          <div id="md-detail" class="md-card md-detail-card">
            <h3 class="md-card-title">รายละเอียดตลาด</h3>
            <p class="md-detail-text md-detail-lead">
              {{ currentMarket.description || 'ตลาดมาตรฐานภายใต้การกำกับดูแลของสำนักงานตลาดกรุงเทพมหานคร ได้รับการรับรองด้านสุขอนามัย ความปลอดภัย และการจัดระเบียบแผงค้าที่เป็นธรรม' }}
            </p>

            <div class="md-detail-body">
              <div class="md-detail-row">
                <span class="md-detail-icon">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
                </span>
                <div class="md-detail-content">
                  <h4 class="md-sub-title">ที่อยู่</h4>
                  <p class="md-detail-text">{{ currentMarket.address }}</p>
                </div>
              </div>

              <div class="md-detail-row">
                <span class="md-detail-icon">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                </span>
                <div class="md-detail-content">
                  <h4 class="md-sub-title">วันและเวลาเปิดทำการ</h4>
                  <p class="md-detail-text">{{ operatingDays }} {{ operatingTime }}</p>
                </div>
              </div>

              <div class="md-detail-row">
                <span class="md-detail-icon">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="3 11 22 2 13 21 11 13 3 11"></polygon></svg>
                </span>
                <div class="md-detail-content">
                  <h4 class="md-sub-title">การเดินทาง</h4>
                  <div class="md-travel-list">
                    <span v-for="item in travelOptions" :key="item" class="md-travel-chip">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="3" width="16" height="15" rx="2"></rect><line x1="4" y1="10" x2="20" y2="10"></line><line x1="7" y1="21" x2="7" y2="18"></line><line x1="17" y1="21" x2="17" y2="18"></line></svg>
                      {{ item }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <aside class="md-main-right">
          <!-- More Photos -->
          <div class="md-card md-media-card">
            <h3 class="md-card-title">รูปภาพเพิ่มเติม</h3>
            <div class="md-media-frame">
              <img :src="currentMarket.image_url || '/images/market_hall_front.jpg'" :alt="currentMarket.name" />
            </div>
          </div>

          <!-- Intro Video -->
          <div class="md-card md-media-card">
            <h3 class="md-card-title">วิดีโอแนะนำตลาด</h3>
            <a
              id="md-video-link"
              class="md-media-frame md-video-frame"
              :href="currentMarket.video_url || `https://www.youtube.com/results?search_query=${encodeURIComponent(currentMarket.name)}`"
              target="_blank"
              rel="noopener"
              title="ดูวิดีโอแนะนำตลาดบน YouTube"
            >
              <img :src="currentMarket.image_url || '/images/market_hall_front.jpg'" :alt="`วิดีโอแนะนำ${currentMarket.name}`" />
              <span class="md-video-overlay"></span>
              <span class="md-play-btn">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor" stroke="none"><polygon points="8 5 19 12 8 19 8 5"></polygon></svg>
              </span>
              <span class="md-video-caption">{{ currentMarket.name }}</span>
            </a>
          </div>
        </aside>
      </div>

      <!-- 4. Map Card (เต็มความกว้าง) -->
      <section id="md-map" class="md-card md-map-card">
        <div class="md-map-header">
          <h3 class="md-card-title" style="margin-bottom: 0;">แผนที่ตลาด</h3>

          <div class="md-map-tools">
            <div class="gis-mini-mode">
              <button
                type="button"
                :class="['gis-mini-tab', { active: currentMapMode === 'street' }]"
                @click="setMapMode('street')"
                title="มุมมองแผนที่ถนน"
              >
                <span>แผนที่</span>
              </button>
              <button
                type="button"
                :class="['gis-mini-tab', { active: currentMapMode === 'satellite' }]"
                @click="setMapMode('satellite')"
                title="มุมมองภาพถ่ายดาวเทียม"
              >
                <span>ดาวเทียม</span>
              </button>
            </div>

            <a
              id="md-google-maps-link"
              :href="`https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(currentMarket.name + ' ' + (currentMarket.district || 'กรุงเทพ'))}`"
              target="_blank"
              rel="noopener"
              class="md-map-external"
            >
              <span>เปิด Google Maps</span>
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
            </a>
          </div>
        </div>

        <div class="md-map-frame">
          <div id="market-overview-map" style="width: 100%; height: 100%;"></div>
        </div>
      </section>
    </div>

    <!-- Floorplan Image Lightbox Modal -->
    <Teleport to="body">
      <transition name="fade">
        <div
          v-if="isFloorplanLightboxOpen"
          class="md-lightbox-backdrop"
          @click="isFloorplanLightboxOpen = false"
        >
          <div class="md-lightbox-dialog" @click.stop>
            <div class="md-lightbox-header">
              <div class="md-lightbox-title-wrap">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"></polygon>
                  <line x1="8" y1="2" x2="8" y2="18"></line>
                  <line x1="16" y1="6" x2="16" y2="22"></line>
                </svg>
                <strong>ภาพแผนผัง{{ currentMarket.name }} (ภาพรวม)</strong>
              </div>
              <button
                type="button"
                class="md-lightbox-close"
                @click="isFloorplanLightboxOpen = false"
                title="ปิดหน้าต่าง"
                aria-label="ปิด"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>

            <div class="md-lightbox-body">
              <img
                :src="currentMarket.floorplan_image_url || '/images/chatuchak_isometric_plan.jpg'"
                :alt="`ภาพแผนผัง${currentMarket.name}`"
                class="md-lightbox-img"
              />
            </div>

            <div class="md-lightbox-footer">
              <span class="md-lightbox-tip">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="16" x2="12" y2="12"></line>
                  <line x1="12" y1="8" x2="12.01" y2="8"></line>
                </svg>
                <span>คลิกหรือเลื่อนดูผังโซนและทางเข้า-ออกต่างๆ</span>
              </span>
              <div class="md-lightbox-actions">
                <NuxtLink
                  :to="bookLink"
                  class="btn btn-primary btn-sm"
                  style="display: inline-flex; align-items: center; gap: 6px;"
                >
                  <span>{{ isOfficer ? 'จัดการผังแผง GIS' : (availNow <= 0 ? 'เข้าชมผังแผงดิจิทัล' : 'เข้าสู่การจองแผง') }}</span>
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                    <polyline points="12 5 19 12 12 19"></polyline>
                  </svg>
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { addStandardMapControls } from '~/utils/mapControls'
const { showWarning } = useAppAlert()
import { loadMarketGisJson } from '../data/market_loader'
import { getMarketStallSummary, type MarketStallSummary, type StallStatus } from '../data/market_stall_layout'
import { getMarketGisConfig } from '../data/markets_registry'

const route = useRoute()
const { get } = useApi()
const { role } = useAuth()
const { applyMarketOverride, getMarketOverride } = useMarketOverrides()
const { resetTestData } = useResetTestData()
const isOfficer = computed(() => role.value === 'admin' || role.value === 'staff')
const isStaff = computed(() => role.value === 'staff')
const isAdmin = computed(() => role.value === 'admin')

const handleClearAllData = async () => {
  await resetTestData(async () => {
    await loadMarketDetailData(route.query.market)
  })
}

const buildInitialMarket = (queryVal?: any) => {
  const cfg = getMarketGisConfig(String(queryVal || '1'))
  return applyMarketOverride({
    id: cfg.id,
    slug: cfg.slug,
    name: cfg.name,
    name_en: cfg.name,
    image_url: '',
    district: cfg.district,
    address: `${cfg.name} ${cfg.district} กรุงเทพมหานคร`,
    operating_hours: 'ทุกวัน 05:00 - 19:00 น.',
    description: cfg.description,
    tags: [],
    lat: cfg.center[1],
    lng: cfg.center[0]
  })
}

const currentMarket = ref<any>(buildInitialMarket(route.query.market))

const gisStallSummary = ref<MarketStallSummary | null>(null)
const marketFullFromApi = ref(false)
const isFloorplanLightboxOpen = ref(false)
const currentMapMode = ref<'street' | 'satellite'>('street')
let maplibregl: any = null
let map: any = null
let marker: any = null

const marketHeroImage = computed(() => {
  if (currentMarket.value?.image_url) return currentMarket.value.image_url
  const slug = currentMarket.value?.slug
  if (slug) {
    const slugMap: Record<string, string> = {
      'chatuchak': '/images/markets/chatuchak.jpg',
      'pracha-niwet-1': '/images/markets/pracha_niwet.jpg',
      'min-buri': '/images/markets/min_buri.jpg',
      'singhanat': '/images/markets/singhanat.jpg',
      'thewarat': '/images/markets/thewarat.jpg',
      'thonburi-sanamluang-2': '/images/markets/thonburi.jpg',
      'bang-kapi': '/images/markets/bang_kapi.jpg',
      'bang-khae-phirom': '/images/markets/bang_khae.jpg',
      'nong-chok': '/images/markets/nong_chok.jpg',
      'rat-burana': '/images/markets/rat_burana.jpg',
      'ratchadaphisek': '/images/markets/ratchada.jpg',
      'wongwian-lek': '/images/markets/wongwian_lek.jpg'
    }
    if (slugMap[slug]) return slugMap[slug]
  }
  return '/images/market_hall_front.jpg'
})

const streetStyle: any = {
  version: 8,
  sources: {
    'osm-street': {
      type: 'raster',
      tiles: [
        'https://tile.openstreetmap.org/{z}/{x}/{y}.png'
      ],
      tileSize: 256,
      attribution: '&copy; OpenStreetMap contributors',
      maxzoom: 19
    }
  },
  layers: [
    {
      id: 'osm-street-layer',
      type: 'raster',
      source: 'osm-street',
      minzoom: 0,
      maxzoom: 19
    }
  ]
}

const satelliteStyle: any = {
  version: 8,
  sources: {
    'esri-satellite': {
      type: 'raster',
      tiles: [
        'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}'
      ],
      tileSize: 256,
      attribution: 'Tiles &copy; Esri'
    },
    'esri-labels': {
      type: 'raster',
      tiles: [
        'https://server.arcgisonline.com/ArcGIS/rest/services/Reference/World_Boundaries_and_Places/MapServer/tile/{z}/{y}/{x}'
      ],
      tileSize: 256,
      attribution: 'Tiles &copy; Esri'
    }
  },
  layers: [
    {
      id: 'esri-satellite-layer',
      type: 'raster',
      source: 'esri-satellite',
      minzoom: 0,
      maxzoom: 20
    },
    {
      id: 'esri-labels-layer',
      type: 'raster',
      source: 'esri-labels',
      minzoom: 0,
      maxzoom: 20
    }
  ]
}

const setMapMode = (mode: 'street' | 'satellite') => {
  if (currentMapMode.value === mode || !map) return
  currentMapMode.value = mode
  map.setStyle(mode === 'street' ? streetStyle : satelliteStyle)
}

const initMap = () => {
  if (!maplibregl || !document.getElementById('market-overview-map')) return

  const lng = currentMarket.value.lng || 100.5501
  const lat = currentMarket.value.lat || 13.7998

  map = new maplibregl.Map({
    container: 'market-overview-map',
    style: currentMapMode.value === 'street' ? streetStyle : satelliteStyle,
    center: [lng, lat],
    zoom: 15
  })

  addStandardMapControls(maplibregl, map, {
    onFit: () => map.flyTo({ center: [lng, lat], zoom: 15, pitch: 0, bearing: 0, essential: true }),
    fitTitle: 'กลับสู่ตำแหน่งตลาด',
    notify: showWarning
  })

  // Create custom marker pin
  const el = document.createElement('div')
  el.className = 'market-overview-pin'
  el.innerHTML = `
    <div style="background: #2563eb; color: #fff; padding: 6px 12px; border-radius: 20px; font-weight: 700; font-size: 11px; box-shadow: 0 4px 12px rgba(37,99,235,0.4); border: 2px solid #ffffff; display: flex; align-items: center; gap: 5px; cursor: pointer;">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
      <span>${currentMarket.value.name}</span>
    </div>
  `

  const popup = new maplibregl.Popup({ offset: 25, closeButton: false }).setHTML(`
    <div style="padding: 4px; font-family: 'Bai Jamjuree', sans-serif;">
      <strong style="color: #0f172a; font-size: 0.9rem; display: block; margin-bottom: 2px;">${currentMarket.value.name}</strong>
      <span style="font-size: 0.78rem; color: #64748b;">${currentMarket.value.address}</span>
      <div style="margin-top: 6px;">
        <button type="button" onclick="window.$navigateToMarket('${currentMarket.value.id}', 'book')" style="background: none; border: none; padding: 0; color: #2563eb; font-weight: 700; font-size: 0.8rem; cursor: pointer; font-family: 'Bai Jamjuree', sans-serif;">จองแผงเช่าตลาดนี้ &rarr;</button>
      </div>
    </div>
  `)

  marker = new maplibregl.Marker({ element: el, anchor: 'bottom' })
    .setLngLat([lng, lat])
    .setPopup(popup)
    .addTo(map)

  // Open popup by default
  popup.addTo(map)
}

// ---------- Display helpers (หน้ารายละเอียดตลาด) ----------
const bookLink = computed(() => isOfficer.value
  ? `/stalls/manage?market=${currentMarket.value.id || 1}`
  : `/stalls/book?market=${currentMarket.value.id || 1}`)
const gisManageLink = computed(() => `/stalls/manage?market=${currentMarket.value.id || route.query.market || 1}`)
const visitScheduleLink = computed(() => `/stalls/visit-schedule?market=${currentMarket.value.id || route.query.market || 1}`)
const editMarketLink = computed(() => `/markets/${currentMarket.value.id || route.query.market || 1}/edit`)

// แยก "เสาร์ - อาทิตย์ 09:00 - 18:00 น." เป็นวัน / เวลา
const operatingDays = computed(() => {
  const h = currentMarket.value.operating_hours || 'ทุกวัน 05:00 - 19:00 น.'
  const idx = h.search(/\d/)
  return idx > 0 ? h.slice(0, idx).trim() : 'ทุกวัน'
})

const operatingTime = computed(() => {
  const h = currentMarket.value.operating_hours || 'ทุกวัน 05:00 - 19:00 น.'
  const idx = h.search(/\d/)
  return idx >= 0 ? h.slice(idx).trim() : h
})

const priceRange = computed(() => {
  const saved = getMarketOverride(currentMarket.value?.id) || getMarketOverride(currentMarket.value?.slug)
  const min = Number(saved?.starting_price_monthly ?? gisStallSummary.value?.minRent ?? currentMarket.value.starting_price_monthly ?? 1800)
  const max = Number(saved?.max_price_monthly ?? gisStallSummary.value?.maxRent ?? currentMarket.value.max_price_monthly ?? min)
  return min === max ? `${min.toLocaleString()}` : `${min.toLocaleString()} - ${max.toLocaleString()}`
})

const areaRange = computed(() => {
  const saved = getMarketOverride(currentMarket.value?.id) || getMarketOverride(currentMarket.value?.slug)
  const min = Number(saved?.min_area_sqm ?? gisStallSummary.value?.minAreaSqm ?? currentMarket.value.min_area_sqm ?? 4)
  const max = Number(saved?.max_area_sqm ?? gisStallSummary.value?.maxAreaSqm ?? currentMarket.value.max_area_sqm ?? min)
  const fmt = (val: number) => Number(val.toFixed(1)).toLocaleString()
  return Math.abs(min - max) < 0.05 ? fmt(min) : `${fmt(min)} - ${fmt(max)}`
})

// ดึงจำนวนที่จอดรถจาก tag เช่น "มีที่จอดรถ 800+ คัน"
const parkingCount = computed(() => {
  if (currentMarket.value.parking_spaces !== undefined && currentMarket.value.parking_spaces !== null) {
    return Number(currentMarket.value.parking_spaces).toLocaleString()
  }
  const tag = (currentMarket.value.tags || []).find((t: string) => t.includes('ที่จอดรถ'))
  const m = tag && tag.match(/([\d,]+\+?)/)
  return m ? m[1] : '120+'
})

const travelOptions = computed<string[]>(() => currentMarket.value.travel_options?.length
  ? currentMarket.value.travel_options
  : ['รถโดยสารประจำทาง', 'รถรับจ้างสาธารณะ', 'รถยนต์ส่วนบุคคล', 'รถไฟฟ้า BTS / MRT'])

// แผงว่าง: ใช้เกณฑ์สีเดียวกับการ์ดตลาดหน้าแรก (ok ≥50%, mid 20-49%, low 1-19%, full)
// The API is the source of truth for market-wide counts. GIS derives the
// floorplan total, but it does not contain every booking record as a stall row
// (for example, a newly generated GIS stall), so preferring it here made the
// card stay at the original full count after booking.
const availNow = computed<number>(() => {
  const fromApi = Number(currentMarket.value.available_stalls)
  if (Number.isFinite(fromApi) && fromApi >= 0) return fromApi
  return gisStallSummary.value?.available ?? 0
})
const totalStalls = computed<number>(() => {
  const fromApi = Number(currentMarket.value.total_stalls)
  if (Number.isFinite(fromApi) && fromApi > 0) return fromApi
  return gisStallSummary.value?.total ?? Math.max(180, availNow.value)
})
const availLevel = computed(() => {
  if (availNow.value <= 0) return 'full'
  const pct = (availNow.value / totalStalls.value) * 100
  if (pct >= 50) return 'ok'
  if (pct >= 20) return 'mid'
  return 'low'
})
const availBarWidth = computed(() => {
  if (availNow.value <= 0) return '100%'
  return `${Math.max(2, Math.round((availNow.value / totalStalls.value) * 100))}%`
})

const loadMarketDetailData = async (marketQueryRaw?: any) => {
  const marketQuery = marketQueryRaw ?? route.query.market ?? '1'
  currentMarket.value = buildInitialMarket(marketQuery)
  try {
    const response = await get('/markets')
    const markets = Array.isArray(response) ? response : (response?.data || [])
    const lookup = String(marketQuery || currentMarket.value.id || '1')
    const found = markets.find((m: any) => String(m.id) === lookup || m.slug === lookup)
    if (found) {
      const hasAvailableCount = found.available_stalls !== undefined && found.available_stalls !== null
      marketFullFromApi.value = found.status === 'full'
        || (hasAvailableCount && Number(found.available_stalls) === 0)
      currentMarket.value = applyMarketOverride(found)
    }
  } catch (e) {
    console.warn('Could not fetch market details:', e)
  }

  let backendStalls: any[] = []
  try {
    const marketParam = String(marketQuery || currentMarket.value.slug || currentMarket.value.id || '1')
    const response = await get(`/stalls?market=${encodeURIComponent(marketParam)}`)
    backendStalls = Array.isArray(response) ? response : (response?.data || [])
  } catch (err) {
    console.error('Error fetching stalls:', err)
  }

  try {
    const gisMarket = await loadMarketGisJson(
      String(marketQuery || currentMarket.value.slug || currentMarket.value.id || '1')
    )
    const statusOverrides: Record<string, StallStatus> = {}
    backendStalls.forEach((stall: any) => {
      const ident = String(stall.ident || stall.stall_ident || stall.stall_number || '').trim()
      // `status` is the record lifecycle flag (for example "A"), while
      // `stall_status` is the availability changed by a booking transaction.
      const stallStatus = String(stall.stall_status || '').trim().toLowerCase()
      if (!ident || !['available', 'reserved', 'occupied'].includes(stallStatus)) return
      statusOverrides[ident] = stallStatus as StallStatus
    })
    gisStallSummary.value = getMarketStallSummary(gisMarket.projects, statusOverrides, marketFullFromApi.value)
  } catch (error) {
    console.warn('Could not calculate GIS stall summary:', error)
  }

  if (map) {
    try { map.remove() } catch {}
    map = null
  }
  if (maplibregl) {
    initMap()
  }
}

onMounted(async () => {
  await loadMarketDetailData(route.query.market)

  // Initialize MapLibre GL
  if (process.client) {
    (window as any).$navigateToMarket = (marketId: number | string, action: 'book' | 'stalls') => {
      if (action === 'book') {
        navigateTo(isOfficer.value ? `/stalls/manage?market=${marketId}` : `/stalls/book?market=${marketId}`)
      } else {
        navigateTo(`/stalls?market=${marketId}`)
      }
    }
    const ml = await import('maplibre-gl')
    maplibregl = ml.default || ml
    if (!map) {
      initMap()
    }
  }
})

watch(() => route.query.market, (newMarket) => {
  void loadMarketDetailData(newMarket)
})

onBeforeUnmount(() => {
  if (map) {
    const mapInstance = map
    map = null
    try {
      mapInstance.remove()
    } catch {}
  }
})
</script>

<style scoped>
/* ============================================================== */
/* MARKET DETAIL PAGE                                             */
/* ============================================================== */
.market-detail-page {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.md-breadcrumb-row {
  margin-bottom: -16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.md-breadcrumb-row :deep(.body-breadcrumb-wrapper) {
  margin-bottom: 0;
}

.md-page-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin: 8px 0 -8px;
}

.md-page-heading h2 {
  margin: 0;
  color: #1e293b;
  font-size: 1.5rem;
  font-weight: 800;
  line-height: 1.3;
}

.md-page-heading p {
  margin: 3px 0 0;
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
}

.md-page-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.md-card {
  background: #ffffff;
  border: 1px solid #eef2f6;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
}

.md-card-title {
  font-size: 1rem;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 14px;
}

/* ---------- Header Card (Cinematic Real-Market Photo Backdrop + Luxury Glassmorphism) ---------- */
.md-header-card {
  position: relative;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  justify-content: space-between;
  align-items: center;
  column-gap: clamp(52px, 7vw, 112px);
  padding: 30px 56px;
  background: #0f172a;
  border: 1px solid #eef2f6;
  border-radius: 16px;
  box-shadow: 0 26px 70px -24px rgba(15, 23, 42, 0.24), 0 10px 28px -20px rgba(37, 99, 235, 0.28);
  overflow: hidden;
  min-height: 310px;
  isolation: isolate;
}

/* Background photo layer */
.md-header-bg-frame {
  position: absolute;
  inset: 0;
  z-index: 0;
  overflow: hidden;
}

.md-header-bg-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center 35%;
  transform: scale(1.01);
  transition: transform 1.2s cubic-bezier(0.16, 1, 0.3, 1), filter 0.8s ease;
  filter: brightness(0.92) saturate(1.18) contrast(1.04);
}

.md-header-card:hover .md-header-bg-img {
  transform: scale(1.04);
}

/* Multi-layered Frosted Glass Gradient Overlay */
.md-header-overlay-gradient {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg,
    rgba(248, 251, 255, 0.99) 0%,
    rgba(248, 251, 255, 0.97) 42%,
    rgba(239, 246, 255, 0.91) 68%,
    rgba(239, 246, 255, 0.82) 100%
  );
  backdrop-filter: blur(1px);
  pointer-events: none;
}

.md-header-overlay-mesh {
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at 10% 12%, rgba(147, 197, 253, 0.34) 0%, transparent 42%),
              radial-gradient(circle at 72% 115%, rgba(37, 99, 235, 0.2) 0%, transparent 42%),
              linear-gradient(180deg, transparent 58%, rgba(255, 255, 255, 0.68) 100%);
  pointer-events: none;
}

.md-header-ambient-glow {
  position: absolute;
  top: -80px;
  right: 15%;
  width: 380px;
  height: 380px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(96, 165, 250, 0.28) 0%, transparent 70%);
  pointer-events: none;
}

.md-header-light-beam {
  position: absolute;
  z-index: 0;
  top: -35%;
  left: 44%;
  width: 18%;
  height: 180%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.32), transparent);
  transform: rotate(13deg);
  pointer-events: none;
}

.md-header-grid-pattern {
  position: absolute;
  z-index: 0;
  inset: 0;
  opacity: 0.18;
  background-image: linear-gradient(rgba(37, 99, 235, 0.12) 1px, transparent 1px),
                    linear-gradient(90deg, rgba(37, 99, 235, 0.12) 1px, transparent 1px);
  background-size: 34px 34px;
  mask-image: linear-gradient(90deg, #000 0%, transparent 58%);
  pointer-events: none;
}

.md-header-left {
  min-width: 0;
  position: relative;
  z-index: 1;
  max-width: 700px;
}

.md-header-eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.md-eyebrow-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 999px;
  background: rgba(239, 246, 255, 0.95);
  border: 1px solid #bfdbfe;
  color: #1d4ed8;
  font-size: 0.74rem;
  font-weight: 700;
  letter-spacing: 0.02em;
  box-shadow: 0 1px 3px rgba(37, 99, 235, 0.08);
}

.md-eyebrow-pill svg {
  color: #2563eb;
}

.md-title-wrap {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 10px;
}

.md-title {
  font-size: 2.15rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.02em;
  line-height: 1.2;
  margin: 0;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.9);
}

.md-hero-summary {
  max-width: 680px;
  margin: 0 0 20px;
  color: #475569;
  font-size: 0.94rem;
  font-weight: 500;
  line-height: 1.7;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.md-title-accent,
.md-section-accent {
  display: block;
  width: 44px;
  height: 4px;
  border-radius: 999px;
  background: #2563eb;
  margin-top: 8px;
}

/* Metadata Bar - Modern Minimalist Style */
.md-meta-row {
  display: inline-flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 24px;
  padding: 16px 28px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  margin: 0 0 16px 0;
  transition: all 0.2s ease;
}

.md-meta-row:hover {
  border-color: #cbd5e1;
  background: #f1f5f9;
}

.md-meta-item {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.md-meta-icon {
  width: 32px;
  height: 32px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #2563eb;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
  transition: transform 0.2s ease, border-color 0.2s ease;
}

.md-meta-item:hover .md-meta-icon {
  transform: scale(1.05);
  border-color: #93c5fd;
}

.md-meta-content {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
}

.md-meta-label {
  font-size: 0.68rem;
  color: #64748b;
  font-weight: 500;
  letter-spacing: 0.01em;
}

.md-meta-val {
  font-size: 0.88rem;
  color: #0f172a;
  font-weight: 700;
}

.md-meta-divider {
  width: 1px;
  height: 24px;
  margin: 0 2px;
  background: #e2e8f0;
}

.md-btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 9px 16px;
  border-radius: 10px;
  background: #2563eb;
  color: #ffffff;
  font-size: 0.86rem;
  font-weight: 600;
  text-decoration: none;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.2);
  transition: background 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease;
}

.md-btn-primary:hover {
  background: #1d4ed8;
  transform: translateY(-1px);
  box-shadow: 0 6px 14px rgba(29, 78, 216, 0.28);
}

.md-btn-lg {
  padding: 12px 24px;
  font-size: 0.92rem;
}

/* Right Side: Quick Stall Status & Booking Card */
.md-header-right {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: center;
  gap: 10px;
  position: relative;
  z-index: 1;
  padding: 0;
}

.md-booking-quick-card {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.md-bqc-status-row {
  grid-column: 1 / -1;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 8px;
  margin: 0 0 2px;
  padding-bottom: 10px;
  border-bottom: 1px solid rgba(203, 213, 225, 0.58);
}

.md-bqc-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.72rem;
  font-weight: 700;
  padding: 0;
  background: transparent;
  border: 0;
  border-radius: 0;
  box-shadow: none;
  letter-spacing: 0.01em;
}

.md-bqc-badge.is-available {
  color: #059669;
}

.md-bqc-badge.is-full {
  color: #dc2626;
}

.md-bqc-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
}

.md-bqc-dot.pulsing {
  box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7);
  animation: bqc-pulse 2s infinite;
}

.md-bqc-simple-summary {
  position: relative;
  display: grid;
  grid-template-columns: 120px minmax(0, 1fr);
  align-items: center;
  gap: 14px 18px;
  min-height: 0;
  padding: 16px 18px 18px 21px;
  border: 1px solid rgba(255, 255, 255, 0.88);
  border-radius: 18px;
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.86), rgba(239, 246, 255, 0.72));
  box-shadow: 0 18px 40px -30px rgba(15, 23, 42, 0.66), 0 12px 28px -25px rgba(37, 99, 235, 0.62), inset 0 1px 0 rgba(255, 255, 255, 0.94);
  backdrop-filter: blur(14px) saturate(1.14);
  box-sizing: border-box;
  transition: transform 0.24s ease, border-color 0.24s ease, box-shadow 0.24s ease;
}

.md-bqc-simple-summary:hover {
  transform: translateY(-1px);
  border-color: rgba(147, 197, 253, 0.92);
  box-shadow: 0 24px 50px -28px rgba(15, 23, 42, 0.66), 0 18px 36px -24px rgba(37, 99, 235, 0.72), inset 0 1px 0 #ffffff;
}

.md-bqc-summary-copy {
  position: relative;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding: 2px 0 2px 17px;
}

.md-bqc-summary-copy::before {
  content: '';
  position: absolute;
  top: 3px;
  bottom: 3px;
  left: 0;
  width: 2px;
  border-radius: 999px;
  background: linear-gradient(180deg, rgba(59, 130, 246, 0.24), rgba(16, 185, 129, 0.68));
}

.md-bqc-copy-title {
  display: flex;
  align-items: baseline;
  flex-wrap: wrap;
  column-gap: 5px;
  row-gap: 0;
  color: #1e293b;
  font-size: 0.96rem;
  font-weight: 700;
  line-height: 1.3;
  letter-spacing: -0.01em;
}

.md-bqc-copy-title strong {
  font: inherit;
  font-weight: 800;
}

.md-bqc-copy-title strong.is-ready {
  color: #059669;
}

.md-bqc-copy-title strong.is-unavailable {
  color: #dc2626;
}

.md-bqc-summary-copy small {
  color: #718096;
  font-size: 0.7rem;
  font-weight: 600;
  white-space: nowrap;
}

.md-bqc-summary-copy small strong {
  color: #475569;
  font: inherit;
  font-weight: 800;
}

.md-bqc-summary-count {
  width: auto;
  height: auto;
  display: flex;
  align-items: baseline;
  justify-content: flex-start;
  flex-shrink: 0;
  gap: 4px;
  overflow: visible;
}

.md-bqc-summary-count strong {
  display: inline-block;
  font-size: 4.35rem;
  font-weight: 850;
  letter-spacing: -0.045em;
  line-height: 0.9;
  padding-inline: 0.06em;
  margin-inline: -0.06em 0;
  font-variant-numeric: tabular-nums;
  color: #2563eb;
  background: linear-gradient(145deg, #1d4ed8 10%, #3b82f6 55%, #4f46e5 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: none;
}

.md-bqc-summary-count.is-three-digits strong {
  font-size: 3rem;
  letter-spacing: -0.04em;
}

.md-bqc-summary-count.is-four-plus-digits strong {
  font-size: 2.25rem;
  letter-spacing: -0.03em;
}

.md-bqc-summary-count small {
  margin: 0 0 8px 2px;
  font-size: 0.72rem;
  font-weight: 700;
}

.md-bqc-summary-count.color-full {
  color: #94a3b8;
}

.md-bqc-summary-count.color-full strong {
  color: #94a3b8;
  background: none;
  -webkit-text-fill-color: #94a3b8;
}

@keyframes bqc-pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7);
  }
  70% {
    box-shadow: 0 0 0 6px rgba(16, 185, 129, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(16, 185, 129, 0);
  }
}

.md-bqc-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 0.72rem;
  font-weight: 600;
  color: #1e40af;
  background: #eff6ff;
  border: 1px solid #dbeafe;
  padding: 3px 8px;
  border-radius: 6px;
}

.md-bqc-tag svg {
  color: #2563eb;
}

/* Inner Showcase Panel */
.md-bqc-showcase-panel {
  background: linear-gradient(145deg, rgba(248, 250, 252, 0.98), rgba(239, 246, 255, 0.78));
  border: 1px solid #dbeafe;
  border-radius: 18px;
  padding: 17px;
  display: flex;
  flex-direction: column;
  gap: 9px;
}

.md-bqc-metric-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.md-bqc-metric-title {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.78rem;
  color: #64748b;
  font-weight: 600;
}

.md-bqc-metric-title svg {
  color: #3b82f6;
}

.md-bqc-count-wrap {
  display: flex;
  align-items: baseline;
  gap: 3px;
}

.md-bqc-count {
  font-size: 1.75rem;
  font-weight: 800;
  line-height: 1;
}

.md-bqc-count.color-available {
  color: #0f172a;
}

.md-bqc-count.color-full {
  color: #94a3b8;
}

.md-bqc-total {
  font-size: 0.76rem;
  color: #94a3b8;
  font-weight: 600;
}

.md-bqc-progress-track {
  width: 100%;
  height: 8px;
  background: #e2e8f0;
  border-radius: 999px;
  overflow: hidden;
}

.md-bqc-progress-fill {
  height: 100%;
  border-radius: 999px;
  transition: width 0.4s ease;
}

.md-bqc-progress-fill.tone-ok {
  background: linear-gradient(90deg, #10b981 0%, #059669 100%);
}

.md-bqc-progress-fill.tone-mid {
  background: linear-gradient(90deg, #3b82f6 0%, #2563eb 100%);
}

.md-bqc-progress-fill.tone-low {
  background: linear-gradient(90deg, #f59e0b 0%, #d97706 100%);
}

.md-bqc-progress-fill.tone-full {
  background: #94a3b8;
}

.md-bqc-progress-meta {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.md-bqc-urgency-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.7rem;
  font-weight: 700;
  color: #b45309;
  background: #fef3c7;
  padding: 1px 7px;
  border-radius: 4px;
}

.md-bqc-urgency-badge svg {
  color: #d97706;
}

.md-bqc-ratio-text,
.md-bqc-full-text {
  font-size: 0.7rem;
  color: #64748b;
  font-weight: 500;
}

.md-bqc-panel-divider {
  height: 1px;
  background: #e2e8f0;
  margin: 2px 0;
}

.md-bqc-price-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.md-bqc-price-box {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.md-bqc-price-label {
  font-size: 0.68rem;
  color: #64748b;
  font-weight: 500;
}

.md-bqc-price-amount {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.md-bqc-price-amount strong {
  font-size: 1.22rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1;
}

.md-bqc-price-amount small {
  font-size: 0.72rem;
  font-weight: 600;
  color: #64748b;
}

.md-bqc-trust-seal {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.7rem;
  font-weight: 600;
  color: #15803d;
  background: #dcfce7;
  border: 1px solid #bbf7d0;
  padding: 2.5px 8px;
  border-radius: 999px;
  white-space: nowrap;
}

.md-bqc-trust-seal svg {
  color: #16a34a;
}

/* Full-Width Action Button */
.md-bqc-action-btn {
  position: relative;
  overflow: hidden;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 30px;
  align-items: center;
  width: 100%;
  height: 46px;
  padding: 0 10px 0 15px;
  border-radius: 14px;
  font-size: 0.86rem;
  font-weight: 700;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.22s cubic-bezier(0.4, 0, 0.2, 1);
  box-sizing: border-box;
}

.md-bqc-action-btn::before {
  content: '';
  position: absolute;
  top: -45%;
  left: -34%;
  width: 22%;
  height: 190%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.42), transparent);
  transform: rotate(18deg);
  animation: bqc-button-sheen 4.8s ease-in-out infinite;
  pointer-events: none;
}

.md-bqc-action-btn.btn-active {
  background: linear-gradient(105deg, #2563eb 0%, #1d4ed8 58%, #4338ca 100%);
  color: #ffffff;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 14px 26px -12px rgba(37, 99, 235, 0.76), inset 0 1px 0 rgba(255, 255, 255, 0.24);
}

.md-bqc-action-btn.btn-active:hover {
  background: linear-gradient(135deg, #1d4ed8 0%, #1e40af 100%);
  transform: translateY(-2px);
  box-shadow: 0 8px 22px rgba(37, 99, 235, 0.45);
}

.md-bqc-action-btn.btn-active:hover .md-bqc-btn-arrow {
  transform: translateX(4px);
}

.md-bqc-action-btn.btn-active:active {
  transform: scale(0.985);
}

.md-bqc-action-btn.btn-full {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #cbd5e1;
}

.md-bqc-action-btn.btn-full:hover {
  background: #e2e8f0;
  color: #0f172a;
}

.md-bqc-btn-label {
  position: relative;
  z-index: 1;
  grid-column: 1;
  text-align: left;
  padding: 0 14px;
  letter-spacing: 0.01em;
}

.md-bqc-btn-arrow {
  position: relative;
  z-index: 1;
  grid-column: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 30px;
  height: 30px;
  padding: 7px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  box-sizing: border-box;
  transition: transform 0.2s ease;
}

@keyframes bqc-button-sheen {
  0%, 58% {
    left: -34%;
    opacity: 0;
  }
  68% {
    opacity: 1;
  }
  86%, 100% {
    left: 116%;
    opacity: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .md-bqc-dot.pulsing,
  .md-bqc-action-btn::before {
    animation: none;
  }
}

.md-bqc-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  font-size: 0.72rem;
  color: #64748b;
  text-align: center;
}

.md-bqc-hint svg {
  color: #3b82f6;
  flex-shrink: 0;
}

/* ---------- Floorplan ---------- */
.md-floorplan-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 28px 30px;
  scroll-margin-top: 80px;
}

.md-section-title {
  font-size: 1.4rem;
  font-weight: 800;
  color: #0f172a;
  text-align: center;
}

.md-section-accent {
  margin: 8px auto 0;
}

.md-section-sub {
  margin-top: 14px;
  font-size: 0.88rem;
  color: #64748b;
  text-align: center;
}

.md-floorplan-img-link {
  position: relative;
  display: block;
  width: 100%;
  max-width: 760px;
  margin: 18px auto 22px;
  border-radius: 14px;
  overflow: hidden;
}

.md-floorplan-img {
  width: 100%;
  display: block;
  mix-blend-mode: multiply;
  transition: transform 0.4s ease;
}

.md-floorplan-img-link:hover .md-floorplan-img {
  transform: scale(1.02);
}

.md-floorplan-hint {
  position: absolute;
  right: 14px;
  bottom: 14px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid #e2e8f0;
  color: #334155;
  font-size: 0.78rem;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.08);
  opacity: 0;
  transform: translateY(4px);
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.md-floorplan-img-link:hover .md-floorplan-hint {
  opacity: 1;
  transform: translateY(0);
}

/* ---------- Main Grid ---------- */
.md-main-grid {
  display: flex;
  gap: 20px;
  align-items: stretch;
}

.md-main-left {
  flex: 1.85;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.md-main-right {
  flex: 1;
  min-width: 320px;
  max-width: 380px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.md-stats-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.md-stat-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.md-stat-card:hover {
  border-color: #dbeafe;
  box-shadow: 0 6px 18px rgba(15, 23, 42, 0.06);
  transform: translateY(-2px);
}

.md-stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: var(--tone-bg, #eff6ff);
  color: var(--tone-fg, #2563eb);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* โทนสีไอคอนแต่ละการ์ด + สีระดับแผงว่าง */
.tone-green { --tone-bg: #ecfdf5; --tone-fg: #059669; }
.tone-blue  { --tone-bg: #eff6ff; --tone-fg: #2563eb; }
.tone-slate { --tone-bg: #f1f5f9; --tone-fg: #475569; }
.tone-ok    { --tone-bg: #dcfce7; --tone-fg: #15803d; --bar: #22c55e; }
.tone-mid   { --tone-bg: #fef9c3; --tone-fg: #a16207; --bar: #eab308; }
.tone-low   { --tone-bg: #ffedd5; --tone-fg: #c2410c; --bar: #f97316; }
.tone-full  { --tone-bg: #fee2e2; --tone-fg: #dc2626; --bar: #ef4444; }

.md-stat-body {
  flex: 1;
  min-width: 0;
}

.md-stat-label {
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 500;
}

.md-stat-value {
  display: flex;
  align-items: baseline;
  flex-wrap: wrap;
  gap: 6px;
  font-size: 1.45rem;
  font-weight: 800;
  color: #1e293b;
  line-height: 1.25;
  margin-top: 2px;
}

.md-stat-unit {
  font-size: 0.8rem;
  font-weight: 500;
  color: #64748b;
  white-space: nowrap;
}

.md-avail-num {
  color: var(--tone-fg, #1e293b);
}

.md-avail-track {
  height: 6px;
  margin-top: 8px;
  border-radius: 999px;
  background: #f1f5f9;
  overflow: hidden;
}

.md-avail-bar {
  display: block;
  height: 100%;
  border-radius: 999px;
  background: var(--bar, #2563eb);
  transition: width 0.6s ease;
}

.md-detail-card {
  padding: 24px 28px;
  scroll-margin-top: 80px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.md-detail-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.md-detail-text {
  font-size: 0.88rem;
  color: #475569;
  line-height: 1.7;
}

.md-detail-lead {
  margin-bottom: 8px;
}

.md-detail-row {
  display: flex;
  gap: 14px;
  padding: 16px 0;
  border-top: 1px solid #f1f5f9;
}

.md-detail-row:first-of-type {
  margin-top: 12px;
}

.md-detail-row:last-child {
  padding-bottom: 0;
}

.md-detail-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.md-detail-content {
  flex: 1;
  min-width: 0;
}

.md-sub-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 4px;
}

.md-travel-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 6px;
}

.md-travel-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  color: #334155;
  background: #f8fafc;
  font-size: 0.82rem;
  font-weight: 500;
}

.md-travel-chip svg {
  color: #2563eb;
}

/* ---------- Media Cards ---------- */
.md-media-card {
  padding: 20px;
}

.md-media-frame {
  position: relative;
  display: block;
  width: 100%;
  aspect-ratio: 16 / 10;
  border-radius: 12px;
  overflow: hidden;
  background: #f1f5f9;
}

.md-media-frame img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.4s ease;
}

.md-media-frame:hover img {
  transform: scale(1.04);
}

.md-video-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.05) 30%, rgba(15, 23, 42, 0.6) 100%);
}

.md-play-btn {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 56px;
  height: 40px;
  transform: translate(-50%, -50%);
  border-radius: 12px;
  background: #ef4444;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6px 18px rgba(239, 68, 68, 0.45);
  transition: transform 0.2s ease;
}

.md-video-frame:hover .md-play-btn {
  transform: translate(-50%, -50%) scale(1.08);
}

.md-video-caption {
  position: absolute;
  left: 14px;
  bottom: 12px;
  color: #ffffff;
  font-size: 0.85rem;
  font-weight: 600;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
}

/* ---------- Map Card ---------- */
.md-map-card {
  padding: 24px;
  scroll-margin-top: 80px;
}

.md-map-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 14px;
}

.md-map-tools {
  display: flex;
  align-items: center;
  gap: 12px;
}

.md-map-external {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 8px;
  border: 1px solid #93c5fd;
  background: #ffffff;
  font-size: 0.825rem;
  font-weight: 600;
  color: #2563eb;
  text-decoration: none;
  transition: background 0.2s ease;
}

.md-map-external:hover {
  background: #eff6ff;
}

.md-map-frame {
  height: 380px;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.gis-mini-mode {
  display: inline-flex;
  background: #f1f5f9;
  border-radius: 8px;
  padding: 2px;
  border: 1px solid #e2e8f0;
}

.gis-mini-tab {
  padding: 4px 12px;
  font-size: 0.78rem;
  font-weight: 600;
  border: none;
  background: transparent;
  color: #64748b;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.gis-mini-tab:hover {
  color: #0f172a;
}

.gis-mini-tab.active {
  background: #ffffff;
  color: #2563eb;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  font-weight: 700;
}

/* แจ้งเตือนเมื่อตลาดเต็ม */
.md-full-market-alert {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
  padding: 12px 18px;
  border-radius: 12px;
  margin: 16px 0 10px;
  max-width: 760px;
  width: 100%;
}

.md-full-market-alert svg {
  flex-shrink: 0;
  color: #dc2626;
}

.md-full-alert-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 0.84rem;
  line-height: 1.5;
}

.md-full-alert-content strong {
  font-weight: 700;
  color: #b91c1c;
}

.md-btn-full {
  background: #475569 !important;
  box-shadow: 0 2px 8px rgba(71, 85, 105, 0.25) !important;
}

.md-btn-full:hover {
  background: #334155 !important;
}

/* ---------- Responsive ---------- */
@media (max-width: 1180px) {
  .md-header-card {
    column-gap: 36px;
    padding: 30px;
  }

  .md-title {
    font-size: clamp(1.9rem, 3vw, 2.15rem);
  }

  .md-booking-quick-card {
    width: 310px;
  }

  .md-meta-item {
    padding-right: 10px;
  }
}

@media (max-width: 960px) {
  .md-header-card {
    display: flex;
    flex-direction: column;
    align-items: stretch;
    gap: 28px;
  }

  .md-header-left {
    max-width: none;
  }

  .md-header-right {
    width: 100%;
    flex-direction: column;
    align-items: center;
  }

  .md-booking-quick-card {
    width: min(100%, 320px);
  }

  .md-main-grid {
    flex-direction: column;
  }
  .md-main-right {
    max-width: 100%;
    min-width: 0;
  }
  .md-detail-card {
    flex: none;
  }
}

@media (max-width: 768px) {
  .md-page-heading {
    align-items: flex-start;
    flex-direction: column;
    gap: 12px;
  }

  .md-header-card {
    align-items: flex-start;
    padding: 24px 20px;
    gap: 22px;
    min-height: 0;
    border-radius: 16px;
  }
  .md-header-overlay-gradient {
    background: linear-gradient(
      180deg,
      rgba(255, 255, 255, 0.98) 0%,
      rgba(255, 255, 255, 0.94) 60%,
      rgba(255, 255, 255, 0.88) 100%
    );
  }
  .md-header-right {
    width: 100%;
    justify-content: center;
  }
  .md-booking-quick-card {
    width: 100%;
    max-width: 320px;
    margin: 0 auto;
    padding: 0;
  }
  .md-bqc-simple-summary {
    grid-template-columns: 116px minmax(0, 1fr);
    gap: 14px;
  }
  .md-bqc-summary-count {
    width: auto;
    height: auto;
  }
  .md-bqc-summary-count strong {
    font-size: 4rem;
  }
  .md-bqc-summary-count.is-three-digits strong {
    font-size: 2.9rem;
  }
  .md-bqc-summary-count.is-four-plus-digits strong {
    font-size: 2.15rem;
  }
  .md-meta-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
    width: 100%;
    padding: 12px 14px;
  }
  .md-meta-item {
    width: auto;
  }
  .md-meta-divider {
    display: none;
  }
  .md-title {
    font-size: 1.65rem;
  }
  .md-hero-summary {
    font-size: 0.84rem;
    line-height: 1.6;
  }
  .md-stats-grid {
    grid-template-columns: 1fr;
  }
  .md-map-frame {
    height: 300px;
  }
}

/* ---------- Floorplan Lightbox Modal ---------- */
.md-lightbox-backdrop {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: rgba(15, 23, 42, 0.78);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.md-lightbox-dialog {
  background: #ffffff;
  border-radius: 18px;
  max-width: 1100px;
  width: 100%;
  max-height: 92vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.35);
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.md-lightbox-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 22px;
  border-bottom: 1px solid #f1f5f9;
  background: #ffffff;
}

.md-lightbox-title-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #0f172a;
  font-size: 1.05rem;
}

.md-lightbox-title-wrap svg {
  color: #2563eb;
}

.md-lightbox-close {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.md-lightbox-close:hover {
  background: #fee2e2;
  color: #dc2626;
  border-color: #fca5a5;
}

.md-lightbox-body {
  padding: 16px;
  overflow: auto;
  background: #faf8f5;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 380px;
}

.md-lightbox-img {
  width: 100%;
  max-width: 100%;
  height: auto;
  border-radius: 10px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
}

.md-lightbox-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 22px;
  border-top: 1px solid #f1f5f9;
  background: #ffffff;
  flex-wrap: wrap;
  gap: 12px;
}

.md-lightbox-tip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8rem;
  color: #64748b;
}

.md-lightbox-tip svg {
  color: #3b82f6;
}

.md-lightbox-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}
</style>
