package pipeline

import (
	"db-sync/config"
	"testing"
)

func TestDryRunMode(t *testing.T) {
	c := &config.Config{
		Source: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "password",
			Database: "test",
		},
		Destination: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "password",
			Database: "test",
		},
	}

	if err := Run(c, true); err != nil {
		t.Errorf("Error running pipeline: %s", err)
	}
}
