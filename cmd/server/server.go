package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

const (
	TRANSPORT_PROTOCOL = "tcp"
	HOST               = "localhost"
	PORT               = "8000"
	CONN               = HOST + ":" + PORT
	MAX_CLIENTS        = 100
)

type Client struct {
	id   string
	name string
	conn net.Conn
	room *Room
}

type Room struct {
	id      string
	clients []Client
}

func (c *Client) Write(message string) {
	c.conn.Write([]byte(message))
}

func NewRoom() *Room {
	clients := make([]Client, MAX_CLIENTS)
	room := &Room{id: uuid.New().String(), clients: clients}
	fmt.Printf("Room with id %s was created\n", room.id)
	return room
}

func NewClient(conn net.Conn, room *Room) *Client {
	client := &Client{id: uuid.New().String(), conn: conn, room: room}
	fmt.Printf("Client %s connected \n", client.id)
	return client
}

func (r *Room) Broadcast(curr *Client, message string) {
	for _, client := range r.clients {
		if curr.id != client.id && client.conn != nil {
			client.Write(message)
		}
	}
}

func processNewClient(listener net.Listener, room *Room) (bool, error) {
	conn, err := listener.Accept()
	if err != nil {
		log.Println("Error: ", err)
		return false, err
	}

	client := NewClient(conn, room)
	room.clients = append(room.clients, *client)

	go handleClientInput(client)

	return true, nil
}

func handleClientInput(client *Client) {
	for {
		reader := bufio.NewReader(client.conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected", err.Error())
			return
		}
		client.room.Broadcast(client, message)
	}
}

func main() {

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

	room := NewRoom()
	for {
		processNewClient(listener, room)
	}
}
