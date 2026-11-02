<template>
  <Teleport to="body">
    <Transition name="app-modal-fade">
      <div v-if="modelValue" class="app-modal-backdrop" @click.self="close">
        <section class="app-modal" :class="`app-modal--${size}`" role="dialog" aria-modal="true" :aria-labelledby="titleId">
          <header class="app-modal-header">
            <div class="app-modal-heading">
              <div class="app-modal-icon" aria-hidden="true">
                <slot name="icon">
                  <svg v-if="icon === 'map'" viewBox="0 0 24 24"><polygon points="1 6 8 2 16 6 23 2 23 18 16 22 8 18 1 22 1 6"/><line x1="8" y1="2" x2="8" y2="18"/><line x1="16" y1="6" x2="16" y2="22"/></svg>
                  <svg v-else-if="icon === 'pin'" viewBox="0 0 24 24"><path d="M21 10c0 7-9 12-9 12S3 17 3 10a9 9 0 1 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
                  <svg v-else-if="icon === 'users'" viewBox="0 0 24 24"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
                  <svg v-else-if="icon === 'tool'" viewBox="0 0 24 24"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94z"/></svg>
                  <svg v-else-if="icon === 'payment'" viewBox="0 0 24 24"><rect x="2" y="5" width="20" height="14" rx="2"/><line x1="2" y1="10" x2="22" y2="10"/></svg>
                  <svg v-else-if="icon === 'settings'" viewBox="0 0 24 24"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.7 1.7 0 0 0 .34 1.88l.06.06-2.83 2.83-.06-.06a1.7 1.7 0 0 0-1.88-.34 1.7 1.7 0 0 0-1.03 1.55V21h-4v-.08A1.7 1.7 0 0 0 8.97 19.4a1.7 1.7 0 0 0-1.88.34l-.06.06-2.83-2.83.06-.06A1.7 1.7 0 0 0 4.6 15 1.7 1.7 0 0 0 3.08 14H3v-4h.08A1.7 1.7 0 0 0 4.6 8.97a1.7 1.7 0 0 0-.34-1.88L4.2 7.03 7.03 4.2l.06.06A1.7 1.7 0 0 0 8.97 4.6 1.7 1.7 0 0 0 10 3.08V3h4v.08a1.7 1.7 0 0 0 1.03 1.52 1.7 1.7 0 0 0 1.88-.34l.06-.06 2.83 2.83-.06.06a1.7 1.7 0 0 0-.34 1.88A1.7 1.7 0 0 0 20.92 10H21v4h-.08A1.7 1.7 0 0 0 19.4 15z"/></svg>
                  <svg v-else viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
                </slot>
              </div>
              <div class="app-modal-heading-copy">
                <h2 :id="titleId">{{ title }}</h2>
                <p v-if="description">{{ description }}</p>
              </div>
            </div>
            <AppButton variant="ghost" size="sm" class="app-modal-close" aria-label="ปิด" @click="close">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 6 12 12M18 6 6 18" /></svg>
            </AppButton>
          </header>

          <div class="app-modal-body"><slot /></div>
          <footer v-if="$slots.footer" class="app-modal-footer"><slot name="footer" /></footer>
        </section>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { toRef } from 'vue'
import AppButton from './AppButton.vue'
import { useBodyScrollLock } from '~/composables/useBodyScrollLock'

const props = withDefaults(defineProps<{
  modelValue: boolean
  title: string
  description?: string
  icon?: 'document' | 'map' | 'pin' | 'users' | 'tool' | 'payment' | 'settings'
  // md = ฟอร์มทั่วไป (480px), xl = ตารางข้อมูลกว้าง (1100px)
  size?: 'md' | 'xl'
}>(), { size: 'md' })

const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>()
const titleId = `app-modal-title-${Math.random().toString(36).slice(2, 9)}`
const close = () => emit('update:modelValue', false)

useBodyScrollLock(toRef(props, 'modelValue'))
</script>

<style scoped>
.app-modal-backdrop{
  position:fixed;
  inset:0;
  z-index:2500;
  display:flex;
  align-items:center;
  justify-content:center;
  padding:16px;
  background:rgba(15,23,42,.55);
  backdrop-filter:blur(4px);
  overflow-y:auto;
  -webkit-overflow-scrolling:touch;
}
.app-modal{
  width:min(480px,100%);
  max-height:min(90vh, 90dvh);
  display:flex;
  flex-direction:column;
  overflow:hidden;
  border:1px solid var(--bd-e2e8f0);
  border-radius:20px;
  background:var(--bg-ffffff);
  box-shadow:0 24px 48px rgba(15,23,42,.22);
  margin:auto;
}
.app-modal--xl{
  width:min(1100px,100%);
}
.app-modal-header{
  flex:0 0 auto;
  display:flex;
  align-items:center;
  justify-content:space-between;
  gap:16px;
  min-height:112px;
  padding:24px 28px;
  border-bottom:1px solid var(--bd-e2e8f0);
}
.app-modal-heading{display:flex;align-items:center;gap:16px;min-width:0}
.app-modal-icon{width:48px;height:48px;display:grid;place-items:center;flex:0 0 48px;border-radius:14px;border:1px solid var(--bd-dbeafe);background:var(--bg-eff6ff);color:var(--tx-2563eb)}
.app-modal-icon svg{width:24px;height:24px;fill:none;stroke:currentColor;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}
.app-modal-heading-copy{min-width:0}
.app-modal-header h2{margin:0;color:var(--tx-0f172a);font-size:var(--font-size-xl);line-height:1.35;font-weight:700}
.app-modal-header p{margin:4px 0 0;color:var(--tx-64748b);font-size:var(--font-size-base);line-height:1.5}
.app-modal-close{flex:0 0 auto!important;padding:0!important;width:38px!important;height:38px!important;color:var(--tx-94a3b8)!important}.app-modal-close svg{width:20px;height:20px}
.app-modal-body{
  flex:1 1 auto;
  min-height:0;
  overflow-y:auto;
  -webkit-overflow-scrolling:touch;
  padding:28px;
}
.app-modal-footer{
  flex:0 0 auto;
  display:flex;
  justify-content:flex-end;
  gap:10px;
  padding:16px 28px;
  border-top:1px solid var(--bd-e2e8f0);
  background:var(--bg-f8fafc);
}
.app-modal-fade-enter-active,.app-modal-fade-leave-active{transition:opacity .16s ease}.app-modal-fade-enter-active .app-modal,.app-modal-fade-leave-active .app-modal{transition:transform .16s ease}.app-modal-fade-enter-from,.app-modal-fade-leave-to{opacity:0}.app-modal-fade-enter-from .app-modal,.app-modal-fade-leave-to .app-modal{transform:translateY(6px) scale(.98)}
@media(max-width:640px){.app-modal-header{min-height:92px;padding:18px 20px}.app-modal-icon{width:44px;height:44px;flex-basis:44px}.app-modal-header h2{font-size:var(--font-size-lg)}.app-modal-header p{font-size:var(--font-size-sm)}.app-modal-body{padding:22px 20px}.app-modal-footer{padding:14px 20px;flex-wrap:wrap}}
</style>
