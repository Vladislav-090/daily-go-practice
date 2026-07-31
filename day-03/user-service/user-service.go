package main

import (
	"errors"
	"fmt"
	"strings"
)

type UserRepository interface {
	GetUserByID(userID int64) (User, error)
	UpdateUser(user User) error
}

type MemoryUserRepository struct {
	users map[int64]User
}

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user id")
	ErrEmailIsEmpty  = errors.New("user email is empty")
	ErrInvalidEmail  = errors.New("invalid user email")
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
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *UserService) ChangeUserEmail(userID int64, email string) (User, error) {
	if userID <= 0 {
		return User{}, ErrInvalidUserID
	}

	if email == "" {
		return User{}, ErrEmailIsEmpty
	}

	if !strings.HasSuffix(email, "@gmail.com") && !strings.HasSuffix(email, "@mail.ru") {
		return User{}, ErrInvalidEmail
	}
	user, err := s.GetUserByID(userID)
	if err != nil {
		return User{}, err
	}

	user.Email = email

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func main() {
	repo := &MemoryUserRepository{
		users: map[int64]User{
			1: {
				ID:    1,
				Name:  "Vlad",
				Email: "vlad@gmail.com",
			},
			2: {
				ID:    2,
				Name:  "Viola",
				Email: "violatest@gmail.com",
			},
		},
	}

	service := NewUserService(repo)

	gotUser, err := service.GetUserByID(2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("got user:", gotUser)

	updatedUser, err := service.ChangeUserEmail(gotUser.ID, "newViola@gmail.com")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated user", updatedUser)
}
