package middleware

import (
	. "TFG/API-REST/lib"
	"TFG/API-REST/structures"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jinzhu/now"
	"strconv"
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

func createClinicalBackground(dni string) bool {
	var data structures.ClinicalBackgroundData
	sqlStatement := "INSERT INTO ClinicalBackground (dni_patients, clinical_background_data) VALUES ($1, $2)"
	return DoInsertClinicalbackgroundQuery(sqlStatement, dni, data)
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
	sqlStatment := "INSERT INTO Appointments (date_time, dni_Employee, dni_Patients) " +
		"VALUES ($1, $2, $3)"
	return InsertAppointmentQuery(sqlStatment, date, employee_dni, patient_dni)
}

func insertExercise( date, patient_dni, name, description string) bool{
	sqlStatment := "INSERT INTO Exercise (exercise_date_time, dni_Patients, name, description) " +
		"VALUES ($1, $2, $3, $4)"
	return InsertExerciseQuery(sqlStatment, date, patient_dni, name, description)
}

func checkIfAvailable(condition, dni, date, table string) bool{
	sqlStatement := "SELECT EXTRACT('epoch' from date_time) FROM " + table  + " WHERE dni_" + condition + " = $1 and date_time = $2"
	return CheckIfIsAvailable(sqlStatement, dni, date)
}

func getAppointmentsFromDB(dni string) (bool, pgx.Rows) {
	sqlStatement := "SELECT EXTRACT('epoch' from date_time), employee.name AS employee_name, patients.name AS patient_name FROM " +
		"appointments, employee, patients WHERE (appointments.dni_employee = $1 or appointments.dni_patients = $1) " +
		"AND (appointments.dni_employee = employee.dni) AND (appointments.dni_patients =  patients.dni)"
	return GetRowsFromADniQuery(sqlStatement, dni)
}

func getAllAppointmentsFromDB() (bool, pgx.Rows) {
	sqlStatement := "SELECT EXTRACT('epoch' from date_time), employee.name AS employee_name, patients.name AS patient_name FROM " +
		"appointments, employee, patients WHERE (appointments.dni_employee = employee.dni) AND (appointments.dni_patients =  patients.dni)"
	return GetRowsFromQuery(sqlStatement)
}

func getExercisesFromDB(dni string) (bool, pgx.Rows) {
	sqlStatement := "SELECT EXTRACT('epoch' from exercise_date_time), description, name FROM Exercise WHERE dni_patients = $1"
	return GetRowsFromADniQuery(sqlStatement, dni)
}

func getAppointmentsDataFromRows(rows pgx.Rows) map[string]map[string]interface{} {
	appointments := make(map[string]map[string]interface{})

	var appointmentDataResponse structures.AppointmentResponse
	rowsCount := 0
	for rows.Next() {

		var date int64
		err := rows.Scan(&date, &appointmentDataResponse.Employee_name, &appointmentDataResponse.Patient_name)
		if err != nil {
			fmt.Println(err)
			return appointments
		}
		dateAsSomething := time.Unix(date, 0)
		finalDate := time.Date(dateAsSomething.Year(), dateAsSomething.Month(), dateAsSomething.Day(), dateAsSomething.Hour()-1, dateAsSomething.Minute(),
			dateAsSomething.Second(), 0, time.UTC)
		appointmentDataResponse.Date = finalDate.String()[:20]
		rowsCount += 1
		appointments["Cita " + strconv.Itoa(rowsCount)] = map[string]interface{}{"date": appointmentDataResponse.Date, "Patient_name": appointmentDataResponse.Patient_name,
			"Employee_name": appointmentDataResponse.Employee_name}

	}

	if rows.Err() != nil {
		fmt.Println(rows.Err())
		return appointments
	}
	return appointments
}

func getExercisesDataFromRows(rows pgx.Rows) map[string]map[string]interface{} {
	exercises := make(map[string]map[string]interface{})

	var exerciseDataResponse structures.ExercisesResponse
	rowsCount := 0
	for rows.Next() {

		var date int64
		err := rows.Scan(&date, &exerciseDataResponse.Description, &exerciseDataResponse.Exercise_name)
		if err != nil {
			fmt.Println(err)
			return exercises
		}
		dateAsSomething := time.Unix(date, 0)
		finalDate := time.Date(dateAsSomething.Year(), dateAsSomething.Month(), dateAsSomething.Day(), dateAsSomething.Hour()-1, dateAsSomething.Minute(),
			dateAsSomething.Second(), 0, time.UTC)
		exerciseDataResponse.Date = finalDate.String()[:20]
		rowsCount += 1
		exercises["Ejercicio " + strconv.Itoa(rowsCount)] = map[string]interface{}{"date": exerciseDataResponse.Date, "Description": exerciseDataResponse.Description,
			"Name": exerciseDataResponse.Exercise_name}

	}

	if rows.Err() != nil {
		fmt.Println(rows.Err())
		return exercises
	}
	return exercises
}

func deleteAppointmentFromDB(dni, date string) bool{
	sqlStatement := "DELETE FROM appointments WHERE (dni_employee = $1 or dni_patients = $1) AND date_time = $2"
	return DeleteAppointmentQuery(sqlStatement, dni, date)
}

func deleteExerciseFromDB(dni, date string) bool{
	sqlStatement := "DELETE FROM exercise WHERE dni_patients = $1 AND date_time = $2"
	return DeleteExerciseQuery(sqlStatement, dni, date)
}

func getClinicalBackground(data structures.ClinicalBackgroundData) (bool, structures.ClinicalBackgroundData) {
	sqlStatement := "SELECT clinical_background_data FROM clinicalbackground WHERE (clinical_background_dni_patients = $1)"
	return GetClinicalBackgroundQuery(sqlStatement, data.Patient_dni)
}

func updateClinicalBackground(data structures.ClinicalBackgroundData) bool{
	sqlStatement := "UPDATE clinicalbackground SET clinical_background_data = $1 WHERE (clinical_background_dni_patients = $2)"
	return UpdateClinicalBackgroundQuery(sqlStatement, data.Patient_dni, data)
}

func getAllEmployeesFromDB() (bool, pgx.Rows) {
	sqlStatement := "SELECT * FROM Employee"
	return GetRowsFromQuery(sqlStatement)
}

func getEmployeeDataFromRows(rows pgx.Rows) map[string]map[string]interface{} {
	employees := make(map[string]map[string]interface{})

	var employeeDataResponse structures.Employee
	rowsCount := 0
	for rows.Next() {

		err := rows.Scan(&employeeDataResponse.User.DNI, &employeeDataResponse.User.Email, &employeeDataResponse.User.Name, &employeeDataResponse.User.Surname,
					&employeeDataResponse.User.Phone, &employeeDataResponse.Active, &employeeDataResponse.Admin)
		if err != nil {
			fmt.Println(err)
			return employees
		}
		rowsCount += 1
		employees["Empleado " + employeeDataResponse.User.Name] = map[string]interface{}{"DNI": employeeDataResponse.User.DNI,
				"Email": employeeDataResponse.User.Email, "Name": employeeDataResponse.User.Name, "Surname": employeeDataResponse.User.Surname,
				"Phone": employeeDataResponse.User.Phone, "Active": employeeDataResponse.Active, "Admin": employeeDataResponse.Admin}

	}

	if rows.Err() != nil {
		fmt.Println(rows.Err())
		return employees
	}
	return employees
}

func getAllPatientsFromDB() (bool, pgx.Rows) {
	sqlStatement := "SELECT * FROM Patients"
	return GetRowsFromQuery(sqlStatement)
}

func getPatientDataFromRows(rows pgx.Rows) map[string]map[string]interface{} {
	patients := make(map[string]map[string]interface{})

	var patientDataResponse structures.Patient
	rowsCount := 0
	for rows.Next() {

		err := rows.Scan(&patientDataResponse.User.DNI, &patientDataResponse.User.Email, &patientDataResponse.User.Name, &patientDataResponse.User.Surname,
			&patientDataResponse.User.Phone, &patientDataResponse.Birthdate)
		if err != nil {
			fmt.Println(err)
			return patients
		}
		rowsCount += 1
		patients["Paciente " + patientDataResponse.User.Name] = map[string]interface{}{"DNI": patientDataResponse.User.DNI,
			"Email": patientDataResponse.User.Email, "Name": patientDataResponse.User.Name, "Surname": patientDataResponse.User.Surname,
			"Phone": patientDataResponse.User.Phone, "Fecha de nacimiento": patientDataResponse.Birthdate}

	}

	if rows.Err() != nil {
		fmt.Println(rows.Err())
		return patients
	}
	return patients
}

func getAllEmployeeDnis() (bool, pgx.Rows) {
	sqlStatement := "SELECT dni, email FROM Employee"
	return GetRowsFromQuery(sqlStatement)
}

func getAllApointmentsOfTheDay(employeeDni string) (bool, pgx.Rows) {
	beginDay := now.BeginningOfDay()
	endDay := now.EndOfDay()
	beginDayString := beginDay.String()[:19]
	endDayString := endDay.String()[:19]
	sqlStatement := "SELECT EXTRACT('epoch' from date_time), patients.name AS patient_name FROM appointments, patients" +
		" WHERE (appointments.dni_employee = '" + employeeDni + "') AND (appointments.dni_patients = patients.dni)" +
		"AND date_time between '" + beginDayString + "' and '" + endDayString + "'"
	return GetRowsFromQuery(sqlStatement)
}