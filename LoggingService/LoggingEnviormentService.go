package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"LoggingService/utils"
)

type LoggingEnviormentService struct {
	LogAccountManagement interface {
		PersistStatusChange(utils.StatusChangeLog) error
		GetAccountStatusHistory(string) error
	}
	LogBalanceManagement interface {
		PersistBalanceChange(utils.BalanceLog) error
	}
	LogInterestManagement interface {
		PersistInterestChange(utils.InterestLog) error
		PersistInterestUserApplication(utils.InterestApplicationLog) error
	}
}

func (repositories *LoggingEnviormentService) PersistStatusChange(c *gin.Context) {
	var PersistStatusChangeRequest utils.StatusChangeLog
	if err := c.ShouldBindJSON(&PersistStatusChangeRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.LogAccountManagement.PersistStatusChange(PersistStatusChangeRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "log successful"})
}

func (repositories *LoggingEnviormentService) PersistBalanceChange(c *gin.Context) {
	var PersistBalanceChangeRequest utils.BalanceLog
	if err := c.ShouldBindJSON(&PersistBalanceChangeRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.LogBalanceManagement.PersistBalanceChange(PersistBalanceChangeRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "log successful"})
}

func (repositories *LoggingEnviormentService) PersistInterestChange(c *gin.Context) {
	var PersistInterestChangeRequest utils.InterestLog
	if err := c.ShouldBindJSON(&PersistInterestChangeRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.LogInterestManagement.PersistInterestChange(PersistInterestChangeRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "log successful"})
}

func (repositories *LoggingEnviormentService) PersistInterestUserApplication(c *gin.Context) {
	var PersistInterestUserApplicationRequest utils.InterestApplicationLog
	if err := c.ShouldBindJSON(&PersistInterestUserApplicationRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.LogInterestManagement.PersistInterestUserApplication(PersistInterestUserApplicationRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "log successful"})
}

func (repositories *LoggingEnviormentService) GetAccountStatusHistory(c *gin.Context) {
	email := c.Query("email")

	err := repositories.LogAccountManagement.GetAccountStatusHistory(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "log successful"})
}

func (repositories *LoggingEnviormentService) ValidateJWT(c *gin.Context) {

	ByteBody, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(ByteBody))

	type emailReader struct {
		Email string `json:"email" binding:"required"`
	}

	var email emailReader

	err := json.Unmarshal(ByteBody, &email)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{"Error on parsing request validation bytes"}})
	}
	// if err := c.ShouldBindJSON(&email); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
	// 	return
	// }

	requestBody, err := json.Marshal(map[string]string{
		"email": email.Email,
	})

	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	HeaderToken := c.GetHeader("Authorization")

	//validate token
	client := &http.Client{}
	url := "http://account-service:8080/auth/validate"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": []string{"Error on validation http setup"}})

	}

	req.Header.Set("Authorization", HeaderToken)

	resp, err := client.Do(req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": []string{"Error on sending validation http"}})
	}

	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": []string{"Invalid token"}})
	}
}
