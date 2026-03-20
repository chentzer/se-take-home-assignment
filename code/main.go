package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var vipQueue []Order
var normalQueue []Order
var completeOrders []Order

var bots []*Bot
var orderID = 1

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Commands: normal, vip, addbot, removebot, status, exit")

	for {
		fmt.Print("> ")
		scanner.Scan()
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
		case "exit":
			return
		}
	}
}

func log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)        // format the message
	now := time.Now().Format("15:04:05")       // get current time in HH:MM:SS
	line := fmt.Sprintf("[%s] %s\n", now, msg) // combine timestamp + message

	fmt.Print(line) // print to console

	// append to result.txt
	f, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()
	f.WriteString(line)
}

func printStatus() {
	fmt.Println("VIP Queue:", vipQueue)
	fmt.Println("Normal Queue:", normalQueue)
	fmt.Println("Completed:", completeOrders)
	fmt.Println("Bots:", len(bots))
}
