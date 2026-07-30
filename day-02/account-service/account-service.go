package main

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

type AccountRepository interface {
	GetAccountByID(accountID int64) (Account, error)
	UpdateAccount(account Account) error
}

type Account struct {
	ID      int64           `json:"id"`
	Owner   string          `json:"owner"`
	Balance decimal.Decimal `json:"balance"`
}

type AccountService struct {
	repo AccountRepository
}

type MemoryAccountRepository struct {
	Accounts map[int64]Account
}

func NewAccountService(accountRepo AccountRepository) *AccountService {
	return &AccountService{
		repo: accountRepo,
	}
}

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInvalidAccountID  = errors.New("invalid account id")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

func (r *MemoryAccountRepository) GetAccountByID(accountID int64) (Account, error) {
	account, ok := r.Accounts[accountID]
	if !ok {
		return Account{}, ErrAccountNotFound
	}

	return account, nil
}

func (s *AccountService) GetAccountByID(accountID int64) (Account, error) {
	if accountID <= 0 {
		return Account{}, ErrInvalidAccountID
	}

	account, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

func (r *MemoryAccountRepository) UpdateAccount(account Account) error {
	_, ok := r.Accounts[account.ID]
	if !ok {
		return ErrAccountNotFound
	}

	r.Accounts[account.ID] = account

	return nil
}

func (s *AccountService) Deposit(accountID int64, amount decimal.Decimal) (Account, error) {
	if accountID <= 0 {
		return Account{}, ErrInvalidAccountID
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return Account{}, ErrInvalidAmount
	}

	account, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		return Account{}, err
	}

	account.Balance = account.Balance.Add(amount)

	err = s.repo.UpdateAccount(account)
	if err != nil {
		return Account{}, err
	}

	return account, nil

}

func (s *AccountService) Withdraw(accountID int64, amount decimal.Decimal) (Account, error) {
	if accountID <= 0 {
		return Account{}, ErrInvalidAccountID
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return Account{}, ErrInvalidAmount
	}

	account, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		return Account{}, err
	}

	if account.Balance.LessThan(amount) {
		return Account{}, ErrInsufficientFunds
	}

	account.Balance = account.Balance.Sub(amount)

	err = s.repo.UpdateAccount(account)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

func main() {
	repo := &MemoryAccountRepository{
		Accounts: map[int64]Account{
			1: {
				ID:      1,
				Owner:   "Vlad",
				Balance: decimal.RequireFromString("1000.00"),
			},
			2: {
				ID:      2,
				Owner:   "Viola",
				Balance: decimal.RequireFromString("20000.00"),
			},
		},
	}

	service := NewAccountService(repo)

	account, err := service.GetAccountByID(1)
	if err != nil {
		fmt.Println("error :", err)
		return
	}

	fmt.Println("account: ", account)
	fmt.Println("balance before:", account.Balance)

	updatedAccount, err := service.Deposit(account.ID, decimal.RequireFromString("25000.00"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("after deposit:", updatedAccount)

	updatedAccount, err = service.Withdraw(account.ID, decimal.RequireFromString("10000.00"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("after withdraw:", updatedAccount)

}
