using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using BaytAlHikmah.Core.Entities;
using BaytAlHikmah.Core.Enums;
using BaytAlHikmah.Infrastructure.Data;

namespace BaytAlHikmah.Api
{
    public static class AdminEndpoints
    {
        public static void MapAdminEndpoints(this WebApplication app)
        {
            var adminGroup = app.MapGroup("/admin").RequireAuthorization("AdminPolicy");

            adminGroup.MapGet("/users", async (ApplicationDbContext dbContext) =>
            {
                var users = await dbContext.Users.ToListAsync();
                return Results.Ok(users);
            });

            adminGroup.MapPut("/users/{userId}/role", async (Guid userId, [FromBody] UpdateUserRoleRequest request, ApplicationDbContext dbContext) =>
            {
                var user = await dbContext.Users.FindAsync(userId);
                if (user == null)
                {
                    return Results.NotFound("User not found.");
                }

                user.Role = request.Role;
                await dbContext.SaveChangesAsync();

                return Results.Ok(new { message = "User role updated successfully." });
            });
        }
    }

    public record UpdateUserRoleRequest(UserRole Role);
}