<template>
  <div>
    <AppHeader />

    <div class="content-body">
      <!-- Breadcrumb Navigation -->
      <div class="page-breadcrumb-nav">
        <NuxtLink to="/applications" class="breadcrumb-back-link">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
          <span>กลับไปยังรายการคำร้อง</span>
        </NuxtLink>
        <span class="breadcrumb-separator">/</span>
        <span class="breadcrumb-current">จัดการประเภทคำร้อง</span>
      </div>

      <!-- Main Card Container -->
      <div class="table-card-container">
        <!-- Header -->
        <div class="table-card-header">
          <div class="table-header-left">
            <div class="header-icon-title-wrap">
              <div class="header-icon-box">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                  <polyline points="14 2 14 8 20 8"></polyline>
                  <line x1="16" y1="13" x2="8" y2="13"></line>
                  <line x1="16" y1="17" x2="8" y2="17"></line>
                  <polyline points="10 9 9 9 8 9"></polyline>
                </svg>
              </div>
              <div>
                <h1 class="table-card-title">ประเภทคำร้อง / เรื่องที่ขอ (Application Types)</h1>
                <div class="table-card-subtitle">
                  จัดการหัวข้อประเภทคำร้องทั่วไปสำหรับให้ผู้ค้าและประชาชนเลือกยื่นคำร้องในระบบ (ทั้งหมด {{ typesList.length }} รายการ)
                </div>
              </div>
            </div>
          </div>

          <div class="table-header-actions">
            <!-- Search Box -->
            <div class="table-search-box">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="11" cy="11" r="8"></circle>
                <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
              </svg>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="ค้นหาชื่อเรื่อง / คำอธิบาย..."
                @input="currentPage = 1"
              />
            </div>

            <!-- Reload Button -->
            <button class="btn-refresh" @click="loadTypes" title="รีเฟรชข้อมูล">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="23 4 23 10 17 10"></polyline>
                <polyline points="1 20 1 14 7 14"></polyline>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
              </svg>
            </button>

            <!-- Add New Type CTA (Admin Only) -->
            <AppButton
              v-if="isAdmin"
              variant="primary"
              size="md"
              @click="openCreateModal"
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 6px;">
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
              </svg>
              <span>เพิ่มประเภทคำร้องใหม่</span>
            </AppButton>
          </div>
        </div>

        <!-- Filter Tabs -->
        <div class="filter-tabs-row">
          <button
            type="button"
            class="tab-btn"
            :class="{ active: statusFilter === 'all' }"
            @click="statusFilter = 'all'; currentPage = 1"
          >
            <span>ทั้งหมด</span>
            <span class="tab-count-pill">{{ typesList.length }}</span>
          </button>
          <button
            type="button"
            class="tab-btn"
            :class="{ active: statusFilter === 'active' }"
            @click="statusFilter = 'active'; currentPage = 1"
          >
            <span>เปิดใช้งานอยู่</span>
            <span class="tab-count-pill pill-success">{{ activeCount }}</span>
          </button>
          <button
            type="button"
            class="tab-btn"
            :class="{ active: statusFilter === 'inactive' }"
            @click="statusFilter = 'inactive'; currentPage = 1"
          >
            <span>ปิดการใช้งาน</span>
            <span class="tab-count-pill">{{ inactiveCount }}</span>
          </button>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="table-loading-wrap">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="spin-icon">
            <line x1="12" y1="2" x2="12" y2="6"></line>
            <line x1="12" y1="18" x2="12" y2="22"></line>
            <line x1="4.93" y1="4.93" x2="7.76" y2="7.76"></line>
            <line x1="16.24" y1="16.24" x2="19.07" y2="19.07"></line>
            <line x1="2" y1="12" x2="6" y2="12"></line>
            <line x1="18" y1="12" x2="22" y2="12"></line>
            <line x1="4.93" y1="19.07" x2="7.76" y2="16.24"></line>
            <line x1="16.24" y1="7.76" x2="19.07" y2="4.93"></line>
          </svg>
          <span>กำลังโหลดข้อมูลประเภทคำร้อง...</span>
        </div>

        <!-- Table View -->
        <div v-else-if="paginatedTypes.length > 0" class="types-table-wrap">
          <table class="types-table">
            <thead>
              <tr>
                <th style="width: 70px; text-align: center;">ลำดับ</th>
                <th>ประเภทคำร้อง / เรื่องที่ขอ</th>
                <th>คำอธิบายเพิ่มเติม</th>
                <th style="width: 140px; text-align: center;">สถานะ</th>
                <th v-if="isAdmin" style="width: 160px; text-align: right;">การจัดการ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="t in paginatedTypes" :key="t.id">
                <td style="text-align: center;">
                  <span class="sort-order-badge">{{ t.sort_order }}</span>
                </td>
                <td>
                  <div class="type-title-text">{{ t.title }}</div>
                </td>
                <td>
                  <div class="type-desc-text">{{ t.description || '-' }}</div>
                </td>
                <td style="text-align: center;">
                  <span
                    class="status-soft-pill"
                    :class="t.is_active ? 'pill-active' : 'pill-inactive'"
                  >
                    <span class="status-dot"></span>
                    <span>{{ t.is_active ? 'เปิดใช้งาน' : 'ปิดใช้งาน' }}</span>
                  </span>
                </td>
                <td v-if="isAdmin" style="text-align: right;">
                  <div class="table-actions-cell">
                    <!-- Toggle Active Button -->
                    <button
                      class="btn-action-icon"
                      :title="t.is_active ? 'คลิกเพื่อปิดใช้งาน' : 'คลิกเพื่อเปิดใช้งาน'"
                      @click="toggleTypeStatus(t)"
                    >
                      <svg v-if="t.is_active" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: #059669;">
                        <path d="M18.36 6.64a9 9 0 1 1-12.73 0"></path>
                        <line x1="12" y1="2" x2="12" y2="12"></line>
                      </svg>
                      <svg v-else width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: #94a3b8;">
                        <circle cx="12" cy="12" r="10"></circle>
                        <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"></line>
                      </svg>
                    </button>

                    <!-- Edit Button -->
                    <button
                      class="btn-action-outline btn-edit"
                      title="แก้ไขข้อมูล"
                      @click="openEditModal(t)"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                      </svg>
                      <span>แก้ไข</span>
                    </button>

                    <!-- Delete Button -->
                    <button
                      class="btn-action-outline btn-delete"
                      title="ลบประเภทคำร้องนี้"
                      @click="confirmDelete(t)"
                    >
                      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="3 6 5 6 21 6"></polyline>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                      </svg>
                      <span>ลบ</span>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Empty State -->
        <div v-else class="empty-state-box">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="16" y1="13" x2="8" y2="13"></line>
            <line x1="16" y1="17" x2="8" y2="17"></line>
          </svg>
          <div class="empty-title">ไม่พบประเภทคำร้องตามเงื่อนไขที่ค้นหา</div>
          <div class="empty-subtitle">ลองเปลี่ยนคำค้นหา หรือกดปุ่มด้านล่างเพื่อเพิ่มประเภทคำร้องใหม่</div>
          <AppButton v-if="isAdmin" variant="primary" size="sm" style="margin-top: 14px;" @click="openCreateModal">
            เพิ่มประเภทคำร้องใหม่
          </AppButton>
        </div>

        <!-- Pagination Footer -->
        <AppPagination
          v-if="filteredTypes.length > 0"
          v-model="currentPage"
          v-model:items-per-page="pageSize"
          :total-items="filteredTypes.length"
          unit-label="ประเภทคำร้อง"
        />
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <div
      v-if="showModal"
      class="modal-overlay"
      @click.self="closeModal"
    >
      <div class="modal-content type-modal-box">
        <!-- Modal Header -->
        <div class="modal-header">
          <div class="modal-header-title-wrap">
            <div class="modal-icon-badge modal-icon-primary">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="12" y1="18" x2="12" y2="12"></line>
                <line x1="9" y1="15" x2="15" y2="15"></line>
              </svg>
            </div>
            <div>
              <h3 class="modal-title">{{ editingItem ? 'แก้ไขประเภทคำร้อง' : 'เพิ่มประเภทคำร้องใหม่' }}</h3>
              <div class="modal-subtitle">{{ editingItem ? 'ปรับปรุงหัวข้อและรายละเอียดประเภทคำร้อง' : 'กำหนดหัวข้อคำร้องสำหรับให้ผู้ใช้เลือกยื่นคำร้องทั่วไป' }}</div>
            </div>
          </div>
          <button class="btn-modal-close" @click="closeModal" aria-label="ปิด">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Modal Body Form -->
        <div class="modal-body">
          <div v-if="formError" class="alert-error-box">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
            <span>{{ formError }}</span>
          </div>

          <!-- Title -->
          <div class="form-group">
            <label class="form-label" for="type-title">
              ชื่อประเภทคำร้อง / เรื่องที่ขอ <span class="text-danger">*</span>
            </label>
            <input
              id="type-title"
              v-model="form.title"
              type="text"
              class="form-input"
              placeholder="เช่น ขออนุญาตปรับปรุง/ต่อเติมแผงค้า"
            />
          </div>

          <!-- Sort Order -->
          <div class="form-group">
            <label class="form-label" for="type-sort">
              ลำดับการแสดงผล (Sort Order)
            </label>
            <input
              id="type-sort"
              v-model.number="form.sort_order"
              type="number"
              min="0"
              class="form-input"
              placeholder="0"
            />
          </div>

          <!-- Description -->
          <div class="form-group">
            <label class="form-label" for="type-desc">
              คำอธิบายรายละเอียดประเภทคำร้อง
            </label>
            <textarea
              id="type-desc"
              v-model="form.description"
              rows="3"
              class="form-textarea"
              placeholder="ระบุคำอธิบาย หรือเงื่อนไขเพิ่มเติมของคำร้องประเภทนี้..."
            ></textarea>
          </div>

          <!-- Active Status Toggle -->
          <div class="form-group status-toggle-wrap">
            <label class="toggle-checkbox-label">
              <AppCheckbox v-model="form.is_active" />
              <div class="toggle-text-block">
                <span class="toggle-title">เปิดใช้งานประเภทคำร้องนี้</span>
                <span class="toggle-desc">เมื่อเปิดใช้งาน ผู้ใช้จะสามารถเลือกหัวข้อนี้ได้ในหน้าต่างยื่นคำร้องทั่วไป</span>
              </div>
            </label>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="modal-footer">
          <AppButton
            variant="secondary"
            size="md"
            :disabled="submitting"
            @click="closeModal"
          >
            ยกเลิก
          </AppButton>
          <AppButton
            variant="primary"
            size="md"
            :loading="submitting"
            @click="saveType"
          >
            {{ editingItem ? 'บันทึกการแก้ไข' : 'สร้างประเภทคำร้อง' }}
          </AppButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'

definePageMeta({
  middleware: () => {
    if (!import.meta.client) return
    const auth = useAuthStore()
    auth.hydrate()
    if (auth.role !== 'admin') {
      return navigateTo(auth.role === 'staff' ? '/staff' : '/applications', { replace: true })
    }
  }
})

const authStore = useAuthStore()
const { isAdmin } = useAuth()
const { get, post, put, del } = useApi()
const { showConfirm, showSuccess, showError } = useAppAlert()

const loading = ref(true)
const submitting = ref(false)
const typesList = ref<any[]>([])
const searchQuery = ref('')
const statusFilter = ref<'all' | 'active' | 'inactive'>('all')
const currentPage = ref(1)
const pageSize = ref(10)

// Modal State
const showModal = ref(false)
const editingItem = ref<any>(null)
const formError = ref('')
const form = ref({
  title: '',
  code: '',
  description: '',
  sort_order: 1,
  is_active: true
})

// Counts
const activeCount = computed(() => typesList.value.filter(t => t.is_active).length)
const inactiveCount = computed(() => typesList.value.filter(t => !t.is_active).length)

// Filtered Types
const filteredTypes = computed(() => {
  return typesList.value.filter((item) => {
    // Status filter
    if (statusFilter.value === 'active' && !item.is_active) return false
    if (statusFilter.value === 'inactive' && item.is_active) return false

    // Search query
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.trim().toLowerCase()
      const title = String(item.title || '').toLowerCase()
      const desc = String(item.description || '').toLowerCase()
      if (!title.includes(q) && !desc.includes(q)) {
        return false
      }
    }
    return true
  })
})

// Paginated Types
const paginatedTypes = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredTypes.value.slice(start, start + pageSize.value)
})

// Fetch Types
const loadTypes = async () => {
  loading.value = true
  try {
    const res = await get<any>('/application-types')
    if (res.status && res.data) {
      typesList.value = Array.isArray(res.data) ? res.data : []
    }
  } catch (err: any) {
    await showError(err?.data?.detail || 'ไม่สามารถโหลดข้อมูลประเภทคำร้องได้', 'ข้อผิดพลาด')
  } finally {
    loading.value = false
  }
}

// Open Create Modal
const openCreateModal = () => {
  editingItem.value = null
  const nextOrder = (typesList.value.length || 0) + 1
  form.value = {
    title: '',
    code: '',
    description: '',
    sort_order: nextOrder,
    is_active: true
  }
  formError.value = ''
  showModal.value = true
}

// Open Edit Modal
const openEditModal = (item: any) => {
  editingItem.value = item
  form.value = {
    title: item.title || '',
    code: item.code || '',
    description: item.description || '',
    sort_order: item.sort_order ?? 0,
    is_active: Boolean(item.is_active)
  }
  formError.value = ''
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingItem.value = null
  formError.value = ''
}

// Save Type (Create or Update)
const saveType = async () => {
  formError.value = ''
  if (!form.value.title.trim()) {
    formError.value = 'กรุณาระบุชื่อประเภทคำร้อง / เรื่องที่ขอ'
    return
  }

  submitting.value = true
  try {
    if (editingItem.value) {
      await put(`/application-types/${editingItem.value.id}`, {
        title: form.value.title.trim(),
        code: form.value.code ? form.value.code.trim() : undefined,
        description: form.value.description.trim(),
        sort_order: Number(form.value.sort_order) || 0,
        is_active: form.value.is_active
      })
      await showSuccess('บันทึกการแก้ไขประเภทคำร้องเรียบร้อยแล้ว', 'สำเร็จ')
    } else {
      await post('/application-types', {
        title: form.value.title.trim(),
        code: form.value.code ? form.value.code.trim() : undefined,
        description: form.value.description.trim(),
        sort_order: Number(form.value.sort_order) || 0,
        is_active: form.value.is_active
      })
      await showSuccess('เพิ่มประเภทคำร้องใหม่เรียบร้อยแล้ว', 'สำเร็จ')
    }
    closeModal()
    await loadTypes()
  } catch (err: any) {
    formError.value = err?.data?.detail || err?.data?.message || 'เกิดข้อผิดพลาดในการบันทึกข้อมูล'
  } finally {
    submitting.value = false
  }
}

// Toggle Active Status
const toggleTypeStatus = async (item: any) => {
  const nextStatusText = item.is_active ? 'ปิดการใช้งาน' : 'เปิดใช้งาน'
  const confirmed = await showConfirm(
    `คุณต้องการ ${nextStatusText} ประเภทคำร้อง "${item.title}" ใช่หรือไม่?`,
    {
      title: `ยืนยันการ${nextStatusText}`,
      type: item.is_active ? 'warning' : 'info',
      confirmText: nextStatusText,
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return

  try {
    await post(`/application-types/${item.id}/toggle`, {})
    item.is_active = !item.is_active
    await showSuccess(`เปลี่ยนสถานะเป็น ${nextStatusText} เรียบร้อยแล้ว`, 'สำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาดในการเปลี่ยนสถานะ', 'ข้อผิดพลาด')
  }
}

// Confirm Delete
const confirmDelete = async (item: any) => {
  const confirmed = await showConfirm(
    `คุณต้องการลบประเภทคำร้อง "${item.title}" (${item.code}) ออกจากระบบอย่างถาวรใช่หรือไม่?`,
    {
      title: 'ยืนยันการลบประเภทคำร้อง',
      type: 'danger',
      confirmText: 'ยืนยันการลบ',
      cancelText: 'ยกเลิก'
    }
  )
  if (!confirmed) return

  try {
    await del(`/application-types/${item.id}`)
    await showSuccess('ลบประเภทคำร้องเรียบร้อยแล้ว', 'ลบสำเร็จ')
    await loadTypes()
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาดในการลบข้อมูล', 'ข้อผิดพลาด')
  }
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    authStore.hydrate()
  }
  loadTypes()
})
</script>

<style scoped>
/* Breadcrumb */
.page-breadcrumb-nav {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-size: 0.8125rem; /* 13px */
  color: #64748b;
}

.breadcrumb-back-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #2563eb;
  font-weight: 500;
  text-decoration: none;
  transition: color 0.15s ease;
}

.breadcrumb-back-link:hover {
  color: #1d4ed8;
  text-decoration: underline;
}

.breadcrumb-separator {
  color: #cbd5e1;
}

.breadcrumb-current {
  color: #1e293b;
  font-weight: 600;
}

/* Card Container */
.table-card-container {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #eef2f6;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
  padding: 24px 28px;
  margin-bottom: 32px;
}

.table-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 20px;
}

.header-icon-title-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.table-card-title {
  font-size: 1.25rem; /* 20px */
  font-weight: 700;
  color: #1e293b;
  margin: 0;
  line-height: 1.3;
}

.table-card-subtitle {
  font-size: 0.8125rem; /* 13px */
  color: #64748b;
  margin-top: 2px;
}

.table-header-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

/* Search Box */
.table-search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 0 12px;
  height: 40px;
  color: #94a3b8;
  min-width: 220px;
}

.table-search-box input {
  border: none;
  outline: none;
  font-size: 0.875rem; /* 14px */
  color: #1e293b;
  width: 100%;
  font-family: inherit;
  background: transparent;
}

.btn-refresh {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  color: #64748b;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-refresh:hover {
  background: #f8fafc;
  color: #1e293b;
  border-color: #cbd5e1;
}

/* Filter Tabs */
.filter-tabs-row {
  display: flex;
  align-items: center;
  gap: 18px;
  border-bottom: 1px solid #f1f5f9;
  margin-bottom: 20px;
}

.tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 4px 12px;
  background: transparent;
  border: none;
  font-size: 0.875rem; /* 14px */
  font-weight: 500;
  color: #64748b;
  cursor: pointer;
  position: relative;
  transition: color 0.15s ease;
}

.tab-btn:hover {
  color: #1e293b;
}

.tab-btn.active {
  color: #2563eb;
  font-weight: 600;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 3px;
  background: #2563eb;
  border-radius: 2px;
}

.tab-count-pill {
  font-size: 0.75rem; /* 12px */
  font-weight: 600;
  padding: 2px 7px;
  border-radius: 9999px;
  background: #f1f5f9;
  color: #64748b;
}

.tab-count-pill.pill-success {
  background: #dcfce7;
  color: #15803d;
}

/* Loading */
.table-loading-wrap {
  text-align: center;
  padding: 50px 0;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 0.875rem;
}

.spin-icon {
  animation: spin 1s linear infinite;
  color: #2563eb;
}

@keyframes spin {
  100% { transform: rotate(360deg); }
}

/* Table */
.types-table-wrap {
  overflow-x: auto;
  border-radius: 12px;
  border: 1px solid #f1f5f9;
}

.types-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem; /* 14px */
}

.types-table thead th {
  background: #ffffff;
  color: #64748b;
  font-size: 0.8125rem; /* 13px */
  font-weight: 600;
  padding: 14px 16px;
  text-align: left;
  border-bottom: 1px solid #f1f5f9;
  white-space: nowrap;
}

.types-table tbody td {
  padding: 16px 16px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  vertical-align: middle;
}

.types-table tbody tr:hover {
  background: #f8fafc;
}

.sort-order-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 8px;
  background: #f1f5f9;
  color: #475569;
  font-weight: 700;
  font-size: 0.75rem; /* 12px */
}

.code-pill {
  font-family: monospace;
  font-size: 0.8125rem; /* 13px */
  background: #f8fafc;
  color: #0f172a;
  border: 1px solid #e2e8f0;
  padding: 3px 8px;
  border-radius: 6px;
}

.type-title-text {
  font-weight: 600;
  color: #1e293b;
  font-size: 0.9375rem; /* 15px */
}

.type-desc-text {
  color: #64748b;
  font-size: 0.8125rem; /* 13px */
  max-width: 380px;
  line-height: 1.4;
}

.status-soft-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 9999px;
  font-size: 0.75rem; /* 12px */
  font-weight: 600;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.pill-active {
  background: #dcfce7;
  color: #15803d;
}

.pill-active .status-dot {
  background: #15803d;
}

.pill-inactive {
  background: #f1f5f9;
  color: #64748b;
}

.pill-inactive .status-dot {
  background: #94a3b8;
}

/* Actions Cell */
.table-actions-cell {
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
}

.btn-action-icon {
  width: 32px;
  height: 32px;
  border-radius: 7px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-action-icon:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.btn-action-outline {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 6px 12px;
  border-radius: 7px;
  font-size: 0.8125rem; /* 13px */
  font-weight: 500;
  cursor: pointer;
  background: #ffffff;
  transition: all 0.15s ease;
}

.btn-edit {
  border: 1px solid #cbd5e1;
  color: #334155;
}

.btn-edit:hover {
  background: #f8fafc;
  border-color: #94a3b8;
  color: #0f172a;
}

.btn-delete {
  border: 1px solid #fecaca;
  color: #dc2626;
}

.btn-delete:hover {
  background: #fef2f2;
  border-color: #f87171;
}

/* Empty State */
.empty-state-box {
  text-align: center;
  padding: 60px 20px;
  color: #94a3b8;
}

.empty-title {
  font-size: 1rem;
  font-weight: 600;
  color: #334155;
  margin-top: 12px;
}

.empty-subtitle {
  font-size: 0.8125rem;
  color: #64748b;
  margin-top: 4px;
}

/* Modal */
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

.type-modal-box {
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

.modal-icon-primary {
  background: #eff6ff;
  color: #2563eb;
}

.modal-title {
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

.alert-error-box {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fef2f2;
  border: 1px solid #fee2e2;
  color: #dc2626;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 0.8125rem;
  margin-bottom: 16px;
}

.form-group {
  margin-bottom: 16px;
}

.form-row-2col {
  display: grid;
  grid-template-columns: 1.4fr 1fr;
  gap: 12px;
}

.form-label {
  display: block;
  font-size: 0.8125rem; /* 13px */
  font-weight: 600;
  color: #334155;
  margin-bottom: 6px;
}

.text-danger {
  color: #dc2626;
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

.form-textarea {
  width: 100%;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 0.875rem; /* 14px */
  color: #1e293b;
  font-family: inherit;
  outline: none;
  box-sizing: border-box;
  resize: vertical;
  min-height: 70px;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.form-textarea:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.status-toggle-wrap {
  background: #f8fafc;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  margin-top: 8px;
  margin-bottom: 8px;
}

.toggle-checkbox-label {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  cursor: pointer;
}

.toggle-text-block {
  display: flex;
  flex-direction: column;
}

.toggle-title {
  font-size: 0.875rem; /* 14px */
  font-weight: 600;
  color: #1e293b;
}

.toggle-desc {
  font-size: 0.75rem; /* 12px */
  color: #64748b;
  margin-top: 2px;
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
</style>
