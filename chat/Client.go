package chat

import (
	"github.com/expproletariy/base_server/types"
	"github.com/gorilla/websocket"
	"time"
)

//Client - type which represented chat clint in messages exchange
type Client struct {
	//Client connection
	conn *websocket.Conn

	types.User
}

//Create new user client
func NewClient(conn *websocket.Conn, user types.User) *Client {
	return &Client{
		User: user,
		conn: conn,
	}
}

//Close client websocket connection
func (clt *Client) Close() error {
	clt.conn.WriteControl(websocket.CloseMessage, []byte(""), time.Now().Add(time.Second))
	return clt.conn.Close()
}

//ReadMessage from user
func (clt *Client) ReadMessage() (types.Message, error) {
	msg := types.Message{}
	err := clt.conn.ReadJSON(&msg)
	return msg, err
}
