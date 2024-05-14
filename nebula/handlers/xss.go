package handlers

import (
	"bytes"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
	"io"
	"net/http"
	"regexp"
)

func XSSProtection(next http.Handler) http.Handler {
	logger := structlog.GetLogger("nebula")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xssPatterns := [...]string{
			`<\s*script\b[^>]*>;[^<]+<;\s*/\s*script\s*>`,
			`[\s]*((delete)|(exec)|(drop\s*table)|(insert)|(shutdown)|(update)|(\bor\b))`,
			`alert|script`,
		}

		body, err := io.ReadAll(r.Body)
		r.Body.Close() //  must close
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		if err != nil {
			logger.Errorv(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		targets := [...]string{r.URL.RequestURI(), r.URL.Path, string(body)}

		for _, target := range targets {
			for _, filter := range xssPatterns {
				regex := regexp.MustCompile(filter)
				matches := regex.FindStringSubmatch(target)
				if len(matches) != 0 {
					http.Error(w, "Invalid Request", http.StatusBadRequest)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
