package middleware

import (
	"TFG/lib"
	"context"
	"github.com/Nerzal/gocloak/v7"
	"github.com/dgrijalva/jwt-go/v4"
)

type Data struct{
	BasePath 	string	`json:"BasePath"`
	AdminName	string	`json:"Username"`
	AdminPwd	string	`json:"Password"`
	AdminRealm	string	`json:"AdminRealm"`
	Secret		string	`json:"Secret"`
	UserRealm	string	`json:"UserRealm"`
}

func getClientAndContext() (gocloak.GoCloak, context.Context) {
	return gocloak.NewClient("http://localhost:8080"), context.Background()
}

func getAdminToken() *gocloak.JWT {
	client, ctx := getClientAndContext()
	token, err := client.LoginAdmin(ctx, "DaniM", "admin", "master")
	if err != nil {
		lib.TerminalLogger.Info("Problems with the admin login in keycloak", err)
		lib.DocuLogger.Info("Problems with the admin login in keycloak", err)
		return nil
	}

	return token
}

func UserCredentialsLogin(username, password string) (bool, string, string) {

	client, ctx := getClientAndContext()

	userToken, err := client.Login(ctx, "front-end", "0acd54f5-4232-4a21-b6e5-fed4a5527835", "Physiotherapy", username, password)

	if err != nil{
		lib.TerminalLogger.Info("Problem with the login in keycloak", err)
		lib.DocuLogger.Info("Problem with the login in keycloak", err)
		return false, "", ""
	}

	return true, string(userToken.AccessToken), string(userToken.RefreshToken)
}

func CreateUser(username, password, role string) ( string) {

	client, ctx := getClientAndContext()

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
		answer, err := client.CreateUser(ctx, adminToken.AccessToken, "Physiotherapy", user)
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

	client, ctx := getClientAndContext()

	_, claims, err := client.DecodeAccessToken(ctx, userToken, "Physiotherapy", "account")

	if err!=nil{
		lib.TerminalLogger.Info("Problem with the decoding", err)
		lib.DocuLogger.Info("Problem with the decoding", err)
		return nil
	}

	return claims

}