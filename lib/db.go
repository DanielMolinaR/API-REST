package lib

import (
	"TFG/API-REST/structures"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
	TerminalLogger.Warn("The DNI has been found in the DDBB")
	DocuLogger.Warn("The DNI has been found in the DDBB")
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

func InsertPatientQuery(sqlStatement string, patient structures.Patient) (bool){
	_, err := db.Exec(context.Background(), sqlStatement, patient.Birthdate, patient.User.DNI, patient.User.Email,
		patient.User.Name, patient.User.Surname, patient.User.Phone)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func DoSelectUuid(sqlStatement, uuid string) (string, string){
	var uuidResponse, expTime string

	rows := db.QueryRow(context.Background(), sqlStatement, uuid).Scan(&uuidResponse, &expTime)
	if rows == nil {
		TerminalLogger.Error("Error with the query: ", rows.Error())
		DocuLogger.Error("Error with the query: ", rows.Error())
		return "",""
	} else {
		return uuidResponse, expTime
	}
}

func DoUuidInsert(sqlStatement string, uuid, expTime string) bool{
	_, err := db.Exec(context.Background(), sqlStatement, uuid, expTime)
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

func updateQuery(db *sql.DB){
	sqlStatement := "UPDATE employee set name=($1), surname=($2) where dni=($3)"

	_, err := db.Exec(sqlStatement, "DANIEL", "MOLINA RUBIO", "09089691E")
	if err != nil {
		panic(err)
	}
}

func deleteQuery(db *sql.DB){
	sqlStatement := "DELETE FROM employee where dni=($1)"

	_, err := db.Exec(sqlStatement,  "09089691E")
	if err != nil {
		panic(err)
	}
}

