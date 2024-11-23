package middleware

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// WithLogger is a middleware to log request
func WithLogger(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log := zap.L()
		log.Info(fmt.Sprintf("request to %s", r.URL.Path), zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("query", r.URL.RawQuery))
		fn(w, r, p)
	}
}
