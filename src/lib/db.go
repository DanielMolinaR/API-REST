package lib

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
	"log"
)

func ConectToDB() *pgxpool.Pool {
	//Set the params to connect to the DB
	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.Host = "localhost"
	config.ConnConfig.Port = 5432
	config.ConnConfig.User = "postgres"
	config.ConnConfig.Password = "postgres"
	config.ConnConfig.Database = "tfg"
	//Create the connection pool
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully connected!")

	return conn
}

func SelectQuery(db *pgxpool.Pool, sqlStatement, dni string) (bool, string) {
	var password string
	//Do the query and if It's correct the password is saved
	err := db.QueryRow(context.Background(), sqlStatement, dni).Scan(&password)
	if err != nil {
		log.Fatal(err)
		return false, ""
	}

	return true, password
}

func InsertQuery(db *pgxpool.Pool, sqlStatement string) (bool){
	_, err := db.Exec(context.Background(), sqlStatement)
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
