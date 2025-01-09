package main

import (
	"InterestService/utils"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func main() {
	r := gin.Default()

	interestEnv := setupEnviorment()

	c := cron.New()

	c.AddFunc("0 0 * * *", func() {
		interestEnv.ProcessInterest()
	})

	interestRates := r.Group("/interestRates")
	interestRates.PATCH("/modify", interestEnv.ModifyDailyInterestRate)

	interstUser := r.Group("/interestUser")
	interstUser.POST("/add", interestEnv.AddInterestUser)
	interstUser.PATCH("/modifyFrequency", interestEnv.ModifyInterestUserFrequency)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Failed to start server")
	}
}

// load dependencies and setup enviorment
func setupEnviorment() *InterestEnviormentService {
	// this should be a db pool ideally, and repositories would be tapping into the connections of the pool
	database, err := GetDatabase()
	if err != nil {
		return nil
	}
	// repositories initialization
	return &InterestEnviormentService{
		InterestManagement:     &utils.InterestManagement{DB: database},
		InterestUserManagement: &utils.InterestUserManagement{DB: database},
	}

}

func GetDatabase() (*sql.DB, error) {
	// change these to env variables
	port := 5432
	url := "interest-db"
	password := "admin"
	db := "interest"
	user := "admin"

	postgresConnection, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", url, port, user, password, db))
	if err != nil {
		fmt.Println("Failed to connect to database " + err.Error())
	}
	return postgresConnection, nil
}
