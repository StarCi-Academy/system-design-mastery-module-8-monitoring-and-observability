// Distributed tracing lab (Go + OpenTelemetry-Go SDK + net/http).
// Initializes the tracer provider before the server serves, wraps the
// /checkout route with otelhttp for auto-instrumentation, and emits a child
// span checkout.inventory_reserve under the inbound HTTP span.
package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// initTracer wires the OTLP HTTP exporter and the global tracer provider.
// It must run before the HTTP server accepts requests so otelhttp and
// tracer.Start resolve a real provider instead of the no-op fallback.
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:4318/v1/traces"
	}
	// OTLP HTTP exporter: where spans are shipped to.
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
	if err != nil {
		return nil, err
	}
	res, _ := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName("go-tracing-app"),
	))
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // batches spans, exports asynchronously
		sdktrace.WithResource(res),
	)
	// Register globally so otelhttp and tracer.Start pick it up.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}

// checkoutHandler emits a child span under the inbound HTTP span and returns
// the active traceId so the caller can search it in Jaeger.
func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	// r.Context() already carries the inbound HTTP span installed by otelhttp.
	ctx := r.Context()
	tracer := otel.Tracer("go-tracing-app-checkout")
	_, span := tracer.Start(ctx, "checkout.inventory_reserve")
	delayMs := 50 + rand.Intn(400)
	span.SetAttributes(attribute.Int("checkout.inventory.delay_ms", delayMs))
	time.Sleep(time.Duration(delayMs) * time.Millisecond)
	span.End()

	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Checkout completed",
		"traceId": traceID,
	})
}

// healthHandler is a liveness probe used by the lab to confirm the server is up.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// newRouter wires OTel once at startup; handlers stay free of SDK setup.
func newRouter() http.Handler {
	mux := http.NewServeMux()
	// Only the checkout feature is registered so a single /checkout flow
	// demonstrates a clean trace + child span without noise.
	mux.Handle("/checkout", otelhttp.NewHandler(
		http.HandlerFunc(checkoutHandler), "GET /checkout"))
	mux.HandleFunc("/health", healthHandler)
	return mux
}

func main() {
	ctx := context.Background()
	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(ctx) // flush remaining batch on exit

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("go-tracing-app listening on :%s", port)
	if err := http.ListenAndServe(":"+port, newRouter()); err != nil {
		log.Fatal(err)
	}
}
