package middleware

import "strings"

func dniIncorrect(newUser users) bool{
	if len(newUser.DNI)!=9{
		return false
	} else {
		c := strings.ToUpper(newUser.DNI[8:])
		asciiValue := int(c[0])
		if asciiValue < 65 || asciiValue > 90 {
			return false
		} else {
			//dniNumbers, _ := strconv.Atoi(newUser.DNI[0:8])


			//verificar que el dni existe, numeros / 23
		}
	}
	//if dni no existe return false
	return true
}