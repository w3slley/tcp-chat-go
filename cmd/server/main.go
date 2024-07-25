package main

import (
	"fmt"
	"log"
	"net"

	"connverse/internal/chat"
)

const (
	TRANSPORT_PROTOCOL = "tcp"
	HOST               = "localhost"
	PORT               = "8000"
	CONN               = HOST + ":" + PORT
)

func main() {
	listener, err := net.Listen(TRANSPORT_PROTOCOL, CONN)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Println("Listenning on " + CONN + " 🚀")

	lobby := chat.NewLobby()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
		}

		client := chat.NewClient(conn)
		lobby.JoinClient(client)

		go chat.HandleClientInput(client, lobby)
	}
}
