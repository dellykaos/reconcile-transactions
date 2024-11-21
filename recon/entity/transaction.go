package entity

import (
	"time"

	"github.com/shopspring/decimal"
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
	ID     string
	Amount decimal.Decimal
	Type   TransactionType
	Time   time.Time
}
