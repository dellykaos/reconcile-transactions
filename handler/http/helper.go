package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	minLimit   = 10
	maxLimit   = 100
	dateFormat = "2006-01-02"
)

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

func writeInternalServerError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "internal server error")
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

func getPagination(r *http.Request) pagination {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = minLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	return pagination{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
}

func logInfoF(logger *log.Logger, format string, v ...interface{}) {
	logger.Printf("[INFO] "+format, v...)
}

func logErrorF(logger *log.Logger, format string, v ...interface{}) {
	logger.Printf("[ERROR] "+format, v...)
}

func parseDate(dateStr string) time.Time {
	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return time.Now()
	}

	return date
}

func parseFloat32(str string) float32 {
	f, _ := strconv.ParseFloat(str, 32)
	return float32(f)
}

func isCSVExtension(filename string) bool {
	return filename[len(filename)-4:] == ".csv"
}
