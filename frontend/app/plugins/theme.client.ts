// ตั้งธีมตั้งแต่เริ่มโหลดฝั่ง client (อ่านจาก localStorage หรือค่าของระบบปฏิบัติการ)
export default defineNuxtPlugin(() => {
  useTheme().init()
})
