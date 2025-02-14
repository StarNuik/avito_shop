package client

import (
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"net/http"
)

type Client interface {
	Auth(username string, password string) (string, error)
	Info(authToken string) (dto.InfoResponse, error)
	SendCoins(authToken string, toUsername string, amount int64) error
	BuyItem(authToken string, itemName string) error
}

var _ Client = (*Impl)(nil)

type Impl struct {
	infra.HttpEngine
	HostUrl string
}

func New(hostUrl string) Client {
	return &Impl{
		HostUrl: hostUrl,
		HttpEngine: infra.HttpEngine{
			ErrHandler: UnmarshalError,
		},
	}
}

func (c *Impl) url(suffix string) string {
	return c.HostUrl + suffix
}

func (c *Impl) Auth(username string, password string) (string, error) {
	out := dto.AuthResponse{}

	url := c.url("/api/auth")
	req := dto.AuthRequest{Username: username, Password: password}
	err := c.Do(http.MethodPost, url, nil, req, &out)
	return out.Token, err
}

func (c *Impl) Info(token string) (dto.InfoResponse, error) {
	out := dto.InfoResponse{}

	url := c.url("/api/info")
	headers := map[string]string{
		"Authorization": "Bearer" + " " + token,
	}
	err := c.Do(http.MethodGet, url, headers, nil, &out)
	return out, err
}

func (c *Impl) SendCoins(token string, toUsername string, amount int64) error {
	url := c.url("/api/sendCoin")
	headers := map[string]string{
		"Authorization": "Bearer" + " " + token,
	}
	req := dto.SendCoinRequest{ToUser: toUsername, Amount: amount}
	err := c.Do(http.MethodPost, url, headers, req, nil)
	return err
}

func (c *Impl) BuyItem(token string, itemName string) error {
	url := c.url("/api/buy/" + itemName)
	headers := map[string]string{
		"Authorization": "Bearer" + " " + token,
	}
	err := c.Do(http.MethodGet, url, headers, nil, nil)
	return err
}
