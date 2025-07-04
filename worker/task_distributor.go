package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeSendWelcomeEmail(ctx context.Context, payload PayloadSendEmail, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{client: client}
}

func (d *RedisTaskDistributor) DistributeSendWelcomeEmail(ctx context.Context, payload PayloadSendEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendWelcomeEmail, jsonPayload)

	opts = []asynq.Option{
		asynq.MaxRetry(3),               // Retry maksimal 5x
		asynq.Queue("default"),          // Kirim ke queue default
		asynq.Timeout(30 * time.Second), // Task timeout 30 detik
		asynq.Retention(24 * time.Hour), // Simpan task selesai/gagal 24 jam
		asynq.ProcessIn(0),              // Delay (bisa diganti kalau perlu)
	}

	info, err := d.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("cannot enqueue task: %w", err)
	}

	fmt.Printf("🎯 Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}
