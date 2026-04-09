package scheduler

import (
	"db-sync/config"
	"db-sync/notifier"
	"db-sync/pipeline"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

// slogAdapter bridges slog with the cron.Logger interface
type slogAdapter struct{}

func (s *slogAdapter) Info(msg string, keysAndValues ...interface{}) {
	slog.Info(msg, keysAndValues...)
}

func (s *slogAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	slog.Error(msg, append([]interface{}{"error", err}, keysAndValues...)...)
}

func Start(cfg *config.Config, schedule string) error {
	c := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(&slogAdapter{}),
	))

	c.AddFunc(schedule, func() {
		slog.Info("Scheduled run starting")

		if config.AppConfig.Notify.SlackWebhook != "" {
			if err := notifier.Notify(true, "Scheduled Database sync started"); err != nil {
				slog.Error("Error sending notification", "error", err)
			}
		}

		if err := pipeline.Run(cfg, false); err != nil {
			slog.Error("Scheduled Database sync failed", "error", err)
			if config.AppConfig.Notify.SlackWebhook != "" {
				if err := notifier.Notify(false, err.Error()); err != nil {
					slog.Error("Error sending notification", "error", err)
				}
			}
			return
		}
		slog.Info("Scheduled Database sync completed successfully")
		if config.AppConfig.Notify.SlackWebhook != "" {
			if err := notifier.Notify(true, "Scheduled Database sync completed successfully"); err != nil {
				slog.Error("Error sending notification", "error", err)
			}
		}
	})

	c.Start()

	// block forever - wait for Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // blocks here until signal received
	slog.Info("shutting down scheduler")
	c.Stop()

	return nil
}
