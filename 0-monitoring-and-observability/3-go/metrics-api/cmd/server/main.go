// Command server wires the demo /cats resource and the /metrics endpoint, then
// serves them on :3000. Prometheus scrapes /metrics on its scrape_interval.
package main

import (
	"log"
	"net/http"

	"metrics-api/internal/cats"
	"metrics-api/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	mux := http.NewServeMux()
	// promhttp serves the shared registry in Prometheus text exposition format.
	mux.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}))
	// The route template "/cats" is bound here, keeping the route label bounded.
	mux.Handle("/cats", metrics.Middleware(cats.Handler(), "/cats"))

	log.Println("metrics-api (go) listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
