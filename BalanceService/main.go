package main

import (
	"BalanceService/utils"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()

	balanceEnv := setupEnviorment()

	currencyAccountRoutes := r.Group("/currencyAccount")
	currencyAccountRoutes.Use(balanceEnv.Validate)
	{
		currencyAccountRoutes.POST("/create", balanceEnv.CreateCurrencyAccount)
		currencyAccountRoutes.POST("/IntraAccountTransaction", balanceEnv.ProcessIntraAccountTransaction)

	}

	transactionRoutes := r.Group("/balance")
	transactionRoutes.Use(balanceEnv.Validate)
	{
		transactionRoutes.GET("/getBalance", balanceEnv.GetBalance)
		transactionRoutes.POST("/process", balanceEnv.ProcessTransaction)
		transactionRoutes.POST("/processInterAccount", balanceEnv.ProcessInterAccountTransaction)

	}

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Failed to start server")
	}
}

// load dependencies and setup enviorment
func setupEnviorment() *BalanceEnviormentService {
	// this should be a db pool ideally, and repositories would be tapping into the connections of the pool
	database, err := GetDatabase()
	if err != nil {
		return nil
	}
	// repositories initialization
	return &BalanceEnviormentService{
		BalanceManagement:         &utils.BalanceManagement{DB: database},
		CurrencyAccountManagement: &utils.CurrencyAccountManagement{DB: database},
		ForexManagement:           &utils.ForexManagement{DB: database},
	}

}

func GetDatabase() (*sql.DB, error) {
	// change these to env variables
	port := 5432
	url := "balance-db"
	password := "admin"
	db := "balance"
	user := "admin"

	postgresConnection, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", url, port, user, password, db))
	if err != nil {
		fmt.Println("Failed to connect to database " + err.Error())
	}
	return postgresConnection, nil
}
