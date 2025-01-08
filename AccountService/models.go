package main

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