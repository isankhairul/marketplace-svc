package logger

import (
	"context"
	"fmt"
	stdhttp "net/http"

	"github.com/go-kit/kit/transport/http"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	TraceIDContextKey       = "KD-Trace-ID"
	TraceIDRequestHeaderKey = "X-Correlation-ID"
)

// TraceIdentifier use in the ServerBefore method
func TraceIdentifier() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		traceId := r.Header.Get(TraceIDRequestHeaderKey)
		if traceId == "" {
			traceId, _ = gonanoid.New()
			r.Header.Set(TraceIDRequestHeaderKey, traceId)
		}
		return context.WithValue(ctx, TraceIDContextKey, traceId) // nolint
	}
}

// TraceIdentifierMiddleware use like the middleware
func TraceIdentifierMiddleware(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		traceId := r.Header.Get(TraceIDRequestHeaderKey)
		if traceId == "" {
			traceId, _ = gonanoid.New()
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, TraceIDContextKey, traceId) // nolint
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTraceIdentifier(ctx context.Context) string {
	traceId := ctx.Value(TraceIDContextKey)
	if traceId == nil {
		traceId, _ = gonanoid.New()
		//return ""
	}
	return fmt.Sprint(traceId)
}
