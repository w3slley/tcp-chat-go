package chat

import (
	"fmt"
	"slices"
)

type Room struct {
	id       string
	name     string
	clients  []*Client
	messages []*Message
	lobby    *Lobby
}

func (r *Room) Broadcast(sender *Client, message string) {
	for _, receiver := range r.clients {
		receiver.Write(message, sender)
	}
}

func (r *Room) BroadcastLog(message string) {
	for _, receiver := range r.clients {
		receiver.Log(message)
	}
}

func (r *Room) JoinClient(client *Client) {
	r.clients = append(r.clients, client)
	client.room = r

	client.Log(fmt.Sprintf(JOINED_ROOM, r.name))
}

func (r *Room) RemoveClient(client *Client) {
	indexToDelete := -1
	for i, clientInRoom := range r.clients {
		if clientInRoom.id == client.id {
			indexToDelete = i
		} else {
			clientInRoom.Log(fmt.Sprintf(USER_LEFT_ROOM, client.username))
		}
	}
	if indexToDelete != -1 {
		r.clients = slices.Delete(r.clients, indexToDelete, indexToDelete+1)
		client.Log(fmt.Sprintf(LEFT_ROOM, r.name))
	}
	if len(r.clients) == 0 {
		r.lobby.RemoveRoom(r.id)
	}
	client.room = nil
}
