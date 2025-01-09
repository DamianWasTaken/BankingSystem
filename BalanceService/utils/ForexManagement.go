package utils

import (
	"database/sql"
	"fmt"
)

type ForexManagement struct {
	DB *sql.DB
}

func (forex *ForexManagement) GetForexRate(fromCurrency string, toCurrency string) (float32, error) {
	fromQuery := fmt.Sprintf("SELECT rate FROM public.currencyList WHERE currency = '%s'", fromCurrency)
	toQuery := fmt.Sprintf("SELECT rate FROM public.currencyList WHERE currency = '%s'", toCurrency)
	sqlRow := forex.DB.QueryRow(fromQuery)
	var fromRate float32
	err := sqlRow.Scan(&fromRate)
	if err != nil {
		return 0, fmt.Errorf("error when getting from forex rate: %w", err)
	}
	sqlRow = forex.DB.QueryRow(toQuery)
	var toRate float32
	err = sqlRow.Scan(&toRate)
	if err != nil {
		return 0, fmt.Errorf("error when getting to forex rate: %w", err)
	}

	rate := fromRate / toRate

	return rate, nil
}
