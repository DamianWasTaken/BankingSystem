package utils

import (
	"database/sql"
	"fmt"
)

type AccountManagement struct {
	DB *sql.DB
}

func (account *AccountManagement) DeactivateAccount(email string) error {
	query := fmt.Sprintf("UPDATE account SET status = 'inactive' WHERE email = '%s'", email)
	_, err := account.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (account *AccountManagement) ReactivateAccount(email string) error {
	query := fmt.Sprintf("UPDATE account SET status = 'active' WHERE email = '%s'", email)
	_, err := account.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
