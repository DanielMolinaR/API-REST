package middleware

import (
	. "TFG/API-REST/src/lib"
	"TFG/API-REST/src/structures"
)

//Conect to the DB
var db = ConectToDB()

func checkIfDniExistsAndPassswordIsCorrect(dni, password string) bool{
	sqlStatement := "SELECT dni, password FROM users WHERE dni = " + dni
	//Do the query which return a bool and rows of data
	if bool, rows := SelectQuery(db, sqlStatement); !bool{
		return false
	} else {
		//Check if the password is correct
		u := structures.Users{}
		for rows.Next(){
			if rows.Scan(&u.Password); u.Password == password {
				return true
			}
		}
		return false
	}
}

func DoInsert(sqlStatement string) bool {
	response := InsertQuery(db, sqlStatement)
	return response
}