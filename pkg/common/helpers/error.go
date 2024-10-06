package helpers

import (
	"errors"
	"fmt"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIfNil(obj interface{}) {
	if obj != nil {
		panic("")
	}
}

// CustomError wraps an error with a custom message.
func CustomError(message string, args ...interface{}) error {
	return fmt.Errorf(message, args...)
}

// IsNotFoundError checks if an error is a "not found" error.
func IsNotFoundError(err error) bool {
	return errors.Is(err, errors.New("not found"))
}
