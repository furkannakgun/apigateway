package limiter

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPKeyFunc(t *testing.T) {
	type rHeaders struct {
		// Create some example headers that could pop up during actual use of Nebula
		// Scope to add more test headers if tests need to be expanded in future.
		x_forwarded_for string
	}
	tests := []struct {
		name           string
		inputHeaders   []string
		requestHeaders rHeaders
		requestAddr    string
		expResult      string
		expErr         error
	}{
		{
			name:           "Retrieve IPKeyFunc OK empty input headers",
			inputHeaders:   []string{},
			requestHeaders: rHeaders{},
			requestAddr:    "127.0.0.1:8080",
			expResult:      "127.0.0.1",
			expErr:         nil,
		},
		{
			name:           "Retrieve IPKeyFunc invalid remote address",
			inputHeaders:   []string{},
			requestHeaders: rHeaders{},
			requestAddr:    "127.0.0.1:8080:8080", // Invalid, double declaration of port
			expResult:      "",
			expErr:         &net.AddrError{Err: "too many colons in address", Addr: "127.0.0.1:8080:8080"},
		},
		{
			name:         "Retrieve IPKeyFunc OK with input headers",
			inputHeaders: []string{"X-Forwarded-For"},
			requestHeaders: rHeaders{
				x_forwarded_for: "1.0.0.1", // Example address
			},
			requestAddr: "127.0.0.1:8080",
			expResult:   "1.0.0.1",
			expErr:      nil,
		},
		{
			name:           "Retrieve IPKeyFunc OK input headers not found",
			inputHeaders:   []string{"X-Forwarded-For"},
			requestHeaders: rHeaders{},
			requestAddr:    "127.0.0.1:8080",
			expResult:      "127.0.0.1",
			expErr:         nil,
		},
		{
			name:         "Retrieve IPKeyFunc OK request headers ignored",
			inputHeaders: []string{},
			requestHeaders: rHeaders{
				x_forwarded_for: "1.0.0.1", // Example address
			},
			requestAddr: "127.0.0.1:8080",
			expResult:   "127.0.0.1",
			expErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.RemoteAddr = tt.requestAddr
			req.Header.Set("X-Forwarded-For", tt.requestHeaders.x_forwarded_for)

			fn := IPKeyFunc(tt.inputHeaders...)

			result, err := fn(req)

			assert.Equal(t, tt.expResult, result)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func TestNewMiddleware(t *testing.T) {
	type args struct {
		s mockStore // mock limiter.Store for sake of creating middleware
		f KeyFunc
	}
	tests := []struct {
		name   string
		args   args
		expMW  *Middleware
		expErr error
	}{
		{
			name: "Create New Middleware OK",
			args: args{
				s: &mockStoreOk{},
				f: mockKeyFuncOk(),
			},
			expMW:  &Middleware{},
			expErr: nil,
		},
		{
			name: "Create New Middleware missing store",
			args: args{
				s: nil,
				f: mockKeyFuncOk(),
			},
			expMW:  nil,
			expErr: errors.New("store cannot be nil"),
		},
		{
			name: "Create New Middleware missing KeyFunc",
			args: args{
				s: &mockStoreOk{},
				f: nil,
			},
			expMW:  nil,
			expErr: errors.New("key function cannot be nil"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw, err := NewMiddleware(tt.args.s, tt.args.f)

			assert.IsType(t, tt.expMW, mw)
			assert.Equal(t, tt.expErr, err)

		})
	}
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name    string
		keyFunc KeyFunc
		store   mockStore
		expCode int
	}{
		{
			name:    "Test handle function KeyFunc error",
			keyFunc: mockKeyFuncErr(),
			store:   &mockStoreOk{},
			expCode: 500,
		},
		{
			name:    "Test handle function empty key",
			keyFunc: mockKeyFuncEmptyStr(),
			store:   &mockStoreOk{},
			expCode: 200,
		},
		{
			name:    "Test handle function store take error",
			keyFunc: mockKeyFuncOk(),
			store:   &mockStoreTakeErr{},
			expCode: 500,
		},
		{
			name:    "Test handle function store no tokens remaining",
			keyFunc: mockKeyFuncOk(),
			store:   &mockStoreNoTokensLeft{},
			expCode: 429,
		},
		{
			name:    "Test handle function all OK.",
			keyFunc: mockKeyFuncOk(),
			store:   &mockStoreOk{},
			expCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := Middleware{store: tt.store, keyFunc: tt.keyFunc}
			request := httptest.NewRequest("GET", "http://example.com", nil)

			recorder := httptest.NewRecorder()
			handlerFunction := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			mw.Handle(handlerFunction).ServeHTTP(
				recorder,
				request,
			)
			assert.Equal(t, tt.expCode, recorder.Code)
		})
	}
}
