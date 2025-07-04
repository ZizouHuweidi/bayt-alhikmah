using System.ComponentModel.DataAnnotations;
using BaytAlHikmah.Core.Enums;

namespace BaytAlHikmah.Core.Entities
{
    public class User
    {
        public Guid Id { get; set; }

        [Required]
        [EmailAddress]
        public string Email { get; set; }

        public string? PasswordHash { get; set; }

        public string? GoogleId { get; set; }

        public string? FullName { get; set; }

        public DateTime CreatedAt { get; set; } = DateTime.UtcNow;

        public UserRole Role { get; set; } = UserRole.NormalUser;
    }
}
