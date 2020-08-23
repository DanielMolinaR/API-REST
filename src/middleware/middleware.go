package middleware

import (

	. "TFG/API-REST/src/structures"
	"encoding/json"
)

func UsersLogin (reqBody []byte) (bool, string){
	var userToLog Users
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &userToLog)
	//verify that the DNI or the Email exists
	if len(userToLog.DNI) == 0 && len(userToLog.Email) != 0{
		if !checkIfExists(userToLog.Email, "email"){
			return false, "El usuario no existe"
			//If exists check the password
		} else if bool, response := checkIfPassswordIsCorrect(userToLog.Email, userToLog.Password); !bool{
			return false, response
		}
	} else if len(userToLog.DNI) != 0 && len(userToLog.Email) == 0 {
		if !checkIfExists(userToLog.DNI, "dni"){
			return false, "El usuario no existe"
		}
		//If exists check the password
		if bool, response := checkIfPassswordIsCorrect(userToLog.DNI, userToLog.Password); !bool{
			return false, response
		}
	}
	return true, "Sesion inciciada"
}

func EmployeeSignInVerification(reqBody []byte) (bool, string){
	var newEmployee Employee
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newEmployee)
	bool, response := signInVerifications(newEmployee.User.DNI, newEmployee.User.Phone, newEmployee.User.Email, newEmployee.User.Password)
	if  !bool{
		return false, response
	}

	if !DoEmployeeInsert(newEmployee){
		return false, ""
	}
	return true, "Usuario creado"
}

func PatientSignInVerification(reqBody []byte) (bool, string){
	var newPatient Patient
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newPatient)
	bool, response := signInVerifications(newPatient.User.DNI, newPatient.User.Phone, newPatient.User.Email, newPatient.User.Password)
	if  !bool{
		return false, response
	}
	if !DoPatientInsert(newPatient){
		return false, ""
	}
	return true, "Usuario creado"
}

func signInVerifications(dni, phone, email, password string) (bool, string){
	//verifyDNI verify if the DNI is correct
	// and if it exists in the DB
	if !verifyDNI(dni){
		return false, "DNI incorrecto"
	} else if checkIfExists(dni, "dni"){
		return false, "Este DNI ya existe"
	}
	//Phone number verification
	if !verifyPhoneNumber(phone){
		return false, "El numero de telefono no existe"
	}
	//Email verification
	if !verifyEmail(email){
		return false, "Email no váido"
	}
	//Verify if the password is strong
	if !verifyPasswordIsSafe(string(password)){
		return false, "La contraseña es muy débil"
	}
	return true, ""
}
