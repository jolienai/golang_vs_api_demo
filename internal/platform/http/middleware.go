package http

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey).(string)
	return requestID
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(data)
	r.bytes += n
	return n, err
}

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(recorder, r)

			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"status", recorder.status,
				"response_bytes", recorder.bytes,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", RequestIDFromContext(r.Context()),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			}

			switch {
			case recorder.status >= http.StatusInternalServerError:
				logger.ErrorContext(r.Context(), "http request", attrs...)
			case recorder.status >= http.StatusBadRequest:
				logger.WarnContext(r.Context(), "http request", attrs...)
			default:
				logger.InfoContext(r.Context(), "http request", attrs...)
			}
		})
	}
}

func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if value := recover(); value != nil {
					logger.ErrorContext(r.Context(), "panic recovered",
						"panic", value,
						"stack", string(debug.Stack()),
						"request_id", RequestIDFromContext(r.Context()),
					)
					InternalServerError(w)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
