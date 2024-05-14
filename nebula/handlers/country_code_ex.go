package handlers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/getkin/kin-openapi/routers"
	"github.com/gorilla/mux"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/ccextractor"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/subextractor"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

func CCHeadersMiddleware(validator routers.Router) mux.MiddlewareFunc {
	logger := structlog.GetLogger("nebula")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract path params and body
			_, pathParams, _ := validator.FindRoute(r)
			body, _ := io.ReadAll(r.Body)
			r.Body.Close() //  must close
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			subscriber := subextractor.GetSubscriberFromRequest(pathParams, body)
			countryCode := ccextractor.GetCountryCodeAsAlpha2(subscriber)
			if countryCode != "" { // Only if Nebula finds an MSISDN or IMSI in the path/body, otherwise do not create a new header
				w.Header().Add("Country-Code", countryCode)
			} else {
				logger.Error("Could not find Country Code from request")
			}
			next.ServeHTTP(w, r)
		})
	}
}
