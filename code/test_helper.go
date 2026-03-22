package main

func resetTestState() {
	mu.Lock()
	defer mu.Unlock()

	// Stop all bots gracefully
	for _, bot := range bots {
		bot.stopOnce.Do(func() {
			close(bot.stopChan)
		})
	}

	vipQueue = []*Order{}
	normalQueue = []*Order{}
	completeOrders = []*Order{}
	bots = []*Bot{}
	orderID = 1
	totalVIP = 0
	totalNormal = 0
	completedOrders = 0
}
