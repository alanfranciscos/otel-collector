package main

import (
	"log"
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

func startServer() {
	engine := gin.New()
	routes(engine)

	go func() {
		log.Println("Starting example API on :8080...")
		if err := engine.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}

func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutdown complete.")
}

func main() {
	startServer()
	gracefulShutdown()
}
