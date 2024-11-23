package reconciliatonjob

import (
	"fmt"

	"github.com/delly/amartha/entity"
)

var (
	errEmptyBuffer = func(filename string) error {
		return fmt.Errorf("file buffer of file %s is empty", filename)
	}
	errInvalidTrxType = func(trxType entity.TransactionType, trxID string) error {
		return fmt.Errorf("invalid transaction type: %s, trx id: %s", trxType, trxID)
	}
)
