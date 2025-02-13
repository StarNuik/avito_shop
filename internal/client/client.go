package client

import (
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"net/http"
)

type client struct {
	hostUrl   string
	roundtrip infra.RoundtripHandler
}

func New(hostUrl string) *client {
	return &client{
		hostUrl: hostUrl,
	}
}

// TODO: refactor?
// type Client interface
// type Impl struct
// func New() Client

func WithRoundtrip(hostUrl string, roundtrip infra.RoundtripHandler) *client {
	return &client{
		hostUrl:   hostUrl,
		roundtrip: roundtrip,
	}
}

func (c *client) url(suffix string) string {
	return c.hostUrl + suffix
}

func (c *client) Auth(req dto.AuthRequest) (dto.AuthResponse, error) {
	out := dto.AuthResponse{}

	url := c.url("/api/auth")
	httpReq := infra.HttpRequest{
		Method:         http.MethodPost,
		Url:            url,
		In:             req,
		Out:            &out,
		UnmarshalError: UnmarshalError,
		HttpRoundtrip:  c.roundtrip,
	}
	err := infra.DoHttp(httpReq)
	return out, err
}

func (c *client) Info(auth dto.AuthResponse) (dto.InfoResponse, error) {
	out := dto.InfoResponse{}

	url := c.url("/api/info")
	headers := map[string]string{
		"Authorization": "Bearer" + " " + auth.Token,
	}
	httpReq := infra.HttpRequest{
		Method:         http.MethodGet,
		Url:            url,
		Headers:        headers,
		In:             nil,
		Out:            &out,
		UnmarshalError: UnmarshalError,
		HttpRoundtrip:  c.roundtrip,
	}
	err := infra.DoHttp(httpReq)
	return out, err
}
