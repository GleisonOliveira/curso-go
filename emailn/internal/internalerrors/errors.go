package internalerrors

import "errors"

type ErrNotFound struct {
	Entity string
}

func (e *ErrNotFound) Error() string {
	return e.Entity + " not found"
}

var ErrInternal error = errors.New("Internal Server Error")

func ProcessError(err error) error {
	var notFoundErr *ErrNotFound

	if errors.As(err, &notFoundErr) {
		return notFoundErr
	}

	return ErrInternal
}
