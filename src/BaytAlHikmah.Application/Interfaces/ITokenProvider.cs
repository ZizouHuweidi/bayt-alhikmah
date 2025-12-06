using BaytAlHikmah.Domain.Entities;

namespace BaytAlHikmah.Application.Interfaces;

public interface ITokenProvider
{
    string GenerateToken(User user);
}
