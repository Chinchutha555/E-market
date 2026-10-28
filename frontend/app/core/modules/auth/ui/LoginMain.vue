<template>
  <div class="login-wrapper">
    <div class="login-box card">
      <div class="login-header">
        <div class="brand-icon-large">
          <svg width="30" height="30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path>
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <path d="M16 10a4 4 0 0 1-8 0"></path>
          </svg>
        </div>
        <div class="bma-header-pill">
          <span>สำนักงานตลาดกรุงเทพมหานคร</span>
        </div>
        <h2>เข้าสู่ระบบ e-Market</h2>
        <p>ระบบบริหารจัดการแผงเช่าตลาด 12 แห่ง และสัญญาเช่าดิจิทัล</p>
      </div>

      <div v-if="errorMessage" class="alert-error">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
        <span>{{ errorMessage }}</span>
      </div>

      <form novalidate @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label class="form-label" for="username">ชื่อผู้ใช้งานหรืออีเมล (Username / Email) <b class="req-star">*</b></label>
          <AppInput
            id="username"
            v-model="username"
            type="text"
            placeholder="เช่น tenant หรือ admin"
            :error="usernameError"
            @blur="validateUsername"
          />
        </div>

        <div class="form-group">
          <label class="form-label" for="password">รหัสผ่าน (Password) <b class="req-star">*</b></label>
          <AppInput
            id="password"
            v-model="password"
            type="password"
            placeholder="••••••••"
            :error="passwordError"
            @blur="validatePassword"
          />
        </div>

        <button type="submit" class="btn btn-primary btn-submit-login" :disabled="isLoading">
          <span v-if="isLoading" style="display: inline-flex; align-items: center; gap: 8px;">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            <span>กำลังตรวจสอบสิทธิ์...</span>
          </span>
          <span v-else style="display: inline-flex; align-items: center; gap: 8px;">
            <span>เข้าสู่ระบบ (Login)</span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
          </span>
        </button>
      </form>

      <!-- Demo Accounts Quick Fill -->
      <div class="quick-credentials">
        <div class="quick-title">เลือกบัญชีทดสอบระบบด่วน (Quick Demo Login):</div>
        <div class="quick-buttons">
          <button type="button" class="btn-quick" @click="fillCredentials('tenant', 'tenant1234')">
            <div class="quick-badge badge-tenant">ผู้เช่า (Tenant)</div>
            <div class="quick-user-info">
              <strong>tenant</strong> / <span>tenant1234</span>
            </div>
          </button>
          <button type="button" class="btn-quick" @click="fillCredentials('tenant2', 'tenant21234')">
            <div class="quick-badge badge-tenant">ผู้เช่าทดสอบ 2 (Tenant)</div>
            <div class="quick-user-info">
              <strong>tenant2</strong> / <span>tenant21234</span>
            </div>
          </button>
          <button type="button" class="btn-quick" @click="fillCredentials('tenant3', 'tenant31234')">
            <div class="quick-badge badge-tenant">ผู้เช่าทดสอบ 3 (Tenant)</div>
            <div class="quick-user-info">
              <strong>tenant3</strong> / <span>tenant31234</span>
            </div>
          </button>
          <button type="button" class="btn-quick" @click="fillCredentials('admin', 'admin1234')">
            <div class="quick-badge badge-admin">ผู้ดูแลระบบ (Admin)</div>
            <div class="quick-user-info">
              <strong>admin</strong> / <span>admin1234</span>
            </div>
          </button>
            <button type="button" class="btn-quick" @click="fillCredentials('market_director_demo', 'marketdemo1234')">
              <div class="quick-badge badge-staff">ผอ.ตลาด (ทดสอบ)</div>
              <div class="quick-user-info"><strong>market_director_demo</strong> / <span>marketdemo1234</span></div>
            </button>
            <button type="button" class="btn-quick" @click="fillCredentials('bureau_director_demo', 'bureaudemo1234')">
              <div class="quick-badge badge-staff">ผอ.สำนักตลาด (ทดสอบ)</div>
              <div class="quick-user-info"><strong>bureau_director_demo</strong> / <span>bureaudemo1234</span></div>
            </button>
          <section class="staff-demo-accounts" aria-label="บัญชีเจ้าหน้าที่ตลาดทั้ง 12 แห่ง">
            <div class="staff-demo-heading">
              <span>บัญชีเจ้าหน้าที่ตลาด 12 แห่ง</span>
              <small>รหัสผ่านทุกบัญชี: <strong>staff1234</strong></small>
            </div>
            <div class="staff-demo-grid">
              <button
                v-for="staff in staffDemoAccounts"
                :key="staff.username"
                type="button"
                class="btn-quick btn-quick-staff"
                @click="fillCredentials(staff.username, 'staff1234')"
              >
                <div class="quick-badge badge-staff">{{ staff.market }}</div>
                <div class="quick-user-info"><strong>{{ staff.username }}</strong></div>
              </button>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import AppInput from "~/components/AppInput.vue"
import { useAuthStore } from "../store/AuthStore"
import { useAppAlert } from "~/composables/useAppAlert"

const router = useRouter()
const authStore = useAuthStore()
const { showSuccess, handleConfirm } = useAppAlert()

const username = ref("tenant")
const password = ref("tenant1234")
const usernameError = ref("")
const passwordError = ref("")
const isLoading = ref(false)
const errorMessage = ref("")
const staffDemoAccounts = [
  { username: "staff01", market: "ตลาดนัดจตุจักร" },
  { username: "staff02", market: "ตลาดประชานิเวศน์ 1" },
  { username: "staff03", market: "ตลาดมีนบุรี" },
  { username: "staff04", market: "ตลาดสิงหา" },
  { username: "staff05", market: "ตลาดเทวราช" },
  { username: "staff06", market: "ตลาดราษฎร์บูรณะ" },
  { username: "staff07", market: "ตลาดบางกะปิ" },
  { username: "staff08", market: "ตลาดวงเวียนเล็ก" },
  { username: "staff09", market: "ตลาดหนองจอก" },
  { username: "staff10", market: "ตลาดรัชดาภิเษก" },
  { username: "staff11", market: "ตลาดธนบุรี สนามหลวง 2" },
  { username: "staff12", market: "ตลาดบางแคภิรมย์" },
]

const validateUsername = () => {
  usernameError.value = username.value.trim() ? "" : "กรุณาระบุชื่อผู้ใช้งานหรืออีเมล"
}

const validatePassword = () => {
  passwordError.value = password.value.trim() ? "" : "กรุณาระบุรหัสผ่าน"
}

const fillCredentials = (u: string, p: string) => {
  username.value = u
  password.value = p
  usernameError.value = ""
  passwordError.value = ""
}

const handleLogin = async () => {
  validateUsername()
  validatePassword()
  if (usernameError.value || passwordError.value) return
  isLoading.value = true
  errorMessage.value = ""

  try {
    await authStore.login({
      username: username.value,
      password: password.value,
    })

    // แสดงสถานะสำเร็จสั้น ๆ แล้วเปลี่ยนหน้าเอง โดยไม่ต้องกดยืนยันซ้ำ
    const successAlert = showSuccess("เข้าสู่ระบบสำเร็จ ยินดีต้อนรับสู่ e-Market")
    const autoCloseTimer = window.setTimeout(handleConfirm, 900)
    await successAlert
    window.clearTimeout(autoCloseTimer)
    await new Promise<void>((resolve) => window.setTimeout(resolve, 220))
    await router.push(authStore.role === "staff" ? "/staff" : "/")
  } catch (err: any) {
    errorMessage.value = err.message || "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง"
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-wrapper {
  min-height: 100vh;
}
.req-star {
  color: #ef4444;
  font-weight: 700;
}
.login-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  background: radial-gradient(circle at 10% 20%, rgba(37, 99, 235, 0.08) 0%, transparent 40%),
              radial-gradient(circle at 90% 80%, rgba(16, 185, 129, 0.06) 0%, transparent 40%),
              #f8fafc;
}

.login-box {
  width: 100%;
  max-width: 440px;
  padding: 36px 32px;
  border-radius: var(--radius-2xl);
  border: 1px solid #e2e8f0;
  box-shadow: 0 20px 40px -10px rgba(15, 23, 42, 0.1), 0 4px 14px rgba(15, 23, 42, 0.04);
  background: #ffffff;
}

.login-header {
  text-align: center;
  margin-bottom: 24px;
}

.brand-icon-large {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
  box-shadow: 0 8px 20px rgba(37, 99, 235, 0.3);
}

.bma-header-pill {
  display: inline-block;
  background: #eff6ff;
  border: 1px solid #dbeafe;
  color: #1d4ed8;
  padding: 3px 12px;
  border-radius: 20px;
  font-size: 0.74rem;
  font-weight: 700;
  margin-bottom: 8px;
}

.login-header h2 {
  font-size: 1.55rem;
  font-weight: 800;
  margin-bottom: 6px;
  color: #0f172a;
  letter-spacing: -0.02em;
}

.login-header p {
  color: #64748b;
  font-size: 0.86rem;
  line-height: 1.5;
}

.alert-error {
  background: #fff1f2;
  border: 1px solid #fecdd3;
  color: #e11d48;
  padding: 10px 14px;
  border-radius: 10px;
  margin-bottom: 18px;
  font-size: 0.84rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-with-icon {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 12px;
  display: flex;
  align-items: center;
  pointer-events: none;
}

.input-with-icon .form-input {
  padding-left: 38px;
}

.btn-submit-login {
  width: 100%;
  padding: 12px;
  font-size: 0.95rem;
  font-weight: 700;
  border-radius: 12px;
  margin-top: 6px;
}

.quick-credentials {
  margin-top: 26px;
  padding-top: 20px;
  border-top: 1px solid #f1f5f9;
}

.quick-title {
  font-size: 0.78rem;
  color: #64748b;
  margin-bottom: 10px;
  text-align: center;
  font-weight: 700;
}

.quick-buttons {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.btn-quick {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 8px 12px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-quick:hover {
  background: #eff6ff;
  border-color: #bfdbfe;
  transform: translateY(-1px);
}

.quick-badge {
  font-size: 0.8125rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 6px;
}

.badge-tenant {
  background: #ecfdf5;
  color: #047857;
  border: 1px solid #a7f3d0;
}

.badge-admin {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.badge-staff {
  background: #faf5ff;
  color: #7e22ce;
  border: 1px solid #e9d5ff;
}

.quick-user-info {
  font-size: 0.78rem;
  color: #475569;
}

.quick-user-info strong {
  color: #0f172a;
}

.staff-demo-accounts {
  padding: 12px;
  border: 1px solid #eadcff;
  border-radius: 12px;
  background: #fcfaff;
}

.staff-demo-heading {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 10px;
  color: #6b21a8;
  font-size: 0.8125rem;
  font-weight: 700;
}

.staff-demo-heading small {
  color: #7c3aed;
  font-size: 0.75rem;
  font-weight: 500;
}

.staff-demo-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.btn-quick-staff {
  min-width: 0;
  padding: 7px 9px;
  gap: 6px;
}

.btn-quick-staff .quick-badge {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 560px) {
  .staff-demo-heading { align-items: flex-start; flex-direction: column; }
  .staff-demo-grid { grid-template-columns: 1fr; }
}
</style>
