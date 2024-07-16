package chat

import "time"

type Message struct {
	sender  *Client
	message string
	time    time.Time
}
