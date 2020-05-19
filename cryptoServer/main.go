package main

import (
	"bufio"
	"cryptoServer/controller"
	"cryptoServer/database"
	"cryptoServer/router"
	"cryptoServer/transactions"
	"cryptoServer/utils"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {

	nwallets := flag.Int("nwallets", 10, "set the number of wallets to be seeded ")
	port := flag.Int("p", 8080, "set the port for the server to listen to ")
	flag.Parse()
	db := database.NewDatabase()
	te := transactions.NewTransactionEngine(db)
	c := controller.NewController(db, *te)
	r := router.NewRouter(c)
	utils.CreateWallets(db, *nwallets)
	p := fmt.Sprintf(":%v", (*port))
	ln, err := net.Listen("tcp", p)
	fmt.Println("Server started listening ")

	if err != nil {
		fmt.Println("Error listening  " + err.Error())
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting  " + err.Error())
			os.Exit(1)
		}
		defer conn.Close()
		fmt.Println("Accepted conncetion")

		req, err := bufio.NewReader(conn).ReadBytes('\r')
		if err != nil {
			fmt.Println("Error reading conn data ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Copied data ")
		fmt.Println(string(req))

		go func(conn net.Conn) {
			fmt.Println("Started go routine ")

			resp, err := r.HandleRequest(req, r.RouteRequest)
			if err != nil {
				conn.Write([]byte("error" + err.Error()))
				fmt.Println("Error in server " + err.Error())
			} else {
				conn.Write(append(resp, '\r'))
			}

		}(conn)
	}
}
