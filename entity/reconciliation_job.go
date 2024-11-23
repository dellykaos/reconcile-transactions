package entity

import (
	"time"
)

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
	BankName string `json:"bank_name"`
	FilePath string `json:"file_path"`
}

// ReconciliationResult hold reconciliation result data
type ReconciliationResult struct {
	TotalTransactionProcessed int                      `json:"total_transaction_processed"`
	TotalTransactionMatched   int                      `json:"total_transaction_matched"`
	TotalTransactionUnmatched int                      `json:"total_transaction_unmatched"`
	TotalDiscrepancyAmount    float64                  `json:"total_discrepancy_amount"`
	MissingTransactions       []Transaction            `json:"missing_transactions"`
	MissingBankTransactions   map[string][]Transaction `json:"missing_bank_transactions"`
}

// ReconciliationJob hold reconciliation job data
type ReconciliationJob struct {
	ID                       int64                   `json:"id"`
	Status                   ReconciliationJobStatus `json:"status"`
	SystemTransactionCsvPath string                  `json:"system_transaction_csv_path"`
	BankTransactionCsvPaths  []BankTransactionCsv    `json:"bank_transaction_csv_paths"`
	DiscrepancyThreshold     float32                 `json:"discrepancy_threshold"`
	ErrorInformation         string                  `json:"error_information"`
	Result                   *ReconciliationResult   `json:"result"`
	StartDate                time.Time               `json:"start_date"`
	EndDate                  time.Time               `json:"end_date"`
	CreatedAt                time.Time               `json:"created_at"`
	UpdatedAt                time.Time               `json:"updated_at"`
}
