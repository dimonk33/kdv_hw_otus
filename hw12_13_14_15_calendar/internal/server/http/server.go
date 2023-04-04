package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	log    Logger
	server *http.Server
}

type Logger interface {
	Info(msg string)
}

type Application interface {
	HelloHandler(w http.ResponseWriter, r *http.Request)
}

func NewServer(addr string, logger Logger, app Application) *Server {
	s := &Server{
		log: logger,
	}

	router := http.NewServeMux()
	router.HandleFunc("/hello", app.HelloHandler)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.LogHandler(router),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}

func (s *Server) LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		s.log.Info(fmt.Sprintf(
			"%s %s %s %s %d %s %s",
			req.RemoteAddr,
			req.Method,
			req.RequestURI,
			req.Proto,
			req.Response.StatusCode,
			time.Since(start),
			req.UserAgent(),
		))
	})
}
