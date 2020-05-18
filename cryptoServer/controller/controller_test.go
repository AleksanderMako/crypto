package controller

import (
	"cryptoServer/database"
	"cryptoServer/database/types"
	"cryptoServer/transactions"
	"cryptoServer/utils"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_PlaceOrder_ShouldExecuteConsecutiveMatchingTransactionForBuyerConcurrenctly(t *testing.T) {
	//Arrange
	db := database.NewDatabase()
	ids := utils.CreateWallets(db, 5)
	transactionEngine := transactions.NewTransactionEngine(db)
	controller := NewController(db, *transactionEngine)

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

	var wg sync.WaitGroup
	wg.Add(3)
	//Act

	go func(wg *sync.WaitGroup) {

		_, err := controller.PlaceOrder(*order2)
		if err != nil {
			t.Error(err.Error())
		}
		wg.Done()

	}(&wg)
	time.Sleep(1)
	go func(wg *sync.WaitGroup) {
		_, err := controller.PlaceOrder(*order3)
		if err != nil {
			t.Error(err.Error())
		}
		wg.Done()

	}(&wg)
	time.Sleep(1)
	go func(wg *sync.WaitGroup) {

		_, err := controller.PlaceOrder(*order1)
		if err != nil {
			t.Error(err.Error())
		}
		wg.Done()
	}(&wg)

	wg.Wait()
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
