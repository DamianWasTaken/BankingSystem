package utils

type ModifyDailyInterestRateRequest struct {
	InterestRate float32 `json:"interestRate" binding:"required"`
}

type AddInterestUserRequest struct {
	Email     string `json:"email" binding:"required,max=200"`
	Currency  string `json:"currency" binding:"required,max=3"`
	Frequency int    `json:"frequency" binding:"required"`
}

type ModifyFrequencyRequest struct {
	InterestId int `json:"interestId" binding:"required"`
	Frequency  int `json:"frequency" binding:"required"`
}

type InterestUser struct {
	InterestId int
	Frequency  int
	Currency   string
	Email      string
}

type ProcessInterest struct {
	Email     string  `json:"email"`
	Currency  string  `json:"currency"`
	Interest  float32 `json:"interest"`
	Frequency int     `json:"frequency"`
}

type InterestLog struct {
	Interest float32 `json:"interest"`
}
