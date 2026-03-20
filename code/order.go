package main

import (
	"fmt"
)

type Order struct {
	ID   int
	Type string // "VIP" or "NORMAL"
}

func addNormalOrder() {
	order := Order{ID: orderID, Type: "NORMAL"}
	orderID++
	normalQueue = append(normalQueue, order)
	log("%s", fmt.Sprintf("New NORMAL Order %d", order.ID))
}

func addVIPOrder() {
	order := Order{ID: orderID, Type: "VIP"}
	orderID++
	vipQueue = append(vipQueue, order)
	log("%s", fmt.Sprintf("New VIP Order %d", order.ID))
}
