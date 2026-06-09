// Entry point — registers a typed HttpClient and the Consul proxy controller.
var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddHttpClient();

var app = builder.Build();

app.MapControllers();

// Bind 0.0.0.0 so the service is reachable inside Docker; PORT defaults to 3000.
var port = Environment.GetEnvironmentVariable("PORT") ?? "3000";
app.Run($"http://0.0.0.0:{port}");
