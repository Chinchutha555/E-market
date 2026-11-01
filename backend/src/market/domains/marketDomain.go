package domains

import (
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
)

type UseCase interface {
	GetMarkets(search string) ([]models.Market, error)
	GetMarketsForUser(search string, userID int, role string) ([]models.Market, error)
	GetMarketByIDOrSlug(idOrSlug string) (*models.Market, error)
	GetMarketZones(marketID int) ([]models.Zone, error)
	GetGISConfig(idOrSlug string, includeDraft bool) (*responses.MarketGISResponse, error)
	SaveGISConfig(idOrSlug string, userID int, req requests.UpdateMarketGISRequest) (*responses.MarketGISResponse, error)
	GetUtilitySettings(idOrSlug string) (*responses.MarketUtilitySettingsResponse, error)
	SaveUtilitySettings(idOrSlug string, userID int, req requests.UpdateMarketUtilitySettingsRequest) (*responses.MarketUtilitySettingsResponse, error)
}

type Repository interface {
	FindAll(search string) ([]models.Market, error)
	FindByStaffUserID(search string, userID int) ([]models.Market, error)
	FindByIDOrSlug(idOrSlug string) (*models.Market, error)
	FindZonesByMarketID(marketID int) ([]models.Zone, error)
	FindGISConfig(marketID int) (*models.MarketGISConfig, error)
	SaveGISConfig(config *models.MarketGISConfig) error
	FindUtilitySettings(marketID int) (*models.MarketUtilitySetting, error)
	SaveUtilitySettings(setting *models.MarketUtilitySetting) error
}
