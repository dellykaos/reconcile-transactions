package http

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/delly/amartha/common/logger"
	"github.com/delly/amartha/entity"
	"github.com/delly/amartha/handler/http/middleware"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	"github.com/h2non/filetype"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	allowedMimeType       = "text/csv"
	humanizeLimitFileSize = "10MB"
)

// ReconciliationJobHandler is a handler for reconciliation job
type ReconciliationJobHandler struct {
	finderService  reconciliatonjob.Finder
	creatorService reconciliatonjob.Creator
	log            *zap.Logger
}

// NewReconciliationJobHandler create new reconciliation job handler, it used to create new reconciliation job, get reconciliation job by id, and get all reconciliation job
func NewReconciliationJobHandler(finderService reconciliatonjob.Finder,
	creatorService reconciliatonjob.Creator) *ReconciliationJobHandler {
	return &ReconciliationJobHandler{
		finderService:  finderService,
		creatorService: creatorService,
		log:            zap.L().With(zap.String("handler", "reconciliation_job")),
	}
}

// Register register reconciliation job handler to router
func (h *ReconciliationJobHandler) Register(router *httprouter.Router) {
	router.GET("/reconciliations", middleware.PrependMiddleware(h.GetAllReconciliationJob, middleware.WithLogger))
	router.GET("/reconciliations/:id", middleware.PrependMiddleware(h.GetReconciliationJobByID, middleware.WithLogger))
	router.POST("/reconciliations", middleware.PrependMiddleware(h.CreateReconciliationJob, middleware.WithLogger))
}

// GetReconciliationJobByID get reconciliation job by id
func (h *ReconciliationJobHandler) GetReconciliationJobByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log := logger.WithMethod(h.log, "GetReconciliationJobByID")
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		log.Error("invalid id", zap.String("id", p.ByName("id")), zap.Error(err))
		writeBadRequest(w, "invalid id")
		return
	}

	rj, err := h.finderService.FindByID(r.Context(), id)
	if err != nil {
		log.Error("failed to get reconciliation job by id", zap.Error(err), zap.Int64("id", id))
		writeInternalServerError(w)
		return
	}
	if rj == nil {
		log.Error("reconciliation job not found", zap.Int64("id", id))
		writeNotFound(w, "reconciliation job not found")
		return
	}

	writeJSON(w, http.StatusOK, rj, nil)
}

// GetAllReconciliationJob get all reconciliation job
func (h *ReconciliationJobHandler) GetAllReconciliationJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log := logger.WithMethod(h.log, "GetAllReconciliationJob")
	pagination := getPagination(r)
	total, err := h.finderService.Count(r.Context())
	if err != nil {
		log.Error("failed to count reconciliation job", zap.Error(err))
		writeInternalServerError(w)
		return
	}
	pagination.Total = int32(total)
	if total == 0 {
		writeJSON(w, http.StatusOK, []*entity.ReconciliationJob{}, pagination)
		return
	}

	rjs, err := h.finderService.FindAll(r.Context(), pagination.Limit, pagination.Offset)
	if err != nil {
		log.Error("failed to list reconciliation jobs", zap.Error(err))
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, rjs, pagination)
}

// CreateReconciliationJob create new reconciliation job
func (h *ReconciliationJobHandler) CreateReconciliationJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log := logger.WithMethod(h.log, "CreateReconciliationJob")
	params, err := h.parseCreateReconciliationJobParams(r)
	if err != nil {
		log.Error("failed to parse create reconciliation job params", zap.Error(err))
		writeBadRequest(w, err.Error())
		return
	}

	rj, err := h.creatorService.Create(r.Context(), params)
	if err != nil {
		log.Error("failed to create reconciliation job", zap.Error(err))
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusCreated, rj, nil)
}

func (h *ReconciliationJobHandler) parseCreateReconciliationJobParams(r *http.Request) (*reconciliatonjob.CreateParams, error) {
	params, err := h.buildFileParams(r)
	if err != nil {
		return nil, err
	}

	startDate := parseDate(r.FormValue("start_date"))
	endDate := parseDate(r.FormValue("end_date"))
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}
	params.StartDate = startDate
	params.EndDate = endDate

	discrepancyThreshold := parseFloat32(r.FormValue("discrepancy_threshold"))
	if discrepancyThreshold < 0 {
		discrepancyThreshold = 0
	}
	params.DiscrepancyThreshold = discrepancyThreshold

	return params, nil
}

func (h *ReconciliationJobHandler) buildFileParams(r *http.Request) (*reconciliatonjob.CreateParams, error) {
	if err := r.ParseMultipartForm(entity.LimitContentSize); err != nil {
		return nil, err
	}
	if r.ContentLength > entity.LimitContentSize {
		return nil, errors.New("content size more than 100mb")
	}

	file, err := h.buildSystemTrxFile(r.MultipartForm)
	if err != nil {
		return nil, err
	}

	bankFiles, err := h.buildBankTrxFiles(r.MultipartForm)
	if err != nil {
		return nil, err
	}

	params := &reconciliatonjob.CreateParams{
		SystemTransactionCsv: file,
		BankTransactionCsvs:  bankFiles,
	}

	return params, nil
}

func (h *ReconciliationJobHandler) buildSystemTrxFile(form *multipart.Form) (*reconciliatonjob.File, error) {
	systemTrxFile := form.File["system_transaction_file"]
	if len(systemTrxFile) == 0 {
		return nil, errors.New("system transaction file is required")
	}

	buf, err := h.validateCSVFile(systemTrxFile[0])
	if err != nil {
		return nil, err
	}

	return &reconciliatonjob.File{
		Name: systemTrxFile[0].Filename,
		Buf:  buf,
	}, nil
}

func (h *ReconciliationJobHandler) buildBankTrxFiles(form *multipart.Form) ([]*reconciliatonjob.BankTransactionFile, error) {
	bankNames := form.Value["bank_names"]
	bankTrxFiles := form.File["bank_transaction_files"]
	if len(bankTrxFiles) == 0 {
		return nil, ErrBankTrxFileEmpty
	}
	if len(bankNames) != len(bankTrxFiles) {
		return nil, ErrBankFileAndNameLengthNotMatch
	}

	result := []*reconciliatonjob.BankTransactionFile{}
	for idx, file := range bankTrxFiles {
		buf, err := h.validateCSVFile(file)
		if err != nil {
			return nil, err
		}
		bankFile := &reconciliatonjob.BankTransactionFile{
			BankName: bankNames[idx],
			File: &reconciliatonjob.File{
				Name: file.Filename,
				Buf:  buf,
			},
		}
		result = append(result, bankFile)
	}

	return result, nil
}

func (h *ReconciliationJobHandler) validateCSVFile(file *multipart.FileHeader) (*bytes.Buffer, error) {
	if !isCSVExtension(file.Filename) {
		return nil, ErrExtensionFileInvalid(file.Filename)
	}
	if file.Size > entity.LimitCSVSize {
		return nil, ErrFileSizeExceedLimit(file.Filename, humanizeLimitFileSize)
	}

	byteFileBuf := bytes.NewBuffer(nil)
	f, err := file.Open()
	if err != nil {
		return nil, ErrFileCannotBeAccessed(file.Filename)
	}
	defer f.Close()

	if _, err = io.Copy(byteFileBuf, f); err != nil {
		return bytes.NewBuffer(nil), ErrFileCannotBeAccessed(file.Filename)
	}

	fileType, err := filetype.Match(byteFileBuf.Bytes())
	if err != nil {
		return bytes.NewBuffer(nil), ErrExtensionFileUnknown(file.Filename)
	}

	if fileType.MIME.Value == allowedMimeType {
		return bytes.NewBuffer(nil), ErrExtensionFileInvalid(file.Filename)
	}

	return byteFileBuf, nil
}
