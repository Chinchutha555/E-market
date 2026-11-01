DROP INDEX IF EXISTS idx_notification_source_recipient;
ALTER TABLE "notification" DROP COLUMN IF EXISTS source_notification_id;