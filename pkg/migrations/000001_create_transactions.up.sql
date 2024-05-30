CREATE TABLE IF NOT EXISTS transactions (
  transaction_id        SERIAL PRIMARY KEY,
  request_id            INTEGER NOT NULL,
  terminal_id           INTEGER NOT NULL,
  partner_object_id     INTEGER NOT NULL,
  amount_total          DOUBLE PRECISION NOT NULL,
  amount_original       DOUBLE PRECISION NOT NULL,
  commission_ps         DOUBLE PRECISION NOT NULL,
  commission_client     DOUBLE PRECISION NOT NULL,
  commission_provider   DOUBLE PRECISION NOT NULL,
  date_input            TIMESTAMPTZ NOT NULL,
  date_post             TIMESTAMPTZ NOT NULL,
  status                VARCHAR(255) NOT NULL,
  payment_type          VARCHAR(255) NOT NULL,
  payment_number        VARCHAR(255) NOT NULL,
  service_id            INTEGER NOT NULL,
  service               VARCHAR(255) NOT NULL,
  payee_id              INTEGER NOT NULL,
  payee_name            VARCHAR(255) NOT NULL,
  payee_bank_mfo        INTEGER NOT NULL,
  payee_bank_account    VARCHAR(255) NOT NULL,
  payment_narrative     TEXT NOT NULL
);

CREATE INDEX idx_transaction_id ON transactions (transaction_id);
CREATE INDEX idx_terminal_id ON transactions (terminal_id);
