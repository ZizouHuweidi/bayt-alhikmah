using Maktba.Infrastructure;
using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddOpenApi();

builder.Services.AddDbContext<CatalogContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("CatalogDatabase")));

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
}

app.UseHttpsRedirection();

using (var scope = app.Services.CreateScope())
{
    var db = scope.ServiceProvider.GetRequiredService<CatalogContext>();
    db.Database.Migrate();
}

app.MapPost("/sources", Maktba.Features.CreateSource.Handle);
app.MapGet("/sources/{id}", Maktba.Features.GetSource.Handle);

app.Run();
