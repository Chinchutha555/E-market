ALTER TABLE application ADD COLUMN finance_signature BLOB;
ALTER TABLE application ADD COLUMN finance_signature_type VARCHAR(50) NOT NULL DEFAULT '';
