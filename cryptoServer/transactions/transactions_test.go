package transactions

import (
	"cryptoServer/database"
	"cryptoServer/database/dbResponses"
	"cryptoServer/database/mocks"
	"cryptoServer/database/types"
	"cryptoServer/utils"
	"reflect"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"testing"
)

func Test_FindMatchingBuyOrder_ShouldReturnHighestApplicable(t *testing.T) {

	// Arrange
	user1 := uuid.New().String()

	user2 := uuid.New().String()
	user3 := uuid.New().String()
	user4 := uuid.New().String()
	user5 := uuid.New().String()

	newOrder := types.Order{
		CurrencyID:  "BTC",
		OrderType:   types.SellOrder,
		Price:       100,
		WalletID:    "W1",
		SumToInvest: 300,
		UserID:      user1,
	}

	candidateMatches := []dbResponses.DBOrder{
		*dbResponses.NewDBOrder("O1", types.Order{
			CurrencyID:  "BTC",
			OrderType:   types.BuyOrder,
			Price:       300,
			SumToInvest: 900,
			UserID:      user2,
		}),
		*dbResponses.NewDBOrder("O2", types.Order{
			CurrencyID:  "BTC",
			OrderType:   types.BuyOrder,
			Price:       200,
			SumToInvest: 800,
			UserID:      user3,
		}),
		*dbResponses.NewDBOrder("O3", types.Order{
			CurrencyID:  "BTC",
			OrderType:   types.BuyOrder,
			Price:       400,
			SumToInvest: 800,
			UserID:      user4,
		}),
		*dbResponses.NewDBOrder("O4", types.Order{
			CurrencyID:  "BTC",
			OrderType:   types.BuyOrder,
			Price:       120,
			SumToInvest: 360,
			UserID:      user5,
		}),
	}
	db := new(mocks.Storage)
	transactionEngine := NewTransactionEngine(db)

	// Act
	matches, found := transactionEngine.FindMatchingBuyOrder(newOrder, candidateMatches)

	//Assert
	if !found {
		t.Error("found shuould be true")
	}
	if matches[0].Order.Price != 400 {
		t.Errorf("the matching order should be the highest due to sort but got %v \n", matches[0].Order.Price)
	}

	if len(matches) < 1 {
		t.Error("should have found at least 1 match ")
	}

}
func Test_FindMatchingBuyOrder_ShouldReturnFalse(t *testing.T) {
	// Arrange

	newOrder := types.Order{
		CurrencyID:  "BTC",
		OrderType:   types.SellOrder,
		Price:       100,
		WalletID:    "W1",
		SumToInvest: 300,
	}

	candidateMatches := []dbResponses.DBOrder{
		*dbResponses.NewDBOrder("O1", types.Order{
			CurrencyID: "BTC",
			OrderType:  types.BuyOrder,
			Price:      90,
		}),
		*dbResponses.NewDBOrder("O2", types.Order{
			CurrencyID: "BTC",
			OrderType:  types.BuyOrder,
			Price:      80,
		}),
		*dbResponses.NewDBOrder("O3", types.Order{
			CurrencyID: "BTC",
			OrderType:  types.BuyOrder,
			Price:      88,
		}),
		*dbResponses.NewDBOrder("O4", types.Order{
			CurrencyID: "BTC",
			OrderType:  types.BuyOrder,
			Price:      33,
		}),
	}
	db := new(mocks.Storage)
	transactionEngine := NewTransactionEngine(db)

	// Act
	matches, found := transactionEngine.FindMatchingBuyOrder(newOrder, candidateMatches)

	//Assert
	if found {
		t.Error("found shuould not be true")
	}
	if len(matches) != 0 {
		t.Errorf("There should not be any matches but got  %v  matches \n", len(matches))
	}

}

func Test_GetBuyOrders_ShouldReturnOnlyBuyOrdersOfCertainCurrency(t *testing.T) {

	//Arrange

	const BTC = "BTC"
	const LTC = "LTC"

	db := database.NewDatabase()
	btcCurrency := types.NewCurrency(BTC)
	ltcCurrency := types.NewCurrency(LTC)

	db.Orders["o1"] = types.Order{CurrencyID: db.CreateCurrency(*btcCurrency), OrderType: types.BuyOrder}
	db.Orders["o2"] = types.Order{CurrencyID: db.CreateCurrency(*btcCurrency), OrderType: types.BuyOrder}
	db.Orders["o3"] = types.Order{CurrencyID: db.CreateCurrency(*btcCurrency), OrderType: types.SellOrder}
	db.Orders["o3"] = types.Order{CurrencyID: db.CreateCurrency(*ltcCurrency), OrderType: types.BuyOrder}

	transactionEngine := NewTransactionEngine(db)

	newOrder := types.Order{
		CurrencyID: "BTC",
		OrderType:  types.SellOrder,
		Price:      100,
		WalletID:   "W1",
	}

	// Act

	responses := transactionEngine.GetBuyOrders(newOrder)

	//Assert

	for _, response := range responses {
		if db.GetCurrency(response.Order.CurrencyID).Name != BTC || response.Order.OrderType != types.BuyOrder {
			t.Error("Expected only buy Orders for BTC ")
		}
	}

}

func Test_determineTransactionEntities(t *testing.T) {

	//Arrange

	buyer := types.Order{
		CurrencyID:  "BTC",
		OrderType:   types.BuyOrder,
		Price:       100,
		WalletID:    "W1",
		SumToInvest: 200,
	}
	seller := *dbResponses.NewDBOrder("sellerID", types.Order{
		CurrencyID:  "BTC",
		OrderType:   types.SellOrder,
		Price:       100,
		SumToInvest: 300,
		WalletID:    "W2",
	})
	db := new(mocks.Storage)
	transactionEngine := NewTransactionEngine(db)

	db.On("GetWallet", mock.Anything).Return(*types.NewWallet([]string{"BTC"}))
	db.On("GetCurrency", mock.Anything).Return(types.Currency{Ammount: 100})

	// Act
	testBuyer, testSeller := transactionEngine.determineTransactionEntities(buyer, seller)

	// Assert

	if !reflect.DeepEqual(testBuyer.Order, buyer) {
		t.Error("Wrong buyer object")
	}

	if !reflect.DeepEqual(testSeller.Order, seller.Order) {
		t.Error("Wrong seller")
	}

}
func Test_transferFunds_transferTokens_ShouldAccountForSmalletSum(t *testing.T) {

	//Arrange

	//buyer
	buyerWalletID := "BuyerWallet"
	currencyName := "BTC"
	buyersCurrencyID := buyerWalletID + currencyName
	buyersCurrency := types.NewCurrency(currencyName)
	buyerWallet := types.NewWallet([]string{buyersCurrencyID})
	buyerOrder := types.NewOrder(buyerWalletID, buyersCurrencyID, 100, types.BuyOrder, 200)
	buyer := types.NewTransactionEntity(*buyerWallet, *buyerOrder, types.EmptyString)

	//seller
	sellerWalletID := "SellerWallet"
	sellersCurrencyID := buyerWalletID + currencyName
	sellersCurrency := types.NewCurrency(currencyName)
	sellerWallet := types.NewWallet([]string{sellersCurrencyID})
	sellerOrder := types.NewOrder(sellerWalletID, sellersCurrencyID, 100, types.SellOrder, 300)
	seller := types.NewTransactionEntity(*sellerWallet, *sellerOrder, types.EmptyString)

	db := new(mocks.Storage)
	transactionEngine := NewTransactionEngine(db)

	//Act

	transactionEngine.transferFunds(buyer, seller)
	transactionEngine.transferTokens(buyersCurrency, sellersCurrency, 200, 100)

	//Assert

	if buyer.Wallet.Balance != 800 {
		t.Errorf("Buyers wallet funds should be 1000 - %v = 800 but balanace is %v \n", 200, buyer.Wallet.Balance)
	}
	if seller.Wallet.Balance != 1200 {
		t.Errorf("Sellers wallet funds should be 100 + %v = 1200 but balanace is %v \n", 200, seller.Wallet.Balance)
	}
	if buyersCurrency.Ammount != 1002 {
		t.Errorf("buyers Currency  amount should be 1000 + %v = 1002 but balanace is %v \n", 2, buyersCurrency.Ammount)
	}
	if sellersCurrency.Ammount != 998 {
		t.Errorf("sellers Currency   amount should be 1000 - %v = 998 but balanace is %v \n", 2, sellersCurrency.Ammount)
	}

}

func Test_ExecuteTransfer_ShouldCorrectlyUpdateAllBalancesWhenBuyerLesThanSeller(t *testing.T) {

	//Arrange

	db := database.NewDatabase()
	transactionEngine := NewTransactionEngine(db)

	WalletID := "BuyerWallet"
	currencyName := "BTC"
	currencyID := WalletID + currencyName
	currency := types.NewCurrency(currencyName)
	wallet := types.NewWallet([]string{currencyID})
	order := types.NewOrder(WalletID, currencyID, 100, types.BuyOrder, 200)

	matchingOrderWalletID := "matchOrderWallet"
	matchingOfferCurrencyID := matchingOrderWalletID + currencyName
	matchinOfferCurrency := types.NewCurrency(currencyName)
	matchinOfferWallet := types.NewWallet([]string{matchingOfferCurrencyID})
	matchingOrder := types.NewOrder(matchingOrderWalletID, matchingOfferCurrencyID, 100, types.SellOrder, 400)

	matchingOffer := dbResponses.NewDBOrder("matchingOffer", *matchingOrder)

	db.Currencies[currencyID] = *currency
	db.Currencies[matchingOfferCurrencyID] = *matchinOfferCurrency
	db.Wallets[matchingOrderWalletID] = *matchinOfferWallet
	db.Wallets[WalletID] = *wallet
	db.Orders["matchingOffer"] = *matchingOrder

	// Act

	transactionEngine.ExecuteTransfer(order, []dbResponses.DBOrder{*matchingOffer})

	// Assert

	if db.Orders["matchingOffer"].SumToInvest != 200 {
		t.Errorf("Matching offer is seller and the seller should have sold 200,remaining selling sum should be 200 but it is %v",
			matchingOffer.Order.SumToInvest)
	}

	if db.Wallets[WalletID].Balance != 800 {
		t.Error("Wrong balance for order ")
	}

	if db.Wallets[matchingOrderWalletID].Balance != 1200 {
		t.Error("Wrong balance for matching offer")
	}

	if db.Currencies[currencyID].Ammount != 1002 {
		t.Errorf("Wrong amount of coin in buyers balance should be 1020 bu it is %v", db.Currencies[currencyID].Ammount)
	}

	if db.Currencies[matchingOfferCurrencyID].Ammount != 998 {
		t.Errorf("Wrong amount of coin in buyers balance should be 980 bu it is %v", db.Currencies[matchingOfferCurrencyID].Ammount)
	}
}

func Test_ExecuteTransfer_ShouldCorrectlyUpdateAllBalancesWhenBuyerEqualSeller(t *testing.T) {
	//Arrange

	db := database.NewDatabase()
	transactionEngine := NewTransactionEngine(db)

	WalletID := "BuyerWallet"
	currencyName := "BTC"
	currencyID := WalletID + currencyName
	currency := types.NewCurrency(currencyName)
	wallet := types.NewWallet([]string{currencyID})
	order := types.NewOrder(WalletID, currencyID, 100, types.BuyOrder, 200)

	matchingOrderWalletID := "matchOrderWallet"
	matchingOfferCurrencyID := matchingOrderWalletID + currencyName
	matchinOfferCurrency := types.NewCurrency(currencyName)
	matchinOfferWallet := types.NewWallet([]string{matchingOfferCurrencyID})
	matchingOrder := types.NewOrder(matchingOrderWalletID, matchingOfferCurrencyID, 100, types.SellOrder, 200)

	matchingOffer := dbResponses.NewDBOrder("matchingOffer", *matchingOrder)

	db.Currencies[currencyID] = *currency
	db.Currencies[matchingOfferCurrencyID] = *matchinOfferCurrency
	db.Wallets[matchingOrderWalletID] = *matchinOfferWallet
	db.Wallets[WalletID] = *wallet
	db.Orders["matchingOffer"] = *matchingOrder

	// Act

	transactionEngine.ExecuteTransfer(order, []dbResponses.DBOrder{*matchingOffer})

	// Assert

	if db.Wallets[WalletID].Balance != 800 {
		t.Error("Wrong balance for order ")
	}

	if db.Wallets[matchingOrderWalletID].Balance != 1200 {
		t.Error("Wrong balance for matching offer")
	}

	if db.Currencies[currencyID].Ammount != 1002 {
		t.Errorf("Wrong amount of coin in buyers balance should be 1020 bu it is %v", db.Currencies[currencyID].Ammount)
	}

	if db.Currencies[matchingOfferCurrencyID].Ammount != 998 {
		t.Errorf("Wrong amount of coin in buyers balance should be 980 bu it is %v", db.Currencies[matchingOfferCurrencyID].Ammount)
	}

	if _, ok := db.Orders["matchingOffer"]; ok {
		t.Error("The matching offer should not be in the db anymore ")
	}
}

func Test_PlaceOrder_ShouldCreateOrderWhenNoMatchingOfferExists(t *testing.T) {

	//Arrange
	db := database.NewDatabase()
	ids := utils.CreateWallets(db, 2)

	transactionEngine := NewTransactionEngine(db)

	//Act
	wallet := db.Wallets[ids[0]]
	order := types.NewOrder(ids[0], wallet.Currencies[0], 100, types.SellOrder, 200)
	transactionEngine.PlaceOrder(*order)

	//Assert
	if len(db.Orders) != 1 {
		t.Errorf("Wrong number of wallets created  %v", len(db.Orders))
	}

}

func Test_PlaceOrder_ShouldExecuteConsecutiveMatchingTransaction(t *testing.T) {
	//Arrange
	db := database.NewDatabase()
	ids := utils.CreateWallets(db, 5)

	transactionEngine := NewTransactionEngine(db)

	id1 := ids[0]
	id2 := ids[1]
	id3 := ids[2]
	user1 := uuid.New().String()
	user2 := uuid.New().String()
	user3 := uuid.New().String()

	wallet1 := db.Wallets[id1]
	wallet2 := db.Wallets[id2]
	wallet3 := db.Wallets[id3]

	// order 1 sell at 100 10 BTC
	order1 := types.NewOrder(id1, wallet1.Currencies[0], 100, types.SellOrder, 1000)
	order1.UserID = user1

	// order 2 buy at 200 3 coins
	order2 := types.NewOrder(id2, wallet2.Currencies[0], 200, types.BuyOrder, 600)
	order2.UserID = user2

	// order 3 buy at 150 4 coins
	order3 := types.NewOrder(id3, wallet3.Currencies[0], 150, types.BuyOrder, 600)
	order3.UserID = user3

	//Act

	_, err := transactionEngine.PlaceOrder(*order2)
	_, err = transactionEngine.PlaceOrder(*order3)
	_, err = transactionEngine.PlaceOrder(*order1)
	if err != nil {
		t.Error(err.Error())
	}

	//Assert
	wallet1 = db.Wallets[id1]
	wallet2 = db.Wallets[id2]
	wallet3 = db.Wallets[id3]

	currency1 := db.Currencies[wallet1.Currencies[0]]
	currency2 := db.Currencies[wallet2.Currencies[0]]
	currency3 := db.Currencies[wallet3.Currencies[0]]

	if wallet1.Balance != 2000 {
		t.Errorf("Wallet 1 shoud be 2000 but it is %v", wallet1.Balance)
	}

	if wallet2.Balance != 400 {
		t.Errorf("Wallet 2 balance should be 400 after spending 600 but balance is %v", wallet2.Balance)
	}

	if wallet3.Balance != 600 {
		t.Errorf("Wallet 3 balance should be 400 after spending 600 but balance is %v", wallet3.Balance)
	}

	if currency1.Ammount != 990 {
		t.Errorf("Currency 1 should be 990 bu it is %v", currency1.Ammount)
	}

	if currency2.Ammount != 1006 {
		t.Errorf("Currency 2 should be 1006 bu it is %v", currency2.Ammount)
	}

	if currency3.Ammount != 1004 {
		t.Errorf("Currency 3 should be 1004 bu it is %v", currency2.Ammount)
	}

}
func Test_PlaceOrder_ShouldFailWhenUserDoesNotHaveEnoughFunds(t *testing.T) {
	//Arrange
	db := database.NewDatabase()
	ids := utils.CreateWallets(db, 5)

	transactionEngine := NewTransactionEngine(db)

	id1 := ids[0]
	id2 := ids[1]
	id3 := ids[2]
	user1 := uuid.New().String()
	user2 := uuid.New().String()
	user3 := uuid.New().String()

	wallet1 := db.Wallets[id1]
	wallet2 := db.Wallets[id2]
	wallet3 := db.Wallets[id3]

	order1 := types.NewOrder(id1, wallet1.Currencies[0], 300, types.BuyOrder, 10000)
	order1.UserID = user1

	order2 := types.NewOrder(id2, wallet2.Currencies[0], 200, types.SellOrder, 600)
	order2.UserID = user2

	order3 := types.NewOrder(id3, wallet3.Currencies[0], 150, types.SellOrder, 600)
	order3.UserID = user3

	//Act

	_, err := transactionEngine.PlaceOrder(*order2)
	_, err = transactionEngine.PlaceOrder(*order3)
	_, err = transactionEngine.PlaceOrder(*order1)

	// Assert

	errormsg := "Buyer does not have enough money"
	if err == nil {
		t.Error("Place Order method should return error ")
	}

	if err.Error() != errormsg {
		t.Errorf("error message should have been %v but it was %v", errormsg, err.Error())
	}
}
func Test_PlaceOrder_ShouldExecuteConsecutiveMatchingTransactionForBuyer(t *testing.T) {
	//Arrange
	db := database.NewDatabase()
	ids := utils.CreateWallets(db, 5)

	transactionEngine := NewTransactionEngine(db)

	id1 := ids[0]
	id2 := ids[1]
	id3 := ids[2]
	user1 := uuid.New().String()
	user2 := uuid.New().String()
	user3 := uuid.New().String()

	wallet1 := db.Wallets[id1]
	wallet2 := db.Wallets[id2]
	wallet3 := db.Wallets[id3]

	order1 := types.NewOrder(id1, wallet1.Currencies[0], 120, types.BuyOrder, 480)
	order1.UserID = user1

	order2 := types.NewOrder(id2, wallet2.Currencies[0], 100, types.SellOrder, 300)
	order2.UserID = user2

	order3 := types.NewOrder(id3, wallet3.Currencies[0], 120, types.SellOrder, 300)
	order3.UserID = user3

	//Act

	_, err := transactionEngine.PlaceOrder(*order2)
	_, err = transactionEngine.PlaceOrder(*order3)
	_, err = transactionEngine.PlaceOrder(*order1)
	if err != nil {
		t.Error(err.Error())
	}

	//Assert
	wallet1 = db.Wallets[id1]
	wallet2 = db.Wallets[id2]
	wallet3 = db.Wallets[id3]

	currency1 := db.Currencies[wallet1.Currencies[0]]
	currency2 := db.Currencies[wallet2.Currencies[0]]
	currency3 := db.Currencies[wallet3.Currencies[0]]

	if wallet1.Balance != 520 {
		t.Errorf("Wallet 1 shoud be 520 but it is %v", wallet1.Balance)
	}

	if wallet2.Balance != 1300 {
		t.Errorf("Wallet 2 balance should be 1300 after selling 300 but balance is %v", wallet2.Balance)
	}

	if wallet3.Balance != 1180 {
		t.Errorf("Wallet 3 balance should be 1180 after spending 180 but balance is %v", wallet3.Balance)
	}

	if currency1.Ammount != 1003 {
		t.Errorf("Currency 1 should be 1003 bu it is %v", currency1.Ammount)
	}

	if currency2.Ammount != 998 {
		t.Errorf("Currency 2 should be 998 bu it is %v", currency2.Ammount)
	}

	if currency3.Ammount != 999 {
		t.Errorf("Currency 3 should be 999 bu it is %v", currency2.Ammount)
	}

}
