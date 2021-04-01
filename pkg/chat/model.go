package chat

type Chat struct {
	ID        uint   `gorm:"primaryKey"`
	URL       string //The url of the chatting window
	DoctorID  uint
	PatientID uint
}

type Message struct {
	ID        uint `gorm:"primaryKey"`
	Chat      Chat //Default foreign key:ID in Chat
	CreatedAt int
	Type      uint   //Texts, pictures....
	Content   string //Use escaped string for PictureURL
}

/*
Question catagories
-------------------------------------------------------------------
If the patient's message contains some keyword specified in DB,
then we can find the matching catagory of questions.
Doctor's work can be reduced by automatically showing right questions.
*/
type Catagory struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Keywords  []Keyword `gorm:"many2many:catagory_keywords;"`
	Questions *string   //Each catagory correspond with a set of fixed questions
}

type Keyword struct {
	ID   uint
	Name string
}
