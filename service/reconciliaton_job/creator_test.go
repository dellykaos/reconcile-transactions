package reconciliatonjob_test

import (
	"context"
	"testing"
	"time"

	"github.com/delly/amartha/entity"
	dbgen "github.com/delly/amartha/repository/postgresql"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	mock_reconciliatonjob "github.com/delly/amartha/test/mock/service/reconciliaton_job"
	"github.com/stretchr/testify/assert"
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

	params := &reconciliatonjob.CreateParams{
		SystemTransactionCsv: &reconciliatonjob.File{
			Path: "/path/to/system/transaction.csv",
		},
		DiscrepancyThreshold: 0.1,
		StartDate:            time.Now(),
		EndDate:              time.Now(),
		BankTransactionCsvs: []*reconciliatonjob.BankTransactionFile{
			{
				BankName: "BCA",
				File: &reconciliatonjob.File{
					Path: "/path/to/bca/transaction.csv",
				},
			},
		},
	}
	bankTrxCsvPaths := []entity.BankTransactionCsv{
		{
			BankName: params.BankTransactionCsvs[0].BankName,
			FilePath: params.BankTransactionCsvs[0].File.Path,
		},
	}
	dbParams := dbgen.ReconciliationJob{
		SystemTransactionCsvPath: params.SystemTransactionCsv.Path,
		DiscrepancyThreshold:     float64(params.DiscrepancyThreshold),
		StartDate:                params.StartDate,
		EndDate:                  params.EndDate,
	}
	dbParams.BankTransactionCsvPaths.Set(bankTrxCsvPaths)
	dbResult := dbParams
	dbResult.Status = "PENDING"
	jrResult := &entity.ReconciliationJob{
		ID:                       dbResult.ID,
		Status:                   entity.ReconciliationJobStatus(dbResult.Status),
		SystemTransactionCsvPath: dbResult.SystemTransactionCsvPath,
		BankTransactionCsvPaths:  bankTrxCsvPaths,
		DiscrepancyThreshold:     float32(dbResult.DiscrepancyThreshold),
		StartDate:                dbResult.StartDate,
		EndDate:                  dbResult.EndDate,
		CreatedAt:                dbResult.CreatedAt,
		UpdatedAt:                dbResult.UpdatedAt,
	}

	s.Run("success", func() {
		s.repo.EXPECT().CreateReconciliationJob(ctx, dbParams).Return(dbResult, nil)

		res, err := s.svc.Create(ctx, params)

		s.Nil(err)
		s.Equal(jrResult, res)
	})

	s.Run("error", func() {
		s.repo.EXPECT().CreateReconciliationJob(ctx, dbParams).Return(dbgen.ReconciliationJob{}, assert.AnError)

		res, err := s.svc.Create(ctx, params)

		s.Nil(res)
		s.Equal(assert.AnError, err)
	})
}
