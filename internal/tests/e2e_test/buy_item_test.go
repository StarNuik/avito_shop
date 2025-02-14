package e2e

import (
	"github.com/avito_shop/internal/client"
	"github.com/avito_shop/internal/dto"
	"github.com/avito_shop/internal/shoptest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuyItem(t *testing.T) {
	// Arrange
	require := require.New(t)

	client := client.NewTestClient()

	// Act
	auth, err := client.Auth(shoptest.Users[0].Username, shoptest.Users[0].Password)
	require.NoError(err)

	err = client.BuyItem(auth, shoptest.Inventory[0].Name)
	require.NoError(err)
	_ = client.BuyItem(auth, shoptest.Inventory[0].Name)
	_ = client.BuyItem(auth, shoptest.Inventory[0].Name)
	_ = client.BuyItem(auth, shoptest.Inventory[0].Name)
	_ = client.BuyItem(auth, shoptest.Inventory[0].Name)

	_ = client.BuyItem(auth, shoptest.Inventory[1].Name)
	_ = client.BuyItem(auth, shoptest.Inventory[1].Name)
	_ = client.BuyItem(auth, shoptest.Inventory[1].Name)

	// Assert
	info, err := client.Info(auth)
	require.NoError(err)

	coinsSpent := 5*shoptest.Inventory[0].Price + 3*shoptest.Inventory[1].Price
	require.Equal(info.Coins, shoptest.DefaultBalance-coinsSpent)

	require.Len(info.Inventory, 2)
	require.Contains(info.Inventory, dto.InventoryInfo{
		Type:     shoptest.Inventory[0].Name,
		Quantity: 5,
	})
	require.Contains(info.Inventory, dto.InventoryInfo{
		Type:     shoptest.Inventory[1].Name,
		Quantity: 3,
	})
}
