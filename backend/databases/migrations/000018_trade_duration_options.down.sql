ALTER TABLE application DROP COLUMN trade_end_date;
ALTER TABLE application DROP COLUMN trade_start_date;
ALTER TABLE application DROP COLUMN trade_option_id;
ALTER TABLE booking DROP COLUMN trade_options_json;
ALTER TABLE market_trade_conditions DROP COLUMN options_json;