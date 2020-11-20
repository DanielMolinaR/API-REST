package middleware

import (
	. "TFG/API-REST/lib"
	"TFG/API-REST/structures"
	"time"
)

func checkIfExists (condition, dataToSelect, data string) bool {

	sqlStatement := "SELECT " + dataToSelect + " FROM " + condition +" WHERE "+ dataToSelect +" = $1"
	//Do the query which return a bool if exists
	if !SelectQuery(sqlStatement, data){
		TerminalLogger.Info("The " + dataToSelect +" doesnt exists in the DDBB")
		DocuLogger.Info("The " + dataToSelect +" doesnt exists in the DDBB")
		return false
	}
	return true
}

func doEmployeeInsert(employee structures.Employee) bool {

	sqlStatement := "INSERT INTO Employee (active, admin, dni, email, name, surname, phone) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if !InsertEmployeeQuery(sqlStatement, employee){
		return false
	} else {
		return createKeycloakUser(employee.User.DNI, employee.User.Password, "EMPLOYEE_ROLE")
	}
}

func doPatientInsert(patient structures.Patient) bool {

	if !createKeycloakUser(patient.User.DNI, patient.User.Password, "PATIENT_ROLE"){
		return false
	} else {
		sqlStatement := "INSERT INTO patients (birthdate, dni, email, name, surname, phone) " +
			"VALUES ($1, $2, $3, $4, $5, $6)"
		response := InsertPatientQuery(sqlStatement, patient)
		return response
	}
}

func getUuid(uuid string) (string, string){
	sqlStatement := "SELECT * FROM UniqueUrl WHERE (uuid = $1)"
	return DoSelectUuid(sqlStatement, uuid)
}

func insertUuidAndExpTime(uuid string) bool{
	sqlStatement := "INSERT INTO UniqueUrl VALUES ($1, $2);"
	expTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day() + 3, time.Now().Hour(), time.Now().Minute(), time.Now().Second(),0,time.Local)
	return DoUuidInsert(sqlStatement, uuid, expTime.String()[:20])
}

func deleteUuidRow(uuid string) {
	sqlStatement := "DELETE FROM UniqueUrl WHERE (uuid = $1)"
	DoDeleteUuidRow(sqlStatement, uuid)
}



