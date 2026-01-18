package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed templates/*
var templatesFS embed.FS

type PageData struct {
	Message     string
	Environment string
	Version     string
	Hostname    string
	Timestamp   string
}

func main() {
	// Read environment variables
	port := getEnv("PORT", "8080")
	message := getEnv("DISPLAY_MESSAGE", "Welcome to Service Template!")
	environment := getEnv("ENVIRONMENT", "production")
	version := getEnv("VERSION", "1.0.0")

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Parse templates
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Handler for the main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Message:     message,
			Environment: environment,
			Version:     version,
			Hostname:    hostname,
			Timestamp:   time.Now().Format("2006-01-02 15:04:05 MST"),
		}
		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","environment":"%s","version":"%s"}`, environment, version)
	})

	// Ready check endpoint
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ready"}`)
	})

	// Start server
	log.Printf("Starting server on port %s", port)
	log.Printf("Environment: %s", environment)
	log.Printf("Version: %s", version)
	log.Printf("Message: %s", message)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
