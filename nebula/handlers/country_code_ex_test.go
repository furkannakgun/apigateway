package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSubExtractor_CCHeadersMiddleware(t *testing.T) {
	type args struct {
		url  string
		body string
	}
	tests := []struct {
		name        string
		args        args
		headerValue string
	}{
		// test with different URLs and bodies (strings converted to []byte) here
		{
			name: "With no MSISDN in the path and no body",
			args: args{
				url:  "/go-test-service/test/123456789",
				body: "",
			},
			headerValue: "",
		},
		{
			name: "With MSISDN in the path and no body",
			args: args{
				url:  "/go-test-service/test/dynamicbandwidth/34600556800",
				body: "",
			},
			headerValue: "ES",
		},
		{
			name: "With MSISDN in the path and in the body",
			args: args{
				url: "/go-test-service/test/dynamicbandwidth/34600556800",
				body: `
				{
					"duration": "PT20M",
					"maxBandwidth": 10000,
					"name": "Promotional offer number 5",
					"msisdn": "442012345678"
				}
				`,
			},
			headerValue: "ES",
		},
		{
			name: "With MSISDN in the body only",
			args: args{
				url: "/go-test-service/test/123456789",
				body: `
				{
					"duration": "PT20M",
					"maxBandwidth": 10000,
					"name": "Promotional offer number 5",
					"msisdn": "442012345678"
				}
				`,
			},
			headerValue: "GB",
		},
	}
	for _, tt := range tests {
		r := mux.NewRouter()
		doc, _ := openapi3.NewLoader().LoadFromFile("../testdata/openAPIspec_test.yaml")
		specValidatorRouter, _ := gorillamux.NewRouter(doc)
		r.Use(CCHeadersMiddleware(specValidatorRouter))

		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		r.Handle(tt.args.url, mockHandler)

		body := strings.NewReader(tt.args.body)
		req, _ := http.NewRequest("GET", tt.args.url, body)
		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)
		header := recorder.Header().Get("Country-Code")
		assert.Equal(t, tt.headerValue, header)
	}
}
