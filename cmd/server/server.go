package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	TRANSPORT_PROTOCOL = "tcp"
	HOST               = "localhost"
	PORT               = "8000"
	CONN               = HOST + ":" + PORT
	MAX_CLIENTS        = 100
)

type Client struct {
	id   int
	conn net.Conn
}

func main() {
	clients := make([]Client, MAX_CLIENTS)
	p := 0

	listener, err := net.Listen(TRANSPORT_PROTOCOL, CONN)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Println(`                                              
  ___ ___  _ __  _ ____   _____ _ __ ___  ___ 
 / __/ _ \| '_ \| '_ \ \ / / _ \ '__/ __|/ _ \
| (_| (_) | | | | | | \ V /  __/ |  \__ \  __/
 \___\___/|_| |_|_| |_|\_/ \___|_|  |___/\___|
  `)
	fmt.Println("Listenning on " + CONN + " 🚀")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		fmt.Println("Client connected!")

		//create client struct here and add it to the lobby
		client := &Client{id: p, conn: conn}
		clients[p] = *client
		p += 1

		go handleConnection(conn, client, clients)
	}
}

func handleConnection(conn net.Conn, curr *Client, clients []Client) {
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected", err.Error())
			return
		}
		for _, client := range clients {
			if curr.id != client.id && client.conn != nil {
				client.conn.Write([]byte(message))
			}
		}
	}
}
