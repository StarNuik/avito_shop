package e2e

import (
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendCoins(t *testing.T) {
	// Arrange
	require := require.New(t)

	shoptest.ClearRepo()
	client := shoptest.NewTestClient()

	// Act
	users := []dto.AuthRequest{
		shoptest.AuthRequest(0),
		shoptest.AuthRequest(1),
		shoptest.AuthRequest(2),
	}
	auth, err := client.Auth(users[0].Username, users[0].Password)
	require.NoError(err)

	err = client.SendCoins(auth, users[1].Username, int64(30))
	require.NoError(err)
	_ = client.SendCoins(auth, users[2].Username, int64(30))
	_ = client.SendCoins(auth, users[1].Username, int64(60))
	_ = client.SendCoins(auth, users[2].Username, int64(60))
	_ = client.SendCoins(auth, users[1].Username, int64(90))
	_ = client.SendCoins(auth, users[2].Username, int64(90))

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

	require.Equal(info.CoinHistory.Sent[0].ToUser, users[2].Username)
	require.Equal(info.CoinHistory.Sent[1].ToUser, users[1].Username)
	require.Equal(info.CoinHistory.Sent[2].ToUser, users[2].Username)
	require.Equal(info.CoinHistory.Sent[3].ToUser, users[1].Username)
	require.Equal(info.CoinHistory.Sent[4].ToUser, users[2].Username)
	require.Equal(info.CoinHistory.Sent[5].ToUser, users[1].Username)

	auth, err = client.Auth(users[1].Username, users[1].Password)
	require.NoError(err)

	info, err = client.Info(auth)
	require.NoError(err)

	require.Equal(info.Coins, shoptest.DefaultBalance+int64(30)+int64(60)+int64(90))
	require.Len(info.CoinHistory.Received, 3)

	require.Equal(info.CoinHistory.Received[0].Amount, int64(90))
	require.Equal(info.CoinHistory.Received[1].Amount, int64(60))
	require.Equal(info.CoinHistory.Received[2].Amount, int64(30))

	require.Equal(info.CoinHistory.Received[0].FromUser, users[0].Username)
	require.Equal(info.CoinHistory.Received[1].FromUser, users[0].Username)
	require.Equal(info.CoinHistory.Received[2].FromUser, users[0].Username)
}
