package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_HealthCheck(t *testing.T) {
	t.Run("Should return 200 OK when calling healtcheck endpoint", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/health", strings.NewReader(""))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, HealthCheck(c)) {
			assert.Equal(t, 200, rec.Code)
			assert.Equal(t, "{\"environment\":\"dev\",\"message\":\"API is healthy\"}\n", rec.Body.String())
		}
	})
}
