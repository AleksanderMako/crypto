package dbResponses

import "cryptoServer/database/types"

type DBWallet struct {
	ID     string
	Wallet types.Wallet
}

func NewDBWallet(ID string, wallet types.Wallet) *DBWallet {

	return &DBWallet{
		ID:     ID,
		Wallet: wallet,
	}
}
