CREATE DATABASE gateway_service_data;


CREATE TABLE IF NOT EXISTS transactions (
                                            id TEXT,
                                            reference id TEXT,
                                            account_id TEXT,
                                            amount NUMERIC(10, 2) CHECK (amount >= 0),
                                            currency TEXT,
                                            status TEXT,
                                            operation TEXT,
                                            timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                            PRIMARY KEY (id, reference_id)
);

CREATE INDEX idx_transactions_account_id ON transactions(account_id);

CREATE INDEX idx_transactions_status ON transactions(status);

