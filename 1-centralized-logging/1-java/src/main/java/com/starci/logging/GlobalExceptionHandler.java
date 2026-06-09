package com.starci.logging;

import java.util.Map;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

/**
 * Global exception handler — normalizes every unhandled exception to the
 * cross-language envelope {statusCode, message} so the HTTP contract matches
 * the TypeScript, C#, and Go implementations. Sentry still captures the
 * exception through its own resolver; this only shapes the HTTP body.
 */
@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(Exception.class)
    public ResponseEntity<Map<String, Object>> handle(Exception ex) {
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
            .body(Map.of("statusCode", 500, "message", "Internal server error"));
    }
}
