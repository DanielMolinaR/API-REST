package middleware

import (
	"TFG/API-REST/src/structures"
	"encoding/json"
)


func UsersLogin (reqBody []byte) bool{
	var newUser structures.Users
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newUser)
	//verifyLogin verify if all params are correct
	// and if they exist in the DB
	if !verifyLogin(newUser){
		return false
	}
	return true
}

