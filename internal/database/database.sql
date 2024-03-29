CREATE TABLE IF NOT EXISTS transactions_history (
  transaction_id VARCHAR PRIMARY KEY,
  status VARCHAR(20) NOT NULL,
  failure_reason VARCHAR(50),
  payment_provider VARCHAR(20) NOT NULL,
  description VARCHAR(100) NOT NULL,
  amount NUMERIC NOT NULL,
  currency CHAR(3) NOT NULL,
  type VARCHAR(10) NOT NULL,
  additional_fields JSONB
);