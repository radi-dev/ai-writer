package ai

import (
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

func GetLlm() *ollama.LLM {
	llm, err := ollama.New(ollama.WithModel("gemma2:2b"))
	if err != nil {
		log.Fatal(err)
	}
	return llm
}

var LLM = GetLlm()
