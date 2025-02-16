package domain

type User struct {
	Id           int64
	Username     string
	PasswordHash string
	//Coins        int64
}

type Transfer struct {
	Id       int64
	Delta    int64
	FromUser int64
	ToUser   int64
}

type InventoryItem struct {
	Id    int64
	Name  string
	Price int64
}

type Purchase struct {
	Id     int64
	Item   int64
	UserId int64
	Price  int64
}

type InventoryInfo struct {
	Name     string
	Quantity int64
}

type TransferInfo struct {
	Delta        int64
	FromUser     int64
	ToUser       int64
	FromUsername string
	ToUsername   string
}
