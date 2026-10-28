CREATE TABLE IF NOT EXISTS tenant_debt (
  id SERIAL PRIMARY KEY,
  tenant_id INTEGER NOT NULL REFERENCES "user"(id),
  market_id INTEGER NOT NULL REFERENCES market(id),
  description VARCHAR(500) NOT NULL,
  amount_satang BIGINT NOT NULL CHECK (amount_satang > 0),
  status VARCHAR(20) NOT NULL DEFAULT 'outstanding',
  created_by INTEGER NOT NULL REFERENCES "user"(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_tenant_debt_tenant_id ON tenant_debt(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tenant_debt_market_id ON tenant_debt(market_id);
CREATE INDEX IF NOT EXISTS idx_tenant_debt_status ON tenant_debt(status);