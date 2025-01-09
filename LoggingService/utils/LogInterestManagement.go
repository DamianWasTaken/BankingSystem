package utils

import (
	"database/sql"
	"fmt"
	"time"
)

type LogInterestManagement struct {
	DB *sql.DB
}

func (interestLogManagement *LogInterestManagement) PersistInterestChange(interestLog InterestLog) error {

	query := fmt.Sprintf("INSERT INTO interestConfiguration (interestRate, created) VALUES (%f, '%s')", interestLog.Interest, time.Now().Format("2006-01-02 15:04:05"))
	_, err := interestLogManagement.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (interestLogManagement *LogInterestManagement) PersistInterestUserApplication(interestApplicationLog InterestApplicationLog) error {
	query := fmt.Sprintf("INSERT INTO interestUserApplication (email, currency, interestRate, frequency, created, outcome) VALUES ('%s', '%s', %f, '%s', '%s', '%s')", interestApplicationLog.Email, interestApplicationLog.Currency, interestApplicationLog.InterestRate, interestApplicationLog.Frequency, time.Now().Format("2006-01-02 15:04:05"), interestApplicationLog.Outcome)
	_, err := interestLogManagement.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
