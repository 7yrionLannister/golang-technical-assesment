package util

import (
	"fmt"
	"strings"

	"github.com/7yrionLannister/golang-technical-assesment/config/logger"
)

// Logs the error and returns a new error with the message and the original error.
// The message is converted to lowercase, it is best practice for error strings.
func HandleError(err error, message string) error {
	e := fmt.Errorf("%s: %w", strings.ToLower(message), err)
	logger.L.Error(message, "error", err)
	return e
}
