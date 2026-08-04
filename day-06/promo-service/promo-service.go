package main

import (
	"errors"
	"fmt"
)

type Promocode struct {
	Code          string `json:"code"`
	Discount      int    `json:"discount"`
	IsActive      bool   `json:"is_active"`
	RemainingUses int    `json:"remaining_uses"`
}

type MemoryPromoRepository struct {
	promocodes map[string]Promocode
}

var (
	ErrEmptyCode          = errors.New("the code is empty")
	ErrPromocodeNotFound  = errors.New("promocode not found")
	ErrPromocodeInactive  = errors.New("promocode is inactive")
	ErrPromoCodeExhausted = errors.New("promo code has no remaining uses")
)

func (r *MemoryPromoRepository) ApplyPromo(code string) (Promocode, error) {
	if code == "" {
		return Promocode{}, ErrEmptyCode
	}

	promocode, ok := r.promocodes[code]
	if !ok {
		return Promocode{}, ErrPromocodeNotFound
	}

	if !promocode.IsActive {
		return Promocode{}, ErrPromocodeInactive
	}

	if promocode.RemainingUses < 1 {
		return Promocode{}, ErrPromoCodeExhausted
	}
	promocode.RemainingUses -= 1

	r.promocodes[code] = promocode

	return promocode, nil
}

func main() {
	repo := &MemoryPromoRepository{
		promocodes: map[string]Promocode{
			"code1": {
				Code:          "code1",
				Discount:      15,
				IsActive:      true,
				RemainingUses: 1,
			},

			"code2": {
				Code:          "code2",
				Discount:      10,
				IsActive:      false,
				RemainingUses: 0,
			},

			"code3": {
				Code:          "code3",
				Discount:      5,
				IsActive:      true,
				RemainingUses: 5,
			},

			"code4": {
				Code:          "code4",
				Discount:      5,
				IsActive:      true,
				RemainingUses: 0,
			},
		},
	}

	result1, err := repo.ApplyPromo("code1")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(result1)

	result2, err := repo.ApplyPromo("code3")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(result2)

	result3, err := repo.ApplyPromo("code4")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(result3)
}
