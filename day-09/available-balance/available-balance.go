package main

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID      int64           `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}

type Hold struct {
	ID       int64           `json:"id"`
	WalletID int64           `json:"wallet_id"`
	Amount   decimal.Decimal `json:"amount"`
	Status   string          `json:"status"`
}

func CalculateAvailableBalance(wallet Wallet, holds []Hold) (decimal.Decimal, error) {
	if wallet.ID <= 0 {
		return decimal.Zero, errors.New("invalid wallet id")
	}

	if wallet.Balance.LessThan(decimal.Zero) {
		return decimal.Zero, errors.New("wallet balance must not be negative")
	}
	blocked := decimal.Zero

	for _, hold := range holds {
		if hold.WalletID != wallet.ID {
			continue
		}
		if hold.Status != "active" {
			continue
		}

		blocked = blocked.Add(hold.Amount)
	}

	if blocked.GreaterThan(wallet.Balance) {
		return decimal.Zero, errors.New("blocked amount exceeds wallet balance")
	}

	available := wallet.Balance.Sub(blocked)

	return available, nil

}

func main() {
	wallet := Wallet{
		ID:      1,
		Balance: decimal.RequireFromString("25000.00"),
	}

	holds := []Hold{
		{
			ID:       1,
			WalletID: 1,
			Amount:   decimal.RequireFromString("5000.00"),
			Status:   "active",
		},

		{
			ID:       2,
			WalletID: 1,
			Amount:   decimal.RequireFromString("2500.00"),
			Status:   "released",
		},

		{
			ID:       3,
			WalletID: 1,
			Amount:   decimal.RequireFromString("1250.00"),
			Status:   "active",
		},

		{
			ID:       4,
			WalletID: 2,
			Amount:   decimal.RequireFromString("22000.00"),
			Status:   "active",
		},
	}

	result, err := CalculateAvailableBalance(wallet, holds)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("available balance for wallet %d: %s\n", wallet.ID, result.StringFixed(2))
}
