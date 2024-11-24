package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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
	id                   = int64(1)
	now                  = time.Now()
	simpleEntityReconJob = &entity.SimpleReconciliationJob{
		ID:                       id,
		Status:                   entity.ReconciliationJobStatus("SUCCESS"),
		SystemTransactionCsvPath: "path_to_file",
		BankTransactionCsvPaths: []entity.BankTransactionCsv{
			{
				BankName: "BCA",
				FilePath: "path_to_file_bca",
			},
		},
		DiscrepancyThreshold: 0.1,
		StartDate:            now,
		EndDate:              now,
	}
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
		Result: &entity.ReconciliationResult{
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
	router             *httprouter.Router
	mockFinderService  *mock_reconciliatonjob.MockFinder
	mockCreatorService *mock_reconciliatonjob.MockCreator
	handler            *handler.ReconciliationJobHandler
}

func (s *ReconciliationJobHandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.mockFinderService = mock_reconciliatonjob.NewMockFinder(ctrl)
	s.mockCreatorService = mock_reconciliatonjob.NewMockCreator(ctrl)
	s.handler = handler.NewReconciliationJobHandler(s.mockFinderService, s.mockCreatorService)

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
		s.mockFinderService.EXPECT().FindByID(ctx, id).Return(entityReconJob, nil)

		resp := s.executeReq(req)

		jsonRecon, _ := json.Marshal(entityReconJob)
		bodyJson := resp.Body.String()
		s.Equal(http.StatusOK, resp.Code)
		s.Contains(bodyJson, string(jsonRecon))
	})

	s.Run("not found", func() {
		s.mockFinderService.EXPECT().FindByID(ctx, id).Return(nil, nil)

		resp := s.executeReq(req)

		bodyJson := resp.Body.String()
		s.Equal(http.StatusNotFound, resp.Code)
		s.Contains(bodyJson, "reconciliation job not found")
	})

	s.Run("internal server error", func() {
		s.mockFinderService.EXPECT().FindByID(ctx, id).Return(nil, assert.AnError)

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
		s.mockFinderService.EXPECT().Count(ctx).Return(total, nil)
		s.mockFinderService.EXPECT().FindAll(ctx, int32(10), int32(0)).Return([]*entity.SimpleReconciliationJob{simpleEntityReconJob}, nil)

		resp := s.executeReq(req)

		jsonRecon, _ := json.Marshal([]*entity.SimpleReconciliationJob{simpleEntityReconJob})
		bodyJson := resp.Body.String()
		s.Equal(http.StatusOK, resp.Code)
		s.Contains(bodyJson, string(jsonRecon))
	})

	s.Run("no data", func() {
		total := int64(0)
		s.mockFinderService.EXPECT().Count(ctx).Return(total, nil)

		resp := s.executeReq(req)

		bodyJson := resp.Body.String()
		s.Equal(http.StatusOK, resp.Code)
		s.Contains(bodyJson, "[]")
	})

	s.Run("invalid offset and limit", func() {
		req, _ := http.NewRequest(http.MethodGet, "/reconciliations?offset=-1&limit=1000", nil)
		total := int64(1)
		s.mockFinderService.EXPECT().Count(ctx).Return(total, nil)
		s.mockFinderService.EXPECT().FindAll(ctx, int32(100), int32(0)).Return([]*entity.SimpleReconciliationJob{simpleEntityReconJob}, nil)

		resp := s.executeReq(req)

		bodyJson := resp.Body.String()
		jsonRecon, _ := json.Marshal([]*entity.SimpleReconciliationJob{simpleEntityReconJob})
		s.Equal(http.StatusOK, resp.Code)
		s.Contains(bodyJson, string(jsonRecon))
	})

	s.Run("error on count", func() {
		s.mockFinderService.EXPECT().Count(ctx).Return(int64(0), assert.AnError)

		resp := s.executeReq(req)

		s.Equal(http.StatusInternalServerError, resp.Code)
	})

	s.Run("error on find all", func() {
		total := int64(1)
		s.mockFinderService.EXPECT().Count(ctx).Return(total, nil)
		s.mockFinderService.EXPECT().FindAll(ctx, int32(10), int32(0)).Return(nil, assert.AnError)

		resp := s.executeReq(req)

		s.Equal(http.StatusInternalServerError, resp.Code)
	})
}

func (s *ReconciliationJobHandlerTestSuite) TestCreateReconciliationJob() {
	ctx := context.Background()

	s.Run("success", func() {
		req := s.buildCreatorReq(func(mw *multipart.Writer) {
			mw.WriteField("discrepancy_threshold", "0.1")
			mw.WriteField("start_date", now.Format("2006-01-02"))
			mw.WriteField("end_date", now.Format("2006-01-02"))
			mw.WriteField("bank_names", "BCA")
			s.createFormFile(mw, "system_transaction_file", "system_trx.csv")
			s.createFormFile(mw, "bank_transaction_files", "bca_trx.csv")
		})
		s.mockCreatorService.EXPECT().Create(ctx, gomock.Any()).Return(entityReconJob, nil)

		resp := s.executeReq(req)

		s.Equal(http.StatusCreated, resp.Code)
	})

	s.Run("error on create", func() {
		req := s.buildCreatorReq(func(mw *multipart.Writer) {
			mw.WriteField("discrepancy_threshold", "0.1")
			mw.WriteField("start_date", now.Format("2006-01-02"))
			mw.WriteField("end_date", now.Format("2006-01-02"))
			mw.WriteField("bank_names", "BCA")
			s.createFormFile(mw, "system_transaction_file", "system_trx.csv")
			s.createFormFile(mw, "bank_transaction_files", "bca_trx.csv")
		})
		s.mockCreatorService.EXPECT().Create(ctx, gomock.Any()).Return(nil, assert.AnError)

		resp := s.executeReq(req)

		s.Equal(http.StatusInternalServerError, resp.Code)
	})

	s.Run("empty payload", func() {
		req, _ := http.NewRequest(http.MethodPost, "/reconciliations", nil)

		resp := s.executeReq(req)

		s.Equal(http.StatusBadRequest, resp.Code)
	})

	s.Run("invalid discrepancy threshold", func() {
		req := s.buildCreatorReq(func(mw *multipart.Writer) {
			mw.WriteField("discrepancy_threshold", "-0.1")
			mw.WriteField("start_date", now.Format("2006-01-02"))
			mw.WriteField("end_date", now.Format("2006-01-02"))
			mw.WriteField("bank_names", "BCA")
			s.createFormFile(mw, "system_transaction_file", "system_trx.csv")
			s.createFormFile(mw, "bank_transaction_files", "bca_trx.csv")
		})
		s.mockCreatorService.EXPECT().Create(ctx, gomock.Any()).Return(entityReconJob, nil)

		resp := s.executeReq(req)

		s.Equal(http.StatusCreated, resp.Code)
	})

	s.Run("invalid start date", func() {
		req := s.buildCreatorReq(func(mw *multipart.Writer) {
			mw.WriteField("discrepancy_threshold", "0.1")
			mw.WriteField("start_date", now.AddDate(0, 0, 1).Format("2006-01-02"))
			mw.WriteField("end_date", now.Format("2006-01-02"))
			mw.WriteField("bank_names", "BCA")
			s.createFormFile(mw, "system_transaction_file", "system_trx.csv")
			s.createFormFile(mw, "bank_transaction_files", "bca_trx.csv")
		})

		resp := s.executeReq(req)

		s.Equal(http.StatusBadRequest, resp.Code)
	})

	s.Run("invalid length of bank names and bank transaction files", func() {
		req := s.buildCreatorReq(func(mw *multipart.Writer) {
			mw.WriteField("discrepancy_threshold", "0.1")
			mw.WriteField("start_date", now.Format("2006-01-02"))
			mw.WriteField("end_date", now.Format("2006-01-02"))
			mw.WriteField("bank_names", "BCA")
			mw.WriteField("bank_names", "BRI")
			s.createFormFile(mw, "system_transaction_file", "system_trx.csv")
			s.createFormFile(mw, "bank_transaction_files", "bca_trx.csv")
		})

		resp := s.executeReq(req)

		s.Equal(http.StatusBadRequest, resp.Code)
	})

	s.Run("empty bank transaction files", func() {
		req := s.buildCreatorReq(func(mw *multipart.Writer) {
			mw.WriteField("discrepancy_threshold", "0.1")
			mw.WriteField("start_date", now.Format("2006-01-02"))
			mw.WriteField("end_date", now.Format("2006-01-02"))
			mw.WriteField("bank_names", "BCA")
			s.createFormFile(mw, "system_transaction_file", "system_trx.csv")
		})

		resp := s.executeReq(req)

		s.Equal(http.StatusBadRequest, resp.Code)
	})

	s.Run("invalid system transaction file", func() {
		req := s.buildCreatorReq(func(mw *multipart.Writer) {
			mw.WriteField("discrepancy_threshold", "0.1")
			mw.WriteField("start_date", now.Format("2006-01-02"))
			mw.WriteField("end_date", now.Format("2006-01-02"))
			mw.WriteField("bank_names", "BCA")
			s.createFormFile(mw, "system_transaction_file", "test.json")
			s.createFormFile(mw, "bank_transaction_files", "bca_trx.csv")
		})

		resp := s.executeReq(req)

		s.Equal(http.StatusBadRequest, resp.Code)
	})

	s.Run("invalid bank transaction file", func() {
		req := s.buildCreatorReq(func(mw *multipart.Writer) {
			mw.WriteField("discrepancy_threshold", "0.1")
			mw.WriteField("start_date", now.Format("2006-01-02"))
			mw.WriteField("end_date", now.Format("2006-01-02"))
			mw.WriteField("bank_names", "BCA")
			s.createFormFile(mw, "system_transaction_file", "system_trx.csv")
			s.createFormFile(mw, "bank_transaction_files", "test.json")
		})

		resp := s.executeReq(req)

		s.Equal(http.StatusBadRequest, resp.Code)
	})
}

func (s *ReconciliationJobHandlerTestSuite) executeReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	return rr
}

func (s *ReconciliationJobHandlerTestSuite) buildCreatorReq(fn func(mw *multipart.Writer)) *http.Request {
	var buf bytes.Buffer
	mWriter := multipart.NewWriter(&buf)
	fn(mWriter)
	mWriter.Close()

	req, _ := http.NewRequest(http.MethodPost, "/reconciliations", &buf)
	req.Header.Set("Content-Type", mWriter.FormDataContentType())

	return req
}

func (s *ReconciliationJobHandlerTestSuite) createFormFile(mw *multipart.Writer, fieldName, fileName string) {
	f, _ := mw.CreateFormFile(fieldName, fileName)
	file, _ := os.Open("../../test/data/" + fileName)
	io.Copy(f, file)
	file.Close()
}
