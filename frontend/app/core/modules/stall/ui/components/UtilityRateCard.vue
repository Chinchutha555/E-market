<template>
  <section class="urc" :class="{ 'is-compact': compact }" :aria-label="title">
    <header class="urc-head">
      <svg v-if="kind === 'water'" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 2.7s6 6.3 6 10.8a6 6 0 0 1-12 0c0-4.5 6-10.8 6-10.8z"/><path d="M9.5 14.5a2.7 2.7 0 0 0 2.5 2"/></svg>
      <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
      <h4>{{ title }}</h4>
    </header>

    <div class="urc-modes" :class="{ 'is-grid': allowInherit }" role="radiogroup" :aria-label="`วิธีคิด${title}`">
      <AppRadio v-if="allowInherit" v-model="rule.mode" :name="groupName" value="inherit" :label="inheritLabel" />
      <AppRadio v-model="rule.mode" :name="groupName" value="meter" label="ตามมิเตอร์" />
      <AppRadio v-model="rule.mode" :name="groupName" value="flat" label="เหมาจ่าย" />
      <AppRadio v-model="rule.mode" :name="groupName" value="none" :label="`ไม่คิด${title}`" />
    </div>
    <p v-if="error && rule.mode === 'inherit'" class="urc-error" role="alert">{{ error }}</p>

    <div v-else-if="rule.mode === 'meter'" class="urc-fields">
      <label class="urc-field"><span>อัตรา (บาท/{{ unit }})<span class="req-star" aria-hidden="true">*</span></span>
        <AppInput v-model="rule.rate" type="number" inputmode="decimal" min="0.01" step="0.01" placeholder="เช่น 8.00" required :error="error" />
      </label>
      <label class="urc-field"><span>ค่าใช้จ่ายเพิ่มต่อหน่วย (บาท/{{ unit }}) ถ้ามี</span>
        <AppInput v-model="rule.surcharge" type="number" inputmode="decimal" min="0" step="0.01" placeholder="0.00" />
      </label>
    </div>

    <div v-else-if="rule.mode === 'flat'" class="urc-fields">
      <label class="urc-field"><span>ยอดเหมา (บาท/เดือน)<span class="req-star" aria-hidden="true">*</span></span>
        <AppInput v-model="rule.flat" type="number" inputmode="decimal" min="0.01" step="0.01" placeholder="เช่น 500.00" required :error="error" />
      </label>
      <label class="urc-field"><span>เหมารายวัน แผงชั่วคราว (บาท/วัน)</span>
        <AppInput v-model="rule.flatDaily" type="number" inputmode="decimal" min="0" step="0.01" placeholder="0.00" />
      </label>
    </div>

  </section>
</template>

<script setup lang="ts">
import AppInput from '~/components/AppInput.vue'
import AppRadio from '~/components/AppRadio.vue'

// การ์ดตั้งค่าอัตราค่าน้ำ/ค่าไฟ: เลือกวิธีคิด (ตามมิเตอร์ / เหมาจ่าย / ไม่คิดค่า) แล้วกรอกอัตรา
// allowInherit = ระดับโครงการ/แผงที่ให้เลือกใช้ค่าจากระดับบนได้
const props = withDefaults(defineProps<{
  rule: { mode: string; rate: string; minMonthly: string; flat: string; flatDaily: string; surcharge: string }
  kind: 'water' | 'electric'
  title: string
  unit: string
  groupName: string
  error?: string
  allowInherit?: boolean
  inheritLabel?: string
  inheritHint?: string
  compact?: boolean
}>(), { error: '', allowInherit: false, compact: false, inheritLabel: 'ค่ามาตรฐาน', inheritHint: '' })


</script>

<style scoped>
.urc { border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-f8fafc); padding: 0 20px 20px; min-width: 0; box-sizing: border-box; }
.urc-head { display: flex; align-items: center; gap: 12px; padding: 18px 0 14px; border-bottom: 1px solid var(--bd-e2e8f0); color: var(--tx-334155); }
.urc-head h4 { margin: 0; color: var(--tx-1e293b); font-size: 1.0625rem; font-weight: 700; }
.urc-modes { display: flex; flex-wrap: wrap; gap: 10px 28px; padding-top: 16px; }
.urc-fields { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 14px; margin-top: 16px; align-items: start; }
.urc-field { display: flex; flex-direction: column; gap: 8px; min-width: 0; color: var(--tx-334155); font-size: 0.875rem; font-weight: 600; }
/* ป้ายกำกับสูงเท่ากันทุกช่อง (2 บรรทัด) เพื่อให้ช่องกรอกอยู่แนวเดียวกันแม้ข้อความยาวไม่เท่ากัน */
.urc-field > span { display: flex; align-items: flex-end; min-height: 2.6em; line-height: 1.3; }
.urc-hint { margin: 14px 0 0; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.urc-error { margin: 12px 0 0; color: var(--tx-dc2626); font-size: 0.8125rem; font-weight: 600; line-height: 1.4; }
.urc.is-compact { padding: 0 16px 16px; }
.urc.is-compact .urc-head { padding: 14px 0 12px; }
.urc.is-compact .urc-modes { padding-top: 10px; gap: 6px 18px; }
.urc.is-compact .urc-fields { margin-top: 14px; }
.urc.is-compact .urc-field > span { min-height: 1.4em; }
.urc-modes.is-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px 20px; }
.urc-modes:not(.is-grid) { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 12px 20px; }
.urc-modes :deep(.app-radio) { width: 100%; }
</style>
