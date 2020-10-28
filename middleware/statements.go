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


func checkIfExists (condition, data string) bool {

	sqlStatement := "SELECT dni FROM " + condition +" WHERE dni = $1"
	//Do the query which return a bool if exists
	if !SelectQuery(sqlStatement, data){
		return false
	}
	return true
}

func DoEmployeeInsert(employee structures.Employee) bool {

	if !CreateUser(employee.User.DNI, employee.User.Password, "EMPLOYEE_ROLE"){
		return false
	} else {
		sqlStatement := "INSERT INTO Employee (active, admin, dni, email, name, surname, phone) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7)"
		response := InsertEmployeeQuery(sqlStatement, employee)
		return response
	}
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


