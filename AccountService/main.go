package main

import (
	"database/sql"
	"fmt"

	"AccountService/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()

	accEnv := setupEnviorment()

	userRoutes := r.Group("/user/")

	userRoutes.POST("/create", accEnv.CreateUser)
	userRoutes.DELETE("/delete", accEnv.DeleteUser)
	userRoutes.POST("/login", accEnv.LoginUser)

	//auth there is for authentication, not authorization, thus why login is on the user routes
	authRoutes := r.Group("/auth/")

	authRoutes.Use(accEnv.ValidateJWT)
	{
		authRoutes.POST("/validate", accEnv.ValidateUser)
	}

	//admin routes
	accountRoutes := r.Group("/account/")

	//unprotected on pourpose, otherwise I'd have to implement user roles
	accountRoutes.PATCH("/deactivate", accEnv.DeactivateAccount)
	accountRoutes.PATCH("/reactivate", accEnv.ReactivateAccount)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Failed to start server")
	}
}

// load dependencies and setup enviorment
func setupEnviorment() *AccountEnviormentSerivce {
	// this should be a db pool ideally, and repositories would be tapping into the connections of the pool
	database, err := GetDatabase()
	if err != nil {
		return nil
	}
	// repositories initialization
	return &AccountEnviormentSerivce{
		Secret:            "w1XyDEXb5/lTYLHiw768tCknKpMPvjqAEPxm0wXjRuw=", // SHA265 hash of "secret", should be an eviorment variable variable
		UserManagement:    &utils.UserManagement{DB: database},
		AccountManagement: &utils.AccountManagement{DB: database},
	}

}

func GetDatabase() (*sql.DB, error) {
	// change these to env variables
	port := 5432
	url := "account-db"
	password := "admin"
	db := "account"
	user := "admin"

	postgresConnection, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", url, port, user, password, db))
	if err != nil {
		fmt.Println("Failed to connect to database " + err.Error())
	}
	return postgresConnection, nil
}
