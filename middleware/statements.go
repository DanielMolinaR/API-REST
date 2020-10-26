package middleware

import (
	. "TFG/API-REST/lib"
	"TFG/API-REST/structures"
)

func checkIfPassswordIsCorrect(dni , insertedPwd string) (bool, string){

	sqlStatement := "SELECT password FROM users WHERE dni=$1"
	//DO the select and return the password
	_ = SelectQueryPwd(sqlStatement, dni)
	//Check if the password is correct
	return false, "Contraseña incorrecta"
}


func checkIfExists (data string) bool {

	sqlStatement := "SELECT dni FROM users WHERE dni = $1"
	//Do the query which return a bool if exists
	if !SelectQuery(sqlStatement, data){
		return false
	}
	return true
}

func DoEmployeeInsert(employee structures.Employee) bool {

	sqlStatement := "INSERT INTO employee (active, admin, dni, email, password, Name, Surname, phone) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	var answer bool
	if !answer{
		return false
	}
	response := InsertEmployeeQuery(sqlStatement, employee)
	return response
}

func DoPatientInsert(patient structures.Patient) bool {

	sqlStatement := "INSERT INTO patients (age, dni, email, password, Name, Surname, phone) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"
	response := InsertPatientQuery(sqlStatement, patient)
	return response
}

func getUserName(data, condition string) string {

	sqlStatement := "SELECT * FROM users WHERE " + condition +" = $1"
	userData := SelectUserDataQuery(sqlStatement, data)
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

	sqlStatement := "SELECT " + condition + " FROM patients WHERE " + condition +" = $1"
	//Do the query which return a bool if exists
	TerminalLogger.Trace("Checking if the user is a patient or not")
	DocuLogger.Trace("Checking if the user is a patient or not")
	return SelectQuery(sqlStatement, data)
}

func isTrabajador(data, condition string) (bool, string) {

	sqlStatement := "SELECT admin, active FROM employee WHERE " + condition +" = $1"
	//Do the query which return a bool if exists
	if active, admin := SelectEmployeeDataQuery(sqlStatement, data); active{
		if admin {
			return true, "admin"
		}
		return true, "employee"
	}
	return false, ""
}

