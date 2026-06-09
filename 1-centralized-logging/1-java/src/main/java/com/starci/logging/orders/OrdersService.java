package com.starci.logging.orders;

import java.util.Map;
import java.util.Optional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

/** Orders service — generates a trace and writes demo logs. */
@Service
public class OrdersService {

    private static final Logger log = LoggerFactory.getLogger(OrdersService.class);

    /** Simulate a successful order; log INFO with a trace for Loki lookup. */
    public Map<String, String> accept() {
        String trace = "trace-" + System.currentTimeMillis();
        // HOSTNAME in containers helps distinguish replicas in demos.
        String host = Optional.ofNullable(System.getenv("HOSTNAME")).orElse("local");
        log.info("[{}] Order accepted trace={}", host, trace);
        return Map.of("status", "accepted", "trace", trace);
    }

    /** Simulate an order failure; log ERROR then throw a real exception for Sentry. */
    public Map<String, String> fail() {
        String host = Optional.ofNullable(System.getenv("HOSTNAME")).orElse("local");
        log.error("[{}] Simulated order failure for Loki + Sentry demo", host);
        // Throw a real exception — the Sentry resolver catches this and reports to Sentry.
        throw new IllegalStateException("Simulated order processing failure");
    }
}
