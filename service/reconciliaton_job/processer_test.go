package reconciliatonjob_test

import (
	"context"
	"testing"

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
}
