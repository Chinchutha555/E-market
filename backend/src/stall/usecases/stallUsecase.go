package usecases

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"emarket-api/logs"
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/stall/domains"
	"gorm.io/gorm"
)

type useCase struct {
	repo domains.Repository
}

func NewUseCase(repo domains.Repository) domains.UseCase {
	return &useCase{repo: repo}
}

func toStallResponse(s *models.Stall) responses.StallResponse {
	return responses.StallResponse{
		ID:                 s.ID,
		MarketID:           s.MarketID,
		MarketSlug:         s.MarketSlug,
		MarketName:         s.MarketName,
		ZoneID:             s.ZoneID,
		ZoneName:           s.ZoneName,
		ProjectName:        s.ProjectName,
		Ident:              s.Ident,
		StallIdent:         s.Ident,
		StallNumber:        s.Ident,
		Name:               s.Name,
		Description:        s.Description,
		Category:           s.Category,
		Width:              s.Width,
		Depth:              s.Depth,
		AreaSqm:            s.AreaSqm,
		DailyRent:          s.DailyRent,
		MonthlyRent:        s.MonthlyRent,
		Status:             s.StallStatus,
		StallStatus:        s.StallStatus,
		Floor:              s.Floor,
		CoordinatesGeoJSON: s.CoordinatesGeoJSON,
		Published:          s.Published,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
}

func maskPDPAName(name string) string {
	clean := strings.TrimSpace(name)
	for _, prefix := range []string{"นางสาว", "น.ส.", "นาย", "นาง", "คุณ"} {
		if strings.HasPrefix(clean, prefix) {
			clean = strings.TrimSpace(strings.TrimPrefix(clean, prefix))
			break
		}
	}
	parts := strings.Fields(clean)
	if len(parts) == 0 {
		return "คุณ ***"
	}
	firstName := parts[0]
	if len(parts) > 1 {
		lastRunes := []rune(parts[1])
		if len(lastRunes) > 0 {
			return fmt.Sprintf("คุณ%s %s.***", firstName, string(lastRunes[0]))
		}
	}
	firstRunes := []rune(firstName)
	if len(firstRunes) > 3 {
		return fmt.Sprintf("คุณ%s***", string(firstRunes[:3]))
	}
	return fmt.Sprintf("คุณ%s***", firstName)
}

func maskPDPAphone(phone string) string {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(phone, "-", ""), " ", "")
	if len(cleaned) == 10 {
		return fmt.Sprintf("%s-XXX-%s", cleaned[:3], cleaned[6:])
	}
	if len(cleaned) >= 6 {
		return cleaned[:3] + "-XXX-" + cleaned[len(cleaned)-2:]
	}
	return "08X-XXX-XXXX"
}

func (u *useCase) GetStalls(filter requests.StallFilter) ([]responses.StallResponse, error) {
	stalls, err := u.repo.FindAll(filter)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	res := make([]responses.StallResponse, len(stalls))
	waitlistsByStall := make(map[string][]models.Waitlist)
	if waitlists, waitlistErr := u.repo.GetActiveWaitlists(filter); waitlistErr == nil {
		for _, item := range waitlists {
			key := strings.ToUpper(strings.TrimSpace(item.StallIdent))
			if key != "" {
				waitlistsByStall[key] = append(waitlistsByStall[key], item)
			}
		}
	}
	for i, s := range stalls {
		resp := toStallResponse(&s)
		if waitlist := waitlistsByStall[strings.ToUpper(strings.TrimSpace(s.Ident))]; len(waitlist) > 0 {
			resp.WaitlistCount = len(waitlist)
			if waitlist[0].Status == "offered" && waitlist[0].ExpiresAt != nil && waitlist[0].ExpiresAt.After(time.Now()) {
				resp.OfferedToTenantID = waitlist[0].TenantID
				resp.OfferedExpiresAt = waitlist[0].ExpiresAt
				if resp.Status == "available" || resp.Status == "" {
					resp.Status = "offered"
					resp.StallStatus = "offered"
				}
			}
		}
		res[i] = resp
	}
	return res, nil
}

func (u *useCase) GetStallByID(id int) (*responses.StallResponse, error) {
	stall, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลแผงค้า")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	res := toStallResponse(stall)
	return &res, nil
}

func (u *useCase) GetWaitingQueue(ident string) (*responses.WaitingQueueResponse, error) {
	stallStatus := "available"
	stallIdent := strings.TrimSpace(ident)

	stall, err := u.repo.ByIdent(stallIdent)
	if err == nil && stall != nil {
		stallStatus = stall.StallStatus
		stallIdent = stall.Ident
	}

	waitlistItems, _ := u.repo.GetActiveWaitlistByStall(stallIdent)
	queueCount := len(waitlistItems)

	var maskedName = "-"
	var currentBookerName = ""
	var maskedPhone = "-"
	var currentBookerPhone = ""
	var activeHolderStatus = ""
	var offeredAtStr *string
	var expiresAtStr *string
	var remainingSec int64 = 0

	activeBooking, err := u.repo.GetActiveBooking(stallIdent)
	if err == nil && activeBooking != nil {
		currentBookerName = activeBooking.ApplicantName
		currentBookerPhone = activeBooking.ApplicantPhone
		maskedName = maskPDPAName(activeBooking.ApplicantName)
		maskedPhone = maskPDPAphone(activeBooking.ApplicantPhone)
		activeHolderStatus = activeBooking.BookingStatus
		if activeBooking.ExpiresAt != nil {
			tStr := activeBooking.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
			expiresAtStr = &tStr
		}
		if stallStatus == "available" {
			stallStatus = "reserved"
		}
	}

	queueList := make([]responses.WaitingQueueItem, 0, len(waitlistItems))
	now := time.Now()
	for idx, w := range waitlistItems {
		var rem int64 = 0
		if w.Status == "offered" && w.ExpiresAt != nil && w.ExpiresAt.After(now) {
			rem = int64(w.ExpiresAt.Sub(now).Seconds())
			if activeBooking == nil {
				stallStatus = "offered"
				activeHolderStatus = "offered"
				currentBookerName = w.TenantName
				currentBookerPhone = w.TenantPhone
				maskedName = maskPDPAName(w.TenantName)
				maskedPhone = maskPDPAphone(w.TenantPhone)
				if w.OfferedAt != nil {
					oStr := w.OfferedAt.Format("2006-01-02T15:04:05Z07:00")
					offeredAtStr = &oStr
				}
				eStr := w.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
				expiresAtStr = &eStr
				remainingSec = rem
			}
		}
		pos := w.Position
		if pos <= 0 {
			pos = idx + 1
		}
		queueList = append(queueList, responses.WaitingQueueItem{
			ID:               w.ID,
			Position:         pos,
			QueueNumber:      pos,
			QueuePosition:    pos,
			TenantID:         w.TenantID,
			TenantName:       w.TenantName,
			MaskedName:       maskPDPAName(w.TenantName),
			TenantPhone:      w.TenantPhone,
			MaskedPhone:      maskPDPAphone(w.TenantPhone),
			Status:           w.Status,
			OfferedAt:        w.OfferedAt,
			ExpiresAt:        w.ExpiresAt,
			RemainingSeconds: rem,
			CreatedAt:        w.CreatedAt,
		})
	}

	res := responses.WaitingQueueResponse{
		StallIdent:          stallIdent,
		QueueCount:          queueCount,
		CurrentPosition:     queueCount + 1,
		MaskedCurrentBooker: maskedName,
		CurrentBookerName:   currentBookerName,
		MaskedBookerPhone:   maskedPhone,
		CurrentBookerPhone:  currentBookerPhone,
		ActiveHolderStatus:  activeHolderStatus,
		Status:              stallStatus,
		OfferedAt:           offeredAtStr,
		ExpiresAt:           expiresAtStr,
		RemainingSeconds:    remainingSec,
		Queues:              queueList,
	}
	return &res, nil
}

func (u *useCase) CreateStall(userID int, req requests.CreateStallRequest) (*responses.StallResponse, error) {
	area := req.Width * req.Depth
	initialStatus := req.StallStatus
	if initialStatus == "" {
		initialStatus = "available"
	}
	stall := models.Stall{
		MarketID:           req.MarketID,
		MarketSlug:         req.MarketSlug,
		MarketName:         req.MarketName,
		ZoneID:             req.ZoneID,
		ZoneName:           req.ZoneName,
		ProjectName:        req.ProjectName,
		Ident:              req.Ident,
		Name:               req.Name,
		Description:        req.Description,
		Category:           req.Category,
		Width:              req.Width,
		Depth:              req.Depth,
		AreaSqm:            area,
		DailyRent:          req.DailyRent,
		MonthlyRent:        req.MonthlyRent,
		Floor:              req.Floor,
		CoordinatesGeoJSON: req.CoordinatesGeoJSON,
		Published:          req.Published,
		StallStatus:        initialStatus,
		CreatedBy:          &userID,
		UpdatedBy:          &userID,
	}

	if err := u.repo.Create(&stall); err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	res := toStallResponse(&stall)
	return &res, nil
}

func (u *useCase) UpdateStall(id int, userID int, req requests.UpdateStallRequest) (*responses.StallResponse, error) {
	stall, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลแผงค้า")
		}
		return nil, responses.NewInternalServerError()
	}
	stall.MarketID, stall.MarketSlug, stall.MarketName = req.MarketID, req.MarketSlug, req.MarketName
	stall.ZoneID, stall.ZoneName, stall.ProjectName = req.ZoneID, req.ZoneName, req.ProjectName
	stall.Ident, stall.Name, stall.Description, stall.Category = req.Ident, req.Name, req.Description, req.Category
	stall.Width, stall.Depth, stall.AreaSqm = req.Width, req.Depth, req.Width*req.Depth
	stall.DailyRent, stall.MonthlyRent = req.DailyRent, req.MonthlyRent
	stall.Floor, stall.CoordinatesGeoJSON, stall.Published = req.Floor, req.CoordinatesGeoJSON, req.Published
	if req.StallStatus != "" {
		stall.StallStatus = req.StallStatus
	}
	stall.UpdatedBy = &userID
	if err := u.repo.Update(stall, userID, "แก้ไขข้อมูล ราคา ขนาด หรือตำแหน่งแผง"); err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	res := toStallResponse(stall)
	return &res, nil
}

func (u *useCase) GetStallHistory(id int) ([]models.StallAuditLog, error) {
	items, err := u.repo.FindHistory(id)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	return items, nil
}
