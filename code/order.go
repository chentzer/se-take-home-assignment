package main

import (
	"errors"
	"time"
)

type Order struct {
	ID        int
	Type      string
	CreatedAt time.Time
}

const (
	OrderTypeNormal = "NORMAL"
	OrderTypeVIP    = "VIP"
)

func ValidateOrderType(orderType string) error {
	if orderType != OrderTypeNormal && orderType != OrderTypeVIP {
		return errors.New("invalid order type: " + orderType)
	}
	return nil
}

func NewOrder(orderType string) (*Order, error) {
	if err := ValidateOrderType(orderType); err != nil {
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()

	order := &Order{
		ID:        orderID,
		Type:      orderType,
		CreatedAt: time.Now(),
	}
	orderID++

	if orderType == OrderTypeVIP {
		totalVIP++
		vipQueue = append(vipQueue, order)
	} else {
		totalNormal++
		normalQueue = append(normalQueue, order)
	}

	return order, nil
}

func getNextOrder() *Order {
	mu.Lock()
	defer mu.Unlock()

	if len(vipQueue) > 0 {
		order := vipQueue[0]
		vipQueue = vipQueue[1:]
		return order
	}

	if len(normalQueue) > 0 {
		order := normalQueue[0]
		normalQueue = normalQueue[1:]
		return order
	}

	return nil
}

func returnOrderToQueue(order *Order) {
	mu.Lock()
	defer mu.Unlock()

	if order.Type == OrderTypeVIP {
		insertIndex := 0
		for i, existingOrder := range vipQueue {
			if existingOrder.CreatedAt.After(order.CreatedAt) {
				insertIndex = i
				break
			}
			insertIndex = i + 1
		}
		vipQueue = append(vipQueue[:insertIndex], append([]*Order{order}, vipQueue[insertIndex:]...)...)
	} else {
		insertIndex := 0
		for i, existingOrder := range normalQueue {
			if existingOrder.CreatedAt.After(order.CreatedAt) {
				insertIndex = i
				break
			}
			insertIndex = i + 1
		}
		normalQueue = append(normalQueue[:insertIndex], append([]*Order{order}, normalQueue[insertIndex:]...)...)
	}
}
