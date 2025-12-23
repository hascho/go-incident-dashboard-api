-- +goose Up
-- +goose StatementBegin
ALTER TABLE incidents ADD COLUMN notification_status VARCHAR(20) DEFAULT 'pending';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE incidents DROP COLUMN notification_status;
-- +goose StatementEnd
