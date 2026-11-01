ALTER TABLE "application" ADD COLUMN IF NOT EXISTS signed_document BYTEA;
ALTER TABLE "application" ADD COLUMN IF NOT EXISTS signed_document_type VARCHAR(50);