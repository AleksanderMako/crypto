package database

import (
	dbResponses "cryptoServer/database/dbResponses"
	types "cryptoServer/database/types"
	"sync"

	"github.com/google/uuid"
)

type Storage interface {
	GetWallets() []dbResponses.DBWallet
	GetOrders(filter func(order types.Order) bool) []dbResponses.DBOrder
	DeleteOrder(key string)
	GetWallet(ID string) (w types.Wallet)
	GetCurrency(ID string) (c types.Currency)
	UpdateOrder(ID string, updatedData types.Order)
	CreateOrder(newOrder types.Order)
	UpdateWallet(ID string, wallet types.Wallet)
	UpdateCurrency(ID string, currency types.Currency)
	GetWalletsAndCurrencies() []dbResponses.DBWalletCurrency
	CreateCurrency(newCurrency types.Currency) string
}

// Database is the encapsulating struct of the wallets, currencies and orders map
type Database struct {
	Wallets          map[string]types.Wallet
	Currencies       map[string]types.Currency
	WalletCurrencies map[string]string
	Orders           map[string]types.Order
	mutex            sync.Mutex
}

// NewDatabase returns a pointer to a new Database struct
func NewDatabase() *Database {

	wallets := make(map[string]types.Wallet)
	currencies := make(map[string]types.Currency)
	walletCurrencies := make(map[string]string)
	orders := make(map[string]types.Order)

	return &Database{
		Wallets:          wallets,
		Currencies:       currencies,
		WalletCurrencies: walletCurrencies,
		Orders:           orders,
	}
}

// GetWallets returns all wallets
func (d *Database) GetWallets() []dbResponses.DBWallet {

	d.mutex.Lock()
	wallets := []dbResponses.DBWallet{}
	for k, v := range d.Wallets {
		wallets = append(wallets, *dbResponses.NewDBWallet(k, v))
	}
	d.mutex.Unlock()
	return wallets

}

// GetOrders returns a list of orers filtered by some condition
func (d *Database) GetOrders(filter func(order types.Order) bool) []dbResponses.DBOrder {
	// d.mutex.Lock()
	mapPtr := d.Orders
	// d.mutex.Unlock()
	orders := []dbResponses.DBOrder{}
	for k, v := range mapPtr {
		if filter(v) {
			orders = append(orders, *dbResponses.NewDBOrder(k, v))
		}
	}

	return orders
}

// DeleteOrder erases an order from the map given its ID
func (d *Database) DeleteOrder(key string) {
	// d.mutex.Lock()

	if _, ok := d.Orders[key]; ok {
		delete(d.Orders, key)
	}
	// d.mutex.Unlock()

}

// GetWallet returns a wallet struct given its ID
func (d *Database) GetWallet(ID string) (w types.Wallet) {
	// d.mutex.Lock()
	if wallet, ok := d.Wallets[ID]; ok {
		// d.mutex.Unlock()
		return wallet
	}
	// d.mutex.Unlock()

	return types.Wallet{}
}

// GetCurrency returns a currency struct given its ID
func (d *Database) GetCurrency(ID string) (c types.Currency) {
	// d.mutex.Lock()
	if currency, ok := d.Currencies[ID]; ok {
		// d.mutex.Unlock()
		return currency
	}
	// d.mutex.Unlock()

	return types.Currency{}
}

// UpdateOrder updates an existing order given the ID and an updated struct
func (d *Database) UpdateOrder(ID string, updatedData types.Order) {
	// d.mutex.Lock()

	if _, ok := d.Orders[ID]; ok {
		d.Orders[ID] = updatedData
	}
	// d.mutex.Unlock()

}

// CreateOrder creates and stores an orer struct in a in-memory map
func (d *Database) CreateOrder(newOrder types.Order) {
	// d.mutex.Lock()
	ID := uuid.New()
	d.Orders[ID.String()] = newOrder
	// d.mutex.Unlock()
}

//UpdateWallet updates an existing wallet given and ID and updated struct
func (d *Database) UpdateWallet(ID string, wallet types.Wallet) {
	// d.mutex.Lock()
	if _, ok := d.Wallets[ID]; ok {
		d.Wallets[ID] = wallet
	}
	// d.mutex.Unlock()

}

//UpdateCurrency updates an existing currency given an ID and updated struct
func (d *Database) UpdateCurrency(ID string, currency types.Currency) {
	// d.mutex.Lock()
	if _, ok := d.Currencies[ID]; ok {
		d.Currencies[ID] = currency
	}
	// d.mutex.Unlock()

}

// GetWalletsAndCurrencies retunrs the wallet struct merged with its corresponding currencies
func (d *Database) GetWalletsAndCurrencies() []dbResponses.DBWalletCurrency {

	wallets := d.GetWallets()
	// d.mutex.Lock()
	walletsAndCurrencies := []dbResponses.DBWalletCurrency{}
	for _, wallet := range wallets {

		currencies := []types.Currency{}
		for _, currency := range wallet.Wallet.Currencies {

			currencies = append(currencies, d.GetCurrency(currency))
		}
		walletsAndCurrencies = append(walletsAndCurrencies, *dbResponses.NewDBWalletCurrency(currencies, wallet.Wallet.Balance, wallet.ID))
	}
	// d.mutex.Unlock()

	return walletsAndCurrencies
}

// CreateCurrency writes a new currency struct into the in-memory Currencies Map
func (d *Database) CreateCurrency(newCurrency types.Currency) string {
	// d.mutex.Lock()
	ID := uuid.New().String()
	d.Currencies[ID] = newCurrency
	// d.mutex.Unlock()

	return ID
}
