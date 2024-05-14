package handlers

import (
	"net/http"

	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

func SecHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			h := map[string]string{
				"X-Content-Type-Options":  "nosniff",
				"X-Frame-Options":         "deny",
				"Content-Security-Policy": "default-src 'none'",
			}
			for k, v := range h {
				w.Header().Set(k, v)
			}
			next.ServeHTTP(w, r)
		},
	)
}

func RequestHeaderValidation(next http.Handler) http.Handler {
	logger := structlog.GetLogger("nebula")
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			acceptHeader := r.Header.Get("Accept")
			contentType := r.Header.Get("Content-Type")
			reqMethod := r.Method
			if acceptHeader != "application/json" && acceptHeader != "*/*" {
				logger.Error("Invalid/Missing Accept Header for Request")
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if reqMethod == "PUT" || reqMethod == "POST" || reqMethod == "PATCH" {
				if contentType != "application/json" {
					logger.Error("Invalid Content-Type for Request")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			next.ServeHTTP(w, r)
		},
	)
}
