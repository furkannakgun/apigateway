package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXSSProtection(t *testing.T) {
	type args struct {
		next http.Handler
		url  string
		body io.Reader
	}
	tests := []struct {
		name    string
		args    args
		expCode int
	}{
		{
			name: "Valid request",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				url:  "http://testing",
				body: nil,
			},
			expCode: 200,
		},
		{
			name: "Invalid URI",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				url:  "http://testing/script",
				body: nil,
			},
			expCode: 400,
		},
		{
			name: "Invalid Body",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				url:  "http://testing/",
				body: strings.NewReader("<script>delete</script>"),
			},
			expCode: 400,
		},
	}
	rr := httptest.NewRecorder()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			XSSProtection(tt.args.next).ServeHTTP(rr, httptest.NewRequest("GET", tt.args.url, tt.args.body))
			assert.Equal(t, tt.expCode, rr.Code)
		})
	}
}
