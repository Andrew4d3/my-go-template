package customerror

import (
	"github.com/go-errors/errors"
)

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

type CustomErrorData struct {
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

func (err *customError) setErrorData(errData []CustomErrorData) {
	if len(errData) > 0 {
		err.Data = errData[0].Data
		err.ErrType = errData[0].ErrType
		err.StatusCode = errData[0].StatusCode
	}
}

func createCustomError(err *errors.Error) CustomError {
	customError := customError{err: err}
	customError.setFunction()

	return customError
}

func New(msg string, errData ...CustomErrorData) CustomError {
	e := errors.Errorf(msg)
	return createCustomError(e)
}

func NewFromErrror(err error, errData ...CustomErrorData) CustomError {
	e := errors.New(err)
	return createCustomError(e)
}
