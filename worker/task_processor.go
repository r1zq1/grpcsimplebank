package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	"github.com/hibiken/asynq"
	"github.com/jordan-wright/email"
	"github.com/r1zq1/grpcsimplebank/config"
)

const TaskSendWelcomeEmail = "task:send_welcome_email"

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type PayloadSendEmail struct {
	AccountID int64  `json:"account_id"`
	Email     string `json:"email"`
	Owner     string `json:"owner"`
}

type RedisTaskProcessor struct {
	config config.Config
}

func NewRedisTaskProcessor(config config.Config) *RedisTaskProcessor {
	return &RedisTaskProcessor{config: config}
}

func (p *RedisTaskProcessor) ProcessTaskSendWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	var payload PayloadSendEmail
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	e := email.NewEmail()
	e.From = p.config.SMTPFrom
	e.To = []string{payload.Email}
	e.Subject = "Welcome to GRPC Bank!"
	e.Text = []byte(fmt.Sprintf("Hello %s,\n\nYour account has been created successfully.\n\nRegards,\nGRPC Bank", payload.Owner))

	addr := fmt.Sprintf("%s:%d", p.config.SMTPHost, p.config.SMTPPort)
	auth := smtp.PlainAuth("", p.config.SMTPUsername, p.config.SMTPPassword, p.config.SMTPHost)

	log.Printf("📧 Sending welcome email to %s (%s)...", payload.Owner, payload.Email)
	if err := e.Send(addr, auth); err != nil {
		log.Printf("❌ Failed to send email: %v", err)
		return err
	}

	log.Printf("📧 Using From: %s", p.config.SMTPFrom)
	log.Printf("✅ Email sent to %s", payload.Email)
	return nil
}
