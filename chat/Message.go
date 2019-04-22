package chat

import "time"

//Message type in chat exchange
type Message struct {
	//Time when message was sent
	Time time.Time `json:"time"`

	//Message text
	Text string `json:"text"`

	//User ID
	UserID string `json:"user_id"`

	//User name
	UserName string `json:"user_name"`
}
