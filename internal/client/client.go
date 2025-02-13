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

func (c *client) Auth(req dto.AuthRequest) (dto.AuthResponse, error) {
	out := dto.AuthResponse{}
	url := c.url("/api/auth")
	err := infra.HttpRequest(http.MethodPost, url, nil, UnmarshalError, req, &out)
	return out, err
}

func (c *client) Info(auth dto.AuthResponse) (dto.InfoResponse, error) {
	out := dto.InfoResponse{}
	url := c.url("/api/info")
	headers := map[string]string{
		"Authorization": "Bearer" + auth.Token,
	}
	err := infra.HttpRequest(http.MethodGet, url, headers, UnmarshalError, nil, &out)
	return out, err
}
