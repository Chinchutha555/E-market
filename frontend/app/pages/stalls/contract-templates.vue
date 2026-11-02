<template>
  <div>
    <AppHeader />
    <main class="content-body">
      <header class="page-heading">
        <div>
          <h1 class="page-title">จัดการรูปแบบสัญญา</h1>
          <p class="page-desc">กำหนดรูปแบบและประเภทสัญญาที่ใช้ในตลาด</p>
        </div>
        <div v-if="!loading && !loadError" class="heading-actions">
          <span class="save-state" :class="{ dirty: isDirty, saved: !isDirty && templates.length }">
            <i></i>{{ isDirty ? 'มีการแก้ไขที่ยังไม่ได้บันทึก' : templates.length ? 'บันทึกแล้ว' : 'ยังไม่มีรูปแบบสัญญา' }}
          </span>
          <AppButton size="md" :variant="isDirty ? 'primary' : 'outline'" :loading="saving" :disabled="!isDirty" @click="save">
            <template #prefix><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><path d="M17 21v-8H7v8M7 3v5h8"/></svg></template>
            บันทึกการตั้งค่า
          </AppButton>
        </div>
      </header>

      <div v-if="loading" class="notice" role="status">กำลังโหลดรูปแบบสัญญา...</div>
      <div v-else-if="loadError" class="error-box" role="alert">{{ loadError }}</div>
      <section v-else class="settings-card">
        <div class="card-header">
          <div>
            <h2>รูปแบบสัญญา</h2>
            <p>{{ templates.length }} รูปแบบ</p>
          </div>
          <AppButton size="md" variant="primary" :disabled="templates.length >= 100" @click="openCreate">
            <template #prefix><svg viewBox="0 0 24 24"><path d="M12 5v14M5 12h14"/></svg></template>
            เพิ่มรูปแบบสัญญา
          </AppButton>
        </div>

        <AppSecondaryTabs v-model="statusFilter" :options="tabs" aria-label="กรองสถานะรูปแบบสัญญา" />

        <div v-if="filteredTemplates.length" class="table-wrap">
          <table>
            <thead><tr><th>รูปแบบสัญญา</th><th>สัญญาที่ใช้</th><th>รายละเอียด</th><th>สถานะ</th><th class="actions-heading">จัดการ</th></tr></thead>
            <tbody>
              <tr v-for="item in filteredTemplates" :key="item.Code">
                <td><strong>{{ item.Title }}</strong><small>{{ item.Code }}</small></td>
                <td><span class="kind-pill">{{ kindLabel(item.ContractKind) }}</span></td>
                <td class="description">{{ item.Description || '-' }}</td>
                <td><span class="status-pill" :class="{ inactive: !item.IsActive }">{{ item.IsActive ? 'เปิดใช้งาน' : 'ปิดใช้งาน' }}</span></td>
                <td>
                  <div class="actions">
                    <AppButton size="sm" variant="outline" @click="previewTemplate = item">
                      <template #prefix>
                        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6M8 13h8M8 17h8"/></svg>
                      </template>
                      ตัวอย่างเอกสาร
                    </AppButton>
                    <AppButton size="sm" variant="outline" @click="openEdit(item)">แก้ไข</AppButton>
                    <AppButton size="sm" variant="danger-outline" @click="remove(item)">ลบ</AppButton>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <AppEmptyState v-else variant="table" icon="document" title="ยังไม่มีรูปแบบสัญญา" description="เพิ่มรูปแบบสัญญาที่ตลาดต้องการใช้งาน">
          <template #actions><AppButton size="sm" variant="primary" @click="openCreate">เพิ่มรูปแบบสัญญา</AppButton></template>
        </AppEmptyState>
        <p v-if="error" class="form-error" role="alert">{{ error }}</p>
      </section>
    </main>

    <AppModal v-model="showModal" :title="editingCode ? 'แก้ไขรูปแบบสัญญา' : 'เพิ่มรูปแบบสัญญา'" description="ระบุชื่อและเลือกประเภทสัญญาที่ต้องการใช้" icon="document">
      <div class="form-grid">
        <label class="field field-wide">
          <span>ชื่อรูปแบบสัญญา <b>*</b></span>
          <AppInput v-model="form.Title" maxlength="255" placeholder="เช่น สัญญาเช่ารายปี" />
        </label>
        <label class="field">
          <span>รหัสรูปแบบ <b>*</b></span>
          <AppInput v-model="form.Code" maxlength="40" placeholder="เช่น annual_contract" :disabled="Boolean(editingCode)" />
        </label>
        <label class="field">
          <span>สัญญาที่ใช้ <b>*</b></span>
          <AppSelect v-model="form.ContractKind" :options="kindOptions" placeholder="เลือกประเภทสัญญา" />
        </label>
        <label class="field">
          <span>ลำดับการแสดง</span>
          <AppInput v-model="form.SortOrder" type="number" min="0" />
        </label>
        <div class="field active-field">
          <span>สถานะ</span>
          <AppCheckbox v-model="form.IsActive" label="เปิดใช้งานรูปแบบนี้" />
        </div>
        <label class="field field-wide">
          <span>รายละเอียด</span>
          <textarea v-model="form.Description" maxlength="1000" rows="4" placeholder="อธิบายการใช้งานรูปแบบสัญญานี้"></textarea>
        </label>
        <p v-if="modalError" class="form-error field-wide" role="alert">{{ modalError }}</p>
      </div>
      <template #footer>
        <AppButton size="md" variant="secondary" @click="showModal = false">ยกเลิก</AppButton>
        <AppButton size="md" variant="primary" @click="applyForm">{{ editingCode ? 'บันทึกการแก้ไข' : 'เพิ่มรูปแบบ' }}</AppButton>
      </template>
    </AppModal>

    <ContractTemplatePreviewModal :template="previewTemplate" @close="previewTemplate = null" />
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/AuthStore'

definePageMeta({
  middleware: (to) => {
    if (!import.meta.client) return
    const auth = useAuthStore()
    auth.hydrate()
    if (!auth.isLoggedIn) return navigateTo('/auth/login')
    if (auth.role === 'staff' && to.query.market) {
      const assigned = auth.currentUser?.assigned_markets?.some((m: any) => String(m.id) === String(to.query.market || ''))
      if (!assigned) return navigateTo('/staff', { replace: true })
    }
  }
})

type Template = { Code: string; Title: string; ContractKind: string; Description: string; IsActive: boolean; SortOrder: number }
const route = useRoute()
const authStore = useAuthStore()
const { get, put } = useApi()
const { showConfirm, showSuccess } = useAppAlert()
const market = computed(() => String(route.query.market || authStore.currentUser?.assigned_markets?.[0]?.id || '1'))
const templates = ref<Template[]>([])
const savedSnapshot = ref('[]')
const loading = ref(true)
const saving = ref(false)
const loadError = ref('')
const error = ref('')
const statusFilter = ref('all')
const showModal = ref(false)
const previewTemplate = ref<Template | null>(null)
const editingCode = ref('')
const modalError = ref('')
const blank = (): Template => ({ Code: '', Title: '', ContractKind: 'standard', Description: '', IsActive: true, SortOrder: templates.value.length + 1 })
const form = ref<Template>(blank())
const kindOptions = [
  { value: 'standard', label: 'สัญญาเช่าพื้นที่มาตรฐาน' },
  { value: 'trial', label: 'สัญญาทดลองขาย' },
  { value: 'temporary', label: 'สัญญาเช่าชั่วคราว' },
  { value: 'custom', label: 'สัญญาแบบกำหนดเอง' }
]
const kindLabel = (value: string) => kindOptions.find(option => option.value === value)?.label || value
const isDirty = computed(() => savedSnapshot.value !== JSON.stringify(templates.value))
const tabs = computed(() => [
  { value: 'all', label: 'ทั้งหมด', count: templates.value.length },
  { value: 'active', label: 'เปิดใช้งาน', count: templates.value.filter(item => item.IsActive).length },
  { value: 'inactive', label: 'ปิดใช้งาน', count: templates.value.filter(item => !item.IsActive).length }
])
const filteredTemplates = computed(() => templates.value.filter(item => statusFilter.value === 'all' || (statusFilter.value === 'active' ? item.IsActive : !item.IsActive)))

const defaultChatuchakTemplate: Template = {
  Code: 'jj_semi_permanent_lease',
  Title: 'สัญญาให้เช่าแผงค้าภายในบริเวณตลาดนัดจตุจักร (แผงค้ากึ่งถาวร)',
  ContractKind: 'standard',
  Description: 'บันทึกข้อตกลงสัญญาให้เช่าแผงค้ากึ่งถาวรภายในบริเวณตลาดนัดจตุจักร สำนักงานตลาดกรุงเทพมหานคร กำหนดระยะเวลา 1 ปี',
  IsActive: true,
  SortOrder: 1
}

const load = async () => {
  loading.value = true
  loadError.value = ''
  if (!market.value) { loadError.value = 'ไม่พบตลาดที่ต้องการตั้งค่า'; loading.value = false; return }
  try {
    const response: any = await get('/markets/' + encodeURIComponent(market.value) + '/contract-templates')
    const list: Template[] = (Array.isArray(response.data) ? response.data : []).map((row: any) => ({
      Code: row.code, Title: row.title, ContractKind: row.contract_kind,
      Description: row.description || '', IsActive: row.is_active !== false, SortOrder: Number(row.sort_order) || 0
    }))
    if (market.value === '1' && !list.some(item => item.Code === defaultChatuchakTemplate.Code)) {
      list.unshift({ ...defaultChatuchakTemplate })
    }
    templates.value = list
    savedSnapshot.value = JSON.stringify(templates.value)
  } catch { loadError.value = 'ไม่สามารถโหลดรูปแบบสัญญาได้ กรุณาลองใหม่' }
  finally { loading.value = false }
}
const openCreate = () => { editingCode.value = ''; form.value = blank(); modalError.value = ''; showModal.value = true }
const openEdit = (item: Template) => { editingCode.value = item.Code; form.value = { ...item }; modalError.value = ''; showModal.value = true }
const applyForm = () => {
  modalError.value = ''
  form.value.Code = form.value.Code.trim().toLowerCase()
  form.value.Title = form.value.Title.trim()
  form.value.Description = form.value.Description.trim()
  form.value.SortOrder = Math.max(0, Number(form.value.SortOrder) || 0)
  if (!form.value.Title || !/^[a-z0-9_-]{1,40}$/.test(form.value.Code) || !kindOptions.some(option => option.value === form.value.ContractKind)) {
    modalError.value = 'กรุณากรอกชื่อ รหัส และประเภทสัญญาให้ครบ'
    return
  }
  const duplicate = templates.value.some(item => item.Code !== editingCode.value && (item.Code === form.value.Code || item.Title.toLowerCase() === form.value.Title.toLowerCase()))
  if (duplicate) { modalError.value = 'ชื่อหรือรหัสรูปแบบสัญญานี้มีอยู่แล้ว'; return }
  const index = templates.value.findIndex(item => item.Code === editingCode.value)
  if (index >= 0) templates.value[index] = { ...form.value }
  else templates.value.push({ ...form.value })
  templates.value.sort((a, b) => a.SortOrder - b.SortOrder)
  showModal.value = false
}
const remove = async (item: Template) => {
  const confirmed = await showConfirm('ต้องการลบรูปแบบ “' + item.Title + '” หรือไม่', { title: 'ลบรูปแบบสัญญา', type: 'danger', icon: 'trash', confirmText: 'ลบ' })
  if (confirmed) templates.value = templates.value.filter(row => row.Code !== item.Code)
}
const save = async () => {
  if (!isDirty.value || saving.value) return
  error.value = ''
  saving.value = true
  try {
    await put('/markets/' + encodeURIComponent(market.value) + '/contract-templates', {
      templates: templates.value.map(item => ({
        code: item.Code, title: item.Title, contract_kind: item.ContractKind,
        description: item.Description, is_active: item.IsActive, sort_order: item.SortOrder
      }))
    })
    savedSnapshot.value = JSON.stringify(templates.value)
    await showSuccess('บันทึกรูปแบบสัญญาแล้ว', 'สำเร็จ')
  } catch (err: any) { error.value = err?.data?.message || 'ไม่สามารถบันทึกรูปแบบสัญญาได้' }
  finally { saving.value = false }
}
onMounted(() => { load(); watch(market, load) })
</script>

<style scoped>
.page-heading{display:flex;align-items:center;justify-content:space-between;gap:20px;flex-wrap:wrap;margin:8px 0 18px}.page-title{margin:0;color:var(--tx-0f172a);font-size:var(--font-size-2xl);font-weight:700}.page-desc{margin:4px 0 0;color:var(--tx-64748b);font-size:var(--font-size-sm);font-weight:500}.heading-actions,.card-header,.actions{display:flex;align-items:center;gap:10px}.save-state{display:inline-flex;align-items:center;gap:7px;padding:8px 12px;border:1px solid var(--bd-e2e8f0);border-radius:999px;color:var(--tx-64748b);font-size:var(--font-size-xs);font-weight:600}.save-state i{width:8px;height:8px;border-radius:50%;background:currentColor}.save-state.dirty{border-color:var(--bd-fed7aa);background:var(--bg-fff7ed);color:var(--tx-c2410c)}.save-state.saved{border-color:var(--bd-bbf7d0);background:var(--bg-ecfdf5);color:var(--tx-15803d)}.settings-card{padding:24px;border:1px solid var(--bd-e2e8f0);border-radius:16px;background:var(--bg-ffffff);box-shadow:0 1px 3px rgba(0,0,0,.03)}.card-header{justify-content:space-between;margin-bottom:16px}.card-header h2{margin:0;color:var(--tx-1e293b);font-size:var(--font-size-lg);font-weight:700}.card-header p{margin:4px 0 0;color:var(--tx-64748b);font-size:var(--font-size-sm)}.table-wrap{overflow-x:auto;margin-top:8px}table{width:100%;border-collapse:collapse}th{padding:14px 16px;border-bottom:1px solid var(--bd-f1f5f9);color:var(--tx-64748b);font-size:var(--font-size-xs);font-weight:500;text-align:left}td{padding:16px;border-bottom:1px solid var(--bd-f1f5f9);color:var(--tx-334155);font-size:var(--font-size-sm);font-weight:500}tbody tr:hover{background:var(--bg-f8fafc)}td strong,td small{display:block}td strong{color:var(--tx-1e293b);font-size:var(--font-size-base)}td small{margin-top:4px;color:var(--tx-94a3b8);font-size:var(--font-size-xs)}.description{max-width:360px}.kind-pill,.status-pill{display:inline-flex;padding:5px 10px;border-radius:7px;background:var(--bg-eff6ff);color:var(--tx-2563eb);font-size:var(--font-size-xs);font-weight:600}.status-pill{background:var(--bg-dcfce7);color:var(--tx-15803d)}.status-pill.inactive{background:var(--bg-f1f5f9);color:var(--tx-64748b)}.actions-heading{text-align:right}.actions{justify-content:flex-end}.form-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:16px}.field{display:flex;min-width:0;flex-direction:column;gap:8px}.field>span{color:var(--tx-334155);font-size:var(--font-size-sm);font-weight:600}.field b,.form-error{color:var(--tx-dc2626)}.field-wide{grid-column:1/-1}.active-field{justify-content:flex-end;padding-bottom:8px}.field textarea{box-sizing:border-box;width:100%;padding:11px 14px;border:1px solid var(--bd-cbd5e1);border-radius:10px;resize:vertical;color:var(--tx-1e293b);font:500 var(--font-size-base)/1.5 inherit;outline:none}.field textarea:focus{border-color:#2563eb;box-shadow:0 0 0 3px rgba(37,99,235,.12)}.form-error,.notice,.error-box{font-size:var(--font-size-sm)}.form-error{margin:0}.error-box{padding:14px;border:1px solid var(--bd-fca5a5);border-radius:10px;background:var(--bg-fef2f2);color:var(--tx-dc2626)}.heading-actions svg,.card-header svg{width:16px;height:16px;fill:none;stroke:currentColor;stroke-width:2;stroke-linecap:round;stroke-linejoin:round}@media(max-width:700px){.form-grid{grid-template-columns:1fr}.field-wide{grid-column:auto}.heading-actions{width:100%}.settings-card{padding:16px}.page-title{font-size:var(--font-size-xl)}}
</style>
