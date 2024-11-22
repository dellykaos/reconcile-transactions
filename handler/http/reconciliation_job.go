package http

import (
	"net/http"
	"strconv"

	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	"github.com/julienschmidt/httprouter"
)

// ReconciliationJobHandler is a handler for reconciliation job
type ReconciliationJobHandler struct {
	finderService reconciliatonjob.FinderService
}

// NewReconciliationJobHandler create new reconciliation job handler, it used to create new reconciliation job, get reconciliation job by id, and get all reconciliation job
func NewReconciliationJobHandler(finderService reconciliatonjob.FinderService) *ReconciliationJobHandler {
	return &ReconciliationJobHandler{finderService: finderService}
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
		writeBadRequest(w, "invalid id")
		return
	}

	rj, err := h.finderService.FindByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if rj == nil {
		writeNotFound(w, "reconciliation job not found")
		return
	}

	writeJSON(w, http.StatusOK, rj, nil)
}

// GetAllReconciliationJob get all reconciliation job
func (h *ReconciliationJobHandler) GetAllReconciliationJob(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pagination := getPagination(p)

	total, err := h.finderService.Count(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pagination.Total = int32(total)
	rjs, err := h.finderService.FindAll(r.Context(), pagination.Limit, pagination.Offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, rjs, pagination)
}
