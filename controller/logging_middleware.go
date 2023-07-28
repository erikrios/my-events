package controller

import (
	"log"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rec := statusRecorder{w, 500}

			next.ServeHTTP(&rec, r)

			log.Printf("method=%s, uri=%s, status=%d\n", r.Method, r.RequestURI, rec.status)
		},
	)
}
