package com.starci.tracing;

/**
 * Response payload for GET /checkout: a fixed message plus the active traceId.
 */
public record CheckoutResult(String message, String traceId) {
}
