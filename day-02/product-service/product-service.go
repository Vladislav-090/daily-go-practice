package main

import (
	"errors"
	"fmt"
)

type ProductRepository interface {
	GetProductByID(id int64) (Product, error)
	UpdateProduct(product Product) error
}

type MemoryProductRepository struct {
	Products map[int64]Product
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(productRepo ProductRepository) *ProductService {
	return &ProductService{
		repo: productRepo,
	}
}

type Product struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}

var (
	ErrInvalidProductID  = errors.New("invalid product id")
	ErrProductNotFound   = errors.New("product not found")
	ErrInvalidQuantity   = errors.New("invalid products quantity")
	ErrInvalidAmount     = errors.New("invalid amount to reserve")
	ErrInsufficientStock = errors.New("insufficient stock")
)

func (r *MemoryProductRepository) GetProductByID(id int64) (Product, error) {
	product, ok := r.Products[id]
	if !ok {
		return Product{}, ErrProductNotFound
	}

	return product, nil
}

func (s *ProductService) GetProductByID(productid int64) (Product, error) {
	if productid <= 0 {
		return Product{}, ErrInvalidProductID
	}
	product, err := s.repo.GetProductByID(productid)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (r *MemoryProductRepository) UpdateProduct(product Product) error {
	_, ok := r.Products[product.ID]
	if !ok {
		return ErrProductNotFound
	}
	r.Products[product.ID] = product

	return nil
}

func (s *ProductService) UpdateProduct(productID int64, quantity int64) (Product, error) {
	if productID <= 0 {
		return Product{}, ErrInvalidProductID
	}

	if quantity <= 0 {
		return Product{}, ErrInvalidQuantity
	}

	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return Product{}, err
	}
	product.Quantity = quantity

	err = s.repo.UpdateProduct(product)
	if err != nil {
		return Product{}, err
	}

	return product, nil

}

func (s *ProductService) ReserveProduct(productID int64, amount int64) (Product, error) {
	if productID <= 0 {
		return Product{}, ErrInvalidProductID
	}

	if amount <= 0 {
		return Product{}, ErrInvalidAmount
	}

	product, err := s.repo.GetProductByID(productID)
	if err != nil {
		return Product{}, err
	}

	if product.Quantity < amount {
		return Product{}, ErrInsufficientStock
	}

	product.Quantity -= amount

	err = s.repo.UpdateProduct(product)
	if err != nil {
		return Product{}, err
	}

	return product, nil

}
func main() {
	repo := &MemoryProductRepository{
		Products: map[int64]Product{
			1: {
				ID:       1,
				Name:     "Apple",
				Quantity: 5,
			},

			2: {
				ID:       2,
				Name:     "Laptop",
				Quantity: 6,
			},
		},
	}

	service := NewProductService(repo)
	gotProduct, err := service.GetProductByID(2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(gotProduct)

	updatedProduct, err := service.UpdateProduct(2, 10)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated product :", updatedProduct)

	reservedProduct, err := service.ReserveProduct(gotProduct.ID, 6)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("product in stock now:", reservedProduct)
}
