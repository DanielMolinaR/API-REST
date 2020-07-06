package middleware

import (
	"encoding/json"
)

type users struct{
	DNI         string `json:"DNI"`
	Email       string `json:"Email"`
	password 	string `json:"password"`
	Name 		string `json:"Name"`
	Surname 	string `json:"Surname""`
	phone		string	`json:"phone"`
}

type client struct{
	user 		users
	age 		int 	`json:"age"`
}

type employee struct{
	user		users
	active		bool	`json:"active"`

}

func UsersLogin (reqBody []byte) bool{
	var newUser users
	json.Unmarshal(reqBody, &newUser)
	if !verifyLogin(newUser){
		return false
	}
	return true
}

