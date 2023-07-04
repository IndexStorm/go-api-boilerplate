package api

import (
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog"
	"os"
)

func (s *server) configureLogging() {
	requestLogger := zerolog.New(os.Stderr).With().Timestamp().Str("cat", "srv").Logger()
	var fmt = "[${ip}] ${status} - ${latency} ${method} ${path} \"${error}\"\n"
	if s.debug {
		fmt = "[${ip}] [${time}] ${status} - ${latency} ${method} ${path} \"${error}\"\n"
	}
	s.app.Use(
		logger.New(
			logger.Config{
				TimeZone: "Etc/GMT-3",
				Output:   requestLogger,
				Format:   fmt,
			},
		),
	)
}
