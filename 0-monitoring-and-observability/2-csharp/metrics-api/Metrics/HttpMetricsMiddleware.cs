using System.Diagnostics;
using Microsoft.AspNetCore.Routing;

namespace MetricsApi.Metrics;

/// <summary>
/// Records the counter + histogram when the response finishes, reading the real
/// status and the bounded route template (never the raw path).
/// </summary>
public class HttpMetricsMiddleware
{
    private readonly RequestDelegate _next;

    public HttpMetricsMiddleware(RequestDelegate next) => _next = next;

    public async Task InvokeAsync(HttpContext context)
    {
        var stopwatch = Stopwatch.StartNew();
        try
        {
            await _next(context);
        }
        finally
        {
            // Runs after the handler so the real status code (incl. 400) is known.
            stopwatch.Stop();
            var route = ResolveRoute(context);
            var method = context.Request.Method;
            var statusCode = context.Response.StatusCode.ToString();
            AppMetrics.HttpRequestsTotal
                .WithLabels(method, route, statusCode)
                .Inc();
            AppMetrics.HttpRequestDurationSeconds
                .WithLabels(method, route)
                .Observe(stopwatch.Elapsed.TotalSeconds);
        }
    }

    // Collapse dynamic params to the route template to keep cardinality bounded.
    private static string ResolveRoute(HttpContext context)
        => context.GetEndpoint() is RouteEndpoint endpoint
            ? "/" + endpoint.RoutePattern.RawText?.TrimStart('/')
            : context.Request.Path.Value ?? "unknown";
}
