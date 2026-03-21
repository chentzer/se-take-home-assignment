package main

import "testing"

// Test that bot picks up order from VIP/Normal queue correctly
func TestBotPickOrder(t *testing.T) {
	// Reset queues
	vipQueue = []Order{}
	normalQueue = []Order{}

	// Add a normal order
	o := NewOrder("NORMAL")
	normalQueue = append(normalQueue, o)

	bot := NewBot()

	// Simulate picking order
	order := getNextOrder()
	bot.CurrentOrder = order

	if bot.CurrentOrder == nil {
		t.Errorf("Bot did not pick an order")
	}
	if bot.CurrentOrder.ID != o.ID {
		t.Errorf("Bot picked wrong order")
	}
}
