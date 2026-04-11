package notifier

import (
	"bytes"
	"db-sync/config"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Notify(success bool, message string) error {

	payload := map[string]string{
		"text": message,
	}
	body, _ := json.Marshal(payload)

	response, err := http.Post(config.AppConfig.Notify.SlackWebhook, "application/json", bytes.NewBuffer(body))

	if err != nil {
		slog.Error("Error sending notification", "error", err)
		return fmt.Errorf("slack notification failed: %w", err)
	}

	defer response.Body.Close()

	return nil
}

func NotifyWithRetry(success bool, message string, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = Notify(success, message); err == nil {
			return nil
		}

		slog.Warn("notification failed, retrying",
			"attempt", i+1,
			"error", err,
		)

		time.Sleep(time.Second * 2)
	}
	return fmt.Errorf("slack notification failed after %d retries", maxRetries)
}
