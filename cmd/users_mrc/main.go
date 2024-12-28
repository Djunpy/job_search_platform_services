package main

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	db "job_search_platform/internal/users_mrc/db/sqlc"
	"job_search_platform/internal/users_mrc/server"
	config2 "job_search_platform/pkg/config"
	"job_search_platform/pkg/database"
	logger2 "job_search_platform/pkg/logger"
	"job_search_platform/pkg/scheduler"
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
	config, err := config2.LoadConfig("../..", "users_mrc")
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
	redisOpt := asynq.RedisClientOpt{Addr: config.RedisAddress}
	taskDistributor := scheduler.NewRedisTaskDistributor(redisOpt)
	waitGroup, ctx := errgroup.WithContext(ctx)
	err = scheduler.RunTaskProcessor(ctx, waitGroup, config, redisOpt, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot run task processor")
	}
	err = server.RunGinServer(config, store, taskDistributor, logger)
	err = waitGroup.Wait()
	if err != nil {
		logger.Fatal().Err(err).Msg("error from wait group")
	}
}
