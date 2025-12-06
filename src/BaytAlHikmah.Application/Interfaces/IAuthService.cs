using System;
using System.Threading.Tasks;

namespace BaytAlHikmah.Application.Interfaces;

public interface IAuthService
{
    Task<Guid> RegisterAsync(string email, string password, string firstName, string lastName);
    Task<string> LoginAsync(string email, string password);
}
