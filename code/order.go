package main

import (
	"fmt"
)

type Order struct {
	ID   int
	Type string // "VIP" or "NORMAL"
}

func addNormalOrder() {
	order := NewOrder("NORMAL")
	normalQueue = append(normalQueue, order)
	log("%s", fmt.Sprintf("New NORMAL Order %d", order.ID))
}

func addVIPOrder() {
	order := NewOrder("VIP")
	vipQueue = append(vipQueue, order)
	log("%s", fmt.Sprintf("New VIP Order %d", order.ID))
}

func NewOrder(orderType string) Order {
	order := Order{ID: orderID, Type: orderType}
	orderID++
	return order
}
