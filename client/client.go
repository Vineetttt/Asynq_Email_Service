package main

import (
	"log"
	"math/rand"
	"time"

	"asynqmon/tasks"

	"github.com/hibiken/asynq"
)

func main() {
	// Create a new Redis connection for the client.
	redisConnection := asynq.RedisClientOpt{
		Addr: "localhost:6379", // Redis server address
	}

	// Create a new Asynq client.
	client := asynq.NewClient(redisConnection)
	defer client.Close()

	for i := 1; i < 10; i++ {
		userID := rand.Intn(1000) + 10
		delay := 2 * time.Minute

		task1 := tasks.NewWelcomeEmailTask(userID)
		task2 := tasks.NewReminderEmailTask(userID, time.Now().Add(delay))

		// Process the task immediately in critical queue.
		if _, err := client.Enqueue(
			task1,
			asynq.Queue("critical"),
		); err != nil {
			log.Fatal(err)
		}

		// Process the task 2 minutes later in low queue.
		if _, err := client.Enqueue(
			task2,
			asynq.Queue("low"),
			asynq.ProcessIn(delay),
		); err != nil {
			log.Fatal(err)
		}
	}
}
