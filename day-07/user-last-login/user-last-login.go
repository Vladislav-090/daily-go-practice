package main

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	LastLogin time.Time `json:"last_login"`
	IsBlocked bool      `json:"is_blocked"`
}

var (
	ErrInvalidDays = errors.New("day must be positive")
)

func GetRecentlyActiveUsers(users []User, days int) ([]User, error) {
	if days <= 0 {
		return []User{}, ErrInvalidDays
	}
	lastLogin := time.Now().AddDate(0, 0, -days)

	result := make([]User, 0)

	for _, user := range users {
		if user.IsBlocked == true {
			continue
		}
		if user.LastLogin.Before(lastLogin) {
			continue
		}
		result = append(result, user)
	}

	return result, nil
}

func main() {
	users := []User{
		{
			ID:        1,
			Name:      "vlad",
			LastLogin: time.Now().AddDate(0, 0, -3),
			IsBlocked: false,
		},

		{
			ID:        2,
			Name:      "viola",
			LastLogin: time.Now().AddDate(0, 0, -5),
			IsBlocked: false,
		},

		{
			ID:        3,
			Name:      "afina",
			LastLogin: time.Now(),
			IsBlocked: false,
		},
	}

	result, err := GetRecentlyActiveUsers(users, 4)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, user := range result {
		fmt.Println("user_id", user.ID)
		fmt.Println("user_name", user.Name)
		fmt.Println("user_last_login", user.LastLogin.Format("02.01.2006 15:05:00"))
		fmt.Println("user_id", user.IsBlocked)
		fmt.Println("------------------------")

	}

}
