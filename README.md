# 🏬 e-Market Management System (ระบบจัดการแผงเช่าตลาดและสัญญาการเช่า)

ระบบบริหารจัดการแผงเช่าตลาดและสัญญาการเช่า พัฒนาด้วยสถาปัตยกรรม **Decoupled Modern Web Application**:
- **Frontend (FE)**: [Nuxt.js (Nuxt 4)](file:///Users/chinchutha/Project/frontend) + Clean Frontend Modular Architecture
- **Backend (BE)**: [Go (Gin Framework)](file:///Users/chinchutha/Project/backend) ตามมาตรฐาน **Clean Architecture**
- **Authentication**: มาตรฐาน **JSON Web Token (JWT Bearer Token)** ทุกการเรียกใช้งาน API
- **Database**: GORM (SQLite / PostgreSQL) พร้อมระบบ Auto-Migrate & Seeds

---

## 🚀 ลิงก์ระบบที่เปิดใช้งานอยู่ (Active Services)

| บริการ | URL | รายละเอียด |
| :--- | :--- | :--- |
| **Frontend Web App** | [http://localhost:3000](http://localhost:3000) | หน้าเว็บ Nuxt 4 (Dashboard, แผงเช่า, สัญญาเช่า, ตั้งค่า) |
| **Backend REST API** | [http://localhost:8000](http://localhost:8000) | Go Gin REST API Server |
| **Interactive API Docs (Swagger UI)** | [http://localhost:8000/swagger/index.html](http://localhost:8000/swagger/index.html) | เอกสารและทดสอบ REST API แบบ Interactive |

---

## 🔑 บัญชีเข้าสู่ระบบทดสอบ (Demo Accounts)

ระบบเตรียมข้อมูลผู้ใช้งานเริ่มต้นไว้ให้พร้อมใช้งานทันที:

| สิทธิ์ (Role) | Username | Password | ความสามารถ |
| :--- | :--- | :--- | :--- |
| **Admin (ผู้ดูแลระบบ)** | `admin` | `admin1234` | จัดการได้ครบทุกฟังก์ชัน (แผงเช่า, สัญญา, โซน, ผู้ใช้) |
| **Staff (เจ้าหน้าที่ตลาด)** | `staff` | `staff1234` | ดูข้อมูลแผงเช่า ทำสัญญาเช่า และอัปเดตสถานะ |

*(มีปุ่มคลิกกรอกข้อมูลอัตโนมัติบนหน้า Login)*

---

## 🌟 ฟีเจอร์หลักของระบบ (Key Features)

### 1. 🔐 ระบบยืนยันตัวตน (Authentication & JWT)
- ทุกการเชื่อมต่อไปยัง Backend ต้องส่ง Header `Authorization: Bearer <jwt_token>`
- รหัสผ่านเข้ารหัสด้วยมาตรฐาน **Bcrypt**
- รองรับการเข้าสู่ระบบผ่าน JSON Body (`POST /api/v1/auth/login`) และ OAuth2 Form (`POST /api/v1/auth/token`)

### 2. 🗺️ แผนที่ GIS สำหรับเลือกและจองแผงเช่า (Interactive GIS Map)
- พัฒนาด้วย **Leaflet.js** รองรับการแสดงผลพิกัดจริงระดับภูมิสารสนเทศ (GIS Geographic Map)
- สลับโหมดแผนที่ถนน (Street Map) และแผนที่ภาพถ่ายดาวเทียม (Satellite View) ได้
- วาดขอบเขตแผงค้าทุกแผงเป็น **GeoJSON Polygons** แยกสีตามสถานะแบบเรียลไทม์:
  - 🟢 **สีเขียว**: ว่างพร้อมจองเข้าดูแผง (Available)
  - 🟡 **สีส้ม**: มีผู้เช่าแล้ว (Occupied)
  - 🟣 **สีม่วง**: มีคำขอจองรอเจ้าหน้าที่อนุมัติ (Pending Approval)
  - 🔴 **สีแดง**: ปิดปรับปรุง (Maintenance)
- คลิกที่แผงเพื่อเปิด Drawer แสดงรายละเอียดแผง ขนาด ค่าเช่า และปุ่ม **"📍 จองดูแผงนี้"**

### 3. 📅 ระบบจองดูแผงและการอนุมัติ (Flow ขั้นตอนที่ 1 - 3)
- **ผู้เช่า**: เลือกแผงบนแผนที่ GIS &rarr; ระบุวัน/เวลาที่สะดวกเข้าชมหน้างาน &rarr; ส่งคำขอจอง
- **เจ้าหน้าที่ตลาด**: ตรวจสอบคำขอจอง &rarr; กด **"อนุมัติ"** วันนัดหมายเข้าชม หรือ **"ไม่อนุมัติ (พร้อมระบุเหตุผล)"**
- **ผู้เช่าเข้าชมแผงหน้างาน**: เมื่อดูแผงแล้ว สามารถกดเลือก:
  - **"✅ สนใจยื่นขอเช่า"**: ก้าวเข้าสู่การยื่นคำร้องขอเช่าแผง (`Lease Application`)
  - **"❌ ไม่สนใจ"**: ปลดแอกแผงคืนกลับสู่สถานะว่างทันที

### 4. 📋 การยื่นและอนุมัติคำร้องขอเช่าแผง (Flow ขั้นตอนที่ 4)
- บันทึกคำร้องขอเช่าระบุชื่อร้านค้า สินค้าที่จะขาย ระยะเวลาขอเช่า
- ระบบอนุมัติ 3 ระดับตามระเบียบตลาด:
  1. เจ้าหน้าที่ตรวจสอบคำร้อง &rarr; 2. ผอ.ตลาด อนุมัติ &rarr; 3. ผอ.สำนักตลาด อนุมัติ
- เมื่ออนุมัติครบทุกระดับ สามารถกดปุ่มคลิกเดียวเพื่อ **"จัดทำสัญญาเช่า (Generate Contract)"**

### 5. 📝 การจัดการสัญญาการเช่า (Rental Contracts)
- สร้างสัญญาเช่าใหม่ รหัสสัญญารันอัตโนมัติ เช่น `CNT-2026-0001`
- บันทึกข้อมูลผู้เช่าครบถ้วน: ชื่อ-นามสกุล, เลขประจำตัวประชาชน, เบอร์โทรศัพท์, ที่อยู่, ชื่อร้านค้า, ประเภทสินค้า
- ระบุระยะเวลาสัญญา, รอบการชำระเงิน (รายเดือน/รายวัน), อัตราค่าเช่า, เงินประกันสัญญา, เงินจ่ายล่วงหน้า
- **ซิงค์สถานะอัตโนมัติ**: เมื่อเปิดใช้งานสัญญา ระบบจะปรับสถานะแผงเป็น **"มีผู้เช่าแล้ว"** ทันที
- **ยกเลิก/สิ้นสุดสัญญาก่อนกำหนด**: ระบบจะปลดแอกแผงเช่ากลับมาเป็นสถานะ **"ว่าง"** ให้อัตโนมัติ
- **พิมพ์สัญญาฉบับทางการ (Print Agreement)**: มีหน้าต่างพรีวิวเอกสาร "หนังสือสัญญาเช่าพื้นที่แผงค้าในตลาด" พร้อมช่องลงลายมือชื่อผู้ให้เช่าและผู้เช่า สั่งพิมพ์หรือบันทึกเป็น PDF ได้ทันที

### 4. 📊 แดชบอร์ดและสถิติภาพรวม (Analytics Dashboard)
- แสดงตัวเลขสถิติสำคัญ (KPIs):
  - จำนวนแผงทั้งหมด, แผงว่าง, แผงมีผู้เช่า, แผงซ่อมบำรุง
  - อัตราการเช่าพื้นที่คิดเป็นร้อยละ (Occupancy Rate %)
  - จำนวนสัญญาที่เปิดใช้งานอยู่
  - ประมาณการรายรับค่าเช่าต่อเดือน และยอดเงินประกันสะสม
  - แจ้งเตือนสัญญาที่ใกล้หมดอายุภายใน 30 วัน
- กราฟแท่งความคืบหน้าอัตราการเช่าแยกตามโซนตลาด

---

## 🛠️ โครงสร้างโปรเจกต์ (Project Structure)

```
Project/
├── backend/                  # FastAPI Application
│   ├── app/
│   │   ├── core/             # JWT Auth, Security, Dependencies
│   │   ├── models/           # SQLAlchemy Models (User, Zone, Stall, Contract)
│   │   ├── routers/          # REST API Routers (/auth, /stalls, /contracts, /dashboard)
│   │   ├── schemas/          # Pydantic v2 Request/Response Schemas
│   │   ├── config.py         # App configuration & CORS settings
│   │   ├── database.py       # SQLite & SQLAlchemy engine
│   │   ├── main.py           # FastAPI entrypoint
│   │   └── seed.py           # Initial market demo dataset
│   ├── requirements.txt
│   └── test_api.py           # Automated test suite
│
├── frontend/                 # Nuxt.js (Nuxt 3) Application
│   ├── assets/css/main.css   # Vanilla CSS Design System (Glassmorphism & Tokens)
│   ├── composables/          # useAuth (JWT store), useApi (REST wrapper)
│   ├── layouts/              # default.vue (with sidebar), auth.vue
│   ├── pages/
│   │   ├── index.vue         # Overview Analytics Dashboard
│   │   ├── login.vue         # Authentication page
│   │   ├── settings.vue      # Zones & System Architecture
│   │   ├── contracts/        # Lease agreements & printable modal
│   │   └── stalls/           # Interactive stalls grid & management
│   ├── components/           # AppSidebar, AppHeader
│   ├── nuxt.config.ts
│   └── package.json
└── README.md
```

---

## 💻 วิธีการรันระบบด้วยตนเอง (Manual Start Commands)

### 1. เริ่มต้น Backend (FastAPI):
```bash
cd backend
# สร้าง Virtualenv และติดตั้ง dependencies (หากยังไม่ได้ติดตั้ง)
~/.local/bin/uv venv .venv --python 3.11
~/.local/bin/uv pip install -r requirements.txt

# รันเซิร์ฟเวอร์
.venv/bin/uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

### 2. เริ่มต้น Frontend (Nuxt 3):
```bash
cd frontend
export PATH="$HOME/.local/node/bin:$PATH"
export NUXT_TELEMETRY_DISABLED=1
npm run dev
```
เข้าใช้งานผ่านเบราว์เซอร์ที่: **http://localhost:82/e-market/**
# E-market
