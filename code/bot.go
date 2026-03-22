package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type Bot struct {
	ID           int
	busy         int32
	CurrentOrder *Order
	stopChan     chan struct{}
	stopOnce     sync.Once
}

func NewBot(id int) *Bot {
	return &Bot{
		ID:       id,
		stopChan: make(chan struct{}),
	}
}

func (b *Bot) start() {
	go func() {
		for {
			select {
			case <-b.stopChan:
				log("Bot #%d stopped", b.ID)
				return
			default:
			}

			order := getNextOrder()
			if order == nil {
				time.Sleep(200 * time.Millisecond)
				continue
			}

			if !atomic.CompareAndSwapInt32(&b.busy, 0, 1) {
				returnOrderToQueue(order)
				continue
			}

			b.CurrentOrder = order

			log("Bot #%d picked up %s Order #%d - Status: PROCESSING", b.ID, order.Type, order.ID)

			processed := b.processOrder(order)

			atomic.StoreInt32(&b.busy, 0)
			b.CurrentOrder = nil

			if !processed {
				return
			}
		}
	}()
}

func (b *Bot) processOrder(order *Order) bool {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		mu.Lock()
		completedOrders++
		completeOrders = append(completeOrders, order)
		mu.Unlock()

		log("Bot #%d completed %s Order #%d - Status: COMPLETE (Processing time: 10s)",
			b.ID, order.Type, order.ID)
		return true

	case <-b.stopChan:
		returnOrderToQueue(order)
		log("Bot #%d destroyed while processing %s Order #%d - returning to queue",
			b.ID, order.Type, order.ID)
		return false
	}
}

func addBot() {
	mu.Lock()
	newID := len(bots) + 1
	bot := NewBot(newID)
	bots = append(bots, bot)
	mu.Unlock()

	bot.start()

	log("Bot #%d created - Status: ACTIVE", bot.ID)
}

func removeBot() {
	mu.Lock()
	if len(bots) == 0 {
		mu.Unlock()
		log("No bots to remove")
		return
	}

	bot := bots[len(bots)-1]
	bots = bots[:len(bots)-1]
	mu.Unlock()

	bot.stopOnce.Do(func() {
		close(bot.stopChan)
	})

	log("Bot #%d removed - Status: INACTIVE", bot.ID)
}
