package main

type Order struct {
	ID   int
	Type string // "VIP" or "NORMAL"
}

func addNormalOrder() {
	order := NewOrder("NORMAL")
	normalQueue = append(normalQueue, order)
	log("Created Normal Order #%d - Status: PENDING", order.ID)
}

func addVIPOrder() {
	order := NewOrder("VIP")
	vipQueue = append(vipQueue, order)
	log("Created VIP Order #%d - Status: PENDING", order.ID)
}

func NewOrder(orderType string) Order {
	order := Order{ID: orderID, Type: orderType}
	orderID++
	return order
}
