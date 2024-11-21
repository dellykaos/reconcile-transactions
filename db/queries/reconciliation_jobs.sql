-- name: ListReconciliationJobs :many
SELECT * FROM reconciliation_jobs WHERE status = 'PENDING' ORDER BY created_at ASC;

-- name: GetReconciliationJobById :one
SELECT * FROM reconciliation_jobs WHERE id = $1;

-- name: CreateReconciliationJob :one
INSERT INTO reconciliation_jobs (status, system_transaction_csv_path, bank_transaction_csv_paths, discrepancy_threshold, start_date, end_date) VALUES ('PENDING', $1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateReconciliationJobStatus :one
UPDATE reconciliation_jobs SET status = $2 WHERE id = $1 RETURNING *;

-- name: FinishReconciliationJob :one
UPDATE reconciliation_jobs SET status = 'SUCCESS', result = $2 WHERE id = $1 RETURNING *;