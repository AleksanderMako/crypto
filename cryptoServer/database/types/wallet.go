package types

// Wallet is the struct representing the Wallet entity
type Wallet struct {
	Currencies []string
	Balance    float64
}

// NewWallet returns a new pointer to Wallet object given a list of currency names that should be contained
func NewWallet(currencyIDs []string) *Wallet {

	return &Wallet{
		Balance:    1000,
		Currencies: currencyIDs,
	}
}
