package middleware

import (
	"TFG/API-REST/lib"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"strings"
)

func getTheRole(token string) int {
	claims_token := decodeToken(token)
	allTheRoles := (*claims_token)["realm_access"].(map[string]interface{})["roles"]
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

func getUserDniFromToken(token string) (string){
	claims_token := decodeToken(token)
	data := (*claims_token)["sub"]
	dni := data.([]interface{})
	return dni[0]
}

func getEmail(token string) string {
	claims_token := decodeToken(token)
	return fmt.Sprintf("%v", (*claims_token)["email"])
}

func decodeToken(userToken string) *jwt.MapClaims {

	_, claims, err := client.DecodeAccessToken(ctx, userToken, data.UserRealm, "account")

	if err!=nil{
		lib.TerminalLogger.Info("Problem with the decoding", err)
		lib.DocuLogger.Info("Problem with the decoding", err)
		return nil
	}

	return claims

}

func VerifyToken(token string) bool{
	//Retrospect the token
	rptResult, err := client.RetrospectToken(ctx, token, data.ClientId, data.Secret, data.UserRealm)
	if err != nil {
		lib.TerminalLogger.Error("Problem retrospecting the token", err)
		lib.DocuLogger.Error("Problem retrospecting the token", err)
	}
	//Check if the token is active
	if rptResult != nil {
		if *rptResult.Active {
			return true
		}
	}
	return false

}

func VerifyTokenIsFromEmployeeOrAdmin(tokenBrearer string) (bool, map[string]interface{}) {

	//Extract the Bearer from the data of the header
	token := strings.Replace(tokenBrearer, "Bearer ", "", -1)

	if !VerifyToken(token){
		return false, map[string]interface{}{"state": "Token no válido"}

		//Verify if the user that is requesting this endpoint is an employee or an admin
	} else if getTheRole(token)< 1 {
		lib.TerminalLogger.Warn("Someone who is not an Amdin or employee is trying to do something restricted")
		lib.DocuLogger.Warn("Someone who is not an Amdin or employee is triying to do something restricted")
		return false, map[string]interface{}{"state": "Acceso restringido"}
	}
	return true, nil
}

func VerifyTokenIsFromAdmin(tokenBrearer string) (bool, map[string]interface{}) {

	//Extract the Bearer from the data of the header
	token := strings.Replace(tokenBrearer, "Bearer ", "", -1)

	if !VerifyToken(token){
		return false, map[string]interface{}{"state": "Token no válido"}

		//Verify if the user that is requesting this endpoint is an employee or an admin
	} else if getTheRole(token) != 2 {

		return false, map[string]interface{}{"state": "Acceso restringido"}
	}
	return true, nil
}