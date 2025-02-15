package repository_test

import (
	"context"
	"github.com/avito_shop/internal/domain"
	"github.com/avito_shop/internal/repository"
	"github.com/avito_shop/internal/shoptest"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

const postgresUrl = "postgres://postgres:password@localhost:5432/shop"

var (
	db *pgx.Conn
)

func TestMain(m *testing.M) {
	var err error
	db, err = pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		log.Panic(err)
	}

	os.Exit(m.Run())
}

func TestUser(t *testing.T) {
	// Arrange
	require := require.New(t)

	testRepo := shoptest.NewShopRepo(db)
	err := testRepo.Clear("Users")
	require.NoError(err)

	want := domain.User{
		Username:     "username",
		PasswordHash: "hash",
	}
	_, err = testRepo.InsertUser(want, 0)
	require.NoError(err)

	// Act
	ctx := context.Background()
	repo := repository.NewShopPostgres(db)
	have, err := repo.User(ctx, want.Username)
	require.NoError(err)

	// Assert
	require.Equal(have.Username, want.Username)
	require.Equal(have.PasswordHash, want.PasswordHash)
}

func TestInventoryItem(t *testing.T) {
	// Arrange
	require := require.New(t)

	testRepo := shoptest.NewShopRepo(db)
	err := testRepo.Clear("Inventory")
	require.NoError(err)

	want := domain.InventoryItem{
		Name:  "inventory-item",
		Price: 10,
	}
	_, err = db.Exec(context.Background(), `
        insert into Inventory (Name, Price)
        values ($1, $2)
    `, want.Name, want.Price)
	require.NoError(err)

	// Act
	ctx := context.Background()
	repo := repository.NewShopPostgres(db)
	have, err := repo.InventoryItem(ctx, want.Name)
	require.NoError(err)

	// Assert
	require.Equal(have.Name, want.Name)
	require.Equal(have.Price, want.Price)
}

func TestUserBalance(t *testing.T) {
	// Arrange
	require := require.New(t)

	testRepo := shoptest.NewShopRepo(db)
	err := testRepo.Clear("Users")
	require.NoError(err)

	want := domain.User{Username: "username"}
	wantBalance := int64(100)
	userId, err := testRepo.InsertUser(want, wantBalance)
	require.NoError(err)

	// Act
	ctx := context.Background()
	repo := repository.NewShopPostgres(db)
	tx, err := repo.Begin(ctx)
	require.NoError(err)

	haveBalance, err := tx.UserBalanceLock(userId)
	require.NoError(err)

	// Assert
	require.Equal(haveBalance, wantBalance)
}

func TestUserPairBalance(t *testing.T) {
	// Arrange
	require := require.New(t)

	testRepo := shoptest.NewShopRepo(db)
	err := testRepo.Clear("Users")
	require.NoError(err)

	wantBalances := []int64{50, 100}
	userIds := []int64{0, 0}
	userIds[0], err = testRepo.InsertUser(domain.User{}, wantBalances[0])
	require.NoError(err)
	userIds[1], err = testRepo.InsertUser(domain.User{}, wantBalances[1])
	require.NoError(err)

	// Act
	ctx := context.Background()
	repo := repository.NewShopPostgres(db)
	tx, err := repo.Begin(ctx)
	require.NoError(err)

	haveBalances := []int64{0, 0}
	haveBalances[0], haveBalances[1], err = tx.UserPairBalanceLock(userIds[0], userIds[1])
	require.NoError(err)

	// Assert
	require.Equal(haveBalances[0], wantBalances[0])
	require.Equal(haveBalances[1], wantBalances[1])
}

func TestUpdateBalance(t *testing.T) {
	// Arrange
	require := require.New(t)

	testRepo := shoptest.NewShopRepo(db)
	err := testRepo.Clear("Users")
	require.NoError(err)

	userId, err := testRepo.InsertUser(domain.User{}, 0)
	require.NoError(err)

	wantBalance := int64(100)

	// Act
	ctx := context.Background()
	repo := repository.NewShopPostgres(db)
	tx, err := repo.Begin(ctx)
	require.NoError(err)

	err = tx.UpdateBalance(userId, wantBalance)
	require.NoError(err)

	// Assert
	haveBalance, err := tx.UserBalanceLock(userId)
	require.NoError(err)

	require.Equal(haveBalance, wantBalance)
}
