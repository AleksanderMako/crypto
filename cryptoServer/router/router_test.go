package router

import (
	"cryptoServer/controller"
	"cryptoServer/database"
	"cryptoServer/database/dbResponses"
	"cryptoServer/database/types"
	requestModels "cryptoServer/reqeuestModels"
	"cryptoServer/transactions"
	"cryptoServer/utils"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func Test_HandleRequest_ForListOfOrders(t *testing.T) {

	// Arrange
	db := database.NewDatabase()
	te := transactions.NewTransactionEngine(db)
	c := controller.NewController(db, *te)
	r := NewRouter(c)

	user := uuid.New().String()

	order1 := types.NewOrder(uuid.New().String(), uuid.New().String(), 100, types.SellOrder, 1000)
	order1.UserID = user

	order2 := types.NewOrder(uuid.New().String(), uuid.New().String(), 100, types.SellOrder, 1000)
	order2.UserID = user

	order3 := types.NewOrder(uuid.New().String(), uuid.New().String(), 100, types.SellOrder, 1000)
	order3.UserID = user

	db.CreateOrder(*order1)
	db.CreateOrder(*order2)
	db.CreateOrder(*order3)
	db.RegisterUser(user)

	// Act

	userIDRequest := requestModels.UserOrdersRequest{
		UserID: user,
	}
	data, err := json.Marshal(userIDRequest)
	if err != nil {
		t.Error(" Marshalling of userIDRequest failed ")
	}
	request := requestModels.Request{
		RequestType: types.ListYourOrders,
		Data:        data,
		UserID:      user,
	}
	req, err := json.Marshal(request)
	if err != nil {
		t.Error("Marshalling of request failed")
	}

	response, err := r.HandleRequest(req, r.IdentifyUser)

	//Assert

	if err != nil {
		t.Error(err.Error())
	}

	var orders []dbResponses.DBOrder

	err = json.Unmarshal(response, &orders)
	if err != nil {
		t.Error(err.Error())
	}

	if len(orders) != 3 {
		t.Error("Should have been 3 orders returned ")
	}

}

func Test_HandleRequest_ForListOfWallets(t *testing.T) {

	// Arrange
	db := database.NewDatabase()
	te := transactions.NewTransactionEngine(db)
	c := controller.NewController(db, *te)
	r := NewRouter(c)

	ids := utils.CreateWallets(db, 3)

	id1 := ids[0]
	id2 := ids[1]
	id3 := ids[2]

	user := uuid.New().String()
	db.RegisterUser(user)

	wallet1 := db.Wallets[id1]
	wallet2 := db.Wallets[id2]
	wallet3 := db.Wallets[id3]

	wallet1Currencies := []types.Currency{}
	for i := 0; i < 3; i++ {
		wallet1Currencies = append(wallet1Currencies, db.Currencies[wallet1.Currencies[i]])
	}

	wallet2Currencies := []types.Currency{}
	for i := 0; i < 3; i++ {
		wallet2Currencies = append(wallet2Currencies, db.Currencies[wallet2.Currencies[i]])
	}

	wallet3Currencies := []types.Currency{}
	for i := 0; i < 3; i++ {
		wallet3Currencies = append(wallet3Currencies, db.Currencies[wallet3.Currencies[i]])
	}

	// Act
	request := requestModels.Request{
		RequestType: types.ListWalletBalances,
		Data:        nil,
		UserID:      user,
	}
	req, err := json.Marshal(request)
	if err != nil {
		t.Error("Marshalling of request failed")
	}
	response, err := r.HandleRequest(req, r.IdentifyUser)

	//Assert

	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(string(response))
	var wallets []dbResponses.DBWalletCurrency

	err = json.Unmarshal(response, &wallets)
	if err != nil {
		t.Error(err.Error())
	}

	if len(wallets) != 3 {
		t.Error("Should have been 3 wallets returned")
	}

}
