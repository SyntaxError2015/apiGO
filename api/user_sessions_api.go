package api

import (
    "apiGO/filter"
    "apiGO/models"
    "apiGO/service"
    "net/http"
)

func (api *Api) GetUserSession(vars *ApiVar, resp *ApiResponse) error {
    token, err, found := filter.GetStringValueFromParams("token", vars.RequestForm)

    if !found {
        return badRequest(resp, "Session token has not been specified")
    }

    if err != nil {
        return badRequest(resp, err.Error())
    }

    userSession, err := service.GetUserSessionByToken(token)
    if err != nil || userSession == nil {
        return notFound(resp, "There is no session with the specified token")
    }

    user, err := service.GetUser(userSession.UserId)
    if err != nil || user == nil {
        service.DeleteUserSession(userSession.Id)
        return notFound(resp, "The user with the current session no longer exists")
    }

    expandedUser := &models.User{}
    expandedUser.Expand(*user)
    userJson, _ := expandedUser.SerializeJson()

    resp.StatusCode = http.StatusOK
    resp.Message = userJson

    return nil
}

// PRIMESC USER + PASS intr-un JSON de User
// PARAMS: 1.) username 2.) password
// Validare key value, daca nu -> 401 UNAUTHORIZED
// STERGE SESIUNILE DEJA EXISTENTE LA USER
// CREEZ O SESIUNE CU USERUL MATCHUIT
// RETURN Token de la user
func (api *Api) PostUserSession(vars *ApiVar, resp *ApiResponse) error {
    // Get URL parameters
    username, userError, userWasFound := filter.GetStringValueFromParams("username", vars.RequestForm)
    password, passwordError, passwordWasFound := filter.GetStringValueFromParams("password", vars.RequestForm)

    if !userWasFound || !passwordWasFound {
        return badRequest(resp, "The username or password was not specified")
    }

    if userError != nil {
        return badRequest(resp, userError.Error())
    }
    if passwordError != nil {
        return badRequest(resp, passwordError.Error())
    }

    // Fetch User entity from database
    user, err := service.GetUserByUsernameAndPassword(username, password)
    if err != nil || user == nil {
        return unauthorized(resp, "Username or password is incorrect")
    }

    // Delete all the existing sessions
    service.DeleteAllSessionsWithUserId(user.Id)

    // Generate a new user session
    userSession, err := service.GenerateAndInsertUserSession(user.Id)
    if err != nil {
        return internalServerError(resp, err.Error())
    }

    resp.StatusCode = http.StatusCreated
    resp.Message = []byte(userSession.Token)

    return nil
}

func (api *Api) DeleteUserSession(vars *ApiVar, resp *ApiResponse) error {
    return nil
}
