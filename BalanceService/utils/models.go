package utils

type GetBalanceRequest struct {
	Email    string `json:"email" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

type ProcessTransactionRequest struct {
	Email           string  `json:"email" binding:"required"`
	Currency        string  `json:"currency" binding:"required"`
	Value           float32 `json:"value" binding:"required"`
	TransactionType string  `json:"process" binding:"required" enum:"credit,debit"`
}

type ProcessIntraAccountTransactionRequest struct {
	Email        string  `json:"email" binding:"required"`
	FromCurrency string  `json:"fromCurrency" binding:"required"`
	ToCurrency   string  `json:"toCurrency" binding:"required"`
	Value        float32 `json:"value" binding:"required"`
}

type CreateCurrencyAccountRequest struct {
	Email    string  `json:"email" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
	Balance  float32 `json:"balance" binding:"required"`
}

type ProcessInterAccountTransactionRequest struct {
	Email        string  `json:"email" binding:"required,max=100"`
	ToEmail      string  `json:"toEmail" binding:"required,max=100"`
	FromCurrency string  `json:"fromCurrency" binding:"required,max=100"`
	ToCurrency   string  `json:"toCurrency" binding:"required,max=100"`
	Value        float32 `json:"value" binding:"required"`
}
type ProcessInterestRequest struct {
	Email     string  `json:"email" binding:"required"`
	Currency  string  `json:"currency" binding:"required"`
	Interest  float32 `json:"interest" binding:"required"`
	Frequency int     `json:"frequency" binding:"required"`
}
