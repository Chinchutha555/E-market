// สถานะการจองใช้ร่วมกันทั้งหน้ารายการและหน้ารายละเอียด เพื่อให้ชื่อและสีตรงกัน
export type BookingStatusTone = 'amber' | 'green' | 'blue' | 'purple' | 'sky' | 'gray'

export const getBookingStatusLabel = (status?: string | null): string => {
  switch (status) {
    case 'pending':
    case 'pending_approval': return 'รอเจ้าหน้าที่อนุมัติ'
    case 'pending_slot': return 'รอเลือกวันนัดชม'
    case 'approved': return 'นัดหมายสำเร็จ'
    case 'interested':
    case 'pending_decision': return 'รอผู้เช่ายื่นคำร้อง'
    case 'waiting_queue': return 'รอคิว'
    case 'lease_applied':
    case 'completed': return 'ยื่นคำร้องสำเร็จ'
    case 'rejected':
    case 'cancelled': return 'ยกเลิก'
    case 'not_interested': return 'สละสิทธิ์'
    case 'no_show': return 'ไม่มาตามนัด (ยกเลิก)'
    case 'terminated': return 'สิ้นสุดสัญญาเช่า'
    default: return 'นัดหมายสำเร็จ'
  }
}

export const getBookingStatusTone = (status?: string | null): BookingStatusTone => {
  switch (status) {
    case 'lease_applied':
    case 'completed': return 'purple'
    case 'approved': return 'green'
    case 'interested':
    case 'pending_decision': return 'blue'
    case 'pending_approval':
    case 'pending':
    case 'pending_slot': return 'amber'
    case 'waiting_queue': return 'sky'
    case 'cancelled':
    case 'rejected':
    case 'no_show':
    case 'not_interested':
    case 'terminated': return 'gray'
    default: return 'purple'
  }
}
