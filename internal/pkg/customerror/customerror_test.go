package customerror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertNoData(t *testing.T, customErr CustomError) {
	assert.Equal(t, customErr.Error(), "Test Custom Error")
	errorData := customErr.GetData()
	assert.Equal(t, 0, errorData["statusCode"])
	assert.NotEqual(t, "", errorData["function"])
	assert.Equal(t, "", errorData["type"])
	assert.Nil(t, errorData["data"])
	assert.Equal(t, 500, customErr.GetStatusCode())
}

func assertWithData(t *testing.T, customErr CustomError) {
	assert.Equal(t, customErr.Error(), "Error with Data")
	errorData := customErr.GetData()
	assert.Equal(t, 400, errorData["statusCode"])
	assert.NotEqual(t, "", errorData["function"])
	assert.Equal(t, "BadRequest", errorData["type"])
	assert.Equal(t, map[string]interface{}{
		"info": "Wrong payload",
	}, errorData["data"])
	assert.Equal(t, 400, customErr.GetStatusCode())

}

func getCustomErrorData() InputData {
	return InputData{
		StatusCode: 400,
		ErrType:    "BadRequest",
		Data: map[string]interface{}{
			"info": "Wrong payload",
		},
	}
}
func Test_New(t *testing.T) {
	t.Run("Should create a custom error without custom Data", func(t *testing.T) {
		customErr := New("Test Custom Error")
		assertNoData(t, customErr)

	})

	t.Run("Should create a custom error with custom Data", func(t *testing.T) {
		customErr := New("Error with Data", getCustomErrorData())
		assertWithData(t, customErr)
	})
}

func Test_NewFromErrror(t *testing.T) {
	t.Run("Should create a custom error without custom Data", func(t *testing.T) {
		customErr := NewFromErrror(errors.New("Test Custom Error"))
		assertNoData(t, customErr)
	})

	t.Run("Should create a custom error with custom Data", func(t *testing.T) {
		customErr := NewFromErrror(errors.New("Error with Data"), getCustomErrorData())
		assertWithData(t, customErr)
	})
}
