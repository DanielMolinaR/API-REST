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

func SelectQueryPwd(sqlStatement, data string) (string) {
	var password string
	//Do the query and if It's correct
	//It means that the password is saved
	rows, err := db.Query(context.Background(), sqlStatement, data)
	rows.Scan(&password)
	if err != nil {
		TerminalLogger.Error("Error with the query: ", err)
		DocuLogger.Error("Error with the query: ", err)
		return ""
	}

	return password
}

func SelectQuery(sqlStatement, data string) (bool) {
	var dni string

	rows := db.QueryRow(context.Background(), sqlStatement, data).Scan(&dni)
	if rows != nil {
		if dni == ""{
			TerminalLogger.Info("The DNI doesnt exists in the DDBB")
			DocuLogger.Info("The DNI doesnt exists in the DDBB")
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

func SelectUserDataQuery(sqlStatement, data string) map[string]interface{} {
	var (
		name string
		email string
		phone string
		surname string
		dni string
		password string
	)
	//Do the query and if It's correct
	//It means that the password is saved
	rows, err := db.Query(context.Background(), sqlStatement, data)
	rows.Scan(&dni, &email, &password, &name, &surname, &phone)
	if err != nil {
		TerminalLogger.Error("Error with the query:", err)
		DocuLogger.Error("Error with the query:", err)
		return map[string]interface{}{"error": err}
	}

	return map[string]interface{}{"name": name, "email": email, "phone": phone, "surname": surname}
}

func SelectEmployeeDataQuery(sqlStatement, data string) (bool, bool) {
	var (
		admin bool
		active bool
	)
	//Do the query and if It's correct
	//It means that the password is saved
	rows, err := db.Query(context.Background(), sqlStatement, data)
	rows.Scan(&admin, &active)
	if err != nil {
		TerminalLogger.Error("Error with the query:", err)
		DocuLogger.Error("Error with the query:", err)
		return false, false
	}
	return active, admin
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
	_, err := db.Exec(context.Background(), sqlStatement, patient.Age, patient.User.DNI, patient.User.Email,
		patient.User.Name, patient.User.Surname, patient.User.Phone)
	if err != nil {
		fmt.Println(err)
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

