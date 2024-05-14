package specvalidator

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateRequest(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		reqMethod string
		expCode   int
	}{
		{
			name:      "URL found in the OpenAPI Spec",
			url:       "/go-test-service/test/123456789",
			reqMethod: "GET",
			expCode:   200,
		},
		{
			name:      "Request Method not found for path",
			url:       "/go-test-service/test/123456789",
			reqMethod: "POST",
			expCode:   404,
		},
		{
			name:      "URL not found in the OpenAPI Spec",
			url:       "/go-test-service/fake/path/test/123456789",
			reqMethod: "GET",
			expCode:   404,
		},
		{
			name:      "Invalid path",
			url:       "/go test service/fake/path/test/123456789",
			reqMethod: "GET",
			expCode:   404,
		},
		{
			name:      "Invalid path params",
			url:       "/go-test-service/test/fake_test_params",
			reqMethod: "GET",
			expCode:   400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mux.NewRouter()
			doc, _ := openapi3.NewLoader().LoadFromFile("../testdata/openAPIspec_test.yaml")
			specValidatorRouter, _ := gorillamux.NewRouter(doc)
			r.Use(ValidateRequest(specValidatorRouter))

			r.Handle(tt.url, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

			req, _ := http.NewRequest(tt.reqMethod, tt.url, nil)
			recorder := httptest.NewRecorder()

			r.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expCode, recorder.Code)

		})
	}
}
