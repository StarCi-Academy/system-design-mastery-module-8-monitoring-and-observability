package com.starci.tracing;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.trace.Tracer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * Exposes the OpenTelemetry instance and a named Tracer as Spring beans.
 * The SDK, OTLP exporter, BatchSpanProcessor, and inbound HTTP auto-instrumentation
 * are all provided by the OpenTelemetry Java agent (JAVA_TOOL_OPTIONS=-javaagent:...)
 * before the JVM loads any application class; this @Configuration only surfaces the
 * GlobalOpenTelemetry singleton so services can inject it via DI.
 */
@Configuration
public class OpenTelemetryConfig {

    // The javaagent registers itself as the GlobalOpenTelemetry provider at JVM startup
    // (before Spring's ApplicationContext is created), so GlobalOpenTelemetry.get()
    // always returns the fully wired SDK — no in-process SDK initialization needed here.
    @Bean
    public OpenTelemetry openTelemetry() {
        // autoConfigure() reads OTEL_* env, wires the OTLP HTTP exporter and a BatchSpanProcessor.
        return GlobalOpenTelemetry.get();
    }

    // Expose a named Tracer so services can open spans under this sub-feature namespace.
    @Bean
    public Tracer checkoutTracer(OpenTelemetry openTelemetry) {
        return openTelemetry.getTracer("spring-tracing-app-checkout");
    }
}
