package utils

type CreateUserRequest struct {
	Email        string  `json:"email" binding:"required"`
	Password     string  `json:"password" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	DepositValue float32 `json:"depositValue"`
}

type DeactivateAccountRequest struct {
	Email int `json:"email" binding:"required"`
}

type ReactivateAccountRequest struct {
	Email int `json:"email" binding:"required"`
}

type DeleteUserRequest struct {
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,max=100"`
	Password string `json:"password" binding:"required,max=100"`
}

type ValidateTokenRequest struct {
	Email string `json:"email" binding:"required,max=100"`
}

type LogStatusChange struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}

type AccountStatusStateChange struct {
	Email string `json:"email"`
}
