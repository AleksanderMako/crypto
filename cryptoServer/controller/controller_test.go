package controller

import (
	"cryptoServer/database"
	"cryptoServer/database/types"
	"cryptoServer/transactions"
	"cryptoServer/utils"
	"encoding/json"
	"sort"
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
	time.Sleep(1 * time.Millisecond)
	go func(wg *sync.WaitGroup) {
		_, err := controller.PlaceOrder(*order3)
		if err != nil {
			t.Error(err.Error())
		}
		wg.Done()

	}(&wg)
	time.Sleep(1 * time.Millisecond)
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
func Test_ListOrderBook(t *testing.T) {

	//Arrange

	db := database.NewDatabase()
	transactionEngine := transactions.NewTransactionEngine(db)
	controller := NewController(db, *transactionEngine)

	order1 := *types.NewOrder("id1", "c1", 120, types.BuyOrder, 480)
	order2 := *types.NewOrder("id2", "c2", 200, types.BuyOrder, 380)
	order3 := *types.NewOrder("id3", "c3", 220, types.BuyOrder, 660)
	order4 := *types.NewOrder("id4", "c4", 500, types.BuyOrder, 1000)
	order5 := *types.NewOrder("id5", "c5", 340, types.BuyOrder, 480)
	order6 := *types.NewOrder("id6", "c6", 600, types.BuyOrder, 700)
	order7 := *types.NewOrder("id7", "c7", 230, types.BuyOrder, 700)
	order8 := *types.NewOrder("id8", "c8", 330, types.BuyOrder, 650)
	order9 := *types.NewOrder("id9", "c9", 300, types.BuyOrder, 720)
	order10 := *types.NewOrder("id10", "c10", 400, types.BuyOrder, 820)
	order11 := *types.NewOrder("id10", "c10", 400, types.BuyOrder, 450)
	buyers := []types.Order{order1, order2, order3, order4, order5, order6, order7, order8, order9, order10, order11}
	sort.Slice(buyers, func(i, j int) bool {
		if buyers[i].Price > buyers[j].Price {
			return true
		}
		return false
	})

	sellorder1 := *types.NewOrder("id1", "c1", 120, types.SellOrder, 480)
	sellorder2 := *types.NewOrder("id2", "c2", 200, types.SellOrder, 380)
	sellorder3 := *types.NewOrder("id3", "c3", 220, types.SellOrder, 660)
	sellorder4 := *types.NewOrder("id4", "c4", 500, types.SellOrder, 1000)
	sellorder5 := *types.NewOrder("id5", "c5", 340, types.SellOrder, 480)
	sellorder6 := *types.NewOrder("id6", "c6", 600, types.SellOrder, 700)
	sellorder7 := *types.NewOrder("id7", "c7", 230, types.SellOrder, 700)
	sellorder8 := *types.NewOrder("id8", "c8", 330, types.SellOrder, 650)
	sellorder9 := *types.NewOrder("id9", "c9", 300, types.SellOrder, 720)
	sellorder10 := *types.NewOrder("id10", "c10", 400, types.SellOrder, 820)
	sellorder11 := *types.NewOrder("id10", "c10", 700, types.SellOrder, 1400)
	sellers := []types.Order{sellorder1, sellorder2, sellorder3, sellorder4, sellorder5, sellorder6,
		sellorder7, sellorder8, sellorder9, sellorder10, sellorder11}

	sort.Slice(sellers, func(i, j int) bool {
		if sellers[i].Price < sellers[j].Price {
			return true
		}
		return false
	})
	db.CreateOrder(order1)
	db.CreateOrder(order2)
	db.CreateOrder(order3)
	db.CreateOrder(order4)
	db.CreateOrder(order5)
	db.CreateOrder(order6)
	db.CreateOrder(order7)
	db.CreateOrder(order8)
	db.CreateOrder(order9)
	db.CreateOrder(order10)
	db.CreateOrder(order11)

	db.CreateOrder(sellorder1)
	db.CreateOrder(sellorder2)
	db.CreateOrder(sellorder3)
	db.CreateOrder(sellorder4)
	db.CreateOrder(sellorder5)
	db.CreateOrder(sellorder6)
	db.CreateOrder(sellorder7)
	db.CreateOrder(sellorder8)
	db.CreateOrder(sellorder9)
	db.CreateOrder(sellorder10)
	db.CreateOrder(sellorder11)

	// Act

	response, err := controller.ListOrderBook()
	if err != nil {
		t.Error(err.Error())
	}
	var orderBook OrderBook

	err = json.Unmarshal(response, &orderBook)
	if err != nil {
		t.Error(err.Error())
	}

	// Assert

	if len(orderBook.HighestBuyOrders) != 10 {
		t.Errorf("Order book should have 10 entires in the highest buying offers but has %v \n", len(orderBook.HighestBuyOrders))
	}

	if len(orderBook.LowestSellOrders) != 10 {
		t.Errorf("Order book should have 10 entires in the lowest selling offers but has %v \n", len(orderBook.LowestSellOrders))
	}

	for i := 0; i < 10; i++ {

		if orderBook.HighestBuyOrders[i].Order.Price != buyers[i].Price {
			t.Error("sorted buyer should match the expected values ")
		}
	}
	for i := 0; i < 10; i++ {

		if orderBook.LowestSellOrders[i].Order.Price != sellers[i].Price {
			t.Error("sorted seller should match the expected values ")
		}
	}

}
