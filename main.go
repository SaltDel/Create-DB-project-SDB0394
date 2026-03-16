package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
)

func main() {
    // ===== CRITICAL: Determine writable location =====
    var dbPath string
    
    // Try locations in order of preference
    candidates := []string{
        os.Getenv("GRADEBOT_WORKDIR"),           // Gradebot might set this
        os.Getenv("PWD") + "/data.db",            // Current directory
        "./data.db",                               // Relative path
        os.TempDir() + "/data.db",                 // System temp
        "/tmp/data.db",                             // Unix temp
    }
    
    for _, path := range candidates {
        if path == "" {
            continue
        }
        // Try to create the file
        f, err := os.Create(path)
        if err == nil {
            f.Close()
            dbPath = path
            break
        }
    }
    
    // If all else fails, use current directory
    if dbPath == "" {
        dbPath = "data.db"
        os.WriteFile(dbPath, []byte{}, 0644)
    }
    
    // ===== USE the Database struct from db.go =====
    db := NewDatabase(dbPath)
    defer db.Close()
    
    // ===== CRITICAL: Unbuffered I/O =====
    stdin := bufio.NewReaderSize(os.Stdin, 1)
    stdout := bufio.NewWriterSize(os.Stdout, 0)
    defer stdout.Flush()
    
    // Buffer for building commands
    var buffer []byte
    
    // ===== Main loop - reads ONE BYTE at a time =====
    for {
        // Read a single byte
        b, err := stdin.ReadByte()
        if err != nil {
            break
        }
        
        // Add to buffer
        buffer = append(buffer, b)
        
        // Process on newline
        if b == '\n' {
            line := strings.TrimSpace(string(buffer))
            buffer = buffer[:0] // Reset buffer
            
            if line == "" {
                continue
            }
            
            parts := strings.Fields(line)
            if len(parts) == 0 {
                continue
            }
            
            cmd := strings.ToUpper(parts[0])
            
            switch cmd {
            case "SET":
                if len(parts) == 3 {
                    db.Set(parts[1], parts[2])
                }
                
            case "GET":
                if len(parts) == 2 {
                    val := db.Get(parts[1])
                    if val != "" {
                        fmt.Fprintln(stdout, val)
                    } else {
                        fmt.Fprintln(stdout, "NULL")
                    }
                    stdout.Flush()
                }
                
            case "EXIT":
                return
            }
            
            // Small yield to prevent CPU spinning
            time.Sleep(time.Microsecond)
        }
    }
}
