package handlers

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
)

func TestServeRequest(t *testing.T) {
	tests := []struct {
		name        string
		mockUrl     string
		url         string
		expErrsHTTP int
	}{
		{
			name:        "Serve request parse URL OK",
			mockUrl:     "http://example.com",
			url:         "/go-test-service/test/123456789",
			expErrsHTTP: 0,
		},
		{
			name:        "Serve request parse URL Error",
			mockUrl:     "ht*tp://example.com", //invalid scheme
			url:         "/go-test-service/test/123456789",
			expErrsHTTP: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config.Conf = viper.New()
			config.Conf.Set("targetURL", tt.mockUrl)

			mockHttp := httpMock{}
			httpErr = mockHttp.Error

			mockHttp.On("Error").Return(nil)

			// Act
			request := httptest.NewRequest("GET", tt.url, nil)
			recorder := httptest.NewRecorder()
			ServeRequest(recorder, request)

			// Assert
			/*
				Test based on number of errors produced, if no errors are called
				it can be assumed that ServeRequest has passed through everything OK
				to Proxy().
			*/
			mockHttp.AssertNumberOfCalls(t, "Error", tt.expErrsHTTP)
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		requestMethod string
		expCode       int
	}{
		{
			name:          "Test Health Endpoint with GET method",
			requestMethod: "GET",
			expCode:       200,
		},
		{
			name:          "Test Health Endpoint with POST method",
			requestMethod: "POST",
			expCode:       405,
		},
		{
			name:          "Test Health Endpoint with PUT method",
			requestMethod: "PUT",
			expCode:       405,
		},
		{
			name:          "Test Health Endpoint with DELETE method",
			requestMethod: "DELETE",
			expCode:       405,
		},
	}
	for _, tt := range tests {
		request := httptest.NewRequest(tt.requestMethod, "/health", nil)
		recorder := httptest.NewRecorder()
		Health(recorder, request)
		assert.Equal(t, tt.expCode, recorder.Code)
	}
}

func TestProxy(t *testing.T) {
	tests := []struct {
		name    string
		mockURL string
	}{
		{
			name:    "With valid URL",
			mockURL: "http://127.0.0.1:8080/go-test-service/test/123456789",
		},
	}
	for _, tt := range tests {
		// url.Parse validates URL for us, so invalid URLs would not pass here before Proxy() is called
		target, _ := url.Parse(tt.mockURL)
		reverseProxy := Proxy(target)
		assert.NotEmpty(t, reverseProxy.Director)
		assert.NotEmpty(t, reverseProxy.ModifyResponse)
	}
}
