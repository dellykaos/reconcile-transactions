package entity

import (
	"time"
)

// TransactionType is a custom type for transaction type
type TransactionType string

const (
	// TxTypeDebit is a constant for debit transaction type
	TxTypeDebit TransactionType = "DEBIT"
	// TxTypeCredit is a constant for credit transaction type
	TxTypeCredit TransactionType = "CREDIT"
)

// Transaction hold transaction data
type Transaction struct {
	ID     string          `json:"id"`
	Amount float64         `json:"amount"`
	Type   TransactionType `json:"type"`
	Time   time.Time       `json:"time"`
}
