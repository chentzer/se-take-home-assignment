package main

import (
	"fmt"
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

func (b *Bot) PickOrder() {
	if len(pendingQueue) > 0 {
		order := pendingQueue[0]
		pendingQueue = pendingQueue[1:]
		b.CurrentOrder = order
		b.busy = true
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
			order := getNextOrder()
			if order == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			b.busy = true
			b.CurrentOrder = order
			log("%s", fmt.Sprintf("Bot %d processing Order %d", b.ID, order.ID))

			select {
			case <-time.After(10 * time.Second):
				completeOrders = append(completeOrders, *order)
				log("%s", fmt.Sprintf("Order %d, Type: %s - COMPLETE", order.ID, order.Type))
				b.busy = false
				b.CurrentOrder = nil
				log("%s", fmt.Sprintf("Bot %d stopped", b.ID))

				// return order back to queue
				if order.Type == "VIP" {
					vipQueue = append([]Order{*order}, vipQueue...)
				} else {
					normalQueue = append([]Order{*order}, normalQueue...)
				}
				return
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
}

func removeBot() {
	if len(bots) == 0 {
		return
	}

	bot := bots[len(bots)-1]
	bots = bots[:len(bots)-1]

	if bot.busy {
		bot.stopChan <- true
	}
}
