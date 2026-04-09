package notifier

import (
	"bytes"
	"db-sync/config"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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
