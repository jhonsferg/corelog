DROP INDEX IF EXISTS idx_users_team_id;

ALTER TABLE users DROP COLUMN IF EXISTS team_id;
