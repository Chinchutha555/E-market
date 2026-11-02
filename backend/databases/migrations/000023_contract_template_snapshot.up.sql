ALTER TABLE contract ADD COLUMN contract_template_id INTEGER;
ALTER TABLE contract ADD COLUMN contract_template_code VARCHAR(40) NOT NULL DEFAULT '';
ALTER TABLE contract ADD COLUMN contract_template_title VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE contract ADD COLUMN contract_template_kind VARCHAR(30) NOT NULL DEFAULT '';
ALTER TABLE contract ADD COLUMN contract_template_description TEXT NOT NULL DEFAULT '';
