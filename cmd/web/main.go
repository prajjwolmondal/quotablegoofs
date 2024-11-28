package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"cloud.google.com/go/cloudsqlconn"
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
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	var missingRequiredEnvVars []string

	dbUser, err := getEnvVar("DB_USER")
	if err != nil {
		missingRequiredEnvVars = append(missingRequiredEnvVars, err.Error())
		logger.Error(err.Error())
	}

	dbPwd, err := getEnvVar("DB_PASSWORD")
	if err != nil {
		missingRequiredEnvVars = append(missingRequiredEnvVars, err.Error())
		logger.Error(err.Error())
	}

	dbName, err := getEnvVar("DB_NAME")
	if err != nil {
		missingRequiredEnvVars = append(missingRequiredEnvVars, err.Error())
		logger.Error(err.Error())
	}

	instanceConnectionName, err := getEnvVar("INSTANCE_CONNECTION_NAME")
	if err != nil {
		missingRequiredEnvVars = append(missingRequiredEnvVars, err.Error())
		logger.Error(err.Error())
	}

	if len(missingRequiredEnvVars) > 0 {
		logger.Info("Please provide missing required environment variables and restart application")
		os.Exit(1)
	}

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)

	dbpool, err := openDbPool(dsn, instanceConnectionName)
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

func openDbPool(dsn, instanceConnectionName string) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	dialer, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		return nil, err
	}

	// Tell the driver to use the Cloud SQL Go Connector to create connections
	config.ConnConfig.DialFunc = func(ctx context.Context, _ string, instance string) (net.Conn, error) {
		return dialer.Dial(ctx, instanceConnectionName)
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		dbpool.Close()
		dialer.Close()
		return nil, err
	}

	return dbpool, nil
}
