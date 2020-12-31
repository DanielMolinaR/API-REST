package middleware

import (
	. "TFG/API-REST/lib"
	"TFG/API-REST/structures"
	"time"
)

func checkIfExists (table, dataToSelect, fieldToCompare, data string) bool {

	sqlStatement := "SELECT " + dataToSelect + " FROM " + table +" WHERE "+ fieldToCompare +" = $1"
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
		if ok, userid := createKeycloakUser(employee.User.DNI, employee.User.Password, employee.User.Email, "EMPLOYEE_ROLE"); !ok{

			//As the user could not be saved in keycloak It must be deleted in the DB
			DeleteUserStatement(employee.User.DNI)
			return ok, userid
		} else {
			return ok, userid
		}
	}
}

func getStringFromField(table, dataToSelect, fieldToCompare, dataToCompare string) (bool, string){
	 sqlStatement := "SELECT " + dataToSelect + " FROM " + table + " WHERE " + fieldToCompare + " = $1"
	 return SelectStringQuery(sqlStatement, dataToCompare)
}

func doPatientInsert(patient structures.Patient) (bool, string) {

	sqlStatement := "INSERT INTO Patients (birthdate, dni, email, name, surname, phone) " +
		"VALUES ($1, $2, $3, $4, $5, $6)"
	if !InsertPatientQuery(sqlStatement, patient){
		return false, ""
	} else {

		//if the patient has been inserted in the DB correctly now is inserted into keycloak
		if ok, userid := createKeycloakUser(patient.User.DNI, patient.User.Password, patient.User.Email, "PATIENT_ROLE"); !ok{

			//As the user couldnt be saved in keycloak It must be deleted in the DB
			DeleteUserStatement(patient.User.DNI)
			return ok, userid
		} else {
			return ok, userid
		}
	}
}

func doPatientUpdateAndInsert(patient structures.Patient) (bool, string) {

	if ok, userid := createKeycloakUser(patient.User.DNI, patient.User.Password, patient.User.Email, "PATIENT_ROLE"); !ok{
		return false, ""
	} else{
		sqlStatement := "UPDATE Patients SET birthdate = $1, dni = $2, email = $3, name = $4, surname = $5, phone = $6 " +
			"WHERE email = $7"
		if !UpdatePatientQuery(sqlStatement, patient){
			DeleteKeycloakUser(userid)
			return false, ""
		} else {
			return true, userid
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
	sqlStatement := "INSERT INTO UniqueUrl (uuid, expiration_date, user_id, email) VALUES ($1, $2, $3, $4);"
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

func insertNewRandomUser(dni, name, email, phone string) bool{
	sqlStatement := "INSERT INTO Patients (dni, email, name, phone) " +
		"VALUES ($1, $2, $3, $4)"
	return InsertNewUserQuery(sqlStatement, dni, email, name, phone)
}

func insertAppointment( date, employee_dni, patient_dni string) bool{

	/* VER LA BASE DE DATOS DE LAS CITAS */

	sqlStatment := "INSERT INTO Appointments (date_time, dni_Employee, dni_Patients) " +
		"VALUES ($1, $2, $3)"

	return InsertAppointmentQuery(sqlStatment, date, employee_dni, patient_dni)

}

func checkIfAvailable(condition, dni, date string) bool{
	sqlStatement := "SELECT EXTRACT('epoch' from date_time) FROM appointments WHERE dni_"+condition+" = $1 and date_time = $2"
	return CheckIfIsAvailable(sqlStatement, dni, date)
}
