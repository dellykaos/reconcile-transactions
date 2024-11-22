package http

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/delly/amartha/entity"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	"github.com/julienschmidt/httprouter"
)

// ReconciliationJobHandler is a handler for reconciliation job
type ReconciliationJobHandler struct {
	finderService reconciliatonjob.Finder
	logger        *log.Logger
}

// NewReconciliationJobHandler create new reconciliation job handler, it used to create new reconciliation job, get reconciliation job by id, and get all reconciliation job
func NewReconciliationJobHandler(finderService reconciliatonjob.Finder) *ReconciliationJobHandler {
	return &ReconciliationJobHandler{
		finderService: finderService,
		logger:        log.New(os.Stdout, "[ReconciliationJobHandler] ", log.LstdFlags),
	}
}

// Register register reconciliation job handler to router
func (h *ReconciliationJobHandler) Register(router *httprouter.Router) {
	router.GET("/reconciliations", h.GetAllReconciliationJob)
	router.GET("/reconciliations/:id", h.GetReconciliationJobByID)
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
