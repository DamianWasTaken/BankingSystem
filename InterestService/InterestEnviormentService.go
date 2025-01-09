package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"InterestService/utils"
)

type InterestEnviormentService struct {
	InterestManagement interface {
		ModifyDailyInterestRate(utils.ModifyDailyInterestRateRequest) error

		GetDailyInterestRate() (float32, error)
	}

	InterestUserManagement interface {
		AddInterestUser(utils.AddInterestUserRequest) error
		ModifyInterestUserFrequency(utils.ModifyFrequencyRequest) error
		GetInterestRateUsers(string) ([]utils.InterestUser, error)
		UpdateInterestUserDate(int, int) error
	}
}

func (repositories *InterestEnviormentService) AddInterestUser(c *gin.Context) {
	var newUser utils.AddInterestUserRequest

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.InterestUserManagement.AddInterestUser(newUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "interest user created"})
}

func (repositories *InterestEnviormentService) ModifyInterestUserFrequency(c *gin.Context) {
	var modifyUser utils.ModifyFrequencyRequest

	if err := c.ShouldBindJSON(&modifyUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.InterestUserManagement.ModifyInterestUserFrequency(modifyUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "interest frequency modified"})
}

func (repositories *InterestEnviormentService) ModifyDailyInterestRate(c *gin.Context) {
	var interestRateRequest utils.ModifyDailyInterestRateRequest

	if err := c.ShouldBindJSON(&interestRateRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.InterestManagement.ModifyDailyInterestRate(interestRateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Interest rate modified"})
}

func (repositories *InterestEnviormentService) ValidateJWT(c *gin.Context) {

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

func (repositories *InterestEnviormentService) ProcessInterest() {
	// get the daily interest rate
	dailyInterestRate, err := repositories.InterestManagement.GetDailyInterestRate()
	if err != nil {
		fmt.Println("Failed to get daily interest rate")
		return
	}

	now := time.Now().Format("20060102")

	// get all users
	users, err := repositories.InterestUserManagement.GetInterestRateUsers(now)
	if err != nil {
		fmt.Println("Failed to get users")
		return
	}

	//validate token
	client := &http.Client{}
	url := "http://balance-service:8080/interest/applyInterest"

	for _, user := range users {
		requestBody := utils.ProcessInterest{
			Email:     user.Email,
			Currency:  user.Currency,
			Interest:  dailyInterestRate,
			Frequency: user.Frequency,
		}

		jsonRequestBody, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Println("Failed to marshal body")

			// logs and event driven events
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequestBody))

		if err != nil {
			fmt.Println("Failed to create request to balance service")
			// logs and event driven events
		}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Failed to send request to balance service")
			// logs and event driven events
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println("Failed to apply interest, code:", resp.StatusCode)
			// logs and event driven events
		}

	}

}
