package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/avito_shop/internal/dto"
	"io"
	"net/http"
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

	if httpResponse.StatusCode/100 != 2 {
		return nil, unmarshalError(httpResponse.Body)
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

func unmarshalError(from io.Reader) error {
	bytes, err := io.ReadAll(from)
	if err != nil {
		return err
	}

	dto := dto.ErrorResponse{}
	err = json.Unmarshal(bytes, &dto)
	if err != nil {
		return err
	}

	return fmt.Errorf("ErrorResponse: %s", dto.Errors)
}
