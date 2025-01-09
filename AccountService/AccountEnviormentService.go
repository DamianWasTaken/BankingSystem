package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"AccountService/utils"
)

// repositories
type AccountEnviormentSerivce struct {
	Secret         string
	UserManagement interface {
		CreateUser(utils.CreateUserRequest) error
		DeleteUser(utils.DeleteUserRequest) error
		LoginUser(utils.LoginRequest) (string, error)
		CheckIfEmailExists(email string, create bool) error
	}
	AccountManagement interface {
		DeactivateAccount(string) error
		ReactivateAccount(string) error
	}
}

// handlers/service methods
func (repositories *AccountEnviormentSerivce) CreateUser(c *gin.Context) {
	var newUser utils.CreateUserRequest

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.UserManagement.CheckIfEmailExists(newUser.Email, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	newUser.Password, err = argon2id.CreateHash(newUser.Password, argon2id.DefaultParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	err = repositories.UserManagement.CreateUser(newUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})

}
func (repositories *AccountEnviormentSerivce) DeleteUser(c *gin.Context) {
	var deleteUserRequest utils.DeleteUserRequest

	if err := c.ShouldBindJSON(&deleteUserRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.UserManagement.DeleteUser(deleteUserRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// this will be needed for other services to validate the user before doing any actions
func (repositories *AccountEnviormentSerivce) ValidateUser(c *gin.Context) {
	// we'll put this behind the jwt check middleware, so if it reaches here, the user is valid
	var ValidateTokenRequest utils.ValidateTokenRequest

	if err := c.ShouldBindJSON(&ValidateTokenRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}
	val, err := c.Get("email")
	if !err {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{"Email not found in claims"}})
		return
	}
	if val != ValidateTokenRequest.Email {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{"Claims email does not match request"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User validated"})
}

func (repositories *AccountEnviormentSerivce) DeactivateAccount(c *gin.Context) {
	var status utils.AccountStatusStateChange

	if err := c.ShouldBindJSON(&status); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.AccountManagement.DeactivateAccount(status.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	LogStatusChange(status.Email, "inactive")
	c.JSON(http.StatusOK, gin.H{"message": "Account deactivated"})

}
func (repositories *AccountEnviormentSerivce) ReactivateAccount(c *gin.Context) {
	var status utils.AccountStatusStateChange

	if err := c.ShouldBindJSON(&status); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.AccountManagement.ReactivateAccount(status.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	LogStatusChange(status.Email, "active")
	c.JSON(http.StatusOK, gin.H{"message": "Account Reactivated"})
}

func (repositories *AccountEnviormentSerivce) LoginUser(c *gin.Context) {
	var loginRequest utils.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	hashedPassword, err := repositories.UserManagement.LoginUser(loginRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	valid, err := argon2id.ComparePasswordAndHash(loginRequest.Password, hashedPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{"Invalid password"}})
		return
	}
	if !valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": []string{"Invalid password"}})
		return
	}

	// create a jwt token
	claims := jwt.RegisteredClaims{
		Subject:   loginRequest.Email,
		Issuer:    "payter",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(repositories.Secret))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "expiresAt": claims.ExpiresAt})
}

func (repositories *AccountEnviormentSerivce) AccountStatusHistory(c *gin.Context) {
	// this will make a call to the log service
}

func (e *AccountEnviormentSerivce) ValidateJWT(c *gin.Context) {

	HeaderToken := c.GetHeader("Authorization")

	token, err := jwt.ParseWithClaims(HeaderToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header)
		}
		return []byte(e.Secret), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
		return
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// Check the expiry on the token
		if claims.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": []string{"Token expired"}})
			return
		}

		email := claims.Subject

		err = e.UserManagement.CheckIfEmailExists(email, false)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": []string{"UserId invalid"}})
			return
		}
		c.Set("email", email)
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": []string{"Token invalid"}})
		return
	}
	c.Next()
}

func LogStatusChange(email string, status string) {
	log := utils.LogStatusChange{Email: email, Status: status}
	jsonData, _ := json.Marshal(log)
	http.Post("http://logging-service:8080/logging/account/status", "application/json", bytes.NewBuffer(jsonData))
}
