using BaytAlHikmah.Application.Interfaces;
using BCrypt.Net;

namespace BaytAlHikmah.Infrastructure.Authentication;

public class PasswordHasher : IPasswordHasher
{
    public string HashPassword(string password, out string salt)
    {
        salt = BCrypt.Net.BCrypt.GenerateSalt();
        return BCrypt.Net.BCrypt.HashPassword(password, salt);
    }

    public bool VerifyPassword(string password, string hash, string salt)
    {
        // BCrypt stores the salt in the hash usually, but since we are extracting it explicitly
        // we essentially just need to verify against the hash. The salt param here is ensuring
        // we are using the one associated with the user, though BCrypt.Verify() parses the salt from the hash string itself.
        // However, to be strictly consistent with the "Salt" property usage if we used it to Generate the hash:
        // We can just call Verify.
        return BCrypt.Net.BCrypt.Verify(password, hash);
    }
}
