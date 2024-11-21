package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"quotablegooofs.prajjmon.net/internal/models"
)

type application struct {
	logger *slog.Logger
	jokes  *models.JokeModel
	quotes *models.QuoteModel
}

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")

	dsn := flag.String("dsn", "postgres://quotablegoof:localdevpassword@localhost:5432/quotablegoofs", "PGX data source name")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	dbpool, err := openDbPool(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("DB ping successful, dbpool is open")

	defer dbpool.Close()

	app := application{
		logger: logger,
		jokes:  &models.JokeModel{DbPool: dbpool},
		quotes: &models.QuoteModel{DbPool: dbpool},
	}

	server := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", slog.String("addr", server.Addr))

	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDbPool(dsn string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}
