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

func convertToDBReconciliationJob(job *entity.ReconciliationJob) dbgen.ReconciliationJob {
	res := dbgen.ReconciliationJob{
		SystemTransactionCsvPath: job.SystemTransactionCsvPath,
		DiscrepancyThreshold:     float64(job.DiscrepancyThreshold),
		StartDate:                job.StartDate,
		EndDate:                  job.EndDate,
		CreatedAt:                job.CreatedAt,
		UpdatedAt:                job.UpdatedAt,
	}
	res.BankTransactionCsvPaths.Set(job.BankTransactionCsvPaths)

	return res
}
