package main

import (
	"fmt"
	"os"
	"strings"
)

// Database struct
type Database struct {
	records map[string]string
	file    *os.File
	path    string
}

// NewDatabase creates database with explicit path
func NewDatabase(path string) *Database {
	db := &Database{
		records: make(map[string]string),
		path:    path,
	}
	
	// Open file for append/create/read
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Cannot open %s: %v\n", path, err)
		return db
	}
	db.file = file
	
	// Load existing data
	db.load()
	
	return db
}

// load reads existing records
func (db *Database) load() {
	if db.file == nil {
		return
	}
	
	// Seek to beginning
	_, err := db.file.Seek(0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Cannot seek file: %v\n", err)
		return
	}
	
	// Read file content
	data := make([]byte, 65536) // 64KB buffer
	n, err := db.file.Read(data)
	if err != nil && err.Error() != "EOF" {
		fmt.Fprintf(os.Stderr, "Warning: Cannot read file: %v\n", err)
		return
	}
	
	if n > 0 {
		content := string(data[:n])
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) == 3 && parts[0] == "SET" {
				db.records[parts[1]] = parts[2]
			}
		}
	}
}

// Set stores key-value
func (db *Database) Set(key, value string) {
	// Update memory
	db.records[key] = value
	
	// Write to file
	if db.file != nil {
		_, err := db.file.WriteString(fmt.Sprintf("SET %s %s\n", key, value))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Cannot write to file: %v\n", err)
			return
		}
		db.file.Sync()
	}
}

// Get retrieves value
func (db *Database) Get(key string) string {
	return db.records[key]
}

// Close closes file
func (db *Database) Close() {
	if db.file != nil {
		db.file.Close()
	}
}
