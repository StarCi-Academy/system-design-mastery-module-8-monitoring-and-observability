using Microsoft.AspNetCore.Mvc;

namespace MetricsApi.Cats;

/// <summary>
/// REST controller for /cats. In-memory store seeded with one row (Tom) so the
/// first GET is meaningful and Prometheus sees real traffic.
/// </summary>
[ApiController]
[Route("cats")]
public class CatsController : ControllerBase
{
    private static readonly List<Cat> Cats = new() { new Cat(1, "Tom", 3) };
    private static int _nextId = 2;
    private static readonly object Lock = new();

    [HttpGet]
    public IActionResult FindAll()
    {
        lock (Lock)
        {
            return Ok(Cats);
        }
    }

    [HttpPost]
    public IActionResult Create([FromBody] CreateCatRequest body)
    {
        lock (Lock)
        {
            var cat = new Cat(_nextId++, body.Name!, body.Age!.Value);
            Cats.Add(cat);
            return StatusCode(StatusCodes.Status201Created, cat);
        }
    }
}
