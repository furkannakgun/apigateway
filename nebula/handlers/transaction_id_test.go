package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTransactionId(t *testing.T) {
	tests := []struct {
		name       string
		headerKeys []string
	}{
		{
			name:       "Create from existing vf-trace-transaction-id",
			headerKeys: []string{"vf-trace-transaction-id"},
		},
		{
			name:       "Create from existing x-vf-trace-transaction-id",
			headerKeys: []string{"x-vf-trace-transaction-id"},
		},
		{
			name:       "Create from existing VF_INT_TRACE_ID",
			headerKeys: []string{"vf-int-trace-id"},
		},
		{
			name:       "Create from multiple transaction ID headers",
			headerKeys: []string{"x-vf-trace-transaction-id", "Vf-Int-Trace-Id"},
		},
		{
			name:       "Create from an Invalid header",
			headerKeys: []string{"Invalid-Header"},
		},
		{
			name:       "Create from an empty header",
			headerKeys: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)

			for i := 0; i < len(tt.headerKeys); i++ {
				request.Header.Add(tt.headerKeys[i], "example-transaction-id")
			}

			recorder := httptest.NewRecorder()

			handlerFunction := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			GetTransactionId(handlerFunction).ServeHTTP(
				recorder,
				request,
			)
			// Always check a value (UUID or not) has been created for the vf-trace-transaction-id header
			assert.NotEmpty(t, request.Header["vf-trace-transaction-id"])
		})
	}
}
