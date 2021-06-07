package chat

import (
	"encoding/json"

	"github.com/AsterNighT/software-engineering-backend/api"
	"github.com/AsterNighT/software-engineering-backend/pkg/cases"
	"github.com/AsterNighT/software-engineering-backend/pkg/utils"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type RoleType = int

const (
	Doctor  RoleType = 0
	Patient RoleType = 1
)

type ClientMsgType = int

const (
	MsgFromClient        ClientMsgType = 1
	CloseChat            ClientMsgType = 2
	RequireMedicalRecord ClientMsgType = 3
	RequirePrescription  ClientMsgType = 4
	RequireQuestions     ClientMsgType = 5
)

type ServerMsgType int

const (
	MsgFromServer     ServerMsgType = 6
	NewPatient        ServerMsgType = 7
	SendMedicalRecord ServerMsgType = 8
	SendPrescription  ServerMsgType = 9
	SendQuestions     ServerMsgType = 10
	ChatTerminated    ServerMsgType = 11
)

type Client struct {
	ID        string
	Role      RoleType
	Conn      *websocket.Conn
	MsgBuffer chan []byte //Message to send to this client
}

var (
	upgrader    = websocket.Upgrader{}
	Clients     = make(map[*Client]bool)
	Connections = make(map[*Client][]*Client)
)

//Add a new client into pool
func AddClient(client *Client) {
	Clients[client] = true
}

//Delete a client from pool
func DeleteClient(client *Client) {
	delete(Clients, client)
	delete(Connections, client)
}

type ChatHandler struct {
}

type CategoryHandler struct {
}

// @Summary create a new connection for patient
// @Description
// @Tags Chat
// @Produce json
// @Param newClient body Client true "id, role, conn and msgbuffer"
// @Success 200 {object} api.ReturnedData{}
// @Router /patient/{patientID}/chat [POST]
func (h *ChatHandler) NewPatientConn(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(400, api.Return("Upgrade Fail", nil))
	}
	defer conn.Close()
	newClient := &Client{
		ID:        c.Param("patientID"),
		Role:      Patient,
		Conn:      conn,
		MsgBuffer: make(chan []byte),
	}
	//Add client to database
	AddClient(newClient)
	go newClient.Read(c)
	go newClient.Send(c)

	c.Logger().Debug("ChatServer$: NewPatientConn")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary create a new connection for doctor
// @Description
// @Tags Chat
// @Produce json
// @Param newClient body Client true "id, role, conn and msgbuffer"
// @Success 200 {object} api.ReturnedData{}
// @Router /doctor/{doctorID}/chat [POST]
func (h *ChatHandler) NewDoctorConn(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(400, api.Return("Upgrade Fail", nil))
	}
	newClient := &Client{
		ID:        c.Param("doctorID"),
		Role:      Doctor,
		Conn:      conn,
		MsgBuffer: make(chan []byte),
	}
	//Add client to database
	AddClient(newClient)
	go newClient.Read(c)
	go newClient.Send(c)

	c.Logger().Debug("ChatServer$: NewDoctorConn")
	return c.JSON(200, api.Return("ok", nil))
}

// @Summary Get questions by department id
// @Description
// @Tags Chat
// @Produce json
// @Param Department path uint true "department ID"
// @Success 200 {object} api.ReturnedData{data=[]string}
// @Router /department/{departmentID}  [GET]
func (h *CategoryHandler) GetQuestionsByDepartmentID(c echo.Context) error {
	db := utils.GetDB()
	db.Where("DepartmentID = ?", c.Param("DepartmentID"))

	var cate Category
	db.Find(&cate)

	c.Logger().Debug("ChatServer$: GetQuestionsByDepartmentID")
	return c.JSON(200, api.Return("ok", cate.Questions))
}

//Read subroutine for client
func (client *Client) Read(c echo.Context) {
	defer func() { //delete the client from pool
		//delete client from database
		DeleteClient(client)
		close(client.MsgBuffer)
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			c.Logger().Debug("ChatServer$: Read: " + err.Error())
			break
		}
		client.ProcessMessage(message, c)
	}
}

//Send subroutine for client
func (client *Client) Send(c echo.Context) {
	for {
		message, ok := <-client.MsgBuffer
		if !ok {
			err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			if err != nil {
				c.Logger().Debug("ChatServer$: Send: " + err.Error())
				return
			}
			break
		}

		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			c.Logger().Debug("ChatServer$: Send: " + err.Error())
			return
		}
	}
}

type Message struct {
	Type       int      `json:"type"`
	Role       RoleType `json:"role,omitempty"`
	Name       string   `json:"name,omitempty"`
	SenderID   string   `json:"senderid,omitempty"`
	ReceiverID string   `json:"receiverid,omitempty"`
	PatientID  string   `json:"petientid,omitempty"`
	DoctorID   string   `json:"doctorid,omitempty"`
	Content    string   `json:"content,omitempty"`
	Time       string   `json:"time,omitempty"`
	CaseID     string   `json:"caseid,omitempty"`
	URL        string   `json:"url,omitempty"`
	Questions  []string `json:"questions,omitempty"`
}

//Process one message
func (client *Client) ProcessMessage(msgBytes []byte, c echo.Context) {
	message := &Message{}
	err := json.Unmarshal(msgBytes, message)
	if err != nil {
		c.Logger().Debug("ChatServer$: ProcessMessage: " + err.Error())
		return
	}
	switch message.Type {
	case MsgFromClient:
		client.MsgFromClient(message, c)
	case CloseChat:
		client.CloseChat(message, c)
	case RequireMedicalRecord:
		client.RequireMedicalRecord(message, c)
	case RequirePrescription:
		client.RequirePrescription(message, c)
	case RequireQuestions:
		client.RequireQuestions(message, c)
	default:
		client.WrongMsgType(message, c)
	}
}

//Process msgfromclient message
func (client *Client) MsgFromClient(message *Message, c echo.Context) {
	receiver := client.FindReceiver(message, c)
	if receiver == nil {
		client.ReceiverNotConnected(message, c)
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		c.Logger().Debug("ChatServer$: MsgFromClient: " + err.Error())
		return
	}
	receiver.MsgBuffer <- msgBytes //add the message to receiver buffer
}

//Process closechat message
func (client *Client) CloseChat(message *Message, c echo.Context) {
	receiver := client.FindReceiver(message, c)
	if receiver == nil {
		client.ReceiverNotConnected(message, c)
	}

	msg := Message{
		Type: int(ChatTerminated),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		c.Logger().Debug("ChatServer$: CloseChat: " + err.Error())
		return
	}
	client.MsgBuffer <- msgBytes //add the message to sender buffer
	receiver.MsgBuffer <- msgBytes

	//terminate the connection of receiver
}

//TODO no medicalrecord
//Process requiremedicalrecord message
func (client *Client) RequireMedicalRecord(message *Message, c echo.Context) {
	msg := Message{
		Type:      int(SendMedicalRecord),
		PatientID: message.PatientID,
		URL:       "url", // get from database
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		c.Logger().Debug("ChatServer$: RequireMedicalRecord: " + err.Error())
		return
	}
	client.MsgBuffer <- msgBytes //add the message to sender buffer
}

//Process requireprescription message
func (client *Client) RequirePrescription(message *Message, c echo.Context) {
	db := utils.GetDB()
	var pres []cases.Prescription
	db.Where("CaseID = ?", message.CaseID).Find(&pres)

	msg := Message{
		Type:      int(SendPrescription),
		PatientID: message.PatientID,
		URL:       "pres url", //TODO how to convert
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		c.Logger().Debug("ChatServer$: RequirePrescription: " + err.Error())
		return
	}
	client.MsgBuffer <- msgBytes //add the message to sender buffer
}

//Process requirequestions message
func (client *Client) RequireQuestions(message *Message, c echo.Context) {
	db := utils.GetDB()

	var cate Category
	db.Where("DepartmentID = ?", db.Select("DepartmentID").Where("DoctorID = ?")).Find(&cate)

	msg := Message{
		Type:      int(SendQuestions),
		Questions: cate.Questions,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		c.Logger().Debug("ChatServer$: RequireQuestions: " + err.Error())
		return
	}
	client.MsgBuffer <- msgBytes //add the message to sender buffer
}

// //Process newpatient message
// func (client *Client) NewPatient(message *Message) {
// }

// //Process msgfromserver message
// func (client *Client) MsgFromServer(message *Message) {
// }

// //Process sendmedicalrecord message
// func (client *Client) SendMedicalRecord(message *Message) {
// }

// //Process sendprescription message
// func (client *Client) SendPrescription(message *Message) {
// }

// //Process sendquestions message
// func (client *Client) SendQuestions(message *Message) {
// }

//Find the receiver of specific message
func (client *Client) FindReceiver(message *Message, c echo.Context) *Client {
	var receiver *Client
	//look up sender in database
	_, ok := Connections[client]

	if !ok { //map result doesn't exist
		for client := range Clients { //search the connected clients for receiver
			if client.ID == message.ReceiverID {
				receiver = client
				break
			}
		}
		if receiver == nil { //receiver not connected to server yet
			client.ReceiverNotConnected(message, c)
			return nil
		}
		slice := make([]*Client, 5)
		slice[0] = receiver
		Connections[client] = slice
	} else { //map result exists
		slice := Connections[client]
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
				client.ReceiverNotConnected(message, c)
				return nil
			}
			Connections[client] = append(slice, receiver) //add receiver to map result
		}
	}
	return receiver
}

//Deal with unknown message type
func (client *Client) WrongMsgType(message *Message, c echo.Context) {
	c.Logger().Debug("ChatServer$: WrongMsgType")
}

//Deal with the case when receiver of the message has't connected to the server
func (client *Client) ReceiverNotConnected(message *Message, c echo.Context) {
	c.Logger().Debug("ChatServer$: ReceiverNotConnected")
}
