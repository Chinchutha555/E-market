<template>
  <div>
    <AppHeader />
    <main class="content-body market-edit-page">
      <AppBreadcrumb
        :market="originalMarket.name"
        :marketId="originalMarket.id || marketKey"
        current="แก้ไขข้อมูลตลาด"
        :backTo="detailLink"
      />

      <div class="edit-title-row">
        <div>
          <h1>แก้ไขข้อมูลตลาด</h1>
          <p>ข้อมูลที่บันทึกจะนำไปแสดงในหน้ารายละเอียดตลาด</p>
        </div>
        <span v-if="lastUpdated" class="saved-state">บันทึกล่าสุด {{ formatUpdatedAt(lastUpdated) }}</span>
      </div>

      <button type="button" class="timeline-edge-trigger" :class="{ active: showEditHistory }" :aria-expanded="showEditHistory" aria-controls="market-edit-history-panel" @click="showEditHistory = !showEditHistory"><svg width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 3-6.7L3 8"/><path d="M3 3v5h5"/><path d="M12 7v5l3 2"/></svg><span>ประวัติการแก้ไข</span></button>

      <Transition name="timeline-slide">
        <div v-if="showEditHistory" class="timeline-backdrop" @click.self="showEditHistory = false">
          <aside id="market-edit-history-panel" class="timeline-panel" role="dialog" aria-modal="true" aria-label="ประวัติการแก้ไขข้อมูลตลาด">
            <header class="timeline-header">
              <div><span class="eyebrow">รายการกิจกรรมทั้งหมด</span><h2>ประวัติการแก้ไข</h2></div>
              <button type="button" class="close-button" @click="showEditHistory = false" aria-label="ปิดประวัติการแก้ไข">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
              </button>
            </header>
            <div class="timeline-body">
              <div v-if="editHistory.length" class="timeline-list">
                <div v-for="entry in editHistory" :key="entry.id" class="timeline-item">
                  <span class="timeline-dot"></span>
                  <div class="timeline-copy"><strong>{{ entry.changedFields.join(' · ') }}</strong><time>{{ formatUpdatedAt(entry.changedAt) }}</time><div class="timeline-actor">{{ entry.actorRole }} · {{ entry.actorName }}</div><p>แก้ไขข้อมูลตลาด</p></div>
                </div>
              </div>
              <AppEmptyState v-else variant="compact" icon="history" title="ยังไม่มีประวัติการแก้ไข" description="การเปลี่ยนแปลงข้อมูลตลาดจะแสดงที่นี่" />
            </div>
          </aside>
        </div>
      </Transition>

      <form v-if="!loading" class="edit-form" novalidate @submit.prevent="saveMarket">
        <section class="form-card">
          <div class="section-heading"><div class="section-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 21h18M5 21V7l7-4 7 4v14M9 9h2m2 0h2m-6 4h2m2 0h2m-6 4h2m2 0h2" /></svg></div><div><h2>ข้อมูลตลาด</h2><p>ชื่อ ที่ตั้ง เวลาเปิดทำการ และข้อมูลสรุปบนหน้า Detail</p></div></div>
          <div class="form-grid three-cols">
            <label><span>ชื่อตลาด <b>*</b></span><AppInput v-model="form.name" placeholder="กรอกชื่อตลาด" :error="errors.name" @blur="validateField('name')" /></label>
            <label><span>เขต <b>*</b></span><AppInput v-model="form.district" placeholder="เช่น เขตจตุจักร" :error="errors.district" @blur="validateField('district')" /></label>
            <label><span>วันที่เปิด–ปิด <b>*</b></span><AppSelect v-model="form.operating_days" :options="dayOptions" /></label>
            <label><span>เวลาเปิด <b>*</b></span><AppTimePicker v-model="form.open_time" placeholder="เลือกเวลาเปิด" /></label>
            <label><span>เวลาปิด <b>*</b></span><AppTimePicker v-model="form.close_time" placeholder="เลือกเวลาปิด" /></label>
            <label class="span-2"><span>ที่อยู่ <b>*</b></span><AppInput v-model="form.address" placeholder="กรอกที่อยู่ตลาด" :error="errors.address" @blur="validateField('address')" /></label>
            <label class="span-3"><span>รายละเอียดตลาด <b>*</b></span><AppInput v-model="form.description" type="textarea" :rows="4" placeholder="กรอกรายละเอียดตลาด" :error="errors.description" @blur="validateField('description')" /></label>
            <label class="span-2"><span>จุดเด่นของตลาด</span><AppInput v-model="form.tags_text" placeholder="คั่นแต่ละรายการด้วยเครื่องหมายจุลภาค" hint="เช่น ใกล้รถไฟฟ้า, มีที่จอดรถ, อาหารและเครื่องดื่ม" /></label>
            <label><span>การเดินทาง</span><AppInput v-model="form.travel_text" placeholder="คั่นแต่ละรายการด้วยเครื่องหมายจุลภาค" /></label>
            <label class="span-3"><span>หมวดสินค้าของตลาด</span><AppInput v-model="form.categories_text" placeholder="คั่นแต่ละหมวดหมู่ด้วยเครื่องหมายจุลภาค (,)" hint="รายการหมวดสินค้าที่จะแสดงให้เลือกตอนเพิ่มโซนและแผงเช่าในหน้าจัดการผังตลาด" /></label>
            <label><span>ละติจูด <b>*</b></span><AppInput v-model="form.lat" type="number" placeholder="เช่น 13.7999" :error="errors.lat" @blur="validateField('lat')" /></label>
            <label><span>ลองจิจูด <b>*</b></span><AppInput v-model="form.lng" type="number" placeholder="เช่น 100.5508" :error="errors.lng" @blur="validateField('lng')" /></label>
          </div>
        </section>

        <section class="form-card">
          <div class="section-heading"><div class="section-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 3h18v18H3zM3 9h18M9 21V9" /></svg></div><div><h2>พื้นที่และราคา</h2><p>ข้อมูลสถิติที่แสดงใต้ภาพผังตลาด</p></div></div>
          <div class="form-grid four-cols">
            <label><span>ค่าเช่าเริ่มต้น/เดือน <b>*</b></span><AppInput v-model="form.starting_price_monthly" type="number" :error="errors.starting_price_monthly" @blur="validateField('starting_price_monthly')" /></label>
            <label><span>ค่าเช่าสูงสุด/เดือน</span><AppInput v-model="form.max_price_monthly" type="number" /></label>
            <label><span>พื้นที่ต่ำสุด (ตร.ม.)</span><AppInput v-model="form.min_area_sqm" type="number" /></label>
            <label><span>พื้นที่สูงสุด (ตร.ม.)</span><AppInput v-model="form.max_area_sqm" type="number" /></label>
            <label><span>จำนวนแผงทั้งหมด</span><AppInput v-model="form.total_stalls" type="number" /></label>
            <label><span>พื้นที่ว่าง</span><AppInput :modelValue="form.available_stalls" disabled type="number" hint="คำนวณจากสถานะแผงในระบบ" /></label>
            <label><span>จำนวนที่จอดรถ (คัน)</span><AppInput v-model="form.parking_spaces" type="number" /></label>
          </div>
        </section>

        <section class="form-card">
          <div class="section-heading"><div class="section-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><path d="m21 15-5-5L5 21"/></svg></div><div><h2>รูปภาพและวิดีโอของตลาด</h2><p>ใช้ URL รูปจากระบบหรืออัปโหลดไฟล์ใหม่</p></div></div>
          <div class="media-grid">
            <MediaField title="รูปภาพแผนผังตลาด" hint="JPG, JPEG, PNG" :preview="form.floorplan_image_url" @file="file => setImage('floorplan_image_url', file)" />
            <MediaField title="รูปภาพตลาด" hint="รูปหลักที่แสดงในหน้ารายละเอียด" :preview="form.image_url" @file="file => setImage('image_url', file)" />
          </div>
          <div class="form-grid two-cols media-links">
            <label><span>URL รูปภาพแผนผัง</span><AppInput v-model="form.floorplan_image_url" placeholder="/images/market-plan.jpg" /></label>
            <label><span>URL รูปภาพตลาด</span><AppInput v-model="form.image_url" placeholder="/images/market.jpg" /></label>
            <label class="span-2"><span>วิดีโอแนะนำตลาด</span><AppInput v-model="form.video_url" placeholder="กรอกลิงก์วิดีโอแนะนำตลาด" /></label>
          </div>
        </section>

        <div class="form-actions"><NuxtLink :to="detailLink" class="cancel-button">ยกเลิก</NuxtLink><button class="save-button" :disabled="saving" type="submit"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2Z"/><path d="M17 21v-8H7v8M7 3v5h8"/></svg>{{ saving ? 'กำลังบันทึก...' : 'บันทึก' }}</button></div>
      </form>
      <div v-else class="loading-card">กำลังโหลดข้อมูลตลาด...</div>
    </main>
  </div>
</template>

<script setup lang="ts">
import AppInput from '~/components/AppInput.vue'
import AppSelect from '~/components/AppSelect.vue'
import AppTimePicker from '~/components/AppTimePicker.vue'
import MediaField from '~/components/market/MediaField.vue'
import { getMarketGisConfig } from '~/core/modules/stall/data/markets_registry'

const route = useRoute()
const router = useRouter()
const { get } = useApi()
const { user, role } = useAuth()
const { showSuccess, showError } = useAppAlert()
const { applyMarketOverride, saveMarketOverride, getMarketEditHistory, saveMarketEditHistory } = useMarketOverrides()
const loading = ref(true)
const saving = ref(false)
const marketKey = computed(() => String(route.params.id || route.query.market || '1'))
const buildInitialMarket = (key: string) => {
  const cfg = getMarketGisConfig(key)
  return applyMarketOverride({
    id: cfg.id,
    slug: cfg.slug,
    name: cfg.name,
    district: cfg.district,
    description: cfg.description,
    lat: cfg.center[1],
    lng: cfg.center[0]
  })
}
const originalMarket = ref<any>(buildInitialMarket(marketKey.value))
const lastUpdated = ref('')
const showEditHistory = ref(false)
const editHistory = ref<any[]>([])
const detailLink = computed(() => `/stalls?market=${originalMarket.value.id || marketKey.value}`)
const dayOptions = ['ทุกวัน', 'จันทร์ - ศุกร์', 'เสาร์ - อาทิตย์', 'อังคาร - อาทิตย์'].map(value => ({ value, label: value }))
const form = reactive<any>({})
const errors = reactive({
  name: '',
  district: '',
  address: '',
  description: '',
  lat: '',
  lng: '',
  starting_price_monthly: ''
})
const validateField = (field: keyof typeof errors) => {
  if (field === 'name') {
    errors.name = String(form.name || '').trim() ? '' : 'กรุณาระบุชื่อตลาด'
  } else if (field === 'district') {
    errors.district = String(form.district || '').trim() ? '' : 'กรุณาระบุเขตที่ตั้งตลาด'
  } else if (field === 'address') {
    errors.address = String(form.address || '').trim() ? '' : 'กรุณาระบุที่อยู่ตลาด'
  } else if (field === 'description') {
    errors.description = String(form.description || '').trim() ? '' : 'กรุณาระบุรายละเอียดตลาด'
  } else if (field === 'lat') {
    const val = Number(form.lat)
    errors.lat = String(form.lat ?? '').trim() === '' || !Number.isFinite(val) ? 'กรุณาระบุพิกัดละติจูด' : ''
  } else if (field === 'lng') {
    const val = Number(form.lng)
    errors.lng = String(form.lng ?? '').trim() === '' || !Number.isFinite(val) ? 'กรุณาระบุพิกัดลองจิจูด' : ''
  } else if (field === 'starting_price_monthly') {
    const val = Number(form.starting_price_monthly)
    errors.starting_price_monthly = String(form.starting_price_monthly ?? '').trim() === '' || !Number.isFinite(val) || val < 0 ? 'กรุณาระบุค่าเช่าเริ่มต้นต่อเดือน' : ''
  }
}
const validateForm = () => {
  validateField('name')
  validateField('district')
  validateField('address')
  validateField('description')
  validateField('lat')
  validateField('lng')
  validateField('starting_price_monthly')
  return !Object.values(errors).some(Boolean)
}
const refreshEditHistory = () => { editHistory.value = getMarketEditHistory(originalMarket.value.id || marketKey.value) }

const splitOperatingHours = (value = '') => {
  const match = value.match(/^(.*?)\s+(\d{1,2}:\d{2})\s*-\s*(\d{1,2}:\d{2})/)
  return { days: match?.[1]?.trim() || 'ทุกวัน', open: match?.[2] || '05:00', close: match?.[3] || '19:00' }
}
const assignForm = (market: any) => {
  const hours = splitOperatingHours(market.operating_hours)
  const gisCfg = getMarketGisConfig(market.id || market.slug || marketKey.value)
  const defaultCategories = Array.from(new Set((gisCfg?.projects || []).map(p => String(p.category || '').trim()).filter(Boolean)))
  const categoriesList = Array.isArray(market.product_categories) && market.product_categories.length
    ? market.product_categories
    : defaultCategories
  Object.assign(form, market, { operating_days: hours.days, open_time: hours.open, close_time: hours.close, tags_text: (market.tags || []).join(', '), categories_text: categoriesList.join(', '), travel_text: (market.travel_options || ['รถโดยสารประจำทาง', 'รถรับจ้างสาธารณะ', 'รถยนต์ส่วนบุคคล', 'รถไฟฟ้า BTS และ MRT']).join(', '), max_price_monthly: market.max_price_monthly || Math.round(Number(market.starting_price_monthly || 1800) * 2.5), min_area_sqm: market.min_area_sqm || 4, max_area_sqm: market.max_area_sqm || 25, parking_spaces: market.parking_spaces || 120, floorplan_image_url: market.floorplan_image_url || '/images/chatuchak_isometric_plan.jpg', video_url: market.video_url || '' })
  lastUpdated.value = market.updated_at || ''
}
const setImage = (field: 'image_url' | 'floorplan_image_url', file: File) => {
  if (!file.type.startsWith('image/')) return showError('กรุณาเลือกไฟล์รูปภาพเท่านั้น')
  if (file.size > 1024 * 1024) return showError('ไฟล์รูปต้องมีขนาดไม่เกิน 1 MB')
  const reader = new FileReader()
  reader.onload = () => { form[field] = String(reader.result || '') }
  reader.readAsDataURL(file)
}
const saveMarket = async () => {
  if (!validateForm()) return
  saving.value = true
  try {
    const payload = { ...originalMarket.value, ...form, operating_hours: `${form.operating_days} ${form.open_time} - ${form.close_time} น.`, tags: String(form.tags_text || '').split(',').map((item: string) => item.trim()).filter(Boolean), product_categories: String(form.categories_text || '').split(',').map((item: string) => item.trim()).filter(Boolean), travel_options: String(form.travel_text || '').split(',').map((item: string) => item.trim()).filter(Boolean) }
    delete payload.tags_text; delete payload.categories_text; delete payload.travel_text; delete payload.operating_days; delete payload.open_time; delete payload.close_time
    saveMarketEditHistory(payload, originalMarket.value, {
      name: user.value?.full_name || user.value?.username,
      role: role.value === 'admin' ? 'ผู้ดูแลระบบ' : role.value === 'staff' ? 'เจ้าหน้าที่ตลาด' : 'ผู้ใช้งาน'
    })
    saveMarketOverride(payload)
    originalMarket.value = payload
    refreshEditHistory()
    await showSuccess('ข้อมูลใหม่จะแสดงในหน้ารายละเอียดตลาดทันที', 'บันทึกข้อมูลตลาดสำเร็จ')
    await router.push(detailLink.value)
  } catch {
    showError('ไม่สามารถบันทึกข้อมูลตลาดได้ กรุณาลองใหม่อีกครั้ง')
  } finally { saving.value = false }
}
const formatUpdatedAt = (value: string) => new Intl.DateTimeFormat('th-TH', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
onMounted(async () => {
  try {
    const res: any = await get(`/markets/${marketKey.value}`)
    const market = res?.data || res
    originalMarket.value = applyMarketOverride(market)
    assignForm(originalMarket.value)
    refreshEditHistory()
  }
  catch { showError('ไม่สามารถโหลดข้อมูลตลาดได้') }
  finally { loading.value = false }
})
</script>

<style scoped>
.timeline-edge-trigger { position: fixed; top: 52%; right: 0; z-index: 900; width: 44px; height: 150px; transform: translateY(-50%); display: block; padding: 0; overflow: hidden; border: 1px solid var(--bd-dbeafe); border-right: 0; border-radius: 12px 0 0 12px; background: var(--bg-rgba-255-255-255-0p98); color: var(--tx-2563eb); box-shadow: -6px 7px 20px rgba(15,23,42,.11); backdrop-filter: blur(10px); font-family: inherit; font-size: var(--font-size-xs); font-weight: 700; white-space: nowrap; cursor: pointer; transition: background .18s ease, width .18s ease, box-shadow .18s ease, transform .18s ease; }
.timeline-edge-trigger > svg { position: absolute; top: 15px; left: 50%; width: 17px; height: 17px; transform: translateX(-50%); }.timeline-edge-trigger > span { position: absolute; top: 56%; left: 50%; transform: translate(-50%,-50%) rotate(-90deg); transform-origin: center; }.timeline-edge-trigger:hover, .timeline-edge-trigger.active { width: 48px; background: var(--bg-eff6ff); color: var(--tx-1d4ed8); box-shadow: -8px 9px 24px rgba(37,99,235,.16); }
.timeline-backdrop { position: fixed; inset: 0; z-index: 1000; display: flex; justify-content: flex-end; background: rgba(15,23,42,.45); backdrop-filter: blur(2px); }.timeline-panel { width: min(480px, 92vw); height: 100%; background: var(--bg-ffffff); box-shadow: -16px 0 40px rgba(15,23,42,.16); overflow-y: auto; }.timeline-header { position: sticky; top: 0; z-index: 2; display: flex; align-items: center; justify-content: space-between; padding: 22px 24px; border-bottom: 1px solid var(--bd-e2e8f0); background: var(--bg-ffffff); }.timeline-header h2 { margin: 0; font-size: var(--font-size-xl); }.close-button { width: 36px; height: 36px; display: grid; place-items: center; border: 1px solid var(--bd-e2e8f0); border-radius: 9px; background: var(--bg-ffffff); color: var(--tx-475569); cursor: pointer; }.timeline-body { padding: 24px; }.timeline-list { display: flex; flex-direction: column; }.timeline-item { position: relative; display: flex; gap: 14px; padding-bottom: 24px; }.timeline-item:not(:last-child)::before { content: ''; position: absolute; left: 5px; top: 13px; bottom: 0; width: 1px; background: var(--bg-e2e8f0); }.timeline-dot { position: relative; z-index: 1; width: 11px; height: 11px; margin-top: 5px; flex: 0 0 11px; border-radius: 50%; background: #3b82f6; box-shadow: 0 0 0 4px #dbeafe; }.timeline-copy { flex: 1; min-width: 0; }.timeline-copy strong { display: block; font-size: var(--font-size-sm); }.timeline-copy time { display: block; color: var(--tx-64748b); font-size: var(--font-size-xs); margin-top: 3px; }.timeline-copy p { margin: 8px 0 0; padding: 8px 10px; border-radius: 8px; background: var(--bg-f8fafc); color: var(--tx-475569); font-size: var(--font-size-xs); line-height: 1.45; }.timeline-empty { padding: 30px 10px; text-align: center; color: var(--tx-94a3b8); font-size: var(--font-size-sm); }
.timeline-slide-enter-active, .timeline-slide-leave-active { transition: background-color .28s ease, backdrop-filter .28s ease; }.timeline-slide-enter-active .timeline-panel, .timeline-slide-leave-active .timeline-panel { transition: transform .32s cubic-bezier(.22,1,.36,1), box-shadow .32s ease; }.timeline-slide-enter-from, .timeline-slide-leave-to { background: var(--bg-rgba-15-23-42-0); backdrop-filter: blur(0); }.timeline-slide-enter-from .timeline-panel, .timeline-slide-leave-to .timeline-panel { transform: translateX(100%); box-shadow: none; }
.market-edit-page{padding-bottom:48px}.edit-breadcrumb{display:flex;align-items:center;gap:10px;color:var(--tx-94a3b8);font-size:var(--font-size-xs);margin-bottom:24px}.back-button{width:38px;height:38px;display:grid;place-items:center;border:1px solid var(--bd-e2e8f0);border-radius:9px;color:var(--tx-475569);margin-right:8px}.back-button svg{width:18px}.crumb-link{color:var(--tx-64748b)}.edit-breadcrumb strong{color:var(--tx-334155)}.edit-title-row{display:flex;justify-content:space-between;align-items:end;margin-bottom:20px}.eyebrow{font-size:var(--font-size-xs);letter-spacing:.08em;color:var(--tx-2563eb);font-weight:800}.edit-title-row h1{font-size:1.55rem;color:var(--tx-1e293b);margin:4px 0}.edit-title-row p{color:var(--tx-64748b);margin:0;font-size:var(--font-size-sm)}.saved-state{font-size:var(--font-size-xs);color:var(--tx-15803d);background:var(--bg-ecfdf5);padding:7px 11px;border-radius:999px}.edit-history-edge-trigger{position:fixed;top:27%;right:0;z-index:900;width:44px;height:188px;transform:translateY(-50%);display:block;padding:0;overflow:hidden;border:1px solid var(--bd-dbeafe);border-right:0;border-radius:12px 0 0 12px;background:var(--bg-rgba-255-255-255-0p98);color:var(--tx-2563eb);box-shadow:-6px 7px 20px rgba(15,23,42,.11);backdrop-filter:blur(10px);font:700 .75rem inherit;white-space:nowrap;cursor:pointer;transition:background .18s ease,width .18s ease,box-shadow .18s ease}.edit-history-edge-trigger>svg{position:absolute;top:15px;left:50%;transform:translateX(-50%)}.edit-history-edge-trigger>span{position:absolute;top:56%;left:50%;transform:translate(-50%,-50%) rotate(-90deg);transform-origin:center}.edit-history-edge-trigger:hover,.edit-history-edge-trigger.active{width:48px;background:var(--bg-eff6ff);color:var(--tx-1d4ed8);box-shadow:-8px 9px 24px rgba(37,99,235,.16)}.edit-history-backdrop{position:fixed;inset:0;z-index:1000;display:flex;justify-content:flex-end;background:rgba(15,23,42,.45);backdrop-filter:blur(2px)}.edit-history-panel{width:min(480px,92vw);height:100%;overflow-y:auto;background:var(--bg-ffffff);box-shadow:-16px 0 40px rgba(15,23,42,.16)}.edit-history-header{position:sticky;top:0;z-index:2;display:flex;align-items:flex-start;justify-content:space-between;gap:16px;padding:22px 24px;border-bottom:1px solid var(--bd-e2e8f0);background:var(--bg-ffffff)}.edit-history-header h2{margin:0;color:var(--tx-1e293b);font-size:var(--font-size-xl)}.edit-history-header p{margin:4px 0 0;color:var(--tx-64748b);font-size:var(--font-size-sm)}.edit-history-close{width:36px;height:36px;display:grid;place-items:center;border:1px solid var(--bd-e2e8f0);border-radius:9px;background:var(--bg-ffffff);color:var(--tx-475569);cursor:pointer}.edit-history-body{padding:24px}.edit-history-list{display:flex;flex-direction:column;margin:0;padding:0;list-style:none}.edit-history-list li{position:relative;display:flex;gap:14px;padding-bottom:24px}.edit-history-list li:not(:last-child)::before{content:'';position:absolute;left:5px;top:13px;bottom:0;width:1px;background:var(--bg-e2e8f0)}.edit-history-dot{position:relative;z-index:1;width:11px;height:11px;margin-top:5px;flex:0 0 11px;border-radius:50%;background:#2563eb;box-shadow:0 0 0 4px #dbeafe}.edit-history-list strong,.edit-history-list time{display:block}.edit-history-list strong{color:var(--tx-334155);font-size:var(--font-size-sm)}.edit-history-list time{margin-top:4px;color:var(--tx-64748b);font-size:var(--font-size-xs)}.edit-history-empty{padding:30px 10px;text-align:center;color:var(--tx-94a3b8);font-size:var(--font-size-sm)}.edit-history-slide-enter-active,.edit-history-slide-leave-active{transition:background-color .28s ease,backdrop-filter .28s ease}.edit-history-slide-enter-active .edit-history-panel,.edit-history-slide-leave-active .edit-history-panel{transition:transform .32s cubic-bezier(.22,1,.36,1),box-shadow .32s ease}.edit-history-slide-enter-from,.edit-history-slide-leave-to{background:var(--bg-rgba-15-23-42-0);backdrop-filter:blur(0)}.edit-history-slide-enter-from .edit-history-panel,.edit-history-slide-leave-to .edit-history-panel{transform:translateX(100%);box-shadow:none}.edit-form{display:grid;gap:20px}.form-card,.loading-card{background:var(--bg-ffffff);border:1px solid var(--bd-e2e8f0);border-radius:16px;padding:26px 30px;box-shadow:0 1px 3px rgba(15,23,42,.03)}.section-heading{display:flex;align-items:center;gap:12px;padding-bottom:18px;margin-bottom:22px;border-bottom:1px solid var(--bd-eef2f6)}.section-icon{width:38px;height:38px;display:grid;place-items:center;border-radius:10px;background:var(--bg-eff6ff);color:var(--tx-2563eb)}.section-icon svg{width:19px}.section-heading h2{font-size:var(--font-size-lg);color:var(--tx-1e293b);margin:0}.section-heading p{margin:3px 0 0;color:var(--tx-64748b);font-size:var(--font-size-xs)}.form-grid{display:grid;gap:18px 22px}.three-cols{grid-template-columns:repeat(3,minmax(0,1fr))}.four-cols{grid-template-columns:repeat(4,minmax(0,1fr))}.two-cols{grid-template-columns:repeat(2,minmax(0,1fr))}.form-grid label{display:flex;flex-direction:column;gap:7px;color:var(--tx-475569);font-size:var(--font-size-xs);font-weight:600}.form-grid b{color:var(--tx-ef4444)}.form-grid input,.form-grid textarea{width:100%;box-sizing:border-box;border:1px solid var(--bd-cbd5e1);border-radius:9px;background:var(--bg-ffffff);color:var(--tx-1e293b);padding:11px 13px;font:inherit;font-weight:500;outline:none}.form-grid input:focus,.form-grid textarea:focus{border-color:#60a5fa;box-shadow:0 0 0 3px rgba(37,99,235,.1)}.form-grid input:disabled{background:var(--bg-f1f5f9);color:var(--tx-94a3b8)}.form-grid small{font-weight:400;color:var(--tx-94a3b8)}.span-2{grid-column:span 2}.span-3{grid-column:span 3}.media-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:22px}.media-links{margin-top:22px}.input-with-icon{position:relative}.input-with-icon svg{position:absolute;right:13px;top:50%;transform:translateY(-50%);width:18px;color:var(--tx-64748b)}.input-with-icon input{padding-right:42px}.form-actions{display:flex;justify-content:flex-end;gap:12px}.cancel-button,.save-button{min-width:120px;padding:11px 20px;border-radius:9px;font-weight:700;display:inline-flex;align-items:center;justify-content:center;gap:8px}.cancel-button{border:1px solid var(--bd-cbd5e1);color:var(--tx-475569);background:var(--bg-ffffff)}.save-button{border:1px solid #2563eb;background:#2563eb;color:#fff;box-shadow:0 4px 12px rgba(37,99,235,.22);cursor:pointer}.save-button svg{width:17px}.loading-card{text-align:center;color:var(--tx-64748b)}@media(max-width:960px){.three-cols,.four-cols{grid-template-columns:repeat(2,minmax(0,1fr))}.span-3{grid-column:span 2}}@media(max-width:640px){.three-cols,.four-cols,.two-cols,.media-grid{grid-template-columns:1fr}.span-2,.span-3{grid-column:auto}.edit-title-row{align-items:flex-start;flex-direction:column;gap:12px}.form-card{padding:20px}.form-actions>*{flex:1}.edit-history-edge-trigger{top:32%;width:40px;height:180px;font-size:var(--font-size-xs)}.edit-history-edge-trigger:hover{width:44px}}
.timeline-actor { display: block; margin-top: 3px; color: var(--tx-64748b); font-size: var(--font-size-xs); }
</style>
