package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   error = fmt.Errorf("not found")
	ErrNotEnough        = fmt.Errorf("not enough")
	ErrNotAllowed       = fmt.Errorf("not allowed")
)

func IsDomainError(err error) bool {
	return errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrNotAllowed) ||
		errors.Is(err, ErrNotEnough)
}
