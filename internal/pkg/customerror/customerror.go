package customerror

import (
	"github.com/go-errors/errors"
)

// CustomError provides a custom error interface
type CustomError interface {
	Error() string
	GetData() map[string]interface{}
	GetStatusCode() int
}

type customError struct {
	err        *errors.Error
	StatusCode int
	ErrType    string
	Function   string
	Data       map[string]interface{}
}

// InputData defines custom error input data
type InputData struct {
	StatusCode int
	ErrType    string
	Data       map[string]interface{}
}

func (c customError) Error() string {
	return c.err.Error()
}

func (c customError) GetData() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": c.StatusCode,
		"type":       c.ErrType,
		"function":   c.Function,
		"data":       c.Data,
	}
}

func (c customError) GetStatusCode() int {
	if c.StatusCode == 0 {
		return 500
	}

	return c.StatusCode
}

func (c *customError) setFunction() {
	c.Function = c.err.StackFrames()[1].String()
}

func (c *customError) setErrorData(errData []InputData) {
	if len(errData) > 0 {
		c.Data = errData[0].Data
		c.ErrType = errData[0].ErrType
		c.StatusCode = errData[0].StatusCode
	}
}

func createCustomError(err *errors.Error, errData []InputData) CustomError {
	customError := customError{err: err}
	customError.setFunction()
	customError.setErrorData(errData)

	return customError
}

// New creates a new custom error instance
func New(msg string, errData ...InputData) CustomError {
	e := errors.Errorf(msg)
	return createCustomError(e, errData)
}

// NewFromErrror creates a new custom error instance based on an error
func NewFromErrror(err error, errData ...InputData) CustomError {
	e := errors.New(err)
	return createCustomError(e, errData)
}
