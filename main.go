package main

import (
	"db-sync/cmd"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {

	//set this up once at startup
	main_logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	logger := main_logger.With(
		"run_id",
		uuid.New().String(),
	)
	slog.SetDefault(logger)

	startTime := time.Now()

	cmd.Execute()

	endTime := time.Now()

	slog.Info("==============================================")
	slog.Info("Start: %s", "Time", startTime.Format("2006-01-02 15:04:05"))
	slog.Info("End: %s", "Time", endTime.Format("2006-01-02 15:04:05"))
	slog.Info("Total %s", "Time Taken:", endTime.Sub(startTime))
	slog.Info("==============================================")

}
