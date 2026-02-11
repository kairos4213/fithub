package validate

import (
	"strings"
	"testing"
)

func TestRequired(t *testing.T) {
	tests := map[string]struct {
		value   string
		wantErr bool
	}{
		"non-empty value":    {value: "hello", wantErr: false},
		"empty string":       {value: "", wantErr: true},
		"whitespace only":    {value: "   ", wantErr: true},
		"tabs and newlines":  {value: "\t\n", wantErr: true},
		"value with padding": {value: "  hello  ", wantErr: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			check := Required(tc.value, "test_field")
			err := check()
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
			if err != nil && err.Field != "test_field" {
				t.Errorf("expected field: test_field, got: %s", err.Field)
			}
		})
	}
}

func TestMinLen(t *testing.T) {
	tests := map[string]struct {
		value   string
		min     int
		wantErr bool
	}{
		"meets minimum":   {value: "abcdefghij", min: 10, wantErr: false},
		"exceeds minimum": {value: "abcdefghijk", min: 10, wantErr: false},
		"below minimum":   {value: "abc", min: 10, wantErr: true},
		"empty string":    {value: "", min: 1, wantErr: true},
		"exact boundary":  {value: "ab", min: 2, wantErr: false},
		"one below":       {value: "a", min: 2, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			check := MinLen(tc.value, tc.min, "password")
			err := check()
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestMaxLen(t *testing.T) {
	tests := map[string]struct {
		value   string
		max     int
		wantErr bool
	}{
		"under maximum":  {value: "abc", max: 10, wantErr: false},
		"at maximum":     {value: "abcde", max: 5, wantErr: false},
		"over maximum":   {value: "abcdef", max: 5, wantErr: true},
		"empty string":   {value: "", max: 5, wantErr: false},
		"long string":    {value: strings.Repeat("a", 1000), max: 100, wantErr: true},
		"exact boundary": {value: "ab", max: 2, wantErr: false},
		"one over":       {value: "abc", max: 2, wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			check := MaxLen(tc.value, tc.max, "title")
			err := check()
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestNumeric(t *testing.T) {
	tests := map[string]struct {
		value   string
		wantErr bool
	}{
		"integer":          {value: "42", wantErr: false},
		"negative":         {value: "-5", wantErr: false},
		"decimal":          {value: "3.14", wantErr: false},
		"zero":             {value: "0", wantErr: false},
		"not a number":     {value: "abc", wantErr: true},
		"empty string":     {value: "", wantErr: true},
		"mixed":            {value: "12abc", wantErr: true},
		"spaces":           {value: " 42 ", wantErr: true},
		"negative decimal": {value: "-0.5", wantErr: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			check := Numeric(tc.value, "weight")
			err := check()
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestFields(t *testing.T) {
	t.Run("all valid", func(t *testing.T) {
		errs := Fields(
			Required("hello", "name"),
			MinLen("password123", 10, "password"),
			MaxLen("short", 100, "title"),
			Numeric("42", "weight"),
		)
		if errs != nil {
			t.Fatalf("expected no errors, got: %v", errs)
		}
	})

	t.Run("single failure", func(t *testing.T) {
		errs := Fields(
			Required("hello", "name"),
			Required("", "email"),
			MinLen("password123", 10, "password"),
		)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got: %d", len(errs))
		}
		if errs[0].Field != "email" {
			t.Errorf("expected field: email, got: %s", errs[0].Field)
		}
	})

	t.Run("multiple failures", func(t *testing.T) {
		errs := Fields(
			Required("", "name"),
			Required("", "email"),
			MinLen("short", 10, "password"),
		)
		if len(errs) != 3 {
			t.Fatalf("expected 3 errors, got: %d", len(errs))
		}
	})

	t.Run("no checks", func(t *testing.T) {
		errs := Fields()
		if errs != nil {
			t.Fatalf("expected no errors, got: %v", errs)
		}
	})
}

func TestFieldErrorMessage(t *testing.T) {
	err := FieldError{Field: "email", Message: "email is required"}
	if err.Error() != "email is required" {
		t.Errorf("expected 'email is required', got: %s", err.Error())
	}
}
