package reconciliatonjob

import (
	"bytes"
	"context"
	"time"

	"github.com/delly/amartha/entity"
	dbgen "github.com/delly/amartha/repository/postgresql"
)

// Creator is a contract to create reconciliation job
type Creator interface {
	Create(ctx context.Context, params *CreateParams) (*entity.ReconciliationJob, error)
}

// CreatorRepository is a contract to create reconciliation job
type CreatorRepository interface {
	CreateReconciliationJob(ctx context.Context, job dbgen.ReconciliationJob) (dbgen.ReconciliationJob, error)
}

// CreatorService is a service to create reconciliation job
type CreatorService struct {
	repo CreatorRepository
}

// File is a struct to hold metadata of csv file
type File struct {
	Name string
	Buf  *bytes.Buffer
	Path string
}

// BankTransactionFile is a struct to hold metadata of bank transaction csv files
type BankTransactionFile struct {
	BankName string
	File     *File
}

// CreateParams is a parameter to create reconciliation job
type CreateParams struct {
	SystemTransactionCsv *File
	BankTransactionCsvs  []*BankTransactionFile
	StartDate            time.Time
	EndDate              time.Time
	DiscrepancyThreshold float32
}

var _ = Creator(&CreatorService{})

// NewCreatorService create new creator service
func NewCreatorService(repo CreatorRepository) *CreatorService {
	return &CreatorService{repo: repo}
}

// Create create reconciliation job
func (s *CreatorService) Create(ctx context.Context, params *CreateParams) (*entity.ReconciliationJob, error) {
	// TODO: store csv files to storage

	rj, err := s.repo.CreateReconciliationJob(ctx, params.convertParamsToDBReconciliationJob())
	if err != nil {
		return nil, err
	}

	return convertToEntityReconciliationJob(rj), nil
}

func (p *CreateParams) convertParamsToDBReconciliationJob() dbgen.ReconciliationJob {
	res := dbgen.ReconciliationJob{
		SystemTransactionCsvPath: p.SystemTransactionCsv.Path,
		DiscrepancyThreshold:     float64(p.DiscrepancyThreshold),
		StartDate:                p.StartDate,
		EndDate:                  p.EndDate,
	}
	res.BankTransactionCsvPaths.Set(p.convertBankTransactionFilesToEntity())

	return res
}

func (p *CreateParams) convertBankTransactionFilesToEntity() []entity.BankTransactionCsv {
	res := make([]entity.BankTransactionCsv, len(p.BankTransactionCsvs))
	for i, v := range p.BankTransactionCsvs {
		res[i] = entity.BankTransactionCsv{
			BankName: v.BankName,
			FilePath: v.File.Path,
		}
	}

	return res
}
