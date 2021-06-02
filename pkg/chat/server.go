package chat

import (
	"fmt"
	"net/http"
	"strconv"
)

//Serve an incoming http connection
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ChatServer:$Websocket endpoint hit")
	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Println("ChatServer:$Error:" + err.Error())
		return 
	}
	newClient := &Client{
		ID:        strconv.Itoa(len(Clients)),
		Conn:      conn,
		MsgBuffer: make(chan []byte),
	}


	Register <- newClient
	go newClient.Read()
	go newClient.Send()
}

//Start the chat server
func StartServer() {
	fmt.Println("ChatServer:$System initializing")
	go PoolLoop()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	fmt.Println("ChatServer:$System on")
	http.ListenAndServe(":8080", nil)
}
