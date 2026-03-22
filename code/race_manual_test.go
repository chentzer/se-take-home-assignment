package main

import (
	"sync"
	"testing"
	"time"
)

func TestConcurrentOperations(t *testing.T) {
	t.Run("ConcurrentOrderCreation", func(t *testing.T) {
		for iter := 0; iter < 3; iter++ {
			resetTestState()

			var wg sync.WaitGroup
			orderCount := 50

			for i := 0; i < orderCount; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					if idx%2 == 0 {
						NewOrder("NORMAL")
					} else {
						NewOrder("VIP")
					}
				}(i)
			}

			wg.Wait()

			mu.Lock()
			total := totalNormal + totalVIP
			mu.Unlock()

			if total != orderCount {
				t.Errorf("Iteration %d: Expected %d orders, got %d", iter, orderCount, total)
			}
		}
	})

	t.Run("ConcurrentBotAddRemove", func(t *testing.T) {
		for iter := 0; iter < 3; iter++ {
			resetTestState()

			// Add some orders
			for i := 0; i < 10; i++ {
				NewOrder("NORMAL")
			}

			var wg sync.WaitGroup

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					addBot()
				}()

				wg.Add(1)
				go func() {
					defer wg.Done()
					time.Sleep(50 * time.Millisecond)
					removeBot()
				}()
			}

			wg.Wait()
			time.Sleep(500 * time.Millisecond)

			// Just verify no panic
			mu.Lock()
			_ = len(bots)
			_ = completedOrders
			mu.Unlock()
		}
	})

	t.Run("OrderIDUniqueness", func(t *testing.T) {
		for iter := 0; iter < 3; iter++ {
			resetTestState()

			var wg sync.WaitGroup
			ids := make(map[int]bool)
			var muIDs sync.Mutex
			duplicates := 0

			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					order, _ := NewOrder("NORMAL")
					if order != nil {
						muIDs.Lock()
						if ids[order.ID] {
							duplicates++
						}
						ids[order.ID] = true
						muIDs.Unlock()
					}
				}()
			}

			wg.Wait()

			if duplicates > 0 {
				t.Errorf("Iteration %d: Found %d duplicate order IDs", iter, duplicates)
			}

			t.Logf("Iteration %d: Created %d unique orders", iter, len(ids))
		}
	})

	t.Run("OrderProcessingWithMultipleBots", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping long test in short mode")
		}

		for iter := 0; iter < 2; iter++ {
			resetTestState()

			// Add bots
			botCount := 3
			for i := 0; i < botCount; i++ {
				addBot()
			}

			time.Sleep(200 * time.Millisecond)

			// Add orders
			orderCount := 6 // 6 orders with 3 bots = 20 seconds total (2 rounds)
			for i := 0; i < orderCount; i++ {
				NewOrder("NORMAL")
			}

			// Wait enough time for all orders to process
			// 6 orders with 3 bots = 2 rounds of 10 seconds = 20 seconds
			// Add 2 seconds buffer
			processingTime := 22 * time.Second
			time.Sleep(processingTime)

			mu.Lock()
			completed := completedOrders
			pending := len(normalQueue) + len(vipQueue)
			total := totalNormal + totalVIP
			mu.Unlock()

			t.Logf("Iteration %d: Total=%d, Completed=%d, Pending=%d", iter, total, completed, pending)

			// All orders should be completed
			if completed != total {
				t.Errorf("Iteration %d: Not all orders completed. Completed=%d, Total=%d, Pending=%d",
					iter, completed, total, pending)
			}

			// No orders should be pending
			if pending != 0 {
				t.Errorf("Iteration %d: Expected 0 pending orders, got %d", iter, pending)
			}
		}
	})

	t.Run("VIPPriorityWithMultipleOrders", func(t *testing.T) {
		resetTestState()

		// Add orders in specific order: Normal, VIP, Normal, VIP
		NewOrder("NORMAL") // Order 1
		NewOrder("VIP")    // Order 2
		NewOrder("NORMAL") // Order 3
		NewOrder("VIP")    // Order 4

		// Add one bot
		addBot()

		// Give bot time to pick up first order
		time.Sleep(500 * time.Millisecond)

		mu.Lock()
		firstOrder := bots[0].CurrentOrder
		mu.Unlock()

		// First order picked should be VIP (order 2), not normal (order 1)
		if firstOrder == nil {
			t.Fatal("Bot did not pick up any order")
		}

		if firstOrder.Type != "VIP" {
			t.Errorf("Expected VIP order first, got %s order #%d", firstOrder.Type, firstOrder.ID)
		}

		// Wait for first order to complete
		time.Sleep(11 * time.Second)

		mu.Lock()
		secondOrder := bots[0].CurrentOrder
		mu.Unlock()

		// Second order should be the other VIP (order 4)
		if secondOrder == nil {
			t.Fatal("Bot did not pick up second order")
		}

		if secondOrder.Type != "VIP" {
			t.Errorf("Expected VIP order second, got %s order #%d", secondOrder.Type, secondOrder.ID)
		}

		t.Log("VIP priority test passed - VIP orders processed before normal orders")
	})
}
