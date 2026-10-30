package seeds

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"emarket-api/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) error {
	hashPassword := func(pwd string) string {
		h, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		return string(h)
	}

	// 1. Seed Users
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		users := []models.User{
			{Username: "admin", Password: hashPassword("admin1234"), FullName: "ผู้ดูแลระบบตลาด กทม. (Admin)", Role: "admin", Email: "admin@emarket.bma.go.th", Phone: "02-123-4567", Status: "A"},
			{Username: "staff", Password: hashPassword("staff1234"), FullName: "เจ้าหน้าที่ตลาด (Staff)", Role: "staff", Email: "staff@emarket.bma.go.th", Phone: "02-123-4568", Status: "A"},
			{Username: "tenant", Password: hashPassword("tenant1234"), FullName: "สมชาย ใจดี", Role: "tenant", Email: "tenant@emarket.bma.go.th", Phone: "081-234-5678", Status: "A"},
			{Username: "tenant2", Password: hashPassword("tenant21234"), FullName: "นภาพร รัตนศิลป์", Role: "tenant", Email: "tenant2@emarket.bma.go.th", Phone: "089-111-2233", Status: "A"},
			{Username: "tenant3", Password: hashPassword("tenant31234"), FullName: "วิชัย เจริญทรัพย์", Role: "tenant", Email: "tenant3@emarket.bma.go.th", Phone: "086-778-9900", Status: "A"},
		}
		if err := db.Create(&users).Error; err != nil {
			return err
		}
		log.Println("Seeded users successfully")
	}

	// Ensure tenant accounts exist even if user table was already created
	for _, t := range []models.User{
		{Username: "tenant", Password: hashPassword("tenant1234"), FullName: "สมชาย ใจดี", Role: "tenant", Email: "tenant@emarket.bma.go.th", Phone: "081-234-5678", Status: "A"},
		{Username: "tenant2", Password: hashPassword("tenant21234"), FullName: "นภาพร รัตนศิลป์", Role: "tenant", Email: "tenant2@emarket.bma.go.th", Phone: "089-111-2233", Status: "A"},
		{Username: "tenant3", Password: hashPassword("tenant31234"), FullName: "วิชัย เจริญทรัพย์", Role: "tenant", Email: "tenant3@emarket.bma.go.th", Phone: "086-778-9900", Status: "A"},
	} {
		var existing models.User
		db.Where("username = ?", t.Username).Attrs(t).FirstOrCreate(&existing)
	}

	// Director demo accounts are only created in an explicitly configured local environment.
	if os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "test" {
		for _, director := range []models.User{
			{Username: "market_director_demo", Password: hashPassword("marketdemo1234"), FullName: "ผู้อำนวยการตลาด (ทดสอบ)", Role: "market_director", Status: "A"},
			{Username: "bureau_director_demo", Password: hashPassword("bureaudemo1234"), FullName: "ผู้อำนวยการสำนักตลาด (ทดสอบ)", Role: "bureau_director", Status: "A"},
		} {
			var existing models.User
			if err := db.Where("username = ?", director.Username).Attrs(director).FirstOrCreate(&existing).Error; err != nil {
				return err
			}
		}
	}
	// Director accounts for any environment (including production). Credentials come
	// from environment variables and fall back to the demo credentials shown on the login
	// page when unset, so the accounts always exist after deploy. Set the *_PASSWORD
	// variables in production to use private credentials instead.
	// Existing accounts are never modified, so passwords can be rotated in the database.
	for _, d := range []struct {
		userEnv, passEnv, nameEnv, defUser, defPass, defName, role string
	}{
		{"MARKET_DIRECTOR_USERNAME", "MARKET_DIRECTOR_PASSWORD", "MARKET_DIRECTOR_FULL_NAME", "market_director_demo", "marketdemo1234", "ผู้อำนวยการตลาด (ทดสอบ)", "market_director"},
		{"BUREAU_DIRECTOR_USERNAME", "BUREAU_DIRECTOR_PASSWORD", "BUREAU_DIRECTOR_FULL_NAME", "bureau_director_demo", "bureaudemo1234", "ผู้อำนวยการสำนักตลาด (ทดสอบ)", "bureau_director"},
	} {
		pwd := os.Getenv(d.passEnv)
		if pwd == "" {
			pwd = d.defPass
		}
		username := os.Getenv(d.userEnv)
		if username == "" {
			username = d.defUser
		}
		fullName := os.Getenv(d.nameEnv)
		if fullName == "" {
			fullName = d.defName
		}
		var existing models.User
		if err := db.Where("username = ?", username).
			Attrs(models.User{Username: username, Password: hashPassword(pwd), FullName: fullName, Role: d.role, Status: "A"}).
			FirstOrCreate(&existing).Error; err != nil {
			return err
		}
	}

	// 2. Seed 12 BMA Markets
	var marketCount int64
	db.Model(&models.Market{}).Count(&marketCount)
	if marketCount == 0 {
		markets := []models.Market{
			{
				ID: 1, Slug: "chatuchak", Code: "MKT-01",
				Name: "ตลาดนัดจตุจักร", NameEn: "Chatuchak Weekend Market",
				District: "เขตจตุจักร", Address: "587/10 ถนนกำแพงเพชร 2 แขวงจตุจักร เขตจตุจักร กรุงเทพฯ 10900",
				Latitude: 13.79996, Longitude: 100.55047, Lat: 13.79996, Lng: 100.55047,
				TotalStalls: 9522, AvailableStalls: 9522, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 120, StartingPriceMonthly: 1600,
				ImageUrl: "/images/markets/chatuchak.jpg", Category: "weekend",
				OperatingHours: "เสาร์ - อาทิตย์ 09:00 - 18:00 น.",
				Description:    "ศูนย์รวมการค้าปลีกกลางแจ้งที่ใหญ่ที่สุดในโลก สินค้าแฟชั่น ของแต่งบ้าน งานคราฟต์ อาหาร และต้นไม้",
				Status:         "open",
			},
			{
				ID: 2, Slug: "pracha-niwet-1", Code: "MKT-02",
				Name: "ตลาดประชานิเวศน์ 1", NameEn: "Pracha Niwet 1 Market",
				District: "เขตจตุจักร", Address: "ถนนเทศบาลสงเคราะห์ แขวงลาดยาว เขตจตุจักร กรุงเทพฯ 10900",
				Latitude: 13.8441, Longitude: 100.5489, Lat: 13.8441, Lng: 100.5489,
				TotalStalls: 528, AvailableStalls: 528, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 150, StartingPriceMonthly: 1800,
				ImageUrl: "/images/markets/pracha_niwet.jpg", Category: "fresh_market",
				OperatingHours: "ทุกวัน 05:00 - 19:00 น.",
				Description:    "ตลาดสดและอาหารสำเร็จรูปชื่อดังย่านประชานิเวศน์ สะอาด ปลอดภัย สินค้าคุณภาพดีเพื่อคนเมือง",
				Status:         "open",
			},
			{
				ID: 3, Slug: "min-buri", Code: "MKT-03",
				Name: "ตลาดมีนบุรี", NameEn: "Min Buri Market",
				District: "เขตมีนบุรี", Address: "ถนนสีหบุรานุกิจ แขวงมีนบุรี เขตมีนบุรี กรุงเทพฯ 10510",
				Latitude: 13.8139, Longitude: 100.7188, Lat: 13.8139, Lng: 100.7188,
				TotalStalls: 2029, AvailableStalls: 2029, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 100, StartingPriceMonthly: 1500,
				ImageUrl: "/images/markets/min_buri.jpg", Category: "metro",
				OperatingHours: "ทุกวัน 04:00 - 20:00 น.",
				Description:    "ศูนย์กลางตลาดค้าส่ง-ปลีกฝั่งกรุงเทพฯ ตะวันออก เชื่อมต่อการเดินทางด้วยรถไฟฟ้าสายสีชมพู",
				Status:         "open",
			},
			{
				ID: 4, Slug: "singhanat", Code: "MKT-04",
				Name: "ตลาดสิงหา", NameEn: "Singhanat Market",
				District: "เขตคลองเตย", Address: "ซอยสุขุมวิท 22 แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
				Latitude: 13.7285, Longitude: 100.5642, Lat: 13.7285, Lng: 100.5642,
				TotalStalls: 91, AvailableStalls: 91, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 180, StartingPriceMonthly: 2500,
				ImageUrl: "/images/markets/singhanat.jpg", Category: "metro",
				OperatingHours: "จันทร์ - ศุกร์ 06:00 - 15:00 น.",
				Description:    "ตลาดชุมชนใจกลางสุขุมวิท 22 รองรับกลุ่มคนทำงานออฟฟิศ นักท่องเที่ยว และผู้อยู่อาศัยในพื้นที่",
				Status:         "open",
			},
			{
				ID: 5, Slug: "thewarat", Code: "MKT-05",
				Name: "ตลาดเทวราช", NameEn: "Thewarat Market",
				District: "เขตดุสิต", Address: "ถนนพิษณุโลก แขวงสี่แยกมหานาค เขตดุสิต กรุงเทพฯ 10300",
				Latitude: 13.7712, Longitude: 100.5050, Lat: 13.7712, Lng: 100.5050,
				TotalStalls: 475, AvailableStalls: 475, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 120, StartingPriceMonthly: 1800,
				ImageUrl: "/images/markets/thewarat.jpg", Category: "fresh_market",
				OperatingHours: "ทุกวัน 05:00 - 18:00 น.",
				Description:    "ตลาดสดเก่าแก่อัตลักษณ์ดั้งเดิมย่านดุสิต ใกล้สะพานเทเวศร์ แหล่งวัตถุดิบอาหารและดอกไม้สด",
				Status:         "open",
			},
			{
				ID: 6, Slug: "rat-burana", Code: "MKT-06",
				Name: "ตลาดราษฎร์บูรณะ", NameEn: "Rat Burana Market",
				District: "เขตราษฎร์บูรณะ", Address: "ถนนราษฎร์บูรณะ แขวงราษฎร์บูรณะ เขตราษฎร์บูรณะ กรุงเทพฯ 10140",
				Latitude: 13.6819, Longitude: 100.5036, Lat: 13.6819, Lng: 100.5036,
				TotalStalls: 140, AvailableStalls: 140, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 100, StartingPriceMonthly: 1500,
				ImageUrl: "/images/markets/rat_burana.jpg", Category: "fresh_market",
				OperatingHours: "ทุกวัน 05:00 - 18:00 น.",
				Description:    "ตลาดริมแม่น้ำเจ้าพระยาฝั่งธนบุรี บริการสินค้าอุปโภคบริโภคสดใหม่เพื่อชุมชนราษฎร์บูรณะ",
				Status:         "open",
			},
			{
				ID: 7, Slug: "bang-kapi", Code: "MKT-07",
				Name: "ตลาดบางกะปิ", NameEn: "Bang Kapi Market",
				District: "เขตบางกะปิ", Address: "ถนนลาดพร้าว แขวงคลองจั่น เขตบางกะปิ กรุงเทพฯ 10240",
				Latitude: 13.7661, Longitude: 100.6441, Lat: 13.7661, Lng: 100.6441,
				TotalStalls: 710, AvailableStalls: 710, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 130, StartingPriceMonthly: 1900,
				ImageUrl: "/images/markets/bang_kapi.jpg", Category: "metro",
				OperatingHours: "ทุกวัน 04:00 - 21:00 น.",
				Description:    "ศูนย์กลางการค้าและคมนาคมย่านรามคำแหง-ลาดพร้าว เชื่อมต่อรถไฟฟ้าสายสีเหลืองและท่าเรือคลองแสนแสบ",
				Status:         "open",
			},
			{
				ID: 8, Slug: "wongwian-lek", Code: "MKT-08",
				Name: "ตลาดวงเวียนเล็ก", NameEn: "Wongwian Lek Market",
				District: "เขตคลองสาน", Address: "ถนนประชาธิปก แขวงสมเด็จเจ้าพระยา เขตคลองสาน กรุงเทพฯ 10600",
				Latitude: 13.7323, Longitude: 100.4952, Lat: 13.7323, Lng: 100.4952,
				TotalStalls: 648, AvailableStalls: 648, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 150, StartingPriceMonthly: 2100,
				ImageUrl: "/images/markets/wongwian_lek.jpg", Category: "night",
				OperatingHours: "ทุกวัน 16:00 - 23:00 น.",
				Description:    "ตลาดของอร่อยยามค่ำคืนฝั่งธนบุรี ใกล้วงเวียนใหญ่ แหล่งรวมสตรีทฟู้ดในตำนานและขนมหวานโบราณ",
				Status:         "open",
			},
			{
				ID: 9, Slug: "nong-chok", Code: "MKT-09",
				Name: "ตลาดหนองจอก", NameEn: "Nong Chok Market",
				District: "เขตหนองจอก", Address: "ถนนมิตรไมตรี แขวงหนองจอก เขตหนองจอก กรุงเทพฯ 10530",
				Latitude: 13.8564, Longitude: 100.8641, Lat: 13.8564, Lng: 100.8641,
				TotalStalls: 268, AvailableStalls: 268, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 80, StartingPriceMonthly: 1200,
				ImageUrl: "/images/markets/nong_chok.jpg", Category: "fresh_market",
				OperatingHours: "ทุกวัน 05:00 - 18:00 น.",
				Description:    "ตลาดชุมชนชานเมืองกรุงเทพฯ ศูนย์กลางอาหารฮาลาลและผลิตผลเกษตรอินทรีย์คุณภาพดี",
				Status:         "open",
			},
			{
				ID: 10, Slug: "ratchadaphisek", Code: "MKT-10",
				Name: "ตลาดรัชดาภิเษก", NameEn: "Ratchadaphisek Market",
				District: "เขตธนบุรี", Address: "ถนนรัชดาภิเษก แขวงดาวคะนอง เขตธนบุรี กรุงเทพฯ 10600",
				Latitude: 13.7088, Longitude: 100.4823, Lat: 13.7088, Lng: 100.4823,
				TotalStalls: 128, AvailableStalls: 128, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 140, StartingPriceMonthly: 2000,
				ImageUrl: "/images/markets/ratchada.jpg", Category: "metro",
				OperatingHours: "ทุกวัน 06:00 - 20:00 น.",
				Description:    "ตลาดสดและสตรีทฟู้ดย่านตลาดพลู แหล่งรวมของอร่อยสำหรับคนทำงานและครอบครัว",
				Status:         "open",
			},
			{
				ID: 11, Slug: "thonburi-sanamluang-2", Code: "MKT-11",
				Name: "ตลาดธนบุรี (สนามหลวง 2)", NameEn: "Thonburi Market (Sanam Luang 2)",
				District: "เขตทวีวัฒนา", Address: "ถนนเลียบคลองทวีวัฒนา แขวงทวีวัฒนา เขตทวีวัฒนา กรุงเทพฯ 10170",
				Latitude: 13.7479, Longitude: 100.3542, Lat: 13.7479, Lng: 100.3542,
				TotalStalls: 4263, AvailableStalls: 4263, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 100, StartingPriceMonthly: 1500,
				ImageUrl: "/images/markets/thonburi.jpg", Category: "popular",
				OperatingHours: "ทุกวัน 08:00 - 18:00 น.",
				Description:    "ตลาดนัดสุดสัปดาห์ฝั่งธนบุรี บรรยากาศร่มรื่น ศูนย์รวมสัตว์เลี้ยง ต้นไม้ ของสะสม และงานฝีมือ",
				Status:         "open",
			},
			{
				ID: 12, Slug: "bang-khae-phirom", Code: "MKT-12",
				Name: "ตลาดบางแคภิรมย์", NameEn: "Bang Khae Phirom Market",
				District: "เขตบางแค", Address: "ซอยเพชรเกษม 69 แขวงหลักสอง เขตบางแค กรุงเทพฯ 10160",
				Latitude: 13.6896, Longitude: 100.3807, Lat: 13.6896, Lng: 100.3807,
				TotalStalls: 624, AvailableStalls: 624, OccupiedStalls: 0, ReservedStalls: 0,
				StartingPriceDaily: 110, StartingPriceMonthly: 1600,
				ImageUrl: "/images/markets/bang_khae.jpg", Category: "popular",
				OperatingHours: "ทุกวัน 06:00 - 20:00 น.",
				Description:    "ตลาดริมสวนสุขภาพบางแคภิรมย์ บรรยากาศโปร่งสบาย เหมาะกับสินค้าเพื่อสุขภาพและอาหารครอบครัว",
				Status:         "open",
			},
		}

		if err := db.Create(&markets).Error; err != nil {
			return err
		}
		log.Println("Seeded 12 markets successfully")
	}

	// Ensure staff account has market assignments so ScopeMarketQuery and CanManageMarket work
	var staffUser models.User
	if err := db.Where("username = ?", "staff").First(&staffUser).Error; err == nil {
		var assignmentCount int64
		db.Model(&models.MarketStaffAssignment{}).Where("user_id = ?", staffUser.ID).Count(&assignmentCount)
		if assignmentCount == 0 {
			var mktIDs []int
			db.Model(&models.Market{}).Pluck("id", &mktIDs)
			var assignments []models.MarketStaffAssignment
			for _, mID := range mktIDs {
				assignments = append(assignments, models.MarketStaffAssignment{
					UserID:   staffUser.ID,
					MarketID: mID,
				})
			}
			if len(assignments) > 0 {
				_ = db.Create(&assignments).Error
				log.Println("Seeded staff market assignments successfully")
			}
		}
	}

	if marketCount == 0 {
		var markets []models.Market
		db.Find(&markets)
		var zones []models.Zone
		for _, m := range markets {
			zones = append(zones,
				models.Zone{
					MarketID:    m.ID,
					Name:        "โครงการ 1 (แฟชั่นและเครื่องแต่งกาย)",
					Code:        fmt.Sprintf("Z-%02d-1", m.ID),
					Color:       "#3b82f6",
					TotalStalls: 50,
					Status:      "A",
				},
				models.Zone{
					MarketID:    m.ID,
					Name:        "โครงการ 2 (อาหารและเครื่องดื่ม)",
					Code:        fmt.Sprintf("Z-%02d-2", m.ID),
					Color:       "#10b981",
					TotalStalls: 40,
					Status:      "A",
				},
				models.Zone{
					MarketID:    m.ID,
					Name:        "โครงการ 3 (ของตกแต่งและเบ็ดเตล็ด)",
					Code:        fmt.Sprintf("Z-%02d-3", m.ID),
					Color:       "#f59e0b",
					TotalStalls: 30,
					Status:      "A",
				},
			)
		}
		if err := db.Create(&zones).Error; err != nil {
			return err
		}
		log.Println("Seeded market zones successfully")

		// 4. Seed Representative Stalls for each market
		var stalls []models.Stall
		stallPrefixes := []string{"A", "B", "C", "D"}
		for _, m := range markets {
			zoneID := (m.ID-1)*3 + 1
			for pIdx, prefix := range stallPrefixes {
				for num := 1; num <= 5; num++ {
					ident := fmt.Sprintf("%s-%02d", prefix, num)
					if m.ID > 1 {
						ident = fmt.Sprintf("M%d-%s%02d", m.ID, prefix, num)
					}
					status := "available"
					if num == 3 {
						status = "reserved"
					} else if num == 4 || num == 5 {
						status = "occupied"
					}

					stalls = append(stalls, models.Stall{
						MarketID:    m.ID,
						MarketSlug:  m.Slug,
						MarketName:  m.Name,
						ZoneID:      zoneID + (pIdx % 3),
						ZoneName:    fmt.Sprintf("โครงการ %d", (pIdx%3)+1),
						ProjectName: fmt.Sprintf("โครงการ %d", (pIdx%3)+1),
						Ident:       ident,
						Name:        fmt.Sprintf("แผงค้า %s (%s)", ident, m.Name),
						Width:       3.0,
						Depth:       3.0,
						AreaSqm:     9.0,
						DailyRent:   m.StartingPriceDaily,
						MonthlyRent: m.StartingPriceMonthly,
						StallStatus: status,
						Floor:       "1",
						Status:      "A",
					})
				}
			}
		}
		if err := db.Create(&stalls).Error; err != nil {
			return err
		}
		log.Println("Seeded stalls successfully")
	}

	// Each market officer has an explicit, single-market scope. This seed is
	// idempotent so existing databases receive the accounts and assignments too.
	var marketsForStaff []models.Market
	if err := db.Where("status != 'D'").Order("id ASC").Limit(12).Find(&marketsForStaff).Error; err != nil {
		return err
	}
	for _, market := range marketsForStaff {
		username := fmt.Sprintf("staff%02d", market.ID)
		staff := models.User{Username: username}
		defaults := models.User{
			Username: username,
			Password: hashPassword("staff1234"),
			FullName: fmt.Sprintf("เจ้าหน้าที่%s", market.Name),
			Role:     "staff",
			Email:    fmt.Sprintf("%s@emarket.bma.go.th", username),
			Phone:    fmt.Sprintf("02-100-%04d", market.ID),
			Status:   "A",
		}
		if err := db.Where("username = ?", username).Attrs(defaults).FirstOrCreate(&staff).Error; err != nil {
			return err
		}
		assignment := models.MarketStaffAssignment{UserID: staff.ID, MarketID: market.ID}
		if err := db.Where("user_id = ? AND market_id = ?", staff.ID, market.ID).
			FirstOrCreate(&assignment).Error; err != nil {
			return err
		}
	}
	// Scope the market director demo accounts to markets so they have assignments to review applications.
	if len(marketsForStaff) > 0 {
		var directors []models.User
		if err := db.Where("role = ?", "market_director").Find(&directors).Error; err == nil {
			for _, director := range directors {
				for _, m := range marketsForStaff {
					assignment := models.MarketStaffAssignment{UserID: director.ID, MarketID: m.ID}
					if err := db.Where("user_id = ? AND market_id = ?", director.ID, m.ID).FirstOrCreate(&assignment).Error; err != nil {
						return err
					}
				}
			}
		}
	}
	// Keep the legacy demo account usable while treating it as a normal
	// single-market officer rather than a global staff account.
	if len(marketsForStaff) > 0 {
		var legacyStaff models.User
		if db.Where("username = ? AND role = ?", "staff", "staff").First(&legacyStaff).Error == nil {
			assignment := models.MarketStaffAssignment{UserID: legacyStaff.ID, MarketID: marketsForStaff[0].ID}
			if err := db.Where("user_id = ? AND market_id = ?", legacyStaff.ID, marketsForStaff[0].ID).FirstOrCreate(&assignment).Error; err != nil {
				return err
			}
		}
	}

	// Backfill scope for records created before market_id was introduced.
	db.Exec("UPDATE application SET market_id = (SELECT market_id FROM stall WHERE stall.id = application.stall_id OR stall.ident = application.stall_ident LIMIT 1) WHERE market_id IS NULL")
	db.Exec("UPDATE application SET market_name = (SELECT name FROM market WHERE market.id = application.market_id LIMIT 1) WHERE (market_name IS NULL OR market_name = '') AND market_id IS NOT NULL")
	db.Exec("UPDATE application SET details = review_note WHERE (details IS NULL OR details = '') AND review_note IS NOT NULL AND review_note != ''")
	db.Exec("UPDATE contract SET market_id = (SELECT market_id FROM stall WHERE stall.id = contract.stall_id OR stall.ident = contract.stall_ident LIMIT 1) WHERE market_id IS NULL")
	db.Exec("UPDATE contract SET tenant_user_id = (SELECT created_by FROM application WHERE application.id = contract.application_id LIMIT 1) WHERE tenant_user_id IS NULL AND application_id IS NOT NULL")
	db.Exec("UPDATE contract SET tenant_user_id = (SELECT created_by FROM booking WHERE booking.stall_ident = contract.stall_ident AND created_by IS NOT NULL ORDER BY booking.id DESC LIMIT 1) WHERE tenant_user_id IS NULL")

	// Ensure accurate GIS totals and starting rates across all markets
	gisTotals := map[string]int{
		"chatuchak":             9522,
		"pracha-niwet-1":        528,
		"min-buri":              2734,
		"singhanat":             91,
		"thewarat":              475,
		"rat-burana":            140,
		"bang-kapi":             710,
		"wongwian-lek":          648,
		"nong-chok":             268,
		"ratchadaphisek":        128,
		"thonburi-sanamluang-2": 4263,
		"bang-khae-phirom":      624,
	}
	startingMonthly := map[string]float64{
		"chatuchak":             1600,
		"pracha-niwet-1":        1800,
		"min-buri":              1500,
		"singhanat":             2500,
		"thewarat":              1800,
		"rat-burana":            1500,
		"bang-kapi":             1900,
		"wongwian-lek":          2100,
		"nong-chok":             1200,
		"ratchadaphisek":        2000,
		"thonburi-sanamluang-2": 1500,
		"bang-khae-phirom":      1600,
	}
	for slug, total := range gisTotals {
		price := startingMonthly[slug]
		db.Model(&models.Market{}).Where("slug = ?", slug).Updates(map[string]interface{}{
			"total_stalls":           total,
			"starting_price_monthly": price,
		})
	}

	// Seed initial market_gis_config from backend/data/markets/<slug>.json if not yet present
	for _, market := range marketsForStaff {
		var existingGIS models.MarketGISConfig
		err := db.Where("market_id = ?", market.ID).First(&existingGIS).Error
		if err == nil && len(existingGIS.ProjectsJSON) > 2 {
			continue
		}
		jsonCandidates := []string{
			filepath.Join("data", "markets", market.Slug+".json"),
			filepath.Join("backend", "data", "markets", market.Slug+".json"),
		}
		for _, p := range jsonCandidates {
			rawBytes, rErr := os.ReadFile(p)
			if rErr != nil {
				continue
			}
			var doc map[string]interface{}
			if json.Unmarshal(rawBytes, &doc) != nil {
				continue
			}
			if projsRaw, ok := doc["projects"]; ok {
				if projsBytes, mErr := json.Marshal(projsRaw); mErr == nil && len(projsBytes) > 2 {
					facilitiesJSON := "[]"
					if facilitiesRaw, facilitiesOK := doc["facilities"]; facilitiesOK {
						if facilitiesBytes, facilitiesErr := json.Marshal(facilitiesRaw); facilitiesErr == nil {
							facilitiesJSON = string(facilitiesBytes)
						}
					}
					cfg := models.MarketGISConfig{
						MarketID:       market.ID,
						ProjectsJSON:   string(projsBytes),
						FacilitiesJSON: facilitiesJSON,
						Published:      true,
					}
					_ = db.Where(models.MarketGISConfig{MarketID: market.ID}).
						Assign(map[string]interface{}{
							"projects_json":   cfg.ProjectsJSON,
							"facilities_json": cfg.FacilitiesJSON,
							"published":       true,
						}).
						FirstOrCreate(&cfg).Error
					break
				}
			}
		}
	}

	// Normalize legacy booking queue_number and waiting_queue / pending_slot states
	var activeBookings []models.Booking
	if err := db.Where("status = 'A' AND booking_status NOT IN ('rejected', 'cancelled', 'not_interested', 'no_show', 'terminated')").
		Order("id ASC").Find(&activeBookings).Error; err == nil {
		holderSeen := map[string]bool{}
		queueIndex := map[string]int{}
		for _, bk := range activeBookings {
			stall := bk.StallIdent
			wasLegacyQueue := bk.BookingStatus == "waiting_queue"
			newStatus := bk.BookingStatus
			if newStatus == "waiting_queue" || newStatus == "pending_slot" {
				newStatus = "pending_approval"
			}
			newQueueNum := bk.QueueNumber
			if !wasLegacyQueue && !holderSeen[stall] {
				holderSeen[stall] = true
				newQueueNum = 0
			} else {
				queueIndex[stall]++
				newQueueNum = queueIndex[stall]
			}
			updates := map[string]interface{}{}
			if newStatus != bk.BookingStatus {
				updates["booking_status"] = newStatus
			}
			if newQueueNum != bk.QueueNumber {
				updates["queue_number"] = newQueueNum
			}
			if bk.VisitDate == "" {
				updates["visit_date"] = "2026-10-26"
			}
			if bk.VisitTime == "" {
				updates["visit_time"] = "10:00 - 11:00"
			}
			if len(updates) > 0 {
				db.Model(&models.Booking{}).Where("id = ?", bk.ID).Updates(updates)
			}
		}
	}

	// 10. Seed Application Types (ประเภทคำร้อง / เรื่องที่ขอ)
	var appTypeCount int64
	db.Model(&models.ApplicationType{}).Count(&appTypeCount)
	if appTypeCount == 0 {
		defaultTypes := []models.ApplicationType{
			{Title: "ขออนุญาตปรับปรุง/ต่อเติมแผงค้า", Code: "stall_renovation", Description: "ขออนุญาตดัดแปลง ตกแต่ง หรือติดตั้งสิ่งอำนวยความสะดวกในแผงค้า", IsActive: true, SortOrder: 1},
			{Title: "ขอเปลี่ยน/เพิ่มประเภทสินค้าที่จำหน่าย", Code: "product_change", Description: "ขอเปลี่ยนแปลงหรือเพิ่มหมวดหมู่สินค้าที่วางจำหน่ายในแผงค้า", IsActive: true, SortOrder: 2},
			{Title: "ขอออกหนังสือรับรองการเช่าแผงค้า", Code: "lease_certificate", Description: "ขอเอกสารหนังสือรับรองสิทธิการเช่าแผงเพื่อนำไปประกอบธุรกรรม", IsActive: true, SortOrder: 3},
			{Title: "ขอผ่อนผันเวลาเปิด-ปิดแผงค้า", Code: "hours_waiver", Description: "ขอเปิดหรือปิดร้านนอกเวลาทำการปกติของตลาด", IsActive: true, SortOrder: 4},
			{Title: "แจ้งปัญหา / ร้องเรียนการใช้งานพื้นที่ตลาด", Code: "complaint_maintenance", Description: "แจ้งซ่อมแซมสิ่งอำนวยความสะดวก ไฟฟ้า น้ำประปา หรือข้อร้องเรียนทั่วไป", IsActive: true, SortOrder: 5},
			{Title: "ขอโอนสิทธิ์การเช่าแผงค้า", Code: "lease_transfer", Description: "ขอเปลี่ยนตัวผู้เช่าหรือโอนสิทธิ์การเช่าแผงค้าให้ผู้อื่น", IsActive: true, SortOrder: 6},
			{Title: "ขอยกเลิกสัญญาเช่าก่อนกำหนด", Code: "early_termination", Description: "ขอยกเลิกสัญญาเช่าแผงค้าและส่งมอบพื้นที่คืนก่อนสิ้นสุดสัญญา", IsActive: true, SortOrder: 7},
			{Title: "เรื่องอื่นๆ", Code: "other", Description: "คำร้องหรือเรื่องติดต่ออื่นๆ นอกเหนือจากหัวข้อข้างต้น", IsActive: true, SortOrder: 99},
		}
		if err := db.Create(&defaultTypes).Error; err != nil {
			return err
		}
		log.Println("Seeded default application types successfully")
	}

	// 11. Seed Contracts
	var contractCount int64
	db.Model(&models.Contract{}).Count(&contractCount)
	if contractCount == 0 {
		mID := 1
		contracts := []models.Contract{
			{
				ContractNumber: "CNT-2026-JJ-16-232-01", StallIdent: "JJ-16-232", MarketID: &mID,
				TenantName: "คุณสมศักดิ์ มั่งคั่ง", TenantPhone: "081-889-9221", TenantIDCard: "3-1006-00421-44-1",
				ShopName: "ร้านค้าทั่วไป", StartDate: "2026-01-01", EndDate: "2026-12-31",
				MonthlyRental: 2000, SecurityDeposit: 4000, ContractState: "active", Status: "A",
			},
			{
				ContractNumber: "CNT-2026-JJ-17-392-01", StallIdent: "JJ-17-392", MarketID: &mID,
				TenantName: "สมชาย ใจดี", TenantPhone: "081-234-5678", TenantIDCard: "1-1004-99882-12-1",
				ShopName: "เสื้อผ้าแฟชั่นและคราฟต์", StartDate: "2026-01-01", EndDate: "2026-12-31",
				MonthlyRental: 1800, SecurityDeposit: 3600, ContractState: "active", Status: "A",
			},
			{
				ContractNumber: "CNT-2026-JJ-17-434-01", StallIdent: "JJ-17-434", MarketID: &mID,
				TenantName: "วิชัย เจริญทรัพย์", TenantPhone: "086-778-9900", TenantIDCard: "2-1005-44321-77-2",
				ShopName: "ร้านกระเป๋าและเครื่องหนัง", StartDate: "2026-02-01", EndDate: "2027-01-31",
				MonthlyRental: 2200, SecurityDeposit: 4400, ContractState: "active", Status: "A",
			},
			{
				ContractNumber: "CNT-2026-JJ-19-267-01", StallIdent: "JJ-19-267", MarketID: &mID,
				TenantName: "นภาพร รัตนศิลป์", TenantPhone: "089-111-2233", TenantIDCard: "1-1004-88421-13-9",
				ShopName: "ร้านเครื่องเงินและจิวเวลรี่", StartDate: "2026-03-01", EndDate: "2027-02-28",
				MonthlyRental: 2500, SecurityDeposit: 5000, ContractState: "active", Status: "A",
			},
			{
				ContractNumber: "CNT-2026-JJ-19-260-01", StallIdent: "JJ-19-260", MarketID: &mID,
				TenantName: "คุณกิตติศักดิ์ เจริญพร", TenantPhone: "085-334-1122", TenantIDCard: "3-1007-88214-99-0",
				ShopName: "ของแต่งบ้านวินเทจ", StartDate: "2026-01-01", EndDate: "2026-12-31",
				MonthlyRental: 2000, SecurityDeposit: 4000, ContractState: "active", Status: "A",
			},
			{
				ContractNumber: "CNT-2025-JJ-16-232-02", StallIdent: "JJ-16-232", MarketID: &mID,
				TenantName: "คุณกมลชนก แสงทอง", TenantPhone: "089-245-7810", TenantIDCard: "1-1008-33214-55-6",
				ShopName: "ร้านเครื่องดื่มสมุนไพร", StartDate: "2025-01-01", EndDate: "2025-12-31",
				MonthlyRental: 1800, SecurityDeposit: 3600, ContractState: "expired", Status: "A",
			},
			{
				ContractNumber: "CNT-2025-JJ-17-392-02", StallIdent: "JJ-17-392", MarketID: &mID,
				TenantName: "คุณประวิทย์ วัฒนกิจ", TenantPhone: "086-441-9082", TenantIDCard: "1-1009-44123-66-7",
				ShopName: "ของเล่นและของสะสม", StartDate: "2025-01-01", EndDate: "2025-12-31",
				MonthlyRental: 1600, SecurityDeposit: 3200, ContractState: "expired", Status: "A",
			},
		}
		if err := db.Create(&contracts).Error; err != nil {
			return err
		}
		log.Println("Seeded default contracts successfully")
	}

	// 12. Seed Tenant Debt
	var debtCount int64
	db.Model(&models.Debt{}).Count(&debtCount)
	if debtCount == 0 {
		debts := []models.Debt{
			{TenantID: 3, MarketID: 1, Description: "ค้างชำระค่าเช่าแผง JJ-16-232 งวดตุลาคม 2569", AmountSatang: 200000, Status: "outstanding", CreatedBy: 1},
			{TenantID: 4, MarketID: 1, Description: "ค้างชำระค่าน้ำ-ค่าไฟแผง JJ-19-267 ประจำเดือนกันยายน 2569", AmountSatang: 85000, Status: "outstanding", CreatedBy: 1},
			{TenantID: 5, MarketID: 1, Description: "ค้างชำระค่าธรรมเนียมขยะมูลฝอยประจำไตรมาส 3", AmountSatang: 45000, Status: "cleared", CreatedBy: 1},
		}
		if err := db.Create(&debts).Error; err != nil {
			return err
		}
		log.Println("Seeded default tenant debts successfully")
	}

	// 13. Seed Additional Bookings & Applications
	var appCount int64
	db.Model(&models.Application{}).Count(&appCount)
	if appCount <= 4 {
		mID := 1
		extraApps := []models.Application{
			{
				ApplicationNo: "APP-20261028-101234", StallIdent: "JJ-16-230", MarketID: &mID, MarketName: "ตลาดนัดจตุจักร",
				ApplicantName: "คุณณภัทร วงศ์สว่าง", ApplicantPhone: "089-555-4321", ApplicantIDCard: "1-1002-33441-22-3",
				ShopName: "ร้านเซรามิกและของแต่งบ้าน", ProductCategory: "ของแต่งบ้าน", LeaseTermMonths: 12, MonthlyRate: 2000,
				ApplicationState: "pending_officer", Status: "A", Subject: "ขอเช่าพื้นที่แผงค้าใหม่",
			},
			{
				ApplicationNo: "APP-20261028-202345", StallIdent: "JJ-16-231", MarketID: &mID, MarketName: "ตลาดนัดจตุจักร",
				ApplicantName: "คุณจิราพร บุญมี", ApplicantPhone: "084-222-3344", ApplicantIDCard: "1-1003-44552-33-4",
				ShopName: "ร้านขนมไทยโบราณ", ProductCategory: "อาหารและเครื่องดื่ม", LeaseTermMonths: 12, MonthlyRate: 1800,
				ApplicationState: "submitted", Status: "A", Subject: "ขอเช่าพื้นที่แผงค้าใหม่",
			},
			{
				ApplicationNo: "APP-20261028-303456", StallIdent: "JJ-17-390", MarketID: &mID, MarketName: "ตลาดนัดจตุจักร",
				ApplicantName: "คุณศิริชัย แซ่ตั้ง", ApplicantPhone: "087-999-8877", ApplicantIDCard: "1-1004-55663-44-5",
				ShopName: "ร้านอุปกรณ์แคมป์ปิ้ง", ProductCategory: "ทั่วไป", LeaseTermMonths: 12, MonthlyRate: 2200,
				ApplicationState: "approved", Status: "A", Subject: "ขอเช่าพื้นที่แผงค้าใหม่",
			},
		}
		_ = db.Create(&extraApps).Error
	}

	var totalBookings int64
	db.Model(&models.Booking{}).Count(&totalBookings)
	if totalBookings <= 6 {
		mID := 1
		extraBookings := []models.Booking{
			{
				BookingNo: "BK-20261028-091100", StallIdent: "JJ-16-230", MarketID: &mID, MarketName: "ตลาดนัดจตุจักร",
				ApplicantName: "คุณณภัทร วงศ์สว่าง", ApplicantPhone: "089-555-4321", BusinessType: "ของแต่งบ้าน",
				VisitDate: "2026-11-02", VisitTime: "10:00 - 11:00", BookingStatus: "approved", Status: "A",
			},
			{
				BookingNo: "BK-20261028-112200", StallIdent: "JJ-16-231", MarketID: &mID, MarketName: "ตลาดนัดจตุจักร",
				ApplicantName: "คุณจิราพร บุญมี", ApplicantPhone: "084-222-3344", BusinessType: "อาหารและเครื่องดื่ม",
				VisitDate: "2026-11-02", VisitTime: "14:00 - 15:00", BookingStatus: "pending_approval", Status: "A",
			},
			{
				BookingNo: "BK-20261028-143300", StallIdent: "JJ-17-390", MarketID: &mID, MarketName: "ตลาดนัดจตุจักร",
				ApplicantName: "คุณศิริชัย แซ่ตั้ง", ApplicantPhone: "087-999-8877", BusinessType: "ทั่วไป",
				VisitDate: "2026-11-03", VisitTime: "11:00 - 12:00", BookingStatus: "lease_applied", Status: "A",
			},
		}
		_ = db.Create(&extraBookings).Error
	}

	return nil
}
