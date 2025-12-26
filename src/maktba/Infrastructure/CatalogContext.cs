using Maktba.Domain;
using Microsoft.EntityFrameworkCore;

namespace Maktba.Infrastructure;

public class CatalogContext : DbContext
{
    public CatalogContext(DbContextOptions<CatalogContext> options) : base(options)
    {
    }

    public DbSet<Source> Sources { get; set; } = null!;
    public DbSet<Author> Authors { get; set; } = null!;
    public DbSet<Taxonomy> Taxonomies { get; set; } = null!;

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // Source - Author (Many-to-Many)
        modelBuilder.Entity<Source>()
            .HasMany(s => s.Authors)
            .WithMany(a => a.Sources)
            .UsingEntity("SourceAuthors");

        // Source - Taxonomy (Many-to-Many)
        modelBuilder.Entity<Source>()
            .HasMany(s => s.Taxonomies)
            .WithMany(t => t.Sources)
            .UsingEntity("SourceTaxonomies");
            
        // Indexes
        modelBuilder.Entity<Source>()
            .HasIndex(s => s.Title);
    }
}
