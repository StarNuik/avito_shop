package e2e_test

import (
	"github.com/avito_shop/internal/client"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_InfoScenario(t *testing.T) {
	// Arrange
	require := require.New(t)

	client := client.New(shoptest.HostUrl)

	// Act
	authReq := dto.AuthRequest{
		Username: shoptest.Username,
		Password: shoptest.Password,
	}
	auth, err := client.Auth(authReq)
	require.NoError(err)

	info, err := client.Info(auth)
	require.NoError(err)

	// Assert
	require.Equal(info.Coins, shoptest.DefaultBalance)
	require.Empty(info.Inventory)
	require.Empty(info.CoinHistory.Sent)
	require.Empty(info.CoinHistory.Received)
}
