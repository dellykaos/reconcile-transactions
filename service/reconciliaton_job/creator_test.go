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
	mockRepo     *mock_reconciliatonjob.MockCreatorRepository
	mockFileRepo *mock_reconciliatonjob.MockFileRepository
	svc          reconciliatonjob.Creator
}

func (s *ReconciliationJobCreatorTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.mockRepo = mock_reconciliatonjob.NewMockCreatorRepository(ctrl)
	s.mockFileRepo = mock_reconciliatonjob.NewMockFileRepository(ctrl)
	s.svc = reconciliatonjob.NewCreatorService(s.mockRepo, s.mockFileRepo)
}

func TestReconciliationJobCreatorTestSuite(t *testing.T) {
	suite.Run(t, new(ReconciliationJobCreatorTestSuite))
}

func (s *ReconciliationJobCreatorTestSuite) TestCreate() {
	ctx := context.Background()

	systemTrxPath := "/path/to/system/transaction.csv"
	bcaTrxPath := "/path/to/bca/transaction.csv"
	params := &reconciliatonjob.CreateParams{
		SystemTransactionCsv: &reconciliatonjob.File{
			Name: "system_transaction.csv",
		},
		DiscrepancyThreshold: 0.1,
		StartDate:            time.Now(),
		EndDate:              time.Now(),
		BankTransactionCsvs: []*reconciliatonjob.BankTransactionFile{
			{
				BankName: "BCA",
				File: &reconciliatonjob.File{
					Name: "bca_transaction.csv",
				},
			},
		},
	}
	bankTrxCsvPaths := []entity.BankTransactionCsv{
		{
			BankName: params.BankTransactionCsvs[0].BankName,
			FilePath: bcaTrxPath,
		},
	}
	dbParams := dbgen.CreateReconciliationJobParams{
		SystemTransactionCsvPath: systemTrxPath,
		DiscrepancyThreshold:     float64(params.DiscrepancyThreshold),
		StartDate:                params.StartDate,
		EndDate:                  params.EndDate,
	}
	dbParams.BankTransactionCsvPaths.Set(bankTrxCsvPaths)
	dbResult := dbgen.ReconciliationJob{
		ID:                       1,
		Status:                   "PENDING",
		SystemTransactionCsvPath: dbParams.SystemTransactionCsvPath,
		BankTransactionCsvPaths:  dbParams.BankTransactionCsvPaths,
		DiscrepancyThreshold:     dbParams.DiscrepancyThreshold,
		StartDate:                dbParams.StartDate,
		EndDate:                  dbParams.EndDate,
		CreatedAt:                now,
		UpdatedAt:                now,
	}
	jrResult := &entity.ReconciliationJob{
		ID:                       dbResult.ID,
		Status:                   entity.ReconciliationJobStatus(dbResult.Status),
		SystemTransactionCsvPath: systemTrxPath,
		BankTransactionCsvPaths:  bankTrxCsvPaths,
		DiscrepancyThreshold:     float32(dbResult.DiscrepancyThreshold),
		StartDate:                dbResult.StartDate,
		EndDate:                  dbResult.EndDate,
		CreatedAt:                dbResult.CreatedAt,
		UpdatedAt:                dbResult.UpdatedAt,
	}

	s.Run("success", func() {
		s.mockFileRepo.EXPECT().Store(ctx, params.SystemTransactionCsv).Return(systemTrxPath, nil)
		s.mockFileRepo.EXPECT().Store(ctx, params.BankTransactionCsvs[0].File).Return(bcaTrxPath, nil)
		s.mockRepo.EXPECT().CreateReconciliationJob(ctx, dbParams).Return(dbResult, nil)

		res, err := s.svc.Create(ctx, params)

		s.Nil(err)
		s.Equal(jrResult, res)
	})

	s.Run("error store system transaction csv", func() {
		s.mockFileRepo.EXPECT().Store(ctx, params.SystemTransactionCsv).Return("", assert.AnError)

		res, err := s.svc.Create(ctx, params)

		s.Nil(res)
		s.Equal(assert.AnError, err)
	})

	s.Run("error store bank transaction csv", func() {
		s.mockFileRepo.EXPECT().Store(ctx, params.SystemTransactionCsv).Return(systemTrxPath, nil)
		s.mockFileRepo.EXPECT().Store(ctx, params.BankTransactionCsvs[0].File).Return("", assert.AnError)

		res, err := s.svc.Create(ctx, params)

		s.Nil(res)
		s.Equal(assert.AnError, err)
	})

	s.Run("error create recon", func() {
		s.mockFileRepo.EXPECT().Store(ctx, params.SystemTransactionCsv).Return(systemTrxPath, nil)
		s.mockFileRepo.EXPECT().Store(ctx, params.BankTransactionCsvs[0].File).Return(bcaTrxPath, nil)
		s.mockRepo.EXPECT().CreateReconciliationJob(ctx, dbParams).Return(dbgen.ReconciliationJob{}, assert.AnError)

		res, err := s.svc.Create(ctx, params)

		s.Nil(res)
		s.Equal(assert.AnError, err)
	})
}
