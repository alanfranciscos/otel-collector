package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func simulateDBCall() {
	time.Sleep(50 * time.Millisecond)
}

func routes(engine *gin.Engine) {
	engine.GET("/api/v1/users", func(c *gin.Context) {
		simulateDBCall()

		c.JSON(200, gin.H{
			"message": "Users fetched correctly",
			"count":   15,
		})
	})
}

func main() {
	engine := gin.New()
	routes(engine)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		log.Println("Starting example API on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// shutdown gracefully on interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Shutdown complete.")
}
