package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	otelcol "github.com/alanfranciscos/otel-collector/pkg/telemetry"
	ginmdw "github.com/alanfranciscos/otel-collector/pkg/telemetry/middleware/gin"
)

func TestIntegration_GinRoutes_LogOutput(t *testing.T) {
	// Setup standard test environment
	setEnvironment()

	// Initialize Telemetry
	ctx := context.Background()
	shutdown, err := otelcol.NewTelemetry(&applicationName).Initialize(ctx)
	assert.NoError(t, err)
	defer shutdown(ctx)

	// Intercept logrus output to a buffer
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	defer logrus.SetOutput(nil) // Reset after test

	// Setup Gin app
	gin.SetMode(gin.TestMode)
	app := gin.New()

	ginMdw := ginmdw.NewGinMiddlewareConfig()
	app.Use(ginMdw.Middleware()...)

	routes(app)

	// Perform request to /users
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	app.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Users fetched correctly")

	// Read intercepted logs
	logOutput := buf.String()
	assert.NotEmpty(t, logOutput)

	// We expect multiple JSON lines (Request Started, Request Completed, etc.)
	lines := strings.Split(strings.TrimSpace(logOutput), "\n")
	foundCompletionLog := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}

		var logData map[string]interface{}
		err := json.Unmarshal([]byte(line), &logData)
		if err != nil {
			continue // Skip non-json lines that happen to start with {
		}

		// Check if it's the request completion log from the middleware
		if msg, ok := logData["message"].(string); ok && msg == "Request Completed" {
			foundCompletionLog = true

			// Assert core fields are present and valid
			assert.Equal(t, applicationName, logData["service_name"])
			assert.NotNil(t, logData["trace_id"])
			assert.NotNil(t, logData["span_id"])
			assert.NotEqual(t, "00000000000000000000000000000000", logData["trace_id"])
			assert.NotEqual(t, "0000000000000000", logData["span_id"])

			// Assert nested HTTP request fields
			reqData, ok := logData["request"].(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, "GET", reqData["method"])
			assert.Equal(t, "/users", reqData["path"])

			// Assert duration is present
			_, ok = logData["duration_ms"].(float64) // JSON unmarshals numbers as float64
			assert.True(t, ok)
		}
	}

	assert.True(t, foundCompletionLog, "Expected to find 'Request Completed' log entry")
}
