package middleware

import (
	"TFG/lib"
	. "TFG/structures"
	"encoding/json"
)

/*func ValidateToken(token string) {
	claims := DecodeToken(token)


}*/

func UsersLogin(reqBody []byte) (bool, map[string]interface{}) {

	var userToLogIn Users

	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &userToLogIn)

	//Verify the credentials login
 	if ok, accessToken, refreshToken := UserCredentialsLogin(userToLogIn.DNI, userToLogIn.Password); !ok{
		return false, map[string]interface{}{"state": "DNI o contraseña incorrecto"}
	} else {
		//Return true with a msg of correct login,
		//the name of the user and the role
		lib.TerminalLogger.Info("User logged with the DNI: ******", userToLogIn.DNI[6:])
		lib.DocuLogger.Info("User logged with the DNI: ******", userToLogIn.DNI[6:])
		return true, map[string]interface{}{"state": "Sesión iniciada", "Access token": accessToken, "Refresh token": refreshToken}
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
	} else if checkIfExists(dni){
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

/*func getUserRole (token string) string{
	claims := DecodeToken(token)
	role := claims
}*/
