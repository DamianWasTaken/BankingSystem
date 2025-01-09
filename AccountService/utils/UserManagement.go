package utils

import (
	"database/sql"
	"fmt"
)

type UserManagement struct {
	DB *sql.DB
	//could implement a map to store emails to check if they exist without db and at O(1), map existing on init
}

func (user *UserManagement) CreateUser(newUser CreateUserRequest) error { // TODO: comeback here once balance service is up
	query := fmt.Sprintf("INSERT INTO public.user (email, name, password, status) VALUES ('%s', '%s', '%s', '%s')", newUser.Email, newUser.Name, newUser.Password, "active")
	_, err := user.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error when creating user: %w", err)
	}
	return nil
}

func (user *UserManagement) DeleteUser(deleteUser DeleteUserRequest) error {
	query := fmt.Sprintf("DELETE FROM public.user WHERE email = '%s'", deleteUser.Email)
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

func (user *UserManagement) LoginUser(loginUser LoginRequest) (string, error) {
	query := fmt.Sprintf("SELECT password FROM public.user WHERE email = '%s'", loginUser.Email)
	sqlRow := user.DB.QueryRow(query)
	var hashedPassword string
	err := sqlRow.Scan(&hashedPassword)
	if err != nil {
		return "", fmt.Errorf("error when checking if user exists: %w", err)
	}
	return hashedPassword, nil
}
