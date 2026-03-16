package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Initialize database
	db := NewDatabase()
	defer db.Close()

	// Create scanner for reading input
	scanner := bufio.NewScanner(os.Stdin)
	
	// Process commands
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		// Parse command
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		
		command := strings.ToUpper(parts[0])
		
		switch command {
		case "SET":
			if len(parts) != 3 {
				continue // Invalid SET command
			}
			key := parts[1]
			value := parts[2]
			db.Set(key, value)
			
		case "GET":
			if len(parts) != 2 {
				continue // Invalid GET command
			}
			key := parts[1]
			if value, ok := db.Get(key); ok {
				fmt.Println(value)
			} else {
				fmt.Println("NULL")
			}
			
		case "EXIT":
			return
			
		default:
			// Unknown command - ignore
			continue
		}
	}
	
	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}