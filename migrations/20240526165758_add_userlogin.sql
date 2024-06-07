-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD userlogin varchar NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN userlogin;
-- +goose StatementEnd
