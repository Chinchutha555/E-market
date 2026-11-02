package models

type MarketTradeCondition struct {
	OptionsJSON      string                `gorm:"column:options_json;type:text" json:"-"`
	Options          []TradeDurationOption `gorm:"-" json:"Options"`
	RequiresContract *bool                 `gorm:"column:requires_contract;default:true"`
	MarketID         int                   `gorm:"primaryKey"`
	TradeType        string                `gorm:"primaryKey;size:20"`
	Title            string                `gorm:"size:255;not null;default:''"`
	MinDuration      int
	MaxDuration      int
	Unit             string
	StartDate        string
	EndDate          string
}

func (MarketTradeCondition) TableName() string { return "market_trade_conditions" }
