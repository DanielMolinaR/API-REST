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
	password := SelectQueryPwd(db, sqlStatement,dni)
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
	sqlStatement = "SELECT " + condition + " FROM users WHERE " + condition +" = $1"
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

func getUserName(data, condition string) string {
	var db = ConectToDB()
	sqlStatement := "SELECT * FROM users WHERE " + condition +" = $1"
	userData := SelectUserDataQuery(db, sqlStatement, data)
	response := userData["name"].(string)
	return response
}

func getUserId(data, condition string) string{
	if isPatient(data, condition){
		TerminalLogger.Info("The user is a patient")
		DocuLogger.Info("The user is a patient")
		return "patient"
	}
	if bool, response := isTrabajador(data, condition); bool{
		TerminalLogger.Info("The user is", response)
		DocuLogger.Info("The user is", response)
		return response
	}
	return "employee suspended"
}

func isPatient(data, condition string) bool {
	var db = ConectToDB()
	sqlStatement := "SELECT " + condition + " FROM patients WHERE " + condition +" = $1"
	//Do the query which return a bool if exists
	TerminalLogger.Trace("Checking if the user is a patient or not")
	DocuLogger.Trace("Checking if the user is a patient or not")
	return SelectQuery(db, sqlStatement, data)
}

func isTrabajador(data, condition string) (bool, string) {
	var db = ConectToDB()
	sqlStatement := "SELECT admin, active FROM employee WHERE " + condition +" = $1"
	//Do the query which return a bool if exists
	if active, admin := SelectEmployeeDataQuery(db, sqlStatement, data); active{
		if admin {
			return true, "admin"
		}
		return true, "employee"
	}
	return false, ""
}
