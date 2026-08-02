package main

import (
	"errors"
	"fmt"
)

type CartRepository interface {
	GetCartByUserID(userid int64) (Cart, error)
	UpdateCart(cart Cart) error
}

type MemoryCartRepository struct {
	carts map[int64]Cart
}

type Cart struct {
	UserID int64      `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type CartService struct {
	repo CartRepository
}

func NewCartService(cartRepo CartRepository) *CartService {
	return &CartService{
		repo: cartRepo,
	}
}

var (
	ErrCartNotFound      = errors.New("cart not found")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidProductID  = errors.New("invalid product id")
	ErrQuantityOfProduct = errors.New("quantity of products must be positive")
)

func (r *MemoryCartRepository) GetCartByUserID(userID int64) (Cart, error) {
	cart, ok := r.carts[userID]
	if !ok {
		return Cart{}, ErrCartNotFound
	}

	return cart, nil
}

func (s *CartService) GetCartByUserID(userID int64) (Cart, error) {
	if userID <= 0 {
		return Cart{}, ErrInvalidUserID
	}

	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return Cart{}, err
	}

	return cart, nil
}

func (r *MemoryCartRepository) UpdateCart(cart Cart) error {
	_, ok := r.carts[cart.UserID]
	if !ok {
		return ErrCartNotFound
	}

	r.carts[cart.UserID] = cart

	return nil
}

func (s *CartService) AddProduct(userID int64, productID int64, quantity int64) (Cart, error) {
	if userID <= 0 {
		return Cart{}, ErrInvalidUserID
	}

	if productID <= 0 {
		return Cart{}, ErrInvalidProductID
	}

	if quantity < 1 {
		return Cart{}, ErrQuantityOfProduct
	}

	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return Cart{}, err
	}
	found := false
	for i := range cart.Items {
		if cart.Items[i].ProductID == productID {
			cart.Items[i].Quantity += quantity
			found = true
			break
		}
	}

	if found == false {
		newCartItem := CartItem{
			ProductID: productID,
			Quantity:  quantity,
		}
		cart.Items = append(cart.Items, newCartItem)
	}

	err = s.repo.UpdateCart(cart)
	if err != nil {
		return Cart{}, err
	}

	return cart, nil
}

func main() {
	repo := &MemoryCartRepository{
		carts: map[int64]Cart{
			1: {
				UserID: 1,
				Items: []CartItem{
					{
						ProductID: 1,
						Quantity:  4,
					},
					{
						ProductID: 2,
						Quantity:  1,
					},
				},
			},
			2: {
				UserID: 2,
				Items: []CartItem{
					{
						ProductID: 1,
						Quantity:  10,
					},
				},
			},
		},
	}

	service := NewCartService(repo)

	gotCart, err := service.GetCartByUserID(1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("user's cart:", gotCart)

	updatedCart1, err := service.AddProduct(gotCart.UserID, 3, 6)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated cart:", updatedCart1)

	updatedCart2, err := service.AddProduct(gotCart.UserID, 2, 6)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated cart:", updatedCart2)
}
