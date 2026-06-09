using Microsoft.AspNetCore.Mvc;
using MetricsApi.Metrics;
using Prometheus;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddControllers();

// Shape model-validation failures into the shared error contract:
// { "statusCode": 400, "message": ["..."], "error": "Bad Request" }.
builder.Services.Configure<ApiBehaviorOptions>(options =>
{
    options.InvalidModelStateResponseFactory = context =>
    {
        var messages = context.ModelState
            .SelectMany(entry => entry.Value!.Errors)
            .Select(error => error.ErrorMessage)
            .Where(message => !string.IsNullOrWhiteSpace(message))
            .ToArray();
        return new BadRequestObjectResult(new
        {
            statusCode = 400,
            message = messages,
            error = "Bad Request",
        });
    };
});

var app = builder.Build();

// Record duration + status for every request via the custom middleware.
app.UseRouting();
app.UseMiddleware<HttpMetricsMiddleware>();
app.MapControllers();

// Expose the default registry in Prometheus text exposition format (pull endpoint).
app.MapMetrics("/metrics");

app.Run("http://0.0.0.0:3000");
