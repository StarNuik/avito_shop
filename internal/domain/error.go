package domain

import "fmt"

var (
	ErrNotFound   error = fmt.Errorf("not found")
	ErrNotEnough        = fmt.Errorf("not enough")
	ErrNotAllowed       = fmt.Errorf("not allowed")
)
