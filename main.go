package main

import (
	"db-sync/cmd"
	"fmt"
	"log/slog"
	"os"
	"time"
)

func main() {

	//set this up once at startup
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	startTime := time.Now()

	cmd.Execute()

	endTime := time.Now()

	fmt.Println("==============================================")
	fmt.Printf("Start time: %s\n", startTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("End time: %s\n", endTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Total time Taken: %s\n", endTime.Sub(startTime))
	fmt.Println("==============================================")

}
