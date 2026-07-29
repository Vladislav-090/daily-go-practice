package main

import (
	"errors"
	"fmt"
)

type Wallet struct {
	ID      int64
	Name    string
	Balance int32
}

func (w *Wallet) WalletDeposit(amount int32) error {
	if amount <= 0 {
		return  errors.New("amount must be positive")
	}
	w.Balance += amount

	return nil
}

func main() {
	wallet := Wallet{
		ID: 1,
		Name: "Kaspi",
		Balance: 1000,
	}
	err := wallet.WalletDeposit(1250)
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	fmt.Println(wallet)

}
