package controller

import (
	"cryptoServer/database"
	"cryptoServer/database/types"
	"cryptoServer/transactions"
	"encoding/json"
	"sync"
)

// Controller encapsulates the API functionalities
type Controller struct {
	db    database.Storage
	te    transactions.TransactionEngine
	mutex sync.Mutex
}

// NewController returns a pointer to a new Controller struct
func NewController(db database.Storage, te transactions.TransactionEngine) *Controller {

	return &Controller{
		db: db,
		te: te,
	}
}

// ListWalletBalances returns wallet balances
func (c *Controller) ListWalletBalances() ([]byte, error) {

	c.mutex.Lock()
	wallets := c.db.GetWalletsAndCurrencies()
	c.mutex.Unlock()
	response, err := json.Marshal(wallets)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// ListOrdersByUser returns user orders given ID
func (c *Controller) ListOrdersByUser(userID string) ([]byte, error) {

	c.mutex.Lock()
	userOrders := c.db.GetOrders(func(order types.Order) bool {
		if order.UserID == userID {
			return true
		}
		return false
	})
	c.mutex.Unlock()
	response, err := json.Marshal(userOrders)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// CancelOrder deletes an order struct given its ID
func (c *Controller) CancelOrder(ID string, userID string) {
	c.mutex.Lock()
	c.db.DeleteOrder(ID)
	c.mutex.Unlock()
}

// PlaceOrder places either a buy or a sell limit order
func (c *Controller) PlaceOrder(order types.Order) ([]byte, error) {

	c.mutex.Lock()
	status, err := c.te.PlaceOrder(order)
	if err != nil {
		return nil, err
	}
	c.mutex.Unlock()
	return []byte(status), nil
}
