package utils

import (
	"database/sql"
	"fmt"
)

type AccountManagement struct {
	DB *sql.DB
}

func (account *AccountManagement) DeactivateAccount(email string) error {
	query := fmt.Sprintf("UPDATE public.user SET status = 'inactive' WHERE email = '%s'", email)
	_, err := account.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (account *AccountManagement) ReactivateAccount(email string) error {
	query := fmt.Sprintf("UPDATE public.user SET status = 'active' WHERE email = '%s'", email)
	_, err := account.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (account *AccountManagement) IsAccountActive(email string) (bool, error) {
	query := fmt.Sprintf("SELECT status FROM public.user WHERE email = '%s'", email)
	sqlRow := account.DB.QueryRow(query)
	var status string
	err := sqlRow.Scan(&status)
	if err != nil {
		return false, err
	}
	if status == "active" {
		return true, nil
	}
	return false, nil
}
