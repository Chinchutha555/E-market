-- 000005_application_types.up.sql

-- 1. Create application_types table
CREATE TABLE IF NOT EXISTS "application_types" (
    id SERIAL PRIMARY KEY,
    code VARCHAR(64) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_application_types_code ON "application_types"(code);

-- 2. Add missing columns to application table
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS application_type VARCHAR(50) DEFAULT 'stall_lease';
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS subject VARCHAR(255);
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS booking_id INT;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS market_id INT;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS market_name VARCHAR(255);
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS proposed_rent_type VARCHAR(50) DEFAULT 'monthly';
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS requested_start_date VARCHAR(50);
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS applicant_email VARCHAR(255);
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS applicant_address TEXT;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS details TEXT;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS officer_comment TEXT;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS market_director_comment TEXT;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS bureau_director_comment TEXT;

-- 3. Make stall_ident and shop_name nullable for general applications
ALTER TABLE "application" ALTER COLUMN stall_ident DROP NOT NULL;
ALTER TABLE "application" ALTER COLUMN shop_name DROP NOT NULL;

-- 4. Create indices for application queries
CREATE INDEX IF NOT EXISTS idx_application_market_id ON "application"(market_id);
CREATE INDEX IF NOT EXISTS idx_application_type ON "application"(application_type);
CREATE INDEX IF NOT EXISTS idx_application_state ON "application"(application_state);

-- 5. Create market_staff_assignment table if not exists
CREATE TABLE IF NOT EXISTS "market_staff_assignment" (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES "user"(id),
    market_id INT NOT NULL REFERENCES "market"(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT idx_market_staff UNIQUE (user_id, market_id)
);

CREATE INDEX IF NOT EXISTS idx_market_staff_assignment_user_id ON "market_staff_assignment"(user_id);
