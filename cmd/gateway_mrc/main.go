package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	db "job_search_platform/internal/gateway_mrc/db/sqlc"
	"job_search_platform/internal/gateway_mrc/server"
	config2 "job_search_platform/pkg/config"
	"job_search_platform/pkg/database"
	logger2 "job_search_platform/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	var configPath string
	if _, err := os.Stat("/app"); err == nil {
		configPath = "/app"
	} else {
		configPath = "../.."
	}
	config, err := config2.LoadConfig(configPath, "gateway_mrc")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	logger, file := logger2.ConfigureLogger(config.Environment)
	defer file.Close()
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	connPool, err := pgxpool.New(ctx, config.PostgresSource)
	defer connPool.Close()
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to db")
	}
	err = connPool.Ping(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to db")
	}
	err = database.RunDBMigration(config.MigrationURL, config.PostgresSource)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot run migration")
	}
	store := db.NewStore(connPool)
	err = server.RunGinServer(config, store, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot run gin server")
	}
}
