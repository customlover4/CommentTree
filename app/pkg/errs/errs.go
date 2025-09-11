package errs

import "errors"

var (
	ErrDBNotAffected = errors.New("no one rows didn't affected")
	ErrDBNotFound = errors.New("not found")
)
