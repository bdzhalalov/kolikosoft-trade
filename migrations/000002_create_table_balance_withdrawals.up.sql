CREATE TABLE balance_withdrawals (
         id             BIGSERIAL PRIMARY KEY,
         user_id        BIGINT NOT NULL REFERENCES users(id),
         amount         BIGINT NOT NULL CHECK (amount > 0),
         balance_before BIGINT NOT NULL,
         balance_after  BIGINT NOT NULL,
         created_at     TIMESTAMP NOT NULL DEFAULT now(),
         request_id     VARCHAR(256) NOT NULL
);

CREATE UNIQUE INDEX ux_balance_withdrawals_user_request
    ON balance_withdrawals(user_id, request_id);