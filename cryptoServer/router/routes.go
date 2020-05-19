package router

import (
	"cryptoServer/controller"
	"cryptoServer/database/types"
	requestModels "cryptoServer/reqeuestModels"
	"encoding/json"
	"errors"
	"fmt"
)

// Router encapsulates the logic needed to route the requests
type Router struct {
	controller *controller.Controller
}

func NewRouter(controller *controller.Controller) *Router {

	return &Router{
		controller: controller,
	}
}

// HandleRequest receives the request and then unmrshalles it into a Request struct
func (r *Router) HandleRequest(data []byte,
	middleware func(req requestModels.Request, route func(req requestModels.Request) ([]byte, error)) ([]byte, error)) ([]byte, error) {

	var request requestModels.Request
	err := json.Unmarshal(data, &request)
	if err != nil {
		//return
	}
	fmt.Printf("Request datais %v", string(request.Data))
	response, err := middleware(request, r.RouteRequest)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *Router) IdentifyUser(req requestModels.Request, route func(req requestModels.Request) ([]byte, error)) ([]byte, error) {

	if len(req.UserID) < 1 && req.RequestType != types.Register {
		return nil, errors.New("User ID is missing from the request ")
	}
	doesUserExist := r.controller.DoesUserExist(req.UserID)
	if doesUserExist == false && req.RequestType != types.Register {
		return nil, errors.New("This is ID does not exist  , please hit the register endpoint to get an ID ")
	}
	response, err := route(req)
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

		response, err = r.controller.ListOrdersByUser(req.UserID)
		if err != nil {
			return nil, err
		}

	case types.CancelOrder:
		var payload requestModels.CancelOrder
		err = json.Unmarshal(req.Data, &payload)
		if err != nil {
			return nil, err
		}
		r.controller.CancelOrder(payload.OrderID)
		return []byte("successfully cancelled the order"), nil

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

	case types.ListOrderBook:
		response, err = r.controller.ListOrderBook()
		if err != nil {
			return nil, err
		}
		return response, nil
	case types.Register:
		ID := r.controller.RegisterUser()
		return []byte(fmt.Sprintf("This is your ID: %v please include it in subsequent requests", ID)), nil
	}
	return response, nil
}
