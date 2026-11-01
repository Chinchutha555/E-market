<template>
  <div>
    <AppHeader />
    <main class="content-body us-page">
      <header class="us-heading">
        <div>
          <h1 class="page-title">{{ section === 'meters' ? 'จัดการมิเตอร์' : 'อัตราค่าน้ำ-ไฟ' }}</h1>
          <p class="page-desc">{{ section === 'meters' ? 'รหัสมิเตอร์น้ำและไฟของแต่ละแผง' : 'ตั้งวิธีคิดค่าน้ำและค่าไฟ' }}</p>
        </div>
        <div v-if="!loading && (section === 'rates' || (section === 'meters' && activeMeterProject))" class="us-heading-actions">
          <span class="draft-status" :class="{ published: !currentDirty, 'is-dirty': currentDirty }">
            <i></i>
            {{ currentDirty ? 'มีการแก้ไขที่ยังไม่บันทึก' : 'บันทึกแล้ว' }}
          </span>
          <AppButton size="md" :variant="currentDirty ? 'primary' : 'outline'" :loading="saving" @click="saveCurrent">
            <template #prefix>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path><polyline points="17 21 17 13 7 13 7 21"></polyline><polyline points="7 3 7 8 15 8"></polyline></svg>
            </template>
            {{ section === 'meters' ? 'บันทึกมิเตอร์' : 'บันทึกการแก้ไข' }}
          </AppButton>
          <span v-if="utilitySyncError" class="utility-sync-error" role="status">{{ utilitySyncError }}</span>
        </div>
      </header>

      <div v-if="loading" class="us-page-skeleton" aria-label="กำลังโหลดข้อมูล">
        <div class="us-skeleton-title"></div>
        <div class="us-skeleton-cards"><i v-for="n in 3" :key="`card-${n}`"></i></div>
        <div class="us-skeleton-workspace">
          <i class="us-skeleton-sidebar"></i>
          <i class="us-skeleton-content"></i>
        </div>
      </div>

      <div v-else>
        <!-- 1) รูปแบบการคิด -->
        <section v-if="section === 'rates'" class="us-default">
        <h2 class="us-h">รูปแบบการคิด</h2>
        <AppOptionCards v-model="scope" :options="scopeOptions" aria-label="รูปแบบการคิด" @change="onScopeChange" />
        </section>

        <!-- 2) ค่ามาตรฐาน -->
        <section v-if="section === 'rates'" class="us-default">
        <h2 class="us-h">ค่ามาตรฐาน</h2>
        <div class="us-default-grid">
          <UtilityRateCard v-for="u in utilityFields" :key="u.key" :rule="defaultForm[u.key]" :kind="u.key" :title="u.label" :unit="u.unit" :group-name="`default-${u.key}`" :error="defaultErrors[u.key]" compact />
        </div>
        </section>

        <section v-if="section === 'meters' || scope !== 'market'" :class="section === 'rates' ? 'us-rates-shell us-rates-section' : 'us-card'">

        <!-- ===== อัตราค่าน้ำ-ไฟ ===== -->
        <div v-if="section === 'rates'" class="us-project-section-head">
          <h2>อัตราตามโครงการ</h2>
          <p>กำหนดค่าไฟฟ้าและค่าน้ำประปาแยกกันในแต่ละโครงการ</p>
        </div>
        <div v-if="section === 'rates'" class="us-split">
          <!-- ซ้าย: ค่ากลาง + รายการโครงการ -->
          <aside class="us-side" :style="{ height: ratesProjectPanelHeight }">
            <div class="us-side-head">
              <h3>โครงการ</h3>
              <span>{{ projects.length }} โครงการ</span>
            </div>
            <div v-if="loading" class="us-skeleton"><i v-for="n in 5" :key="n"></i></div>
            <div v-else class="us-item-list" :style="{ maxHeight: ratesProjectListHeight }">
              <button v-for="project in projects" :key="project.id" type="button" class="us-item" :class="{ active: selected === project.id }" :aria-pressed="selected === project.id" @click="select(project.id)">
                <span class="us-item-bar"></span>
                <span class="us-item-main">
                  <span class="us-item-top">
                    <span class="us-item-name">{{ project.name }}</span>
                    <span class="us-item-badge" :class="{ 'is-default': !isCustomProject(project.id) }">{{ isCustomProject(project.id) ? 'ตั้งเอง' : 'ค่ามาตรฐาน' }}</span>
                  </span>
                </span>
                <svg class="us-item-arrow" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="9 18 15 12 9 6"></polyline></svg>
              </button>
              <p v-if="!projects.length" class="us-side-empty">ไม่พบโครงการ</p>
            </div>
          </aside>

          <!-- ขวา: รายละเอียดอัตราของรายการที่เลือก -->
          <main v-if="selectedProject" ref="ratesDetailRef" class="us-detail">
            <div class="us-detail-head"><h2>อัตราค่าน้ำ-ไฟ</h2><p>{{ selectedProject.name }}</p></div>

            <div class="us-cards">
              <UtilityRateCard
                v-for="u in utilityFields" :key="`${selected}-${u.key}`"
                :rule="detailForm[u.key]" :kind="u.key" :title="u.label" :unit="u.unit" :group-name="`detail-${u.key}`" :error="detailErrors[u.key]"
                allow-inherit compact inherit-label="ค่ามาตรฐาน" :inherit-hint="describeRule(getDefaults()[u.key], u.key)"
              />
            </div>

            <!-- แผงที่ใช้อัตราต่างจากโครงการ (เฉพาะรูปแบบ โครงการ + แผง) -->
            <section v-if="scope !== 'market'" class="us-exceptions">
              <div class="us-exceptions-head">
                <div class="us-stalls-info"><strong>แผงที่ตั้งอัตราเอง</strong><small>ตั้งค่าเฉพาะแผงนี้แทนการใช้อัตราของโครงการ</small></div>
                <AppButton variant="outline" size="sm" @click="openBulk">เพิ่มแผง</AppButton>
              </div>
              <AppEmptyState v-if="!exceptionStalls.length" class="us-exceptions-empty" variant="card" icon="document" title="ยังไม่มีแผงที่ตั้งอัตราเอง" description="ทุกแผงใช้อัตราของโครงการ" />
              <div v-else class="us-table-wrap">
                <table class="us-table">
                  <thead><tr><th>แผง</th><th>ค่าไฟฟ้า</th><th>ค่าน้ำ</th><th class="us-right">จัดการ</th></tr></thead>
                  <tbody>
                    <tr v-for="stall in exceptionStalls" :key="stall.code">
                      <td><strong>แผง {{ stall.number }}</strong></td>
                      <td v-for="u in utilityFields" :key="u.key"><span class="us-chip" :class="getStallRates(stall.code)[u.key].mode !== 'inherit' ? 'is-stall' : 'is-default'">{{ getStallRates(stall.code)[u.key].mode !== 'inherit' ? ruleShortText(getStallRates(stall.code)[u.key], u.key) : 'ตามโครงการ' }}</span></td>
                      <td class="us-right us-actions">
                        <AppButton variant="outline" size="sm" @click="openStallEditor(stall)">แก้ไข</AppButton>
                        <AppButton variant="danger-outline" size="sm" @click="removeException(stall)">ลบ</AppButton>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </section>
          </main>
          <AppEmptyState v-else variant="card" icon="document" title="ยังไม่มีโครงการ" description="ตลาดนี้ยังไม่มีโครงการ" />
        </div>

        <!-- ===== มิเตอร์แผง ===== -->
        <template v-else>
          <section class="meter-overview is-compact" aria-labelledby="meter-overview-title">
            <div class="meter-overview-head">
              <div><h2 id="meter-overview-title">สรุปการตั้งค่ามิเตอร์</h2><p>ภาพรวมความพร้อมของมิเตอร์ในตลาด</p></div>
            </div>
            <div class="meter-kpis">
              <div class="meter-kpi is-all"><span>ทั้งหมด</span><strong>{{ meterMetrics.total.toLocaleString('th-TH') }}</strong><small>แผง</small></div>
              <div class="meter-kpi is-complete"><span>ตั้งค่าครบ</span><strong>{{ meterMetrics.complete.toLocaleString('th-TH') }}</strong><small>แผง</small></div>
              <div class="meter-kpi is-incomplete"><span>ยังไม่ครบ</span><strong>{{ (meterMetrics.incomplete + meterMetrics.error).toLocaleString('th-TH') }}</strong><small>แผง</small></div>
            </div>
          </section>

          <div class="meter-workspace">
            <aside class="meter-project-panel">
              <div class="meter-project-panel-head"><h3>โครงการ</h3><span>{{ visibleMeterProjects.length }}/{{ projects.length }}</span></div>
              <div class="meter-project-list">
                <button v-for="project in visibleMeterProjects" :key="project.id" type="button" class="meter-project-item" :class="{ active: meterProject === project.id }" @click="changeMeterProject(project.id)">
                  <span class="meter-project-item-top"><strong>{{ project.name }}</strong><b>{{ project.percent }}%</b></span>
                  <span class="meter-progress-track"><i :style="{ width: `${project.percent}%` }"></i></span>
                  <small>{{ project.complete.toLocaleString('th-TH') }} จาก {{ project.total.toLocaleString('th-TH') }} แผง</small>
                </button>
                <p v-if="!visibleMeterProjects.length" class="meter-progress-empty">ไม่พบโครงการ</p>
              </div>
            </aside>

            <main v-if="activeMeterProject" class="meter-workspace-detail">
              <div class="meter-manage-head">
                <div>
                  <h2>{{ activeMeterProject.name }}</h2>
                </div>
                <input v-model.trim="search" class="us-search meter-header-search" type="search" placeholder="ค้นหาแผง" aria-label="ค้นหาแผง" />
              </div>
              <AppSecondaryTabs v-model="meterFilter" :options="meterFilterOptions" aria-label="กรองสถานะการตั้งค่ามิเตอร์" />
              <div v-if="loading" class="us-skeleton"><i v-for="n in 6" :key="n"></i></div>
              <AppEmptyState v-else-if="!filteredStalls.length" class="meter-filter-empty" variant="card" icon="search" :title="meterFilter === 'complete' ? 'ยังไม่มีแผงที่ตั้งค่าครบ' : meterFilter === 'incomplete' ? 'ไม่มีแผงที่ต้องแก้ไข' : 'ไม่พบแผง'" :description="meterFilter === 'complete' ? 'ยังไม่มีแผงที่ตั้งค่ามิเตอร์ครบในโครงการนี้' : meterFilter === 'incomplete' ? 'ทุกแผงในโครงการนี้ตั้งค่ามิเตอร์ครบแล้ว' : 'ลองค้นหาหรือเปลี่ยนตัวกรอง'" />
              <template v-else>
              <div class="meter-card-list">
              <section v-for="stall in pagedStalls" :key="stall.code" class="meter-card" :data-stall="stall.code">
                <header class="meter-card-head">
                  <div><h3>แผง {{ stall.number }}</h3></div>
                  <div class="meter-card-badges">
                    <span v-if="activeRows(stall.code).length" class="meter-count-badge">{{ activeRows(stall.code).length }} มิเตอร์</span>
                    <span v-if="meterState(stall.code) === 'complete'" class="meter-state is-complete">ตั้งค่าครบ</span>
                  </div>
                </header>

                <div v-if="!meterRows(stall.code).length" class="meter-empty-note">
                  <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 8v4M12 16h.01"/></svg>
                  <div><strong>{{ meterEmptyTitle(stall) }}</strong><small>{{ meterEmptyDescription(stall) }}</small></div>
                  <AppButton v-if="requiredKinds(stall).length" variant="outline" size="sm" @click="beginMeterEditing(stall.code)">+ เพิ่มมิเตอร์</AppButton>
                </div>
                <div v-if="meterRows(stall.code).length" class="meter-columns" aria-hidden="true">
                  <span>ชื่อเรียก <b class="meter-required-mark">*</b></span><span>ประเภท</span><span>รหัสมิเตอร์ <b class="meter-required-mark">*</b></span><span>เลขอ่านล่าสุด <b class="meter-required-mark">*</b></span><span></span>
                </div>
                <div v-for="device in meterRows(stall.code)" :key="device.id" class="meter-row" :class="{ 'is-inactive': device.status === 'inactive' }">
                  <div class="meter-field"><AppInput v-model="device.name" :placeholder="meterNamePlaceholder(meterDisplayKind(stall, device.kind))" required :error="meterError(stall.code, device.id, 'name')" @update:model-value="markMeterDirty" /></div>
                  <AppSelect :model-value="meterDisplayKind(stall, device.kind)" :options="meterKindOptionsFor(stall)" :disabled="requiredKinds(stall).length === 1" @update:model-value="setMeterKind(device, $event)" />
                  <div class="meter-field"><AppInput v-model="device.code" :placeholder="meterCodePlaceholder(meterDisplayKind(stall, device.kind))" required :error="meterError(stall.code, device.id, 'code')" @update:model-value="markMeterDirty" /></div>
                  <div class="meter-field"><AppInput v-model="device.currentReading" type="number" inputmode="decimal" min="0" step="0.01" placeholder="0.00" required :error="meterError(stall.code, device.id, 'currentReading')" @update:model-value="markMeterDirty" /></div>
                  <button v-if="device.persisted" type="button" class="meter-status-toggle" :class="{ 'is-active': device.status === 'active' }" role="switch" :aria-checked="device.status === 'active'" :aria-label="device.status === 'active' ? 'ปิดใช้งานมิเตอร์' : 'เปิดใช้งานมิเตอร์'" @click="device.status === 'active' ? deactivateMeter(stall.code, device) : reactivateMeter(device)">
                    <span class="meter-status-track"><i></i></span>
                    <strong>{{ device.status === 'active' ? 'เปิด' : 'ปิดใช้งาน' }}</strong>
                  </button>
                  <AppButton v-else variant="danger-outline" size="md" full-width class="meter-action-btn" @click="deactivateMeter(stall.code, device)">ลบ</AppButton>
                </div>
                <AppButton v-if="meterRows(stall.code).length" variant="ghost" size="sm" class="meter-add" @click="addMeter(stall.code)">
                  <template #prefix><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M12 5v14M5 12h14" /></svg></template>
                  เพิ่มมิเตอร์
                </AppButton>
              </section>
              </div>
              <AppPagination v-model="page" v-model:items-per-page="perPage" :total-items="filteredStalls.length" />
              </template>
            </main>
            <AppEmptyState v-else class="meter-project-empty" variant="card" icon="document" title="เลือกโครงการเพื่อเริ่มตั้งค่ามิเตอร์" description="เลือกโครงการทางซ้ายเพื่อดูรายการแผงและจัดการมิเตอร์รายแผง">
              <template #details>
                <div class="meter-project-empty-steps">
                  <span><b>1</b> เลือกโครงการ</span>
                  <span><b>2</b> ตรวจสอบแผง</span>
                  <span><b>3</b> เพิ่มมิเตอร์ตามอัตราน้ำไฟ</span>
                </div>
              </template>
            </AppEmptyState>
          </div>
        </template>
        </section>
      </div>
    </main>

    <!-- แก้ไขอัตราของโครงการ / แผง -->
    <AppModal v-model="editorOpen" icon="settings" :title="editor ? editor.title : ''" :description="editor ? editor.description : ''">
      <div class="us-editor-grid">
        <UtilityRateCard
          v-for="u in utilityFields" :key="u.key"
          :rule="editorForm[u.key]" :kind="u.key" :title="u.label" :unit="u.unit" :group-name="`editor-${u.key}`" :error="editorErrors[u.key]"
          allow-inherit compact :inherit-label="editor?.type === 'stall' ? 'ตามโครงการ' : 'ค่ามาตรฐาน'"
          :inherit-hint="editor ? describeRule(editor.inherited[u.key], u.key) : ''"
        />
      </div>
      <template #footer>
        <AppButton variant="secondary" size="md" @click="editorOpen = false">ยกเลิก</AppButton>
        <AppButton variant="primary" size="md" @click="saveEditor">บันทึก</AppButton>
      </template>
    </AppModal>


    <!-- เพิ่มแผงที่ใช้อัตราต่างจากโครงการ (2 ขั้น: เลือกแผง -> ตั้งอัตรา) -->
    <AppModal v-model="bulkOpen" icon="settings" title="ตั้งอัตราเฉพาะแผง" :description="selectedProject?.name || ''">
      <nav class="us-steps" aria-label="ขั้นตอน">
        <span class="us-step" :class="{ 'is-active': bulkStep === 1, 'is-done': bulkStep > 1 }"><b>{{ bulkStep > 1 ? '✓' : '1' }}</b>เลือกแผง</span>
        <span class="us-step-line" :class="{ 'is-done': bulkStep > 1 }"></span>
        <span class="us-step" :class="{ 'is-active': bulkStep === 2 }"><b>2</b>ตั้งอัตรา</span>
      </nav>
      <template v-if="bulkStep === 1">
        <div class="us-pick-label">
          <span>เลือกแผง</span>
          <small>กดเลือกแผง กดซ้ำเพื่อยกเลิก</small>
        </div>
        <AppSelect v-model="bulkSelected" multiple searchable inline :options="bulkOptions" :error="bulkErrors.stalls" />
        <p class="us-picked">
          <template v-if="bulkSelected.length"><strong>เลือกแล้ว {{ bulkSelected.length }} แผง</strong> · {{ pickedText }}</template>
          <template v-else>ยังไม่ได้เลือกแผง</template>
        </p>
        <p v-if="bulkErrors.stalls" class="us-err">{{ bulkErrors.stalls }}</p>
      </template>
      <template v-else>
        <p class="us-bulk-summary">ตั้งอัตราให้ <strong>{{ bulkSelected.length }} แผง</strong></p>
        <div class="us-editor-grid">
          <UtilityRateCard
            v-for="u in utilityFields" :key="u.key"
            :rule="bulkForm[u.key]" :kind="u.key" :title="u.label" :unit="u.unit" :group-name="`bulk-${u.key}`" :error="bulkErrors[u.key]"
            allow-inherit compact inherit-label="ตามโครงการ"
            :inherit-hint="selectedProject ? describeRule(resolveProject(selectedProject.id, u.key).rule, u.key) : ''"
          />
        </div>
      </template>
      <template #footer>
        <AppButton v-if="bulkStep === 1" variant="secondary" size="md" @click="bulkOpen = false">ยกเลิก</AppButton>
        <AppButton v-else variant="secondary" size="md" @click="bulkStep = 1">ย้อนกลับ</AppButton>
        <AppButton v-if="bulkStep === 1" variant="primary" size="md" @click="goBulkStep2">ถัดไป{{ bulkSelected.length ? ` (${bulkSelected.length} แผง)` : '' }}</AppButton>
        <AppButton v-else variant="primary" size="md" @click="saveBulk">บันทึก</AppButton>
      </template>
    </AppModal>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { onBeforeRouteLeave, onBeforeRouteUpdate } from 'vue-router'
import AppButton from '~/components/AppButton.vue'
import AppEmptyState from '~/components/AppEmptyState.vue'
import AppInput from '~/components/AppInput.vue'
import AppModal from '~/components/AppModal.vue'
import AppPagination from '~/components/AppPagination.vue'
import AppSecondaryTabs from '~/components/AppSecondaryTabs.vue'
import AppSelect from '~/components/AppSelect.vue'
import { loadMarketGisJson } from '../data/market_loader'
import { getAllFloorLayoutsForProject, isOverviewBoundaryProject } from '../data/market_stall_layout'
import AppOptionCards from '~/components/AppOptionCards.vue'
import UtilityRateCard from './components/UtilityRateCard.vue'
import { RULE_MODE_LABEL, describeRule, utilitySyncError, useUtilitySettings, type MeterDevice, type RuleMode, type UtilityKind, type UtilityRates, type UtilityRule } from '../composables/useUtilitySettings'

interface ProjectItem { id: string; name: string; stallCount: number }
interface StallItem { code: string; number: string; projectId: string; projectName: string }

// แปลงและแยกเลขโครงการสำหรับจัดเรียงลำดับน้อยไปมากอย่างแม่นยำ (รองรับทั้ง 1, 2, 10, 20/1, 20/2 รวมถึงเลขไทย)
const parseProjectNumber = (name: string) => {
  const normalized = String(name || '').replace(/[๐-๙]/g, ch => String('๐๑๒๓๔๕๖๗๘๙'.indexOf(ch)))
  const match = normalized.match(/(\d+)(?:\/(\d+))?/)
  if (!match) return { main: Infinity, sub: Infinity }
  return {
    main: parseInt(match[1], 10),
    sub: match[2] ? parseInt(match[2], 10) : 0
  }
}

const compareProjects = (a: { name: string }, b: { name: string }) => {
  const numA = parseProjectNumber(a.name)
  const numB = parseProjectNumber(b.name)
  if (numA.main !== numB.main) {
    return numA.main - numB.main
  }
  if (numA.sub !== numB.sub) {
    return numA.sub - numB.sub
  }
  return a.name.localeCompare(b.name, 'th', { numeric: true })
}

const route = useRoute()
const { showSuccess, showConfirm } = useAppAlert()
const marketId = computed(() => String(route.query.market || '1'))
const { ready, getScope, setScope, getDefaults, setDefaults, getProjectRates, setProjectRates, getStallRates, setStallRates, customStallCount, resolve, resolveProject, getMeter, setMeter, duplicateOf, syncNow } = useUtilitySettings(marketId.value)

// แยกเป็น 2 หน้าตามเมนู: อัตราค่าน้ำ-ไฟ และ มิเตอร์แผง
const section = computed(() => (route.query.section === 'meters' ? 'meters' : 'rates'))

const loading = ref(true)
const projects = ref<ProjectItem[]>([])
const stalls = ref<StallItem[]>([])
const ratesDetailRef = ref<HTMLElement | null>(null)
const ratesDetailHeight = ref(0)
const ratesProjectPanelHeight = computed(() => `${Math.max(520, ratesDetailHeight.value)}px`)
const ratesProjectListHeight = computed(() => `${Math.max(180, Math.max(520, ratesDetailHeight.value) - 96)}px`)
let ratesDetailObserver: ResizeObserver | null = null
const observeRatesDetail = () => {
  ratesDetailObserver?.disconnect()
  ratesDetailObserver = null
  if (!ratesDetailRef.value || typeof ResizeObserver === 'undefined') return
  const update = () => { ratesDetailHeight.value = ratesDetailRef.value?.getBoundingClientRect().height || 0 }
  ratesDetailObserver = new ResizeObserver(update)
  ratesDetailObserver.observe(ratesDetailRef.value)
  update()
}

onMounted(async () => {
  try {
    await ready
    const cfg: any = await loadMarketGisJson(marketId.value)
    const all = cfg?.projects || []
    const visible = all.filter((p: any) => !isOverviewBoundaryProject(p, all))
    const cleanName = (name: string) => String(name || '').trim()
        const seen = new Set<string>()
    const list: StallItem[] = []
    visible.forEach((p: any) => {
      getAllFloorLayoutsForProject(p).forEach((layout: any) => {
        layout.stalls.forEach((s: any) => {
          const code = String(s?.code || s?.ident || s?.stall_ident || '').trim().toUpperCase()
          if (!code || seen.has(code)) return
          seen.add(code)
          const number = String(s.sourceProperties?.number ?? s.sourceProperties?.Number ?? '').trim() || code
          list.push({ code, number, projectId: String(p.id), projectName: cleanName(p.name) })
        })
      })
    })
    projects.value = visible
      .map((p: any) => ({ id: String(p.id), name: cleanName(p.name), stallCount: list.filter(s => s.projectId === String(p.id)).length }))
      .sort(compareProjects)
    stalls.value = list.sort((a, b) => compareProjects({ name: a.projectName }, { name: b.projectName }) || a.number.localeCompare(b.number, 'th', { numeric: true }))
    initialiseMeterDrafts()
    if (section.value === 'meters' && meterProject.value === 'all' && projects.value[0]) meterProject.value = projects.value[0].id
  } finally {
    loading.value = false
  }
})

// ----- อัตรา -----
const utilityFields = [
  { key: 'electric' as const, label: 'ค่าไฟฟ้า', unit: 'kWh' },
  { key: 'water' as const, label: 'ค่าน้ำประปา', unit: 'ลบ.ม.' }
]
const sourceLabel = (source: string) => (source === 'stall' ? 'รายแผง' : source === 'project' ? 'โครงการ' : 'ค่ามาตรฐาน')

type RuleForm = { mode: string; rate: string; minMonthly: string; flat: string; flatDaily: string; surcharge: string }
const toForm = (rule: UtilityRule): RuleForm => ({ mode: rule.mode, rate: String(rule.rate || ''), minMonthly: String(rule.minMonthly || ''), flat: String(rule.flat || ''), flatDaily: String(rule.flatDaily || ''), surcharge: String(rule.surcharge || '') })
const nonNegative = (value: string) => Math.max(0, Number(value) || 0)
const fromForm = (form: RuleForm): UtilityRule => ({ mode: form.mode as RuleMode, rate: nonNegative(form.rate), minMonthly: nonNegative(form.minMonthly), flat: nonNegative(form.flat), flatDaily: nonNegative(form.flatDaily), surcharge: nonNegative(form.surcharge) })
const validateForm = (form: RuleForm): string => {
  if (form.mode === 'meter' && !(Number(form.rate) > 0)) return 'กรุณากรอกอัตราต่อหน่วย'
  if (form.mode === 'meter' && Number(form.surcharge) < 0) return 'ค่าใช้จ่ายเพิ่มต่อหน่วยต้องไม่ติดลบ'
  if (form.mode === 'flat' && !(Number(form.flat) > 0)) return 'กรุณากรอกยอดเหมารายเดือน'
  if (form.mode === 'flat' && Number(form.flatDaily) < 0) return 'ยอดเหมารายวันต้องไม่ติดลบ'
  return ''
}
const emptyForms = () => ({ electric: toForm({ mode: 'meter', rate: 0, minMonthly: 0, flat: 0, flatDaily: 0 }), water: toForm({ mode: 'meter', rate: 0, minMonthly: 0, flat: 0, flatDaily: 0 }) })

// รูปแบบการตั้งอัตรา: ทั้งตลาดใช้แบบเดียวกัน หรือแยกตามโครงการ (พร้อมข้อยกเว้นรายแผง)
const scope = ref<string>(getScope() === 'market' ? 'market' : 'project')
const scopeOptions = [
  { value: 'market', label: 'ทั้งตลาดแบบเดียวกัน', description: 'ทุกโครงการใช้อัตราเดียวกัน' },
  { value: 'project', label: 'แยกตามโครงการ', description: 'แต่ละโครงการตั้งอัตราเอง และกำหนดข้อยกเว้นรายแผงได้' }
]
const onScopeChange = (value: string) => {
  if (value !== 'market' && projects.value[0] && !selected.value) select(projects.value[0].id)
}

// ค่ากลางของตลาด (การ์ดด้านบน)
const defaultForm = reactive(emptyForms())
const defaultErrors = reactive({ electric: '', water: '' })
const loadDefaults = () => {
  const d = getDefaults()
  defaultForm.electric = toForm(d.electric)
  defaultForm.water = toForm(d.water)
}
const snapshot = (form: Record<string, RuleForm>) => JSON.stringify([form.electric, form.water])
const defaultSnap = ref('')
loadDefaults()
defaultSnap.value = snapshot(defaultForm)
// ค่าจาก backend มาถึงทีหลังค่าที่แคชไว้ในเบราว์เซอร์: โหลดฟอร์มใหม่เมื่อพร้อม (ถ้าผู้ใช้ยังไม่ได้แก้ไขอะไร)
ready.then(() => {
  if (isDirty.value) return
  scope.value = getScope() === 'market' ? 'market' : 'project'
  loadDefaults()
  defaultSnap.value = snapshot(defaultForm)
  if (selectedProject.value) loadDetail()
})
const validateDefaults = () => {
  defaultErrors.electric = validateForm(defaultForm.electric)
  defaultErrors.water = validateForm(defaultForm.water)
  return !defaultErrors.electric && !defaultErrors.water
}

// เลือกโครงการทางซ้าย (เลือกโครงการแรกให้อัตโนมัติเมื่อโหลดเสร็จ)
const selected = ref<string>('')
const selectedProject = computed(() => projects.value.find(p => p.id === selected.value) || null)
const isCustomProject = (id: string) => {
  const r = getProjectRates(id)
  return r.electric.mode !== 'inherit' || r.water.mode !== 'inherit' || customStallCount(id) > 0
}
const customProjectCount = computed(() => projects.value.filter(p => isCustomProject(p.id)).length)
const utilityUnit = (kind: UtilityKind) => (kind === 'electric' ? 'kWh' : 'ลบ.ม.')
const ruleShortText = (rule: UtilityRule, kind: UtilityKind) => (rule.mode === 'none' ? `ไม่คิด${kind === 'electric' ? 'ค่าไฟฟ้า' : 'ค่าน้ำประปา'}` : rule.mode === 'flat' ? `เหมา ฿${rule.flat.toLocaleString('th-TH')}/เดือน` : `฿${rule.rate}/${utilityUnit(kind)}${(rule.surcharge || 0) > 0 ? ` + ฿${rule.surcharge}/${utilityUnit(kind)}` : ''}`)
const shortRule = (rule: UtilityRule, kind: UtilityKind) => (rule.mode === 'none' ? `ไม่คิด${kind === 'electric' ? 'ค่าไฟฟ้า' : 'ค่าน้ำประปา'}` : rule.mode === 'flat' ? `เหมา ฿${rule.flat}` : `฿${rule.rate}/${utilityUnit(kind)}${(rule.surcharge || 0) > 0 ? ` + ฿${rule.surcharge}/${utilityUnit(kind)}` : ''}`)
const summaryText = (rates: UtilityRates) => `ไฟ ${shortRule(rates.electric, 'electric')} · น้ำ ${shortRule(rates.water, 'water')}`
const projectEffective = (id: string): UtilityRates => ({ electric: resolveProject(id, 'electric').rule, water: resolveProject(id, 'water').rule })

const detailForm = reactive(emptyForms())
const detailErrors = reactive({ electric: '', water: '' })
const detailSnap = ref('')
// ฉบับร่างของโครงการอื่นที่แก้ไว้แล้วยังไม่บันทึก (สลับโครงการได้โดยไม่ถามและไม่หาย)
const drafts = reactive<Record<string, { electric: RuleForm; water: RuleForm }>>({})
const stashDetail = () => {
  const id = selectedProject.value?.id
  if (!id) return
  if (detailDirty.value) drafts[id] = JSON.parse(JSON.stringify({ electric: detailForm.electric, water: detailForm.water }))
  else delete drafts[id]
}
const loadDetail = () => {
  if (!selectedProject.value) return
  const id = selectedProject.value.id
  const current = getProjectRates(id)
  const saved = { electric: toForm(current.electric), water: toForm(current.water) }
  const draft = drafts[id]
  detailSnap.value = snapshot(saved as any)
  detailForm.electric = draft ? draft.electric : saved.electric
  detailForm.water = draft ? draft.water : saved.water
  delete drafts[id]
  detailErrors.electric = ''
  detailErrors.water = ''
}
const detailDirty = computed(() => !!selectedProject.value && snapshot(detailForm) !== detailSnap.value)
const normalizedScope = (value: string) => value === 'market' ? 'market' : 'project'
const isDirty = computed(() => normalizedScope(scope.value) !== normalizedScope(getScope()) || snapshot(defaultForm) !== defaultSnap.value || detailDirty.value || Object.keys(drafts).length > 0)
const saving = ref(false)
// ตรวจสอบ/ยืนยันก่อนออกจากหน้าหรือเปลี่ยน section ผ่าน query (การสลับโครงการภายในหน้านี้เก็บเป็นฉบับร่างโดยไม่ถาม)
const confirmUtilityLeave = async () => {
  if (!currentDirty.value) return true
  const subject = section.value === 'meters' ? 'ข้อมูลมิเตอร์' : 'อัตราค่าน้ำ-ไฟ'
  return await showConfirm(`${subject}มีการแก้ไขที่ยังไม่ได้บันทึก ออกจากหน้านี้และทิ้งการแก้ไขหรือไม่`, { title: 'มีการแก้ไขที่ยังไม่ได้บันทึก', confirmText: 'ทิ้งและออกจากหน้านี้' })
}
onBeforeRouteLeave(confirmUtilityLeave)
onBeforeRouteUpdate(async (to, from) => {
  const sameView = to.path === from.path && String(to.query.market || '') === String(from.query.market || '') && String(to.query.section || '') === String(from.query.section || '')
  return sameView ? true : confirmUtilityLeave()
})
const validateDetail = () => {
  const check = (form: RuleForm) => (form.mode === 'inherit' ? '' : validateForm(form))
  detailErrors.electric = check(detailForm.electric)
  detailErrors.water = check(detailForm.water)
  return !detailErrors.electric && !detailErrors.water
}
// บันทึกทุกอย่างที่แก้ไขในหน้านี้พร้อมกัน: รูปแบบ ค่ากลาง และอัตราของโครงการที่เปิดอยู่
const saveAll = () => {
  if (!isDirty.value) { showSuccess('ไม่มีการแก้ไขที่ต้องบันทึก'); return }
  if (!validateDefaults()) return
  if (scope.value !== 'market') {
    // ตรวจอัตราของทุกโครงการที่แก้ไว้ ถ้าพบที่ไม่ครบจะพาไปโครงการนั้นพร้อมแสดงข้อผิดพลาด
    stashDetail()
    const invalid = Object.entries(drafts).find(([, d]) => [d.electric, d.water].some(f => f.mode !== 'inherit' && validateForm(f)))
    if (invalid) {
      select(invalid[0])
      validateDetail()
      return
    }
    // stashDetail ย้ายโครงการที่เปิดอยู่เข้าฉบับร่างแล้ว ดึงกลับมาไว้ในฟอร์มเพื่อบันทึกรวมกัน
    loadDetail()
  }
  saving.value = true
  setScope(scope.value as 'market' | 'project')
  setDefaults({ electric: fromForm(defaultForm.electric), water: fromForm(defaultForm.water) })
  defaultSnap.value = snapshot(defaultForm)
  Object.entries(drafts).forEach(([id, d]) => setProjectRates(id, { electric: fromForm(d.electric), water: fromForm(d.water) }))
  Object.keys(drafts).forEach(id => delete drafts[id])
  if (selectedProject.value && detailDirty.value) {
    setProjectRates(selectedProject.value.id, { electric: fromForm(detailForm.electric), water: fromForm(detailForm.water) })
    detailSnap.value = snapshot(detailForm)
  }
  saving.value = false
  showSuccess('บันทึกอัตราค่าน้ำ-ไฟแล้ว')
}
const select = (id: string) => {
  if (id === selected.value) return
  stashDetail()
  selected.value = id
  loadDetail()
}
watch(projects, (list) => { const first = list[0]; if (!selected.value && first) select(first.id) })
watch([selectedProject, section, loading, scope], () => nextTick(observeRatesDetail))
watch(selected, () => {
  nextTick(() => {
    observeRatesDetail()
    if (typeof requestAnimationFrame === 'function') requestAnimationFrame(observeRatesDetail)
  })
})
onBeforeUnmount(() => ratesDetailObserver?.disconnect())

// ตัวแก้ไขอัตรา (ใช้ร่วมกันทั้งระดับโครงการและแผง)
const editorOpen = ref(false)
const editor = ref<{ type: 'project' | 'stall'; id: string; projectId: string; title: string; description: string; inherited: Record<UtilityKind, UtilityRule> } | null>(null)
const editorForm = reactive(emptyForms())
const editorErrors = reactive({ electric: '', water: '' })
const openEditor = (target: NonNullable<typeof editor.value>, current: UtilityRates) => {
  editor.value = target
  editorErrors.electric = ''
  editorErrors.water = ''
  editorForm.electric = toForm(current.electric)
  editorForm.water = toForm(current.water)
  editorOpen.value = true
}
const openStallEditor = (stall: StallItem) => openEditor({
  type: 'stall', id: stall.code, projectId: stall.projectId, title: `อัตราค่าน้ำ-ไฟ: แผง ${stall.number}`,
  description: stall.projectName,
  inherited: { electric: resolveProject(stall.projectId, 'electric').rule, water: resolveProject(stall.projectId, 'water').rule }
}, getStallRates(stall.code))
const saveEditor = () => {
  const e = editor.value
  if (!e) return
  const check = (form: RuleForm) => (form.mode === 'inherit' ? '' : validateForm(form))
  editorErrors.electric = check(editorForm.electric)
  editorErrors.water = check(editorForm.water)
  if (editorErrors.electric || editorErrors.water) return
  const value: UtilityRates = { electric: fromForm(editorForm.electric), water: fromForm(editorForm.water) }
  if (e.type === 'project') setProjectRates(e.id, value)
  else setStallRates(e.id, e.projectId, value)
  editorOpen.value = false
  showSuccess('บันทึกอัตราค่าน้ำ-ไฟแล้ว')
}

// รายแผงในโครงการ
const projectStalls = computed(() => stalls.value.filter(st => st.projectId === selected.value))
const INHERIT_RULE = (): UtilityRule => ({ mode: 'inherit', rate: 0, minMonthly: 0, flat: 0, flatDaily: 0 })
// อ่านจากข้อมูลรายแผงโดยตรง (resolve ขึ้นกับรูปแบบที่บันทึกแล้ว ซึ่งยังไม่ตรงกับที่เลือกบนหน้าจอก่อนกดบันทึก)
const hasStallOverride = (st: StallItem) => { const r = getStallRates(st.code); return r.electric.mode !== 'inherit' || r.water.mode !== 'inherit' }
const exceptionStalls = computed(() => projectStalls.value.filter(hasStallOverride))
const removeException = (st: StallItem) => {
  setStallRates(st.code, st.projectId, { electric: INHERIT_RULE(), water: INHERIT_RULE() })
  showSuccess(`แผง ${st.number} กลับไปใช้อัตราของโครงการแล้ว`)
}

// เพิ่มแผงที่ใช้อัตราต่างออกไป: เลือกหลายแผงแล้วกำหนดอัตราครั้งเดียว
const bulkOpen = ref(false)
const bulkStep = ref<1 | 2>(1)
const bulkSelected = ref<string[]>([])
const bulkForm = reactive(emptyForms())
const bulkErrors = reactive({ electric: '', water: '', stalls: '' })
const bulkCandidates = computed(() => projectStalls.value.filter(st => !hasStallOverride(st)))
const bulkOptions = computed(() => bulkCandidates.value.map(st => ({ value: st.code, label: `แผง ${st.number}` })))
// ข้อความบอกแผงที่เลือก เรียงตามเลขแผง แสดงสูงสุด 8 เลข ที่เหลือเป็นจำนวน
const pickedText = computed(() => {
  const picked = bulkCandidates.value.filter(st => bulkSelected.value.includes(st.code))
  const numbers = picked.map(st => st.number)
  if (numbers.length <= 8) return `แผง ${numbers.join(', ')}`
  return `แผง ${numbers.slice(0, 8).join(', ')} และอีก ${numbers.length - 8} แผง`
})
const goBulkStep2 = () => {
  if (!bulkSelected.value.length) { bulkErrors.stalls = 'กรุณาเลือกอย่างน้อย 1 แผง'; return }
  bulkErrors.stalls = ''
  bulkStep.value = 2
}
const openBulk = () => {
  bulkStep.value = 1
  bulkSelected.value = []
  bulkErrors.electric = bulkErrors.water = bulkErrors.stalls = ''
  const inherit = toForm(INHERIT_RULE())
  bulkForm.electric = { ...inherit }
  bulkForm.water = { ...inherit }
  bulkOpen.value = true
}
const saveBulk = () => {
  const project = selectedProject.value
  if (!project) return
  bulkErrors.stalls = bulkSelected.value.length ? '' : 'กรุณาเลือกอย่างน้อย 1 แผง'
  const noRate = bulkForm.electric.mode === 'inherit' && bulkForm.water.mode === 'inherit'
  const check = (form: RuleForm) => (form.mode === 'inherit' ? '' : validateForm(form))
  bulkErrors.electric = check(bulkForm.electric) || (noRate ? 'เลือกวิธีคิดของค่าไฟหรือค่าน้ำอย่างน้อย 1 อย่าง' : '')
  bulkErrors.water = check(bulkForm.water)
  if (bulkErrors.stalls || bulkErrors.electric || bulkErrors.water) return
  const value: UtilityRates = { electric: fromForm(bulkForm.electric), water: fromForm(bulkForm.water) }
  bulkSelected.value.forEach(code => setStallRates(code, project.id, value))
  bulkOpen.value = false
  showSuccess(`ตั้งอัตราเฉพาะให้ ${bulkSelected.value.length} แผงแล้ว`)
}

// ----- มิเตอร์ -----
type MeterDraft = Omit<MeterDevice, 'currentReading' | 'initialReading' | 'startDate'> & { currentReading: string; persisted: boolean }
type MeterField = 'name' | 'kind' | 'code' | 'currentReading'
const meterProject = ref('all')
const search = ref('')
const meterFilter = ref('all')
const page = ref(1)
const perPage = ref(10)
const meterDrafts = reactive<Record<string, MeterDraft[]>>({})
const meterErrors = reactive<Record<string, string>>({})
const meterDirty = ref(false)
const meterKindOptions = [
  { value: 'electric', label: 'ไฟฟ้า' },
  { value: 'water', label: 'น้ำประปา' }
]
const activeMeterProject = computed(() => projects.value.find(project => project.id === meterProject.value) || null)
const makeMeterDraft = (stallCode: string, kind: UtilityKind, index: number, source?: MeterDevice): MeterDraft => ({
  id: source?.id || `${stallCode}-${kind}-${Date.now()}-${index}`,
  name: source?.name || '',
  kind: source?.kind || kind,
  code: source?.code || '',
  currentReading: source?.currentReading == null ? '' : String(source.currentReading),
  status: source?.status || 'active',
  persisted: Boolean(source)
})
const requiredKinds = (stall: StallItem) => (['electric', 'water'] as UtilityKind[]).filter(kind => resolve(stall.projectId, stall.code, kind).rule.mode === 'meter')
const meterKindOptionsFor = (stall: StallItem) => {
  const allowed = requiredKinds(stall)
  return meterKindOptions.filter(option => allowed.includes(option.value as UtilityKind))
}
const meterDisplayKind = (stall: StallItem, currentKind: UtilityKind): UtilityKind => requiredKinds(stall).length === 1 ? requiredKinds(stall)[0] : currentKind
const setMeterKind = (device: MeterDraft, value: string | number) => { device.kind = String(value) as UtilityKind; markMeterDirty() }
const meterNamePlaceholder = (kind: UtilityKind) => kind === 'electric' ? 'เช่น ไฟหน้าร้าน' : 'เช่น มิเตอร์น้ำหลัก'
const meterCodePlaceholder = (kind: UtilityKind) => kind === 'electric' ? 'เช่น EM-00458' : 'เช่น WM-00201'
const meterKindLabel = (kind: UtilityKind) => kind === 'electric' ? 'ไฟฟ้า' : 'น้ำ'
const meterRequiredLabel = (stall: StallItem) => {
  const kinds = requiredKinds(stall)
  if (kinds.length === 2) return 'มิเตอร์ไฟฟ้าและมิเตอร์น้ำ'
  return kinds.length === 1 ? `มิเตอร์${meterKindLabel(kinds[0])}` : ''
}
const meterEmptyTitle = (stall: StallItem) => requiredKinds(stall).length ? `ต้องเพิ่ม${meterRequiredLabel(stall)}` : 'ไม่ต้องเพิ่มมิเตอร์'
const meterEmptyDescription = (stall: StallItem) => {
  const kinds = requiredKinds(stall)
  if (kinds.length === 2) return 'อัตราค่าไฟฟ้าและค่าน้ำของแผงนี้คิดตามมิเตอร์ กรุณาเพิ่มให้ครบทั้ง 2 ประเภท'
  if (kinds.length === 1) return `อัตราค่า${meterKindLabel(kinds[0])}ของแผงนี้คิดตามมิเตอร์ กรุณาเพิ่ม${meterRequiredLabel(stall)}`
  return 'แผงนี้ใช้เหมาจ่ายหรือไม่คิดค่าน้ำค่าไฟ จึงไม่ต้องเพิ่มมิเตอร์'
}
const initialiseMeterDrafts = () => {
  Object.keys(meterDrafts).forEach(code => delete meterDrafts[code])
  stalls.value.forEach((stall) => {
    const saved = getMeter(stall.code).devices || []
    const allowed = requiredKinds(stall)
    const rows = saved.map((device, index) => makeMeterDraft(stall.code, allowed.length === 1 ? allowed[0] : device.kind, index, device))
    if (rows.length) meterDrafts[stall.code] = rows
  })
  meterDirty.value = false
}
const meterRows = (code: string) => meterDrafts[code] || []
const editableMeterRows = (code: string) => meterDrafts[code] || (meterDrafts[code] = [])
const markMeterDirty = () => { meterDirty.value = true }
const meterKey = (stallCode: string, deviceId: string, field: MeterField) => `${stallCode}:${deviceId}:${field}`
const meterError = (stallCode: string, deviceId: string, field: MeterField) => meterErrors[meterKey(stallCode, deviceId, field)] || ''
const activeRows = (code: string) => meterRows(code).filter(device => device.status === 'active')
const meterCountsByProject = computed<Record<string, number>>(() => {
  const counts: Record<string, number> = {}
  stalls.value.forEach(stall => { counts[stall.projectId] = (counts[stall.projectId] || 0) + activeRows(stall.code).length })
  return counts
})
const rowHasInput = (row: MeterDraft) => Boolean(row.name.trim() || row.code.trim() || row.currentReading !== '')
const stallByCode = computed(() => new Map(stalls.value.map(stall => [stall.code, stall])))
const duplicateIndex = computed(() => {
  const index = new Map<string, Array<{ stall: StallItem; deviceId: string }>>()
  Object.entries(meterDrafts).forEach(([stallCode, rows]) => {
    const stall = stallByCode.value.get(stallCode)
    if (!stall) return
    rows.filter(row => row.status === 'active' && row.code.trim()).forEach((row) => {
      const code = row.code.trim().toUpperCase()
      const owners = index.get(code) || []
      owners.push({ stall, deviceId: row.id })
      index.set(code, owners)
    })
  })
  return index
})
const duplicateMeter = (stallCode: string, deviceId: string, code: string) => {
  const normalized = code.trim().toUpperCase()
  if (!normalized) return null
  return duplicateIndex.value.get(normalized)?.find(owner => owner.stall.code !== stallCode || owner.deviceId !== deviceId)?.stall || null
}
const rowIsValid = (stallCode: string, row: MeterDraft) => Boolean(
  row.status === 'inactive' || (
    row.name.trim() && row.code.trim() && row.currentReading !== '' && Number(row.currentReading) >= 0 && !duplicateMeter(stallCode, row.id, row.code)
  )
)
type MeterStateKey = 'complete' | 'incomplete' | 'error'
const computeMeterState = (code: string): MeterStateKey => {
  const stall = stallByCode.value.get(code)
  const rows = activeRows(code)
  if (rows.some(row => !rowIsValid(code, row))) {
    return rows.some(row => row.code && duplicateMeter(code, row.id, row.code)) ? 'error' : 'incomplete'
  }
  if (stall && requiredKinds(stall).some(kind => !rows.some(row => row.kind === kind))) return 'incomplete'
  return 'complete'
}
// คำนวณสถานะของทุกแผงครั้งเดียวต่อการเปลี่ยนข้อมูล (แผงหลักพัน ห้ามค้นหา/คำนวณซ้ำรายแผงในทุกจุดที่เรียกใช้)
const meterStateMap = computed(() => {
  const map = new Map<string, MeterStateKey>()
  stalls.value.forEach(stall => map.set(stall.code, computeMeterState(stall.code)))
  return map
})
const meterState = (code: string): MeterStateKey => meterStateMap.value.get(code) ?? 'incomplete'
const meterStateLabel = (code: string) => ({ complete: 'ตั้งค่าครบ', incomplete: 'ข้อมูลไม่ครบ', error: 'มีข้อมูลซ้ำ' }[meterState(code)])
const meterOverviewData = computed(() => {
  const metrics = { total: stalls.value.length, complete: 0, incomplete: 0, error: 0 }
  const grouped = new Map(projects.value.map(project => [project.id, { total: 0, complete: 0, incomplete: 0, error: 0 }]))
  stalls.value.forEach((stall) => {
    const state = meterState(stall.code)
    metrics[state] += 1
    const project = grouped.get(stall.projectId)
    if (project) {
      project.total += 1
      if (state === 'complete') project.complete += 1
      if (state === 'incomplete') project.incomplete += 1
      if (state === 'error') project.error += 1
    }
  })
  const progress = projects.value.map((project) => {
    const stats = grouped.get(project.id) || { total: 0, complete: 0, incomplete: 0, error: 0 }
    return {
      id: project.id,
      name: project.name,
      ...stats,
      percent: stats.total ? Math.round((stats.complete / stats.total) * 100) : 0
    }
  }).sort(compareProjects)
  return { metrics, progress }
})
const meterMetrics = computed(() => meterOverviewData.value.metrics)
const meterProjectProgress = computed(() => meterOverviewData.value.progress)
const visibleMeterProjects = computed(() => meterProjectProgress.value)
const discardMeterDrafts = () => initialiseMeterDrafts()
const changeMeterProject = async (value: string | number | (string | number)[]) => {
  const nextProject = String(value)
  if (nextProject === meterProject.value) return
  // การสลับโครงการยังอยู่ในหน้าเดิม จึงไม่ควรถามทิ้งการแก้ไข
  // draft ของแต่ละโครงการถูกเก็บไว้ใน meterDrafts จนกว่าจะกดบันทึกหรือออกจากหน้า
  meterProject.value = nextProject
  meterFilter.value = 'all'
  search.value = ''
  page.value = 1
}
const baseFilteredStalls = computed(() => {
  const q = search.value.toLowerCase()
  return stalls.value.filter((s) => {
    if (meterProject.value !== 'all' && s.projectId !== meterProject.value) return false
    if (!q) return true
    const terms = meterRows(s.code).map(row => `${row.name} ${row.code}`).join(' ')
    return `${s.number} ${s.code} ${s.projectName} ${terms}`.toLowerCase().includes(q)
  })
})
const filteredStalls = computed(() => baseFilteredStalls.value.filter(stall => {
  if (meterFilter.value === 'all') return true
  if (meterFilter.value === 'complete') return meterState(stall.code) === 'complete'
  if (meterFilter.value === 'incomplete') return meterState(stall.code) !== 'complete'
  return true
}))
const pagedStalls = computed(() => filteredStalls.value.slice((page.value - 1) * perPage.value, page.value * perPage.value))
const boundCount = computed(() => filteredStalls.value.filter(stall => meterState(stall.code) === 'complete').length)
const countState = (state: 'complete' | 'incomplete') => baseFilteredStalls.value.filter(stall => (state === 'complete' ? meterState(stall.code) === 'complete' : meterState(stall.code) !== 'complete')).length
const meterFilterOptions = computed(() => [
  { value: 'all', label: 'ทั้งหมด', count: baseFilteredStalls.value.length },
  { value: 'incomplete', label: 'ยังไม่ครบ', count: countState('incomplete') },
  { value: 'complete', label: 'ครบแล้ว', count: countState('complete') }
])
watch([meterProject, search, meterFilter, perPage], () => { page.value = 1 })
const addMeter = (code: string) => {
  const rows = editableMeterRows(code)
  const stall = stalls.value.find(item => item.code === code)
  const allowed = stall ? requiredKinds(stall) : []
  const kind: UtilityKind = allowed.slice().sort((a, b) => rows.filter(row => row.status === 'active' && row.kind === a).length - rows.filter(row => row.status === 'active' && row.kind === b).length)[0] || 'electric'
  rows.push(makeMeterDraft(code, kind, rows.length))
  markMeterDirty()
}
const beginMeterEditing = (code: string) => addMeter(code)
const deactivateMeter = async (code: string, device: MeterDraft) => {
  if (!device.persisted) {
    meterDrafts[code] = meterRows(code).filter(row => row.id !== device.id)
    markMeterDirty()
    return
  }
  const meterIdentity = device.code.trim() || device.name.trim()
  const meterLabel = meterIdentity ? `มิเตอร์ ${meterIdentity}` : 'มิเตอร์นี้'
  const confirmed = await showConfirm(`ต้องการปิดใช้งาน ${meterLabel} ใช่หรือไม่ มิเตอร์นี้จะไม่ถูกนำไปใช้คำนวณ แต่ข้อมูลเดิมยังคงอยู่`, { title: 'ปิดใช้งานมิเตอร์', confirmText: 'ปิดใช้งาน' })
  if (confirmed) {
    device.status = 'inactive'
    markMeterDirty()
  }
}
const reactivateMeter = (device: MeterDraft) => {
  device.status = 'active'
  markMeterDirty()
}
const validateMeters = () => {
  Object.keys(meterErrors).forEach(key => delete meterErrors[key])
  const targetStalls = meterProject.value === 'all' ? stalls.value : stalls.value.filter(stall => stall.projectId === meterProject.value)
  targetStalls.forEach((stall) => {
    activeRows(stall.code).forEach((row) => {
      const allowed = requiredKinds(stall)
      if (allowed.length === 1) row.kind = allowed[0]
      if (!row.name.trim()) meterErrors[meterKey(stall.code, row.id, 'name')] = 'กรุณากรอกชื่อมิเตอร์'
      if (!row.code.trim()) meterErrors[meterKey(stall.code, row.id, 'code')] = 'กรุณากรอกรหัสมิเตอร์'
      const duplicate = duplicateMeter(stall.code, row.id, row.code)
      if (duplicate) meterErrors[meterKey(stall.code, row.id, 'code')] = `รหัสซ้ำกับแผง ${duplicate.number}`
      if (row.currentReading === '' || Number(row.currentReading) < 0) meterErrors[meterKey(stall.code, row.id, 'currentReading')] = 'กรุณากรอกเลขมิเตอร์ปัจจุบันตั้งแต่ 0 ขึ้นไป'
    })
  })
  const firstKey = Object.keys(meterErrors)[0]
  if (firstKey) {
    const stallCode = firstKey.split(':')[0]
    const invalidStall = stalls.value.find(stall => stall.code === stallCode)
    if (invalidStall) meterProject.value = invalidStall.projectId
    meterFilter.value = 'all'
    const index = filteredStalls.value.findIndex(stall => stall.code === stallCode)
    page.value = Math.floor(Math.max(0, index) / perPage.value) + 1
    nextTick(() => document.querySelector(`[data-stall="${stallCode}"] .has-error, [data-stall="${stallCode}"] [aria-invalid="true"]`)?.scrollIntoView({ behavior: 'smooth', block: 'center' }))
  }
  return !firstKey
}
const saveMeters = async () => {
  if (!validateMeters()) return
  if (!meterDirty.value) { showSuccess('ไม่มีการแก้ไขที่ต้องบันทึก'); return }
  saving.value = true
  const targetStalls = meterProject.value === 'all' ? stalls.value : stalls.value.filter(stall => stall.projectId === meterProject.value)
  targetStalls.forEach((stall) => {
    const devices: MeterDevice[] = meterRows(stall.code)
      .filter(row => row.persisted || row.status === 'inactive' || rowHasInput(row))
      .map(({ persisted: _persisted, currentReading, ...device }) => {
        const allowed = requiredKinds(stall)
        return { ...device, kind: allowed.length === 1 ? allowed[0] : device.kind, code: device.code.trim().toUpperCase(), currentReading: currentReading === '' ? null : Number(currentReading) }
      })
    if (meterDrafts[stall.code]) setMeter(stall.code, { electric: '', water: '', devices })
  })
  const synced = await syncNow()
  if (!synced) {
    saving.value = false
    return
  }
  initialiseMeterDrafts()
  saving.value = false
  showSuccess('บันทึกการตั้งค่ามิเตอร์แล้ว')
}
const currentDirty = computed(() => section.value === 'meters' ? meterDirty.value : isDirty.value)
const saveCurrent = () => (section.value === 'meters' ? saveMeters() : saveAll())
</script>

<style scoped>
.us-heading { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 12px 20px; margin: 8px 0 18px; }
.us-heading-actions { display: inline-flex; flex-wrap: wrap; align-items: center; gap: 10px; }
.utility-sync-error { color: var(--tx-dc2626); font-size: 0.8125rem; font-weight: 600; }
.draft-status { display: inline-flex; align-items: center; gap: 7px; height: 34px; box-sizing: border-box; padding: 0 12px; border-radius: 999px; background: var(--bg-f1f5f9); color: var(--tx-475569); font-size: 0.8125rem; font-weight: 600; transition: all .2s ease; }
.draft-status i { width: 8px; height: 8px; border-radius: 50%; background: currentColor; }
.draft-status.published { background: var(--bg-ecfdf5); color: var(--tx-15803d); border: 1px solid var(--bd-bbf7d0); }
.draft-status.is-dirty { background: var(--bg-fff7ed); color: var(--tx-c2410c); border: 1px solid var(--bd-fed7aa); }
.draft-status.is-dirty i { background: #ea580c; box-shadow: 0 0 0 3px rgba(234, 88, 12, .18); }
.page-title { margin: 0; font-size: 1.5rem; font-weight: 700; color: var(--tx-0f172a); }
.page-desc { margin: 4px 0 0; font-size: 0.875rem; font-weight: 500; color: var(--tx-64748b); }
.us-card { max-width: 100%; box-sizing: border-box; padding: 24px 28px; border: 1px solid var(--bd-eef2f6); border-radius: 16px; background: var(--bg-ffffff); box-shadow: 0 1px 3px rgba(0,0,0,.03); }
.us-rates-shell { max-width: 100%; box-sizing: border-box; }
.us-rates-section { margin-top: 20px; padding-top: 20px; border-top: 1px solid var(--bd-e2e8f0); }
.us-project-section-head { margin: 0 0 14px; }
.us-project-section-head h2 { margin: 0; color: var(--tx-1e293b); font-size: 1.125rem; font-weight: 700; }
.us-project-section-head p { margin: 3px 0 0; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; line-height: 1.35; }
.us-rates-shell .us-split { align-items: start; }
.us-rates-shell .us-side { box-sizing: border-box; min-height: 0; overflow: hidden; }
.us-rates-shell .us-item-list { flex: 1 1 auto; min-height: 0; overflow-y: auto; }
.us-rates-shell .us-detail { box-sizing: border-box; min-height: 520px; padding: 24px; border: 1px solid var(--bd-e2e8f0); border-radius: 16px; background: var(--bg-ffffff); box-shadow: 0 1px 3px rgba(15, 23, 42, .03); }
.us-split { display: grid; grid-template-columns: 330px minmax(0, 1fr); gap: 30px; align-items: start; }
.us-side { display: flex; flex-direction: column; gap: 12px; padding: 20px 18px; border: 1px solid var(--bd-e2e8f0); border-radius: 16px; background: var(--bg-ffffff); box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03); min-width: 0; }
.us-side-head { padding-bottom: 12px; border-bottom: 1px solid var(--bd-f1f5f9); }
.us-side-head h3 { margin: 0; color: var(--tx-1e293b); font-size: 1.0625rem; font-weight: 700; }
.us-side-head span { display: block; margin-top: 3px; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.us-item-list { display: flex; flex-direction: column; gap: 8px; max-height: 52vh; overflow-y: auto; padding: 2px; margin: -2px; }
.us-item { display: flex; align-items: center; gap: 10px; width: 100%; padding: 12px 14px; border: 1.5px solid var(--bd-f1f5f9); border-radius: 12px; background: var(--bg-ffffff); font: inherit; text-align: left; cursor: pointer; transition: all .15s ease; }
.us-item:hover { border-color: var(--bd-cbd5e1); background: var(--bg-f8fafc); }
.us-item.active { border-color: #3b82f6; background: var(--bg-eff6ff); box-shadow: 0 2px 8px rgba(37, 99, 235, 0.1); }
.us-item-bar { flex: 0 0 auto; width: 4px; height: 28px; border-radius: 2px; background: transparent; }
.us-item.active .us-item-bar { background: #2563eb; }
.us-item-main { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 4px; }
.us-item-top { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.us-item-name { color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 700; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.us-item.active .us-item-name { color: var(--tx-1d4ed8); }
.us-item-arrow { flex: 0 0 auto; color: var(--tx-cbd5e1); transition: all .15s ease; }
.us-item.active .us-item-arrow { color: var(--tx-2563eb); transform: translateX(2px); }
.us-badge { flex: 0 0 auto; padding: 2px 10px; border-radius: 8px; background: var(--bg-eff6ff); color: var(--tx-2563eb); font-size: 0.8125rem; font-weight: 600; white-space: nowrap; }
.us-side-empty { margin: 4px 0; color: var(--tx-64748b); font-size: 0.875rem; text-align: center; }
.us-tip { display: flex; align-items: center; gap: 8px; padding: 11px 14px; border: 1px solid var(--bd-e2e8f0); border-radius: 10px; background: var(--bg-f8fafc); color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; line-height: 1.45; }
.us-tip svg { flex: 0 0 auto; color: var(--tx-3b82f6); }
.us-detail { display: flex; flex-direction: column; gap: 16px; min-width: 0; }
.us-detail-head { display: flex; flex-direction: column; gap: 2px; padding: 4px 2px 0; }
.us-detail-title-row { display: flex; align-items: center; flex-wrap: wrap; gap: 10px; }
.us-detail-head h2 { margin: 0; color: var(--tx-1e293b); font-size: 1.125rem; font-weight: 700; line-height: 1.3; }
.us-item { padding: 12px 14px; }
.us-detail-head p { margin: 0; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; line-height: 1.3; }
.us-cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(340px, 1fr)); gap: 14px; align-items: start; }
.us-detail-foot { display: flex; flex-wrap: wrap; align-items: center; justify-content: flex-end; gap: 10px; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.us-stalls { padding: 18px 20px; border: 1px solid var(--bd-e2e8f0); border-radius: 16px; background: var(--bg-ffffff); box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03); }
.us-stalls-head { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 10px; }
.us-stalls-head h3 { margin: 0; color: var(--tx-1e293b); font-size: 1.0625rem; font-weight: 700; }
.us-stalls-head span { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.us-stalls-head > div { display: flex; flex-direction: column; gap: 2px; }
@media (max-width: 992px) {
  .us-split { grid-template-columns: 1fr; gap: 18px; }
  .us-item-list { max-height: 260px; }
}
.us-default { padding: 16px 24px 18px; margin: 0 0 14px; box-sizing: border-box; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04); border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-ffffff); }
.us-h { margin: 0 0 6px; color: var(--tx-1e293b); font-size: 1.0625rem; font-weight: 700; }
.us-dot { flex: 0 0 auto; width: 8px; height: 8px; border-radius: 50%; background: #2563eb; }
.us-scope { margin-bottom: 16px; }
.us-scope .us-default-head { margin-bottom: 12px; }
.us-default-head { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 10px 16px; }
.us-default-title { display: flex; flex-direction: column; gap: 2px; }
.us-default-head h2, .us-list-head h2 { margin: 0; color: var(--tx-1e293b); font-size: 1.125rem; font-weight: 700; }
.us-default-head span, .us-list-head span { color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.us-default-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(340px, 1fr)); gap: 14px; margin-top: 12px; }
.us-editor-grid { display: flex; flex-direction: column; gap: 14px; }
.us-default-foot { display: flex; flex-wrap: wrap; align-items: center; justify-content: flex-end; gap: 10px; margin-top: 14px; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.us-list-head { display: flex; flex-wrap: wrap; align-items: flex-end; justify-content: space-between; gap: 8px; }
.us-list-head > div:first-child { display: flex; flex-direction: column; gap: 2px; }
.us-legend { display: inline-flex; align-items: center; gap: 6px; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.us-legend i { width: 10px; height: 10px; border-radius: 3px; }
.us-legend i + i, .us-legend i:not(:first-child) { margin-left: 10px; }
.us-legend .is-custom { background: #2563eb; }
.us-legend .is-default { background: var(--bg-cbd5e1); }
.us-chip { display: inline-block; padding: 4px 10px; border-radius: 8px; background: var(--bg-f1f5f9); color: var(--tx-475569); font-size: 0.8125rem; font-weight: 600; }
.us-chip.is-project, .us-chip.is-stall { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.us-actions { white-space: nowrap; }
.us-actions > * + * { margin-left: 8px; }
.us-hint { margin: 12px 0 0; padding: 10px 12px; border-radius: 10px; background: var(--bg-f8fafc); color: var(--tx-475569); font-size: 0.875rem; font-weight: 500; }
.us-table-wrap { margin-top: 16px; overflow-x: auto; -webkit-overflow-scrolling: touch; }
.us-table { width: 100%; border-collapse: collapse; }
.us-table th { padding: 14px 16px; border-bottom: 1px solid var(--bd-f1f5f9); color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; text-align: left; white-space: nowrap; }
.us-table td { padding: 14px 16px; border-bottom: 1px solid var(--bd-f1f5f9); color: var(--tx-334155); font-size: 0.875rem; font-weight: 500; vertical-align: top; }
.us-table tbody tr:hover { background: var(--bg-f8fafc); }
.us-table td strong { display: block; color: var(--tx-1e293b); font-weight: 700; }
.us-table td small { display: block; color: var(--tx-64748b); font-size: 0.8125rem; }
.us-right { text-align: right !important; }
.us-toolbar { display: flex; flex-wrap: wrap; align-items: center; gap: 10px; margin-top: 16px; }
.us-filter { min-width: 220px; }
.us-search, .us-meter { height: 40px; padding: 0 12px; box-sizing: border-box; border: 1px solid var(--bd-e2e8f0); border-radius: 8px; background: var(--bg-ffffff); color: var(--tx-334155); font: inherit; font-size: 0.9375rem; }
.us-search { flex: 1 1 220px; max-width: 320px; }
.us-meter { width: 100%; max-width: 200px; }
.us-search:focus, .us-meter:focus { outline: 2px solid #93c5fd; outline-offset: 0; border-color: #2563eb; }
.us-meter.is-dup { border-color: #fca5a5; }
.us-count { margin-left: auto; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 600; }
.us-err { color: var(--tx-dc2626) !important; font-size: 0.8125rem; font-weight: 500; }
.meter-overview { margin-bottom: 22px; }
.meter-overview.is-compact { padding-bottom: 16px; border-bottom: 1px solid var(--bd-f1f5f9); }
.meter-overview.is-compact .meter-overview-head { margin-bottom: 10px; }
.meter-overview.is-compact .meter-kpis { grid-template-columns: repeat(3, minmax(0, 1fr)); }
.meter-overview.is-compact .meter-kpi { min-height: 88px; padding: 13px 16px; }
.meter-overview-head { display: flex; align-items: flex-end; justify-content: space-between; gap: 16px; margin-bottom: 14px; }
.meter-overview-head h2 { margin: 0; color: var(--tx-1e293b); font-size: 1.125rem; font-weight: 700; }
.meter-overview-head p { margin: 3px 0 0; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.meter-overview-head > span { color: var(--tx-64748b); font-size: 0.875rem; font-weight: 600; }
.meter-overview-search { display: flex; align-items: center; gap: 12px; margin-bottom: 14px; }
.meter-overview-search .us-search { max-width: 440px; }
.meter-overview-search > span { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.meter-overview-grid { display: grid; grid-template-columns: minmax(280px, .72fr) minmax(440px, 1.65fr); gap: 16px; }
.meter-kpis { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.meter-kpi { position: relative; min-height: 112px; padding: 16px 18px; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-ffffff); color: var(--tx-1e293b); font-family: inherit; text-align: left; }
.meter-kpi span { display: block; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.meter-kpi strong { display: inline-block; margin-top: 8px; color: currentColor; font-size: 1.5rem; font-weight: 800; line-height: 1; }
.meter-kpi small { display: inline-flex; margin-left: 7px; padding: 2px 8px; border-radius: 999px; background: var(--bg-f1f5f9); color: var(--tx-64748b); font-size: 0.75rem; font-weight: 600; }
.meter-kpi.is-all { color: var(--tx-2563eb); }
.meter-kpi.is-complete { color: var(--tx-15803d); }
.meter-kpi.is-incomplete { color: var(--tx-d97706); }
.meter-kpi.is-error { color: var(--tx-dc2626); }
.meter-progress-card { min-width: 0; padding: 18px 20px; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-ffffff); }
.meter-progress-card h3 { margin: 0 0 10px; color: var(--tx-1e293b); font-size: 1rem; font-weight: 700; }
.meter-project-progress { display: flex; flex-direction: column; gap: 3px; max-height: 286px; overflow-y: auto; padding-right: 4px; }
.meter-progress-row { display: grid; grid-template-columns: minmax(110px, .7fr) minmax(180px, 1.6fr) minmax(150px, auto); align-items: center; gap: 12px; width: 100%; min-height: 42px; padding: 6px 8px; border: 1px solid transparent; border-radius: 9px; background: transparent; color: var(--tx-475569); font-family: inherit; text-align: left; cursor: pointer; }
.meter-progress-row:hover, .meter-progress-row.active { border-color: var(--bd-dbeafe); background: var(--bg-eff6ff); }
.meter-progress-name { overflow: hidden; color: var(--tx-334155); font-size: 0.8125rem; font-weight: 600; text-overflow: ellipsis; white-space: nowrap; }
.meter-progress-track { height: 10px; overflow: hidden; border-radius: 999px; background: var(--bg-e2e8f0); }
.meter-progress-track i { display: block; height: 100%; border-radius: inherit; background: #2563eb; transition: width .25s ease; }
.meter-progress-value { color: var(--tx-475569); font-size: 0.8125rem; font-weight: 500; text-align: right; white-space: nowrap; }
.meter-progress-value b { color: var(--tx-2563eb); font-weight: 700; }
.meter-progress-empty { margin: 20px 0; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; text-align: center; }
.meter-manage-head { display: flex; align-items: center; gap: 14px; padding-bottom: 16px; border-bottom: 1px solid var(--bd-f1f5f9); }
.meter-manage-head h2 { margin: 0; color: var(--tx-1e293b); font-size: 1.125rem; font-weight: 700; }
.meter-manage-head p { margin: 3px 0 0; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.meter-header-search { flex: 0 1 240px; width: 240px; height: 38px; margin-left: auto; }
.meter-workspace { display: grid; grid-template-columns: 290px minmax(0, 1fr); gap: 30px; align-items: start; }
.meter-project-panel { position: static; display: flex; flex-direction: column; gap: 12px; min-width: 0; padding: 16px; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-ffffff); }
.meter-project-panel-head { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.meter-project-panel-head h3 { margin: 0; color: var(--tx-1e293b); font-size: 1rem; font-weight: 700; }
.meter-project-panel-head span { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.meter-project-list { display: flex; flex-direction: column; gap: 8px; max-height: 520px; overflow-y: auto; padding: 2px; }
.meter-project-item { display: flex; flex-direction: column; gap: 8px; width: 100%; padding: 12px; border: 1px solid var(--bd-e2e8f0); border-radius: 11px; background: var(--bg-ffffff); color: var(--tx-334155); font-family: inherit; text-align: left; cursor: pointer; transition: border-color .15s, background .15s, box-shadow .15s; }
.meter-project-item:hover { border-color: var(--bd-cbd5e1); background: var(--bg-f8fafc); }
.meter-project-item.active { border-color: #3b82f6; background: var(--bg-eff6ff); box-shadow: 0 0 0 2px rgba(37, 99, 235, .08); }
.meter-project-item-top { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.meter-project-item-top strong { overflow: hidden; color: var(--tx-1e293b); font-size: 0.875rem; font-weight: 700; text-overflow: ellipsis; white-space: nowrap; }
.meter-project-item-top b { flex: 0 0 auto; color: var(--tx-2563eb); font-size: 0.8125rem; font-weight: 700; }
.meter-project-item small { color: var(--tx-64748b); font-size: 0.75rem; font-weight: 500; }
.meter-project-item small .is-error { color: var(--tx-dc2626); font-weight: 600; }
.meter-workspace-detail { width: 100%; min-width: 0; }
.meter-workspace-detail :deep(.secondary-tabs-scroll) { margin-top: 18px; }
.meter-card-list { display: flex; flex-direction: column; gap: 16px; width: 100%; margin-top: 18px; }
.meter-card { width: 100%; box-sizing: border-box; padding: 18px 20px; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-ffffff); box-shadow: 0 1px 3px rgba(15, 23, 42, .03); }
.meter-card-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; padding-bottom: 14px; border-bottom: 1px solid var(--bd-f1f5f9); }
.meter-card-head h3 { margin: 0; color: var(--tx-1e293b); font-size: 1.0625rem; font-weight: 700; }
.meter-card-head p { margin: 3px 0 0; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.meter-card-badges { display: inline-flex; align-items: center; gap: 8px; }
.meter-count-badge { display: inline-flex; align-items: center; min-height: 26px; padding: 3px 9px; border-radius: 999px; background: var(--bg-eff6ff); color: var(--tx-2563eb); font-size: 0.8125rem; font-weight: 600; white-space: nowrap; }
.meter-empty-note { display: flex; align-items: center; gap: 12px; margin-top: 14px; padding: 18px 20px; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: var(--bg-f8fafc); color: var(--tx-64748b); }
.meter-empty-note > svg { flex: 0 0 auto; color: var(--tx-94a3b8); }
.meter-empty-note > div { display: flex; flex: 1; flex-direction: column; gap: 2px; min-width: 0; }
.meter-empty-note strong { color: var(--tx-334155); font-size: 0.875rem; font-weight: 700; }
.meter-empty-note small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.meter-workspace-detail :deep(.meter-filter-empty.app-empty-state) { min-height: 300px; margin-top: 16px; padding: 44px 20px; box-sizing: border-box; border: 1px dashed var(--bd-d5e0ed); border-radius: 20px; background: var(--bg-ffffff); }
.meter-workspace > .meter-project-empty { min-height: 420px; box-sizing: border-box; border: 1px dashed var(--bd-d5e0ed); border-radius: 20px; background: var(--bg-f8fafc); }
.meter-project-empty-steps { display: flex; flex-wrap: wrap; justify-content: center; gap: 8px 18px; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 600; }
.meter-project-empty-steps span { display: inline-flex; align-items: center; gap: 7px; }
.meter-project-empty-steps b { display: inline-grid; place-items: center; width: 22px; height: 22px; border-radius: 50%; background: var(--bg-eff6ff); color: var(--tx-2563eb); font-size: 0.75rem; }
.meter-state { display: inline-flex; align-items: center; min-height: 28px; padding: 3px 10px; border-radius: 8px; font-size: 0.8125rem; font-weight: 600; white-space: nowrap; }
.meter-state.is-complete { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.meter-state.is-incomplete { background: var(--bg-fef3c7); color: var(--tx-d97706); }
.meter-state.is-error { background: var(--bg-fee2e2); color: var(--tx-dc2626); }
.meter-columns, .meter-row { display: grid; grid-template-columns: minmax(150px, 1.1fr) minmax(130px, .8fr) minmax(150px, 1fr) minmax(150px, .9fr) 100px; gap: 12px; align-items: start; }
.meter-columns { padding: 14px 0 7px; color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.meter-required-mark { color: var(--tx-ef4444); font-weight: 700; }
.meter-row { padding: 8px 0; border-bottom: 1px solid var(--bd-f1f5f9); }
.meter-row.is-inactive { opacity: .62; background: var(--bg-f8fafc); }
.meter-field { min-width: 0; }
.meter-action-btn { height: 42px !important; min-height: 42px !important; max-height: 42px !important; border-radius: 10px !important; font-size: 0.875rem !important; }
.meter-status-toggle { display: inline-flex; align-items: center; justify-content: center; gap: 6px; width: 100%; min-height: 42px; padding: 6px 2px; border: 0; border-radius: 10px; background: transparent; color: var(--tx-64748b); font: inherit; cursor: pointer; white-space: nowrap; }
.meter-status-toggle strong { flex: 0 0 auto; font-size: 0.8125rem; font-weight: 700; white-space: nowrap; }
.meter-status-toggle.is-active { color: var(--tx-16a34a); }
.meter-status-toggle:focus-visible { outline: 3px solid rgba(37, 99, 235, .2); outline-offset: 2px; }
.meter-status-track { position: relative; display: inline-flex; align-items: center; width: 48px; height: 28px; padding: 3px; box-sizing: border-box; border-radius: 999px; background: var(--bg-cbd5e1); transition: background .18s ease; }
.meter-status-track i { display: block; width: 22px; height: 22px; border-radius: 50%; background: var(--bg-ffffff); box-shadow: 0 1px 3px rgba(15, 23, 42, .2); transform: translateX(0); transition: transform .18s ease; }
.meter-status-toggle.is-active .meter-status-track { background: #16a34a; }
.meter-status-toggle.is-active .meter-status-track i { transform: translateX(20px); }
.meter-add { margin-top: 10px; color: var(--tx-2563eb) !important; }
.us-rule { margin-bottom: 18px; }
.us-rule h4 { margin: 0 0 8px; color: var(--tx-0f172a); font-size: 1rem; font-weight: 700; }
.us-field { display: flex; flex-direction: column; gap: 6px; margin-top: 12px; color: var(--tx-334155); font-size: 0.9375rem; font-weight: 600; }
.us-skeleton { display: flex; flex-direction: column; gap: 10px; margin-top: 16px; }
.us-skeleton i { height: 48px; border-radius: 10px; background: linear-gradient(90deg, var(--bg-f1f5f9) 25%, var(--bg-e2e8f0) 50%, var(--bg-f1f5f9) 75%); background-size: 200% 100%; animation: us-shimmer 1.4s infinite; }
.us-page-skeleton { display: flex; flex-direction: column; gap: 18px; }
.us-skeleton-title { width: 220px; height: 28px; border-radius: 8px; background: linear-gradient(90deg, var(--bg-f1f5f9) 25%, var(--bg-e2e8f0) 50%, var(--bg-f1f5f9) 75%); background-size: 200% 100%; animation: us-shimmer 1.4s infinite; }
.us-skeleton-cards { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 14px; }
.us-skeleton-cards i { height: 116px; border: 1px solid var(--bd-e2e8f0); border-radius: 14px; background: linear-gradient(90deg, var(--bg-f8fafc) 25%, var(--bg-e2e8f0) 50%, var(--bg-f8fafc) 75%); background-size: 200% 100%; animation: us-shimmer 1.4s infinite; }
.us-skeleton-workspace { display: grid; grid-template-columns: 290px minmax(0, 1fr); gap: 30px; }
.us-skeleton-sidebar, .us-skeleton-content { display: block; min-height: 420px; border: 1px solid var(--bd-e2e8f0); border-radius: 16px; background: linear-gradient(90deg, var(--bg-f8fafc) 25%, var(--bg-e2e8f0) 50%, var(--bg-f8fafc) 75%); background-size: 200% 100%; animation: us-shimmer 1.4s infinite; }
@keyframes us-shimmer { to { background-position: -200% 0; } }
@media (max-width: 640px) {
  .us-card { padding: 14px; }
  .us-rates-shell .us-detail { padding: 14px; }
  .us-search { max-width: none; flex-basis: 100%; }
  .us-count { margin-left: 0; }
  .page-title { font-size: 1.2rem; }
}
@media (max-width: 1100px) {
  .meter-overview-grid { grid-template-columns: 1fr; }
  .meter-workspace { grid-template-columns: 240px minmax(0, 1fr); gap: 18px; }
  .meter-columns { display: none; }
  .meter-row { grid-template-columns: repeat(2, minmax(0, 1fr)); padding: 14px 0; }
  .meter-row > :last-child { justify-self: start; }
}
@media (max-width: 768px) {
  .meter-overview.is-compact .meter-kpis,
  .us-skeleton-cards { grid-template-columns: 1fr; }
}
@media (max-width: 640px) {
  .meter-overview-head { align-items: flex-start; flex-direction: column; }
  .meter-overview-search { align-items: flex-start; flex-direction: column; }
  .meter-kpis { grid-template-columns: 1fr; }
  .meter-overview.is-compact .meter-kpis { grid-template-columns: 1fr; }
  .meter-kpi { min-height: 90px; padding: 14px; }
  .meter-progress-row { grid-template-columns: minmax(90px, .8fr) minmax(110px, 1.2fr); }
  .meter-progress-value { grid-column: 1 / -1; text-align: left; }
  .meter-manage-head { align-items: flex-start; flex-direction: column; }
  .meter-workspace { grid-template-columns: 1fr; }
  .meter-project-panel { position: static; }
  .meter-project-list { max-height: 360px; }
  .meter-workspace > .meter-project-empty { min-height: 300px; }
  .us-skeleton-cards { grid-template-columns: 1fr; }
  .us-skeleton-workspace { grid-template-columns: 1fr; gap: 14px; }
  .us-skeleton-sidebar { min-height: 180px; }
  .us-skeleton-content { min-height: 300px; }
  .meter-header-search { flex: 1; width: 100%; margin-left: 0; }
  .meter-card { padding: 16px; }
  .meter-row { grid-template-columns: 1fr; }
}
.us-stalls-info { display: flex; flex-direction: column; gap: 2px; }
.us-stalls-info strong { color: var(--tx-1e293b); font-size: 1.125rem; font-weight: 700; line-height: 1.3; }
.us-stalls-info small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; line-height: 1.35; }
.us-stalls-info span { color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.us-h-sub { margin-left: 8px; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 500; }
.us-exceptions { margin-top: 20px; padding-top: 16px; border-top: 1px solid var(--bd-e2e8f0); }
.us-exceptions-head { display: flex; flex-wrap: wrap; align-items: center; justify-content: space-between; gap: 10px; }
.us-bulk-summary { margin: 0 0 14px; color: var(--tx-475569); font-size: 0.9375rem; font-weight: 500; }
.us-steps { display: flex; align-items: center; justify-content: center; gap: 10px; margin: 0 0 18px; }
.us-step { display: inline-flex; align-items: center; gap: 8px; color: var(--tx-64748b); font-size: 0.875rem; font-weight: 600; }
.us-step b { display: inline-flex; align-items: center; justify-content: center; width: 24px; height: 24px; border-radius: 50%; background: var(--bg-e2e8f0); color: var(--tx-475569); font-size: 0.8125rem; font-weight: 700; }
.us-step.is-active { color: var(--tx-1d4ed8); }
.us-step.is-active b { background: #2563eb; color: #ffffff; }
.us-step.is-done b { background: var(--bg-dcfce7); color: var(--tx-15803d); }
.us-step-line { flex: 0 0 48px; height: 2px; border-radius: 2px; background: var(--bg-e2e8f0); }
.us-step-line.is-done { background: #93c5fd; }
.us-pick-label { display: flex; flex-direction: column; gap: 2px; margin-bottom: 8px; }
.us-pick-label span { color: var(--tx-1e293b); font-size: 0.9375rem; font-weight: 700; }
.us-pick-label small { color: var(--tx-64748b); font-size: 0.8125rem; font-weight: 500; }
.us-picked { margin: 12px 0 0; color: var(--tx-475569); font-size: 0.875rem; font-weight: 500; }
.us-picked strong { color: var(--tx-1e293b); font-weight: 700; }
.us-exceptions :deep(.us-exceptions-empty.app-empty-state) { min-height: 180px; margin-top: 12px; padding: 26px 20px; box-sizing: border-box; border: 1px dashed var(--bd-d5e0ed); border-radius: 20px; background: var(--bg-ffffff); }
.us-exceptions :deep(.us-exceptions-empty .app-empty-state__icon) { width: 52px; height: 52px; margin-bottom: 12px; border-radius: 15px; }
.us-exceptions :deep(.us-exceptions-empty .app-empty-state__icon svg) { width: 26px; height: 26px; }
.us-exceptions :deep(.us-exceptions-empty .app-empty-state__title) { font-size: var(--font-size-lg); }
.us-exceptions :deep(.us-exceptions-empty .app-empty-state__description) { margin-top: 6px; font-size: var(--font-size-sm); }
.us-item-badge { flex: 0 0 auto; padding: 2px 8px; border-radius: 6px; background: var(--bg-eff6ff); color: var(--tx-2563eb); font-size: 0.75rem; font-weight: 600; }
.us-item-badge.is-default { background: var(--bg-f1f5f9); color: var(--tx-64748b); }
.us-item.active .us-item-badge { background: var(--bg-ffffff); }
.us-item.active .us-item-badge.is-default { color: var(--tx-475569); }
</style>
