package reconciliatonjob_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/delly/amartha/entity"
	filestorage "github.com/delly/amartha/repository/file_storage"
	dbgen "github.com/delly/amartha/repository/postgresql"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	mock_reconciliatonjob "github.com/delly/amartha/test/mock/service/reconciliaton_job"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ReconciliationJobProcessorTestSuite struct {
	suite.Suite

	mockRepo       *mock_reconciliatonjob.MockProcesserRepository
	mockFileGetter *mock_reconciliatonjob.MockFileGetter
	svc            *reconciliatonjob.ProcesserService
}

func (s *ReconciliationJobProcessorTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.mockRepo = mock_reconciliatonjob.NewMockProcesserRepository(ctrl)
	s.mockFileGetter = mock_reconciliatonjob.NewMockFileGetter(ctrl)
	s.svc = reconciliatonjob.NewProcesserService(s.mockRepo, s.mockFileGetter)
}

func TestReconciliationJobProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ReconciliationJobProcessorTestSuite))
}

func (s *ReconciliationJobProcessorTestSuite) TestProcess_Failed() {
	ctx := context.Background()

	s.Run("error get list pending reconciliation jobs", func() {
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{}, assert.AnError)

		err := s.svc.Process(ctx)

		s.Error(err)
	})

	s.Run("error get file system trx", func() {
		rj := dbReconJob
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(nil, assert.AnError)
		s.mockRepo.EXPECT().SaveFailedReconciliationJob(ctx, gomock.Any()).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.Nil(err)
	})

	s.Run("error get file bank trx", func() {
		rj := dbReconJob
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
		}
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(nil, assert.AnError)
		s.mockRepo.EXPECT().SaveFailedReconciliationJob(ctx, gomock.Any()).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.Nil(err)
	})

	s.Run("error read file system trx", func() {
		rj := dbReconJob
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  nil,
		}
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(&filestorage.File{}, nil)
		s.mockRepo.EXPECT().SaveFailedReconciliationJob(ctx, gomock.Any()).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.Nil(err)
	})

	s.Run("error convert amount file system trx to entity", func() {
		rj := dbReconJob
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  bytes.NewBuffer([]byte("\"\",\"abc\",\n")),
		}
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(&filestorage.File{}, nil)
		s.mockRepo.EXPECT().SaveFailedReconciliationJob(ctx, gomock.Any()).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.Nil(err)
	})

	s.Run("error convert transaction type file system trx to entity", func() {
		rj := dbReconJob
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  bytes.NewBuffer([]byte("\"\",\"10000\",\"abc\",\n")),
		}
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(&filestorage.File{}, nil)
		s.mockRepo.EXPECT().SaveFailedReconciliationJob(ctx, gomock.Any()).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.Nil(err)
	})

	s.Run("error convert time file system trx to entity", func() {
		rj := dbReconJob
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  bytes.NewBuffer([]byte("\"\",\"10000\",\"DEBIT\",\"abc\"\n")),
		}
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(&filestorage.File{}, nil)
		s.mockRepo.EXPECT().SaveFailedReconciliationJob(ctx, gomock.Any()).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.Nil(err)
	})

	s.Run("error read file bank trx", func() {
		rj := dbReconJob
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  bytes.NewBuffer([]byte("\"\",\"10000\",\"DEBIT\",\"" + strLastWeek + "\"\n")),
		}
		fsBankTrx := &filestorage.File{
			Name: "bank_transaction.csv",
			Buf:  nil,
		}
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(fsBankTrx, nil)
		s.mockRepo.EXPECT().SaveFailedReconciliationJob(ctx, gomock.Any()).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.Nil(err)
	})

	s.Run("error convert amount file bank trx to entity", func() {
		rj := dbReconJob
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  bytes.NewBuffer([]byte("\"\",\"10000\",\"DEBIT\",\"" + strLastWeek + "\"\n")),
		}
		fsBankTrx := &filestorage.File{
			Name: "bank_transaction.csv",
			Buf:  bytes.NewBuffer([]byte("\"\",\"abc\",\n")),
		}
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(fsBankTrx, nil)
		s.mockRepo.EXPECT().SaveFailedReconciliationJob(ctx, gomock.Any()).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.Nil(err)
	})

	s.Run("error convert time file bank trx to entity", func() {
		rj := dbReconJob
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  bytes.NewBuffer([]byte("\"\",\"10000\",\"DEBIT\",\"2021-01-01T01:01:01+07:00\"\n")),
		}
		fsBankTrx := &filestorage.File{
			Name: "bank_transaction.csv",
			Buf:  bytes.NewBuffer([]byte("\"\",\"10000\",\"abc\"\n")),
		}
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(fsBankTrx, nil)
		s.mockRepo.EXPECT().SaveFailedReconciliationJob(ctx, gomock.Any()).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.Nil(err)
	})
}

func (s *ReconciliationJobProcessorTestSuite) TestProcess_Success() {
	ctx := context.Background()

	s.Run("success nothing to process", func() {
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.NoError(err)
	})

	s.Run("success process reconciliation job with system trx matched and unmatched", func() {
		rj := dbReconJob
		rj.StartDate = time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
		rj.EndDate = time.Date(2024, 11, 23, 0, 0, 0, 0, time.UTC)
		rj.DiscrepancyThreshold = 0
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  fetchSystemFile("system_trx.csv"),
		}
		fsBankTrx := &filestorage.File{
			Name: "bank_transaction.csv",
			Buf:  fetchSystemFile("bca_trx.csv"),
		}
		expectedResult := entity.ReconciliationResult{
			TotalTransactionProcessed: 13,
			TotalTransactionMatched:   9,
			TotalTransactionUnmatched: 4,
			TotalDiscrepancyAmount:    212963022,
			MissingTransactions: []entity.Transaction{
				{
					ID:     "ABC-127",
					Amount: 123000000,
					Type:   entity.TxTypeDebit,
					Time:   parseTime("2024-11-05T11:24:00Z"),
				},
				{
					ID:     "ABC-128",
					Amount: 54200000,
					Type:   entity.TxTypeCredit,
					Time:   parseTime("2024-11-06T11:22:03Z"),
				},
				{
					ID:     "ABC-129",
					Amount: 23450000,
					Type:   entity.TxTypeCredit,
					Time:   parseTime("2024-11-07T12:33:22Z"),
				},
				{
					ID:     "ABC-130",
					Amount: 12313022,
					Type:   entity.TxTypeDebit,
					Time:   parseTime("2024-11-07T14:00:23Z"),
				},
			},
			MissingBankTransactions: map[string][]entity.Transaction{},
		}
		saveParams := dbgen.SaveSuccessReconciliationJobParams{
			ID: rj.ID,
		}
		saveParams.Result.Set(expectedResult)
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(fsBankTrx, nil)
		s.mockRepo.EXPECT().SaveSuccessReconciliationJob(ctx, saveParams).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.NoError(err)
	})

	s.Run("success process reconciliation job without discrepancy", func() {
		rj := dbReconJob
		rj.StartDate = time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
		rj.EndDate = time.Date(2024, 11, 23, 0, 0, 0, 0, time.UTC)
		rj.DiscrepancyThreshold = 0
		bankCsvs := []entity.BankTransactionCsv{
			{
				BankName: "BCA",
				FilePath: "path/to/bca_transaction.csv",
			},
			{
				BankName: "BRI",
				FilePath: "path/to/bri_transaction.csv",
			},
		}
		rj.BankTransactionCsvPaths.Set(bankCsvs)
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  fetchSystemFile("system_trx.csv"),
		}
		fsBankBcaTrx := &filestorage.File{
			Name: "bca_transaction.csv",
			Buf:  fetchSystemFile("bca_trx.csv"),
		}
		fsBankBriTrx := &filestorage.File{
			Name: "bri_transaction.csv",
			Buf:  fetchSystemFile("bri_trx.csv"),
		}
		expectedResult := entity.ReconciliationResult{
			TotalTransactionProcessed: 13,
			TotalTransactionMatched:   13,
			TotalTransactionUnmatched: 0,
			TotalDiscrepancyAmount:    0,
			MissingTransactions:       []entity.Transaction{},
			MissingBankTransactions:   map[string][]entity.Transaction{},
		}
		saveParams := dbgen.SaveSuccessReconciliationJobParams{
			ID: rj.ID,
		}
		saveParams.Result.Set(expectedResult)
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, bankCsvs[0].FilePath).Return(fsBankBcaTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, bankCsvs[1].FilePath).Return(fsBankBriTrx, nil)
		s.mockRepo.EXPECT().SaveSuccessReconciliationJob(ctx, saveParams).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.NoError(err)
	})

	s.Run("success process reconciliation job with system and bank unmatched", func() {
		rj := dbReconJob
		rj.StartDate = time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
		rj.EndDate = time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC)
		rj.DiscrepancyThreshold = 0
		fsSystemTrx := &filestorage.File{
			Name: "system_transaction.csv",
			Buf:  fetchSystemFile("system_trx_1.csv"),
		}
		fsBankTrx := &filestorage.File{
			Name: "bank_transaction.csv",
			Buf:  fetchSystemFile("bca_trx_1.csv"),
		}
		expectedResult := entity.ReconciliationResult{
			TotalTransactionProcessed: 1,
			TotalTransactionMatched:   0,
			TotalTransactionUnmatched: 1,
			TotalDiscrepancyAmount:    192131,
			MissingTransactions: []entity.Transaction{
				{
					ID:     "ABC-123",
					Amount: 150000,
					Type:   entity.TxTypeCredit,
					Time:   parseTime("2024-11-01T02:00:00Z"),
				},
			},
			MissingBankTransactions: map[string][]entity.Transaction{
				"BCA": {
					{
						ID:     "BCA-133",
						Amount: 42131,
						Type:   entity.TxTypeDebit,
						Time:   parseTime("2024-11-25T00:00:00Z"),
					},
				},
			},
		}
		saveParams := dbgen.SaveSuccessReconciliationJobParams{
			ID: rj.ID,
		}
		saveParams.Result.Set(expectedResult)
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{rj}, nil)
		s.mockFileGetter.EXPECT().Get(ctx, rj.SystemTransactionCsvPath).Return(fsSystemTrx, nil)
		s.mockFileGetter.EXPECT().Get(ctx, entityReconJob.BankTransactionCsvPaths[0].FilePath).Return(fsBankTrx, nil)
		s.mockRepo.EXPECT().SaveSuccessReconciliationJob(ctx, saveParams).Return(dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.NoError(err)
	})
}

func parseTime(t string) time.Time {
	res, _ := time.Parse(time.RFC3339, t)
	return res
}

func fetchSystemFile(filename string) *bytes.Buffer {
	f, _ := os.ReadFile("../../test/data/" + filename)
	return bytes.NewBuffer(f)
}
