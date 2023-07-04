package db

import (
	"crypto/x509"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"os"
	"runtime"
	"time"
)

func NewConnection(conn string) (*bun.DB, error) {
	config, err := pgx.ParseConfig(conn)
	if err != nil {
		return nil, err
	}
	config.DefaultQueryExecMode = pgx.QueryExecModeExec
	config.ConnectTimeout = time.Second * 20
	if certPath := os.Getenv("POSTGRES_CERT_PATH"); certPath != "" {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(certPath)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, errors.New("failed to append PEM")
		}
		config.TLSConfig.RootCAs = rootCertPool
		config.TLSConfig.InsecureSkipVerify = false
	}
	db := bun.NewDB(stdlib.OpenDB(*config), pgdialect.New())
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxOpenConns)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}
