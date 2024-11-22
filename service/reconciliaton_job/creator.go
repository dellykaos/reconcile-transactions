package reconciliatonjob

import (
	"context"

	"github.com/delly/amartha/entity"
	dbgen "github.com/delly/amartha/repository/postgresql"
)

// Creator is a contract to create reconciliation job
type Creator interface {
	Create(ctx context.Context, job *entity.ReconciliationJob) (*entity.ReconciliationJob, error)
}

// CreatorRepository is a contract to create reconciliation job
type CreatorRepository interface {
	CreateReconciliationJob(ctx context.Context, job dbgen.ReconciliationJob) (dbgen.ReconciliationJob, error)
}

// CreatorService is a service to create reconciliation job
type CreatorService struct {
	repo CreatorRepository
}

var _ = Creator(&CreatorService{})

// NewCreatorService create new creator service
func NewCreatorService(repo CreatorRepository) *CreatorService {
	return &CreatorService{repo: repo}
}

// Create create reconciliation job
func (s *CreatorService) Create(ctx context.Context, job *entity.ReconciliationJob) (*entity.ReconciliationJob, error) {
	// TODO: store csv files to storage

	rj, err := s.repo.CreateReconciliationJob(ctx, convertToDBReconciliationJob(job))
	if err != nil {
		return nil, err
	}

	return convertToEntityReconciliationJob(rj), nil
}
