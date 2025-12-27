using Maktba.Infrastructure;
using Microsoft.EntityFrameworkCore;
using OpenTelemetry.Logs;
using OpenTelemetry.Metrics;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using Serilog;
using Serilog.Events;

// Configure Serilog early for startup logging
Log.Logger = new LoggerConfiguration()
    .MinimumLevel.Override("Microsoft", LogEventLevel.Information)
    .MinimumLevel.Override("Microsoft.EntityFrameworkCore", LogEventLevel.Warning)
    .Enrich.FromLogContext()
    .Enrich.WithEnvironmentName()
    .Enrich.WithThreadId()
    .WriteTo.Console(outputTemplate: "[{Timestamp:HH:mm:ss} {Level:u3}] {Message:lj} {Properties:j}{NewLine}{Exception}")
    .WriteTo.OpenTelemetry(options =>
    {
        options.Endpoint = "http://tempo:4317";
        options.Protocol = Serilog.Sinks.OpenTelemetry.OtlpProtocol.Grpc;
        options.ResourceAttributes = new Dictionary<string, object>
        {
            ["service.name"] = "maktba"
        };
    })
    .CreateLogger();

try
{
    Log.Information("Starting Maktba service...");
    
    var builder = WebApplication.CreateBuilder(args);
    
    // Use Serilog for all logging
    builder.Host.UseSerilog();

    builder.Services.AddOpenApi();

    // Configure OpenTelemetry with tracing and metrics
    var serviceName = "maktba";
    var serviceVersion = "1.0.0";
    
    builder.Services.AddOpenTelemetry()
        .ConfigureResource(resource => resource
            .AddService(serviceName: serviceName, serviceVersion: serviceVersion)
            .AddAttributes(new Dictionary<string, object>
            {
                ["deployment.environment"] = builder.Environment.EnvironmentName
            }))
        .WithTracing(tracing => tracing
            .AddAspNetCoreInstrumentation(options =>
            {
                options.RecordException = true;
            })
            .AddHttpClientInstrumentation()
            .AddEntityFrameworkCoreInstrumentation()
            .AddOtlpExporter(options =>
            {
                options.Endpoint = new Uri("http://tempo:4317");
                options.Protocol = OpenTelemetry.Exporter.OtlpExportProtocol.Grpc;
            }))
        .WithMetrics(metrics => metrics
            .AddAspNetCoreInstrumentation()
            .AddHttpClientInstrumentation()
            .AddRuntimeInstrumentation()
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
        
        // Apply migrations at startup in development with retry logic
        var maxRetries = 10;
        var delay = TimeSpan.FromSeconds(3);
        
        for (int i = 0; i < maxRetries; i++)
        {
            try
            {
                using var scope = app.Services.CreateScope();
                var db = scope.ServiceProvider.GetRequiredService<CatalogContext>();
                db.Database.Migrate();
                Log.Information("Database migrated successfully");
                break;
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Migration attempt {Attempt} failed", i + 1);
                if (i == maxRetries - 1) throw;
                Thread.Sleep(delay);
            }
        }
    }

    app.UseSerilogRequestLogging(options =>
    {
        options.EnrichDiagnosticContext = (diagnosticContext, httpContext) =>
        {
            diagnosticContext.Set("RequestHost", httpContext.Request.Host.Value);
            diagnosticContext.Set("UserAgent", httpContext.Request.Headers["User-Agent"].FirstOrDefault());
        };
    });

    app.MapGet("/healthz", () => Results.Ok("Healthy"));

    app.MapPost("/sources", Maktba.Features.CreateSource.Handle);
    app.MapGet("/sources/{id}", Maktba.Features.GetSource.Handle);

    Log.Information("Maktba service started successfully");
    app.Run();
}
catch (Exception ex)
{
    Log.Fatal(ex, "Application terminated unexpectedly");
}
finally
{
    Log.CloseAndFlush();
}
