package main

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

type Transfer struct {
	ID           int64           `json:"id"`
	FromWalletID int64           `json:"from_wallet_id"`
	ToWalletID   int64           `json:"to_wallet_id"`
	Amount       decimal.Decimal `json:"amount"`
	Status       string          `json:"status"`
}

var (
	ErrInvalidWalletID      = errors.New("invalid wallet id")
	ErrTransferStatusFailed = errors.New("transfer status is failed")
)

func CalculateTransferredAmount(transfers []Transfer, walletID int64) (decimal.Decimal, error) {
	if walletID <= 0 {
		return decimal.Decimal{}, ErrInvalidWalletID
	}

	summ := decimal.Zero

	for _, transfer := range transfers {
		if transfer.FromWalletID != walletID {
			continue
		}
		if transfer.Status != "completed" {
			continue
		}

		summ = summ.Add(transfer.Amount)
	}
	return summ, nil

}

func main() {
	transfers := []Transfer{
		{
			ID:           1,
			FromWalletID: 1,
			ToWalletID:   2,
			Amount:       decimal.RequireFromString("400.00"),
			Status:       "completed",
		},

		{
			ID:           2,
			FromWalletID: 1,
			ToWalletID:   3,
			Amount:       decimal.RequireFromString("350.00"),
			Status:       "pending",
		},

		{
			ID:           3,
			FromWalletID: 1,
			ToWalletID:   4,
			Amount:       decimal.RequireFromString("1000.00"),
			Status:       "failed",
		},

		{
			ID:           3,
			FromWalletID: 1,
			ToWalletID:   4,
			Amount:       decimal.RequireFromString("1000.00"),
			Status:       "completed",
		},

		{
			ID:           3,
			FromWalletID: 1,
			ToWalletID:   5,
			Amount:       decimal.RequireFromString("935.00"),
			Status:       "completed",
		},
	}

	result, err := CalculateTransferredAmount(transfers, 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("summ of transfers:", result)
}
