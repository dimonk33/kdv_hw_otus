package internalgrpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type Middleware struct {
	Logger Logger
}

func (m *Middleware) unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	r, ok := req.(http.Request)
	if !ok {
		m.Logger.Info(fmt.Sprintf("%s - %d", info.FullMethod, time.Since(start)))
	} else {
		m.Logger.Info(fmt.Sprintf(
			"%s %s %s %s %d %s %s",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			r.Proto,
			r.Response.StatusCode,
			time.Since(start),
			r.UserAgent(),
		))
	}

	return h, err
}
