package utils

import (
	"database/sql"
)

type AccountManagement struct {
	DB *sql.DB
}

func (account *AccountManagement) DeactivateAccount() {
	return
}

func (account *AccountManagement) ReactivateAccount() {
	return
}
