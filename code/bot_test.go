package main

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestAddBot(t *testing.T) {
	resetTestState()

	initialBotCount := len(bots)
	addBot()

	time.Sleep(100 * time.Millisecond)

	if len(bots) != initialBotCount+1 {
		t.Errorf("Expected %d bots, got %d", initialBotCount+1, len(bots))
	}
}

func TestRemoveBot(t *testing.T) {
	resetTestState()

	addBot()
	addBot()
	time.Sleep(100 * time.Millisecond)

	initialBotCount := len(bots)
	removeBot()

	time.Sleep(100 * time.Millisecond)

	if len(bots) != initialBotCount-1 {
		t.Errorf("Expected %d bots, got %d", initialBotCount-1, len(bots))
	}
}

func TestRemoveBotWhenEmpty(t *testing.T) {
	resetTestState()

	initialBotCount := len(bots)
	removeBot()

	if len(bots) != initialBotCount {
		t.Errorf("Bot count should remain %d, got %d", initialBotCount, len(bots))
	}
}

func TestBotProcessesOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long test in short mode")
	}

	resetTestState()

	addBot()
	order, err := NewOrder("NORMAL")
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	hasOrder := bots[0].CurrentOrder != nil
	mu.Unlock()

	if !hasOrder {
		t.Error("Bot should have an order")
	}

	time.Sleep(11 * time.Second)

	mu.Lock()
	if completedOrders != 1 {
		t.Errorf("Expected 1 completed order, got %d", completedOrders)
	}
	if completeOrders[0].ID != order.ID {
		t.Errorf("Expected order ID %d, got %d", order.ID, completeOrders[0].ID)
	}
	mu.Unlock()
}

func TestBotReturnsOrderWhenRemoved(t *testing.T) {
	resetTestState()

	addBot()
	order, err := NewOrder("NORMAL")
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	queueLen := len(normalQueue)
	hasOrder := bots[0].CurrentOrder != nil
	mu.Unlock()

	if queueLen != 0 {
		t.Fatal("Order not taken from queue")
	}
	if !hasOrder {
		t.Fatal("Bot does not have order")
	}

	removeBot()
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	if len(normalQueue) == 0 {
		t.Error("Order was not returned to queue")
	}
	if len(normalQueue) > 0 && normalQueue[0].ID != order.ID {
		t.Errorf("Expected order %d at front, got %d", order.ID, normalQueue[0].ID)
	}
	mu.Unlock()
}

func TestBotDoesNotTakeOrderWhenBusy(t *testing.T) {
	resetTestState()

	addBot()
	NewOrder("NORMAL")

	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	isBusy := atomic.LoadInt32(&bots[0].busy) == 1
	mu.Unlock()

	if !isBusy {
		t.Fatal("Bot is not busy")
	}

	NewOrder("NORMAL")
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	if len(normalQueue) != 1 {
		t.Errorf("Expected 1 order in queue, got %d", len(normalQueue))
	}
	mu.Unlock()
}

func TestVIPOrderPriority(t *testing.T) {
	resetTestState()

	NewOrder("NORMAL")
	NewOrder("VIP")
	addBot()

	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	bot := bots[0]
	if bot.CurrentOrder != nil && bot.CurrentOrder.Type != "VIP" {
		t.Errorf("Bot should pick VIP first, got %s", bot.CurrentOrder.Type)
	}
	mu.Unlock()
}
