CREATE TABLE IF NOT EXISTS tickets (
    id            UUID PRIMARY KEY,
    title         VARCHAR(200) NOT NULL,
    description   TEXT NOT NULL,
    status        VARCHAR(20) NOT NULL CHECK (status IN ('open', 'in_progress', 'resolved', 'closed')),
    priority      VARCHAR(20) NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    requester_id  UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    assignee_id   UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets (status);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets (priority);
CREATE INDEX IF NOT EXISTS idx_tickets_requester_id ON tickets (requester_id);
CREATE INDEX IF NOT EXISTS idx_tickets_assignee_id ON tickets (assignee_id);
