// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Goal struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Description    string
	GoalDate       time.Time
	CompletionDate sql.NullTime
	Notes          sql.NullString
	Status         string
	UserID         uuid.UUID
}

type RefreshToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FirstName      string
	MiddleName     sql.NullString
	LastName       string
	Email          string
	HashedPassword string
	ProfileImage   sql.NullString
	Preferences    pqtype.NullRawMessage
}
