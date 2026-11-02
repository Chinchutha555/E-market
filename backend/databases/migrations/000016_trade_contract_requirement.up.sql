ALTER TABLE market_trade_conditions ADD COLUMN requires_contract BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE booking ADD COLUMN requires_contract BOOLEAN;