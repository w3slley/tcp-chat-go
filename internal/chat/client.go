package chat

import (
	"fmt"
	"net"

	"github.com/google/uuid"
)

type Client struct {
	id       string
	username string
	color    []byte //hex
	conn     net.Conn
	room     *Room
}

func (c *Client) Write(message string, sender *Client) {
	if c.conn != nil {
		c.conn.Write([]byte(fmt.Sprintf("%s: %s", sender.username, message)))
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
	if c.IsInLobby() {
		lobby.RemoveClient(c)
	} else if c.room != nil && len(c.room.clients) == 1 {
		room.lobby.RemoveRoom(c.room.id)
	}
	room.BroadcastLog(fmt.Sprintf(USER_JOINED_ROOM, c.username))
	room.JoinClient(c)
}

func (c *Client) LeaveRoom() {
	lobby := c.room.lobby
	if c.IsInLobby() {
		c.Log(NOT_IN_ROOM)
		lobby.Help(c)
	} else {
		c.room.RemoveClient(c)
		lobby.JoinClient(c)
	}
}

func (c *Client) ChangeUsername(username string) {
	c.username = username
	c.Log(fmt.Sprintf(NEW_USERNAME, username))
}

func NewClient(conn net.Conn) *Client {
	client := &Client{id: uuid.New().String(), username: DEFAULT_USERNAME, conn: conn}
	fmt.Printf(CLIENT_CONNECTED, client.id)
	return client
}
