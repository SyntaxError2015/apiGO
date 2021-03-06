package filter

import (
    "apiGO/dbmodels"
    "apiGO/models"
    "apiGO/service"
)

// Checks if the user entity has all the mandatory
// fields populated with any kind of data
func CheckUserIntegrity(user *models.User) bool {
    switch {
    case len(user.Username) == 0:
        return false
    case len(user.Password) == 0:
        return false
    case len(user.Email) == 0:
        return false
    }

    return true
}

// Checks whether the user already exists or not
func CheckUserExists(user *models.User) bool {
    dbUser, err := service.GetUserByUsernameAndPassword(user.Username, user.Password)

    if err == nil || !dbUser.Equal(dbmodels.User{}) {
        return true
    }

    return false
}
