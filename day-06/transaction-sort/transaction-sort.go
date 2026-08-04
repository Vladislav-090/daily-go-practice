package main

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID     int64           `json:"id"`
	UserID int64           `json:"user_id"`
	Amount decimal.Decimal `json:"amount"`
}

func SortTransactionsByAmount(transactions []Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Amount.LessThan(transactions[j].Amount)
	})
}

func main() {
	transactions := []Transaction{
		{
			ID:     1,
			UserID: 1,
			Amount: decimal.RequireFromString("1000.00"),
		},

		{
			ID:     2,
			UserID: 1,
			Amount: decimal.RequireFromString("400.00"),
		},

		{
			ID:     3,
			UserID: 1,
			Amount: decimal.RequireFromString("1500"),
		},

		{
			ID:     4,
			UserID: 1,
			Amount: decimal.RequireFromString("200.00"),
		},
	}

	fmt.Println("before sorting:")

	for _, transaction := range transactions {
		fmt.Println(
			"id", transaction.ID,
			"amount", transaction.Amount.StringFixed(2),
		)
	}

	SortTransactionsByAmount(transactions)

	fmt.Println("after sorting:")

	for _, transaction := range transactions {
		fmt.Println(
			"id", transaction.ID,
			"amount", transaction.Amount.StringFixed(2),
		)
	}
}
