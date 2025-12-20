-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Payload for the worker
    payload JSONB NOT NULL,
    incident_id UUID NOT NULL,

    -- State Machine and Retry Logic
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    retries INTEGER NOT NULL DEFAULT 0,

    -- Auditing
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notification_jobs_status ON notification_jobs (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_jobs;
-- +goose StatementEnd
