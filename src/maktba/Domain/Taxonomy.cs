using System.ComponentModel.DataAnnotations;

namespace Maktba.Domain;

public enum TaxonomyType
{
    Topic,
    Tag,
    Person, // as subject
    Era,
    Region
}

public class Taxonomy
{
    public Guid Id { get; set; }
    
    [MaxLength(100)]
    public string Name { get; set; } = string.Empty;
    
    public TaxonomyType Type { get; set; }
    
    public ICollection<Source> Sources { get; set; } = new List<Source>();
}
