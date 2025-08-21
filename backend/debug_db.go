package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/junk_journal?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check what elements exist
	fmt.Println("=== Elements in database ===")
	rows, err := db.Query("SELECT id, page_id, kind, visible, locked FROM elements")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, pageId, kind string
		var visible, locked *bool
		err := rows.Scan(&id, &pageId, &kind, &visible, &locked)
		if err != nil {
			log.Fatal(err)
		}

		visibleStr := "NULL"
		if visible != nil {
			if *visible {
				visibleStr = "true"
			} else {
				visibleStr = "false"
			}
		}

		lockedStr := "NULL"
		if locked != nil {
			if *locked {
				lockedStr = "true"
			} else {
				lockedStr = "false"
			}
		}

		fmt.Printf("ID: %s | Page: %s | Kind: %s | Visible: %s | Locked: %s\n", id, pageId, kind, visibleStr, lockedStr)
	}

	// Check pages
	fmt.Println("\n=== Pages in database ===")
	rows2, err := db.Query("SELECT id, board_id, title FROM pages")
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var id, boardId, title string
		err := rows2.Scan(&id, &boardId, &title)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Page ID: %s | Board: %s | Title: %s\n", id, boardId, title)
	}
}
