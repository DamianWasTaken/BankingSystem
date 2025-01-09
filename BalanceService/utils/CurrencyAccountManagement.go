package utils

import (
	"database/sql"
	"fmt"
)

type CurrencyAccountManagement struct {
	DB *sql.DB
}

func (currencyAccount *CurrencyAccountManagement) CreateCurrencyAccount(newCurrencyAccount CreateCurrencyAccountRequest) error {
	query := fmt.Sprintf("INSERT INTO public.account_currencies (email, currency, balance) VALUES ('%s', '%s', %f)", newCurrencyAccount.Email, newCurrencyAccount.Currency, newCurrencyAccount.Balance)

	_, err := currencyAccount.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error when creating currency account: %w", err)
	}
	return nil
}

func (currencyAccount *CurrencyAccountManagement) ProcessIntraAccountTransaction(process ProcessIntraAccountTransactionRequest, value float32) error {
	transaction, err := currencyAccount.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	debitQuery := fmt.Sprintf("UPDATE public.account_currencies SET balance = balance - %f WHERE email = '%s' AND currency = '%s'", process.Value, process.Email, process.FromCurrency)
	creditQuery := fmt.Sprintf("UPDATE public.account_currencies SET balance = balance + %f WHERE email = '%s' AND currency = '%s'", value, process.Email, process.ToCurrency)

	_, err = transaction.Exec(debitQuery)
	if err != nil {
		return fmt.Errorf("error when debiting account: %w", err)
	}
	_, err = transaction.Exec(creditQuery)
	if err != nil {
		return fmt.Errorf("error when crediting account: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}
	return nil
}

func (currencyAccount *CurrencyAccountManagement) CheckIfAccountExists(email string, currency string) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM public.account_currencies WHERE email = '%s' AND currency = '%s')", email, currency)
	sqlRow := currencyAccount.DB.QueryRow(query)
	var exists bool
	err := sqlRow.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error when checking if currency account exists: %w", err)
	}
	return exists, nil
}
