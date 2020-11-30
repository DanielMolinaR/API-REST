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

func doEmployeeInsert(employee structures.Employee) (bool, string) {

	sqlStatement := "INSERT INTO Employee (active, admin, dni, email, name, surname, phone) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if !InsertEmployeeQuery(sqlStatement, employee){
		return false, ""
	} else {

		//if the employee has been inserted in the DB correctly now is inserted into keycloak
		if ok, userid := createKeycloakUser(employee.User.DNI, employee.User.Password, "EMPLOYEE_ROLE"); !ok{

			//As the user could not be saved in keycloak It must be deleted in the DB
			DeleteUserStatement(employee.User.DNI)
			return ok, userid
		} else {
			return ok, userid
		}
	}
}

func doPatientInsert(patient structures.Patient) (bool, string) {

	sqlStatement := "INSERT INTO Patients (birthdate, dni, email, name, surname, phone) " +
		"VALUES ($1, $2, $3, $4, $5, $6)"
	if !InsertPatientQuery(sqlStatement, patient){
		return false, ""
	} else {

		//if the patient has been inserted in the DB correctly now is inserted into keycloak
		if ok, userid := createKeycloakUser(patient.User.DNI, patient.User.Password, "PATIENT_ROLE"); !ok{

			//As the user couldnt be saved in keycloak It must be deleted in the DB
			DeleteUserStatement(patient.User.DNI)
			return ok, userid
		} else {
			return ok, userid
		}
	}
}

func DeleteUserStatement(dni string) {
	sqlStatemente := "DELETE FROM users WHERE dni = $1"
	DoDeleteUserQuery(sqlStatemente, dni)
}

func getExpTimeFromUuid(uuid string) (int64){
	sqlStatement := "SELECT EXTRACT('epoch' from expiration_date) FROM uniqueUrl WHERE uuid = $1"
	return DoSelectExpTimeFromUniqueUrl(sqlStatement, uuid)
}

func insertUuidExpTimeAndUserId(uuid, userId, email string) bool{
	sqlStatement := "INSERT INTO UniqueUrl VALUES ($1, $2, $3, $4);"
	expTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day() + 3, time.Now().Hour(), time.Now().Minute(), time.Now().Second(),0, time.UTC)
	return DoUniqueUrlTableInsert(sqlStatement, uuid, expTime.String()[:20], userId, email)
}

func DeleteUuidRow(uuid string) {
	sqlStatement := "DELETE FROM UniqueUrl WHERE (uuid = $1)"
	DoDeleteUuidRow(sqlStatement, uuid)
}

func UpdateExpTimeFromUuid(uuid string) bool {
	sqlStatement := "UPDATE UniqueUrl SET expiration_date = $1 WHERE (uuid = $2)"
	expTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day() + 3, time.Now().Hour(), time.Now().Minute(), time.Now().Second(),0, time.UTC)
	 return DoUpdateExpTime(sqlStatement, expTime.String()[:20], uuid)
}

func getUserIdFromUuid(uuid string) string {
	sqlStatement := "SELECT user_id FROM uniqueUrl WHERE uuid = $1"
	return DoSelectOneString(sqlStatement, uuid)
}

func getEmailFromUuid(uuid string) string {
	sqlStatement := "SELECT email FROM uniqueUrl WHERE uuid = $1"
	return DoSelectOneString(sqlStatement, uuid)
}

func getNameFromUsersWithEmail(email string) string {
	sqlStatement := "SELECT name FROM users WHERE email = $1"
	return DoSelectOneString(sqlStatement, email)
}


