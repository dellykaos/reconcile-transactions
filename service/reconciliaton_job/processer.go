package reconciliatonjob

import (
	"context"

	filestorage "github.com/delly/amartha/repository/file_storage"
	dbgen "github.com/delly/amartha/repository/postgresql"
)

// Processer is a contract to process pending reconciliation job
type Processer interface {
	Process() error
}

// ProcesserRepository is a dependency of repository that needed to process reconciliation job
type ProcesserRepository interface {
	ListPendingReconciliationJobs(ctx context.Context) ([]dbgen.ReconciliationJob, error)
	UpdateReconciliationJobStatus(ctx context.Context, arg dbgen.UpdateReconciliationJobStatusParams) (dbgen.ReconciliationJob, error)
}

// StorageGetter is a dependency of repository that needed to get file from storage
type StorageGetter interface {
	Get(filePath string) (*filestorage.File, error)
}

// ProcesserService is an implementation of Processer to process
// pending reconciliation job
type ProcesserService struct {
	repo    ProcesserRepository
	storage StorageGetter
}

var _ = Processer(&ProcesserService{})

// NewProcesserService create new processer service
func NewProcesserService(repo ProcesserRepository, storage StorageGetter) *ProcesserService {
	return &ProcesserService{repo: repo, storage: storage}
}

// Process process pending reconciliation job
func (s *ProcesserService) Process() error {
	return nil
}
