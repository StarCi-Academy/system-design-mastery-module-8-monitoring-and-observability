using System.ComponentModel.DataAnnotations;

namespace MetricsApi.Cats;

/// <summary>Demo resource returned by GET /cats and created by POST /cats.</summary>
public record Cat(int Id, string Name, int Age);

/// <summary>Request body for POST /cats. Name is required; Age must be present.</summary>
public class CreateCatRequest
{
    [Required(ErrorMessage = "name should not be empty")]
    [MinLength(1, ErrorMessage = "name should not be empty")]
    public string? Name { get; set; }

    [Required(ErrorMessage = "age must be an integer number")]
    public int? Age { get; set; }
}
