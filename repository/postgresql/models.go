// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package dbgen

import (
	"time"

	"github.com/jackc/pgtype"
)

type ReconciliationJob struct {
	ID                       int64        `db:"id"`
	Status                   string       `db:"status"`
	SystemTransactionCsvPath string       `db:"system_transaction_csv_path"`
	BankTransactionCsvPaths  pgtype.JSONB `db:"bank_transaction_csv_paths"`
	DiscrepancyThreshold     float64      `db:"discrepancy_threshold"`
	StartDate                time.Time    `db:"start_date"`
	EndDate                  time.Time    `db:"end_date"`
	Result                   pgtype.JSONB `db:"result"`
	CreatedAt                time.Time    `db:"created_at"`
	UpdatedAt                time.Time    `db:"updated_at"`
}
