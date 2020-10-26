package middleware

import (
	"TFG/API-REST/lib"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nerzal/gocloak/v7"
	"github.com/dgrijalva/jwt-go/v4"
	"io/ioutil"
	"os"
)

type Data struct{
	ClientId	string	`json:"ClientId"`
	BasePath 	string	`json:"BasePath"`
	AdminName	string	`json:"AdminName"`
	AdminPwd	string	`json:"AdminPwd"`
	AdminRealm	string	`json:"AdminRealm"`
	Secret		string	`json:"Secret"`
	UserRealm	string	`json:"UserRealm"`
}

var data Data
var client gocloak.GoCloak
var ctx context.Context

func init() {

	dataconfig, err := os.Open("./API-REST/middleware/keycloak.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonBody, _ := ioutil.ReadAll(dataconfig)
	json.Unmarshal(jsonBody, &data)

	client = gocloak.NewClient(data.BasePath)

	ctx = context.Background()

}

func getAdminToken() *gocloak.JWT {

	token, err := client.LoginAdmin(ctx, data.AdminName, data.AdminPwd, data.AdminRealm)

	if err != nil {
		lib.TerminalLogger.Info("Problems with the admin login in keycloak", err)
		lib.DocuLogger.Info("Problems with the admin login in keycloak", err)
		return nil
	}

	return token
}

func UserCredentialsLogin(username, password string) (bool, string, string) {

	userToken, err := client.Login(ctx, data.ClientId, data.Secret, data.UserRealm, username, password)

	if err != nil{
		lib.TerminalLogger.Info("Problem with the login in keycloak", err)
		lib.DocuLogger.Info("Problem with the login in keycloak", err)
		return false, "", ""
	}

	return true, userToken.AccessToken, userToken.RefreshToken
}

func CreateUser(username, password, role string) string {

	adminToken := getAdminToken()

	if adminToken != nil{

		user := gocloak.User{
			Email:     gocloak.StringP("something@really.wrong"),
			Enabled:   gocloak.BoolP(true),
			Username:  gocloak.StringP(username),
			Credentials: &[]gocloak.CredentialRepresentation{
				{
					Temporary: gocloak.BoolP(false),
					Type:      gocloak.StringP("password"),
					Value:     gocloak.StringP(password),
				},
			},
			RealmRoles: &[]string{role},
		}
		answer, err := client.CreateUser(ctx, adminToken.AccessToken, data.UserRealm, user)
		if err != nil {
			lib.TerminalLogger.Info("Problem sign in the user", err)
			lib.DocuLogger.Info("Problem sign in the user", err)
			return answer
		}
	} else {
		return "Imposible to login as administrator"
	}
	return ""
}

func DecodeToken(userToken string) *jwt.MapClaims {

	token, claims, err := client.DecodeAccessToken(ctx, userToken, data.UserRealm, "account")

	lib.TerminalLogger.Info(token.Raw)

	if err!=nil{
		lib.TerminalLogger.Info("Problem with the decoding", err)
		lib.DocuLogger.Info("Problem with the decoding", err)
		return nil
	}

	return claims

}

func Verify(token string) bool{
	//Retrospect the token
	rptResult, err := client.RetrospectToken(ctx, token, data.ClientId, data.Secret, data.UserRealm)
	if err != nil {
		lib.TerminalLogger.Info("Problem retrospecting the token", err)
		lib.DocuLogger.Info("Problem retrospecting the token", err)
	}
	//Check if the token is active
	if *rptResult.Active {
		return true
	}

	return false

}