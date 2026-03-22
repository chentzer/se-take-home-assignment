package code

import (
	"testing"
	"time"
)

func TestAddBot(t *testing.T) {
	c := NewTestController()

	initialBotCount := len(c.Bots)
	c.AddBot()

	time.Sleep(100 * time.Millisecond)

	if len(c.Bots) != initialBotCount+1 {
		t.Errorf("Expected %d bots, got %d", initialBotCount+1, len(c.Bots))
	}

	c.StopAllBots()
}

func TestRemoveBot(t *testing.T) {
	c := NewTestController()

	c.AddBot()
	c.AddBot()
	time.Sleep(100 * time.Millisecond)

	initialBotCount := len(c.Bots)
	c.RemoveBot()

	time.Sleep(100 * time.Millisecond)

	if len(c.Bots) != initialBotCount-1 {
		t.Errorf("Expected %d bots, got %d", initialBotCount-1, len(c.Bots))
	}

	c.StopAllBots()
}

func TestRemoveBotWhenEmpty(t *testing.T) {
	c := NewTestController()

	initialBotCount := len(c.Bots)
	c.RemoveBot()

	if len(c.Bots) != initialBotCount {
		t.Errorf("Bot count should remain %d, got %d", initialBotCount, len(c.Bots))
	}
}

func TestBotProcessesOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long test in short mode")
	}

	c := NewTestController()

	c.AddBot()
	order, err := c.NewOrder("NORMAL")
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	hasOrder := c.Bots[0].GetCurrentOrder() != nil

	if !hasOrder {
		t.Error("Bot should have an order")
	}

	time.Sleep(11 * time.Second)

	if c.CompletedOrders != 1 {
		t.Errorf("Expected 1 completed order, got %d", c.CompletedOrders)
	}
	if c.CompleteOrders[0].ID != order.ID {
		t.Errorf("Expected order ID %d, got %d", order.ID, c.CompleteOrders[0].ID)
	}

	c.StopAllBots()
}

func TestBotReturnsOrderWhenRemoved(t *testing.T) {
	c := NewTestController()

	c.AddBot()
	order, err := c.NewOrder("NORMAL")
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	hasOrder := c.Bots[0].GetCurrentOrder() != nil

	if !hasOrder {
		t.Fatal("Bot does not have order")
	}

	c.RemoveBot()
	time.Sleep(500 * time.Millisecond)

	_, normalCount := c.GetPendingCount()
	if normalCount == 0 {
		t.Error("Order was not returned to queue")
	}

	nextOrder := c.GetNextOrder()
	if nextOrder == nil || nextOrder.ID != order.ID {
		t.Errorf("Expected order %d at front", order.ID)
	}
}

func TestBotDoesNotTakeOrderWhenBusy(t *testing.T) {
	c := NewTestController()

	c.AddBot()
	c.NewOrder("NORMAL")

	time.Sleep(500 * time.Millisecond)

	isBusy := c.Bots[0].IsBusy()

	if !isBusy {
		t.Fatal("Bot is not busy")
	}

	c.NewOrder("NORMAL")
	time.Sleep(500 * time.Millisecond)

	_, normalCount := c.GetPendingCount()
	if normalCount != 1 {
		t.Errorf("Expected 1 order in queue, got %d", normalCount)
	}

	c.StopAllBots()
}

func TestVIPOrderPriority(t *testing.T) {
	c := NewTestController()

	c.NewOrder("NORMAL")
	c.NewOrder("VIP")
	c.AddBot()

	time.Sleep(500 * time.Millisecond)

	bot := c.Bots[0]
	currentOrder := bot.GetCurrentOrder()
	if currentOrder != nil && currentOrder.Type != "VIP" {
		t.Errorf("Bot should pick VIP first, got %s", currentOrder.Type)
	}

	c.StopAllBots()
}
