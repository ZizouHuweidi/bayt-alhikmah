package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	issuer    string
	audience  string
	private   ed25519.PrivateKey
	public    ed25519.PublicKey
	accessTTL time.Duration
}

type AccessClaims struct {
	Username string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

func NewTokenManager(issuer, audience, encodedPrivateKey string, accessTTL time.Duration) (*TokenManager, error) {
	var private ed25519.PrivateKey
	if encodedPrivateKey == "" {
		_, generated, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		private = generated
	} else {
		decoded, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
		if err != nil {
			return nil, err
		}
		switch len(decoded) {
		case ed25519.SeedSize:
			private = ed25519.NewKeyFromSeed(decoded)
		case ed25519.PrivateKeySize:
			private = ed25519.PrivateKey(decoded)
		default:
			return nil, errors.New("ed25519 private key must be a base64-encoded 32-byte seed or 64-byte private key")
		}
	}

	public, ok := private.Public().(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("failed to derive ed25519 public key")
	}

	return &TokenManager{
		issuer:    issuer,
		audience:  audience,
		private:   private,
		public:    public,
		accessTTL: accessTTL,
	}, nil
}

func (m *TokenManager) CreateAccessToken(user User) (string, error) {
	jti, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	claims := AccessClaims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{m.audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        jti.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(m.private)
}

func (m *TokenManager) VerifyAccessToken(rawToken string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodEdDSA {
			return nil, errors.New("unexpected JWT signing method")
		}
		return m.public, nil
	}, jwt.WithIssuer(m.issuer), jwt.WithAudience(m.audience))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func NewRefreshToken() (string, []byte, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	hash := HashRefreshToken(token)
	return token, hash, nil
}

func HashRefreshToken(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}
