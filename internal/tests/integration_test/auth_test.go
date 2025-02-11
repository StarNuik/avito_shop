package integration_test

import (
	"github.com/avito_shop/internal/client"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/infra"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const (
	hostUrl         = "http://localhost:8080"
	usernameCorrect = "0x186A0"
	passwordCorrect = "password"
)

func Test_Auth_HappyPath_ReturnsToken(t *testing.T) {
	// Arrange
	require := require.New(t)

	client := client.New(hostUrl)

	// Act
	req := dto.AuthRequest{
		Username: usernameCorrect,
		Password: passwordCorrect,
	}

	resp, err := client.Auth(req)

	// Assert
	require.NoError(err)
	require.NotNil(resp)
}

func Test_Auth_IncorrectUsername_Unauthorized(t *testing.T) {
	// Arrange
	require := require.New(t)

	shop := client.New(hostUrl)

	// Act
	req := dto.AuthRequest{
		Username: "incorrect username",
		Password: passwordCorrect,
	}

	_, err := shop.Auth(req)

	// Assert
	require.Error(err)
	require.ErrorIs(err, client.ErrUnauthorized)
}

func Test_Auth_IncorrectPassword_Unauthorized(t *testing.T) {
	// Arrange
	require := require.New(t)

	shop := client.New(hostUrl)

	// Act
	req := dto.AuthRequest{
		Username: usernameCorrect,
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
	err := infra.HttpRequest(http.MethodPost, hostUrl+"/api/auth", nil, client.UnmarshalError, req, nil)

	// Assert
	require.Error(err)
	require.ErrorIs(err, client.ErrBadRequest)
}
