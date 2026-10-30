<template>
  <Teleport to="body">
    <div v-if="open" class="rent-overlay" @click.self="$emit('cancel')">
      <div class="rent-dialog" role="dialog" aria-modal="true" aria-labelledby="rent-confirm-title">
        <div class="rent-icon"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M12 11v5M12 8h.01"/></svg></div>
        <h3 id="rent-confirm-title">ยืนยันการเช่า</h3>
        <p>ยืนยันการเช่าแผง {{ stall || '-' }} ของ {{ booking?.tenant_name || booking?.applicant_name || '-' }}</p>
        <div class="rent-fields">
          <p v-if="!stall" class="rent-error">กรุณาเลือกรายการแผงที่ต้องการยืนยันจากหน้ารายการจอง</p>
          <label>ประเภทการค้า</label>
          <div class="rent-choices">
            <AppCheckbox :model-value="rentType === 'trial'" label="ทดลองขาย" @update:model-value="rentType = $event ? 'trial' : ''" />
            <AppCheckbox :model-value="rentType === 'temporary'" label="ทำการค้าแบบชั่วคราว" @update:model-value="rentType = $event ? 'temporary' : ''" />
            <AppCheckbox :model-value="rentType === 'contract'" label="ทำการค้าแบบมีสัญญา" @update:model-value="rentType = $event ? 'contract' : ''" />
          </div>
        </div>
        <div class="rent-actions"><AppButton variant="secondary" @click="$emit('cancel')">ยกเลิก</AppButton><AppButton variant="primary" :disabled="!stall || !rentType || loading" :loading="loading" @click="$emit('confirm', { stallIdent: stall, rentType })">ยืนยันการเช่า</AppButton></div>
      </div>
    </div>
  </Teleport>
</template>
<script setup lang="ts">
const props = defineProps<{ open: boolean; booking: any; loading?: boolean }>()
defineEmits<{ cancel: []; confirm: [value: { stallIdent: string; rentType: string }] }>()
const stall = ref('')
const rentType = ref('')
const stallOptions = computed(() => {
  const b = props.booking
  if (!b) return []
  const items = Array.isArray(b.sub_bookings) && b.sub_bookings.length ? b.sub_bookings : [b]
  return items.filter((item: any) => !['cancelled', 'rejected', 'not_interested', 'no_show', 'terminated', 'interested', 'lease_applied', 'completed'].includes(String(item.status || '')))
    .flatMap((item: any) => String(item.stall_ident || item.stall_number || '').split(',').map((code: string) => code.trim()).filter(Boolean))
    .map((code: string) => ({ value: code, label: code }))
})
watch(() => props.open, (open) => { if (open) { rentType.value = ''; const code = String(props.booking?.stall_ident || props.booking?.stall_number || '').trim(); stall.value = stallOptions.value.length === 1 ? stallOptions.value[0].value : (stallOptions.value.some((option: any) => option.value === code) ? code : '') } })
</script>
<style scoped>
.rent-overlay { position: fixed; inset: 0; z-index: 2000; display: grid; place-items: center; padding: 16px; background: rgba(15,23,42,.55); }
.rent-dialog { width: min(100%, 520px); padding: 28px; border-radius: 16px; background: #fff; box-shadow: 0 20px 35px rgba(15,23,42,.2); }
.rent-icon { width: 54px; height: 54px; display: grid; place-items: center; margin: 0 auto 16px; color: #2563eb; background: #eff6ff; border-radius: 14px; }
h3 { margin: 0; text-align: center; font-size: 1.25rem; color: #0f172a; } p { color: #64748b; text-align: center; font-size: .9375rem; line-height: 1.5; }
.rent-error { margin: 0; color: #dc2626; text-align: left; }
.rent-fields { display: grid; gap: 10px; margin: 20px 0; } label { font-size: .875rem; font-weight: 600; color: #334155; }
.rent-choices { display: flex; flex-wrap: wrap; gap: 12px 20px; } .rent-actions { display: flex; justify-content: flex-end; gap: 10px; }
</style>
