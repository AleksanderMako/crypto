package router

import (
	"cryptoServer/controller"
	"cryptoServer/database/types"
	requestModels "cryptoServer/reqeuestModels"
	"encoding/json"
	"fmt"
)

// Router encapsulates the logic needed to route the requests
type Router struct {
	controller controller.Controller
}

func NewRouter(controller controller.Controller) *Router {

	return &Router{
		controller: controller,
	}
}

func (r *Router) HandleRequest(data []byte, route func(req requestModels.Request) ([]byte, error)) ([]byte, error) {

	var request requestModels.Request
	err := json.Unmarshal(data, &request)
	if err != nil {
		//return
	}
	fmt.Printf("Request datais %v", string(request.Data))
	response, err := route(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// RouteRequest determines which controller function to call given a request object
func (r *Router) RouteRequest(req requestModels.Request) ([]byte, error) {

	var err error
	var response []byte
	switch req.RequestType {

	case types.ListWalletBalances:
		response, err = r.controller.ListWalletBalances()
		if err != nil {
			return nil, err
		}
		return response, nil

	case types.ListYourOrders:
		var payload requestModels.UserOrdersRequest
		err = json.Unmarshal(req.Data, &payload)
		if err != nil {
			return nil, err
		}
		response, err = r.controller.ListOrdersByUser(payload.UserID)
		if err != nil {
			return nil, err
		}

	case types.CancelOrder:
		var payload requestModels.UserOrdersRequest
		err = json.Unmarshal(req.Data, &payload)
		if err != nil {
			return nil, err
		}
		r.controller.CancelOrder(payload.UserID)
		return []byte("successfully deleted wallet"), nil

	case types.PlaceOrder:
		var payload types.Order
		err = json.Unmarshal(req.Data, &payload)
		if err != nil {
			return nil, err
		}
		response, err = r.controller.PlaceOrder(payload)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return response, nil
}
