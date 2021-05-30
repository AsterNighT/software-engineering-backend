package chat

import "fmt";

var Clients map[*Client]bool = make(map[*Client]bool)
var Connections map[*Client][]*Client = make(map[*Client][]*Client)

var Register chan *Client = make(chan *Client)
var Unregister chan *Client = make(chan *Client)

//Pool main loop, deal with client addition and deletion
func PoolLoop() {
	for {
		select {
		case client := <-Register:
			AddClient(client)
			fmt.Println("Size of Connection Pool: ", len(Clients))
		case client := <-Unregister:
			DeleteClient(client)
		}
	}
}

//Add a new client into pool
func AddClient(client *Client) {
	Clients[client] = true
}

//Delete a client from pool
func DeleteClient(client *Client) {
	var ok bool = false
	_, ok = Clients[client]
	if ok {
		delete(Clients, client)
	}
	_, ok = Connections[client]
	if ok {
		delete(Connections, client)
	}
	// _, ok = DoctorConns[client]
	// if ok {
	// 	delete(DoctorConns, client)
	// }
	// _, ok = PatientConns[client]
	// if ok {
	// 	delete(PatientConns, client)
	// }

}
