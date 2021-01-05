package middleware

import (
	"TFG/API-REST/lib"
	"TFG/API-REST/structures"
	"github.com/badoux/checkmail"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var letters = []string{"t", "r", "w", "a", "g", "m", "y", "f", "p", "d", "x", "b", "n",	"j", "z", "s", "q",	"v", "h", "l", "c", "k", "e"}

func userDataVerifications(condition, dni, phone, email, password string) (bool, map[string]interface{}){

	//verify if the DNI is correct and if it exists in the DB
	if !verifyDNI(dni){
		return false, map[string]interface{}{"state": "DNI incorrecto"}
	} else if checkIfExists(condition,"dni", "dni", dni){
		return false, map[string]interface{}{"state": "Ya existe este DNI"}
	}

	//Verify if the password is strong
	if !verifyPasswordIsSafe(password){
		return false, map[string]interface{}{"state": "Contraseña débil"}
	}

	//verify if the email is correct and if it exists in the DB
	if !verifyEmail(email){
		return false, map[string]interface{}{"state": "Correo no válido"}
	} else if checkIfExists(condition,"email", "email", email){
		return false, map[string]interface{}{"state": "Ya existe este correo"}
	}

	//verify if the phone is correct and if it exists in the DB
	if !verifyPhoneNumber(phone){
		return false, map[string]interface{}{"state": "teléfono no válido"}
	}

	//If everything is correct return true
	return true, map[string]interface{}{"state": ""}

}

func verifyIfDniIsRandom(email string) bool{
	if ok, dni := getStringFromField("patients", "dni", "email", email); !ok{
		return false
	} else {

		//if the dni has a dni format return false
		if verifyDNI(dni){
			return false
		} else {
			return true
		}
	}
}

func existingPatientVerification(condition, dni, phone, password string) (bool, map[string]interface{}) {

	//verify if the new DNI is correct and if it exists in the DB
	if !verifyDNI(dni){
		return false, map[string]interface{}{"state": "DNI incorrecto"}
	} else if checkIfExists(condition,"dni", "dni", dni){
		return false, map[string]interface{}{"state": "Ya existe este DNI"}
	}

	//Verify if the password is strong
	if !verifyPasswordIsSafe(password){
		return false, map[string]interface{}{"state": "Contraseña débil"}
	}

	//verify if the phone is correct and if it exists in the DB
	if !verifyPhoneNumber(phone){
		return false, map[string]interface{}{"state": "teléfono no válido"}
	}

	//If everything is correct return true
	return true, map[string]interface{}{"state": ""}

}

func verifyDNI(dni string) bool{
	//The DNI must has 9 characters
	if len(dni)!=9{
		return false
	} else {
		//The last char of the DNI must be a Letter
		if !verifyLastCharIsALetter(dni){
			return false
		} else {
			//Verify if the Letter is correct with the numbers of the DNI
			if !verifyLetterIsCorrect(dni){
				return false
			}
		}
	}
	return true
}

func verifyLastCharIsALetter(dni string) bool{
	//Take the last char
	c := strings.ToUpper(dni[8:])
	//Verified if the last char is a Letter
	// parsing it to and int and using ASCII
	asciiValue := int(c[0])
	if asciiValue < 97 || asciiValue > 122 {
		return false
	} else {
		return true
	}
}

func verifyLetterIsCorrect (dni string) bool {
	//Parse to int the DNI except the last char
	dniNumber, err := strconv.Atoi(dni[0:8])
	if err!=nil{
		return false
	}
	//The module of the division of the number of the DNI
	// by 23, must be the position of the Letter in dniLetter[]
	//This rule is established by Spain
	if letters[dniNumber%23] != dni[8:]{
		return false
	}
	return true
}

func verifyPhoneNumber(phone string) bool {

	if len(phone)!=9 {
		return false
	} else if !allCharAreNumbers(phone){
		return false
	}
	return true
}

func allCharAreNumbers(phone string) bool{
	//The range of a string return an int32
	//because It iterates over UTF-8-encoded
	//code points in the string
	for i, ch := range phone{
		if int(ch) < 48 || int(ch) > 57{
			return false
		}
		if i == 0{
			//Verify if the first digit of the number
			//matches with one of the three types of
			//phone numbers in Spain (6,7 or 9)
			if !verifyDigit(int(ch)){
				return false
			}
		}
	}
	return true
}

func verifyDigit(c int) bool{
	if c == 54 || c == 55 || c == 57{
		return true
	}
	return false
}

func verifyEmail (email string) bool {
	//Validate Format
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return false
	}
	//Validate Domain
	err = checkmail.ValidateHost(email)
	if err != nil {
		return false
	}
	//Validate User
	err = checkmail.ValidateHost(email)
	if _, ok := err.(checkmail.SmtpError); ok && err != nil {
		return false
	}
	return true
}

func verifyPasswordIsSafe(s string) bool {
	//Validate if the password has at least
	//one letter in upper case, another one
	//in lower case, a special character,
	//a number and if It's longer than 6
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
		hasntSpace = true
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		//
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		case int(char) == 32:
			hasntSpace = false
		}
	}
	//If every value is true the password is safe
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial && hasntSpace
}

func VerifyUuidAndGetExpTime(uuid string) int64 {

	//If the uuid doesnt exists the time is 0
	var unixExpTime int64 = 0
	unixExpTime = getExpTimeFromUuid(uuid)
	return unixExpTime
}

func VerifyExpTime(unixExpTime int64) (bool) {
	expTime := time.Unix(unixExpTime, 0)
	if time.Now().After(expTime){
		lib.TerminalLogger.Error("The uuid has expired. Expiration date: ", expTime)
		lib.DocuLogger.Error("The uuid has expired. Expriration Date: ", expTime)
		return false
	} else {
		lib.TerminalLogger.Error("The slug is correct and it is not expired")
		lib.DocuLogger.Error("The slug is correct and it is not expired")
		return true
	}
}

func verifyAppointmentData(appointmentData structures.Appointment, newUser bool) (bool, map[string]interface{}) {

	if newUser {
		//verify patient email
		if !verifyEmail(appointmentData.Patient_email) {
			return false, map[string]interface{}{"state": "Correo no válido"}
		} else if checkIfExists("patients", "email", "email", appointmentData.Patient_email) {
			return false, map[string]interface{}{"state": "Ya existe este correo"}
		}

		//verify patient phone number
		if !verifyPhoneNumber(appointmentData.Patient_phone){
			return false, map[string]interface{}{"state": "Número de telefono no válido"}
		} else if checkIfExists("patients","phone", "phone", appointmentData.Patient_phone){
			return false, map[string]interface{}{"state": "Ya existe este correo"}
		}

	} else {
		if !checkIfExists("patients", "email", "email", appointmentData.Patient_email) {
			return false, map[string]interface{}{"state": "El correo del paciente no existe"}
		}
	}

	//verify employee email
	if !checkIfExists("employee", "email", "email", appointmentData.Employee_email){
		return false, map[string]interface{}{"state": "El correo del empleado no existe"}
	}

	//verify date
	if !verifyTime(appointmentData.Year, appointmentData.Month, appointmentData.Day, appointmentData.Hour, appointmentData.Minute){
		return false, map[string]interface{}{"state": "Fecha no válida"}
	}

	return true, nil
}

func verifyExerciseData(exerciseData structures.Exercise) (bool, map[string]interface{}){
	if !verifyEmail(exerciseData.Patient_email){
		return false, map[string]interface{}{"state": "Correo no válido"}
	} else if !checkIfExists("patients", "email", "email", exerciseData.Patient_email) {
		return false, map[string]interface{}{"state": "Este correo no existe"}
	}

	//verify date
	if !verifyTime(exerciseData.Year, exerciseData.Month, exerciseData.Day, exerciseData.Hour, exerciseData.Minute){
		return false, map[string]interface{}{"state": "Fecha no válida"}
	}

	return true, nil
}

func verifyTime(year, month, day, hour, minute int) bool{
	date := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)

	if time.Now().After(date){
		return false
	}
	return true
}

func verifyAppointmentAvailableness(patient_dni, employee_dni, date string) (bool, map[string]interface{}) {
	if !verifyEmployeeAvaliableness(employee_dni, date, "appointments"){
		return false, map[string]interface{}{"state": "El empleado ya tiene una cita en esta fecha"}
		lib.TerminalLogger.Error("The employee has a appointments at the same date")
		lib.DocuLogger.Error("The employee has a appointments at the same date")
	} else if !verifyPatientAvaliableness(patient_dni, date, "appointments"){
		return false, map[string]interface{}{"state": "El paciente ya tiene una cita en esta fecha"}
		lib.TerminalLogger.Error("The patient has a appointments at the same date")
		lib.DocuLogger.Error("The patient has a appointments at the same date")
	}
	return true, nil
}

func verifyExerciseAvailableness(patient_dni, date string) (bool, map[string]interface{}) {
	if !verifyPatientAvaliableness(patient_dni, date, "exercises"){
		return false, map[string]interface{}{"state": "El paciente ya tiene una cita en esta fecha"}
		lib.TerminalLogger.Error("The patient has a appointments at the same date")
		lib.DocuLogger.Error("The patient has a appointments at the same date")
	}
	return true, nil
}

func verifyEmployeeAvaliableness(employee_dni, date, table string) bool{
	return checkIfAvailable("employee", employee_dni, date, table)
}

func verifyPatientAvaliableness(patient_dni, date, table string) bool{
	return checkIfAvailable("patients", patient_dni, date, table)
}

