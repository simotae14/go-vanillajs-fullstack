package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection string for your remote PostgreSQL database
	connStr := "postgres://username:password@remote-host:port/database?sslmode=disable"

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test the connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Read the SQL file
	sqlFilePath := "database-dump.sql" // Adjust this to your .sql file path
	sqlContent, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatal("Failed to read SQL file:", err)
	}

	// Split the SQL content into individual statements
	// This is a basic split; it assumes statements end with semicolons
	statements := strings.Split(string(sqlContent), ";\n")

	// Execute each statement
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue // Skip empty statements
		}

		// Remove comment lines from the statement
		lines := strings.Split(stmt, "\n")
		var cleanedLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "--") {
				continue // Skip comment lines
			}
			cleanedLines = append(cleanedLines, trimmed)
		}

		cleanedStmt := strings.Join(cleanedLines, " ")
		if cleanedStmt == "" {
			continue // Skip if the statement is empty after cleaning
		}

		// Execute the cleaned statement
		_, err := db.Exec(cleanedStmt)
		if err != nil {
			log.Printf("Failed to execute statement: %v\nStatement: %s\n", err, cleanedStmt)
			// Continue with next statement even if one fails
			return
		}
		fmt.Printf("Executed: %s\n", cleanedStmt[:min(50, len(cleanedStmt))]+"...") // Log first 50 chars
	}

	fmt.Println("SQL script execution completed.")
}

// Helper function to get min of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
