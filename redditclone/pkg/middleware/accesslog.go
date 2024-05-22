package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.status = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

func AccessLogMiddleware(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
		}
		logger.Infow("access log middleware")
		start := time.Now()
		next.ServeHTTP(lrw, r)
		logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
			"status", lrw.status,
		)
	})
}
