package restore

import (
	"strings"
	"testing"
)

func TestRestoreFailsIfDumpFileNotFound(t *testing.T) {
	// this file definitely doesn't exist
	fakePath := "/tmp/this-file-does-not-exist-12345.sql"

	_, err := Restore(fakePath, nil) // pass nil config - shouldn't reach that far
	if err == nil {
		t.Fatal("expected an error for missing dump file, got nil")
	}

	if !strings.Contains(err.Error(), "dump file not found") {
		t.Errorf("expected 'dump file not found' in error, got: %s", err.Error())
	}
}
