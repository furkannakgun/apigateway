package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeadersMiddleware(t *testing.T) {
	type args struct {
		header http.Header
		next   http.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add security headers",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				header: http.Header{
					"Content-Security-Policy": []string{"default-src 'none'"},
					"X-Content-Type-Options":  []string{"nosniff"},
					"X-Frame-Options":         []string{"deny"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://testing", nil)
			recorder := httptest.NewRecorder()

			SecHeadersMiddleware(tt.args.next).ServeHTTP(recorder, request)

			assert.Equal(t, tt.args.header, recorder.Header())
		})
	}
}

func TestRequestHeaderValidation(t *testing.T) {
	type args struct {
		header http.Header
		next   http.Handler
	}
	tests := []struct {
		name      string
		args      args
		reqMethod string
		expCode   int
	}{
		{
			name: "Invalid Accept Header",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				header: http.Header{
					"Accept": {"invalid/type"},
				},
			},
			expCode: 400,
		},
		{
			name: "Missing Accept Header",
			args: args{
				next:   http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				header: http.Header{
					// No Accept Header to be seen here!
				},
			},
			expCode: 400,
		},
		{
			name: "Valid Accept Header + Method and Invalid Content-Type",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				header: http.Header{
					"Accept":       {"*/*"},
					"Content-Type": {"invalid_content_type"},
				},
			},
			reqMethod: "PUT",
			expCode:   400,
		},
		{
			name: "Valid Accept Header + Method and valid Content-Type",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				header: http.Header{
					"Accept":       {"application/json"},
					"Content-Type": {"application/json"},
				},
			},
			reqMethod: "PUT",
			expCode:   200,
		},
		{
			name: "Valid Accept Header + GET Method",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				header: http.Header{
					"Accept":       {"application/json"},
					"Content-Type": {"i_dont_matter"}, // Function won't check for GET method so shouldn't read this header
				},
			},
			reqMethod: "GET",
			expCode:   200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.reqMethod, "http://testing", nil)
			recorder := httptest.NewRecorder()
			request.Header = tt.args.header

			RequestHeaderValidation(tt.args.next).ServeHTTP(recorder, request)
			assert.Equal(t, tt.expCode, recorder.Code)
		})
	}
}
