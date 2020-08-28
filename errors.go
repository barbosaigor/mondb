package mondb

import "errors"

var (
	// ErrDBNotConnected error when there isn't connection to database
	ErrDBNotConnected = errors.New("No connection to database")
	// ErrDocumentNotFound error when no Document was found
	ErrDocumentNotFound = errors.New("No Document found")
	// ErrEmptyObject error when no object were defined
	ErrEmptyObject = errors.New("Empty Object")
	// ErrInvalidID error when the ID is not valid
	ErrInvalidID = errors.New("Invalid ID")
)
