package config

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
)

// Database configuration file path
var dbConfigFileName string = "config/db.json"

// Struct for modelling the configuration json representing
// The database connection details
type DbConfig struct {
    DatabaseName             string `json:"databaseName"`
    DatabaseConnectionString string `json:"databaseConnectionString"`
}

// The database connection string variable
// This variable needs to be initialized
var DbConnectionString string

// The name of the database that will be used
// This variable needs to be initialized
var DbName string

// Reads the configuration file and places all the data
// into the configuration holder struct
func fetchAndDeserializeDbData(filePath string) *DbConfig {
    var configEntity DbConfig

    data, err := ioutil.ReadFile(filePath)

    if err != nil {
        log.Fatal(err)
    }

    err = json.Unmarshal(data, &configEntity)

    if err != nil {
        log.Fatal(err)
    }

    return &configEntity
}

// Initialization of production database
// Production database details is kept within
// special configuration files
func InitDatabase(configFile string) {
    if len(configFile) != 0 {
        dbConfigFileName = configFile
    }

    data := fetchAndDeserializeDbData(dbConfigFileName)

    DbName = data.DatabaseName
    DbConnectionString = data.DatabaseConnectionString
}

// Initialization of tests database
// The tests database configuration details are
// kept in environment variables. This is done
// in order to easily use a CI service
func InitTestsDatabase() {
    dbName := os.Getenv("APIGO_TESTDB_NAME")
    dbConn := os.Getenv("APIGO_TESTDB_CONN")

    if len(dbName) == 0 || len(dbConn) == 0 {
        log.Fatal("Environment variables for the test database are not set!")
    }

    DbName = dbName
    DbConnectionString = dbConn
}
