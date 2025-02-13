package integration_test

import (
	"github.com/avito_shop/internal/client"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_Auth_HappyPath_ReturnsToken(t *testing.T) {
	// Arrange
	require := require.New(t)

	client := client.New(shoptest.HostUrl)

	// Act
	req := dto.AuthRequest{
		Username: shoptest.Username,
		Password: shoptest.Password,
	}

	resp, err := client.Auth(req)

	// Assert
	require.NoError(err)
	require.NotEmpty(resp.Token)
}

func Test_Auth_IncorrectUsername_Unauthorized(t *testing.T) {
	// Arrange
	require := require.New(t)

	shop := client.New(shoptest.HostUrl)

	// Act
	req := dto.AuthRequest{
		Username: "incorrect username",
		Password: shoptest.Password,
	}

	_, err := shop.Auth(req)

	// Assert
	require.ErrorIs(err, client.ErrUnauthorized)
}

func Test_Auth_IncorrectPassword_Unauthorized(t *testing.T) {
	// Arrange
	require := require.New(t)

	shop := client.New(shoptest.HostUrl)

	// Act
	req := dto.AuthRequest{
		Username: shoptest.Username,
		Password: "bad password",
	}

	_, err := shop.Auth(req)

	// Assert
	require.Error(err)
	require.ErrorIs(err, client.ErrUnauthorized)
}

func Test_Auth_IncorrectDto_BadRequest(t *testing.T) {
	// Arrange
	require := require.New(t)

	req := struct {
		Nameuser string `json:"nameuser"`
		Wordpass string `json:"wordpass"`
	}{"user", "password"}

	// Act
	err := infra.HttpRequest(http.MethodPost,
		shoptest.HostUrl+"/api/auth",
		nil,
		client.UnmarshalError,
		req,
		nil)

	// Assert
	require.Error(err)
	require.ErrorIs(err, client.ErrBadRequest)
}
