package reconciliatonjob

import (
	"context"

	"github.com/delly/amartha/entity"
	dbgen "github.com/delly/amartha/repository/postgresql"
)

// Finder is a contract to find reconciliation job
type Finder interface {
	FindByID(ctx context.Context, id int64) (*entity.ReconciliationJob, error)
	FindAll(ctx context.Context, limit, offset int32) ([]*entity.ReconciliationJob, error)
}

// FinderRepository is a contract to find reconciliation job
type FinderRepository interface {
	CountReconciliationJobs(ctx context.Context) (int64, error)
	GetReconciliationJobById(ctx context.Context, id int64) (dbgen.ReconciliationJob, error)
	ListReconciliationJobs(ctx context.Context, arg dbgen.ListReconciliationJobsParams) ([]dbgen.ReconciliationJob, error)
}

// FinderService is a service to find reconciliation job
type FinderService struct {
	repo FinderRepository
}

// NewFinderService create new finder service
func NewFinderService(repo FinderRepository) *FinderService {
	return &FinderService{repo: repo}
}

// FindByID find reconciliation job by id
func (s *FinderService) FindByID(ctx context.Context, id int64) (*entity.ReconciliationJob, error) {
	rj, err := s.repo.GetReconciliationJobById(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertToEntityReconciliationJob(rj), nil
}

// FindAll find all reconciliation job
func (s *FinderService) FindAll(ctx context.Context, limit, offset int32) ([]*entity.ReconciliationJob, error) {
	params := dbgen.ListReconciliationJobsParams{
		Limit:  limit,
		Offset: offset,
	}
	rjs, err := s.repo.ListReconciliationJobs(ctx, params)
	if err != nil {
		return nil, err
	}

	res := []*entity.ReconciliationJob{}
	for _, rj := range rjs {
		res = append(res, s.convertToEntityReconciliationJob(rj))
	}

	return res, nil
}

func (s *FinderService) convertToEntityReconciliationJob(rj dbgen.ReconciliationJob) *entity.ReconciliationJob {
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
