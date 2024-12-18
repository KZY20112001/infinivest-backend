package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/app"
	"github.com/KZY20112001/infinivest-backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	app.LoadEnv()
	app.SetupDB()
	app.InjectDependencies()
}

func main() {
	r := gin.Default()
	routes.RegisterRoutes(r)

	// Channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s\n", err)
		}
	}()

	// Block until a signal is received
	<-quit
	log.Println("Shutting down server...")

	// Gracefully shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shut down: %s\n", err)
	}

	log.Println("Disconnecting from DB...")
	app.CloseDB()

	log.Println("Server exited gracefully")
}
