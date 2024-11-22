package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const minLimit = 10

type pagination struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
	Total  int32 `json:"total"`
}

type dataResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}, meta interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dataResponse{
		Data: data,
		Meta: meta,
	})
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeBadRequest(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func writeNotFound(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func getPagination(p httprouter.Params) pagination {
	limit, _ := strconv.Atoi(p.ByName("limit"))
	if limit <= 0 {
		limit = minLimit
	}

	offset, _ := strconv.Atoi(p.ByName("offset"))
	if offset < 0 {
		offset = 0
	}

	return pagination{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
}
