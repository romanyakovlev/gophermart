-- +goose Up
-- +goose StatementBegin

-- Create the users table with userid as UUID
CREATE TABLE users (
                       userid UUID PRIMARY KEY,
                       pass_hash VARCHAR(255) NOT NULL,
                       current_balance NUMERIC NOT NULL,
                       withdrawn_balance NUMERIC NOT NULL
);

-- Create the orders table with user_id as UUID
CREATE TABLE orders (
                        number VARCHAR(255),
                        status VARCHAR(255),
                        accrual NUMERIC,
                        uploaded_at TIMESTAMP,
                        user_id UUID,
                        FOREIGN KEY (user_id) REFERENCES users(userid)
);

-- Create the withdrawals table with user_id as UUID
CREATE TABLE withdrawals (
                             number VARCHAR(255),
                             sum NUMERIC,
                             processed_at TIMESTAMP,
                             user_id UUID,
                             FOREIGN KEY (user_id) REFERENCES users(userid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop table orders;
drop table withdrawals;
-- +goose StatementEnd
