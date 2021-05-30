package chat

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID        string
	Role      string //0:Doctor, 1:Patient
	Conn      *websocket.Conn
	MsgBuffer chan []byte //Message to send to this client
}

type Doctor struct {
	Client
}

type Patient struct {
	Client
}

//Read subroutine for client
func (client *Client) Read() {
	defer func() { //delete the client from pool
		Unregister <- client
		close(client.MsgBuffer)
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println("ChatServer:$Error:" + err.Error())
			break
		}

		client.ProcessMessage(message)
	}

}

//Send subroutine for client
func (client *Client) Send() {
	for {
		message, ok := <-client.MsgBuffer
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}
		client.Conn.WriteMessage(websocket.TextMessage, message)
	}
}
