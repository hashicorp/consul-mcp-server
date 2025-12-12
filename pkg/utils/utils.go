// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

// GetEnv retrieves the value of an environment variable or returns a fallback value if not set
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// LogAndReturnError logs the error with context and returns a formatted error.
func LogAndReturnError(logger *log.Logger, context string, err error) error {
	err = fmt.Errorf("%s, %w", context, err)
	if logger != nil {
		logger.Errorf("Error in %s, %v", context, err)
	}
	return err
}
