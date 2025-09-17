package errs

import "errors"

var (
	// From db realisation to upper layer use ErrDB preffix.
	// ErrDBNotAffected on update/delete when no one fields not affected.
	ErrDBNotAffected = errors.New("no one rows didn't affected")
	// ErrDBNotFound on get when we can't find our field.
	ErrDBNotFound = errors.New("not found")
	// ErrDBViolatesForeignKey when foreign key is wrong.
	ErrDBViolatesForeignKey = errors.New("violates foreign key")
)
