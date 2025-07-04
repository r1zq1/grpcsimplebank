package main

import (
	"github.com/hibiken/asynq"
	"github.com/r1zq1/grpcsimplebank/config"
	"github.com/r1zq1/grpcsimplebank/server"
	"github.com/r1zq1/grpcsimplebank/worker"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

func main() {
	go server.StartGatewayServer()
	go server.StartGRPCServer()

	// Load configuration
	config, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config: %v")
	}

	// Setup Redis connection options
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress, // e.g., "localhost:6379"
	}

	// Create task processor
	processor := worker.NewRedisTaskProcessor(config)

	// Create Asynq server
	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  4,
		},
	})

	// Register task handler(s)
	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TaskSendWelcomeEmail, processor.ProcessTaskSendWelcomeEmail)

	log.Info().Msg("👷‍♂️ Starting Asynq worker server...")

	// Start processing tasks
	if err := srv.Run(mux); err != nil {
		log.Fatal().Err(err).Msg("could not run Asynq server: %v")
	}
	select {}
}
