package middleware

import (
	"TFG/API-REST/src/lib"
	"strconv"
	"strings"
)

func verifyLogin(newUser users) bool{
	if len(newUser.DNI)!=9{
		return false
	} else {
		if !verifyLastCharIsALetter(newUser){
			return false
		} else {
			if !verifyLetterIsCorrect(newUser){
				return false
			} else {
				if !checkIfDniExistsAndPassswordIsCorrect(newUser){
					return false
				}
			}
		}
	}
	return true
}

func verifyLastCharIsALetter(newUser users) bool{
	c := strings.ToUpper(newUser.DNI[8:])
	asciiValue := int(c[0])
	if asciiValue < 65 || asciiValue > 90 {
		return false
	} else {
		return true
	}
}

func verifyLetterIsCorrect (newUser users) bool {
	_, err := strconv.Atoi(newUser.DNI[0:8])
	if err!=nil{
		return false
	}
	return true
}

func checkIfDniExistsAndPassswordIsCorrect(newUser users) bool{
	db := lib.ConectToDB()
	sqlStatement := "SELECT dni, password FROM clients WHERE dni = " + newUser.DNI
	if bool, rows := lib.SelectQuery2(db, sqlStatement); !bool{
		return false
	} else {
		u := users{}
		for rows.Next(){
			if rows.Scan(&u.password); u.password == newUser.password {
				return true
			}
		}
	}
	return false
}