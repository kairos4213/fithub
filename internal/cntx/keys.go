package cntx

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"

// UserID extracts the authenticated user's ID from the request context.
// Returns false if the key is missing or not a valid UUID, which indicates
// the auth middleware is misconfigured or the route is unprotected.
func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}
