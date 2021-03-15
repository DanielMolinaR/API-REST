package middleware

import (
	"TFG/lib"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nerzal/gocloak/v7"
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

	dataconfig, err := os.Open("./middleware/keycloak.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonBody, _ := ioutil.ReadAll(dataconfig)
	json.Unmarshal(jsonBody, &data)

	client = gocloak.NewClient(data.BasePath)

	ctx = context.Background()

}

func getAdminToken() string {

	token, err := client.LoginAdmin(ctx, data.AdminName, data.AdminPwd, data.AdminRealm)

	if err != nil {
		lib.TerminalLogger.Info("Problems with the admin login in keycloak", err)
		lib.DocuLogger.Info("Problems with the admin login in keycloak", err)
		return ""
	}

	return token.AccessToken
}

func UserCredentialsLogin(userDni, password string) (bool, string, string) {

	userToken, err := client.Login(ctx, data.ClientId, data.Secret, data.UserRealm, userDni, password)

	if err != nil{
		lib.TerminalLogger.Error("Problem with the login in keycloak", err)
		lib.DocuLogger.Error("Problem with the login in keycloak", err)
		return false, "", ""
	}

	return true, userToken.AccessToken, userToken.RefreshToken
}

func createKeycloakUser(userDni, password, email, role string) (bool, string) {

	user := gocloak.User{
		Enabled:   gocloak.BoolP(false),
		Username:  gocloak.StringP(userDni),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP(password),
			},
		},
		Email: gocloak.StringP(email),
	}

	userId, err := client.CreateUser(ctx, getAdminToken(), data.UserRealm, user)
	if err != nil {
		lib.TerminalLogger.Info("Something went wrong", err)
		lib.DocuLogger.Info("Something went wrong", err)
		return false, ""
	}
	if !updateUserRole(role, userId){
		lib.TerminalLogger.Info("Impossible to update the role", err)
		lib.DocuLogger.Info("Impossible to update the role", err)
		DeleteKeycloakUser(userId)
		return false, ""
	}
	lib.TerminalLogger.Info("User created, answer: ")
	lib.DocuLogger.Info("User created, answer: ")
	return true, userId
}

func updateUserRole(role, userId string) bool{
	//As we cant create an user with a role for problems that keycloak has
	//we have to update the user. We need to take the ID from the user.
	roleObject, err := client.GetRealmRole(ctx, getAdminToken(), data.UserRealm, role)

	var roles []gocloak.Role
	roles = append(roles, *roleObject)

	err = client.AddRealmRoleToUser(ctx, getAdminToken(), data.UserRealm, userId, roles)
	if err != nil{
		lib.TerminalLogger.Info("Something went wrong updating the user role: ", err)
		lib.DocuLogger.Info("Something went wrong updating the user role: ", err)
		return false
	}
	return true
}

func updateUserEnabled(userId string) bool {
	user := gocloak.User{
		ID: gocloak.StringP(userId),
		Enabled:   gocloak.BoolP(true),
	}
	err := client.UpdateUser(ctx, getAdminToken(), data.UserRealm, user)
	if err != nil{
		lib.TerminalLogger.Info("Something went wrong updating the user enabled: ", err)
		lib.DocuLogger.Info("Something went wrong updating the user enabled: ", err)
		return false
	} else {
		return true
	}
}

func getEmployeesFromKeycloak() ([]*gocloak.User, error) {
	return client.GetUsersByRoleName(ctx, getAdminToken(), data.UserRealm, "EMPLOYEE_ROLE")
}

func updateEmployeeToAdmin(userId string) (bool, map[string]interface{}){
	//As we cant create an user with a role for problems that keycloak has
	//we have to update the user. We need to take the ID from the user.
	employeeRole, err := client.GetRealmRole(ctx, getAdminToken(), data.UserRealm, "EMPLOYEE_ROLE")
	if err != nil{
		lib.TerminalLogger.Info("Couldn't get the employee role", err)
		lib.DocuLogger.Info("Couldn't get the employee role", err)
		return false, map[string]interface{}{"state": "No se ha podido dar permisos de administrador a este empleado"}
	}

	var eRole []gocloak.Role
	eRole = append(eRole, *employeeRole)
	client.DeleteRealmRoleFromUser(ctx, getAdminToken(), data.UserRealm, userId, eRole)

	adminRole, err := client.GetRealmRole(ctx, getAdminToken(), data.UserRealm, "ADMIN_ROLE")
	if err != nil{
		lib.TerminalLogger.Info("Couldn't get the admin role", err)
		lib.DocuLogger.Info("Couldn't get the admin role", err)
		return false, map[string]interface{}{"state": "No se ha podido dar permisos de administrador a este empleado"}
	}

	var aRole []gocloak.Role
	aRole = append(aRole, *adminRole)
	err = client.AddRealmRoleToUser(ctx, getAdminToken(), data.UserRealm, userId, aRole)
	if err != nil{
		lib.TerminalLogger.Info("Couldn't update the user. ", err)
		lib.DocuLogger.Info("Couldn't update the user. ", err)
		return false, map[string]interface{}{"state": "No se ha podido dar permisos de administrador a este empleado"}
	} else {
		lib.TerminalLogger.Info("The employee has been upgraded. ")
		lib.DocuLogger.Info("The employee has been upgraded. ")
		return true, map[string]interface{}{"state": "Este empleado ya tiene permisos de administrador"}
	}

}

func DeleteKeycloakUser(userId string) bool{

	//we never are going to delete someone from keycloak or the database
	//we need them and their data so what we do is making disabled the keycloak user

	err := client.DeleteUser(ctx, getAdminToken(), data.UserRealm, userId)
	if err != nil{
		return false
	}
	return true
}

