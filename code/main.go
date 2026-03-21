package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var vipQueue []Order
var normalQueue []Order
var pendingQueue []*Order
var completeOrders []Order

var bots []*Bot
var orderID = 1
var totalVIP = 0
var totalNormal = 0
var completedOrders = 0

func main() {
	// Reset state
	vipQueue = []Order{}
	normalQueue = []Order{}
	bots = []*Bot{}
	completeOrders = []Order{}
	orderID = 1
	totalVIP = 0
	totalNormal = 0
	completedOrders = 0

	log("McDonald's Order Management System - Simulation Started")
	log("System initialized with %d bots", len(bots))

	fmt.Println("Commands: normal, vip, addbot, removebot, status, exit")

	scanner := bufio.NewScanner(os.Stdin)

	// If stdin is not a terminal, run non-interactive automatically
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Non-interactive mode: read all commands from stdin
		for scanner.Scan() {
			cmd := scanner.Text()
			processCommand(cmd)
		}
		return
	}

	// Interactive mode
	for {
		fmt.Print("> ")
		scanner.Scan()
		cmd := scanner.Text()
		processCommand(cmd)
	}
}

func processCommand(cmd string) {
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
		printFinalStatus()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Unknown command")
	}
}

func log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)        // format the message
	now := time.Now().Format("15:04:05")       // get current time in HH:MM:SS
	line := fmt.Sprintf("[%s] %s\n", now, msg) // combine timestamp + message

	fmt.Print(line) // print to console

	// append to result.txt
	f, err := os.OpenFile("../scripts/result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()
	f.WriteString(line)
}

func printFinalStatus() {
	total := totalVIP + totalNormal

	fmt.Println("\nFinal Status:")
	fmt.Printf("- Total Orders Processed: %d (%d VIP, %d Normal)\n",
		total, totalVIP, totalNormal)
	fmt.Printf("- Orders Completed: %d\n", completedOrders)
	fmt.Printf("- Active Bots: %d\n", len(bots))
	fmt.Printf("- Pending Orders: %d\n", len(vipQueue)+len(normalQueue))
}
