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

	connverseAscii := `                                              
  ___ ___  _ __  _ ____   _____ _ __ ___  ___ 
 / __/ _ \| '_ \| '_ \ \ / / _ \ '__/ __|/ _ \
| (_| (_) | | | | | | \ V /  __/ |  \__ \  __/
 \___\___/|_| |_|_| |_|\_/ \___|_|  |___/\___|
  `
	fmt.Println(connverseAscii)
	fmt.Println("Listenning on " + CONN + " 🚀")

	lobby := chat.NewLobby()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
		}

		client := chat.NewClient(conn)
		client.Log(fmt.Sprintf(chat.WELCOME, connverseAscii))
		lobby.JoinClient(client)

		go chat.HandleClientInput(client, lobby)
	}
}
