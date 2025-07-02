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
            app.MapGet("/admin/users", async ([FromServices] ApplicationDbContext dbContext) =>
            {
                var users = await dbContext.Users.ToListAsync();
                return Results.Ok(users);
            }).RequireAuthorization("AdminPolicy");

            app.MapPut("/admin/users/{id}/role", async (Guid id, [FromBody] UserRole newRole, [FromServices] ApplicationDbContext dbContext) =>
            {
                var user = await dbContext.Users.FindAsync(id);
                if (user == null)
                {
                    return Results.NotFound("User not found.");
                }

                user.Role = newRole;
                await dbContext.SaveChangesAsync();

                return Results.Ok(new { message = $"User {user.Email} role updated to {newRole}" });
            }).RequireAuthorization("AdminPolicy");
        }
    }
}