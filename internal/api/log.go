package api

import (
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (s *server) configureLogging() {
	var fmt = "[${ip}] ${status} - ${latency} ${method} ${path} \"${error}\"\n"
	if s.debug {
		fmt = "[${ip}] [${time}] ${status} - ${latency} ${method} ${path} \"${error}\"\n"
	}
	s.app.Use(
		logger.New(
			logger.Config{
				TimeZone: "Etc/GMT-3",
				Output:   s.logger,
				Format:   fmt,
			},
		),
	)
}
