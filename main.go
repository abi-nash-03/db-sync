package main

import (
	// "db-sync/cmd"
	"database/sql"
	"db-sync/config"
	"db-sync/tunnel"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// cmd.Execute()

	config.LoadConfig("config.yaml")

	tunnel := &tunnel.Tunnel{
		SSHHost:    config.AppConfig.Destination.SSH.Host,
		SSHPort:    config.AppConfig.Destination.SSH.Port,
		SSHUser:    config.AppConfig.Destination.SSH.User,
		KeyPath:    config.AppConfig.Source.SSH.KeyPath,
		RemotePort: config.AppConfig.Source.Port,
		LocalPort:  0,
	}

	if err := tunnel.Open(); err != nil {
		fmt.Println("Error opening tunnel:", err)
		return
	}

	defer tunnel.Close()
	fmt.Println("Tunnel opened on port:", tunnel.LocalPort)

	// --- 3. Connect MySQL through the tunnel ---
	// DSN format: user:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf(
		"%s:%s@tcp(127.0.0.1:%d)/%s?timeout=10s",
		config.AppConfig.Destination.User,     // e.g. "root"
		config.AppConfig.Destination.Password, // e.g. "secret"
		tunnel.LocalPort,
		config.AppConfig.Destination.Database, // e.g. "myapp" or leave blank as "information_schema" for testing
	)

	fmt.Println("DSN:", dsn)

	fmt.Println("Connecting to MySQL through tunnel...")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// --- 4. Run a simple test query ---
	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		log.Fatalf("query failed: %v", err)
	}
	fmt.Printf("✓ MySQL is alive! Version: %s\n", version)

	// --- 5. List databases (sanity check) ---
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		log.Fatalf("show databases failed: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nDatabases on remote server:")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Printf("  - %s\n", name)
	}

	fmt.Println("\n✓ All good! Tunnel + MySQL working correctly.")
}
