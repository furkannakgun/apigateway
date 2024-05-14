package router

import (
	"net/http"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gorilla/mux"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/handlers"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/limiter"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/specvalidator"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

func LoadSpec(path string) routers.Router {
	logger := structlog.GetLogger("nebula")
	doc, err := openapi3.NewLoader().LoadFromFile(path)
	if err != nil {
		logger.Errorf("Error while loading OpenAPI specification: %s", err)
		panic(err)
	}

	specRouter, err := gorillamux.NewRouter(doc)
	if err != nil {
		logger.Errorf("Invalid OpenAPI specification: %s", err)
		panic(err)
	}
	return specRouter
}

func Configure(r *mux.Router) {
	logger := structlog.GetLogger("nebula")
	// Get sub-router for the /health endpoint
	healthRouter := r.Path("/health").Subrouter()
	healthRouter.PathPrefix("").HandlerFunc(handlers.Health)
	// Get sub-router for microservice endpoints
	microserviceRouter := r.PathPrefix("/").Subrouter()
	microserviceSpec := LoadSpec(config.OpenAPIPath)
	// Logging middleware
	microserviceRouter.Use(logger.LoggingHandler)
	// IP whitelisting
	if config.Conf.Get("IPWhitelist") != nil {
		IPs := make(map[string]bool)
		config.Conf.UnmarshalKey("IPWhitelist", &IPs)
		// If 0.0.0.0 is set, we allow all traffic
		if !IPs["0.0.0.0"] {
			handlers.LoadIPs(IPs)
			microserviceRouter.Use(handlers.IPFilterMiddleware)
		}
	}
	// Spike Arrest
	err := limiter.CreateSpikeArrest(microserviceSpec, microserviceRouter)
	if err != nil {
		logger.Errorf("Error initializing Spike Arrest: %s", err)
		panic(err)
	}
	// Header Validation
	microserviceRouter.Use(handlers.RequestHeaderValidation)
	// Security Headers
	microserviceRouter.Use(handlers.SecHeadersMiddleware)
	// XSS Protection
	microserviceRouter.Use(handlers.XSSProtection)
	// Transaction ID generation
	microserviceRouter.Use(handlers.GetTransactionId)
	// Validate request against microservice spec
	microserviceRouter.Use(specvalidator.ValidateRequest(microserviceSpec))
	// Country Code Extraction
	countryCodeExtraction := config.Conf.GetBool("countryCodeEx")
	if countryCodeExtraction {
		microserviceRouter.Use(handlers.CCHeadersMiddleware(microserviceSpec))
	}
	err = limiter.CreateAndLoadLimiters(microserviceSpec, microserviceRouter)
	if err != nil {
		logger.Errorf("Error initializing limiter middleware: %s", err)
		panic(err)
	}
	microserviceRouter.PathPrefix("/").HandlerFunc(handlers.ServeRequest)
}

type RouteSwapper struct {
	mu     sync.Mutex
	Router *mux.Router
}

func (rs *RouteSwapper) Init() {
	rs.Router = mux.NewRouter()
}

func (rs *RouteSwapper) Swap(newRouter *mux.Router) {
	rs.mu.Lock()
	rs.Router = newRouter
	rs.mu.Unlock()
}

func (rs *RouteSwapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rs.mu.Lock()
	router := rs.Router
	rs.mu.Unlock()
	router.ServeHTTP(w, r)
}
