<template>
  <header class="modal-header app-standard-modal-header" :class="`variant-${variant}`">
    <div class="app-standard-modal-heading">
      <div class="app-standard-modal-icon" :class="`tone-${tone}`" aria-hidden="true">
        <slot name="icon">
          <svg v-if="icon === 'edit'" viewBox="0 0 24 24"><path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"/></svg>
          <svg v-else-if="icon === 'calendar'" viewBox="0 0 24 24"><rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>
          <svg v-else-if="icon === 'pin'" viewBox="0 0 24 24"><path d="M21 10c0 7-9 12-9 12S3 17 3 10a9 9 0 1 1 18 0Z"/><circle cx="12" cy="10" r="3"/></svg>
          <svg v-else-if="icon === 'clock'" viewBox="0 0 24 24"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg>
          <svg v-else-if="icon === 'grid'" viewBox="0 0 24 24"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
          <svg v-else-if="icon === 'settings'" viewBox="0 0 24 24"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.7 1.7 0 0 0 .34 1.88l.06.06-2.83 2.83-.06-.06a1.7 1.7 0 0 0-2.91 1.21V21h-4v-.08a1.7 1.7 0 0 0-2.91-1.21l-.06.06-2.83-2.83.06-.06A1.7 1.7 0 0 0 3.08 14H3v-4h.08a1.7 1.7 0 0 0 1.18-2.91l-.06-.06L7.03 4.2l.06.06A1.7 1.7 0 0 0 10 3.08V3h4v.08a1.7 1.7 0 0 0 2.91 1.18l.06-.06 2.83 2.83-.06.06A1.7 1.7 0 0 0 20.92 10H21v4h-.08A1.7 1.7 0 0 0 19.4 15Z"/></svg>
          <svg v-else-if="icon === 'warning'" viewBox="0 0 24 24"><path d="m21.7 18-8-14a2 2 0 0 0-3.4 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.7-3Z"/><path d="M12 9v4M12 17h.01"/></svg>
          <svg v-else viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8Z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg>
        </slot>
      </div>
      <div class="app-standard-modal-copy">
        <div v-if="$slots.eyebrow" class="app-standard-modal-eyebrow"><slot name="eyebrow" /></div>
        <h2>{{ title }}</h2>
        <p v-if="subtitle">{{ subtitle }}</p>
        <div v-if="$slots.meta" class="app-standard-modal-meta"><slot name="meta" /></div>
      </div>
    </div>
    <div class="app-standard-modal-actions">
      <slot name="actions" />
      <AppButton v-if="closable" variant="ghost" size="sm" class="app-standard-modal-close" aria-label="ปิด" @click="emit('close')">
        <template #prefix><svg viewBox="0 0 24 24"><path d="M18 6 6 18M6 6l12 12"/></svg></template>
      </AppButton>
    </div>
  </header>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  title: string
  subtitle?: string
  icon?: 'document' | 'edit' | 'calendar' | 'pin' | 'clock' | 'grid' | 'settings' | 'warning'
  tone?: 'primary' | 'danger' | 'warning' | 'success'
  variant?: 'standard' | 'booking-pass'
  closable?: boolean
}>(), { subtitle: '', icon: 'document', tone: 'primary', variant: 'standard', closable: true })

const emit = defineEmits<{ close: [] }>()
</script>

<style scoped>
.app-standard-modal-header{box-sizing:border-box;display:flex;align-items:center;justify-content:space-between;gap:20px;min-height:112px;padding:24px 28px;border-bottom:1px solid var(--bd-e2e8f0);background:var(--bg-ffffff)}
.app-standard-modal-heading{display:flex;align-items:center;gap:16px;min-width:0}
.app-standard-modal-icon{display:grid;place-items:center;width:48px;height:48px;flex:0 0 48px;border:1px solid var(--bd-dbeafe);border-radius:14px;background:var(--bg-eff6ff);color:var(--tx-2563eb)}
.app-standard-modal-icon.tone-danger{border-color:var(--bd-fecaca);background:var(--bg-fee2e2);color:var(--tx-dc2626)}.app-standard-modal-icon.tone-warning{border-color:var(--bd-fde68a);background:var(--bg-fef3c7);color:var(--tx-d97706)}.app-standard-modal-icon.tone-success{border-color:var(--bd-bbf7d0);background:var(--bg-dcfce7);color:var(--tx-16a34a)}
.app-standard-modal-icon svg{width:24px;height:24px;fill:none;stroke:currentColor;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}
.app-standard-modal-copy{min-width:0}.app-standard-modal-copy h2{margin:0;color:var(--tx-0f172a);font-size:1.25rem;line-height:1.35;font-weight:700}.app-standard-modal-copy p{margin:4px 0 0;color:var(--tx-64748b);font-size:.9375rem;line-height:1.5;font-weight:500}
.app-standard-modal-actions{display:flex;align-items:center;gap:10px;flex:0 0 auto}.app-standard-modal-close{width:38px!important;height:38px!important;min-height:38px!important;max-height:38px!important;padding:0!important;color:var(--tx-94a3b8)!important}.app-standard-modal-close :deep(.app-btn-prefix){margin:0}.app-standard-modal-close svg{width:20px;height:20px;fill:none;stroke:currentColor;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}
.app-standard-modal-header.variant-booking-pass{position:relative;min-height:0;padding:24px 28px;border-bottom:0;background:linear-gradient(135deg,#1e3a8a 0%,#2563eb 50%,#3b82f6 100%);color:#fff;text-align:left}
.variant-booking-pass .app-standard-modal-heading{width:100%;align-items:flex-start}.variant-booking-pass .app-standard-modal-copy{flex:1}.variant-booking-pass .app-standard-modal-actions{position:absolute;top:24px;right:28px}.variant-booking-pass .app-standard-modal-icon{width:42px;height:42px;flex-basis:42px;border-color:rgba(255,255,255,.28);background:rgba(255,255,255,.14);color:#dbeafe}.variant-booking-pass .app-standard-modal-icon svg{width:25px;height:25px}.variant-booking-pass .app-standard-modal-copy h2{color:#fff;font-size:1.25rem;font-weight:800}.app-standard-modal-eyebrow{margin-bottom:2px;color:#bfdbfe;font-size:.75rem;line-height:1.4;font-weight:700;letter-spacing:.08em}.app-standard-modal-meta{margin-top:10px}.app-standard-modal-meta :deep(.booking-pass-reference){display:inline-flex;align-items:center;gap:8px;max-width:100%;padding:5px 12px;border:1px solid rgba(255,255,255,.3);border-radius:999px;background:rgba(255,255,255,.16);color:#fff;font-size:.8125rem;font-weight:500;white-space:nowrap}.app-standard-modal-meta :deep(.booking-pass-reference strong){color:#fef08a;letter-spacing:.03em;white-space:nowrap}.app-standard-modal-meta :deep(.booking-pass-copy){display:inline-flex;align-items:center;gap:4px;padding:0;border:0;background:none;color:#fff;font:inherit;font-weight:600;text-decoration:underline;white-space:nowrap;cursor:pointer}.variant-booking-pass .app-standard-modal-close{border:0!important;background:rgba(255,255,255,.16)!important;color:#fff!important;backdrop-filter:blur(4px)}.variant-booking-pass .app-standard-modal-close:hover{background:rgba(255,255,255,.26)!important}
@media(max-width:640px){.app-standard-modal-header{min-height:92px;padding:18px 20px}.app-standard-modal-heading{gap:12px}.app-standard-modal-icon{width:44px;height:44px;flex-basis:44px}.app-standard-modal-copy h2{font-size:1.125rem}.app-standard-modal-copy p{font-size:.875rem}.variant-booking-pass .app-standard-modal-actions{top:18px;right:20px}.variant-booking-pass .app-standard-modal-meta{margin-left:-56px}.app-standard-modal-meta :deep(.booking-pass-reference){font-size:.75rem;gap:6px}}
</style>
