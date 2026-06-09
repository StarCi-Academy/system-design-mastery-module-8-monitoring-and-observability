using System.Net.Http.Json;
using Microsoft.AspNetCore.Mvc;

namespace ConsulApi;

/// <summary>
/// Consul proxy controller — register / health / deregister via the Consul Agent HTTP API.
/// </summary>
[ApiController]
[Route("consul")]
public class ConsulController : ControllerBase
{
    private readonly HttpClient _http;
    private readonly string _baseUrl;

    public ConsulController(HttpClient http, IConfiguration config)
    {
        _http = http;
        _baseUrl = config["Consul:BaseUrl"]
            ?? Environment.GetEnvironmentVariable("CONSUL_BASE_URL")
            ?? "http://consul:8500";
    }

    /// <summary>Register a service instance into the Consul catalog via the Agent API.</summary>
    [HttpPost("register")]
    public async Task<IActionResult> Register([FromBody] RegisterDto body)
    {
        // Consul requires PascalCase fields; camelCase registers silently empty.
        var payload = new { ID = body.Id, Name = body.Name, Address = body.Address, Port = body.Port };
        await _http.PutAsJsonAsync($"{_baseUrl}/v1/agent/service/register", payload);
        return StatusCode(201, new { status = "registered" });
    }

    /// <summary>List passing instances of a service by logical name.</summary>
    [HttpGet("health/{service}")]
    public async Task<IActionResult> Health(string service)
    {
        // Query by logical service name; the health endpoint joins catalog + check state.
        var response = await _http.GetStringAsync($"{_baseUrl}/v1/health/service/{service}");
        return Content(response, "application/json");
    }

    /// <summary>Deregister a service instance from the Agent by its exact instance ID.</summary>
    [HttpPost("deregister/{id}")]
    public async Task<IActionResult> Deregister(string id)
    {
        // Deregister by exact instance ID, removing it from the catalog immediately.
        await _http.PutAsync($"{_baseUrl}/v1/agent/service/deregister/{id}", null);
        return StatusCode(201, new { status = "deregistered" });
    }
}
