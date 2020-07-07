package structures

type Users struct{
	DNI         string `json:"DNI"`
	Email       string `json:"Email"`
	Password 	string `json:"password"`
	Name 		string `json:"Name"`
	Surname 	string `json:"Surname""`
	Phone		string	`json:"phone"`
}

type Patients struct{
	User Users
	Age  int 	`json:"age"`
}

type Employee struct{
	User   Users
	Active string	`json:"active"`
	Admin  string	`json:"active"`
}
