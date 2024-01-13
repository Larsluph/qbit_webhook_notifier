package helpers

import (
	"fmt"
	"runtime"
)

type ErrorPayload struct {
	CodeLocation string
	Description  string
}

func NewErrorPayload(err error) *ErrorPayload {
	file, line := getErrorLocation()

	o := new(ErrorPayload)
	o.CodeLocation = fmt.Sprintf("%s | L%d", file, line)
	o.Description = err.Error()
	return o
}

func getErrorLocation() (string, int) {
	_, filePath, line, _ := runtime.Caller(2)
	fmt.Printf("%s:%d", filePath, line)
	return filePath, line
}
