// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: reconciliation_jobs.sql

package dbgen

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
)

const createReconciliationJob = `-- name: CreateReconciliationJob :one
INSERT INTO reconciliation_jobs (status, system_transaction_csv_path, bank_transaction_csv_paths, discrepancy_threshold, start_date, end_date) VALUES ('PENDING', $1, $2, $3, $4, $5) RETURNING id, status, system_transaction_csv_path, bank_transaction_csv_paths, discrepancy_threshold, start_date, end_date, result, created_at, updated_at
`

type CreateReconciliationJobParams struct {
	SystemTransactionCsvPath string       `db:"system_transaction_csv_path"`
	BankTransactionCsvPaths  pgtype.JSONB `db:"bank_transaction_csv_paths"`
	DiscrepancyThreshold     float64      `db:"discrepancy_threshold"`
	StartDate                time.Time    `db:"start_date"`
	EndDate                  time.Time    `db:"end_date"`
}

func (q *Queries) CreateReconciliationJob(ctx context.Context, arg CreateReconciliationJobParams) (ReconciliationJob, error) {
	row := q.db.QueryRow(ctx, createReconciliationJob,
		arg.SystemTransactionCsvPath,
		arg.BankTransactionCsvPaths,
		arg.DiscrepancyThreshold,
		arg.StartDate,
		arg.EndDate,
	)
	var i ReconciliationJob
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.SystemTransactionCsvPath,
		&i.BankTransactionCsvPaths,
		&i.DiscrepancyThreshold,
		&i.StartDate,
		&i.EndDate,
		&i.Result,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const finishReconciliationJob = `-- name: FinishReconciliationJob :one
UPDATE reconciliation_jobs SET status = 'SUCCESS', result = $2 WHERE id = $1 RETURNING id, status, system_transaction_csv_path, bank_transaction_csv_paths, discrepancy_threshold, start_date, end_date, result, created_at, updated_at
`

type FinishReconciliationJobParams struct {
	ID     int64        `db:"id"`
	Result pgtype.JSONB `db:"result"`
}

func (q *Queries) FinishReconciliationJob(ctx context.Context, arg FinishReconciliationJobParams) (ReconciliationJob, error) {
	row := q.db.QueryRow(ctx, finishReconciliationJob, arg.ID, arg.Result)
	var i ReconciliationJob
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.SystemTransactionCsvPath,
		&i.BankTransactionCsvPaths,
		&i.DiscrepancyThreshold,
		&i.StartDate,
		&i.EndDate,
		&i.Result,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getReconciliationJobById = `-- name: GetReconciliationJobById :one
SELECT id, status, system_transaction_csv_path, bank_transaction_csv_paths, discrepancy_threshold, start_date, end_date, result, created_at, updated_at FROM reconciliation_jobs WHERE id = $1
`

func (q *Queries) GetReconciliationJobById(ctx context.Context, id int64) (ReconciliationJob, error) {
	row := q.db.QueryRow(ctx, getReconciliationJobById, id)
	var i ReconciliationJob
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.SystemTransactionCsvPath,
		&i.BankTransactionCsvPaths,
		&i.DiscrepancyThreshold,
		&i.StartDate,
		&i.EndDate,
		&i.Result,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listReconciliationJobs = `-- name: ListReconciliationJobs :many
SELECT id, status, system_transaction_csv_path, bank_transaction_csv_paths, discrepancy_threshold, start_date, end_date, result, created_at, updated_at FROM reconciliation_jobs WHERE status = 'PENDING' ORDER BY created_at ASC
`

func (q *Queries) ListReconciliationJobs(ctx context.Context) ([]ReconciliationJob, error) {
	rows, err := q.db.Query(ctx, listReconciliationJobs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReconciliationJob
	for rows.Next() {
		var i ReconciliationJob
		if err := rows.Scan(
			&i.ID,
			&i.Status,
			&i.SystemTransactionCsvPath,
			&i.BankTransactionCsvPaths,
			&i.DiscrepancyThreshold,
			&i.StartDate,
			&i.EndDate,
			&i.Result,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateReconciliationJobStatus = `-- name: UpdateReconciliationJobStatus :one
UPDATE reconciliation_jobs SET status = $2 WHERE id = $1 RETURNING id, status, system_transaction_csv_path, bank_transaction_csv_paths, discrepancy_threshold, start_date, end_date, result, created_at, updated_at
`

type UpdateReconciliationJobStatusParams struct {
	ID     int64  `db:"id"`
	Status string `db:"status"`
}

func (q *Queries) UpdateReconciliationJobStatus(ctx context.Context, arg UpdateReconciliationJobStatusParams) (ReconciliationJob, error) {
	row := q.db.QueryRow(ctx, updateReconciliationJobStatus, arg.ID, arg.Status)
	var i ReconciliationJob
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.SystemTransactionCsvPath,
		&i.BankTransactionCsvPaths,
		&i.DiscrepancyThreshold,
		&i.StartDate,
		&i.EndDate,
		&i.Result,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
