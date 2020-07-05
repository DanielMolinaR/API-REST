package middleware

import "encoding/json"

func UsersLogin (newUser users, reqBody []byte) bool{
	json.Unmarshal(reqBody, &newUser)
	if !dniIncorrect(newUser){
		return false
	}
}

func dniIncorrect(newUser users) bool{
	if len(newUser.DNI)!=9{
		panic("DNI incorrecto")
		return false
	}
	//if dni no existe return false
	return true
}