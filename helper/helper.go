package helper

import (
	"errors"
	"fmt"
	"os"
)

func validateKeyFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("key file does not exist: %s", path)
		}
		return fmt.Errorf("cannot access key file: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory, not a key file", path)
	}
	if info.Mode().Perm()&0077 != 0 {
		return fmt.Errorf("key file permissions too open — run: chmod 600 %s", path)
	}
	return nil
}
