package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var mu sync.Mutex

var vipQueue []*Order
var normalQueue []*Order
var completeOrders []*Order

var bots []*Bot

var orderID = 1
var totalVIP = 0
var totalNormal = 0
var completedOrders = 0

var shutdownChan = make(chan struct{})
var logFile *os.File
var logMutex sync.Mutex

func main() {
	setupGracefulShutdown()

	scanner := bufio.NewScanner(os.Stdin)

	log("McDonald's Order Management System - Simulation Started")
	log("System initialized with %d bots", len(bots))

	printHelp()

	for {
		select {
		case <-shutdownChan:
			gracefulShutdown()
			return
		default:
		}

		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		cmd := scanner.Text()

		switch cmd {
		case "normal":
			addNormalOrder()
		case "vip":
			addVIPOrder()
		case "addbot":
			addBot()
		case "removebot":
			removeBot()
		case "status":
			printStatus()
		case "help":
			printHelp()
		case "exit":
			log("Exit command received")
			gracefulShutdown()
			return
		default:
			if cmd != "" {
				fmt.Println("Unknown command. Type 'help' for available commands.")
			}
		}
	}
}

func printHelp() {
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  normal    - Add a normal order")
	fmt.Println("  vip       - Add a VIP order")
	fmt.Println("  addbot    - Add a new cooking bot")
	fmt.Println("  removebot - Remove the newest cooking bot")
	fmt.Println("  status    - Show system status")
	fmt.Println("  help      - Show this help message")
	fmt.Println("  exit      - Shutdown the system")
	fmt.Println()
}

func setupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log("\nReceived shutdown signal")
		close(shutdownChan)
	}()
}

func gracefulShutdown() {
	log("Initiating graceful shutdown...")

	mu.Lock()
	botsToStop := make([]*Bot, len(bots))
	copy(botsToStop, bots)
	bots = []*Bot{}
	mu.Unlock()

	for _, bot := range botsToStop {
		bot.stopOnce.Do(func() {
			close(bot.stopChan)
		})
	}

	time.Sleep(2 * time.Second)
	log("System shutdown complete")
}

func addNormalOrder() {
	order, err := NewOrder("NORMAL")
	if err != nil {
		log("Error creating order: %v", err)
		return
	}
	log("Created Normal Order #%d - Status: PENDING", order.ID)
}

func addVIPOrder() {
	order, err := NewOrder("VIP")
	if err != nil {
		log("Error creating order: %v", err)
		return
	}
	log("Created VIP Order #%d - Status: PENDING", order.ID)
}

func printStatus() {
	mu.Lock()
	defer mu.Unlock()

	total := totalVIP + totalNormal
	pending := len(vipQueue) + len(normalQueue)

	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║     MCDONALD'S ORDER STATUS           ║")
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Printf("║ Total Orders Created:  %-20d ║\n", total)
	fmt.Printf("║   ├─ VIP Orders:        %-20d ║\n", totalVIP)
	fmt.Printf("║   └─ Normal Orders:     %-20d ║\n", totalNormal)
	fmt.Printf("║ Orders Completed:       %-20d ║\n", completedOrders)
	fmt.Printf("║ Orders Pending:         %-20d ║\n", pending)
	fmt.Printf("║   ├─ VIP Queue:         %-20d ║\n", len(vipQueue))
	fmt.Printf("║   └─ Normal Queue:      %-20d ║\n", len(normalQueue))
	fmt.Printf("║ Active Bots:            %-20d ║\n", len(bots))
	if total > 0 {
		fmt.Printf("║ Completion Rate:        %-20.1f%% ║\n",
			float64(completedOrders)/float64(total)*100)
	}
	fmt.Println("╚════════════════════════════════════════╝")
}

func init() {
	var err error
	logFile, err = os.OpenFile("../scripts/result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logFile = nil
	}
}

func log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	now := time.Now().Format("15:04:05")
	line := fmt.Sprintf("[%s] %s\n", now, msg)

	fmt.Print(line)

	if logFile != nil {
		logMutex.Lock()
		defer logMutex.Unlock()
		logFile.WriteString(line)
	}
}
