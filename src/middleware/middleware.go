package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func Homelink (w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome home!")
}

func Login (w http.ResponseWriter, r *http.Request) {
	var newUser users
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the user dni, email, password and name only in order to update")
	} else{
		UsersLogin(newUser, reqBody)
	}

}


func SignIn (w http.ResponseWriter, r *http.Request){

}