package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"BalanceService/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// repositories
type BalanceEnviormentService struct {
	BalanceManagement interface {

		//composite key of email and currency
		GetAccountBalance(utils.GetBalanceRequest) (float32, error)

		// this needs to go through a queue
		ProcessTransaction(utils.ProcessTransactionRequest) error

		ProcessInterAccountTransaction(utils.ProcessInterAccountTransactionRequest, float32) error

		ProcessInterest(utils.ProcessInterestRequest) (float32, error)
	}
	CurrencyAccountManagement interface {
		CreateCurrencyAccount(utils.CreateCurrencyAccountRequest) error

		ProcessIntraAccountTransaction(utils.ProcessIntraAccountTransactionRequest, float32) error

		CheckIfAccountExists(string, string) (bool, error)
	}
	ForexManagement interface {
		GetForexRate(string, string) (float32, error)
	}
}

func (repositories *BalanceEnviormentService) GetBalance(c *gin.Context) {
	var GetBalanceRequest utils.GetBalanceRequest
	if err := c.ShouldBindJSON(&GetBalanceRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	balance, err := repositories.BalanceManagement.GetAccountBalance(GetBalanceRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"currency": GetBalanceRequest.Currency, "balance": balance})
}

func (repositories *BalanceEnviormentService) ProcessTransaction(c *gin.Context) {
	var ProcessTransactionRequest utils.ProcessTransactionRequest

	if err := c.ShouldBindJSON(&ProcessTransactionRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	err := repositories.BalanceManagement.ProcessTransaction(ProcessTransactionRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	logTransaction(ProcessTransactionRequest.Currency, ProcessTransactionRequest.Email, ProcessTransactionRequest.TransactionType, ProcessTransactionRequest.Value)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction processed"})
}

func (repositories *BalanceEnviormentService) CreateCurrencyAccount(c *gin.Context) {
	var CreateCurrencyAccountRequest utils.CreateCurrencyAccountRequest
	// once validated we can keep on with the request
	if err := c.ShouldBindBodyWith(&CreateCurrencyAccountRequest, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors1": []string{err.Error()}})
		return
	}

	err := repositories.CurrencyAccountManagement.CreateCurrencyAccount(CreateCurrencyAccountRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors2": []string{err.Error()}})
		return
	}
	logTransaction(CreateCurrencyAccountRequest.Currency, CreateCurrencyAccountRequest.Email, "created", CreateCurrencyAccountRequest.Balance)

	c.JSON(http.StatusCreated, gin.H{"message": "Currency account created"})
}

func (repositories *BalanceEnviormentService) ProcessIntraAccountTransaction(c *gin.Context) {
	var ProcessIntraAccountTransactionRequest utils.ProcessIntraAccountTransactionRequest

	if err := c.ShouldBindJSON(&ProcessIntraAccountTransactionRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	fromValue, err := repositories.BalanceManagement.GetAccountBalance(utils.GetBalanceRequest{Email: ProcessIntraAccountTransactionRequest.Email, Currency: ProcessIntraAccountTransactionRequest.FromCurrency})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if fromValue < ProcessIntraAccountTransactionRequest.Value {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{"Insufficient funds"}})
		return
	}

	rate, err := repositories.ForexManagement.GetForexRate(ProcessIntraAccountTransactionRequest.ToCurrency, ProcessIntraAccountTransactionRequest.FromCurrency)
	valueToAdd := ProcessIntraAccountTransactionRequest.Value * rate
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	err = repositories.CurrencyAccountManagement.ProcessIntraAccountTransaction(ProcessIntraAccountTransactionRequest, valueToAdd)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}
	logTransaction(ProcessIntraAccountTransactionRequest.FromCurrency, ProcessIntraAccountTransactionRequest.Email, "debit", ProcessIntraAccountTransactionRequest.Value)
	logTransaction(ProcessIntraAccountTransactionRequest.ToCurrency, ProcessIntraAccountTransactionRequest.Email, "credit", valueToAdd)
	c.JSON(http.StatusOK, gin.H{"message": "Intra account transaction processed"})
}

func (repositories *BalanceEnviormentService) ProcessInterAccountTransaction(c *gin.Context) {
	var ProcessInterAccountTransactionRequest utils.ProcessInterAccountTransactionRequest

	if err := c.ShouldBindJSON(&ProcessInterAccountTransactionRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	fromValue, err := repositories.BalanceManagement.GetAccountBalance(utils.GetBalanceRequest{Email: ProcessInterAccountTransactionRequest.Email, Currency: ProcessInterAccountTransactionRequest.FromCurrency})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if fromValue < ProcessInterAccountTransactionRequest.Value {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{"Insufficient funds"}})
		return
	}

	//check if account exists
	exists, err := repositories.CurrencyAccountManagement.CheckIfAccountExists(ProcessInterAccountTransactionRequest.ToEmail, ProcessInterAccountTransactionRequest.ToCurrency)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{"Account does not exist"}})
		return
	}
	if ProcessInterAccountTransactionRequest.FromCurrency != ProcessInterAccountTransactionRequest.ToCurrency {
		rate, err := repositories.ForexManagement.GetForexRate(ProcessInterAccountTransactionRequest.ToCurrency, ProcessInterAccountTransactionRequest.FromCurrency)
		valueToAdd := ProcessInterAccountTransactionRequest.Value * rate
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
			return
		}

		err = repositories.BalanceManagement.ProcessInterAccountTransaction(ProcessInterAccountTransactionRequest, valueToAdd)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
			return
		}
		logTransaction(ProcessInterAccountTransactionRequest.FromCurrency, ProcessInterAccountTransactionRequest.Email, "debit", ProcessInterAccountTransactionRequest.Value)
		logTransaction(ProcessInterAccountTransactionRequest.ToCurrency, ProcessInterAccountTransactionRequest.ToEmail, "credit", valueToAdd)

		c.JSON(http.StatusOK, gin.H{"message": "Inter account transaction processed"})
		return
	}

	err = repositories.BalanceManagement.ProcessInterAccountTransaction(ProcessInterAccountTransactionRequest, ProcessInterAccountTransactionRequest.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		return
	}

	// logging
	logTransaction(ProcessInterAccountTransactionRequest.FromCurrency, ProcessInterAccountTransactionRequest.Email, "debit", ProcessInterAccountTransactionRequest.Value)
	logTransaction(ProcessInterAccountTransactionRequest.ToCurrency, ProcessInterAccountTransactionRequest.ToEmail, "credit", ProcessInterAccountTransactionRequest.Value)

	c.JSON(http.StatusOK, gin.H{"message": "Inter account transaction processed"})
}

func (repositories *BalanceEnviormentService) ApplyInterest(c *gin.Context) {
	var interestRequest utils.ProcessInterestRequest

	if err := c.ShouldBindJSON(&interestRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		logInterestApplication(interestRequest.Email, interestRequest.Currency, "failed", interestRequest.Frequency, interestRequest.Interest)
		return
	}

	valueAdded, err := repositories.BalanceManagement.ProcessInterest(interestRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": []string{err.Error()}})
		logInterestApplication(interestRequest.Email, interestRequest.Currency, "failed", interestRequest.Frequency, interestRequest.Interest)
		return
	}
	logInterestApplication(interestRequest.Email, interestRequest.Currency, "success", interestRequest.Frequency, interestRequest.Interest)
	logTransaction(interestRequest.Currency, interestRequest.Email, "interest credit", valueAdded)
	c.JSON(http.StatusOK, gin.H{"message": "Interest applied"})
}

func (repositories *BalanceEnviormentService) ValidateJWT(c *gin.Context) {

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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": []string{"Please check your token or email, and make sure your account is active"}})
	}
}

func logTransaction(currency, email, transactionType string, value float32) error {
	log := utils.BalanceLog{Email: email, Currency: currency, Value: value, TransactionType: transactionType}
	jsonData, err := json.Marshal(log)
	if err != nil {
		return err
	}

	// correct way would be create transaction, send to db, log, then commit the transaction if log is successful otherwise rollback
	// too late now to refactor though
	http.Post("http://logging-service:8080/logging/balance", "application/json", bytes.NewBuffer(jsonData))
	return nil
}

func logInterestApplication(email, currency, outcome string, frequency int, interestRate float32) error {
	log := utils.InterestAppliedLog{Email: email, Currency: currency, Frequency: frequency, InterestRate: interestRate, Outcome: outcome}
	jsonData, err := json.Marshal(log)
	if err != nil {
		return err
	}

	http.Post("http://logging-service:8080/logging/interest/interestUserApplication", "application/json", bytes.NewBuffer(jsonData))
	return nil
}
