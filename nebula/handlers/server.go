package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

// Wrap http.Error so that calls to it can be monitored during testing
var httpErr = http.Error

// Serve requests
func ServeRequest(respWriter http.ResponseWriter, req *http.Request) {
	logger := structlog.GetLogger("nebula")
	target := config.Conf.GetString("targetURL") + req.URL.Path
	logger.LogRequestOut("Making request to the microservice...", target, &req.Header)
	if targetURL, err := url.Parse(target); err != nil {
		httpErr(respWriter, "Internal Server Error", http.StatusInternalServerError)
	} else {
		Proxy(targetURL).ServeHTTP(respWriter, req)
	}
}

// Health endpoint: return 200 OK to GET
func Health(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(resp, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Proxy requests to the microservice / API endpoint
func Proxy(target *url.URL) *httputil.ReverseProxy {
	logger := structlog.GetLogger("nebula")
	p := httputil.NewSingleHostReverseProxy(target)
	p.Director = func(req *http.Request) {
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}

	p.ModifyResponse = func(resp *http.Response) error {
		logger.LogResponseIn("Recieved response from the microservice", resp.StatusCode, &resp.Header)
		return nil
	}
	return p
}
