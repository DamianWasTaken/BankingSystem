package utils

type StatusChangeLog struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}

type BalanceLog struct {
	Email           string  `json:"email"`
	Currency        string  `json:"currency"`
	Value           float32 `json:"value"`
	TransactionType string  `json:"transactionType"`
}

type InterestLog struct {
	Interest float32 `json:"interest"`
}

type InterestApplicationLog struct {
	Email        string  `json:"email"`
	Currency     string  `json:"currency"`
	InterestRate float32 `json:"interest"`
	Frequency    string  `json:"frequency"`
	Outcome      string  `json:"outcome"`
}
