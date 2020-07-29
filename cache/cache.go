package cache

import (
	"errors"
	"time"
)

// ErrEntryNotFound is an error type struct which is returned when entry was not found for provided key
var ErrEntryNotFound = errors.New("Entry not found")

type Cache interface {
	// Get reads entry for the key.
	// It returns an ErrEntryNotFound when
	// no entry exists for the given key.
	Get(key string) ([]byte, error)
	// Set saves entry under the key
	Set(key string, entry []byte, ex time.Duration) error
	// Delete removes the key
	Delete(key string) error
	// Reset empties all cache shards
	Reset() error
}
