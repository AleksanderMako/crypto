package types

// Currency represents the cryptocurrency entity
type Currency struct {
	Name     string `json:"name"`
	Ammount  int    `json:"amount"`
	WalletID string `json:"walletID"`
}

// NewCurrency returns a new pointer to Currency object given the name of a crypto
func NewCurrency(name string) *Currency {

	return &Currency{
		Ammount:  1000,
		Name:     name,
		WalletID: "",
	}
}
