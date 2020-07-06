package structures

type Users struct{
	DNI         string `json:"DNI"`
	Email       string `json:"Email"`
	Password 	string `json:"password"`
	Name 		string `json:"Name"`
	Surname 	string `json:"Surname""`
	Phone		string	`json:"phone"`
}

type patients struct{
	User Users
	Age  int 	`json:"age"`
}

type employee struct{
	User   Users
	Active bool	`json:"active"`
}
