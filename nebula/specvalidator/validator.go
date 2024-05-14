package specvalidator

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/gorilla/mux"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

func ValidateRequest(validator routers.Router) mux.MiddlewareFunc {
	logger := structlog.GetLogger("nebula")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route, pathParams, err := validator.FindRoute(r)
			if err != nil {
				logger.Errorf("Route not found while validating against OpenAPI spec: %v", err)
				http.Error(w, "Route not found", http.StatusNotFound)
				return
			}
			requestValidationInput := &openapi3filter.RequestValidationInput{
				Request:    r,
				PathParams: pathParams,
				Route:      route,
			}
			validationError := openapi3filter.ValidateRequest(r.Context(), requestValidationInput)
			if validationError != nil {
				logger.Errorf("Error while validating request against spec: %v", validationError)
				http.Error(w, "Error while validating request against spec", http.StatusBadRequest)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
