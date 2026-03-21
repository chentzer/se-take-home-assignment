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
				log("Bot #%d stopped", b.ID)
				return
			default:
				order := getNextOrder()
				if order == nil {
					time.Sleep(1 * time.Second)
					continue
				}

				b.busy = true
				b.CurrentOrder = order
				log("Bot #%d picked up %s Order #%d - Status: PROCESSING", b.ID, order.Type, order.ID)

				select {
				case <-time.After(10 * time.Second):
					// Only complete if stopChan is not triggered
					completeOrders = append(completeOrders, *order)
					log("Bot #%d completed %s Order #%d - Status: COMPLETE (Processing time: 10s)", b.ID, order.Type, order.ID)
					b.busy = false
					b.CurrentOrder = nil
				case <-b.stopChan:
					// Immediately return order to queue
					log("Bot #%d destroyed while processing %s Order #%d - returning order to queue", b.ID, order.Type, order.ID)
					if order.Type == "VIP" {
						vipQueue = append([]Order{*order}, vipQueue...)
					} else {
						normalQueue = append([]Order{*order}, normalQueue...)
					}
					b.busy = false
					b.CurrentOrder = nil
					return
				}
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
		log("No bots to remove")
		return
	}

	// Remove the newest bot
	bot := bots[len(bots)-1]
	bots = bots[:len(bots)-1]

	if bot.busy && bot.CurrentOrder != nil {
		log("Bot #%d removed while processing %s Order #%d", bot.ID, bot.CurrentOrder.Type, bot.CurrentOrder.ID)
		// Put order back to the front of the appropriate queue
		if bot.CurrentOrder.Type == "VIP" {
			vipQueue = append([]Order{*bot.CurrentOrder}, vipQueue...)
		} else {
			normalQueue = append([]Order{*bot.CurrentOrder}, normalQueue...)
		}
	} else {
		log("Bot #%d removed while IDLE", bot.ID)
	}

	// Signal bot to stop
	close(bot.stopChan)
}
