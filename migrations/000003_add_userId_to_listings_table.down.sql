DROP INDEX IF EXISTS idx_listings_user_id;

ALTER TABLE listings
DROP COLUMN IF EXISTS user_id;