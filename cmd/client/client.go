package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, err := d.DialContext(ctx, "tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	Write(conn)
}

func Write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: ", err)
			os.Exit(1)
		}
		if _, err := conn.Write([]byte(line)); err != nil {
			log.Fatal(err)
		}
	}

}
