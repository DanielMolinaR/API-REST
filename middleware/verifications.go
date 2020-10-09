package middleware

import (
	"github.com/badoux/checkmail"
	"strconv"
	"strings"
	"unicode"
)

var letters = []string{"T", "R", "W", "A", "G", "M", "Y", "F", "P", "D", "X", "B", "N",	"J", "Z", "S", "Q",	"V", "H", "L", "C", "K", "E"}

func verifyDNI(dni string) bool{
	//The DNI must has 9 characters
	if len(dni)!=9{
		return false
	} else {
		//The last char of the DNI must be a Letter
		if !verifyLastCharIsALetter(dni){
			return false
		} else {
			//Verify if the Letter is correct
			//with the numbers of the DNI
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
	if asciiValue < 65 || asciiValue > 90 {
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
	//It checks if the phone number has 9 characters
	if len(phone)!=9 {
		return false
		//It checks if all the characters are numbers
	} else if !allAreNumbers(phone){
		return false
	}
	return true
}

func allAreNumbers(phone string) bool{
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
