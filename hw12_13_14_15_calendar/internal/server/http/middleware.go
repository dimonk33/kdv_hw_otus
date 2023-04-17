package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type Middleware struct {
	logger Logger
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.logger.Info(fmt.Sprintf(
			"%s %s %s %s %d %s %s",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			r.Proto,
			r.Response.StatusCode,
			time.Since(start),
			r.UserAgent(),
		))
	})
}
