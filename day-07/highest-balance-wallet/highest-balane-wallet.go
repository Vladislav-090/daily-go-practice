package main

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID      int64           `json:"id"`
	UserID  int64           `json:"user_id"`
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
}

func FindHighestBalance(wallets []Wallet) (Wallet, error) {

	if len(wallets) == 0 {
		return Wallet{}, errors.New("dont have any wallets")
	}
	richestWallet := wallets[0]

	for _, wallet := range wallets[1:] {
		if wallet.Balance.GreaterThan(richestWallet.Balance) {
			richestWallet = wallet
		}
	}
	return richestWallet, nil
}

func main() {
	wallets := []Wallet{
		{
			ID:      1,
			UserID:  1,
			Name:    "Kaspi",
			Balance: decimal.RequireFromString("20000.00"),
		},

		{
			ID:      2,
			UserID:  1,
			Name:    "Halyk",
			Balance: decimal.RequireFromString("15000.00"),
		},

		{
			ID:      3,
			UserID:  1,
			Name:    "Freedom",
			Balance: decimal.RequireFromString("25000.00"),
		},

		{
			ID:      4,
			UserID:  2,
			Name:    "Freedom",
			Balance: decimal.RequireFromString("99000.00"),
		},
	}

	richestWallet, err := FindHighestBalance(wallets)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("id:", richestWallet.ID)
	fmt.Println("user_id:", richestWallet.UserID)
	fmt.Println("name:", richestWallet.Name)
	fmt.Println("balance:", richestWallet.Balance.StringFixed(2))

}
