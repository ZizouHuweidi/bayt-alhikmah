using Maktba.Domain;
using Maktba.Infrastructure;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

namespace Maktba.Features;

public static class CreateSource
{
    public record Request(string Title, SourceType Type, string? Description, string? Url);
    public record Response(Guid Id, string Title);

    public static async Task<Created<Response>> Handle(
        [FromBody] Request request,
        CatalogContext db,
        CancellationToken ct)
    {
        var source = new Source
        {
            Id = Guid.CreateVersion7(),
            Title = request.Title,
            Type = request.Type,
            Description = request.Description,
            Url = request.Url,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        db.Sources.Add(source);
        await db.SaveChangesAsync(ct);

        return TypedResults.Created($"/sources/{source.Id}", new Response(source.Id, source.Title));
    }
}
