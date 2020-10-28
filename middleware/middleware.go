package middleware

import (
	"TFG/API-REST/lib"
	. "TFG/API-REST/structures"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
)


func UsersLogin(reqBody []byte) (bool, map[string]interface{}) {

	var userToLogIn Users

	//The data from reqBody is filled in the newUser
	_ = json.Unmarshal(reqBody, &userToLogIn)

	//Verify the credentials login
 	if ok, accessToken, refreshToken := UserCredentialsLogin(userToLogIn.DNI, userToLogIn.Password); !ok{
		return false, map[string]interface{}{"state": "DNI o contraseña incorrecto"}
	} else {
		//Return true with a msg of correct login,
		//the name of the user, the tokens and the role
		roles := DecodeToken(accessToken)
		lib.TerminalLogger.Info("User logged with the DNI: ******", userToLogIn.DNI[6:])
		lib.DocuLogger.Info("User logged with the DNI: ******", userToLogIn.DNI[6:])
		return true, map[string]interface{}{"state": "Sesión iniciada", "Access token": accessToken,
			"Refresh token": refreshToken, "Roles": GetTheRole(roles)}
	}
}

func GetTheRole(roles *jwt.MapClaims) int {
	allTheRoles := (*roles)["realm_access"].(map[string]interface{})["roles"]
	data := allTheRoles.([]interface{})
	answer := 0
	for i := 0; i < len(data); i++{
		userRole := fmt.Sprintf("%v", data[i])
		if userRole == "EMPLOYEE_ROLE"{
			answer = 1
		}
		if userRole == "ADMIN_ROLE"{
			answer = 2
		}
	}
	return answer
}

func EmployeeSignInVerification(reqBody []byte) (bool, string){
	var newEmployee Employee
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newEmployee)
	newEmployee.Admin = false
	newEmployee.Active = true
	bool, response := signInVerifications("employee", newEmployee.User.DNI, newEmployee.User.Phone, newEmployee.User.Email, newEmployee.User.Password)
	if  !bool{
		return false, response
	}else if !DoEmployeeInsert(newEmployee){
		return false, ""
	}
	return true, "Usuario creado"
}

func PatientSignInVerification(reqBody []byte) (bool, string){
	var newPatient Patient
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newPatient)
	bool, response := signInVerifications("patients", newPatient.User.DNI, newPatient.User.Phone, newPatient.User.Email, newPatient.User.Password)
	if  !bool{
		return false, response
	}
	if !DoPatientInsert(newPatient){
		return false, ""
	}
	return true, "Usuario creado"
}

func signInVerifications(condition, dni, phone, email, password string) (bool, string){
	//verifyDNI verify if the DNI is correct
	// and if it exists in the DB
	if !verifyDNI(dni){
		return false, "DNI incorrecto"
	} else if checkIfExists(condition, dni){
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
	if !verifyPasswordIsSafe(password){
		return false, "La contraseña es muy débil"
	}
	return true, ""
}

