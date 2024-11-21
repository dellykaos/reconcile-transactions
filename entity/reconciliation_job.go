package entity

import "time"

// ReconciliationJobStatus is a custom type for reconciliation job status
type ReconciliationJobStatus string

const (
	// ReconciliationJobStatusPending is a pending status of reconciliation job
	ReconciliationJobStatusPending ReconciliationJobStatus = "PENDING"
	// ReconciliationJobStatusProcessing is a processing status of reconciliation job
	ReconciliationJobStatusProcessing ReconciliationJobStatus = "PROCESSING"
	// ReconciliationJobStatusSuccess is a success status of reconciliation job
	ReconciliationJobStatusSuccess ReconciliationJobStatus = "SUCCESS"
	// ReconciliationJobStatusFailed is a failed status of reconciliation job
	ReconciliationJobStatusFailed ReconciliationJobStatus = "FAILED"
)

// BankTransactionCsv hold bank transaction csv data
type BankTransactionCsv struct {
	BankName string
	FilePath string
}

// ReconciliationJob hold reconciliation job data
type ReconciliationJob struct {
	ID                       int
	Status                   ReconciliationJobStatus
	SystemTransactionCsvPath string
	BankTransactionCsvPaths  []BankTransactionCsv
	DiscrepancyThreshold     float32
	StartDate                time.Time
	EndDate                  time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
