package dumper

import (
	"db-sync/config"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func Dump(c *config.Config) (string, error) {
	fmt.Println("Dumping database...")
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	file_name := "dump_" + timestamp + ".sql"
	path := "/tmp/" + file_name

	// execute the mysqldump command
	cmd := exec.Command(
		"mysqldump",
		"-u", c.Destination.User,
		"-h", c.Destination.Host,
		"-P", strconv.Itoa(c.Destination.Port),
		c.Destination.Database)

	// adding the password to env so that command wont prompt for password
	cmd.Env = append(os.Environ(),
		"MYSQL_PWD="+c.Destination.Password,
	)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// write the output to a file
	if err := os.WriteFile(path, out, 0644); err != nil {
		return "", err
	}

	return path, nil
}
