package handlers

import (
	"net/http"
	"strings"

	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

var IPWhitelist map[string]bool

func LoadIPs(IPs map[string]bool) {
	IPWhitelist = IPs
}

func IPFilterMiddleware(next http.Handler) http.Handler {
	logger := structlog.GetLogger("nebula")
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			IP, _, _ := strings.Cut(r.RemoteAddr, ":")
			if IPWhitelist[IP] {
				next.ServeHTTP(w, r)
			} else {
				logger.Warnf("IP %s not allowed!", IP)
				w.WriteHeader(http.StatusUnauthorized)
			}
		},
	)
}
