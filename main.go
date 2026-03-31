package main

import (
	"db-sync/config"
	"db-sync/dumper"
	"db-sync/restore"
	"fmt"
)

func main() {

	config.LoadConfig("config.yaml")

	dumpPath, err := dumper.Dump(config.AppConfig)
	if err != nil {
		fmt.Println("Error while generating dump %w", err)
		return
	}

	_, err = restore.Restore(dumpPath, config.AppConfig)
	if err != nil {
		fmt.Println("Error while restoring dump %w", err)
	}
}
