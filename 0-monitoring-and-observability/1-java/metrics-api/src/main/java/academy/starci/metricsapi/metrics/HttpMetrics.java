package academy.starci.metricsapi.metrics;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import java.time.Duration;
import org.springframework.stereotype.Component;

/**
 * Declares the HTTP counter and duration timer on the shared MeterRegistry.
 * Each metric type answers a different question, so it needs its own primitive.
 */
@Component
public class HttpMetrics {

    private final MeterRegistry registry;

    public HttpMetrics(MeterRegistry registry) {
        this.registry = registry;
    }

    // Counter: monotonically increasing total of HTTP requests, tagged by bounded labels.
    public void incRequests(String method, String route, String statusCode) {
        Counter.builder("http_requests_total")
                .description("Total HTTP requests")
                .tags("method", method, "route", route, "status_code", statusCode)
                .register(registry)
                .increment();
    }

    // Timer/histogram: distributes latency into fixed buckets so quantiles are computed at query time.
    public void observeDuration(String method, String route, Duration duration) {
        Timer.builder("http_request_duration_seconds")
                .description("HTTP request latency in seconds")
                .tags("method", method, "route", route)
                .serviceLevelObjectives(
                        Duration.ofMillis(5), Duration.ofMillis(10), Duration.ofMillis(25),
                        Duration.ofMillis(50), Duration.ofMillis(100), Duration.ofMillis(250),
                        Duration.ofMillis(500), Duration.ofSeconds(1), Duration.ofMillis(2500),
                        Duration.ofSeconds(5))
                .register(registry)
                .record(duration);
    }
}
