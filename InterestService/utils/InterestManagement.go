package utils

import (
	"database/sql"
	"fmt"
)

type InterestManagement struct {
	DB *sql.DB
}

func (interest *InterestManagement) ModifyDailyInterestRate(rate ModifyDailyInterestRateRequest) error {
	query := fmt.Sprintf("UPDATE public.interestRate SET interest = %f", rate.InterestRate)
	_, err := interest.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (interest *InterestManagement) GetDailyInterestRate() (float32, error) {
	query := "SELECT interest FROM public.interestRate"
	sqlRow := interest.DB.QueryRow(query)
	var interestRate float32
	err := sqlRow.Scan(&interestRate)
	if err != nil {
		return 0, err
	}
	return interestRate, nil
}
