package api

import (
	"github.com/gofiber/fiber/v2"
)

func (s *server) getCurrentUser(c *fiber.Ctx) error {
	claims := s.extractAuthToken(c)
	u, err := s.repo.FetchUserByID(claims.UserID, c.Context())
	if err != nil {
		return err
	}
	return c.JSON(u)
}
