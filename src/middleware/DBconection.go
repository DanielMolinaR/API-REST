package middleware

import (
	. "TFG/API-REST/src/lib"
	"TFG/API-REST/src/structures"
)

func checkIfPassswordIsCorrect(dni, p string) (bool, string){
	//Conect to the DB
	var db = ConectToDB()

	sqlStatement := "SELECT password FROM users WHERE dni=$1"
	//DO the select and return the password
	_, password := SelectQuery(db, sqlStatement,dni)
	//Check if the password is correct
	if ComparePwdAndHash([]byte(p), []byte(password)) {
		return true, "Sesion iniciada"
	}
	return false, "Contraseña incorrecta"
}


func checkIfExists (data string) bool {
	var db = ConectToDB()
	sqlStatement := "SELECT dni FROM users WHERE dni = $1"
	//Do the query which return a bool and rows of data
	if bool, _ := SelectQuery(db, sqlStatement, data); !bool{
		return false
	}
	return true
}

func DoEmployeeInsert(employee structures.Employee) bool {
	var db = ConectToDB()
	sqlStatement := "INSERT INTO employee (active, admin, dni, email, password, Name, Surname, phone) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	response := InsertEmployeeQuery(db, sqlStatement, employee)
	return response
}

func DoPatientInsert(patient structures.Patient) bool {
	var db = ConectToDB()
	sqlStatement := "INSERT INTO patients (age, dni, email, password, Name, Surname, phone) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"
	response := InsertPatientQuery(db, sqlStatement, patient)
	return response
}