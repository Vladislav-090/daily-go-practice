package main

import "fmt"

type Order struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Status string `json:"status"`
}

func GroupOrdersByStatus(orders []Order) map[string][]Order {
	grouped := make(map[string][]Order)

	for _, order := range orders {
		grouped[order.Status] = append(grouped[order.Status], order)
	}
	return grouped
}

func main() {
	orders := []Order{
		{
			ID:     1,
			UserID: 1,
			Status: "new",
		},

		{
			ID:     2,
			UserID: 1,
			Status: "paid",
		},
		{
			ID:     3,
			UserID: 1,
			Status: "new",
		},

		{
			ID:     4,
			UserID: 2,
			Status: "cancelled",
		},

		{
			ID:     5,
			UserID: 2,
			Status: "paid",
		},

		{
			ID:     6,
			UserID: 1,
			Status: "cancelled",
		},
	}

	groupedOrders := GroupOrdersByStatus(orders)

	for status, statusOrders := range groupedOrders {
		fmt.Println("status:", status)
		fmt.Println("orders:", statusOrders)
	}
}
