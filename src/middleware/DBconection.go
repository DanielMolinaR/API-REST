package middleware

import (
	. "TFG/API-REST/src/lib"
)

func checkIfDniExistsAndPassswordIsCorrect(dni, p string) (bool, string){
	//Conect to the DB
	var db = ConectToDB()

	sqlStatement := "SELECT password FROM users WHERE dni=$1"
	//Do the query which return a bool if
	//the DNI exists and the password saved
	if bool, password := SelectQuery(db, sqlStatement, dni); !bool{
		return false, "El DNI no existe"
	} else {
		//Check if the password is correct
		if password == p {
			return true, "Sesion iniciada"
		}
		return false, "Contraseña incorrecta"
	}
}

func checkIfEmailExistsAndPassswordIsCorrect (email, p string) (bool, string) {
	//Conect to the DB
	var db = ConectToDB()

	sqlStatement := "SELECT email, password FROM users WHERE email = " + email
	//Do the query which return a bool and rows of data
	if bool, password := SelectQuery(db, sqlStatement, email); !bool{
		return false, "El email no existe"
	} else {
		//Check if the password is correct
		if password == p {
			return true, "Sesion iniciada"
		}
		return false, "Contraseña incorrecta"
	}
}

func checkIfDniExists (dni string) bool {
	var db = ConectToDB()
	sqlStatement := "SELECT dni FROM users WHERE dni = " + dni
	//Do the query which return a bool and rows of data
	if bool, _ := SelectQuery(db, sqlStatement, dni); !bool{
		return false
	}
	return true
}

func DoInsert(sqlStatement string) bool {
	var db = ConectToDB()
	response := InsertQuery(db, sqlStatement)
	return response
}