package academy.starci.metricsapi.metrics;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.time.Duration;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;
import org.springframework.web.servlet.HandlerMapping;

/**
 * Records the counter + histogram after the response finishes, reading the real
 * status and the bounded route template (never the raw path).
 */
@Component
public class HttpMetricsFilter extends OncePerRequestFilter {

    private final HttpMetrics metrics;

    public HttpMetricsFilter(HttpMetrics metrics) {
        this.metrics = metrics;
    }

    @Override
    protected void doFilterInternal(HttpServletRequest req, HttpServletResponse res,
                                    FilterChain chain) throws ServletException, IOException {
        long start = System.nanoTime();
        try {
            chain.doFilter(req, res);
        } finally {
            // Read status and resolve the route only after the handler ran.
            String route = resolveRoute(req);
            Duration duration = Duration.ofNanos(System.nanoTime() - start);
            String statusCode = String.valueOf(res.getStatus());
            metrics.incRequests(req.getMethod(), route, statusCode);
            metrics.observeDuration(req.getMethod(), route, duration);
        }
    }

    // Collapse dynamic params to a route template to keep cardinality bounded.
    private String resolveRoute(HttpServletRequest req) {
        Object pattern = req.getAttribute(HandlerMapping.BEST_MATCHING_PATTERN_ATTRIBUTE);
        if (pattern instanceof String s && !s.isBlank()) {
            return s;
        }
        return req.getRequestURI();
    }
}
