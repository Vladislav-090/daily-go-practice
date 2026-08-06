package main

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID      int64           `json:"id"`
	UserID  int64           `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
}

type TransferService struct{}

var (
	ErrSameWallet         = errors.New("transfer to same wallet")
	ErrInvalidFromWalletID = errors.New("invalid from wallet id")
	ErrInvalidToWalletID   = errors.New("invalid to wallet id")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrInvalidAmount       = errors.New("amount must be positive")
	ErrWalletIsNil         = errors.New("wallet is nil")
)

func (t *TransferService) Transfer(from *Wallet, to *Wallet, amount decimal.Decimal) error {
	if from == nil || to == nil {
		return ErrWalletIsNil
	}

	if from.ID <= 0 {
		return ErrInvalidFromWalletID
	}

	if to.ID <= 0 {
		return ErrInvalidToWalletID
	}

	if from.ID == to.ID {
		return ErrSameWallet
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	if from.Balance.LessThan(amount) {
		return ErrInsufficientFunds
	}

	from.Balance = from.Balance.Sub(amount)

	to.Balance = to.Balance.Add(amount)

	return nil
}

func main() {
	fromWallet := Wallet{
		ID:      1,
		UserID:  1,
		Balance: decimal.RequireFromString("25000.00"),
	}

	toWallet := Wallet{
		ID:      2,
		UserID:  2,
		Balance: decimal.RequireFromString("1000"),
	}

	fmt.Println("from wallet balance:", fromWallet.Balance)
	fmt.Println("to wallet balance:", toWallet.Balance)

	service := TransferService{}

	err := service.Transfer(
		&fromWallet,
		&toWallet,
		decimal.RequireFromString("8000.00"),
	)

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("--------------\nafter transfer\n---------------")

	fmt.Println("updated from wallet balance:", fromWallet.Balance)
	fmt.Println("updated to wallet balance:", toWallet.Balance)
}
