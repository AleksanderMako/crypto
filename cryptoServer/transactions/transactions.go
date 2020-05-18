package transactions

import (
	"cryptoServer/database"
	"cryptoServer/database/dbResponses"
	"cryptoServer/database/types"
	"errors"
	"math"
	"sort"
)

// TransactionEngine encapsulates transaction logic for buy and sell
type TransactionEngine struct {
	db database.Storage
}

// NewTransactionEngine returns a pointer to new  TransactionEngine
func NewTransactionEngine(db database.Storage) *TransactionEngine {

	return &TransactionEngine{
		db: db,
	}
}

// PlaceOrder attempts to execute either a buy or a sell Limit Order
func (t *TransactionEngine) PlaceOrder(order types.Order) (string, error) {
	if order.OrderType == types.SellOrder {

		status, err := t.TryPlaceSellOrder(order)
		if err != nil {
			return "", err

		}
		return status, nil
	}

	status, err := t.TryPlaceBuyOrder(order)
	if err != nil {
		return "", err
	}
	return status, nil

}

// TryPlaceSellOrder tries to find matching buyer otherwise places this order with the rest
func (t *TransactionEngine) TryPlaceSellOrder(newOrder types.Order) (string, error) {
	// get buy orders
	buyOrders := t.GetBuyOrders(newOrder)
	matchingOffers, found := t.FindMatchingBuyOrder(newOrder, buyOrders)
	if found {
		// execute transaction

		err := t.ExecuteTransfer(&newOrder, matchingOffers)
		if err != nil {
			return "", err
		}
	} else {
		// add sell order to exisitng order
		t.db.CreateOrder(newOrder)
		return types.OrderHold, nil
	}
	return types.OrderExecutedSuccessfully, nil

}

// TryPlaceBuyOrder tries to find matching seller otherwiser places this order with the rest
func (t *TransactionEngine) TryPlaceBuyOrder(newOrder types.Order) (string, error) {

	//get sell orders
	sellOrders := t.GetSellOrders(newOrder)
	mathchingOffers, found := t.FindMatchingSellOrder(newOrder, sellOrders)

	if found {

		err := t.ExecuteTransfer(&newOrder, mathchingOffers)
		if err != nil {
			return "", err
		}

	} else {
		t.db.CreateOrder(newOrder)
		return types.OrderHold, nil

	}
	return types.OrderExecutedSuccessfully, nil
}

// GetBuyOrders returns a list of buyer DBOrder structs which have an order struct and the corresponding ID for the map
func (t *TransactionEngine) GetBuyOrders(newOrder types.Order) []dbResponses.DBOrder {
	buyOrders := t.db.GetOrders(func(order types.Order) bool {
		// TODO: check this cID to maybe be ctype
		if order.OrderType == types.BuyOrder && t.db.GetCurrency(order.CurrencyID).Name == t.db.GetCurrency(newOrder.CurrencyID).Name {
			return true
		}
		return false
	})
	return buyOrders
}

// GetSellOrders returns a list of  seller DBOrder structs which have an order struct and the corresponding ID for the map
func (t *TransactionEngine) GetSellOrders(newOrder types.Order) []dbResponses.DBOrder {

	sellOrders := t.db.GetOrders(func(order types.Order) bool {
		if order.OrderType == types.SellOrder && t.db.GetCurrency(order.CurrencyID).Name == t.db.GetCurrency(newOrder.CurrencyID).Name {
			return true
		}
		return false
	})
	return sellOrders
}

// FindMatchingBuyOrder search for a buyer matching the sellers price given the seller
func (t *TransactionEngine) FindMatchingBuyOrder(
	newOrder types.Order,
	candidateMatches []dbResponses.DBOrder) (matches []dbResponses.DBOrder, found bool) {

	if len(candidateMatches) < 1 {
		return nil, false
	}

	sort.Slice(candidateMatches, func(i, j int) bool {
		return candidateMatches[i].Order.Price > candidateMatches[j].Order.Price
	})

	if candidateMatches[0].Order.Price < newOrder.Price {
		return nil, false
	}

	return t.FindMatchingOrders(newOrder, candidateMatches, func(orderPrice float64, candidatePrice float64) bool {
		if candidatePrice >= orderPrice {
			return true
		}
		return false
	})
}

// FindMatchingSellOrder search for a seller matching the buyers price given the buyer
func (t *TransactionEngine) FindMatchingSellOrder(
	newOrder types.Order,
	candidateMatches []dbResponses.DBOrder) (matches []dbResponses.DBOrder, found bool) {

	if len(candidateMatches) < 1 {
		return nil, false
	}

	sort.Slice(candidateMatches, func(i, j int) bool {
		return candidateMatches[i].Order.Price < candidateMatches[j].Order.Price
	})

	if candidateMatches[0].Order.Price > newOrder.Price {
		return nil, false
	}

	return t.FindMatchingOrders(newOrder, candidateMatches, func(orderPrice float64, candidatePrice float64) bool {

		if candidatePrice <= orderPrice {
			return true
		}
		return false
	})
}

func (t *TransactionEngine) FindMatchingOrders(
	newOrder types.Order,
	candidateMatches []dbResponses.DBOrder,
	compare func(orderPrice float64, candidatePrice float64) bool) (matches []dbResponses.DBOrder, found bool) {

	newOrderCoins := int(newOrder.SumToInvest / newOrder.Price)
	matches = []dbResponses.DBOrder{}
	for _, order := range candidateMatches {
		if compare(newOrder.Price, order.Order.Price) && newOrderCoins != 0 && newOrder.UserID != order.Order.UserID {
			matches = append(matches, order)
			newOrderCoins -= int(order.Order.SumToInvest / newOrder.Price)

		}
	}
	if len(matches) > 0 {
		return matches, true
	}
	return nil, false
}

func (t *TransactionEngine) ExecuteTransfer(order *types.Order, matchingOffers []dbResponses.DBOrder) error {

	for _, matchingOffer := range matchingOffers {

		if order.SumToInvest == 0 {
			return nil
		}

		exchangeRate := order.Price
		buyer, seller := t.determineTransactionEntities(*order, matchingOffer)
		buyCurrencty, sellCurrency := t.determineTransactionCurrencies(*order, matchingOffer)

		err := t.verifyTransactionEconomics(buyer, seller, sellCurrency)
		if err != nil {
			return err
		}

		t.transferFunds(buyer, seller)

		if buyer.Order.SumToInvest == seller.Order.SumToInvest {

			t.transferTokens(buyCurrencty, sellCurrency, buyer.Order.SumToInvest, exchangeRate)
			t.db.DeleteOrder(maxIDLen(buyer.OrderID, seller.OrderID))
			order.SumToInvest = 0
		}

		if buyer.Order.SumToInvest < seller.Order.SumToInvest {

			t.transferTokens(buyCurrencty, sellCurrency, buyer.Order.SumToInvest, exchangeRate)
			seller.Order.SumToInvest -= buyer.Order.SumToInvest
			buyer.Order.SumToInvest -= buyer.Order.SumToInvest

			if buyer.OrderID != types.EmptyString {
				// just update seller
				t.db.DeleteOrder(buyer.OrderID)
				*order = seller.Order
			} else {
				*order = buyer.Order
			}
			t.db.UpdateOrder(seller.OrderID, seller.Order)

			// maybe recur here
		}

		if buyer.Order.SumToInvest > seller.Order.SumToInvest {
			t.transferTokens(buyCurrencty, sellCurrency, seller.Order.SumToInvest, exchangeRate)

			buyer.Order.SumToInvest -= seller.Order.SumToInvest
			seller.Order.SumToInvest -= seller.Order.SumToInvest
			if buyer.OrderID == types.EmptyString {
				// create new order for buyer
				// delete selling order
				t.db.CreateOrder(buyer.Order)
				t.db.DeleteOrder(seller.OrderID)
				*order = buyer.Order

			} else {
				// just update existing buying order
				t.db.UpdateOrder(buyer.OrderID, buyer.Order)
				*order = seller.Order
			}
		}

		// update wallets and currencies
		t.db.UpdateWallet(buyer.Order.WalletID, buyer.Wallet)
		t.db.UpdateWallet(seller.Order.WalletID, seller.Wallet)

		t.db.UpdateCurrency(seller.Order.CurrencyID, *sellCurrency)
		t.db.UpdateCurrency(buyer.Order.CurrencyID, *buyCurrencty)
	}
	return nil
}

func (t *TransactionEngine) determineTransactionEntities(order types.Order, matchingOffer dbResponses.DBOrder) (*types.TransactionEntity, *types.TransactionEntity) {

	if order.OrderType == types.SellOrder {

		seller := types.NewTransactionEntity(t.db.GetWallet(order.WalletID), order, types.EmptyString)
		seller.NewOrder = true

		buyer := types.NewTransactionEntity(t.db.GetWallet(matchingOffer.Order.WalletID), matchingOffer.Order, matchingOffer.ID)
		buyer.NewOrder = false
		return buyer, seller
	}

	seller := types.NewTransactionEntity(t.db.GetWallet(matchingOffer.Order.WalletID), matchingOffer.Order, matchingOffer.ID)
	seller.NewOrder = false

	buyer := types.NewTransactionEntity(t.db.GetWallet(order.WalletID), order, types.EmptyString)
	buyer.NewOrder = true
	return buyer, seller

}
func (t *TransactionEngine) determineTransactionCurrencies(order types.Order, matchingOffer dbResponses.DBOrder) (*types.Currency, *types.Currency) {
	if order.OrderType == types.SellOrder {

		buyCurrencty := t.db.GetCurrency(matchingOffer.Order.CurrencyID)
		sellCurrency := t.db.GetCurrency(order.CurrencyID)
		return &buyCurrencty, &sellCurrency
	}
	buyCurrencty := t.db.GetCurrency(order.CurrencyID)
	sellCurrency := t.db.GetCurrency(matchingOffer.Order.CurrencyID)
	return &buyCurrencty, &sellCurrency
}

func (t *TransactionEngine) transferFunds(buyer *types.TransactionEntity, seller *types.TransactionEntity) {
	buyer.Wallet.Balance -= math.Min(buyer.Order.SumToInvest, seller.Order.SumToInvest)
	seller.Wallet.Balance += math.Min(buyer.Order.SumToInvest, seller.Order.SumToInvest)
}

func (t *TransactionEngine) transferTokens(buyCurrencty *types.Currency, sellCurrency *types.Currency, sumToInvest float64, exchangeRate float64) {
	sellCurrency.Ammount -= int(sumToInvest / exchangeRate)
	buyCurrencty.Ammount += int(sumToInvest / exchangeRate)
}

func (t *TransactionEngine) doesBuyerHaveFunds(buyer *types.TransactionEntity) bool {

	if buyer.Wallet.Balance < buyer.Order.SumToInvest {
		return false
	}
	return true
}

func (t *TransactionEngine) doesSellerHaveCoins(seller *types.TransactionEntity, coins int) bool {

	if coins < int(seller.Order.SumToInvest/types.BTCXchangeRate) {
		return false
	}
	return true
}

func maxIDLen(a, b string) string {

	if len(a) > len(b) {
		return a
	}
	return b
}

func (t *TransactionEngine) verifyTransactionEconomics(buyer *types.TransactionEntity,
	seller *types.TransactionEntity, sellCurrency *types.Currency) error {

	if !t.doesBuyerHaveFunds(buyer) {
		return errors.New("Buyer does not have enough money")
	}

	if !t.doesSellerHaveCoins(seller, sellCurrency.Ammount) {
		return errors.New("Seller does not have enough coins ")
	}
	return nil
}
