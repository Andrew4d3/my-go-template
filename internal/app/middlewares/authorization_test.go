package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"template-go-api/configs"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func getTestContext() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func stubController(c echo.Context) error {
	return c.JSON(200, "OK")
}

var testJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtc2ciOiJ0ZXN0In0.aMwCwG_SaTYTfzI6WVtnN-LlpLx-IfKLpfjBO6OFcEM"

func Test_AuthorizationMiddleware(t *testing.T) {
	ogGetJWTSecret := getJWTSecret
	defer func() {
		getJWTSecret = ogGetJWTSecret
	}()

	getJWTSecret = func() configs.AppConfig {
		return configs.AppConfig{
			JWTSecret: "mySecret",
		}
	}

	t.Run("Should return 401 if authorization header is empty", func(t *testing.T) {
		c, _ := getTestContext()
		err := AuthorizationMiddleware(stubController)(c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "code=401")
	})

	t.Run("Should continue if JWT is valid", func(t *testing.T) {
		c, rec := getTestContext()
		c.Request().Header.Set("authorization", fmt.Sprintf("Bearer %s", testJWT))
		err := AuthorizationMiddleware(stubController)(c)
		if assert.NoError(t, err) {
			assert.Equal(t, 200, rec.Code)
		}
	})

	t.Run("Should throw 401 if JWT is not valid", func(t *testing.T) {
		c, _ := getTestContext()
		c.Request().Header.Set("authorization", "Bearer BadToken")
		err := AuthorizationMiddleware(stubController)(c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "code=401")
	})

}
