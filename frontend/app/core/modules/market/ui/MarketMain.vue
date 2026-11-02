<template>
  <div class="homepage-container">
    <!-- Top Header Navigation Bar (Clean & Streamlined) -->
    <AppHeader />

    <div class="content-body">
      <!-- ============================================================== -->
      <!-- 1. TOP HERO PROMOTIONAL BANNER (Modern High-Impact Banner)     -->
      <!-- ============================================================== -->
      <div
        class="hero-banner-wrapper"
        @mouseenter="pauseHeroCarousel"
        @mouseleave="resumeHeroCarousel"
        @focusin="pauseHeroCarousel"
        @focusout="resumeHeroCarousel"
      >
        <div class="hero-banner-card card">
          <!-- Left Side: Market Photo with Diagonal Red Discount Ribbon -->
          <div class="hero-left-media">
            <img
              :key="heroSlides[currentHeroSlide - 1].image"
              :src="heroSlides[currentHeroSlide - 1].image"
              :alt="heroSlides[currentHeroSlide - 1].alt"
              class="hero-media-img"
            />
            <div class="hero-media-overlay"></div>
            <!-- Optional creative badge configured per advertisement -->
            <div v-if="heroSlide.ribbon" class="hero-discount-ribbon" :class="`hero-ribbon-${heroSlide.ribbon.tone}`">
              <div class="ribbon-text-top">{{ heroSlide.ribbon.top }}</div>
              <div class="ribbon-text-main">{{ heroSlide.ribbon.main }} <span class="ribbon-off">{{ heroSlide.ribbon.suffix }}</span></div>
            </div>
          </div>

          <!-- Right Side: Clean Callout Content -->
          <div class="hero-right-content">
            <div class="campaign-tag-pill" :class="`campaign-tag-${heroSlide.tone}`">
              <span class="campaign-dot"></span>
              <span>{{ heroSlide.tag }}</span>
            </div>
            <div class="hero-subheading">{{ heroSlide.subheading }}</div>
            <h1 class="hero-main-heading">
              {{ heroSlide.title }} <span class="heading-gradient">{{ heroSlide.highlight }}</span>
            </h1>
            <p class="hero-description">{{ heroSlide.description }}</p>

            <!-- Value Proposition Bullets -->
            <div class="hero-features-list">
              <div v-for="feature in heroSlide.features" :key="feature" class="hero-feature-item">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span>{{ feature }}</span>
              </div>
            </div>

            <!-- CTA Buttons Row -->
            <div class="hero-cta-row">
              <button class="btn btn-primary btn-hero-cta" @click="scrollToMarkets">
                <span>{{ heroSlide.cta }}</span>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
              </button>
              <div class="hero-note">
                {{ heroSlide.note }}
              </div>
            </div>

            <!-- Decorative Navy Accent Shape at bottom-right -->
            <div class="hero-accent-shape"></div>
          </div>
        </div>

        <!-- Carousel Slide Indicators (Blue pill + 3 Gray dots) -->
        <div class="hero-carousel-dots">
          <button
            v-for="i in 4"
            :key="i"
            type="button"
            class="carousel-dot"
            :class="{ active: currentHeroSlide === i }"
            @click="currentHeroSlide = i"
            :aria-label="`แสดงภาพสไลด์ที่ ${i}`"
          ></button>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- 2. SECTION: อีเวนต์ตลาด (Market Events)                         -->
      <!-- ============================================================== -->
      <section class="events-card-wrapper">
        <header class="events-promo-header">
          <div class="events-promo-copy">
            <div class="events-badge-pill">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
              <span>อีเวนต์ที่เปิดรับเร็ว ๆ นี้</span>
            </div>
            <h2 class="events-main-title">อีเวนต์ตลาด</h2>
            <p class="events-desc-text">ค้นหากิจกรรมที่เหมาะกับสินค้า ตรวจสอบวัน เวลา และจำนวนแผงว่างก่อนจองพื้นที่</p>
            <button type="button" class="events-view-all" @click="eventPeriod = 'month'">
              <span>ดูกิจกรรมทั้งหมด</span>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
            </button>
          </div>
          <img src="/images/megaphone_3d_transparent.png?v=bma-theme" alt="ประชาสัมพันธ์อีเวนต์ตลาด" class="events-promo-image" />
        </header>

        <div class="events-list-panel">
          <AppSecondaryTabs v-model="eventPeriod" :options="eventPeriodTabs" aria-label="เลือกช่วงเวลาอีเวนต์" />

          <div class="events-list">
            <article v-for="event in filteredMarketEvents" :key="event.name" class="event-list-row">
              <img :src="event.image" :alt="event.name" class="event-list-image" />
              <div class="event-list-content">
                <h3>{{ event.name }}</h3>
                <div class="event-list-meta">
                  <span><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>{{ event.date }}</span>
                  <span class="event-time-pill"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="9"></circle><polyline points="12 7 12 12 16 14"></polyline></svg>{{ event.time }}</span>
                </div>
                <div class="event-location-line">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
                  <span>{{ event.location }}</span>
                </div>
              </div>
              <div class="event-list-action">
                <div class="event-capacity" :class="{ 'is-full': event.available === 0, 'is-low': event.available > 0 && event.available <= 3 }">
                  <strong>{{ event.available }}</strong>
                  <span>{{ event.available === 0 ? 'เต็ม' : event.available <= 3 ? 'ใกล้เต็ม' : 'ว่าง' }}</span>
                </div>
                <NuxtLink
                  v-if="isOfficer || event.available > 0"
                  :to="isOfficer ? `/stalls/manage?market=${event.marketId}` : `/stalls/book?market=${event.marketId}`"
                  class="event-primary-link"
                >
                  <span>{{ isOfficer ? 'จัดการพื้นที่' : 'จองลานกิจกรรม' }}</span>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
                </NuxtLink>
                <span v-else class="event-disabled-link">ปิดรับจอง</span>
              </div>
            </article>
          </div>
        </div>
      </section>

      <!-- ============================================================== -->
      <!-- 3. SECTION: 12 ตลาด (The 12 Markets Directory)                 -->
      <!-- ============================================================== -->
      <div id="markets-12-directory" class="markets-directory-wrap">
        <!-- Center Header: Title & Subtitle -->
        <div class="directory-title-header">
          <div class="directory-header-pill">
            <span>สำนักงานตลาดกรุงเทพมหานคร</span>
          </div>
          <h2 class="directory-big-heading">
            <span class="heading-underline">12 ตลาดในสังกัด กทม.</span>
          </h2>
          <p class="directory-sub-caption">
            ค้นหาข้อมูล ทำเล และยื่นขอจองแผงค้าออนไลน์ได้ทันที พร้อมระบบผังแผงดิจิทัล 3 มิติ และสัญญามาตรฐาน
          </p>
        </div>

        <!-- Search Bar Centered -->
        <div class="search-bar-centered">
          <div class="search-input-pill-wrap">
            <span class="search-icon-left">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
            </span>
            <input
              v-model="searchQuery"
              type="text"
              class="search-input-field"
              placeholder="ค้นหาตลาดที่ต้องการจอง, ชื่อเขต, หรือทำเล (เช่น จตุจักร, มีนบุรี, ธนบุรี)..."
              @input="applyMarketSearch"
            />
            <button
              v-if="searchQuery"
              type="button"
              class="search-clear-btn"
              @click="clearSearch"
              aria-label="ล้างคำค้นหา"
            >
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
            </button>
          </div>
        </div>

        <!-- 5 Category Filter Buttons -->
        <div class="category-icon-buttons-row">
          <!-- Button 0: ทั้งหมด -->
          <div
            class="cat-icon-item"
            :class="{ active: selectedCategory === 'all' }"
            @click="toggleCategory('all')"
          >
            <div class="cat-icon-circle cat-all">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
            </div>
            <div class="cat-label">ทั้งหมด (12 ตลาด)</div>
          </div>

          <!-- Button 1: % โปรโมชั่น -->
          <div
            class="cat-icon-item"
            :class="{ active: selectedCategory === 'promo' }"
            @click="toggleCategory('promo')"
          >
            <div class="cat-icon-circle cat-promo">
              <span class="cat-symbol">%</span>
            </div>
            <div class="cat-label">โปรโมชั่น</div>
          </div>

          <!-- Button 2: ตลาดเปิดใหม่ -->
          <div
            class="cat-icon-item"
            :class="{ active: selectedCategory === 'new' }"
            @click="toggleCategory('new')"
          >
            <div class="cat-icon-circle cat-new">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
            </div>
            <div class="cat-label">แผงว่างเยอะ</div>
          </div>

          <!-- Button 3: ใกล้รถไฟฟ้า -->
          <div
            class="cat-icon-item"
            :class="{ active: selectedCategory === 'metro' }"
            @click="toggleCategory('metro')"
          >
            <div class="cat-icon-circle cat-metro">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="3" width="16" height="16" rx="2"></rect><path d="M4 11h16"></path><path d="M12 3v8"></path><path d="m8 19-2 3"></path><path d="m16 19 2 3"></path><circle cx="8" cy="15" r="1"></circle><circle cx="16" cy="15" r="1"></circle></svg>
            </div>
            <div class="cat-label">ใกล้รถไฟฟ้า</div>
          </div>

          <!-- Button 4: เดินทางสะดวก -->
          <div
            class="cat-icon-item"
            :class="{ active: selectedCategory === 'parking' }"
            @click="toggleCategory('parking')"
          >
            <div class="cat-icon-circle cat-transit">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 17h2V8a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v9h2"></path><circle cx="7" cy="17" r="2"></circle><path d="M9 17h6"></path><circle cx="17" cy="17" r="2"></circle><path d="M5 11h14"></path></svg>
            </div>
            <div class="cat-label">เดินทางสะดวก</div>
          </div>
        </div>

        <!-- Segmented View Mode Toggle: Grid Cards vs Interactive GIS Map -->
        <div class="market-view-switcher-bar">
          <div class="segmented-view-control">
            <button
              class="svc-tab-btn"
              :class="{ active: viewMode === 'grid' }"
              @click="viewMode = 'grid'"
              title="ดูรายการตลาดทั้งหมด"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="3" width="7" height="7"></rect>
                <rect x="14" y="3" width="7" height="7"></rect>
                <rect x="14" y="14" width="7" height="7"></rect>
                <rect x="3" y="14" width="7" height="7"></rect>
              </svg>
              <span>รายการตลาด</span>
              <span v-if="loadingMarkets" class="svc-counter is-loading" aria-label="กำลังโหลดจำนวนตลาด"></span>
              <span v-else class="svc-counter">{{ filteredMarkets.length }}</span>
            </button>

            <button
              class="svc-tab-btn"
              :class="{ active: viewMode === 'map' }"
              @click="viewMode = 'map'"
              title="ดูแผนที่ตำแหน่งตลาดทั้งหมด"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"></polygon>
                <line x1="8" y1="2" x2="8" y2="18"></line>
                <line x1="16" y1="6" x2="16" y2="22"></line>
              </svg>
              <span>แผนที่ตลาด</span>
            </button>
          </div>
        </div>

        <!-- VIEW MODE 1: Interactive GIS Map Section (When Map is Selected - No Cards Below) -->
        <div v-show="viewMode === 'map'" class="market-gis-view-container" style="margin-bottom: 36px;">
          <BangkokMarketsGisMap
            :markets="filteredMarkets"
            :active="viewMode === 'map'"
            @selectMarket="onSelectMarketOnMap"
          />
        </div>

        <!-- VIEW MODE 2: 12 MARKET CARDS GRID (When Grid is Selected) -->
        <template v-if="viewMode === 'grid'">
          <div v-if="showMarketSkeleton" class="market-cards-3col-grid market-skeleton-grid" role="status" aria-live="polite" aria-label="กำลังโหลดรายการตลาด">
            <div v-for="index in 3" :key="index" class="market-box-card market-skeleton-card" aria-hidden="true">
              <div class="market-skeleton-shape market-skeleton-image"></div>
              <div class="market-box-body market-skeleton-body">
                <span class="market-skeleton-shape market-skeleton-title"></span>
                <span class="market-skeleton-shape market-skeleton-address"></span>
                <span class="market-skeleton-shape market-skeleton-address is-short"></span>
                <div class="market-skeleton-availability">
                  <span class="market-skeleton-shape market-skeleton-label"></span>
                  <span class="market-skeleton-shape market-skeleton-count"></span>
                  <span class="market-skeleton-shape market-skeleton-progress"></span>
                </div>
                <div class="market-skeleton-price-row">
                  <span class="market-skeleton-shape market-skeleton-price-label"></span>
                  <span class="market-skeleton-shape market-skeleton-price"></span>
                </div>
                <span class="market-skeleton-shape market-skeleton-button"></span>
              </div>
            </div>
          </div>

          <div v-else class="market-cards-3col-grid">
            <div
              v-for="m in filteredMarkets"
              :key="m.id"
              class="market-box-card"
              :class="{ 'is-officer-card': isOfficer }"
            >
              <!-- Card Top Image with Rounded Corners & Green "ว่าง" Badge -->
              <NuxtLink :to="marketActionLink(m)" class="market-img-container" title="คลิกเพื่อดูภาพแผนผังและข้อมูลตลาด">
                <img :src="m.image_url" :alt="m.name" class="market-box-img" />
                <div class="market-img-gradient"></div>

                <!-- Status Badge (Top-Right Pill: "● ว่าง" or "● เต็ม") -->
                <div
                  class="market-avail-badge"
                  :class="m.status === 'open' && getAvailableCount(m) > 0 ? 'badge-green' : 'badge-red'"
                >
                  <span class="badge-dot-live"></span>
                  <span>{{ m.status === 'open' && getAvailableCount(m) > 0 ? 'ว่าง' : 'เต็ม' }}</span>
                </div>
              </NuxtLink>

              <!-- Card Body Content -->
              <div class="market-box-body">
                <!-- Market Title & District Pill Row -->
                <div class="market-title-views-row">
                  <div class="market-title-col">
                    <NuxtLink :to="marketActionLink(m)" class="market-title-link">
                      <h3 class="market-card-name">{{ m.name }}</h3>
                    </NuxtLink>
                  </div>
                  <span v-if="m.district" class="officer-district-pill">
                    {{ m.district }}
                  </span>
                </div>

                <!-- Address with Map Pin Icon -->
                <div class="market-card-address" :title="m.address">
                  <span class="pin-icon">
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                      <circle cx="12" cy="10" r="3"></circle>
                    </svg>
                  </span>
                  <span class="address-text">{{ m.address }}</span>
                </div>

                <!-- ==================================================== -->
                <!-- OFFICER VIEW: 3-Status Stall Breakdown               -->
                <!-- ==================================================== -->
                <template v-if="isOfficer">
                  <div class="officer-occupancy-box">
                    <div class="officer-occupancy-head">
                      <span class="officer-occupancy-title">
                        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                        <span>สรุปสถานะแผงเช่า</span>
                      </span>
                      <span class="officer-total-stalls">ทั้งหมด <strong>{{ (m.total_stalls || 0).toLocaleString() }}</strong> แผง</span>
                    </div>

                    <div class="officer-stall-stats-grid">
                      <div class="officer-stat-cell stat-occupied">
                        <div class="stat-copy">
                          <span class="stat-label">มีผู้เช่า</span>
                          <strong class="stat-val">{{ getOccupiedCount(m).toLocaleString() }}</strong>
                        </div>
                      </div>
                      <div class="officer-stat-cell stat-reserved">
                        <div class="stat-copy">
                          <span class="stat-label">ติดจอง</span>
                          <strong class="stat-val">{{ getReservedCount(m).toLocaleString() }}</strong>
                        </div>
                      </div>
                      <div class="officer-stat-cell stat-available">
                        <div class="stat-copy">
                          <span class="stat-label">แผงว่าง</span>
                          <strong class="stat-val">{{ getAvailableCount(m).toLocaleString() }}</strong>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Officer Single Action Button -->
                  <div class="market-card-actions">
                    <NuxtLink
                      :to="marketActionLink(m)"
                      class="btn-book-action-bar"
                      title="ดูและจัดการข้อมูลตลาด"
                    >
                      <span>ข้อมูลตลาด</span>
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" class="btn-arrow-icon">
                        <line x1="5" y1="12" x2="19" y2="12"></line>
                        <polyline points="12 5 19 12 12 19"></polyline>
                      </svg>
                    </NuxtLink>
                  </div>
                </template>

                <!-- ==================================================== -->
                <!-- TENANT VIEW: Stall Availability & Starting Price     -->
                <!-- ==================================================== -->
                <template v-else>
                  <!-- Stall Availability Box (มุมมองแผงว่าง) -->
                  <div class="market-occupancy-row">
                    <div class="occupancy-info">
                      <div class="occupancy-title-col">
                        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                        <span class="occupancy-label">แผงว่างพร้อมให้เช่า</span>
                      </div>
                      <span class="occupancy-count">
                        <template v-if="getAvailableCount(m) > 0">
                          <strong :class="getAvailabilityTextClass(m)">{{ getAvailableCount(m).toLocaleString() }}</strong>
                          <span class="occupancy-total">/ {{ (m.total_stalls || 0).toLocaleString() }} แผง</span>
                        </template>
                        <strong v-else class="avail-text-full">เต็ม</strong>
                      </span>
                    </div>

                    <!-- Progress Bar showing availability -->
                    <div class="occupancy-progress-bar" :title="`ว่าง ${getAvailabilityRate(m)}% จากทั้งหมด ${(m.total_stalls || 0).toLocaleString()} แผง`">
                      <div
                        class="occupancy-progress-fill"
                        :class="getAvailabilityBarColor(m)"
                        :style="{ width: `${getAvailabilityBarWidth(m)}%` }"
                      ></div>
                    </div>
                  </div>

                  <!-- Starting Price -->
                  <div class="market-price-row">
                    <span class="price-prefix-text">เริ่มต้น</span>
                    <span class="price-currency">฿</span>
                    <span class="price-big-number">{{ getStartingPrice(m) }}</span>
                    <span class="price-suffix-text">/ {{ getPriceUnit(m) }}</span>
                  </div>

                  <!-- Action Button: Full-Width Primary "จองแผงเช่า →" -->
                  <div class="market-card-actions">
                    <NuxtLink
                      :to="marketActionLink(m)"
                      class="btn-book-action-bar"
                      title="เข้าสู่ขั้นตอนการจองแผงเช่าทันที"
                    >
                      <span>จองแผงเช่า</span>
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" class="btn-arrow-icon">
                        <line x1="5" y1="12" x2="19" y2="12"></line>
                        <polyline points="12 5 19 12 12 19"></polyline>
                      </svg>
                    </NuxtLink>
                  </div>
                </template>
              </div>
            </div>
          </div>
        </template>

        <!-- Empty Results Message -->
        <div v-if="filteredMarkets.length === 0 && !loadingMarkets" class="card empty-search-card">
          <AppEmptyState
            variant="card"
            icon="search"
            title="ไม่พบตลาดตามเงื่อนไข"
            description="ลองเปลี่ยนคำค้นหาหรือล้างตัวกรองเพื่อดูตลาดทั้งหมด"
            action-label="แสดงตลาดทั้งหมด"
            action-variant="primary"
            @action="resetFilters"
          />
        </div>
      </div>
    </div>

    <!-- Official BMA e-Market Footer -->
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
import AppFooter from '~/components/AppFooter.vue'
import { onBeforeRouteLeave } from 'vue-router'

const { get, apiFetch } = useApi()
const { role } = useAuth()
const { applyMarketOverride } = useMarketOverrides()
const route = useRoute()
const isOfficer = computed(() => role.value === 'admin' || role.value === 'staff')

const marketActionLink = (market: any) => {
  const marketId = market.slug === 'chatuchak' ? 'chatuchak' : market.id
  return `/stalls?market=${marketId}`
}

const marketInspectStallsLink = (market: any) => {
  const marketId = market.slug === 'chatuchak' ? 'chatuchak' : market.id
  return `/stalls/manage?market=${marketId}&mode=inspect`
}

// Search & Filter State
const searchQuery = ref('')
const selectedCategory = ref<string>('all')
const appliedSearchQuery = ref('')
const appliedCategory = ref<string>('all')
const viewMode = ref<'grid' | 'map'>('grid')
const showGisMap = computed(() => viewMode.value === 'map')
const currentHeroSlide = ref(1)
const heroCarouselPaused = ref(false)
const heroSlides = [
  {
    image: '/images/market_promo_hero.jpg', alt: 'บรรยากาศตลาดนัดจตุจักร', tone: 'green',
    tag: 'โครงการสิทธิประโยชน์ผู้ค้า กทม. 2568 - 2569', subheading: 'หากคุณกำลังมองหา',
    title: 'ทำเลพื้นที่เช่า', highlight: '12 ตลาด กทม.',
    description: 'เปิดรับพ่อค้าแม่ค้า ประชาชนทั่วไป และผู้ประกอบการรายย่อย สำหรับประกอบการค้าปลีก สินค้าแฟชั่น อาหารสด อาหารสำเร็จรูป และสินค้าชุมชน ทำเลทองใจกลางเมือง',
    features: ['ค่าเช่าเริ่มต้น 500 บ./วัน', 'ติดรถไฟฟ้า & ขนส่งมวลชน', 'จองออนไลน์ รู้ผลอนุมัติทันที'],
    cta: 'สำรวจแผงเช่าว่าง 12 ตลาด', note: '* เงื่อนไขเป็นไปตามที่สำนักงานตลาดกรุงเทพมหานครกำหนด',
    ribbon: { tone: 'red', top: 'รับส่วนลดทันที', main: '50%', suffix: 'OFF' }
  },
  {
    image: '/images/market_hall_front.jpg', alt: 'บรรยากาศตลาดในอาคาร', tone: 'blue',
    tag: 'ตลาดแนะนำประจำเดือน', subheading: 'เปิดโอกาสให้ธุรกิจของคุณ',
    title: 'เริ่มต้นขายง่าย', highlight: 'ทำเลพร้อมเช่า',
    description: 'เลือกพื้นที่ค้าขายที่เหมาะกับธุรกิจ พร้อมข้อมูลแผง ราคา และสิ่งอำนวยความสะดวกครบในที่เดียว',
    features: ['ดูข้อมูลแผงแบบ 3 มิติ', 'เปรียบเทียบตลาดได้ทันที', 'จองวันนัดหมายออนไลน์'],
    cta: 'ดูตลาดแนะนำ', note: 'ข้อมูลพื้นที่อัปเดตตามรอบการจัดสรร', ribbon: { tone: 'blue', top: 'ตลาดแนะนำ', main: 'NEW', suffix: '' }
  },
  {
    image: '/images/chatuchak_entrance.jpg', alt: 'ทางเข้าตลาดนัดจตุจักร', tone: 'orange',
    tag: 'โอกาสพิเศษสำหรับผู้ค้าใหม่', subheading: 'ทำเลดี ลูกค้าเห็นง่าย',
    title: 'จองพื้นที่ขาย', highlight: 'ให้ธุรกิจเติบโต',
    description: 'พบพื้นที่ว่างในตลาดสังกัดกรุงเทพมหานคร พร้อมเงื่อนไขการเช่าที่ชัดเจนและขั้นตอนสมัครที่เข้าใจง่าย',
    features: ['มีพื้นที่หลายรูปแบบ', 'ค่าเช่าตามประเภทแผง', 'เจ้าหน้าที่ให้คำแนะนำ'],
    cta: 'ค้นหาพื้นที่ที่เหมาะกับคุณ', note: 'สิทธิ์และเงื่อนไขเป็นไปตามประกาศของตลาด', ribbon: null
  },
  {
    image: '/images/market_stall_sample.jpg', alt: 'ตัวอย่างแผงค้าในตลาด', tone: 'purple',
    tag: 'ประชาสัมพันธ์จากสำนักงานตลาด', subheading: 'พร้อมให้คุณเริ่มต้น',
    title: 'พื้นที่ค้าขาย', highlight: 'ที่ใช่สำหรับคุณ',
    description: 'ติดตามข่าวสาร กิจกรรม และพื้นที่เปิดใหม่จากตลาดในสังกัด กทม. ได้จากหน้าเดียว',
    features: ['ข่าวสารและกิจกรรมตลาด', 'ประกาศพื้นที่เปิดใหม่', 'ติดต่อเจ้าหน้าที่ได้สะดวก'],
    cta: 'ดูประกาศและกิจกรรม', note: 'ติดตามข้อมูลใหม่ได้ทุกสัปดาห์', ribbon: null
  }
]
const heroSlide = computed(() => heroSlides[currentHeroSlide.value - 1] || heroSlides[0])
let heroCarouselTimer: ReturnType<typeof setInterval> | null = null
const marketEvents = [
  { name: 'ช้อป ชิม ชิลล์ สัญจร 2026', date: '25–26 ต.ค. 2568', time: '10:00–22:00 น.', location: 'ตลาดนัดจตุจักร โครงการ 23', available: 0, marketId: 1, periods: ['today', 'week', 'month'], image: '/images/chatuchak_entrance.jpg' },
  { name: 'เทศกาลอาหารและดนตรี', date: '1–3 พ.ย. 2568', time: '11:00–21:00 น.', location: 'ตลาดธนบุรี (สนามหลวง 2)', available: 18, marketId: 11, periods: ['week', 'month'], image: '/images/market_hall_front.jpg' },
  { name: 'Neighborhood Market', date: '8–10 พ.ย. 2568', time: '10:00–22:00 น.', location: 'ตลาดเทวราช เขตดุสิต', available: 3, marketId: 5, periods: ['month'], image: '/images/market_promo_hero.jpg' }
]
const eventPeriod = ref<'today' | 'week' | 'month'>('week')
const eventPeriodTabs = computed(() => [
  { value: 'today', label: 'วันนี้', count: marketEvents.filter(event => event.periods.includes('today')).length },
  { value: 'week', label: 'สัปดาห์นี้', count: marketEvents.filter(event => event.periods.includes('week')).length },
  { value: 'month', label: 'เดือนนี้', count: marketEvents.filter(event => event.periods.includes('month')).length }
])
const filteredMarketEvents = computed(() => marketEvents.filter(event => event.periods.includes(eventPeriod.value)))

// Data State
const loadingMarkets = ref(true)
const marketsList = ref<any[]>([])
const filteredMarkets = ref<any[]>([])
let marketsRequestController: AbortController | null = null
let marketsLoadRun = 0
let marketsLoadedAt = 0
const showMarketSkeleton = computed(() => loadingMarkets.value)

const loadData = async () => {
  if (marketsRequestController || (marketsList.value.length > 0 && Date.now() - marketsLoadedAt < 30_000)) return
  const run = ++marketsLoadRun
  loadingMarkets.value = true
  try {
    let marketsRes: any = null
    let lastError: unknown = null

    // Retry once so a request opened while the backend is restarting cannot
    // leave the tenant home page on a permanent skeleton.
    for (let attempt = 0; attempt < 2; attempt += 1) {
      const controller = new AbortController()
      marketsRequestController = controller
      const timeoutId = window.setTimeout(() => controller.abort(), 8_000)
      try {
        marketsRes = await apiFetch('/markets', {
          method: 'GET',
          signal: controller.signal
        })
        lastError = null
        break
      } catch (err) {
        lastError = err
        if (run !== marketsLoadRun) return
      } finally {
        window.clearTimeout(timeoutId)
      }
    }

    if (lastError) throw lastError
    if (run !== marketsLoadRun) return
    const list = Array.isArray(marketsRes) ? marketsRes : (Array.isArray(marketsRes?.data) ? marketsRes.data : [])
    marketsList.value = list.map((market: any) => applyMarketOverride(market))
    filteredMarkets.value = marketsList.value
    marketsLoadedAt = Date.now()
  } catch (err) {
    if ((err as any)?.name === 'AbortError') return
    console.error('Failed to load markets:', err)
  } finally {
    if (run === marketsLoadRun) {
      loadingMarkets.value = false
      marketsRequestController = null
    }
  }
}

const filterMarkets = () => {
  let list = [...marketsList.value]

  // Filter by text search
  if (appliedSearchQuery.value.trim()) {
    const q = appliedSearchQuery.value.toLowerCase().trim()
    list = list.filter(m =>
      (m.name && m.name.toLowerCase().includes(q)) ||
      (m.name_en && m.name_en.toLowerCase().includes(q)) ||
      (m.district && m.district.toLowerCase().includes(q)) ||
      (m.address && m.address.toLowerCase().includes(q)) ||
      (m.tags && m.tags.some((t: string) => t.toLowerCase().includes(q)))
    )
  }

  // Filter by category icon
  if (appliedCategory.value === 'promo') {
    list = list.filter(m => m.starting_price_daily <= 120 || m.has_event)
  } else if (appliedCategory.value === 'new') {
    list = list.filter(m => m.available_stalls > 5)
  } else if (appliedCategory.value === 'metro') {
    list = list.filter(m => m.tags && m.tags.some((t: string) => t.includes('BTS') || t.includes('MRT') || t.includes('รถไฟฟ้า')))
  } else if (appliedCategory.value === 'parking') {
    list = list.filter(m => m.tags && m.tags.some((t: string) => t.includes('จอดรถ') || t.includes('รถเมล์')))
  }

  filteredMarkets.value = list
}

const toggleCategory = (cat: string) => {
  selectedCategory.value = cat
  appliedCategory.value = cat
  filterMarkets()
}

const applyMarketSearch = () => {
  appliedSearchQuery.value = searchQuery.value
  filterMarkets()
}

const clearSearch = () => {
  searchQuery.value = ''
  appliedSearchQuery.value = ''
  filterMarkets()
}

const resetFilters = () => {
  searchQuery.value = ''
  selectedCategory.value = 'all'
  appliedSearchQuery.value = ''
  appliedCategory.value = 'all'
  filterMarkets()
}

const scrollToMarkets = () => {
  const el = document.getElementById('markets-12-directory')
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' })
  }
}

const nextHeroSlide = () => {
  currentHeroSlide.value = currentHeroSlide.value === heroSlides.length
    ? 1
    : currentHeroSlide.value + 1
}

const pauseHeroCarousel = () => {
  heroCarouselPaused.value = true
}

const resumeHeroCarousel = () => {
  heroCarouselPaused.value = false
}

const startHeroCarousel = () => {
  if (heroCarouselTimer || heroSlides.length < 2) return
  heroCarouselTimer = setInterval(() => {
    if (!heroCarouselPaused.value && document.visibilityState === 'visible') {
      nextHeroSlide()
    }
  }, 5000)
}

const onSelectMarketOnMap = (market: any) => {
  // Option to focus market
}

const cleanTag = (tag?: string) => {
  if (!tag) return ''
  return tag.replace(/[\u{1F300}-\u{1F9FF}\u{2600}-\u{26FF}\u{2700}-\u{27BF}\u{1F1E0}-\u{1F1FF}]/gu, '').trim()
}

const getAvailableCount = (m: any) => {
  const total = m.total_stalls || 0
  return Math.max(0, Math.min(total, m.available_stalls ?? 0))
}

const getOccupiedCount = (m: any) => {
  const total = m.total_stalls || 0
  const explicit = Number(m.occupied_stalls ?? 0)
  if (explicit > 0) return Math.min(total, explicit)
  const reserved = Math.max(0, Number(m.reserved_stalls ?? 0))
  const avail = getAvailableCount(m)
  return Math.max(0, total - avail - reserved)
}

const getReservedCount = (m: any) => {
  const total = m.total_stalls || 0
  return Math.max(0, Math.min(total, Number(m.reserved_stalls ?? 0)))
}

const getOccupiedRate = (m: any) => {
  const total = m.total_stalls || 0
  if (!total) return 0
  const used = getOccupiedCount(m) + getReservedCount(m)
  return Math.min(100, Math.round((used / total) * 100))
}

const getSegmentPercent = (count: number, total?: number) => {
  const totalCount = Number(total || 0)
  if (!totalCount || count <= 0) return 0
  return Math.max(2, (count / totalCount) * 100)
}

const getAvailabilityRate = (m: any) => {
  const total = m.total_stalls || 0
  if (!total) return 0
  const rate = (getAvailableCount(m) / total) * 100
  // ถ้ายังมีแผงว่างแต่ต่ำกว่า 1% ให้แสดงอย่างน้อย 1% เพื่อไม่ให้ดูเหมือนเต็ม
  return getAvailableCount(m) > 0 ? Math.max(1, Math.round(rate)) : 0
}

// ความกว้างแถบ = สัดส่วนแผงว่างจริง (เช่น 8/180 ≈ 4.4%)
// กรณีเต็ม ให้แถบเต็มสีแดงเพื่อสื่อสถานะ "เต็ม"
const getAvailabilityBarWidth = (m: any) => {
  const total = m.total_stalls || 0
  const avail = getAvailableCount(m)
  if (!total || avail === 0) return 100
  return Math.max(2, (avail / total) * 100)
}

// ระดับความใกล้เต็ม: ว่าง >= 50% เขียว | 20-49% เหลือง | 1-19% ส้ม | 0 แดง
const getAvailabilityLevel = (m: any) => {
  const total = m.total_stalls || 0
  const avail = getAvailableCount(m)
  if (!total || avail === 0) return 'full'
  const pct = (avail / total) * 100
  if (pct >= 50) return 'ok'
  if (pct >= 20) return 'mid'
  return 'low'
}

const getAvailabilityBarColor = (m: any) => `bar-${getAvailabilityLevel(m)}`

const getAvailabilityTextClass = (m: any) => `avail-text-${getAvailabilityLevel(m)}`

const getMarketViews = (marketId: number | string) => {
  const seed = typeof marketId === 'number' ? marketId : parseInt(String(marketId)) || 1
  const views = [342, 480, 215, 620, 390, 540, 290, 430, 310, 670, 450, 380]
  return views[(seed - 1) % views.length] || 342
}

const getStartingPrice = (m: any) => {
  if (m.starting_price_monthly) {
    return Number(m.starting_price_monthly).toLocaleString()
  }
  if (m.starting_price_daily) {
    return Number(m.starting_price_daily).toLocaleString()
  }
  return '500'
}

const getPriceUnit = (m: any) => (m.starting_price_monthly || !m.starting_price_daily ? 'เดือน' : 'วัน')

const MARKET_RETURN_STATE_KEY = 'e-market:market-return-state'

onBeforeRouteLeave((to) => {
  if (typeof window === 'undefined') return
  if (to.path.startsWith('/stalls') || to.path.startsWith('/markets/')) {
    try {
      sessionStorage.setItem(MARKET_RETURN_STATE_KEY, JSON.stringify({
        scrollY: window.scrollY,
        viewMode: viewMode.value,
        searchQuery: searchQuery.value,
        selectedCategory: selectedCategory.value
      }))
    } catch {}
  } else {
    sessionStorage.removeItem(MARKET_RETURN_STATE_KEY)
    sessionStorage.removeItem('e-market:bma-map-state')
  }
})

onMounted(async () => {
  startHeroCarousel()

  let savedState: {
    scrollY?: number
    viewMode?: 'grid' | 'map'
    searchQuery?: string
    selectedCategory?: string
  } | null = null

  if (typeof window !== 'undefined') {
    sessionStorage.removeItem('e-market:market-list-scroll-y')
    try {
      const raw = sessionStorage.getItem(MARKET_RETURN_STATE_KEY)
      if (raw) {
        savedState = JSON.parse(raw)
        sessionStorage.removeItem(MARKET_RETURN_STATE_KEY)
      }
    } catch {}

    if (savedState) {
      if (savedState.viewMode === 'map' || savedState.viewMode === 'grid') {
        viewMode.value = savedState.viewMode
      }
      if (typeof savedState.searchQuery === 'string') {
        searchQuery.value = savedState.searchQuery
        appliedSearchQuery.value = savedState.searchQuery
      }
      if (typeof savedState.selectedCategory === 'string') {
        selectedCategory.value = savedState.selectedCategory
        appliedCategory.value = savedState.selectedCategory
      }
    } else {
      sessionStorage.removeItem('e-market:bma-map-state')
      window.scrollTo({ top: 0, left: 0, behavior: 'auto' })
    }
  }

  await loadData()
  if (savedState) {
    filterMarkets()
    if (typeof window !== 'undefined' && typeof savedState.scrollY === 'number') {
      const targetY = savedState.scrollY
      await nextTick()
      window.scrollTo({ top: targetY, left: 0, behavior: 'auto' })
      requestAnimationFrame(() => {
        window.scrollTo({ top: targetY, left: 0, behavior: 'auto' })
      })
      setTimeout(() => {
        window.scrollTo({ top: targetY, left: 0, behavior: 'auto' })
      }, 60)
    }
  }

  if (typeof window !== 'undefined') {
    window.addEventListener('focus', loadData)
  }
})

onBeforeUnmount(() => {
  marketsLoadRun += 1
  marketsRequestController?.abort()
  marketsRequestController = null
  if (heroCarouselTimer) {
    clearInterval(heroCarouselTimer)
    heroCarouselTimer = null
  }
  if (typeof window !== 'undefined') {
    window.removeEventListener('focus', loadData)
  }
})
</script>

<style scoped>
/* Main Container */
.homepage-container {
  background-color: var(--bg-main);
  min-height: 100vh;
  --hero-to-events-gap: 80px;
  --events-to-markets-gap: 128px;
}

/* ============================================================== */
/* 1. HERO PROMOTIONAL BANNER                                     */
/* ============================================================== */
.hero-banner-wrapper {
  margin-bottom: var(--hero-to-events-gap);
}

.hero-banner-card {
  display: flex;
  padding: 0;
  overflow: hidden;
  border-radius: var(--radius-2xl);
  height: clamp(360px, 28vw, 460px);
  background: var(--bg-ffffff);
  border: 1px solid var(--border-card);
  box-shadow: 0 10px 30px -5px rgba(37, 99, 235, 0.08), 0 4px 12px rgba(15, 23, 42, 0.04);
  position: relative;
}

.hero-left-media {
  position: relative;
  flex: 1.08 1 0;
  min-width: 0;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.hero-media-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: opacity 0.35s ease, transform 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

.hero-banner-card:hover .hero-media-img {
  transform: scale(1.03);
}

.hero-media-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to right, var(--bg-rgba-0-0-0-0p1) 0%, var(--bg-rgba-0-0-0-0) 100%);
}

/* Diagonal Cut Red 50% OFF Badge */
.hero-discount-ribbon {
  position: absolute;
  top: 50%;
  left: 0;
  transform: translateY(-50%);
  background: linear-gradient(135deg, #e11d48 0%, #be123c 100%);
  color: #ffffff;
  padding: 12px 30px 12px 22px;
  clip-path: polygon(0 0, 100% 0, 85% 100%, 0% 100%);
  box-shadow: 0 8px 24px rgba(225, 29, 72, 0.45);
  z-index: 5;
}

.hero-ribbon-blue {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
}

.ribbon-text-top {
  font-size: var(--font-size-xs);
  font-weight: 700;
  letter-spacing: 0.04em;
  opacity: 0.95;
}

.ribbon-text-main {
  font-size: 1.85rem;
  font-weight: 900;
  line-height: 1.05;
  font-family: var(--font-display);
}

.ribbon-off {
  font-size: var(--font-size-base);
  font-weight: 800;
}

.hero-right-content {
  flex: 0.92 1 0;
  min-width: 0;
  padding: 28px 36px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
  background: var(--bg-ffffff);
}

.campaign-tag-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-ecfdf5);
  color: var(--tx-047857);
  border: 1px solid var(--bd-a7f3d0);
  padding: 3px 12px;
  border-radius: 20px;
  font-size: var(--font-size-xs);
  font-weight: 700;
  width: fit-content;
  margin-bottom: 8px;
}

.campaign-tag-blue {
  color: var(--tx-1d4ed8);
  background: var(--bg-eff6ff);
  border-color: var(--bd-bfdbfe);
}

.campaign-tag-orange {
  color: var(--tx-c2410c);
  background: var(--bg-fff7ed);
  border-color: var(--bd-fed7aa);
}

.campaign-tag-purple {
  color: var(--tx-7e22ce);
  background: var(--bg-faf5ff);
  border-color: var(--bd-e9d5ff);
}

.campaign-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #10b981;
  animation: pulse-dot 2s infinite;
}

.hero-subheading {
  font-size: var(--font-size-md);
  font-weight: 700;
  color: var(--tx-059669);
  margin-bottom: 2px;
}

.hero-main-heading {
  font-size: 2.25rem;
  font-weight: 900;
  color: var(--tx-0f172a);
  line-height: 1.22;
  margin-bottom: 12px;
  letter-spacing: -0.02em;
}

.heading-gradient {
  background: linear-gradient(135deg, #1e40af 0%, #2563eb 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-description {
  font-size: var(--font-size-sm);
  color: var(--tx-475569);
  line-height: 1.6;
  max-width: 480px;
  margin-bottom: 16px;
}

.hero-features-list {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.hero-feature-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--tx-334155);
  background: var(--bg-f8fafc);
  padding: 4px 10px;
  border-radius: 8px;
  border: 1px solid var(--bd-e2e8f0);
}

.hero-cta-row {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.btn-hero-cta {
  padding: 11px 22px;
  border-radius: 12px;
  font-weight: 700;
  font-size: var(--font-size-base);
  box-shadow: 0 4px 16px rgba(37, 99, 235, 0.3);
}

.hero-note {
  font-size: var(--font-size-xs);
  color: var(--tx-94a3b8);
}

.hero-accent-shape {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 140px;
  height: 8px;
  background: linear-gradient(90deg, transparent, #1d4ed8);
  border-top-left-radius: 8px;
}

.hero-carousel-dots {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
}

.carousel-dot {
  appearance: none;
  border: 0;
  padding: 0;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--bg-cbd5e1);
  transition: all 0.25s ease;
  cursor: pointer;
}

.carousel-dot.active {
  width: 26px;
  border-radius: 8px;
  background: var(--primary);
}

/* ============================================================== */
/* 2. SECTION: อีเวนต์ตลาด                                        */
/* ============================================================== */
.events-card-wrapper {
  margin-bottom: var(--events-to-markets-gap);
  padding: 0;
  background: transparent;
}

.events-promo-header {
  width: min(100%, 860px);
  min-height: 240px;
  display: grid;
  grid-template-columns: minmax(0, 470px) 320px;
  align-items: center;
  justify-content: center;
  gap: 36px;
  margin: 0 auto;
  padding: 12px 16px 8px;
  box-sizing: border-box;
  background: transparent;
}

.events-promo-copy {
  max-width: 680px;
}

.events-promo-copy .events-badge-pill {
  margin-bottom: 12px;
}

.events-promo-copy .events-main-title {
  margin: 0;
  font-size: clamp(1.75rem, 2.2vw, 2.35rem);
}

.events-promo-copy .events-main-title::after {
  content: '';
  display: block;
  width: 54px;
  height: 4px;
  margin-top: 9px;
  border-radius: 999px;
  background: #2563eb;
}

.events-promo-copy .events-desc-text {
  max-width: 610px;
  margin: 18px 0 20px;
  color: var(--tx-52637a);
  font-size: var(--font-size-base);
  line-height: 1.75;
}

.events-promo-image {
  width: min(100%, 380px);
  max-height: 310px;
  justify-self: center;
  object-fit: contain;
  filter: drop-shadow(0 14px 28px rgba(37, 99, 235, 0.14));
}

.events-section-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 22px;
}

.events-section-header .events-desc-text {
  margin: 6px 0 0;
}

.events-view-all {
  min-height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 16px;
  border: 1px solid #2563eb;
  border-radius: 9px;
  color: var(--tx-2563eb);
  background: transparent;
  font: inherit;
  font-size: var(--font-size-sm);
  font-weight: 700;
  cursor: pointer;
  transition: border-color 0.2s ease, color 0.2s ease, background 0.2s ease;
}

.events-view-all:hover {
  color: var(--tx-1d4ed8);
  border-color: #1d4ed8;
  background: var(--bg-eff6ff);
}

.events-period-tabs {
  width: fit-content;
  display: grid;
  grid-template-columns: repeat(3, minmax(92px, 1fr));
  margin: 16px auto 0;
  padding: 4px;
  border: 1px solid var(--bd-dbe4f0);
  border-radius: 12px;
  background: var(--bg-f8fafc);
  position: relative;
  z-index: 1;
}

.events-list-panel {
  max-width: 880px;
  width: calc(100% - 48px);
  margin: 48px auto 0;
  padding-top: 0;
  overflow: hidden;
  border: 1px solid var(--bd-dfe7f1);
  border-radius: 18px;
  background: var(--bg-ffffff);
  box-shadow: 0 10px 28px rgba(15, 23, 42, 0.045);
}

.events-period-tabs button {
  min-height: 38px;
  padding: 7px 16px;
  border: 0;
  border-radius: 8px;
  color: var(--tx-64748b);
  background: transparent;
  font: inherit;
  font-size: var(--font-size-sm);
  font-weight: 600;
  cursor: pointer;
  transition: color 0.18s ease, background 0.18s ease, box-shadow 0.18s ease;
}

.events-period-tabs button:hover:not(.active) {
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
}

.events-period-tabs button.active {
  color: #ffffff;
  background: #2563eb;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.2);
}

.events-list {
  margin-top: 22px;
  padding: 0 20px 18px;
}

.event-list-row {
  min-height: 118px;
  display: grid;
  grid-template-columns: 108px minmax(0, 1fr) 220px;
  align-items: center;
  gap: 22px;
  padding: 18px 8px;
  border-bottom: 1px solid var(--bd-e8edf4);
  transition: background 0.18s ease;
}

.event-list-row:first-child {
  border-top: 1px solid var(--bd-e8edf4);
}

.event-list-row:hover {
  background: var(--bg-fbfdff);
}

.event-list-image {
  width: 108px;
  height: 82px;
  display: block;
  object-fit: cover;
  border-radius: 12px;
}

.event-list-content {
  min-width: 0;
}

.event-list-content h3 {
  margin: 0 0 9px;
  color: var(--tx-172033);
  font-size: var(--font-size-lg);
  font-weight: 700;
  line-height: 1.4;
}

.event-list-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px 16px;
  margin-bottom: 7px;
  color: var(--tx-64748b);
  font-size: var(--font-size-sm);
}

.event-list-meta > span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.event-list-meta svg {
  color: var(--tx-2563eb);
}

.event-list-meta .event-time-pill {
  padding: 4px 10px;
  border-radius: 999px;
  color: var(--tx-0284c7);
  background: var(--bg-f0f9ff);
}

.event-list-row .event-location-line {
  min-height: 0;
  font-size: var(--font-size-sm);
}

.event-list-action {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 16px;
}

.event-capacity {
  min-width: 68px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1px;
  color: var(--tx-059669);
}

.event-capacity strong {
  font-size: 1.75rem;
  font-weight: 800;
  line-height: 1;
}

.event-capacity span {
  padding: 2px 8px;
  border-radius: 999px;
  color: var(--tx-047857);
  background: var(--bg-ecfdf5);
  font-size: var(--font-size-xs);
  font-weight: 700;
}

.event-capacity.is-low {
  color: var(--tx-d97706);
}

.event-capacity.is-low span {
  color: var(--tx-b45309);
  background: var(--bg-fffbeb);
}

.event-capacity.is-full {
  color: var(--tx-ef4444);
}

.event-capacity.is-full span {
  color: var(--tx-dc2626);
  background: var(--bg-fef2f2);
}

.events-cards-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.event-showcase-card {
  min-width: 0;
  overflow: hidden;
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 16px;
  background: var(--bg-ffffff);
  transition: transform 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

.event-showcase-card:hover {
  transform: translateY(-2px);
  border-color: var(--bd-bfdbfe);
  box-shadow: 0 10px 24px rgba(37, 99, 235, 0.09);
}

.event-showcase-media {
  position: relative;
  height: 170px;
  overflow: hidden;
}

.event-showcase-media img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  transition: transform 0.35s ease;
}

.event-showcase-card:hover .event-showcase-media img {
  transform: scale(1.035);
}

.event-date-pill {
  position: absolute;
  left: 14px;
  bottom: 12px;
  padding: 5px 10px;
  border-radius: 999px;
  color: var(--tx-1e3a8a);
  background: var(--bg-rgba-255-255-255-0p94);
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.12);
  font-size: var(--font-size-xs);
  font-weight: 700;
}

.event-showcase-body {
  padding: 16px;
}

.event-showcase-body h3 {
  margin: 0 0 8px;
  color: var(--tx-0f172a);
  font-size: var(--font-size-lg);
  font-weight: 700;
  line-height: 1.35;
}

.event-location-line {
  min-height: 40px;
  display: flex;
  align-items: flex-start;
  gap: 7px;
  color: var(--tx-64748b);
  font-size: var(--font-size-sm);
  line-height: 1.45;
}

.event-location-line svg {
  flex-shrink: 0;
  margin-top: 2px;
  color: var(--tx-ef4444);
}

.event-showcase-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid var(--bd-f1f5f9);
}

.event-availability {
  color: var(--tx-15803d);
  font-size: var(--font-size-sm);
  font-weight: 700;
}

.event-availability.is-full {
  color: var(--tx-64748b);
}

.event-primary-link,
.event-disabled-link {
  min-height: 36px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 7px 12px;
  border-radius: 9px;
  font-size: var(--font-size-sm);
  font-weight: 700;
  text-decoration: none;
}

.event-primary-link {
  color: #ffffff;
  background: #2563eb;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.2);
}

.event-primary-link:hover {
  background: #1d4ed8;
}

.event-disabled-link {
  color: var(--tx-94a3b8);
  background: var(--bg-f1f5f9);
}

.events-grid-layout {
  display: grid;
  grid-template-columns: 280px 200px 1fr;
  gap: 28px;
  align-items: center;
}

.events-left-col {
  padding: 12px 16px 12px 32px;
}

.events-badge-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: var(--font-size-xs);
  font-weight: 700;
  color: var(--tx-2563eb);
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-dbeafe);
  padding: 3px 10px;
  border-radius: 20px;
  margin-bottom: 8px;
}

.events-main-title {
  font-size: 1.85rem;
  font-weight: 800;
  color: var(--tx-0f172a);
  letter-spacing: -0.02em;
}

.title-blue-bar {
  width: 44px;
  height: 4px;
  background: var(--primary);
  border-radius: var(--radius-full);
  margin: 6px 0 14px;
}

.events-desc-text {
  font-size: var(--font-size-sm);
  color: var(--tx-64748b);
  line-height: 1.6;
  margin-bottom: 18px;
}

.btn-events-all {
  border-radius: 10px;
  font-weight: 700;
}

.events-center-col {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 220px;
}

.megaphone-img {
  width: 150px;
  height: auto;
  object-fit: contain;
  filter: drop-shadow(0 14px 20px rgba(0, 0, 0, 0.12));
}

.floating-animation {
  animation: float 4s ease-in-out infinite;
}

.events-right-col {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.event-row-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 14px;
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  transition: all 0.22s cubic-bezier(0.4, 0, 0.2, 1);
}

.event-row-card:hover {
  background: var(--bg-ffffff);
  border-color: var(--bd-bfdbfe);
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.08);
  transform: translateY(-2px);
}

.event-thumb-box {
  width: 76px;
  height: 64px;
  border-radius: 10px;
  overflow: hidden;
  flex-shrink: 0;
}

.event-thumb-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.event-row-card:hover .event-thumb-img {
  transform: scale(1.08);
}

.event-details-box {
  flex: 1;
  min-width: 0;
}

.event-name {
  font-size: var(--font-size-base);
  font-weight: 700;
  color: var(--tx-0f172a);
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.event-meta-line {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 2px;
}

.event-action-box {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.event-stall-stat {
  padding: 4px 10px;
  border-radius: 8px;
  display: flex;
  align-items: baseline;
  gap: 3px;
  font-weight: 700;
}

.stat-zero {
  background: var(--bg-f1f5f9);
  color: var(--tx-94a3b8);
  border: 1px solid var(--bd-e2e8f0);
}

.stat-green {
  background: var(--bg-ecfdf5);
  color: var(--tx-047857);
  border: 1px solid var(--bd-a7f3d0);
}

.stat-orange {
  background: var(--bg-fffbeb);
  color: var(--tx-b45309);
  border: 1px solid var(--bd-fde68a);
}

.stat-number {
  font-size: var(--font-size-lg);
  font-family: var(--font-display);
}

.stat-unit {
  font-size: var(--font-size-xs);
}

.btn-event-act {
  border-radius: 9px;
  padding: 6px 14px;
  font-weight: 700;
}

.events-page-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 10px;
}

.page-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: 1px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  color: var(--tx-475569);
  font-size: var(--font-size-xs);
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
}

.page-btn:hover:not(:disabled) {
  background: var(--bg-f1f5f9);
  color: var(--primary);
  border-color: var(--bd-cbd5e1);
}

.page-btn.active {
  background: var(--primary);
  color: #ffffff;
  border-color: var(--primary);
  font-weight: 700;
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* ============================================================== */
/* 3. SECTION: 12 ตลาด DIRECTORY                                   */
/* ============================================================== */
.directory-title-header {
  text-align: center;
  margin-bottom: 24px;
}

/* ระยะห่างระหว่าง 3 ส่วนหลัก (Banner / อีเวนต์ / ตลาด) บนจอเล็ก */
@media (max-width: 768px) {
  .homepage-container {
    --hero-to-events-gap: 48px;
    --events-to-markets-gap: 72px;
  }

  .hero-banner-wrapper {
    margin-bottom: var(--hero-to-events-gap);
  }

  .events-card-wrapper {
    margin-bottom: var(--events-to-markets-gap);
  }
}

.directory-header-pill {
  display: inline-block;
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-dbeafe);
  color: var(--tx-1d4ed8);
  padding: 4px 14px;
  border-radius: 20px;
  font-size: var(--font-size-xs);
  font-weight: 700;
  margin-bottom: 8px;
}

.directory-big-heading {
  font-size: 2.25rem;
  font-weight: 900;
  color: var(--tx-0f172a);
  letter-spacing: -0.02em;
}

.heading-underline {
  position: relative;
  display: inline-block;
}

/* ให้เส้นใต้เหมือน .title-blue-bar ของหัวข้ออีเวนต์ (44x4px, ห่าง 6px) แต่จัดกึ่งกลาง */
.heading-underline::after {
  content: '';
  position: absolute;
  bottom: -10px;
  left: 50%;
  transform: translateX(-50%);
  width: 44px;
  height: 4px;
  background: var(--primary);
  border-radius: var(--radius-full);
}

.directory-sub-caption {
  font-size: var(--font-size-base);
  color: var(--tx-64748b);
  margin-top: 24px;
  max-width: 620px;
  margin-left: auto;
  margin-right: auto;
}

.search-bar-centered {
  max-width: 580px;
  margin: 0 auto 28px;
}

.search-input-pill-wrap {
  position: relative;
  display: flex;
  align-items: center;
  background: var(--bg-ffffff);
  border: 1.5px solid var(--bd-cbd5e1);
  border-radius: 30px;
  padding: 4px 16px;
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.05);
  transition: all 0.22s ease;
}

.search-input-pill-wrap:focus-within {
  border-color: var(--primary);
  box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.15), 0 4px 16px rgba(15, 23, 42, 0.08);
}

.search-icon-left {
  display: flex;
  align-items: center;
  color: var(--tx-64748b);
  margin-right: 10px;
}

.search-input-field {
  flex: 1;
  border: none;
  background: transparent;
  padding: 10px 0;
  font-size: var(--font-size-base);
  color: var(--tx-0f172a);
  outline: none;
}

.search-input-field::placeholder {
  color: var(--tx-94a3b8);
}

.search-clear-btn {
  border: none;
  background: var(--bg-e2e8f0);
  color: var(--tx-475569);
  width: 22px;
  height: 22px;
  border-radius: 50%;
  font-size: var(--font-size-xs);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
}

.search-clear-btn:hover {
  background: var(--bg-cbd5e1);
}


/* Category Filter Buttons Row */
.category-icon-buttons-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
  margin-bottom: 28px;
  flex-wrap: wrap;
}

.cat-icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.cat-icon-item:hover {
  transform: translateY(-3px);
}

.cat-icon-circle {
  width: 58px;
  height: 58px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.375rem;
  border: 2px solid var(--bd-e2e8f0);
  background: var(--bg-ffffff);
  color: var(--tx-64748b);
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.05);
  transition: all 0.22s ease;
}

.cat-icon-circle.cat-all {
  color: var(--tx-475569);
}

.cat-icon-circle.cat-promo {
  color: var(--tx-e11d48);
}

.cat-icon-circle.cat-new {
  color: var(--tx-059669);
}

.cat-icon-circle.cat-metro {
  color: var(--tx-2563eb);
}

.cat-icon-circle.cat-transit {
  color: var(--tx-d97706);
}

.cat-icon-item:hover .cat-icon-circle {
  border-color: var(--bd-cbd5e1);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.1);
}

.cat-icon-item.active .cat-icon-circle {
  background: var(--primary-gradient);
  border-color: var(--primary);
  color: #ffffff;
  box-shadow: 0 8px 20px rgba(37, 99, 235, 0.35);
  transform: scale(1.05);
}

.cat-symbol {
  font-weight: 900;
  font-family: var(--font-display);
}

.cat-label {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--tx-475569);
  transition: color 0.15s ease;
}

.cat-icon-item.active .cat-label {
  color: var(--primary);
  font-weight: 800;
}

.market-view-switcher-bar,
.map-view-toggle-bar {
  display: flex;
  justify-content: center;
  margin-bottom: 32px;
}

.segmented-view-control {
  display: inline-flex;
  align-items: center;
  padding: 4px;
  background: var(--bg-f1f5f9);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 9999px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  gap: 4px;
}

.svc-tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 9px 22px;
  border-radius: 9999px;
  border: none;
  background: transparent;
  color: var(--tx-64748b);
  font-size: var(--font-size-sm);
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
}

.svc-tab-btn:hover:not(.active) {
  color: var(--tx-1e293b);
  background: var(--bg-rgba-255-255-255-0p6);
}

.svc-tab-btn.active {
  background: var(--bg-ffffff);
  color: var(--tx-2563eb);
  font-weight: 700;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.svc-counter {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-xs);
  font-weight: 700;
  padding: 1px 7px;
  border-radius: 12px;
  background: var(--bg-e2e8f0);
  color: var(--tx-475569);
  transition: all 0.2s ease;
}

.svc-counter.is-loading {
  width: 30px;
  height: 22px;
  padding: 0;
  overflow: hidden;
  position: relative;
  background: var(--bg-e2e8f0);
}

.svc-counter.is-loading::after {
  content: '';
  position: absolute;
  inset: 0;
  transform: translateX(-100%);
  background: linear-gradient(90deg, transparent, var(--bg-rgba-255-255-255-0p8), transparent);
  animation: marketSkeletonShimmer 1.35s ease-in-out infinite;
}

.svc-tab-btn.active .svc-counter {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
}

.market-gis-view-container {
  border-radius: 18px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
}

/* ============================================================== */
/* 12 MARKET CARDS GRID                                           */
/* ============================================================== */
.market-cards-3col-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.market-box-card {
  padding: 16px;
  border-radius: 20px;
  border: 1px solid var(--bd-eef2f6);
  background: var(--bg-ffffff);
  box-shadow: 0 4px 20px rgba(15, 23, 42, 0.04);
  transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
}

.market-box-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 16px 32px -4px rgba(15, 23, 42, 0.08), 0 4px 12px rgba(15, 23, 42, 0.04);
  border-color: var(--bd-cbd5e1);
}

.market-skeleton-card {
  pointer-events: none;
}

.market-skeleton-shape {
  display: block;
  position: relative;
  overflow: hidden;
  background: var(--bg-eef2f7);
  border-radius: 8px;
}

.market-skeleton-shape::after {
  content: '';
  position: absolute;
  inset: 0;
  transform: translateX(-100%);
  background: linear-gradient(90deg, transparent, var(--bg-rgba-255-255-255-0p8), transparent);
  animation: marketSkeletonShimmer 1.35s ease-in-out infinite;
}

.market-skeleton-image {
  width: 100%;
  aspect-ratio: 16 / 10;
  border-radius: 14px;
}

.market-skeleton-body {
  gap: 10px;
}

.market-skeleton-title { width: 48%; height: 25px; }
.market-skeleton-address { width: 92%; height: 14px; }
.market-skeleton-address.is-short { width: 62%; }

.market-skeleton-availability {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 10px;
  padding: 18px 20px;
  margin-top: 10px;
  border: 1px solid var(--bd-eef2f6);
  border-radius: 14px;
  background: var(--bg-f8fafc);
}

.market-skeleton-label { width: 136px; height: 15px; }
.market-skeleton-count { width: 84px; height: 20px; justify-self: end; }
.market-skeleton-progress { grid-column: 1 / -1; width: 100%; height: 10px; border-radius: 999px; }

.market-skeleton-price-row {
  display: flex;
  align-items: baseline;
  justify-content: flex-end;
  gap: 8px;
  margin: 16px 0 8px;
}

.market-skeleton-price-label { width: 52px; height: 16px; }
.market-skeleton-price { width: 126px; height: 38px; }
.market-skeleton-button { width: 100%; height: 48px; border-radius: 12px; }

@keyframes marketSkeletonShimmer {
  to { transform: translateX(100%); }
}

@media (prefers-reduced-motion: reduce) {
  .market-skeleton-shape::after { animation: none; }
}

.market-img-container {
  position: relative;
  width: 100%;
  aspect-ratio: 16 / 10;
  overflow: hidden;
  background: var(--bg-f1f5f9);
  border-radius: 14px;
  display: block;
}

.market-box-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}

.market-box-card:hover .market-box-img {
  transform: scale(1.05);
}

.market-img-gradient {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, var(--bg-rgba-0-0-0-0p12) 0%, transparent 35%, transparent 70%, var(--bg-rgba-0-0-0-0p18) 100%);
}

.market-avail-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 5px 14px;
  border-radius: 20px;
  font-size: var(--font-size-sm);
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.badge-green {
  background: #10b981;
  color: #ffffff;
}

.badge-red {
  background: #ef4444;
  color: #ffffff;
}

.badge-dot-live {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--bg-ffffff);
}

.market-box-body {
  padding: 16px 4px 4px 4px;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.market-title-views-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.market-title-col {
  display: flex;
  flex-direction: column;
}

.market-title-link {
  text-decoration: none;
}

.market-card-name {
  font-size: var(--font-size-xl);
  font-weight: 700;
  color: var(--tx-1e293b);
  line-height: 1.3;
  transition: color 0.15s ease;
}

.market-card-name:hover {
  color: var(--tx-2563eb);
}

.market-views-pill {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  padding: 5px 12px;
  border-radius: 10px;
  font-size: var(--font-size-sm);
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  flex-shrink: 0;
}

.market-card-address {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  font-size: var(--font-size-sm);
  color: var(--tx-64748b);
  line-height: 1.5;
  margin-bottom: 14px;
  min-height: calc(0.84rem * 1.5 * 2);
}

.address-text {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.pin-icon {
  display: inline-flex;
  align-items: center;
  color: var(--tx-94a3b8);
  flex-shrink: 0;
  margin-top: 2px;
}

.market-occupancy-row {
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-eef2f6);
  border-radius: 12px;
  padding: 12px 14px;
  margin-bottom: 16px;
}

.occupancy-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.occupancy-title-col {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  font-weight: 600;
}

.occupancy-count {
  display: inline-flex;
  align-items: baseline;
  gap: 4px;
  font-size: var(--font-size-xs);
}

.occupancy-count strong {
  font-size: var(--font-size-base);
  font-weight: 700;
}

.occupancy-total {
  color: var(--tx-94a3b8);
  font-weight: 500;
}

.avail-text-ok {
  color: var(--tx-15803d);
}

.avail-text-mid {
  color: var(--tx-ca8a04);
}

.avail-text-low {
  color: var(--tx-ea580c);
}

.avail-text-full {
  color: var(--tx-dc2626);
}

.occupancy-progress-bar {
  width: 100%;
  height: 6px;
  background: var(--bg-e2e8f0);
  border-radius: 999px;
  overflow: hidden;
}

.occupancy-progress-fill {
  height: 100%;
  border-radius: 999px;
  transition: width 0.4s ease;
}

.occupancy-progress-fill.bar-ok {
  background: linear-gradient(90deg, #22c55e 0%, #16a34a 100%);
}

.occupancy-progress-fill.bar-mid {
  background: linear-gradient(90deg, #facc15 0%, #eab308 100%);
}

.occupancy-progress-fill.bar-low {
  background: linear-gradient(90deg, #fb923c 0%, #ea580c 100%);
}

.occupancy-progress-fill.bar-full {
  background: linear-gradient(90deg, #ef4444 0%, #b91c1c 100%);
}

.occupancy-meta-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: var(--font-size-xs);
  color: var(--tx-94a3b8);
  font-weight: 500;
}

.meta-occupied-text {
  color: var(--tx-64748b);
  font-weight: 600;
}

.meta-total-text {
  color: var(--tx-94a3b8);
}

.market-price-row {
  display: flex;
  justify-content: flex-end;
  align-items: baseline;
  gap: 2px;
  margin-bottom: 14px;
  margin-top: auto;
}

.price-prefix-text {
  font-size: var(--font-size-xs);
  color: var(--tx-94a3b8);
  font-weight: 500;
  margin-right: 6px;
}

.price-currency {
  font-size: var(--font-size-lg);
  font-weight: 700;
  color: var(--tx-1e293b);
  margin-right: 1px;
}

.price-big-number {
  font-size: 1.9rem;
  font-weight: 800;
  color: var(--tx-1e293b);
  line-height: 1;
  letter-spacing: -0.02em;
}

.price-suffix-text {
  font-size: var(--font-size-sm);
  color: var(--tx-64748b);
  font-weight: 500;
  margin-left: 4px;
}

.market-card-actions {
  margin-top: auto;
}

.btn-book-action-bar {
  width: 100%;
  background: #2563eb;
  color: #ffffff;
  border-radius: 12px;
  padding: 12px 18px;
  font-weight: 600;
  font-size: var(--font-size-base);
  display: flex;
  justify-content: space-between;
  align-items: center;
  text-decoration: none;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.22);
  transition: all 0.2s ease;
}

.btn-book-action-bar:hover {
  background: #1d4ed8;
  box-shadow: 0 6px 16px rgba(29, 78, 216, 0.3);
  transform: translateY(-1px);
}

.btn-book-action-bar .btn-arrow-icon {
  transition: transform 0.2s ease;
}

.btn-book-action-bar:hover .btn-arrow-icon {
  transform: translateX(4px);
}

.btn-book-stall:hover .btn-arrow-icon {
  transform: translateX(4px);
}

/* ============================================================== */
/* OFFICER MARKET CARD VIEW                                       */
/* ============================================================== */
.officer-district-pill {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 999px;
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-dbeafe);
  color: var(--tx-2563eb);
  font-size: var(--font-size-xs);
  font-weight: 600;
  white-space: nowrap;
  flex-shrink: 0;
}

.officer-occupancy-box {
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 14px;
  padding: 12px 14px;
  margin-bottom: 22px;
}

.officer-occupancy-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 10px;
  margin-bottom: 10px;
  border-bottom: 1px solid var(--bd-e2e8f0);
}

.officer-occupancy-title {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--tx-475569);
}

.officer-total-stalls {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  font-weight: 500;
}

.officer-total-stalls strong {
  color: var(--tx-0f172a);
  font-weight: 700;
}

.officer-stall-stats-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  text-align: center;
  padding: 2px 0;
}

.officer-stat-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px 8px;
}

.officer-stat-cell + .officer-stat-cell {
  border-left: 1px solid var(--bd-e2e8f0);
}

.officer-stat-cell .stat-copy {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  min-width: 0;
}

.officer-stat-cell .stat-label {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  font-weight: 600;
  line-height: 1.4;
}

.officer-stat-cell .stat-val {
  font-size: var(--font-size-md);
  font-weight: 800;
  color: var(--tx-0f172a);
  line-height: 1.3;
}

.officer-stat-cell.stat-occupied .stat-val {
  color: var(--tx-0f172a);
}

.officer-stat-cell.stat-reserved .stat-val {
  color: var(--tx-d97706);
}

.officer-stat-cell.stat-available .stat-val {
  color: var(--tx-10b981);
}

.empty-search-card {
  text-align: center;
  padding: 48px 24px;
  max-width: 500px;
  margin: 32px auto;
}

/* Responsive adjustments */
@media (max-width: 1100px) {
  .events-promo-header {
    grid-template-columns: minmax(0, 1fr) 310px;
    gap: 24px;
    padding-inline: 44px;
  }
  .events-list-panel {
    max-width: 840px;
    width: calc(100% - 48px);
    margin: 40px auto 0;
  }
  .event-list-row {
    grid-template-columns: 96px minmax(0, 1fr) 200px;
    gap: 16px;
  }
  .event-list-image {
    width: 96px;
    height: 76px;
  }
  .events-cards-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .market-cards-3col-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .events-grid-layout {
    grid-template-columns: 1fr 1fr;
  }
  .events-center-col {
    display: none;
  }
}

@media (max-width: 768px) {
  .hero-banner-card {
    flex-direction: column;
    height: auto;
  }
  .hero-left-media {
    height: 230px;
  }
  .hero-right-content {
    padding: 24px 20px;
  }
  .hero-main-heading {
    font-size: 1.75rem;
  }
  .events-grid-layout {
    grid-template-columns: 1fr;
  }
  .events-card-wrapper {
    padding: 0;
  }
  .events-promo-header {
    min-height: auto;
    grid-template-columns: minmax(0, 1fr) 210px;
    gap: 16px;
    padding: 26px 24px;
  }
  .events-promo-image {
    width: 210px;
    max-height: 195px;
  }
  .events-period-tabs {
    width: calc(100% - 40px);
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
  .events-period-tabs button {
    padding-inline: 8px;
  }
  .events-list {
    margin-top: 16px;
    padding: 0 20px 18px;
  }
  .event-list-row {
    grid-template-columns: 88px minmax(0, 1fr);
    gap: 14px;
    padding: 16px 0;
  }
  .event-list-image {
    width: 88px;
    height: 76px;
  }
  .event-list-action {
    grid-column: 1 / -1;
    justify-content: space-between;
    padding-left: 102px;
  }
  .event-capacity {
    min-width: 0;
    flex-direction: row;
    gap: 7px;
  }
  .event-capacity strong {
    font-size: 1.375rem;
  }
  .events-section-header {
    align-items: flex-start;
    flex-direction: column;
  }
  .events-cards-grid {
    grid-template-columns: 1fr;
  }
  .event-showcase-media {
    height: 190px;
  }
  .market-cards-3col-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 520px) {
  .events-promo-header {
    display: block;
    padding: 24px 20px;
  }
  .events-promo-copy .events-desc-text {
    margin-bottom: 16px;
  }
  .events-promo-image {
    display: none;
  }
  .events-list {
    padding-inline: 16px;
  }
  .events-list-panel {
    width: calc(100% - 20px);
    margin: 20px auto 0;
    border-radius: 14px;
  }
  .event-list-row {
    grid-template-columns: 76px minmax(0, 1fr);
    gap: 12px;
  }
  .event-list-image {
    width: 76px;
    height: 70px;
  }
  .event-list-content h3 {
    font-size: var(--font-size-md);
  }
  .event-list-meta {
    gap: 6px;
  }
  .event-list-action {
    padding-left: 88px;
  }
  .event-primary-link,
  .event-disabled-link {
    min-height: 38px;
    padding-inline: 10px;
    font-size: var(--font-size-xs);
  }
}
</style>
