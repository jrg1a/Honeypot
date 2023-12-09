package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error

	db, err = sql.Open("sqlite3", "./mydatabase.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	const tableCreationQuery = `CREATE TABLE IF NOT EXISTS logs (
        id INTEGER PRIMARY KEY,
        time TEXT,
        method TEXT,
        path TEXT,
        ip TEXT
    )`

	_, err = db.Exec(tableCreationQuery)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	srv := StartHTTPServer(db)
	StartSSHServer()
	StartFTPServer()

	// Graceful Shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
}
