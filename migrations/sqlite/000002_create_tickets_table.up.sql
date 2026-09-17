CREATE TABLE IF NOT EXISTS tickets (
    id            TEXT PRIMARY KEY,
    title         TEXT NOT NULL,
    description   TEXT NOT NULL,
    status        TEXT NOT NULL CHECK (status IN ('open', 'in_progress', 'resolved', 'closed')),
    priority      TEXT NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    requester_id  TEXT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    assignee_id   TEXT REFERENCES users (id) ON DELETE SET NULL,
    created_at    TEXT NOT NULL,
    updated_at    TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets (status);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets (priority);
CREATE INDEX IF NOT EXISTS idx_tickets_requester_id ON tickets (requester_id);
CREATE INDEX IF NOT EXISTS idx_tickets_assignee_id ON tickets (assignee_id);
