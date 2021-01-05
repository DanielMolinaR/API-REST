package lib

import (
	"TFG/API-REST/structures"
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
}

var db *pgxpool.Pool

func init() {
	dataconfig, err := os.Open("./API-REST/lib/conection.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonBody, _ := ioutil.ReadAll(dataconfig)
	var conection Conection
	json.Unmarshal(jsonBody, &conection)
	//Set the params to connect to the DB
	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.Host = conection.Host
	config.ConnConfig.Port = conection.Port
	config.ConnConfig.User = conection.User
	config.ConnConfig.Password = conection.Password
	config.ConnConfig.Database = conection.Database
	//Create the connection pool
	db, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Println(err)
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

func SelectEmployeeDataQuery(sqlStatement, data string) structures.Employee {
	var employee structures.Employee
	//Do the query and if It's correct
	//It means that the password is saved
	err := db.QueryRow(context.Background(), sqlStatement, data).Scan(&employee.User.DNI, &employee.User.Name, &employee.User.Email,
		&employee.User.Phone, &employee.User.Surname, &employee.Active, employee.Admin)
	if err != nil {
		TerminalLogger.Error("Error with the query:", err)
		DocuLogger.Error("Error with the query:", err)
		return employee
	}

	return employee
}

/*func selectAppointmentsQuery(db *pgxpool.Pool, sqlStatement, data string) (bool, pgx.Rows) {

	rows, err := db.Query(context.Background(), sqlStatement, data)

	if err != nil{
		TerminalLogger.Warn("Error with the query:", err)
		DocuLogger.Warn("Error with the query:", err)
		return false, ""
	}

	//    for rows.Next() {
	//        var n int32
	//        err = rows.Scan(&n) escaneamos cada parametro de la base de datos
	//        if err != nil {
	//            return err
	//        }
	//		  la fecha (timestamp) la descomponemos en dia y hora
	//        sum += n  añadimos los datos en una biblioteca o array y devolvemos esta
	//    }

}*/

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

	rows := db.QueryRow(context.Background(), sqlStatement, data).Scan(&response)
	if rows != nil {
		return false, ""
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

func InsertNewUserQuery(sqlStatement, dni, email, name, phone string) bool {
	_, err := db.Exec(context.Background(), sqlStatement, dni, email, name, phone)
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

func GetRowsFromADniQuery(sqlStatement, dni string) (bool, pgx.Rows) {
	rows, err := db.Query(context.Background(), sqlStatement, dni)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, rows
}

func GetAllAppointmentsAndNamesFromQuery(sqlStatement string) (bool, pgx.Rows) {
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
		TerminalLogger.Error("Something went wrong deleting the appointment", err)
		DocuLogger.Error("Something went wrong deleting the appointment", err)
		return false
	}
	return true
}