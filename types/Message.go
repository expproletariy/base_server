package types

//Message type
type Message struct {
	RoomID   string `json:"room_id" db:"rooms_id"`
	UserName string `json:"user_name"`
	SimpleMessage
}
