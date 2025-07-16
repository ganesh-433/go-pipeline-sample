package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	log.Printf("Received request for %s\n", name)
	w.Write([]byte(CreateGreeting(name)))
}

func CreateGreeting(name string) string {
	if name == "" {
		name = "Guest"
	}
	return "Hello, " + name + "\n"
}

func main() {
	// Create Server and Route Handlers
	r := mux.NewRouter()
	r.HandleFunc("/", handler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080", // Listen on all interfaces on port 8080
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Create a channel to listen for OS signals (e.g., Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start Server in a goroutine
	go func() {
		log.Println("Starting Server on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log fatal error if server fails to start (e.g., port in use)
			log.Fatalf("Could not listen on %s: %v\n", srv.Addr, err)
		}
	}()

	// Block until an OS signal is received
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
	// os.Exit(0) is not strictly needed here as the main goroutine will exit
	// once srv.Shutdown completes and no other goroutines are running.
	// However, if you have other goroutines, os.Exit(0) might be necessary
	// to ensure the process truly terminates. For this simple app, it's fine.
}
