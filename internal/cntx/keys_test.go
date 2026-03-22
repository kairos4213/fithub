package cntx

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestUserID(t *testing.T) {
	validUserID := uuid.New()
	userIDContextExistsAndValid := context.WithValue(t.Context(), UserIDKey, validUserID)
	userIDContextExistsAndNotValidType := context.WithValue(t.Context(), UserIDKey, "1234567890")
	userIDContextDoesNotExist := context.Background()

	tests := map[string]struct {
		userIDContext context.Context
		wantOk        bool
		wantUserID    uuid.UUID
	}{
		"exists and valid": {
			userIDContext: userIDContextExistsAndValid,
			wantOk:        true,
			wantUserID:    validUserID,
		},
		"invalid id": {
			userIDContext: userIDContextExistsAndNotValidType,
			wantOk:        false,
			wantUserID:    uuid.Nil,
		},
		"no id": {
			userIDContext: userIDContextDoesNotExist,
			wantOk:        false,
			wantUserID:    uuid.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			userID, ok := UserID(tc.userIDContext)
			if ok != tc.wantOk {
				t.Fatalf("expected ok: %t, got %t", tc.wantOk, ok)
			}
			if userID != tc.wantUserID {
				t.Fatalf("expected userID: %q, got %q", tc.wantUserID, userID)
			}
		})
	}
}
