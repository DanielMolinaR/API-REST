package middleware

import (
	"TFG/API-REST/lib"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
)

func GetTheRole(token string) int {
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