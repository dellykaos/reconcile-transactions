package http

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/delly/amartha/entity"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	"github.com/julienschmidt/httprouter"
)

// ReconciliationJobHandler is a handler for reconciliation job
type ReconciliationJobHandler struct {
	finderService  reconciliatonjob.Finder
	creatorService reconciliatonjob.Creator
	logger         *log.Logger
}

// NewReconciliationJobHandler create new reconciliation job handler, it used to create new reconciliation job, get reconciliation job by id, and get all reconciliation job
func NewReconciliationJobHandler(finderService reconciliatonjob.Finder,
	creatorService reconciliatonjob.Creator) *ReconciliationJobHandler {
	return &ReconciliationJobHandler{
		finderService:  finderService,
		creatorService: creatorService,
		logger:         log.New(os.Stdout, "[ReconciliationJobHandler] ", log.LstdFlags),
	}
}

// Register register reconciliation job handler to router
func (h *ReconciliationJobHandler) Register(router *httprouter.Router) {
	router.GET("/reconciliations", h.GetAllReconciliationJob)
	router.GET("/reconciliations/:id", h.GetReconciliationJobByID)
	router.POST("/reconciliations", h.CreateReconciliationJob)
}

// GetReconciliationJobByID get reconciliation job by id
func (h *ReconciliationJobHandler) GetReconciliationJobByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		logErrorF(h.logger, "invalid id: %s", p.ByName("id"))
		writeBadRequest(w, "invalid id")
		return
	}

	rj, err := h.finderService.FindByID(r.Context(), id)
	if err != nil {
		logErrorF(h.logger, "error find by id: %v", err)
		writeInternalServerError(w)
		return
	}
	if rj == nil {
		logInfoF(h.logger, "reconciliation job %d not found", id)
		writeNotFound(w, "reconciliation job not found")
		return
	}

	writeJSON(w, http.StatusOK, rj, nil)
}

// GetAllReconciliationJob get all reconciliation job
func (h *ReconciliationJobHandler) GetAllReconciliationJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pagination := getPagination(r)

	total, err := h.finderService.Count(r.Context())
	if err != nil {
		h.logger.Printf("error count: %v", err)
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
		h.logger.Printf("error find all: %v", err)
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusOK, rjs, pagination)
}

// CreateReconciliationJob create new reconciliation job
func (h *ReconciliationJobHandler) CreateReconciliationJob(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params, err := h.parseCreateReconciliationJobParams(r)
	if err != nil {
		logErrorF(h.logger, "error parse create reconciliation job params: %v", err)
		writeBadRequest(w, err.Error())
		return
	}

	rj, err := h.creatorService.Create(r.Context(), params)
	if err != nil {
		logErrorF(h.logger, "error create reconciliation job: %v", err)
		writeInternalServerError(w)
		return
	}

	writeJSON(w, http.StatusCreated, rj, nil)
}

func (h *ReconciliationJobHandler) parseCreateReconciliationJobParams(r *http.Request) (*reconciliatonjob.CreateParams, error) {
	if err := h.validateBodyRequest(r); err != nil {
		return nil, err
	}

	startDate := parseDate(r.FormValue("start_date"))
	endDate := parseDate(r.FormValue("end_date"))
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}

	discrepancyThreshold := parseFloat32(r.FormValue("discrepancy_threshold"))
	if discrepancyThreshold < 0 {
		discrepancyThreshold = 0
	}

	return &reconciliatonjob.CreateParams{
		StartDate:            startDate,
		EndDate:              endDate,
		DiscrepancyThreshold: discrepancyThreshold,
	}, nil
}

func (h *ReconciliationJobHandler) validateBodyRequest(r *http.Request) error {
	if err := r.ParseMultipartForm(entity.LimitContentSize); err != nil {
		return err
	}

	if r.ContentLength > entity.LimitContentSize {
		return errors.New("content size more than 100mb")
	}

	systemTrxFile := r.MultipartForm.File["system_transaction_file"]
	if err := h.validateSystemTrxFile(systemTrxFile); err != nil {
		return err
	}

	bankNames := r.MultipartForm.Value["bank_names"]
	bankTrxFiles := r.MultipartForm.File["bank_transaction_files"]
	if len(bankTrxFiles) == 0 {
		return errors.New("bank transaction files is required, at least provide one")
	}
	if len(bankNames) != len(bankTrxFiles) {
		return errors.New("bank names and bank transaction files length must be same")
	}

	for _, file := range bankTrxFiles {
		h.validateCSVFile(file)
	}

	return nil
}

func (h *ReconciliationJobHandler) validateSystemTrxFile(mp []*multipart.FileHeader) error {
	if len(mp) == 0 {
		return errors.New("system transaction file is required")
	}

	return h.validateCSVFile(mp[0])
}

func (h *ReconciliationJobHandler) validateCSVFile(file *multipart.FileHeader) error {
	// we can improve this by using mime type, but for the sake of simplicity we just check the extension
	if !isCSVExtension(file.Filename) {
		return fmt.Errorf("file %s must have a .csv extension", file.Filename)
	}
	if file.Size > entity.LimitCSVSize {
		return fmt.Errorf("file size %s more than 10mb", file.Filename)
	}

	return nil
}
