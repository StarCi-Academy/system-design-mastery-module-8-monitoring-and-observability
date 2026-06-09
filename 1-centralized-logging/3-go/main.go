// Package main — a small net/http logging-api that emits structured JSON logs via zerolog
// to stdout (shipped to Loki by a Promtail sidecar) and reports unhandled panics to Sentry.
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

// logger writes structured JSON to stdout; Promtail tails it and ships to Loki.
var logger zerolog.Logger

func main() {
	// Configure zerolog to emit JSON with a level field and timestamp.
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", envOr("LOG_APP_NAME", "logging-api")).
		Str("env", envOr("LOG_ENV", "lab")).
		Logger()

	// Initialize Sentry only when SENTRY_DSN is set (skip if empty — not mandatory for the lab).
	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		_ = sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			TracesSampleRate: 1.0,
		})
		defer sentry.Flush(2 * time.Second)
	}

	// sentryhttp middleware recovers panics and reports them to Sentry.
	sentryHandler := sentryhttp.New(sentryhttp.Options{Repanic: false})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/orders", acceptHandler)
	mux.HandleFunc("/api/orders/error", errorHandler)

	port := envOr("PORT", "3000")
	logger.Info().Msgf("logging-api listening on :%s", port)
	if err := http.ListenAndServe(":"+port, sentryHandler.Handle(mux)); err != nil {
		logger.Fatal().Err(err).Msg("server failed")
	}
}

// acceptHandler — GET /api/orders; logs INFO with a trace then returns the accepted envelope.
func acceptHandler(w http.ResponseWriter, _ *http.Request) {
	trace := "trace-" + time.Now().Format("20060102150405.000")
	host := envOr("HOSTNAME", "local")
	logger.Info().Str("trace", trace).Msgf("[%s] Order accepted trace=%s", host, trace)
	writeJSON(w, http.StatusOK, map[string]any{"status": "accepted", "trace": trace})
}

// errorHandler — GET /api/orders/error; logs ERROR then panics so Sentry captures it,
// while the response is normalized to the cross-language {statusCode, message} envelope.
func errorHandler(w http.ResponseWriter, _ *http.Request) {
	host := envOr("HOSTNAME", "local")
	logger.Error().Msgf("[%s] Simulated order failure for Loki + Sentry demo", host)
	// Report to Sentry explicitly, then return the normalized error envelope.
	sentry.CaptureException(&simulatedError{})
	writeJSON(w, http.StatusInternalServerError, map[string]any{
		"statusCode": 500,
		"message":    "Internal server error",
	})
}

// simulatedError — a real error type so Sentry groups the demo failure as an issue.
type simulatedError struct{}

func (e *simulatedError) Error() string { return "Simulated order processing failure" }

func writeJSON(w http.ResponseWriter, status int, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
