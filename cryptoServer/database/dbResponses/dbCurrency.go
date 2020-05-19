package dbResponses

import "cryptoServer/database/types"

type DBCurrency struct {
	ID       string         `json:"id"`
	Currency types.Currency `json:"currency"`
}

func NewDBCurrency(currency types.Currency, ID string) *DBCurrency {
	return &DBCurrency{
		ID:       ID,
		Currency: currency,
	}
}
