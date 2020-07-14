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
	host string `json:host`
	port uint16 `json:port`
	user string `json:user`
	password string `json:password`
	database string `json:database`
}

func ConectToDB() *pgxpool.Pool {
	dataconfig, err := os.Open("conection.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonBody, _ := ioutil.ReadAll(dataconfig)
	var conection Conection
	json.Unmarshal(jsonBody, &conection)
	//Set the params to connect to the DB
	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.Host = conection.host
	config.ConnConfig.Port = conection.port
	config.ConnConfig.User = conection.user
	config.ConnConfig.Password = conection.password
	config.ConnConfig.Database = conection.database
	//Create the connection pool
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully connected!")

	return conn
}

func SelectQuery(db *pgxpool.Pool, sqlStatement, data string) (bool, string) {
	var password string
	//Do the query and if It's correct the password is saved
	err := db.QueryRow(context.Background(), sqlStatement, data).Scan(&password)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	return true, password
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
