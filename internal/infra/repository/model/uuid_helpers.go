package model

import "github.com/google/uuid"

// nilIfZeroUUID returns nil if the UUID is the zero value, otherwise returns a pointer to it.
func nilIfZeroUUID(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

// derefUUID safely dereferences a *uuid.UUID, returning uuid.Nil if the pointer is nil.
func derefUUID(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}
