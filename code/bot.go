package code

import (
	"sync"
	"sync/atomic"
	"time"
)

// Bot represents a cooking bot
type Bot struct {
	ID           int
	busy         int32
	currentOrder *Order
	orderMu      sync.Mutex
	stopChan     chan struct{}
	StopOnce     sync.Once
	controller   *Controller
}

// NewBot creates a new bot
func NewBot(id int, controller *Controller) *Bot {
	return &Bot{
		ID:         id,
		stopChan:   make(chan struct{}),
		controller: controller,
	}
}

// GetCurrentOrder safely returns the current order being processed
func (b *Bot) GetCurrentOrder() *Order {
	b.orderMu.Lock()
	defer b.orderMu.Unlock()
	return b.currentOrder
}

// SetCurrentOrder safely sets the current order
func (b *Bot) SetCurrentOrder(order *Order) {
	b.orderMu.Lock()
	defer b.orderMu.Unlock()
	b.currentOrder = order
}

// IsBusy returns whether the bot is currently processing an order
func (b *Bot) IsBusy() bool {
	return atomic.LoadInt32(&b.busy) == 1
}

// Start begins the bot's processing loop
func (b *Bot) Start() {
	go func() {
		for {
			select {
			case <-b.stopChan:
				b.controller.Log("Bot #%d stopped", b.ID)
				return
			default:
			}

			order := b.controller.GetNextOrder()
			if order == nil {
				time.Sleep(200 * time.Millisecond)
				continue
			}

			if !atomic.CompareAndSwapInt32(&b.busy, 0, 1) {
				b.controller.ReturnOrderToQueue(order)
				continue
			}

			b.SetCurrentOrder(order)

			b.controller.Log("Bot #%d picked up %s Order #%d - Status: PROCESSING", b.ID, order.Type, order.ID)

			processed := b.processOrder(order)

			atomic.StoreInt32(&b.busy, 0)
			b.SetCurrentOrder(nil)

			if !processed {
				return
			}
		}
	}()
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.StopOnce.Do(func() {
		close(b.stopChan)
	})
}

func (b *Bot) processOrder(order *Order) bool {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		b.controller.CompleteOrder(order)
		b.controller.Log("Bot #%d completed %s Order #%d - Status: COMPLETE (Processing time: 10s)",
			b.ID, order.Type, order.ID)
		return true

	case <-b.stopChan:
		b.controller.ReturnOrderToQueue(order)
		b.controller.Log("Bot #%d destroyed while processing %s Order #%d - returning to queue",
			b.ID, order.Type, order.ID)
		return false
	}
}

// AddBot creates and starts a new bot
func (c *Controller) AddBot() *Bot {
	c.mu.Lock()
	newID := c.nextBotID
	c.nextBotID++
	bot := NewBot(newID, c)
	c.Bots = append(c.Bots, bot)
	c.mu.Unlock()

	bot.Start()

	c.Log("Bot #%d created - Status: ACTIVE", bot.ID)
	return bot
}

// RemoveBot removes and stops the newest bot
func (c *Controller) RemoveBot() {
	c.mu.Lock()
	if len(c.Bots) == 0 {
		c.mu.Unlock()
		c.Log("No bots to remove")
		return
	}

	bot := c.Bots[len(c.Bots)-1]
	c.Bots = c.Bots[:len(c.Bots)-1]
	c.mu.Unlock()

	bot.Stop()

	c.Log("Bot #%d removed - Status: INACTIVE", bot.ID)
}

// StopAllBots stops all bots gracefully
func (c *Controller) StopAllBots() {
	c.mu.Lock()
	botsToStop := make([]*Bot, len(c.Bots))
	copy(botsToStop, c.Bots)
	c.Bots = []*Bot{}
	c.mu.Unlock()

	for _, bot := range botsToStop {
		bot.Stop()
	}
}
