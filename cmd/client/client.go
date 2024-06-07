package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

const (
	HOST     = "192.168.1.126"
	PORT     = "8000"
	PROTOCOL = "tcp"
	CONN     = HOST + ":" + PORT
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	conn, err := net.Dial(PROTOCOL, CONN)
	if err != nil {
		fmt.Println(err)
	}

	go Read(conn)
	go Write(conn)

	wg.Wait()
}

func Read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected")
			return
		}
		fmt.Print(message)
	}
}

func Write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = writer.WriteString(message)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err = writer.Flush(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
