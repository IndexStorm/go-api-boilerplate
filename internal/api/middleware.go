package api

import (
	"github.com/gofiber/fiber/v2"
)

func (s *server) newJWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if len(token) < 64 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		claims, err := s.verifyAuthToken(token, c.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		c.Context().SetUserValue(AuthClaimsContextKey, claims)
		return c.Next()
	}
}
