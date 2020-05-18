package utils

import (
	"cryptoServer/database"
	"cryptoServer/database/types"

	"github.com/google/uuid"
)

func CreateCurrencies(db *database.Database) []string {

	names := []string{"BTC", "LTC", "ETH"}
	namesIDX := 0
	ids := []string{}
	for i := 0; i < 3; i++ {
		ID := uuid.New()
		db.Currencies[ID.String()] = *types.NewCurrency(names[namesIDX])
		namesIDX++
		ids = append(ids, ID.String())
	}
	return ids
}

func CreateWallets(db *database.Database, numberOfWallets int) []string {

	walletIDs := []string{}
	for i := 0; i < numberOfWallets; i++ {
		ID := uuid.New()
		ids := CreateCurrencies(db)
		db.Wallets[ID.String()] = *types.NewWallet(ids)
		walletIDs = append(walletIDs, ID.String())
	}
	return walletIDs

}
