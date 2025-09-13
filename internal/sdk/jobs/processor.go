package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/0xOnah/bank/internal/db/repo"
	"github.com/0xOnah/bank/internal/entity"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	JobSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type UserStore interface {
	GetUser(ctx context.Context, username string) (*entity.User, error)
}

type WorkerService struct {
	server    *asynq.Server
	userStore UserStore
	logger    *zerolog.Logger
}

func NewWorkerService(redisOpt asynq.RedisClientOpt, usStore UserStore, logger *zerolog.Logger) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
		},
	)
	return &WorkerService{server: server, userStore: usStore, logger: logger}
}

func (rt *WorkerService) JobSendVerifyEmail(ctx context.Context, t *asynq.Task) error {
	var payload VerifyEmailPayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		rt.logger.Error().
			Err(err).
			Str("username", payload.Username).
			Msg("JobSendVerifyEmail: failed to unmarshal payload")
		return fmt.Errorf("bad payload: %w", asynq.SkipRetry)
	}

	user, err := rt.userStore.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotExist) {
			rt.logger.Warn().
				Err(err).
				Str("username", payload.Username).
				Msg("JobSendVerifyEmail: user does not exist")
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		rt.logger.Error().
			Err(err).Str("username", payload.Username).
			Msg("failed to get user")
		return fmt.Errorf("get user: %w, %w", err, asynq.SkipRetry)
	}

	//Todo: send emai to user

	_ = user
	rt.logger.Info().
		Str("type", t.Type()).
		Str("to", user.Email.String()).
		Msg("JobSendVerifyEmail: successfully sent verification email")
	return nil
}

func (rt *WorkerService) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailVerify, rt.JobSendVerifyEmail)

	return rt.server.Run(mux)
}
