using System;
using System.Threading.Tasks;
using BaytAlHikmah.Application.Interfaces;
using BaytAlHikmah.Domain.Entities;
using BaytAlHikmah.Domain.Repositories;

namespace BaytAlHikmah.Infrastructure.Services;

public class AuthService : IAuthService
{
    private readonly IUserRepository _userRepository;
    private readonly IPasswordHasher _passwordHasher;
    private readonly ITokenProvider _tokenProvider;

    public AuthService(IUserRepository userRepository, IPasswordHasher passwordHasher, ITokenProvider tokenProvider)
    {
        _userRepository = userRepository;
        _passwordHasher = passwordHasher;
        _tokenProvider = tokenProvider;
    }

    public async Task<Guid> RegisterAsync(string email, string password, string firstName, string lastName)
    {
        if (!await _userRepository.IsEmailUniqueAsync(email))
        {
            throw new InvalidOperationException("Email already exists.");
        }

        var hash = _passwordHasher.HashPassword(password, out var salt);
        var user = User.Create(email, hash, salt, firstName, lastName);

        await _userRepository.AddAsync(user);
        await _userRepository.SaveChangesAsync(default);
        
        return user.Id;
    }

    public async Task<string> LoginAsync(string email, string password)
    {
        var user = await _userRepository.GetByEmailAsync(email);
        if (user == null)
        {
            throw new InvalidOperationException("Invalid credentials.");
        }

        if (!_passwordHasher.VerifyPassword(password, user.PasswordHash, user.Salt))
        {
             throw new InvalidOperationException("Invalid credentials.");
        }

        return _tokenProvider.GenerateToken(user);
    }
}
