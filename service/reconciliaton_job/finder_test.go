package reconciliatonjob_test

import (
	"context"
	"testing"
	"time"

	"github.com/delly/amartha/entity"
	dbgen "github.com/delly/amartha/repository/postgresql"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	mock_reconciliatonjob "github.com/delly/amartha/test/mock/service/reconciliaton_job"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type FinderTestSuite struct {
	suite.Suite
	repo *mock_reconciliatonjob.MockFinderRepository
	svc  *reconciliatonjob.FinderService
}

func (s *FinderTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.repo = mock_reconciliatonjob.NewMockFinderRepository(ctrl)
	s.svc = reconciliatonjob.NewFinderService(s.repo)
}

func TestFinderTestSuite(t *testing.T) {
	suite.Run(t, new(FinderTestSuite))
}

func (s *FinderTestSuite) TestFindByID() {
	ctx := context.Background()
	id := int64(1)
	lastMonth := time.Now().AddDate(0, -1, 0)
	yesterday := time.Now().AddDate(0, 0, -1)

	s.Run("success", func() {
		rj := dbgen.ReconciliationJob{
			ID:                       id,
			Status:                   "SUCCESS",
			SystemTransactionCsvPath: "path_to_file",
			DiscrepancyThreshold:     0.1,
			StartDate:                lastMonth,
			EndDate:                  yesterday,
			CreatedAt:                time.Now(),
			UpdatedAt:                time.Now(),
			BankTransactionCsvPaths: pgtype.JSONB{
				Status: pgtype.Present,
				Bytes:  []byte(`[{"bank_name": "BCA", "file_path": "path_to_file_bca"}]`),
			},
			Result: pgtype.JSONB{
				Status: pgtype.Present,
				Bytes: []byte(`
				{
					"total_transaction_processed": 10,
					"total_transaction_matched": 5,
					"total_transaction_unmatched": 5,
					"total_discrepancy_amount": 1000,
					"missing_transactions": [
						{
							"id": "1",
							"amount": 1000,
							"type":
							"DEBIT",
							"time": "2021-01-02T00:00:00Z"
						}
					],
					"missing_bank_transactions": {
						"BCA": [
							{
								"id": "1",
								"amount": 1000,
								"type": "DEBIT",
								"time": "2021-01-01T00:00:00Z"
							}
						]
					}
				}`),
			},
		}
		expectedRJ := &entity.ReconciliationJob{
			ID:                       id,
			Status:                   entity.ReconciliationJobStatus("SUCCESS"),
			SystemTransactionCsvPath: "path_to_file",
			DiscrepancyThreshold:     0.1,
			StartDate:                lastMonth,
			EndDate:                  yesterday,
			CreatedAt:                rj.CreatedAt,
			UpdatedAt:                rj.UpdatedAt,
			BankTransactionCsvPaths: []entity.BankTransactionCsv{
				{
					BankName: "BCA",
					FilePath: "path_to_file_bca",
				},
			},
			Result: entity.ReconciliationResult{
				TotalTransactionProcessed: 10,
				TotalTransactionMatched:   5,
				TotalTransactionUnmatched: 5,
				TotalDiscrepancyAmount:    1000,
				MissingTransactions: []entity.Transaction{
					{
						ID:     "1",
						Amount: 1000,
						Type:   entity.TxTypeDebit,
						Time:   time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					},
				},
				MissingBankTransactions: map[string][]entity.Transaction{
					"BCA": {
						{
							ID:     "1",
							Amount: 1000,
							Type:   entity.TxTypeDebit,
							Time:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		}

		s.repo.EXPECT().GetReconciliationJobById(ctx, id).Return(rj, nil)

		res, err := s.svc.FindByID(ctx, id)
		s.NoError(err)
		s.Equal(expectedRJ, res)
	})

	s.Run("error", func() {
		s.repo.EXPECT().GetReconciliationJobById(ctx, id).Return(dbgen.ReconciliationJob{}, assert.AnError)

		res, err := s.svc.FindByID(ctx, id)
		s.Error(err)
		s.Nil(res)
	})
}
