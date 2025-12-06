using BaytAlHikmah.Domain.Entities;

namespace BaytAlHikmah.Domain.Repositories;

public interface IUserRepository
{
    Task<bool> IsEmailUniqueAsync(string email, CancellationToken cancellationToken = default);
    Task<User?> GetByEmailAsync(string email, CancellationToken cancellationToken = default);
    Task AddAsync(User user, CancellationToken cancellationToken = default);
    Task SaveChangesAsync(CancellationToken cancellationToken);
}
