package api

import (
	"errors"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"indexstorm/go-api-boilerplate/internal/repository"
	"time"
)

type Server interface {
	Start(address string)
	Shutdown()
}

type server struct {
	debug  bool
	app    *fiber.App
	logger zerolog.Logger
	repo   repository.Repository
}

func NewServer(isDebug bool, repo repository.Repository, logger zerolog.Logger) Server {
	app := fiber.New(
		fiber.Config{
			Prefork:               !isDebug,
			JSONEncoder:           json.Marshal,
			JSONDecoder:           json.Unmarshal,
			IdleTimeout:           time.Second * 30,
			ProxyHeader:           fiber.HeaderXForwardedFor,
			DisableStartupMessage: true,
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError
				var e *fiber.Error
				if errors.As(err, &e) {
					code = e.Code
				}
				if code != fiber.StatusNotFound {
					logger.Err(err).Msg("")
				}
				return ctx.Status(code).JSON(fiber.Map{"error": err.Error()})
			},
		},
	)
	app.Use(recover2.New())

	return &server{
		debug:  isDebug,
		app:    app,
		logger: logger,
		repo:   repo,
	}
}

func (s *server) Start(address string) {
	s.configureLogging()
	s.configureEndpoints()
	s.logger.Fatal().Err(s.app.Listen(address)).Msg("server:start")
}

func (s *server) Shutdown() {
	s.app.Shutdown()
}
