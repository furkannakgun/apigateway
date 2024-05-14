package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/router"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

var Router router.RouteSwapper
var server *http.Server

func Init() {
	// Initialize Router
	Router.Init()
	// Configure Router
	router.Configure(Router.Router)
	// Configure server
	server = &http.Server{
		Addr:           ":8080",
		Handler:        http.TimeoutHandler(corsMiddleware(&Router), 10*time.Second, "Request timed out"),
		ReadTimeout:    9 * time.Second,
		WriteTimeout:   8 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func Run() {
	logger := structlog.GetLogger("nebula")
	logger.Info("Starting config watcher go routines")
	logger.Info("Serving requests...")
	logger.Fatalv(server.ListenAndServe())
}

func ReloadRouter() {
	config.Load()
	r := mux.NewRouter()
	router.Configure(r)
	Router.Swap(r)
}

func corsMiddleware(next http.Handler) http.Handler {
	// CORS management
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
