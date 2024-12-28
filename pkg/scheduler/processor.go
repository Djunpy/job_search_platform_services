package scheduler

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"job_search_platform/pkg/config"
	"job_search_platform/pkg/mail_sender"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type RedisTaskProcessor struct {
	server *asynq.Server
	config config.Config
	mailer mail_sender.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, mailer mail_sender.EmailSender, config config.Config) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: logger,
		},
	)

	return &RedisTaskProcessor{
		server: server,
		mailer: mailer,
		config: config,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}

func RunTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config config.Config,
	redisOpt asynq.RedisClientOpt,
	logger zerolog.Logger,
) error {
	mailer := mail_sender.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := NewRedisTaskProcessor(redisOpt, mailer, config)

	logger.Info().Msg("start task processor")
	err := taskProcessor.Start()

	if err != nil {
		return err
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		logger.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		logger.Info().Msg("task processor is stopped")
		return nil
	})
	return nil
}
