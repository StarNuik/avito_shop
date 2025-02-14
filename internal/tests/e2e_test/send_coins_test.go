package e2e

import (
	"github.com/avito_shop/internal/client"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendCoins(t *testing.T) {
	// Arrange
	require := require.New(t)

	client := client.NewTestClient()

	// Act
	auth, err := client.Auth(shoptest.Users[0].Username, shoptest.Users[0].Password)
	require.NoError(err)

	err = client.SendCoins(auth, shoptest.Users[1].Username, int64(30))
	require.NoError(err)
	_ = client.SendCoins(auth, shoptest.Users[2].Username, int64(30))
	_ = client.SendCoins(auth, shoptest.Users[1].Username, int64(60))
	_ = client.SendCoins(auth, shoptest.Users[2].Username, int64(60))
	_ = client.SendCoins(auth, shoptest.Users[1].Username, int64(90))
	_ = client.SendCoins(auth, shoptest.Users[2].Username, int64(90))

	// Assert
	info, err := client.Info(auth)
	require.NoError(err)

	require.Equal(info.Coins, shoptest.DefaultBalance-2*int64(30)-2*int64(60)-2*int64(90))
	require.Len(info.CoinHistory.Sent, 6)

	require.Equal(info.CoinHistory.Sent[0].Amount, int64(90))
	require.Equal(info.CoinHistory.Sent[1].Amount, int64(90))
	require.Equal(info.CoinHistory.Sent[2].Amount, int64(60))
	require.Equal(info.CoinHistory.Sent[3].Amount, int64(60))
	require.Equal(info.CoinHistory.Sent[4].Amount, int64(30))
	require.Equal(info.CoinHistory.Sent[5].Amount, int64(30))

	require.Equal(info.CoinHistory.Sent[0].ToUser, shoptest.Users[2].Username)
	require.Equal(info.CoinHistory.Sent[1].ToUser, shoptest.Users[1].Username)
	require.Equal(info.CoinHistory.Sent[2].ToUser, shoptest.Users[2].Username)
	require.Equal(info.CoinHistory.Sent[3].ToUser, shoptest.Users[1].Username)
	require.Equal(info.CoinHistory.Sent[4].ToUser, shoptest.Users[2].Username)
	require.Equal(info.CoinHistory.Sent[5].ToUser, shoptest.Users[1].Username)

	auth, err = client.Auth(shoptest.Users[1].Username, shoptest.Users[1].Password)
	require.NoError(err)

	info, err = client.Info(auth)
	require.NoError(err)

	require.Equal(info.Coins, shoptest.DefaultBalance+int64(30)+int64(60)+int64(90))
	require.Len(info.CoinHistory.Received, 3)

	require.Equal(info.CoinHistory.Received[0].Amount, int64(90))
	require.Equal(info.CoinHistory.Received[1].Amount, int64(60))
	require.Equal(info.CoinHistory.Received[2].Amount, int64(30))

	require.Equal(info.CoinHistory.Received[0].FromUser, shoptest.Users[0].Username)
	require.Equal(info.CoinHistory.Received[1].FromUser, shoptest.Users[0].Username)
	require.Equal(info.CoinHistory.Received[2].FromUser, shoptest.Users[0].Username)
}
