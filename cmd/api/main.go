package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"indexstorm/go-api-boilerplate/internal/api"
	"indexstorm/go-api-boilerplate/internal/repository"
	"indexstorm/go-api-boilerplate/pkg/db"
	"indexstorm/go-api-boilerplate/pkg/env"
	"os"
	"strconv"
)

type queryHook struct {
	enabled bool
	debug   bool
	logger  zerolog.Logger
}

type application struct {
	debug    bool
	database *bun.DB
	repo     repository.Repository
	server   api.Server
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	app := newApplication(logger)
	defer app.database.Close()

	app.repo.Bootstrap()
	app.database.AddQueryHook(
		&queryHook{
			enabled: app.debug,
			debug:   env.Bool("APP_BUN_DEBUG"),
			logger:  logger.With().Str("cat", "sql").Logger(),
		},
	)
	defer app.server.Shutdown()
	app.server.Start(":1337")
}

func newApplication(logger zerolog.Logger) *application {
	isDebug := env.Bool("APP_DEBUG")
	database, err := db.NewConnection(
		fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_SSL_MODE"),
		),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("db:connect")
	}
	if err = database.Ping(); err != nil {
		database.Close()
		logger.Fatal().Err(err).Msg("db:ping")
	}
	repo := repository.NewPostgresRepo(database, logger)
	server := api.NewServer(
		isDebug,
		repo,
		logger,
	)
	return &application{
		debug:    isDebug,
		database: database,
		repo:     repo,
		server:   server,
	}
}

func (h *queryHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (h *queryHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	if !h.enabled {
		return
	}
	if !h.debug {
		switch {
		case event.Err == nil,
			errors.Is(event.Err, sql.ErrNoRows),
			errors.Is(event.Err, sql.ErrTxDone):
			return
		}
	}
	var logEvent *zerolog.Event
	if event.Err != nil {
		logEvent = h.logger.Err(event.Err)
	} else {
		logEvent = h.logger.Info()
	}
	logEvent.Str("op", event.Operation()).Str("q", event.Query).Msg("")
}
