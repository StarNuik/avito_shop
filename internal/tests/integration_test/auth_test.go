package integration_test

import (
	"github.com/avito_shop/internal/client"
	"github.com/avito_shop/internal/dto"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Auth_HappyPath(t *testing.T) {
	require := require.New(t)

	// Arrange
	client := client.New("http://localhost:8080")
	req := dto.AuthRequest{
		Username: "user0x186A0",
		Password: "password",
	}

	// Act
	resp, err := client.Auth(req)

	// Assert
	require.NoError(err)
	require.NotNil(resp)
}
