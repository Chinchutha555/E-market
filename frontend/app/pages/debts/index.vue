<template>
  <div class="debt-page-wrapper">
    <!-- Top Global Navigation Header -->
    <AppHeader />

    <div class="content-body">
      <!-- Page Header Row -->
      <header class="page-header-row">
        <div class="header-left">
          <div class="header-icon-box">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="2" y="4" width="20" height="16" rx="2"></rect>
              <line x1="2" y1="10" x2="22" y2="10"></line>
              <line x1="6" y1="15" x2="6.01" y2="15"></line>
              <line x1="10" y1="15" x2="14" y2="15"></line>
            </svg>
          </div>
          <div>
            <h1 class="page-title">หนี้ค้างชำระ (Tenant Debts)</h1>
            <p class="page-subtitle">บันทึกและตรวจสอบรายการหนี้ค้างชำระของผู้เช่าแผงค้าในระบบ</p>
          </div>
        </div>

        <div class="header-actions">
          <AppButton
            variant="outline"
            size="md"
            :loading="loading"
            @click="loadDebts"
            title="รีเฟรชข้อมูล"
          >
            <template #prefix>
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="23 4 23 10 17 10"></polyline>
                <polyline points="1 20 1 14 7 14"></polyline>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
              </svg>
            </template>
            รีเฟรชข้อมูล
          </AppButton>
        </div>
      </header>

      <!-- Record Debt Form Card (For Staff / Admin) -->
      <section v-if="canCreate" class="debt-card form-card">
        <div class="card-section-header">
          <div class="section-title-wrap">
            <div class="section-icon-badge">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
              </svg>
            </div>
            <div>
              <h2 class="section-card-title">บันทึกหนี้ค้างชำระใหม่</h2>
              <p class="section-card-desc">กรอกข้อมูลผู้เช่าและยอดหนี้เพื่อบันทึกลงในระบบ</p>
            </div>
          </div>
        </div>

        <form @submit.prevent="saveDebt" class="form-grid">
          <div class="form-group">
            <label class="form-label" for="tenant-id-input">
              รหัสผู้เช่า (Tenant ID) <span class="required-star">*</span>
            </label>
            <input
              id="tenant-id-input"
              v-model.trim="form.tenantId"
              type="number"
              min="1"
              class="form-input"
              placeholder="เช่น 1001"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label" for="market-id-input">
              รหัสตลาด (Market ID) <span class="required-star">*</span>
            </label>
            <input
              id="market-id-input"
              v-model.trim="form.marketId"
              type="number"
              min="1"
              class="form-input"
              placeholder="เช่น 1"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label" for="amount-input">
              ยอดหนี้ (บาท) <span class="required-star">*</span>
            </label>
            <div class="input-with-suffix">
              <input
                id="amount-input"
                v-model.trim="form.amount"
                inputmode="decimal"
                class="form-input"
                placeholder="0.00"
                required
              />
              <span class="input-suffix">บาท</span>
            </div>
          </div>

          <div class="form-group wide">
            <label class="form-label" for="desc-input">
              รายละเอียดหนี้ <span class="required-star">*</span>
            </label>
            <input
              id="desc-input"
              v-model.trim="form.description"
              maxlength="500"
              class="form-input"
              placeholder="ระบุรายละเอียด เช่น ค่าเช่าแผงค้างชำระประจำงวด ม.ค. 2569"
              required
            />
          </div>

          <div class="form-action-row wide">
            <AppButton type="submit" variant="primary" size="md" :loading="saving">
              <template #prefix>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path>
                  <polyline points="17 21 17 13 7 13 7 21"></polyline>
                  <polyline points="7 3 7 8 15 8"></polyline>
                </svg>
              </template>
              บันทึกรายการหนี้
            </AppButton>
          </div>
        </form>
      </section>

      <!-- Debts Table Card -->
      <section class="debt-card table-card-container">
        <!-- Table Header: Title + Subtitle on Left, Search on Right -->
        <div class="table-card-header">
          <div class="table-header-left">
            <h2 class="table-card-title">รายการหนี้ค้างชำระทั้งหมด</h2>
            <div class="table-card-subtitle">ทั้งหมด {{ filtered.length }} รายการ</div>
          </div>

          <div class="table-header-actions">
            <!-- Search Box with Search Icon Left -->
            <div class="table-search-box">
              <span class="search-icon">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                </svg>
              </span>
              <input
                v-model.trim="search"
                type="text"
                class="search-input"
                placeholder="ค้นหารหัสผู้เช่า, ตลาด, รายละเอียด..."
                aria-label="ค้นหารายการหนี้"
              />
              <button
                v-if="search"
                type="button"
                class="clear-search-btn"
                @click="search = ''"
                title="ล้างการค้นหา"
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Filter Tabs -->
        <div class="table-nav-tabs">
          <button
            type="button"
            class="table-nav-tab"
            :class="{ active: statusFilter === '' }"
            @click="statusFilter = ''"
          >
            ทั้งหมด ({{ debts.length }})
          </button>
          <button
            type="button"
            class="table-nav-tab"
            :class="{ active: statusFilter === 'outstanding' }"
            @click="statusFilter = 'outstanding'"
          >
            ค้างชำระ ({{ countOutstanding }})
          </button>
          <button
            type="button"
            class="table-nav-tab"
            :class="{ active: statusFilter === 'paid' }"
            @click="statusFilter = 'paid'"
          >
            ชำระแล้ว ({{ countPaid }})
          </button>
        </div>

        <!-- Loading State Skeleton -->
        <div v-if="loading" class="skeleton-wrap" aria-label="กำลังโหลดข้อมูล">
          <div v-for="n in 5" :key="n" class="skeleton-row"></div>
        </div>

        <!-- Table View -->
        <div v-else class="table-wrap">
          <table class="debt-table">
            <thead>
              <tr>
                <th style="width: 80px;">ลำดับ</th>
                <th style="width: 130px;">รหัสผู้เช่า</th>
                <th style="width: 120px;">รหัสตลาด</th>
                <th>รายละเอียดหนี้</th>
                <th style="width: 150px;">วันที่บันทึก</th>
                <th style="width: 160px; text-align: right;">ยอดหนี้</th>
                <th style="width: 130px; text-align: center;">สถานะ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, idx) in paged" :key="item.id">
                <td class="col-num">{{ (page - 1) * pageSize + idx + 1 }}</td>
                <td>
                  <span class="tenant-id-pill">#{{ item.tenant_id }}</span>
                </td>
                <td>
                  <span class="market-id-pill">ตลาด {{ item.market_id }}</span>
                </td>
                <td class="col-desc">{{ item.description }}</td>
                <td class="col-date">{{ formatDate(item.created_at) }}</td>
                <td class="col-amount">
                  {{ (item.amount_satang / 100).toLocaleString('th-TH', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }} <span class="currency-unit">บาท</span>
                </td>
                <td style="text-align: center;">
                  <span
                    class="status-pill"
                    :class="item.status === 'outstanding' || item.status === 'ค้างชำระ' ? 'status-outstanding' : 'status-paid'"
                  >
                    {{ item.status === 'outstanding' ? 'ค้างชำระ' : (item.status === 'paid' ? 'ชำระแล้ว' : item.status) }}
                  </span>
                </td>
              </tr>
              <tr v-if="!filtered.length">
                <td colspan="7" class="empty-cell">
                  <div class="empty-state">
                    <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="empty-icon">
                      <rect x="2" y="4" width="20" height="16" rx="2"></rect>
                      <line x1="2" y1="10" x2="22" y2="10"></line>
                    </svg>
                    <div class="empty-title">ไม่พบรายการหนี้ค้างชำระ</div>
                    <div class="empty-desc">ยังไม่มีประวัติหนี้ค้างที่ตรงตามเงื่อนไขการค้นหา</div>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <AppPagination
          v-if="!loading && filtered.length > 0"
          v-model="page"
          v-model:items-per-page="pageSize"
          :total-items="filtered.length"
        />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'auth' })

const { get, post } = useApi()
const { role } = useAuth()
const { showError, showSuccess } = useAppAlert()

const canCreate = computed(() => ['admin', 'staff'].includes(role.value || ''))
const loading = ref(true)
const saving = ref(false)
const search = ref('')
const statusFilter = ref('')
const page = ref(1)
const pageSize = ref(10)

const debts = ref<any[]>([])
const form = reactive({
  tenantId: '',
  marketId: '',
  amount: '',
  description: ''
})

const countOutstanding = computed(() => debts.value.filter(d => d.status === 'outstanding').length)
const countPaid = computed(() => debts.value.filter(d => d.status === 'paid').length)

const filtered = computed(() => {
  return debts.value.filter(item => {
    // Status filter
    if (statusFilter.value && item.status !== statusFilter.value) {
      return false
    }
    // Search query
    if (search.value) {
      const q = search.value.toLowerCase()
      const match = `${item.tenant_id} ${item.market_id} ${item.description}`.toLowerCase().includes(q)
      if (!match) return false
    }
    return true
  })
})

const paged = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filtered.value.slice(start, start + pageSize.value)
})

watch([search, statusFilter, pageSize], () => {
  page.value = 1
})

function formatDate(iso?: string) {
  if (!iso) return '-'
  try {
    const d = new Date(iso)
    if (isNaN(d.getTime())) return '-'
    return d.toLocaleDateString('th-TH', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  } catch {
    return '-'
  }
}

async function loadDebts() {
  loading.value = true
  try {
    const res: any = await get('/debts')
    debts.value = Array.isArray(res?.data) ? res.data : []
  } catch {
    await showError('ไม่สามารถโหลดรายการหนี้ได้', 'เกิดข้อผิดพลาด')
  } finally {
    loading.value = false
  }
}

async function saveDebt() {
  if (!/^(?:0|[1-9]\d*)(?:\.\d{1,2})?$/.test(form.amount)) {
    await showError('กรุณาระบุยอดหนี้เป็นตัวเลขไม่เกินสองตำแหน่ง', 'ข้อมูลไม่ถูกต้อง')
    return
  }
  const [whole, decimal = ''] = form.amount.split('.')
  const amount = Number(whole) * 100 + Number((decimal + '00').slice(0, 2))
  if (
    !Number.isSafeInteger(amount) ||
    amount <= 0 ||
    !Number.isSafeInteger(Number(form.tenantId)) ||
    !Number.isSafeInteger(Number(form.marketId))
  ) {
    await showError('กรุณาตรวจสอบยอดหนี้และรหัสที่ระบุ', 'ข้อมูลไม่ถูกต้อง')
    return
  }

  saving.value = true
  try {
    await post('/debts', {
      tenant_id: Number(form.tenantId),
      market_id: Number(form.marketId),
      description: form.description,
      amount_satang: amount
    })
    form.amount = ''
    form.description = ''
    await loadDebts()
    await showSuccess('บันทึกรายการหนี้เรียบร้อยแล้ว', 'สำเร็จ')
  } catch (err: any) {
    await showError(
      err?.data?.error?.message ||
        err?.data?.message ||
        'ไม่สามารถบันทึกหนี้ได้ กรุณาตรวจสอบสิทธิ์และข้อมูลผู้เช่า',
      'เกิดข้อผิดพลาด'
    )
  } finally {
    saving.value = false
  }
}

onMounted(loadDebts)
</script>

<style scoped>
.debt-page-wrapper {
  min-height: 100vh;
  background-color: var(--bg-main, #f8fafc);
}

.content-body {
  padding: 24px 32px 48px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* Page Header Row */
.page-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.header-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.08);
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  line-height: 1.3;
}

.page-subtitle {
  font-size: 0.875rem;
  color: #64748b;
  margin: 4px 0 0;
}

/* Card Container Standards */
.debt-card {
  background: #ffffff;
  border: 1px solid #eef2f6;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
  padding: 24px 28px;
}

/* Form Card */
.form-card {
  border-left: 4px solid #2563eb;
}

.card-section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.section-icon-badge {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
}

.section-card-title {
  font-size: 1.0625rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.section-card-desc {
  font-size: 0.8125rem;
  color: #64748b;
  margin: 2px 0 0;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group.wide {
  grid-column: 1 / -1;
}

.form-label {
  font-size: 0.84375rem;
  font-weight: 600;
  color: #334155;
}

.required-star {
  color: #ef4444;
}

.form-input {
  height: 40px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 0 12px;
  font-family: inherit;
  font-size: 0.875rem;
  color: #1e293b;
  background: #ffffff;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
  width: 100%;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.12);
}

.input-with-suffix {
  position: relative;
  display: flex;
  align-items: center;
}

.input-with-suffix .form-input {
  padding-right: 48px;
}

.input-suffix {
  position: absolute;
  right: 12px;
  font-size: 0.8125rem;
  color: #64748b;
  font-weight: 500;
  pointer-events: none;
}

.form-action-row {
  display: flex;
  justify-content: flex-end;
  margin-top: 6px;
}

/* Table Card Header */
.table-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 18px;
}

.table-card-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.table-card-subtitle {
  font-size: 0.8125rem;
  color: #64748b;
  margin: 3px 0 0;
}

.table-header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.table-search-box {
  position: relative;
  display: flex;
  align-items: center;
  width: 280px;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: #94a3b8;
  display: flex;
  align-items: center;
  pointer-events: none;
}

.search-input {
  width: 100%;
  height: 38px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 0 32px 0 36px;
  font-family: inherit;
  font-size: 0.84375rem;
  color: #1e293b;
  background: #ffffff;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
  box-sizing: border-box;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.12);
}

.clear-search-btn {
  position: absolute;
  right: 10px;
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.clear-search-btn:hover {
  color: #475569;
}

/* Filter Tabs */
.table-nav-tabs {
  display: flex;
  gap: 24px;
  border-bottom: 1px solid #f1f5f9;
  margin-bottom: 16px;
}

.table-nav-tab {
  background: none;
  border: none;
  padding: 10px 4px;
  font-size: 0.875rem;
  font-weight: 500;
  color: #64748b;
  cursor: pointer;
  position: relative;
  transition: color 0.15s ease;
  font-family: inherit;
}

.table-nav-tab:hover {
  color: #1e293b;
}

.table-nav-tab.active {
  color: #2563eb;
  font-weight: 600;
}

.table-nav-tab.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 3px;
  background: #2563eb;
  border-radius: 2px;
}

/* Table Design Standards */
.table-wrap {
  overflow-x: auto;
  border-radius: 8px;
}

.debt-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.debt-table thead th {
  background: transparent;
  font-size: 0.8125rem;
  font-weight: 500;
  color: #64748b;
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
  white-space: nowrap;
}

.debt-table tbody td {
  font-size: 0.875rem;
  color: #334155;
  padding: 16px 16px;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}

.debt-table tbody tr:hover {
  background: #f8fafc;
}

.col-num {
  color: #94a3b8;
  font-size: 0.8125rem;
}

.tenant-id-pill {
  font-weight: 600;
  color: #1e293b;
  background: #f1f5f9;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 0.8125rem;
}

.market-id-pill {
  color: #475569;
  font-size: 0.84375rem;
  font-weight: 500;
}

.col-desc {
  max-width: 320px;
  word-break: break-word;
  color: #334155;
}

.col-date {
  color: #64748b;
  font-size: 0.8125rem;
  white-space: nowrap;
}

.col-amount {
  font-weight: 700;
  color: #1e293b;
  text-align: right;
  white-space: nowrap;
  font-size: 0.9375rem;
}

.currency-unit {
  font-size: 0.8125rem;
  font-weight: 500;
  color: #64748b;
}

/* Status Soft Pills */
.status-pill {
  display: inline-block;
  border-radius: 6px;
  padding: 4px 12px;
  font-size: 0.8125rem;
  font-weight: 600;
  white-space: nowrap;
}

.status-outstanding {
  background: #fee2e2;
  color: #dc2626;
}

.status-paid {
  background: #dcfce7;
  color: #15803d;
}

/* Empty State */
.empty-cell {
  text-align: center;
  padding: 48px 16px !important;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.empty-icon {
  color: #cbd5e1;
  margin-bottom: 4px;
}

.empty-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #475569;
}

.empty-desc {
  font-size: 0.8125rem;
  color: #94a3b8;
}

/* Skeleton Loading */
.skeleton-wrap {
  padding: 8px 0;
}

.skeleton-row {
  height: 48px;
  margin: 10px 0;
  border-radius: 8px;
  background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

@media (max-width: 840px) {
  .content-body {
    padding: 16px;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }

  .page-header-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .table-card-header {
    flex-direction: column;
    align-items: stretch;
  }

  .table-search-box {
    width: 100%;
  }
}
</style>