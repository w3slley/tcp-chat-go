package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	TRANSPORT_PROTOCOL = "tcp"
	HOST               = "192.168.1.126"
	PORT               = "8000"
	CONN               = HOST + ":" + PORT
)

func main() {
	connections := make([]net.Conn, 2)

	listener, err := net.Listen(TRANSPORT_PROTOCOL, CONN)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	fmt.Println("Listenning on " + CONN)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		fmt.Println("Client connected!")

		connections = append(connections, conn)

		go Read(conn, connections)
		//fmt.Println(connections)
	}
}

func Read(conn net.Conn, connections []net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected", err.Error())
			return
		}
		for _, connection := range connections {
			if connection != nil {
				connection.Write([]byte(message))
			}
		}
	}
}
