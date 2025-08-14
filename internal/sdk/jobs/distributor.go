package jobs

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type TaskDistributor interface {
	DistributeTaskVerifyEmail(context.Context, *VerifyEmailPayload) error
}

type TaskQueue struct {
	client *asynq.Client
	logger *zerolog.Logger
}

func NewTaskQueue(redisOpt asynq.RedisClientOpt, logger *zerolog.Logger) TaskDistributor {
	redisClient := asynq.NewClient(redisOpt)
	return &TaskQueue{client: redisClient, logger: logger}
}

// taskcreation and distribution
func (jd *TaskQueue) DistributeTaskVerifyEmail(ctx context.Context, payload *VerifyEmailPayload) error {
	taskJob, err := NewVerifyEmailTask(payload.Username)
	if err != nil {
		jd.logger.Error().
			Err(err).
			Str("username", payload.Username).
			Str("task_type", TypeEmailVerify).
			Msg("failed to create email verification task")
		return fmt.Errorf("create email verification task: %w", err)
	}

	info, err := jd.client.Enqueue(taskJob)
	if err != nil {
		jd.logger.Error().
			Err(err).
			Str("queue", info.Queue).
			Str("username", payload.Username).
			Str("task_type", TypeEmailVerify).
			Msg("failed to enqueue email verification task")
		return fmt.Errorf("enqueue email verification task: %w", err)
	}
	jd.logger.Info().
		Str("username", payload.Username).
		Str("task_type", TypeEmailVerify).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("succesfylly enqueued email verification task")
	return nil
}
