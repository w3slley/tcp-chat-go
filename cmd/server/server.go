package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"slices"
	"strings"
	"time"

	"connverse/helpers"
	"github.com/google/uuid"
)

const (
	TRANSPORT_PROTOCOL = "tcp"
	HOST               = "localhost"
	PORT               = "8000"
	CONN               = HOST + ":" + PORT
)

const (
	DEFAULT_USERNAME        = "anonymous"
	JOIN_ROOM_COMMAND       = "/join"
	LEAVE_ROOM_COMMAND      = "/leave"
	LIST_ROOMS_COMMAND      = "/list"
	CHANGE_USERNAME_COMMAND = "/username"
	SEND_MESSAGE_COMMAND    = "/send"
	QUIT_COMMAND            = "/quit"
	HELP_COMMAND            = "/help"
)

type Client struct {
	id       string
	username string
	color    string //store hex
	conn     net.Conn
	room     *Room
}

type Room struct {
	id       string
	name     string
	clients  []*Client
	messages []*Message
	lobby    *Lobby
}

type Lobby struct {
	clients []*Client
	rooms   []*Room
}

type Message struct {
	sender  *Client
	message string
	time    time.Time
}

func (c *Client) Write(message string) {
	if c.conn != nil {
		c.conn.Write([]byte(fmt.Sprintf("%s: %s", c.username, message)))
	}
}

func (c *Client) Log(message string) {
	if c.conn != nil {
		c.conn.Write([]byte(message))
	}
}

func (c *Client) IsInLobby() bool {
	return c.room == nil
}

func (c *Client) Quit() {
	c.conn.Close()
}

func (c *Client) JoinRoom(lobby *Lobby, roomName string) {
	room := lobby.GetRoomByName(roomName)
	if room == nil {
		room = lobby.NewRoom(roomName)
	}
	lobby.RemoveClient(c)
	room.Join(c)
}

func (c *Client) LeaveRoom() {
	lobby := c.room.lobby
	if c.IsInLobby() {
		c.Log(messages.NOT_IN_ROOM)
		lobby.Help(c)
	} else {
		c.room.Remove(c)
		lobby.JoinClient(c)
	}
}

func (r *Room) Broadcast(sender *Client, message string) {
	for _, receiver := range r.clients {
		receiver.Write(message)
	}
}

func (r *Room) Join(client *Client) {
	r.clients = append(r.clients, client)
	client.room = r

	client.Log(fmt.Sprintf(messages.JOINED_ROOM, r.name))
}

func (r *Room) Remove(client *Client) {
	indexToDelete := -1
	for i, clientInRoom := range r.clients {
		if clientInRoom.id == client.id {
			indexToDelete = i
		} else {
			clientInRoom.Log(fmt.Sprintf(messages.USER_LEFT_ROOM, client.username))
		}
	}
	if indexToDelete != -1 {
		r.clients = slices.Delete(r.clients, indexToDelete, indexToDelete+1)
		client.Log(fmt.Sprintf(messages.LEFT_ROOM, r.name))
	}
	if len(r.clients) == 0 {
		r.lobby.RemoveRoom(r.id)
	}
	client.room = nil
}

func (l *Lobby) Broadcast(sender *Client, message string) {
	for _, receiver := range l.clients {
		receiver.Write(message)
	}
}

func (l *Lobby) GetRoomByName(name string) *Room {
	for _, room := range l.rooms {
		if room.name == name {
			return room
		}
	}
	return nil
}

func (l *Lobby) RemoveClient(client *Client) {
	for i, clientInLobby := range l.clients {
		if clientInLobby.id == client.id {
			l.clients = slices.Delete(l.clients, i, i+1)
		}
	}
}

func (l *Lobby) JoinClient(client *Client) {
	l.clients = append(l.clients, client)
}

func (l *Lobby) NewRoom(name string) *Room {
	var clients []*Client
	room := &Room{id: uuid.New().String(), name: name, clients: clients, lobby: l}
	fmt.Printf(messages.ROOM_CREATED, room.id)
	l.rooms = append(l.rooms, room)
	return room
}

func (l *Lobby) RemoveRoom(id string) {
	for i, room := range l.rooms {
		if room.id == id {
			l.rooms = slices.Delete(l.rooms, i, i+1)
		}
	}
}

func (l *Lobby) ListRooms(client *Client) {
	if l.rooms == nil {
		client.Log(messages.NO_ROOMS)
	}
	for _, room := range l.rooms {
		if client.room.id == room.id {
			client.Log(messages.CURRENT_ROOM_ICON)
		}
		client.Log(fmt.Sprintf("%s\n", room.name))
	}
}

func (l *Lobby) Help(client *Client) {
	client.Log("\n")
	client.Log("Commands:\n")
	client.Log("/help - list of all commands\n")
	client.Log("/list - lists all the rooms\n")
	client.Log("/join <room> - joins room named <room> if it exists or creates it if not\n")
	client.Log("/username <username> - changes username to <username>\n")
	client.Log("/quit - exits the program\n")
}

func NewClient(conn net.Conn) *Client {
	client := &Client{id: uuid.New().String(), username: DEFAULT_USERNAME, conn: conn}
	fmt.Printf(messages.CLIENT_CONNECTED, client.id)
	return client
}

func NewLobby() *Lobby {
	return &Lobby{
		clients: []*Client{},
	}
}

func GetCommandArgument(message string, cmd string) string {
	return strings.TrimSuffix(strings.TrimPrefix(message, JOIN_ROOM_COMMAND+" "), "\n")
}

func HandleClientInput(client *Client, lobby *Lobby) {
	for {
		reader := bufio.NewReader(client.conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf(messages.USER_DISCONNECTED, client.id)
			break
		}
		switch {
		case strings.HasPrefix(message, JOIN_ROOM_COMMAND):
			roomName := GetCommandArgument(message, JOIN_ROOM_COMMAND)
			client.JoinRoom(lobby, roomName)

		case strings.HasPrefix(message, LEAVE_ROOM_COMMAND):
			client.LeaveRoom()

		case strings.HasPrefix(message, LIST_ROOMS_COMMAND):
			lobby.ListRooms(client)

		case strings.HasPrefix(message, SEND_MESSAGE_COMMAND):
		case strings.HasPrefix(message, CHANGE_USERNAME_COMMAND):
		case strings.HasPrefix(message, HELP_COMMAND):
			lobby.Help(client)
		case strings.HasPrefix(message, QUIT_COMMAND):
			client.Quit()

		default:
			if client.IsInLobby() {
				lobby.Broadcast(client, message)
			} else {
				client.room.Broadcast(client, message)
			}
		}
	}
}

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

	lobby := NewLobby()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
		}

		client := NewClient(conn)
		client.Log(fmt.Sprintf("%s \n Welcome to connverse, your TCP chat application accessed via SSH! \n", connverseAscii))
		lobby.JoinClient(client)

		go HandleClientInput(client, lobby)
	}
}
