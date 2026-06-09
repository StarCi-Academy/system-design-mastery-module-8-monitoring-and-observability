namespace ConsulApi;

/// <summary>JSON body accepted by POST /consul/register.</summary>
public record RegisterDto(string Id, string Name, string Address, int Port);
