package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/avito_shop/internal/dto"
	"io"
	"net/http"
)

var (
	ErrUnauthorized   = fmt.Errorf("unauthorized")
	ErrInternalServer = fmt.Errorf("internal server error")
	ErrBadRequest     = fmt.Errorf("bad request")
)

type client struct {
	hostUri string
}

func New(hostUri string) *client {
	return &client{
		hostUri: hostUri,
	}
}

func (c *client) Auth(req dto.AuthRequest) (*dto.AuthResponse, error) {
	endpointPrefix := "/api/auth"

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(reqBody)
	httpResponse, err := http.Post(c.hostUri+endpointPrefix, "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	statusCode := httpResponse.StatusCode
	if statusCode/100 != 2 {
		err := unmarshalError(statusCode, httpResponse.Body)
		return nil, fmt.Errorf("client.Auth: %w", err)
	}

	var resp dto.AuthResponse
	err = unmarshal(httpResponse.Body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func unmarshal(from io.Reader, v any) error {
	bytes, err := io.ReadAll(from)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

func unmarshalError(code int, from io.Reader) error {
	bytes, err := io.ReadAll(from)
	if err != nil {
		return err
	}

	dto := dto.ErrorResponse{}
	err = json.Unmarshal(bytes, &dto)
	if err != nil {
		return err
	}

	switch code {
	case 400:
		err = ErrBadRequest
	case 401:
		err = ErrUnauthorized
	case 500:
		err = ErrInternalServer
	default:
		return fmt.Errorf("unmarshalError: http code not supported")
	}

	return fmt.Errorf("%w: %s", err, dto.Errors)
}
