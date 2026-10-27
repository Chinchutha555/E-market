-- 000005_application_types.down.sql

DROP TABLE IF EXISTS "market_staff_assignment";
DROP TABLE IF EXISTS "application_types";
ALTER TABLE "application" DROP COLUMN IF EXISTS application_type;
ALTER TABLE "application" DROP COLUMN IF EXISTS subject;
ALTER TABLE "application" DROP COLUMN IF EXISTS booking_id;
ALTER TABLE "application" DROP COLUMN IF EXISTS market_id;
ALTER TABLE "application" DROP COLUMN IF EXISTS market_name;
ALTER TABLE "application" DROP COLUMN IF EXISTS proposed_rent_type;
ALTER TABLE "application" DROP COLUMN IF EXISTS requested_start_date;
ALTER TABLE "application" DROP COLUMN IF EXISTS applicant_email;
ALTER TABLE "application" DROP COLUMN IF EXISTS applicant_address;
ALTER TABLE "application" DROP COLUMN IF EXISTS details;
ALTER TABLE "application" DROP COLUMN IF EXISTS officer_comment;
ALTER TABLE "application" DROP COLUMN IF EXISTS market_director_comment;
ALTER TABLE "application" DROP COLUMN IF EXISTS bureau_director_comment;
