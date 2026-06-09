using Prometheus;

namespace MetricsApi.Metrics;

/// <summary>
/// Declares the HTTP counter and histogram on the prometheus-net default registry.
/// Each metric type answers a different question, so it needs its own primitive.
/// </summary>
public static class AppMetrics
{
    // Counter only increases; used to compute rate (events per second).
    public static readonly Counter HttpRequestsTotal = Prometheus.Metrics.CreateCounter(
        "http_requests_total",
        "Total HTTP requests",
        new CounterConfiguration
        {
            LabelNames = new[] { "method", "route", "status_code" },
        });

    // Histogram distributes latency into fixed cumulative buckets for quantiles.
    public static readonly Histogram HttpRequestDurationSeconds = Prometheus.Metrics.CreateHistogram(
        "http_request_duration_seconds",
        "HTTP request latency in seconds",
        new HistogramConfiguration
        {
            LabelNames = new[] { "method", "route" },
            Buckets = new[] { 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5 },
        });
}
