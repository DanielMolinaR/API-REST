package structures

type Users struct{
	DNI         string `json:"DNI"`
	Email       string `json:"Email"`
	Password 	string `json:"Password"`
	Name 		string `json:"Name"`
	Surname 	string `json:"Surname""`
	Phone		string	`json:"Phone"`
}

type Patient struct{
	User Users
	Age  string 	`json:"Age"`
}

type Employee struct{
	User   Users	`json:"data"`
	Active string	`json:"Active"`
	Admin  string	`json:"Admin"`
}
