package dbResponses

import "cryptoServer/database/types"

type DBWalletCurrency struct {
	Balance    float64          `json:"balance"`
	Currencies []types.Currency `json:"currencies"`
	WalletID   string           `json:"walletID"`
}

func NewDBWalletCurrency(currencies []types.Currency, balance float64, walletID string) *DBWalletCurrency {

	return &DBWalletCurrency{
		Balance:    balance,
		Currencies: currencies,
		WalletID:   walletID,
	}
}
