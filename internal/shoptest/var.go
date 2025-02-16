package shoptest

import (
	"fmt"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/dto"
)

const (
	HostUrl        = "http://localhost:8080"
	DefaultBalance = int64(1000)
	UserCount      = 100_000
)

type user dto.AuthRequest

var Users = []user{
	{Username: "user#0", Password: "user#0"},
	{Username: "user#1", Password: "user#1"},
	{Username: "user#2", Password: "user#2"},
}

var Inventory = []domain.InventoryItem{
	{Name: "t-shirt", Price: 80},
	{Name: "cup", Price: 20},
	{Name: "book", Price: 50},
	{Name: "pen", Price: 10},
	{Name: "powerbank", Price: 200},
	{Name: "hoody", Price: 300},
	{Name: "umbrella", Price: 200},
	{Name: "socks", Price: 10},
	{Name: "wallet", Price: 50},
	{Name: "pink-hoody", Price: 500},
}

func User(idx int) user {
	return user{
		Username: fmt.Sprintf("user#%d", idx),
		Password: fmt.Sprintf("password#%d", idx),
	}
}

func AuthRequest(idx int) dto.AuthRequest {
	return dto.AuthRequest(User(idx))
}
