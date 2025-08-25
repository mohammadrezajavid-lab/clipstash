package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.design/x/clipboard"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	if err := clipboard.Init(); err != nil {
		log.Fatalf("Failed to initialize clipboard: %v", err)
	}

	dbPath, err := getDBPath()
	if err != nil {
		log.Fatalf("Could not determine DB path: %v", err)
	}

	db := InitDB(dbPath)
	defer db.Close()

	if len(os.Args) == 1 {
		runAgent(db)
	} else {
		handleCliCommand(db, os.Args[1:])
	}
}

func runAgent(db *sql.DB) {
	log.Println("Agent is running and watching clipboard...")
	ch := clipboard.Watch(context.Background(), clipboard.FmtText)

	var lastContent string

	err := db.QueryRow("SELECT content FROM history ORDER BY id DESC LIMIT 1").Scan(&lastContent)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Could not get last history item:", err)
	}

	for data := range ch {
		content := string(data)
		if content != lastContent && content != "" {
			InsertData(content, db)
			lastContent = content
		}
	}
}

func InsertData(content string, db *sql.DB) {
	log.Println("New item copied: ", content)
	stmt, err := db.Prepare("INSERT INTO history(content, created_at) VALUES(?, ?)")
	if err != nil {
		log.Println("Error preparing statement: ", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(content, time.Now())
	if err != nil {
		log.Println("Error executing statement: ", err)
	}
}

func InitDB(filepath string) *sql.DB {

	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		log.Println("Failed to set WAL mode:", err)
	}
	if _, err := db.Exec(`PRAGMA synchronous=NORMAL;`); err != nil {
		log.Println("Failed to set synchronous=NORMAL:", err)
	}
	if _, err := db.Exec(`PRAGMA busy_timeout=5000;`); err != nil {
		log.Println("Failed to set busy_timeout:", err)
	}
	db.SetMaxOpenConns(1)

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS history (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        created_at DATETIME NOT NULL
    );`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt)
	}
	return db
}

func handleCliCommand(db *sql.DB, args []string) {
	if len(args) == 0 {
		return
	}
	command := args[0]
	switch command {
	case "list":
		fmt.Println("Last 10 clipboard items:")
		rows, err := db.Query("SELECT id, content FROM history ORDER BY id DESC LIMIT 10")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var content string
			if err := rows.Scan(&id, &content); err != nil {
				log.Fatal(err)
			}
			if len(content) > 80 {
				content = content[:80] + "..."
			}
			fmt.Printf("%d: %s\n", id, content)
		}

	case "get":
		if len(args) < 2 {
			log.Fatal("Error: Please provide an ID. Usage: clipstash get <id>")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("Error: Invalid ID '%s'. It must be a number.", args[1])
		}

		var content string
		err = db.QueryRow("SELECT content FROM history WHERE id = ?", id).Scan(&content)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Fatalf("No item found with ID: %d", id)
			}
			log.Fatal(err)
		}

		clipboard.Write(clipboard.FmtText, []byte(content))
		fmt.Printf("Copied item #%d to clipboard.\n", id)

	case "search":
		if len(args) < 2 {
			log.Fatal("Error: Please provide a search term. Usage: clipstash search <term>")
		}
		term := args[1]
		query := "SELECT id, content FROM history WHERE content LIKE ? ORDER BY id DESC LIMIT 20"
		rows, err := db.Query(query, "%"+term+"%")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		fmt.Printf("Search results for \"%s\":\n", term)
		found := false
		for rows.Next() {
			found = true
			var id int
			var content string
			if err := rows.Scan(&id, &content); err != nil {
				log.Fatal(err)
			}
			if len(content) > 80 {
				content = content[:80] + "..."
			}
			fmt.Printf("%d: %s\n", id, content)
		}
		if !found {
			fmt.Println("No items found.")
		}

	case "clear":
		res, eErr := db.Exec("DELETE FROM history")
		if eErr != nil {
			log.Fatal(eErr)
		}

		if _, err := db.Exec("VACUUM"); err != nil {
			log.Println("VACUUM failed:", err)
		}

		n, _ := res.RowsAffected()
		fmt.Printf("Cleared %d items.\n", n)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: list, get, search")
	}
}

func getDBPath() (string, error) {

	if v := os.Getenv("CLIPSTASH_DB_PATH"); v != "" {
		if err := os.MkdirAll(filepath.Dir(v), 0700); err != nil {
			return "", fmt.Errorf("could not create parent directory for CLIPSTASH_DB_PATH: %w", err)
		}
		return v, nil
	}

	// - Linux:   ~/.config
	// - macOS:   ~/Library/Application Support
	// - Windows: %APPDATA%
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not get user config directory: %w", err)
	}

	toolConfigDir := filepath.Join(configDir, "clipstash")

	if err := os.MkdirAll(toolConfigDir, 0700); err != nil {
		return "", fmt.Errorf("could not create tool config directory: %w", err)
	}

	return filepath.Join(toolConfigDir, "clipstash.db"), nil
}
