package com.starci.tracing;

import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.Scope;
import org.springframework.stereotype.Service;

import java.util.concurrent.ThreadLocalRandom;

/**
 * Emits a child span under the inbound HTTP span and returns the active traceId.
 */
@Service
public class CheckoutService {

    private final Tracer tracer;

    public CheckoutService(Tracer tracer) {
        this.tracer = tracer;
    }

    public CheckoutResult simulateCheckout() {
        // startSpan() creates the span; makeCurrent() installs it as the active context.
        Span inventorySpan = tracer.spanBuilder("checkout.inventory_reserve").startSpan();
        try (Scope scope = inventorySpan.makeCurrent()) {
            int delayMs = 50 + ThreadLocalRandom.current().nextInt(400);
            inventorySpan.setAttribute("checkout.inventory.delay_ms", delayMs);
            Thread.sleep(delayMs);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } finally {
            inventorySpan.end();
        }
        // After the child span closes, the active span is the inbound HTTP span again.
        String traceId = Span.current().getSpanContext().getTraceId();
        return new CheckoutResult("Checkout completed", traceId);
    }
}
