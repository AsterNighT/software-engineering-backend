package chat

import "time"

type Chat struct {
	ID			uint				
	URL	string						//The url of the chatting window
	DoctorID	uint
	PatientID	uint

	//primary key: ID
}


type Message struct{
	ID			uint
	ChatID		uint
	Timestamp	Time
	Type		uint				//Texts, pictures....
	Content		string
	PictureURL  string 

	//primary key: ID&chatID
}


//Question catagories
//For standard query,each catagory correspond with a set of fixed questions
type Catagory struct{
	ID			uint
	Name		string
	Keywords	[]string
	Questions	[]string	

	//primary key: ID
}

