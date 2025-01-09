package utils

import (
	"database/sql"
	"fmt"
	"time"
)

type InterestUserManagement struct {
	DB *sql.DB
}

func (interestUser *InterestUserManagement) AddInterestUser(request AddInterestUserRequest) error {
	now := time.Now().AddDate(0, 0, request.Frequency).Format("20060102")
	query := fmt.Sprintf("INSERT INTO public.interestUser (email, accountcurrency, interestFrequency, nextInterestDate) VALUES ('%s', '%s', %d, '%s')", request.Email, request.Currency, request.Frequency, now)
	_, err := interestUser.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (interestUser *InterestUserManagement) ModifyInterestUserFrequency(rate ModifyFrequencyRequest) error {
	nextInterest := time.Now().AddDate(0, 0, rate.Frequency).Format("20060102")
	query := fmt.Sprintf("UPDATE public.interestUser SET interestFrequency = %d, nextInterestDate= '%s' where interestId=%d", rate.Frequency, nextInterest, rate.InterestId)
	_, err := interestUser.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (interestUser *InterestUserManagement) UpdateInterestUserDate(interestId int, frequency int) error {
	date := time.Now().AddDate(0, 0, frequency).Format("20060102")
	query := fmt.Sprintf("UPDATE public.interestUser SET nextInterestDate = '%s' WHERE interestId = %d", date, interestId)
	_, err := interestUser.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (interestUser *InterestUserManagement) GetInterestRateUsers(date string) ([]InterestUser, error) {
	query := fmt.Sprintf("SELECT interestId, email, accountCurrency, interestFrequency FROM public.interestUser WHERE nextInterestDate = '%s'", date)
	sqlRows, err := interestUser.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	var interestUsers []InterestUser
	for sqlRows.Next() {
		var interestUser InterestUser
		err := sqlRows.Scan(&interestUser.InterestId, &interestUser.Email, &interestUser.Currency, &interestUser.Frequency)
		if err != nil {
			return nil, err
		}
		interestUsers = append(interestUsers, interestUser)
	}
	return interestUsers, nil
}
