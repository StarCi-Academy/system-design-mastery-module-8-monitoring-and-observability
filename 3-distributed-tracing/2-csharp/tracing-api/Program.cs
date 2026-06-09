// Distributed tracing lab (ASP.NET Core minimal API + OpenTelemetry).
// Registers the tracing pipeline before app.Run() so AspNetCore auto-instrumentation
// wraps the middleware pipeline, then emits a manual child span via an ActivitySource.
using System.Diagnostics;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using TracingApi;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddOpenTelemetry()
    .ConfigureResource(r => r.AddService(serviceName: "nestjs-tracing-app"))
    .WithTracing(tracing => tracing
        // Auto-instrument inbound HTTP so every request emits a Span.
        .AddAspNetCoreInstrumentation()
        .AddHttpClientInstrumentation()
        // Register our own ActivitySource so manual child spans are recorded.
        .AddSource("nestjs-tracing-app-checkout")
        .AddOtlpExporter(o =>
        {
            // Reads the collector endpoint from env; no hardcoded URL.
            o.Endpoint = new Uri(builder.Configuration["OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"]
                ?? "http://localhost:4318/v1/traces");
            o.Protocol = OpenTelemetry.Exporter.OtlpExportProtocol.HttpProtobuf;
        }));

builder.Services.AddSingleton<CheckoutService>();

var app = builder.Build();

// Minimal API endpoint: business handler only, tracing wired up at startup.
app.MapGet("/checkout", async (CheckoutService service) =>
{
    var result = await service.SimulateCheckoutAsync();
    return Results.Ok(result);
});

app.MapGet("/health", () => Results.Ok(new { status = "ok" }));

app.Run();

namespace TracingApi
{
    // Holds the single ActivitySource shared by the app; its name must be
    // registered with AddSource so manual activities are recorded.
    public static class CheckoutTracing
    {
        public static readonly ActivitySource CheckoutSource = new("nestjs-tracing-app-checkout");
    }

    public record CheckoutResult(string Message, string TraceId);

    // Emits a child span under the inbound HTTP span and returns the active traceId.
    public class CheckoutService
    {
        public async Task<CheckoutResult> SimulateCheckoutAsync()
        {
            // StartActivity installs the new Activity as Activity.Current (ambient parent).
            using (var inventory = CheckoutTracing.CheckoutSource.StartActivity("checkout.inventory_reserve"))
            {
                var delayMs = 50 + Random.Shared.Next(400);
                inventory?.SetTag("checkout.inventory.delay_ms", delayMs);
                await Task.Delay(delayMs);
            }
            // After the child closes, Activity.Current is the inbound HTTP span again.
            var traceId = Activity.Current?.TraceId.ToString() ?? "";
            return new CheckoutResult("Checkout completed", traceId);
        }
    }
}
