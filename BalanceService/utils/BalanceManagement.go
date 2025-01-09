package utils

import (
	"database/sql"
	"fmt"
)

type BalanceManagement struct {
	DB *sql.DB
}

func (balance *BalanceManagement) GetAccountBalance(balanceRequest GetBalanceRequest) (float32, error) {
	query := fmt.Sprintf("SELECT balance FROM public.account_currencies WHERE email = '%s' AND currency = '%s'", balanceRequest.Email, balanceRequest.Currency)
	sqlRow := balance.DB.QueryRow(query)
	var dbBalance float32
	err := sqlRow.Scan(&dbBalance)
	if err != nil {
		return 0, fmt.Errorf("balance account does not exist: %w", err)
	}
	return dbBalance, nil
}

func (balance *BalanceManagement) ProcessTransaction(processRequest ProcessTransactionRequest) error {
	query := ""
	if processRequest.TransactionType == "credit" {
		query = fmt.Sprintf("UPDATE public.account_currencies SET balance = balance + %f WHERE email = '%s' AND currency = '%s'", processRequest.Value, processRequest.Email, processRequest.Currency)
	} else {
		query = fmt.Sprintf("UPDATE public.account_currencies SET balance = balance - %f WHERE email = '%s' AND currency = '%s'", processRequest.Value, processRequest.Email, processRequest.Currency)
	}
	_, err := balance.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error when processing transaction: %w", err)
	}
	return nil
}

func (balance *BalanceManagement) ProcessInterAccountTransaction(processRequest ProcessInterAccountTransactionRequest, value float32) error {
	transaction, err := balance.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	debitQuery := fmt.Sprintf("UPDATE public.account_currencies SET balance = balance - %f WHERE email = '%s' AND currency = '%s'", processRequest.Value, processRequest.Email, processRequest.FromCurrency)
	creditQuery := fmt.Sprintf("UPDATE public.account_currencies SET balance = balance + %f WHERE email = '%s' AND currency = '%s'", value, processRequest.ToEmail, processRequest.ToCurrency)

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

func (balance *BalanceManagement) ProcessInterest(InterestObject ProcessInterestRequest) (float32, error) {
	overallInterest := float32(InterestObject.Frequency) * (InterestObject.Interest / 100)
	selectQuery := fmt.Sprintf("SELECT balance FROM public.account_currencies WHERE email = '%s' AND currency = '%s'", InterestObject.Email, InterestObject.Currency)
	sqlRow := balance.DB.QueryRow(selectQuery)
	var dbBalance float32
	err := sqlRow.Scan(&dbBalance)
	if err != nil {
		return 0, fmt.Errorf("error when processing interest: %w", err)
	}
	valueAdded := overallInterest * dbBalance
	query := fmt.Sprintf("UPDATE public.account_currencies SET balance = balance * %f WHERE email = '%s' AND currency = '%s'", overallInterest, InterestObject.Email, InterestObject.Currency)
	_, err = balance.DB.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("error when processing interest: %w", err)
	}
	return valueAdded, nil
}
