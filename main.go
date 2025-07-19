package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Record struct {
	Prompt, Response string
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "memory.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS conv(id INTEGER PRIMARY KEY, role TEXT, text TEXT)`)
	return db, err
}

func save(db *sql.DB, role, text string) {
	_, _ = db.Exec("INSERT INTO conv(role, text) VALUES (?, ?)", role, text)
}

func history(db *sql.DB) ([]Record, error) {
	rows, err := db.Query("SELECT role, text FROM conv ORDER BY id")
	if err != nil {
		return nil, err
	}
	var recs []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.Prompt, &r.Response); err != nil {
			return nil, err
		}
		recs = append(recs, r)
	}
	return recs, nil
}

func main() {
	ctx := context.Background()
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	llm, err := ollama.New(ollama.WithModel("gemma2:2b"))
	if err != nil {
		log.Fatal(err)
	}

	prompt := `Write a blog outline titled "Why Go is great for LLM apps", with bullet headings`
	save(db, "user", prompt)
	completion, err := llms.GenerateFromSinglePrompt(
		ctx,
		llm, prompt,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	save(db, "assistant", completion)

	// Expand into full post:
	prompt2 := fmt.Sprintf("Expand this outline into a detailed blog post:\n\n%s", completion)
	save(db, "user", prompt2)
	resp2, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt2, llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}))
	if err != nil {
		log.Fatal(err)
	}
	save(db, "assistant", resp2)

	fmt.Println("\nHistory:")
	conv, _ := history(db)
	for _, r := range conv {
		fmt.Printf("[%s] %s\n", r.Prompt, r.Response)
	}

	_ = completion
}
