package api

import "github.com/gofiber/fiber/v2/middleware/cors"

func (s *server) configureEndpoints() {
	api := s.app.Group("/api", cors.New(), s.newJWTMiddleware())
	api.Get("/users/me", s.getCurrentUser)
}
