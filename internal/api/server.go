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
	srvLogger := logger.With().Str("cat", "srv").Logger()
	app := fiber.New(
		fiber.Config{
			Prefork:               !isDebug,
			JSONEncoder:           json.Marshal,
			JSONDecoder:           json.Unmarshal,
			IdleTimeout:           time.Second * 30,
			ProxyHeader:           fiber.HeaderXForwardedFor,
			BodyLimit:             1 * 1024 * 1024,
			DisableStartupMessage: true,
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				var a *apiError
				if errors.As(err, &a) {
					srvLogger.Info().
						Err(err).
						Int("code", a.Code).
						Str("uri", string(ctx.Request().RequestURI())).
						Msg("api_error")
					return ctx.Status(a.Code).JSON(a)
				}
				var f *fiber.Error
				if errors.As(err, &f) {
					srvLogger.Info().
						Err(err).
						Int("code", f.Code).
						Str("uri", string(ctx.Request().RequestURI())).
						Msg("fiber_error")
					return ctx.Status(f.Code).JSON(f)
				}
				srvLogger.Err(err).
					Any("err", err).
					Str("uri", string(ctx.Request().RequestURI())).
					Msg("internal_error")
				if isDebug {
					return ctx.Status(fiber.StatusInternalServerError).JSON(
						fiber.Map{
							"message": err.Error(),
							"error":   err,
						},
					)
				}
				return ctx.SendStatus(fiber.StatusInternalServerError)
			},
		},
	)
	app.Use(
		recover2.New(
			recover2.Config{
				EnableStackTrace: true,
			},
		),
	)
	return &server{
		debug:  isDebug,
		app:    app,
		logger: srvLogger,
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
