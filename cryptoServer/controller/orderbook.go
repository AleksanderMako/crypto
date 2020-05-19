package controller

import (
	"cryptoServer/database/dbResponses"
)

// OrderBook holds top 10 lowest sell orders and top 10 highest buys
type OrderBook struct {
	LowestSellOrders []dbResponses.DBOrder `json:"lowestSellOrders"`
	HighestBuyOrders []dbResponses.DBOrder `json:"highestBuy"`
}
