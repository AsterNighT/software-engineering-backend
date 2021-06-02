package chat

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Type       int      `json:"type"`
	Role       string   `json:"role,omitempty"`
	Name       string   `json:"name,omitempty"`
	SenderID   string   `json:"senderid,omitempty"`
	ReceiverID string   `json:"receiverid,omitempty"`
	PatientID  string   `json:"petientid,omitempty"`
	DoctorID   string   `json:"doctorid,omitempty"`
	Content    string   `json:"content,omitempty"`
	Time       string   `json:"time,omitempty"`
	URL        string   `json:"url,omitempty"`
	Questions  []string `json:"questions,omitempty"`
}

//Process one message
func (sender *Client) ProcessMessage(msgBytes []byte) {
	message := &Message{}
	json.Unmarshal(msgBytes, message)
	switch message.Type {
	//client to server
	case 0:
		sender.Hello(message)
	case 1:
		sender.MsgFromClient(message)
	case 2:
		sender.CloseChat(message)
	case 3:
		sender.RequireMedicalRecord(message)
	case 4:
		sender.RequirePrescription(message)
	case 5:
		sender.RequireQuestions(message)
	//server to client
	case 6:
		sender.NewPatient(message)
	case 7:
		sender.MsgFromServer(message)
	case 8:
		sender.SendMedicalRecord(message)
	case 9:
		sender.SendPrescription(message)
	case 10:
		sender.SendQuestions(message)
	default:
		sender.WrongMsgType(message)
	}
}

//Process hello message
func (sender *Client) Hello(message *Message) {
	//strconv.Itoa(len(Clients))//message.SenderID
	sender.Role = message.Role

	fmt.Println("in hello " + sender.ID + " *** " + sender.Role)
}

//Process msgfromclient message
func (sender *Client) MsgFromClient(message *Message) {
	receiver := sender.FindReceiver(message)
	if receiver == nil {
		sender.ReceiverNotConnected(message)
	}
	fmt.Println("in msgfromclient " + sender.ID + " *** " + message.Content)

	msgBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("ChatServer:$Error:" + err.Error())
		return
	}

	receiver.MsgBuffer <- msgBytes //add the message to receiver buffer
}

//Process closechat message
func (sender *Client) CloseChat(message *Message) {
}

//Process requiremedicalrecord message
func (sender *Client) RequireMedicalRecord(message *Message) {
	msg := Message{
		Type:      8,
		PatientID: message.PatientID,
		URL:       "url", // get from database?
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("ChatServer:$Error:" + err.Error())
		return
	}

	sender.MsgBuffer <- msgBytes //add the message to sender buffer
}

//Process requireprescription message
func (sender *Client) RequirePrescription(message *Message) {
	msg := Message{
		Type:      9,
		PatientID: message.PatientID,
		URL:       "url", // get from database?
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("ChatServer:$Error:" + err.Error())
		return
	}

	sender.MsgBuffer <- msgBytes //add the message to sender buffer
}

//Process requirequestions message
func (sender *Client) RequireQuestions(message *Message) {
	msg := Message{
		Type:      10,
		Questions: []string{"aaa", "bbb", "ccc"},
		// get from database based on message.DoctorID?
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("ChatServer:$Error:" + err.Error())
		return
	}

	sender.MsgBuffer <- msgBytes //add the message to sender buffer
}

//Process newpatient message
func (sender *Client) NewPatient(message *Message) {
}

//Process msgfromserver message
func (sender *Client) MsgFromServer(message *Message) {
}

//Process sendmedicalrecord message
func (sender *Client) SendMedicalRecord(message *Message) {
}

//Process sendprescription message
func (sender *Client) SendPrescription(message *Message) {
}

//Process sendquestions message
func (sender *Client) SendQuestions(message *Message) {
}

//Find the receiver of specific message
func (sender *Client) FindReceiver(message *Message) *Client {
	var receiver *Client = nil
	_, ok := Connections[sender]

	if !ok { //map result doesn't exist
		for client := range Clients { //search the connected clients for receiver
			if client.ID == message.ReceiverID {
				receiver = client
				break
			}
		}
		if receiver == nil { //receiver not connected to server yet
			sender.ReceiverNotConnected(message)
			return nil
		}
		slice := make([]*Client, 5)
		slice[0] = receiver
		Connections[sender] = slice
	} else { //map result exists
		slice := Connections[sender]
		for _, client := range slice { //search the map result for receiver
			if client.ID == message.ReceiverID {
				receiver = client
				break
			}
		}
		if receiver == nil { //receiver not in map result
			for client := range Clients { //search the connected clients for receiver
				if client.ID == message.ReceiverID {
					receiver = client
					break
				}
			}
			if receiver == nil { //receiver not connected to server yet
				sender.ReceiverNotConnected(message)
				return nil
			}
			Connections[sender] = append(slice, receiver) //add receiver to map result
		}

	}
	return receiver
}

//Deal with unknown message type
func (sender *Client) WrongMsgType(message *Message) {

}

//Deal with the case when receiver of the message has't connected to the server
func (sender *Client) ReceiverNotConnected(message *Message) {

}
