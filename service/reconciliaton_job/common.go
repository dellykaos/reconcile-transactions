package reconciliatonjob

import (
	"github.com/delly/amartha/entity"
	dbgen "github.com/delly/amartha/repository/postgresql"
)

func convertToEntityReconciliationJob(rj dbgen.ReconciliationJob) *entity.ReconciliationJob {
	res := &entity.ReconciliationJob{
		ID:                       rj.ID,
		Status:                   entity.ReconciliationJobStatus(rj.Status),
		SystemTransactionCsvPath: rj.SystemTransactionCsvPath,
		DiscrepancyThreshold:     float32(rj.DiscrepancyThreshold),
		StartDate:                rj.StartDate,
		EndDate:                  rj.EndDate,
		CreatedAt:                rj.CreatedAt,
		UpdatedAt:                rj.UpdatedAt,
	}
	rj.BankTransactionCsvPaths.AssignTo(&res.BankTransactionCsvPaths)
	rj.Result.AssignTo(&res.Result)

	return res
}
