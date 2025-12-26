using Maktba.Infrastructure;
using Microsoft.EntityFrameworkCore;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddOpenApi();

builder.Services.AddOpenTelemetry()
    .ConfigureResource(resource => resource.AddService("maktba"))
    .WithTracing(tracing => tracing
        .AddAspNetCoreInstrumentation()
        .AddOtlpExporter(options =>
        {
            options.Endpoint = new Uri("http://tempo:4317");
            options.Protocol = OpenTelemetry.Exporter.OtlpExportProtocol.Grpc;
        }));

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
