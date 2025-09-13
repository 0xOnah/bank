package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const TypeEmailVerify = "task:verify_email"

type VerifyEmailPayload struct {
	Username string
}

func TaskVerifyEmail(username string) (*asynq.Task, error) {
	payload, err := json.Marshal(VerifyEmailPayload{Username: username})
	if err != nil {
		return nil, fmt.Errorf("failed to marshall payload %w", err)
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(QueueCritical),
	}
	return asynq.NewTask(TypeEmailVerify, payload, opts...), nil
}
