ALTER TABLE users ADD COLUMN team_id TEXT REFERENCES teams (id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_users_team_id ON users (team_id);
