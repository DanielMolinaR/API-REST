package middleware

import (

	. "TFG/API-REST/src/structures"
	"encoding/json"
)

func UsersLogin (reqBody []byte) (bool, string){
	var newUser Users
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newUser)
	//verifyDNI verify if the DNI is correct
	// and if it exists in the DB
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
	//verifyDNI verify if the DNI is correct
	// and if it exists in the DB
	if !verifyDNI(newEmployee.User.DNI){
		return false, "DNI incorrecto"
	} else if checkIfDniExistsAndPassswordIsCorrect(newEmployee.User.DNI, newEmployee.User.Password){
		return false, "Este DNI ya existe"
	}
	//Phone number verification
	if !verifyPhoneNumber(newEmployee.User.Phone){
		return false, "El numero de telefono no existe"
	}
	//Email verification
	if !verifyEmail(newEmployee.User.Email){
		return false, "Email no váido"
	}
	if !verifyPasswordIsSafe(newEmployee.User.Password){
		return false, "La contraseña es muy débil"
	}
	sqlStatement := "INSERT INTO employee (active, admin, dni, email, password, Name, Surname, phone) " +
		"VALUES (" + newEmployee.Active + ", " +  newEmployee.Admin + ", " + newEmployee.User.DNI + ", " +
		newEmployee.User.Email + ", " + newEmployee.User.Password + ", " + newEmployee.User.Name + ", " +
		newEmployee.User.Surname + ", " + newEmployee.User.Phone + ")"
	if !DoInsert(sqlStatement){
		return false, "No se ha podido crear el usuario"
	}
	return true, "Usuario creado"
}

