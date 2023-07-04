package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"indexstorm/go-api-boilerplate/internal/api"
	"indexstorm/go-api-boilerplate/internal/repository"
	"indexstorm/go-api-boilerplate/pkg/db"
	"indexstorm/go-api-boilerplate/pkg/env"
	"os"
)

type queryHook struct {
	enabled bool
	debug   bool
	logger  zerolog.Logger
}

type application struct {
	debug    bool
	logger   *zerolog.Logger
	database *bun.DB
	repo     repository.Repository
	server   api.Server
}

func main() {
	app := newApplication()
	defer app.database.Close()

	app.repo.Bootstrap()
	defer app.server.Shutdown()
	app.server.Start(":1337")
}

func newApplication() *application {
	logger := newLogger()
	isDebug := env.Bool("APP_DEBUG")
	dbHook := &queryHook{
		enabled: isDebug,
		debug:   env.Bool("APP_BUN_DEBUG"),
		logger:  zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
	sslMode := "require"
	if certPath := os.Getenv("POSTGRES_CERT_PATH"); certPath != "" {
		sslMode = "verify-full"
	}
	database, err := db.NewConnection(
		fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DB"),
			sslMode,
		),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("db:connect")
	}
	if err = database.Ping(); err != nil {
		database.Close()
		logger.Fatal().Err(err).Msg("db:ping")
	}
	database.AddQueryHook(dbHook)
	repo := repository.NewPostgresRepo(database, logger)
	server := api.NewServer(
		isDebug,
		repo,
		logger.With().Str("cat", "srv").Logger(),
	)
	return &application{
		debug:    isDebug,
		logger:   logger,
		database: database,
		repo:     repo,
		server:   server,
	}
}

func newLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	return &logger
}

func (h *queryHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (h *queryHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	if !h.enabled {
		return
	}
	if !h.debug {
		switch event.Err {
		case nil, sql.ErrNoRows, sql.ErrTxDone:
			return
		}
	}
	var logEvent *zerolog.Event
	if event.Err != nil {
		logEvent = h.logger.Err(event.Err)
	} else {
		logEvent = h.logger.Info()
	}
	logEvent.Str("op", event.Operation()).Str("q", event.Query).Str("cat", "sql").Msg("")
}
