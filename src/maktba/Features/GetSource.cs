using Maktba.Domain;
using Maktba.Infrastructure;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.EntityFrameworkCore;

namespace Maktba.Features;

public static class GetSource
{
    public record Response(Guid Id, string Title, SourceType Type, string? Description, string? Url, List<AuthorDto> Authors, List<TaxonomyDto> Taxonomies);
    public record AuthorDto(Guid Id, string Name);
    public record TaxonomyDto(Guid Id, string Name, TaxonomyType Type);

    public static async Task<Results<Ok<Response>, NotFound>> Handle(
        Guid id,
        CatalogContext db,
        CancellationToken ct)
    {
        var source = await db.Sources
            .Include(s => s.Authors)
            .Include(s => s.Taxonomies)
            .FirstOrDefaultAsync(s => s.Id == id, ct);

        if (source == null)
        {
            return TypedResults.NotFound();
        }

        var response = new Response(
            source.Id,
            source.Title,
            source.Type,
            source.Description,
            source.Url,
            source.Authors.Select(a => new AuthorDto(a.Id, a.Name)).ToList(),
            source.Taxonomies.Select(t => new TaxonomyDto(t.Id, t.Name, t.Type)).ToList()
        );

        return TypedResults.Ok(response);
    }
}
