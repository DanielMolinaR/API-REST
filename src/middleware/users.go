package middleware

import (
	"TFG/API-REST/src/lib"
	"TFG/API-REST/src/structures"
	"strconv"
	"strings"
)

func verifyLogin(newUser structures.Users) bool{
	//The DNI must has 9 char
	if len(newUser.DNI)!=9{
		return false
	} else {
		//The last char of the DNI must be a Letter
		if !verifyLastCharIsALetter(newUser){
			return false
		} else {
			//Verify if the Letter is correct
			//with the numbers of the DNI
			if !verifyLetterIsCorrect(newUser){
				return false
			} else {
				//Check if the DNI exists in the DB
				if !checkIfDniExistsAndPassswordIsCorrect(newUser){
					return false
				}
			}
		}
	}
	return true
}

func verifyLastCharIsALetter(newUser structures.Users) bool{
	//Take the last char
	c := strings.ToUpper(newUser.DNI[8:])
	//Verified if the last char is a Letter
	// parsing it to and int and using ASCII
	asciiValue := int(c[0])
	if asciiValue < 65 || asciiValue > 90 {
		return false
	} else {
		return true
	}
}

func verifyLetterIsCorrect (newUser structures.Users) bool {
	//Parse to int the DNI except the last char
	_, err := strconv.Atoi(newUser.DNI[0:8])
	if err!=nil{
		return false
	}
	//The module of the division of thenumber of the DNI
	// by 23, must be the postion of the Letter in dniLetter[]
	//This is a rule established by Spain
	//dniNumber%23 matcheado con la tabla de los modulos, arrays de letras con sus restos
	return true
}

func checkIfDniExistsAndPassswordIsCorrect(newUser structures.Users) bool{
	//Conect to the DB
	db := lib.ConectToDB()
	sqlStatement := "SELECT dni, password FROM employee WHERE dni = " + newUser.DNI
	//Do the query which return a bool and rows of data
	if bool, rows := lib.SelectQuery(db, sqlStatement); !bool{
		return false
	} else {
		//Check if the password is correct
		u := structures.Users{}
		for rows.Next(){
			if rows.Scan(&u.Password); u.Password == newUser.Password {
				return true
			}
		}
	}
	db.Close()
	return false
}