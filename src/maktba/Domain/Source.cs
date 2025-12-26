using System.ComponentModel.DataAnnotations;

namespace Maktba.Domain;

public enum SourceType
{
    Book,
    Paper,
    Article,
    Video,
    Podcast
}

public class Source
{
    public Guid Id { get; set; }
    
    [MaxLength(500)]
    public string Title { get; set; } = string.Empty;
    
    public SourceType Type { get; set; }
    
    public string? Description { get; set; }
    public string? CoverUrl { get; set; }
    public string? Url { get; set; }
    
    public DateTime? PublishedDate { get; set; }
    
    public ICollection<Author> Authors { get; set; } = new List<Author>();
    public ICollection<Taxonomy> Taxonomies { get; set; } = new List<Taxonomy>();
    
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
