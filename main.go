package main

import (
	"db-sync/config"
	"db-sync/dumper"
	"db-sync/restore"
	"fmt"
	"time"
)

func main() {

	startTime := time.Now()
	config.LoadConfig("config.yaml")
	fmt.Printf("Dump Started at %s\n", startTime.Format("2006-01-02 15:04:05"))

	dumpPath, err := dumper.Dump(config.AppConfig)
	if err != nil {
		fmt.Println("Error while generating dump %w", err)
		return
	}

	_, err = restore.Restore(dumpPath, config.AppConfig)
	if err != nil {
		fmt.Println("Error while restoring dump %w", err)
	}

	endTime := time.Now()

	fmt.Println("==============================================")
	fmt.Printf("Start time: %s\n", startTime)
	fmt.Printf("End time: %s\n", endTime)
	fmt.Printf("Total time Taken: %s\n", endTime.Sub(startTime))
	fmt.Println("==============================================")

}
