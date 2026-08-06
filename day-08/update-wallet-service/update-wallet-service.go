package main

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID      int64           `json:"id"`
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
}

func ApplyBonus(wallets []Wallet, minBalance decimal.Decimal, bonus decimal.Decimal) int {
	if bonus.LessThanOrEqual(decimal.RequireFromString("0")) {
		return 0
	}

	count := 0

	for i := range wallets {
		if wallets[i].Balance.GreaterThanOrEqual(minBalance) {
			wallets[i].Balance = wallets[i].Balance.Add(bonus)
			count++
		}
	}

	return count
}

func main() {
	wallets := []Wallet{
		{
			ID:      1,
			Name:    "kaspi",
			Balance: decimal.RequireFromString("1000.00"),
		},

		{
			ID:      2,
			Name:    "halyk",
			Balance: decimal.RequireFromString("2000"),
		},

		{
			ID:      3,
			Name:    "freedom",
			Balance: decimal.RequireFromString("400"),
		},
	}

	bonusWalletCount := ApplyBonus(wallets, decimal.RequireFromString("450.00"), decimal.RequireFromString("250"))

	fmt.Println("count of wallets which got bonuses:", bonusWalletCount)

	for _, wallet := range wallets {
		fmt.Println(wallet.Name, wallet.Balance)
	}
}
