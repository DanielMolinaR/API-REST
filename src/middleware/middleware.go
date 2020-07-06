package middleware

import (

	. "TFG/API-REST/src/structures"
	"TFG/API-REST/src/lib"
	"encoding/json"
)

//Conect to the DB
var db = lib.ConectToDB()

func UsersLogin (reqBody []byte) (bool, string){
	var newUser Users
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newUser)
	//verifyLogin verify if all params are correct
	// and if they exist in the DB
	if !verifyDNI(newUser.DNI){
		return false, "DNI incorrecto"
	} else if !checkIfDniExistsAndPassswordIsCorrect(newUser.DNI, newUser.Password){
		return false, "El DNI no existe"
	}
	return true, "Sesion inciciada"
}

func EmployeeSignInVerification(reqBody []byte) (bool, string){
	var newEmployee Employee
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newEmployee)
	//verifyLogin verify if all params are correct
	// and if they exist in the DB
	if !verifyDNI(newEmployee.User.DNI){
		return false, "DNI incorrecto"
	} else if checkIfDniExistsAndPassswordIsCorrect(newEmployee.User.DNI, newEmployee.User.Password){
		return false, "DNI y/o contraseña incorrectos"
	}
	return true, "Usuario creado"
}

func checkIfDniExistsAndPassswordIsCorrect(dni, password string) bool{
	sqlStatement := "SELECT dni, password FROM employee, patients WHERE dni = " + dni
	//Do the query which return a bool and rows of data
	if bool, rows := lib.SelectQuery(db, sqlStatement); !bool{
		return false
	} else {
		//Check if the password is correct
		u := Users{}
		for rows.Next(){
			if rows.Scan(&u.Password); u.Password == password {
				return true
			}
		}
		return false
	}
	db.Close()
}
