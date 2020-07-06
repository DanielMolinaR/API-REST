package lib

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "tfg"
)

type Employee struct{
	dni string
	email string
	password string
	name string
	surname string
}

func ConectToDB() *sql.DB{

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db

}

func SelectQuery(db *sql.DB){
	sqlStatement := "select * from employee;"

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	employee := []Employee{}

	for rows.Next() {
		b := Employee{}
		if err := rows.Scan(&b.dni, &b.email, &b.password, &b.name, &b.surname); err != nil {
			log.Fatal(err)
		}
		employee = append(employee, b)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(employee)
}

func SelectQuery2(db *sql.DB, sqlStatement string) (bool, *sql.Rows) {

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return false, nil
	}
	return true, rows
}

func insertQuery(db *sql.DB){
	sqlStatement := "Insert into employee (dni, email, password, name, surname) values ($1, $2, $3, $4, $5)"

	_, err := db.Exec(sqlStatement, "09089691E", "dani99mr31@gmail.com", "12345678", "Dani", "Molina")
	if err != nil {
		panic(err)
	}
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
