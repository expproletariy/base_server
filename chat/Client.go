package chat

import "github.com/gorilla/websocket"

//Client - type which represented chat clint in messages exchange
type Client struct {
	//Client connection
	conn *websocket.Conn

	//Client ID
	ID string

	//User name
	Name string
}

//Create new user client
func NewClient(conn *websocket.Conn, id, name string) *Client {
	return &Client{
		Name: name,
		ID:   id,
		conn: conn,
	}
}

//Close client websocket connection
func (clt *Client) Close() error {
	return clt.conn.Close()
}

//ReadMessage from user
func (clt *Client) ReadMessage() (Message, error) {
	msg := Message{}
	err := clt.conn.ReadJSON(&msg)
	return msg, err
}
