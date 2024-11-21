package reconciliatonjob

import (
	"context"

	"github.com/delly/amartha/entity"
	dbgen "github.com/delly/amartha/repository/postgresql"
)

// Finder is a contract to find reconciliation job
type Finder interface {
	FindByID(ctx context.Context, id int64) (*entity.ReconciliationJob, error)
	FindAll(ctx context.Context, limit, offset int) ([]*entity.ReconciliationJob, error)
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
