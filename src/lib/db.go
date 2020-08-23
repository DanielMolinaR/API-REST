package lib

import (
	"TFG/API-REST/src/structures"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
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

func ConectToDB() *pgxpool.Pool {
	dataconfig, err := os.Open("./API-REST/src/lib/conection.json")
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
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully connected!")

	return conn
}

func SelectQueryPwd(db *pgxpool.Pool, sqlStatement string, data string) (bool, string) {
	var password string
	//Do the query and if It's correct
	//It means that the password is saved
	err := db.QueryRow(context.Background(), sqlStatement, data).Scan(&password)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	return true, password
}

func SelectQuery(db *pgxpool.Pool, sqlStatement, data string) (bool) {
	_, err := db.Query(context.Background(), sqlStatement, data)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func InsertEmployeeQuery(db *pgxpool.Pool, sqlStatement string, employee structures.Employee) (bool){
	_, err := db.Exec(context.Background(), sqlStatement, employee.Active, employee.Admin, employee.User.DNI,
		employee.User.Email, employee.User.Password, employee.User.Name, employee.User.Surname, employee.User.Phone)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func InsertPatientQuery(db *pgxpool.Pool, sqlStatement string, patient structures.Patient) (bool){
	_, err := db.Exec(context.Background(), sqlStatement, patient.Age, patient.User.DNI, patient.User.Email,
		patient.User.Password, patient.User.Name, patient.User.Surname, patient.User.Phone)
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
