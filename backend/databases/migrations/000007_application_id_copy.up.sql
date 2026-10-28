ALTER TABLE "application" ADD COLUMN IF NOT EXISTS applicant_age INTEGER;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS id_card_copy BYTEA;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS id_card_copy_type VARCHAR(50);
