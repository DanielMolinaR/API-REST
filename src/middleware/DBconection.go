package middleware

import (
	. "TFG/API-REST/src/lib"
	"TFG/API-REST/src/structures"
)

func checkIfPassswordIsCorrect(dni , insertedPwd string) (bool, string){
	//Conect to the DB
	var db = ConectToDB()

	sqlStatement := "SELECT password FROM users WHERE dni=$1"
	//DO the select and return the password
	_, password := SelectQueryPwd(db, sqlStatement,dni)
	//Check if the password is correct
	if ComparePwdAndHash(insertedPwd, password) {
		return true, "Sesion iniciada"
	}
	return false, "Contraseña incorrecta"
}


func checkIfExists (data, condition string) bool {
	var (
		db = ConectToDB()
		sqlStatement string
	)
	if condition == "dni"{
		sqlStatement = "SELECT dni FROM users WHERE dni = $1"
	} else {
		sqlStatement = "SELECT email FROM users WHERE email = $1"
	}
	//Do the query which return a bool if exists
	if !SelectQuery(db, sqlStatement, data){
		return false
	}
	return true
}

func DoEmployeeInsert(employee structures.Employee) bool {
	var db = ConectToDB()
	sqlStatement := "INSERT INTO employee (active, admin, dni, email, password, Name, Surname, phone) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	var err error
	employee.User.Password, err = encryptPwd(employee.User.Password)
	if err != nil{
		return false
	}
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