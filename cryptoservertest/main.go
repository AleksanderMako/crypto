package main

import (
	"bufio"
	"cryptoServer/database/types"
	requestModels "cryptoServer/reqeuestModels"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Request struct {
	RequestType int    `json:"requestType"`
	Data        []byte `json:"data"`
}
type UserOrdersRequest struct {
	UserID string `json:"userID"`
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Failed connecting" + err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	// make request
	request := requestModels.Request{
		RequestType: types.ListWalletBalances,
		Data:        nil,
	}
	req, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Marshalling of request failed" + err.Error())
	}

	fmt.Println(string(req))

	_, err = conn.Write(append(req, '\r'))
	if err != nil {
		fmt.Println("Error writing to conn  " + err.Error())
	}

	resp, err := bufio.NewReader(conn).ReadBytes('\r')
	if err != nil {
		fmt.Println("Error reading conn data ", err.Error())
	}
	fmt.Println(string(resp))

}
