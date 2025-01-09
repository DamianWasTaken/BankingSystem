package utils

import (
	"database/sql"
	"fmt"
	"time"
)

type LogAccountManagement struct {
	DB *sql.DB
}

func (accountLog *LogAccountManagement) PersistStatusChange(statusLog StatusChangeLog) error {
	query := fmt.Sprintf("INSERT INTO accountStatus (email, status, created) VALUES ('%s', '%s', '%s')", statusLog.Email, statusLog.Status, time.Now().Format("2006-01-02 15:04:05"))
	_, err := accountLog.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (accountLog *LogAccountManagement) GetAccountStatusHistory(email string) error {
	// query := fmt.Sprintf("SELECT status, created FROM accountStatus WHERE email = '%s' order by created", email)

	return nil
}
