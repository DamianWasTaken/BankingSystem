package utils

import (
	"database/sql"
	"fmt"
	"time"
)

type LogBalanceManagement struct {
	DB *sql.DB
}

func (balanceLogManagement *LogBalanceManagement) PersistBalanceChange(balanceLog BalanceLog) error {
	query := fmt.Sprintf("INSERT INTO public.balanceTransactions (email, currency, amount, transactionType, created) VALUES ('%s', '%s', %f, '%s', '%s')", balanceLog.Email, balanceLog.Currency, balanceLog.Value, balanceLog.TransactionType, time.Now().Format(time.RFC3339))
	_, err := balanceLogManagement.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
