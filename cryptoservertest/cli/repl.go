package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Order struct {
	WalletID    string  `json:"walletID"`
	CurrencyID  string  `json:"currencyID"`
	Price       float64 `json:"limitPrice"`
	OrderType   int     `json:"OrderType"`
	SumToInvest float64 `json:"sumToInvest"`
	UserID      string  `json:"userID"`
	Deleted     bool
}

type CancelOrder struct {
	OrderID string `json:"orderID"`
}

type Request struct {
	RequestType int    `json:"requestType"`
	UserID      string `json:"userID"`
	Data        []byte `json:"data"`
}

func Run() {

	for {
		fmt.Println("1--Register")
		fmt.Println("2--ListWallets")
		fmt.Println("3--List top 10 lowest sells and top 10 highest buys ")
		fmt.Println("4--Place order ")
		fmt.Println("5--List Your orders ")
		fmt.Println("6--Cancel orders ")
		fmt.Println("e-- to exit ")
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Use one of the numbers preceding the commands to chose functionality followed by enter")
		choice := ""
		for scanner.Scan() {
			choice = scanner.Text()
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		i, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Println("Non intiger input detected ")
			os.Exit(1)
		}
		fmt.Println(i)

		switch i {
		case 1:
			register()
			break
		case 2:
			listwallets()
			break
		case 3:
			listOrderRanking()
			break
		case 4:
			placeOrder()
			break
		case 5:
			listUserOrders()
			break
		case 6:
			cancelOrder()
			break

		}
	}
}

func cancelOrder() {

	cancel := CancelOrder{}
	cancel.OrderID = readConsole("Please enter the Order ID you got from the server")
	data, err := json.Marshal(cancel)
	if err != nil {
		fmt.Println("Marshalling of request failed" + err.Error())
	}
	request := Request{
		RequestType: 6,
		Data:        data,
	}
	request.UserID = readConsole("Please enter the User ID you got from the server  ")
	req, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Marshalling of request failed" + err.Error())
	}
	conn := write(req)
	read(conn)
	defer conn.Close()
}

func register() {
	request := Request{
		RequestType: 1,
		Data:        nil,
	}
	req, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Marshalling of request failed" + err.Error())
	}
	conn := write(req)
	read(conn)
	defer conn.Close()

}
func listUserOrders() {
	request := Request{
		RequestType: 5,
		Data:        nil,
	}
	request.UserID = readConsole("Please enter the User ID you got from the server  ")

	req, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Marshalling of request failed" + err.Error())
	}

	conn := write(req)
	read(conn)
	defer conn.Close()
}

func listwallets() {

	request := Request{
		RequestType: 2,
		Data:        nil,
	}
	request.UserID = readConsole("Please enter the ID you got from the server")
	req, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Marshalling of request failed" + err.Error())
	}
	conn := write(req)
	read(conn)
	defer conn.Close()
}

func listOrderRanking() {
	request := Request{
		RequestType: 3,
		Data:        nil,
	}
	request.UserID = readConsole("Please enter the ID you got from the server ")
	req, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Marshalling of request failed" + err.Error())
	}
	conn := write(req)
	read(conn)
	defer conn.Close()
}
func write(req []byte) net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Failed connecting" + err.Error())
		os.Exit(1)
	}
	//defer conn.Close()
	_, err = conn.Write(append(req, '\r'))
	if err != nil {
		fmt.Println("Error writing to conn  " + err.Error())
	}
	return conn
}
func read(conn net.Conn) {
	resp, err := bufio.NewReader(conn).ReadBytes('\r')
	if err != nil {
		fmt.Println("Error reading conn data ", err.Error())
	}
	fmt.Println(string(resp))

}
func placeOrder() {
	order := Order{}
	order.WalletID = readConsole("Please enter one of the wallet IDs you got from the list ")
	order.CurrencyID = readConsole("Please enter one of the currency IDs you got from the list ")
	order.UserID = readConsole("Please enter the User ID you got from the server  ")
	i, err := strconv.ParseFloat(readConsole("Please enter the sum to invest "), 64)
	if err != nil {
		fmt.Println("Non intiger input detected ")
		os.Exit(1)
	}
	order.SumToInvest = i
	orderType, err := strconv.Atoi(readConsole("Please enter 1 for SellOrder or 2 for BuyOrder "))
	if err != nil {
		fmt.Println("Non intiger input detected ")
		os.Exit(1)
	}
	for {
		if orderType == 1 {
			break
		} else if orderType == 2 {
			break
		}
		orderType, err = strconv.Atoi(readConsole("Please enter 1 for SellOrder or 2 for BuyOrder not another number  "))
		fmt.Println(orderType)
		if err != nil {
			fmt.Println("Non intiger input detected ")
			os.Exit(1)
		}
	}
	order.OrderType = orderType

	limitPrice, err := strconv.ParseFloat(readConsole("Please enter your limit price "), 64)
	if err != nil {
		fmt.Println("Non intiger input detected ")
		os.Exit(1)
	}
	order.Price = limitPrice

	data, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Marshalling of order failed" + err.Error())
	}

	request := Request{
		RequestType: 4,
		Data:        data,
		UserID:      order.UserID,
	}
	req, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Marshalling of order failed" + err.Error())
	}
	conn := write(req)
	read(conn)
	defer conn.Close()

}
func readConsole(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(prompt)
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func listYourOrder() {
	request := Request{
		RequestType: 5,
		Data:        nil,
	}
	request.UserID = readConsole("Please enter the ID you got from the server")
	req, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Marshalling of request failed" + err.Error())
	}
	conn := write(req)
	read(conn)
	defer conn.Close()
}
