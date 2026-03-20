package main

import "testing"

// Test that NewOrder creates unique IDs
func TestNewOrder(t *testing.T) {
	orderID = 1 // reset for test

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
