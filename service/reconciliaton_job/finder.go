package reconciliatonjob

import (
	"context"
	"errors"

	"github.com/delly/amartha/entity"
	dbgen "github.com/delly/amartha/repository/postgresql"
	"github.com/jackc/pgx/v4"
)

// Finder is a contract to find reconciliation job
type Finder interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, limit, offset int32) ([]*entity.SimpleReconciliationJob, error)
	FindByID(ctx context.Context, id int64) (*entity.ReconciliationJob, error)
}

// FinderRepository is a contract to find reconciliation job
type FinderRepository interface {
	CountReconciliationJobs(ctx context.Context) (int64, error)
	GetReconciliationJobById(ctx context.Context, id int64) (dbgen.ReconciliationJob, error)
	ListReconciliationJobs(ctx context.Context, arg dbgen.ListReconciliationJobsParams) ([]dbgen.ListReconciliationJobsRow, error)
}

// FinderService is a service to find reconciliation job
type FinderService struct {
	repo FinderRepository
}

var _ = Finder(&FinderService{})

// NewFinderService create new finder service
func NewFinderService(repo FinderRepository) *FinderService {
	return &FinderService{repo: repo}
}

// FindByID find reconciliation job by id
func (s *FinderService) FindByID(ctx context.Context, id int64) (*entity.ReconciliationJob, error) {
	rj, err := s.repo.GetReconciliationJobById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return convertToEntityReconciliationJob(rj), nil
}

// FindAll find all reconciliation job
func (s *FinderService) FindAll(ctx context.Context, limit, offset int32) ([]*entity.SimpleReconciliationJob, error) {
	params := dbgen.ListReconciliationJobsParams{
		Limit:  limit,
		Offset: offset,
	}
	rjs, err := s.repo.ListReconciliationJobs(ctx, params)
	if err != nil {
		return nil, err
	}

	res := []*entity.SimpleReconciliationJob{}
	for _, rj := range rjs {
		res = append(res, convertRowListDbToEntitySimpleReconciliationJob(rj))
	}

	return res, nil
}

// Count count all reconciliation job
func (s *FinderService) Count(ctx context.Context) (int64, error) {
	return s.repo.CountReconciliationJobs(ctx)
}
