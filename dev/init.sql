CREATE TABLE IF NOT EXISTS transactions (   id TEXT,
                                            reference_id TEXT PRIMARY KEY,
                                            account_id TEXT,
                                            amount NUMERIC(10, 2) CHECK (amount >= 0),
                                            currency TEXT,
                                            status TEXT,
                                            operation TEXT,
                                            ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP);

CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_reference_id ON transactions(reference_id);
