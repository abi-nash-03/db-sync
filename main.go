package main

import (
	"db-sync/cmd"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

// Version is set at build time via -ldflags "-X main.Version=<tag>"
var version = "dev"

func main() {

	main_logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	logger := main_logger.With(
		"run_id",
		uuid.New().String(),
	)
	slog.SetDefault(logger)


	cmd.SetVersion(version)
	cmd.Execute()
}
