package main

import (
	"errors"
	"fmt"
)

type UserRepository interface {
	GetUserByID(userID int64) (User, error)
	UpdateUser(user User) error
}

type Device struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type User struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Devices []Device `json:"devices"`
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

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidDeviceID     = errors.New("invalid device id")
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDeviceAlreadyActive = errors.New("device already active")
	ErrActiveDeviceLimit   = errors.New("active device limit")
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

func (s *UserService) ActivateDevice(userID int64, deviceID int64) (User, error) {
	if userID <= 0 {
		return User{}, ErrInvalidUserID
	}

	if deviceID <= 0 {
		return User{}, ErrInvalidDeviceID
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return User{}, err
	}

	deviceIndex := -1

	for i := range user.Devices {
		if user.Devices[i].ID == deviceID {
			deviceIndex = i
			break
		}
	}
	if deviceIndex == -1 {
		return User{}, ErrDeviceNotFound
	}
	if user.Devices[deviceIndex].IsActive {
		return User{}, ErrDeviceAlreadyActive
	}

	activeCount := 0

	for i := range user.Devices {
		if user.Devices[i].IsActive {
			activeCount++
		}
	}
	if activeCount >= 2 {
		return User{}, ErrActiveDeviceLimit
	}

	user.Devices[deviceIndex].IsActive = true

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
				ID:   1,
				Name: "Vlad",
				Devices: []Device{
					{
						ID:       1,
						Name:     "phone",
						IsActive: false,
					},

					{
						ID:       2,
						Name:     "laptop",
						IsActive: false,
					},

					{
						ID:       3,
						Name:     "Tv",
						IsActive: false,
					},
				},
			},
		},
	}

	service := NewUserService(repo)

	gotUser, err := service.repo.GetUserByID(1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("got user:", gotUser)

	updatedUser1, err := service.ActivateDevice(1, 2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated user:", updatedUser1)

	updatedUser2, err := service.ActivateDevice(1, 3)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated user ", updatedUser2)

	updatedUser3, err := service.ActivateDevice(1, 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated user:", updatedUser3)

}
