package sse

import (
	"context"
	"fmt"
	"net/http"
)

func SseHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, chunk string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	select {
	case <-ctx.Done():
		fmt.Println("Client disconnected")
		return
	default:
		// Write the chunk as an SSE event
		// _, err := fmt.Fprint(w, chunk)
		fmt.Fprintf(w, "event: message\n")
		fmt.Fprintf(w, "data: %s\n\n", chunk)

		flusher.Flush()
	}
}
