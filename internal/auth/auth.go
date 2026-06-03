package auth

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v5"
)

const userIDContextKey = "user_id"

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RefreshToken struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	TokenHash []byte     `db:"token_hash"`
	FamilyID  uuid.UUID  `db:"family_id"`
	ExpiresAt time.Time  `db:"expires_at"`
	RevokedAt *time.Time `db:"revoked_at"`
	CreatedAt time.Time  `db:"created_at"`
}

type RefreshTokenRotation struct {
	CurrentTokenID uuid.UUID
	NewToken       RefreshToken
	ReplacedByID   uuid.UUID
}

func SetUserID(c *echo.Context, userID uuid.UUID) {
	c.Set(userIDContextKey, userID)
}

func UserID(c *echo.Context) (uuid.UUID, bool) {
	userID, ok := c.Get(userIDContextKey).(uuid.UUID)
	return userID, ok
}
