<template>
  <footer class="app-bma-footer">
    <!-- Quick Assistance Banner -->
    <div class="footer-assistance-bar">
      <div class="footer-container assistance-inner">
        <div class="assistance-info">
          <div class="assistance-icon">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
            </svg>
          </div>
          <div>
            <div class="assistance-title">ต้องการสอบถามข้อมูลพื้นที่เช่าหรือขอรับความช่วยเหลือ?</div>
            <div class="assistance-subtitle">ติดต่อสำนักงานตลาดกรุงเทพมหานคร หรือ โทรสายด่วน กทม. 1555 (บริการ 24 ชั่วโมง)</div>
          </div>
        </div>
        <div class="assistance-actions">
          <a href="tel:021598888" class="btn-footer-call">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
            </svg>
            <span>โทร 02-159-8888</span>
          </a>
        </div>
      </div>
    </div>

    <!-- Main Footer Body -->
    <div class="footer-main-content">
      <div class="footer-container footer-columns-grid">
        <!-- Col 1: Organization & Identity -->
        <div class="footer-col col-brand">
          <div class="footer-brand">
            <div class="brand-emblem-badge">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path>
                <line x1="3" y1="6" x2="21" y2="6"></line>
                <path d="M16 10a4 4 0 0 1-8 0"></path>
              </svg>
            </div>
            <div>
              <div class="brand-org-name">สำนักงานตลาดกรุงเทพมหานคร</div>
              <div class="brand-org-sub">Bangkok Metropolitan Administration Market Office</div>
            </div>
          </div>
          <p class="brand-desc">
            ระบบบริการประชาชนและผู้ประกอบการออนไลน์ (e-Market) แพลตฟอร์มบริหารจัดการแผงค้าและตลาดนัดในสังกัดกรุงเทพมหานคร 12 แห่ง โปร่งใส สะดวก รวดเร็ว และเป็นธรรม
          </p>
          <div class="footer-office-address">
            <div class="office-item">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                <circle cx="12" cy="10" r="3"></circle>
              </svg>
              <span>8 ถนนเทศบาลสงเคราะห์ แขวงลาดยาว เขตจตุจักร กรุงเทพมหานคร 10900</span>
            </div>
            <div class="office-item">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
              </svg>
              <span>วันจันทร์ - ศุกร์: 08:30 - 16:30 น. (เว้นวันหยุดราชการ)</span>
            </div>
          </div>
        </div>

        <!-- Col 2: Services -->
        <div class="footer-col">
          <h4 class="footer-heading">ระบบบริการออนไลน์</h4>
          <ul class="footer-nav-list">
            <li>
              <NuxtLink to="/#markets-12-directory" class="footer-nav-link">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
                <span>ค้นหาแผงค้า 12 ตลาด กทม.</span>
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/stalls/book?market=chatuchak" class="footer-nav-link">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
                <span>ผังดิจิทัล GIS ตลาดนัดจตุจักร</span>
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/bookings" class="footer-nav-link">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
                <span>ตรวจสอบประวัติการจอง</span>
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/applications" class="footer-nav-link">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
                <span>ติดตามสถานะคำร้องขอเช่า</span>
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/contracts" class="footer-nav-link">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
                <span>สัญญาเช่าและใบแจ้งหนี้</span>
              </NuxtLink>
            </li>
          </ul>
        </div>

        <!-- Col 3: 12 BMA Markets -->
        <div class="footer-col">
          <h4 class="footer-heading">ตลาดในสังกัด กทม.</h4>
          <ul class="footer-nav-list markets-list-grid">
            <li>
              <NuxtLink to="/stalls/book?market=chatuchak" class="footer-nav-link">
                <span class="mkt-bullet"></span>
                <span>ตลาดนัดจตุจักร</span>
              </NuxtLink>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดมีนบุรี</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดธนบุรี (สนามหลวง 2)</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดประชานิเวศน์ 1</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดเทวราช</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดบางกะปิ</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดหนองจอก</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดพระเครื่องวงเวียนเล็ก</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดราษฎร์บูรณะ</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดรัชดาภิเษก</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดสิงหา</span>
              </span>
            </li>
            <li>
              <span class="footer-nav-static">
                <span class="mkt-bullet"></span>
                <span>ตลาดบางแคภิรมย์</span>
              </span>
            </li>
          </ul>
        </div>

        <!-- Col 4: Contact & Security Badges -->
        <div class="footer-col">
          <h4 class="footer-heading">ติดต่อสำนักงาน & ความปลอดภัย</h4>
          <div class="footer-contact-items">
            <div class="contact-line">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
              </svg>
              <span>โทร: 02-159-8888, 02-272-4440</span>
            </div>
            <div class="contact-line">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
                <polyline points="22,6 12,13 2,6"></polyline>
              </svg>
              <span>อีเมล: contact@market.bangkok.go.th</span>
            </div>
            <div class="contact-line">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
              </svg>
              <span>สายด่วน กทม. 1555 (ตลอด 24 ชม.)</span>
            </div>
          </div>

          <!-- Security Badges -->
          <div class="footer-security-box">
            <div class="security-badge-item">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
              <span>เข้ารหัสข้อมูล 256-bit SSL</span>
            </div>
            <div class="security-badge-item">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#10b981" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              <span>มาตรฐานความปลอดภัย PDPA</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Copyright Bar -->
    <div class="footer-bottom-bar">
      <div class="footer-container bottom-bar-inner">
        <div class="copyright-text">
          © 2026 สำนักงานตลาดกรุงเทพมหานคร (Bangkok Metropolitan Administration Market Office). สงวนลิขสิทธิ์
        </div>
        <div class="legal-links">
          <span class="legal-link">นโยบายความเป็นส่วนตัว</span>
          <span class="legal-sep">•</span>
          <span class="legal-link">ข้อกำหนดและเงื่อนไข</span>
          <span class="legal-sep">•</span>
          <span class="legal-link">ความคุ้มครองข้อมูลส่วนบุคคล</span>
          <span class="legal-sep">•</span>
          <span class="system-version-pill">v2.5.0 Production</span>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
// Standalone BMA e-Market Footer Component
</script>

<style scoped>
.app-bma-footer {
  width: 100%;
  background: var(--bg-ffffff);
  color: var(--tx-475569);
  font-family: 'Bai Jamjuree', sans-serif;
  margin-top: 56px;
  border-top: 1px solid var(--bd-e2e8f0);
  box-shadow: 0 -4px 20px rgba(15, 23, 42, 0.02);
}

.footer-container {
  max-width: 1320px;
  margin: 0 auto;
  padding: 0 32px;
}

/* 1. Quick Assistance Bar */
.footer-assistance-bar {
  background: linear-gradient(135deg, var(--bg-eff6ff) 0%, var(--bg-f8fafc) 100%);
  color: var(--tx-0f172a);
  padding: 22px 0;
  border-bottom: 1px solid var(--bd-e2e8f0);
}

.assistance-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  flex-wrap: wrap;
}

.assistance-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.assistance-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: var(--bg-dbeafe);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--tx-2563eb);
  border: 1px solid var(--bd-bfdbfe);
}

.assistance-title {
  font-size: var(--font-size-md);
  font-weight: 700;
  color: var(--tx-0f172a);
  margin-bottom: 3px;
}

.assistance-subtitle {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
}

.assistance-actions {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.btn-footer-call {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #2563eb;
  color: #ffffff;
  padding: 9px 20px;
  border-radius: 9999px;
  font-size: var(--font-size-sm);
  font-weight: 700;
  text-decoration: none;
  transition: all 0.18s ease;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.25);
}

.btn-footer-call:hover {
  background: #1d4ed8;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.35);
}

/* 2. Main Content Grid */
.footer-main-content {
  background: var(--bg-ffffff);
  padding: 52px 0 44px;
}

.footer-columns-grid {
  display: grid;
  grid-template-columns: 2fr 1.2fr 1.4fr 1.4fr;
  gap: 40px;
}

.footer-col {
  display: flex;
  flex-direction: column;
}

.col-brand {
  padding-right: 12px;
}

.footer-brand {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 16px;
}

.brand-emblem-badge {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: var(--bg-eff6ff);
  border: 1.5px solid var(--bd-bfdbfe);
  color: var(--tx-2563eb);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 10px rgba(37, 99, 235, 0.12);
}

.brand-org-name {
  font-size: var(--font-size-lg);
  font-weight: 800;
  color: var(--tx-0f172a);
  letter-spacing: -0.01em;
}

.brand-org-sub {
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
  font-weight: 600;
  margin-top: 1px;
}

.brand-desc {
  font-size: var(--font-size-sm);
  color: var(--tx-64748b);
  line-height: 1.6;
  margin: 0 0 20px;
}

.footer-office-address {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.office-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font-size: var(--font-size-xs);
  color: var(--tx-475569);
  line-height: 1.5;
}

.office-item svg {
  color: var(--tx-2563eb);
  flex-shrink: 0;
  margin-top: 3px;
}

.footer-heading {
  font-size: var(--font-size-base);
  font-weight: 700;
  color: var(--tx-0f172a);
  margin: 0 0 18px;
  position: relative;
  padding-bottom: 10px;
}

.footer-heading::after {
  content: '';
  position: absolute;
  left: 0;
  bottom: 0;
  width: 28px;
  height: 3px;
  background: #2563eb;
  border-radius: 2px;
}

.footer-nav-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.footer-nav-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--tx-475569);
  text-decoration: none;
  font-size: var(--font-size-sm);
  transition: all 0.16s ease;
}

.footer-nav-link svg {
  color: var(--tx-2563eb);
  transition: transform 0.16s ease;
}

.footer-nav-link:hover {
  color: var(--tx-2563eb);
  transform: translateX(3px);
}

.footer-nav-static {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--tx-64748b);
  font-size: var(--font-size-xs);
}

.markets-list-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px 12px;
}

.mkt-bullet {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--bg-93c5fd);
  flex-shrink: 0;
}

.footer-contact-items {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 18px;
}

.contact-line {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: var(--font-size-xs);
  color: var(--tx-334155);
}

.contact-line svg {
  color: var(--tx-2563eb);
  flex-shrink: 0;
}

.footer-security-box {
  background: var(--bg-f8fafc);
  border: 1px solid var(--bd-e2e8f0);
  border-radius: 10px;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.security-badge-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: var(--font-size-xs);
  color: var(--tx-334155);
  font-weight: 500;
}

/* 3. Bottom Bar */
.footer-bottom-bar {
  background: var(--bg-f8fafc);
  border-top: 1px solid var(--bd-e2e8f0);
  padding: 18px 0;
  font-size: var(--font-size-xs);
  color: var(--tx-64748b);
}

.bottom-bar-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.copyright-text {
  color: var(--tx-64748b);
}

.legal-links {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.legal-link {
  color: var(--tx-64748b);
  cursor: pointer;
  transition: color 0.15s ease;
}

.legal-link:hover {
  color: var(--tx-0f172a);
}

.legal-sep {
  color: var(--tx-cbd5e1);
}

.system-version-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 9px;
  border-radius: 6px;
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-bfdbfe);
  color: var(--tx-2563eb);
  font-size: var(--font-size-xs);
  font-weight: 700;
}

/* Responsive adjustments */
@media (max-width: 1080px) {
  .footer-columns-grid {
    grid-template-columns: 1fr 1fr;
    gap: 36px;
  }
}

@media (max-width: 640px) {
  .footer-columns-grid {
    grid-template-columns: 1fr;
    gap: 32px;
  }
  .footer-container {
    padding: 0 20px;
  }
  .assistance-inner {
    flex-direction: column;
    align-items: flex-start;
  }
  .bottom-bar-inner {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
