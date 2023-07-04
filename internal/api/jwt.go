package api

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"indexstorm/go-api-boilerplate/internal/repository"
	"indexstorm/go-api-boilerplate/pkg/fs"
	"indexstorm/go-api-boilerplate/pkg/nanoid"
	"os"
	"time"
)

const (
	AuthClaimsContextKey    = "AuthClaims"
	authTokenExpirationTime = 12 * time.Hour
)

var (
	jwtPublicKey  crypto.PublicKey
	jwtPrivateKey crypto.PrivateKey
)

type AuthClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func init() {
	jwtPrivateKey = ed25519.NewKeyFromSeed(fs.ReadFile(os.Getenv("APP_JWT_PATH")))
	jwtPublicKey = jwtPrivateKey.(crypto.Signer).Public()
}

func (s *server) generateAuthToken(user *repository.User, isMobile bool) (string, error) {
	return s.generateAuthTokenWithExpiration(user, isMobile, authTokenExpirationTime)
}

func (s *server) generateAuthTokenWithExpiration(user *repository.User, isMobile bool, exp time.Duration) (
	string, error,
) {
	now := time.Now()
	claims := &AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        nanoid.RandomID(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-time.Minute * 1)),
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
		},
		UserID: user.PublicID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signed, err := token.SignedString(jwtPrivateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *server) verifyAuthToken(str string, ctx context.Context) (*AuthClaims, error) {
	claims := new(AuthClaims)
	token, err := jwt.ParseWithClaims(
		str, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, errors.New("JWT signing method mismatch")
			}
			return jwtPublicKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("JWT is not valid")
	}
	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, errors.New("JWT malformed")
	}
	return claims, nil
}

func (s *server) extractAuthToken(c *fiber.Ctx) *AuthClaims {
	claims, ok := c.Context().Value(AuthClaimsContextKey).(*AuthClaims)
	if !ok {
		s.logger.Fatal().Any("ctx_value", c.Context().Value(AuthClaimsContextKey)).Msg("AuthClaims:type assert failed")
	}
	return claims
}
