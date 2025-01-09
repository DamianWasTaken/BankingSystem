package main

import (
	"LoggingService/utils"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()

	loggingEnv := setupEnviorment()

	logRoutes := r.Group("/logging")

	accountRoutes := logRoutes.Group("/account")
	accountRoutes.POST("/status", loggingEnv.PersistStatusChange)
	accountRoutes.GET("/status", loggingEnv.GetAccountStatusHistory)

	logRoutes.POST("/balance", loggingEnv.PersistBalanceChange)
	logRoutes.POST("/interest", loggingEnv.PersistInterestChange)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Failed to start server")
	}
}

// load dependencies and setup enviorment
func setupEnviorment() *LoggingEnviormentService {
	// this should be a db pool ideally, and repositories would be tapping into the connections of the pool
	database, err := GetDatabase()
	if err != nil {
		return nil
	}
	// repositories initialization
	return &LoggingEnviormentService{
		LogAccountManagement:  &utils.LogAccountManagement{DB: database},
		LogBalanceManagement:  &utils.LogBalanceManagement{DB: database},
		LogInterestManagement: &utils.LogInterestManagement{DB: database},
	}

}

func GetDatabase() (*sql.DB, error) {
	// change these to env variables
	port := 5432
	url := "logging-db"
	password := "admin"
	db := "logging"
	user := "admin"

	postgresConnection, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", url, port, user, password, db))
	if err != nil {
		fmt.Println("Failed to connect to database " + err.Error())
	}
	return postgresConnection, nil
}
