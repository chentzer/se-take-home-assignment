package main

import (
	"time"
)

type Bot struct {
	ID           int
	busy         bool
	CurrentOrder *Order
	stopChan     chan bool
}

func NewBot() *Bot {
	return &Bot{
		ID:       len(bots) + 1,
		stopChan: make(chan bool),
	}
}

func getNextOrder() *Order {
	if len(vipQueue) > 0 {
		order := vipQueue[0]
		vipQueue = vipQueue[1:]
		return &order
	}
	if len(normalQueue) > 0 {
		order := normalQueue[0]
		normalQueue = normalQueue[1:]
		return &order
	}
	return nil
}

func (b *Bot) start() {
	go func() {
		for {
			select {
			case <-b.stopChan:
				// Bot is stopped
				if b.busy && b.CurrentOrder != nil {
					order := b.CurrentOrder

					// Return order to front of queue
					if order.Type == "VIP" {
						vipQueue = append([]Order{*order}, vipQueue...)
					} else {
						normalQueue = append([]Order{*order}, normalQueue...)
					}

					log("Bot #%d stopped while processing Order #%d", b.ID, order.ID)
				} else {
					log("Bot #%d destroyed while IDLE", b.ID)
				}
				return

			default:
				order := getNextOrder()

				if order == nil {
					if b.busy {
						b.busy = false
						log("Bot #%d is now IDLE - No pending orders", b.ID)
					}
					time.Sleep(1 * time.Second)
					continue
				}

				// Start processing
				b.busy = true
				b.CurrentOrder = order

				log("Bot #%d picked up %s Order #%d - Status: PROCESSING",
					b.ID, order.Type, order.ID)

				time.Sleep(10 * time.Second)

				// Complete order
				completeOrders = append(completeOrders, *order)
				completedOrders++

				// Update counters
				if order.Type == "VIP" {
					totalVIP++
				} else {
					totalNormal++
				}

				log("Bot #%d completed %s Order #%d - Status: COMPLETE (Processing time: 10s)",
					b.ID, order.Type, order.ID)

				b.CurrentOrder = nil
			}
		}
	}()
}

func addBot() {
	bot := &Bot{
		ID:       len(bots) + 1,
		stopChan: make(chan bool),
	}
	bots = append(bots, bot)
	bot.start()
	log("Bot #%d created - Status: ACTIVE", bot.ID)
}

func removeBot() {
	if len(bots) == 0 {
		return
	}

	bot := bots[len(bots)-1]
	bots = bots[:len(bots)-1]

	if bot.busy {
		close(bot.stopChan)
	}
	log("Bot #%d removed - Status: INACTIVE", bot.ID)
}
