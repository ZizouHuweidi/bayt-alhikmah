using System.ComponentModel.DataAnnotations;

namespace Maktba.Domain;

public class Author
{
    public Guid Id { get; set; }
    
    [MaxLength(200)]
    public string Name { get; set; } = string.Empty;
    
    public string? Bio { get; set; }

    public ICollection<Source> Sources { get; set; } = new List<Source>();
}
