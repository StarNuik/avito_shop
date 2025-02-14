package e2e_test

import (
	"github.com/avito_shop/internal/client"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInfo(t *testing.T) {
	// Arrange
	require := require.New(t)

	client := client.NewTestClient()

	// Act
	auth, err := client.Auth(shoptest.Users[0].Username, shoptest.Users[0].Password)
	require.NoError(err)

	info, err := client.Info(auth)
	require.NoError(err)

	// Assert
	require.Equal(info.Coins, shoptest.DefaultBalance)
	require.Empty(info.Inventory)
	require.Empty(info.CoinHistory.Sent)
	require.Empty(info.CoinHistory.Received)
}
