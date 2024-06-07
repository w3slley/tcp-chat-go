package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	TRANSPORT_PROTOCOL = "tcp"
	HOST               = "localhost"
	PORT               = "8000"
)

func main() {
	listener, err := net.Listen(TRANSPORT_PROTOCOL, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Client connected")
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	buf := make([]byte, 1024)
	c.Read(buf)
	os.Stdout.Write(buf)
}
