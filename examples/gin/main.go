package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var applicationName string = "EXAMPLE-API"

func simulateDBCall() {
	time.Sleep(50 * time.Millisecond)
}

func routes(app *gin.Engine) {
	app.GET("/users", func(c *gin.Context) {
		simulateDBCall()

		c.JSON(200, gin.H{
			"message": "Users fetched correctly",
			"count":   15,
		})
	})

	app.GET("/", func(ctx *gin.Context) {
		traceID, spanID, isSampled := GetTraceInfo(ctx)
		fmt.Printf("traceID: %v; spanID: %v; isSampled: %v\n", traceID, spanID, isSampled)
	})
}

func GetTraceInfo(ctx context.Context) (traceID string, spanID string, isSampled bool) {
	spanCtx := trace.SpanContextFromContext(ctx)

	if spanCtx.HasTraceID() {
		traceID = spanCtx.TraceID().String()
	}
	if spanCtx.HasSpanID() {
		spanID = spanCtx.SpanID().String()
	}

	isSampled = spanCtx.IsSampled()

	return traceID, spanID, isSampled
}

func main() {

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
		),
	)

	app := gin.Default()
	app.ContextWithFallback = true
	app.Use(otelgin.Middleware(applicationName))

	routes(app)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app,
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
