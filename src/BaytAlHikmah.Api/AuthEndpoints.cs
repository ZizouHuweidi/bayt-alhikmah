using System.Security.Claims;
using System.IdentityModel.Tokens.Jwt;
using System.Text;
using Microsoft.AspNetCore.Authentication;
using Microsoft.IdentityModel.Tokens;
using Microsoft.EntityFrameworkCore;
using BaytAlHikmah.Core.Entities;
using BaytAlHikmah.Infrastructure.Data;

namespace BaytAlHikmah.Api
{
    public static class AuthEndpoints
    {
        public static void MapAuthEndpoints(this WebApplication app)
        {
            app.MapPost("/register", async (UserDto userDto, ApplicationDbContext dbContext) =>
            {
                var existingUser = await dbContext.Users.SingleOrDefaultAsync(u => u.Email == userDto.Email);
                if (existingUser != null)
                {
                    return Results.Conflict("User with this email already exists.");
                }

                var user = new User
                {
                    Email = userDto.Email,
                    PasswordHash = BCrypt.Net.BCrypt.HashPassword(userDto.Password),
                    FullName = userDto.FullName
                };

                // Make the first user an admin
                if (await dbContext.Users.CountAsync() == 0)
                {
                    user.Role = BaytAlHikmah.Core.Enums.UserRole.Admin;
                }

                dbContext.Users.Add(user);
                await dbContext.SaveChangesAsync();

                return Results.Ok(new { message = "User registered successfully" });
            });

            app.MapPost("/login", async (UserDto userDto, ApplicationDbContext dbContext, IConfiguration config) =>
            {
                var user = await dbContext.Users.SingleOrDefaultAsync(u => u.Email == userDto.Email);

                if (user == null || !BCrypt.Net.BCrypt.Verify(userDto.Password, user.PasswordHash))
                {
                    return Results.Unauthorized();
                }

                var token = GenerateJwtToken(user, config);
                return Results.Ok(new { token });
            });

            app.MapGet("/account/login-google", (string returnUrl = "/") =>
            {
                var props = new AuthenticationProperties { RedirectUri = $"/account/google-callback?returnUrl={returnUrl}" };
                return Results.Challenge(props, new[] { "Google" });
            });

            app.MapGet("/account/google-callback", async (HttpContext httpContext, ApplicationDbContext dbContext, IConfiguration config) =>
            {
                var result = await httpContext.AuthenticateAsync("Google");
                if (result?.Succeeded != true)
                {
                    return Results.Unauthorized();
                }

                var email = result.Principal.FindFirstValue(ClaimTypes.Email);
                var googleId = result.Principal.FindFirstValue(ClaimTypes.NameIdentifier);
                var fullName = result.Principal.FindFirstValue(ClaimTypes.Name);

                var user = await dbContext.Users.SingleOrDefaultAsync(u => u.Email == email);
                if (user == null)
                {
                    user = new User
                    {
                        Email = email,
                        GoogleId = googleId,
                        FullName = fullName,
                    };
                    dbContext.Users.Add(user);
                    await dbContext.SaveChangesAsync();
                }

                var token = GenerateJwtToken(user, config);
                return Results.Ok(new { token });
            });
        }

        private static string GenerateJwtToken(User user, IConfiguration config)
        {
            var claims = new[]
            {
                new Claim(JwtRegisteredClaimNames.Sub, user.Id.ToString()),
                new Claim(JwtRegisteredClaimNames.Email, user.Email),
                new Claim(JwtRegisteredClaimNames.Name, user.FullName ?? string.Empty),
                new Claim(ClaimTypes.Role, user.Role.ToString())
            };

            var key = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(config["Jwt:Key"]));
            var creds = new SigningCredentials(key, SecurityAlgorithms.HmacSha256);

            var token = new JwtSecurityToken(
                issuer: config["Jwt:Issuer"],
                audience: config["Jwt:Audience"],
                claims: claims,
                expires: DateTime.Now.AddMinutes(30),
                signingCredentials: creds
            );

            return new JwtSecurityTokenHandler().WriteToken(token);
        }
    }
}
