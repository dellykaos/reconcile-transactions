package entity

import "time"

// ReconcileJobStatus is a custom type for reconcile job status
type ReconcileJobStatus string

const (
	// ReconcileJobStatusPending is a constant for pending reconcile job status
	ReconcileJobStatusPending ReconcileJobStatus = "PENDING"
	// ReconcileJobStatusProcessing is a constant for processing reconcile job status
	ReconcileJobStatusProcessing ReconcileJobStatus = "PROCESSING"
	// ReconcileJobStatusSuccess is a constant for success reconcile job status
	ReconcileJobStatusSuccess ReconcileJobStatus = "SUCCESS"
	// ReconcileJobStatusFailed is a constant for failed reconcile job status
	ReconcileJobStatusFailed ReconcileJobStatus = "FAILED"
)

// BankTransactionCsv hold bank transaction csv data
type BankTransactionCsv struct {
	BankName string
	FilePath string
}

// ReconcileJob hold reconcile job data
type ReconcileJob struct {
	ID                       int
	Status                   ReconcileJobStatus
	SystemTransactionCsvPath string
	BankTransactionCsvPaths  []BankTransactionCsv
	DiscrepancyThreshold     float32
	StartDate                time.Time
	EndDate                  time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
