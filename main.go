package main

import (
	"fmt"
	"time"
	"db-sync/cmd"
)




func main() {
    fmt.Println("Welcome to db-sync")
    cmd.Execute()
    time.Sleep(10 * time.Second)
}
