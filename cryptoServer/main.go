package main

import (
	"bufio"
	"cryptoServer/controller"
	"cryptoServer/database"
	"cryptoServer/router"
	"cryptoServer/transactions"
	"cryptoServer/utils"
	"fmt"
	"net"
)

func main() {

	db := database.NewDatabase()
	te := transactions.NewTransactionEngine(db)
	c := controller.NewController(db, *te)
	r := router.NewRouter(*c)
	utils.CreateWallets(db, 10)

	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	fmt.Println("Server started listening ")

	if err != nil {
		fmt.Println("Error listening  " + err.Error())

	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting  " + err.Error())
		}
		defer conn.Close()
		fmt.Println("Accepted conncetion")

		req, err := bufio.NewReader(conn).ReadBytes('\r')
		if err != nil {
			fmt.Println("Error reading conn data ", err.Error())
		}

		fmt.Println("Copied data ")
		messages := make(chan []byte)
		fmt.Println(string(req))

		go func() {
			fmt.Println("Started go routine ")

			resp, err := r.HandleRequest(req, r.RouteRequest)
			if err != nil {
				messages <- []byte("error" + err.Error())
				fmt.Println("Error in server " + err.Error())
			}
			messages <- resp
		}()
		response := <-messages
		conn.Write(append(response, '\r'))

	}
}
