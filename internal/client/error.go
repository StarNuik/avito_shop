package client

import (
	"encoding/json"
	"fmt"
	"github.com/avito_shop/internal/dto"
)

var (
	ErrUnauthorized   = fmt.Errorf("unauthorized")
	ErrInternalServer = fmt.Errorf("internal server error")
	ErrBadRequest     = fmt.Errorf("bad request")
)

func UnmarshalError(status int, body []byte) error {
	if status/100 == 2 {
		return nil
	}

	dto := dto.ErrorResponse{}
	err := json.Unmarshal(body, &dto)
	if err != nil {
		return err
	}

	switch status {
	case 400:
		err = ErrBadRequest
	case 401:
		err = ErrUnauthorized
	case 500:
		err = ErrInternalServer
	default:
		return fmt.Errorf("UnmarshalError: http status code not supported")
	}

	return fmt.Errorf("%w: %s", err, dto.Errors)
}
