ALTER TABLE market_trade_conditions ADD COLUMN options_json TEXT NOT NULL DEFAULT '';
ALTER TABLE booking ADD COLUMN trade_options_json TEXT NOT NULL DEFAULT '';
ALTER TABLE application ADD COLUMN trade_option_id VARCHAR(20) NOT NULL DEFAULT '';
ALTER TABLE application ADD COLUMN trade_start_date VARCHAR(10) NOT NULL DEFAULT '';
ALTER TABLE application ADD COLUMN trade_end_date VARCHAR(10) NOT NULL DEFAULT '';