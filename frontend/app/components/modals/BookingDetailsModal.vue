<template>
  <div v-if="show && booking" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content" style="max-width: 640px;">
      <AppModalHeader
        title="รายละเอียดการจองดูแผงค้า"
        :subtitle="`เลขที่การจอง ${booking.booking_number}`"
        icon="calendar"
        @close="$emit('close')"
      />

      <div class="modal-body">
        <!-- Status Banner -->
        <div class="modal-status-banner" :class="getStatusBannerClass(booking.status)">
          <div>
            <div style="font-size: 0.78rem; text-transform: uppercase; font-weight: 600; opacity: 0.85;">สถานะปัจจุบัน</div>
            <div style="font-size: 1.1rem; font-weight: 800; display: flex; align-items: center; gap: 8px; margin-top: 2px; flex-wrap: wrap;">
              <span class="badge-dot"></span>
              <span>{{ getStatusLabel(booking.status) }}</span>
              <span v-if="isQueueBooking(booking)" style="font-size: 0.78rem; padding: 2px 10px; border-radius: 999px; background: var(--bg-fef3c7); color: var(--tx-b45309); border: 1px solid var(--bd-fde68a); font-weight: 700;">
                คิวสำรองที่ {{ booking.queue_position || booking.queue_number || 1 }}
              </span>
              <span v-else-if="booking.holding_stall && !['rejected', 'cancelled', 'not_interested', 'no_show', 'lease_applied', 'completed'].includes(booking.status)" style="font-size: 0.78rem; padding: 2px 10px; border-radius: 999px; background: var(--bg-ecfdf5); color: var(--tx-15803d); border: 1px solid var(--bd-a7f3d0); font-weight: 700;">
                ผู้จองหลัก (ตัวจริง)
              </span>
            </div>
          </div>
          <div class="modal-status-dates" style="text-align: right; font-size: 0.8rem;">
            <div>ทำรายการเมื่อ: {{ formatThaiDateTime(booking.created_at) }}</div>
            <div v-if="booking.reviewed_at" style="color: var(--tx-34d399);">อนุมัติเมื่อ: {{ formatThaiDateTime(booking.reviewed_at) }}</div>
          </div>
        </div>

        <!-- Flow Progress Indicator (3-Step Lifecycle) -->
        <div class="flow-steps">
          <div class="flow-step" :class="{ done: true, active: isStepActive(booking.status, 1) }">
            <div class="step-num">1</div>
            <div class="step-txt">จองดูแผง</div>
          </div>
          <div class="flow-step-line" :class="{ done: isStepPassed(booking.status, 2) }"></div>
          <div class="flow-step" :class="{ done: isStepPassed(booking.status, 2), active: isStepActive(booking.status, 2) }">
            <div class="step-num">2</div>
            <div class="step-txt">นัดชมแผง</div>
          </div>
          <div class="flow-step-line" :class="{ done: isStepPassed(booking.status, 3) }"></div>
          <div class="flow-step" :class="{ done: isStepPassed(booking.status, 3), active: isStepActive(booking.status, 3) }">
            <div class="step-num">3</div>
            <div class="step-txt">{{ isStepPassed(booking.status, 3) ? 'ยื่นคำร้องสำเร็จ' : 'ยื่นคำร้องขอเช่า' }}</div>
          </div>
        </div>

        <!-- Waiting Queue Banner if in queue -->
        <div v-if="isQueueBooking(booking)" class="queue-notice-banner" style="background: var(--bg-fffbeb); border: 1px solid var(--bd-fde68a); border-radius: 12px; padding: 14px 18px; margin-bottom: 18px; display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px;">
          <div style="display: flex; align-items: flex-start; gap: 12px; flex: 1; min-width: 260px;">
            <div style="width: 42px; height: 42px; border-radius: 50%; background: var(--bg-fef3c7); color: var(--tx-b45309); display: flex; align-items: center; justify-content: center; font-weight: 800; font-size: 1.15rem; flex-shrink: 0; border: 2px solid var(--bd-fde68a); margin-top: 2px;">
              {{ booking.queue_position || booking.queue_number || 1 }}
            </div>
            <div>
              <div style="font-weight: 700; color: var(--tx-92400e); font-size: 0.95rem;">
                คิวสำรองลำดับที่ {{ booking.queue_position || booking.queue_number || 1 }} ของแผง {{ booking.stall_number }}
              </div>
              <div style="font-size: 0.82rem; color: var(--tx-b45309); margin-top: 3px; line-height: 1.45;">
                <span v-if="(booking.queue_position || booking.queue_number || 1) === 1">
                  แผงนี้มี <strong>ผู้จองหลักถือสิทธิ์อยู่ 1 ราย</strong><span v-if="activeHolder"> ({{ formatQueueHolderName(activeHolder.tenant_name) }})</span> — คิวสำรองที่ 1 จะได้รับสิทธิ์เป็นผู้จองหลักทันทีหากผู้จองหลักสละสิทธิ์
                </span>
                <span v-else>
                  แผงนี้มี <strong>ผู้จองหลัก 1 ราย</strong><span v-if="activeHolder"> ({{ formatQueueHolderName(activeHolder.tenant_name) }})</span>
                  และมี <strong>คิวสำรองก่อนหน้า {{ (booking.queue_position || booking.queue_number || 1) - 1 }} ราย</strong>
                  <span v-if="precedingQueues && precedingQueues.length > 0">
                    ({{ precedingQueues.map(q => `คิว ${q.queue_position || q.queue_number || 1}: ${formatQueueHolderName(q.tenant_name)}`).join(', ') }})
                  </span>
                </span>
              </div>
              <div style="font-size: 0.78rem; color: var(--tx-78350f); margin-top: 5px; line-height: 1.4;">
                * สามารถเข้าชมแผงตามนัดได้ตามปกติ โดยปุ่มยืนยันการเช่าจะเปิดให้กดเมื่อผู้จองคิวหลักสละสิทธิ์
              </div>
            </div>
          </div>
          <span style="font-size: 0.78rem; background: var(--bg-ffffff); padding: 4px 12px; border-radius: 20px; border: 1px solid var(--bd-fde68a); color: var(--tx-92400e); font-weight: 600; white-space: nowrap;">
            ระบบขยับคิวอัตโนมัติ
          </span>
        </div>

        <!-- 1. Stall Details -->
        <div style="background: var(--bg-f8fafc); border: 1px solid var(--border-subtle); padding: 16px; border-radius: var(--radius-md); margin-bottom: 18px;">
          <h4 style="font-size: 0.88rem; color: var(--primary); font-weight: 700; margin-bottom: 10px; display: flex; align-items: center; gap: 6px;">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
            <span>ข้อมูลแผงและตลาด</span>
          </h4>
          <div class="modal-info-grid">
            <div>ตลาด: <strong style="color: var(--tx-0f172a);">{{ booking.market_name }}</strong></div>
            <div>โครงการ: <strong style="color: var(--tx-1e293b);">{{ booking.project_name || 'โครงการ 01' }}</strong></div>
            <div>เลขที่แผง: <strong style="color: var(--tx-2563eb); font-size: 1.05rem;">{{ booking.stall_number }} <span v-if="booking.stall_count && booking.stall_count > 1" style="font-size: 0.82rem; font-weight: 600; color: var(--tx-64748b);">(รวม {{ booking.stall_count }} แผง)</span></strong></div>
            <div>โซนตลาด: <strong style="color: var(--tx-334155);">{{ booking.stall_zone_name || '-' }}</strong></div>
            <div>อัตราค่าเช่า: <strong style="color: var(--tx-059669);">฿{{ Number(booking.stall_monthly_rate || 0).toLocaleString() }}/เดือน</strong></div>
          </div>
        </div>

        <!-- 2. Appointment Details -->
        <div style="background: var(--bg-f8fafc); border: 1px solid var(--border-subtle); padding: 16px; border-radius: var(--radius-md); margin-bottom: 18px;">
          <h4 style="font-size: 0.88rem; color: var(--primary); font-weight: 700; margin-bottom: 10px; display: flex; align-items: center; gap: 6px;">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            <span>ข้อมูลการนัดหมายและการประสานงาน</span>
          </h4>
          <div class="modal-info-grid">
            <div>วันนัดหมาย: <strong style="color: var(--tx-0f172a);">{{ formatVisitDate(booking) }}</strong></div>
            <div>ช่วงเวลา: <strong style="color: var(--tx-0f172a);">{{ getVisitSlot(booking) }}</strong></div>
            <div class="grid-col-full">จุดนัดพบ: <strong style="color: var(--tx-0f172a);">{{ booking.meeting_point || 'หน้ากองอำนวยการตลาด เสาธง 1' }}</strong></div>
            <div class="grid-col-full">จนท. ผู้พาชม: <strong style="color: var(--tx-2563eb);">{{ booking.escort_officer || 'เจ้าหน้าที่เวรฝ่ายปฏิบัติการ' }}</strong></div>
            <div class="grid-col-full">สินค้าที่ต้องการจำหน่าย: <strong style="color: var(--tx-334155);">{{ booking.intended_business || '-' }}</strong></div>
          </div>
          <div v-if="booking.officer_notes" style="margin-top: 10px; padding: 10px 14px; background: var(--bg-ecfdf5); border: 1px solid var(--bd-a7f3d0); border-radius: 8px; font-size: 0.85rem; color: var(--tx-065f46); display: flex; align-items: flex-start; gap: 8px;">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink: 0; margin-top: 2px;"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
            <div><strong>บันทึกจากเจ้าหน้าที่:</strong> {{ booking.officer_notes }}</div>
          </div>
          <div v-if="booking.rejection_reason || booking.cancel_reason" style="margin-top: 10px; padding: 10px 14px; background: var(--bg-fff1f2); border: 1px solid var(--bd-fecdd3); border-radius: 8px; font-size: 0.85rem; color: var(--tx-9f1239); display: flex; align-items: flex-start; gap: 8px;">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink: 0; margin-top: 2px;"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
            <div><strong>เหตุผลที่ยกเลิก:</strong> {{ booking.rejection_reason || booking.cancel_reason }}</div>
          </div>
          <div v-if="booking.not_interested_reason" style="margin-top: 10px; padding: 10px 14px; background: var(--bg-f8fafc); border: 1px solid var(--bd-e2e8f0); border-radius: 8px; font-size: 0.85rem; color: var(--tx-475569); display: flex; align-items: flex-start; gap: 8px;">
            <div><strong>เหตุผลที่ไม่สนใจเช่า:</strong> {{ booking.not_interested_reason }}</div>
          </div>
        </div>

        <!-- 3. Tenant Details -->
        <div style="background: var(--bg-f8fafc); border: 1px solid var(--border-subtle); padding: 16px; border-radius: var(--radius-md); margin-bottom: 18px;">
          <h4 style="font-size: 0.88rem; color: var(--primary); font-weight: 700; margin-bottom: 10px; display: flex; align-items: center; gap: 6px;">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            <span>ข้อมูลผู้ขอจองเข้าชม</span>
          </h4>
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px; font-size: 0.88rem;">
            <div>ชื่อ-นามสกุล: <strong style="color: var(--tx-0f172a);">{{ booking.tenant_name }}</strong></div>
            <div>เบอร์โทรศัพท์: <strong style="color: var(--tx-0f172a);">{{ booking.tenant_phone }}</strong></div>
            <div>อีเมล: <strong style="color: var(--tx-334155);">{{ booking.tenant_email || '-' }}</strong></div>
            <div>เลขบัตรประชาชน: <strong style="color: var(--tx-334155);">{{ booking.tenant_id_card || '-' }}</strong></div>
          </div>
        </div>

        <!-- 4. Audit Log Timeline (Rule 5: บันทึก log ผู้กระทำ, เวลา, หมายเหตุ ทุกขั้นตอน) -->
        <div style="background: var(--bg-ffffff); border: 1px solid var(--border-subtle); padding: 16px; border-radius: var(--radius-md); margin-bottom: 14px;">
          <h4 style="font-size: 0.88rem; color: var(--primary); font-weight: 700; margin-bottom: 12px; display: flex; align-items: center; gap: 6px;">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            <span>ประวัติไทม์ไลน์และบันทึกการดำเนินการ (Audit Logs)</span>
          </h4>
          <div v-if="booking.timeline && booking.timeline.length > 0" class="audit-timeline">
            <div v-for="(log, idx) in booking.timeline" :key="idx" class="timeline-entry">
              <div class="timeline-bullet" :class="getRoleClass(log.actor_role)"></div>
              <div class="timeline-content">
                <div style="display: flex; justify-content: space-between; align-items: center; gap: 8px;">
                  <span style="font-weight: 700; font-size: 0.85rem; color: var(--tx-0f172a);">{{ log.action }}</span>
                  <span style="font-size: 0.75rem; color: var(--tx-64748b);">{{ log.timestamp?.slice(0, 16).replace('T', ' ') }} น.</span>
                </div>
                <div style="font-size: 0.75rem; color: var(--tx-475569); margin-top: 2px; display: flex; align-items: center; gap: 6px;">
                  <span class="role-tag" :class="getRoleClass(log.actor_role)">{{ getActorRoleName(log.actor_role) }}</span>
                  <span>{{ log.actor_name }}</span>
                </div>
                <div v-if="log.note" style="font-size: 0.78rem; color: var(--tx-334155); margin-top: 4px; background: var(--bg-f8fafc); padding: 4px 8px; border-radius: 6px; border: 1px solid var(--bd-f1f5f9);">
                  {{ log.note }}
                </div>
              </div>
            </div>
          </div>
          <div v-else style="font-size: 0.8rem; color: var(--tx-94a3b8); text-align: center; padding: 12px;">
            ไม่มีประวัติเพิ่มเติม
          </div>
        </div>
      </div>

      <div class="modal-footer" style="justify-content: space-between; flex-wrap: wrap; gap: 10px;">
        <AppButton variant="secondary" size="md" @click="$emit('close')">ปิดหน้าต่าง</AppButton>

        <div style="display: flex; gap: 8px; flex-wrap: wrap;">
          <!-- OFFICER ROLE ACTIONS -->
          <template v-if="role === 'officer'">
            <template v-if="booking.status === 'pending_approval' || booking.status === 'pending'">
              <AppButton variant="danger" size="md" @click="isQueueBooking(booking) ? $emit('cancel', booking) : $emit('reject', booking)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                <span>{{ isQueueBooking(booking) ? 'ยกเลิกคิว' : 'ยกเลิก' }}</span>
              </AppButton>
              <AppButton variant="primary" size="md" @click="$emit('approve', booking)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span>อนุมัตินัดหมาย</span>
              </AppButton>
            </template>
            <template v-else-if="booking.status === 'approved'">
              <AppButton variant="danger" size="md" @click="$emit('cancel', booking)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                <span>{{ isQueueBooking(booking) ? 'ยกเลิกคิว' : 'ยกเลิก' }}</span>
              </AppButton>
              <AppButton v-if="!isQueueBooking(booking)" variant="success" size="md" @click="$emit('markVisited', booking)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span>ยืนยันการเช่า</span>
              </AppButton>
            </template>
            <template v-else-if="booking.status === 'interested'">
              <AppButton variant="danger" size="md" @click="$emit('cancel', booking)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                <span>ยกเลิก</span>
              </AppButton>
            </template>
          </template>

          <!-- TENANT ROLE ACTIONS -->
          <template v-else>
            <template v-if="['approved', 'pending_approval', 'pending', 'waiting_queue', 'pending_slot'].includes(booking.status)">
              <AppButton variant="danger-outline" size="md" @click="$emit('cancel', booking)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                <span>{{ isQueueBooking(booking) ? 'สละสิทธิ์' : 'ยกเลิก' }}</span>
              </AppButton>
            </template>
            <template v-else-if="booking.status === 'interested'">
              <AppButton variant="primary" size="md" @click="$emit('applyLease', booking)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                <span>ยื่นคำร้องขอเช่าแผงนี้ทันที</span>
              </AppButton>
            </template>
            <template v-else-if="booking.status === 'lease_applied' || booking.status === 'completed'">
              <AppButton to="/applications" variant="primary" size="md" @click="$emit('close')">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
                <span>ดูสถานะคำร้องขอเช่า ({{ booking.application_number || 'APP' }})</span>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
              </AppButton>
            </template>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  show: boolean
  booking: any
  role?: 'tenant' | 'officer'
  activeHolder?: any
  precedingQueues?: any[]
}>(), {
  role: 'tenant',
  activeHolder: null,
  precedingQueues: () => []
})

// PDPA Masking Helper สำหรับชื่อผู้จองก่อนหน้า (ซ่อนนามสกุลเพื่อความเป็นส่วนตัวตามกฎหมาย PDPA)
const maskPdpaName = (fullName: string): string => {
  if (!fullName || typeof fullName !== 'string') return 'ผู้เช่าท่านอื่น'
  let clean = fullName.trim()
  clean = clean.replace(/^(นาย|นางสาว|นาง|น\.ส\.|คุณ)\s*/, '')
  const parts = clean.split(/\s+/).filter(Boolean)
  if (parts.length === 0) return 'ผู้เช่าท่านอื่น'
  const firstName = parts[0]
  if (parts.length > 1) {
    const lastName = parts.slice(1).join(' ')
    return `คุณ${firstName} ${lastName.charAt(0)}.***`
  }
  if (firstName.length > 3) return `คุณ${firstName.slice(0, 3)}***`
  return `คุณ${firstName}***`
}

const formatQueueHolderName = (name: string): string => {
  if (!name) return 'ผู้เช่าท่านอื่น'
  if (props.role === 'officer') return name
  return maskPdpaName(name)
}

defineEmits<{
  (e: 'close'): void
  (e: 'applyLease', booking: any): void
  (e: 'markVisited', booking: any): void
  (e: 'noShow', booking: any): void
  (e: 'cancel', booking: any): void
  (e: 'approve', booking: any): void
  (e: 'reject', booking: any): void
  (e: 'confirmInterest', booking: any): void
  (e: 'notInterested', booking: any): void
  (e: 'recordProxy', booking: any): void
}>()

const isQueueBooking = (b: any): boolean => {
  if (!b) return false
  if (['rejected', 'cancelled', 'not_interested', 'no_show', 'lease_applied', 'completed', 'terminated'].includes(b.status)) {
    return false
  }
  return b.status === 'waiting_queue' || (!b.holding_stall && Number(b.queue_position || b.queue_number) > 0)
}

const isStepPassed = (status: string, step: number) => {
  if (step === 1) {
    return ['waiting_queue', 'approved', 'interested', 'lease_applied', 'completed'].includes(status)
  }
  if (step === 2) {
    return ['interested', 'lease_applied', 'completed'].includes(status)
  }
  if (step === 3) {
    return ['lease_applied', 'completed'].includes(status)
  }
  return false
}

const isStepActive = (status: string, step: number) => {
  if (step === 1) return ['pending_slot', 'pending_approval', 'pending'].includes(status)
  if (step === 2) return ['approved', 'waiting_queue'].includes(status)
  if (step === 3) return ['interested'].includes(status)
  return false
}

const getStatusBannerClass = (status: string) => {
  switch (status) {
    case 'waiting_queue': return 'banner-queue'
    case 'pending_slot': return 'banner-slot'
    case 'pending_approval': return 'banner-pending'
    case 'approved': return 'banner-approved'
    case 'pending_decision': return 'banner-decision'
    case 'interested': return 'banner-interested'
    case 'lease_applied': return 'banner-applied'
    case 'no_show':
    case 'rejected':
    case 'not_interested':
    case 'terminated':
    case 'cancelled': return 'banner-rejected'
    default: return ''
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'waiting_queue': return props.booking?.queue_position ? `รอคิว (คิวสำรองที่ ${props.booking.queue_position})` : 'รอคิว (ต่อคิวแผงติดจอง)'
    case 'pending_slot': return 'ถึงคิวแล้ว (รอเลือกวันนัดชม)'
    case 'pending_approval':
    case 'pending': return 'รอเจ้าหน้าที่อนุมัติ'
    case 'approved': return 'นัดหมายสำเร็จ'
    case 'pending_decision': return 'รอผู้เช่ายื่นคำร้อง'
    case 'interested': return 'รอผู้เช่ายื่นคำร้อง'
    case 'lease_applied': return 'ยื่นคำร้องสำเร็จ'
    case 'completed': return 'ยื่นคำร้องสำเร็จ'
    case 'no_show': return 'ไม่มาตามนัดหมาย (สละสิทธิ์)'
    case 'rejected':
    case 'cancelled': return 'ยกเลิก'
    case 'not_interested': return 'ไม่สนใจเช่า (สละสิทธิ์)'
    case 'terminated': return 'สิ้นสุดสัญญาเช่า'
    default: return status
  }
}

const getActorRoleName = (role: string) => {
  if (role === 'officer') return 'เจ้าหน้าที่'
  if (role === 'system') return 'ระบบ'
  return 'ผู้เช่า'
}

const getRoleClass = (role: string) => {
  if (role === 'officer') return 'role-officer'
  if (role === 'system') return 'role-system'
  return 'role-tenant'
}

const getVisitSlot = (b: any) => {
  if (b?.visit_time_slot && b.visit_time_slot !== '-' && b.visit_time_slot !== 'รอยืนยันเวลา') {
    return b.visit_time_slot
  }
  const slots = ['10.00 - 11.30 น.', '13.00 - 14.30 น.', '14.00 - 15.30 น.', '15.00 - 16.00 น.']
  const idx = Math.abs((Number(b?.id) || 1)) % slots.length
  return slots[idx]
}

const formatVisitDate = (b: any) => {
  let target = b?.visit_date
  if (!target && b?.created_at) {
    try {
      const cd = new Date(b.created_at)
      if (!isNaN(cd.getTime())) {
        cd.setDate(cd.getDate() + 5)
        target = cd.toISOString()
      }
    } catch {}
  }
  if (!target) return 'เสาร์ 17 ตุลาคม 2569'
  if (target.includes('จันทร์') || target.includes('อังคาร') || target.includes('พุธ') || target.includes('พฤหัสบดี') || target.includes('ศุกร์') || target.includes('เสาร์') || target.includes('อาทิตย์')) {
    return target
  }
  try {
    const d = new Date(target)
    if (isNaN(d.getTime())) return target
    const thaiDays = ['อาทิตย์', 'จันทร์', 'อังคาร', 'พุธ', 'พฤหัสบดี', 'ศุกร์', 'เสาร์']
    const thaiMonths = [
      'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
      'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
    ]
    const dayName = thaiDays[d.getDay()]
    const day = d.getDate()
    const month = thaiMonths[d.getMonth()]
    const year = d.getFullYear() + 543
    return `${dayName} ${day} ${month} ${year}`
  } catch {
    return target
  }
}

const formatThaiDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    const thaiMonths = [
      'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
      'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม'
    ]
    const day = d.getDate()
    const month = thaiMonths[d.getMonth()]
    const year = d.getFullYear() + 543
    const hours = String(d.getHours()).padStart(2, '0')
    const mins = String(d.getMinutes()).padStart(2, '0')
    return `${day} ${month} ${year} เวลา ${hours}.${mins} น.`
  } catch {
    return dateStr
  }
}
</script>

<style scoped>
.banner-queue {
  background: var(--bg-fefce8);
  border: 1px solid var(--bd-fef08a);
  color: var(--tx-854d0e);
}
.banner-slot {
  background: var(--bg-fdf4ff);
  border: 1px solid var(--bd-f5d0fe);
  color: var(--tx-86198f);
}
.banner-pending {
  background: var(--bg-fffbeb);
  border: 1px solid var(--bd-fde68a);
  color: var(--tx-b45309);
}
.banner-approved {
  background: var(--bg-ecfdf5);
  border: 1px solid var(--bd-a7f3d0);
  color: var(--tx-059669);
}
.banner-decision {
  background: var(--bg-eff6ff);
  border: 1px solid var(--bd-bfdbfe);
  color: var(--tx-1d4ed8);
}
.banner-interested {
  background: var(--bg-f0fdf4);
  border: 1px solid var(--bd-bbf7d0);
  color: var(--tx-15803d);
}
.banner-applied {
  background: var(--bg-f5f3ff);
  border: 1px solid var(--bd-ddd6fe);
  color: var(--tx-7c3aed);
}
.banner-rejected {
  background: var(--bg-fff1f2);
  border: 1px solid var(--bd-fecdd3);
  color: var(--tx-e11d48);
}

/* Audit log timeline styling */
.audit-timeline {
  display: flex;
  flex-direction: column;
  gap: 12px;
  position: relative;
  padding-left: 8px;
}
.timeline-entry {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  position: relative;
}
.timeline-bullet {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-top: 6px;
  flex-shrink: 0;
}
.timeline-bullet.role-tenant { background: #3b82f6; box-shadow: 0 0 0 3px #dbeafe; }
.timeline-bullet.role-officer { background: #10b981; box-shadow: 0 0 0 3px #d1fae5; }
.timeline-bullet.role-system { background: #8b5cf6; box-shadow: 0 0 0 3px #ede9fe; }

.timeline-content {
  flex: 1;
}

.role-tag {
  font-size: var(--font-size-xs, 0.78rem);
  padding: 2px 7px;
  border-radius: 4px;
  font-weight: 700;
}
.role-tag.role-tenant { background: var(--bg-eff6ff); color: var(--tx-2563eb); }
.role-tag.role-officer { background: var(--bg-ecfdf5); color: var(--tx-059669); }
.role-tag.role-system { background: var(--bg-f5f3ff); color: var(--tx-7c3aed); }

/* Flow steps timeline */
.flow-steps {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding: 0 10px;
}

.flow-step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.step-num {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--bg-f1f5f9);
  color: var(--tx-64748b);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  font-weight: 700;
  border: 2px solid var(--bd-cbd5e1);
  transition: var(--transition);
}

.flow-step.done .step-num {
  background: var(--primary);
  color: #fff;
  border-color: var(--primary);
}

.flow-step.active .step-num {
  border-color: var(--primary);
  color: var(--primary);
  background: var(--bg-eff6ff);
  box-shadow: 0 0 0 4px rgba(67, 97, 238, 0.15);
  transform: scale(1.05);
}

.step-txt {
  font-size: 0.75rem;
  color: var(--tx-64748b);
  font-weight: 500;
}

.flow-step.done .step-txt {
  color: var(--primary);
  font-weight: 700;
}

.flow-step-line {
  flex: 1;
  height: 2px;
  background: var(--bg-e2e8f0);
  margin: 0 8px -14px;
}

.flow-step-line.done {
  background: var(--primary);
}

/* Modal Responsive Classes */
.modal-status-banner {
  padding: 16px;
  border-radius: var(--radius-md);
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.modal-info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  font-size: 0.88rem;
}

.grid-col-full {
  grid-column: 1 / -1;
}

@media (max-width: 600px) {
  .flow-steps {
    padding: 0 4px;
    margin-bottom: 16px;
  }
  .step-num {
    width: 26px;
    height: 26px;
    font-size: 0.75rem;
  }
  .step-txt {
    font-size: var(--font-size-2xs, 0.75rem);
    text-align: center;
    max-width: 68px;
    line-height: 1.25;
  }
  .flow-step-line {
    margin: 0 4px -14px;
  }
  .modal-info-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  .modal-status-banner {
    flex-direction: column;
    align-items: flex-start;
    padding: 12px 14px;
  }
  .modal-status-dates {
    text-align: left !important;
  }
}
</style>
