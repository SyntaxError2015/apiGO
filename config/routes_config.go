package config

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "log"
    "os"
)

// Routes configuration file path
var routesConfigFile = "config/routes.json"

// Variable for storing all the routes that the api will have
var Routes []Route

// Initialize the routes by deserializing them from the configuration file
func InitRoutes(routesConfigPath string) {
    if len(routesConfigPath) != 0 {
        routesConfigFile = routesConfigPath
    }

    data, err := ioutil.ReadFile(routesConfigFile)

    if err != nil {
        log.Fatal(err)
    }

    err = json.Unmarshal(data, &Routes)

    if err != nil {
        log.Fatal(err)
    }
}

// Save all the active routes (Routes slice) in json format
// into the configuration file
func SaveRoutesConfiguration() error {
    if len(Routes) == 0 {
        return errors.New("There are no routes configured in order to be saved")
    }

    data, err := json.MarshalIndent(Routes, "", "  ")

    if err != nil {
        return errors.New("Encoding routes slice to json failed!")
    }

    err = ioutil.WriteFile(routesConfigFile, data, os.ModeDevice)
    if err != nil {
        return err
    }

    return nil
}

// Add a new route and make it active
func AddRoute(route *Route, saveChangesToConfigFile bool) error {
    initialLength := len(Routes)

    for _, r := range Routes {
        if r.Id == route.Id {
            return errors.New("Route already exists!")
        }
    }

    Routes = append(Routes, *route)

    err := checkCollectionModification(route, initialLength)

    return saveChanges(err, saveChangesToConfigFile, SaveRoutesConfiguration)
}

// Disable and remove a certain route
func RemoveRoute(routeId string, saveChangesToConfigFile bool) error {
    initialLength := len(Routes)
    index := -1

    for ind, route := range Routes {
        if route.Id == routeId {
            index = ind
            break
        }
    }

    if index == -1 {
        return errors.New("Route was not found for deletion!")
    }

    removedRoute := Routes[index]
    Routes = append(Routes[:index], Routes[index+1:]...)

    err := checkCollectionModification(&removedRoute, initialLength)

    return saveChanges(err, saveChangesToConfigFile, SaveRoutesConfiguration)
}

// Modify the state and information of a certain route
func ModifyRoute(routeId string, newRouteData Route, saveChangesToConfigFile bool) error {
    for i := 0; i < len(Routes); i++ {
        if Routes[i].Id == routeId {
            Routes[i] = newRouteData

            return saveChanges(nil, saveChangesToConfigFile, SaveRoutesConfiguration)
        }
    }

    return errors.New("Route was not found for modification!")
}

// Get a Route entity from the active routes list, base on its ID
func GetRoute(routeId string) *Route {
    for _, route := range Routes {
        if route.Id == routeId {
            return &route
        }
    }

    return nil
}

// Get a Route entity from the active routes list, base on its pattern
func GetRouteByPattern(pattern string) *Route {
    for _, route := range Routes {
        if route.Pattern == pattern {
            return &route
        }
    }

    return nil
}

// This function checks whether there has been a modification in the
// array of active API routes. If there is no modification, then
// this means that an error has occured when adding, deleting or
// modifying an existing route
func checkCollectionModification(route *Route, initialLength int) error {
    if initialLength == len(Routes) {
        return errors.New("The route couldn't be processed for the collection!")
    }

    return nil
}

// Check whether the new routes configuration can be saved to the config file.
// If yes, then the function of choice is used in order to serialize all the routes.
func saveChanges(err error, saveChangesToConfigFile bool, saverFunction func() error) error {
    if err != nil {
        return err
    }

    if saveChangesToConfigFile {
        return saverFunction()
    }

    return nil
}
