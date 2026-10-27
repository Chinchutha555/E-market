<template>
  <div>
    <AppHeader />

    <div class="content-body">
      <!-- Zones Management Card -->
      <div class="card" style="margin-bottom: 30px;">
        <div class="card-header">
          <div class="card-title" style="display: flex; align-items: center; gap: 8px;">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--primary);"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
            <span>โซนพื้นที่ตลาด (Market Zones)</span>
          </div>
          <div style="display: flex; align-items: center; gap: 10px;">
            <button class="btn btn-primary btn-sm" @click="openCreateZoneModal" style="display: inline-flex; align-items: center; gap: 6px;">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
              <span>เพิ่มโซนตลาดใหม่</span>
            </button>
            <button class="btn btn-secondary btn-sm" style="display: inline-flex; align-items: center; gap: 5px;" @click="loadZones">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path></svg>
              <span>รีเฟรช</span>
            </button>
          </div>
        </div>

        <div class="table-container">
          <table class="custom-table">
            <thead>
              <tr>
                <th>รหัสโซน</th>
                <th>ชื่อโซนตลาด</th>
                <th>สีประจำโซน</th>
                <th>จำนวนแผงทั้งหมด</th>
                <th>คำอธิบาย</th>
                <th style="text-align: right;">การจัดการ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="z in zones" :key="z.id">
                <td>
                  <span class="badge" :style="{ backgroundColor: z.color_code + '20', color: z.color_code, borderColor: z.color_code + '40', fontSize: '0.85rem' }">
                    Zone {{ z.code }}
                  </span>
                </td>
                <td><strong>{{ z.name }}</strong></td>
                <td>
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <span :style="{ backgroundColor: z.color_code, width: '16px', height: '16px', borderRadius: '4px', display: 'inline-block' }"></span>
                    <span style="font-family: monospace; font-size: 0.8rem;">{{ z.color_code }}</span>
                  </div>
                </td>
                <td><strong>{{ z.stalls_count }} แผง</strong></td>
                <td style="color: var(--text-muted); font-size: 0.82rem;">{{ z.description || '-' }}</td>
                <td style="text-align: right;">
                  <div style="display: inline-flex; gap: 6px;">
                    <button class="btn btn-sm btn-secondary" @click="openEditZoneModal(z)">แก้ไข</button>
                    <button class="btn btn-sm btn-danger" @click="deleteZone(z)" :disabled="z.stalls_count > 0">ลบ</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Application Types Management Card -->
      <div class="card" style="margin-bottom: 30px;">
        <div class="card-header" style="display: flex; justify-content: space-between; align-items: center;">
          <div class="card-title" style="display: flex; align-items: center; gap: 8px;">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--primary);">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <polyline points="10 9 9 9 8 9"></polyline>
            </svg>
            <span>ประเภทคำร้อง / เรื่องที่ขอ (Application Types)</span>
          </div>
          <NuxtLink
            to="/applications/types"
            class="btn btn-primary btn-sm"
            style="display: inline-flex; align-items: center; gap: 6px; text-decoration: none;"
          >
            <span>จัดการประเภทคำร้อง</span>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
          </NuxtLink>
        </div>
        <div style="padding: 16px 20px; color: #475569; font-size: 0.875rem; line-height: 1.6; background: #f8fafc; border-radius: 8px; margin: 0 20px 20px; border: 1px solid var(--border-subtle);">
          ระบบจัดการหัวข้อ/เรื่องที่ขอในการยื่นคำร้องทั่วไป เช่น ขอเปลี่ยนชื่อผู้ค้าประจำแผง, ขอติดตั้งมิเตอร์ไฟฟ้า, ขอยกเลิกการเช่า ฯลฯ สามารถกำหนดรหัส, ลำดับการแสดงผล และเปิด-ปิดการใช้งานเพื่อให้หน้าบ้านดึงไปแสดงผลโดยอัตโนมัติ
        </div>
      </div>

      <!-- System Architecture & API Info -->
      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px;">
        <!-- REST API Info -->
        <div class="card">
          <div class="card-header">
            <div class="card-title" style="display: flex; align-items: center; gap: 8px;">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--primary);"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon></svg>
              <span>FastAPI REST API & Security Standards</span>
            </div>
          </div>

          <div style="display: flex; flex-direction: column; gap: 14px; font-size: 0.88rem;">
            <div style="background: #f8fafc; padding: 14px; border-radius: 8px; border: 1px solid var(--border-subtle);">
              <div style="color: var(--primary); font-weight: 700; margin-bottom: 4px; display: flex; align-items: center; gap: 6px;">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
                <span>JWT Authentication:</span>
              </div>
              <div style="color: #334155;">เข้าถึงทุก API ผ่าน Authorization Header: <code>Bearer &lt;token&gt;</code> มาตรฐาน RFC 7519</div>
            </div>

            <div style="background: #f8fafc; padding: 14px; border-radius: 8px; border: 1px solid var(--border-subtle);">
              <div style="color: #2563eb; font-weight: 700; margin-bottom: 4px; display: flex; align-items: center; gap: 6px;">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path></svg>
                <span>Interactive Documentation:</span>
              </div>
              <div style="color: #334155;">
                Swagger UI: <a href="http://localhost:8000/docs" target="_blank" style="color: var(--primary); text-decoration: underline; font-weight: 600;">http://localhost:8000/docs</a><br>
                ReDoc: <a href="http://localhost:8000/redoc" target="_blank" style="color: #7c3aed; text-decoration: underline; font-weight: 600;">http://localhost:8000/redoc</a>
              </div>
            </div>

            <div style="background: #f8fafc; padding: 14px; border-radius: 8px; border: 1px solid var(--border-subtle);">
              <div style="color: #d97706; font-weight: 700; margin-bottom: 4px; display: flex; align-items: center; gap: 6px;">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
                <span>Endpoints Security:</span>
              </div>
              <div style="color: #334155;">CORS Allowed Origins, Pydantic Schema Validation, Direct Bcrypt Hashing, Password Expiry Policy.</div>
            </div>
          </div>
        </div>

        <!-- Tech Stack Info -->
        <div class="card">
          <div class="card-header">
            <div class="card-title" style="display: flex; align-items: center; gap: 8px;">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--primary);"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>
              <span>Technology Stack</span>
            </div>
          </div>

          <div style="display: flex; flex-direction: column; gap: 12px; font-size: 0.88rem;">
            <div style="display: flex; justify-content: space-between; padding: 12px; background: #f8fafc; border: 1px solid var(--border-subtle); border-radius: 8px;">
              <span style="color: var(--text-muted);">Frontend Framework:</span>
              <strong style="color: #059669;">Nuxt.js (Nuxt 3.12 + Vue 3)</strong>
            </div>
            <div style="display: flex; justify-content: space-between; padding: 12px; background: #f8fafc; border: 1px solid var(--border-subtle); border-radius: 8px;">
              <span style="color: var(--text-muted);">Backend Framework:</span>
              <strong style="color: #2563eb;">Python FastAPI (REST API Standard)</strong>
            </div>
            <div style="display: flex; justify-content: space-between; padding: 12px; background: #f8fafc; border: 1px solid var(--border-subtle); border-radius: 8px;">
              <span style="color: var(--text-muted);">Database & ORM:</span>
              <strong style="color: #7c3aed;">SQLAlchemy 2.0 + SQLite</strong>
            </div>
            <div style="display: flex; justify-content: space-between; padding: 12px; background: #f8fafc; border: 1px solid var(--border-subtle); border-radius: 8px;">
              <span style="color: var(--text-muted);">Authentication:</span>
              <strong style="color: #d97706;">OAuth2 Bearer JWT Token (HS256)</strong>
            </div>
            <div style="display: flex; justify-content: space-between; padding: 12px; background: #f8fafc; border: 1px solid var(--border-subtle); border-radius: 8px;">
              <span style="color: var(--text-muted);">UI Design System:</span>
              <strong style="color: #4361ee;">Modern Royal Blue Light Design</strong>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal: Create / Edit Zone -->
    <div v-if="showZoneModal" class="modal-overlay" @click.self="showZoneModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3 style="font-size: 1.15rem;">{{ isEditing ? 'แก้ไขข้อมูลโซน' : 'เพิ่มโซนตลาดใหม่' }}</h3>
          <button class="btn-logout" @click="showZoneModal = false" aria-label="ปิด">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>

        <form @submit.prevent="saveZone">
          <div class="modal-body">
            <div class="form-row">
              <div class="form-group">
                <label class="form-label" for="z_code">รหัสโซน (Zone Code) *</label>
                <input id="z_code" v-model="formZone.code" type="text" required placeholder="เช่น E, F" class="form-input" maxlength="5" />
              </div>

              <div class="form-group">
                <label class="form-label" for="z_color">สีประจำโซน *</label>
                <div style="display: flex; gap: 8px;">
                  <input id="z_color" v-model="formZone.color_code" type="color" class="form-input" style="width: 50px; padding: 2px;" />
                  <input v-model="formZone.color_code" type="text" class="form-input" style="flex: 1;" />
                </div>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="z_name">ชื่อโซนตลาด *</label>
              <input id="z_name" v-model="formZone.name" type="text" required placeholder="เช่น โซน E: สินค้าโอทอปและของฝาก" class="form-input" />
            </div>

            <div class="form-group">
              <label class="form-label" for="z_desc">คำอธิบายโซน</label>
              <textarea id="z_desc" v-model="formZone.description" rows="2" class="form-textarea" placeholder="รายละเอียดสินค้าและลักษณะพื้นที่"></textarea>
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="showZoneModal = false">ยกเลิก</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              <span v-if="saving">กำลังบันทึก...</span>
              <span v-else>บันทึก</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'

definePageMeta({
  middleware: () => {
    if (!import.meta.client) return
    const auth = useAuthStore()
    auth.hydrate()
    if (auth.role !== 'admin') return navigateTo(auth.role === 'staff' ? '/staff' : '/', { replace: true })
  }
})
const { get, post, put, del } = useApi()
const { showSuccess, showError, showConfirm } = useAppAlert()

const zones = ref<any[]>([])
const loading = ref(true)
const saving = ref(false)

const showZoneModal = ref(false)
const isEditing = ref(false)
const editZoneId = ref<number | null>(null)
const formZone = ref({
  code: '',
  name: '',
  color_code: '#3B82F6',
  description: ''
})

const loadZones = async () => {
  loading.value = true
  try {
    const res: any = await get('/zones')
    zones.value = Array.isArray(res) ? res : (Array.isArray(res?.data) ? res.data : [])
  } catch (err) {
    console.error('Failed to load zones:', err)
  } finally {
    loading.value = false
  }
}

const openCreateZoneModal = () => {
  isEditing.value = false
  editZoneId.value = null
  formZone.value = {
    code: '',
    name: '',
    color_code: '#10B981',
    description: ''
  }
  showZoneModal.value = true
}

const openEditZoneModal = (z: any) => {
  isEditing.value = true
  editZoneId.value = z.id
  formZone.value = {
    code: z.code,
    name: z.name,
    color_code: z.color_code,
    description: z.description || ''
  }
  showZoneModal.value = true
}

const saveZone = async () => {
  saving.value = true
  try {
    if (isEditing.value && editZoneId.value) {
      await put(`/zones/${editZoneId.value}`, formZone.value)
    } else {
      await post('/zones', formZone.value)
    }
    showZoneModal.value = false
    await loadZones()
    await showSuccess('บันทึกข้อมูลโซนเรียบร้อยแล้ว', 'สำเร็จ')
  } catch (err: any) {
    await showError(err?.data?.detail || 'เกิดข้อผิดพลาดในการบันทึกข้อมูล', 'ข้อผิดพลาด')
  } finally {
    saving.value = false
  }
}

const deleteZone = async (z: any) => {
  const confirmed = await showConfirm(
    `คุณต้องการลบโซน ${z.code} - ${z.name} หรือไม่?`,
    {
      title: 'ยืนยันลบโซน',
      type: 'error',
      icon: 'trash',
      confirmText: 'ลบโซน',
      cancelText: 'ยกเลิก'
    }
  )
  if (confirmed) {
    try {
      await del(`/zones/${z.id}`)
      await loadZones()
      await showSuccess('ลบโซนเรียบร้อยแล้ว', 'สำเร็จ')
    } catch (err: any) {
      await showError(err?.data?.detail || 'ไม่สามารถลบโซนได้', 'ข้อผิดพลาด')
    }
  }
}

onMounted(() => {
  loadZones()
})
</script>
