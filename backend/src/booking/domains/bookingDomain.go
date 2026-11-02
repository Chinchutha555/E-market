package domains

import (
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
)

type UseCase interface {
	GetBookings(filter requests.BookingFilter, actorID int, actorRole string) ([]responses.BookingResponse, error)
	GetBookingByID(id int, actorID int, actorRole string) (*responses.BookingResponse, error)
	CreateBooking(userID int, req requests.CreateBookingRequest) (*responses.BookingResponse, error)
	CreateWaitlist(userID int, req requests.CreateWaitlistRequest) (*responses.WaitlistResponse, error)
	GetMyWaitlist(userID int) ([]responses.WaitlistResponse, error)
	GetWaitlistByStall(stallIdent string, actorID int, actorRole string) ([]responses.WaitlistResponse, error)
	WithdrawWaitlist(id int, userID int, actorRole string) error
	AdvanceWaitlistByStall(stallIdent string, actorID int, actorRole string) ([]responses.WaitlistResponse, error)
	CancelBookingStall(id int, stallIdent string, userID int, actorRole string, reason string) (*responses.BookingResponse, error)
	UpdateStatus(id int, userID int, req requests.UpdateBookingStatusRequest) (*responses.BookingResponse, error)
	UpdateVisit(id int, actorID int, actorRole string, visitDate string, visitTime string) (*responses.BookingResponse, error)
	GetPolicySettings(marketID int, marketName string) (*responses.BookingPolicySettingsResponse, error)
	UpdatePolicySettings(actorRole string, req requests.UpdateBookingPolicySettingsRequest) (*responses.BookingPolicySettingsResponse, error)
	ClearAll(actorRole string) error
}

type Repository interface {
	FindTradeCondition(marketID int, code string) (*models.MarketTradeCondition, error)
	ProcessWaitlistLifecycle() error
	GetPolicySettings(marketID int, marketName string) (leaseDays int, offerHours int, resolvedMarketID int, resolvedMarketName string)
	UpdatePolicySettings(marketID int, marketName string, leaseDays int, offerHours int) (*responses.BookingPolicySettingsResponse, error)
	FindAll(filter requests.BookingFilter, ownerID *int) ([]models.Booking, error)
	FindByID(id int) (*models.Booking, error)
	FindByGroupID(groupID string) ([]models.Booking, error)
	FindActiveDuplicate(stallIdent string, userID int, applicantName string, applicantPhone string) (*models.Booking, error)
	HasActiveHolder(stallIdent string) (bool, error)
	Create(booking *models.Booking) error
	Update(booking *models.Booking) error
	CreateNotification(notif *models.Notification) error
	CreateNotifications(notifs []models.Notification) error
	CountActiveQueue(stallIdent string) (int64, error)
	CountWaitlistByStall(stallIdent string) (int64, error)
	UpdateStallStatus(stallIdent string, status string) error
	ClearAll() error
	CreateWaitlist(item *models.Waitlist) error
	FindWaitlistByTenant(userID int) ([]models.Waitlist, error)
	FindWaitlistByStall(stallIdent string) ([]models.Waitlist, error)
	FindWaitlistByID(id int) (*models.Waitlist, error)
	WithdrawWaitlist(item *models.Waitlist) error
	AdvanceWaitlistByStall(stallIdent string) error
	FindActiveWaitlistOffer(stallIdent string) (*models.Waitlist, error)
	AcceptWaitlistOffer(stallIdent string, tenantID int) error
	LookupUserByID(userID int) (*models.User, error)
	LookupStallRate(stallIdent string, marketID *int) (monthlyRent float64, depositFee float64, projectName string, marketName string)
	SyncBookingStalls(parentBookingID int, parentBookingNo string, groupID string, bookings []models.Booking)
}
