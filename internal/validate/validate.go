// Package validate allows validation checks on user input
package validate

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldError represents a validation failure for a specific field.
type FieldError struct {
	Field   string
	Message string
}

func (e FieldError) Error() string {
	return e.Message
}

// Check is a validation function that returns nil if valid.
type Check func() *FieldError

// Fields runs all checks and returns any errors found.
func Fields(checks ...Check) []FieldError {
	var errs []FieldError
	for _, check := range checks {
		if err := check(); err != nil {
			errs = append(errs, *err)
		}
	}
	return errs
}

// Required checks that a string is non-empty after trimming whitespace.
func Required(value, field string) Check {
	return func() *FieldError {
		if strings.TrimSpace(value) == "" {
			return &FieldError{Field: field, Message: fmt.Sprintf("%s is required", field)}
		}
		return nil
	}
}

// MinLen checks that a string has at least min characters.
func MinLen(value string, min int, field string) Check {
	return func() *FieldError {
		if len(value) < min {
			return &FieldError{Field: field, Message: fmt.Sprintf("%s must be at least %d characters", field, min)}
		}
		return nil
	}
}

// MaxLen checks that a string has at most max characters.
func MaxLen(value string, max int, field string) Check {
	return func() *FieldError {
		if len(value) > max {
			return &FieldError{Field: field, Message: fmt.Sprintf("%s must be at most %d characters", field, max)}
		}
		return nil
	}
}

// Numeric checks that a string can be parsed as a number.
func Numeric(value, field string) Check {
	return func() *FieldError {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return &FieldError{Field: field, Message: fmt.Sprintf("%s must be a number", field)}
		}
		return nil
	}
}
