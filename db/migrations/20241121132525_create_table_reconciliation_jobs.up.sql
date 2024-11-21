BEGIN;

CREATE TABLE IF NOT EXISTS reconciliation_jobs (
    id BIGSERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    system_transaction_csv_path VARCHAR NOT NULL,
    bank_transaction_csv_paths JSONB NOT NULL,
    discrepancy_threshold FLOAT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    result JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_reconciliation_jobs_status_created_at ON reconciliation_jobs(status, created_at);

END;
