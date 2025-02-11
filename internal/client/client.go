package client

import (
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"net/http"
)

type client struct {
	hostUrl string
}

func New(hostUri string) *client {
	return &client{
		hostUrl: hostUri,
	}
}

func (c *client) url(suffix string) string {
	return c.hostUrl + suffix
}

func (c *client) Auth(req dto.AuthRequest) (*dto.AuthResponse, error) {
	resp := &dto.AuthResponse{}
	url := c.url("/api/auth")
	err := infra.HttpRequest(http.MethodPost, url, nil, UnmarshalError, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
