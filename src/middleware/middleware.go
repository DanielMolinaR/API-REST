package middleware

import (

	. "TFG/API-REST/src/structures"
	"encoding/json"
)

func UsersLogin(reqBody []byte) (bool, map[string]interface{}) {
	var userToLog Users
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &userToLog)
	//verify that the DNI or the Email exists
	if len(userToLog.DNI) == 0 && len(userToLog.Email) != 0{
		if !checkIfExists(userToLog.Email, "email"){
			return false, map[string]interface{}{"response": "El usuario no existe"}
			//If exists check the password
		} else if bool, response := checkIfPassswordIsCorrect(userToLog.Email, userToLog.Password); !bool{
			return false, map[string]interface{}{"response": response}
		}
	} else if len(userToLog.DNI) != 0 && len(userToLog.Email) == 0 {
		if !checkIfExists(userToLog.DNI, "dni"){
			return false, map[string]interface{}{"response": "El usuario no existe"}
		}
		//If exists check the password
		if bool, response := checkIfPassswordIsCorrect(userToLog.DNI, userToLog.Password); !bool{
			return false, map[string]interface{}{"response": response}
		}
	}
	if len(userToLog.DNI) == 0 && len(userToLog.Email) != 0 {
		return true, map[string]interface{}{"response": "Sesión inicada", "name": getUserName(userToLog.Email, "email"),
			"userId": getUserId(userToLog.Email, "email"), "token": generateToken()}
	} else {
		return true, map[string]interface{}{"response": "Sesión inicada", "name": getUserName(userToLog.DNI, "dni"),
			"userId": getUserId(userToLog.DNI, "dni"), "token": generateToken()}
	}
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
