<template>
  <Teleport to="body">
    <div v-if="open" class="rent-overlay" @click.self="$emit('cancel')">
      <div class="rent-dialog" role="dialog" aria-modal="true" aria-labelledby="rent-confirm-title">
        <AppModalHeader title="ยืนยันการเช่า" :subtitle="`แผง ${stall || '-'} · ${booking?.tenant_name || booking?.applicant_name || '-'}`" icon="document" @close="$emit('cancel')" />
        <div class="rent-body"><div class="rent-fields">
          <p v-if="submitted && !stall" class="rent-error">กรุณาเลือกรายการแผงที่ต้องการยืนยันจากหน้ารายการจอง</p>
          <label>ประเภทการค้า <span class="req-star" aria-hidden="true">*</span></label>
          <div class="rent-choices">
            <AppCheckbox :model-value="rentType === 'trial'" label="ทดลองขาย" @update:model-value="rentType = $event ? 'trial' : ''" />
            <AppCheckbox :model-value="rentType === 'temporary'" label="ทำการค้าแบบชั่วคราว" @update:model-value="rentType = $event ? 'temporary' : ''" />
            <AppCheckbox :model-value="rentType === 'contract'" label="ทำการค้าแบบมีสัญญา" @update:model-value="rentType = $event ? 'contract' : ''" />
          </div>
          <p v-if="submitted && !rentType" class="rent-error">กรุณาเลือกประเภทการค้า</p>
        </div></div>
        <div class="rent-actions"><AppButton variant="secondary" size="md" @click="$emit('cancel')">ยกเลิก</AppButton><AppButton variant="primary" size="md" :loading="loading" @click="confirmRental">ยืนยันการเช่า</AppButton></div>
      </div>
    </div>
  </Teleport>
</template>
<script setup lang="ts">
const props = defineProps<{ open: boolean; booking: any; loading?: boolean }>()
const emit = defineEmits<{ cancel: []; confirm: [value: { stallIdent: string; rentType: string }] }>()
const stall = ref('')
const rentType = ref('')
const submitted = ref(false)
const stallOptions = computed(() => {
  const b = props.booking
  if (!b) return []
  const items = Array.isArray(b.sub_bookings) && b.sub_bookings.length ? b.sub_bookings : [b]
  return items.filter((item: any) => !['cancelled', 'rejected', 'not_interested', 'no_show', 'terminated', 'interested', 'lease_applied', 'completed'].includes(String(item.status || '')))
    .flatMap((item: any) => String(item.stall_ident || item.stall_number || '').split(',').map((code: string) => code.trim()).filter(Boolean))
    .map((code: string) => ({ value: code, label: code }))
})
const confirmRental = () => {
  submitted.value = true
  if (!stall.value || !rentType.value || props.loading) return
  emit('confirm', { stallIdent: stall.value, rentType: rentType.value })
}
watch(() => props.open, (open) => { if (open) { submitted.value = false; rentType.value = ''; const code = String(props.booking?.stall_ident || props.booking?.stall_number || '').trim(); stall.value = stallOptions.value.length === 1 ? stallOptions.value[0].value : (stallOptions.value.some((option: any) => option.value === code) ? code : '') } })
</script>
<style scoped>
.rent-overlay { position: fixed; inset: 0; z-index: 2000; display: grid; place-items: center; padding: 16px; background: rgba(15,23,42,.55); }
.rent-dialog { width: min(100%, 520px); max-height:min(90vh,90dvh); display:flex; flex-direction:column; overflow:hidden; border-radius: 20px; background: var(--bg-ffffff); box-shadow: 0 20px 35px rgba(15,23,42,.2); }
.rent-body{flex:1 1 auto;min-height:0;overflow-y:auto;overscroll-behavior:contain;padding:24px 28px}
.rent-icon { width: 54px; height: 54px; display: grid; place-items: center; margin: 0 auto 16px; color: var(--tx-2563eb); background: var(--bg-eff6ff); border-radius: 14px; }
h3 { margin: 0; text-align: center; font-size: 1.25rem; color: var(--tx-0f172a); } p { color: var(--tx-64748b); text-align: center; font-size: .9375rem; line-height: 1.5; }
.rent-error { margin: 0; color: var(--tx-dc2626); text-align: left; }
.rent-fields { display: grid; gap: 10px; } label { font-size: .875rem; font-weight: 600; color: var(--tx-334155); }
.rent-choices { display: flex; flex-wrap: wrap; gap: 12px 20px; } .rent-actions { flex:0 0 auto;display: flex; justify-content: flex-end; gap: 10px; padding:16px 28px; border-top:1px solid var(--bd-e2e8f0); background:var(--bg-f8fafc); }
</style>
