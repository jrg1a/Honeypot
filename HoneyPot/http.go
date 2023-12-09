package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// RequestData structure to capture more information from each request
type RequestData struct {
	Time      string
	Method    string
	Path      string
	IP        string
	UserAgent string
	Referrer  string
	Query     string
}

func RequestLogger(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	currentTime := time.Now().Format(time.RFC3339)
	data := RequestData{
		Time:      currentTime,
		Method:    req.Method,
		Path:      req.URL.Path,
		IP:        req.RemoteAddr,
		UserAgent: req.UserAgent(),
		Referrer:  req.Referer(),
		Query:     req.URL.RawQuery,
	}

	query := "INSERT INTO logs (time, method, path, ip, user_agent, referrer, query) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, data.Time, data.Method, data.Path, data.IP, data.UserAgent, data.Referrer, data.Query)

	if err != nil {
		log.Printf("Error inserting log: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(data)
	fmt.Fprintf(w, string(response))
}

func StartHTTPServer(db *sql.DB) *http.Server {
	// Rate limiting setup
	limiter := rate.NewLimiter(rate.Every(time.Second), 10)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		RequestLogger(w, r, db)
	})

	// Simulated endpoints
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		fmt.Fprintf(w, "Fake login page")
	})
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		fmt.Fprintf(w, "Fake admin page")
	})
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		fmt.Fprintf(w, "Fake API data")
	})

	// Custom error handling
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		http.Error(w, "Custom error message", http.StatusInternalServerError)
	})

	// HTTP Server setup
	port := "8080"
	srv := &http.Server{Addr: ":" + port}

	go func() {
		log.Printf("Starting HTTP honeypot server on port %s...\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	return srv
}
