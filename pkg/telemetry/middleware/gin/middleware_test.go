package gin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewGinMiddlewareConfig(t *testing.T) {
	mdw := NewGinMiddlewareConfig()
	assert.NotNil(t, mdw)
}

func TestMakeFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test-path?q=1", nil)
	req.RemoteAddr = "127.0.0.1:5000"
	req.Header.Set("User-Agent", "Test-Agent")
	c.Request = req

	// Simulate db fields and errors
	c.Set("num_db_calls", 5)
	c.Set("num_db_failures", 1)
	c.Error(errors.New("simulated error"))
	c.Writer.WriteHeader(500)

	fields := makeFields(c, 150)

	assert.NotNil(t, fields)

	reqFields := fields["request"]
	assert.NotNil(t, reqFields)

	dbFields := fields["database"]
	assert.NotNil(t, dbFields)

	errFields := fields["errors"]
	assert.NotNil(t, errFields)

	assert.Equal(t, int64(150), fields["duration_ms"])
}

func TestLoggerMiddleware_Execution(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mdw := NewGinMiddlewareConfig()
	handler := mdw.Middleware()

	assert.Len(t, handler, 2) // Should contain otelgin and loggerMiddleware

	// Execute just the logger middleware to ensure it doesn't panic
	r := gin.New()
	r.Use(handler[1])
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/bad", func(c *gin.Context) {
		c.JSON(400, gin.H{"error": "bad request"})
	})
	r.GET("/err", func(c *gin.Context) {
		c.JSON(500, gin.H{"error": "server error"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w400 := httptest.NewRecorder()
	req400, _ := http.NewRequest("GET", "/bad", nil)
	r.ServeHTTP(w400, req400)
	assert.Equal(t, 400, w400.Code)

	w500 := httptest.NewRecorder()
	req500, _ := http.NewRequest("GET", "/err", nil)
	r.ServeHTTP(w500, req500)
	assert.Equal(t, 500, w500.Code)
}
