package conversation

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tmc/langchaingo/llms"

	"github.com/radi-dev/ai-writer/ai"
	"github.com/radi-dev/ai-writer/database"
	"github.com/radi-dev/ai-writer/database/messages"
	"github.com/radi-dev/ai-writer/services/sse"
)

func WriteLinkedInArticle(ctx context.Context, w http.ResponseWriter, r *http.Request, topic string, length int) string {
	// ctx := r.Context()
	prompt := "Write a random LinkedIn article"
	if topic != "" {
		prompt = "Write a LinkedIn article based on the topic: " + topic
	}
	if length > 0 {
		prompt += fmt.Sprintf(" with a maximum length of %d words", length)
	}

	messages.Create(database.DB, "user", prompt)

	input := strings.TrimSpace(prompt)

	fmt.Print(input)

	response, err := ai.LLM.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}, llms.WithTemperature(0.4), llms.WithMaxLength(length),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {

			fmt.Print(string(chunk))
			sse.SseHandler(ctx, w, r, string(chunk))
			return nil
		}),
	)
	fmt.Fprintf(w, "event: done\ndata: [END]\n\n")
	w.(http.Flusher).Flush()

	// GenerateFromSinglePrompt(
	// 	ctx,
	// 	ai.LLM, prompt, llms.WithMaxLength(100),
	// 	llms.WithTemperature(0.8),
	// 	llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	// 		// fmt.Print(string(chunk))
	// 		sse.SseHandler(w, r, string(chunk))
	// 		return nil
	// 	}),
	// )

	// completion, err := llms.GenerateFromSinglePrompt(
	// 	ctx,
	// 	ai.LLM, prompt, llms.WithMaxLength(100),
	// 	llms.WithTemperature(0.8),
	// 	llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	// 		// fmt.Print(string(chunk))
	// 		sse.SseHandler(w, r, string(chunk))
	// 		return nil
	// 	}),
	// )

	if err != nil {
		log.Fatal(err)
	}

	completion := response.Choices[0].Content

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
