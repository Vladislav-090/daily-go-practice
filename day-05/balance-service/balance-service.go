package main

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID     int64           `json:"id"`
	UserID int64           `json:"user_id"`
	Amount decimal.Decimal `json:"amount"`
	Type   string          `json:"type"`
}

func CalculateUserBalance(transactions []Transaction, userID int64) (decimal.Decimal, error) {
	if userID <= 0 {
		return decimal.Zero, errors.New("invalid user id")
	}

	balance := decimal.Zero

	for _, transaction := range transactions {
		if transaction.UserID != userID {
			continue
		}
		switch transaction.Type {
		case "income":
			balance = balance.Add(transaction.Amount)
		case "expense":
			balance = balance.Sub(transaction.Amount)
		default:
			return decimal.Zero, errors.New("invalid transaction type")
		}

	}

	return balance, nil
}

func main() {
	transactions := []Transaction{
		{
			ID:     1,
			UserID: 1,
			Amount: decimal.RequireFromString("1000.00"),
			Type:   "income",
		},

		{
			ID:     2,
			UserID: 1,
			Amount: decimal.RequireFromString("995.00"),
			Type:   "expense",
		},
	}

	balance, err := CalculateUserBalance(transactions, 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated balance:", balance.StringFixed(2))
}
