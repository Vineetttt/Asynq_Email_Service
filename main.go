package main

import (
	"log"
	"net/http"
	"time"

	"asynqmon/tasks"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

type TaskRequest struct {
	TaskType string `json:"task_type"`
	UserID   int    `json:"user_id"`
}

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})

	router := gin.Default()
	router.POST("/tasks", func(c *gin.Context) {
		var request TaskRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		var taskToEnqueue *asynq.Task
		var queueName string

		switch request.TaskType {
		case "welcome-email":
			taskToEnqueue = tasks.NewWelcomeEmailTask(request.UserID)
			queueName = "critical"
		case "reminder-email":
			taskToEnqueue = tasks.NewReminderEmailTask(request.UserID, time.Now().Add(2*time.Minute))
			queueName = "low"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task type"})
			return
		}

		info, err := client.Enqueue(
			taskToEnqueue,
			asynq.Queue(queueName),
			asynq.ProcessIn(2*time.Minute),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue task"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"info": info})
	})

	h := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitoring", // RootPath specifies the root for asynqmon app
		RedisConnOpt: asynq.RedisClientOpt{Addr: "localhost:6379"},
	})

	mux := http.NewServeMux()
	mux.Handle(h.RootPath()+"/", h)
	mux.Handle("/", router)

	log.Fatal(http.ListenAndServe(":4000", mux))
}
