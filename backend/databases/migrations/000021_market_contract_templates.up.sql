CREATE TABLE IF NOT EXISTS market_contract_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    market_id INTEGER NOT NULL,
    code VARCHAR(40) NOT NULL,
    title VARCHAR(255) NOT NULL,
    contract_kind VARCHAR(30) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT 1,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME,
    updated_at DATETIME,
    CONSTRAINT fk_market_contract_templates_market FOREIGN KEY (market_id) REFERENCES market(id),
    CONSTRAINT uq_market_contract_template_code UNIQUE (market_id, code)
);

CREATE INDEX IF NOT EXISTS idx_market_contract_templates_market_id
ON market_contract_templates(market_id);
