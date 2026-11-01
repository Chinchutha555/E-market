<template>
  <Teleport to="body">
    <transition name="modal-bounce">
      <div
        v-if="isOpen"
        class="app-alert-backdrop"
        @click.self="handleCancel"
        role="dialog"
        aria-modal="true"
      >
        <div class="app-alert-card">
          <!-- Icon Header -->
          <div class="app-alert-icon-wrapper" :class="`icon-${alertType === 'danger' ? 'error' : alertType}`">
            <!-- Trash Icon (Destructive / Delete / Clear) -->
            <svg
              v-if="alertIcon === 'trash'"
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              <line x1="10" y1="11" x2="10" y2="17"></line>
              <line x1="14" y1="11" x2="14" y2="17"></line>
            </svg>

            <!-- Warning Icon -->
            <svg
              v-else-if="alertType === 'warning'"
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.3"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
              <line x1="12" y1="9" x2="12" y2="13"></line>
              <line x1="12" y1="17" x2="12.01" y2="17"></line>
            </svg>

            <!-- Error / Danger Icon -->
            <svg
              v-else-if="alertType === 'error' || alertType === 'danger'"
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.3"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="15" y1="9" x2="9" y2="15"></line>
              <line x1="9" y1="9" x2="15" y2="15"></line>
            </svg>

            <!-- Success Icon -->
            <svg
              v-else-if="alertType === 'success'"
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>

            <!-- Info Icon -->
            <svg
              v-else
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.3"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="16" x2="12" y2="12"></line>
              <line x1="12" y1="8" x2="12.01" y2="8"></line>
            </svg>
          </div>

          <!-- Title & Message -->
          <div class="app-alert-content">
            <h3 class="app-alert-title">{{ alertTitle }}</h3>
            <p class="app-alert-message">{{ alertMessage }}</p>
          </div>

          <!-- Actions -->
          <div class="app-alert-actions">
            <button
              v-if="isConfirm"
              type="button"
              class="btn-alert-cancel"
              @click="handleCancel"
            >
              {{ cancelText }}
            </button>
            <button
              type="button"
              class="btn-alert-confirm"
              :class="alertType === 'danger' ? 'btn-error' : `btn-${alertType}`"
              @click="handleConfirm"
              autofocus
            >
              {{ confirmText }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useAppAlert } from '~/composables/useAppAlert'
import { useBodyScrollLock } from '~/composables/useBodyScrollLock'

const {
  isOpen,
  alertTitle,
  alertMessage,
  alertType,
  alertIcon,
  confirmText,
  cancelText,
  isConfirm,
  handleConfirm,
  handleCancel
} = useAppAlert()

useBodyScrollLock(isOpen)

const handleKeydown = (e: KeyboardEvent) => {
  if (!isOpen.value) return
  if (e.key === 'Escape') {
    handleCancel()
  } else if (e.key === 'Enter') {
    handleConfirm()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)

  // Globally intercept window.alert so third-party or residual calls use this modal
  if (process.client) {
    (window as any).nativeAlert = window.alert
    window.alert = (msg?: any) => {
      const { showAlert } = useAppAlert()
      showAlert(String(msg || ''), { title: 'แจ้งเตือน', type: 'warning' })
    }
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.app-alert-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  z-index: 99999;
}

.app-alert-card {
  background: var(--bg-ffffff);
  border-radius: 18px;
  width: 100%;
  max-width: 420px;
  padding: 28px 24px 22px;
  text-align: center;
  box-shadow: 0 24px 48px -12px rgba(15, 23, 42, 0.28), 0 0 0 1px rgba(226, 232, 240, 0.8);
  position: relative;
  overflow: hidden;
}

.app-alert-icon-wrapper {
  width: 58px;
  height: 58px;
  border-radius: 16px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s ease;
}

.app-alert-icon-wrapper.icon-warning {
  background: var(--bg-fef3c7);
  color: var(--tx-d97706);
  border: 1px solid var(--bd-fde68a);
  box-shadow: 0 4px 12px rgba(217, 119, 6, 0.15);
}

.app-alert-icon-wrapper.icon-error {
  background: var(--bg-fee2e2);
  color: var(--tx-dc2626);
  border: 1px solid var(--bd-fecaca);
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.15);
}

.app-alert-icon-wrapper.icon-success {
  background: var(--bg-ecfdf5);
  color: var(--tx-10b981);
  border: 1px solid var(--bd-a7f3d0);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.15);
}

.app-alert-icon-wrapper.icon-info {
  background: var(--bg-eff6ff);
  color: var(--tx-2563eb);
  border: 1px solid var(--bd-bfdbfe);
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.15);
}

.app-alert-title {
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--tx-0f172a);
  margin: 0 0 8px 0;
  line-height: 1.35;
}

.app-alert-message {
  font-size: 0.92rem;
  color: var(--tx-475569);
  line-height: 1.55;
  margin: 0 0 24px 0;
  word-break: break-word;
}

.app-alert-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.btn-alert-cancel {
  flex: 1;
  padding: 11px 20px;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 600;
  background: var(--bg-f1f5f9);
  color: var(--tx-475569);
  border: 1px solid var(--bd-cbd5e1);
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-alert-cancel:hover {
  background: var(--bg-e2e8f0);
  color: var(--tx-1e293b);
}

.btn-alert-confirm {
  flex: 1;
  padding: 11px 24px;
  border-radius: 10px;
  font-size: 0.92rem;
  font-weight: 700;
  border: none;
  cursor: pointer;
  transition: all 0.15s ease;
  color: #ffffff;
}

.btn-alert-confirm.btn-warning {
  background: linear-gradient(135deg, #d97706 0%, #b45309 100%);
  box-shadow: 0 4px 14px rgba(217, 119, 6, 0.35);
}
.btn-alert-confirm.btn-warning:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 18px rgba(217, 119, 6, 0.45);
}

.btn-alert-confirm.btn-error {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  box-shadow: 0 4px 14px rgba(239, 68, 68, 0.35);
}
.btn-alert-confirm.btn-error:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 18px rgba(239, 68, 68, 0.45);
}

.btn-alert-confirm.btn-success {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  box-shadow: 0 4px 14px rgba(16, 185, 129, 0.35);
}
.btn-alert-confirm.btn-success:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 18px rgba(16, 185, 129, 0.45);
}

.btn-alert-confirm.btn-info {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.35);
}
.btn-alert-confirm.btn-info:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 18px rgba(37, 99, 235, 0.45);
}

/* Animations */
.modal-bounce-enter-active,
.modal-bounce-leave-active {
  transition: opacity 0.22s ease;
}

.modal-bounce-enter-from,
.modal-bounce-leave-to {
  opacity: 0;
}

.modal-bounce-enter-active .app-alert-card {
  animation: alertPopIn 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.modal-bounce-leave-active .app-alert-card {
  animation: alertPopOut 0.18s ease-in;
}

@keyframes alertPopIn {
  from {
    transform: scale(0.92) translateY(10px);
    opacity: 0;
  }
  to {
    transform: scale(1) translateY(0);
    opacity: 1;
  }
}

@keyframes alertPopOut {
  from {
    transform: scale(1) translateY(0);
    opacity: 1;
  }
  to {
    transform: scale(0.92) translateY(10px);
    opacity: 0;
  }
}
</style>
