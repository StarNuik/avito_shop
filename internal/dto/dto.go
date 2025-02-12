package dto

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type JwtPayload struct {
	UserId int64
}

// auto-generated
type InfoResponse struct {
	Coins       int64           `json:"coins"`
	Inventory   []InventoryInfo `json:"inventory"`
	CoinHistory struct {
		Received []BalanceDebitInfo  `json:"received"`
		Sent     []BalanceCreditInfo `json:"sent"`
	} `json:"coinHistory"`
}

type InventoryInfo struct {
	Type     string `json:"type"`
	Quantity int64  `json:"quantity"`
}

type BalanceDebitInfo struct {
	FromUser string `json:"fromUser"`
	Amount   int64  `json:"amount"`
}

type BalanceCreditInfo struct {
	ToUser string `json:"toUser"`
	Amount int64  `json:"amount"`
}
