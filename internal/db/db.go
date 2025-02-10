package db

type User struct {
	Id           int64
	Username     string
	PasswordHash string
}

type Transaction struct {
	Id       int64
	Delta    int64
	Result   int64
	UserFrom int64
	UserTo   int64
}

type InventoryEntry struct {
	Id    int64
	Name  string
	Price int64
}

type Purchase struct {
	Id   int64
	Item int64
	User int64
}

type Shop interface {
	User(username string) User
	//PurchasesOf(userId int64) []Purchase
	//TransactionsOf(userId int64) []Transaction
}
