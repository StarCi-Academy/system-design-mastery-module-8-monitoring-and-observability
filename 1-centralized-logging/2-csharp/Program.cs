// ASP.NET Core logging-api — Serilog ships structured logs to Loki; Sentry captures exceptions.
using Serilog;
using Serilog.Sinks.Grafana.Loki;

var builder = WebApplication.CreateBuilder(args);

var lokiUrl = Environment.GetEnvironmentVariable("LOKI_URL") ?? "http://loki:3100";
var appName = Environment.GetEnvironmentVariable("LOG_APP_NAME") ?? "logging-api";
var env = Environment.GetEnvironmentVariable("LOG_ENV") ?? "lab";

// Serilog: Console for stdout + Grafana Loki sink with a bounded label set {app, env}.
Log.Logger = new LoggerConfiguration()
    .Enrich.FromLogContext()
    .WriteTo.Console()
    .WriteTo.GrafanaLoki(
        lokiUrl,
        labels: new[]
        {
            new LokiLabel { Key = "app", Value = appName },
            new LokiLabel { Key = "env", Value = env },
        })
    .CreateLogger();

builder.Host.UseSerilog();

// Sentry — initialized via env DSN; disables itself when the DSN is empty.
builder.WebHost.UseSentry(o =>
{
    o.Dsn = Environment.GetEnvironmentVariable("SENTRY_DSN") ?? string.Empty;
    o.TracesSampleRate = 1.0;
});

// Bind 0.0.0.0:3000 inside the container.
builder.WebHost.UseUrls("http://0.0.0.0:3000");

var app = builder.Build();

// GET /api/orders — accepted order; logs INFO with a trace for the Loki query demo.
app.MapGet("/api/orders", () =>
{
    var trace = $"trace-{DateTimeOffset.UtcNow.ToUnixTimeMilliseconds()}";
    var host = Environment.GetEnvironmentVariable("HOSTNAME") ?? "local";
    Log.Information("[{Host}] Order accepted trace={Trace}", host, trace);
    return Results.Ok(new { status = "accepted", trace });
});

// GET /api/orders/error — deliberate failure; logs ERROR, reports to Sentry, returns the
// cross-language {statusCode, message} envelope.
app.MapGet("/api/orders/error", () =>
{
    var host = Environment.GetEnvironmentVariable("HOSTNAME") ?? "local";
    Log.Error("[{Host}] Simulated order failure for Loki + Sentry demo", host);
    var ex = new InvalidOperationException("Simulated order processing failure");
    SentrySdk.CaptureException(ex);
    return Results.Json(new { statusCode = 500, message = "Internal server error" }, statusCode: 500);
});

app.Run();
