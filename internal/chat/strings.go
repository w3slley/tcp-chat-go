package chat

const (
	DEFAULT_USERNAME        = "anonymous"
	JOIN_ROOM_COMMAND       = "/join"
	LEAVE_ROOM_COMMAND      = "/leave"
	LIST_ROOMS_COMMAND      = "/list"
	CHANGE_USERNAME_COMMAND = "/username"
	SEND_MESSAGE_COMMAND    = "/send"
	QUIT_COMMAND            = "/quit"
	HELP_COMMAND            = "/help"

	CLIENT_CONNECTED  = "Client %s connected \n"
	JOINED_ROOM       = "Welcome to the room: %s. \n"
	LEFT_ROOM         = "You have left the room %s. Now you are in the lobby!\n"
	USER_LEFT_ROOM    = "%s left the room.\n"
	ROOM_CREATED      = "Room with id %s was created!\n"
	USER_DISCONNECTED = "User with id %s disconnected!\n"
	NOT_IN_ROOM       = "You are not in a room!\n"
	NO_ROOMS          = "There are no rooms.\n"
	NEW_USERNAME      = "Your username was changed to %s.\n"
	CURRENT_ROOM_ICON = "* "
)
