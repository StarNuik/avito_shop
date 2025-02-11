package domain

type User struct {
	Id           int64
	Username     string
	PasswordHash string
}

type BalanceOperation struct {
	Id     int64
	User   int64
	Delta  int64
	Result int64
}

type Transfer struct {
	Id       int64
	SourceOp int64
	TargetOp int64
}

type InventoryEntry struct {
	Id    int64
	Name  string
	Price int64
}

type Purchase struct {
	Id        int64
	Item      int64
	User      int64
	Operation int64
}
