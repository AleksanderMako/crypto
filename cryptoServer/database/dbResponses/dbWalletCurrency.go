package dbResponses

type DBWalletCurrency struct {
	Balance    float64      `json:"balance"`
	Currencies []DBCurrency `json:"currencies"`
	WalletID   string       `json:"walletID"`
}

func NewDBWalletCurrency(currencies []DBCurrency, balance float64, walletID string) *DBWalletCurrency {

	return &DBWalletCurrency{
		Balance:    balance,
		Currencies: currencies,
		WalletID:   walletID,
	}
}
