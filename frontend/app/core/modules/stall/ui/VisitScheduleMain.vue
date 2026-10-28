<template>
  <div class="visit-page-root">
    <AppHeader />

    <main class="content-body visit-page">
      <AppBreadcrumb
        v-if="!isStaff"
        :market="marketName"
        :marketId="marketId"
        current="จัดการรอบเข้าชมแผง"
        :backTo="`/stalls?market=${marketId}`"
      />

      <!-- Page Heading แบบรูปแรก (Pattern เดียวกันทั้งระบบ) -->
      <header class="visit-page-heading">
        <div class="heading-left">
          <h1 class="page-title">{{ isHolidaySection ? 'จัดการวันหยุด' : 'จัดการรอบเข้าชม' }}</h1>
          <p class="page-desc">{{ isHolidaySection ? 'กำหนดวันที่งดรับนัดเป็นกรณีพิเศษ' : 'ตั้งค่ารอบเวลาเข้าชมและการยื่นคำร้อง' }}</p>
        </div>
      </header>

      <section v-if="isHolidaySection" class="schedule-card card glass-card holiday-placeholder-card">
        <div class="holiday-manager">
          <div class="holiday-manager-head">
            <div><h2>วันหยุดของตลาด</h2><p>วันหยุดที่ตั้งค่าจะปิดรับนัดเฉพาะวันนั้น</p></div>
            <span class="holiday-count">{{ holidayDates.length }} วัน</span>
          </div>
          <div class="holiday-add-row">
            <AppInput v-model="holidayDateDraft" type="date" aria-label="วันที่หยุด" />
            <AppInput v-model.trim="holidayNameDraft" placeholder="ชื่อวันหยุด" aria-label="ชื่อวันหยุด" />
            <AppButton variant="primary" size="md" @click="addHolidayDate">เพิ่มวันหยุด</AppButton>
          </div>
          <div v-if="holidayDates.length" class="holiday-list">
            <div v-for="holiday in holidayDates" :key="holiday.date" class="holiday-row">
              <div><strong>{{ formatHolidayDate(holiday.date) }}</strong><span>{{ holiday.name }}</span></div>
              <button type="button" class="btn btn-outline btn-sm holiday-remove" @click="removeHolidayDate(holiday.date)">ลบ</button>
            </div>
          </div>
          <div v-else class="holiday-empty">ยังไม่มีวันหยุดที่ตั้งค่าไว้</div>
        </div>
      </section>

      <section v-else class="schedule-card card glass-card">

        <!-- ======================================================== -->
        <!-- WEEKLY MASTER (ตั้งครั้งเดียว หมุนเวียนใช้ได้ตลอดไป)         -->
        <!-- ======================================================== -->
        <div class="mode-container">

          <!-- Master-Detail 2-Column Layout for Weekly Schedule -->
          <div class="weekly-master-split">
            <!-- LEFT: Master Sidebar (7 Days List) -->
            <aside class="weekly-master-sidebar">
              <div class="sidebar-head">
                <div>
                  <h3 class="sidebar-title">วันในสัปดาห์ (7 วัน)</h3>
                  <span class="sidebar-sub">เปิดรับนัด {{ openDaysCount }} จาก 7 วัน</span>
                </div>
              </div>

              <!-- 7 Day List -->
              <div class="days-master-list" role="tablist">
                <button
                  v-for="dayObj in weeklyDayList"
                  :key="dayObj.dayOfWeek"
                  type="button"
                  class="day-master-item"
                  :class="{
                    active: selectedDayOfWeek === dayObj.dayOfWeek,
                    closed: !weeklyMaster[dayObj.dayOfWeek]?.enabled
                  }"
                  role="tab"
                  :aria-selected="selectedDayOfWeek === dayObj.dayOfWeek"
                  @click="selectedDayOfWeek = dayObj.dayOfWeek"
                >
                  <div class="dmi-indicator"></div>
                  <div class="dmi-content">
                    <div class="dmi-top-row">
                      <span class="dmi-name">{{ dayObj.dayName }}</span>
                      <span
                        v-if="weeklyMaster[dayObj.dayOfWeek]?.badge"
                        class="date-badge"
                        :class="weeklyMaster[dayObj.dayOfWeek]?.badgeType"
                      >
                        {{ weeklyMaster[dayObj.dayOfWeek]?.badge }}
                      </span>
                    </div>

                    <div class="dmi-bot-row">
                      <span v-if="weeklyMaster[dayObj.dayOfWeek]?.enabled" class="dmi-status-open">
                        เปิด {{ openWeeklySlotCount(dayObj.dayOfWeek) }} รอบ · {{ getWeeklyTotalCapacity(dayObj.dayOfWeek) }} คิว
                      </span>
                      <span v-else class="dmi-status-closed">
                        ปิดรับนัด
                      </span>
                    </div>
                  </div>

                  <span class="dmi-arrow">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="9 18 15 12 9 6"></polyline>
                    </svg>
                  </span>
                </button>
              </div>

              <div class="sidebar-tip-box">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="16" x2="12" y2="12"></line>
                  <line x1="12" y1="8" x2="12.01" y2="8"></line>
                </svg>
                <span>ตารางนี้ใช้ซ้ำทุกสัปดาห์</span>
              </div>
            </aside>

            <!-- RIGHT: Detail Editor Panel -->
            <main class="weekly-detail-panel">
              <!-- Active Day Header -->
              <div class="detail-day-header">
                <div class="ddh-left">
                  <div class="ddh-title-row">
                    <h2 class="ddh-day-title">กำหนดตาราง{{ activeWeeklyDay.dayName }}</h2>
                    <span
                      v-if="activeWeeklyDay.badge"
                      class="date-badge"
                      :class="activeWeeklyDay.badgeType"
                    >
                      {{ activeWeeklyDay.badge }}
                    </span>
                  </div>
                  <p class="ddh-subtitle">
                    รอบเวลานี้จะเปิดให้ผู้เช่าเลือกจองในทุกๆ {{ activeWeeklyDay.dayName }} ของทุกสัปดาห์
                  </p>
                </div>

                <!-- Toggle Open/Closed for this day -->
                <div class="ddh-right">
                  <label class="day-toggle-switch">
                    <input
                      v-model="activeWeeklyDay.enabled"
                      type="checkbox"
                      @change="saveWeeklyMaster"
                    />
                    <span class="slider"></span>
                  </label>
                  <span class="toggle-text" :class="{ 'is-active': activeWeeklyDay.enabled }">
                    {{ activeWeeklyDay.enabled ? `เปิดรับนัด` : `ปิดรับนัด` }}
                  </span>
                </div>
              </div>

              <!-- Active Day Config Box (Badge & Advice) -->
              <div v-if="activeWeeklyDay.enabled" class="detail-day-config">
                <div class="config-field badge-field">
                  <label class="field-label">ป้ายกำกับประจำวัน (Badge)</label>
                  <AppSelect
                    v-model="activeWeeklyDay.badge"
                    :options="badgeOptions"
                    placeholder="เลือกป้ายกำกับ"
                    @change="handleWeeklyBadgeChange"
                  />
                </div>

                <div class="config-field advice-field">
                  <label class="field-label">คำแนะนำสำหรับผู้เช่าใน{{ activeWeeklyDay.dayName }}</label>
                  <input
                    v-model.trim="activeWeeklyDay.description"
                    type="text"
                    class="styled-input"
                    placeholder="เช่น เปิดเต็มวัน คนเดินหนาแน่น เหมาะดูทำเลและการค้าขายจริง"
                    @change="saveWeeklyMaster"
                  />
                </div>
              </div>

              <!-- Slots Management Box -->
              <div v-if="activeWeeklyDay.enabled" class="detail-slots-card">
                <div class="slots-card-header">
                  <div class="sch-left">
                    <h3 class="sch-title">รอบเวลาเข้าชม (Time Slots)</h3>
                    <span class="sch-sub">{{ activeWeeklyDay.slots?.length || 0 }} รอบ | รองรับรวม {{ getWeeklyTotalCapacity(activeWeeklyDay.dayOfWeek) }} คิว</span>
                  </div>

                  <button
                    type="button"
                    class="btn-add-slot-top"
                    @click="addNewWeeklySlot"
                  >
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <line x1="12" y1="5" x2="12" y2="19"></line>
                      <line x1="5" y1="12" x2="19" y2="12"></line>
                    </svg>
                    <span>เพิ่มรอบเวลา</span>
                  </button>
                </div>

                <!-- Slot List Table -->
                <div class="slot-list">
                  <div class="slot-list-head">
                    <span>ช่วงเวลา (เริ่ม - สิ้นสุด)</span>
                    <span>โควตารับผู้เช่า</span>
                    <span>คำอธิบาย / จุดเด่นของรอบ</span>
                    <span style="text-align: center;">สถานะ</span>
                    <span></span>
                  </div>

                  <article
                    v-for="(slot, index) in activeWeeklyDay.slots"
                    :key="slot.id || `slot-${index}`"
                    class="slot-row"
                    :class="{
                      disabled: !slot.enabled
                    }"
                  >
                    <!-- Time Range Pickers (Start - End) -->
                    <div class="slot-col-time">
                      <div class="time-range-picker-group">
                        <AppTimePicker
                          v-model="slot.startTime"
                          placeholder="เริ่ม"
                          size="sm"
                          :aria-label="`เวลาเริ่มรอบ ${index + 1}`"
                          @change="handleWeeklySlotTimeChange(slot)"
                        />
                        <span class="time-range-separator">-</span>
                        <AppTimePicker
                          v-model="slot.endTime"
                          placeholder="สิ้นสุด"
                          size="sm"
                          :aria-label="`เวลาสิ้นสุดรอบ ${index + 1}`"
                          @change="handleWeeklySlotTimeChange(slot)"
                        />
                      </div>
                    </div>

                    <!-- Capacity Numeric Stepper Input -->
                    <div class="slot-col-cap">
                      <div class="cap-stepper-box">
                        <button
                          type="button"
                          class="cap-step-btn"
                          :disabled="slot.capacity <= 1 || !slot.enabled"
                          title="ลดจำนวนคิว"
                          @click="decrementSlotCapacity(slot)"
                        >
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <line x1="5" y1="12" x2="19" y2="12"></line>
                          </svg>
                        </button>
                        <div class="cap-input-wrap">
                          <input
                            v-model.number="slot.capacity"
                            type="number"
                            min="1"
                            max="99"
                            class="cap-number-input"
                            :aria-label="`โควตารับผู้เช่ารอบที่ ${index + 1}`"
                            @input="handleWeeklyCapacityChange(slot)"
                            @blur="normalizeSlotCapacity(slot)"
                          />
                          <span class="cap-suffix">คิว</span>
                        </div>
                        <button
                          type="button"
                          class="cap-step-btn"
                          :disabled="slot.capacity >= 99 || !slot.enabled"
                          title="เพิ่มจำนวนคิว"
                          @click="incrementSlotCapacity(slot)"
                        >
                          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <line x1="12" y1="5" x2="12" y2="19"></line>
                            <line x1="5" y1="12" x2="19" y2="12"></line>
                          </svg>
                        </button>
                      </div>
                    </div>

                    <!-- Note Input -->
                    <div class="slot-col-note">
                      <input
                        v-model.trim="slot.note"
                        type="text"
                        class="note-input"
                        aria-label="คำอธิบายรอบ"
                        placeholder="เช่น ช่วงเช้า อากาศสบาย เหมาะดูโครงสร้าง"
                        @change="saveWeeklyMaster"
                      />
                    </div>

                    <!-- Status Toggle Switch (เปิด / ปิด) -->
                    <div class="slot-col-status">
                      <label
                        class="status-toggle-control"
                        :title="slot.enabled ? 'คลิกเพื่อปิดรอบเวลานี้' : 'คลิกเพื่อเปิดรอบเวลานี้'"
                      >
                        <input
                          v-model="slot.enabled"
                          type="checkbox"
                          class="status-toggle-input"
                          @change="handleWeeklySlotToggle(slot)"
                        />
                        <span class="status-toggle-slider"></span>
                        <span
                          class="status-toggle-label"
                          :class="{ 'is-active': slot.enabled }"
                        >
                          {{ slot.enabled ? 'เปิด' : 'ปิด' }}
                        </span>
                      </label>
                    </div>

                    <!-- Delete Action -->
                    <div class="slot-col-actions">
                      <button
                        class="delete-button"
                        type="button"
                        title="ลบรอบเวลานี้"
                        @click="removeWeeklySlot(index)"
                      >
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <polyline points="3 6 5 6 21 6"></polyline>
                          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                        </svg>
                      </button>
                    </div>
                  </article>

                </div>
              </div>

              <!-- Closed Day Note -->
              <div v-else class="closed-day-banner">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"></line>
                </svg>
                <div class="cdb-text">
                  <strong>ปิดรับนัดหมายทุก{{ activeWeeklyDay.dayName }}</strong>
                  <p>ในหน้าจองของผู้เช่า {{ activeWeeklyDay.dayName }}ของทุกสัปดาห์จะแสดงสถานะเป็น “ปิดรับนัดหมาย”</p>
                </div>
                <button type="button" class="reopen-day-btn" @click="activeWeeklyDay.enabled = true">
                  เปิดรับนัดทุก{{ activeWeeklyDay.dayName }}
                </button>
              </div>
            </main>
          </div>
        </div>


        <!-- Footer Bar with Save -->
        <footer class="schedule-footer">
          <div class="footer-right">
            <AppButton variant="primary" size="md" @click="handleSaveAll">
              <template #prefix>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2Z"></path>
                  <polyline points="17 21 17 13 7 13 7 21"></polyline>
                  <polyline points="7 3 7 8 15 8"></polyline>
                </svg>
              </template>
              บันทึกตารางเวลา
            </AppButton>
          </div>
        </footer>
      </section>
    </main>

    <!-- Copy Modal: คัดลอกการตั้งค่าไปยังวันอื่นๆ ในสัปดาห์ -->
    <div v-if="showWeeklyCopyModal" class="copy-modal-overlay" @click.self="showWeeklyCopyModal = false">
      <div class="copy-modal-card">
        <div class="cm-head">
          <div class="cm-title-wrap">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" style="color: #2563eb;">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
            <h3>คัดลอกรอบเวลาของ{{ activeWeeklyDay.dayName }}</h3>
          </div>
          <button type="button" class="cm-close" @click="showWeeklyCopyModal = false">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>
        <div class="cm-body">
          <p class="cm-desc">เลือกวันที่ต้องการใช้รอบเวลา {{ activeWeeklyDay.slots.length }} รอบ</p>

          <div class="copy-target-options">
            <label class="target-option-card">
              <input v-model="weeklyCopyTargetType" type="radio" value="weekends" />
              <div>
                <strong>ทุกวันเสาร์และอาทิตย์ (สุดสัปดาห์)</strong>
                <p>ใช้กับวันเสาร์และอาทิตย์</p>
              </div>
            </label>
            <label class="target-option-card">
              <input v-model="weeklyCopyTargetType" type="radio" value="all_week" />
              <div>
                <strong>ทุกวันในสัปดาห์ (จันทร์ - อาทิตย์)</strong>
                <p>ใช้กับทุกวัน</p>
              </div>
            </label>
            <label class="target-option-card">
              <input v-model="weeklyCopyTargetType" type="radio" value="weekdays" />
              <div>
                <strong>วันจันทร์ - ศุกร์ (วันทำการ)</strong>
                <p>นำรอบเวลานี้ไปใช้กับทุกวันธรรมดา</p>
              </div>
            </label>
          </div>
        </div>
        <div class="cm-foot">
          <button type="button" class="outline-btn" @click="showWeeklyCopyModal = false">ยกเลิก</button>
          <button type="button" class="primary-btn" @click="executeWeeklyCopy">
            <span>ยืนยันคัดลอกรอบเวลา</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppSelect from '~/components/AppSelect.vue'
import AppCheckbox from '~/components/AppCheckbox.vue'
import AppTimePicker from '~/components/AppTimePicker.vue'
import {
  type IVisitSlot,
  type IWeeklyDayMaster,
  type IWeeklyMasterMap,
  getStoredWeeklyMaster,
  saveStoredWeeklyMaster,
  createDefaultWeeklyMaster,
  createDefaultSlotList,
  parseSlotTimeRange,
  formatSlotTime,
  setCachedScheduleBookings,
  isBookingOccupyingVisitSlot,
  doesBookingMatchMarket,
  doesBookingMatchSlotTime,
  normalizeVisitDateToIso,
  getAppointmentRoundKey
} from '../composables/useMarketVisitSchedule'
import { getMarketGisConfig } from '../data/markets_registry'

const route = useRoute()
const { get, post } = useApi()
const { user } = useAuth()
const { showSuccess, showError } = useAppAlert()
const { applyMarketOverride } = useMarketOverrides()
const isStaff = computed(() => user.value?.role === 'staff')
const isHolidaySection = computed(() => route.query.section === 'holidays')
const holidayDateDraft = ref('')
const holidayNameDraft = ref('')
const holidayDates = ref<Array<{ date: string; name: string }>>([])

const marketId = computed(() => String(route.query.market || '1'))

const holidayStorageKey = computed(() => `emarket-holidays-${marketId.value}`)
const loadHolidayDates = () => {
  if (!import.meta.client) return
  try { holidayDates.value = JSON.parse(localStorage.getItem(holidayStorageKey.value) || '[]') } catch { holidayDates.value = [] }
}
const persistHolidayDates = () => {
  if (import.meta.client) localStorage.setItem(holidayStorageKey.value, JSON.stringify(holidayDates.value))
}
const addHolidayDate = () => {
  if (!holidayDateDraft.value) return
  if (!holidayDates.value.some(item => item.date === holidayDateDraft.value)) {
    holidayDates.value.push({ date: holidayDateDraft.value, name: holidayNameDraft.value || 'วันหยุดนักขัตฤกษ์' })
    holidayDates.value.sort((a, b) => a.date.localeCompare(b.date))
    persistHolidayDates()
  }
  holidayDateDraft.value = ''
  holidayNameDraft.value = ''
}
const removeHolidayDate = (date: string) => { holidayDates.value = holidayDates.value.filter(item => item.date !== date); persistHolidayDates() }
const formatHolidayDate = (date: string) => new Intl.DateTimeFormat('th-TH', { dateStyle: 'long' }).format(new Date(`${date}T00:00:00`))
onMounted(loadHolidayDates)
watch(marketId, loadHolidayDates)
const apiMarketName = ref('')
const marketName = computed(() => {
  const cfg = getMarketGisConfig(marketId.value)
  const overridden = applyMarketOverride({ id: cfg.id, slug: cfg.slug, name: apiMarketName.value || cfg.name })
  return overridden?.name || cfg.name
})
const scheduleBookings = ref<any[]>([])

const getDayBookedCount = (dayOfWeek: number): number => {
  const matched = scheduleBookings.value.filter(b => {
    if (!isBookingOccupyingVisitSlot(b)) return false
    if (!doesBookingMatchMarket(b, marketId.value, marketName.value)) return false
    const iso = normalizeVisitDateToIso(b.visit_date)
    if (!iso) return false
    const d = new Date(iso)
    return !isNaN(d.getTime()) && d.getDay() === dayOfWeek
  })
  return matched.length
}

const showPresetMenu = ref(false)
const showWeeklyCopyModal = ref(false)
const weeklyCopyTargetType = ref<'weekends' | 'all_week' | 'weekdays'>('weekends')

// 7 Weekdays Master (0 = Sun, 1 = Mon, ..., 6 = Sat)
const weeklyDayList = [
  { dayOfWeek: 1, dayName: 'วันจันทร์' },
  { dayOfWeek: 2, dayName: 'วันอังคาร' },
  { dayOfWeek: 3, dayName: 'วันพุธ' },
  { dayOfWeek: 4, dayName: 'วันพฤหัสบดี' },
  { dayOfWeek: 5, dayName: 'วันศุกร์' },
  { dayOfWeek: 6, dayName: 'วันเสาร์' },
  { dayOfWeek: 0, dayName: 'วันอาทิตย์' }
]

const weeklyMaster = reactive<IWeeklyMasterMap>(createDefaultWeeklyMaster())
const selectedDayOfWeek = ref<number>(6) // Default to Saturday (Prime market day)

const openDaysCount = computed(() => {
  return weeklyDayList.filter(d => weeklyMaster[d.dayOfWeek]?.enabled).length
})

const activeWeeklyDay = computed<IWeeklyDayMaster>(() => {
  return weeklyMaster[selectedDayOfWeek.value] || weeklyMaster[6]
})

const badgeOptions = [
  { value: '', label: 'ไม่มีป้ายกำกับ' },
  { value: 'แนะนำ', label: 'แนะนำ (คนเดินเยอะ)' },
  { value: 'เปิดเต็มวัน', label: 'เปิดเต็มวัน' },
  { value: 'ตลาดกลางคืน', label: 'ตลาดกลางคืน' },
  { value: 'วันทำการ', label: 'วันทำการ (ดูโครงสร้าง)' },
  { value: 'รอบพิเศษ', label: 'รอบพิเศษ' }
]

const capacityOptions = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map(val => ({
  value: val,
  label: `รับได้ ${val} คิว`
}))

const initAllData = () => {
  // Load Weekly Master
  const storedWeekly = getStoredWeeklyMaster(marketId.value)
  Object.keys(weeklyMaster).forEach(k => delete (weeklyMaster as any)[k])
  Object.assign(weeklyMaster, storedWeekly)

  // Ensure all slots have startTime and endTime
  Object.values(weeklyMaster).forEach(day => {
    day.slots?.forEach(slot => ensureSlotTimeRange(slot))
  })
}

const openWeeklySlotCount = (dayOfWeek: number): number => {
  const day = weeklyMaster[dayOfWeek]
  if (!day || !day.enabled) return 0
  return (day.slots || []).filter(s => s.enabled).length
}

const getWeeklyTotalCapacity = (dayOfWeek: number): number => {
  const day = weeklyMaster[dayOfWeek]
  if (!day || !day.enabled) return 0
  return (day.slots || []).filter(s => s.enabled).reduce((sum, s) => sum + (s.capacity || 0), 0)
}

const handleWeeklyBadgeChange = (val: string) => {
  if (val === 'แนะนำ') activeWeeklyDay.value.badgeType = 'badge-primary'
  else if (val === 'เปิดเต็มวัน') activeWeeklyDay.value.badgeType = 'badge-success'
  else if (val === 'ตลาดกลางคืน') activeWeeklyDay.value.badgeType = 'badge-purple'
  else activeWeeklyDay.value.badgeType = 'badge-neutral'
  saveWeeklyMaster()
}

const saveWeeklyMaster = () => {
  saveStoredWeeklyMaster(marketId.value, weeklyMaster)
}

const handleWeeklySlotToggle = (slot: IVisitSlot) => {
  if (!slot.enabled) {
    slot.isFull = true
    slot.remainingQueues = 0
  } else {
    slot.isFull = false
    slot.remainingQueues = Number(slot.capacity) || 0
  }
  saveWeeklyMaster()
}

const handleWeeklyCapacityChange = (slot: IVisitSlot) => {
  slot.remainingQueues = Number(slot.capacity) || 0
  if (slot.capacity > 0 && slot.enabled) slot.isFull = false
  saveWeeklyMaster()
}

const normalizeSlotCapacity = (slot: IVisitSlot) => {
  const num = parseInt(String(slot.capacity), 10)
  if (isNaN(num) || num < 1) {
    slot.capacity = 1
  } else if (num > 99) {
    slot.capacity = 99
  } else {
    slot.capacity = num
  }
  handleWeeklyCapacityChange(slot)
}

const incrementSlotCapacity = (slot: IVisitSlot) => {
  const current = parseInt(String(slot.capacity), 10) || 1
  slot.capacity = Math.min(current + 1, 99)
  handleWeeklyCapacityChange(slot)
}

const decrementSlotCapacity = (slot: IVisitSlot) => {
  const current = parseInt(String(slot.capacity), 10) || 1
  slot.capacity = Math.max(current - 1, 1)
  handleWeeklyCapacityChange(slot)
}

const toggleWeeklySlotFull = (slot: IVisitSlot) => {
  slot.isFull = !slot.isFull
  slot.remainingQueues = slot.isFull ? 0 : slot.capacity
}

const ensureSlotTimeRange = (slot: IVisitSlot) => {
  if (!slot.id) {
    slot.id = Date.now() + Math.floor(Math.random() * 1000)
  }
  if (!slot.startTime || !slot.endTime) {
    const parsed = parseSlotTimeRange(slot.time)
    if (!slot.startTime) slot.startTime = parsed.startTime
    if (!slot.endTime) slot.endTime = parsed.endTime
  }
  if (!slot.time) {
    slot.time = formatSlotTime(slot.startTime, slot.endTime)
  }
}

const handleWeeklySlotTimeChange = (slot: IVisitSlot) => {
  slot.time = formatSlotTime(slot.startTime, slot.endTime)
  saveWeeklyMaster()
}

watch(
  selectedDayOfWeek,
  () => {
    activeWeeklyDay.value?.slots?.forEach(slot => ensureSlotTimeRange(slot))
  },
  { immediate: true }
)

const addNewWeeklySlot = () => {
  const slots = activeWeeklyDay.value.slots
  const defaultPresets = [
    { start: '09:30', end: '10:30' },
    { start: '11:00', end: '12:00' },
    { start: '13:00', end: '14:00' },
    { start: '15:30', end: '16:30' },
    { start: '17:00', end: '18:00' },
    { start: '18:30', end: '19:30' }
  ]
  const nextPreset = defaultPresets[slots.length % defaultPresets.length] || { start: '17:00', end: '18:00' }

  slots.push({
    id: Date.now() + Math.floor(Math.random() * 1000),
    startTime: nextPreset.start,
    endTime: nextPreset.end,
    time: formatSlotTime(nextPreset.start, nextPreset.end),
    capacity: 3,
    enabled: true,
    note: 'รอบพิเศษ มีเจ้าหน้าที่พาชม',
    remainingQueues: 3,
    isFull: false
  })
  saveWeeklyMaster()
}

const removeWeeklySlot = (index: number) => {
  if (!activeWeeklyDay.value?.slots) return
  activeWeeklyDay.value.slots.splice(index, 1)
  saveWeeklyMaster()
}

const openWeeklyCopyModal = () => {
  showWeeklyCopyModal.value = true
}

const executeWeeklyCopy = () => {
  const sourceSlots = JSON.parse(JSON.stringify(activeWeeklyDay.value.slots))
  const sourceBadge = activeWeeklyDay.value.badge
  const sourceBadgeType = activeWeeklyDay.value.badgeType
  const sourceDesc = activeWeeklyDay.value.description

  let count = 0
  weeklyDayList.forEach(dayObj => {
    if (dayObj.dayOfWeek === selectedDayOfWeek.value) return

    let shouldCopy = false
    if (weeklyCopyTargetType.value === 'weekends' && (dayObj.dayOfWeek === 0 || dayObj.dayOfWeek === 6)) shouldCopy = true
    else if (weeklyCopyTargetType.value === 'weekdays' && dayObj.dayOfWeek >= 1 && dayObj.dayOfWeek <= 5) shouldCopy = true
    else if (weeklyCopyTargetType.value === 'all_week') shouldCopy = true

    if (shouldCopy) {
      weeklyMaster[dayObj.dayOfWeek] = {
        dayOfWeek: dayObj.dayOfWeek,
        dayName: dayObj.dayName,
        enabled: true,
        badge: sourceBadge,
        badgeType: sourceBadgeType,
        description: sourceDesc,
        slots: JSON.parse(JSON.stringify(sourceSlots))
      }
      count++
    }
  })

  showWeeklyCopyModal.value = false
  showSuccess(
    `คัดลอกตารางรอบเวลาไปยัง ${count} วัน เรียบร้อยแล้ว (กด "บันทึกตารางเวลา" เพื่อยืนยันใช้งานตลอดไป)`,
    'คัดลอกสำเร็จ'
  )
}

const applyWeeklyPreset = (type: 'weekend' | 'everyday' | 'weekdays' | 'reset') => {
  showPresetMenu.value = false

  if (type === 'reset') {
    const def = createDefaultWeeklyMaster()
    Object.keys(weeklyMaster).forEach(k => delete (weeklyMaster as any)[k])
    Object.assign(weeklyMaster, def)
    refreshCalendarProjection()
    showSuccess('คืนค่าตารางประจำสัปดาห์เริ่มต้นเรียบร้อยแล้ว', 'รีเซ็ตสำเร็จ')
    return
  }

  weeklyDayList.forEach(dayObj => {
    const dow = dayObj.dayOfWeek
    if (type === 'weekend') {
      const isWeekend = dow === 6 || dow === 0 || dow === 5
      weeklyMaster[dow] = {
        dayOfWeek: dow,
        dayName: dayObj.dayName,
        enabled: isWeekend,
        badge: dow === 6 ? 'แนะนำ' : dow === 0 ? 'เปิดเต็มวัน' : dow === 5 ? 'ตลาดกลางคืน' : 'ปิดรับนัด',
        badgeType: dow === 6 ? 'badge-primary' : dow === 0 ? 'badge-success' : dow === 5 ? 'badge-purple' : 'badge-neutral',
        description: isWeekend ? 'เปิดตลาดรอบสุดสัปดาห์ เหมาะดูทำเลและการค้าขาย' : 'วันปิดปรับปรุงตลาด',
        slots: isWeekend ? createDefaultSlotList() : createDefaultSlotList().map(s => ({ ...s, enabled: false }))
      }
    } else if (type === 'everyday') {
      weeklyMaster[dow] = {
        dayOfWeek: dow,
        dayName: dayObj.dayName,
        enabled: true,
        badge: dow === 6 || dow === 0 ? 'แนะนำ' : 'วันทำการ',
        badgeType: dow === 6 || dow === 0 ? 'badge-primary' : 'badge-neutral',
        description: 'เปิดรับเข้าชมทุกวัน มีเจ้าหน้าที่ฝ่ายอาคารพาเดินชม',
        slots: createDefaultSlotList()
      }
    } else if (type === 'weekdays') {
      const isWeekday = dow >= 1 && dow <= 5
      weeklyMaster[dow] = {
        dayOfWeek: dow,
        dayName: dayObj.dayName,
        enabled: isWeekday,
        badge: isWeekday ? 'วันทำการ' : 'ปิดรับนัด',
        badgeType: isWeekday ? 'badge-neutral' : 'badge-neutral',
        description: isWeekday ? 'รอบเวลาพาชมในวันทำการราชการ' : 'วันหยุดประจำสัปดาห์',
        slots: isWeekday ? createDefaultSlotList() : createDefaultSlotList().map(s => ({ ...s, enabled: false }))
      }
    }
  })

  showSuccess('ปรับใช้แม่แบบสำเร็จ กด "บันทึกตารางเวลา" เพื่อให้มีผลถาวร', 'ปรับใช้แม่แบบแล้ว')
}

const handleSaveAll = () => {
  saveStoredWeeklyMaster(marketId.value, weeklyMaster)
  showSuccess(
    'บันทึกตารางหลักประจำสัปดาห์สำเร็จแล้ว! ระบบจะนำตารางเวลานี้ไปใช้งานวนรอบทุกสัปดาห์ตลอดไปโดยอัตโนมัติ โดยไม่ต้องกลับมาเซ็ตซ้ำ',
    'บันทึกตารางเวลาแล้ว'
  )
}

const loadMarketAndBookings = async () => {
  initAllData()
  const cfg = getMarketGisConfig(marketId.value)
  apiMarketName.value = cfg.name

  try {
    const [marketsRes, bookingsRes] = await Promise.allSettled([
      get('/markets'),
      get('/bookings')
    ])
    if (marketsRes.status === 'fulfilled') {
      const markets = Array.isArray(marketsRes.value) ? marketsRes.value : (marketsRes.value?.data || [])
      const market = markets.find((m: any) => String(m.id) === marketId.value || m.slug === marketId.value || String(m.id) === String(cfg.id))
      if (market) {
        const resolved = applyMarketOverride(market)
        apiMarketName.value = resolved.name || cfg.name
      }
    }
    if (bookingsRes.status === 'fulfilled') {
      const bList = Array.isArray(bookingsRes.value)
        ? bookingsRes.value
        : (Array.isArray((bookingsRes.value as any)?.data) ? (bookingsRes.value as any).data : [])
      scheduleBookings.value = bList
      setCachedScheduleBookings(bList)
    }
  } catch (e) {
    console.warn('Could not fetch market or bookings:', e)
  }
}

const handleDocumentClick = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest?.('.preset-dropdown-wrap')) {
    showPresetMenu.value = false
  }
}

onMounted(() => {
  void loadMarketAndBookings()
  window.addEventListener('click', handleDocumentClick)
})

onBeforeUnmount(() => {
  window.removeEventListener('click', handleDocumentClick)
})

watch(marketId, () => {
  void loadMarketAndBookings()
})
</script>

<style scoped>
.visit-page-root {
  min-height: 100vh;
  background: #f8fafc;
}

.visit-page {
  width: 100%;
  max-width: 100%;
  padding: 24px 32px 48px;
  box-sizing: border-box;
}

.schedule-card {
  margin-top: 0;
  width: 100%;
  box-sizing: border-box;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 20px 24px 24px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04);
}

.visit-page-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin: 8px 0 16px;
  flex-wrap: wrap;
}

.heading-left {
  flex: 1;
  min-width: 280px;
}

.visit-page-heading .page-title {
  margin: 0;
  color: #1e293b;
  font-size: 1.5rem;
  font-weight: 800;
  line-height: 1.3;
}

.visit-page-heading .page-desc {
  margin: 4px 0 0;
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
  line-height: 1.5;
}

.heading-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.preview-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 9px 14px;
  border-radius: 9px;
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
  font-size: 0.84rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.15s ease;
}

.preview-link-btn:hover {
  background: #dbeafe;
  transform: translateY(-1px);
}

.outline-btn,
.primary-btn,
.secondary-btn {
  border-radius: 9px;
  padding: 9px 16px;
  font-weight: 700;
  font-size: 0.86rem;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-family: inherit;
  transition: all 0.15s ease;
}

.outline-btn {
  border: 1px solid #cbd5e1;
  background: #ffffff;
  color: #334155;
}

.outline-btn:hover {
  border-color: #94a3b8;
  background: #f8fafc;
}

.primary-btn {
  background: #2563eb;
  border: 1px solid #2563eb;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.22);
}

.primary-btn:hover {
  background: #1d4ed8;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.3);
}

.secondary-btn {
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  color: #334155;
  text-decoration: none;
}

.secondary-btn:hover {
  background: #e2e8f0;
}

/* Mode Navigation Bar */
.mode-navigation-bar {
  padding: 16px 0 0;
  border-bottom: 1px solid #f1f5f9;
}

.mode-tabs {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.mode-tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-bottom: none;
  border-radius: 12px 12px 0 0;
  color: #64748b;
  font-size: 0.88rem;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s ease;
  position: relative;
  bottom: -1px;
}

.mode-tab-btn:hover {
  color: #2563eb;
  background: #ffffff;
}

.mode-tab-btn.active {
  background: #ffffff;
  color: #2563eb;
  font-weight: 700;
  border-color: #e2e8f0;
  border-bottom: 2px solid #ffffff;
  box-shadow: 0 -2px 8px rgba(15, 23, 42, 0.04);
}

.mode-active-pill {
  font-size: 0.68rem;
  font-weight: 700;
  background: #dcfce7;
  color: #15803d;
  padding: 2px 7px;
  border-radius: 6px;
}

.override-badge {
  font-size: 0.68rem;
  font-weight: 700;
  background: #fef3c7;
  color: #b45309;
  padding: 2px 7px;
  border-radius: 6px;
}

/* Mode Container */
.mode-container {
  padding-top: 0;
}

/* Standing Schedule Banner */
.standing-schedule-banner {
  display: flex;
  align-items: center;
  gap: 14px;
  background: linear-gradient(135deg, #eff6ff 0%, #f0fdf4 100%);
  border: 1px solid #bfdbfe;
  border-radius: 12px;
  padding: 14px 18px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.ssb-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: #ffffff;
  color: #2563eb;
  display: grid;
  place-items: center;
  flex-shrink: 0;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.12);
}

.ssb-content {
  flex: 1;
  min-width: 260px;
}

.ssb-content strong {
  display: block;
  font-size: 0.92rem;
  color: #1e3a8a;
}

.ssb-content p {
  margin: 3px 0 0;
  font-size: 0.78rem;
  color: #3b82f6;
  line-height: 1.45;
}

/* Preset Dropdown */
.preset-dropdown-wrap {
  position: relative;
}

.preset-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  color: #1e293b;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.preset-btn:hover {
  background: #f8fafc;
  border-color: #94a3b8;
}

.preset-menu {
  position: absolute;
  top: calc(100% + 6px);
  right: 0;
  z-index: 50;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 6px;
  box-shadow: 0 10px 25px -5px rgba(15, 23, 42, 0.12);
  width: 280px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.preset-menu-item {
  width: 100%;
  text-align: left;
  border: none;
  background: transparent;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.12s ease;
  font-family: inherit;
}

.preset-menu-item:hover {
  background: #eff6ff;
}

.preset-menu-item .pmi-title {
  font-size: 0.82rem;
  font-weight: 700;
  color: #1e293b;
}

.preset-menu-item .pmi-desc {
  font-size: 0.72rem;
  color: #64748b;
  margin-top: 2px;
  line-height: 1.35;
}

.preset-menu-item.is-reset:hover {
  background: #fef2f2;
}

.preset-menu-item.is-reset .pmi-title {
  color: #dc2626;
}

.preset-menu-divider {
  height: 1px;
  background: #f1f5f9;
  margin: 2px 0;
}

/* ======================================================== */
/* MASTER-DETAIL 2-COLUMN LAYOUT                            */
/* ======================================================== */
.weekly-master-split {
  display: grid;
  grid-template-columns: 330px 1fr;
  gap: 20px;
  align-items: start;
  margin-top: 14px;
}

/* LEFT: Master Sidebar (7 Days List) */
.weekly-master-sidebar {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 20px 18px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.sidebar-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 12px;
  border-bottom: 1px solid #f1f5f9;
}

.sidebar-title {
  font-size: 1.08rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.sidebar-sub {
  font-size: 0.85rem;
  font-weight: 500;
  color: #64748b;
  display: block;
  margin-top: 3px;
}

.quick-copy-all-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 6px 12px;
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
  border-radius: 7px;
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
}

.quick-copy-all-btn:hover {
  background: #dbeafe;
}

.days-master-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.day-master-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  background: #ffffff;
  border: 1.5px solid #f1f5f9;
  border-radius: 12px;
  cursor: pointer;
  text-align: left;
  transition: all 0.15s ease;
  font-family: inherit;
}

.day-master-item:hover {
  border-color: #cbd5e1;
  background: #f8fafc;
}

.day-master-item.active {
  background: #eff6ff;
  border-color: #3b82f6;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.1);
}

.day-master-item.closed {
  opacity: 0.78;
}

.dmi-indicator {
  width: 4px;
  height: 28px;
  border-radius: 2px;
  background: transparent;
  flex-shrink: 0;
  transition: background 0.15s ease;
}

.day-master-item.active .dmi-indicator {
  background: #2563eb;
}

.dmi-content {
  flex: 1;
  min-width: 0;
}

.dmi-top-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.dmi-name {
  font-size: 0.96rem;
  font-weight: 700;
  color: #1e293b;
}

.day-master-item.active .dmi-name {
  color: #1d4ed8;
}

.dmi-bot-row {
  margin-top: 4px;
  font-size: 0.84rem;
}

.dmi-status-open {
  color: #16a34a;
  font-weight: 600;
}

.dmi-status-closed {
  color: #94a3b8;
  font-weight: 500;
}

.dmi-arrow {
  color: #cbd5e1;
  display: inline-flex;
  align-items: center;
  transition: all 0.15s ease;
}

.day-master-item.active .dmi-arrow {
  color: #2563eb;
  transform: translateX(2px);
}

.sidebar-tip-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 11px 14px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  color: #64748b;
  font-size: 0.84rem;
  font-weight: 500;
  line-height: 1.45;
}

.sidebar-tip-box svg {
  flex-shrink: 0;
  color: #3b82f6;
}

/* RIGHT: Detail Panel */
.weekly-detail-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

/* Active Day Header */
.detail-day-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
  flex-wrap: wrap;
}

.ddh-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ddh-day-title {
  font-size: 1.32rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}

.ddh-subtitle {
  margin: 5px 0 0;
  font-size: 0.92rem;
  font-weight: 500;
  color: #64748b;
}

.ddh-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.day-toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 25px;
}

.day-toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.day-toggle-switch .slider {
  position: absolute;
  cursor: pointer;
  inset: 0;
  background-color: #cbd5e1;
  transition: 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 25px;
}

.day-toggle-switch .slider:before {
  position: absolute;
  content: "";
  height: 19px;
  width: 19px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 50%;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

.day-toggle-switch input:checked + .slider {
  background-color: #2563eb;
}

.day-toggle-switch input:checked + .slider:before {
  transform: translateX(19px);
}

.toggle-text {
  font-size: 0.92rem;
  font-weight: 700;
  color: #64748b;
}

.toggle-text.is-active {
  color: #1e293b;
}

/* Day Config (Badge & Advice) */
.detail-day-config {
  display: flex;
  align-items: flex-end;
  gap: 14px;
  padding: 18px 22px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  flex-wrap: wrap;
}

.config-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.config-field.badge-field {
  flex: 1;
  min-width: 0;
}

.config-field.advice-field {
  flex: 1;
  min-width: 0;
}

/* Slots Management Box */
.detail-slots-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 20px 22px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.03);
}

.slots-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 14px;
  border-bottom: 1px solid #f1f5f9;
  margin-bottom: 14px;
}

.sch-title {
  font-size: 1.06rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.sch-sub {
  font-size: 0.85rem;
  font-weight: 500;
  color: #64748b;
  display: block;
  margin-top: 3px;
}

.btn-add-slot-top {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: #eff6ff;
  border: 1px solid #93c5fd;
  color: #2563eb;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-add-slot-top:hover {
  background: #dbeafe;
}

.field-label {
  font-size: 0.875rem;
  font-weight: 700;
  color: #334155;
}

.holiday-placeholder-card { min-height: 0; display: block; background: #ffffff; border: 1px solid #e2e8f0; border-radius: 16px; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.04); }
.holiday-placeholder-card .empty-state-panel { display: grid; justify-items: center; max-width: 420px; text-align: center; }
.holiday-placeholder-card .empty-icon { display: grid; place-items: center; width: 56px; height: 56px; margin-bottom: 14px; border-radius: 16px; background: #eff6ff; color: #2563eb; }
.holiday-placeholder-card .empty-title { color: #1e293b; font-size: 1rem; font-weight: 700; }
.holiday-placeholder-card .empty-desc { margin-top: 6px; color: #64748b; font-size: .875rem; line-height: 1.5; }
.holiday-manager { padding: 0; }
.holiday-manager-head { display: flex; justify-content: space-between; gap: 16px; align-items: flex-start; padding: 4px 0 18px; border-bottom: 1px solid #eef2f6; }
.holiday-manager-head h2 { margin: 0; color: #1e293b; font-size: 1.125rem; }
.holiday-manager-head p { margin: 5px 0 0; color: #64748b; font-size: .875rem; }
.holiday-count { color: #2563eb; background: #eff6ff; border-radius: 999px; padding: 5px 10px; font-size: .8rem; font-weight: 700; white-space: nowrap; }
.holiday-add-row { display: grid; grid-template-columns: 180px 1fr auto; gap: 10px; margin-top: 16px; align-items: center; }
.holiday-list { display: grid; gap: 8px; margin-top: 16px; }
.holiday-row { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 12px 14px; border: 1px solid #e2e8f0; border-radius: 10px; background: #f8fafc; }
.holiday-row div { display: flex; flex-direction: column; gap: 3px; }.holiday-row strong { color: #1e293b; }.holiday-row span { color: #64748b; font-size: .8125rem; }.holiday-remove { color: #dc2626; border-color: #fecaca; }
.holiday-empty { margin-top: 18px; padding: 22px; border: 1px dashed #cbd5e1; border-radius: 10px; color: #64748b; text-align: center; font-size: .875rem; }

.day-booked-badge { display: inline-flex; align-items: center; margin-left: 6px; padding: 2px 7px; border-radius: 999px; background: #dbeafe; color: #1d4ed8; font-size: .75rem; font-weight: 700; white-space: nowrap; }

.styled-input {
  width: 100%;
  height: 44px;
  box-sizing: border-box;
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 0.9375rem;
  color: #1e293b;
  font-family: inherit;
  background: #ffffff;
  transition: border-color 0.15s ease;
}

.styled-input:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.15);
}

.copy-day-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 9px 14px;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  white-space: nowrap;
  font-family: inherit;
  transition: all 0.15s ease;
}

.copy-day-btn:hover {
  background: #eff6ff;
  border-color: #93c5fd;
  color: #2563eb;
}

/* Slots Step & List */
.slots-step {
  border-top: 1px solid #f1f5f9;
  margin-top: 14px;
}

.slot-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 10px;
}

.slot-list-head {
  display: grid;
  grid-template-columns: 245px 138px 1fr 90px 38px;
  align-items: center;
  gap: 12px;
  padding: 0 14px;
  font-size: 0.825rem;
  font-weight: 600;
  color: #64748b;
}

.slot-row {
  display: grid;
  grid-template-columns: 245px 138px 1fr 90px 38px;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #ffffff;
  transition: all 0.15s ease;
}

.time-range-picker-group {
  display: flex;
  align-items: center;
  gap: 6px;
}

.time-range-separator {
  color: #94a3b8;
  font-size: 0.875rem;
  font-weight: 600;
  flex-shrink: 0;
}

/* Quota Capacity Stepper Input */
.cap-stepper-box {
  display: inline-flex;
  align-items: center;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  height: 38px;
  width: 100%;
  box-sizing: border-box;
  transition: all 0.15s ease;
  overflow: hidden;
}

.cap-stepper-box:focus-within {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.16);
}

.cap-step-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 100%;
  border: none;
  background: #f8fafc;
  color: #64748b;
  cursor: pointer;
  transition: all 0.12s ease;
  flex-shrink: 0;
  padding: 0;
}

.cap-step-btn:hover:not(:disabled) {
  background: #eff6ff;
  color: #2563eb;
}

.cap-step-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.cap-input-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  min-width: 0;
  height: 100%;
  padding: 0 2px;
}

.cap-number-input {
  width: 32px;
  border: none;
  background: transparent;
  text-align: right;
  font-size: 0.9375rem;
  font-weight: 700;
  color: #1e293b;
  padding: 0;
  margin: 0;
  outline: none;
  font-family: inherit;
  -moz-appearance: textfield;
}

.cap-number-input::-webkit-outer-spin-button,
.cap-number-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.cap-suffix {
  font-size: 0.84rem;
  color: #64748b;
  font-weight: 500;
  margin-left: 4px;
  user-select: none;
  white-space: nowrap;
}

.slot-row:hover {
  border-color: #cbd5e1;
  box-shadow: 0 2px 6px rgba(15, 23, 42, 0.04);
}

.slot-row.disabled {
  background: #f8fafc;
  opacity: 0.65;
}

.slot-row.is-sim-full {
  background: #fffdf5;
  border-color: #fde68a;
}

.custom-check {
  width: 18px;
  height: 18px;
  accent-color: #2563eb;
  cursor: pointer;
}

.time-input,
.note-input {
  width: 100%;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 12px;
  color: #1e293b;
  font-family: inherit;
  font-size: 0.92rem;
}

.time-input:focus,
.note-input:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.12);
}

/* Status Toggle Switch (เปิด / ปิด) */
.status-toggle-control {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
}

.status-toggle-input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.status-toggle-slider {
  position: relative;
  width: 36px;
  height: 20px;
  background-color: #cbd5e1;
  border-radius: 20px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
}

.status-toggle-slider::before {
  content: "";
  position: absolute;
  height: 14px;
  width: 14px;
  left: 3px;
  bottom: 3px;
  background-color: #ffffff;
  border-radius: 50%;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
}

.status-toggle-input:checked + .status-toggle-slider {
  background-color: #16a34a;
}

.status-toggle-input:checked + .status-toggle-slider::before {
  transform: translateX(16px);
}

.status-toggle-label {
  font-size: 0.875rem;
  font-weight: 700;
  color: #94a3b8;
  min-width: 26px;
  transition: color 0.15s ease;
}

.status-toggle-label.is-active {
  color: #16a34a;
}

.delete-button {
  border: none;
  background: transparent;
  color: #94a3b8;
  padding: 6px;
  cursor: pointer;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.delete-button:hover {
  color: #dc2626;
  background: #fee2e2;
}

.add-slot-row-btn {
  margin-top: 4px;
  padding: 10px 14px;
  border: 1.5px dashed #cbd5e1;
  border-radius: 10px;
  background: #f8fafc;
  color: #2563eb;
  font-size: 0.875rem;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s ease;
}

.add-slot-row-btn:hover {
  background: #eff6ff;
  border-color: #93c5fd;
}

.closed-day-banner {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 22px;
  border-radius: 14px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  color: #92400e;
  margin-top: 4px;
}

.closed-day-banner svg {
  flex-shrink: 0;
  color: #d97706;
}

.closed-day-banner strong {
  display: block;
  font-size: 1rem;
  font-weight: 700;
  color: #78350f;
}

.closed-day-banner p {
  margin: 4px 0 0;
  font-size: 0.9rem;
  font-weight: 500;
  color: #92400e;
}

.reopen-day-btn {
  margin-left: auto;
  padding: 9px 18px;
  background: #2563eb;
  color: #ffffff;
  border: none;
  border-radius: 9px;
  font-size: 0.9rem;
  font-weight: 700;
  cursor: pointer;
  white-space: nowrap;
  font-family: inherit;
  transition: background 0.15s ease, transform 0.15s ease;
}

.reopen-day-btn:hover {
  background: #1d4ed8;
  transform: translateY(-1px);
}

/* Mode 2: Calendar & Overrides Styles */
.calendar-horizon-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 14px 18px;
  margin-bottom: 16px;
  gap: 14px;
  flex-wrap: wrap;
}

.calendar-horizon-banner strong {
  display: block;
  font-size: 0.92rem;
  color: #0f172a;
}

.calendar-horizon-banner p {
  margin: 3px 0 0;
  font-size: 0.78rem;
  color: #64748b;
}

.horizon-capsule {
  display: inline-flex;
  padding: 3px;
  background: #ffffff;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
}

.horizon-tab {
  border: none;
  background: transparent;
  padding: 5px 12px;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 600;
  color: #64748b;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}

.horizon-tab.active {
  background: #2563eb;
  color: #ffffff;
  font-weight: 700;
}

.date-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
  gap: 10px;
  margin-top: 8px;
  overflow-x: auto;
  padding-bottom: 6px;
}

.date-card {
  min-width: 125px;
  text-align: left;
  padding: 12px 14px;
  border: 1.5px solid #e2e8f0;
  border-radius: 12px;
  background: #ffffff;
  color: #334155;
  cursor: pointer;
  transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
}

.date-card:hover {
  border-color: #93c5fd;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.08);
}

.date-card.active {
  border-color: #2563eb;
  background: #eff6ff;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.18);
}

.date-card.closed {
  opacity: 0.65;
  background: #f8fafc;
}

.date-card.has-override {
  border-color: #f59e0b;
}

.date-card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 4px;
  margin-bottom: 8px;
  min-height: 20px;
}

.date-status {
  font-size: 0.72rem;
  font-weight: 700;
}

.date-status.status-open { color: #16a34a; }
.date-status.status-closed { color: #94a3b8; }

.date-badge {
  font-size: 0.78rem;
  font-weight: 700;
  padding: 3px 9px;
  border-radius: 6px;
  white-space: nowrap;
}

.badge-primary { background: #dbeafe; color: #1e40af; }
.badge-success { background: #dcfce7; color: #166534; }
.badge-purple { background: #f3e8ff; color: #6b21a8; }
.badge-neutral { background: #f1f5f9; color: #475569; }
.badge-override { background: #fef3c7; color: #b45309; }

.date-day-name { font-size: 0.8rem; font-weight: 700; color: #475569; }
.date-card.active .date-day-name { color: #1e40af; }
.date-day-num { font-size: 1.6rem; color: #0f172a; line-height: 1.1; margin: 2px 0 3px; font-weight: 800; }
.date-card.active .date-day-num { color: #1d4ed8; }
.date-month { font-size: 0.74rem; color: #64748b; font-weight: 500; }
.date-card-bot { margin-top: 8px; padding-top: 6px; border-top: 1px dashed #e2e8f0; font-size: 0.72rem; }
.day-cap-tag { color: #2563eb; font-weight: 700; }
.day-closed-tag { color: #94a3b8; font-weight: 500; }

.single-override-box {
  background: #f8fafc;
  border: 1.5px solid #cbd5e1;
  border-radius: 12px;
  padding: 18px 20px;
  margin-top: 18px;
}

.sob-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.sob-head h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 700;
  color: #0f172a;
}

.sob-head p {
  margin: 4px 0 0;
  font-size: 0.8rem;
  color: #64748b;
}

.sob-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.revert-override-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  background: #ffffff;
  border-radius: 8px;
  font-size: 0.78rem;
  font-weight: 600;
  color: #64748b;
  cursor: pointer;
}

.revert-override-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.quick-toggle-day-btn {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  border: none;
  font-family: inherit;
  transition: all 0.15s ease;
}

.quick-toggle-day-btn.btn-red {
  background: #fee2e2;
  color: #b91c1c;
}

.quick-toggle-day-btn.btn-red:hover {
  background: #fecaca;
  color: #991b1b;
}

.quick-toggle-day-btn.btn-green {
  background: #dcfce7;
  color: #15803d;
}

.quick-toggle-day-btn.btn-green:hover {
  background: #bbf7d0;
  color: #166534;
}

/* Footer Bar */
.schedule-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-top: 28px;
  padding-top: 20px;
  border-top: 1px solid #f1f5f9;
  flex-wrap: wrap;
  gap: 12px;
}

.footer-left .footer-hint {
  font-size: 0.8rem;
  color: #475569;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* Copy Modal */
.copy-modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(4px);
  display: grid;
  place-items: center;
  padding: 20px;
}

.copy-modal-card {
  width: 100%;
  max-width: 480px;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 20px 40px -10px rgba(15, 23, 42, 0.2);
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.cm-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
}

.cm-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cm-title-wrap h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 700;
  color: #0f172a;
}

.cm-close {
  border: none;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
}

.cm-body {
  padding: 18px 20px;
}

.cm-desc {
  margin: 0 0 14px;
  font-size: 0.84rem;
  color: #475569;
}

.copy-target-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.target-option-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 14px;
  border: 1.5px solid #e2e8f0;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.target-option-card:hover {
  border-color: #93c5fd;
  background: #f8fafc;
}

.target-option-card input {
  margin-top: 3px;
  accent-color: #2563eb;
}

.target-option-card strong {
  display: block;
  font-size: 0.86rem;
  color: #1e293b;
}

.target-option-card p {
  margin: 2px 0 0;
  font-size: 0.76rem;
  color: #64748b;
}

.cm-foot {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 20px;
  background: #f8fafc;
  border-top: 1px solid #f1f5f9;
}

@media (max-width: 992px) {
  .weekly-master-split {
    grid-template-columns: 1fr;
  }

  .days-master-list {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 8px;
  }
}

@media (max-width: 768px) {
  .visit-page {
    padding: 16px 16px 36px;
  }

  .days-master-list {
    grid-template-columns: repeat(2, 1fr);
  }

  .slot-list-head {
    display: none;
  }

  .slot-row {
    grid-template-columns: 1fr 36px;
    gap: 10px;
  }

  .slot-col-time,
  .slot-col-cap,
  .slot-col-note {
    grid-column: 1;
  }

  .slot-col-status {
    grid-column: 1;
    justify-self: start;
  }

  .slot-col-actions {
    grid-column: 2;
    grid-row: 1;
  }
}

@media (max-width: 640px) {
  .visit-page {
    padding: 12px 12px 28px;
  }

  .schedule-card {
    padding: 18px;
  }

  .visit-page-heading {
    flex-direction: column;
    align-items: flex-start;
  }

  .detail-day-config {
    flex-direction: column;
    align-items: stretch;
  }

  .config-field.badge-field {
    width: 100%;
  }

  .date-grid {
    grid-template-columns: repeat(auto-fill, minmax(115px, 1fr));
  }
}
</style>
