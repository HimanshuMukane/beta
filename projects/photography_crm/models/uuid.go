package models

import (
	"fmt"

	"github.com/google/uuid"
)

// UUIDToString converts a UUID to its canonical string representation
// for storage in SQLite TEXT columns.
func UUIDToString(id uuid.UUID) string {
	return id.String()
}

// UUIDFromString converts a SQLite TEXT value to a UUID.
func UUIDFromString(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("decode UUID: %w", err)
	}

	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("UUID cannot be nil")
	}

	return id, nil
}
