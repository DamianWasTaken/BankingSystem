package utils

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type AuthManagment struct {
	DB *sql.DB
}

func (auth *AuthManagment) ValidateUser(c *gin.Context) {}
