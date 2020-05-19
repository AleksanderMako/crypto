package controller

import (
	"cryptoServer/database"
	"cryptoServer/database/types"
	"cryptoServer/transactions"
	"encoding/json"
	"sort"
	"sync"

	"github.com/google/uuid"
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
func (c *Controller) CancelOrder(ID string) {
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

// ListOrderBook returns the top 10 highest buy and lowest sell prices
func (c *Controller) ListOrderBook() ([]byte, error) {
	c.mutex.Lock()
	groups := c.db.GetOrdersByType()
	orderBook := OrderBook{}
	if sellorders, ok := (*groups)[types.SellOrder]; ok {
		sort.Slice(sellorders, func(i, j int) bool {

			if sellorders[i].Order.Price < sellorders[j].Order.Price {
				return true
			}
			return false
		})

		sellorders = sellorders[:min(10, len(sellorders))]
		orderBook.LowestSellOrders = sellorders
	}
	if buyOrders, ok := (*groups)[types.BuyOrder]; ok {

		sort.Slice(buyOrders, func(i, j int) bool {

			if buyOrders[i].Order.Price > buyOrders[j].Order.Price {
				return true
			}
			return false
		})
		buyOrders = buyOrders[:min(10, len(buyOrders))]
		orderBook.HighestBuyOrders = buyOrders
	}
	c.mutex.Unlock()
	response, err := json.Marshal(orderBook)
	if err != nil {
		return nil, err
	}
	return response, nil

}

// RegisterUser generates a user ID to be included in all requests
func (c *Controller) RegisterUser() string {
	ID := uuid.New().String()
	c.mutex.Lock()
	c.db.RegisterUser(ID)
	c.mutex.Unlock()
	return ID
}

//DoesUserExist check the map of users to see if this ID has already been entered in it
func (c *Controller) DoesUserExist(ID string) bool {

	c.mutex.Lock()
	exist := c.db.VerifyUser(ID)
	c.mutex.Unlock()
	return exist
}

func min(a int, b int) int {
	if a >= b {
		return b
	}
	return a
}
