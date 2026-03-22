package code

import (
	"errors"
	"sync"
	"time"
)

// Order represents a customer order
type Order struct {
	ID        int
	Type      string
	CreatedAt time.Time
}

const (
	OrderTypeNormal = "NORMAL"
	OrderTypeVIP    = "VIP"
)

// ValidateOrderType checks if the order type is valid
func ValidateOrderType(orderType string) error {
	if orderType != OrderTypeNormal && orderType != OrderTypeVIP {
		return errors.New("invalid order type: " + orderType)
	}
	return nil
}

// Controller manages the order system state
type Controller struct {
	mu sync.Mutex

	vipQueue       []*Order
	normalQueue    []*Order
	CompleteOrders []*Order

	Bots []*Bot

	orderID         int
	nextBotID       int // Monotonically increasing bot ID
	TotalVIP        int
	TotalNormal     int
	CompletedOrders int

	LogFunc func(format string, args ...interface{})
}

// NewController creates a new order controller
func NewController(logFunc func(format string, args ...interface{})) *Controller {
	return &Controller{
		vipQueue:       []*Order{},
		normalQueue:    []*Order{},
		CompleteOrders: []*Order{},
		Bots:           []*Bot{},
		orderID:        1,
		nextBotID:      1,
		LogFunc:        logFunc,
	}
}

// NewOrder creates a new order and adds it to the appropriate queue
func (c *Controller) NewOrder(orderType string) (*Order, error) {
	if err := ValidateOrderType(orderType); err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	order := &Order{
		ID:        c.orderID,
		Type:      orderType,
		CreatedAt: time.Now(),
	}
	c.orderID++

	if orderType == OrderTypeVIP {
		c.TotalVIP++
		c.vipQueue = append(c.vipQueue, order)
	} else {
		c.TotalNormal++
		c.normalQueue = append(c.normalQueue, order)
	}

	return order, nil
}

// GetNextOrder returns the next order to process (VIP first, then normal)
func (c *Controller) GetNextOrder() *Order {
	c.mu.Lock()
	defer c.mu.Unlock()

	// VIP orders have priority
	if len(c.vipQueue) > 0 {
		order := c.vipQueue[0]
		c.vipQueue = c.vipQueue[1:]
		return order
	}

	// Then normal orders
	if len(c.normalQueue) > 0 {
		order := c.normalQueue[0]
		c.normalQueue = c.normalQueue[1:]
		return order
	}

	return nil
}

// ReturnOrderToQueue returns an order to its original position maintaining priority and FIFO
func (c *Controller) ReturnOrderToQueue(order *Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if order.Type == OrderTypeVIP {
		insertIndex := 0
		for i, existingOrder := range c.vipQueue {
			if existingOrder.CreatedAt.After(order.CreatedAt) {
				insertIndex = i
				break
			}
			insertIndex = i + 1
		}
		c.vipQueue = append(c.vipQueue[:insertIndex], append([]*Order{order}, c.vipQueue[insertIndex:]...)...)
	} else {
		insertIndex := 0
		for i, existingOrder := range c.normalQueue {
			if existingOrder.CreatedAt.After(order.CreatedAt) {
				insertIndex = i
				break
			}
			insertIndex = i + 1
		}
		c.normalQueue = append(c.normalQueue[:insertIndex], append([]*Order{order}, c.normalQueue[insertIndex:]...)...)
	}
}

// CompleteOrder marks an order as completed
func (c *Controller) CompleteOrder(order *Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CompletedOrders++
	c.CompleteOrders = append(c.CompleteOrders, order)
}

// GetPendingCount returns the number of pending orders
func (c *Controller) GetPendingCount() (vip, normal int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.vipQueue), len(c.normalQueue)
}

// GetStats returns current system statistics
func (c *Controller) GetStats() (totalVIP, totalNormal, completed, pendingVIP, pendingNormal, activeBots int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.TotalVIP, c.TotalNormal, c.CompletedOrders, len(c.vipQueue), len(c.normalQueue), len(c.Bots)
}

// Log calls the log function if set
func (c *Controller) Log(format string, args ...interface{}) {
	if c.LogFunc != nil {
		c.LogFunc(format, args...)
	}
}
