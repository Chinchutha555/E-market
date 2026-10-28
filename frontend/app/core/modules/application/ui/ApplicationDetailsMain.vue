<template>
  <div class="app-detail-page">
    <AppHeader />

    <main class="content-body app-detail-container">
      <AppBreadcrumb
        :items="crumbItems"
        :current="application ? `รายละเอียด ${application.application_number || application.application_no || ''}` : 'รายละเอียดคำร้อง'"
        :back-to="backTo"
      />

      <div v-if="loading" class="app-detail-state">กำลังโหลดรายละเอียดคำร้อง...</div>

      <div v-else-if="!application" class="app-detail-state">
        <svg width="38" height="38" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        <h1>ไม่พบคำร้อง</h1>
        <p>คำร้องนี้อาจถูกลบหรือหมายเลขคำร้องไม่ถูกต้อง</p>
      </div>

      <template v-else>
        <section class="app-detail-hero">
          <div class="hero-left">
            <span class="hero-type-tag">{{ isStall ? 'คำร้องขอเช่าแผงค้า' : 'คำร้องทั่วไป' }}</span>
            <h1>{{ application.application_number || application.application_no }}</h1>
            <span class="hero-sub">ยื่นโดย {{ application.applicant_name || '-' }}</span>
          </div>
          <span class="status-badge" :class="statusClass">{{ statusLabel }}</span>
        </section>

        <div class="app-detail-grid">
          <section class="detail-card">
            <h2>ผู้ยื่นคำร้อง</h2>
            <dl>
              <div><dt>ชื่อ-นามสกุล</dt><dd>{{ application.applicant_name || '-' }}</dd></div>
              <div><dt>เบอร์โทรศัพท์</dt><dd>{{ application.applicant_phone || '-' }}</dd></div>
              <div v-if="application.applicant_id_card"><dt>เลขประจำตัวประชาชน</dt><dd>{{ application.applicant_id_card }}</dd></div>
              <div v-if="application.applicant_email"><dt>อีเมล</dt><dd>{{ application.applicant_email }}</dd></div>
              <div v-if="application.applicant_address && application.applicant_address !== '-'" class="full"><dt>ที่อยู่</dt><dd>{{ application.applicant_address }}</dd></div>
            </dl>
          </section>

          <section class="detail-card">
            <h2>{{ isStall ? 'ข้อมูลแผงที่ขอเช่า' : 'ข้อมูลคำร้อง' }}</h2>
            <dl>
              <div v-if="application.market_name"><dt>ตลาด</dt><dd>{{ application.market_name }}</dd></div>
              <div v-if="isStall"><dt>แผงค้า</dt><dd>{{ application.stall_number || application.stall_ident || '-' }}</dd></div>
              <div v-if="isStall"><dt>ชื่อร้านค้า</dt><dd>{{ application.shop_name || '-' }}</dd></div>
              <div v-if="isStall"><dt>ประเภทสินค้า</dt><dd>{{ application.product_type || application.product_category || 'ทั่วไป' }}</dd></div>
              <div v-if="!isStall"><dt>หัวข้อคำร้อง</dt><dd>{{ application.subject || application.shop_name || 'คำร้องทั่วไป' }}</dd></div>
              <div v-if="application.details" class="full"><dt>รายละเอียด</dt><dd>{{ application.details }}</dd></div>
            </dl>
          </section>
        </div>

        <section class="detail-card">
          <h2>บันทึกความเห็นการพิจารณา</h2>
          <ul v-if="comments.length" class="comment-list">
            <li v-for="(c, idx) in comments" :key="idx">
              <div class="comment-head"><strong>{{ c.role }}</strong><span>{{ c.step }}</span></div>
              <p>{{ c.text }}</p>
            </li>
          </ul>
          <p v-else class="empty-text">ยังไม่มีบันทึกความเห็น</p>
        </section>

        <div class="app-detail-actions">
          <NuxtLink :to="backTo" class="back-link">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
            <span>{{ returnUrl ? 'กลับไปหน้าเดิม' : 'กลับไปรายการคำร้อง' }}</span>
          </NuxtLink>
          <NuxtLink :to="`/applications?search=${encodeURIComponent(application.application_number || '')}`" class="manage-link">
            พิจารณา / จัดการคำร้อง
          </NuxtLink>
        </div>
      </template>
    </main>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const { get } = useApi()
const { returnUrl, backTo, crumbItems } = useReturnTo({ label: 'รายการคำร้อง', to: '/applications' })

const loading = ref(true)
const application = ref<any | null>(null)

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
    case 'pending_officer':
    case 'pending_debt_check':
    case 'debt_pending':
    case 'pending_market_director':
    case 'pending_bureau_director': return 'รอตรวจสอบ'
    case 'approved': return 'อนุมัติ'
    case 'revision_required': return 'ส่งกลับแก้ไข'
    case 'rejected': return 'ไม่อนุมัติ'
    default: return status.value || '-'
  }
})
const statusClass = computed(() => {
  if (status.value === 'approved') return 'is-success'
  if (status.value === 'rejected') return 'is-danger'
  return 'is-warning'
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

onMounted(async () => {
  try {
    const res: any = await get(`/applications/${route.params.id}`)
    application.value = res?.id ? res : (res?.data || null)
  } catch {
    application.value = null
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.app-detail-container { display: flex; flex-direction: column; gap: 16px; }
.app-detail-state { display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 60px 16px; color: #64748b; font-size: 0.9375rem; font-weight: 500; }
.app-detail-state h1 { margin: 0; font-size: 1.25rem; color: #1e293b; }
.app-detail-state p { margin: 0; font-size: 0.875rem; }

.app-detail-hero, .detail-card { background: #fff; border: 1px solid #eef2f6; border-radius: 16px; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03); padding: 22px 26px; }
.app-detail-hero { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.hero-left { display: flex; flex-direction: column; gap: 4px; }
.hero-left h1 { margin: 0; font-size: 1.5rem; font-weight: 800; color: #1e293b; }
.hero-type-tag { align-self: flex-start; padding: 3px 10px; border-radius: 6px; background: #f5f3ff; color: #6d28d9; font-size: 0.8125rem; font-weight: 600; }
.hero-sub { font-size: 0.875rem; color: #64748b; font-weight: 500; }

.status-badge { padding: 4px 14px; border-radius: 8px; font-size: 0.8125rem; font-weight: 600; white-space: nowrap; }
.status-badge.is-warning { background: #fef3c7; color: #d97706; }
.status-badge.is-success { background: #dcfce7; color: #15803d; }
.status-badge.is-danger { background: #fee2e2; color: #dc2626; }

.app-detail-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 16px; }
.detail-card h2 { margin: 0 0 14px; font-size: 1.0625rem; font-weight: 700; color: #1e293b; }
.detail-card dl { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px 16px; margin: 0; }
.detail-card dl .full { grid-column: 1 / -1; }
.detail-card dt { font-size: 0.8125rem; color: #64748b; font-weight: 500; }
.detail-card dd { margin: 3px 0 0; font-size: 0.9375rem; color: #334155; font-weight: 600; word-break: break-word; }

.comment-list { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 10px; }
.comment-list li { padding: 10px 14px; border: 1px solid #f1f5f9; border-radius: 10px; background: #f8fafc; }
.comment-head { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.comment-head strong { font-size: 0.875rem; color: #1e293b; }
.comment-head span { font-size: 0.8125rem; color: #64748b; }
.comment-list p { margin: 6px 0 0; font-size: 0.875rem; color: #475569; line-height: 1.5; }
.empty-text { margin: 0; font-size: 0.875rem; color: #94a3b8; }

.app-detail-actions { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.back-link, .manage-link { display: inline-flex; align-items: center; gap: 6px; height: 40px; padding: 0 16px; border-radius: 8px; font-size: 0.9375rem; font-weight: 600; text-decoration: none; transition: all 0.15s ease; }
.back-link { border: 1px solid #cbd5e1; background: #f1f5f9; color: #475569; }
.back-link:hover { background: #e2e8f0; }
.manage-link { border: 1px solid #2563eb; background: #2563eb; color: #fff; box-shadow: 0 2px 8px rgba(37, 99, 235, 0.25); }
.manage-link:hover { background: #1d4ed8; }

@media (max-width: 768px) {
  .app-detail-grid, .detail-card dl { grid-template-columns: 1fr; }
}
</style>
