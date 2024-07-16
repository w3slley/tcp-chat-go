package chat

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type Lobby struct {
	clients []*Client
	rooms   []*Room
}

func (l *Lobby) Broadcast(sender *Client, message string) {
	for _, receiver := range l.clients {
		receiver.Write(message, sender)
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
			break
		}
	}
}

func (l *Lobby) JoinClient(client *Client) {
	l.clients = append(l.clients, client)
}

func (l *Lobby) NewRoom(name string) *Room {
	var clients []*Client
	room := &Room{id: uuid.New().String(), name: name, clients: clients, lobby: l}
	fmt.Printf(ROOM_CREATED, room.id)
	l.rooms = append(l.rooms, room)
	return room
}

func (l *Lobby) RemoveRoom(id string) {
	for i, room := range l.rooms {
		if room.id == id {
			l.rooms = slices.Delete(l.rooms, i, i+1)
			break
		}
	}

}

func (l *Lobby) ListRooms(client *Client) {
	if l.rooms == nil {
		client.Log(NO_ROOMS)
		return
	}
	for _, room := range l.rooms {
		if client.room.id == room.id {
			client.Log(CURRENT_ROOM_ICON)
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

func NewLobby() *Lobby {
	return &Lobby{
		clients: []*Client{},
	}
}
