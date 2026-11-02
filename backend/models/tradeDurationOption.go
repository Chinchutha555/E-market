package models

import "encoding/json"

type TradeDurationOption struct {
	ID        string
	Duration  int
	Unit      string
	StartDate string
	EndDate   string
}

func (c *MarketTradeCondition) DurationOptions() []TradeDurationOption {
	if c.OptionsJSON != "" {
		var options []TradeDurationOption
		if json.Unmarshal([]byte(c.OptionsJSON), &options) == nil && len(options) > 0 {
			return options
		}
	}
	if c.MaxDuration > 0 && c.Unit != "" {
		return []TradeDurationOption{{ID: "legacy", Duration: c.MaxDuration, Unit: c.Unit, StartDate: c.StartDate, EndDate: c.EndDate}}
	}
	return []TradeDurationOption{}
}

func BookingDurationOptions(b *Booking) []TradeDurationOption {
	if b.TradeOptionsJSON != "" {
		var options []TradeDurationOption
		if json.Unmarshal([]byte(b.TradeOptionsJSON), &options) == nil && len(options) > 0 {
			return options
		}
	}
	if b.TradeDuration > 0 && b.TradeDurationUnit != "" {
		return []TradeDurationOption{{ID: "legacy", Duration: b.TradeDuration, Unit: b.TradeDurationUnit}}
	}
	return []TradeDurationOption{}
}
