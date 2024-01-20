package repository

import "errors"

// ErrNotFound is returned when the requested resource is not found.
var ErrNotFound = errors.New("not found")
