CREATE TABLE IF NOT EXISTS market_trade_conditions (
 market_id INTEGER NOT NULL REFERENCES market(id),
 trade_type VARCHAR(20) NOT NULL,
 min_duration INTEGER NOT NULL,
 max_duration INTEGER NOT NULL,
 unit VARCHAR(10) NOT NULL,
 PRIMARY KEY (market_id, trade_type)
);