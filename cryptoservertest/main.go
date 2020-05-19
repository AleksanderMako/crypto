package main

import "cryptoservertest/cli"

type Request struct {
	RequestType int    `json:"requestType"`
	Data        []byte `json:"data"`
}
type UserOrdersRequest struct {
	UserID string `json:"userID"`
}

func main() {

	cli.Run()
}
