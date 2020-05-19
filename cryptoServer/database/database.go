package database

import (
	dbResponses "cryptoServer/database/dbResponses"
	types "cryptoServer/database/types"
	"sync"

	"github.com/google/uuid"
)

// Storage is an interface defining the storage behavior for mocking and testing
type Storage interface {
	GetWallets() []dbResponses.DBWallet
	GetOrders(filter func(order types.Order) bool) []dbResponses.DBOrder
	DeleteOrder(key string)
	GetWallet(ID string) (w types.Wallet)
	GetCurrency(ID string) (c types.Currency)
	UpdateOrder(ID string, updatedData types.Order)
	CreateOrder(newOrder types.Order) string
	UpdateWallet(ID string, wallet types.Wallet)
	UpdateCurrency(ID string, currency types.Currency)
	GetWalletsAndCurrencies() []dbResponses.DBWalletCurrency
	CreateCurrency(newCurrency types.Currency) string
	GetOrdersByType() *map[int][]dbResponses.DBOrder
}

// Database is the encapsulating struct of the wallets, currencies and orders map
type Database struct {
	Wallets    map[string]types.Wallet
	Currencies map[string]types.Currency
	Orders     map[string]types.Order
	mutex      sync.Mutex
}

// NewDatabase returns a pointer to a new Database struct
func NewDatabase() *Database {

	wallets := make(map[string]types.Wallet)
	currencies := make(map[string]types.Currency)
	orders := make(map[string]types.Order)

	return &Database{
		Wallets:    wallets,
		Currencies: currencies,
		Orders:     orders,
	}
}

// GetWallets returns all wallets
func (d *Database) GetWallets() []dbResponses.DBWallet {

	wallets := []dbResponses.DBWallet{}
	for k, v := range d.Wallets {
		wallets = append(wallets, *dbResponses.NewDBWallet(k, v))
	}
	return wallets

}

// GetOrders returns a list of orers filtered by some condition
func (d *Database) GetOrders(filter func(order types.Order) bool) []dbResponses.DBOrder {
	orders := []dbResponses.DBOrder{}
	for k, v := range d.Orders {
		if filter(v) {
			orders = append(orders, *dbResponses.NewDBOrder(k, v))
		}
	}

	return orders
}

// GetOrdersByType returns a pointer to a map of orders grouped by type of order
func (d *Database) GetOrdersByType() *map[int][]dbResponses.DBOrder {

	orderGroups := make(map[int][]dbResponses.DBOrder)
	for k, v := range d.Orders {
		orderGroups[v.OrderType] = append(orderGroups[v.OrderType], *dbResponses.NewDBOrder(k, v))
	}
	return &orderGroups
}

// DeleteOrder erases an order from the map given its ID
func (d *Database) DeleteOrder(key string) {

	if order, ok := d.Orders[key]; ok {
		order.Deleted = true
	}
}

// GetWallet returns a wallet struct given its ID
func (d *Database) GetWallet(ID string) (w types.Wallet) {
	if wallet, ok := d.Wallets[ID]; ok {
		return wallet
	}
	return types.Wallet{}
}

// GetCurrency returns a currency struct given its ID
func (d *Database) GetCurrency(ID string) (c types.Currency) {
	if currency, ok := d.Currencies[ID]; ok {
		return currency
	}
	return types.Currency{}
}

// UpdateOrder updates an existing order given the ID and an updated struct
func (d *Database) UpdateOrder(ID string, updatedData types.Order) {

	if _, ok := d.Orders[ID]; ok {
		d.Orders[ID] = updatedData
	}
}

// CreateOrder creates and stores an order struct in a in-memory map
func (d *Database) CreateOrder(newOrder types.Order) string {
	ID := uuid.New().String()
	d.Orders[ID] = newOrder
	return ID
}

//UpdateWallet updates an existing wallet given and ID and updated struct
func (d *Database) UpdateWallet(ID string, wallet types.Wallet) {
	if _, ok := d.Wallets[ID]; ok {
		d.Wallets[ID] = wallet
	}

}

//UpdateCurrency updates an existing currency given an ID and updated struct
func (d *Database) UpdateCurrency(ID string, currency types.Currency) {
	if _, ok := d.Currencies[ID]; ok {
		d.Currencies[ID] = currency
	}
}

// GetWalletsAndCurrencies retunrs the wallet struct merged with its corresponding currencies
func (d *Database) GetWalletsAndCurrencies() []dbResponses.DBWalletCurrency {

	wallets := d.GetWallets()
	walletsAndCurrencies := []dbResponses.DBWalletCurrency{}
	for _, wallet := range wallets {

		currencies := []types.Currency{}
		for _, currency := range wallet.Wallet.Currencies {

			currencies = append(currencies, d.GetCurrency(currency))
		}
		walletsAndCurrencies = append(walletsAndCurrencies, *dbResponses.NewDBWalletCurrency(currencies, wallet.Wallet.Balance, wallet.ID))
	}

	return walletsAndCurrencies
}

// CreateCurrency writes a new currency struct into the in-memory Currencies Map
func (d *Database) CreateCurrency(newCurrency types.Currency) string {
	ID := uuid.New().String()
	d.Currencies[ID] = newCurrency
	return ID
}
