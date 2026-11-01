ALTER TABLE "application" ADD COLUMN IF NOT EXISTS market_director_signature BYTEA;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS market_director_signature_type VARCHAR(50);
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS bureau_director_signature BYTEA;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS bureau_director_signature_type VARCHAR(50);