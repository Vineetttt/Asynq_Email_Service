package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// HandleWelcomeEmailTask handler for welcome email task.
func HandleWelcomeEmailTask(c context.Context, t *asynq.Task) error {
	payload := make(map[string]interface{})
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	id, ok := payload["user_id"].(float64)
	if !ok {
		return fmt.Errorf("user_id field not found")
	}
	userID := int(id)
	fmt.Printf("Send Welcome Email to User ID %d\n", userID)

	return nil
}

// HandleReminderEmailTask for reminder email task.
func HandleReminderEmailTask(c context.Context, t *asynq.Task) error {
	type ReminderPayload struct {
		UserID int    `json:"user_id"`
		SentIn string `json:"sent_in"`
	}

	var payload ReminderPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	id := payload.UserID
	time := payload.SentIn

	fmt.Printf("Send Reminder Email to User ID %d\n", id)
	fmt.Printf("Reason: time is up (%v)\n", time)

	return nil
}
