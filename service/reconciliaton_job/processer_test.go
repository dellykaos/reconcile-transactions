package reconciliatonjob_test

import (
	"bytes"
	"context"
	"os"
	"testing"

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

	systemFile     *bytes.Buffer
	bcaFile        *bytes.Buffer
	briFile        *bytes.Buffer
	mockRepo       *mock_reconciliatonjob.MockProcesserRepository
	mockFileGetter *mock_reconciliatonjob.MockFileGetter
	svc            *reconciliatonjob.ProcesserService
}

func (s *ReconciliationJobProcessorTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.mockRepo = mock_reconciliatonjob.NewMockProcesserRepository(ctrl)
	s.mockFileGetter = mock_reconciliatonjob.NewMockFileGetter(ctrl)
	s.svc = reconciliatonjob.NewProcesserService(s.mockRepo, s.mockFileGetter)
	f, _ := os.ReadFile("../test/data/system_trx.csv")
	s.systemFile = bytes.NewBuffer(f)
	f, _ = os.ReadFile("../test/data/bca_trx.csv")
	s.bcaFile = bytes.NewBuffer(f)
	f, _ = os.ReadFile("../test/data/bri_trx.csv")
	s.briFile = bytes.NewBuffer(f)
}

func TestReconciliationJobProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ReconciliationJobProcessorTestSuite))
}

func (s *ReconciliationJobProcessorTestSuite) TestProcess() {
	ctx := context.Background()

	s.Run("success nothing to process", func() {
		s.mockRepo.EXPECT().ListPendingReconciliationJobs(ctx).Return([]dbgen.ReconciliationJob{}, nil)

		err := s.svc.Process(ctx)

		s.NoError(err)
	})

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
}
