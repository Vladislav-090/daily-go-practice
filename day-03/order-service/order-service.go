package main

import (
	"errors"
	"fmt"
)

type OrderRepository interface {
	GetOrderByID(orderID int64) (Order, error)
	UpdateOrder(order Order) error
}

type Order struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Status string `json:"status"`
}

type OrderService struct {
	orderRepo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: repo,
	}
}

type MemoryOrderRepository struct {
	Orders map[int64]Order
}

var (
	ErrInvalidOrderID          = errors.New("invalid order id")
	ErrInvalidOrderStatus      = errors.New("invalid order status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrOrderNotFound           = errors.New("order not found")
)

var allowedStatus = map[string]bool{
	"new":       true,
	"paid":      true,
	"cancelled": true,
}

func (r *MemoryOrderRepository) GetOrderByID(orderID int64) (Order, error) {
	order, ok := r.Orders[orderID]
	if !ok {
		return Order{}, ErrOrderNotFound
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(orderID int64) (Order, error) {
	if orderID <= 0 {
		return Order{}, ErrInvalidOrderID
	}

	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (r *MemoryOrderRepository) UpdateOrder(order Order) error {
	_, ok := r.Orders[order.ID]
	if !ok {
		return ErrOrderNotFound
	}

	r.Orders[order.ID] = order

	return nil
}

func (s *OrderService) ChangeOrderStatus(orderID int64, newStatus string) (Order, error) {
	if orderID <= 0 {
		return Order{}, ErrInvalidOrderID
	}

	if newStatus == "" {
		return Order{}, ErrInvalidOrderStatus
	}

	_, ok := allowedStatus[newStatus]
	if !ok {
		return Order{}, ErrInvalidOrderStatus
	}

	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return Order{}, err
	}

	if order.Status == newStatus {
		return Order{}, ErrInvalidStatusTransition
	}

	if order.Status == "cancelled" && newStatus == "new" {
		return Order{}, ErrInvalidStatusTransition
	}

	order.ID = orderID
	order.Status = newStatus

	err = s.orderRepo.UpdateOrder(order)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func main() {
	repo := &MemoryOrderRepository{
		Orders: map[int64]Order{
			1: {
				ID:     1,
				UserID: 1,
				Status: "new",
			},
			2: {
				ID:     2,
				UserID: 1,
				Status: "new",
			},
			3: {
				ID:     3,
				UserID: 1,
				Status: "cancelled",
			},
		},
	}

	service := NewOrderService(repo)

	gotOrder1, err := service.GetOrderByID(2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("order:", gotOrder1)

	gotOrder2, err := service.GetOrderByID(1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("order:", gotOrder2)

	updatedOrder1, err := service.ChangeOrderStatus(gotOrder1.ID, "cancelled")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated order:", updatedOrder1)

	updateOrder2, err := service.ChangeOrderStatus(gotOrder2.ID, "paid")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated order:", updateOrder2)

	updateOrder3, err := service.ChangeOrderStatus(gotOrder2.ID, "in the kitchen")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("updated order:", updateOrder3)

}
