package main

import (
	"bufio"
	"flag"
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

// Command line flags
var demoMode bool
var outputFile string

func main() {
	flag.BoolVar(&demoMode, "demo", false, "Run in demo mode with predefined commands")
	flag.StringVar(&outputFile, "output", "", "Output file path for logs (default: ../scripts/result.txt)")
	flag.Parse()

	initLogFile()
	setupGracefulShutdown()

	if demoMode {
		runDemoMode()
		return
	}

	runInteractiveMode()
}

func runInteractiveMode() {
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
		processCommand(cmd)
		if cmd == "exit" {
			return
		}
	}
}

func runDemoMode() {
	log("McDonald's Order Management System - Demo Mode Started")
	log("System initialized with %d bots", len(bots))

	// Demo sequence that showcases all functionality
	commands := []struct {
		cmd   string
		delay time.Duration
	}{
		{"normal", 100 * time.Millisecond},
		{"vip", 100 * time.Millisecond},
		{"normal", 100 * time.Millisecond},
		{"addbot", 100 * time.Millisecond},
		{"addbot", 100 * time.Millisecond},
		{"status", 2 * time.Second},
		{"status", 12 * time.Second}, // Wait for orders to complete
		{"removebot", 100 * time.Millisecond},
		{"status", 100 * time.Millisecond},
		{"exit", 0},
	}

	for _, c := range commands {
		log("Command: %s", c.cmd)
		processCommand(c.cmd)
		if c.cmd == "exit" {
			return
		}
		time.Sleep(c.delay)
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
		printStatus()
	case "help":
		printHelp()
	case "exit":
		log("Exit command received")
		gracefulShutdown()
	default:
		if cmd != "" {
			fmt.Println("Unknown command. Type 'help' for available commands.")
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

func initLogFile() {
	var logPath string

	if outputFile != "" {
		// Use the specified output file path
		logPath = outputFile
	} else {
		// Default path: try to find scripts directory relative to executable or current directory
		execPath, err := os.Executable()
		if err == nil {
			// Try relative to executable (for built binary)
			scriptsDir := execPath + "/../scripts"
			if _, err := os.Stat(scriptsDir); err == nil {
				logPath = scriptsDir + "/result.txt"
			}
		}

		if logPath == "" {
			// Fallback: relative to current working directory
			logPath = "../scripts/result.txt"
		}
	}

	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		// Try current directory as last resort
		logFile, err = os.OpenFile("result.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not open log file: %v\n", err)
			logFile = nil
		}
	}
}

func log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	now := time.Now().Format("15:04:05")
	line := fmt.Sprintf("[%s] %s\n", now, msg)

	fmt.Print(line)

	if logFile != nil {
		logMutex.Lock()
		_, err := logFile.WriteString(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to write to log file: %v\n", err)
		}
		logFile.Sync() // Ensure data is flushed to disk
		logMutex.Unlock()
	}
}
