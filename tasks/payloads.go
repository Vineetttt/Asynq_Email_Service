package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeWelcomeEmail  = "email:welcome"
	TypeReminderEmail = "email:reminder"
)

// NewWelcomeEmailTask task payload for a new welcome email.
func NewWelcomeEmailTask(id int) *asynq.Task {
	payload := map[string]interface{}{
		"user_id": id,
	}

	// Convert payload to a byte slice.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return asynq.NewTask(TypeWelcomeEmail, payloadBytes)
}

// NewReminderEmailTask task payload for a reminder email.
func NewReminderEmailTask(id int, ts time.Time) *asynq.Task {
	payload := map[string]interface{}{
		"user_id": id,
		"sent_in": ts.String(),
	}

	// Convert payload to a byte slice.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return asynq.NewTask(TypeReminderEmail, payloadBytes)
}
