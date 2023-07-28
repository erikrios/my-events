package persistence

import "errors"

var (
	ErrNotFound = errors.New("Resources not found.")
	ErrDatabase = errors.New("Something went wrong with database.")
)
