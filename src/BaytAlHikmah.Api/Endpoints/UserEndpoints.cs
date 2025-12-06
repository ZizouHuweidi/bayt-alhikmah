using BaytAlHikmah.Api.Contracts;
using BaytAlHikmah.Application.Interfaces;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Routing;

namespace BaytAlHikmah.Api.Endpoints;

public static class UserEndpoints
{
    public static void MapUserEndpoints(this IEndpointRouteBuilder app)
    {
        var group = app.MapGroup("/api/users")
            .WithTags("Users");

        group.MapPost("/register", Register)
            .WithName("RegisterUser");

        group.MapPost("/login", Login)
            .WithName("LoginUser");
    }

    private static async Task<IResult> Register(RegisterRequest request, IAuthService authService)
    {
        var userId = await authService.RegisterAsync(request.Email, request.Password, request.FirstName, request.LastName);
        return Results.Created($"/api/users/{userId}", new { Id = userId });
    }

    private static async Task<IResult> Login(LoginRequest request, IAuthService authService)
    {
        var token = await authService.LoginAsync(request.Email, request.Password);
        return Results.Ok(new { Token = token });
    }
}
