package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	otelcol "github.com/alanfranciscos/otel-collector/pkg/telemetry"
	ginmdw "github.com/alanfranciscos/otel-collector/pkg/telemetry/middleware/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var applicationName string = "EXAMPLE-API"
var tracer = otel.Tracer("example-handler")
var meter = otel.Meter("example-meter")
var requestCounter, _ = meter.Int64Counter("api_requests_total", metric.WithDescription("Total number of API requests"))

func simulateDBCall(ctx context.Context) {
	_, span := tracer.Start(ctx, "database.query.simulate")
	defer span.End()
	time.Sleep(50 * time.Millisecond)
}

func routes(app *gin.Engine) {
	app.GET("/users", func(c *gin.Context) {
		requestCounter.Add(c.Request.Context(), 1)
		simulateDBCall(c.Request.Context())

		c.JSON(200, gin.H{
			"message": "Users fetched correctly",
			"count":   15,
		})
	})

	app.GET("/", func(ctx *gin.Context) {
		requestCounter.Add(ctx.Request.Context(), 1)
		ctx.Error(errors.New("fake error (simulated vem antes)"))
		simulateDBCall(ctx.Request.Context())
		ctx.Error(errors.New("fake error (simulated)"))
		logrus.WithContext(ctx.Request.Context()).Info(
			"Info test",
		)
		ctx.JSON(500, gin.H{
			"message": "Hello, World!",
		})
	})
}

func setEnvironment() {
	os.Setenv("OTEL_SERVICE_NAME", applicationName)
	os.Setenv("ENVIRONMENT", "local")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "http")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
}

func main() {
	setEnvironment()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := otelcol.NewTelemetry(&applicationName).Initialize(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	app := gin.New()

	ginMdw := ginmdw.NewGinMiddlewareConfig()
	app.Use(ginMdw.Middleware()...)

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

	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Shutdown complete.")
}
