package main

import (
	"errors"
	"fmt"
)

type UserRepository interface {
	GetUserByID(userID int64) (User, error)
	UpdateUser(user User) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		repo: userRepo,
	}
}

type MemoryUserRepository struct {
	users map[int64]User
}

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	IsBlocked bool   `json:"is_blocked"`
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrUserAlreadyBlocked = errors.New("user already blocked")
)

func (r *MemoryUserRepository) GetUserByID(userID int64) (User, error) {
	user, ok := r.users[userID]
	if !ok {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (r *MemoryUserRepository) UpdateUser(user User) error {
	_, ok := r.users[user.ID]
	if !ok {
		return ErrUserNotFound
	}
	r.users[user.ID] = user

	return nil
}

func (s *UserService) GetUserByID(userID int64) (User, error) {
	if userID <= 0 {
		return User{}, ErrInvalidUserID
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserService) BlockUser(userID int64) (User, error) {
	if userID <= 0 {
		return User{}, ErrInvalidUserID
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return User{}, err
	}

	if user.IsBlocked == true {
		return User{}, ErrUserAlreadyBlocked
	}

	user.IsBlocked = true

	err = s.repo.UpdateUser(user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func main() {
	repo := &MemoryUserRepository{
		users: map[int64]User{
			1: {
				ID:        1,
				Name:      "VLad",
				Email:     "vlad@gmail.com",
				IsBlocked: false,
			},
			2: {
				ID:        2,
				Name:      "Viola",
				Email:     "viola@gmail.com",
				IsBlocked: true,
			},
		},
	}

	service := NewUserService(repo)

	gotUser1, err := service.GetUserByID(1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("got user:", gotUser1)

	updatedUser1, err := service.BlockUser(gotUser1.ID)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("user hasbeen blocked", updatedUser1)

	blockedUser, err := service.GetUserByID(gotUser1.ID)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("blocked user:", blockedUser)

	updatedUser2, err := service.BlockUser(2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("user hasbeen blocked", updatedUser2)
}
