package utils

import (
	"database/sql"
	"fmt"
)

type UserManagement struct {
	DB *sql.DB
	//could implement a map to store emails to check if they exist without db and at O(1), map existing on init
}

func (user *UserManagement) CreateUser(email, name, password string, depositValue float32) error { // TODO: comeback here once balance service is up
	query := fmt.Sprintf("INSERT INTO public.user (email, name, password, active) VALUES ('%s', '%s', '%s', %t)", email, name, password, true)
	_, err := user.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error when creating user: %w", err)
	}
	return nil
}

func (user *UserManagement) DeleteUser(email string) error {
	query := fmt.Sprintf("DELETE FROM public.user WHERE email = '%s'", email)
	_, err := user.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error when deleting user: %w", err)
	}

	return nil
}

func (user *UserManagement) CheckIfEmailExists(email string, create bool) error {
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM public.user WHERE email = '%s')", email)
	sqlRow := user.DB.QueryRow(query)
	var exists bool
	err := sqlRow.Scan(&exists)
	if err != nil {
		return fmt.Errorf("error when checking if user exists: %w", err)
	}
	if exists && create {
		return fmt.Errorf("email already exists")
	}

	return nil
}

func (user *UserManagement) CheckIfIdExists(id int) error {
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM public.user WHERE userId = %d)", id)
	sqlRow := user.DB.QueryRow(query)
	var exists bool
	err := sqlRow.Scan(&exists)
	if err != nil {
		return fmt.Errorf("error when checking if userId exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("userId does not exist")
	}
	return nil
}

func (user *UserManagement) LoginUser(email, password string) (string, error) {
	query := fmt.Sprintf("SELECT password FROM public.user WHERE email = '%s'", email)
	sqlRow := user.DB.QueryRow(query)
	var hashedPassword string
	err := sqlRow.Scan(&hashedPassword)
	if err != nil {
		return "", fmt.Errorf("error when checking if user exists: %w", err)
	}
	return hashedPassword, nil
}
