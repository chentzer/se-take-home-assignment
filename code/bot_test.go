package main

import "testing"

// Test that bot picks up order from queue
func TestBotPickOrder(t *testing.T) {
	pendingQueue = []*Order{} // reset queue
	o := NewOrder("NORMAL")
	pendingQueue = append(pendingQueue, &o)

	bot := NewBot()
	bot.PickOrder()

	if bot.CurrentOrder == nil {
		t.Errorf("Bot did not pick an order")
	}
	if bot.CurrentOrder.ID != o.ID {
		t.Errorf("Bot picked wrong order")
	}
}
