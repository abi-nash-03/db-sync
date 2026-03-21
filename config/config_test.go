package config

import (
	"os"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {
	// write a temp config file
	content := `
source:
  host: "localhost"
  port: 3306
  user: "root"
  password: "secret"
  database: "mydb"
destination:
  host: "remotehost"
  port: 3306
  user: "dev"
  password: "devpass"
  database: "devdb"
`
	tmpFile, _ := os.CreateTemp("", "config-*.yaml")
	tmpFile.WriteString(content)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Source.Host != "localhost" {
		t.Errorf("expected localhost, got %s", cfg.Source.Host)
	}
}

func TestValidateMissingField(t *testing.T) {
	cfg := &Config{} // empty config
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
