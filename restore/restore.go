package restore

import (
	"bytes"
	"db-sync/config"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func Restore(dumpPath string, c *config.Config) (string, error) {
	fmt.Println("Restoring database...")

	_, error := os.Stat(dumpPath)
	if os.IsNotExist(error) {
		return "", fmt.Errorf("dump file not found: %s and error : %w", dumpPath, error)
	}

	// build mysql restore command
	cmd := exec.Command(
		"mysql",
		"-h", c.Destination.Host,
		"-u", c.Destination.User,
		"-P", strconv.Itoa(c.Destination.Port),
		c.Destination.Database,
	)

	cmd.Env = append(os.Environ(),
		"MYSQL_PWD="+c.Destination.Password,
	)

	// Piping the dump file as stdin
	dumpFile, err := os.Open(dumpPath)
	if err != nil {
		return "", fmt.Errorf("Error reading the dump file %w", err)
	}
	defer dumpFile.Close()

	// point cmd.Stdin at the file
	// mysql will read from it as if you ran: mysql < dump.sql
	cmd.Stdin = dumpFile

	// capture stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run mysql restore: %w\n%s", err, stderr.String())
	}

	fmt.Printf("✓ Restore complete: %s is up to date\n", c.Destination.Database)

	if err := os.Remove(dumpPath); err != nil {
		fmt.Printf("Error removing dump file: %s\n", err)
	}

	return "", nil
}
