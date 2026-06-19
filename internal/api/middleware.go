package api

import (
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// responseRecorder capture le status HTTP et le batch_id pour le logging.
type responseRecorder struct {
	http.ResponseWriter
	status  int
	batchID string
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Chain enchaîne plusieurs middlewares (pattern vu en TP Gin).
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// loggingMiddleware journalise chaque requête en JSON via slog (sauf /healthz).
func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" {
				// fmt.Println("healthz ignoré par slog") // curl /healthz ne doit pas loguer
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
			}
			if rec.batchID != "" {
				attrs = append(attrs, "batch_id", rec.batchID)
			} else if id := extractCheckID(r.URL.Path); id != "" {
				attrs = append(attrs, "batch_id", id)
			}

			// fmt.Println("log requête", r.Method, r.URL.Path, rec.status)
			// fmt.Println("duration_ms :", time.Since(start).Milliseconds(), "batch_id :", rec.batchID)
			logger.Info("request", attrs...)
		})
	}
}

// recoveryMiddleware intercepte les panic et renvoie une 500 propre (bonus).
func recoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// test recovery bonus
					// fmt.Println("PANIC interceptée :", rec, "sur", r.URL.Path)
					logger.Error("panic",
						"err", rec,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)
					writeAPIError(w, http.StatusInternalServerError, "internal", "erreur interne")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// NewJSONLogger configure slog avec un handler JSON et le niveau LOG_LEVEL.
func NewJSONLogger(level string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(level),
	}))
}

// parseLogLevel lit la variable d'environnement LOG_LEVEL.
func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		// fmt.Println("LOG_LEVEL inconnu, fallback INFO :", level)
		return slog.LevelInfo
	}
}

func extractCheckID(path string) string {
	return strings.TrimPrefix(path, "/v1/checks/")
}
