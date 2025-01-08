package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// repositories
type AccountEnviormentSerivce struct {
	Secret         string
	UserManagement interface {
		CreateUser(string, string, string, float32) error
		DeleteUser(email string) error
		LoginUser(email, password string) (string, error)
		CheckIfEmailExists(email string, create bool) error
		CheckIfIdExists(id int) error
	}
	AuthManagment interface {
		ValidateUser(c *gin.Context)
	}
	AccountManagment interface {
		DeactivateAccount()
		ReactivateAccount()
	}
}

// handlers/service methods
func (repositories *AccountEnviormentSerivce) CreateUser(c *gin.Context) {
	var newUser CreateUserRequest

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.UserManagement.CheckIfEmailExists(newUser.Email, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	hashedPassword, err := argon2id.CreateHash(newUser.Password, argon2id.DefaultParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	err = repositories.UserManagement.CreateUser(newUser.Email, newUser.Name, hashedPassword, newUser.DepositValue)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created"})

}
func (repositories *AccountEnviormentSerivce) DeleteUser(c *gin.Context) {
	var deleteUserRequest struct {
		UserId string `json:"email"`
	}

	if err := c.ShouldBindJSON(&deleteUserRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.UserManagement.DeleteUser(deleteUserRequest.UserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// this will be needed for other services to validate the user before doing any actions
func (repositories *AccountEnviormentSerivce) ValidateUser(c *gin.Context) {
	// we'll put this behind the jwt check middleware, so if it reaches here, the user is valid
	c.JSON(http.StatusOK, gin.H{"message": "User validated"})
}

func (repositories *AccountEnviormentSerivce) DeactivateAccount(c *gin.Context) {

}
func (repositories *AccountEnviormentSerivce) ReactivateAccount(c *gin.Context) {

}

func (repositories *AccountEnviormentSerivce) LoginUser(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	hashedPassword, err := repositories.UserManagement.LoginUser(loginRequest.Email, loginRequest.Password)
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

func (e *AccountEnviormentSerivce) CheckJWT(c *gin.Context) {

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
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": []string{"Token invalid"}})
		return
	}
	c.Next()
}
