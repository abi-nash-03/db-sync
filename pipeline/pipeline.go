package pipeline

import (
	"db-sync/config"
	"db-sync/dumper"
	"db-sync/restore"
	"fmt"
	"log/slog"
)

func Run(c *config.Config, dryRun bool) error {

	if dryRun {
		slog.Info("[dry-run] would dump source database",
			"host", c.Source.Host,
			"database", c.Source.Database,
		)
		slog.Info("[dry-run] would restore to destination",
			"host", c.Destination.Host,
			"database", c.Destination.Database,
		)
		return nil
	}

	// actual pipeline
	dumpPath, err := dumper.Dump(c)
	if err != nil {
		slog.Error("Dump failed: %s\n", "error", err)
		return fmt.Errorf("dump step failed: %w", err)
	}

	slog.Info("✓ Dump created successfully: %s\n", "info", dumpPath)

	_, err = restore.Restore(dumpPath, c)
	if err != nil {
		slog.Error("Restore failed: %s\n", "error", err)
		return fmt.Errorf("restore step failed: %w", err)
	}
	slog.Info("✓ Restore complete: %s is up to date\n", "info", c.Destination.Database)

	return nil
}
