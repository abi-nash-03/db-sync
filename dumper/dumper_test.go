package dumper

import (
	"regexp"
	"strings"
	"testing"
)

func TestGenerateDumpPath(t *testing.T) {
	path := generateDumpPath()

	// should start with /tmp/dump_
	if !strings.HasPrefix(path, "/tmp/dump_") {
		t.Errorf("expected path to start with '/tmp/dump_', got: %s", path)
	}

	// should end with .sql
	if !strings.HasSuffix(path, ".sql") {
		t.Errorf("expected path to end with '.sql', got: %s", path)
	}

	// full path should match the pattern: /tmp/dump_YYYY-MM-DD_HH-MM-SS.sql
	matched, err := regexp.MatchString(
		`^.*/dump_\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}\.sql$`,
		path,
	)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Errorf("path does not match expected format, got: %s", path)
	}
}
