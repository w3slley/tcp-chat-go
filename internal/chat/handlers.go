package chat

import (
	"bufio"
	"fmt"

	"connverse/pkg/utils"
)

func HandleClientInput(client *Client, lobby *Lobby) {
	for {
		reader := bufio.NewReader(client.conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf(USER_DISCONNECTED, client.id)
			break
		}

		command := utils.GetCommandFromMessage(message)

		switch command {
		case JOIN_ROOM_COMMAND:
			roomName := utils.GetCommandArgument(message, JOIN_ROOM_COMMAND)
			client.JoinRoom(lobby, roomName)

		case LEAVE_ROOM_COMMAND:
			client.LeaveRoom()

		case LIST_ROOMS_COMMAND:
			lobby.ListRooms(client)

		case SEND_MESSAGE_COMMAND:
			//TODO: Implement

		case CHANGE_USERNAME_COMMAND:
			newUsername := utils.GetCommandArgument(message, CHANGE_USERNAME_COMMAND)
			client.ChangeUsername(newUsername)

		case HELP_COMMAND:
			lobby.Help(client)

		case QUIT_COMMAND:
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
