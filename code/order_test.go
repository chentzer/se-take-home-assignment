package main

import (
	"testing"
)

func TestCreateNormalOrder(t *testing.T) {
	resetTestState()

	order, err := NewOrder("NORMAL")

	if err != nil {
		t.Fatalf("Order creation failed: %v", err)
	}

	if order.Type != "NORMAL" {
		t.Errorf("Expected NORMAL, got %s", order.Type)
	}

	mu.Lock()
	if totalNormal != 1 {
		t.Errorf("Expected totalNormal=1, got %d", totalNormal)
	}
	if len(normalQueue) != 1 {
		t.Errorf("Expected normalQueue length 1, got %d", len(normalQueue))
	}
	mu.Unlock()
}

func TestCreateVIPOrder(t *testing.T) {
	resetTestState()

	order, err := NewOrder("VIP")

	if err != nil {
		t.Fatalf("Order creation failed: %v", err)
	}

	if order.Type != "VIP" {
		t.Errorf("Expected VIP, got %s", order.Type)
	}

	mu.Lock()
	if totalVIP != 1 {
		t.Errorf("Expected totalVIP=1, got %d", totalVIP)
	}
	if len(vipQueue) != 1 {
		t.Errorf("Expected vipQueue length 1, got %d", len(vipQueue))
	}
	mu.Unlock()
}

func TestInvalidOrderType(t *testing.T) {
	resetTestState()

	order, err := NewOrder("INVALID")

	if err == nil {
		t.Error("Expected error for invalid order type")
	}
	if order != nil {
		t.Error("Expected nil order for invalid type")
	}
}

func TestOrderIDIncrements(t *testing.T) {
	resetTestState()

	order1, _ := NewOrder("NORMAL")
	order2, _ := NewOrder("NORMAL")

	if order2.ID != order1.ID+1 {
		t.Errorf("Order IDs not sequential: %d then %d", order1.ID, order2.ID)
	}
}

func TestGetNextOrderFromEmptyQueue(t *testing.T) {
	resetTestState()

	order := getNextOrder()

	if order != nil {
		t.Error("Expected nil from empty queue")
	}
}

func TestGetNextOrderPriority(t *testing.T) {
	resetTestState()

	NewOrder("NORMAL")
	NewOrder("VIP")

	order := getNextOrder()

	if order.Type != "VIP" {
		t.Errorf("Expected VIP first, got %s", order.Type)
	}
}

func TestGetNextOrderNormalOnly(t *testing.T) {
	resetTestState()

	NewOrder("NORMAL")
	NewOrder("NORMAL")

	order1 := getNextOrder()
	order2 := getNextOrder()

	if order1.Type != "NORMAL" {
		t.Errorf("Expected NORMAL, got %s", order1.Type)
	}
	if order2.Type != "NORMAL" {
		t.Errorf("Expected NORMAL, got %s", order2.Type)
	}
}

func TestGetNextOrderVIPOnly(t *testing.T) {
	resetTestState()

	NewOrder("VIP")
	NewOrder("VIP")

	order1 := getNextOrder()
	order2 := getNextOrder()

	if order1.Type != "VIP" {
		t.Errorf("Expected VIP, got %s", order1.Type)
	}
	if order2.Type != "VIP" {
		t.Errorf("Expected VIP, got %s", order2.Type)
	}
}

func TestFIFOWithinSameType(t *testing.T) {
	resetTestState()

	order1, _ := NewOrder("NORMAL")
	order2, _ := NewOrder("NORMAL")

	first := getNextOrder()
	second := getNextOrder()

	if first.ID != order1.ID {
		t.Errorf("Expected first order ID %d, got %d", order1.ID, first.ID)
	}
	if second.ID != order2.ID {
		t.Errorf("Expected second order ID %d, got %d", order2.ID, second.ID)
	}
}

func TestMixedQueueOrdering(t *testing.T) {
	resetTestState()

	NewOrder("VIP")
	NewOrder("NORMAL")
	NewOrder("VIP")
	NewOrder("NORMAL")

	order1 := getNextOrder()
	order2 := getNextOrder()
	order3 := getNextOrder()
	order4 := getNextOrder()

	if order1.Type != "VIP" {
		t.Error("First order should be VIP")
	}
	if order2.Type != "VIP" {
		t.Error("Second order should be VIP")
	}
	if order3.Type != "NORMAL" {
		t.Error("Third order should be NORMAL")
	}
	if order4.Type != "NORMAL" {
		t.Error("Fourth order should be NORMAL")
	}
}
