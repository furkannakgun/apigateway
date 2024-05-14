package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

func GetTransactionId(next http.Handler) http.Handler {
	logger := structlog.GetLogger("nebula")
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			traceIdHeaders := []string{"Vf-Trace-Transaction-Id", "X-Vf-Trace-Transaction-Id", "Vf-Int-Trace-Id"}

			//	Create a UUID and assign it to the vf-trace-transaction-id header

			generatedUuid, err := uuid.NewRandom()
			if err != nil {
				logger.Errorf("Error generating UUID for request %v", err)
				return
			}
			r.Header["vf-trace-transaction-id"] = []string{generatedUuid.String()}

			/*
				Loop through headers to be looked for, if a value is found,
				overwrite the UUID in the vf-trace-transaction-id header with
				the aformentioned value and the original header is deleted.
			*/

			for i := 0; i < len(traceIdHeaders); i++ {
				headerValue := r.Header.Get(traceIdHeaders[i])
				if headerValue != "" {
					r.Header["vf-trace-transaction-id"] = []string{headerValue}
					r.Header.Del(traceIdHeaders[i])
				}
			}
			next.ServeHTTP(w, r)
		},
	)
}
