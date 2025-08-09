package conversation

import (
	"context"
	"fmt"
	"log"

	"github.com/radi-dev/ai-writer/ai"
	"github.com/radi-dev/ai-writer/database"
	"github.com/radi-dev/ai-writer/database/messages"
	"github.com/tmc/langchaingo/llms"
)

func WriteLinkedInArticle(ctx context.Context, prompt string) string {
	// sysPrompt := `Write a blog outline titled "Why Go is great for LLM apps", with bullet headings`
	if prompt == "" {
		prompt = "Write a LinkedIn article"
	}

	messages.Create(database.DB, "user", prompt)

	completion, err := llms.GenerateFromSinglePrompt(
		ctx,
		ai.LLM, prompt,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	messages.Create(database.DB, "assistant", completion)
	return completion
}

func RefineArticle(ctx context.Context, article string) string {

	prompt := fmt.Sprintf("Humanize the following article so it doesn't feel like it was written by AI:\n\n%s", article)
	messages.Create(database.DB, "user", prompt)
	completion, err := llms.GenerateFromSinglePrompt(ctx, ai.LLM, prompt, llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}))
	if err != nil {
		log.Fatal(err)
	}
	messages.Create(database.DB, "assistant", completion)

	fmt.Println("\nHistory:")
	conv, _ := messages.FetchHistory(database.DB)
	for _, r := range conv {
		fmt.Printf("[%s] %s\n", r.Prompt, r.Response)
	}

	return completion
}
