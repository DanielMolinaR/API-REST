package lib

import (
	"TFG/structures"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"os"
)

type Conection struct{
	Host string `json:"Host"`
	Port uint16 `json:"Port"`
	User string `json:"User"`
	Password string `json:"Password"`
	Database string `json:"Database"`
	Uri		 string `json:"Uri"`
}

var db *pgxpool.Pool

func init() {
	dataconfig, err := os.Open("./lib/conection.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonBody, _ := ioutil.ReadAll(dataconfig)
	var conection Conection
	json.Unmarshal(jsonBody, &conection)

	//Create the connection pool
	db, err = pgxpool.Connect(context.Background(), conection.Uri)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec(context.Background(), "SET time zone gmt")
	if err != nil {
		TerminalLogger.Info("Something went wrong setting the time zone", err)
		DocuLogger.Info("Something went wrong setting the time zone", err)
	}
}

func SelectQuery(sqlStatement, data string) (bool) {
	var response string

	rows := db.QueryRow(context.Background(), sqlStatement, data).Scan(&response)
	if rows != nil {
		if response == ""{
			return false
		} else {
			TerminalLogger.Warn("Error with the query: ", rows.Error())
			DocuLogger.Warn("Error with the query: ", rows.Error())
			return true
		}
	}
	TerminalLogger.Warn("The data has been found in the DDBB")
	DocuLogger.Warn("The data has been found in the DDBB")
	return true
}

func InsertEmployeeQuery(sqlStatement string, employee structures.Employee) (bool){
	_, err := db.Exec(context.Background(), sqlStatement, employee.Active, employee.Admin, employee.User.DNI,
		employee.User.Email, employee.User.Name, employee.User.Surname, employee.User.Phone)
	if err != nil {
		TerminalLogger.Info("Something went wrong", err)
		DocuLogger.Info("Something went wrong", err)
		return false
	}
	return true
}

func SelectStringQuery(sqlStatement, data string) (bool, string){
	var response string

	err := db.QueryRow(context.Background(), sqlStatement, data).Scan(&response)
	if err != nil {
		return false, " "
	}
	TerminalLogger.Warn("The data has been found in the DDBB")
	DocuLogger.Warn("The data has been found in the DDBB")
	return true, response
}

func InsertPatientQuery(sqlStatement string, patient structures.Patient) (bool){
	_, err := db.Exec(context.Background(), sqlStatement, patient.Birthdate, patient.User.DNI, patient.User.Email,
		patient.User.Name, patient.User.Surname, patient.User.Phone)
	if err != nil {
		TerminalLogger.Error("Something went wrong inserting the patient", err)
		DocuLogger.Error("Something went wrong inserting the patient", err)
		return false
	}
	return true
}

func UpdatePatientQuery(sqlStatement string, patient structures.Patient) (bool){
	_, err := db.Exec(context.Background(), sqlStatement, patient.Birthdate, patient.User.DNI, patient.User.Email,
		patient.User.Name, patient.User.Surname, patient.User.Phone, patient.User.Email)
	if err != nil {
		TerminalLogger.Error("Something went wrong updating the patient", err)
		DocuLogger.Error("Something went wrong updating the patient", err)
		return false
	}
	return true
}

func DoInsertClinicalbackgroundQuery(sqlStatement, dni string, data structures.ClinicalBackgroundData) bool {
	_, err := db.Exec(context.Background(), sqlStatement, dni, data)
	if err != nil {
		TerminalLogger.Info("Something went wrong", err)
		DocuLogger.Info("Something went wrong", err)
		return false
	}
	return true
}

func DoDeleteUserQuery(sqlStatement, dni string) bool{
	_, err := db.Exec(context.Background(), sqlStatement, dni)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func DoSelectExpTimeFromUniqueUrl(sqlStatement, uuid string) (int64){

	var expTime int64

	rows := db.QueryRow(context.Background(), sqlStatement, uuid).Scan(&expTime)
	if rows == nil {
		return expTime
	} else {
		TerminalLogger.Error("Error with the query: ", rows.Error())
		DocuLogger.Error("Error with the query: ", rows.Error())
		return 0
	}
}

func DoSelectOneString(sqlStatement, uuid string) string{

	var value string

	rows := db.QueryRow(context.Background(), sqlStatement, uuid).Scan(&value)
	if rows == nil {
		return value
	} else {
		TerminalLogger.Error("Error with the query: ", rows.Error())
		DocuLogger.Error("Error with the query: ", rows.Error())
		return ""
	}
}

func DoUniqueUrlTableInsert(sqlStatement, uuid, expTime, userId, email string) bool{
	_, err := db.Exec(context.Background(), sqlStatement, uuid, expTime, userId, email)
	if err != nil {
		TerminalLogger.Error("Something went wrong inserting the UUID", err)
		DocuLogger.Error("Something went wrong inserting the UUID", err)
		return false
	}
	return true
}

func DoDeleteUuidRow(sqlStatement, uuid string) bool{
	_, err := db.Exec(context.Background(), sqlStatement, uuid)
	if err != nil {
		TerminalLogger.Error("Something went wrong deleting the UUID", err)
		DocuLogger.Error("Something went wrong deleting the UUID", err)
		return false
	}
	return true
}

func DoUpdateExpTime(sqlStatement, newExpTime, uuid string) bool{
	_, err := db.Exec(context.Background(), sqlStatement, newExpTime, uuid)
	if err != nil {
		TerminalLogger.Error("Something went wrong updating the expiration_date", err)
		DocuLogger.Error("Something went wrong updating the expiration_date", err)
		return false
	}
	return true
}

func InsertNewUserQuery(sqlStatement, dni, email, name, phone, surname string) bool {
	_, err := db.Exec(context.Background(), sqlStatement, dni, email, name, phone, surname)
	if err != nil {
		TerminalLogger.Error("Something went wrong inserting the new random user", err)
		DocuLogger.Error("Something went wrong inserting the new random user", err)
		return false
	}
	TerminalLogger.Error("New random user successfully inserted", err)
	DocuLogger.Error("New random user successfully inserted", err)
	return true
}

func InsertAppointmentQuery(sqlStatement, date, employee_dni, patient_dni string) bool {
	_, err := db.Exec(context.Background(), sqlStatement, date, employee_dni, patient_dni )
	if err != nil {
		TerminalLogger.Error("Something went wrong inserting the appointment", err)
		DocuLogger.Error("Something went wrong inserting the appointment", err)
		return false
	}
	return true
}

func InsertExerciseQuery(sqlStatement, date, patient_dni, name, description string) bool {
	_, err := db.Exec(context.Background(), sqlStatement, date, patient_dni, name, description)
	if err != nil {
		TerminalLogger.Error("Something went wrong inserting the appointment", err)
		DocuLogger.Error("Something went wrong inserting the appointment", err)
		return false
	}
	return true
}

func CheckIfIsAvailable(sqlStatement, dni, date string) bool {
	var dateRetrieved int64

	rows := db.QueryRow(context.Background(), sqlStatement, dni, date).Scan(&dateRetrieved)
	if rows == nil {
		TerminalLogger.Warn("An appointment as the same hour has been found in the DDBB")
		DocuLogger.Warn("An appointment as the same hour has been found in the DDBB")
		return false
	}
	return true
}

func GetRowsFromADniQuery(sqlStatement, dni, date string) (bool, pgx.Rows) {
	rows, err := db.Query(context.Background(), sqlStatement, dni, date)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, rows
}

func GetRowsFromQuery(sqlStatement string) (bool, pgx.Rows) {
	rows, err := db.Query(context.Background(), sqlStatement)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, rows
}

func DeleteAppointmentQuery(sqlStatement, dni, date string) bool {
	_, err := db.Exec(context.Background(), sqlStatement, dni, date)
	if err != nil {
		TerminalLogger.Error("Something went wrong deleting the appointment", err)
		DocuLogger.Error("Something went wrong deleting the appointment", err)
		return false
	}
	return true
}

func DeleteExerciseQuery(sqlStatement, dni, date string) bool {
	_, err := db.Exec(context.Background(), sqlStatement, dni, date)
	if err != nil {
		TerminalLogger.Error("Something went wrong deleting the exercise", err)
		DocuLogger.Error("Something went wrong deleting the exercise", err)
		return false
	}
	return true
}

func GetClinicalBackgroundQuery(sqlStatement, dni string) (bool, structures.ClinicalBackgroundData) {

	var clinicalBackground structures.ClinicalBackgroundData

	rows := db.QueryRow(context.Background(), sqlStatement, dni).Scan(&clinicalBackground)
	if rows == nil {
		return true, clinicalBackground
	} else {
		TerminalLogger.Error("Error with the query: ", rows.Error())
		DocuLogger.Error("Error with the query: ", rows.Error())
		return false, clinicalBackground
	}
}

func UpdateClinicalBackgroundQuery(sqlStatement, dni string, data structures.ClinicalBackgroundData) bool {
	_, err := db.Exec(context.Background(), sqlStatement, data, dni)
	if err != nil {
		TerminalLogger.Error("Something went wrong updating the clinical background", err)
		DocuLogger.Error("Something went wrong updating the clinical background", err)
		return false
	}
	return true
}

func DoEmployeeUpgrade(sqlStatement, dni string) (bool, map[string]interface{}){
	_, err := db.Exec(context.Background(), sqlStatement, dni)
	if err != nil {
		TerminalLogger.Error("Something went wrong upgrading the employee", err)
		DocuLogger.Error("Something went wrong upgrading the employee", err)
		return false, map[string]interface{}{"state": "No se ha podido dar permisos de administrador al empleado"}
	}
	return true, map[string]interface{}{"state": ""}
}

func DoEmployeeLayOff(sqlStatement, dni string) (bool, map[string]interface{}){
	_, err := db.Exec(context.Background(), sqlStatement, dni)
	if err != nil {
		TerminalLogger.Error("Something went wrong laying off the employee", err)
		DocuLogger.Error("Something went wrong lYING OFF the employee", err)
		return false, map[string]interface{}{"state": "No se ha podido dar de baja al empleado"}
	}
	return true, map[string]interface{}{"state": "Empleado suspendido"}
}

func DoEmployeeRenew(sqlStatement, dni string) (bool, map[string]interface{}){
	_, err := db.Exec(context.Background(), sqlStatement, dni)
	if err != nil {
		TerminalLogger.Error("Something went wrong renewing the employee", err)
		DocuLogger.Error("Something went wrong renewing the employee", err)
		return false, map[string]interface{}{"state": "No se ha podido renovar al empleado"}
	}
	return true, map[string]interface{}{"state": "Empleado Renovado"}
}