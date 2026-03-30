package restore

import (
	"os"
	"testing"
	"time"
)

func TestDumpFileExists(t *testing.T) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	dumpPath := "/tmp/dump_" + timestamp + ".sql"
	_, err := os.Stat(dumpPath)
	if os.IsNotExist(err) {
		t.Errorf("dump file not found: %s", dumpPath)
	}
}
