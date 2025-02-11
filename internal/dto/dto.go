package dto

import "fmt"

var (
	ErrUnauthorized   = fmt.Errorf("unauthorized")
	ErrInternalServer = fmt.Errorf("internal server error")
	ErrBadRequest     = fmt.Errorf("bad request")
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type JwtPayload struct {
	UserId int64
}
