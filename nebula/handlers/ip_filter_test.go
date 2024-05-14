package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadIPs(t *testing.T) {
	type args struct {
		IPs map[string]bool
	}
	tests := []struct {
		name     string
		args     args
		wantList map[string]bool
	}{
		{
			name:     "IP whitelist Loaded",
			args:     args{IPs: map[string]bool{"127.0.0.1": true}},
			wantList: map[string]bool{"127.0.0.1": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadIPs(tt.args.IPs)
			assert.Equal(t, IPWhitelist, tt.wantList)
		})
	}
}

func TestIPList_IPFilterMiddleware(t *testing.T) {
	type fields struct {
		white    map[string]bool
		sourceIP string
	}
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCode int
	}{
		{
			name:     "IP whitelisted.",
			fields:   fields{white: map[string]bool{"127.0.0.1": true}, sourceIP: "127.0.0.1:62444"},
			args:     args{next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) })},
			wantCode: 201,
		},
		{
			name:     "IP not whitelisted.",
			fields:   fields{white: map[string]bool{"127.0.0.1": true}, sourceIP: "192.168.0.1:62444"},
			args:     args{next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) })},
			wantCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IPWhitelist = tt.fields.white

			request := httptest.NewRequest("GET", "/", nil)
			request.RemoteAddr = tt.fields.sourceIP
			recorder := httptest.NewRecorder()
			IPFilterMiddleware(tt.args.next).ServeHTTP(
				recorder,
				request,
			)
			assert.Equal(t, tt.wantCode, recorder.Code)
		})
	}
}
