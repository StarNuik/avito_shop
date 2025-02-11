package integration_test

import (
	"github.com/avito_shop/internal/client"
	"github.com/avito_shop/internal/dto"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	hostUri         = "http://localhost:8080"
	usernameCorrect = "0x186A0"
	passwordCorrect = "password"
)

func Test_Auth_HappyPath_ReturnsToken(t *testing.T) {
	// Arrange
	require := require.New(t)

	client := client.New(hostUri)

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

	shop := client.New(hostUri)

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

	shop := client.New(hostUri)

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

//func Test_Auth_IncorrectDto_BadRequest(t *testing.T) {
//    // Arrange
//    require := require.New(t)
//
//    body := []byte("{}")
//    reader := bytes.NewReader(body)
//    response, err := http.Post(hostUri+"/api/auth", "application/json", reader)
//
//    // Act
//    req := dto.AuthRequest{
//        Username: usernameCorrect,
//        Password: "bad password",
//    }
//
//    _, err := client.Auth(req)
//
//    // Assert
//    require.Error(err)
//    require.ErrorIs(err, client.ErrUnauthorized)
//}
