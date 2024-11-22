package reconciliatonjob_test

import (
	"context"
	"testing"
	"time"

	"github.com/delly/amartha/entity"
	dbgen "github.com/delly/amartha/repository/postgresql"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	mock_reconciliatonjob "github.com/delly/amartha/test/mock/service/reconciliaton_job"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ReconciliationJobCreatorTestSuite struct {
	suite.Suite
	repo *mock_reconciliatonjob.MockCreatorRepository
	svc  reconciliatonjob.Creator
}

func (s *ReconciliationJobCreatorTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.repo = mock_reconciliatonjob.NewMockCreatorRepository(ctrl)
	s.svc = reconciliatonjob.NewCreatorService(s.repo)
}

func TestReconciliationJobCreatorTestSuite(t *testing.T) {
	suite.Run(t, new(ReconciliationJobCreatorTestSuite))
}

func (s *ReconciliationJobCreatorTestSuite) TestCreate() {
	ctx := context.Background()

	job := &entity.ReconciliationJob{
		SystemTransactionCsvPath: "/path/to/system/transaction.csv",
		DiscrepancyThreshold:     0.1,
		StartDate:                time.Now(),
		EndDate:                  time.Now(),
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
		BankTransactionCsvPaths: []entity.BankTransactionCsv{
			{
				BankName: "BCA",
				FilePath: "/path/to/bca/transaction.csv",
			},
		},
	}
	dbParams := dbgen.ReconciliationJob{
		SystemTransactionCsvPath: job.SystemTransactionCsvPath,
		DiscrepancyThreshold:     float64(job.DiscrepancyThreshold),
		StartDate:                job.StartDate,
		EndDate:                  job.EndDate,
		CreatedAt:                job.CreatedAt,
		UpdatedAt:                job.UpdatedAt,
	}
	dbParams.BankTransactionCsvPaths.Set(job.BankTransactionCsvPaths)
	dbResult := dbParams
	dbResult.Status = "PENDING"
	jrResult := &entity.ReconciliationJob{
		ID:                       dbResult.ID,
		Status:                   entity.ReconciliationJobStatus(dbResult.Status),
		SystemTransactionCsvPath: dbResult.SystemTransactionCsvPath,
		BankTransactionCsvPaths:  job.BankTransactionCsvPaths,
		DiscrepancyThreshold:     float32(dbResult.DiscrepancyThreshold),
		StartDate:                dbResult.StartDate,
		EndDate:                  dbResult.EndDate,
		CreatedAt:                dbResult.CreatedAt,
		UpdatedAt:                dbResult.UpdatedAt,
	}

	s.Run("success", func() {
		s.repo.EXPECT().CreateReconciliationJob(ctx, dbParams).Return(dbResult, nil)

		res, err := s.svc.Create(ctx, job)

		s.Nil(err)
		s.Equal(jrResult, res)
	})
}
