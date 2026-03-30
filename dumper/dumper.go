package dumper

import (
	"bytes"
	"db-sync/config"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

/*
Generate the dump, using mysqldump command

path :
	/tmp/dump_<timestamp>.sql
*/

// generateDumpPath returns the file path for a dump file based on the current timestamp.
func generateDumpPath() string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return "/tmp/dump_" + timestamp + ".sql"
}

func Dump(c *config.Config) (string, error) {
	fmt.Println("Dumping database...")
	path := generateDumpPath()

	// check mysqldump exists before even trying
	_, err := exec.LookPath("mysqldump")
	if err != nil {
		return "", fmt.Errorf("mysqldump not found in PATH — please install mysql-client")
	}

	// execute the mysqldump command
	cmd := exec.Command(
		"mysqldump",
		"-u", c.Source.User,
		"-h", c.Source.Host,
		"-P", strconv.Itoa(c.Source.Port),
		"--single-transaction", // safe for InnoDB, no table locks
		"--routines",           // include stored procedures and functions
		"--triggers",           // include triggers
		"--events",             // include scheduled events
		"--add-drop-database",  // ← adds DROP DATABASE IF EXISTS
		"--add-drop-table",     // ← adds DROP TABLE IF EXISTS
		"--databases",
		c.Source.Database)

	// this will add the DROP "DATABASE IF EXISTS" and "DROP TABLE IF EXISTS" command in the dump file
	// so that when we restore the database, it will drop the database and tables first and then restore
	// we dont neeed to handle that in restore function

	// adding the password to env so that command won't prompt for password
	cmd.Env = append(os.Environ(),
		"MYSQL_PWD="+c.Source.Password,
	)

	// create the dump file
	outFile, err := os.Create(path)
	if err != nil {
		os.Remove(path) // clean up the partial file
		return "", fmt.Errorf("failed to create dump file: %w", err)
	}
	defer outFile.Close()

	// redirect mysqldump output to the file
	// incase of large dump, we dont need to store in-memory
	var stderr bytes.Buffer
	cmd.Stdout = outFile
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run mysqldump: %w\n%s", err, stderr.String())
	}

	fmt.Printf("✓ Dump created successfully: %s\n", path)

	return path, nil
}
