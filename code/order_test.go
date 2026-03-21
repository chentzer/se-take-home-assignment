package main

import "testing"

// Test that NewOrder creates unique IDs
func TestNewOrder(t *testing.T) {
	orderID = 1 // reset

	o1 := NewOrder("NORMAL")
	o2 := NewOrder("VIP")

	if o1.ID != 1 {
		t.Errorf("Expected ID 1, got %d", o1.ID)
	}
	if o2.ID != 2 {
		t.Errorf("Expected ID 2, got %d", o2.ID)
	}
	if o1.Type != "NORMAL" {
		t.Errorf("Expected type NORMAL, got %s", o1.Type)
	}
	if o2.Type != "VIP" {
		t.Errorf("Expected type VIP, got %s", o2.Type)
	}
}

func TestVIPPriority(t *testing.T) {
	// Reset state
	orderID = 1
	vipQueue = []Order{}
	normalQueue = []Order{}

	n := NewOrder("NORMAL")
	v := NewOrder("VIP")

	normalQueue = append(normalQueue, n)
	vipQueue = append(vipQueue, v)

	order := getNextOrder()

	if order == nil {
		t.Fatalf("Expected order, got nil")
	}
	if order.Type != "VIP" {
		t.Errorf("Expected VIP order first, got %s", order.Type)
	}
}

func TestFIFOWithinSameType(t *testing.T) {
	// Reset state
	orderID = 1
	vipQueue = []Order{}
	normalQueue = []Order{}

	o1 := NewOrder("VIP")
	o2 := NewOrder("VIP")

	vipQueue = append(vipQueue, o1, o2)

	first := getNextOrder()
	second := getNextOrder()

	if first == nil || second == nil {
		t.Fatalf("Expected two orders, got nil")
	}

	if first.ID != o1.ID {
		t.Errorf("Expected first VIP order ID %d, got %d", o1.ID, first.ID)
	}
	if second.ID != o2.ID {
		t.Errorf("Expected second VIP order ID %d, got %d", o2.ID, second.ID)
	}
}

// TestMixedOrders simulates multiple normal and VIP orders
// and validates VIP priority and FIFO behavior
func TestMixedOrders(t *testing.T) {
	// Reset global state
	orderID = 1
	vipQueue = []Order{}
	normalQueue = []Order{}
	completeOrders = []Order{}
	totalVIP = 0
	totalNormal = 0
	completedOrders = 0

	// Add orders in realistic sequence
	n1 := NewOrder("NORMAL") // #1
	n2 := NewOrder("NORMAL") // #2
	v1 := NewOrder("VIP")    // #3
	n3 := NewOrder("NORMAL") // #4
	v2 := NewOrder("VIP")    // #5

	// Add to queues
	normalQueue = append(normalQueue, n1, n2, n3)
	vipQueue = append(vipQueue, v1, v2)

	// Pick orders in sequence using getNextOrder()
	expectedSequence := []int{v1.ID, v2.ID, n1.ID, n2.ID, n3.ID}
	actualSequence := []int{}

	for i := 0; i < 5; i++ {
		o := getNextOrder()
		if o == nil {
			t.Fatalf("Expected order at position %d, got nil", i+1)
		}
		actualSequence = append(actualSequence, o.ID)
	}

	for i := 0; i < len(expectedSequence); i++ {
		if actualSequence[i] != expectedSequence[i] {
			t.Errorf("Expected order ID %d at position %d, got %d",
				expectedSequence[i], i+1, actualSequence[i])
		}
	}
}
