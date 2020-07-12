package middleware

import (

	. "TFG/API-REST/src/structures"
	"encoding/json"
)

func UsersLogin (reqBody []byte) (bool, string){
	var userToLog Users
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &userToLog)
	// and if it exists in the DB
	// ************ Como hacer la verificacion del login con el correo
	// ************ Cuando rellene el campo en l front, al ser comun
	// ************ el cuadro de relleno como diferencia si lo que le
	// ************ paso es para . DNI o . Email
	if len(userToLog.DNI) == 0 && len(userToLog.Email) != 0{
		if bool, response := checkIfEmailExistsAndPassswordIsCorrect(userToLog.Email, userToLog.Password); !bool{
			return false, response
		}
	} else if len(userToLog.DNI) != 0 && len(userToLog.Email) == 0 {
		if bool, response := checkIfDniExistsAndPassswordIsCorrect(userToLog.DNI, userToLog.Password); !bool{
			return false, response
		}
	}
	return true, "Sesion inciciada"
}

/*func EmployeeSignInVerification(reqBody []byte) (bool, string){
	var newEmployee Employee
	//The data from reqBody is filled in the newUser
	json.Unmarshal(reqBody, &newEmployee)
	bool, response := signInVerifications(newEmployee.User.DNI, newEmployee.User.Phone, newEmployee.User.Email, newEmployee.User.Password)
	if  !bool{
		return false, response
	}
	sqlStatement := "INSERT INTO employee (active, admin, dni, email, password, Name, Surname, phone) " +
		"VALUES (" + newEmployee.Active + ", " +  newEmployee.Admin + ", '" + newEmployee.User.DNI + "', '" +
		newEmployee.User.Email + "', '" + newEmployee.User.Password + "', '" + newEmployee.User.Name + "', '" +
		newEmployee.User.Surname + "', '" + newEmployee.User.Phone + "')"
	if !DoInsert(sqlStatement){
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
	sqlStatement := "INSERT INTO patients ( age, dni, email, password, Name, Surname, phone) " +
		"VALUES (" + newPatient.Age + ", '" + newPatient.User.DNI + "', '" + newPatient.User.Email +
		"', '" + newPatient.User.Password + "', '" + newPatient.User.Name + "', '" +
		newPatient.User.Surname + "', '" + newPatient.User.Phone + "')"
	if !DoInsert(sqlStatement){
		return false, "No se ha podido crear el usuario"
	}
	return true, "Usuario creado"
}

func signInVerifications(dni, phone, email, password string) (bool, string){
	//verifyDNI verify if the DNI is correct
	// and if it exists in the DB
	if !verifyDNI(dni){
		return false, "DNI incorrecto"
	} else if checkIfDniExists(dni){
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
}*/
