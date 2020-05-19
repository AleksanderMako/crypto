package types

// Order represents a buy or sell Limit Order
type Order struct {
	WalletID    string  `json:"walletID"`
	CurrencyID  string  `json:"currencyID"`
	Price       float64 `json:"limitPrice"`
	OrderType   int     `json:"OrderType"`
	SumToInvest float64 `json:"sumToInvest"`
	UserID      string  `json:"userID"`
	Deleted     bool
}

// NewOrder returns a pointer to a new Order struct
func NewOrder(walletID, currencyID string, price float64, ordeType int, sumToInvest float64) *Order {

	return &Order{
		CurrencyID:  currencyID,
		WalletID:    walletID,
		Price:       price,
		SumToInvest: sumToInvest,
		OrderType:   ordeType,
	}
}
