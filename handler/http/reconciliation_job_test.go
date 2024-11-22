package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/delly/amartha/entity"
	handler "github.com/delly/amartha/handler/http"
	mock_reconciliatonjob "github.com/delly/amartha/test/mock/service/reconciliaton_job"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	id             = int64(1)
	now            = time.Now()
	entityReconJob = &entity.ReconciliationJob{
		ID:                       id,
		Status:                   entity.ReconciliationJobStatus("SUCCESS"),
		SystemTransactionCsvPath: "path_to_file",
		DiscrepancyThreshold:     0.1,
		StartDate:                now,
		EndDate:                  now,
		CreatedAt:                now,
		UpdatedAt:                now,
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
)

type ReconciliationJobHandlerTestSuite struct {
	suite.Suite
	router      *httprouter.Router
	mockService *mock_reconciliatonjob.MockFinder
	handler     *handler.ReconciliationJobHandler
}

func (s *ReconciliationJobHandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.mockService = mock_reconciliatonjob.NewMockFinder(ctrl)
	s.handler = handler.NewReconciliationJobHandler(s.mockService)

	s.router = httprouter.New()
	s.handler.Register(s.router)
}

func TestReconciliationJobHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ReconciliationJobHandlerTestSuite))
}

func (s *ReconciliationJobHandlerTestSuite) TestGetReconciliationJobByID() {
	ctx := context.Background()

	req, _ := http.NewRequest(http.MethodGet, "/reconciliations/1", nil)
	s.Run("success", func() {
		s.mockService.EXPECT().FindByID(ctx, id).Return(entityReconJob, nil)

		resp := s.executeReq(req)

		jsonRecon, _ := json.Marshal(entityReconJob)
		bodyJson := resp.Body.String()
		s.Equal(http.StatusOK, resp.Code)
		s.Contains(bodyJson, string(jsonRecon))
	})

	s.Run("not found", func() {
		s.mockService.EXPECT().FindByID(ctx, id).Return(nil, nil)

		resp := s.executeReq(req)

		bodyJson := resp.Body.String()
		s.Equal(http.StatusNotFound, resp.Code)
		s.Contains(bodyJson, "reconciliation job not found")
	})

	s.Run("internal server error", func() {
		s.mockService.EXPECT().FindByID(ctx, id).Return(nil, assert.AnError)

		resp := s.executeReq(req)

		s.Equal(http.StatusInternalServerError, resp.Code)
	})

	s.Run("invalid id", func() {
		req, _ := http.NewRequest(http.MethodGet, "/reconciliations/invalid", nil)

		resp := s.executeReq(req)

		bodyJson := resp.Body.String()
		s.Equal(http.StatusBadRequest, resp.Code)
		s.Contains(bodyJson, "invalid id")
	})
}

func (s *ReconciliationJobHandlerTestSuite) TestGetAllReconciliationJob() {
	ctx := context.Background()

	req, _ := http.NewRequest(http.MethodGet, "/reconciliations", nil)
	s.Run("success", func() {
		total := int64(1)
		s.mockService.EXPECT().Count(ctx).Return(total, nil)
		s.mockService.EXPECT().FindAll(ctx, int32(10), int32(0)).Return([]*entity.ReconciliationJob{entityReconJob}, nil)

		resp := s.executeReq(req)

		jsonRecon, _ := json.Marshal([]*entity.ReconciliationJob{entityReconJob})
		bodyJson := resp.Body.String()
		s.Equal(http.StatusOK, resp.Code)
		s.Contains(bodyJson, string(jsonRecon))
	})

	s.Run("no data", func() {
		total := int64(0)
		s.mockService.EXPECT().Count(ctx).Return(total, nil)

		resp := s.executeReq(req)

		bodyJson := resp.Body.String()
		s.Equal(http.StatusOK, resp.Code)
		s.Contains(bodyJson, "[]")
	})

	s.Run("invalid offset and limit", func() {
		req, _ := http.NewRequest(http.MethodGet, "/reconciliations?offset=-1&limit=1000", nil)
		total := int64(1)
		s.mockService.EXPECT().Count(ctx).Return(total, nil)
		s.mockService.EXPECT().FindAll(ctx, int32(100), int32(0)).Return([]*entity.ReconciliationJob{entityReconJob}, nil)

		resp := s.executeReq(req)

		bodyJson := resp.Body.String()
		jsonRecon, _ := json.Marshal([]*entity.ReconciliationJob{entityReconJob})
		s.Equal(http.StatusOK, resp.Code)
		s.Contains(bodyJson, string(jsonRecon))
	})

	s.Run("error on count", func() {
		s.mockService.EXPECT().Count(ctx).Return(int64(0), assert.AnError)

		resp := s.executeReq(req)

		s.Equal(http.StatusInternalServerError, resp.Code)
	})

	s.Run("error on find all", func() {
		total := int64(1)
		s.mockService.EXPECT().Count(ctx).Return(total, nil)
		s.mockService.EXPECT().FindAll(ctx, int32(10), int32(0)).Return(nil, assert.AnError)

		resp := s.executeReq(req)

		s.Equal(http.StatusInternalServerError, resp.Code)
	})
}

func (s *ReconciliationJobHandlerTestSuite) executeReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	return rr
}
