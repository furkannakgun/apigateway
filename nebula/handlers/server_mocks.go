package handlers

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// Create mock for HTTP object
type httpMock struct {
	mock.Mock
}

// Create mock for http.Error function that will do nothing
func (m *httpMock) Error(w http.ResponseWriter, error string, code int) {
	m.Called()
}
